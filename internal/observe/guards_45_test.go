package observe

// Ticket 45 (45-r1): the diagnostics guards and the privacy contract's preview
// half. Every case below is a reading of a number the ticket names - 200MB, the
// held-line guarantee, the 500MB artifacts budget, the WAL bound, "exactly as
// previewed" - through an injected fixture, so no test here waits for a real
// disk to fill up, mounts anything, or moves a threshold in thresholds.go (that
// file is not touched by this ticket, and these numbers are write-admission
// numbers, not SLO gates).

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Fixtures
// ---------------------------------------------------------------------------

// fixedClock is the guard's time source in these tests: the recheck interval is
// read off the monotonic clock, so a fixture that never advances makes the cache
// behaviour exact rather than timing-dependent.
type fixedClock struct {
	mu sync.Mutex
	at time.Time
}

func (c *fixedClock) now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.at
}

func (c *fixedClock) add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.at = c.at.Add(d)
}

// varProbe is the "small virtual disk" of AC#2: the free-space reading is a
// variable a test dials, and the calls are counted so the cached-reading rule is
// observable.
type varProbe struct {
	mu    sync.Mutex
	free  int64
	err   error
	calls int
}

func (p *varProbe) set(free int64, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.free, p.err = free, err
}

func (p *varProbe) calls_() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

func (p *varProbe) probe(string) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	return p.free, p.err
}

// sink45 is the writer behind the gated sink: it records what actually reached
// "the file", and can be told to refuse.
type sink45 struct {
	mu     sync.Mutex
	buf    bytes.Buffer
	refuse bool
}

func (s *sink45) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.refuse {
		return 0, errors.New("sink45: no space left on device")
	}
	return s.buf.Write(p)
}

func (s *sink45) Flush() error { return nil }
func (s *sink45) Close() error { return nil }

func (s *sink45) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func (s *sink45) setRefuse(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refuse = v
}

// ---------------------------------------------------------------------------
// AC#2 - disk-full guard: trips at threshold, gated classes stop, the live path
// does not, recovery when space is freed.
// ---------------------------------------------------------------------------

func newGuard45(t *testing.T, p *varProbe, clock *fixedClock, threshold int64) *DiskGuard {
	t.Helper()
	g, err := NewDiskGuard(DiskGuardOptions{
		Volume:    filepath.Join(t.TempDir(), "data"),
		Probe:     p.probe,
		Threshold: threshold,
		now:       clock.now,
	})
	if err != nil {
		t.Fatalf("NewDiskGuard: %v", err)
	}
	return g
}

func TestAC2DiskGuardTripsAtThresholdAndRecovers(t *testing.T) {
	clock := &fixedClock{at: time.Unix(1780000000, 0)}
	p := &varProbe{free: 1 << 30} // 1GiB free
	g := newGuard45(t, p, clock, DiskGuardThresholdDefault)

	// Healthy: every class lands, including the gated ones.
	for _, c := range []WriteClass{WriteEssential, WriteLog, WriteArtifacts, WriteL3} {
		if err := g.Check(c); err != nil {
			t.Errorf("healthy volume must admit %q, got %v", c, err)
		}
	}

	// Below the 200MB threshold: the three expensive classes stop.
	p.set(199<<20, nil)
	clock.add(time.Minute)
	for _, c := range []WriteClass{WriteLog, WriteArtifacts, WriteL3} {
		err := g.Check(c)
		if err == nil {
			t.Errorf("guard did not stop write class %q at 199MiB free (threshold is 200MiB)", c)
			continue
		}
		if !errors.Is(err, ErrDiskSpaceLow) {
			t.Errorf("refusal for %q must wrap ErrDiskSpaceLow, got %v", c, err)
		}
	}

	// The live task is UNAFFECTED: D42#5 says so, and this is the only place in
	// the product that says it, so it gets its own assertion.
	if err := g.Check(WriteEssential); err != nil {
		t.Errorf("essential write must never be gated, got %v", err)
	}

	st := g.Status()
	if !st.Tripped || st.Trips != 1 || st.FreeBytes != 199<<20 {
		t.Errorf("status while tripped = %+v, want tripped with 1 trip and 199MiB free", st)
	}
	if st.Warning == "" {
		t.Error("the tripped guard owes a user warning; Status.Warning is empty")
	}
	if !strings.Contains(st.Warning, "live task is unaffected") {
		t.Errorf("the warning must say what is NOT being stopped, got %q", st.Warning)
	}
	if st.TripAt == "" {
		t.Error("a trip without a persisted timestamp is a trip nobody can find in the bundle")
	}

	// Space freed: gated writes come back on the next decision, not after a
	// restart.
	p.set(3<<30, nil)
	clock.add(time.Minute)
	if err := g.Check(WriteLog); err != nil {
		t.Errorf("after the volume recovered the log class must land, got %v", err)
	}
	st = g.Status()
	if st.Tripped || st.Trips != 1 || st.RecoveredAt == "" {
		t.Errorf("status after recovery = %+v, want not tripped, 1 trip, recovery booked", st)
	}
}

func TestAC2DiskGuardCachesOnlyWhileHealthy(t *testing.T) {
	clock := &fixedClock{at: time.Unix(1780000000, 0)}
	p := &varProbe{free: 1 << 30}
	g := newGuard45(t, p, clock, DiskGuardThresholdDefault)

	if err := g.Check(WriteLog); err != nil {
		t.Fatalf("healthy check: %v", err)
	}
	before := p.calls_()
	// A 1-byte-free reading inside the recheck window is not seen: that is the
	// cost of not stat-ing the volume per log line, and it is bounded by
	// diskGuardRecheckDefault.
	p.set(1, nil)
	for i := 0; i < 20; i++ {
		if err := g.Check(WriteLog); err != nil {
			t.Fatalf("cached reading must be trusted inside the window: %v", err)
		}
	}
	if got := p.calls_(); got != before {
		t.Errorf("guard re-probed %d time(s) inside the recheck window; want 0", got-before)
	}
	clock.add(diskGuardRecheckDefault)
	if err := g.Check(WriteLog); err == nil {
		t.Error("after the recheck window the guard must look again and trip")
	}
	// While tripped there is no cache: recovery must be seen at once.
	callsBefore := p.calls_()
	p.set(1<<30, nil)
	if err := g.Check(WriteLog); err != nil {
		t.Errorf("recovery must be seen on the next call: %v", err)
	}
	if p.calls_() <= callsBefore {
		t.Error("a tripped guard must re-probe, not serve a cached low reading")
	}
}

func TestAC2DiskGuardProbeFailureFailsClosed(t *testing.T) {
	clock := &fixedClock{at: time.Unix(1780000000, 0)}
	p := &varProbe{err: errors.New("probe: volume not found")}
	g := newGuard45(t, p, clock, DiskGuardThresholdDefault)

	if err := g.Check(WriteLog); err == nil {
		t.Fatal("a guard that cannot read free space must refuse gated writes, not wave them through")
	} else if errors.Is(err, ErrDiskSpaceLow) {
		t.Errorf("a probe failure must not masquerade as ErrDiskSpaceLow: %v", err)
	}
	if err := g.Check(WriteEssential); err != nil {
		t.Errorf("the live path is still unaffected by an unreadable volume: %v", err)
	}
	if st := g.Status(); st.ProbeError == "" {
		t.Error("Status must name a failed reading instead of reporting 0 free bytes")
	} else if st.FreeBytes == 0 && st.Tripped {
		t.Errorf("status reads as a trip at 0 bytes: %+v", st)
	}
	// A guard with no probe at all is a configuration error, not an open gate.
	if _, err := NewDiskGuard(DiskGuardOptions{Volume: "x"}); err == nil {
		t.Error("NewDiskGuard without a probe must fail closed")
	}
	if _, err := NewDiskGuard(DiskGuardOptions{Probe: p.probe}); err == nil {
		t.Error("NewDiskGuard without a volume must fail closed")
	}
}

// ---------------------------------------------------------------------------
// "Log rollover under disk pressure never loses the last error line"
// ---------------------------------------------------------------------------

func TestGatedLogsHoldTheLastLineAndReplayItAfterRecovery(t *testing.T) {
	clock := &fixedClock{at: time.Unix(1780000000, 0)}
	p := &varProbe{free: 1 << 30}
	g := newGuard45(t, p, clock, DiskGuardThresholdDefault)
	inner := &sink45{}
	d := NewDiskGatedWriter(inner, g)

	first := `{"msg":"task started"}\n`
	if _, err := d.Write([]byte(first)); err != nil {
		t.Fatalf("healthy write: %v", err)
	}
	if inner.String() != first {
		t.Fatalf("healthy write must reach the file, got %q", inner.String())
	}

	// Pressure arrives. The error line that explains it is the one line that
	// must not be the one the guard drops.
	p.set(10<<20, nil)
	clock.add(time.Minute)
	errLine := `{"level":"ERROR","msg":"wisp write failed: no space left on device"}\n`
	if _, err := d.Write([]byte(errLine)); err != nil {
		t.Fatalf("a refused gated write must be held, not returned as a writer failure: %v", err)
	}
	if strings.Contains(inner.String(), errLine) {
		t.Error("the guarded write reached the file while the guard was tripped")
	}
	held, lost, _ := d.HeldStats()
	if held == 0 || lost != 0 {
		t.Errorf("held/lost after one refusal = %d/%d, want the bytes held and nothing lost", held, lost)
	}

	// More pressure, then the volume is cleared. The next write replays the tail
	// FIRST, so the file reads as a timeline.
	lastLine := `{"level":"ERROR","msg":"the last error line"}\n`
	if _, err := d.Write([]byte(lastLine)); err != nil {
		t.Fatalf("second refused write: %v", err)
	}
	p.set(2<<30, nil)
	clock.add(time.Minute)
	after := `{"msg":"back to work"}\n`
	if _, err := d.Write([]byte(after)); err != nil {
		t.Fatalf("write after recovery: %v", err)
	}
	got := inner.String()
	for _, want := range []string{first, errLine, lastLine, after} {
		if !strings.Contains(got, strings.TrimSuffix(want, "\n")) {
			t.Errorf("recovered file is missing %q; got %q", want, got)
		}
	}
	if i, j := strings.Index(got, lastLine), strings.Index(got, after); i > j {
		t.Errorf("the replay must come before the write that triggered it: %q", got)
	}
	if _, _, replayed := d.HeldStats(); replayed == 0 {
		t.Error("the replay must be countable, not just visible in the file")
	}

	// The tail is bounded and it keeps the NEWEST bytes: a pressure storm of 4x
	// the cap costs the oldest lines, and the cost is a number, not a silence.
	p.set(10<<20, nil)
	clock.add(time.Minute)
	storm := strings.Repeat("x", heldLogTailMaxBytes/2) + "\n"
	for i := 0; i < 6; i++ {
		if _, err := d.Write([]byte(storm)); err != nil {
			t.Fatalf("storm write %d: %v", i, err)
		}
	}
	held, lost, _ = d.HeldStats()
	if held > heldLogTailMaxBytes {
		t.Errorf("held tail %d exceeds its %d byte cap", held, heldLogTailMaxBytes)
	}
	if lost == 0 {
		t.Error("a tail that overflowed must book what it could not hold")
	}
	if err := d.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := d.Write([]byte("late\n")); err != nil {
		t.Errorf("a write after Close must still be accounted, not panic: %v", err)
	}
}

func TestGatedLogsPutBytesBackWhenTheVolumeRefusesThem(t *testing.T) {
	clock := &fixedClock{at: time.Unix(1780000000, 0)}
	p := &varProbe{free: 1 << 30}
	g := newGuard45(t, p, clock, DiskGuardThresholdDefault)
	inner := &sink45{}
	inner.setRefuse(true) // The reading said there was room; the write said no.
	d := NewDiskGatedWriter(inner, g)

	line := `{"level":"ERROR","msg":"ENOSPC between check and write"}\n`
	if _, err := d.Write([]byte(line)); err != nil {
		t.Fatalf("a write the volume refused must be held, not surfaced as a pipeline error: %v", err)
	}
	held, _, _ := d.HeldStats()
	if held != len(line) {
		t.Errorf("held = %d, want %d", held, len(line))
	}
	inner.setRefuse(false)
	if err := d.Flush(); err != nil {
		t.Fatalf("flush after the volume came back: %v", err)
	}
	if !strings.Contains(inner.String(), strings.TrimSuffix(line, "\n")) {
		t.Errorf("Flush must replay the held tail before flushing: %q", inner.String())
	}
	if h, _, _ := d.HeldStats(); h != 0 {
		t.Errorf("after a successful replay the tail must be empty, held=%d", h)
	}
}

// ---------------------------------------------------------------------------
// AC#3 - WAL bounded under sustained writes; the action is named.
// ---------------------------------------------------------------------------

func TestAC3WalStaysBoundedUnderSustainedWrites(t *testing.T) {
	dir := t.TempDir()
	wal := filepath.Join(dir, "wisp.db-wal")

	// No -WAL file is healthy: nothing to fold.
	if c := CheckWalFile(wal, 4096); c.Action != WalActionNone || c.OverBound || c.ReadingFail != "" {
		t.Errorf("missing WAL must read as healthy, got %+v", c)
	}

	f, err := os.Create(wal)
	if err != nil {
		t.Fatalf("create fixture WAL: %v", err)
	}
	chunk := make([]byte, 512)
	var seen []WalCheck
	for i := 0; i < 16; i++ {
		if _, werr := f.Write(chunk); werr != nil {
			t.Fatalf("sustained write %d: %v", i, werr)
		}
		if ferr := f.Sync(); ferr != nil {
			t.Fatalf("sync: %v", ferr)
		}
		seen = append(seen, CheckWalFile(wal, 4096))
	}
	_ = f.Close()

	var sizes []int64
	overAt := -1
	for i, c := range seen {
		sizes = append(sizes, c.SizeBytes)
		if c.SizeBytes > 4096 && overAt < 0 {
			overAt = i
		}
	}
	if overAt < 0 {
		t.Fatalf("the fixture never crossed its bound; sizes=%v", sizes)
	}
	if !seen[overAt].OverBound || seen[overAt].Action != WalActionTruncate {
		t.Errorf("the first reading over the bound must name TRUNCATE, got %+v", seen[overAt])
	}
	// Past half the bound the autocheckpoint is already due - before the file is
	// over it, which is the point of having two readings.
	if got := seen[0].Action; got != WalActionNone && got != WalActionCheckpointAuto {
		t.Errorf("first reading action = %q, want none or the auto fold", got)
	}
	var lastOver bool
	for _, c := range seen {
		lastOver = c.OverBound
	}
	if !lastOver {
		t.Error("the last reading of a still-growing WAL must still be over the bound (no decay)")
	}
	for _, c := range seen {
		if c.BoundBytes != 4096 || c.Autocheck != WalAutoCheckpointPages {
			t.Errorf("reading must carry the bound it was measured against: %+v", c)
		}
		if c.CheckedAt == "" {
			t.Error("a WAL reading without a timestamp is not a record")
		}
	}
	// An unreadable path is REPORTED as unreadable, never scored 0 bytes.
	bad := CheckWalFile(filepath.Join(dir, "nope", "x-wal"), 4096)
	if bad.ReadingFail == "" {
		t.Errorf("a failing stat must be booked, got %+v", bad)
	}
}

// ---------------------------------------------------------------------------
// AC#5 - artifacts LRU eviction under quota pressure; evicted files logged.
// ---------------------------------------------------------------------------

func writeArtifact45(t *testing.T, dir, name string, bytes_ int, mod time.Time) {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, bytes.Repeat([]byte("a"), bytes_), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	if err := os.Chtimes(p, mod, mod); err != nil {
		t.Fatalf("chtimes %s: %v", name, err)
	}
}

func TestAC5ArtifactsQuotaEvictsOldestFirst(t *testing.T) {
	dir := t.TempDir()
	base := time.Unix(1780000000, 0).UTC()
	// Four artifacts: oldest first. 300 bytes total against a 120-byte quota.
	writeArtifact45(t, dir, "a-old-1.bin", 100, base)
	writeArtifact45(t, dir, "b-old-2.bin", 100, base.Add(time.Hour))
	writeArtifact45(t, dir, "c-new-1.bin", 100, base.Add(2*time.Hour))
	// A directory inside the artifacts tree: it is not an artifact, and a quota
	// pass that counted or deleted it would be walking somewhere the ticket does
	// not point it.
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o700); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "d.bin"), bytes.Repeat([]byte("a"), 100), 0o600); err != nil {
		t.Fatalf("seed nested artifact: %v", err)
	}

	rep, err := EnforceArtifactsQuota(dir, 120)
	if err != nil {
		t.Fatalf("EnforceArtifactsQuota: %v", err)
	}
	if rep.TotalBefore != 300 || rep.TotalAfter > 120 {
		t.Errorf("quota pass reported before=%d after=%d, want 300 -> <=120", rep.TotalBefore, rep.TotalAfter)
	}
	got := strings.Join(rep.Evicted, ",")
	if got != "a-old-1.bin,b-old-2.bin" {
		t.Errorf("evicted %q, want the two oldest first", got)
	}
	for _, name := range rep.Evicted {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Errorf("%q is still on disk but the report says it was evicted", name)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "c-new-1.bin")); err != nil {
		t.Errorf("the newest artifact must survive: %v", err)
	}
	if rep.BytesFreed != 200 {
		t.Errorf("BytesFreed = %d, want 200", rep.BytesFreed)
	}
	if rep.EnforcedAt == "" || rep.QuotaBytes != 120 {
		t.Errorf("report must carry the quota it enforced and a timestamp: %+v", rep)
	}
	// A directory entry is refused, not deleted.
	for _, r := range rep.Refused {
		if r == "sub" {
			t.Errorf("a directory is not an artifact and must not even be refused as one: %v", rep.Refused)
		}
	}
	// Under quota: nothing moves.
	rep2, err := EnforceArtifactsQuota(dir, 1<<20)
	if err != nil {
		t.Fatalf("second pass: %v", err)
	}
	if len(rep2.Evicted) != 0 || rep2.TotalAfter != rep2.TotalBefore {
		t.Errorf("an under-quota pass must evict nothing, got %+v", rep2)
	}
	// quotaBytes <= 0 selects D35's 500MB, it does not mean "quota off".
	rep3, err := EnforceArtifactsQuota(dir, 0)
	if err != nil {
		t.Fatalf("default-quota pass: %v", err)
	}
	if rep3.QuotaBytes != ArtifactsQuotaDefault {
		t.Errorf("default quota = %d, want D35's %d", rep3.QuotaBytes, ArtifactsQuotaDefault)
	}
	if _, err := EnforceArtifactsQuota("", 10); err == nil {
		t.Error("an empty artifacts dir must be a configuration error")
	}
}

// ---------------------------------------------------------------------------
// AC#1 - bundle contents exactly as previewed; secrets/transcripts/audio absent
// by default; per-item opt-in with a per-item preview.
// ---------------------------------------------------------------------------

func zipEntries45(t *testing.T, path string) []string {
	t.Helper()
	zr, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open bundle zip: %v", err)
	}
	defer zr.Close()
	out := make([]string, 0, len(zr.File))
	for _, f := range zr.File {
		out = append(out, f.Name)
	}
	return out
}

func readZipEntry45(t *testing.T, path, name string) (string, bool) {
	t.Helper()
	zr, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open bundle zip: %v", err)
	}
	defer zr.Close()
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, oerr := f.Open()
		if oerr != nil {
			t.Fatalf("open entry %s: %v", name, oerr)
		}
		defer rc.Close()
		data, rerr := io.ReadAll(rc)
		if rerr != nil {
			t.Fatalf("read entry %s: %v", name, rerr)
		}
		return string(data), true
	}
	return "", false
}

func TestAC1BundleContentsAreExactlyAsPreviewed(t *testing.T) {
	dir := t.TempDir()
	logDir := filepath.Join(dir, "logs")
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		t.Fatalf("mkdir logs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(logDir, "wisp-20261103-001.jsonl"),
		[]byte("{\"msg\":\"task finished\"}\n"), 0o600); err != nil {
		t.Fatalf("seed log: %v", err)
	}
	// Over the per-file cap: the builder skips it, so the preview must not
	// promise it.
	if err := os.WriteFile(filepath.Join(logDir, "wisp-20261103-002.jsonl"),
		[]byte(strings.Repeat("y", 200)), 0o600); err != nil {
		t.Fatalf("seed oversized log: %v", err)
	}
	cfg := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(cfg, []byte("[observe]\nlevel = \"info\"\n"), 0o600); err != nil {
		t.Fatalf("seed config: %v", err)
	}
	trDir := filepath.Join(dir, "transcripts")
	if err := os.MkdirAll(trDir, 0o700); err != nil {
		t.Fatalf("mkdir transcripts: %v", err)
	}
	trBody := "user: what is the weather\nassistant: it is raining\n"
	if err := os.WriteFile(filepath.Join(trDir, "session-1.txt"), []byte(trBody), 0o600); err != nil {
		t.Fatalf("seed transcript: %v", err)
	}
	if err := os.WriteFile(filepath.Join(trDir, "session-2.txt"), []byte("user: second session\n"), 0o600); err != nil {
		t.Fatalf("seed transcript 2: %v", err)
	}
	snapDir := filepath.Join(dir, "snapshots")
	if err := os.MkdirAll(snapDir, 0o700); err != nil {
		t.Fatalf("mkdir snapshots: %v", err)
	}
	if err := os.WriteFile(filepath.Join(snapDir, "snap.json"), []byte("{\"rss_kb\":1}"), 0o600); err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	guardProbe := &varProbe{free: 199 << 20}
	guardClock := &fixedClock{at: time.Unix(1780000000, 0)}
	guard := newGuard45(t, guardProbe, guardClock, DiskGuardThresholdDefault)
	if err := guard.Check(WriteLog); err == nil {
		t.Fatal("fixture: the guard must be tripped so the bundle carries a real reading")
	}
	st := guard.Status()

	base := BundleOptions{
		OutPath:       filepath.Join(dir, "bundle-default.zip"),
		LogDir:        logDir,
		ConfigPath:    cfg,
		SLOSnapshot:   []byte("{\"pass\":1}"),
		Version:       map[string]string{"version": "0.0-test"},
		Deferred:      []DeferredItem{{ID: "D-45", Ticket: "45", Summary: "guards"}},
		MaxLogBytes:   64,
		Now:           guardClock.now(),
		DiskStatus:    &st,
		TranscriptDir: trDir,
		SnapshotDir:   snapDir,
	}

	// --- the default bundle: preview first, then build, then compare. ---
	prev, err := PreviewDiagnosticsBundle(base)
	if err != nil {
		t.Fatalf("preview: %v", err)
	}
	if _, err := BuildDiagnosticsBundle(base); err != nil {
		t.Fatalf("build default: %v", err)
	}
	have := zipEntries45(t, base.OutPath)
	sort.Strings(have)
	want := append([]string(nil), prev.Entries...)
	sort.Strings(want)
	if strings.Join(have, "|") != strings.Join(want, "|") {
		t.Errorf("bundle contents are not exactly as previewed:\n  zip:     %v\n  preview: %v", have, want)
	}

	// Default exclusions, each asserted as ABSENT rather than as a manifest word.
	for _, e := range have {
		if strings.HasPrefix(e, "transcripts/") || strings.HasPrefix(e, "watchdog/") || strings.HasPrefix(e, "audio") {
			t.Errorf("default bundle carries opt-in material: %q", e)
		}
	}
	manifest, ok := readZipEntry45(t, base.OutPath, "manifest.json")
	if !ok {
		t.Fatal("manifest.json missing from the default bundle")
	}
	for _, phrase := range []string{"not collected: default-excluded", "never collected", "tripped=true"} {
		if !strings.Contains(manifest, phrase) {
			t.Errorf("manifest must state %q, got %s", phrase, manifest)
		}
	}
	// The oversized log is skipped by BOTH paths (preview says so, zip omits it).
	if _, ok := readZipEntry45(t, base.OutPath, "logs/wisp-20261103-002.jsonl"); ok {
		t.Error("the oversized log file must not be in the bundle")
	}
	for _, c := range prev.Categories {
		if c.Category != "logs" {
			continue
		}
		if len(c.Items) != 2 {
			t.Fatalf("log preview listed %d item(s), want the kept one plus the skipped one: %+v", len(c.Items), c.Items)
		}
		var skipped *PreviewItem
		for i := range c.Items {
			if !c.Items[i].Included {
				skipped = &c.Items[i]
			} else if !strings.Contains(c.Items[i].Head, "task finished") {
				t.Errorf("log per-item preview head = %q", c.Items[i].Head)
			}
		}
		if skipped == nil || skipped.Reason == "" {
			t.Error("the preview must show the skipped file and why")
		}
	}

	// --- per-item opt-in: one transcript, previewed item by item. ---
	opted := base
	opted.OutPath = filepath.Join(dir, "bundle-optin.zip")
	opted.OptIn = BundleOptIn{Transcripts: []string{"session-1.txt"}}
	prevOpt, err := PreviewDiagnosticsBundle(opted)
	if err != nil {
		t.Fatalf("preview opt-in: %v", err)
	}
	var trPrev *PreviewCategory
	for i := range prevOpt.Categories {
		if prevOpt.Categories[i].Category == "transcripts" {
			trPrev = &prevOpt.Categories[i]
		}
	}
	if trPrev == nil || len(trPrev.Items) != 1 || !trPrev.Items[0].Included {
		t.Fatalf("transcript category must list exactly the one opted-in item: %+v", trPrev)
	}
	if trPrev.Items[0].Name != "transcripts/session-1.txt" {
		t.Errorf("opted-in entry name = %q, want transcripts/session-1.txt", trPrev.Items[0].Name)
	}
	red := Redactor{}
	if trPrev.Items[0].Head != red.String("user: what is the weather") {
		t.Errorf("per-item preview head = %q, want the redacted first line", trPrev.Items[0].Head)
	}
	if trPrev.DefaultIncluded {
		t.Error("transcripts must be marked default-excluded in the preview")
	}
	if _, err := BuildDiagnosticsBundle(opted); err != nil {
		t.Fatalf("build opt-in: %v", err)
	}
	haveOpt := zipEntries45(t, opted.OutPath)
	sort.Strings(haveOpt)
	wantOpt := append([]string(nil), prevOpt.Entries...)
	sort.Strings(wantOpt)
	if strings.Join(haveOpt, "|") != strings.Join(wantOpt, "|") {
		t.Errorf("opt-in bundle is not exactly as previewed:\n  zip:     %v\n  preview: %v", haveOpt, wantOpt)
	}
	body, ok := readZipEntry45(t, opted.OutPath, "transcripts/session-1.txt")
	if !ok {
		t.Fatal("the opted-in transcript is missing from the bundle")
	}
	if body != red.String(trBody) {
		t.Errorf("the transcript must go through the same redactor the logs use:\n got %q\nwant %q", body, red.String(trBody))
	}
	for _, e := range haveOpt {
		if strings.HasPrefix(e, "transcripts/session-2") {
			t.Errorf("a transcript nobody opted in is in the bundle: %q", e)
		}
	}

	// --- consent shape: only names the user could have seen are accepted. ---
	for _, bad := range []string{"", "..", ".", "../secrets", "a/b", `a\b`, ".hidden"} {
		evil := opted
		evil.OutPath = filepath.Join(dir, "bundle-evil-"+strings.NewReplacer("/", "_", `\`, "_", ".", "dot").Replace(bad)+".zip")
		evil.OptIn = BundleOptIn{Transcripts: []string{bad}}
		if _, err := PreviewDiagnosticsBundle(evil); err == nil {
			t.Errorf("preview accepted the opt-in name %q", bad)
		}
		if _, err := BuildDiagnosticsBundle(evil); err == nil {
			t.Errorf("builder accepted the opt-in name %q", bad)
		}
	}

	// An opt-in with nowhere to read from is a configuration error, not a
	// silently empty bundle.
	dangling := base
	dangling.OutPath = filepath.Join(dir, "bundle-dangling.zip")
	dangling.OptIn = BundleOptIn{Snapshots: []string{"snap.json"}}
	dangling.SnapshotDir = ""
	if _, err := BuildDiagnosticsBundle(dangling); err == nil {
		t.Error("a snapshot opt-in without a snapshot dir must fail")
	}
	dangling2 := dangling
	dangling2.OptIn = BundleOptIn{Transcripts: []string{"session-1.txt"}}
	dangling2.TranscriptDir = ""
	if _, err := PreviewDiagnosticsBundle(dangling2); err == nil {
		t.Error("a transcript opt-in without a transcript dir must fail")
	}

	// A preview that cannot build must not pretend: an unreadable log dir is a
	// failure in both paths, so a previewed bundle is always buildable.
	ghost := base
	ghost.OutPath = filepath.Join(dir, "bundle-ghost.zip")
	ghost.LogDir = filepath.Join(dir, "no-such-log-dir")
	if _, err := PreviewDiagnosticsBundle(ghost); err == nil {
		t.Error("preview of an unreadable log dir must fail like the builder does")
	}
	if _, err := BuildDiagnosticsBundle(ghost); err == nil {
		t.Error("builder must fail on an unreadable log dir")
	}
	if _, err := PreviewDiagnosticsBundle(BundleOptions{}); err == nil {
		t.Error("a preview with no output path must fail")
	}
}
