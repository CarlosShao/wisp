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
	return newRollingWriterClock(dir, sizeMB, days, time.Now)
}

// newRollingWriterClock is the seam newRollingWriter runs through: the clock is
// a PARAMETER so a day-roll test can install its fake day BEFORE the eager
// open. With the wall clock baked in, the boot file is named after the real
// today and any fixture day is an extra file - which is what made
// TestRollingWriterDayRoll green only on the two dates its fixture hardcoded.
func newRollingWriterClock(dir string, sizeMB, days int, now func() time.Time) (*rollingWriter, error) {
	if sizeMB <= 0 {
		sizeMB = 10 // schema default ([observe] roll.size_mb)
	}
	if days <= 0 {
		days = 7 // schema default ([observe] roll.days)
	}
	if now == nil {
		now = time.Now
	}
	w := &rollingWriter{
		dir:     dir,
		maxSize: int64(sizeMB) << 20,
		days:    days,
		now:     now,
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

// ---------------------------------------------------------------------------
// The early log buffer (ticket 130, ruling "a-with-mirror" in
// docs/reports/pending-and-issues.md A110 row 3).
//
// WHAT IT FIXES. Go runs a package's init() before main(), so a record made
// from an init() is emitted while the process default logger is still the stock
// one and no JSONL pipeline exists yet. Today that is one path:
// internal/risk/winsec_c26.go's init() calls internal/winsec's
// SetPathResolver, and that function logs the conformance verdict - either
// "resolver installed" or "refusing to install". Both branches were reaching
// the terminal and nothing else, so the single most auditable fact about the
// sealing seam survived only as far as a console that the resident leg does not
// have.
//
// WHY IT LIVES IN THIS PACKAGE. The main package's init() runs after every
// dependency's, so nothing under cmd/ can catch a record made by another
// package's init() - that is the trap this ticket's first draft fell into. Go
// initializes dependencies before dependents, and internal/risk imports
// internal/observe, so this package's init() is necessarily earlier than the
// emitting one. internal/winsec must not import this package (the graph runs
// observe -> secret -> winsec, so the edge would close a cycle - see
// cmd/wisp/logsink.go's header), which is exactly why the buffer is pushed onto
// the default logger here rather than pulled by the emitter there.
//
// WHY THE CONSOLE COPY IS NEVER SWITCHED OFF. The buffer is a SECOND copy, not
// a relocation: every record still reaches the console, at the same level gate
// and exactly once, from the moment this package initializes. (The console
// handler is a TextHandler rather than the one Go had installed, which is a
// shape change of the timestamp prefix and nothing else - see init() for the
// cycle that capturing the stock handler creates.) The failure set of this
// design is therefore today's failure set plus a record on disk: if the process
// dies before a listener is installed, or runs a leg that never installs one,
// everything it would have printed today still prints.
// ---------------------------------------------------------------------------

const (
	// earlyLogMaxRecords and earlyLogMaxBytes bound the in-memory copy. The
	// measured scale is one record per boot (the sealing verdict); 64 records
	// leaves room for further packages to speak from their init() without
	// touching a budget. The byte cap is the same 64 KiB the flush buffer uses
	// (bufferedLogBytes). Neither cap costs filesystem activity at boot, which
	// is what keeps this out of the D32 cold-read budget.
	earlyLogMaxRecords = 64
	earlyLogMaxBytes   = bufferedLogBytes

	// earlyLogOverflowMsg is booked by a replay that hit a cap: "dropped N
	// records" has to be a fact on disk, not a silent subtraction. The OLDEST
	// kept message is named because in a boot ordering the earliest record is
	// usually the cause and the newest one only the effect.
	earlyLogOverflowMsg = "observe: early log buffer overflow"
)

// earlyLogBuffer holds the second copy. Package-level because the point of the
// exercise is records made before any object exists to hold them.
type earlyLogBuffer struct {
	mu      sync.Mutex
	recs    []slog.Record
	size    int64
	dropped int
	first   string // message of the first kept record, named by the overflow record
	closed  bool   // set by the first replay; later records are the pipeline's job
}

// earlyLogs is the buffer installed as the buffering half of the process default
// by the init() below.
var earlyLogs = &earlyLogBuffer{}

// earlyTeeInstalled keeps the tee this package's init() installed addressable,
// so a test can pin what its console mirror is: the live process default may have
// been replaced by any test that calls slog.SetDefault, and the property worth
// pinning is the one decided at init(). See internal/observe/earlylog_130_test.go.
var earlyTeeInstalled earlyTee

func init() {
	// The console copy is a concrete TextHandler, NOT the handler Go had
	// installed. That is not a style choice and it cost a deadlock to learn:
	// Go's own default handler is a bridge into the log package, whose output is
	// in turn bridged back into slog - resolved through slog.Default() at WRITE
	// time, not at construction time. Capturing it here and then replacing the
	// default is therefore a cycle: log.Print -> slog.Default -> this tee -> the
	// captured stock handler -> log.Print -> the same non-reentrant log mutex,
	// forever. internal/observe's own TestUnknownNameIsLeakSymptom is what
	// catches it, because Registry.Spawn logs through the default while holding
	// the registry (it hung this package's gate for 600s the first time).
	//
	// The visible cost of the TextHandler is the SHAPE of a pre-install console
	// line: Go's stock form prints "2026-09-23 11:49:39 INFO msg", the
	// TextHandler form prints "time=... level=INFO msg=...". Every record this
	// ticket installs already has the second shape (cmd/wisp/logsink.go's mirror
	// is the same handler), so the effect is that one line stops being the odd
	// one out. What is unchanged is the promise ticket 130 is about: the console
	// still receives the record, exactly once, and the level gate is the same
	// LevelInfo the stock default applied.
	tee := earlyTee{
		buf:    &earlyBufferHandler{b: earlyLogs},
		mirror: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}),
	}
	earlyTeeInstalled = tee
	slog.SetDefault(slog.New(tee))
}

// FlushEarlyLogRecords replays the buffered early records into next, in the
// order they were captured, and closes the buffer: after the first call nothing
// is buffered any more, because the caller has installed a listener that carries
// records itself. Idempotent - a second call replays nothing.
//
// It returns the number of records the handler accepted and the number that will
// never be seen on disk (over the cap, or refused by the handler). The caller
// books both counts; a reader holding one file off a disk should be able to tell
// "no early record happened" apart from "the replay did not run".
//
// This is called by cmd/wisp/logsink.go immediately after slog.SetDefault and
// BEFORE the install record, which is what puts a record made before the
// listener ahead of the listener's own booking in the same file.
func FlushEarlyLogRecords(next slog.Handler) (flushed, dropped int) {
	return earlyLogs.drain(next)
}

// drain takes the whole buffer under the lock and writes it out unlocked: a
// handler that calls back into logging must not deadlock against the buffer.
func (b *earlyLogBuffer) drain(next slog.Handler) (int, int) {
	b.mu.Lock()
	recs, dropped, first := b.recs, b.dropped, b.first
	b.recs, b.size, b.dropped, b.first, b.closed = nil, 0, 0, "", true
	b.mu.Unlock()

	if len(recs) == 0 && dropped == 0 {
		return 0, 0
	}
	ctx := context.Background()
	flushed := 0
	for _, r := range recs {
		if !next.Enabled(ctx, r.Level) {
			// Below the installed sink's level: this record was never going to
			// be a persisted one. Counted, not written.
			dropped++
			continue
		}
		if err := next.Handle(ctx, r.Clone()); err != nil {
			dropped++
			continue
		}
		flushed++
	}
	if dropped > 0 {
		rec := slog.NewRecord(time.Now(), slog.LevelWarn, earlyLogOverflowMsg, 0)
		rec.AddAttrs(
			slog.Int("dropped", dropped),
			slog.String("first_kept", first),
			slog.Int("capacity", earlyLogMaxRecords),
		)
		if next.Enabled(ctx, rec.Level) {
			_ = next.Handle(ctx, rec)
		}
	}
	return flushed, dropped
}

// add keeps one record if there is room, and reports whether it was kept.
func (b *earlyLogBuffer) add(r slog.Record) bool {
	n := earlyRecordSize(r)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return false
	}
	if len(b.recs) >= earlyLogMaxRecords || b.size+n > earlyLogMaxBytes {
		b.dropped++
		return false
	}
	if len(b.recs) == 0 {
		b.first = r.Message
	}
	b.recs = append(b.recs, r.Clone())
	b.size += n
	return true
}

// earlyRecordSize is the memory a record costs the buffer. An estimate on
// purpose - the cap exists to bound retention, not to account exactly.
func earlyRecordSize(r slog.Record) int64 {
	n := int64(len(r.Message)) + int64(r.NumAttrs())*8
	r.Attrs(func(a slog.Attr) bool {
		n += int64(len(a.Key)) + int64(len(a.Value.String())) + 8
		return true
	})
	return n
}

// earlyTee is the handler the init() above installs as the process default: one
// record, the buffer copy first and the console copy second. Buffering first is
// the order that matters - a console that cannot be written (a detached GUI
// process, a closed stderr) must not cost the copy that exists to outlive it.
//
// Enabled delegates to the console, so the level gate in front of this tee is
// exactly the gate Go's own default applied: nothing that reaches it today gets
// dropped by it, and nothing that was dropped before gets kept now.
type earlyTee struct {
	buf    *earlyBufferHandler
	mirror slog.Handler
}

func (t earlyTee) Enabled(ctx context.Context, level slog.Level) bool {
	return t.mirror.Enabled(ctx, level)
}

func (t earlyTee) Handle(ctx context.Context, r slog.Record) error {
	_ = t.buf.Handle(ctx, r)
	if t.mirror.Enabled(ctx, r.Level) {
		_ = t.mirror.Handle(ctx, r)
	}
	// No error is propagated: Go's stock default handler never returned one
	// either, and a record nobody could print is still a record.
	return nil
}

func (t earlyTee) WithAttrs(attrs []slog.Attr) slog.Handler {
	return earlyTee{buf: t.buf.withAttrs(attrs), mirror: t.mirror.WithAttrs(attrs)}
}

func (t earlyTee) WithGroup(name string) slog.Handler {
	return earlyTee{buf: t.buf.withGroup(name), mirror: t.mirror.WithGroup(name)}
}

// earlyBufferHandler is the buffering half of the tee. It keeps the attributes
// and groups added by With/WithGroup so a logger captured off the default before
// any listener exists - internal/memory/open.go and internal/agent/loop.go both
// do exactly that - still contributes the same record shape to the buffer that
// it contributes to the console.
//
// It deliberately does NOT satisfy slog.Handler on its own: WithAttrs and
// WithGroup are reached through earlyTee, which is the only thing ever installed
// as a handler, so a chain can never buffer on one side and console on the other.
type earlyBufferHandler struct {
	b      *earlyLogBuffer
	prefix []earlyAttr // chain attributes, in order, each with the groups open when it was added
	group  []string    // group names open right now, outermost first; a "" name is a no-op
}

// earlyAttr is one attribute plus the group path it was added under.
type earlyAttr struct {
	path []string
	attr slog.Attr
}

func (h *earlyBufferHandler) Handle(_ context.Context, r slog.Record) error {
	own := make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		own = append(own, a)
		return true
	})
	out := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	out.AddAttrs(materializeEarlyAttrs(h.prefix, h.group, own)...)
	h.b.add(out)
	return nil
}

func (h *earlyBufferHandler) withAttrs(attrs []slog.Attr) *earlyBufferHandler {
	if len(attrs) == 0 {
		return h
	}
	merged := make([]earlyAttr, 0, len(h.prefix)+len(attrs))
	merged = append(merged, h.prefix...)
	for _, a := range attrs {
		merged = append(merged, earlyAttr{path: h.group, attr: a})
	}
	return &earlyBufferHandler{b: h.b, prefix: merged, group: h.group}
}

func (h *earlyBufferHandler) withGroup(name string) *earlyBufferHandler {
	if name == "" {
		return h
	}
	group := make([]string, 0, len(h.group)+1)
	group = append(group, h.group...)
	return &earlyBufferHandler{b: h.b, prefix: h.prefix, group: append(group, name)}
}

// materializeEarlyAttrs turns a With-chain plus the attributes of one call into
// the attribute list a handler would print: attributes that are contiguous and
// share an open group become ONE group attribute, in order, nested as deep as
// their paths say.
//
// The contiguity is not a shortcut, it is what log/slog itself does - a group
// reopened later prints as a second object with the same key - and getting it
// wrong here is a silent data loss rather than a shape difference: two
// same-keyed group attributes in one JSON object are resolved by most parsers by
// keeping the last, so the first group's attributes would vanish from the file
// while still being on the console.
func materializeEarlyAttrs(prefix []earlyAttr, group []string, own []slog.Attr) []slog.Attr {
	items := make([]earlyAttr, 0, len(prefix)+len(own))
	items = append(items, prefix...)
	for _, a := range own {
		items = append(items, earlyAttr{path: group, attr: a})
	}
	if len(items) == 0 {
		return nil
	}
	type openGroup struct {
		name  string
		attrs []slog.Attr
	}
	var (
		top   []slog.Attr
		stack []openGroup
	)
	closeTo := func(depth int) {
		for len(stack) > depth {
			f := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			values := make([]any, 0, len(f.attrs))
			for _, a := range f.attrs {
				values = append(values, a)
			}
			nested := slog.Group(f.name, values...)
			if len(stack) > 0 {
				stack[len(stack)-1].attrs = append(stack[len(stack)-1].attrs, nested)
				continue
			}
			top = append(top, nested)
		}
	}
	for _, it := range items {
		depth := 0
		for depth < len(stack) && depth < len(it.path) && stack[depth].name == it.path[depth] {
			depth++
		}
		closeTo(depth)
		for i := depth; i < len(it.path); i++ {
			stack = append(stack, openGroup{name: it.path[i]})
		}
		if len(stack) > 0 {
			stack[len(stack)-1].attrs = append(stack[len(stack)-1].attrs, it.attr)
			continue
		}
		top = append(top, it.attr)
	}
	closeTo(0)
	return top
}

var _ slog.Handler = earlyTee{}
