package handler

import (

	
	"net/http"

	
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

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, "User aniqlanmadi", 401)
		return
	}

	profile, err := h.profileUC.GetMyProfile(userID.(int))
	if err != nil {
		response.Error(c, "Profil topilmadi", 404)
		return
	}

	response.Success(c, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
    userIDInterface, exists := c.Get("userID")
    if !exists {
        response.Error(c, "Token topilmadi", http.StatusUnauthorized)
        return
    }

    userID, ok := userIDInterface.(int)
    if !ok {
        response.Error(c, "User ID noto‘g‘ri formatda", http.StatusBadRequest)
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
