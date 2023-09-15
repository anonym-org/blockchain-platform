package blockchain

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/internal/domain"
)

type Usecase interface {
	InitBlockchain(ctx context.Context) *domain.Blockchain
	AddBlock(ctx context.Context, data string) error
	GetBlock(ctx context.Context, blockchain *domain.Blockchain) (*domain.Block, error)
}
