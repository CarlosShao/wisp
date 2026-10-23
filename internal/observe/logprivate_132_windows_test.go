//go:build windows

package observe

// The Windows half of ticket 132's instrument: the reading is `icacls` on the
// object the pipeline made, parsed into trustees and ACE flags. Two facts are
// read out of it, and they answer different questions:
//
//   - no trustee beyond the current user, SYSTEM and Administrators: the seal
//     happened at all;
//   - at least one ACE with no (I) flag: the seal on THIS object is its own
//     descriptor rather than a grant inherited from the directory above.
//
// The second fact is what keeps the file case from measuring air. A file born
// inside an already-sealed directory inherits a narrow ACE, so "no foreign
// trustee" is green with no file seal at all - ticket 132's c3 leg booked exactly
// that shape ("按逐枚写", not "at least one object in the tree"). An explicit ACE
// on the file itself can only come from a seal applied to the file.

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

// ace132 is one trustee line out of icacls.
type ace132 struct {
	trustee   string
	inherited bool
}

// judge132 reads the descriptor and returns whether it is private plus a detail
// string that names every trustee and every ACE kind it saw, so a red says what
// was wrong rather than only that something was.
func judge132(t *testing.T, path string) (bool, string) {
	t.Helper()
	acl := icacls132(t, path)
	t.Logf("ticket132 reading: icacls %s\n%s", path, acl)
	aces := parseACL132(acl, path)
	if len(aces) == 0 {
		return false, fmt.Sprintf("%s: icacls named no trustee at all, so this reading cannot tell sealed from unreadable:\n%s", path, acl)
	}
	account := currentAccount132(t)
	var foreign, explicit []string
	for _, ace := range aces {
		if !ace.inherited {
			explicit = append(explicit, ace.trustee)
		}
		if !allowedTrustee132(ace.trustee, account) {
			foreign = append(foreign, ace.trustee)
		}
	}
	if len(foreign) > 0 {
		return false, fmt.Sprintf("%s: these trustees still hold grants: %v", path, foreign)
	}
	if len(explicit) == 0 {
		var all []string
		for _, ace := range aces {
			all = append(all, ace.trustee+"(inherited)")
		}
		return false, fmt.Sprintf("%s: every ACE is inherited (%v), so nothing was sealed on this object - its narrowness, if any, was born from the level above", path, all)
	}
	return true, fmt.Sprintf("%s: private, explicit ACEs for %v", path, explicit)
}

func assertSealedDir132(t *testing.T, path string) { assertSealed132(t, path) }

func assertSealedFile132(t *testing.T, path string) { assertSealed132(t, path) }

func assertSealed132(t *testing.T, path string) {
	t.Helper()
	if ok, detail := judge132(t, path); !ok {
		t.Errorf("ticket 132: not sealed: %s", detail)
	} else {
		t.Logf("ticket 132: sealed: %s", detail)
	}
}

// widenDir132 grants Authenticated Users the rights the inherited DACL of an
// ordinary profile directory already carries, so the fixture reproduces "a log
// tree an older build left wide" instead of inventing a hostile shape. The SID is
// spelled numerically so the fixture does not depend on the account database's
// localized names.
func widenDir132(t *testing.T, dir string) {
	t.Helper()
	icacls132(t, dir, "/grant:r", `*S-1-5-11:(OI)(CI)(M)`)
	if ok, _ := judge132(t, dir); ok {
		t.Fatalf("fixture is wrong: %s still reads as sealed after the widening, so the repair case would prove nothing", dir)
	}
}

// dirIsSealed132 reports the current verdict without failing the test, for
// fixtures that have to prove their own precondition.
func dirIsSealed132(t *testing.T, dir string) bool {
	t.Helper()
	ok, _ := judge132(t, dir)
	return ok
}

func parseACL132(acl, path string) []ace132 {
	var out []ace132
	acl = strings.ReplaceAll(acl, "\r\n", "\n")
	for _, line := range strings.Split(acl, "\n") {
		line = strings.ReplaceAll(line, path, "")
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue
		}
		trustee := strings.TrimSpace(line[:i])
		if trustee == "" {
			continue
		}
		out = append(out, ace132{trustee: trustee, inherited: strings.Contains(line[i:], "(I)")})
	}
	return out
}

// allowedTrustee132 is the set privateACL grants: the current user, SYSTEM and
// Administrators. The account may come back qualified ("MACHINE\user") or bare,
// so both spellings and the NT AUTHORITY prefix are accepted.
func allowedTrustee132(trustee, account string) bool {
	switch {
	case strings.EqualFold(trustee, `NT AUTHORITY\SYSTEM`),
		strings.EqualFold(trustee, `NT AUTHORITY\LOCAL SERVICE`),
		strings.EqualFold(trustee, `BUILTIN\Administrators`):
		return true
	case strings.EqualFold(trustee, account):
		return true
	case strings.HasSuffix(strings.ToLower(trustee), `\`+strings.ToLower(account)):
		return true
	}
	return false
}

func icacls132(t *testing.T, args ...string) string {
	t.Helper()
	var buf bytes.Buffer
	cmd := exec.Command("icacls", args...)
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("icacls %v: %v\n%s", args, err, buf.String())
	}
	return buf.String()
}

func currentAccount132(t *testing.T) string {
	t.Helper()
	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	return u.Username
}

// TestUnsealableLogDirIsRefusedLoudly132r1 - the directory half of AC#2's failure
// direction, measured on a real filesystem state rather than through a seam. The
// log directory arrives through a junction, which is the one shape the C26
// resolver refuses to traverse: sealing "A" through a link would modify B and then
// report A as private (winsec.PrivateDirAll's doc). The reading is that InitLog
// fails AND that nothing was created behind the link - a refusal that still
// writes is the silent degradation the ticket forbids.
func TestUnsealableLogDirIsRefusedLoudly132r1(t *testing.T) {
	realDir := filepath.Join(t.TempDir(), "someone-elses-tree")
	if err := os.MkdirAll(realDir, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "logs-through-a-junction")
	if out, err := exec.Command("cmd", "/c", "mklink", "/J", link, realDir).CombinedOutput(); err != nil {
		t.Skipf("this box cannot create a junction, so the refusal route has no world to be measured in: %v\n%s", err, out)
	}

	before := entries132(t, realDir)
	p, err := InitLogWithRegistry(LogConfig{Dir: filepath.Join(link, "logs")}, NewRegistry())
	if err == nil {
		_ = p.Close()
		t.Fatalf("InitLogWithRegistry sealed %s through a junction and reported success; the bytes would land in %s while the seal names the link", link, realDir)
	}
	if after := entries132(t, realDir); after != before {
		t.Errorf("the refusal still created something behind the link: %s now holds %d entries, had %d", realDir, after, before)
	}
	t.Logf("ticket132 refusal: %v", err)
}

func entries132(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	return len(entries)
}
