package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env fayl topilmadi yoki yuklanmadi")
	}
}

// config/gorm.go (yoki shu faylga qo‘sh)
func ConnectGormDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("✅ GORM DB muvaffaqiyatli ulandi")
	return gormDB, nil
}


// 🔽 GORM bilan faqat migratsiya qilish
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		return err
	}

	log.Println("✅ AutoMigrate muvaffaqiyatli bajarildi")
	return nil
}
