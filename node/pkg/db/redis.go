package db

import (
	"context"
	"os"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(config *config.Config) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_DSN"),
		Password: "",
		DB:       0,
	})

	if _, err := db.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return db, nil
}
