package repository

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	ch = "current_hash"
)

type repository struct {
	db *redis.Client
}

func NewRepository(db *redis.Client) blockchain.Repository {
	return &repository{db: db}
}

func (r *repository) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	logrus.Info(v)
	return []byte(v), nil
}

func (r *repository) Add(ctx context.Context, block *domain.Block) error {
	tx := r.db.TxPipeline()

	tx.Set(ctx, string(block.Hash), block.Serialize(), 0)
	tx.Set(ctx, ch, block.Hash, 0)

	_, err := tx.Exec(ctx)
	return err
}
