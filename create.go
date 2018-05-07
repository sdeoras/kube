package kube

import (
	"context"
	"sync"
)

// Create creates coders sequentially or in an async way
func Create(ctx context.Context, syncType SyncType, coders ...Coder) (context.Context, error) {
	var trigger context.Context
	trigger = ctx

	switch syncType {
	case Sync:
		for i := 0; i < len(coders); i++ {
			trigger = coders[i].Create(trigger)
		}
		return trigger, nil
	case Async:
		var wg sync.WaitGroup
		for _, coder := range coders {
			coder := coder
			wg.Add(1)
			go func(w *sync.WaitGroup, trig context.Context, cdr Coder) {
				select {
				case <-cdr.Create(trig).Done():
					w.Done()
				}
			}(&wg, ctx, coder)
		}

		done, cancel := context.WithCancel(context.Background())
		go func(c context.CancelFunc, w *sync.WaitGroup) {
			w.Wait()
			c()
		}(cancel, &wg)
		return done, nil
	default:
		return nil, UnsupportedSync
	}
}
