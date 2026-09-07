package service

import (
	"context"

	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	wf "agora-backend/internal/workflow"

	"go.temporal.io/sdk/client"
)

type CreateTopicReq struct {
	CategoryID int32  `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required,min=3,max=255"`
	Content    string `json:"content" binding:"required,min=5"`
}

type TopicService struct {
	topicDAO       *dao.TopicDAO
	temporalClient client.Client
}

func NewTopicService(topicDAO *dao.TopicDAO, temporalClient client.Client) *TopicService {
	return &TopicService{
		topicDAO:       topicDAO,
		temporalClient: temporalClient,
	}
}

func (s *TopicService) CreateTopic(ctx context.Context, authorID int64, req *CreateTopicReq) (*model.Topic, error) {
	topic := &model.Topic{
		CategoryID:        req.CategoryID,
		AuthorID:          authorID,
		Title:             req.Title,
		Content:           req.Content,
		StructuredContent: "{}",
		Status:            "published",
	}

	// 写入数据库
	if err := s.topicDAO.Create(ctx, topic); err != nil {
		return nil, err
	}

	// 触发 Temporal 异步审核工作流
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: wf.TaskQueueName,
	}
	auditInput := wf.TopicAuditInput{
		TopicID: topic.ID,
		Title:   topic.Title,
		Content: topic.Content,
	}

	_, err := s.temporalClient.ExecuteWorkflow(ctx, workflowOptions, wf.TopicAuditWorkflow, auditInput)
	if err != nil {
		// 注意：此处仅记录日志或报警，不影响主流程响应
	}

	return topic, nil
}
