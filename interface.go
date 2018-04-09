package kube

import (
	"context"

	"github.com/sdeoras/configio"
	"k8s.io/client-go/kubernetes"
)

type Kind string

type Coder interface {
	// Kind returns kind of the object
	Kind() Kind
	//Context returns context
	Context() context.Context
	// Init initializes object based on config read from config reader
	Init(clientset *kubernetes.Clientset, configReader configio.ConfigReader) error
	// Create deploys after receiving done signal from input context
	// It will output a context for downstream processes to use
	Create(ctx context.Context) context.Context
	// Delete deletes objects after receiving done signal from input context
	// It will output a context for downstream processes to use
	Delete(ctx context.Context) context.Context
	// Satisfies error interface
	error
}
