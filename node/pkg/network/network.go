package network

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/BakuPukul/blockchain-platform/proto"
	"google.golang.org/grpc"
)

type Network struct {
	conf config.Config
	log  logger.Logger
}

func NewNetwork(conf config.Config, log logger.Logger) *Network {
	return &Network{conf: conf, log: log}
}

func (n *Network) Broadcast(block *domain.Block) {
	n.log.Info(block.PrevHash)
	n.log.Info(block.Data)
	for _, v := range n.conf.Nodes {
		cc, err := grpc.Dial(v, grpc.WithInsecure())
		if err != nil {
			n.log.Error("err dial node: ", err)
			continue
		}
		defer cc.Close()

		node := proto.NewBlockchainClient(cc)

		_, err = node.SendBlock(context.Background(), &proto.SendBlockRequest{
			PrevHash: block.PrevHash,
			Data:     block.Data,
		})
		if err != nil {
			n.log.Warn("fail to broadcast data: ", err)
			continue
		}

		n.log.Info("broadcast data sent")
	}
}
