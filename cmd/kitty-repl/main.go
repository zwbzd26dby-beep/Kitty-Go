// Command kitty-repl is the interactive Kitty REPL entry point (Phase 0:
// history-echo mock provider).
package main

import (
	"os"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/repl"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func main() {
	engine := repl.NewEngine(llm.NewHistoryEchoClient(), os.Stdin, os.Stdout)
	if err := engine.Run(); err != nil {
		os.Exit(1)
	}
}
