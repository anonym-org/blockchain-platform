package config

import "github.com/spf13/viper"

type Config struct {
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(configPath)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var c *Config
	err := v.Unmarshal(&c)
	return c, err
}
