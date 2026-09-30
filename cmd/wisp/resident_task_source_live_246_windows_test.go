//go:build windows && winlive

package main

// Ticket 246 AC#7, live tier: the four steps on a real desktop, in the process
// the owner actually launches.
//
//	起任务  -> a shipped wisp.exe with no arguments takes ONE task text from its
//	           task entry and runs a real agent loop against a real HTTP endpoint;
//	举卡    -> the model's own tool call reaches the bridge, the bridge asks the
//	           ONE gate this process owns, and the resident surface puts it on the
//	           orb (D43's Confirming) while Win32 borrows the cancel key;
//	Esc 否决 -> a SECOND PROCESS on the desktop injects the bare Esc and Wisp's
//	           hot key fires: the card answers veto, the gate books
//	           ANSWER-VETO channel=esc, and the written file never lands;
//	任务被取消 -> measured two ways, because D43's row 22 and D38(e)'s step 3 are
//	           two different sentences: the vetoed CALL is cancelled and the loop
//	           relays it back to the model (TestLive246ResidentPipeline...), and a
//	           RUNNING TASK is cancelled by the exit path
//	           (TestLive246ExitCancelsARunningTask).
//
// WHAT THIS IS AND IS NOT ABOUT CREDENTIALS. The provider here is a local
// httptest endpoint replaying the repository's own golden SSE format through
// internal/llm/golden - the C5 seam AGENTS.md §1.3 names - reached over real HTTP
// by a real openai-chat adapter with a real DPAPI blob. It is NOT a real model:
// the reading that needs an owner-entered credential is owed and is booked as
// owed in docs/evidence/s1/246-resident-task-source-r2.md §2. Nothing in this
// file may be read as "the task pipeline has been proven against a live model",
// and no key material appears in it: the blob value is this package's existing
// fakeStoreKey constant, named by identifier, never printed.
//
// The task text arrives through WISP_TEST_TASK_TEXT, the narrow injection ruling
// 2.4 asked for. The console path cannot be driven from a subprocess at all (its
// stdin is a pipe, and interactiveStdin() correctly says so), which is the whole
// reason the injection exists; TestAC246DevLegIgnoresTheTestTaskInjection in the
// default tier is the ruler that keeps it narrow.
//
// Run it with no other Wisp or balldebug process on the desktop:
//
//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -tags winlive ./cmd/wisp \
//	  -run 'TestLive246ResidentPipeline|TestLive246ExitCancels' -v

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/secret"
)

// cardRaisedClaim is the resident surface's own "a card is on the orb" line
// (ballCardUI.Prompt). The veto case's whole budget is the 3s L1 window, so this
// line is what the case waits for rather than a sleep.
const cardRaisedClaim = "wisp: 卡片挂起："

// vetoDoneClaim is the sentence the injected cancel executor returns and the ball
// host prints (residentApproval.vetoByEsc through recordCancelHotkey).
const vetoDoneClaim = "按 Esc 否决了卡片"

// residentCardFileName is the file the scripted model asks to write. Its being
// absent after the veto is the reading that the vetoed call never executed.
const residentCardFileName = "resident-246r2-card.txt"

// TestLive246ResidentPipelineRaisesACardAndEscVetoesIt is AC#7's four steps on a
// shipped resident process: task -> card -> Esc veto -> the call cancelled.
func TestLive246ResidentPipelineRaisesACardAndEscVetoesIt(t *testing.T) {
	rig := buildEscListener246(t)

	// The ruler's own positive control, before any of our code runs: a window
	// that registers bare Esc itself must NOT receive the keydown.
	ctrl := observeEsc246(t, rig, "-steal")
	if ctrl.keydown != 0 || ctrl.hotkey != 1 {
		t.Fatalf("the observer rig is blind: with its OWN global bare-Esc registration it still received "+
			"keydown_esc=%d (want 0), wm_hotkey=%d (want 1); no reading below means anything",
			ctrl.keydown, ctrl.hotkey)
	}
	base := observeEsc246(t, rig, "-watch")
	if base.keydown != 1 {
		t.Fatalf("AC#7 LIVE RED (desktop state, not our code): an injected bare Esc reached no foreground "+
			"window before Wisp started (keydown_esc=%d, wm_hotkey=%d): somebody else owns the key on this "+
			"desktop, so the borrow cannot be measured here", base.keydown, base.hotkey)
	}

	dataDir := t.TempDir()
	target := filepath.Join(dataDir, residentCardFileName)
	startResidentHarness(t, dataDir, residentCardGolden(t, target))

	// Arm the observer while nothing is pending: it takes the foreground, pumps
	// its window and waits for this case to say "go" on stdin. Triggering from
	// the test rather than from a stop watch is what keeps a 2-3s L1 window from
	// racing the rig's own startup.
	arm := armEscObserver246r2(t, rig)

	exe := buildWispForTest(t)
	leg := bootResidentLegWithEnv(t, exe, []string{
		"WISP_ENV=test",
		procTestDataDirEnv + "=" + dataDir,
		testTaskTextEnv + "=把这段说明写进 note.txt",
	})
	sinkDir := logSinkDir(dataDir)

	// Steps 1 and 2: the entry took the task, the loop ran, and the model's own
	// tool call reached the gate. The card line is the only thing that can
	// produce it, and the case says which of the two readings failed.
	sawCard := pollUntil127(400, func() bool { return leg.stdout.has(cardRaisedClaim) })
	out := leg.stdout.String()
	if !sawCard {
		leg.stop()
		t.Fatalf("AC#7 LIVE RED (step 1 起任务 / step 2 举卡): the shipped resident process never raised a "+
			"card from its own task source.\n%s\nThat is the exact gap AC#7 was filed for - the gate, the orb "+
			"and the veto route were all landed, and nothing asked.", leg.console())
	}
	if !strings.Contains(out, taskEntryInjectedClaim) {
		t.Errorf("AC#7 LIVE RED: the task came from the injection but the boot never said so:\n%s", out)
	}
	t.Logf("AC#7 LIVE steps 1-2: %s", tailLine(out, cardRaisedClaim))

	// Step 3: the veto. One injected Esc, on command, while the card counts down.
	arm.trigger()
	sawVeto := pollUntil127(200, func() bool { return leg.stdout.has(vetoDoneClaim) })
	reading := arm.reading(t)
	if !sawVeto {
		leg.stop()
		t.Fatalf("AC#7 LIVE RED (step 3 Esc 否决): the injected Esc arrived while a card was counting down and "+
			"nothing vetoed it. observer: %s\n%s", reading.raw, leg.console())
	}
	if reading.keydown != 0 {
		t.Fatalf("AC#7 LIVE RED: the second process received the injected Esc (keydown_esc=%d) while a card was "+
			"waiting, so the key was not borrowed even though the veto fired: %s", reading.keydown, reading.raw)
	}

	// The vetoed call never executed - the reading the whole ticket is about.
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("AC#7 LIVE RED: the vetoed fs.write landed anyway (stat err %v): the veto did not reach the "+
			"tool call", err)
	}

	// Step 4, first shape: D43 row 22 - the vetoed CALL is cancelled, the loop
	// relays it back to the model, and the task then closes.
	sawEnd := pollUntil127(400, func() bool {
		o := leg.stdout.String()
		return strings.Contains(o, "wisp run: 任务") && strings.Contains(o, "结束")
	})
	if !sawEnd {
		leg.stop()
		t.Fatalf("AC#7 LIVE RED (step 4): the task never closed after the veto; D43 row 22 says the vetoed call "+
			"goes back to the LLM, so the loop must still finish the task.\n%s", leg.console())
	}
	t.Logf("AC#7 LIVE step 4 (D43 row 22): %s", tailLine(leg.stdout.String(), "wisp run: 任务"))

	// The gate's own record, read off the shipped process's on-disk ledger.
	recs := readResidentSink(t, sinkDir)
	if i := indexOfMsgContaining246(recs, "ANSWER-VETO"); i < 0 {
		t.Fatalf("AC#7 LIVE RED: no ANSWER-VETO record in the shipped process's log (of %d records: %v)",
			len(recs), msgsOf127(recs))
	} else if !strings.Contains(recs[i].Msg, "channel=esc") {
		t.Errorf("AC#7 LIVE RED: the veto record does not name the esc channel: %q", recs[i].Msg)
	} else {
		t.Logf("AC#7 LIVE veto record: %s", recs[i].Msg)
	}
	card := cardLedgerRecord(t, sinkDir)
	if card == nil {
		t.Errorf("AC#7 LIVE RED: the shipped process's ledger holds no card record: %v", msgsOf127(recs))
	} else if card["esc_borrowed"] != true {
		// The attribute, not the sentence: ballCardUI.Prompt books esc_borrowed as
		// a structured field of the card record. The desktop-level proof is the
		// observer's keydown=0 above, which no log line can fake; this is the
		// ledger's own statement of what the surface asked Win32 for.
		t.Errorf("AC#7 LIVE RED: the card record carries esc_borrowed=%v, want true: %v", card["esc_borrowed"], card)
	} else if card["orb_state"] != "Confirming" {
		t.Errorf("AC#7 LIVE RED: the card put the orb in %v, want D43's Confirming: %v", card["orb_state"], card)
	}

	// The two exit steps the task source owns have to be in the boot report's
	// roster, taken off the registration rather than off a comment.
	boot := leg.stdout.String()
	for _, claim := range []string{"1:scheduler-close", "7:flush-logs-close-db"} {
		if !strings.Contains(boot, claim) {
			t.Errorf("AC#7 LIVE RED: the boot report does not own %s (the task source registered nothing): %q",
				claim, tailLine(boot, "D38(e) steps"))
		}
	}

	// (c) RETURNED: with nothing waiting any more, the very next injected Esc
	// reaches the second process again.
	returned := observeEsc246(t, rig, "-watch")
	if returned.keydown != 1 {
		t.Fatalf("AC#7 LIVE RED: after the veto the second process received keydown_esc=%d, want 1 - Wisp is "+
			"still swallowing Esc (wm_hotkey=%d on the observer means the observer itself holds it): %s",
			returned.keydown, returned.hotkey, returned.raw)
	}

	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("break: %v", err)
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("AC#7 LIVE RED: the leg with a finished task did not leave through its own path.\n%s", leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("AC#7 LIVE RED: exit %v on a clean leg.\n%s", leg.waitErr, leg.console())
	}
	recs = readResidentSink(t, sinkDir)
	for _, bad := range []string{"shutdown step failed", "abandoned wait", "到点仍有"} {
		if idx := indexOfMsgContaining246(recs, bad); idx >= 0 {
			t.Errorf("AC#7 LIVE RED: the shipped process booked %q: %q", bad, msgsOf127(recs)[idx])
		}
	}
	if !strings.Contains(leg.stdout.String(), "10 steps, 0 failed") {
		t.Errorf("AC#7 LIVE RED: no clean 10-step exit: %q", tailLine(leg.stdout.String(), "shutdown order"))
	}
	t.Logf("AC#7 LIVE: observer keydown %d->%d->%d (steal control %d)",
		base.keydown, reading.keydown, returned.keydown, ctrl.keydown)
}

// TestLive246ExitCancelsARunningTask is AC#7's fourth step read the way D38(e)
// step 3 and ruling 2.3 mean it: an exit request while a task is RUNNING cancels
// that task through the root the gate owns, and the cancelled status is booked on
// the console and in the store ahead of the sink closing.
//
// The scripted turn streams slowly on purpose (the golden format's own @latency
// directive through internal/llm/golden's pacing replayer), because the reading
// needs a task genuinely mid-flight when the break arrives. A task that finished
// first would answer "nothing to cancel", which is the false green this case
// exists to refuse.
func TestLive246ExitCancelsARunningTask(t *testing.T) {
	dataDir := t.TempDir()
	srv := startGoldenServer(t, residentSlowGolden(t))
	prepareResidentHarness(t, dataDir, srv.URL)

	exe := buildWispForTest(t)
	leg := bootResidentLegWithEnv(t, exe, []string{
		"WISP_ENV=test",
		procTestDataDirEnv + "=" + dataDir,
		testTaskTextEnv + "=慢慢说一段话，别停",
	})
	sinkDir := logSinkDir(dataDir)

	// The task is running once the streamed text starts arriving.
	sawStream := pollUntil127(400, func() bool { return leg.stdout.has("第一段 0") })
	if !sawStream {
		leg.stop()
		t.Fatalf("AC#7 LIVE RED: no streamed text ever arrived, so there was never a running task to cancel.\n%s",
			leg.console())
	}
	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("break: %v", err)
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("AC#7 LIVE RED: the leg took the exit request and stayed.\n%s", leg.console())
	}
	out := leg.stdout.String()
	if !strings.Contains(out, "结束（cancelled") {
		t.Fatalf("AC#7 LIVE RED: the running task was not reported as cancelled by the exit path.\n%s\n"+
			"That is ruling 2.3's whole point: the task loop's ctx hangs under the root D38(e) step 3 cancels, "+
			"so leaving must end the task as cancelled and say so.", leg.console())
	}
	recs := readResidentSink(t, sinkDir)
	if i := indexOfMsgContaining246(recs, cancelStepBookedMsg); i < 0 {
		t.Errorf("AC#7 LIVE RED: step 3 left no completion record: %v", msgsOf127(recs))
	}
	for _, bad := range []string{"shutdown step failed", "abandoned wait", "到点仍有"} {
		if idx := indexOfMsgContaining246(recs, bad); idx >= 0 {
			t.Errorf("AC#7 LIVE RED: booked %q: %q", bad, msgsOf127(recs)[idx])
		}
	}
	if i := indexOfMsgContaining246(recs, "退出第 7 步完成"); i < 0 {
		t.Errorf("AC#7 LIVE RED: the task source's step 7 never booked its close: %v", msgsOf127(recs))
	}
	t.Logf("AC#7 LIVE (exit cancels a running task): %s", tailLine(out, "wisp run: 任务"))
}

// ------------------------------------------------------------------ harness

// startGoldenServer serves golden-format bytes over a real HTTP listener through
// internal/llm/golden's own replayer (the unit-test runner of the one format).
func startGoldenServer(t *testing.T, goldenText string) *httptest.Server {
	t.Helper()
	srv, rep, err := golden.ServerFromBytes([]byte(goldenText))
	if err != nil {
		t.Fatalf("golden.ServerFromBytes: %v", err)
	}
	rep.Pace = true
	t.Cleanup(srv.Close)
	return srv
}

// prepareResidentHarness writes what a shipped resident process needs to actually
// assemble: config.toml naming the local endpoint, and the DPAPI blob that the
// api_key_ref resolves through. The key value is this package's existing fake
// constant and is never printed by anything here.
func prepareResidentHarness(t *testing.T, dataDir, baseURL string) {
	t.Helper()
	cfg := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]
timeout_ms = 120000

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
context_window = 128000

[fs]
allowed_dirs = [%q]
`, baseURL+"/v1", filepath.ToSlash(dataDir))
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(dataDir)
	if err != nil {
		t.Fatalf("secret store: %v", err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatalf("store the fake blob: %v", err)
	}
}

// startResidentHarness is both halves in the order the veto case needs.
func startResidentHarness(t *testing.T, dataDir, goldenText string) *httptest.Server {
	t.Helper()
	srv := startGoldenServer(t, goldenText)
	prepareResidentHarness(t, dataDir, srv.URL)
	return srv
}

// residentCardGolden is the two-turn script the veto case runs: the model asks
// for one fs.write INSIDE the allowlist (which the risk gate raises to a real L1
// window), and after the veto travels back it closes the task with a sentence.
func residentCardGolden(t *testing.T, target string) string {
	t.Helper()
	args, err := json.Marshal(map[string]string{
		"path":    filepath.ToSlash(target),
		"content": "这张卡片被否决了，这段文字不该落盘",
	})
	if err != nil {
		t.Fatal(err)
	}
	toolCall := sseChunk(t, map[string]any{"index": 0, "delta": map[string]any{
		"tool_calls": []any{map[string]any{
			"index": 0, "id": "call_246r2", "type": "function",
			"function": map[string]any{"name": "fs.write", "arguments": string(args)},
		}},
	}})
	finish := sseChunk(t, map[string]any{"index": 0, "delta": map[string]any{}, "finish_reason": "tool_calls"})
	usage := `data: {"choices":[],"usage":{"prompt_tokens":31,"completion_tokens":12,"total_tokens":43}}`
	reply := sseChunk(t, map[string]any{"index": 0, "delta": map[string]any{
		"content": "已按否决结果收尾，未写入任何文件",
	}})
	stop := sseChunk(t, map[string]any{"index": 0, "delta": map[string]any{}, "finish_reason": "stop"})
	return strings.Join([]string{
		"# wisp golden sse v1",
		"# @scenario 246-r2 resident pipeline: one L1 fs.write tool call, then the veto relayed back",
		"# @response 200",
		toolCall, "", finish, "", usage, "", "data: [DONE]", "",
		"# @response 200",
		reply, "", stop, "", usage, "", "data: [DONE]", "",
	}, "\n")
}

// residentSlowGolden is one long, slow first turn: enough chunks at enough pacing
// that a task is still streaming when the exit request arrives.
func residentSlowGolden(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("# wisp golden sse v1\n")
	b.WriteString("# @scenario 246-r2 exit path: a task still running when D38(e) starts\n")
	b.WriteString("# @latency 700\n")
	b.WriteString("# @response 200\n")
	for i := 0; i < 24; i++ {
		b.WriteString(sseChunk(t, map[string]any{"index": 0, "delta": map[string]any{
			"content": fmt.Sprintf("第一段 %d ", i),
		}}))
		b.WriteString("\n\n")
	}
	b.WriteString(sseChunk(t, map[string]any{"index": 0, "delta": map[string]any{}, "finish_reason": "stop"}))
	b.WriteString("\ndata: [DONE]\n")
	return b.String()
}

// sseChunk renders one OpenAI chat-completion SSE data line.
func sseChunk(t *testing.T, choice map[string]any) string {
	t.Helper()
	body, err := json.Marshal(map[string]any{"choices": []any{choice}})
	if err != nil {
		t.Fatal(err)
	}
	return "data: " + string(body)
}

// tailLine returns the last line containing needle, for the failure messages.
func tailLine(all, needle string) string {
	lines := strings.Split(all, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], needle) {
			return lines[i]
		}
	}
	return "(no line containing " + needle + ")"
}

// cardLedgerRecord returns the attributes of the shipped process's own card
// record from its JSONL ledger. sinkInstallRecord (ticket 127's decoder) carries
// only the fields that ticket's claims are about, and AC#7's claims are about the
// card record's OWN attributes - which orb state the surface booked and whether it
// asked Win32 for the cancel key - so the record is decoded as a map here rather
// than guessed at from its message text.
func cardLedgerRecord(t *testing.T, dir string) map[string]any {
	t.Helper()
	for _, name := range mustGlob117(t, filepath.Join(dir, "wisp-*.jsonl")) {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			if !strings.Contains(line, "常驻进程显示一张确认卡片") {
				continue
			}
			var rec map[string]any
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("the card record is not a JSON object: %v\n%s", err, line)
			}
			return rec
		}
	}
	return nil
}

// ------------------------------------------------------------------ observer

// armedObserver is the in-repo Esc rig started, ready, and NOT yet triggered.
// observeEsc246 runs one whole measurement; AC#7's card has a 2-3 second window
// and the press has to land inside it, so the case needs the halves split: bring
// the observer window up, wait for the card, then trigger.
type armedObserver struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	ready  chan struct{}
	result chan observerReading246
	errBuf *bytes.Buffer
	raw    string
	t      *testing.T
}

// armEscObserver246r2 starts the rig in -watch mode and waits for its READY line.
// The rig pumps its window while it waits for stdin, so it stays a live recipient
// for as long as this case needs it (a frozen window would "receive" nothing and
// that would read as a pass).
func armEscObserver246r2(t *testing.T, exe string) *armedObserver {
	t.Helper()
	cmd := exec.Command(exe, "-watch")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	a := &armedObserver{
		cmd: cmd, stdin: stdin, ready: make(chan struct{}),
		result: make(chan observerReading246, 1), errBuf: &bytes.Buffer{}, t: t,
	}
	cmd.Stderr = a.errBuf
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the armed observer: %v", err)
	}
	t.Cleanup(func() {
		_ = a.stdin.Close()
		if a.cmd.Process != nil {
			_ = a.cmd.Process.Kill()
			_, _ = a.cmd.Process.Wait()
		}
	})
	lines := make(chan string, 32)
	go func() {
		sc := bufio.NewScanner(stdout)
		sc.Buffer(make([]byte, 0, 4096), 1<<16)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()
	foreground := false
	go func() {
		var got observerReading246
		for line := range lines {
			a.raw += line + "\n"
			switch {
			case strings.HasPrefix(line, "READY"):
				foreground = strings.Contains(line, "foreground_is_mine=true")
				close(a.ready)
			case strings.HasPrefix(line, "READING"):
				got.raw = a.raw
				got.keydown = mustField246(t, line, "keydown_esc=")
				got.syskey = mustField246(t, line, "syskeydown_esc=")
				got.hotkey = mustField246(t, line, "wm_hotkey=")
				got.other = mustField246(t, line, "other_keys=")
				got.fgMine = foreground
				a.result <- got
			}
		}
	}()
	tm := observe.NewTimeout(25 * time.Second)
	for {
		select {
		case <-a.ready:
			return a
		case <-time.After(50 * time.Millisecond):
			if tm.Expired() {
				t.Fatalf("the armed observer never said READY (stdout: %s, stderr: %s)", a.raw, a.errBuf.String())
			}
		}
	}
}

// trigger is the moment the key is pressed: the case calls it only once the card
// line is on screen.
func (a *armedObserver) trigger() {
	if _, err := fmt.Fprintln(a.stdin, "go"); err != nil {
		a.t.Fatalf("trigger write: %v", err)
	}
}

// reading waits for the observer's READING line through observe's monotonic
// timeout, so no wall-clock subtraction decides when this gives up (D22 ban #4 /
// D42#9).
func (a *armedObserver) reading(t *testing.T) observerReading246 {
	t.Helper()
	tm := observe.NewTimeout(30 * time.Second)
	for {
		select {
		case got := <-a.result:
			if !got.fgMine {
				t.Fatalf("the armed observer never owned the foreground, so its delivery count says nothing "+
					"about Wisp: %s", got.raw)
			}
			t.Logf("armed observer: %s", strings.TrimSpace(got.raw))
			return got
		case <-time.After(50 * time.Millisecond):
			if tm.Expired() {
				t.Fatalf("the armed observer produced no READING within the budget (stdout: %s, stderr: %s)",
					a.raw, a.errBuf.String())
			}
		}
	}
}
