package workflow

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

const TaskQueueName = "AGORA_TASK_QUEUE"

type TopicAuditInput struct {
	TopicID int64  `json:"topic_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type TopicAuditOutput struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

// TopicAuditWorkflow 异步审核工作流入口
func TopicAuditWorkflow(ctx workflow.Context, input TopicAuditInput) (TopicAuditOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result TopicAuditOutput
	err := workflow.ExecuteActivity(ctx, AuditActivity, input).Get(ctx, &result)
	if err != nil {
		return TopicAuditOutput{Approved: false, Reason: err.Error()}, err
	}

	return result, nil
}

// AuditActivity 实际执行审核的具体步骤（例如调用 LLM Mock 或接口）
func AuditActivity(ctx context.Context, input TopicAuditInput) (TopicAuditOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Executing AuditActivity for Topic", "TopicID", input.TopicID)

	// Mock 逻辑：直接通过审核
	return TopicAuditOutput{
		Approved: true,
		Reason:   "Content is safe (mock audit passed)",
	}, nil
}
