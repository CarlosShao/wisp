//go:build !windows

// Ticket 125 AC#2's judgment, on the only side of the seam guard that was
// untested on this platform: what the install-time probes are built *on*, and
// whether the guard's refusals survive a change to that material.
//
// Three legs, deliberately different in what they promise:
//
//  1. the fixture leg, which is the defect. Before this ticket's fix both probe
//     sites handed the guard an unresolved os.TempDir(), so on a machine whose
//     temp dir is spelled through a symlink the guard built its probes on a root
//     the built-in floor refuses for its own sake, read the candidate's honest
//     answer as dishonesty, and left the whole C26 pipeline out of the process.
//     RED against that tree, GREEN against the fix - and it reads only
//     resolverProbeShapes, an entry point that exists on both sides, so the red
//     run is a measurement and not a compile error;
//  2. the candidate leg, same shape, asked of the guard instead of the fixture: a
//     candidate that answers the way internal/risk's POSIX pipeline answers has
//     to be installable on a resolved root and demonstrably was not on an
//     unresolved one;
//  3. the refusal leg, which is the boundary condition. The probes have to still
//     refuse every candidate shape ever refused here - pass through, constant,
//     tree-moving, no-account - in BOTH temp spellings. That leg is green before
//     and after by design: it is the nail that stops leg 1 and 2 from being paid
//     for with a softer guard, the trade ticket 107 and the acceptance round of
//     ticket 119 (MUT-119-E) were both sent back for making.
//
// The candidates here are fakes, and fakes on purpose: this binary cannot link
// internal/risk (risk -> observe -> secret -> winsec is already a path, so the
// reverse edge is a cycle), and R-126-1 books the consequence - the install-time
// pair is two spellings under one root, so only a fake can make it say anything
// interesting. The verdict for the real pipeline is pinned by
// TestAC1POSIXSeamHoldsC26Pipeline125 in internal/config, which runs the same
// seam decision in a process that does link internal/risk.
package winsec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// plantedTemps125 is the two temp spellings this ticket turns on: one a real
// directory, one the same directory reached through a symlink.
type plantedTemps125 struct {
	plain  string
	linked string
	real   string
}

// plantTemps125 builds both spellings inside the harness's own temp base and
// asserts the link really is a link. acceptor-ticket119's first attempt built a
// real directory where it meant to plant a symlink and reported a false green off
// it, so the assertions below are part of the instrument, not a courtesy: a
// reading taken through a fixture that silently stopped being a link is worse than
// no reading.
func plantTemps125(t *testing.T) plantedTemps125 {
	t.Helper()
	base := t.TempDir()
	if link := firstLinkInPath125(base); link != "" {
		t.Skipf("fixture: the harness's own temp base %q already reaches itself through the link at %q, "+
			"so neither spelling below is the control side and every reading here would measure one shape", base, link)
	}
	realRoot := filepath.Join(base, "real125", "tmproot125")
	if err := os.MkdirAll(realRoot, 0o700); err != nil {
		t.Fatalf("fixture mkdir %s: %v", realRoot, err)
	}
	linkRoot := filepath.Join(base, "varlink125")
	if err := os.Symlink(filepath.Join(base, "real125"), linkRoot); err != nil {
		t.Skipf("cannot build the symlink this case measures (%v): without it the shape does not exist", err)
	}
	info, err := os.Lstat(linkRoot)
	if err != nil {
		t.Fatalf("fixture Lstat %s: %v", linkRoot, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("fixture: %q is not a symlink (mode %v) - every reading taken through it would be a false green", linkRoot, info.Mode())
	}
	linked := filepath.Join(linkRoot, "tmproot125")
	if firstLinkInPath125(linked) == "" {
		t.Fatalf("fixture: the planted TMPDIR %q has no link in its ancestor chain, so it is the control shape, not the measured one", linked)
	}
	if clean, err := filepath.EvalSymlinks(linked); err != nil || clean != realRoot {
		t.Fatalf("fixture: %q resolves to %q (err %v), want %q", linked, clean, err, realRoot)
	}
	return plantedTemps125{plain: realRoot, linked: linked, real: realRoot}
}

// firstLinkInPath125 names the first prefix of a spelling the filesystem says is a
// link, or "". It asks the filesystem rather than the function under test, which
// is the rule internal/proc/envfork_test.go writes down for fixtures in this
// family (and the reason MUT-119-R's all-green run meant nothing).
func firstLinkInPath125(path string) string {
	cur := path
	for {
		if info, err := os.Lstat(cur); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return ""
		}
		cur = parent
	}
}

// probeRootOf125 takes the root a probe shape was built on back out of the shape
// by undoing its construction (root + separator + ".." + separator + name). It is
// deliberately textual: filepath.Dir Cleans, so walking it would collapse the very
// ".." the shape is made of and report a root the guard never asked about.
func probeRootOf125(shape string) string {
	sep := string(filepath.Separator)
	rest := strings.TrimSuffix(shape, sep+"wisp-103-conformance-probe")
	return strings.TrimSuffix(rest, sep+"..")
}

// TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125 is leg 1, the defect.
func TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125(t *testing.T) {
	planted := plantTemps125(t)
	for _, tc := range []struct{ name, tmpdir string }{
		{"control_plain_temp", planted.plain},
		{"measured_symlink_spelled_temp", planted.linked},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TMPDIR", tc.tmpdir)
			shapes := resolverProbeShapes()
			t.Logf("reading: TMPDIR=%q shapes=%q", tc.tmpdir, shapes)
			for _, shape := range shapes {
				if !filepath.IsAbs(shape) {
					continue // the relative leg is relative on purpose (ticket 103's second hostile shape)
				}
				if link := firstLinkInPath125(shape); link != "" {
					t.Errorf("AC#2 RED: the seam guard built probe %q on a root that reaches itself through the link "+
						"at %s: the floor refuses that spelling for its own sake, so the guard is asking a candidate "+
						"to be trusted about a root nobody resolved (TMPDIR=%q)", shape, link, tc.tmpdir)
				}
				if got := probeRootOf125(shape); got != planted.real {
					t.Errorf("AC#2 RED: with TMPDIR=%q the probe stands on %q, want the resolved root %q", tc.tmpdir, got, planted.real)
				}
			}
		})
	}
}

// passThroughResolver125 answers every input with the input: ticket 103's original
// bypass, and the shape this seam exists to close.
type passThroughResolver125 struct{}

func (passThroughResolver125) Resolve(in string) (string, error) { return in, nil }
func (passThroughResolver125) ResolveAccounted(in string) (string, bool, error) {
	return in, false, nil
}

// constantResolver125 answers every input with one fixed clean spelling: ticket
// 108's P1b "constant" fake, which the containment leg - not the spelling leg -
// has to catch.
type constantResolver125 struct{ fixed string }

func (c constantResolver125) Resolve(string) (string, error) { return c.fixed, nil }
func (c constantResolver125) ResolveAccounted(string) (string, bool, error) {
	return c.fixed, false, nil
}

// treeMovingResolver125 answers both probes by moving them into another tree, and
// moves that tree's own parent into a third: the fake ticket 126 measured on the
// Windows leg, on the shape POSIX can reach.
type treeMovingResolver125 struct{}

func (treeMovingResolver125) Resolve(in string) (string, error) {
	return filepath.Clean(filepath.Join("/wisp125-foreign-tree", filepath.Base(in))), nil
}

func (treeMovingResolver125) ResolveAccounted(in string) (string, bool, error) {
	p, err := treeMovingResolver125{}.Resolve(in)
	return p, false, err
}

// noAccountResolver125 can answer but cannot say whether it moved: R-103-1's
// second half, refused before any probe runs.
type noAccountResolver125 struct{}

func (noAccountResolver125) Resolve(in string) (string, error) { return filepath.Clean(in), nil }

// lexicalFoldResolver125 answers the way internal/risk's pipeline answers on this
// platform: fold the parent pointer out of the spelling, hand back the rest of the
// string it was given, account no rewrite. That is verbatim what the ERROR line
// this ticket quotes records, and it is the candidate whose install the unresolved
// probe root blocked.
type lexicalFoldResolver125 struct{}

func (lexicalFoldResolver125) Resolve(in string) (string, error) {
	return filepath.Clean(in), nil
}

func (lexicalFoldResolver125) ResolveAccounted(in string) (string, bool, error) {
	return filepath.Clean(in), false, nil
}

// linkSpellingResolver125 answers every probe by re-spelling the root it was given
// through the planted link and is otherwise honest: the child's answer does sit
// inside the parent's answer, no answer is a pass-through, and it does carry the
// rewrite account. So the only leg that can refuse it is the guard's re-run of the
// floor on the candidate's answer (resolverConformanceFailure's "a spelling the
// built-in floor itself refuses" branch). The case below asserts that phrase,
// which is what makes the attribution a reading rather than a guess, and MUT-125B
// on the ticket deletes that branch and turns this leg red - the proof that leg 3
// does not rest on the other four legs for its refusals.
type linkSpellingResolver125 struct{ realRoot, linkedRoot string }

func (l linkSpellingResolver125) Resolve(in string) (string, error) { return l.answer(in), nil }
func (l linkSpellingResolver125) ResolveAccounted(in string) (string, bool, error) {
	return l.answer(in), false, nil
}

func (l linkSpellingResolver125) answer(in string) string {
	return strings.Replace(in, l.realRoot, l.linkedRoot, 1)
}

// TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125 is leg 3, the boundary
// condition: changing what the probes stand on may not change what they refuse.
func TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125(t *testing.T) {
	planted := plantTemps125(t)
	fixtures := t.TempDir() // created before any TMPDIR is planted, so it is one tree in both shapes
	for _, tc := range []struct{ name, tmpdir string }{
		{"control_plain_temp", planted.plain},
		{"measured_symlink_spelled_temp", planted.linked},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TMPDIR", tc.tmpdir)
			cases := []struct {
				name  string
				cand  C26Resolver
				whyIt string
			}{
				{"pass_through", passThroughResolver125{}, "it answers a hostile shape by passing it through"},
				{"constant", constantResolver125{fixed: filepath.Join(fixtures, "constant-leaf")}, "it answers every input with one spelling (ticket 108's P1b)"},
				{"tree_moving", treeMovingResolver125{}, "it answers the child outside the tree named for the parent"},
				{"no_account", noAccountResolver125{}, "it cannot answer the rewrite-account question at all"},
				{
					"answer_through_the_link",
					linkSpellingResolver125{realRoot: planted.real, linkedRoot: planted.linked},
					"it answers with a spelling that reaches its tree through a link; the leg that must refuse it is " +
						"the guard's re-run of the floor, which is what the logged reason shows and what MUT-125B " +
						"removes (that run is the attribution proof - this leg only promises a refusal, so that it " +
						"stays green on the pre-fix tree, where the unresolved probe root makes the reason differ)",
				},
			}
			for _, c := range cases {
				reason := resolverConformanceFailure(c.cand)
				t.Logf("reading: shape=%s candidate=%s refused=%v reason=%q", tc.name, c.name, reason != "", reason)
				if reason == "" {
					t.Errorf("AC#2 RED (the refusal side moved): the seam guard accepted %s with TMPDIR=%q - %s. "+
						"Nothing in this ticket may buy the install-survives-the-shape legs with a candidate that was "+
						"previously refused: that is the trade ticket 107 and MUT-119-E were both sent back for: %s",
						c.name, tc.tmpdir, c.whyIt, c.whyIt)
					continue
				}
			}
			// The leg above is only worth anything if the guard can also say yes:
			// the built-in floor is a legal incumbent by design, so "everything is
			// refused" cannot be what makes it green.
			if reason := resolverConformanceFailure(builtinVerifier{}); reason != "" {
				t.Errorf("AC#2: the guard refused its own built-in floor (%q), so the refusals above prove nothing", reason)
			}
		})
	}
}

// TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125 is leg 2, the fix in the guard's
// own terms: a candidate that answers the way this platform's pipeline answers has
// to survive the probes on a resolved root, and demonstrably did not survive them
// on an unresolved one. Both halves are asserted because the second one is what
// keeps the first from being a test that could not have failed.
func TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125(t *testing.T) {
	planted := plantTemps125(t)
	t.Run("control_unresolved_root_still_refused_by_the_floor_itself", func(t *testing.T) {
		// What the guard was actually objecting to, taken apart from the guard: the
		// floor refuses a link-spelled answer on its own, and must keep doing so.
		t.Setenv("TMPDIR", planted.linked)
		if _, err := (builtinVerifier{}).Resolve(planted.linked); err == nil {
			t.Errorf("AC#2 RED: the built-in floor accepted %q, a spelling that reaches itself through a link - the "+
				"leg below is only meaningful while this refusal stands", planted.linked)
		}
	})
	t.Run("measured_symlink_spelled_temp", func(t *testing.T) {
		t.Setenv("TMPDIR", planted.linked)
		if got := os.Getenv("TMPDIR"); got != planted.linked {
			t.Fatalf("fixture: TMPDIR is %q, want %q - the reading below would be attributed to the wrong shape", got, planted.linked)
		}
		reason := resolverConformanceFailure(lexicalFoldResolver125{})
		shapes := resolverProbeShapes()
		t.Logf("reading: TMPDIR=%q shapes=%q refused=%v reason=%q", planted.linked, shapes, reason != "", reason)
		if reason != "" {
			t.Errorf("AC#2 RED: with the process's temp dir spelled through a symlink the seam guard still refuses a "+
				"candidate whose only sin is answering honestly about a root nobody resolved: %q. This is the reading "+
				"acceptor-ticket119 took as winsec.PathResolverInstalled()=<nil> and every seal on the built-in floor.", reason)
		}
	})
	t.Run("control_plain_temp", func(t *testing.T) {
		t.Setenv("TMPDIR", planted.plain)
		if reason := resolverConformanceFailure(lexicalFoldResolver125{}); reason != "" {
			t.Errorf("AC#2: the guard refused a POSIX-honest candidate on the plain temp shape (%q), so the fix did "+
				"not make the answer acceptable - it made the candidate unacceptable", reason)
		}
	})
}
