package handler

import (
	"todolist/internal/repository/mysql"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Repos
	userRepo := mysql.NewUserRepo(db)
	profileRepo := mysql.NewProfileRepo(db)
	settingRepo := mysql.NewSettingRepo(db)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, profileRepo, settingRepo)
	profileUsecase := usecase.NewProfileUsecase(profileRepo)
	settingUsecase := usecase.NewSettingUsecase(settingRepo)

	// Handlers
	userHandler := NewUserHandler(userUsecase)
	profileHandler := NewProfileHandler(profileUsecase)
	settingHandler := NewSettingHandler(settingUsecase)

	// Middleware
	authMiddleware := AuthMiddleware(userUsecase)

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "golang todo-app"})
	})

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/logout", authMiddleware, userHandler.Logout)

	r.GET("/profile/", authMiddleware, profileHandler.GetMyProfile)
	r.PUT("/profile/", authMiddleware, profileHandler.UpdateProfile)

	r.GET("/setting/", authMiddleware, settingHandler.GetSetting)
	r.PUT("/setting/", authMiddleware, settingHandler.UpdateSetting)

	r.Static("/storage/images", "pkg/storage/images")
}
