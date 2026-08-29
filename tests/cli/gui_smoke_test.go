package clitest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/gui"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func guiStack() agent.Agent {
	exec := execution.NewLocalExecutor()
	orch := orchestrator.New(exec, llm.NewHistoryEchoClient())
	return agent.New(orch, agent.DefaultOptions{})
}

func TestGUIServesPage(t *testing.T) {
	srv := gui.NewServer(guiStack())
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Kitty GUI") {
		t.Fatalf("expected page title, got: %s", rec.Body.String())
	}
}

func TestGUIPostResponse(t *testing.T) {
	srv := gui.NewServer(guiStack())
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/agent/respond", strings.NewReader(`{"message":"hej"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "mock działa") {
		t.Fatalf("expected mock response, got: %s", rec.Body.String())
	}
}

func TestGUIEmptyMessageBadRequest(t *testing.T) {
	srv := gui.NewServer(guiStack())
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/agent/respond", strings.NewReader(`{"message":""}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
