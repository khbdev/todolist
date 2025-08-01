package main

import (
	"log"
	"todolist/internal/config"
	"todolist/internal/handler"
	"todolist/internal/repository/models"
	"todolist/internal/repository/mysql"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)


func main() {
	config.LoadEnv()

	db, err := config.ConnectGormDB()
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

	// AutoMigrate
	err = config.AutoMigrate(db, &models.User{}, &models.Profile{})
	if err != nil {
		log.Fatalf("AutoMigrate xatolik: %v", err)
	}

	// Repositories
	userRepo := mysql.NewUserRepo(db)
	profileRepo := mysql.NewProfileRepo(db)

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, profileRepo)
	profileUsecase := usecase.NewProfileUsecase(profileRepo) 

	// Handlers
	userHandler := handler.NewUserHandler(userUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase) 

	// Gin
	r := gin.Default()

	// Routes
	handler.SetupRoutes(r, userHandler, profileHandler, userUsecase)

	log.Println("🚀 Server running on :8002")
	r.Run(":8002")
}

