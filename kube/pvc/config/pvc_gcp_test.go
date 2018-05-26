package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/pvc"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLoadDefaults(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "pvc_gcp"
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
	config.PersistentVolumeClaim.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany}
	config.PersistentVolumeClaim.ObjectMeta.Name = "gcp-pvc"
	config.PersistentVolumeClaim.Spec.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
	config.PersistentVolumeClaim.Spec.Resources.Requests[v1.ResourceStorage] = resource.MustParse("256Gi")
	config.PersistentVolumeClaim.Spec.VolumeName = "gcp-pv"

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
