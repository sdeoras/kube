package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/svc"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLoadDefaults(t *testing.T) {
	key := "service_token_server"
	config := new(parent.Config).Init(key)

	// initialize params
	myService := config.Svc
	myService.Name = "token-server"
	myService.ObjectMeta.Name = "token-server"
	myService.Spec.Selector = make(map[string]string)
	myService.Spec.Selector["app"] = "token-server"
	myService.Spec.Ports = []v1.ServicePort{
		{
			Protocol:   v1.ProtocolTCP,
			Port:       7001,
			TargetPort: intstr.FromInt(7001),
		},
	}

	b, err := kube.YAMLMarshal(config.Svc)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
