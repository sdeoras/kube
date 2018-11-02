package main

import (
	"fmt"
	"io"

	"github.com/sdeoras/kube"
	"github.com/sdeoras/kube/kube/svc"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getService(w io.Writer) error {
	config := new(svc.Config).Init(consulService)

	config.Svc.Name = consulService
	config.Svc.Labels = make(map[string]string)
	config.Svc.Labels["name"] = consulService
	config.Svc.Spec.ClusterIP = "None"
	config.Svc.Spec.Ports = []v1.ServicePort{
		{
			Name:       "http",
			Port:       8500,
			TargetPort: intstr.FromInt(8500),
		},
		{
			Name:       "https",
			Port:       8443,
			TargetPort: intstr.FromInt(8443),
		},
		{
			Name:       "rpc",
			Port:       8400,
			TargetPort: intstr.FromInt(8400),
		},
		{
			Name:       "serflan-tcp",
			Protocol:   "TCP",
			Port:       8301,
			TargetPort: intstr.FromInt(8301),
		},
		{
			Name:       "serflan-udp",
			Protocol:   "UDP",
			Port:       8301,
			TargetPort: intstr.FromInt(8301),
		},
		{
			Name:       "serfwan-tcp",
			Protocol:   "TCP",
			Port:       8302,
			TargetPort: intstr.FromInt(8302),
		},
		{
			Name:       "serfwan-udp",
			Protocol:   "UDP",
			Port:       8302,
			TargetPort: intstr.FromInt(8302),
		},
		{
			Name:       "server",
			Port:       8300,
			TargetPort: intstr.FromInt(8300),
		},
		{
			Name:       "consuldns",
			Port:       8600,
			TargetPort: intstr.FromInt(8600),
		},
	}
	config.Svc.Spec.Selector = make(map[string]string)
	config.Svc.Spec.Selector["app"] = consulService

	b, err := kube.YAMLMarshal(config.Svc)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(b))
	fmt.Fprintln(w, "---")

	return nil
}
