package main

// Ticket 273, wave 1: the nail under the D43 machine's first real Sink.
//
// WHAT IS BEING ASSERTED, and why each clause is a fact about the product
// rather than a fact about this file.
//
//  1. THE MACHINE IS THE PRODUCT'S. The case does not build a D43 machine. It
//     calls modelStore.handOffModel - the only production walker of
//     models.WireDownloading outside internal/models - and that function calls
//     handOffMachine273 (models.go), the same named constructor the product
//     uses. Ticket 273's R1 exists because the previous shape was an anonymous
//     literal with no grip: a case that could only rebuild the same shape would
//     be asserting on the machine it just made, which passes no matter what
//     the product does. Deleting the Sink line inside handOffMachine273 changes
//     the machine this case drives.
//
//  2. THE RECORD IS READ OFF THE DISK, not off a captured variable. The
//     positive control installs the leg's real persistent listener
//     (installLogSink, the same call models.go:339 makes), hands handOffModel
//     the same writer modelsEnsure hands it, and reads back
//     <data dir>\logs\wisp-*.jsonl. The assertion is: EXACTLY ONE effect record,
//     with d43=37 and name=model.verify-sha256-signature. Nothing in it reads a
//     field off a struct this test could have filled in.
//
//  3. THE WALK REALLY SUCCEEDED. Row #37's FAILED branch carries no SideEffects
//     (internal/statemachine/table.go:267), so "a record exists" must not be
//     satisfied by a failed walk. The case pins code=0 and the state=FirstRun
//     sentence on stdout first, and only then counts records: a walk that ends
//     in Error would have to fake both to pass.
//
// WHAT THIS FILE DELIBERATELY DOES NOT CLAIM.
//
//   - It does not claim C29's signature gate is exercised by this fixture. The
//     production assembly (openModelStore) loads the committed manifest under
//     buildinfo's key and cannot be pointed at a fixture - that key's private
//     half is not in this repository - so the manifest here is signed with a
//     THROWAWAY key and loaded through the REAL LoadSignedManifest, which still
//     runs the full minisign verification (it just verifies against this
//     fixture's key). Signature coverage itself is internal/models
//     manifest_real_test.go's claim, not this file's.
//   - It does not claim the effect EXECUTES. `wisp models ensure` already
//     re-verified the tree at the hand-off (internal/models/bridge.go:58); this
//     ticket books that as a record. See effect_audit_273.go's header.
//   - It does not claim the reverse-control walk below is the product's. That
//     case exists to prove the ASSERTION has teeth (the "empty can" shape), and
//     it says so in its own failure message.
//
// WHY THE logf CLOSURE IS RE-SPELLED HERE (three lines, models.go:344-346).
// modelsEnsure cannot be driven in this package: it resolves the manifest
// through models.ResolveManifestPath and verifies it against buildinfo's
// embedded key, so the only manifest it can read is the committed one, and no
// model in it can be installed from bytes a test owns. What the nail therefore
// holds on to is the machine (point 1) and the listener (point 2); the closure
// in between is modelsEnsure's own line copied verbatim, and if that line ever
// moves to slog.Default() this case goes on passing on its copy. That gap is
// registered here rather than papered over: ticket 273 R2 rules the channel, and
// the mutation that this file DOES catch is the one that matters ("no Sink at
// all"), which is measured red in the ticket's evidence.

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/models"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

const (
	// effectAuditModelID273 is a fixture id, not a manifest id: the whole
	// manifest is signed by this file (see newEffectAuditRig273).
	effectAuditModelID273 = "effect-audit-fixture-273"
	// effectAuditName273 / effectAuditRow273 are row #37's success branch as
	// the frozen table spells it (internal/statemachine/table.go:264-266). If
	// either moves, this case must go red and somebody has to say where the
	// ticket's headline assertion went.
	effectAuditName273 = "model.verify-sha256-signature"
	effectAuditRow273  = 37
	// effectAuditHandOffOK273 is the production walker's own record, spelled
	// the way models.go logs it. It is the second half of "this file is this
	// drive's": an effect record alone could come from any machine.
	effectAuditHandOffOK273 = "models: hand-off OK id=" + effectAuditModelID273
)

// effectRecord273 is one audit line parsed back into the three fields ticket
// 273 books. `at` is uint64 to match the counter's width, so a negative or a
// clock-shaped value cannot be parsed into it.
type effectRecord273 struct {
	D43  int
	Name string
	At   uint64
}

func (r effectRecord273) String() string {
	return fmt.Sprintf("{d43=%d name=%s at=%d}", r.D43, r.Name, r.At)
}

// parseEffectRecord273 turns one log message back into the audit record it is
// supposed to be, and refuses anything that is not exactly that shape. Parsing
// the fields (instead of searching the raw line for "d43=37") is what makes the
// count a statement about the FIELDS: a message that mentions 37 in prose, or
// carries a fourth key, is not an audit record and does not count.
func parseEffectRecord273(msg string) (effectRecord273, bool) {
	rest, ok := strings.CutPrefix(msg, "models: "+effectAuditRecordPrefix273)
	if !ok {
		return effectRecord273{}, false
	}
	var (
		rec                       effectRecord273
		seenD43, seenName, seenAt bool
	)
	fields := strings.Fields(rest)
	if len(fields) != 3 {
		return effectRecord273{}, false
	}
	for _, f := range fields {
		key, value, ok := strings.Cut(f, "=")
		if !ok {
			return effectRecord273{}, false
		}
		switch key {
		case "d43":
			n, err := strconv.Atoi(value)
			if err != nil {
				return effectRecord273{}, false
			}
			rec.D43, seenD43 = n, true
		case "name":
			if value == "" {
				return effectRecord273{}, false
			}
			rec.Name, seenName = value, true
		case "at":
			n, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return effectRecord273{}, false
			}
			rec.At, seenAt = n, true
		default:
			return effectRecord273{}, false
		}
	}
	if !seenD43 || !seenName || !seenAt {
		return effectRecord273{}, false
	}
	return rec, true
}

// auditLines273 reads every message the leg's own listener wrote under dataDir,
// in file order, and refuses to return an empty set: "the directory exists and
// says nothing" is a failure of the instrument, not a pass for the assertion.
func auditLines273(t *testing.T, dataDir string) []string {
	t.Helper()
	dir := logSinkDir(dataDir)
	names, err := filepath.Glob(filepath.Join(dir, "wisp-*.jsonl"))
	if err != nil {
		t.Fatalf("glob %s: %v", dir, err)
	}
	if len(names) != 1 {
		t.Fatalf("pipeline files in %s = %d, want 1 (one drive of one leg, one file); names: %v",
			dir, len(names), names)
	}
	raw, err := os.ReadFile(names[0])
	if err != nil {
		t.Fatalf("read %s: %v", names[0], err)
	}
	var out []string
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var rec struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("line %q in %s is not a JSON object: %v", line, names[0], err)
		}
		out = append(out, rec.Msg)
		// The raw line goes into the case's own log on purpose: the ticket's
		// evidence is a record read back off a disk, and a reader of the test
		// log should be able to see the JSONL line itself rather than this
		// file's summary of it.
		t.Logf("raw audit line in %s: %s", filepath.Base(names[0]), line)
	}
	if len(out) == 0 {
		t.Fatalf("%s holds a pipeline file with no record in it", dir)
	}
	return out
}

// effectRecords273 keeps only the messages that parse as audit records. The
// count of what comes back is the number the ticket's assertion counts.
func effectRecords273(t *testing.T, lines []string) []effectRecord273 {
	t.Helper()
	var out []effectRecord273
	for _, line := range lines {
		if rec, ok := parseEffectRecord273(line); ok {
			out = append(out, rec)
		}
	}
	return out
}

// assertExactlyOneEffect273 is the ONE assertion this ticket is judged by, and
// every case below applies this same function:
//
//	the leg's JSONL holds exactly one record whose d43 is 37 and whose name is
//	model.verify-sha256-signature.
//
// It takes a testing.TB so the reverse control can hand it an instrument that
// records a failure instead of aborting - the same call, the same words, so
// "this assertion is red on an empty sink" is a measured statement and not a
// second, weaker assertion written to pass.
func assertExactlyOneEffect273(tb testing.TB, recs []effectRecord273, d43 int, name string) {
	tb.Helper()
	var hits []effectRecord273
	for _, rec := range recs {
		if rec.D43 == d43 && rec.Name == name {
			hits = append(hits, rec)
		}
	}
	if len(hits) != 1 {
		tb.Fatalf("audit records with d43=%d name=%s = %d, want exactly 1; every effect record read back off the disk: %v",
			d43, name, len(hits), recs)
	}
}

// redExpectation273 stands in for *testing.T where a FAILURE is the expected
// outcome. Embedding the real *testing.T supplies the interface's unexported
// methods; Fatalf/Errorf are replaced so a case can observe that the assertion
// under test would have failed without the case itself going red.
type redExpectation273 struct {
	*testing.T
	failed bool
	why    string
}

func (r *redExpectation273) Fatalf(format string, args ...any) {
	r.failed = true
	r.why = fmt.Sprintf(format, args...)
}

func (r *redExpectation273) Errorf(format string, args ...any) {
	r.failed = true
	r.why = fmt.Sprintf(format, args...)
}

// effectAuditRig273 is one downloadable model nobody has to download: a
// fixture-signed manifest whose single artifact is pinned to bytes this file
// owns, plus a local_override directory holding exactly those bytes. Ensure's
// override branch is a production branch (SPEC-04 7.1 sideload), it re-verifies
// the tree before handing it back, and it touches no network - the same
// no-network discipline ticket 131's model case states for the same reason.
type effectAuditRig273 struct {
	dataDir string
	store   *modelStore
	id      string
}

func newEffectAuditRig273(t *testing.T) effectAuditRig273 {
	t.Helper()

	body := []byte("wisp ticket 273 effect-audit fixture bytes\n")
	sum := sha256.Sum256(body)
	manifestJSON, err := json.Marshal(&models.Manifest{
		ManifestVersion: 1,
		Models: []models.ModelEntry{{
			ID:        effectAuditModelID273,
			Purpose:   "vad",
			Quant:     "int8",
			License:   "Apache-2.0",
			URLs:      []string{"https://example.invalid/wisp-273/vad.onnx"},
			SHA256:    hex.EncodeToString(sum[:]),
			SizeBytes: int64(len(body)),
		}},
	})
	if err != nil {
		t.Fatalf("marshal fixture manifest: %v", err)
	}

	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestJSON, 0o600); err != nil {
		t.Fatalf("write fixture manifest: %v", err)
	}
	// The seed is 32 bytes, checked rather than trusted: a typo here would
	// otherwise surface as ed25519's "bad seed length" from inside the fixture.
	seed := []byte("wisp-273-effect-audit-fixture!!!")
	if len(seed) != ed25519.SeedSize {
		t.Fatalf("fixture seed is %d bytes, want %d", len(seed), ed25519.SeedSize)
	}
	priv := ed25519.NewKeyFromSeed(seed)
	pubKey := priv.Public().(ed25519.PublicKey)
	fixturePub := models.MinisignPublicKey{KeyID: models.DeriveKeyID(pubKey), Key: pubKey}
	sig, err := models.SignPayload(priv, fixturePub, manifestJSON,
		"ticket 273 fixture", "ticket 273 effect-audit fixture")
	if err != nil {
		t.Fatalf("sign fixture manifest: %v", err)
	}
	if err := os.WriteFile(manifestPath+".minisig", sig, 0o600); err != nil {
		t.Fatalf("write fixture signature: %v", err)
	}
	// Real verification path, fixture key: see the header's "WHAT THIS FILE
	// DELIBERATELY DOES NOT CLAIM".
	manifest, err := models.LoadSignedManifest(manifestPath, manifestPath+".minisig",
		fixturePub.String("ticket 273 fixture key"))
	if err != nil {
		t.Fatalf("LoadSignedManifest over the fixture: %v", err)
	}

	overrideDir := filepath.Join(dir, "override")
	if err := os.MkdirAll(overrideDir, 0o700); err != nil {
		t.Fatalf("create override dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(overrideDir, "vad.onnx"), body, 0o600); err != nil {
		t.Fatalf("write override bytes: %v", err)
	}

	dataDir := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		t.Fatalf("create data dir: %v", err)
	}
	root := filepath.Join(dataDir, "models")
	mgr, err := models.NewManager(models.Options{
		DataDir:         root,
		Manifest:        manifest,
		VerifySignature: true,
		LocalOverride:   map[string]string{effectAuditModelID273: overrideDir},
		Attempts:        1,
	})
	if err != nil {
		t.Fatalf("NewManager over the fixture: %v", err)
	}
	return effectAuditRig273{
		dataDir: dataDir,
		store:   &modelStore{mgr: mgr, manifest: manifest, root: root, src: manifestPath},
		id:      effectAuditModelID273,
	}
}

// auditLogf273 is the leg's persistent listener installed exactly as
// modelsEnsure installs it, plus the writer that listener gives the walk.
// stop is explicit and idempotent because the caller reads the file it flushes:
// the ordering comes from the code under test closing its own sink, not from a
// wait (ticket 131's rule for the same legs).
type auditLogf273 struct {
	logf func(string, ...any)
	stop func()
}

func installAuditLogf273(t *testing.T, dataDir string) auditLogf273 {
	t.Helper()
	oldDefault := slog.Default()
	sink, err := installLogSink(dataDir)
	if err != nil {
		slog.SetDefault(oldDefault)
		t.Fatalf("installLogSink(%s): %v", dataDir, err)
	}
	sinkLogger := sink.logger()
	stopped := false
	stop := func() {
		if stopped {
			return
		}
		stopped = true
		sink.close()
		slog.SetDefault(oldDefault)
	}
	t.Cleanup(stop)
	return auditLogf273{
		// models.go:344-346, verbatim: the tag is what makes one JSONL file
		// readable as one leg's records.
		logf: func(format string, args ...any) {
			sinkLogger.Info("models: " + fmt.Sprintf(format, args...))
		},
		stop: stop,
	}
}

// TestEffectAudit273HandOffBooksExactlyOneRow37 is the positive control. It
// drives the PRODUCT's walker (modelStore.handOffModel), which calls the
// product's machine constructor, and reads the verdict off the disk.
func TestEffectAudit273HandOffBooksExactlyOneRow37(t *testing.T) {
	rig := newEffectAuditRig273(t)
	audit := installAuditLogf273(t, rig.dataDir)

	var stdout, stderr strings.Builder
	code := rig.store.handOffModel(
		modelsIO{stdout: &stdout, stderr: &stderr, dataDir: rig.dataDir}, rig.id, audit.logf)
	audit.stop() // flush before reading; the leg closes its own listener

	// The walk's verdict first: row #37's success branch is the ONLY branch of
	// that row that fires an effect, so a case that counted records without
	// pinning this could be satisfied by a walk that never reached it.
	if code != 0 {
		t.Fatalf("handOffModel(%s) exited %d, want 0 (successful hand-off); stderr:\n%s",
			rig.id, code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "state=FirstRun") {
		t.Fatalf("the walk's own verdict does not end in FirstRun:\n%s", stdout.String())
	}

	lines := auditLines273(t, rig.dataDir)
	recs := effectRecords273(t, lines)

	// Exactly one effect record in the whole file: "at least one" is the shape
	// a double-installed listener or a double-fired row would sneak past.
	if len(recs) != 1 {
		t.Fatalf("effect records on disk = %d, want exactly 1 (row #2 has no SideEffects and row #37 fires once); messages: %v",
			len(recs), lines)
	}
	assertExactlyOneEffect273(t, recs, effectAuditRow273, effectAuditName273)
	if recs[0].At != 1 {
		t.Errorf("the walk's single effect recorded at=%d, want 1 (first effect of this machine's counter)", recs[0].At)
	}
	// The record must sit in the file the leg's own writer produced, next to the
	// walk's verdict - not in some other file that happens to be in the tree.
	foundOK := false
	for _, line := range lines {
		if strings.HasPrefix(line, effectAuditHandOffOK273) {
			foundOK = true
		}
	}
	if !foundOK {
		t.Fatalf("no %q record among the messages this drive wrote: the file under test is not this walk's; messages: %v",
			effectAuditHandOffOK273, lines)
	}
	t.Logf("D43 AUDIT RECORD ON DISK: %v (file dir=%s)", recs[0], logSinkDir(rig.dataDir))
}

// TestEffectAudit273EmptySinkIsRed is the reverse control ticket 273 asks for:
// the same walk, the same assertion, but the machine's Sink does nothing with
// the effect. The assertion MUST go red, or it is a sentence that would pass
// for a machine wired to a jar.
//
// This case says nothing about the product (handOffMachine273 is not used
// here); it is here to prove the assertion above has teeth, which is the one
// property a green control cannot show on its own.
func TestEffectAudit273EmptySinkIsRed(t *testing.T) {
	rig := newEffectAuditRig273(t)
	audit := installAuditLogf273(t, rig.dataDir)

	// The can: it receives the effect and writes nowhere. It keeps the effects
	// in memory so this case can prove the walk DID fire row #37 - otherwise
	// "no record on disk" would also be explained by "the walk never got there",
	// and the control would pass for the wrong reason.
	var fired []statemachine.Effect
	machine := statemachine.New(statemachine.Options{
		Initial: statemachine.StateFirstRun,
		Sink:    func(e statemachine.Effect) { fired = append(fired, e) },
	})
	defer machine.Close()
	bridge := models.WireDownloading(rig.store.mgr, machine)
	if _, err := bridge.Run(context.Background(), rig.id); err != nil {
		t.Fatalf("walk over the fixture failed, so the control would be vacuous: %v", err)
	}
	audit.stop()

	if len(fired) != 1 || fired[0].D43 != effectAuditRow273 || fired[0].Name != effectAuditName273 {
		t.Fatalf("the walk fired %v, want exactly row #37's %s; the control needs the effect to have fired so that "+
			"'no record' means 'the sink wrote nothing'", fired, effectAuditName273)
	}

	recs := effectRecords273(t, auditLines273(t, rig.dataDir))
	want := &redExpectation273{T: t}
	assertExactlyOneEffect273(want, recs, effectAuditRow273, effectAuditName273)
	if !want.failed {
		t.Fatalf("the ticket's assertion PASSED on a sink that writes nothing (records: %v): "+
			"it is a tautology and proves nothing about the product", recs)
	}
	t.Logf("empty-can control red as required: %s", want.why)
}

// TestEffectAudit273SinkFieldsAndCounter pins what the audit sink itself books,
// because the on-disk case can only ever see ONE effect (one walk, one row):
// the counter's monotonicity and the exact field spelling need more than one
// call, and this is that measurement. The nil-logf branch is the degraded
// listener path - it must count and write nothing, not panic and not invent a
// second channel.
func TestEffectAudit273SinkFieldsAndCounter(t *testing.T) {
	var lines []string
	sink := effectAuditSink273(func(format string, args ...any) {
		lines = append(lines, "models: "+fmt.Sprintf(format, args...))
	})
	for i := 0; i < 3; i++ {
		sink(statemachine.Effect{D43: effectAuditRow273, Name: effectAuditName273})
	}

	want := []effectRecord273{
		{D43: effectAuditRow273, Name: effectAuditName273, At: 1},
		{D43: effectAuditRow273, Name: effectAuditName273, At: 2},
		{D43: effectAuditRow273, Name: effectAuditName273, At: 3},
	}
	if len(lines) != len(want) {
		t.Fatalf("sink wrote %d line(s), want %d: %v", len(lines), len(want), lines)
	}
	for i, line := range lines {
		got, ok := parseEffectRecord273(line)
		if !ok {
			t.Fatalf("line %d %q does not parse as an audit record", i, line)
		}
		if got != want[i] {
			t.Errorf("line %d = %v, want %v", i, got, want[i])
		}
	}

	// Degraded listener branch: effectAuditSink273(nil) is what the walk gets
	// when installLogSink refused (modelsEnsure has already said so on stderr).
	// The only thing measurable from here is that the branch is CALLABLE - with
	// no writer there is nowhere a record could land, so "it wrote nothing" is
	// not a measurement. What this call kills is the mutation that drops the nil
	// guard: a nil func call panics, and the degraded leg would take the walk
	// down with it instead of degrading.
	effectAuditSink273(nil)(statemachine.Effect{D43: effectAuditRow273, Name: effectAuditName273})
}
