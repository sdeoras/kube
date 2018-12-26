package pvc

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cdr *coder) verifyCreate() error {
	log := cdr.log.WithField("func", "Create/verify")
	opts := new(meta_v1.GetOptions)
	for {
		obj, err := cdr.clientset.CoreV1().PersistentVolumeClaims(cdr.namespace).Get(cdr.config.PersistentVolumeClaim.Name, *opts)
		if err != nil {
			log.Error(err)
			return err
		}

		if obj.Status.Phase == v1.ClaimBound {
			log.Info(obj.Status.Phase)
			return nil
		}

		log.Warn(obj.Status.Phase)

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
		_, err := cdr.clientset.CoreV1().PersistentVolumeClaims(cdr.namespace).Get(cdr.config.PersistentVolumeClaim.Name, *opts)
		if err != nil {
			log.Info(err)
			return nil
		}

		log.Warn("deleting")

		select {
		case <-cdr.ctx.Done():
			return fmt.Errorf("coder context done")
		case <-time.After(time.Second * 5):
		}
	}

	return nil
}
