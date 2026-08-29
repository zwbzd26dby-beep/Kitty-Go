package execution

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// FallbackDevice couples a remote executor with the device it targets.
type FallbackDevice struct {
	Device compute.Device
	Exec   DistributedExecutor
}

// ExecuteOptions configures the full execution flow.
type ExecuteOptions struct {
	// Client drives the local path (required for local execution).
	Client LLMClient
	// Model selected for the job.
	Model string
	// MaxRetries per target before failing over.
	MaxRetries int
	// Backoff base between local retries.
	Backoff time.Duration
	// Fallbacks are ordered remote targets tried after the local path.
	Fallbacks []FallbackDevice
}

// FullExecutor implements the complete execution flow from Phase 11
// (Master Architecture §36): retry with backoff, then fallback to remote
// devices, then error.
type FullExecutor struct {
	jobs *jobTracker
}

// NewFullExecutor creates a FullExecutor.
func NewFullExecutor() *FullExecutor {
	return &FullExecutor{jobs: newJobTracker()}
}

// Execute runs task through the unified execution flow.
func (f *FullExecutor) Execute(ctx context.Context, task types.Task, opts ExecuteOptions) (*Result, error) {
	job := Job{ID: task.ID, Model: opts.Model, Prompt: task.Input}
	f.jobs.put(job.ID, &JobStatus{JobID: job.ID, State: JobRunning, StartedAt: time.Now()})

	var runners []Runner
	if opts.Client != nil {
		runners = append(runners, &retryRunner{
			base:    localRunner{client: opts.Client, history: task.History},
			max:     opts.MaxRetries,
			backoff: opts.Backoff,
		})
	}
	for _, fb := range opts.Fallbacks {
		runners = append(runners, remoteRunner{exec: fb.Exec, device: fb.Device})
	}
	if len(runners) == 0 {
		f.jobs.put(job.ID, &JobStatus{JobID: job.ID, State: JobFailed})
		return nil, fmt.Errorf("execution failed: no target available")
	}

	var errs []string
	for _, r := range runners {
		res, err := r.Run(ctx, job)
		if err == nil {
			f.jobs.put(job.ID, &JobStatus{JobID: job.ID, State: JobCompleted, StartedAt: time.Now(), EndedAt: time.Now()})
			return res, nil
		}
		errs = append(errs, err.Error())
	}

	f.jobs.put(job.ID, &JobStatus{JobID: job.ID, State: JobFailed, StartedAt: time.Now(), EndedAt: time.Now()})
	return nil, fmt.Errorf("execution failed on all targets: %s", strings.Join(errs, "; "))
}

// ExecuteLocal delegates to the Phase 1 local executor.
func (f *FullExecutor) ExecuteLocal(ctx context.Context, task types.Task, client LLMClient) (*Result, error) {
	le := NewLocalExecutor()
	return le.ExecuteLocal(ctx, task, client)
}

// GetStatus returns the current status of a job.
func (f *FullExecutor) GetStatus(jobID string) (*JobStatus, error) {
	return f.jobs.get(jobID)
}

// Cancel cancels a running job.
func (f *FullExecutor) Cancel(jobID string) error {
	return f.jobs.cancel(jobID)
}

var _ Executor = (*FullExecutor)(nil)