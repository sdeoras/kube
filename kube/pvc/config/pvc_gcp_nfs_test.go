package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"io/ioutil"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube"
	parent "github.com/sdeoras/kube/kube/pvc"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func Test_NFS_PVC(t *testing.T) {
	log := logrus.WithField("func", "TestLoadDefaults").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "pvc_gcp_nfs"
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
	config.PersistentVolumeClaim.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteMany}
	config.PersistentVolumeClaim.ObjectMeta.Name = "gcp-nfs-pvc"
	config.PersistentVolumeClaim.Spec.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
	config.PersistentVolumeClaim.Spec.Resources.Requests[v1.ResourceStorage] = resource.MustParse("1T")
	config.PersistentVolumeClaim.Spec.VolumeName = "gcp-nfs-pv"

	config.PersistentVolumeClaim.Kind = parent.Kind
	config.PersistentVolumeClaim.APIVersion = parent.APIVersion

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	b, err := kube.YAMLMarshal(config.PersistentVolumeClaim)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
