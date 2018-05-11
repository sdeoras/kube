package cm

import (
	"fmt"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cdr *coder) verifyCreate() error {
	log := cdr.log.WithField("func", "Create/verify")
	opts := new(meta_v1.GetOptions)
	for {
		_, err := cdr.clientset.CoreV1().ConfigMaps(cdr.namespace).Get(cdr.config.ConfigMap.Name, *opts)
		if err != nil {
			log.Error(err)
			return err
		}

		select {
		case <-cdr.ctx.Done():
			return fmt.Errorf("coder context done")
		default:
			return nil
		}
	}

	return nil
}

func (cdr *coder) verifyDelete() error {
	log := cdr.log.WithField("func", "Delete/verify")
	opts := new(meta_v1.GetOptions)
	for {
		_, err := cdr.clientset.CoreV1().ConfigMaps(cdr.namespace).Get(cdr.config.ConfigMap.Name, *opts)
		if err != nil {
			log.Info(err)
			return nil
		}

		select {
		case <-cdr.ctx.Done():
			return fmt.Errorf("coder context done")
		case <-time.After(time.Second * 5):
		}
	}

	return nil
}
