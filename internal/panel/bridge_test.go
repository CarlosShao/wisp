package panel

// The composer request envelope (ticket 92): one pipe, one identity claim, no
// capability smuggled through either.
//
// The point of these cases is the shape of the door, not the payload: the
// approval decision must stay unaddressable from the panel side. A request that
// names a decision verb, a request with no requestId (so nothing can be traced
// back to a user action), or a request that claims a different source all get the
// same treatment - a refusal whose reason is readable, and no state change.

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// envelopeJSON builds one composer request the way the frontend does: the
// identity fields are added through the same json.Marshal path, so a case that
// means "omit requestId" omits the KEY rather than inventing a second decoder.
func envelopeJSON(t *testing.T, method, requestID, source string, extra map[string]any) string {
	t.Helper()
	body := map[string]any{}
	if method != "" {
		body["method"] = method
	}
	if requestID != "" {
		body["requestId"] = requestID
	}
	if source != "" {
		body["source"] = source
	}
	for k, v := range extra {
		body[k] = v
	}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(raw)
}

func okEnvelope(t *testing.T, m string, extra map[string]any) string {
	t.Helper()
	return envelopeJSON(t, m, NewRequestID(), ComposerRequestSource, extra)
}

func TestComposerEnvelopeAcceptsItsFourRequests(t *testing.T) {
	cases := []struct {
		method string
		extra  map[string]any
	}{
		{MethodModeRequest, map[string]any{"to": "auto_approve"}},
		{MethodWorkspaceRequest, map[string]any{"path": `D:\work\Wisp`}},
		{MethodAttachmentAdd, map[string]any{
			"name": "a.png", "declaredMime": "image/png",
			"sizeBytes": 16, "dataBase64": "iVBORw0KGgo=",
		}},
		{MethodMessageSend, map[string]any{"text": "看看这个"}},
	}
	for _, tc := range cases {
		raw := okEnvelope(t, tc.method, tc.extra)
		r, err := ParseComposerRequest(raw)
		if err != nil {
			t.Errorf("%s refused: %v", tc.method, err)
			continue
		}
		if r.Source != ComposerRequestSource || r.RequestID == "" {
			t.Errorf("%s parsed without its identity fields: %+v", tc.method, r)
		}
	}
}

func TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests(t *testing.T) {
	t.Run("a request that names an approval decision is not a composer method", func(t *testing.T) {
		// A hostile renderer's most valuable string: the decision verb, aimed at
		// the card. It must be refused at the door, and refused as "not a
		// capability of this route" rather than routed anywhere.
		raw := okEnvelope(t, "approval.decide", map[string]any{"outcome": "grant"})
		if _, err := ParseComposerRequest(raw); err == nil ||
			!errors.Is(err, ErrComposerRequest) ||
			!strings.Contains(err.Error(), "不是面板 composer 通路") {
			t.Fatalf("ParseComposerRequest(%s) = %v, want a refusal naming the route", raw, err)
		}
	})

	t.Run("a missing requestId is refused", func(t *testing.T) {
		raw := `{"method":"panel.mode.request","source":"panel-composer","to":"auto_approve"}`
		if _, err := ParseComposerRequest(raw); err == nil ||
			!strings.Contains(err.Error(), "requestId") {
			t.Fatalf("a request with no requestId parsed as %v", err)
		}
	})

	t.Run("a forged or absent source is refused", func(t *testing.T) {
		for _, src := range []string{"", "ball", "native", "panel-approval", "PANEL-COMPOSER"} {
			body := `{"method":"panel.workspace.request","requestId":"r1","source":` +
				jsonString(src) + `,"path":"D:\\work"}`
			if _, err := ParseComposerRequest(body); err == nil {
				t.Errorf("source %q was accepted - only %q may use this route", src, ComposerRequestSource)
			} else if !strings.Contains(err.Error(), "伪造") {
				t.Errorf("source %q refused for the wrong reason: %v", src, err)
			}
		}
	})

	t.Run("undecodable and non-object payloads are refused, not defaulted", func(t *testing.T) {
		for _, raw := range []string{"", "{", `[]`, `null`, `"panel.mode.request"`} {
			if _, err := ParseComposerRequest(raw); err == nil {
				t.Errorf("payload %q parsed without error", raw)
			}
		}
	})

	t.Run("the refusal text keeps the request identifiable", func(t *testing.T) {
		line := RefusedEnvelopeForUser(ComposerRequest{Method: "panel.mode.request", RequestID: "pc-7-x"},
			errors.New("来源不合法"))
		if !strings.Contains(line, "pc-7-x") || !strings.Contains(line, "panel.mode.request") {
			t.Errorf("refusal %q does not name what was refused", line)
		}
	})
}

// TestFrontendComposerRequestsMatchTheEnvelope reads the real frontend source
// and checks the four method names, the source claim and the requestId are
// emitted on the ONE envelope - so the parser above and the page below cannot
// drift apart without a red test, and no second IPC can appear unnoticed.
func TestFrontendComposerRequestsMatchTheEnvelope(t *testing.T) {
	root := panelRepoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "frontend", "src", "lib", "panel.ts"))
	if err != nil {
		t.Fatalf("read frontend/src/lib/panel.ts: %v", err)
	}
	text := string(data)
	for _, m := range []string{
		MethodModeRequest, MethodWorkspaceRequest,
		MethodAttachmentAdd, MethodMessageSend,
	} {
		if !strings.Contains(text, `"`+m+`"`) {
			t.Errorf("frontend never emits %q - the Go side parses a method nobody sends", m)
		}
	}
	for _, key := range []string{`requestId`, `"source"`, `ComposerRequestSource`, `"panel-composer"`} {
		if key == `ComposerRequestSource` || key == `"panel-composer"` {
			if !strings.Contains(text, "panel-composer") {
				t.Errorf("frontend does not claim the source %q on its envelope", ComposerRequestSource)
			}
			continue
		}
		if !strings.Contains(text, strings.Trim(key, `"`)) {
			t.Errorf("frontend envelope is missing %s", key)
		}
	}
	if strings.Contains(text, "window.postMessage(") || strings.Contains(text, "new WebSocket") ||
		strings.Contains(text, "fetch(") {
		t.Error("frontend opened a second channel besides the WebView2 host bridge - ticket 92 forbids it")
	}
	if n := strings.Count(text, "bridge.postMessage"); n != 2 {
		t.Errorf("postMessage call sites = %d, want 2 (the approval request and sendRequest) - "+
			"a third one means a second envelope was invented", n)
	}
}

// TestEveryWhitelistedInboundNameIsAlsoSentByThePage is ticket 114 AC#8's
// forward direction as a gate: adding a name to the Go-side whitelist must make
// the page carry the matching literal, or this test goes red. Its subject is
// measured, not guessed: the cousin above checks four names written out by hand,
// so when ticket 248 added MethodConfigGet/MethodConfigSet ("config.get",
// "config.set") the whitelist grew and no instrument tied either name to the page
// (a2 census section 2, AC#8). Ticket 253's closed-roster ruler gates the Go-side
// registration spine - declared, guard, router - and reads nothing about the
// page, so it does not close this direction.
//
// The name list is read off bridge.go through the same prefix-agnostic extractor
// the closed-roster test uses, so a new constant cannot slip past by not being
// copied into a list here; the page is frontend/src/lib/panel.ts, the one door
// file the structural nails pin as the only place a host call may live
// (composer_test.go, TestTheRendererHoldsExactlyOneDoorToTheHost). The planted
// half runs the SAME check against a bridge.go copy carrying an unmirrored name,
// so a green above cannot come from a check that looks at nothing.
func TestEveryWhitelistedInboundNameIsAlsoSentByThePage(t *testing.T) {
	root := panelRepoRoot(t)
	bridgePath := filepath.Join(root, "internal", "panel", "bridge.go")
	pageData, err := os.ReadFile(filepath.Join(root, "frontend", "src", "lib", "panel.ts"))
	if err != nil {
		t.Fatalf("read frontend/src/lib/panel.ts: %v", err)
	}
	page := string(pageData)

	declared := parseFileInbound(t, bridgePath)
	if declared.parseError != "" {
		t.Fatalf("bridge.go: %s", declared.parseError)
	}
	names := declared.namesOf()
	if len(names) < 6 {
		t.Fatalf("bridge.go declares %d inbound names (%v) - an extractor reading almost nothing must not read as a pass", len(names), names)
	}
	if missing := namesAbsentFromPage(page, names); len(missing) > 0 {
		t.Errorf("the Go side answers %s and frontend/src/lib/panel.ts never spells %s as a literal: adding a route on one side only is the drift ticket 114 AC#8 turns into a gate - land the page literal in the same change or take the Go name out",
			strings.Join(names, ", "), strings.Join(missing, ", "))
	}

	src, err := os.ReadFile(bridgePath)
	if err != nil {
		t.Fatalf("read bridge.go: %v", err)
	}
	dir := t.TempDir()
	plantedPath := filepath.Join(dir, "bridge.go")
	plant := "\nconst ConfigUnmirrored114Test = \"config.unmirrored114\"\n"
	if err := os.WriteFile(plantedPath, append(append([]byte{}, src...), []byte(plant)...), 0o644); err != nil {
		t.Fatalf("plant bridge.go: %v", err)
	}
	planted := parseFileInbound(t, plantedPath)
	if planted.parseError != "" {
		t.Fatalf("planted bridge.go: %s", planted.parseError)
	}
	missing := namesAbsentFromPage(page, planted.namesOf())
	if len(missing) != 1 || missing[0] != "config.unmirrored114" {
		t.Errorf("a whitelist name the page never sends read as missing=%v, want exactly [\"config.unmirrored114\"] - this gate has no teeth for the shape ticket 114 AC#8 was filed over", missing)
	}
	t.Logf("page mirrors all %d whitelisted inbound names; the planted unmirrored name was reported (%v), so this check reads bridge.go rather than a list in this file",
		len(names), missing)
}

// namesAbsentFromPage reports which names never appear as a double-quoted string
// literal in the page source, keeping the input order (namesOf returns sorted).
func namesAbsentFromPage(page string, names []string) []string {
	var out []string
	for _, name := range names {
		if !strings.Contains(page, `"`+name+`"`) {
			out = append(out, name)
		}
	}
	return out
}

func jsonString(s string) string { b, _ := json.Marshal(s); return string(b) }
