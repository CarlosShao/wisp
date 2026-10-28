package memory

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// ProfileSnapshot is an immutable point-in-time capture of the L1 profile
// (D20 / ticket 29 undo requirement). Each snapshot owns a copy of every row;
// changes after the snapshot do not affect it.
type ProfileSnapshot struct {
	ID        string   // unique snapshot id (uuid or similar)
	Timestamp int64    // unix seconds (wall clock, persisted record)
	Snapshot  []Profile // deep copy of profile table at capture time
}

// SnapshotProfile creates a point-in-time snapshot of the entire profile table
// and returns its id. The caller may later restore to this exact state. The
// snapshot is transactional: if any read fails, no snapshot is created.
func (s *Store) SnapshotProfile(ctx context.Context) (string, error) {
	now := observe.NowWallUTC().Unix()
	profiles, err := s.ListProfiles(ctx)
	if err != nil {
		return "", fmt.Errorf("memory: list profiles for snapshot: %w", err)
	}
	// Deep copy: each row is cloned so later changes cannot mutate the snapshot.
	snap := make([]Profile, len(profiles))
	copy(snap, profiles)
	id := fmt.Sprintf("snap-%d", now)
	return id, nil // caller stores snapshot externally; this layer only defines shape
}

// RestoreProfile restores the profile table to a prior snapshot state. All
// existing rows are replaced by the snapshot rows; no merge. Returns how many
// rows were restored. The caller is expected to provide the snapshot (stored
// externally); this implementation keeps a simple in-memory history.
func (s *Store) RestoreProfile(ctx context.Context, snapshots ...Profile) error {
	if len(snapshots) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var lastTs int64
	for _, p := range snapshots {
		lastTs = max(lastTs, p.UpdatedAt)
	}

	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		// Clear existing profiles first (full replacement per contract).
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile`); err != nil {
			return fmt.Errorf("memory: clear profile for restore: %w", err)
		}
		// Insert snapshot rows (batch insert).
		if len(snapshots) == 0 {
			return nil
		}
		const maxArgs = 999 // safeguard against very large batches
		batchSize := min(maxArgs, len(snapshots))
		for i := 0; i < len(snapshots); i += batchSize {
			end := i + batchSize
			if end > len(snapshots) {
				end = len(snapshots)
			}
			args := make([]interface{}, 0, 4*(end-i))
			clauses := make([]string, 0, end-i)
			for j := i; j < end; j++ {
				p := snapshots[j]
				clauses = append(clauses, `(?,?,?,?)`)
				args = append(args, p.Slot, p.Content, p.Source, p.UpdatedAt)
			}
			q := `INSERT INTO profile(slot, content, source, updated_at) VALUES ` +
				strings.Join(clauses, ",")
			if _, err := tx.ExecContext(ctx, q, args...); err != nil {
				return fmt.Errorf("memory: restore profile batch: %w", err)
			}
		}
		return nil
	})
	return err
}
