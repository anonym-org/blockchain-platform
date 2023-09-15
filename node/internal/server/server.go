package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain/controller"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain/repository"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain/service"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/BakuPukul/blockchain-platform/pkg/network"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	config  *config.Config
	log     logger.Logger
	db      *redis.Client
	handler *http.ServeMux
}

func NewServer(c *config.Config, l logger.Logger, db *redis.Client) *Server {
	return &Server{
		config:  c,
		log:     l,
		db:      db,
		handler: http.NewServeMux(),
	}
}

func (s *Server) Run() error {
	s.log.Info("running server on port ", s.config.Port)

	repository := repository.NewRepository(s.db)
	service := service.NewService(s.log, repository)
	chain := service.InitBlockchain(context.TODO())
	network := network.NewNetwork()
	controller.NewController(*s.config, chain, s.handler, service, network)

	server := &http.Server{
		Addr:         s.config.Port,
		Handler:      s.handler,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	go func() {
		s.log.Info("http server running on port", s.config.Port)
		if err := server.ListenAndServe(); err != nil {
			s.log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)

	sig := <-quit
	s.log.Warn("received terminate, graceful shutdown ", sig)

	tc, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdown()

	return server.Shutdown(tc)
}
