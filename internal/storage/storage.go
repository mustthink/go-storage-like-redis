package storage

import (
	"time"

	"github.com/mustthink/go-storage-like-redis/config"
	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

const (
	defaultCollection = "default"
)

type Storage interface {
	NewCollection(name string) (err error)
	GetCollection(name string) (collection Collection, err error)
	DeleteCollection(name string) (err error)

	refreshing()
	defaultTimeout() time.Duration
}

func GetObject(s Storage, collectionName, objectKey string) (object.Object, error) {
	collection, err := s.GetCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.Get(objectKey)
}

func SetObject(s Storage, collectionName, objectKey string, objSettings object.RequestSettings) error {
	collection, err := s.GetCollection(collectionName)
	if err != nil {
		return err
	}

	obj := objSettings.New(s.defaultTimeout())
	collection.Set(objectKey, obj)
	return nil
}

func DeleteObject(s Storage, collectionName, objectKey string) error {
	collection, err := s.GetCollection(collectionName)
	if err != nil {
		return err
	}

	return collection.Delete(objectKey)
}

// storage is simple implementation of Storage
type storage struct {
	collections map[string]Collection
	config      config.StorageConfig
}

func New(config config.StorageConfig) Storage {
	// create default collection
	collections := make(map[string]Collection)
	collections[defaultCollection] = NewCollection()

	storage := &storage{
		collections: collections,
		config:      config,
	}

	// start refreshing storage collections
	go storage.refreshing()
	return storage
}

func (s *storage) NewCollection(name string) error {
	if len(s.collections) == s.config.MaxCollectionsCount {
		return errors.ErrMaxCollectionsCount
	}

	if _, ok := s.collections[name]; ok {
		return errors.ErrCollectionAlreadyExist(name)
	}

	s.collections[name] = NewCollection()
	return nil
}

func (s *storage) GetCollection(name string) (Collection, error) {
	if name == "" {
		return s.collections[defaultCollection], nil
	}

	if collection, ok := s.collections[name]; ok {
		return collection, nil
	}

	return nil, errors.ErrNoCollection(name)
}

func (s *storage) DeleteCollection(name string) error {
	if name == defaultCollection {
		return errors.ErrDeleteDefaultCollection
	}

	if _, ok := s.collections[name]; !ok {
		return errors.ErrNoCollection(name)
	}

	delete(s.collections, name)
	return nil
}

// parallel refreshing collections
func (s *storage) refreshing() {
	ticker := time.NewTicker(s.config.RefreshTime * time.Second)
	for ; ; <-ticker.C {
		for _, collection := range s.collections {
			go collection.Refresh()
		}
	}
}

func (s *storage) defaultTimeout() time.Duration {
	return s.config.DefaultTTL * time.Second
}
