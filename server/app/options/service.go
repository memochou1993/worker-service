package options

// ServiceOptions contains options to configure a Service instance.
type ServiceOptions struct {
	MaxWorkers *int
}

// SetMaxWorkers specifies the maximum amount of workers.
func (s *ServiceOptions) SetMaxWorkers(max int) *ServiceOptions {
	s.MaxWorkers = &max
	return s
}

// Service creates a new ServiceOptions instance.
func Service() *ServiceOptions {
	return &ServiceOptions{
		MaxWorkers: new(int),
	}
}

// MergeServiceOptions combines the given ServiceOptions instances into a single ServiceOptions in a last-one-wins fashion.
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
