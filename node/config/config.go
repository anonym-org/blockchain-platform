package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	*Server
	*DB
}

type Server struct {
	Port string
}

type DB struct {
	DSN string
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(configPath)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	c := Config{
		Server: &Server{},
		DB:     &DB{},
	}
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	// env config
	c.Port = os.Getenv("PORT")
	c.DSN = os.Getenv("DB_DSN")

	return &c, nil
}
