package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/ds"
	apps_v1beta2 "k8s.io/api/apps/v1beta2"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestBusyBoxDS(t *testing.T) {
	// config init
	key := "busybox-ds"

	config := new(parent.Config).Init(key)

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

	myDs := config.DaemonSet
	myDs.Name = "busybox"
	myDs.Spec = apps_v1beta2.DaemonSetSpec{}
	myDs.Spec.Template = *podTemplateSpec
	myDs.Spec.Selector = labelSelector

	b, err := kube.YAMLMarshal(config.DaemonSet)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
