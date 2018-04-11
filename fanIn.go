package kube

// Fan in collects coder contexts and errors and funnels them through
// single endpoints
func FanIn(coders ...Coder) (<-chan struct{}, <-chan error) {
	done := make(chan struct{})
	err := make(chan error)
	for _, coder := range coders {
		coder := coder

		// fan in coder context into done
		go func(cdr Coder) {
			done <- <-cdr.Context().Done()
		}(coder)

		// fan in coder error into errChan
		go func(cdr Coder) {
			err <- <-cdr.Error()
		}(coder)
	}

	return done, err
}
