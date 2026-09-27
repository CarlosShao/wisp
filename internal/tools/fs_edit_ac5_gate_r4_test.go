package tools

// Ticket 162 AC#5 — 风险档与门控：fs.edit 的判定必须与 fs.write 同族，且越界路径
// 必须被 C26 那一道判定拒掉。
//
// 票面 09:4x 的 `>` 更正①把这里的档位从"与 fs.write 同族（L1）"改判成 **L2**
// （理由＝R8 覆盖已存在 → L2，与 fs.write 的**覆盖**那一行同级；票面原句一字未抹），
// 并留下"越界路径必须被拒、不许新增豁免、不动 allowlist.txt"三句原样有效。
// 所以本文件钉的是**两件事**，不是档位字符串：
//
//	① 同一枚越界样本分别喂 fs.write 与 fs.edit，两边的档位 / 命中的规则 /
//	   走到的那一条通道 / 审计行，逐字对上——"同族"是这枚并排读数，不是一句注释；
//	② 越界那一发必须**真的过不去**：判定在桥上（C19 的 R2 读 C26 的 InAllowlist），
//	   拒在写入之前，目标字节一个也没动。
//
// ⚠ 本程没有动 `internal/risk/**`（owner 给 162 的射程里没有那块地）。档位由
// `FSEditDecl` 的声明（`fs_edit.go`，`Declared: risk.L2`）+ 共用的 R2 达成，
// 摘掉哪一行会红写在证据件里，判据本体在下面的断言里。

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// r4DeniedBridge is fsEditBridgeAudited with the ONE knob this cell needs
// changed: the human side says **no**. A fixture that answers Allow cannot say
// anything about "门控真让调用过不了", so this is the same builder shape
// (real C26 over a real [fs] allowed_dirs, real C19 assessor, real audit sink)
// with the gate answering Reject at the approval channel and Veto at the L1
// window - whichever channel a call lands on, it stops there.
func r4DeniedBridge(t *testing.T, root string, audit *[]string) (*Bridge, *gateSpy) {
	t.Helper()
	g := &gateSpy{approveAns: AnswerReject, windowAns: AnswerVeto}
	paths := NewPathCanonicalizer([]string{root}, nil)
	if u := paths.UnusableRoots(); len(u) > 0 {
		t.Fatalf("fixture premise broken, unusable roots: %v", u)
	}
	d := FSDeps{Paths: paths}
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(d) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return New(Options{
		Registry: reg, Paths: paths, Gate: g,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       func(format string, args ...any) { *audit = append(*audit, fmt.Sprintf(format, args...)) },
	}), g
}

// r4AuditLine is auditLineFor's sibling for a tool OTHER than fs.edit: this
// cell's claim is a comparison, so the fs.write line has to be caught the same
// way - by matching content, never by counting lines. wantSubs all have to
// appear in the SAME line.
func r4AuditLine(t *testing.T, audit []string, tool string, wantSubs ...string) string {
	t.Helper()
	var hits []string
	for _, l := range audit {
		if !strings.Contains(l, "tool="+tool) || !strings.Contains(l, "outcome=") {
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
	if len(hits) == 0 {
		t.Fatalf("wanted an audit line for %s carrying %v, got none in %q", tool, wantSubs, audit)
	}
	return hits[len(hits)-1]
}

func r4WriteArgs(t *testing.T, path, content string) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{"path": path, "content": content})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// TestFSEditAndFSWriteShareTheOutOfScopeVerdict is AC#5's 并排读数.
//
// ONE sample outside the authorized tree, fed to both tools. The sample is an
// EXISTING file so fs.edit has something to locate (its own precondition) and
// fs.write's overwrite fact (R8) is on the table too - otherwise the two calls
// would not be the same question, and a reading that differs for that reason
// would prove nothing about the family.
//
// What must come out identical: the level (L2), the rule that put it there
// (R2), the channel (the approval queue, never the L1 pre-execution window),
// the disposition (refused, class user_rejected), and the state of the disk
// (byte-for-byte what the caller wrote by hand before the call).
func TestFSEditAndFSWriteShareTheOutOfScopeVerdict(t *testing.T) {
	root := sealableTempCanonical124(t)
	outside := sealableTempCanonical124(t)
	const original = "outside the tree\n"
	// The target lives in a tree [fs] allowed_dirs never named.
	target := filepath.Join(outside, "note.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	var audit []string
	b, g := r4DeniedBridge(t, root, &audit)

	type reading struct {
		tool        string
		level       risk.Level
		rules       []risk.RuleID
		reason      string
		class       string
		text        string
		windowCalls int
		askedCalls  int
	}
	var got []reading

	// --- fs.write, out of scope -----------------------------------------------
	if out, err := b.Execute(t.Context(), req("fs.write",
		r4WriteArgs(t, target, "overwritten by the model"))); err != nil {
		t.Fatalf("fs.write Execute: %v", err)
	} else {
		if !out.IsError {
			t.Fatalf("an out-of-scope fs.write must not report success: %+v", out)
		}
		w, a := g.counts()
		dec := g.approvalDecision()
		got = append(got, reading{
			"fs.write", dec.Level, dec.RulesHit, dec.Reason,
			out.ErrorClass, out.Text, w, a,
		})
		if out.ErrorClass != "user_rejected" {
			t.Errorf("fs.write out of scope: ErrorClass=%q, want user_rejected (a refusal is not a host fault, SPEC-07 §2)", out.ErrorClass)
		}
		if dec.Level != risk.L2 {
			t.Errorf("fs.write out of scope: level=%v, want L2", dec.Level)
		}
		if !contains(dec.RulesHit, risk.R2) {
			t.Errorf("fs.write out of scope: rules_hit=%v, want R2 (C26's allowlist verdict)", dec.RulesHit)
		}
		if w != 0 || a != 1 {
			t.Errorf("fs.write out of scope routing: window=%d approval=%d, want 0/1", w, a)
		}
	}

	// --- fs.edit, the SAME sample --------------------------------------------
	if out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "outside", "new": "INSIDE"}))); err != nil {
		t.Fatalf("fs.edit Execute: %v", err)
	} else {
		if !out.IsError {
			t.Fatalf("an out-of-scope fs.edit must not report success: %+v", out)
		}
		w, a := g.counts()
		dec := g.approvalDecision()
		got = append(got, reading{
			"fs.edit", dec.Level, dec.RulesHit, dec.Reason,
			out.ErrorClass, out.Text, w, a,
		})
		if out.ErrorClass != "user_rejected" {
			t.Errorf("fs.edit out of scope: ErrorClass=%q, want user_rejected", out.ErrorClass)
		}
		if dec.Level != risk.L2 {
			t.Errorf("fs.edit out of scope: level=%v, want L2 (票面 09:4x 更正① 改判的那一枚)", dec.Level)
		}
		if !contains(dec.RulesHit, risk.R2) {
			t.Errorf("fs.edit out of scope: rules_hit=%v, want R2 - 越界判定必须由 C26 出，不是由本工具自报", dec.RulesHit)
		}
		// The delta is exactly the one call this test added, so both readings
		// land on the same channel: window stays 0, approval goes to 2.
		if w != 0 || a != 2 {
			t.Errorf("fs.edit out of scope routing: window=%d approval=%d cumulative, want 0/2 (the L1 window must never be what stops an edit)", w, a)
		}
	}

	// --- the side-by-side claim, stated as a claim ----------------------------
	for _, r := range got {
		t.Logf("并排读数 %s: level=%v rules_hit=%v window=%d approval=%d class=%q reason=%q",
			r.tool, r.level, r.rules, r.windowCalls, r.askedCalls, r.class, r.reason)
		t.Logf("并排读数 %s: 回执=%q", r.tool, r.text)
	}
	if got[0].level != got[1].level {
		t.Errorf("同族判失：fs.write 判 %v 而 fs.edit 判 %v", got[0].level, got[1].level)
	}
	if !contains(got[0].rules, risk.R2) || !contains(got[1].rules, risk.R2) {
		t.Errorf("同族判失：R2 必须两边都在，got %v / %v", got[0].rules, got[1].rules)
	}

	// The audit line is the half the reader who is not in this process sees.
	// Both tools must say the same thing about the same tree, and both must say
	// it with in_allowlist_scope=false.
	for _, tool := range []string{"fs.write", "fs.edit"} {
		line := r4AuditLine(t, audit, tool, "risk=L2", "R2", "in_allowlist_scope=false")
		t.Logf("并排审计行 %s：%s", tool, line)
		if !strings.Contains(line, "outcome=error") {
			t.Errorf("%s: the booked outcome must be error, got %q", tool, line)
		}
	}

	// --- and the disk says nothing happened ----------------------------------
	if data, err := os.ReadFile(target); err != nil {
		t.Fatal(err)
	} else if string(data) != original {
		t.Errorf("越界的两发调用动了盘：got %q want %q", data, original)
	}
	// No staging file may be left either: the refusal happens BEFORE the D31
	// writer, so a leftover .wisp-tmp-* under an unauthorized tree would mean
	// the gate let the writer run and only the rename failed.
	noStagingFilesLeft(t, outside)
}

// TestFSEditRoutesToApprovalNotTheL1Window pins the half a level string cannot
// say: WHICH channel a call lands on.
//
// fs.write's in-scope NEW file is the family's L1 shape - a 2-3s pre-execution
// window that expires INTO execution (SPEC-06 §2). fs.edit is D34's L2 row, so
// the same bridge must put it in the approval queue, where silence rejects. Both
// calls here answer Allow (the gate is not what this case is measuring - the
// ROUTING is), and each must touch exactly one channel.
func TestFSEditRoutesToApprovalNotTheL1Window(t *testing.T) {
	root := sealableTempCanonical124(t)
	const original = "alpha\nbravo\n"
	existing := filepath.Join(root, "has.txt")
	if err := os.WriteFile(existing, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
	b, _ := fsEditBridgeAuditedPair(t, root, g, nil)

	fresh := filepath.Join(root, "new.txt")
	if out, err := b.Execute(t.Context(), req("fs.write",
		r4WriteArgs(t, fresh, "created\n"))); err != nil {
		t.Fatalf("fs.write Execute: %v", err)
	} else if out.IsError {
		t.Fatalf("an approved in-scope fs.write must execute: %+v", out)
	}
	if w, a := g.counts(); w != 1 || a != 0 {
		t.Fatalf("in-scope fs.write (new file) must use the L1 window only: window=%d approval=%d", w, a)
	}

	if out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, existing,
		map[string]any{"old": "bravo", "new": "BRAVO"}))); err != nil {
		t.Fatalf("fs.edit Execute: %v", err)
	} else if out.IsError {
		t.Fatalf("an approved in-scope fs.edit must execute: %+v", out)
	}
	if w, a := g.counts(); w != 1 || a != 1 {
		t.Fatalf("fs.edit must reach the approval channel and never the L1 window: window=%d approval=%d (want 1/1 cumulative)", w, a)
	}
	if dec := g.approvalDecision(); dec.Level != risk.L2 {
		t.Errorf("the approval card carried level=%v, want the D34 L2 row", dec.Level)
	}
	if got := readString(t, existing); got != "alpha\nBRAVO\n" {
		t.Errorf("the approved edit did not land: %q", got)
	}
}

// fsEditBridgeAuditedPair is fsEditBridgeDeps with the gate handed in instead
// of hardcoded to "yes", and an optional audit sink (nil = discard, the shape
// fsEditBridge already is). One builder, one knob, so the routing case above
// and the refusal cases share every other property.
func fsEditBridgeAuditedPair(t *testing.T, root string, g Gate, audit *[]string) (*Bridge, *gateSpy) {
	t.Helper()
	spy, _ := g.(*gateSpy)
	paths := NewPathCanonicalizer([]string{root}, nil)
	d := FSDeps{Paths: paths}
	reg := NewRegistry()
	for _, e := range BuiltinFSEntries(d) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	logf := func(string, ...any) {}
	if audit != nil {
		logf = func(format string, args ...any) { *audit = append(*audit, fmt.Sprintf(format, args...)) }
	}
	return New(Options{
		Registry: reg, Paths: paths, Gate: g,
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       logf,
	}), spy
}

// TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord is the
// "摘掉那一道判定会红" case, written so a reader can see WHY it is load-bearing.
//
// fs.edit has no scope check of its own - it hands `path` to the judge through
// FSEditDecl's PathParams, and C19's R2 reads C26's InAllowlist. Drop that one
// declaration line and the judge sees NO path for this tool at all: the call
// still carries the declared L2 floor (so the level alone would stay green), but
// the out-of-scope fact, the reason, and the audit line's
// in_allowlist_scope=false all disappear - which is precisely the shape AC#5
// refuses to accept as "same family".
//
// The gate answers Allow here on purpose: with a "no" the call would stop at the
// gate and the tool's own fail-closed line would never be reached. The assertion
// that must survive the mutation is the VERDICT the card got, and the fact that
// the tool still refuses to open a path the judge never vouched for.
func TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord(t *testing.T) {
	root := sealableTempCanonical124(t)
	outside := sealableTempCanonical124(t)
	const original = "not authorized\n"
	target := filepath.Join(outside, "quiet.txt")
	if err := os.WriteFile(target, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}
	var audit []string
	g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
	b, _ := fsEditBridgeAuditedPair(t, root, g, &audit)

	if _, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, target,
		map[string]any{"old": "not", "new": "NOT"}))); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	dec := g.approvalDecision()
	if !contains(dec.RulesHit, risk.R2) {
		t.Errorf("the card for an out-of-scope edit carried rules_hit=%v, want R2 - 没有它就没有「越界」这一判", dec.RulesHit)
	}
	if !strings.Contains(dec.Reason, "授权目录之外") {
		t.Errorf("the reason shown to the human must name the allowlist verdict, got %q", dec.Reason)
	}
	r4AuditLine(t, audit, "fs.edit", "rules_hit=[", "R2", "in_allowlist_scope=false")
	// The tool then does write it - because a human said yes to the card that
	// named the tree. That is the whole design: the gate is the refusal, the
	// tool's own fail-closed leg is the second one, and the two are different
	// claims (see the next case).
	if got := readString(t, target); got != "NOT authorized\n" {
		t.Fatalf("premise: the approved edit must land, got %q", got)
	}
	t.Logf("读数：越界 fs.edit 经 R2 判 L2、卡片 reason=%q，人批准后落盘", dec.Reason)
}

// TestFSEditRefusesAPathC26CannotCanonicalize is the other direction of the same
// cell: 越界 is not only "outside the allowlist", it is also "we cannot answer
// what tree this is" - and the call must fail CLOSED, not fall through.
//
// The reading this case pins is the measurement, not a guess: a blank path never
// reaches the judge at all (pathArgs' stringValues filter drops it, so the
// booked line carries rules_hit=[R1] from the declared floor and
// in_allowlist_scope=true, which is "nothing was judged about a path", NOT
// "the path was judged in scope"). What stops the call is therefore the tool's
// own line (fs_edit.go's TrimSpace refusal, before any read or write), and this
// case is the pin that the second leg exists at all.
//
// The rules_gateway.go R2 "路径无法规范化，按越界处理（fail-closed）" leg is a
// DIFFERENT door - it needs a raw path that survives the filter and still fails
// C26, which on this tree means a reparse traversal, and that family already
// has its own尺 (bridge_junction_windows_test.go). Pinning it here would be a
// second ruler for one job, which is what AC#4b forbids for delete primitives.
func TestFSEditRefusesAPathC26CannotCanonicalize(t *testing.T) {
	root := sealableTempCanonical124(t)
	var audit []string
	g := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
	b, _ := fsEditBridgeAuditedPair(t, root, g, &audit)

	// A whitespace path: PathCanonicalizer.Canonicalize refuses it outright
	// ("tools: empty path"), so there is no canonical form to judge or open.
	if out, err := b.Execute(t.Context(), req("fs.edit", editArgs(t, "   ",
		map[string]any{"old": "a", "new": "b"}))); err != nil {
		t.Fatalf("Execute: %v", err)
	} else {
		if !out.IsError {
			t.Fatalf("a path C26 cannot answer must never report success: %+v", out)
		}
		if out.ErrorClass != "tool" {
			t.Errorf("ErrorClass=%q, want tool (the model can fix its own argument, D37)", out.ErrorClass)
		}
		if !strings.Contains(out.Text, "缺少 path 参数") {
			t.Errorf("the refusal must name what was missing, got %q", out.Text)
		}
		t.Logf("读数（未修码即响）：空白 path ⇒ IsError=true class=%q 回执=%q", out.ErrorClass, out.Text)
	}
	// And the judge never saw a path to vouch for: the booked line carries the
	// declared floor, not a scope verdict.
	line := r4AuditLine(t, audit, "fs.edit", "risk=L2")
	t.Logf("审计行：%s", line)
	// No file was created anywhere under the root by a call that named none.
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("a refused call left %d entries under the root: %v", len(entries), entries)
	}
}
