//go:build windows

package risk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 72: the A-tier verdict must not depend on the spelling a path
// arrived in (R17 invariant). These cases sit NEXT TO
// pathresolver_junction_windows_test.go, which is the failing instrument and
// is deliberately left byte-for-byte untouched: its expected values, its
// anchors and its build tags are the measurement, not the thing under edit.
//
// Nothing here needs 8.3 short names, admin rights or a CI environment to
// reproduce the defect: the runner's RUNNER~1-vs-runneradmin discrepancy is
// the general shape "anchor spelling != handle-resolved spelling", and a
// junction as the profile spelling produces that shape on any volume.

// TestClassifyAnchorSpellingIsNotVerdict seeds a real profile tree, points
// USERPROFILE at it through a junction (so the anchor's spelling is not the
// handle's spelling), and classifies the handle-resolved canonical. Every A
// rule must still win. Pre-fix, all four of these degrade: ~/.ssh/** to B
// (id_*), ~/.git-credentials, ~/.aws/credentials and ~/.kube/config to
// ClassNone (Allow, no confirm at all).
func TestClassifyAnchorSpellingIsNotVerdict(t *testing.T) {
	root := t.TempDir()
	profile := filepath.Join(root, "profile")
	alias := filepath.Join(root, "alias")
	if err := os.MkdirAll(profile, 0o700); err != nil {
		t.Fatal(err)
	}
	mkJunction(t, alias, profile)

	t.Setenv("USERPROFILE", alias)
	t.Setenv("HOME", alias)
	t.Setenv("APPDATA", filepath.Join(alias, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(alias, "AppData", "Local"))

	targets := []string{
		filepath.Join(".ssh", "id_testkey"), // ~/.ssh/** (prefix rule) AND id_*
		".git-credentials",                  // exact rule
		filepath.Join(".aws", "credentials"),
		filepath.Join(".kube", "config"),
	}
	for _, rel := range targets {
		p := filepath.Join(profile, rel)
		mustWrite(t, p, "secret")
		res, err := Resolve(p, nil)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", p, err)
		}
		if !strings.EqualFold(res.Canonical, p) {
			t.Fatalf("Resolve(%q) canonical %q: test premise broken (candidate must be the real path)", p, res.Canonical)
		}
		if got := Classify(res.Canonical); got != ClassA {
			t.Errorf("A-list miss for %q: canonical %q classified %v under USERPROFILE=%q, want ClassA",
				rel, res.Canonical, got, alias)
		}
		if d := Gate(res.Canonical, nil); d.Allow || d.NeedL2 {
			t.Errorf("A-list %q through a junction-spelled anchor became approvable: %+v", rel, d)
		}
	}
}

// TestClassifySpellingInvariance states the invariant directly: one and the
// same file, classified through every spelling the OS accepts, yields one
// class. The 8.3 spelling is only folded in when this volume generates short
// names (same probe the ticket-18 case uses); every other spelling is
// unconditional, so the test never quietly stops measuring.
func TestClassifySpellingInvariance(t *testing.T) {
	home := tmpHome(t)
	long := filepath.Join(home, ".ssh", "id_testkey")
	vol := filepath.VolumeName(long)
	share := strings.ToLower(vol[:1]) + "$"
	rest := strings.TrimPrefix(long, vol)

	spellings := map[string]string{
		"long":           long,
		"upper":          strings.ToUpper(long),
		"lower":          strings.ToLower(long),
		"extended":       `\\?\` + long,
		"unc-admin":      `\\localhost\` + share + rest,
		"unc-extended":   `\\?\UNC\localhost\` + share + rest,
		"slashes":        strings.ReplaceAll(long, `\`, "/"),
		"dot-and-dotdot": filepath.Join(home, ".", ".ssh", "..", ".ssh", "id_testkey"),
		"not-yet-there":  filepath.Join(home, ".ssh", "brand-new-key"),
	}
	if s := getShortPath(t, filepath.Join(home, ".ssh", "id_testkey")); s != "" {
		spellings["8.3-short"] = s
	}
	for name, sp := range spellings {
		cands := []string{sp}
		if res, err := Resolve(sp, nil); err == nil {
			cands = append(cands, res.Canonical)
		}
		for _, c := range cands {
			if got := Classify(c); got != ClassA {
				t.Errorf("spelling %s: %q classified %v, want ClassA (classification may not depend on spelling)", name, c, got)
			}
		}
	}
	// Control in the same shape: a non-listed file must stay ClassNone under
	// both spellings, or the loop above would only be proving "everything is A".
	for _, sp := range []string{
		filepath.Join(home, "plain.txt"),
		strings.ToUpper(filepath.Join(home, "plain.txt")),
		filepath.Join(home, "no-such-dir", "notes.txt"),
	} {
		if got := Classify(sp); got != ClassNone {
			t.Errorf("control %q classified %v, want ClassNone", sp, got)
		}
	}
}

// TestCanonicalInputGainsNoSecondForm bounds the blast radius of ticket 72:
// for a path that already IS the handle's spelling (which is what Resolve
// hands to every production caller), formsOf must produce exactly one
// comparison form and call it certain. That is the precise statement of "this
// fix does not change the verdict for a path that has no second spelling";
// the tests above cover the other half.
func TestCanonicalInputGainsNoSecondForm(t *testing.T) {
	home := tmpHome(t)
	for _, p := range []string{
		filepath.Join(home, ".ssh", "id_testkey"),
		filepath.Join(home, ".git-credentials"),
		filepath.Join(home, "plain.txt"),
		home,
		filepath.Join(home, "not-created-yet", "deep.txt"),
	} {
		res, err := Resolve(p, nil)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", p, err)
		}
		f := formsOf(res.Canonical)
		if got := len(f.spellings()); got != 1 {
			t.Errorf("canonical %q produced %d comparison forms %q, want 1 (verdict must not move for a path with one spelling)", res.Canonical, got, f.spellings())
		}
		if !f.certain() {
			t.Errorf("canonical %q is not certain: root=%q real=%q partial=%v", res.Canonical, f.root, f.real, f.partial)
		}
	}
}

// TestOverrideOnlyAcceptsTheResolvedForm pins the asymmetry the deny side does
// NOT share (编排者插单 2026-09-21 10:08): classification may walk every
// spelling because a miss there downgrades A to B, but a B-tier single-file
// override is an ALLOW and may only be satisfied by the one form the OS
// vouches for. The alias here is a junction, i.e. exactly the "one spelling,
// two files" shape an 8.3 short name can take on a busy volume.
//
// Its mutation ("let the override walk all forms again") must go red on the
// first assertion below.
func TestOverrideOnlyAcceptsTheResolvedForm(t *testing.T) {
	root := t.TempDir()
	profile := filepath.Join(root, "profile")
	alias := filepath.Join(root, "alias")
	mustWrite(t, filepath.Join(profile, "proj", "id_demo.txt"), "secret")
	mkJunction(t, alias, profile)

	res, err := Resolve(filepath.Join(profile, "proj", "id_demo.txt"), nil)
	if err != nil {
		t.Fatal(err)
	}
	canonical := res.Canonical                              // the resolved spelling
	viaAlias := filepath.Join(alias, "proj", "id_demo.txt") // a second spelling of it
	if strings.EqualFold(canonical, viaAlias) {
		t.Fatalf("test premise broken: alias %q must be spelled differently from %q", viaAlias, canonical)
	}
	if got := Classify(viaAlias); got != ClassB {
		t.Fatalf("alias spelling classified %v, want ClassB (premise)", got)
	}

	// (1) the override was granted for the ALIAS spelling itself: not a form
	// the OS vouches for, so it must not unlock anything.
	if d := Gate(viaAlias, map[string]bool{strings.ToLower(viaAlias): true}); d.Allow {
		t.Errorf("override recorded on a non-resolved spelling unlocked the file: %+v", d)
	}
	// (2) an override recorded for the resolved path, asked through a
	// different spelling: confirm again, do not wave it through on a string.
	if d := Gate(viaAlias, map[string]bool{strings.ToLower(canonical): true}); d.Allow {
		t.Errorf("non-canonical spelling was allowed on the strength of a canonical override: %+v", d)
	}
	// (3) the reverse: the canonical path presented with an override recorded
	// against some other spelling stays denied.
	if d := Gate(canonical, map[string]bool{strings.ToLower(viaAlias): true}); d.Allow {
		t.Errorf("canonical path unlocked by an override for a different spelling: %+v", d)
	}
	// (4) control, so this is a rule and not a blanket deny: canonical override
	// + canonical path is exactly the L2-confirmed single file the tier allows.
	d := Gate(canonical, map[string]bool{strings.ToLower(canonical): true})
	if !d.Allow || d.Class != ClassB {
		t.Fatalf("canonical override for the canonical path must allow: %+v", d)
	}
}

// TestAListWinsWhereBothTablesHit is AC#3's A∩B shape: one path that matches
// an A rule and a B rule at once must classify A, and B's single-file
// override must not unlock it. The control proves the A verdict comes from the
// anchor and not from the whole tree being denied.
func TestAListWinsWhereBothTablesHit(t *testing.T) {
	home := tmpHome(t)
	both := []string{
		filepath.Join(home, ".ssh", "id_bothkey"),            // ~/.ssh/**  +  id_*
		filepath.Join(home, ".ssh", "prod-credentials.json"), // ~/.ssh/**  +  *credentials*.json
		filepath.Join(home, ".ssh", "deep", "server.pem"),    // ~/.ssh/**  +  *.pem
	}
	for _, p := range both {
		mustWrite(t, p, "secret")
		res, err := Resolve(p, nil)
		if err != nil {
			t.Fatalf("Resolve(%q): %v", p, err)
		}
		for _, sp := range []string{p, res.Canonical} {
			if got := Classify(sp); got != ClassA {
				t.Fatalf("%q classified %v, want ClassA (A must win wherever both tables hit)", sp, got)
			}
			d := Gate(sp, map[string]bool{strings.ToLower(sp): true})
			if d.Allow || d.NeedL2 {
				t.Fatalf("A-list override attempt succeeded for %q: %+v", sp, d)
			}
			if d.Class != ClassA {
				t.Fatalf("%q decision lost its A class: %+v", sp, d)
			}
		}
	}

	// Control: the very same B patterns outside the A tree stay B (approvable
	// after one L2 confirm), and an ordinary file stays ClassNone.
	out := filepath.Join(home, "proj", "id_bothkey")
	if got := Classify(out); got != ClassB {
		t.Fatalf("B-only control %q classified %v, want ClassB", out, got)
	}
	if d := Gate(out, map[string]bool{strings.ToLower(out): true}); !d.Allow {
		t.Fatalf("B-only control must stay overridable by design: %+v", d)
	}
	if got := Classify(filepath.Join(home, "proj", "main.go")); got != ClassNone {
		t.Fatalf("ordinary file classified %v, want ClassNone", got)
	}
}

// TestClassifyFailClosedWhenSpellingUnprovable covers R17's other half: when a
// side cannot be PROVEN expanded, the answer is A, never "looks clean". The
// instrument is a dangling junction — it exists (Lstat sees the reparse
// point) but refuses to open (its target is gone), so no handle can tell us
// what the path routed through it really is called.
func TestClassifyFailClosedWhenSpellingUnprovable(t *testing.T) {
	home := tmpHome(t)
	dangling := filepath.Join(home, "dangling")
	target := filepath.Join(home, "dangling-target")
	mustWrite(t, filepath.Join(target, "here.txt"), "x")
	mkJunction(t, dangling, target)
	if err := os.RemoveAll(target); err != nil {
		t.Fatalf("orphan the junction: %v", err)
	}
	if _, err := os.Lstat(dangling); err != nil {
		t.Fatalf("junction must still be there for the walk to be fooled: %v", err)
	}
	if _, ok := resolveHandle(dangling); ok {
		t.Fatal("test premise broken: a junction whose target is gone must not yield a handle real path")
	}
	if _, err := Resolve(dangling, nil); err == nil {
		t.Fatal("test premise broken: Resolve must still refuse the reparse traversal by default")
	}

	cand := filepath.Join(dangling, "sub", "anything.txt")
	if got := Classify(cand); got != ClassA {
		t.Fatalf("%q classified %v, want ClassA: an unprovable expansion may not fail open", cand, got)
	}
	d := Gate(cand, map[string]bool{strings.ToLower(cand): true})
	if d.Allow || d.NeedL2 || d.Class != ClassA {
		t.Fatalf("unprovable path became approvable through the gate: %+v", d)
	}

	// Controls: the same volume, fully explicable paths, keep their classes.
	for _, tc := range []struct {
		path string
		want Class
	}{
		{filepath.Join(home, "plain.txt"), ClassNone},
		{filepath.Join(home, "no-such-dir", "notes.txt"), ClassNone},
		{filepath.Join(home, "proj", "id_testkey"), ClassB},
		{filepath.Join(home, ".ssh", "id_testkey"), ClassA},
	} {
		if got := Classify(tc.path); got != tc.want {
			t.Errorf("control %q classified %v, want %v", tc.path, got, tc.want)
		}
	}
}
