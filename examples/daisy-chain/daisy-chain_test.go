package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sdeoras/configio/configfile"
	"github.com/sdeoras/kube"
	"github.com/sdeoras/kube/ds"
	"github.com/sdeoras/kube/jobs"
	"github.com/sdeoras/kube/pods"
	"github.com/sdeoras/kube/pv"
	"github.com/sdeoras/kube/pvc"
	"github.com/sdeoras/kube/svc"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDaisyChain(t *testing.T) {
	log := logrus.WithField("func", "main").WithField("package", "main")

	// config file locations
	configFile := filepath.Join(os.Getenv("HOME"), ".config", "tf", "inception", "config.json")
	kubeConfigFile := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// keys for config data... config managers can pull data from config files based on these keys
	// config data is able to populate k8s objects
	keys := make(map[string]string)
	keys["pv"] = "b773bde3-e73b-4914-8fac-3513ca76a596"
	keys["pvc"] = "e614bfad-436e-4e26-b6b5-41384f2260e6"
	keys["pods"] = "86f96730-44c7-4942-ba59-9cc711143ffa"
	keys["svc"] = "58a60855-ca76-45c8-a8cf-ad2b03362db9"
	keys["ds"] = "4d415f93-0a19-4037-839e-00bd7b049eae"
	keys["jobs"] = "567605cf-3ae7-41ef-8a49-31bb1f6bb3cc"

	// create a global context
	// global context flows everywhere
	globalCtx, globalCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer globalCancel()

	// create a context to trigger with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when startTrigger() is called, so it is best to init it with context.Background()
	trigger, startTrigger := context.WithCancel(context.Background())

	// kubernetes clientset init
	var clientset *kubernetes.Clientset
	// use the current context in kubeconfig
	if kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile); err != nil {
		log.Error(err)
		t.Fatal(err)
	} else {
		// create the clientset
		clientset, err = kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			log.Error(err)
			t.Fatal(err)
		}
	}

	// get config manager for fetching config data
	configManager, err := configfile.NewManager(globalCtx, configfile.OptFilePath, configFile)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}
	defer configManager.Close()

	// build a list of kube coders
	// kube coders can be daisy chained used context and have methods to
	// create and delete k8s objects
	coders := make([]kube.Coder, 0, 0)

	// get coder for pv
	if coder, err := pv.NewCoder(keys["pv"], configManager, globalCtx); err != nil {
		log.Error(err)
		t.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for pvc with a key and config manager
	if coder, err := pvc.NewCoder(keys["pvc"], configManager, globalCtx); err != nil {
		log.Error(err)
		t.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for daemon set with a key and config manager
	if coder, err := ds.NewCoder(keys["ds"], configManager, globalCtx); err != nil {
		log.Error(err)
		t.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for pods with a key and config manager
	if coder, err := pods.NewCoder(keys["pods"], configManager, globalCtx); err != nil {
		log.Error(err)
		t.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for svc with a key and config manager
	if coder, err := svc.NewCoder(keys["svc"], configManager, globalCtx); err != nil {
		log.Error(err)
		t.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for jobs with a key and config manager
	if coder, err := jobs.NewCoder(keys["jobs"], configManager, globalCtx); err != nil {
		log.Error(err)
		t.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// assign clientset to kube coders
	done := make(chan struct{})
	for _, coder := range coders {
		coder := coder

		// set clientset and namespace for coder
		coder.Clientset(clientset, kube.DefaultNamespace)

		// fan in coder context into done
		go func(cdr kube.Coder) {
			done <- <-cdr.Context().Done()
		}(coder)
	}

	// boot processes by passing parent context as a trigger
	for i := 0; i < len(coders); i++ {
		trigger = coders[i].Create(trigger)
	}

	// wait for a few seconds
	trigger, trigCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer trigCancel()

	// shutdown processes in a sequence by passing parent context as trigger
	for i := len(coders) - 1; i >= 0; i-- {
		trigger = coders[i].Delete(trigger)
	}

	// send trigger to initiate k8s actions
	startTrigger()

	// wait for various events
	// at least one of them will surely happen
	select {
	case <-trigger.Done():
		log.Info("all done")
	case <-done:
		log.Error("at least one of the coder contexts is done")
	case <-globalCtx.Done():
		log.Fatal("global context cancelled")
	}
}
