package main

import (
	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/server"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
)

func main() {
	l := logger.NewLogger()

	c, err := config.LoadConfig("config/config")
	if err != nil {
		l.Fatal(err)
	}

	s := server.NewServer(c, l)
	s.Run()
}
