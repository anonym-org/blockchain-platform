package main

import (
	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/server"
	"github.com/BakuPukul/blockchain-platform/pkg/db"
	"github.com/BakuPukul/blockchain-platform/pkg/logger"
)

func main() {
	l := logger.NewLogger()

	c, err := config.LoadConfig("config/config")
	if err != nil {
		l.Fatal(err)
	}

	db, err := db.NewRedis(c)
	if err != nil {
		l.Fatal(err)
	}

	s := server.NewServer(c, l, db)
	err = s.Run()
	if err != nil {
		l.Fatal(err)
	}
}
