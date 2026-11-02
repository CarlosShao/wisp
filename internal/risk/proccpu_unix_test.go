//go:build unix

package risk

import (
	"syscall"
	"time"
)

// processCPUTime reports the user+kernel CPU this process has actually burned
// (getrusage RUSAGE_SELF), the load-insensitive instrument the C26 per-call
// budget reads against. See the Windows sibling for why CPU time and not the
// wall clock is what makes the budget reproducible under parallel packages.
func processCPUTime() (time.Duration, error) {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		return 0, err
	}
	user := time.Duration(ru.Utime.Sec)*time.Second + time.Duration(ru.Utime.Usec)*time.Microsecond
	kernel := time.Duration(ru.Stime.Sec)*time.Second + time.Duration(ru.Stime.Usec)*time.Microsecond
	return user + kernel, nil
}
