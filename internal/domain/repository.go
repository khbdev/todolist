package domain





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
	Create(category *Category) error
	GetByID(id int64) (*Category, error)
	Update(category *Category) error
	Delete(id int64) error
	GetAllByUserID(userID int64) ([]*Category, error)
}