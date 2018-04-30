package svc

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestNewCoder(t *testing.T) {
	log := logrus.WithField("func", "TestNewCoder").WithField("package", PackageName)

	globalCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	clientset, err := kube.GetDefaultClientSet()
	if err != nil {
		t.Fatal(err)
	}

	// config init
	key := "e1d2a37e-ba01-4b57-a35b-eefc82913e42"
	configFilePath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com/sdeoras",
		PackageName, "defaults", DefaultConfigDir, DefaultConfigFile)
	configManager, err := configfile.NewManager(globalCtx, configfile.OptFilePath, configFilePath)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// initialize new kube coder
	// key is needed because coder works with a config manager to retrieve config data
	// and config manager requires a key to pull config data from the backend
	coder, err := NewCoder(globalCtx, configManager, key)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}
	coder.Clientset(clientset, kube.DefaultNamespace)

	// create a context to start with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when startFunc() is called
	start, startFunc := context.WithCancel(context.Background())

	// create kube obj (akin to kubectl create -f file)
	created := coder.Create(start)

	// delete kube object (akin to kubectl delete -f file)
	done := coder.Delete(created)

	// start booting
	startFunc()
	// wait for done
	select {
	case err := <-coder.Error():
		t.Fatal(err)
	case <-coder.Context().Done():
		t.Fatal("coder context cancelled on error")
	case <-done.Done():
	case <-globalCtx.Done():
		t.Fatal("global context cancelled")
	}
}
