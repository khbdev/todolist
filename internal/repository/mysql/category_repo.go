package mysql

import (
    "context"
    "errors"
    "todolist/internal/domain"
    "todolist/internal/repository/models"

    "gorm.io/gorm"
)

type CategoryRepo struct {
    db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
    return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, category *domain.Category) error {
    model := &models.Category{
        UserID: category.UserID,
        Name:   category.Name,
    }
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *CategoryRepo) GetByID(ctx context.Context, id, userID int64) (*domain.Category, error) {
    var model models.Category

    if err := r.db.WithContext(ctx).
        Preload("Todos"). // bu yerda todos larni ham yuklab oladi
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

    if result.Error != nil { return result.Error }
    if result.RowsAffected == 0 {
        return errors.New("kategoriya topilmadi yoki ruxsat yo'q")
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
    return nil
}

func (r *CategoryRepo) GetAllByUserID(ctx context.Context, userID int64) ([]*domain.Category, error) {
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
    return result, nil
}