package orchestrator

// Coordinator is responsible for chaining the orchestration stages. In
// Phase 1 the coordination logic lives entirely in simpleOrchestrator; this
// type documents the interface that will be expanded in later phases.
type Coordinator interface {
	// Coordinate runs the full orchestration pipeline for a task.
	Coordinate(task interface{}) (interface{}, error)
}
