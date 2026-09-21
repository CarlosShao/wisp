//go:build windows

package models

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 95 AC#3's other half: this class is left inherit-wide ON PURPOSE, and
// a test that measures it is what separates that from an unwired site. The
// premise (public content, per-file verification at read time) is pinned in
// no_seal_ruling_test.go.
//
// If somebody later routes archive.go's extraction through winsec, this test
// goes red. That is intended: sealing here is a product decision (it is what
// would stop a second instance running under another account, e.g. a service
// account, from reusing a GB-sized cache), and decisions like that do not get
// made by an unrelated hardening sweep.

// usersSID is BUILTIN\Users, the class "any other local account".
const usersSID = "*S-1-5-32-545"

func TestAC3ExtractionIsDeliberatelyNotSealed(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "staging")
	if err := os.Mkdir(dir, 0o777); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	icaclsRun(t, dir, "/grant:r", usersSID+`:(OI)(CI)(RX)`)
	acl := icaclsRun(t, dir)
	t.Logf("wide staging parent %s\n%s", dir, acl)
	if !namesPrincipal(acl, dir, "BUILTIN\\Users") {
		t.Fatalf("icacls seeding did not land, this test would prove nothing:\n%s", acl)
	}

	archive := filepath.Join(dir, "tiny-model.tar.bz2")
	body := mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2"))
	if err := os.WriteFile(archive, body, 0o644); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "kws-fixture")
	if err := os.Mkdir(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ExtractTarBz2(archive, dest); err != nil {
		t.Fatalf("ExtractTarBz2: %v", err)
	}
	for _, want := range tinyArchiveMembers {
		p := filepath.Join(dest, filepath.FromSlash(want.Path))
		got := icaclsRun(t, p)
		t.Logf("icacls %s\n%s", p, got)
		if !namesPrincipal(got, p, "BUILTIN\\Users") {
			t.Errorf("%s stopped being inherit-wide: something sealed this class. Re-open the "+
				"ticket 95 AC#1 ruling before landing it:\n%s", want.Path, got)
		}
		// The same file must still be readable and hash-correct, i.e. the
		// measurement above is about a file the process itself uses.
		if err := verifyFileHash(p, want.SHA256, want.SizeBytes); err != nil {
			t.Errorf("extracted %s: %v", want.Path, err)
		}
	}
}

func namesPrincipal(acl, path, name string) bool {
	for _, line := range strings.Split(strings.ReplaceAll(acl, "\r\n", "\n"), "\n") {
		line = strings.ReplaceAll(line, path, "")
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(line[:i]), name) {
			return true
		}
	}
	return false
}

func icaclsRun(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	cmd := exec.Command("icacls", args...)
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("icacls %v: %v\n%s", args, err, buf.String())
	}
	return buf.String()
}
