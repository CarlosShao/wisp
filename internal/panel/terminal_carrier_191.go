package panel

import "sync"

// Ticket 191 wave 2a's carrier: the snapshot's top-level "terminal" key - the
// panel-embedded terminal's two read-only readouts, the INPUT PORT (the
// command line the user typed, display only) and the OUTPUT STREAM (a rolling
// window over the most recent output the terminal session produced).
//
// THE PANEL TERMINAL DOES NOT EXECUTE ANYTHING (ticket 191 AC#0's 甲 ruling,
// ledger A379): a command typed at the terminal travels to the SAME gated
// tool leg every other command takes (shell.session, ticket 163, streaming
// face delivered by wave one 6512d5a73's SessionOptions.OnChunk) - the gate,
// the audit line and the L2 cards are not bypassed, there is no panel-side
// process creation, and this file is the shape proof of that sentence: its
// import list is exactly one package (sync) and the AST nail in
// terminal_carrier_191_test.go refuses any process constructor or second
// import here. The feed's methods take text and hand text back; there is no
// method a page could call to make anything run, and the method-set ruler in
// the same test file turns the arrival of one red.
//
// READER-GATED, THE 213/187/65/193 PATTERN, NO SECOND PARADIGM. The pump
// assembled without the Terminal reader sends no "terminal" key at all (the
// shape the four-key byte nails pin), and a pump that HAS one always sends
// the key - with an empty section rendered as "terminal":{} rather than a
// vanished key, because "the terminal carrier is wired" and "nothing has
// been typed or streamed yet" are two different facts a page must not
// render the same. No C17 method arrives with it: the section is a display
// surface, its write half (if any is ever ruled) rides the composer's
// existing send doors, and the cmd/wisp pump-source wiring is wave 2b's to
// land - until then the running process carries no terminal key at all.
//
// THE TWO FIELDS ARE ARRIVAL FACTS, NOT STATE FILLS (tooth 2). Input is
// absent until a line arrived; Output is absent until a chunk arrived. An
// empty string on the wire would read as "the command printed nothing",
// which is a DIFFERENT fact from "no output ever reached this carrier" -
// absent is the truth, never a blank dressed as data. Within a present
// Output row, Truncated=false means "the window still holds everything the
// stream sent so far", which is information, not absence - that is why the
// flags live inside the rows and the rows themselves carry omitempty.
//
// THE WINDOW IS BOUNDED, AND THE BOUND IS COUNTED (tooth 1). The output
// window keeps the most recent tail and evicts the head past
// terminalWindowKeepRunes; every evicted rune is counted in DroppedRunes and
// Truncated flips true the first time anything was evicted, for the feed's
// whole life. A bound without its count would turn "the head fell off" into
// silence - the same honesty rule RunFeed's toolsDropped and StreamLog's
// dropped ledger state. The input line is bounded the same way
// (terminalInputKeepRunes): most-recent-line-wins never accumulates, but one
// pasted line can still be arbitrarily long, so it gets the same tail window
// and the same counted eviction.
//
// RUNES, NOT BYTES - the runReasoningKeepRunes rule verbatim: a byte cut
// would split the Chinese this project streams and hand the panel mojibake.
// The tool side's own 64MiB transport cap (sessionPendingCap) bounds the
// pipe; this window bounds the DISPLAY, and the two bounds are different
// jobs that must not be confused.
//
// THE LEAK HALF IS NOT THIS FILE'S (ticket 191 AC#2, Q-457 pending the
// owner): command output naturally carries absolute paths, tokens and
// environment variables, and the C25 stamping / leak-detection pipe (tickets
// 175/177/183's pipe) is the ruled route for that. Q-457 is open, so this
// carrier neither rides that pipe nor pretends to: the fields carry verbatim
// what arrived, and the wave-2b pump-source wiring may not attach a source
// that sidesteps the pipe once it is ruled. Scrubbing or redacting here
// would be inventing a second, weaker pipe - not done.
//
// The wire shape is registered as a judged page-side gap in
// composerTSGapLedger's PanelSnapshot entry (the page half belongs to the
// frontend leg; frontend/** is outside this wave's write face) - the same
// registration procedure every Go-only wave before it used.

// terminalWindowKeepRunes is the output window's tail bound. 8192 runes is
// twice the reasoning bound (runReasoningKeepRunes): the terminal is the
// flood face of the two, and the bound has to survive a few screens of
// scrollback while staying a bounded price.
const terminalWindowKeepRunes = 8192

// terminalInputKeepRunes is the input line's tail bound, matched to the
// reasoning bound: a command line lives under the same per-region volume
// rule, and a pasted megabyte gets its tail kept and its head counted.
const terminalInputKeepRunes = 4096

// TerminalSection is the snapshot's top-level "terminal" key. A section that
// exists always renders, even when both rows are still absent - "terminal":
// {} is the honest "wired, nothing yet" state, and the rows' own absence is
// the honest "nothing arrived" state. The key's presence on the wire is the
// reader's existence; nothing here opens a door and no C17 method arrives
// with it.
type TerminalSection struct {
	// Input is the INPUT PORT: the command line the user typed, verbatim,
	// display only. Most recent line wins - a terminal input port shows the
	// line the user is on, and the session's own echo carries the history in
	// the output window. Absent (never an empty string) until a line arrived.
	Input *TerminalInputRow `json:"input,omitempty"`
	// Output is the OUTPUT STREAM: the rolling tail window over what the
	// session produced, with the window's own honesty flags. Absent (never
	// an empty string) until a chunk arrived.
	Output *TerminalOutputRow `json:"output,omitempty"`
}

// TerminalInputRow is one command line the user typed at the panel terminal.
// The line is the producer's own text verbatim; this row derives nothing.
type TerminalInputRow struct {
	// Text is the retained tail of the line, verbatim.
	Text string `json:"text"`
	// Truncated is the honest half of the input bound: true once any part of
	// the line was evicted by terminalInputKeepRunes. A present row with
	// false means the line fits whole.
	Truncated bool `json:"truncated"`
	// DroppedRunes counts the head runes the bound evicted. A bound without
	// its count would turn "the head fell off" into silence.
	DroppedRunes int `json:"droppedRunes"`
}

// TerminalOutputRow is the output stream's rolling window. Text is the
// retained tail, in arrival order; the two flags are the window's own
// honesty ledger.
type TerminalOutputRow struct {
	// Text is the retained tail of the stream, verbatim, in arrival order.
	Text string `json:"text"`
	// Truncated is true once any part of the stream was evicted by
	// terminalWindowKeepRunes, for the feed's whole life - a window that ate
	// bytes once must not quietly pretend it did not.
	Truncated bool `json:"truncated"`
	// DroppedRunes counts the head runes the bound evicted, cumulatively.
	// The pair (Truncated, DroppedRunes) is the same honesty rule
	// RunSection.ToolsDropped states: eviction is counted, never silent.
	DroppedRunes int `json:"droppedRunes"`
}

// TerminalFeed is the holder the wave-2b pump source feeds and the pump's
// reader reads. It is a RunFeed-shaped dependency of the snapshot pump: the
// host that owns the terminal session feeds the lines and chunks it already
// has, the assembly root hands Section to PumpSources.Terminal, and the
// mutex covers the two goroutines that genuinely meet here - whatever
// delivers session output and whatever drives the pump.
//
// The feed methods take plain strings, not tool events, on purpose: this
// package does not import internal/tools, and the narrow seams keep the
// carrier honest about exactly what it was given - text that arrived,
// nothing that executes.
type TerminalFeed struct {
	mu sync.Mutex

	// input state: the most recent line, tail-windowed, with its own seen
	// count so an evicted head stays counted.
	inputTail    []rune
	inputSeen    int
	inputDropped int

	// output state: the rolling window, tail-windowed, cumulative.
	windowTail    []rune
	windowSeen    int
	windowDropped int
}

// NewTerminalFeed builds an empty feed.
func NewTerminalFeed() *TerminalFeed { return &TerminalFeed{} }

// Input records the command line the user typed, verbatim, display only.
// The most recent line replaces the previous one - an input port shows the
// line the user is on, and this method never accumulates history. An empty
// string writes nothing: no line arrived, and an empty line is not a fact
// this carrier invents. The line is bounded by terminalInputKeepRunes with
// the head counted, never silently dropped.
func (f *TerminalFeed) Input(line string) {
	if f == nil || line == "" {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r := []rune(line)
	f.inputTail = append(f.inputTail[:0], r...)
	f.inputSeen += len(r)
	if over := len(f.inputTail) - terminalInputKeepRunes; over > 0 {
		f.inputTail = f.inputTail[over:]
		f.inputDropped += over
	}
}

// OutputChunk appends one chunk of session output to the rolling window, in
// arrival order. The window keeps the tail past terminalWindowKeepRunes and
// counts every evicted rune; the first eviction flips the window's
// Truncated flag for the feed's whole life. An empty string writes nothing -
// no bytes arrived, and an empty window must not be manufactured out of
// nothing (tooth 2).
func (f *TerminalFeed) OutputChunk(text string) {
	if f == nil || text == "" {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r := []rune(text)
	f.windowTail = append(f.windowTail, r...)
	f.windowSeen += len(r)
	if over := len(f.windowTail) - terminalWindowKeepRunes; over > 0 {
		f.windowTail = f.windowTail[over:]
		f.windowDropped += over
	}
}

// Section builds the wire shape. A feed that exists always yields a non-nil
// section - the key's presence is the reader's existence - and the two rows
// appear exactly when their own sources delivered something: an input row
// only after a line arrived, an output row only after a chunk arrived.
// Absent rows are the truth "nothing arrived yet", never an empty string
// dressed as "the command printed nothing" (tooth 2).
func (f *TerminalFeed) Section() *TerminalSection {
	if f == nil {
		return nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	sec := &TerminalSection{}
	if f.inputSeen > 0 {
		sec.Input = &TerminalInputRow{
			Text:         string(f.inputTail),
			Truncated:    f.inputDropped > 0,
			DroppedRunes: f.inputDropped,
		}
	}
	if f.windowSeen > 0 {
		sec.Output = &TerminalOutputRow{
			Text:         string(f.windowTail),
			Truncated:    f.windowDropped > 0,
			DroppedRunes: f.windowDropped,
		}
	}
	return sec
}
