package svc

import (
	"context"

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
	err       error
	log       *logrus.Entry
}

func (cdr *coder) Kind() kube.Kind {
	return kube.KindOfSvc
}

func (cdr *coder) Context() context.Context {
	return cdr.ctx
}

func (cdr *coder) Error() string {
	return cdr.err.Error()
}

func (cdr *coder) Init(clientset *kubernetes.Clientset, configReader configio.ConfigReader) error {

	config := new(Config).Init(cdr.key)
	if err := configReader.Unmarshal(config); err != nil {
		return err
	} else {
		cdr.config = config
	}

	cdr.clientset = clientset

	cdr.log = logrus.WithField("package", PackageName)

	return nil
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
				cdr.err = err
				cdr.cancel()
				log.Info("self context cancelled")
				return
			} else {
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
				cdr.err = err
				cdr.cancel()
				log.Info("self context cancelled")
				return
			} else {
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
