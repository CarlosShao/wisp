//go:build windows

package main

// Ticket 149 probe (OUT OF REPO COPY — the file that is committed under
// .scratch/wisp/probes/149/ is the same bytes; this file is never part of the
// shipped test roster of the package it sits in, it is an instrument).
//
// It asks, for shapes that classify as reportCorrupt, three different numbers
// about the same failure and prints them side by side:
//
//	z  = subjectReportRead.offset as today's code fills it (dec.InputOffset())
//	eo = the offset the decoder's OWN error carries (*json.SyntaxError.Offset /
//	     *json.UnmarshalTypeError.Offset), recovered through obs.err with
//	     errors.As - which is exactly the route a future fill would take
//
// plus the two sentences a reader actually sees (obs.err.Error(), obs.summary()).
// That is the reading AC#1 finding 2 and AC#2's fix need: which corrupt shapes
// report z=0 while the decoder had plainly begun reading a value, and whether eo
// names a position inside those bytes for them.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

func p149Fixture(t *testing.T) []byte {
	t.Helper()
	doc, err := json.MarshalIndent(&sloRun{
		Mode:       "subject-in-tree",
		State:      "Sleeping",
		Posture:    "skeleton",
		StartedAt:  "2026-09-25T00:00:00Z",
		Seconds:    5,
		Pass:       true,
		SubjectPID: 4242,
		Report: &observe.StateReport{
			State:          observe.SLOSleeping,
			StartedAt:      "2026-09-25T00:00:00Z",
			DurationSec:    5,
			IntervalSec:    0.25,
			Basis:          observe.BasisInTree,
			ObserverCost:   true,
			MemMedianBytes: 7 << 20,
			CPUMeanPercent: 0.52,
			Samples: []observe.Sample{
				{At: "2026-09-25T00:00:00Z", TreePrivateBytes: 7 << 20, CPUPercent: 0.4, Handles: 210, TCPConnections: 1},
			},
			Verdicts: []observe.Verdict{
				{Metric: "mem_private_working_set", Measured: "7864320", Limit: "<= 30 MB", Pass: true, Gate: true},
			},
		},
	}, "", "  ")
	if err != nil {
		t.Fatalf("fixture marshal: %v", err)
	}
	if len(doc) < 300 {
		t.Fatalf("fixture only %d bytes", len(doc))
	}
	return doc
}

// p149ErrOffset recovers the offset the decoder's own error names, if any.
func p149ErrOffset(err error) (string, int64) {
	var syn *json.SyntaxError
	if errors.As(err, &syn) {
		return "SyntaxError", syn.Offset
	}
	var typ *json.UnmarshalTypeError
	if errors.As(err, &typ) {
		return "UnmarshalTypeError", typ.Offset
	}
	return "none", -1
}

func TestP149CorruptShapesOffsetCensus(t *testing.T) {
	doc := p149Fixture(t)
	t.Logf("P149 FIXTURE len=%d", len(doc))

	// A prefix with one byte of the real document replaced, so the corruption
	// sits deep inside a document whose head is genuinely the subject's.
	deepBreak := func(at int, ch byte) []byte {
		b := append([]byte(nil), doc...)
		b[at] = ch
		return b
	}

	cases := []struct {
		name string
		body []byte
		// want is where the byte that the decoder objects to actually is, stated
		// by the test (not read back from the decoder) so the two can disagree.
		want int
	}{
		{"double-comma-39B", []byte(`{"mode":"subject-in-tree",,"pass":true}`), 25},
		{"html-head", []byte("<html>the runner wrote an error page</html>"), 0},
		{"wrong-type", []byte(`{"mode": 123}`), 10},
		{"empty-object", []byte(`{}`), -1},
		{"no-state-report", []byte(`{"mode":"subject-in-tree","state":"Sleeping","pass":true}`), -1},
		{"trailing-garbage", append(append([]byte(nil), doc...), []byte("\nnot part of the document")...), -1},
		{"second-document", append(append([]byte(nil), doc...), doc...), -1},
		{"trunc40-plus-0xff", append(append([]byte(nil), doc[:40]...), 0xff), 40},
		{"deep-at-100-bang", deepBreak(100, '!'), 100},
		{"deep-at-500-at", deepBreak(500, '@'), 500},
		{"deep-at-1000-hash", deepBreak(1000, '#'), 1000},
		{"deep-late-2-0xff", append(append([]byte(nil), doc[:len(doc)-2]...), 0xff, '\n'), len(doc) - 2},
		{"bad-escape-in-string", []byte(`{"mode":"\q","pass":true}`), 10},
		{"colon-then-comma", []byte(`{"mode":"subject-in-tree":,"pass":true}`), 26},
		{"array-into-struct", []byte(`[1,2,3]`), 6},
		{"number-into-struct", []byte(`123`), 2},
		{"null-doc", []byte(`null`), 3},
		{"string-doc", []byte(`"just a string"`), 14},
		{"brace-first-then-good", append([]byte{'}'}, doc...), 0},
		{"late-type-break", []byte(`{"mode":"subject-in-tree","seconds":"not a number"}`), 50},
	}

	printedNoOffsetWord := 0
	printedZero := 0
	printedPosition := 0
	rows := make([]string, 0, len(cases))
	for _, tc := range cases {
		obs := readSubjectReport(tc.body)
		kind, eo := p149ErrOffset(obs.err)
		sentence := obs.summary()
		hasOffsetWord := strings.Contains(sentence, "offset")
		switch {
		case obs.state != reportCorrupt:
			t.Logf("P149 %-22s len=%-5d state=%s  (not corrupt - classified elsewhere)", tc.name, len(tc.body), obs.state)
			continue
		case !hasOffsetWord:
			printedNoOffsetWord++
		case obs.offset == 0:
			printedZero++
		default:
			printedPosition++
		}
		line := fmt.Sprintf("P149 %-22s len=%-5d state=%-7s z(offset)=%-5d errkind=%-19s eo=%-5d want=%-5d offsetWord=%v",
			tc.name, len(tc.body), obs.state, obs.offset, kind, eo, tc.want, hasOffsetWord)
		rows = append(rows, line)
		t.Logf("%s", line)
		t.Logf("     SENT %s", sentence)
	}
	t.Logf("P149 CENSUS corrupt-ish shapes=%d prints-no-offset-word=%d prints-offset-0=%d prints-a-position=%d",
		len(cases), printedNoOffsetWord, printedZero, printedPosition)

	// Second block: the all-prefix census, to show no fixture prefix is corrupt
	// (so every corrupt reading below comes from a hand-built shape).
	states := map[string]int{}
	for i := 0; i <= len(doc); i++ {
		states[readSubjectReport(doc[:i]).state.String()]++
	}
	t.Logf("P149 ALL-PREFIX state census (0..len=%d) = %v", len(doc), states)

	// Third block: what dec.InputOffset() reports versus what the error reports,
	// for the two shapes AC#1 finding 2 names as counterexamples to the comment.
	for _, name := range []string{"double-comma-39B", "trunc40-plus-0xff"} {
		for _, tc := range cases {
			if tc.name != name {
				continue
			}
			dec := json.NewDecoder(bytes.NewReader(tc.body))
			var rep sloRun
			err := dec.Decode(&rep)
			kind, eo := p149ErrOffset(err)
			t.Logf("P149 RAWDECODE %s: err=%v InputOffset=%d %s.Offset=%d bytes=%d",
				name, err, dec.InputOffset(), kind, eo, len(tc.body))
		}
	}
}
