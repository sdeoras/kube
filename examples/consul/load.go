package main

import (
	"fmt"
	"io"

	"github.com/sdeoras/kube"
	parent "github.com/sdeoras/kube/kube/pods"
	"k8s.io/api/core/v1"
)

func getLoad(w io.Writer) error {
	config := new(parent.Config).Init(consulLoadApp)

	consulClient := new(v1.Container)
	consulClient.Name = "consul-client"
	consulClient.Image = "sdeoras/consul"
	consulClient.Lifecycle = getLifeCycle()
	consulClient.ImagePullPolicy = v1.PullIfNotPresent
	consulClient.Env = getEnv()
	consulClient.Command = []string{"consul", "agent",
		"-config-file", "/etc/consul.d/client/config.json",
		"-datacenter", datacenter,
		"-encrypt", "$(GOSSIP_ENCRYPTION_KEY)",
		"-join", "consul-0.consul.$(NAMESPACE).svc.cluster.local",
		"-join", "consul-1.consul.$(NAMESPACE).svc.cluster.local",
		"-join", "consul-2.consul.$(NAMESPACE).svc.cluster.local",
	}

	consulRequest := new(v1.Container)
	consulRequest.Name = "consul-request"
	consulRequest.Image = "sdeoras/consul"
	consulRequest.ImagePullPolicy = v1.PullIfNotPresent
	consulRequest.Lifecycle = getLifeCycle()
	consulRequest.Command = []string{
		"/kvdb", "load", "--count", "100", "--key", "bar/baz", "--prefix", kvPrefix, "--leave",
	}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = consulLoadApp

	myPod := config.Pod
	myPod.Spec.Containers = []v1.Container{*consulClient, *consulRequest}
	myPod.ObjectMeta.Name = consulLoadApp
	myPod.ObjectMeta.Labels = myPodLabels
	myPod.Spec.RestartPolicy = v1.RestartPolicyNever
	myPod.Spec.Affinity = getAffinity()

	b, err := kube.YAMLMarshal(config.Pod)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(b))
	fmt.Fprintln(w, "---")

	return nil
}
