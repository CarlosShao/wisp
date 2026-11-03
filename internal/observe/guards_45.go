package observe

// ---------------------------------------------------------------------------
// The disk-full / WAL / artifacts-quota guards (ticket 45, D42#5 + D35
// artifacts quota).
//
// WHAT THESE ARE FOR. Ticket 45's key constraint is that running out of space
// must degrade the expensive writers and touch the live one not at all: "below
// -> stop L3/log/artifacts writes + user warning, main response path
// UNAFFECTED". That decision has three parts, held here where they can be read
// and tested together:
//
//   - DiskGuard decides, per write class, whether this write may land now. It
//     reads free space on the data volume through an injected probe: the cmd
//     wiring owns the platform call (golang.org/x/sys/windows already carries
//     GetDiskFreeSpaceEx and is in go.mod:12), and this package stays
//     dependency-clean - the same reasoning that makes BundleOptions.Version
//     caller-supplied in diagnostics.go. No probe means no verdict, and no
//     verdict means the gated classes are refused (fail-closed), never waved
//     through.
//   - The 200MB number is D42#5's own, carried as DiskGuardThresholdDefault.
//     It is NOT one of the frozen SLO thresholds: thresholds.go is untouched by
//     this file. This is a write-admission threshold, which is why a fixture can
//     dial it and no acceptance gate moves.
//   - A gated log write must not silently eat the record that explains why the
//     disk is full. diskGatedWriter therefore HOLDS the refused bytes (bounded)
//     and replays them once the volume recovers: "log rollover under disk
//     pressure never loses the last error line". A naive refusal drops exactly
//     the lines a reader needs.
//
// The WAL and artifacts halves are the other two bounded-writer claims: the
// runtime re-check that the WAL is still inside its bound (ticket 04 sets the
// startup TRUNCATE + wal_autocheckpoint=1000; this is the reading that notices
// when the file stopped obeying), and the 500MB artifacts LRU whose evictions
// have to be logged to be "visible in diagnostics".
//
// One rule runs through all of it: the guard never logs while holding its lock.
// A warning written here can be routed back into the very pipeline that is
// asking for permission, and a second Check call on the same mutex from the same
// goroutine is a deadlock, not a warning.
// ---------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// WriteClass is what a writer is trying to do. The disk guard gates everything
// except WriteEssential: D42#5 never stops the main response path, so a caller
// asking about WriteEssential gets nil even with zero bytes free.
type WriteClass int

const (
	// WriteEssential is the main response path (live task, user-facing
	// completion). Never gated.
	WriteEssential WriteClass = iota
	// WriteLog is the JSONL log pipeline.
	WriteLog
	// WriteArtifacts is an artifact write (D35 artifacts quota).
	WriteArtifacts
	// WriteL3 is an L3 (highest-cost, non-live) write.
	WriteL3
)

func (c WriteClass) String() string {
	switch c {
	case WriteEssential:
		return "essential"
	case WriteLog:
		return "log"
	case WriteArtifacts:
		return "artifacts"
	case WriteL3:
		return "l3"
	}
	return "unknown"
}

const (
	// DiskGuardThresholdDefault is D42#5's pre-write free-space threshold:
	// 200MB on the data volume.
	DiskGuardThresholdDefault int64 = 200 << 20

	// diskGuardRecheckDefault is how long an "there is room" reading stays
	// trusted before the guard looks again, so a hot log path does not stat the
	// volume per line while things are healthy. A tripped guard re-probes on
	// every call so recovery is picked up at once.
	diskGuardRecheckDefault = 5 * time.Second

	// heldLogTailMaxBytes bounds the refused-but-kept log bytes. It is the same
	// 64KiB the log pipeline buffers (bufferedLogBytes), so the worst case this
	// guard adds to the process is one more log-sized buffer.
	heldLogTailMaxBytes = bufferedLogBytes
)

// ErrDiskSpaceLow is the identity of a refusal made by the disk guard. Callers
// match it with errors.Is; the wording a user sees comes from DiskStatus, not
// from this string.
var ErrDiskSpaceLow = errors.New("observe: data volume below the disk-full guard threshold")

// SpaceProbe reads the free bytes on a volume. It is injected: the platform
// call belongs to the composition root, and the fixture AC#2 asks for ("small
// virtual disk / quota") is a probe, not a mounted disk.
type SpaceProbe func(volumePath string) (freeBytes int64, err error)

// DiskGuardOptions configures NewDiskGuard.
type DiskGuardOptions struct {
	// Volume is the data volume to watch (the directory the data tree lives
	// under; the probe resolves the volume itself).
	Volume string
	// Probe reads free space. Required: a guard that can never read is a guard
	// that never trips.
	Probe SpaceProbe
	// Threshold overrides DiskGuardThresholdDefault (tests dial it).
	Threshold int64
	// RecheckInterval overrides diskGuardRecheckDefault.
	RecheckInterval time.Duration
	// now reads the monotonic clock for the recheck interval - never a
	// wall-clock subtraction (package clock rules, ban #4).
	now func() time.Time
}

// DiskGuard is the pre-write free-space decision for the non-essential write
// classes. Safe for concurrent use.
type DiskGuard struct {
	mu        sync.Mutex
	volume    string
	probe     SpaceProbe
	threshold int64
	recheck   time.Duration
	now       func() time.Time

	tripped   bool
	lastFree  int64
	lastProbe time.Time
	tripAt    time.Time
	recovered time.Time
	trips     int
	warning   string
}

// NewDiskGuard builds the guard; an empty volume or a missing probe is a
// configuration error.
func NewDiskGuard(o DiskGuardOptions) (*DiskGuard, error) {
	if o.Volume == "" {
		return nil, New(ClassConfig, "observe: disk guard volume is empty")
	}
	if o.Probe == nil {
		return nil, New(ClassConfig, "observe: disk guard needs a free-space probe")
	}
	if o.Threshold <= 0 {
		o.Threshold = DiskGuardThresholdDefault
	}
	if o.RecheckInterval <= 0 {
		o.RecheckInterval = diskGuardRecheckDefault
	}
	if o.now == nil {
		o.now = time.Now
	}
	return &DiskGuard{
		volume: o.Volume, probe: o.Probe, threshold: o.Threshold,
		recheck: o.RecheckInterval, now: o.now,
	}, nil
}

// Check reports whether a write of class c may land right now. WriteEssential
// is always allowed. A gated class returns an error wrapping ErrDiskSpaceLow
// when the volume is below threshold - or when the free-space reading failed,
// which is the fail-closed direction.
func (g *DiskGuard) Check(c WriteClass) error {
	if c == WriteEssential {
		return nil
	}
	warning := ""
	g.mu.Lock()
	free, err := g.probeLocked()
	if err != nil {
		g.mu.Unlock()
		return Wrap(ClassResource, err, "observe: disk guard cannot read free space; gated writes refused")
	}
	low := free < g.threshold
	switch {
	case low && !g.tripped:
		g.tripped = true
		g.trips++
		g.tripAt = g.now()
		warning = g.buildWarning(free)
		g.warning = warning
	case low:
		// Still tripped: refresh the reading, so a status line is never a stale
		// number dressed as a live one.
		g.warning = g.buildWarning(free)
	case !low && g.tripped:
		g.tripped = false
		g.recovered = g.now()
		g.warning = ""
	}
	tripped, warned := g.tripped, g.warning
	g.mu.Unlock()

	// Outside the lock: this record may travel through the log pipeline, which
	// asks this same guard for permission.
	if warning != "" {
		slog.Warn(warning)
	} else if !tripped && warned != "" {
		slog.Warn("observe: disk-full guard recovered, gated writes allowed again",
			"volume", g.volume, "free_bytes", free, "threshold_bytes", g.threshold)
	}
	if tripped {
		return fmt.Errorf("%w: volume %s free %d < threshold %d, write class %q stopped",
			ErrDiskSpaceLow, g.volume, free, g.threshold, c.String())
	}
	return nil
}

func (g *DiskGuard) buildWarning(free int64) string {
	return fmt.Sprintf(
		"disk-full guard: %s has %d bytes free, below the %d byte threshold - "+
			"log, artifact and L3 writes are stopped; the live task is unaffected",
		g.volume, free, g.threshold)
}

// probeLocked reads free space, honouring the recheck interval while healthy.
// It emits nothing.
func (g *DiskGuard) probeLocked() (int64, error) {
	if !g.tripped && !g.lastProbe.IsZero() && g.now().Sub(g.lastProbe) < g.recheck {
		return g.lastFree, nil
	}
	free, err := g.probe(g.volume)
	if err != nil {
		return 0, err
	}
	g.lastFree, g.lastProbe = free, g.now()
	return free, nil
}

// DiskStatus is the readable state of the guard: the user warning's numbers and
// the record a diagnostics bundle carries.
type DiskStatus struct {
	Volume      string `json:"volume"`
	Tripped     bool   `json:"tripped"`
	FreeBytes   int64  `json:"free_bytes"`
	Threshold   int64  `json:"threshold_bytes"`
	Trips       int    `json:"trips"`
	Warning     string `json:"warning,omitempty"`
	TripAt      string `json:"trip_at,omitempty"`
	RecoveredAt string `json:"recovered_at,omitempty"`
	ProbeError  string `json:"probe_error,omitempty"`
}

// Status reads the guard without changing its state and without emitting a
// record: a status read must never be able to deadlock the pipeline that calls
// it. ProbeError names a reading that failed instead of reporting 0 free bytes,
// which would read as "infinitely full".
func (g *DiskGuard) Status() DiskStatus {
	g.mu.Lock()
	defer g.mu.Unlock()
	st := DiskStatus{
		Volume: g.volume, Tripped: g.tripped, FreeBytes: g.lastFree,
		Threshold: g.threshold, Trips: g.trips, Warning: g.warning,
	}
	if !g.tripAt.IsZero() {
		st.TripAt = WallTimestampUTC(g.tripAt.UTC())
	}
	if !g.recovered.IsZero() {
		st.RecoveredAt = WallTimestampUTC(g.recovered.UTC())
	}
	free, err := g.probeLocked()
	if err != nil {
		st.ProbeError = err.Error()
	} else {
		st.FreeBytes = free
	}
	return st
}

// ---------------------------------------------------------------------------
// The log-side coordination: refuse, hold, replay.
// ---------------------------------------------------------------------------

// diskGatedWriter is the sink the log pipeline writes through when a DiskGuard
// is installed (LogConfig.Disk). For a refused write it keeps the bytes in a
// bounded tail instead of dropping them, and hands them to the file once the
// volume recovers.
type diskGatedWriter struct {
	w     io.WriteCloser
	guard *DiskGuard

	mu      sync.Mutex
	held    [][]byte
	heldLen int
	lost    int // bytes the tail could not even hold: counted, never silent
	kept    int // bytes replayed to the file after a hold
	closed  bool
}

// NewDiskGatedWriter wraps w so writes are admitted only when the guard allows
// the log class. Exported for the cmd wiring leg; InitLog installs it when
// LogConfig.Disk is set.
func NewDiskGatedWriter(w io.WriteCloser, g *DiskGuard) *diskGatedWriter {
	return &diskGatedWriter{w: w, guard: g}
}

func (d *diskGatedWriter) Write(p []byte) (int, error) {
	if err := d.guard.Check(WriteLog); err != nil {
		d.hold(append([]byte(nil), p...))
		return len(p), nil
	}
	// Room again: replay the held bytes, oldest first, ahead of this write, so
	// the file still reads as a timeline and the last error line is in it.
	if err := d.drain(); err != nil {
		return 0, err
	}
	n, err := d.w.Write(p)
	if err != nil {
		// The volume refused us between the reading and the write - the real
		// ENOSPC shape. The bytes go back on the tail rather than into the void.
		d.hold(append([]byte(nil), p...))
		return len(p), nil
	}
	return n, err
}

// hold appends to the bounded tail, evicting from the FRONT (the oldest refused
// bytes), so the newest lines - including the last error line - are the kept
// ones.
func (d *diskGatedWriter) hold(b []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed || len(b) > heldLogTailMaxBytes {
		d.lost += len(b)
		return
	}
	d.held = append(d.held, b)
	d.heldLen += len(b)
	for d.heldLen > heldLogTailMaxBytes && len(d.held) > 0 {
		d.lost += len(d.held[0])
		d.heldLen -= len(d.held[0])
		d.held = d.held[1:]
	}
}

// drain hands the held bytes to the wrapped writer and empties the tail.
func (d *diskGatedWriter) drain() error {
	d.mu.Lock()
	if len(d.held) == 0 {
		d.mu.Unlock()
		return nil
	}
	take := d.held
	d.held, d.heldLen = nil, 0
	d.mu.Unlock()

	for i, b := range take {
		if _, err := d.w.Write(b); err != nil {
			d.mu.Lock()
			if !d.closed {
				rest := take[i:]
				d.held = append([][]byte(nil), rest...)
				d.heldLen = sumHeld(d.held)
			} else {
				d.lost += sumHeld(take[i:])
			}
			d.mu.Unlock()
			return Wrap(ClassResource, err, "observe: replay held log bytes")
		}
	}
	d.mu.Lock()
	d.kept += sumHeld(take)
	d.mu.Unlock()
	return nil
}

// Flush replays the tail (so a held line never waits past the next flush) and
// then flushes the wrapped writer.
func (d *diskGatedWriter) Flush() error {
	if err := d.drain(); err != nil {
		return err
	}
	if f, ok := d.w.(interface{ Flush() error }); ok {
		return f.Flush()
	}
	return nil
}

// Close drains best-effort and closes the wrapped writer. Bytes still held here
// are a number on the console and in HeldStats, not a silence.
func (d *diskGatedWriter) Close() error {
	_ = d.drain()
	d.mu.Lock()
	held, lost := d.heldLen, d.lost
	d.closed = true
	d.mu.Unlock()
	if held > 0 {
		slog.Warn("observe: log bytes still held at close (disk-full guard still tripped)",
			"held_bytes", held, "lost_bytes", lost)
	}
	return d.w.Close()
}

// HeldStats reports the tail: bytes still waiting, bytes the tail could not
// hold, and bytes already replayed.
func (d *diskGatedWriter) HeldStats() (held, lost, replayed int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.heldLen, d.lost, d.kept
}

func sumHeld(chunks [][]byte) int {
	n := 0
	for _, b := range chunks {
		n += len(b)
	}
	return n
}

var _ io.WriteCloser = (*diskGatedWriter)(nil)

// ---------------------------------------------------------------------------
// The bounded WAL: the runtime re-check (D42#5).
// ---------------------------------------------------------------------------

const (
	// WalAutoCheckpointPages is ticket 04's startup setting (wal_autocheckpoint
	// = 1000 pages); the re-check measures against the same bound.
	WalAutoCheckpointPages = 1000
	// WalPageSizeBytesDefault is SQLite's default page size. The storage leg
	// passes its own if it builds with a different one.
	WalPageSizeBytesDefault = 4096
	// WalMaxBytesDefault is the bounded ceiling: twice the autocheckpoint
	// target, so a WAL that reached it has ignored the checkpoint policy rather
	// than merely grown between two checkpoints.
	WalMaxBytesDefault int64 = int64(WalAutoCheckpointPages) * int64(WalPageSizeBytesDefault) * 2
)

// The actions a reader of the check can take, spelled as the PRAGMA's own
// words: the storage leg runs the SQL, this package names the bound and says
// whether it is broken.
const (
	WalActionNone = "none"
	// WalActionCheckpointAuto is the autocheckpoint the policy already asks
	// for: the file is past half its bound, so the fold is due.
	WalActionCheckpointAuto = "wal_checkpoint(AUTO)"
	// WalActionTruncate is startup's TRUNCATE, re-asked at runtime.
	WalActionTruncate = "wal_checkpoint(TRUNCATE)"
)

// WalCheck is one reading of the WAL file against its bound.
type WalCheck struct {
	Path        string `json:"path"`
	SizeBytes   int64  `json:"size_bytes"`
	BoundBytes  int64  `json:"bound_bytes"`
	Autocheck   int    `json:"autocheckpoint_pages"`
	OverBound   bool   `json:"over_bound"`
	Action      string `json:"action"`
	CheckedAt   string `json:"checked_at"`
	ReadingFail string `json:"reading_failure,omitempty"`
}

// CheckWalFile reads the WAL size against the bound (boundBytes <= 0 selects
// WalMaxBytesDefault). A missing -WAL file is healthy: there is nothing to fold.
// A failing stat is REPORTED as failing, never scored 0 bytes - an unreadable
// WAL is exactly the record a reader must see as unreadable.
func CheckWalFile(walPath string, boundBytes int64) WalCheck {
	if boundBytes <= 0 {
		boundBytes = WalMaxBytesDefault
	}
	c := WalCheck{
		Path: walPath, BoundBytes: boundBytes,
		Autocheck: WalAutoCheckpointPages, Action: WalActionNone,
		CheckedAt: WallTimestampUTC(NowWallUTC()),
	}
	info, err := os.Stat(walPath)
	if err != nil {
		// A genuinely absent -WAL beside an existing data dir is healthy: there is
		// nothing to fold. A path that cannot be reached because a parent is missing is
		// a READING failure, never a clean bill - on Windows both arrive as "not
		// exist", so the parent directory is what separates them.
		if os.IsNotExist(err) && dirReadable(filepath.Dir(walPath)) {
			return c
		}
		c.ReadingFail = err.Error()
		return c
	}
	c.SizeBytes = info.Size()
	switch {
	case c.SizeBytes > boundBytes:
		c.OverBound = true
		c.Action = WalActionTruncate
		slog.Warn("observe: WAL over its bound under sustained writes",
			"path", walPath, "size_bytes", c.SizeBytes, "bound_bytes", boundBytes,
			"action", WalActionTruncate)
	case c.SizeBytes > boundBytes/2:
		// Past the autocheckpoint target: name the fold, do not warn.
		c.Action = WalActionCheckpointAuto
	}
	return c
}

// ---------------------------------------------------------------------------
// Artifacts quota: the LRU eviction point (D35, 500MB).
// ---------------------------------------------------------------------------

// ArtifactsQuotaDefault is D35's artifacts budget: 500MB.
const ArtifactsQuotaDefault int64 = 500 << 20

// ArtifactsEvictionReport is one enforcement pass. It is JSON-shaped so the same
// reading can be attached to a diagnostics bundle (the ticket's "visible in
// diagnostics" claim) and asserted without a log parser.
type ArtifactsEvictionReport struct {
	Dir         string   `json:"dir"`
	QuotaBytes  int64    `json:"quota_bytes"`
	TotalBefore int64    `json:"total_before_bytes"`
	TotalAfter  int64    `json:"total_after_bytes"`
	Evicted     []string `json:"evicted"`
	BytesFreed  int64    `json:"bytes_freed"`
	Refused     []string `json:"refused_entries"`
	EnforcedAt  string   `json:"enforced_at"`
}

// EnforceArtifactsQuota walks dir (non-recursive - the artifacts tree is flat by
// D35) and evicts oldest-first until the total fits quotaBytes
// (quotaBytes <= 0 selects ArtifactsQuotaDefault). Oldest is mtime, not atime:
// NTFS does not reliably update access times, so "last written" is the only
// order this package can defend on the shipped platform. Every eviction is
// logged - a quota that deletes silently is the shape ticket 45 exists to
// prevent - and a file that will not go is booked in Refused, so the report
// still says the quota is over.
func EnforceArtifactsQuota(dir string, quotaBytes int64) (ArtifactsEvictionReport, error) {
	if dir == "" {
		return ArtifactsEvictionReport{}, New(ClassConfig, "observe: artifacts quota dir is empty")
	}
	if quotaBytes <= 0 {
		quotaBytes = ArtifactsQuotaDefault
	}
	rep := ArtifactsEvictionReport{
		Dir: dir, QuotaBytes: quotaBytes, EnforcedAt: WallTimestampUTC(NowWallUTC()),
	}
	type item struct {
		name string
		size int64
		mod  time.Time
	}
	var files []item
	total := int64(0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return rep, Wrap(ClassResource, err, "observe: read artifacts dir for quota")
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			// A directory inside the artifacts tree is not an artifact: it is walked
			// around, not refused. Listing it in Refused would make the report read
			// as "something here was rejected", which is not what happened.
			continue
		}
		if strings.ContainsAny(name, `/\`) || name == ".." {
			// A name carrying a separator cannot come from ReadDir: refusing it is
			// cheaper than reasoning about how it got there.
			rep.Refused = append(rep.Refused, name)
			continue
		}
		info, ierr := e.Info()
		if ierr != nil {
			rep.Refused = append(rep.Refused, name)
			continue
		}
		files = append(files, item{name: name, size: info.Size(), mod: info.ModTime().UTC()})
		total += info.Size()
	}
	rep.TotalBefore = total
	rep.TotalAfter = total
	if total <= quotaBytes {
		return rep, nil
	}
	sort.SliceStable(files, func(i, j int) bool {
		if files[i].mod.Equal(files[j].mod) {
			return files[i].name < files[j].name
		}
		return files[i].mod.Before(files[j].mod)
	})
	for _, f := range files {
		if total <= quotaBytes {
			break
		}
		if rmErr := os.Remove(filepath.Join(dir, f.name)); rmErr != nil {
			rep.Refused = append(rep.Refused, f.name)
			continue
		}
		total -= f.size
		rep.Evicted = append(rep.Evicted, f.name)
		rep.BytesFreed += f.size
		slog.Warn("observe: artifacts quota eviction",
			"file", f.name, "bytes", f.size, "quota_bytes", quotaBytes,
			"total_before_bytes", rep.TotalBefore, "mod_utc", WallTimestampUTC(f.mod))
	}
	rep.TotalAfter = total
	return rep, nil
}

// dirReadable is the WAL check's parent-directory question: does the data // directory the -WAL file should sit in actually exist? Only a "no" here turns // an absent WAL into a reported reading failure.
func dirReadable(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
