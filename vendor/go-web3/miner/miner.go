package miner

import (
	"go-web3/providers"
)

// Net - The Net Module
type Miner struct {
	provider providers.ProviderInterface
}

// NewNet - Net Module constructor to set the default provider
func NewMiner(provider providers.ProviderInterface) *Miner {
	miner := new(Miner)
	miner.provider = provider
	return miner
}