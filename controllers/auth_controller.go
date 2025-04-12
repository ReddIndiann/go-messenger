package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ReddIndiann/go-messanger/database"
	"github.com/ReddIndiann/go-messanger/helpers"
	"github.com/ReddIndiann/go-messanger/middleware"
	"github.com/ReddIndiann/go-messanger/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser() fiber.Handler {
	// RegisterUser handles user registration
	return func(c *fiber.Ctx) error {
		var user model.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// fmt.Printf("Parsed user: %+v\n", user)
		var requestData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		password, passwordExists := requestData["password"].(string)
		role, roleExists := requestData["role"].(string)

		if user.Name == "" || user.Email == "" || !passwordExists || password == "" || user.Phone == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "All fields are required",
			})
		}

		// Set default role if not provided
		if !roleExists || role == "" {
			role = string(model.RoleUser)
		}

		// Validate role
		if role != string(model.RoleUser) && role != string(model.RoleAdmin) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid role. Must be either 'user' or 'admin'",
			})
		}

		collection := database.GetCollection("users")
		existingUser := collection.FindOne(context.Background(), bson.M{"email": user.Email})
		if existingUser.Err() == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"error":  "User already exists",
				"email":  user.Email,
			})
		}

		// Hash the password
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error hashing password",
			})
		}
		hashedPassword := string(hashedPasswordBytes)

		newUser := model.User{
			ID:        primitive.NewObjectID(),
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Password:  hashedPassword,
			Verified:  false,
			Role:      role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := collection.InsertOne(c.Context(), newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating user",
			})
		}
		go helpers.SendOTP(user.Email, "register", user.Name)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "User registered successfully. Check your email to verify your account.",
			"userId":  result.InsertedID,
			"info":    newUser,
		})

	}
}
func VerifyMail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Request struct {
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}

		var request Request
		var user model.User

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
			})
		}
		storedOTP, found := helpers.GetOTP(request.Email)
		if !found {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "OTP expired or not found",
			})
		}

		if storedOTP != request.OTP {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid OTP",
			})
		}

		if user.Verified {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User already verified",
			})
		}
		collection := database.GetCollection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		filter := bson.M{"email": request.Email}
		update := bson.M{"$set": bson.M{"verified": true}}

		result := collection.FindOneAndUpdate(ctx, filter, update)

		if result.Err() != nil {
			if result.Err() == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"message": "User not found",
				})
			}
			log.Println("Error updating user:", result.Err())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to verify user",
			})
		}
		helpers.DeleteOTP(request.Email)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "User verified successfully",
		})
	}

}

func ResendOTP() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type ResendRequest struct {
			Email string `json:"email"`
		}

		var request ResendRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		collection := database.GetCollection("users")

		var user model.User
		err := collection.FindOne(context.Background(), bson.M{"email": request.Email}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		if user.Verified {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User already verified",
			})
		}

		newCode := helpers.GenerateOTP()

		go helpers.SendOTP(request.Email, newCode, user.Name)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "A new verification code has been sent to your email",
		})
	}

}

func GetAllUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		collection := database.GetCollection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opts := options.Find().SetProjection(bson.M{"password": 0})

		cursor, err := collection.Find(ctx, bson.M{}, opts)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error fetching users",
				"error":   err.Error(),
			})
		}
		defer cursor.Close(ctx)

		var users []model.User
		if err := cursor.All(ctx, &users); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error parsing users",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"count":  len(users),
			"users":  users,
		})
	}
}

func GetUserByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid ID format",
			})
		}

		collection := database.GetCollection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opts := options.FindOne().SetProjection(bson.M{"password": 0})

		var user model.User
		err = collection.FindOne(ctx, bson.M{"_id": objectID}, opts).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":  "error",
					"message": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error fetching user",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"user":   user,
		})
	}
}

func UpdateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid ID format",
			})
		}

		type UpdateRequest struct {
			FullName string `json:"name"`
			Phone    string `json:"phone"`
			Email    string `json:"email"`
		}

		var request UpdateRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		update := bson.M{
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		}

		// Only update fields that are provided
		if request.FullName != "" {
			update["$set"].(bson.M)["name"] = request.FullName
		}
		if request.Phone != "" {
			update["$set"].(bson.M)["phone"] = request.Phone
		}
		if request.Email != "" {
			update["$set"].(bson.M)["email"] = request.Email
		}

		collection := database.GetCollection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		opts.SetProjection(bson.M{"password": 0})

		var updatedUser model.User
		err = collection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": objectID},
			update,
			opts,
		).Decode(&updatedUser)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":  "error",
					"message": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error updating user",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "User updated successfully",
			"user":    updatedUser,
		})
	}
}

func DeleteUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid ID format",
			})
		}

		collection := database.GetCollection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error deleting user",
				"error":   err.Error(),
			})
		}

		if result.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "User deleted successfully",
		})
	}
}

func LoginUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var request LoginRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		if request.Email == "" || request.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Email and password are required",
			})
		}

		collection := database.GetCollection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user model.User
		err := collection.FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid credentials",
			})
		}

		if !user.Verified {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Account not verified. Please verify your email before logging in",
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid credentials, password does not match",
			})
		}

		tokens, err := middleware.GenerateTokens(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error generating authentication token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Login successful",
			"user": fiber.Map{
				"id":       user.ID,
				"name":     user.Name,
				"email":    user.Email,
				"phone":    user.Phone,
				"verified": user.Verified,
			},
			"tokens": fiber.Map{
				"access_token":  tokens.AccessToken,
				"refresh_token": tokens.RefreshToken,
				"expires_in":    tokens.AtExpires,
			},
		})
	}
}

func RefreshToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type RefreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}

		var request RefreshRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid refresh token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid token claims",
			})
		}

		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid user ID in token",
			})
		}

		userId, err := primitive.ObjectIDFromHex(userIdStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid user ID format",
			})
		}

		newTokens, err := middleware.GenerateTokens(userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error generating tokens",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Token refreshed successfully",
			"tokens": fiber.Map{
				"access_token":  newTokens.AccessToken,
				"refresh_token": newTokens.RefreshToken,
				"expires_in":    newTokens.AtExpires,
			},
		})
	}
}

func LogoutUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenMetadata, err := middleware.ExtractTokenMetadata(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Unauthorized",
			})
		}

		tokenString := middleware.ExtractToken(c)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to parse token claims",
			})
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to get token expiration",
			})
		}

		middleware.BlacklistToken(tokenMetadata.AccessUuid, int64(exp))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Successfully logged out",
		})
	}
}
