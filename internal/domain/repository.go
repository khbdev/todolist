package domain

type UserRepository interface {
    CreateUser(user *User) error
    GetUserByEmail(email string) (*User, error)
}