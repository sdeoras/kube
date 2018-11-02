package config

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube"
	parent "github.com/sdeoras/kube/kube/endpoints"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestConsulDeployment(t *testing.T) {
	log := logrus.WithField("func", "TestServer_fio_GCP").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "consul-endpoints"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	endpoints := config.Endpoints
	endpoints.Name = "consul-endpoints"
	endpoints.Subsets = []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "10.138.0.2",
				},
				{
					IP: "10.138.0.4",
				},
				{
					IP: "10.138.0.6",
				},
			},
			Ports: []v1.EndpointPort{
				{
					Port: 8500,
					Name: "http",
				},
				{
					Port: 8300,
					Name: "gossip",
				},
				{
					Port: 8301,
					Name: "gossip2",
				},
			},
		},
	}

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if b, err := kube.YAMLMarshal(config.Endpoints); err != nil {
		t.Fatal(err)
	} else {
		if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
