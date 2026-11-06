//go:build !windows

package scheduler

// goosIsWindows folds path case in the lock's comparability pass. Windows
// filesystems are case-insensitive (one spelling per path); other GOOS
// builds are case-sensitive, so no folding happens there.
const goosIsWindows = false
