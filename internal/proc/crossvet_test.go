package proc

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// crossVetPackages are the main-module packages this check compiles for
// GOOS=linux. Two rules decide membership, and both come from ticket 78:
//
//  1. the package must be one whose untagged files can legitimately be
//     referenced from a non-Windows build (internal/proc and internal/ball
//     both had untagged files reaching into //go:build windows declarations,
//     which is exactly the shape this catches); and
//  2. the package must be cgo-free. A cross-compile from windows forces
//     CGO_ENABLED=0, so a package that needs cgo fails to LOAD (its whole
//     dependency resolves to "build constraints exclude all Go files") long
//     before it fails to type-check, and the real error never surfaces.
//     cmd/wisp is excluded for that reason - it imports sherpa-onnx and can
//     only be cross-checked by the `go vet ./...` that the lint job already
//     runs on ubuntu-latest, where cgo is on.
//
// The list is therefore narrow and every entry is verified green; widening it
// is welcome, but only with a package that is actually clean under this exact
// command, or the check becomes the thing that is always red.
var crossVetPackages = []string{"./internal/proc/", "./internal/ball/"}

// TestCrossVetForLinux catches the ticket 78 regression class on the runners
// where it is otherwise invisible.
//
// The defect: statevisual.go (untagged) called mulA, which was declared only
// in renderer_windows.go (//go:build windows), and cmd/wisp's untagged files
// reached for proc.Boot / proc.Runtime the same way. Both compile and vet
// cleanly on a Windows developer machine and on the Windows CI jobs, because
// those build for windows. They fail only for another GOOS - and until this
// test, nothing in the repo ever asked for one from a Windows runner.
//
// This test is that something: it type-checks the packages above as the linux
// compiler sees them, so an untagged file referencing a Windows-only symbol
// goes red before push, on every host. It is also what CI's ubuntu lint job
// does after the fact, which is how the defect was finally caught - but there
// it aborted the whole lint job at `go vet`, downstream of the D22 seven-ban
// scan, so the scan produced no verdict for as long as these errors stood.
//
// There is no skip path. `go` not being resolvable is a loud failure, not a
// silent pass (a permanently-skipping gate is the false green this repo has
// been burned by twice today).
func TestCrossVetForLinux(t *testing.T) {
	root := repoRootForTest(t)

	goExe := "go"
	if p, err := exec.LookPath("go"); err == nil {
		goExe = p
	} else {
		rooted := filepath.Join(runtime.GOROOT(), "bin", "go.exe")
		if st, err := os.Stat(rooted); err != nil || st.IsDir() {
			t.Fatalf("cannot resolve the go tool: LookPath failed (%v) and %s is not a file - "+
				"this check must fail loudly, never skip", err, rooted)
		}
		goExe = rooted
	}

	cmd := exec.Command(goExe, "vet")
	cmd.Args = append(cmd.Args, crossVetPackages...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return
	}
	t.Fatalf("%s vet %s for GOOS=linux failed (%v):\n%s\n\n"+
		"An untagged file in one of these packages references a symbol that exists only in a\n"+
		"//go:build windows file. Fix the shape, do not add a zero-value non-windows stub:\n"+
		"  - the symbol is genuinely platform-specific -> build-tag the referencing file, and\n"+
		"    give any caller that must stay untagged a fail-closed !_windows counterpart;\n"+
		"  - the symbol is plain portable code that merely lives in the wrong file -> move it\n"+
		"    to an untagged file (that is what ticket 78 did for ball.mulA).",
		goExe, strings.Join(crossVetPackages, " "), err, out)
}

// repoRootForTest walks up from this file to the module root, so the check can
// address packages by their ./internal/... paths from whichever directory the
// test binary happens to run in.
func repoRootForTest(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed; cannot locate the module root")
	}
	root := filepath.Join(filepath.Dir(thisFile), "..", "..")
	if st, err := os.Stat(filepath.Join(root, "go.mod")); err != nil || st.IsDir() {
		t.Fatalf("no go.mod at %s: this test must run from inside the module", root)
	}
	return root
}
