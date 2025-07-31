package usecase

import (
	"errors"
	"todolist/internal/domain"
	"todolist/pkg/token"
)

type UserUsecase struct {
	userRepo    domain.UserRepository
	profileRepo domain.ProfileRepository
}

func NewUserUsecase(userRepo domain.UserRepository, profileRepo domain.ProfileRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

func (uc *UserUsecase) Register(user *domain.User) error {
	tok, err := token.GenerateToken(16)
	if err != nil {
		return err
	}
	user.Token = tok

	// 1. Foydalanuvchini yaratamiz
	err = uc.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	// 2. User ID avtomatik shakllanmagan bo‘lsa, uni qayta olish kerak bo'lishi mumkin
	// lekin agar `CreateUser` ichida `user.ID` set bo‘lsa, bevosita ishlatamiz:

	// 3. Default profile yaratamiz (bo‘sh yoki default qiymatlar bilan)
	profile := &domain.Profile{
		Firstname: "Azizbek",
		LastName:  "Xasanov",
		Username:  "khbdev",
		Image:     "rasm.jpg",
	}
	err = uc.profileRepo.CreateProfile(user.ID, profile)
	if err != nil {
		return err
	}

	return nil
}


func (uc *UserUsecase) Login(email, password string) (*domain.User, string, error) {
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if user.Password != password {
		return nil, "", errors.New("noto‘g‘ri parol yoki email")
	}

	newToken, err := token.GenerateToken(16)
	if err != nil {
		return nil, "", err
	}

	err = uc.userRepo.UpdateToken(user.ID, newToken)
	if err != nil {
		return nil, "", err
	}

	user.Token = newToken
	return user, newToken, nil
}

func (uc *UserUsecase) GetUserIDByToken(tok string) (int, error) {
	user, err := uc.userRepo.GetUserByToken(tok)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

