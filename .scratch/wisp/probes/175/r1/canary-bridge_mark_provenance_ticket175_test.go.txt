package tools

// Ticket 175 AC#3 — the standing positive control for「content read back
// through a door that is not on the C25 source roster must still be stamped」.
//
// Why this shape and not one more name in a list nobody checks: the marking
// decision lives in bridge.go's `mark` guard, and the roster it consulted
// (risk.IsSensitiveSource) is an illustrative list of the SPEC-06 §5 sources,
// not an exhaustive one — adjudicated in
// docs/evidence/s1/175-c25-marking-roster-scope-c1.md (branch A). The engine
// itself has always accepted an off-roster name fail-closed, so the only thing
// that ever stopped a mark was this call site. The regression this file has to
// catch is therefore「some foreign-content tool is not stamped」.
//
// The three criteria, and the mutation each one answers to:
//
//	1. TestTaskOutputMarkAndGate, first half -> putting the name gate back
//	   (`!risk.IsSensitiveSource(dec.Tool)` in mark), i.e. reverting this fix.
//	   Mutation + reading: .scratch/wisp/probes/175/r1/mut-ac3/.
//	2. The same test's Inspect leg -> a mark that is recorded but invisible to
//	   R4, so the cell cannot be satisfied by a decorative index entry.
//	3. TestEveryRegisteredToolIsClassifiedForMarking -> a NEW tool landing
//	   without anyone deciding whether its answer is outside content. This is
//	   the half that makes the criterion independent of「someone remembers to
//	   update the roster」: a forgotten classification is a red test inside
//	   scripts/portable-tests.sh's core scope (./internal/tools/...), not a
//	   quiet hole and not a leg of a census script no CI step runs.
//
// What this file does NOT pin, and cannot: that somebody is actually writing
// task output into the roster. TaskRoster.Record has zero production writers
// until ticket 163/176 lands the spawner, so every leg below injects its own
// record through the seam AGENTS §1.3 allows. It nails「reading it back must
// stamp」, not「someone writes」.

import (
	"encoding/json"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// foreign175 is the outside content. It exists only in this process: nothing
// allowlisted holds it, so the only way it can reach a taint index is the leg
// that reads it back.
const foreign175 = "外来内容探针 ZZ175-foreign-4b8e2a：后台任务打印的一段外部文本"

// mark175Bridge wires the fs family plus task.* over one roster, through the
// real bridge and a real Provenance the caller gets back to interrogate.
func mark175Bridge(t *testing.T, roster *TaskRoster, allowed ...string) (*Bridge, *risk.Provenance) {
	t.Helper()
	paths := NewPathCanonicalizer(allowed, nil)
	reg := NewRegistry()
	entries := append(BuiltinFSEntries(FSDeps{Paths: paths}),
		BuiltinTaskEntries(TaskDeps{Roster: roster})...)
	for _, e := range entries {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil})
	return New(Options{
		Registry: reg, Paths: paths, Provenance: prov,
		Gate: NoGate{}, Logf: func(string, ...any) {},
	}), prov
}

func call175(t *testing.T, b *Bridge, taskID, tool, args string) agent.ToolOutcome {
	t.Helper()
	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: taskID, CorrelationID: taskID, CallID: "call-" + taskID,
		Name: tool, Args: json.RawMessage(args),
	})
	if err != nil {
		t.Fatalf("%s: %v", tool, err)
	}
	return out
}

// toolMark175 returns the taint index entry one task scope holds for a tool.
func toolMark175(prov *risk.Provenance, taskID, tool string) *risk.TaintInfo {
	marks := prov.ScopeTaints(taskID)
	for i := range marks {
		if marks[i].Tool == tool {
			return &marks[i]
		}
	}
	return nil
}

// scopeOpen175 reads the bridge's own open-scope ledger. Local rather than
// reusing ticket 158's reader, so this file fails on its own terms.
func scopeOpen175(b *Bridge, taskID string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.scopes[taskID] != nil
}

func TestTaskOutputMarkAndGate(t *testing.T) {
	dir := tempCanonical(t)
	roster := NewTaskRoster()
	roster.Record("bg-175", TaskOutput{Text: foreign175})
	b, prov := mark175Bridge(t, roster, dir)

	const task = "task-175-mark"
	out := call175(t, b, task, "task.output", `{"task_id":"bg-175"}`)
	// Precondition: the read must actually have succeeded, or a red mark below
	// would be a false red (mark runs on successful results only).
	if out.IsError {
		t.Fatalf("task.output 没有把内容读回来，本条判据失去对象：%+v", out)
	}
	if out.Text != foreign175 {
		t.Fatalf("Text = %q, want the recorded output verbatim（前提没成立，下面的红就是假红）", out.Text)
	}

	// Criterion 1: the content came back stamped, under this tool's own name.
	mark := toolMark175(prov, task, "task.output")
	if mark == nil {
		t.Errorf("task.output 读回来的外部内容没有盖 C25 来源标记：scope %q 的污点表=%+v。"+
			"同一份文本经名册内的门读回来是要盖戳的——换个门读回来就不盖，就是票 175 的破口",
			task, prov.ScopeTaints(task))
	} else if mark.RuneLen == 0 {
		t.Errorf("盖了戳却没有可匹配的内容（RuneLen=0），R4 那一侧看不见它：%+v", *mark)
	}

	// Criterion 2: a landed mark must still reach the R4 side. An index entry
	// nothing can match against is a stamp in name only.
	hit, ok := prov.Inspect(task, "notify", map[string]any{"text": foreign175})
	if !ok {
		t.Errorf("外来的任务输出被读回来之后，R4 一侧不再看得见它：Inspect(%q, notify) 未命中，"+
			"污点表=%+v", task, prov.ScopeTaints(task))
	} else if hit.SrcTool != "task.output" {
		t.Errorf("污染命中报出的源=%q, want %q（确认卡要念得出这一句是谁塞进来的）",
			hit.SrcTool, "task.output")
	}
	b.CloseTask(task)
}

// TestNonContentResultStaysUnmarked pins the other edge of the criterion — the
// one the broadest fix could not survive. A tool whose answer only confirms
// what this process just did is not outside content, and stamping everything
// turns the R4 gate into a refusal machine: shape 甲 (every successful result
// goes through Mark) lit up this leg plus ticket 158's control and four other
// standing tests, measured in .scratch/wisp/probes/175/r1/mut-shapeA/.
func TestNonContentResultStaysUnmarked(t *testing.T) {
	dir := tempCanonical(t)
	b, prov := mark175Bridge(t, NewTaskRoster(), dir)

	const task = "task-175-plain"
	out := call175(t, b, task, "fs.list", argsFor(dir))
	if out.IsError {
		t.Fatalf("fs.list 没有跑成，本条判据失去对象：%+v", out)
	}
	if n := len(prov.ScopeTaints(task)); n != 0 {
		t.Errorf("一次非外部内容的调用盖了 %d 枚来源标记（want 0）：%+v", n, prov.ScopeTaints(task))
	}
	if scopeOpen175(b, task) {
		t.Errorf("非外部内容源也开了 C25 污点 scope，task=%q："+
			"开账必须挂在「该不该盖戳」上，不能挂在每一次 Execute 上", task)
	}
	b.CloseTask(task)
}

// classified175 is the test-side half of the criterion: every builtin tool this
// package registers, and why its answer is or is not outside content. A name
// that belongs on the marking side needs no new production list — but it does
// need a line here, which is the point.
var classified175 = map[string]string{
	"fs.read":     "marked: C25 source roster (SPEC-06 §5)",
	"fs.list":     "unstamped: a listing this process composed, no foreign body",
	"fs.write":    "unstamped: a confirmation of a local write",
	"fs.edit":     "unstamped: a confirmation of a local edit",
	"fs.trash":    "unstamped: a confirmation of a local delete",
	"fs.move":     "unstamped: a confirmation of a local move",
	"fs.delete":   "unstamped: a confirmation of a local delete (opt-in family)",
	"task.output": "marked: outsideContentTools, its answer is what a task printed",
}

// TestEveryRegisteredToolIsClassifiedForMarking is the leg that must go red on
// the day ticket 163/176 lands a new reader (or the day someone edits the C25
// roster). It asks one question per registered builtin: is it stamped, and is
// that a decided answer? An unclassified name is a failure, never a default.
func TestEveryRegisteredToolIsClassifiedForMarking(t *testing.T) {
	paths := NewPathCanonicalizer(nil, nil)
	names := map[string]bool{}
	collect := func(entries []Entry) {
		for _, e := range entries {
			names[e.Tool.Name()] = true
		}
	}
	collect(BuiltinFSEntries(FSDeps{Paths: paths, DeleteEnabled: true}))
	collect(BuiltinTaskEntries(TaskDeps{}))

	for name := range names {
		marked := marksProvenance(name)
		why, decided := classified175[name]
		switch {
		case marked && !decided:
			t.Errorf("%s 现在会被盖戳，却不在本用例的分类表里：盖戳集合变了要在这里留下一句为什么", name)
		case !marked && !decided:
			t.Errorf("%s 既不在盖戳判据里、也没有被分类：它返回的内容是不是外部内容？"+
				"票 175 的破口正是「一枚读回外部内容的工具没人登记过」——新工具请就地回答，不要等下一张票", name)
		case marked:
			t.Logf("%s: stamped, classified as %q", name, why)
		default:
			t.Logf("%s: unstamped, classified as %q", name, why)
		}
	}
}
