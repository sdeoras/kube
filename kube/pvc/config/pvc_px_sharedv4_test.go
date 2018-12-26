package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"encoding/json"
	"io/ioutil"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/pvc"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestPxPvcSharedV4(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "pvc_px_sharedv4"
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
	config.PersistentVolumeClaim.Name = "px-pvc-1"
	config.PersistentVolumeClaim.Annotations = make(map[string]string)
	config.PersistentVolumeClaim.Annotations["volume.beta.kubernetes.io/storage-class"] = "portworx-sharedv4"
	config.PersistentVolumeClaim.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteMany}
	requests := make(map[v1.ResourceName]resource.Quantity)
	requests["storage"] = resource.MustParse("1Gi")
	config.PersistentVolumeClaim.Spec.Resources = v1.ResourceRequirements{
		Requests: requests,
	}

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if b, err := json.MarshalIndent(config.PersistentVolumeClaim, "", "  "); err != nil {
		t.Fatal(err)
	} else {
		if err := ioutil.WriteFile(key+".json", b, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
