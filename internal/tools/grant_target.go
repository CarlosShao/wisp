package tools

import (
	"context"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// Ticket 49: the input.type target-process binding (D45-2's §15 item-7 variant)
//
// The file-path grants ride GrantSource.Covering (grant.go, ticket 224). This
// file is the SECOND dialect of the same (tool, pattern, session) row: for
// input.type the pattern is the RESOLVED TARGET PROCESS NAME (§15 ruling item
// 7: 「按目标进程分会话授权」), and the downgrade it buys is the one the frozen
// D34 row names - input.type is L2 (首次) → 可按 D45 会话授权降为 L1
// (PLAN.md:2557). Nothing here widens anything else: an L2 verdict reached
// WITHOUT this variant (another L2 tool's R1 floor, R8, R4, R5...) is never
// covered, because the downgrade is keyed on exactly one tool and refuses on
// every red line.
//
// WHY A SECOND SEAM AND NOT OVERLOADING Covering: the read seam
// (tools.GrantSource) is pinned read-only and one-method by
// TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority. Adding a method to
// that interface re-shapes the frozen seam; a SEPARATE seam keeps the file-path
// source untouched, is satisfied by the same *session.Ledger, and lets the
// bridge fail closed when its host cannot answer process questions at all (nil
// TargetGrants = always ask, exactly the posture grant.go documents for a nil
// Grants).
//
// WHY THE RESOLVER IS INJECTED TOO: the process name the card must name, and
// the one a grant binds, is what the OS backend resolves from the call's own
// target (hwnd / title / foreground → pid → image name). Judging happens
// BEFORE execution, so the bridge asks an injected resolver instead of the
// backend - the same resolution the backend will repeat at Execute time, and
// the execution half re-checks it (the swap guard below). nil resolver = the
// variant is dormant and every input.type call asks, which is the fail-closed
// posture every pre-49 assembly keeps.
// ---------------------------------------------------------------------------

// TargetGrantSource is the read side of the input.type binding rows. It is
// satisfied by internal/session's Ledger (TargetCovering) and by nothing in
// this package. One method, read-only, fail-closed by contract - the same
// three obligations Covering states, keyed on the process name instead of
// paths.
//
// It is NOT a second field on Bridge on purpose: the ticket-224 pin
// (TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter) freezes the bridge at
// exactly ONE grant-shaped field, and the reason survives ticket 49 unchanged -
// two channels would be two places an authorization source can come from, with
// only one in the audit log. The input.type read goes through the SAME
// injected source, discovered by an interface assertion at call time
// (sessionTargetID below): a ledger that implements TargetCovering (every
// post-49 *session.Ledger) answers it, and anything else answers "not
// covered", which keeps the question asking - the fail-closed posture.
type TargetGrantSource interface {
	// TargetCovering reports whether the host session already holds a live
	// grant for (tool, target process), and returns the id of the row that
	// covers it (tool_call.grant_id, SPEC-02 §3). It must answer false, not
	// guess true, whenever it cannot look.
	TargetCovering(ctx context.Context, tool, process string) (grantID int64, ok bool)
}

// targetCoveringReader is the assertion face sessionTargetID reads. Declared
// unexported so the single public seam story stays simple: the bridge holds
// one GrantSource; whether that source ALSO answers the target question is
// discovered here, at call time, from the same object.
type targetCoveringReader interface {
	TargetCovering(ctx context.Context, tool, process string) (int64, bool)
}

// TargetResolver resolves the injection target one input.type call names, from
// the call's own parameters. The composition root wires it to the SAME
// SystemBridge the tools execute through, so the name a card shows, the name a
// grant binds and the name the backend injects into are all one resolver's
// output.
type TargetResolver func(params map[string]any) (process string, ok bool)

// TargetGrantTool is the ONE tool the binding applies to. The §15 ruling scopes
// the variant to input.type and nothing else; a second tool would be a contract
// move, not a line in this file.
const TargetGrantTool = "input.type"

// inputTypePidFallbackPrefix marks the name a backend answers when the image
// name could not be read. A pid-keyed rule is dead the moment the pid is
// reused by another program, and a dead-looking-live row is exactly the shape
// the session ledger refuses to store - so a fallback name is refused HERE,
// where its provenance is known, and the call keeps asking.
const inputTypePidFallbackPrefix = "pid:"

// isGrantableTarget reports whether a resolved process name may key a grant
// lookup at all.
func isGrantableTarget(name string) bool {
	if name == "" {
		return false
	}
	if strings.HasPrefix(name, inputTypePidFallbackPrefix) {
		return false
	}
	return true
}

// sessionTargetID asks the seam about one input.type call, recovering from
// anything it throws - the same fail-closed shape sessionGrantID has (grant.go:
// 57-99). The read goes through the ONE injected grant source (b.grants),
// discovering the target capability with an interface assertion; a source
// without it (or a nil one) answers 0: ask. That keeps the bridge at exactly
// one grant-shaped field (ticket 224's pin) while the same ledger serves both
// dialects.
func (b *Bridge) sessionTargetID(ctx context.Context, process string) (id int64) {
	if b.grants == nil {
		return 0
	}
	tr, ok := b.grants.(targetCoveringReader)
	if !ok {
		return 0
	}
	if !isGrantableTarget(process) {
		b.log("tools: GRANT-TARGET-SKIP tool=%s target=%q (pid 回退名不落授权，按未授权处理)", TargetGrantTool, process)
		return 0
	}
	defer func() {
		if rec := recover(); rec != nil {
			b.log("tools: GRANT-READ-FAILED recovering=%v; 按未授权处理（继续询问）", rec)
			id = 0
		}
	}()
	id, ok = tr.TargetCovering(ctx, TargetGrantTool, process)
	if !ok {
		return 0
	}
	return id
}

// targetBindingBlocked reports why one L2 verdict may NOT ride the input.type
// downgrade. The downgrade exists for exactly one shape: R1's own declared
// floor (the D34 row's L2 for input.type). Any red line in RulesHit or the R4
// block flag refuses it - SPEC-06 §8.3's three bullets, with the §15 item-7 ⑤
// sentence (「C25 污染命中（R4）不受本授权覆盖」) as the loudest.
func targetBindingBlocked(verdict risk.Decision) (string, bool) {
	if verdict.SessionOverrideBlocked {
		return "R4 污染升级不受任何会话授权覆盖（SPEC-06 §8.3 第 2 条）", true
	}
	for _, r := range verdict.RulesHit {
		switch r {
		case risk.R2:
			return "R2 越界判定不适用会话授权", true
		case risk.R3:
			return "R3 敏感路径判定不适用会话授权", true
		case risk.R4:
			return "R4 污染升级不受会话授权覆盖（§15 第 7 项 ⑤）", true
		case risk.R5:
			return "R5 网络目标判定不适用会话授权", true
		case risk.R6:
			return "R6 shell 判定不适用会话授权", true
		case risk.R7:
			return "R7 批量规模判定不适用会话授权", true
		case risk.R8:
			return "R8 不可逆判定不适用会话授权", true
		case risk.R9:
			return "R9 fail-closed 判定不适用会话授权", true
		}
	}
	return "", false
}
