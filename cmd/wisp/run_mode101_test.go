package main

// Ticket 101 AC#2 + AC#3: the permission-mode WIRE, tested from the composition
// root rather than from inside the package that owns the storage.
//
// Ticket 90 proved internal/perm works. What it could not prove - and what this
// file exists to prove - is that `wisp run` CALLS IT. The difference is the
// whole ticket: R20/M3 ("手动选过哪档，之后一直按那档，重启也不回默认") was
// implemented, tested, green, and FALSE on a real machine until the mode was
// injected into the decision chain at this file's assembly root.
//
// What each case pins:
//
//	TestTicket101ManualSwitchSurvivesRestart      AC#2(a) a手动 switch is the next
//	                                              process's starting档, and the
//	                                              decision chain obeys it
//	TestTicket101UntouchedConfigRestartsAtDefault AC#2(b) a config that never
//	                                              mentions the key stays on the
//	                                              strictest档 across restarts
//	TestTicket101SessionGrantDoesNotCrossRestart  AC#2(c) a live D45 session grant
//	                                              changes NOTHING after a restart
//	TestTicket101ModeSwitchUsesTheRealL2Gate      M4's confirmation is wired to the
//	                                              C18 gate, not to a test double
//	TestTicket101UnreadableModeFailsLoudlyAndStrict AC#3 corrupt / unknown-version /
//	                                              unreadable storage: loud, audited,
//	                                              never relaxed
//
// The three AC#2 cases are three functions on purpose. The R20 boundary clause
// ("改模式能跨重启，会话授权不能跨重启，两条用例不许合并") is the one thing in
// this file a future edit must not "simplify" into a table: one merged
// persistence case can only ever be satisfied by the looser half, which is how a
// permanent免审通行证 gets in.
//
// The AC#2(a)/(b) contrast is deliberate and is the reason both exist: (a) shows
// an L1 write asking ZERO cards after a restart onto auto_approve, (b) shows the
// same dispatch asking ONE card for a config that was never touched. Neither
// half is vacuous: each also asserts the opposite-shaped control inside the same
// boot (a red-line R2 read still asks under auto_approve; a plain write still
// executes), so a bridge that simply stopped reaching the gate passes neither.
//
// WARNING: this package's test binary links sherpa-onnx and dies at LOAD time
// (0xc0000135) unless third_party/sherpa-onnx is on PATH - ticket 98's account,
// and the reason AC#5's judge for these cases is the injected command:
//
//	PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp/ -run Ticket101

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/secret"
)

// t101Key is the same fixed fake constant run_test.go uses: not a real
// credential, and no assertion here prints it.
const t101Key = "wire-fake-key-0123456789abcdef"

// t101host is one data dir holding one config.toml. "Restart" in every case
// below is a second runTextTask over the same dir: a new process surface, a new
// config Manager, a new bridge, the same bytes on disk - which is what a real
// restart is, and the only thing a persisted preference can legitimately be
// measured against.
type t101host struct {
	t    *testing.T
	dir  string
	base string
}

// t101boot writes a config whose [risk] block carries the given mode ("" = the
// key is ABSENT, i.e. a fresh install) and short windows, and wires it to
// mockllm + a stored fake key so a full `wisp run` succeeds.
func t101boot(t *testing.T, mode string) *t101host {
	t.Helper()
	srv := adaptertest.StartMockllm(t)
	h := &t101host{t: t, dir: t.TempDir(), base: srv.Base}
	h.writeConfig(t, mode)
	st, err := secret.NewStore(h.dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", t101Key); err != nil {
		t.Fatal(err)
	}
	return h
}

// writeConfig lays down config.toml. A nil-ish base ("" ) makes the provider
// point at a closed port, which the mode-failure cases exploit: they never get
// as far as an HTTP request, and if a future change LET them get that far the
// run fails for a different reason and the assertions notice.
func (h *t101host) writeConfig(t *testing.T, mode string) {
	t.Helper()
	if h.base == "" {
		h.base = "http://127.0.0.1:9"
	}
	riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"
	if mode != "" {
		riskLines += fmt.Sprintf("permission_mode = %q\n", mode)
	}
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

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

%s`, h.base+"/v1", filepath.ToSlash(h.dir), riskLines)
	if err := os.WriteFile(h.configPath(), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// start runs one `wisp run` over this data dir and returns the exit code plus
// both streams merged (the audit lines are what most assertions read). hook
// receives the assembled stack while the store is still open.
func (h *t101host) start(t *testing.T, confirm perm.ConfirmFunc, hook func(*agentRuntime)) (int, string) {
	t.Helper()
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	code := runTextTask(runSpec{
		argv:        []string{"总结一下 这份笔记"},
		stdout:      out,
		stderr:      errb,
		dataDir:     h.dir,
		modeConfirm: confirm,
		onRuntime: func(rt *agentRuntime) {
			if hook != nil {
				hook(rt)
			}
		},
		notify: func(string, string) error { return nil },
	})
	return code, out.String() + errb.String()
}

// configPath is the file the mode lives in.
func (h *t101host) configPath() string { return filepath.Join(h.dir, configFileName) }

// readMode returns the档 as the NEXT process would read it straight off disk -
// a fresh Manager, no runtime state - so "the file says auto_approve" stays a
// separate fact from "the running process does".
func (h *t101host) readMode(t *testing.T) risk.Mode {
	t.Helper()
	mgr, err := config.NewManager(h.configPath(), nil)
	if err != nil {
		t.Fatalf("the persisted config no longer loads: %v", err)
	}
	return mgr.Config().PermissionMode()
}

// t101call dispatches one host-side tool call through the bridge THIS boot
// assembled, registering the task first the way the loop does (D47 has no
// host-initiated exception). n keeps the ids distinct per call so the D31
// bookkeeping cannot cross-match two dispatches from the same hook.
//
// It returns how many cards this run's own UI has displayed so far, the
// user-visible text/error, and whether the call came back an error outcome.
func t101call(rt *agentRuntime, n int, tool, path string) (int, string, bool) {
	args := json.RawMessage(fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(path)))
	if tool == "fs.write" {
		args = json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"touched through the chain"}`,
			filepath.ToSlash(path)))
	}
	task, corr := fmt.Sprintf("t101-task-%d", n), fmt.Sprintf("t101-corr-%d", n)
	revoke := rt.gate.AdmitTextTask(task)
	defer revoke()
	res, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
		TaskID: task, CorrelationID: corr, CallID: fmt.Sprintf("t101-call-%d", n),
		Name: tool, Args: args,
	})
	if err != nil {
		return rt.windowCount(), err.Error(), true
	}
	return rt.windowCount(), res.Text, res.IsError
}

// ---------------------------------------------------------------------------
// AC#2(a): the manually chosen档 is what the next process starts with
// ---------------------------------------------------------------------------

func TestTicket101ManualSwitchSurvivesRestart(t *testing.T) {
	h := t101boot(t, "") // fresh install: the key is not in the file at all
	var (
		mu       sync.Mutex
		asked    []risk.Mode // which switches the confirmation channel was called for
		lastSwch perm.Switch // the record the last switch produced
	)
	// The stub stands for "the native side clicked allow once" - a console run
	// has no native channel, so the production adapter always refuses (pinned by
	// TestTicket101ModeSwitchUsesTheRealL2Gate, not faked away here).
	confirm := func(_ context.Context, sw perm.Switch) error {
		mu.Lock()
		defer mu.Unlock()
		asked = append(asked, sw.To)
		lastSwch = sw
		return nil
	}
	hook1 := func(rt *agentRuntime) {
		if got := rt.modes.PermissionMode(); got != risk.ModeAskEveryStep {
			t.Errorf("boot 1 starts on %v, want the R20/M2 default ask_every_step", got)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// R20/M4 is asymmetric: only the switch INTO auto_approve costs a card.
		if err := rt.modes.Set(ctx, risk.ModeAskHighRisk, "cli", "t101"); err != nil {
			t.Errorf("Set(ask_high_risk): %v", err)
		}
		if err := rt.modes.Set(ctx, risk.ModeAutoApprove, "cli", "t101"); err != nil {
			t.Errorf("Set(auto_approve): %v", err)
		}
		if rt.modes.PermissionMode() != risk.ModeAutoApprove {
			t.Fatalf("boot 1 did not apply the switch: %v", rt.modes.PermissionMode())
		}
	}
	code, log := h.start(t, confirm, hook1)
	if code != 0 {
		t.Fatalf("boot 1 exit %d\n%s", code, log)
	}
	mu.Lock()
	gotAsked := append([]risk.Mode(nil), asked...)
	mu.Unlock()
	if len(gotAsked) != 1 || gotAsked[0] != risk.ModeAutoApprove {
		t.Errorf("L2 confirmations asked = %v, want exactly [auto_approve] (R20/M4)", gotAsked)
	}
	if lastSwch.Origin != "cli" || lastSwch.Actor != "t101" {
		t.Errorf("the switch record has origin=%q actor=%q; a mode change must be attributable",
			lastSwch.Origin, lastSwch.Actor)
	}
	if !strings.Contains(log, "MODE-SWITCH from=ask_high_risk to=auto_approve") ||
		!strings.Contains(log, "result=applied") {
		t.Errorf("no audited switch line in:\n%s", log)
	}
	// The persisted form is the config key, on disk, as a string a restart can
	// read - not a runtime value the next process would have to be told about.
	// (SaveFile writes canonical TOML, which quotes with single quotes, so the
	// assertion names the key and the value separately rather than one spelling.)
	raw, err := os.ReadFile(h.configPath())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "permission_mode") ||
		!strings.Contains(string(raw), "auto_approve") {
		t.Errorf("config.toml does not carry the chosen档:\n%s", raw)
	}

	// ---- restart ----
	// No confirm stub in boot 2: nothing here switches, it only reads. This is
	// the half AC#4's mutation targets - delete the assembly line and boot 2 is
	// back on ask_every_step, which fails the first assertion below.
	var (
		boot2Mode  risk.Mode
		cards      int
		text       string
		isErr      bool
		redLineTxt string
		redLineErr bool
	)
	hook2 := func(rt *agentRuntime) {
		boot2Mode = rt.modes.PermissionMode()
		// (1) the L1 pre-execution window: auto_approve must silence it.
		cards, text, isErr = t101call(rt, 1, "fs.write", filepath.Join(h.dir, "note.txt"))
		// (2) NON-VACUITY CONTROL: silence is not "the chain stopped asking
		// anything". A red-line class (R2, outside the allowlist) still asks, so
		// a bridge that never reached the gate cannot pass this test.
		_, redLineTxt, redLineErr = t101call(rt, 2, "fs.read", "C:/Windows/win.ini")
	}
	code, log = h.start(t, nil, hook2)
	if code != 0 {
		t.Fatalf("boot 2 exit %d\n%s", code, log)
	}
	if boot2Mode != risk.ModeAutoApprove {
		t.Errorf("AC#2(a) FAILS: after the restart the chain is on %v, want the manually "+
			"chosen auto_approve (that sentence IS R20/M3)", boot2Mode)
	}
	if got := h.readMode(t); got != risk.ModeAutoApprove {
		t.Errorf("a cold read of config.toml says %v, want auto_approve", got)
	}
	if cards != 0 {
		t.Errorf("auto_approve still opened %d card(s) for an L1 write; the档 was read but "+
			"never reached the decision chain:\n%s", cards, log)
	}
	if isErr {
		t.Errorf("the silenced L1 write should have executed, got: %s", text)
	}
	if !redLineErr || !strings.Contains(redLineTxt, "拒绝") {
		t.Errorf("auto_approve got past a red-line R2 read (out=%q)", redLineTxt)
	}
	// The startup line is the audit witness that THIS boot read the档; without
	// it a restart would be silent about the posture it came up in.
	if !strings.Contains(log, "MODE-READ origin=startup mode=auto_approve") {
		t.Errorf("boot 2 wrote no startup mode line:\n%s", log)
	}
	if fi, statErr := os.Stat(filepath.Join(h.dir, "note.txt")); statErr != nil || fi.Size() == 0 {
		t.Errorf("the silenced write never landed on disk: %v", statErr)
	}
}

// ---------------------------------------------------------------------------
// AC#2(b): never manually chosen => a restart lands on the default档
// ---------------------------------------------------------------------------

func TestTicket101UntouchedConfigRestartsAtDefault(t *testing.T) {
	h := t101boot(t, "") // the key is ABSENT, not written as the default
	raw, err := os.ReadFile(h.configPath())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "permission_mode") {
		t.Fatalf("the fixture wrote the key, so this case would not be measuring a "+
			"never-touched config:\n%s", raw)
	}
	for boot := 1; boot <= 3; boot++ {
		var (
			got   risk.Mode
			cards int
			isErr bool
			text  string
		)
		code, log := h.start(t, nil, func(rt *agentRuntime) {
			got = rt.modes.PermissionMode()
			// The default档 ASKS about an L1 write: the window opens and the call
			// waits out its one second (ticket 12's composed gate). A fresh path
			// per boot, because overwriting an existing file is R8 (irreversible)
			// and R8 is a red line in EVERY档 - that escalation is real, and it is
			// not what this case is measuring.
			cards, text, isErr = t101call(rt, boot, "fs.write",
				filepath.Join(h.dir, fmt.Sprintf("note-%d.txt", boot)))
		})
		if code != 0 {
			t.Fatalf("boot %d exit %d\n%s", boot, code, log)
		}
		if got != risk.ModeAskEveryStep {
			t.Fatalf("cold start %d is on %v, want ask_every_step (R20/M2)", boot, got)
		}
		if cards != 1 {
			t.Errorf("cold start %d showed %d cards, want the L1 window to have been asked "+
				"(this is the contrast against AC#2(a)'s zero):\n%s", boot, cards, log)
		}
		if isErr {
			t.Errorf("cold start %d: an unvetoed L1 window means EXECUTE, got: %s", boot, text)
		}
		if !strings.Contains(log, "MODE-READ origin=startup mode=ask_every_step") {
			t.Errorf("cold start %d did not audit the default it read:\n%s", boot, log)
		}
	}
	// Three restarts of a never-touched config must not have written the key:
	// reading a default is not the same act as choosing one.
	raw2, err := os.ReadFile(h.configPath())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw2), "permission_mode") {
		t.Logf("the key is still absent after 3 cold starts, as expected:\n%s", raw2)
	} else {
		t.Errorf("the host persisted a permission_mode nobody asked it to persist:\n%s", raw2)
	}
}

// ---------------------------------------------------------------------------
// AC#2(c): the session grant does NOT cross the boundary - opposite expectation
// ---------------------------------------------------------------------------

func TestTicket101SessionGrantDoesNotCrossRestart(t *testing.T) {
	// auto_approve on purpose: it is the one档 that could turn a "remembered
	// grant" into a silent write, so this is where the boundary is load-bearing.
	// AC#2(a) owns the question of how the档 got there; here it is the operator's
	// own persisted declaration, written before the process started.
	h := t101boot(t, risk.ModeAutoApproveName)
	env := filepath.Join(h.dir, ".env") // B-tier per SPEC-06 §4.1 (.env*)
	if err := os.WriteFile(env, []byte("OLD_SECRET=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	const session = "session-before-restart"

	var (
		sameSession int
		firstErr    bool
		firstWhy    string
		id          int64
	)
	code, log := h.start(t, nil, func(rt *agentRuntime) {
		now := time.Now().Unix()
		var err error
		id, err = rt.store.InsertGrant(ctx, memory.ApprovalGrant{
			Scope:     memory.GrantScopeSession,
			Tool:      "fs.write",
			Pattern:   filepath.ToSlash(env),
			SessionID: session,
			CreatedAt: now,
			ExpiresAt: now + 3600, // live by its OWN clock: only the session may die
		})
		if err != nil {
			t.Fatalf("InsertGrant: %v", err)
		}
		rows, err := rt.store.ListGrantsBySession(ctx, session)
		if err != nil {
			t.Fatalf("ListGrantsBySession: %v", err)
		}
		sameSession = len(rows)
		// Refused with the grant live IN ITS OWN SESSION, too: the assembly hands
		// the bridge a mode, never a grant source (Options.Confirmations nil).
		_, firstWhy, firstErr = t101call(rt, 1, "fs.write", env)
	})
	if code != 0 {
		t.Fatalf("boot 1 exit %d\n%s", code, log)
	}
	if sameSession != 1 {
		t.Fatalf("the fixture's grant row is not live and readable (%d rows): every assertion "+
			"below would be vacuous", sameSession)
	}
	if !firstErr {
		t.Fatalf("boot 1 already let the B-tier .env write through (%q): the grant boundary "+
			"cannot be measured against a chain that never asked", firstWhy)
	}

	// ---- restart: a new process, therefore a new session ----
	var (
		after    int
		refused  bool
		why      string
		plainErr bool
		plainTxt string
		mode     risk.Mode
	)
	code, log = h.start(t, nil, func(rt *agentRuntime) {
		rows, err := rt.store.ListGrantsBySession(ctx, "session-after-restart")
		if err != nil {
			t.Fatalf("ListGrantsBySession after restart: %v", err)
		}
		after = len(rows)
		mode = rt.modes.PermissionMode()
		_, why, refused = t101call(rt, 1, "fs.write", env)
		// NON-VACUITY CONTROL, other side of the same boot: the ordinary
		// in-allowlist write IS silenced by auto_approve and executes, so the
		// .env refusal cannot be explained by "everything asks anyway" or by a
		// bridge that stopped routing calls.
		_, plainTxt, plainErr = t101call(rt, 2, "fs.write", filepath.Join(h.dir, "plain.txt"))
	})
	if code != 0 {
		t.Fatalf("boot 2 exit %d\n%s", code, log)
	}
	if after != 0 {
		t.Errorf("%d grant rows are visible to the new session, want 0 (PLAN.md:1640)", after)
	}
	if mode != risk.ModeAutoApprove {
		t.Errorf("boot 2 is on %v, want auto_approve: AC#2(c) must measure the grant boundary "+
			"under the loosest档, and AC#2(a) owns the mode question", mode)
	}
	if !refused {
		t.Errorf("AC#2(c) FAILS: after the restart the session grant changed what the chain "+
			"asks - the B-tier write went through (%q). Reading a grant back off disk is the "+
			"permanent免审通行证 this case exists to keep out.", why)
	}
	if plainErr {
		t.Errorf("control write should NOT error under auto_approve (%s); the control half is "+
			"broken, so the refusal above proves nothing", plainTxt)
	}
	if !strings.Contains(why, "R3") && !strings.Contains(why, "敏感") &&
		!strings.Contains(why, "审批") && !strings.Contains(why, "拒绝") {
		t.Errorf("the .env refusal names no sensitive-path reason: %q", why)
	}
	// The row itself survives for audit, keyed to the dead session.
	st, err := memory.Open(h.dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	all, err := st.ListGrants(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || all[0].SessionID != session || all[0].ID != id {
		t.Errorf("grant rows after restart = %+v, want the one dead-session row kept (id %d)", all, id)
	}
}

// ---------------------------------------------------------------------------
// M4's confirmation channel is the real gate, not a seam the test brought
// ---------------------------------------------------------------------------

func TestTicket101ModeSwitchUsesTheRealL2Gate(t *testing.T) {
	h := t101boot(t, "")
	var (
		setErr error
		cards  int
		after  risk.Mode
	)
	code, log := h.start(t, nil, func(rt *agentRuntime) {
		// modeConfirm is nil here, so perm.Store calls the PRODUCTION adapter:
		// one L2 card through approval.Gate, on a task the host admitted itself.
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()
		setErr = rt.modes.Set(ctx, risk.ModeAutoApprove, "cli", "t101")
		cards = rt.windowCount()
		after = rt.modes.PermissionMode()
	})
	if code != 0 {
		t.Fatalf("exit %d\n%s", code, log)
	}
	if setErr == nil {
		t.Error("a console run granted auto_approve without a native click; SPEC-06:1588 says " +
			"只有原生侧点击才是\"允许\"")
	}
	if cards != 1 {
		t.Errorf("cards shown = %d, want exactly one L2 card per switch (R20/M4)", cards)
	}
	if after != risk.ModeAskEveryStep {
		t.Errorf("an unanswered confirmation still moved the档 to %v", after)
	}
	line := ""
	for _, l := range strings.Split(log, "\n") {
		if strings.Contains(l, "MODE-SWITCH") && strings.Contains(l, "to=auto_approve") {
			line = l
		}
	}
	if !strings.Contains(line, "result=refused-confirm") &&
		!strings.Contains(line, "result=confirm-failed") {
		t.Errorf("the unanswered L2 left no refused audit line, got %q", line)
	}
	if !strings.Contains(line, modeSwitchToolName) && !strings.Contains(log, modeSwitchToolName) {
		t.Errorf("nothing in the audit names the card that was raised (%s):\n%s",
			modeSwitchToolName, log)
	}
	if got := h.readMode(t); got != risk.ModeAskEveryStep {
		t.Errorf("an unconfirmed switch persisted %v; only a confirmed one may move the档", got)
	}
}

// ---------------------------------------------------------------------------
// AC#3: the three ways the mode cannot be read
// ---------------------------------------------------------------------------

func TestTicket101UnreadableModeFailsLoudlyAndStrict(t *testing.T) {
	// NON-VACUITY CONTROL for all three states below, run first and cheaply: an
	// intact store boots, assembles a decision chain and runs relaxed. Without
	// this, "exit 2 + no chain + strictest档" would pass for a fixture that could
	// never have worked in the first place.
	ok := t101boot(t, risk.ModeAutoApproveName)
	var (
		okMode      risk.Mode
		okAssembled bool
	)
	code, log := ok.start(t, nil, func(rt *agentRuntime) {
		okAssembled, okMode = true, rt.modes.PermissionMode()
	})
	if code != 0 || !okAssembled || okMode != risk.ModeAutoApprove {
		t.Fatalf("control boot over an intact store: exit %d assembled=%v mode=%v\n%s",
			code, okAssembled, okMode, log)
	}
	if !strings.Contains(log, "mode=ask_every_step") {
		t.Logf("control log for reference (no fail-closed line expected here):\n%.400s", log)
	}

	cases := []struct {
		name   string
		mutate func(t *testing.T, dir string)
	}{
		{"存储损坏", func(t *testing.T, dir string) {
			if err := os.WriteFile(filepath.Join(dir, configFileName),
				[]byte("schema_version = 2\n[[[\nthis is not a config\n"), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{"版本不认识", func(t *testing.T, dir string) {
			if err := os.WriteFile(filepath.Join(dir, configFileName),
				[]byte("schema_version = 99\n"), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{"权限读不到", func(t *testing.T, dir string) {
			p := filepath.Join(dir, configFileName)
			if err := os.Remove(p); err != nil {
				t.Fatal(err)
			}
			// A DIRECTORY where the file must be: opening it for reading fails on
			// every platform, which is the honest portable stand-in for "the ACL
			// denied us" (that axis is ticket 89's, and icacls in a test would
			// leave t.TempDir's cleanup unable to delete the tree).
			if err := os.Mkdir(p, 0o700); err != nil {
				t.Fatal(err)
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// auto_approve LOOSE on disk first, then the store breaks underneath
			// it: the assertion is about what the failed read answers with.
			h := &t101host{t: t, dir: t.TempDir()}
			h.writeConfig(t, risk.ModeAutoApproveName)
			if _, err := secret.NewStore(h.dir); err != nil {
				t.Fatal(err)
			}
			tc.mutate(t, h.dir)
			reached := false
			code, log := h.start(t, nil, func(*agentRuntime) { reached = true })
			if code == 0 {
				t.Errorf("exit 0 over an unreadable mode store; AC#3 wants a loud failure")
			}
			if code != 2 {
				t.Errorf("exit %d, want 2 (Unconfigured: SPEC-03 §4.1 forbids a half-configured run)\n%s",
					code, log)
			}
			if reached {
				t.Error("a decision chain was assembled over an unreadable mode store")
			}
			// Loud AND strict: the fallback档 is named in the audit, in the same
			// "[audit] perm:" family the store writes, so a log reader never has
			// to guess which posture a failed read answered with.
			if !strings.Contains(log, "MODE-READ-FAILED") {
				t.Errorf("no MODE-READ-FAILED audit line for %s:\n%s", tc.name, log)
			}
			if !strings.Contains(log, "mode="+risk.ModeAskEveryStepName) {
				t.Errorf("the failure did not name the strictest档 it falls back to:\n%s", log)
			}
			if !strings.Contains(log, "配置未就绪") {
				t.Errorf("the user-visible line is missing for %s:\n%s", tc.name, log)
			}
			// The forbidden alternative: answering from a cached looser value.
			// auto_approve WAS on disk here, and it is nowhere in what this boot
			// claims.
			if strings.Contains(log, "mode=auto_approve") {
				t.Errorf("a broken mode store still claimed auto_approve:\n%s", log)
			}
		})
	}
}
