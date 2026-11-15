//go:build windows

package main

// Ticket 46's disk-gate case: observe.NewDiskGatedWriter's production assembly.
//
// BEFORE THIS CASE the constructor had two call sites in the whole tree, both
// inside internal/observe/guards_45.go's own fixtures - so the D42#5 log half
// ("below the threshold log writes stop, the live path is unaffected") was
// implemented, tested and installed nowhere, and 46-s6-acceptance.md §1 AC#2
// booked that as "no production assembly point". The wiring now sits in
// logsink.go's installLogSink (LogConfig.Disk), and this case drives THAT
// function: a tripped guard must HOLD the install record instead of writing it,
// and the same record must land once the volume admits writes again.
//
// Red under each of the ways this could regress:
//   - the guard not assembled (LogConfig.Disk left nil): the install record is
//     in the file while the volume is "full", i.e. the first assertion;
//   - the guard assembled but the gated writer not installed by InitLog: same
//     assertion, because the record reaches the file through no admission gate;
//   - the held bytes dropped rather than replayed: the recovery half goes red
//     (the record never lands).
//
// The positive control (a healthy guard, record lands at once) is what
// separates "the guard held it" from "this sink writes nothing at all".

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// diskProbe46 is a flippable free-space reading: the fixture owns the number so
// "the volume is full" and "the volume recovered" are two states this case can
// walk in one process, with no real disk pressure and no mounted image.
type diskProbe46 struct{ free atomic.Int64 }

func (p *diskProbe46) probe(string) (int64, error) { return p.free.Load(), nil }

func swapDiskGuard46(f func(string) (*observe.DiskGuard, string)) func() {
	prev := diskGuardForDataDir
	diskGuardForDataDir = f
	return func() { diskGuardForDataDir = prev }
}

func readSinkFile46(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("46: read log dir %s: %v", dir, err)
	}
	var b strings.Builder
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "wisp-") || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("46: read %s: %v", e.Name(), err)
		}
		b.Write(raw)
	}
	return b.String()
}

// Test46DiskGuardIsAssembledWithTheRealVolumeProbe is the assembly's own
// reading: the production seam answers a guard, and the guard's status carries
// a real free-space number for the volume the data root lives on.
func Test46DiskGuardIsAssembledWithTheRealVolumeProbe(t *testing.T) {
	dir := t.TempDir()
	guard, note := newDataVolumeDiskGuard(dir)
	if guard == nil {
		t.Fatalf("46: the production disk-guard assembly refused a real data dir (%s): the log writer would go in ungated", note)
	}
	st := guard.Status()
	if st.ProbeError != "" {
		t.Fatalf("46: the platform probe failed on a real directory: %s", st.ProbeError)
	}
	if st.FreeBytes <= 0 {
		t.Fatalf("46: free-space reading = %d bytes, want a real positive number", st.FreeBytes)
	}
	if st.Threshold != observe.DiskGuardThresholdDefault {
		t.Fatalf("46: threshold = %d, want D42#5's %d (this leg must not dial the frozen number)", st.Threshold, observe.DiskGuardThresholdDefault)
	}
	if st.Volume != dir {
		t.Fatalf("46: the guard watches %q, want the data root %q", st.Volume, dir)
	}
	t.Logf("46-READ disk guard: %s", note)
}

// Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink is the end-to-end
// case: installLogSink, tripped guard, held record, recovery, replayed record.
func Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink(t *testing.T) {
	prevDefault := slog.Default()
	t.Cleanup(func() { slog.SetDefault(prevDefault) })

	probe := &diskProbe46{}
	probe.free.Store(1 << 20) // 1MB free: far below D42#5's 200MB threshold
	guard, err := observe.NewDiskGuard(observe.DiskGuardOptions{
		Volume: t.TempDir(),
		Probe:  probe.probe,
	})
	if err != nil {
		t.Fatalf("46: fixture guard: %v", err)
	}
	restore := swapDiskGuard46(func(string) (*observe.DiskGuard, string) {
		return guard, "46 fixture probe (1MB free)"
	})
	t.Cleanup(restore)

	// (a) tripped: the install record must not reach the file.
	trippedDir := t.TempDir()
	sink, err := installLogSink(trippedDir)
	if err != nil {
		t.Fatalf("46: installLogSink (tripped): %v", err)
	}
	t.Cleanup(func() { sink.close() })
	logDir := logSinkDir(trippedDir)
	// The reading that separates "held" from "not written yet": a buffered
	// record and a held record both leave an empty file. HeldStats is the
	// gate's own count, so a pipeline with no gated writer answers ok=false
	// here instead of quietly looking healthy.
	held, lost, _, gated := sink.pipeline.DiskHeldStats()
	if !gated {
		t.Fatalf("46: this pipeline has no gated writer at all (LogConfig.Disk never reached InitLog): nothing is holding bytes, so the log class is ungated")
	}
	if held == 0 {
		t.Fatalf("46: a 1MB volume was refused and nothing is held (held=%d): the refused record was dropped, not parked", held)
	}
	if lost != 0 {
		t.Fatalf("46: the first refused record was lost to the bounded tail (%d bytes lost): the install record is what explains the gap", lost)
	}
	if got := readSinkFile46(t, logDir); strings.Contains(got, "persistent log sink installed") {
		t.Fatalf("46: the install record reached the file while the volume was below the disk guard's threshold - the log writer is not gated (%q)", got)
	}
	if got := guard.Status(); !got.Tripped {
		t.Fatalf("46: the guard did not trip on a 1MB reading: %+v", got)
	}
	slog.Info("46 marker before recovery")

	// (b) recovery: the first ADMITTED write replays the held tail ahead of
	// itself, so both the install record and the new marker land, in order.
	probe.free.Store(2 << 30)
	slog.Info("46 marker after recovery")
	if err := sink.pipeline.FlushNow(); err != nil {
		t.Fatalf("46: flush after recovery: %v", err)
	}
	got := readSinkFile46(t, logDir)
	if !strings.Contains(got, "persistent log sink installed") {
		t.Fatalf("46: the record held while the volume was full was never replayed after recovery: %q", got)
	}
	if !strings.Contains(got, "46 marker after recovery") {
		t.Fatalf("46: the admitted record after recovery is missing: %q", got)
	}
	if !strings.Contains(got, `"disk_gated":true`) {
		t.Fatalf("46: the install record does not say the writer is gated: %q", got)
	}
	if i, j := strings.Index(got, "persistent log sink installed"), strings.Index(got, "46 marker after recovery"); i > j {
		t.Fatalf("46: the replayed tail landed AFTER the record that followed it (timeline broken): %q", got)
	}
	if held, _, replayed, _ := sink.pipeline.DiskHeldStats(); held != 0 || replayed == 0 {
		t.Fatalf("46: after recovery the gate holds %d byte(s) and counts %d replayed: the parked tail never went out", held, replayed)
	}

	// Positive control: with a healthy guard the same record lands at once, so
	// (a) above is "the guard held it" and not "this sink writes nothing".
	probe.free.Store(2 << 30)
	healthyDir := t.TempDir()
	healthy, err := installLogSink(healthyDir)
	if err != nil {
		t.Fatalf("46: installLogSink (healthy): %v", err)
	}
	t.Cleanup(func() { healthy.close() })
	// FlushNow is the control's reading step, not a second feature: without it
	// the record is still in the rolling writer's buffer and the file reads
	// empty for every sink, gated or not.
	if err := healthy.pipeline.FlushNow(); err != nil {
		t.Fatalf("46: control flush: %v", err)
	}
	if got := readSinkFile46(t, logSinkDir(healthyDir)); !strings.Contains(got, "persistent log sink installed") {
		t.Fatalf("46: control: a healthy guard did not let the install record land: %q", got)
	}
}
