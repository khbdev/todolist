package models

import "gorm.io/gorm"


type User struct {
    ID          int64      `gorm:"primaryKey;autoIncrement"`
    Email       string     `gorm:"unique;not null"`
    Password    string     `gorm:"not null"`
    Token       string
    Role        string

    Profile     Profile    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
    Setting     Setting    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
    Categories  []Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
    Todos       []Todo     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
    var users []User
    if err := db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}