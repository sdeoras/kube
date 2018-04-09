package pv

import (
	"context"

	"github.com/sdeoras/kube"
)

func NewCoder(key string, ctx context.Context) kube.Coder {
	return newCoder(key, ctx)
}

func newCoder(key string, ctx context.Context) *coder {
	coder := new(coder)
	coder.key = key
	coder.ctx, coder.cancel = context.WithCancel(ctx)
	return coder
}
