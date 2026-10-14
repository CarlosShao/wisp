// Probe for acceptance leg 174-v3 (ticket 174 AC#3(iii)): measure the [fs]
// factory default in-process instead of quoting SPEC-03. D36 rule 3 makes the
// struct tags the single source of defaults, so NewDefaults() is where
// "allowed_dirs[]=[]" lives in code; this probe reads that table directly.
//
// It compiles into package config through `go test -overlay` (see the overlay
// json next to it), so internal/config/ stays byte-clean.
package config

import "testing"

func TestV3FSDefaultsStillEmpty174(t *testing.T) {
	c := NewDefaults()
	t.Logf("NewDefaults: FS.AllowedDirs=%v (len=%d) ReparsePointExceptions=%v DeleteEnabled=%v",
		c.FS.AllowedDirs, len(c.FS.AllowedDirs), c.FS.ReparsePointExceptions, c.FS.DeleteEnabled)
	if len(c.FS.AllowedDirs) != 0 {
		t.Fatalf("the factory default allowlist must stay empty (SPEC-03:35 allowed_dirs[]=[]), got %v",
			c.FS.AllowedDirs)
	}
	if len(c.FS.ReparsePointExceptions) != 0 {
		t.Fatalf("reparse_point_exceptions must stay empty by default, got %v", c.FS.ReparsePointExceptions)
	}
	if c.FS.DeleteEnabled {
		t.Fatalf("delete_enabled must stay false by default")
	}
}
