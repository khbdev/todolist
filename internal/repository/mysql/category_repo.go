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

    cat := &domain.Category{
        ID:     model.ID,
        UserID: model.UserID,
        Name:   model.Name,
    }

    // 🔹 Individual cache
    key := fmt.Sprintf("category:%d:%d", model.UserID, model.ID)
    val, _ := json.Marshal(cat)
    _ = r.cache.Set(ctx, key, string(val))

    // 🔹 List cache update
    listKey := fmt.Sprintf("categories:user:%d", model.UserID)
    if val, err := r.cache.Get(ctx, listKey); err == nil {
        var cats []*domain.Category
        if json.Unmarshal([]byte(val), &cats) == nil {
            cats = append(cats, cat) // yangi category qo‘shamiz
            newVal, _ := json.Marshal(cats)
            _ = r.cache.Set(ctx, listKey, string(newVal))
        }
    }

    return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id, userID int64) (*domain.Category, error) {
	var model models.Category
	if err := r.db.WithContext(ctx).
		Select("id", "user_id", "name").
		Preload("Todos", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id", "category_id", "title", "description")
		}).
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

	return &domain.Category{
		ID:     model.ID,
		UserID: model.UserID,
		Name:   model.Name,
		Todos:  todos,
	}, nil
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

    // 🔹 Individual cache update
    key := fmt.Sprintf("category:%d:%d", category.UserID, category.ID)
    val, _ := json.Marshal(category)
    _ = r.cache.Set(ctx, key, string(val))

    // 🔹 List cache update
    listKey := fmt.Sprintf("categories:user:%d", category.UserID)
    if val, err := r.cache.Get(ctx, listKey); err == nil {
        var cats []*domain.Category
        if json.Unmarshal([]byte(val), &cats) == nil {
            for i, c := range cats {
                if c.ID == category.ID {
                    cats[i] = category // yangilash
                    break
                }
            }
            newVal, _ := json.Marshal(cats)
            _ = r.cache.Set(ctx, listKey, string(newVal))
        }
    }

    return nil
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

    // 🔹 Individual cache o‘chir
    key := fmt.Sprintf("category:%d:%d", userID, id)
    _ = r.cache.Delete(ctx, key)

    // 🔹 List cache update (o‘chirish)
    listKey := fmt.Sprintf("categories:user:%d", userID)
    if val, err := r.cache.Get(ctx, listKey); err == nil {
        var cats []*domain.Category
        if json.Unmarshal([]byte(val), &cats) == nil {
            var newCats []*domain.Category
            for _, c := range cats {
                if c.ID != id {
                    newCats = append(newCats, c)
                }
            }
            newVal, _ := json.Marshal(newCats)
            _ = r.cache.Set(ctx, listKey, string(newVal))
        }
    }

    return nil
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
