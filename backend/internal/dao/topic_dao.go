package dao

import (
	"context"

	"agora-backend/internal/model"
	"gorm.io/gorm"
)

type TopicDAO struct {
	db *gorm.DB
}

func NewTopicDAO(db *gorm.DB) *TopicDAO {
	return &TopicDAO{db: db}
}

func (d *TopicDAO) Create(ctx context.Context, topic *model.Topic) error {
	return d.db.WithContext(ctx).Create(topic).Error
}

func (d *TopicDAO) GetByID(ctx context.Context, id int64) (*model.Topic, error) {
	var topic model.Topic
	if err := d.db.WithContext(ctx).First(&topic, id).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}
