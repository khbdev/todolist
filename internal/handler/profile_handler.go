package handler

import (
	"net/http"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileUC *usecase.ProfileUsecase
}

func NewProfileHandler(profileUC *usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{profileUC: profileUC}
}

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User aniqlanmadi"})
		return
	}

	profile, err := h.profileUC.GetMyProfile(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil topilmadi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}
