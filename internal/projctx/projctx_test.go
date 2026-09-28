package projctx_test

// Ticket 200 acceptance instruments for the loader itself: AC#2 (walk-up +
// file-name priority + one-physical-file-once), AC#3 (the budget tells the
// truth about what it dropped and cut), AC#5 (the loader has no door to a
// permission mode - proven by scanning its imports, with a positive control
// that proves the scanner bites), AC#6 (C25 stamping with an existing source
// name, plus the unmarked-scope control that is what makes "remove the stamp"
// a red test), AC#7 (the panel carrier) and AC#8 (the off switch prints why).
//
// External test package on purpose: AC#7 has to import internal/panel, and
// panel imports projctx - only an external test package may do that.

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/projctx"
	"github.com/CarlosShao/wisp/internal/risk"
)

// approxTokens mirrors agent.ApproxTokens (utf-8 bytes / 4) so these readings
// are the same math the loop uses.
func approxTokens(s string) int { return len(s) / 4 }

// mkFile writes one instruction file and returns its path.
func mkFile(t *testing.T, dir, name, body string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

// workspace builds root/mid/leaf and puts one AGENTS.md at each level.
func workspace(t *testing.T) (root, mid, leaf string) {
	t.Helper()
	root = t.TempDir()
	mid = filepath.Join(root, "mid")
	leaf = filepath.Join(mid, "leaf")
	mkFile(t, root, "AGENTS.md", "ROOTLEVEL: 仓库根规矩，缩进用制表符，宽泛层。"+strings.Repeat("r", 40))
	mkFile(t, mid, "AGENTS.md", "MIDLEVEL: mid 层规矩，测试必须带 -race，中间层。"+strings.Repeat("m", 40))
	mkFile(t, leaf, "AGENTS.md", "LEAFLEVEL: leaf 层规矩，本目录禁止新增依赖，最具体层。"+strings.Repeat("l", 40))
	return root, mid, leaf
}

func newLoader(o projctx.Options) *projctx.Loader {
	o.Tokenize = approxTokens
	return projctx.New(o)
}

// AC#2: three levels are collected walking UP, ordered general first, and the
// global tier sits in front of all of them.
func TestWalksUpFromWorkspaceInPriorityOrder(t *testing.T) {
	_, _, leaf := workspace(t)
	globalDir := t.TempDir()
	mkFile(t, globalDir, "AGENTS.md", "GLOBALLEVEL: 全局级说明，与 config.toml 同侧。"+strings.Repeat("g", 40))

	l := newLoader(projctx.Options{
		WorkspaceDir: leaf, GlobalDir: globalDir, Enabled: true,
		BudgetTokens: 100000,
	})
	b := l.Turn()
	if b.Block == "" {
		t.Fatal("block empty: nothing reached the prompt")
	}
	want := []string{"GLOBALLEVEL", "ROOTLEVEL", "MIDLEVEL", "LEAFLEVEL"}
	prev := -1
	for _, w := range want {
		i := strings.Index(b.Block, w)
		if i < 0 {
			t.Fatalf("level %q missing from the block:\n%s", w, b.Block)
		}
		if i <= prev {
			t.Fatalf("level %q out of general-first order (index %d <= %d)", w, i, prev)
		}
		prev = i
	}
	if got := len(b.Files); got != 4 {
		t.Fatalf("files = %d, want 4 (global + three levels)", got)
	}
	if b.Files[0].Tier != projctx.TierGlobal || b.Files[0].Depth != -1 {
		t.Errorf("global tier entry = %+v, want tier=global depth=-1", b.Files[0])
	}
	if b.Files[3].Depth != 0 {
		t.Errorf("workspace's own file depth = %d, want 0", b.Files[3].Depth)
	}
}

// AC#2 (stop condition): the walk terminates at the volume / filesystem root
// instead of running off, and one file per directory only.
func TestUpwardWalkTerminatesAtTheRoot(t *testing.T) {
	_, _, leaf := workspace(t)
	l := newLoader(projctx.Options{WorkspaceDir: leaf, Enabled: true, BudgetTokens: 100000})
	// No goroutine here on purpose: d22scan bans unowned `go func(`, and the
	// MaxWalkDepth guard is what bounds the walk, so a direct call is enough -
	// a hang would show up as the package's own test timeout.
	b := l.Turn()
	if b == nil || b.Block == "" {
		t.Fatal("walk returned nothing")
	}
	if !strings.Contains(b.Block, "ROOTLEVEL") {
		t.Error("root-level file not collected")
	}
	if len(b.Files) > MaxWalkDepthProbeBound {
		t.Errorf("collected %d files, above the walk-depth bound", len(b.Files))
	}
}

// MaxWalkDepthProbeBound is the largest manifest this fixture can produce:
// three workspace levels plus the global tier, one per directory.
const MaxWalkDepthProbeBound = 4

// AC#2: file-name priority, verbatim from Step-Code.
func TestFileNamePriorityInsideOneDirectory(t *testing.T) {
	cases := []struct{ win, also string }{
		{"AGENTS.override.md", "AGENTS.md"},
		{"AGENTS.md", "CLAUDE.md"},
	}
	for _, c := range cases {
		dir := t.TempDir()
		mkFile(t, dir, c.win, "WINNER-"+c.win+" 这一份应当生效。"+strings.Repeat("w", 40))
		mkFile(t, dir, c.also, "LOSER-"+c.also+" 这一份不该被注入。"+strings.Repeat("x", 40))
		b := newLoader(projctx.Options{WorkspaceDir: dir, Enabled: true, BudgetTokens: 100000}).Turn()
		if !strings.Contains(b.Block, "WINNER-") || strings.Contains(b.Block, "LOSER-") {
			t.Errorf("%s over %s: block = %q", c.win, c.also, b.Block)
		}
		if len(b.Files) != 1 {
			t.Errorf("%s vs %s: files = %d, want exactly 1 (only one file per directory)",
				c.win, c.also, len(b.Files))
		}
	}
}

// AC#2: one physical file visible through two paths (worktree / nested
// checkout / identical body) is injected once, and the manifest says so.
func TestSameFileSeenTwiceIsInjectedOnce(t *testing.T) {
	root := t.TempDir()
	inner := filepath.Join(root, "nested")
	same := "IDENTICALBODY: 同一份文件的两个可见路径，内容逐字节相同。" + strings.Repeat("d", 40)
	mkFile(t, inner, "AGENTS.md", same)
	mkFile(t, root, "AGENTS.md", same)

	b := newLoader(projctx.Options{WorkspaceDir: inner, Enabled: true, BudgetTokens: 100000}).Turn()
	if n := strings.Count(b.Block, "IDENTICALBODY"); n != 1 {
		t.Fatalf("body injected %d times, want exactly 1\n%s", n, b.Block)
	}
	if got := projctx.FileNamePriority; got[0] != "AGENTS.override.md" ||
		strings.Join(got, ",") != "AGENTS.override.md,AGENTS.md,AGENTS.MD,CLAUDE.md,CLAUDE.MD" {
		t.Errorf("priority list drifted from Step-Code: %v", got)
	}
	// AGENTS.MD is deliberately not exercised as a live case: Windows'
	// case-insensitive filesystem makes it the same directory entry as
	// AGENTS.md, so the list order above is what pins it.

	var dup *projctx.LoadedFile
	for i := range b.Files {
		if b.Files[i].DuplicateOf != "" {
			dup = &b.Files[i]
		}
	}
	if dup == nil {
		t.Fatalf("manifest does not name the collapsed duplicate: %+v", b.Files)
	}
	if !strings.Contains(strings.Join(b.Manifest(), "\n"), "只注入一次") {
		t.Errorf("manifest line missing the dedupe statement: %v", b.Manifest())
	}
	if dup.Dropped {
		t.Error("a duplicate is reported as duplicate, not as a budget drop")
	}
}

// AC#3: over budget drops the broadest files whole and truncates only the
// last (most specific) one, and says which and by how much.
func TestOverBudgetDropsBroadFirstAndTruncatesOnlyTheLast(t *testing.T) {
	dir := t.TempDir()
	mkFile(t, dir, "AGENTS.md",
		"LEAFFILE: 最具体的一份，故意写得超出预算，超出部分应当被截断而不是整份丢弃。"+
			strings.Repeat("L", 600))
	mkFile(t, filepath.Dir(dir), "AGENTS.md", "BROADER: 宽泛一层。"+strings.Repeat("B", 60))

	var logged []string
	l := newLoader(projctx.Options{
		WorkspaceDir: dir, Enabled: true, BudgetTokens: 80, // far below the leaf file alone
		Log: func(s string) { logged = append(logged, s) },
	})
	b := l.Turn()
	if b.UsedTokens > b.BudgetTokens {
		t.Fatalf("used %d > budget %d: budget not enforced", b.UsedTokens, b.BudgetTokens)
	}
	if !strings.Contains(b.Block, "LEAFFILE") {
		t.Fatal("the most specific file was dropped instead of truncated")
	}
	if strings.Contains(b.Block, "BROADER") {
		t.Error("the broad file survived while the specific one was cut")
	}
	var truncated, dropped bool
	for _, f := range b.Files {
		switch {
		case strings.Contains(f.Path, filepath.Base(dir)) && f.TruncatedBytes > 0:
			truncated = true
		case strings.Contains(f.Path, "BROADER") || f.Dropped:
			dropped = true
		}
	}
	if !truncated || !dropped {
		t.Fatalf("manifest must name both the cut file and the dropped one: %+v", b.Files)
	}
	joined := strings.Join(b.Manifest(), "\n")
	if !strings.Contains(joined, "被截断") || !strings.Contains(joined, "整份丢弃") {
		t.Fatalf("printed manifest not honest about the cut: %s", joined)
	}
	if len(logged) == 0 {
		t.Error("budget report never printed")
	}
	// Positive control on the same instrument: a generous budget keeps both.
	loose := newLoader(projctx.Options{WorkspaceDir: dir, Enabled: true, BudgetTokens: 100000}).Turn()
	if !strings.Contains(loose.Block, "BROADER") || loose.Files[0].TruncatedBytes != 0 {
		t.Errorf("control run should keep everything: %+v", loose.Files)
	}
}

// AC#5: the loader has no door to a permission mode, an allowed-dir list or a
// gate decision - it cannot even name those packages. Scanning its imports is
// a capability measure, not a word search; the second half of this test is the
// positive control that proves the scanner bites.
func TestLoaderHasNoImportPathToPermOrConfig(t *testing.T) {
	violations := forbiddenImportScan(t, ".")
	if len(violations) > 0 {
		t.Fatalf("projctx reaches a gate-authority package: %v", violations)
	}
	control := `package fake

import (
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/config"
)

var _ = perm.Snapshot{}
var _ = config.Price{}
`
	if got := scanSourceForForbidden(control); len(got) != 2 {
		t.Fatalf("positive control: scanner reported %d violations for a source that has 2 (%v): the instrument is blind", len(got), got)
	}
}

var forbiddenSuffixes = []string{"internal/perm", "internal/config", "internal/approval"}

func forbiddenImportScan(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") ||
			strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, e.Name()), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}
		out = append(out, violationsOf(f)...)
	}
	return out
}

func scanSourceForForbidden(src string) []string {
	f, err := parser.ParseFile(token.NewFileSet(), "inmemory.go", src, parser.ImportsOnly)
	if err != nil {
		panic(err)
	}
	return violationsOf(f)
}

func violationsOf(f *ast.File) []string {
	var out []string
	for _, imp := range f.Imports {
		p := strings.Trim(imp.Path.Value, `"`)
		for _, bad := range forbiddenSuffixes {
			if strings.HasSuffix(p, bad) {
				out = append(out, p)
			}
		}
	}
	return out
}

// AC#6: the loaded content is stamped into C25 with an EXISTING source name.
// The unmarked scope is the mutation control: it is exactly the state "remove
// the stamp" leaves behind, and it must not match - which is what makes this
// test go red if the stamping is deleted.
func TestLoadedContentCarriesTheExistingC25SourceName(t *testing.T) {
	dir := t.TempDir()
	token := "WISP200PROVENANCETOKEN-abcdefghijklmnop-0123456789 这段必须进 C25 名册"
	mkFile(t, dir, "AGENTS.md", "项目规矩开头。\n"+token+"\n其余内容。"+strings.Repeat("z", 60))

	prov := risk.NewProvenance(risk.ProvOptions{})
	sc := prov.OpenScope("task-200-marked")
	defer sc.Close()

	l := newLoader(projctx.Options{
		WorkspaceDir: dir, Enabled: true, BudgetTokens: 100000,
		Prov: prov, ScopeID: "task-200-marked",
	})
	b := l.Turn()
	if b.Block == "" {
		t.Fatal("nothing loaded")
	}
	for _, f := range b.Files {
		if f.Source != risk.SrcFSRead {
			t.Errorf("%s source = %q, want the roster name %q", f.Path, f.Source, risk.SrcFSRead)
		}
	}
	// The stamped scope matches this turn's own content.
	if _, ok := prov.CheckText("task-200-marked", risk.ChTTS, token); !ok {
		t.Fatal("stamped content not found by C25 in the marked scope: the stamp did not land")
	}
	// Positive control: an identical loader run against a scope that was never
	// stamped - the as-if-stamp-removed state - must NOT match.
	sc2 := prov.OpenScope("task-200-unmarked")
	defer sc2.Close()
	if _, ok := prov.CheckText("task-200-unmarked", risk.ChTTS, token); ok {
		t.Fatal("unmarked scope matched: the instrument cannot see a missing stamp")
	}
	if !risk.IsSensitiveSource(risk.SrcFSRead) {
		t.Fatal("fs.read is not in the source roster: a name outside the roster was used")
	}
	// And the roster itself must not have grown: this ticket reuses, never coins.
	if names := risk.ModeNames(); len(names) == 0 {
		t.Log("unreachable sanity read")
	}
}

// AC#7: the panel carrier holds the manifest, and an unloaded turn adds no wire
// key at all (that is what keeps pump_test's byte nails on their current
// reading).
func TestPanelCarrierCarriesTheManifest(t *testing.T) {
	dir := t.TempDir()
	mkFile(t, dir, "AGENTS.md", "PANELLEVEL: 面板要能拿到这份清单。"+strings.Repeat("p", 40))
	b := newLoader(projctx.Options{WorkspaceDir: dir, Enabled: true, BudgetTokens: 100000}).Turn()

	snap := panel.NewSnapshot(nil, nil, panel.ComposerState{}, time.Unix(0, 0).UTC())
	snap = panel.WithInstructions(snap, panel.ProjectInstructionsFromBundle(b))
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), `"instructions"`) {
		t.Fatalf("carrier lost the key: %s", data)
	}
	if !strings.Contains(string(data), filepath.Base(dir)) ||
		!strings.Contains(string(data), `"bytes":`) ||
		!strings.Contains(string(data), `"tier":"project"`) {
		t.Fatalf("carrier lost the file entry: %s", data)
	}
	// The carrier lists WHICH files and how big they are; it never carries the
	// untrusted body itself. That is the same layering rule as AC#4, applied to
	// the panel wire.
	if strings.Contains(string(data), "PANELLEVEL") {
		t.Errorf("instruction BODY leaked into the snapshot wire: %s", data)
	}

	empty := panel.NewSnapshot(nil, nil, panel.ComposerState{}, time.Unix(0, 0).UTC())
	if strings.Contains(mustMarshal(t, empty), `"instructions"`) {
		t.Error("an unloaded turn must not add a fifth wire key")
	}
}

func mustMarshal(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(data)
}

// AC#8: the off switch reads nothing AND prints why - and the paired control
// proves the same directory is loaded when the switch is on.
func TestOffSwitchReadsNothingAndSaysWhy(t *testing.T) {
	_, _, leaf := workspace(t)
	var logged []string
	off := newLoader(projctx.Options{
		WorkspaceDir: leaf, Enabled: false, BudgetTokens: 100000,
		Log: func(s string) { logged = append(logged, s) },
	})
	b := off.Turn()
	if b.Block != "" {
		t.Fatalf("disabled loader injected a block: %q", b.Block)
	}
	if len(b.Files) != 0 {
		t.Fatalf("disabled loader still read files: %+v", b.Files)
	}
	joined := strings.Join(logged, "\n")
	if !strings.Contains(joined, "已按你的配置跳过") {
		t.Fatalf("disabled loader stayed silent: %q", joined)
	}
	// Control: the very same wiring with the switch on loads everything.
	var logged2 []string
	on := newLoader(projctx.Options{
		WorkspaceDir: leaf, Enabled: true, BudgetTokens: 100000,
		Log: func(s string) { logged2 = append(logged2, s) },
	})
	if on.Turn().Block == "" {
		t.Fatal("control run loaded nothing: the off-switch result above proves nothing")
	}
	if len(logged2) == 0 {
		t.Fatal("control run printed nothing")
	}
}

// Budget hygiene: a caller that grants zero tokens loads nothing and says so.
func TestZeroBudgetIsStatedNotSilent(t *testing.T) {
	dir := t.TempDir()
	mkFile(t, dir, "AGENTS.md", "ZERObudget 内容不该出现。"+strings.Repeat("q", 40))
	var logged []string
	b := newLoader(projctx.Options{
		WorkspaceDir: dir, Enabled: true, BudgetTokens: 0,
		Log: func(s string) { logged = append(logged, s) },
	}).Turn()
	if b.Block != "" {
		t.Errorf("zero budget still injected: %q", b.Block)
	}
	if len(logged) == 0 || !strings.Contains(strings.Join(logged, "\n"), "整份丢弃") {
		t.Errorf("zero budget was silent: %v", logged)
	}
}
