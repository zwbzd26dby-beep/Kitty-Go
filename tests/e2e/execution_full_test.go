package e2etest

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// flakyClient fails the first n calls, then succeeds.
type flakyClient struct {
	fails   int
	calls   int
	content string
}

func (f *flakyClient) Generate(ctx context.Context, _ types.Message, _ []types.Turn) (types.Response, error) {
	f.calls++
	if f.calls <= f.fails {
		return types.Response{}, errors.New("transient failure")
	}
	return types.NewResponse(f.content)
}

// fakeRemoteExec simulates a remote worker.
type fakeRemoteExec struct {
	fail bool
	done string
}

func (r *fakeRemoteExec) Execute(_ context.Context, d compute.Device, job execution.Job) (*execution.JobResult, error) {
	if r.fail {
		return nil, errors.New("remote down")
	}
	return &execution.JobResult{JobID: job.ID, Content: "remote[" + r.done + "]", Success: true}, nil
}
func (r *fakeRemoteExec) Stream(_ context.Context, d compute.Device, job execution.Job) (<-chan *execution.JobResult, error) {
	ch := make(chan *execution.JobResult, 1)
	ch <- &execution.JobResult{JobID: job.ID, Content: "remote", Success: true}
	close(ch)
	return ch, nil
}
func (r *fakeRemoteExec) Ping(context.Context, compute.Device) error { return nil }
func (r *fakeRemoteExec) Authenticate(context.Context, compute.Device, string) error {
	return nil
}

func task(input string) types.Task {
	return types.Task{ID: "t1", Input: input}
}

func TestFullExecutionLocalSuccess(t *testing.T) {
	fe := execution.NewFullExecutor()
	res, err := fe.Execute(context.Background(), task("hi"), execution.ExecuteOptions{
		Client:     &flakyClient{content: "hello"},
		MaxRetries: 1,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.Content != "hello" {
		t.Fatalf("unexpected content %q", res.Content)
	}
}

func TestFullExecutionRetryThenSuccess(t *testing.T) {
	fe := execution.NewFullExecutor()
	fc := &flakyClient{fails: 2, content: "recovered"}
	res, err := fe.Execute(context.Background(), task("hi"), execution.ExecuteOptions{
		Client:     fc,
		MaxRetries: 4,
		Backoff:    time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.Content != "recovered" {
		t.Fatalf("unexpected content %q", res.Content)
	}
	if fc.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", fc.calls)
	}
}

func TestFullExecutionLocalFailsFallbackRemote(t *testing.T) {
	fe := execution.NewFullExecutor()
	dev := compute.Device{ID: "node-1", Address: "127.0.0.1:9000"}
	res, err := fe.Execute(context.Background(), task("hi"), execution.ExecuteOptions{
		Client:     &flakyClient{fails: 10, content: "never"},
		MaxRetries: 2,
		Backoff:    time.Millisecond,
		Fallbacks: []execution.FallbackDevice{
			{Device: dev, Exec: &fakeRemoteExec{done: "a"}},
		},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.HasPrefix(res.Content, "remote[") {
		t.Fatalf("expected remote fallback content, got %q", res.Content)
	}
}

func TestFullExecutionTotalFailure(t *testing.T) {
	fe := execution.NewFullExecutor()
	dev := compute.Device{ID: "node-1", Address: "127.0.0.1:9000"}
	_, err := fe.Execute(context.Background(), task("hi"), execution.ExecuteOptions{
		Client:     &flakyClient{fails: 10, content: "never"},
		MaxRetries: 2,
		Backoff:    time.Millisecond,
		Fallbacks: []execution.FallbackDevice{
			{Device: dev, Exec: &fakeRemoteExec{fail: true}},
		},
	})
	if err == nil {
		t.Fatal("expected total failure")
	}
	if !strings.Contains(err.Error(), "all targets") {
		t.Fatalf("expected aggregated error, got %v", err)
	}
}

func TestFullExecutionNoTarget(t *testing.T) {
	fe := execution.NewFullExecutor()
	_, err := fe.Execute(context.Background(), task("hi"), execution.ExecuteOptions{})
	if err == nil {
		t.Fatal("expected error with no targets")
	}
}

func TestFullExecutionSatisfiesExecutorInterface(t *testing.T) {
	fe := execution.NewFullExecutor()
	var _ execution.Executor = fe
	st, err := fe.GetStatus("missing")
	if err == nil {
		t.Fatal("expected error for missing job")
	}
	_ = st
}
