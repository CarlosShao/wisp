//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 307 AC#3: the world this instrument needs is "the profile directory
// has a genuine 8.3 alias". Today that world only exists on machines whose
// TEMP/profile chain happens to carry one, so the five syncdirs_test.go cases
// go silently GREEN (no SKIP, no failure, no self-report) on a machine without
// aliases: the asymmetry they probe never materialises before the assertion is
// even reached. These cases therefore build the alias themselves - a directory
// with a component longer than eight characters, whose short spelling is read
// back from the OS (tryShortPath -> GetShortPathNameW, never typed) - so the
// verdict no longer depends on the shape of the machine.
//
// If the volume assigns no alias, that is a broken premise, not a pass: the
// case fails (ticket 115 - "种不出即 Fatal, the instrument measured nothing").
// It must not SKIP, and the five shared cases in syncdirs_test.go stay
// tag-free on purpose (they are part of the POSIX-evaluated denominator).

const (
	// alias307LongName is > 8 characters and 8.3-legal, so an alias exists
	// whenever the volume generates them at all.
	alias307LongName  = "ThirtySevenAliasProfileHome"
	alias307EvilName  = "ThirtySevenAliasProfileHomeevil"
	alias307CloudName = "ThirtySevenAliasCloudOutside"
	alias307RootName  = "ThirtySevenAliasMissingRoot"
)

// alias307Pair returns the long spelling of a freshly created directory under
// base together with the OS-assigned 8.3 short spelling of the same directory.
// A volume that hands back no alias fails the test rather than skipping it.
func alias307Pair(t *testing.T, base, name string) (long, short string) {
	t.Helper()
	long = filepath.Join(base, name)
	if err := os.MkdirAll(long, 0o755); err != nil {
		t.Fatal(err)
	}
	short, ok := tryShortPath(long)
	if !ok {
		t.Fatalf("AC#3 premise unbuildable: the volume generated no 8.3 alias for %q (SHORT == LONG), so this instrument measured nothing", long)
	}
	return long, short
}

// TestSync307AliasHomeUnderProfileStaysSuspect is the fail-open shape in its
// purest form: the profile anchor is handed to the gate in its SHORT spelling
// while the write target is the LONG spelling of a real file under it.
// Before ticket 307 s.home stayed short, cand was C26-expanded long, isUnder
// compared the two and answered "not under any sync root".
func TestSync307AliasHomeUnderProfileStaysSuspect(t *testing.T) {
	base := t.TempDir()
	longHome, shortHome := alias307Pair(t, base, alias307LongName)
	longEvil, _ := alias307Pair(t, base, alias307EvilName)
	if err := os.MkdirAll(filepath.Join(longHome, "Documents"), 0o755); err != nil {
		t.Fatal(err)
	}

	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: shortHome})
	if p.SyncDetectionComplete() {
		t.Fatal("precondition: nothing here may count as a confirmed sync location")
	}

	st := p.IsSyncPath(filepath.Join(longHome, "Documents", "report.txt"))
	if !st.Sync || st.Root.Provider != "sync-suspect" || st.Root.Source != "suspect-fallback" {
		t.Fatalf("ticket 307 fail-open: write under an alias-homed profile read as not-sync, home=%q got %+v", shortHome, st)
	}

	// The mirrored spelling pair is the same write and must stay suspect too.
	if st := p.IsSyncPath(filepath.Join(shortHome, "Documents", "report.txt")); !st.Sync {
		t.Fatalf("short-spelled target under a short-spelled home must be suspect, got %+v", st)
	}

	// N-5 bounding survives the form comparison: a sibling whose name merely
	// starts with the profile's name is not inside it, in either spelling.
	if st := p.IsSyncPath(filepath.Join(longEvil, "x.txt")); st.Sync {
		t.Fatalf("name-prefix sibling must not be swept into the suspect net, got %+v", st)
	}
	if st := p.IsSyncPath(filepath.Join(base, "nowhere", "x.txt")); st.Sync {
		t.Fatalf("a path outside the profile must not become suspect just because the anchor is short, got %+v", st)
	}
}

// TestSync307AliasShortRootMembership is the nail that tells the two branches
// apart: a sync root that is NOT under the profile and does NOT exist on disk
// (so syncdirs.go add() stores its canon in the short spelling), written to in
// its long spelling. Only the :422 root-membership comparison can flag this -
// the under-profile fallback at :428 is structurally out of reach because the
// path is outside the profile - so a green here is attributable to root
// membership and not to the fallback (the mistake syncdirs_other_test.go:99
// pins for !windows).
func TestSync307AliasShortRootMembership(t *testing.T) {
	base := t.TempDir()
	_, shortHome := alias307Pair(t, base, alias307LongName)
	longCloud, shortCloud := alias307Pair(t, base, alias307CloudName)

	// The root itself is never created: syncdirs.go add() drops it into the
	// expandInput-only branch and stores the short spelling verbatim.
	longRoot := filepath.Join(longCloud, alias307RootName)
	if _, err := os.Stat(longRoot); err == nil {
		t.Fatalf("precondition broken: the fixture root must not exist on disk: %s", longRoot)
	}
	shortRoot := filepath.Join(shortCloud, alias307RootName)

	p := NewProvenance(ProvOptions{NoProbe: true, HomeDir: shortHome, SyncRoots: []SyncRoot{
		{Provider: "Nutstore", Path: shortRoot, Source: "fixture"},
	}})
	if p.SyncDetectionComplete() {
		t.Fatal("precondition: an unverified fixture root may not disarm the fallback")
	}

	target := filepath.Join(longRoot, "out", "secret.md")
	st := p.IsSyncPath(target)
	if !st.Sync {
		t.Fatalf("ticket 307 fail-open at the root-membership branch: write under a short-spelled missing root read as not-sync, root=%q got %+v", shortRoot, st)
	}
	if st.Root.Source == "suspect-fallback" {
		t.Fatalf("verdict must come from ROOT MEMBERSHIP, not from the fallback: %+v", st)
	}
	if st.Root.Provider != "Nutstore" {
		t.Fatalf("root-membership verdict must name the injected provider, got %+v", st)
	}

	// The same directory in its own short spelling is the same write.
	if st := p.IsSyncPath(filepath.Join(shortRoot, "out", "secret.md")); !st.Sync || st.Root.Source == "suspect-fallback" {
		t.Fatalf("short spelling of the same target must match the same root, got %+v", st)
	}

	// Negative control on the same engine: a plain sibling tree with no root
	// anywhere near it stays a plain non-sync write, so the widened comparison
	// is not "everything is a sync root".
	plain := filepath.Join(longCloud, "work", "brand-new.md")
	if st := p.IsSyncPath(plain); st.Sync {
		t.Fatalf("false positive: plain write outside every root flagged, got %+v", st)
	}
}
