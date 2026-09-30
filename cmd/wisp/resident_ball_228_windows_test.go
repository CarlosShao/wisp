//go:build windows

package main

// Ticket 228 AC#1, runtime half: what the shipped no-args process SAYS about
// its ball, and what it BOOKS while saying it.
//
// The shape walk (resident_ball_228_test.go) proves the call site exists. This
// file proves the call runs, in the binary an operator actually launches, and
// that its three sentences agree with each other:
//
//   - the console line about the event loop names this process's ball posture,
//     in exactly one of the two forms the host can produce;
//   - the persistent log holds the matching booking (created XOR refused -
//     never both, never neither), behind the listener's own install record;
//   - when the ball came up, the D38(e)-step-2 teardown booked too, between the
//     booking and the shutdown trail.
//
// No skip and no relaxed assertion: a headless host takes the refusal branch,
// and the refusal branch is still an asserted pair (booking present, teardown
// absent, console line consistent). Deleting the wiring in
// cmd/wisp/resident_windows.go makes this case red because the ball posture is
// then said by nobody in either place.
//
// The three literals below are copied on purpose, the same rule
// resident_sink_nail_127_windows_test.go states for its own sentences: reword
// one and a case here goes red and somebody has to say out loud where the
// record that proves the ball was here now lives.

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

const (
	// ballCreatedMsg is resident_ball_windows.go's booking of a live window.
	ballCreatedMsg = "ball: the resident leg created the floating ball window"
	// ballRefusedMsg is the same file's booking of a failed creation. The leg
	// boots either way; only the sentence changes.
	ballRefusedMsg = "ball: the resident leg could not create the floating ball window"
	// ballStoppedMsg is the teardown booking, i.e. D38(e) step 2's work
	// (internal/proc/shutdown.go:18) done by the leg that owns the hot keys.
	ballStoppedMsg = "ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
	// ballUpClaim / ballAbsentClaim are the two verdict forms the host prints.
	ballUpClaim     = "the floating ball window is up in this process"
	ballAbsentClaim = "this process has NO floating ball window"
)

// TestAC228ResidentLegReportsAndBooksItsBall drives the shipped wisp.exe with no
// arguments - the double-click command line - and reads both of its records.
func TestAC228ResidentLegReportsAndBooksItsBall(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	sinkDir := logSinkDir(dataDir)

	// The console line first: the boot report is printed before the event loop
	// is entered, so its arrival is also the signal that the ball host ran.
	sawVerdict := pollUntil127(200, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
	})
	out := leg.stdout.String()
	up := strings.Count(out, ballUpClaim)
	absent := strings.Count(out, ballAbsentClaim)
	if !sawVerdict || up+absent == 0 {
		leg.stop()
		t.Fatalf("AC#1 RED: the resident process never said which ball posture it has (matches: up=%d absent=%d).\n%s\n"+
			"That sentence is AC#1's readable half: a leg that hosts a window has to be able to name it, and a leg "+
			"that deleted the host call says nothing at all.", up, absent, leg.console())
	}
	if up != 0 && absent != 0 {
		t.Errorf("AC#1 RED: the console claims both %q and %q; the two verdicts are written in one statement and "+
			"cannot both be true of one process.", ballUpClaim, ballAbsentClaim)
	}
	if up+absent != 1 {
		t.Errorf("AC#1 RED: the ball posture is stated %d times, want exactly 1 (up=%d absent=%d).\n%s",
			up+absent, up, absent, leg.console())
	}
	if !strings.Contains(out, residentReachedLoop) {
		t.Errorf("AC#1 RED: the ball sentence did not arrive on the same boot as the event-loop sentence (%q), "+
			"so the leg that reports a ball is not the leg that parks in the loop.\n%s", residentReachedLoop, leg.console())
	}

	// Ask this child to leave the way Ctrl+C does, so the teardown booking and
	// the D38(e) trail are both in the file the same process closes.
	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("GenerateConsoleCtrlEvent(CTRL_BREAK_EVENT, %d) on the leg's own process group: %v\n%s",
			leg.pid(), err, leg.console())
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("the child took the break and did not exit through its own shutdown path with a ball attached.\n%s",
			leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("the child exited with %v, want 0: a ball that cannot be torn down is a shutdown defect, not a "+
			"cosmetic one.\n%s", leg.waitErr, leg.console())
	}

	recs := readResidentSink(t, sinkDir)
	if len(recs) == 0 {
		t.Fatalf("%s holds a pipeline file with no record in it, so the claims below have nothing to be read out of", sinkDir)
	}
	installIdx := installRecordIndex127(recs)
	if installIdx < 0 {
		t.Fatalf("no %q record in the file, so the ball booking has no landmark to be read against (records: %v)",
			residentInstallMsg, msgsOf127(recs))
	}
	createdIdx := indexOfMsg228(recs, ballCreatedMsg)
	refusedIdx := indexOfMsg228(recs, ballRefusedMsg)
	stoppedIdx := indexOfMsg228(recs, ballStoppedMsg)

	switch {
	case createdIdx >= 0 && refusedIdx >= 0:
		t.Errorf("AC#1 RED: the log books both a created ball (record %d) and a refused one (record %d).", createdIdx, refusedIdx)
	case createdIdx < 0 && refusedIdx < 0:
		t.Fatalf("AC#1 RED: neither %q nor %q is in the file this process wrote (records: %v). The console said "+
			"something about a ball; the listener kept no record of it, which is the ticket 117 shape again.",
			ballCreatedMsg, ballRefusedMsg, msgsOf127(recs))
	}

	// The booking must sit behind the listener's own install record: everything
	// this process means to be auditable has to be auditable by itself.
	for _, pair := range []struct {
		what string
		idx  int
	}{
		{"creation booking", createdIdx},
		{"refusal booking", refusedIdx},
		{"teardown booking", stoppedIdx},
	} {
		if pair.idx >= 0 && pair.idx < installIdx {
			t.Errorf("AC#1 RED: the ball %s at record %d precedes the log listener's own install booking at %d, "+
				"so the record is in a file that cannot vouch for it", pair.what, pair.idx, installIdx)
		}
	}

	// Console and log must tell the same story about the same process.
	if up > 0 && createdIdx < 0 {
		t.Errorf("AC#1 RED: the console claims %q while the log holds no %q record.", ballUpClaim, ballCreatedMsg)
	}
	if absent > 0 && refusedIdx < 0 {
		t.Errorf("AC#1 RED: the console claims %q while the log holds no %q record.", ballAbsentClaim, ballRefusedMsg)
	}

	// Teardown, and where it sits: D38(e) step 2's work is done by the host that
	// owns the hot keys, ahead of the sequence's own records.
	switch {
	case createdIdx >= 0 && stoppedIdx < 0:
		t.Fatalf("AC#1 RED: this process created a ball (record %d) and never booked destroying it, so the hot keys "+
			"and the tray icon were left to process death. See internal/proc/shutdown.go:18 (D38(e) step 2).", createdIdx)
	case createdIdx < 0 && stoppedIdx >= 0:
		t.Errorf("AC#1 RED: a teardown was booked (record %d) for a ball that was never created.", stoppedIdx)
	case createdIdx >= 0 && stoppedIdx >= 0 && stoppedIdx < createdIdx:
		t.Errorf("AC#1 RED: the teardown (record %d) is booked before the creation (record %d).", stoppedIdx, createdIdx)
	}
	if stoppedIdx >= 0 {
		firstShutdown := firstIndexOfContains228(recs[stoppedIdx+1:], residentShutdownRecord)
		if firstShutdown < 0 {
			t.Errorf("AC#1 RED: no %q record after the ball teardown, so the D38(e) trail cannot be read as running "+
				"after the listening stopped (records: %v)", residentShutdownRecord, msgsOf127(recs))
		}
	}
	t.Logf("AC#1 READING: console up=%d absent=%d; records created=%d refused=%d stopped=%d install=%d",
		up, absent, createdIdx, refusedIdx, stoppedIdx, installIdx)
}

// TestAC228ExitRequestDuringBootStillLeavesThroughD38E is the other half of
// hosting a ball in this leg: the ball path costs a few hundred milliseconds IN
// FRONT of the point where internal/proc installs the handler that stops the
// process, so an operator who presses Ctrl+C while the window is being created
// used to get 0xc000013a - no D38(e) trail, no tray removal, no hot key
// release, and a log file with nothing flushed in it.
//
// The trigger is the sink FILE appearing, which is the earliest observable
// moment inside that window (the rolling writer opens its file at install time,
// ~270ms before the loop is entered on the bench machine). Delete the
// signal.Notify block in cmd/wisp/resident_windows.go and this case is the red
// name: the child dies on the break and the file it wrote holds no records.
func TestAC228ExitRequestDuringBootStillLeavesThroughD38E(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	sinkDir := logSinkDir(dataDir)

	if !pollUntil127(200, func() bool {
		n, err := observe.CountLogFiles(sinkDir)
		return err == nil && n > 0
	}) {
		leg.stop()
		t.Fatalf("no sink file appeared, so the boot window has nothing to be addressed to.\n%s", leg.console())
	}
	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("GenerateConsoleCtrlEvent(CTRL_BREAK_EVENT, %d): %v\n%s", leg.pid(), err, leg.console())
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("a break sent during boot left the child running: it was neither handled nor survived.\n%s", leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("AC#1 RED: the child answered a boot-time Ctrl+C with %v, want a clean exit 0 through the "+
			"D38(e) order. 0xc000013a here means nothing was registered for the signal while the ball came up.\n%s",
			leg.waitErr, leg.console())
	}

	recs := readResidentSink(t, sinkDir)
	installIdx := installRecordIndex127(recs)
	if installIdx < 0 {
		t.Fatalf("AC#1 RED: a boot-time exit produced no %q record (records: %v)", residentInstallMsg, msgsOf127(recs))
	}
	shutdownIdx := firstIndexOfContains228(recs, residentShutdownRecord)
	if shutdownIdx < 0 {
		t.Fatalf("AC#1 RED: the child exited a boot-time Ctrl+C without booking the D38(e) sequence (records: %v).\n%s",
			msgsOf127(recs), leg.console())
	}
	t.Logf("AC#1 READING: boot-time break exited clean; install=%d shutdown trail starts at %d of %d record(s); ball records created=%d stopped=%d",
		installIdx, shutdownIdx, len(recs),
		indexOfMsg228(recs, ballCreatedMsg), indexOfMsg228(recs, ballStoppedMsg))
}

// indexOfMsg228 finds the first record whose msg is exactly want, or -1.
func indexOfMsg228(recs []sinkInstallRecord, want string) int {
	for i, r := range recs {
		if r.Msg == want {
			return i
		}
	}
	return -1
}

// firstIndexOfContains228 is indexOfMsg228's substring twin.
func firstIndexOfContains228(recs []sinkInstallRecord, sub string) int {
	for i, r := range recs {
		if strings.Contains(r.Msg, sub) {
			return i
		}
	}
	return -1
}
