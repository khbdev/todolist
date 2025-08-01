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
	// .env ni yuklash
	config.LoadEnv()

	// DB ulanish
	db, err := config.ConnectGormDB()
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

	// AutoMigrate qilish
	err = config.AutoMigrate(db, &models.User{}, &models.Profile{}, &models.Setting{})
	if err != nil {
		log.Fatalf("AutoMigrate xatolik: %v", err)
	}

	// Repositorylar
	userRepo := mysql.NewUserRepo(db)
	profileRepo := mysql.NewProfileRepo(db)
	settingRepo := mysql.NewSettingRepo(db)

	// Usecaselar
	userUsecase := usecase.NewUserUsecase(userRepo, profileRepo, settingRepo)
	profileUsecase := usecase.NewProfileUsecase(profileRepo)
	settingUsecase := usecase.NewSettingUsecase(settingRepo)

	// Handlerlar
	userHandler := handler.NewUserHandler(userUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	settingHandler := handler.NewSettingHandler(settingUsecase)

	// Gin router
	r := gin.Default()

	// Marshrutlar (routes)
	handler.SetupRoutes(r, userHandler, profileHandler, settingHandler, userUsecase)

	// Serverni ishga tushirish
	log.Println("🚀 Server running on :8002")
	r.Run(":8002")
}
