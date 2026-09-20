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

// channelNames are the display labels the card and the ball strip show.
var channelNames = map[Channel]string{
	ChannelBall:  "单击悬浮球",
	ChannelEsc:   "按 Esc 键",
	ChannelPanel: "面板拒绝",
	ChannelKWS:   "说取消词",
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
		return "Esc 取消不可用"
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
			st.Text = channelNames[ch]
		} else {
			st.Text = channelNames[ch] + "：" + unavailableText(ch)
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

// spend validates and consumes. It reports why a rejection happened in
// user-safe terms: never "close, you got 3 bytes right".
func (s *grantStore) spend(nonce, bind string) bool {
	if nonce == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for v, stored := range s.values {
		if !equalSecret(v, nonce) {
			continue
		}
		delete(s.values, v) // consumed whether or not the binding matched
		return equalSecret(stored, bind)
	}
	return false
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
