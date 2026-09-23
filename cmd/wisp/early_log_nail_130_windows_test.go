//go:build windows

package main

// Ticket 130 AC#3: a record that speaks before the listener exists must still
// reach the file, and getting it there may not be bought by silencing the
// console.
//
// What is broken (two independent reproductions, R-117-2 and R-125-3): the
// sealing-seam verdict of internal/winsec is logged from
// internal/risk/winsec_c26.go's package init(), which runs before main() and
// therefore before any listener exists at all. Go's stock default logger put it
// on stderr, the rolling JSONL pipeline did not exist yet, and the one file the
// process later wrote held exactly the records made after the install. A reader
// holding one JSONL off a disk could not tell whether the resolver was accepted
// or refused - the single most auditable fact about the sealing seam.
//
// The ruling this case pins (a-with-mirror, docs/reports/pending-and-issues.md
// A110 row 3): early records are kept a SECOND time in an in-memory buffer that
// internal/observe installs from its own package init() - earlier than any
// package that logs from init(), because observe is a dependency of risk - and
// cmd/wisp/logsink.go replays that buffer into the pipeline at install time,
// ahead of the install record. The stderr pass is never switched off, so the
// worst case of this design is today's behaviour, not something below it.
//
// Three assertions, and why each one is here rather than a looser form:
//
//  1. The record EXISTS in the JSONL the shipped process wrote. This is the
//     assertion that is red before the fix: on this leg the file held one
//     record (the install booking) and the resolver verdict held none. It is
//     not enough for the sentence to be somewhere else: it has to be in the
//     file that outlives the process.
//  2. Its INDEX and its TIME both precede the install record's. Existence alone
//     would be satisfied by re-printing the same sentence after the listener is
//     up, which is a different fact - "the listener decided to say this later"
//     instead of "this happened before the listener existed". Reading the
//     timestamp off the persisted record (not off the file's mtime) is what
//     makes that distinction checkable after the fact.
//  3. stderr still carries it EXACTLY ONCE, and so does the install record.
//     "Exactly" is the whole point of this row: zero means the mirror was
//     dropped to make assertion 1 cheap, two means the replay double-printed
//     onto the console. Ticket 117's fan-out promise (the run leg's console saw
//     these lines before any file existed) and this ticket's buffer are kept by
//     the same word.
//
// Deliberately not claimed here: this covers the legs that install a listener.
// A leg that never installs one (wisp providers, wisp models list/verify,
// wisp doctor) is not rescued by this fix - it is ticket 131's enumeration
// domain, and this file stays silent about it rather than asserting a promise
// nothing made.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// The two sentences this case is about, copied as literals on purpose (the same
// rule as residentInstallMsg and sealNoticeMsg in the sibling files): if either
// is reworded, a case in this file goes red and somebody has to say out loud
// where the record that proves the early verdict now lives.
const (
	// early130ResolverMsg is internal/winsec/resolve.go's success branch, the
	// record emitted from internal/risk's package init(). The refusal branch is
	// the same gap with a different level; it cannot be produced on a host whose
	// conformance probes pass, so this literal names the reachable one.
	early130ResolverMsg = "winsec: sealing path resolver installed"
	// early130InstallMsg is installLogSink's own booking record; it is the marker
	// the early record has to precede.
	early130InstallMsg = "wisp: persistent log sink installed"
)

// earlyRecord130 is one JSONL line, decoded only for the fields the three
// assertions read. ProbesPassed is a pointer so "the attribute is absent" and
// "the attribute says zero" are two different readings of the file rather than
// one value a template defaulted.
type earlyRecord130 struct {
	Time         string `json:"time"`
	Level        string `json:"level"`
	Msg          string `json:"msg"`
	Resolver     string `json:"resolver"`
	ProbesPassed *int   `json:"probes_passed"`
	Dir          string `json:"dir"`
}

func (r earlyRecord130) String() string {
	return fmt.Sprintf("{time=%s level=%s msg=%q resolver=%q probes_passed=%v}",
		r.Time, r.Level, r.Msg, r.Resolver, r.ProbesPassed)
}

func msgsOf130(recs []earlyRecord130) []string {
	out := make([]string, 0, len(recs))
	for _, r := range recs {
		out = append(out, r.Msg)
	}
	return out
}

// runSecretListLeg drives the shipped binary on the leg that installs a
// listener (wisp secret list, cmd/wisp/secret.go) against a data root handed in
// from outside the repository, and returns what that process left on disk plus
// what it printed.
//
// The subprocess is the point: the record under test is produced by a package
// init(), which no in-process test can re-run, and the version boundary that
// makes the fix real is the init-order between internal/observe and
// internal/risk - visible only in a binary that actually boots from scratch.
// WISP_ENV=test plus WISP_TEST_DATA_DIR keep every byte inside the harness' own
// directory (internal/proc/envfork.go), so the owner's real data root is not
// touched by this case.
func runSecretListLeg(t *testing.T, exe, dataDir string) (recs []earlyRecord130, stderr string, stdout string) {
	t.Helper()
	cmd := exec.Command(exe, "secret", "list")
	cmd.Dir = filepath.Dir(exe)
	cmd.Env = append(os.Environ(), "WISP_ENV=test", procTestDataDirEnv+"="+dataDir)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr, cmd.Stdin = &out, &errb, nil
	if err := cmd.Run(); err != nil {
		t.Fatalf("run %s secret list against %s: %v\nstdout:\n%s\nstderr:\n%s",
			exe, dataDir, err, out.String(), errb.String())
	}
	return readEarlySink130(t, logSinkDir(dataDir)), errb.String(), out.String()
}

// readEarlySink130 returns the shipped process's records in file order,
// refusing a directory that holds no pipeline file. The naming rule comes from
// observe itself (CountLogFiles), so a renamed pipeline is a red case here
// rather than a silently empty read - the same discipline readResidentSink uses.
func readEarlySink130(t *testing.T, dir string) []earlyRecord130 {
	t.Helper()
	n, err := observe.CountLogFiles(dir)
	if err != nil {
		t.Fatalf("reading the sink dir %s: %v", dir, err)
	}
	if n == 0 {
		t.Fatalf("no wisp-<day>-<seq>.jsonl in %s: the shipped process wrote nothing there", dir)
	}
	if n != 1 {
		t.Errorf("pipeline files in %s = %d, want 1 (one boot, one file)", dir, n)
	}
	var out []earlyRecord130
	for _, name := range mustGlob117(t, filepath.Join(dir, "wisp-*.jsonl")) {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var rec earlyRecord130
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("line %q in %s is not a JSON object: %v", line, name, err)
			}
			out = append(out, rec)
		}
	}
	t.Logf("early-record sink %s: %d file(s), %d record(s), msgs=%v", dir, n, len(out), msgsOf130(out))
	return out
}

// TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord is the headline of this
// file and the reading that is red before the fix.
//
// Mutation states this case was written to separate:
//   - delete the replay call in cmd/wisp/logsink.go (M-1) and assertions 1 and 2
//     go red while 3 stays green: the early record is back to being console-only.
//   - install the buffer from main() instead of a package init() (M-2) and the
//     same two go red for the same reason, which is the proof that the buffer has
//     to live in internal/observe: main-package init() runs last, so nothing in
//     cmd/wisp can catch a record made during another package's init().
//   - drop the stderr side of the early handler (M-3) and assertion 3 goes red
//     with a count of zero, which is the proof that the file copy was not paid
//     for by taking the console away.
//   - replay through the process default instead of into the pipeline and
//     assertion 3 goes red with a count of two.
func TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()

	recs, stderr, stdout := runSecretListLeg(t, exe, dataDir)
	console := func() string { return "stdout:\n" + stdout + "\nstderr:\n" + stderr }

	if len(recs) == 0 {
		t.Fatalf("%s holds a pipeline file with no record in it, so there is nothing to read the ordering out of.\n%s",
			logSinkDir(dataDir), console())
	}

	earlyIdx, installIdx := -1, -1
	for i, r := range recs {
		switch r.Msg {
		case early130ResolverMsg:
			if earlyIdx < 0 {
				earlyIdx = i
			}
		case early130InstallMsg:
			if installIdx < 0 {
				installIdx = i
			}
		}
	}

	// Assertion 1: the early verdict is in the file that outlives the process.
	if earlyIdx < 0 {
		t.Fatalf("assertion 1: %q is in no record of the file the shipped process wrote, "+
			"which is ticket 130's gap exactly: the record was made before the listener existed, "+
			"so it went to the console and nowhere else.\nfile records: %v\n%s",
			early130ResolverMsg, msgsOf130(recs), console())
	}
	early := recs[earlyIdx]
	if early.Resolver == "" {
		t.Errorf("assertion 1: the early record carries no resolver attribute: %v", early)
	}
	if early.ProbesPassed == nil {
		t.Errorf("assertion 1: the early record carries no probes_passed attribute, so it is not the "+
			"conformance verdict this ticket is about: %v", early)
	} else if *early.ProbesPassed < 1 {
		t.Errorf("assertion 1: probes_passed = %d, want at least 1 (the value internal/winsec's "+
			"resolverProbeShapes returns on this platform)", *early.ProbesPassed)
	}
	if early.Level != "INFO" {
		t.Errorf("assertion 1: early record level = %q, want INFO: %v", early.Level, early)
	}
	if installIdx < 0 {
		t.Fatalf("no %q record in the file, so the ordering below has no landmark: %v",
			early130InstallMsg, msgsOf130(recs))
	}

	// Assertion 2: index first, and the persisted time agrees.
	if earlyIdx >= installIdx {
		t.Errorf("assertion 2: the early record is index %d and the install record is index %d, "+
			"want the early one first; a replay that lands after the install is a different fact "+
			"(records: %v)", earlyIdx, installIdx, recs)
	}
	earlyTime, err := time.Parse(time.RFC3339Nano, early.Time)
	if err != nil {
		t.Errorf("assertion 2: early record time %q is not the RFC3339Nano stamp the pipeline writes: %v", early.Time, err)
	}
	installTime, err := time.Parse(time.RFC3339Nano, recs[installIdx].Time)
	if err != nil {
		t.Errorf("assertion 2: install record time %q is not the RFC3339Nano stamp the pipeline writes: %v", recs[installIdx].Time, err)
	}
	if err == nil && !earlyTime.Before(installTime) {
		t.Errorf("assertion 2: early record stamp %s is not before the install stamp %s", earlyTime, installTime)
	}

	// Assertion 3: the console copy is still there, exactly once, and the
	// replay did not double it.
	if n := strings.Count(stderr, early130ResolverMsg); n != 1 {
		if n == 0 {
			t.Errorf("assertion 3: the early record never reached the child's stderr (count 0): the mirror was "+
				"switched off to make the file copy cheap, which is the failure mode (a-with-mirror) exists to rule out.\n%s", console())
		} else {
			t.Errorf("assertion 3: the early record reached the child's stderr %d times, want exactly 1: the replay is "+
				"double-printing onto the console.\n%s", n, console())
		}
	}
	if n := strings.Count(stderr, early130InstallMsg); n != 1 {
		t.Errorf("assertion 3: the install record reached the child's stderr %d times, want exactly 1 (ticket 117's fan-out).\n%s",
			n, console())
	}

	t.Logf("EARLY RECORD ON DISK: index %d of %d, stamp %s, resolver=%q probes_passed=%v; install at index %d stamp %s",
		earlyIdx, len(recs), early.Time, early.Resolver, early.ProbesPassed, installIdx, recs[installIdx].Time)
}

// TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot is the placement half: the
// record that ticket 130 rescues must not become a byte that lands outside the
// env's own tree. Same rule the run and resident legs are already held to, read
// off this leg instead of inferred from it.
func TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()

	recs, stderr, stdout := runSecretListLeg(t, exe, dataDir)
	if len(recs) == 0 {
		t.Fatalf("the shipped process wrote no record under %s.\nstdout:\n%s\nstderr:\n%s",
			logSinkDir(dataDir), stdout, stderr)
	}
	for _, p := range jsonlFilesUnder(t, dataDir) {
		if rel, err := filepath.Rel(logSinkDir(dataDir), p); err != nil || strings.HasPrefix(rel, "..") {
			t.Errorf("the replayed early record was written outside %s: %q", logSinkDir(dataDir), p)
		}
	}
}
