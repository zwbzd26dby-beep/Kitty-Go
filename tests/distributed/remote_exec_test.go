package distributedtest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/compute"
	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
)

// newRemoteWorker spins up a fake worker answering /execute and /health.
func newRemoteWorker(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/health":
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
		case "/execute":
			var job execution.Job
			_ = json.NewDecoder(r.Body).Decode(&job)
			if job.ID == "" {
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(execution.JobResult{Success: false, Error: "missing id"})
				return
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(execution.JobResult{
				JobID:   job.ID,
				Content: "worker processed: " + job.Prompt,
				Success: true,
			})
		default:
			w.WriteHeader(404)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func device(addr, token string) compute.Device {
	return compute.Device{ID: "node-1", Address: trimHTTP(addr), AuthToken: token}
}

func trimHTTP(addr string) string {
	return addr[len("http://"):]
}

func TestRemoteExecute(t *testing.T) {
	srv := newRemoteWorker(t)
	exec := execution.NewRemoteExecutor(compute.NewHTTPTransport())
	res, err := exec.Execute(context.Background(), device(srv.URL, ""), execution.Job{
		ID: "j1", Model: "llama3", Prompt: "hello",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.Success || res.Content == "" {
		t.Fatalf("unexpected result %+v", res)
	}
}

func TestRemoteExecuteErrorPayload(t *testing.T) {
	exec := execution.NewRemoteExecutor(compute.NewHTTPTransport())
	_, err := exec.Execute(context.Background(), device("127.0.0.1:1", ""), execution.Job{ID: "j2", Prompt: "x"})
	if err == nil {
		t.Fatal("expected error for unreachable worker")
	}
}

func TestRemotePing(t *testing.T) {
	srv := newRemoteWorker(t)
	exec := execution.NewRemoteExecutor(compute.NewHTTPTransport())
	if err := exec.Ping(context.Background(), device(srv.URL, "")); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestRemoteAuthenticate(t *testing.T) {
	exec := execution.NewRemoteExecutor(compute.NewHTTPTransport())
	dev := device("127.0.0.1:1", "secret-token")
	if err := exec.Authenticate(context.Background(), dev, "secret-token"); err != nil {
		t.Fatalf("expected auth success, got %v", err)
	}
	if err := exec.Authenticate(context.Background(), dev, "wrong"); err == nil {
		t.Fatal("expected auth failure")
	}
}

func TestRemoteStreamSingleChunk(t *testing.T) {
	srv := newRemoteWorker(t)
	exec := execution.NewRemoteExecutor(compute.NewHTTPTransport())
	ch, err := exec.Stream(context.Background(), device(srv.URL, ""), execution.Job{ID: "j3", Prompt: "hi"})
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	count := 0
	for range ch {
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 chunk, got %d", count)
	}
}
