package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testConfig = "test.json"

func createTestConfig(config Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("couldn't marshal test config w err: %s", err.Error())
	}

	file, err := os.Create(testConfig)
	if err != nil {
		return fmt.Errorf("couldn't create test config file w err: %s", err.Error())
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("couldn't write data to test config file w err: %s", err.Error())
	}
	return nil
}

func deleteTestConfig() error {
	return os.Remove(testConfig)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		testConfig Config
		wantError  bool
	}{
		{
			name: "without any error",
			testConfig: Config{
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
		},
		{
			name:       "validation fail",
			testConfig: Config{},
			wantError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			err = createTestConfig(test.testConfig)
			require.Nil(t, err, "couldn't create test config")

			defer func() {
				err = deleteTestConfig()
				require.Nil(t, err, "couldn't delete test config")
			}()

			gotConfig, err := New(testConfig)
			assert.Equal(t, err != nil, test.wantError)

			if !test.wantError {
				assert.Equal(t, *gotConfig, test.testConfig)
			}
		})
	}
}
