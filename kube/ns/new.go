// ns implements kube.Coder interface for deployment of namespaces
package ns

import (
	"context"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func NewCoder(ctx context.Context, clientset *kubernetes.Clientset, namespace string) (kube.Coder, error) {
	return newCoder(ctx, clientset, namespace)
}

func newCoder(ctx context.Context, clientset *kubernetes.Clientset, namespace string) (*coder, error) {
	cdr := new(coder)
	cdr.key = namespace
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)
	cdr.clientset = clientset
	cdr.namespace = namespace
	cdr.log = logrus.WithField("package", PackageName).WithField("namespace", cdr.namespace)

	cdr.config = new(Config).Init(cdr.key)
	cdr.err = make(chan error)

	return cdr, nil
}
