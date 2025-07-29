package main

import (
	"log"
	"todolist/internal/config"
	"todolist/internal/domain"
	"todolist/internal/handler"
	"todolist/internal/repository/mysql"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)


func main(){
	config.LoadEnv()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("❌ DB ulanishda xatolik: %v", err)
	}

	// 🔽 AutoMigrate
	err = config.AutoMigrate(db, &domain.User{}) // kerakli modelni yoz
	if err != nil {
		log.Fatalf("❌ AutoMigrate xatolik: %v", err)
	}
		userRepo := mysql.NewUserRepo(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

		r := gin.Default()
	handler.SetupRoutes(r, userHandler)

	// 5. Serverni ishga tushirish
	log.Println("🚀 Server running on :8080")
	r.Run(":8002")
}