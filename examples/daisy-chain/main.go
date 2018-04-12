package main

import (
	"context"
	"flag"
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
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	shutdown := flag.Bool("shutdown", false, "(Bool)\tshutdown")
	flag.Parse()

	log := logrus.WithField("func", "main").WithField("package", "main")

	// keys for config data... config managers can pull data from config files based on these keys
	// config data is able to populate k8s objects
	keys := make(map[string]string)
	configFile := initKeys(keys)

	// create a global context with half hour timeout
	// global context flows everywhere
	globalCtx, globalCancel := context.WithTimeout(context.Background(), time.Second*60*30)
	defer globalCancel()

	// create a context to trigger with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when start() is called, so it is best to init it with context.Background()
	trigger, start := context.WithCancel(context.Background())
	errCtx, doCleanup := context.WithCancel(context.Background())

	// get new kubernetes clientset from default config file
	clientset, err := newClientset()
	if err != nil {
		log.Fatal(err)
	}

	// get config manager for fetching config data
	configManager, err := configfile.NewManager(globalCtx, configfile.OptFilePath, configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer configManager.Close()

	// build a list of kube coders
	coders := make([]kube.Coder, 0, 0)

	// get coder for pv
	if coder, err := pv.NewCoder(keys["pv"], configManager, globalCtx); err != nil {
		log.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for pvc with a key and config manager
	if coder, err := pvc.NewCoder(keys["pvc"], configManager, globalCtx); err != nil {
		log.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for daemon set with a key and config manager
	if coder, err := ds.NewCoder(keys["ds"], configManager, globalCtx); err != nil {
		log.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for pods with a key and config manager
	if coder, err := pods.NewCoder(keys["pods"], configManager, globalCtx); err != nil {
		log.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for svc with a key and config manager
	if coder, err := svc.NewCoder(keys["svc"], configManager, globalCtx); err != nil {
		log.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// get coder for jobs with a key and config manager
	if coder, err := jobs.NewCoder(keys["jobs"], configManager, globalCtx); err != nil {
		log.Fatal(err)
	} else {
		coders = append(coders, coder)
	}

	// assign clientset to each of the coders and fan in their contexts and errors
	kube.AssignClientSet(clientset, kube.DefaultNamespace, coders...)
	coderCtx, coderErrChan := kube.FanIn(coders...)

	// boot processes by passing parent context as a trigger
	if !*shutdown {
		trigger = kube.Bootup(trigger, kube.Forward, coders...)
	}

	/*// wait for a few seconds
	trigger, trigCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer trigCancel()*/

	// shutdown processes in a sequence by passing parent context as trigger
	if *shutdown {
		trigger = kube.Shutdown(trigger, kube.Backward, coders...)
	}

	// arm cleanup action on coders on a trigger
	// cleanup does not happen till context triggers
	doneCleanup := kube.Cleanup(errCtx, coders...)

	// start sends triggers to processes waiting to start
	start()

	// wait for various events
	// at least one of them will surely happen
	select {
	// at least one of the coder errors out
	case err := <-coderErrChan:
		// in which case trigger cleanup
		doCleanup()
		// and wait for cleanup done
		select {
		// either cleanup succeeds
		case <-doneCleanup:
			log.Info("all cleanup done")
			// or global context is over
		case <-globalCtx.Done():
			log.Error("global context cancelled")
		}
		log.Fatal(err)
		// everything completes successfully
	case <-trigger.Done():
		log.Info("all done")
	case <-coderCtx:
		log.Error("at least one of the coder contexts is done")
	case <-globalCtx.Done():
		log.Fatal("global context cancelled")
	}
}
