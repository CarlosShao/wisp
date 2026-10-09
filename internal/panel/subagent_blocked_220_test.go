package panel

// Ticket 220's join cases: the roster's blockedOnApproval cell must read TRUE for
// a task that is actually waiting on an approval, and it must keep reading true
// for the shape production silently drops today.
//
// WHAT THE DEFECT WAS, in the two shapes the ticket names:
//
//  1. The set of waiting ids was keyed by the card's correlation id
//     (subagent_roster_197.go:190-195) and looked up by the row's task id (:210).
//     Those two agree ONLY because a decision whose CorrelationID is empty falls
//     back to the task id (internal/agent/approval/gate.go:275 and :471). The
//     queue's clash rewrite breaks even that: a second card for the SAME task is
//     filed as "<original id>#<seq>" (internal/agent/approval/queue.go:152-154),
//     and a lookup by "task-7" cannot find "task-7#9". The busier a task is, the
//     more its waiting reads as "not waiting".
//  2. The card list the join reads comes from Queue().LiveApprovals() - the L2
//     queue only (cmd/wisp's liveVerdicts). An L1 confirmation window is not a
//     queue item (see gate.go's own note at :309-313), so a task blocked in that
//     2-3s block had no producer at all until ticket 220 AC#2's read-only
//     enumeration. PumpSources.L1Windows is the panel-side half of that hop: a
//     READER, and a nil reader says "this host cannot see L1 windows" instead of
//     saying "nothing is waiting".
//
// WHAT THESE CASES DELIBERATELY DO NOT DO: they do not turn the boolean into a
// state name (AC#4 - D43's transition table is frozen, and blockedOnApproval
// stays a bool), and they do not let the blocked dimension REPLACE the running
// dimension. The ticket's 外部对照 section names a reference implementation whose
// roster priority (需要批准 > 问了问题 > 还在跑 > Done) hides "still running" behind
// "needs approval"; the last case below pins that we do the opposite - the row
// carries its filed D43 state AND the fact that one card is waiting.

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// card197 is one pending L2 card as the host reads it off its own queue: the
// correlation id is the only id the card carries (ApprovalCardView has no task
// field, and giving it one would widen the frozen card JSON keys).
func card197(corr string) NativeVerdict {
	return NativeVerdict{
		CorrelationID: corr, Tool: "fs.write", Level: risk.L2,
		RulesHit: []risk.RuleID{risk.R2}, Reason: "写在工作区之外",
	}
}

// rows197 are the two rows every case below shares: a root and the child it
// derived, with the D43 name the host filed for each.
func rows197() []TaskRow {
	return []TaskRow{
		{
			TaskID: "root-1", Label: "总结这份笔记", Kind: "root",
			Status: "Thinking", StatusKnown: true, StreamKey: "root-1",
		},
		{
			TaskID: "child-1", ParentTaskID: "root-1", Label: "读三份文件",
			Kind: "subagent", Status: "Thinking", StatusKnown: true,
			StreamKey: SubagentStreamKey("child-1"),
		},
	}
}

// blockedPacket is the published bytes for one assembly: these rows, these cards,
// these L1 windows. Everything goes through Publish so the assertion reads the
// wire, not a Go struct (the rule this package's roster cases were written on).
func blockedPacket(t *testing.T, cards []NativeVerdict, waits []L1WindowWait) []byte {
	t.Helper()
	pump := NewSnapshotPump(PumpSources{
		Out:      okExit197,
		Verdicts: func() []NativeVerdict { return cards },
		Mode:     func() risk.Mode { return risk.ModeAskHighRisk },
		Now:      func() time.Time { return time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC) },
		Tasks: func() TaskRosterState {
			return TaskRosterState{Rows: rows197(), InFlightSlots: 1, PoolCap: poolCapGiven197}
		},
		L1Windows: func() []L1WindowWait { return waits },
	})
	_, data, err := pump.Publish()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// assertBlocked reads one row off the wire and states what it found in BOTH
// directions, so a case cannot pass because the row it is looking for is absent.
func assertBlocked(t *testing.T, data []byte, taskID string, want bool) {
	t.Helper()
	sect := packetRoster(t, data)
	row := wireRow(t, sect, taskID)
	if row.BlockedOnApproval != want {
		t.Errorf("row %q blockedOnApproval=%v，期望 %v； packets=%s",
			taskID, row.BlockedOnApproval, want, data)
	}
}

// TestASecondCardForTheSameTaskStillLightsItsRow is AC#1's positive control and
// AC#3's second shape: the queue rewrote this card's correlation id to
// "<task id>#<seq>" because the task's first card still holds the plain name, and
// the task IS waiting right now. Before the join fix this row read false, which
// is the silent drop the ticket was filed for.
func TestASecondCardForTheSameTaskStillLightsItsRow(t *testing.T) {
	data := blockedPacket(t, []NativeVerdict{card197("child-1#9")}, nil)
	assertBlocked(t, data, "child-1", true)
	// Same assembly, both cards pending: one waiting card is enough, and the row
	// stays true rather than flickering with the card count.
	data = blockedPacket(t, []NativeVerdict{card197("child-1"), card197("child-1#9")}, nil)
	assertBlocked(t, data, "child-1", true)
	// And the task nobody filed a card for still reads false - the widened match
	// must not become "anything pending blocks everything".
	assertBlocked(t, data, "root-1", false)
}

// TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt is the sensitivity half: the
// widened key must be exactly the queue's own rewrite ("<id>#<decimal seq>",
// queue.go:152-154), not "anything that starts with my task id". Drop the digit
// test and a prefix match passes; drop the other-task test and a blanket match
// passes.
func TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt(t *testing.T) {
	cases := []struct {
		name string
		corr string
		want bool
	}{
		{"the queue's own rewrite lights", "child-1#9", true},
		{"a multi-digit seq lights", "child-1#147", true},
		{"another task's rewrite does not light", "other-1#9", false},
		{"a non-numeric suffix is not the queue's form", "child-1#retry", false},
		{"a task id that merely prefixes does not light", "child-1x", false},
		{"an empty correlation id names no row", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertBlocked(t, blockedPacket(t, []NativeVerdict{card197(c.corr)}, nil), "child-1", c.want)
		})
	}
}

// TestAL1WindowInWaitingLightsItsRow is AC#2's landing (甲: the read-only
// enumeration) at the pump, and AC#3's first shape. The window is reported by the
// host's enumeration under BOTH of its ids, and the row that matters here is named
// by the task id while the window's own correlation id is a value the host filed
// for another purpose - the corr/task 双查 AC#1 asks for, exercised on the side of
// the join that never had a producer.
func TestAL1WindowInWaitingLightsItsRow(t *testing.T) {
	data := blockedPacket(t, nil, []L1WindowWait{{CorrelationID: "corr-host-9", TaskID: "child-1"}})
	assertBlocked(t, data, "child-1", true)
	assertBlocked(t, data, "root-1", false)

	// Since bd124b2a the loop's corr is per call (<taskID>#<callID>, loop.go:603,
	// used :676), never == taskID: the corr-only pairing below is this fixture's.
	data = blockedPacket(t, nil, []L1WindowWait{{CorrelationID: "child-1", TaskID: ""}})
	assertBlocked(t, data, "child-1", true)

	// A host assembled WITHOUT the enumeration is a statement about the assembly,
	// not a claim that nothing is waiting: the same window, unread, reads false.
	// This is the reading the pre-fix tree gave for every L1 window in the tree.
	noread := NewSnapshotPump(PumpSources{
		Out:      okExit197,
		Verdicts: func() []NativeVerdict { return nil },
		Mode:     func() risk.Mode { return risk.ModeAskHighRisk },
		Now:      func() time.Time { return time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC) },
		Tasks: func() TaskRosterState {
			return TaskRosterState{Rows: rows197(), InFlightSlots: 1, PoolCap: poolCapGiven197}
		},
	})
	_, data, err := noread.Publish()
	if err != nil {
		t.Fatal(err)
	}
	assertBlocked(t, data, "child-1", false)
}

// TestABlockedRowKeepsItsRunningDimension is the ticket's 外部对照 constraint as
// code: the blocked cell is an ADDED dimension, never a replacement for the state
// the host filed. The reference implementation this constraint was drawn against
// reports only one value per row and its priority (需要批准 > 问了问题 > 还在跑)
// swallows "still running"; here both halves travel together, and neither one is
// derived from the other.
func TestABlockedRowKeepsItsRunningDimension(t *testing.T) {
	data := blockedPacket(t, []NativeVerdict{card197("child-1#9")},
		[]L1WindowWait{{CorrelationID: "corr-root", TaskID: "root-1"}})
	sect := packetRoster(t, data)
	child := wireRow(t, sect, "child-1")
	if !child.BlockedOnApproval {
		t.Error("the second card is pending and the row does not say so")
	}
	if child.Status != "Thinking" || !child.StatusKnown {
		t.Errorf("blockedOnApproval took the running dimension with it: status=%q known=%v, "+
			"want the host's filed Thinking carried untouched (AC#4: this cell is a boolean, "+
			"not a state, and it never rewrites one)", child.Status, child.StatusKnown)
	}
	root := wireRow(t, sect, "root-1")
	if !root.BlockedOnApproval || root.Status != "Thinking" || !root.StatusKnown {
		t.Errorf("root row = %+v, want blocked AND still carrying its filed state", root)
	}

	// The wire keys are unchanged: a boolean gaining a producer is not a new key,
	// and the four-way / two-way reconciliations stay the only rulers of the shape.
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatalf("packet is not JSON: %v", err)
	}
	raw, ok := top["tasks"]
	if !ok {
		t.Fatal("the packet carries no tasks key")
	}
	if strings.Contains(string(raw), `"blockedState"`) || strings.Contains(string(raw), `"waitingOn"`) {
		t.Fatalf("the join invented a new wire key instead of filling the existing boolean: %s", raw)
	}
}
