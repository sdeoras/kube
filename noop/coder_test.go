package noop

import (
	"context"
	"testing"
	"time"

	"fmt"

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
		coder, err := NewCoder(fmt.Sprintf("coder_%d", i), globalCtx)
		if err != nil {
			log.Error(err)
			t.Fatal(err)
		}

		coders = append(coders, coder)
	}

	spacer, err := NewCoder("-------", globalCtx)
	if err != nil {
		log.Error(err)
		t.Fatal(err)
	}

	// create a context to start with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when startFunc() is called
	start, startFunc := context.WithCancel(context.Background())

	start = coders[0].Create(start)
	start = spacer.Create(start)

	coders[1].Create(start)
	coders[2].Create(start)
	start = coders[3].Create(start)
	start = spacer.Create(start)

	coders[4].Create(start)
	coders[5].Create(start)
	start = coders[6].Create(start)
	start = spacer.Create(start)

	start = coders[7].Create(start)
	start = coders[8].Create(start)
	done := coders[9].Create(start)

	// start booting
	startFunc()
	// wait for done
	select {
	case <-done.Done():
	case <-globalCtx.Done():
		t.Fatal("global context cancelled")
	}

	trigger, stop := context.WithCancel(context.Background())
	if done, err := kube.Shutdown(trigger, kube.Async, coders...); err != nil {
		t.Fatal(err)
	} else {
		stop()
		select {
		case <-done.Done():
		case <-globalCtx.Done():
			t.Fatal("global context cancelled")
		}
	}

	time.Sleep(time.Second)
}
