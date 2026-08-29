package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
)

// Job is a serialisable unit of work sent to a remote worker
// (Master Architecture §17).
type Job struct {
	ID        string `json:"id"`
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens,omitempty"`
}

// JobResult is the serialisable outcome of a remote Job.
type JobResult struct {
	JobID   string        `json:"job_id"`
	Content string        `json:"content"`
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
	Elapsed time.Duration `json:"elapsed"`
}

// DistributedExecutor executes jobs on remote devices
// (Master Architecture §17).
type DistributedExecutor interface {
	// Execute runs a job on device and waits for the result.
	Execute(ctx context.Context, device compute.Device, job Job) (*JobResult, error)
	// Stream runs a job on device and returns a channel of partial results.
	Stream(ctx context.Context, device compute.Device, job Job) (<-chan *JobResult, error)
	// Ping verifies device is reachable.
	Ping(ctx context.Context, device compute.Device) error
	// Authenticate verifies device credentials (token).
	Authenticate(ctx context.Context, device compute.Device, token string) error
}

// RemoteExecutor is a DistributedExecutor over a compute.Transport.
type RemoteExecutor struct {
	transport compute.Transport
}

// NewRemoteExecutor builds a RemoteExecutor over transport.
func NewRemoteExecutor(transport compute.Transport) *RemoteExecutor {
	return &RemoteExecutor{transport: transport}
}

// Execute posts the job and decodes the result.
func (r *RemoteExecutor) Execute(ctx context.Context, device compute.Device, job Job) (*JobResult, error) {
	payload, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}
	body, err := r.transport.SendJSON(ctx, device.Address, "/execute", payload)
	if err != nil {
		return nil, err
	}
	var res JobResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("decode job result: %w", err)
	}
	if !res.Success {
		return &res, fmt.Errorf("remote job failed: %s", res.Error)
	}
	return &res, nil
}

// Stream is not implemented for remote in Phase 9 (single-shot).
func (r *RemoteExecutor) Stream(ctx context.Context, device compute.Device, job Job) (<-chan *JobResult, error) {
	ch := make(chan *JobResult, 1)
	res, err := r.Execute(ctx, device, job)
	if err != nil {
		close(ch)
		return nil, err
	}
	ch <- res
	close(ch)
	return ch, nil
}

// Ping delegates to the transport.
func (r *RemoteExecutor) Ping(ctx context.Context, device compute.Device) error {
	return r.transport.Ping(ctx, device.Address)
}

// Authenticate verifies the device token matches.
func (r *RemoteExecutor) Authenticate(ctx context.Context, device compute.Device, token string) error {
	if device.AuthToken == "" {
		return fmt.Errorf("device %q has no auth token configured", device.ID)
	}
	if device.AuthToken != token {
		return fmt.Errorf("authentication failed for device %q", device.ID)
	}
	return nil
}

var _ DistributedExecutor = (*RemoteExecutor)(nil)
