package handler

import (
	"net/http"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userUC *usecase.UserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token topilmadi"})
			c.Abort()
			return
		}

	
		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Noto‘g‘ri Authorization formati"})
			c.Abort()
			return
		}

		token := authHeader[len(bearerPrefix):]

		userID, err := userUC.GetUserIDByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Noto‘g‘ri token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
