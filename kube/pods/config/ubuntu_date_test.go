package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/pods"
	v1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func Test_UbuntuDate(t *testing.T) {
	// config init
	key := "ubuntu_date"
	config := new(parent.Config).Init(key)

	// initialize params
	myVolSecret := new(v1.Volume)
	myVolSecret.Name = "secret-volume"
	myVolSecret.Secret = new(v1.SecretVolumeSource)
	myVolSecret.Secret.SecretName = "gcs-auth"

	myVolMtSecret := new(v1.VolumeMount)
	myVolMtSecret.Name = myVolSecret.Name
	myVolMtSecret.MountPath = "/mnt/secret"

	myContainer := new(v1.Container)
	myContainer.Name = "sleeping"
	myContainer.Image = "ubuntu"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"sleep", "10000"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtSecret}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = "sleeping"

	myPod := config.Pod
	myPod.Spec.Containers = []v1.Container{*myContainer}
	myPod.Spec.Volumes = []v1.Volume{*myVolSecret}
	myPod.ObjectMeta.Name = "sleeping"
	myPod.ObjectMeta.Labels = myPodLabels
	myPod.Spec.RestartPolicy = v1.RestartPolicyNever

	b, err := kube.YAMLMarshal(config.Pod)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
