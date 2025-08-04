package handler

import (
	"net/http"
	"strconv"
	"todolist/internal/domain"
	"todolist/internal/resource"
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
    usecase *usecase.CategoryUsecase
}

func NewCategoryHandler(uc *usecase.CategoryUsecase) *CategoryHandler {
    return &CategoryHandler{usecase: uc}
}

// Helper: userID olish
func getUserID(c *gin.Context) (int64, bool) {
    userIDInterface, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User aniqlanmadi"})
        return 0, false
    }

    switch v := userIDInterface.(type) {
    case int:
        return int64(v), true
    case int64:
        return v, true
    case float64:
        return int64(v), true
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID noto‘g‘ri formatda"})
        return 0, false
    }
}

func (h *CategoryHandler) Create(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }

    var cat domain.Category
    if err := c.ShouldBindJSON(&cat); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Noto'g'ri JSON: " + err.Error()})
        return
    }
    cat.UserID = userID

    if err := h.usecase.Create(c.Request.Context(), &cat); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Yaratishda xatolik: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Kategoriya yaratildi"})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }

    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID noto‘g‘ri"})
        return
    }

    category, err := h.usecase.GetByID(c.Request.Context(), userID, id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topilmadi: " + err.Error()})
        return
    }

   c.JSON(http.StatusOK, resource.NewCategoryWithTodosResponse(category))
}

func (h *CategoryHandler) Update(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok { return }

    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID noto‘g‘ri"})
        return
    }

    var cat domain.Category
    if err := c.ShouldBindJSON(&cat); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Noto'g'ri JSON: " + err.Error()})
        return
    }
    cat.ID = id
    cat.UserID = userID

    if err := h.usecase.Update(c.Request.Context(), &cat); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Yangilashda xatolik: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Yangilandi"})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }

    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID noto‘g‘ri"})
        return
    }

    if err := h.usecase.Delete(c.Request.Context(), userID, id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "O'chirishda xatolik: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Kategoriya o‘chirildi"})
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
    userID, ok := getUserID(c)
    if !ok {
        return
    }

    categories, err := h.usecase.GetAllByUserID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Xatolik: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, categories)
}