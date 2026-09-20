package risk

import (
	"fmt"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// FROZEN CONTRACT — C19 rule set (SPEC-06 §3). Rule IDs and semantics below
// are contract-level: adding, removing or renumbering a rule is a contract
// change and requires human approval. The confirmation card (tickets 21/37)
// renders Decision.RulesHit and Decision.Reason verbatim.
//
//	R1  declared RiskLevel            lower bound, never the conclusion
//	R2  target path outside allowlist L2 (C26-canonicalized; resolver wired by 18)
//	R3  sensitive path tier A / B     A -> Deny; B -> L2 (per-file exemption later)
//	R4  C25 taint hit                 L2, never overridable by session auth
//	R5  network target                deny-listed protocol -> Deny; private
//	                                  range / non-web protocol / domain outside
//	                                  allowlist / URL > 2048 chars -> L2
//	R6  shell argv metacharacters     metachars or unlisted program -> L2;
//	                                  allowlisted clean argv -> L1
//	R7  batch scale >= 50 files       L2
//	R8  irreversibility               delete/overwrite/power/close-window/send -> L2
//	R9  judge panic / invalid input   fail-closed L2
//
// Absent dependencies (PathCanonicalizer / SensitiveClassifier / TaintDetector
// not yet wired by tickets 18/19) leave the corresponding rule DORMANT — an
// integration gap, not a judged result. A WIRED dependency that panics or
// errors is a judged failure: panic -> R9 fail-closed L2 (recovered here);
// canonicalization error -> R2 fail-closed L2.
//
// Fusion = max severity over all contributions (built-in rules + registered
// sub-assessors + the R1 declared floor); a single Deny beats any allow.
// ---------------------------------------------------------------------------

// Level is the risk gate level (SPEC-06 §2). Severity order is
// L0 < L1 < L2 < Deny: fusion takes the max, and Deny is terminal.
type Level int

const (
	// L0: read-only, no side effects — execute directly.
	L0 Level = iota
	// L1: reversible write — pre-execution block window (2-3s).
	L1
	// L2: irreversible / high risk — C18 ApprovalQueue.
	L2
	// Deny: rejected outright (tier-A sensitive path, deny-listed protocol).
	Deny
)

// String renders the level the way the confirmation card shows it.
func (l Level) String() string {
	switch l {
	case L0:
		return "L0"
	case L1:
		return "L1"
	case L2:
		return "L2"
	case Deny:
		return "Deny"
	default:
		return fmt.Sprintf("Level(%d)", int(l))
	}
}

// RuleID identifies a frozen C19 rule (R1-R9). Encoded as data so the rule set
// stays stable, reviewable and diffable.
type RuleID string

// The frozen rule IDs (SPEC-06 §3). Do not renumber; do not add without human
// approval — this is a contract change.
const (
	R1 RuleID = "R1" // declared level = lower bound
	R2 RuleID = "R2" // path outside authorized allowlist
	R3 RuleID = "R3" // sensitive path tier A (Deny) / tier B (L2)
	R4 RuleID = "R4" // C25 taint hit (exfil channel carries sensitive content)
	R5 RuleID = "R5" // network target (protocol / private range / allowlist)
	R6 RuleID = "R6" // shell argv metacharacters / allowlist
	R7 RuleID = "R7" // batch scale >= 50 files
	R8 RuleID = "R8" // irreversibility
	R9 RuleID = "R9" // fail-closed (panic, invalid declared level)
)

// Decision is the ONLY risk verdict the gate hands out. rules_hit and reason
// are rendered verbatim on the confirmation card, so reasons are human-readable
// and stable.
type Decision struct {
	Level    Level
	RulesHit []RuleID
	Reason   string
	// SessionOverrideBlocked is true when the verdict must not be covered by a
	// D45 scope session authorization (R4 taint hits, SPEC-06 §8.3).
	SessionOverrideBlocked bool
}

// Facts is the per-call context the rules judge. It is an INPUT bundle: the
// tool provider fills it before calling Assess. All fields are optional; a
// rule whose inputs are absent does not fire.
type Facts struct {
	// Declared is the tool manifest's declared RiskLevel (R1: the lower bound).
	Declared Level
	// Paths are the raw affected target paths (R2 allowlist, R3 sensitive tiers,
	// R7 fallback count). Unresolved — the wired PathCanonicalizer owns C26.
	Paths []string
	// BatchCount is how many files this single call affects (R7). Zero falls
	// back to len(Paths).
	BatchCount int
	// Network describes the outbound target (R5). Nil = not a network call.
	Network *NetTarget
	// ShellArgv is the forced argv vector (R6). Nil/empty = not a shell call.
	ShellArgv []string
	// ShellString marks string-mode shell ([risk] allow_shell_string=true):
	// always L2 (SPEC-06 §10).
	ShellString bool
	// ShellAllowlist lists allowlisted shell programs (R6: clean argv -> L1).
	ShellAllowlist []string
	// Irreversible lists irreversibility classes (R8): "delete", "overwrite",
	// "power", "close-window", "send". Unknown classes fail closed to L2.
	Irreversible []string
	// OverwriteExisting marks that an existing file will be overwritten (R8).
	OverwriteExisting bool
}

// NetTarget is the R5 input for one outbound call.
type NetTarget struct {
	Protocol  string   // scheme: "http", "https", ... empty = unknown
	Host      string   // bare host or IP literal (port optional, stripped)
	URL       string   // full URL when the call has one (length cap 2048)
	Allowlist []string // allowlisted domains (config-fed); empty = no domain gate
}

// PathCanonicalizer is the consumer-side interface for C26 (ticket 18 wires
// the real PathResolver). Canonicalize must be the ONLY canonicalization
// entry (SPEC-06 §4); InAllowlist judges canonical paths against the
// authorized-directory allowlist.
type PathCanonicalizer interface {
	Canonicalize(raw string) (string, error)
	InAllowlist(canonical string) bool
}

// Tier classifies a canonical path against the sensitive-path blacklist
// (SPEC-06 §4.1): tier A is absolute deny, tier B is L2 with per-file
// exemption (exemption flow lands with the approval queue, ticket 21).
type Tier int

const (
	TierNone Tier = iota
	TierB         // default-deny + single-file exemption -> L2
	TierA         // absolute deny, not overridable by any authorization
)

// SensitiveClassifier is the consumer-side interface for the A/B sensitive
// path lists (ticket 18 wires the real blacklist). Input is a C26-canonical
// path.
type SensitiveClassifier interface {
	Classify(canonical string) Tier
}

// TaintDetector is the consumer-side interface for C25 provenance (ticket 19
// wires the real taint engine). It inspects the outgoing call parameters for
// sensitive-source fragments on an exfil channel.
type TaintDetector interface {
	TaintHit(params map[string]any) (source string, hit bool)
}

// SubAssessor is the pluggable assessor interface. RESERVED for guardian
// scoring (attaches later) — do NOT implement it in this ticket. Any panic
// inside a registered sub-assessor is recovered and fail-closes to R9 L2.
type SubAssessor interface {
	Assess(tool string, params map[string]any, facts Facts) Decision
}

// RiskAssessor is the C19 central risk assessor — the ONLY authority for risk
// decisions in the system (SPEC-06 §3). Zero values are unusable; use New.
type RiskAssessor struct {
	canonicalizer PathCanonicalizer   // wired by ticket 18
	classifier    SensitiveClassifier // wired by ticket 18
	taint         TaintDetector       // wired by ticket 19
	subassessors  []SubAssessor
}

// NewRiskAssessor builds an assessor with the built-in R1-R9 rule set. Until
// tickets 18/19 wire their dependencies, R2/R3/R4 stay dormant (documented
// integration gap, see the frozen-contract block above).
func NewRiskAssessor() *RiskAssessor {
	return &RiskAssessor{}
}

// WithCanonicalizer wires the C26 resolver (ticket 18). Returns the receiver.
func (a *RiskAssessor) WithCanonicalizer(p PathCanonicalizer) *RiskAssessor {
	a.canonicalizer = p
	return a
}

// WithSensitiveClassifier wires the A/B blacklist (ticket 18). Returns the
// receiver.
func (a *RiskAssessor) WithSensitiveClassifier(c SensitiveClassifier) *RiskAssessor {
	a.classifier = c
	return a
}

// WithTaintDetector wires the C25 taint engine (ticket 19). Returns the
// receiver.
func (a *RiskAssessor) WithTaintDetector(t TaintDetector) *RiskAssessor {
	a.taint = t
	return a
}

// RegisterSubAssessor attaches a pluggable judge (guardian scoring later). Its
// verdict fuses at max severity; its panic fail-closes to R9 L2.
func (a *RiskAssessor) RegisterSubAssessor(s SubAssessor) {
	a.subassessors = append(a.subassessors, s)
}

// contribution is one judged input to fusion.
type contribution struct {
	rules            []RuleID
	level            Level
	reason           string
	blockSessionAuth bool
}

// Assess classifies one tool call. It never panics and never returns a level
// below the declared floor; any judge panic fail-closes to R9 L2.
func (a *RiskAssessor) Assess(tool string, params map[string]any, facts Facts) Decision {
	ctx := &assessCtx{
		tool:       tool,
		params:     params,
		facts:      &facts,
		canon:      a.canonicalizer,
		classifier: a.classifier,
		taint:      a.taint,
	}
	var contribs []*contribution
	for _, rf := range builtinRules {
		if c := runBuiltin(rf, ctx); c != nil {
			contribs = append(contribs, c)
		}
	}
	for _, sub := range a.subassessors {
		if c := runSubAssessor(sub, tool, params, facts); c != nil {
			contribs = append(contribs, c)
		}
	}
	return fuse(contribs)
}

// assessCtx carries what the built-in rules need.
type assessCtx struct {
	tool       string
	params     map[string]any
	facts      *Facts
	canon      PathCanonicalizer
	classifier SensitiveClassifier
	taint      TaintDetector
}

// builtinRules is the frozen judge pipeline, in R1-R8 order (R9 is the
// fail-closed wrapper, not a pipeline stage).
var builtinRules = []func(*assessCtx) *contribution{
	ruleDeclared,
	rulePathAllowlist,
	ruleSensitivePath,
	ruleTaint,
	ruleNetwork,
	ruleShell,
	ruleBatchScale,
	ruleIrreversible,
}

// runBuiltin executes one built-in rule under recover: a panicking judge is a
// judged failure and fail-closes to R9 L2.
func runBuiltin(rf func(*assessCtx) *contribution, ctx *assessCtx) (c *contribution) {
	defer func() {
		if rec := recover(); rec != nil {
			c = &contribution{
				rules:  []RuleID{R9},
				level:  L2,
				reason: fmt.Sprintf("R9: 判定器 panic，fail-closed 升 L2（内置规则 panic: %v）", rec),
			}
		}
	}()
	return rf(ctx)
}

// runSubAssessor executes one registered sub-assessor under recover.
func runSubAssessor(s SubAssessor, tool string, params map[string]any, facts Facts) (c *contribution) {
	defer func() {
		if rec := recover(); rec != nil {
			c = &contribution{
				rules:  []RuleID{R9},
				level:  L2,
				reason: fmt.Sprintf("R9: 判定器 panic，fail-closed 升 L2（插件判定器 panic: %v）", rec),
			}
		}
	}()
	d := s.Assess(tool, params, facts)
	return &contribution{
		rules:            d.RulesHit,
		level:            d.Level,
		reason:           d.Reason,
		blockSessionAuth: d.SessionOverrideBlocked,
	}
}

// fuse merges contributions at max severity. rules_hit is deduped and sorted
// so the confirmation card renders a stable, reviewable list; the reason is
// the reasons of every hit rule joined with "; " (empty when nothing fired).
func fuse(contribs []*contribution) Decision {
	level := L0
	reasonByRule := map[RuleID]string{}
	blocked := false
	for _, c := range contribs {
		if c == nil {
			continue
		}
		if c.level > level {
			level = c.level
		}
		if c.blockSessionAuth {
			blocked = true
		}
		for _, r := range c.rules {
			if _, seen := reasonByRule[r]; !seen {
				reasonByRule[r] = c.reason
			}
		}
	}
	rules := make([]RuleID, 0, len(reasonByRule))
	for r := range reasonByRule {
		rules = append(rules, r)
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i] < rules[j] })
	parts := make([]string, 0, len(rules))
	for _, r := range rules {
		if reasonByRule[r] != "" {
			parts = append(parts, reasonByRule[r])
		}
	}
	reason := strings.Join(parts, "; ")
	if len(rules) == 0 {
		reason = "无规则命中（L0 直接执行）"
	}
	return Decision{
		Level:                  level,
		RulesHit:               rules,
		Reason:                 reason,
		SessionOverrideBlocked: blocked,
	}
}
