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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	_ "embed"
)

//go:embed assets/orb-siri.html
var orbPageHTML []byte

// jarvisPageHTML is the SECOND style family's page (ticket 65, owner-ordered
// second visual style): the cyber1443/jarvis-ai-orb-web-animation component
// (upstream MIT) built by tools/jarvis-export into ONE self-contained HTML -
// react + react-dom + three + the component, inlined by inline.mjs. The
// component source is unmodified; entry.tsx (the only integration-owned code)
// mounts it and exposes window.jarvisOrb = { setState, getState } with the
// page's own four state names (idle|thinking|success|alert). Built artifact,
// build method and the three build-time pitfalls are recorded in
// docs/evidence/s1/65-jarvis-style-integration.md; the sha pin lives in
// orb_test.go (jarvisUpstreamSHA256), same discipline as orb-siri.html.
//
// Unlike the siri page this one is WebGL (not WebGPU), so it would render
// even off the NavigateToString fallback; it is served through the SAME
// virtual-host seam anyway so both styles share one no-socket shape.
//
// Audio: the upstream component has NO audio-band input (verified
// states.ts/renderer.ts/JarvisOrb.tsx); none is invented - pushAudio skips
// this style (orb_windows.go).
//
//go:embed assets/jarvis-orb.html
var jarvisPageHTML []byte

const (
	// orbSiteHost is the synthetic host name the site folder is mapped onto
	// (SetVirtualHostNameToFolderMapping takes the host WITHOUT the scheme).
	// ".orb" is a TLD no DNS resolver answers, so even a request that escaped
	// the mapping fails OFFLINE instead of reaching a real host - the same
	// no-socket shape the panel host uses
	// (cmd/wisp/panel_control_windows.go panelOrigin).
	orbSiteHost = "wisp.orb"

	// orbOrigin / orbEntryName / orbEntryURL: the SIRI page's origin under the
	// mapping and its entry path. The jarvis page lives in the same folder
	// under its own entry name (jarvisEntryName), so a runtime style switch is
	// a pure navigation - no site rewrite, no browser rebuild.
	orbOrigin    = "https://" + orbSiteHost
	orbEntryName = "orb.html"

	orbEntryURL = orbOrigin + "/" + orbEntryName

	// jarvisEntryName / jarvisEntryURL: the jarvis page's entry under the same
	// mapped folder (orbWriteSite materializes BOTH entries every attach).
	jarvisEntryName = "jarvis.html"

	jarvisEntryURL = orbOrigin + "/" + jarvisEntryName
)

// orbCSPHeaderValue is the policy the SERVED document carries as an injected
// <meta http-equiv="Content-Security-Policy"> (orbPageDocument below). The
// first cut of this leg carried the same policy as a response HEADER on a
// WebResourceRequested-constructed response; that seam measured black (the
// jchv CreateWebResourceResponse wrapper misreports success off the raw
// GetLastError, orb_windows.go header), so delivery moved to a virtual-host
// mapping - and neither that nor NavigateToString offers a response-header
// channel. The meta form is the CSP the HTML standard itself defines for
// this case: enforced at parse time, and a document API can TIGHTEN it but
// never LOOSEN it, so the page cannot strip the exfiltration cut at runtime.
// The export page is one self-contained document: one inline module script,
// one inline style block, a data: favicon, and WebGPU - which needs no CSP
// grant. No network verb of any kind: default-src 'none' with connect-src
// 'none' is the exfiltration cut, same shape as the panel's serving policy.
const orbCSPHeaderValue = "default-src 'none'; script-src 'unsafe-inline'; " +
	"style-src 'unsafe-inline'; img-src data:; connect-src 'none'; " +
	"frame-src 'none'; object-src 'none'; base-uri 'none'"

// orbPageDocument returns the SIRI document actually served: the upstream page
// byte-for-byte PLUS the two injected head entries (orbServeDocument below).
func orbPageDocument() []byte { return orbServeDocument(orbPageHTML) }

// orbJarvisDocument returns the JARVIS document actually served: the same
// injection applied to the built single-file page (jarvisPageHTML). The jarvis
// asset ships with a transparent body background of its own (index.html
// template); the injected transparent style stays as the same belt both pages
// get, and the CSP meta is what makes the page offline by policy.
func orbJarvisDocument() []byte { return orbServeDocument(jarvisPageHTML) }

// orbServeDocument injects the two head entries into one embedded page right
// after <head> - the only transformations the serving seam applies on top of
// the embed, applied in memory so the asset files stay untouched upstream (or
// untouched build output, for jarvis; their shas are pinned in orb_test.go).
// The first is the CSP <meta> (a response-header CSP is not available to
// either delivery shape). The second is the transparent background <style>:
// the script bootstrap (orbPageBootstrap) clears the body fill through
// AddScriptToExecuteOnDocumentCreated, but that seam is a RACE - the jchv
// wrapper passes a NULL completion handler and the registration can land
// after the first document already loaded (measured 2026-11-07 on real
// desktop runs: page alive, WebGPU running, corners showing the upstream dark
// body fill for the whole process lifetime). A declarative CSS override has
// no such seam: `transparent !important` beats any opaque body rule
// regardless of rule order, and the pixels outside the orb become the
// window's transparency. If a page's head marker were ever missing the
// entries are injected after the doctype instead: the policy is never
// silently dropped.
func orbServeDocument(page []byte) []byte {
	meta := []byte("<head>\n  <meta http-equiv=\"Content-Security-Policy\" content=\"" +
		orbCSPHeaderValue + "\" />" +
		"\n  <style>html,body{background:transparent !important}</style>")
	if doc := bytes.Replace(page, []byte("<head>"), meta, 1); len(doc) != len(page) {
		return doc
	}
	// Fallback (head tag removed): inject after the doctype rather than serve
	// an unpoliced document.
	return bytes.Replace(page, []byte(">"), []byte(">\n<head>"+orbCSPHeaderValue+"</head>"), 1)
}

// orbWriteSite materializes BOTH style pages into dir as the mapping's entry
// files: the same bytes the NavigateToString fallback would hand (orb.html =
// orbPageDocument, jarvis.html = orbJarvisDocument), so both delivery shapes
// serve ONE document shape per style, and a runtime style switch is a pure
// navigation between the two entries. The folder lives under the orb's
// per-user user-data path (orb_windows.go orbDataPath); writing it is
// idempotent, done once per attach.
func orbWriteSite(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("ball: orb site mkdir: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, orbEntryName), orbPageDocument(), 0o600); err != nil {
		return fmt.Errorf("ball: orb site write: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, jarvisEntryName), orbJarvisDocument(), 0o600); err != nil {
		return fmt.Errorf("ball: orb site write: %w", err)
	}
	return nil
}

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

// orbJarvisStateScript is the jarvis page's twin of orbStateScript: the entry
// (tools/jarvis-export/entry.tsx) exposes window.jarvisOrb = { setState,
// getState } and IGNORES unknown names instead of throwing, so a stray push
// can never kill the RAF loop; the mapping table (orbstate.go jarvisStateFor)
// is pinned to the page's own four-name vocabulary by tests all the same.
func orbJarvisStateScript(name string) string {
	return fmt.Sprintf("window.jarvisOrb&&window.jarvisOrb.setState(%q)", name)
}

// orbPageForStyle resolves one style family to the page it serves: the siri
// glass page for StyleGlass, the jarvis particle page for StyleJarvis. This
// is the whole "which asset" decision - newOrbWindow navigates to entryURL
// (or falls back to doc), orbWindow.setStyle re-navigates on a runtime switch.
func orbPageForStyle(s Style) (doc []byte, entryName string, entryURL string) {
	if s == StyleJarvis {
		return orbJarvisDocument(), jarvisEntryName, jarvisEntryURL
	}
	return orbPageDocument(), orbEntryName, orbEntryURL
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
