package kube

import (
	"context"
	"sync"
)

func Cleanup(ctx context.Context, coders ...Coder) <-chan struct{} {
	// shutdown all if error occurs
	var wg sync.WaitGroup
	for _, coder := range coders {
		coder := coder
		go func(coder Coder) {
			wg.Add(1)

			delete := coder.Delete(ctx)
			select {
			case <-coder.Error():
			case <-delete.Done(): // if no error, block till delete
			}

			wg.Done()
		}(coder)
	}

	cleanup := make(chan struct{})
	go func() {
		wg.Wait()
		cleanup <- struct{}{}
	}()

	return cleanup
}
