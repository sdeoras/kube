package config

import (
	"io/ioutil"
	"testing"

	"github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	"github.com/sdeoras/kube"
	parent "github.com/sdeoras/kube/kube/tfjob"
	v1 "k8s.io/api/core/v1"
)

func TestLoadDefaults(t *testing.T) {
	key := "tfjob_worker"
	config := new(parent.Config).Init(key)

	// initialize params
	job := config.Job
	job.Name = "example-job"

	TFReplicaSpecs := make(map[v1alpha2.TFReplicaType]*v1alpha2.TFReplicaSpec)

	worker := new(v1alpha2.TFReplicaSpec)
	TFReplicaSpecs[v1alpha2.TFReplicaTypeWorker] = worker
	job.Spec = v1alpha2.TFJobSpec{
		TFReplicaSpecs: TFReplicaSpecs,
	}

	worker.Template = v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "tensorflow",
					Image: "gcr.io/tf-on-k8s-dogfood/tf_sample:dc944ff",
				},
			},
		},
	}

	b, err := kube.YAMLMarshal(config.Job)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
