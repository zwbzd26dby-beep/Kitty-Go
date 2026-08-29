// Command kitty-gui serves a minimal web chat on the shared Core.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/gui"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func main() {
	addr := getenv("KITTY_GUI_ADDR", "127.0.0.1:8081")

	client := llm.NewHistoryEchoClient()
	exec := execution.NewLocalExecutor()
	orch := orchestrator.New(exec, client)
	a := agent.New(orch, agent.DefaultOptions{})

	srv := gui.NewServer(a)
	log.Printf("kitty-gui listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
