//go:build windows

package main

// Ticket 246 AC#7, default tier: the shapes that do NOT need a desktop.
//
// AC#7's measured gap was a caller count - AskOnTaskRoot and askConfirmation had
// zero product callers, so the resident gate had no source of cards. The first
// four legs closed every other half of that ticket and left exactly that hole,
// which is why the cases here are wired to "who asks, through what, and what this
// process says when it cannot":
//
//   - the pipeline this leg now assembles must use the ONE gate the resident leg
//     built at boot (ruling 2.2), which is measured by pointer identity and by the
//     resident surface's own sentinel error arriving back through the bridge - not
//     by a string this file typed;
//   - the task's context must hang under the root D38(e) step 3 cancels (ruling
//     2.3), measured by what the provider ever sees: a run whose root was cancelled
//     issues zero chat requests, and the same call on a live root issues at least
//     one (that pair is this ruler's positive control);
//   - the console gate is interactiveStdin() and nothing else (ruling 2.1), so a
//     shipped resident process with no console says its task entry is off AND does
//     not open the pipeline at all - the absence of a wisp.db is the reading that
//     says the refusal was the real branch, not a print;
//   - the test-only injection is narrow on purpose (ruling 2.4): accepted in test
//     with the harness' own data root, refused everywhere else, and a refusal that
//     is printed rather than dropped.
//
// What this file does NOT claim, and cannot: the four steps AC#7 asks for with a
// card that is really on an orb and a cancel key that is really borrowed live in
// resident_task_source_live_246_windows_test.go behind -tags winlive, because only
// Win32 can answer those. Nothing here says "任务管线已跑通"; the split is in the
// evidence table (§2), and the real-credential reading is booked as owed.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/proc"
	"golang.org/x/sys/windows"
)

// TestAC246ResidentPipelineAsksThroughTheOneGate is ruling 2.2's capability half.
// A fixture that can assemble (config + credential + provider) is handed to
// assembleRuntime with ra's three pieces injected, and the questions are:
//
//  1. is the gate the bridge asks through the SAME object the resident leg built
//     at boot (pointer identity - a second gate would pass every string assertion
//     in this repository and still be the bug);
//  2. did the injected console surface get built (it must NOT: a resident card is
//     not a console card);
//  3. when a real call reaches a real window and this process has no ball to put
//     it on, does the refusal come back carrying the resident surface's own
//     sentinel error - which is the only way to know the gate's UI really is
//     ra.ui, from outside the package.
func TestAC246ResidentPipelineAsksThroughTheOneGate(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	ra := newResidentApproval()
	// No bindBallHost: the no-desktop / no-window posture is the one this case is
	// about, and it is the fail-closed direction (a card nobody can be shown must
	// never run the call).
	ar, code := assembleRuntime(runSpec{
		stdout:  f.out,
		stderr:  f.err,
		dataDir: f.dir,
		gate:    ra.gate,
		ui:      ra.ui,
		cards:   ra.cards,
		taskCtx: ra.root,
		notify:  swallowPoster,
	})
	if code != 0 {
		legFail(t, ra, f, "assemble with the resident gate injected: exit %d", code)
	}
	defer ar.close()

	if ar.gate != ra.gate {
		t.Fatalf("the assembled pipeline runs on a DIFFERENT gate than the resident leg built: %p vs %p "+
			"(two gates in one process is the shape ruling 2.2 and ledger A481 refuse)", ar.gate, ra.gate)
	}
	if ar.ui != nil {
		t.Error("assembleRuntime built a console approval surface next to an injected gate: a resident card " +
			"would be printed to a terminal nobody is watching while the orb stayed asleep")
	}
	if ar.spec.taskCtx != ra.root {
		t.Error("runSpec.taskCtx did not reach the runtime, so execute() would park the task off the root " +
			"D38(e) step 3 cancels")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// A real L1 fs.write through the bridge this assembly built. The task is
	// admitted the way the loop admits it (D47), so the only thing that can stop
	// this call is the gate - and the gate's UI is the resident surface.
	target := filepath.Join(f.dir, "resident-pipeline-write.txt")
	revoke := ar.gate.AdmitTextTask("host:246-one-gate")
	defer revoke()
	out, err := ar.bridge.Execute(ctx, agent.ToolRequest{
		TaskID: "host:246-one-gate", CorrelationID: "host-corr-1", CallID: "w1",
		Name: "fs.write",
		Args: json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"must not land"}`, filepath.ToSlash(target))),
	})
	if err != nil {
		t.Fatalf("bridge.Execute on the assembled stack: %v", err)
	}
	if !out.IsError {
		t.Fatalf("a resident leg with no ball executed a write it had nowhere to confirm: %+v", out)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Errorf("the vetoed-by-absence write landed on disk anyway: %v", statErr)
	}
	if !strings.Contains(out.Text, errResidentNoBall.Error()) {
		t.Fatalf("the refusal did not come from the resident surface: text = %q, want it to carry %q "+
			"(that sentinel is returned only by ballCardUI.Prompt, so its absence means the gate is not "+
			"showing cards on the surface this process injected)", out.Text, errResidentNoBall)
	}
	if n := ra.ui.displayedCards(); n != 0 {
		t.Errorf("the resident surface counted %d displayed cards with no window attached", n)
	}
	if n := ar.windowCount(); n != 0 {
		t.Errorf("windowCount() = %d for a run whose surface displayed nothing (the injected-UI fallback "+
			"reads the host surface's own counter)", n)
	}
	if got := ra.vetoByEsc(); !strings.Contains(got, "没有可否决的确认项") {
		t.Errorf("vetoByEsc with a refused card said %q", got)
	}
}

// TestAC246ResidentGateInjectionIsRefusedHalfAssembled pins the guard inside
// assembleRuntime: a gate handed in without the surface and the ledger it was
// built with is a second truth source waiting to happen, and the assembly says so
// instead of quietly composing a console surface next to it.
func TestAC246ResidentGateInjectionIsRefusedHalfAssembled(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	ra := newResidentApproval()
	cases := []struct {
		name string
		spec runSpec
	}{
		{"gate without ledger", runSpec{stdout: f.out, stderr: f.err, dataDir: f.dir, gate: ra.gate, ui: ra.ui}},
		{"gate without surface", runSpec{stdout: f.out, stderr: f.err, dataDir: f.dir, gate: ra.gate, cards: ra.cards}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ar, code := assembleRuntime(c.spec)
			if code == 0 {
				defer ar.close()
				t.Fatalf("assembleRuntime accepted %s: a half-injected gate is how a process ends up with two "+
					"approval truths", c.name)
			}
			if ar != nil {
				ar.close()
			}
			if !strings.Contains(f.err.String(), "注入的审批门必须连同它自己的界面与会话账本一起递进来") {
				t.Errorf("%s: exit %d but the refusal was never said: stderr=%q", c.name, code, f.err.String())
			}
			if ra.gate.Channels().Loaded(approval.ChannelEsc) {
				t.Errorf("%s: a refused assembly still loaded a channel", c.name)
			}
		})
	}
}

// TestAC246ResidentTaskRootCancelStopsTheModelCallIs ruling 2.3 measured as a
// capability rather than as a comment: when the resident leg's task root is
// cancelled, a task started through the assembled pipeline must not reach the
// model at all. The pair below is the ruler's positive control - the SAME call on
// a live root does reach it - so a chat counter that is stuck at zero for an
// unrelated reason cannot buy the pass.
func TestAC246ResidentTaskRootCancelStopsTheModelCall(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	ra := newResidentApproval()
	ar, code := assembleRuntime(runSpec{
		stdout: f.out, stderr: f.err, dataDir: f.dir,
		gate: ra.gate, ui: ra.ui, cards: ra.cards, taskCtx: ra.root,
		notify: swallowPoster,
	})
	if code != 0 {
		legFail(t, ra, f, "assemble for the root-cancel reading: exit %d", code)
	}
	defer ar.close()

	// (control) a live root: the loop talks to the provider.
	liveCode := ar.execute(context.Background(), "总结一下 这份笔记")
	liveReqs := f.srv.RouteCount(t, "chat")
	t.Logf("AC#7 positive control: live root -> execute code %d, provider chat requests %d", liveCode, liveReqs)
	if liveReqs == 0 {
		legFail(t, ra, f, "the live-root run reached the provider 0 times, so the reading below proves nothing"+
			" (exit %d)", liveCode)
	}

	// D38(e) step 3 as the resident leg registers it: cancel the task root.
	hookCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ra.cancelTaskRoots(hookCtx); err != nil {
		t.Fatalf("cancelTaskRoots with nothing hanging: %v", err)
	}
	if ra.root.Err() == nil {
		t.Fatal("cancelTaskRoots returned without cancelling the root it owns")
	}

	before := f.srv.RouteCount(t, "chat")
	code2 := ar.execute(ra.root, "总结一下 这份笔记")
	after := f.srv.RouteCount(t, "chat")
	if after != before {
		t.Errorf("a cancelled task root still issued %d provider requests, want 0: the loop is not parked on "+
			"the root D38(e) step 3 cancels", after-before)
	}
	if code2 == 0 {
		t.Errorf("a task cancelled before it started exited 0; SPEC-05 §3.4 wants a classified failure")
	}
	t.Logf("AC#7 READING (ruling 2.3): after step 3's cancel, execute -> exit %d, provider requests %d -> %d",
		code2, before, after)
}

// TestAC246TestTaskInjectionPredicate walks the escape hatch's own predicate.
// residentTestTaskText is the whole rule (ruling 2.4: test env AND the harness'
// own data root), and every rejected branch owes a printed reason - the case
// below the one that shows the shipped process really calls this function.
func TestAC246TestTaskInjectionPredicate(t *testing.T) {
	const (
		harnessDir = `C:\temp\wisp-harness-dir`
		otherDir   = `C:\temp\someone-elses-root`
	)
	layout := func(dir string) proc.Layout { return proc.Layout{DataDir: dir} }

	t.Run("accepted in the harness' own root", func(t *testing.T) {
		t.Setenv(testTaskTextEnv, "把说明写进 note.txt")
		t.Setenv(proc.TestDataDirEnv, harnessDir)
		text, refusal := residentTestTaskText(buildinfo.EnvTest, layout(harnessDir))
		if text == "" || refusal != "" {
			t.Fatalf("test env + harness-owned root must be accepted: text=%q refusal=%q", text, refusal)
		}
	})
	t.Run("refused outside test", func(t *testing.T) {
		t.Setenv(testTaskTextEnv, "这条在生产里不该被吃")
		t.Setenv(proc.TestDataDirEnv, harnessDir)
		for _, env := range []buildinfo.Env{buildinfo.EnvDev, buildinfo.EnvProd} {
			text, refusal := residentTestTaskText(env, layout(harnessDir))
			if text != "" {
				t.Errorf("env %s accepted a task-text injection", env)
			}
			if !strings.Contains(refusal, taskEntryRefusedClaim) {
				t.Errorf("env %s refused silently: %q", env, refusal)
			}
		}
	})
	t.Run("refused when the root is not the harness' own", func(t *testing.T) {
		t.Setenv(testTaskTextEnv, "这条也不该被吃")
		t.Setenv(proc.TestDataDirEnv, "")
		if _, refusal := residentTestTaskText(buildinfo.EnvTest, layout(harnessDir)); !strings.Contains(refusal, proc.TestDataDirEnv) {
			t.Errorf("missing %s produced no named refusal: %q", proc.TestDataDirEnv, refusal)
		}
		t.Setenv(proc.TestDataDirEnv, otherDir)
		if _, refusal := residentTestTaskText(buildinfo.EnvTest, layout(harnessDir)); !strings.Contains(refusal, taskEntryRefusedClaim) {
			t.Errorf("a data root that is not the injected one was not refused: %q", refusal)
		}
	})
	t.Run("unset is silent", func(t *testing.T) {
		t.Setenv(testTaskTextEnv, "")
		text, refusal := residentTestTaskText(buildinfo.EnvTest, layout(harnessDir))
		if text != "" || refusal != "" {
			t.Errorf("an unset injection said something: %q / %q", text, refusal)
		}
	})
	t.Run("a refusal never carries the whole text", func(t *testing.T) {
		long := strings.Repeat("很", 40)
		t.Setenv(testTaskTextEnv, long)
		t.Setenv(proc.TestDataDirEnv, "")
		_, refusal := residentTestTaskText(buildinfo.EnvDev, layout(harnessDir))
		if strings.Contains(refusal, long) {
			t.Error("the refusal line printed the entire task text instead of a bounded echo")
		}
	})
}

// TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry is the gate read
// where it matters - the command line the owner double clicks. A child with no
// interactive console must (a) say the entry is off, (b) NOT open the pipeline
// (proved by the data root holding no wisp.db: the branch returned before
// assembleRuntime), (c) still host the ball, and (d) still leave through the
// D38(e) order with nothing failed.
func TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	sinkDir := logSinkDir(dataDir)

	sawDisabled := pollUntil127(200, func() bool { return leg.stdout.has(taskEntryDisabledClaim) })
	sawLoop := pollUntil127(200, func() bool { return leg.stdout.has(residentReachedLoop) })
	out := leg.stdout.String()
	if !sawDisabled || !sawLoop {
		leg.stop()
		t.Fatalf("AC#7 RED: the shipped leg with no console said %v that its entry is off and %v that it "+
			"reached the loop.\n%s\nOne of those two sentences is the downgrade made visible; the other is the "+
			"proof the ball still came up after it.", sawDisabled, sawLoop, leg.console())
	}
	if strings.Contains(out, taskEntryInjectedClaim) || strings.Contains(out, taskEntryConsoleClaim) {
		t.Errorf("AC#7 RED: the leg claims an enabled task entry with stdin disconnected:\n%s", out)
	}
	if !strings.Contains(out, "任务来源："+taskPostureAbsent) {
		t.Errorf("AC#7 RED: the boot report does not name the absent task source (%q):\n%s",
			"任务来源："+taskPostureAbsent, out)
	}
	// (b) the reading that the refusal is a BRANCH and not a print: no store was
	// ever opened, because the entry returned before the assembly ran.
	if _, err := os.Stat(filepath.Join(dataDir, "wisp.db")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("AC#7 RED: the leg opened a data store with no task source (stat err %v): the disabled "+
			"branch still assembled a pipeline, which means the console gate is decoration", err)
	}
	if strings.Contains(out, pipelineAbsentClaim) {
		t.Errorf("AC#7 RED: the leg printed %q when it never tried to assemble: the two degraded branches "+
			"collapsed into one sentence:\n%s", pipelineAbsentClaim, out)
	}

	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("GenerateConsoleCtrlEvent on the leg's own group: %v", err)
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("AC#7 RED: the leg never left through its own shutdown path.\n%s", leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("AC#7 RED: exit %v on a leg with no task running.\n%s", leg.waitErr, leg.console())
	}
	recs := readResidentSink(t, sinkDir)
	if i := indexOfMsgContaining246(recs, cancelStepBookedMsg); i < 0 {
		t.Errorf("AC#7 RED: the on-disk trail lost the step-3 record (%v)", msgsOf127(recs))
	}
	for _, bad := range []string{"shutdown step failed", "abandoned wait"} {
		if idx := indexOfMsgContaining246(recs, bad); idx >= 0 {
			t.Errorf("AC#7 RED: the shipped process booked %q: %q", bad, msgsOf127(recs)[idx])
		}
	}
}

// TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline is the
// positive control of the case above: the SAME harness, one env var apart. Here
// the injection is accepted, so the leg really calls assembleRuntime - and in a
// data root with no config.toml that attempt fails loudly and the boot still
// reaches its loop. Two readings, both required: the injection claim, and the
// pipeline-refused claim, which together can only be printed by a process that
// actually tried.
func TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLegWithEnv(t, exe, []string{
		"WISP_ENV=test",
		procTestDataDirEnv + "=" + dataDir,
		testTaskTextEnv + "=把说明写进 note.txt",
	})
	sinkDir := logSinkDir(dataDir)

	sawInjected := pollUntil127(200, func() bool { return leg.stdout.has(taskEntryInjectedClaim) })
	sawPipeline := pollUntil127(200, func() bool { return leg.stdout.has(pipelineAbsentClaim) })
	sawLoop := pollUntil127(200, func() bool { return leg.stdout.has(residentReachedLoop) })
	if !sawInjected || !sawPipeline || !sawLoop {
		leg.stop()
		t.Fatalf("AC#7 RED: injection said=%v pipeline=%v loop=%v.\n%s\n"+
			"Those three sentences are the whole of \"the entry was accepted and the pipeline really was "+
			"attempted and refused\"; any one of them missing means the branch is not what it reports.",
			sawInjected, sawPipeline, sawLoop, leg.console())
	}
	out := leg.stdout.String()
	if strings.Contains(out, taskEntryDisabledClaim) {
		t.Errorf("AC#7 RED: an accepted injection still printed the no-console refusal:\n%s", out)
	}
	if strings.Contains(out, taskEntryRefusedClaim) {
		t.Errorf("AC#7 RED: the predicate refused its own harness root:\n%s", out)
	}
	if !strings.Contains(out, "任务来源："+taskPostureRefused) {
		t.Errorf("AC#7 RED: the boot report does not separate \"pipeline refused\" from \"no entry\" (%q):\n%s",
			"任务来源："+taskPostureRefused, out)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "config.toml")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("AC#7 RED: the leg created a config.toml it was never asked for (err %v)", err)
	}

	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("break: %v", err)
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("AC#7 RED: the leg did not leave through its own path.\n%s", leg.console())
	}
	recs := readResidentSink(t, sinkDir)
	for _, bad := range []string{"shutdown step failed", "abandoned wait"} {
		if idx := indexOfMsgContaining246(recs, bad); idx >= 0 {
			t.Errorf("AC#7 RED: booked %q: %q", bad, msgsOf127(recs)[idx])
		}
	}
}

// TestAC246DevLegIgnoresTheTestTaskInjection is the escape hatch's other edge, on
// the shipped binary: WISP_ENV=dev with APPDATA pointed at a throwaway root (so
// the leg resolves a data dir without touching the owner's), the injection set,
// and no console. The process must refuse the injection by name, fall through to
// the no-console branch, and host its ball as usual.
func TestAC246DevLegIgnoresTheTestTaskInjection(t *testing.T) {
	exe := buildWispForTest(t)
	appData := t.TempDir()
	dataDir := filepath.Join(appData, "wisp-dev")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		t.Fatal(err)
	}
	leg := bootResidentLegWithEnv(t, exe, []string{
		"WISP_ENV=dev",
		"APPDATA=" + appData,
		"LOCALAPPDATA=" + appData,
		testTaskTextEnv + "=这条 env 在生产里不该被受理",
	})

	sawRefusal := pollUntil127(200, func() bool { return leg.stdout.has(taskEntryRefusedClaim) })
	sawDisabled := pollUntil127(200, func() bool { return leg.stdout.has(taskEntryDisabledClaim) })
	sawLoop := pollUntil127(200, func() bool { return leg.stdout.has(residentReachedLoop) })
	out := leg.stdout.String()
	if strings.Contains(out, "another instance is running") {
		leg.stop()
		t.Fatalf("AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session, so "+
			"this leg handed its activation over and exited. Stop it and re-run.\n%s", leg.console())
	}
	if !sawRefusal || !sawDisabled || !sawLoop {
		leg.stop()
		t.Fatalf("AC#7 RED: dev leg: refusal=%v entry-off=%v loop=%v.\n%s\nA production-shaped environment that "+
			"eats a test injection is the failure ruling 2.4 names; a refusal that is not printed is the same "+
			"failure wearing silence.", sawRefusal, sawDisabled, sawLoop, leg.console())
	}
	if strings.Contains(out, taskEntryInjectedClaim) {
		t.Errorf("AC#7 RED: the dev leg opened the task entry through the test injection:\n%s", out)
	}
	leg.stop()
}

// ------------------------------------------------------------------ helpers

// bootResidentLegWithEnv is bootResidentLeg with the child's environment spelled
// out, because three of AC#7's readings are about which ENVIRONMENT the shipped
// process was started in. The stdin stays nil: every one of these cases is the
// no-interactive-console posture.
func bootResidentLegWithEnv(t *testing.T, exe string, env []string) *residentLeg {
	t.Helper()
	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	cmd.Env = append(os.Environ(), env...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
	leg := &residentLeg{cmd: cmd, stdout: &lockedBuf{}, stderr: &lockedBuf{}}
	cmd.Stdout, cmd.Stderr, cmd.Stdin = leg.stdout, leg.stderr, nil
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the resident leg (%s, env %v): %v", exe, env, err)
	}
	leg.done = make(chan error, 1)
	go func() { leg.done <- cmd.Wait() }()
	t.Cleanup(leg.stop)
	return leg
}

// swallowPoster keeps an in-process AC#7 case from posting real toasts while it
// drives a task through execute(): D10's notification is the run leg's result
// path, and a test that runs it ten times a day is not a reason to fill the
// desktop. It is a recorder-shaped hole in the injection surface, not a stub
// standing in for a subsystem - the poster itself is pinned in notify_windows.go's
// own leg.
func swallowPoster(string, string) error { return nil }

// legFail prints both consoles and then fails, so the reading a case could not
// take is in the log rather than described after the fact.
func legFail(t *testing.T, ra *residentApproval, f *runFixture, format string, args ...any) {
	t.Helper()
	t.Fatalf(format+"\nstatus line: %s\nstdout:\n%s\nstderr:\n%s",
		append(args, ra.residentStatusLine(), f.out.String(), f.err.String())...)
}
