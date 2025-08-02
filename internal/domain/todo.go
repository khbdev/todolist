package domain

type Todo struct {
    ID          int64
    UserID      int64
    CategoryID  *int64  // agar kerak bo'lsa
    Title       string
    Description string
}