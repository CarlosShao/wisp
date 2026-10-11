package tools

// 票 309 AC#1 —— "Bridge.assessorFor 的 taskID == \"\" 那一肢把 R4 探测器静默摘掉"
// 这一件事在本包必须有**恒真形必须响**的仪器（形 Ⓑ：写明故意、行为逐字节不变；
// 所以"响"只能靠这枚仪器＋ toolrequest_nail_309_test.go 的枚数钉，别处没有牙）。
//
// 判据（两向都有牙，缺一即那格没交）：
//   ① 非空 task id：assessorFor 走 WithTaintDetector 那一支 ⇒ 对同一句敏感文本必须
//      判出 R4（探测器真挂上了）。把 `.WithTaintDetector(b.prov.Detector(taskID))`
//      换成 `.WithTaintDetector(nil)`（探测器没挂）时，这一腿必须红。
//   ② 空 task id：assessorFor 命中 `taskID == ""` 那一肢、返回 b.assess（⛔ 探测器）⇒
//      对同一句敏感文本必须判不出 R4（今天对空 id 的实际防御＝零）。把守卫
//      `b.injected || b.prov == nil || taskID == ""` 里的 `|| taskID == ""` 摘掉（或整枚
//      摘掉守卫）时，空 id 也去挂探测器、撞成 unbound-scope fail-closed R4，这一腿必须红。
//
// 判据锚⛔ 落在"返回值里的字样"（读返回值会被句子里写了成功骗过去），落在**探测器挂没挂**
// 这件事实上：现成可观察面＝经 assessorFor 拿回的判级器对同一句污点文本的 R4 差异。
// ⛔ 新造导出 API、⛔ 加字段；只在包内直接驱动 assessorFor+Assess（无新增生产可达路径）。

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// hasRule309 reports whether a fused decision carries a specific rule id.
func hasRule309(d risk.Decision, want risk.RuleID) bool {
	for _, r := range d.RulesHit {
		if r == want {
			return true
		}
	}
	return false
}

func Test309AssessorForEmptyTaskIDDropsTaintDetector(t *testing.T) {
	root := sealableTempCanonical124(t)
	b, prov := fsBridgeWith(t, nil, nil, root)

	// One real taint mark in the engine, under a NON-EMPTY scope id. Mark stamps
	// p.scopes[id] even with no OpenScope (provenance.go MarkWithHostPath: a
	// never-opened scope still records the taint), which is all the detector needs.
	const taintedTask = "task-309-tainted"
	const sentinel = "sentinel309taintsource"
	if !prov.Mark(taintedTask, risk.SrcWebFetch, "https://leak.example/report.md", "本机凭据摘录 "+sentinel+" 收尾") {
		t.Fatalf("前置失败：Mark 没有落污点，本判据失去前提")
	}

	params := map[string]any{"text": "外发正文包含 " + sentinel}

	// ① 非空 id —— 探测器挂上 ⇒ R4 必须响。
	gotTainted := b.assessorFor(taintedTask).Assess("notify", params, risk.Facts{})
	if !hasRule309(gotTainted, risk.R4) {
		t.Errorf("非空 task id 的判级器没判出 R4（level=%v rules=%v）："+
			"assessorFor 的 WithTaintDetector 那一支若被摘掉，R4 就整条不发生；"+
			"① 这一腿的存在就是为了：探测器真没挂上时这里必须红",
			gotTainted.Level, gotTainted.RulesHit)
	}

	// ② 空 id —— 命中 taskID == "" 那一肢、返回 b.assess（⛔ 探测器）⇒ R4 必须不响。
	// 这记录的是形 Ⓑ 的实话：对空 id 的实防＝零，门照走、只是 R4 那一支不发生。
	gotEmpty := b.assessorFor("").Assess("notify", params, risk.Facts{})
	if hasRule309(gotEmpty, risk.R4) {
		t.Errorf("空 task id 竟判出了 R4（rules=%v）：assessorFor 的 `taskID == \"\"` 那一肢"+
			"被摘掉／改成\"空 id 也挂探测器\"时才会这样（空 id 撞 unbound-scope fail-closed）。"+
			"② 这一腿的存在就是为了：恒真形（守卫被削弱）必须在这里红",
			gotEmpty.RulesHit)
	}
	if hasRule309(gotEmpty, risk.R4) == hasRule309(gotTainted, risk.R4) {
		t.Errorf("两向没有拉开（空/非空同为 %v 或同不为 R4）：这枚仪器对\"探测器挂没挂\"是瞎的",
			gotEmpty.RulesHit)
	}
}
