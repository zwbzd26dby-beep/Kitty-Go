// Command kitty-repl is the interactive Kitty REPL entry point.
//
// It runs the layered stack explicitly:
//
//	REPL → Agent → Orchestrator → Executor(Local) → provider
//
// The provider used here is the history-echo mock.
package main

import (
	"os"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/repl"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func main() {
	core := repl.NewCore(llm.NewHistoryEchoClient())
	engine := repl.NewEngineWithCore(core, os.Stdin, os.Stdout)
	if err := engine.Run(); err != nil {
		os.Exit(1)
	}
}
