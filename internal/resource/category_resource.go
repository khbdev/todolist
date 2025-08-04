package resource

import "todolist/internal/domain"

type TodoResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CategoryWithTodosResponse struct {
	ID     int64          `json:"id"`
	UserID int64          `json:"user_id"`
	Name   string         `json:"name"`
	Todos  []TodoResponse `json:"todos"`
}


func NewCategoryWithTodosResponse(category *domain.Category) *CategoryWithTodosResponse {
	var todos []TodoResponse

	for _, todo := range category.Todos {
		todos = append(todos, TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		})
	}

	return &CategoryWithTodosResponse{
		ID:     category.ID,
		UserID: category.UserID,
		Name:   category.Name,
		Todos:  todos,
	}
}
