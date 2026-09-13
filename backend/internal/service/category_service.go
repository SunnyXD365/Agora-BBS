package service

import (
	"context"
	"time"

	"agora-backend/internal/cache"
	"agora-backend/internal/dao"
	"agora-backend/internal/model"
)

type CategoryService struct {
	categoryDAO *dao.CategoryDAO
	cache       *cache.Store
}

func NewCategoryService(categoryDAO *dao.CategoryDAO, store *cache.Store) *CategoryService {
	return &CategoryService{categoryDAO: categoryDAO, cache: store}
}

// ListCategories 获取所有板块列表
func (s *CategoryService) ListCategories(ctx context.Context) ([]*model.Category, error) {
	items := make([]*model.Category, 0)
	if err := s.cache.GetJSON(ctx, "categories:active", &items); err == nil {
		return items, nil
	}
	items, err := s.categoryDAO.ListCategories(ctx)
	if err == nil {
		_ = s.cache.SetJSON(ctx, "categories:active", items, 5*time.Minute)
	}
	return items, err
}
