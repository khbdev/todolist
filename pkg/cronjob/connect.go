package cronjob

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var RDB *redis.Client
var DB *gorm.DB
var Interval time.Duration

func InitConnection(db *gorm.DB) {
	
	DB = db

	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	
	Interval = 24 * time.Hour

	log.Println("✅ DB va Redis ulanishi tayyor, cron interval:", Interval)
}
