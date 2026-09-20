//go:build windows

package secret

import (
	"os"
	"syscall"
	"time"
)

// createdTime reads the NTFS creation time Go already fetched in the find
// data (ticket 63: `wisp secret list` shows blob creation times, SPEC-02 §6
// layout). A missing/zero FILETIME falls back to the modification time so the
// listing never shows the epoch.
func createdTime(fi os.FileInfo) time.Time {
	sys, ok := fi.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return fi.ModTime()
	}
	ns := sys.CreationTime.Nanoseconds()
	if ns == 0 {
		return fi.ModTime()
	}
	return time.Unix(0, ns).UTC()
}
