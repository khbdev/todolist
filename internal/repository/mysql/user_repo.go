package mysql

import (
	"todolist/internal/domain"
	"todolist/internal/repository/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *domain.User) error {
	newUser := models.User{
		Email:    user.Email,
		Password: user.Password,
		Token:    user.Token,
	}
	result := r.db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}
	user.ID = int(newUser.ID) // 🔥 Muhim qadam!
	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       int(user.ID),
		Email:    user.Email,
		Password: user.Password,
		Token:    user.Token,
	}, nil
}

func (r *UserRepo) GetUserByToken(token string) (*domain.User, error) {
	var user models.User
	err := r.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       int(user.ID),
		Email:    user.Email,
		Password: user.Password,
		Token:    user.Token,
	}, nil
}

func (r *UserRepo) UpdateToken(userID int, token string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("token", token).Error
}
