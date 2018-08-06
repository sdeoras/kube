package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/cm"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestExampleCM(t *testing.T) {
	log := logrus.WithField("func", "TestBusyBoxDS").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "sample-cm"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	config.ConfigMap.Name = "example-cm"
	config.ConfigMap.Data = make(map[string]string)
	config.ConfigMap.Annotations = make(map[string]string)
	config.ConfigMap.BinaryData = make(map[string][]byte)

	config.ConfigMap.Data["a"] = "A"
	config.ConfigMap.Data["b"] = "B"
	config.ConfigMap.Annotations["aa"] = "AA"
	config.ConfigMap.Annotations["ab"] = "AB"
	config.ConfigMap.BinaryData["ba"] = []byte{2, 4, 6, 8}

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
