package memory

// Ticket 76 — the artifacts path invariant, turned into an executed test.
//
// The claim being pinned (and the reason ticket 20's "attempting a user-dir
// write via the artifacts API -> denied" box could not be ticked as written) is
// that this API has NO caller-controlled destination path: the only string a
// caller ever supplies is a bare artifact NAME, and the directory it is resolved
// against is the one memory.Open fixed at construction. So there is no "attempt"
// to make, and the honest replacement for that box is a positive invariant plus
// the four hostile component shapes being proven rejected here.
//
// Three layers of proof, deliberately redundant:
//   1. an AST audit of the source itself (no path-shaped parameter, every
//      filesystem-mutating call rooted at s.artifactsDir, every internal
//      dir-taking helper called only with s.artifactsDir);
//   2. the four hostile shapes (separator / .. / drive letter / UNC) each in a
//      named subtest, each asserted rejected BY THE GUARD (sentinel error, not a
//      filesystem error) with the file the caller named still on disk;
//   3. containment by real recursive directory listing (no string comparison),
//      including a positive control so the listing diff cannot pass vacuously.

import (
	"bytes"
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

// inv76Shape is one hostile caller-supplied component plus the on-disk file it
// names. The target really exists so that "rejected" is observable: if the guard
// were gone the file would be gone too.
type inv76Shape struct {
	sub    string // subtest name — one per AC#1 shape
	name   string // what the caller hands the API
	target string // absolute path the caller is naming
}

// inv76Shapes builds the four AC#1 shapes against a concrete layout:
//
//	<root>\data\artifacts\           the store's artifacts dir
//	<root>\data\                     the store's data dir
//	<root>\Users\carlos\Documents\   a user dir the caller aims at
func inv76Shapes(artifactsDir, dataDir, userDir string) []inv76Shape {
	driveTarget := filepath.Join(userDir, "canary-drive-letter.txt")
	uncOf := func(p string) string {
		// C:\Users\... -> \\localhost\C$\Users\...
		return `\\localhost\` + string(p[0]) + `$` + p[2:]
	}
	return []inv76Shape{
		{
			sub:    "separator",
			name:   `nested\canary-separator.txt`,
			target: filepath.Join(artifactsDir, `nested\canary-separator.txt`),
		},
		{
			sub:    "dotdot",
			name:   `..\canary-dotdot.txt`,
			target: filepath.Join(dataDir, "canary-dotdot.txt"),
		},
		{
			sub:    "drive_letter",
			name:   driveTarget, // an absolute path, drive letter included
			target: driveTarget,
		},
		{
			sub:    "unc",
			name:   uncOf(filepath.Join(userDir, "canary-unc.txt")),
			target: filepath.Join(userDir, "canary-unc.txt"),
		},
	}
}

// inv76Fixture lays out root\data (store), root\Users\carlos\Documents (user dir)
// and plants one canary per hostile shape plus two legitimate artifacts.
func inv76Fixture(t *testing.T) (*Store, string, string, string, []inv76Shape) {
	t.Helper()
	root := t.TempDir()
	dataDir := filepath.Join(root, "data")
	userDir := filepath.Join(root, "Users", "carlos", "Documents")
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		t.Fatal(err)
	}
	s, err := Open(dataDir, WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	artifactsDir := s.ArtifactsDir()

	shapes := inv76Shapes(artifactsDir, dataDir, userDir)
	for _, sh := range shapes {
		if err := os.MkdirAll(filepath.Dir(sh.target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(sh.target, []byte("CANARY "+sh.sub), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, n := range []string{"real-1.txt", "real-2.txt"} {
		if err := os.WriteFile(filepath.Join(artifactsDir, n), []byte("artifact"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return s, root, dataDir, userDir, shapes
}

// ---------------------------------------------------------------------------
// Layer 1 — the API surface itself takes no destination.
// ---------------------------------------------------------------------------

// TestArtifactsAPITakesNoCallerControlledDestinationPath audits artifacts.go
// and its call sites syntactically: no artifact entry point accepts a path/dir/
// dest-shaped parameter, and every filesystem-mutating call in the file is
// rooted at the store's own artifactsDir. This is the invariant ticket 20's box
// needs, stated as something the compiler-visible source either does or does
// not do.
func TestArtifactsAPITakesNoCallerControlledDestinationPath(t *testing.T) {
	src := inv76ParseFile(t, "artifacts.go")

	// Every parameter name that would mean "the caller chooses where this goes".
	pathLike := []string{"path", "dir", "dirname", "dest", "dst", "target", "root",
		"basedir", "folder", "filename", "file", "full", "prefix", "joined", "abs", "where"}

	// listArtifactsDir is the one internal helper that does take a directory, so
	// the audit must (a) allow it and (b) prove nobody can feed it anything but
	// the store's own field.
	dirTakingHelpers := map[string]bool{"listArtifactsDir": true}

	exportedBad := 0
	ast.Inspect(src, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if !ok || fd.Name == nil {
			return true
		}
		for _, fl := range fd.Type.Params.List {
			for _, nm := range fl.Names {
				if !dirTakingHelpers[fd.Name.Name] {
					if slices.Contains(pathLike, nm.Name) {
						t.Errorf("%s: parameter %q is a caller-controlled destination (AC#1 invariant violated)",
							fd.Name.Name, nm.Name)
						exportedBad++
					}
				}
			}
		}
		return true
	})
	if exportedBad != 0 {
		t.Fatalf("%d path-shaped parameters in artifacts.go", exportedBad)
	}

	// No exported artifact method takes MORE than one string: a (dir, name) pair
	// is exactly how a caller-controlled destination would be spelled.
	ast.Inspect(src, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if !ok || !fd.Name.IsExported() {
			return true
		}
		if !strings.Contains(fd.Name.Name, "Artifact") {
			return true
		}
		var strs []string
		ast.Inspect(fd.Type.Params, func(x ast.Node) bool {
			if id, ok := x.(*ast.Ident); ok && id.Name == "string" {
				strs = append(strs, "<string>")
			}
			return true
		})
		if len(strs) > 1 {
			t.Errorf("%s takes %d string parameters, want at most 1 (one component, no destination)",
				fd.Name.Name, len(strs))
		}
		return true
	})

	// Every filesystem-mutating call in the file must be rooted at
	// s.artifactsDir, either directly or through a local assigned from a
	// filepath.Join(s.artifactsDir, ...) in the same function.
	mutators := map[string]bool{"Remove": true, "RemoveAll": true, "MkdirAll": true,
		"WriteFile": true, "Create": true, "CreateTemp": true, "OpenFile": true, "Rename": true,
		"Truncate": true, "Chmod": true, "Chtimes": true}
	checked := 0
	for _, fn := range inv76Funcs(src) {
		derived := inv76DirDerivedVars(fn, "s.artifactsDir")
		ast.Inspect(fn, func(n ast.Node) bool {
			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := ce.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok || pkg.Name != "os" || !mutators[sel.Sel.Name] {
				return true
			}
			checked++
			if len(ce.Args) == 0 {
				t.Errorf("%s.%s: no path argument?", pkg.Name, sel.Sel.Name)
				return true
			}
			arg := ce.Args[0]
			argTxt := inv76Render(arg)
			base := argTxt
			if id, ok := arg.(*ast.Ident); ok {
				base = id.Name
			}
			rooted := strings.Contains(argTxt, "s.artifactsDir") || derived[base]
			if !rooted {
				t.Errorf("%s(%s) mutates a path not derived from s.artifactsDir (AC#1 invariant violated); dir-derived locals here: %v",
					sel.Sel.Name, argTxt, inv76Keys(derived))
			}
			return true
		})
	}
	if checked == 0 {
		t.Fatal("audit found no filesystem-mutating call in artifacts.go: the probe is vacuous")
	}
	t.Logf("audited %d filesystem-mutating call(s), all rooted at s.artifactsDir", checked)

	// The dir-taking helper must be fed only the store field, package-wide.
	callsites := 0
	for _, f := range inv76PackageFiles(t) {
		ast.Inspect(f.src, func(n ast.Node) bool {
			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok || !dirTakingHelpers[id.Name] {
				return true
			}
			callsites++
			if len(ce.Args) != 1 {
				t.Errorf("%s: %d args, want 1", id.Name, len(ce.Args))
				return true
			}
			if txt := inv76Render(ce.Args[0]); txt != "s.artifactsDir" {
				t.Errorf("%s(%s): helper dir argument is not the store's own artifactsDir", id.Name, txt)
			}
			return true
		})
	}
	if callsites == 0 {
		t.Fatal("no listArtifactsDir call sites found: the probe is vacuous")
	}
	t.Logf("audited %d listArtifactsDir call site(s), all fed s.artifactsDir", callsites)
}

// ---------------------------------------------------------------------------
// Layer 2 — the four hostile shapes, one named subtest each.
// ---------------------------------------------------------------------------

// TestDeleteArtifactRejectsTheFourHostileShapes is AC#1 on the delete route:
// separator / .. / drive letter / UNC each get a named subtest, each must be
// rejected by the name guard (not by the filesystem), and the file the caller
// named must still be there afterwards.
func TestDeleteArtifactRejectsTheFourHostileShapes(t *testing.T) {
	s, _, _, _, shapes := inv76Fixture(t)
	ctx := context.Background()

	for _, sh := range shapes {
		t.Run(sh.sub, func(t *testing.T) {
			inv76MustReject(t, "DeleteArtifact", sh.name, s.DeleteArtifact(ctx, sh.name), sh.target)
		})
	}
	// The degenerate spellings of the same idea, in the shape that owns them.
	t.Run("dotdot/bare_and_empty", func(t *testing.T) {
		for _, n := range []string{"..", ".", "", "...."} {
			err := s.DeleteArtifact(ctx, n)
			if err == nil {
				t.Errorf("DeleteArtifact(%q) = nil, want rejection", n)
				continue
			}
			if !errors.Is(err, ErrInvalidArtifactName) {
				t.Errorf("DeleteArtifact(%q) = %v, want the name guard's ErrInvalidArtifactName", n, err)
			}
		}
		if _, err := os.Stat(s.Dir()); err != nil {
			t.Errorf("the data dir itself must survive every name: %v", err)
		}
	})
}

// TestDeletePrivacyItemRejectsTheFourHostileShapes is AC#1 on the route the
// privacy page actually calls: the wrapper must not swallow the guard and must
// not add a path component of its own.
func TestDeletePrivacyItemRejectsTheFourHostileShapes(t *testing.T) {
	s, _, _, _, shapes := inv76Fixture(t)
	ctx := context.Background()

	for _, sh := range shapes {
		t.Run(sh.sub, func(t *testing.T) {
			inv76MustReject(t, "DeletePrivacyItem", sh.name,
				s.DeletePrivacyItem(ctx, PrivacyArtifacts, sh.name), sh.target)
		})
	}
}

func inv76MustReject(t *testing.T, route, name string, err error, target string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s(%q) = nil: a caller-named path was accepted (AC#1 shape not rejected)", route, name)
	}
	// Rejected BY THE GUARD, not by the filesystem: an ErrNotFound would mean the
	// path was built and touched.
	if !errors.Is(err, ErrInvalidArtifactName) {
		t.Errorf("%s(%q) = %v, want ErrInvalidArtifactName (guard said no before the FS was touched)",
			route, name, err)
	}
	if errors.Is(err, ErrNotFound) {
		t.Errorf("%s(%q) reached the filesystem: %v", route, name, err)
	}
	body, statErr := os.ReadFile(target)
	if statErr != nil {
		t.Fatalf("the file the caller named is gone (%v) even though the call failed: %v", statErr, err)
	}
	if !strings.HasPrefix(string(body), "CANARY") {
		t.Errorf("canary at %s was rewritten: %q", target, body)
	}
}

// ---------------------------------------------------------------------------
// Layer 3 — containment by real directory listing.
// ---------------------------------------------------------------------------

// TestArtifactsContainmentByDirectoryListing is AC#2: after driving every
// artifacts entry point with the four hostile shapes, a recursive listing of the
// whole fixture root must be byte-identical to the snapshot taken before, and
// nothing may exist at any path the caller named. The second half is a positive
// control: one legitimate delete must show up in the same diff, so an empty diff
// can never be the result of a broken probe.
func TestArtifactsContainmentByDirectoryListing(t *testing.T) {
	s, root, dataDir, userDir, shapes := inv76Fixture(t)
	ctx := context.Background()
	artifactsDir := s.ArtifactsDir()

	before := inv76Tree(t, root)
	userBefore := inv76Tree(t, userDir)
	if len(before) < 7 {
		t.Fatalf("fixture listing has only %d entries, the probe is vacuous", len(before))
	}

	for _, sh := range shapes {
		if err := s.DeleteArtifact(ctx, sh.name); err == nil {
			t.Errorf("DeleteArtifact(%q) accepted a caller-named path", sh.name)
		}
		if err := s.DeletePrivacyItem(ctx, PrivacyArtifacts, sh.name); err == nil {
			t.Errorf("DeletePrivacyItem(artifacts, %q) accepted a caller-named path", sh.name)
		}
	}
	// Read-side entry points too: they cannot name a path either.
	if _, err := s.ListArtifacts(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ListPrivacy(ctx, PrivacyArtifacts); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ExportPrivacy(ctx, PrivacyArtifacts); err != nil {
		t.Fatal(err)
	}

	// (a) nothing appeared, nothing disappeared, anywhere under root.
	after := inv76Tree(t, root)
	added, removed := inv76Diff(before, after)
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("hostile names changed the tree: + %v / - %v", added, removed)
	}
	// (b) the user dir in particular is untouched, entry for entry.
	if ua, ur := inv76Diff(userBefore, inv76Tree(t, userDir)); len(ua) != 0 || len(ur) != 0 {
		t.Errorf("the user dir changed under the artifacts API: + %v / - %v", ua, ur)
	}
	if len(userBefore) != 2 {
		t.Errorf("fixture broke: user dir holds %v, want the 2 planted canaries", inv76Keys(userBefore))
	}
	// (c) no file exists that the caller named and did not already plant: each
	// caller-named target is still exactly the canary the fixture wrote (this is
	// the delete-route form of "nothing appeared at the path the caller named").
	for _, sh := range shapes {
		body, err := os.ReadFile(sh.target)
		if err != nil {
			t.Errorf("caller-named target %s (%s) is gone: %v", sh.sub, sh.target, err)
			continue
		}
		if want := "CANARY " + sh.sub; string(body) != want {
			t.Errorf("caller-named target %s holds %q, want %q", sh.sub, body, want)
		}
	}

	// (d) POSITIVE CONTROL: a legitimate bare name does delete exactly one file,
	// inside the artifacts dir, and the diff sees it.
	if err := s.DeleteArtifact(ctx, "real-1.txt"); err != nil {
		t.Fatalf("legitimate delete failed: %v", err)
	}
	added2, removed2 := inv76Diff(before, inv76Tree(t, root))
	if len(added2) != 0 {
		t.Errorf("legitimate delete added %v", added2)
	}
	want := inv76Rel(root, filepath.Join(artifactsDir, "real-1.txt"))
	if !slices.Equal(removed2, []string{want}) {
		t.Errorf("legitimate delete removed %v, want exactly [%s] (the diff probe is not seeing the tree)",
			removed2, want)
	}
	for _, r := range removed2 {
		if !inv76InsideDir(r, artifactsDir, root) {
			t.Errorf("removal %s escaped the artifacts dir", r)
		}
	}

	// (e) purge must be containment-safe with a nested canary on disk, and
	// TICKET 79 changed what it is allowed to leave behind: a stray subdirectory
	// used to be invisible to the listing (and so to the quota and to this purge,
	// which is how a 500MB cap becomes 2GB and why "one-click clear" kept bytes).
	// It is now reclaimed too - so the containment question gets a sharper form:
	// not "did purge skip the subtree" but "did every path purge took away sit
	// inside the artifacts dir, at any depth, and was the set exactly this".
	//
	// inv76InsideDir only accepts a DIRECT child, which is the wrong shape for a
	// recursive reclaim, so the loop below prefixes on the artifacts tree and the
	// set equality right after it is what carries the teeth: an unexpected removal
	// anywhere - the data dir, the DB, the user dir - fails on the name.
	n, err := s.PurgeArtifacts(ctx)
	if err != nil {
		t.Fatalf("PurgeArtifacts: %v", err)
	}
	if n != 2 {
		t.Errorf("purge removed %d entries, want 2 (the 1 remaining file + the stray nested dir)", n)
	}
	added3, removed3 := inv76Diff(before, inv76Tree(t, root))
	if len(added3) != 0 {
		t.Errorf("purge added %v", added3)
	}
	const artTree = "data/artifacts/"
	for _, r := range removed3 {
		if !strings.HasPrefix(r, artTree) {
			t.Errorf("purge removed %s, outside the artifacts dir", r)
		}
	}
	wantRemoved := []string{
		artTree + "nested",
		artTree + "nested/canary-separator.txt",
		artTree + "real-1.txt", // the (d) positive control, still in this diff
		artTree + "real-2.txt",
	}
	if !slices.Equal(removed3, wantRemoved) {
		t.Errorf("purge removed %v, want exactly %v", removed3, wantRemoved)
	}
	if _, err := os.Stat(filepath.Join(artifactsDir, "nested")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the stray nested dir survived a purge (err=%v); ticket 79 says it must not", err)
	}
	if entries, err := os.ReadDir(artifactsDir); err != nil || len(entries) != 0 {
		t.Errorf("artifacts dir is not empty after a purge: %d entries left, err=%v", len(entries), err)
	}
	if _, err := os.Stat(dataDir); err != nil {
		t.Errorf("the data dir must survive a purge: %v", err)
	}
}

// ---------------------------------------------------------------------------
// listing / parsing helpers
// ---------------------------------------------------------------------------

// inv76Tree returns every path under root (files AND directories), relative to
// root, cleaned and separator-normalised, as a set.
func inv76Tree(t *testing.T, root string) map[string]bool {
	t.Helper()
	out := map[string]bool{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		out[filepath.ToSlash(filepath.Clean(rel))] = d.IsDir()
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return out
}

func inv76Diff(before, after map[string]bool) (added, removed []string) {
	for k := range after {
		if _, ok := before[k]; !ok {
			added = append(added, k)
		}
	}
	for k := range before {
		if _, ok := after[k]; !ok {
			removed = append(removed, k)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	return added, removed
}

func inv76Rel(root, abs string) string {
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return abs
	}
	return filepath.ToSlash(rel)
}

func inv76Abs(root, rel string) string {
	return filepath.Join(root, filepath.FromSlash(rel))
}

type inv76File struct {
	name string
	src  *ast.File
}

// inv76AuditFset is shared by every parse in this file: printer.Fprint needs the
// FileSet a node actually came from, so a per-call fset silently renders nothing.
var inv76AuditFset = token.NewFileSet()

func inv76ParseFile(t *testing.T, name string) *ast.File {
	t.Helper()
	f, err := parser.ParseFile(inv76AuditFset, name, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	return f
}

// inv76PackageFiles parses every non-test .go file in the package directory so
// call-site audits cover the whole package, not just one file.
func inv76PackageFiles(t *testing.T) []inv76File {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var out []inv76File
	for _, e := range entries {
		n := e.Name()
		if !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		out = append(out, inv76File{name: n, src: inv76ParseFile(t, n)})
	}
	if len(out) == 0 {
		t.Fatal("no package files found")
	}
	return out
}

func inv76Funcs(f *ast.File) []*ast.FuncDecl {
	var out []*ast.FuncDecl
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
			out = append(out, fd)
		}
	}
	return out
}

// inv76DirDerivedVars reports the local variables of fn whose initializer
// mentions root (e.g. full := filepath.Join(s.artifactsDir, name)).
func inv76DirDerivedVars(fn *ast.FuncDecl, root string) map[string]bool {
	out := map[string]bool{}
	ast.Inspect(fn, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		if !strings.Contains(inv76RenderExprs(as.Rhs), root) {
			return true
		}
		for _, lhs := range as.Lhs {
			if id, ok := lhs.(*ast.Ident); ok {
				out[id.Name] = true
			}
		}
		return true
	})
	return out
}

func inv76Render(n any) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, inv76AuditFset, n)
	return buf.String()
}

// inv76RenderExprs renders each expression separately: printer.Fprint is not
// guaranteed to accept a bare []ast.Expr, and a silently empty render here would
// turn the whole audit into a false negative.
func inv76RenderExprs(list []ast.Expr) string {
	parts := make([]string, 0, len(list))
	for _, e := range list {
		parts = append(parts, inv76Render(e))
	}
	return strings.Join(parts, ", ")
}

// inv76InsideDir reports whether the root-relative listing entry rel is a direct
// child of absDir (which is itself inside root).
func inv76InsideDir(rel, absDir, root string) bool {
	dir := filepath.ToSlash(inv76Rel(root, absDir)) + "/"
	if !strings.HasPrefix(rel, dir) {
		return false
	}
	return !strings.Contains(strings.TrimPrefix(rel, dir), "/")
}

func inv76Keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
