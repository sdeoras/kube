package kube

import (
	"context"
)

// Cleanup sends delete commands to all coders
func Cleanup(coders ...Coder) context.Context {
	trigger, startFunc := context.WithCancel(context.Background())
	trigger, _ = Delete(trigger, Async, coders...)
	go startFunc()
	return trigger
}
