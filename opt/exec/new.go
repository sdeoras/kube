// exec implements kube.Coder interface for executing a shell command
package exec

import (
	"context"
	"io"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
)

func NewCoder(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, key, cmd string, args ...string) (kube.Coder, error) {
	return newCoder(ctx, stdin, stdout, stderr, key, cmd, args...)
}

func newCoder(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, key, cmd string, args ...string) (*coder, error) {
	cdr := new(coder)
	cdr.key = key
	cdr.config = new(Config).Init(key, stdin, stdout, stderr, cmd, args...)
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)
	cdr.log = logrus.WithField("coder", PackageName)
	cdr.err = make(chan error)

	return cdr, nil
}
