// pods implements kube.Coder interface for deployment of pods
package pods

import (
	"context"

	"github.com/sdeoras/configio"
	"github.com/sdeoras/kube"
)

func NewCoder(key string, configReader configio.ConfigReader, ctx context.Context) (kube.Coder, error) {
	return newCoder(key, configReader, ctx)
}

func newCoder(key string, configReader configio.ConfigReader, ctx context.Context) (*coder, error) {
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
