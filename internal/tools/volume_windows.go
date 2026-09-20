//go:build windows

package tools

import (
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Volume identity for the cross-volume half of D34's move row. Comparing drive
// LETTERS would be both too weak (a folder mounted from another volume shares
// its parent's letter) and too strong (two spellings of one volume), so the
// mount point is resolved to its persistent volume GUID and only that is
// compared.

var procGetVolumeNameForVolumeMountPointW = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetVolumeNameForVolumeMountPointW")

// volumeID returns the stable identity of the volume one canonical path sits
// on. It never errors, and it fails in the safe direction: an unresolvable
// mount point yields a string that differs from every resolved GUID, so the
// move is escalated to L2 (cross-volume) instead of being waved through as a
// cheap same-volume rename.
func volumeID(canonical string) string {
	vol := strings.ToLower(filepath.VolumeName(canonical))
	if len(vol) != 2 || vol[1] != ':' {
		// A UNC share, or no volume at all: the spelling is the identity.
		return vol
	}
	root := vol + `\`
	if procGetVolumeNameForVolumeMountPointW.Find() != nil {
		return root
	}
	buf := make([]uint16, 64)
	p, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return root
	}
	r, _, _ := procGetVolumeNameForVolumeMountPointW.Call(
		uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	runtime.KeepAlive(buf)
	if r == 0 {
		return root
	}
	return strings.ToLower(windows.UTF16ToString(buf))
}

// volumeRootOf gives the mount root a per-volume Recycle Bin is queried
// against (SHQueryRecycleBinW wants a path like `C:\`, not a GUID).
func volumeRootOf(canonical string) string {
	vol := filepath.VolumeName(canonical)
	if len(vol) == 2 && vol[1] == ':' {
		return vol + `\`
	}
	return vol
}

// sameVolume is the D34 predicate behind fs.move's level split.
func sameVolume(a, b string) bool {
	ia, ib := volumeID(a), volumeID(b)
	return ia != "" && ia == ib
}
