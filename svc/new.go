package svc

import (
	"context"

	"github.com/sdeoras/kube"
)

func NewCoder(key, namespace string, ctx context.Context) kube.Coder {
	return newCoder(key, namespace, ctx)
}

func newCoder(key, namespace string, ctx context.Context) *coder {
	coder := new(coder)
	coder.key = key
	coder.namespace = namespace
	coder.ctx, coder.cancel = context.WithCancel(ctx)
	return coder
}
