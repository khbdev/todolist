package models

type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex"`
	Firstname string `gorm:"column:firstname"`
	LastName  string `gorm:"column:lastname"` // 👈 shu joy muhim!
	Username  string `gorm:"column:username"`
	Image     string `gorm:"column:image"`
	User      *User  `gorm:"foreignKey:UserID"`
}
