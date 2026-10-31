package watchdog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// --- fixtures -------------------------------------------------------------

// fakeTree feeds the loop scripted tree metrics; the fixtures mutate `cur`
// to model leaks that grow and drops that shrink the tree.
type fakeTree struct {
	mu  sync.Mutex
	cur observe.TreeMetrics
}

func (f *fakeTree) ReadTree() (observe.TreeMetrics, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.cur, nil
}

func (f *fakeTree) set(m observe.TreeMetrics) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cur = m
}

// fakeComponent is a fixture Unloadable. dropped/detail model the two
// honest outcomes; panicWith turns the next Unload into a panic (the
// panic-isolation fixture); onUnload runs the drop's side effect (the
// fixture tree shrinking when the component lets go).
type fakeComponent struct {
	mu        sync.Mutex
	calls     int
	dropped   bool
	detail    string
	panicWith any
	onUnload  func()
}

func (c *fakeComponent) Unload(ctx context.Context) (bool, string) {
	c.mu.Lock()
	c.calls++
	pw, detail, drop, onUnload := c.panicWith, c.detail, c.dropped, c.onUnload
	c.mu.Unlock()
	if pw != nil {
		panic(pw)
	}
	if onUnload != nil {
		onUnload()
	}
	return drop, detail
}

func (c *fakeComponent) unloadCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls
}

func captureLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(h), &buf
}

func overCapTree(bytes int64) observe.TreeMetrics {
	return observe.TreeMetrics{
		PIDs:                   1,
		PrivateWorkingSetBytes: bytes,
		CommitBytes:            bytes,
		Handles:                100,
		GDIObjects:             50,
		USERObjects:            40,
		Threads:                20,
	}
}

type alertLog struct {
	mu     sync.Mutex
	alerts []Alert
}

func (a *alertLog) add(al Alert) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.alerts = append(a.alerts, al)
}

func (a *alertLog) of(kind AlertKind) []Alert {
	a.mu.Lock()
	defer a.mu.Unlock()
	var out []Alert
	for _, al := range a.alerts {
		if al.Kind == kind {
			out = append(out, al)
		}
	}
	return out
}

// logEntries parses the captured JSONL log into maps, one per line.
func logEntries(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	dec := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	for dec.More() {
		var m map[string]any
		if err := dec.Decode(&m); err != nil {
			t.Fatalf("log line is not JSON: %v\nbuffer:\n%s", err, buf.String())
		}
		out = append(out, m)
	}
	return out
}

func logMsgs(entries []map[string]any, msg string) []map[string]any {
	var out []map[string]any
	for _, e := range entries {
		if e["msg"] == msg {
			out = append(out, e)
		}
	}
	return out
}

// --- AC#1 (Sleeping half): sustained over-cap drops unloadables -----------

// TestSleepingOverCapDropsUnloadablesThenWithinCap drives the D32 revision 3
// Sleeping row end to end: the fixture leak is over the 25MB cap for the
// sustain window, the pass drops the panel WebView (the user is not viewing
// it) and the BM25 cache, and the tree re-read lands within cap - so no
// alert. Non-whitelisted registrations are never called (AC#4's structural
// half is asserted here too).
func TestSleepingOverCapDropsUnloadablesThenWithinCap(t *testing.T) {
	const capB = int64(25 << 20)
	tree := &fakeTree{}
	tree.set(overCapTree(30 << 20)) // over the Sleeping cap

	panel := &fakeComponent{dropped: true, detail: "panel hidden for unload (user not viewing)"}
	bm25 := &fakeComponent{dropped: true, detail: "bm25 index cache evicted"}
	kws := &fakeComponent{dropped: true, detail: "MUST NEVER BE CALLED"}
	dbw := &fakeComponent{dropped: true, detail: "MUST NEVER BE CALLED"}

	log, buf := captureLogger()
	als := &alertLog{}
	l := New(Options{
		TreeReader:   tree,
		Registry:     observe.NewRegistry(),
		StateFn:      func() observe.SLOState { return observe.SLOSleeping },
		OnAlert:      als.add,
		Log:          log,
		TickEvery:    time.Millisecond,
		SustainTicks: 3,
	})
	l.Register(ComponentPanelWebView, panel)
	l.Register(ComponentBM25IndexCache, bm25)
	l.Register("kws", kws)              // never whitelisted: D32 revision 3
	l.Register("sqlite-db-writer", dbw) // never whitelisted: D32 revision 3

	// The drops shrink the tree below the cap: the fixture leak is "fixed"
	// by dropping, which is what the real owners' Unload does.
	panel.onUnload = func() { tree.set(overCapTree(20 << 20)) }
	bm25.onUnload = func() { tree.set(overCapTree(20 << 20)) }

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		l.tick(ctx)
	}

	if got := panel.unloadCalls(); got != 1 {
		t.Fatalf("panel-webview unload calls = %d, want 1 (one pass over the sustained breach)", got)
	}
	if got := bm25.unloadCalls(); got != 1 {
		t.Errorf("bm25-index-cache unload calls = %d, want 1", got)
	}
	if got := kws.unloadCalls(); got != 0 {
		t.Errorf("kws unload calls = %d, want 0: a non-whitelisted component was unloaded", got)
	}
	if got := dbw.unloadCalls(); got != 0 {
		t.Errorf("sqlite-db-writer unload calls = %d, want 0", got)
	}
	if n := len(als.alerts); n != 0 {
		t.Errorf("alerts = %d, want 0: the pass brought the tree within cap, so the episode closes without an alert", n)
	}

	entries := logEntries(t, buf)
	passes := logMsgs(entries, "watchdog: unload pass complete")
	if len(passes) != 1 {
		t.Fatalf("unload pass complete lines = %d, want 1", len(passes))
	}
	p := passes[0]
	if before, _ := p["before_bytes"].(float64); before != float64(30<<20) {
		t.Errorf("before_bytes = %v, want the over-cap 30MB reading", p["before_bytes"])
	}
	if after, _ := p["after_bytes"].(float64); after != float64(20<<20) {
		t.Errorf("after_bytes = %v, want 20MB - but the fixture never shrank the tree; wire the drop hook", p["after_bytes"])
	}
	actions := logMsgs(entries, "watchdog: unload action")
	dropped := map[string]bool{}
	for _, a := range actions {
		if a["outcome"] == "dropped" {
			dropped[a["component"].(string)] = true
		}
	}
	if !dropped[ComponentPanelWebView] || !dropped[ComponentBM25IndexCache] {
		t.Errorf("dropped components = %v, want panel-webview and bm25-index-cache", dropped)
	}
}

// TestSleepingOverCapStillOverAfterPassAlertsOnce is the other arm of
// "within cap or alert": the whitelist components all skip (panel viewed,
// caches busy), the tree stays over cap, and the overrun alert fires once
// per episode no matter how long the breach continues.
func TestSleepingOverCapStillOverAfterPassAlertsOnce(t *testing.T) {
	tree := &fakeTree{}
	tree.set(overCapTree(30 << 20)) // over the Sleeping cap, never drops

	panel := &fakeComponent{dropped: false, detail: "user is viewing the panel"}
	bm25 := &fakeComponent{dropped: false, detail: "cache hot"}

	log, _ := captureLogger()
	als := &alertLog{}
	l := New(Options{
		TreeReader:   tree,
		Registry:     observe.NewRegistry(),
		StateFn:      func() observe.SLOState { return observe.SLOSleeping },
		OnAlert:      als.add,
		Log:          log,
		TickEvery:    time.Millisecond,
		SustainTicks: 2,
	})
	l.Register(ComponentPanelWebView, panel)
	l.Register(ComponentBM25IndexCache, bm25)

	ctx := context.Background()
	for i := 0; i < 8; i++ { // four passes worth of ticks
		l.tick(ctx)
	}

	if panel.unloadCalls() != 4 || bm25.unloadCalls() != 4 {
		t.Errorf("each pass must consult both components: panel=%d bm25=%d, want 4/4",
			panel.unloadCalls(), bm25.unloadCalls())
	}
	overruns := als.of(AlertOverrun)
	if len(overruns) != 1 {
		t.Fatalf("overrun alerts = %d, want 1 (deduped per episode)", len(overruns))
	}
	a := overruns[0]
	if a.State != observe.SLOSleeping || a.Metric != "tree_private_bytes" {
		t.Errorf("overrun alert = %s/%s, want Sleeping/tree_private_bytes", a.State, a.Metric)
	}
	if a.Measured == "" || a.Limit == "" {
		t.Errorf("overrun alert must carry measured (%q) and limit (%q)", a.Measured, a.Limit)
	}
}

// --- AC#1 (Armed half): KWS RETAINED + warning (rule-critical) -------------

// TestArmedOverCapRetainsKWSAndWarns is the rule-critical test: in Armed the
// sustained breach must NOT unload KWS (it is not whitelisted - the pass
// still drops whitelisted caches), and the keep-alive warning with the
// disable-KWS suggestion goes out. The host's D43 row 9 mapping keeps the
// machine in Armed.
func TestArmedOverCapRetainsKWSAndWarns(t *testing.T) {
	tree := &fakeTree{}
	tree.set(overCapTree(100 << 20)) // over the Armed cap (90MB)

	kws := &fakeComponent{dropped: true, detail: "MUST NEVER BE CALLED"}
	bm25 := &fakeComponent{dropped: true, detail: "bm25 index cache evicted"}

	log, buf := captureLogger()
	als := &alertLog{}
	l := New(Options{
		TreeReader:   tree,
		Registry:     observe.NewRegistry(),
		StateFn:      func() observe.SLOState { return observe.SLOArmed },
		OnAlert:      als.add,
		Log:          log,
		TickEvery:    time.Millisecond,
		SustainTicks: 3,
	})
	l.Register("kws", kws)
	l.Register(ComponentBM25IndexCache, bm25)

	ctx := context.Background()
	for i := 0; i < 6; i++ { // two sustain windows inside one over-cap episode
		l.tick(ctx)
	}

	if got := kws.unloadCalls(); got != 0 {
		t.Fatalf("KWS unload calls = %d, want 0: unloading KWS in Armed is the D32 revision 3 harmful behavior", got)
	}
	if bm25.unloadCalls() != 2 {
		t.Errorf("bm25 unload calls = %d, want 2 (the pass still drops whitelisted components in Armed)", bm25.unloadCalls())
	}
	overruns := als.of(AlertOverrun)
	if len(overruns) != 1 {
		t.Fatalf("overrun alerts = %d, want 1: the Armed warning is the state's overrun response, deduped per episode", len(overruns))
	}
	if !strings.Contains(overruns[0].Detail, "KWS") || !strings.Contains(overruns[0].Detail, "disabling") {
		t.Errorf("Armed warning must name KWS retention and the disable suggestion, got: %q", overruns[0].Detail)
	}
	if !strings.Contains(buf.String(), "outcome\":\"dropped") || !strings.Contains(buf.String(), "bm25-index-cache") {
		t.Errorf("the Armed pass must still log its dropped whitelisted component; log:\n%s", buf.String())
	}
}

// TestArmedWarningMapsToD43Row9 stays Armed: the host maps AlertOverrun in
// Armed to EvWatchdogOverrun (D43 row 9), whose side effect is the
// keep-alive alert and whose destination is Armed itself.
func TestArmedWarningMapsToD43Row9(t *testing.T) {
	// The mapping this test pins is the host seam: alert -> D43 event. The
	// machine side is frozen (table.go row 9); what the watchdog guarantees
	// is that the Armed warning arrives as AlertOverrun on State Armed.
	tree := &fakeTree{}
	tree.set(overCapTree(100 << 20))
	als := &alertLog{}
	l := New(Options{
		TreeReader:   tree,
		Registry:     observe.NewRegistry(),
		StateFn:      func() observe.SLOState { return observe.SLOArmed },
		OnAlert:      als.add,
		TickEvery:    time.Millisecond,
		SustainTicks: 1,
	})
	l.tick(context.Background())
	overruns := als.of(AlertOverrun)
	if len(overruns) != 1 {
		t.Fatalf("overrun alerts = %d, want 1", len(overruns))
	}
	if overruns[0].Kind != AlertOverrun || overruns[0].State != observe.SLOArmed {
		t.Errorf("the host mapping needs AlertOverrun on Armed, got %s/%s",
			overruns[0].Kind, overruns[0].State)
	}
}

// --- AC#3: leak indicators -> alert + snapshot -----------------------------

// TestLeakIndicatorsAlertWithSnapshot crosses the handle, GDI and goroutine
// caps with fixtures, asserts one alert each with a non-empty diagnostic
// snapshot, dedup per episode, re-arm after the metric falls back, and that
// leak alerts never unload anything.
func TestLeakIndicatorsAlertWithSnapshot(t *testing.T) {
	tree := &fakeTree{}
	base := overCapTree(10 << 20) // within the Sleeping memory cap
	base.Handles = 601            // over HandleLimit (600)
	base.GDIObjects = 201         // over GDILimit (200)
	tree.set(base)

	reg := observe.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < 7; i++ { // over the frozen Sleeping goroutine cap (6)
		reg.Spawn("kws-infer", "test", nil, func(c context.Context) { <-c.Done() })
	}
	deadline := observe.NewTimeout(2 * time.Second)
	for reg.Count() < 7 {
		if deadline.Expired() {
			t.Fatalf("goroutine fixtures did not come up: count=%d", reg.Count())
		}
		time.Sleep(time.Millisecond)
	}

	log, buf := captureLogger()
	als := &alertLog{}
	l := New(Options{
		TreeReader: tree,
		Registry:   reg,
		StateFn:    func() observe.SLOState { return observe.SLOSleeping },
		OnAlert:    als.add,
		Log:        log,
		TickEvery:  time.Millisecond,
	})
	// A registered component: leak alerts must not touch it.
	panel := &fakeComponent{dropped: true, detail: "must stay"}
	l.Register(ComponentPanelWebView, panel)

	l.tick(ctx) // crossing tick
	if got := len(als.alerts); got != 3 {
		t.Fatalf("alerts = %d, want 3 (handles, gdi_objects, goroutines)", got)
	}
	want := map[string]int{"handles": 601, "gdi_objects": 201, "goroutines": 7}
	for _, a := range als.alerts {
		if a.Kind != AlertLeak {
			t.Errorf("alert kind = %s, want leak", a.Kind)
		}
		if a.Class != observe.ClassResource {
			t.Errorf("alert %s class = %q, want resource (D37)", a.Metric, a.Class)
		}
		if a.Measured != fmt.Sprint(want[a.Metric]) {
			t.Errorf("alert %s measured = %q, want %d", a.Metric, a.Measured, want[a.Metric])
		}
		if len(a.Snapshot) == 0 {
			t.Errorf("alert %s carries no diagnostic snapshot", a.Metric)
			continue
		}
		var snap map[string]any
		if err := json.Unmarshal(a.Snapshot, &snap); err != nil {
			t.Errorf("alert %s snapshot is not JSON: %v", a.Metric, err)
			continue
		}
		if snap["metric"] != a.Metric || snap["state"] != string(observe.SLOSleeping) {
			t.Errorf("snapshot identity = %v/%v, want %s/Sleeping", snap["metric"], snap["state"], a.Metric)
		}
		if _, ok := snap["tree"].(map[string]any); !ok {
			t.Errorf("snapshot must carry the crossing tree sample; got keys %v", snap)
		}
	}
	if panel.unloadCalls() != 0 {
		t.Errorf("leak indicators unloaded a component (%d calls): D42#10 forbids that", panel.unloadCalls())
	}

	l.tick(ctx) // same episode: deduped
	l.tick(ctx)
	if got := len(als.alerts); got != 3 {
		t.Fatalf("alerts after dedup ticks = %d, want 3", got)
	}

	base.Handles = 100 // falls back under
	tree.set(base)
	l.tick(ctx) // re-arms handles
	base.Handles = 700
	tree.set(base)
	l.tick(ctx) // crosses again
	leaks := als.of(AlertLeak)
	handles := 0
	for _, a := range leaks {
		if a.Metric == "handles" {
			handles++
		}
	}
	if handles != 2 {
		t.Errorf("handles leak alerts = %d, want 2 (one per episode, re-armed after falling back)", handles)
	}
	if !strings.Contains(buf.String(), "no auto-unload") {
		t.Errorf("the leak log line must state that nothing auto-unloads; log:\n%s", buf.String())
	}
	cancel()
}

// --- AC#4: log audit - no unload of non-whitelisted components, ever -------

// TestLogAuditNoUnloadOfNonWhitelistedEver registers the full whitelist plus
// every NEVER component the ticket names, drives sustained over-cap breaches
// in two states, and audits the structured log: every unload action names a
// whitelisted component, and no non-whitelisted name appears in any unload
// line in any outcome.
func TestLogAuditNoUnloadOfNonWhitelistedEver(t *testing.T) {
	tree := &fakeTree{}
	tree.set(overCapTree(400 << 20)) // over Warm and Sleeping caps alike

	comps := map[string]*fakeComponent{}
	for _, name := range whitelist {
		comps[name] = &fakeComponent{dropped: true, detail: "dropped"}
	}
	never := []string{"kws", "audio-capture", "sqlite-db-writer", "session-engines", "agent-context"}
	for _, name := range never {
		comps[name] = &fakeComponent{dropped: true, detail: "MUST NEVER BE CALLED"}
	}

	log, buf := captureLogger()
	als := &alertLog{}
	state := observe.SLOSleeping
	l := New(Options{
		TreeReader:   tree,
		Registry:     observe.NewRegistry(),
		StateFn:      func() observe.SLOState { return state },
		OnAlert:      als.add,
		Log:          log,
		TickEvery:    time.Millisecond,
		SustainTicks: 1,
	})
	for name, c := range comps {
		l.Register(name, c)
	}

	ctx := context.Background()
	states := []observe.SLOState{observe.SLOSleeping, observe.SLOWarm, observe.SLOArmed, observe.SLOConversation}
	for _, st := range states {
		state = st
		l.tick(ctx) // one sustain window per state -> one pass per state
	}

	entries := logEntries(t, buf)
	actions := logMsgs(entries, "watchdog: unload action")
	if len(actions) == 0 {
		t.Fatalf("no unload action lines found; the audit has nothing to read (log):\n%s", buf.String())
	}
	whitelisted := map[string]bool{}
	for _, name := range whitelist {
		whitelisted[name] = true
	}
	for _, a := range actions {
		name, _ := a["component"].(string)
		if !whitelisted[name] {
			t.Errorf("unload action for non-whitelisted component %q (outcome %v): the whitelist is the only unload set", name, a["outcome"])
		}
		if a["outcome"] != "dropped" && a["outcome"] != "skipped" {
			t.Errorf("unload action %q carries outcome %v, want dropped or skipped", name, a["outcome"])
		}
	}
	for _, name := range never {
		for _, a := range actions {
			if a["component"] == name {
				t.Errorf("%q appears in an unload action line: the audit forbids it", name)
			}
		}
		if got := comps[name].unloadCalls(); got != 0 {
			t.Errorf("%q Unload called %d times across %d states: never unload a non-whitelisted component", name, got, len(states))
		}
	}
	passes := logMsgs(entries, "watchdog: unload pass complete")
	if len(passes) != len(states) {
		t.Errorf("unload pass lines = %d, want %d (one per state)", len(passes), len(states))
	}
	for _, p := range passes {
		for _, key := range []string{"before_bytes", "after_bytes", "cap_bytes", "dropped"} {
			if _, ok := p[key]; !ok {
				t.Errorf("unload pass line is missing %q: actions must be logged with metrics", key)
			}
		}
	}
}

// --- panic isolation --------------------------------------------------------

// TestUnloadPanicIsolated proves a panicking component neither kills the
// pass, the loop, nor the other components: the pass logs the error outcome
// and keeps dropping the rest, and the loop keeps ticking afterwards.
func TestUnloadPanicIsolated(t *testing.T) {
	tree := &fakeTree{}
	tree.set(overCapTree(30 << 20))

	panel := &fakeComponent{dropped: true, panicWith: "boom: panel host exploded"}
	bm25 := &fakeComponent{dropped: true, detail: "evicted"}

	log, buf := captureLogger()
	als := &alertLog{}
	l := New(Options{
		TreeReader:   tree,
		Registry:     observe.NewRegistry(),
		StateFn:      func() observe.SLOState { return observe.SLOSleeping },
		OnAlert:      als.add,
		Log:          log,
		TickEvery:    time.Millisecond,
		SustainTicks: 1,
	})
	l.Register(ComponentPanelWebView, panel)
	l.Register(ComponentBM25IndexCache, bm25)

	ctx := context.Background()
	for i := 0; i < 2; i++ { // two passes; the panicking component panics twice
		l.tick(ctx)
	}

	if bm25.unloadCalls() != 2 {
		t.Errorf("bm25 unload calls = %d, want 2: the panic must not stop the rest of the pass", bm25.unloadCalls())
	}
	if panel.unloadCalls() != 2 {
		t.Errorf("panel unload calls = %d, want 2: the loop must survive and consult it again", panel.unloadCalls())
	}
	entries := logEntries(t, buf)
	if fails := logMsgs(entries, "watchdog: unload action failed"); len(fails) != 2 {
		t.Errorf("unload-failure lines = %d, want 2; log:\n%s", len(fails), buf.String())
	}
	if !strings.Contains(buf.String(), "panicked") {
		t.Errorf("the failure line must name the panic; log:\n%s", buf.String())
	}
	// The loop itself is alive: the next tick still runs the full body.
	tree.set(overCapTree(20 << 20)) // within cap
	l.tick(ctx)
}

// --- the named goroutine ----------------------------------------------------

// TestStartRunsOneNamedWatchdogGoroutineAndExitsOnCancel proves the loop
// runs as the D38b resident "watchdog" goroutine (named, counted by the
// registry) and exits cleanly on ctx cancel with no panic record.
func TestStartRunsOneNamedWatchdogGoroutineAndExitsOnCancel(t *testing.T) {
	tree := &fakeTree{}
	tree.set(overCapTree(10 << 20))
	reg := observe.NewRegistry()
	l := New(Options{
		TreeReader: tree,
		Registry:   reg,
		TickEvery:  2 * time.Millisecond,
	})
	ctx, cancel := context.WithCancel(context.Background())
	h := l.Start(ctx)
	deadline := observe.NewTimeout(5 * time.Second)
	for reg.CountByName("watchdog") != 1 {
		if deadline.Expired() {
			t.Fatalf("the watchdog goroutine never appeared in the registry")
		}
		time.Sleep(time.Millisecond)
	}
	cancel()
	select {
	case <-h.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("the loop did not exit after ctx cancel")
	}
	if err := h.Err(); err != nil {
		t.Fatalf("loop goroutine ended with %v", err)
	}
	if got := reg.CountByName("watchdog"); got != 0 {
		t.Errorf("watchdog goroutines still registered after exit: %d", got)
	}
}
