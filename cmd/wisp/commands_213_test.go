package main

// Ticket 213 AC#4 - the command roster's production caller, driven where it runs.
//
// The premise this file pins is NOT "cmdCommands prints something". It is the pair:
// the printed rows ARE internal/command's production rows (same names, same
// summaries, same arguments, same doors), and a command that has no capability is
// absent from the output because it is absent from the roster. A leg that kept a
// private table would pass the first half and fail the second, which is the failure
// ticket 213's AC#2 exists for ("菜单里有、点了没反应").
//
// The numbers below come from the roster the test reads at run time, never from
// literals written here: a row added to internal/command has to appear in this leg's
// output without an edit in this file, and a row removed has to disappear.

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/command"
)

// rows213 matches one printed roster row: two spaces, the slash-less name, two
// spaces, the summary. Continuation lines are indented four spaces, so they cannot
// be read as rows - and a row that loses its summary stops matching, which is the
// "clicked and nothing happened" shape this leg must not print.
var rows213 = regexp.MustCompile(`(?m)^  ([a-z][a-z0-9-]*)  (.+)$`)

func Test213CommandsLegPrintsTheProductionRoster(t *testing.T) {
	roster, err := command.Production()
	if err != nil {
		t.Fatalf("command.Production(): %v", err)
	}
	entries := roster.Entries()

	var out, errOut bytes.Buffer
	rc := cmdCommands(nil, commandsIO{stdout: &out, stderr: &errOut})
	if rc != 0 {
		t.Fatalf("cmdCommands() = %d, want 0; stderr:\n%s", rc, errOut.String())
	}
	text := out.String()

	// (1) The header counts the roster, not a literal.
	if want := fmt.Sprintf("%d 条命令", len(entries)); !strings.Contains(text, want) {
		t.Errorf("output does not carry %q; it opens with:\n%s", want, firstLines213(text, 3))
	}

	// (2) Row set equality: every printed name is a roster row and every roster row
	// is printed. A hand-written row in this leg reds here; a missing row reds here.
	got := map[string]string{}
	var order []string
	for _, m := range rows213.FindAllStringSubmatch(text, -1) {
		if _, dup := got[m[1]]; dup {
			t.Errorf("command %q printed twice", m[1])
		}
		if _, known := roster.Lookup(m[1]); !known {
			t.Errorf("row %q is printed and is not in the roster: this leg must not carry a table of its own (AC#2)", m[1])
		}
		got[m[1]] = m[2]
		order = append(order, m[1])
	}
	if len(order) != len(entries) {
		t.Errorf("printed %d rows, roster has %d: %v", len(order), len(entries), order)
	}
	for _, name := range roster.Names() {
		if _, ok := got[name]; !ok {
			t.Errorf("roster row %q is missing from the output", name)
			continue
		}
		e, _ := roster.Lookup(name)
		if got[name] != e.Summary {
			t.Errorf("row %q: printed summary = %q, roster says %q", name, got[name], e.Summary)
		}
	}

	// (3) Each row's three derived facts come from the row itself.
	for _, e := range entries {
		if !strings.Contains(text, e.Door()) {
			t.Errorf("row %q: the output never carries its door %q", e.Name, e.Door())
		}
		if !strings.Contains(text, "参数："+commandArgs(e.Args)) {
			t.Errorf("row %q: the output never carries its argument shape %q", e.Name, commandArgs(e.Args))
		}
		if !strings.Contains(text, "人批："+e.Approval.String()) {
			t.Errorf("row %q: the output never carries its approval %q", e.Name, e.Approval.String())
		}
	}

	// (4) The footer names the command-line spelling of each door present, derived
	// from the roster's own entries - so the footer cannot drift from the table.
	for _, e := range entries {
		if e.Channel != command.ChannelControlWord {
			continue
		}
		if want := fmt.Sprintf("wisp run %q", e.Utterance); !strings.Contains(text, want) {
			t.Errorf("footer does not name the control row %q's command-line spelling %q", e.Name, want)
		}
		if e.Method != "" {
			t.Errorf("control row %q carries a panel method %q: one command, one door", e.Name, e.Method)
		}
	}
	if method := panelMethod213(entries); method != "" && !strings.Contains(text, method) {
		t.Errorf("footer does not name the panel door %q", method)
	}
	if strings.Contains(text, "wisp panel-inbound") && !strings.Contains(text, "-data") {
		t.Errorf("footer names the panel-inbound leg without its -data requirement: that leg refuses to start without a root (ticket 128 AC#2)")
	}

	// (5) AC#4's named half. The terminal state is "visible on the command line,
	// not clickable in the page", and this leg says so out loud rather than leaving
	// the operator to infer it from a missing menu entry.
	if !strings.Contains(text, "界面点不了") {
		t.Errorf("output does not say the page cannot ring these commands yet: AC#4's terminal state is 名册可见 + 命令行可执行, 界面点不了")
	}

	// (6) AC#2 read off the production output: capabilities this machine does not
	// have are absent. These are the ticket's own eight minus the ones that exist,
	// and each absence is argued in internal/command's roster comment.
	for _, absent := range []string{"compact", "goal", "plan", "model", "export", "feedback"} {
		if _, ok := got[absent]; ok {
			t.Errorf("row %q is printed, and this machine has no door for it: a listed command with no capability is the exact failure AC#2 forbids", absent)
		}
	}
}

// Test213CommandsLegRefusesArguments pins the operator-mistake path: this leg takes
// no arguments, says so, prints its usage on stderr, and puts no row on stdout (an
// exit code alone would let a caller parse a half-printed table as success).
func Test213CommandsLegRefusesArguments(t *testing.T) {
	var out, errOut bytes.Buffer
	rc := cmdCommands([]string{"list"}, commandsIO{stdout: &out, stderr: &errOut})
	if rc != 2 {
		t.Errorf("cmdCommands(list) = %d, want 2 (usage error, the code the other console legs use)", rc)
	}
	if out.Len() != 0 {
		t.Errorf("a refused invocation wrote %d bytes to stdout:\n%s", out.Len(), out.String())
	}
	if !strings.Contains(errOut.String(), "不接受参数") {
		t.Errorf("stderr never says the argument was refused:\n%s", errOut.String())
	}
	if !strings.Contains(errOut.String(), "wisp commands") {
		t.Errorf("stderr does not print this leg's usage block:\n%s", errOut.String())
	}
}

// panelMethod213 returns the roster's panel-door method, or "" when the roster has
// no such row (the footer then has no envelope to show and the check above is
// skipped rather than vacuously satisfied).
func panelMethod213(entries []command.Command) string {
	for _, e := range entries {
		if e.Channel == command.ChannelPanelMode {
			return e.Method
		}
	}
	return ""
}

// firstLines213 quotes at most n lines of s, for a red that a reader can read
// without scrolling the whole listing.
func firstLines213(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}
