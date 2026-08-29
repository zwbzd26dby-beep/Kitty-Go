package decision

import "errors"

// ErrNoRuleMatched is returned when no rule matched the task.
var ErrNoRuleMatched = errors.New("no decision rule matched the task")

// Engine evaluates rules for a task and produces a Decision.
type Engine struct {
	rules []Rule
}

// New creates an Engine with the built-in rule set.
func New() *Engine {
	return &Engine{rules: rules()}
}

// Analyze produces a Decision for the task.
func (e *Engine) Analyze(t TaskInfo) (Decision, error) {
	best := -1
	var bestRule *Rule
	for i := range e.rules {
		r := &e.rules[i]
		if r.Match(t) {
			if bestRule == nil || r.Weight > best || (r.Weight == best && weighted(r, t) > weighted(bestRule, t)) {
				best = r.Weight
				bestRule = r
			}
		}
	}
	if bestRule == nil {
		return Decision{}, ErrNoRuleMatched
	}

	kind := bestRule.Kind
	if t.Category != "" {
		kind = t.Category
	}

	d := Decision{
		TaskID:   t.ID,
		Kind:     kind,
		Budget:   t.Budget,
		Priority: priorityFrom(kind, t),
	}
	d.Requirements = requirementsFor(kind)

	if d.Priority > 100 {
		d.Priority = 100
	}
	return d, nil
}

// AnalyzeOrDefault degrades to a default "chat" decision when no rule matches,
// so the orchestrator always has a usable Decision.
func (e *Engine) AnalyzeOrDefault(t TaskInfo) Decision {
	d, err := e.Analyze(t)
	if err == nil {
		return d
	}
	return Decision{
		TaskID:       t.ID,
		Kind:         "chat",
		Budget:       t.Budget,
		Priority:     priorityFrom("chat", t),
		Requirements: requirementsFor("chat"),
	}
}

func weighted(r *Rule, t TaskInfo) float64 {
	if r.WeightedBy != nil {
		return r.WeightedBy(t)
	}
	return 0
}

func priorityFrom(kind string, t TaskInfo) int {
	p, _ := prioritySignals(t)
	return p
}
