package pv

import (
	"context"

	"os"
	"path/filepath"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// coder implements kube.Coder interface for a pv
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

func GetDefaultConfigFile() string {
	return filepath.Join(os.Getenv("HOME"), DefaultConfigDir, DefaultConfigFile)
}

func (m *coder) Kind() kube.Kind {
	return kube.KindOfPv
}

func (m *coder) Context() context.Context {
	return m.ctx
}

func (m *coder) Error() string {
	return m.err.Error()
}

func (cdr *coder) Clientset(clientset *kubernetes.Clientset, namespace string) {
	cdr.clientset = clientset
	cdr.namespace = namespace
	cdr.log = logrus.WithField("package", PackageName).WithField("namespace", cdr.namespace)
}

func (m *coder) Create(ctx context.Context) context.Context {
	log := m.log.WithField("func", "Create")
	out := context.Background()
	out, done := context.WithCancel(out)

	go func(input context.Context, done context.CancelFunc) {
		select {
		case <-input.Done():
			if _, err := m.clientset.CoreV1().PersistentVolumes().Create(m.config.PersistentVolume); err != nil {
				log.Error(err)
				m.err = err
				m.cancel()
				log.Info("self context cancelled")
				return
			} else {
				log.Info("done")
				done()
				return
			}
		case <-m.ctx.Done():
			log.Info("self context done")
			return
		}
	}(ctx, done)

	return out
}

func (m *coder) Delete(ctx context.Context) context.Context {
	log := m.log.WithField("func", "Delete")
	out := context.Background()
	out, done := context.WithCancel(out)

	go func(parent context.Context, f context.CancelFunc) {
		select {
		case <-parent.Done():
			options := new(meta_v1.DeleteOptions)
			if err := m.clientset.CoreV1().PersistentVolumes().Delete(m.config.PersistentVolume.Name, options); err != nil {
				log.Error(err)
				m.err = err
				m.cancel()
				log.Info("self context cancelled")
				return
			} else {
				log.Info("done")
				f()
				return
			}
		case <-m.ctx.Done():
			log.Info("self context done")
			return
		}
	}(ctx, done)

	return out
}
