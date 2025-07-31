package models

type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex"` // Har bir foydalanuvchiga bitta profil
	Firstname string
	LastName  string
	Username  string
	Image     string
	User      *User `gorm:"foreignKey:UserID"` // 🔥 mana shu qator qo‘shilsin

	
}
