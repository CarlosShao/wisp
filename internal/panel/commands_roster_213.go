package panel

// Ticket 213 交付 3's carrier: the command roster as a snapshot section.
//
// WHAT THIS IS. internal/command owns the roster (one enumerable table of the
// commands a person can ring, every row proved against a real door by its own
// AC#2 gate). What the chat box's plus menu needs is those rows ON THE WIRE,
// and this file is the wire shape: one row per command, carrying the four
// facts the ticket's 交付 1 names - a slash-less name, one sentence of what it
// does, the argument shape, whether a human approval is involved - plus the
// door that executes it, in this package's own words.
//
// IT CARRIES, IT DOES NOT DECIDE. No row here re-derives a verdict: Approval
// and Door arrive as the strings internal/command's own row renders, so a
// menu and `wisp commands` cannot disagree about what a command costs. The
// per-request judgement (whether THIS switch costs a card under THIS 档)
// stays with the doors - perm.JudgeCard behind the mode route - exactly as
// ticket 213 AC#3 landed it; a section that judged would be a second policy,
// and R9's direction rule already has one home.
//
// NO EXECUTION PATH ARRIVES WITH IT. The section is display data: it names
// doors, it opens none. No C17 method is added by this carrier, nothing here
// answers a postMessage, and the control words it lists still have no page
// entry at all (a page's text reaches no agent.Loop today) - the menu being
// able to DRAW a row has never been the row being RINGABLE, and the wire
// keeps that distinction by carrying nothing but the row's own facts.
//
// THE CONTRACT MOVE IS EXPLICIT. This key joins the reader-gated sections
// (instructions, 200-r2; tasks, 197-r3) in moving the snapshot's top-level
// key set: a pump assembled without the reader sends no key (the three
// four-key byte nails keep their exact shape), and a pump with one always
// sends it. The two-way reconciliations against frontend/src/lib/panel.ts
// name "commands" as a key the interface does not declare until the page
// side writes its declaration - that red is ticket 213 交付 3's own
// procedure (红因写进票面，由 owner 带给界面那支，不改尺、不加豁免), and
// the page-side declaration is the UI leg's to write (Q-51).

// CommandRow is one roster row on the wire. Name carries no slash (the menu
// draws it); Args is the roster's own argument shape, where empty means the
// command takes none - the roster's documented semantics, carried verbatim,
// not re-rendered into a display string here. Approval and Door are the
// row's own sentences (command.Approval.String() and command.Command.Door()):
// one spelling of what the command costs and which door runs it, shared with
// the console leg instead of forked for the page.
type CommandRow struct {
	Name     string `json:"name"`
	Summary  string `json:"summary"`
	Args     string `json:"args"`
	Approval string `json:"approval"`
	Door     string `json:"door"`
}

// CommandRosterSection is the snapshot's commands key. Rows is the roster
// this machine can execute today; empty means this machine can execute none,
// never "unknown" - the rows come from a table whose own gate refuses any
// row without a door, so an empty roster is a measured fact.
//
// Reason is the one case where the section says something other than rows:
// the roster's production table itself refused to assemble (its AC#2 gate
// fired on a row that has no door). That is a defect, not an empty menu, and
// the wire says which - an empty rows array with no Reason would invite the
// page to render "此机无命令" over a table that is broken, which is the
// 文案在、控件不在 family ledger A408 names. Reason empty = the roster
// loaded; Reason set = rows is empty and the sentence names the refusal.
type CommandRosterSection struct {
	Rows   []CommandRow `json:"rows"`
	Reason string       `json:"reason,omitempty"`
}
