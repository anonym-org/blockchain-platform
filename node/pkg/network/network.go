package network

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/BakuPukul/blockchain-platform/proto"
	"google.golang.org/grpc"
)

type Network struct{}

func NewNetwork() *Network {
	return &Network{}
}

func (n *Network) Broadcast(conf config.Config, log logger.Logger, block *domain.Block) {
	for _, n := range conf.Nodes {
		cc, err := grpc.Dial(n, grpc.WithTransportCredentials(nil))
		if err != nil {
			log.Error("err dial node: ", err)
			continue
		}
		defer cc.Close()

		node := proto.NewBlockchainClient(cc)

		_, err = node.SendBlock(context.Background(), &proto.SendBlockRequest{
			PrevHash: string(block.PrevHash),
			Data:     block.Data,
		})
		if err != nil {
			log.Warn("fail to broadcast data: ", err)
			continue
		}

		log.Info("broadcast data sent")
	}
}
