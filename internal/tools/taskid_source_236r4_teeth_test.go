package tools

// ticket 236-r4 — AC#3, and ONLY its 键的来源 (host mint vs model give) dimension:
// the two places that must carry the SAME host-minted task id, pinned on a real
// run. "两者相等" had no ruler anywhere before this file (236-a4 measured zero
// in internal/tools + cmd/wisp, and named why the one shape ruler that WOULD see
// a swap — internal/agent/compress_trace_test.go:526 — stayed invisible: a
// per-package ruler cannot see the cross-package one).
//
// The two places, re-read at 2026-10-14 (the ticket's own :892 and a4's :1099
// have drifted since they were written; re-check before quoting):
//
//	internal/agent/loop.go:329   id := newTaskID()                <- the mint
//	internal/agent/loop.go:330   t := &RunningTask{ID: id, ...}   <- RunningTask.ID
//	internal/agent/loop.go:333   t.result = l.run(c, id, input)   <- the same id
//	internal/agent/loop.go:382   res := Result{TaskID: taskID, ...}
//	cmd/wisp/run.go:1105         rt.tasks.MarkRoot(bg.ID, task)
//	cmd/wisp/run.go:1110         Backfill(bg.Root().Ctx, res.TaskID, res.Text)
//	cmd/wisp/run.go:1124/1135/1145   task_log 与状态行，全部用 res.TaskID
//
// THIS pin is written in internal/tools but its射程 includes internal/agent on
// purpose: it drives the REAL loop (agent.New + RunAsync) over the frozen C5
// fixture, so a mutation at loop.go:329/333 is what turns it red — a swap at the
// mint point or a break of the mint-to-use link is invisible to any ruler that
// stays inside one package.
//
// The three segments are copied from internal/session's minted-shape lock
// (grants_test.go:152 + session.go:67/:94 + cmd/wisp/ticket224_assembly_test.go:433):
//
//	1. the lock       r236r4MintedShape below: the shape newTaskID mints
//	                  (loop.go:1153). The session lock's Valid() is exported
//	                  because production code consults it; this one is test-local
//	                  because no production caller consults a task-id predicate —
//	                  copying the LOCK, not adding a tool surface.
//	2. cheap direction  the id values the frozen C5 seam really carries, and the
//	                  hand-typed task ids this repo's fixtures use, are OUTSIDE
//	                  that shape and never equal a real mint. Ceiling, stated
//	                  because it is load-bearing: this direction is enumerated,
//	                  not structural — a model CAN hand over a 36-char uuid-shaped
//	                  value (cmd/wisp/subagent_stream_key_197_test.go:55 and
//	                  internal/tools/subagent_197_test.go:886 already use one as a
//	                  call id), and no shape lock can exclude that.
//	3. expensive direction  the id a REAL mint produced, read back out of the
//	                  real loop, passes the lock, is not any literal, AND is the
//	                  key the host's own write side files the output under — the
//	                  roster key and the backfill key are one row, not two.
//
// What each assertion is for (measured, 交件 §3 m-1/m-2):
//
//	loop.go:329 id := <a model-suppliable literal>  -> shape + literal + seam red
//	loop.go:333 run(c, newTaskID(), ...)            -> equality + roster-chain red
//
// The entropy-failure degenerate branch (loop.go:1155-1161) is NOT covered by
// this file and cannot be: rand.Read is a package-level call with no injection
// seam, and the fallback keeps the exact same 36-char shape (the formatting sits
// outside the if), so neither the lock nor the equality asserts can tell whether
// entropy failed. Opening a seam for it would be a test-only backdoor in
// production code (AGENTS §1.3), which this leg is forbidden to add.

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat" // protocol registration
)

// r236r4ModelCallID is the one model-supplied id on the frozen fixture this file
// replays (internal/agent/testdata/golden/tool-then-text.sse, verbatim), and the
// literal the source-swap mutation plants at the mint point (交件 §3 m-1). It is
// a string the MODEL gives, off the same wire production reads.
const r236r4ModelCallID = "call_e1"

// r236r4ModelWireIDs names id values the frozen C5 seam really carries. Ten of
// the twelve were re-read on 2026-10-14 from internal/agent/testdata/golden and
// internal/llm/testdata/golden (grep -rhoE '"id":"[^"]+"' over both directories,
// 33 unique ids). The other two sit in inline wire payloads elsewhere in this
// repo, not in those directories: "call_w1" (internal/tools/loop_approval_test.go:288)
// and "chatcmpl-fake" (cmd/wisp/secret_test.go:1133). None of them may validate
// as a host-minted task id: every one is a value the model supplies.
var r236r4ModelWireIDs = []string{
	"call_e1", "call_w1", "call_b2", "call_a1", "call_x1",
	"fc_1", "fc_2", "rs_1",
	"msg_mock1", "msg_out1", "resp_mock1", "chatcmpl-fake",
}

// r236r4FixtureTaskIDs are hand-typed task ids this repo's own fixtures really
// use as roster keys (parent197, r183BackgroundID, the 164/175/185 seam ids and
// the Chinese one from subagent_197_test.go). They must stay outside the minted
// shape, or a fixture could one day be mistaken for a production mint in this
// dimension's own legs.
var r236r4FixtureTaskIDs = []string{
	"parent-task-197", "c66c0634183r2a1b", "c66c0635183r2b2c", "183r2-cli",
	"164-ac2-before", "175r2-j3", "185r1-second", "bg-1", "ghost-9",
	"task-197-a", "带中文的id",
}

// r236r4MintedShape reports whether s has the shape agent.newTaskID mints
// (internal/agent/loop.go:1153): 36 characters, dashes at positions 8/13/18/23,
// lowercase hex everywhere else. This is the "plain lowercase hex uuid" the
// task_backfill.go doc block names; it is deliberately the SAME claim the one
// existing shape ruler makes (compress_trace_test.go:526: 36 chars, 4 dashes),
// tightened only to the fixed dash positions.
//
// Not checked on purpose: the v4 version/variant nibbles newTaskID sets at
// loop.go:1162-1163. Pinning implementation detail adds nothing here because
// the entropy-failure fallback sets the same bits (the formatting sits outside
// the if), so the nibbles can never tell the branches apart.
func r236r4MintedShape(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
			continue
		}
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

// r236r4Carrier builds the composition this file measures, in the shape
// pointer_183_cli_seam_test.go uses: real components, no stand-in for any of
// them. C5 is the real openai-chat adapter replaying the frozen fixture through
// an httptest server (loop_approval_test.go:60-85's wiring, SPEC-05 §9: no
// function mocks at the C5 seam). The tool surface is the loop's own default
// (agent.NewEchoProvider), declared L0 so the fixture's one echo call really
// dispatches. The roster + spill pair is wired the way cmd/wisp/run.go:1107-1110
// wires it, because the consumption half of this pin IS that wiring.
func r236r4Carrier(t *testing.T) (*agent.Loop, *agent.EchoProvider, *TaskRoster, *agent.Spiller) {
	t.Helper()
	fixture := filepath.Join("..", "agent", "testdata", "golden", "tool-then-text.sse")
	srv, _, err := golden.ServerFromFile(fixture)
	if err != nil {
		t.Fatalf("冻结夹具 %s 读不回来：%v —— 这一发量的是真 C5 seam，夹具不在就没有对象可量", fixture, err)
	}
	t.Cleanup(srv.Close)
	prov, err := llm.NewProvider(llm.EndpointOptions{
		Provider: "mock", Model: "236r4-mock", Protocol: "openai-chat",
		BaseURL: srv.URL, ContextWindow: 128000,
		HTTPClient: &http.Client{Transport: &http.Transport{Proxy: nil}},
	})
	if err != nil {
		t.Fatalf("provider: %v", err)
	}
	echo := agent.NewEchoProvider(agent.ToolInfo{
		Name: "echo", Description: "回显参数（本钉只借一次真派发看三个 id 并排）",
		Parameters: json.RawMessage(`{"type":"object","properties":{"text":{"type":"string"}},"required":["text"]}`),
		Resident:   true, RiskLevel: agent.RiskL0,
	})
	dir := t.TempDir()
	loop, err := agent.New(agent.Options{
		Provider: prov,
		Tools:    echo,
		Config: agent.Config{
			Model: "236r4-mock", ContextWindow: 128000,
			ArtifactsDir:   filepath.Join(dir, "artifacts"),
			PerToolTimeout: 10 * time.Second,
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	spills := agent.NewSpiller(filepath.Join(dir, "artifacts"), loop.Budgets())
	return loop, echo, NewTaskRoster(), spills
}

// Test236R4RunningTaskIDAndResultTaskIDAreOneHostMintedKey is the pin: one real
// run, then the three segments of the session-shaped lock over the ids it yields.
func Test236R4RunningTaskIDAndResultTaskIDAreOneHostMintedKey(t *testing.T) {
	loop, echo, roster, spills := r236r4Carrier(t)

	// The mint happens inside RunAsync, before any model interaction exists
	// (loop.go:329), and this is the only exported reading of it.
	bg := loop.RunAsync(t.Context(), "外面天气怎么样？")
	if bg.ID == "" {
		t.Fatal("RunAsync 回来的 RunningTask.ID 是空串：internal/agent/loop.go:329 的铸造点没有给出 id，下面每条断言都失去对象")
	}
	// Segment 1: the value at the mint must be the shape newTaskID mints.
	if !r236r4MintedShape(bg.ID) {
		t.Errorf("RunningTask.ID = %q 不是宿主铸的形状（36 字符、8/13/18/23 位是连字符、其余小写 hex）："+
			"internal/agent/loop.go:329 的 id := newTaskID() 填的就是这一格，形状由 :1153 决定。"+
			"把来源换成模型可给的值时这一句必须红", bg.ID)
	}

	res := bg.Wait()
	if res.Status != agent.StatusCompleted {
		t.Fatalf("run status = %s（%s）：这一发要一次真跑完的任务，否则 id 的读数不可归因", res.Status, res.Message)
	}

	// Segment 3, the expensive direction: the id the REAL mint produced is the
	// SAME string at both places the ticket named. This is the assertion that
	// had no ruler before this file.
	if res.TaskID != bg.ID {
		t.Errorf("Result.TaskID = %q 而 RunningTask.ID = %q：这两处必须是同一枚宿主铸的键。"+
			"铸造与传递在 internal/agent/loop.go:329（id := newTaskID()）与 :333（同一个 id 交给 run），"+
			"结果在 :382 落进 Result.TaskID；消费在 cmd/wisp/run.go:1105（MarkRoot(bg.ID)）与 "+
			":1110/:1124/:1135/:1145（回填、task_log、状态行全部用 res.TaskID）。"+
			"两处分叉即「名册登记的键」和「输出写进的键」不是同一枚", res.TaskID, bg.ID)
	}

	// The seam where a host id and a model id sit side by side (loop.go:676:
	// TaskID: taskID, CallID: p.call.ID). Anti-vacuity first: the run must
	// really have carried the model's own id, or the inequality below proves
	// nothing.
	calls := echo.Calls()
	if len(calls) != 1 {
		t.Fatalf("precondition: 夹具 tool-then-text.sse 要恰一次 echo 派发，got %d —— 派发次数不是 1 时下面的接缝读数没有对象", len(calls))
	}
	req := calls[0].Req
	if req.CallID != r236r4ModelCallID {
		t.Fatalf("precondition: 接缝上模型给的 call id = %q, want 夹具里那枚 %q —— 夹具一改，这一句读到的是别的世界，必须一起改", req.CallID, r236r4ModelCallID)
	}
	if req.TaskID != bg.ID {
		t.Errorf("工具接缝上 TaskID = %q，不是环路铸的那枚 RunningTask.ID = %q："+
			"internal/agent/loop.go:676 把 run 收到的 taskID 填进 ToolRequest.TaskID，"+
			"它与 cmd/wisp/run.go:1105/:1110 后来当名册键、回填键用的是同一个值", req.TaskID, bg.ID)
	}
	if req.TaskID == req.CallID {
		t.Errorf("接缝上 TaskID 与 CallID 同值（%q）：宿主铸的键与模型给的键合流了，"+
			"「键的来源」这一维当场失效（loop.go:676 的两格必须各是各的）", req.TaskID)
	}

	// The consumption the composition root really performs (run.go:1105 then
	// :1110): if the two ids ever diverge, the roster gains a second row under
	// the result's key and the root row stops carrying the text — so "两者相等"
	// has an observable consequence, not just a string comparison.
	roster.MarkRoot(bg.ID, "236r4 根任务")
	rec, why := (TaskBackfill{Roster: roster, Spills: spills}).
		Backfill(bg.Root().Ctx, res.TaskID, res.Text)
	if why != "" {
		t.Fatalf("回填被拒（%s）：这条链只在该键不是铸来的、或任务已取消时才拒，两种都说明「键的来源」这一维断了", why)
	}
	row, ok := roster.Look(bg.ID)
	if !ok {
		t.Fatalf("名册里没有以 RunningTask.ID（%q）为键的行：cmd/wisp/run.go:1105 用 bg.ID 登的行不见了", bg.ID)
	}
	if rec.Text != res.Text || row.Text != res.Text {
		t.Errorf("以铸来的键读回的正文 = %q / %q，不是 Result 的 %q：输出写的键与名册登的键分叉了（res.TaskID %q vs RunningTask.ID %q）",
			rec.Text, row.Text, res.Text, res.TaskID, bg.ID)
	}
	if row.Kind != TaskKindRoot {
		t.Errorf("以铸来的键读回的行 Kind = %q, want %q：MarkRoot 的行和回填的行落到两枚键上了",
			row.Kind, TaskKindRoot)
	}
	if got := roster.Count(); got != 1 {
		t.Errorf("名册有 %d 行, want 1：MarkRoot(bg.ID) 与 Backfill(res.TaskID) 一旦分叉就会多出一行，"+
			"这正是「两处相等」的可观察后果", got)
	}

	// Segment 2, the cheap direction: the model-side vocabulary and this repo's
	// hand-typed fixture ids are outside the lock, and this run's mint is not
	// one of them. Enumerated, not structural — see the file header's ceiling.
	for _, lit := range r236r4ModelWireIDs {
		if r236r4MintedShape(lit) {
			t.Errorf("模型侧 id %q 通过了宿主铸的形状锁：锁不再把模型给的值排除在外", lit)
		}
		if res.TaskID == lit {
			t.Errorf("铸出来的 id 等于模型侧字面量 %q：模型给的值冒充了宿主铸的键", lit)
		}
	}
	for _, lit := range r236r4FixtureTaskIDs {
		if r236r4MintedShape(lit) {
			t.Errorf("台件手写的 task id %q 落在宿主铸的形状里：台件与生产铸点就此可以相撞", lit)
		}
	}
}
