package approval

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// The injected presentation surface (segment 2 implements this)
// ---------------------------------------------------------------------------

// Prompt is one confirmation as the gate hands it to the UI. The UI is an
// INJECTED interface: this segment ships no native card and no ball visuals,
// and the tests fake this type.
//
// Two rules are structural here, not conventional:
//
//   - Params carries the FULL argument object (C27 fallback-grade: "tool name
//   - L2 badge + full params, not a summary"), and RulesHit/Reason travel as
//     data because ticket 17 froze them as verbatim card content.
//   - Grant is populated ONLY on the native surface. A prompt the panel can
//     reach is rendered from PanelItem, which has no such field.
type Prompt struct {
	CorrelationID string
	TaskID        string
	Tool          string
	Level         string // 'L1' | 'L2'
	Params        map[string]any
	Args          json.RawMessage
	RulesHit      []risk.RuleID
	Reason        string
	Paths         []string
	Batch         *BatchView // nil unless D45-1 aggregation applies
	Window        time.Duration
	Deadline      time.Duration
	Channels      []ChannelStatus
	Depth         int
	Grant         string
}

// PanelItem is everything the panel (ticket 37) is allowed to see about one
// pending approval. There is deliberately no grant field, no Allow capability
// and no way to name a source: SPEC-06 §9 layer 3 says the panel may only
// 「拒绝 / 查看完整参数」, and a struct that cannot express "allow" is a
// stronger guarantee than a check that hopes the caller meant it.
type PanelItem struct {
	CorrelationID string
	Tool          string
	Level         string
	Reason        string
	Paths         []string
	Depth         int
}

// EventKind names the transient states the strip/card must show.
type EventKind string

// The events.
const (
	EventTick      EventKind = "tick"      // L1 countdown still running
	EventWarning   EventKind = "warning"   // C18 last-30s prominent warning
	EventDismissed EventKind = "dismissed" // prompt gone (answered/expired)
	EventStarted   EventKind = "started"   // handoff to execution happened
)

// Event is one transient update about a live prompt.
type Event struct {
	Kind          EventKind
	CorrelationID string
	Remaining     time.Duration
	Text          string
}

// UI is the injected presentation seam. Prompt is called once per displayed
// confirmation, Update for each transient state.
//
// A nil UI, or a UI that returns an error, is treated as "the host could not
// be reached": the gate fails closed and the operation does not run. There is
// no path through this package where an unanswered question executes.
type UI interface {
	Prompt(ctx context.Context, p Prompt) error
	Update(ctx context.Context, e Event) error
}

// UIFuncs adapts two callbacks to UI for hosts and tests that only have
// functions. A nil callback fails closed like a nil UI.
type UIFuncs struct {
	PromptFn func(context.Context, Prompt) error
	UpdateFn func(context.Context, Event) error
}

// Prompt implements UI.
func (u UIFuncs) Prompt(ctx context.Context, p Prompt) error {
	if u.PromptFn == nil {
		return errNoUI
	}
	return u.PromptFn(ctx, p)
}

// Update implements UI. Update failures are not fatal: a strip that cannot
// repaint must not retroactively cancel a question the user can still answer.
func (u UIFuncs) Update(ctx context.Context, e Event) error {
	if u.UpdateFn == nil {
		return nil
	}
	return u.UpdateFn(ctx, e)
}

var errNoUI = errors.New("approval: 确认界面未接入，已 fail-closed 拒绝执行")

// ---------------------------------------------------------------------------
// The two API faces (F2 layer 3)
// ---------------------------------------------------------------------------

// Routing errors. They are values, not strings, because the caller's branch
// differs: an unknown correlation id is a mis-routed reply, a bad grant is a
// forgery attempt, a spent grant is a replay.
var (
	// ErrUnknownCorrelation: no pending item carries that id.
	ErrUnknownCorrelation = errors.New("approval: correlation_id 无对应待审批项")
	// ErrNotPending: the item exists but already left the queue.
	ErrNotPending = errors.New("approval: 该审批项已结束，不再接受决定")
	// ErrBadGrant: the native proof is missing, spent, or bound to another
	// item. The message never says which, so the API cannot be used to probe
	// for valid nonces.
	ErrBadGrant = errors.New("approval: 原生令牌无效（缺失/已用/与本次请求不绑定）")
	// ErrPanelAllow: an allow arrived on the panel surface. Rejected on the
	// route alone, before any grant is looked at.
	ErrPanelAllow = errors.New("approval: 面板来源不得允许（F2 第三层：允许只接受原生侧）")
	// ErrNotAdmitted: the task never entered the approval layer through the
	// TEXT loop, which D47 says is the only path that can have tools.
	ErrNotAdmitted = errors.New("approval: 该任务未经文本循环登记，不得进入审批队列（D47）")
)

// NativeAPI is the surface the native side holds: the ball click handler, the
// native card's buttons and the global hotkey (SPEC-06 §2). It is the ONLY
// place an allow can arrive, and an allow must present the grant that
// UI.Prompt was handed.
type NativeAPI interface {
	// Allow spends grant for correlationID. grant is not optional and is not
	// bypassable: an empty or stale value returns ErrBadGrant.
	Allow(ctx context.Context, correlationID, grant string) error
	// Reject needs no proof: refusing is always safe, so the panel may do it
	// too through PanelAPI.
	Reject(correlationID, reason string) error
}

// PanelAPI is the ENTIRE surface a panel/WebView host may be handed. It has no
// Allow method: "panel-sourced allow is structurally rejected" is then not a
// check that could be forgotten, it is a method that does not exist.
type PanelAPI interface {
	Reject(correlationID, reason string) error
	Head() (PanelItem, bool)
	View(correlationID string) (PanelItem, bool)
}

// ---------------------------------------------------------------------------
// D45-1 batch view
// ---------------------------------------------------------------------------

// BatchView is the aggregated summary of a homogeneous L1 batch: total +
// affected dirs + first 5 (D45-1). It is data, so the card can render the
// expandable list without re-deriving what the gate already decided.
type BatchView struct {
	Total        int
	AffectedDirs []string // unique, sorted
	First        []string // at most BatchPreviewCount, in call order
	Deferred     int      // items beyond the preview, still approved by this one confirm
	Summary      string
}

// BatchPreviewCount is D45-1's "前 5 条明细".
const BatchPreviewCount = 5
