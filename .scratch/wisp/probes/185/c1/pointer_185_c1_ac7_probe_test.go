package risk

import (
	"strings"
	"testing"
)

// 185-c1 probe leg — ticket 185 AC#7（同名不同目录那一形）.
//
// This file is NOT tracked by the repo: it exists only as a `-overlay` copy at
// internal/risk/pointer_185_c1_ac7_probe_test.go (see
// .scratch/wisp/probes/185/c1/overlay-ac7.json for the unmutated leg and
// overlay-ac7-mutant.json for the basename-widening mutation).
//
// Rulers for the shipped face (line numbers drift, symbol names do not):
//
//	grep -n "^func Test" internal/risk/pointer_183_test.go
//	grep -n "spellsDeclaredPath\|attachDeclaredPath" internal/risk/taintmatch.go | grep -v _test.go
//
// Shape: the host declared ONE artifact path; the marked body ALSO spells,
// verbatim, a sibling path that carries the SAME file name under a different
// directory. Ticket 177 W-2's boundary is that a foreign spelling of a path
// must stay evidence, so a reread of the sibling must still hit R4. This is the
// shape ticket 183 AC#4 handed over as AC#7: today the landed exemption compares
// whole normalized strings (spellsDeclaredPath is rune-exact over the declared
// path), so it is expected to PASS — and the mutation is what shows whether any
// judgement in the repo would notice a basename-level widening.

const (
	// p185Declared is the path the host minted and declared on this mark.
	p185Declared = `C:/Users/swq/AppData/Roaming/wisp/artifacts/001/tool-output-agent-task-c66c0634.txt`
	// p185Sibling is a DIFFERENT file the host never declared: same basename,
	// another directory.
	p185Sibling = `C:/Users/swq/AppData/Roaming/wisp/artifacts/002/tool-output-agent-task-c66c0634.txt`
	// p185Secret is ordinary tainted body that spells nothing of either path.
	p185Secret = "SUPERSECRET-185C1-BODY"
	// p185Body is the host-written stub: the declared path once, plus the
	// sibling path spelled verbatim at a different position.
	p185Body = "报告头部 " + p185Secret + " 中段 […输出已落文件：省略 17200 字符，" +
		"总长 20000 字节，全文见 " + p185Declared +
		"…]\n另一份兄弟产物 " + p185Sibling + " 在这次运行正文里逐字出现 结束"
)

// p185Tail is local on purpose: the leg must compile against the SHIPPED
// taintmatch.go (unmutated leg) as well as against the mutant copy.
func p185Tail(s string) string {
	if i := strings.LastIndexAny(s, `/\`); i >= 0 {
		return s[i+1:]
	}
	return s
}

func TestPointer185C1SameNameDifferentDirectoryStillHits(t *testing.T) {
	dn, sn := normalizeTaint(p185Declared), normalizeTaint(p185Sibling)
	if dn == "" || sn == "" || dn == sn {
		t.Fatalf("fixture broken: declared=%q sibling=%q", dn, sn)
	}
	if p185Tail(dn) != p185Tail(sn) {
		t.Fatalf("fixture broken: AC#7 is about the SAME basename, got %q vs %q", p185Tail(dn), p185Tail(sn))
	}
	if len([]rune(sn)) < 8 {
		t.Fatal("fixture broken: sibling must be longer than the fragment floor")
	}

	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", p185Body, p185Declared) {
		t.Fatal("host stub must be marked")
	}
	h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": p185Sibling})
	if !ok {
		t.Fatal("AC#7: a same-name sibling under another directory, spelled verbatim in the body, must still hit R4 (basename-level exemption = widening)")
	}
	if h.SrcTool != SrcFSRead {
		t.Fatalf("AC#7: the hit must come from the host mark's own non-exempt windows, got %q", h.SrcTool)
	}
}

// Control leg: the host's OWN declared path stays re-readable (this is what
// ticket 183 landed). It keeps the AC#7 leg from passing for a boring reason —
// if nothing were exempted at all, the forward leg would still be green.
func TestPointer185C1DeclaredPathItselfStaysClean(t *testing.T) {
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", p185Body, p185Declared) {
		t.Fatal("host stub must be marked")
	}
	if _, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": p185Declared}); ok {
		t.Fatal("control: the host's own declared pointer must not be evidence against a reread")
	}
}
