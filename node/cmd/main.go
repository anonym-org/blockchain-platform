package main

import (
	"github.com/anonym-org/blockchain-platform/config"
	"github.com/anonym-org/blockchain-platform/internal/server"
	"github.com/anonym-org/blockchain-platform/pkg/logger"
	"github.com/dgraph-io/badger/v4"
)

func main() {
	log := logger.NewLogger()

	conf, err := config.LoadConfig("config/config")
	if err != nil {
		log.Fatal(err)
	}

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(conf, log, db)
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
