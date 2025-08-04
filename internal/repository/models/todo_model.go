package models

type Todo struct {
    ID          int64  `gorm:"primaryKey;autoIncrement"`
    UserID      int64  `gorm:"not null"`
    CategoryID  *int64 `gorm:"index"`

    Title       string `gorm:"size:255;not null"`
    Description string `gorm:"type:text"`

    Category    *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
    User        *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
