//go:build windows

package winsec

// Ticket 104 (source: acceptor-ticket89b's R-89b-5).
//
// Ticket 89's fourth fix made *explicit* out-of-band grants loud. This is the
// other half it never covered: a child whose only foreign grant was the one it
// **inherited** from its parent. Narrowing that child shrinks the set of
// principals that can reach it, and the old code reported nothing at all,
// because the detection face skipped every ACE carrying the INHERITED_ACE bit.
//
// The judgment this file pins is "louder", never "quieter" and never "clear
// nothing": the seal still removes the grant (verifyPrivate below proves it),
// and the notice now says *which kind* of ACE went away.
//
// AC#1 - sealing one child that carries an inherited foreign grant must emit
//
//	one WARN whose payload distinguishes `inherited` from `explicit`.
//
// AC#2 - noise bound, both shapes measured:
//
//	leg 1 a {me,SY,BA}-only tree  -> 0 notices, still
//	leg 2 parent re-widened, whole tree propagated -> <=1 notice per child
//	leg 3 parent left wide, children sealed singly -> exactly 1 per child
//
// AC#3 leg (b) - the whitelist is keyed by SID, not by a display name, so our
//
//	own three trustees stay silent whichever way the OS spells them (a
//	name-string whitelist goes red here on a localized machine).

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

// icaclsRaw and currentSID already exist in this directory's *external* test
// package (winsec_test, ticket 89's suite); the cases below have to reach
// narrowNotice, so they run in the internal one and print the same evidence.
func icaclsRaw(t *testing.T, path string) string {
	t.Helper()
	var buf bytes.Buffer
	cmd := exec.Command("icacls", path)
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("icacls %s: %v\n%s", path, err, buf.String())
	}
	return strings.ReplaceAll(buf.String(), "\r\n", "\n")
}

func currentSID(t *testing.T) string {
	t.Helper()
	u, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current: %v", err)
	}
	if !strings.HasPrefix(u.Uid, "S-") {
		t.Fatalf("os/user gave no SID for the current user: %q", u.Uid)
	}
	return u.Uid
}

// sealNotices installs a capture seam for the duration of one test body.
func sealNotices(t *testing.T) *[]narrowNotice {
	t.Helper()
	got := &[]narrowNotice{}
	orig := noticeNarrowed
	noticeNarrowed = func(n narrowNotice) { *got = append(*got, n) }
	t.Cleanup(func() { noticeNarrowed = orig })
	return got
}

func noticesFor(got []narrowNotice, path string) []narrowNotice {
	var out []narrowNotice
	for _, n := range got {
		if strings.EqualFold(n.Path, path) {
			out = append(out, n)
		}
	}
	return out
}

func joinedBuckets(rs []narrowNotice) (explicit, inherited string) {
	for _, n := range rs {
		explicit += strings.Join(n.Principals, ",") + ";"
		inherited += strings.Join(n.Inherited, ",") + ";"
	}
	return explicit, inherited
}

// TestAC1SealFileReportsTheInheritedGrantItCleared is the must-be-red-before-the-
// fix case: the child never held an explicit foreign ACE, only the copy it
// inherited from its parent directory, and sealing that single child drops it.
func TestAC1SealFileReportsTheInheritedGrantItCleared(t *testing.T) {
	got := sealNotices(t)

	root := filepath.Join(t.TempDir(), "store")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	// The parent is widened by an operator's hand, inheritably; the child then
	// picks that grant up at creation time and never has one of its own.
	mustExec(t, "icacls", root, "/grant", "*"+everyoneSID+":(OI)(CI)(RX)")
	child := filepath.Join(root, "readable-by-inheritance.txt")
	if err := os.WriteFile(child, []byte("artifact"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Before reading: the grant really is on the child, and really is inherited.
	beforeACL := icaclsRaw(t, child)
	t.Logf("icacls BEFORE SealFile %s\n%s", child, beforeACL)
	if !namesEveryone(beforeACL) {
		t.Fatalf("the fixture did not produce the inherited grant this test is about:\n%s", beforeACL)
	}
	if !strings.Contains(strings.ToUpper(beforeACL), "(I)") {
		t.Logf("note: icacls renders the inherited copy as: %s", strings.TrimSpace(beforeACL))
	}

	*got = nil
	if err := SealFile(child); err != nil {
		t.Fatalf("SealFile: %v", err)
	}
	afterACL := icaclsRaw(t, child)
	t.Logf("icacls AFTER SealFile %s\n%s", child, afterACL)

	rs := noticesFor(*got, child)
	if len(rs) != 1 {
		t.Fatalf("AC#1: sealing one child that lost an *inherited* foreign grant reported %d notice(s), want exactly 1; all notices: %+v",
			len(rs), *got)
	}
	explicit, inherited := joinedBuckets(rs)
	if !strings.Contains(inherited, everyoneSID) {
		t.Fatalf("AC#1: the notice did not name the cleared *inherited* principal by SID: inherited=%q explicit=%q (icacls before:\n%s)",
			inherited, explicit, beforeACL)
	}
	if strings.Contains(explicit, everyoneSID) {
		t.Fatalf("AC#1: the message must distinguish inherited from explicit, but the cleared grant was filed under explicit: explicit=%q inherited=%q",
			explicit, inherited)
	}
	// Direction of this fix is "louder", not "gentler": the seal still happened.
	if err := verifyPrivate(child); err != nil {
		t.Errorf("the child is not private after the seal, i.e. the notice papered over a seal that cleared nothing: %v", err)
	}
	if namesEveryone(afterACL) {
		t.Errorf("icacls still shows the foreign principal on %s after the seal:\n%s", child, afterACL)
	}
}

// TestAC1DefaultLogSaysInherited pins the same distinction on the *rendered*
// default line, which is what an operator greps. It installs no message builder
// of its own: the default renderer is the thing under test.
func TestAC1DefaultLogSaysInherited(t *testing.T) {
	var buf bytes.Buffer
	origDefault := slog.Default()
	orig := noticeNarrowed // the untouched default renderer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() {
		noticeNarrowed = orig
		slog.SetDefault(origDefault)
	})

	root := filepath.Join(t.TempDir(), "store")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	mustExec(t, "icacls", root, "/grant", "*"+everyoneSID+":(OI)(CI)(RX)")
	child := filepath.Join(root, "log-leg.txt")
	if err := os.WriteFile(child, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := SealFile(child); err != nil {
		t.Fatalf("SealFile: %v", err)
	}
	out := buf.String()
	t.Logf("default rendering: %s", strings.TrimSpace(out))
	if !strings.Contains(out, "level=WARN") || !strings.Contains(out, "inherited") {
		t.Fatalf("the default notifier does not distinguish an inherited clearing: %q", out)
	}
	if !strings.Contains(out, everyoneSID) || !strings.Contains(out, filepath.Base(child)) {
		t.Errorf("the line names neither the principal nor the path an operator has to grep for: %q", out)
	}
}

// TestAC2InheritedNoticeHasANoiseBound is the reason this ticket is not just
// "delete the ID filter": the bound is measured in all three shapes.
func TestAC2InheritedNoticeHasANoiseBound(t *testing.T) {
	t.Run("a_private_tree_with_only_the_private_set_stays_quiet", func(t *testing.T) {
		got := sealNotices(t)
		root := filepath.Join(t.TempDir(), "private")
		if err := PrivateDirAll(root, 0o700); err != nil {
			t.Fatal(err)
		}
		var kids []string
		for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
			p := filepath.Join(root, name)
			if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
				t.Fatal(err)
			}
			kids = append(kids, p)
		}
		sub := filepath.Join(root, "sub")
		if err := PrivateDirAll(sub, 0o700); err != nil {
			t.Fatal(err)
		}
		kids = append(kids, sub)

		*got = nil
		for _, p := range kids {
			if err := SealFile(p); err != nil && !os.IsNotExist(err) {
				t.Fatalf("SealFile(%s): %v", p, err)
			}
		}
		if err := SealDir(sub); err != nil {
			t.Fatal(err)
		}
		if err := SealDir(root); err != nil {
			t.Fatal(err)
		}
		t.Logf("leg 1 (only {me,SY,BA}) WARN count = %d", len(*got))
		if len(*got) != 0 {
			t.Fatalf("AC#2 leg 1: a tree that only ever held the private set produced %d WARN(s), want 0: %+v", len(*got), *got)
		}
	})

	t.Run("parent_policy_change_propagates_without_a_per_child_storm", func(t *testing.T) {
		got := sealNotices(t)
		root := filepath.Join(t.TempDir(), "propagated")
		if err := PrivateDirAll(root, 0o700); err != nil {
			t.Fatal(err)
		}
		var kids []string
		for i := 0; i < 6; i++ {
			p := filepath.Join(root, "k"+string(rune('a'+i))+".txt")
			if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
				t.Fatal(err)
			}
			kids = append(kids, p)
		}
		// The operator widens the parent; the OS pushes the grant into the kids.
		mustExec(t, "icacls", root, "/grant", "*"+everyoneSID+":(OI)(CI)(RX)")
		for _, k := range kids {
			if !strings.Contains(mustExec(t, "icacls", k), everyoneSID) {
				t.Logf("note: %s did not pick up the parent's grant:\n%s", k, mustExec(t, "icacls", k))
			}
		}

		*got = nil
		if err := SealDir(root); err != nil {
			t.Fatalf("SealDir(parent): %v", err)
		}
		parentNotices := len(*got)
		for _, k := range kids {
			if err := SealFile(k); err != nil {
				t.Fatalf("SealFile(child): %v", err)
			}
		}
		perChild := map[string]int{}
		for _, n := range *got {
			perChild[strings.ToLower(n.Path)]++
		}
		t.Logf("leg 2 (parent sealed first, then %d children): total WARN = %d (parent=%d, children=%d), per-child=%v",
			len(kids), len(*got), parentNotices, len(*got)-parentNotices, perChild)
		for _, k := range kids {
			if n := perChild[strings.ToLower(k)]; n > 1 {
				t.Errorf("AC#2 leg 2: %s emitted %d WARN(s), the bound is 1 per child", filepath.Base(k), n)
			}
		}
		if len(*got) > 1 {
			t.Errorf("AC#2 leg 2: whole-tree propagation over a parent the ticket-89 leg already covers must not re-report every child: got %d notice(s): %+v", len(*got), *got)
		}
	})

	t.Run("sealing_only_children_reports_each_of_them_once", func(t *testing.T) {
		got := sealNotices(t)
		parent := filepath.Join(t.TempDir(), "widewide")
		if err := os.MkdirAll(parent, 0o700); err != nil {
			t.Fatal(err)
		}
		mustExec(t, "icacls", parent, "/grant", "*"+everyoneSID+":(OI)(CI)(RX)")
		var kids []string
		for i := 0; i < 4; i++ {
			p := filepath.Join(parent, "c"+string(rune('a'+i))+".txt")
			if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
				t.Fatal(err)
			}
			kids = append(kids, p)
		}
		*got = nil
		for _, k := range kids {
			if err := SealFile(k); err != nil {
				t.Fatalf("SealFile: %v", err)
			}
		}
		perChild := map[string]int{}
		for _, n := range *got {
			perChild[strings.ToLower(n.Path)]++
		}
		t.Logf("leg 3 (parent left wide, %d children sealed singly): total WARN = %d", len(kids), len(*got))
		for _, k := range kids {
			if n := perChild[strings.ToLower(k)]; n != 1 {
				t.Errorf("AC#2 leg 3 / AC#1's shape: %s got %d WARN(s), want exactly 1: %+v", filepath.Base(k), n, *got)
			}
		}
	})
}

// TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem is the anti-localization
// pin: the whitelist is a set of SIDs. Mutating it to compare display-name
// strings makes this red on any machine where "Administrators"/"SYSTEM" is not
// spelled the way the constant was typed, which is exactly the failure the
// ticket refuses to reintroduce.
func TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem(t *testing.T) {
	got := sealNotices(t)
	set := privateSetSIDsBySID(t)
	root := filepath.Join(t.TempDir(), "byname")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(root, "ours.txt")
	if err := os.WriteFile(child, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Re-grant our own trustees on the child by name spelling, explicit and
	// inherited: none of these is a foreign principal, so none may be reported.
	for _, spelling := range []string{"*S-1-5-18", "*S-1-5-32-544", "*" + currentSID(t)} {
		mustExec(t, "icacls", child, "/grant", spelling+":(RX)")
	}
	mustExec(t, "icacls", root, "/grant", "*S-1-5-18:(OI)(CI)(RX)")

	*got = nil
	if err := SealFile(child); err != nil {
		t.Fatalf("SealFile: %v", err)
	}
	t.Logf("own-trustee seal WARN count = %d (%+v)", len(*got), *got)
	for _, n := range noticesFor(*got, child) {
		for _, p := range append(append([]string{}, n.Principals...), n.Inherited...) {
			sid := p
			if i := strings.IndexByte(sid, '('); i > 0 {
				sid = sid[:i]
			}
			if inPrivateSet(set, sid) {
				t.Errorf("a private-set principal was reported as cleared from %s: %q - the whitelist stopped being keyed by SID (all: %+v)",
					filepath.Base(n.Path), p, n)
			}
		}
	}
	if len(*got) != 0 {
		t.Errorf("sealing a child that only ever held our own trustees + their inherited copies emitted %d WARN(s), want 0: %+v", len(*got), *got)
	}
}
