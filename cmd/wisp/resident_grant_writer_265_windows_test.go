//go:build windows

package main

// Ticket 265-r1 (form ⓐ-Ⅰ, orchestrator ruling A601 §4, named unfreeze A602) -
// the pins for "the resident leg's approval gate has a session-grant writer, and
// that writer is honest about the window in which no ledger is attached yet".
//
// Why these cases had to be written at all: 265-a1 census §5.2 measured that the
// resident leg's grant half was guarded by ZERO cases. The only instrument
// mentioning it was ruler ④ of the 256 file, and it guarded the opposite polarity
// - it demanded Options.Grants stay ABSENT until this ticket was adjudicated, and
// A602 turned it over in the same commit as this file. The two log rulers that do
// exist (internal/agent/approval/ticket224_reply_grant_test.go wants GRANT-DROPPED
// on a recorder-less host, cmd/wisp/ticket224_assembly_test.go refuses it on the
// console host) cannot see this leg: neither builds the gate this process builds.
//
// What is claimed, in the order the cases run it:
//
//	POSTURE     - while the holder is unbound, Record returns an error, the gate's
//	              own line becomes GRANT-RECORD-FAILED, and no id is invented.
//	              This is A601 §4 hard constraint 1: the one way this form could
//	              have been made WORSE than the shape it replaces, because a nil
//	              error would let approval_reply.go:277-281 print 「记入本会话」 and
//	              book an ANSWERED audit row for a rule nobody wrote.
//	CONNECTION  - both wiring sites are read off the syntax tree: the value behind
//	              Options.Grants is this process's holder, and the bind runs after
//	              `src.run = run` and before any task can be submitted. Delete
//	              either and a name here reddens - that is AC#3's ⓐ arm ("摘掉那处
//	              接线必须让指名用例红"), which no pre-existing case could answer.
//	NO 2nd MINT - the ledger still has exactly one production construction site,
//	              which is what separates ⓐ-Ⅰ from the ⓐ-Ⅲ shape the ruling
//	              refused by name (gate writes under id-A, bridge reads under id-B).
//
// NOT claimed here, stated so the absence is not read as a pass: a real
// double-clicked resident process answering 「本会话内允许」 on a real card (that
// shape has no answer entry at all - census §2.1 - and no rig here reaches it),
// and any WebView2 page, since this tree has no Go-to-page channel. The gate in
// the fixture below is a real approval.Gate carrying the real holder, but it is
// not the object newResidentApprovalWithConfig returns; that one needs a ball
// window before it will display a card, and its own construction is what
// TestTicket265ResidentApprovalConstructsAnUnboundHolder runs.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

const (
	residentApprovalFile265Name   = "resident_approval_windows.go"
	residentTaskSourceFile265Name = "resident_task_source_windows.go"
	grant265Task                  = "task-265"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// grant265Recorder is the stand-in ledger. Ids start at 21 so a grant_id=0, or a
// small invented one, cannot be read as a real row.
type grant265Recorder struct {
	mu      sync.Mutex
	calls   [][2]string
	next    int64
	err     error
	sawCtx  bool
	wantCtx string
}

type ctxKey265 struct{}

func newGrant265Recorder() *grant265Recorder { return &grant265Recorder{next: 20} }

func (r *grant265Recorder) Record(ctx context.Context, tool, pattern string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return 0, r.err
	}
	if v, ok := ctx.Value(ctxKey265{}).(string); ok && v == r.wantCtx {
		r.sawCtx = true
	}
	r.next++
	r.calls = append(r.calls, [2]string{tool, pattern})
	return r.next, nil
}

func (r *grant265Recorder) got() [][2]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([][2]string(nil), r.calls...)
}

// grant265Log collects the gate's audit lines exactly as the resident leg's
// residentAuditf would have received them, so the sentences asserted below are
// the ones an operator (or a reader of <dataDir>\logs) would have.
type grant265Log struct {
	mu    sync.Mutex
	lines []string
}

func (l *grant265Log) write(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, args...))
}

func (l *grant265Log) all() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.lines, "\n")
}

func (l *grant265Log) has(needle string) bool { return strings.Contains(l.all(), needle) }

// grant265UI is the surface the card is displayed on. The single-use nonce the
// answer must carry travels on this Prompt only, so a fixture without a UI has no
// answer to give.
type grant265UI struct {
	mu       sync.Mutex
	prompts  []approval.Prompt
	prompted chan approval.Prompt
}

func newGrant265UI() *grant265UI {
	return &grant265UI{prompted: make(chan approval.Prompt, 8)}
}

func (u *grant265UI) Prompt(_ context.Context, p approval.Prompt) error {
	u.mu.Lock()
	u.prompts = append(u.prompts, p)
	u.mu.Unlock()
	u.prompted <- p
	return nil
}

func (u *grant265UI) Update(_ context.Context, _ approval.Event) error { return nil }

func (u *grant265UI) wait(t *testing.T) approval.Prompt {
	t.Helper()
	select {
	case p := <-u.prompted:
		return p
	case <-time.After(3 * time.Second):
		t.Fatal("the gate never displayed a card within 3s")
		return approval.Prompt{}
	}
}

type answer265 struct {
	a   tools.Answer
	why string
}

// gate265Fixture is one real approval.Gate carrying one real residentGrantHolder.
// Options carries four of the fields the resident literal carries (UI / Channels /
// Logf / Grants); Window and ApprovalTimeout are left at zero, which
// approval.New clamps to its own defaults - this file pins the grant half, not
// ticket 256's [risk] half, and the 256 file owns those two.
type gate265Fixture struct {
	g      *approval.Gate
	holder *residentGrantHolder
	log    *grant265Log
	ui     *grant265UI
	ord    int
}

func newGate265Fixture(t *testing.T) *gate265Fixture {
	t.Helper()
	f := &gate265Fixture{holder: &residentGrantHolder{}, log: &grant265Log{}, ui: newGrant265UI()}
	f.g = approval.New(approval.Options{
		UI:       f.ui,
		Channels: approval.NewChannels(),
		Logf:     f.log.write,
		Grants:   f.holder,
	})
	revoke := f.g.AdmitTextTask(grant265Task)
	t.Cleanup(revoke)
	t.Cleanup(func() {
		if n := f.g.Queue().Depth(); n != 0 {
			t.Errorf("%d pending approvals left at exit", n)
		}
	})
	return f
}

// showCard routes one L2 decision through the real queue and returns the Prompt
// the UI got - its CorrelationID, Grant and Paths are the gate's own, never
// literals typed here - plus the channel the routed answer arrives on.
func (f *gate265Fixture) showCard(t *testing.T, paths ...string) (approval.Prompt, <-chan answer265) {
	t.Helper()
	f.ord++
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	out := make(chan answer265, 1)
	go func() {
		defer close(out)
		a, why := f.g.PendingApproval(ctx, l2Decision265(f.ord, paths))
		out <- answer265{a: a, why: why}
	}()
	p := f.ui.wait(t)
	if p.Level != "L2" {
		t.Fatalf("card level=%q, want L2: an L1 window has no allow verb, so there would be no "+
			"answer to record", p.Level)
	}
	if p.Grant == "" || !strings.HasPrefix(p.Grant, "grant_") {
		t.Fatalf("card carries no native nonce (Grant=%q): there is no answer to test", p.Grant)
	}
	return p, out
}

func l2Decision265(n int, paths []string) tools.Decision {
	items := make([]any, 0, len(paths))
	for _, p := range paths {
		items = append(items, p)
	}
	corr := fmt.Sprintf("265-corr-%d", n)
	return tools.Decision{
		Tool:          "fs.write",
		Provider:      tools.KindBuiltin,
		Params:        map[string]any{"paths": items},
		Args:          json.RawMessage(`{"paths":["a"]}`),
		Level:         risk.L2,
		RulesHit:      []risk.RuleID{risk.R1, risk.R2},
		Reason:        "目标越出 [fs] 授权目录（R2），需 L2 强确认",
		Paths:         paths,
		CorrelationID: corr,
		TaskID:        grant265Task,
		CallID:        fmt.Sprintf("call-265-%d", n),
	}
}

// ---------------------------------------------------------------------------
// POSTURE: the unbound holder, on its own
// ---------------------------------------------------------------------------

// TestTicket265ResidentGrantHolderUnboundFailsLoudly is hard constraint 1 run
// rather than written. Seed the quiet shape - make the unbound branch return
// (0, nil), or hand back an invented id - and the first two assertions redden.
func TestTicket265ResidentGrantHolderUnboundFailsLoudly(t *testing.T) {
	h := &residentGrantHolder{}
	id, err := h.Record(context.Background(), "fs.write", "C:/work/unbound.txt")
	if err == nil {
		t.Errorf("an UNBOUND grant holder returned a nil error. That is the shape A601 §4 forbids: " +
			"approval.Gate would log GRANT-RECORDED, replySurface.session would print the " +
			"「记入本会话」 receipt (approval_reply.go:279-281) and book the ANSWERED audit row " +
			"(:277-278), all of it about a rule nobody wrote. Today's GRANT-DROPPED line at least " +
			"told the truth; this would be a lie, and a persisted one.")
	}
	if id != 0 {
		t.Errorf("unbound holder returned grant id %d. An id must come from a store holding a row; "+
			"a fabricated one feeds the audit line a number no reader can look up.", id)
	}
	if err != nil && strings.TrimSpace(err.Error()) == "" {
		t.Errorf("the unbound error carries no sentence, and the gate interpolates it straight into "+
			"GRANT-RECORD-FAILED: got %#v", err)
	}
	if h.bound() {
		t.Errorf("a fresh holder reports itself bound")
	}

	// A nil *session.Ledger must not become a non-nil interface holding nil.
	h.bindSessionLedger(nil)
	if h.bound() {
		t.Errorf("bindSessionLedger(nil) marked the holder bound: the typed-nil guard is the same " +
			"trap cmd/wisp/run.go:570 names for the read side, and a bound holder whose target is a " +
			"nil ledger dereferences on the first answer.")
	}
	if _, err := h.Record(context.Background(), "fs.write", "C:/work/still-unbound.txt"); err == nil {
		t.Errorf("after bindSessionLedger(nil) the holder still claimed success")
	}

	h.bind(nil)
	if h.bound() {
		t.Errorf("bind(nil) marked the holder bound")
	}
}

// TestTicket265ResidentGrantHolderBoundDelegatesEveryField is the other half of
// the same object: once bound, the call reaches the ledger with this call's own
// context, tool and pattern, and the ledger's id is what comes back.
func TestTicket265ResidentGrantHolderBoundDelegatesEveryField(t *testing.T) {
	rec := newGrant265Recorder()
	rec.wantCtx = "carry-me"
	h := &residentGrantHolder{}
	h.bind(rec)
	if !h.bound() {
		t.Fatal("bound() says false after bind, so nothing below can mean anything")
	}
	ctx := context.WithValue(context.Background(), ctxKey265{}, "carry-me")
	id, err := h.Record(ctx, "fs.write", "C:/work/a.txt")
	if err != nil {
		t.Fatalf("Record through a bound holder: %v", err)
	}
	if id != 21 {
		t.Errorf("Record returned id %d, want the id the recorder handed back (21)", id)
	}
	if !rec.sawCtx {
		t.Errorf("the recorder did not see this call's context: a holder that drops ctx breaks " +
			"whatever cancellation or trace the ledger relies on")
	}
	if got := rec.got(); len(got) != 1 || got[0][0] != "fs.write" || got[0][1] != "C:/work/a.txt" {
		t.Errorf("recorder calls = %v, want [[fs.write C:/work/a.txt]]", got)
	}
	if _, err := h.Record(ctx, "fs.write", "C:/work/b.txt"); err != nil {
		t.Fatalf("second Record: %v", err)
	}
	if got := rec.got(); len(got) != 2 {
		t.Errorf("recorder saw %d calls after two Records, want 2", len(got))
	}
}

// TestTicket265LedgerErrorReachesTheOperatorThroughTheHolder is the third shape:
// bound, and the store refuses. The holder passes the error through rather than
// swallowing it, or the failure the user is told about disappears.
func TestTicket265LedgerErrorReachesTheOperatorThroughTheHolder(t *testing.T) {
	rec := newGrant265Recorder()
	rec.err = errors.New("fixture: 账本写不进去")
	h := &residentGrantHolder{}
	h.bind(rec)
	if _, err := h.Record(context.Background(), "fs.write", "C:/work/err.txt"); err == nil {
		t.Errorf("the holder swallowed the ledger's error")
	} else if !strings.Contains(err.Error(), "账本写不进去") {
		t.Errorf("Record error = %v, want the ledger's own sentence passed through", err)
	}
}

// ---------------------------------------------------------------------------
// POSTURE, through the gate: the sentence an operator actually gets
// ---------------------------------------------------------------------------

// TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow is the
// between-construction-and-bind window, and the whole life of a boot whose
// session mint failed: the answer releases the call, is NOT reported as recorded,
// and names itself a failure.
func TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow(t *testing.T) {
	f := newGate265Fixture(t)
	p, ch := f.showCard(t, "C:/work/unbound-card.txt")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := f.g.Native().AllowSession(ctx, p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("Native().AllowSession: %v: a holder with no ledger behind it must still answer "+
			"the call the user is answering", err)
	}
	select {
	case a := <-ch:
		if a.a != tools.AnswerAllow {
			t.Errorf("the routed call came back %v (%s): a writer that cannot write must not revoke "+
				"the answer the user already gave", a.a, a.why)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("the answered call never came back within 3s")
	}

	if f.holder.bound() {
		t.Errorf("the holder became bound without a bind call, so the assertions below are theatre")
	}
	log := f.log.all()
	if !f.log.has("GRANT-RECORD-FAILED") {
		t.Errorf("audit is missing gate.go's GRANT-RECORD-FAILED line, which is the whole of this "+
			"form's honesty while the holder is unbound. got:\n%s", log)
	}
	if f.log.has("GRANT-RECORDED") || f.log.has("grant_id=") {
		t.Errorf("audit claims a stored row while the holder was unbound:\n%s", log)
	}
	if f.log.has("GRANT-DROPPED") {
		t.Errorf("audit fell back to GRANT-DROPPED, i.e. the gate saw grants == nil, so Options.Grants "+
			"was never passed at all. That is the wiring this ticket exists to pin.\ngot:\n%s", log)
	}
}

// TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed is the same card
// one bind later: the promise the two sentences at approval_reply.go:277-281 make,
// kept.
func TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed(t *testing.T) {
	f := newGate265Fixture(t)
	rec := newGrant265Recorder()
	f.holder.bind(rec)
	if !f.holder.bound() {
		t.Fatal("bind did not take: the holder still reports unbound")
	}

	p, ch := f.showCard(t, "C:/work/r1.txt", "C:/work/r2.txt")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := f.g.Native().AllowSession(ctx, p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("Native().AllowSession: %v", err)
	}
	select {
	case a := <-ch:
		if a.a != tools.AnswerAllow {
			t.Errorf("routed answer = %v (%s), want allow", a.a, a.why)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("the answered call never came back within 3s")
	}

	got := rec.got()
	if len(got) != len(p.Paths) {
		t.Fatalf("recorder saw %d calls, want one per path the card printed (%d): calls=%v. This is "+
			"AC#1's other half: a bound writer that records nothing puts the same false 「记入本会话」 "+
			"receipt on screen as the unbound one did, and books the audit row too",
			len(got), len(p.Paths), got)
	}
	for i, c := range got {
		if c[0] != p.Tool {
			t.Errorf("call %d recorded tool %q, want the card's own %q", i, c[0], p.Tool)
		}
		if c[1] != p.Paths[i] {
			t.Errorf("call %d recorded pattern %q, want the printed path %q", i, c[1], p.Paths[i])
		}
		line := fmt.Sprintf("approval: GRANT-RECORDED corr=%s grant_id=%d tool=%s pattern=%q",
			p.CorrelationID, 21+i, c[0], c[1])
		if !f.log.has(line) {
			t.Errorf("audit is missing the recording line %q. got:\n%s", line, f.log.all())
		}
	}
	for _, banned := range []string{"GRANT-RECORD-FAILED", "GRANT-DROPPED", "grant_id=0 "} {
		if f.log.has(banned) {
			t.Errorf("audit carries %q on a host whose holder IS bound:\n%s", banned, f.log.all())
		}
	}
}

// ---------------------------------------------------------------------------
// CONNECTION: the two wiring sites, read off the syntax tree
// ---------------------------------------------------------------------------

// TestTicket265GateLiteralCarriesTheResidentHolder reads the resident leg's
// approval.New literal and demands the value behind Options.Grants is this
// process's own holder - not a per-call literal no bind site can reach. Ruler ④
// of the 256 file proves the FIELD is passed; this proves what is passed INTO it,
// which is the half that decides whether the bind can ever reach the gate.
func TestTicket265GateLiteralCarriesTheResidentHolder(t *testing.T) {
	path := filepath.Join(packageDir265(t), residentApprovalFile265Name)
	sel := optionsValueSel265(t, path, "Grants")
	if sel != "ra.grants" {
		t.Errorf("Options.Grants in %s is %q, want \"ra.grants\": anything else is a holder no bind "+
			"site can reach - a fresh per-call holder leaves the gate writing into an object that is "+
			"never bound, which is GRANT-RECORD-FAILED for the life of the process.", path, sel)
	}
}

// TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask is AC#3's second
// control: delete the bind call in startResidentTaskSource, move it below
// submitTask, or point it at something other than run.session, and this name
// reddens. Position is the claim, so position is what is read.
func TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask(t *testing.T) {
	path := filepath.Join(packageDir265(t), residentTaskSourceFile265Name)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	fn := funcDecl265(file, "startResidentTaskSource")
	if fn == nil {
		t.Fatalf("no func startResidentTaskSource in %s: the resident leg's task entry moved or was "+
			"renamed, and this ruler's whole claim is about where its bind line sits", path)
	}
	var (
		assignRun   token.Pos
		bindCall    token.Pos
		bindArg     string
		submitCalls []token.Pos
	)
	ast.Inspect(fn, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			if len(x.Lhs) == 1 && len(x.Rhs) == 1 && exprString265(x.Lhs[0]) == "src.run" {
				assignRun = x.Pos()
			}
		case *ast.CallExpr:
			switch exprString265(x.Fun) {
			case "ra.bindResidentGrantLedger":
				bindCall = x.Pos()
				if len(x.Args) == 1 {
					bindArg = exprString265(x.Args[0])
				}
			case "src.submitTask":
				submitCalls = append(submitCalls, x.Pos())
			}
		}
		return true
	})
	if !bindCall.IsValid() {
		t.Errorf("%s never calls ra.bindResidentGrantLedger inside startResidentTaskSource: the "+
			"resident gate's holder would stay unbound for the life of the process. That is exactly "+
			"what A601 §4's hard constraint 1 is written against, and exactly what AC#3's 「摘掉那处"+
			"接线」 arm asks to see red.", path)
		return
	}
	if bindArg != "run.session" {
		t.Errorf("the bind passes %q, want \"run.session\": that is this process's ONE minted ledger "+
			"(run.go:473-487). Anything else is the ⓐ-Ⅲ shape the ruling refused by name.", bindArg)
	}
	if !assignRun.IsValid() || bindCall <= assignRun {
		t.Errorf("the bind at line %v is not after `src.run = run` at line %v: binding before the "+
			"assembly returned means there is no ledger to bind yet, so the line would be decoration.",
			fset.Position(bindCall).Line, fset.Position(assignRun).Line)
	}
	if len(submitCalls) == 0 {
		t.Errorf("no src.submitTask call inside startResidentTaskSource any more: this ruler's " +
			"ordering claim (\"before any task can be submitted\") needs that landmark to mean " +
			"anything. If the entry moved, re-site this ruler rather than dropping the claim.")
	}
	for _, s := range submitCalls {
		if bindCall >= s {
			t.Errorf("the bind at line %v runs AFTER a submitTask at line %v: the first card of that "+
				"task could be answered while the holder is still unbound, and the user would be "+
				"handed GRANT-RECORD-FAILED one line after the ledger came to exist.",
				fset.Position(bindCall).Line, fset.Position(s).Line)
		}
	}
}

// TestTicket265SessionLedgerStillHasOneProductionConstructionSite is the ⓐ-Ⅲ
// refusal, run: the mint and the ledger must still be built in exactly one place,
// because a second one means the gate writes rows under one session id while the
// bridge reads under another - AC#1's rejected 「记到了别的会话」 arm, arriving with
// a GRANT-RECORDED line that names a real, unqueryable grant_id.
func TestTicket265SessionLedgerStillHasOneProductionConstructionSite(t *testing.T) {
	dir := packageDir265(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	var hits []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for i, line := range strings.Split(string(b), "\n") {
			if strings.Contains(line, "session.Mint(") || strings.Contains(line, "session.NewLedger(") {
				hits = append(hits, fmt.Sprintf("%s:%d", name, i+1))
			}
		}
	}
	if len(hits) == 0 {
		t.Errorf("no production file in %s mints a session any more: the ledger this leg binds is the "+
			"one cmd/wisp/run.go's assembleRuntime mints, so a ruler that cannot find it has lost the "+
			"object both halves are supposed to share.", dir)
	}
	for _, h := range hits {
		if !strings.HasPrefix(h, "run.go:") {
			t.Errorf("a second production construction site appeared: %s. Ticket 265's ruling (A601 §4) "+
				"refuses form ⓐ-Ⅲ by name precisely for this: grant rows would be written under one "+
				"session id and read under another.", h)
		}
	}
}

// ---------------------------------------------------------------------------
// the shipped constructor's own state
// ---------------------------------------------------------------------------

// TestTicket265ResidentApprovalConstructsAnUnboundHolder runs the production
// constructor: a holder exists (so Options.Grants is never a nil interface at
// construction), it starts unbound (fail-closed is the default, not prose), and
// bindResidentGrantLedger is what moves it - including the case that must NOT
// move it, a nil ledger from a boot that could not mint.
func TestTicket265ResidentApprovalConstructsAnUnboundHolder(t *testing.T) {
	ra := newResidentApproval()
	t.Cleanup(ra.cancel)
	if ra.gate == nil {
		t.Fatal("newResidentApproval built no gate")
	}
	if ra.grants == nil {
		t.Fatal("the resident approval carries no grant holder, so the gate's Options.Grants is a " +
			"nil interface and this leg is back to the pre-265 shape with nothing wired at all")
	}
	if ra.grants.bound() {
		t.Errorf("a freshly built resident approval reports its grant writer bound: nothing has been " +
			"assembled yet, so that claim could only be false")
	}
	if _, err := ra.grants.Record(context.Background(), residentCardTool, "C:/work/pre-bind.txt"); err == nil {
		t.Errorf("the shipped holder recorded before any bind. Every sentence downstream of that " +
			"success - the console receipt and the ANSWERED audit row - would be unearned.")
	}

	ra.bindResidentGrantLedger(nil)
	if ra.grants.bound() {
		t.Errorf("bindResidentGrantLedger(nil) bound the holder to nothing")
	}

	rec := newGrant265Recorder()
	ra.grants.bind(rec)
	if !ra.grants.bound() {
		t.Fatal("bind through the shipped holder did not take")
	}
	if _, err := ra.grants.Record(context.Background(), residentCardTool, "C:/work/post-bind.txt"); err != nil {
		t.Errorf("Record through the bound shipped holder: %v", err)
	}
	if got := rec.got(); len(got) != 1 {
		t.Errorf("recorder saw %v, want the one call", got)
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// packageDir265 resolves this package's own directory off the caller's file, the
// way the 256 file's risk265Dir does: a prebuilt test binary runs from elsewhere,
// and a ruler that reads sibling sources must not turn its own path handling into
// an apparent code defect.
func packageDir265(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed for the 265 ruler")
	}
	return filepath.Dir(thisFile)
}

// optionsValueSel265 returns the source text of one field of the single
// approval.Options literal in path, or fails: two literals would mean the resident
// leg builds a second gate, which ruler ④ and the 246 pointer-identity pin already
// refuse.
func optionsValueSel265(t *testing.T, path, field string) string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	var found []string
	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		if exprString265(lit.Type) != "approval.Options" {
			return true
		}
		for _, el := range lit.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if exprString265(kv.Key) == field {
				found = append(found, exprString265(kv.Value))
			}
		}
		return true
	})
	if len(found) == 0 {
		t.Fatalf("no Options.%s in %s: the field this ruler reads is gone, which is the same defect "+
			"ruler ④ of the 256 file reports from the other direction.", field, path)
	}
	if len(found) > 1 {
		t.Fatalf("%d Options.%s values in %s, want 1", len(found), field, path)
	}
	return found[0]
}

func funcDecl265(file *ast.File, name string) *ast.FuncDecl {
	for _, d := range file.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if ok && fn.Recv == nil && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}

// exprString265 renders the handful of expression shapes these rulers look for.
// Anything unrecognised prints its Go type, so a refactor shows up as a mismatched
// string rather than as a silent miss.
func exprString265(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return exprString265(x.X) + "." + x.Sel.Name
	case *ast.StarExpr:
		return "*" + exprString265(x.X)
	case *ast.CallExpr:
		return exprString265(x.Fun) + "(...)"
	case nil:
		return "<nil>"
	default:
		return fmt.Sprintf("<%T>", e)
	}
}
