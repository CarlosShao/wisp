package memory

// Ticket 76 — the artifacts path invariant, turned into an executed test.
//
// The claim being pinned (and the reason ticket 20's "attempting a user-dir
// write via the artifacts API -> denied" box could not be ticked as written) is
// that this API has NO caller-controlled destination path: the only string a
// caller ever supplies is a bare artifact NAME, and the directory it is resolved
// against is the one memory.Open fixed at construction. So there is no "attempt"
// to make, and the honest replacement for that box is a positive invariant plus
// every hostile component shape being proven rejected here.
//
// Three layers of proof, deliberately redundant:
//   1. an AST audit of the source itself (no path-shaped parameter, every
//      filesystem-mutating call rooted at s.artifactsDir, every internal
//      dir-taking helper called only with s.artifactsDir);
//   2. every hostile shape (separator / .. / drive letter / UNC) each in its
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
	"runtime"
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

// inv76Sep names a nested component the way the RUNNING platform spells
// nesting, so "a file one directory below the artifacts dir" is expressed once
// and means the same physical thing on Windows and on Linux (ticket 81 AC#1:
// the fixtures used to hard-code `\`, which is a directory boundary only on
// Windows and an ordinary file-name character on Linux).
func inv76Sep(parts ...string) string {
	return strings.Join(parts, string(filepath.Separator))
}

// inv76Shapes builds the hostile caller-supplied components against a concrete
// layout:
//
//	<root>/data/artifacts/          the store's artifacts dir
//	<root>/data/                    the store's data dir
//	<root>/Users/carlos/Documents/  a user dir the caller aims at
//
// Every target is `filepath.Join(<dir the caller aims at>, <the name the caller
// supplied>)`, i.e. the path the running OS actually resolves that string to —
// that is what "the file the caller named" means, and it is the only honest
// form on both platforms.
func inv76Shapes(artifactsDir, dataDir, userDir string) []inv76Shape {
	driveTarget := filepath.Join(userDir, "canary-drive-letter.txt")
	uncOf := func(p string) string {
		// C:\Users\... -> \\localhost\C$\Users\...
		// A UNC path is a Windows-shaped STRING, and that is fine: the claim
		// pinned here is "the guard refuses it", which holds identically on
		// both platforms (validArtifactName rejects `\` and `:` outright).
		return `\\localhost\` + string(p[0]) + `$` + p[2:]
	}
	return []inv76Shape{
		{
			sub:    "separator",
			name:   inv76Sep("nested", "canary-separator.txt"),
			target: filepath.Join(artifactsDir, "nested", "canary-separator.txt"),
		},
		{
			// The literal `\` spelling, KEPT ON PURPOSE. It means "nested" on
			// Windows and "one flat file whose name contains a backslash" on
			// Linux, and that asymmetry is itself under test — see
			// TestArtifactsLiteralBackslashIsNotAFixedSetOfBytes.
			sub:    "separator_backslash_literal",
			name:   `nested-backslash\canary-literal.txt`,
			target: filepath.Join(artifactsDir, `nested-backslash\canary-literal.txt`),
		},
		{
			sub:    "dotdot",
			name:   inv76Sep("..", "canary-dotdot.txt"),
			target: filepath.Join(dataDir, "canary-dotdot.txt"),
		},
		{
			sub:    "dotdot_backslash_literal",
			name:   `..\canary-dotdot-literal.txt`,
			target: filepath.Join(artifactsDir, `..\canary-dotdot-literal.txt`),
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
	pathLike := []string{
		"path", "dir", "dirname", "dest", "dst", "target", "root",
		"basedir", "folder", "filename", "file", "full", "prefix", "joined", "abs", "where",
	}

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
	mutators := map[string]bool{
		"Remove": true, "RemoveAll": true, "MkdirAll": true,
		"WriteFile": true, "Create": true, "CreateTemp": true, "OpenFile": true, "Rename": true,
		"Truncate": true, "Chmod": true, "Chtimes": true,
	}
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
// Layer 2 — every hostile shape, one named subtest each: six of them, because
// AC#1 (ticket 81) keeps the literal-`\` spelling of the two separator kinds
// beside the platform-neutral one. The two test names below still say "Four",
// which docs/evidence/s1/76-adversarial-acceptance.md cites verbatim; renaming
// them would rewrite someone else's evidence trail, so inv76Shapes is the source
// of truth and both tests just loop over whatever it returns.
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
// artifacts entry point with every hostile shape, a recursive listing of the
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
	// The expectation is derived from the paths the fixture actually planted,
	// not from a hard-coded list of slash-keys that only one platform can
	// produce: a canary the OS resolved into a subdirectory contributes that
	// directory as well, a canary it resolved to one flat name contributes only
	// itself, and a `..\`-shaped canary is only inside this tree on the platform
	// that does not read the backslash as a boundary (Linux - ticket 79's
	// encoder is precisely about that). Which of those each shape is happens
	// per-platform, so the membership test below is the prefix, not a name list.
	// The set equality still has the same teeth: an unexpected removal anywhere -
	// the DB, the data dir, the user dir - fails on the name.
	planted := []string{
		inv76Rel(root, filepath.Join(artifactsDir, "real-1.txt")), // the (d) control, already gone
		inv76Rel(root, filepath.Join(artifactsDir, "real-2.txt")),
	}
	for _, sh := range shapes {
		if k := inv76Rel(root, sh.target); strings.HasPrefix(k, artTree) {
			planted = append(planted, k)
		}
	}
	if len(planted) < 4 {
		t.Fatalf("only %d canaries resolve inside the artifacts tree, the set equality is too small to bite", len(planted))
	}
	wantRemoved := inv76ReclaimKeys(artTree, planted...)
	if !slices.Equal(removed3, wantRemoved) {
		t.Errorf("purge removed %v, want exactly %v", removed3, wantRemoved)
	}
	// PurgeArtifacts counts FILES it reclaimed, not directories: the listing
	// already proved which entries left the tree, so the count is asserted
	// against the file entries the fixture planted in the artifacts tree (minus
	// the one (d) deleted), which is again derived, not hard-coded per platform.
	wantFiles := 0
	for k, isDir := range before {
		if isDir || !strings.HasPrefix(k, artTree) {
			continue
		}
		wantFiles++
	}
	if wantFiles < 3 {
		t.Fatalf("the fixture planted only %d artifact files, the count cannot discriminate", wantFiles)
	}
	if n != wantFiles-1 {
		t.Errorf("purge removed %d files, want %d (every artifact-tree file in the fixture except the one (d) deleted; entries %v)",
			n, wantFiles-1, wantRemoved)
	}
	for _, sh := range []string{"nested", "nested-backslash"} {
		if _, err := os.Stat(filepath.Join(artifactsDir, sh)); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("the stray %q dir survived a purge (err=%v); ticket 79 says it must not", sh, err)
		}
	}
	if entries, err := os.ReadDir(artifactsDir); err != nil || len(entries) != 0 {
		t.Errorf("artifacts dir is not empty after a purge: %d entries left, err=%v", len(entries), err)
	}
	if _, err := os.Stat(dataDir); err != nil {
		t.Errorf("the data dir must survive a purge: %v", err)
	}
}

// TestArtifactsLiteralBackslashKeepsItsAsymmetry is ticket 81 AC#1's other half:
// the literal-`\` fixture must STAY platform-specific, deliberately. The same
// caller-supplied bytes are a directory boundary on Windows and an ordinary
// character inside one file name on Linux — where ticket 79's encoder percent-
// escapes them, so a long flat name is the CORRECT artifact of that route.
// Pinning the asymmetry is what stops a future edit from "fixing portability" by
// deleting the backslash case (AC#1 forbids that), and it is what stops the
// platform-neutral separator case from quietly being the only thing measured.
//
// There is no build tag on this file and there must not be one: this test is the
// proof that both platforms run it and see different, asserted, physics.
func TestArtifactsLiteralBackslashKeepsItsAsymmetry(t *testing.T) {
	s, root, _, _, shapes := inv76Fixture(t)
	artifactsDir := s.ArtifactsDir()
	const artTree = "data/artifacts/"

	neu := inv76ShapeNamed(t, shapes, "separator")

	// Each target is defined as "where this OS resolves the caller's string", so
	// every canary must be exactly there on every platform — the part that used to
	// be hard-coded with a `\` and therefore only ever held on Windows.
	for _, sh := range shapes {
		if _, err := os.Stat(sh.target); err != nil {
			t.Fatalf("shape %q: the file the caller named is not at %s on this platform: %v",
				sh.sub, sh.target, err)
		}
	}
	// The separator shape nests on BOTH platforms; that is the whole point of
	// building its name from filepath.Separator.
	if d := filepath.Dir(neu.target); !strings.HasSuffix(d, "nested") {
		t.Errorf("the platform-neutral shape resolved to %s, not one level under a nested dir", d)
	}
	if _, err := os.Stat(filepath.Dir(neu.target)); err != nil {
		t.Errorf("the platform-neutral separator shape did not create a directory: %v", err)
	}

	// A name carrying `\` that this OS does not read as a boundary is ONE flat
	// file name; the same bytes on Windows are two path components. Derived from
	// filepath.Separator, so "which side we are on" is asserted, not assumed.
	var flats []string
	for _, sh := range shapes {
		if strings.Contains(sh.name, `\`) && !strings.Contains(sh.name, string(filepath.Separator)) {
			flats = append(flats, sh.name)
		}
	}
	slices.Sort(flats)
	if filepath.Separator == '\\' {
		if len(flats) != 0 {
			t.Fatalf("Windows reads `\\` as a separator, so no shape may be flat here; got %v", flats)
		}
		if _, err := os.Stat(filepath.Join(artifactsDir, "nested-backslash")); err != nil {
			t.Errorf("the literal-backslash separator shape made no directory on Windows: %v", err)
		}
		if _, err := os.Stat(filepath.Join(artifactsDir, "nested-backslash", "canary-literal.txt")); err != nil {
			t.Errorf("the literal-backslash separator shape is not nested on Windows: %v", err)
		}
	} else {
		want := []string{`..\canary-dotdot-literal.txt`, `nested-backslash\canary-literal.txt`}
		if !slices.Equal(flats, want) {
			t.Fatalf("on Linux both literal-backslash shapes are flat names; this fixture produced %v, want %v", flats, want)
		}
		for _, f := range flats {
			if _, err := os.Lstat(filepath.Join(artifactsDir, f)); err != nil {
				t.Errorf("flat literal-backslash entry %q is missing: %v", f, err)
			}
			if k := inv76Rel(root, filepath.Join(artifactsDir, f)); !strings.HasPrefix(k, artTree) {
				t.Errorf("flat entry %q is not inside the artifacts tree (%s)", f, k)
			}
		}
		if _, err := os.Stat(filepath.Join(artifactsDir, "nested-backslash")); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("a backslash-shaped name created a DIRECTORY on a platform where the backslash is not a separator (err=%v)", err)
		}
	}

	// And the real listing agrees with the claim, entry by entry: no stray
	// backslash-named file beyond the pinned ones, on either platform.
	entries, err := os.ReadDir(artifactsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.Contains(e.Name(), `\`) {
			continue
		}
		if !slices.Contains(flats, e.Name()) {
			t.Errorf("a top-level entry name contains a backslash but is not a pinned fixture: %q (pinned: %v)", e.Name(), flats)
		}
	}
	t.Logf("GOOS=%s separator=%q: literal-backslash names living as ONE flat file name in the artifacts dir: %v",
		runtime.GOOS, filepath.Separator, flats)
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

// inv76ShapeNamed pulls one shape out of the fixture by subtest name, so a
// statement about a specific shape can say which one instead of re-deriving its
// path and drifting away from inv76Shapes.
func inv76ShapeNamed(t *testing.T, shapes []inv76Shape, sub string) inv76Shape {
	t.Helper()
	for _, sh := range shapes {
		if sh.sub == sub {
			return sh
		}
	}
	t.Fatalf("fixture has no shape %q: the assertion naming it is measuring nothing", sub)
	return inv76Shape{}
}

// inv76ReclaimKeys expands root-relative slash keys into the full set of
// entries a recursive reclaim of exactly those paths must take away: every key
// itself, plus every parent directory up to but not including artTree (the
// artifacts dir itself survives a purge - ticket 79 reclaims entries under it,
// never the root it was given).
//
// Deriving the expectation this way is what makes the containment claim
// platform-neutral: on Windows the literal-backslash canary brings its
// directory with it, on Linux it is one flat name and brings nothing, and both
// answers fall out of the same rule rather than out of a hard-coded list.
func inv76ReclaimKeys(artTree string, keys ...string) []string {
	seen := map[string]bool{}
	for _, k := range keys {
		if !strings.HasPrefix(k, artTree) {
			panic("inv76ReclaimKeys: key " + k + " is outside the artifacts tree")
		}
		for cur := k; ; {
			seen[cur] = true
			i := strings.LastIndex(cur, "/")
			if i < 0 {
				break
			}
			parent := cur[:i]
			if parent+"/" == artTree || seen[parent] {
				break // the artifacts dir itself is never reclaimed
			}
			cur = parent
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
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
