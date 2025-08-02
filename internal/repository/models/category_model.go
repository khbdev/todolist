package models



type Category struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"not null"`
	Name      string    `gorm:"size:255;not null"`
}
