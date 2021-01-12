package options

// WorkerOptions contains options to configure a Worker instance.
type WorkerOptions struct {
	MaxDelay *int
}

// SetMaxDelay specifies the maximum amount of delay.
func (s *WorkerOptions) SetMaxDelay(max int) *WorkerOptions {
	s.MaxDelay = &max
	return s
}

// Worker creates a new WorkerOptions instance.
func Worker() *WorkerOptions {
	return &WorkerOptions{
		MaxDelay: new(int),
	}
}

// MergeWorkerOptions combines the given WorkerOptions instances into a single WorkerOptions in a last-one-wins fashion.
func MergeWorkerOptions(opts ...*WorkerOptions) *WorkerOptions {
	wOpts := Worker()
	for _, wo := range opts {
		if wo == nil {
			continue
		}
		if wo.MaxDelay != nil {
			wOpts.MaxDelay = wo.MaxDelay
		}
	}
	return wOpts
}
