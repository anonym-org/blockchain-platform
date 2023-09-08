package service

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/redis/go-redis/v9"
)

const (
	lastHashKey = "last_hash"
)

type service struct {
	log logger.Logger
	db  *redis.Client
}

func NewService(log logger.Logger, db *redis.Client) blockchain.Usecase {
	return &service{
		log: log,
		db:  db,
	}
}

func (s *service) InitBlockchain(ctx context.Context) *domain.Blockchain {
	var lastHash []byte

	val, err := s.db.Get(ctx, lastHashKey).Result()
	if err != nil {
		if err != redis.Nil {
			s.log.Fatal(err)
		}

		genesis := domain.Genesis()

		err = s.db.Set(ctx, string(genesis.Hash), genesis.Serialize(), 0).Err()
		if err != nil {
			s.log.Fatal(err)
		}

		err = s.db.Set(ctx, lastHashKey, genesis.Hash, 0).Err()
		if err != nil {
			s.log.Fatal(err)
		}

		lastHash = genesis.Hash
	}

	lastHash = append(lastHash, []byte(val)...)
	blockchain := domain.Blockchain{
		CurrentHash: lastHash,
	}
	return &blockchain
}

func (s *service) AddBlock(ctx context.Context, data string) error {
	var lastHash []byte

	val, err := s.db.Get(ctx, lastHashKey).Result()
	if err != nil {
		s.log.Error(err)
		return err
	}

	lastHash = append(lastHash, []byte(val)...)
	newBlock := domain.NewBlock(data, lastHash)
	if err := s.db.Set(ctx, string(newBlock.Hash), newBlock.Serialize(), 0).Err(); err != nil {
		s.log.Error(err)
		return err
	}
	if err := s.db.Set(ctx, lastHashKey, newBlock.Hash, 0).Err(); err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

func (s *service) Next(ctx context.Context, blockchain *domain.Blockchain) (*domain.Block, error) {
	val, err := s.db.Get(ctx, string(blockchain.CurrentHash)).Result()
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	block := domain.Deserialize([]byte(val))
	blockchain.CurrentHash = block.PrevHash

	return block, nil
}
