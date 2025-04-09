package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	NodeID string   `json:"nodeID"`
	Port   string   `json:"port"`
	Peers  []string `json:"peers"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not decode config: %v", err)
	}

	return &cfg, nil
}
