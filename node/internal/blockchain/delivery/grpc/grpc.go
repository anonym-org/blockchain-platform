package grpc

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/BakuPukul/blockchain-platform/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	conf       config.Config
	log        logger.Logger
	blockchain *domain.Blockchain
	service    blockchain.Usecase
	proto.UnimplementedBlockchainServer
}

func NewHandler(conf config.Config, log logger.Logger, blockchain *domain.Blockchain, service blockchain.Usecase) proto.BlockchainServer {
	return &handler{
		conf:       conf,
		log:        log,
		blockchain: blockchain,
		service:    service,
	}
}

func (h *handler) SendBlock(ctx context.Context, r *proto.SendBlockRequest) (*proto.SendBlockResponse, error) {
	block, err := h.service.GetBlock(ctx, h.blockchain)
	if err != nil {
		h.log.Error("fail to get node blockchain")
		return &proto.SendBlockResponse{}, status.Error(codes.Internal, "fail to get node blockchain")
	}

	h.log.Info("compare: ", block.Hash, " with req: ", r.PrevHash)
	if block.Hash != r.PrevHash {
		h.log.Error("node blockchain were invalid, will copy full blockchain from node")
		// TODO: copy full blockchain from other node
		return &proto.SendBlockResponse{}, status.Error(codes.Internal, "node blockchain were invalid, will copy full blockchain from node")
	}

	if _, err = h.service.AddBlock(ctx, h.blockchain, r.Data); err != nil {
		h.log.Error("fail to add block into node blockchain")
		return &proto.SendBlockResponse{}, status.Error(codes.Internal, "fail to add block into node blockchain")
	}

	return &proto.SendBlockResponse{}, nil
}
