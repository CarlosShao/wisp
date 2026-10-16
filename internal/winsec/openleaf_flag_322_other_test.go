//go:build !windows

// Ticket 322, leg 322-r1b. Judged here: the kernel effect of the leaf flag itself,
// deterministically, by name - the half of ticket 322's POSIX leg that does not need a
// racer. This case moved verbatim out of openleaf_322_other_test.go when that file came
// in at 262 lines, two over this project's 260-line budget; the split is by content
// (a deterministic by-name pin here, the probabilistic race there), and no assertion,
// message or threshold was edited on the way.
//
// SCOPE OF THIS CASE, stated because the teeth run measured it: it calls os.OpenFile by
// name and never goes through privateFile, so it pins what oNoFollowPlatform buys at the
// kernel and nothing about the wiring. The two mutations that remove the flag from the
// open inside privateFile, and the mutation that makes the pre-open probe refuse every
// name, both leave THIS case green while the race case next door goes red
// (readings: .scratch/wisp/probes/322/r1b/logs/, reasoning: docs/evidence/s1/322-openleaf.md).
// That is the correct division of labour, not a gap: the wiring is judged by
// TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1.
package winsec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1 pins the constant's kernel
// effect deterministically, including two shapes only this can show: a refusal must not
// delete the link it was told about (322-a1's nearest miss, the unlink on privateFile's
// seal-failure path), and a DANGLING leaf must stop creating an object outside the tree.
func TestAC322LeafFlagRefusesALeafLinkWithoutTouchingTheTarget322r1(t *testing.T) {
	dir, outside := SealableTempDirForTest124(t), SealableTempDirForTest124(t)
	witness := filepath.Join(outside, "someone-elses-file-322r1.bin")
	const body = "the bytes a truncating open would have taken away"
	if err := os.WriteFile(witness, []byte(body), 0o666); err != nil {
		t.Fatal(err)
	}
	before, err := os.Lstat(witness)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ label, link, target string }{
		{"relative sibling", filepath.Join(dir, "sib-322r1.lnk"), filepath.Base(witness)},
		{"dangling, target never created", filepath.Join(dir, "dangling-322r1.lnk"), "never-created-322r1.bin"},
	} {
		if err := os.Symlink(tc.target, tc.link); err != nil {
			t.Fatalf("%s: stage the leaf: %v", tc.label, err)
		}
		f, err := os.OpenFile(tc.link, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|oNoFollowPlatform, 0o600)
		if err == nil {
			_ = f.Close()
			t.Errorf("AC#1 RED %s: the open carrying oNoFollowPlatform took a leaf link, want ELOOP", tc.label)
		} else if !errors.Is(err, syscall.ELOOP) || !strings.Contains(err.Error(), whyOpenELOOP322r1) {
			t.Errorf("AC#1 RED %s: the open answered %v, want the leaf leg carrying %q", tc.label, err, whyOpenELOOP322r1)
		}
		if _, sErr := os.Lstat(tc.link); sErr != nil {
			t.Errorf("AC#1 RED %s: the refusal removed the link itself (%v), so a refusal deleted the entry it was told about", tc.label, sErr)
		}
	}
	if _, err := os.Lstat(filepath.Join(dir, "never-created-322r1.bin")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("AC#1 RED: the dangling-leaf spelling created an object outside the tree")
	}
	after, err := os.Lstat(witness)
	if err != nil {
		t.Fatalf("the witness became unreadable: %v", err)
	}
	if got, rErr := os.ReadFile(witness); rErr != nil || string(got) != body {
		t.Errorf("AC#1 RED: the witness holds %q, want its own bytes untouched", got)
	}
	if after.Size() != before.Size() || !after.ModTime().Equal(before.ModTime()) {
		t.Errorf("AC#1 RED: the witness moved from size %d mtime %s to size %d mtime %s", before.Size(),
			before.ModTime().Format(time.RFC3339Nano), after.Size(), after.ModTime().Format(time.RFC3339Nano))
	}
}
