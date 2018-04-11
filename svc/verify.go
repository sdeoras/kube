package svc

import (
	"fmt"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cdr *coder) verifyCreate() error {
	log := cdr.log.WithField("func", "verifyCreate")
	opts := new(meta_v1.GetOptions)
	for {
		obj, err := cdr.clientset.CoreV1().Services(cdr.namespace).Get(cdr.config.Svc.Name, *opts)
		if err != nil {
			log.Error(err)
			return err
		}

		for _, status := range obj.Status.LoadBalancer.Ingress {
			if len(status.IP) > 0 {
				return nil
			}
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
	log := cdr.log.WithField("func", "verifyDelete")
	opts := new(meta_v1.GetOptions)
	for {
		_, err := cdr.clientset.CoreV1().Services(cdr.namespace).Get(cdr.config.Svc.Name, *opts)
		if err != nil {
			log.Info(err)
			return nil
		}

		log.Warn("deleting ", cdr.config.Svc.Name)

		select {
		case <-cdr.ctx.Done():
			return fmt.Errorf("coder context done")
		case <-time.After(time.Second * 5):
		}
	}

	return nil
}
