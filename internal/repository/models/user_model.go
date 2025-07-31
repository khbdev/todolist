package models

type User struct {
	ID       uint    `gorm:"primaryKey"`
	Email    string  `gorm:"unique;not null"`
	Password string  `gorm:"not null"`
    Token string 

	// ✅ One-to-One aloqasi
	Profile  Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
}
