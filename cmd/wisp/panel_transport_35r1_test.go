//go:build windows

package main

// Ticket 35 AC#6 (`:52`) - the page<->host transport, judged on the Go side.
//
// WHAT THIS CASE IS AND WHAT IT DELIBERATELY IS NOT. The page posts its composer
// envelope on window.chrome.webview.postMessage (frontend/src/lib/panel.ts:179/:218,
// READ-ONLY for this fleet). Ticket 35:52-59 records that, with only the library's
// RPC pipe behind that call, the raw envelope reaches go-webview2's msgcb, parses as
// an RPC frame with id=0 and an unbound method, and dies with no Go-side log. The
// chosen fix (orchestrator shape 甲 sub-shape ①, ticket 35:150) is a Go-side
// forwarding hook, installed at installPanelTransport, that reroutes the page's
// postMessage into the already-bound window.wispDispatch(raw) so the envelope
// arrives at dispatchRaw verbatim.
//
// So this case does NOT call dispatchRaw (ticket 35 AC#7 (b) forbids that for this
// edge) and does NOT run w.Eval("window.wispDispatch(...)") (that is N6, the test
// playing the page, and it stays green with no forwarding at all - zero discriminating
// power). Instead it hands the page's postMessage to the transport the way a browser
// would, and the ONLY thing that carries the raw string to the door is the forwarding
// installPanelTransport registered via Init. fakePageControl below models the library
// honestly: msgcb parses {id, method, params} and routes ONLY to a bound method, so a
// raw page envelope with no forwarding lands exactly where ticket 35:59 says it does
// (id=0, method "panel.mode.request", nothing bound -> the door is never invoked).
//
// ★FALSIFIABILITY (the tooth): delete or no-op installPanelTransport's w.Init(...)
// forwarding statement and this case goes red, because the door is then reached only
// by the mis-parsed raw envelope, which never fires it. N1 (Go roster count) and
// N2/N3 (page-wording counters) cannot stand in for it (ticket 35:67).
//
// SCOPE: AC#6 only. The (d)-1/-2/-3 guard-wording assertions (bridge.go:133/:136-137/
// :140-141 each named) are AC#7 (`:63`) and belong to the next leg; this case does not
// assert them, and its positive-control envelope is a roster method that PASSES the
// guard order (bridge.go:132), so it exercises none of the three reject branches.

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// pageModeRequestEnvelope is the verbatim string the page's
// sendRequest("panel.mode.request", { to }) posts (frontend/src/lib/panel.ts:246 and
// the :211 sendRequest that stamps requestId "pc-<seq>-<uuid>" + source
// "panel-composer"). It is a roster method (bridge.go:132 answers it), so it is the
// shape ticket 35:70 designates as AC#6/AC#7's positive control - unlike the
// panel.approval.request at :181, which the roster guard refuses first.
const (
	pageModeRequestEnvelope = `{"method":"panel.mode.request","requestId":"pc-1-3f2b1c0d-9e7a-4c1b-8f14-e45fceea469a","source":"panel-composer","to":"ask_every_step"}`
	pageModeRequestID       = "pc-1-3f2b1c0d-9e7a-4c1b-8f14-e45fceea469a"
)

// modeSpy35r1 is the downstream of dispatchRaw for a mode request: it records the
// parsed ComposerRequest the router hands the mode handler, so the case can assert the
// page's requestId actually arrived at the door rather than trusting a JS-side claim.
type modeSpy35r1 struct {
	mu   sync.Mutex
	reqs []panel.ComposerRequest
}

func (s *modeSpy35r1) HandleModeRequest(_ context.Context, req panel.ComposerRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reqs = append(s.reqs, req)
	return nil
}

func (s *modeSpy35r1) snapshot() []panel.ComposerRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]panel.ComposerRequest(nil), s.reqs...)
}

// fakePageControl implements pageTransport (Bind + Init) and models the browser + the
// go-webview2 message channel that sit between the page's postMessage and dispatchRaw.
// It is deliberately minimal: Bind records the door callback, Init records the
// forwarding script, and msgcb re-derives the library's routing so nothing here calls
// dispatchRaw itself.
type fakePageControl struct {
	door        func(raw string) string // the wispDispatch binding callback production installed
	hooks       []string                // every Init string the host registered
	doorRounds  int
	lastReply   string
	rpcDeadSlot int // raw / unbound frames that found no bound method (the ticket 35:59 death)
	seq         int
}

// Bind mirrors webview.Bind: it stores the callback under the binding name and (like
// the real library) registers a per-binding Init stub. Only func(string) string doors
// are exercised here.
func (fc *fakePageControl) Bind(name string, f interface{}) error {
	cb, ok := f.(func(string) string)
	if !ok {
		// pageTransport in the library only ever binds this shape for the panel door;
		// anything else means installPanelTransport changed contract, which this case
		// must refuse loudly rather than mis-model.
		return &bindShapeError{name: name}
	}
	if name == panelDispatchBinding {
		fc.door = cb
	}
	// The library's own stub for the binding: note it does NOT literally name
	// window.wispDispatch or chrome.webview.postMessage, so it never trips the
	// forwarding detector below.
	fc.hooks = append(fc.hooks, "binding-stub:"+name)
	return nil
}

// Init mirrors webview.Init: the page-side forwarding hook is one of these strings.
func (fc *fakePageControl) Init(js string) { fc.hooks = append(fc.hooks, js) }

type bindShapeError struct{ name string }

func (e *bindShapeError) Error() string {
	return "fake Bind saw a non (string)->string door: " + e.name
}

// forwardingInstalled reports whether the host installed a hook that reroutes the
// page's window.chrome.webview.postMessage into window.wispDispatch. It reads the
// PRODUCTION Init scripts only, so the verdict flips when installPanelTransport's w.Init
// forwarding statement is deleted.
func (fc *fakePageControl) forwardingInstalled() bool {
	for _, js := range fc.hooks {
		if strings.Contains(js, "chrome.webview") &&
			strings.Contains(js, "postMessage") &&
			strings.Contains(js, "window.wispDispatch") {
			return true
		}
	}
	return false
}

// deliverPostMessage is the browser's hop. A page calls chrome.webview.postMessage(raw).
// WITH the forwarding hook installed, that raw string is what window.wispDispatch is fed,
// and the library stub wraps it as an RPC frame {id, method:"wispDispatch", params:[raw]}
// before msgcb sees it. WITHOUT the hook, the native pipe hands msgcb the page's raw
// envelope verbatim - the ticket 35:59 death. This method is the ONLY entry point the
// test uses to "post" a request; it never calls dispatchRaw directly.
func (fc *fakePageControl) deliverPostMessage(message string) {
	if fc.forwardingInstalled() {
		fc.seq++
		frame, err := json.Marshal(map[string]any{
			"id":     fc.seq,
			"method": panelDispatchBinding,
			"params": []string{message},
		})
		if err != nil {
			panic("fake: cannot build the forwarded RPC frame: " + err.Error())
		}
		fc.msgcb(string(frame))
		return
	}
	// No forwarding: the browser delivers the page's own envelope straight to msgcb.
	fc.msgcb(message)
}

// msgcb is go-webview2's webview.msgcb (webview.go:139-160) reproduced honestly: parse
// {id, method, params}; route ONLY to a bound method; otherwise reply to a window._rpc
// slot the page does not have. The door callback is never invoked for an unbound method,
// which is exactly why an envelope that never got forwarded cannot reach dispatchRaw.
func (fc *fakePageControl) msgcb(msg string) {
	var d struct {
		ID     int               `json:"id"`
		Method string            `json:"method"`
		Params []json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal([]byte(msg), &d); err != nil {
		fc.rpcDeadSlot++
		return
	}
	if d.Method != panelDispatchBinding || fc.door == nil {
		// Unknown/unbound method (the raw page envelope): the library "resolves" a
		// window._rpc[id] slot the page has no listener for. The request dies here.
		fc.rpcDeadSlot++
		return
	}
	if len(d.Params) == 0 {
		fc.rpcDeadSlot++
		return
	}
	var raw string
	if err := json.Unmarshal(d.Params[0], &raw); err != nil {
		fc.rpcDeadSlot++
		return
	}
	fc.doorRounds++
	fc.lastReply = fc.door(raw) // -> dispatchRaw -> ComposerDispatch.Handle -> the mode handler
}

func TestPagePostMessageEnvelopeReachesDispatchRawViaTransport(t *testing.T) {
	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")

	ctx := context.Background()
	fc := &fakePageControl{}
	if err := mgr.installPanelTransport(fc, ctx); err != nil {
		t.Fatalf("installPanelTransport: %v (the shipping transport wiring itself failed to install)", err)
	}

	// Positive-control gate #1: the forwarding statement really registered. If this is
	// false the transport is not connected and everything below is the ticket 35:59 shape.
	if !fc.forwardingInstalled() {
		t.Fatal("AC#6 RED: installPanelTransport registered no page->wispDispatch forwarding hook " +
			"(no Init string reroutes window.chrome.webview.postMessage into window.wispDispatch). " +
			"Without it the page's envelope is mis-parsed by msgcb and never reaches dispatchRaw (ticket 35:52-59).")
	}

	// The page posts its real composer envelope through the transport the page uses.
	fc.deliverPostMessage(pageModeRequestEnvelope)

	if fc.doorRounds != 1 {
		t.Fatalf("AC#6 RED: the page's postMessage did not drive the door exactly once (doorRounds=%d, rpcDeadSlot=%d). "+
			"The forwarding statement is what carries the raw string to dispatchRaw; if it is missing the envelope "+
			"stays an unbound RPC frame and the door is never fired.", fc.doorRounds, fc.rpcDeadSlot)
	}

	// (c) dispatchRaw's downstream was reached carrying that correlationId/requestId.
	reqs := spy.snapshot()
	if len(reqs) != 1 {
		t.Fatalf("AC#6 RED: the mode handler downstream of dispatchRaw ran %d times, want 1 (the request reached the door but not the router: %+v)", len(reqs), reqs)
	}
	got := reqs[0]
	if got.RequestID != pageModeRequestID {
		t.Errorf("AC#6 RED: dispatchRaw's downstream did not carry the page's requestId: got %q want %q", got.RequestID, pageModeRequestID)
	}
	if got.Method != panel.MethodModeRequest {
		t.Errorf("AC#6 RED: the routed method is %q, want %q (the raw envelope was altered in transit)", got.Method, panel.MethodModeRequest)
	}
	if got.Source != panel.ComposerRequestSource {
		t.Errorf("AC#6 RED: the source carried to the router is %q, want %q", got.Source, panel.ComposerRequestSource)
	}
	// The door returns dispatchRaw's receipt; for an accepted mode request that is the
	// empty receipt with no refusal, i.e. the request was not silently dropped.
	if fc.rpcDeadSlot != 0 {
		t.Errorf("AC#6 RED: %d page frame(s) died at msgcb's unbound-method branch even though the door fired", fc.rpcDeadSlot)
	}
}
