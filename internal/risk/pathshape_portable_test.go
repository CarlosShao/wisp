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
	if filepath.Separator == '\\' {
		t.Skip("a backslash is not an ordinary filename character on Windows")
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
