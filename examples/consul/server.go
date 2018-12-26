package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/sdeoras/kube"
	"github.com/sdeoras/kube/kube/statefulset"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getServer(w io.Writer) error {
	config := new(statefulset.Config).Init(consulServerName)

	var replicas int32
	var gracePeriod int64
	var fsGroup int64

	Lifecycle := new(v1.Lifecycle)
	Lifecycle.PreStop = new(v1.Handler)
	Lifecycle.PreStop.Exec = new(v1.ExecAction)
	Lifecycle.PreStop.Exec.Command = []string{"consul", "leave"}

	// initialize params
	selectorRequirement := new(meta_v1.LabelSelectorRequirement)
	selectorRequirement.Key = "app"
	selectorRequirement.Operator = meta_v1.LabelSelectorOpIn
	selectorRequirement.Values = []string{consulServerName, consulLoadApp, consulWatchApp}

	labelSelector := new(meta_v1.LabelSelector)
	labelSelector.MatchExpressions = []meta_v1.LabelSelectorRequirement{*selectorRequirement}

	affinityTerm := new(v1.PodAffinityTerm)
	affinityTerm.LabelSelector = labelSelector
	affinityTerm.TopologyKey = "kubernetes.io/hostname"

	affinity := new(v1.Affinity)
	affinity.PodAntiAffinity = new(v1.PodAntiAffinity)
	affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []v1.PodAffinityTerm{*affinityTerm}

	consulServer := new(v1.Container)
	consulServer.Name = consulServerName
	consulServer.Image = consulImage
	consulServer.Env = getEnv()
	consulServer.Args = []string{
		"agent",
		"-advertise=$(POD_IP)",
		"-bind=0.0.0.0",
		"-bootstrap-expect=" + strconv.FormatInt(int64(serverReplicas), 10),
		"-retry-join=consul-0.consul.$(NAMESPACE).svc.cluster.local",
		"-retry-join=consul-1.consul.$(NAMESPACE).svc.cluster.local",
		"-retry-join=consul-2.consul.$(NAMESPACE).svc.cluster.local",
		"-client=0.0.0.0",
		"-datacenter=" + datacenter,
		"-data-dir=/consul/data",
		"-domain=cluster.local",
		"-encrypt=$(GOSSIP_ENCRYPTION_KEY)",
		"-server",
		"-disable-host-node-id",
	}
	consulServer.Lifecycle = Lifecycle
	consulServer.Ports = []v1.ContainerPort{
		{
			ContainerPort: 8500,
			Name:          "ui-port",
		},
		{
			ContainerPort: 8400,
			Name:          "alt-port",
		},
		{
			ContainerPort: 53,
			Name:          "udp-port",
		},
		{
			ContainerPort: 8443,
			Name:          "https-port",
		},
		{
			ContainerPort: 8080,
			Name:          "http-port",
		},
		{
			ContainerPort: 8301,
			Name:          "serflan",
		},
		{
			ContainerPort: 8302,
			Name:          "serfwan",
		},
		{
			ContainerPort: 8600,
			Name:          "consuldns",
		},
		{
			ContainerPort: 8300,
			Name:          "server",
		},
	}

	var podTSpec v1.PodTemplateSpec
	podTSpec.Labels = make(map[string]string)
	podTSpec.Labels["app"] = consulServerName
	podTSpec.Spec.Affinity = affinity
	gracePeriod = 10
	podTSpec.Spec.TerminationGracePeriodSeconds = &gracePeriod
	podTSpec.Spec.SecurityContext = new(v1.PodSecurityContext)
	fsGroup = 1000
	podTSpec.Spec.SecurityContext.FSGroup = &fsGroup
	podTSpec.Spec.Containers = []v1.Container{*consulServer}

	ss := config.StatefulSet
	ss.Name = consulServerName
	ss.Spec.ServiceName = consulServerName
	replicas = serverReplicas
	ss.Spec.Replicas = &replicas
	ss.Spec.Template = podTSpec

	b, err := kube.YAMLMarshal(config.StatefulSet)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(b))
	fmt.Fprintln(w, "---")

	return nil
}
