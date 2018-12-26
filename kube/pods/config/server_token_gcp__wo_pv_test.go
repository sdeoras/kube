package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/pods"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestServer_token_GCP_WO_PV(t *testing.T) {
	log := logrus.WithField("func", "TestServer_token_GCP").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "pods_server_token_gcp_wo_pv"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// initialize params
	myContainer := new(v1.Container)
	myContainer.Name = "token-server"
	myContainer.Image = "sdeoras/token"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"/token/bin/server", "--dir", "/tmp"}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = "token-server"

	myPod := new(v1.Pod)
	myPod.Spec.Containers = []v1.Container{*myContainer}
	myPod.ObjectMeta.Name = "token-server"
	myPod.ObjectMeta.Labels = myPodLabels
	myPod.Spec.RestartPolicy = v1.RestartPolicyAlways

	config.Pod = myPod

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
