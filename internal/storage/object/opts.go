package object

import "time"

// for timeless objects
var interstellar, _ = time.Parse(time.DateOnly, "2067-01-01")

// WithTimeout set timeout for object
func WithTimeout(timeout time.Duration) Opt {
	return func(o object) object {
		o.expires = time.Now().Add(timeout)
		return o
	}
}

// WithDeadline set deadline for object
func WithDeadline(deadline time.Time) Opt {
	return func(o object) object {
		o.expires = deadline
		return o
	}
}

// WithoutTimeout set expire date to interstellar event time
func WithoutTimeout() Opt {
	return func(o object) object {
		o.expires = interstellar
		return o
	}
}
