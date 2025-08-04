package usecase

import (
    "todolist/internal/domain"
)

type TodoUsecase struct {
    repo domain.TodoRepository
}

func NewTodoUsecase(repo domain.TodoRepository) *TodoUsecase {
    return &TodoUsecase{repo: repo}
}

func (uc *TodoUsecase) CreateTodo(todo domain.Todo) (int64, error) {
    return uc.repo.CreateTodo(todo)
}

func (uc *TodoUsecase) GetTodoByID(id int64) (*domain.Todo, error) {
    return uc.repo.GetTodoByID(id)
}

func (uc *TodoUsecase) GetTodosByUserID(userID int64) ([]domain.Todo, error) {
    return uc.repo.GetTodosByUserID(userID)
}

func (uc *TodoUsecase) UpdateTodo(todo domain.Todo) error {
    return uc.repo.UpdateTodo(todo)
}

func (uc *TodoUsecase) DeleteTodo(id int64) error {
    return uc.repo.DeleteTodo(id)
}
