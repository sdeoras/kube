package defaults

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/cm"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
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

	myCm := new(v1.ConfigMap)
	myCm.Name = "example-cm"
	myCm.Data = make(map[string]string)
	myCm.Annotations = make(map[string]string)
	myCm.BinaryData = make(map[string][]byte)

	myCm.Data["a"] = "A"
	myCm.Data["b"] = "B"
	myCm.Annotations["aa"] = "AA"
	myCm.Annotations["ab"] = "AB"
	myCm.BinaryData["ba"] = []byte{2, 4, 6, 8}

	// assign to config
	config.ConfigMap = myCm

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}
}
