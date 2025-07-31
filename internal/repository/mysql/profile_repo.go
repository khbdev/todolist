package mysql

import (
	"todolist/internal/domain"
	"todolist/internal/repository/models"

	"gorm.io/gorm"
)

type ProfileRepo struct {
	db *gorm.DB
}

func NewProfileRepo(db *gorm.DB) domain.ProfileRepository {
	return &ProfileRepo{db: db}
}

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

func (r *ProfileRepo) GetProfileByUserID(userID int) (*domain.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}

	return &domain.Profile{
		ID:        int(profile.ID),
		Firstname: profile.Firstname,
		LastName:  profile.LastName,
		Username:  profile.Username,
		Image:     profile.Image,
	}, nil
}

func (r *ProfileRepo) UpdateProfile(userID int, profile *domain.Profile) error {
	return r.db.Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"firstname": profile.Firstname,
			"lastname":  profile.LastName,
			"username":  profile.Username,
			"image":     profile.Image,
		}).Error
}
