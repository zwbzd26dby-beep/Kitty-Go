// Package orchestrator coordinates a task through decision, model routing,
// compute routing and execution (Master Architecture §32).
package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/budget"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/cost"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/decision"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/limit"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/observability"
	compute "github.com/zwbzd26dby-beep/Kitty-Go/internal/router/compute"
	modelrouter "github.com/zwbzd26dby-beep/Kitty-Go/internal/router/model"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// ErrEmptyTask is returned when a task has no input.
var ErrEmptyTask = errors.New("task input is empty")

// Orchestrator coordinates the execution of a task.
type Orchestrator interface {
	// Orchestrate runs a task and returns its result.
	Orchestrate(ctx context.Context, task types.Task) (*execution.Result, error)
}

// FullOptions wires the complete orchestration pipeline.
type FullOptions struct {
	Exec       execution.Executor
	Client     execution.LLMClient
	Decision   *decision.Engine
	Models     modelrouter.ModelRouter
	Compute    compute.ComputeRouter
	Cost       cost.CostManager
	Limits     *limit.Manager
	Budget     *budget.Manager
	Distributed execution.DistributedExecutor
	// Fallbacks are explicit remote targets tried after the local path.
	Fallbacks []execution.FallbackDevice
	// MaxRetries per target before failing over.
	MaxRetries int
	// Backoff base between retries.
	Backoff time.Duration
	// Obs optionally receives metrics and traces.
	Obs *observability.Observability
}

// pipelinedOrchestrator is the full Phase 17 implementation.
type pipelinedOrchestrator struct {
	exec        execution.Executor
	client      execution.LLMClient
	engine      *decision.Engine
	models      modelrouter.ModelRouter
	compute     compute.ComputeRouter
	cost        cost.CostManager
	limits      *limit.Manager
	budget      *budget.Manager
	distributed execution.DistributedExecutor
	fallbacks   []execution.FallbackDevice
	maxRetries  int
	backoff     time.Duration
	obs         *observability.Observability
}

// NewFull creates a full pipeline orchestrator.
func NewFull(opts FullOptions) Orchestrator {
	engine := opts.Decision
	if engine == nil {
		engine = decision.New()
	}
	return &pipelinedOrchestrator{
		exec:        opts.Exec,
		client:      opts.Client,
		engine:      engine,
		models:      opts.Models,
		compute:     opts.Compute,
		cost:        opts.Cost,
		limits:      opts.Limits,
		budget:      opts.Budget,
		distributed: opts.Distributed,
		fallbacks:   opts.Fallbacks,
		maxRetries:  opts.MaxRetries,
		backoff:     opts.Backoff,
		obs:         opts.Obs,
	}
}

// New creates a pass-through orchestrator (Phase 1 contract, kept for tests
// and simple entry points).
func New(exec execution.Executor, client execution.LLMClient) Orchestrator {
	return &pipelinedOrchestrator{exec: exec, client: client, engine: decision.New()}
}

// Orchestrate runs the full pipeline.
func (o *pipelinedOrchestrator) Orchestrate(ctx context.Context, task types.Task) (*execution.Result, error) {
	if strings.TrimSpace(task.Input) == "" {
		return nil, ErrEmptyTask
	}
	if o.obs != nil {
		defer o.obs.End()
		o.obs.Start("orchestrate")
		o.obs.Metrics.Inc("tasks_total")
	}

	dec := o.engine.AnalyzeOrDefault(decision.TaskInfo{
		ID:    task.ID,
		Input: task.Input,
	})
	modelName := o.routeModel(ctx, task, dec)
	if o.obs != nil {
		o.obs.Metrics.Inc("model_routed")
	}

	estTokens := estimateTokens(task.Input) + 2048
	estCost := o.estimateCost(dec, modelName, task.Input)

	if o.limits != nil {
		if err := o.limits.CheckAndTrack(uint64(estTokens)); err != nil {
			return nil, fmt.Errorf("rate limit: %w", err)
		}
	}
	if o.budget != nil {
		if err := o.budget.Check(estCost); err != nil {
			return nil, fmt.Errorf("budget: %w", err)
		}
	}

	fallbacks := append([]execution.FallbackDevice(nil), o.fallbacks...)
	if fb, ok := o.deviceFallback(ctx, dec); ok {
		fallbacks = append(fallbacks, fb)
	}

	res, err := o.execute(ctx, task, modelName, fallbacks)
	if err != nil {
		return nil, fmt.Errorf("orchestrate: %w", err)
	}

	if o.cost != nil {
		_, _ = o.cost.TrackCost(cost.Record{
			Model: modelName,
			Cost:  estCost,
		})
	}
	if o.budget != nil {
		_ = o.budget.Spend(estCost)
	}
	if o.obs != nil {
		o.obs.Metrics.Inc("tasks_completed")
	}
	return res, nil
}

func (o *pipelinedOrchestrator) routeModel(ctx context.Context, task types.Task, dec decision.Decision) string {
	if o.models != nil {
		req := modelrouter.RouteRequest{
			Task: task.Input,
			Decision: modelrouter.Decision{
				RequiresCodeGeneration: dec.NeedsCapability("CodeGeneration"),
				MaxCost:                dec.Budget,
			},
		}
		if m, err := o.models.Select(req); err == nil {
			return m.Model
		}
	}
	return ""
}

func (o *pipelinedOrchestrator) estimateCost(dec decision.Decision, modelName, input string) float64 {
	// Char-count proxy; refined pricing arrives when model metadata carries it.
	return float64(len(input)) * 0.00003 * float64(dec.Priority)/100
}

// fullExecutor is the superset of Executor offered by the unified
// execution layer (Phase 11).
type fullExecutor interface {
	execution.Executor
	Execute(ctx context.Context, task types.Task, opts execution.ExecuteOptions) (*execution.Result, error)
}

func (o *pipelinedOrchestrator) execute(ctx context.Context, task types.Task, modelName string, fallbacks []execution.FallbackDevice) (*execution.Result, error) {
	if fe, ok := o.exec.(fullExecutor); ok {
		return fe.Execute(ctx, task, execution.ExecuteOptions{
			Client:     o.client,
			Model:      modelName,
			MaxRetries: o.maxRetries,
			Backoff:    o.backoff,
			Fallbacks:  fallbacks,
		})
	}
	return o.exec.ExecuteLocal(ctx, task, o.client)
}

// deviceFallback selects a compute device and wraps it as a remote fallback
// target when a distributed executor is configured.
func (o *pipelinedOrchestrator) deviceFallback(ctx context.Context, dec decision.Decision) (execution.FallbackDevice, bool) {
	if o.compute == nil || o.distributed == nil {
		return execution.FallbackDevice{}, false
	}
	dev, err := o.compute.Select(compute.Request{})
	if err != nil {
		return execution.FallbackDevice{}, false
	}
	_ = ctx
	return execution.FallbackDevice{Device: dev, Exec: o.distributed}, true
}

// SimpleCoordinator documents the interface that chains pipeline stages.
type SimpleCoordinator interface {
	Coordinate(task interface{}) (interface{}, error)
}

func estimateTokens(s string) int {
	return len(s) / 4
}