//go:build !windows

// Ticket 125 AC#1's nail, plus AC#2's reading of the same nail in a second shape:
// the one POSIX claim about the sealing seam that no test on this platform made.
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
	"os/exec"
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

// childShape125 is how this file's subprocess passes which temp spelling it is
// supposed to be standing in. The install decision is taken by internal/risk's
// init(), i.e. once per process and before any test body runs, so the only honest
// way to read it for a second shape is a second process.
const childShape125 = "WISP_125_SEAM_SHAPE"

// TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125 runs AC#1's nail again in
// a child process whose TMPDIR reaches itself through a symlink, which is the
// shape acceptor-ticket119 measured as winsec.PathResolverInstalled()=<nil>: the
// whole C26 layer missing, every seal in the process on the built-in floor, and
// the only witness an ERROR line nobody on POSIX had a test for.
//
// The child runs TestAC1POSIXSeamHoldsC26Pipeline125 unchanged, so the reading
// this takes is the same assertion the plain shape already passes - one nail, two
// shapes, no second copy of the claim. A plain-TMPDIR child is run first as the
// instrument's own control: if it fails, the link-shaped child's failure says
// nothing about the shape.
func TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125(t *testing.T) {
	base := t.TempDir()
	if link := firstLink125(base); link != "" {
		t.Skipf("fixture: the harness's temp base %q reaches itself through the link at %q, so neither child below "+
			"is the control shape and the pair would measure one thing twice", base, link)
	}
	realRoot := filepath.Join(base, "real125", "tmproot125")
	if err := os.MkdirAll(realRoot, 0o700); err != nil {
		t.Fatalf("fixture mkdir %s: %v", realRoot, err)
	}
	linkRoot := filepath.Join(base, "varlink125")
	if err := os.Symlink(filepath.Join(base, "real125"), linkRoot); err != nil {
		t.Skipf("cannot build the symlink this case measures (%v): without it the shape does not exist", err)
	}
	info, err := os.Lstat(linkRoot)
	if err != nil {
		t.Fatalf("fixture Lstat %s: %v", linkRoot, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("fixture: %q is not a symlink (mode %v) - every reading taken through it would be a false green", linkRoot, info.Mode())
	}
	linked := filepath.Join(linkRoot, "tmproot125")
	if firstLink125(linked) == "" {
		t.Fatalf("fixture: the planted TMPDIR %q has no link in its ancestor chain, so it is the control shape, not the measured one", linked)
	}
	if clean, err := filepath.EvalSymlinks(linked); err != nil || clean != realRoot {
		t.Fatalf("fixture: %q resolves to %q (err %v), want %q", linked, clean, err, realRoot)
	}

	for _, tc := range []struct{ name, tmpdir string }{
		{"control_plain_temp", realRoot},
		{"measured_symlink_spelled_temp", linked},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, rc := runSeamChild125(t, tc.tmpdir)
			t.Logf("reading: child TMPDIR=%q rc=%d", tc.tmpdir, rc)
			for _, line := range splitLines125(out) {
				if strings.Contains(line, "resolver installed") || strings.Contains(line, "refusing to install") ||
					strings.Contains(line, "AC#1 RED") || strings.Contains(line, "--- ") {
					t.Log("child: " + line)
				}
			}
			if rc != 0 {
				t.Errorf("AC#2 RED: in a process whose temp dir is spelled %q (it reaches %q through the link at %q) "+
					"the sealing seam does not hold C26 - rc=%d. The claim under test is winsec.PathResolverInstalled()"+
					"'s type, so this is the whole pipeline missing, not a rewrite being unavailable.",
					tc.tmpdir, realRoot, linkRoot, rc)
			}
		})
	}
}

// runSeamChild125 re-executes this test binary with TMPDIR planted and only AC#1's
// nail selected, returning its combined output and exit code.
func runSeamChild125(t *testing.T, tmpdir string) (string, int) {
	t.Helper()
	env := append(os.Environ(), childShape125+"=1", "TMPDIR="+tmpdir)
	cmd := exec.Command(os.Args[0],
		"-test.run=^TestAC1POSIXSeamHoldsC26Pipeline125$",
		"-test.count=1", "-test.v")
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return string(out), ee.ExitCode()
		}
		t.Fatalf("cannot re-exec this test binary (%v): the install decision is taken at init, so no in-process "+
			"call can measure a second temp shape without it: %s", err, out)
	}
	return string(out), 0
}

// firstLink125 names the first prefix of a spelling the filesystem says is a link,
// or "". It asks the filesystem, not the function under test.
func firstLink125(path string) string {
	cur := path
	for {
		if info, err := os.Lstat(cur); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return ""
		}
		cur = parent
	}
}

// splitLines125 splits a child's output on either newline convention.
func splitLines125(out string) []string {
	out = strings.ReplaceAll(out, "\r\n", "\n")
	return strings.Split(out, "\n")
}
