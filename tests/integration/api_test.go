package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/interface/api"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/observability"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/orchestrator"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/llm"
)

func apiStack() (agent.Agent, *observability.Observability) {
	exec := execution.NewLocalExecutor()
	obs := observability.New(nil)
	orch := orchestrator.NewFull(orchestrator.FullOptions{
		Exec:   exec,
		Client: llm.NewHistoryEchoClient(),
		Obs:    obs,
	})
	return agent.New(orch, agent.DefaultOptions{}), obs
}

func TestAPIHealth(t *testing.T) {
	a, obs := apiStack()
	srv := api.NewServer(a, obs, "")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"ok"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestAPIMetrics(t *testing.T) {
	a, obs := apiStack()
	obs.Metrics.Inc("tasks_total")
	srv := api.NewServer(a, obs, "")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "tasks_total") {
		t.Fatalf("expected tasks_total in metrics, got %s", rec.Body.String())
	}
}

func TestAPIRespondOK(t *testing.T) {
	a, obs := apiStack()
	srv := api.NewServer(a, obs, "")
	body := `{"message":"Witaj Kitty"}`
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/agent/respond", bytes.NewReader([]byte(body))))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Content string
		Intent  string
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Content, "mock działa") {
		t.Fatalf("unexpected content: %q", resp.Content)
	}
}

func TestAPIRespondBadMethod(t *testing.T) {
	a, obs := apiStack()
	srv := api.NewServer(a, obs, "")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/agent/respond", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestAPIRespondEmptyMessage(t *testing.T) {
	a, obs := apiStack()
	srv := api.NewServer(a, obs, "")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/agent/respond", bytes.NewReader([]byte(`{"message":"  "}`))))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestAPIBearerAuth(t *testing.T) {
	a, obs := apiStack()
	srv := api.NewServer(a, obs, "secret-token")

	// Without token → 401.
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/agent/respond", bytes.NewReader([]byte(`{"message":"hi"}`))))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", rec.Code)
	}

	// Wrong token → 401.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/agent/respond", bytes.NewReader([]byte(`{"message":"hi"}`)))
	req.Header.Set("Authorization", "Bearer wrong")
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with wrong token, got %d", rec.Code)
	}

	// Correct token → 200.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/agent/respond", bytes.NewReader([]byte(`{"message":"hi"}`)))
	req.Header.Set("Authorization", "Bearer secret-token")
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with correct token, got %d", rec.Code)
	}
}
