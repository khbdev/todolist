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

func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
    return &UserHandler{userUsecase: usecase}
}

func (h *UserHandler) Register(c *gin.Context) {
    var user struct {
        Username string `json:"username" binding:"required,min=3,max=20"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&user); err != nil {
        response.Error(c, err.Error(), 400)
        return
    }

   domainUser := &domain.User{
    Username: user.Username,
    Email:    user.Email,
    Password: user.Password, // ✅ Password qo‘shildi
}

    if err := h.userUsecase.Register(domainUser); err != nil {
        response.Error(c, "Registration failed", 500)
        return
    }

    response.Success(c, domainUser)
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
        response.Error(c, "Invalid email or password", 401)
        return
    }

    response.Success(c, gin.H{
        "user":  user,
        "token": token,
        "msg":   "Logged in successfully",
    })
}
