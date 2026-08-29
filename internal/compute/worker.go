package compute

// Worker runs a job on a device. Phase 9 will implement a real remote worker
// over the wire; this stub defines the contract used by the Compute Router.
type Worker interface {
	// Execute runs a job and returns a result string.
	Execute(job interface{}) (string, error)
}
