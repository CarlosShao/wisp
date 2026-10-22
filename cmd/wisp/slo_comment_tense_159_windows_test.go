//go:build windows

package main

// Ticket 159 AC#2 (C2): the ruler that keeps a "the next reader will find X"
// promise from silently outliving X. The defect it is aimed at is not a wrong
// number - ticket 156's acceptance run re-counted all six numbers in case 14's
// comment block and found each one still true - it is a sentence whose TENSE says
// the defect is live after the commit that killed it. 1d53526c made :739 true and,
// in the same edit, rewrote one paragraph from future to past; the three sentences
// around it kept saying "is", "at this anchor", "the bytes this file is in today",
// and nothing in this repository could notice.
//
// WHY A TEST AND NOT A GREP (the two questions the ticket forces before a ruler is
// allowed to exist: what is its smallest visible unit, and does the evidence still
// survive being split across two of them)
//   The smallest visible unit here is ONE (pointer, fragment) pair. The earlier
//   shape of this case would have been "grep the file for 'at this anchor' and
//   fail" - a unit of one token, which goes red the day somebody writes the phrase
//   into a sentence that is CORRECT, and stays green the day somebody keeps the
//   phrase and deletes the measurement it referred to. The problem with the deictic
//   is not the words, it is that the words and their referent can drift apart. So
//   each unit asserts the PAIR: the line the historicising block claims still
//   carries fragment F, byte for byte, AND the block still makes that claim. Split
//   the evidence across the two and both halves are still readable; drop either
//   half and this case names which one went missing.
//
// WHY "UNIQUE" DOES NOT MEAN "APPEARS ONCE IN THE FILE"
//   A correction that quotes the sentence it corrects puts that sentence on disk
//   twice, and it must: a paraphrase would let the correction drift away from its
//   subject without anything noticing. So the rule this case enforces is narrower
//   and stronger than a count: the fragment may appear on exactly one line OUTSIDE
//   the correction block, and any occurrence inside the block must be pill-delimited
//   (|...|), which is this file's visible mark for "this is a quotation of a
//   sentence that has been historicised, not a live claim". A second unmarked copy
//   is what a reader would take for a live sentence, and that is the failure.
//
// WHAT "TURNS RED ON THE NEXT DRIFT" MEANS HERE, stated so it can be checked
//   (a) edit a historicised line so the quoted fragment no longer sits on the named
//       line  -> fragment-not-on-its-line
//   (b) delete a "159 (n) :NNN - |fragment|" row from the block -> claim-missing
//   (c) restore a historicised sentence to the present tense while the correction
//       still calls it past -> correction-lost-its-tense
//   (d) insert lines above the Case 13 block without re-measuring the pointer table
//       -> pointer-row-target-mismatch
//   (e) edit one of the six numbers in the block without re-counting the logs
//       -> number-disagrees-with-its-log
//   (f) copy a stale sentence to a second unmarked line -> fragment-has-a-second-home
//   An assertion that could not go red was not written down.

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// slo159Self is the file the pointers live in, relative to cmd/wisp.
const slo159Self = "slo_report_144_windows_test.go"

// slo159BlockStart marks the first line of the correction block, and
// slo159BlockEnd the first line after it.
const (
	slo159BlockStart = "Ticket 159 AC#2 (C2), appended by the record-level cleanup leg"
	slo159BlockEnd   = "Who was misled, and by which sentence, until this commit"
)

// slo159Claim is one row of the historicising block's pointer claim: the line
// numbered want in slo159Self must still carry fragment. fragment is quoted the way
// the block quotes it - a byte-equal tail of the line, without the leading "// " -
// so the block and the file cannot agree by paraphrase.
type slo159Claim struct {
	want     int
	fragment string
}

// slo159Claims is the (pointer, fragment) set the case 14 block asserts. The line
// numbers are the ones on this commit; when a later edit moves them, both the block
// and this table have to be updated in one commit, and the failure names both.
var slo159Claims = []slo159Claim{
	{742, "Correction of that triple, measured at this anchor instead of copied"},
	{749, `block names says 20/13/0 as well. So "=== RUN=21, PASS=14, FAIL=0" is neither`},
	{750, "reading: it is this selector's UNMUTATED baseline, 21/14/0, pasted into the"},
	{752, "21/14/0 again at this anchor"},
}

// slo159Numbers is the second half: every reading the block states must still equal
// what its named log contains, re-counted on the spot. This is what stops a future
// "the comment is stale" repair from pasting in numbers nobody measured.
var slo159Numbers = []struct {
	what string
	log  string
	run  int
	pass int
	fail int
}{
	{"case-off (the mutation with case 14 renamed away)", ".scratch/wisp/probes/152/mut-post/g1-restored-borrowed-value-case14off.log", 20, 13, 0},
	{"case-live (the same mutation with case 14 in the roster)", ".scratch/wisp/probes/152/mut-post/g1-restored-borrowed-value.log", 21, 13, 1},
	{"unmutated baseline (what the stale sentence had pasted)", ".scratch/wisp/probes/152/mut-post/asis.log", 21, 14, 0},
	{"case-off, re-measured by r4's own run", ".scratch/wisp/probes/156/mut-156-r2/r2-152g1-case14off.log", 20, 13, 0},
	{"case-live, re-measured by r4's own run", ".scratch/wisp/probes/156/mut-156-r2/r2-152g1-case14live.log", 21, 13, 1},
	{"unmutated baseline, re-measured by r4's own run", ".scratch/wisp/probes/156/mut-156-r2/r2-asis-152-surface.log", 21, 14, 0},
}

// slo159PtrRow is one row of the pointer table the Case 13 block carries: a pointer
// whose number was written before this ticket's insertion must still be recoverable,
// and at its new number must carry the named content.
type slo159PtrRow struct {
	what   string
	said   int
	saysAt int
	token  string
}

var slo159PtrRows = []slo159PtrRow{
	{"Case 13 header", 835, 872, "Case 13 - AC#3's ruling is"},
	{"the elided evidence path", 851, 888, "docs/evidence/s1/152-...-accept-r1.md"},
	{"the \"which letter\" block", 849, 886, "Which"},
	{"the same block's last line", 857, 894, "Neither ruling moved one assertion"},
}

// slo159HeaderMD5 is the hash the Case 13 block publishes for its own first line,
// together with the pointer row it belongs to. The hash is of the line INCLUDING its
// newline, which is how the repo's own reading was taken (sed -n '835p' | md5sum).
const (
	slo159HeaderMD5       = "a756e044366a5c957bef470027a86bc0"
	slo159HeaderMD5SaidAt = 835
	slo159HeaderMD5SaysAt = 872
)

// slo159Lines reads the file this test lives in, as lines. A read failure is fatal,
// not a skip: the file is the subject, and "the subject is missing" is not a pass.
func slo159Lines(t *testing.T) []string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("instrument: runtime.Caller cannot name this file, so the subject cannot be located")
	}
	path := filepath.Join(filepath.Dir(thisFile), slo159Self)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("instrument: reading %s: %v", path, err)
	}
	if len(raw) == 0 {
		t.Fatalf("instrument: %s is empty, so every assertion below would be vacuous", path)
	}
	return strings.Split(string(raw), "\n")
}

// slo159Block returns the half-open [start, end) line-index range of the correction
// block, or (-1, -1) when either marker is gone.
func slo159Block(lines []string) (int, int) {
	start, end := -1, -1
	for i, line := range lines {
		if start < 0 && strings.Contains(line, slo159BlockStart) {
			start = i
			continue
		}
		if start >= 0 && strings.Contains(line, slo159BlockEnd) {
			end = i
			break
		}
	}
	return start, end
}

// slo159RepoRoot walks up from this file to the directory holding go.mod.
func slo159RepoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("instrument: runtime.Caller cannot name this file, so the logs cannot be located")
	}
	dir := filepath.Dir(thisFile)
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatalf("instrument: no go.mod above %s", thisFile)
	return ""
}

// TestSLO159Case14CommentKeepsItsTenseAndItsNumbers is the anti-drift ruler AC#2
// asks for. One subtest per question, each named so a red run says which half of
// the evidence went missing.
func TestSLO159Case14CommentKeepsItsTenseAndItsNumbers(t *testing.T) {
	lines := slo159Lines(t)
	blockStart, blockEnd := slo159Block(lines)
	at := func(n int) (string, bool) {
		if n < 1 || n > len(lines) {
			return "", false
		}
		return lines[n-1], true
	}

	t.Run("the-block-itself-is-still-here", func(t *testing.T) {
		if blockStart < 0 {
			t.Fatalf("claim-missing: the historicising block is gone - %q is nowhere in %s. The sentences it corrects are back to reading as live claims, which is exactly the state this ticket closed.", slo159BlockStart, slo159Self)
		}
		if blockEnd < 0 || blockEnd <= blockStart {
			t.Fatalf("claim-missing: the correction block starts at line %d but never ends with %q, so this case cannot tell a quotation from a live sentence.", blockStart+1, slo159BlockEnd)
		}
	})

	t.Run("the-corrected-sentences-are-still-on-their-lines", func(t *testing.T) {
		for _, c := range slo159Claims {
			line, ok := at(c.want)
			if !ok {
				t.Errorf("fragment-not-on-its-line: the block still names :%d, but this file has only %d lines. Either the block or this table is stale; both must move in one commit.", c.want, len(lines))
				continue
			}
			if !strings.Contains(line, c.fragment) {
				t.Errorf("fragment-not-on-its-line: :%d no longer carries %q.\n  the line reads: %q\n  The block above still historicises it by that fragment, so either the sentence was edited (a reading taken at a line keeps its words) or the block must be re-pinned. This is the drift AC#2 exists to catch: a correction whose subject moved out from under it is a correction of nothing.", c.want, c.fragment, line)
			}
		}
	})

	t.Run("the-correction-sits-after-what-it-corrects", func(t *testing.T) {
		if blockStart < 0 {
			t.Skip("the block-is-still-here subtest already failed; nothing to order")
		}
		for _, c := range slo159Claims {
			if c.want > blockStart {
				t.Errorf("claim-out-of-order: the block starts at line %d and claims to historicise :%d, which is at or below it. The correction must sit AFTER the sentences it corrects, so a reader meets the live sentence and its correction in that order.", blockStart+1, c.want)
			}
		}
	})

	t.Run("the-block-quotes-its-subjects-and-does-not-paraphrase-them", func(t *testing.T) {
		if blockStart < 0 {
			t.Skip("the block-is-still-here subtest already failed; nothing to compare")
		}
		block := strings.Join(lines[blockStart:blockEnd], "\n")
		for _, c := range slo159Claims {
			if !strings.Contains(block, "|"+c.fragment+"|") {
				t.Errorf("claim-missing: the block no longer quotes :%d with its pill-delimited fragment %q. A correction that paraphrases instead of quoting can drift from its subject without anything noticing - which is the defect this whole ticket is about.", c.want, c.fragment)
			}
		}
	})

	t.Run("a-stale-sentence-has-exactly-one-unmarked-home", func(t *testing.T) {
		if blockStart < 0 {
			t.Skip("the block-is-still-here subtest already failed; nothing to scope")
		}
		for _, c := range slo159Claims {
			var unmarked []int
			for i, line := range lines {
				if !strings.Contains(line, c.fragment) {
					continue
				}
				if i >= blockStart && i < blockEnd && strings.Contains(line, "|"+c.fragment+"|") {
					continue // a marked quotation, which is what the block is for
				}
				unmarked = append(unmarked, i+1)
			}
			switch len(unmarked) {
			case 1:
				if unmarked[0] != c.want {
					t.Errorf("fragment-has-a-second-home: %q now sits unmarked at :%d, and the block says :%d. Two unmarked homes for one stale sentence is how a reader ends up reading the copy nobody corrects.", c.fragment, unmarked[0], c.want)
				}
			case 0:
				t.Errorf("fragment-not-on-its-line: %q is nowhere outside the correction block, so the pointer at :%d designates nothing", c.fragment, c.want)
			default:
				t.Errorf("fragment-has-a-second-home: %q sits unmarked on lines %v, so :%d no longer designates one line. A second unmarked copy lets a stale sentence hide behind the correction.", c.fragment, unmarked, c.want)
			}
		}
	})

	t.Run("the-numbers-still-match-their-logs", func(t *testing.T) {
		root := slo159RepoRoot(t)
		for _, n := range slo159Numbers {
			raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(n.log)))
			if err != nil {
				t.Errorf("instrument: reading the tracked log %s: %v. If a checkout does not carry it, this case must not report a green it did not measure - that is what it is asserting about other people's numbers.", n.log, err)
				continue
			}
			body := strings.Split(string(raw), "\n")
			count := func(prefix string) int {
				c := 0
				for _, line := range body {
					if strings.HasPrefix(line, prefix) {
						c++
					}
				}
				return c
			}
			got := [3]int{count("=== RUN"), count("--- PASS:"), count("--- FAIL:")}
			want := [3]int{n.run, n.pass, n.fail}
			if got != want {
				t.Errorf("number-disagrees-with-its-log: the case 14 comment states %d/%d/%d for %s, and %s re-counts %d/%d/%d. Either the comment was edited without a re-measure, or the log changed under it. Both are the drift this case exists for.",
					n.run, n.pass, n.fail, n.what, n.log, got[0], got[1], got[2])
			}
		}
	})

	t.Run("the-block-states-the-numbers-its-logs-carry", func(t *testing.T) {
		// The subtest above proves the logs still say 20/13/0 and 21/13/1. This one
		// proves the BLOCK says so too - otherwise a block that stopped naming any
		// number would pass both.
		joined := strings.Join(lines, "\n")
		for _, want := range []string{
			"case-off 20/13/0",
			"case-live 21/13/1",
			"the unmutated baseline 21/14/0 twice",
		} {
			if !strings.Contains(joined, want) {
				t.Errorf("claim-missing: the historicising block no longer states %q. The numbers are the half of the paragraph above that did NOT expire - a reader who follows this correction must still be told which readings stand.", want)
			}
		}
	})

	t.Run("the-correction-still-says-these-sentences-are-past", func(t *testing.T) {
		if blockStart < 0 {
			t.Skip("the block-is-still-here subtest already failed; nothing to read")
		}
		block := strings.Join(lines[blockStart:blockEnd], "\n")
		for _, want := range []string{
			`"This anchor" was the head at the moment of writing`,
			"went stale the moment 1d53526c put the re-measured numbers into",
			"Those two lines were true when written and are false now",
		} {
			if !strings.Contains(block, want) {
				t.Errorf("correction-lost-its-tense: the historicising block no longer says %q, so the three sentences above read as live claims again. This is the exact drift AC#2 named: r4 made :739 true and the sentences around it kept the present tense.", want)
			}
		}
	})

	t.Run("the-pointer-table-still-designates-its-lines", func(t *testing.T) {
		joined := strings.Join(lines, "\n")
		for _, row := range slo159PtrRows {
			line, ok := at(row.saysAt)
			if !ok {
				t.Errorf("pointer-row-target-mismatch: the table says %q moved to :%d, but this file has only %d lines.", row.what, row.saysAt, len(lines))
				continue
			}
			if !strings.Contains(line, row.token) {
				t.Errorf("pointer-row-target-mismatch: the table says %q (said :%d) now sits at :%d, and that line reads %q - it does not carry %q. An insertion above the Case 13 block moves every one of these at once, which is why the table exists; re-measure it in the same commit that moves them.", row.what, row.said, row.saysAt, line, row.token)
			}
			if !strings.Contains(joined, "159 ptr said |"+row.what+"|"+strconv.Itoa(row.said)+"|"+strconv.Itoa(row.saysAt)+"|") {
				t.Errorf("claim-missing: the pointer table no longer carries the row %q (%d -> %d), so nothing in this file says where that pointer went.", row.what, row.said, row.saysAt)
			}
		}
	})

	t.Run("the-published-md5-still-hashes-the-line-it-names", func(t *testing.T) {
		// The Case 13 block publishes a hash as proof that the header line was
		// restored byte for byte. A hash nobody recomputes is decoration, and it is
		// exactly the kind of number that survives an edit that moves the line. The
		// reading is taken the way the repo's own reading was: line INCLUDING its
		// newline, which is what `sed -n '835p' | md5sum` feeds.
		line, ok := at(slo159HeaderMD5SaysAt)
		if !ok {
			t.Errorf("pointer-row-target-mismatch: the block says its own header now sits at :%d, and this file has only %d lines.", slo159HeaderMD5SaysAt, len(lines))
			return
		}
		sum := md5.Sum([]byte(line + "\n"))
		if got := hex.EncodeToString(sum[:]); got != slo159HeaderMD5 {
			t.Errorf("md5-disagrees-with-its-line: the block publishes %s for its own header line (said :%d, now :%d), and that line hashes to %s. Either the header was edited - which would make the block's 'nothing was rewritten' claim false - or the pointer moved to a different line. The line reads %q.",
				slo159HeaderMD5, slo159HeaderMD5SaidAt, slo159HeaderMD5SaysAt, got, line)
		}
	})

	t.Run("the-old-numbers-record-a-real-move", func(t *testing.T) {
		// A pointer table whose every "said" equals its "says" would be decoration:
		// the correction block sits above the Case 13 block, so those pointers must
		// have moved. This fails if somebody rewrites the table by copying the new
		// numbers into both columns, which would erase the fact that they moved.
		if len(slo159PtrRows) == 0 {
			t.Fatalf("instrument: no pointer rows, so this subtest would be vacuous")
		}
		unchanged := 0
		for _, row := range slo159PtrRows {
			if row.said == row.saysAt {
				unchanged++
			}
		}
		if unchanged == len(slo159PtrRows) {
			t.Errorf("table-is-decoration: every row says its pointer did not move, yet the correction block sits above the Case 13 block. Copying the new number into the 'said' column erases the record of what the reader used to see.")
		}
	})

	t.Run("the-four-historicised-lines-are-four-lines", func(t *testing.T) {
		// Append-only: the correction quotes four distinct lines. If a later edit
		// merged two of them, every fragment above could still match while the
		// block's own count of what it corrects became wrong.
		seen := map[int]bool{}
		for _, c := range slo159Claims {
			if seen[c.want] {
				t.Errorf("instrument: two claims name the same line :%d, so this table cannot tell four sentences from three", c.want)
			}
			seen[c.want] = true
		}
		if got := len(seen); got != 4 {
			t.Errorf("claim-count: this table holds %d distinct lines, and the correction historicises four (:742, :749, :750, :752)", got)
		}
	})
}
