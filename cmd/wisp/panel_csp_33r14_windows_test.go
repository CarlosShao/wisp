//go:build windows

package main

// Ticket 33 AC#6's rulers (33-r14): the CSP response header on every document
// the panel host serves. The policy itself is SPEC-06 sec 9's strict CSP
// (docs/specs/SPEC-06-security-gatekeeping.md:126) verbatim, and the serving
// surface is the AddWebResourceRequestedFilter arm in panel_control_windows.go
// - the one place this tree answers a resource request, so a policy on it
// covers the entry document, every hashed asset and the not-found body.
// Everything here is window-free on purpose: the header string is built by
// panelResponseHeaders and handed to CreateWebResourceResponse, so these
// rulers ask the filter's own construction, not a real browser.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// wantCSPPolicy is the ONE spelling the AC#6 value may carry. It is compared
// against panelCSPHeaderValue below; a drift in either the constant or the
// spec line fails here, which is the two-way reconciliation the ticket's CSP
// row asks for ("provide the injection hook + default strict CSP now ...
// SPEC-06 sec 9").
const wantCSPPolicy = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' data: blob:; font-src 'self'; connect-src 'none'; frame-src 'none'; object-src 'none'; base-uri 'none'"

func TestAC6CSPHeaderIsOnEveryServedResponse33r14(t *testing.T) {
	if panelCSPHeaderValue != wantCSPPolicy {
		t.Errorf("panelCSPHeaderValue drifted from SPEC-06 sec 9's strict CSP:\n got %q\nwant %q",
			panelCSPHeaderValue, wantCSPPolicy)
	}
	ok200 := panelResponseHeaders(200, "text/html; charset=utf-8")
	if !strings.Contains(ok200, "Content-Security-Policy: "+wantCSPPolicy) {
		t.Errorf("the 200 header set carries no CSP line: %q", ok200)
	}
	if !strings.Contains(ok200, "Content-Type: text/html; charset=utf-8") {
		t.Errorf("adding the CSP line dropped or mangled the Content-Type the filter already served: %q", ok200)
	}
	// Exactly two header fields, CRLF-separated: the string this host hands to
	// CreateWebResourceResponse is one header set per response, no more.
	if got := strings.Split(ok200, "\r\n"); len(got) != 2 {
		t.Errorf("the 200 header set carries %d CRLF-separated field(s), want exactly 2 (Content-Type + CSP): %q", len(got), got)
	}
	// The miss answer is empty-bodied, not header-less: the not-found body the
	// filter still serves lives under the same policy as the documents around
	// it, and it still names no content type it does not have.
	miss := panelResponseHeaders(404, "")
	if !strings.Contains(miss, "Content-Security-Policy: "+wantCSPPolicy) {
		t.Errorf("the 404 header set carries no CSP line: %q", miss)
	}
	if strings.Contains(miss, "Content-Type:") {
		t.Errorf("the 404 header set invents a Content-Type for a body it does not have: %q", miss)
	}
}

// TestAC6CSPHeaderHasOneBuildSite33r14 keeps the single construction point:
// the literal "Content-Security-Policy" may be spelled exactly once in
// panel_control_windows.go - in panelCSPHeader's own definition - and every
// response gets the policy through that constant. A second place spelling the
// header by hand is a second place for the served document and the tested one
// to drift apart. Comments do not count (the file's prose names the header):
// the walk is over string literals only, the same comment-exempt reading
// d22scan's emoji ban takes.
func TestAC6CSPHeaderHasOneBuildSite33r14(t *testing.T) {
	src, err := os.ReadFile("panel_control_windows.go")
	if err != nil {
		t.Fatalf("read panel_control_windows.go: %v", err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "panel_control_windows.go", src, 0)
	if err != nil {
		t.Fatalf("parse panel_control_windows.go: %v", err)
	}
	hits := 0
	ast.Inspect(file, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING &&
			strings.Contains(lit.Value, "Content-Security-Policy") {
			hits++
		}
		return true
	})
	if hits != 1 {
		t.Errorf("panel_control_windows.go spells the Content-Security-Policy literal %d time(s), want exactly 1 (panelCSPHeader's definition): every extra one is a second build site", hits)
	}
}
