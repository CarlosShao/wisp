package agent

// Ticket 76 (agent side) — the spill route's path invariant, executed.
//
// spill.go:97-106 derives its whole file name from artifactName(callID, seq):
// the only caller-controlled piece is the model-supplied tool-call id, and the
// destination directory is fixed by NewSpiller. So, exactly as on the memory
// side, there is no caller-supplied DESTINATION - and the box in ticket 20 that
// says "attempting a user-dir write via the artifacts API -> denied" describes a
// threat this API cannot express. What this file pins instead:
//   1. the four hostile component shapes (separator / .. / drive letter / UNC)
//      each in a named subtest, each proven SANITIZED, with the resulting
//      on-disk name asserted exactly;
//   2. containment proven by a real recursive directory listing, with a positive
//      control that deliberately performs the unsanitized join so the listing
//      provably can see an escape if one happened;
//   3. the same again end to end against a real memory.Store data dir;
//   4. an AST audit of spill.go: no destination-shaped parameter outside the
//      host-side constructor, and the only path fed to os.WriteFile is the one
//      derived from s.dir + artifactName.

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
	"unicode"

	"github.com/CarlosShao/wisp/internal/memory"
)

// inv76aBudget is a hand-built small budget: 800 bytes of payload crosses the
// 40-token spill threshold, so a spill is guaranteed without a 16KB string.
var inv76aBudget = Budgets{RawOutputCapBytes: 1 << 20, SpillTokens: 40, SpillHeadTokens: 10, SpillTailTokens: 10}

const inv76aPayload = "PAYLOAD" // repeated below; content is irrelevant to the path

// inv76aShape is one hostile caller-supplied component plus the name the
// ENCODER must produce from it.
//
// TICKET 79 changed these literals, and only these literals: the names below are
// what the ids USED to collapse onto ("p/q" and "pq" both became
// tool-output-pq.txt), and that collapse was the clobbering defect this repo is
// fixing. The four shapes are still asserted to come out BARE, with no / \ : ..
// anywhere and with the full payload inside - containment is exactly as pinned as
// it was; what changed is that nothing is silently dropped on the way.
type inv76aShape struct {
	sub      string // AC#1 subtest name
	callID   string // model-supplied tool-call id
	wantName string // the ONLY file name it may become
}

func inv76aShapes() []inv76aShape {
	return []inv76aShape{
		{
			sub: "separator",
			// Both spellings, because the encoder escapes them alike: the
			// forward-slash form first (a nested path), then a backslash form
			// that must not become a second directory level either.
			callID:   "p/q",
			wantName: "tool-output-p%2Fq.txt",
		},
		{
			sub: "dotdot",
			// The pure escape spelling. Every dot escapes too, which is what
			// kills the class outright: a dot-component cannot appear in an
			// artifact name at any position, so neither can the Windows
			// trailing-dot alias ticket 76 fixed on the delete side.
			callID:   "../../escape",
			wantName: "tool-output-%2E%2E%2F%2E%2E%2Fescape.txt",
		},
		{
			sub:      "drive_letter",
			callID:   `C:\Windows\System32\drop`,
			wantName: "tool-output-%43%3A%5C%57indows%5C%53ystem32%5Cdrop.txt",
		},
		{
			sub:      "unc",
			callID:   `\\fileserver\share\payload`,
			wantName: "tool-output-%5C%5Cfileserver%5Cshare%5Cpayload.txt",
		},
	}
}

// TestSpillCallIDHostileShapesSanitizedToBareNames is AC#1 on the spill route:
// every shape is ENCODED (not rejected - a spill must still happen), and the
// exact on-disk name is asserted, plus that the name is a bare file name in the
// dir the constructor fixed and holds the full payload.
func TestSpillCallIDHostileShapesSanitizedToBareNames(t *testing.T) {
	root := t.TempDir()

	for _, sh := range inv76aShapes() {
		// A per-shape dir that does not exist yet: MkdirAll(s.dir) is on this
		// route too, so the created directory is part of what is asserted.
		dir := filepath.Join(root, "data", sh.sub, "artifacts")
		t.Run(sh.sub, func(t *testing.T) {
			sp, err := NewSpiller(dir, inv76aBudget).Prepare(sh.callID, strings.Repeat(inv76aPayload, 200))
			if err != nil {
				t.Fatalf("Prepare(%q) failed: %v", sh.callID, err)
			}
			if !sp.Spilled {
				t.Fatal("payload did not spill, so no path was exercised at all")
			}
			if sp.Name != sh.wantName {
				t.Errorf("on-disk name = %q, want %q", sp.Name, sh.wantName)
			}
			if sp.Path != filepath.Join(dir, sh.wantName) {
				t.Errorf("path = %q, want %q (the dir must be the constructor's, unchanged)", sp.Path, filepath.Join(dir, sh.wantName))
			}
			// The resulting name must be a bare file name: no separator, no
			// colon, no dot component that Windows could re-absorb.
			if filepath.Base(sp.Name) != sp.Name {
				t.Errorf("name %q is not bare", sp.Name)
			}
			for _, bad := range []string{"/", "\\", ":", ".."} {
				if strings.Contains(sp.Name, bad) {
					t.Errorf("sanitized name %q still contains %q", sp.Name, bad)
				}
			}
			body, err := os.ReadFile(sp.Path)
			if err != nil {
				t.Fatalf("the artifact did not appear at the asserted path: %v", err)
			}
			if len(body) == 0 {
				t.Fatal("empty artifact")
			}
			// Real listing, not a string comparison: the dir holds exactly one
			// file and it is the asserted name.
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 || entries[0].Name() != sh.wantName {
				t.Errorf("artifacts dir listing = %v, want exactly [%s]", inv76aNames(entries), sh.wantName)
			}
		})
	}

	// The degenerate case the fallback exists for. TICKET 79 narrowed it and the
	// subtest is renamed to say so: an id made only of formerly-stripped
	// characters now ENCODES to something non-empty, so the sequence fallback is
	// left with exactly one job - the empty id, the only input that still encodes
	// to nothing. What this subtest used to be asserting was itself an instance of
	// the defect: all four ids stripped to "", all four fresh Spillers restarted at
	// sequence 1, and their artifacts overwrote each other in this one directory
	// while the test congratulated them on "the fallback name". Bare and
	// non-degenerate is still what is pinned - plus that five ids get five names.
	t.Run("encoded_shapes_stay_bare_and_empty_id_falls_back_to_sequence", func(t *testing.T) {
		seen := map[string]string{}
		for _, id := range []string{"", "....", "..", `\\`, "!!!"} {
			sp, err := NewSpiller(filepath.Join(root, "fallback"), inv76aBudget).Prepare(id, strings.Repeat(inv76aPayload, 200))
			if err != nil {
				t.Fatalf("Prepare(%q): %v", id, err)
			}
			if id == "" && !strings.HasPrefix(sp.Name, "tool-output-seq") {
				t.Errorf("empty id produced %q, want the sequence fallback name", sp.Name)
			}
			if sp.Name == "tool-output-.txt" || sp.Name == "tool-output-..txt" {
				t.Errorf("id %q produced the degenerate name %q", id, sp.Name)
			}
			if filepath.Base(sp.Name) != sp.Name {
				t.Errorf("id %q produced the non-bare name %q", id, sp.Name)
			}
			for _, bad := range []string{"/", "\\", ":", ".."} {
				if strings.Contains(sp.Name, bad) {
					t.Errorf("id %q produced a name containing %q: %q", id, bad, sp.Name)
				}
			}
			if prev, dup := seen[sp.Name]; dup {
				t.Errorf("ids %q and %q collide on the name %q", prev, id, sp.Name)
			}
			seen[sp.Name] = id
		}
		if len(seen) != 5 {
			t.Errorf("%d distinct names for 5 distinct ids, want 5 (name -> id: %v)", len(seen), seen)
		}
	})
}

// TestSpillContainmentByDirectoryListing is AC#2 on the spill route: after
// spilling with every hostile shape plus one id engineered to walk out of the
// artifacts dir, a recursive listing of the fixture root must show that nothing
// appeared outside the data dir, and the positive control at the end proves the
// listing would have seen an escape (so a green run is not an artifact of a
// blind probe).
func TestSpillContainmentByDirectoryListing(t *testing.T) {
	root := t.TempDir()
	dataDir := filepath.Join(root, "data")
	artifacts := filepath.Join(dataDir, "artifacts")
	userDir := filepath.Join(root, "Users", "carlos", "Documents")
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(userDir, "diary.txt"), []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}

	before := inv76aTree(t, root)
	userBefore := inv76aTree(t, userDir)
	sp := NewSpiller(artifacts, inv76aBudget)
	// Errors are collected, not fatal: under the AC#3 mutation the sanitizer is
	// gone and some of these calls fail on a broken path. Aborting there would
	// hide the thing this test exists to measure - the LISTING after the escape.
	var errs []error
	var callerNamed []string
	for _, sh := range inv76aShapes() {
		_, err := sp.Prepare(sh.callID, strings.Repeat(inv76aPayload, 200))
		if err != nil {
			errs = append(errs, fmt.Errorf("Prepare(%q): %w", sh.callID, err))
		}
		callerNamed = append(callerNamed, sh.callID)
	}
	// The nastiest accepted input: an id that, joined raw, resolves to
	// <root>\ESCAPE-API-<n>.txt - out of the artifacts dir, out of the data dir.
	for i := 0; i < 3; i++ {
		marker := "ESCAPE-API-" + string(rune('a'+i))
		id := inv76aEscapeID(t, artifacts, root, marker)
		callerNamed = append(callerNamed, id)
		if _, err := sp.Prepare(id, strings.Repeat(inv76aPayload, 200)); err != nil {
			errs = append(errs, fmt.Errorf("Prepare(escape id %q): %w", id, err))
		}
	}
	for _, e := range errs {
		t.Error(e)
	}

	added, removed := inv76aDiff(before, inv76aTree(t, root))
	if len(removed) != 0 {
		t.Errorf("the spill route removed %v, it has no business deleting anything", removed)
	}
	if len(added) == 0 {
		t.Fatal("the spill route wrote nothing at all: the probe is vacuous")
	}
	for _, a := range added {
		// The data dir itself may appear (the route MkdirAlls its own dir);
		// nothing outside it may.
		if a != "data" && !strings.HasPrefix(a, "data/") {
			t.Errorf("spill wrote OUTSIDE the data dir: %s", a)
		}
		if !inv76aIsBareArtifact(a) && a != "data" && a != "data/artifacts" {
			t.Errorf("spill wrote something that is not a bare artifact of the artifacts dir: %s", a)
		}
	}
	// (a) the caller's user dir is byte-identical, entry for entry.
	if ua, ur := inv76aDiff(userBefore, inv76aTree(t, userDir)); len(ua) != 0 || len(ur) != 0 {
		t.Errorf("the spill route disturbed the user dir: + %v / - %v", ua, ur)
	}
	// (b) nothing exists at any path a caller named, resolved exactly the way the
	// API resolves it - minus the sanitizer.
	for _, raw := range callerNamed {
		guess := filepath.Join(artifacts, "tool-output-"+raw+".txt")
		if _, err := os.Lstat(guess); err == nil {
			t.Errorf("a file appeared at the caller-named path %s", guess)
		}
	}
	for _, m := range []string{"ESCAPE-API-a.txt", "ESCAPE-API-b.txt", "ESCAPE-API-c.txt"} {
		if _, err := os.Lstat(filepath.Join(root, m)); err == nil {
			t.Errorf("an escape landed above the data dir at %s", filepath.Join(root, m))
		}
	}

	// (c) POSITIVE CONTROL: perform the very same join with the sanitizer
	// bypassed by hand. If the listing does not see THAT, the whole test above
	// is measuring nothing.
	control := inv76aEscapeID(t, artifacts, root, "CONTROL-escape")
	unsanitized := filepath.Join(artifacts, "tool-output-"+control+".txt")
	t.Logf("control: escape id %q resolves to %q (root is %q)", control, unsanitized, root)
	if err := os.WriteFile(unsanitized, []byte("deliberate"), 0o644); err != nil {
		t.Fatalf("control write failed (%s): %v - the escape id does not actually escape, so the control proves nothing", unsanitized, err)
	}
	afterControl := inv76aTree(t, root)
	if _, ok := afterControl["CONTROL-escape.txt"]; !ok {
		t.Fatalf("the listing never saw the deliberate escape outside the data dir (it went to %q); AC#2 would be vacuous",
			unsanitized)
	}
	if err := os.Remove(unsanitized); err != nil {
		t.Errorf("cleaning up the control file: %v", err)
	}
}

// TestSpillIntoRealStoreThenDeleteStaysUnderDataDir is the two-route end to end
// proof on one real data dir: the spill route writes (hostile ids), the memory
// artifacts route then deletes (hostile names), and a full listing of the root
// in between shows that neither route could touch anything outside the data dir.
func TestSpillIntoRealStoreThenDeleteStaysUnderDataDir(t *testing.T) {
	root := t.TempDir()
	store, err := memory.Open(filepath.Join(root, "data"),
		memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()
	userDir := filepath.Join(root, "Users", "carlos", "Documents")
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"diary.txt", "notes.md"} {
		if err := os.WriteFile(filepath.Join(userDir, n), []byte("mine"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	artifacts := store.ArtifactsDir()

	// Warm the DB files so they are inside the "before" snapshot.
	if _, err := store.ListArtifacts(context.Background()); err != nil {
		t.Fatal(err)
	}
	before := inv76aTree(t, root)

	sp := NewSpiller(artifacts, inv76aBudget)
	payload := strings.Repeat(inv76aPayload, 200)
	var produced []string
	for _, sh := range append(inv76aShapes(), inv76aShape{
		sub: "escape", callID: inv76aEscapeID(t, artifacts, root, "ESCAPE-E2E"), wantName: "x",
	}) {
		s, err := sp.Prepare(sh.callID, payload)
		if err != nil {
			t.Fatalf("Prepare(%q): %v", sh.callID, err)
		}
		produced = append(produced, s.Name)
		if !strings.HasPrefix(s.Path, artifacts+string(filepath.Separator)) {
			t.Errorf("spill path %q is not under %q", s.Path, artifacts)
		}
	}

	ctx := context.Background()
	abs := filepath.Join(userDir, "diary.txt")
	unc := `\\localhost\` + abs[:1] + "$" + abs[2:]
	for _, bad := range []string{`nested\canary.txt`, `..\canary.txt`, abs, unc} {
		if err := store.DeleteArtifact(ctx, bad); err == nil {
			t.Errorf("DeleteArtifact(%q) accepted a caller-named path", bad)
		}
	}
	// The legitimate deletes go through the same guard and must work.
	for _, n := range produced[:1] {
		if err := store.DeleteArtifact(ctx, n); err != nil {
			t.Errorf("legitimate delete of %q: %v", n, err)
		}
	}

	added, removed := inv76aDiff(before, inv76aTree(t, root))
	for _, a := range added {
		if !strings.HasPrefix(a, "data/artifacts/") {
			t.Errorf("the two routes together created %s outside data/artifacts", a)
		}
	}
	for _, r := range removed {
		if r != "data/artifacts/"+produced[0] {
			t.Errorf("unexpected removal %s", r)
		}
	}
	for _, n := range []string{"diary.txt", "notes.md"} {
		if _, err := os.Stat(filepath.Join(userDir, n)); err != nil {
			t.Errorf("user dir file %s did not survive the artifacts API: %v", n, err)
		}
	}
}

// TestSpillAPITakesNoCallerControlledDestinationPath audits spill.go the same way
// the memory side audits artifacts.go: the destination is only ever the
// constructor's dir, the model-supplied id is the single caller component, and
// no write call is fed a path built from anything else.
func TestSpillAPITakesNoCallerControlledDestinationPath(t *testing.T) {
	fset := token.NewFileSet()
	src, err := parser.ParseFile(fset, "spill.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	pathLike := []string{
		"path", "dir", "dirname", "dest", "dst", "target", "root",
		"basedir", "folder", "filename", "file", "prefix", "where", "output",
	}
	const constructor = "NewSpiller" // the one host-side place a dir is supplied

	render := func(n any) string {
		var buf bytes.Buffer
		_ = printer.Fprint(&buf, fset, n)
		return buf.String()
	}

	// Only the CALLER-FACING surface is audited for destination-shaped
	// parameters: writeFileExclusive(path, data) is the file's internal write
	// primitive and its argument is audited at the call site instead.
	for _, fd := range inv76aFuncDecls(src) {
		if !fd.Name.IsExported() || fd.Name.Name == constructor {
			continue
		}
		for _, fl := range fd.Type.Params.List {
			for _, nm := range fl.Names {
				if slices.Contains(pathLike, strings.ToLower(nm.Name)) {
					t.Errorf("%s: exported parameter %q is destination-shaped; only %s may name a directory",
						fd.Name.Name, nm.Name, constructor)
				}
			}
		}
	}

	// Sinks: the os mutators plus this file's own write primitive. Every sink
	// call's first argument must be a local whose initializer mentions BOTH
	// s.dir and the sanitizer artifactName(.
	mutators := map[string]bool{
		"WriteFile": true, "Create": true, "OpenFile": true,
		"MkdirAll": true, "Remove": true, "Rename": true,
	}
	localSinks := map[string]bool{"writeFileExclusive": true}
	checked := 0
	for _, fd := range inv76aFuncDecls(src) {
		if localSinks[fd.Name.Name] {
			continue // its body's `path` is a parameter; the call site is what matters
		}
		derived := map[string]string{}
		ast.Inspect(fd, func(n ast.Node) bool {
			as, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			rhs := make([]string, 0, len(as.Rhs))
			for _, e := range as.Rhs {
				rhs = append(rhs, render(e))
			}
			for _, lhs := range as.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					derived[id.Name] = strings.Join(rhs, ", ")
				}
			}
			return true
		})
		ast.Inspect(fd, func(n ast.Node) bool {
			ce, ok := n.(*ast.CallExpr)
			if !ok || len(ce.Args) == 0 {
				return true
			}
			var sinkName string
			switch fn := ce.Fun.(type) {
			case *ast.SelectorExpr:
				if pkg, ok := fn.X.(*ast.Ident); ok && pkg.Name == "os" && mutators[fn.Sel.Name] {
					sinkName = "os." + fn.Sel.Name
				}
			case *ast.Ident:
				if localSinks[fn.Name] {
					sinkName = fn.Name
				}
			}
			if sinkName == "" {
				return true
			}
			checked++
			arg := render(ce.Args[0])
			if id, ok := ce.Args[0].(*ast.Ident); ok {
				arg = derived[id.Name]
				if arg == "" {
					t.Errorf("%s(%s): path local has no initializer in this function", sinkName, id.Name)
					return true
				}
				// One hop is not enough on this route: path :=
				// filepath.Join(s.dir, name) and name := artifactName(...) are two
				// separate statements, so expand identifiers to their
				// initializers before judging where the destination came from.
				arg = inv76aExpand(arg, derived)
			}
			if !strings.Contains(arg, "s.dir") {
				t.Errorf("%s(%s): destination is not derived from the constructor's s.dir", sinkName, arg)
			}
			if (sinkName == "os.WriteFile" || sinkName == "writeFileExclusive") &&
				!strings.Contains(arg, "artifactName(") {
				t.Errorf("%s(%s): the name half is not produced by the sanitizer artifactName()", sinkName, arg)
			}
			return true
		})
	}
	if checked == 0 {
		t.Fatal("audit found no filesystem-mutating call in spill.go: the probe is vacuous")
	}
	t.Logf("audited %d filesystem-mutating call(s) in spill.go", checked)
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// inv76aEscapeID builds a tool-call id that - if it ever reached a join
// unsanitized - resolves to <root>\<marker>.txt.
//
// The arithmetic is measured, not assumed: artifactName glues "tool-output-"
// onto the first "..", so that component becomes a LITERAL partial
// ("tool-output-..") which the next ".." cancels - one level is consumed by the
// prefix, then each remaining ".." pops one real component. dir sits `depth`
// components below root, so an escape needs depth+2 groups. Verified by this
// test's own positive control, which fails loudly if the guess is off by one.
//
// TICKET 79 note: the encoder now escapes every dot and separator, so feeding
// this id through Prepare lands a bare name INSIDE the artifacts dir - the
// escape below is only reachable through the hand-rolled join in the positive
// control, which is exactly where it belongs. It is kept as live input for the
// Prepare loop because "an id engineered to walk out" is the thing a containment
// test should still hand to the real route.
func inv76aEscapeID(t *testing.T, dir, root, marker string) string {
	t.Helper()
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		t.Fatal(err)
	}
	depth := strings.Count(rel, string(filepath.Separator)) + 1
	return strings.Repeat(`..\`, depth+2) + marker
}

func inv76aIsBareArtifact(rel string) bool {
	return strings.HasPrefix(rel, "data/artifacts/") &&
		strings.Count(strings.TrimPrefix(rel, "data/artifacts/"), "/") == 0
}

func inv76aTree(t *testing.T, root string) map[string]bool {
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

func inv76aDiff(before, after map[string]bool) (added, removed []string) {
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

func inv76aNames(entries []fs.DirEntry) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}

// inv76aExpand rewrites an initializer expression into itself plus the
// initializers of every identifier it references, to a bounded fixed point.
func inv76aExpand(text string, derived map[string]string) string {
	seen := map[string]bool{}
	out := text
	for round := 0; round < 4; round++ {
		grew := false
		for _, k := range slices.Sorted(maps.Keys(derived)) {
			if seen[k] || !inv76aUsesIdent(out, k) {
				continue
			}
			seen[k] = true
			out += " /*" + k + "=*/ " + derived[k]
			grew = true
		}
		if !grew {
			break
		}
	}
	return out
}

// inv76aUsesIdent reports whether name appears as a whole identifier in text.
func inv76aUsesIdent(text, name string) bool {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return !(r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r))
	})
	return slices.Contains(fields, name)
}

func inv76aFuncDecls(f *ast.File) []*ast.FuncDecl {
	var out []*ast.FuncDecl
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
			out = append(out, fd)
		}
	}
	return out
}
