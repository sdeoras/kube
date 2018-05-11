// nl implements kube.Coder interface for nodelist
package nl

import (
	"context"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
)

func NewCoder(ctx context.Context, configReader configio.ConfigReader, key string) (kube.Coder, error) {
	return newCoder(ctx, configReader, key)
}

func newCoder(ctx context.Context, configReader configio.ConfigReader, key string) (*coder, error) {
	cdr := new(coder)
	cdr.key = key
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)

	cdr.err = make(chan error)
	cdr.config = new(Config).Init(key)

	return cdr, nil
}
