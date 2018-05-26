package nl

import (
	"context"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// coder implements kube.Coder interface
type coder struct {
	key       string
	namespace string
	config    *Config
	clientset *kubernetes.Clientset
	ctx       context.Context
	cancel    context.CancelFunc
	err       chan error
	log       *logrus.Entry
}

func (cdr *coder) Kind() kube.Kind {
	return kube.KindOfDs
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

	go func(input context.Context, done context.CancelFunc) {
		select {
		case <-input.Done():
			var err error
			listOpts := new(meta_v1.ListOptions)
			if cdr.config.NodeList, err = cdr.clientset.CoreV1().Nodes().List(*listOpts); err != nil {
				log.Error(err)
				cdr.err <- err
			} else {
				done()
			}
		case <-cdr.ctx.Done():
			log.Info("self context done")
			return
		}
	}(ctx, done)

	return out
}

func (cdr *coder) Delete(ctx context.Context) context.Context {
	// log := cdr.log.WithField("func", "Delete")
	out := context.Background()
	out, done := context.WithCancel(out)
	defer done()

	return out
}
