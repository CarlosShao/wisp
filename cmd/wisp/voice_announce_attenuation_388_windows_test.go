//go:build windows

package main

// Ticket 388's host-side wiring: the live-device interrupt readings
// (Test381r2RealDevicePlayingPhaseAbortWithin400ms and
// Test390r1BallClickRealDeviceWithin400ms) render the REAL ticket sentence
// through the REAL waveOut device and the REAL matcha synthesis, so before
// this file every round - control included - reached the owner's speakers at
// whatever volume the session currently had (the 2026-10-31 startle the
// ticket was opened from).
//
// The attenuation rides internal/audio/aecspike's guard (wave121's delivery):
// the process render session's master volume is written through the proven
// cgo FLOAT call only (59-r1's finding #2: the integer path once persisted a
// near-zero session volume), the waveOut stream itself is never touched, and
// the owner's original value is restored afterwards - so the measured object
// (the abort-latency path) keeps its exact semantics and only loudness
// changes (ticket 388's boundary: the volume axis must not bend the
// <=400ms reading).
//
// A host where the guard cannot drive the session volume skips rather than
// runs loud: a live reading at full volume is exactly the startle 388 exists
// to stop, and claiming the reading without the前置 would be the false green
// this repo refuses.

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio/aecspike"
)

// attenuateAnnounce388 applies the session-volume attenuation a real-device
// announce test must carry before its first announce, and hands the restore
// to t.Cleanup - the aecspike guard's own contract: one Apply, exactly one
// Restore, the owner's value verified by read-back on both sides (the
// ATTENUATE-388 line carries the full original/attenuated/restored triple).
// The value axis is parameterized (AC#1): callers pass the volume, and the
// default 0.15 is aecspike's own constant, not a number copied here.
func attenuateAnnounce388(t *testing.T, attenuated float64) {
	t.Helper()
	g, err := aecspike.ApplyAnnounceAttenuation(attenuated)
	if err != nil {
		t.Skipf("announce attenuation guard unavailable on this host: %v (a live announce without the前置 would run at full volume - exactly the startle ticket 388 stops)", err)
	}
	t.Cleanup(func() {
		readings, err := g.Restore()
		if err != nil {
			t.Errorf("the announce attenuation restore failed: %v (readings %+v)", err, readings)
			return
		}
		t.Logf("ATTENUATE-388 original=%.6f attenuated=%.6f restored=%.6f",
			readings.Original, readings.Attenuated, readings.Restored)
	})
}

// interruptStopBound388 is the same 400ms bound the per-round asserts in the
// two live tests enforce; the distribution form exists so AC#4's "the
// distribution did not move" reading is a line in the log, not an impression.
const interruptStopBound388 = 400 * time.Millisecond

// assertInterruptDistribution388 is AC#4's distribution reading over the stop
// latencies the live rounds collected (unstick and waveOutReset): it logs
// n/min/median/max in one line and asserts every sample inside the SAME 400ms
// bound the per-round asserts enforce - no loosening, no new threshold, the
// volume axis only rides loudness. An empty sample set is a failure, not a
// pass: a distribution reading over zero readings would be the vacuous green
// this repo has been burned by before.
func assertInterruptDistribution388(t *testing.T, label string, ds []time.Duration) {
	t.Helper()
	if len(ds) == 0 {
		t.Fatalf("the %s distribution reading collected zero samples: the rounds did not record, so AC#4's distribution has nothing to stand on", label)
	}
	sorted := append([]time.Duration(nil), ds...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	median := sorted[len(sorted)/2]
	max := sorted[len(sorted)-1]
	t.Logf("%s distribution: n=%d min=%s median=%s max=%s (bound %s)",
		label, len(sorted), sorted[0].Round(time.Millisecond), median.Round(time.Millisecond),
		max.Round(time.Millisecond), interruptStopBound388)
	if max > interruptStopBound388 {
		t.Fatalf("the %s distribution answered in %s at worst, want every sample within %s", label, max, interruptStopBound388)
	}
}

// Test388AnnounceAttenuationWired is the wiring nail: BOTH live-device
// interrupt readings must carry the attenuation前置 BEFORE their first
// announce, and the helper must pair the apply with exactly one t.Cleanup
// restore. Read from disk (runtime.Caller-anchored, the 253r1 roster pin's
// shape) so the wiring cannot rot quietly out of either file - the same
// bidirectional discipline the transport roster holds: a live test without
// the前置 is red, a helper without the Restore pairing is red.
func Test388AnnounceAttenuationWired(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for _, tc := range []struct {
		file   string
		testFn string
	}{
		{"resident_voice_interrupt_381r2_windows_test.go", "Test381r2RealDevicePlayingPhaseAbortWithin400ms"},
		{"resident_voice_interrupt_390r1_windows_test.go", "Test390r1BallClickRealDeviceWithin400ms"},
	} {
		src, err := os.ReadFile(filepath.Join(dir, tc.file))
		if err != nil {
			t.Fatalf("cannot read %s: %v", tc.file, err)
		}
		mark := "func " + tc.testFn + "("
		start := strings.Index(string(src), mark)
		if start < 0 {
			t.Fatalf("%s no longer declares %s: the live reading this nail pins has moved", tc.file, tc.testFn)
		}
		body := string(src)[start:]
		if end := strings.Index(body, "\nfunc "); end >= 0 {
			body = body[:end]
		}
		att := strings.Index(body, "attenuateAnnounce388(t,")
		ann := strings.Index(body, "rv.announce(sentence)")
		if ann < 0 {
			t.Errorf("%s/%s no longer announces the sentence through rv.announce: the reading this nail pins has changed shape", tc.file, tc.testFn)
			continue
		}
		if att < 0 {
			t.Errorf("AC#1 RED: %s/%s announces the real sentence with NO attenuation前置 - the owner's speakers carry every round again (ticket 388's startle is back)", tc.file, tc.testFn)
			continue
		}
		if att > ann {
			t.Errorf("AC#1 RED: %s/%s announces BEFORE applying the attenuation前置 (attenuateAnnounce388 at offset %d, first rv.announce at %d): the control round would run loud", tc.file, tc.testFn, att, ann)
		}
	}
	helperSrc, err := os.ReadFile(filepath.Join(dir, "voice_announce_attenuation_388_windows_test.go"))
	if err != nil {
		t.Fatalf("cannot read the 388 helper: %v", err)
	}
	for _, want := range []string{"ApplyAnnounceAttenuation(", "t.Cleanup(", "Restore()"} {
		if !strings.Contains(string(helperSrc), want) {
			t.Errorf("the 388 helper lost %q: the apply/restore pairing is the aecspike guard's own contract and the owner's volume rides on it", want)
		}
	}
}
