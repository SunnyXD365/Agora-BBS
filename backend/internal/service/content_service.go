package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
)

var ErrInvalidContentQuery = errors.New("invalid content query")

type ContentService struct {
	dao *dao.ContentDAO
}

func NewContentService(contentDAO *dao.ContentDAO) *ContentService {
	return &ContentService{dao: contentDAO}
}

func normalizePublicPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize != 10 && pageSize != 20 && pageSize != 50 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *ContentService) ListMine(ctx context.Context, userID int64, contentType, status string, page, pageSize int) (*model.Page[*model.MyContentItem], error) {
	if contentType == "" {
		contentType = "topic"
	}
	if contentType != "topic" && contentType != "post" {
		return nil, ErrInvalidContentQuery
	}
	allowedStatuses := map[string]bool{"": true, "draft": true, "cooling": true, "pending_review": true, "published": true, "rejected": true, "recalled": true, "hidden": true}
	if !allowedStatuses[status] {
		return nil, ErrInvalidContentQuery
	}
	page, pageSize = normalizePublicPage(page, pageSize)
	items, total, err := s.dao.ListMine(ctx, userID, contentType, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &model.Page[*model.MyContentItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *ContentService) Search(ctx context.Context, query, contentType string, page, pageSize int) (*model.Page[*model.SearchResult], error) {
	query = strings.TrimSpace(query)
	if utf8.RuneCountInString(query) < 2 || utf8.RuneCountInString(query) > 100 {
		return nil, ErrInvalidContentQuery
	}
	if contentType != "" && contentType != "topic" && contentType != "post" {
		return nil, ErrInvalidContentQuery
	}
	page, pageSize = normalizePublicPage(page, pageSize)
	pattern := "%" + escapeLike(query) + "%"
	items, total, err := s.dao.Search(ctx, pattern, contentType, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &model.Page[*model.SearchResult]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "%", `\%`)
	return strings.ReplaceAll(value, "_", `\_`)
}
