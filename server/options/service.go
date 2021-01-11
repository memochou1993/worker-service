package options

// ServiceOptions 服務選項
type ServiceOptions struct {
	MaxWorkers *int
}

// SetMaxWorkers 設置服務最大工人數量
func (s *ServiceOptions) SetMaxWorkers(max int) *ServiceOptions {
	s.MaxWorkers = &max
	return s
}

// Service 服務選項
func Service() *ServiceOptions {
	return &ServiceOptions{
		MaxWorkers: new(int),
	}
}

// MergeServiceOptions 合併服務選項
func MergeServiceOptions(opts ...*ServiceOptions) *ServiceOptions {
	sOpts := Service()
	for _, so := range opts {
		if so == nil {
			continue
		}
		if so.MaxWorkers != nil {
			sOpts.MaxWorkers = so.MaxWorkers
		}
	}
	return sOpts
}
