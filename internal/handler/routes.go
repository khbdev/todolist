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
	categoryRepo := mysql.NewCategoryRepo(db)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, profileRepo, settingRepo)
	profileUsecase := usecase.NewProfileUsecase(profileRepo)
	settingUsecase := usecase.NewSettingUsecase(settingRepo)
	categoryUscase := usecase.NewCategoryUsecase(categoryRepo)

	// Handlers
	userHandler := NewUserHandler(userUsecase)
	profileHandler := NewProfileHandler(profileUsecase)
	settingHandler := NewSettingHandler(settingUsecase)
	categoryHandler := NewCategoryHandler(categoryUscase)
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

	r.GET("/categories", authMiddleware, categoryHandler.GetAll)
r.POST("/categories", authMiddleware, categoryHandler.Create)
r.GET("/categories/:id", authMiddleware, categoryHandler.GetByID)
r.PUT("/categories/:id", authMiddleware, categoryHandler.Update)
r.DELETE("/categories/:id", authMiddleware, categoryHandler.Delete)

	r.Static("/storage/images", "pkg/storage/images")
}
