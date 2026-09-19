package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// approval_grant / cost_daily / plugin_state / schema_meta DAO (D35).
// Lifecycle POLICY (when grants are created/revoked, when cost days roll up)
// belongs to tickets 49/44 — this layer only stores rows.

// InsertGrant records a D45 session grant.
func (s *Store) InsertGrant(ctx context.Context, g ApprovalGrant) (int64, error) {
	if g.Scope != GrantScopeSession {
		return 0, fmt.Errorf("memory: invalid grant scope %q (want session)", g.Scope)
	}
	if g.Tool == "" || g.Pattern == "" || g.SessionID == "" {
		return 0, fmt.Errorf("memory: grant tool/pattern/session_id are required")
	}
	if g.CreatedAt == 0 {
		g.CreatedAt = observe.NowWallUTC().Unix()
	}
	if g.ExpiresAt == 0 {
		return 0, fmt.Errorf("memory: grant expires_at is required")
	}
	var id int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO approval_grant(scope, tool, pattern, session_id, created_at, expires_at, revoked_at)
			 VALUES(?,?,?,?,?,?,?)`,
			g.Scope, g.Tool, g.Pattern, g.SessionID, g.CreatedAt, g.ExpiresAt, nullableInt64(g.RevokedAt))
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("memory: insert grant: %w", err)
	}
	return id, nil
}

// RevokeGrant stamps revoked_at on a grant (idempotent).
func (s *Store) RevokeGrant(ctx context.Context, id int64) error {
	now := observe.NowWallUTC().Unix()
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`UPDATE approval_grant SET revoked_at=? WHERE id=? AND revoked_at IS NULL`, now, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			// Already revoked, or never existed: distinguish for the caller.
			var exists int
			if err := tx.QueryRowContext(ctx,
				`SELECT count(*) FROM approval_grant WHERE id=?`, id).Scan(&exists); err != nil {
				return err
			}
			if exists == 0 {
				return fmt.Errorf("approval_grant %d: %w", id, ErrNotFound)
			}
			return nil // idempotent re-revoke
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("memory: revoke grant %d: %w", id, err)
	}
	return nil
}

// ListGrantsBySession returns a session's grants (newest first).
func (s *Store) ListGrantsBySession(ctx context.Context, sessionID string) ([]ApprovalGrant, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, scope, tool, pattern, session_id, created_at, expires_at, revoked_at
		 FROM approval_grant WHERE session_id=? ORDER BY created_at DESC, id DESC`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("memory: list grants by session: %w", err)
	}
	defer rows.Close()
	return scanGrants(rows)
}

// ListGrants returns all grants (audit view, newest first).
func (s *Store) ListGrants(ctx context.Context) ([]ApprovalGrant, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, scope, tool, pattern, session_id, created_at, expires_at, revoked_at
		 FROM approval_grant ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("memory: list grants: %w", err)
	}
	defer rows.Close()
	return scanGrants(rows)
}

// DeleteGrant removes one grant row.
func (s *Store) DeleteGrant(ctx context.Context, id int64) error {
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM approval_grant WHERE id=?`, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("approval_grant %d: %w", id, ErrNotFound)
		}
		return nil
	})
}

func scanGrants(rows *sql.Rows) ([]ApprovalGrant, error) {
	var out []ApprovalGrant
	for rows.Next() {
		var g ApprovalGrant
		var revoked sql.NullInt64
		if err := rows.Scan(&g.ID, &g.Scope, &g.Tool, &g.Pattern, &g.SessionID,
			&g.CreatedAt, &g.ExpiresAt, &revoked); err != nil {
			return nil, err
		}
		g.RevokedAt = nullInt64Ptr(revoked)
		out = append(out, g)
	}
	return out, rows.Err()
}

// BumpCostDay adds one task's accounting to the C23 daily aggregate row
// (upsert + add; micros are integers, no floats).
func (s *Store) BumpCostDay(ctx context.Context, day string, tokensIn, tokensOut, costMicros int64) error {
	if !validDay(day) {
		return fmt.Errorf("memory: invalid cost_daily.day %q (want YYYY-MM-DD)", day)
	}
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO cost_daily(day, tokens_in, tokens_out, cost_micros, tasks) VALUES(?,?,?,?,1)
			 ON CONFLICT(day) DO UPDATE SET
			   tokens_in=tokens_in+excluded.tokens_in,
			   tokens_out=tokens_out+excluded.tokens_out,
			   cost_micros=cost_micros+excluded.cost_micros,
			   tasks=tasks+1`,
			day, tokensIn, tokensOut, costMicros)
		return err
	})
	if err != nil {
		return fmt.Errorf("memory: bump cost_daily %s: %w", day, err)
	}
	return nil
}

// CostDay returns one aggregate row (zero-value row when absent).
func (s *Store) CostDay(ctx context.Context, day string) (CostDaily, error) {
	var c CostDaily
	err := s.reader.QueryRowContext(ctx,
		`SELECT day, tokens_in, tokens_out, cost_micros, tasks FROM cost_daily WHERE day=?`, day).
		Scan(&c.Day, &c.TokensIn, &c.TokensOut, &c.CostMicros, &c.Tasks)
	if errors.Is(err, sql.ErrNoRows) {
		return CostDaily{Day: day}, nil
	}
	if err != nil {
		return CostDaily{}, fmt.Errorf("memory: cost_daily by day: %w", err)
	}
	return c, nil
}

// ListCostDaily returns aggregates, newest day first.
func (s *Store) ListCostDaily(ctx context.Context) ([]CostDaily, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT day, tokens_in, tokens_out, cost_micros, tasks FROM cost_daily ORDER BY day DESC`)
	if err != nil {
		return nil, fmt.Errorf("memory: list cost_daily: %w", err)
	}
	defer rows.Close()
	var out []CostDaily
	for rows.Next() {
		var c CostDaily
		if err := rows.Scan(&c.Day, &c.TokensIn, &c.TokensOut, &c.CostMicros, &c.Tasks); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// UpsertPluginState stores a §14.7 install record (id is the PK).
func (s *Store) UpsertPluginState(ctx context.Context, p PluginState) error {
	if p.ID == "" || p.Version == "" || p.Hash == "" {
		return fmt.Errorf("memory: plugin_state id/version/hash are required")
	}
	if p.CapabilitiesJSON == "" {
		return fmt.Errorf("memory: plugin_state.capabilities_json is required")
	}
	if p.InstalledAt == 0 {
		p.InstalledAt = observe.NowWallUTC().Unix()
	}
	enabled := int64(0)
	if p.Enabled {
		enabled = 1
	}
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO plugin_state(id, version, enabled, capabilities_json, net_allowlist_json, installed_at, hash, exe_hash)
			 VALUES(?,?,?,?,?,?,?,?)
			 ON CONFLICT(id) DO UPDATE SET
			   version=excluded.version, enabled=excluded.enabled,
			   capabilities_json=excluded.capabilities_json,
			   net_allowlist_json=excluded.net_allowlist_json,
			   installed_at=excluded.installed_at, hash=excluded.hash, exe_hash=excluded.exe_hash`,
			p.ID, p.Version, enabled, p.CapabilitiesJSON, nullIfEmpty(p.NetAllowlistJSON),
			p.InstalledAt, p.Hash, nullIfEmpty(p.ExeHash))
		return err
	})
	if err != nil {
		return fmt.Errorf("memory: upsert plugin_state %s: %w", p.ID, err)
	}
	return nil
}

// PluginStateByID returns one install record.
func (s *Store) PluginStateByID(ctx context.Context, id string) (PluginState, error) {
	var p PluginState
	var enabled int64
	var netAllow, exeHash sql.NullString
	err := s.reader.QueryRowContext(ctx,
		`SELECT id, version, enabled, capabilities_json, net_allowlist_json, installed_at, hash, exe_hash
		 FROM plugin_state WHERE id=?`, id).
		Scan(&p.ID, &p.Version, &enabled, &p.CapabilitiesJSON, &netAllow, &p.InstalledAt, &p.Hash, &exeHash)
	if errors.Is(err, sql.ErrNoRows) {
		return PluginState{}, fmt.Errorf("plugin_state %s: %w", id, ErrNotFound)
	}
	if err != nil {
		return PluginState{}, fmt.Errorf("memory: plugin_state by id: %w", err)
	}
	p.Enabled = enabled != 0
	p.NetAllowlistJSON = netAllow.String
	p.ExeHash = exeHash.String
	return p, nil
}

// ListPluginStates returns all install records ordered by id.
func (s *Store) ListPluginStates(ctx context.Context) ([]PluginState, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, version, enabled, capabilities_json, net_allowlist_json, installed_at, hash, exe_hash
		 FROM plugin_state ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("memory: list plugin_state: %w", err)
	}
	defer rows.Close()
	var out []PluginState
	for rows.Next() {
		var p PluginState
		var enabled int64
		var netAllow, exeHash sql.NullString
		if err := rows.Scan(&p.ID, &p.Version, &enabled, &p.CapabilitiesJSON, &netAllow,
			&p.InstalledAt, &p.Hash, &exeHash); err != nil {
			return nil, err
		}
		p.Enabled = enabled != 0
		p.NetAllowlistJSON = netAllow.String
		p.ExeHash = exeHash.String
		out = append(out, p)
	}
	return out, rows.Err()
}

// DeletePluginState removes one install record.
func (s *Store) DeletePluginState(ctx context.Context, id string) error {
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM plugin_state WHERE id=?`, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("plugin_state %s: %w", id, ErrNotFound)
		}
		return nil
	})
}

func validDay(day string) bool {
	if len(day) != 10 || day[4] != '-' || day[7] != '-' {
		return false
	}
	for i, r := range day {
		if i == 4 || i == 7 {
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
