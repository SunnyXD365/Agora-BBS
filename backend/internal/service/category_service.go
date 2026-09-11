package service

import (
	"context"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
)

type CategoryService struct {
	categoryDAO *dao.CategoryDAO
}

func NewCategoryService(categoryDAO *dao.CategoryDAO) *CategoryService {
	return &CategoryService{categoryDAO: categoryDAO}
}

// ListCategories 获取所有板块列表
func (s *CategoryService) ListCategories(ctx context.Context) ([]*model.Category, error) {
	return s.categoryDAO.ListCategories(ctx)
}
