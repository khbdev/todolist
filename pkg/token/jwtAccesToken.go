package token

import (
	"errors"
	"os"
	"time"
	"todolist/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user *domain.User) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET env not set")
    }

    expireMinutesStr := os.Getenv("JWT_EXPIRE_MINUTES")
    if expireMinutesStr == "" {
        expireMinutesStr = "15"
    }

    expireMin, err := time.ParseDuration(expireMinutesStr + "m")
    if err != nil {
        return "", err
    }

    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(expireMin).Unix(),
        "iat":     time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
