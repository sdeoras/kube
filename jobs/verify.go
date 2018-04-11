package jobs

import (
	"fmt"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cdr *coder) verifyCreate() error {
	log := cdr.log.WithField("func", "verifyCreate")
	opts := new(meta_v1.GetOptions)
	for {
		obj, err := cdr.clientset.BatchV1().Jobs(cdr.namespace).Get(cdr.config.Job.Name, *opts)
		if err != nil {
			log.Error(err)
			return err
		}

		if obj.Status.Active != 0 {
			log.Warn("active:", obj.Status.Active)
		}

		if obj.Status.Failed != 0 {
			log.Error("failed:", obj.Status.Failed)
			return fmt.Errorf("%d: %s", obj.Status.Failed, "jobs failed")
		}

		if obj.Status.Succeeded == *cdr.config.Job.Spec.Parallelism {
			log.Info("%d: %s", obj.Status.Succeeded, "jobs succeeded")
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
	log := cdr.log.WithField("func", "verifyCreate")
	opts := new(meta_v1.GetOptions)
	for {
		_, err := cdr.clientset.BatchV1().Jobs(cdr.namespace).Get(cdr.config.Job.Name, *opts)
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
