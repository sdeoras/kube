package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/sc"
	"github.com/sirupsen/logrus"
	"k8s.io/api/storage/v1beta1"
)

func TestPxSharedv4(t *testing.T) {
	log := logrus.WithField("func", "TestBusyBoxDS").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "storage-class-px-sharedv4"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	storageClass := new(v1beta1.StorageClass)
	storageClass.Kind = parent.Kind
	storageClass.APIVersion = parent.APIVersion

	storageClass.Name = "portworx-sharedv4"
	storageClass.Provisioner = parent.PortworxVolume
	storageClass.Parameters = make(map[string]string)
	storageClass.Parameters["repl"] = "1"
	storageClass.Parameters["sharedv4"] = "true"

	// assign to config
	config.StorageClass = storageClass

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if b, err := json.MarshalIndent(config.StorageClass, "", "  "); err != nil {
		t.Fatal(err)
	} else {
		if err := ioutil.WriteFile(key+".json", b, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
