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
	cmd := exec.Command("icacls", args...)
	var out, errb strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		t.Fatalf("icacls %v: %v\n%s%s", args, err, out.String(), errb.String())
	}
	return out.String()
}

func TestTicket223PermissionDeniedSitsInItsOwnSentence(t *testing.T) {
	r := newReloadRun223(t, 40*time.Second)
	r.live(t, func() {
		path := r.cfgPath()
		// The whole trail so far, so the case can prove the tick is what failed
		// and not the assembly.
		r.awaitAudit(t, "config: HOT-RELOAD state=armed")

		// Deny read AND attribute-read. Attribute-read is the point: what the
		// poll does first is stat the file (manager.go:118), and a fingerprint
		// that has not moved is a file it never opens - so a plain (R) deny
		// leaves this sentence unreachable without also rewriting the file,
		// which a process that cannot read it must not do. Measured: the first
		// draft of this case denied (R) alone and os.Chtimes then failed with
		// "Access is denied", so the plant could not move the fingerprint either.
		runIcacls223(t, path, "/deny", everyoneSID+":(RX)")
		t.Cleanup(func() {
			// Best-effort: TempDir removal must not trip over the ACE this case
			// installed. A failure here is logged, not fatal - the denial itself
			// was already proven below.
			var out, errb strings.Builder
			cmd := exec.Command("icacls", path, "/remove:d", everyoneSID)
			cmd.Stdout, cmd.Stderr = &out, &errb
			if err := cmd.Run(); err != nil {
				t.Logf("icacls /remove:d: %v\n%s%s", err, out.String(), errb.String())
			}
		})

		// Positive control on the plant: the file is there, and this process can
		// no longer stat or read it. That pair is exactly what separates 权限不够
		// from 缺失, so the assertion is made before the product is asked.
		if _, err := os.Stat(path); err == nil {
			t.Fatal("icacls denied nothing: os.Stat still succeeds, so this case cannot show the permission sentence")
		} else if !strings.Contains(err.Error(), "Access is denied") {
			t.Fatalf("stat failed for a reason other than a permission denial: %v", err)
		}
		if _, err := os.ReadFile(path); err == nil {
			t.Fatal("the file is still readable, so this is not the condition under test")
		}
		// The directory listing is the proof of presence: it reads the parent,
		// whose ACL this case never touched.
		entries, err := os.ReadDir(r.h.dir)
		if err != nil {
			t.Fatalf("cannot list the data dir to prove the file is present: %v", err)
		}
		var listed bool
		for _, e := range entries {
			if e.Name() == configFileName {
				listed = true
			}
		}
		if !listed {
			t.Fatalf("config.toml is not in the directory listing, which would make this the missing-file case")
		}

		trail := r.awaitAudit(t, "config: HOT-RELOAD state=not-applied cause=permission")
		for _, other := range []string{"cause=missing", "cause=syntax", "cause=unknown-key",
			"cause=invalid", "cause=unclassified"} {
			if strings.Contains(trail, other) {
				t.Errorf("the permission sentence carries %q:\n%s", other, trail)
			}
		}
		// And the running config stayed live and unchanged: a denied read is not
		// an empty config.
		if got := r.rt.mgr.Config().FS.AllowedDirs; len(got) != 1 {
			t.Errorf("the permission failure replaced the running config: %v", got)
		}
	})
}
