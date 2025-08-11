package usecase

import (
	"errors"
	"todolist/internal/domain"
	"todolist/pkg/token"
)

type UserUsecase struct {
	userRepo    domain.UserRepository
	profileRepo domain.ProfileRepository
	settingRepo domain.SettingRepository

}

func NewUserUsecase(userRepo domain.UserRepository, profileRepo domain.ProfileRepository,settingRepo domain.SettingRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		settingRepo: settingRepo,
	}
}

func (uc *UserUsecase) Register(user *domain.User) error {
	// Foydalanuvchini yaratish
	err := uc.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	// Profile yaratish
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

	// Setting yaratish
	setting := &domain.Setting{
		BgColor:   "#FFD701",
		TextColor: "#FFD700",
	}
	err = uc.settingRepo.CreateSetting(user.ID, setting)
	if err != nil {
		return err
	}

	return nil
}




func (uc *UserUsecase) Login(email, password string) (*domain.User, string, string, error) {
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", "", err
	}
	if user.Password != password {
		return nil, "", "", errors.New("noto‘g‘ri parol yoki email")
	}

	// Access token yaratish
	accessToken, err := token.GenerateJWT(user)
	if err != nil {
		return nil, "", "", err
	}

	// Refresh token yaratish (token paketida shunday funksiya bo‘lishi kerak)
	refreshToken, err := token.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}



func (uc *UserUsecase) Logout(token string) error {
	user, err := uc.userRepo.GetUserByToken(token)
	if err != nil {
		return err
	}

	err = uc.userRepo.UpdateToken(user.ID, "")
	if err != nil {
		return err
	}

	return nil
}


func (uc *UserUsecase) GetUserIDByToken(tok string) (int, error) {
	user, err := uc.userRepo.GetUserByToken(tok)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

