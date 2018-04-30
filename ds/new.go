// ds implements kube.Coder interface for deployment of daemon set
package ds

import (
	"context"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
)

func NewCoder(ctx context.Context, configReader configio.ConfigReader, key string) (kube.Coder, error) {
	return newCoder(ctx, configReader, key)
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
