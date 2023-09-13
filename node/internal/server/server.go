package server

import (
	"context"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain/repository"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain/service"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	config *config.Config
	log    logger.Logger
	db     *redis.Client
}

func NewServer(c *config.Config, l logger.Logger, db *redis.Client) *Server {
	return &Server{
		config: c,
		log:    l,
		db:     db,
	}
}

func (s *Server) Run() {
	s.log.Info("running server on port ", s.config.Port)

	repository := repository.NewRepository(s.db)
	service := service.NewService(s.log, repository)

	chain := service.InitBlockchain(context.TODO())
	block, err := service.Next(context.Background(), chain)
	if err != nil {
		s.log.Error(err)
	}
	s.log.Info(block.Data)

	// ex: add new block
	service.AddBlock(context.TODO(), "First block after genesis")

	// ex: validate block
	// pow := domain.NewProof(block)
	// s.log.Info(pow.Validate())
}
