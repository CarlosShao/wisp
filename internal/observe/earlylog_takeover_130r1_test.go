package observe

// Ticket 130, second landing (130-r1): the handover of the boot records is a
// duty of the CALLER today - cmd/wisp/logsink.go calls FlushEarlyLogRecords and
// every other leg has to remember to do the same. One leg in the shipped tree
// already does not: `wisp slo` (cmd/wisp/slo_windows.go) opens the same pipeline
// through observe.InitLog and installs it with InstallAsDefault, and its JSONL
// holds zero records while its own console printed the sealing verdict. That is
// the measured proof that "the caller remembers" is the wrong seam, and it is
// what these cases pin: the pipeline takes the boot records over at the instant
// it becomes the process listener, so no leg can forget.
//
// The instant is InstallAsDefault and NOT InitLog, for two reasons that are
// readings rather than preferences:
//
//  1. the reading a construction-time replay would spend is this file's own, not
//     ticket 132's: moving the replay into InitLog turns
//     TestInstallAsDefaultTakesOverTheBootRecord130R1 red on its "after InitLog
//     (and a flush) the boot file holds 0 bytes" assertion. Ticket 132's six
//     cases stay green under the same move (6 PASS / 0 FAIL, measured in
//     docs/evidence/s1/130-v1-verdict.md §6): they stat the file without
//     flushing, so they read the writer's bufio.
//  2. "the pipeline became the listener" is exactly what InstallAsDefault means -
//     it is this package's spelling of slog.SetDefault, the act every mainstream
//     logger treats as the installation (zap.ReplaceGlobals,
//     otel.SetLoggerProvider are the same act).
//
// Every case below replaces the package-level buffer with its own, so it cannot
// read or drain what the rest of the suite left there; see swapEarlyLogs130r1.

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	// takeoverBootMsg130r1 is the sentence internal/risk's init() really emits
	// (internal/winsec/resolve.go). Copied as a literal for the same reason the
	// sibling files copy theirs: a reworded record is a red case here rather than
	// a silent change of subject.
	takeoverBootMsg130r1 = "winsec: sealing path resolver installed"
	// takeoverAfterMsg130r1 is a record the listener makes AFTER it is installed;
	// it is the landmark the boot record has to precede on disk.
	takeoverAfterMsg130r1 = "130r1: listener speaking after the handover"
)

// takeoverRecord130r1 is one JSONL line, decoded only for the fields these cases
// read. ProbesPassed is a pointer so "the attribute is absent" and "the
// attribute is zero" stay two different readings.
type takeoverRecord130r1 struct {
	Time         string `json:"time"`
	Level        string `json:"level"`
	Msg          string `json:"msg"`
	Resolver     string `json:"resolver"`
	ProbesPassed *int   `json:"probes_passed"`
	APIKey       string `json:"api_key"`
	Dir          string `json:"dir"`
}

func (r takeoverRecord130r1) String() string {
	return r.Msg
}

// swapEarlyLogs130r1 gives one case its own boot buffer. The variable, not a
// copy of it, is what InstallAsDefault and FlushEarlyLogRecords read, so
// replacing the pointer is the same seam the production code uses - and the
// restore keeps the package-level buffer (the one init() installed into the
// process default's tee) out of this file's way.
func swapEarlyLogs130r1(t *testing.T) *earlyLogBuffer {
	t.Helper()
	old := earlyLogs
	b := &earlyLogBuffer{}
	earlyLogs = b
	t.Cleanup(func() { earlyLogs = old })
	return b
}

// fileBytes130r1 is the number of bytes in the pipeline's files in dir. Zero is
// the reading ticket 132 takes at the construction instant, so it is stated as
// bytes rather than as "no records": a pipeline that opened a file and wrote
// nothing has 0 bytes, and that is the distinction this number makes.
func fileBytes130r1(t *testing.T, dir string) int64 {
	t.Helper()
	var n int64
	for _, name := range logFileNames(t, dir) {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("stat %s: %v", name, err)
		}
		n += info.Size()
	}
	return n
}

// records130r1 reads every record of every pipeline file in dir, in file order.
func records130r1(t *testing.T, dir string) []takeoverRecord130r1 {
	t.Helper()
	var out []takeoverRecord130r1
	for _, name := range logFileNames(t, dir) {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var rec takeoverRecord130r1
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("line %q in %s is not a JSON object: %v", line, name, err)
			}
			out = append(out, rec)
		}
	}
	return out
}

func msgs130r1(recs []takeoverRecord130r1) []string {
	out := make([]string, 0, len(recs))
	for _, r := range recs {
		out = append(out, r.Msg)
	}
	return out
}

// newTakeoverPipeline130r1 makes a pipeline over a fresh temp dir and restores
// the process default on the way out: InstallAsDefault really does install it,
// and a closed pipeline left as the default would silently drop the next case's
// records (logging_test.go's cleanup does the same thing for the same reason).
func newTakeoverPipeline130r1(t *testing.T) (*LogPipeline, string) {
	t.Helper()
	dir := t.TempDir()
	prev := slog.Default()
	p, err := InitLogWithRegistry(LogConfig{Dir: dir, Level: "info"}, NewRegistry())
	if err != nil {
		t.Fatalf("InitLogWithRegistry over %s: %v", dir, err)
	}
	// Registered first, so it runs LAST: the sink has to be closed only after
	// the process default no longer points at it.
	t.Cleanup(func() { _ = p.Close() })
	t.Cleanup(func() { slog.SetDefault(prev) })
	return p, dir
}

// TestInstallAsDefaultTakesOverTheBootRecord130R1 is the headline and the case
// that is red before this landing: a record kept by the early buffer reaches the
// file when the pipeline becomes the listener, ahead of the listener's own first
// record, with the stamp it was spoken with.
//
// The three assertions are separate facts on purpose:
//   - construction writes nothing (ticket 132's seal reading is taken there, and
//     it stays taken: the handover is not moved into InitLog);
//   - the record is IN the file, first, with its attributes - not just "somewhere
//     in the process";
//   - its stamp is older than the listener's own record: a handover that
//     re-dates the record says the listener decided to speak, which is the fact
//     this ticket exists to distinguish from "the boot spoke".
func TestInstallAsDefaultTakesOverTheBootRecord130R1(t *testing.T) {
	b := swapEarlyLogs130r1(t)
	emit130(b, takeoverBootMsg130r1,
		slog.String("resolver", "risk.c26Pipeline"),
		slog.Int("probes_passed", 2))
	if n := b.held130(); n != 1 {
		t.Fatalf("fixture holds %d record(s), want 1: the case would read a buffer it did not fill", n)
	}

	p, dir := newTakeoverPipeline130r1(t)
	// The size is read AFTER a flush, and that is not tidiness: a record handed
	// to the pipeline sits in the writer's 64 KiB buffer, so a bare stat reports
	// 0 bytes whether or not a replay happened. The first draft of this case read
	// the raw size and stayed GREEN when the hop was moved into InitLog - it was
	// measuring the buffer, not the handover. FlushNow is what makes the number
	// stand for "the pipeline held no boot record at construction".
	if err := p.FlushNow(); err != nil {
		t.Fatalf("FlushNow right after InitLog: %v", err)
	}
	if got := fileBytes130r1(t, dir); got != 0 {
		t.Fatalf("after InitLog (and a flush) the boot file holds %d bytes, want 0: construction is not the "+
			"handover instant - a replay moved here is caught by this case's own flushed reading, not by "+
			"ticket 132's six, which stay green under the same move (6 PASS / 0 FAIL, "+
			"docs/evidence/s1/130-v1-verdict.md section 6): the one case there that reads the boot file's "+
			"size at construction stats it without flushing, so it reads the writer's bufio rather than the "+
			"handover, and the other five never take that reading", got)
	}

	p.InstallAsDefault()
	slog.Info(takeoverAfterMsg130r1, "dir", dir)
	if err := p.FlushNow(); err != nil {
		t.Fatalf("FlushNow: %v", err)
	}

	recs := records130r1(t, dir)
	if len(recs) < 2 {
		t.Fatalf("the pipeline's file holds %d record(s) (%v), want the boot record and then the "+
			"listener's own: a sink installed after the first event loses exactly the record it exists "+
			"to keep", len(recs), msgs130r1(recs))
	}
	if recs[0].Msg != takeoverBootMsg130r1 {
		t.Fatalf("record 0 = %q, want the boot record %q (records: %v): the handover did not happen, "+
			"so the record made before the listener existed is still console-only", recs[0].Msg,
			takeoverBootMsg130r1, msgs130r1(recs))
	}
	if recs[0].Resolver != "risk.c26Pipeline" {
		t.Errorf("the taken-over record carries resolver = %q, want risk.c26Pipeline: its attributes did "+
			"not survive the handover: %v", recs[0].Resolver, recs[0])
	}
	if recs[0].ProbesPassed == nil || *recs[0].ProbesPassed != 2 {
		t.Errorf("the taken-over record carries probes_passed = %v, want 2: %v", recs[0].ProbesPassed, recs[0])
	}
	if recs[0].Level != "INFO" {
		t.Errorf("the taken-over record level = %q, want INFO", recs[0].Level)
	}
	if recs[1].Msg != takeoverAfterMsg130r1 {
		t.Errorf("record 1 = %q, want the listener's own record %q: the boot records go first in the file "+
			"(records: %v)", recs[1].Msg, takeoverAfterMsg130r1, msgs130r1(recs))
	}
	bootTime, bootErr := time.Parse(time.RFC3339Nano, recs[0].Time)
	afterTime, afterErr := time.Parse(time.RFC3339Nano, recs[1].Time)
	if bootErr != nil || afterErr != nil {
		t.Fatalf("stamps are not the RFC3339Nano form the pipeline writes: %q (%v), %q (%v)",
			recs[0].Time, bootErr, recs[1].Time, afterErr)
	}
	if !bootTime.Before(afterTime) {
		t.Errorf("the taken-over stamp %s is not older than the listener's own %s: a handover that "+
			"re-dates the record loses the one thing it was kept for", bootTime, afterTime)
	}
	if n := b.held130(); n != 0 {
		t.Errorf("the buffer still holds %d record(s) after the handover, want 0: the same record must not "+
			"be replayable twice", n)
	}
}

// TestHandoverBelongsToTheFirstListener130R1 pins the one-shot half: a second
// pipeline installed later gets nothing, and the first one's file keeps exactly
// one copy. "Each pipeline replays the buffer" would put the same boot verdict
// into every file a process opens, which reads as several boots.
func TestHandoverBelongsToTheFirstListener130R1(t *testing.T) {
	b := swapEarlyLogs130r1(t)
	emit130(b, takeoverBootMsg130r1, slog.String("resolver", "risk.c26Pipeline"), slog.Int("probes_passed", 2))

	first, firstDir := newTakeoverPipeline130r1(t)
	first.InstallAsDefault()
	if err := first.FlushNow(); err != nil {
		t.Fatalf("first FlushNow: %v", err)
	}

	second, secondDir := newTakeoverPipeline130r1(t)
	second.InstallAsDefault()
	slog.Info(takeoverAfterMsg130r1)
	if err := second.FlushNow(); err != nil {
		t.Fatalf("second FlushNow: %v", err)
	}

	firstRecs := records130r1(t, firstDir)
	if n := countMsg130r1(firstRecs, takeoverBootMsg130r1); n != 1 {
		t.Errorf("the first listener's file holds the boot record %d time(s), want exactly 1 (records: %v)",
			n, msgs130r1(firstRecs))
	}
	secondRecs := records130r1(t, secondDir)
	if n := countMsg130r1(secondRecs, takeoverBootMsg130r1); n != 0 {
		t.Errorf("the second listener's file holds the boot record %d time(s), want 0: the boot belongs to "+
			"the listener that arrived first, and a copy per pipeline reads as a boot per pipeline "+
			"(records: %v)", n, msgs130r1(secondRecs))
	}
}

// TestInstallThenFlushDeliversOnce130R1 is the "two blessed paths, one record"
// case: a caller that installs the pipeline AND then asks for the early records
// by hand (the shape cmd/wisp/logsink.go would take if it ever called
// InstallAsDefault) must not get a second copy, and the hand flush has to report
// 0 because there was nothing left for it to do.
func TestInstallThenFlushDeliversOnce130R1(t *testing.T) {
	b := swapEarlyLogs130r1(t)
	emit130(b, takeoverBootMsg130r1, slog.String("resolver", "risk.c26Pipeline"), slog.Int("probes_passed", 2))

	p, dir := newTakeoverPipeline130r1(t)
	p.InstallAsDefault()
	flushed, dropped := FlushEarlyLogRecords(p.Handler())
	if flushed != 0 || dropped != 0 {
		t.Errorf("FlushEarlyLogRecords after InstallAsDefault returned %d/%d, want 0/0: the pipeline already "+
			"took the records over, and a number here claims a second delivery", flushed, dropped)
	}
	if err := p.FlushNow(); err != nil {
		t.Fatalf("FlushNow: %v", err)
	}
	if n := countMsg130r1(records130r1(t, dir), takeoverBootMsg130r1); n != 1 {
		t.Errorf("the boot record is in the file %d time(s), want exactly 1: the two handover paths must "+
			"deliver one copy, not one each", n)
	}
}

// TestHandoverKeepsTheRedaction130R1 pins that the new path into the file is the
// REDACTING handler and not a shortcut around it: the handover is a write path
// added to a file that outlives the process, which is exactly where a bypass
// would be invisible (the console copy is unredacted by design, so "it looked
// fine on the terminal" proves nothing about the file).
func TestHandoverKeepsTheRedaction130R1(t *testing.T) {
	const seeded = "sk-this-is-a-long-seeded-key-9999"
	b := swapEarlyLogs130r1(t)
	emit130(b, takeoverBootMsg130r1, slog.String("api_key", seeded))

	p, dir := newTakeoverPipeline130r1(t)
	p.InstallAsDefault()
	if err := p.FlushNow(); err != nil {
		t.Fatalf("FlushNow: %v", err)
	}

	var raw strings.Builder
	for _, name := range logFileNames(t, dir) {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		raw.Write(body)
	}
	if !strings.Contains(raw.String(), takeoverBootMsg130r1) {
		t.Fatalf("the taken-over record is not in the file at all, so this case would be green on nothing: %s",
			raw.String())
	}
	if strings.Contains(raw.String(), seeded) {
		t.Fatalf("the taken-over record's secret is in the file verbatim: the handover bypassed the "+
			"redacting handler.\n%s", raw.String())
	}
	if want := secretRedactForTest(seeded); !strings.Contains(raw.String(), want) {
		t.Errorf("the masked form %q is not in the file, so the record was rewritten rather than redacted: %s",
			want, raw.String())
	}
}

func countMsg130r1(recs []takeoverRecord130r1, msg string) int {
	n := 0
	for _, r := range recs {
		if r.Msg == msg {
			n++
		}
	}
	return n
}
