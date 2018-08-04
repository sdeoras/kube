package config

import (
	"context"
	"encoding/json"
	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/deployment"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	"path/filepath"
	"testing"
)

func TestServerSpecDeployment(t *testing.T) {
	log := logrus.WithField("func", "TestServer_fio_GCP").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "ServerSpecDeployment"
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
	myVolGCP := new(v1.Volume)
	myVolGCP.Name = "gcp-volume"
	myVolGCP.PersistentVolumeClaim = new(v1.PersistentVolumeClaimVolumeSource)
	myVolGCP.PersistentVolumeClaim.ReadOnly = true
	myVolGCP.PersistentVolumeClaim.ClaimName = "gcp-pvc"

	myVolMtGCP := new(v1.VolumeMount)
	myVolMtGCP.Name = myVolGCP.Name
	myVolMtGCP.ReadOnly = true
	myVolMtGCP.MountPath = "/mnt/gcp"

	myContainer := new(v1.Container)
	myContainer.Name = "transferfile"
	myContainer.Image = "lynned/server:v15"
	myContainer.Ports = []v1.ContainerPort{{ContainerPort: 17361}}
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"./server", "-p", "17361", "-P", "portworx", "-d", "transferfile-deployment"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtGCP}

	myDeploymentLabels := make(map[string]string)
	myDeploymentLabels["app"] = "transferfile"

	myPod := config.Deployment
	myPod.Spec.Template.Spec.Containers = []v1.Container{*myContainer}
	myPod.Spec.Template.Spec.Volumes = []v1.Volume{*myVolGCP}
	myPod.ObjectMeta.Name = "transferfile-deployment"
	myPod.ObjectMeta.Labels = myDeploymentLabels
	myPod.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyAlways

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if b, err := json.MarshalIndent(config.Deployment, "", "  "); err != nil {
		t.Fatal(err)
	} else {
		if err := ioutil.WriteFile(key+".json", b, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
