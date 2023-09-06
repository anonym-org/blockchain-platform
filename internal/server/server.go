package server

import (
	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
)

type Server struct {
	config *config.Config
	log    logger.Logger
}

func NewServer(c *config.Config, l logger.Logger) *Server {
	return &Server{
		config: c,
		log:    l,
	}
}

func (s *Server) Run() {
	s.log.Info("running server on port ", s.config.Port)
}
