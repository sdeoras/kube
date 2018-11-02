package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/sc"
)

func TestPxStorageClass(t *testing.T) {
	key := "storage-class-px-high-rf2"
	config := new(parent.Config).Init(key)

	storageClass := config.StorageClass
	storageClass.Kind = parent.Kind
	storageClass.APIVersion = parent.APIVersion

	storageClass.Name = "px-high-rf2"
	storageClass.Provisioner = parent.PortworxVolume
	storageClass.Parameters = map[string]string{
		"fs":            "ext4",
		"block_size":    "4k",
		"shared":        "true",
		"repl":          "2",
		"snap_interval": "0",
		"priority_io":   "high",
	}

	b, err := kube.YAMLMarshal(config.StorageClass)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}

}
