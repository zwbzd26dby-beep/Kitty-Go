package compute

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Transport is a generic byte/JSON transport to a remote worker. It is kept
// independent of the execution domain to avoid an import cycle between
// compute and execution; execution encodes its Job payloads to JSON.
type Transport interface {
	// SendJSON posts payload to addr+path and returns the response body.
	SendJSON(ctx context.Context, addr, path string, payload []byte) ([]byte, error)
	// Ping verifies a worker is reachable.
	Ping(ctx context.Context, addr string) error
}

// HTTPTransport is a Transport over net/http with pluggable client.
type HTTPTransport struct {
	Client *http.Client
}

// NewHTTPTransport creates an HTTPTransport.
func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{Client: &http.Client{Timeout: 30 * time.Second}}
}

// SendJSON posts payload to addr+path and returns the response body.
func (t *HTTPTransport) SendJSON(ctx context.Context, addr, path string, payload []byte) ([]byte, error) {
	url := "http://" + addr + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := t.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("remote worker status %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// Ping verifies a worker is reachable via a GET to its root.
func (t *HTTPTransport) Ping(ctx context.Context, addr string) error {
	url := "http://" + addr + "/health"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := t.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("worker health status %d", resp.StatusCode)
	}
	return nil
}

var _ Transport = (*HTTPTransport)(nil)
