package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// ticket 162 AC#1 - the READING, taken on code that ticket 162 has not touched
// ---------------------------------------------------------------------------
//
// The premise this ticket is built on is a sentence nobody had measured:
// "if the model re-copies a whole file and drops a block, the loss is silent."
// This file measures it against today's fs.write and nothing else - fs.write's
// own behaviour is NOT under test here and is not modified by it.
//
// Three things get a number:
//
//	(a) the bytes really did disappear;
//	(b) the call answered as a success (no error, no warning, no counter);
//	(c) what any AFTER-THE-FACT instrument can say about it - the tool_call row,
//	    the audit line, and a re-read through fs.read.
//
// The assertions below pin the ABSENT signals, not just the present ones: each
// one is written so that it fails loudly the day someone adds the missing
// instrument, which is the difference between a reading and a shrug.

// ac1Block is the region the simulated model drops on the floor.
var ac1Block = []string{
	"beta-01 保留段第一行：这一段共六行，模型抄写时把它整段漏掉了",
	"beta-02 保留段第二行",
	"beta-03 保留段第三行",
	"beta-04 保留段第四行",
	"beta-05 保留段第五行",
	"beta-06 保留段第六行",
}

// ac1SpelledFile builds head + block + tail and returns the bytes plus the
// same file with the block dropped (what a model that "re-typed the whole
// thing" and lost a section would hand back).
func ac1SpelledFile() (original, mutated string) {
	var o, m strings.Builder
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(&o, "alpha-%02d 头部保留行\n", i)
		fmt.Fprintf(&m, "alpha-%02d 头部保留行\n", i)
	}
	for _, l := range ac1Block {
		fmt.Fprintf(&o, "%s\n", l)
	}
	for i := 1; i <= 8; i++ {
		fmt.Fprintf(&o, "gamma-%02d 尾部保留行\n", i)
		fmt.Fprintf(&m, "gamma-%02d 尾部保留行\n", i)
	}
	return o.String(), m.String()
}

// ac1Bridge wires the fs family over a real C26 resolver, a real SQLite
// journal and a capturing audit sink - the same composition helpers_test.go
// describes, with the two sinks this reading needs added.
func ac1Bridge(t *testing.T, root string, j agent.Journal, audit *[]string) *Bridge {
	t.Helper()
	paths := NewPathCanonicalizer([]string{root}, nil)
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths, Journal: j,
		Gate:       &gateSpy{approveAns: AnswerAllow}, // an overwrite is the L2 row
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       func(format string, args ...any) { *audit = append(*audit, fmt.Sprintf(format, args...)) },
	})
}

// TestFSWriteSilentLossIsNotReported is AC#1: one fs.write whose content is
// the file with a block left out.
func TestFSWriteSilentLossIsNotReported(t *testing.T) {
	original, mutated := ac1SpelledFile()
	if strings.Contains(mutated, ac1Block[0]) {
		t.Fatal("the mutation shape is wrong: the dropped block is still there")
	}

	root := sealableTempDir124(t)
	canonicalRoot := mustCanonical(t, root)
	store := openStore(t, filepath.Join(root, "data"))
	mustStartTask(t, store, "task-1")
	var audit []string
	b := ac1Bridge(t, canonicalRoot, store, &audit)

	target := filepath.Join(root, "notes.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	// (a) BEFORE: what was there.
	beforeBytes, beforeLines := len(original), countLines(original)

	out, err := b.Execute(t.Context(), req("fs.write",
		args(t, map[string]any{"path": slash(target), "content": mutated})))
	if err != nil {
		t.Fatalf("a tool-level failure must travel as an outcome, not a Go error: %v", err)
	}

	after := readString(t, target)
	afterBytes, afterLines := len(after), countLines(after)
	lostLines := beforeLines - afterLines
	lostBytes := beforeBytes - afterBytes

	// -----------------------------------------------------------------------
	// (a) the loss is real
	// -----------------------------------------------------------------------
	t.Logf("AC#1(a) 落盘前 %d 字节 / %d 行；落盘后 %d 字节 / %d 行；少了 %d 行 %d 字节",
		beforeBytes, beforeLines, afterBytes, afterLines, lostLines, lostBytes)
	if lostLines != len(ac1Block) {
		t.Fatalf("the mutation dropped %d lines, want %d - the reading is measuring the wrong thing",
			lostLines, len(ac1Block))
	}
	for _, l := range ac1Block {
		if strings.Contains(after, l) {
			t.Fatalf("(a) violated: the block is still in the file, nothing was lost")
		}
	}
	if after != mutated {
		t.Fatal("(a) violated: the file does not hold what was sent")
	}

	// -----------------------------------------------------------------------
	// (b) the call reports nothing about it
	// -----------------------------------------------------------------------
	t.Logf("AC#1(b) IsError=%v Truncated=%v RiskLevel=%s ErrorClass=%q",
		out.IsError, out.Truncated, out.RiskLevel, out.ErrorClass)
	t.Logf("AC#1(b) Result.Text=%q", out.Text)
	for i, s := range out.AppliedSteps {
		t.Logf("AC#1(b) AppliedSteps[%d]=%q", i, s)
	}
	if out.IsError || out.ErrorClass != "" {
		t.Fatalf("(b) would be falsified: the call came back as a failure: %+v", out)
	}
	// The only number in the answer is the size of what was WRITTEN. Nothing in
	// it is measured against what was there before, so no reader can compute a
	// delta out of it - which is exactly why "silent" is the right word.
	if !strings.Contains(out.Text, fmt.Sprintf("%d 字节", afterBytes)) {
		t.Fatalf("(b): the answer no longer carries its byte count, re-read this test: %q", out.Text)
	}
	for _, leak := range []string{"少", "丢", "缺", "删除", "行差", "警告", "warn", "missing", "dropped"} {
		if strings.Contains(out.Text, leak) || containsStep(out.AppliedSteps, leak) {
			t.Fatalf("(b) would be falsified: the answer mentions %q: %q / %v", leak, out.Text, out.AppliedSteps)
		}
	}

	// -----------------------------------------------------------------------
	// (c) after the fact, can any instrument tell "N lines are gone"?
	// -----------------------------------------------------------------------
	rows, err := store.ListToolCallsByTask(context.Background(), "task-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("want exactly the one booked fs.write row, got %d", len(rows))
	}
	r := rows[0]
	t.Logf("AC#1(c) tool_call 行：tool=%s risk=%s decision=%s outcome=%s error_class=%q",
		r.Tool, r.RiskLevel, r.Decision, r.Outcome, r.ErrorClass)
	t.Logf("AC#1(c) tool_call 行的全部列名：%v", structFieldsOf(memory.ToolCall{}))
	if r.Outcome != agent.OutcomeSuccess || r.ErrorClass != "" {
		t.Fatalf("(c) would be falsified: the row itself reports a problem: %+v", r)
	}
	// The row is the forensic record. It cannot carry the loss because it has no
	// field that could hold it - not a size, not a line count, not a diff.
	for _, f := range structFieldsOf(memory.ToolCall{}) {
		for _, want := range []string{"size", "bytes", "line", "diff", "sha", "hash"} {
			if strings.Contains(strings.ToLower(f), want) {
				t.Errorf("(c) would be falsified: tool_call now has a %q-bearing column %q", want, f)
			}
		}
	}
	// The audit line is the operator's record. Same conclusion, checked the same
	// way: the pre-write size is the number a delta needs, and it appears
	// nowhere in it.
	for _, l := range audit {
		t.Logf("AC#1(c) audit: %q", l)
	}
	if auditMentions(audit, beforeBytes) {
		t.Errorf("(c) would be falsified: the audit line carries the pre-write size %d, so a reader could compute the loss", beforeBytes)
	}

	// And a re-read through the same tool family: fs.read answers with the bytes
	// and nothing else - no count, no "this file is 6 lines shorter than before".
	rd, err := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, target)))
	if err != nil || rd.IsError {
		t.Fatalf("(c) the re-read itself failed: %+v err=%v", rd, err)
	}
	t.Logf("AC#1(c) fs.read 事后复读：%d 字节，Truncated=%v，文本里没有任何行数/大小字段（Text 就是文件内容本身）",
		len(rd.Text), rd.Truncated)

	// -----------------------------------------------------------------------
	// the verdict, in one place
	// -----------------------------------------------------------------------
	t.Logf("AC#1 读数结论：丢了 %d 行 %d 字节；Result.Text / AppliedSteps / tool_call 行 / "+
		"audit 行 / fs.read 事后复读 五处信号全部为零 ⇒ 今天的 fs.write 对\"整文件写回时漏抄\"没有任何检测",
		lostLines, lostBytes)
}

// countLines counts newline-terminated lines (the file's own shape, not a
// strings.Fields guess).
func countLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n")
}

// structFieldsOf lists a row struct's field names - the cheapest way to ask
// "could this record even hold the number we are looking for?".
func structFieldsOf(v any) []string {
	rt := reflect.TypeOf(v)
	out := make([]string, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		out = append(out, rt.Field(i).Name)
	}
	return out
}

// containsStep greps the ledger for one word.
func containsStep(steps []string, word string) bool {
	for _, s := range steps {
		if strings.Contains(s, word) {
			return true
		}
	}
	return false
}

// auditMentions reports whether any audit line carries one number as a token.
func auditMentions(lines []string, n int) bool {
	token := json.Number(fmt.Sprint(n)).String()
	for _, l := range lines {
		if strings.Contains(l, token) {
			return true
		}
	}
	return false
}
