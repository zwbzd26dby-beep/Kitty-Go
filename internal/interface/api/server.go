// Package api implements the HTTP API surface bound to the shared Core
// (GET /health, GET /metrics, POST /agent/respond) with Bearer auth.
package api

import (
	"compress/gzip"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/observability"
)

// Server serves the persistent HTTP API over the shared Core.
type Server struct {
	agent      agent.Agent
	obs        *observability.Observability
	token      string
	mux        *http.ServeMux
	httpServer *http.Server
}

// NewServer builds an API server. An empty token disables auth.
func NewServer(a agent.Agent, obs *observability.Observability, token string) *Server {
	s := &Server{
		agent: a,
		obs:   obs,
		token: token,
		mux:   http.NewServeMux(),
	}
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/metrics", s.handleMetrics)
	s.mux.HandleFunc("/agent/respond", s.withAuth(s.handleRespond))
	return s
}

// Handler returns the underlying http.Handler (useful for tests).
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe starts the server on addr. An error is returned when the
// process stops unexpectedly.
func (s *Server) ListenAndServe(addr string) error {
	s.httpServer = &http.Server{Addr: addr, Handler: s.mux}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.token == "" {
			next(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) ||
			subtle.ConstantTimeCompare([]byte(auth[len(prefix):]), []byte(s.token)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

type respondReq struct {
	Message string `json:"message"`
}

type respondResp struct {
	Content string `json:"content"`
	Intent  string `json:"intent,omitempty"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

func (s *Server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	out := map[string]any{}
	if s.obs != nil {
		cs, gs := s.obs.Metrics.Snapshot()
		out["counters"] = cs
		out["gauges"] = gs
		out["uptime_seconds"] = s.obs.Uptime()
	} else {
		out["counters"] = map[string]int64{}
		out["gauges"] = map[string]float64{}
	}
	_ = json.NewEncoder(w).Encode(out)
}

func (s *Server) handleRespond(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var reader io.Reader = r.Body
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "bad gzip body", http.StatusBadRequest)
			return
		}
		defer gz.Close()
		reader = gz
	}
	var req respondReq
	if err := json.NewDecoder(reader).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("bad request: %v", err), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}
	resp, err := s.agent.Process(r.Context(), req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if s.obs != nil {
		s.obs.Metrics.Inc("api_requests")
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(respondResp{Content: resp.Content, Intent: resp.CurrentIntent})
}
