package rbac

import (
	"fmt"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cdr *coder) verifyCreate() error {
	log := cdr.log.WithField("func", "Create/verify")
	opts := new(meta_v1.GetOptions)
	for {
		obj, err := cdr.clientset.AppsV1beta2().DaemonSets(cdr.namespace).Get(cdr.config.ClusterRoleBinding.Name, *opts)
		if err != nil {
			log.Error(err)
			return err
		}

		if obj.Status.DesiredNumberScheduled != 0 && obj.Status.DesiredNumberScheduled == obj.Status.NumberReady {
			log.Info("desired:", obj.Status.DesiredNumberScheduled, ", ready:", obj.Status.NumberReady)
			return nil
		} else {
			log.Warn("desired:", obj.Status.DesiredNumberScheduled, ", ready:", obj.Status.NumberReady)
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
		_, err := cdr.clientset.RbacV1beta1().ClusterRoleBindings().Get(cdr.config.ClusterRoleBinding.Name, *opts)
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
