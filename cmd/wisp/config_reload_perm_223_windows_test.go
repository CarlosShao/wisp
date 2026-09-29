//go:build windows

package main

// Ticket 223 AC#4's third sentence: 权限不够. It gets its own file because the
// only way to make a real read-permission failure on this platform is a real
// ACL, and icacls is not available anywhere else. This is the shape
// internal/config/private_acl_windows_test.go already uses for the same reason.
//
// WHY AN ACL AND NOT A MOCK. The classifier in config_reload.go decides between
// four sentences off the *error the filesystem actually returned* - so the case
// plants the real condition and reads the real line back. os.Chmod cannot do it
// on Windows (a cleared bit becomes the read-only attribute, which still reads
// fine), so this denies read to Everyone through icacls and then PROVES the
// denial landed before asking the product about it. A seeding that did not land
// is a t.Fatal, not a skip: a case that cannot plant its condition cannot pass
// either.

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// everyoneSID is the well-known group every token carries, so one deny ACE on it
// is enough to take read away from this process.
const everyoneSID = "*S-1-1-0"

func runIcacls223(t *testing.T, args ...string) string {
	t.Helper()
	out, err := icaclsRaw223(args...)
	if err != nil {
		t.Fatalf("icacls %v: %v\n%s", args, err, out)
	}
	return out
}

// runIcaclsClear223 removes this case's deny ACE if one is there. Failure is
// logged, never fatal: on the first attempt there is nothing to remove, and
// icacls says so in its own words.
func (r *reloadRun223) runIcaclsClear223(t *testing.T, path string) {
	t.Helper()
	out, err := icaclsRaw223(path, "/remove:d", everyoneSID)
	if err != nil {
		t.Logf("icacls /remove:d (nothing to remove is fine): %v\n%s", err, out)
	}
}

func icaclsRaw223(args ...string) (string, error) {
	cmd := exec.Command("icacls", args...)
	var out, errb strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	return out.String() + errb.String(), err
}

func TestTicket223PermissionDeniedSitsInItsOwnSentence(t *testing.T) {
	r := newReloadRun223(t, 60*time.Second)
	r.live(t, func() {
		path := r.cfgPath()
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")

		// WHAT THIS PLATFORM ACTUALLY DOES, measured on the first two drafts of
		// this case: a (R) deny on the file does not stop os.Stat (Go reads the
		// directory entry instead), and it DOES stop os.WriteFile ("Access is
		// denied" on the O_WRONLY open). So "the file's fingerprint moved and
		// this process can no longer read it" has to be planted as
		// write-then-deny, and the 1s tick can in principle slip between the two
		// statements - if it does, the case says so and plants again instead of
		// timing out on a sentence that can never arrive.
		var trail string
		t.Cleanup(func() { r.runIcaclsClear223(t, path) })
		for attempt := 1; attempt <= 3; attempt++ {
			// Re-enable write (no-op on the first attempt), move the fingerprint,
			// then deny again. The content is deliberately the same shape every
			// time; only its length changes, which is enough for mtime+size.
			r.runIcaclsClear223(t, path)
			body := "schema_version = 2\n\n[llm]\ntext_chain = [\"acme/m1\"]\n\n" +
				"# padding to move the fingerprint: " + strings.Repeat("x", attempt*5) + "\n"
			if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
				t.Fatalf("attempt %d cannot move the fingerprint: %v", attempt, err)
			}
			runIcacls223(t, path, "/deny", everyoneSID+":(R)")

			// Positive control on the plant, before the product is asked: stat
			// answers, read does not. That pair is what separates 权限不够 from
			// 缺失, so it is asserted rather than assumed.
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("attempt %d: stat failed too, which is the missing-file case: %v", attempt, err)
			}
			if _, err := os.ReadFile(path); err == nil {
				t.Fatalf("attempt %d: icacls denied nothing, so this case cannot show the permission sentence", attempt)
			}

			got, adopted := r.awaitFailureOrAdoption(t)
			if got != "" {
				trail = got
				break
			}
			if adopted {
				t.Logf("attempt %d: the tick read the file inside the write/deny window and adopted it; planting again", attempt)
				continue
			}
			t.Fatalf("neither the permission sentence nor an adoption arrived; stderr:\n%s", r.h.err.String())
		}
		if trail == "" {
			t.Fatal("three plants all landed in the tick's read window; this case cannot be decided on this machine")
		}
		for _, other := range []string{
			"cause=missing", "cause=syntax", "cause=unknown-key",
			"cause=invalid", "cause=unclassified", "cause=migration",
		} {
			if strings.Contains(trail, other) {
				t.Errorf("the permission sentence carries %q:\n%s", other, trail)
			}
		}
		// The sentence names the cause in the operator's own terms, and it does
		// not claim the file is gone or broken.
		if !strings.Contains(trail, "没有读它的权限") {
			t.Errorf("the permission sentence does not say what is wrong in words: %q", trail)
		}
	})
}

// awaitFailureOrAdoption gives the tick one cycle to either report the read
// failure (returns the trail) or to have already re-read the file (returns
// adopted=true, so the caller plants again).
func (r *reloadRun223) awaitFailureOrAdoption(t *testing.T) (trail string, adopted bool) {
	t.Helper()
	deadline := time.NewTimer(6 * time.Second)
	defer deadline.Stop()
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()
	for {
		got := r.h.err.String()
		if strings.Contains(got, "config: HOT-RELOAD state=not-applied cause=permission") {
			return got, false
		}
		if strings.Contains(got, "config: HOT-RELOAD state=applied") ||
			strings.Contains(got, "state=not-applied cause=invalid") {
			return "", true
		}
		select {
		case <-deadline.C:
			return "", false
		case <-tick.C:
		}
	}
}
