package deployment

import (
	"fmt"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cdr *coder) verifyCreate() error {
	log := cdr.log.WithField("func", "Create/verify")
	opts := new(meta_v1.GetOptions)
	for {
		obj, err := cdr.clientset.CoreV1().Pods(cdr.namespace).Get(cdr.config.Deployment.Name, *opts)
		if err != nil {
			log.Error(err)
			return err
		}

		ready := true
		for _, status := range obj.Status.ContainerStatuses {
			if !status.Ready {
				log.Warn("not all containers are ready")
				ready = false
				break
			}
		}

		if ready {
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

func (cdr *coder) verifyDelete() error {
	log := cdr.log.WithField("func", "Delete/verify")
	opts := new(meta_v1.GetOptions)
	for {
		_, err := cdr.clientset.AppsV1().Deployments(cdr.namespace).Get(cdr.config.Deployment.Name, *opts)
		if err != nil {
			log.Info(err)
			return nil
		} else {
			log.Warn("deleting")
		}

		select {
		case <-cdr.ctx.Done():
			return fmt.Errorf("coder context done")
		case <-time.After(time.Second * 5):
		}
	}

	return nil
}
