package panel

// Ticket 216 AC#2 - the menu's display base: ONE pure function that washes
// control characters out of any string the panel prints that the user did not
// type themselves, applied ONCE at the display boundary instead of once per
// producer.
//
// WHY THIS TICKET EXISTS AT ALL. A file name is somebody else's string. It cannot
// change a permission verdict - ticket 200's hostile-file judgement owns that
// half - but it CAN change what a row of the menu looks like: an ESC sequence can
// blank the tail of the line a human is reading, or print a second half-sentence
// ("已允许") that no code ever decided. That is a display-side injection surface,
// and the ticket's rule is that every row the user did not type is untrusted input
// on its way to a human eye.
//
// WHY ONE WASH AT THE DISPLAY BOUNDARY, NOT ONE PER PRODUCER (the ticket's own
// warning, and the census's headline). This repository already grew two mutually
// incompatible control-character handling implementations before this file -
// cmd/wisp's sanitizeReplyText replaces every unicode.IsControl rune with a SPACE
// and collapses whitespace, internal/risk's findMetachars turns only \n \r \0 into
// visible escapes - and neither covers the menu. A third "each producer scrubs its
// own output" rule would have produced a fourth copy with a fourth character set.
// So this file holds the display rule, and it is applied where the display packet
// is finished: NewSnapshot, the one push the panel receives per state change
// (composer.go's own words). Washing here is also what keeps the functional halves
// intact: ProjectInstructionLoadRequestFor reads WorkspaceView.Canonical as a PATH
// it walks, and workspace.go resolves it; washing at those producers would hand a
// directory walker an escaped string. Here the escaped string only ever reaches a
// human.
//
// THE ROSTER THIS FILE WASHES TODAY, and why nothing else is here. AC#1's census
// (docs/evidence/s1/216-census.md 问 2) measured which menu fields are (a) not the
// user's own typing, (b) actually on the wire today, (c) unwashed:
//
//	AttachmentRef.Name        picker-supplied; the ingress guard rejects only
//	                          NUL/CR/LF/TAB, so ESC and the whole C1 range pass
//	WorkspaceView.Canonical   C26's resolved tree, not the user's spelling
//
// Plugin names, skill names and command descriptions are on AC#2's list in the
// ticket but measure ZERO SOURCE / ZERO OUTBOUND at this anchor (tickets 213/215
// are not landed), so this file invents no key and no field for them; when they
// arrive they arrive inside a Snapshot and are washed by the same call below, with
// no new copy of the rule.
//
// WHAT THIS FILE DELIBERATELY DOES NOT DO, named rather than quietly skipped:
//   - it does not merge cmd/wisp's sanitizeReplyText into here. That is the "两份
//     合一" residual ticket 216 hands back to the orchestrator: reusing this
//     function from cmd/wisp needs an EXPORTED name on package panel, which is a
//     contract-face change (AGENTS.md 1.1 / this leg's write face is internal/panel
//     only), and the two sides disagree on the OUTPUT SHAPE (space vs escape), so
//     the merge is a ruling, not an edit;
//   - it does not wash Cf (format) characters such as U+202E RTL OVERRIDE or
//     U+200B ZERO WIDTH SPACE. unicode.IsControl is category Cc only, and "控制字符"
//     is what AC#2 says. Widening to Cf is a real question - U+202E is exactly the
//     "盖掉后半段名字" attack the ticket's rationale names - and the honest answer
//     for this leg is to record it as residual for the verifying leg to rule on,
//     not to silently widen a security-relevant character set in the same commit
//     that claims to implement the narrower one;
//   - it does not touch ApprovalCardView.Tool/Args[]/Reason. Those also go out
//     verbatim today (census 问 2's 相邻面 note) but sit outside this ticket's
//     composer roster, so naming them here is all this leg does;
//   - it changes NO verdict and NO permission judgement. Everything below runs
//     after the verdicts are made, on the way to a display string.

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// scrubForDisplay turns the characters a human cannot see into characters a human
// can read, and changes nothing else.
//
// The shape is VISIBLE ESCAPE, not deletion, for one reason: deletion splices.
// "photo\x1b已允许.png" with the control character removed reads
// "photo已允许.png" - still a claim no code made, and now with the only evidence
// that it was an injected claim gone. Escaping keeps every other byte in its own
// place and prints the invisible one as `\x1b`, so the row cannot say anything the
// file name does not already say, and a reader (or a later audit of the packet's
// bytes) can see that something was there.
//
// The escapes are this function's own tiny vocabulary, not strconv.Quote, and that
// is deliberate: Quote doubles every backslash, which would render every canonical
// Windows path in the menu as C:\\Users\\..., i.e. it would make the common honest
// case unreadable to protect the rare hostile one. Instead:
//
//	\n \r \t          -> the two-character mnemonic Go itself prints
//	other C0 and C1    -> \xNN, lowercase hex, the same spelling Go uses
//	invalid UTF-8 byte -> \xNN of the raw byte, never U+FFFD
//
// The last one matters for this input class: a picker that returns truncated or
// non-UTF-8 bytes would otherwise be rendered as the replacement character, which
// the user reads as "a name with a strange glyph in it" rather than "these bytes
// are not a string at all". Showing the byte is the honest version.
//
// Idempotence holds because every escape this function emits is printable ASCII:
// washing an already-washed string returns it unchanged, which is what lets the
// same call sit at the snapshot boundary without double-escaping across pumps.
func scrubForDisplay(s string) string {
	if !needsDisplayWash(s) {
		// The honest fast path: an ordinary name or path comes back as the very
		// same string value, so "washed" never means "re-encoded".
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 16)
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == utf8.RuneError && size == 1:
			fmt.Fprintf(&b, `\x%02x`, s[i])
			i++
		case unicode.IsControl(r):
			writeControlEscape(&b, r)
			i += size
		default:
			b.WriteRune(r)
			i += size
		}
	}
	return b.String()
}

// writeControlEscape is the vocabulary above, kept apart from the scan so the
// character set has exactly one place to look at.
func writeControlEscape(b *strings.Builder, r rune) {
	switch r {
	case '\n':
		b.WriteString(`\n`)
	case '\r':
		b.WriteString(`\r`)
	case '\t':
		b.WriteString(`\t`)
	case '\x00':
		// Named explicitly because NUL is the truncation attack every later
		// C-string consumer falls for, and an audit of the packet's bytes should
		// see it said out loud rather than as an unremarkable \x00.
		b.WriteString(`\0`)
	default:
		if r <= 0xff {
			fmt.Fprintf(b, `\x%02x`, r)
			return
		}
		fmt.Fprintf(b, `\u%04x`, r)
	}
}

// needsDisplayWash reports whether scrubForDisplay would change s. It is the same
// character set as the writer, read once: unicode.IsControl (C0 U+0000-U+001F and
// C1 U+007F-U+009F) plus any byte that is not a valid UTF-8 sequence.
func needsDisplayWash(s string) bool {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if (r == utf8.RuneError && size == 1) || unicode.IsControl(r) {
			return true
		}
		i += size
	}
	return false
}

// scrubComposerForDisplay is the display boundary's whole roster, in one call. It
// returns a COPY: the attachment slice a caller handed NewComposerState is the
// broker's own array, and rewriting a row in place would silently rename an
// attachment in whatever still holds the original slice - including the message
// the agent is about to receive (OutgoingMessage reads the same rows).
func scrubComposerForDisplay(c ComposerState) ComposerState {
	c.Workspace.Canonical = scrubForDisplay(c.Workspace.Canonical)
	if len(c.Attachments) == 0 {
		return c
	}
	rows := make([]AttachmentRef, len(c.Attachments))
	for i, a := range c.Attachments {
		a.Name = scrubForDisplay(a.Name)
		rows[i] = a
	}
	c.Attachments = rows
	return c
}
