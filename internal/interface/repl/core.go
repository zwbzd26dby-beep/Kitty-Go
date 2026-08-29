package repl

import (
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/budget"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/cost"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/memory"
	reg "github.com/zwbzd26dby-beep/Kitty-Go/internal/model"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/observability"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/security"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/tools"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

// Core bundles the subsystem handles the REPL commands surface.
type Core struct {
	Agent    agent.Agent
	Models   *reg.Registry
	Budget   *budget.Manager
	Cost     *cost.Manager
	Tools    *tools.Registry
	Devices  *compute.Registry
	Memory   memory.ShortTermMemory
	Obs      *observability.Observability
	Security *security.Manager
}

// NewCore builds a Core with the subsystem handles wired up.
func NewCore(client llm.Client) *Core {
	models := reg.NewRegistry()
	_ = reg.RegisterKnown(models)
	budgetMgr := budget.NewManager(0.0, 0.0)
	costMgr := cost.NewManager()
	obs := observability.New(nil)
	toolsReg := tools.NewRegistry()
	tools.RegisterBuiltins(toolsReg)
	devices := compute.NewRegistry()

	exec := execution.NewLocalExecutor()
	orch := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:   exec,
		Client: client,
		Cost:   costMgr,
		Budget: budgetMgr,
		Obs:    obs,
	})
	a := agent.New(orch, agent.DefaultOptions{})

	return &Core{
		Agent:    a,
		Models:   models,
		Budget:   budgetMgr,
		Cost:     costMgr,
		Tools:    toolsReg,
		Devices:  devices,
		Memory:   memory.NewShortTermMemory(64),
		Obs:      obs,
		Security: security.NewManager(nil, nil),
	}
}
