package tools

// ticket 162 AC#3b (dispatch r3 cell ①) plus the half-cell gap 162-v1 registered
// against it: the two shapes that the ticket names at 票面 :37 measured SILENT on
// the code at this cell's anchor, and "CRLF 文件改完还是 CRLF" had only a forward
// ruler. Readings of the silent state are in
// docs/evidence/s1/162-fs-edit-r3.md §0a (baseline log
// .scratch/wisp/probes/162/r3/baseline-gotest-verbose.txt, lines 126 and 128):
//
//	shape ①  14 字节带 BOM 的文件、old 从偏移 0 起 ⇒ BOM 被当普通字节换掉，落盘 11 字节，
//	        回执 IsError=false、文案零提示
//	shape ②  纯 CRLF 文件里一枚 LF 写法的 new ⇒ 落盘 3 组 CRLF + 1 条裸 LF，
//	        回执 IsError=false、文案里没有任何行尾提示
//
// Both are now REFUSALS, and both happen before the first byte is written, so
// they inherit this tool's all-or-nothing property (the ledger stays empty, the
// target keeps its bytes). Why refuse rather than "land it with a loud warning"
// (票面 :38 lets either be chosen):
//
//   - a warning has to land the damage first. The tool's stated principle is
//     响亮失败、绝不安静改坏, and 改坏 with a footnote is still 改坏 - the BOM is
//     gone and the file is mixed before anyone reads the warning.
//   - both shapes have a one-line fix the model can act on (抄 old 从 BOM 之后；
//     new 按文件的行尾拼写), which is what 判据形状 2/3 require a refusal to carry.
//     A case with no remedial action would be a wall, and is written below as a
//     control instead.
//   - neither refusal widens or narrows WHAT matches (that is 判据形状 1 的第二半,
//     the normalization layer this file still does not implement). They only
//     refuse a caller-supplied new/old that would damage a property of the file
//     nobody named - so the合法对照 cases below all still land.
//
// 恒真那一问，逐枚（两向读数在证据件 §3）：本件两枚拒绝用例、r2 那两枚翻转用例，
// 加上四枚合法对照（BOM 之后抄起的 old、按 CRLF 拼写的 new、无 BOM 文件第 0 字节起
// 的 old、已混合文件的 new），一起在 **未修码**（base overlay＝本 cell 锚点那一版
// fs_edit.go，命令见 §3）上跑：对照全绿、断言拒绝的那几枚全红 - 那一版根本没有这两
// 道检查，"今天静默"这一前提就是这一格能立项的理由。反向尺那一枚
// （TestFSEdit…MixedFileRewritten）在未修码上是绿的 - 它不测修码，它测 M-3：那一味
// 之下它必须红，证据件 §4。

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// fsEditBridgeAudited is fsEditBridge with exactly one addition: the bridge's own
// audit sink is captured, because AC#3b's "外部可见" claim has to be pointed at a
// reading somebody else can see, not at a test-local string. Bridge.Logf is the
// same sink production wires to the app logger; the line it writes here is the
// one that carries outcome= for the call.
func fsEditBridgeAudited(t *testing.T, root string, audit *[]string) *Bridge {
	t.Helper()
	paths := NewPathCanonicalizer([]string{root}, nil)
	d := FSDeps{Paths: paths}
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(d) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths,
		Gate:       &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow},
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       func(format string, args ...any) { *audit = append(*audit, fmt.Sprintf(format, args...)) },
	})
}

// auditLineFor returns the i-th audit line booking an fs.edit call, or fails: an
// assertion about "what the operator can see" that quietly matched nothing would
// be the weakest kind of green this repo has produced. wantSubs are checked
// against that line; pass index 0 for the only call so far.
func auditLineFor(t *testing.T, audit []string, idx int, wantSubs ...string) string {
	t.Helper()
	var hits []string
	for _, l := range audit {
		if !strings.Contains(l, "tool=fs.edit") || !strings.Contains(l, "outcome=") {
			continue
		}
		ok := true
		for _, w := range wantSubs {
			if !strings.Contains(l, w) {
				ok = false
			}
		}
		if ok {
			hits = append(hits, l)
		}
	}
	if len(hits) <= idx {
		t.Fatalf("wanted audit line #%d carrying %v, got %d matches in %q", idx, wantSubs, len(hits), audit)
	}
	return hits[idx]
}

// ---------------------------------------------------------------------------
// AC#3b shape ①: an old that swallows the file's BOM
// ---------------------------------------------------------------------------

// TestFSEditRefusesAnOldThatSwallowsTheFileBOM is AC#3b ①: the matched region
// overlaps EF BB BF at the front of the file, so the tool used to replace the
// encoding marker as an incidental part of the edit (14 bytes in, 11 out, no
// signal). The case pins four things: the call refuses; the refusal names the
// three bytes and where they are; NOTHING lands, not even partially; and the
// operator-visible reading flips with it (the receipt's IsError/ErrorClass and
// the bridge's audit line outcome= both change value).
func TestFSEditRefusesAnOldThatSwallowsTheFileBOM(t *testing.T) {
	root := sealableTempDir124(t)
	original := utf8BOM + "alpha\nbeta\n" // 14 bytes, BOM owns [0,3)
	target := filepath.Join(root, "bom-zero.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if len(original) != 14 || !strings.HasPrefix(readString(t, target), utf8BOM) {
		t.Fatalf("the fixture is not a 14-byte BOM file: %d bytes, % x", len(original), []byte(original))
	}
	var audit []string
	b := fsEditBridgeAudited(t, mustCanonical(t, root), &audit)

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": utf8BOM + "alpha", "new": "ALPHA"})))
	if err != nil {
		t.Fatalf("a refusal is an outcome, not a host fault: %v", err)
	}
	assertRefusal(t, out, "UTF-8 BOM", "EF BB BF", "BOM 之后的第一个字节")
	assertUntouched(t, target, original)

	// The external reading, both halves of it, in the same run.
	line := auditLineFor(t, audit, 0, "outcome=error")
	t.Logf("AC#3b① 改后读数（回执）：IsError=%v ErrorClass=%q AppliedSteps=%d",
		out.IsError, out.ErrorClass, len(out.AppliedSteps))
	t.Logf("AC#3b① 改后读数（回执文案）：%q", out.Text)
	t.Logf("AC#3b① 改后读数（审计行，外部可见）：%s", line)
	if strings.Contains(out.Text, "已改写") {
		t.Errorf("a refusal must never also claim the edit landed: %q", out.Text)
	}

	// 合法对照（同文件、同一次调用之后的第二发）：old 从 BOM 之后抄起，就必须命中并
	// 落盘，BOM 还在。这一枚是"门没被焊死"的凭据 - 没有它，上面的拒绝可以靠
	// "凡是带 BOM 的文件一律不回"做到，而那是另一回事。
	out2, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha", "new": "ALPHA"})))
	if err != nil || out2.IsError {
		t.Fatalf("the control edit (old copied AFTER the BOM) must land: out=%+v err=%v", out2, err)
	}
	const want = utf8BOM + "ALPHA\nbeta\n"
	if got := readString(t, target); got != want {
		t.Fatalf("the control edit changed the wrong bytes:\n got=%q\nwant=%q", got, want)
	}
	t.Logf("AC#3b① 合法对照：%d 字节 → %d 字节，前三字节仍是 EF BB BF，文案=%q",
		len(original), len(want), out2.Text)
}

// TestFSEditStillLandsAnOldStartingAtByteZeroWithoutABOM is ①'s other boundary:
// the check is "the match overlaps a BOM that is THERE", not "old starts at
// offset 0". A file with no BOM whose first bytes are exactly what the caller
// named must keep landing, or the refusal would be a wall instead of a ruler.
func TestFSEditStillLandsAnOldStartingAtByteZeroWithoutABOM(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\nbeta\n" // no BOM: offset 0 is ordinary content
	target := filepath.Join(root, "no-bom.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(original, utf8BOM) {
		t.Fatal("the control fixture must not carry a BOM")
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "alpha", "new": "ALPHA"})))
	if err != nil || out.IsError {
		t.Fatalf("an offset-0 old in a BOM-less file must land: out=%+v err=%v", out, err)
	}
	const want = "ALPHA\nbeta\n"
	if got := readString(t, target); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	t.Logf("AC#3b① 合法对照（无 BOM 文件、old 从第 0 字节起）：%d 字节 → %d 字节，未被拒绝",
		len(original), len(want))
}

// TestFSEditCannotBeFedAPartialBOMThroughJSON records how far ① actually reaches.
// The code checks "match overlaps the first three bytes", which reads as if an
// old carrying ONE or TWO of those bytes were in scope. It is not, and not because
// of this tool: a JSON string cannot carry a lone continuation byte, so the
// transport replaces it (Go's encoder answers U+FFFD) and the old arrives as
// something the file does not contain. This case measures that instead of
// asserting it: the partial-BOM old is refused as a 0-hit miss, NOT as a BOM
// swallow, so the BOM check's only reachable input stays "old covers EF BB BF".
func TestFSEditCannotBeFedAPartialBOMThroughJSON(t *testing.T) {
	root := sealableTempDir124(t)
	original := utf8BOM + "alpha\n" // 7 bytes; byte index 2 is the BOM's last byte
	target := filepath.Join(root, "bom-partial.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	partial := original[2:5] // "\xbf" + "al": not a valid UTF-8 sequence
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": partial, "new": "X"})))
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, out, "0 命中")
	if strings.Contains(out.Text, "BOM") {
		t.Errorf("the BOM check must not be what fired here: %q", out.Text)
	}
	assertUntouched(t, target, original)
	t.Logf("AC#3b① 射程读数：old 只带 BOM 的第 3 个字节时，落盘的是 %d 字节、第 0 字节是 %#U，"+
		"传进来的已经不是那三个字节：%q", len(original), []rune(original)[0], out.Text)
}

// ---------------------------------------------------------------------------
// AC#3b shape ②: a new that breaks the file's own line-ending convention
// ---------------------------------------------------------------------------

// TestFSEditRefusesALFSpelledNewInAPureCRLFFile is AC#3b ②: the old matches
// byte-exactly (so nothing refuses for that reason) and the new brings bare LFs
// into a pure-CRLF file. That used to land 3 CRLF pairs + 1 bare LF, unsignalled.
// The refusal must say which edit did it and what the file would become, and the
// three controls below must all still land: the same batch spelled CRLF, an
// LF-spelled new in an LF file, and any new at all in a file that has no single
// convention to break.
func TestFSEditRefusesALFSpelledNewInAPureCRLFFile(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\r\nbravo\r\ncharlie\r\n" // 23 bytes, pure CRLF
	target := filepath.Join(root, "crlf-mixed.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if crlf, cr, lf := endings(original); crlf != 3 || cr != 0 || lf != 0 {
		t.Fatalf("fixture is not pure CRLF: %d/%d/%d", crlf, cr, lf)
	}
	var audit []string
	b := fsEditBridgeAudited(t, mustCanonical(t, root), &audit)

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "bravo", "new": "bravo\nbravo2"})))
	if err != nil {
		t.Fatal(err)
	}
	assertRefusal(t, out, "行尾约定", "第 1 枚编辑的 new", "裸 LF 1 条", "CRLF 3 组")
	assertUntouched(t, target, original)
	line := auditLineFor(t, audit, 0, "outcome=error")
	t.Logf("AC#3b② 改后读数（回执）：IsError=%v ErrorClass=%q AppliedSteps=%d",
		out.IsError, out.ErrorClass, len(out.AppliedSteps))
	t.Logf("AC#3b② 改后读数（回执文案）：%q", out.Text)
	t.Logf("AC#3b② 改后读数（审计行，外部可见）：%s", line)

	// control 1 - the ticket's own named 对照: the same edit, new spelled the way
	// the file spells it, lands and the file stays pure.
	out2, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "bravo", "new": "bravo\r\nbravo2"})))
	if err != nil || out2.IsError {
		t.Fatalf("a CRLF-spelled new in a CRLF file must land: out=%+v err=%v", out2, err)
	}
	const want = "alpha\r\nbravo\r\nbravo2\r\ncharlie\r\n"
	after := readString(t, target)
	if after != want {
		t.Fatalf("control 1 changed the wrong bytes:\n got=%q\nwant=%q", after, want)
	}
	if crlf, cr, lf := endings(after); crlf != 4 || cr != 0 || lf != 0 {
		t.Fatalf("control 1 lost the file's purity: crlf=%d bareCR=%d bareLF=%d", crlf, cr, lf)
	}
	auditLineFor(t, audit, 0, "outcome=success")
	t.Logf("AC#3b② 合法对照一（同一枚编辑、new 按 CRLF 拼写）：%d 字节 → %d 字节，CRLF 3→4 组、裸 CR/裸 LF 均 0",
		len(original), len(after))
}

// TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak is the other two
// controls of ②, and the honest boundary of the check: an LF file edited in LF
// keeps landing (the door is not one-way), and a file that ALREADY mixes is not
// refused - it has no single convention, so "which shape is the right one" has no
// answer this tool may invent. Its endings are instead pinned byte-exactly by
// TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten.
func TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak(t *testing.T) {
	t.Run("lf_file_lands_an_lf_spelled_new", func(t *testing.T) {
		root := sealableTempDir124(t)
		const original = "alpha\nbravo\ncharlie\n"
		target := filepath.Join(root, "lf.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "bravo", "new": "bravo\nbravo2"})))
		if err != nil || out.IsError {
			t.Fatalf("the same new that ② refuses in a CRLF file must land in an LF file: out=%+v err=%v", out, err)
		}
		const want = "alpha\nbravo\nbravo2\ncharlie\n"
		if got := readString(t, target); got != want {
			t.Fatalf("got %q want %q", got, want)
		}
		t.Logf("AC#3b② 合法对照二（LF 文件 + LF 写法 new）：%d 字节 → %d 字节，未被拒绝",
			len(original), len(want))
	})

	t.Run("mixed_file_is_not_refused", func(t *testing.T) {
		root := sealableTempDir124(t)
		const original = "alpha\r\nbravo\ncharlie\r\n" // 22 bytes: 2 CRLF + 1 bare LF
		target := filepath.Join(root, "already-mixed.txt")
		if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
			t.Fatal(err)
		}
		b, _ := fsEditBridge(t, mustCanonical(t, root))
		out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
			map[string]any{"old": "bravo", "new": "BRAVO\nbravo2"})))
		if err != nil || out.IsError {
			t.Fatalf("a file with no single convention must not be refused: out=%+v err=%v", out, err)
		}
		const want = "alpha\r\nBRAVO\nbravo2\ncharlie\r\n"
		if got := readString(t, target); got != want {
			t.Fatalf("got %q want %q", got, want)
		}
		t.Logf("AC#3b② 合法对照三（已混合的文件不被告知该长成什么样）：%d 字节 → %d 字节，行尾形 %s",
			len(original), len(want), describeEndings(want))
	})
}

// ---------------------------------------------------------------------------
// the gap half-cell: 162-v1 measured that "改完还是 CRLF" has NO reverse ruler
// ---------------------------------------------------------------------------

// TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten is the missing
// reverse ruler. 162-v1 ran M-3 (force the bytes about to land into CRLF,
// mut/m3: `updated := strings.ReplaceAll(strings.ReplaceAll(b.String(), "\r\n",
// "\n"), "\n", "\r\n")`) and measured that BOTH
// TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF and
// TestFSEditLineEndingForensicsRefusesRealSpellings stayed GREEN - because
// forcing CRLF onto a pure-CRLF file is the identity transform, so no fixture
// made of only-pure-CRLF bytes can ever tell "the tool preserves CRLF" apart from
// "the tool stamps CRLF". The discriminator has to be a file that CONTAINS CRLF
// and is not made of nothing else: the bare LF in the middle is what a forced
// rewrite touches, and the two CRLF pairs are what it must leave alone.
//
// So this case asserts the whole blob byte-for-byte AND each ending count, and it
// is the case that must go RED under M-3 (evidence: probes/162/r3/logs/m3-on-r3.log,
// named in 证据件 §4). It is deliberately NOT the AC#3b refusal path: the edit
// brings no break of its own, so the tool is being measured as a byte mover here,
// not as a judge.
//
// 恒真那一问：未修码（base overlay）上它绿 - 那一版也不改行尾，这一枚不是"前提尺"
// 而是"反向尺"，它的红来自 M-3 那一味，不来自修码与否（两向读数见证据件 §4）。
func TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten(t *testing.T) {
	root := sealableTempDir124(t)
	const original = "alpha\r\nbravo\ncharlie\r\n" // 22 bytes: 2 CRLF pairs + 1 bare LF
	target := filepath.Join(root, "mixed-reverse-ruler.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	if crlf, cr, lf := endings(original); crlf != 2 || cr != 0 || lf != 1 {
		t.Fatalf("the ruler needs a CRLF-bearing, LF-touched fixture: %d/%d/%d", crlf, cr, lf)
	}
	b, _ := fsEditBridge(t, mustCanonical(t, root))

	out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "bravo", "new": "BRAVO"})))
	if err != nil || out.IsError {
		t.Fatalf("an edit that brings no line break of its own must land: out=%+v err=%v", out, err)
	}
	const want = "alpha\r\nBRAVO\ncharlie\r\n"
	after := readString(t, target)
	if after != want {
		t.Fatalf("the落盘 bytes are not the original bytes with ONE region swapped (a forced-ending rewrite looks exactly like this): got %q want %q",
			after, want)
	}
	crlf, cr, lf := endings(after)
	if crlf != 2 || cr != 0 || lf != 1 {
		t.Fatalf("反向尺 violated: 改完还是 CRLF 这一侧没有尺了 - crlf=%d bareCR=%d bareLF=%d in %q",
			crlf, cr, lf, after)
	}
	if strings.Contains(after, "\r\r") || strings.Contains(after, "\n\n") {
		t.Fatalf("the edit invented a doubled break: %q", after)
	}
	noStagingFilesLeft(t, root)
	t.Logf("反向尺读数：%d 字节混合文件（%s）改完 %d 字节，CRLF 仍 2 组、裸 LF 仍 1 条、裸 CR 0 枚，文案=%q",
		len(original), describeEndings(original), len(after), out.Text)
}
