package domain

type Todo struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"-"` // token’dan
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	CategoryID  *int64  `json:"category_id"` // optional
}
