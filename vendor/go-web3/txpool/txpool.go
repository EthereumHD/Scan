package txpool

import (
	"fmt"
	"go-web3/dto"
	"go-web3/providers"
)

// Net - The Net Module
type Txpool struct {
	provider providers.ProviderInterface
}

// NewNet - Net Module constructor to set the default provider
func NewTxpool(provider providers.ProviderInterface) *Txpool {
	txpool := new(Txpool)
	txpool.provider = provider
	return txpool
}

func (txpool *Txpool) Content() (ret dto.Content, err error) {
	params := []string{}
	//pointer := &dto.Content{}
	pointer := &dto.RequestResult{Result: &dto.Content{}}

	err = txpool.provider.SendRequest(&pointer, "txpool_content", params)

	if err != nil {
		return ret, err
	}

	if c, ok := pointer.Result.(*dto.Content); ok {
		ret = *c
	} else {
		err = fmt.Errorf("UNREACHABLE CODE, maybe is not valide poc server")
	}
	return ret, nil

}
