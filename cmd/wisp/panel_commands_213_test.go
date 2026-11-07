package main

// Ticket 213 交付 3's host-half tests, and the teeth the carrier rides on:
//
//   - The wire section is DERIVED from command.Production(), row by row, field
//     by field - this leg keeps no table of its own. A mutation that swaps any
//     field's source, drops a row, or appends a row the roster never proved
//     reds the first case (the same AC#2 failure shape the console leg's test
//     pins, one hop deeper: on the wire the page would read).
//   - The assembly root actually carries the reader (a disk reading, the same
//     shape ticket 214's test pins for its two readers): removing the run.go
//     line reds the second case, and the key would silently vanish from every
//     snapshot the process pushes.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/command"
)

// Test213CommandRosterSectionIsTheProductionRoster pins the derivation: the
// section's rows equal the roster's own rows - name, summary, argument shape,
// the approval sentence and the door - with no Reason (the gate did not fire).
func Test213CommandRosterSectionIsTheProductionRoster(t *testing.T) {
	roster, err := command.Production()
	if err != nil {
		t.Fatalf("command.Production(): %v", err)
	}
	entries := roster.Entries()
	if len(entries) == 0 {
		t.Fatal("the production roster is empty; this case would be vacuous")
	}

	got := commandRosterSection213()
	if got.Reason != "" {
		t.Fatalf("section carries a refusal sentence with a loaded roster: %q", got.Reason)
	}
	if len(got.Rows) != len(entries) {
		t.Fatalf("section rows = %d, roster rows = %d: a row was dropped or invented on the wire", len(got.Rows), len(entries))
	}
	for i, e := range entries {
		row := got.Rows[i]
		if row.Name != e.Name {
			t.Errorf("row %d: name = %q, roster says %q", i, row.Name, e.Name)
		}
		if row.Summary != e.Summary {
			t.Errorf("row %d (%s): summary = %q, roster says %q", i, e.Name, row.Summary, e.Summary)
		}
		if row.Args != e.Args {
			t.Errorf("row %d (%s): args = %q, roster says %q", i, e.Name, row.Args, e.Args)
		}
		if want := e.Approval.String(); row.Approval != want {
			t.Errorf("row %d (%s): approval = %q, the row's own rendering is %q", i, e.Name, row.Approval, want)
		}
		if want := e.Door(); row.Door != want {
			t.Errorf("row %d (%s): door = %q, the row's own rendering is %q", i, e.Name, row.Door, want)
		}
	}
}

// Test213AssemblyRootCarriesTheCommandsReader is the disk reading: run.go's
// PumpSources literal hands the pump the commands reader, and the file that
// derives the rows derives them from command.Production() rather than keeping
// a private table (a literal command name spelled in the builder is a private
// table by definition).
func Test213AssemblyRootCarriesTheCommandsReader(t *testing.T) {
	runGo, err := os.ReadFile(filepath.Join("run.go"))
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	if !strings.Contains(string(runGo), "Commands: commandRosterSection213") {
		t.Error("run.go's PumpSources literal no longer carries the commands reader: the snapshot's commands key would silently vanish from every packet this process pushes (ticket 213 交付 3)")
	}

	builder, err := os.ReadFile(filepath.Join("panel_commands_213.go"))
	if err != nil {
		t.Fatalf("read panel_commands_213.go: %v", err)
	}
	text := string(builder)
	if !strings.Contains(text, "command.Production()") {
		t.Error("the commands section no longer derives from command.Production(): a wire roster that is not the production roster is the exact private-table failure ticket 213 AC#2 forbids")
	}
	if strings.Contains(text, `Name: "`) {
		t.Error("the builder spells a literal command name: that is a private table riding under the roster's name")
	}
}
