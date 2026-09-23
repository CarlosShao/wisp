package main

// Ticket 131 AC#2 and AC#3: the two legs of this command line that had no nail
// of their own get one, in the shape the ticket demanded - in-process, through
// the same entry point main.go dispatches to, with the verdict read off the disk
// the leg's own listener writes.
//
// WHY THESE CASES EXIST, stated once because it is the same sentence twice.
// `acceptor-ticket127` mutation 4 deleted cmd/wisp/models.go's nine install lines
// (the block ticket 121's 7fe5e73 put there) and the package did not notice:
// `go build` rc=0, `go vet` rc=0, `go test -count=1 -v ./cmd/wisp/` ok with 76
// cases and zero red. The same deletion on the resident leg (ticket 127 AC#1)
// took three cases red, and on the run leg three more - which is what proves the
// missing thing was never the code, only a nail. `wisp secret` is the same
// reading measured the same way by `acceptor-ticket117` as R-117-1, with the
// extra twist that the leg had not even been ruled (AC#3): the install now lives
// in cmdSecret, so this file is what keeps it there.
//
// WHAT IS DRIVEN, which is the whole content of AC#2's "可像 run 腿那样进程内
// 驱动" note:
//
//	models  cmdModels([]string{"ensure", id}, modelsIO{dataDir: <root>})
//	        - main.go's "models" case calls exactly this function with the same
//	        argv shape, and modelsIO.dataDir is the injection seam ticket 121
//	        left in place for this.
//	secret  cmdSecret([]string{"set", name, "--from-stdin"})
//	        - main.go's "secret" case calls exactly this function. cmdSecret
//	        resolves the data root from the OS rather than from an injected
//	        field, so the two variables below stand for "the test env's own data
//	        root": WISP_ENV picks the env, WISP_TEST_DATA_DIR is an identity
//	        contract (internal/proc/envfork.go) that comes back verbatim, which
//	        is what lets a case compare strings with what the process wrote.
//
// NO NETWORK, written because both drives sit next to code that downloads: the
// models drive names a model id the committed signed manifest does not carry, so
// Manager.Ensure refuses at Manifest.FindModel - one statement before the
// local_override branch, two before any transport exists. A wrong-digest
// local_override would be the deeper refusal and would cost a real 300 MB
// archive to produce one red hash, so this case drives the shallower branch,
// which still walks LoadSignedManifest's C29 signature check, config.LoadFile,
// models.NewManager, the D43 FirstRun/Downloading/Error walk, and the record
// under test. AC#2's criterion is "盘上 jsonl 有 models: 记录", and a refused
// hand-off is exactly the record this leg exists to keep - models.go's own
// comment: "a refused hand-off nobody records is a security verdict that leaves
// no trace".
//
// WHAT IS NOT CLAIMED here: that the success branch of the hand-off books a
// record too. It is the same logf closure over the same sink, and the byte that
// decides which of the two records is written is a model nobody ships into a
// test tree; ticket 109's success-path hand-off tests live in internal/models.
//
// NO t.Skip, no relaxed assertion, no wall-clock budget and no polling: every
// record is read after the driven function returns, and that function closes its
// own sink on the way out (models.go's `defer sink.close()`, secret.go's defer
// beside its install), so the flush is ordered by the code under test instead of
// waited for. The three drives in this file together cost milliseconds, which is
// also why none of them is a subprocess case like ticket 127's.
//
// PLATFORM, because these two files are untagged and reach across to
// windows-tagged ones: sinkInstallRecord and residentInstallMsg come from
// resident_sink_nail_127_windows_test.go, and readSink, readResidentSink and
// jsonlFilesUnder are the instruments declared there and in
// logsink_windows_test.go. That is the choice, not an accident. package main has
// no Linux leg (main.go's sherpa import is what makes `GOOS=linux go vet
// ./cmd/wisp/` rc=1, and scripts/wisp-cli-tests.sh says the same in its header),
// so today both directions compile or nothing does; on the day a Linux leg lands,
// an untagged file that fails to *compile* is a louder reading than a second
// windows-tagged file that silently drops the enumeration gate's denominator on
// the other leg - which is exactly the hole ticket 121 AC#3 wrote out at length.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// The sentences the two install blocks are responsible for, copied as literals
// on purpose (the rule sealNoticeMsg and residentInstallMsg follow): if any of
// them is reworded, a case in this file goes red and somebody has to say out
// loud where the record that proves the listener was there now lives.
const (
	// modelsRefusalNotice is modelsEnsure's loud fallback when the sink cannot be
	// installed (the Fprintf inside the nine lines). It is the half of the block
	// that never reaches the log file, so the console is the only place to catch
	// it.
	modelsRefusalNotice = "wisp models: 持久日志未启用"
	// secretRefusalNotice is the same sentence on the credential leg.
	secretRefusalNotice = "wisp secret: 持久日志未启用"
	// degradationPromise is the half of the run and resident refusals that says
	// what is lost, and the half secret.go's sentence (added by this ticket)
	// copies. Ticket 121 wrote models.go's refusal as "只会到 stderr" without it,
	// and this ticket does not rewrite another ticket's user-visible sentence to
	// make an assertion fit, so the two legs are pinned against their own wording
	// and the common half is pinned for both.
	degradationPromise = "不会落盘"
	// degradationConsole is the half every one of the four refusals promises: the
	// notices still reach a terminal.
	degradationConsole = "只会到 stderr"
	// secretStoredMsg is runSet's audit record: that a blob now exists, named by
	// ref, env and portability - never by value (C28, and the ruling comment
	// beside secretCmd.runSet).
	secretStoredMsg = "wisp secret: stored dpapi blob"
	// secretUnsetMsg is runUnset's audit record of a delete that left config
	// references dangling; --force books it at WARN.
	secretUnsetMsg = "wisp secret unset: deleted"
	// modelsRecordPrefix is the tag handOffModel puts on every record this leg
	// books, and the string AC#2's criterion is about.
	modelsRecordPrefix = "models: hand-off "
)

// legSink131 is one driven leg's residue: the records its listener wrote, and
// the two console streams it printed to. Keeping the streams apart is what makes
// the fan-out assertion mean something - the leg's own stdout carries the
// operator-facing sentence, so a merged buffer could satisfy "the record reached
// the console" without the console handler ever having run.
type legSink131 struct {
	dir     string
	records []sinkInstallRecord
	stdout  string
	stderr  string
	console string
}

func (s legSink131) String() string {
	return fmt.Sprintf("{dir=%s records=%d}", s.dir, len(s.records))
}

// captureLeg131 runs act with the process-wide state the production wrappers
// read set aside for the drive and restored before it returns: slog's default
// logger, and the three standard handles.
//
// The restore is not cosmetics. installLogSink replaces slog's default with a
// handler over a pipeline that the drive closes on its way out, so a later case
// in this package that logs while that default is still installed would write
// into a closed pipeline - a green nobody should have to reason about.
func captureLeg131(t *testing.T, stdin string, act func() int) (code int, stdout, stderr string) {
	t.Helper()
	oldDefault, oldOut, oldErr, oldIn := slog.Default(), os.Stdout, os.Stderr, os.Stdin
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = outW, errW
	if stdin != "" {
		f, err := os.CreateTemp(t.TempDir(), "wisp131-stdin")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString(stdin); err != nil {
			t.Fatal(err)
		}
		if _, err := f.Seek(0, 0); err != nil {
			t.Fatal(err)
		}
		os.Stdin = f
		defer f.Close()
	}

	code = act()

	if err := outW.Close(); err != nil {
		t.Fatal(err)
	}
	if err := errW.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr, os.Stdin = oldOut, oldErr, oldIn
	slog.SetDefault(oldDefault)

	var ob, eb bytes.Buffer
	if _, err := ob.ReadFrom(outR); err != nil {
		t.Fatal(err)
	}
	if _, err := eb.ReadFrom(errR); err != nil {
		t.Fatal(err)
	}
	return code, ob.String(), eb.String()
}

// readLegSink131 is the instrument every nail in this file shares, and the one
// the AC#4 gate names in a registered nail. It refuses an empty tree - "the
// directory exists and says nothing" is not a pass - and it refuses an empty
// record set, which is R-127-4's shape (a panic eats the red name and truncates
// the package, so an empty read has to be a named failure here rather than an
// index out of range in the caller).
func readLegSink131(t *testing.T, dataDir, stdout, stderr string) legSink131 {
	t.Helper()
	s := legSink131{
		dir:     logSinkDir(dataDir),
		stdout:  stdout,
		stderr:  stderr,
		console: "stdout:\n" + stdout + "\nstderr:\n" + stderr,
	}
	n, err := observe.CountLogFiles(s.dir)
	if err != nil || n == 0 {
		// Two shapes, one reading: the directory never appeared (the install
		// line is gone, so nothing ever asked for it) and the directory is there
		// but holds no pipeline file. Both mean this leg ran to completion with no
		// listener, and the message has to say so with the leg's own console in
		// hand - "counting pipeline files: The system cannot find the file
		// specified" names a syscall, not the claim that broke.
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("counting pipeline files in %s: %v", s.dir, err)
		}
		if err == nil && n != 0 {
			t.Fatalf("unreachable: CountLogFiles reported %d files with no error", n)
		}
		t.Fatalf("no wisp-<day>-<seq>.jsonl in %s: the leg ran to completion and its listener wrote nothing there.\n%s",
			s.dir, s.console)
	}
	if n != 1 {
		t.Errorf("pipeline files in %s = %d, want 1 (one drive of one leg, one file)", s.dir, n)
	}
	names, err := filepath.Glob(filepath.Join(s.dir, "wisp-*.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var rec sinkInstallRecord
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("line %q in %s is not a JSON object: %v", line, name, err)
			}
			s.records = append(s.records, rec)
		}
	}
	if len(s.records) == 0 {
		t.Fatalf("%s holds a pipeline file with no record in it", s.dir)
	}
	t.Logf("leg sink %s: %d record(s) %v", s.dir, len(s.records), msgsOf131(s.records))
	return s
}

func msgsOf131(recs []sinkInstallRecord) []string {
	out := make([]string, 0, len(recs))
	for _, r := range recs {
		out = append(out, fmt.Sprintf("%s/%s", r.Level, r.Msg))
	}
	return out
}

// assertInstallRecordFirst131 is AC#2's second clause, shared by every nail
// here: the listener's own booking record must precede the event this file is
// about, must name THIS data root, and nothing but a record ticket 130 replays
// from before the listener existed may sit in front of it. Ordering is not
// decoration - installLogSink states the invariant ("a sink installed after the
// first event loses exactly the record it exists to keep"), and a leg that
// installed its sink after the event under test would satisfy a "the record
// exists" assertion while breaking the one about when.
//
// WHY THE FIRST CLAUSE CHANGED SHAPE (ticket 130; the long version of this
// reasoning is in cmd/wisp/resident_sink_nail_127_windows_test.go's header, for
// the nail that ticket 130's ruling names explicitly). This used to require the
// booking to be record 0. Ticket 130's fix moves a package init()'s sealing
// verdict into the same file AHEAD of the booking, so "record 0 is the booking"
// stopped being a claim about ordering and became a claim that the replay is
// missing - after the fix it passes exactly when the defect is present. The
// comparison therefore moved onto what the old one was protecting: the booking is
// the landmark, everything in front of it is a replayed boot record, and because
// every caller here asserts an event record that is neither the booking nor a
// boot record, that partition is the same "the listener was installed before the
// event" claim, stated without the incidental index.
//
// What this helper deliberately does NOT claim: that a boot record exists. These
// nails run in-process, and one test binary installs many sinks, so which sink
// receives the replay is decided by test order. The replay itself is pinned where
// the process is real - cmd/wisp/early_log_nail_130_windows_test.go and the 127
// file's subprocess cases.
func assertInstallRecordFirst131(t *testing.T, s legSink131, dataDir string) {
	t.Helper()
	if len(s.records) == 0 {
		t.Fatal("no records to read: the instrument returned an empty set")
	}
	installIdx := -1
	for i, r := range s.records {
		if r.Msg == residentInstallMsg {
			installIdx = i
			break
		}
	}
	if installIdx < 0 {
		t.Fatalf("no %q record in the file this leg wrote, so the leg booked no listener at all (records: %v)",
			residentInstallMsg, msgsOf131(s.records))
	}
	for i, r := range s.records[:installIdx] {
		if r.Msg != residentEarlyResolverMsg {
			t.Errorf("record %d sits before the listener's booking and is %q, want %q or nothing: a leg that logs "+
				"something before installLogSink loses exactly the record it exists to keep (records: %v)",
				i, r.Msg, residentEarlyResolverMsg, msgsOf131(s.records))
		}
	}
	first := s.records[installIdx]
	if want := logSinkDir(dataDir); first.Dir != want {
		t.Errorf("install record names dir = %q, want %q: AC#2's judgment is exactly this comparison, and it is the data root of the env the leg was driven with",
			first.Dir, want)
	}
	if first.MinLevel != logSinkLevel {
		t.Errorf("install record min_level = %q, want the schema default %q", first.MinLevel, logSinkLevel)
	}
	// Placement, the same rule ticket 117 AC#3 states for the other legs: no
	// byte of this listener's tree lives outside <data root>\logs.
	for _, p := range jsonlFilesUnder(t, dataDir) {
		if rel, err := filepath.Rel(s.dir, p); err != nil || strings.HasPrefix(rel, "..") {
			t.Errorf("the listener wrote %q outside %s", p, s.dir)
		}
	}
}

// TestAC2ModelsLegBooksItsHandOffVerdictOnDisk is AC#2's nail.
func TestAC2ModelsLegBooksItsHandOffVerdictOnDisk(t *testing.T) {
	dataDir := t.TempDir()
	writeModelStoreConfig131(t, dataDir)
	t.Setenv("WISP_MODELS_MANIFEST", repoSignedManifest131(t))

	code, stdout, stderr := captureLeg131(t, "", func() int {
		return cmdModels([]string{"ensure", modelsAbsentID131}, modelsIO{
			stdout:  os.Stdout,
			stderr:  os.Stderr,
			dataDir: dataDir,
		})
	})
	s := readLegSink131(t, dataDir, stdout, stderr)
	// A refused hand-off is exit 1 and must not be softened into a warning
	// (cmdModels' documented contract). Checked first, so a change of verdict is
	// reported as the verdict change it is rather than as a missing record.
	if code != 1 {
		t.Errorf("`wisp models ensure %s` exited %d, want 1 (refused hand-off)\n%s",
			modelsAbsentID131, code, s.console)
	}
	assertInstallRecordFirst131(t, s, dataDir)

	var hits []sinkInstallRecord
	for _, r := range s.records {
		if strings.HasPrefix(r.Msg, modelsRecordPrefix+"REFUSED") {
			hits = append(hits, r)
		}
	}
	if len(hits) != 1 {
		t.Fatalf("records naming a refused hand-off = %d, want exactly 1; all records: %v",
			len(hits), msgsOf131(s.records))
	}
	hit := hits[0]
	if !strings.Contains(hit.Msg, "id="+modelsAbsentID131) {
		t.Errorf("the refusal record does not name the model it refused: %q", hit.Msg)
	}
	// The end state is part of the record on purpose - "Error" and "still
	// Downloading" are two different failures (models.go's own comment) - so the
	// record is only worth having if it still says which one this was.
	if !strings.Contains(hit.Msg, "state=Error") {
		t.Errorf("the refusal record does not name the D43 end state the walk reached: %q", hit.Msg)
	}
	if hit.Level != "INFO" {
		t.Errorf("refusal record level = %q, want INFO (sinkLogger.Info is what the nine lines install)", hit.Level)
	}
	// The same sentence on both outlets: stderr is where the operator standing at
	// a terminal reads it, and the file is where it survives the terminal.
	if !strings.Contains(stderr, "wisp models ensure: 交还被拒绝") {
		t.Errorf("the refusal never reached the leg's stderr; the fan-out went one way only.\n%s", s.console)
	}
	t.Logf("MODELS LEG RECORDED ON DISK: %q (install dir=%s, exit=%d)", hit.Msg, s.dir, code)
}

// TestAC3SecretLegBooksItsAuditRecordsOnDisk is AC#3's nail, and it is what the
// AC#3 ruling ("install") costs: the credential leg's two audit records - an
// INFO for storing a blob, a WARN for a --force delete that left references
// dangling - now have to reach the same file as every other leg's.
//
// The planted value is a fake key, and the case asserts its ABSENCE from every
// byte the new outlet wrote: C28 re-asserted against the new file rather than
// around it.
func TestAC3SecretLegBooksItsAuditRecordsOnDisk(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("WISP_ENV", "test")
	t.Setenv("WISP_TEST_DATA_DIR", dataDir)
	const (
		name    = "nail131"
		ref     = "dpapi:nail131"
		fakeKey = "sk-FAKE-131-nail-plaintext-0123456789"
	)
	// config.toml naming the ref is what makes the later --force delete the
	// dangling-reference shape, and it is the line runSet's footer tells the
	// operator to write.
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(
		"schema_version = 2\n\n[llm.providers.acme]\napi_key_ref = \""+ref+"\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	setCode, setOut, setErr := captureLeg131(t, fakeKey, func() int {
		return cmdSecret([]string{"set", name, "--from-stdin"})
	})
	s := readLegSink131(t, dataDir, setOut, setErr)
	if setCode != 0 {
		t.Fatalf("`wisp secret set %s --from-stdin` exited %d, want 0\n%s", name, setCode, s.console)
	}
	assertInstallRecordFirst131(t, s, dataDir)

	var stored []sinkInstallRecord
	for _, r := range s.records {
		if r.Msg == secretStoredMsg {
			stored = append(stored, r)
		}
	}
	if len(stored) != 1 {
		t.Fatalf("records naming %q = %d, want 1; all records: %v", secretStoredMsg, len(stored), msgsOf131(s.records))
	}
	if stored[0].Level != "INFO" {
		t.Errorf("store record level = %q, want INFO", stored[0].Level)
	}
	if !strings.Contains(setErr, secretStoredMsg) {
		t.Errorf("the audit record never reached the console side of the fan-out.\n%s", s.console)
	}
	assertNoPlaintextInSink131(t, s, fakeKey)
	if _, err := os.Stat(filepath.Join(dataDir, "secrets", name)); err != nil {
		t.Errorf("the blob the record announces is not on disk: %v", err)
	}

	// The destructive half, over the same listener: --force past a referenced
	// blob, which runUnset books at WARN because config.toml now points at
	// nothing.
	forceCode, forceOut, forceErr := captureLeg131(t, "", func() int {
		return cmdSecret([]string{"unset", name, "--force"})
	})
	s2 := readLegSink131(t, dataDir, forceOut, forceErr)
	if forceCode != 0 {
		t.Fatalf("`wisp secret unset %s --force` exited %d, want 0\n%s", name, forceCode, s2.console)
	}
	assertInstallRecordFirst131(t, s2, dataDir)
	var warned []sinkInstallRecord
	for _, r := range s2.records {
		if strings.HasPrefix(r.Msg, secretUnsetMsg) {
			warned = append(warned, r)
		}
	}
	if len(warned) != 1 {
		t.Fatalf("records naming %q = %d, want 1; all records: %v", secretUnsetMsg, len(warned), msgsOf131(s2.records))
	}
	if warned[0].Level != "WARN" {
		t.Errorf("dangling-reference delete record level = %q, want WARN (runUnset logs --force at Warn on purpose)", warned[0].Level)
	}
	if !strings.Contains(warned[0].Msg, ref) {
		t.Errorf("the delete record does not name the ref it deleted: %q", warned[0].Msg)
	}
	t.Logf("SECRET LEG RECORDED ON DISK: %q and %q (install dir=%s)", stored[0].Msg, warned[0].Msg, s2.dir)
}

// TestAC2AC3DegradedLegsStillDeliverTheirVerdict is the fallback branch of both
// install blocks, driven the way ticket 127 drove the resident leg: stand a
// plain file where <data root>\logs must go, so the sink's first MkdirAll fails,
// and require each leg to (a) say so on the console, naming what is lost,
// (b) still deliver the verdict the leg exists to deliver, and (c) write nothing
// anywhere. Without (b) a full disk would lock the owner out of `wisp secret`,
// and a broken log tree would turn a refused model hand-off into an exit code
// that means something else.
func TestAC2AC3DegradedLegsStillDeliverTheirVerdict(t *testing.T) {
	t.Run("models", func(t *testing.T) {
		dataDir := t.TempDir()
		writeModelStoreConfig131(t, dataDir)
		t.Setenv("WISP_MODELS_MANIFEST", repoSignedManifest131(t))
		blocker := logSinkDir(dataDir)
		if err := os.WriteFile(blocker, []byte("not a directory"), 0o600); err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		code := cmdModels([]string{"ensure", modelsAbsentID131}, modelsIO{
			stdout:  &stdout,
			stderr:  &stderr,
			dataDir: dataDir,
		})
		if !strings.Contains(stderr.String(), modelsRefusalNotice) {
			t.Errorf("a log directory that cannot be opened produced no named refusal (%q on stderr).\n"+
				"stderr:\n%s\nstdout:\n%s\nThat sentence is the half of the nine lines which is not a comment.",
				modelsRefusalNotice, stderr.String(), stdout.String())
		}
		if !strings.Contains(stderr.String(), degradationConsole) {
			t.Errorf("the models refusal does not say where the notices still go (%q): %s", degradationConsole, stderr.String())
		}
		if code != 1 {
			t.Errorf("a broken log tree changed the verdict to %d, want 1 (the refusal must survive its own listener failing)", code)
		}
		if files := jsonlFilesUnder(t, dataDir); len(files) != 0 {
			t.Errorf("the refused install produced %d jsonl file(s) anyway: %v", len(files), files)
		}
		if st, err := os.Stat(blocker); err != nil || st.IsDir() {
			t.Errorf("the blocker at %s is gone or became a directory (stat err %v)", blocker, err)
		}
	})

	t.Run("secret", func(t *testing.T) {
		dataDir := t.TempDir()
		t.Setenv("WISP_ENV", "test")
		t.Setenv("WISP_TEST_DATA_DIR", dataDir)
		blocker := logSinkDir(dataDir)
		if err := os.WriteFile(blocker, []byte("not a directory"), 0o600); err != nil {
			t.Fatal(err)
		}
		const fakeKey = "sk-FAKE-131-degraded-plaintext-0123456789"
		code, stdout, stderr := captureLeg131(t, fakeKey, func() int {
			return cmdSecret([]string{"set", "nail131-degraded", "--from-stdin"})
		})
		if !strings.Contains(stderr, secretRefusalNotice) {
			t.Errorf("the credential leg booted over an unopenable log directory and said nothing about it.\nstdout:\n%s\nstderr:\n%s",
				stdout, stderr)
		}
		if !strings.Contains(stderr, degradationPromise) {
			t.Errorf("the credential leg's refusal does not say what is lost: %s", stderr)
		}
		if code != 0 {
			t.Errorf("`wisp secret set` exited %d over a broken log tree, want 0 (a full disk must not lock the owner out of their own keys)\nstdout:\n%s\nstderr:\n%s",
				code, stdout, stderr)
		}
		if files := jsonlFilesUnder(t, dataDir); len(files) != 0 {
			t.Errorf("the refused install produced %d jsonl file(s) anyway: %v", len(files), files)
		}
		if strings.Contains(stdout, fakeKey) || strings.Contains(stderr, fakeKey) {
			t.Error("the degraded branch printed the plaintext it was storing")
		}
	})
}

// ---------------------------------------------------------------------------
// Fixtures, local to this file on purpose: ticket 131 owns cmd/wisp and nothing
// under internal/, so the signed manifest is read where it is committed rather
// than re-signed into a temp dir - which would need a private key this repository
// does not carry (buildinfo.MinisignPublicKey is only the public half).
// ---------------------------------------------------------------------------

// modelsAbsentID131 is an id the committed signed manifest does not carry. That
// is the whole fixture: Manager.Ensure refuses at FindModel, before the
// local_override branch and before any transport exists.
const modelsAbsentID131 = "absent-in-signed-manifest-131"

// repoSignedManifest131 returns the committed manifest's path in the form
// WISP_MODELS_MANIFEST wants it (models.ResolveManifestPath appends .minisig).
// It is the same file internal/models/manifest_real_test.go verifies against
// buildinfo's key, so the leg under test reads a manifest whose signature really
// was checked by C29 rather than a fixture this test signed itself.
func repoSignedManifest131(t *testing.T) string {
	t.Helper()
	for _, dir := range []string{"../..", "."} {
		p := filepath.Join(dir, "models", "manifest.json")
		if _, err := os.Stat(p); err == nil {
			abs, err := filepath.Abs(p)
			if err != nil {
				t.Fatal(err)
			}
			return abs
		}
	}
	t.Fatal("committed models/manifest.json not found from the test working directory")
	return ""
}

// writeModelStoreConfig131 puts a config.toml in the data root naming the store a
// subdirectory of it. modelsEnsure reaches the nine lines only after
// config.LoadFile succeeds, so without this file the drive would stop at
// Unconfigured and never touch the install.
func writeModelStoreConfig131(t *testing.T, dataDir string) {
	t.Helper()
	body := "schema_version = 2\n\n[models]\ndir = \"store\"\n"
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// assertNoPlaintextInSink131 points C28 at the new outlet: every byte the driven
// leg wrote into its log tree, and everything it printed, must be free of the
// value. scanForPlaintext is secret_test.go's instrument, so this assertion and
// that file's positive control share exactly one pair of eyes.
func assertNoPlaintextInSink131(t *testing.T, s legSink131, wantAbsent string) {
	t.Helper()
	if hits := scanForPlaintext(t, []string{s.dir}, wantAbsent, 1<<20); len(hits) != 0 {
		t.Errorf("the credential leg's new persistent outlet holds the plaintext it is forbidden to hold: %v", hits)
	}
	if strings.Contains(s.console, wantAbsent) {
		t.Error("the driven leg printed the plaintext to the console")
	}
}

// The four nails below are the registrations AC#4's gate enumerates: a leg key
// taken from main.go's own dispatch labels, and a function value rather than a
// string, so a rename or a build tag that hides a nail is a compile error
// instead of a quiet green. registerLegNail131 (leg_sink_gate_131_test.go) reads
// the case's name back out of the compiled function, so the key is a claim about
// main.go and the gate is what checks it.
var (
	nailRun131      = registerLegNail131("run", TestAC2SealNoticeLandsInTheRunLegLogFile)
	nailResident131 = registerLegNail131("no-args", TestAC1ResidentLegInstallsItsLogListenerOnDisk)
	nailModels131   = registerLegNail131("models", TestAC2ModelsLegBooksItsHandOffVerdictOnDisk)
	nailSecret131   = registerLegNail131("secret", TestAC3SecretLegBooksItsAuditRecordsOnDisk)
)

var _ = []any{nailRun131, nailResident131, nailModels131, nailSecret131}
