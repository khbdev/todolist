package domain

type UserRepository interface {
    CreateUser(user *User) error
    GetUserByEmail(email string) (*User, error)
}

type ProfileRepository interface {
	CreateProfile(userID int, profile *Profile) error
	GetProfileByUserID(userID int) (*Profile, error)
	UpdateProfile(userID int, profile *Profile) error
}