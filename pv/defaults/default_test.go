package defaults

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube/pv"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLoadDefaults(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").WithField("package", "kube/pv/defaults")

	// config init
	key := uuid.New().String()
	config := new(pv.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("HOME"), pv.DefaultConfigDir, pv.DefaultConfigFile)
	configManager, err := configfile.NewManager(context.Background(), configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := configManager.Unmarshal(config); err != nil {
		t.Fatal(err)
	} else {
		log.Info("using pv key: ", key)
	}

	// initialize pv params
	config.PersistentVolume.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany}
	config.PersistentVolume.Spec.GCEPersistentDisk = new(v1.GCEPersistentDiskVolumeSource)
	config.PersistentVolume.Spec.GCEPersistentDisk.PDName = "tf-data-disk-1"
	config.PersistentVolume.Spec.GCEPersistentDisk.ReadOnly = true
	config.PersistentVolume.ObjectMeta.Name = "my-pv"
	config.PersistentVolume.Spec.Capacity = make(map[v1.ResourceName]resource.Quantity)
	config.PersistentVolume.Spec.Capacity[v1.ResourceStorage] = resource.MustParse("256Gi")
	config.PersistentVolume.Spec.StorageClassName = "standard"

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
