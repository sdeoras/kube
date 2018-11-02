package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/pods"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func Test_UbuntuDate(t *testing.T) {
	// config init
	key := "ubuntu_date"
	config := new(parent.Config).Init(key)

	// initialize params
	myVolPWX := new(v1.Volume)
	myVolPWX.Name = "px-volume"
	myVolPWX.PersistentVolumeClaim = new(v1.PersistentVolumeClaimVolumeSource)
	myVolPWX.PersistentVolumeClaim.ClaimName = "px-pvc-1"

	myVolMtPWX := new(v1.VolumeMount)
	myVolMtPWX.Name = myVolPWX.Name
	myVolMtPWX.MountPath = "/mnt/pwx"

	myContainer := new(v1.Container)
	myContainer.Name = "ubuntu-date"
	myContainer.Image = "ubuntu"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"ls", "-la", "/tmp"}
	//myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtPWX}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = "ubuntu-date"

	myPod := config.Pod
	myPod.Spec.Containers = []v1.Container{*myContainer}
	//myPod.Spec.Volumes = []v1.Volume{*myVolPWX}
	myPod.ObjectMeta.Name = "ubuntu-date"
	myPod.ObjectMeta.Labels = myPodLabels
	myPod.Spec.RestartPolicy = v1.RestartPolicyAlways

	b, err := kube.YAMLMarshal(config.Pod)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
