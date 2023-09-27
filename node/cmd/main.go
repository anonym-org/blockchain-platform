package main

import (
	"github.com/anonym-org/blockchain-platform/config"
	"github.com/anonym-org/blockchain-platform/internal/server"
	"github.com/anonym-org/blockchain-platform/pkg/db"
	"github.com/anonym-org/blockchain-platform/pkg/logger"
)

func main() {
	log := logger.NewLogger()

	conf, err := config.LoadConfig("config/config")
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.NewRedis(conf)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(conf, log, db)
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
