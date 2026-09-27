package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ticket 162 AC#2: the four refusals of 判据形状 1-4, one case each, plus the
// control that proves the tool can also do its job.
//
// "这一发在未修码上响不响" is answered per case below. The short version, for
// all four: NO. Unmodified code has no fs.edit to call, so none of these four
// errors can be raised; the same intent expressed through the only writer that
// does exist (fs.write, whole file) is accepted and answered with a success -
// which is what AC#1 measured (see fswrite_silentloss_ac1_test.go: 6 lines,
// 222 bytes gone, five instruments, zero signal). These cases therefore buy
// nothing unless the tool also works, hence TestFSEditAppliesABatchOnDisk.
//
// Every refusal asserts the file's bytes are IDENTICAL afterwards, because
// "响亮失败且绝不落盘" is one claim, not two.

// fsEditBridge wires the family over one root with an approval gate that lets
// the L2 fs.edit calls through to Execute, so what is under test is the tool's
// own refusal and not the gate's.
func fsEditBridge(t *testing.T, root string) (*Bridge, *gateSpy) {
	t.Helper()
	g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
	paths := NewPathCanonicalizer([]string{root}, nil)
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths, Gate: g,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       func(string, ...any) {},
	}), g
}

// editArgs renders one fs.edit call.
func editArgs(t *testing.T, path string, edits ...map[string]any) string {
	t.Helper()
	return args(t, map[string]any{"path": slash(path), "edits": edits})
}

// assertUntouched is the 绝不落盘 half: same bytes, and nothing staged left
// behind in the directory.
func assertUntouched(t *testing.T, path, want string) {
	t.Helper()
	if got := readString(t, path); got != want {
		t.Fatalf("a refused call must not move one byte:\n before=%q\n  after=%q", want, got)
	}
	noStagingFilesLeft(t, filepath.Dir(path))
}

// assertRefusal is the shared shape of a loud failure: IsError, a remedial
// action in the text, and an empty ledger (nothing was ever applied).
func assertRefusal(t *testing.T, out agent.ToolOutcome, wantSubs ...string) {
	t.Helper()
	if !out.IsError {
		t.Fatalf("expected a hard error, got a success: %+v", out)
	}
	if out.ErrorClass != "tool" {
		t.Errorf("a tool-level refusal must carry D37 class %q, got %q (the model has to be able to self-correct)",
			"tool", out.ErrorClass)
	}
	for _, w := range wantSubs {
		if !strings.Contains(out.Text, w) {
			t.Errorf("the refusal must state the remedial action; missing %q in:\n%s", w, out.Text)
		}
	}
	if len(out.AppliedSteps) != 0 {
		t.Errorf("a refusal must have an empty ledger, got %q", out.AppliedSteps)
	}
}

// ---------------------------------------------------------------------------
// 判据形状 1 + 2: locating is a literal byte string, and missing it is loud
// ---------------------------------------------------------------------------

// TestFSEditRefusesANonLiteralOld is AC#2 case 1+2: the call carries something
// that is NOT a byte-exact substring - here, the same line with the indentation
// re-flowed, the shape a model produces when it "remembers" a paragraph. The
// refusal is hard, it names the remedial action in the upstream wording, and
// the file keeps its bytes.
//
// 未修码上响不响：不响。fs.edit 不存在，于是"缩进记错了"这一发根本走不到定位
// 那一步；今天唯一的通路是 fs.write 把整文件重抄一遍，缩进错了也算成功。
func TestFSEditRefusesANonLiteralOld(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "package main\n\nfunc main() {\n\tprintln(\"keep me\")\n}\n"
	target := filepath.Join(root, "main.go")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	// The old below uses spaces where the file uses a tab: a diff-format or
	// line-number based locator would still find line 4. A literal one must not.
	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "    println(\"keep me\")", "new": "    println(\"changed\")"})))
	if err != nil {
		t.Fatalf("a refusal is an outcome, not a host fault: %v", err)
	}
	assertRefusal(t, out, "0 命中", "must match exactly including all whitespace and newlines",
		"fs.read")
	assertUntouched(t, target, original)
	t.Logf("判据 1+2 读数：命中前 %d 字节 / 命中失败后 %d 字节（一字未变），文案=%q",
		len(original), len(readString(t, target)), out.Text)
}

// ---------------------------------------------------------------------------
// 判据形状 3: several hits are a DIFFERENT hard error, counted before refusing
// ---------------------------------------------------------------------------

// TestFSEditRefusesAnAmbiguousOldIsADistinctError is AC#2 case 3: the old
// appears twice. The refusal must be the ambiguity error (not the exact-match
// one), it must carry the count it got from counting first, and it must NOT
// quietly take the first occurrence.
//
// 未修码上响不响：不响。今天的 fs.write 没有"命中几处"这个概念 - 它不问文件
// 里有什么，所以也无从歧义；歧义在今天的通路上表现为"改了第一处、第二处原样
// 留着"，而且成功。
func TestFSEditRefusesAnAmbiguousOldIsADistinctError(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "item: apple\nnote: x\nitem: apple\nnote: y\n"
	target := filepath.Join(root, "list.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "item: apple", "new": "item: apricot"})))
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, out, "命中 2 处", "provide more surrounding context to make it unique")
	// The two error shapes are told apart on purpose: "找不到" and "找到多处"
	// need different fixes from the model.
	if strings.Contains(out.Text, "0 命中") {
		t.Errorf("an ambiguous old must not report as a miss: %q", out.Text)
	}
	// And it is NOT "first one wins".
	assertUntouched(t, target, original)
	if n := strings.Count(readString(t, target), "item: apricot"); n != 0 {
		t.Fatalf("the refusal took an occurrence anyway: %d", n)
	}
	t.Logf("判据 3 读数：old 在 %d 字节文件里命中 2 处，先数后拒，落盘前后一字未变", len(original))
}

// ---------------------------------------------------------------------------
// 判据形状 4 (first half): an empty old is refused
// ---------------------------------------------------------------------------

// TestFSEditRefusesAnEmptyOldIsAWholeFileInsert is AC#2 case 4a: old:"" is not
// "match nothing", it is "insert this into every position", i.e. the one call
// shape whose blast radius is the whole file.
//
// 未修码上响不响：不响，而且方向相反 - 今天要做"往文件里插一段"只能把整个文件
// 加上那一段重抄一遍交回来（fs.write 接受），插没插对没有任何读数。
func TestFSEditRefusesAnEmptyOldIsAWholeFileInsert(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "一行既有内容\n第二行\n"
	target := filepath.Join(root, "notes.md")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "", "new": "凭空加的一段\n"})))
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, out, "old 是空字符串", "非空")
	assertUntouched(t, target, original)

	// The same batch with a missing old field is the other malformed shape and
	// must not be read as an empty one.
	out2, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"new": "x"})))
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, out2, "没有给 old 字段")
	assertUntouched(t, target, original)
	t.Logf("判据 4 空 old 读数：两形（old:\"\" 与缺 old 字段）都拒，文件保持 %d 字节", len(original))
}

// ---------------------------------------------------------------------------
// 判据形状 4 (second half): all or nothing across a batch
// ---------------------------------------------------------------------------

// TestFSEditIsAllOrNothingAcrossABatch is AC#2 case 4b. Three batches that all
// contain one good edit and one bad one: a no-op (old == new), an overlapping
// pair, and a miss. In none of them may the good edit reach the disk - that
// half-written state is exactly the "silently lost something" class AC#1
// measured, arriving from the other direction.
//
// 未修码上响不响：不响。fs.write 只会看到"整文件少了一段/多了一段"，然后成功。
func TestFSEditIsAllOrNothingAcrossABatch(t *testing.T) {
	const original = "alpha\nbravo\ncharlie\ndelta\n"
	cases := []struct {
		name     string
		edits    []map[string]any
		wantSubs []string
	}{
		{
			name: "no_op_edit_in_the_batch",
			edits: []map[string]any{
				{"old": "alpha", "new": "ALPHA"},
				{"old": "charlie", "new": "charlie"},
			},
			wantSubs: []string{"old 与 new 相同"},
		},
		{
			name: "overlapping_edits",
			edits: []map[string]any{
				{"old": "alpha\nbravo", "new": "A-B"},
				{"old": "bravo\ncharlie", "new": "B-C"},
			},
			wantSubs: []string{"重叠的区间"},
		},
		{
			name: "second_edit_misses",
			edits: []map[string]any{
				{"old": "alpha", "new": "ALPHA"},
				{"old": "zulu", "new": "ZULU"},
			},
			wantSubs: []string{"0 命中"},
		},
		{
			name:     "empty_batch",
			edits:    []map[string]any{},
			wantSubs: []string{"edits 是空的"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := sealableTempDir124(t)
			target := filepath.Join(root, "batch.txt")
			if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
				t.Fatal(err)
			}
			b, _ := fsEditBridge(t, mustCanonical(t, root))

			var editsJSON []any
			for _, e := range tc.edits {
				editsJSON = append(editsJSON, e)
			}
			out, err := b.Execute(t.Context(), req("fs.edit",
				args(t, map[string]any{"path": slash(target), "edits": editsJSON})))
			if err != nil {
				t.Fatal(err)
			}
			assertRefusal(t, out, tc.wantSubs...)
			// The good edit did NOT land either: all or nothing.
			assertUntouched(t, target, original)
			t.Logf("全有或全无读数：本批 %d 枚编辑，其中一枚不成立（空批＝整批不成立）⇒ 一个字也不落盘，"+
				"文件仍是 %d 字节", len(tc.edits), len(readString(t, target)))
		})
	}
}

// ---------------------------------------------------------------------------
// the control: the tool also works
// ---------------------------------------------------------------------------

// TestFSEditAppliesABatchOnDisk is the reading that keeps AC#2 from being
// satisfiable by a tool that always says no: a two-edit batch lands, in the
// model's listing order irrelevant, through D31's staged writer, and the answer
// names the before/after size - the one number fs.write's answer never has.
func TestFSEditAppliesABatchOnDisk(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\nbravo\ncharlie\ndelta\n"
	target := filepath.Join(root, "ok.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	b, g := fsEditBridge(t, mustCanonical(t, root))

	// Listed second-then-first on purpose: the applied bytes must not depend on
	// that order.
	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "charlie", "new": "CHARLIE\nextra"},
		map[string]any{"old": "alpha", "new": "ALPHA"})))
	if err != nil || out.IsError {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	const want = "ALPHA\nbravo\nCHARLIE\nextra\ndelta\n"
	if got := readString(t, target); got != want {
		t.Fatalf("the batch did not land as expected:\n got=%q\nwant=%q", got, want)
	}
	noStagingFilesLeft(t, root)
	t.Logf("对照读数：%d 字节 / %d 行 → %s", len(original), strings.Count(original, "\n"), out.Text)
	t.Logf("对照读数：AppliedSteps=%q", out.AppliedSteps)

	// Declared L2 means the approval route, and the refusal cases above all ran
	// after that route said yes - so the tool's own errors are being tested,
	// not the gate's.
	if w, a := g.counts(); w != 0 || a != 1 {
		t.Errorf("fs.edit must reach the L2 approval channel exactly once: window=%d approval=%d", w, a)
	}
	dec := g.approvalDecision()
	if dec.Level != risk.L2 {
		t.Errorf("fs.edit was judged at %v, want the D34 declared L2 floor", dec.Level)
	}
	for _, s := range out.AppliedSteps {
		if strings.Contains(s, "原子重命名") {
			return
		}
	}
	t.Errorf("the write must go through the D31 staged writer, ledger=%q", out.AppliedSteps)
}
