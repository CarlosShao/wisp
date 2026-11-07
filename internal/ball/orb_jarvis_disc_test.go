package ball

// Ticket 65 jarvis dark-disc leg (2026-11-07): the owner's ruling - the jarvis
// orb gets its OWN dark circular backdrop (a disc) so the cyan particles keep
// the upstream demo's contrast over the owner's light desktop and windows,
// WITHOUT breaking the transparency architecture (the window corners must keep
// showing the desktop). The disc lives INSIDE the served page (#jarvis-disc,
// border-radius 50%, the demo's dark-blue #050913 family) UNDER the particle
// canvas (#root). These tests pin the disc's static shape in the EMBED - the
// bytes the WebView2 actually loads; the live composited-pixel proof rides
// winlive (orb_jarvis_disc_live_windows_test.go).
//
// Why the disc can only live in the page (the 甲 ruling, recorded in the
// evidence doc): the z-law (orb_windows.go follow()) slots the ball window
// ABOVE the orb window, so a D2D-painted disc would sit ON TOP of the
// particles - the page is the only layer beneath them.

import (
	"regexp"
	"strings"
	"testing"
)

// The disc's geometry constants, shared with the winlive pixel proof. The
// template ships the disc at 96% of the window edge (radius 48%): the corner
// probes (3 px inside the rect) sit on the diagonal at 66.5% of the edge on
// the 100 px debug window (70.7% asymptotically) - far outside the disc's
// 48%-of-edge radius - so the corners keep reading the desktop through real
// transparency. Coverage math from the upstream component
// (jarvis-ai-web-animation dist): PerspectiveCamera(35deg) at z=4.1 with
// sceneScale 0.82 puts the widest particle element (orbit radius 1.17 rings;
// particle shell radius 1.05 with its 0.88-1.14 jitter) at 76% of the canvas
// half-edge - the disc radius (96% of the half-edge) covers it all, so the
// AdditiveBlended cyan never composites over the desktop again.
const (
	jarvisDiscDiameterPct = 96   // #jarvis-disc width/height, % of the window edge
	jarvisDiscRadiusFrac  = 0.48 // disc radius as a fraction of the window edge
	jarvisDiscProbeFrac   = 0.42 // live disc-sample ring radius, fraction of the edge
	jarvisDiscProbeCount  = 8    // live disc-sample ring points
	jarvisDiscMinSamples  = 4    // live: at least this many non-particle ring samples required
	jarvisDiscLumFloor    = 60   // live: accepted disc samples must average darker than this
)

// TestJarvisAssetShipsDarkDisc pins the disc in the embedded asset: the exact
// bytes the serving seam (orbJarvisDocument) hands the WebView2.
func TestJarvisAssetShipsDarkDisc(t *testing.T) {
	page := string(jarvisPageHTML)
	if !strings.Contains(page, `id="jarvis-disc"`) {
		t.Fatal(`embedded jarvis page lost the #jarvis-disc element (the dark disc backdrop)`)
	}
	if !strings.Contains(page, "border-radius: 50%") {
		t.Fatal("embedded jarvis page disc lost its circular shape (border-radius: 50%)")
	}
	if !strings.Contains(page, "#050913") {
		t.Fatal("embedded jarvis page disc lost the demo's dark-blue family (#050913)")
	}
	// The disc must NOT come back as a flat body backdrop: that form breaks the
	// transparency architecture (the corners would stop reading the desktop).
	if strings.Contains(page, "background: #050913") {
		t.Error("embedded jarvis page carries a flat #050913 body backdrop; " +
			"the dark base may only exist inside the circular #jarvis-disc")
	}
	// The body rule stays transparent by construction (the earlier rework's
	// served-page tooth depends on this exact form).
	if !strings.Contains(page, "background: transparent;") {
		t.Error("embedded jarvis page body rule is no longer transparent by construction")
	}
	// The disc keeps its 96% window-edge diameter - the geometry the live disc
	// proof's ring and corner probes are calibrated against.
	if !regexp.MustCompile(`(?s)#jarvis-disc\s*\{[^}]*width:\s*96%`).MatchString(page) {
		t.Errorf("#jarvis-disc rule lost its %d%% window-edge diameter", jarvisDiscDiameterPct)
	}
	// Document order: the disc element BEFORE #root - the disc is the backdrop
	// layer, the React particle canvas paints above it.
	discIdx := strings.Index(page, `id="jarvis-disc"`)
	rootIdx := strings.Index(page, `id="root"`)
	if discIdx < 0 || rootIdx < 0 || discIdx > rootIdx {
		t.Error("#jarvis-disc must appear in the document BEFORE #root: " +
			"the disc is the backdrop layer, the particle canvas paints above it")
	}
}

// TestJarvisServedDocumentKeepsDiscAndBelt pins the SERVED document (the
// injection applied): the transparent belt both orb pages share must survive
// next to the disc, because the belt governs the pixels OUTSIDE the disc
// (html,body) while the disc paints only inside its circle.
func TestJarvisServedDocumentKeepsDiscAndBelt(t *testing.T) {
	doc := string(orbJarvisDocument())
	const injectedBelt = "<style>html,body{background:transparent !important}</style>"
	if !strings.Contains(doc, injectedBelt) {
		t.Error("served jarvis document lost the injected transparent belt")
	}
	if !strings.Contains(doc, `id="jarvis-disc"`) {
		t.Error("served jarvis document lost the dark disc element")
	}
}
