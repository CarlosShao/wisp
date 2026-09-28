package tools

// Ticket 183-r2 AC#2b + AC#4 - the permanent CLI-seam half that AC#2 lacked.
//
// Why these legs sit here and not one level down: every permanent 183 criterion
// lived in internal/risk (package ruler at the start of this ticket:
// `grep -rln "183" internal/tools/*_test.go` => zero files). The end-to-end
// reading that ticket 183-r1 did get lived only in
// .scratch/wisp/probes/183/r1/**, i.e. a one-shot overlay probe - which is
// precisely the shape tickets 164 and 175 kept re-earning: a green package-level
// suite that does not stop the real CLI. 183-v2 judged AC#2 "partially met" for
// exactly this hole, and named the shape to close it: the pointer the model
// follows must be the one the BACKFILL writer produced (task_backfill.go ->
// TaskRoster -> task.output reads it back through hostPathBoxFromCtx + box.set
// -> bridge.mark(..., hostPath)), not a path a leg injected by hand.
//
// The two criteria, and what turns each one red:
//
//	AC#2b TestPointer183BackfilledArtifactRereadStaysClean
//	    -> the exemption dying on the CLI seam: box.set removed from task.go,
//	       "task.output" taken off the C25 marking roster (ticket 175's gate in
//	       bridge.mark), or the declared string no longer reaching
//	       MarkWithHostPath. The preconditions below are the anti-vacuity halves:
//	       this leg is only measuring 183's family if the roster really handed
//	       back a spilled artifact path, the stub really printed it, and the
//	       marked body really spells a >=8-rune window of that path OUTSIDE the
//	       pointer itself (175-r2's pair pins the carrier and cannot see that
//	       dimension - its payload is an undifferentiated ascii run).
//	AC#4  TestPointer183SiblingArtifactInTheSameDirectoryStillHits
//	    -> widening the exemption from "the one path this call declared" to its
//	       directory or prefix: the sibling lives in the same artifacts dir and
//	       the host's own stub spells it verbatim, so any directory-level rule
//	       launders it and this leg goes red. Ticket 177 W-3's family on the CLI
//	       shape; the host minted path1, never path2.
//
// Deliberately NOT asserted here: the SECOND bounded reread of the same path
// (ticket 185 owns that one, and 183-v2's L4 measured that the blocker is the
// undeclared read-back mark - quieting it from here would wash ticket 183 AC#3),
// and anything that reads the confirmation window's seconds (AC#6 is a
// attribution-only cell; no approval constant is touched by this file).

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

const (
	// r183BackgroundID is the task id the host mints for the background task
	// whose answer this leg re-reads. It is a plain hex id because that is what
	// agent.newTaskID() gives in production (task_backfill.go's doc block names
	// it), and the artifact key is "agent-task-"+id, so the spilled file is
	// named tool-output-agent-task-<id>.txt - the same name family the real CLI
	// of ticket 183 produced.
	r183BackgroundID  = "c66c0634183r2a1b"
	r183SiblingTaskID = "c66c0635183r2b2c"

	// r183FillerLine is the spilled body's repeating line, pinned to the exact
	// payload 179-r2 and 183-r1 used on the real CLI: its "-output-" is one of
	// the artifact path's own 8-rune windows, which is the whole defect family
	// of ticket 183. An undifferentiated payload would let this leg pass on the
	// carrier alone.
	r183FillerLine = "WISP179R2-background-output-line."

	r183BodyBytes = 20000
)

// r183PointerRe reads the pointer the host itself printed inside the stub.
// Local on purpose: task.go owns the stub wording (frozen by PLAN.md:2564) and
// this leg has to fail on its own terms, not on a regex borrowed from another
// file.
var r183PointerRe = regexp.MustCompile(`全文见 (\S+)…`)

// r183Payload returns a body longer than the D15(3) spill budget, so the host
// really files it and really points at the file.
func r183Payload(extraPrefix string) string {
	body := strings.Repeat(r183FillerLine, r183BodyBytes/len(r183FillerLine)+1)[:r183BodyBytes]
	return extraPrefix + body[len(extraPrefix):]
}

// r183Seam wires the roster plus the real D15(3) spill layer the composition
// root hands to TaskBackfill, over one temporary directory that is also the
// allowed root, so the artifact the pointer names is readable through C26.
func r183Seam(t *testing.T) (*TaskRoster, *agent.Spiller, string) {
	t.Helper()
	dir := t.TempDir()
	// BudgetsFor(128000) is the CLI's own 128k row: SpillTokens 4000, so the
	// 20000-byte answer below really spills instead of riding along whole.
	spills := agent.NewSpiller(filepath.Join(dir, "artifacts"), agent.BudgetsFor(128000))
	return NewTaskRoster(), spills, dir
}

// r183Bridge runs fs.read and task.output over one roster through the real
// bridge with the real C25 provenance, and hands the provenance back so a leg
// can interrogate the index instead of trusting a log line. NoGate is the CLI's
// own un-wired-panel shape (票 21 未接入): a judgement still happens, it just
// cannot ask anybody, which is exactly the side that must not be papered over.
func r183Bridge(t *testing.T, roster *TaskRoster, allowed string) (*Bridge, *risk.Provenance) {
	t.Helper()
	paths := NewPathCanonicalizer([]string{allowed}, nil)
	reg := NewRegistry()
	entries := append(BuiltinFSEntries(FSDeps{Paths: paths}),
		BuiltinTaskEntries(TaskDeps{Roster: roster, Paths: paths})...)
	for _, e := range entries {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil})
	return New(Options{
		Registry: reg, Paths: paths, Provenance: prov, Gate: NoGate{},
		Logf: func(string, ...any) {},
	}), prov
}

func r183Call(t *testing.T, b *Bridge, taskID, tool, args string) agent.ToolOutcome {
	t.Helper()
	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: taskID, CorrelationID: taskID, CallID: "call-" + taskID + "-" + tool,
		Name: tool, Args: json.RawMessage(args),
	})
	if err != nil {
		t.Fatalf("%s: %v", tool, err)
	}
	return out
}

// r183Args encodes one fs.read call for a Windows path, backslashes and all.
func r183Args(t *testing.T, path string, maxBytes int) string {
	t.Helper()
	if maxBytes > 0 {
		return `{"path":"` + jsonEscape(path) + `","max_bytes":` + strconv.Itoa(maxBytes) + `}`
	}
	return `{"path":"` + jsonEscape(path) + `"}`
}

// r183SpelledTwice reports whether `needle` occurs in hay more times than the
// pointer occurrence alone accounts for - the twin-fragment precondition that
// separates this family from 175-r2's carrier-only pair.
func r183SpelledOutsidePointer(hay, pointer, needle string) (int, bool) {
	total := strings.Count(hay, needle)
	inPointer := strings.Count(pointer, needle)
	if extra := total - inPointer; extra > 0 {
		return extra, true
	}
	return 0, false
}

// ---------------------------------------------------------------------------
// AC#2b - the model follows the pointer the host wrote for a BACKFILLED
// artifact and reads every byte of it back, on a permanent leg, on the CLI seam.
// ---------------------------------------------------------------------------

func TestPointer183BackfilledArtifactRereadStaysClean(t *testing.T) {
	roster, spills, dir := r183Seam(t)
	b, prov := r183Bridge(t, roster, dir)
	const task = "183r2-cli"

	// The WRITE side, through the writer the ticket named: no hand-made record.
	rec, why := TaskBackfill{Roster: roster, Spills: spills}.
		Backfill(context.Background(), r183BackgroundID, r183Payload(""))
	if why != "" {
		t.Fatalf("precondition: 回填没写进名册：%s", why)
	}
	if rec.ArtifactPath == "" {
		t.Fatal("precondition: 超长正文没有落盘，这一发测不到回填产生的那条 ArtifactPath")
	}
	if name := filepath.Base(rec.ArtifactPath); !strings.HasPrefix(name, "tool-output-agent-task-") {
		t.Fatalf("precondition: 产物名不是 CLI 那一族（got %q），twin 片段这一维就不成立了", name)
	}
	want, err := os.ReadFile(rec.ArtifactPath)
	if err != nil || len(want) != r183BodyBytes {
		t.Fatalf("precondition: 名册登记的路径读不回来：%v (%d bytes)", err, len(want))
	}

	// task.output reads that record back on the spot and prints the pointer.
	stub := r183Call(t, b, task, "task.output", `{"task_id":"`+r183BackgroundID+`"}`)
	if stub.IsError || !stub.Truncated {
		t.Fatalf("precondition: 需要一个带指针的截断桩，got %+v", stub)
	}
	m := r183PointerRe.FindStringSubmatch(stub.Text)
	if m == nil {
		t.Fatalf("precondition: 桩里没有宿主写下的指针：%q", tail183(stub.Text))
	}
	pointer := m[1]
	if p := filepath.ToSlash(rec.ArtifactPath); pointer != rec.ArtifactPath && pointer != p {
		t.Errorf("模型续读用的路径 %q 不是名册登记的那条 %q", pointer, rec.ArtifactPath)
	}
	if extra, ok := r183SpelledOutsidePointer(stub.Text, pointer, "-output-"); !ok {
		t.Fatalf("precondition: 正文里没有指针之外的 %q（总共 %d 次），这一发退化成票 175-r2 的载具腿",
			"-output-", strings.Count(stub.Text, "-output-")+extra)
	}
	if mark := prov.ScopeTaints(task); len(mark) == 0 {
		t.Fatal("precondition: task.output 没盖戳，绿色就不是豁免给的（票 175 的破口本身）")
	}

	// The claim: following the host's own pointer is not R4 evidence.
	got := r183Call(t, b, task, "fs.read", r183Args(t, pointer, 0))
	if got.RiskLevel != "L0" {
		t.Errorf("照宿主指针的续读判成 %q, want L0（R4 又咬住宿主自己写下的指针了）\n桩：%s",
			got.RiskLevel, tail183(stub.Text))
	}
	if got.IsError {
		t.Errorf("续读那一发被拒了：%+v", got)
	} else if len(got.Text) != len(want) {
		t.Errorf("续读没逐字节回全：产物 %d 字节 / 读回 %d 字节", len(want), len(got.Text))
	} else if got.Text != string(want) {
		t.Errorf("续读回来的字节与产物不相等（前 80 字节 %q vs %q）", got.Text[:80], string(want[:80]))
	}
	b.CloseTask(task)
}

// ---------------------------------------------------------------------------
// AC#4 - ticket 177 W-3's family on the CLI shape: a SIBLING artifact in the
// very same directory, spelled verbatim by the host's own stub, that the host
// never declared. Directory/prefix-level exemptions die here.
// ---------------------------------------------------------------------------

func TestPointer183SiblingArtifactInTheSameDirectoryStillHits(t *testing.T) {
	roster, spills, dir := r183Seam(t)
	b, prov := r183Bridge(t, roster, dir)
	const task = "183r2-sibling"

	// Two real background tasks, so the sibling is not an invention of this leg:
	// both artifacts really landed, in the same artifacts directory.
	sibRec, sibWhy := TaskBackfill{Roster: roster, Spills: spills}.
		Backfill(context.Background(), r183SiblingTaskID, r183Payload(""))
	if sibWhy != "" || sibRec.ArtifactPath == "" {
		t.Fatalf("precondition: 兄弟产物没落盘：%q", sibWhy)
	}
	// The host's own answer text mentions the sibling path, so the stub the model
	// reads spells it - while the only path this call DECLARES is its own.
	mention := "上一次运行也落过一份同类产物 " + filepath.ToSlash(sibRec.ArtifactPath) + " 在同一目录下\n"
	rec, why := TaskBackfill{Roster: roster, Spills: spills}.
		Backfill(context.Background(), r183BackgroundID, r183Payload(mention))
	if why != "" || rec.ArtifactPath == "" {
		t.Fatalf("precondition: 声明方产物没落盘：%q", why)
	}
	if filepath.Dir(rec.ArtifactPath) != filepath.Dir(sibRec.ArtifactPath) {
		t.Fatalf("precondition: 两枚产物必须同目录（got %q / %q）", rec.ArtifactPath, sibRec.ArtifactPath)
	}
	if rec.ArtifactPath == sibRec.ArtifactPath {
		t.Fatal("precondition: 兄弟产物与声明产物同一条，这一发测不到兄弟形状")
	}

	stub := r183Call(t, b, task, "task.output", `{"task_id":"`+r183BackgroundID+`"}`)
	if stub.IsError || !stub.Truncated {
		t.Fatalf("precondition: 需要一个带指针的截断桩，got %+v", stub)
	}
	if m := r183PointerRe.FindStringSubmatch(stub.Text); m == nil || m[1] != rec.ArtifactPath {
		t.Fatalf("precondition: 桩里的指针不是名册登记那条：%+v vs %q", m, rec.ArtifactPath)
	}
	sibling := filepath.ToSlash(sibRec.ArtifactPath)
	if !strings.Contains(stub.Text, sibling) {
		t.Fatalf("precondition: 宿主正文没逐字拼出兄弟路径，这一发就测不到 W-3 那一族（桩：%s）", tail183(stub.Text))
	}

	// The control half, so this cannot pass on "block everything": the declared
	// pointer of the very same stub stays re-readable.
	if declared := r183Call(t, b, task, "fs.read", r183Args(t, rec.ArtifactPath, 0)); declared.RiskLevel != "L0" {
		t.Errorf("同一发桩里宿主自己声明的那条路径判成 %q, want L0：豁免连声明方都不放了", declared.RiskLevel)
	}

	// And the sibling, spelled verbatim by the same host-written stub but never
	// declared, must still be evidence.
	sib := r183Call(t, b, task, "fs.read", r183Args(t, sibling, 0))
	if sib.RiskLevel == "L0" {
		t.Errorf("同目录下那枚兄弟产物判成 L0：豁免从「这一条路径」扩到了「这个目录」，票 177 W-3 那一族被洗掉了\n声明方：%s\n兄弟：%s", rec.ArtifactPath, sibRec.ArtifactPath)
	}
	h, ok := prov.Inspect(task, "fs.read", map[string]any{"path": sibling})
	if !ok {
		t.Fatalf("兄弟路径必须在 C25 里命中（Inspect 没命中＝豁免按目录生效了）：%+v", prov.ScopeTaints(task))
	}
	if h.SrcTool != "task.output" {
		t.Errorf("兄弟路径的命中来源判成 %q, want task.output：命中的不是宿主正文里的兄弟路径那一维", h.SrcTool)
	}
	b.CloseTask(task)
}

// tail183 clips a stub for a failure message: head and tail only, so the
// reading stays readable without pasting 20 KB into a test log.
func tail183(s string) string {
	r := []rune(s)
	if len(r) <= 240 {
		return s
	}
	return string(r[:120]) + " …[剪]… " + string(r[len(r)-120:])
}
