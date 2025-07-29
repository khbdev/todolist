package domain


type User struct {
    ID       int    `json:"id"`
    Username string `json:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"-"`
}
