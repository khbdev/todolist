package handler

import (
	"net/http"

	"todolist/internal/admin/handler"
	"todolist/internal/repository/mysql"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Middleware’ni global qo‘shamiz
	r.Use(CORSMiddleware())
// end
	// Repos
	userRepo := mysql.NewUserRepo(db)
	profileRepo := mysql.NewProfileRepo(db)
	settingRepo := mysql.NewSettingRepo(db)
	categoryRepo := mysql.NewCategoryRepo(db)
	todoRepo := mysql.NewTodoRepository(db)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, profileRepo, settingRepo)
	profileUsecase := usecase.NewProfileUsecase(profileRepo)
	settingUsecase := usecase.NewSettingUsecase(settingRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	todoUsecase := usecase.NewTodoUsecase(todoRepo)

	// Handlers
	userHandler := NewUserHandler(userUsecase)
	profileHandler := NewProfileHandler(profileUsecase)
	settingHandler := NewSettingHandler(settingUsecase)
	categoryHandler := NewCategoryHandler(categoryUsecase)
	todoHandler := NewTodoHandler(todoUsecase)

	// Auth middleware faqat kerakli route’larda
	authMiddleware := AuthMiddleware()

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "golang todo-app"})
	})

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/refresh", userHandler.RefreshHandler)
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

	r.POST("/todo", authMiddleware, todoHandler.Create)
	r.GET("/todo", authMiddleware, todoHandler.GetAll)
	r.GET("/todo/:id", authMiddleware, todoHandler.GetByID)
	r.PUT("/todo/:id", authMiddleware, todoHandler.Update)
	r.DELETE("/todo/:id", authMiddleware, todoHandler.Delete)
	r.GET("/salom", authMiddleware, AdminOnly(), handler.Salom)

	// Static files uchun
	r.Static("/storage/images", "pkg/storage/images")
}