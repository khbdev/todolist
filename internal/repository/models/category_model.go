package models

type Category struct {
    ID     uint   `gorm:"primaryKey"`
    Name   string
    UserID uint // yoki uint64, agar User.ID uint64 bo‘lsa
    User   User  `gorm:"foreignKey:UserID"`
}