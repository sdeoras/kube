package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sdeoras/configio/configfile"
	parent "github.com/sdeoras/kube/kube/rbac"
	"github.com/sirupsen/logrus"
	"k8s.io/api/rbac/v1beta1"
)

func TestClusterRbac(t *testing.T) {
	log := logrus.WithField("func", "TestBusyBoxDS").
		WithField("package", filepath.Join(parent.PackageName, "defaults"))

	// config init
	key := "cluster-rbac"
	log.Info(parent.PackageName, " using key: ", key)
	config := new(parent.Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(context.Background(),
		configfile.OptFilePath, configFilePath)
	if err != nil {
		t.Fatal(err)
	}

	rbac := new(v1beta1.ClusterRoleBinding)
	rbac.Kind = parent.Kind
	rbac.APIVersion = parent.APIVersion

	rbac.Name = "cluster-rbac"
	rbac.Subjects = []v1beta1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "default",
			Namespace: "default",
		},
	}
	rbac.RoleRef = v1beta1.RoleRef{
		Kind:     "ClusterRole",
		Name:     "cluster-admin",
		APIGroup: "rbac.authorization.k8s.io",
	}

	// assign to config
	config.ClusterRoleBinding = rbac

	// write params to disk as a config file
	if err := configManager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if b, err := json.MarshalIndent(config.ClusterRoleBinding, "", "  "); err != nil {
		t.Fatal(err)
	} else {
		if err := ioutil.WriteFile(key+".json", b, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
