package tools

// 票 287 AC#2 — the fixture that makes internal/tools/subagent_197.go's ROSTER depth
// check reachable, and the three legs that read it. AC#1's rulers (.scratch/wisp/
// probes/287/a1/00-three-rulers-re-run.md) measured that production cannot reach it:
// the spawner hands the child newSubagentToolProvider (:316, my line) and its name gate
// refuses task.spawn before the bridge is asked (:604-610). 这一支不是缺陷：深度 1 是
// 设计，缺的是能让它走到的仪器。The fixture hands the caller's loop the bridge UNWRAPPED
// - the one shape where :276-280 alone stands between a subagent and a second generation
// - and ASSERTS AC#2's three ingredients: the caller row filed Kind == TaskKindSubagent
// by the real PublishSubagent from an agent.Loop AdmitTask hook (:333, production's own
// site); a REAL loop, so corr = callCorr(...) != taskID (票 282/283's precondition); the
// task.spawn request. MaxSubagentDepth stays 1; injection = the frozen C5 seam only.

import (
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// refuse287Tail: the depth branch's own sentence, at ZERO hits in any *_test.go (AC#1).
const (
	refuse287Tail    = "子代理不许再派子代理。"
	refuse287HeadFmt = "拒绝派生：深度上限是 "
	poolRefuse287    = "拒绝派生：同时在跑的子代理已达上限"
	callerSub287     = "287 号调用者：名册行本身就是子代理"
	callerRoot287    = "287 号调用者：宿主自己的根任务（反形正控）"
	grand287Answer   = "孙代理的结论（深度 1 的树里不该出现这一行）"
)

var (
	re287Who   = regexp.MustCompile(`派生者 (\S+) 自己就是一枚子代理`)
	re287Kind  = regexp.MustCompile(`父 kind=([^）]+)）`)
	re287Child = regexp.MustCompile(`子代理 (\S+?)（`)
)

// new287Harness wires the real bridge, roster and journal; provider() is the child the
// control leg may derive, and it never asks for a spawn itself.
func new287Harness(t *testing.T) *harness221 {
	t.Helper()
	return build221(t, windowGate221(), func() llm.LlmProvider {
		return &fake197Provider{name: "grandchild287", answer: grand287Answer}
	}, true)
}

// runCallerLoop287 runs ONE real agent.Loop whose first model turn asks for task.spawn,
// filing the loop's own task id in the roster as asKind before that turn can be
// dispatched. The id is minted by the loop, never by this fixture.
func runCallerLoop287(t *testing.T, x *harness221, asKind string) (agent.Result, string) {
	t.Helper()
	var (
		mu       sync.Mutex
		admits   int
		callerID string
	)
	loop, err := agent.New(agent.Options{
		Provider: &fake197Provider{name: "caller287", answer: "287 号调用者收工", askSpawnFirstTurn: true},
		Tools:    x.h.bridge, // UNWRAPPED, on purpose: see the header.
		Sink:     &agent.RecordSink{},
		Registry: observe.NewRegistry(),
		AdmitTask: func(taskID string) func() {
			mu.Lock()
			admits++
			callerID = taskID
			mu.Unlock()
			// 前置一：走产码自己那枚 API、产码自己那个发布点。
			if asKind == TaskKindSubagent {
				x.h.roster.PublishSubagent(taskID, parent197, callerSub287)
			} else {
				x.h.roster.MarkRoot(taskID, callerRoot287)
			}
			return func() {}
		},
		Config: agent.Config{
			Model: "probe-287", ContextWindow: 8192,
			ArtifactsDir:                filepath.Join(sealableTempDir124(t), "artifacts"),
			PerToolTimeout:              10 * time.Second,
			PassThroughUnclassifiedRisk: false,
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	res := loop.Run(t.Context(), "自己就是子代理，还要再派一枚子代理")
	mu.Lock()
	defer mu.Unlock()
	if admits != 1 || callerID != res.TaskID {
		t.Fatalf("夹具自己没站住：AdmitTask 跑了 %d 次，callerID=%q, res.TaskID=%q", admits, callerID, res.TaskID)
	}
	return res, callerID
}

// spawnLog287 requires the loop to have run exactly one task.spawn and hands back that
// call's log line: the text the model actually got.
func spawnLog287(t *testing.T, res agent.Result) agent.ToolResultLog {
	t.Helper()
	if res.Status != agent.StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if res.ToolCalls != 1 {
		t.Fatalf("ToolCalls = %d, want 1（第一程只该发一枚 task.spawn）：%+v", res.ToolCalls, res.ToolLog)
	}
	for _, l := range res.ToolLog {
		if l.Name == "task.spawn" {
			return l
		}
	}
	t.Fatalf("工具回执里没有 task.spawn 这一发（前置三没跑到）：%+v", res.ToolLog)
	return agent.ToolResultLog{}
}

// corrOf287 reads the persisted tool_call row of one call and returns its correlation
// id, after requiring the real loop's per-call shape. Without it a leg could be reading
// a direct bridge.Execute fixture, where corr collapses onto the task id and the
// identity question has no answer（前置二）.
func corrOf287(t *testing.T, x *harness221, taskID, tool string) string {
	t.Helper()
	rows, err := x.mem.ListToolCallsByTask(t.Context(), taskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	hits := make([]memory.ToolCall, 0, len(rows))
	for _, r := range rows {
		if r.Tool == tool {
			hits = append(hits, r)
		}
	}
	if len(hits) != 1 {
		t.Fatalf("%s 的 tool_call 行数 = %d, want 1（真 loop 每枚调用一行）：%+v", tool, len(hits), rows)
	}
	corr := hits[0].CorrelationID
	if corr == "" || corr == taskID || !strings.HasPrefix(corr, taskID) {
		t.Errorf("%s 的行 corr = %q, want 非空、以 task id %q 为前缀且不等于它（派发没经真 loop）", tool, corr, taskID)
	}
	return corr
}

func rowIDs287(rows []TaskRecord) []string {
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.TaskID)
	}
	return ids
}

// 主用例：调用者那一行是 subagent ⇒ 名册深度判定那一支必须自己拦住。
func Test287SubagentCallerRefusedAtTheRosterDepthCheckThroughTheRealLoop(t *testing.T) {
	x := new287Harness(t)
	res, callerID := runCallerLoop287(t, x, TaskKindSubagent)
	rec, ok := x.h.roster.Look(callerID)
	if !ok || rec.Kind != TaskKindSubagent {
		t.Fatalf("前置一没落地：调用者 %s 的行 ok=%v Kind=%q, want %q", callerID, ok, rec.Kind, TaskKindSubagent)
	}
	log := spawnLog287(t, res)
	// 拒绝发生在 Execute 里面的名册判定，不是闸门：Rejected 只由执行前的拦下置位。
	if log.Rejected {
		t.Fatalf("task.spawn 被闸门拦下了（gate 读数，不是名册判定读数）：%+v", log)
	}
	// 判据一：拒语整句的两枚承重点（:278-279 现量号，正文照盘上）。
	if !strings.HasPrefix(log.Text, refuse287HeadFmt) {
		t.Errorf("subagent_197.go:278 那句「%s%d」没有出现在 task.spawn 的回执里＝深度判定那一支没跑到"+
			"（种坏 M1/M2 就是这个读数形）：text=%q", refuse287HeadFmt, MaxSubagentDepth, log.Text)
	}
	if !strings.Contains(log.Text, refuse287Tail) {
		t.Errorf("subagent_197.go:279 那句 %q（AC#1 尺①数出来全部 *_test.go 零命中）没被说出口：text=%q", refuse287Tail, log.Text)
	}
	// 判据二＝身份那一维：两枚 id 同屏，才写得出"读错了会红"的那一句。
	corr := corrOf287(t, x, res.TaskID, "task.spawn")
	who := re287Who.FindStringSubmatch(log.Text)
	if who == nil {
		t.Errorf(":278 那句里的「派生者 <task id>」抽不出来＝身份这一维没有读数：text=%q", log.Text)
	} else if who[1] != res.TaskID {
		t.Errorf("深度判定读的身份 = %q, want task id %q；这一发的 per-call corr = %q（corr != taskID）"+
			"：把 :271 改读 CorrelationID(ctx) 就是这个读数", who[1], res.TaskID, corr)
	}
	if m := re287Kind.FindStringSubmatch(log.Text); m == nil || m[1] != TaskKindSubagent {
		t.Errorf("父 kind 那一枚 = %v, want %q（判定读的不是调用者自己的子代理行）：text=%q", m, TaskKindSubagent, log.Text)
	}
	// 判据三＝后果，不只是句子：拒绝之后第二代不许存在。
	if grand := x.h.roster.Descendants(callerID); len(grand) != 0 {
		t.Errorf("名册里出现第二代（%d 枚 %v）＝深度上限 %d 被走到并破掉；调用者=%s", len(grand), rowIDs287(grand), MaxSubagentDepth, callerID)
	}
	if kids := x.h.roster.Descendants(parent197); len(kids) != 1 {
		t.Errorf("parent197 的孩子行 = %d 枚 %v, want 1（夹具只登记调用者自己那一枚）", len(kids), rowIDs287(kids))
	}
	if log.Outcome != agent.OutcomeError {
		t.Errorf("task.spawn 的 outcome = %q, want %q：被判定拒掉的那一发必须是错误回执，模型才改得动", log.Outcome, agent.OutcomeError)
	}
	if n := x.h.roster.InFlightSubagents(); n != 0 {
		t.Errorf("深度拒绝发生之后还占着 %d 枚池位：那一支必须在 TryAcquireSubagentSlot 之前返回", n)
	}
}

// 反形正控：同一形状、调用者那一行是 root ⇒ 不许拒（打掉恒真检测器）。
func Test287RootCallerWithTheSameShapeStillDerives(t *testing.T) {
	x := new287Harness(t)
	res, callerID := runCallerLoop287(t, x, TaskKindRoot)
	rec, ok := x.h.roster.Look(callerID)
	if !ok || rec.Kind != TaskKindRoot {
		t.Fatalf("反形前置没落地：调用者 %s 的行 ok=%v Kind=%q, want %q", callerID, ok, rec.Kind, TaskKindRoot)
	}
	log := spawnLog287(t, res)
	for _, refusal := range []string{refuse287HeadFmt, refuse287Tail} {
		if strings.Contains(log.Text, refusal) {
			t.Errorf("父 kind=root 的调用者被深度判定拒了（那一支已经宽到认不出 kind）：refusal=%q 出现在 text=%q", refusal, log.Text)
		}
	}
	if log.Outcome != agent.OutcomeSuccess {
		t.Errorf("task.spawn 的 outcome = %q, want %q：正控必须真的派出来一枚", log.Outcome, agent.OutcomeSuccess)
	}
	kids := x.h.roster.Descendants(callerID)
	if len(kids) != 1 {
		t.Fatalf("孩子行 = %d 枚 %v, want 1：%+v", len(kids), rowIDs287(kids), log)
	}
	if kids[0].Out.Kind != TaskKindSubagent {
		t.Errorf("派生出来的那一行 Kind = %q, want %q", kids[0].Out.Kind, TaskKindSubagent)
	}
	if m := re287Child.FindStringSubmatch(log.Text); m == nil || m[1] != kids[0].TaskID {
		t.Errorf("回执里报的孩子 id = %v, want %q（派生与回执不是同一枚）", m, kids[0].TaskID)
	}
	corrOf287(t, x, res.TaskID, "task.spawn")
}

// 次序尺：深度判定在前，池位判定在后（两枚拒绝同屏，读哪一枚才算数）。
func Test287DepthRefusalWinsOverThePoolRefusal(t *testing.T) {
	x := new287Harness(t)
	// 把池占满（产码自己那枚 API），让"派生不了"有两个可能的理由摆在同一屏上。
	releases := make([]func(), 0, MaxConcurrentSubagents)
	defer func() {
		for _, release := range releases {
			release()
		}
	}()
	for i := 0; i < MaxConcurrentSubagents; i++ {
		release, ok := x.h.roster.TryAcquireSubagentSlot(MaxConcurrentSubagents)
		if !ok {
			t.Fatalf("第 %d 枚池位没占上（want %d 枚全占）", i+1, MaxConcurrentSubagents)
		}
		releases = append(releases, release)
	}
	if n := x.h.roster.InFlightSubagents(); n != MaxConcurrentSubagents {
		t.Fatalf("池位读数 = %d, want %d：这一发的分母不是它声称的那个", n, MaxConcurrentSubagents)
	}
	res, _ := runCallerLoop287(t, x, TaskKindSubagent)
	log := spawnLog287(t, res)
	if !strings.HasPrefix(log.Text, refuse287HeadFmt) {
		t.Errorf("深度判定没有在池位判定之前拦住（池已满，回执是另一句）：text=%q", log.Text)
	}
	if strings.Contains(log.Text, poolRefuse287) {
		t.Errorf("调用者自己是子代理，却先读到池满那句＝:276 那一支没跑：%q", log.Text)
	}
}
