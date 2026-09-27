package tools

// 162-v1 ACCEPTANCE rigs (independent non-implementer cells). This file is
// NEVER in the tree: it reaches the compiler through `go test -overlay` only
// (see probes/162/v1/*.json). Everything below is re-derived from the ticket
// face; none of the implementer's readings are taken as input.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
)

// v1Fields is this cell's OWN reflection helper (not r1's structFieldsOf).
func v1Fields(v any) []string {
	rt := reflect.TypeOf(v)
	out := make([]string, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		out = append(out, rt.Field(i).Name)
	}
	return out
}

// v1SizeBearing is the token set the ticket's "the ledger structurally cannot
// hold the minuend" claim rests on.
func v1SizeBearing(names []string) []string {
	hits := []string{}
	for _, f := range names {
		lf := strings.ToLower(f)
		for _, tok := range []string{"size", "bytes", "line", "diff", "sha", "hash"} {
			if strings.Contains(lf, tok) {
				hits = append(hits, f+"/"+tok)
			}
		}
	}
	return hits
}

// v1NegativeControl proves the matcher above is not a dead switch: a struct
// that DOES carry such a field must be reported.
func v1NegativeControl(t *testing.T) {
	t.Helper()
	type bearing struct {
		SizeBytes int
		LineCount int
		DiffText  string
		SHA       string
	}
	got := v1SizeBearing(v1Fields(bearing{}))
	if len(got) < 4 {
		t.Fatalf("matcher is dead: want >=4 hits on the control struct, got %d (%v)", len(got), got)
	}
	type clean struct {
		Tool   string
		Reason string
	}
	if n := len(v1SizeBearing(v1Fields(clean{}))); n != 0 {
		t.Fatalf("matcher over-reports: clean struct yielded %d hits", n)
	}
	t.Logf("V1-AC#1 反射尺自证：控制结构体命中 %d 枚 %v（非死开关），无名结构体命中 0 枚", len(got), got)
	if n := len(v1SizeBearing(v1Fields(memory.ToolCall{}))); n != 0 {
		t.Logf("NOTE memory.ToolCall now carries %d size-bearing field(s): %v", n,
			v1SizeBearing(v1Fields(memory.ToolCall{})))
	}
}

// ---------------------------------------------------------------------------
// AC#1, my OWN dropped-block shot (different file, different block, different
// position than the implementer's 672/24 -> 450/18).
// ---------------------------------------------------------------------------

func v1Spelled() (string, string) {
	var o, m strings.Builder
	for i := 1; i <= 5; i++ {
		fmt.Fprintf(&o, "h-%02d 验收自有头部行\n", i)
		fmt.Fprintf(&m, "h-%02d 验收自有头部行\n", i)
	}
	for i := 1; i <= 4; i++ {
		fmt.Fprintf(&o, "drop-%02d 这一行在写回时被漏抄\n", i)
	}
	for i := 1; i <= 4; i++ {
		fmt.Fprintf(&o, "t-%02d 尾部保留行\n", i)
		fmt.Fprintf(&m, "t-%02d 尾部保留行\n", i)
	}
	return o.String(), m.String()
}

func TestV1AC1WholeFileRewriteDropsMiddleBlockUnreported(t *testing.T) {
	original, mutated := v1Spelled()
	beforeBytes, beforeLines := len(original), strings.Count(original, "\n")
	if strings.Contains(mutated, "drop-01") {
		t.Fatal("rig is wrong: nothing was dropped")
	}
	if beforeLines != 13 {
		t.Fatalf("rig is wrong: want 13 lines, got %d", beforeLines)
	}
	v1NegativeControl(t)

	root := sealableTempDir124(t)
	canonicalRoot := mustCanonical(t, root)
	store := openStore(t, filepath.Join(root, "data"))
	mustStartTask(t, store, "task-1") // req() books under this fixed task id (bridge_test.go:147)
	var audit []string
	b := ac1Bridge(t, canonicalRoot, store, &audit)

	target := filepath.Join(root, "v1notes.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	out, err := b.Execute(t.Context(), req("fs.write",
		args(t, map[string]any{"path": slash(target), "content": mutated})))
	if err != nil {
		t.Fatalf("a tool failure must travel as an outcome: %v", err)
	}
	after := readString(t, target)
	afterLines := strings.Count(after, "\n")

	// (a) the loss is real and it is MY loss shape (middle block, 4 lines).
	t.Logf("V1-AC#1(a) 落盘前 %d 字节 / %d 行 → 落盘后 %d 字节 / %d 行（少了 %d 行 %d 字节，中段四行）",
		beforeBytes, beforeLines, len(after), afterLines, beforeLines-afterLines, beforeBytes-len(after))
	if beforeLines-afterLines != 4 || after != mutated {
		t.Fatalf("(a) violated: %d lines lost, after==mutated? %v", beforeLines-afterLines, after == mutated)
	}
	// (b) the call answers success, and the answer's only number is what went OUT.
	t.Logf("V1-AC#1(b) IsError=%v Truncated=%v RiskLevel=%s ErrorClass=%q Text=%q",
		out.IsError, out.Truncated, out.RiskLevel, out.ErrorClass, out.Text)
	if out.IsError || out.ErrorClass != "" {
		t.Fatalf("(b) would be falsified: the call failed: %+v", out)
	}
	if !strings.Contains(out.Text, fmt.Sprintf("%d 字节", len(mutated))) {
		t.Fatalf("(b): the answer no longer carries the sent size, re-read: %q", out.Text)
	}
	if strings.Contains(out.Text, fmt.Sprintf("%d", beforeBytes)) {
		t.Fatalf("(b) would be falsified: the pre-write size %d appears in the answer", beforeBytes)
	}
	for _, leak := range []string{"少", "丢", "缺", "删除", "行差", "警告", "warn", "missing", "dropped", "truncat"} {
		if strings.Contains(out.Text, leak) || containsStep(out.AppliedSteps, leak) {
			t.Fatalf("(b) would be falsified: the answer mentions %q: %q / %v", leak, out.Text, out.AppliedSteps)
		}
	}
	// (c) after-the-fact instruments: my own reflection, my own audit scan.
	rows, err := store.ListToolCallsByTask(context.Background(), "task-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("want one booked row, got %d", len(rows))
	}
	r := rows[0]
	t.Logf("V1-AC#1(c) row tool=%s risk=%s decision=%s outcome=%s error_class=%q",
		r.Tool, r.RiskLevel, r.Decision, r.Outcome, r.ErrorClass)
	if r.Outcome != agent.OutcomeSuccess || r.ErrorClass != "" {
		t.Fatalf("(c) would be falsified: the row reports a problem: %+v", r)
	}
	t.Logf("V1-AC#1(c) memory.ToolCall 字段名=%v", v1Fields(memory.ToolCall{}))
	if hits := v1SizeBearing(v1Fields(memory.ToolCall{})); len(hits) != 0 {
		t.Errorf("(c) would be falsified: tool_call carries size-bearing column(s) %v", hits)
	}
	t.Logf("V1-AC#1(c) agent.ToolOutcome 字段名=%v size-bearing=%v",
		v1Fields(agent.ToolOutcome{}), v1SizeBearing(v1Fields(agent.ToolOutcome{})))
	t.Logf("V1-AC#1(c) Result 字段名=%v size-bearing=%v",
		v1Fields(Result{}), v1SizeBearing(v1Fields(Result{})))
	for _, l := range audit {
		t.Logf("V1-AC#1(c) audit: %q", l)
	}
	if auditMentions(audit, beforeBytes) {
		t.Errorf("(c) would be falsified: audit carries the pre-write size %d", beforeBytes)
	}
	rd, rerr := b.Execute(t.Context(), req("fs.read", pathArgsOf(t, target)))
	if rerr != nil || rd.IsError {
		t.Fatalf("re-read failed: %+v %v", rd, rerr)
	}
	t.Logf("V1-AC#1(c) 事后 fs.read：%d 字节 Truncated=%v ErrorClass=%q", len(rd.Text), rd.Truncated, rd.ErrorClass)
	t.Logf("V1-AC#1 结论：漏抄中段 4 行 %d 字节，Result.Text / AppliedSteps / tool_call 行 / audit 行 / "+
		"事后 fs.read 五处读数全为安静 ⇒ 票面 :32 的前提由验收位自己复算成立", beforeBytes-len(after))
}

// TestV1AC4TargetIsNeverHalfWritten is the acceptance cell for AC#4's load bit,
// written as an INVARIANT instead of as a ledger probe: after a mid-write stop
// the target must hold EITHER the pre-edit bytes OR the complete new bytes -
// never a prefix of either. Run twice: unmutated (staged writer) it must pass;
// under MUT-4 (temp+rename removed from fs.edit's own write path) it must go
// red on the BYTE READING, which is the assertion their kill case never reaches
// because its boundary probe fires first.
func TestV1AC4TargetIsNeverHalfWritten(t *testing.T) {
	const original = "alpha\nbravo\ncharlie\ndelta\n"                           // 26 bytes
	const intended = "ALPHA-LONGER\nbravo\nCHARLIE-LONGER-THAN-BEFORE\ndelta\n" // 53 bytes
	var boundaries []string
	root := sealableTempDir124(t)
	target := filepath.Join(root, "inv.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ := fsEditBridgeDeps(t, mustCanonical(t, root), func(d *FSDeps) {
		d.WriteChunk = 8
		d.Hooks = Hooks{
			AtStep: func(step string) { boundaries = append(boundaries, step) },
			Kill: func(step string) error {
				if step == "write:16" {
					return errors.New("验收位模拟：进程在此刻被杀")
				}
				return nil
			},
		}
	})
	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha", "new": "ALPHA-LONGER"},
		map[string]any{"old": "charlie", "new": "CHARLIE-LONGER-THAN-BEFORE"})))
	if err != nil {
		t.Fatal(err)
	}
	got := readString(t, target)
	t.Logf("V1-AC#4 不变式读数：杀点边界序列=%q IsError=%v 目标 %d 字节（原文 %d／全量 %d）内容=%q",
		boundaries, out.IsError, len(got), len(original), len(intended), got)
	if !out.IsError {
		t.Fatalf("the killed call must report a failure: %+v", out)
	}
	if got != original && got != intended {
		t.Fatalf("AC#4 invariant broken: the target holds a PREFIX (%d of %d bytes) - a half-written file, "+
			"neither the old content nor the new one: %q", len(got), len(intended), got)
	}
	t.Logf("V1-AC#4 不变式成立：目标要么是原文要么是完整新内容（本次=%d 字节）", len(got))
}

func v1Endings(s string) (int, int, int) {
	crlf := strings.Count(s, "\r\n")
	return crlf, strings.Count(s, "\r") - crlf, strings.Count(s, "\n") - crlf
}

func TestV1AC3ThreeShotsOnMyOwnFixtures(t *testing.T) {
	t.Run("crlf_file_stays_crlf", func(t *testing.T) {
		root := sealableTempDir124(t)
		const original = "one\r\ntwo\r\nthree\r\nfour\r\n"
		target := filepath.Join(root, "a.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		if c, cr, lf := v1Endings(original); c != 4 || cr != 0 || lf != 0 {
			t.Fatalf("fixture not pure CRLF: %d/%d/%d", c, cr, lf)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "two\r\nthree", "new": "dos\r\ntres\r\ndos2"})))
		if err != nil || out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		after := readString(t, target)
		c, cr, lf := v1Endings(after)
		if cr != 0 || lf != 0 {
			t.Fatalf("AC#3 violated: CRLF file gained bare endings cr=%d lf=%d in %q", cr, lf, after)
		}
		if c != 5 || after != "one\r\ndos\r\ntres\r\ndos2\r\nfour\r\n" {
			t.Fatalf("wrong bytes landed: %d pairs %q", c, after)
		}
		t.Logf("V1-AC#3 正向①：%d 字节纯 CRLF → %d 字节，CRLF 对 4→%d，裸 CR/裸 LF 0/0",
			len(original), len(after), c)
	})

	t.Run("bom_survives", func(t *testing.T) {
		root := sealableTempDir124(t)
		original := "\xef\xbb\xbf" + "name: v1\nstate: open\n"
		target := filepath.Join(root, "b.md")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "state: open", "new": "state: closed"})))
		if err != nil || out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		after := readString(t, target)
		if string([]byte(after)[:3]) != "\xef\xbb\xbf" {
			t.Fatalf("AC#3 violated: BOM gone, file starts % x", []byte(after)[:3])
		}
		if after != "\xef\xbb\xbf"+"name: v1\nstate: closed\n" {
			t.Fatalf("changed beyond the match: %q", after)
		}
		t.Logf("V1-AC#3 正向②：%d 字节带 BOM → %d 字节，前三字节仍 EF BB BF，出现 %d 次",
			len(original), len(after), strings.Count(after, "\xef\xbb\xbf"))
	})

	t.Run("lf_file_endings_untouched", func(t *testing.T) {
		root := sealableTempDir124(t)
		const original = "one\ntwo\nthree\n"
		target := filepath.Join(root, "c.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "two", "new": "2"})))
		if err != nil || out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		after := readString(t, target)
		if strings.ContainsRune(after, '\r') {
			t.Fatalf("AC#3 violated: CR appeared in an LF file: %q", after)
		}
		if after != "one\n2\nthree\n" {
			t.Fatalf("wrong bytes: %q", after)
		}
		t.Logf("V1-AC#3 反向：%d 字节 → %d 字节，裸 LF 3 条、CR 0 枚", len(original), len(after))
	})
}

// TestV1AC3bPremiseTheTwoShapesClaimedSilent are MY readings of the two shapes
// the orchestrator wrote into ticket face :36/:37 as "今天不响". The ticket face
// says these numbers came from the implementer and were NOT recomputed by the
// orchestrator, so this cell recomputes them independently: if either one is
// actually loud today, the ticket face's premise is wrong and must be fixed.
func TestV1AC3bPremiseTheTwoShapesClaimedSilent(t *testing.T) {
	t.Run("shape1_old_spanning_the_bom", func(t *testing.T) {
		root := sealableTempDir124(t)
		original := "\xef\xbb\xbf" + "alpha\n" // 9 bytes: 3 BOM + "alpha\n"
		target := filepath.Join(root, "s1.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		if len(original) != 9 {
			t.Fatalf("fixture arithmetic: %d", len(original))
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "\xef\xbb\xbfalpha", "new": "ALPHA"})))
		if err != nil {
			t.Fatal(err)
		}
		after := readString(t, target)
		t.Logf("V1-AC#3b① 读数：IsError=%v ErrorClass=%q %d 字节 → %d 字节，前三字节 % x，Text=%q",
			out.IsError, out.ErrorClass, len(original), len(after), []byte(after)[:3], out.Text)
		// The claim under test: this lands, silently. If it refuses, the ticket
		// face's AC#3b① premise is FALSE and the orchestrator must fix :36.
		if out.IsError {
			t.Errorf("PREMISE CHECK: shape1 is LOUD today (IsError=true) - ticket face :37 ① is wrong")
		}
		if strings.HasPrefix(after, "\xef\xbb\xbf") {
			t.Errorf("PREMISE CHECK: the BOM survived - ticket face :37 ① (BOM 被一起换掉) is wrong")
		}
		if len(after) != 6 {
			t.Errorf("expected 9 -> 6 bytes (BOM's three bytes consumed), got %d (%q)", len(after), after)
		}
		for _, probe := range []string{"BOM", "字节序", "编码"} {
			if strings.Contains(out.Text, probe) {
				t.Errorf("PREMISE CHECK: the answer now signals the encoding change via %q", probe)
			}
		}
	})

	t.Run("shape2_lf_spelled_new_into_a_pure_crlf_file", func(t *testing.T) {
		root := sealableTempDir124(t)
		const original = "one\r\ntwo\r\nthree\r\n"
		target := filepath.Join(root, "s2.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "two", "new": "two\ntwo-b"})))
		if err != nil {
			t.Fatal(err)
		}
		after := readString(t, target)
		c, cr, lf := v1Endings(after)
		t.Logf("V1-AC#3b② 读数：IsError=%v ErrorClass=%q 落盘后 CRLF=%d 裸CR=%d 裸LF=%d，Text=%q",
			out.IsError, out.ErrorClass, c, cr, lf, out.Text)
		if out.IsError {
			t.Errorf("PREMISE CHECK: shape2 is LOUD today - ticket face :37 ② is wrong")
		}
		if lf != 1 {
			t.Errorf("expected exactly one bare LF (mixed endings), got %d", lf)
		}
		for _, probe := range []string{"行尾", "不一致", "mixed endings", "endings"} {
			if strings.Contains(out.Text, probe) {
				t.Errorf("PREMISE CHECK: the answer signals the endings change via %q", probe)
			}
		}
	})
}
