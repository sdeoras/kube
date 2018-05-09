package sleep

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

	globalCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	coders := make([]kube.Coder, 0, 0)

	for i := 0; i < 10; i++ {
		coder, err := NewCoder(globalCtx, fmt.Sprintf("coder_%d", i), time.Second)
		if err != nil {
			log.Error(err)
			t.Fatal(err)
		}

		coders = append(coders, coder)
	}

	spacer, err := NewCoder(globalCtx, "-------", 0)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// create a context to trigger with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when startTrigger() is called
	trigger, startTrigger := context.WithCancel(context.Background())

	trigger = coders[0].Create(trigger)
	trigger = spacer.Create(trigger)

	// fan in various execution contexts
	trigger = kube.FanIn(
		coders[1].Create(trigger),
		coders[2].Create(trigger),
		coders[3].Create(trigger))
	trigger = spacer.Create(trigger)

	/*
		// alternatively execute using kube.Create
		trigger, _ = kube.Create(trigger, kube.Async, coders[1:4]...)
		trigger = spacer.Create(trigger)
	*/

	trigger, _ = kube.Create(trigger, kube.Async, coders[4:7]...)
	trigger = spacer.Create(trigger)

	trigger, _ = kube.Create(trigger, kube.Sync,
		coders[7], spacer,
		coders[8], spacer,
		coders[9], spacer)

	// then delete em all asynchronously
	trigger, _ = kube.Delete(trigger, kube.Async, coders...)

	// start trigger
	// without this function nothing will execute
	startTrigger()

	select {
	case <-trigger.Done():
	case <-globalCtx.Done():
		t.Fatal("global context cancelled")
	}

	time.Sleep(time.Second)
}
