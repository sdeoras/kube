// pods implements kube.Coder interface for deployment of pods
package noop

import (
	"context"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
)

func NewCoder(name string, ctx context.Context) (kube.Coder, error) {
	return newCoder(name, ctx)
}

func newCoder(name string, ctx context.Context) (*coder, error) {
	cdr := new(coder)
	cdr.key = ""
	cdr.config = new(Config).Init(name)
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)
	cdr.log = logrus.WithField("coder", "noop")
	cdr.err = make(chan error)

	return cdr, nil
}
