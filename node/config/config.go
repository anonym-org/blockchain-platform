package config

import (
	"os"
	"strings"
)

type Config struct {
	*Server
}

type Server struct {
	Port     string
	GrpcPort string
	Nodes    []string
}

func LoadConfig(configPath string) (*Config, error) {
	c := Config{
		Server: &Server{},
	}

	// env config
	c.Port = os.Getenv("PORT")
	c.GrpcPort = os.Getenv("GRPC_PORT")
	c.Nodes = strings.Split(os.Getenv("NODES"), ",")

	return &c, nil
}
