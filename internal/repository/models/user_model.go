package models


type User struct {
    ID          int64      `gorm:"primaryKey;autoIncrement"`
    Email       string     `gorm:"unique;not null"`
    Password    string     `gorm:"not null"`
    Token       string

    Profile     Profile    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
    Setting     Setting    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
    Categories  []Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
    Todos       []Todo     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}