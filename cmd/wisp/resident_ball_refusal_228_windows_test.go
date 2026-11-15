//go:build windows

package main

// Ticket 228 AC#9: the ball-refused sentence was said TWICE and only by the arm
// nobody had walked.
//
// The defect, on the pre-ticket tree: startResidentBall's failure branch printed
// the verdict itself (resident_ball_windows.go, the fmt.Printf this slice deleted)
// and the boot report printed rb.statusLine() a second time further down
// runResident - which is the SAME sentence, because statusLine() is the verdict.
// The success arm prints it only in the boot report, so the fix is not symmetry
// for its own sake: the invariant is "the console states this process's ball
// posture exactly once", it holds on the arm the owner's desktop always takes,
// and the failing arm broke it.
//
// THE ARM IS NOW WALKABLE (that is the other half of the AC): the factory hook in
// resident_ball_seam_228_windows.go answers the error ball.New would answer on a
// host with no desktop, so this case drives the REAL production path - same
// verdict assignment, same booking, same boot report - with no desktop and no
// mutation. Before the seam the only way to see the branch was to edit the source.

import (
	"errors"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/observe"
)

func TestAC9BallRefusalIsSaidOnceByTheBootReport(t *testing.T) {
	reg := observe.NewRegistry()
	refusal := errors.New("injected by AC#9: this session has no desktop to draw on")

	var rb *residentBall
	said := captureStdout128(t, func() {
		rb = startResidentBallBounded(t, reg, nil, nil, nil,
			withBallNew228(func(ball.Options) (*ball.Ball, error) { return nil, refusal }))
	})
	if rb == nil {
		t.Fatal("AC#9: startResidentBall answered nil; the refused arm still owes a verdict")
	}
	if rb.b != nil {
		t.Fatalf("AC#9 setup: the injected factory's error produced a live ball (%v)", rb.b)
	}

	// 1. The refusing author says NOTHING on the console. This is the exact line
	//    this slice deleted: restore it and the count goes to 1 and the case reds.
	if n := strings.Count(said, ballAbsentClaim); n != 0 {
		t.Errorf("AC#9 RED: the ball host itself printed the absence claim %d time(s); the boot report is the "+
			"single author of that sentence, and two authors is how the failure arm came to say it twice.\n"+
			"captured:\n%s", n, said)
	}

	// 2. The boot report's own line - the SAME function runResident prints, in both
	//    of its mutually exclusive arms - carries the sentence exactly once, and
	//    still carries the reason the ball is absent.
	for _, entered := range []bool{false, true} {
		line := residentLoopReport228(entered, "stop request arrived during boot", "console entry", rb)
		if n := strings.Count(line, ballAbsentClaim); n != 1 {
			t.Errorf("AC#9 RED: the boot report's line (loop entered=%v) states this process's ball posture %d "+
				"time(s), want exactly 1:\n%s", entered, n, line)
		}
		if !strings.Contains(line, refusal.Error()) {
			t.Errorf("AC#9 RED: the boot report's line (loop entered=%v) dropped the refusal reason; the one "+
				"sentence this process gets about the ball has to name why it is not there:\n%s", entered, line)
		}
	}

	// 3. The verdict itself is still the sentence statusLine answers with, so
	//    nothing that asks this host about its ball can get a different reading.
	if !strings.Contains(rb.statusLine(), ballAbsentClaim) {
		t.Errorf("AC#9 RED: statusLine() no longer carries the absence claim: %q", rb.statusLine())
	}
	if !strings.Contains(rb.statusLine(), refusal.Error()) {
		t.Errorf("AC#9 RED: statusLine() dropped the reason: %q", rb.statusLine())
	}
}
