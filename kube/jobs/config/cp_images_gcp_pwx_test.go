package config

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/jobs"
	"github.com/sirupsen/logrus"
	batch_v1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestCopy_images_GCP_PWX(t *testing.T) {
	log := logrus.WithField("func", "TestCopy_images_GCP_PWX").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "jobs_cp_images_gcp_pwx"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// params to come from outside
	parallel := 1
	jobId := key
	batchSize := 100
	numBatches := 100

	// initialize params
	selectorRequirement := new(meta_v1.LabelSelectorRequirement)
	selectorRequirement.Key = "app"
	selectorRequirement.Operator = meta_v1.LabelSelectorOpIn
	selectorRequirement.Values = []string{"cp-images-gcp-pwx"}

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

	myVolTMP := new(v1.Volume)
	myVolTMP.Name = "tmp-volume"
	myVolTMP.HostPath = new(v1.HostPathVolumeSource)
	myVolTMP.HostPath.Path = "/tmp"

	myVolMtTMP := new(v1.VolumeMount)
	myVolMtTMP.Name = myVolTMP.Name
	myVolMtTMP.MountPath = "/mnt/host"

	myContainer := new(v1.Container)
	myContainer.Name = "cp-images-gcp-pwx"
	myContainer.Image = "sdeoras/token"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"/token/bin/cp",
		"--host", "token-server:7001",
		"--job-id", jobId,
		"--batch-size", strconv.FormatInt(int64(batchSize), 10),
		"--num-batches", strconv.FormatInt(int64(numBatches), 10),
		"--source-dir", "/mnt/gcp/images",
		"--destination-dir", "/mnt/pwx/images",
		"--out-dir", "/mnt/pwx/token/cp/out"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolMtGCP, *myVolMtPWX}

	podTemplateSpec := new(v1.PodTemplateSpec)
	podTemplateSpec.ObjectMeta.Labels = make(map[string]string)
	podTemplateSpec.ObjectMeta.Labels["app"] = "cp-images-gcp-pwx"
	podTemplateSpec.Spec.Containers = []v1.Container{*myContainer}
	podTemplateSpec.Spec.Volumes = []v1.Volume{*myVolGCP, *myVolPWX}
	podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicyNever
	podTemplateSpec.Spec.Affinity = affinity

	myJob := new(batch_v1.Job)
	myJob.Name = "cp-images-gcp-pwx"
	parallelism := new(int32)
	*parallelism = int32(parallel)
	myJob.Spec.Parallelism = parallelism
	myJob.Spec.Template = *podTemplateSpec

	// assign to config
	config.Job = myJob

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
