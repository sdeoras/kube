package ns

import (
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cdr *coder) verifyCreate() error {
	log := cdr.log.WithField("func", "Create/verify")
	opts := new(meta_v1.GetOptions)
	for {
		obj, err := cdr.clientset.CoreV1().Namespaces().Get(cdr.config.Namespace.Name, *opts)
		if err != nil {
			log.Error(err)
			return err
		}

		log.Info(obj.Status.Phase)
		if obj.Status.Phase == v1.NamespaceActive {
			break
		}

		select {
		case <-cdr.ctx.Done():
			return fmt.Errorf("coder context done")
		case <-time.After(time.Second * 5):
		}
	}

	return nil
}

func (cdr *coder) verifyDelete() error {
	log := cdr.log.WithField("func", "Delete/verify")
	opts := new(meta_v1.GetOptions)
	for {
		obj, err := cdr.clientset.CoreV1().Namespaces().Get(cdr.config.Namespace.Name, *opts)
		if err != nil {
			log.Info(err)
			return nil
		}

		log.Info(obj.Status.Phase)

		select {
		case <-cdr.ctx.Done():
			return fmt.Errorf("coder context done")
		case <-time.After(time.Second * 5):
		}
	}

	return nil
}
