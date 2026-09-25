//go:build windows

package main

// Ticket 144 AC#1: the shape that made `wisp slo` red on CI run 36094258734 -
// the subject's report file EXISTS but its bytes had not all landed, and
// collectReport read that as "the report is corrupt" and killed the run on one
// json.Unmarshal error instead of re-looking inside its budget.
//
// These cases do not wait on a real writer racing a real reader (a race like
// that is not replayable, and this host's timings are not trustworthy anyway -
// see scripts/wisp-cli-tests.sh for why `go test ./cmd/wisp/` cannot even start
// without the sherpa DLLs on PATH). They ask the two things AC#1 asks:
//
//   - does one look at these bytes classify them as "tail still in flight" or
//     as "this contradicts a subject report"? readSubjectReport is a pure
//     function of a byte slice, so every prefix of a real report and every
//     broken shape gets asked directly (cases 1 and 2).
//   - what does the loop DO with each classification? That is asked of the real
//     loop over a real file (cases 3, 4) and of the loop with its file touch
//     scripted (cases 5, 6, 7), which is the only way to count reads and so the
//     only way to show "fails closed on the spot" without measuring seconds.
//
// The fixture bytes are produced by the same call the subject itself uses
// (json.MarshalIndent of a sloRun, see writeSLO), because the whole claim rests
// on what a PREFIX of that exact payload looks like to a decoder.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// slo144Report is the payload writeSLO would put on disk for a self-sampling
// subject: a complete sloRun with a state report in it.
//
// The fixture is built by CALLING THE SAME SERIALIZER the subject calls
// (json.MarshalIndent of a sloRun with a real observe.StateReport inside, with
// samples and verdicts, like runSubject's run.Report), and it is checked to be
// non-trivial before it is used: an empty or hand-written toy document would
// make the all-prefix scan below true for the wrong reason.
func slo144Report(t *testing.T) []byte {
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
			GDIMax:         12,
			WriteOpsTotal:  141,
			SampleErrors:   0,
			Pass:           true,
			Samples: []observe.Sample{
				{At: "2026-09-25T00:00:00Z", TreePrivateBytes: 7 << 20, CPUPercent: 0.4, Handles: 210, TCPConnections: 1},
				{At: "2026-09-25T00:00:00.250Z", TreePrivateBytes: 8 << 20, CPUPercent: 0.6, Handles: 211, TCPConnections: 2},
			},
			Verdicts: []observe.Verdict{
				{Metric: "mem_private_working_set", Measured: "7864320", Limit: "<= 30 MB", Pass: true, Gate: true},
				{Metric: "cpu_percent_all_core", Measured: "0.52", Limit: "<= 0.5", Pass: false, Gate: false, ObserverCost: true},
			},
		},
	}, "", "  ")
	if err != nil {
		t.Fatalf("fixture marshal (the same call writeSLO makes): %v", err)
	}
	// Non-triviality of the fixture, so the all-prefix scan cannot be true
	// because the document is short or one-valued.
	if len(doc) < 300 {
		t.Fatalf("fixture is only %d bytes - too short for the all-prefix claim to mean anything: %s", len(doc), doc)
	}
	for _, want := range []string{`"report"`, `"samples"`, `"verdicts"`, `"pass": true`, `"cpu_percent": 0.4`} {
		if !strings.Contains(string(doc), want) {
			t.Fatalf("fixture does not look like a real subject report (missing %s):\n%s", want, doc)
		}
	}
	return doc
}

// Case 1 - the retryable half of the split, across the WHOLE shape rather than
// a hand-picked example: every proper prefix of a real report, including the
// empty one, must classify as "tail has not arrived". One prefix that came back
// corrupt would be a report the loop kills mid-write, which is this ticket's
// bug; one prefix that came back complete would be a lie about the file.
func TestSLO144EveryPrefixOfARealReportIsUnwrittenNotCorrupt(t *testing.T) {
	doc := slo144Report(t)
	for i := 0; i < len(doc); i++ {
		obs := readSubjectReport(doc[:i])
		switch {
		case obs.state != reportUnwritten:
			t.Errorf("prefix of %d/%d bytes classified as %s, want unwritten (%s)",
				i, len(doc), obs.state, obs.summary())
		case obs.report != nil:
			t.Errorf("prefix of %d/%d bytes handed back a report: %v", i, len(doc), obs.report)
		case obs.err != nil:
			t.Errorf("prefix of %d/%d bytes returned an error though it is only unfinished: %v", i, len(doc), obs.err)
		}
		if obs.bytes != i {
			t.Errorf("prefix of %d bytes reported %d bytes read", i, obs.bytes)
		}
	}
	// And the whole document is the third state, not a retryable one.
	whole := readSubjectReport(doc)
	if whole.state != reportComplete || whole.report == nil || !whole.report.Pass {
		t.Errorf("complete report classified %s (report %v)", whole.state, whole.report)
	}
}

// Case 2 - the fail-closed half. AC#2 refuses to let the retry leg swallow
// these: none of them is a prefix of what writeSLO writes, so waiting cannot
// fix any of them and the loop must say so on the spot.
func TestSLO144ReportsThatContradictThemselvesAreCorruptNow(t *testing.T) {
	doc := string(slo144Report(t))
	cases := []struct{ name, body, wantReason string }{
		{"html-head", "<html>the runner wrote an error page</html>", "offset"},
		{"stray-comma", `{"mode":"subject-in-tree",,"pass":true}`, "contradict"},
		{"wrong-type", `{"mode": 123}`, "contradict"},
		{"trailing-garbage", doc + "\nnot part of the document", "more than one document"},
		{"second-document", doc + doc, "more than one document"},
		{"complete-but-no-state-report", `{"mode":"subject-in-tree","state":"Sleeping","pass":true}`, "no state report"},
		{"empty-object-is-a-finished-lie", `{}`, "no state report"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obs := readSubjectReport([]byte(tc.body))
			if obs.state != reportCorrupt {
				t.Fatalf("state %s (%s), want corrupt", obs.state, obs.summary())
			}
			if obs.err == nil {
				t.Fatal("corrupt without an error: the loop would retry this forever")
			}
			if obs.report != nil {
				t.Errorf("corrupt handed back a report: %v", obs.report)
			}
			if !strings.Contains(obs.err.Error(), tc.wantReason) {
				t.Errorf("error %q does not name %q", obs.err.Error(), tc.wantReason)
			}
			if !strings.Contains(obs.err.Error(), fmt.Sprintf("%d bytes", len(tc.body))) {
				t.Errorf("error %q does not name the %d bytes it read", obs.err.Error(), len(tc.body))
			}
		})
	}
}

// writeSLOFile puts bytes on disk the way the subject does: one os.WriteFile of
// the whole payload, same call, same mode.
func writeSLOFile(t *testing.T, path string, body []byte) {
	t.Helper()
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// Case 3 - AC#1(a) end to end over a real file: the report starts as an empty
// file, then a prefix, then the whole payload, and collectReport must come back
// with the parsed report instead of killing the run.
func TestSLO144LoopRetriesAnUnfinishedFileAndReadsTheWholeReport(t *testing.T) {
	doc := slo144Report(t)
	path := filepath.Join(t.TempDir(), "subject-report.json")
	writeSLOFile(t, path, nil) // the state that red on CI: present, empty

	written := make(chan error, 1)
	go func() {
		// The writer half of the race, replayed by hand: empty file, then half
		// the payload, then all of it. Nothing here asserts WHEN the loop looks;
		// case 5 and case 7 are where the reading count is pinned.
		defer func() { written <- nil }()
		time.Sleep(10 * time.Millisecond)
		if err := os.WriteFile(path, doc[:len(doc)/2], 0o644); err != nil {
			written <- err
			return
		}
		time.Sleep(10 * time.Millisecond)
		written <- os.WriteFile(path, doc, 0o644)
	}()

	s := &sloSubject{pid: 4242, outPath: path}
	rep, err := s.collectReportWithin(20*time.Second, time.Millisecond)
	if werr := <-written; werr != nil {
		t.Fatalf("fixture writer: %v", werr)
	}
	if err != nil {
		t.Fatalf("collectReport on a file that finishes writing: %v", err)
	}
	if rep == nil || rep.Report == nil {
		t.Fatalf("no report came back: %+v", rep)
	}
	if rep.Mode != "subject-in-tree" || !rep.Pass {
		t.Errorf("parsed the wrong document: mode=%q pass=%v", rep.Mode, rep.Pass)
	}
}

// Case 4 - the two give-up sentences, over a real file.
//
// 4a: a prefix whose tail never arrives (the CI shape, frozen at half a
// document). The run is still red - fail-closed is not relaxed - but the
// sentence now names the budget it spent and the byte count it stopped at,
// which is what AC#1(a) asks for instead of "unexpected end of JSON input".
// 4b: bytes that contradict a subject report, and must be red on the spot.
func TestSLO144LoopGiveUpSentencesOnRealFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "subject-report.json")

	head := `{"mode":"subject-in-tree","report":{"state":`
	writeSLOFile(t, path, []byte(head))
	s := &sloSubject{pid: 4243, outPath: path}
	_, err := s.collectReportWithin(50*time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("a report that never finished writing passed the loop")
	}
	msg := err.Error()
	for _, want := range []string{"within 50ms", fmt.Sprintf("%d bytes", len(head)), "tail had not arrived"} {
		if !strings.Contains(msg, want) {
			t.Errorf("unfinished give-up %q does not name %q", msg, want)
		}
	}
	if strings.Contains(msg, "unexpected end of JSON input") {
		t.Errorf("the give-up sentence is still the bare decoder complaint: %q", msg)
	}

	notADoc := []byte("<html>not a report</html>")
	writeSLOFile(t, path, notADoc)
	s2 := &sloSubject{pid: 4244, outPath: path}
	_, err = s2.collectReportWithin(20*time.Second, time.Millisecond)
	if err == nil {
		t.Fatal("corrupt report bytes passed the loop")
	}
	for _, want := range []string{"corrupt", fmt.Sprintf("%d bytes", len(notADoc))} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("corrupt give-up %q does not name %q", err.Error(), want)
		}
	}
	if strings.Contains(err.Error(), "within 20s") {
		t.Errorf("corrupt bytes were treated as a wait, not a verdict: %q", err.Error())
	}
}

// scriptedReader hands the loop a fixed sequence of file contents and counts how
// many times it looked - the credential behind "on the spot". With repeatLast it
// keeps serving the final scripted reading (a poll waiting for bytes that never
// arrive); without it a reading past the end of the script is a failure, because
// the loop looked at something it should already have judged.
type scriptedReader struct {
	t          *testing.T
	got        []func() ([]byte, error)
	repeatLast bool
	calls      int
}

func (r *scriptedReader) read(path string) ([]byte, error) {
	r.calls++
	if r.calls > len(r.got) {
		if !r.repeatLast {
			r.t.Fatalf("%s: read #%d of a %d-reading script - the loop is retrying something it should have judged on the first look",
				path, r.calls, len(r.got))
		}
		return r.got[len(r.got)-1]()
	}
	return r.got[r.calls-1]()
}

func scriptBytes(body []byte) func() ([]byte, error) {
	return func() ([]byte, error) { return body, nil }
}

func scriptMissing() ([]byte, error) {
	return nil, os.ErrNotExist
}

// Case 5 - the load-bearing count: corrupt bytes are read EXACTLY ONCE, so the
// run is red because the report contradicts itself and not because a budget ran
// out. AC#2's "不许改成解析失败就 continue 到预算尽头" is what this nails: a loop
// that retried would reach past the end of the one-reading script and be
// reported by the reader itself.
func TestSLO144CorruptReportIsJudgedOnTheFirstRead(t *testing.T) {
	body := `{"mode":"subject-in-tree","pass":true}` // complete, no state report
	r := &scriptedReader{t: t, got: []func() ([]byte, error){scriptBytes([]byte(body))}}
	s := &sloSubject{pid: 4245, outPath: "scripted.json", readReportFile: r.read}
	_, err := s.collectReportWithin(30*time.Second, time.Millisecond)
	if err == nil {
		t.Fatal("a subject report with no state report passed")
	}
	if !strings.Contains(err.Error(), "corrupt") || !strings.Contains(err.Error(), "no state report") {
		t.Errorf("error %q does not name the corruption", err.Error())
	}
	if r.calls != 1 {
		t.Errorf("corrupt report was read %d times, want exactly 1 (fail closed on the spot)", r.calls)
	}
}

// Case 6 - AC#1(a)'s other half: a prefix keeps the loop polling, and when it
// gives up the sentence names the budget AND the last byte count.
func TestSLO144UnfinishedReportKeepsPollingThenNamesBudgetAndBytes(t *testing.T) {
	head := []byte(`{"mode":`) // 8 bytes of a document whose tail never arrives
	r := &scriptedReader{t: t, repeatLast: true, got: []func() ([]byte, error){
		scriptMissing, scriptMissing, scriptBytes(head),
	}}
	s := &sloSubject{pid: 4246, outPath: "scripted.json", readReportFile: r.read}
	_, err := s.collectReportWithin(60*time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("a report that never grew past a prefix passed the loop")
	}
	if r.calls <= 3 {
		t.Errorf("the loop took %d reads; a missing file and an unfinished file must both be re-read inside the budget", r.calls)
	}
	msg := err.Error()
	for _, want := range []string{"within 60ms", "8 bytes", "tail had not arrived"} {
		if !strings.Contains(msg, want) {
			t.Errorf("give-up error %q does not name %q", msg, want)
		}
	}
}

// Case 7 - the same loop, and the payload that arrives late is read: proves the
// polling ends in a parsed report, not just in a longer error.
func TestSLO144ReportThatArrivesAfterMissingReadingsIsCollected(t *testing.T) {
	doc := slo144Report(t)
	r := &scriptedReader{t: t, got: []func() ([]byte, error){
		scriptMissing, scriptBytes(nil), scriptBytes(doc[:3]), scriptBytes(doc),
	}}
	s := &sloSubject{pid: 4247, outPath: "scripted.json", readReportFile: r.read}
	rep, err := s.collectReportWithin(5*time.Second, time.Millisecond)
	if err != nil {
		t.Fatalf("collectReportWithin: %v", err)
	}
	if r.calls != 4 {
		t.Errorf("the loop took %d readings, want the scripted 4", r.calls)
	}
	if rep == nil || rep.Report == nil || rep.Report.State != observe.SLOSleeping {
		t.Fatalf("collected the wrong report: %+v", rep)
	}
}

// ---------------------------------------------------------------------------
// Ticket 147, cases 8-10. Same fixture, same helpers, one different subject:
// the NUMBER inside the unwritten sentence rather than the retry it describes.
//
// subjectReportRead.offset on the reportUnwritten branch used to be the
// decoder's own InputOffset(), which a probe measured over the fixture's whole
// prefix range: 1978 readings, one distinct value, 0. The give-up sentence
// printed "1949 bytes read, document still open at offset 0" - a byte position
// that was not a byte position, and one that points a reader at "nothing was
// written" when the truth is "all of it but the tail is there". AC#1 resolves
// this to (a), the field carries the information, so the sentence's two numbers
// each answer one question: how many bytes landed, and where inside them the
// document stopped.
//
// These three are the teeth ticket 144's acceptance run measured as missing
// then: deleting "at offset %d" from that sentence, and putting the raw decoder
// value back, both turn them red (each measured, transcribed in
// docs/evidence/s1/147-offset-naming-r1.md). Nothing here waits on wall-clock
// agreement about the budget, which is ticket 144's ground; the only new
// assertion is about an offset and the bytes that sentence was made from.
// ---------------------------------------------------------------------------

// slo147UnwrittenRe captures BOTH numbers of the unwritten sentence. Case 8
// requires it to span the whole string, so the sentence cannot lose either half
// or grow a tail; case 9 requires only a match, because there the sentence sits
// inside the loop's longer give-up error.
var slo147UnwrittenRe = regexp.MustCompile(`([0-9]+) bytes read, document still open at offset ([0-9]+): the tail had not arrived`)

// slo147PairOf reads the two numbers an unwritten sentence promises, and says
// whether they are the whole sentence or a fragment inside a longer one. It is
// a parser and nothing else: no expectation of its own, so an expectation that
// fails has somewhere to be wrong.
func slo147PairOf(sentence string) (bytesRead, offset int, whole bool, ok bool) {
	loc := slo147UnwrittenRe.FindStringSubmatchIndex(sentence)
	if loc == nil {
		return 0, 0, false, false
	}
	b, berr := strconv.Atoi(sentence[loc[2]:loc[3]])
	o, oerr := strconv.Atoi(sentence[loc[4]:loc[5]])
	if berr != nil || oerr != nil {
		return 0, 0, false, false
	}
	return b, o, loc[0] == 0 && loc[1] == len(sentence), true
}

// Case 8 - every prefix of the real fixture, aggregated into one report so the
// claim is "all of them", not "the three I picked". For a truncated document
// the position the decoder stopped at IS the last byte on disk, so the two
// numbers must agree with each other and with the length handed to the
// classifier.
func TestSLO147UnwrittenSentenceNamesTheOffsetTheDocumentStoppedAt(t *testing.T) {
	doc := slo144Report(t)
	prefixes := 0
	wrong := 0
	first := ""
	for i := 0; i <= len(doc); i++ {
		obs := readSubjectReport(doc[:i])
		if obs.state != reportUnwritten {
			continue // case 1 owns the classification; the whole document is case 10's
		}
		prefixes++
		sentence := obs.summary()
		b, o, whole, ok := slo147PairOf(sentence)
		switch {
		case obs.offset != int64(i):
			// The field leg: printing the byte count twice would keep the
			// sentence true while the value it is made from rots back to the
			// constant this ticket is about.
			wrong++
			if first == "" {
				first = fmt.Sprintf("prefix of %d bytes: subjectReportRead.offset is %d", i, obs.offset)
			}
		case !ok:
			wrong++
			if first == "" {
				first = fmt.Sprintf("prefix of %d bytes: sentence carries no byte/offset pair: %q", i, sentence)
			}
		case !whole:
			wrong++
			if first == "" {
				first = fmt.Sprintf("prefix of %d bytes: %q is more or less than the unwritten sentence", i, sentence)
			}
		case b != i || o != i:
			wrong++
			if first == "" {
				first = fmt.Sprintf("prefix of %d bytes: sentence says %d bytes and offset %d", i, b, o)
			}
		}
	}
	if prefixes < len(doc) {
		t.Fatalf("only %d of %d readings were unwritten - case 1's denominator is not what this asserts over",
			prefixes, len(doc))
	}
	if wrong > 0 {
		t.Errorf("%d of %d unwritten sentences do not name the offset the document stopped at; first: %s",
			wrong, prefixes, first)
	}
}

// Case 9 - the same two numbers as a human meets them: the give-up error the
// real loop returns, not the classifier's return value. The byte count is the
// one the scripted reader handed the loop, so "the offset in that sentence
// equals the bytes that sentence counted" is checked on the production path.
func TestSLO147LoopGiveUpSentenceAgreesWithTheBytesItRead(t *testing.T) {
	head := []byte(`{"mode":`) // 8 bytes whose tail never arrives
	r := &scriptedReader{t: t, repeatLast: true, got: []func() ([]byte, error){
		scriptMissing, scriptBytes(head),
	}}
	s := &sloSubject{pid: 4248, outPath: "scripted.json", readReportFile: r.read}
	_, err := s.collectReportWithin(60*time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("a report that never grew past a prefix passed the loop")
	}
	b, o, _, ok := slo147PairOf(err.Error())
	if !ok {
		t.Fatalf("give-up error %q does not read as bytes-then-offset", err.Error())
	}
	if b != len(head) || o != len(head) {
		t.Errorf("give-up error says %d bytes read, document open at offset %d; want both %d (the prefix the loop was handed)",
			b, o, len(head))
	}
}

// Case 10 - the field's three readings have to stay tellable apart from the
// sentences alone: -1 means no decoder ran at all and must not print a position
// it does not have, N on an unfinished document is where it stopped inside
// exactly N bytes, N on a complete one is the end of the document. Three
// different shapes, three different sentences; a reader should never have to
// guess which of the three a bare number is.
func TestSLO147OffsetSemanticsRenderThreeDifferentSentences(t *testing.T) {
	doc := slo144Report(t)

	// -1: the loop never got bytes to classify (file absent the whole budget).
	missing := &scriptedReader{t: t, repeatLast: true, got: []func() ([]byte, error){scriptMissing}}
	_, err := (&sloSubject{pid: 4249, outPath: "scripted.json", readReportFile: missing.read}).
		collectReportWithin(40*time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("a report file that never appeared passed the loop")
	}
	if strings.Contains(err.Error(), "offset") {
		t.Errorf("the no-decoder-ran give-up prints a byte position it does not have: %q", err.Error())
	}

	// N on an unfinished document, taken from the middle of the prefix range so
	// it cannot be the constant 0 the old shape reported.
	mid := len(doc) / 2
	unfinished := readSubjectReport(doc[:mid]).summary()
	if b, o, _, ok := slo147PairOf(unfinished); !ok {
		t.Errorf("unfinished sentence %q is not the unwritten shape", unfinished)
	} else if b != mid || o != mid {
		t.Errorf("unfinished sentence says %d bytes, offset %d; want both %d", b, o, mid)
	}

	// N on the complete document, which is a different sentence about the same
	// two numbers.
	whole := readSubjectReport(doc).summary()
	if !strings.Contains(whole, fmt.Sprintf("%d bytes read, complete document at offset %d", len(doc), len(doc))) {
		t.Errorf("complete sentence %q does not name %d bytes and offset %d", whole, len(doc), len(doc))
	}

	for _, pair := range [][2]string{{err.Error(), unfinished}, {err.Error(), whole}, {unfinished, whole}} {
		if pair[0] == pair[1] {
			t.Errorf("two of the three offset meanings render the same sentence: %q", pair[0])
		}
	}
}
