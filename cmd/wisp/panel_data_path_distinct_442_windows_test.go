//go:build windows

package main

// Ticket 442's primary criterion, as ruled on 2026-11-14 15:5x (甲＋落①): a
// fixed-name WebView2 user-data leaf is caught by a RUNTIME observation, not by
// the shape of the source that spells it. The seam this file reads is
// panelDataPathSeam (panel_data_path_seam_windows.go), installed into
// NewPanelManager, so every construction point in the process is visible without
// creating a window.
//
// WHAT IT ASSERTS, AND WHY THAT IS THE PROPERTY. A fixture that derives the
// host's data path from this run's pid/nanotime hands over a different value
// every time it constructs one; a reverted fixed name hands over the same value
// twice. So each covered construction point is driven TWICE in one process and
// the two values must differ. The in-process pair is the proxy for the defect a
// fixed name causes - two concurrent `go test` runs landing on one WebView2
// user-data folder (internal/ball/orb_windows.go:432-458, ticket 314's audit) -
// and, just as important, the pair is what makes this NOT the tautology the
// ticket's ③-a candidate was rejected for: the assertion never asks "does
// panelProfileLeaf look unique", it asks what the construction point handed the
// constructor.
//
// SCOPE, named rather than implied (the static ruler's own honesty rule):
//
//	covered - the construction points a resident panel host can be built at
//	  WITHOUT a window, each driven through its own fixture or the production
//	  assembly function, so the revert that this test must catch is a revert in
//	  those functions. Today: 253r6's windowed fixture, 33r5's windowed fixture,
//	  and the resident boot's own assembly root (the one production site).
//	not covered, on purpose - (1) the 319 holder fixture hands a deliberately
//	  SHARED folder and never asks bringUp (panel_holder_319_windows_test.go:47-53);
//	  that exemption's justification is owned by the static ruler's exemption
//	  table - one entry, with the noWindow clause asserted there - and is NOT
//	  recopied here, which is exactly the cost candidate ③-b was rejected for;
//	  (2) sites that hand "" (the documented "let go-webview2 default it"
//	  contract) or a t.TempDir() from a test body that creates no window;
//	  (3) construction points inside test bodies that cannot be driven without
//	  that test's own window - they remain the static ruler's, which is why it is
//	  kept (frozen) instead of retired.
//
// DIVISION OF LABOUR (the ruling requires each file to state it): this file
// judges the VALUE, the ruler judges the SPELLING. Either alone leaves the other's
// hole open - a fixed name in code this test never drives is invisible here, and a
// per-run leaf assembled in a shape the ruler cannot fold is invisible there.
// Neither file may claim the other's coverage, and the ruler no longer grows
// (ticket 442's freeze): a new shape belongs to a new construction point covered
// here, or to another ticket.

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// panelDataPathReading442 is one value the seam saw, with the construction point
// that handed it over, named by the file:line the recorder read off its OWN
// caller (see panelDataPathCaller442) rather than by a label typed here: a point
// that moves is then named where it now is, and a point that stops calling the
// constructor simply stops appearing - which the count assertions below refuse to
// read as a pass.
type panelDataPathReading442 struct {
	point string
	value string
}

// panelDataPathRecorder442 collects one test's readings. It is the only state the
// seam ever has: the package-level hook stays nil unless this test installs it,
// and t.Cleanup puts it back.
type panelDataPathRecorder442 struct {
	mu    sync.Mutex
	reads []panelDataPathReading442
}

func (r *panelDataPathRecorder442) observe(dataPath string) {
	reading := panelDataPathReading442{point: panelDataPathCaller442(), value: dataPath}
	r.mu.Lock()
	r.reads = append(r.reads, reading)
	r.mu.Unlock()
}

func (r *panelDataPathRecorder442) snapshot() []panelDataPathReading442 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]panelDataPathReading442(nil), r.reads...)
}

// panelDataPathCaller442 names the construction point: the first frame that is
// neither this recorder nor NewPanelManager itself. It walks frames instead of
// taking a fixed depth because the seam is reached through a method value (one
// generated wrapper frame) and because the reader must not start crediting
// panel_host_windows.go - the seam's own send line - as the construction point.
// Inlining cannot fool it either: CallersFrames expands inline frames.
func panelDataPathCaller442() string {
	var pcs [16]uintptr
	frames := runtime.CallersFrames(pcs[:runtime.Callers(2, pcs[:])])
	for {
		frame, more := frames.Next()
		if !panelDataPathOwnFrame442(frame.Function) {
			file := frame.File
			if i := strings.LastIndexAny(file, `\/`); i >= 0 {
				file = file[i+1:]
			}
			return fmt.Sprintf("%s:%d", file, frame.Line)
		}
		if !more {
			break
		}
	}
	return "unknown:0"
}

// panelDataPathOwnFrame442 answers whether fn is a frame of the observation
// itself rather than of the code being observed.
func panelDataPathOwnFrame442(fn string) bool {
	return strings.Contains(fn, "panelDataPathRecorder442") || strings.Contains(fn, "NewPanelManager")
}

// armPanelDataPathRecorder442 installs this test's recorder on the seam and
// restores nil afterwards. A seam that already has an observer is a leak from
// another test, not something to overwrite: overwriting would let a later test
// inherit readings taken on someone else's behalf.
func armPanelDataPathRecorder442(t *testing.T) *panelDataPathRecorder442 {
	t.Helper()
	if panelDataPathSeam != nil {
		t.Fatalf("ticket 442's data-path seam already carries an observer before this test armed it: a previous test leaked its recorder, so any reading here would be someone else's")
	}
	rec := &panelDataPathRecorder442{}
	panelDataPathSeam = rec.observe
	t.Cleanup(func() { panelDataPathSeam = nil })
	return rec
}

// Test442EveryDrivenConstructionPointHandsItsOwnDataPath is the criterion itself.
func Test442EveryDrivenConstructionPointHandsItsOwnDataPath(t *testing.T) {
	rec := armPanelDataPathRecorder442(t)

	// rounds = 2 on purpose: one construction per point cannot tell a per-run
	// value from a fixed name, and a fixed name is the whole subject.
	const rounds = 2

	points := []struct {
		// file is where the construction point lives. It is asserted, not just
		// documented: a fixture that starts constructing from somewhere else (or
		// constructs twice) must not silently pass under this label.
		file  string
		name  string
		build func(t *testing.T)
	}{
		{
			file: "panel_unbound_253r6_windows_test.go",
			name: "253r6's windowed fixture",
			build: func(t *testing.T) {
				t.Helper()
				// The fixture the 2026-11-14 audit mutated back to a fixed name
				// (docs/evidence/s1/today-four-fixes-audit.md section 4, A3).
				startPanelFor253r6(t, &panel.ComposerDispatch{Mode: &recordingModeHandler{}})
			},
		},
		{
			file: "panel_resident_windows_test.go",
			name: "33r5's windowed fixture",
			build: func(t *testing.T) {
				t.Helper()
				startPanelForTest(t)
			},
		},
		{
			file: "panel_resident_windows.go",
			name: "the resident boot's assembly root (the one production site)",
			build: func(t *testing.T) {
				t.Helper()
				// newResidentPanelManager is the row cmd/wisp/resident_windows.go's
				// boot calls; ticket 323 gave it this harness
				// (panel_resident_assembly_323_test.go:44), and the data path it
				// derives is the one a shipped panel actually uses.
				host, _, _, _, err := newResidentPanelManager(newResidentAssemblyDir323(t), nil)
				if err != nil {
					t.Fatalf("newResidentPanelManager refused to assemble, so this construction point produced no reading: %v", err)
				}
				if host == nil {
					t.Fatalf("newResidentPanelManager assembled no host at all, so this construction point produced no reading")
				}
			},
		},
	}

	for round := 0; round < rounds; round++ {
		for _, point := range points {
			point.build(t)
		}
	}

	reads := rec.snapshot()
	if want := rounds * len(points); len(reads) != want {
		t.Fatalf("ticket 442's seam recorded %d data path(s) for %d construction point(s) x %d round(s), want %d. "+
			"A shorter count is a failed measurement, not a pass: either a construction point stopped going through "+
			"NewPanelManager (and this criterion is now blind to it), or the seam was removed from the constructor. "+
			"readings = %v", len(reads), len(points), rounds, want, reads)
	}

	// Group by construction point, keeping first-seen order so the failure text
	// reads in the order the test drove them.
	var order []string
	byPoint := map[string][]string{}
	for _, reading := range reads {
		if _, seen := byPoint[reading.point]; !seen {
			order = append(order, reading.point)
		}
		byPoint[reading.point] = append(byPoint[reading.point], reading.value)
	}
	if len(order) != len(points) {
		t.Errorf("the %d readings came from %d distinct construction point(s), want %d: %v. "+
			"Either a fixture began constructing from a second place, or two fixtures collapsed onto one point.",
			len(reads), len(order), len(points), order)
	}

	// The expected file must carry exactly one construction point, and it must be
	// the one this test drove.
	for _, point := range points {
		var found []string
		for _, seen := range order {
			if strings.HasPrefix(seen, point.file+":") {
				found = append(found, seen)
			}
		}
		if len(found) != 1 {
			t.Errorf("%s: the seam recorded %d construction point(s) in %s (%v), want exactly 1 - the value this "+
				"test judges must come from that file's own constructor call, or the label above is a lie",
				point.name, len(found), point.file, found)
		}
	}

	for _, seen := range order {
		values := byPoint[seen]
		if len(values) != rounds {
			t.Errorf("construction point %s handed the host %d value(s), want %d (one per round): %v - "+
				"a point that no longer constructs on every round is a measurement gap, not a pass",
				seen, len(values), rounds, values)
		}
		firstSeenAt := map[string]int{}
		for i, value := range values {
			if value == "" {
				t.Errorf("construction point %s handed the host an EMPTY data path (round %d). Empty is the documented "+
					"'let go-webview2 default it under %%AppData%%' contract, which is ONE shared folder - for a fixture "+
					"that creates a window that is the fixed-name defect wearing the default's clothes", seen, i+1)
				continue
			}
			if first, dup := firstSeenAt[value]; dup {
				t.Errorf("construction point %s handed the host the SAME data path in two constructions of one run: "+
					"round %d and round %d both got %q. That is the fixed-name shape ticket 442 exists to catch: two "+
					"concurrent runs of this package would land on one WebView2 user-data folder, and the two measured "+
					"deaths for that shape are internal/ball/orb_windows.go:432-458 (a contended CreateEnvironment "+
					"answering S_OK with a NULL environment, and a new environment JOINing a dying browser process)",
					seen, first+1, i+1, value)
				continue
			}
			firstSeenAt[value] = i
		}
	}

	if !t.Failed() {
		t.Logf("ticket 442 runtime criterion: %d construction point(s) x %d round(s) = %d data path(s), "+
			"all distinct within the run and none empty - points %v", len(order), rounds, len(reads), order)
	}
}
