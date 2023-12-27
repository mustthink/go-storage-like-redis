package storage

import (
	"sync"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

type (
	Collection interface {
		Get(key string) (object object.Object, err error)
		Set(key string, object object.Object)
		Delete(key string) error
		Refresh()
	}

	// collection is simple implementation of Collection
	collection struct {
		objects map[string]object.Object
		mu      *sync.RWMutex
	}
)

func NewCollection() Collection {
	collection := collection{
		objects: make(map[string]object.Object),
		mu:      &sync.RWMutex{},
	}
	return collection
}

func (c collection) Get(key string) (object.Object, error) {
	c.mu.RLock()
	obj, ok := c.objects[key]
	c.mu.RUnlock()
	if !ok {
		return nil, errors.ErrNoObject(key)
	}

	if obj.IsExpired() {
		c.Delete(key)
		return nil, errors.ErrNoObject(key)
	}

	return obj, nil
}

func (c collection) Set(key string, object object.Object) {
	c.mu.Lock()
	c.objects[key] = object
	c.mu.Unlock()
}

func (c collection) Delete(key string) error {
	_, ok := c.objects[key]
	if !ok {
		return errors.ErrNoObject(key)
	}

	c.mu.Lock()
	delete(c.objects, key)
	c.mu.Unlock()
	return nil
}

func (c collection) Refresh() {
	for key, obj := range c.objects {
		if obj.IsExpired() {
			c.Delete(key)
		}
	}
}
