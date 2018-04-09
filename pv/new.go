package pv

import (
	"context"

	"github.com/sdeoras/scratchpad/kube"
)

func NewCoder(key string, ctx context.Context) kube.Coder {
	return newCoder(key, ctx)
}

func newCoder(key string, ctx context.Context) *manager {
	coder := new(manager)
	coder.key = key
	coder.ctx, coder.cancel = context.WithCancel(ctx)
	return coder
}
