package handler

import (
	"todolist/internal/domain"
	"todolist/internal/usecase"
	"todolist/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error(), 400)
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.userUsecase.Register(user); err != nil {
		response.Error(c, "Ro‘yxatdan o‘tishda xatolik", 500)
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error(), 400)
		return
	}

	user, token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, "Email yoki parol xato", 401)
		return
	}

	response.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, "Token yuborilmagan", 401)
		return
	}

	// Bearer prefiksni ajratib olish
	const prefix = "Bearer "
	var token string
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		token = authHeader[len(prefix):]
	} else {
		token = authHeader
	}

	err := h.userUsecase.Logout(token)
	if err != nil {
		response.Error(c, "Logout amalga oshmadi", 500)
		return
	}

	response.Success(c, gin.H{
		"message": "Muvaffaqiyatli chiqildi",
	})
}


