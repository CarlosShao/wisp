package tools

// 票 309 —— leg 309b：给 assessorFor 的守卫三枚前提各装一枚**现状仪器**。
//
// 这枚件只描述 HEAD 的真实行为并把它标成缺口，⛔ 断言"生产应然"：
// 票面（`.scratch/wisp/issues/309-assessorfor-empty-taskid-limb-silently-drops-r4-detector.md`）
// 现量节逐字写着 `bridge.go:925` ＝ "静默摘掉探测器、零注释"，而 `AC#0` 裁形 Ⓑ
// ＝ "写明故意、行为逐字节不变" ⇒ 对空 id 的实际防御＝零。下面三腿把这件事
// 钉成读数，每一腿的红句都指向缺口本身，⛔ 指向"生产该有的样子"。
//
// 三腿共用的可观测点＝**assessorFor 拿回来的判级器对同一句污点文本的 R4 差异**，
// 与 assessorfor_taint_instrument_309_test.go 同一把尺（⛔ 读返回值里的字样，
// 读它会被句子里写了成功骗过去）。两枚分支在一件事上必然不同：
//   - 走 `WithTaintDetector(b.prov.Detector(taskID))` 那一支：探测器**在**，
//     所以空 scope ‡ 未登记、而别的 scope 手里有污点时，R4 以 unbound-scope
//     fail-closed 响（`internal/risk/provenance.go` scopeMarks 的第三支）；
//   - 走 `return b.assess` 那一支：探测器**没挂**，R4 整条不发生。
//
// ⇒ "R4 响 ⇔ 走了挂探测器那一支"，这是两条分支的可区分面，⛔ 是猜测。

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// probe309GuardLimb 是一条守卫肢的读数：同一次 wiring 下，assessorFor(空 id)
// 与 assessorFor(非空 id) 各自判出没判出 R4。
type probe309GuardLimb struct {
	emptyR4    bool // assessorFor("") 是否判出 R4
	nonEmptyR4 bool // assessorFor(非空 id) 是否判出 R4
	emptyRules []risk.RuleID
	emptyNote  string
}

// probe309Guard drives one wiring (a b *Bridge whose prov already holds a taint
// mark under a DIFFERENT, non-empty scope id) and reports what the two limbs do.
//
// 前提的构造（两条分支都吃同一份）：
//   - 引擎里真有一枚污点，落在另一枚非空 scope 上 ⇒ 任何**未登记** scope 的
//     Inspect 都会 fail-closed 回 unbound-scope；
//   - 被问的那个 id（空串或另一枚非空串）**从没登记过** ⇒ 谁挂了探测器，
//     谁就一定判出 R4。于是"判出 R4"＝"探测器挂着"，"判不出"＝"探测器被摘掉"。
func probe309Guard(t *testing.T, b *Bridge, prov *risk.Provenance, emptyID string, attached bool) probe309GuardLimb {
	t.Helper()
	const taintedScope = "task-309-guard-tainted"
	const sentinel = "sentinel309guardsource"
	if !prov.Mark(taintedScope, risk.SrcWebFetch, "https://leak.example/guard.md", "本机凭据摘录 "+sentinel+" 收尾") {
		t.Fatalf("前置失败：Mark 没有落污点，本判据失去前提")
	}
	params := map[string]any{"text": "外发正文包含 " + sentinel}

	gotEmpty := b.assessorFor(emptyID).Assess("notify", params, risk.Facts{})
	// 非空对照：一枚同样没登记、同样非空的 id。它走挂探测器那一支时也必判 R4，
	// 所以它证明"这台面本身能判出 R4"＝空 id 那边判不出⛔ 是台面哑了。
	gotNonEmpty := b.assessorFor("task-309-guard-unregistered").Assess("notify", params, risk.Facts{})

	out := probe309GuardLimb{
		emptyR4:    hasRule309(gotEmpty, risk.R4),
		nonEmptyR4: hasRule309(gotNonEmpty, risk.R4),
		emptyRules: gotEmpty.RulesHit,
	}
	if attached {
		out.emptyNote = "探测器挂着（R4 会响）"
	} else {
		out.emptyNote = "探测器被摘掉（R4 不发生）"
	}
	return out
}

// Test309bAssessorForGuardLimbsReadAsBaseline 三腿一次跑完，每腿一枚 t.Run，
// 每条断言的失败文案都自带票面行号与 EXPECTED-GAP 标记。
func Test309bAssessorForGuardLimbsReadAsBaseline(t *testing.T) {
	root := sealableTempCanonical124(t)

	// ---- 腿 1：空 taskID + prov 非空 + 非 injected ----------------------
	// 票面依据（逐字）：`bridge.go:925` ＝ `if b.injected || b.prov == nil || taskID == "" {`
	// / `return b.assess`，`109` 那行"三枚前提任一成立就整枚摘掉 R4 探测器
	// （`:931 WithTaintDetector(b.prov.Detector(taskID))` 那一支不走）"。
	t.Run("empty-taskID-drops-detector-EXPECTED-GAP", func(t *testing.T) {
		b, prov := fsBridgeWith(t, nil, nil, root)
		if b.injected || b.prov == nil {
			t.Fatalf("前置失败：这一腿要的是 injected=false 且 prov 非空（现量 injected=%v prov!=nil=%v）",
				b.injected, b.prov != nil)
		}
		got := probe309Guard(t, b, prov, "", false)

		// 台面自证：同一台面上非空 id 判得出 R4 ⇒ 空 id 判不出⛔ 是"探测器没挂"。
		if !got.nonEmptyR4 {
			t.Fatalf("前置失败：非空对照 id 没判出 R4，这台面读不出探测器挂没挂（读不到区分面）")
		}
		// ★EXPECTED-GAP（本票待补结构保证）：空 taskID 走 `return b.assess`
		// ⇒ R4 静默不发生。把这一肢改成"空 id 也挂探测器"时这条必红——
		// 红就是本票 AC#1 要的那声；今天它是绿的，绿＝缺口在那里。
		if got.emptyR4 {
			t.Errorf("EXPECTED-GAP 被翻掉了：空 taskID 竟判出 R4（rules=%v）-> `taskID == \"\"` "+
				"那一肢不再走 b.assess、探测器被挂上了。本票（票 309）正是要暴露这一肢今天"+
				"静默摘探测器；若这是有意修的，请同步改票面现量与 AC#1 的仪器", got.emptyRules)
		}
		t.Logf("读数［腿1］空 taskID：R4=%v rules=%v（%s）；非空对照：R4=%v",
			got.emptyR4, got.emptyRules, got.emptyNote, got.nonEmptyR4)
	})

	// ---- 腿 2：injected=true -------------------------------------------
	// 票面依据（逐字）：`:925` 那个布尔式的第一枚前缀 `b.injected`；
	// 票面现量节另记"`b.injected` 那一肢同样不挂探测器，但生产面 `Assessor:` 赋值点
	// ＝0 枚 ⇒ `injected` 今天恒 false，归装配那族另裁"。
	// 装配方式：`New(Options{Assessor: ...})` 是**唯一**置 injected 的写点
	// （`bridge.go:209` `b.assess, b.injected = o.Assessor, true`），所以这里照它装。
	t.Run("injected-assessor-drops-detector-EXPECTED-GAP", func(t *testing.T) {
		injected := risk.NewRiskAssessor().
			WithCanonicalizer(NewPathCanonicalizer(nil, nil)).
			WithSensitiveClassifier(NewSensitiveClassifier())
		prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil})
		b := New(Options{
			Registry: NewRegistry(), Paths: NewPathCanonicalizer(nil, nil),
			Provenance: prov, Assessor: injected, Gate: NoGate{},
			Logf: func(string, ...any) {},
		})
		if !b.injected {
			t.Fatalf("前置失败：Options.Assessor 没把 injected 置起来（bridge.go:209）")
		}
		got := probe309Guard(t, b, prov, "", false)

		// 非空对照在这台面上同样走 b.assess（injected 那一肢对任何 id 都成立），
		// 所以它判不出 R4 ⛔ 是台面哑，而是同一肢的第三个读数；断言它也是缺口的一部分。
		if got.nonEmptyR4 {
			t.Errorf("injected=true 时非空 id 竟判出 R4（rules=%v）-> `b.injected` 那一肢已不再整枚摘探测器",
				got.emptyRules)
		}
		// ★EXPECTED-GAP（本票待补结构保证）：injected=true ⇒ 无论 id 空不空，一律
		// 走 `return b.assess` ⇒ R4 不发生。今天生产面没有写点，所以这条是"装配坏了
		// 就静默"那一族的地板；票 309 AC#0 裁⛔ 给这一肢造第二枚仪器，这行是登记。
		if got.emptyR4 {
			t.Errorf("EXPECTED-GAP 被翻掉了：injected=true 时空 taskID 判出 R4（rules=%v）-> "+
				"`b.injected` 那一肢的形状变了", got.emptyRules)
		}
		t.Logf("读数［腿2］injected=true：空 id R4=%v rules=%v；非空 id R4=%v（%s）",
			got.emptyR4, got.emptyRules, got.nonEmptyR4, got.emptyNote)
	})

	// ---- 腿 3：prov == nil ---------------------------------------------
	// 票面依据（逐字）：`:925` 那个布尔式的第二枚前缀 `b.prov == nil`；
	// 票面现量节"`b.prov == nil` … 同族同处置：属'装配坏了'"。
	// 装配方式：Options.Provenance 留空 ⇒ `b.prov == nil`；其余照 fsBridgeWith 装。
	t.Run("nil-provenance-drops-detector-EXPECTED-GAP", func(t *testing.T) {
		reg := NewRegistry()
		b := New(Options{
			Registry: reg, Paths: NewPathCanonicalizer(nil, nil), Provenance: nil,
			Gate: NoGate{}, Logf: func(string, ...any) {},
		})
		if b.prov != nil {
			t.Fatalf("前置失败：Options.Provenance 留空后 b.prov 却非空")
		}
		// 这一台面**装不了**"别的 scope 持污点"那枚前提（没有引擎），所以
		// "挂上探测器"在这里不可观测：探测器只能挂在 nil 引擎上，第一次问就 panic
		// -> fail-closed R9。★这正是缺口本身的一部分：prov==nil 那一肢**没有
		// 可区分面**（探测绑定不可从外观测），所以本腿只能断言"返回非 nil 且
		// 不会因此判出 R4"这一件可观测的事，⛔ 能区分两条分支。
		got := b.assessorFor("").Assess("notify", map[string]any{"text": "任意正文"}, risk.Facts{})
		assign := b.assessorFor("task-309-nonempty")
		if assign == nil {
			t.Fatalf("assessorFor 返回 nil：调用点 bridge.go:297 会立刻解引用， 这是今天的行为")
		}
		// ★EXPECTED-GAP（本票待补结构保证）：prov==nil ⇒ 走 `return b.assess`，
		// 而那枚 b.assess 从来就没挂过探测器（`bridge.go:211-213` 只挂 canonicalizer
		// 与 classifier）-> R4 不发生。把 `b.prov == nil` 摘掉时，这一腿会撞成
		// nil 引擎上的 Detector → R9，红在那条；今天它是绿的。
		if hasRule309(got, risk.R4) {
			t.Errorf("EXPECTED-GAP 被翻掉了：prov==nil 时空 taskID 判出 R4（rules=%v）", got.RulesHit)
		}
		if hasRule309(got, risk.R9) {
			t.Errorf("prov==nil 走 b.assess 却撞出 R9（rules=%v）：说明这一肢已不再返回 b.assess、"+
				"改去碰 nil 引擎了", got.RulesHit)
		}
		t.Logf("读数［腿3］prov=nil：空 id R4=%v rules=%v；assessorFor(非空) 非 nil=%v；"+
			"注意：探测绑定不可从外观测＝缺口本身（仅文档化， 能机器区分两分支）",
			hasRule309(got, risk.R4), got.RulesHit, assign != nil)
	})
}
