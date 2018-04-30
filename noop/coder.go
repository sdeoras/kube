package noop

import (
	"context"

	"sync"
	"time"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type data struct {
	sync.Mutex
	value int
}

var ind data

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
	return kube.KindOfNoop
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

func (cdr *coder) Clientset(clientset *kubernetes.Clientset, namespace string) {
	// no op here
}

func (cdr *coder) Create(ctx context.Context) context.Context {
	log := cdr.log.WithField("func", "Create")
	out := context.Background()
	out, done := context.WithCancel(out)

	go func(input context.Context, done context.CancelFunc) {
		select {
		case <-input.Done():
			ind.Lock()
			cdr.config.Value = ind.value
			ind.value += 1
			ind.Unlock()
			time.Sleep(time.Second)
			log.Info("started ", cdr.config.Name)
			done()
		case <-cdr.ctx.Done():
			log.Info("self context done")
			return
		}
	}(ctx, done)

	return out
}

func (cdr *coder) Delete(ctx context.Context) context.Context {
	log := cdr.log.WithField("func", "Delete")
	out := context.Background()
	out, done := context.WithCancel(out)

	go func(parent context.Context, f context.CancelFunc) {
		select {
		case <-parent.Done():
			ind.Lock()
			ind.value -= 1
			ind.Unlock()
			time.Sleep(time.Second)
			log.Info("stopped ", cdr.config.Name)
			f()
		case <-cdr.ctx.Done():
			log.Info("self context done")
			return
		}
	}(ctx, done)

	return out
}
