package usecase

import (
	
	"todolist/pkg/token"
)

type UserUsecases struct {
	// kerak bo‘lsa repo maydonlari
}

func (uc *UserUsecase) RefreshAccessToken(refreshToken string) (string, error) {
	return token.RefreshAccessToken(refreshToken)
}
