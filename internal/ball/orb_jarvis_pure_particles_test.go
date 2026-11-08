package ball

// Ticket 65 pure-particle leg (2026-11-08): the owner's ruling VERBATIM and
// highest-priority - the jarvis orb must have NO sphere body, NO visible
// frame, NO backdrop: "直接就是粒子态" (particles only, like Jarvis). This leg
// REVERSES the previous dark-disc ruling: the in-page dark circular backdrop
// the disc leg shipped is GONE from the template and the asset (the owner
// rejected the disc itself; the superseded form lives in git history only).
//
// The visibility fix is the blending mode, applied at build time by the
// scaffold's wisp-jarvis-normal-blending vite plugin
// (tools/jarvis-export/vite.config.mjs): the upstream's AdditiveBlending ADDS
// particle color onto the backdrop, so over white surfaces low-alpha particle
// edges clamp toward white (the owner's JARVIS-PRODUCT-REAL.png wash-out);
// NormalBlending composites each particle by its OWN alpha over whatever sits
// behind it - cyan on white, cyan on dark - with no backdrop needed. That is
// the mainstream three.js shape for a particle system on an unknown
// background: normal blending + the component's own saturated cyan palette,
// NOT additive + a dark disc.
//
// These tests pin the pure-particle shape in the EMBED - the exact bytes the
// WebView2 loads; the live composited-pixel proof (white floor AND black
// floor, both must show cyan particles, neither may show a disc) rides
// winlive (orb_jarvis_pure_particles_live_windows_test.go).

import (
	"strings"
	"testing"
)

// TestJarvisAssetShipsPureParticles pins the NO-DISC shape in the embedded
// asset (the flip of the previous leg's disc tooth): no backdrop element, no
// disc CSS, none of the disc's dark-blue family anywhere in the document.
func TestJarvisAssetShipsPureParticles(t *testing.T) {
	page := string(jarvisPageHTML)
	if strings.Contains(page, "jarvis-disc") {
		t.Fatal("embedded jarvis page still carries a disc marker; " +
			"the owner rejected the in-page backdrop - particles only")
	}
	if strings.Contains(page, "#050913") {
		t.Fatal("embedded jarvis page still carries the disc's dark-blue family (#050913)")
	}
	if strings.Contains(page, "radial-gradient(circle at 50% 42%") {
		t.Fatal("embedded jarvis page still carries the disc's radial-gradient backdrop rule")
	}
	// The body rule stays transparent by construction (the earlier rework's
	// served-page tooth depends on this exact form): no painted backdrop of
	// ANY kind - the pixels outside the particles are the window's
	// transparency.
	if !strings.Contains(page, "background: transparent;") {
		t.Error("embedded jarvis page body rule is no longer transparent by construction")
	}
	// The React mount point is still the ONLY body child: with the disc gone
	// there is no backdrop element before it (a restored backdrop would add a
	// second <div between <body> and #root).
	bodyIdx := strings.Index(page, "<body>")
	rootIdx := strings.Index(page, `id="root"`)
	if bodyIdx < 0 || rootIdx < 0 || bodyIdx > rootIdx {
		t.Fatal("embedded jarvis page lost the body/#root document shape")
	}
	if n := strings.Count(page[bodyIdx:rootIdx], "<div"); n != 1 {
		t.Errorf("embedded jarvis page has %d divs before #root (want exactly 1, the root mount): "+
			"a backdrop element is back in front of the particles", n)
	}
}

// TestJarvisServedDocumentKeepsBeltNoDisc pins the SERVED document (the
// injection applied): the transparent belt both orb pages share must survive,
// and the disc must stay gone from the exact bytes the serving seam
// (orbJarvisDocument) hands the WebView2.
func TestJarvisServedDocumentKeepsBeltNoDisc(t *testing.T) {
	doc := string(orbJarvisDocument())
	const injectedBelt = "<style>html,body{background:transparent !important}</style>"
	if !strings.Contains(doc, injectedBelt) {
		t.Error("served jarvis document lost the injected transparent belt")
	}
	if strings.Contains(doc, "jarvis-disc") {
		t.Error("served jarvis document still carries a disc marker; " +
			"the owner rejected the in-page backdrop - particles only")
	}
}
