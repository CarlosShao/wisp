package main

// Platform-neutral gates for the panel host (ticket 33). These run in the ubuntu
// core scope too (J7: the L1 half must have a real CI denominator), which is why
// they read the Windows host SOURCE FILE as text rather than executing a WebView2
// window: an import that would open a socket is forbidden at the source level,
// independent of GOOS.

import (
	"archive/zip"
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
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

// TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 is AC#12's ruler, re-tooled
// by 33-r4. The version it replaces
// (TestEmbeddedDistCleanCheckoutHasPlaceholderOnly_AC12) contained ZERO t.Errorf:
// every finding - the non-anchor entries and the count - went to t.Logf, so it
// could not go red under any bundle at all (33-v1 §A#31: "the fails-closed claim
// in its comment has no counterpart in the code").
//
// AC#12's question is "can this instrument tell 'the embed only matched the
// placeholder' apart from 'there really is a page bundle'?". Both shapes are asked
// through Go-side capability doors (panel.Assets: Built / Resolve / Check /
// Manifest). Nothing here asks a build exit code, prints marketing text, or opens
// a frontend file: frontend/** is on the two-layer ban list for this leg, and the
// tracked/ignored axes are read as ENTRY NAMES from git only.
//
// The assertions are the two invariants that make the shapes separable:
//   - anchor-only shape  -> every door must fail closed (a placeholder served as a
//     page is the false green AC#12 exists to stop);
//   - page-bundle shape  -> the entry must resolve as html and the entry's own
//     references must resolve from the same tree (the "half-bundle" shape, which
//     the pre-33-r4 file only promised in a comment);
//   - provenance axis    -> a page cannot exist in the binary with nothing in this
//     tree to build it from, and committed bundle content cannot produce a
//     not-built binary.
//
// The reading today (docs/evidence/s1/33-panel-host-c27-r4.md §①格 5): this tree
// carries a page built from GITIGNORED working-tree products, so the ticked answer
// is still "not yet" - that is the orchestrator's归口 call (ticket 33 AC#12 says
// the box cannot be ticked until whoever fills dist is settled), not a red this
// leg gets to file on the frontend side.
func TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12(t *testing.T) {
	assets, err := panel.BuiltinAssets()
	if err != nil {
		t.Fatalf("panel.BuiltinAssets() cannot answer what this binary carries: %v", err)
	}
	built := assets.Built()
	entry, entryCT, entryErr := assets.Resolve(panel.EntryFile)
	refs, checkErr := assets.Check()
	manifest, manifestErr := assets.Manifest()

	root, haveGit := gitRepoRootForTest()
	var tracked, ignoredOrUntracked []string
	if haveGit {
		tracked, _ = gitLinesInDir(root, "ls-files", "frontend/dist")
		ignoredOrUntracked, _ = gitLinesInDir(root, "status", "--porcelain", "--ignored", "--", "frontend/dist")
	}
	trackedBeyondAnchor := 0
	trackedHasEntry := false
	for _, e := range tracked {
		if filepath.Base(filepath.ToSlash(e)) == ".gitkeep" {
			continue
		}
		trackedBeyondAnchor++
		if filepath.ToSlash(e) == "frontend/dist/"+panel.EntryFile {
			trackedHasEntry = true
		}
	}

	shape := "anchor-only"
	if built {
		shape = "page-bundle"
	}
	t.Logf("AC#12 reading (head %s): shape=%s built=%t entry-bytes=%d entry-ctype=%q entry-err=%v refs=%d check-err=%v manifest-entries=%d manifest-err=%v | git-metadata=%t tracked=%d tracked-beyond-anchor=%d tracked-has-entry=%t ignored-or-untracked=%d",
		gitHeadShortForTest(t), shape, built, len(entry), entryCT, entryErr, len(refs), checkErr, len(manifest), manifestErr,
		haveGit, len(tracked), trackedBeyondAnchor, trackedHasEntry, len(ignoredOrUntracked))

	if built {
		if entryErr != nil {
			t.Fatalf("the embed reports built=true yet Resolve(%q) failed: %v", panel.EntryFile, entryErr)
		}
		if len(entry) == 0 {
			t.Errorf("built=true but the embedded entry file is 0 bytes - that is a bundle that renders nothing while still claiming to be built")
		}
		if !strings.Contains(strings.ToLower(entryCT), "text/html") {
			t.Errorf("the embedded entry resolves as content type %q, want an html document", entryCT)
		}
		if checkErr != nil {
			t.Errorf("half-bundle: built=true and the entry resolves, but the entry's own asset references do not all resolve from the same embed tree: %v (resolved %d of them before failing)", checkErr, len(refs))
		}
		if manifestErr != nil || len(manifest) == 0 {
			t.Errorf("built=true yet Manifest() carried nothing (err=%v)", manifestErr)
		}
	} else {
		if entryErr == nil {
			t.Errorf("the embed reports built=false (placeholder only), yet Resolve(%q) handed back %d bytes - a placeholder-only bundle is being served as a page", panel.EntryFile, len(entry))
		}
		if checkErr == nil {
			t.Errorf("the embed reports built=false, yet Check() passed and resolved %d asset reference(s) - the not-built fail-closed door is open", len(refs))
		}
		if manifestErr == nil {
			t.Errorf("the embed reports built=false, yet Manifest() listed %d entr(y/ies) - a placeholder-only tree must not hand out a bundle manifest", len(manifest))
		}
	}

	if haveGit {
		if built && trackedBeyondAnchor == 0 && len(ignoredOrUntracked) == 0 {
			t.Errorf("AC#12's two shapes have just been conflated: the binary carries a page (built=true) while this tree tracks nothing under frontend/dist beyond the anchor and reports no untracked/ignored entry there either - there is nothing this embed could have been built from, so 'built' is not being decided by the bundle")
		}
		if !built && trackedBeyondAnchor > 0 {
			if trackedHasEntry {
				t.Errorf("%d bundle file(s) including the entry are committed under frontend/dist, yet the embed reports built=false - the go:embed pattern is not matching committed content", trackedBeyondAnchor)
			} else {
				t.Errorf("committed half-bundle: %d file(s) are tracked under frontend/dist but the entry %s is not - the binary would ship a bundle that can never render, which is the exact shape the pre-33-r4 file only promised to fail closed on", trackedBeyondAnchor, panel.EntryFile)
			}
		}
	} else {
		t.Logf("AC#12 provenance axis not measurable in this tree (no git metadata - e.g. a git-archive copy); the capability assertions above still ran")
	}
}

// TestAC1SessionDisposeHasAProductionTrigger_AC1 turns 33-v1 §A#14's reading
// ("NewPanelManager has exactly one call site in the whole repo, and it is a
// test") into an instrument, and 33-r5's wiring is what flipped its branch: the
// resident leg now constructs the host (cmd/wisp/panel_resident_windows.go), so the
// case stops skipping and asserts AC#1's second clause for real - a session that
// builds a panel host must also hold the place that tears the window down.
//
// Tightened per the orchestrator's ruling 4 (10-01 12:12), which named the exact
// false green to avoid: 33-r4's version counted ANY method call named Destroy in a
// production file, and the two sites it found were the host letting go of the
// WebView2 CONTROL inside its own file (the bind-failure branch, and PanelManager's
// own Destroy body). Those would have satisfied the rule even if no session ever
// called it. Now the receiver has to be the manager: either an identifier this file
// bound from a *PanelManager factory call, or a field named mgr/manager (the
// spelling the resident thread's teardown uses). A Destroy on anything else is still
// counted and printed, it just no longer passes.
//
// The positive control that proves the tightening has teeth is a deletion, not an
// addition: take the manager Destroy out of the resident teardown and this case goes
// red naming what it saw. Reading in docs/evidence/s1/33-panel-host-c27-r5.md.
func TestAC1SessionDisposeHasAProductionTrigger_AC1(t *testing.T) {
	ctorHits, teardownHits, controlDestroys := panelHostProductionSites(t)
	t.Logf("AC#1 dispose scan: %d production constructor(s) %v | %d manager-teardown site(s) %v | %d other .Destroy() call site(s), WebView2 control teardown and NOT a session dispose %v",
		len(ctorHits), ctorHits, len(teardownHits), teardownHits, len(controlDestroys), controlDestroys)
	if len(ctorHits) == 0 {
		t.Skipf("AC#1's 'session dispose destroys' clause is unreachable in the product path: this package has 0 non-test call sites of NewPanelManager (other Destroy sites found: %d). Explicit Destroy is asserted in TestPanelHostRealWindowHopAndLifecycle; the resident wiring is ticket 33's.", len(controlDestroys))
	}
	if len(teardownHits) == 0 {
		t.Errorf("a production site now constructs the panel host (%s) but no non-test file calls Destroy ON THE MANAGER. The %d .Destroy() call(s) this package does have [%s] release the WebView2 control inside the host's own file, which is not a session tearing its window down: AC#1's second clause would still be unmet and the panel would outlive the session",
			strings.Join(ctorHits, ", "), len(controlDestroys), strings.Join(controlDestroys, ", "))
	}
	t.Logf("AC#1 dispose reachability: %d production constructor(s) [%s], %d manager teardown site(s) [%s]",
		len(ctorHits), strings.Join(ctorHits, ", "), len(teardownHits), strings.Join(teardownHits, ", "))
}

// panelHostFactories are the functions that hand back a *PanelManager. An
// identifier bound from either of them means "this file holds the host".
var panelHostFactories = map[string]bool{
	"NewPanelManager":         true,
	"newResidentPanelManager": true,
}

// panelManagerFields are the field names this package stores the manager under, so
// rp.mgr.Destroy() reads as a manager teardown and w.Destroy() does not.
var panelManagerFields = map[string]bool{"mgr": true, "manager": true}

// panelHostProductionSites walks the NON-TEST Go files of this package (the go test
// working directory is the package dir) and returns file:line for (a) calls that
// construct the host, (b) Destroy calls whose receiver is the manager, and (c)
// Destroy calls on anything else. A parse failure is a failed measurement, not a
// zero.
func panelHostProductionSites(t *testing.T) (ctorHits, managerDestroyHits, otherDestroyHits []string) {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir for the dispose scan: %v", err)
	}
	fset := token.NewFileSet()
	scanned := 0
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatalf("parse production file %s for the dispose scan: %v", name, err)
		}
		scanned++
		holderNames := map[string]bool{}
		var ctors, destroys []*ast.CallExpr
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.AssignStmt:
				for _, rhs := range node.Rhs {
					call, ok := rhs.(*ast.CallExpr)
					if !ok {
						continue
					}
					fn, ok := call.Fun.(*ast.Ident)
					if !ok || !panelHostFactories[fn.Name] {
						continue
					}
					for _, lhs := range node.Lhs {
						if id, ok := lhs.(*ast.Ident); ok && id.Name != "_" {
							holderNames[id.Name] = true
						}
					}
				}
			case *ast.CallExpr:
				switch fun := node.Fun.(type) {
				case *ast.Ident:
					if panelHostFactories[fun.Name] {
						ctors = append(ctors, node)
					}
				case *ast.SelectorExpr:
					if fun.Sel != nil && fun.Sel.Name == "Destroy" {
						destroys = append(destroys, node)
					}
				}
			}
			return true
		})
		for _, call := range ctors {
			ctorHits = append(ctorHits, fmtPos(fset, name, call.Fun))
		}
		for _, call := range destroys {
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			pos := fmtPos(fset, name, sel.Sel)
			if managerReceiver(sel.X, holderNames) {
				managerDestroyHits = append(managerDestroyHits, pos)
				continue
			}
			otherDestroyHits = append(otherDestroyHits, pos)
		}
	}
	if scanned == 0 {
		t.Fatalf("the dispose scan parsed 0 production files - the instrument is not looking at the package it claims to cover")
	}
	return ctorHits, managerDestroyHits, otherDestroyHits
}

// managerReceiver reports whether x names the panel manager: an identifier this file
// bound from a factory call, or a field access whose field is one of the names the
// package stores the manager under.
func managerReceiver(x ast.Expr, holderNames map[string]bool) bool {
	switch rec := x.(type) {
	case *ast.Ident:
		return holderNames[rec.Name]
	case *ast.SelectorExpr:
		return rec.Sel != nil && panelManagerFields[rec.Sel.Name]
	default:
		return false
	}
}

func fmtPos(fset *token.FileSet, file string, pos ast.Node) string {
	p := fset.Position(pos.Pos())
	return file + ":" + strconv.Itoa(p.Line)
}

// TestEmbeddedDistCleanCheckoutHasPlaceholderOnly_AC12 was the pre-33-r4 AC#12
// ruler. It is replaced by
// TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12 above; the reason is in
// 33-v1 §A#31 - the body had no t.Errorf at all, so it could not go red whatever
// the embed carried.

// repoRootForTest asks git for the toplevel of the working tree the tests live in.
func repoRootForTest(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatalf("git rev-parse --show-toplevel: %v", err)
	}
	return strings.TrimSpace(string(out))
}

// gitRepoRootForTest is the non-fatal form: a git-archive copy of this tree (which
// is what TestCleanCheckoutBuilds_AC11 builds, and what the 33-r4 reverse controls
// run in) has no .git, and "no git metadata" is a different fact from "no tracked
// files". Callers must handle the false case explicitly.
func gitRepoRootForTest() (string, bool) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", false
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", false
	}
	return root, true
}

// gitHeadShortForTest reads the current commit at run time. The pre-33-r4 latency
// log line carried the literal text "HEAD 7a41db9b", so eleven runs taken on a
// later commit each claimed to be taken on that one (33-v1 §A#21). When there is
// no git metadata it returns HEAD-unknown rather than a value that would read like
// a real anchor.
func gitHeadShortForTest(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "HEAD-unknown"
	}
	head := strings.TrimSpace(string(out))
	if head == "" {
		return "HEAD-unknown"
	}
	return head
}

// gitLinesInDir runs git in dir and splits its stdout into lines. ok=false means
// git itself failed, which the caller must not read as "empty list".
func gitLinesInDir(dir string, args ...string) ([]string, bool) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, false
	}
	lines := strings.FieldsFunc(string(out), func(r rune) bool { return r == '\n' || r == '\r' })
	keep := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		// git status --porcelain --ignored prefixes each row with two columns
		// (`!! path`, `?? path`, ` M path`); strip to the path.
		if len(l) > 3 && (strings.HasPrefix(l, "!! ") || strings.HasPrefix(l, "?? ")) {
			l = l[3:]
		}
		keep = append(keep, l)
	}
	return keep, true
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
