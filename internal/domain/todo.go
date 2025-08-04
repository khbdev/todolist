package domain

type Todo struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"-"` // JSON'da chiqmaydi
	CategoryID  *int64 `json:"-"` // JSON'da chiqmaydi
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=5,max=1000"`
}
