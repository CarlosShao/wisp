package approval

// The host-facing reply seam (ticket 201 AC#1's last Go-side inch).
//
// WHY THIS FILE EXISTS, precisely. Before it, the ANSWER side of an approval had
// exactly one shape in this tree: two exported routers on the Gate
// (DecideFromNative / DecideFromPanel) plus the two faces (Native() / Panel()),
// and one caller - the console listener in cmd/wisp. That caller had to carry its
// own ledger of displayed cards, because the single-use grant an allow needs is
// handed to the injected UI (ui.go:99-111, Prompt.Grant) and nothing downstream
// could read it back. The ledger therefore lived in `package main`, where no
// other package can compile against it: a native surface that is not a terminal
// (the ball's buttons, the tray menu, a native card) had nothing to be handed.
// That is the difference between "an entry point exists" and "an entry point a
// host can call", and ticket 201 AC#1 is about the second one.
//
// WHAT THIS IS NOT. It opens no third authority and no third decision route -
// the two routers above are still the only doors, and this type only ever calls
// them. It cannot mint a grant (only Queue/Gate does, per card, exactly once), it
// cannot bypass one (an empty or spent grant returns ErrBadGrant like the router
// it delegates to), and it cannot carry an allow on the panel route
// (PanelAllow hands the request to DecideFromPanel and reports what came back,
// which is ErrPanelAllow by route alone - AGENTS.md §1.2 ban #6, F2 layer 3).
// The command surface is a host's click, not a wire method name: nothing here
// joins C17's inbound roster, and the panel roster stays untouched.
//
// ONE COPY. cmd/wisp's console listener answers through this type rather than
// alongside it. The reason is stated in that file's own words: a second copy of
// the grant-spending and route-choosing logic is how one of them grows a
// leniency the other does not have.

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// MaxTrackedCards bounds the host-side ledger. The queue bounds its own pending
// items at DefaultMaxPending, and an L1 window leaves this ledger on its
// dismissal or handoff event, so the ceiling only covers the case where an event
// never arrives (a UI that returned an error mid-flight). Eviction is the safe
// direction: an evicted card has no grant held for it, so it can be refused but
// never allowed through this seam.
const MaxTrackedCards = 16

// Routing errors of this seam - the two cases the routers cannot see, because a
// router only knows what arrives on it, not what this host ever displayed.
var (
	// ErrNoTrackedCard: this surface never displayed that correlation (or it
	// already left the screen), so there is no grant to spend and no card to
	// address.
	ErrNoTrackedCard = errors.New("approval: 本机账上没有这张卡（从未显示、已答复，或已经离开屏幕）")
	// ErrRouteHasNoAllow: the card is an L1 pre-execution window, which has a
	// veto and no allow (SPEC-06 §2 B1).
	ErrRouteHasNoAllow = errors.New("approval: 这一路线只有否决、没有允许")
	// ErrNoGateAttached: the ledger was built but never bound to a Gate, so no
	// router exists to carry a decision. It fails closed with a refusal, never
	// with a silent yes.
	ErrNoGateAttached = errors.New("approval: 答复面未绑定审批 gate，任何决定都送不出去")
)

// ReplyCard is one confirmation this process displayed and has not yet
// retracted, as an answering surface needs it: where it is, what it is, which
// paths it named, and - native side only - the proof an allow has to present.
type ReplyCard struct {
	CorrelationID string
	TaskID        string
	Tool          string
	Level         string // "L1" | "L2", verbatim off the Prompt
	// Grant is the single-use native nonce Gate.PendingApproval minted for THIS
	// card and handed to exactly one recipient, the injected UI. An L1 window's
	// entry carries "" because that route has no allow verb to answer for.
	Grant string
	// Paths is Prompt.Paths verbatim: what the card itself named. The 长期 branch
	// reads it to state which rule a widening would store (ticket 201 AC#4), so
	// the rule text comes off the card rather than off a re-derivation.
	Paths []string
	// Window is the L1 block length, zero on an L2 card. Deadline is the C18
	// ceiling an unanswered L2 card dies at.
	Window   time.Duration
	Deadline time.Duration
}

// Replies is the entry point a host surface is handed.
//
// The method set is the whole boundary, and it is deliberately asymmetric:
//
//	Allow / AllowSession / Reject / Veto     the native route (DecideFromNative,
//	                           Gate.Veto); Allow and AllowSession both spend the
//	                           card's single-use nonce, AllowSession additionally
//	                           records a D45 session rule (ticket 224)
//	PanelReject / Head / View        the panel route (DecideFromPanel, PanelAPI)
//	PanelAllow                       reaches DecideFromPanel and is refused THERE
//	Pending / AwaitingHuman / WaitingState   read-only, no authority at all
//
// There is no method that lets a panel-sourced decision execute anything, and no
// path that mints a grant or forwards one to the panel route. That is not a check
// this type remembers to run; it is the absence of a door.
type Replies struct {
	mu    sync.Mutex
	byID  map[string]ReplyCard
	order []string

	// gate is bound after construction: the ledger is handed to the injected UI
	// (which needs it at Prompt time) while the Gate that owns it is built one
	// statement later, and the UI is one of that Gate's own options. A nil gate
	// makes every route return ErrNoGateAttached, so a half-wired host cannot
	// answer anything.
	gate        *Gate
	vetoChannel Channel
	// nativeSource / panelSource are the transport names this host wants booked
	// in the audit line. They are LABELS, not authority (gate.go says no branch
	// reads Request.Source), and they exist so the ledger can hold "what the
	// transport claimed" next to "what the route proved". Empty values fall back
	// to the seam's own names, so a host that declares nothing still gets an
	// honest line rather than a blank one.
	nativeSource string
	panelSource  string
}

// HostBinding is what a host hands the ledger once the Gate exists.
type HostBinding struct {
	Gate *Gate
	// VetoChannel is which of SPEC-06 §2's four channels this host's cancel
	// transport really is. Empty = none, which Gate.Veto then says out loud.
	VetoChannel Channel
	// NativeSource / PanelSource name this host's two transports in the audit.
	NativeSource string
	PanelSource  string
}

// Host labels the seam stamps when a host declares none.
const (
	defaultNativeSource = "approval-host-native"
	defaultPanelSource  = "approval-host-panel-route"
)

// NewReplies makes an unbound ledger. It is safe to Record into before Attach:
// displaying a card and answering it are two different moments of a run, and a
// host must not lose the first because the second is not wired yet.
func NewReplies() *Replies {
	return &Replies{byID: map[string]ReplyCard{}}
}

// Attach binds this ledger to the Gate whose UI it was injected into, and
// declares which of SPEC-06 §2's four veto channels this host's cancel transport
// really is. An empty channel is the honest value for a host that wired none of
// them: Gate.Veto's own ChannelRegistry then answers 「未知取消通道」, which is a
// refusal the user hears rather than a channel name this file invented (B1).
func (r *Replies) Attach(b HostBinding) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gate = b.Gate
	r.vetoChannel = b.VetoChannel
	r.nativeSource = b.NativeSource
	r.panelSource = b.PanelSource
}

// Record books one displayed card. The host's UI face calls it from Prompt with
// the Prompt it was handed, and nothing here shapes a field: a card the seam
// cannot describe is a card the seam never showed.
func (r *Replies) Record(p Prompt) {
	if r == nil || p.CorrelationID == "" {
		return
	}
	card := ReplyCard{
		CorrelationID: p.CorrelationID,
		TaskID:        p.TaskID,
		Tool:          p.Tool,
		Level:         p.Level,
		Grant:         p.Grant,
		Paths:         append([]string(nil), p.Paths...),
		Window:        p.Window,
		Deadline:      p.Deadline,
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, dup := r.byID[p.CorrelationID]; !dup {
		r.order = append(r.order, p.CorrelationID)
		for len(r.order) > MaxTrackedCards {
			drop := r.order[0]
			r.order = r.order[1:]
			delete(r.byID, drop)
		}
	}
	r.byID[p.CorrelationID] = card
}

// Look reads one card back. A missing entry means this surface never displayed
// it, or it already left the screen - in both cases there is no grant to spend.
func (r *Replies) Look(corr string) (ReplyCard, bool) {
	if r == nil {
		return ReplyCard{}, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.byID[corr]
	return c, ok
}

// Forget drops one entry (answered, dismissed, or handed off to execution), so no
// grant stays spendable for a card that is no longer on screen.
func (r *Replies) Forget(corr string) {
	if r == nil || corr == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.byID, corr)
	for i, c := range r.order {
		if c == corr {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
}

// Pending lists the cards this host is holding, oldest displayed first.
func (r *Replies) Pending() []ReplyCard {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]ReplyCard, 0, len(r.order))
	for _, corr := range r.order {
		if c, ok := r.byID[corr]; ok {
			out = append(out, c)
		}
	}
	return out
}

// AwaitingHuman answers ticket 201 AC#6's first half: is a human being waited on
// right now, and which question. It prefers the L2 branch when both are live,
// because an L2 card has a C18 deadline it will die at while an L1 window
// expires into execution.
//
// When a Gate is bound, the queue decides the L2 half: it is the object that
// actually holds the pending items and it drops them the instant one settles.
// The host ledger alone decides the L1 half, because an L1 window is not a queue
// item at all (Gate.PendingWindow owns it) - which is exactly why the ball could
// not know it was in Confirming before this method existed.
func (r *Replies) AwaitingHuman() (ReplyCard, bool) {
	if r == nil {
		return ReplyCard{}, false
	}
	if r.queueAwaitsHuman() {
		for _, c := range r.Pending() {
			if c.Level == "L2" {
				return c, true
			}
		}
		return ReplyCard{}, true
	}
	for _, c := range r.Pending() {
		if c.Level != "" && c.Level != "L1" {
			return c, true
		}
	}
	for _, c := range r.Pending() {
		if c.Level == "L1" {
			return c, true
		}
	}
	return ReplyCard{}, false
}

// WaitingState is AC#6's second half, in the only vocabulary this project has
// for "what the orb should be doing": one of D43's 20 frozen names, produced by
// the run itself rather than by cmd/balldebug's scripted tour.
//
// It never invents a state name and never returns a word outside the table -
// statemachine.Valid is the judge on both branches. A host surface calls this to
// decide what to draw (internal/ball/statevisual.go:177 and :181 already render
// exactly these two rows); a host with no ball calls it to log or to fill a
// snapshot field. ("", false) means nothing is being waited on.
func (r *Replies) WaitingState() (statemachine.State, bool) {
	if r == nil {
		return "", false
	}
	card, awaiting := r.AwaitingHuman()
	if awaiting && card.Level != "L1" {
		// A settled L2 answer, or a queue that holds items the host has not
		// finished displaying: either way there is a C18 card being waited on.
		return statemachine.StateAwaitingApproval, true
	}
	for _, c := range r.Pending() {
		if c.Level == "L1" {
			return statemachine.StateConfirming, true
		}
	}
	return "", false
}

// queueAwaitsHuman reads the bound queue's own pending depth, and answers false
// when no Gate is bound (a ledger alone cannot assert that anyone is waiting).
func (r *Replies) queueAwaitsHuman() bool {
	g, _, _, _ := r.routes()
	if g == nil {
		return false
	}
	return g.Queue().Depth() > 0
}

// Allow spends this card's native grant through the native router. It is the
// method a native button, a tray menu item or a ball click handler calls.
func (r *Replies) Allow(ctx context.Context, corr string) error {
	card, ok := r.Look(corr)
	if !ok {
		return ErrNoTrackedCard
	}
	if card.Grant == "" {
		return ErrRouteHasNoAllow
	}
	g, native, _, _ := r.routes()
	if g == nil {
		return ErrNoGateAttached
	}
	if err := g.DecideFromNative(ctx, Request{
		CorrelationID: corr, Allow: true, Grant: card.Grant, Source: native,
	}); err != nil {
		return err
	}
	r.Forget(corr)
	return nil
}

// AllowSession is the 「本会话内允许」 answer (ticket 224, D45-2's third card
// option): the same tracked card, the same single-use native nonce, and the
// card's own (tool, path) pair recorded as a rule for the rest of this session.
//
// The guardrails are the ones Allow runs, and for the same reasons: an
// untracked correlation is ErrNoTrackedCard, an L1 window has no allow verb at
// all so it is ErrRouteHasNoAllow (SPEC-06 §2 B1 - a countdown window cannot be
// answered, and that does not change because the answer would be remembered),
// and a missing or spent grant is ErrBadGrant from the queue.
//
// There is deliberately no PanelAllowSession. PanelAPI has no Allow method, so
// a panel host has no door to stand in front of; adding a session-flavoured
// sibling of the one method the panel route exists to refuse would be the
// widening this seam is built to prevent.
func (r *Replies) AllowSession(ctx context.Context, corr string) error {
	card, ok := r.Look(corr)
	if !ok {
		return ErrNoTrackedCard
	}
	if card.Grant == "" {
		return ErrRouteHasNoAllow
	}
	g, _, _, _ := r.routes()
	if g == nil {
		return ErrNoGateAttached
	}
	if err := g.Native().AllowSession(ctx, corr, card.Grant); err != nil {
		return err
	}
	r.Forget(corr)
	return nil
}

// Reject refuses through the native router. Refusing needs no proof, which is
// why the panel may do the same thing through PanelReject.
func (r *Replies) Reject(ctx context.Context, corr, reason string) error {
	return r.decide(ctx, corr, reason, false)
}

// PanelReject is the panel route's refusal - SPEC-06 §9 layer 3's 「拒绝」 button,
// and the only decision a page may send.
func (r *Replies) PanelReject(ctx context.Context, corr, reason string) error {
	return r.decide(ctx, corr, reason, true)
}

// routes reads the bound transport under the ledger's lock: one place, so no
// route can be assembled from a half-updated binding.
func (r *Replies) routes() (g *Gate, native, panel string, veto Channel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, veto = r.gate, r.vetoChannel
	native, panel = r.nativeSource, r.panelSource
	if native == "" {
		native = defaultNativeSource
	}
	if panel == "" {
		panel = defaultPanelSource
	}
	return g, native, panel, veto
}

// decide is the one place the two routers are named, so a reviewer can point at a
// line and say which transport a decision came in on.
func (r *Replies) decide(ctx context.Context, corr, reason string, viaPanel bool) error {
	g, native, panel, _ := r.routes()
	if g == nil {
		return ErrNoGateAttached
	}
	req := Request{CorrelationID: corr, Allow: false, Reason: reason}
	if viaPanel {
		req.Source = panel
		if err := g.DecideFromPanel(ctx, req); err != nil {
			return err
		}
	} else {
		req.Source = native
		if err := g.DecideFromNative(ctx, req); err != nil {
			return err
		}
	}
	r.Forget(corr)
	return nil
}

// PanelAllow exists so the panel route's refusal is reachable from a host, not to
// soften it. A WebView host that is handed {allow:true} for approval.decide
// forwards it here and gets back ErrPanelAllow plus a burned nonce: the answer is
// refused on the ROUTE, before any grant is looked at, and a grant that surfaced
// on this path is treated as leaked (gate.go's DecideFromPanel branch,
// queue.revokeGrants). The bool it returns is that statement: true means a grant
// WAS offered on the untrusted route, which is the reading that costs the card
// its spendability. Nothing on this method's happy path exists, because it has
// none - a nil error out of here is a security fault, and the caller is expected
// to say so out loud.
func (r *Replies) PanelAllow(ctx context.Context, corr string) (bool, error) {
	card, _ := r.Look(corr)
	g, _, panel, _ := r.routes()
	if g == nil {
		return card.Grant != "", ErrNoGateAttached
	}
	err := g.DecideFromPanel(ctx, Request{
		CorrelationID: corr, Allow: true, Grant: card.Grant, Source: panel,
	})
	if err == nil {
		return card.Grant != "", errors.New("approval: 面板路线的「允许」没有被路由拒绝")
	}
	if card.Grant != "" {
		// The card STAYS in the ledger - it is still on screen and still pending,
		// and the operator may still refuse it. What a native attempt now meets is
		// ErrBadGrant, which is the point: a token that passed through the panel
		// route is not spendable anywhere else.
		return true, err
	}
	return false, err
}

// Veto cancels an L1 window through the gate's veto funnel, on the channel THIS
// host declared at Attach. It never claims a channel the assembly did not wire.
func (r *Replies) Veto(corr string) error {
	if _, ok := r.Look(corr); !ok {
		return ErrNoTrackedCard
	}
	g, _, _, ch := r.routes()
	if g == nil {
		return ErrNoGateAttached
	}
	if err := g.Veto(Veto{CorrelationID: corr, Channel: ch}); err != nil {
		return err
	}
	r.Forget(corr)
	return nil
}

// Head and View read the queue through the panel-safe projection: what the page
// may show, which by construction carries no grant and no allow (PanelAPI has no
// Allow method at all).
func (r *Replies) Head() (PanelItem, bool) {
	g, _, _, _ := r.routes()
	if g == nil {
		return PanelItem{}, false
	}
	return g.Panel().Head()
}

// View reads one card the same way, by the name the operator was given.
func (r *Replies) View(corr string) (PanelItem, bool) {
	g, _, _, _ := r.routes()
	if g == nil {
		return PanelItem{}, false
	}
	return g.Panel().View(corr)
}

// WideningRule derives the rule text the 长期 branch would store (ticket 201
// AC#4: 「卡上出现规则文本」). It reads the card's own Paths and shapes nothing
// else: exactly one distinct directory yields one line, and anything wider is
// refused rather than guessed at, because a stored rule the operator did not read
// is the failure this ticket was filed for.
//
// The parent segment comes from this package's existing dirOf (batch.go:91),
// which splits an already-canonical path by pure string work and deliberately
// calls no filepath.* — AGENTS.md §1.2 ban #2 is on making path DECISIONS outside
// risk.PathResolver, and Prompt.Paths are already the C26 output. What comes back
// is one line of TEXT for [fs] allowed_dirs, which the next start canonicalises
// again through C26 like any hand edit of the same section, and which a second,
// L2-level confirmation has to approve before anything is written at all.
func (c ReplyCard) WideningRule() (dir, rule string, ok bool) {
	dirs := make([]string, 0, len(c.Paths))
	for _, p := range c.Paths {
		d := strings.ReplaceAll(dirOf(p), "\\", "/")
		// A bare drive root ("D:", ".") and an empty segment are refusals, not
		// offers: allowlisting a whole volume is not what one card asked for.
		if d == "" || d == "." || isBareDrive(d) || slices.Contains(dirs, d) {
			continue
		}
		dirs = append(dirs, d)
	}
	if len(dirs) != 1 {
		return "", "", false
	}
	dir = dirs[0]
	return dir, "[fs] allowed_dirs += \"" + dir + "\"", true
}

// isBareDrive reports a Windows drive root written without its trailing slash
// ("D:", "c:"), the one shape this package's dirOf can return that must not
// become a stored allowlist entry.
func isBareDrive(s string) bool {
	return len(s) == 2 && s[1] == ':' &&
		((s[0] >= 'A' && s[0] <= 'Z') || (s[0] >= 'a' && s[0] <= 'z'))
}
