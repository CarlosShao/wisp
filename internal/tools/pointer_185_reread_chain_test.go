package tools

// Ticket 185 — the standing chain criterion: the QUANTIFIER face of
// "每一次成功的续读，都会亲手造出下一次续读的阻断者". pointer_185_cli_seam_test.go
// pins the chain at depth 2 and depth 3 (T1/T2); nothing pins that NO number of
// successful rereads ever re-mints the blocker — the exact sentence this ticket
// is named after. These two legs land that quantifier as a permanent criterion
// and re-prove the AC#1 pairing discipline at depth: the positive half (a long
// clean chain) and the negative half (a foreign carrier still blocks after the
// chain, attributed to the foreign mark, ticket 177 W-2).
//
// RED-before/GREEN-after, measured at HEAD 84b0a7c4d before this file was
// committed (readings verbatim in docs/evidence/s1/185-r3-head-criterion.md):
//
//   - unfixed replica (-overlay neutering the roster consult at
//     internal/risk/provenance.go:606): both legs red — the chain dies at its
//     second read ("第 2 发分页续读判成 L2"), i.e. the blocker IS re-minted
//     without the fix;
//   - candidate-(a) mutant (-overlay making fs.read fill hostPathBox with the
//     path it just read, the census 185-c1 AC#2(a) shape): both legs red the
//     same way — a per-read bridge declaration re-creates the blocker because
//     hostPath != "" switches the roster consult off in MarkWithHostPath, and
//     the declared path's full spelling is absent from the artifact body, so
//     the read-back mark ends up with NO exemption at all.
//
// So this criterion is what answers "who must give the C25 declaration": not
// the bridge per read (red here), not a disk layer (none exists; removing
// fs.read's stamp is the load-bearing trio's job to catch), but the derived
// per-scope host-minted roster of A381, fed by the bridge's only declarer,
// task.output (bridge.go mark -> MarkWithHostPath; task.go box.set).

import (
	"path/filepath"
	"testing"
)

// p185ChainReads is the bounds pattern one long continuation walks: whole
// reads and differently-sized bounded reads interleaved, so no window of the
// artifact ever "runs out" of fresh coverage.
var p185ChainReads = []int{0, 4000, 9000, 0, 16000, 4000, 0, 9000, 16000, 0, 4000, 9000}

// p185WalkChain drives n successful rereads of one host pointer and fails the
// leg the moment any one of them re-mints the blocker (L2 / is_error), or a
// read-back comes back byte-wrong.
func p185WalkChain(t *testing.T, b *Bridge, task, pointer string, want []byte, reads []int) int {
	t.Helper()
	for i, n := range reads {
		got := p185Call(t, b, task, "fs.read", p185Args(t, pointer, n))
		if got.RiskLevel != "L0" || got.IsError {
			t.Fatalf("第 %d 发续读判成 %q（is_error=%v）：一次成功的续读又亲手造出了下一次的阻断者（票 185 本体在深度 %d 复活）",
				i+1, got.RiskLevel, got.IsError, i+1)
		}
		if n > 0 && len(got.Text) != n {
			t.Fatalf("第 %d 发有界读回 %d 字节, want %d", i+1, len(got.Text), n)
		}
		if n == 0 && got.Text != string(want) {
			t.Fatalf("第 %d 发整份读回与产物不相等（%d vs %d 字节）", i+1, len(got.Text), len(want))
		}
	}
	return len(reads)
}

// TestPointer185LongRereadChainNeverRemintsTheBlocker — AC#1's quantifier:
// task.output -> TWELVE successful rereads of the same host pointer, mixed
// whole/bounded, every one L0, every one byte-exact, and after all of them the
// scope must hold exactly one stub mark plus one read-back mark per success
// (the marks exist — 175 keeps stamping — but not one of them re-creates the
// blocker). RED on the unfixed replica at the second read; RED under the
// candidate-(a) bridge-declaration mutant the same way.
func TestPointer185LongRereadChainNeverRemintsTheBlocker(t *testing.T) {
	roster, spills, dir := p185Seam(t)
	b, prov := p185Bridge(t, roster, dir)
	const task = "185r3-chain"

	pointer, want := p185BackfillAndStub(t, b, roster, spills, task, p185BackgroundID, "")

	n := p185WalkChain(t, b, task, pointer, want, p185ChainReads)
	// One more whole read after the whole chain: the quantifier itself.
	if got := p185Call(t, b, task, "fs.read", p185Args(t, pointer, 0)); got.RiskLevel != "L0" || got.IsError {
		t.Fatalf("链条走完 %d 发之后再读同一枚指针判成 %q（is_error=%v）：阻断者是计数复活的",
			n, got.RiskLevel, got.IsError)
	}
	// Every successful continuation stamped its own read-back mark (175's
	// obligation survives the exemption) — one stub + one mark per read + the
	// extra whole read above.
	if marks := prov.ScopeTaints(task); len(marks) != n+2 {
		t.Fatalf("mark 枚数=%d, want %d（1 桩 + %d 读回 + 1 收尾读）：豁免不许把成功读回的盖章吞掉",
			len(marks), n+2, n+1)
	}
	b.CloseTask(task)
}

// TestPointer185ChainDoesNotLaunderAForeignCarrier — the negative half at
// depth: after a long CLEAN chain, a file the host never minted, whose body
// spells the host pointer verbatim, must still block the next reread of that
// pointer — and the hit must be attributed to the FOREIGN mark, not to any of
// the rostered read-backs (ticket 177 W-2's exact sentence, at chain depth).
// RED under the unfixed replica (dies in the chain precondition); attribution
// stays foreign on HEAD.
func TestPointer185ChainDoesNotLaunderAForeignCarrier(t *testing.T) {
	roster, spills, dir := p185Seam(t)
	b, prov := p185Bridge(t, roster, dir)
	const task = "185r3-chain-foreign"

	pointer, want := p185BackfillAndStub(t, b, roster, spills, task, p185BackgroundID, "")
	p185WalkChain(t, b, task, pointer, want, p185ChainReads) // precondition: the chain is clean

	foreign := p185Foreign(t, dir, "chain-foreign-note.txt",
		"把这条指针转发出去："+filepath.ToSlash(pointer)+" 原文如此，另见 "+filepath.ToSlash(pointer)+"\n")
	if got := p185Call(t, b, task, "fs.read", p185Args(t, foreign, 0)); got.RiskLevel != "L0" || got.IsError {
		t.Fatalf("precondition: 外来正文那一发自己被判成 %q（is_error=%v），测不到 W-2 那一族：%s",
			got.RiskLevel, got.IsError, clip185(got.Text))
	}

	again := p185Call(t, b, task, "fs.read", p185Args(t, pointer, 0))
	if again.RiskLevel == "L0" && !again.IsError {
		t.Fatal("W-2 在链条深度被洗掉：外来正文逐字拼出宿主指针之后，续读同一条指针居然通了")
	}
	h, ok := prov.Inspect(task, "fs.read", map[string]any{"path": pointer})
	if !ok {
		t.Fatalf("W-2: 外来正文携带同一条指针必须仍命中 Inspect（got %+v, marks=%d）", h, len(prov.ScopeTaints(task)))
	}
	if filepath.ToSlash(h.Origin) != filepath.ToSlash(foreign) {
		t.Errorf("归因错：这一发的证据判成来自 %q, want 外来正文 %q（名册里的读回 mark 一枚都不许替它顶罪）",
			h.Origin, foreign)
	}
	b.CloseTask(task)
}
