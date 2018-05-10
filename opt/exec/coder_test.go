package exec

import (
	"context"
	"testing"
	"time"

	"os"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestNewCoder(t *testing.T) {
	log := logrus.WithField("func", "TestNewCoder").WithField("package", PackageName)

	globalCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	coders := make([]kube.Coder, 0, 0)

	for i := 0; i < 5; i++ {
		coder, err := NewCoder(globalCtx,
			os.Stdin, os.Stdout, os.Stderr,
			"execCommand", "ls", []string{"-la", "/tmp/"}...)
		if err != nil {
			log.Error(err)
			t.Fatal(err)
		}
		coders = append(coders, coder)
	}

	// create a context to trigger with
	// note, that it is being used to trigger action when it _ends_
	// i.e., when startTrigger() is called
	trigger, startTrigger := context.WithCancel(context.Background())

	trigger, _ = kube.Create(trigger, kube.Async, coders...)

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
