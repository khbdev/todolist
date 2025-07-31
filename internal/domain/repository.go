package domain

type UserRepository interface {
CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
    
    // 🔽 Shu ikki metodni qo‘shing:
	UpdateToken(userID int, token string) error
	GetUserByToken(token string) (*User, error)

}

type ProfileRepository interface {
	CreateProfile(userID int, profile *Profile) error
	GetProfileByUserID(userID int) (*Profile, error)
	UpdateProfile(userID int, profile *Profile) error
}
