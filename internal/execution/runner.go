package execution

import (
	"context"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/utils"
)

// Runner unifies local and remote execution paths behind one interface
// (Master Architecture §36).
type Runner interface {
	Run(ctx context.Context, job Job) (*Result, error)
}

// retryRunner wraps a Runner with retry + exponential backoff.
type retryRunner struct {
	base    Runner
	max     int
	backoff time.Duration
}

// Run executes base with retry/backoff until success or MaxAttempts.
func (r *retryRunner) Run(ctx context.Context, job Job) (*Result, error) {
	max := r.max
	if max <= 0 {
		max = 1
	}
	backoff := r.backoff
	if backoff <= 0 {
		backoff = 100 * time.Millisecond
	}
	var res *Result
	err := utils.RetryWith(ctx, utils.RetryPolicy{
		MaxAttempts: max,
		Initial:     backoff,
		MaxBackoff:  backoff * 4,
		Multiplier:  2,
	}, nil, func() error {
		var e error
		res, e = r.base.Run(ctx, job)
		return e
	})
	return res, err
}

// localRunner bridges an LLMClient to the Runner contract.
type localRunner struct {
	client  LLMClient
	history []types.Turn
}

// Run generates a completion via the local client.
func (l localRunner) Run(ctx context.Context, job Job) (*Result, error) {
	msg, err := types.NewMessage(job.Prompt)
	if err != nil {
		return nil, err
	}
	resp, err := l.client.Generate(ctx, msg, l.history)
	if err != nil {
		return nil, err
	}
	return &Result{TaskID: job.ID, Content: resp.Content()}, nil
}

// remoteRunner bridges a DistributedExecutor + Device to the Runner contract.
type remoteRunner struct {
	exec   DistributedExecutor
	device compute.Device
}

// Run executes the job on the remote device.
func (r remoteRunner) Run(ctx context.Context, job Job) (*Result, error) {
	jr, err := r.exec.Execute(ctx, r.device, job)
	if err != nil {
		return nil, err
	}
	return &Result{TaskID: job.ID, Content: jr.Content}, nil
}
