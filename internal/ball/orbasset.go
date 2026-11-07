package ball

// Ticket-65 rework, integration leg (A1594): the glass-quality problem is
// solved by SHIPPING the open-source liquid glass orb, not by re-authoring a
// shader. The renderer page is LerSent001/orb's own single-file export
// (orb-siri.html, upstream MIT, package.json "license": "MIT"), embedded
// byte-for-byte below (the only transformation is CRLF -> LF). The Go side of
// this file is the serving seam and the two pure script builders the WebView2
// host (orb_windows.go) pushes through ExecuteScript.
//
// Upstream interface, verified from the page itself (the ONLY seam this
// integration talks to):
//
//	window.liquidOrb = { getState, setState, setAudioBands }  (frozen object)
//	setState(name) - names are exactly the stateSeeds keys: "idle", "thinking"
//	setAudioBands({low, mid, high, all}) - normalised 0..1 bands, "zero them
//	  on stop" per the page's own contract
//
// The page throws on an unknown state name and the shader paints exactly zero
// alpha outside the orb (premultiplied canvas), which is what makes the
// host-side transparency and the mapping table in orbstate.go the entire
// integration surface.

import (
	"fmt"
	"strconv"

	_ "embed"
)

//go:embed assets/orb-siri.html
var orbPageHTML []byte

const (
	// orbOrigin is the synthetic origin the page is served under, spelled in
	// ONE place. ".orb" is a name no DNS resolver can answer, so a request
	// that ever escaped the WebResourceRequested filter fails OFFLINE instead
	// of reaching a real host - the same no-socket shape the panel host uses
	// (cmd/wisp/panel_control_windows.go panelOrigin).
	orbOrigin = "https://wisp.orb"

	// orbEntryName is the one path this embed serves; everything else 404s.
	orbEntryName = "orb.html"

	orbEntryURL           = orbOrigin + "/" + orbEntryName
	orbAssetFilterPattern = orbOrigin + "/*"
)

// orbCSPHeaderValue is the policy every served response carries as a response
// HEADER (never an injected <meta> the page could strip). The export page is
// one self-contained document: one inline module script, one inline style
// block, a data: favicon, and WebGPU - which needs no CSP grant. No network
// verb of any kind: default-src 'none' with connect-src 'none' is the
// exfiltration cut, same shape as the panel's serving policy.
const orbCSPHeaderValue = "default-src 'none'; script-src 'unsafe-inline'; " +
	"style-src 'unsafe-inline'; img-src data:; connect-src 'none'; " +
	"frame-src 'none'; object-src 'none'; base-uri 'none'"

// orbCSPHeader is the header LINE handed to CreateWebResourceResponse. It is
// on every response this filter answers - 200 and 404 alike.
const orbCSPHeader = "Content-Security-Policy: " + orbCSPHeaderValue

// orbPageBootstrap is injected with AddScriptToExecuteOnDocumentCreated, so
// it runs before the page's own module script. The upstream page ships an
// opaque dark body background because its export target is a browser tab;
// the Wisp host needs the pixels OUTSIDE the orb fully transparent - which
// the shader already provides (alpha exactly 0 there) - so the host clears
// the body background WITHOUT touching the page's shader, seeds, or
// interface. This bootstrap is the only byte the integration adds on top of
// the upstream export, and it lives in Go, not in the asset.
const orbPageBootstrap = `(function () {
  "use strict";
  function clear() {
    var b = document.body;
    if (b) { b.style.background = "transparent"; }
  }
  document.documentElement.style.background = "transparent";
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", clear);
  } else {
    clear();
  }
})();`

// orbStateScript builds the ExecuteScript body that pushes one mapped orb
// state. The mapping table in orbstate.go is pinned by a test that walks ALL
// 20 frozen ball states, because the page throws on an unknown name.
func orbStateScript(name string) string {
	return fmt.Sprintf("window.liquidOrb&&window.liquidOrb.setState(%q)", name)
}

// orbAudioScript builds the ExecuteScript body for one audio-band push. The
// page clamps every band into 0..1 itself, but a NaN slips through JS's
// Math.max/Math.min chain, so the Go side clamps FIRST and hands the page
// only finite literals.
func orbAudioScript(low, mid, high, all float64) string {
	return fmt.Sprintf(
		"window.liquidOrb&&window.liquidOrb.setAudioBands({low:%s,mid:%s,high:%s,all:%s})",
		orbBand(low), orbBand(mid), orbBand(high), orbBand(all),
	)
}

// orbBand clamps one band into [0,1] (NaN to 0) and renders it as a finite
// fixed-point JS number literal.
func orbBand(v float64) string {
	if v != v {
		v = 0
	}
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return strconv.FormatFloat(v, 'f', 4, 64)
}
