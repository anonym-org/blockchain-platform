package grpc

import (
	"context"

	"github.com/anonym-org/blockchain-platform/config"
	"github.com/anonym-org/blockchain-platform/internal/blockchain"
	"github.com/anonym-org/blockchain-platform/internal/domain"
	"github.com/anonym-org/blockchain-platform/pkg/logger"
	"github.com/anonym-org/blockchain-platform/pkg/network"
	"github.com/anonym-org/blockchain-platform/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	conf       config.Config
	log        logger.Logger
	blockchain *domain.Blockchain
	service    blockchain.Usecase
	network    *network.Network
	proto.UnimplementedBlockchainServer
}

func NewHandler(conf config.Config, log logger.Logger, blockchain *domain.Blockchain, service blockchain.Usecase, network *network.Network) proto.BlockchainServer {
	return &handler{
		conf:       conf,
		log:        log,
		blockchain: blockchain,
		service:    service,
		network:    network,
	}
}

func (h *handler) SendBlock(ctx context.Context, r *proto.SendBlockRequest) (*proto.SendBlockResponse, error) {
	isValid, err := h.ValidateBlockchain(ctx)
	if err != nil {
		return &proto.SendBlockResponse{}, nil
	}

	if !isValid {
		h.log.Error("node blockchain were invalid, will copy full blockchain from nodes")

		chain, err := h.network.DownloadBlockchain(h.conf, h.log)
		if err != nil {
			h.log.Error("node blockchain invalid, fail to repair blockchain")
			return &proto.SendBlockResponse{}, status.Error(codes.Internal, "node blockchain invalid, fail to repair blockchain")
		}
		h.blockchain = &domain.Blockchain{CurrentHash: chain.CurrentHash}
		if err = h.service.ReplaceBlockchain(ctx, chain); err != nil {
			return &proto.SendBlockResponse{}, status.Error(codes.Internal, "node blockchain invalid, fail to repair blockchain")
		}

		return &proto.SendBlockResponse{}, status.Error(codes.Internal, "node blockchain were invalid, will copy full blockchain from nodes")
	}

	if _, err := h.service.AddBlock(ctx, h.blockchain, r.Data); err != nil {
		h.log.Error("fail to add block into node blockchain")
		return &proto.SendBlockResponse{}, status.Error(codes.Internal, "fail to add block into node blockchain")
	}

	return &proto.SendBlockResponse{}, nil
}

func (h *handler) CopyBlockchain(context.Context, *proto.CopyBlockchainRequest) (*proto.CopyBlockchainResponse, error) {
	currentHash, blocks, err := h.service.ListBlocks(context.TODO(), h.blockchain)
	if err != nil {
		h.log.Error("fail to get node blockchain")
		return &proto.CopyBlockchainResponse{}, status.Error(codes.Internal, "fail to get full copy of blockchain")
	}

	return &proto.CopyBlockchainResponse{
		CurrentHash: currentHash,
		Blocks:      blocks,
	}, nil
}

func (h *handler) ValidateBlockchain(ctx context.Context) (bool, error) {
	_, blocks, err := h.service.ListBlocks(ctx, h.blockchain)
	if err != nil {
		return false, err
	}

	reversed := make([]*proto.Block, len(blocks))
	for i, block := range blocks {
		reversed[len(blocks)-(i+1)] = &proto.Block{
			Hash:     block.Hash,
			Data:     block.Data,
			PrevHash: block.PrevHash,
		}
	}

	hash := h.blockchain.GenesisHash
	for i, block := range reversed {
		// skip genesis block
		if i == 0 {
			continue
		}
		if block.PrevHash != hash {
			return false, nil
		}
		hash = block.Hash
	}
	return true, nil
}
