package kube

import (
	"context"
	"sync"
)

// Fan in returns a context that ends when each of the input contexts are done
func FanIn(ctxs ...context.Context) context.Context {
	var wg sync.WaitGroup
	for _, ctx := range ctxs {
		ctx := ctx

		wg.Add(1)

		// fan in coder context into done
		go func(ctx context.Context, w *sync.WaitGroup) {
			select {
			case <-ctx.Done():
				w.Done()
			}
		}(ctx, &wg)
	}

	done, cancel := context.WithCancel(context.Background())
	go func(c context.CancelFunc, w *sync.WaitGroup) {
		w.Wait()
		c()
	}(cancel, &wg)

	return done
}
