package util

import (
	"time"
)

func Retry(fn func() (bool, error), attempts int, delay time.Duration) (bool, error) {
	for i := 0; i < attempts; i++ {
		if p, err := fn(); err != nil {
			return false, err
		} else if p {
			return p, nil
		}
		time.Sleep(delay)
	}
	return false, nil
}
