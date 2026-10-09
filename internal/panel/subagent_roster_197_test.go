package panel

// Ticket 197's 载体层, measured at the pump (leg C cell (a) and the wire half of
// cell (b)).
//
// Legs A and B built a subagent and its stream; what they left unconnected is
// named in 197-r2's own §6.2: "截断在界面上不可见：没有任何用例断言面板显示
// 已截断/谁被丢了 … 载体那一层（197-r3）落地之前，一行被截断的文字在页面上和
// 全文长得一样". So every assertion below reads the BYTES the pump publishes, not
// a Go struct: a field that exists on a view type and never reaches the wire is
// the shape ledger A408 calls 文案在、控件不在, and it is the exact failure this
// file is here to make red.
//
// What stays out of this file on purpose:
//   - the production assembly. cmd/wisp owns the "a real run spawned a real
//     subagent and its packet carries it" judgment
//     (TestRunPacketCarriesTheSubagentItsRosterRowFed), because only a real
//     assembly has a roster worth reading. A package-local roster proves the
//     mapping, not the wiring.
//   - the four-key byte nails (pump_test.go:123, :291). Untouched, and case 2
//     below is the reason they still pass: the new key is a pointer with
//     omitempty, so a pump with no roster reader sends exactly four keys.

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// rosterWire is the packet's tasks section as the renderer reads it: no Go-side
// help, and every key the page has to declare in frontend/src/lib/panel.ts.
type rosterWire struct {
	Rows              []rosterRowWire `json:"rows"`
	InFlightSlots     int             `json:"inFlightSlots"`
	PoolCap           int             `json:"poolCap"`
	StreamTruncated   bool            `json:"streamTruncated"`
	StreamElidedRunes int             `json:"streamElidedRunes"`
	DroppedStreamKeys []string        `json:"droppedStreamKeys"`
}

// rosterRowWire is one row of the section as the page reads it. It is a named
// type because the assertion helper below hands rows back, and an anonymous
// struct there is a thing nobody can read.
type rosterRowWire struct {
	TaskID            string `json:"taskId"`
	Label             string `json:"label"`
	Kind              string `json:"kind"`
	ParentTaskID      string `json:"parentTaskId"`
	Status            string `json:"status"`
	StatusKnown       bool   `json:"statusKnown"`
	StatusReason      string `json:"statusReason"`
	StreamKey         string `json:"streamKey"`
	BlockedOnApproval bool   `json:"blockedOnApproval"`
	StreamTruncated   bool   `json:"streamTruncated"`
	StreamElidedRunes int    `json:"streamElidedRunes"`
	StreamDropped     bool   `json:"streamDropped"`
}

// packetRoster reads the tasks section back out of published bytes. Its absence
// is a Fatal: "the carrier exists but nothing sent it" is the finding this leg
// was written to close, so it must not be able to pass by being skipped.
func packetRoster(t *testing.T, data []byte) rosterWire {
	t.Helper()
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("the packet is not JSON: %v", err)
	}
	raw, ok := wire["tasks"]
	if !ok {
		t.Fatalf("the packet carries no tasks key at all: %s", data)
	}
	var sect rosterWire
	if err := json.Unmarshal(raw, &sect); err != nil {
		t.Fatalf("tasks section unreadable: %v (%s)", err, raw)
	}
	return sect
}

func wireRow(t *testing.T, sect rosterWire, taskID string) rosterRowWire {
	t.Helper()
	for _, r := range sect.Rows {
		if r.TaskID == taskID {
			return r
		}
	}
	t.Fatalf("the packet carries no row for %q: %+v", taskID, sect.Rows)
	panic("unreachable")
}

// okExit197 is the pump's far end for these cases. Publish refuses to report
// success when there is no exit at all (pump.go:271-274 - the honest version of
// "the last mile does not exist yet"), and the bytes asserted below are exactly
// the ones it hands that exit, i.e. what the transport would carry.
func okExit197(Snapshot, []byte) error { return nil }

// TestTheRosterReaderPutsSubagentsOnTheWire is cell (a) at the pump: the user can
// list which subagents exist, which D43 state each is in, and who derived it - and
// each row names the key its own work page is filed under.
func TestTheRosterReaderPutsSubagentsOnTheWire(t *testing.T) {
	pump := NewSnapshotPump(PumpSources{
		Out: okExit197,
		Verdicts: func() []NativeVerdict {
			return []NativeVerdict{{
				CorrelationID: "child-1", Tool: "fs.write", Level: risk.L2,
				RulesHit: []risk.RuleID{risk.R2}, Reason: "写在工作区之外",
			}}
		},
		Mode: func() risk.Mode { return risk.ModeAskHighRisk },
		Now:  func() time.Time { return time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC) },
		Tasks: func() TaskRosterState {
			return TaskRosterState{
				Rows: []TaskRow{
					{
						TaskID: "root-1", Label: "总结这份笔记", Kind: "root",
						Status: "Thinking", StatusKnown: true, StreamKey: "root-1",
					},
					{
						TaskID: "child-1", ParentTaskID: "root-1", Label: "读三份文件",
						Kind: "subagent", Status: string(statemachine.StateThinking), StatusKnown: true,
						StreamKey: SubagentStreamKey("child-1"),
					},
				},
				InFlightSlots: 1, PoolCap: poolCapGiven197,
			}
		},
	})
	_, data, err := pump.Publish()
	if err != nil {
		t.Fatal(err)
	}
	sect := packetRoster(t, data)
	if len(sect.Rows) != 2 {
		t.Fatalf("rows = %d, want the two the roster held: %+v", len(sect.Rows), sect.Rows)
	}
	if sect.Rows[0].TaskID > sect.Rows[1].TaskID {
		t.Errorf("rows arrived unordered (%s before %s): the ledger books a sha256 of these bytes, "+
			"so an unordered section is a packet that cannot be re-derived",
			sect.Rows[0].TaskID, sect.Rows[1].TaskID)
	}

	child := wireRow(t, sect, "child-1")
	if child.ParentTaskID != "root-1" {
		t.Errorf("parentTaskId = %q, want the task that derived it", child.ParentTaskID)
	}
	if child.Label != "读三份文件" || child.Kind != "subagent" {
		t.Errorf("label/kind = %q/%q, want what the host filed verbatim", child.Label, child.Kind)
	}
	if !child.StatusKnown || child.Status != string(statemachine.StateThinking) {
		t.Errorf("status = %q known=%v, want Thinking: the roster's own D43 name has to reach the page",
			child.Status, child.StatusKnown)
	}
	// The field that answers "这一行是谁的" without touching ResultChunk's three
	// keys: the page takes streamKey and finds this agent's row in the same packet.
	if child.StreamKey != "subagent:child-1" {
		t.Errorf("streamKey = %q, want subagent:child-1 (the row's own work page address)",
			child.StreamKey)
	}

	// AC#4 / PLAN.md:1127 - a subagent blocked on a card is visible as blocked.
	if !child.BlockedOnApproval {
		t.Error("the child has an L2 card on the queue and the packet does not say so: " +
			"「被阻塞的子代理在所有界面上都不可见」 is somebody else's recorded incident")
	}
	root := wireRow(t, sect, "root-1")
	if root.BlockedOnApproval {
		t.Error("the root was marked blocked although no card names it - the join is by correlation id, " +
			"not by anything being in flight")
	}
	if sect.InFlightSlots != 1 {
		t.Errorf("inFlightSlots = %d, want the slot the pool is holding", sect.InFlightSlots)
	}
	if sect.PoolCap != poolCapGiven197 {
		t.Errorf("poolCap = %d, want the cap the host handed over (%d) carried verbatim - this "+
			"package neither reads nor re-decides that number", sect.PoolCap, poolCapGiven197)
	}
}

// poolCapGiven197 is the number THIS file hands the pump. The point of carrying it
// through a constant is that the assertion below is about the carrier passing the
// host's value unchanged, not about what the real pool cap is: that number is
// tools.MaxConcurrentSubagents, its tie to the bridge's D38d ceiling is ticket
// 211's, and cmd/wisp's carrier case is where the real one gets measured.
const poolCapGiven197 = 4

// TestAPumpWithoutARosterReaderSendsFourKeysIsStillTheContractCase is the half the
// two byte nails cannot cover: those nails read a fixture that never fills an
// omitempty section, so they pass whether or not a new key exists. This case pins
// the OTHER direction - a pump that has the reader always sends it, and one that
// has it not sends exactly the four keys panel.ts declares today.
func TestAPumpWithoutARosterReaderSendsFourKeys(t *testing.T) {
	pump := NewSnapshotPump(PumpSources{
		Mode: func() risk.Mode { return risk.ModeAskEveryStep },
		Now:  func() time.Time { return time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC) },
	})
	data, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	if _, ok := wire["tasks"]; ok {
		t.Fatalf("a pump assembled without a roster reader sent the key anyway: %s", data)
	}
	keys := make([]string, 0, len(wire))
	for k := range wire {
		keys = append(keys, k)
	}
	if got := strings.Join(sortedCopy(keys), ","); got != "composer,generatedAt,pending,results" {
		t.Errorf("keys without the reader = %q, want the four", got)
	}

	// And the paired direction: with a reader, the key arrives even when the run
	// has spawned nothing - "no subagents yet" and "this host cannot see
	// subagents" must stay two readings (the reason 200-r2 replaced a bare list).
	empty := NewSnapshotPump(PumpSources{
		Out:   okExit197,
		Tasks: func() TaskRosterState { return TaskRosterState{} },
		Now:   func() time.Time { return time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC) },
	})
	_, emptyData, err := empty.Publish()
	if err != nil {
		t.Fatal(err)
	}
	sect := packetRoster(t, emptyData)
	if len(sect.Rows) != 0 {
		t.Errorf("rows = %d, want the empty array a run with no children has", len(sect.Rows))
	}
	if sect.DroppedStreamKeys == nil {
		t.Error("droppedStreamKeys marshalled as null: a renderer should be able to read the " +
			"list without a nil check, and 'nothing was dropped' has to be a stated fact")
	}
}

// TestAStatusOutsideD43NeverReachesTheWire is cell (c): the roster's status
// dimension carries D43's names and nothing else. statemachine.Valid is the judge
// - and it is the judge in the HOST (tools.TaskOutput.StateAnswer), which
// TestTaskState188AC3PanelHasNoWriteLeg keeps from being otherwise, so this case
// pins what the carrier does with a reading that was refused: the name is dropped
// and the reason travels, never softened.
func TestAStatusOutsideD43NeverReachesTheWire(t *testing.T) {
	// The 20 names themselves must all survive the carrier untouched.
	all := make([]TaskRow, 0, 20)
	for _, s := range []statemachine.State{
		statemachine.StateFirstRun, statemachine.StateSleeping, statemachine.StateArmed,
		statemachine.StateMuted, statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateActing, statemachine.StateSpeaking, statemachine.StateWarm,
		statemachine.StateConversation, statemachine.StateConfirming,
		statemachine.StateAwaitingApproval, statemachine.StateSettling,
		statemachine.StateDownloading, statemachine.StateUnconfigured,
		statemachine.StateNoNetwork, statemachine.StateError, statemachine.StateQueued,
		statemachine.StateStuck, statemachine.StateWatchdogAlert,
	} {
		if !statemachine.Valid(s) {
			t.Fatalf("fixture: %q is not D43-valid", s)
		}
		all = append(all, TaskRow{
			TaskID: string(s), Status: string(s), StatusKnown: statemachine.Valid(s),
			StreamKey: string(s), Kind: "subagent",
		})
	}
	pump := NewSnapshotPump(PumpSources{Out: okExit197, Tasks: func() TaskRosterState {
		return TaskRosterState{Rows: all}
	}})
	_, data, err := pump.Publish()
	if err != nil {
		t.Fatal(err)
	}
	sect := packetRoster(t, data)
	if len(sect.Rows) != 20 {
		t.Fatalf("rows = %d, want all 20 D43 names carried", len(sect.Rows))
	}
	for _, r := range sect.Rows {
		if !r.StatusKnown || !statemachine.Valid(statemachine.State(r.Status)) {
			t.Errorf("row %q carries status %q, which is not one of D43's 20 names", r.TaskID, r.Status)
		}
	}

	// The refused reading: a legacy word (ticket 196's vocabulary) handed over as
	// if it were a status must not reach the page as a status.
	legacy := NewSnapshotPump(PumpSources{Out: okExit197, Tasks: func() TaskRosterState {
		return TaskRosterState{Rows: []TaskRow{{
			TaskID: "x-1", Kind: "subagent", Status: "running", StatusKnown: false,
			StatusReason: "登记的「running」不是 D43 状态机表里的 20 个名字之一",
		}}}
	}})
	_, legacyData, err := legacy.Publish()
	if err != nil {
		t.Fatal(err)
	}
	row := wireRow(t, packetRoster(t, legacyData), "x-1")
	if row.Status != "" {
		t.Errorf("status = %q: a name the host's judge refused still reached the wire, "+
			"which is the second vocabulary ticket 197 §2(c) forbids", row.Status)
	}
	if row.StatusKnown {
		t.Error("statusKnown = true for a refused name")
	}
	if !strings.Contains(row.StatusReason, "D43") {
		t.Errorf("statusReason = %q, want the host's own sentence for why nothing was filed", row.StatusReason)
	}

	// And the third state, which is neither: a host that filed nothing at all.
	blank := NewSnapshotPump(PumpSources{Out: okExit197, Tasks: func() TaskRosterState {
		return TaskRosterState{Rows: []TaskRow{{TaskID: "x-2", Kind: "root"}}}
	}})
	_, blankData, err := blank.Publish()
	if err != nil {
		t.Fatal(err)
	}
	blankRow := wireRow(t, packetRoster(t, blankData), "x-2")
	if blankRow.Status != "" || blankRow.StatusKnown {
		t.Errorf("an unread dimension read as %q/known=%v: ticket 92 AC#4's rule, inherited",
			blankRow.Status, blankRow.StatusKnown)
	}
	if blankRow.StatusReason == "" {
		t.Error("a blank status with no reason is the shape StateAnswer exists to prevent")
	}
}

// TestTruncationFactsRideThePacket is the wire half of cell (b): when the shared
// stream log had to cut a subagent's own text, the packet says HOW MUCH that one
// stream lost, and it says it on the row it applies to. Before this leg the only
// readers of those facts were Go-side (197-r2 §6.2: zero production callers), so
// a truncated line and a full line reached the page looking identical.
func TestTruncationFactsRideThePacket(t *testing.T) {
	log := NewStreamLog(2)
	long := strings.Repeat("漢", 600) // runes, not bytes: the bound is counted in runes
	log.Append("subagent:big", long)
	log.Append("subagent:small", "一句短话")
	// A third key passes the bound, which is what puts the WHOLE log into
	// truncating mode - and the short row must still lose nothing.
	log.Append("subagent:third", "第三条")

	if !log.Truncated() {
		t.Fatal("fixture did not reach the bound, so there is no truncation to carry")
	}
	big := log.TruncationFor("subagent:big")
	if !big.Truncated || big.ElidedRunes != 600-2*streamTruncateKeepRunes {
		t.Errorf("big row: truncated=%v elided=%d, want true/%d",
			big.Truncated, big.ElidedRunes, 600-2*streamTruncateKeepRunes)
	}
	small := log.TruncationFor("subagent:small")
	if small.Truncated || small.ElidedRunes != 0 {
		t.Errorf("the short row paid for somebody else's overflow: truncated=%v elided=%d, want false/0",
			small.Truncated, small.ElidedRunes)
	}
	// The marker in the text and the number on the row are the same fact told
	// twice by two different halves of the packet; they must not diverge.
	var rowText string
	for _, c := range log.Chunks() {
		if c.CorrelationID == "subagent:big" {
			rowText = c.Text
		}
	}
	if !strings.Contains(rowText, "[truncated: ") {
		t.Fatalf("the big row carries no inline marker, so the numbers above describe a fold that "+
			"never happened: %q", rowText)
	}

	pump := NewSnapshotPump(PumpSources{
		Out:     okExit197,
		Results: log.Chunks,
		Tasks: func() TaskRosterState {
			return TaskRosterState{
				Rows: []TaskRow{
					{
						TaskID: "big", Kind: "subagent", Status: "Thinking", StatusKnown: true,
						StreamKey:       "subagent:big",
						StreamTruncated: true, StreamElidedRunes: big.ElidedRunes,
					},
					{
						TaskID: "small", Kind: "subagent", Status: "Thinking", StatusKnown: true,
						StreamKey:       "subagent:small",
						StreamTruncated: small.Truncated, StreamElidedRunes: small.ElidedRunes,
					},
				},
				StreamTruncated:   log.Truncated(),
				StreamElidedRunes: log.ElidedRunes(),
			}
		},
	})
	_, data, err := pump.Publish()
	if err != nil {
		t.Fatal(err)
	}
	sect := packetRoster(t, data)
	if !sect.StreamTruncated {
		t.Error("tasks.streamTruncated = false although the log said it is truncating")
	}
	if sect.StreamElidedRunes != log.ElidedRunes() {
		t.Errorf("section total = %d, want the log's own %d", sect.StreamElidedRunes, log.ElidedRunes())
	}
	if got := wireRow(t, sect, "big"); got.StreamElidedRunes != 600-2*streamTruncateKeepRunes || !got.StreamTruncated {
		t.Errorf("big row on the wire = elided %d/truncated %v, want %d/true",
			got.StreamElidedRunes, got.StreamTruncated, 600-2*streamTruncateKeepRunes)
	}
	if got := wireRow(t, sect, "small"); got.StreamTruncated || got.StreamElidedRunes != 0 {
		t.Errorf("small row on the wire = truncated %v/elided %d, want false/0: the section must not "+
			"make every row look damaged because one row was", got.StreamTruncated, got.StreamElidedRunes)
	}
	if !strings.Contains(string(data), `"streamElidedRunes":88`) {
		t.Errorf("the bytes do not state the loss as a number: %s", data)
	}
}

// TestDroppedStreamsAreNamedOnTheWire is the other half of "不许静默丢": past the
// hard ceiling a stream stops being tracked at all, and the packet must name which
// ones - "nobody is showing this one" is a different sentence from "somebody
// merged it into a row that is still here" (197-r2 §1 乙).
func TestDroppedStreamsAreNamedOnTheWire(t *testing.T) {
	log := NewStreamLog(2) // hard ceiling = 2 * StreamKeyHardCeilingMultiple = 4 keys
	for i := 0; i < 6; i++ {
		log.Append(string(rune('a'+i))+strings.Repeat("x", 600), "delta")
	}
	dropped := log.DroppedKeys()
	if len(dropped) != 2 {
		t.Fatalf("dropped keys = %v, want the two oldest past the ceiling", dropped)
	}
	for _, k := range dropped {
		if tr := log.TruncationFor(k); !tr.Dropped {
			t.Errorf("TruncationFor(%q).Dropped = false although DroppedKeys names it: the per-row "+
				"reader and the log-level reader disagree about the same fact", k)
		}
	}
	pump := NewSnapshotPump(PumpSources{Out: okExit197, Tasks: func() TaskRosterState {
		return TaskRosterState{
			DroppedStreamKeys: log.DroppedKeys(),
			StreamTruncated:   log.Truncated(),
			StreamElidedRunes: log.ElidedRunes(),
		}
	}})
	_, data, err := pump.Publish()
	if err != nil {
		t.Fatal(err)
	}
	sect := packetRoster(t, data)
	if len(sect.DroppedStreamKeys) != len(dropped) {
		t.Errorf("packet names %v, want the log's own %v", sect.DroppedStreamKeys, dropped)
	}
	if sect.InFlightSlots != 0 || sect.PoolCap != 0 {
		t.Errorf("a state the host never read marshalled as %d/%d instead of 0/0",
			sect.InFlightSlots, sect.PoolCap)
	}
}

// WHAT THIS FILE DID NOT MEASURE.
//
//  1. Nothing here is the production assembly. Every pump above was built by this
//     file with hand-shaped TaskRow values, which proves the mapping and the
//     encoder and NOT that a running `wisp run` fills them - that reading belongs
//     to cmd/wisp (TestRunPacketCarriesTheSubagentItsRosterRowFed). A regression
//     that only breaks the wiring in run.go passes every case above.
//  2. The status judge itself is not tested here. This file passes StatusKnown in;
//     the thing that decides it (tools.TaskOutput.StateAnswer via
//     statemachine.Valid) is judged in internal/tools, and the reason the carrier
//     may not run that check is a resident nail of its own
//     (TestTaskState188AC3PanelHasNoWriteLeg). If the host ever stopped calling
//     StateAnswer and started passing StatusKnown: true for anything, no case in
//     this file would notice.
//  3. The blockedOnApproval join is asserted against a card this file made. On
//     the production path it pairs the row's task id, not the corr: since
//     bd124b2a callCorr mints one corr per call (loop.go:603, used :676),
//     shape <taskID>#<callID>, never == taskID. A host that paired on the
//     corr alone would miss every row; only cmd/wisp's real gate catches it.
//  4. No byte-size bound: nothing marshals a saturated roster and measures the
//     packet against the ledger's 512-rune log bound. A run with 4 children each
//     at the retained floor is still a packet the exit can carry and no ruler here
//     says so.
//  5. Reasoning deltas and a child's tool-call lines are deliberately NOT on its
//     page, and no case asserts they are absent - the parity argument is in
//     internal/tools' subagentTextSink's comment, not in a judgment.
