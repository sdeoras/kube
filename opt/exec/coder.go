package exec

import (
	"context"
	"os/exec"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
)

// coder implements kube.Coder interface
type coder struct {
	key    string
	config *Config
	ctx    context.Context
	cancel context.CancelFunc
	err    chan error
	log    *logrus.Entry
}

func (cdr *coder) Kind() kube.Kind {
	return kube.KindOfSleep
}

func (cdr *coder) SetConfig(config configio.Config) error {
	if config, ok := config.(*Config); !ok {
		return kube.TypeAssertionError
	} else {
		cdr.config = config
	}

	return nil
}

func (cdr *coder) GetConfig() configio.Config {
	return cdr.config
}

func (cdr *coder) Context() context.Context {
	return cdr.ctx
}

func (cdr *coder) Error() <-chan error {
	return cdr.err
}

func (cdr *coder) Create(ctx context.Context) context.Context {
	log := cdr.log.WithField("func", "Create")
	out := context.Background()
	out, done := context.WithCancel(out)

	go func(input context.Context, f context.CancelFunc) {
		select {
		case <-input.Done():
			cmd := exec.Command(cdr.config.Cmd, cdr.config.Args...)
			cmd.Stdin = cdr.config.stdin
			cmd.Stdout = cdr.config.stdout
			cmd.Stderr = cdr.config.stderr
			log.Info("started ", cdr.config.Cmd)
			if err := cmd.Start(); err != nil {
				cdr.err <- err
			} else {
				if err := cmd.Wait(); err != nil {
					cdr.err <- err
				} else {
					log.Info("finished ", cdr.config.Cmd)
					f()
				}
			}
		case <-cdr.ctx.Done():
			log.Info("self context done")
			return
		}
	}(ctx, done)

	return out
}

func (cdr *coder) Delete(ctx context.Context) context.Context {
	out := context.Background()
	out, done := context.WithCancel(out)
	defer done()
	return out
}
