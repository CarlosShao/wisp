//go:build !windows

package main

import "path/filepath"

// t123dirSpellingForms is the non-windows half of the ticket-123 spelling
// guard: outside Windows there is no 8.3 alias machinery, so the canonical
// slash form is the only spelling a target can arrive in. See the windows
// half for why the fixtures feed more than one form at all.
func t123dirSpellingForms(dir string) []string {
	return []string{filepath.ToSlash(dir)}
}
