package service

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/redis/go-redis/v9"
)

const (
	REDIS_KEY_PREVIOUS_HASH = "prev_hash"
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

	val, err := s.repository.Get(ctx, REDIS_KEY_PREVIOUS_HASH)
	if err != nil {
		if err != redis.Nil {
			s.log.Fatal(err)
		}

		genesis := domain.Genesis()
		if err = s.repository.Add(ctx, REDIS_KEY_PREVIOUS_HASH, genesis); err != nil {
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

func (s *service) AddBlock(ctx context.Context, blockchain *domain.Blockchain, data string) (*domain.Block, error) {
	var lastHash []byte

	val, err := s.repository.Get(ctx, REDIS_KEY_PREVIOUS_HASH)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	lastHash = append(lastHash, val...)
	newBlock := domain.NewBlock(data, lastHash)
	if err := s.repository.Add(ctx, REDIS_KEY_PREVIOUS_HASH, newBlock); err != nil {
		s.log.Error(err)
		return nil, err
	}
	blockchain.CurrentHash = newBlock.Hash

	return newBlock, nil
}

func (s *service) GetBlock(ctx context.Context, blockchain *domain.Blockchain) (*domain.Block, error) {
	val, err := s.repository.Get(ctx, string(blockchain.CurrentHash))
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	block := domain.Deserialize(val)
	return block, nil
}
