package middleware

import (
	"fmt"

	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type TokenBlacklist struct {
	blacklist map[string]int64 // Maps access_uuid to expiration time
	mutex     sync.RWMutex
}

var tokenBlacklist = &TokenBlacklist{

	blacklist: make(map[string]int64),
}

func (tb *TokenBlacklist) IsBlacklisted(accessUuid string) bool {

	tb.mutex.RLock()
	defer tb.mutex.RUnlock()

	expTime, exist := tb.blacklist[accessUuid]
	if !exist {
		return false
	}

	if expTime < time.Now().Unix() {
		tb.mutex.RUnlock()
		tb.mutex.Lock()
		delete(tb.blacklist, accessUuid)
		tb.mutex.Unlock()
		tb.mutex.RLocker()
	}
	return true
}

func GenerateTokens(userId primitive.ObjectID) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 24).Unix()
	td.AccessUuid = primitive.NewObjectID().Hex()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = primitive.NewObjectID().Hex()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId.Hex()
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, err
	}

	rtCliams := jwt.MapClaims{}
	rtCliams["refresh_uuid"] = td.RefreshUuid
	rtCliams["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))

	if err != nil {
		return nil, err

	}
	return td, nil
}

func (tb *TokenBlacklist) AddToBlacklist(accessUuid string, expiresAt int64) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.blacklist[accessUuid] = expiresAt

}

func (tb *TokenBlacklist) CleanupBlacklist() {
	tb.mutex.Lock()

	defer tb.mutex.Lock()

	now := time.Now().Unix()
	for uuid, expTime := range tb.blacklist {
		if expTime < now {
			delete(tb.blacklist, uuid)
		}
	}
}

func StartCleanupRoutine() {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			tokenBlacklist.CleanupBlacklist()
		}
	}()

}

func BlacklistToken(accessUuid string, expiresAt int64) {
	tokenBlacklist.AddToBlacklist(accessUuid, expiresAt)
}

func IsTokenBlaclisted(accessUuid string) bool {
	return tokenBlacklist.IsBlacklisted(accessUuid)
}

func ExtractTokenMetadata(c *fiber.Ctx) (*AccessDetails, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing access_uuid in token")
	}

	// Check if token is blacklisted
	if IsTokenBlaclisted(accessUuid) {
		return nil, fmt.Errorf("token has been revoked")
	}

	userIdStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing user_id in token")
	}

	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format in token")
	}

	return &AccessDetails{
		AccessUuid: accessUuid,
		UserId:     userId,
	}, nil
}

type AccessDetails struct {
	AccessUuid string
	UserId     primitive.ObjectID
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := ExtractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing in method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil

	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		return bearerToken[7:]
	}
	return ""
}
func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := ExtractTokenMetadata(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized",
				"error":   err.Error(),
			})
		}
		return c.Next()
	}
}
