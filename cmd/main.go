package main

import (
	"log"

	"todolist/internal/admin"
	"todolist/internal/config"
	"todolist/internal/handler"
	"todolist/internal/repository/models"

	"todolist/pkg/cronjob"
	"todolist/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	db, err := config.ConnectGormDB()
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

	cronjob.InitConnection(db)

	cronjob.RunCronJob()

	err = config.AutoMigrate(db,
		&models.User{},
		&models.Profile{},
		&models.Setting{},
		&models.Category{},
		&models.Todo{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate Amalga oshmadi: %v", err)
	}

	admin.CreateAdmin(db)



	r := gin.Default()

	r.Use(handler.CORSMiddleware())

	reteLimiting := handler.NewRedisRateLimiter(config.NewRedis())
	r.Use(reteLimiting.RateLimit())

	handler.SetupRoutes(r, db)

	
	rmq := rabbitmq.GetInstance()  
	go rmq.Consume()         

	

	log.Println("🚀 Server ishga tushdi 8082")
	r.Run(":8082")
}
