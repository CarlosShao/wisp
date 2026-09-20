//go:build !windows

package secret

import (
	"os"
	"time"
)

// createdTime falls back to the modification time off Windows: the creation
// timestamp is a NTFS/Win32 find-data field (SPEC-02 §6). This shim exists so
// the linux test-core job compiles the package (ticket 08), not to provide
// create times there.
func createdTime(fi os.FileInfo) time.Time { return fi.ModTime() }
