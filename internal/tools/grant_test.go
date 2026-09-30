package tools

// Ticket 224 tests: the D45 scoped session grant inside the real decision chain.
//
// These are the internal/tools half of AC#2's 读 leg. The 写 and 失效 halves have
// their own cases elsewhere on purpose (internal/session/grants_test.go books the
// row and measures invalidation against a second minted identity; cmd/wisp's
// ticket224 production case runs a whole boot). AC#2's own wording is why they
// are not one test: "三格各自有独立用例，不许合并成一枚持久化用例", and
// internal/perm/ticket90_persist_test.go's header says the same thing about the
// restart pair - a merged persistence case can only ever be satisfied by the
// looser half, which is how a permanent免审通行证 gets in.
//
// Every positive assertion here carries its negative control in the same
// function, because the failure this feature can have is not "it broke" but "it
// silently grants nothing and everything still asks" - which reads exactly like
// a passing test unless something asserts the opposite shape.
//
// AC#4 anchors, named per case so a mutation can be pointed at one branch:
//
//	Bridge.Execute's grant gate      `if !sil.Silenced && sil.Level == risk.L1`
//	Bridge.route's L1 grant branch   `if grantID != 0 { ... DecisionAllowGrant }`
//	Bridge.sessionGrantID            the nil-source and recover fail-closed arms
//	Bridge.book                      the grant_id / DecisionAllowGrant row
//
// These cases need no sherpa PATH; cmd/wisp's do.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// t224Grants is a GrantSource the test controls. Pointer receiver for the same
// reason t90Modes is one: the bridge holds the interface, so a mutation that
// parked the grant source in a package global would show up as two bridges
// agreeing after they were told to disagree.
type t224Grants struct {
	mu sync.Mutex
	// live maps tool -> the set of canonical paths a grant covers.
	live map[string]map[string]bool
	// grantID is what the source hands back as the covering row's id.
	grantID int64
	// calls records every (tool, paths) the bridge asked about, in order. The
	// "an L2 never even reaches the check" and "a pathless call is never quoted"
	// assertions read this rather than a counter, because consulting-and-refusing
	// and never-consulting are different futures for a mutation.
	calls []string
	// panicOn makes the source blow up mid-lookup; errOn makes it report it
	// could not look. Both are the fail-closed cases.
	panicOn bool
	errOn   bool
}

func newT224Grants(grantID int64) *t224Grants {
	return &t224Grants{live: map[string]map[string]bool{}, grantID: grantID}
}

func (g *t224Grants) add(tool, path string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.live[tool] == nil {
		g.live[tool] = map[string]bool{}
	}
	g.live[tool][path] = true
}

func (g *t224Grants) Covering(_ context.Context, tool string, paths []string) (int64, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.calls = append(g.calls, tool+" "+strings.Join(paths, ","))
	if g.panicOn {
		panic("ticket 224 fixture: simulated grant lookup blow-up")
	}
	if g.errOn || len(paths) == 0 {
		return 0, false
	}
	for _, p := range paths {
		if !g.live[tool][p] {
			return 0, false
		}
	}
	return g.grantID, true
}

func (g *t224Grants) asked() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]string(nil), g.calls...)
}

// t224Journal is the smallest agent.Journal that records what the bridge booked.
// The row is the point: SPEC-06 §8.3 / D45-2 say "每次使用被授权通道都写
// tool_call 日志（可取证）", and tool_call.grant_id has been in the frozen DDL
// (memory/schema.go:83) with nothing able to fill it until this ticket.
type t224Journal struct {
	mu     sync.Mutex
	calls  []memory.ToolCall
	nextID int64
}

func (j *t224Journal) StartTaskLog(context.Context, memory.TaskLog) error { return nil }

func (j *t224Journal) FinishTaskLog(context.Context, string, string, string, int64, int64, int64, string) error {
	return nil
}

func (j *t224Journal) InsertToolCall(_ context.Context, tc memory.ToolCall) (int64, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.nextID++
	tc.ID = j.nextID
	j.calls = append(j.calls, tc)
	return j.nextID, nil
}

func (j *t224Journal) DecideToolCall(context.Context, int64, string, *int64) error { return nil }

func (j *t224Journal) FinishToolCall(context.Context, int64, string, string) error { return nil }

func (j *t224Journal) last() (memory.ToolCall, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if len(j.calls) == 0 {
		return memory.ToolCall{}, false
	}
	return j.calls[len(j.calls)-1], true
}

// t224Bridge is a bridge over one allowlist directory, on ask_every_step (the
// strictest档, so a question definitely exists to be answered) with a controllable
// grant source. allowlist may be empty to model "this call is outside [fs]".
func t224Bridge(t *testing.T, declared risk.Level, grants GrantSource, j agent.Journal,
	allowlist ...string,
) (*Bridge, *t90Gate, *t90Tool, *PathCanonicalizer) {
	t.Helper()
	paths := NewPathCanonicalizer(allowlist, nil)
	e, tool := t90Entry("t224.write", declared, []string{"path"}, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	gate := &t90Gate{}
	b := New(Options{
		Registry: reg,
		Paths:    paths,
		Gate:     gate,
		Journal:  j,
		Modes:    &t90Modes{},
		Grants:   grants,
		Logf:     func(string, ...any) {},
	})
	return b, gate, tool, paths
}

// t224Seed registers a live grant for one path through the SAME canonicalizer the
// bridge uses, so a fixture can never "grant" a spelling the judge would not
// produce.
func t224Seed(t *testing.T, paths *PathCanonicalizer, grants *t224Grants, target string) string {
	t.Helper()
	c, err := paths.Canonicalize(target)
	if err != nil {
		t.Fatalf("Canonicalize(%q): %v", target, err)
	}
	grants.add("t224.write", c)
	return c
}

// ---------------------------------------------------------------------------
// AC#2 读: a live grant in the caller's own session stops the L1 question
// ---------------------------------------------------------------------------

// TestTicket224LiveGrantStopsTheL1Question is the read half of AC#2: the next
// 同类 request hits the stored rule and stops asking.
//
// The control half is not decoration. Without the "no source wired" run, a bridge
// that never consults grants and a bridge whose check silently never matches both
// produce "zero windows", which is also what a working feature produces.
func TestTicket224LiveGrantStopsTheL1Question(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	// ---- control half: nothing on record, so the window is shown ----
	base, gateBase, toolBase, _ := t224Bridge(t, risk.L1, nil, nil, dir)
	if _, err := base.Execute(context.Background(), t90Req("t224.write", target)); err != nil {
		t.Fatalf("Execute (no grant source): %v", err)
	}
	if w, _ := gateBase.counts(); w != 1 {
		t.Fatalf("with no grant source the L1 window ran %d times, want 1: the control half "+
			"is broken, so a zero below would prove nothing", w)
	}
	if toolBase.runs != 1 {
		t.Fatalf("control tool ran %d times, want 1 (t90Gate answers timeout = execute)",
			toolBase.runs)
	}

	// ---- measured half: the same call, with a live grant for this session ----
	grants := newT224Grants(77)
	b, gate, tool, paths := t224Bridge(t, risk.L1, grants, nil, dir)
	canon := t224Seed(t, paths, grants, target)

	if _, err := b.Execute(context.Background(), t90Req("t224.write", target)); err != nil {
		t.Fatalf("Execute (live grant): %v", err)
	}
	if w, _ := gate.counts(); w != 0 {
		t.Errorf("a live session grant left the L1 window on screen (%d asks): AC#2's "+
			"「命中授权就不再弹卡」 is not happening", w)
	}
	if tool.runs != 1 {
		t.Errorf("granted call ran the tool %d times, want 1", tool.runs)
	}
	asked := grants.asked()
	if len(asked) == 0 {
		t.Fatal("the bridge never consulted the grant source, so the zero above is an accident")
	}
	// The canonical form is what must travel, never the raw argument: matching on
	// a raw path is how a normalization bypass becomes an authorization.
	if !strings.Contains(asked[0], canon) {
		t.Errorf("grant lookup asked with %q, want the C26 canonical form %q", asked[0], canon)
	}
}

// TestTicket224PartialPathCoverageStillAsks is the multi-path half of the same
// rule: a grant covering two of a call's three paths covers nothing. A partial
// match executing the call would let one remembered file authorize an
// unremembered neighbour, which is wider than what the user clicked.
func TestTicket224PartialPathCoverageStillAsks(t *testing.T) {
	dir := t.TempDir()
	var targets []string
	for _, n := range []string{"a.txt", "b.txt", "c.txt"} {
		p := filepath.Join(dir, n)
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		targets = append(targets, p)
	}
	grants := newT224Grants(5)
	b, gate, _, paths := t224Bridge(t, risk.L1, grants, nil, dir)
	for _, p := range targets[1:] { // two of three covered
		t224Seed(t, paths, grants, p)
	}
	args, err := json.Marshal(map[string]any{"path": targets})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t224", CorrelationID: "t224-corr", CallID: "t224-call",
		Name: "t224.write", Args: args,
	}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if w, _ := gate.counts(); w != 1 {
		t.Errorf("a call with one ungranted path asked %d windows, want 1: partial coverage "+
			"must not read as coverage", w)
	}
	// Control inside the same boot: cover the third path too and the question
	// goes away, so the reading above is about the missing path and not about a
	// multi-path call the check cannot see.
	t224Seed(t, paths, grants, targets[0])
	if _, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t224", CorrelationID: "t224-corr2", CallID: "t224-call2",
		Name: "t224.write", Args: args,
	}); err != nil {
		t.Fatalf("Execute (fully covered): %v", err)
	}
	if w, _ := gate.counts(); w != 1 {
		t.Errorf("after all three paths were covered the call asked %d windows in total, "+
			"want still 1 (no new question)", w)
	}
}

// ---------------------------------------------------------------------------
// The scope line: SPEC-06 §8.3 bullet 1 / PLAN.md:2148 D45-3
// ---------------------------------------------------------------------------

// TestTicket224SessionGrantNeverCoversL2 is D45's "L2 永不进入任何持久授权
// （含会话级）" as a behavioural assertion, and it is the case that keeps ticket
// 101's guard in cmd/wisp honest rather than accidental.
//
// It is written as a POSITIVE control on purpose: the grant source is told, in so
// many words, that this exact tool and path are covered. If the bridge ever
// learns to honour a grant on an L2, this test is the one that says so.
func TestTicket224SessionGrantNeverCoversL2(t *testing.T) {
	dir := t.TempDir()
	other := t.TempDir()
	outside := filepath.Join(other, "elsewhere.txt")
	inside := filepath.Join(dir, "inside.txt")
	for _, p := range []string{outside, inside} {
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	// An out-of-allowlist write is R2 = L2 even on a tool declaring L1, and a
	// tool declaring L2 stays L2 inside the allowlist. Both halves matter: the
	// second is "the tool's own admission", the first is "the judge decided it".
	grD := newT224Grants(9)
	bDeclared, gateD, _, pathsD := t224Bridge(t, risk.L2, grD, nil, dir)
	t224Seed(t, pathsD, grD, inside)
	if _, err := bDeclared.Execute(context.Background(), t90Req("t224.write", inside)); err != nil {
		t.Fatalf("Execute (declared L2): %v", err)
	}
	if _, a := gateD.counts(); a != 1 {
		t.Errorf("declared-L2 call with a live grant produced %d approval cards, want 1: a "+
			"session grant covered an L2, which SPEC-06 §8.3 bullet 1 forbids", a)
	}
	if got := grD.asked(); len(got) != 0 {
		t.Errorf("the bridge queried the grant source %d times for an L2 verdict (%v); the L2 "+
			"branch must not reach the check at all, so a future edit cannot make it answer "+
			"yes by accident", len(got), got)
	}

	bJudged, gateJ, _, pathsJ := t224Bridge(t, risk.L1, newT224Grants(9), nil, dir)
	grJ := grantsOf(t, bJudged)
	t224Seed(t, pathsJ, grJ, outside)
	if _, err := bJudged.Execute(context.Background(), t90Req("t224.write", outside)); err != nil {
		t.Fatalf("Execute (R2 out-of-scope): %v", err)
	}
	if _, a := gateJ.counts(); a != 1 {
		t.Errorf("R2 out-of-allowlist call with a live grant produced %d approval cards, want 1", a)
	}
	if got := grJ.asked(); len(got) != 0 {
		t.Errorf("the bridge queried the grant source for an R2 call (%v), want no query", got)
	}
}

// TestTicket224SessionGrantNeverCoversDeny is the A-tier half of the same
// sentence: an absolute path is beyond every authorization, remembered or not
// (PLAN.md:1636 「读 ~/.git-credentials … 必须拒绝，且任何授权与会话授权都无法覆盖」).
func TestTicket224SessionGrantNeverCoversDeny(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	key := filepath.Join(home, ".ssh", "id_ed25519")
	if err := os.MkdirAll(filepath.Dir(key), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(key, []byte("private"), 0o600); err != nil {
		t.Fatal(err)
	}
	grants := newT224Grants(11)
	b, gate, tool, paths := t224Bridge(t, risk.L1, grants, nil, home)
	t224Seed(t, paths, grants, key)

	if _, err := b.Execute(context.Background(), t90Req("t224.write", key)); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if tool.runs != 0 {
		t.Errorf("an A-tier path ran the tool %d times, want 0: a session grant executed "+
			"something no authorization can cover", tool.runs)
	}
	if w, a := gate.counts(); w != 0 || a != 0 {
		t.Errorf("A-tier denial reached a gate (windows=%d approvals=%d), want 0/0: Deny "+
			"never becomes a question", w, a)
	}
	if got := grants.asked(); len(got) != 0 {
		t.Errorf("the grant source was consulted on a Deny (%v); Deny must return before the "+
			"check, otherwise the refusal is one flipped comparison from silence", got)
	}
}

// ---------------------------------------------------------------------------
// fail-closed shapes
// ---------------------------------------------------------------------------

// TestTicket224GrantSourceThatCannotAnswerMeansAsking covers the two ways a read
// seam can fail: it reports it could not look, and it blows up. Both must put the
// question back on the user rather than assume the answer.
func TestTicket224GrantSourceThatCannotAnswerMeansAsking(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		fix  func(*t224Grants)
	}{
		{"reports that it could not look", func(g *t224Grants) { g.errOn = true }},
		{"panics mid-lookup", func(g *t224Grants) { g.panicOn = true }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			grants := newT224Grants(3)
			b, gate, _, paths := t224Bridge(t, risk.L1, grants, nil, dir)
			// Seed a rule that WOULD cover it, so the reading below can only come
			// from the source failing and not from "nothing was on record".
			t224Seed(t, paths, grants, target)
			tc.fix(grants)
			if _, err := b.Execute(context.Background(), t90Req("t224.write", target)); err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if w, _ := gate.counts(); w != 1 {
				t.Errorf("a grant source that %s left the call asking %d windows, want 1",
					tc.name, w)
			}
		})
	}
}

// TestTicket224PathlessCallIsNeverGranted states the middle term of the
// (工具, 路径模式, 会话) triple: with no path there is nothing a stored rule could
// have bound to, so the question stands and the ledger is not even asked.
func TestTicket224PathlessCallIsNeverGranted(t *testing.T) {
	dir := t.TempDir()
	grants := newT224Grants(4)
	b, gate, _, paths := t224Bridge(t, risk.L1, grants, nil, dir)
	t224Seed(t, paths, grants, dir)

	if _, err := b.Execute(context.Background(), t90Req("t224.write", "")); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got := grants.asked(); len(got) != 0 {
		t.Errorf("a pathless call consulted the grant source (%v), want no query at all", got)
	}
	if w, _ := gate.counts(); w != 1 {
		t.Errorf("pathless L1 call asked %d windows, want 1", w)
	}
}

// ---------------------------------------------------------------------------
// forensics: D45-2's "每次使用被授权通道都写 tool_call 日志（可取证）"
// ---------------------------------------------------------------------------

// TestTicket224GrantedCallBooksAllowSessionGrantWithItsRow is the audit half: the
// row must say BOTH that this was not a click on this card
// (decision=allow_session_grant, the value agent/journal.go:32 and
// memory/schema.go:76 already freeze) and WHICH rule paid for it
// (tool_call.grant_id, SPEC-02 §3's column).
func TestTicket224GrantedCallBooksAllowSessionGrantWithItsRow(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	j := &t224Journal{}
	grants := newT224Grants(4242)
	b, _, _, paths := t224Bridge(t, risk.L1, grants, j, dir)
	t224Seed(t, paths, grants, target)

	if _, err := b.Execute(context.Background(), t90Req("t224.write", target)); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	row, ok := j.last()
	if !ok {
		t.Fatal("no tool_call row was booked for a call that ran")
	}
	if row.Decision != agent.DecisionAllowGrant {
		t.Errorf("booked decision = %q, want %q: a remembered authorization and a click on "+
			"this card must not land on the same value", row.Decision, agent.DecisionAllowGrant)
	}
	if row.GrantID == nil || *row.GrantID != 4242 {
		t.Errorf("booked grant_id = %v, want the covering row's 4242", row.GrantID)
	}

	// Control: the same call with nothing on record books a plain allow, so the
	// value above is not this bridge's default.
	j2 := &t224Journal{}
	b2, _, _, _ := t224Bridge(t, risk.L1, nil, j2, dir)
	if _, err := b2.Execute(context.Background(), t90Req("t224.write", target)); err != nil {
		t.Fatalf("Execute (control): %v", err)
	}
	ctrl, ok := j2.last()
	if !ok {
		t.Fatal("no tool_call row was booked for the control call")
	}
	if ctrl.Decision == agent.DecisionAllowGrant {
		t.Error("the control call booked allow_session_grant with no grant source wired: " +
			"that value would then be a constant, not a reading")
	}
	if ctrl.GrantID != nil {
		t.Errorf("control row carries grant_id %v, want NULL", *ctrl.GrantID)
	}
}

// ---------------------------------------------------------------------------
// the shape property ticket 90's reflection pin protects, restated for 224
// ---------------------------------------------------------------------------

// TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority is ticket 224's share of
// the boundary TestTicket90ModeCarriesNoAllowAuthority guards. That case already
// forbids allow/approve/grant-shaped fields on tools.Decision and mode/allow/
// approve/auto-shaped fields on agent.ToolRequest; it does NOT look at the one
// shape this ticket could get wrong, so this case looks at it and nothing else.
//
// The forbidden shape is a SESSION field on the caller's request: Execute's caller
// (the loop, a panel host, a future spill re-read) naming which session it is
// asking as is the M-7/C-3 bug again - a security decision keyed on a
// caller-controlled selector. The identity lives on the injected source, which is
// a construction-time act, and nothing else.
func TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority(t *testing.T) {
	rTyp := reflect.TypeOf(agent.ToolRequest{})
	for i := 0; i < rTyp.NumField(); i++ {
		name := rTyp.Field(i).Name
		low := strings.ToLower(name)
		if strings.Contains(low, "session") || strings.Contains(low, "grant") {
			t.Errorf("agent.ToolRequest field %q lets the CALLER name the session or hand in "+
				"a grant: the identity has exactly one owner (the host minted it, "+
				"internal/session/session.go) and Execute's caller supplies arguments", name)
		}
	}

	// And the seam stays read-only: one method, no revoke, no setter.
	src := reflect.TypeOf((*GrantSource)(nil)).Elem()
	if src.NumMethod() != 1 {
		t.Errorf("tools.GrantSource has %d methods, want exactly 1 (read-only, like ModeSource)",
			src.NumMethod())
	}
	if src.NumMethod() > 0 {
		if name := src.Method(0).Name; name != "Covering" {
			t.Errorf("tools.GrantSource's only method is %s, want Covering", name)
		}
	}
	for i := 0; i < src.NumMethod(); i++ {
		n := strings.ToLower(src.Method(i).Name)
		for _, bad := range []string{"record", "revoke", "set", "insert", "add", "delete"} {
			if strings.Contains(n, bad) {
				t.Errorf("tools.GrantSource method %q writes: the read side of an authorization "+
					"must not be able to create or remove one", src.Method(i).Name)
			}
		}
	}

	// The bridge must not own a setter either - same rule modes already have.
	bTyp := reflect.TypeOf(Bridge{})
	for i := 0; i < bTyp.NumField(); i++ {
		if strings.HasPrefix(strings.ToLower(bTyp.Field(i).Name), "grant") &&
			bTyp.Field(i).Type.Kind() == reflect.Func {
			t.Errorf("Bridge field %q is a function-typed grant channel: the bridge may read a "+
				"grant, never install one", bTyp.Field(i).Name)
		}
	}
}

// grantsOf pulls the fixture source back out of a bridge. It exists because two
// of the cases below need the bridge's canonicalizer to compute the path a grant
// must be stored as, and the source has to be handed in at construction: build
// first, then seed through the same resolver the judge used.
//
// It reaches for the unexported field rather than taking a parameter because this
// file is in package tools and the alternative is a setter, which is exactly the
// shape TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority asserts does not
// exist.
func grantsOf(t *testing.T, b *Bridge) *t224Grants {
	t.Helper()
	g, ok := b.grants.(*t224Grants)
	if !ok {
		t.Fatalf("grantsOf: bridge.grants is %T, want *t224Grants", b.grants)
	}
	return g
}
