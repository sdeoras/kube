// pods implements kube.Coder interface for deployment of pods
package pods

import (
	"context"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func NewCoder(ctx context.Context, configReader configio.ConfigReader, key string, clientset *kubernetes.Clientset, namespace string) (kube.Coder, error) {
	return newCoder(ctx, configReader, key, clientset, namespace)
}

func newCoder(ctx context.Context, configReader configio.ConfigReader, key string, clientset *kubernetes.Clientset, namespace string) (*coder, error) {
	cdr := new(coder)
	cdr.key = key
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)
	cdr.clientset = clientset
	cdr.namespace = namespace
	cdr.log = logrus.WithField("package", PackageName).WithField("namespace", cdr.namespace)

	cdr.err = make(chan error)

	if configReader != nil && len(key) != 0 {
		config := new(Config).Init(cdr.key)
		if err := configReader.Unmarshal(config); err != nil {
			return nil, err
		} else {
			cdr.config = config
		}
	}

	return cdr, nil
}
