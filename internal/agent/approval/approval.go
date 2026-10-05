package approval

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Clock: the layer's only source of time
// ---------------------------------------------------------------------------

// Clock is the injected monotonic time source (D42#9 / D22 ban 5). Every
// timeout in this package - the 2-3s L1 window, the 300s queue deadline, the
// last-30s warning - is a duration handed to After, never a difference between
// two wall-clock readings. Ticket 11's rate limiter established the pattern
// (internal/llm/ratelimit.go:33-45): production passes time.Now/time.After,
// whose readings carry the monotonic clock; tests pass a fake they can advance.
type Clock interface {
	// Now is for stamps and audit lines ONLY. No timeout decision may be made
	// by comparing two Now readings.
	Now() time.Time
	// After returns a channel that receives once d has elapsed.
	After(d time.Duration) <-chan time.Time
}

// SystemClock is the production Clock.
type SystemClock struct{}

// Now implements Clock.
func (SystemClock) Now() time.Time { return time.Now() }

// After implements Clock.
func (SystemClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// ---------------------------------------------------------------------------
// Veto channels (SPEC-06 §2, B1)
// ---------------------------------------------------------------------------

// Channel is one of the four L1 veto channels. The set is closed: a fifth
// value means a caller invented a selector, which is exactly the M-7/C-3
// failure shape (a security decision keyed on an unvalidated label).
type Channel string

// The four veto channels of SPEC-06 §2.
const (
	ChannelBall  Channel = "ball"  // click the floating ball
	ChannelEsc   Channel = "esc"   // global Esc hotkey
	ChannelPanel Channel = "panel" // panel reject button (ticket 37)
	ChannelKWS   Channel = "kws"   // spoken veto word (ticket 41)
)

// allChannels is the closed set, in display order.
var allChannels = []Channel{ChannelBall, ChannelEsc, ChannelPanel, ChannelKWS}

// AllChannels returns the four veto channels in display order.
func AllChannels() []Channel {
	out := make([]Channel, len(allChannels))
	copy(out, allChannels)
	return out
}

// Veto is one cancel attempt. CorrelationID routing is a C18 requirement
// (「回复按 correlationId 路由，点击不可能落到别的请求上」), so a veto that
// names no live window is not a cancel - it is an error, and the window
// keeps counting down.
type Veto struct {
	CorrelationID string
	Channel       Channel
}

// ChannelStatus is the honest, user-visible availability line for one channel.
// An unloaded channel is NEVER rendered as available (B1).
type ChannelStatus struct {
	Channel Channel
	Loaded  bool
	Text    string
}

// channelNames are the display labels the card and the ball strip show, and
// EVERY entry here names a CHANNEL, never a key. The cancel channel's live
// label - the one that names the key this host actually borrows - is resolved
// at read time by channelLabel, NOT frozen in this map (ticket 260 AC#4,
// ledger A598 §2). Two facts made that true:
//
//  1. since ticket 260 形ⓐ the borrowed key follows [hotkey] cancel, and ticket
//     258's reload bridge rebinds it while the process runs - a label computed
//     once into a package var is the stale half of the exact defect ticket 260
//     was filed over ("the card says one key, the desktop holds another");
//  2. the key's spelling is owned by internal/ball's hotkey receipt, and this
//     package must not import internal/ball (a new package-level dependency
//     edge = human-approval face), so the spelling arrives through the
//     composition-root seam below instead.
//
// ChannelEsc's entry is therefore the SLOT name: it is what the availability
// line prints when the channel is NOT loaded, i.e. when there is no key to
// name. That is A595 §1 boundary ④ / 档一, and it is why this line and the
// 「取消键无处可借」 the resident leg prints for the same branch agree.
var channelNames = map[Channel]string{
	ChannelBall:  "单击悬浮球",
	ChannelEsc:   "取消键通道",
	ChannelPanel: "面板拒绝",
	ChannelKWS:   "说取消词",
}

// defaultCancelKeySpelling is the fallback this package prints when its host
// installed no reader (a host with no hotkey chain at all - cmd/wisp's console
// leg, and every test that never wired one). It is the SHIPPED DEFAULT KEY and
// nothing else, and it is duplicated here on purpose rather than imported:
// internal/ball's DefaultHotkeys().Cancel is the authority, and two cmd/wisp
// rulers keep the two from drifting silently -
// TestTicket260R4ShippedConstructorInstallsTheReader opens with the premise
// check `ball.DefaultHotkeys().Cancel == "Esc"` (move that ball default and the
// cell goes red instead of the card quietly naming the old key), and
// TestTicket260R3CancelKeyComesFromTheBallChain pins that the reader with no
// window answers with that same ball constant, never a copied literal. So a
// drift between them is a red test, not a sentence on the card naming a key
// nobody holds.
const defaultCancelKeySpelling = "Esc"

// cancelKeyLabelPattern wraps a key spelling in the words the card used before
// ticket 260 ever moved: with the default spelling it renders 「按 Esc 键」,
// byte-for-byte the string this package printed for years. Only the spelling
// moves; the shape around it does not (AC#4 ① zero drift).
const cancelKeyLabelPattern = "按 %s 键"

// CancelKeySpelling answers "which key does this process call the cancel key",
// in the host's own spelling (「Esc」, 「Ctrl+Alt+Q」, ...). It is the whole of
// what the composition root injects: the value comes from the ball's hotkey
// report - the same line internal/ball fills from the cancelBorrow it hands to
// RegisterHotKey - so the printed name and the borrowed key cannot disagree.
//
// It is a display seam and NOTHING more: no decision in this package reads it.
// Which veto lands is decided by ChannelRegistry (loadedness) and the grant
// proof, so a host that installed a reader answering「F13」 still cannot cancel
// through a channel it never loaded, and cannot make an unloaded channel look
// available. Pinned by TestTicket260R4SpellingSeamCarriesNoAuthority.
type CancelKeySpelling func() string

var (
	// cancelKeyMu guards the incumbent reader, not any decision. The read faces
	// (Statuses / gate.go's two audit lines / report.go's channel attribution)
	// are called from the gate's own goroutines AND from the UI thread, so the
	// map this replaces could be read unlocked while the root re-wired it.
	cancelKeyMu sync.RWMutex
	// cancelKeySrc is nil until a composition root installs one.
	cancelKeySrc CancelKeySpelling
)

// SetCancelKeySpelling installs (or, with nil, removes) the host's reader and
// returns the incumbent so a caller can restore it. Callers are composition
// roots only; the shape is the one winsec.SetPathResolver established for a
// package-level seam the provider package must not import.
func SetCancelKeySpelling(src CancelKeySpelling) CancelKeySpelling {
	cancelKeyMu.Lock()
	defer cancelKeyMu.Unlock()
	prev := cancelKeySrc
	cancelKeySrc = src
	return prev
}

// cancelKeySpelling resolves the spelling at read time: no installed reader and
// an empty answer from one both land on the shipped default, because 「按  键」
// is a sentence that names nothing and reads like a broken card.
func cancelKeySpelling() string {
	cancelKeyMu.RLock()
	src := cancelKeySrc
	cancelKeyMu.RUnlock()
	if src == nil {
		return defaultCancelKeySpelling
	}
	if s := src(); s != "" {
		return s
	}
	return defaultCancelKeySpelling
}

// cancelKeyLabel is the loaded form of the cancel channel's label: the key this
// host really holds, in the words the card has always used.
func cancelKeyLabel() string { return fmt.Sprintf(cancelKeyLabelPattern, cancelKeySpelling()) }

// channelLabel renders a channel whose label is being printed as a thing that IS
// live or DID happen: the availability line of a loaded channel, and the three
// attribution faces (gate.go's L1 and L2 veto sentences, report.go's
// 「（否决通道：…）」). A veto that reached those lines passed the registry check,
// so naming a key there is honest - as long as the key it names is THIS host's,
// which is the whole of what this function adds.
func channelLabel(ch Channel) string {
	if ch == ChannelEsc {
		return cancelKeyLabel()
	}
	return channelNames[ch]
}

// channelStatusLabel renders the label half of one availability line. The
// cancel channel splits in two, and the split is by the only fact the line
// claims: loaded -> name the key this host holds; not loaded -> name the SLOT,
// because there is no borrowed key to point at and the second half of the line
// says so out loud. Before this split the pair rendered 「按 Esc 键：快捷键取消
// 不可用」 - the first clause naming a key the machine is not holding, in the
// one moment (card would not come up / key never got borrowed) a user reads it.
func channelStatusLabel(ch Channel, loaded bool) string {
	if ch == ChannelEsc && !loaded {
		return channelNames[ch]
	}
	return channelLabel(ch)
}

// unavailableText is the wording B1 mandates. The KWS line is the exact string
// the spec requires (「语音取消不可用」) and a test asserts it verbatim: the
// strip must say this, never stay silent about a channel that cannot fire.
//
// DEFERRED(kws-veto, B1): the veto word is functional at ticket 41 and the
// panel button at ticket 37 -> SPEC-12 §5 row「DEFERRED | 快捷键路径语音否决
// （B1）| KWS 未加载+ASR 加载 1–3s > 2–3s 窗口，物理不可能 | 快捷键会话中说
// 「取消」能否决 L1 | AEC | Confirming 明示「语音取消不可用」」.
func unavailableText(ch Channel) string {
	switch ch {
	case ChannelKWS:
		return "语音取消不可用"
	case ChannelPanel:
		return "面板取消不可用（票 37 未接入）"
	case ChannelBall:
		return "悬浮球取消不可用"
	case ChannelEsc:
		// Ticket 260 AC#1 condition 2, ledger A595 §1 boundary ④: this line is
		// reached ONLY when the channel is not loaded, so there is no key to
		// name. Saying 「Esc」 here was a claim about a key this host is not
		// holding at all - and since the borrow now follows [hotkey] cancel, the
		// named key could not even be inferred from the default. The honest form
		// names the SLOT, never a key, matching the "无处可借／不可用" wording
		// the resident leg uses for the same two branches.
		return "快捷键取消不可用"
	default:
		return "该取消通道不可用"
	}
}

// ErrChannelUnavailable means a veto arrived on a channel the host never
// loaded. It is returned, not swallowed: a silently ignored cancel attempt is
// how a fake channel becomes a user-visible promise.
var ErrChannelUnavailable = errors.New("approval: 取消通道未加载")

// ChannelError wraps ErrChannelUnavailable with the user-visible wording.
type ChannelError struct {
	Channel Channel
	Msg     string
}

// Error implements error.
func (e *ChannelError) Error() string { return e.Msg + ": " + string(e.Channel) }

// Unwrap lets errors.Is(err, ErrChannelUnavailable) answer for any unloaded
// channel without string matching.
func (e *ChannelError) Unwrap() error { return ErrChannelUnavailable }

// ChannelRegistry is the host's authoritative statement of which veto channels
// are loaded. Loadedness is NOT taken from the veto request: the caller cannot
// talk its way into being treated as available, and the composition root flips
// the flag only when the backing subsystem is really up (KWS model loaded,
// panel window created).
type ChannelRegistry struct {
	mu     sync.Mutex
	loaded map[Channel]bool
}

// NewChannels builds a registry from the channels that are loaded today.
func NewChannels(loaded ...Channel) *ChannelRegistry {
	m := make(map[Channel]bool, len(allChannels))
	for _, ch := range allChannels {
		m[ch] = false
	}
	for _, ch := range loaded {
		if _, ok := channelNames[ch]; ok {
			m[ch] = true
		}
	}
	return &ChannelRegistry{loaded: m}
}

// DefaultChannels is today's honest answer: the ball and the Esc hotkey are
// wired, the panel (ticket 37) and KWS (ticket 41) are not.
func DefaultChannels() *ChannelRegistry {
	return NewChannels(ChannelBall, ChannelEsc)
}

// Loaded reports whether one channel is live.
func (r *ChannelRegistry) Loaded(ch Channel) bool {
	if r == nil {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loaded[ch]
}

// SetLoaded flips one channel's availability (the KWS loader calls this when
// its model is actually resident).
func (r *ChannelRegistry) SetLoaded(ch Channel, ok bool) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := channelNames[ch]; exists {
		r.loaded[ch] = ok
	}
}

// Statuses renders the four channels with their availability lines.
func (r *ChannelRegistry) Statuses() []ChannelStatus {
	out := make([]ChannelStatus, 0, len(allChannels))
	for _, ch := range allChannels {
		loaded := r.Loaded(ch)
		st := ChannelStatus{Channel: ch, Loaded: loaded}
		if loaded {
			st.Text = channelStatusLabel(ch, true)
		} else {
			st.Text = channelStatusLabel(ch, false) + "：" + unavailableText(ch)
		}
		out = append(out, st)
	}
	return out
}

// check returns a ChannelError when ch is not one of the four, or is not
// loaded. Every veto path in this package goes through here.
func (r *ChannelRegistry) check(ch Channel) error {
	if _, ok := channelNames[ch]; !ok {
		return &ChannelError{Channel: ch, Msg: "未知取消通道，已忽略"}
	}
	if r.Loaded(ch) {
		return nil
	}
	return &ChannelError{Channel: ch, Msg: unavailableText(ch)}
}

// ---------------------------------------------------------------------------
// Grant: the unforgeable native-source proof (F2 layer 3)
// ---------------------------------------------------------------------------

// A grant is a single-use nonce minted by the gate AT THE MOMENT a native
// prompt is displayed and handed to exactly one recipient: the native UI
// argument of UI.Prompt. The panel-facing surface (PanelItem) has no grant
// field and PanelAPI has no Allow method, so the untrusted path cannot carry
// one - and cannot mint one, because minting is unexported and the value comes
// from crypto/rand.
//
// Authority for an allow answer is therefore: (1) the answer arrived on the
// native API surface, (2) it presented a live grant for THIS pending item, and
// (3) the grant's binding digest matches the item it is spent on. The
// caller-settable Request.Source string is logged and never consulted
// (rulings M-7 / C-3 in docs/reports/pending-and-issues.md: a security
// decision keyed on a caller-controlled selector fails open).
const grantBytes = 32

// grantPrefix exists so an audit line can name a grant without carrying it.
const grantPrefix = "grant_"

// mintGrant returns a fresh nonce. An error here means the RNG broke, which is
// a broken queue: the caller fails closed rather than displaying a prompt
// nobody can authorize.
func mintGrant() (string, error) {
	buf := make([]byte, grantBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("approval: 无法生成审批令牌（随机源故障），已 fail-closed: %w", err)
	}
	return grantPrefix + hex.EncodeToString(buf), nil
}

// bindDigest is the domain separator that ties a grant to ONE pending item.
// Two items with the same correlation id at different times (a replay) still
// differ, because the sequence number is folded in.
func bindDigest(corr, taskID, tool, level string, seq uint64, args []byte) string {
	sum := sha256.New()
	for _, s := range []string{corr, taskID, tool, level, strconv.FormatUint(seq, 10)} {
		_, _ = sum.Write([]byte{byte(len(s))})
		_, _ = sum.Write([]byte(s))
	}
	_, _ = sum.Write(args)
	return hex.EncodeToString(sum.Sum(nil))
}

// equalSecret is a constant-time comparison. A plain == on a token would make
// the reject reason depend on how many leading bytes matched, which is a
// timing oracle over a value the whole design says the attacker cannot know.
func equalSecret(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// grantStore holds the live nonces of one pending item. Every entry is
// single-use: spending deletes it, and so does answering or expiring, so a
// screenshot of a dismissed card cannot authorize anything afterwards.
type grantStore struct {
	mu     sync.Mutex
	values map[string]string // nonce -> binding digest
}

func newGrantStore() *grantStore { return &grantStore{values: map[string]string{}} }

// issue registers one minted nonce for this item.
func (s *grantStore) issue(nonce, bind string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[nonce] = bind
}

// grantDenial is the INTERNAL answer to one question: why did this allow not
// spend a grant (ticket 259 AC#2). Before this type existed the answer was a
// bool, so "no proof presented", "proof not live on this card", "proof live but
// bound to something else" and "the card already left pending" folded into one
// reply, and the audit line could only list them ("missing/spent/misbound")
// instead of naming which one happened.
//
// Scope, stated because the two halves are NOT the same decision:
//
//   - audit side (this type, and Queue.allowScoped's GRANT-DENY line): the four
//     causes are named apart, so an auditor attributes a refusal by reading a
//     line instead of inferring it from the absence of three other lines;
//   - API and UI side: unchanged on purpose. Every grant denial still leaves
//     this package as the one error value ErrBadGrant with the one merged
//     sentence (ui.go), and no caller - native, panel, host router or a page -
//     can read a grantDenial from out here. That merge is the anti-probe
//     property the ErrBadGrant comment states, and ticket 259 does not move it.
//
// It is unexported for that reason: exporting it would hand the distinction to
// exactly the faces the merge exists to keep it from.
type grantDenial uint8

// The denials. denialNone is the zero value and means "no denial": the spend
// succeeded. The others are the four causes of ticket 259 AC#2, one apiece.
const (
	denialNone grantDenial = iota
	// denialMissingNonce: no proof was presented at all (the empty value).
	denialMissingNonce
	// denialSpentNonce: the presented value is not live in THIS card's own
	// store. One label covers "already spent" and "never issued" deliberately:
	// spending deletes the row, so telling them apart needs a kept roster of
	// dead nonces, which is a new state face and not what AC#2 asked for.
	denialSpentNonce
	// denialMisbound: the proof IS live on this card, but the binding digest
	// it was issued under does not match the one it is presented against.
	denialMisbound
	// denialNotPending: the item exists but already left the pending set, so no
	// proof of any quality could be spent on it. Produced by the queue (the
	// store has no view of item state), not by spend.
	denialNotPending
)

// label is the token the audit line prints. It is a machine-readable name, not
// user-facing prose, and it never appears in any error value or card text.
func (d grantDenial) label() string {
	switch d {
	case denialNone:
		return "none"
	case denialMissingNonce:
		return "missing-nonce"
	case denialSpentNonce:
		return "spent-or-never-live-nonce"
	case denialMisbound:
		return "misbound"
	case denialNotPending:
		return "card-not-pending"
	default:
		return "unclassified-denial"
	}
}

// spend validates and consumes, and returns WHICH denial it reached (ticket
// 259 AC#2): denialNone means the grant opened the card, anything else means it
// did not and says why in audit terms. The control flow is the pre-259 one,
// byte-for-byte equivalent in outcome: the nonce is consumed whether or not the
// binding matched (a rejected caller cannot retry), and an unknown value
// deletes nothing. What changed is only that the branch taken is now reportable.
func (s *grantStore) spend(nonce, bind string) grantDenial {
	if nonce == "" {
		return denialMissingNonce
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for v, stored := range s.values {
		if !equalSecret(v, nonce) {
			continue
		}
		delete(s.values, v) // consumed whether or not the binding matched
		if !equalSecret(stored, bind) {
			return denialMisbound
		}
		return denialNone
	}
	return denialSpentNonce
}

// revoke drops every live nonce for the item (answered, expired, host gone).
func (s *grantStore) revoke() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = map[string]string{}
}

// live counts unconsumed nonces - an assertion seam for the tests.
func (s *grantStore) live() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.values)
}

// sortedStrings is a tiny deterministic-render helper for audit text.
func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
