package execution

import "time"

// JobState is the lifecycle state of a job.
type JobState string

// Job lifecycle states.
const (
	JobPending   JobState = "pending"
	JobRunning   JobState = "running"
	JobCompleted JobState = "completed"
	JobFailed    JobState = "failed"
	JobCancelled JobState = "cancelled"
)

// JobStatus describes the current state of a job.
type JobStatus struct {
	JobID     string
	State     JobState
	StartedAt time.Time
	EndedAt   time.Time
}
