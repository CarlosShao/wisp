package ball

// Ticket-65 integration leg (A1594): deterministic tests for the seam that
// can be tested without a window - the asset, the state mapping, the script
// builders, and the debug face. The window itself is exercised by the winlive
// suite (orb_live_windows_test.go), which is where the real WebView2 runtime
// and the real GPU live. None of these tests skip, ever.

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// orbUpstreamSHA256 is the sha256 of the embedded upstream export page
// (assets/orb-siri.html, LF-normalised as shipped), pinned so no later leg
// can edit the upstream file and keep the marker test above green.
const orbUpstreamSHA256 = "d93da649d11bc4bc3ab91c006777bbadce1c2eb7d0ffc9d91608ad8d866e3bc6"

// TestOrbAssetShippedUnmodified pins the integration's core claim: the embed
// is LerSent001/orb's own export page, byte-for-byte (CRLF->LF aside), NOT a
// rewrite. The owner's rejection of the hand-written D2D orb ("拉裤兜子的复刻")
// is answered by shipping the upstream page; this test is the ruler that keeps
// a later leg from quietly swapping the asset for a local imitation.
func TestOrbAssetShippedUnmodified(t *testing.T) {
	if len(orbPageHTML) == 0 {
		t.Fatal("embedded orb page is empty")
	}
	if !bytes.HasPrefix(orbPageHTML, []byte("<!doctype html>")) {
		t.Fatalf("embedded orb page does not start with the upstream doctype: %q", orbPageHTML[:32])
	}
	if bytes.Contains(orbPageHTML, []byte("\r\n")) {
		t.Fatal("embedded orb page still contains CRLF; the asset must be LF")
	}
	page := string(orbPageHTML)
	for _, marker := range []string{
		"defineProperty(window, \"liquidOrb\"", // the upstream interface object
		"function setState(nextState)",         // the upstream state entry
		"function setAudioBands(bands",         // the upstream audio entry
		`const stateSeeds = {"idle":`,          // the idle seed
		`"thinking":[`,                         // the thinking seed
		"@fragment",                            // real WGSL, not a stand-in
		"requestAdapter",                       // real WebGPU boot path
		"Animated liquid glass orb",            // upstream canvas label
	} {
		if !strings.Contains(page, marker) {
			t.Errorf("upstream marker %q missing from the embedded orb page", marker)
		}
	}
	// The page must stay offline: no external reference of any kind (the CSP
	// with connect-src 'none' enforces it at runtime; this is the asset-side
	// half of the same rule).
	for _, bad := range []string{"http://", "https://", "src=\"http", "fetch("} {
		if strings.Contains(page, bad) {
			t.Errorf("embedded orb page contains an external reference %q", bad)
		}
	}
}

// TestOrbStateMappingCoversAll20States walks every frozen D43 state and
// requires the mapping to land on one of the page's own two state names -
// the page throws on anything else, so a miss here is a runtime error in
// front of the user, not a cosmetic one.
func TestOrbStateMappingCoversAll20States(t *testing.T) {
	all := []statemachine.State{
		statemachine.StateFirstRun, statemachine.StateSleeping, statemachine.StateArmed,
		statemachine.StateMuted, statemachine.StateListening, statemachine.StateThinking,
		statemachine.StateActing, statemachine.StateSpeaking, statemachine.StateWarm,
		statemachine.StateConversation, statemachine.StateConfirming,
		statemachine.StateAwaitingApproval, statemachine.StateSettling,
		statemachine.StateDownloading, statemachine.StateUnconfigured,
		statemachine.StateNoNetwork, statemachine.StateError, statemachine.StateQueued,
		statemachine.StateStuck, statemachine.StateWatchdogAlert,
	}
	if len(all) != 20 {
		t.Fatalf("expected the 20 frozen states, got %d", len(all))
	}
	for _, s := range all {
		if !statemachine.Valid(s) {
			t.Fatalf("state %q is not valid per the frozen table", s)
		}
		got := orbStateFor(s)
		if got != "idle" && got != "thinking" {
			t.Errorf("orbStateFor(%q) = %q, want one of the page's own names (idle|thinking)", s, got)
		}
	}
}

// TestOrbStateMappingWorkShape pins the mapping DECISIONS, not just the
// closure: the work family lands on "thinking" (the orb's activation profile)
// and the standby/failure family on "idle". The failure family maps to idle
// by design - the D2D fallback carries the dedicated failure recipes, and the
// orb has no error face of its own (A1594 keeps upstream unmodified).
func TestOrbStateMappingWorkShape(t *testing.T) {
	thinking := []statemachine.State{
		statemachine.StateListening, statemachine.StateThinking, statemachine.StateActing,
		statemachine.StateSpeaking, statemachine.StateConversation, statemachine.StateConfirming,
		statemachine.StateAwaitingApproval, statemachine.StateQueued,
	}
	for _, s := range thinking {
		if got := orbStateFor(s); got != "thinking" {
			t.Errorf("orbStateFor(%q) = %q, want thinking", s, got)
		}
	}
	idle := []statemachine.State{
		statemachine.StateSleeping, statemachine.StateArmed, statemachine.StateMuted,
		statemachine.StateWarm, statemachine.StateSettling, statemachine.StateDownloading,
		statemachine.StateFirstRun, statemachine.StateUnconfigured, statemachine.StateNoNetwork,
		statemachine.StateError, statemachine.StateStuck, statemachine.StateWatchdogAlert,
	}
	for _, s := range idle {
		if got := orbStateFor(s); got != "idle" {
			t.Errorf("orbStateFor(%q) = %q, want idle", s, got)
		}
	}
}

// TestOrbAudioFamilyIsAudioDriven pins WHICH states let the mic move the orb:
// the voice family only, and never Sleeping (a hot mic must not wake a
// sleeping orb - the same law the D2D liquid already enforces).
func TestOrbAudioFamilyIsAudioDriven(t *testing.T) {
	driven := []statemachine.State{
		statemachine.StateListening, statemachine.StateSpeaking, statemachine.StateConversation,
	}
	for _, s := range driven {
		if !orbAudioFor(s) {
			t.Errorf("orbAudioFor(%q) = false, want true (voice family drives the bands)", s)
		}
	}
	notDriven := []statemachine.State{
		statemachine.StateSleeping, statemachine.StateThinking, statemachine.StateError,
		statemachine.StateAwaitingApproval,
	}
	for _, s := range notDriven {
		if orbAudioFor(s) {
			t.Errorf("orbAudioFor(%q) = true, want false", s)
		}
	}
}

// TestOrbStateScript pins the ExecuteScript body against the page's frozen
// object name, including the guard for a page that never came up.
func TestOrbStateScript(t *testing.T) {
	got := orbStateScript("thinking")
	want := `window.liquidOrb&&window.liquidOrb.setState("thinking")`
	if got != want {
		t.Errorf("orbStateScript = %q, want %q", got, want)
	}
}

// TestOrbAudioScript pins the band payload shape and the clamping: NaN to 0,
// negatives to 0, over-range to 1, and every band a finite fixed-point
// literal (the page clamps again, but a NaN would slip through its chain).
func TestOrbAudioScript(t *testing.T) {
	got := orbAudioScript(0.5, -2, 7, nan64())
	want := "window.liquidOrb&&window.liquidOrb.setAudioBands({low:0.5000,mid:0.0000,high:1.0000,all:0.0000})"
	if got != want {
		t.Errorf("orbAudioScript = %q, want %q", got, want)
	}
}

// nan64 keeps the math/bits-free NaN literal out of the table above.
func nan64() float64 { return zero64() / zero64() }

func zero64() float64 { return 0 }

// TestOrbDebugLine pins the debug face the harness prints.
func TestOrbDebugLine(t *testing.T) {
	if got := OrbDebugLine(true, statemachine.StateThinking); got != "orb-surface=webview2-orb orb-state=thinking" {
		t.Errorf("OrbDebugLine(attached) = %q", got)
	}
	if got := OrbDebugLine(false, statemachine.StateSleeping); got != "orb-surface=d2d-fallback orb-state=idle" {
		t.Errorf("OrbDebugLine(fallback) = %q", got)
	}
	if OrbSurfaceName(true) != "webview2-orb" || OrbSurfaceName(false) != "d2d-fallback" {
		t.Error("OrbSurfaceName names drifted")
	}
}

// TestOrbServingConstants pins the serving seam AFTER the black-orb fix
// (ticket 65): the site is served under the synthetic https://wisp.orb origin
// (host + entry URL spelled in orbasset.go), the document handed to BOTH
// delivery shapes (virtual-host folder + NavigateToString fallback) carries
// the CSP as an injected <meta http-equiv> with the exfiltration cut intact,
// and the upstream embed itself stays byte-identical to the pinned sha256
// (the injection happens at serve time, never in the asset file).
func TestOrbServingConstants(t *testing.T) {
	if orbSiteHost != "wisp.orb" {
		t.Errorf("orbSiteHost = %q", orbSiteHost)
	}
	if orbOrigin != "https://wisp.orb" {
		t.Errorf("orbOrigin = %q", orbOrigin)
	}
	if orbEntryURL != "https://wisp.orb/orb.html" {
		t.Errorf("orbEntryURL = %q", orbEntryURL)
	}
	doc := string(orbPageDocument())
	if got := strings.Count(doc, "http-equiv=\"Content-Security-Policy\""); got != 1 {
		t.Errorf("served document carries %d CSP meta tags, want exactly 1", got)
	}
	for _, must := range []string{"connect-src 'none'", "default-src 'none'"} {
		if !strings.Contains(doc, must) {
			t.Errorf("served document CSP missing %q", must)
		}
	}
	// The meta must sit inside the head, before the upstream style block.
	headEnd := strings.Index(doc, "<style>")
	metaEnd := strings.Index(doc, "Content-Security-Policy")
	if headEnd < 0 || metaEnd < 0 || metaEnd > headEnd {
		t.Error("CSP meta is not inside the served document's head")
	}
	// The RAW embed must NOT already contain the meta: the injection is this
	// repo's serving seam, and the asset file remains the upstream export.
	if strings.Contains(string(orbPageHTML), "http-equiv") {
		t.Error("the embedded asset itself was modified with the CSP meta")
	}
	sum := sha256.Sum256(orbPageHTML)
	if got := fmt.Sprintf("%x", sum); got != orbUpstreamSHA256 {
		t.Errorf("embedded orb page sha256 = %s, want the pinned upstream %s", got, orbUpstreamSHA256)
	}
	if !strings.Contains(orbPageBootstrap, `b.style.background = "transparent"`) {
		t.Error("orb bootstrap does not clear the body background")
	}
}

// TestOrbWriteSite pins the materialization seam: the folder the virtual-host
// mapping points at receives exactly the document the fallback path would
// hand NavigateToString - one file, the entry name, meta CSP included.
func TestOrbWriteSite(t *testing.T) {
	dir := t.TempDir()
	if err := orbWriteSite(dir); err != nil {
		t.Fatalf("orbWriteSite: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, orbEntryName))
	if err != nil {
		t.Fatalf("site entry file missing: %v", err)
	}
	if !bytes.Equal(got, orbPageDocument()) {
		t.Error("site entry file differs from the served document")
	}
}
