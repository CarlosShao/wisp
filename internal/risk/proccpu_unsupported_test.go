//go:build !windows && !unix

package risk

import (
	"errors"
	"time"
)

// processCPUTime has no reading on this platform. The budget test names the
// missing instrument and skips instead of silently passing, so a port to a new
// GOOS cannot turn the C26 gate into an empty assertion.
func processCPUTime() (time.Duration, error) {
	return 0, errors.New("process CPU time is not implemented for this GOOS")
}
