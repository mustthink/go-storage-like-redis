package object

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type (
	Object interface {
		Binary() []byte
		IsExpired() bool
	}

	// simple implementation of object with expiration logic
	object struct {
		data    []byte
		expires time.Time
	}

	Opt func(object) object

	// RequestSettings is settings for add new object
	// priority of opts timeless > deadline > timeout
	RequestSettings struct {
		Data []byte `json:"data"`

		// timeout in seconds
		Timeout time.Duration `json:"timeout"`

		// expire date
		Deadline time.Time `json:"deadline"`

		// without expiration
		Timeless bool `json:"timeless"`
	}
)

func New(data []byte, opts ...Opt) Object {
	object := object{
		data: data,
	}

	for _, opt := range opts {
		object = opt(object)
	}

	return object
}

func (s RequestSettings) New(defaultTimeout time.Duration) Object {
	switch {
	case s.Timeout != 0:
		return New(s.Data, WithTimeout(s.Timeout*time.Second))
	case !s.Deadline.IsZero():
		return New(s.Data, WithDeadline(s.Deadline))
	case s.Timeless:
		return New(s.Data, WithoutTimeout())
	}
	return New(s.Data, WithTimeout(defaultTimeout))
}

func (s RequestSettings) NewKey() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("couldn't marshal request settings w err: %s", err.Error())
	}

	hash := sha256.Sum256(data)
	return string(hash[:]), nil
}

func (o object) Binary() []byte {
	return o.data
}

func (o object) IsExpired() bool {
	now := time.Now()
	return o.expires.Before(now)
}
