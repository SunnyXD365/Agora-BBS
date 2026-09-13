package service

import (
	"context"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
)

type BookmarkService struct{ dao *dao.BookmarkDAO }

func NewBookmarkService(bookmarkDAO *dao.BookmarkDAO) *BookmarkService {
	return &BookmarkService{dao: bookmarkDAO}
}

func (s *BookmarkService) Create(ctx context.Context, userID, topicID int64) error {
	return s.dao.Create(ctx, userID, topicID)
}

func (s *BookmarkService) Delete(ctx context.Context, userID, topicID int64) error {
	return s.dao.Delete(ctx, userID, topicID)
}

func (s *BookmarkService) List(ctx context.Context, userID int64, page, pageSize int) (*model.Page[*model.Bookmark], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	items, total, err := s.dao.List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &model.Page[*model.Bookmark]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}
