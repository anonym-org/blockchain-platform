package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BakuPukul/blockchain-platform/config"
	_grpc "github.com/BakuPukul/blockchain-platform/internal/blockchain/delivery/grpc"
	_http "github.com/BakuPukul/blockchain-platform/internal/blockchain/delivery/http"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain/repository"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain/service"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
	"github.com/BakuPukul/blockchain-platform/pkg/network"
	"github.com/BakuPukul/blockchain-platform/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Server struct {
	config  *config.Config
	log     logger.Logger
	db      *redis.Client
	handler *http.ServeMux
	gs      *grpc.Server
}

func NewServer(c *config.Config, l logger.Logger, db *redis.Client) *Server {
	return &Server{
		config:  c,
		log:     l,
		db:      db,
		handler: http.NewServeMux(),
		gs:      grpc.NewServer(),
	}
}

func (s *Server) Run() error {
	repository := repository.NewRepository(s.db)
	service := service.NewService(s.log, repository)
	chain := service.InitBlockchain(context.TODO())
	network := network.NewNetwork(*s.config, s.log)

	_http.NewController(*s.config, chain, s.handler, service, network)
	srv := _grpc.NewHandler(*s.config, s.log, chain, service)

	proto.RegisterBlockchainServer(s.gs, srv)
	grpcListener, err := net.Listen("tcp", s.config.GrpcPort)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:         s.config.Port,
		Handler:      s.handler,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	// GRPC server
	go func() {
		s.log.Info("grpc server running on port", s.config.GrpcPort)
		if err := s.gs.Serve(grpcListener); err != nil {
			s.log.Fatal(err)
		}
	}()

	// HTTP (REST API) server
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
