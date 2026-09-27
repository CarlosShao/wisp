package v1

import "time"

func left(deadline time.Time) time.Duration {
	return time.Now().Sub(deadline)
}
