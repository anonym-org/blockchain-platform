package blockchain

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/internal/domain"
)

type Repository interface {
	Get(ctx context.Context, key string) (string, error)
	Add(ctx context.Context, currentHashKey string, block *domain.Block) error
}
