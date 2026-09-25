package main

// Ticket 35's data half, measured from the outside: does a real `wisp run` build
// the packet the panel renders, out of the state it is actually running?
//
// The judgement these cases exist to answer is ticket 145 AC#2's: a view model
// that only a test can build is a view model production never sends
// (panel/composer.go:41-43). So nothing here constructs a panel.Snapshot, a
// panel.ApprovalCardView or a panel.ComposerState. Every packet read below came
// out of rt.pump.Publish(), and the only calls to the pump are the ones the run
// itself makes (consoleApprovalUI.Prompt / .Update, consoleSink.Publish).
//
// Each value is then compared against a DIFFERENT production surface of the same
// run, never against a literal this file typed:
//
//	pending's level / tool / reason / rules  vs  the console card the gate printed
//	                                            from approval.Prompt
//	composer.mode.current                    vs  the "perm: MODE-READ" audit line
//	                                            the assembly wrote at boot
//	results[].correlationId                  vs  the task id in the run's status line
//	results[].text                           vs  the bytes the run streamed to stdout
//	composer.workspace.canonical             vs  the view panel.RequestWorkspaceSwitch
//	                                            resolved through C26
//
// The ledger record is a bounded summary, not the packet: internal/observe's
// redacting handler caps any single logged string at 512 characters (rule 4 of
// its header), and the smallest packet this ticket produced was 534 bytes. The
// summary therefore carries the packet's byte length and sha256, and the cases
// below check those against the retained bytes - so "the durable record and the
// packet I am asserting on are the same event" is itself asserted, not assumed.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ledgerSummaries145 returns the bounded panel records this run booked, oldest
// first. An empty read is a failure, not a pass: "assembled but never driven" is
// exactly the shape this ticket is named for.
func ledgerSummaries145(t *testing.T, dataDir string) []map[string]string {
	t.Helper()
	lines := sinkLines117(t, filepath.Join(dataDir, "logs"))
	const tag = "panel: SNAPSHOT "
	var out []map[string]string
	for _, line := range lines {
		var rec struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		if !strings.HasPrefix(rec.Msg, tag) {
			continue
		}
		fields := map[string]string{}
		for _, f := range strings.Fields(strings.TrimPrefix(rec.Msg, tag)) {
			if k, v, ok := strings.Cut(f, "="); ok {
				fields[k] = v
			}
		}
		out = append(out, fields)
	}
	if len(out) == 0 {
		t.Fatalf("the ledger holds no panel record in %d lines, so nothing was published", len(lines))
	}
	return out
}

// assertPacketMatchesLedger145 pins that the packet a case is asserting on is the
// one the durable record describes.
func assertPacketMatchesLedger145(t *testing.T, rec map[string]string, data []byte, snap panel.Snapshot) {
	t.Helper()
	if got, want := fmt.Sprint(len(data)), rec["bytes"]; got != want {
		t.Errorf("ledger says %s bytes, retained packet is %d", want, len(data))
	}
	sum := sha256.Sum256(data)
	if rec["sha256"] != hex.EncodeToString(sum[:])[:16] {
		t.Errorf("ledger digest %q is not the retained packet's %q",
			rec["sha256"], hex.EncodeToString(sum[:])[:16])
	}
	if rec["depth"] != fmt.Sprint(len(snap.Pending)) {
		t.Errorf("ledger depth %q, retained packet holds %d cards", rec["depth"], len(snap.Pending))
	}
	if rec["mode"] != snap.Composer.Mode.Current {
		t.Errorf("ledger mode %q, retained packet says %q", rec["mode"], snap.Composer.Mode.Current)
	}
}

var modeRead145 = regexp.MustCompile(`perm: MODE-READ origin=startup mode=(\S+) source=`)

// executeOn145 runs one host-side tool call on the bridge this root assembled,
// on its own goroutine, and waits until the gate has displayed a card. The card
// is what makes the queue non-empty, which is the whole subject of these cases.
func executeOn145(t *testing.T, rt *agentRuntime, taskID, corr, callID, tool, args string) (agent.ToolOutcome, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return rt.bridge.Execute(ctx, agent.ToolRequest{
		TaskID: taskID, CorrelationID: corr, CallID: callID,
		Name: tool, Args: json.RawMessage(args),
	})
}

// nonDefaultConfig145 gives this run two values the fixture's config does not
// carry, and both are there to make an assertion falsifiable rather than
// accidentally true:
//
//   - confirm_timeout_sec = 2, so the case can watch an L2 card open and close
//     (tightening the C18 window is not a loosening);
//   - permission_mode = ask_high_risk, because the default档 is exactly what a
//     pump that hardcoded "ask_every_step" would print. With a non-default档 in
//     the file, the packet's mode is only right if something read it.
func nonDefaultConfig145(t *testing.T, dataDir string) {
	t.Helper()
	cfgPath := filepath.Join(dataDir, configFileName)
	old, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, []byte(string(old)+"\n[risk]\nconfirm_timeout_sec = 2\npermission_mode = \"ask_high_risk\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRunBooksWithASnapshotOfItsLiveQueue(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	nonDefaultConfig145(t, f.dir)
	target := "C:/Windows/win.ini"
	var (
		live      panel.Snapshot
		liveBytes []byte
		driven    *agentRuntime
		cards     int
	)
	f.rtHook = func(rt *agentRuntime) {
		driven = rt
		revoke := rt.gate.AdmitTextTask("pump-task")
		defer revoke()
		done := make(chan struct{})
		var out agent.ToolOutcome
		var err error
		go func() {
			defer close(done)
			// An out-of-allowlist read is R2's own example: the verdict is L2, so
			// the item goes into the approval QUEUE (an L1 window does not), which
			// is what snapshot.pending is made of.
			out, err = executeOn145(t, rt, "pump-task", "pump-corr", "p1", "fs.read",
				fmt.Sprintf(`{"path":%q}`, target))
		}()
		for i := 0; i < 300 && rt.ui.shown() == 0; i++ {
			time.Sleep(10 * time.Millisecond)
		}
		if rt.ui.shown() == 0 {
			t.Fatal("the gate never displayed a card, so the queue never held anything to report")
		}
		// Sample the packet the production readers produce while the card is up.
		// The record this publish books is checked below against the records the
		// run's own triggers booked, so the sampling cannot be the only evidence.
		rt.publishPanelSnapshot()
		live, liveBytes, _ = rt.lastPanelSnapshot()
		<-done
		if err != nil {
			t.Errorf("the L1 route failed: %v", err)
		}
		if !out.IsError {
			t.Errorf("an unanswered L2 card must not execute, got: %s", out.Text)
		}
		cards = rt.windowCount()
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if cards != 1 {
		t.Fatalf("the gate displayed %d cards, want the one L2 card the packet reports", cards)
	}
	if len(live.Pending) != 1 {
		t.Fatalf("sampled pending = %+v, want the one item the queue held", live.Pending)
	}
	card := live.Pending[0]

	// The run's OWN trigger must also have booked a depth-1 record: ui.Prompt
	// publishes, and that publish is what proves the pump is driven by the state
	// change rather than only sampled by a test.
	summaries := ledgerSummaries145(t, f.dir)
	var cardRecords []map[string]string
	for _, rec := range summaries {
		if rec["depth"] == "1" && strings.Contains(rec["pending"], "pump-corr") {
			cardRecords = append(cardRecords, rec)
		}
	}
	// Two of them, and the count is the point: the FIRST was booked by
	// consoleApprovalUI.Prompt, i.e. by the run itself, without this file asking.
	// The second is the sample the field checks below read.
	if len(cardRecords) < 2 {
		t.Fatalf("depth=1 records among %d: %v, want the run's own publish plus this case's sample",
			len(summaries), cardRecords)
	}
	assertPacketMatchesLedger145(t, cardRecords[len(cardRecords)-1], liveBytes, live)

	// Keep the artifact readable from a -v run: this is the first packet
	// production ever built, and the evidence file quotes it verbatim.
	t.Logf("booked record: %v", cardRecords[len(cardRecords)-1])
	t.Logf("packet bytes (%d): %s", len(liveBytes), liveBytes)

	if card.CorrelationID != "pump-corr" {
		t.Errorf("correlationId = %q, want the id the queue issued for this call", card.CorrelationID)
	}
	if card.DecidedBy != "native" {
		t.Errorf("decidedBy = %q, want native", card.DecidedBy)
	}
	if len(card.CallChain) != 0 {
		t.Errorf("callChain = %v, want the honest empty: no production code records a chain", card.CallChain)
	}

	// Cross-check against the other rendering of the same native verdict.
	m := consoleCard145.FindStringSubmatch(f.out.String())
	if m == nil {
		t.Fatalf("the run printed no confirmation card:\n%s", f.out.String())
	}
	if card.Level != m[1] {
		t.Errorf("level: packet=%q console=%q, want the same native level on both surfaces", card.Level, m[1])
	}
	if card.Tool != m[2] {
		t.Errorf("tool: packet=%q console=%q", card.Tool, m[2])
	}
	if got, want := strings.TrimSpace(m[3]), card.Reason; got != want {
		t.Errorf("reason: console=%q packet=%q, want the same native reason on both surfaces", got, want)
	}
	if !card.ReasonKnown {
		t.Error("reasonKnown is false although the gate produced a reason")
	}
	var consoleRules []string
	for _, line := range strings.Split(strings.TrimRight(m[4], "\n"), "\n") {
		if rule := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "规则")); rule != "" {
			consoleRules = append(consoleRules, rule)
		}
	}
	if len(consoleRules) == 0 {
		t.Fatal("the console card named no rule, so the rulesHit check below would prove nothing")
	}
	for _, rule := range consoleRules {
		var found bool
		for _, r := range card.RulesHit {
			if r == rule {
				found = true
			}
		}
		if !found {
			t.Errorf("the card showed rule %s on the console but the packet carries %v", rule, card.RulesHit)
		}
	}
	if len(card.Args) == 0 || !strings.Contains(strings.Join(card.Args, " "), filepath.ToSlash(target)) {
		t.Errorf("args = %v, want the arguments this call carried (target %q)",
			card.Args, filepath.ToSlash(target))
	}
	if live.Composer.Workspace.Set {
		t.Errorf("workspace = %+v, want unset for a run that never switched", live.Composer.Workspace)
	}
	if live.Composer.Workspace.Canonical != "" {
		t.Errorf("canonical = %q, want empty while nothing is narrowed", live.Composer.Workspace.Canonical)
	}
	if live.Composer.Workspace.Reason == "" {
		t.Error("an unset workspace arrived without the reason that makes it readable")
	}
	if len(live.Composer.AttachmentMIMEs) == 0 ||
		live.Composer.MaxAttachmentB != panel.MaxAttachmentBytes {
		t.Errorf("attachment ceiling = %v / %d, want the broker's accepted list and its default",
			live.Composer.AttachmentMIMEs, live.Composer.MaxAttachmentB)
	}
	if live.GeneratedAt == "" {
		t.Error("generatedAt is empty, so the packet cannot say when it was true")
	}

	// The mode section, against the boot's own audit line.
	mm := modeRead145.FindStringSubmatch(f.err.String())
	if mm == nil {
		t.Fatal("no perm: MODE-READ startup line; the mode this packet reports would be unverifiable")
	}
	if live.Composer.Mode.Current != mm[1] {
		t.Errorf("composer.mode.current = %q, want the boot read %q", live.Composer.Mode.Current, mm[1])
	}
	if live.Composer.Mode.Current != risk.ModeAskHighRiskName {
		t.Errorf("composer.mode.current = %q, want the non-default档 this config set",
			live.Composer.Mode.Current)
	}
	if strings.Join(live.Composer.Mode.Names, ",") != "ask_every_step,ask_high_risk,auto_approve" {
		t.Errorf("mode.names = %v, want risk's vocabulary", live.Composer.Mode.Names)
	}

	// The results section: the last record of the run, and the last packet.
	if driven == nil {
		t.Fatal("the hook never saw the assembled runtime")
	}
	finalSnap, finalBytes, sawFinal := driven.lastPanelSnapshot()
	if !sawFinal {
		t.Fatal("no packet was retained, so nothing was published")
	}
	assertPacketMatchesLedger145(t, summaries[len(summaries)-1], finalBytes, finalSnap)
	if len(finalSnap.Pending) != 0 {
		t.Errorf("the closing packet still reports %d pending cards after the answer landed: %+v",
			len(finalSnap.Pending), finalSnap.Pending)
	}
	if len(finalSnap.Results) == 0 {
		t.Fatal("the closing packet has no results section, though this run streamed a reply")
	}
	chunk := finalSnap.Results[0]
	if chunk.CorrelationID != f.taskID() {
		t.Errorf("results correlationId = %q, want the task id the run printed %q",
			chunk.CorrelationID, f.taskID())
	}
	if !chunk.Done {
		t.Error("the reply's chunk is not closed although the loop published EvDone")
	}
	if !strings.Contains(chunk.Text, "echo:") {
		t.Errorf("chunk text %q, want the provider's reply the console also showed", chunk.Text)
	}
	if !strings.Contains(f.out.String(), chunk.Text) {
		t.Errorf("the booked text %q is not the text this run streamed to the console", chunk.Text)
	}
}

// TestSnapshotWorkspaceSectionReportsTheNarrowing drives the real workspace
// handler (panel.RequestWorkspaceSwitch, what a panel request is answered by) and
// then checks what the packet says about scope against what that handler
// resolved through C26.
func TestSnapshotWorkspaceSectionReportsTheNarrowing(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	sub := filepath.Join(f.dir, "scoped")
	note := filepath.Join(sub, "note.txt")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(note, []byte("scoped body"), 0o600); err != nil {
		t.Fatal(err)
	}
	var (
		handled panel.WorkspaceView
		live    panel.Snapshot
	)
	f.rtHook = func(rt *agentRuntime) {
		view, err := panel.RequestWorkspaceSwitch(rt.paths, filepath.ToSlash(sub), rt.auditf)
		if err != nil {
			t.Fatalf("a workspace inside the allowlist must be accepted: %v", err)
		}
		handled = view
		rt.publishPanelSnapshot()
		live, _, _ = rt.lastPanelSnapshot()
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if !handled.Set {
		t.Fatalf("handler view = %+v, want the accepted narrowing", handled)
	}
	if !live.Composer.Workspace.Set {
		t.Fatalf("packet workspace = %+v, want the narrowing this run applied", live.Composer.Workspace)
	}
	// The same tree, and only the same tree. It is compared case/separator-wise
	// ON PURPOSE and the reason is a finding, written up in
	// docs/evidence/s1/35-panel-snapshot-pump-r1.md §3: the value native holds for
	// the narrowed root is internal/tools' fold key (paths_workspace.go's
	// SetWorkspaceRoot stores foldPath(root)), so the packet reports a lowercased
	// path rather than the spelling C26 authorized. Asserting byte equality here
	// would encode that as intended; asserting equality-of-tree without noting it
	// would hide it. Both are worse than naming it.
	if !strings.EqualFold(
		filepath.ToSlash(live.Composer.Workspace.Canonical),
		filepath.ToSlash(handled.Canonical)) {
		t.Errorf("packet canonical = %q, want the tree the handler authorized %q (case-insensitively)",
			live.Composer.Workspace.Canonical, handled.Canonical)
	}
	if live.Composer.Workspace.Canonical == handled.Canonical {
		t.Log("fold-key finding did not reproduce on this run: the two spellings came out identical")
	}
	if live.Composer.Workspace.Reparse || live.Composer.Workspace.Rewritten {
		t.Errorf("workspace = %+v, want a plain in-scope switch reported as such",
			live.Composer.Workspace)
	}
	if live.Composer.Workspace.Reason == "" {
		t.Error("a set workspace arrived without the reason that says what it narrows")
	}
	var sawSet bool
	for _, rec := range ledgerSummaries145(t, f.dir) {
		if rec["ws"] == "set" {
			sawSet = true
		}
	}
	if !sawSet {
		t.Error("no durable panel record reported the narrowed scope")
	}
}

// TestThePumpIsDrivenNotJustAssembled: a runtime that builds the pump and never
// publishes is the exact shape this ticket exists to end, and every field
// assertion above would still pass if someone kept the assembly and dropped the
// calls. The count is read off the pump, which only Publish() moves.
func TestThePumpIsDrivenNotJustAssembled(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	var (
		seenAtHook int
		captured   *agentRuntime
	)
	f.rtHook = func(rt *agentRuntime) {
		seenAtHook = rt.snapshotCount()
		captured = rt
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, f.out.String(), f.err.String())
	}
	if seenAtHook != 0 {
		t.Errorf("snapshotCount at the hook = %d, want 0: nothing had happened yet", seenAtHook)
	}
	if captured == nil {
		t.Fatal("the hook never saw the assembled runtime")
	}
	if got := captured.snapshotCount(); got < 1 {
		t.Errorf("snapshotCount after the run = %d, want at least the closing publish", got)
	}
	if got := len(ledgerSummaries145(t, f.dir)); got < 1 {
		t.Errorf("ledger holds %d panel records, want at least one", got)
	}
}

// consoleCard145 reads the confirmation the gate itself printed: level, tool,
// reason, then the rule lines under it. Another surface of the same verdict is
// what turns the comparisons above into source checks instead of self-assertions.
var consoleCard145 = regexp.MustCompile(
	`(?s)\n\[确认 (\S+) (\S+)\] ([^\n]*)\n((?:  规则 [^\n]*\n)*)`)
