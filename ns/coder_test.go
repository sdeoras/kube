package ns

import (
	"context"
	"testing"
	"time"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestNewCoder(t *testing.T) {
	log := logrus.WithField("func", "TestNewCoder").WithField("package", PackageName)

	globalCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	clientset, err := kube.GetDefaultClientSet()
	if err != nil {
		t.Fatal(err)
	}

	// initialize new kube coder
	coder, err := NewCoder("", nil, globalCtx)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}
	coder.Clientset(clientset, kube.DefaultNamespace)

	config := new(Config).Init("")
	config.Namespace.Name = "test"

	if err := coder.SetConfig(config); err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// create a context to trigger with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when startFunc() is called
	trigger, startFunc := context.WithCancel(context.Background())

	// create kube obj (akin to kubectl create -f file)
	trigger = coder.Create(trigger)

	// delete kube object (akin to kubectl delete -f file)
	trigger = coder.Delete(trigger)

	// trigger it
	startFunc()
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
