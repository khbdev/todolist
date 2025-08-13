package admin

import (
    "log"
    "os"
    "todolist/internal/repository/models"
    "todolist/pkg/hash"

    "gorm.io/gorm"
)

func CreateAdmin(db *gorm.DB) {
    adminEmail := os.Getenv("ADMIN_EMAIL")
    adminPassword := os.Getenv("ADMIN_PASSWORD")
    adminRole := os.Getenv("ADMIN_ROLE")

    if adminEmail == "" || adminPassword == "" || adminRole == "" {
        log.Println("ADMIN credentials not set in .env, skipping admin creation")
        return
    }

    var existing models.User
    result := db.Where("email = ?", adminEmail).First(&existing)
    if result.Error == nil {
        log.Println("Admin user already exists")
        return
    }

    // Passwordni hash qilish pkg/hash yordamida
    hashedPassword, err := hash.HashPassword(adminPassword)
    if err != nil {
        log.Fatalf("Failed to hash admin password: %v", err)
    }

    admin := models.User{
        Email:    adminEmail,
        Password: hashedPassword,
        Role:     adminRole,
    }

    if err := db.Create(&admin).Error; err != nil {
        log.Fatalf("Failed to create admin user: %v", err)
    }

    log.Println("Admin user created successfully")
}
