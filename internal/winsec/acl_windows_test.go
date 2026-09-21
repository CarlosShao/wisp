//go:build windows

// AC#1 / AC#2 / AC#3 criteria for ticket 89. Everything here is measured with
// `icacls` and the text is printed verbatim: FileInfo.Mode() is not evidence on
// Windows, because the bits it reports are exactly the bits the platform
// ignored.
package winsec_test

import (
	"bytes"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/winsec"
)

const (
	everyoneSID = "S-1-1-0"
	systemSID   = "S-1-5-18"
	adminsSID   = "S-1-5-32-544"
	// Mandatory Label SIDs are integrity levels, not access grants; icacls
	// prints them as ACE-looking lines and they are not principals that can
	// read anything.
	integrityPrefix = "S-1-16-"
	// Well-known SIDs need no name resolution.
	wellKnown = systemSID + "," + adminsSID + "," + everyoneSID
)

func run(t *testing.T, name string, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out.String())
	}
	return out.String()
}

// icaclsRaw returns the security descriptor icacls prints for path - the only
// evidence this ticket accepts. The raw text is logged verbatim, so every
// verdict below can be re-checked against the real ACL instead of against this
// file's parsing of it.
func icaclsRaw(t *testing.T, path string) string {
	t.Helper()
	raw := strings.ReplaceAll(run(t, "icacls", path), "\r\n", "\n")
	t.Logf("icacls %q\n%s", path, raw)
	return raw
}

// sidCache keeps the name->SID resolutions (each is a subprocess) for the
// whole run; a handful of distinct principals covers every file probed.
var sidCache = map[string]string{}

func sidOf(t *testing.T, account string) string {
	t.Helper()
	if strings.HasPrefix(account, "S-1-") {
		return account // icacls could not resolve it to a name and printed the SID.
	}
	if s, ok := sidCache[account]; ok {
		return s
	}
	out := strings.TrimSpace(run(t, "powershell", "-NoProfile", "-Command",
		"(New-Object System.Security.Principal.NTAccount '"+account+
			"').Translate([System.Security.Principal.SecurityIdentifier]).Value"))
	if !strings.HasPrefix(out, "S-") {
		t.Fatalf("cannot resolve principal %q to a SID: %q", account, out)
	}
	sidCache[account] = out
	return out
}

// aclSIDs returns the SIDs (and the printed names) of the principals holding an
// ACE on path. Names are resolved rather than string-matched because icacls
// prints "swq\swq" while the account name is "swq", and it prints the *path*
// glued to the first ACE: the first version of this parser skipped the header
// line and thereby ignored exactly one grant - the foreign Modify ACE that
// leads the ACL of every file in this repository today. That false green is
// why the path is stripped and every remaining ":(" line is counted.
func aclSIDs(t *testing.T, path string) (sids, names []string) {
	t.Helper()
	raw := strings.ReplaceAll(icaclsRaw(t, path), path, "")
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue // blank, or the "Successfully processed" trailer
		}
		principal := line[:i]
		names = append(names, principal)
		sids = append(sids, sidOf(t, principal))
	}
	return sids, names
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

// privateACLError reports every SID holding a grant on path that should not:
// anything but the current user (read/write), SYSTEM and Administrators. That
// is AC#3's accepted substitute for a second account - it does not prove a
// foreign principal *failed* to read, it proves no foreign principal was ever
// given the right to, which is the claim "only I can read this" carries.
func privateACLError(sids []string, me string) string {
	var bad []string
	grantsMe := false
	for _, s := range sids {
		switch {
		case s == me:
			grantsMe = true
		case s == systemSID, s == adminsSID:
		case strings.HasPrefix(s, integrityPrefix):
		default:
			bad = append(bad, s)
		}
	}
	if len(bad) > 0 {
		return "foreign SID(s) " + strings.Join(bad, ", ")
	}
	if !grantsMe {
		return "the current user's own SID is absent, so nothing here is readable at all"
	}
	return ""
}

// assertPrivateACL is the AC#3 criterion.
func assertPrivateACL(t *testing.T, path string) {
	t.Helper()
	sids, names := aclSIDs(t, path)
	if msg := privateACLError(sids, currentSID(t)); msg != "" {
		t.Errorf("NOT PRIVATE %s: %s (principals: %v)", filepath.Base(path), msg, names)
	}
}

// wideParent makes a data directory whose DACL is explicit and *carries a
// foreign read grant* - the shape of a data root on a shared volume, a roaming
// profile, or any parent an operator touched. Inheritance stays on, because
// inheritance is the mechanism under test: what a child gets is decided at
// creation time from this parent.
func wideParent(t *testing.T, name string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	// os.MkdirAll's 0o700 is what internal/memory and internal/agent do today;
	// on Windows it does nothing to access control, which is the point.
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	run(t, "icacls", dir, "/inheritance:r", "/grant:r",
		"*"+currentSID(t)+":(OI)(CI)(F)",
		"*"+systemSID+":(OI)(CI)(F)",
		"*"+adminsSID+":(OI)(CI)(F)",
		"*"+everyoneSID+":(OI)(CI)(RX)", // the foreign read
	)
	return dir
}

// ---------------------------------------------------------------- AC#1 -----

// TestAC1BaselineProductionPaths is the baseline the ticket demands: the four
// private-data classes written the way the repository writes them today,
// through the production entry points, each followed by its real icacls text.
func TestAC1BaselineProductionPaths(t *testing.T) {
	t.Logf("current user SID: %s", currentSID(t))

	// 1. artifact: internal/agent/spill.go:245 - O_CREATE|O_EXCL, 0o600 under
	//    memory/open.go:184's MkdirAll(0o755) artifacts dir.
	artDir := filepath.Join(t.TempDir(), "data", "artifacts")
	if err := os.MkdirAll(artDir, 0o755); err != nil {
		t.Fatal(err)
	}
	art := filepath.Join(artDir, "artifact-89.txt")
	f, err := os.OpenFile(art, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("tool output that only I am supposed to read")); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	icaclsRaw(t, art)
	icaclsRaw(t, artDir)

	// 2. DPAPI blob: internal/secret/store.go:74 - os.WriteFile(..., 0o600)
	//    into store.go:44's MkdirAll(0o700).
	dataDir := filepath.Join(t.TempDir(), "data")
	st, err := secret.NewStore(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:baseline", "hunter2-not-a-real-secret"); err != nil {
		t.Fatal(err)
	}
	icaclsRaw(t, st.Dir())
	for _, e := range mustReadDir(t, st.Dir()) {
		icaclsRaw(t, filepath.Join(st.Dir(), e))
	}

	// 3. wisp.db (+ -wal/-shm, created by SQLite itself): memory.Open, probed
	//    while the store is still open so the sidecar files exist.
	dbDir := filepath.Join(t.TempDir(), "data")
	store, err := memory.Open(dbDir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	for _, name := range []string{"wisp.db", "wisp.db-wal", "wisp.db-shm"} {
		p := filepath.Join(dbDir, name)
		if _, err := os.Stat(p); err != nil {
			t.Logf("BASELINE %s: not present while open (%v)", name, err)
			continue
		}
		icaclsRaw(t, p)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	// 4. staging temp: internal/models/downloader.go:224 MkdirAll(0o755) +
	//    :389 OpenFile(0o644).
	stage := filepath.Join(dataDir, "staging")
	if err := os.MkdirAll(stage, 0o755); err != nil {
		t.Fatal(err)
	}
	dl := filepath.Join(stage, "model.bin")
	g, err := os.OpenFile(dl, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Write([]byte("bytes")); err != nil {
		t.Fatal(err)
	}
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
	icaclsRaw(t, stage)
	icaclsRaw(t, dl)
}

func mustReadDir(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites pins *how* the hole
// works, so the fix aims at the right mechanism: with a parent that grants
// Everyone read, the bare mode idiom hands that grant to the child. This test
// asserts the leak is present - if it ever stops being, the inheritance model
// this ticket is built on changed underneath us and AC#3's criterion has to be
// re-derived rather than quietly kept.
func TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites(t *testing.T) {
	dir := wideParent(t, "data")
	// The exact idiom of internal/agent/spill.go's writeFileExclusive and
	// internal/secret/store.go's Store.
	art := filepath.Join(dir, "artifact.txt")
	if err := os.WriteFile(art, []byte("private"), 0o600); err != nil {
		t.Fatal(err)
	}
	sids, names := aclSIDs(t, art)
	if msg := privateACLError(sids, currentSID(t)); msg == "" {
		t.Fatalf("the baseline is not leaking: the foreign ACE never reached %s (principals %v); "+
			"a criterion that has nothing to catch here proves nothing either", art, names)
	} else {
		t.Logf("baseline leak confirmed: %s carries %s", filepath.Base(art), msg)
	}
}

// ---------------------------------------------------------------- AC#3 -----

// TestAC3SealedWritesCarryNoForeignSID is the criterion the fix has to satisfy,
// written against the winsec API. It is RED while winsec_windows.go still only
// calls os.Chmod - which is the proof the ticket asked for that the criterion
// measures the filesystem and not the code.
func TestAC3SealedWritesCarryNoForeignSID(t *testing.T) {
	dir := wideParent(t, "data")
	cases := []struct {
		name string
		make func(string) error
	}{
		{"artifact-exclusive", func(p string) error {
			return winsec.PrivateFileExclusive(p, []byte("tool output"))
		}},
		{"secret-blob", func(p string) error {
			return winsec.PrivateFile(p, []byte("dpapi-shaped blob"), 0o600)
		}},
		{"staging-temp", func(p string) error {
			return winsec.PrivateFile(p, []byte("downloaded bytes"), 0o644)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := filepath.Join(dir, tc.name)
			if err := tc.make(p); err != nil {
				t.Fatalf("write refused: %v", err)
			}
			assertPrivateACL(t, p)
		})
	}
}

// ---------------------------------------------------------------- AC#2 -----

// TestAC2SealedDirCoversFilesItNeverTouched is the ordering criterion: sealing
// the parent must be enough for children created by code that knows nothing
// about winsec (SQLite's -wal/-shm, a rename temp file). If this passed only
// because winsec also sealed each child it would prove nothing about the
// direction chosen in AC#2 - so the children below are written with bare os
// calls, and winsec is asked only about the directories it created.
func TestAC2SealedDirCoversFilesItNeverTouched(t *testing.T) {
	root := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(root, 0o700); err != nil {
		t.Fatalf("PrivateDirAll: %v", err)
	}
	assertPrivateACL(t, root)

	nested := filepath.Join(root, "artifacts")
	if err := os.Mkdir(nested, 0o700); err != nil { // a dir winsec never saw
		t.Fatal(err)
	}
	run(t, "icacls", nested, "/grant", "*"+everyoneSID+":(OI)(CI)(RX)")
	stray := filepath.Join(nested, "by-someone-else.txt")
	if err := os.WriteFile(stray, []byte("bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Files winsec never touched, created inside the sealed root, must inherit
	// privacy from it - that is the whole coverage argument.
	untouched := filepath.Join(root, "sqlite-like-sidecar.tmp")
	if err := os.WriteFile(untouched, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	assertPrivateACL(t, untouched)

	// Sealing a directory also narrows what is already inside it: children keep
	// the ACL they were *born* with, so a tree that predates this package (or a
	// directory somebody widened underneath it) is only repaired if the seal
	// propagates. Without this leg the "seal the data root once" decision in
	// AC#2 would leave every artifact written before the upgrade wide forever.
	// Somebody widened a directory *underneath* the sealed root (a work dir from
	// before this package existed, an operator's share fix). Asking the public
	// API for privacy again has to take it back - and narrow what is already
	// inside it, because children keep the ACL they were *born* with: without
	// propagation the "seal the data root once" decision in AC#2 would leave
	// every artifact written before the upgrade wide forever.
	if err := winsec.PrivateDirAll(nested, 0o700); err != nil {
		t.Fatalf("PrivateDirAll of a widened existing dir: %v", err)
	}
	assertPrivateACL(t, nested)
	assertPrivateACL(t, stray)
	// And the per-file repair path has to work on its own, for a file nobody
	// asks winsec to place.
	if err := winsec.SealFile(stray); err != nil {
		t.Fatalf("SealFile: %v", err)
	}
	assertPrivateACL(t, stray)
}
