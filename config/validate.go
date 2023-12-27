package config

import "github.com/mustthink/go-storage-like-redis/internal/errors"

type validateItem interface {
	Validate() error
}

func (s StorageConfig) Validate() error {
	switch {
	case s.DefaultTTL == 0:
		return errors.ErrEmptyField("default_ttl")
	case s.RefreshTime == 0:
		return errors.ErrEmptyField("refresh_time")
	case s.MaxCollectionsCount <= 0:
		return errors.ErrEmptyField("max_collections_count")
	default:
		return nil
	}
}

func (s ServerConfig) Validate() error {
	switch {
	case s.Host == "":
		return errors.ErrEmptyField("host")
	case s.Port == "":
		return errors.ErrEmptyField("port")
	case s.ReadTimeout == 0:
		return errors.ErrEmptyField("read_timeout")
	case s.WriteTimeout == 0:
		return errors.ErrEmptyField("write_timeout")
	default:
		return nil
	}
}

func (c Config) validation() error {
	var configs = []validateItem{c.ServerConfig, c.StorageConfig}
	for _, config := range configs {
		if err := config.Validate(); err != nil {
			return err
		}
	}
	return nil
}
