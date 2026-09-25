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
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// slo144Report is the payload writeSLO would put on disk for a self-sampling
// subject: a complete sloRun with a state report in it.
func slo144Report(t *testing.T) []byte {
	t.Helper()
	doc, err := json.MarshalIndent(&sloRun{
		Mode:      "subject-in-tree",
		State:     "Sleeping",
		Posture:   "skeleton",
		StartedAt: "2026-09-25T00:00:00Z",
		Seconds:   5,
		Pass:      true,
		Report: &observe.StateReport{
			State:       observe.SLOSleeping,
			StartedAt:   "2026-09-25T00:00:00Z",
			DurationSec: 5,
			IntervalSec: 0.25,
			Basis:       observe.BasisInTree,
		},
	}, "", "  ")
	if err != nil {
		t.Fatalf("fixture marshal (the same call writeSLO makes): %v", err)
	}
	if !strings.Contains(string(doc), `"report"`) {
		t.Fatalf("fixture has no report field, so it proves nothing: %s", doc)
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
