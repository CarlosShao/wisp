//go:build windows

package agent

import (
	"bytes"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

// TestAC3SpillArtifactLandsPrivate is the production call site: the artifact is
// written by Spiller.Prepare (internal/agent/spill.go), which is the path a
// model's tool output actually takes to disk. The claim under test is not that
// winsec can seal things - that is covered in internal/winsec with full SID
// resolution - but that *this* write stopped being decorative.
func TestAC3SpillArtifactLandsPrivate(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data", "artifacts")
	sp := NewSpiller(dir, spill79Budget)
	out, err := sp.Prepare("call_acl", spill79Payload("ACLA"))
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if out.Path == "" {
		t.Fatal("Prepare spilled no artifact, so this test measured nothing")
	}
	assertNoForeignACE(t, out.Path)
	assertNoForeignACE(t, dir)

	// The retry swap goes through a temp file and a rename; the temp must be
	// private too, because for its whole life it holds the artifact's bytes.
	if _, err := sp.Prepare("call_acl", spill79Payload("ACLB")); err != nil {
		t.Fatalf("retry Prepare: %v", err)
	}
	assertNoForeignACE(t, out.Path)
}

// assertNoForeignACE fails if icacls names any principal beyond the current
// user, SYSTEM and Administrators. Names rather than SIDs here because this is
// a smoke check on a call site; internal/winsec's acl_windows_test.go resolves
// every principal to a SID and is the strict form of the same criterion.
func assertNoForeignACE(t *testing.T, path string) {
	t.Helper()
	out := icaclsOut(t, path)
	account := currentAccountName(t)
	var principals, foreign []string
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n") {
		line = strings.ReplaceAll(line, path, "")
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue
		}
		principal := strings.TrimSpace(line[:i])
		principals = append(principals, principal)
		switch {
		case strings.EqualFold(principal, "NT AUTHORITY\\SYSTEM"),
			strings.EqualFold(principal, `BUILTIN\Administrators`):
		case strings.EqualFold(principal, account),
			strings.HasSuffix(strings.ToLower(principal), "\\"+strings.ToLower(account)):
		default:
			foreign = append(foreign, principal)
		}
	}
	t.Logf("icacls %s -> %v", filepath.Base(path), principals)
	if len(principals) == 0 {
		t.Fatalf("no ACE parsed out of %q - the check is not checking anything", out)
	}
	if len(foreign) > 0 {
		t.Errorf("%s is not private, these principals hold grants: %v\n%s",
			filepath.Base(path), foreign, out)
	}
}

func icaclsOut(t *testing.T, path string) string {
	t.Helper()
	var buf bytes.Buffer
	cmd := exec.Command("icacls", path)
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("icacls %s: %v\n%s", path, err, buf.String())
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
