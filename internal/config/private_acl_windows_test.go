//go:build windows

package config

import (
	"bytes"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 95 AC#2. The claim under test is not that winsec can seal something -
// internal/winsec measures that with full SID resolution - but that *these two*
// config writes stopped being decorative: the migration backup and the
// pre-rename temp behind config.toml. Both hold the same bytes a user would
// call sensitive (endpoints, model paths, and on a machine that has not been
// through D33 yet, plaintext api_key values).
//
// Every case below seeds its parent directory with a real grant for
// BUILTIN\Users, which is the class "some other local account": inherited
// read+execute is exactly what a mode argument of 0o600 was supposed to prevent
// and never did on NTFS. The parent is always under t.TempDir(), never a real
// data root.

// usersSID is BUILTIN\Users.
const usersSID = "*S-1-5-32-545"

// wideDir returns a fresh data directory whose children inherit read+execute
// for every local account, and proves the seeding landed before returning.
func wideDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "data")
	if err := os.Mkdir(dir, 0o777); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	icacls(t, dir, "/grant:r", usersSID+`:(OI)(CI)(RX)`)
	acl := icacls(t, dir)
	t.Logf("seeded wide parent %s\n%s", dir, acl)
	if !hasPrincipal(acl, dir, "BUILTIN\\Users") {
		t.Fatalf("icacls seeding did not land, this test would prove nothing:\n%s", acl)
	}
	return dir
}

// TestAC2BaselineModeIsDecorative is the "before" reading AC#2 asks for: it
// writes the backup the way the code did before this ticket wired winsec
// (os.WriteFile with a 0o600 mode argument) and measures who can actually read
// it. The point of this test is that it PASSES while naming a foreign
// principal - it documents the hole, it does not close it. Reverting either
// production write to the mode-argument form turns TestAC2MigrationBackup...
// and ...SaveFile red, which is AC#4.
func TestAC2BaselineModeIsDecorative(t *testing.T) {
	dir := wideDir(t)
	backup := filepath.Join(dir, "config.toml.bak-1")
	raw := []byte("schema_version = 1\n[app]\ntheme = \"dark\"\n")
	if err := os.WriteFile(backup, raw, 0o600); err != nil {
		t.Fatalf("baseline write: %v", err)
	}
	acl := icacls(t, backup)
	t.Logf("BASELINE (pre-wiring shape) icacls %s\n%s", backup, acl)
	if !hasPrincipal(acl, backup, "BUILTIN\\Users") {
		t.Errorf("expected the 0o600 mode argument to have landed nothing, got a private file: %s", acl)
	}
	if _, err := os.ReadFile(backup); err != nil {
		t.Fatalf("read-back of the baseline file: %v", err)
	}
}

// TestAC2MigrationBackupLandsPrivate drives the real read-back path: LoadFile
// on a v1 config, which is what the process does at every startup, so
// applyMigrations writes config.toml.bak-1 and rewrites config.toml. Both must
// come back private.
func TestAC2MigrationBackupLandsPrivate(t *testing.T) {
	dir := wideDir(t)
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(v1Fixture), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := LoadFile(path, nil); err != nil {
		t.Fatalf("migrating load: %v", err)
	}
	backup := path + ".bak-1"
	if _, err := os.Stat(backup); err != nil {
		t.Fatalf("migration wrote no backup, so the seal was never exercised: %v", err)
	}
	assertPrivate(t, backup)
	assertPrivate(t, path)

	// The bytes still read back exactly as the pre-migration original, i.e. the
	// seal did not eat the backup.
	got, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("read-back: %v", err)
	}
	if string(got) != v1Fixture {
		t.Errorf("backup differs from the original: %d bytes vs %d", len(got), len(v1Fixture))
	}
}

// TestAC2SaveFileLandsPrivate covers the other leg: SaveFile writes a temp in
// the config directory and renames it over config.toml. Rename carries the
// security descriptor, so sealing the temp is what makes the *final* file
// private - the claim this pins.
func TestAC2SaveFileLandsPrivate(t *testing.T) {
	dir := wideDir(t)
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(validMinimal), 0o600); err != nil {
		t.Fatal(err)
	}
	c, _, err := LoadFile(path, nil)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	// "Before", on the same file the save is about to replace: LoadFile on an
	// up-to-date schema does not rewrite, so what is on disk here is the wide
	// os.WriteFile above. Naming it here is what makes the reading below a
	// before/after pair rather than a single favourable sample.
	before := icacls(t, path)
	t.Logf("BASELINE (pre-wiring shape) icacls %s\n%s", path, before)
	if !hasPrincipal(before, path, "BUILTIN\\Users") {
		t.Errorf("expected the 0o600 mode argument to have landed nothing, got: %s", before)
	}

	if err := SaveFile(path, c); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}
	assertPrivate(t, path)
	// Read-back through the production reader, so the assertion is about the
	// file the process actually uses, not a leftover temp.
	if _, _, err := LoadFile(path, nil); err != nil {
		t.Fatalf("re-load after save: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".wisp-config-") {
			t.Errorf("temp file survived the rename: %s", e.Name())
		}
	}
}

// principals lists the trustees icacls names for path. The first ACE shares its
// line with the echoed path and the rest are aligned under it, so the path is
// stripped before the trustee name is taken; names like "NT AUTHORITY\SYSTEM"
// contain blanks and must survive intact.
func principals(acl, path string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(acl, "\r\n", "\n"), "\n") {
		line = strings.ReplaceAll(line, path, "")
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue
		}
		if p := strings.TrimSpace(line[:i]); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// assertPrivate fails if icacls names any principal beyond the current user,
// SYSTEM and Administrators, and always logs the raw text so the evidence
// carries the principal and its rights, not just a verdict.
func assertPrivate(t *testing.T, path string) {
	t.Helper()
	acl := icacls(t, path)
	t.Logf("icacls %s\n%s", path, acl)
	account := currentAccountName(t)
	var foreign []string
	for _, principal := range principals(acl, path) {
		switch {
		case strings.EqualFold(principal, "NT AUTHORITY\\SYSTEM"),
			strings.EqualFold(principal, `BUILTIN\Administrators`):
		case strings.EqualFold(principal, account),
			strings.EqualFold(principal, `NT AUTHORITY\`+account),
			strings.HasSuffix(strings.ToLower(principal), "\\"+strings.ToLower(account)):
		default:
			foreign = append(foreign, principal)
		}
	}
	if len(foreign) > 0 {
		t.Errorf("%s is not private, these principals hold grants: %v\n%s",
			filepath.Base(path), foreign, acl)
	}
}

func hasPrincipal(acl, path, name string) bool {
	for _, p := range principals(acl, path) {
		if strings.EqualFold(p, name) {
			return true
		}
	}
	return false
}

func icacls(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	cmd := exec.Command("icacls", args...)
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("icacls %v: %v\n%s", args, err, buf.String())
	}
	return buf.String()
}

func currentAccountName(t *testing.T) string {
	t.Helper()
	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	return u.Username
}
