package models


// models/category.go
type Category struct {
    ID     int64  `gorm:"primaryKey;autoIncrement"`
    UserID int64  `gorm:"not null"`
    Name   string `gorm:"size:255;not null"`

    Todos []Todo `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
}