package handler

import (
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *UserHandler, profileHandler *ProfileHandler, userUC *usecase.UserUsecase) {
	// Auth
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Profile
	authMiddleware := AuthMiddleware(userUC)
	r.GET("/profile/", authMiddleware, profileHandler.GetMyProfile)
}


