package service

import (
	"context"

	"Agora-BBS/internal/dao"
	"Agora-BBS/internal/model"
)

type CategoryService struct {
	categoryDAO *dao.CategoryDAO
}

func NewCategoryService(categoryDAO *dao.CategoryDAO) *CategoryService {
	return &CategoryService{categoryDAO: categoryDAO}
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*model.Category, error) {
	return s.categoryDAO.ListAll(ctx)
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *model.CreateCategoryReq) (*model.Category, error) {
	c := &model.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
	}
	if err := s.categoryDAO.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func CategoryServiceFactory(categoryDAO *dao.CategoryDAO) *CategoryService {
	return &CategoryService{categoryDAO: categoryDAO}
}
