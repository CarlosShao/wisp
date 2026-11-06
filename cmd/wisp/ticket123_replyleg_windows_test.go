//go:build windows

package main

// Ticket 123: the allowlist spelling the hosted runner meets.
//
// GitHub windows runners hand TMP out as an 8.3 short path
// (C:\Users\RUNNER~1\AppData\Local\Temp\...). A fixture that allowlists only
// the spelling t.TempDir returned can therefore miss the write target when
// the judge reads the target back in its expanded long form: the target
// leaves the allowlist, the call is judged outside it and takes the L2
// route, and the ticket-123 failure shape follows. This volume cannot
// reproduce that locally (8.3 aliases are disabled here - ShortPath equals
// the long path), so the guard is belt and braces, not a measured repair.
//
// Feeding both spellings into [fs].allowed_dirs keeps the allowlist verdict
// spelling-independent without touching any judge or gate code: the frozen
// C18/C26 semantics are untouched, and the fixture simply stops betting on
// which spelling the environment hands back.

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// t123dirSpellingForms returns dir in every spelling the judge may meet it
// in: the canonical slash form the fixtures already wrote, plus the long
// (8.3-expanded) form when the volume actually generated an alias. Deduped,
// so an alias-free volume yields exactly the one form it always did.
func t123dirSpellingForms(dir string) []string {
	forms := []string{filepath.ToSlash(dir)}
	if p16, err := windows.UTF16PtrFromString(dir); err == nil {
		if n, err := windows.GetLongPathName(p16, nil, 0); err == nil && n > 0 {
			buf := make([]uint16, n)
			if n2, err := windows.GetLongPathName(p16, &buf[0], n); err == nil && n2 > 0 {
				long := windows.UTF16ToString(buf[:n2])
				if long != "" && !strings.EqualFold(long, dir) {
					forms = append(forms, filepath.ToSlash(long))
				}
			}
		}
	}
	return forms
}
