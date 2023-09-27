package blockchain

import (
	"context"

	"github.com/anonym-org/blockchain-platform/internal/domain"
	"github.com/anonym-org/blockchain-platform/proto"
)

type Usecase interface {
	InitBlockchain(ctx context.Context) *domain.Blockchain
	AddBlock(ctx context.Context, blockchain *domain.Blockchain, data string) (*domain.Block, error)
	GetBlock(ctx context.Context, blockchain *domain.Blockchain) (*domain.Block, error)
	ListBlocks(ctx context.Context, blockchain *domain.Blockchain) (currentHash string, blocks []*proto.Block, err error)
	ReplaceBlockchain(ctx context.Context, blockchain *proto.CopyBlockchainResponse) error
}
