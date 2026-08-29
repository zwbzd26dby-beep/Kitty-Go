// Command kitty-repl is the interactive Kitty REPL entry point.
//
// Since Phase 1 it runs the layered stack explicitly:
//
//	REPL → Agent → Orchestrator → Executor(Local) → provider
//
// The provider used here is the history-echo mock.
package main

import (
	"os"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/repl"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func main() {
	client := llm.NewHistoryEchoClient()
	exec := execution.NewLocalExecutor()
	orch := orchestrator.New(exec, client)
	a := agent.New(orch, agent.DefaultOptions{})
	engine := repl.NewEngineWithAgent(a, os.Stdin, os.Stdout)
	if err := engine.Run(); err != nil {
		os.Exit(1)
	}
}
