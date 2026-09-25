//go:build windows

package main

// PROBE (ticket 147, not a delivered test): measures the shape the ticket
// claims, with its own positive control, so nothing downstream has to trust a
// previous run's reading.
//
//	claim 2 of the ticket: on the reportUnwritten branch subjectReportRead.offset
//	carries no information (the fixture's whole prefix range gives one value),
//	while reportComplete reports the document length.
//
// It reads two things per reading so the two are never confused:
//   - obs.offset, the field the red sentence prints, and
//   - the decoder's OWN InputOffset() for the same bytes, measured here
//     independently, which is what readSubjectReport would have stored.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func tallyString(counts map[int64]int) string {
	keys := make([]int64, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(a, b int) bool { return keys[a] < keys[b] })
	out := ""
	for _, k := range keys {
		out += fmt.Sprintf("%d(%d) ", k, counts[k])
	}
	return out
}

func TestP147OffsetShapeAcrossEveryPrefix(t *testing.T) {
	doc := slo144Report(t)
	t.Logf("P147 fixture length = %d bytes (positive control: must be > 0 and non-trivial)", len(doc))

	stateCounts := map[string]int{}
	fieldOffset := map[string]map[int64]int{} // state -> obs.offset values
	rawOffset := map[string]map[int64]int{}   // state -> dec.InputOffset() of the same bytes

	for i := 0; i <= len(doc); i++ {
		obs := readSubjectReport(doc[:i])
		st := obs.state.String()
		stateCounts[st]++
		if fieldOffset[st] == nil {
			fieldOffset[st] = map[int64]int{}
			rawOffset[st] = map[int64]int{}
		}
		fieldOffset[st][obs.offset]++

		dec := json.NewDecoder(bytes.NewReader(doc[:i]))
		var rep sloRun
		_ = dec.Decode(&rep)
		rawOffset[st][dec.InputOffset()]++

		if obs.bytes != i {
			t.Errorf("prefix %d: obs.bytes = %d", i, obs.bytes)
		}
	}

	t.Logf("P147 readings over 0..len inclusive = %d, state tally: %v", len(doc)+1, stateCounts)
	for _, st := range []string{"unwritten", "corrupt", "complete"} {
		if fieldOffset[st] == nil {
			t.Logf("P147 state=%-9s never reached", st)
			continue
		}
		t.Logf("P147 state=%-9s obs.offset   distinct=%d -> %s", st, len(fieldOffset[st]), tallyString(fieldOffset[st]))
		t.Logf("P147 state=%-9s decoder own =%s", st, tallyString(rawOffset[st]))
	}

	// Positive control: the ruler reports a non-zero value somewhere, so "only
	// ever 0" cannot be a broken instrument.
	whole := readSubjectReport(doc)
	t.Logf("P147 positive control: whole document -> state=%s offset=%d bytes=%d", whole.state, whole.offset, whole.bytes)
	if whole.offset != int64(len(doc)) {
		t.Errorf("positive control broken: complete document reports offset %d, want %d", whole.offset, len(doc))
	}
	// A mid-size prefix, spelled out because that is the sentence a human reads.
	for _, i := range []int{1, 44, 100, len(doc) / 2, len(doc) - 1} {
		obs := readSubjectReport(doc[:i])
		t.Logf("P147 prefix %-5d state=%-9s bytes=%-5d obs.offset=%-5d summary=%q",
			i, obs.state, obs.bytes, obs.offset, obs.summary())
	}
}

func TestP147OffsetShapeOnTheCorruptShapes(t *testing.T) {
	doc := string(slo144Report(t))
	bodies := []string{
		"<html>the runner wrote an error page</html>",
		`{"mode":"subject-in-tree",,"pass":true}`,
		`{"mode": 123}`,
		doc + "\nnot part of the document",
		`{"mode":"subject-in-tree","state":"Sleeping","pass":true}`,
		`{}`,
		`{"mode":"subject-in-tree","report":{"state":"Sleeping","samples":[{"at":`,
	}
	for _, body := range bodies {
		obs := readSubjectReport([]byte(body))
		dec := json.NewDecoder(bytes.NewReader([]byte(body)))
		var rep sloRun
		derr := dec.Decode(&rep)
		t.Logf("P147 corrupt %-4d bytes state=%-9s obs.offset=%-5d decoder=%-5d decoderErr=%v",
			len(body), obs.state, obs.offset, dec.InputOffset(), derr)
		if obs.err != nil {
			t.Logf("    summary=%q", obs.summary())
			t.Logf("    err    =%q", obs.err.Error())
		}
	}
}

// The note-carrying observations (no file / unreadable / nothing read yet) put
// offset at -1; does that number ever reach a sentence a human reads?
func TestP147DoesMinusOneEverPrint(t *testing.T) {
	for _, note := range []string{"nothing read yet", "no report file yet", "report file unreadable: open x: The system cannot find the file specified."} {
		obs := subjectReportRead{state: reportUnwritten, bytes: 0, offset: -1, note: note}
		t.Logf("P147 note=%-40q summary=%q contains-(-1)=%v", note, obs.summary(), bytes.Contains([]byte(obs.summary()), []byte("-1")))
		if strings.Contains(obs.summary(), "-1") {
			t.Errorf("a negative offset reached a sentence: %q", obs.summary())
		}
	}
}
