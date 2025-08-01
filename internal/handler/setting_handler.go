package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todolist/internal/domain"
	"todolist/pkg/response"
)

type SettingHandler struct {
	settingUC domain.SettingRepository
}

func NewSettingHandler(settingUC domain.SettingRepository) *SettingHandler {
	return &SettingHandler{
		settingUC: settingUC,
	}
}

func (h *SettingHandler) GetSetting(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Error(c, "Token topilmadi", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		response.Error(c, "User ID noto‘g‘ri formatda", http.StatusInternalServerError)
		return
	}

	setting, err := h.settingUC.GetSettingByUserID(userID)
	if err != nil {
		response.Error(c, "Sozlamani olishda xatolik", http.StatusInternalServerError)
		return
	}

	response.Success(c, setting)
}

func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Error(c, "Token topilmadi", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		response.Error(c, "User ID noto‘g‘ri formatda", http.StatusInternalServerError)
		return
	}

	var setting domain.Setting
	if err := c.ShouldBindJSON(&setting); err != nil {
		response.Error(c, "JSON format xato", http.StatusBadRequest)
		return
	}

	err := h.settingUC.UpdateSetting(userID, &setting)
	if err != nil {
		response.Error(c, "Yangilashda xatolik", http.StatusInternalServerError)
		return
	}

	response.Success(c, "Sozlama muvaffaqiyatli yangilandi")
}

