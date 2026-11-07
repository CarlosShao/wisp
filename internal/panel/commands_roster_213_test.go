package panel

// Ticket 213 交付 3's carrier tests. The commands key is a reader-gated
// section like instructions and tasks, so its tests pin the same two
// directions those sections' tests pin:
//
//   - WITHOUT the reader the key is absent from the bytes a real pump
//     marshals - its presence elsewhere is caused by the reader, never by a
//     field that is always on the wire with nothing behind it (the three
//     four-key byte nails keep pinning that shape; this case names the
//     commands key itself).
//   - WITH the reader the key always arrives, rows travel verbatim (deep
//     equality, not a spot check), nil rows render as [] and never null, and
//     the defect shape (Reason set, rows empty) arrives as written instead of
//     being dressed down to an empty menu.
//
// A mutation that drops the pump wiring, drops a field from the carrier, or
// folds the defect shape into an empty menu reds here.

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTheCommandsKeyTravelsOnlyFromItsReader213(t *testing.T) {
	// (1) No reader, no key: the four-key shape the byte nails pin, and the
	// commands key named so a regression here says which key came back.
	data, err := NewSnapshotPump(PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"commands"`) {
		t.Fatalf("a pump with no commands reader sent the key: %s", data)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if len(generic) != 4 {
		t.Errorf("top-level keys = %d, want the four pump_test pins for a reader-less pump", len(generic))
	}

	// (2) With a reader the key arrives and the rows travel verbatim: every
	// field of every row is the value the reader returned, in order.
	want := CommandRosterSection{Rows: []CommandRow{
		{Name: "cancel", Summary: "取消当前正在跑的任务", Args: "", Approval: "不需要", Door: "D11(1) 控制词「取消」"},
		{Name: "mode", Summary: "切换权限预设", Args: "目标档位", Approval: "看参数", Door: "C17 panel.mode.request"},
	}}
	pump := NewSnapshotPump(PumpSources{
		Commands: func() CommandRosterSection { return want },
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	raw, ok := generic["commands"]
	if !ok {
		t.Fatalf("a pump with a commands reader sent no commands key: %s", data)
	}
	var got CommandRosterSection
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Rows) != len(want.Rows) {
		t.Fatalf("wire rows = %d, want %d: %s", len(got.Rows), len(want.Rows), raw)
	}
	for i, w := range want.Rows {
		if got.Rows[i] != w {
			t.Errorf("wire row %d = %+v, want %+v (every field travels verbatim, or the menu lies about a command's cost)", i, got.Rows[i], w)
		}
	}

	// (3) Nil rows from a reader render as [] - the page reads .length on
	// this list like it does on pending/results, and null breaks that.
	pump = NewSnapshotPump(PumpSources{
		Commands: func() CommandRosterSection { return CommandRosterSection{} },
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"rows":[]`) {
		t.Errorf("nil rows did not render as []: %s", data)
	}

	// (4) The defect shape travels as written: a roster whose own AC#2 gate
	// refused arrives with the refusal sentence and empty rows - never
	// silently folded into an empty menu.
	defect := CommandRosterSection{
		Rows:   []CommandRow{},
		Reason: "名册装配被拒（AC#2 的闸拒了生产表本身）：命令 \"x\" 的门不认识",
	}
	pump = NewSnapshotPump(PumpSources{
		Commands: func() CommandRosterSection { return defect },
	})
	data, err = pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	raw, ok = generic["commands"]
	if !ok {
		t.Fatalf("the defect shape lost its key: %s", data)
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.Reason != defect.Reason || len(got.Rows) != 0 {
		t.Errorf("defect shape = %+v, want the refusal sentence carried verbatim with empty rows", got)
	}
}
