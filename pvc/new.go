// pvc implements kube.Coder interface for deployment of persistent volume claims
package pvc

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

	config := new(Config).Init(cdr.key)
	if err := configReader.Unmarshal(config); err != nil {
		return nil, err
	} else {
		cdr.config = config
	}
	return cdr, nil
}
