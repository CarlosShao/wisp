//go:build windows

package memory

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestAC2ReclaimWalkNeverDescendsIntoAJunction is ticket 103's tripwire for the
// one fact that keeps PROBE F (RemoveUnlinked deleting a foreign tree's file) out
// of reach of production code today: internal/memory's reclaim walk measures its
// own subtree with filepath.WalkDir, and WalkDir reports a junction as a *file*
// entry and never descends through it. Nothing in the repository asserts that, so
// the next person who "improves" the walk into a hand-rolled recursion over
// os.ReadDir would silently gain the ability to delete somebody else's tree.
//
// Both readings are pinned here: the walk shape itself (leg 1, the tripwire) and
// the outcome of a real reclaim of a junction entry (leg 2/3, the foreign file has
// to be standing at the end). Leg 4 is the positive control that keeps the test
// from passing because nothing was removed at all.
func TestAC2ReclaimWalkNeverDescendsIntoAJunction(t *testing.T) {
	s := openTestStore(t)
	arts := s.ArtifactsDir()

	outside := filepath.Join(s.Dir(), "someone-elses-tree")
	if err := os.MkdirAll(filepath.Join(outside, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	victimFile := filepath.Join(outside, "sub", "keep-me.txt")
	stray79Plant(t, victimFile, 32, time.Hour)

	link := filepath.Join(arts, "junklink")
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, outside).CombinedOutput()
	t.Logf("mklink /J: %s", strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatalf("mklink /J %s -> %s: %v", link, outside, err)
	}

	// Leg 1, the tripwire: what the reclaim walk actually sees.
	var seen []string
	descended := false
	if wErr := filepath.WalkDir(link, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		seen = append(seen, p)
		if p != link {
			descended = true
		}
		return nil
	}); wErr != nil {
		t.Fatalf("WalkDir(%s): %v", link, wErr)
	}
	t.Logf("walk over the junction entry reported %v", seen)
	if descended {
		t.Errorf("tripwire RED: the walk descended into the junction, so removeStray can now build a path into somebody else's tree: %v", seen)
	}

	// Leg 2: reclaiming the junction entry must succeed - it is a stray like any
	// other, and winsec.RemoveUnlinked is the call that unlinks it.
	if rErr := s.removeStray("junklink"); rErr != nil {
		t.Errorf("removeStray(junklink): %v", rErr)
	}
	// Leg 3: the foreign tree behind it is intact, file and directory alike.
	if _, statErr := os.Stat(victimFile); statErr != nil {
		t.Errorf("tripwire RED: reclaiming a junction inside artifacts deleted %s: %v", victimFile, statErr)
	}
	if _, statErr := os.Stat(outside); statErr != nil {
		t.Errorf("reclaiming a junction removed the tree it pointed at: %v", statErr)
	}

	// Leg 4: positive control - the same walk does reclaim an ordinary stray
	// subtree, so legs 2/3 above cannot be scored by a walk that did nothing.
	stray79Plant(t, filepath.Join(arts, "stray-dir", "in.bin"), 128, time.Hour)
	if rErr := s.removeStray("stray-dir"); rErr != nil {
		t.Fatalf("removeStray(stray-dir): %v", rErr)
	}
	if _, statErr := os.Stat(filepath.Join(arts, "stray-dir")); !os.IsNotExist(statErr) {
		t.Errorf("the positive control did not run: stray-dir survived, so this test measured nothing (%v)", statErr)
	}
}
