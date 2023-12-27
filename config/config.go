package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const DefaultConfig = "config/default.json"

type (
	BaseAuthConfig struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	StorageConfig struct {
		DefaultTTL          time.Duration `json:"ttl_in_seconds"`
		MaxCollectionsCount int           `json:"max_collections_count"`
		RefreshTime         time.Duration `json:"refresh_time_in_seconds"`
	}

	ServerConfig struct {
		Host string         `json:"host"`
		Port string         `json:"port"`
		Auth BaseAuthConfig `json:"auth"`

		ReadTimeout  time.Duration `json:"read_timeout_in_ms"`
		WriteTimeout time.Duration `json:"write_timeout_in_ms"`
	}

	Config struct {
		StorageConfig StorageConfig `json:"storage"`
		ServerConfig  ServerConfig  `json:"server"`
	}
)

func New(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read configuration file w err: %s", err.Error())
	}

	var configuration Config
	if err := json.Unmarshal(data, &configuration); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal configuration w err: %s", err.Error())
	}

	if err := configuration.validation(); err != nil {
		return nil, fmt.Errorf("couldn't validate config w err: %s", err.Error())
	}

	return &configuration, nil
}

func (s ServerConfig) URL() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
