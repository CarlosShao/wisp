package config

// Ticket 360 - the hop this case is named for: a guarded config writeback puts
// the operator's hand edit on disk and does NOT adopt it into memory.
//
// WHY THIS FILE EXISTS. Ticket 226 pinned that a guarded write keeps a hand
// edited key in the FILE (writeguard_226_test.go, TestAC1...). Ticket 227 pinned
// that the reload tick still sees a hand edit made after a write. Neither case is
// about the moment measured here: after a writeback the file says one thing and
// the process holds another, and that divergence is the thing that has to survive
// until the reload path (D36 rule 1: ConfirmLocked first) takes its verdict on it.
// A probing leg measured that shape (file.ball.size=60 while mem.ball.size=56) and
// that nine existing cases go red if the write path adopts memory instead - but
// none of those nine is named for this hop, so the hop was only ever pinned as a
// side effect of somebody else's assertion. This case is the direct pin.
//
// THE TWO STATES, ON ONE TABLE. State 1 runs the writeback alone and reads both
// sides: the file carries the hand edit, memory carries what the process loaded.
// State 2 runs the Adopt path (CheckAndReload) once more on the SAME manager and
// reads both sides again: now memory is the file's value and the report names the
// hot section. The second state is not decoration - it is what makes the first
// reading a measurement rather than a constant. If the table could only ever
// print one answer, state 1 would prove nothing, so the case fails loudly when
// the control arm did not run or when both arms print the same number.
//
// WHAT THIS CASE DOES NOT SAY. It does not say the divergence is desirable, and
// it does not authorize changing the write path to adopt memory: that would be a
// behaviour change colliding with D36 rule 1 and needs its own ticket. This case
// pins today's shape so a later edit cannot change it silently.

import (
	"slices"
	"testing"
)

// TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit pins the hop directly:
// after a guarded writeback the file holds the hand-edited value and memory still
// holds the value this process loaded; one Adopt pass on the same table then makes
// both sides equal, which is what proves the readings above can tell them apart.
func TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit(t *testing.T) {
	// One table, used by both states: the process loads ball.size=56, the operator
	// then hand-edits the file to 60 while it runs.
	m, path, mutate := newTestManager(t, matrixBase)
	mutate(handEdited)

	// STATE 1 - the writeback alone. No reload runs before these two reads.
	if err := m.AddAllowedDir(testDir); err != nil {
		t.Fatalf("AddAllowedDir: %v", err)
	}
	fileSize := readLoaded(t, path).Ball.Size
	memSize := m.Config().Ball.Size
	// Printed so the two states are readable in the -v log, not only asserted:
	// this is the reading the ticket names (file=60 mem=56).
	t.Logf("STATE 1 after writeback: file.ball.size=%d mem.ball.size=%d", fileSize, memSize)
	if fileSize != 60 {
		t.Errorf("after the writeback the file no longer carries the hand edit: file.ball.size=%d, want 60 "+
			"(the write re-read the file as its base, so a foreign key must survive it)", fileSize)
	}
	if memSize != 56 {
		t.Errorf("the writeback adopted the hand edit into memory: mem.ball.size=%d, want the 56 this "+
			"process loaded; adopting here is the reload path's job, behind ConfirmLocked (D36 rule 1)", memSize)
	}
	// The divergence this case is about, spelled the way the write reports it.
	if diverged := diffKeyPaths(m.Config(), readLoaded(t, path)); !slices.Equal(diverged, []string{"ball.size"}) {
		t.Errorf("the memory/file divergence after the writeback = %v, want exactly [ball.size]", diverged)
	}

	// STATE 2 - the positive control, on the same table: one Adopt pass. The
	// readings start at an impossible value and a flag records that the arm body
	// executed, so deleting this arm leaves state 1 as the only reading and fails
	// the case below instead of passing quietly.
	ranAdopt := false
	adoptedMem, adoptedFile := -1, -1
	var hot []string
	if rep, err := m.CheckAndReload(); err != nil {
		t.Fatalf("CheckAndReload (the Adopt path this control drives): %v", err)
	} else {
		ranAdopt = true
		if rep != nil {
			hot = rep.Hot
		}
		adoptedFile = readLoaded(t, path).Ball.Size
		adoptedMem = m.Config().Ball.Size
	}
	t.Logf("STATE 2 after one Adopt pass: mem.ball.size=%d file.ball.size=%d report.Hot=%v",
		adoptedMem, adoptedFile, hot)
	// THE PRE-CONDITION (ticket 360-r2). Every assertion below reads the control
	// arm's output, so first prove the arm ran at all. This is a separate statement
	// with its own Fatalf and it reads nothing but the sentinels above, so it sits
	// before the value assertions on purpose: neutering those assertions, or letting
	// their two-state guard fall silent, cannot silence this one. Remove the arm and
	// ranAdopt stays false while both sentinels stay at -1, so this fires first.
	if !ranAdopt || adoptedMem < 0 || adoptedFile < 0 {
		t.Fatalf("the Adopt arm did not run: ranAdopt=%v mem.ball.size=%d file.ball.size=%d; "+
			"the positive control never executed, so the state 1 readings above are unbacked",
			ranAdopt, adoptedMem, adoptedFile)
	}
	if adoptedMem != 60 || adoptedFile != 60 || !slices.Contains(hot, "ball") {
		t.Fatalf("the Adopt control did not produce the other state: mem.ball.size=%d file.ball.size=%d "+
			"report.Hot=%v, want 60 / 60 / [ball]; without it the readings above cannot tell two states apart",
			adoptedMem, adoptedFile, hot)
	}
	if adoptedMem == memSize {
		t.Fatalf("both states read mem.ball.size=%d, so this table cannot distinguish 'the writeback kept "+
			"the old value' from 'the Adopt path applied the new one'; the case would be vacuous", adoptedMem)
	}
}
