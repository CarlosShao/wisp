package main

// The persistent log sink, installed by the two production entry points
// (ticket 117).
//
// What this file fixes is a listener, not a judgment: the security decisions in
// internal/winsec (a seal clearing a principal that stood on an object out of
// band), internal/config (a locked section loosened only after L2), the D33
// credential migration and the whole `[audit]` trail of `wisp run` - the path
// rewrites booked by D31 above all - were already being written to a logger.
// They just went to the process-default slog logger, which is stderr, and to
// `rt.stderr`, which is also a stream nobody keeps. Owner double-clicks the
// icon, no terminal is attached, and "your out-of-band authorization was
// cleared" is an event that leaves no trace.
//
// Before this file the JSONL pipeline in internal/observe had exactly ONE
// caller in the whole non-test tree: `wisp slo` (slo_windows.go, the diagnostic
// command). Every test that proved these notices "are recorded" installed its
// own sink, which is the shape this repository has now booked six times as
// "capability finished, nobody listening" - and the most expensive of the six,
// because it downgrades every earlier ticket's "loud failure" to "loud in a
// test binary".
//
// Two invariants this file exists to keep, both stated as grep-able claims in
// the ticket:
//
//  1. Ordering: install before anything that can log. On the `run` leg that
//     means before assembleRuntime, whose FIRST statement
//     (secret.NewStore -> winsec.PrivateDirAll) is already a sealing site; on
//     the resident leg it means immediately after proc.Boot, which performs no
//     sealing at all (internal/proc has zero winsec imports). A sink installed
//     after the first event loses exactly the record it exists to keep.
//  2. Landing spot: <the env's data root>\logs - the same spelling `wisp slo`
//     already writes, and the tree ticket 95 put under private-data discipline
//     (config.toml, the DPAPI blob dir, memory.db and the backups all land
//     sealed inside it). Never a world-shared temp default: a security outlet
//     that anyone on the machine can read is a leak path, not a listener.
//
// INTERIM (ticket 117 AC#3): whether a log file counts as private data is the
// one question of ticket 95's that owner has not answered, so nothing here
// seals or un-seals the log file - it neither adds nor removes that decision.
// What this file does guarantee is *placement*: the sink is created only under
// the resolved data dir of the running env, so the question owner has to answer
// is about one directory, not about a temp file. See logSinkDir and its
// invariant test.

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/observe"
)

// logDirName is the pipeline's directory inside the env's data root (SPEC-02 §6
// puts logs at <data>\logs\). Identical to the spelling `wisp slo` uses, so one
// env has one log tree and `wisp slo`'s retention sweep sees these files too.
const logDirName = "logs"

// logSinkLevel is the level filter the installed handler runs at: the
// [observe] schema default ("info", internal/config/schema.go
// ObserveSection.Level), which is also the value `wisp slo` passes today. The
// notices this exists to keep are WARN, so "info" is not a loosening.
const logSinkLevel = "info"

// logSinkDir is where the persistent sink for a given data root writes. Kept
// separate from installLogSink so the placement rule can be asserted without
// touching a filesystem.
func logSinkDir(dataDir string) string {
	return filepath.Join(dataDir, logDirName)
}

// logSink is the running pipeline plus the directory it resolved to, so a
// caller can name the landing spot in whatever it prints.
type logSink struct {
	pipeline *observe.LogPipeline
	dir      string
}

// logDir returns the log directory the sink writes to.
func (s *logSink) logDir() string {
	if s == nil {
		return ""
	}
	return s.dir
}

// close stops the flusher, flushes and syncs. Not propagating an error: the
// only caller is a process-exit path whose job is to hand the bytes to the
// disk, and the pipeline already refuses later writes internally.
func (s *logSink) close() {
	if s == nil || s.pipeline == nil {
		return
	}
	_ = s.pipeline.Close()
}

// logger returns a writer that goes to the file and ONLY the file. The run
// leg's audit trail uses it (agentRuntime.auditf) because that trail already
// has its own console form, "[audit] ..." on stderr; sending the same sentence
// through the process default as well would put every audit line on the
// terminal twice. With no sink it degrades to the process default, which is
// the state before this ticket: on the console, and nowhere else.
func (s *logSink) logger() *slog.Logger {
	if s == nil || s.pipeline == nil {
		return slog.Default()
	}
	return slog.New(s.pipeline.Handler())
}

// installLogSink starts the rolling JSONL pipeline over the given data root and
// makes it the process-wide slog default, which is the one channel every
// security notice in internal/* already writes to (they use the package-level
// slog functions, and internal/winsec must not import internal/observe - the
// graph runs observe -> secret -> winsec, so that edge would close a cycle).
//
// The default it installs is a fan-out, not a swap. Ticket 89's promise was that
// a cleared authorization is loud, and one of the places it was loud was the
// terminal; routing the log to a file and leaving the terminal silent would buy
// AC#2 by re-creating R-104-5 one screen over. So a record now goes to the
// redacting JSONL writer AND to stderr, which is what an attached console saw
// before this file existed.
//
// An empty dataDir is a refusal, not a fallback: "no data root resolved" must
// never turn into "so write the security log somewhere else".
func installLogSink(dataDir string) (*logSink, error) {
	if dataDir == "" {
		return nil, errors.New("wisp: no data dir resolved, so the log sink has nowhere to write")
	}
	dir := logSinkDir(dataDir)
	p, err := observe.InitLog(observe.LogConfig{Dir: dir, Level: logSinkLevel})
	if err != nil {
		return nil, err
	}
	slog.SetDefault(slog.New(teeHandler{
		// primary is the pipeline's own handler, which is redactHandler on
		// top of the JSON writer: the file side keeps the non-disableable
		// redaction exactly as `wisp slo`'s sink does.
		primary: p.Handler(),
		// mirror carries the record the way Go's stock default logger did, so
		// the console content does not change shape when the sink arrives.
		mirror: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}),
	}))
	// Book the install itself: a reader holding one file off a disk has to be
	// able to tell which data root wrote it, and a record of "the listener was
	// here from this moment" is what makes the ordering claim checkable after
	// the fact instead of only in the source. The threshold attribute is named
	// min_level because "level" is the record's own field - a second one in the
	// same JSON object would be a key collision a parser resolves by taking the
	// last value, which is how the first version of this line read as
	// "level":"info" next to "level":"INFO".
	slog.Info("wisp: persistent log sink installed", "dir", dir, "min_level", logSinkLevel)
	return &logSink{pipeline: p, dir: dir}, nil
}

// teeHandler is the fan-out installLogSink installs: one record, two handlers,
// in that order (file first, console second) so a record that makes the primary
// fail still gets a chance to reach the terminal.
//
// The mirror deliberately carries no redaction of its own. That is not a new
// exposure: before a sink existed at all, every one of these lines went to the
// terminal through Go's stock default logger, unredacted, and this keeps the
// console at exactly that content. Redaction is asserted where the promise is
// made - the file that outlives the process (see D33 and
// internal/observe/redact.go).
type teeHandler struct {
	primary slog.Handler
	mirror  slog.Handler
}

func (t teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return t.primary.Enabled(ctx, level) || t.mirror.Enabled(ctx, level)
}

func (t teeHandler) Handle(ctx context.Context, r slog.Record) error {
	err := t.primary.Handle(ctx, r)
	if t.mirror.Enabled(ctx, r.Level) {
		_ = t.mirror.Handle(ctx, r)
	}
	return err
}

func (t teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return teeHandler{primary: t.primary.WithAttrs(attrs), mirror: t.mirror.WithAttrs(attrs)}
}

func (t teeHandler) WithGroup(name string) slog.Handler {
	return teeHandler{primary: t.primary.WithGroup(name), mirror: t.mirror.WithGroup(name)}
}

var _ slog.Handler = teeHandler{}
