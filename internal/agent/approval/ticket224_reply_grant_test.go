package approval_test

// Ticket 224-r2 cell N#1 - the permanent judge for "the answer writes the row".
//
// Why these cases exist at all: the adversarial table
// docs/evidence/s1/224-session-grant-v1.md §2 ran a targeted mutation (M4) that
// deletes the Record call inside Gate.allowSession while KEEPING the
// "GRANT-RECORDED" audit line. Nothing went red in any of the four packages: the
// row the user paid for with a click was never written, and the log claimed it
// was. That is exactly the failure internal/tools/grant_test.go:14-17 names -
// "the failure this feature can have is not 'it broke' but 'it silently grants
// nothing and everything still asks'" - which reads like a green suite unless
// something asserts the opposite shape. Before this file, ticket 224's two
// commits had landed ZERO test files in internal/agent/approval, so the whole
// answer-to-row hop had no instrument on it.
//
// Two claims, kept in separate assertions so a mutation cannot satisfy one by
// faking the other:
//
//	RECORD  - one answered card puts one Record call per path the card itself
//		  printed onto the recorder, carrying the card's own tool name.
//	HONESTY - the grant_id in the audit line is the id the recorder actually
//		  returned, so "GRANT-RECORDED grant_id=0" over a recorder that saw
//		  nothing cannot pass either.
//
// The controls run in the fail-closed direction, named because M4 is the shape
// to be afraid of: a forged answer records nothing, a host with no recorder says
// out loud that the scope was dropped, and a recorder that errors still gets its
// call released (a storage fault does not revoke the user's answer) while the
// line for the row it did not write stays unsaid.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// recCall is one Record call as the recorder saw it.
type recCall struct{ tool, pattern string }

// recGrantRecorder is the GrantRecorder the gate is held to. It hands back ids
// starting at 11 so a grant_id=0 in an audit line is unmistakable, and failOn
// makes one pattern's write blow up (the partial-failure branch).
type recGrantRecorder struct {
	mu     sync.Mutex
	calls  []recCall
	next   int64
	failOn map[string]error
}

func newRecRecorder() *recGrantRecorder {
	return &recGrantRecorder{next: 10, failOn: map[string]error{}}
}

func (r *recGrantRecorder) Record(_ context.Context, tool, pattern string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err, bad := r.failOn[pattern]; bad {
		return 0, err
	}
	r.next++
	r.calls = append(r.calls, recCall{tool: tool, pattern: pattern})
	return r.next, nil
}

func (r *recGrantRecorder) got() []recCall {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recCall(nil), r.calls...)
}

// logBox collects the audit lines the gate wrote, so an assertion reads the same
// sentence an operator would have read in wisp.log.
type logBox struct {
	mu    sync.Mutex
	lines []string
}

func (b *logBox) write(format string, args ...any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, fmt.Sprintf(format, args...))
}

func (b *logBox) all() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return strings.Join(b.lines, "\n")
}

func (b *logBox) has(needle string) bool { return strings.Contains(b.all(), needle) }

// grantFixture is one gate plus the three objects these cases read: the UI the
// card is displayed on (the only carrier of the nonce), the recorder, the audit
// box. withRecorder=false builds the pre-224 host: a gate that cannot store a
// scope at all.
type grantFixture struct {
	rec *recGrantRecorder
	log *logBox
	ui  *fakeUI
	g   *approval.Gate
}

func newGrantFixture(t *testing.T, withRecorder bool) *grantFixture {
	t.Helper()
	f := &grantFixture{ui: newFakeUI(), log: &logBox{}}
	o := approval.Options{
		UI:       f.ui,
		Clock:    newFakeClock(),
		Channels: approval.NewChannels(),
		Logf:     f.log.write,
	}
	// The recorder is assigned into the interface only when one exists. Written
	// that way on purpose: `Grants: f.rec` with a nil *recGrantRecorder produces a
	// NON-nil interface holding a nil pointer, which is the exact typed-nil shape
	// cmd/wisp/run.go:511 names as load-bearing in production. This fixture tripped
	// over it first, and the missing GRANT-DROPPED line is what showed it.
	if withRecorder {
		f.rec = newRecRecorder()
		o.Grants = f.rec
	}
	f.g = approval.New(o)
	t.Cleanup(func() {
		if n := f.g.Queue().Depth(); n != 0 {
			t.Errorf("%d pending approvals left at exit", n)
		}
	})
	return f
}

// showCard admits this task through D47's gate, routes one L2 decision through
// the real queue, and returns the Prompt the UI got (its Grant and Paths are the
// gate's own, never literals typed here) plus the channel the routed answer
// arrives on.
func (f *grantFixture) showCard(t *testing.T, paths ...string) (approval.Prompt, <-chan answer) {
	t.Helper()
	revoke := f.g.AdmitTextTask(testTask)
	t.Cleanup(revoke)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	ch := runApproval(t, f.g, ctx, l2Decision(paths...))
	p := f.ui.wait(t)
	if p.Level != "L2" {
		t.Fatalf("card level=%q, want L2: an L1 window has no allow verb, so there would be "+
			"no answer to record", p.Level)
	}
	return p, ch
}

// answerSession spends the card's nonce through the native session route.
func (f *grantFixture) answerSession(t *testing.T, p approval.Prompt) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := f.g.Native().AllowSession(ctx, p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("Native().AllowSession: %v", err)
	}
}

// ---------------------------------------------------------------------------
// N#1 claim 1 + claim 2: one answer, one row per printed path, named honestly
// ---------------------------------------------------------------------------

// TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNamed is the judge M4
// defeated. Deleting Gate.allowSession's Record call - or letting it return a
// zero id - makes this case go red by naming it, because both the recorder's call
// list and the audit line's grant_id are read back here.
func TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNamed(t *testing.T) {
	f := newGrantFixture(t, true)
	p, ch := f.showCard(t, "C:/work/a.txt", "C:/work/b.txt")

	if p.Grant == "" || !strings.HasPrefix(p.Grant, "grant_") {
		t.Fatalf("card carries no native nonce (Grant=%q): there is no answer to test", p.Grant)
	}
	f.answerSession(t, p)

	a := mustAnswer(t, ch)
	if a.a != tools.AnswerAllow {
		t.Fatalf("the answered call came back %v (%s): the answer has to release the call it "+
			"is answering before it thinks about the scope", a.a, a.why)
	}

	// CLAIM 1 - the row. One call per path the card printed, carrying the card's
	// own tool. A recorder that saw nothing is what M4 produced, and nothing else
	// in the suite could tell.
	got := f.rec.got()
	if len(got) != len(p.Paths) {
		t.Fatalf("recorder saw %d Record calls, want %d (the card printed %d paths, calls=%v): "+
			"answering 「本会话内允许」 without a row per path IS the silently-grants-nothing "+
			"shape - the user clicked remember, the disk remembered nothing, the next similar "+
			"call asks again, and the log says otherwise",
			len(got), len(p.Paths), len(p.Paths), got)
	}
	for i, c := range got {
		if c.tool != p.Tool {
			t.Errorf("call %d recorded tool %q, want the card's own %q", i, c.tool, p.Tool)
		}
		if c.pattern != p.Paths[i] {
			t.Errorf("call %d recorded pattern %q, want the path the card printed %q: a rule "+
				"wider than what was shown is an authorization nobody clicked", i, c.pattern, p.Paths[i])
		}
	}

	// CLAIM 2 - the log must not lie. The grant_id in the audit line is the id the
	// recorder actually returned. This half stays red even if someone re-adds a
	// Record call that hands back 0.
	for i, c := range got {
		line := fmt.Sprintf("approval: GRANT-RECORDED corr=%s grant_id=%d tool=%s pattern=%q",
			p.CorrelationID, 11+i, c.tool, c.pattern)
		if !f.log.has(line) {
			t.Errorf("audit is missing the recording line %q (got:\n%s): a row that was written "+
				"has to be nameable with the id the recorder returned, or the line is claiming a "+
				"row it cannot point at", line, f.log.all())
		}
	}
	if f.log.has("grant_id=0 ") {
		t.Errorf("audit claims a recorded row with grant_id=0, i.e. with no row at all:\n%s", f.log.all())
	}
}

// TestTicket224ForgedSessionAnswerRecordsNothing keeps "the answer" the only
// entrance to the row: an answer that never spent this card's live nonce must not
// put anything on disk. Without this half, "make Record get called" could be
// satisfied by a gate that writes a rule for any request naming a correlation id.
func TestTicket224ForgedSessionAnswerRecordsNothing(t *testing.T) {
	f := newGrantFixture(t, true)
	p, ch := f.showCard(t, "C:/work/forged.txt")
	ctx := context.Background()

	for _, bad := range []string{"", "grant_00", "sess_deadbeef", p.Grant + "x"} {
		if err := f.g.Native().AllowSession(ctx, p.CorrelationID, bad); err == nil {
			t.Errorf("AllowSession(nonce=%q) succeeded, want a refusal: a card the native side "+
				"never proved must not create a standing rule", bad)
		}
	}
	if n := len(f.rec.got()); n != 0 {
		t.Fatalf("%d rows recorded across 4 refused answers: a refused answer is not an answer, "+
			"and this is the widening AGENTS.md §1.2 ban #6 (a panel-sourced allow) exists to keep out", n)
	}
	if f.log.has("GRANT-RECORDED") {
		t.Errorf("audit claims a row for a card that was never answered:\n%s", f.log.all())
	}

	// The card is still live: the honest answer afterwards does land, so the four
	// refusals above are not "this gate never records".
	f.answerSession(t, p)
	mustAnswer(t, ch)
	if got := f.rec.got(); len(got) != 1 || got[0].pattern != "C:/work/forged.txt" {
		t.Errorf("recorder calls = %v, want the one path of the real answer", got)
	}
}

// ---------------------------------------------------------------------------
// N#1 controls: the two ways the row can honestly be absent, and what gets said
// ---------------------------------------------------------------------------

// TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped is the nil
// recorder's contract: the call is released as a single-use allow and the audit
// says nothing was stored. It is the same claim as M4's, read from the other
// side - a future change that prints GRANT-RECORDED while holding no recorder
// cannot pass this either.
func TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped(t *testing.T) {
	f := newGrantFixture(t, false)
	p, ch := f.showCard(t, "C:/work/dropped.txt")

	if err := f.g.Native().AllowSession(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("a host with no ledger must still be answerable, got %v", err)
	}
	if a := mustAnswer(t, ch); a.a != tools.AnswerAllow {
		t.Errorf("answer = %v, want allow: dropping the SCOPE must not drop the call", a.a)
	}
	if !f.log.has("approval: GRANT-DROPPED") {
		t.Errorf("audit is missing the drop line. The operator was told 「本会话内」 and nothing was "+
			"stored, so this sentence is the only trace of that. got:\n%s", f.log.all())
	}
	if f.log.has("GRANT-RECORDED") || f.log.has("grant_id=") {
		t.Errorf("audit claims a stored grant on a host that has no recorder:\n%s", f.log.all())
	}
}

// TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow is the partial
// branch: two paths, one row lands, one write fails. The call was already released
// (the answer is proved by the nonce spend, which happens first), so the failing
// scope has to be reported as a failure - and the path that DID land must still be
// named with its own id.
func TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow(t *testing.T) {
	f := newGrantFixture(t, true)
	const broke, fine = "C:/work/broke.txt", "C:/work/fine.txt"
	f.rec.failOn[broke] = errors.New("fixture: 盘写不进去")

	p, ch := f.showCard(t, broke, fine)
	f.answerSession(t, p)

	if a := mustAnswer(t, ch); a.a != tools.AnswerAllow {
		t.Errorf("answer = %v, want allow: a recording fault is not a refusal of the user's click", a.a)
	}
	if !f.log.has("GRANT-RECORD-FAILED") {
		t.Errorf("the failed write is invisible in the audit, so the operator would believe the "+
			"scope stuck. got:\n%s", f.log.all())
	}
	if n := strings.Count(f.log.all(), "GRANT-RECORDED"); n != 1 {
		t.Errorf("audit has %d GRANT-RECORDED lines, want exactly 1 for a 2-path answer whose "+
			"other write failed:\n%s", n, f.log.all())
	}
	// The one line that exists names the successful path with the recorder's real
	// id (11, because the failed call consumed none).
	if !f.log.has(fmt.Sprintf("GRANT-RECORDED corr=%s grant_id=11 tool=%s pattern=%q",
		p.CorrelationID, p.Tool, fine)) {
		t.Errorf("the path that landed is not named with its id (this recorder starts at 11):\n%s",
			f.log.all())
	}
	if got := f.rec.got(); len(got) != 1 || got[0].pattern != fine {
		t.Errorf("recorder calls = %v, want exactly the one path that succeeded", got)
	}
}
