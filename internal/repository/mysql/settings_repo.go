package mysql

import (
	"encoding/json"
	"fmt"
	"todolist/internal/domain"
	"todolist/pkg/cache"

	"gorm.io/gorm"
)

type SettingRepo struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewSettingRepo(db *gorm.DB, cache *cache.Cache) domain.SettingRepository {
	return &SettingRepo{db: db, cache: cache}
}

// CreateSetting - write-through cache
func (r *SettingRepo) CreateSetting(userID int, setting *domain.Setting) error {
	setting.UserID = uint(userID)

	if err := r.db.Create(setting).Error; err != nil {
		return err
	}

	// Redisga yozish
	key := fmt.Sprintf("setting:user:%d", userID)
	val, _ := json.Marshal(setting)
	return r.cache.Set(r.db.Statement.Context, key, string(val)) // context DB dan olinadi
}

// GetSettingByUserID - read-through cache
func (r *SettingRepo) GetSettingByUserID(userID int) (*domain.Setting, error) {
	key := fmt.Sprintf("setting:user:%d", userID)

	// Avval cache’dan olishga urinish
	if val, err := r.cache.Get(r.db.Statement.Context, key); err == nil {
		var setting domain.Setting
		if err := json.Unmarshal([]byte(val), &setting); err == nil {
			return &setting, nil
		}
	}

	// Agar cache’da bo‘lmasa → DB’dan olish
	var setting domain.Setting
	if err := r.db.
		Where("user_id = ?", userID).
		First(&setting).Error; err != nil {
		return nil, err
	}

	// Cache’ga yozib qo‘yish
	val, _ := json.Marshal(&setting)
	_ = r.cache.Set(r.db.Statement.Context, key, string(val))

	return &setting, nil
}

// UpdateSetting - write-through cache
func (r *SettingRepo) UpdateSetting(userID int, setting *domain.Setting) error {
	if err := r.db.
		Model(&domain.Setting{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"bg_color":   setting.BgColor,
			"text_color": setting.TextColor,
		}).Error; err != nil {
		return err
	}

	// Cache’ni yangilash
	key := fmt.Sprintf("setting:user:%d", userID)
	val, _ := json.Marshal(setting)
	return r.cache.Set(r.db.Statement.Context, key, string(val))
}
