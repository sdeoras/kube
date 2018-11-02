package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"io/ioutil"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube"
	parent "github.com/sdeoras/kube/kube/pv"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func Test_Local_PV(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "pv_gcp_local"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// initialize params
	config.PersistentVolume.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteMany}
	config.PersistentVolume.Spec.Local = new(v1.LocalVolumeSource)
	config.PersistentVolume.Spec.Local.Path = "/tmp"
	config.PersistentVolume.ObjectMeta.Name = "gcp-local-pv"
	config.PersistentVolume.Spec.Capacity = make(map[v1.ResourceName]resource.Quantity)
	config.PersistentVolume.Spec.Capacity[v1.ResourceStorage] = resource.MustParse("2G")
	config.PersistentVolume.Spec.StorageClassName = "standard"

	config.PersistentVolume.Annotations = make(map[string]string)

	config.PersistentVolume.Kind = parent.Kind
	config.PersistentVolume.APIVersion = parent.APIVersion

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	b, err := kube.YAMLMarshal(config.PersistentVolume)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
