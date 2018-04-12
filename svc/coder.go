package svc

import (
	"context"

	"errors"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// coder implements kube.Coder interface for a pvc
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
	return kube.KindOfSvc
}

func (cdr *coder) SetConfig(config configio.Config) error {
	if config, ok := config.(*Config); !ok {
		return errors.New(kube.TypeAssertionError)
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
	cdr.clientset = clientset
	cdr.namespace = namespace
	cdr.log = logrus.WithField("package", PackageName).WithField("namespace", cdr.namespace)
}

func (cdr *coder) Create(ctx context.Context) context.Context {
	log := cdr.log.WithField("func", "Create")
	out := context.Background()
	out, done := context.WithCancel(out)

	go func(input context.Context, done context.CancelFunc) {
		select {
		case <-input.Done():
			if _, err := cdr.clientset.CoreV1().Services(cdr.namespace).Create(cdr.config.Svc); err != nil {
				log.Error(err)
				cdr.err <- err
				return
			} else {
				// verifyCreate if running
				/*if err := cdr.verifyCreate(); err != nil {
					cdr.err <- err
				}*/
				log.Info("done")
				done()
				return
			}
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
			options := new(meta_v1.DeleteOptions)
			if err := cdr.clientset.CoreV1().Services(cdr.namespace).Delete(cdr.config.Svc.Name, options); err != nil {
				log.Error(err)
				cdr.err <- err
				return
			} else {
				// verifyCreate if deleted
				/*if err := cdr.verifyDelete(); err != nil {
					cdr.err <- err
				}*/
				log.Info("done")
				f()
				return
			}
		case <-cdr.ctx.Done():
			log.Info("self context done")
			return
		}
	}(ctx, done)

	return out
}
