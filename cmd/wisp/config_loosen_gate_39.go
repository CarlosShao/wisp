package main

// Ticket 39 AC#2's assembly half (write leg 2026-11-18): the PRODUCTION
// adapter that turns Manager.LoosenQueue's seam - ticket 195's C18 hand-off,
// nil-refused fail-closed since then - into a REAL native L2 card, and turns
// the card's answer into the write it was raised for.
//
// WHY THIS FILE EXISTS. Ticket 195 deliberately stopped at the seam: a native
// loosening was queued and nothing was written (fail-closed), and its evidence
// (docs/evidence/s1/195-r1-direction-verdict.md §7) named the missing piece as
// exactly this wiring - "把 internal/agent/approval 的 C18 队列适配成
// config.LoosenQueue" at the assembly root. Until it exists, ticket 39 AC#2's
// sentence (security loosening -> native L2 confirm -> accept applies + logs;
// reject keeps old) has an engine with no wire: SetPluginEnabled's native ON
// (ticket 215) answered ErrLoosenQueueUnwired in every production assembly.
//
// THE CARD IS THE EXISTING ONE, NOT A SECOND KIND. One card through the same
// C18 gate a tool call goes through (gate.PendingApproval, risk.L2), the same
// synthetic-task admission D47 requires (gate.AdmitTextTask, revoked on the
// way out), the same correlation-id minting shape confirmModeSwitchVia uses -
// except the id is minted UPSTREAM (config.Manager mints
// LoosenRequest.RequestID and registers the pending write under it before the
// queue is called), so the card's correlation id, the pending-write registry
// key and the audit join are one value by construction. The card can only be
// answered from the native side (SPEC-06:1588, ban #6); nothing here adds a
// panel-side allow, and the panel-sourced loosening refusal upstream
// (ErrPanelLoosenRefused) is untouched.
//
// WHY NOTHING BLOCKS (ruling N9). QueueLoosen spawns the card goroutine and
// returns the id immediately: the 300s C18 deadline must not sit inside a
// config write. The card's answer lands on the spawned goroutine, which calls
// Manager.ResolveQueuedLoosen - accepted: the real write, one audit line with
// the card id and the old->new values; refused or timed out (the gate resolves
// an unanswered card to a reject): nothing written, the old values stand.
//
// EVERY GOROUTINE GOES THROUGH observe.Registry.Spawn (D22 ban #1), parked on
// this run's own root so close() cancels it and the join at close waits for
// it - an abandoned card resolves as a deny, so shutdown can never be the
// moment a loosening slips in (the same posture config_reload.go's tick keeps).

import (
	"context"
	"fmt"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The card identity for a queued config loosening. Neither is a C4
// registration and neither is a C17 method - the same shape configReloadTaskID
// / modeSwitchTaskID keep.
const (
	configLoosenTaskID   = "host:config-loosen"
	configLoosenToolName = "config.loosen"
)

// startConfigLoosenGate wires Manager.LoosenQueue to this run's gate. Called
// from the assembly root AFTER the answer side exists (a card raised before
// the reply listener is up is a card nobody can answer - the same ordering
// rule startConfigReload keeps), and only when there is a gate at all: a host
// with no card surface stays on the fail-closed refusal, which is the honest
// state, not a bug to route around.
func (rt *agentRuntime) startConfigLoosenGate() {
	if rt == nil || rt.mgr == nil || rt.gate == nil {
		return
	}
	rt.loosenRoot = observe.NewRoot("config-loosen")
	rt.mgr.LoosenQueue = configLoosenQueue{rt: rt}
	rt.auditf("config: LOOSEN-GATE state=armed gate=on detail=%q",
		"受保护写路（WriteSection）遇到 [risk]/[fs]/[net]/[plugins] 的放宽：先弹一张原生 L2 卡，"+
			"接受才落盘并记一条审计（老值→新值+卡id），拒绝或超时保持旧值；面板来源的放宽仍然结构性拒绝")
}

// configLoosenQueue is the production config.LoosenQueue: hand the card out,
// return at once, resolve on the answer.
type configLoosenQueue struct{ rt *agentRuntime }

// QueueLoosen implements config.LoosenQueue. It never waits for a human
// (ruling N9): the goroutine that parks on the card is spawned through
// observe.Registry.Spawn and returns through channels, not through this call.
// A run with no gate refuses the hand-off - the Manager then answers
// ErrLoosenQueueUnwired upstream, the same fail-closed direction a nil queue
// carries.
func (q configLoosenQueue) QueueLoosen(req config.LoosenRequest) (string, error) {
	if q.rt == nil || q.rt.gate == nil {
		return "", fmt.Errorf("config: no approval gate is assembled in this run; the %v loosening card cannot be raised",
			req.Keys)
	}
	if req.RequestID == "" {
		return "", fmt.Errorf("config: the loosening request carries no RequestID; the card would have nothing to resolve under")
	}
	if q.rt.loosenRoot == nil {
		return "", fmt.Errorf("config: the loosen-gate root is not armed in this run; refusing to raise the %v card",
			req.Keys)
	}
	// "approval-waiter" is the D38b roster's own name for a goroutine that
	// parks waiting on a card answer (observe/goroutine.go's PerTaskPrefixes);
	// a made-up fifth name would be a leak symptom by the registry's own
	// classification, and this goroutine IS an approval waiter - it lives and
	// dies with one card.
	observe.Default.Spawn("approval-waiter", req.RequestID, q.rt.loosenRoot, func(ctx context.Context) {
		q.rt.resolveLoosenCard(ctx, req)
	})
	return req.RequestID, nil
}

// resolveLoosenCard raises ONE native L2 card for one queued loosening and
// hands the verdict to Manager.ResolveQueuedLoosen. Every exit except an
// explicit allow resolves as a refusal that keeps the old values - the same
// direction D33 requires of the reload path's card.
func (rt *agentRuntime) resolveLoosenCard(ctx context.Context, req config.LoosenRequest) {
	if rt == nil || rt.gate == nil {
		// No gate to ask: the same fail-closed answer a nil queue carries.
		_, _ = rt.mgr.ResolveQueuedLoosen(req.RequestID, false)
		return
	}
	revoke := rt.gate.AdmitTextTask(configLoosenTaskID)
	defer revoke()
	ans, why := rt.gate.PendingApproval(ctx, tools.Decision{
		TaskID:        configLoosenTaskID,
		CorrelationID: req.RequestID,
		Tool:          configLoosenToolName,
		Level:         risk.L2,
		Reason: fmt.Sprintf("这次配置写入会放宽 [%s]：%v（来源 %s，origin %s）。"+
			"D33/D36 要求放宽必须先重新确认一次才生效；拒绝或不做声都保留旧的严格值。"+
			"（这张卡由配置写路签发，与任何一次工具调用无关）",
			req.Section, req.Keys, req.Source, req.Origin),
		Mode: rt.permissionMode(),
	})
	if ans != tools.AnswerAllow {
		if why == "" {
			why = string(ans)
		}
		rt.auditf("config: LOOSEN-GATE state=answered card=%s section=%s keys=%v answer=%s result=deny detail=%q",
			req.RequestID, req.Section, req.Keys, ans, why)
		// The refused/timeout arm: ResolveQueuedLoosen writes nothing and
		// books the denied line; the old values are what any read-back finds.
		_, _ = rt.mgr.ResolveQueuedLoosen(req.RequestID, false)
		return
	}
	rt.auditf("config: LOOSEN-GATE state=answered card=%s section=%s keys=%v answer=%s result=allow",
		req.RequestID, req.Section, req.Keys, ans)
	if _, err := rt.mgr.ResolveQueuedLoosen(req.RequestID, true); err != nil {
		// The accepted write refused itself upstream (unreadable file, drift,
		// validation) - it booked its own audit line; nothing to add here
		// beyond naming that the card WAS answered, so the join is readable.
		rt.auditf("config: LOOSEN-GATE state=resolve-failed card=%s err=%v", req.RequestID, err)
	}
}

// compiled-only guard: the adapter must keep satisfying the seam it exists for.
var _ config.LoosenQueue = configLoosenQueue{}
