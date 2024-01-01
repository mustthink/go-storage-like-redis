package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
)

func TestConfig_validation(t *testing.T) {
	tests := []struct {
		name       string
		haveConfig Config
		wantError  error
	}{
		{
			name: "valid config",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					DefaultTTL:          1,
					MaxCollectionsCount: 1,
					RefreshTime:         1,
				},
				ServerConfig: ServerConfig{
					Host:         "host",
					Port:         "port",
					ReadTimeout:  1,
					WriteTimeout: 1,
				},
			},
			wantError: nil,
		},
		{
			name: "StorageConfig: default ttl is empty",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					MaxCollectionsCount: 1,
					RefreshTime:         1,
				},
				ServerConfig: ServerConfig{
					Host:         "host",
					Port:         "port",
					ReadTimeout:  1,
					WriteTimeout: 1,
				},
			},
			wantError: errors.ErrEmptyField("default_ttl"),
		},
		{
			name: "StorageConfig: max collections count is empty",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					DefaultTTL:  1,
					RefreshTime: 1,
				},
				ServerConfig: ServerConfig{
					Host:         "host",
					Port:         "port",
					ReadTimeout:  1,
					WriteTimeout: 1,
				},
			},
			wantError: errors.ErrEmptyField("max_collections_count"),
		},
		{
			name: "StorageConfig: refresh time is empty",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					DefaultTTL:          1,
					MaxCollectionsCount: 1,
				},
				ServerConfig: ServerConfig{
					Host:         "host",
					Port:         "port",
					ReadTimeout:  1,
					WriteTimeout: 1,
				},
			},
			wantError: errors.ErrEmptyField("refresh_time"),
		},
		{
			name: "ServerConfig: host is empty",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					DefaultTTL:          1,
					MaxCollectionsCount: 1,
					RefreshTime:         1,
				},
				ServerConfig: ServerConfig{
					Port:         "port",
					ReadTimeout:  1,
					WriteTimeout: 1,
				},
			},
			wantError: errors.ErrEmptyField("host"),
		},
		{
			name: "ServerConfig: port is empty",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					DefaultTTL:          1,
					MaxCollectionsCount: 1,
					RefreshTime:         1,
				},
				ServerConfig: ServerConfig{
					Host:         "host",
					ReadTimeout:  1,
					WriteTimeout: 1,
				},
			},
			wantError: errors.ErrEmptyField("port"),
		},
		{
			name: "ServerConfig: ReadTimeout is empty",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					DefaultTTL:          1,
					MaxCollectionsCount: 1,
					RefreshTime:         1,
				},
				ServerConfig: ServerConfig{
					Host:         "host",
					Port:         "port",
					WriteTimeout: 1,
				},
			},
			wantError: errors.ErrEmptyField("read_timeout"),
		},
		{
			name: "ServerConfig: WriteTimeout is empty",
			haveConfig: Config{
				StorageConfig: StorageConfig{
					DefaultTTL:          1,
					MaxCollectionsCount: 1,
					RefreshTime:         1,
				},
				ServerConfig: ServerConfig{
					Host:        "host",
					Port:        "port",
					ReadTimeout: 1,
				},
			},
			wantError: errors.ErrEmptyField("write_timeout"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.haveConfig.validation()
			assert.Equal(t, err, test.wantError)
		})
	}
}
