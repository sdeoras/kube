package config

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/pv"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLoadDefaults(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "pv_gcp"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)

	// initialize params
	config.PersistentVolume.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany}
	config.PersistentVolume.Spec.GCEPersistentDisk = new(v1.GCEPersistentDiskVolumeSource)
	config.PersistentVolume.Spec.GCEPersistentDisk.PDName = "tf-data-disk-1"
	config.PersistentVolume.Spec.GCEPersistentDisk.ReadOnly = true
	config.PersistentVolume.ObjectMeta.Name = "gcp-pv"
	config.PersistentVolume.Spec.Capacity = make(map[v1.ResourceName]resource.Quantity)
	config.PersistentVolume.Spec.Capacity[v1.ResourceStorage] = resource.MustParse("256Gi")
	config.PersistentVolume.Spec.StorageClassName = "standard"

	b, err := kube.YAMLMarshal(config.PersistentVolume)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
