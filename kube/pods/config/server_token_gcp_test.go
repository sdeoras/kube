package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/pods"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestServer_token_GCP(t *testing.T) {
	log := logrus.WithField("func", "TestServer_token_GCP").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "pods_server_token_gcp"
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

	myVolPWX := new(v1.Volume)
	myVolPWX.Name = "pwx-vol-1"
	myVolPWX.PortworxVolume = new(v1.PortworxVolumeSource)
	myVolPWX.PortworxVolume.VolumeID = myVolPWX.Name

	myVolMtPWX := new(v1.VolumeMount)
	myVolMtPWX.Name = myVolPWX.Name
	myVolMtPWX.MountPath = "/mnt/pwx/"

	myContainer := new(v1.Container)
	myContainer.Name = "token-server"
	myContainer.Image = "sdeoras/token"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"/token/bin/server", "--dir", "/mnt/gcp/images"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtGCP}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = "token-server"

	myPod := new(v1.Pod)
	myPod.Spec.Containers = []v1.Container{*myContainer}
	myPod.Spec.Volumes = []v1.Volume{*myVolGCP}
	myPod.ObjectMeta.Name = "token-server"
	myPod.ObjectMeta.Labels = myPodLabels
	myPod.Spec.RestartPolicy = v1.RestartPolicyAlways

	config.Pod = myPod

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
