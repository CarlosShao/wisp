package panel

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// The pump's own cases (ticket 35's data half). What each one pins:
//
//	TestThePumpBuildsThePacketFromWhatTheHostHolds  every section of the packet
//	    carries what the reader returned, and the card carries the native verdict
//	    verbatim (rules, reason, R4 flag, decidedBy).
//	TestAPumpWithNoReadersSaysSo                   a missing reader is reported
//	    as unknown / unset / empty, never as a value the host did not have.
//	TestAnUnreadableModeNeverRendersAsASafeOne     AC#4 across the pump, not just
//	    across NewSnapshot: an invalid mode from a live reader is still "unknown".
//	TestTheStreamLogTruncatesInsteadOfMerging      the results section is bounded by
//	    eliding each stream's own middle - ticket 197 leg B's re-cut of the fold.
//	TestPublishWithoutAnExitReports                no exit is a named failure, not
//	    a successful-looking drop.
//	TestPublishHandsTheBytesToTheAttachedExit      the exit gets exactly the bytes
//	    Marshal would have returned.

// fixedClock is a clock, not a deadline: generatedAt is display-only
// (panel.ts:132-134) and nothing here derives state from it.
func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestThePumpBuildsThePacketFromWhatTheHostHolds(t *testing.T) {
	pump := NewSnapshotPump(PumpSources{
		Verdicts: func() []NativeVerdict {
			return []NativeVerdict{{
				CorrelationID:          "corr-1",
				Tool:                   "shell.run",
				Args:                   []string{"git", "status"},
				Level:                  risk.L2,
				RulesHit:               []risk.RuleID{"R6", "R2"},
				Reason:                 "命令串按 SPEC-06 §10 恒为 L2",
				SessionOverrideBlocked: true,
			}}
		},
		Mode: func() risk.Mode { return risk.ModeAskHighRisk },
		Workspace: func() WorkspaceView {
			return WorkspaceView{
				Set: true, Spelling: "D:\\work\\Wisp", Canonical: `D:\work\Wisp`,
				Reason: "已收窄到该工作区：范围外的路径按 R2 判定",
			}
		},
		Results: func() []ResultChunk {
			return []ResultChunk{
				{CorrelationID: "task-1", Text: "前半", Done: false},
				{CorrelationID: "task-1", Text: "后半", Done: true},
			}
		},
		AttachmentMax: 1234,
		Now:           fixedClock(time.Date(2026, 9, 25, 12, 0, 0, 0, time.UTC)),
	})

	snap := pump.Snapshot()

	if len(snap.Pending) != 1 {
		t.Fatalf("pending = %+v, want the one verdict the host reported", snap.Pending)
	}
	card := snap.Pending[0]
	if card.CorrelationID != "corr-1" || card.Tool != "shell.run" {
		t.Errorf("card = %+v, want the queue's own correlation id and tool", card)
	}
	if card.Level != "L2" {
		t.Errorf("level = %q, want the assessed level rendered by risk, not by this file", card.Level)
	}
	if strings.Join(card.RulesHit, ",") != "R6,R2" {
		t.Errorf("rulesHit = %v, want the rules verbatim and in the assessor's order", card.RulesHit)
	}
	if card.Reason != "命令串按 SPEC-06 §10 恒为 L2" || !card.ReasonKnown {
		t.Errorf("reason = %q known=%v, want the native reason carried as-is", card.Reason, card.ReasonKnown)
	}
	if !card.SessionOverrideBlocked {
		t.Error("sessionOverrideBlocked was lost between the verdict and the card (R4's flag)")
	}
	if card.DecidedBy != "native" {
		t.Errorf("decidedBy = %q, want native: the pump cannot make a decision, only carry one", card.DecidedBy)
	}
	if len(card.CallChain) != 0 {
		t.Errorf("callChain = %v, want the honest empty: no code in this tree records a chain", card.CallChain)
	}

	if snap.Composer.Mode.Current != risk.ModeAskHighRiskName {
		t.Errorf("composer.mode.current = %q, want what the store actually holds", snap.Composer.Mode.Current)
	}
	if !snap.Composer.Workspace.Set || snap.Composer.Workspace.Canonical != `D:\work\Wisp` {
		t.Errorf("composer.workspace = %+v, want the resolved workspace", snap.Composer.Workspace)
	}
	if snap.Composer.MaxAttachmentB != 1234 {
		t.Errorf("maxAttachmentBytes = %d, want the ceiling the assembly reported", snap.Composer.MaxAttachmentB)
	}
	if len(snap.Results) != 2 || !snap.Results[1].Done {
		t.Errorf("results = %+v, want both chunks the stream produced", snap.Results)
	}
	if snap.GeneratedAt != "2026-09-25T12:00:00Z" {
		t.Errorf("generatedAt = %q, want the injected clock, RFC3339 UTC", snap.GeneratedAt)
	}

	// And the bytes: the four keys are the contract, so a fifth arriving here is
	// a contract change dressed as a bug fix.
	data, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	got := make([]string, 0, len(generic))
	for k := range generic {
		got = append(got, k)
	}
	if strings.Join(sortedCopy(got), ",") != "composer,generatedAt,pending,results" {
		t.Errorf("snapshot JSON keys = %v, want exactly the four PanelSnapshot declares", got)
	}
}

func TestAPumpWithNoReadersSaysSoInsteadOfInventingState(t *testing.T) {
	snap := NewSnapshotPump(PumpSources{}).Snapshot()

	if snap.Composer.Mode.Current != "unknown" {
		t.Errorf("mode.current = %q, want unknown - no reader means nobody read the档", snap.Composer.Mode.Current)
	}
	if strings.Join(snap.Composer.Mode.Names, ",") != strings.Join(risk.ModeNames(), ",") {
		t.Errorf("mode.names = %v, want the vocabulary from risk regardless of the read", snap.Composer.Mode.Names)
	}
	if snap.Composer.Workspace.Set || snap.Composer.Workspace.Reason == "" {
		t.Errorf("workspace = %+v, want unset WITH a reason: 'no workspace' is a fact to show", snap.Composer.Workspace)
	}
	if snap.Composer.MaxAttachmentB != MaxAttachmentBytes {
		t.Errorf("maxAttachmentBytes = %d, want the broker's own default %d",
			snap.Composer.MaxAttachmentB, MaxAttachmentBytes)
	}
	if snap.Pending == nil || snap.Results == nil {
		t.Error("pending/results serialise as null; the renderer reads .length on both")
	}
	if len(snap.Pending) != 0 || len(snap.Results) != 0 {
		t.Errorf("pending=%v results=%v, want empty without inventing", snap.Pending, snap.Results)
	}
}

func TestAnUnreadableModeFromALiveReaderStillRendersUnknown(t *testing.T) {
	snap := NewSnapshotPump(PumpSources{
		Mode: func() risk.Mode { return risk.Mode(99) },
	}).Snapshot()
	if snap.Composer.Mode.Current != "unknown" {
		t.Fatalf("mode.current = %q, want unknown for a mode the store could not name",
			snap.Composer.Mode.Current)
	}
	if snap.Composer.Mode.Current == risk.ModeAutoApproveName {
		t.Error("an unreadable mode rendered as the loosest one")
	}
}

func TestTheStreamLogTruncatesInsteadOfMerging(t *testing.T) {
	// Ticket 197 leg B re-cut the overflow rule, so this case's old assertions -
	// "two rows, both texts kept, folded" - pin exactly the behaviour the ticket
	// called a lie. What it asserted that has nothing to do with folding is
	// asserted here still, in the stronger form: the bound holds, deltas
	// accumulate into their own key, Close ends THAT key, and no row is ever two
	// producers' text.
	sl := NewStreamLog(2)
	sl.Append("a", "one ")
	sl.Append("b", "two ")
	sl.Append("c", "three")
	sl.Close("c")

	chunks := sl.Chunks()
	if len(chunks) != 3 {
		t.Fatalf("chunks = %+v, want one row per key: folding two keys into one row is "+
			"a panel showing one agent doing two agents' work", chunks)
	}
	byKey := map[string]ResultChunk{}
	for _, c := range chunks {
		byKey[c.CorrelationID] = c
	}
	for key, want := range map[string]string{"a": "one ", "b": "two ", "c": "three"} {
		if byKey[key].Text != want {
			t.Errorf("row %q = %q, want exactly %q with no neighbour's text", key, byKey[key].Text, want)
		}
	}
	if !byKey["c"].Done || byKey["a"].Done || byKey["b"].Done {
		t.Errorf("done = c:%v a:%v b:%v, want only the key Close() named: the old fold OR-ed "+
			"a neighbour's Done onto the row", byKey["c"].Done, byKey["a"].Done, byKey["b"].Done)
	}
	if !sl.Truncated() {
		t.Error("3 keys through a 2-key bound reported no truncation")
	}
	if sl.ElidedRunes() != 0 {
		t.Errorf("elided = %d, want 0: short streams keep full fidelity past the bound", sl.ElidedRunes())
	}

	// The bound still has to hold under pressure, just not by joining rows: 50
	// keys through a 2-key log keeps at most the hard ceiling of rows, names every
	// key it stopped tracking, and never puts two deltas in one row.
	big := NewStreamLog(2)
	for i := 0; i < 50; i++ {
		big.Append("k"+strconv.Itoa(i), "x"+strconv.Itoa(i))
	}
	got := big.Chunks()
	if len(got) > 2*StreamKeyHardCeilingMultiple {
		t.Errorf("after 50 keys the log holds %d chunks, want the hard ceiling kept", len(got))
	}
	for _, c := range got {
		if want := "x" + strings.TrimPrefix(c.CorrelationID, "k"); c.Text != want {
			t.Fatalf("row %q = %q, want its own delta %q", c.CorrelationID, c.Text, want)
		}
	}
	if len(big.DroppedKeys()) != 50-len(got) {
		t.Errorf("dropped = %v, want the %d keys no row covers", big.DroppedKeys(), 50-len(got))
	}

	// Deltas accumulate into the key's own chunk; Close ends that stream.
	sl.Append("d", "alpha")
	sl.Append("d", "beta")
	sl.Close("d")
	var d ResultChunk
	found := false
	for _, c := range sl.Chunks() {
		if c.CorrelationID == "d" {
			d, found = c, true
		}
	}
	if !found || d.Text != "alphabeta" || !d.Done {
		t.Errorf("chunk d = %+v found=%v, want the appended deltas and done=true", d, found)
	}
}

func TestPublishWithoutAnExitReportsInsteadOfPretending(t *testing.T) {
	pump := NewSnapshotPump(PumpSources{
		Verdicts: func() []NativeVerdict {
			return []NativeVerdict{{CorrelationID: "c", Tool: "fs.write", Level: risk.L2}}
		},
	})
	snap, data, err := pump.Publish()
	if err == nil {
		t.Fatal("Publish reported success with no exit attached; that is a silent drop")
	}
	if !strings.Contains(err.Error(), "没有出口") {
		t.Errorf("err = %v, want it to name the missing last mile", err)
	}
	if len(data) == 0 {
		t.Error("no bytes: the packet itself should still be deliverable to a caller")
	}
	if len(snap.Pending) != 1 {
		t.Errorf("pending = %+v, want the verdict the reader had", snap.Pending)
	}
	if pump.Publishes() != 1 {
		t.Errorf("publishes = %d, want the attempt counted even though the exit refused it", pump.Publishes())
	}
}

func TestPublishHandsTheBytesToTheAttachedExit(t *testing.T) {
	var got []byte
	var gotSnap Snapshot
	pump := NewSnapshotPump(PumpSources{
		Mode: func() risk.Mode { return risk.ModeAskEveryStep },
		Now:  fixedClock(time.Date(2026, 9, 25, 13, 0, 0, 0, time.UTC)),
		Out: func(s Snapshot, data []byte) error {
			got, gotSnap = data, s
			return nil
		},
	})
	if _, _, err := pump.Publish(); err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("the exit received no bytes")
	}
	// The exit gets the same bytes Marshal produces for the same state, and the
	// packet is the four-key shape the renderer reads - not a Go-side struct
	// dump with unexported fields or an extra key.
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(got, &wire); err != nil {
		t.Fatal(err)
	}
	keys := make([]string, 0, len(wire))
	for k := range wire {
		keys = append(keys, k)
	}
	if strings.Join(sortedCopy(keys), ",") != "composer,generatedAt,pending,results" {
		t.Errorf("exit bytes carry keys %v, want the four PanelSnapshot declares", keys)
	}
	var composer map[string]json.RawMessage
	if err := json.Unmarshal(wire["composer"], &composer); err != nil {
		t.Fatal(err)
	}
	var mimes []string
	if err := json.Unmarshal(composer["acceptedAttachmentMimes"], &mimes); err != nil || len(mimes) == 0 {
		t.Errorf("acceptedAttachmentMimes = %s (err %v), want the broker's list, not null: "+
			"NewComposerState owns it and the pump must go through that constructor",
			composer["acceptedAttachmentMimes"], err)
	}
	var maxBytes int64
	if err := json.Unmarshal(composer["maxAttachmentBytes"], &maxBytes); err != nil || maxBytes != MaxAttachmentBytes {
		t.Errorf("maxAttachmentBytes = %s, want %d", composer["maxAttachmentBytes"], MaxAttachmentBytes)
	}
	if gotSnap.Composer.Mode.Current != risk.ModeAskEveryStepName {
		t.Errorf("exit saw %+v, want the same snapshot those bytes encode", gotSnap)
	}
	if !strings.Contains(string(got), `"generatedAt":"2026-09-25T13:00:00Z"`) {
		t.Errorf("exit bytes %s, want the injected clock's stamp", got)
	}
	if pump.Publishes() != 1 {
		t.Errorf("publishes = %d", pump.Publishes())
	}
}

// sortedCopy is local so the key-set nail above does not depend on anything but
// the standard library it already uses.
func sortedCopy(in []string) []string {
	out := append([]string(nil), in...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
