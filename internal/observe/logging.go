package observe

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// JSONL log pipeline (ticket 08; SPEC-01 §3, [observe] config section).
//
// Structure: slog JSONHandler -> redactHandler (non-disableable, D33) ->
// rollingWriter (size + daily rolls, 7-day retention) -> <data>\logs\.
// A resident "log-flusher" goroutine (D38b roster, spawned through the
// registry so it carries the recover boundary) flushes the buffer every
// flushInterval and sweeps expired files once per retention sweep.
//
// Log line format: one JSON object per line (slog JSONHandler); timestamps
// are wall-clock UTC - a persisted record, per the package clock rules.

const (
	flushInterval        = 500 * time.Millisecond
	retentionSweepPeriod = 1 * time.Hour
	logFilePrefix        = "wisp-"
	logFileExt           = ".jsonl"
	logDayLayout         = "20060102"
	bufferedLogBytes     = 64 << 10
)

// LogConfig is the [observe] slice of the config model (schema.go
// ObserveSection) plus the [privacy] redact_paths flag. It is filled by the
// caller (cmd wiring / tickets 05+); zero values fall back to the schema
// defaults.
type LogConfig struct {
	// Dir is the log directory (<data>\logs\, SPEC-02 §6).
	Dir string
	// Level is debug|info|warn|error ("" = info).
	Level string
	// RollSizeMB is the size roll threshold (0 = 10MB, schema default).
	RollSizeMB int
	// RollDays is the retention window (0 = 7 days, schema default).
	RollDays int
	// RedactPaths is [privacy] redact_paths.
	RedactPaths bool
}

// LogPipeline owns the running log pipeline. Create with InitLog; Close
// during shutdown (after the workers stop logging, before process exit).
type LogPipeline struct {
	writer    *rollingWriter
	handler   slog.Handler
	root      *Root
	flushDone <-chan struct{}
}

// InitLog starts the log pipeline on the process registry (Default).
func InitLog(cfg LogConfig) (*LogPipeline, error) { return InitLogWithRegistry(cfg, Default) }

// InitLogWithRegistry starts the log pipeline on an explicit registry
// (isolated tests).
func InitLogWithRegistry(cfg LogConfig, reg *Registry) (*LogPipeline, error) {
	if cfg.Dir == "" {
		return nil, New(ClassConfig, "observe: log dir is empty")
	}
	if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {
		return nil, Wrap(ClassConfig, err, "observe: create log dir")
	}
	w, err := newRollingWriter(cfg.Dir, cfg.RollSizeMB, cfg.RollDays)
	if err != nil {
		return nil, err
	}
	level := slog.LevelInfo
	if cfg.Level != "" {
		if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Level))); err != nil {
			_ = w.Close()
			return nil, Wrap(ClassConfig, err, "observe: log level")
		}
	}
	// Redaction is non-disableable: there is no constructor switch for it.
	handler := &redactHandler{
		red: Redactor{RedactPaths: cfg.RedactPaths},
		next: slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
		}),
	}

	p := &LogPipeline{writer: w, handler: handler}
	p.root = NewRootFrom(context.Background(), "observe-logs")
	h := reg.Spawn("log-flusher", "observe", p.root, w.flushLoop)
	p.flushDone = h.Done()
	return p, nil
}

// Handler returns the redacting JSON handler. Wire it with
// slog.SetDefault(slog.New(p.Handler())) so every slog call in the process
// (including the ticket-03 panic sink) flows through redaction.
func (p *LogPipeline) Handler() slog.Handler { return p.handler }

// InstallAsDefault makes the pipeline the process-wide slog default.
func (p *LogPipeline) InstallAsDefault() { slog.SetDefault(slog.New(p.handler)) }

// FlushNow flushes buffered log lines to the current file (no fsync; the
// periodic flusher also runs while the pipeline is open).
func (p *LogPipeline) FlushNow() error { return p.writer.Flush() }

// Close stops the flusher goroutine, flushes, syncs and closes the current
// file. Logging after Close returns an error (and is dropped by slog).
func (p *LogPipeline) Close() error {
	p.root.Cancel()
	<-p.flushDone
	return p.writer.Close()
}

// redactHandler wraps a slog.Handler and forces every record through the
// Redactor, including records with pre-baked attributes.
type redactHandler struct {
	red  Redactor
	next slog.Handler
}

func (h *redactHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.next.Enabled(ctx, l)
}

func (h *redactHandler) Handle(ctx context.Context, rec slog.Record) error {
	out := slog.NewRecord(rec.Time, rec.Level, h.red.Message(rec.Message), rec.PC)
	rec.Attrs(func(a slog.Attr) bool {
		out.AddAttrs(h.red.Attr(a.Key, a.Value))
		return true
	})
	return h.next.Handle(ctx, out)
}

func (h *redactHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := make([]slog.Attr, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, h.red.Attr(a.Key, a.Value))
	}
	return &redactHandler{red: h.red, next: h.next.WithAttrs(out)}
}

func (h *redactHandler) WithGroup(name string) slog.Handler {
	return &redactHandler{red: h.red, next: h.next.WithGroup(name)}
}

// rollingWriter is an io.WriteCloser that rolls its file on size and on UTC
// day and prunes files older than the retention window. Safe for concurrent
// use (slog handlers serialize per record, but the flusher touches the same
// buffer).
type rollingWriter struct {
	mu      sync.Mutex
	dir     string
	maxSize int64
	days    int
	now     func() time.Time // wall clock, injectable for tests

	day     string // UTC day the current file belongs to
	seq     int    // per-day roll sequence, 1-based
	file    *os.File
	bw      *bufio.Writer
	written int64
	closed  bool
}

func newRollingWriter(dir string, sizeMB, days int) (*rollingWriter, error) {
	if sizeMB <= 0 {
		sizeMB = 10 // schema default ([observe] roll.size_mb)
	}
	if days <= 0 {
		days = 7 // schema default ([observe] roll.days)
	}
	w := &rollingWriter{
		dir:     dir,
		maxSize: int64(sizeMB) << 20,
		days:    days,
		now:     time.Now,
	}
	// Open eagerly so a broken log dir fails at init, not on first log line,
	// and run the retention sweep once at boot.
	if err := w.ensureFileLocked(true); err != nil {
		return nil, err
	}
	if err := w.sweepLocked(); err != nil {
		// Retention is best-effort at boot; the periodic sweep retries.
		slog.Warn("observe: log retention sweep failed", "err", err.Error())
	}
	return w, nil
}

func (w *rollingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return 0, errors.New("observe: log writer closed")
	}
	if err := w.ensureFileLocked(false); err != nil {
		return 0, err
	}
	n, err := w.bw.Write(p)
	w.written += int64(n)
	return n, err
}

// ensureFileLocked rolls to a new file when needed: first write, UTC day
// change (new sequence), or size threshold (next sequence, same day).
func (w *rollingWriter) ensureFileLocked(newDayCheck bool) error {
	day := w.now().UTC().Format(logDayLayout)
	if w.file != nil && day == w.day && w.written < w.maxSize {
		return nil
	}
	if err := w.closeFileLocked(); err != nil {
		return err
	}
	if w.day == day || w.file == nil && w.day == "" {
		// Same day: size roll (or first open today continues the sequence).
		w.seq++
	} else {
		w.seq = 1
	}
	if w.day != day {
		w.seq = 1
		w.day = day
	}
	name := filepath.Join(w.dir, fmt.Sprintf("%s%s-%03d%s", logFilePrefix, w.day, w.seq, logFileExt))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return Wrap(ClassResource, err, "observe: open log file")
	}
	w.file = f
	w.bw = bufio.NewWriterSize(f, bufferedLogBytes)
	w.written = 0
	return nil
}

func (w *rollingWriter) closeFileLocked() error {
	if w.file == nil {
		return nil
	}
	ferr := w.bw.Flush()
	serr := w.file.Sync()
	cerr := w.file.Close()
	w.file, w.bw, w.written = nil, nil, 0
	if ferr != nil {
		return Wrap(ClassResource, ferr, "observe: flush log file")
	}
	if serr != nil {
		return Wrap(ClassResource, serr, "observe: sync log file")
	}
	return cerr
}

// Flush drains the buffer into the current file (no fsync).
func (w *rollingWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.bw == nil {
		return nil
	}
	return w.bw.Flush()
}

func (w *rollingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
	return w.closeFileLocked()
}

// flushLoop is the log-flusher goroutine body: periodic buffer flush plus a
// once-per-period retention sweep.
func (w *rollingWriter) flushLoop(ctx context.Context) {
	t := time.NewTicker(flushInterval)
	defer t.Stop()
	lastSweep := w.now()
	for {
		select {
		case <-ctx.Done():
			_ = w.Flush()
			return
		case <-t.C:
			_ = w.Flush()
			if w.now().Sub(lastSweep) >= retentionSweepPeriod {
				lastSweep = w.now()
				w.mu.Lock()
				err := w.sweepLocked()
				w.mu.Unlock()
				if err != nil {
					slog.Warn("observe: log retention sweep failed", "err", err.Error())
				}
			}
		}
	}
}

// sweepLocked deletes log files older than the retention window. Files are
// named wisp-<day>-<seq>.jsonl; the day is parsed from the name first and
// the mtime is the fallback (rename/copy preserves names, not mtimes).
func (w *rollingWriter) sweepLocked() error {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return err
	}
	cutoff := w.now().UTC().Add(-time.Duration(w.days) * 24 * time.Hour)
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasPrefix(name, logFilePrefix) || !strings.HasSuffix(name, logFileExt) {
			continue
		}
		day := strings.TrimSuffix(strings.TrimPrefix(name, logFilePrefix), logFileExt)
		if i := strings.LastIndexByte(day, '-'); i >= 0 {
			day = day[:i]
		}
		t, terr := time.ParseInLocation(logDayLayout, day, time.UTC)
		if terr != nil {
			info, ierr := e.Info()
			if ierr != nil {
				continue
			}
			t = info.ModTime().UTC()
		}
		if t.Before(cutoff) {
			if rmErr := os.Remove(filepath.Join(w.dir, name)); rmErr != nil && !os.IsNotExist(rmErr) {
				return rmErr
			}
		}
	}
	return nil
}

// CountLogFiles returns the number of log files currently in dir (tests,
// diagnostics). It is a read-only helper, not part of the pipeline.
func CountLogFiles(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), logFilePrefix) && strings.HasSuffix(e.Name(), logFileExt) {
			n++
		}
	}
	return n, nil
}

var _ io.WriteCloser = (*rollingWriter)(nil)
