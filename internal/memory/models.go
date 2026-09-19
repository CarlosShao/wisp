package memory

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Row models for the eight D35 tables (SPEC-02 §3). Field-for-field the
// columns; timestamps are unix seconds (wall clock, D42#9). Enum-shaped
// fields are validated at the DAO boundary — the SQL DDL stays free of CHECK
// constraints so it remains byte-identical to the spec.

// Profile is one L1 user-profile slot (D20). Slots are a finite enum
// ('pref.language', 'habit.work_hours', ...); at most L1MaxProfiles rows
// exist, enforced app-layer with logged LRU eviction (D35).
type Profile struct {
	ID        int64
	Slot      string
	Content   string
	Source    string // 'extracted' | 'manual'
	UpdatedAt int64  // unix seconds
}

// L1MaxProfiles is the D35 profile row cap (D20 L1, [memory] l1_max default).
const L1MaxProfiles = 20

// ProfileSource enumerates the profile.source column values.
type ProfileSource string

const (
	ProfileExtracted ProfileSource = "extracted"
	ProfileManual    ProfileSource = "manual"
)

// Valid reports whether s is one of the two allowed sources.
func (s ProfileSource) Valid() bool { return s == ProfileExtracted || s == ProfileManual }

// Memory is one L2 explicit memory row (D20). Kind is reserved at 'explicit'
// per the spec comment (no additional kinds are implemented).
type Memory struct {
	ID        int64
	Kind      string
	Content   string
	Keywords  string // space-separated, lowercase-normalized
	CreatedAt int64
	LastHitAt int64
	HitCount  int64
}

// MemoryKindExplicit is the only kind value (SPEC-02 §3 memory.kind).
const MemoryKindExplicit = "explicit"

// TaskLog is one L3 task row (D20). QueryText is expected pre-redacted per
// §14.4 — this layer stores what it is given and never claims the mask is
// complete.
type TaskLog struct {
	ID              string // task UUID
	StartedAt       int64  // unix seconds
	EndedAt         *int64 // nil while running
	State           string // terminal state (incl. 'interrupted')
	QueryText       string // redacted upstream
	SummaryText     *string
	CostTokensIn    int64
	CostTokensOut   int64
	CostAmountMicro int64  // 1e-6 currency units (no floats)
	Currency        string // default 'CNY'
	ErrorClass      string // D37 enum or "" (NULL)
}

// ToolCall is one forensics row: a tool invocation plus its approval decision
// (D35). This is what answers "what did it actually do to my machine".
type ToolCall struct {
	ID            int64
	TaskID        string
	Seq           int64
	Tool          string
	ArgsJSON      string
	RiskLevel     string // 'L0'|'L1'|'L2'
	Decision      string // '' (pending) | allow|allow_session_grant|reject|timeout|batch_aggregated
	DecidedAt     *int64
	StartedAt     *int64
	EndedAt       *int64
	Outcome       string // '' (pending) | success|error|cancelled|truncated
	ErrorClass    string // D37 enum or ""
	CorrelationID string // C18
	GrantID       *int64
}

// ApprovalGrant is one D45 scoped session grant. Expired/revoked rows stay
// for audit for GrantAuditTTL (30 days) before RetentionJob removes them.
type ApprovalGrant struct {
	ID        int64
	Scope     string // 'session'
	Tool      string
	Pattern   string
	SessionID string
	CreatedAt int64
	ExpiresAt int64
	RevokedAt *int64
}

// GrantScopeSession is the only scope value (SPEC-02 §3 approval_grant.scope).
const GrantScopeSession = "session"

// CostDaily is the C23 per-day aggregation row.
type CostDaily struct {
	Day        string // 'YYYY-MM-DD'
	TokensIn   int64
	TokensOut  int64
	CostMicros int64
	Tasks      int64
}

// PluginState is the §14.7 install record incl. integrity hashes.
type PluginState struct {
	ID               string
	Version          string
	Enabled          bool
	CapabilitiesJSON string
	NetAllowlistJSON string // "" = NULL
	InstalledAt      int64
	Hash             string // manifest hash; mismatch -> refuse to load (upstream)
	ExeHash          string // D46 command plugin exe sha256; "" = NULL
}

// Risk levels (tool_call.risk_level).
const (
	RiskL0 = "L0"
	RiskL1 = "L1"
	RiskL2 = "L2"
)

// toolCallDecisions are the tool_call.decision enum values (DDL comment).
var toolCallDecisions = map[string]bool{
	"allow":               true,
	"allow_session_grant": true,
	"reject":              true,
	"timeout":             true,
	"batch_aggregated":    true,
}

// toolCallOutcomes are the tool_call.outcome enum values (DDL comment).
var toolCallOutcomes = map[string]bool{
	"success":   true,
	"error":     true,
	"cancelled": true,
	"truncated": true,
}

func validRiskLevel(s string) error {
	switch s {
	case RiskL0, RiskL1, RiskL2:
		return nil
	}
	return fmt.Errorf("memory: invalid risk_level %q (want L0|L1|L2)", s)
}

func validateToolCall(tc ToolCall) error {
	if tc.TaskID == "" {
		return fmt.Errorf("memory: tool_call.task_id is required")
	}
	if tc.Tool == "" {
		return fmt.Errorf("memory: tool_call.tool is required")
	}
	if tc.CorrelationID == "" {
		return fmt.Errorf("memory: tool_call.correlation_id is required (C18)")
	}
	if err := validRiskLevel(tc.RiskLevel); err != nil {
		return err
	}
	if tc.Decision != "" && !toolCallDecisions[tc.Decision] {
		return fmt.Errorf("memory: invalid tool_call.decision %q", tc.Decision)
	}
	if tc.Outcome != "" && !toolCallOutcomes[tc.Outcome] {
		return fmt.Errorf("memory: invalid tool_call.outcome %q", tc.Outcome)
	}
	if tc.ErrorClass != "" {
		return observe.ValidateErrorClass(tc.ErrorClass)
	}
	return nil
}

func validateTaskLog(tl TaskLog) error {
	if tl.ID == "" {
		return fmt.Errorf("memory: task_log.id is required")
	}
	if tl.State == "" {
		return fmt.Errorf("memory: task_log.state is required")
	}
	if tl.ErrorClass != "" {
		// D37: task_log.error_class accepts exactly the 17 observe classes.
		return observe.ValidateErrorClass(tl.ErrorClass)
	}
	return nil
}

// normalizeKeywords lowercases, splits on whitespace, dedupes and re-joins
// with single spaces (SPEC-02 §3 memory.keywords: space-separated,
// lowercase-normalized).
func normalizeKeywords(s string) string {
	fields := strings.Fields(strings.ToLower(s))
	seen := make(map[string]bool, len(fields))
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	return strings.Join(out, " ")
}

// escapeLike escapes LIKE wildcards for use with the ESCAPE '\' clause.
func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}
