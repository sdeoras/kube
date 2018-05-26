// kube defines interface for k8s object deployment API.
// The intent is to keep this interface as small as possible.
// This is work in progress and API may change.
package kube

import (
	"context"

	"github.com/sdeoras/configio"
)

// Coder defines an interface that is primarily to be implemented by
// kubernetes actions such as pv's, pvc's etc.
type Coder interface {
	// Kind returns kind of the object
	Kind() Kind
	// SetConfig configures coder with provided config object
	// It returns type assertion error if the config object type
	// does not match the accepted type of interface implementor
	SetConfig(config configio.Config) error
	// GetConfig retrieves implementor's internal config object
	GetConfig() configio.Config
	//Context returns context of the object implementing this interface
	Context() context.Context
	// Error returns a channel on which internal errors are reported
	Error() <-chan error
	// Create deploys after receiving done signal from input context
	// It will output a context for downstream processes to use
	Create(ctx context.Context) context.Context
	// Delete deletes objects after receiving done signal from input context
	// It will output a context for downstream processes to use
	Delete(ctx context.Context) context.Context
	// Satisfies error interface
}
