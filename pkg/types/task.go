package types

// Task represents a unit of work handed to the orchestrator.
// In Phase 0 it is a minimal carrier; later phases extend it with
// requirements, priority, budget and deadline.
type Task struct {
	ID      string
	Input   string
	Session string
}

// TaskResult is the outcome of executing a task.
type TaskResult struct {
	TaskID  string
	Content string
	Error   error
}
