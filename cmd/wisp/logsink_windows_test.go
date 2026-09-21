//go:build windows

package main

// Ticket 117 AC#2/AC#4 in the shape the ticket demanded: not "a logger is
// wired up", but "this named WARN, with these fields, is greppable in this
// file on disk, produced by the production assembly root".
//
// Both cases below drive `runTextTask`, which is what `wisp run` calls
// (main.go, case "run"). Neither builds a sink by hand: the only sink either
// case reads is the one the production path installed itself, and the only
// trigger of the notice is a sealing site the assembly root reaches on its own
// first statement (secret.NewStore -> winsec.PrivateDirAll). That distinction
// is the whole point of the ticket: R-104-5 counted the existing proofs as
// zero because each one installed its own handler, and a sink built in a test
// is not a production outlet.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/observe"
)

// sealNoticeMsg is the event name internal/winsec writes when a seal clears a
// principal that stood on an object out of band. It is copied here as a string
// on purpose: that message is the promise this ticket exists to keep, so if the
// text ever changes this case goes red and somebody has to say where the record
// that now says "you lost an authorization" lives. That is the outcome a
// grep-able AC is supposed to produce.
const sealNoticeMsg = "winsec: seal cleared principals that stood on this object"

// sinkRecord is one JSONL line as the pipeline writes it (slog JSONHandler
// behind redactHandler). Only the fields asserted here are decoded.
type sinkRecord struct {
	Time             string `json:"time"`
	Level            string `json:"level"`
	Msg              string `json:"msg"`
	Path             string `json:"path"`
	Kind             string `json:"kind"`
	Cleared          string `json:"cleared"`
	ClearedInherited string `json:"cleared_inherited"`
	Policy           string `json:"policy"`
}

// readSink decodes the pipeline's own files in dir and refuses if the directory
// does not hold exactly what observe itself counts as pipeline files: a test
// that globbed any *.jsonl would still pass if the pipeline renamed its files,
// and "which file, which field" is the whole content of AC#2.
func readSink(t *testing.T, dir string) []sinkRecord {
	t.Helper()
	n, err := observe.CountLogFiles(dir)
	if err != nil {
		t.Fatalf("counting pipeline files in %s: %v", dir, err)
	}
	if n == 0 {
		t.Fatalf("no wisp-<day>-<seq>.jsonl in %s: the sink wrote nothing where the pipeline names it", dir)
	}
	if n != 1 {
		t.Errorf("pipeline files in %s = %d, want 1 (one boot, one file: no roll, no second writer)", dir, n)
	}
	var out []sinkRecord
	total := 0
	for _, name := range mustGlob117(t, filepath.Join(dir, "wisp-*.jsonl")) {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		total += len(b)
		for _, line := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var rec sinkRecord
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("line %q in %s is not a JSON object: %v", line, name, err)
			}
			if rec.Msg == "" {
				t.Errorf("record %q carries no msg field", line)
			}
			out = append(out, rec)
		}
	}
	t.Logf("sink dir %s: %d file(s), %d bytes, %d record(s)", dir, n, total, len(out))
	return out
}

func mustGlob117(t *testing.T, pattern string) []string {
	t.Helper()
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	return matches
}

// icacls117 runs the ACL tool the only way this repository accepts as evidence
// (ticket 89's ruling: FileInfo.Mode() reports the fiction the mode argument
// was supposed to control). The name is ticket-qualified because package main's
// other test files belong to other tickets.
func icacls117(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("icacls", args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		t.Fatalf("icacls %v: %v\n%s", args, err, errb.String())
	}
	return out.String()
}

func namesEveryone117(s string) bool {
	for _, tok := range []string{"S-1-1-0", "Everyone", "WD"} {
		if strings.Contains(s, tok) {
			return true
		}
	}
	return false
}

// TestAC2SealNoticeLandsInTheRunLegLogFile is AC#2's headline row and AC#4's
// in-process half: an out-of-band grant on the DPAPI blob directory is cleared
// by the production assembly root, and the record of it comes off disk.
//
// Landing spot, verbatim: <data root>\logs\wisp-<YYYYMMDD>-001.jsonl, the line
// whose "msg" field is exactly sealNoticeMsg, carrying "level":"WARN" plus the
// path / kind / cleared / cleared_inherited / policy fields.
func TestAC2SealNoticeLandsInTheRunLegLogFile(t *testing.T) {
	dataDir := t.TempDir()
	secrets := filepath.Join(dataDir, "secrets")
	if err := os.MkdirAll(secrets, 0o700); err != nil {
		t.Fatal(err)
	}
	// Seed the out-of-band grant an operator would have added ("a service
	// account needs to read the store"), then prove the seed landed: asserting
	// on a grant icacls never applied would prove nothing.
	icacls117(t, secrets, "/grant", "*S-1-1-0:(OI)(CI)(RX)")
	if seeded := icacls117(t, secrets); !namesEveryone117(seeded) {
		t.Fatalf("the seed grant did not land, this case would prove nothing:\n%s", seeded)
	}

	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	code := runTextTask(runSpec{
		argv:    []string{"总结一下 这份笔记"},
		stdout:  out,
		stderr:  errb,
		dataDir: dataDir,
		notify:  func(string, string) error { return nil },
	})
	// There is no config.toml here, so the run stops at Unconfigured (exit 2).
	// That is the point of the case rather than an accident of it: the notice
	// has to survive a boot that gets no further than opening the store, which
	// is exactly the record a terminal-only outlet throws away.
	if code != 2 {
		t.Fatalf("exit %d, want 2 (Unconfigured)\nstdout:\n%s\nstderr:\n%s", code, out.String(), errb.String())
	}

	recs := readSink(t, logSinkDir(dataDir))
	var hits []sinkRecord
	for _, r := range recs {
		if r.Msg == sealNoticeMsg {
			hits = append(hits, r)
		}
	}
	if len(hits) != 1 {
		t.Fatalf("records naming %q = %d, want 1; all records: %+v", sealNoticeMsg, len(hits), recs)
	}
	hit := hits[0]
	if hit.Level != "WARN" {
		t.Errorf("level = %q, want WARN (the pipeline's level filter must not swallow it)", hit.Level)
	}
	assertNamesTree117(t, hit.Path, secrets)
	if !namesEveryone117(hit.Cleared) {
		t.Errorf("cleared = %q, want the removed principal named", hit.Cleared)
	}
	// The bucket the seed fills is "explicit" on every host; whether the temp
	// root also hands the directory an inherited copy varies with the machine
	// (ticket 95 measured hosts where %TEMP%'s parent grants BUILTIN\Users),
	// so the claim pinned here is "the grant standing on this object in its own
	// right got named", which is the sentence the operator needs.
	if !strings.HasPrefix(hit.Kind, "explicit") {
		t.Errorf("kind = %q, want it to start with explicit (the grant stood on this object in its own right)", hit.Kind)
	}
	if !strings.Contains(hit.Policy, "winsec owns the grants") {
		t.Errorf("policy = %q, want the standing policy sentence", hit.Policy)
	}
	t.Logf("RECORDED ON DISK: level=%s msg=%s path=%s kind=%s cleared=%q",
		hit.Level, hit.Msg, hit.Path, hit.Kind, hit.Cleared)

	// And the removal really happened, so the record is not describing a no-op.
	if after := icacls117(t, secrets); namesEveryone117(after) {
		t.Errorf("the grant is still there after the notice said it was cleared:\n%s", after)
	}
}

// assertNamesTree117 checks that the notice is about the directory the test
// widened. It compares components, not spellings, because the path in the
// record is ResolvePath's answer and not the caller's text (ticket 115's
// ruling): on a host whose profile directory is spelled 8.3 (C:\Users\RUNNER~1,
// measured on CI run 35599458439) an exact == and a strings.EqualFold both
// fail on objects that are plainly the same tree. Comparing the leaf, its
// parent and the volume keeps the assertion honest above the temp root without
// reimplementing winsec's tree rule in a test.
func assertNamesTree117(t *testing.T, got, want string) {
	t.Helper()
	base, wantBase := filepath.Clean(got), filepath.Clean(want)
	if !strings.EqualFold(filepath.Base(base), filepath.Base(wantBase)) {
		t.Fatalf("notice path %q does not name the sealed dir %q", got, want)
	}
	if !strings.EqualFold(filepath.Base(filepath.Dir(base)), filepath.Base(filepath.Dir(wantBase))) {
		t.Fatalf("notice path %q is not a child of the temp root %q", got, wantBase)
	}
	if !strings.EqualFold(filepath.VolumeName(base), filepath.VolumeName(wantBase)) {
		t.Fatalf("notice path %q is on another volume than %q", got, want)
	}
}

// TestAC2AuditTrailLandsInTheRunLegLogFile answers R-105-1, the other half of
// the same root: ticket 105 booked the path-rewrite ledger into
// agentRuntime.auditf, and R-105-1 was right that auditf was an
// fmt.Fprintf(rt.stderr, ...) - so "the rewrite ledger is on disk" was true
// only of a test's own sink. This case drives one complete run and then requires
// EVERY audit line the run printed to stderr to have a twin record in the file,
// which is a stronger claim than picking two favourites.
func TestAC2AuditTrailLandsInTheRunLegLogFile(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	f.rtHook = func(rt *agentRuntime) {
		source := filepath.Join(f.dir, "note.txt")
		if err := os.WriteFile(source, []byte("notes body"), 0o600); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// One host-side call through the assembled bridge: its argument goes
		// through C26, so the path account has something to report and the
		// "tools: PATH-ACCOUNT" line is produced for real, on the production
		// funnel (tools.Options.Logf -> agentRuntime.auditf).
		out, err := rt.bridge.Execute(ctx, agent.ToolRequest{
			TaskID: "host-audit-task", CorrelationID: "host-audit-corr", CallID: "c1",
			Name: "fs.read",
			Args: json.RawMessage(fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(source))),
		})
		if err != nil {
			t.Fatal(err)
		}
		if out.IsError {
			t.Fatalf("an allowlisted read must execute, got: %s", out.Text)
		}
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, f.out.String(), f.err.String())
	}

	var printed []string
	for _, line := range strings.Split(f.err.String(), "\n") {
		if rest, ok := strings.CutPrefix(line, "[audit] "); ok {
			printed = append(printed, rest)
		}
	}
	if len(printed) < 2 {
		t.Fatalf("only %d audit line(s) reached stderr, this case needs a real trail to compare:\n%s",
			len(printed), f.err.String())
	}
	seen := map[string]int{}
	for _, r := range readSink(t, logSinkDir(f.dir)) {
		seen[r.Msg]++
	}
	for _, line := range printed {
		if seen["audit: "+line] == 0 {
			t.Errorf("audit line printed to stderr has no record in the log file: %q", line)
		}
	}
	// Two named rows out of that set, because AC#2 wants event names and not
	// only a count: the mode this boot came up in, and the rewrite ledger.
	if !anyPrefix117(seen, "audit: perm: MODE-READ") {
		t.Error("no \"perm: MODE-READ\" record in the log file: the boot posture is not auditable")
	}
	if !anyPrefix117(seen, "audit: tools: PATH-ACCOUNT") {
		t.Error("no \"tools: PATH-ACCOUNT\" record in the log file: the rewrite ledger is still only on stderr")
	}
}

func anyPrefix117(seen map[string]int, prefix string) bool {
	for msg, n := range seen {
		if n > 0 && strings.HasPrefix(msg, prefix) {
			return true
		}
	}
	return false
}

// TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite is AC#2's ordering half
// stated as its own case: the first thing assembleRuntime does is seal the blob
// directory, so that seal's notice can only be in the file if the install
// happened earlier. Moving installLogSink below assembleRuntime keeps every
// other test in this package green and takes this one red, which is the point -
// a sink installed after the first event loses exactly the record it exists to
// keep.
func TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite(t *testing.T) {
	dataDir := t.TempDir()
	secrets := filepath.Join(dataDir, "secrets")
	if err := os.MkdirAll(secrets, 0o700); err != nil {
		t.Fatal(err)
	}
	icacls117(t, secrets, "/grant", "*S-1-1-0:(OI)(CI)(RX)")
	if err := os.WriteFile(filepath.Join(secrets, "blob"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	runTextTask(runSpec{
		argv:    []string{"x"},
		stdout:  out,
		stderr:  errb,
		dataDir: dataDir,
		notify:  func(string, string) error { return nil },
	})
	recs := readSink(t, logSinkDir(dataDir))
	first := -1
	for i, r := range recs {
		if r.Msg == sealNoticeMsg {
			first = i
			break
		}
	}
	if first < 0 {
		t.Fatalf("no seal notice recorded at all; stderr:\n%s\nrecords: %+v", errb.String(), recs)
	}
	if first != 1 {
		t.Errorf("seal notice is record %d, want 1: only the install record may precede it; records: %+v", first, recs)
	}
	if first >= 1 && !strings.Contains(recs[0].Msg, "log sink installed") {
		t.Errorf("record 0 = %q, want the install record that names the sink dir", recs[0].Msg)
	}
}
