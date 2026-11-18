package main

// Readface wave 4 (2026-11-18), ticket 40 AC#2's Go half: the real leg behind
// internal/panel's session.rules.list / session.rules.revoke doors - the
// panel-side live view of the D45-2 session ledger and its one-click revoke.
//
// WHY THIS DOOR EXISTS. Ticket 49's grants rows have had a console exit since
// the verbs landed (approval_reply.go grantsList / revokeGrant, teethed in
// ticket49_verbs_test.go), but the security page the specs promise (SPEC-06
// §8.2 「面板安全页实时列出并可一键撤销」, SPEC-00 #34) had nothing to read:
// the snapshot carries no grants section and the spec's own route spellings
// (grants.list / grants.revoke, SPEC-08 §5.2:169) die on the frozen
// grant-boundary instrument (bridge.go's note on the two wave-4 constants
// records the naming decision). AC#2 needs the view to reflect 49's grants
// LIVE and the revoke to bite on the very NEXT covered call.
//
// THE SEAM IS THE LEDGER, NOT A COPY. The door holds no grant state of its
// own - the seam hands it the run's own session ledger (agentRuntime.session,
// the internal/session.Ledger the answer route writes through). Three doors,
// all the ledger's own:
//
//   - ListBySession  (grants.go:375): the read face ticket 49's console verb
//     walks - the same rows, the same order (newest first).
//   - Revoke         (grants.go:338): the ONE revoke entrance the console
//     verb spends. It stamps revoked_at on disk; rows are retained for the
//     audit window (SPEC-02 §4), and the NEXT Covering/TargetCovering call
//     answers false because grantLive reads revoked_at - which is exactly
//     why "撤销立刻生效" holds without any invalidation code: the gate asks
//     the store again, and the store now says no.
//   - SessionID      : the identity the listing names, so a zero-row answer
//     says WHICH session has none.
//
// THE SEAM IS LATE-BOUND, the roster-cancel pattern: the panel chain is
// assembled before the task pipeline exists (resident_windows.go's ordering),
// so the door leaves its seam nil at construction and the boot binds it on
// the first line where the runtime exists (bindPanelSessionRules). A nil seam
// is a legal state (the CLI seam) and refuses by name, never drops.
//
// THE PAYLOAD IS THE ROWS, NOT A SUMMARY. Every row is the store's own
// approval_grant fields verbatim - id/scope/tool/pattern/session/createdAt/
// expiresAt - exactly the field set the dispatch names. Two display-only
// judgements ride the door, both stated by ticket 49's operator view already:
// the LIVE filter (grantLive: not revoked, not expired - the door shows what
// is in force, and a revoked row LEAVES this list while its audited row stays
// on disk) and the zero-row sentence (an empty ledger is said OUT LOUD in the
// payload's zero field, so the page renders an explicit 还没有授权 and never
// an empty table dressed up as a reading).
//
// WHY NO CONFIRM BIT ON THE REVOKE. Revoking REMOVES an authorization: the
// worst outcome is the gate asking again - the safe direction, and the design
// (SPEC-06 §8.2 一键撤销, 纠错成本必须≈0). There is no panel-side allow
// anywhere on this door (AGENTS.md §1.2): the door can only make the gate
// ask more, never less.
//
// D42#9 NOTE (the clock). The live filter compares persisted absolute
// timestamps (expires_at) with the current wall clock - no duration
// difference is measured, the same shape approval_reply.go's grantsList
// already documents. Persisted cross-restart absolute moments have no
// monotonic representation; wall clock is what the timestamps are.
//
// SHAPE. The settings-route shape, like every read face in this family: one
// prose sentence plus one payload line behind SessionRulesPayloadMarker, the
// page finds the marker and parses the rest of the line as JSON. The prose is
// never parsed by the page (the same rule cost/privacy/history state).

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// SessionRulesPayloadMarker starts the session-rules read payload. Same
// naming properties as every marker in the family: no `panel.` prefix and no
// dot - it is neither a method name nor a route.
const SessionRulesPayloadMarker = "wisp-session-rules:"

// sessionRulesRow40 is ONE ledger row as the wire carries it: the store's own
// fields, verbatim (memory.ApprovalGrant), minus the revoked_at stamp - this
// face lists LIVE rows, and revoked_at non-nil is precisely what the live
// filter removes. A revoked row keeps its audited row on disk; the listing
// simply stops carrying it (tooth ③'s red is a read face that still lists a
// spent row).
type sessionRulesRow40 struct {
	ID        int64  `json:"id"`
	Scope     string `json:"scope"`
	Tool      string `json:"tool"`
	Pattern   string `json:"pattern"`
	SessionID string `json:"session"`
	CreatedAt int64  `json:"createdAt"`
	ExpiresAt int64  `json:"expiresAt"`
}

// sessionRulesPayload40 is the read door's wire shape, separate from the view
// type so a field added to the view does not silently change the page's JSON
// contract (the costPayload rule). zero is the EXPLICIT empty state: an empty
// string on a non-empty ledger, the door's own sentence on an empty one - the
// page renders it verbatim and never invents an empty-table reading.
type sessionRulesPayload40 struct {
	Session string              `json:"session"`
	Rows    []sessionRulesRow40 `json:"rows"`
	Zero    string              `json:"zero"`
}

// sessionRulesSeams is the one thing the door needs from the runtime: the
// ledger's own three doors plus the clock. It is a struct of functions (not
// an interface) because the ledger has exactly these doors and the 453/452
// precedents set the function shape; the boot binds ONE struct.
type sessionRulesSeams struct {
	session func() string
	list    func(ctx context.Context) ([]memory.ApprovalGrant, error)
	revoke  func(ctx context.Context, id int64) error
	now     func() time.Time
}

// sessionRulesDoor40 is the wave-4 socket of the panel chain.
type sessionRulesDoor40 struct {
	audit panel.AuditFunc
	actor string

	mu   sync.RWMutex
	seam sessionRulesSeams
}

func newSessionRulesDoor40(audit panel.AuditFunc, actor string) *sessionRulesDoor40 {
	return &sessionRulesDoor40{audit: audit, actor: actor}
}

// bindSeam is the late bind. nil struct fields refuse by name, per door.
func (d *sessionRulesDoor40) bindSeam(seam sessionRulesSeams) {
	if d == nil {
		return
	}
	d.mu.Lock()
	d.seam = seam
	d.mu.Unlock()
}

func (d *sessionRulesDoor40) seamOf() sessionRulesSeams {
	if d == nil {
		return sessionRulesSeams{}
	}
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.seam
}

// sessionRulesSeam40 is the production seam, on the composition root: the
// closures read THIS runtime's session ledger - the same object the answer
// route writes through, so what the page lists and what the gate honours
// cannot diverge. A runtime with no ledger (the CLI seam) answers with empty
// closures the door refuses by name.
func (rt *agentRuntime) sessionRulesSeam40() sessionRulesSeams {
	return sessionRulesSeams{
		session: func() string {
			if rt == nil || rt.session == nil {
				return ""
			}
			return rt.session.SessionID().String()
		},
		list: func(ctx context.Context) ([]memory.ApprovalGrant, error) {
			if rt == nil || rt.session == nil {
				return nil, fmt.Errorf("这一装配没有接入会话授权记账（rt.session=nil）")
			}
			return rt.session.ListBySession(ctx)
		},
		revoke: func(ctx context.Context, id int64) error {
			if rt == nil || rt.session == nil {
				return fmt.Errorf("这一装配没有接入会话授权记账（rt.session=nil），无从撤销")
			}
			return rt.session.Revoke(ctx, id)
		},
		now: time.Now,
	}
}

// liveRows is the LIVE filter: not revoked, not expired (grantLive's two
// conditions, read off the store's own rows). The unfiltered read stays
// available one call away (the seam's list) so the audit view never loses its
// source; this face shows what is in force.
func liveRows40(rows []memory.ApprovalGrant, now int64) []memory.ApprovalGrant {
	live := make([]memory.ApprovalGrant, 0, len(rows))
	for _, g := range rows {
		if g.RevokedAt != nil {
			continue
		}
		if g.ExpiresAt <= now {
			continue
		}
		live = append(live, g)
	}
	return live
}

func toWireRows40(rows []memory.ApprovalGrant) []sessionRulesRow40 {
	out := make([]sessionRulesRow40, 0, len(rows))
	for _, g := range rows {
		out = append(out, sessionRulesRow40{
			ID:        g.ID,
			Scope:     g.Scope,
			Tool:      g.Tool,
			Pattern:   g.Pattern,
			SessionID: g.SessionID,
			CreatedAt: g.CreatedAt,
			ExpiresAt: g.ExpiresAt,
		})
	}
	return out
}

// HandleSessionRulesList answers session.rules.list: THIS session's live
// rows, verbatim, plus the explicit zero state. Outcomes, in order:
//
//	nil seams           -> refused by name (no ledger in this process)
//	list read failed    -> the ledger's own refusal, verbatim
//	otherwise           -> receipt + payload line (rows, or the zero sentence)
func (d *sessionRulesDoor40) HandleSessionRulesList(ctx context.Context, req panel.ComposerRequest) (string, error) {
	seam := d.seamOf()
	if seam.list == nil {
		err := fmt.Errorf("%w: 会话授权读面未接入（本进程没有会话授权记账，或任务入口没有起来）：本会话的授权无从列出，也不会被静默丢掉",
			panel.ErrNoHandlerAttached)
		d.record(req, "list", "refused", err)
		return "", err
	}
	rows, err := seam.list(ctx)
	if err != nil {
		err := fmt.Errorf("panel: 会话授权读取失败（requestId=%q）：%w", req.RequestID, err)
		d.record(req, "list", "read-failed", err)
		return "", err
	}
	sid := ""
	if seam.session != nil {
		sid = seam.session()
	}
	now := int64(0)
	if seam.now != nil {
		now = seam.now().Unix()
	}
	live := liveRows40(rows, now)
	if live == nil {
		live = []memory.ApprovalGrant{}
	}
	payload := sessionRulesPayload40{Session: sid, Rows: toWireRows40(live)}
	if len(live) == 0 {
		// The explicit zero state, the door's own words - the same sentence
		// ticket 49's console verb says, so the two exits cannot drift into
		// two truths about an empty ledger.
		payload.Zero = fmt.Sprintf("本会话（%s）当前没有生效中的授权。", sid)
	}
	b, merr := json.Marshal(payload)
	if merr != nil {
		err := fmt.Errorf("panel: 会话授权载荷编不出来（requestId=%q）：%w", req.RequestID, merr)
		d.record(req, "list", "marshal-failed", err)
		return "", err
	}
	d.record(req, "list", "answered", nil)
	var sentence string
	if len(live) == 0 {
		sentence = payload.Zero
	} else {
		sentence = fmt.Sprintf("本会话（%s）生效中的授权共 %d 行；行内字段逐字在下面这一行 JSON（真行，不是摘要）。撤销一枚：授权页每行的撤销键，或 console 的 revoke-grant <编号>。", sid, len(live))
	}
	return sentence + "\n" + SessionRulesPayloadMarker + string(b), nil
}

// HandleSessionRulesRevoke answers session.rules.revoke: ONE row spent
// through the ledger's own Revoke entrance, the receipt naming the row.
// Outcomes, in order:
//
//	empty selector      -> refused before the seam is reached
//	non-numeric / <=0   -> refused before the seam is reached
//	nil seams           -> refused by name (no ledger in this process)
//	ledger refused      -> the ledger's own refusal, verbatim
//	otherwise           -> receipt naming the row
func (d *sessionRulesDoor40) HandleSessionRulesRevoke(ctx context.Context, req panel.ComposerRequest) (string, error) {
	raw := strings.TrimSpace(req.SessionRuleID)
	if raw == "" {
		err := fmt.Errorf("%w: sessionRuleId 是空的：revoke 只撤一行，不给它编号就什么都不撤（编号来自授权页读面的 id 列）",
			panel.ErrComposerRequest)
		d.record(req, "revoke", "refused", err)
		return "", err
	}
	id, perr := strconv.ParseInt(raw, 10, 64)
	if perr != nil || id <= 0 {
		err := fmt.Errorf("%w: sessionRuleId %q 不是读面列出的授权编号（要一个正整数，收到 %q）",
			panel.ErrComposerRequest, raw, raw)
		d.record(req, "revoke", "refused", err)
		return "", err
	}
	seam := d.seamOf()
	if seam.revoke == nil {
		err := fmt.Errorf("%w: 会话授权撤销通路未接入（本进程没有会话授权记账，或任务入口没有起来）：授权 #%d 未被触碰，也不会被静默丢掉",
			panel.ErrNoHandlerAttached, id)
		d.record(req, "revoke", "refused", err)
		return "", err
	}
	if err := seam.revoke(ctx, id); err != nil {
		err := fmt.Errorf("panel: 撤销会话授权 #%d 被拒绝（requestId=%q）：%w", id, req.RequestID, err)
		d.record(req, "revoke", "revoke-failed", err)
		return "", err
	}
	d.record(req, "revoke", "revoked", nil)
	return fmt.Sprintf("已撤销会话授权 #%d：本会话内对应工具的下一次调用会重新弹卡确认（行保留在审计里，授权页的生效清单不再列出它）。", id), nil
}

func (d *sessionRulesDoor40) record(req panel.ComposerRequest, verb, outcome string, err error) {
	if d == nil || d.audit == nil {
		return
	}
	line := "ok"
	if err != nil {
		line = err.Error()
	}
	d.audit("panel: SESSION-RULES request=%q method=%q verb=%q selector=%q actor=%q outcome=%q err=%s",
		req.RequestID, req.Method, verb, strings.TrimSpace(req.SessionRuleID), d.actor, outcome, line)
}
