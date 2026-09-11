package workflow

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Register is the single registration point shared by the development and
// production workers. Governance workflows are added here as they are built.
func Register(w worker.Worker, activities *Activities) {
	w.RegisterWorkflow(SystemHealthWorkflow)
	w.RegisterWorkflow(ContentCoolingWorkflow)
	w.RegisterWorkflow(FeedbackAuditWorkflow)
	w.RegisterWorkflow(ClusterTopicWorkflow)
	w.RegisterWorkflow(BlindReviewWorkflow)
	w.RegisterWorkflow(ReviewAuditWorkflow)
	w.RegisterWorkflow(FinalizeBlindReviewWorkflow)
	registerActivities(w, activities)
	w.RegisterActivityWithOptions(activities.AuditFeedback, activity.RegisterOptions{Name: "AuditFeedback"})
	w.RegisterActivityWithOptions(activities.ClusterTopic, activity.RegisterOptions{Name: "ClusterTopic"})
	w.RegisterActivityWithOptions(activities.CreateBlindReviewBatch, activity.RegisterOptions{Name: "CreateBlindReviewBatch"})
	w.RegisterActivityWithOptions(activities.FinalizeBlindReview, activity.RegisterOptions{Name: "FinalizeBlindReview"})
	w.RegisterActivityWithOptions(activities.AuditReview, activity.RegisterOptions{Name: "AuditReview"})
}

// SystemHealthWorkflow gives deployment smoke tests a deterministic way to
// verify that the Worker can receive and execute a Temporal workflow.
func SystemHealthWorkflow(_ workflow.Context) (string, error) {
	return "ok", nil
}
