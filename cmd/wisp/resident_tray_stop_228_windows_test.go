//go:build windows

package main

// Ticket 228 AC#11, the half that needs no desktop: the tray's 退出 item has an
// executor now, and that executor is internal/proc's stop-request bit (form 乙-1,
// orchestrator ruling A493).
//
// What these four cases read, one per clause of the ticket:
//
//   - TestAC11StopBitIsBuiltAtBoot - ⓐ's first clause. The bit exists because
//     proc.Boot built it, and the handle taken before anybody asks is the handle
//     the request closes; a bit created inside RunEventLoop fails this reading,
//     which is the whole reason the tray callback could not use one.
//   - TestAC11TrayExitExecutorIsMountedAtTheBallHost - the dispatch's extra pin:
//     the value production injects at the assembly root, driven directly, plus a
//     walk proving that is the value runResident hands the host that builds
//     Events.OnTrayExit.
//   - TestAC11TrayExitTriggerWalksTheTenSteps - ⓑ-1: real proc.Boot, real
//     RegisterShutdownHook, the tray's own trigger, and the ten frozen steps read
//     in the process that ran them (the shape resident_approval_246_windows_test.go
//     sets).
//   - TestAC11StopRequestDuringTheBootWindowIsNotSwallowed - ⓐ's second clause,
//     in residentBootWindowExit228, the code the assembly root really runs.
//
// ⛔ NONE OF THIS IS "the tray was clicked". The menu's selection value comes only
// from TrackPopupMenu's synchronous return inside internal/ball's showMenu, and
// that file's own header says WM_COMMAND is NOT used - so nothing here can inject
// a selection, and ledger A495 closed the "post a WM_COMMAND with the menu id"
// road by name. Only the winlive companion
// (resident_tray_stop_live_228_windows_test.go) may speak about a request reaching
// a live process, and it pins the exit code and the audit tail, never a 1..10
// sequence: internal/proc/shutdown.go books a line for skipped, failed and fast
// steps only, so the successful steps are structurally unreadable on disk.
//
// The step names below are copied as literals on purpose, the rule
// resident_ball_228_windows_test.go states for its own sentences: if the order or
// a label moves, a case here goes red and somebody has to say out loud where the
// proof of the D38(e) trail lives instead.

import (
	"context"
	"go/ast"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
)

// frozenShutdownSteps228 is D38(e)'s audit roster in the order
// RunShutdownSequence walks it (internal/proc/shutdown.go's shutdownStepNames).
var frozenShutdownSteps228 = []string{
	"scheduler-close",
	"stop-hotkey-kws",
	"cancel-task-roots",
	"stop-audio",
	"release-speech-sessions",
	"destroy-panel-webview",
	"flush-logs-close-db",
	"free-os-memory",
	"close-job-object",
	"exit",
}

// registeredSteps228 are the three hookable slots this case hands real owners, so
// "the trail says the step happened" is checked against "its owner executed". All
// three are inside the closed eight-slot roster; this ticket adds no slot.
var registeredSteps228 = []proc.ShutdownStep{
	proc.StepSchedulerClose,
	proc.StepStopHotkeyKWS,
	proc.StepCancelTasks,
}

// loopReturnBudget228 bounds the wait for RunEventLoop to answer a recorded
// request. The loop's own tick is 50 ms, so this is two orders of magnitude over
// the mechanism; it exists because a loop that lost its stop check does not fail,
// it blocks, and a blocked test binary is not a red sentence naming a clause.
const loopReturnBudget228 = 5 * time.Second

// bootRuntime228 boots a real runtime in the test env: no session mutex (parallel
// tests), a fresh registry, and the stop bit Boot is required to build.
func bootRuntime228(t *testing.T) *proc.Runtime {
	t.Helper()
	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("proc.Boot(test): %v", err)
	}
	return rt
}

// trayTrigger228 fires the tray's Exit item the way the assembly root wired it:
// the production accessor supplies the function value, recordTrayExit is the body
// Events.OnTrayExit holds. Nothing here re-types either half.
func trayTrigger228(rt *proc.Runtime) {
	recordTrayExit(residentTrayStop228(rt))
}

// runLoopBounded228 runs rt.RunEventLoop and gives it loopReturnBudget228 to come
// back, so the mutation that removes the loop's stop check reddens this line
// instead of the package timeout.
func runLoopBounded228(t *testing.T, rt *proc.Runtime) string {
	t.Helper()
	returned := make(chan string, 1)
	go func() { returned <- rt.RunEventLoop() }()
	select {
	case reason := <-returned:
		return reason
	case <-time.After(loopReturnBudget228):
		t.Fatalf("AC#11 RED: RunEventLoop never returned, although a stop request is recorded on this runtime "+
			"(waited %s). internal/proc's loop select is the only thing in this package that turns a recorded "+
			"request into a return, so this is the clause that reads that check.", loopReturnBudget228)
		return ""
	}
}

// namesOf228 renders a trail for the reading line.
func namesOf228(records []proc.StepRecord) []string {
	out := make([]string, 0, len(records))
	for _, r := range records {
		out = append(out, r.Name)
	}
	return out
}

// TestAC11StopBitIsBuiltAtBoot is ⓐ's first clause, read off a runtime nobody has
// asked to stop yet.
func TestAC11StopBitIsBuiltAtBoot(t *testing.T) {
	rt := bootRuntime228(t)

	early := rt.StopSignal()
	if early == nil {
		t.Fatal("AC#11 RED: a runtime built through proc.Boot has no stop-request bit (StopSignal() answers nil), " +
			"so the tray's Exit item has nothing to record into before the event loop exists")
	}
	if rt.StopRequested() {
		t.Fatalf("AC#11 RED: a freshly booted runtime already reports a stop request (%q)", rt.StopReason())
	}
	if rt.StopReason() != "" {
		t.Fatalf("AC#11 RED: a freshly booted runtime carries a stop reason %q", rt.StopReason())
	}

	// The latch, on the handle taken BEFORE the request: a second channel closing
	// later would satisfy "a request was recorded" while losing this one.
	if !rt.RequestStop("first") {
		t.Fatal("AC#11 RED: the first RequestStop on a booted runtime was refused")
	}
	select {
	case <-early:
	default:
		t.Fatal("AC#11 RED: the stop handle taken at boot did not ready on the first request, so a request " +
			"recorded before the loop is entered is not the request the loop would see")
	}
	if rt.RequestStop("second") {
		t.Error("AC#11 RED: a second RequestStop claimed it was the request on record")
	}
	if got := rt.StopReason(); got != "first" {
		t.Errorf("AC#11 RED: StopReason() = %q, want the first request's reason %q", got, "first")
	}

	// The two shapes that must not panic on a UI thread: a runtime nobody booted,
	// and no runtime at all.
	handBuilt := &proc.Runtime{}
	if !handBuilt.RequestStop("no boot") {
		t.Error("AC#11 RED: a Runtime assembled without Boot refused to record a stop request")
	}
	var noRuntime *proc.Runtime
	if noRuntime.RequestStop("nil") || noRuntime.StopRequested() || noRuntime.StopSignal() != nil {
		t.Error("AC#11 RED: the nil runtime answers like a runtime that can be stopped")
	}
}

// stopHookMount228 is one reading of the assembly hop the tray's Exit item rides.
type stopHookMount228 struct {
	installs  int
	hooks     int
	hookShape string // what the hook's argument was spelled, for the red line
}

// readStopHookMount228 walks runResident's body. Shape only, no line numbers: an
// unrelated edit above the assembly must not move this reading.
func readStopHookMount228(body *ast.BlockStmt) stopHookMount228 {
	out := stopHookMount228{}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if id, ok := call.Fun.(*ast.Ident); !ok || id.Name != "startResidentBall" {
			return true
		}
		out.installs++
		for _, arg := range call.Args {
			hop, ok := arg.(*ast.CallExpr)
			if !ok {
				continue
			}
			if hid, ok := hop.Fun.(*ast.Ident); !ok || hid.Name != "withStopRequest" {
				continue
			}
			out.hooks++
			if len(hop.Args) != 1 {
				out.hookShape = "a withStopRequest call with " + strconv.Itoa(len(hop.Args)) + " argument(s)"
				continue
			}
			inner, ok := hop.Args[0].(*ast.CallExpr)
			if !ok {
				if id, ok := hop.Args[0].(*ast.Ident); ok {
					out.hookShape = "the bare identifier " + id.Name
					continue
				}
				out.hookShape = "an argument that is not a call"
				continue
			}
			fid, ok := inner.Fun.(*ast.Ident)
			if !ok {
				out.hookShape = "a call whose function is not a name"
				continue
			}
			out.hookShape = fid.Name + "(rt)"
		}
		return true
	})
	return out
}

// TestAC11TrayExitExecutorIsMountedAtTheBallHost is the pin the dispatch asks for
// beside ⓑ-1: the tray branch has an executor, it is this runtime's stop bit, and
// the value driven here is the value production injects.
func TestAC11TrayExitExecutorIsMountedAtTheBallHost(t *testing.T) {
	rt := bootRuntime228(t)

	// (1) the production accessor, called exactly as resident_windows.go calls it.
	exec := residentTrayStop228(rt)
	if exec == nil {
		t.Fatal("AC#11 RED: residentTrayStop228 handed the ball host no executor")
	}
	recordTrayExit(exec)
	if !rt.StopRequested() {
		t.Fatal("AC#11 RED: recordTrayExit ran with the injected executor and the runtime still reports no stop request")
	}
	if got := rt.StopReason(); got != trayExitStopWhy {
		t.Errorf("AC#11 RED: the recorded reason is %q, want the tray's own sentence %q", got, trayExitStopWhy)
	}

	// (2) the mount: the host that builds Events.OnTrayExit has to be handed that
	// value by the assembly root. A walk, because the claim is about the call.
	body, file := runResidentBody258(t)
	if body == nil {
		t.Fatal("AC#11 RED (the instrument, not the code): no runResident found in resident_windows.go")
	}
	mount := readStopHookMount228(body)
	if mount.installs != 1 {
		t.Fatalf("AC#11 RED: runResident (%s) calls startResidentBall %d times, want 1 - this walk has lost its subject", file, mount.installs)
	}
	if mount.hooks != 1 {
		t.Fatalf("AC#11 RED: runResident (%s) hands the ball host %d withStopRequest hook(s), want 1. Zero is the "+
			"assembly root un-wiring the tray's Exit item; two is two owners of one gesture.", file, mount.hooks)
	}
	if mount.hookShape != "residentTrayStop228(rt)" {
		t.Fatalf("AC#11 RED: the stop hook installed at this line is %s, want residentTrayStop228(rt) - the hoisted "+
			"accessor is the one value a case may drive, and anything else here is a copy that reads green however "+
			"production changes (the ruling tickets 258 and 296 set for the [hotkey] chain)", mount.hookShape)
	}

	// (3) the nil branch, on the callback that is actually mounted: a host the
	// assembly root handed no executor to must report that without panicking on
	// the thread that fired it, and without touching anybody's stop bit.
	fresh := bootRuntime228(t)
	hs := panelHostHooks{}
	withStopRequest(residentTrayStop228(fresh))(&hs)
	hs.requestStop = nil // the pre-AC#11 bag: no executor handed over
	hs.trayExit()()
	if fresh.StopRequested() {
		t.Error("AC#11 RED: a bag with no stop executor recorded a stop request anyway")
	}
	// The same bag with the executor wired is the production shape, and it records.
	withStopRequest(residentTrayStop228(fresh))(&hs)
	hs.trayExit()()
	if !fresh.StopRequested() || fresh.StopReason() != trayExitStopWhy {
		t.Errorf("AC#11 RED: the mounted callback did not reach the stop bit (requested=%v reason=%q)",
			fresh.StopRequested(), fresh.StopReason())
	}
}

// TestAC11TrayExitTriggerWalksTheTenSteps is ⓑ-1: after the tray's trigger, the
// loop ends and the ten frozen D38(e) steps really run, in order, read inside the
// process that ran them.
func TestAC11TrayExitTriggerWalksTheTenSteps(t *testing.T) {
	rt := bootRuntime228(t)

	ran := map[int]int{}
	for _, step := range registeredSteps228 {
		s := step
		if err := rt.RegisterShutdownHook(s, func(context.Context) error {
			ran[int(s)]++
			return nil
		}); err != nil {
			t.Fatalf("RegisterShutdownHook(step %d): %v", int(s), err)
		}
	}

	// The trigger, and nothing between it and the loop but the loop's own check.
	trayTrigger228(rt)

	reason := runLoopBounded228(t, rt)
	if !strings.Contains(reason, "stop request") || !strings.Contains(reason, trayExitStopWhy) {
		t.Errorf("AC#11 RED: RunEventLoop returned %q, want a reason naming the stop request and who asked for it. "+
			"That reason is the only sentence this process writes about which trigger fired.", reason)
	}

	records := rt.Shutdown(false)
	if len(records) != len(frozenShutdownSteps228) {
		t.Fatalf("AC#11 RED: the audit trail holds %d records, want the frozen %d", len(records), len(frozenShutdownSteps228))
	}
	for i, rec := range records {
		if int(rec.Step) != i+1 {
			t.Errorf("AC#11 RED: record %d carries step number %d; the trail must count 1..10", i, int(rec.Step))
		}
		if rec.Name != frozenShutdownSteps228[i] {
			t.Errorf("AC#11 RED: record %d is named %q, want %q (the frozen D38(e) label)", i, rec.Name, frozenShutdownSteps228[i])
		}
		if rec.Err != nil {
			t.Errorf("AC#11 RED: step %d (%s) failed: %v", int(rec.Step), rec.Name, rec.Err)
		}
		if ran[int(rec.Step)] > 0 && rec.Skipped {
			t.Errorf("AC#11 RED: step %d (%s) is recorded skipped although this process registered an owner for it", int(rec.Step), rec.Name)
		}
	}
	for _, step := range registeredSteps228 {
		if ran[int(step)] != 1 {
			t.Errorf("AC#11 RED: the hook registered for step %d ran %d times, want exactly 1 - the trail says the "+
				"step happened, this counts says whether its owner executed it", int(step), ran[int(step)])
		}
	}
	t.Logf("AC#11 READING (ⓑ-1): reason=%q; trail=%s", reason, strings.Join(namesOf228(records), " > "))
}

// TestAC11StopRequestDuringTheBootWindowIsNotSwallowed is ⓐ's second clause. The
// ball exists before internal/proc installs its own signal handler, so a request
// recorded in that window has to be answered by the code in front of the loop -
// residentBootWindowExit228, the function runResident really calls. The console
// twin of this window is TestAC228ExitRequestDuringBootStillLeavesThroughD38E.
func TestAC11StopRequestDuringTheBootWindowIsNotSwallowed(t *testing.T) {
	// The clear case first: nothing asked, so the loop is entered. A check that
	// always answered "already asked" would pass the second half of this file and
	// break every boot.
	quiet := bootRuntime228(t)
	quietExit := make(chan os.Signal, 1)
	if why, left := residentBootWindowExit228(quiet, quietExit, quiet.StopSignal()); left {
		t.Fatalf("AC#11 RED: a boot nobody asked to stop was answered %q; the event loop would never be entered", why)
	}

	rt := bootRuntime228(t)
	bootExit := make(chan os.Signal, 1)
	bootStop := rt.StopSignal()
	if bootStop == nil {
		t.Fatal("AC#11 RED: proc.Boot built no stop bit, so this window has no tray door to read")
	}

	trayTrigger228(rt)

	why, left := residentBootWindowExit228(rt, bootExit, bootStop)
	if !left {
		t.Fatalf("AC#11 RED: a stop request recorded during boot was not seen by the window in front of the loop "+
			"(reason %q); the report printed before this line would say the loop is running, and RunEventLoop is "+
			"the only thing that would then honour it", why)
	}
	if !strings.Contains(why, "stop request") || !strings.Contains(why, trayExitStopWhy) {
		t.Errorf("AC#11 RED: the boot-window answer is %q, want it to name a stop request and the tray item that "+
			"asked for it", why)
	}
	if !strings.Contains(why, "the event loop was never entered") {
		t.Errorf("AC#11 RED: the boot-window answer %q does not say the loop was skipped, which is the sentence "+
			"AC#7 exists to keep true", why)
	}

	// And the request is still on record for whoever asks next: this window reads
	// a latch, it does not consume it. Entering RunEventLoop after a check that
	// said "leave" is the same few instructions away, and the loop must still see
	// the same answer.
	if !rt.StopRequested() {
		t.Error("AC#11 RED: reading the boot window used the request up")
	}
	if why2, left2 := residentBootWindowExit228(rt, bootExit, bootStop); !left2 {
		t.Errorf("AC#11 RED: the same request answered the check only once (second reading %q)", why2)
	}
	t.Logf("AC#11 READING (ⓐ gap): %q", why)
}
