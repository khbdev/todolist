package domain


type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"-"`
}
