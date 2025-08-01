package models

type Setting struct {
	ID        int    `gorm:"primaryKey"`
	BgColor   string
	TextColor string

	UserID uint   `gorm:"unique"` // One-to-One uchun unique bo'lishi kerak
	User   *User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}