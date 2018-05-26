package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/ds"
	"github.com/sirupsen/logrus"
	apps_v1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestBusyBoxDS(t *testing.T) {
	log := logrus.WithField("func", "TestBusyBoxDS").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "busybox-ds"
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
	myVolume := new(v1.Volume)
	myVolume.Name = "gcp-pv"
	myVolume.PersistentVolumeClaim = new(v1.PersistentVolumeClaimVolumeSource)
	myVolume.PersistentVolumeClaim.ReadOnly = true
	myVolume.PersistentVolumeClaim.ClaimName = "gcp-pvc"

	myVolumeMount := new(v1.VolumeMount)
	myVolumeMount.Name = myVolume.Name
	myVolumeMount.ReadOnly = true
	myVolumeMount.MountPath = "/mnt/gcp"

	myContainer := new(v1.Container)
	myContainer.Name = "busybox"
	myContainer.Image = "busybox"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"sleep", "10000"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolumeMount}

	podTemplateSpec := new(v1.PodTemplateSpec)
	podTemplateSpec.ObjectMeta.Labels = make(map[string]string)
	podTemplateSpec.ObjectMeta.Labels["app"] = "busybox"
	podTemplateSpec.Spec.Containers = []v1.Container{*myContainer}
	podTemplateSpec.Spec.Volumes = []v1.Volume{*myVolume}
	podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicyAlways

	labelRequirement := new(meta_v1.LabelSelectorRequirement)
	labelRequirement.Key = "app"
	labelRequirement.Values = []string{"busybox"}
	labelRequirement.Operator = meta_v1.LabelSelectorOpIn

	labelSelector := new(meta_v1.LabelSelector)
	labelSelector.MatchExpressions = []meta_v1.LabelSelectorRequirement{*labelRequirement}

	myDs := new(apps_v1beta2.DaemonSet)
	myDs.Name = "busybox"
	myDs.Spec = apps_v1beta2.DaemonSetSpec{}
	myDs.Spec.Template = *podTemplateSpec
	myDs.Spec.Selector = labelSelector

	// assign to config
	config.DaemonSet = myDs

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
