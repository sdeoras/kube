package pvc

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube"
	"github.com/sdeoras/kube/kube/pv"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestNewCoder(t *testing.T) {
	log := logrus.WithField("func", "TestNewCoder").WithField("package", PackageName)

	globalCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	clientset, err := kube.GetDefaultClientSet()
	if err != nil {
		t.Fatal(err)
	}

	// config init
	keyPV := "pv_gcp"
	keyPVC := "pvc_gcp"
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src",
		"github.com", "sdeoras", "kube", ".config", "config.json")
	configManager, err := configfile.NewManager(globalCtx, configfile.OptFilePath, configFilePath)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// initialize new kube coder
	// key is needed because coder works with a config manager to retrieve config data
	// and config manager requires a key to pull config data from the backend
	coderPV, err := pv.NewCoder(globalCtx, configManager, keyPV, clientset, kube.DefaultNamespace)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}
	coder, err := NewCoder(globalCtx, configManager, keyPVC, clientset, kube.DefaultNamespace)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// create a context to trigger with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when triggerFunc() is called
	trigger, triggerFunc := context.WithCancel(context.Background())

	// create kube obj (akin to kubectl create -f file)
	trigger, err = kube.Create(trigger, kube.Sync, coderPV, coder)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// delete kube obj (akin to kubectl delete -f file)
	trigger, err = kube.Delete(trigger, kube.Async, coderPV, coder)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// trigger booting
	triggerFunc()
	// wait for done
	select {
	case err := <-coder.Error():
		t.Fatal(err)
	case <-coder.Context().Done():
		t.Fatal("coder context cancelled on error")
	case <-trigger.Done():
	case <-globalCtx.Done():
		t.Fatal("global context cancelled")
	}
}
