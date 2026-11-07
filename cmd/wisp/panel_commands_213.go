package main

// Ticket 213 交付 3's host half: the command roster on the snapshot wire.
//
// WHAT THIS IS. internal/command owns the roster and its AC#2 gate; the pump
// carries the rows (internal/panel/commands_roster_213.go). This file is the
// one hop between them, and it is deliberately the ONLY place that touches
// both: the roster is the table, the section is the wire shape, and what
// lives here is the row-by-row derivation - nothing else.
//
// IT DERIVES, IT DOES NOT KEEP. Every row below is read out of
// command.Production() at build time: name, summary, argument shape, the
// approval sentence and the door are the roster's own row rendered. A private
// table here would keep printing rows the roster refused (or hide rows it
// proved) - exactly the "菜单里有、点了没反应" shape ticket 213's AC#2
// exists for - so this file spells no command name of its own and the
// disk-reading test pins the derivation against the live table.
//
// BUILT ONCE, BECAUSE THE TABLE IS IMMUTABLE. The roster is fixed for the
// life of the process (its own type comment: "It is immutable once built"),
// so the section is built on first read and cached; every snapshot push in
// this process carries the same rows rather than re-deriving them per push.
//
// THE FAILURE PATH IS THE ROSTER'S OWN SENTENCE. Production() errors only
// when the AC#2 gate refuses the PRODUCTION table itself - a row with no
// door. That is a defect, not an empty menu: the wire carries the refusal
// sentence verbatim (Reason) with empty rows, so the page can say the roster
// is broken instead of rendering "此机无命令" over a broken table. Nothing
// here invents a row to keep the menu looking alive.

import (
	"fmt"
	"sync"

	"github.com/CarlosShao/wisp/internal/command"
	"github.com/CarlosShao/wisp/internal/panel"
)

// commandRoster213 caches the section: one build per process, over an
// immutable table.
var commandRoster213 struct {
	sync.Once
	section panel.CommandRosterSection
}

// commandRosterSection213 is the pump's Commands reader (PumpSources.Commands,
// wired at run.go's assembly root). It is a plain function value, not a
// runtime method, on purpose: the roster is the table, not this process's
// running state - no live object is read, and no nil-tolerance dance belongs
// to it.
func commandRosterSection213() panel.CommandRosterSection {
	commandRoster213.Do(buildCommandRoster213)
	return commandRoster213.section
}

// buildCommandRoster213 derives the wire rows from the production roster.
func buildCommandRoster213() {
	roster, err := command.Production()
	if err != nil {
		commandRoster213.section = panel.CommandRosterSection{
			Rows: []panel.CommandRow{},
			Reason: fmt.Sprintf("名册装配被拒（AC#2 的闸拒了生产表本身，这是缺陷不是没有命令）：%v",
				err),
		}
		return
	}
	entries := roster.Entries()
	rows := make([]panel.CommandRow, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, panel.CommandRow{
			Name:     e.Name,
			Summary:  e.Summary,
			Args:     e.Args,
			Approval: e.Approval.String(),
			Door:     e.Door(),
		})
	}
	commandRoster213.section = panel.CommandRosterSection{Rows: rows}
}
