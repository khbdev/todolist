package handler

import (
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *UserHandler, profileHandler *ProfileHandler, userUC *usecase.UserUsecase) {
 
  r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "golang todo-app",
    })
  })
	
	// Auth
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Profile
	authMiddleware := AuthMiddleware(userUC)
	r.GET("/profile/", authMiddleware, profileHandler.GetMyProfile)
	r.PUT("/profile/", authMiddleware, profileHandler.UpdateProfile)

	// Image Path
	 r.Static("/storage/images", "pkg/storage/images")
}


