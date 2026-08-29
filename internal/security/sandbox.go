package security

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// SandboxPolicy restricts what commands a tool wrapper may run.
type SandboxPolicy struct {
	// Allow lists executables the sandbox permits (bare names or abs paths).
	Allow []string
	// Deny blocks matching executables even if in Allow.
	Deny []string
	// Timeout bounds exec lifetime.
	Timeout time.Duration
	// MaxOutput caps captured output bytes.
	MaxOutput int
}

// DefaultSandboxPolicy allows only read-only inspection commands.
func DefaultSandboxPolicy() SandboxPolicy {
	return SandboxPolicy{
		Allow:     []string{"ls", "cat", "pwd", "wc", "sha256sum"},
		Deny:      []string{"rm", "sh", "bash", "curl"},
		Timeout:   5 * time.Second,
		MaxOutput: 1 << 16,
	}
}

// Sandbox runs commands safely.
type Sandbox struct {
	policy SandboxPolicy
}

// NewSandbox creates a Sandbox with the given policy.
func NewSandbox(policy SandboxPolicy) *Sandbox {
	return &Sandbox{policy: policy}
}

// ErrCommandDenied is returned for disallowed commands.
var ErrCommandDenied = fmt.Errorf("command denied by sandbox policy")

// Run executes cmd with args (no shell) under the policy, returning
// captured stdout or an error.
func (s *Sandbox) Run(ctx context.Context, cmd string, args ...string) (string, error) {
	if !s.allowed(cmd) {
		return "", ErrCommandDenied
	}
	timeout := s.policy.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	c := exec.CommandContext(ctx, cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}
	if s.policy.MaxOutput > 0 && len(out) > s.policy.MaxOutput {
		return "", fmt.Errorf("output exceeded %d bytes", s.policy.MaxOutput)
	}
	return string(out), nil
}

func (s *Sandbox) allowed(cmd string) bool {
	base := cmd
	if strings.Contains(cmd, "/") {
		base = cmd[strings.LastIndex(cmd, "/")+1:]
	}
	for _, d := range s.policy.Deny {
		if d == base || d == cmd {
			return false
		}
	}
	for _, a := range s.policy.Allow {
		if a == base || a == cmd {
			return true
		}
	}
	return false
}
