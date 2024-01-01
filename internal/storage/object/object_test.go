package object

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testDefaultTimeout = time.Second * 2

func TestRequestSettings_NewKey(t *testing.T) {
	tests := []struct {
		name      string
		request   RequestSettings
		wantError bool
	}{
		{
			name: "shouldn't fail",
			request: RequestSettings{
				Data: []byte("1"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.request.NewKey()
			assert.Equal(t, test.wantError, err != nil)
		})
	}
}

func TestObject_IsExpired(t *testing.T) {
	tests := []struct {
		name        string
		request     RequestSettings
		sleep       time.Duration
		wantExpired bool
	}{
		{
			name: "default timeout: not expired",
			request: RequestSettings{
				Data: []byte("1"),
			},
			sleep: 1,
		},
		{
			name: "default timeout: expired",
			request: RequestSettings{
				Data: []byte("1"),
			},
			sleep:       5,
			wantExpired: true,
		},
		{
			name: "timeout: not expired",
			request: RequestSettings{
				Data:    []byte("1"),
				Timeout: 5,
			},
			sleep: 1,
		},
		{
			name: "timeout: expired",
			request: RequestSettings{
				Data:    []byte("1"),
				Timeout: 1,
			},
			sleep:       5,
			wantExpired: true,
		},
		{
			name: "deadline: not expired",
			request: RequestSettings{
				Data:     []byte("1"),
				Deadline: time.Now().Add(time.Hour),
			},
			sleep: 1,
		},
		{
			name: "deadline: expired",
			request: RequestSettings{
				Data:     []byte("1"),
				Deadline: time.Now().Add(time.Second),
			},
			sleep:       5,
			wantExpired: true,
		},
		{
			name: "timeless",
			request: RequestSettings{
				Data:     []byte("1"),
				Timeless: true,
			},
			sleep: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			obj := test.request.New(testDefaultTimeout)
			time.Sleep(test.sleep * time.Second)
			assert.Equal(t, test.wantExpired, obj.IsExpired())
		})
	}
}
