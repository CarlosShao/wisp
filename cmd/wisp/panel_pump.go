package main

// The panel snapshot pump, wired (ticket 35's data half, the assembly-root leg).
//
// THIS IS WHERE THE PACKET STARTS EXISTING. Before these files, panel.Snapshot
// had zero non-test constructors and the panel had never received a packet - not
// because the view model lacked fields (ticket 145's census measured 12 of 14
// state rows as already sourced) but because nothing built it. composer.go:41-43
// named the disease this closes: "a view model that only a test can build is a
// view model production never sends".
//
// WHAT IT IS NOT: the transport. There is no Go -> page channel in this tree -
// no WebView2 host (ticket 33 is unclaimed), no postMessage writer, no local
// HTTP/SSE server, and D29 through assets.go:8 forbids the last one. So the
// bytes this assembles are booked into the run's persistent ledger, which is the
// one in-tree production destination for a record this process emits about
// itself (internal/agent's audit lines, the winsec notices, ticket 117's whole
// reason for existing). Re-pointing that one line at ticket 33's host is the
// last mile, and it is still open. Nothing here pretends otherwise, and nothing
// here invents a substitute channel: a file the renderer would poll is exactly
// the second channel composer.go:14 refuses.
//
// EVERY VALUE COMES FROM AN OBJECT THIS PROCESS IS RUNNING. pending is read off
// the live approval queue (approval.Queue.LiveApprovals), which carries the
// verdict internal/risk already reached; composer.mode is perm.Store's own read;
// composer.workspace is the C26 canonicalizer's current root; results is the
// stream this run is writing. The one field with no producer is named in
// liveVerdicts, and no constant stands in for it.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/tools"
)

// liveVerdicts reads the approvals that are pending right now.
//
// The mapping is a copy, not a judgement: Level/RulesHit/Reason/
// SessionOverrideBlocked are the fields internal/risk put on the queue item
// when it was admitted, and panel renders them through the same
// CardViewFromDecision the CLI diagnostic uses, so a card on the panel and the
// card `wisp panel-assets -l2` prints cannot disagree about what L2 meant.
//
// CallChain stays empty. ApprovalCardView's header (panel/approval.go:9-11)
// says the card owes the user "the chain that got there", and no code in this
// tree records one for a live call - tools.Decision has TaskID/CallID/Provider
// and no chain field, and the only caller that ever filled CallChain is the CLI
// probe, which names its own argv door. Filling it with a plausible
// ["session","agent-loop",tool] would put a sentence no producer wrote onto a
// security card, so the packet reports the gap instead and the gap is registered
// in this ticket's evidence file.
func (rt *agentRuntime) liveVerdicts() []panel.NativeVerdict {
	if rt == nil || rt.gate == nil {
		return nil
	}
	items := rt.gate.Queue().LiveApprovals()
	out := make([]panel.NativeVerdict, 0, len(items))
	for _, it := range items {
		out = append(out, panel.NativeVerdict{
			CorrelationID:          it.CorrelationID,
			Tool:                   it.Decision.Tool,
			Args:                   panelArgs(it.Decision),
			Level:                  it.Decision.Level,
			RulesHit:               it.Decision.RulesHit,
			Reason:                 it.Decision.Reason,
			SessionOverrideBlocked: it.Decision.SessionOverrideBlocked,
		})
	}
	return out
}

// workspaceView reports the narrowing native resolved, through panel's own
// renderer of the two states (set / unset-with-a-reason). The spelling and the
// reparse/rewrite account live on the request path (RequestWorkspaceSwitch),
// which the tree still has no caller for, so the snapshot carries what the
// canonicalizer holds: the canonical root, or an honest "no workspace".
func (rt *agentRuntime) workspaceView() panel.WorkspaceView {
	if rt == nil || rt.paths == nil {
		return panel.UnsetWorkspaceView()
	}
	return panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())
}

// panelArgs shapes one live call's argument bytes into the argv-shaped list
// ApprovalCardView.args shows. Two real shapes exist: a call that carries an
// argv vector (what risk.Facts.ShellArgv is built from, and what
// `wisp panel-assets -l2` already prints), and every other call, whose argument
// object is carried verbatim. The second branch is not a summary - it is the
// bytes, which is the standard the card is held to ("full params, not a
// summary", approval.Prompt's comment).
func panelArgs(d tools.Decision) []string {
	raw := strings.TrimSpace(string(d.Args))
	if raw == "" {
		return nil
	}
	var params map[string]any
	if err := json.Unmarshal([]byte(raw), &params); err == nil {
		if argv, ok := stringList(params["argv"]); ok && len(argv) > 0 {
			return argv
		}
		if cmd, ok := params["command"].(string); ok && strings.TrimSpace(cmd) != "" {
			return []string{cmd}
		}
	}
	return []string{raw}
}

// stringList reads a JSON array of strings back from a decoded map.
func stringList(v any) ([]string, bool) {
	list, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(list))
	for _, e := range list {
		s, ok := e.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

// bookPanelSnapshot is the exit's far end, and it is where this repository's
// missing transport becomes a measured fact rather than an assumption.
//
// WHAT IT CAN AND CANNOT CARRY. The one durable, machine-readable destination a
// `wisp run` has is its own ledger (cmd/wisp/logsink.go), and every string that
// reaches it goes through internal/observe's redacting handler, which bounds a
// single string at observe.MaxLoggedString = 512 characters by rule 4 of its own
// header ("long argument strings are truncated (bounded log lines)"). A panel
// packet is routinely longer than that. Measured on this wiring, both units
// written down (票 35 fix r1, docs/evidence/s1/35-panel-snapshot-pump-fix-r1.md
// §1.5): the smallest packet this shape produces is 552 bytes / 518 runes, and
// one carrying a single L2 card is 792 bytes / 734 runes (793/759 with an ASCII
// reason). The bound is counted in runes by the handler, so the smallest case
// clears it by 6 - thin, and thin on purpose: it is the empty-pending,
// empty-results, unset-workspace packet, which no `wisp run` books. The pair
// this comment used to carry (534 / 997) is not what any measurement on record
// produced, which is why it is now stated with its units and its source.
// So the ledger physically cannot
// carry the packet. That is not a reason to raise a global log bound for one
// ticket's convenience, and it is definitely not a reason to start writing
// files for a renderer to poll.
//
// So the booking is what fits: a bounded summary line whose fields are the facts
// an operator needs to know something was sent (when, how deep, which
// correlation ids, which mode, whether scope was narrowed, how many result
// chunks) plus the packet's byte length and sha256. The digest makes any later
// claim about "what the panel would have shown" checkable against the retained
// packet below, and a mismatch is visible instead of silent. The full bytes stay
// with the process until the transport that delivers them exists.
func (rt *agentRuntime) bookPanelSnapshot(snap panel.Snapshot, data []byte) error {
	rt.snapMu.Lock()
	rt.lastSnap, rt.lastSnapBytes, rt.snapSeen = snap, data, true
	rt.snapMu.Unlock()
	rt.spec.sink.logger().Info(panelSnapshotSummary(snap, data))
	return nil
}

// panelSnapshotSummary renders the bounded ledger record. It clamps its own
// length rather than relying on the redactor to cut it mid-sentence: a record
// that got silently shortened by someone else's rule is a record that lies about
// what it holds.
func panelSnapshotSummary(snap panel.Snapshot, data []byte) string {
	ids := make([]string, 0, len(snap.Pending))
	for i, c := range snap.Pending {
		if i >= 3 {
			ids = append(ids, fmt.Sprintf("+%d", len(snap.Pending)-3))
			break
		}
		ids = append(ids, c.CorrelationID)
	}
	workspace := "unset"
	if snap.Composer.Workspace.Set {
		workspace = "set"
	}
	done := 0
	for _, r := range snap.Results {
		if r.Done {
			done++
		}
	}
	line := fmt.Sprintf("panel: SNAPSHOT at=%s depth=%d pending=%s mode=%s ws=%s results=%d/%d done bytes=%d sha256=%s",
		snap.GeneratedAt, len(snap.Pending), strings.Join(ids, ","),
		snap.Composer.Mode.Current, workspace,
		len(snap.Results), done, len(data), sha256Short(data))
	if len(line) > summaryClamp {
		// Drop the ids before anything else: the count and the digest still pin
		// which packet this was, and the ids are in the packet.
		line = fmt.Sprintf("panel: SNAPSHOT at=%s depth=%d pending=- mode=%s ws=%s results=%d/%d done bytes=%d sha256=%s",
			snap.GeneratedAt, len(snap.Pending), snap.Composer.Mode.Current, workspace,
			len(snap.Results), done, len(data), sha256Short(data))
	}
	return line
}

// summaryClamp is how much of observe.MaxLoggedString this line may use. The
// margin is the redaction layer's own rewriting (a matched secret pattern widens
// the text it replaces), not a guess at formatting.
const summaryClamp = 440

// sha256Short is the packet's fingerprint, long enough to collide on purpose
// being improbable and short enough to read off a log line.
func sha256Short(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])[:16]
}

// lastPanelSnapshot returns the most recent packet this run built and the bytes
// it encoded to. The shape is deliberate and it copies cmd/wisp's existing one:
// consoleApprovalUI.cards plus shown() plus agentRuntime.windowCount() - a
// production path records what it did, and a reader asks it afterwards. It is
// not a seam a test uses to feed values in: everything in the returned packet was
// produced by rt.pump.Publish(), which only the run's own approval surface and
// the loop's own sink call.
func (rt *agentRuntime) lastPanelSnapshot() (panel.Snapshot, []byte, bool) {
	if rt == nil {
		return panel.Snapshot{}, nil, false
	}
	rt.snapMu.Lock()
	defer rt.snapMu.Unlock()
	return rt.lastSnap, rt.lastSnapBytes, rt.snapSeen
}

// publishPanelSnapshot puts one snapshot on the wire. Called where a state the
// panel shows actually moved: a card opened, a card went away, a tool call
// started or ended, the stream closed.
//
// Not on streamed text deltas. Every delta would book a record into the ledger
// for one more word of a reply, and the merge/backpressure rule that decides how
// often a panel is pushed belongs to the transport (D38d, ticket 35 AC#5), which
// does not exist yet. The deltas still reach the packet, because the stream log
// accumulates them and the next publish carries what has streamed so far.
func (rt *agentRuntime) publishPanelSnapshot() {
	if rt == nil || rt.pump == nil {
		return
	}
	snap, _, err := rt.pump.Publish()
	if err != nil {
		// The packet was built and did not get out. That is never silent.
		fmt.Fprintf(rt.stderr, "wisp run: 面板快照未送达（待审批 %d 条）：%v\n", len(snap.Pending), err)
	}
}

// snapshotCount reports how many packets this run put through the exit, so a
// test (and an operator reading a log) can tell "the pump was assembled" from
// "the pump was driven".
func (rt *agentRuntime) snapshotCount() int {
	if rt == nil || rt.pump == nil {
		return 0
	}
	return rt.pump.Publishes()
}
