package panel

// Ticket 33 slice A (H4 + H5) - the three "take it away and it must go red"
// judgement nails 33-a1 §7.4 asked for, plus the two honesty nails that keep
// this slice from being read as a working panel.
//
// Counting rule, stated once because it is the whole point of the file: every
// judgement below is taken from a CALL COUNT or an error IDENTITY, never from a
// log string. A router that writes a pretty sentence while never calling the
// handler passes a string assertion and fails the property (preflight §⑤, and
// 33-a1 restated it as "看被调次数，不看日志字符串").
//
// What is real here and what is not: the chain raw -> ParseComposerRequest ->
// dispatch -> *ModeWriteHandler.HandleModeRequest -> ModeWriter.Set is production
// code end to end. The only fake sits at ModeWriter, which is the seam
// composer_handlers.go already declares for perm.Store (AGENTS.md §1.3). That is
// NOT evidence that the panel can click anything: TestSliceAAttachesNoHostAnd
// NamesTheOpenWindowHops pins the opposite reading in the same file.

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
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// hop33Raw builds one postMessage envelope the way a page would send it. Method
// names come from bridge.go's constants, not from literals, so a rename on the
// whitelist moves these tests instead of quietly leaving them green.
func hop33Raw(t *testing.T, fields map[string]any) string {
	t.Helper()
	data, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	return string(data)
}

// hop33ModeWriter is a stand-in for *perm.Store at the seam the handler already
// declares. setCalls is the number 33-a1's first judgement asks about.
type hop33ModeWriter struct {
	current    risk.Mode
	readable   bool
	setCalls   int
	lastTo     risk.Mode
	lastOrigin string
	lastActor  string
}

func (w *hop33ModeWriter) PermissionMode() risk.Mode {
	if !w.readable {
		// Off-ladder: the handler must fail closed rather than guess a档.
		return risk.Mode(909)
	}
	return w.current
}

func (w *hop33ModeWriter) Set(_ context.Context, to risk.Mode, origin, actor string) error {
	w.setCalls++
	w.lastTo = to
	w.lastOrigin = origin
	w.lastActor = actor
	return nil
}

// hop33Capture records the ComposerRequest exactly as the router handed it over,
// which is the only way to see whether the router edited the attribution.
type hop33Capture struct {
	calls int
	got   ComposerRequest
	err   error
}

func (c *hop33Capture) HandleModeRequest(_ context.Context, req ComposerRequest) error {
	c.calls++
	c.got = req
	return c.err
}

// hop33RefusingHandler stands in for a handler that exists and says no.
type hop33RefusingHandler struct{ calls int }

func (h *hop33RefusingHandler) HandleWorkspaceRequest(context.Context, ComposerRequest) error {
	h.calls++
	return errors.New("hop33: 处理器自己拒绝")
}

func (h *hop33RefusingHandler) HandleAttachmentRequest(context.Context, ComposerRequest) error {
	h.calls++
	return errors.New("hop33: 处理器自己拒绝")
}

func (h *hop33RefusingHandler) HandleMessageRequest(context.Context, ComposerRequest) error {
	h.calls++
	return errors.New("hop33: 处理器自己拒绝")
}

// hop33Audit collects the audit lines so "refused AND audited" is one assertion
// and "silently dropped" has a number attached to it.
type hop33Audit struct{ lines []string }

func (a *hop33Audit) fn() AuditFunc {
	return func(format string, args ...any) {
		a.lines = append(a.lines, strings.TrimSpace(fmt.Sprintf(format, args...)))
	}
}

// ------------------------------------------------------------- AC#A: the hop

// TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce is judgement 1: the
// mode door really opens. Vehicle: in composer_dispatch.go, replace the body of
// the mode case's call - `return d.Mode.HandleModeRequest(ctx, req)` - with
// `return nil`, i.e. routing removed while the rest of the file stands, and this
// test goes red on setCalls (not on any prose).
func TestRawEnvelopeReachesTheAttachedModeHandlerExactlyOnce(t *testing.T) {
	writer := &hop33ModeWriter{current: risk.ModeAskHighRisk, readable: true}
	handler := &ModeWriteHandler{Modes: writer, Audit: (&hop33Audit{}).fn(), Actor: "hop33"}
	d := &ComposerDispatch{Mode: handler}

	raw := hop33Raw(t, map[string]any{
		"method":    MethodModeRequest,
		"requestId": "rid-ac-a",
		"source":    ComposerRequestSource,
		// Narrowing, so R20/M4 needs no confirmation leg and the write leg runs.
		"to": risk.ModeAskEveryStepName,
	})
	reply, err := d.Handle(context.Background(), raw)
	if err != nil {
		t.Fatalf("a well-formed narrowing request was refused: %v (reply=%q)", err, reply)
	}
	if writer.setCalls != 1 {
		t.Fatalf("write leg called %d time(s), want exactly 1: routing was removed or doubled, "+
			"and the档 did not move once (composer_handlers.go:141)", writer.setCalls)
	}
	if writer.lastTo != risk.ModeAskEveryStep {
		t.Errorf("write leg received %v, want %v", writer.lastTo, risk.ModeAskEveryStep)
	}
	if writer.lastOrigin != PanelModeOrigin {
		t.Errorf("write leg received origin %q, want %q: the router invented a second door",
			writer.lastOrigin, PanelModeOrigin)
	}
}

// TestAcceptedRequestInventsNoReplyLine keeps this slice from encroaching on
// H10: the reply the page sees is the host's job, so the accepted path is empty.
func TestAcceptedRequestInventsNoReplyLine(t *testing.T) {
	writer := &hop33ModeWriter{current: risk.ModeAskHighRisk, readable: true}
	d := &ComposerDispatch{Mode: &ModeWriteHandler{Modes: writer}}
	raw := hop33Raw(t, map[string]any{
		"method": MethodModeRequest, "requestId": "rid-reply", "source": ComposerRequestSource,
		"to": risk.ModeAskEveryStepName,
	})
	reply, err := d.Handle(context.Background(), raw)
	if err != nil || reply != "" {
		t.Fatalf("accepted path returned (%q, %v), want (\"\", nil): a router that answers for the host is H10, not H4",
			reply, err)
	}
}

// TestRoutingDoesNotBypassTheNativeModeGate is the negative of AC#A: opening the
// door must not weaken what stands behind it. A widening request with no L2
// confirmation leg attached is ticket 114 AC#2's refusal, reached through the
// router.
func TestRoutingDoesNotBypassTheNativeModeGate(t *testing.T) {
	writer := &hop33ModeWriter{current: risk.ModeAskHighRisk, readable: true}
	audit := &hop33Audit{}
	d := &ComposerDispatch{Mode: &ModeWriteHandler{Modes: writer, Audit: audit.fn()}}
	raw := hop33Raw(t, map[string]any{
		"method": MethodModeRequest, "requestId": "rid-widen", "source": ComposerRequestSource,
		"to": risk.ModeAutoApproveName,
	})
	reply, err := d.Handle(context.Background(), raw)
	if !errors.Is(err, ErrNoL2Confirm) {
		t.Fatalf("widening through the router returned %v, want %v: the hop must not skip ticket 114's gate",
			err, ErrNoL2Confirm)
	}
	if writer.setCalls != 0 {
		t.Errorf("write leg called %d time(s) on an unconfirmed widening request, want 0", writer.setCalls)
	}
	if reply == "" {
		t.Error("refusal reached the user as the empty string: a dropped request that looks handled")
	}
}

// ---------------------------------------------------- AC#B: default is refuse

// TestUnlistedMethodNameIsRefusedAndAudited is judgement 2's outside half: a
// method the roster never named must come back as a refusal the user can read,
// with an audit line, and must run nothing.
func TestUnlistedMethodNameIsRefusedAndAudited(t *testing.T) {
	spy := &hop33Capture{}
	audit := &hop33Audit{}
	d := &ComposerDispatch{Mode: spy, Audit: audit.fn()}
	raw := hop33Raw(t, map[string]any{
		"method": "panel.not.on.the.roster", "requestId": "rid-ac-b", "source": ComposerRequestSource,
	})
	reply, err := d.Handle(context.Background(), raw)
	if !errors.Is(err, ErrComposerRequest) {
		t.Fatalf("unlisted method returned %v, want a refusal wrapping %v", err, ErrComposerRequest)
	}
	if !strings.Contains(reply, "rid-ac-b") {
		t.Errorf("refusal %q does not name the requestId, so the user cannot trace their own click", reply)
	}
	if spy.calls != 0 {
		t.Errorf("an unlisted method reached a handler %d time(s), want 0", spy.calls)
	}
	if len(audit.lines) != 1 {
		t.Fatalf("audit lines = %d, want 1: a refusal with no line is a silent drop wearing a return value",
			len(audit.lines))
	}
}

// TestRosterMismatchBackstopRefusesInsteadOfAccepting is judgement 2's actual
// mutation target. The default branch is unreachable from Handle today (parse
// refuses the name first), so this drives dispatch directly - the only way to
// make "change the default to allow" observable at all. Vehicle: in
// composer_dispatch.go replace the default branch's `return d.rosterMismatch(req)`
// with `return nil`, and this test goes red while every other nail in this file
// still passes, which is what makes it a real nail and not a decoration.
func TestRosterMismatchBackstopRefusesInsteadOfAccepting(t *testing.T) {
	audit := &hop33Audit{}
	d := &ComposerDispatch{Audit: audit.fn()}
	req := ComposerRequest{Method: "panel.a.future.fifth.method", RequestID: "rid-backstop", Source: ComposerRequestSource}
	err := d.dispatch(context.Background(), req)
	if !errors.Is(err, ErrRosterMismatch) {
		t.Fatalf("default branch returned %v, want %v: an unhandled whitelisted name must never be accepted",
			err, ErrRosterMismatch)
	}
	if len(audit.lines) != 1 {
		t.Fatalf("default branch wrote %d audit line(s), want 1", len(audit.lines))
	}
}

// TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped covers the three doors
// this tree has no handler behind yet. The router must not treat "I have no code
// for that" as success.
func TestEveryListedMethodWithoutAHandlerIsRefusedNotDropped(t *testing.T) {
	cases := []struct {
		name   string
		method string
		build  func(audit AuditFunc, mode ModeRequestHandler) *ComposerDispatch
	}{
		{"workspace", MethodWorkspaceRequest, func(a AuditFunc, m ModeRequestHandler) *ComposerDispatch {
			return &ComposerDispatch{Mode: m, Audit: a}
		}},
		{"attachment", MethodAttachmentAdd, func(a AuditFunc, m ModeRequestHandler) *ComposerDispatch {
			return &ComposerDispatch{Mode: m, Audit: a}
		}},
		{"message", MethodMessageSend, func(a AuditFunc, m ModeRequestHandler) *ComposerDispatch {
			return &ComposerDispatch{Mode: m, Audit: a}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spy := &hop33Capture{}
			d := tc.build((&hop33Audit{}).fn(), spy)
			raw := hop33Raw(t, map[string]any{
				"method": tc.method, "requestId": "rid-" + tc.name, "source": ComposerRequestSource,
			})
			reply, err := d.Handle(context.Background(), raw)
			if !errors.Is(err, ErrNoHandlerAttached) {
				t.Fatalf("unhandled %s request returned %v, want %v", tc.method, err, ErrNoHandlerAttached)
			}
			if reply == "" {
				t.Errorf("unhandled %s request produced no sentence for the user", tc.method)
			}
			if spy.calls != 0 {
				t.Errorf("the mode handler ran %d time(s) for a %s request, want 0: misrouting", spy.calls, tc.method)
			}
		})
	}
}

// TestAttachedHandlerIsWhatRunsForItsOwnMethod proves the four doors are four,
// not one: an attached workspace handler must not land in the mode handler.
func TestAttachedHandlerIsWhatRunsForItsOwnMethod(t *testing.T) {
	modeSpy := &hop33Capture{}
	other := &hop33RefusingHandler{}
	d := &ComposerDispatch{Mode: modeSpy, Workspace: other, Attachment: other, Message: other}
	for _, m := range []string{MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend} {
		raw := hop33Raw(t, map[string]any{"method": m, "requestId": "rid-" + m, "source": ComposerRequestSource})
		if _, err := d.Handle(context.Background(), raw); err == nil {
			t.Errorf("%s: the stub handler refused but Handle reported success", m)
		}
	}
	if modeSpy.calls != 0 {
		t.Errorf("mode handler called %d time(s) for three non-mode requests, want 0", modeSpy.calls)
	}
	if other.calls != 3 {
		t.Errorf("attached handlers called %d time(s), want 3 (one per door)", other.calls)
	}
}

// ------------------------------------------------------------ AC#C: source

// TestTamperedSourceIsRefusedByTheSameDoor is judgement 3. It pins two things:
// the refusal, and the IDENTITY of the door - the error wraps bridge.go's own
// ErrComposerRequest, the same sentinel the whitelist refusal wraps, and the
// reason names the source. That is the proof this file did not build a second
// gate: a second gate would carry its own error.
func TestTamperedSourceIsRefusedByTheSameDoor(t *testing.T) {
	spy := &hop33Capture{}
	audit := &hop33Audit{}
	d := &ComposerDispatch{Mode: spy, Audit: audit.fn()}
	// "panel-composer " (trailing space) is deliberately NOT here: bridge.go:93
	// TrimSpace's the source before comparing, so that spelling is accepted by
	// the existing gate. Tolerating it is bridge.go's decision, not one this
	// slice gets to quietly tighten or widen.
	for _, src := range []string{"cli", "ball", "", "  ", "PANEL-COMPOSER"} {
		raw := hop33Raw(t, map[string]any{
			"method": MethodModeRequest, "requestId": "rid-ac-c", "source": src,
			"to": risk.ModeAskEveryStepName,
		})
		_, err := d.Handle(context.Background(), raw)
		if !errors.Is(err, ErrComposerRequest) {
			t.Fatalf("source %q accepted or refused by the wrong door: %v", src, err)
		}
		if !strings.Contains(err.Error(), "来源") {
			t.Errorf("source %q refused for %v, which does not name the source: a whitelist refusal "+
				"would be indistinguishable from an attribution refusal", src, err)
		}
		if spy.calls != 0 {
			t.Fatalf("source %q reached a handler %d time(s), want 0", src, spy.calls)
		}
	}
	if len(audit.lines) != 5 {
		t.Fatalf("audit lines = %d, want 5 (one per tampered source): refusal must not be silent", len(audit.lines))
	}
}

// TestAttributionFieldsReachTheHandlerUnchanged is judgement 3's other half: the
// fields that say who asked arrive as they were sent. Vehicle: insert
// `req.Source = ComposerRequestSource` as the first line of dispatch (a router
// that "fixes up" attribution), and this test goes red.
func TestAttributionFieldsReachTheHandlerUnchanged(t *testing.T) {
	spy := &hop33Capture{}
	d := &ComposerDispatch{Mode: spy}
	raw := hop33Raw(t, map[string]any{
		"method": MethodModeRequest, "requestId": "rid-77", "source": ComposerRequestSource,
		"to": risk.ModeAskEveryStepName,
	})
	if _, err := d.Handle(context.Background(), raw); err != nil {
		t.Fatalf("valid request refused: %v", err)
	}
	if spy.calls != 1 {
		t.Fatalf("handler called %d time(s), want 1", spy.calls)
	}
	if spy.got.Source != ComposerRequestSource {
		t.Errorf("handler saw source %q, want %q verbatim", spy.got.Source, ComposerRequestSource)
	}
	if spy.got.RequestID != "rid-77" {
		t.Errorf("handler saw requestId %q, want %q verbatim: the router rewrote the correlation key the audit trail depends on",
			spy.got.RequestID, "rid-77")
	}
}

// TestDispatcherSpellsNoRouteLiteralOfItsOwn is AC#C in static form, and it is
// also the nail that keeps this package's own AST instrument honest: a route
// string written here that knownComposerMethod answers is reported by
// poolJudgedByRealGuard as a second dispatch chain, so the vocabulary must come
// from bridge.go's constants alone.
func TestDispatcherSpellsNoRouteLiteralOfItsOwn(t *testing.T) {
	root := panelRepoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "internal", "panel", "composer_dispatch.go"))
	if err != nil {
		t.Fatalf("read composer_dispatch.go: %v", err)
	}
	routeLiteral := regexp.MustCompile(`"panel\.[^"]*"`)
	for i, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		if hit := routeLiteral.FindString(trimmed); hit != "" {
			t.Errorf("composer_dispatch.go:%d spells route %q as a literal instead of naming bridge.go's "+
				"constant: this is a second whitelist, and the guard's own case list would never have named it",
				i+1, hit)
		}
	}
}

// ----------------------------------------------- the honesty nails (AC#D / #E)

// TestSliceAAttachesNoHostAndNamesTheOpenWindowHops is 33-a1 §7.4's fourth
// judgement, restated (ruling A386, ticket 33 slice B) as a CAPABILITY ruler.
//
// What it asks now, and only this: does this tree attach a native host / WebView2
// message channel? It answers by looking for the symbol family such an attachment
// cannot be written without (CoreWebView2 / WebView2 / WebMessage identifiers in
// production Go sources, and a webview module in go.mod / go.sum), never for a
// type's NAME. The predecessor scanned the repo's text for "ComposerDispatch", so
// it read "a CLI seam calls the router" and "somebody wired the window and passed
// it off as a working panel" as the same event - a wording ruler that blocked the
// very listener ticket 33 slice B was dispatched to land.
//
// What green here means, stated exactly: no host is attached. It says NOTHING
// about whether the panel can click - the inbound router has a production caller
// today (cmd/wisp/panel_inbound.go, `wisp panel-inbound`) and this test is green
// with it, which is precisely the reading the ticket carries: 票 33 AC#9 says the
// panel still cannot be clicked, because H2 (the control is created), H3
// (WebMessageReceived takes the page's bytes into Go) and H10 (the reply goes back
// to the page) are all still open. When one of those lands, the first leg below
// goes red and the window's own ticket has to say so - "do not let a green here
// stand in for a panel that can click" is still the sentence this nail is for.
//
// It is not a door unless it can go red, so the test carries its own positive
// control: the SAME predicate is run over throwaway trees written OUTSIDE this
// repository (t.TempDir; a scan of the real disk cannot see a -overlay file, and
// a carrier that the ruler cannot see would make this a nail with no door).
// Two legs must fire (a fake host source, a fake webview dependency) and one leg
// must stay silent (the CLI-seam shape the old wording ruler reddened on) - see
// the sub-tests.
func TestSliceAAttachesNoHostAndNamesTheOpenWindowHops(t *testing.T) {
	root := panelRepoRoot(t)
	hits := hostChannelCapabilityHits(t, root)
	t.Logf("native-host capability hits in %s: %d -> %v", root, len(hits), hits)
	if len(hits) != 0 {
		t.Errorf("a native host / WebView2 message channel is attached in this tree (%v) while ticket 33's "+
			"real-window ACs are still unticked: that is H2/H3/H10 landing - say so on the ticket, do not let a "+
			"green here stand in for a panel that can click", hits)
	}

	t.Run("positive-control", func(t *testing.T) {
		// Carrier A: a production source that hosts the control. Nothing here is
		// imported by the product and nothing here is a *_test.go file, which is
		// what makes it the shape the predicate is supposed to catch.
		a := t.TempDir()
		writeHostCarrier(t, a, filepath.Join("internal", "panel"), "host_windows.go", `package panel

type CoreWebView2 struct{}

func (w *CoreWebView2) WebMessageReceived(sender, args any) error { return nil }

func (w *CoreWebView2) PostWebMessageAsJson(json string) error { return nil }
`)
		if got := hostChannelCapabilityHits(t, a); len(got) == 0 {
			t.Errorf("POSITIVE CONTROL RED (the ruler, not the product): a tree whose production source hosts a "+
				"CoreWebView2 and takes WebMessageReceived answered 0 capability hits, so the check above is a "+
				"door with no latch. Hits would have named the file and the symbol: got %v", got)
		} else {
			t.Logf("carrier A (fake host source) -> %d hit(s) %v", len(got), got)
		}

		// Carrier B: the dependency half. A host is normally bought, not written,
		// so go.mod / go.sum are read too.
		b := t.TempDir()
		writeHostCarrier(t, b, "", "go.mod", "module github.com/CarlosShao/wisp\n\nrequire github.com/jchv/go-webview2 v0.0.0-20220111092826-5ea30c1aaa37\n")
		writeHostCarrier(t, b, "", "go.sum", "github.com/jchv/go-webview2 v0.0.0-20220111092826-5ea30c1aaa37 h1:8vXupf01+8Q7j94hb8jLZ4J0K+zxuu6W1JxG2rHpNzA=\n")
		if got := hostChannelCapabilityHits(t, b); len(got) == 0 {
			t.Errorf("POSITIVE CONTROL RED (the ruler, not the product): a tree whose go.mod requires a webview "+
				"binding answered 0 capability hits, so the dependency half of the check is not wired. got %v", got)
		} else {
			t.Logf("carrier B (webview dependency) -> %d hit(s) %v", len(got), got)
		}

		// Carrier C is the discrimination leg, and the reason this rewrite is a
		// narrower instrument rather than a loosened one: the shape that reddened
		// the wording ruler - a non-test production file that constructs a
		// ComposerDispatch and sends it Handle, with no window anywhere - must
		// stay silent under the capability ruler. Reddening on it again would put
		// this file back in the way of ticket 33's own inbound leg.
		c := t.TempDir()
		writeHostCarrier(t, c, filepath.Join("cmd", "wisp"), "panel_inbound.go", `package main

import "github.com/CarlosShao/wisp/internal/panel"

func run() {
	disp := &panel.ComposerDispatch{}
	_, _ = disp.Handle(nil, "{}")
}
`)
		if got := hostChannelCapabilityHits(t, c); len(got) != 0 {
			t.Errorf("NEGATIVE CONTROL RED (the ruler is still a wording ruler): an inbound CLI seam with no host "+
				"answered %d capability hit(s) %v - only the presence of a native host / WebView2 message channel "+
				"may redden this test", len(got), got)
		} else {
			t.Logf("carrier C (CLI seam, no host) -> 0 hits, as required")
		}
	})
}

// hostChannelSymbols are the identifier anchors of a WebView2 host attachment.
// They are substrings because the generated COM bindings carry them inside longer
// names (add_WebMessageReceived, CreateCoreWebView2EnvironmentWithOptions,
// get_WebMessageAsJson), and a host written through those spellings is exactly
// what this ruler is for. Case-sensitive on purpose: the camel spelling is what
// the binding uses, while prose in this repository says "WebView2 hosting" in
// comments - and comments are not identifiers, which is the whole reason this
// predicate parses sources instead of grepping them (measured 2026-09-28: 27
// production .go files describe the missing host in prose; reddening on prose
// would have made this test red since before slice A existed).
var hostChannelSymbols = []string{"CoreWebView2", "WebView2", "WebMessage"}

// hostModuleTokens are the go.mod / go.sum half of the same question: a WebView2
// host is normally imported, not hand-written. docs/specs and the local spike
// (scripts/spike/webview2-latency) name the binding as github.com/jchv/go-webview2,
// so a module path containing "webview" is the signal. Case-insensitive here,
// because module paths are lower-case by convention and a renamed fork still says
// what it is.
var hostModuleTokens = []string{"webview", "msedge"}

// hostChannelCapabilityHits reads root's PRODUCTION sources for the host family
// and returns "relpath:line:what" per hit. Out of scope, and named so nobody has
// to guess: test files (a host is a shipping thing, and this very file spells the
// symbols in its own carrier strings), the directory set below, comments, and
// string literals. A host reached through a hand-rolled COM vtable that spells
// none of these names in an identifier is outside this ruler's sight - the
// dependency half is what would catch it, and 票 33's window leg is where that
// shape would actually arrive.
func hostChannelCapabilityHits(t *testing.T, root string) []string {
	t.Helper()
	var hits []string
	fset := token.NewFileSet()
	walk := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			// .git and the build trees carry no product source; scripts/ holds this
			// repository's own WebView2 latency spike, which is deliberately not the
			// shipped host, and the set below is the one the predecessor skipped -
			// frontend/ and design/ are another session's ground and this slice reads
			// neither. .scratch/ is the ticket pool: other slices keep probe copies
			// there (several of them renamed copies of run.go), and a red from one of
			// those would accuse the product of something a scratch tree did.
			case ".git", ".scratch", "node_modules", "dist", "third_party", "scripts", "docs", "frontend", "design", "build":
				return filepath.SkipDir
			}
			return nil
		}
		name := d.Name()
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			src, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			f, err := parser.ParseFile(fset, p, src, 0)
			if err != nil {
				// Unparseable production source is itself a finding: the alternative
				// is a ruler that quietly skips the file it cannot read.
				hits = append(hits, rel+":0:unparseable: "+err.Error())
				return nil
			}
			for _, imp := range f.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				if hostModulePathHit(path) {
					hits = append(hits, fmt.Sprintf("%s:%d:imports %s", rel, fset.Position(imp.Pos()).Line, path))
				}
			}
			ast.Inspect(f, func(n ast.Node) bool {
				id, ok := n.(*ast.Ident)
				if !ok {
					return true
				}
				if sym := hostChannelSymbolHit(id.Name); sym != "" {
					hits = append(hits, fmt.Sprintf("%s:%d:identifier %s", rel, fset.Position(id.Pos()).Line, sym))
				}
				return true
			})
			return nil
		}
		if name == "go.mod" || name == "go.sum" {
			src, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			for i, line := range strings.Split(string(src), "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "//") {
					continue
				}
				for _, tok := range hostModuleTokens {
					if strings.Contains(strings.ToLower(line), tok) {
						hits = append(hits, fmt.Sprintf("%s:%d:dependency %q", rel, i+1, tok))
						break
					}
				}
			}
		}
		return nil
	})
	if walk != nil {
		t.Fatalf("scan %s for host capability: %v", root, walk)
	}
	sort.Strings(hits)
	return hits
}

// hostChannelSymbolHit returns the anchor an identifier carries, or "" when the
// identifier says nothing about hosting.
func hostChannelSymbolHit(name string) string {
	for _, sym := range hostChannelSymbols {
		if strings.Contains(name, sym) {
			return sym
		}
	}
	return ""
}

func hostModulePathHit(path string) bool {
	lower := strings.ToLower(path)
	for _, tok := range hostModuleTokens {
		if strings.Contains(lower, tok) {
			return true
		}
	}
	return false
}

// writeHostCarrier puts one throwaway file into a repository-external temp tree.
// t.TempDir hands back a directory outside the repository and removes it itself,
// which is why this file needs no delete step and why no carrier can be mistaken
// for product source.
func writeHostCarrier(t *testing.T, root, dir, name, content string) {
	t.Helper()
	where := filepath.Join(root, dir)
	if err := os.MkdirAll(where, 0o755); err != nil {
		t.Fatalf("carrier dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(where, name), []byte(content), 0o600); err != nil {
		t.Fatalf("carrier file: %v", err)
	}
}

// TestTheInboundHopAddsNoSwitchingCapability is AC#E in runnable form: ticket
// 92 AC#7's ruler scans this package's production files for the cut feature's
// entry points, and slice A must not be the thing that reintroduces them.
func TestTheInboundHopAddsNoSwitchingCapability(t *testing.T) {
	root := panelRepoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "internal", "panel", "composer_dispatch.go"))
	if err != nil {
		t.Fatalf("read composer_dispatch.go: %v", err)
	}
	for _, banned := range []string{"checkoutBranch", "changeRepo", "repoPicker", "branchSelect", "vcs.switch"} {
		if strings.Contains(string(data), banned) {
			t.Errorf("composer_dispatch.go contains %q: the switching dimension is ticket 186's, gated on Q-69", banned)
		}
	}
}
