//go:build windows

package main

// Ticket 35 :44 (Acceptance: "Concurrency routing: 10 parallel invokes with
// shuffled replies -> correct pairing") - G1's ruler for the RPC reply-pairing
// half, plus the named finding the dispatch asked for: WHAT IS THE PAIRING KEY.
//
// ============================================================================
// THE FINDING, WRITTEN DOWN BEFORE THE RULER (dispatch step 2's precondition)
// ============================================================================
//
// The ticket wording says "correlationId routing" / "requestId". The code does
// NOT pair on requestId. It pairs on the RPC FRAME's own integer `id`, which is
// the page-side `window._rpc` sequence number the bound stub mints:
//
//   - cmd/wisp/panel_control_windows.go:793   id := strconv.Itoa(d.ID)
//   - cmd/wisp/panel_control_windows.go:823   c.Eval("window._rpc[" + id + "].resolve(" + string(b) + "); window._rpc[" + id + "] = undefined")
//   - cmd/wisp/panel_control_windows.go:907   window.external.invoke(JSON.stringify({ id: seq, method: name, params: ... }))
//
// `requestId` is NOT on the reply path at all: it is a log/audit correlation key.
// bridge.go:179-182 only requires it be non-empty; NewRequestID (bridge.go:120)
// has zero non-test callers in this tree; the id the host resolves is the frame
// id, full stop. So a ruler that "pairs by requestId" would be testing a field
// the product never consults - which the dispatch explicitly forbids faking.
//
// ============================================================================
// WHY A HEADLESS RULER AND NOT A REAL WINDOW (the honest limit)
// ============================================================================
//
// The concrete `*panelControl` gives the reply path no seam a headless test can
// observe: `msgcb`'s reply is literally a string handed to `c.Eval`, and
// `panelControl.Eval` is a COM call through `*edge.Chromium` (unexported fields,
// no injectable core). Adding a seam would be a PRODUCTION CHANGE, which this
// dispatch forbids ("这两格票面要的是尺"). The alternative - a second real panel
// window - is measured in this very package to flake: panel_unbound_253r6_windows_test.go
// header records "a second panel window created right after a previous one
// intermittently dies mid-create (measured ... on the pristine tree)". So this
// file drives:
//   - the page side: an async `_rpc[seq]` SLOT TABLE (the thing the recon said the
//     current fixture lacks - a synchronous stub makes any pairing look correct);
//   - the host side: `g1Pipe.msgcb`, a line-for-line port of
//     panel_control_windows.go:782-826's bound-door success branch, including
//     `id := strconv.Itoa(d.ID)` and `_rpc[id]`;
//   - the ANSWERS: the REAL production door `panel.ComposerDispatch.Handle`, so
//     the payloads are product output, not invented strings.
//
// Teeth are therefore two-layered, and both are stated plainly rather than
// implied:
//   (a) BEHAVIOURAL, on the pipe model + slot table: the same 10 frames are
//       replayed through two deliberately-WRONG pairing rules (`byArrival`,
//       `singleSlot`) and the test asserts those MUST mismatch. If the harness
//       could not tell the right pairing from a wrong one this control fails -
//       that is exactly the "sync stub is always green" trap the recon named.
//   (b) ANCHOR, on production: `assertProductionPairsByFrameID35g1` reads
//       panel_control_windows.go and fails if the pairing expression ever stops
//       being `_rpc[id]` keyed off `d.ID`. This is a source guard (A1880 calls
//       that a literal nail, not behaviour) and is labelled as such - it is the
//       tie between the model above and the line it mirrors, not a substitute
//       for behaviour.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// ----------------------------------------------------------------------------
// the page-side slot table (async, out-of-order settle)
// ----------------------------------------------------------------------------

type g1SlotState struct {
	settled bool
	value   string
	settles int // how many times resolve() was called; a Promise settles once
}

type g1Page struct {
	mu      sync.Mutex
	slots   map[int]*g1SlotState
	nextSeq int
}

// newG1Page starts nextSeq at 1, the same seed the production Bind stub writes
// (panel_control_windows.go:898 `window._rpc = window._rpc || {nextSeq: 1}` and
// :900 `var seq = RPC.nextSeq++`).
func newG1Page() *g1Page {
	return &g1Page{slots: map[int]*g1SlotState{}, nextSeq: 1}
}

// openCall is one page-side invoke: mint the next seq and register a slot keyed by
// it. This is the "10 parallel invokes are all outstanding" state: every slot is
// open and unresolved before any reply is applied.
func (p *g1Page) openCall() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	id := p.nextSeq
	p.nextSeq++
	p.slots[id] = &g1SlotState{}
	return id
}

// settleResolve models `window._rpc[id].resolve(value)`: a browser Promise settles
// at most once, so later calls are counted but do not overwrite.
func (p *g1Page) settleResolve(id int, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	s, ok := p.slots[id]
	if !ok {
		return
	}
	s.settles++
	if !s.settled {
		s.settled = true
		s.value = value
	}
}

func (p *g1Page) reading(id int) (settled bool, value string, settles int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	s := p.slots[id]
	if s == nil {
		return false, "", 0
	}
	return s.settled, s.value, s.settles
}

// ----------------------------------------------------------------------------
// the host-side RPC pipe: a port of panel_control_windows.go:782-826
// ----------------------------------------------------------------------------

// g1Reply is one queued reply: the slot id the host will resolve and the answer.
// Production queues these as `Dispatch(func(){ Eval(...) })` closures; modelling
// them as data is what lets the test apply them OUT OF ORDER (a real browser can
// run its queue in any order; the pairing key must survive that).
type g1Reply struct {
	id    int
	value string
}

type g1Pipe struct {
	page     *g1Page
	bindings map[string]func(string) string
	replies  []g1Reply
}

// g1Frame mirrors ctlRPCMessage (panel_control_windows.go:764-768).
type g1Frame struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

// msgcb is the bound-door success branch of panel_control_windows.go:782-826,
// faithfully: parse the frame, refuse a frame naming no bound door (the same
// guard), call the binding with params[0] as a string, and QUEUE a reply that
// resolves `_rpc[<frame id>]`. The queue (not an immediate settle) is the async
// slot the recon demanded.
func (h *g1Pipe) msgcb(msg string) {
	var d g1Frame
	if err := json.Unmarshal([]byte(msg), &d); err != nil {
		return
	}
	door, ok := h.bindings[d.Method]
	if !ok {
		return
	}
	if len(d.Params) != 1 {
		return
	}
	var raw string
	if err := json.Unmarshal(d.Params[0], &raw); err != nil {
		return
	}
	answer := door(raw)
	// panel_control_windows.go:793+823: the reply is addressed by the FRAME ID.
	h.replies = append(h.replies, g1Reply{id: d.ID, value: answer})
}

// applyByFrameID is the CORRECT pairing: each queued reply resolves the slot named
// by its OWN frame id, regardless of the order the replies are applied in.
func (p *g1Page) applyByFrameID(replies []g1Reply) {
	for _, r := range replies {
		p.settleResolve(r.id, r.value)
	}
}

// applyByArrival is a WRONG pairing used as a control: the kth reply is forced
// into the kth slot, ignoring the frame id. If the real pairing is correct, this
// must NOT also be correct when arrival order differs from id order.
func (p *g1Page) applyByArrival(replies []g1Reply) {
	for pos, r := range replies {
		p.settleResolve(pos+1, r.value)
	}
}

// applySingleSlot is the second WRONG pairing: every reply is forced into slot 1.
func (p *g1Page) applySingleSlot(replies []g1Reply) {
	for _, r := range replies {
		p.settleResolve(1, r.value)
	}
}

// ----------------------------------------------------------------------------
// fixtures
// ----------------------------------------------------------------------------

const g1N = 10

// g1RawEnvelope builds a composer envelope whose REFUSAL names a token unique to
// this call. `panel.bogus.<id>` is off the roster, so the real production door
// (ComposerDispatch.Handle) refuses it and the refusal sentence carries both the
// method and the requestId - which is what makes the answers distinguishable, so
// "which answer landed where" is decidable rather than assumed.
func g1RawEnvelope(id int) string {
	m := map[string]any{
		"method":    fmt.Sprintf("panel.bogus.%d", id),
		"requestId": fmt.Sprintf("rid-k%d", id),
		"source":    panel.ComposerRequestSource,
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func g1FrameJSON(id int) string {
	params, _ := json.Marshal([]string{g1RawEnvelope(id)})
	return fmt.Sprintf(`{"id":%d,"method":%q,"params":%s}`, id, panelDispatchBinding, params)
}

// g1ShuffledIDs is a fixed permutation in which the replies are "delivered": it is
// deliberately NOT 1..10, so a pairing that quietly uses arrival position goes red.
func g1ShuffledIDs() []int {
	return []int{7, 2, 9, 1, 10, 4, 6, 3, 8, 5}
}

// ----------------------------------------------------------------------------
// the ruler
// ----------------------------------------------------------------------------

func TestRPCReplyPairsByFrameIDUnder10WayOutOfOrder35g1(t *testing.T) {
	assertProductionPairsByFrameID35g1(t)

	// The real door: production Handle. Its refusals for panel.bogus.<id> carry
	// rid-k<id>, so every answer is self-identifying.
	disp := &panel.ComposerDispatch{}
	door := func(raw string) string {
		out, _ := disp.Handle(context.Background(), raw)
		return out
	}

	// --- the correct run ------------------------------------------------
	page := newG1Page()
	ids := make([]int, g1N)
	for i := 0; i < g1N; i++ {
		ids[i] = page.openCall() // all 10 slots open, before any reply
	}
	pipe := &g1Pipe{page: page, bindings: map[string]func(string) string{panelDispatchBinding: door}}
	// Deliver the 10 frames, but collect the replies in arrival order and then
	// apply them in a SCRAMBLED order: that is "replies shuffled".
	for _, id := range ids {
		pipe.msgcb(g1FrameJSON(id))
	}
	if len(pipe.replies) != g1N {
		t.Fatalf("the pipe queued %d replies, want %d: every bound frame must produce one reply", len(pipe.replies), g1N)
	}
	byID := map[int]g1Reply{}
	for _, r := range pipe.replies {
		byID[r.id] = r
	}
	shuffled := make([]g1Reply, 0, g1N)
	for _, id := range g1ShuffledIDs() {
		r, ok := byID[id]
		if !ok {
			t.Fatalf("no reply for frame id %d - the pipe dropped a frame", id)
		}
		shuffled = append(shuffled, r)
	}
	page.applyByFrameID(shuffled)

	for _, id := range ids {
		settled, value, settles := page.reading(id)
		if !settled {
			t.Fatalf("slot %d never settled: with pairing by frame id every reply lands in its own slot even out of order", id)
		}
		if settles != 1 {
			t.Errorf("slot %d was resolved %d times, want exactly 1 (a Promise settles once; a second resolve against the same slot is a mis-route)", id, settles)
		}
		rids := g1RidsIn(value)
		if len(rids) == 0 {
			t.Errorf("slot %d received %q, which names no rid-k<id> at all - the answer is not self-identifying", id, value)
		}
		for _, n := range rids {
			if n != id {
				t.Errorf("slot %d received %q, which names rid-k%d - that reply belongs to another caller (pairing crossed streams)", id, value, n)
			}
		}
	}

	// --- control A: pairing by ARRIVAL position must be caught -----------
	if !g1ControlMismatches(t, "byArrival", (*g1Page).applyByArrival) {
		t.Fatal("control byArrival did NOT mismatch: the harness cannot tell frame-id pairing from arrival-order pairing, so this ruler is tautological and must not be trusted (the exact trap recon 2-c names)")
	}

	// --- control B: single-slot pairing must be caught -------------------
	if !g1ControlMismatches(t, "singleSlot", (*g1Page).applySingleSlot) {
		t.Fatal("control singleSlot did NOT mismatch: forcing every reply into one slot must leave the other nine unsettled/wrong, and the harness must see it")
	}

	t.Logf("10-way out-of-order pairing: all %d slots settled by frame id, each naming its own rid-k<id>; control byArrival and singleSlot both mismatched as required", g1N)
}

// g1ridRe extracts the rid-k<id> tokens an answer names. Matching the NUMBER
// rather than a substring is load-bearing: "rid-k1" is a substring of "rid-k10",
// so a Contains check would call a rid-k10 answer "belonging to slot 1".
var g1ridRe = regexp.MustCompile(`rid-k(\d+)`)

func g1RidsIn(value string) []int {
	var out []int
	for _, m := range g1ridRe.FindAllStringSubmatch(value, -1) {
		n, err := strconv.Atoi(m[1])
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

// g1ControlMismatches replays the same 10 frames through one WRONG pairing rule
// and reports whether at least one slot ended up with the wrong caller's answer.
// It returns true when the control is caught (which is what the test wants).
func g1ControlMismatches(t *testing.T, name string, apply func(*g1Page, []g1Reply)) bool {
	t.Helper()
	disp := &panel.ComposerDispatch{}
	door := func(raw string) string { out, _ := disp.Handle(context.Background(), raw); return out }

	page := newG1Page()
	ids := make([]int, g1N)
	for i := 0; i < g1N; i++ {
		ids[i] = page.openCall()
	}
	pipe := &g1Pipe{page: page, bindings: map[string]func(string) string{panelDispatchBinding: door}}
	for _, id := range ids {
		pipe.msgcb(g1FrameJSON(id))
	}
	byID := map[int]g1Reply{}
	for _, r := range pipe.replies {
		byID[r.id] = r
	}
	shuffled := make([]g1Reply, 0, g1N)
	for _, id := range g1ShuffledIDs() {
		shuffled = append(shuffled, byID[id])
	}
	apply(page, shuffled)

	mismatch := 0
	for _, id := range ids {
		settled, value, _ := page.reading(id)
		rids := g1RidsIn(value)
		ok := settled && len(rids) > 0
		if ok {
			for _, n := range rids {
				if n != id {
					ok = false
				}
			}
		}
		if !ok {
			mismatch++
		}
	}
	if mismatch == 0 {
		t.Logf("control %q: 0/%d slots mismatched (WRONG pairing looked correct)", name, g1N)
	} else {
		t.Logf("control %q: %d/%d slots mismatched (caught, as required)", name, mismatch, g1N)
	}
	return mismatch > 0
}

// ----------------------------------------------------------------------------
// production anchor: the model above mirrors these exact lines
// ----------------------------------------------------------------------------

// assertProductionPairsByFrameID35g1 is the tie between the model in this file and
// the shipped code: panel_control_windows.go must still (a) derive the reply
// address from the frame id and (b) resolve `_rpc[id]`. A production change from
// id-pairing to anything else reddens this, which is the signal that the model
// above has drifted from the line it claims to mirror. It is a source guard, and
// is labelled one on purpose (A1880): the behavioural teeth live in the model +
// the two controls, not here.
func assertProductionPairsByFrameID35g1(t *testing.T) {
	t.Helper()
	path := ""
	var src string
	for _, cand := range []string{
		filepath.Join("panel_control_windows.go"),
		filepath.Join("cmd", "wisp", "panel_control_windows.go"),
		filepath.Join("..", "cmd", "wisp", "panel_control_windows.go"),
	} {
		if b, err := os.ReadFile(cand); err == nil && len(b) > 0 {
			path = cand
			src = string(b)
			break
		}
	}
	if path == "" {
		t.Fatalf("could not locate panel_control_windows.go from cwd %q - the production anchor cannot run, which is a failed measurement", mustGetwd35g1())
	}
	_ = path
	if !strings.Contains(src, "id := strconv.Itoa(d.ID)") {
		t.Errorf("panel_control_windows.go no longer derives the reply address as `id := strconv.Itoa(d.ID)`: the pairing key changed and the model in this file no longer mirrors production")
	}
	if !strings.Contains(src, `window._rpc[" + id + "].resolve(`) {
		t.Errorf(`panel_control_windows.go no longer resolves window._rpc["+id+"].resolve(...): the model's pairing (frame id -> _rpc[id]) has drifted from production`)
	}
}

func mustGetwd35g1() string {
	wd, err := os.Getwd()
	if err != nil {
		return "(getwd failed)"
	}
	return wd
}
