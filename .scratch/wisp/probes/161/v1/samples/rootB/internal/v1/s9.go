package v1

import "time"

func deltaB(now time.Time) time.Duration {
	return now.Sub(time.Now())
}
