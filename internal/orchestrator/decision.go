package orchestrator

// Decision is a placeholder carrying task routing intent in Phase 1.
// The full Decision type is produced by the Decision Engine in Phase 16.
type Decision struct {
	TaskID string
	Kind   string
}
