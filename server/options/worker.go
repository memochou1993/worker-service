package options

// WorkerOptions 工人選項
type WorkerOptions struct {
	MaxDelay *int
}

// SetMaxDelay 設置工人最大延遲時間
func (s *WorkerOptions) SetMaxDelay(max int) *WorkerOptions {
	s.MaxDelay = &max
	return s
}

// Worker 工人選項
func Worker() *WorkerOptions {
	return &WorkerOptions{
		MaxDelay: new(int),
	}
}

// MergeWorkerOptions 合併工人選項
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
