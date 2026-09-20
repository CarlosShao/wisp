//go:build windows

package risk

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// GetFinalPathNameByHandle flags (winbase.h); both default to 0x0 and are
// not exported by every golang.org/x/sys release, so they are defined here.
const (
	fileNameNormalized = 0x0
	volumeNameDOS      = 0x0
)

// resolveHandle opens the path (following reparse points, dirs included via
// FILE_FLAG_BACKUP_SEMANTICS) and returns its true final path via
// GetFinalPathNameByHandle with VOLUME_NAME_DOS. This step inherently expands
// 8.3 short names (SPEC-06 §4 pipeline step "展开 8.3 短名") and resolves
// junction/symlink targets. The result keeps any \\?\ / \\?\UNC\ prefix; the
// caller strips and normalizes it.
func resolveHandle(p string) (string, bool) {
	p16, err := windows.UTF16PtrFromString(p)
	if err != nil {
		return "", false
	}
	h, err := windows.CreateFile(p16, 0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return "", false
	}
	defer windows.CloseHandle(h)

	buf := make([]uint16, 1024)
	n, err := windows.GetFinalPathNameByHandle(h, &buf[0], uint32(len(buf)),
		fileNameNormalized|volumeNameDOS)
	if err != nil {
		return "", false
	}
	if n > uint32(len(buf)) { //nolint:gosec // n is the required buffer size in uint16 units
		buf = make([]uint16, n)
		if _, err = windows.GetFinalPathNameByHandle(h, &buf[0], n,
			fileNameNormalized|volumeNameDOS); err != nil {
			return "", false
		}
	}
	return windows.UTF16ToString(buf[:n]), true
}

// reparseComponents walks every component of p from the volume root down to
// the leaf and returns those that are reparse points (junction/symlink
// mount points). A FILE_ATTRIBUTE_REPARSE_POINT on ANY traversed component
// triggers the default-deny rule of SPEC-06 §4. GetFileAttributes does not
// follow the reparse point for the attribute itself, so each component is
// inspected exactly as spelled.
func reparseComponents(p string) []string {
	vol := filepath.VolumeName(p)
	rest := strings.TrimPrefix(p, vol)
	rest = strings.TrimPrefix(rest, `\`)
	if rest == "" {
		return nil
	}
	var comps []string
	base := vol
	if strings.HasPrefix(rest, `?\`) { // \\?\X: form: volume includes the prefix
		base = vol + `\`
		rest = strings.TrimPrefix(rest, `?\`)
	}
	for _, part := range strings.Split(rest, `\`) {
		if part == "" {
			continue
		}
		if base != "" && !strings.HasSuffix(base, `\`) {
			base += `\`
		}
		base += part
		attrs, err := windows.GetFileAttributes(windows.StringToUTF16Ptr(base))
		if err != nil {
			// Component missing or inaccessible: nothing further can be a
			// traversed reparse point that exists; stop the walk.
			break
		}
		if attrs&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
			comps = append(comps, base)
		}
	}
	return comps
}
