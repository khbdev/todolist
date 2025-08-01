package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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

func (up *ProfileUsecase) UpdateProfileUsecase(userID int, form *multipart.Form) (*domain.Profile, error) {
    profile := &domain.Profile{
        Firstname: form.Value["firstname"][0],
        LastName:  form.Value["lastname"][0],
        Username:  form.Value["username"][0],
    }

    // Faylni olish
    files := form.File["image"]
    if len(files) > 0 {
        fileHeader := files[0]
        file, err := fileHeader.Open()
        if err != nil {
            return nil, fmt.Errorf("faylni ochishda xato: %v", err)
        }
        defer file.Close()

        filename := filepath.Base(fileHeader.Filename)
        destination := filepath.Join("pkg", "storage", "images", filename)

        if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
            return nil, fmt.Errorf("papka yaratishda xato: %v", err)
        }

        dst, err := os.Create(destination)
        if err != nil {
            return nil, fmt.Errorf("fayl yaratishda xato: %v", err)
        }
        defer dst.Close()

        _, err = io.Copy(dst, file)
        if err != nil {
            return nil, fmt.Errorf("faylni nusxalashda xato: %v", err)
        }

        profile.Image = "/storage/images/" + filename
    } else {
        profile.Image = form.Value["image"][0]
    }

    if err := up.profileRepo.UpdateProfile(userID, profile); err != nil {
        return nil, err
    }

    return profile, nil
}
