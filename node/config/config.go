package config

import (
	"os"
	"strings"
)

type Config struct {
	*Server
	*DB
}

type Server struct {
	Port     string
	GrpcPort string
	Nodes    []string
}

type DB struct {
	DSN string
}

func LoadConfig(configPath string) (*Config, error) {
	c := Config{
		Server: &Server{},
		DB:     &DB{},
	}

	// env config
	c.Port = os.Getenv("PORT")
	c.GrpcPort = os.Getenv("GRPC_PORT")
	c.DSN = os.Getenv("DB_DSN")
	c.Nodes = strings.Split(os.Getenv("NODES"), ",")

	return &c, nil
}
