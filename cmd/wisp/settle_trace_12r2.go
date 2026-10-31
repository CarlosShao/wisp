package main

import (
	"sync"
	"sync/atomic"
	"time"
)

// 票 12 AC#9（registry A25）的「最后 token 时间」载体——3s 回落实测的三个单调钟戳。
//
// AC#9 要的量是「从最后一个 token 送达起算，到回落（Settling→Sleeping 一族）的时延」，
// ≥5 个样本、单调钟、逐样本毫秒全报。状态机事件面是冻结面（events.go 里没有任何
// 时间载体，本缝一个字节不加），所以载体放 resident/session 侧，三个戳各挂一条
// 已经存在的缝：
//
//   - T0 最后一个 token 送达 = execute() 的流式 sink 收到本任务最后一条 delta 的
//     时刻（consoleSink.Publish 的 EvTextDelta 分支，run.go）。这是 LLM 流式回复
//     交付给消费面的最后一跳，也是 AC#9 字面的「最后一个 token 送达」。
//   - T1 进 Settling = session scope 自己宣布 SCOPE-SETTLING 的那一刻（openScope
//     的 Logf 钩子）。这是 D43#31 的 session 半边（Warm 驻留到期 → 3s 静置窗开）。
//     球那台机器的 Warm 驻留走冻结 DefaultTimeouts（90s），且机器 Sink 的效果投递
//     面归票 273 收口（273 AC#2 三形未选形）——本缝两头都不碰。
//   - T2 回落完成 = scope 的 dispose 级联跑 FreeOSMemory 钩子的那一刻（openScope
//     已接线的级联第三步）。onSettleExpired 是 D43#32 的触发器，3s 静置窗是冻结
//     默认 SettleStaticDefault，本缝只观测、不改窗。
//
// 时间全部取 time.Since(settleTraceBase)：time.Now 的单调钟读数（AC#9 指名的
// 单调钟路；这是测量 stamps，不是超时判定，D22 禁令 #4 管的是后者）。
//
// 本件只记录发生了什么，判线（≤3000ms、3s 窗）归测试与证据表——载体不带阈值。
type settleTrace struct {
	// lastTokenNS is T0, monotonic ns; 0 means "no reply streamed yet".
	lastTokenNS atomic.Int64

	mu      sync.Mutex
	pending settlePending
	samples []settleSample
}

// settlePending is the T0+T1 pair waiting for its T2. A SCOPE-SETTLING without a
// streamed reply (a voice turn that never reached the text path) stays invalid:
// no reply, no sample.
type settlePending struct {
	valid       bool
	lastTokenNS int64
	settledNS   int64
}

// settleSample is one measured fall-back, in monotonic-clock durations.
type settleSample struct {
	// LastTokenToSettling is T1-T0: AC#9①'s leg, judged against 3000ms.
	LastTokenToSettling time.Duration
	// SettlingToFallback is T2-T1: the settle window itself (the frozen 3s
	// static plus dispatch overhead) — the「3s 回落」the record is about.
	SettlingToFallback time.Duration
	// LastTokenToFallback is T2-T0: the whole after-reply fall-back.
	LastTokenToFallback time.Duration
}

// settleTraceBase is the epoch every stamp is read against. Initialized once at
// package load and never written again; time.Since reads the monotonic clock.
var settleTraceBase = time.Now()

func monoNS() int64 { return int64(time.Since(settleTraceBase)) }

// markLastToken stamps T0. Every streamed delta passes here; only the last one
// before T1 survives into a sample, so no per-task bookkeeping is needed.
func (t *settleTrace) markLastToken() {
	if t == nil {
		return
	}
	t.lastTokenNS.Store(monoNS())
}

// markSettling stamps T1 and opens the pending sample. Called from the scope's
// own settling announcement; a previous pending (a settle a wake cancelled) is
// replaced, never completed from stale halves.
func (t *settleTrace) markSettling() {
	if t == nil {
		return
	}
	now := monoNS()
	t0 := t.lastTokenNS.Load()
	t.mu.Lock()
	if t0 > 0 && t0 <= now {
		t.pending = settlePending{valid: true, lastTokenNS: t0, settledNS: now}
	} else {
		t.pending = settlePending{}
	}
	t.mu.Unlock()
}

// markFellBack stamps T2 and closes the pending sample, if one is open. Called
// from the dispose cascade; a cascade without a settling before it (a family
// swap) closes nothing.
func (t *settleTrace) markFellBack() {
	if t == nil {
		return
	}
	now := monoNS()
	t.mu.Lock()
	p := t.pending
	t.pending = settlePending{}
	if p.valid && p.settledNS <= now && p.lastTokenNS <= p.settledNS {
		const settleSampleKeep = 64
		t.samples = append(t.samples, settleSample{
			LastTokenToSettling: time.Duration(p.settledNS - p.lastTokenNS),
			SettlingToFallback:  time.Duration(now - p.settledNS),
			LastTokenToFallback: time.Duration(now - p.lastTokenNS),
		})
		if len(t.samples) > settleSampleKeep {
			t.samples = t.samples[len(t.samples)-settleSampleKeep:]
		}
	}
	t.mu.Unlock()
}

// samplesSnapshot hands out the recorded samples in order.
func (t *settleTrace) samplesSnapshot() []settleSample {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]settleSample(nil), t.samples...)
}
