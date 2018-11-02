package config

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube"
	parent "github.com/sdeoras/kube/kube/deployment"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestConsulDeployment(t *testing.T) {
	log := logrus.WithField("func", "TestServer_fio_GCP").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "consul-deployment"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	appName := "consul"

	// initialize params
	selectorRequirement := new(meta_v1.LabelSelectorRequirement)
	selectorRequirement.Key = "app"
	selectorRequirement.Operator = meta_v1.LabelSelectorOpIn
	selectorRequirement.Values = []string{appName}

	labelSelector := new(meta_v1.LabelSelector)
	labelSelector.MatchExpressions = []meta_v1.LabelSelectorRequirement{*selectorRequirement}

	affinityTerm := new(v1.PodAffinityTerm)
	affinityTerm.LabelSelector = labelSelector
	affinityTerm.TopologyKey = "kubernetes.io/hostname"

	affinity := new(v1.Affinity)
	affinity.PodAntiAffinity = new(v1.PodAntiAffinity)
	affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []v1.PodAffinityTerm{*affinityTerm}

	consulClient := new(v1.Container)
	consulClient.Name = "consul-client"
	consulClient.Image = "sdeoras/consul"
	consulClient.Lifecycle = new(v1.Lifecycle)
	consulClient.Lifecycle.PreStop = new(v1.Handler)
	consulClient.Lifecycle.PreStop.Exec = new(v1.ExecAction)
	consulClient.Lifecycle.PreStop.Exec.Command = []string{"consul", "leave"}
	consulClient.Ports = []v1.ContainerPort{
		{ContainerPort: 8500},
		{ContainerPort: 8400},
		{ContainerPort: 53},
		{ContainerPort: 8443},
		{ContainerPort: 8080},
		{ContainerPort: 8301},
		{ContainerPort: 8302},
		{ContainerPort: 8600},
		{ContainerPort: 8300},
	}
	consulClient.ImagePullPolicy = v1.PullIfNotPresent
	consulClient.Command = []string{"consul", "agent",
		"-config-file", "/etc/consul.d/client/config.json",
		/*
			// for colo
			"-datacenter", "colo",
			"-encrypt", "A7RQdchfff2gR4dtiQcEWg==",
			"-join", "70.0.42.203",
			"-join", "70.0.42.204",
			"-join", "70.0.42.205",
		*/
		// for gcp
		"-datacenter", "gcp",
		"-encrypt", "Bt8oLlKBztQKa8XtiJSqCQ==",
		"-join", "10.138.0.2",
		"-join", "10.138.0.4",
		"-join", "10.138.0.6",
	}

	consulRequest := new(v1.Container)
	consulRequest.Name = "consul-request"
	consulRequest.Image = "sdeoras/consul"
	consulRequest.ImagePullPolicy = v1.PullIfNotPresent
	consulRequest.Command = []string{"curl",
		"http://127.0.0.1:8500/v1/kv/myKey",
	}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = appName

	myDeployment := config.Deployment
	myDeployment.Labels = make(map[string]string)
	myDeployment.Labels["run"] = appName
	myDeployment.Spec.Selector = new(meta_v1.LabelSelector)
	myDeployment.Spec.Selector.MatchLabels = make(map[string]string)
	myDeployment.Spec.Selector.MatchLabels["run"] = appName
	myDeployment.Spec.Template.Spec.Containers = []v1.Container{*consulClient}
	myDeployment.ObjectMeta.Name = appName
	myDeployment.ObjectMeta.Labels = myPodLabels
	myDeployment.Spec.Template.Labels = make(map[string]string)
	myDeployment.Spec.Template.Labels["run"] = appName
	myDeployment.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyAlways
	myDeployment.Spec.Template.Spec.Affinity = affinity

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if b, err := kube.YAMLMarshal(config.Deployment); err != nil {
		t.Fatal(err)
	} else {
		if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
