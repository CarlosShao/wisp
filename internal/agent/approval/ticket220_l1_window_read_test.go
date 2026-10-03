package approval

// Ticket 220 AC#2's read-side cases for LiveL1Windows (window_read.go), the
// enumeration the L1 short window never had: gate.go:302-316 opens and closes
// g.windows, and until this ticket the only reader of that map by name was Veto
// (:377ff), which needs the id it is looking for. So "which tasks are waiting in
// an L1 window right now" had no producer anywhere, and the panel's roster cell
// could only ever see the L2 queue.
//
// Why this file is `package approval` (like pending_read_test.go, not like the
// eight black-box files): AC#2's ruling bounds the new name to READ-ONLY -
// "which correlation/task is waiting", nothing else. That bound is only
// measurable from inside the package, because the state it must not touch
// (g.windows' own entries, g.running, g.admitted, g.grants) is unexported. A
// black-box case can show the output is stable; only a white-box case can show
// that producing it moved nothing.

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// a UI that holds a window open, and a clock that never lets it expire
// ---------------------------------------------------------------------------

// holdUI is the window's presentation surface for these cases: Prompt reports the
// prompt and then parks until the test releases it, so the window stays in
// g.windows exactly as long as the test says. It counts Update calls because one
// of AC#2's bounds is that reading the enumeration emits nothing.
type holdUI struct {
	mu       sync.Mutex
	prompted chan Prompt
	release  chan struct{}
	updates  int
	released sync.Once
}

func newHoldUI() *holdUI {
	return &holdUI{prompted: make(chan Prompt, 4), release: make(chan struct{})}
}

func (u *holdUI) Prompt(ctx context.Context, p Prompt) error {
	select {
	case u.prompted <- p:
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case <-u.release:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (u *holdUI) Update(_ context.Context, _ Event) error {
	u.mu.Lock()
	u.updates++
	u.mu.Unlock()
	return nil
}

func (u *holdUI) updateCount() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.updates
}

func (u *holdUI) releaseWindow() { u.released.Do(func() { close(u.release) }) }

// neverClock fires nothing: the L1 window's own deadline is a real 2-3s budget,
// and a case that parks a window on purpose must not be racing it (the repo ban
// on wall-clock timeouts applies to the layer under test, not to how it is
// pinned).
type neverClock struct{ now time.Time }

func (c neverClock) Now() time.Time { return c.now }

func (c neverClock) After(time.Duration) <-chan time.Time {
	return make(chan time.Time)
}

// countingGrants is the D45 write side handed to the gate; its count is what the
// "grants unchanged" half of the read-only bound reads.
type countingGrants struct {
	mu      sync.Mutex
	records int
}

func (r *countingGrants) Record(_ context.Context, _, _ string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records++
	return 1, nil
}

func (r *countingGrants) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.records
}

// ---------------------------------------------------------------------------
// state snapshot, read from inside the package
// ---------------------------------------------------------------------------

type gateStateProbe struct {
	windowKeys   []string
	windowItems  []*window
	runningKeys  []string
	admittedKeys []string
	pendingCorrs []string
	queueDepth   int
	grants       int
	uiUpdates    int
}

// probeState reads everything the new enumeration is forbidden to move. It takes
// the gate's own lock for the maps it reads, so a concurrent answer cannot make
// the diff lie.
func (g *Gate) probeState(grants *countingGrants, ui *holdUI) gateStateProbe {
	var p gateStateProbe
	g.mu.Lock()
	for k, w := range g.windows {
		p.windowKeys = append(p.windowKeys, k)
		p.windowItems = append(p.windowItems, w)
	}
	for k := range g.running {
		p.runningKeys = append(p.runningKeys, k)
	}
	for k := range g.admitted {
		p.admittedKeys = append(p.admittedKeys, k)
	}
	g.mu.Unlock()
	sort.Strings(p.windowKeys)
	sort.Strings(p.runningKeys)
	sort.Strings(p.admittedKeys)

	g.q.mu.Lock()
	p.queueDepth = len(g.q.pending)
	for _, it := range g.q.pending {
		p.pendingCorrs = append(p.pendingCorrs, it.Corr)
	}
	g.q.mu.Unlock()
	sort.Strings(p.pendingCorrs)

	p.grants = grants.count()
	p.uiUpdates = ui.updateCount()
	return p
}

func sameProbe(a, b gateStateProbe) bool {
	return reflect.DeepEqual(a.windowKeys, b.windowKeys) &&
		reflect.DeepEqual(a.windowItems, b.windowItems) &&
		reflect.DeepEqual(a.runningKeys, b.runningKeys) &&
		reflect.DeepEqual(a.admittedKeys, b.admittedKeys) &&
		reflect.DeepEqual(a.pendingCorrs, b.pendingCorrs) &&
		a.queueDepth == b.queueDepth &&
		a.grants == b.grants &&
		a.uiUpdates == b.uiUpdates
}

// ---------------------------------------------------------------------------
// cases
// ---------------------------------------------------------------------------

func windowDecision(corr, task string) tools.Decision {
	return tools.Decision{
		Tool: "fs.write", Provider: tools.KindBuiltin,
		Params:        map[string]any{"paths": []any{"C:/data/a.txt"}},
		Args:          json.RawMessage(`{"paths":["C:/data/a.txt"]}`),
		Level:         risk.L1,
		RulesHit:      []risk.RuleID{risk.R1},
		Reason:        "声明 L1：可逆写",
		Paths:         []string{"C:/data/a.txt"},
		CorrelationID: corr,
		TaskID:        task,
		CallID:        "call-" + corr,
	}
}

// holdOpenWindow starts one L1 window, waits until its prompt has been delivered,
// and returns the gate plus the answer channel. The window is left open on
// purpose: releasing it is the case's last step, never a cleanup the case relies
// on for its reading.
func holdOpenWindow(t *testing.T, corr, task string) (*Gate, *holdUI, *countingGrants, <-chan struct {
	a   tools.Answer
	why string
},
) {
	t.Helper()
	ui := newHoldUI()
	grants := &countingGrants{}
	g := New(Options{UI: ui, Clock: neverClock{now: time.Unix(0, 0)}, Window: 3 * time.Second, Grants: grants})
	g.AdmitTextTask(task)

	out := make(chan struct {
		a   tools.Answer
		why string
	}, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("window goroutine panicked: %v", r)
			}
		}()
		a, why := g.PendingWindow(context.Background(), windowDecision(corr, task))
		out <- struct {
			a   tools.Answer
			why string
		}{a, why}
	}()
	t.Cleanup(func() {
		ui.releaseWindow()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Errorf("派生的窗口协程未在 2s 内退出（闸门存在无限等待）")
		}
	})

	select {
	case p := <-ui.prompted:
		if p.CorrelationID != corr || p.Level != "L1" {
			t.Fatalf("prompt = %+v, want the L1 window for %q", p, corr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("确认界面从未收到提示（等待 2s 后放弃）")
	}
	return g, ui, grants, out
}

// TestLiveL1WindowsSeesAWindowInTheLoop names the gap AC#2 closes: a task waiting
// inside the 2-3s L1 block is now enumerable, under both ids the join needs and
// the tool it is waiting on. Before window_read.go this call did not exist, and
// the only producer of "what is waiting" was the L2 queue.
func TestLiveL1WindowsSeesAWindowInTheLoop(t *testing.T) {
	g, ui, _, out := holdOpenWindow(t, "corr-220", "task-220")

	if got := g.LiveL1Windows(); len(got) != 1 {
		t.Fatalf("LiveL1Windows() = %+v, want exactly the one window in the loop", got)
	}
	w := got0(t, g.LiveL1Windows())
	if w.CorrelationID != "corr-220" || w.TaskID != "task-220" || w.Tool != "fs.write" {
		t.Errorf("reported window = %+v, want corr-220 / task-220 / fs.write verbatim", w)
	}

	// A second, unrelated approval route is untouched by the enumeration: an L1
	// window is not a queue item, and the enumeration must not quietly make it one
	// (that would move cards into the section the user answers, ticket 219's
	// buttons included).
	if items := g.Queue().LiveApprovals(); len(items) != 0 {
		t.Errorf("LiveApprovals() = %+v while only an L1 window is open; an L1 window "+
			"is not a queue item and the enumeration must not publish it as one", items)
	}

	// And once the window is gone the dimension goes back to empty: it reports
	// "waiting", not "this correlation was ever seen".
	if err := g.Veto(Veto{CorrelationID: "corr-220", Channel: ChannelBall}); err != nil {
		t.Fatalf("Veto after the enumeration read the window: %v", err)
	}
	// The gate is parked in Prompt by design, so the veto is queued and the select
	// loop only sees it once the prompt is released - releasing here is the test
	// letting the window run its own veto path, never a way to answer it.
	ui.releaseWindow()
	select {
	case a := <-out:
		if a.a != tools.AnswerVeto {
			t.Errorf("answer = %q (%s), want veto - the enumeration must not have "+
				"consumed or re-answered the window it just read", a.a, a.why)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("否决后闸门未在 2s 内给出答复")
	}
	if got := g.LiveL1Windows(); len(got) != 0 {
		t.Errorf("after the window closed LiveL1Windows() = %+v, want empty", got)
	}
}

// TestReadingL1WindowsMovesNoGateState is AC#2's bound as an instrument: the new
// name may read which ids are waiting and nothing else. Every half of the gate's
// own state - the window entries themselves (pointer identity, not just keys),
// the D31 running set, D47's admission registry, the L2 queue, the D45 grant
// recorder, the UI's event stream - is compared across the call.
func TestReadingL1WindowsMovesNoGateState(t *testing.T) {
	g, ui, grants, out := holdOpenWindow(t, "corr-220b", "task-220b")

	before := g.probeState(grants, ui)
	first := g.LiveL1Windows()
	second := g.LiveL1Windows()
	after := g.probeState(grants, ui)

	if !sameProbe(before, after) {
		t.Errorf("reading the enumeration moved gate state:\n before=%+v\n after =%+v", before, after)
	}
	if !reflect.DeepEqual(first, second) {
		t.Errorf("two calls disagree: %+v vs %+v (and an unordered reply would make the "+
			"packet's bytes unreproducible, which the ledger hashes)", first, second)
	}
	if len(g.Queue().LiveApprovals()) != 0 {
		t.Error("the enumeration created a pending L2 item")
	}

	// The window the enumeration read twice is still answerable exactly once.
	if err := g.Veto(Veto{CorrelationID: "corr-220b", Channel: ChannelEsc}); err != nil {
		t.Fatalf("Veto after two reads: %v", err)
	}
	ui.releaseWindow()
	select {
	case a := <-out:
		if a.a != tools.AnswerVeto {
			t.Errorf("answer=%q, want veto", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("否决后闸门未在 2s 内给出答复")
	}
	if got := g.probeState(grants, ui); len(got.windowKeys) != 0 {
		t.Errorf("the window outlived its veto: %+v", got.windowKeys)
	}
}

// TestL1WindowIsPlainDataWithNoRouteToAnAnswer pins the shape of the bound rather
// than one behaviour of it: what comes back is three id strings and no method. A
// type that could allow, veto, or answer would have to be a method on it (or on
// the gate through it), and this is the case that goes red the day someone adds
// one.
func TestL1WindowIsPlainDataWithNoRouteToAnAnswer(t *testing.T) {
	typ := reflect.TypeOf(L1Window{})
	if typ.Kind() != reflect.Struct {
		t.Fatalf("L1Window is %v, want plain struct data", typ.Kind())
	}
	if typ.NumMethod() != 0 {
		t.Errorf("L1Window has %d method(s); a reported window must carry no behaviour "+
			"at all - it is one dimension of a read, not a door", typ.NumMethod())
	}
	var names []string
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		names = append(names, f.Name)
		if f.Type.Kind() != reflect.String {
			t.Errorf("field %s is %v; only id strings may travel out of g.windows - a "+
				"channel, context, or cancel func would turn the read side into a veto door",
				f.Name, f.Type)
		}
	}
	sort.Strings(names)
	if want := []string{"CorrelationID", "TaskID", "Tool"}; !reflect.DeepEqual(names, want) {
		t.Errorf("L1Window fields = %v, want exactly %v (the one dimension AC#2 was approved for)", names, want)
	}
	gateType := reflect.TypeOf(&Gate{})
	for i := 0; i < gateType.NumMethod(); i++ {
		m := gateType.Method(i)
		if m.Type.NumOut() > 0 {
			if m.Type.Out(m.Type.NumOut()-1) == typ ||
				m.Type.Out(m.Type.NumOut()-1) == reflect.SliceOf(typ) {
				// Exactly one method may hand this type out, and it must be the read.
				if m.Name != "LiveL1Windows" {
					t.Errorf("%s returns the window type; the enumeration must be the only "+
						"producer of it", m.Name)
				}
			}
		}
	}
}

// TestStartedCallsAreNotReportedAsWaiting is the other direction of the same
// bound: a call that already went live is not "waiting", and lighting the roster
// cell for it would be the opposite lie (「需要批准」 on something already writing).
func TestStartedCallsAreNotReportedAsWaiting(t *testing.T) {
	ui := newHoldUI()
	grants := &countingGrants{}
	g := New(Options{UI: ui, Clock: neverClock{now: time.Unix(0, 0)}, Window: 3 * time.Second, Grants: grants})
	g.markStarted("corr-running", "fs.write")
	if got := g.LiveL1Windows(); len(got) != 0 {
		t.Errorf("a call that already started is reported as waiting: %+v", got)
	}
	if len(g.Queue().LiveApprovals()) != 0 {
		t.Error("markStarted must not create a pending card either")
	}
}

// got0 is the "exactly one, or the case has no premise" reader.
func got0(t *testing.T, got []L1Window) L1Window {
	t.Helper()
	if len(got) != 1 {
		t.Fatalf("LiveL1Windows() = %+v, want exactly one entry", got)
	}
	return got[0]
}
