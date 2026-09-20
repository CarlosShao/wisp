//go:build !windows

package tools

import (
	"errors"
	"path/filepath"
)

// The non-Windows half of the two platform seams. Wisp ships Windows-only
// (SPEC-01 §2), but the fs family must still COMPILE here so the risk matrix
// and the atomic writer stay testable, which is why this file answers rather
// than being absent.
//
// What it must never do is pretend: with no shell recycle bin on this platform
// fs.trash refuses and deletes nothing. An unlink dressed up as a trash would
// turn D34's "trash L1 because the bin can give it back" into a lie.

// trashDetail is what the platform backend reports about one successful call.
type trashDetail struct {
	API    string
	Record string
	Bin    string
}

// recycleBinSupported is the honest platform answer.
func recycleBinSupported() bool { return false }

// shellTrash refuses. It returns an error before touching the disk.
func shellTrash(canonical string) (trashDetail, error) {
	return trashDetail{}, errors.New(
		"本平台无 Shell 回收站 API（Wisp 的 trash 只走回收站，绝不做删除），已拒绝执行：" + canonical)
}

// volumeID is the volume identity used for D34's cross-volume move rule.
func volumeID(canonical string) string { return filepath.VolumeName(canonical) }

// volumeRootOf is the mount root a per-volume query is answered against.
func volumeRootOf(canonical string) string { return filepath.VolumeName(canonical) }

// sameVolume compares two volume identities. On a platform with one namespace,
// absolute paths with no volume name compare as one volume, which is what
// os.Rename's own atomicity guarantees are written against.
func sameVolume(a, b string) bool { return volumeID(a) == volumeID(b) }
