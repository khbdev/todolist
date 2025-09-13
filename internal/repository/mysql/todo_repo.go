// mysql/todo_repository.go
package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"todolist/internal/domain"
	"todolist/internal/repository/models"
	"todolist/pkg/cache"
	"todolist/pkg/notifactions"

	"gorm.io/gorm"
)

type todoRepository struct {
	db    *gorm.DB
	cache *cache.Cache
	notifier *notifactions.Notifier
}

// NewTodoRepository - DB ham, cache ham qabul qiladi
func NewTodoRepository(db *gorm.DB, cache *cache.Cache, notifier *notifactions.Notifier) domain.TodoRepository {
	return &todoRepository{
		db:    db,
		cache: cache,
		notifier: notifier,
		
	}
}

// === CreateTodo ===
// Write-through: DB ga yozilgandan keyin cache tozalanadi
// CreateTodo
func (r *todoRepository) CreateTodo(todo domain.Todo) (int64, error) {
	modelTodo := models.Todo{
		UserID:      todo.UserID,
		CategoryID:  todo.CategoryID,
		Title:       todo.Title,
		Description: todo.Description,
	}

	result := r.db.Create(&modelTodo)
	if result.Error != nil {
		return 0, result.Error
	}

	// ✅ YANGI: Agar category cache'langan bo'lsa — TOZALASH
	ctx := context.Background()
	categoryKey := fmt.Sprintf("category:%d:%d", todo.UserID, todo.CategoryID)
	r.cache.Delete(ctx, categoryKey)

	// boshqa cachelarni ham tozalash
	r.cache.Delete(ctx, fmt.Sprintf("todo:%d", modelTodo.ID))
	r.cache.Delete(ctx, fmt.Sprintf("todos:user:%d", todo.UserID))

		_ = r.notifier.Publish(context.Background(), "created", todo)

	return modelTodo.ID, nil
}

// === GetTodoByID ===
// Read-through: avval cache, keyin DB
func (r *todoRepository) GetTodoByID(id int64) (*domain.Todo, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("todo:%d", id)

	// 1. Cache'dan o'qish
	cached, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		var todo domain.Todo
		if json.Unmarshal([]byte(cached), &todo) == nil {
			return &todo, nil
		}
		// Agar unmarshal muvaffaqiyatsiz bo'lsa, DB dan davom etamiz
	}

	// 2. Cache yo'q — DB dan olish
	var modelTodo models.Todo
	if err := r.db.First(&modelTodo, id).Error; err != nil {
		return nil, err
	}

	// 3. Domain modelga o'tkazish
	domainTodo := &domain.Todo{
		ID:          modelTodo.ID,
		UserID:      modelTodo.UserID,
		CategoryID:  modelTodo.CategoryID,
		Title:       modelTodo.Title,
		Description: modelTodo.Description,
	}

	// 4. Cache'ga saqlash
	if data, marshalErr := json.Marshal(domainTodo); marshalErr == nil {
		r.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return domainTodo, nil
}

// === GetTodosByUserID ===
// Read-through
func (r *todoRepository) GetTodosByUserID(userID int64) ([]domain.Todo, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("todos:user:%d", userID)

	// 1. Cache'dan o'qish
	cached, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		var todos []domain.Todo
		if json.Unmarshal([]byte(cached), &todos) == nil {
			return todos, nil
		}
	}

	// 2. DB dan olish
	var modelTodos []models.Todo
	if err := r.db.Where("user_id = ?", userID).Find(&modelTodos).Error; err != nil {
		return nil, err
	}

	// 3. Domainga o'tkazish
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

	// 4. Cache'ga saqlash
	if data, marshalErr := json.Marshal(domainTodos); marshalErr == nil {
		r.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return domainTodos, nil
}

// === UpdateTodo ===
// Write-through: DB yangilanadi, cache tozalanadi
func (r *todoRepository) UpdateTodo(todo domain.Todo) error {
	modelTodo := models.Todo{
		ID:          todo.ID,
		UserID:      todo.UserID,
		CategoryID:  todo.CategoryID,
		Title:       todo.Title,
		Description: todo.Description,
	}

	result := r.db.Model(&models.Todo{}).
		Where("id = ?", todo.ID).
		Updates(modelTodo)

	if result.Error != nil {
		return result.Error
	}

	// Cache tozalash
	ctx := context.Background()
	r.cache.Delete(ctx, fmt.Sprintf("category:%d:%d", todo.UserID, todo.CategoryID))
	r.cache.Delete(ctx, fmt.Sprintf("todo:%d", todo.ID))
	r.cache.Delete(ctx, fmt.Sprintf("todos:user:%d", todo.UserID))

	_ = r.notifier.Publish(context.Background(), "updated", todo)
	return nil
}

// === DeleteTodo ===
// Write-through: DB o'chiriladi, cache tozalanadi
// DeleteTodo
func (r *todoRepository) DeleteTodo(id int64) error {
	// Avval ma'lumotni olish
	var modelTodo models.Todo
	if err := r.db.First(&modelTodo, id).Error; err != nil {
		return err
	}

	// O'chirish
	if err := r.db.Delete(&models.Todo{}, id).Error; err != nil {
		return err
	}

	// ✅ Category cache TOZALANADI
	ctx := context.Background()
	r.cache.Delete(ctx, fmt.Sprintf("category:%d:%d", modelTodo.UserID, modelTodo.CategoryID))
	r.cache.Delete(ctx, fmt.Sprintf("todo:%d", id))
	r.cache.Delete(ctx, fmt.Sprintf("todos:user:%d", modelTodo.UserID))

	_ = r.notifier.Publish(context.Background(), "deleted", map[string]any{"id": id})

	return nil
}