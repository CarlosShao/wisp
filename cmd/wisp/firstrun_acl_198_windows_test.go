//go:build windows

package main

// Ticket 198 AC#3 as a live measurement, not a quote: the file first-run
// created is checked with the only ACL reading this repository accepts as
// evidence (icacls, ticket 89's ruling - FileInfo.Mode() reports the fiction
// the mode argument was supposed to control). This mirrors
// internal/config/private_acl_windows_test.go:103-128 (TestAC2SaveFileLands
// Private, tickets 89/131/132's family) one layer up: that pin proves
// SaveFile's output lands private in the config package; this one proves the
// CLI first-run actually goes through SaveFile instead of a raw write -
// 198-a1 §4-N10 names exactly that fork ("若绕过它用 os.WriteFile 直写就没人钉了").

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTicket198AC3CreatedFileLandsPrivate(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, configFileName)
	code, _, stderrText := run198(t, dir)
	if code != 2 {
		t.Fatalf("exit %d, want 2, stderr:\n%s", code, stderrText)
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("first run did not create the file whose ACL this case checks: %v", err)
	}
	acl := icacls117(t, cfgPath)
	if namesEveryone117(acl) {
		t.Errorf("the first-run config.toml carries a world-readable grant (票 132 一形要求私有):\n%s", acl)
	}
	t.Logf("AC#3 icacls(cfgPath) verbatim:\n%s", acl)
}
