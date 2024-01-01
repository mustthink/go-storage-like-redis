package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mustthink/go-storage-like-redis/config"
	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

var (
	testConfig = config.StorageConfig{
		DefaultTTL:          1,
		MaxCollectionsCount: 1,
		RefreshTime:         1000,
	}
	testCollection = ""

	testKey = "1"

	testRequestSettings = object.RequestSettings{
		Data:     []byte("1"),
		Timeless: true,
	}

	testObject = object.New([]byte("1"), object.WithoutTimeout())
)

func TestGetObject(t *testing.T) {
	tests := []struct {
		name       string
		storage    Storage
		collection string
		setObject  object.RequestSettings
		wantObject object.Object
		wantError  error
	}{
		{
			name:       "valid storage",
			storage:    New(testConfig),
			collection: testCollection,
			setObject:  testRequestSettings,
			wantObject: testObject,
		},
		{
			name:       "unknown collection",
			storage:    New(testConfig),
			collection: "unknown",
			wantError:  errors.ErrNoCollection("unknown"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SetObject(test.storage, testCollection, testKey, test.setObject)
			require.Nil(t, err)

			getObject, err := GetObject(test.storage, test.collection, testKey)
			assert.Equal(t, err, test.wantError)
			assert.Equal(t, getObject, test.wantObject)
		})
	}
}

func TestSetObject(t *testing.T) {
	tests := []struct {
		name       string
		storage    Storage
		collection string
		wantError  error
	}{
		{
			name:       "valid storage",
			storage:    New(testConfig),
			collection: testCollection,
		},
		{
			name:       "unknown collection",
			storage:    New(testConfig),
			collection: "unknown",
			wantError:  errors.ErrNoCollection("unknown"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SetObject(test.storage, test.collection, testKey, testRequestSettings)
			require.Equal(t, err, test.wantError)
		})
	}
}
