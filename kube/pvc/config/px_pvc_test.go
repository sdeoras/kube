package config

import (
	"io/ioutil"
	"testing"

	"github.com/sdeoras/kube"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	parent "github.com/sdeoras/kube/kube/pvc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestPxPVC(t *testing.T) {
	// config init
	key := "px_pvc"
	config := new(parent.Config).Init(key)

	// initialize params
	config.PersistentVolumeClaim.ObjectMeta.Name = "px-pvc-1"
	config.PersistentVolumeClaim.Annotations = map[string]string{
		"volume.beta.kubernetes.io/storage-class": "px-high-rf2",
	}
	config.PersistentVolumeClaim.Spec = v1.PersistentVolumeClaimSpec{
		AccessModes: []v1.PersistentVolumeAccessMode{
			v1.ReadWriteMany,
		},
		Resources: v1.ResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceStorage: resource.MustParse("2Gi"),
			},
		},
	}

	b, err := kube.YAMLMarshal(config.PersistentVolumeClaim)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(key+".yaml", b, 0644); err != nil {
		t.Fatal(err)
	}
}
