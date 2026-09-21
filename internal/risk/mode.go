package risk

// ---------------------------------------------------------------------------
// User-facing permission modes (ticket 90, ruling R20).
//
// This file is NOT part of the frozen C19 contract block in assessor.go: it
// adds no rule, renumbers nothing and changes no verdict. The assessor still
// produces exactly the Decision it always did; this layer answers a different
// question - "given the operator's chosen mode, must this call still be asked
// about?" - and it may only ever make the answer STRICTER than the mode would
// otherwise allow, never looser than a red line.
//
// The three modes (R20/M1, exactly three - "switch everything off" was
// refused):
//
//	ask_every_step  R20/M2 default: every L1 window and every L2 card is shown
//	ask_high_risk   L1 (reversible write) runs silently; L2 still asks
//	auto_approve    L1 runs silently; L2 runs silently EXCEPT for the red-line
//	                classes below, which no mode can silence
//
// Red lines this layer is forbidden to cross, each with its contract source:
//
//	PLAN.md:1629  irreversible operations must be refused or asked about in
//	              EVERY mode          -> R8, and R1 when the tool itself
//	              declares L2 (its own admission of irreversibility)
//	PLAN.md:1640  a C25 taint escalation is covered by no authorization  -> R4
//	              (Decision.SessionOverrideBlocked)
//	SPEC-06:1588  "allow" only ever comes from a native click; a silenced
//	              verdict is not an allow, it is the ABSENCE of a question -
//	              which is why Deny is never silenced here (R3 tier A)
//	SPEC-06 §4.1  tier A is absolute, tier B needs one L2 per file      -> R3
//	PLAN.md:3143  an L2 never waits forever (300s -> reject, C18); this
//	              layer does not own that clock, it only decides whether a
//	              card exists at all
//
// Silence-immune rule set: R1 (declared L2), R2 (outside the authorized
// allowlist = "writes outside the workspace"), R3 (sensitive path A/B),
// R4 (taint), R5 (network target = "外网"), R8 (irreversibility), R9
// (fail-closed). What auto_approve therefore silences is L1 plus the two
// bulk/ergonomics verdicts R6 and R7 - and nothing else.
// ---------------------------------------------------------------------------

// Mode is the operator-selected permission mode (R20/M1).
type Mode int

// The three modes. ModeAskEveryStep is the zero value on purpose: an
// uninitialized or unrecognized setting lands on the strictest档, not on the
// loosest one. Do not reorder these constants - the zero value is the default.
const (
	// ModeAskEveryStep is R20/M2's default: ask for every L1 and every L2.
	ModeAskEveryStep Mode = iota
	// ModeAskHighRisk silences the L1 pre-execution window only.
	ModeAskHighRisk
	// ModeAutoApprove silences L1 and non-red-line L2 verdicts.
	ModeAutoApprove
)

// The config-vocabulary spellings of the modes (config.toml
// risk.permission_mode, ticket 90). These strings are the persisted form.
const (
	ModeAskEveryStepName  = "ask_every_step"
	ModeAskHighRiskName   = "ask_high_risk"
	ModeAutoApproveName   = "auto_approve"
	defaultModeName       = ModeAskEveryStepName
	modeVocabularyPattern = "ask_every_step|ask_high_risk|auto_approve"
)

// DefaultMode is R20/M2: the strictest档 is the default, so a config that
// never mentions the key behaves like a fresh install.
func DefaultMode() Mode { return ModeAskEveryStep }

// ModeNames returns the persisted spellings in mode order (error messages and
// the panel's read-only picker, ticket 92).
func ModeNames() []string {
	return []string{ModeAskEveryStepName, ModeAskHighRiskName, ModeAutoApproveName}
}

// ParseMode reads a persisted mode. An unknown spelling is an ERROR, never a
// silent fallback: a permission setting that fails to parse must stop the load
// (ticket 83's lying-key rule), because "it fell back to the default" and "it
// was never written" are indistinguishable to the operator afterwards.
func ParseMode(s string) (Mode, error) {
	switch s {
	case "", defaultModeName:
		return ModeAskEveryStep, nil
	case ModeAskHighRiskName:
		return ModeAskHighRisk, nil
	case ModeAutoApproveName:
		return ModeAutoApprove, nil
	default:
		return ModeAskEveryStep, &ModeError{Value: s}
	}
}

// ModeError reports an unparseable permission_mode value. It is a distinct
// type so a config error can name the vocabulary without string matching.
type ModeError struct{ Value string }

func (e *ModeError) Error() string {
	return "permission mode " + quote(e.Value) +
		" is not one of " + modeVocabularyPattern +
		" (risk.permission_mode): refusing to guess"
}

func quote(s string) string { return `"` + s + `"` }

// Valid reports whether m is one of the three defined modes.
func (m Mode) Valid() bool { return m >= ModeAskEveryStep && m <= ModeAutoApprove }

// String renders the persisted spelling, so an audit line and the config file
// agree on one shape.
func (m Mode) String() string {
	switch m {
	case ModeAskHighRisk:
		return ModeAskHighRiskName
	case ModeAutoApprove:
		return ModeAutoApproveName
	case ModeAskEveryStep:
		return ModeAskEveryStepName
	default:
		// An out-of-range Mode (a memory corruption, a bad cast) renders as the
		// default and Screen() fail-closes on the same value: never as "".
		return ModeAskEveryStepName
	}
}

// LoosestOf returns the more permissive of two modes (used by the audit line to
// describe what a switch did without re-deriving the ordering).
func LoosestOf(a, b Mode) Mode {
	if a > b {
		return a
	}
	return b
}

// Silenced is what a mode did to one verdict: whether the question was removed,
// and when it was kept, why.
type Silenced struct {
	// Level is the effective level after the mode was applied.
	Level Level
	// Silenced is true when the caller must NOT ask the user.
	Silenced bool
	// Kept names the red line that refused to silence the verdict ("" when the
	// mode applied as asked). It is audit text: "mode=auto_approve but R8
	// irreversible still asks" must be readable in the log, not inferred.
	Kept string
	// Mode is the mode this verdict was screened under, echoed for the log line.
	Mode Mode
}

// Screen applies the mode to one assessed verdict and reports what the caller
// must do. It never raises a level above the verdict it was given (that would
// be a second judge) and it never lowers a red line (that is the point).
//
// The shape is deliberate: Deny and any unknown level are handled first, so
// every later branch can assume it is looking at L1 or L2 only.
func (m Mode) Screen(d Decision) Silenced {
	// Fail-closed floor: an impossible level is not something a mode may wave
	// through, whatever the mode says.
	if d.Level == Deny {
		return Silenced{Level: Deny, Kept: "拒绝级 verdict（Deny 级不受任何模式影响）", Mode: m}
	}
	if d.Level != L1 && d.Level != L2 {
		return Silenced{Level: d.Level, Kept: "无法识别的风险级，fail-closed 保留询问", Mode: m}
	}
	// R20/M2: the default mode asks about everything, so nothing below it can
	// matter. Checked before the red-line scan because it changes no verdict.
	if m == ModeAskEveryStep || !m.Valid() {
		return Silenced{Level: d.Level, Mode: ModeAskEveryStep}
	}
	// L1: the reversible-write pre-execution window. Both relaxed modes run it
	// silently (R20/M1) - no red line is expressed as an L1 verdict, because
	// every red-line rule in SPEC-06 §3 concludes L2 or Deny.
	if d.Level == L1 {
		return Silenced{Level: L0, Silenced: true, Mode: m}
	}
	// L2 and the strictest relaxed mode: high risk still asks.
	if m == ModeAskHighRisk {
		return Silenced{Level: L2, Mode: m}
	}
	// auto_approve, L2. The red lines are the only thing left standing.
	if why, immune := redLine(d); immune {
		return Silenced{Level: L2, Kept: why, Mode: m}
	}
	return Silenced{Level: L0, Silenced: true, Mode: m}
}

// redLine reports whether one L2 verdict carries a class no mode may silence,
// and names it for the audit line.
//
// It is written as an EXCLUSION list on purpose: a new rule added to the
// frozen judge later is only silenced by auto_approve if somebody adds it here
// deliberately. The default for an unknown rule is "still ask".
func redLine(d Decision) (string, bool) {
	if d.SessionOverrideBlocked {
		return "R4 污染升级（C25，任何授权都不可覆盖，含权限模式）", true
	}
	for _, r := range d.RulesHit {
		switch r {
		case R1:
			return "R1 工具自身声明的 L2（不可逆自认，PLAN.md:1629）", true
		case R2:
			return "R2 授权目录之外（写工作区外）", true
		case R3:
			return "R3 敏感路径 A/B 档（SPEC-06 §4.1）", true
		case R4:
			return "R4 污染升级（C25）", true
		case R5:
			return "R5 网络目标（外网）", true
		case R8:
			return "R8 不可逆操作（任何档都必须拒/问）", true
		case R9:
			return "R9 fail-closed（判定器异常）", true
		}
	}
	return "", false
}
