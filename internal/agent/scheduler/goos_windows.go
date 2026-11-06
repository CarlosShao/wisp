//go:build windows

package scheduler

// goosIsWindows folds path case in the lock's comparability pass on the
// Windows product (D7). Windows filesystems are case-insensitive, so two
// spellings of one path are one path here; on other GOOS builds they are
// not, and the constant says so honestly instead of pretending.
const goosIsWindows = true
