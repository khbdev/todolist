package usecase

import (
	"todolist/internal/domain"
)

type ProfileUsecase struct {
	profileRepo domain.ProfileRepository
}

func NewProfileUsecase(profileRepo domain.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{profileRepo: profileRepo}
}

func (u *ProfileUsecase) GetMyProfile(userID int) (*domain.Profile, error) {
	return u.profileRepo.GetProfileByUserID(userID)
}
