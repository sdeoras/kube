package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"

	parent "github.com/sdeoras/kube/kube/pvc"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLoadDefaults(t *testing.T) {
	// config init
	key := "pvc_gcp"
	config := new(parent.Config).Init(key)

	// initialize params
	config.PersistentVolumeClaim.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadOnlyMany}
	config.PersistentVolumeClaim.ObjectMeta.Name = "gcp-pvc"
	config.PersistentVolumeClaim.Spec.Resources.Requests = make(map[v1.ResourceName]resource.Quantity)
	config.PersistentVolumeClaim.Spec.Resources.Requests[v1.ResourceStorage] = resource.MustParse("256Gi")
	config.PersistentVolumeClaim.Spec.VolumeName = "gcp-pv"

	b, err := kube.YAMLMarshal(config.PersistentVolumeClaim)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
