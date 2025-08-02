package domain

import (
	"context"
	"errors"
)

type UserRepository interface {
    CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateToken(userID int, token string) error
	GetUserByToken(token string) (*User, error)
}

type ProfileRepository interface {
	CreateProfile(userID int, profile *Profile) error
	GetProfileByUserID(userID int) (*Profile, error)
	UpdateProfile(userID int, profile *Profile) error
}


type SettingRepository interface {
	CreateSetting(userID int, setting *Setting) error
	GetSettingByUserID(userID int) (*Setting, error)
	UpdateSetting(userID int, setting *Setting) error
}


type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id, userID int64) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id, userID int64) error
	GetAllByUserID(ctx context.Context, userID int64) ([]*Category, error)
}


type TodoRepository interface {
    CreateTodo(todo Todo) (int64, error)
    GetTodoByID(id int64) (*Todo, error)
    GetTodosByUserID(userID int64) ([]Todo, error)
    UpdateTodo(todo Todo) error
    DeleteTodo(id int64) error
}

// error uchuN



var (
	ErrCategoryNameRequired = errors.New("category name is required")
)