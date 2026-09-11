package workflow

import (
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Register is the single registration point shared by the development and
// production workers. Governance workflows are added here as they are built.
func Register(w worker.Worker) {
	w.RegisterWorkflow(SystemHealthWorkflow)
}

// SystemHealthWorkflow gives deployment smoke tests a deterministic way to
// verify that the Worker can receive and execute a Temporal workflow.
func SystemHealthWorkflow(_ workflow.Context) (string, error) {
	return "ok", nil
}
