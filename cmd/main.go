package main

import (
	"log"
	"todolist/internal/config"
	"todolist/internal/handler"
	"todolist/internal/repository/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, err := config.ConnectGormDB()
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

err = config.AutoMigrate(db, 
    &models.User{}, 
    &models.Profile{}, 
    &models.Setting{}, 
    &models.Category{}, 
    &models.Todo{}, // ✅
)


	// Hamma narsa bu yerda qilinadi, faqat router uzatiladi
	r := gin.Default()
	handler.SetupRoutes(r, db)

	log.Println("🚀 Server running on :8002")
	r.Run(":8001")
}
