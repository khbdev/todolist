package models

type User struct {
    ID         uint        `gorm:"primaryKey"`
    Email      string      `gorm:"unique;not null"`
    Password   string      `gorm:"not null"`
    Token      string

  Profile  Profile     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
Setting  Setting     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
Category []Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`

}
