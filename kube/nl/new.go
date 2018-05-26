// nl implements kube.Coder interface for nodelist
package nl

import (
	"context"

	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func NewCoder(ctx context.Context, key string, clientset *kubernetes.Clientset) (kube.Coder, error) {
	return newCoder(ctx, key, clientset)
}

func newCoder(ctx context.Context, key string, clientset *kubernetes.Clientset) (*coder, error) {
	cdr := new(coder)
	cdr.key = key
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)
	cdr.clientset = clientset
	cdr.log = logrus.WithField("package", PackageName)

	cdr.err = make(chan error)
	cdr.config = new(Config).Init(key)

	return cdr, nil
}
