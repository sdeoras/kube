package kube

import "context"

// Bootup sequentially boots a list of coders
func Bootup(ctx context.Context, order Order, coders ...Coder) context.Context {
	var trigger context.Context
	trigger = ctx

	switch order {
	case OrderForward:
		for i := 0; i < len(coders); i++ {
			trigger = coders[i].Create(trigger)
		}
	case OrderBackward:
		for i := len(coders) - 1; i >= 0; i-- {
			trigger = coders[i].Create(trigger)
		}
	}

	return trigger
}
