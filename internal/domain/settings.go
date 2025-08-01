package domain

type Setting struct {
	ID        int    `json:"id"`
	BgColor   string `json:"bg_color" validate:"required,min=4,max=20"`  // misol uchun: "#fff", "red", "#ffffff"
	TextColor string `json:"text_color" validate:"required,min=4,max=20"` // misol uchun: "#000", "black", "#333333"
UserID uint `json:"-"`

}