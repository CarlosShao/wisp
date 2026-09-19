package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Privacy operations API (D35 rule 2, SPEC-02 §5): for profile / memory /
// task_log / tool_call / artifacts — list, delete-one, purge-all, export
// JSON. The panel "data & privacy" page (ticket 40) only RENDERS this layer;
// nothing privacy-relevant may bypass it.

// PrivacyDomain enumerates the five privacy-relevant data domains.
type PrivacyDomain string

const (
	PrivacyProfile   PrivacyDomain = "profile"
	PrivacyMemory    PrivacyDomain = "memory"
	PrivacyTaskLog   PrivacyDomain = "task_log"
	PrivacyToolCall  PrivacyDomain = "tool_call"
	PrivacyArtifacts PrivacyDomain = "artifacts"
)

// PrivacyDomains returns the domains in SPEC-02 §5 order.
func PrivacyDomains() []PrivacyDomain {
	return []PrivacyDomain{
		PrivacyProfile, PrivacyMemory, PrivacyTaskLog, PrivacyToolCall, PrivacyArtifacts,
	}
}

// PrivacyItem is one row as the privacy page lists it: a stable ID (row id
// or artifact file name), a human label, and the record timestamp.
type PrivacyItem struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Detail    string `json:"detail,omitempty"`
	Timestamp int64  `json:"timestamp"` // unix seconds wall clock; 0 = n/a
	SizeBytes int64  `json:"size_bytes,omitempty"`
}

// ListPrivacy lists one domain's items (oldest-first for artifacts,
// newest-first for DB rows — matching the respective DAO listings).
func (s *Store) ListPrivacy(ctx context.Context, domain PrivacyDomain) ([]PrivacyItem, error) {
	switch domain {
	case PrivacyProfile:
		rows, err := s.ListProfiles(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]PrivacyItem, 0, len(rows))
		for _, r := range rows {
			out = append(out, PrivacyItem{
				ID: strconv.FormatInt(r.ID, 10), Label: r.Slot, Detail: r.Content,
				Timestamp: r.UpdatedAt,
			})
		}
		return out, nil
	case PrivacyMemory:
		rows, err := s.ListMemories(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]PrivacyItem, 0, len(rows))
		for _, r := range rows {
			out = append(out, PrivacyItem{
				ID: strconv.FormatInt(r.ID, 10), Label: firstLine(r.Content), Detail: r.Keywords,
				Timestamp: r.CreatedAt,
			})
		}
		return out, nil
	case PrivacyTaskLog:
		rows, err := s.ListTaskLogs(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]PrivacyItem, 0, len(rows))
		for _, r := range rows {
			detail := r.State
			if r.SummaryText != nil && *r.SummaryText != "" {
				detail = r.State + " · " + *r.SummaryText
			}
			out = append(out, PrivacyItem{
				ID: r.ID, Label: firstLine(r.QueryText), Detail: detail,
				Timestamp: r.StartedAt,
			})
		}
		return out, nil
	case PrivacyToolCall:
		rows, err := s.ListToolCalls(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]PrivacyItem, 0, len(rows))
		for _, r := range rows {
			detail := fmt.Sprintf("%s · %s", r.RiskLevel, r.Tool)
			if r.Decision != "" {
				detail += " · " + r.Decision
			}
			ts := int64(0)
			if r.StartedAt != nil {
				ts = *r.StartedAt
			} else if r.DecidedAt != nil {
				ts = *r.DecidedAt
			}
			out = append(out, PrivacyItem{
				ID: strconv.FormatInt(r.ID, 10), Label: r.Tool, Detail: detail,
				Timestamp: ts,
			})
		}
		return out, nil
	case PrivacyArtifacts:
		rows, err := s.ListArtifacts(ctx)
		if err != nil {
			return nil, err
		}
		out := make([]PrivacyItem, 0, len(rows))
		for _, r := range rows {
			out = append(out, PrivacyItem{
				ID: r.Name, Label: r.Name,
				Timestamp: r.ModTime.Unix(), SizeBytes: r.SizeBytes,
			})
		}
		return out, nil
	default:
		return nil, fmt.Errorf("memory: unknown privacy domain %q (valid: %v)", domain, PrivacyDomains())
	}
}

// DeletePrivacyItem removes one item: DB row id (profile/memory/tool_call)
// or task UUID (task_log) or artifact file name.
func (s *Store) DeletePrivacyItem(ctx context.Context, domain PrivacyDomain, id string) error {
	if id == "" {
		return fmt.Errorf("memory: privacy item id is required")
	}
	switch domain {
	case PrivacyProfile:
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("memory: profile item id %q is not a row id", id)
		}
		return s.DeleteProfile(ctx, n)
	case PrivacyMemory:
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("memory: memory item id %q is not a row id", id)
		}
		return s.DeleteMemory(ctx, n)
	case PrivacyTaskLog:
		return s.DeleteTaskLog(ctx, id)
	case PrivacyToolCall:
		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("memory: tool_call item id %q is not a row id", id)
		}
		return s.DeleteToolCall(ctx, n)
	case PrivacyArtifacts:
		return s.DeleteArtifact(ctx, id)
	default:
		return fmt.Errorf("memory: unknown privacy domain %q (valid: %v)", domain, PrivacyDomains())
	}
}

// PurgePrivacy clears a whole domain and returns how many items were removed
// (privacy: one-click clear, must be countable so the UI can show what
// happened).
func (s *Store) PurgePrivacy(ctx context.Context, domain PrivacyDomain) (int64, error) {
	switch domain {
	case PrivacyProfile:
		return s.PurgeProfiles(ctx)
	case PrivacyMemory:
		return s.PurgeMemories(ctx)
	case PrivacyTaskLog:
		return s.PurgeTaskLogs(ctx)
	case PrivacyToolCall:
		return s.PurgeToolCalls(ctx)
	case PrivacyArtifacts:
		n, err := s.PurgeArtifacts(ctx)
		return int64(n), err
	default:
		return 0, fmt.Errorf("memory: unknown privacy domain %q (valid: %v)", domain, PrivacyDomains())
	}
}

// privacyExport is the JSON export envelope (SPEC-02 §5: 导出). Records are
// the full domain rows — the export is the user's data, not a summary.
type privacyExport struct {
	Domain     string `json:"domain"`
	ExportedAt string `json:"exported_at"` // RFC3339 UTC wall clock (persisted record)
	Count      int    `json:"count"`
	Records    any    `json:"records"`
}

// ExportPrivacy renders one domain as indented JSON bytes.
func (s *Store) ExportPrivacy(ctx context.Context, domain PrivacyDomain) ([]byte, error) {
	var records any
	switch domain {
	case PrivacyProfile:
		rows, err := s.ListProfiles(ctx)
		if err != nil {
			return nil, err
		}
		records = nonNilSlice(rows)
	case PrivacyMemory:
		rows, err := s.ListMemories(ctx)
		if err != nil {
			return nil, err
		}
		records = nonNilSlice(rows)
	case PrivacyTaskLog:
		rows, err := s.ListTaskLogs(ctx)
		if err != nil {
			return nil, err
		}
		records = nonNilSlice(rows)
	case PrivacyToolCall:
		rows, err := s.ListToolCalls(ctx)
		if err != nil {
			return nil, err
		}
		records = nonNilSlice(rows)
	case PrivacyArtifacts:
		rows, err := s.ListArtifacts(ctx)
		if err != nil {
			return nil, err
		}
		records = nonNilSlice(rows)
	default:
		return nil, fmt.Errorf("memory: unknown privacy domain %q (valid: %v)", domain, PrivacyDomains())
	}
	return json.MarshalIndent(privacyExport{
		Domain:     string(domain),
		ExportedAt: observe.WallTimestampUTC(observe.NowWallUTC()),
		Records:    records,
	}, "", "  ")
}

// nonNilSlice makes empty exports render as [] instead of null.
func nonNilSlice[T any](rows []T) any {
	if rows == nil {
		return []T{}
	}
	return rows
}

// firstLine truncates a user-facing label to one line (labels only; the
// export always carries full records).
func firstLine(s string) string {
	for i, r := range s {
		if r == '\n' || r == '\r' {
			return s[:i]
		}
	}
	return s
}

// purgeAllDomains is the "wipe everything the privacy page exposes" helper
// (used by diagnostics and tests; the UI clears per domain).
func (s *Store) purgeAllDomains(ctx context.Context) error {
	for _, d := range PrivacyDomains() {
		if _, err := s.PurgePrivacy(ctx, d); err != nil {
			return err
		}
	}
	return nil
}
