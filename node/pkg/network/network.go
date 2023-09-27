package network

import (
	"context"
	"errors"

	"github.com/anonym-org/blockchain-platform/config"
	"github.com/anonym-org/blockchain-platform/internal/domain"
	"github.com/anonym-org/blockchain-platform/pkg/logger"
	"github.com/anonym-org/blockchain-platform/proto"
	"google.golang.org/grpc"
)

type Network struct {
	conf config.Config
	log  logger.Logger
}

func NewNetwork(conf config.Config, log logger.Logger) *Network {
	return &Network{conf: conf, log: log}
}

func (n *Network) Broadcast(block *domain.Block) (err error) {
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
	return
}

func (n *Network) DownloadBlockchain(conf config.Config, log logger.Logger) (*proto.CopyBlockchainResponse, error) {
	if len(n.conf.Nodes) < 1 {
		n.log.Error("download blockchain: no peer nodes")
		return nil, errors.New("download blockchain: no peer nodes")
	}

	cc, err := grpc.Dial(n.conf.Nodes[0], grpc.WithInsecure())
	if err != nil {
		n.log.Error("err dial node: ", err)
		return nil, err
	}
	defer cc.Close()

	node := proto.NewBlockchainClient(cc)

	resp, err := node.CopyBlockchain(context.TODO(), &proto.CopyBlockchainRequest{})
	if err != nil {
		n.log.Error("download blockchain: ", err)
		return nil, err
	}

	return resp, nil
}
