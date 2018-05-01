// sleep implements kube.Coder interface for sleeping (this is a kube no op)
package sleep

import (
	"context"

	"time"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
)

func NewCoder(ctx context.Context, name string, dur time.Duration) (kube.Coder, error) {
	return newCoder(ctx, name, dur)
}

func newCoder(ctx context.Context, name string, dur time.Duration) (*coder, error) {
	cdr := new(coder)
	cdr.key = ""
	cdr.config = new(Config).Init(name, dur)
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)
	cdr.log = logrus.WithField("coder", PackageName)
	cdr.err = make(chan error)

	return cdr, nil
}
