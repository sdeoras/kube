package config

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/jobs"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestCopy_GCP_PX_TMP(t *testing.T) {
	// config init
	key := "jobs_cp_gcp_pwx"
	config := new(parent.Config).Init(key)

	// params to come from outside
	parallel := 1
	jobId := key
	batchSize := 100
	numBatches := 100

	// initialize params
	selectorRequirement := new(meta_v1.LabelSelectorRequirement)
	selectorRequirement.Key = "app"
	selectorRequirement.Operator = meta_v1.LabelSelectorOpIn
	selectorRequirement.Values = []string{"cp-gcp-pwx"}

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
	myVolPWX.Name = "px-pvc-1"
	myVolPWX.PersistentVolumeClaim = new(v1.PersistentVolumeClaimVolumeSource)
	myVolPWX.PersistentVolumeClaim.ClaimName = "px-pvc-1"

	myVolMtPWX := new(v1.VolumeMount)
	myVolMtPWX.Name = myVolPWX.Name
	myVolMtPWX.MountPath = "/mnt/pwx/"

	myContainer := new(v1.Container)
	myContainer.Name = "cp-gcp-pwx"
	myContainer.Image = "sdeoras/token"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"/token/bin/cp",
		"--host", "token-server:7001",
		"--job-id", jobId,
		"--batch-size", strconv.FormatInt(int64(batchSize), 10),
		"--num-batches", strconv.FormatInt(int64(numBatches), 10),
		"--source-dir", "/mnt/gcp/imagenet/input_900",
		"--destination-dir", "/mnt/pwx/images",
		"--out-dir", "/mnt/pwx/out"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtGCP, *myVolMtPWX}

	podTemplateSpec := new(v1.PodTemplateSpec)
	podTemplateSpec.ObjectMeta.Labels = make(map[string]string)
	podTemplateSpec.ObjectMeta.Labels["app"] = "cp-gcp-pwx"
	podTemplateSpec.Spec.Containers = []v1.Container{*myContainer}
	podTemplateSpec.Spec.Volumes = []v1.Volume{*myVolGCP, *myVolPWX}
	podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicyNever
	podTemplateSpec.Spec.Affinity = affinity

	myJob := config.Job
	myJob.Name = "cp-gcp-pwx"
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
