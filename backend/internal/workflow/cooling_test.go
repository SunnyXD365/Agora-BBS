package workflow

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
)

func TestContentCoolingWorkflowRunsPublishActivity(t *testing.T) {
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterActivityWithOptions(func(context.Context, CoolingInput) (string, error) {
		return "published", nil
	}, activity.RegisterOptions{Name: "PublishContent"})
	env.ExecuteWorkflow(ContentCoolingWorkflow, CoolingInput{Kind: "topic", ID: 7, EndsAt: time.Time{}})
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result string
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, "published", result)
}
