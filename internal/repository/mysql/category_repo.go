package mysql

import (
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

func (r *CategoryRepo) Create(category *domain.Category) error {
	model := &models.Category{
		UserID: category.UserID,
		Name:   category.Name,
	}
	return r.db.Create(model).Error
}

func (r *CategoryRepo) GetByID(id int64) (*domain.Category, error) {
	var model models.Category
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &domain.Category{
		ID:     model.ID,
		UserID: model.UserID,
		Name:   model.Name,
	}, nil
}

func (r *CategoryRepo) Update(category *domain.Category) error {
	return r.db.Model(&models.Category{}).
		Where("id = ?", category.ID).
		Updates(map[string]interface{}{
			"name": category.Name,
		}).Error
}

func (r *CategoryRepo) Delete(id int64) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *CategoryRepo) GetAllByUserID(userID int64) ([]*domain.Category, error) {
	var modelsCat []models.Category
	if err := r.db.Where("user_id = ?", userID).Find(&modelsCat).Error; err != nil {
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
