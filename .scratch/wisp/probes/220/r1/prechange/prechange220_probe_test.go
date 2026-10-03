package panel

// PRE-CHANGE PROBE for ticket 220 (220-r1). This file exists to produce the
// reading the ticket demands on the spot - "改前的读数要当场跑出来" - using only
// the API that exists before the fix. It is relocated to
// .scratch/wisp/probes/220/r1/prechange/ once that reading is booked; the
// resident cases live in subagent_blocked_220_test.go.

import (
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

func prePacket220(t *testing.T, corrs ...string) []byte {
	t.Helper()
	cards := make([]NativeVerdict, 0, len(corrs))
	for _, c := range corrs {
		cards = append(cards, NativeVerdict{
			CorrelationID: c, Tool: "fs.write", Level: 2, Reason: "写在工作区之外",
		})
	}
	pump := NewSnapshotPump(PumpSources{
		Out:      okExit197,
		Verdicts: func() []NativeVerdict { return cards },
		Mode:     func() risk.Mode { return risk.ModeAskHighRisk },
		Now:      func() time.Time { return time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC) },
		Tasks: func() TaskRosterState {
			return TaskRosterState{
				Rows: []TaskRow{
					{TaskID: "root-1", Label: "总结这份笔记", Kind: "root",
						Status: "Thinking", StatusKnown: true, StreamKey: "root-1"},
					{TaskID: "child-1", ParentTaskID: "root-1", Label: "读三份文件",
						Kind: "subagent", Status: "Thinking", StatusKnown: true,
						StreamKey: SubagentStreamKey("child-1")},
				},
				InFlightSlots: 1, PoolCap: poolCapGiven197,
			}
		},
	})
	_, data, err := pump.Publish()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestPRECHANGESecondCardOfOneTask(t *testing.T) {
	sect := packetRoster(t, prePacket220(t, "child-1#9"))
	row := wireRow(t, sect, "child-1")
	if !row.BlockedOnApproval {
		t.Errorf("PRE-CHANGE READING: the same task's second card is pending as %q and the "+
			"row still reads blockedOnApproval=%v (want true)", sect.Rows[1].StreamKey, row.BlockedOnApproval)
	}
	// The one-card shape is what the fallback pairing makes work today, so this
	// half must already read true - the defect is only the second card.
	sect = packetRoster(t, prePacket220(t, "child-1"))
	if !wireRow(t, sect, "child-1").BlockedOnApproval {
		t.Errorf("PRE-CHANGE READING: even the plain corr==taskID card fails to light the row")
	}
}
