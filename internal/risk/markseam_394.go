// File: internal/risk/markseam_394.go — 票 394: mark 观察缝（第三形，机主 A1362 批「394 按第三形批」）
//
// Status: 实施中（394-r1 写腿 2026-11-01）
// Blocks: nothing · Blocked by: nothing
//
// 第三形＝risk 层 mark 观察缝：RecordTaintInheritance 的契约（tts_gate.go:245
// 「调用方在 mark 时累积归一化内容」）在生产里没有可履行的调用方——生产 4 处
// 带内容的 mark 点（internal/tools/bridge.go:687、internal/tools/subagent_197.go:512、
// cmd/wisp/panel_assets.go:245、internal/projctx/projctx.go:349）跨 3 包，而
// normalizeTaint 是本包未导出函数，任何外部调用方都无法履行归一化义务
// （382-w-r2 的测试被迫手选已归一化常量即此痛）。本缝把义务收回 risk 内部：
// MarkWithHostPath（全部 mark 点的汇流口）在 mark 落账后把【已归一化】的内容
// 通知给本引擎上注册的观察者，TTSGate 在 NewTTSGate 时借引擎自注册（382-w-r2
// 已落的 attachProvenance 接线零改动即生效）。
//
// 三条硬边界（票 394 派单）：
// 1. ⛔ 改 CheckText/Scope.Close/RecordTaintInheritance 既有契约 —— 本文件只新增面；
//    RecordTaintInheritance 保持逐字节不动（显式追加语义归它的既有调用方/测试）。
// 2. ⛔ internal/tools/** —— 第三形的全部意义＝不动 tools 契约（4 处 mark 点零触碰）。
// 3. 契约面只在 risk —— 观察者接口与注册路由全部未导出，包外无法接线也无法误用。

package risk

import (
	"github.com/CarlosShao/wisp/internal/observe"
)

// taintMarkObserver is the ticket 394 mark-observation seam: one engine's
// notified view of every mark it just committed. The interface is unexported
// on purpose — the seam's contract face is this package, and a type outside
// risk can neither register itself nor be registered by the composition root,
// which is what keeps the third form's promise ("契约面只在 risk").
//
// Implementors are called AFTER the mark is committed to its scope (never
// under p.mu) with the content exactly as the engine indexed it: already
// through normalizeTaint, truncated to MaxSourceRunes. The obligation
// RecordTaintInheritance documents for its caller — "accumulate normalized
// content at mark time" — is fulfilled here, inside risk, where the
// normalizer actually lives.
//
// An observer must not call Mark/MarkWithHostPath back on the notifying
// engine from inside the callback: the notify runs after the engine lock is
// released, but a re-entering mark would recurse mark→notify→mark.
type taintMarkObserver interface {
	observeTaintMark(scopeID, tool, origin, normContent string)
}

// addMarkObserver registers o on THIS engine only. Engine-local is the point:
// a gate observes the engine it borrowed (382-w-r2's "same engine the bridge
// marks into" rule), so a private engine — panel-assets' probe, a fail-closed
// boot engine — can never leak marks into another engine's gate.
func (p *Provenance) addMarkObserver(o taintMarkObserver) {
	if p == nil || o == nil {
		return
	}
	p.mu.Lock()
	p.markObs = append(p.markObs, o)
	p.mu.Unlock()
}

// notifyTaintMark hands one committed mark to every observer of this engine.
// The observer slice is copied under the lock and called OUTSIDE it, so an
// observer taking its own lock (TTSGate's) can never invert the lock order
// against the engine's. Zero observers — every engine without a gate, the
// common case — costs one len() read and no allocation.
func (p *Provenance) notifyTaintMark(scopeID, tool, origin, normContent string) {
	if p == nil {
		return
	}
	p.mu.RLock()
	var obs []taintMarkObserver
	if n := len(p.markObs); n > 0 {
		obs = make([]taintMarkObserver, n)
		copy(obs, p.markObs)
	}
	p.mu.RUnlock()
	for _, o := range obs {
		o.observeTaintMark(scopeID, tool, origin, normContent)
	}
}

// observeTaintMark implements taintMarkObserver on the TTSGate: every mark the
// borrowed engine commits is folded into the inheritance chain this gate scans
// at broadcast time (CheckTextWithInheritance's form B). This is the production
// caller side of RecordTaintInheritance's contract — with the content already
// normalized by the mark path, which no caller outside risk could deliver.
//
// Retention posture (the owner-approved cost of the third form: "留存策略从显式
// opt-in 变引擎侧姿势"): the chain holds ONE link per scope id, appended into
// as that scope keeps marking. Per-mark appends were rejected on the landed
// scan contract's own terms — checkAncestorScopes scans only the FIRST
// MaxAncestorDepth links, so one link per mark would push every mark after the
// 16th out of the scan window for good; one link per scope keeps a whole
// session's worth of tasks inside it. The 382-w-v2 suggestion to collect only
// task-scope marks is answered the same way: risk has no task-scope
// discriminator (production ids are loop-minted, "host:*", "run:*" — no
// shared shape), and an id-prefix table here would be a new fragile contract;
// upsert-per-scope bounds the chain to O(distinct marked scopes), which is the
// retention bound the suggestion was after. Scope ids that never close (the
// project-instructions scope) keep their taints live in their own scope, but
// the announce gate's direct check only reads the announce scope — this chain
// is the ONLY route by which those taints can reach a broadcast decision at
// all, which is why nothing is filtered out here.
//
// normContent arrives normalized (the notify contract above); it is stored
// as-is, the same coordinate space matchTaintsInLink compares in.
func (g *TTSGate) observeTaintMark(scopeID, tool, origin, normContent string) {
	if g == nil || !g.config.EnableTaintInheritance || scopeID == "" {
		return
	}
	if len([]rune(normContent)) == 0 {
		return
	}
	g.mu.Lock()
	var link *taintInheritanceLink
	for _, l := range g.inherited {
		if l.TaskScopeID == scopeID {
			link = l
			break
		}
	}
	if link == nil {
		link = &taintInheritanceLink{TaskScopeID: scopeID}
		g.inherited = append(g.inherited, link)
	}
	link.TaintContent = append(link.TaintContent, normContent)
	link.ToolSources = append(link.ToolSources, tool)
	link.MarkedAt = observe.NowWallUTC().Unix()
	n := len(link.TaintContent)
	g.mu.Unlock()
	logf("risk/C25: 394 seam recorded mark %s origin=%q into scope %q's inheritance link (%d taints)", tool, origin, scopeID, n)
}
