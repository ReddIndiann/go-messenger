package helpers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/patrickmn/go-cache"
)

var otpCache = cache.New(5*time.Minute, 10*time.Minute)

func GenerateOTP() string {
	return fmt.Sprintf("%04d", rand.Intn(100000))
}

func StoreOTP(email string, otp string) {
	otpCache.Set(email, otp, cache.DefaultExpiration)
	fmt.Println("OTP stored in cache for:", email)
}

func GetOTP(email string) (string, bool) {
	otp, found := otpCache.Get(email)
	if found {
		return otp.(string), true
	}
	return "", false
}

func DeleteOTP(email string) {
	otpCache.Delete(email)
	fmt.Println("OTP deleted for:", email)
}