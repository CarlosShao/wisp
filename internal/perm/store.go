// Package perm owns the user-facing permission mode: the read the enforcement
// chain consults, the switch that persists, the one confirmation a switch costs,
// and the audit record every switch leaves.
//
// Ticket 90, ruling R20. The shape is three narrow jobs, and the boundaries are
// the point of the package:
//
//   - READ (PermissionMode): what the bridge asks once per tool call. It never
//     mutates, and it fail-closes to the strictest档 when the config is missing
//     or unreadable.
//   - WRITE (Set): the only function in the system that changes the mode at
//     runtime. A switch INTO auto_approve costs one L2 strong confirmation
//     (R20/M4); the other two switches cost nothing but are still audited.
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
)

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
}

// ConfirmFunc asks the operator for one L2 strong confirmation and reports
// whether it was granted. A nil return means "allowed"; any error means the
// switch did not happen, including "could not ask" - which is a refusal, never
// an assumed yes.
//
// The only production implementation allowed is the C18 approval route
// (internal/agent/approval's L2 queue, native-click-only, 300s deadline that
// resolves to reject). This package takes it as a function so the boundary is
// visible: perm cannot manufacture a confirmation, it can only be handed one.
type ConfirmFunc func(ctx context.Context, req Switch) error

// Options configures a Store. Every zero value is safe: no confirmation source
// means auto_approve can never be reached, no log sink means the record is
// dropped but the switch still cannot be silent to the operator who reads the
// config file back.
type Options struct {
	// Manager is required; without a config to read there is no mode.
	Manager ConfigManager
	// Confirm is the L2 strong confirmation for the switch INTO auto_approve.
	// nil = auto_approve is unreachable (fail closed).
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
	s.record(Switch{
		From: s.PermissionMode(), To: s.PermissionMode(), At: s.now(),
		Origin: "startup", Actor: "process", Result: ResultNoChange,
		Detail: "读取到的启动档位（config.toml risk.permission_mode）",
	})
	return s, nil
}

// PermissionMode implements tools.ModeSource and *config.Config's own shape:
// the mode currently in effect.
//
// Fail-closed twice over: a nil manager, a nil config or an unparsable value
// all answer "ask_every_step". The one thing this function must never do is
// guess toward silence.
func (s *Store) PermissionMode() risk.Mode {
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
// switch INTO auto_approve costs an L2 strong confirmation, because that is the
// only档 that can remove a question about a call the operator would otherwise
// have been asked about. The audit log, by contrast, records ALL THREE档
// switches - a security posture change you cannot reconstruct from the log is
// not an audited one.
func (s *Store) Set(ctx context.Context, to risk.Mode, origin, actor string) error {
	if !to.Valid() {
		s.audit(Switch{
			To: to, At: s.stamp(), Origin: origin, Actor: actor,
			Result: ResultBadMode, Detail: "不是三档之一，拒绝切换（R20/M1：没有第四档）",
		})
		return fmt.Errorf("perm: undefined permission mode %d", int(to))
	}
	from := s.PermissionMode()
	sw := Switch{From: from, To: to, At: s.stamp(), Origin: origin, Actor: actor}
	if from == to {
		sw.Result, sw.Detail = ResultNoChange, "档位未变化"
		s.audit(sw)
		return nil
	}

	if to == risk.ModeAutoApprove {
		if s.confirm == nil {
			sw.Result = ResultRefused
			sw.Detail = "切到 auto_approve 需要一次 L2 强确认，但本机没有接入确认通道（fail-closed 未切）"
			s.audit(sw)
			return fmt.Errorf("perm: %s", sw.Detail)
		}
		started := s.now()
		if err := s.confirm(ctx, sw); err != nil {
			sw.Result, sw.Detail = classifyConfirmFailure(err), "L2 强确认未通过："+err.Error()
			sw.Secs = int64(s.now().Sub(started).Seconds())
			s.audit(sw)
			return fmt.Errorf("perm: 切到 auto_approve 被拒绝: %w", err)
		}
		sw.Secs = int64(s.now().Sub(started).Seconds())
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

// classifyConfirmFailure separates "the operator answered no" from "we could not
// ask". Both refuse the switch; only the second one is a fault worth its own
// audit value, because it is the one that says the confirmation path is broken.
func classifyConfirmFailure(err error) Result {
	if err == nil {
		return ResultApplied
	}
	switch err {
	case context.Canceled, context.DeadlineExceeded:
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
func (s *Store) audit(sw Switch) {
	s.record(sw)
	if s.logf == nil {
		return
	}
	s.logf("perm: MODE-SWITCH from=%s to=%s at=%s origin=%q actor=%q result=%s detail=%q",
		sw.From, sw.To, sw.At.UTC().Format(time.RFC3339Nano), sw.Origin, sw.Actor,
		sw.Result, sw.Detail)
}

func (s *Store) stamp() time.Time {
	if s == nil || s.now == nil {
		return time.Now()
	}
	return s.now()
}
