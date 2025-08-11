package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"todolist/internal/domain"
	"todolist/internal/usecase"

	"todolist/pkg/response"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	usecase *usecase.TodoUsecase
}

func NewTodoHandler(uc *usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{usecase: uc}
}

// userID olish helper
func getUserIDs(c *gin.Context) (int64, bool) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		response.Error(c, "User aniqlanmadi", http.StatusUnauthorized)
		return 0, false
	}

	switch v := userIDInterface.(type) {
	case string:
		userID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			response.Error(c, "User ID noto‘g‘ri formatda", http.StatusBadRequest)
			return 0, false
		}
		return userID, true
	case int:
		return int64(v), true
	case int64:
		return v, true
	case float64:
		return int64(v), true
	default:
		response.Error(c, "User ID noto‘g‘ri formatda", http.StatusBadRequest)
		return 0, false
	}
}

func (h *TodoHandler) Create(c *gin.Context) {
	userID, ok := getUserIDs(c)
	if !ok {
		return
	}

	var todo domain.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		response.Error(c, "Noto'g'ri JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("DEBUG category_id:", todo.CategoryID)
	todo.UserID = userID

	id, err := h.usecase.CreateTodo(todo)
	if err != nil {
		response.Error(c, "Yaratishda xatolik: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{"id": id, "message": "Todo yaratildi"})
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	userID, ok := getUserIDs(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "ID noto‘g‘ri", http.StatusBadRequest)
		return
	}

	todo, err := h.usecase.GetTodoByID(id)
	if err != nil {
		response.Error(c, "Topilmadi: "+err.Error(), http.StatusNotFound)
		return
	}

	if todo.UserID != userID {
		response.Error(c, "Ruxsat yo'q", http.StatusForbidden)
		return
	}

	response.Success(c, todo)
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	userID, ok := getUserIDs(c)
	if !ok {
		return
	}

	todos, err := h.usecase.GetTodosByUserID(userID)
	if err != nil {
		response.Error(c, "Xatolik: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Agar todos nil bo'lsa, bo'sh slice qilib qo'yamiz
	if todos == nil {
		todos = make([]domain.Todo, 0)
	}

	response.Success(c, todos)
}

func (h *TodoHandler) Update(c *gin.Context) {
	userID, ok := getUserIDs(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "ID noto‘g‘ri", http.StatusBadRequest)
		return
	}

	var todo domain.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		response.Error(c, "Noto'g'ri JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	todo.ID = id
	todo.UserID = userID

	if err := h.usecase.UpdateTodo(todo); err != nil {
		response.Error(c, "Yangilashda xatolik: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{"message": "Yangilandi"})
}

func (h *TodoHandler) Delete(c *gin.Context) {
	userID, ok := getUserIDs(c)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, "ID noto‘g‘ri", http.StatusBadRequest)
		return
	}

	todo, err := h.usecase.GetTodoByID(id)
	if err != nil {
		response.Error(c, "Topilmadi: "+err.Error(), http.StatusNotFound)
		return
	}

	if todo.UserID != userID {
		response.Error(c, "Ruxsat yo'q", http.StatusForbidden)
		return
	}

	if err := h.usecase.DeleteTodo(id); err != nil {
		response.Error(c, "O'chirishda xatolik: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{"message": "Todo o‘chirildi"})
}
