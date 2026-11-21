//go:build windows

package main

// Ticket 33 AC#10's geometry seam: the place where [panel] width, height and
// scale stop being three config values and become the numbers Win32 actually
// receives.
//
// WHY A SEAM. AC#10's disease is a producer that exists only in a test: the
// value is read, nothing is set, and the reading cannot tell the two apart.
// The cure is a seam the production path MUST cross - one pure mapping
// (panelWindowGeometry), one injection point (PanelManager.applyGeometry,
// production value posts ONE SetWindowPos, a test injects a recorder) - so a
// test can watch the real path deliver real numbers without opening a window.
//
// WHAT scale MEANS HERE, SAID PLAINLY. The schema calls [panel] scale "the UI
// zoom factor; 0 = follow system DPI". The WebView2 CONTENT zoom lives on
// ICoreWebView2Controller2.PutZoomFactor, which this repository's fork declares
// only as a raw vtbl slot with no Go wrapper (forks/go-webview2/pkg/edge/
// ICoreWebView2Controller2.go - the nail TestAC10ScaleHasNoProducerAndTheReason
// IsNamed pins that), and forks/** is not this leg's write surface. What IS
// reachable, and what the dispatch names as the sanctioned path, is the
// WINDOW-bounds lever: SetWindowPos. So a written scale multiplies the resolved
// bounds - the panel window is created/sized at the zoomed frame - and scale 0
// keeps exactly today's frame (follow the system DPI as the window already
// does). That is a real, observable producer through the SetWindowPos path; it
// is NOT the controller zoom, and nothing here may be quoted as one.

import (
	"math"

	"github.com/CarlosShao/wisp/internal/config"
)

// panelWindowGeometry is THE pure mapping from [panel] to the geometry triple
// the window is created and resized from. The fallback rules are the ones
// windowOptions has carried since 票 255 AC#4, now single-sourced here:
//
//   - width <= 0 answers panelWidthPx and height <= 0 answers panelHeightPx -
//     the host's own constants, "the value when config is absent". A NEGATIVE
//     answer is garbage, not a second spelling of auto, and falls back the
//     same way rather than inventing a rule for itself.
//   - height 0 is NOT auto-height: the standing 255 ruling answers the
//     constant with "由内容定高" registered as NOT DONE; this function keeps
//     that answer and invents no content metric.
//   - scale <= 0 answers 0 - the schema's own "follow system DPI", which is
//     today's no-zoom frame; only a WRITTEN positive factor moves anything.
//
// A triple this function answers is FINAL: the create consumes w/h verbatim
// (the outer frame), the re-show folds the pair through panelScaledSize, and
// the applier hop multiplies through the same helper - one mapping, no second
// spelling anywhere downstream.
func panelWindowGeometry(cfg config.PanelSection) (w, h int, scale float64) {
	w, h = panelWidthPx, panelHeightPx
	if cfg.Width > 0 {
		w = cfg.Width
	}
	if cfg.Height > 0 {
		h = cfg.Height
	}
	if cfg.Scale > 0 {
		scale = cfg.Scale
	}
	return w, h, scale
}

// panelScaledSize applies the written zoom factor to a bounds pair: the
// window's frame is sized at round(w*scale) x round(h*scale). scale <= 0 is
// "follow system DPI" and answers the pair unchanged - the default configs'
// frames stay byte-identical to what they were before AC#10.
func panelScaledSize(w, h int, scale float64) (int, int) {
	if scale <= 0 {
		return w, h
	}
	return int(math.Round(float64(w) * scale)), int(math.Round(float64(h) * scale))
}

// SetWindowPos flags this host's scale hop needs (winuser.h). NOMOVE keeps the
// user's placement, NOZORDER and NOACTIVATE keep the raise/focus behaviour the
// show path already owns - the hop changes SIZE and nothing else.
const (
	pnlSwpNoMove     = 0x0002
	pnlSwpNoZOrder   = 0x0004
	pnlSwpNoActivate = 0x0010
)

// panelResolved is the triple windowOptions last resolved. It is kept on the
// manager so the create tail's applier hop and the re-show's scale fold
// consume the SAME read: one read per create and one per re-show, never two
// versions of one file deciding one window's geometry (票 255's own note).
type panelResolved struct {
	w, h  int
	scale float64
	ok    bool
}

// applyResolvedGeometry is the create tail's hop (bringUp, right after the
// control is registered). It hands the triple windowOptions just resolved to
// the seam - the production value posts ONE SetWindowPos onto the live window,
// a test's recorder watches the delivery. Nothing here opens, destroys or
// reorders a window: the create made the frame, this hop may only resize it.
func (m *PanelManager) applyResolvedGeometry() {
	m.mu.Lock()
	r := m.resolved
	apply := m.applyGeometry
	m.mu.Unlock()
	if !r.ok || apply == nil {
		return
	}
	apply(r.w, r.h, r.scale)
}

// applyGeometryWin32 is applyGeometry's production value. A written scale
// reaches Win32 HERE - ONE SetWindowPos onto the frame the create made, sized
// at the zoomed bounds. scale <= 0 answers nothing: the create's frame IS the
// follow-DPI frame, and a second call that re-states it would be a syscall
// pretending to be a producer. No window (or a destroyed one) also answers
// nothing - the next show resolves and applies fresh.
func (m *PanelManager) applyGeometryWin32(w, h int, scale float64) {
	if scale <= 0 {
		return
	}
	hwnd := m.windowHandle()
	if hwnd == 0 {
		return
	}
	fw, fh := panelScaledSize(w, h, scale)
	_, _, _ = pnlSetWindowPos.Call(hwnd, 0, 0, 0, uintptr(fw), uintptr(fh),
		uintptr(pnlSwpNoMove|pnlSwpNoZOrder|pnlSwpNoActivate))
}
