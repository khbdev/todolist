package domain

type Profile struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname" validate:"required,min=2,max=50"`
	LastName  string `json:"lastname" validate:"required,min=2,max=50"`
	Username  string `json:"username" validate:"required,alphanum,min=3,max=30"`
	Image string `json:"image" validate:"omitempty"`
	User      *User 
}
