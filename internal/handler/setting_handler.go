package handler

import (
	"net/http"
	"strconv"
	"todolist/internal/domain"
	"todolist/pkg/response"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingUC domain.SettingRepository
}

func NewSettingHandler(settingUC domain.SettingRepository) *SettingHandler {
	return &SettingHandler{
		settingUC: settingUC,
	}
}

// userID ni kontekstdan xavfsiz olish helper
func getUserIDFromContext(c *gin.Context) (int64, bool) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Error(c, "Token topilmadi", http.StatusUnauthorized)
		return 0, false
	}

	switch v := userIDInterface.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case float64:
		return int64(v), true
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			response.Error(c, "User ID noto‘g‘ri formatda", http.StatusInternalServerError)
			return 0, false
		}
		return id, true
	default:
		response.Error(c, "User ID noto‘g‘ri formatda", http.StatusInternalServerError)
		return 0, false
	}
}

func (h *SettingHandler) GetSetting(c *gin.Context) {
	userID64, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	userID := int(userID64) // int64 dan int ga o‘tkazish

	setting, err := h.settingUC.GetSettingByUserID(userID)
	if err != nil {
		response.Error(c, "Sozlamani olishda xatolik", http.StatusInternalServerError)
		return
	}

	response.Success(c, setting)
}
func (h *SettingHandler) UpdateSetting(c *gin.Context) {
	userID64, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	userID := int(userID64) // int ga o‘tkazish

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
