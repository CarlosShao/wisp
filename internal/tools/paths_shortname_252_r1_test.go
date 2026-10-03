//go:build windows

package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Ticket 252, AC#2 + AC#3 (implementation leg 252-r1).
//
// AC#1's readings already exist and are NOT rebuilt here: the carrier, the leg
// breakdown and the short-name prerequisite all come from
// paths_shortname_252_probe_test.go (measure252 / run252 / shortName252 /
// repoRoot252, same package, same build tag). That file measures the shape and
// fixes nothing - it deliberately leaves the diverging cells with wantFinal -1.
// This file is where those cells get an expectation.
//
// What is judged, in AC#2's words: the allowlist verdict must give ONE answer
// to two spellings of one physical path. The 252-p1 reading pinned the failing
// leg first (paths.go:140, the lexical one) and showed resolvedForm answers
// leniently (long name, ok=true) for a missing leaf, so the alignment the
// orchestrator cut in is: put the asked side into that same form BEFORE the two
// containments run, and run both containments anyway.
//
// Two hard bans from the ticket are the reason some cases below expect false
// and are labelled "adds no authorization":
//   - the two containments may not be merged into one (ticket 107's hole),
//   - "cannot be resolved means not authorized" may not be softened.
// Those two are also covered by paths_ticket107b_probes_test.go (probes A/B/C,
// link inside the root, link as the root) which this leg must keep green
// instead of re-implementing.
//
// No file and no directory is created anywhere: every ask below is a name that
// does not exist yet, which is the shape the ticket says is misjudged. Where a
// short name cannot be obtained the case FAILS (shortName252's own rule); this
// file holds zero t.Skip calls.

// verdict252r1 runs one full production round trip: configure the roots, hand
// the raw ask to Canonicalize (the step immediately before the judgment), then
// read InAllowlist's answer and its leg breakdown through the 252-p1 meter.
func verdict252r1(t *testing.T, rootsCfg []string, asked string) (string, legs252) {
	t.Helper()
	pc := NewPathCanonicalizer(rootsCfg, nil)
	if bad := pc.UnusableRoots(); len(bad) != 0 {
		t.Fatalf("carrier broken: NewPathCanonicalizer dropped a configured root: %q", bad)
	}
	canonical, err := pc.Canonicalize(asked)
	if err != nil {
		t.Fatalf("carrier broken: Canonicalize(%q): %v", asked, err)
	}
	return canonical, measure252(pc, canonical)
}

// repoSubdirs252r1 are real, existing directories of this checkout, used as the
// trees a new file is asked under. They exist, so their 8.3 aliases exist or the
// carrier fails loudly - see shortName252.
func repoSubdirs252r1(t *testing.T, long string) []string {
	t.Helper()
	dirs := []string{long, filepath.Join(long, "docs"), filepath.Join(long, "internal", "tools")}
	for _, d := range dirs {
		if _, err := os.Stat(d); err != nil {
			t.Fatalf("carrier broken: expected tree %q is not there: %v", d, err)
		}
	}
	return dirs
}

// TestTicket252R1ShortSpellingOfNewFileIsAuthorized is AC#2's verdict cell and
// AC#3's carrier: the runner shape (a root that carries a real 8.3 alias, a
// leaf that does not exist yet) replayed on this machine with the alias the OS
// itself reports. Red before the alignment, green after it.
func TestTicket252R1ShortSpellingOfNewFileIsAuthorized(t *testing.T) {
	long := repoRoot252(t)
	short := shortName252(t, long)
	const leaf = "q252r1-new-note.md"
	t.Logf("root long : %q", long)
	t.Logf("root short: %q", short)

	// The four combinations of how each side is spelled. All four name one
	// physical tree and ask for one name that does not exist yet, so after the
	// alignment all four must answer the same way: authorized.
	cases := []one252{
		// p1's Q1: the cell whose first failing leg was the lexical one.
		{name: "A1 roots=LONG asked=SHORT+missing-leaf", rootsCfg: []string{long}, asked: short + pathSep + leaf, wantFinal: 1},
		// p1's Q2: already true, kept as a control so the fix cannot pay for
		// A1 by breaking the opposite direction.
		{name: "A2 roots=SHORT asked=LONG+missing-leaf", rootsCfg: []string{short}, asked: long + pathSep + leaf, wantFinal: 1},
		// p1's Q3b: both sides CONFIGURED short - p1 logged it false.
		{name: "A3 roots=SHORT asked=SHORT+missing-leaf", rootsCfg: []string{short}, asked: short + pathSep + leaf, wantFinal: 1},
		// p1's P1: the same-shape positive control that was already green.
		{name: "A4 roots=LONG asked=LONG+missing-leaf", rootsCfg: []string{long}, asked: long + pathSep + leaf, wantFinal: 1},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) { run252(t, c) })
	}
}

// TestTicket252R1BothSpellingsAnswerTheSame states AC#2's sentence directly, on
// every tree the carrier can name: for one physical new file, the long and the
// short spelling must produce one canonical, one verdict, and - the shape p1
// measured as "两段各说各话" - the two containments must no longer disagree.
func TestTicket252R1BothSpellingsAnswerTheSame(t *testing.T) {
	long := repoRoot252(t)
	short := shortName252(t, long)
	const leaf = "q252r1-round-trip.md"
	for _, dir := range repoSubdirs252r1(t, long) {
		dir := dir
		shortDir := shortName252(t, dir)
		askLong := dir + pathSep + leaf
		askShort := shortDir + pathSep + leaf
		t.Run(filepath.Base(dir), func(t *testing.T) {
			t.Logf("tree long=%q short=%q", dir, shortDir)
			for _, rootsCfg := range [][]string{{long}, {shortDir}, {short}} {
				rootsCfg := rootsCfg
				cLong, lLong := verdict252r1(t, rootsCfg, askLong)
				cShort, lShort := verdict252r1(t, rootsCfg, askShort)
				t.Logf("roots=%v canonical(LONG)=%q", rootsCfg, cLong)
				t.Logf("roots=%v canonical(SHORT)=%q", rootsCfg, cShort)
				t.Logf("roots=%v legs(LONG): leg1=%v leg2=%v final=%v", rootsCfg, lLong.leg1, lLong.leg2, lLong.final)
				t.Logf("roots=%v legs(SHORT): leg1=%v leg2=%v final=%v", rootsCfg, lShort.leg1, lShort.leg2, lShort.final)

				if foldPath(cLong) != foldPath(cShort) {
					t.Errorf("one physical path still reaches the judge as two shapes: %q vs %q (roots=%v)",
						cLong, cShort, rootsCfg)
				}
				if lLong.final != lShort.final {
					t.Errorf("AC#2 RED: InAllowlist gives two answers to one physical path: "+
						"LONG=%v (roots=%v) SHORT=%v (roots=%v); canonicals %q / %q",
						lLong.final, rootsCfg, lShort.final, rootsCfg, cLong, cShort)
				}
				if !lLong.final || !lShort.final {
					t.Errorf("AC#2 RED: both spellings must be authorized inside their own root, got LONG=%v SHORT=%v (roots=%v)",
						lLong.final, lShort.final, rootsCfg)
				}
				// The leg breakdown must still EXPLAIN the boolean after the
				// alignment - p1's consistency identity, now with teeth: if the
				// fix had merged or dropped a containment, leg1 && leg2 would no
				// longer be the shape InAllowlist answers with.
				if lLong.final != (lLong.leg1 && lLong.leg2) || lShort.final != (lShort.leg1 && lShort.leg2) {
					t.Errorf("leg breakdown no longer explains the boolean: LONG %v/%v/%v SHORT %v/%v/%v",
						lLong.leg1, lLong.leg2, lLong.final, lShort.leg1, lShort.leg2, lShort.final)
				}
				if lLong.leg1 != lLong.leg2 || lShort.leg1 != lShort.leg2 {
					t.Errorf("AC#1's split is still live: the two containments disagree on one input "+
						"(LONG leg1=%v leg2=%v, SHORT leg1=%v leg2=%v)",
						lLong.leg1, lLong.leg2, lShort.leg1, lShort.leg2)
				}
			}
		})
	}
}

// TestTicket252R1AlignmentAddsNoAuthorization is the other half of AC#2: the
// same-form step may only ever remove a spelling artifact, never a refusal.
// Every case here is expected false BEFORE and AFTER, so it stays green across
// the change and goes red the moment the alignment starts vouching for trees it
// was not given.
func TestTicket252R1AlignmentAddsNoAuthorization(t *testing.T) {
	long := repoRoot252(t)
	short := shortName252(t, long)
	parent := filepath.Dir(long)
	shortParent := shortName252(t, parent)
	const leaf = "q252r1-out-of-scope.md"

	cases := []struct {
		name     string
		rootsCfg []string
		asked    string
		why      string
	}{
		{
			name: "N1 sibling tree, long", rootsCfg: []string{long},
			asked: parent + pathSep + "q252r1-no-such-tree" + pathSep + leaf,
			why:   "a sibling of the allowed root is not the allowed root",
		},
		{
			name: "N2 sibling tree, short spelling of the parent", rootsCfg: []string{long},
			asked: shortParent + pathSep + "q252r1-no-such-tree" + pathSep + leaf,
			why:   "an 8.3 alias above the root must not become a way in",
		},
		{
			name: "N3 the root's own parent", rootsCfg: []string{long},
			asked: parent + pathSep + leaf,
			why:   "a root authorizes its subtree, not its parent (ticket 107's boundary guard)",
		},
		{
			name: "N4 the root's parent spelled short", rootsCfg: []string{short},
			asked: shortParent + pathSep + leaf,
			why:   "same question with the root configured by its alias",
		},
		{
			name: "N5 one character off the root's tail", rootsCfg: []string{long},
			asked: long + "-evil" + pathSep + leaf,
			why:   "component boundaries may not fold into a prefix match",
		},
		{
			name: "N6 a name on a volume that is not there", rootsCfg: []string{long},
			asked: `Z:\q252r1-dead-volume\` + leaf,
			why:   "unresolvable means unauthorized: paths.go:153's !ok leg must keep refusing",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			_, l := verdict252r1(t, c.rootsCfg, c.asked)
			l.log(t, "  ")
			if l.final {
				t.Errorf("AC#2 RED (widening): InAllowlist authorized %q (%s); leg1=%v leg2=%v folded=%q resolved=%q",
					c.asked, c.why, l.leg1, l.leg2, l.folded, l.rf)
			}
		})
	}

	// N7 isolates the fail-closed leg from the lexical one. This ask sits INSIDE
	// the allowed root, so leg 1 says yes, and it is refused only because the
	// resolver cannot answer for it. If the !ok branch at paths.go:153 were ever
	// softened, or the alignment learned to invent a resolved form, this is the
	// case that notices - so it also asserts that leg1 is true here.
	t.Run("N7 inside the root but unresolvable", func(t *testing.T) {
		unresolvable := long + pathSep + "q252r1<invalid>:*?.md"
		_, l := verdict252r1(t, []string{long}, unresolvable)
		l.log(t, "  ")
		if !l.leg1 {
			t.Fatalf("carrier moved: leg1 = false, so this case no longer isolates the !ok leg (folded=%q)", l.folded)
		}
		if l.rfOK {
			t.Fatalf("carrier moved: resolvedForm answered ok for %q, so the !ok branch has nothing to refuse here", unresolvable)
		}
		if l.final {
			t.Errorf("AC#2 RED (widening): an unresolvable path inside the root was authorized "+
				"(leg1=%v leg2=%v ok=%v rf=%q): \"cannot be resolved means not authorized\" must stay",
				l.leg1, l.leg2, l.rfOK, l.rf)
		}
	})

	// The authorization book itself is untouched by the alignment: Roots() is
	// what the audit line and the operator's config are reconciled against, so
	// the same configured root must fold to the same entry however the ask is
	// later spelled.
	t.Logf("short ask roots check")
	_, lShort := verdict252r1(t, []string{long}, short+pathSep+leaf)
	pc := NewPathCanonicalizer([]string{long}, nil)
	if got, want := pc.Roots(), []string{foldPath(long)}; len(got) != len(want) || got[0] != want[0] {
		t.Errorf("Roots() = %q, want the one folded configured root %q: the alignment must not touch the book",
			got, want)
	}
	if !lShort.final {
		t.Errorf("carrier moved under the test: %q should now be authorized under %q", short+pathSep+leaf, long)
	}
}

// TestTicket252R1CanonicalStillNamesOneTree guards the consumer side of the
// alignment: fs.go / fs_write.go open exactly what Canonicalize returns, so the
// aligned string must still reach the same file the original spelling does.
func TestTicket252R1CanonicalStillNamesOneTree(t *testing.T) {
	long := repoRoot252(t)
	short := shortName252(t, long)
	pc := NewPathCanonicalizer([]string{long}, nil)

	const rel = "docs" + pathSep + "q252r1-round-trip.md"
	cLong, err := pc.Canonicalize(long + pathSep + rel)
	if err != nil {
		t.Fatalf("Canonicalize(long): %v", err)
	}
	cShort, err := pc.Canonicalize(short + pathSep + rel)
	if err != nil {
		t.Fatalf("Canonicalize(short): %v", err)
	}
	if cLong != cShort {
		t.Errorf("two spellings of one path still canonicalize to two strings: %q vs %q", cLong, cShort)
	}
	// The parent the OS opens must be the docs directory itself, and the leaf
	// name must survive the alignment verbatim: only the existing prefix may be
	// re-spelled, never the name the caller asked for.
	if got := filepath.Dir(cLong); !strings.EqualFold(got, filepath.Join(long, "docs")) {
		t.Errorf("aligned canonical %q sits under parent %q, want the docs tree %q",
			cLong, got, filepath.Join(long, "docs"))
	}
	if got := filepath.Base(cLong); got != "q252r1-round-trip.md" {
		t.Errorf("the alignment rewrote the requested name: %q", got)
	}
	if _, err := os.Lstat(cLong); err == nil {
		t.Fatalf("carrier broken: %q exists, this case is supposed to ask for a new file", cLong)
	}
	// What the tool would do with the aligned string: the directory it lands in
	// is one the OS can open, and it is the allowed tree.
	if _, err := os.Stat(filepath.Dir(cLong)); err != nil {
		t.Errorf("fs.write would open a directory that is not there: %v (%q)", err, filepath.Dir(cLong))
	}
}
