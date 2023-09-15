package network

import (
	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
)

type Network struct{}

func NewNetwork() *Network {
	return &Network{}
}

// TODO: implement gossip protocol
func (n *Network) Broadcast(conf config.Config, block *domain.Block) {}
