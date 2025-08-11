package token

import (
	"errors"
	"os"
	"time"
	"todolist/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

// Secret va expire vaqtlarni .env dan o‘qish
func getRefreshSecret() ([]byte, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_REFRESH_SECRET env not set")
	}
	return []byte(secret), nil
}

func getRefreshExpireDuration() (time.Duration, error) {
	expireDays := os.Getenv("JWT_REFRESH_DAYS")
	if expireDays == "" {
		expireDays = "7" // default 7 kun
	}
	return time.ParseDuration(expireDays + "24h") // kunni soatga aylantiramiz
}

// Refresh token yaratish
func GenerateRefreshToken(user *domain.User) (string, error) {
	secret, err := getRefreshSecret()
	if err != nil {
		return "", err
	}

	expireDuration, err := getRefreshExpireDuration()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(expireDuration).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh", // token turini belgilaymiz
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

// Refresh tokenni tekshirish va yangi access token yaratish
func RefreshAccessToken(refreshTokenString string) (string, error) {
	secret, err := getRefreshSecret()
	if err != nil {
		return "", err
	}

	// Tokenni parse qilish
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		// Algoritm tekshirish
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}

	// Claimsni olish va tekshirish
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	// Token turi refresh ekanligini tekshirish
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		return "", errors.New("invalid token type")
	}

	// Foydalanuvchi ma'lumotlarini olish
	userIDFloat, ok1 := claims["user_id"].(float64) // JWT da raqamlar float64 bo‘ladi
	email, ok2 := claims["email"].(string)
	if !ok1 || !ok2 {
		return "", errors.New("invalid token claims")
	}

	userID := int(userIDFloat)

	// Yangi access token yaratamiz
	// (bu yerda token paketidagi GenerateJWT funksiyasidan foydalanamiz)
	user := &domain.User{
		ID:    userID,
		Email: email,
	}

	return GenerateJWT(user)
}
