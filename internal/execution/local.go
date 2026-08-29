package execution

import (
	"context"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// localExecutor is the Phase 1 synchronous local execution path.
type localExecutor struct {
	jobs *jobTracker
}

// NewLocalExecutor creates a local executor.
func NewLocalExecutor() *localExecutor {
	return &localExecutor{jobs: newJobTracker()}
}

// ExecuteLocal generates a response for the task via the provided client.
func (l *localExecutor) ExecuteLocal(ctx context.Context, task types.Task, client LLMClient) (*Result, error) {
	msg, err := types.NewMessage(task.Input)
	if err != nil {
		return nil, err
	}
	resp, err := client.Generate(ctx, msg, task.History)
	if err != nil {
		return nil, err
	}
	return &Result{
		TaskID:  task.ID,
		Content: resp.Content(),
	}, nil
}

// GetStatus returns job status or an error if the job is unknown.
func (l *localExecutor) GetStatus(jobID string) (*JobStatus, error) {
	return l.jobs.get(jobID)
}

// Cancel cancels a job if it is tracked.
func (l *localExecutor) Cancel(jobID string) error {
	return l.jobs.cancel(jobID)
}
