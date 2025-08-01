package usecase

import (

	"todolist/internal/domain"
)

type SettingUsecase struct {
	settingRepo domain.SettingRepository
}

func NewSettingUsecase(repo domain.SettingRepository) *SettingUsecase {
	return &SettingUsecase{
		settingRepo: repo,
	}
}

// ✅ GET Setting
func (s *SettingUsecase) GetSettingByUserID(userID int) (*domain.Setting, error) {
	return s.settingRepo.GetSettingByUserID(userID)
}

// ✅ UPDATE Setting
func (s *SettingUsecase) UpdateSetting(userID int, setting *domain.Setting) error {
	return s.settingRepo.UpdateSetting(userID, setting)
}

// ✅ CREATE Setting
func (s *SettingUsecase) CreateSetting(userID int, setting *domain.Setting) error {
    return s.settingRepo.CreateSetting(userID, setting)
}