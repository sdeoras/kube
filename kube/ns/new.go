// ns implements kube.Coder interface for deployment of namespaces
package ns

import (
	"context"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
)

func NewCoder(ctx context.Context, configReader configio.ConfigReader, key string) (kube.Coder, error) {
	return newCoder(ctx, configReader, key)
}

func New(ctx context.Context, namespace string) (kube.Coder, error) {
	coder, err := newCoder(ctx, nil, "")
	if err != nil {
		return nil, err
	}
	config := new(Config).Init("")
	config.Namespace.Name = namespace
	if err := coder.SetConfig(config); err != nil {
		return nil, err
	}
	return coder, nil
}

func newCoder(ctx context.Context, configReader configio.ConfigReader, key string) (*coder, error) {
	cdr := new(coder)
	cdr.key = key
	cdr.ctx, cdr.cancel = context.WithCancel(ctx)

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
