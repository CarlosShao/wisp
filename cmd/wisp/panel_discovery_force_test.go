package main

// Readface wave 3 teeth (2026-11-17), ticket 38 AC#3 + AC#5, against the REAL
// composition root: real config.toml -> real secret.Store -> real approval.Gate
// -> real tools.Bridge -> real agent.Loop; only the model's bytes and the
// notifier are injected (the fk378/run-fixture shape, not a mock of the gate).
//
// What each ruler pins, and what would go red if it lied:
//
//	Test38ToolsListRosterPillsAreTheRiskTruth  the tools.list rows carry the
//	                                           registry's own declared level
//	                                           rendered by internal/risk's
//	                                           vocabulary - a pill invented
//	                                           here cannot survive the
//	                                           row-by-row comparison.
//	Test38ToolLaunchL2RaisesARealCard          the launch door's task really
//	                                           dispatches fs.delete through
//	                                           the loop's one ToolRequest
//	                                           construction point, and the
//	                                           REAL gate displays exactly one
//	                                           L2 card (the ticket49 shape:
//	                                           count the card the host's own
//	                                           ledger booked, not a stub's
//	                                           callback).
//	Test38ToolLaunchL1RaisesNoCard             the same launch path with
//	                                           fs.write (a NEW file, L1 by
//	                                           declaration) displays no L2
//	                                           card, and the unopposed window
//	                                           means EXECUTE - the file
//	                                           lands, the row books L1/allow.
//	Test38ForceRulesFacePrecedence             the resolver's three states:
//	                                           TOML rule steers the mode and
//	                                           the source reads toml, the
//	                                           override beats the TOML and
//	                                           reads override, neither means
//	                                           the default mode and the
//	                                           default source. A broken table
//	                                           refuses whole, never half.
//	Test38ForceInstructionRidesThe200Seam      the force block is the ONE
//	                                           resolver's instruction, riding
//	                                           the ticket-200
//	                                           InstructionSource seam
//	                                           (guidance, never a frozen D39
//	                                           section), and a real task with
//	                                           a matching TOML rule books the
//	                                           FORCE-MODE audit line naming
//	                                           mode and source.
//	Test38DiscoveryDoorsRefuseByName           the door-level refusals: empty
//	                                           selector, invalid args JSON,
//	                                           unknown tool (against the REAL
//	                                           registry), nil seam - each
//	                                           names what was refused and
//	                                           touches nothing.
//	Test38ForceChipRidesTheMessageDoor         the palette chip rides
//	                                           panel.message.send's own
//	                                           payload: valid word stages,
//	                                           unknown word refuses, refused
//	                                           submission rolls back.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// the fixture: fk378's shape, plus [fs] delete_enabled so the D34 L2 row is
// really registered, and a golden script built per case.
// ---------------------------------------------------------------------------

// w38Server is the one method the fixture needs off the golden server;
// naming it keeps the struct free of a transport import.
type w38Server interface{ Close() }

type w38Flow struct {
	dir    string
	srv    w38Server
	out    *strings.Builder
	errb   *strings.Builder
	rtHook func(*agentRuntime)

	rtMu sync.Mutex
	rt   *agentRuntime
}

func w38Data(t *testing.T, payload map[string]any) string {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal sse payload: %v", err)
	}
	return "data: " + string(raw) + "\n\n"
}

// w38TurnSSE renders one openai-chat reply: role chunk, optional tool_call,
// optional text, finish, usage, [DONE] (fk378's byte shape). Turn N is served
// to request N, which is why every launch case runs its tool turn FIRST (from
// the onRuntime hook, joined before it returns) and the outer text turn LAST.
func w38TurnSSE(t *testing.T, tool string, toolArgs map[string]any, text string) string {
	t.Helper()
	var b strings.Builder
	b.WriteString(w38Data(t, map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{"role": "assistant", "content": ""}, "finish_reason": nil,
	}}}))
	finish := "stop"
	if tool != "" {
		args, merr := json.Marshal(toolArgs)
		if merr != nil {
			t.Fatalf("marshal tool args: %v", merr)
		}
		b.WriteString(w38Data(t, map[string]any{"choices": []any{map[string]any{
			"index": 0,
			"delta": map[string]any{"tool_calls": []any{map[string]any{
				"index": 0, "id": "call-w38", "type": "function",
				"function": map[string]any{"name": tool, "arguments": string(args)},
			}}},
			"finish_reason": nil,
		}}}))
		finish = "tool_calls"
	}
	if text != "" {
		b.WriteString(w38Data(t, map[string]any{"choices": []any{map[string]any{
			"index": 0, "delta": map[string]any{"content": text}, "finish_reason": nil,
		}}}))
	}
	b.WriteString(w38Data(t, map[string]any{"choices": []any{map[string]any{
		"index": 0, "delta": map[string]any{}, "finish_reason": finish,
	}}}))
	b.WriteString(`data: {"choices":[],"usage":{"prompt_tokens":11,"completion_tokens":3,"total_tokens":14}}` + "\n\n")
	b.WriteString("data: [DONE]\n")
	return b.String()
}

func w38Script(t *testing.T, turns ...string) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("# wisp golden sse v1\n")
	b.WriteString("# @scenario w38 readface wave 3 teeth\n")
	for _, turn := range turns {
		b.WriteString("# @response 200\n")
		b.WriteString(turn)
		b.WriteString("\n")
	}
	return b.String()
}

func w38WriteConfig(t *testing.T, dir, baseURL string) {
	t.Helper()
	keyRef := fixtureKeyRef181(t, dir, "dpapi:acme", fakeStoreKey)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = %q

[llm.providers.acme.models.m1]
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]
delete_enabled = true

[agent]
per_tool_timeout_ms = 60000

[risk]
confirm_timeout_sec = 31
`, baseURL+"/v1", keyRef, filepath.ToSlash(dir))
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func w38Setup(t *testing.T, build func(dir string) string) *w38Flow {
	t.Helper()
	f := &w38Flow{dir: t.TempDir(), out: &strings.Builder{}, errb: &strings.Builder{}}
	srv, _, err := golden.ServerFromBytes([]byte(build(f.dir)))
	if err != nil {
		t.Fatalf("golden.ServerFromBytes: %v", err)
	}
	t.Cleanup(srv.Close)
	f.srv = srv
	w38WriteConfig(t, f.dir, srv.URL)
	return f
}

// run drives the composition root exactly as `wisp run` does. reply makes the
// run attended (ticket 375: only an attended run may resolve an unopposed L1
// window as execute).
func (f *w38Flow) run(task string, reply io.Reader) int {
	return runTextTask(runSpec{
		argv:    []string{task},
		stdout:  f.out,
		stderr:  f.errb,
		dataDir: f.dir,
		reply:   reply,
		onRuntime: func(rt *agentRuntime) {
			f.rtMu.Lock()
			f.rt = rt
			f.rtMu.Unlock()
			if f.rtHook != nil {
				f.rtHook(rt)
			}
		},
		notify: func(string, string) error { return nil },
	})
}

func (f *w38Flow) runtime() *agentRuntime {
	f.rtMu.Lock()
	defer f.rtMu.Unlock()
	return f.rt
}

// w38L2CardsOf counts the L2 cards the host's own ledger is holding right now
// (the production reading: consoleApprovalUI.Prompt -> nativeCards.record).
func w38L2CardsOf(rt *agentRuntime) int {
	n := 0
	for _, c := range rt.liveCards.pending() {
		if c.Level == "L2" {
			n++
		}
	}
	return n
}

// w38LaunchThroughTheRealPath is the launch seam's production body minus only
// the single-slot bookkeeping that lives on the windows-only residentTaskSource:
// compose against the REAL registry (an unknown name refuses HERE), then run
// the composed book through the production execute() - the loop's dispatch is
// the one agent.ToolRequest construction point (ticket 309's nail), so the
// gate below is THE gate, not a copy. It runs SYNCHRONOUSLY: the case's tool
// turn is golden request #1, so the launch must consume it before the outer
// text turn (the goroutine join is the caller's, together with the card
// counting that needs the run blocked on the card).
func (rt *agentRuntime) w38LaunchThroughTheRealPath(t *testing.T, tool, args string) (context.CancelFunc, <-chan struct{}) {
	t.Helper()
	text, err := rt.composeToolLaunch(tool, args)
	if err != nil {
		t.Fatalf("composeToolLaunch(%s): %v", tool, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	done := make(chan struct{})
	go func() {
		defer close(done)
		code := rt.execute(ctx, text)
		if code != 0 && code != 1 {
			t.Errorf("the launch task exited %d", code)
		}
	}()
	return cancel, done
}

// ---------------------------------------------------------------------------
// AC#3 tooth 1: the roster's pills are the risk truth
// ---------------------------------------------------------------------------

func Test38ToolsListRosterPillsAreTheRiskTruth(t *testing.T) {
	f := w38Setup(t, func(dir string) string {
		return w38Script(t, w38TurnSSE(t, "", nil, "好。"))
	})
	if code := f.run("总结一下 这份笔记", nil); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.errb.String())
	}
	rt := f.runtime()
	if rt == nil {
		t.Fatal("onRuntime 没有交出这发的 runtime")
	}
	levels := map[string]risk.Level{}
	for _, e := range rt.bridge.Registry().List() {
		levels[e.Tool.Name()] = e.Decl.Declared
	}
	var (
		rows    = rt.toolsRosterRows()
		pillOK  = true
		aboveL0 = 0
	)
	for _, row := range rows {
		if row.RiskPill != levels[row.Name].String() {
			pillOK = false
			t.Errorf("row %s: riskPill=%q, want the registry's own declared level rendered by internal/risk (%q) - a pill invented here is a lie",
				row.Name, row.RiskPill, levels[row.Name].String())
		}
		if levels[row.Name] > risk.L0 {
			aboveL0++
		}
		if row.Title == "" {
			t.Errorf("row %s: the title is empty - the roster must carry the tool's own description line", row.Name)
		}
	}
	if len(rows) == 0 || !pillOK || aboveL0 == 0 {
		t.Fatalf("roster rows=%d pillOK=%v rowsAboveL0=%d - a ruler that reads nothing (or only L0) cannot tell a truthful pill from an invented one", len(rows), pillOK, aboveL0)
	}
	// The two teeth the dispatch names, verbatim: the L2 tool's pill is the
	// L2 truth word, the L1 tool's is the L1 word - both read off
	// internal/risk's own rendering, never spelled in the handler.
	wantL2, wantL1 := risk.L2.String(), risk.L1.String()
	if wantL2 == "" || wantL2 == wantL1 {
		t.Fatalf("risk vocabulary sanity: L2=%q L1=%q", wantL2, wantL1)
	}
	var sawL2Tool, sawL1Tool bool
	for _, row := range rows {
		switch {
		case row.Name == "fs.delete":
			sawL2Tool = true
			if row.RiskPill != wantL2 {
				t.Errorf("fs.delete riskPill=%q, want %q (the D34 L2 row)", row.RiskPill, wantL2)
			}
		case row.Name == "fs.write":
			sawL1Tool = true
			if row.RiskPill != wantL1 {
				t.Errorf("fs.write riskPill=%q, want %q (the D34 L1 row)", row.RiskPill, wantL1)
			}
			if !strings.Contains(row.ArgSummary, "path") {
				t.Errorf("fs.write argSummary=%q does not name the path parameter", row.ArgSummary)
			}
		}
	}
	if !sawL2Tool || !sawL1Tool {
		t.Fatalf("fixture sanity: fs.delete present=%v fs.write present=%v ([fs] delete_enabled must be on)", sawL2Tool, sawL1Tool)
	}
}

// ---------------------------------------------------------------------------
// AC#3 tooth 2: an L2 launch raises the REAL approval card
// ---------------------------------------------------------------------------

func Test38ToolLaunchL2RaisesARealCard(t *testing.T) {
	f := w38Setup(t, func(dir string) string {
		return w38Script(t,
			w38TurnSSE(t, "fs.delete", map[string]any{"path": filepath.ToSlash(filepath.Join(dir, "victim.txt"))}, ""),
			w38TurnSSE(t, "", nil, "好。"),
		)
	})
	victim := filepath.Join(f.dir, "victim.txt")
	if err := os.WriteFile(victim, []byte("to be confirmed"), 0o600); err != nil {
		t.Fatal(err)
	}
	f.rtHook = func(rt *agentRuntime) {
		// The launch door's seam, bound to the REAL compose + the production
		// execute path (the windows-only single-slot wrapper aside). The
		// compose refuses a name the registry does not hold, and the loop's
		// dispatch is the one ToolRequest construction point - the same path
		// a model-driven call travels, so the card below is the gate's own,
		// not a fixture's.
		cancel, done := rt.w38LaunchThroughTheRealPath(t, "fs.delete",
			fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(victim)))
		// Poll the host ledger until the launch's card is on it.
		deadline := time.Now().Add(15 * time.Second)
		var cardLevel string
		for time.Now().Before(deadline) {
			for _, c := range rt.liveCards.pending() {
				if c.Tool == "fs.delete" {
					cardLevel = c.Level
				}
			}
			if cardLevel != "" {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		if cardLevel != "L2" {
			t.Errorf("the displayed card for fs.delete reads level %q, want L2 (an empty level means no real card was ever displayed)", cardLevel)
		}
		if l2 := w38L2CardsOf(rt); l2 != 1 {
			t.Errorf("the host ledger held %d L2 cards, want exactly one - a second card would be a second gate", l2)
		}
		if _, err := os.Stat(victim); err != nil {
			t.Errorf("the L2 card must GATE the delete, not perform it: %v", err)
		}
		// Refuse by abandoning the launch task, then join: the outer text
		// turn is golden request #2 and must be consumed after this.
		cancel()
		<-done
	}
	if code := f.run("总结一下 这份笔记", nil); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.errb.String())
	}
}

// ---------------------------------------------------------------------------
// AC#3 tooth 3: an L1 launch displays no L2 card, and the unopposed window
// executes
// ---------------------------------------------------------------------------

func Test38ToolLaunchL1RaisesNoCard(t *testing.T) {
	f := w38Setup(t, func(dir string) string {
		// Three turns: the launch's request gets the tool call, the launch's
		// post-tool request gets its text reply, and only the THIRD request is
		// the outer task's own - the launch runs first (inside the hook), so
		// an exhausted script would read as a provider fault, not a verdict.
		return w38Script(t,
			w38TurnSSE(t, "fs.write", map[string]any{"path": filepath.ToSlash(filepath.Join(dir, "l1-target.txt")), "content": "launched from the palette"}, ""),
			w38TurnSSE(t, "", nil, "已写入。"),
			w38TurnSSE(t, "", nil, "好。"),
		)
	})
	target := filepath.Join(f.dir, "l1-target.txt")
	// The attended stream (375): an open, silent reply reader is what lets an
	// UNOPPOSED window mean execute; without it the L1 route refuses.
	pr, pw := io.Pipe()
	defer pw.Close()
	f.rtHook = func(rt *agentRuntime) {
		cancel, done := rt.w38LaunchThroughTheRealPath(t, "fs.write",
			fmt.Sprintf(`{"path":%q,"content":"launched from the palette"}`, filepath.ToSlash(target)))
		// Wait out the L1 window (2s default) plus margin, then read the
		// ledger: an L1 must never take the approval route.
		time.Sleep(4 * time.Second)
		if l2 := w38L2CardsOf(rt); l2 != 0 {
			t.Errorf("the host ledger held %d L2 cards for an L1 launch, want 0 - an L1 must never take the approval route", l2)
		}
		cancel()
		<-done
	}
	if code := f.run("总结一下 这份笔记", pr); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.errb.String())
	}
	body, err := os.ReadFile(target)
	if err != nil || string(body) != "launched from the palette" {
		t.Fatalf("the L1 launch did not execute (unopposed window must mean EXECUTE): %v %q", err, body)
	}
	rt := f.runtime()
	if rt == nil {
		t.Fatal("onRuntime 没有交出这发的 runtime")
	}
	rows := rt.taskRosterState().Rows
	if len(rows) == 0 {
		t.Fatal("the launch admitted no roster row - the launch task never ran")
	}
	// The launch task was admitted FIRST (its execute ran inside the hook,
	// before the outer text task was admitted), so its row is the first one.
	launchID := rows[0].TaskID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	st, serr := memory.Open(f.dir)
	if serr != nil {
		t.Fatal(serr)
	}
	defer func() { _ = st.Close() }()
	trows, rerr := st.ListToolCallsByTask(ctx, launchID)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(trows) == 0 {
		t.Fatal("the launch's tool call booked no forensics row")
	}
	last := trows[len(trows)-1]
	if last.Tool != "fs.write" || last.RiskLevel != memory.RiskL1 || last.Decision != agent.DecisionAllow {
		t.Errorf("the launch's row = tool=%s level=%s decision=%s, want fs.write/L1/allow", last.Tool, last.RiskLevel, last.Decision)
	}
}

// ---------------------------------------------------------------------------
// AC#5 teeth: the resolver's precedence and the face's payload
// ---------------------------------------------------------------------------

func Test38ForceRulesFacePrecedence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ForceRulesFileName)
	rt := &agentRuntime{force: newForceStore(path)}

	// (1) no table: the default mode and the default source, no instruction.
	dec, err := rt.resolveForce("清理临时文件")
	if err != nil {
		t.Fatalf("no-table resolve: %v", err)
	}
	if dec.Mode != "" || dec.Source != "default" || dec.Instruction != "" {
		t.Errorf("no-table decision = %+v, want the default (mode=\"\", source=default, no instruction)", dec)
	}

	// (2) a TOML rule steers the mode; the source reads toml.
	if err := os.WriteFile(path, []byte("[[force_rule]]\npattern = \"清理临时文件\"\nmode = \"force_tool\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	dec, err = rt.resolveForce("帮我清理临时文件")
	if err != nil {
		t.Fatalf("toml resolve: %v", err)
	}
	if dec.Mode != ForceModeTool || dec.Source != "toml" || dec.Pattern != "清理临时文件" || dec.Instruction != forceToolInstruction {
		t.Errorf("toml decision = %+v, want force_tool/toml with the D11③ instruction", dec)
	}
	// A text the pattern does not match stays default.
	if dec, _ := rt.resolveForce("讲个笑话"); dec.Mode != "" || dec.Source != "default" {
		t.Errorf("unmatched text = %+v, want the default", dec)
	}

	// (3) the override beats the TOML rule, and the source reads override.
	rt.stageForceOverride(ForceModeChat)
	dec, err = rt.resolveForce("帮我清理临时文件")
	if err != nil {
		t.Fatalf("override resolve: %v", err)
	}
	if dec.Mode != ForceModeChat || dec.Source != "override" || dec.Instruction != forceChatInstruction {
		t.Errorf("override decision = %+v, want force_chat/override (the chip must beat the table)", dec)
	}
	// The override is a one-slot chip: consumed at admission, it does not
	// outlive the task it rode.
	rt.takeForceOverride("task-w38")
	if dec, _ := rt.resolveForce("帮我清理临时文件"); dec.Source != "toml" {
		t.Errorf("after admission the decision = %+v, want the TOML rule again (a chip must not outlive its task)", dec)
	}

	// (4) a broken table refuses WHOLE (never a half-table): one unknown mode
	// word poisons every rule, and the error names it.
	if err := os.WriteFile(path, []byte("[[force_rule]]\npattern = \"a\"\nmode = \"force_tool\"\n\n[[force_rule]]\npattern = \"b\"\nmode = \"force_banana\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := rt.resolveForce("a"); err == nil {
		t.Fatal("a table with an unknown mode word resolved - it must refuse whole (不许猜)")
	}

	// (5) the face's view: rows carry their own source, and the mode is the
	// resolver's answer for the requested text.
	if err := os.WriteFile(path, []byte("[[force_rule]]\npattern = \"清理临时文件\"\nmode = \"force_tool\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	rt2 := &agentRuntime{force: newForceStore(path)}
	rt2.stageForceOverride(ForceModeChat)
	rt2.takeForceOverride("task-override")
	view, verr := rt2.forceRulesView("帮我清理临时文件")
	if verr != nil {
		t.Fatalf("forceRulesView: %v", verr)
	}
	if view.Mode != ForceModeTool || view.Source != "toml" {
		t.Errorf("view = %+v, want force_tool/toml (the consumed override must not leak into the view)", view)
	}
	var sawOverride, sawToml bool
	for _, r := range view.Rules {
		switch r.Source {
		case "override":
			sawOverride = r.Pattern == "task-override" && r.Mode == ForceModeChat
		case "toml":
			sawToml = r.Pattern == "清理临时文件" && r.Mode == ForceModeTool
		default:
			t.Errorf("rule row %q carries source %q - only override/toml exist", r.Pattern, r.Source)
		}
	}
	if !sawOverride || !sawToml {
		t.Errorf("view rules = %+v, want one override row (task-override/force_chat) and one toml row (清理临时文件/force_tool)", view.Rules)
	}
}

func Test38ForceInstructionRidesThe200Seam(t *testing.T) {
	// The wrapper is the only place the block meets the prompt: the force
	// text is prepended, the inner block follows, and the manifest names the
	// mode - guidance vocabulary, never a frozen D39 section.
	inner := &w38InstrSource{block: "项目说明正文。", manifest: []string{"AGENTS.md"}}
	f := &forceInstructionSource{inner: inner, block: forceToolInstruction}
	got := f.BlockForTurn()
	want := forceToolInstruction + "\n" + "项目说明正文。"
	if got != want {
		t.Errorf("BlockForTurn = %q, want %q (force block first, inner block intact)", got, want)
	}
	manifest := f.ManifestForTurn()
	if len(manifest) != 2 || manifest[0] != "force-mode: "+ForceModeTool || manifest[1] != "AGENTS.md" {
		t.Errorf("ManifestForTurn = %v, want [force-mode: force_tool, AGENTS.md]", manifest)
	}
	// No force block: the seam renders exactly the inner source - the default
	// path is byte-identical to before this file existed.
	bare := &forceInstructionSource{inner: inner}
	if bare.BlockForTurn() != inner.BlockForTurn() {
		t.Errorf("the bare wrapper changed the inner block: %q vs %q", bare.BlockForTurn(), inner.BlockForTurn())
	}
	if len(bare.ManifestForTurn()) != len(inner.ManifestForTurn()) {
		t.Errorf("the bare wrapper added a manifest line: %v", bare.ManifestForTurn())
	}

	// And the real run books the FORCE-MODE line: a TOML rule matching the
	// task text steers the mode, and the audit line names mode and source.
	flow := w38Setup(t, func(dir string) string {
		return w38Script(t, w38TurnSSE(t, "", nil, "好。"))
	})
	if err := os.WriteFile(filepath.Join(flow.dir, ForceRulesFileName), []byte("[[force_rule]]\npattern = \"整理这份笔记\"\nmode = \"force_tool\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if code := flow.run("请整理这份笔记", nil); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, flow.out.String(), flow.errb.String())
	}
	if !strings.Contains(flow.errb.String(), "FORCE-MODE mode=force_tool source=toml") {
		t.Errorf("the run booked no FORCE-MODE line (mode=force_tool source=toml):\n%s", flow.errb.String())
	}
}

type w38InstrSource struct {
	block    string
	manifest []string
}

func (s *w38InstrSource) BlockForTurn() string      { return s.block }
func (s *w38InstrSource) ManifestForTurn() []string { return s.manifest }

// ---------------------------------------------------------------------------
// door-level refusals
// ---------------------------------------------------------------------------

func Test38DiscoveryDoorsRefuseByName(t *testing.T) {
	f := w38Setup(t, func(dir string) string {
		return w38Script(t, w38TurnSSE(t, "", nil, "好。"))
	})
	f.run("总结一下 这份笔记", nil)
	rt := f.runtime()
	if rt == nil {
		t.Fatal("onRuntime 没有交出这发的 runtime")
	}
	envelope := func(method string, extra map[string]string) panel.ComposerRequest {
		m := map[string]string{"method": method, "requestId": "w38-door", "source": panel.ComposerRequestSource}
		for k, v := range extra {
			m[k] = v
		}
		raw, _ := json.Marshal(m)
		req, perr := panel.ParseComposerRequest(string(raw))
		if perr != nil {
			t.Fatalf("parse %s envelope: %v", method, perr)
		}
		return req
	}

	// nil seams refuse by name (the CLI-seam state).
	bare := newDiscoveryDoors(nil, "test")
	if _, err := bare.HandleToolsList(context.Background(), envelope(panel.MethodToolsList, nil)); err == nil {
		t.Error("tools.list with a nil seam answered - a machine with no pipeline must refuse by name")
	}
	if _, err := bare.HandleToolLaunch(context.Background(), envelope(panel.MethodToolLaunch, map[string]string{"toolName": "fs.read"})); err == nil {
		t.Error("tool.launch with a nil seam answered - the launch must refuse by name, never drop")
	}
	if _, err := bare.HandleRulesForce(context.Background(), envelope(panel.MethodRulesForce, nil)); err == nil {
		t.Error("rules.force with a nil seam answered - a machine with no pipeline must refuse by name")
	}

	// The launch door's validators, over the REAL registry: empty name,
	// invalid JSON args, unknown tool - and the known-tool path answers.
	d := newDiscoveryDoors(nil, "test")
	d.bind(rt.toolsRosterRows, func(ctx context.Context, tool, args string) error {
		_, err := rt.composeToolLaunch(tool, args)
		return err
	}, rt.forceRulesView)
	if _, err := d.HandleToolLaunch(context.Background(), envelope(panel.MethodToolLaunch, map[string]string{"toolName": ""})); err == nil {
		t.Error("an empty toolName was accepted - launch must name its target")
	}
	if _, err := d.HandleToolLaunch(context.Background(), envelope(panel.MethodToolLaunch, map[string]string{"toolName": "fs.read", "toolArgs": "{not json"})); err == nil {
		t.Error("invalid toolArgs JSON was accepted")
	}
	if _, err := d.HandleToolLaunch(context.Background(), envelope(panel.MethodToolLaunch, map[string]string{"toolName": "not.a.tool"})); err == nil {
		t.Error("a tool the registry does not hold was accepted - the roster is the only launch menu")
	}
	if _, err := d.HandleToolLaunch(context.Background(), envelope(panel.MethodToolLaunch, map[string]string{"toolName": "fs.read"})); err != nil {
		t.Errorf("a known tool with valid args was refused: %v", err)
	}

	// The rules face answers through the bound seam, and the payload carries
	// the view's own mode/source words.
	raw, err := d.HandleRulesForce(context.Background(), envelope(panel.MethodRulesForce, map[string]string{"forceTask": "清理临时文件"}))
	if err != nil {
		t.Fatalf("rules.force with the seam bound: %v", err)
	}
	if !strings.Contains(raw, rulesForcePayloadMarker) || !strings.Contains(raw, `"source":"default"`) {
		t.Errorf("rules.force payload = %q, want the marker and the default source", raw)
	}

	// The roster face answers with the marker and the count sentence.
	raw, err = d.HandleToolsList(context.Background(), envelope(panel.MethodToolsList, nil))
	if err != nil {
		t.Fatalf("tools.list with the seam bound: %v", err)
	}
	if !strings.Contains(raw, toolsListPayloadMarker) || !strings.Contains(raw, "已安装工具共") {
		t.Errorf("tools.list payload = %q, want the marker and the count sentence", raw)
	}
}

// ---------------------------------------------------------------------------
// the force chip's ride on panel.message.send
// ---------------------------------------------------------------------------

func Test38ForceChipRidesTheMessageDoor(t *testing.T) {
	var staged []string
	door := newPanelMessageDoor453(nil, nil, nil, "test")
	door.bindRelay453(func(ctx context.Context, rendered string, marks []panel.AttachmentRef) error {
		return nil
	})
	door.bindForceStageW3(func(mode string) { staged = append(staged, mode) })
	call := func(forceMode string) error {
		raw, _ := json.Marshal(map[string]any{
			"method": panel.MethodMessageSend, "requestId": "w38-chip",
			"source": panel.ComposerRequestSource, "text": "整理这份笔记",
			"forceMode": forceMode,
		})
		req, err := panel.ParseComposerRequest(string(raw))
		if err != nil {
			t.Fatalf("parse envelope: %v", err)
		}
		return door.HandleMessageRequest(context.Background(), req)
	}
	if err := call(ForceModeChat); err != nil {
		t.Fatalf("a valid chip was refused: %v", err)
	}
	if len(staged) != 1 || staged[0] != ForceModeChat {
		t.Fatalf("staged = %v, want exactly [force_chat]", staged)
	}
	if err := call("force_banana"); err == nil {
		t.Fatal("an unknown chip word was accepted - the vocabulary is force_tool|force_chat, refuse, never guess")
	}
	if len(staged) != 1 {
		t.Errorf("a refused chip staged something: %v", staged)
	}
	// A refused submission rolls the chip back: the next task must not inherit
	// a force it never asked for (the 453 rule, restated for the chip).
	door2 := newPanelMessageDoor453(nil, nil, nil, "test")
	door2.bindForceStageW3(func(mode string) { staged = append(staged, mode) })
	door2.bindRelay453(func(ctx context.Context, rendered string, marks []panel.AttachmentRef) error {
		return fmt.Errorf("slot busy")
	})
	raw, _ := json.Marshal(map[string]any{
		"method": panel.MethodMessageSend, "requestId": "w38-chip2",
		"source": panel.ComposerRequestSource, "text": "整理这份笔记",
		"forceMode": ForceModeTool,
	})
	req, perr := panel.ParseComposerRequest(string(raw))
	if perr != nil {
		t.Fatalf("parse envelope: %v", perr)
	}
	if err := door2.HandleMessageRequest(context.Background(), req); err == nil {
		t.Fatal("a refused submission answered success")
	}
	if len(staged) != 3 || staged[2] != "" {
		t.Errorf("rollback sequence = %v, want the chip staged then cleared (%q at the end)", staged, "")
	}
}
