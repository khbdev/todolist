package usecase

import (
	"errors"
	"fmt"
	"todolist/internal/domain"
	"todolist/pkg/hash"
	"todolist/pkg/token"
)

type UserUsecase struct {
	userRepo    domain.UserRepository
	profileRepo domain.ProfileRepository // 🔥 yangi qo‘shiladi
}

func NewUserUsecase(userRepo domain.UserRepository, profileRepo domain.ProfileRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

func (u *UserUsecase) Register(user *domain.User) error {
	// 1. Email tekshiruvi
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("this email is already registered")
	}

	// 2. Parolni hash qilish
	hashed, err := hash.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	// 3. Userni yaratish
	err = u.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	// 4. Bo‘sh Profile yaratish (user ID asosida)
	profile := &domain.Profile{
		Firstname: "",
		LastName:  "",
		Username:  "",
		Image:     "",
	}
	err = u.profileRepo.CreateProfile(user.ID, profile)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) Login(email, password string) (*domain.User, string, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if !hash.CheckPassword(password, user.Password) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	tokenStr, err := token.GenerateToken(16)
	if err != nil {
		return nil, "", err
	}

	return user, tokenStr, nil
}
