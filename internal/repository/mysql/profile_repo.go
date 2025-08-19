package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"todolist/internal/domain"
	"todolist/internal/repository/models"
	"todolist/pkg/cache"

	"gorm.io/gorm"
)

type ProfileRepo struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewProfileRepo(db *gorm.DB, cache *cache.Cache) domain.ProfileRepository {
	return &ProfileRepo{db: db, cache: cache}
}

// CreateProfile
func (r *ProfileRepo) CreateProfile(userID int, profile *domain.Profile) error {
	newProfile := models.Profile{
		UserID:    uint(userID),
		Firstname: profile.Firstname,
		LastName:  profile.LastName,
		Username:  profile.Username,
		Image:     profile.Image,
	}

	return r.db.Create(&newProfile).Error
}

// GetProfileByUserID → READ-THROUGH
func (r *ProfileRepo) GetProfileByUserID(userID int) (*domain.Profile, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("profile:%d", userID)

	// 1️⃣ Avval cache dan qidiramiz
	if data, err := r.cache.Get(ctx, cacheKey); err == nil {
		var profile domain.Profile
		if err := json.Unmarshal([]byte(data), &profile); err == nil {
			return &profile, nil
		}
	}

	// 2️⃣ Agar cache’da bo‘lmasa, DB’dan olamiz
	var profile models.Profile
	err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}

	result := &domain.Profile{
		ID:        int(profile.ID),
		Firstname: profile.Firstname,
		LastName:  profile.LastName,
		Username:  profile.Username,
		Image:     profile.Image,
		User: &domain.User{
			ID:    int(profile.User.ID),
			Email: profile.User.Email,
		},
	}

	// 3️⃣ Cache ga yozamiz
	if bytes, err := json.Marshal(result); err == nil {
		_ = r.cache.Set(ctx, cacheKey, string(bytes))
	}

	return result, nil
}

// UpdateProfile → WRITE-THROUGH
// UpdateProfile → WRITE-THROUGH
// UpdateProfile → WRITE-THROUGH
func (r *ProfileRepo) UpdateProfile(userID int, profile *domain.Profile) error {
	// 1️⃣ DB’ni yangilaymiz
	err := r.db.Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"firstname": profile.Firstname,
			"lastname":  profile.LastName,
			"username":  profile.Username,
			"image":     profile.Image,
		}).Error
	if err != nil {
		return err
	}

	// 2️⃣ Yangilangan profilni DB’dan olib kelamiz (cache ishlatmasdan!)
	var updated models.Profile
	err = r.db.Preload("User").Where("user_id = ?", userID).First(&updated).Error
	if err != nil {
		return err
	}

	result := &domain.Profile{
		ID:        int(updated.ID),
		Firstname: updated.Firstname,
		LastName:  updated.LastName,
		Username:  updated.Username,
		Image:     updated.Image,
		User: &domain.User{
			ID:    int(updated.User.ID),
			Email: updated.User.Email,
		},
	}

	// 3️⃣ Cache’ni yangilaymiz → write-through
	ctx := context.Background()
	cacheKey := fmt.Sprintf("profile:%d", userID)

	if bytes, err := json.Marshal(result); err == nil {
		_ = r.cache.Set(ctx, cacheKey, string(bytes))
	}

	return nil
}
