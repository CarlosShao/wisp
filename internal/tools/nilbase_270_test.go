// nilbase_270_test.go is ticket 270's instrument for form 甲 (ledger A990,
// decision 3): the call site that reads SubagentDeps.BaseOptions as a function
// refuses for itself when the field holds no function.
//
// What is measurable here and what is not, stated up front because the ticket's
// AC#1 judgement is a mutation and no test in this file can fake it:
//
//   - In the tree as committed, the assembly guard upstream of this call still
//     fires first, so a nil BaseOptions never reaches the call site. These legs
//     therefore pin the OBSERVABLES that both sites owe the caller - the one
//     named sentence, no child row, no slot held, no admission - and each of
//     those assertions is what turns red if the upstream guard is deleted and
//     the call site's own check is not there to take over (the pre-甲 reading of
//     that mutation is a panic, which this package's spawnGuarded reports as one
//     red whose text starts with "panic").
//   - The deletion itself is done by rewriting the production file on disk and
//     restoring it by sha256, not by `go test -overlay`: overlay replaces only
//     the bytes the COMPILER sees, and 321-r1 measured that the same trick is
//     invisible to every ruler that reads the source at runtime. The harness and
//     its readings live in .scratch/wisp/probes/270/r1/.
package tools

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// sentence270r1 is the assembly sentence as production writes it at both of its
// sites today. It is repeated here on purpose: the leg below compares it with
// the literal ticket 236 pinned, so the two copies in the tree cannot drift
// apart without one of the two rulers naming it.
const sentence270r1 = "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，fail-closed：不起第二套运行时）"

// AC#0's precondition, re-verified by this leg instead of inherited from 270-a2:
// the sentence 甲 reuses at the call site IS the sentence the upstream guard
// speaks, which in turn IS the string票 236 pinned. Byte equality, not "both
// mention 装配".
func Test270R1AssemblySentenceIsTheOne236Pinned(t *testing.T) {
	if sentence270r1 != lit236SpawnNoAssembly {
		t.Errorf("甲 形复用的那句与票 236 钉住的那句不是同一串：got %q, want %q",
			sentence270r1, lit236SpawnNoAssembly)
	}
}

// AC#1's reachable half: a BaseOptions that holds no function is refused with
// the named assembly sentence, not with a panic, not with a degraded runtime.
// Under the guard-removal mutation this same leg is the red-to-green ruler: with
// the call site still unprotected it reads as a panic, with 甲 in place it reads
// as the sentence below.
func Test270R1NilBaseOptionsRefusesWithTheAssemblySentence(t *testing.T) {
	h := newSub197Harness(t)
	h.buildWith(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true, func(d *SubagentDeps) { d.BaseOptions = nil })

	res := spawnGuarded(t, h, "装配快照没接线", "读一段")
	if strings.HasPrefix(res.Text, "panic") {
		t.Fatalf("调用点没自己判空（摘掉上游 guard 后的旧形）：%q", res.Text)
	}
	if !res.IsError {
		t.Errorf("BaseOptions 是空函数值时 task.spawn 回了非拒绝：%q", res.Text)
	}
	if res.Text != lit236SpawnNoAssembly {
		t.Errorf("空函数值没有走到那具名拒绝（got %q, want %q）", res.Text, lit236SpawnNoAssembly)
	}
	// Nothing about the refusal may be a degraded success: no child row, no
	// pool slot held past the return, no admission asked of the host.
	if n := h.roster.Count(); n != 1 {
		t.Errorf("名册行数 = %d, want 1（只有 harness 那枚根任务）：%v", n, h.childRows())
	}
	if got := h.roster.InFlightSubagents(); got != 0 {
		t.Errorf("被拒的这一枚占着池位：in-flight = %d, want 0（判空那支必须先把位还回去）", got)
	}
	h.muAdmit.Lock()
	admits := len(h.admits)
	h.muAdmit.Unlock()
	if admits != 0 {
		t.Errorf("空函数值却把子环路送去准入 %d 次，want 0（拒绝不能退化成另起一套运行时）", admits)
	}
}

// AC#3: the three refusals must each speak their own sentence and never stand in
// for one another. 236-v1 measured the substitution shape this pins against - a
// nil roster reaching the pool gate and being told "你排满了" - so the roster
// case is driven here too, four inputs, four reasons, pairwise distinct.
func Test270R1ThreeRefusalSentencesDoNotSubstitute(t *testing.T) {
	assembly := newSub197Harness(t)
	assembly.buildWith(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true, func(d *SubagentDeps) { d.BaseOptions = nil })
	gotAssembly := spawnGuarded(t, assembly, "装配没接线", "读一段").Text

	pool := newSub197Harness(t)
	pool.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true)
	releases := make([]func(), 0, MaxConcurrentSubagents)
	defer func() {
		for _, rel := range releases {
			rel()
		}
	}()
	for i := 0; i < MaxConcurrentSubagents; i++ {
		rel, ok := pool.roster.TryAcquireSubagentSlot(MaxConcurrentSubagents)
		if !ok {
			t.Fatalf("池没占满就只能占 %d 枚，want %d：本枚用例没有对象可派", i, MaxConcurrentSubagents)
		}
		releases = append(releases, rel)
	}
	gotPool := spawnGuarded(t, pool, "池已满", "读一段").Text

	orphan := newSub197Harness(t)
	orphan.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true)
	gotOrphan := spawnWithoutHostTaskID(t, orphan).Text

	noRoster := newSub197Harness(t)
	noRoster.buildWith(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true, func(d *SubagentDeps) { d.Roster = nil })
	gotNoRoster := spawnGuarded(t, noRoster, "名册没接线", "读一段").Text

	for _, c := range []struct{ name, got, want string }{
		{"装配未接线", gotAssembly, lit236SpawnNoAssembly},
		{"父任务 id 没给", gotOrphan, lit236SpawnNoHostTaskID},
		{"名册未接线", gotNoRoster, lit236SpawnRosterUnwired},
	} {
		if c.got != c.want {
			t.Errorf("%s 那一支说的不是自己那句（got %q, want %q）", c.name, c.got, c.want)
		}
	}
	for _, want := range []string{"拒绝派生", "上限", "不是排队"} {
		if !strings.Contains(gotPool, want) {
			t.Errorf("池满那一支说的不是池满那句（%q 里没有 %q）——"+
				"若 got 是装配或父任务那句，说明池满门在替别人说话", gotPool, want)
		}
	}
	// Pairwise distinctness: the four readings above each pass their own leg,
	// but a sentence rewritten to a shared one would still pass those; this is
	// the half that says the four may not be one.
	seen := map[string]string{gotAssembly: "装配未接线", gotOrphan: "父任务 id 没给", gotNoRoster: "名册未接线"}
	for _, other := range []struct{ text, name string }{{gotPool, "池满"}} {
		if prev, dup := seen[other.text]; dup {
			t.Errorf("%s 那一支与 %s 那一支说了同一句：%q（三句互不充当替身）", other.name, prev, other.text)
		}
		seen[other.text] = other.name
	}
	for _, pair := range []struct{ a, b string }{
		{gotAssembly, gotPool}, {gotAssembly, gotOrphan}, {gotPool, gotOrphan}, {gotNoRoster, gotPool},
	} {
		if pair.a == pair.b {
			t.Errorf("两枚不同的输入说了同一句拒绝：%q", pair.a)
		}
	}
	// The substitution 236-v1 actually measured: a nil roster must never be
	// answered by the pool gate.
	if strings.Contains(gotNoRoster, "上限") {
		t.Errorf("名册没接线却被池满门回答了：%q（want %q）", gotNoRoster, lit236SpawnRosterUnwired)
	}
	if strings.Contains(gotAssembly, "上限") || strings.Contains(gotAssembly, "父任务") {
		t.Errorf("装配没接线却被别的门回答了：%q（want %q）", gotAssembly, lit236SpawnNoAssembly)
	}
}

// AC#1's "not a silent degradation" clause, read from the other side: with the
// assembly wired, the same call must NOT take the refusal path. Without this
// control a judge cannot tell the legs above apart from a spawn that refuses
// everything.
func Test270R1WiredBaseOptionsIsNotRefused(t *testing.T) {
	h := newSub197Harness(t)
	h.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true)

	res := spawnGuarded(t, h, "正控那枚", "读完就好")
	if res.IsError {
		t.Errorf("装配齐全时 task.spawn 却拒了：%q", res.Text)
	}
	if strings.Contains(res.Text, lit236SpawnNoAssembly) {
		t.Errorf("正控读到了装配未接线那句，上面那几枚尺就可能是恒真：%q", res.Text)
	}
	if rows := len(h.childRows()); rows != 1 {
		t.Errorf("名册子代理行数 = %d, want 1：%+v", rows, h.childRows())
	}
	if got := h.roster.InFlightSubagents(); got != 0 {
		t.Errorf("正控收尾后 in-flight = %d, want 0", got)
	}
}
