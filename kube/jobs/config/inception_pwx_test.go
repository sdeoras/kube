package config

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/sdeoras/kube"

	"github.com/sdeoras/kube/kube/jobs"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInception_PWX(t *testing.T) {
	key := "jobs_inception_pwx"
	config := new(jobs.Config).Init(key)

	// params to come from outside
	parallel := 1
	jobId := key
	batchSize := 100
	numBatches := 100

	// initialize params
	selectorRequirement := new(meta_v1.LabelSelectorRequirement)
	selectorRequirement.Key = "app"
	selectorRequirement.Operator = meta_v1.LabelSelectorOpIn
	selectorRequirement.Values = []string{"pwx-inception"}

	labelSelector := new(meta_v1.LabelSelector)
	labelSelector.MatchExpressions = []meta_v1.LabelSelectorRequirement{*selectorRequirement}

	affinityTerm := new(v1.PodAffinityTerm)
	affinityTerm.LabelSelector = labelSelector
	affinityTerm.TopologyKey = "kubernetes.io/hostname"

	affinity := new(v1.Affinity)
	affinity.PodAntiAffinity = new(v1.PodAntiAffinity)
	affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []v1.PodAffinityTerm{*affinityTerm}

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
	myContainer.Name = "pwx-inception"
	myContainer.Image = "sdeoras/token"
	myContainer.ImagePullPolicy = v1.PullAlways
	myContainer.Command = []string{"/tensorflow/inception",
		"--host", "token-server:7001",
		"--job-id", jobId,
		"--batch-size", strconv.FormatInt(int64(batchSize), 10),
		"--num-batches", strconv.FormatInt(int64(numBatches), 10),
		"--input-dir", "/mnt/pwx/images",
		"--out-dir", "/mnt/pwx/out"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtPWX}

	podTemplateSpec := new(v1.PodTemplateSpec)
	podTemplateSpec.ObjectMeta.Labels = make(map[string]string)
	podTemplateSpec.ObjectMeta.Labels["app"] = "pwx-inception"
	podTemplateSpec.Spec.Containers = []v1.Container{*myContainer}
	podTemplateSpec.Spec.Volumes = []v1.Volume{*myVolPWX}
	podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicyNever
	podTemplateSpec.Spec.Affinity = affinity

	myJob := config.Job
	myJob.Name = "pwx-inception"
	parallelism := new(int32)
	*parallelism = int32(parallel)
	myJob.Spec.Parallelism = parallelism
	myJob.Spec.Template = *podTemplateSpec

	b, err := kube.YAMLMarshal(config.Job)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
