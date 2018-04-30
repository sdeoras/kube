// pods implements kube.Coder interface for deployment of pods
package sleep

import (
	"context"

	"time"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
)

func NewCoder(name string, dur time.Duration, ctx context.Context) (kube.Coder, error) {
	return newCoder(name, dur, ctx)
}

func newCoder(name string, dur time.Duration, ctx context.Context) (*coder, error) {
	cdr := new(coder)
	cdr.key = ""
	cdr.config = new(Config).Init(name, dur)
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)
	cdr.log = logrus.WithField("coder", PackageName)
	cdr.err = make(chan error)

	return cdr, nil
}
