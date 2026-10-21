//go:build !windows

// Ticket 359's denominator, on the only side where the absoluteness conjunct
// had none.
//
// Ticket 129 put sameAbsoluteness into sameTree and answerInsideTree so that a
// tail would stop vouching for a root it does not hang off. On Windows the pair
// that forced it is "C:wisp\p" next to "C:\wisp\p", and
// absoluteness_attribution_129_windows_test.go pins it there. Here
// filepath.VolumeName is the empty string for every input, so sameVolume is true
// for every pair and the volume leg cannot carry anything: the same sentence
// rests on the absoluteness leg alone, and the pair that states it is "a/b" next
// to "/a/b". The paragraph in resolve.go ending "... the same correction, one
// root over" says those two do not read as one tree here either. Until this file
// nothing in this package asked for that verdict: the ticket 359 census measured
// the whole POSIX package green with the conjunct removed.
//
// Why the component walk cannot answer instead: pathComponents drops every
// leading separator and filepath.VolumeName is empty here, so
// pathComponents("a/b") and pathComponents("/a/b") are the same two names. The
// legs below assert that identity before they assert any verdict, so a later
// change to the walk cannot leave this file quietly measuring nothing.
//
// This is one case with two faces because the sentence appears twice in
// resolve.go - once in sameTree's guard, once in answerInsideTree's guard - and
// each removal leaves the other comparison still answering. A case that asked
// only sameTree would stay green when the answerInsideTree half is removed,
// which is exactly the shape the census booked as "every POSIX case green with
// the conjunct gone".
package winsec

import (
	"path/filepath"
	"slices"
	"testing"
)

// TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359 is ticket 359's
// denominator: on POSIX, "a/b" and "/a/b" are two roots and not one tree.
func TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359(t *testing.T) {
	const rel, abs = "a/b", "/a/b"

	// Instrument sanity first. If the pair ever stops disagreeing in
	// absoluteness, or the walk ever starts telling the two spellings apart,
	// every verdict below stops measuring the leg this ticket is about.
	if filepath.IsAbs(rel) {
		t.Fatalf("the instrument planted nothing: filepath.IsAbs(%q) = true, so this is not a relative-vs-absolute pair", rel)
	}
	if !filepath.IsAbs(abs) {
		t.Fatalf("the instrument planted nothing: filepath.IsAbs(%q) = false, so this is not a relative-vs-absolute pair", abs)
	}
	if v1, v2 := filepath.VolumeName(rel), filepath.VolumeName(abs); v1 != v2 {
		t.Fatalf("this is a volume pair (%q vs %q), not the absoluteness pair this case is about", v1, v2)
	}
	if sameAbsoluteness(rel, abs) {
		t.Fatalf("the instrument planted nothing: sameAbsoluteness(%q, %q) = true, so the conjunct under test has nothing to refuse", rel, abs)
	}
	if !sameVolume(rel, abs) {
		t.Fatalf("the volume leg already refuses this pair (%q, %q), so no verdict below rests on absoluteness", rel, abs)
	}
	if compsRel, compsAbs := pathComponents(rel), pathComponents(abs); !slices.Equal(compsRel, compsAbs) {
		t.Fatalf("the component walk tells these two spellings apart (%q vs %q), so no verdict below rests on absoluteness", compsRel, compsAbs)
	}

	// The verdict the comment states, on both comparisons and in both orders.
	if sameTree(rel, abs) {
		t.Errorf("AC#2 RED: sameTree(%q, %q) = true, but a relative spelling hangs its tail off the process's own directory and an absolute one off the filesystem root - two objects. This is the verdict the resolve.go paragraph states for this platform", rel, abs)
	}
	if sameTree(abs, rel) {
		t.Errorf("AC#2 RED: sameTree(%q, %q) = true (mirror order of the leg above), so the refusal is not symmetric in its arguments", abs, rel)
	}
	if answerInsideTree("a/b/c", abs) {
		t.Errorf("AC#2 RED: answerInsideTree(%q, %q) = true, so a relative child is read as sitting inside the absolute tree %q and the seam guard's containment witness would vouch for a seal across two roots", "a/b/c", abs, abs)
	}
	if answerInsideTree("/a/b/c", rel) {
		t.Errorf("AC#2 RED: answerInsideTree(%q, %q) = true (mirror direction): an absolute child read as sitting inside the relative parent %q", "/a/b/c", rel, rel)
	}

	// CONTROL legs. The rule pinned above is "two roots", not "refuse relative
	// input": without these, a comparison that simply returned false would pass
	// this file.
	if !sameTree(rel, rel) {
		t.Errorf("CONTROL RED: sameTree(%q, itself) = false, so the legs above are refusing rather than comparing", rel)
	}
	if !sameTree(abs, abs) {
		t.Errorf("CONTROL RED: sameTree(%q, itself) = false", abs)
	}
	if !sameTree(rel, "a//b") {
		t.Errorf("CONTROL RED: sameTree denied two separator spellings of one relative path")
	}
	if !answerInsideTree("a/b/c", rel) {
		t.Errorf("CONTROL RED: answerInsideTree denied a relative child inside its own relative tree")
	}
	if !answerInsideTree("/a/b/c", abs) {
		t.Errorf("CONTROL RED: answerInsideTree denied an absolute child inside its own absolute tree")
	}
	// An ordinary tail mismatch is still refused by the walk itself, so the
	// verdicts above are not being carried by anything but the absoluteness leg.
	if sameTree(rel, "/a/c") {
		t.Errorf("CONTROL RED: sameTree read two different tails as one tree; the legs above would then measure the walk, not absoluteness")
	}
}
