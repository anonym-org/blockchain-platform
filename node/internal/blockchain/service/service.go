package service

import (
	"context"

	"github.com/anonym-org/blockchain-platform/internal/blockchain"
	"github.com/anonym-org/blockchain-platform/internal/domain"
	"github.com/anonym-org/blockchain-platform/pkg/logger"
	"github.com/anonym-org/blockchain-platform/proto"
	"github.com/dgraph-io/badger/v4"
)

const (
	CurrentHashKey = "current_hash"
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
	var prevHash string

	genesis := domain.Genesis()
	val, err := s.repository.Get(ctx, CurrentHashKey)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			s.log.Fatal(err)
		}

		s.log.Info(genesis)
		if err = s.repository.Add(ctx, CurrentHashKey, genesis); err != nil {
			s.log.Fatal(err)
		}

		prevHash = genesis.Hash
	} else {
		prevHash = val
	}

	return &domain.Blockchain{
		GenesisHash: genesis.Hash,
		CurrentHash: prevHash,
	}
}

func (s *service) AddBlock(ctx context.Context, blockchain *domain.Blockchain, data string) (*domain.Block, error) {
	prevHash, err := s.repository.Get(ctx, CurrentHashKey)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	newBlock := domain.NewBlock(data, prevHash)
	if err := s.repository.Add(ctx, CurrentHashKey, newBlock); err != nil {
		s.log.Error(err)
		return nil, err
	}
	blockchain.CurrentHash = newBlock.Hash
	return newBlock, nil
}

func (s *service) GetBlock(ctx context.Context, blockchain *domain.Blockchain) (*domain.Block, error) {
	val, err := s.repository.Get(ctx, blockchain.CurrentHash)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	block := domain.Deserialize([]byte(val))
	return block, nil
}

func (s *service) ReplaceBlockchain(ctx context.Context, blockchain *proto.CopyBlockchainResponse) error {
	s.log.Info("REPLACING BLOCKCHAIN....")
	s.log.Info(blockchain.Blocks)

	var currentHash string

	if err := s.repository.Clear(ctx); err != nil {
		s.log.Error(err)
		return err
	}

	for i, b := range blockchain.Blocks {
		if err := s.repository.Add(ctx, CurrentHashKey, &domain.Block{
			Hash:     b.Hash,
			Data:     b.Data,
			PrevHash: b.PrevHash,
		}); err != nil {
			s.log.Error(err)
			return err
		}

		if i == 0 {
			currentHash = b.Hash
		}
	}

	s.repository.Set(ctx, CurrentHashKey, currentHash)

	return nil
}

func (s *service) ListBlocks(ctx context.Context, blockchain *domain.Blockchain) (string, []*proto.Block, error) {
	var currentHash string
	prevHash := blockchain.CurrentHash
	blocks := []*proto.Block{}

	for prevHash != "" {
		val, err := s.repository.Get(ctx, prevHash)
		if err != nil {
			s.log.Error(err)
			return "", nil, err
		}

		if currentHash == "" {
			currentHash = val
		}

		block := domain.Deserialize([]byte(val))
		prevHash = block.PrevHash
		blocks = append(blocks, &proto.Block{
			Hash:     block.Hash,
			Data:     block.Data,
			PrevHash: block.PrevHash,
			Nounce: block.Nounce,
		})
	}

	return currentHash, blocks, nil
}
