//go:build windows

package tools

// Ticket 252, acceptance leg 252-v1 (second round, non-implementer).
//
// Two cells this file settles, both as READINGS (it judges nothing and fixes
// nothing):
//
//   - AC#1: which containment fails first. Post-fix (2b1a3071) the reference
//     input must show BOTH legs passing; the unaligned spelling fed straight
//     to InAllowlist (no Canonicalize in front) is the shape where the legs
//     can still disagree, and the reading names which one fails first there.
//   - AC#6: does the owner's machine actually hit the defect - answered at
//     the gate level with a NON-EXISTING leaf under the REAL repo root (the
//     gate is what R2/L2 rides on), plus an OS-level equivalence cell inside
//     t.TempDir() only. No file is ever created under the real repo root and
//     none under %APPDATA% proper: the only creations are inside t.TempDir(),
//     the same fixture surface the existing rulers in this package use
//     (TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses et al.), and every
//     created name is removed again.
//
// Zero production bytes touched: the file is mounted over the package with
// `go test -overlay` from .scratch, so the tracked tree gains no file. It
// reuses the existing meter (measure252 / legs252 / shortName252 / aliasOf252
// / repoRoot252 from paths_shortname_252_probe_test.go - same package, same
// build tag), whose standing rules apply unchanged: short names come from the
// OS via GetShortPathNameW, never typed by hand; a missing alias on a fixture
// that needs one FAILS LOUDLY instead of skipping.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// probe252v1 runs one ask through the full production round (Canonicalize,
// then the leg meter) and prints everything verbatim.
func probe252v1(t *testing.T, pc *PathCanonicalizer, asked string, label string) {
	t.Helper()
	canonical, err := pc.Canonicalize(asked)
	if err != nil {
		t.Fatalf("carrier broken: Canonicalize(%q): %v", asked, err)
	}
	l := measure252(pc, canonical)
	t.Logf("CELL %s", label)
	t.Logf("  asked raw           : %q", asked)
	t.Logf("  Canonicalize(asked) : %q", canonical)
	l.log(t, "  ")
}

// TestTicket252V1NonImplAc1Readings - AC#1 as a non-implementer reading.
func TestTicket252V1NonImplAc1Readings(t *testing.T) {
	root := repoRoot252(t)
	short := shortName252(t, root)
	t.Logf("REPO LONG : %q", root)
	t.Logf("REPO SHORT: %q", short)
	t.Logf("EQUAL? %v", strings.EqualFold(root, short))

	pc := NewPathCanonicalizer([]string{root}, nil)
	if u := pc.UnusableRoots(); len(u) != 0 {
		t.Fatalf("carrier broken: the configured root was dropped: %q", u)
	}
	t.Logf("Roots() = %q", pc.Roots())

	// C1: the fix's own reference shape - allowed root long, ask spelled by
	// its 8.3 alias, leaf does not exist yet (a create).
	probe252v1(t, pc, short+pathSep+"q252v1-never-created.txt",
		"C1 short-root + missing-leaf, through Canonicalize (post-fix reference shape)")
	// C2: same but one existing directory deeper.
	probe252v1(t, pc, short+pathSep+"docs"+pathSep+"q252v1-never-created.txt",
		"C2 short + docs + missing-leaf, through Canonicalize")
	// C3: long-name control.
	probe252v1(t, pc, root+pathSep+"docs"+pathSep+"q252v1-never-created.txt",
		"C3 long-name control, through Canonicalize")

	// D1: the UNALIGNED spelling handed straight to InAllowlist - no
	// Canonicalize in front. This is the only input shape where the two legs
	// can disagree at all after the fix, and the meter names which leg fails
	// first on it.
	l := measure252(pc, short+pathSep+"q252v1-never-created.txt")
	t.Logf("CELL D1 raw short + missing-leaf, STRAIGHT to InAllowlist (no Canonicalize)")
	l.log(t, "  ")
}

// TestTicket252V1OwnerReach - AC#6: does an in-root create actually hit.
func TestTicket252V1OwnerReach(t *testing.T) {
	root := repoRoot252(t)
	short := shortName252(t, root)
	pc := NewPathCanonicalizer([]string{root}, nil)
	if u := pc.UnusableRoots(); len(u) != 0 {
		t.Fatalf("carrier broken: the configured root was dropped: %q", u)
	}

	// O1: one create-inside-allowed-root ask arriving with the 8.3 spelling,
	// judged exactly the way rules_gateway.go judges it (Canonicalize, then
	// InAllowlist). The leaf is never created.
	asked := short + pathSep + "q252v1-owner-reach-probe.txt"
	canonical, err := pc.Canonicalize(asked)
	if err != nil {
		t.Fatalf("carrier broken: Canonicalize(%q): %v", asked, err)
	}
	l := measure252(pc, canonical)
	t.Logf("CELL O1 in-root create via 8.3 spelling (full production round)")
	t.Logf("  asked raw           : %q", asked)
	t.Logf("  Canonicalize        : %q", canonical)
	l.log(t, "  ")
	if !l.final {
		t.Errorf("AC#6 HIT: an in-root create ask spelled by its 8.3 alias is refused by the gate "+
			"(canonical=%q leg1=%v leg2=%v roots=%q): the owner's machine DOES hit the defect shape",
			canonical, l.leg1, l.leg2, l.roots)
	}

	// O2: long-name control for the same ask.
	probe252v1(t, pc, root+pathSep+"q252v1-owner-reach-probe.txt",
		"O2 long-name control for the same in-root create ask")

	// O3: negative control - the same kind of ask pointing OUTSIDE the root
	// must stay refused (the fix must not have widened anything).
	tmp := t.TempDir()
	outside := filepath.Join(tmp, "q252v1-outside-ask.txt")
	lo := measure252(pc, outside)
	t.Logf("CELL O3 outside-root ask (long spelling), straight reading")
	lo.log(t, "  ")
	if lo.final {
		t.Errorf("WIDENING RED: an outside-root ask is authorized (asked=%q roots=%q)", outside, lo.roots)
	}

	// O4: OS-level equivalence, inside t.TempDir() only. A newly made dir
	// with a long name either gets an 8.3 alias (then the equivalence is
	// measured for real) or the volume generates none for fresh names (then
	// that itself is the reading: the temp side has no structural feeder).
	scratch := filepath.Join(tmp, "q252v1-scratch-root")
	if err := os.Mkdir(scratch, 0o755); err != nil {
		t.Fatalf("mkdir scratch root: %v", err)
	}
	t.Cleanup(func() { os.Remove(scratch) })
	sShort, why := aliasOf252(t, scratch)
	t.Logf("CELL O4 scratch root=%q alias=%q (%s)", scratch, sShort, why)
	if strings.Contains(why, "same-as-long") {
		t.Logf("O4 READING: this volume generated NO 8.3 alias for a freshly created name; " +
			"the temp side of this machine has no structural short-name feeder (the CI-shape " +
			"feeder RUNNER~1 cannot arise here), so any short-spelled ask can only come from " +
			"outside the machine's own spellings")
		return
	}
	// The scratch root HAS an alias: create one file via the SHORT spelling,
	// then verify the LONG spelling names the same file, then remove it.
	viaShort := filepath.Join(sShort, "q252v1-via-short.txt")
	viaLong := filepath.Join(scratch, "q252v1-via-short.txt")
	if err := os.WriteFile(viaShort, []byte("ticket252v1"), 0o644); err != nil {
		t.Fatalf("create via short spelling %q: %v", viaShort, err)
	}
	t.Cleanup(func() { os.Remove(viaLong) })
	st, err := os.Stat(viaLong)
	if err != nil {
		t.Errorf("OS EQUIVALENCE RED: file created via the short spelling %q is not visible under the long spelling %q: %v",
			viaShort, viaLong, err)
		return
	}
	body, rerr := os.ReadFile(viaLong)
	t.Logf("O4 READING: created via %q, visible via %q (size=%d read-back=%q err=%v): the two spellings name ONE file",
		viaShort, viaLong, st.Size(), string(body), rerr)
}
