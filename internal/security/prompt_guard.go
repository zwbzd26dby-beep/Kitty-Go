package security

import "strings"

// InjectionAdvice describes a detected prompt-injection risk.
type InjectionAdvice struct {
	Detected  bool
	RiskLevel string // "low" | "medium" | "high"
	Reason    string
}

// PromptGuard scans untrusted content (documents, tool output) for common
// prompt-injection patterns.
type PromptGuard struct {
	patterns []string
}

// NewPromptGuard creates a PromptGuard with the default patterns.
func NewPromptGuard() *PromptGuard {
	return &PromptGuard{patterns: []string{
		"ignore previous instructions",
		"ignore all previous",
		"disregard prior",
		"new system prompt",
		"you are now",
		"act as system",
		"reveal your system prompt",
		"forget everything above",
	}}
}

// Inspect returns advice for the given untrusted text.
func (g *PromptGuard) Inspect(text string) InjectionAdvice {
	lower := strings.ToLower(text)
	var hits []string
	for _, p := range g.patterns {
		if strings.Contains(lower, p) {
			hits = append(hits, p)
		}
	}
	if len(hits) == 0 {
		return InjectionAdvice{Detected: false, RiskLevel: "low"}
	}
	level := "medium"
	if len(hits) >= 2 {
		level = "high"
	}
	return InjectionAdvice{
		Detected:  true,
		RiskLevel: level,
		Reason:    "matched: " + strings.Join(hits, ", "),
	}
}

// IsSafe reports whether a document is safe to include verbatim.
func (g *PromptGuard) IsSafe(text string) bool {
	return !g.Inspect(text).Detected
}
