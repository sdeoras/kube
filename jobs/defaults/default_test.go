package defaults

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"strconv"

	"github.com/google/uuid"
	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/jobs"
	"github.com/sirupsen/logrus"
	batch_v1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLoadDefaults(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := uuid.New().String()
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("HOME"), parent.DefaultConfigDir, parent.DefaultConfigFile)
	configManager, err := configfile.NewManager(context.Background(), configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// params to come from outside
	parallel := 1
	job_id := key
	batch_size := 100
	num_batches := 100

	// initialize params
	selectorRequirement := new(meta_v1.LabelSelectorRequirement)
	selectorRequirement.Key = "app"
	selectorRequirement.Operator = meta_v1.LabelSelectorOpIn
	selectorRequirement.Values = []string{"token-client-inception"}

	labelSelector := new(meta_v1.LabelSelector)
	labelSelector.MatchExpressions = []meta_v1.LabelSelectorRequirement{*selectorRequirement}

	affinityTerm := new(v1.PodAffinityTerm)
	affinityTerm.LabelSelector = labelSelector
	affinityTerm.TopologyKey = "kubernetes.io/hostname"

	affinity := new(v1.Affinity)
	affinity.PodAntiAffinity = new(v1.PodAntiAffinity)
	affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = []v1.PodAffinityTerm{*affinityTerm}

	myVolume := new(v1.Volume)
	myVolume.Name = "my-volume"
	myVolume.PersistentVolumeClaim = new(v1.PersistentVolumeClaimVolumeSource)
	myVolume.PersistentVolumeClaim.ReadOnly = true
	myVolume.PersistentVolumeClaim.ClaimName = "my-pvc"

	myVolume2 := new(v1.Volume)
	myVolume2.Name = "tmp"
	myVolume2.HostPath = new(v1.HostPathVolumeSource)
	myVolume2.HostPath.Path = "/tmp"

	myVolumeMount := new(v1.VolumeMount)
	myVolumeMount.Name = myVolume.Name
	myVolumeMount.ReadOnly = true
	myVolumeMount.MountPath = "/tf"

	myVolumeMount2 := new(v1.VolumeMount)
	myVolumeMount2.Name = myVolume2.Name
	myVolumeMount2.ReadOnly = false
	myVolumeMount2.MountPath = "/tmp/tf"

	myContainer := new(v1.Container)
	myContainer.Name = "token-client-inception"
	myContainer.Image = "sdeoras/token-inception:1.0.0"
	myContainer.ImagePullPolicy = v1.PullIfNotPresent
	myContainer.Command = []string{"/tensorflow/client",
		"--host", "token-server",
		"--jobid", job_id,
		"--batchsize", strconv.FormatInt(int64(batch_size), 10),
		"--numbatches", strconv.FormatInt(int64(num_batches), 10),
		"--inputdir", "/tf/images",
		"--outdir", "/tmp/tf/out"}
	myContainer.VolumeMounts = []v1.VolumeMount{*myVolumeMount, *myVolumeMount2}

	podTemplateSpec := new(v1.PodTemplateSpec)
	podTemplateSpec.ObjectMeta.Labels = make(map[string]string)
	podTemplateSpec.ObjectMeta.Labels["app"] = "token-client-inception"
	podTemplateSpec.Spec.Containers = []v1.Container{*myContainer}
	podTemplateSpec.Spec.Volumes = []v1.Volume{*myVolume, *myVolume2}
	podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicyNever
	podTemplateSpec.Spec.Affinity = affinity

	myJob := new(batch_v1.Job)
	myJob.Name = "token-client-inception"
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
