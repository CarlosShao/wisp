package main

// Ticket 39 AC#2's production e2e (write leg 2026-11-18): the section-write
// loosening gate, driven end to end through the REAL approval gate, the REAL
// native card answer, and a REAL config.toml on disk - the same assembly
// `wisp run` wires (startConfigLoosenGate over Manager.LoosenQueue), the same
// card shape ticket 21's surface shows, judged on the FILE by reading it back.
//
// WHAT MAKES THESE READINGS PRODUCTION. No field in this file fakes the gate:
// the card is raised by the adapter cmd/wisp wires (config_loosen_gate_39.go),
// the allow arrives through Gate.Native().Allow with the grant the prompt
// carried (the only door an allow has - ban #6), the reject arrives through
// the gate's own Reject, the timeout is the gate's own C18 deadline
// (ApprovalTimeout), and every assertion about "applied" / "kept old" reads
// config.toml back from disk. Deleting the wiring (rt.startConfigLoosenGate or
// the LoosenQueue assignment inside it) reddens every case here with the
// unwired-queue refusal - the mutation reading the ticket asks for.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
)

// loosen39Fixture is a config.toml seed with one disabled plugin row: the
// native ON toggle is a real locked-section loosening.
const loosen39Fixture = "[risk]\npermission_mode = \"ask_every_step\"\n" +
	"[plugins]\ntier2_enabled = false\n" +
	"[plugins.zed]\nenabled = false\n"

// newLoosen39Runtime builds the smallest honest production assembly: a real
// config.toml on a real Manager, a real gate over the scripted native surface,
// and the wiring cmd/wisp's assembly root runs - nothing assigned by hand that
// startConfigLoosenGate assigns in production.
func newLoosen39Runtime(t *testing.T, timeout time.Duration) (*agentRuntime, string, *nativeCard371, *auditCapture39) {
	t.Helper()
	dir := sealableTempDir124(t)
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(loosen39Fixture), 0o600); err != nil {
		t.Fatalf("cannot seed %s: %v", path, err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	ui := newNativeCard371()
	// The Manager's own adjudication lines (SECTION-WRITE) and this runtime's
	// lines share one capture, shaped the way auditf spells them - the same
	// audit trail a run keeps, one buffer for the joins the cases read. The
	// capture is lock-serialized (the deflake half of the 2026-11-18 leg):
	// the Manager audit runs on the caller's goroutine while the resolver's
	// lines arrive on the spawned approval-waiter goroutine, and a shared
	// unsynchronized bytes.Buffer was both a race detector red and the reason
	// a torn line could be unread while a case was already judging the dump.
	audit := newAuditCapture39()
	mgr.Audit = audit.auditf
	rt := &agentRuntime{
		spec:   runSpec{dataDir: dir},
		stderr: audit,
		mgr:    mgr,
		gate: approval.New(approval.Options{
			UI:              ui,
			ApprovalTimeout: timeout,
			Logf:            func(string, ...any) {},
		}),
	}
	rt.startConfigLoosenGate()
	if mgr.LoosenQueue == nil {
		t.Fatal("startConfigLoosenGate did not wire the loosening queue; the gate is unwired")
	}
	return rt, path, ui, audit
}

// auditCapture39 is this file's shared audit sink. fmt.Fprintf to an
// io.Writer is not atomic and a bytes.Buffer is not thread-safe: two
// goroutines write here (the Manager's SECTION-WRITE lines on the test's own
// goroutine, the adapter's LOOSEN-GATE lines on the approval-waiter spawned
// by QueueLoosen - cmd/wisp/config_loosen_gate_39.go), so every access takes
// the lock. This is the deterministic-arrival fix, not an assertion
// relaxation: nothing is skipped, nothing is reworded, every line still
// lands, and the waiters below now read whole lines instead of racing the
// writer mid-Flush.
type auditCapture39 struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func newAuditCapture39() *auditCapture39 { return &auditCapture39{} }

// auditf is the Manager.Audit hook shape (internal/config sectionwrite.go's
// own call shape, [audit]-prefixed the way auditf spells it).
func (c *auditCapture39) auditf(f string, a ...any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Fprintf(&c.buf, "[audit] "+f+"\n", a...)
}

// Write is the io.Writer half agentRuntime.auditf rides (fmt.Fprint(rt.stderr, ...)).
func (c *auditCapture39) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.Write(p)
}

// String snapshots the capture under the same lock.
func (c *auditCapture39) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.String()
}

// readZedEnabled reads the ONE field this ticket owns back off the disk - the
// assertion a reviewer can re-run against the bytes.
func readZedEnabled(t *testing.T, path string) bool {
	t.Helper()
	m, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("read-back of %s: %v", path, err)
	}
	return m.Config().Plugins.Entries["zed"].Enabled
}

// waitForZed polls the file read-back until the predicate holds or the budget
// ends - the spawned resolver writes on its own goroutine, so the test waits
// for a fact instead of sleeping a guess.
func waitForZed(t *testing.T, path string, want bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if readZedEnabled(t, path) == want {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("config.toml never reached plugins.zed.enabled=%v within the budget", want)
}

// auditLines39 returns the audit lines matching every needle.
func auditLines39(audit *auditCapture39, needles ...string) []string {
	var out []string
	for _, l := range strings.Split(audit.String(), "\n") {
		ok := true
		for _, n := range needles {
			if !strings.Contains(l, n) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, l)
		}
	}
	return out
}

// auditWaitBudget39 is the explicit upper bound of every audit wait in this
// file. The waiters are polls-with-deadline, never a bare check, never a
// sleep: the timeout IS the determinism ceiling the dispatch names, and on
// expiry the dump is printed whole so the red names what arrived.
const auditWaitBudget39 = 5 * time.Second

// waitForAudit39 polls the capture until one line matches every needle or the
// budget ends - the resolver writes its lines on its own goroutine, so the
// test waits for the record instead of sleeping a guess.
func waitForAudit39(t *testing.T, audit *auditCapture39, needles ...string) []string {
	t.Helper()
	deadline := time.Now().Add(auditWaitBudget39)
	for {
		if rows := auditLines39(audit, needles...); len(rows) > 0 {
			return rows
		}
		if !time.Now().Before(deadline) {
			t.Fatalf("no audit line matching %v within the budget:\n%s", needles, audit.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// waitForAuditCount39 is the count-shaped wait: poll until at least `want`
// lines match every needle, or the budget ends. The accepted case uses this
// for its exactly-one assertion: the resolver goroutine books the accepted
// SECTION-WRITE line after waitForZed's disk read-back can already see the
// applied value (the write commits before the audit line is printed), so
// waiting on the disk alone was the race the whole-package runs hit
// (docs/evidence/s1/37-depth-timeout-fallback-20261118.md: the dump showed
// outcome=queued while the case was already counting lines). Asserting
// >= want here and exactly-one after is the same assertion, arrived at
// deterministically - not a relaxation.
func waitForAuditCount39(t *testing.T, audit *auditCapture39, want int, needles ...string) []string {
	t.Helper()
	deadline := time.Now().Add(auditWaitBudget39)
	for {
		rows := auditLines39(audit, needles...)
		if len(rows) >= want {
			return rows
		}
		if !time.Now().Before(deadline) {
			t.Fatalf("audit lines matching %v reached %d, want at least %d, within the budget:\n%s",
				needles, len(rows), want, audit.String())
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// AC#2, accepted arm, production seam: the native ON toggle queues, the card
// reaches the native surface, the operator's native allow turns into a REAL
// write, and the one SECTION-WRITE audit line carries the card id and the
// old->new value.
func Test39NativeLooseningCardAcceptedAppliesAndAudits(t *testing.T) {
	rt, path, ui, audit := newLoosen39Runtime(t, 0)

	written, err := rt.mgr.SetPluginEnabled("zed", true, config.SourceNative, "cli.test")
	if err != nil || written != nil {
		t.Fatalf("queueing must return no receipt and no error: written=%v err=%v", written, err)
	}
	// Before the answer: the file still carries the old value (read back).
	if readZedEnabled(t, path) {
		t.Fatal("config.toml moved before any card was answered")
	}
	p := ui.prompt(t)
	if !strings.HasPrefix(p.CorrelationID, configLoosenTaskID+"#") {
		t.Errorf("the card id = %q, want the %q namespace the resolver joins on", p.CorrelationID, configLoosenTaskID)
	}
	if p.Tool != configLoosenToolName || p.Level != "L2" {
		t.Errorf("the card is tool=%q level=%q, want %q/L2", p.Tool, p.Level, configLoosenToolName)
	}
	cardID := p.CorrelationID
	if err := rt.gate.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("the native allow was refused: %v", err)
	}
	waitForZed(t, path, true)

	// The resolver books the accepted SECTION-WRITE line on its own goroutine
	// AFTER the commit, so the read-back above can already see the applied
	// value while the line is still in flight - the whole-package flake
	// (the 37-leg close run: the dump showed outcome=queued at judge time).
	// Wait for the record deterministically (poll with the explicit budget),
	// then judge exactly-one on the same set.
	rows := waitForAuditCount39(t, audit, 1, "SECTION-WRITE", "outcome=written", "card="+cardID)
	if len(rows) != 1 {
		t.Fatalf("accepted SECTION-WRITE lines = %d, want exactly 1:\n%s",
			len(rows), audit.String())
	}
	for _, want := range []string{
		"section=plugins", "direction=loosen", "source=native", "origin=cli.test",
		"changes=plugins.zed.enabled:false->true",
	} {
		if !strings.Contains(rows[0], want) {
			t.Errorf("the accepted write's audit line is missing %q:\n%s", want, rows[0])
		}
	}
	if lines := waitForAuditCount39(t, audit, 1, "LOOSEN-GATE state=answered card="+cardID+" ", "result=allow"); len(lines) != 1 {
		t.Errorf("no single LOOSEN-GATE allow line joins card %s:\n%s", cardID, audit.String())
	}
}

// AC#2, rejected arm: the operator's native REJECT writes nothing. The file is
// read back and still carries the old value, and the audit names the denial.
func Test39NativeLooseningCardRejectedKeepsOldValue(t *testing.T) {
	rt, path, ui, audit := newLoosen39Runtime(t, 0)

	if _, err := rt.mgr.SetPluginEnabled("zed", true, config.SourceNative, "cli.test"); err != nil {
		t.Fatalf("queueing: %v", err)
	}
	p := ui.prompt(t)
	if err := rt.gate.Native().Reject(p.CorrelationID, "测试拒绝"); err != nil {
		t.Fatalf("the native reject was refused: %v", err)
	}
	// Give the resolver its turn, then read the file back: the old value must
	// be what the disk still says.
	waitForZed(t, path, false)
	if readZedEnabled(t, path) {
		t.Fatal("a REJECTED card still wrote the loosening: plugins.zed.enabled is true on disk")
	}
	if rows := waitForAudit39(t, audit, "SECTION-WRITE", "outcome=denied", "card="+p.CorrelationID); len(rows) != 1 {
		t.Fatalf("denied SECTION-WRITE lines = %d, want exactly 1", len(rows))
	}
}

// AC#2, timed-out arm: nobody answers, the gate's own C18 deadline fires, and
// the timeout resolves exactly like a reject - the old values stand.
func Test39NativeLooseningCardTimeoutKeepsOldValue(t *testing.T) {
	rt, path, ui, audit := newLoosen39Runtime(t, 300*time.Millisecond)

	if _, err := rt.mgr.SetPluginEnabled("zed", true, config.SourceNative, "cli.test"); err != nil {
		t.Fatalf("queueing: %v", err)
	}
	_ = ui.prompt(t) // the card is up; nobody answers it
	if readZedEnabled(t, path) {
		t.Fatal("a timed-out card wrote the loosening before its deadline even fired")
	}
	waitForAudit39(t, audit, "outcome=denied card=host:config-loosen#")
	waitForZed(t, path, false)
	if readZedEnabled(t, path) {
		t.Fatal("a timed-out card still wrote the loosening")
	}
}

// The unwired twin: with the gate absent from the runtime, the queue is never
// wired and the Manager answers the fail-closed unwired error - the mutation
// reading for "摘掉 L2 门": a run that cannot ask cannot loosen.
func Test39LooseningWithoutGateFailsClosed(t *testing.T) {
	dir := sealableTempDir124(t)
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte(loosen39Fixture), 0o600); err != nil {
		t.Fatalf("cannot seed %s: %v", path, err)
	}
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	rt := &agentRuntime{spec: runSpec{dataDir: dir}, stderr: os.Stderr, mgr: mgr}
	_ = rt // no gate, and startConfigLoosenGate is never called: the queue stays nil
	_, err = mgr.SetPluginEnabled("zed", true, config.SourceNative, "cli.test")
	if err == nil {
		t.Fatal("a loosening with no queue wired must be refused")
	}
	if readZedEnabled(t, path) {
		t.Fatal("the unwired refusal still wrote the loosening")
	}
}
