package repository

import (
	"context"
	"encoding/json"

	"github.com/anonym-org/blockchain-platform/config"
	"github.com/anonym-org/blockchain-platform/internal/blockchain"
	"github.com/anonym-org/blockchain-platform/internal/domain"
	"github.com/dgraph-io/badger/v4"
)

type repository struct {
	conf *config.Config
	db   *badger.DB
}

func NewRepository(conf *config.Config, db *badger.DB) blockchain.Repository {
	return &repository{conf: conf, db: db}
}

func (r *repository) Get(ctx context.Context, key string) (string, error) {
	var result string
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			result = string(val[:])
			return nil
		})
		return err
	})
	return result, err
}

func (r *repository) Add(ctx context.Context, currentHashKey string, block *domain.Block) error {
	return r.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(block.Hash), block.Serialize()); err != nil {
			return err
		}
		err := txn.Set([]byte(currentHashKey), []byte(block.Hash))
		return err
	})
}

func (r *repository) Set(ctx context.Context, key string, val any) error {
	return r.db.Update(func(txn *badger.Txn) error {
		v, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), v)
	})
}

func (r *repository) Clear(ctx context.Context) error {
	return r.db.DropAll()
}
