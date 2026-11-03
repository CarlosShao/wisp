package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// search_files_test.go + search_apps_test.go in one file, the way ticket 24's
// delivered half put doc.read and search.content tests together (see
// document_text_test.go). Ticket 24 AC#3 (`search.files`) and AC#4
// (`search.apps` feeding `app.launch`). Reuses the r2a seams already in the
// package: r2aCanon / r2aParams / r2aCorpus / mapResolver.

func r3None(string) risk.Class { return risk.ClassNone }

// r3Files builds one file under root and returns its path.
func r3Files(t *testing.T, root, rel string) string {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// r3Lines are the listing lines of a result, i.e. everything after the head.
func r3Lines(text string) []string {
	parts := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(parts) == 0 {
		return nil
	}
	return parts[1:]
}

func r3Exec(t *testing.T, d SearchDeps, args map[string]any) Result {
	t.Helper()
	res, err := searchFiles{d: d}.Execute(context.Background(), r2aParams(t, args), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	return res
}

// Test24R3SearchFilesMatchesNameGlobCaseInsensitive is AC#3's positive half:
// basename glob, case-insensitive, recursive inside the root, directories
// traversed but never returned.
func Test24R3SearchFilesMatchesNameGlobCaseInsensitive(t *testing.T) {
	root := t.TempDir()
	hit := []string{r3Files(t, root, "a.md"), r3Files(t, root, "B.MD"), r3Files(t, root, "sub/c.md")}
	skip := []string{r3Files(t, root, "notes.txt"), filepath.Join(root, "d.md")}
	if err := os.Mkdir(filepath.Join(root, "d.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	p, canon := r2aCanon(t, []string{root}, root)
	res := r3Exec(t, SearchDeps{Paths: p, Classify: r3None}, map[string]any{"path": canon, "pattern": "*.md"})
	if res.IsError {
		t.Fatalf("a valid search must not error: %q", res.Text)
	}
	if res.Origin != canon {
		t.Errorf("Origin must be the C26 canonical root, got %q want %q", res.Origin, canon)
	}
	listing := strings.Join(r3Lines(res.Text), "\n")
	for _, want := range hit {
		if !strings.Contains(listing, filepath.Base(want)) {
			t.Errorf("%q missing from listing:\n%s", want, res.Text)
		}
	}
	for _, name := range []string{filepath.Base(skip[0]), filepath.Base(skip[1])} {
		for _, line := range r3Lines(res.Text) {
			if strings.EqualFold(filepath.Base(line), name) {
				t.Errorf("%s must not be listed (dir / non-matching ext): %q", name, line)
			}
		}
	}
	if got := len(r3Lines(res.Text)); got != 3 {
		t.Errorf("want 3 hits, got %d in:\n%s", got, res.Text)
	}
}

// Test24R3SearchFilesChineseNameSubstring covers the ticket's own scenario
// wording ("找一下上周的那个合同") - a *substring* glob on a CJK name.
func Test24R3SearchFilesChineseNameSubstring(t *testing.T) {
	root := t.TempDir()
	yes := r3Files(t, root, "合同-2026Q4.pdf")
	no := r3Files(t, root, "invoice.pdf")
	p, canon := r2aCanon(t, []string{root}, root)
	res := r3Exec(t, SearchDeps{Paths: p, Classify: r3None}, map[string]any{"path": canon, "pattern": "*合同*"})
	listing := strings.Join(r3Lines(res.Text), "\n")
	if !strings.Contains(listing, filepath.Base(yes)) {
		t.Errorf("the CJK match is missing:\n%s", res.Text)
	}
	if strings.Contains(listing, filepath.Base(no)) {
		t.Errorf("a non-matching file was listed:\n%s", res.Text)
	}
}

// Test24R3SearchFilesEmptyPatternIsNotAListing is the fail-closed reverse: an
// empty pattern would turn one call into a whole-tree dump.
func Test24R3SearchFilesEmptyPatternIsNotAListing(t *testing.T) {
	root := t.TempDir()
	r3Files(t, root, "a.txt")
	p, canon := r2aCanon(t, []string{root}, root)
	res := r3Exec(t, SearchDeps{Paths: p, Classify: r3None}, map[string]any{"path": canon, "pattern": "  "})
	if !res.IsError {
		t.Errorf("empty pattern must be refused, got %+v", res)
	}
	if len(r3Lines(res.Text)) > 1 {
		t.Errorf("a refused search must return no listing lines: %q", res.Text)
	}
}

// Test24R3SearchFilesRejectsSeparatorPattern pins "NAME glob": a pattern with a
// path separator is not a name and is refused instead of becoming a path.
func Test24R3SearchFilesRejectsSeparatorPattern(t *testing.T) {
	root := t.TempDir()
	p, canon := r2aCanon(t, []string{root}, root)
	for _, pat := range []string{"sub/*.md", `sub\*.md`, "../etc/passwd"} {
		res := r3Exec(t, SearchDeps{Paths: p, Classify: r3None}, map[string]any{"path": canon, "pattern": pat})
		if !res.IsError {
			t.Errorf("separator pattern %q must be refused, got %+v", pat, res)
		}
	}
}

// Test24R3SearchFilesRejectsMalformedGlob: a broken bracket class is reported,
// not silently turned into "matched nothing".
func Test24R3SearchFilesRejectsMalformedGlob(t *testing.T) {
	root := t.TempDir()
	p, canon := r2aCanon(t, []string{root}, root)
	res := r3Exec(t, SearchDeps{Paths: p, Classify: r3None}, map[string]any{"path": canon, "pattern": "[a-z"})
	if !res.IsError || !strings.Contains(res.Text, "语法") {
		t.Errorf("a malformed glob must be reported as such, got %+v", res)
	}
}

// Test24R3SearchFilesRefusesOutsideAllowlist is AC#3's R2 half, in-tool: a root
// outside the authorized dirs is refused and echoes nothing back.
func Test24R3SearchFilesRefusesOutsideAllowlist(t *testing.T) {
	inside := t.TempDir()
	outside := r2aCorpus(t)
	p := NewPathCanonicalizer([]string{inside}, nil)
	if p == nil {
		t.Fatal("nil canonicalizer")
	}
	res, err := searchFiles{d: SearchDeps{Paths: p}}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": outside, "pattern": "*.txt"}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "fail-closed") {
		t.Errorf("a root outside the allowlist must be refused fail-closed, got %+v", res)
	}
	if strings.Contains(res.Text, "SECRET") {
		t.Error("the refused search leaked result lines")
	}
}

// Test24R3SearchFilesFailClosedWithoutResolver: no C26 wiring means no walk.
func Test24R3SearchFilesFailClosedWithoutResolver(t *testing.T) {
	res, err := searchFiles{}.Execute(context.Background(),
		r2aParams(t, map[string]any{"path": `C:\x`, "pattern": "*.txt"}), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "fail-closed") {
		t.Errorf("no C26 resolver must fail closed, got %+v", res)
	}
}

// Test24R3SearchFilesNeedsADirectoryRoot stops the "path is a file" shape.
func Test24R3SearchFilesNeedsADirectoryRoot(t *testing.T) {
	root := t.TempDir()
	f := r3Files(t, root, "a.txt")
	p, canon := r2aCanon(t, []string{root}, f)
	res := r3Exec(t, SearchDeps{Paths: p, Classify: r3None}, map[string]any{"path": canon, "pattern": "*.txt"})
	if !res.IsError || !strings.Contains(res.Text, "不是目录") {
		t.Errorf("a file root must be refused as a non-directory, got %+v", res)
	}
}

// Test24R3SearchFilesNeverListsBlacklistedFiles is ticket:25-26 with teeth: the
// blacklisted name never reaches the listing, only an anonymous count.
func Test24R3SearchFilesNeverListsBlacklistedFiles(t *testing.T) {
	root := t.TempDir()
	r3Files(t, root, "SECRET.docx")
	visible := r3Files(t, root, "public.docx")
	p, canon := r2aCanon(t, []string{root}, root)
	d := SearchDeps{Paths: p, Classify: func(path string) risk.Class {
		if strings.Contains(filepath.Base(path), "SECRET") {
			return risk.ClassA
		}
		return risk.ClassNone
	}}
	res := r3Exec(t, d, map[string]any{"path": canon, "pattern": "*.docx"})
	if strings.Contains(res.Text, "SECRET") {
		t.Errorf("a blacklisted name leaked into the listing: %q", res.Text)
	}
	if !strings.Contains(res.Text, "敏感清单排除 1 项") {
		t.Errorf("the anonymous exclusion count is missing: %q", res.Text)
	}
	if !strings.Contains(res.Text, filepath.Base(visible)) {
		t.Errorf("the non-blacklisted hit disappeared: %q", res.Text)
	}
}

// Test24R3SearchFilesBlacklistedRootIsRefused: the veto also applies to the root.
func Test24R3SearchFilesBlacklistedRootIsRefused(t *testing.T) {
	root := t.TempDir()
	p, canon := r2aCanon(t, []string{root}, root)
	d := SearchDeps{Paths: p, Classify: func(string) risk.Class { return risk.ClassB }}
	res := r3Exec(t, d, map[string]any{"path": canon, "pattern": "*"})
	if !res.IsError || !strings.Contains(res.Text, "敏感清单") {
		t.Errorf("a blacklisted root must be refused, got %+v", res)
	}
}

// Test24R3SearchFilesSkipsSymlinks: a symlink is never followed and never
// printed, even when its own name matches.
func Test24R3SearchFilesSkipsSymlinks(t *testing.T) {
	root := t.TempDir()
	target := r3Files(t, root, "target.txt")
	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable on this machine (Windows needs dev mode or admin): %v", err)
	}
	p, canon := r2aCanon(t, []string{root}, root)
	res := r3Exec(t, SearchDeps{Paths: p, Classify: r3None}, map[string]any{"path": canon, "pattern": "*.txt"})
	listing := strings.Join(r3Lines(res.Text), "\n")
	if strings.Contains(listing, "link.txt") {
		t.Errorf("a symlink entry was listed: %q", listing)
	}
	if !strings.Contains(listing, "target.txt") {
		t.Errorf("the real file disappeared: %q", res.Text)
	}
	if !strings.Contains(res.Text, "符号链接跳过 1 项") {
		t.Errorf("the skipped-link count is missing: %q", res.Text)
	}
}

// Test24R3SearchFilesDepthCapBothWays pins the depth budget, with a control that
// proves the deeper file is reachable one level up (so a red cap is not a
// broken walk).
func Test24R3SearchFilesDepthCapBothWays(t *testing.T) {
	root := t.TempDir()
	deep := filepath.Join(root, "l1", "l2", "l3", "deep.md")
	if err := os.MkdirAll(filepath.Dir(deep), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(deep, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	p, canon := r2aCanon(t, []string{root}, root)
	// depth cap 3 stops the walk before l3 is descended.
	capped := r3Exec(t, SearchDeps{Paths: p, MaxDepth: 3, Classify: r3None},
		map[string]any{"path": canon, "pattern": "deep.md"})
	if strings.Contains(capped.Text, "deep.md") {
		t.Errorf("depth cap 3 must not reach the file under l3: %q", capped.Text)
	}
	open := r3Exec(t, SearchDeps{Paths: p, MaxDepth: 8, Classify: r3None},
		map[string]any{"path": canon, "pattern": "deep.md"})
	if !strings.Contains(open.Text, "deep.md") {
		t.Errorf("the control (cap 8) must find it, so the cap is what excluded it: %q", open.Text)
	}
}

// Test24R3SearchFilesEntryCapTruncates is the second of the three budgets: the
// entries-examined ceiling, marked truncated rather than reported as complete.
func Test24R3SearchFilesEntryCapTruncates(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 100; i++ {
		r3Files(t, root, fmt.Sprintf("dir%02d/file.md", i))
	}
	p, canon := r2aCanon(t, []string{root}, root)
	res := r3Exec(t, SearchDeps{Paths: p, MaxResults: 500, MaxFiles: 20, Classify: r3None},
		map[string]any{"path": canon, "pattern": "file.md"})
	if !res.Truncated {
		t.Errorf("hitting the entry cap must mark truncated: %q", res.Text)
	}
	if !strings.Contains(res.Text, "truncated") {
		t.Errorf("the head must say truncated for the reader too: %q", res.Text)
	}
	if got := len(r3Lines(res.Text)); got >= 20 {
		t.Errorf("entry cap 20 must bound the scan, got %d result lines", got)
	}
}

// Test24R3SearchFilesResultCapTruncates is the third budget: the hit ceiling.
func Test24R3SearchFilesResultCapTruncates(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 20; i++ {
		r3Files(t, root, "doc"+strconv.Itoa(i)+".md")
	}
	p, canon := r2aCanon(t, []string{root}, root)
	res := r3Exec(t, SearchDeps{Paths: p, MaxResults: 5, Classify: r3None},
		map[string]any{"path": canon, "pattern": "*.md"})
	lines := r3Lines(res.Text)
	if len(lines) != 5 {
		t.Errorf("want exactly 5 hits at cap 5, got %d:\n%s", len(lines), res.Text)
	}
	if !res.Truncated || !strings.Contains(res.Text, "达到上限，truncated") {
		t.Errorf("hitting the result cap must be reported as truncated: %q", res.Text)
	}
}

// Test24R3SearchFilesCancelMarksPartialResult is ticket:27 for this row: a walk
// cut short reports what it had and says truncated - never a clean scan, never
// a failed call.
func Test24R3SearchFilesCancelMarksPartialResult(t *testing.T) {
	root := t.TempDir()
	r3Files(t, root, "top.md")
	r3Files(t, root, "deeper/inner.md")
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	p, canon := r2aCanon(t, []string{root}, root)
	d := SearchDeps{Paths: p, Classify: func(path string) risk.Class {
		if filepath.Base(path) == "deeper" {
			cancel()
		}
		return risk.ClassNone
	}}
	res, err := searchFiles{d: d}.Execute(ctx, r2aParams(t, map[string]any{"path": canon, "pattern": "*.md"}), nil)
	if err != nil {
		t.Fatalf("a cancelled search must not fail the call: %v", err)
	}
	if !res.Truncated || !strings.Contains(res.Text, "已中止") {
		t.Errorf("a cancelled search must be marked truncated: %q", res.Text)
	}
}

// Test24R3SearchFilesDeclMatchesD34Row pins the D34 row (L0 / fs.read) plus the
// C26 path parameter the bridge needs for R2/R3.
func Test24R3SearchFilesDeclMatchesD34Row(t *testing.T) {
	decl := SearchFilesDecl()
	if len(decl.Capabilities) != 1 || decl.Capabilities[0] != CapFSRead {
		t.Errorf("D34 row says capability fs.read, got %v", decl.Capabilities)
	}
	if len(decl.Needs) != 1 || decl.Needs[0] != CapFSRead {
		t.Errorf("Needs must match Capabilities (a mismatch is a hard reject), got %v", decl.Needs)
	}
	if decl.Declared != risk.L0 {
		t.Errorf("D34 row says L0, got %v", decl.Declared)
	}
	if strings.Join(decl.PathParams, ",") != "path" {
		t.Errorf("only `path` is a path parameter (pattern is a name glob), got %v", decl.PathParams)
	}
	if !decl.Resident || decl.Provider != KindBuiltin {
		t.Errorf("builtin resident row, got %+v", decl)
	}
	if e := BuiltinSearchFilesEntries(SearchDeps{}); len(e) != 1 || e[0].Tool.Name() != "search.files" {
		t.Errorf("entry shape wrong: %+v", e)
	}
}

// Test24R3SearchFilesRegistersInRegistry proves the entry passes C1 validation.
func Test24R3SearchFilesRegistersInRegistry(t *testing.T) {
	r := NewRegistry()
	for _, e := range BuiltinSearchFilesEntries(SearchDeps{}) {
		if err := r.Register(e); err != nil {
			t.Fatalf("Register(search.files): %v", err)
		}
	}
	if _, ok := r.Lookup("search.files"); !ok {
		t.Fatal("search.files not resolvable after registration")
	}
}

// ---- search.apps -----------------------------------------------------------------

// r3AppSource is the AC#4 "fixture-installed app" seam.
type r3AppSource struct {
	label  string
	apps   []string
	aumids []string
	err    error
}

// Label implements AppSource.
func (s r3AppSource) Label() string { return s.label }

// Apps implements AppSource.
func (s r3AppSource) Apps(context.Context) ([]InstalledAppCandidate, error) {
	if s.err != nil {
		return nil, s.err
	}
	out := make([]InstalledAppCandidate, 0, len(s.apps)+len(s.aumids))
	for _, n := range s.apps {
		out = append(out, InstalledAppCandidate{Name: n, Source: s.label})
	}
	for _, a := range s.aumids {
		out = append(out, InstalledAppCandidate{Name: a, AUMID: a, Source: s.label})
	}
	return out, nil
}

func r3ExecApps(t *testing.T, d SearchAppsDeps, args map[string]any) Result {
	t.Helper()
	res, err := searchApps{d: d}.Execute(context.Background(), r2aParams(t, args), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	return res
}

// Test24R3SearchAppsEnumeratesFixtureAndFeedsLaunch is AC#4 in one place: the
// fixture app is enumerated, its listed target is the SAME absolute path the
// launch-side resolver produced, an unresolvable candidate is dropped rather
// than advertised, and nothing was launched.
func Test24R3SearchAppsEnumeratesFixtureAndFeedsLaunch(t *testing.T) {
	_exe := `C:\Fixtures\fixtureapp.exe`
	res := r3ExecApps(t, SearchAppsDeps{
		Sources:  []AppSource{r3AppSource{label: "fixture", apps: []string{"fixtureapp", "ghostapp"}}},
		Resolver: mapResolver{"fixtureapp": filepath.FromSlash(_exe)},
	}, map[string]any{})
	lines := r3Lines(res.Text)
	if len(lines) != 1 {
		t.Fatalf("only the resolvable candidate may be listed:\n%s", res.Text)
	}
	if !strings.HasPrefix(lines[0], "fixtureapp\tpath\t") {
		t.Errorf("want a kind=path line for fixtureapp, got %q", lines[0])
	}
	if !strings.Contains(lines[0], filepath.Base(_exe)) {
		t.Errorf("the listed target must be the resolver's path, got %q", lines[0])
	}
	if strings.Contains(res.Text, "ghostapp") {
		t.Error("an app app.launch could not resolve must not be advertised")
	}
	if !strings.Contains(res.Text, "未能解析 1") {
		t.Errorf("the dropped candidate must be counted: %q", res.Text)
	}
	// Enumeration only (ticket:24). The launch posture strings are produced by
	// startApp; their absence proves this call never handed anything off.
	for _, banned := range []string{"已纳入 Job", "未纳入", "启动载体"} {
		if strings.Contains(res.Text, banned) {
			t.Errorf("search.apps launched something (%q present in the result): %q", banned, res.Text)
		}
	}
	if !strings.Contains(res.Text, "未启动任何应用") {
		t.Errorf("the result must state it did not launch: %q", res.Text)
	}
}

// Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm: a UWP activation id is
// emitted as kind=aumid, which is exactly the shape launcher_windows.go accepts.
func Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("the AUMID shape gate is the Windows launcher")
	}
	aumid := "Microsoft.WindowsNotepad_8wekyb3d8bbwe!App"
	res := r3ExecApps(t, SearchAppsDeps{
		Sources:  []AppSource{r3AppSource{label: "fixture", aumids: []string{aumid}}},
		Resolver: mapResolver{},
	}, map[string]any{"query": "notepad"})
	if !strings.Contains(res.Text, aumid) {
		t.Errorf("the AUMID must be handed over verbatim: %q", res.Text)
	}
	if !strings.Contains(res.Text, "\taumid\t") {
		t.Errorf("want kind=aumid: %q", res.Text)
	}
}

// Test24R3SearchAppsQueryIsLiteralSubstring mirrors search.content's semantics:
// a case-insensitive literal subset, never a regex.
func Test24R3SearchAppsQueryIsLiteralSubstring(t *testing.T) {
	src := r3AppSource{label: "fixture", apps: []string{"Notepad", "Calculator", "NoteJam"}}
	d := SearchAppsDeps{Sources: []AppSource{src}, Resolver: mapResolver{
		"Notepad":    "C:\\a\\notepad.exe",
		"Calculator": "C:\\a\\calc.exe",
		"NoteJam":    "C:\\a\\notejam.exe",
	}}
	filtered := r3ExecApps(t, d, map[string]any{"query": "note"})
	if !strings.Contains(filtered.Text, "Notepad") || !strings.Contains(filtered.Text, "NoteJam") {
		t.Errorf("substring filter (case-insensitive) missed a hit:\n%s", filtered.Text)
	}
	if strings.Contains(filtered.Text, "Calculator") {
		t.Errorf("the filter let a non-match through:\n%s", filtered.Text)
	}
	if filtered.Truncated {
		t.Error("a 2-of-3 subset is not a truncated answer")
	}
	all := r3ExecApps(t, d, map[string]any{})
	if got := len(r3Lines(all.Text)); got != 3 {
		t.Errorf("an empty query enumerates (bounded), got %d lines:\n%s", got, all.Text)
	}
	regex := r3ExecApps(t, d, map[string]any{"query": "note.*"})
	if strings.Contains(regex.Text, "Notepad") {
		t.Errorf("query must NOT be a regex:\n%s", regex.Text)
	}
}

// Test24R3SearchAppsDedupeAcrossSources: the same app from two hives is one line.
func Test24R3SearchAppsDedupeAcrossSources(t *testing.T) {
	res := r3ExecApps(t, SearchAppsDeps{
		Sources: []AppSource{
			r3AppSource{label: "HKLM", apps: []string{"SharedApp"}},
			r3AppSource{label: "HKCU", apps: []string{"sharedapp"}},
		},
		Resolver: mapResolver{"SharedApp": "C:\\a\\shared.exe"},
	}, map[string]any{})
	if got := len(r3Lines(res.Text)); got != 1 {
		t.Errorf("want one deduped line, got %d:\n%s", got, res.Text)
	}
}

// Test24R3SearchAppsResultCapTruncates: the cap is reported, not hidden.
func Test24R3SearchAppsResultCapTruncates(t *testing.T) {
	names := make([]string, 0, 12)
	fixed := mapResolver{}
	for i := 0; i < 12; i++ {
		n := fmt.Sprintf("app%02d", i)
		names = append(names, n)
		fixed[n] = "C:\\a\\" + n + ".exe"
	}
	res := r3ExecApps(t, SearchAppsDeps{
		Sources:    []AppSource{r3AppSource{label: "fixture", apps: names}},
		Resolver:   fixed,
		MaxResults: 5,
	}, map[string]any{})
	if len(r3Lines(res.Text)) != 5 {
		t.Errorf("want 5 lines at cap 5:\n%s", res.Text)
	}
	if !res.Truncated || !strings.Contains(res.Text, "达到上限，truncated") {
		t.Errorf("the app cap must be marked truncated: %q", res.Text)
	}
}

// Test24R3SearchAppsAllSourcesErrorIsNotAnEmptyMachine: total enumeration failure
// is an error, never a report that nothing is installed.
func Test24R3SearchAppsAllSourcesErrorIsNotAnEmptyMachine(t *testing.T) {
	res := r3ExecApps(t, SearchAppsDeps{
		Sources: []AppSource{r3AppSource{label: "broken", err: context.DeadlineExceeded}},
	}, map[string]any{})
	if !res.IsError || !strings.Contains(res.Text, "检索失败") {
		t.Errorf("every source failing must be reported as a failure, got %+v", res)
	}
	if strings.Contains(res.Text, "检索到 0 个应用\n") {
		t.Error("a failed enumeration must not read like an honest zero")
	}
}

// Test24R3SearchAppsRejectsBadParams closes the input side.
func Test24R3SearchAppsRejectsBadParams(t *testing.T) {
	d := SearchAppsDeps{Sources: []AppSource{r3AppSource{label: "fixture", apps: []string{"x"}}}, Resolver: mapResolver{}}
	res, err := searchApps{d: d}.Execute(context.Background(), json.RawMessage(`{"query":`), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "参数解析失败") {
		t.Errorf("malformed params must be an error result, got %+v", res)
	}
}

// Test24R3SearchAppsDeclIsSysinfoL0NoShell pins the D34 row (L0 / sysinfo /
// "只检索不启动") and the structural part of "never launches": no shell
// capability is declared, so a launch call inside this tool could not be wired.
func Test24R3SearchAppsDeclIsSysinfoL0NoShell(t *testing.T) {
	decl := SearchAppsDecl()
	if len(decl.Capabilities) != 1 || decl.Capabilities[0] != CapSysinfo {
		t.Errorf("D34 row says capability sysinfo, got %v", decl.Capabilities)
	}
	if len(decl.Needs) != 1 || decl.Needs[0] != CapSysinfo {
		t.Errorf("Needs must match Capabilities, got %v", decl.Needs)
	}
	for _, c := range decl.Capabilities {
		if c == CapShell {
			t.Fatal("search.apps must not declare shell (it enumerates, never launches)")
		}
	}
	if decl.Declared != risk.L0 {
		t.Errorf("D34 row says L0, got %v", decl.Declared)
	}
	if len(decl.PathParams) != 0 {
		t.Errorf("no user-supplied path parameter, got %v", decl.PathParams)
	}
	if e := BuiltinSearchAppsEntries(SearchAppsDeps{}); len(e) != 1 || e[0].Tool.Name() != "search.apps" {
		t.Errorf("entry shape wrong: %+v", e)
	}
}

// Test24R3SearchAppsRegistersInRegistry proves C1 validation accepts both rows.
func Test24R3SearchAppsRegistersInRegistry(t *testing.T) {
	r := NewRegistry()
	for _, e := range BuiltinSearchAppsEntries(SearchAppsDeps{}) {
		if err := r.Register(e); err != nil {
			t.Fatalf("Register(search.apps): %v", err)
		}
	}
	if _, ok := r.Lookup("search.apps"); !ok {
		t.Fatal("search.apps not resolvable after registration")
	}
}

// Test24R3DefaultAppSourcesHasAtLeastThePATHFloor proves the un-injected path is
// never an empty enumeration set (an empty set would silently list no apps).
func Test24R3DefaultAppSourcesHasAtLeastThePATHFloor(t *testing.T) {
	srcs := DefaultAppSources()
	if len(srcs) == 0 {
		t.Fatal("no default enumeration source")
	}
	var sawPATH bool
	for _, s := range srcs {
		if s.Label() == "PATH 目录" {
			sawPATH = true
		}
		if s.Label() == "" {
			t.Error("an enumeration source must be labelled")
		}
	}
	if !sawPATH {
		t.Error("the PATH floor is always present")
	}
}
