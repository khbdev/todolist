package handler

import (
	"todolist/internal/usecase"
	"todolist/pkg/response"

	"github.com/gin-gonic/gin"
)

type CUserHandler struct {
	userUsecase *usecase.UserUsecases
}

func CreateUserHandler(uuc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uuc}
}

func (h *UserHandler) RefreshHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "refresh_token is required", 400)
		return
	}

	accessToken, err := h.userUsecase.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		response.Error(c, "Invalid or expired refresh token", 401)
		return
	}

	response.Success(c, gin.H{
		"access_token": accessToken,
	})
}
