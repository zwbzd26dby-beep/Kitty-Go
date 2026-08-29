package types

// Task represents a unit of work handed to the orchestrator.
// In Phase 1 it carries the input, session and conversation history.
type Task struct {
	ID      string
	Input   string
	Session string
	History []Turn
}

// TaskResult is the outcome of executing a task.
type TaskResult struct {
	TaskID  string
	Content string
	Error   error
}
