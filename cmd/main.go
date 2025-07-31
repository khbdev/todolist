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

	db, err := config.ConnectGormDB() // ✅ to‘g‘ri nom
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

	// ✅ AutoMigrate
	err = config.AutoMigrate(db, &models.User{}, &models.Profile{}) // models dan
	if err != nil {
		log.Fatalf("AutoMigrate xatolik: %v", err)
	}

	userRepo := mysql.NewUserRepo(db)
	profileRepo := mysql.NewProfileRepo(db)

	userUsecase := usecase.NewUserUsecase(userRepo, profileRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	r := gin.Default()
	handler.SetupRoutes(r, userHandler)

	log.Println("🚀 Server running on :8002")
	r.Run(":8002")
}
