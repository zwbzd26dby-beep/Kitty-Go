package decision

import "strings"

// Rule matches tasks and contributes weights toward a kind/scoring.
type Rule struct {
	Name       string
	Kind       string
	Weight     int
	Match      func(t TaskInfo) bool
	WeightedBy func(t TaskInfo) float64
}

// rules assembles the built-in rule set.
func rules() []Rule {
	return []Rule{
		{
			Name:   "code",
			Kind:   "code",
			Weight: 3,
			Match: func(t TaskInfo) bool {
				return hasAny(t.Input, "code", "coding", "function", "bugfix", "refactor", "test")
			},
			WeightedBy: func(t TaskInfo) float64 {
				return scoreTerms(t.Input, "code", "coding", "function", "debug", "refactor", "test")
			},
		},
		{
			Name:   "math",
			Kind:   "math",
			Weight: 2,
			Match: func(t TaskInfo) bool {
				return hasAny(t.Input, "add", "sum", "multiply", "divide", "equation", "calculate")
			},
			WeightedBy: func(t TaskInfo) float64 {
				return scoreTerms(t.Input, "sum", "multiply", "divide", "equation", "calculate")
			},
		},
		{
			Name:   "vision",
			Kind:   "vision",
			Weight: 2,
			Match: func(t TaskInfo) bool {
				return hasAny(t.Input, "image", "photo", "picture", "vision", "diagram", "screenshot")
			},
			WeightedBy: func(t TaskInfo) float64 {
				return scoreTerms(t.Input, "image", "photo", "picture", "diagram", "screenshot")
			},
		},
	}
}

// PriorityRule scores task priority from urgency/complexity signals.
func prioritySignals(t TaskInfo) (int, float64) {
	score := 1.0
	if hasAny(t.Input, "urgent", "asap", "critical", "immediately", "blocking") {
		score += 2
	}
	if hasAny(t.Input, "complex", "large", "deep", "production", "critical") {
		score += 1
	}
	if len(t.Input) > 200 {
		score += 1
	}
	p := int(score * 10)
	if p > 100 {
		p = 100
	}
	return p, score
}

func hasAny(s string, terms ...string) bool {
	lower := strings.ToLower(s)
	for _, t := range terms {
		if strings.Contains(lower, t) {
			return true
		}
	}
	return false
}

func scoreTerms(s string, terms ...string) float64 {
	lower := strings.ToLower(s)
	n := 0.0
	for _, t := range terms {
		if strings.Contains(lower, t) {
			n++
		}
	}
	return n
}
