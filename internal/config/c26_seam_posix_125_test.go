//go:build !windows

// Ticket 125 AC#1: the POSIX half of the one claim about the sealing seam that
// no test on this platform made.
//
// Where the claim lives today, measured at the commit this file was written in
// (81b4d5f): "C26 is installed in winsec" has exactly one assertion, in
// internal/winsec/resolve_windows_test.go behind //go:build windows. The only
// non-Windows reader of winsec.PathResolverInstalled() is
// TestAC4POSIXFloorAnswersInsideTheNamedTree, which calls t.Skipf the moment the
// seam holds something - so on POSIX nothing pinned "the pipeline is in place",
// and nothing pinned the failure this ticket is about either: acceptor-ticket119
// measured the seam falling back to <nil> (whole C26 layer absent, every seal on
// the built-in floor) in a container whose TMPDIR was spelled through a symlink,
// and the ERROR line that says so had no test, no listener and no CI step.
//
// Why this nail is hosted in internal/config and not in internal/winsec, which
// is where it would naturally live: the nail needs a process that LINKS
// internal/risk, because the install is that package's init() reaching this
// package's seam (winsec cannot import risk - risk -> observe -> secret -> winsec
// is already a path). internal/winsec's own POSIX test binary does not link risk
// (measured: go list -deps -test ./internal/winsec/ names no internal/risk), and
// adding the import there was measured, not guessed: `go test -count=1 -v
// ./internal/winsec/` in a Linux container goes 26 PASS/0 SKIP -> 25 PASS/1 SKIP,
// and the leg that turns into a skip is TestAC4POSIXFloorAnswersInsideTheNamedTree
// - tickets 103/108/113's POSIX floor leg. Pinning the wiring there would buy
// one running leg by silencing another, so it is pinned in the nearest process
// that already has the edge: internal/config imports internal/risk in
// manager.go/validate.go, so this binary links the pipeline with no new edge,
// and the winsec floor legs keep their subject.
//
// Both halves below ask for a fact only the real pipeline can satisfy, so a
// stand-in does not pass: the built-in floor can refuse and pass through, never
// rewrite (see internal/winsec/resolve.go's note on why it is not a second
// normalizer), so an answer that folds a ".." shape into a clean spelling is
// evidence of C26 specifically.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// TestAC1POSIXSeamHoldsC26Pipeline125 is the nail: on POSIX, in a binary that
// links internal/risk, winsec's sealing seam has to hold the C26 pipeline. It
// deliberately does not skip on any condition, because "POSIX has no red-capable
// assertion here" is the defect being closed.
func TestAC1POSIXSeamHoldsC26Pipeline125(t *testing.T) {
	t.Logf("fixture: TMPDIR=%q resolves=%q", os.TempDir(), evalOrSelf(t, os.TempDir()))
	r := winsec.PathResolverInstalled()
	if r == nil {
		t.Fatalf("AC#1 RED: winsec.PathResolverInstalled() = <nil>: this process links internal/risk, " +
			"whose init() is supposed to hand C26 to the sealing seam, so either the install never ran or the " +
			"seam guard refused it and every seal in this process is on the built-in floor")
	}
	if got := fmt.Sprintf("%T", r); got != "risk.c26Pipeline" {
		t.Fatalf("AC#1 RED: the seam holds %s, not risk.c26Pipeline: the wiring under test is internal/risk's "+
			"init installing its own pipeline, and a stand-in answers a different question", got)
	}

	// The behavioral half: what only C26 can do on this platform. The floor
	// refuses a spelling that traverses ".."; the pipeline answers it with the
	// tree's own name. If this call refuses, the installed resolver is not
	// rewriting anything and the seam is floor-equivalent whatever its type says.
	base := t.TempDir()
	sep := string(filepath.Separator)
	hostile := base + sep + ".." + sep + "wisp-125-probe"
	got, err := r.Resolve(hostile)
	t.Logf("reading: installed %T answered %q with %q", r, hostile, got)
	if err != nil {
		t.Fatalf("AC#1 RED: the installed %T refused its own probe-shaped input %q: %v - the built-in floor "+
			"refuses exactly this, so nothing here distinguishes C26 from no C26", r, hostile, err)
	}
	if got == hostile {
		t.Fatalf("AC#1 RED: the installed %T passed the hostile shape %q through unchanged", r, hostile)
	}
	if strings.Contains(got, sep+".."+sep) || strings.HasSuffix(got, sep+"..") {
		t.Fatalf("AC#1 RED: the installed %T answered %q with %q, which still traverses a parent pointer: an "+
			"answer of that shape is not the tree the caller named", r, hostile, got)
	}
}

// evalOrSelf reports the resolved spelling of a path, or the path itself when it
// cannot be resolved, so the log line is a reading and not a t.Fatal.
func evalOrSelf(t *testing.T, path string) string {
	t.Helper()
	real, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Sprintf("<unresolved: %v>", err)
	}
	return real
}
