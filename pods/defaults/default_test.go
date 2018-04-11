package defaults

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/pods"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLoadDefaults(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := uuid.New().String()
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com/sdeoras",
		parent.PackageName, "defaults", parent.DefaultConfigDir, parent.DefaultConfigFile)
	configManager, err := configfile.NewManager(context.Background(), configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// initialize params
	myVolume := new(v1.Volume)
	myVolume.Name = "my-volume"
	myVolume.PersistentVolumeClaim = new(v1.PersistentVolumeClaimVolumeSource)
	myVolume.PersistentVolumeClaim.ReadOnly = true
	myVolume.PersistentVolumeClaim.ClaimName = "my-pvc"

	myVolumeMount := new(v1.VolumeMount)
	myVolumeMount.Name = myVolume.Name
	myVolumeMount.ReadOnly = true
	myVolumeMount.MountPath = "/tf"

	myContainer := new(v1.Container)
	myContainer.Name = "token-server"
	myContainer.Image = "sdeoras/token-server:1.0.0"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"/server", "--dir", "/tf/images"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolumeMount}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = "token-server"

	myPod := new(v1.Pod)
	myPod.Spec.Containers = []v1.Container{*myContainer}
	myPod.Spec.Volumes = []v1.Volume{*myVolume}
	myPod.ObjectMeta.Name = "token-server"
	myPod.ObjectMeta.Labels = myPodLabels
	myPod.Spec.RestartPolicy = v1.RestartPolicyAlways

	config.Pod = myPod

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
