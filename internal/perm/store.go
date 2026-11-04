// Package perm owns the user-facing permission mode: the read the enforcement
// chain consults, the switch that persists, the one confirmation a switch costs,
// and the audit record every switch leaves.
//
// Ticket 90, ruling R20. The shape is three narrow jobs, and the boundaries are
// the point of the package:
//
//   - READ (PermissionMode): what the bridge asks once per tool call. It never
//     mutates, and it fail-closes to the strictest档 when the config is missing
//     or unreadable. It answers the EFFECTIVE档: when this store has no
//     confirmation channel, a configured auto_approve / ask_high_risk reads as
//     ask_every_step (ticket 372-r2). The configured value is never rewritten.
//   - WRITE (Set): the only function in the system that changes the mode at
//     runtime. Its gate is the TARGET, not the direction: a switch INTO
//     auto_approve is refused outright when no confirmation channel was wired
//     (fail-closed), and every other target - wider or narrower - passes
//     (R20/M4, owner 2026-09-21: only the switch into全自动 costs a card).
//     On top of that gate the auto_approve switch still costs one ANSWERED L2
//     strong confirmation: a channel that is merely attached is not enough, the
//     limb has to come back with a Card (ticket 371 - "non-nil" was what let a
//     never-asking leg through silently, docs/evidence/s1/371-a1p-impl.md §5).
//     Ticket 367 judged the DIRECTION instead; the survey at
//     docs/evidence/s1/m4-survey.md:96 found no mainstream harness gating on
//     "wider", so ticket 372 reverted the judgement while keeping 367's
//     single-definition discipline (JudgeCard/RequiresCard: one definition, both
//     transition points).
//   - RECORD: every switch writes one audit line through the sink the project
//     already uses for security decisions (the Logf family that cmd/wisp renders
//     as "[audit] ..."), carrying from/to, time, origin and actor. No second
//     audit subsystem is introduced here, and no switch is silent.
//
// What this package deliberately does NOT do:
//
//   - It does not touch authorizations. R20/M3 persists the MODE, one global
//     preference; it does not make a D45 session grant persistent. A grant stays
//     session-scoped and dies with the process (PLAN.md:1640), and ticket 49's
//     GrantScopeSession keeps that meaning. The two are pinned by two separate
//     tests (AC#3a and AC#3b) so nobody can merge them into one "persistence"
//     case later - a permanent免审通行证 is exactly what this boundary exists to
//     keep out.
//   - It does not enforce anything. Enforcement is internal/tools' route(),
//     which receives the mode as a parameter.
package perm

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ConfigManager is the slice of *config.Manager this package needs. Keeping it
// an interface is what lets the persistence tests run against a real temp
// config.toml and a real reload, rather than against a mock that agrees to
// remember things it never had to store.
type ConfigManager interface {
	// Config returns the effective config (the mode's single source of truth).
	Config() *config.Config
	// SetPermissionMode applies and persists the mode to config.toml.
	SetPermissionMode(mode risk.Mode) error
}

// Result is how one switch attempt ended. It is part of the audit record: "the
// user tried to go全自动 and the card was refused" and "the mode is now auto"
// must be distinguishable from the log alone.
type Result string

// The switch outcomes.
const (
	ResultApplied      Result = "applied"
	ResultRefused      Result = "refused-confirm"
	ResultConfirmError Result = "confirm-failed"
	ResultPersistError Result = "persist-failed"
	ResultNoChange     Result = "no-change"
	ResultBadMode      Result = "invalid-mode"
	// ResultCardNotAnswered is ticket 371's own refusal (AC#2'): the switch
	// needed an answered L2 card, a channel was attached, and the limb came back
	// WITHOUT one - no card was shown, or the one that was shown never came back
	// allow. It is its own value because "refused" alone cannot be told apart from
	// ResultRefused's "this machine has no channel at all", and those are two
	// different operator stories: one is a hole in the assembly, the other is a
	// card nobody answered.
	ResultCardNotAnswered Result = "card-not-answered"
	// ResultDowngraded is the READ side of ticket 372-r2: the configured档 is
	// looser than ask_every_step and this store has no confirmation channel, so
	// PermissionMode answers ask_every_step while config.toml keeps what the
	// operator set. It is its own value because "the operator's file says
	// auto_approve but this process runs strict" must be greppable, not prose.
	ResultDowngraded Result = "downgraded-no-channel"
)

// ErrCardNotAnswered is the error a carded switch earns when the confirmation
// limb reported no answered allow. The wording is load-bearing: AC#2' requires
// the refusal to NAME this, because "refused" is indistinguishable from the
// no-channel branch and would hide the hole ticket 371 exists to close.
var ErrCardNotAnswered = errors.New(
	"perm: card was not answered (no L2 card was shown and answered allow)")

// Switch is the audit record of one mode switch attempt.
type Switch struct {
	From   risk.Mode // the档 in effect before
	To     risk.Mode // the档 asked for
	At     time.Time // when it was asked
	Origin string    // where the request came from: "panel" / "ball" / "cli" / "config-file"
	Actor  string    // who triggered it (login name / session id / "unknown")
	Result Result    // see Result
	Detail string    // the human-readable why, verbatim in the audit line
	Secs   int64     // seconds the confirmation took (0 when none was asked)
	// Asked reports whether this switch's card was actually SHOWN and answered
	// allow (ticket 371-b, the 〔丙-b〕 half of the ruling). Before it, the audit
	// row could not answer "did this switch really put a card in front of
	// somebody?" - the hole docs/evidence/s1/371-a1p-impl.md §5 measured was
	// silent in the log as well as in the gate.
	Asked bool
	// Card is the correlation id of the card that was answered, empty when none
	// was (or when the surface that answered had no id). The approval layer's own
	// `ANSWER-ALLOW corr=<id>` line carries the same value, so the two records can
	// be joined by it; that is what makes "asked, but which card?" answerable.
	Card string
}

// Card is the evidence that one L2 card was actually put in front of an operator
// AND that operator answered ALLOW (ticket 371, AC#1'〔丙-a〕). It is what a
// confirmation limb hands back to the layer that asked, and the only thing that
// settles the one target R20/M4 cards.
//
// THE ZERO VALUE MEANS "NOT ANSWERED". That is the whole point of a return value
// separate from the error: before this ticket the limb reported `err == nil`, and
// that is also exactly what a limb that never raised a card reported, so the gate
// could not tell "the operator answered allow" from "nothing happened" and a
// non-nil but never-asking leg was applied silently (the D/F readings in
// docs/evidence/s1/371-a1p-impl.md).
//
// THE LIMIT, stated where a reader can see it: a credential is minted inside this
// process, so it defeats ACCIDENTAL fake limbs - the shape 371-a1p measured - and
// not a caller that deliberately calls AnsweredCard without asking. That is the
// same trust boundary Options.Confirm already had: perm cannot manufacture an
// ask, it can only be handed one, and the production mint site is the one place
// that has seen approval.Gate answer AnswerAllow (cmd/wisp/run.go's
// confirmModeSwitchVia).
type Card struct {
	answered bool
	corr     string
}

// AnsweredCard mints the credential for the card whose correlation id is corr -
// the id the approval queue issued, i.e. the value its own ANSWER-ALLOW line
// carries. The empty string is allowed for a surface that has no correlation id
// (the audit row then prints card=""), and it still proves the answered allow.
func AnsweredCard(corr string) Card { return Card{answered: true, corr: corr} }

// Answered reports whether this credential proves an answered allow. Only this
// bit settles a carded switch.
func (c Card) Answered() bool { return c.answered }

// Corr is the correlation id of the card that was answered ("" when none).
func (c Card) Corr() string { return c.corr }

// ConfirmFunc asks the operator for one L2 strong confirmation and reports the
// credential: a non-zero Card means "one card was shown and this operator
// answered ALLOW". A ZERO Card with a nil error means "the limb ran and no
// answered allow came back" - Set refuses that with ResultCardNotAnswered
// (ticket 371), because a nil error on its own must never buy the档 that removes
// questions. Any error means the switch did not happen, including "could not
// ask" - a refusal, never an assumed yes.
//
// The only production implementation allowed is the C18 approval route
// (internal/agent/approval's L2 queue, native-click-only, 300s deadline that
// resolves to reject). This package takes it as a function so the boundary is
// visible: perm cannot manufacture a confirmation, it can only be handed one.
type ConfirmFunc func(ctx context.Context, req Switch) (Card, error)

// Options configures a Store. Every zero value is safe: no confirmation source
// means auto_approve can never be reached, no log sink means the record is
// dropped but the switch still cannot be silent to the operator who reads the
// config file back.
type Options struct {
	// Manager is required; without a config to read there is no mode.
	Manager ConfigManager
	// Confirm is the L2 strong confirmation the store asks for the one target
	// that costs it (RequiresCard/JudgeCard). nil = auto_approve is unreachable
	// (fail closed) AND the read side collapses to ask_every_step: a composition
	// with no way to ask must not run a档 that removes questions (ticket 372-r2).
	// A non-nil Confirm is not by itself permission either: Set applies the switch
	// only when the limb hands back an answered-allow Card (ticket 371).
	Confirm ConfirmFunc
	// Logf is the audit sink (the project's existing "[audit]" family).
	Logf func(format string, args ...any)
	// Now is the clock for the audit stamps. nil = time.Now.
	Now func() time.Time
	// History caps the in-memory switch records a read-only view can ask for.
	History int
}

// Store is the owner of the permission mode. Safe for concurrent use.
type Store struct {
	mgr     ConfigManager
	confirm ConfirmFunc
	logf    func(string, ...any)
	now     func() time.Time

	mu   sync.Mutex
	hist []Switch
	cap  int
}

// New composes a Store and records what mode the process started with, so the
// audit trail of a long-running host has a first line to compare against.
func New(o Options) (*Store, error) {
	if o.Manager == nil {
		return nil, fmt.Errorf("perm: New requires a ConfigManager (the mode has no other source)")
	}
	hist := o.History
	if hist <= 0 {
		hist = 64
	}
	s := &Store{
		mgr:     o.Manager,
		confirm: o.Confirm,
		logf:    o.Logf,
		now:     o.Now,
		cap:     hist,
	}
	if s.now == nil {
		s.now = time.Now
	}
	configured := s.configuredMode()
	s.record(Switch{
		From: configured, To: configured, At: s.now(),
		Origin: "startup", Actor: "process", Result: ResultNoChange,
		Detail: "读取到的启动档位（config.toml risk.permission_mode）",
	})
	// Ticket 372-r2: a read-side downgrade is a security-posture change, so it
	// leaves the same kind of record a switch does - one line, typed, naming both
	// ends. It is not a MODE-SWITCH, because config.toml did not move.
	if effective := s.PermissionMode(); effective != configured {
		s.auditAs("MODE-DOWNGRADE", Switch{
			From: configured, To: effective, At: s.now(),
			Origin: "startup", Actor: "process", Result: ResultDowngraded,
			Detail: fmt.Sprintf(
				"本进程没有接入确认通道：config.toml 仍是 %s，有效档位降为 %s（372 读侧降级，配置值不动）",
				configured, effective),
		})
	}
	return s, nil
}

// PermissionMode implements tools.ModeSource and *config.Config's own shape: the
// mode THIS PROCESS enforces, which is config.toml's档 except when this store was
// handed no confirmation channel (ticket 372-r2).
//
// WHY THE READ SIDE. A store with Confirm == nil cannot ask for the one card
// R20/M4 requires, and the two looser档 are exactly the ones that remove a
// question: auto_approve silences every non-red-line call and ask_high_risk
// silences every L1 call (risk.Mode.Screen). An assembly with no way to ask must
// therefore run neither of them, so the effective read collapses to
// ask_every_step. The configured value is NOT rewritten: config.toml keeps what
// the operator chose, and the downgrade disappears the day a channel is attached.
// That is Step-Code's effectiveState shape (m4-survey.md:106) and openchamber's
// displayedPermissionMode reading (m4-survey.md:88), not a new invention.
//
// Fail-closed twice over: a nil store, a nil manager, a nil config or an
// unparsable value all answer "ask_every_step". The one thing this function must
// never do is guess toward silence.
func (s *Store) PermissionMode() risk.Mode {
	if s == nil {
		return risk.DefaultMode()
	}
	configured := s.configuredMode()
	if s.confirm == nil && configured != risk.ModeAskEveryStep {
		return risk.ModeAskEveryStep
	}
	return configured
}

// configuredMode is config.toml's own value, WITHOUT the read-side downgrade.
//
// Set's from/to accounting reads it deliberately. If Set asked PermissionMode
// instead, "configured auto_approve + no channel + request ask_every_step" would
// look like a no-change (the effective read is already ask_every_step), so the
// operator's request to make the file stricter would return nil and persist
// nothing. The audit record also stays the PERSISTED pair, which is what a reader
// comparing it with config.toml can verify.
func (s *Store) configuredMode() risk.Mode {
	if s == nil || s.mgr == nil {
		return risk.DefaultMode()
	}
	c := s.mgr.Config()
	if c == nil {
		return risk.DefaultMode()
	}
	return c.PermissionMode()
}

// Set switches to the档 and persists it (R20/M3: the choice survives a new
// session and a restart; nothing else does).
//
// The confirmation rule is R20/M4 and it is asymmetric on purpose: only the
// switch INTO auto_approve costs an answered L2 strong confirmation, because
// that is the only档 that can remove a question about a call the operator would
// otherwise have been asked about. The audit log, by contrast, records ALL THREE
// 档 switches - a security posture change you cannot reconstruct from the log is
// not an audited one.
//
// The GATE in front of it is the TARGET (ticket 372 reverted 367's direction
// rule): with no confirmation channel wired, a switch INTO auto_approve is
// refused and nothing else is. Ticket 367 refused every widening, which made
// ask_every_step -> ask_high_risk unreachable on a host with no card surface - a
// switch R20/M4 does not card, and no surveyed mainstream harness gates on "wider"
// (docs/evidence/s1/m4-survey.md:96). The missing-channel case is instead made
// safe on the READ side: PermissionMode answers ask_every_step for the two looser
// 档, so a host assembled without a leg asks about everything whatever config.toml
// says. JudgeCard remains the ONE definition of the (target, channel) rule, called
// by both transition points (this function and panel.ModeWriteHandler).
//
// AND THE CARD MUST BE ANSWERED (ticket 371). A channel that is attached is only
// the capability to ask; the switch is applied only when the limb hands back a
// Card that says an allow came back. A limb that returns a nil error without ever
// showing a card is refused with ResultCardNotAnswered, which is a different
// story from ResultRefused and names itself in the audit row.
func (s *Store) Set(ctx context.Context, to risk.Mode, origin, actor string) error {
	if !to.Valid() {
		s.audit(Switch{
			To: to, At: s.stamp(), Origin: origin, Actor: actor,
			Result: ResultBadMode, Detail: "不是三档之一，拒绝切换（R20/M1：没有第四档）",
		})
		return fmt.Errorf("perm: undefined permission mode %d", int(to))
	}
	from := s.configuredMode()
	sw := Switch{From: from, To: to, At: s.stamp(), Origin: origin, Actor: actor}
	if from == to {
		sw.Result, sw.Detail = ResultNoChange, "档位未变化"
		s.audit(sw)
		return nil
	}

	// The gate is the TARGET, not the direction (ticket 372; R20/M4). A store
	// with no confirmation channel cannot reach the one档 that removes the most
	// questions. Judging the target is the owner's 2026-09-21 ruling and what dsh
	// / minimax / openchamber do (m4-survey.md:104-106); judging the direction
	// (ticket 367) also refused ask_every_step -> ask_high_risk, which M4 says
	// costs no card and which the read-side downgrade already makes harmless.
	// JudgeCard is the same function panel.ModeWriteHandler calls (ticket 371-b).
	switch JudgeCard(to, s.confirm != nil) {
	case GateNoChannel:
		sw.Result = ResultRefused
		sw.Detail = fmt.Sprintf(
			"从 %s 切到 %s 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）", from, to)
		s.audit(sw)
		return fmt.Errorf("perm: %s", sw.Detail)
	case GateAsk:
		// R20/M4's asymmetry stands on top of that gate: the auto_approve target
		// costs one ANSWERED L2 card. There is no separate nil test here on
		// purpose - JudgeCard refused every auto_approve switch whose channel was
		// nil, so this call cannot dereference a nil channel. That invariant is
		// what the panic shot pins: change the gate to anything that lets
		// auto_approve past a nil channel and this line panics instead of
		// refusing.
		started := s.now()
		card, err := s.confirm(ctx, sw)
		sw.Secs = int64(s.now().Sub(started).Seconds())
		if err != nil {
			sw.Result, sw.Detail = RefusalResult(err), "L2 强确认未通过："+err.Error()
			s.audit(sw)
			return fmt.Errorf("perm: 切到 auto_approve 被拒绝: %w", err)
		}
		if !card.Answered() {
			// Ticket 371's half of the gate: the channel was called and came back
			// without an answered allow. Nothing was shown, or what was shown was
			// never answered, and a nil error used to be enough. The refusal names
			// that, so it cannot be confused with GateNoChannel's.
			sw.Result = ResultCardNotAnswered
			sw.Detail = "切到 auto_approve 需要一张被回答过的 L2 强确认卡，" +
				"但确认通道没有交回任何“有人答过允许”的凭据（card was not answered）：未切"
			s.audit(sw)
			return fmt.Errorf("perm: 切到 auto_approve 被拒绝: %w", ErrCardNotAnswered)
		}
		sw.Asked, sw.Card = true, card.Corr()
	}

	if err := s.mgr.SetPermissionMode(to); err != nil {
		sw.Result, sw.Detail = ResultPersistError, "写入 config.toml 失败，内存档位保持原值："+err.Error()
		s.audit(sw)
		return err
	}
	sw.Result, sw.Detail = ResultApplied, ""
	s.audit(sw)
	return nil
}

// RefusalResult maps one refusal to the audit's Result vocabulary. It is the ONE
// classification the two transition points share (ticket 371-b): Set uses it for
// the error a limb returned, and panel.ModeWriteHandler uses it to name the
// refusal the write leg made on its behalf - so the same refusal cannot be booked
// as two different things at the two gates.
//
// It separates "the operator answered no" from "we could not ask", and both from
// ticket 371's "a card was asked for and never answered": all refuse the switch,
// but only the last one says the confirmation path produced nothing, and only
// ResultConfirmError says the path itself is broken.
func RefusalResult(err error) Result {
	switch {
	case err == nil:
		return ResultApplied
	case errors.Is(err, ErrCardNotAnswered):
		return ResultCardNotAnswered
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return ResultConfirmError
	}
	return ResultRefused
}

// Snapshot is the read-only view ticket 92's panel renders: the档 in effect and
// the switch records, as copies. It carries no setter, because the panel is a
// display surface and an input that claims authority over permissions is the
// exact shape PLAN.md:1588 forbids.
type Snapshot struct {
	Mode    risk.Mode
	History []Switch
}

// Snapshot returns the current mode plus the switch history, newest last.
func (s *Store) Snapshot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Switch, len(s.hist))
	copy(out, s.hist)
	return Snapshot{Mode: s.PermissionMode(), History: out}
}

// Last returns the most recent switch record, if any.
func (s *Store) Last() (Switch, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.hist) == 0 {
		return Switch{}, false
	}
	return s.hist[len(s.hist)-1], true
}

// record appends one switch to the bounded history (oldest evicted).
func (s *Store) record(sw Switch) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hist = append(s.hist, sw)
	if len(s.hist) > s.cap {
		s.hist = s.hist[len(s.hist)-s.cap:]
	}
}

// audit writes one switch to the bounded history and to the project's audit
// sink. The line names its own kind ("MODE-SWITCH") so it can be grepped out of
// an interleaved log, and it prints the mode NAMES rather than the integers: an
// audit reader must not need this source file to understand what changed.
func (s *Store) audit(sw Switch) { s.auditAs("MODE-SWITCH", sw) }

// auditAs is audit with the kind as a parameter, so the read-side downgrade can
// carry the same fields without calling itself a switch: config.toml did not
// move, the effective read did.
//
// asked / card are ticket 371-b's 〔丙-b〕 reading: every row says whether a card
// was really answered on this attempt, and which one. On a MODE-SWITCH into
// auto_approve, card is the same correlation id the approval layer printed on its
// `ANSWER-ALLOW corr=<id>` line, so "the switch said applied" and "somebody
// clicked allow" can be joined instead of assumed. A downgrade, a no-change and a
// refusal that never asked print asked=0 card="".
func (s *Store) auditAs(kind string, sw Switch) {
	s.record(sw)
	if s.logf == nil {
		return
	}
	asked := 0
	if sw.Asked {
		asked = 1
	}
	s.logf("perm: %s from=%s to=%s at=%s origin=%q actor=%q asked=%d card=%q result=%s detail=%q",
		kind, sw.From, sw.To, sw.At.UTC().Format(time.RFC3339Nano), sw.Origin, sw.Actor,
		asked, sw.Card, sw.Result, sw.Detail)
}

func (s *Store) stamp() time.Time {
	if s == nil || s.now == nil {
		return time.Now()
	}
	return s.now()
}
