// mysql/todo_repository.go
package mysql

import (
	"todolist/internal/domain"
	"todolist/internal/repository/models"

	"gorm.io/gorm"
)


type todoRepository struct {
    db *gorm.DB
}


func NewTodoRepository(db *gorm.DB) domain.TodoRepository {
    return &todoRepository{db: db}
}

func (r *todoRepository) CreateTodo(todo domain.Todo) (int64, error) {
    modelTodo := models.Todo{
        UserID:      todo.UserID,
        CategoryID:  todo.CategoryID,
        Title:       todo.Title,
        Description: todo.Description,
    }

    result := r.db.Create(&modelTodo)
    return modelTodo.ID, result.Error
}

func (r *todoRepository) GetTodoByID(id int64) (*domain.Todo, error) {
    var modelTodo models.Todo
    err := r.db.First(&modelTodo, id).Error
    if err != nil {
        return nil, err
    }

    domainTodo := domain.Todo{
        ID:          modelTodo.ID,
        UserID:      modelTodo.UserID,
        CategoryID:  modelTodo.CategoryID,
        Title:       modelTodo.Title,
        Description: modelTodo.Description,
    }

    return &domainTodo, nil
}

func (r *todoRepository) GetTodosByUserID(userID int64) ([]domain.Todo, error) {
    var modelTodos []models.Todo
    err := r.db.Where("user_id = ?", userID).Find(&modelTodos).Error
    if err != nil {
        return nil, err
    }

    // Model → Domain
    var domainTodos []domain.Todo
    for _, m := range modelTodos {
        domainTodos = append(domainTodos, domain.Todo{
            ID:          m.ID,
            UserID:      m.UserID,
            CategoryID:  m.CategoryID,
            Title:       m.Title,
            Description: m.Description,
        })
    }

    return domainTodos, nil
}

func (r *todoRepository) UpdateTodo(todo domain.Todo) error {
    modelTodo := models.Todo{
        ID:          todo.ID,
        UserID:      todo.UserID,
        CategoryID:  todo.CategoryID,
        Title:       todo.Title,
        Description: todo.Description,
    }

    // ID bo'yicha yangilash
    result := r.db.Model(&models.Todo{}).
        Where("id = ?", todo.ID).
        Updates(modelTodo)

    return result.Error
}

func (r *todoRepository) DeleteTodo(id int64) error {
    result := r.db.Delete(&models.Todo{}, id)
    return result.Error
}