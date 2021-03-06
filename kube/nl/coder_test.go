package nl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestNewCoder(t *testing.T) {
	log := logrus.WithField("func", "TestNewCoder").WithField("package", PackageName)

	globalCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// config init
	key := "nl"

	clientset, err := kube.GetDefaultClientSet()
	if err != nil {
		t.Fatal(err)
	}

	// initialize new kube coder
	// key is needed because coder works with a config manager to retrieve config data
	// and config manager requires a key to pull config data from the backend
	coder, err := NewCoder(globalCtx, key, clientset)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// create a context to start with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when start() is called
	trigger, start := context.WithCancel(context.Background())

	// create kube obj (akin to kubectl create -f file)
	trigger = coder.Create(trigger)

	// start booting
	start()

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

	config := coder.GetConfig().(*Config)
	for _, node := range config.NodeList.Items {
		fmt.Println(node.Name)
	}
}
