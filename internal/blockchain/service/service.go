package service

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/redis/go-redis/v9"
)

const (
	redisKeyCurrentHash = "current_hash"
)

type service struct {
	log        logger.Logger
	repository blockchain.Repository
}

func NewService(log logger.Logger, repository blockchain.Repository) blockchain.Usecase {
	return &service{
		log:        log,
		repository: repository,
	}
}

func (s *service) InitBlockchain(ctx context.Context) *domain.Blockchain {
	var lastHash []byte

	val, err := s.repository.Get(ctx, redisKeyCurrentHash)
	if err != nil {
		if err != redis.Nil {
			s.log.Fatal(err)
		}

		genesis := domain.Genesis()
		if err = s.repository.Add(ctx, redisKeyCurrentHash, genesis); err != nil {
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

	val, err := s.repository.Get(ctx, redisKeyCurrentHash)
	if err != nil {
		s.log.Error(err)
		return err
	}

	lastHash = append(lastHash, val...)
	newBlock := domain.NewBlock(data, lastHash)
	if err := s.repository.Add(ctx, redisKeyCurrentHash, newBlock); err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

func (s *service) Next(ctx context.Context, blockchain *domain.Blockchain) (*domain.Block, error) {
	val, err := s.repository.Get(ctx, string(blockchain.CurrentHash))
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	block := domain.Deserialize(val)
	blockchain.CurrentHash = block.PrevHash

	return block, nil
}
