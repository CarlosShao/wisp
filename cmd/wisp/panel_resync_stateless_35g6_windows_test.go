//go:build windows

package main

// Ticket 35 cell 票面 :49 ("Resync", today's line number) - the behavioural
// teeth the dispatch order names:
//
//	① one test rides the REAL bridge (not a mock internal function): hide,
//	  then the tool state changes, then show - and the page receives a fresh
//	  packet whose changed field carries the NEW value (field by field);
//	② the reverse tooth: pull the "re-take the snapshot" out of the resync
//	  publisher (serve the cached packet instead) and this test MUST go red;
//	  restoring greens it. The mutation ran via `go test -c -overlay` with the
//	  mutated copy outside the repository; readings in
//	  docs/evidence/s1/35-g5-forged-allow-and-resync-20261116.md;
//	③ the SDK side (frontend/src/lib/panel.ts family) must not cache across
//	  show: the file:line evidence is booked in the same evidence file - the
//	  receiver replaces whole packets (panel.ts:560-578), the page keeps no
//	  copy of an answer (panel.ts:611-617, :13), and the architecture lint
//	  banning the storage faces is TestPanelFrontendIsStateless
//	  (internal/panel/frontend_hygiene_test.go:169). No frontend byte was
//	  touched by this leg.
//
// THE REAL BRIDGE CHAIN THIS TEST DRIVES (every hop a production function):
//
//	residentPanel.RequestShow      the product's gesture entry
//	  -> post -> drainTasks        the panel thread, modelled the way 35r9's
//	                               blockedPage resume() models it: the closures
//	                               run on the test goroutine IN PANEL-THREAD
//	                               ORDER, nothing is stubbed out
//	  -> showOnThread              the product trigger
//	  -> mgr.Show(ctx)             the product hot path (a stand-in page is
//	                               attached at the WebView2 COM boundary, the
//	                               same seam attachPanelWindow35r9 established
//	                               for the push-queue rulers)
//	  -> publishOnShow             the frozen resync-on-show trigger
//	                               (PLAN.md:2982, SPEC-08 §5.1)
//	  -> the boot-selected publisher  resident_windows.go's default branch
//	                               wiring line, verbatim: the resident pump
//	                               over the chain's REAL sources
//	  -> SnapshotPump.Publish      re-reads every source fresh per call
//	  -> evalPanelJS               the production push queue (panelPushPendingCap)
//	  -> deliverPanelPush -> Eval  the page boundary, where this test reads
//
// WHAT IS FRESH AND WHAT IS COMPARED. The tool state moved between the two
// shows is the ball style, written through the chain Manager's own setter -
// the same object the config.set door writes through - so the state move is a
// real write on the real config truth, never a test constant. The page's two
// packets are compared field by field: the moved field carries the new value
// in packet 2 and the old value in packet 1; generatedAt moves (fresh stamp
// per publish); an UNCHANGED dimension (modelSwitch.current) reads the same in
// both packets, which proves the comparison is per-field and not just "the
// bytes differ".
//
// NO GOROUTINE is spawned by this test (the panel thread is the synchronous
// drainTasks fiction; the publisher and the push queue run inline), so no
// observe.Registry.Spawn bookkeeping is owed, and no window is opened.

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
	"unsafe"

	webview2 "github.com/jchv/go-webview2"

	"github.com/CarlosShao/wisp/internal/panel"
)

// resyncPage35g6 is the page boundary stand-in: Dispatch queues (the panel
// thread's pump is not running in this rig), Eval records EVERY packet in
// order - unlike blockedPage35r9, which keeps only the last Eval, this rig
// needs both shows' packets side by side.
type resyncPage35g6 struct {
	mu      sync.Mutex
	pending []func()
	evals   []string
}

var _ webview2.WebView = (*resyncPage35g6)(nil)

func (p *resyncPage35g6) Dispatch(f func()) {
	if f == nil {
		return
	}
	p.mu.Lock()
	p.pending = append(p.pending, f)
	p.mu.Unlock()
}

func (p *resyncPage35g6) Eval(js string) {
	p.mu.Lock()
	p.evals = append(p.evals, js)
	p.mu.Unlock()
}

func (p *resyncPage35g6) Run()                                  {}
func (p *resyncPage35g6) Terminate()                            {}
func (p *resyncPage35g6) Destroy()                              {}
func (p *resyncPage35g6) Window() unsafe.Pointer                { return nil }
func (p *resyncPage35g6) SetTitle(title string)                 {}
func (p *resyncPage35g6) SetSize(w, h int, hint webview2.Hint)  {}
func (p *resyncPage35g6) Navigate(url string)                   {}
func (p *resyncPage35g6) SetHtml(html string)                   {}
func (p *resyncPage35g6) Init(js string)                        {}
func (p *resyncPage35g6) Bind(name string, f interface{}) error { return nil }

// resume drains the queued delivery closures in order - the pump not running
// is this rig's one modelled fact, the same fiction blockedPage35r9.resume()
// documents ("the panel thread for the purpose of this ruler").
func (p *resyncPage35g6) resume() int {
	p.mu.Lock()
	q := append([]func(){}, p.pending...)
	p.pending = nil
	p.mu.Unlock()
	for _, f := range q {
		f()
	}
	return len(q)
}

func (p *resyncPage35g6) packets() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.evals...)
}

// TestResyncOnShowDeliversFreshSnapshotFieldByField35g6 is tooth ①: hide ->
// the tool state moves -> show, and the packet the page receives on the second
// show carries the NEW value in the moved field, a fresh stamp, and the same
// values in the dimensions that did not move.
func TestResyncOnShowDeliversFreshSnapshotFieldByField35g6(t *testing.T) {
	dir := settingsDir248(t)
	host, cfgMgr, _, drafts, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("newResidentPanelManager refused to assemble: %v", err)
	}

	page := &resyncPage35g6{}
	// The 35r9 seam: the stand-in page attached at the WebView2 boundary
	// (m.w + created), so Show takes its real hot path without a window.
	attachPanelWindow35r9(t, host, page)

	rp := &residentPanel{mgr: host, tasks: make(chan func(), 16)}
	// The boot's own publisher selection, default branch verbatim
	// (resident_windows.go): the resident pump over the chain's real sources,
	// pushed into the host's real push queue. Wiring it any other way would
	// make this test measure a private publisher instead of the shipped one.
	rp.setSnapshotPublish(newResidentPanelPump(nil, cfgMgr, host.disp, drafts, host).publish)

	show := func(via string) {
		t.Helper()
		if !rp.RequestShow(via) {
			t.Fatalf("RequestShow(%s) was refused by the panel holder (startUp: %v)", via, rp.startUpErr())
		}
		rp.drainTasks()
	}

	// SHOW #1: the packet must arrive and carry the config's CURRENT style.
	show("35g6-show-1")
	if n := page.resume(); n == 0 {
		t.Fatalf("35g6 RED: show #1 queued no delivery closure - the resync-on-show trigger fired nothing")
	}
	pk1 := page.packets()
	if len(pk1) != 1 {
		t.Fatalf("35g6 RED: after one show the page holds %d packet(s), want exactly 1: %v", len(pk1), pk1)
	}
	top1 := decodedSnapshot65(t, pk1[0])
	raw1, ok := top1["ballStyle"]
	if !ok {
		t.Fatalf("35g6 RED: packet 1 carries no ballStyle key; top-level keys: %v", keysOf65(top1))
	}
	var bs1 panel.BallStyleSection
	if err := json.Unmarshal(raw1, &bs1); err != nil {
		t.Fatalf("35g6 RED: packet 1's ballStyle does not decode as the carrier: %v", err)
	}
	before := cfgMgr.BallStyleOptions().Style
	if bs1.Style != before {
		t.Fatalf("35g6 RED: packet 1's ballStyle.style = %q, want the config's current %q - "+
			"the packet is not the state as of show #1", bs1.Style, before)
	}
	var ms1 panel.ModelSwitchSection
	if raw, ok := top1["modelSwitch"]; ok {
		if err := json.Unmarshal(raw, &ms1); err != nil {
			t.Fatalf("35g6 RED: packet 1's modelSwitch does not decode as the carrier: %v", err)
		}
	}
	gen1 := string(top1["generatedAt"])

	// HIDE through the product gesture path (toggle hides a shown panel).
	if !rp.RequestToggle("35g6-toggle") {
		t.Fatalf("RequestToggle was refused by the panel holder")
	}
	rp.drainTasks()
	if host.IsShown() {
		t.Fatalf("35g6 RED: the hide did not take IsShown down - the second show would not be a resync")
	}

	// THE TOOL STATE MOVES, through the real settings-door writer: the chain
	// Manager's own setter, validated against the schema's own vocabulary.
	voc := cfgMgr.BallStyleOptions()
	next := ""
	for _, s := range voc.Vocabulary {
		if s != voc.Style {
			next = s
			break
		}
	}
	if next == "" {
		t.Fatalf("35g6: the schema's vocabulary offers no second style to move to; the state move cannot be driven")
	}
	if _, err := cfgMgr.SetBallStyle(next); err != nil {
		t.Fatalf("SetBallStyle(%q) refused: %v", next, err)
	}
	if got := cfgMgr.BallStyleOptions().Style; got != next {
		t.Fatalf("35g6: the setter did not move the config truth: style = %q, want %q", got, next)
	}

	// generatedAt is RFC3339 at second granularity: one second of separation
	// makes the stamps comparable field-by-field as well.
	time.Sleep(1100 * time.Millisecond)

	// SHOW #2: the packet must arrive with the NEW state, field by field.
	show("35g6-show-2")
	if n := page.resume(); n == 0 {
		t.Fatalf("35g6 RED: show #2 queued no delivery closure - the resync-on-show trigger fired nothing on the second show")
	}
	pk2 := page.packets()
	if len(pk2) != 2 {
		t.Fatalf("35g6 RED: after hide + state move + show the page holds %d packet(s), want exactly 2 "+
			"(one whole packet per show, none cached, none merged): %v", len(pk2), pk2)
	}
	top2 := decodedSnapshot65(t, pk2[1])
	raw2, ok := top2["ballStyle"]
	if !ok {
		t.Fatalf("35g6 RED: packet 2 carries no ballStyle key; top-level keys: %v", keysOf65(top2))
	}
	var bs2 panel.BallStyleSection
	if err := json.Unmarshal(raw2, &bs2); err != nil {
		t.Fatalf("35g6 RED: packet 2's ballStyle does not decode as the carrier: %v", err)
	}

	// The moved field, field by field: packet 2 carries the NEW value, packet 1
	// kept the old one, and the two differ. A cached packet fails all three.
	if bs2.Style != next {
		t.Errorf("35g6 RED: packet 2's ballStyle.style = %q, want the NEW value %q - "+
			"the page received a stale state across hide/mutate/show (cross-show cache or a stale packet)",
			bs2.Style, next)
	}
	if bs1.Style != before {
		t.Errorf("35g6 RED: packet 1's ballStyle.style = %q, want the pre-move value %q - "+
			"the first packet moved after the fact, which is not a snapshot", bs1.Style, before)
	}
	if bs2.Style == bs1.Style {
		t.Errorf("35g6 RED: both packets carry ballStyle.style = %q - the second show did not deliver the state move", bs1.Style)
	}
	if len(bs2.Vocabulary) == 0 || len(bs1.Vocabulary) == 0 {
		t.Errorf("35g6 RED: the vocabulary must travel verbatim in every packet (packet 1: %d row(s), packet 2: %d)",
			len(bs1.Vocabulary), len(bs2.Vocabulary))
	}

	// Fresh stamp: a resync re-takes the snapshot, so the packet is stamped at
	// publish time, not at first-show time.
	if gen2 := string(top2["generatedAt"]); gen2 == gen1 {
		t.Errorf("35g6 RED: generatedAt is identical across the two shows (%q) - "+
			"the second packet was not built at show time", gen1)
	}

	// The UNCHANGED dimension must read the SAME in both packets: the
	// comparison is per field, not "the bytes differ somehow".
	if raw, ok := top2["modelSwitch"]; ok {
		var ms2 panel.ModelSwitchSection
		if err := json.Unmarshal(raw, &ms2); err != nil {
			t.Fatalf("35g6 RED: packet 2's modelSwitch does not decode as the carrier: %v", err)
		}
		if ms1.Current != ms2.Current {
			t.Errorf("35g6 RED: modelSwitch.current moved without a state move: %v -> %v",
				ms1.Current, ms2.Current)
		}
	}
	t.Logf("35g6 resync field-by-field: ballStyle.style %q -> %q, generatedAt %s -> %s",
		bs1.Style, bs2.Style, gen1, string(top2["generatedAt"]))
}
