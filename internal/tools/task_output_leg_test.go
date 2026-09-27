package tools

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

// Ticket 164 AC#3 - a long task output must carry a pointer that RECOVERS the
// full text, and AC#2's after-half (the identical call that landed on
// "未知工具" now reaches a real tool).
//
// Injection surface: the real bridge and the real fs.read - the two seams
// AGENTS §1.3 allows - never a stubbed provider standing in for a real one.
//
// The reverse criteria this file must be able to turn red (the dispatch names
// both; a cell that cannot answer "which mutation lights it up" is decoration):
//
//	只截不指      -> TestLongOutputPointerRecoversEveryByte
//	形状不符 D15 -> TestTruncationShapeIsTheD15Triple
//
// and the third one is a LIMITATION pinned on purpose:
//
//	产物超 256 KiB 时后半段读不到 -> TestPointerPast256KiBIsNotFullyReadable

// taskBridge wires fs.* plus task.* over one roster, through the real bridge.
func taskBridge(t *testing.T, roster *TaskRoster, allowed ...string) *Bridge {
	t.Helper()
	paths := NewPathCanonicalizer(allowed, nil)
	reg := NewRegistry()
	for _, e := range append(BuiltinFSEntries(FSDeps{Paths: paths}),
		BuiltinTaskEntries(TaskDeps{Roster: roster})...) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		Gate:       NoGate{}, Logf: func(string, ...any) {},
	})
}

func TestTaskOutputAC2AfterLegIsReachable(t *testing.T) {
	dir := tempCanonical(t)
	roster := NewTaskRoster()
	roster.Record("bg-1", TaskOutput{Text: "后台任务打印的第一行"})
	b := taskBridge(t, roster, dir)

	// The byte-for-byte identical request that task_output_ac2_before_test.go
	// sends against today's roster.
	r := agent.ToolRequest{
		TaskID: "164-ac2-after",
		CallID: "call-164-ac2-after",
		Name:   "task.output",
		Args:   json.RawMessage(`{"task_id":"bg-1"}`),
	}
	out, err := b.Execute(t.Context(), r)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("AC#2 AFTER verbatim: IsError=%v ErrorClass=%q Text=%q", out.IsError, out.ErrorClass, out.Text)

	if out.IsError || strings.Contains(out.Text, "未知工具") {
		t.Fatalf("the same call that read 未知工具 before must now reach the tool: %+v", out)
	}
	if out.Text != "后台任务打印的第一行" {
		t.Fatalf("Text = %q, want the recorded output verbatim", out.Text)
	}
	if out.Truncated {
		t.Errorf("a short output must not be truncated")
	}
}

// TestUnknownTaskIDIsLoud pins ticket 164 ruling 2: v1 keeps no roster across
// restarts, so an id that is not on the table must come back as a refusal that
// names the id - never as an empty success, which is the shape that reads to
// the model as "the task printed nothing".
func TestUnknownTaskIDIsLoudNotEmpty(t *testing.T) {
	dir := tempCanonical(t)
	b := taskBridge(t, NewTaskRoster(), dir)

	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "164-r2", CallID: "c1", Name: "task.output",
		Args: json.RawMessage(`{"task_id":"ghost-9"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("an unknown task id must be an error outcome, got: %+v", out)
	}
	if out.Text == "" || !strings.Contains(out.Text, "ghost-9") {
		t.Fatalf("the refusal must name the id it could not find, got %q", out.Text)
	}
	if !strings.Contains(out.Text, "没有这条记录") {
		t.Errorf("the refusal must distinguish 无记录 from 无输出, got %q", out.Text)
	}
	if out.ErrorClass != "tool" {
		t.Errorf("error_class = %q, want tool (self-correctable, D37)", out.ErrorClass)
	}
}

// TestMissingArgsAndUnwiredRoster both fail closed.
func TestMissingArgsAndUnwiredRoster(t *testing.T) {
	dir := tempCanonical(t)
	b := taskBridge(t, NewTaskRoster(), dir)
	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "a", CallID: "c", Name: "task.output", Args: json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError || !strings.Contains(out.Text, "task_id") {
		t.Fatalf("missing task_id must be refused, got %+v", out)
	}

	bare := New(Options{
		Registry: mustRegisterTaskEntries(t, TaskDeps{Roster: nil}),
		Gate:     NoGate{}, Logf: func(string, ...any) {},
	})
	out2, err := bare.Execute(t.Context(), agent.ToolRequest{
		TaskID: "b", CallID: "c", Name: "task.output", Args: json.RawMessage(`{"task_id":"x"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !out2.IsError || !strings.Contains(out2.Text, "未接线") {
		t.Fatalf("an unwired roster must fail closed, got %+v", out2)
	}
}

func mustRegisterTaskEntries(t *testing.T, d TaskDeps) *Registry {
	t.Helper()
	reg := NewRegistry()
	for _, e := range BuiltinTaskEntries(d) {
		if err := reg.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	return reg
}

// pointerRe pulls the 全文见 <path> out of the stub. The wording is the D15(3)
// one internal/agent/spill.go:141 already puts in front of the model, so this
// test reads the same sentence the model would.
var pointerRe = regexp.MustCompile(`全文见 (\S+)…`)

// TestLongOutputPointerRecoversEveryByte is AC#3's cell, and it is the reverse
// criterion for 只截不指: drop the pointer from taskOutput.Execute and this test
// goes red, because everything it asserts is downstream of that path.
//
// It asserts READ-BACK, not "a path was mentioned": the file the stub names is
// re-read twice - once raw, and once through a real fs.read call - and the
// bytes must equal the original output the roster holds. The announced total
// length must equal the real artifact length, not the stub's length.
func TestLongOutputPointerRecoversEveryByte(t *testing.T) {
	dir := tempCanonical(t)
	full := asciiRun(20000) // 5000 tokens under the 4-bytes-per-token heuristic

	// The artifact is landed by the host, exactly as D15(3) does it, and the
	// roster is then filed with the bytes plus where they went.
	artName := "tool-output-164-bg7.txt"
	artPath := filepath.Join(dir, artName)
	if err := os.WriteFile(artPath, []byte(full), 0o600); err != nil {
		t.Fatal(err)
	}
	roster := NewTaskRoster()
	roster.Record("bg-7", TaskOutput{Text: full, ArtifactPath: mustCanonical(t, artPath)})

	b := taskBridge(t, roster, dir)
	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "164-ac3", CallID: "c1", Name: "task.output",
		Args: json.RawMessage(`{"task_id":"bg-7"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.IsError {
		t.Fatalf("a recorded task must answer, not refuse: %+v", out)
	}
	if !out.Truncated {
		t.Fatal("a 5000-token output must be marked truncated")
	}

	m := pointerRe.FindStringSubmatch(out.Text)
	if m == nil {
		t.Fatalf("the stub must carry a recoverable pointer, got: %q", out.Text)
	}
	pointer := m[1]

	// (1) the pointer names a file whose bytes ARE the whole output.
	raw, err := os.ReadFile(pointer)
	if err != nil {
		t.Fatalf("the pointer names no readable file: %v", err)
	}
	if string(raw) != full {
		t.Fatalf("the artifact is not the full text (%d bytes on disk, %d in the roster)",
			len(raw), len(full))
	}

	// (2) and the pointer is walkable by the mechanism PLAN.md:431 names: a real
	// fs.read of that path returns the same bytes.
	reread, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "164-ac3", CallID: "c2", Name: "fs.read",
		Args: json.RawMessage(`{"path":"` + jsonEscape(filepath.ToSlash(pointer)) + `"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if reread.IsError {
		t.Fatalf("the pointer must be re-readable with fs.read, got: %+v", reread)
	}
	if reread.Text != full {
		t.Fatalf("fs.read returned %d of %d bytes and did not say why", len(reread.Text), len(full))
	}

	// (3) the announced total is the REAL length, not the stub's.
	if !strings.Contains(out.Text, "总长 "+itoa(len(full))+" 字节") {
		t.Fatalf("stub must announce 总长 %d 字节, got: %q", len(full), out.Text)
	}
	if strings.Contains(out.Text, "总长 "+itoa(len(out.Text))) {
		t.Errorf("the announced length is the stub's own, not the full text's")
	}
	// and head+pointer+tail must not silently eat bytes the announcement denies:
	if n := omittedAnnounced(out.Text); n >= 0 {
		if got := len(full) - n; got != len(headOf(out.Text))+len(tailOf(out.Text)) {
			t.Errorf("省略 %d 与实际保留 %d 字节对不上（全文 %d）", n, got, len(full))
		}
	}
}

// TestTruncationShapeIsTheD15Triple is the second reverse criterion: the shape
// itself (head 500 / tail 200 / total present) must be pinned, so a wrong
// constant in task.go turns this red.
func TestTruncationShapeIsTheD15Triple(t *testing.T) {
	dir := tempCanonical(t)
	full := asciiRun(20000)
	roster := NewTaskRoster()
	roster.Record("bg-8", TaskOutput{Text: full, ArtifactPath: dir})

	b := taskBridge(t, roster, dir)
	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "shape", CallID: "c1", Name: "task.output",
		Args: json.RawMessage(`{"task_id":"bg-8"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	head, tail := headOf(out.Text), tailOf(out.Text)
	// The expected sizes are LITERALS, not task.go's own constants: a test that
	// reads the production constant cannot catch that constant being wrong, and
	// "截断形状与 D15 不符" must be catchable. 500 / 200 tokens at the shared
	// 4 bytes-per-token heuristic is what PLAN.md:431 freezes.
	if len(head) != 500*4 {
		t.Errorf("head = %d bytes, want %d (D15 head 500 tokens x 4)", len(head), 500*4)
	}
	if len(tail) != 200*4 {
		t.Errorf("tail = %d bytes, want %d (D15 tail 200 tokens x 4)", len(tail), 200*4)
	}
	if !strings.Contains(out.Text, "约 "+itoa(agent.ApproxTokens(full))+" token") {
		t.Errorf("the stub must announce the total in tokens, got %q", out.Text)
	}
	if head != full[:2000] || tail != full[len(full)-800:] {
		t.Error("head/tail are not the leading and trailing windows of the full text")
	}
}

// TestTruncationBudgetsAreRead proves the D15 numbers are wired, not echoed: a
// host that scales its budget gets a different cut.
func TestTruncationBudgetsAreRead(t *testing.T) {
	dir := tempCanonical(t)
	full := asciiRun(20000)
	roster := NewTaskRoster()
	roster.Record("bg-9", TaskOutput{Text: full, ArtifactPath: dir})

	paths := NewPathCanonicalizer([]string{dir}, nil)
	reg := NewRegistry()
	if err := reg.Register(Entry{Tool: taskOutput{d: TaskDeps{
		Roster: roster, SpillTokens: 100, SpillHeadTokens: 10, SpillTailTokens: 5,
	}}, Decl: TaskOutputDecl()}); err != nil {
		t.Fatal(err)
	}
	b := New(Options{
		Registry: reg, Paths: paths, Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		Gate: NoGate{}, Logf: func(string, ...any) {},
	})

	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "scaled", CallID: "c1", Name: "task.output",
		Args: json.RawMessage(`{"task_id":"bg-9"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(headOf(out.Text)); got != 40 {
		t.Errorf("head = %d bytes with a 10-token budget, want 40", got)
	}
	if got := len(tailOf(out.Text)); got != 20 {
		t.Errorf("tail = %d bytes with a 5-token budget, want 20", got)
	}
}

// TestNoCopyFileAnnouncesLoss covers the branch that must stay loud: a long
// output the host never landed. Truncating without a pointer is 不合格, so the
// only acceptable answer is to say the rest cannot be recovered.
func TestNoCopyFileAnnouncesLoss(t *testing.T) {
	dir := tempCanonical(t)
	roster := NewTaskRoster()
	roster.Record("bg-10", TaskOutput{Text: asciiRun(20000)}) // ArtifactPath: ""
	b := taskBridge(t, roster, dir)

	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "noloss", CallID: "c1", Name: "task.output",
		Args: json.RawMessage(`{"task_id":"bg-10"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Truncated {
		t.Fatal("it is truncated, and must say so")
	}
	if strings.Contains(out.Text, "全文见") {
		t.Fatalf("it must not point at a file that does not exist: %q", out.Text)
	}
	if !strings.Contains(out.Text, "不可找回") {
		t.Fatalf("a pointerless truncation must announce the loss, got %q", out.Text)
	}
}

// TestPointerPast256KiBIsNotFullyReadable is the 后半段 cell, and the
// disposition is 断言它响 (see the evidence file for why not 明写): fs.read has
// no offset (internal/tools/fs.go:113-116) and a 256 KiB ceiling (:36
// defaultMaxReadBytes), so a pointer to anything past that returns the FIRST
// 256 KiB and nothing else. This test pins that as a canary: the day Q-59 is
// approved and fs.read grows an offset, this test goes red and forces this
// ticket's cell to be re-read. Nothing in it is widened to look better.
func TestPointerPast256KiBIsNotFullyReadable(t *testing.T) {
	dir := tempCanonical(t)
	full := asciiRun(300 * 1024) // > the 256 KiB single-read ceiling
	name := "tool-output-164-big.txt"
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(full), 0o600); err != nil {
		t.Fatal(err)
	}
	roster := NewTaskRoster()
	roster.Record("bg-11", TaskOutput{Text: full, ArtifactPath: mustCanonical(t, p)})

	b := taskBridge(t, roster, dir)
	out, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "big", CallID: "c1", Name: "task.output",
		Args: json.RawMessage(`{"task_id":"bg-11"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	pm := pointerRe.FindStringSubmatch(out.Text)
	if pm == nil {
		// Fail, never panic: a swallowed panic here would take the whole
		// package's remaining readings with it.
		t.Fatalf("a >256 KiB artifact must still be pointed at, got %q", out.Text)
	}
	pointer := pm[1]

	head, err := b.Execute(t.Context(), agent.ToolRequest{
		TaskID: "big", CallID: "c2", Name: "fs.read",
		Args: json.RawMessage(`{"path":"` + jsonEscape(filepath.ToSlash(pointer)) + `"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if head.IsError {
		t.Fatalf("the first segment must read back, got %+v", head)
	}
	if len(head.Text) != defaultMaxReadBytes {
		t.Fatalf("read back %d bytes, want exactly the %d-byte ceiling", len(head.Text), defaultMaxReadBytes)
	}
	if !head.Truncated {
		t.Error("fs.read must mark the capped read truncated")
	}
	if head.Text != full[:defaultMaxReadBytes] {
		t.Error("the returned segment is not the leading bytes of the artifact")
	}
	// THE LIMITATION, asserted: asking for more through the parameters fs.read
	// actually has cannot reach past the ceiling, so the tail of this artifact
	// is unreachable today.
	if r, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "big", CallID: "c3", Name: "fs.read",
		Args: json.RawMessage(`{"path":"` + jsonEscape(filepath.ToSlash(pointer)) + `","max_bytes":9999999}`),
	}); err != nil {
		t.Fatal(err)
	} else if len(r.Text) >= len(full) {
		t.Fatalf("the whole artifact came back - fs.read grew an offset, update this cell")
	}
}

// TestRosterIsNotAPathSurface is the C26/AGENTS §1.2 boundary, stated as a
// behaviour instead of a comment: a task_id shaped like a path is a roster key,
// and must be answered as "no such task" - not by reaching for a file.
func TestRosterIsNotAPathSurface(t *testing.T) {
	dir := tempCanonical(t)
	roster := NewTaskRoster()
	roster.Record("bg-12", TaskOutput{Text: "真输出"})
	b := taskBridge(t, roster, dir)

	for _, evil := range []string{"../../etc/passwd", `..\..\windows\win.ini`, dir} {
		out, err := b.Execute(t.Context(), agent.ToolRequest{
			TaskID: "evil", CallID: "c1", Name: "task.output",
			Args: mustArgs(t, map[string]any{"task_id": evil}),
		})
		if err != nil {
			t.Fatal(err)
		}
		if !out.IsError || !strings.Contains(out.Text, "查不到这个任务") {
			t.Fatalf("path-shaped id %q must be a miss, not a read: %+v", evil, out)
		}
		if strings.Contains(out.Text, filepath.Base(dir)) && evil != dir {
			t.Fatalf("the refusal echoed a canonicalized path: %q", out.Text)
		}
	}
}

// ---------------------------------------------------------------------------
// small helpers (no new production surface)
// ---------------------------------------------------------------------------

func asciiRun(n int) string {
	var b strings.Builder
	b.Grow(n)
	for i := 0; b.Len() < n; i++ {
		b.WriteRune(rune('a' + i%26))
		if i%40 == 39 {
			b.WriteByte('\n')
		}
	}
	return b.String()[:n]
}

func mustArgs(t *testing.T, m map[string]any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}

// headOf / tailOf pull the two kept windows out of a stub, i.e. everything
// before / after the bracketed announcement line.
func headOf(stub string) string {
	i := strings.Index(stub, "\n[…")
	if i < 0 {
		return stub
	}
	return stub[:i]
}

func tailOf(stub string) string {
	i := strings.Index(stub, "…]\n")
	if i < 0 {
		return ""
	}
	return stub[i+len("…]\n"):]
}

// omittedRe reads the 省略 N 字符 field the stub announces.
var omittedRe = regexp.MustCompile(`省略 (\d+) 字符`)

func omittedAnnounced(stub string) int {
	m := omittedRe.FindStringSubmatch(stub)
	if m == nil {
		return -1
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return -1
	}
	return n
}

func itoa(n int) string { return strconv.Itoa(n) }
