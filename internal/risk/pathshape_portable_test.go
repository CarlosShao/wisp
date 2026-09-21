package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 75: C26's canonical form is the REAL PATH of the real file
// (SPEC-06 §4 step 4), and on a POSIX system that shape uses '/'. Until the
// normalizeLocalUNC fix, the resolver rewrote every '/' into '\' on every
// platform, so on Linux the pipeline handed back
// "<cwd>/\home\u\.ssh\id_ed25519": a string that is neither the input nor an
// openable path. Windows could never show this (filepath.Clean had already
// folded the separators, so the rewrite was the identity there), which is the
// only reason it survived.
//
// These tests carry NO build tag on purpose. The shape rule is "platform-shaped",
// so it is assertable on every platform, and a windows-only or linux-only guard
// would leave exactly the hole that let the defect through.

// foreignSep is the separator that is NOT this platform's. A canonical must
// never contain it.
func foreignSep() string {
	if filepath.Separator == '\\' {
		return "/"
	}
	return `\`
}

func TestResolveCanonicalIsPlatformShaped(t *testing.T) {
	home, _, _ := sandbox(t)
	target := filepath.Join(home, "notes.txt")
	if err := os.WriteFile(target, []byte("wisp ticket 75 shape probe"), 0o600); err != nil {
		t.Fatal(err)
	}

	res, err := Resolve(target, nil)
	if err != nil {
		t.Fatalf("Resolve(%s): %v", target, err)
	}
	if !res.Resolved {
		t.Logf("note: handle resolution unavailable on this platform (lexical fallback path)")
	}
	if got := res.Canonical; strings.Contains(got, foreignSep()) {
		t.Errorf("canonical %q carries the foreign separator %q: C26 must hand back the "+
			"platform shape, not the Windows one", got, foreignSep())
	}
	// The load-bearing half: a canonical is used to OPEN things, not only to
	// compare. If the shape is wrong this stat fails.
	st, err := os.Lstat(res.Canonical)
	if err != nil {
		t.Fatalf("the canonical C26 handed back cannot be statted (%v): %q - a canonical that "+
			"the OS cannot open is not a real path", err, res.Canonical)
	}
	if st.IsDir() {
		t.Fatalf("canonical %q resolves to a directory, not the file %q", res.Canonical, target)
	}
}

func TestResolveKeepsSeparatorsOutOfPosixNames(t *testing.T) {
	// The assertion below is about a character POSIX treats as an ordinary
	// filename byte and Windows does not, so it only has content on POSIX.
	// Written as a guard rather than a t.Skip so the SKIP ledger of this
	// package stays exactly what it was before ticket 75 (one entry,
	// TestSyncRegistryProbeLive) and no red can hide behind a skip.
	if filepath.Separator == '\\' {
		return
	}
	home, _, _ := sandbox(t)
	// '\' is a legal character inside a POSIX file name. The old unconditional
	// rewrite made this directory indistinguishable from <home>/we/ird, which
	// is a comparison-collision bug, not a cosmetic one.
	target := filepath.Join(home, `we\ird.txt`)
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	res, err := Resolve(target, nil)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Canonical != target {
		t.Errorf("canonical = %q, want %q (the embedded backslash is part of the name)",
			res.Canonical, target)
	}
	if _, err := os.Lstat(res.Canonical); err != nil {
		t.Errorf("canonical %q is not openable: %v", res.Canonical, err)
	}
}

// TestDoubleSlashSpellingIsNotUNC pins the one ordering the fix introduces: a
// POSIX path that starts with two slashes is an ordinary absolute path, not a
// UNC candidate, and must come back with its own shape. lexCanonical (Clean)
// collapses "//a/b" to "/a/b" before normalizeLocalUNC ever sees it, which is
// what keeps the `//` branch of the UNC guard dead on POSIX. If that order is
// ever swapped, this is the line that says so.
func TestDoubleSlashSpellingIsNotUNC(t *testing.T) {
	if filepath.Separator == '\\' {
		return // guard, not t.Skip: see TestResolveKeepsSeparatorsOutOfPosixNames
	}
	res, err := Resolve(`//localhost/c$/x`, nil)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if strings.Contains(res.Canonical, `\`) {
		t.Errorf("canonical = %q: a POSIX double-slash spelling was rewritten into UNC shape",
			res.Canonical)
	}
	if res.Canonical != "/localhost/c$/x" {
		t.Errorf("canonical = %q, want %q", res.Canonical, "/localhost/c$/x")
	}
}

// TestTierAAnchorsMatchNativeSpelling is the ticket-75 twin of ticket 72's
// anchor invariant: an A-tier path must be classified A on the platform it runs
// on. On Linux every rule below used to fall through to the B tier (id_* /
// config / .git segment) or to ClassNone, because the candidate's canonical and
// the env anchors were built on opposite sides of the separator rewrite.
func TestTierAAnchorsMatchNativeSpelling(t *testing.T) {
	home, lapp, appdata := sandbox(t)
	t.Setenv("HOME", home)        // POSIX home
	t.Setenv("USERPROFILE", home) // Windows home
	t.Setenv("APPDATA", appdata)
	t.Setenv("LOCALAPPDATA", lapp)

	for _, tc := range []struct {
		name string
		path string
	}{
		{"~/.ssh/**", filepath.Join(home, ".ssh", "id_platformkey")},
		{"~/.ssh/ dir itself", filepath.Join(home, ".ssh")},
		{"~/.git-credentials", filepath.Join(home, ".git-credentials")},
		{"~/.aws/credentials", filepath.Join(home, ".aws", "credentials")},
		{"~/.kube/config", filepath.Join(home, ".kube", "config")},
		{".git/config", filepath.Join(home, "repo", ".git", "config")},
		{"%APPDATA%\\wisp\\config.toml", filepath.Join(appdata, "wisp", "config.toml")},
		{"%APPDATA%\\Microsoft\\Protect\\**", filepath.Join(appdata, "Microsoft", "Protect", "S-1-5", "masterkey")},
		{"%LOCALAPPDATA%\\Microsoft\\Credentials\\**", filepath.Join(lapp, "Microsoft", "Credentials", "cred")},
		{"browser credential store", filepath.Join(home, "AppData", "Login Data")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.MkdirAll(filepath.Dir(tc.path), 0o755); err != nil {
				t.Fatal(err)
			}
			res, err := Resolve(tc.path, nil)
			if err != nil {
				t.Fatalf("Resolve(%s): %v", tc.path, err)
			}
			if got := Classify(res.Canonical); got != ClassA {
				t.Fatalf("classify(%q) = %s, want A (rule %s); the A anchor and the candidate "+
					"never met in one shape", res.Canonical, got, tc.name)
			}
		})
	}

	control := filepath.Join(home, "workspace", "plain.txt")
	if err := os.MkdirAll(filepath.Dir(control), 0o755); err != nil {
		t.Fatal(err)
	}
	res, err := Resolve(control, nil)
	if err != nil {
		t.Fatalf("Resolve(control): %v", err)
	}
	if got := Classify(res.Canonical); got != ClassNone {
		t.Fatalf("control %q classified %s, want none: the A rules must not have gone blanket",
			res.Canonical, got)
	}
}

// TestOverrideKeyKeepsPosixBackslashesDistinct is ticket 75 item P4 on the
// ALLOW side of the blacklist. Gate's bOverrides map — the L2 single-file
// confirmation, and the shape config's risk.blacklist_overrides will feed — is
// keyed on normPath output. normPath used to fold '/' into '\' unconditionally,
// and '\' is an ordinary character inside a POSIX file name, so TWO DIFFERENT
// FILES SHARED ONE COMPARISON KEY and a confirmation recorded for either one
// unlocked the other. A merge on the deny side is an annoyance; a merge on the
// allow side is fail-open, which is the same invariant ticket 72 states as
// "the basis for letting something through must be the one spelling the OS
// vouches for".
//
// The pair below is that collision in minimal form, chosen so BOTH members are
// still B-tier after the fix — that is what makes the bleed reachable at all:
//
//	<home>/id_k/x.pem   a directory plus a file           -> B via "*.pem"
//	<home>/id_k\x.pem   ONE file name holding a backslash  -> B via "id_*"
//
// Pre-fix both folded to "\home\u\id_k\x.pem" (baseName cut on the '\' too, so
// the second one read as "x.pem" and inherited the first one's key).
func TestOverrideKeyKeepsPosixBackslashesDistinct(t *testing.T) {
	home, _, _ := sandbox(t)
	if filepath.Separator == '\\' {
		// The same function, in the direction that is real on Windows: the two
		// separator spellings of ONE file must keep folding to ONE key, or the
		// override an operator confirmed would stop applying to its own file.
		// Asserted here so this test has content on both platforms instead of
		// being a silent no-op on one of them.
		a := normPath(filepath.Join(home, "id_k", "x.pem"))
		b := normPath(home + `\id_k/x.pem`)
		if a != b {
			t.Fatalf("Windows must fold both separators into one key: %q vs %q", a, b)
		}
		return
	}

	nested := filepath.Join(home, "id_k", "x.pem")
	flat := filepath.Join(home, `id_k\x.pem`)
	for _, p := range []string{nested, flat} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("key material"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if normPath(nested) == normPath(flat) {
		t.Fatalf("two different POSIX files fold to one comparison key %q", normPath(nested))
	}

	resN, err := Resolve(nested, nil)
	if err != nil {
		t.Fatalf("Resolve(%s): %v", nested, err)
	}
	resF, err := Resolve(flat, nil)
	if err != nil {
		t.Fatalf("Resolve(%s): %v", flat, err)
	}
	for _, tc := range []struct {
		path, canon string
	}{
		{nested, resN.Canonical},
		{flat, resF.Canonical},
	} {
		if got := Classify(tc.canon); got != ClassB {
			t.Fatalf("classify(%q) = %s, want B: the fixture must stay in the tier that has an "+
				"override seam (%s)", tc.canon, got, tc.path)
		}
	}

	// The bleed itself: the operator confirmed exactly ONE of the two files.
	ovr := map[string]bool{normPath(resF.Canonical): true}
	if d := Gate(resF.Canonical, ovr); !d.Allow {
		t.Fatalf("the override must apply to the file it was recorded for: %+v", d)
	}
	if d := Gate(resN.Canonical, ovr); d.Allow {
		t.Fatalf("CROSS-FILE BLEED (fail-open): the override recorded for %q unlocked the different "+
			"file %q (verdict %+v)", resF.Canonical, resN.Canonical, d)
	}
}

// TestSyncAncestorWalkVerifiesNativeChain pins the syncdirs half of the defect:
// deepestExistingAncestor hands each prefix to os.Lstat, so a '\' joined walk
// verified NOTHING on Linux and every write fail-closed to sync-suspect.
func TestSyncAncestorWalkVerifiesNativeChain(t *testing.T) {
	home, _, _ := sandbox(t)
	docs := filepath.Join(home, "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	canon, err := Resolve(docs, nil)
	if err != nil {
		t.Fatalf("Resolve(%s): %v", docs, err)
	}
	anc, rest := deepestExistingAncestor(canon.Canonical + string(filepath.Separator) + "fresh.md")
	if anc == "" {
		t.Fatalf("the ancestor walk verified nothing below %q, although every component of it "+
			"exists: on POSIX this made every single write sync-suspect", canon.Canonical)
	}
	if len(rest) != 1 || rest[0] != "fresh.md" {
		t.Fatalf("rest = %v, want [fresh.md]", rest)
	}
	if strings.Contains(anc, foreignSep()) {
		t.Errorf("verified ancestor %q carries the foreign separator", anc)
	}
}
