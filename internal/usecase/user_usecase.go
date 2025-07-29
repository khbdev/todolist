package usecase

import (
	"errors"
	"fmt"
	"todolist/internal/domain"
	"todolist/pkg/hash"
	"todolist/pkg/token"
)

type UserUsecase struct {
    userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
    return &UserUsecase{userRepo: repo}
}

func (u *UserUsecase) Register(user *domain.User) error {
    // Email mavjudligini tekshiramiz
    existingUser, err := u.userRepo.GetUserByEmail(user.Email)
    if err == nil && existingUser != nil {
        return errors.New("this email is already registered")
    }

    // Parolni hash qilamiz
    hashed, err := hash.HashPassword(user.Password)
    if err != nil {
        return err
    }

    user.Password = hashed
    return u.userRepo.CreateUser(user)
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

    // Optional: tokenni DBga saqlash yoki cache’ga

    return user, tokenStr, nil
}
