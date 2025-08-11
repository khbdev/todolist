package handler

import (
	"net/http"
	"strconv"
	"todolist/internal/usecase"
	"todolist/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileUC *usecase.ProfileUsecase
}

func NewProfileHandler(profileUC *usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{profileUC: profileUC}
}

// yordamchi: kontekstdan userID ni xavfsiz olish
func GetUserIDFromContext(c *gin.Context) (int, bool) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Error(c, "User aniqlanmadi", http.StatusUnauthorized)
		return 0, false
	}

	switch v := userIDInterface.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			response.Error(c, "User ID noto‘g‘ri formatda", http.StatusBadRequest)
			return 0, false
		}
		return id, true
	default:
		response.Error(c, "User ID noto‘g‘ri formatda", http.StatusBadRequest)
		return 0, false
	}
}

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)  // bosh harf bilan chaqirish
	if !ok {
		return
	}

	profile, err := h.profileUC.GetMyProfile(userID)
	if err != nil {
		response.Error(c, "Profil topilmadi", 404)
		return
	}

	response.Success(c, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := GetUserIDFromContext(c)  // bosh harf bilan chaqirish
	if !ok {
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		response.Error(c, "Formni o'qishda xato: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedProfile, err := h.profileUC.UpdateProfileUsecase(userID, c.Request.MultipartForm)
	if err != nil {
		response.Error(c, "Yangilashda xatolik: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{
		"message": "Profil muvaffaqiyatli yangilandi",
		"data":    updatedProfile,
	})
}
