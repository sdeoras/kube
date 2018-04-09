package pv

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sdeoras/configio/configfile"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func TestNewManager(t *testing.T) {
	log := logrus.WithField("func", "TestNewManager").WithField("package", "pv")

	globalCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// kubernetes clientset init
	var clientset *kubernetes.Clientset
	kubeConfigFile := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// use the current context in kubeconfig
	if kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile); err != nil {
		t.Fatal(err)
	} else {
		// create the clientset
		clientset, err = kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			t.Fatal(err)
		}
	}

	// config init
	key := "b773bde3-e73b-4914-8fac-3513ca76a596"
	config := new(Config).Init(key)
	configFilePath := filepath.Join(os.Getenv("HOME"), DefaultConfigDir, DefaultConfigFile)
	configManager, err := configfile.NewManager(globalCtx, configfile.OptFilePath, configFilePath)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// initialize new kube coder
	// key is needed because coder works with a config manager to retrieve config data
	// and config manager requires a key to pull config data from the backend
	coder := NewCoder(config.Key(), globalCtx)
	if err := coder.Init(clientset, configManager); err != nil {
		log.Error(err)
		t.Fatal(err)
	}

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
	case <-coder.Context().Done():
		t.Fatal("coder context cancelled on error")
	case <-done.Done():
	case <-globalCtx.Done():
		t.Fatal("global context cancelled")
	}
}
