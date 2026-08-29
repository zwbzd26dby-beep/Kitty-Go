// Package decision implements the Decision Engine that analyses a task into
// routing and execution requirements (Master Architecture §7).
package decision

// Requirement names the capabilities a task may need.
type Requirement struct {
	Name        string
	Capability  string
	Description string
}

// TaskInfo is the input to the Decision Engine.
type TaskInfo struct {
	ID       string
	Input    string
	Category string // user-provided hint: "chat" | "code" | "vision" | ""
	Budget   float64
}

// Decision is the engine's structured analysis of a task.
type Decision struct {
	TaskID       string
	Kind         string
	Requirements []Requirement
	Priority     int
	Budget       float64
}

// HasRequirement reports whether a requirement is present.
func (d Decision) HasRequirement(name string) bool {
	for _, r := range d.Requirements {
		if r.Name == name {
			return true
		}
	}
	return false
}

// NeedsCapability reports whether the decision requires a model capability.
func (d Decision) NeedsCapability(cap string) bool {
	for _, r := range d.Requirements {
		if r.Capability == cap {
			return true
		}
	}
	return false
}
