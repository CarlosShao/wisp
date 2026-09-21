package approval

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/tools"
)

// itemState is one queue entry's lifecycle. An item that left pending can
// never be answered again - that is what makes a dismissed card, an expired
// request and a spent token all refuse the same way.
type itemState int

const (
	statePending itemState = iota
	stateAnswered
	stateDropped
)

// answer is what the queue delivers to the waiting call.
type answer struct {
	a   tools.Answer
	why string
}

// qitem is one pending L2 approval (C18 at trivial scale).
type qitem struct {
	Dec  tools.Decision
	Corr string
	Seq  uint64

	// names are the OTHER strings this one card can legitimately be addressed
	// by: the correlation id as it arrived before the queue re-stamped it, and
	// the task id the host keys its own bookkeeping on (tools/bridge.go routes
	// cancels by orDefault(CorrelationID, TaskID)). Ticket 87: a card that is
	// on screen is addressable by whatever the display was built from, and a
	// reply carrying one of those names is not a reply for a different request.
	names []string

	// bind is the digest a grant must match to authorize THIS item. It covers
	// the correlation id, the task, the tool, the level, the argument bytes
	// and the sequence number, so a nonce cannot be replayed onto a different
	// request or onto a later re-issue of the same one.
	bind   string
	grants *grantStore

	state  itemState
	answer chan answer

	// replayOf records which request re-opened this one (一键重放).
	replayOf string
}

// Queue is the single-task ApprovalQueue: a bounded FIFO whose HEAD is the one
// item displayed (C18 「队头单显」+ the depth badge), with a 300s deadline that
// always resolves to REJECT. It is safe for concurrent use: the native API,
// the panel API and up to MaxToolConcurrency waiting calls touch it at once.
type Queue struct {
	mu      sync.Mutex
	timeout time.Duration
	warn    time.Duration
	maxPend int
	maxRepl int
	seq     uint64
	pending []*qitem
	byID    map[string]*qitem
	// alias is the reject-direction index: every extra name a live card answers
	// to (qitem.names) plus the item that owns it. The allow side never reads it
	// - a grant is only ever spendable on the exact key the queue issued - so
	// the index can only ever make a REFUSAL land, never an approval.
	alias   map[string]map[*qitem]bool
	history []*qitem
	logf    func(string, ...any)
}

// NewQueue builds the trivial C18 queue. Zero durations get the contract
// defaults (300s / last-30s warning); a negative timeout is a bug, not a
// disable switch, so it is replaced by the default too.
func NewQueue(timeout, warnBefore time.Duration, maxPending int, logf func(string, ...any)) *Queue {
	if timeout <= 0 {
		timeout = DefaultApprovalTimeout
	}
	if warnBefore <= 0 || warnBefore >= timeout {
		warnBefore = DefaultApprovalWarning
	}
	if maxPending <= 0 {
		maxPending = DefaultMaxPending
	}
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &Queue{
		timeout: timeout, warn: warnBefore, maxPend: maxPending,
		maxRepl: DefaultReplayHistory, byID: map[string]*qitem{},
		alias: map[string]map[*qitem]bool{}, logf: logf,
	}
}

// Contract defaults (SPEC-06 §7 / ticket 21).
const (
	// DefaultApprovalTimeout: 「超时 300s 一律判拒绝」. Never an infinite wait.
	DefaultApprovalTimeout = 300 * time.Second
	// DefaultApprovalWarning: 「超时前 30s 醒目提示」.
	DefaultApprovalWarning = 30 * time.Second
	// DefaultMaxPending bounds the FIFO. Overflow is refused fail-closed: a
	// broken/overflowed queue must never be a reason to run something.
	DefaultMaxPending = 8
	// DefaultReplayHistory is how many answered requests stay replayable.
	DefaultReplayHistory = 8
	// DefaultL1Window is the 2-3s pre-execution block (SPEC-06 §2, B1).
	DefaultL1Window = 3 * time.Second
	// MinL1Window / MaxL1Window bound a configured window. A countdown the
	// user cannot read is not a gate, and one long enough to become a modal
	// wall is not a window either.
	MinL1Window = 2 * time.Second
	// MaxL1Window bounds the L1 block window.
	MaxL1Window = 3 * time.Second
)

// Timeout exposes the queue deadline for prompts and audit lines.
func (q *Queue) Timeout() time.Duration { return q.timeout }

// WarningLead exposes the last-30s lead.
func (q *Queue) WarningLead() time.Duration { return q.warn }

// Depth is the number of pending items; the ball's AwaitingApproval badge
// shows it (depth=1 at trivial scale).
func (q *Queue) Depth() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.pending)
}

// push admits one decision into the queue and assigns the correlation id the
// reply is routed by (C18). Overflow fails closed.
func (q *Queue) push(d tools.Decision) (*qitem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.pending) >= q.maxPend {
		return nil, fmt.Errorf("审批队列已满（%d 项待决），无法登记新的 L2 请求，已 fail-closed 拒绝", q.maxPend)
	}
	q.seq++
	corr := d.CorrelationID
	if corr == "" {
		corr = "approval-" + strconv.FormatUint(q.seq, 10)
	}
	if _, clash := q.byID[corr]; clash {
		corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)
	}
	it := &qitem{
		Dec: d, Corr: corr, Seq: q.seq,
		names:  otherNames(d, corr),
		grants: newGrantStore(),
		state:  statePending,
		answer: make(chan answer, 1),
	}
	it.bind = bindDigest(corr, d.TaskID, d.Tool, d.LevelString(), q.seq, d.Args)
	q.pending = append(q.pending, it)
	q.byID[corr] = it
	q.indexLocked(it)
	q.logf("approval: queued corr=%s task=%s tool=%s depth=%d",
		corr, d.TaskID, d.Tool, len(q.pending))
	return it, nil
}

// otherNames are the reject-direction names of one card: everything it may be
// addressed by that is NOT the key the queue just issued. An empty or
// duplicate name is not an alias - it is the same string.
func otherNames(d tools.Decision, corr string) []string {
	var out []string
	for _, n := range []string{d.CorrelationID, d.TaskID} {
		if n == "" || n == corr {
			continue
		}
		dup := false
		for _, have := range out {
			if have == n {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, n)
		}
	}
	return out
}

// indexLocked registers one item's reject-direction names. Caller holds q.mu.
func (q *Queue) indexLocked(it *qitem) {
	for _, n := range it.names {
		set := q.alias[n]
		if set == nil {
			set = map[*qitem]bool{}
			q.alias[n] = set
		}
		set[it] = true
	}
}

// unindexLocked removes one item from the alias index (it left the queue).
// Caller holds q.mu.
func (q *Queue) unindexLocked(it *qitem) {
	for _, n := range it.names {
		set := q.alias[n]
		if set == nil {
			continue
		}
		delete(set, it)
		if len(set) == 0 {
			delete(q.alias, n)
		}
	}
}

// resolveLocked finds the live pending item a reply names. The queue key always
// wins. strict is the posture of the side that must not be forgiving: an allow
// passes true and gets nothing but its exact key. The reject side passes
// false, which additionally reads the alias index - the only direction a
// borrowed name can move a call in is 「do not run it」, so leniency here can
// make a refusal land sooner and can never let something through.
//
// An alias naming more than one live item is NOT guessed: with two cards on
// screen that a reply cannot tell apart, the reply means neither of them and
// stays an unknown correlation id. Caller holds q.mu.
func (q *Queue) resolveLocked(name string, strict bool) *qitem {
	if name == "" {
		return nil
	}
	if it := q.byID[name]; it != nil && it.state == statePending {
		return it
	}
	if strict {
		return nil
	}
	if set := q.alias[name]; len(set) == 1 {
		for it := range set {
			if it.state == statePending {
				return it
			}
		}
	}
	return nil
}

// position returns one item's 1-based place in the FIFO (the depth badge).
func (q *Queue) position(it *qitem) int {
	for i, p := range q.pending {
		if p == it {
			return i + 1
		}
	}
	return 1
}

// drop removes an item from the pending set and records it for replay.
// Caller holds q.mu.
func (q *Queue) dropLocked(it *qitem, keepForReplay bool) {
	for i, p := range q.pending {
		if p == it {
			q.pending = append(q.pending[:i], q.pending[i+1:]...)
			break
		}
	}
	delete(q.byID, it.Corr)
	q.unindexLocked(it)
	it.grants.revoke()
	if keepForReplay {
		q.history = append(q.history, it)
		for len(q.history) > q.maxRepl {
			q.history = q.history[1:]
		}
	}
}

// deliver settles a pending item exactly once.
func (q *Queue) deliver(it *qitem, a answer) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if it.state != statePending {
		return false
	}
	it.state = stateAnswered
	q.dropLocked(it, true)
	select {
	case it.answer <- a:
	default:
	}
	return true
}

// takeAnswer drains an already-settled answer without waiting.
func (q *Queue) takeAnswer(it *qitem) (answer, bool) {
	select {
	case a := <-it.answer:
		return a, true
	default:
		return answer{}, false
	}
}

// grantNonce mints the single-use native proof for one item. An error means
// the RNG failed: no prompt may be displayed, because a displayed prompt
// nobody can authorize is worse than an honest refusal.
func (q *Queue) grantNonce(it *qitem) (string, error) {
	nonce, err := mintGrant()
	if err != nil {
		q.mu.Lock()
		it.state = stateDropped
		q.dropLocked(it, false)
		q.mu.Unlock()
		return "", err
	}
	q.mu.Lock()
	it.grants.issue(nonce, it.bind)
	q.mu.Unlock()
	return nonce, nil
}

// allow is the native-only answer path. It spends a grant and refuses anything
// else, including an empty one.
func (q *Queue) allow(corr, nonce string) error {
	q.mu.Lock()
	it, ok := q.byID[corr]
	if !ok {
		q.mu.Unlock()
		return ErrUnknownCorrelation
	}
	if it.state != statePending {
		q.mu.Unlock()
		return ErrNotPending
	}
	spent := it.grants.spend(nonce, it.bind)
	q.mu.Unlock()
	if !spent {
		q.logf("approval: FORGED-OR-STALE allow rejected corr=%s (native grant missing/spent/misbound)", corr)
		return ErrBadGrant
	}
	if !q.deliver(it, answer{a: tools.AnswerAllow, why: "用户在原生侧批准了本次操作"}) {
		return ErrNotPending
	}
	return nil
}

// revokeGrants burns every live nonce of one item without answering it. A
// grant that surfaced on the untrusted route is treated as leaked: the honest
// card then has to be re-displayed (which mints a fresh one), so a token that
// passed through a compromised panel can never be spent.
func (q *Queue) revokeGrants(corr string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if it, ok := q.byID[corr]; ok {
		it.grants.revoke()
	}
}

// reject closes one item as refused. No proof is required: refusing is the
// fail-closed direction, which is why the panel may do it (SPEC-06 §9), and it
// is the single funnel every refusal route goes through - native, panel, and
// since ticket 87 the veto channel too. That is deliberate: the AnswerReject
// literal below is the one place the human-said-no direction is written, so
// turning it into an allow is a mutation the R7/C18 fail-closed family has to
// survive rather than a per-route decision that can be re-forgotten.
//
// The lookup is the lenient one (see resolveLocked): a reply that cannot be
// answered is a lost vote, and the only thing a borrowed name can buy here is
// that the refusal lands sooner.
func (q *Queue) reject(corr, reason string) error {
	q.mu.Lock()
	it := q.resolveLocked(corr, false)
	q.mu.Unlock()
	if it == nil {
		return ErrUnknownCorrelation
	}
	why := reason
	if why == "" {
		why = "用户拒绝了本次操作"
	}
	if !q.deliver(it, answer{a: tools.AnswerReject, why: why}) {
		return ErrNotPending
	}
	return nil
}

// expire auto-rejects a deadline that nobody answered.
func (q *Queue) expire(it *qitem) answer {
	a := answer{
		a:   tools.AnswerTimeout,
		why: fmt.Sprintf("审批超时（%d 秒未确认），C18 一律判拒绝，已自动拒绝", int(q.timeout.Seconds())),
	}
	if q.deliver(it, a) {
		return a
	}
	if got, ok := q.takeAnswer(it); ok {
		return got
	}
	return a
}

// abandon closes an item because the waiting call went away (task cancelled).
// The task's own context is the caller's; this package never cancels it.
func (q *Queue) abandon(it *qitem, why string) answer {
	a := answer{a: tools.AnswerReject, why: why}
	if q.deliver(it, a) {
		return a
	}
	if got, ok := q.takeAnswer(it); ok {
		return got
	}
	return a
}

// view renders the panel-safe projection of one item. It is built from the
// stored decision and has no way to carry a grant even if one were live.
func (q *Queue) view(corr string) (PanelItem, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	it, ok := q.byID[corr]
	if !ok || it.state != statePending {
		return PanelItem{}, false
	}
	return q.viewLocked(it), true
}

func (q *Queue) viewLocked(it *qitem) PanelItem {
	return PanelItem{
		CorrelationID: it.Corr,
		Tool:          it.Dec.Tool,
		Level:         it.Dec.LevelString(),
		Reason:        it.Dec.Reason,
		Paths:         append([]string(nil), it.Dec.Paths...),
		Depth:         q.position(it),
	}
}

// head is the single displayed item (C18 队头单显).
func (q *Queue) head() (PanelItem, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.pending) == 0 {
		return PanelItem{}, false
	}
	return q.viewLocked(q.pending[0]), true
}

// pendingCount reports live pending items - the D47 assertion seam.
func (q *Queue) pendingCount() int { return q.Depth() }

// replay re-opens an answered/rejected request and returns it under a fresh
// correlation id. Re-displaying a question is not answering it: the caller
// must run the approval flow again, which mints a NEW grant, so a token spent
// (or refused) on the first display cannot carry over.
func (q *Queue) replay(corr string) (tools.Decision, string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, live := q.byID[corr]; live {
		return tools.Decision{}, "", errors.New("该审批项仍在等待答复，无需重放")
	}
	for i, it := range q.history {
		if it.Corr != corr {
			continue
		}
		fresh := *it
		fresh.grants = newGrantStore()
		fresh.state = statePending
		fresh.answer = make(chan answer, 1)
		fresh.replayOf = it.Corr
		d := it.Dec
		d.CorrelationID = corr + "#replay"
		if _, clash := q.byID[d.CorrelationID]; clash {
			d.CorrelationID = fmt.Sprintf("%s#%d", corr, q.seq+uint64(i)+1)
		}
		return d, d.CorrelationID, nil
	}
	return tools.Decision{}, "", fmt.Errorf("无可重放记录: %s", corr)
}
