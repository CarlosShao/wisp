package approval_test

// Ticket 340 (form 链, owner-ruled 2026-11-02, write surface internal/agent/
// approval/**): the queue's clash rewrite of a correlation the loop ACTUALLY
// minted has exactly one real producer in this tree — Queue.push
// (approval/queue.go:159, rewrite at :171) — and exactly one real minter —
// callCorr (agent/loop.go:634, sole dispatch call site :707). Before this
// case no test ever ran the two chained: the stacked form
// "<taskID>#<callID>#<seq>" existed only in hand-built fixtures, so a change
// to the rewrite format string could not turn anything red anywhere.
//
// THE CHAIN THIS CASE RUNS, mint side fully production, queue side fully
// production, joined through the queue's OWN exported read face (form 甲's
// "把真 corr 喂进真队列", zero new exported names, zero production edits):
//
//	golden SSE (C5 seam, the sanctioned injection point, SPEC-05 §9) replays
//	  ONE fs.edit call carrying provider id "call_e340" — fs.edit is
//	  Declared risk.L2 (fs_edit.go), so the call takes the queue route.
//	-> real agent.Loop (agent.New) -> executeCalls
//	-> loop.go:707 CorrelationID: callCorr(taskID, "call_e340", 0)
//	   = X, the production arm-A shape "<taskID>#<callID>", minted by the
//	   real minter, never written in this file.
//	-> real tools.Bridge (assessed L2) -> real Gate.PendingApproval (task
//	   admitted through AdmitTextTask, D47) -> real Queue.push #1: no clash,
//	   the item is filed under X itself.
//	-> THIS case reads X back off the production queue face —
//	   Queue.LiveApprovals()[0].Decision.CorrelationID — the corr as the
//	   queue RECORDED the loop's mint — and hands that same decision to a
//	   second real Gate.PendingApproval call. No corr string is typed by
//	   this test at any point; the second push's input IS the first push's
//	   booked value, byte for byte.
//	-> real Queue.push #2: byID[X] clash -> corr = fmt.Sprintf("%s#%d", X, 2)
//	   => the stacked form, minted by the real minter AND rewritten by the
//	   real clash branch, both running for real.
//
// WHY the second push arrives through this hop and not through a second call
// in the SAME loop turn (named census reading, ticket 340-r1, AC#0):
// internal/llm/events.go:204 accumulates tool-call fragments by ID
// (byID map[string]int), so ONE turn can never yield TWO llm.ToolCalls that
// share a provider id — callCorr therefore cannot mint one corr twice within
// one task from the wire side, and every same-task re-dispatch across turns
// lands only AFTER the first card left byID (deliver/dropLocked delete the
// key, queue.go:312), so no clash. The stacked shape is production-reachable
// through re-issued decisions (hosts restamping the C18 id, replays, retry
// gaps), and the queue's rewrite is the component under test here — this
// case feeds it the minted corr through the gate's own exported route, which
// is exactly the "真 corr 喂进真队列" form 甲 the ticket priced.
//
// What the case then reads off production faces ONLY:
//   - Queue.LiveApprovals: the two issued keys, and each item's OWN incoming
//     Decision.CorrelationID — the clash witness.
//   - memory tool_call rows (written by the bridge, bridge.go): the corr X
//     as booked by the real dispatch, never rebuilt here.
//   - panel.TaskRosterSectionFrom fed through panel.NativeVerdict.CardView
//     (the same pure hop the host's liveVerdicts uses): the stacked key must
//     light the ORIGINAL task row — and a sibling row must stay dark.
//   - the two cards' displayed CorrelationIDs handed to the UI: what the
//     queue issued is what a human would be shown and reply to.
//
// Both cards are then refused on their EXACT issued keys (DecideFromNative,
// reject needs no grant); the task completes, the file is untouched, the
// queue drains. The teeth (ticket AC#3 去牙, run on a temp-root tree so the
// shared worktree never carried a mutated queue.go): mutate queue.go:171's
// rewrite to not stack and this case goes RED at "the two issued keys are
// distinct"; mutate the separator and it goes RED at a DIFFERENT assertion
// ("the rewritten key carries the stacked form"). Both red lines were
// captured verbatim in the leg's delivery report.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat" // protocol registration
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// dupCallID340 is the provider-side call id the replay gives the ONE loop
// dispatch. It is deliberately NOT shaped like a queue rewrite tail (no
// all-decimal remainder): the incoming corr X must carry exactly one '#'.
// It is not a fixture corr — the bytes under test are what callCorr builds
// from it inside the real loop.
const dupCallID340 = "call_e340"

// e2eUI340 records the prompts the real gate displays and returns nil — an
// unreachable UI would fail closed BEFORE anybody could read the queue, so
// the refusals in this case must travel the production answer route
// (DecideFromNative), not a broken presentation seam.
type e2eUI340 struct {
	mu   sync.Mutex
	card []approval.Prompt
}

func (u *e2eUI340) Prompt(_ context.Context, p approval.Prompt) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.card = append(u.card, p)
	return nil
}

func (u *e2eUI340) Update(context.Context, approval.Event) error { return nil }

func (u *e2eUI340) cards() []approval.Prompt {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]approval.Prompt(nil), u.card...)
}

// composeE2E340 wires the real loop + real bridge + real gate + one real
// fs.edit target, the same composition internal/tools/loop_approval_test.go
// proved, re-pointed at the L2 tool this chain needs.
func composeE2E340(t *testing.T) *e2e340Fixture {
	t.Helper()
	dir := sealableTempDir124(t) // ticket 124's in-package discipline
	target := filepath.Join(dir, "edit-target.txt")
	if err := os.WriteFile(target, []byte("alpha\n"), 0o600); err != nil {
		t.Fatalf("seed the edit target: %v", err)
	}

	sse := editThenTextSSE340(t, target)
	res, err := golden.Parse(sse)
	if err != nil {
		t.Fatalf("parse sse: %v", err)
	}
	rep := golden.NewReplayer(res)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))
		rep.ServeHTTP(w, r)
	}))
	t.Cleanup(srv.Close)
	prov, err := llm.NewProvider(llm.EndpointOptions{
		Provider: "mock", Model: "probe-model", Protocol: "openai-chat",
		BaseURL: srv.URL, ContextWindow: 128000,
		HTTPClient: &http.Client{Transport: &http.Transport{Proxy: nil}},
	})
	if err != nil {
		t.Fatalf("provider: %v", err)
	}

	mem, err := memory.Open(filepath.Join(dir, "db"))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = mem.Close() })

	ui := &e2eUI340{}
	gate := approval.New(approval.Options{UI: ui, Channels: approval.NewChannels()})
	// The allowlist entry is folded by the PRODUCTION canonicalizer itself —
	// the exact call internal/tools/wiring_test.go's mustCanon makes; this
	// file hand-rolls no filepath.Clean/Abs of its own (ban 2's shape).
	allow, err := tools.NewPathCanonicalizer(nil, nil).Canonicalize(dir)
	if err != nil {
		t.Fatalf("Canonicalize(%q): %v", dir, err)
	}
	deps := tools.FSDeps{Paths: tools.NewPathCanonicalizer([]string{allow}, nil)}
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(deps) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	b := tools.New(tools.Options{
		Registry: reg, Paths: deps.Paths, Gate: gate,
		Cancel: gate.ToolsCancelBus(), Journal: mem, Provenance: nil,
		Logf: func(string, ...any) {},
	})
	loop, err := agent.New(agent.Options{
		Provider: prov, Tools: b, Sink: &agent.RecordSink{},
		Registry:  observe.NewRegistry(),
		Journal:   nil, // the bridge owns the tool_call rows in this composition
		AdmitTask: gate.AdmitTextTask,
		Config: agent.Config{
			Model: "probe-model", ContextWindow: 128000,
			ArtifactsDir: filepath.Join(dir, "artifacts"),
			// The per-tool ceiling covers the whole L2 wait here (the bridge
			// dispatch ctx is the PendingApproval wait too), so it must be
			// generous against assertion time, not the 10s loop_approval
			// used for its fast window.
			PerToolTimeout: 60 * time.Second,
			// fs.edit classifies itself; the fail-closed default stays on.
			PassThroughUnclassifiedRisk: false,
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	return &e2e340Fixture{loop: loop, gate: gate, ui: ui, mem: mem, target: target}
}

type e2e340Fixture struct {
	loop   *agent.Loop
	gate   *approval.Gate
	ui     *e2eUI340
	mem    *memory.Store
	target string
}

// editThenTextSSE340 renders the two golden responses: ONE turn asking for
// fs.edit, then a final text turn. The bytes are the golden format,
// unmodified — the provider replay is the only seam (C5), and it supplies a
// call id, never a corr.
func editThenTextSSE340(t *testing.T, target string) []byte {
	t.Helper()
	argObj, err := json.Marshal(map[string]any{
		"path": filepath.ToSlash(target),
		"edits": []map[string]string{
			{"old": "alpha", "new": "beta"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	argStr, err := json.Marshal(string(argObj))
	if err != nil {
		t.Fatal(err)
	}
	return []byte(fmt.Sprintf(`# wisp golden sse v1
# @scenario ticket 340 chain: response 1 = one L2 fs.edit call, response 2 = final text
# @response 200
data: {"choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":%[2]q,"type":"function","function":{"name":"fs.edit","arguments":%[1]s}}]},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}

data: {"choices":[],"usage":{"prompt_tokens":30,"completion_tokens":8,"total_tokens":38}}

data: [DONE]
# @response 200
data: {"choices":[{"index":0,"delta":{"content":"编辑没有被批准，未执行。"},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: {"choices":[],"usage":{"prompt_tokens":50,"completion_tokens":4,"total_tokens":54}}

data: [DONE]
`, argStr, dupCallID340))
}

// waitDepth polls the PRODUCTION queue face until the depth is reached. The
// bound is generous (2 per 25ms tick = 50s, well under the queue's own 300s
// deadline); an unreachable depth means the chain broke upstream, and the
// caller's failure line says so with the live state attached.
func waitDepth(t *testing.T, q *approval.Queue, want int) []approval.LiveApproval {
	t.Helper()
	for spin := 0; spin < 2000; spin++ {
		if q.Depth() >= want {
			return q.LiveApprovals()
		}
		time.Sleep(25 * time.Millisecond)
	}
	live := q.LiveApprovals()
	t.Fatalf("queue depth never reached %d (last LiveApprovals=%+v): the chain broke "+
		"before Queue.push, not in the rewrite this case measures", want, live)
	return nil
}

// TestQueueClashRewriteOfARealLoopCorrelationEndToEnd is the ticket-340 shot:
// the real Queue.push clash rewrite acting on the real callCorr product, fed
// back through the real Gate.PendingApproval route, read back through the
// production queue face, the production journal rows and the production
// roster join. Nothing in the asserted corr chain is hand-built here.
func TestQueueClashRewriteOfARealLoopCorrelationEndToEnd(t *testing.T) {
	f := composeE2E340(t)

	ctx, cancel := context.WithTimeout(t.Context(), 4*time.Minute)
	defer cancel()
	resCh := make(chan agent.Result, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("loop goroutine panicked (D37b): %v", r)
			}
		}()
		resCh <- f.loop.Run(ctx, "把 alpha 改成 beta")
	}()

	// Stage 1 — the real mint, the first real push. Wait for the loop's ONE
	// fs.edit call to sit in the queue.
	live := waitDepth(t, f.gate.Queue(), 1)
	X := live[0].CorrelationID
	booked := live[0].Decision.CorrelationID
	taskID := live[0].Decision.TaskID
	if taskID == "" {
		t.Fatal("the queued decision carries no task id — the join below would prove nothing")
	}
	// (1a) The first card is filed under the incoming corr itself (no clash
	// yet) and the queue RECORDED that same corr on its item: X is then the
	// corr the real callCorr minted, read off the production queue face —
	// this value, byte for byte, is what stage 2 feeds back.
	if X == "" || booked != X {
		t.Fatalf("the first card must be filed under its own incoming corr: issued=%q booked=%q", X, booked)
	}
	// (1b) The minter's arm A ran for real: X is the production shape
	// <taskID>#<callID> over the provider id the replay carried, with exactly
	// one '#'. ⛔ a bare task id here means the mint point moved and the
	// clash below would be built on something else.
	if X == taskID || !strings.HasPrefix(X, taskID+"#") {
		t.Fatalf("minted corr %q is not the production shape <taskID>#<callID> "+
			"for task %q: callCorr's mint point moved", X, taskID)
	}
	if strings.Count(X, "#") != 1 || !strings.HasSuffix(X, "#"+dupCallID340) {
		t.Fatalf("minted corr %q does not end in the call's own id %q with a single "+
			"separator: the minter is no longer producing what the rewrite stacks onto",
			X, dupCallID340)
	}

	// Stage 2 — the real clash rewrite of the real corr. The decision handed
	// to the second PendingApproval IS the first push's booked decision read
	// off the queue face (cloneDecision already copied it), so the corr the
	// rewrite stacks onto is production-minted, never typed here.
	reDispatched := live[0].Decision
	if reDispatched.CorrelationID != X {
		t.Fatalf("the queue face handed back a different incoming corr (%q vs %q) — the "+
			"read side moved and the chain this case names would no longer be the one it runs",
			reDispatched.CorrelationID, X)
	}
	answerCh := make(chan tools.Answer, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("re-dispatch goroutine panicked (D37b): %v", r)
			}
		}()
		a, why := f.gate.PendingApproval(ctx, reDispatched)
		if why == "" {
			t.Errorf("the refused re-dispatch returned no reason sentence")
		}
		answerCh <- a
	}()

	live = waitDepth(t, f.gate.Queue(), 2)
	key0, key1 := live[0].CorrelationID, live[1].CorrelationID
	in0, in1 := live[0].Decision.CorrelationID, live[1].Decision.CorrelationID

	// (2a) The clash witness: BOTH items record the SAME incoming corr — the
	// production X — which is the only condition the rewrite branch keys on.
	if in0 == "" || in0 != in1 || in0 != X {
		t.Fatalf("the two queued items must record the same incoming corr %q (the clash the "+
			"rewrite keys on): got %q vs %q", X, in0, in1)
	}
	// (2b) The rewrite: the second card is keyed "<X>#<seq>" with a distinct
	// key and seq suffix 2 (this queue saw exactly two pushes; the clash
	// branch stamps the item's own seq). This assertion is the FIRST pair of
	// teeth: a rewrite that stopped stacking (mutation M1) lands here.
	if key0 == key1 {
		t.Fatalf("the two issued keys are NOT distinct: both cards live under %q — "+
			"Queue.push no longer rewrites the clash, so <corr>#<seq> is not produced", key0)
	}
	if key0 != X {
		t.Fatalf("the first card should keep the plain corr (no clash on the empty queue): key0=%q, incoming=%q", key0, X)
	}
	// (2c) The SECOND pair of teeth, a different assertion line: the stacked
	// form itself. Changing the separator (mutation M2) lands here, not at
	// (2b) — the two failure points must stay distinct (ticket AC#3 成对).
	if strings.Count(key1, "#") != 2 || !strings.HasPrefix(key1, X+"#") {
		t.Fatalf("the rewritten key %q does not carry the stacked form %s#<seq>: "+
			"the queue's rewrite format moved", key1, X)
	}
	// (2d) The stripped segment is still the PRODUCTION corr, never a bare
	// task id (ticket AC#3): cut at the last '#', the remainder must be
	// decimal and the base must byte-equal X.
	last := strings.LastIndexByte(key1, '#')
	tail := key1[last+1:]
	if last <= 0 || tail == "" || strings.TrimLeft(tail, "0123456789") != "" {
		t.Fatalf("the rewritten key %q does not end in \"#\" + decimal seq (got tail=%q): "+
			"the queue rewrites with something %s#%d no longer is", key1, tail, "corr", 1)
	}
	seq, err := strconv.Atoi(tail)
	if err != nil || seq != 2 {
		t.Fatalf("the rewritten key %q carries tail %q, want decimal seq 2 (the clash branch "+
			"stamps the item's own push number; Atoi err=%v)", key1, tail, err)
	}
	base := key1[:strings.LastIndexByte(key1, '#')]
	if base != X || X == taskID {
		t.Fatalf("stripping the rewrite must leave the production corr intact: base=%q want=%q "+
			"(and never the bare task id %q)", base, X, taskID)
	}

	// (2e) The cards a human is shown carry the issued keys — display rides
	// the same string the queue filed.
	for spin := 0; spin < 200 && len(f.ui.cards()) < 2; spin++ {
		time.Sleep(25 * time.Millisecond)
	}
	cards := f.ui.cards()
	if len(cards) != 2 {
		t.Fatalf("the UI saw %d prompts, want 2 (the loop's card + the re-dispatched card)", len(cards))
	}
	shown := map[string]bool{cards[0].CorrelationID: true, cards[1].CorrelationID: true}
	if !shown[key0] || !shown[key1] {
		t.Fatalf("the displayed card ids %v are not the queue's issued keys {%q,%q}: a reply "+
			"routed by the card would miss the item the rewrite renamed", cards, key0, key1)
	}

	// Stage 3 — the join through the host's own pure hop (NativeVerdict.
	// CardView, the half panel_pump.liveVerdicts uses): both live keys, fed
	// verbatim, light the ORIGINAL task row and no sibling row. The stacked
	// key is the one this ticket exists for — before 340 it only ever lit up
	// as a hand-built fixture byte sequence.
	state := panel.TaskRosterState{
		Rows: []panel.TaskRow{
			{TaskID: taskID, Label: "本轮任务", Kind: "root", Status: "Thinking", StatusKnown: true, StreamKey: taskID},
			{
				TaskID: "sibling-of-" + taskID, ParentTaskID: taskID, Label: "另一枚任务", Kind: "subagent",
				Status: "Thinking", StatusKnown: true, StreamKey: panel.SubagentStreamKey("sibling-of-" + taskID),
			},
		},
		InFlightSlots: 1, PoolCap: 4,
	}
	verdict := func(corr string) panel.ApprovalCardView {
		return panel.NativeVerdict{
			CorrelationID: corr, Tool: "fs.edit", Args: []string{filepath.ToSlash(f.target)},
			Level: risk.L2, RulesHit: []risk.RuleID{risk.R2}, Reason: "端到端用例：真队列改写真 corr",
		}.CardView()
	}
	rowLit := func(t *testing.T, name string, cs []panel.ApprovalCardView, want, dark string) {
		t.Helper()
		sect := panel.TaskRosterSectionFrom(state, cs)
		var w, d *panel.TaskRowView
		for i := range sect.Rows {
			switch sect.Rows[i].TaskID {
			case want:
				w = &sect.Rows[i]
			case dark:
				d = &sect.Rows[i]
			}
		}
		if w == nil || d == nil {
			t.Fatalf("%s: the render returned no row for %q/%+v", name, want, sect.Rows)
		}
		if !w.BlockedOnApproval {
			t.Errorf("%s: row %q reads NOT blocked although the live card(s) %+v name its "+
				"production corr %q — the stacked form does not light its own task row",
				name, want, cs, X)
		}
		if d.BlockedOnApproval {
			t.Errorf("%s: row %q lit although no card names it: the join widened past the "+
				"filed task id (ticket 328's boundaries carry into 340)", name, dark)
		}
	}
	// Both pending cards together — the live host picture.
	rowLit(t, "both-cards", []panel.ApprovalCardView{verdict(key0), verdict(key1)}, taskID, "sibling-of-"+taskID)
	// And the stacked key ALONE: the shape this ticket exists for must light
	// the row on its own merits, not inherit the plain card's lighting.
	rowLit(t, "stacked-only", []panel.ApprovalCardView{verdict(key1)}, taskID, "sibling-of-"+taskID)

	// Stage 4 — answer on the EXACT issued keys through the native
	// production route; refusals need no grant. Both items must settle.
	for _, k := range []string{key0, key1} {
		if err := f.gate.DecideFromNative(ctx, approval.Request{
			CorrelationID: k, Allow: false, Reason: "端到端用例拒绝（票 340）", Source: "340-r1",
		}); err != nil {
			t.Fatalf("refusing the card keyed %q over the native route: %v — a rewritten "+
				"card that cannot be answered by its issued key is a worse defect than the "+
				"shape this ticket measures", k, err)
		}
	}
	if d := f.gate.Queue().Depth(); d != 0 {
		t.Fatalf("queue depth after answering both cards = %d, want 0", d)
	}
	select {
	case a := <-answerCh:
		if a == tools.AnswerAllow {
			t.Fatalf("the re-dispatch answered ALLOW while refused by exact key: %+v", a)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("the re-dispatched call did not settle after its card was refused")
	}

	// Stage 5 — the task runs to completion with its call refused; the bytes
	// never executed.
	var res agent.Result
	select {
	case res = <-resCh:
	case <-time.After(2 * time.Minute):
		t.Fatal("the loop did not return after the L2 refusal")
	}
	if res.TaskID != taskID {
		t.Fatalf("the loop reported task %q while the queue filed the card under %q — the "+
			"corr prefix this case pinned is not this run's own task", res.TaskID, taskID)
	}
	if res.Status != agent.StatusCompleted {
		t.Fatalf("status = %s (%s), want completed (an L2 refusal is a tool failure, not a "+
			"task death, D37): %+v", res.Status, res.Message, res.ToolLog)
	}
	if res.ToolCalls != 1 || len(res.ToolLog) != 1 {
		t.Fatalf("ToolCalls=%d ToolLog=%d, want the one fs.edit the replay asked for: %+v",
			res.ToolCalls, len(res.ToolLog), res.ToolLog)
	}
	// The refusal vocabulary of a gate-rejected call (bridge.go:391,1108): class
	// user_rejected with the refusing sentence carried verbatim. Rejected=true is
	// the loop's OWN pre-execution refusal and stays false for anything that
	// reached the gate — asserting THIS split is part of proving the refusal came
	// from the queue route, not from an early loop bail.
	l := res.ToolLog[0]
	if l.ErrorClass != string(observe.ClassUserRejected) || !strings.Contains(l.Text, "票 340") {
		t.Fatalf("the call must read as a gate refusal (class=%q text=%q outcome=%q): %+v",
			l.ErrorClass, l.Text, l.Outcome, l)
	}
	if l.Rejected {
		t.Fatalf("Rejected is the loop's own pre-execution flag — a call that reached the queue "+
			"and was refused there must not read as an early bail: %+v", l)
	}
	if body, err := os.ReadFile(f.target); err != nil || string(body) != "alpha\n" {
		t.Fatalf("the refused edit touched the file anyway: err=%v body=%q", err, body)
	}

	// Stage 6 — the journal: the bridge booked the call under the INCOMING
	// corr X, so "<booked corr>#<seq>" IS the stacked live key byte for
	// byte, read off two production faces. This closes "剥出的段仍是生产
	// corr" (⛔ 裸任务号) from the record side.
	rows, err := f.mem.ListToolCallsByTask(ctx, res.TaskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("tool_call rows = %d, want the one dispatched call: %+v", len(rows), rows)
	}
	r := rows[0]
	if r.CorrelationID != X {
		t.Fatalf("row seq=%d books corr %q, want the minted corr %q — if the bridge "+
			"restamped rows with the rewritten key, the roster join this case feeds would "+
			"join a different string than the queue ever received", r.Seq, r.CorrelationID, X)
	}
	if r.RiskLevel != memory.RiskL2 || r.Tool != "fs.edit" {
		t.Fatalf("the call must be the L2 fs.edit the minter dispatched: %+v", r)
	}
	if r.Decision != agent.DecisionReject {
		t.Fatalf("the refusal must be booked as reject, got %q: %+v", r.Decision, r)
	}
}
