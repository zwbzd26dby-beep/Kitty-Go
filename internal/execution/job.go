package execution

import (
	"fmt"
	"sync"
)

// JobID is the type of a job identifier.
type JobID string

// jobTracker is a minimal in-memory store of job state. Remote/distributed
// tracking arrives in Phase 9.
type jobTracker struct {
	mu   sync.Mutex
	jobs map[JobID]*JobStatus
}

func newJobTracker() *jobTracker {
	return &jobTracker{jobs: map[JobID]*JobStatus{}}
}

func (t *jobTracker) put(jobID string, st *JobStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.jobs[JobID(jobID)] = st
}

func (t *jobTracker) get(jobID string) (*JobStatus, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	st, ok := t.jobs[JobID(jobID)]
	if !ok {
		return nil, fmt.Errorf("job %q not found", jobID)
	}
	return st, nil
}

func (t *jobTracker) cancel(jobID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	st, ok := t.jobs[JobID(jobID)]
	if !ok {
		return fmt.Errorf("job %q not found", jobID)
	}
	st.State = JobCancelled
	return nil
}
