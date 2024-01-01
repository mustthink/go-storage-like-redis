package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

func TestCollection_Get(t *testing.T) {
	collection := NewCollection()
	tests := []struct {
		name       string
		toAdd      object.Object
		getKey     string
		wantObject object.Object
		wantError  error
	}{
		{
			name:       "shouldn't fail",
			toAdd:      testObject,
			getKey:     testKey,
			wantObject: testObject,
		},
		{
			name:      "no object",
			toAdd:     testObject,
			getKey:    "unknown",
			wantError: errors.ErrNoObject("unknown"),
		},
		{
			name:      "object expired",
			toAdd:     object.New([]byte("1"), object.WithTimeout(-time.Second)),
			getKey:    testKey,
			wantError: errors.ErrNoObject(testKey),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			collection.Set(testKey, test.toAdd)

			obj, err := collection.Get(test.getKey)
			assert.Equal(t, err, test.wantError)
			assert.Equal(t, obj, test.wantObject)
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	collection := NewCollection()
	tests := []struct {
		name      string
		toAdd     object.Object
		getKey    string
		wantError error
	}{
		{
			name:   "shouldn't fail",
			toAdd:  testObject,
			getKey: testKey,
		},
		{
			name:      "no object",
			toAdd:     testObject,
			getKey:    "unknown",
			wantError: errors.ErrNoObject("unknown"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			collection.Set(testKey, test.toAdd)

			err := collection.Delete(test.getKey)
			assert.Equal(t, err, test.wantError)
		})
	}
}
