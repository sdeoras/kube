package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"
	parent "github.com/sdeoras/kube/kube/net-policy"
	v1 "k8s.io/api/networking/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNetwork_Policy(t *testing.T) {
	// config init
	key := "net_access_nginx"
	config := new(parent.Config).Init(key, "ns-nginx")

	net := config.NetworkPolicy
	net.Name = "access-nginx"
	net.Namespace = "nginx"
	ls := v12.LabelSelector{
		MatchLabels: map[string]string{
			"run": "nginx",
		},
	}
	net.Spec = v1.NetworkPolicySpec{
		PodSelector: ls,
		Ingress: []v1.NetworkPolicyIngressRule{
			{
				From: []v1.NetworkPolicyPeer{
					{
						PodSelector: &v12.LabelSelector{
							MatchLabels: map[string]string{
								"access": "true",
							},
						},
					},
				},
			},
		},
	}

	b, err := kube.YAMLMarshal(config.NetworkPolicy)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
