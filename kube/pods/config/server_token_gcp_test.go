package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/pods"
	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestServer_token_GCP(t *testing.T) {
	// config init
	key := "pods_server_token_gcp"
	config := new(parent.Config).Init(key)

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
	myContainer.Name = "token-server"
	myContainer.Image = "sdeoras/token"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"/token/bin/server", "--dir", "/mnt/gcp/imagenet/input_900"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtGCP}

	myPodLabels := make(map[string]string)
	myPodLabels["app"] = "token-server"

	myPod := config.Pod
	myPod.Spec.Containers = []v1.Container{*myContainer}
	myPod.Spec.Volumes = []v1.Volume{*myVolGCP}
	myPod.ObjectMeta.Name = "token-server"
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
