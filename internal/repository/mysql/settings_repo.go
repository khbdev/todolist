package mysql

import (
	"todolist/internal/domain"

	"gorm.io/gorm"
)

type SettingRepo struct {
	db *gorm.DB
}

func NewSettingRepo(db *gorm.DB) domain.SettingRepository {
	return &SettingRepo{db: db}
}


func (r *SettingRepo) CreateSetting(userID int, setting *domain.Setting) error {
	setting.UserID = uint(userID)
	return r.db.Create(setting).Error
}


func (r *SettingRepo) GetSettingByUserID(userID int) (*domain.Setting, error) {
	var setting domain.Setting
	err := r.db.Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}


func (r *SettingRepo) UpdateSetting(userID int, setting *domain.Setting) error {
	return r.db.Model(&domain.Setting{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"bg_color":   setting.BgColor,
			"text_color": setting.TextColor,
		}).Error
}
