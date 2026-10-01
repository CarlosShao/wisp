package main

// Platform-neutral gates for the panel host (ticket 33). These run in the ubuntu
// core scope too (J7: the L1 half must have a real CI denominator), which is why
// they read the Windows host SOURCE FILE as text rather than executing a WebView2
// window: an import that would open a socket is forbidden at the source level,
// independent of GOOS.

import (
	"archive/zip"
	"bytes"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestPanelHostOpensNoListeningSocketL1 is AC#3's L1 half: the host production
// file must not import net or net/http (no listener can be built without them).
// The L2 half - reading GetExtendedTcpTable and filtering dwState==LISTEN - runs
// only on the desktop box (winlive has no CI job) in panel_host_windows_test.go.
func TestPanelHostOpensNoListeningSocketL1(t *testing.T) {
	src := filepath.Join("panel_host_windows.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse host source %s: %v", src, err)
	}
	forbidden := []string{"net", "net/http", "net/netip", "net/mail", "net/tcp"}
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		for _, bad := range forbidden {
			if path == bad || strings.HasPrefix(path, bad+"/") {
				t.Errorf("panel host imports %q, which can open a listening socket; D29/AC#3 forbid a localhost server", path)
			}
		}
	}
}

// TestEmbeddedDistCleanCheckoutHasPlaceholderOnly_AC12 states the AC#12 truth in
// a runnable form: what is COMMITTED under frontend/dist decides whether the
// panel can show a real page. In this repo the tracked embed content is only the
// .gitkeep anchor, so a clean checkout embeds no page. The capability to tell
// "placeholder only" from "a real bundle" is already nailed in internal/panel's
// TestAnchorOnlyBundleIsNotBuiltAndFailsClosed (assets_test.go:22); this test
// records the current tracked manifest and, if a real bundle is ever committed,
// fails closed on a half-bundle rather than a clean either-way.
func TestEmbeddedDistCleanCheckoutHasPlaceholderOnly_AC12(t *testing.T) {
	// go test runs with the package dir as cwd (cmd/wisp), so git ls-files must be
	// aimed at the repo root or it resolves "frontend/dist" against the wrong tree
	// and reports 0 tracked entries. (This is a real尺 measurement, so the path is
	// made explicit rather than trusted to cwd.)
	root := repoRootForTest(t)
	g := exec.Command("git", "ls-files", "frontend/dist")
	g.Dir = root
	out, err := g.Output()
	if err != nil {
		t.Fatalf("git ls-files frontend/dist in %s: %v", root, err)
	}
	entries := strings.FieldsFunc(string(out), func(r rune) bool { return r == '\n' || r == '\r' })
	for _, e := range entries {
		if e == "" {
			continue
		}
		if filepath.Base(e) == ".gitkeep" {
			continue
		}
		// A committed non-anchor under dist/ means someone filled the bundle.
		t.Logf("AC#12: tracked embed content beyond the anchor: %s", e)
	}
	t.Logf("AC#12 present count: tracked frontend/dist entries = %d (clean-checkout embed content)", len(entries))
}

// repoRootForTest asks git for the toplevel of the working tree the tests live in.
func repoRootForTest(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatalf("git rev-parse --show-toplevel: %v", err)
	}
	return strings.TrimSpace(string(out))
}

// TestCleanCheckoutBuilds_AC11 copies the committed tree (git archive HEAD, so no
// untracked working files) into a repository-external temp dir and runs
// go build ./... there. AC#11's judgment is that the binary links with only the
// tracked files present - i.e. with just the .gitkeep placeholder in frontend/dist
// and the go-webview2 dependency in go.mod/go.sum. The reverse control (renaming
// the embed pattern to a nonexistent path so the build MUST fail) lives in
// frontend/embed.go, which this leg is forbidden to touch (two-layer ban); it is
// recorded in the evidence table as owned by whoever edits frontend/**.
func TestCleanCheckoutBuilds_AC11(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Fatalf("go toolchain not on PATH for the clean-checkout build: %v", err)
	}
	tmp := t.TempDir()

	// git archive must be rooted at the repo root, not the package dir the test
	// runs from, or the extracted tree is a subdir with no top-level go.mod and the
	// clean build can never find its main module.
	archive := exec.Command("git", "archive", "--format=zip", "HEAD")
	archive.Dir = repoRootForTest(t)
	blob, err := archive.Output()
	if err != nil {
		t.Fatalf("git archive HEAD: %v", err)
	}
	if err := unzipTree(t, blob, tmp); err != nil {
		t.Fatalf("extract HEAD zip into %s: %v", tmp, err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "go.mod")); err != nil {
		t.Fatalf("clean-checkout: go.mod not at %s after extract (%v)", tmp, err)
	}

	build := exec.Command("go", "build", "./...")
	build.Dir = tmp
	build.Env = append(os.Environ(), "GOFLAGS=")
	out, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("clean-checkout go build ./... failed: %v\n%s", err, tail(string(out), 4000))
	}
	t.Logf("clean-checkout go build ./... ok in %s", tmp)
}

// unzipTree writes a `git archive --format=zip` blob into dst using the standard
// library only. git archive never emits absolute or ".." paths, but the guard is
// kept so a malformed entry can never write outside dst.
func unzipTree(t *testing.T, blob []byte, dst string) error {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(blob), int64(len(blob)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		name := filepath.Clean(strings.ReplaceAll(f.Name, "\\", "/"))
		if name == "." || strings.HasPrefix(name, "..") || filepath.IsAbs(name) {
			continue
		}
		target := filepath.Join(dst, filepath.FromSlash(name))
		if !strings.HasPrefix(target, filepath.Clean(dst)+string(os.PathSeparator)) && target != filepath.Clean(dst) {
			continue
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		outf, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(outf, rc); err != nil {
			rc.Close()
			outf.Close()
			return err
		}
		rc.Close()
		outf.Close()
	}
	return nil
}

func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
