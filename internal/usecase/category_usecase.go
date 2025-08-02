package usecase

import (
    "context"
    "todolist/internal/domain"
)

type CategoryUsecase struct {
    repo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) *CategoryUsecase {
    return &CategoryUsecase{repo: repo}
}

func (uc *CategoryUsecase) Create(ctx context.Context, category *domain.Category) error {
    return uc.repo.Create(ctx, category)
}

func (uc *CategoryUsecase) GetByID(ctx context.Context, userID, id int64) (*domain.Category, error) {
    return uc.repo.GetByID(ctx, id, userID)
}

func (uc *CategoryUsecase) Update(ctx context.Context, category *domain.Category) error {
    return uc.repo.Update(ctx, category)
}

func (uc *CategoryUsecase) Delete(ctx context.Context, userID, id int64) error {
    return uc.repo.Delete(ctx, id, userID)
}

func (uc *CategoryUsecase) GetAllByUserID(ctx context.Context, userID int64) ([]*domain.Category, error) {
    return uc.repo.GetAllByUserID(ctx, userID)
}