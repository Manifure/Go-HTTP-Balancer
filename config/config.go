package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ListenPort string   `json:"listen_port"`
	Backends   []string `json:"backends"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}
