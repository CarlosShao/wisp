package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	_ "modernc.org/sqlite"
)

// SQLite storage core (D35, SPEC-02 §2/§7).
//
// Layout under the data dir (SPEC-02 §6):
//
//	wisp.db (+ -wal/-shm)      the database
//	backup\                    pre-migration snapshots wisp.db.bak-<from>-<to>
//	artifacts\                 tool output artifacts (500MB LRU quota)
//
// Concurrency (D35 rule 3 / §14.11): every write is serialized through the
// single db-writer goroutine; reads run concurrently on their own pool.
// Connections are opened on demand and closed after an idle period (D20);
// nothing here keeps a file handle permanently.
const (
	// DriverName is the modernc.org/sqlite registration name (pure Go, no cgo
	// — SPEC-02 §2; the cgo budget stays with sherpa-onnx).
	DriverName = "sqlite"

	dbFileName       = "wisp.db"
	backupDirName    = "backup"
	artifactsDirName = "artifacts"

	// PRAGMA values (D35 rule 3 / SPEC-02 §2, contract-level).
	pragmaBusyTimeoutMS      = 3000
	pragmaWALAutocheckpoint  = 1000
	startupCheckpointQuery   = "PRAGMA wal_checkpoint(TRUNCATE)"
	integrityCheckQuery      = "PRAGMA integrity_check"
	walPageSizeQuery         = "PRAGMA page_size"
	walCheckpointPassiveOpts = "PASSIVE"

	// defaultIdleConnTTL bounds how long pooled connections stay open after
	// their last use (D20: connections on demand, not resident).
	defaultIdleConnTTL = 30 * time.Second
)

// ErrStoreClosed is returned by operations after Close.
var ErrStoreClosed = errors.New("memory: store is closed")

// ErrSchemaUnmigratable marks a database this binary cannot migrate (SPEC-02
// §7): the caller must surface it into the Unconfigured state with a clear
// message. The database is NEVER modified in that case — silently overwriting
// it would reset directory allowlists and permission settings (a security
// bug, not a data bug).
var ErrSchemaUnmigratable = errors.New(
	"memory: database schema cannot be migrated; refusing to touch wisp.db " +
		"(restore the file from backup\\ or remove it, or downgrade the binary)")

// SchemaError details an unmigratable database. errors.Is(err,
// ErrSchemaUnmigratable) holds for every instance.
type SchemaError struct {
	CurrentVersion int    // -1 when unknown
	TargetVersion  int    // SchemaVersionTarget of this binary
	Reason         string // what exactly blocked the migration
}

func (e *SchemaError) Error() string {
	return fmt.Sprintf("schema_version=%d (target %d): %s", e.CurrentVersion, e.TargetVersion, e.Reason)
}

func (e *SchemaError) Unwrap() error { return ErrSchemaUnmigratable }

// IsSchemaUnmigratable reports whether err is (or wraps) an unmigratable-
// schema failure. Callers map it to the Unconfigured ball state.
func IsSchemaUnmigratable(err error) bool { return errors.Is(err, ErrSchemaUnmigratable) }

// Store owns wisp.db and the artifacts directory. Create with Open; all
// methods are safe for concurrent use.
type Store struct {
	dir          string
	dbPath       string
	backupDir    string
	artifactsDir string

	logger *slog.Logger
	reg    *observe.Registry
	target int // schema version this store migrates to

	reader *sql.DB // concurrent read pool (WAL readers never block)
	writer *sql.DB // single-connection write pool behind db-writer
	queue  *writeQueue

	closeOnce sync.Once
	closeErr  error
}

// Option configures a Store.
type Option func(*storeOptions)

type storeOptions struct {
	logger      *slog.Logger
	reg         *observe.Registry
	idleConnTTL time.Duration
	target      int // schema version this store migrates to (tests override)
}

func defaultStoreOptions() storeOptions {
	return storeOptions{
		logger:      slog.Default(),
		reg:         observe.Default,
		idleConnTTL: defaultIdleConnTTL,
		target:      SchemaVersionTarget,
	}
}

// WithLogger overrides the logger (tests inject a capturing handler).
func WithLogger(l *slog.Logger) Option { return func(o *storeOptions) { o.logger = l } }

// WithRegistry spawns db-writer through a specific registry (tests).
func WithRegistry(r *observe.Registry) Option { return func(o *storeOptions) { o.reg = r } }

// WithIdleConnTTL overrides the on-demand connection idle lifetime (tests).
func WithIdleConnTTL(d time.Duration) Option {
	return func(o *storeOptions) { o.idleConnTTL = d }
}

// withSchemaTarget overrides the migration target version. Unexported: the
// production target is SchemaVersionTarget; only the migration-chain test
// raises it (together with a test-only chain entry).
func withSchemaTarget(v int) Option { return func(o *storeOptions) { o.target = v } }

// dsn builds the modernc DSN. Pragmas are attached per connection so every
// pooled connection (reader or writer) inherits the contract settings
// (D35 rule 3). The writer DSN additionally takes the immediate tx lock:
// with a single writer this removes SQLITE_BUSY upgrade deadlocks entirely,
// so no ad-hoc busy-retry logic exists anywhere (D35 rule 3).
func dsn(path string, forWriter bool) string {
	d := "file:" + filepath.ToSlash(path) +
		"?_pragma=busy_timeout(" + itoa(pragmaBusyTimeoutMS) + ")" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=synchronous(NORMAL)" +
		"&_pragma=wal_autocheckpoint(" + itoa(pragmaWALAutocheckpoint) + ")"
	if forWriter {
		d += "&_txlock=immediate"
	}
	return d
}

func itoa(n int) string { return fmt.Sprintf("%d", n) }

// Open creates (or opens) wisp.db under dir, runs the schema migration chain
// (with pre-backups, SPEC-02 §7) and the startup WAL checkpoint. An
// unmigratable database returns a *SchemaError (errors.Is
// ErrSchemaUnmigratable) and leaves the file untouched.
func Open(dir string, opts ...Option) (*Store, error) {
	o := defaultStoreOptions()
	for _, opt := range optList(opts) {
		opt(&o)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("memory: resolve data dir: %w", err)
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("memory: create data dir: %w", err)
	}
	s := &Store{
		dir:          abs,
		dbPath:       filepath.Join(abs, dbFileName),
		backupDir:    filepath.Join(abs, backupDirName),
		artifactsDir: filepath.Join(abs, artifactsDirName),
		logger:       o.logger,
		reg:          o.reg,
		target:       o.target,
	}
	if err := os.MkdirAll(s.artifactsDir, 0o755); err != nil {
		return nil, fmt.Errorf("memory: create artifacts dir: %w", err)
	}
	if err := s.migrate(); err != nil {
		return nil, err
	}

	// Pools: reader pool may hold several concurrent connections; the writer
	// pool is capped at exactly one connection so the single-writer invariant
	// holds at the SQL layer too (defense in depth behind the db-writer
	// goroutine). Both close idle connections after idleConnTTL (D20).
	s.reader, err = sql.Open(DriverName, dsn(s.dbPath, false))
	if err != nil {
		return nil, fmt.Errorf("memory: open reader pool: %w", err)
	}
	s.reader.SetMaxOpenConns(8)
	s.reader.SetMaxIdleConns(8)
	s.reader.SetConnMaxIdleTime(o.idleConnTTL)
	s.writer, err = sql.Open(DriverName, dsn(s.dbPath, true))
	if err != nil {
		_ = s.reader.Close()
		return nil, fmt.Errorf("memory: open writer pool: %w", err)
	}
	s.writer.SetMaxOpenConns(1)
	s.writer.SetMaxIdleConns(1)
	s.writer.SetConnMaxIdleTime(o.idleConnTTL)
	s.queue = newWriteQueue(s.writer, s.reg, s.logger)

	// Startup WAL checkpoint (D35 rule 3): bound the WAL after the previous
	// run — this is part of the "disk full" defense.
	if err := s.startupCheckpoint(); err != nil {
		_ = s.Close()
		return nil, err
	}
	return s, nil
}

func optList(opts []Option) []Option { return opts }

// startupCheckpoint runs PRAGMA wal_checkpoint(TRUNCATE) and logs the result.
func (s *Store) startupCheckpoint() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var busy, logPages, checkpointed int
	if err := s.writer.QueryRowContext(ctx, startupCheckpointQuery).Scan(&busy, &logPages, &checkpointed); err != nil {
		return fmt.Errorf("memory: startup wal_checkpoint(TRUNCATE): %w", err)
	}
	if busy != 0 {
		s.logger.Warn("startup WAL checkpoint did not complete (readers active)",
			"component", "memory", "busy", busy, "wal_pages", logPages)
		return nil
	}
	s.logger.Info("startup WAL checkpoint (TRUNCATE) done",
		"component", "memory", "wal_pages_before", logPages, "pages_moved", checkpointed)
	return nil
}

// Dir returns the data directory backing the store.
func (s *Store) Dir() string { return s.dir }

// ArtifactsDir returns the artifacts directory path.
func (s *Store) ArtifactsDir() string { return s.artifactsDir }

// Close flushes the write queue, checkpoints the WAL, and closes all
// connections. Idempotent.
func (s *Store) Close() error {
	s.closeOnce.Do(func() {
		if s.queue != nil {
			s.queue.close()
		}
		// Best-effort final checkpoint so a clean exit leaves no WAL growth
		// behind (the startup checkpoint is the contractual one).
		if s.writer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, _ = s.writer.ExecContext(ctx, startupCheckpointQuery)
			cancel()
		}
		if s.reader != nil {
			_ = s.reader.Close()
		}
		if s.writer != nil {
			_ = s.writer.Close()
		}
		s.closeErr = nil
	})
	return s.closeErr
}

// write enqueues one write transaction on the db-writer goroutine and waits
// for it. fn runs inside a single transaction; returning an error rolls it
// back. This is the ONLY write path — no code may Exec against s.writer
// directly (D35 rule 3: zero scattered SQLITE_BUSY retry).
func (s *Store) write(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	if s.queue == nil {
		return ErrStoreClosed
	}
	return s.queue.enqueue(ctx, fn)
}

// ---------------------------------------------------------------------------
// Migration chain (SPEC-02 §7): hand-written, no framework. Each step runs in
// one transaction after a pre-backup; unmigratable databases fail without
// modification.
// ---------------------------------------------------------------------------

func (s *Store) migrate() error {
	boot, err := sql.Open(DriverName, dsn(s.dbPath, true))
	if err != nil {
		return fmt.Errorf("memory: open for migration: %w", err)
	}
	version, fresh, err := readSchemaVersion(boot)
	if err != nil {
		_ = boot.Close()
		return err
	}
	if version == s.target {
		_ = boot.Close()
		return nil
	}
	if version > s.target {
		_ = boot.Close()
		return &SchemaError{
			CurrentVersion: version,
			TargetVersion:  s.target,
			Reason:         "database was written by a newer Wisp binary; downgrading is not supported",
		}
	}
	steps, err := planMigration(version, fresh, s.target)
	if err != nil {
		_ = boot.Close()
		return err
	}

	// Apply each step: pre-backup (files, no open handles), then the single
	// transaction. Connection is closed around the file copy so the snapshot
	// is consistent.
	for _, step := range steps {
		if err := boot.Close(); err != nil {
			return fmt.Errorf("memory: close before backup: %w", err)
		}
		if err := s.backupDatabase(step.from, step.to); err != nil {
			return err
		}
		boot, err = sql.Open(DriverName, dsn(s.dbPath, true))
		if err != nil {
			return fmt.Errorf("memory: reopen for migration: %w", err)
		}
		if err := applyMigrationStep(boot, step, s.logger); err != nil {
			_ = boot.Close()
			return fmt.Errorf("memory: migration %d->%d failed (database restored from backup is at %s): %w",
				step.from, step.to, s.backupPath(step.from, step.to), err)
		}
		version = step.to
	}
	_ = boot.Close()
	s.logger.Info("schema migration complete",
		"component", "memory", "schema_version", version)
	return nil
}

// planMigration picks the chain steps from the current version. A fresh
// database (no user tables) starts at 0 even without a schema_meta row; a
// database with tables but an unusable watermark is unmigratable.
func planMigration(version int, fresh bool, target int) ([]migrationStep, error) {
	if fresh {
		version = 0
	}
	if version > target {
		return nil, &SchemaError{
			CurrentVersion: version,
			TargetVersion:  target,
			Reason:         "newer schema than this binary",
		}
	}
	var steps []migrationStep
	v := version
	for v < target {
		next, ok := nextStep(v)
		if !ok {
			return nil, &SchemaError{
				CurrentVersion: version,
				TargetVersion:  target,
				Reason:         fmt.Sprintf("no migration step %d->%d in this binary", v, v+1),
			}
		}
		steps = append(steps, next)
		v = next.to
	}
	return steps, nil
}

func nextStep(from int) (migrationStep, bool) {
	for _, st := range migrationChain {
		if st.from == from {
			return st, true
		}
	}
	return migrationStep{}, false
}

// applyMigrationStep runs one step in a single transaction and updates the
// watermark inside it. Start/end times are recorded (SPEC-02 §7).
func applyMigrationStep(db *sql.DB, step migrationStep, logger *slog.Logger) error {
	ctx := context.Background()
	start := observe.NewTimeout(0)
	startWall := observe.NowWallUTC()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration tx: %w", err)
	}
	if err := step.apply(tx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("apply: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO schema_meta(key, value) VALUES(?, ?)
		 ON CONFLICT(key) DO UPDATE SET value=excluded.value`,
		schemaMetaVersionKey, itoa(step.to)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update schema_version: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration: %w", err)
	}
	logger.Info("schema migration step applied",
		"component", "memory",
		"from", step.from, "to", step.to,
		"started_at", observe.WallTimestampUTC(startWall),
		"duration_ms", start.Elapsed().Milliseconds())
	return nil
}

// readSchemaVersion reads the migration watermark. fresh=true means the file
// has no user tables (brand-new or 0-byte file): a v0 database.
func readSchemaVersion(db *sql.DB) (version int, fresh bool, err error) {
	ctx := context.Background()
	var nTables int
	if err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).
		Scan(&nTables); err != nil {
		return 0, false, fmt.Errorf("memory: inspect sqlite_master: %w", err)
	}
	if nTables == 0 {
		return 0, true, nil
	}
	// Tables exist: the watermark must be readable, otherwise the database is
	// foreign or damaged and MUST NOT be touched (SPEC-02 §7).
	var hasMeta int
	if err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='schema_meta'`).
		Scan(&hasMeta); err != nil {
		return 0, false, fmt.Errorf("memory: inspect schema_meta: %w", err)
	}
	if hasMeta == 0 {
		return 0, false, &SchemaError{
			CurrentVersion: -1,
			TargetVersion:  SchemaVersionTarget,
			Reason:         "database has tables but no schema_meta watermark table",
		}
	}
	var raw string
	err = db.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key=?`, schemaMetaVersionKey).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, &SchemaError{
			CurrentVersion: -1,
			TargetVersion:  SchemaVersionTarget,
			Reason:         "schema_meta exists but has no schema_version row",
		}
	}
	if err != nil {
		return 0, false, fmt.Errorf("memory: read schema_version: %w", err)
	}
	v, convErr := parseVersion(raw)
	if convErr != nil {
		return 0, false, &SchemaError{
			CurrentVersion: -1,
			TargetVersion:  SchemaVersionTarget,
			Reason:         "schema_version is not an integer: " + convErr.Error(),
		}
	}
	return v, false, nil
}

// parseVersion parses the schema_version watermark STRICTLY: digits only
// (whitespace-trimmed). Lenient prefix parsing would accept garbage like
// "1abc" as 1 and silently migrate a foreign database (adversarial MINOR-5).
func parseVersion(raw string) (int, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, fmt.Errorf("empty schema_version")
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("non-numeric schema_version %q", raw)
		}
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("schema_version %q out of range: %w", raw, err)
	}
	return v, nil
}

// backupPath renders the SPEC-02 §7 backup name: backup\wisp.db.bak-<from>-<to>.
func (s *Store) backupPath(from, to int) string {
	return filepath.Join(s.backupDir, fmt.Sprintf("%s.bak-%d-%d", dbFileName, from, to))
}

// backupDatabase snapshots wisp.db (plus -wal/-shm when present) into
// backup\wisp.db.bak-<from>-<to>. An existing backup is kept (the oldest
// pre-migration snapshot is the most valuable one; a retry after a failed
// migration must not destroy it).
func (s *Store) backupDatabase(from, to int) error {
	dst := s.backupPath(from, to)
	if err := os.MkdirAll(s.backupDir, 0o755); err != nil {
		return fmt.Errorf("memory: create backup dir: %w", err)
	}
	if _, err := os.Stat(dst); err == nil {
		s.logger.Warn("migration backup already exists, keeping it",
			"component", "memory", "backup", dst)
		return nil
	}
	for _, suffix := range []string{"", "-wal", "-shm"} {
		src := s.dbPath + suffix
		if suffix != "" {
			if _, err := os.Stat(src); err != nil {
				continue // no WAL/SHM sidecar: nothing to copy
			}
		}
		if err := copyFile(src, dst+suffix); err != nil {
			return fmt.Errorf("memory: backup %s: %w", filepath.Base(src), err)
		}
	}
	s.logger.Info("pre-migration backup written",
		"component", "memory", "backup", dst, "from", from, "to", to)
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
