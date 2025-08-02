package models

type Todo struct {
    ID          int64  `gorm:"primaryKey;autoIncrement"`
    UserID      int64  `gorm:"not null"`
    CategoryID  *int64 `gorm:"index"`
    Title       string `gorm:"size:255;not null"`
    Description string `gorm:"type:text"`
}