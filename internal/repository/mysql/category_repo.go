package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"todolist/internal/domain"
	"todolist/internal/repository/models"
	"todolist/pkg/cache"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewCategoryRepo(db *gorm.DB, cache *cache.Cache) *CategoryRepo {
	return &CategoryRepo{db: db, cache: cache}
}

func (r *CategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	model := &models.Category{
		UserID: category.UserID,
		Name:   category.Name,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Redisga yozish (write-through)
	key := fmt.Sprintf("category:%d:%d", model.UserID, model.ID)
	val, _ := json.Marshal(domain.Category{
		ID:     model.ID,
		UserID: model.UserID,
		Name:   model.Name,
	})
	return r.cache.Set(ctx, key, string(val))
}

func (r *CategoryRepo) GetByID(ctx context.Context, id, userID int64) (*domain.Category, error) {
	key := fmt.Sprintf("category:%d:%d", userID, id)

	// Redis’dan qidirish
	if val, err := r.cache.Get(ctx, key); err == nil {
		var cat domain.Category
		if err := json.Unmarshal([]byte(val), &cat); err == nil {
			return &cat, nil
		}
	}

	// Agar cache’da bo‘lmasa → DB’dan olish
	var model models.Category
	if err := r.db.WithContext(ctx).
		Preload("Todos").
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error; err != nil {
		return nil, err
	}

	var todos []domain.Todo
	for _, t := range model.Todos {
		todos = append(todos, domain.Todo{
			ID:          t.ID,
			UserID:      t.UserID,
			CategoryID:  t.CategoryID,
			Title:       t.Title,
			Description: t.Description,
		})
	}

	cat := &domain.Category{
		ID:     model.ID,
		UserID: model.UserID,
		Name:   model.Name,
		Todos:  todos,
	}

	// Redisga saqlash
	val, _ := json.Marshal(cat)
	_ = r.cache.Set(ctx, key, string(val))

	return cat, nil
}

func (r *CategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	result := r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("id = ? AND user_id = ?", category.ID, category.UserID).
		Updates(map[string]interface{}{"name": category.Name})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("kategoriya topilmadi yoki ruxsat yo'q")
	}

	// Redisni yangilash
	key := fmt.Sprintf("category:%d:%d", category.UserID, category.ID)
	val, _ := json.Marshal(category)
	return r.cache.Set(ctx, key, string(val))
}

func (r *CategoryRepo) Delete(ctx context.Context, id, userID int64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Category{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("kategoriya topilmadi yoki ruxsat yo'q")
	}

	// Redisdan o‘chirish
	key := fmt.Sprintf("category:%d:%d", userID, id)
	return r.cache.Delete(ctx, key)
}

func (r *CategoryRepo) GetAllByUserID(ctx context.Context, userID int64) ([]*domain.Category, error) {
	key := fmt.Sprintf("categories:user:%d", userID)

	if val, err := r.cache.Get(ctx, key); err == nil {
		var cats []*domain.Category
		if err := json.Unmarshal([]byte(val), &cats); err == nil {
			return cats, nil
		}
	}

	// Agar cache’da bo‘lmasa DB’dan olish
	var modelsCat []models.Category
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&modelsCat).Error; err != nil {
		return nil, err
	}

	var result []*domain.Category
	for _, m := range modelsCat {
		result = append(result, &domain.Category{
			ID:     m.ID,
			UserID: m.UserID,
			Name:   m.Name,
		})
	}

	val, _ := json.Marshal(result)
	_ = r.cache.Set(ctx, key, string(val))

	return result, nil
}
