package repository

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/BakuPukul/blockchain-platform/pkg/db"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	conf *config.Config
	db   *redis.Client
}

func NewRepository(conf *config.Config, db *redis.Client) blockchain.Repository {
	return &repository{conf: conf, db: db}
}

func (r *repository) Get(ctx context.Context, key string) (string, error) {
	v, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return v, nil
}

func (r *repository) Add(ctx context.Context, currentHashKey string, block *domain.Block) error {
	tx := r.db.TxPipeline()

	tx.Set(ctx, block.Hash, block.Serialize(), 0)
	tx.Set(ctx, currentHashKey, block.Hash, 0)

	_, err := tx.Exec(ctx)
	return err
}

func (r *repository) Set(ctx context.Context, key string, val any) error {
	_, err := r.db.Set(ctx, key, val, 0).Result()
	return err
}

func (r *repository) Clear(ctx context.Context) error {
	if err := r.db.FlushDB(ctx).Err(); err != nil {
		return err
	}
	if err := r.db.Close(); err != nil {
		return err
	}

	client, err := db.NewRedis(r.conf)
	if err != nil {
		return err
	}

	r.db = client

	return r.db.Do(context.Background(), "SELECT", 0).Err()
}
