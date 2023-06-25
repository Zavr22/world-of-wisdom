package config

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"os"
)

// Config - Configuration for app (both client and server for convenience)
type Config struct {
	ServerHost            string `envconfig:"SERVER_HOST"`
	ServerPort            int    `envconfig:"SERVER_PORT"`
	HashcashZerosCount    int
	HashcashDuration      int64
	HashcashMaxIterations int
}

// Load - loads config from file on path
func Load(path string) (*Config, error) {
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return &config, err
	}
	err = envconfig.Process("", &config)
	return &config, err
}
