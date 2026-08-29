// Command kitty-api serves the persistent HTTP API on the shared Core.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/api"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/observability"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
)

func main() {
	host := getenv("KITTY_HOST", "127.0.0.1")
	port := getenv("KITTY_PORT", "8080")
	token := os.Getenv("KITTY_API_TOKEN")

	client := llm.NewHistoryEchoClient()
	exec := execution.NewLocalExecutor()
	obs := observability.New(os.Stderr)
	orch := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:   exec,
		Client: client,
		Obs:    obs,
	})
	a := agent.New(orch, agent.DefaultOptions{})

	srv := api.NewServer(a, obs, token)
	addr := host + ":" + port
	go func() {
		log.Printf("kitty-api listening on %s", addr)
		if err := srv.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
