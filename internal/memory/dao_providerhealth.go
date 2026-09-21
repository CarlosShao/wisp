package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// provider_health DAO primitive (SPEC-02 §3 note; ticket 09 defines the
// primitive, ticket 11 writes real probe results, ticket 40 renders it).
// Storage-split rule: runtime observation state lives HERE and never in
// config.toml; the catalog (truth source) never learns about health.

// Quota states of provider_health.quota_state.
const (
	QuotaOK        = "ok"
	QuotaThrottled = "throttled"
	QuotaExhausted = "exhausted"
	QuotaUnknown   = "unknown"
)

// ProbeFlags is the measured capability bit-set (the JSON payload of
// provider_health.probe_json). Keys mirror config.Capabilities but the
// source of truth is the PROBE, not the user's declaration.
type ProbeFlags struct {
	Text     *bool `json:"text,omitempty"`
	Vision   *bool `json:"vision,omitempty"`
	AudioIn  *bool `json:"audio_in,omitempty"`
	AudioOut *bool `json:"audio_out,omitempty"`
	Thinking *bool `json:"thinking,omitempty"`
	FC       *bool `json:"fc,omitempty"`
}

// Set returns a copy of f with the given capability set to v.
func (f ProbeFlags) Set(cap string, v bool) ProbeFlags {
	switch cap {
	case "text":
		f.Text = &v
	case "vision":
		f.Vision = &v
	case "audio_in":
		f.AudioIn = &v
	case "audio_out":
		f.AudioOut = &v
	case "thinking":
		f.Thinking = &v
	case "fc":
		f.FC = &v
	}
	return f
}

// Get returns the flag for a capability name (nil = never probed).
func (f ProbeFlags) Get(cap string) *bool {
	switch cap {
	case "text":
		return f.Text
	case "vision":
		return f.Vision
	case "audio_in":
		return f.AudioIn
	case "audio_out":
		return f.AudioOut
	case "thinking":
		return f.Thinking
	case "fc":
		return f.FC
	}
	return nil
}

// ProviderHealth is one row of provider_health.
type ProviderHealth struct {
	Provider     string
	Model        string
	Probe        ProbeFlags // decoded probe_json
	LastProbeAt  time.Time  // zero = never probed
	LastError    string     // redacted detail; empty when none
	LastErrorAt  time.Time
	LatencyMSP50 int64
	QuotaState   string // ok | throttled | exhausted | unknown
}

// validQuotaStates guards the quota_state vocabulary at the DAO layer.
var validQuotaStates = map[string]bool{
	QuotaOK: true, QuotaThrottled: true, QuotaExhausted: true, QuotaUnknown: true,
}

// UpsertProviderProbe writes a probe outcome for provider/model (ticket 11's
// write path). lastProbeAt is the wall-clock instant of the probe.
func (s *Store) UpsertProviderProbe(ctx context.Context, provider, model string,
	flags ProbeFlags, ok bool, latencyMS int64, lastProbeAt time.Time,
) error {
	if provider == "" || model == "" {
		return errors.New("memory: provider and model are required")
	}
	probeJSON, err := json.Marshal(flags)
	if err != nil {
		return fmt.Errorf("memory: encode probe flags: %w", err)
	}
	probe := string(probeJSON)
	if probe == "" || probe == "null" {
		probe = "{}"
	}
	quota := QuotaOK
	if !ok {
		quota = QuotaUnknown
	}
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
INSERT INTO provider_health(provider, model, probe_json, last_probe_at, latency_ms_p50, quota_state)
VALUES(?, ?, ?, ?, ?, ?)
ON CONFLICT(provider, model) DO UPDATE SET
  probe_json     = excluded.probe_json,
  last_probe_at  = excluded.last_probe_at,
  latency_ms_p50 = excluded.latency_ms_p50,
  quota_state    = excluded.quota_state`,
			provider, model, probe, lastProbeAt.UTC().Unix(), latencyMS, quota)
		return err
	})
}

// RecordProviderError stores the latest classified provider failure (D37
// class + redacted detail). quota_state, when non-empty, is updated in the
// same write (e.g. exhausted after a budget-class failure).
func (s *Store) RecordProviderError(ctx context.Context, provider, model, class, detail, quotaState string) error {
	if provider == "" || model == "" {
		return errors.New("memory: provider and model are required")
	}
	if class != "" {
		if err := observe.ValidateErrorClass(class); err != nil {
			return err
		}
	}
	if quotaState != "" && !validQuotaStates[quotaState] {
		return fmt.Errorf("memory: invalid quota_state %q", quotaState)
	}
	detail = truncate(detail, 2000)
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		// Ensure the row exists, then stamp the error (upsert-then-update
		// keeps a probe result intact while recording the failure).
		if _, err := tx.ExecContext(ctx, `
INSERT INTO provider_health(provider, model) VALUES(?, ?)
ON CONFLICT(provider, model) DO NOTHING`, provider, model); err != nil {
			return err
		}
		if quotaState == "" {
			_, err := tx.ExecContext(ctx, `
UPDATE provider_health SET last_error = ?, last_error_at = ? WHERE provider = ? AND model = ?`,
				detail, time.Now().UTC().Unix(), provider, model)
			return err
		}
		_, err := tx.ExecContext(ctx, `
UPDATE provider_health SET last_error = ?, last_error_at = ?, quota_state = ? WHERE provider = ? AND model = ?`,
			detail, time.Now().UTC().Unix(), quotaState, provider, model)
		return err
	})
}

// ProviderHealthRow reads one row; sql.ErrNoRows propagates when the pair
// was never probed.
func (s *Store) ProviderHealthRow(ctx context.Context, provider, model string) (ProviderHealth, error) {
	row := s.reader.QueryRowContext(ctx, `
SELECT provider, model, probe_json, last_probe_at, last_error, last_error_at, latency_ms_p50, quota_state
FROM provider_health WHERE provider = ? AND model = ?`, provider, model)
	return scanProviderHealth(row)
}

// ListProviderHealth returns all rows ordered by provider then model.
func (s *Store) ListProviderHealth(ctx context.Context) ([]ProviderHealth, error) {
	rows, err := s.reader.QueryContext(ctx, `
SELECT provider, model, probe_json, last_probe_at, last_error, last_error_at, latency_ms_p50, quota_state
FROM provider_health ORDER BY provider, model`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProviderHealth
	for rows.Next() {
		h, err := scanProviderHealth(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

// scanner is the row/rows subset both readers need.
type scanner interface{ Scan(dest ...any) error }

func scanProviderHealth(sc scanner) (ProviderHealth, error) {
	var h ProviderHealth
	var probeJSON string
	var lastProbeAt, lastErrorAt sql.NullInt64
	var lastError sql.NullString
	if err := sc.Scan(&h.Provider, &h.Model, &probeJSON, &lastProbeAt,
		&lastError, &lastErrorAt, &h.LatencyMSP50, &h.QuotaState); err != nil {
		return ProviderHealth{}, err
	}
	if probeJSON != "" && probeJSON != "null" {
		if err := json.Unmarshal([]byte(probeJSON), &h.Probe); err != nil {
			return ProviderHealth{}, fmt.Errorf("memory: decode probe_json for %s/%s: %w",
				h.Provider, h.Model, err)
		}
	}
	if lastProbeAt.Valid {
		h.LastProbeAt = time.Unix(lastProbeAt.Int64, 0).UTC()
	}
	h.LastError = lastError.String
	if lastErrorAt.Valid {
		h.LastErrorAt = time.Unix(lastErrorAt.Int64, 0).UTC()
	}
	return h, nil
}

// truncate bounds a stored detail string (rune-safe).
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
