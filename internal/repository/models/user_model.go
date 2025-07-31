package models

type User struct {
	ID       uint    `gorm:"primaryKey"`
	Email    string  `gorm:"unique;not null"`
	Password string  `gorm:"not null"`

	// ✅ One-to-One aloqasi
	Profile  Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
}
