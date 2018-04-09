package pv

import (
	"context"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/scratchpad/kube"
	"github.com/sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type manager struct {
	key       string
	config    *Config
	clientset *kubernetes.Clientset
	ctx       context.Context
	cancel    context.CancelFunc
	err       error
	log       *logrus.Entry
}

func (m *manager) Kind() kube.Kind {
	return kube.KindOfPv
}

func (m *manager) Context() context.Context {
	return m.ctx
}

func (m *manager) Error() string {
	return m.err.Error()
}

func (m *manager) Init(clientset *kubernetes.Clientset, configReader configio.ConfigReader) error {

	config := new(Config).Init(m.key)
	if err := configReader.Unmarshal(config); err != nil {
		return err
	} else {
		m.config = config
	}

	m.clientset = clientset

	m.log = logrus.WithField("package", "kube/pv")

	return nil
}

func (m *manager) Create(ctx context.Context) context.Context {
	log := m.log.WithField("func", "Deploy")
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

func (m *manager) Delete(ctx context.Context) context.Context {
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
