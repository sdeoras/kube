package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/jobs"
	"github.com/sirupsen/logrus"
	batch_v1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestCmd_GCP_WO_PV(t *testing.T) {
	log := logrus.WithField("func", "TestCmd_GCP").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "jobs_cmd_date_3"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// params to come from outside
	parallel := 3

	// initialize params
	selectorRequirement := new(meta_v1.LabelSelectorRequirement)
	selectorRequirement.Key = "app"
	selectorRequirement.Operator = meta_v1.LabelSelectorOpIn
	selectorRequirement.Values = []string{"cmd-date"}

	labelSelector := new(meta_v1.LabelSelector)
	labelSelector.MatchExpressions = []meta_v1.LabelSelectorRequirement{*selectorRequirement}

	affinityTerm := new(v1.PodAffinityTerm)
	affinityTerm.LabelSelector = labelSelector
	affinityTerm.TopologyKey = "kubernetes.io/hostname"

	affinity := new(v1.Affinity)
	affinity.PodAntiAffinity = new(v1.PodAntiAffinity)
	affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []v1.PodAffinityTerm{*affinityTerm}

	myContainer := new(v1.Container)
	myContainer.Name = "ubuntu"
	myContainer.Image = "ubuntu"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"date"}

	podTemplateSpec := new(v1.PodTemplateSpec)
	podTemplateSpec.ObjectMeta.Labels = make(map[string]string)
	podTemplateSpec.ObjectMeta.Labels["app"] = "cmd-date"
	podTemplateSpec.Spec.Containers = []v1.Container{*myContainer}
	podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicyNever
	podTemplateSpec.Spec.Affinity = affinity

	myJob := new(batch_v1.Job)
	myJob.Name = "cmd-date"
	parallelism := new(int32)
	*parallelism = int32(parallel)
	myJob.Spec.Parallelism = parallelism
	myJob.Spec.Template = *podTemplateSpec

	// assign to config
	config.Job = myJob

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
