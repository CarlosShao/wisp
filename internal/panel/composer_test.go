package panel

// Ticket 92 AC#1 and AC#4 and AC#7: the composer's permission surface.
//
// AC#1 asks for a POSITIVE nail, not a comment: plant a "change the mode from
// the composer" call into a composer component and show an instrument goes red.
// Two instruments are asserted here, because the two shapes of that mistake are
// caught by different doors:
//
//	door 1 (ban #6, D33/F2, PLAN.md:1588)   `approval.decide` in frontend/ is
//	   banned outright. A composer that reaches for the decision verb - which is
//	   exactly what "set the mode from the page" degenerates into - goes red on
//	   the existing scanner rule. This test plants that call into a fixture tree
//	   and asserts the rule fires, then asserts the real frontend/src is clean.
//	door 2 (this ticket's gate)             `approval.decide` is not the only
//	   spelling of a mode write. A renderer could post `panel.mode.set` or call
//	   setMode(...) and never touch the banned token, so ban #6 alone would let
//	   it through - which is the argument for hanging the rest of the shape on a
//	   second door rather than pretending one regex covers a permission input.
//
// AC#7 is the negative criterion owner asked for: the git branch/repo switcher
// was cut from the product, so "nothing in the panel surface can switch a
// repository" is checked here as a fact rather than remembered as a decision.

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// ----------------------------------------------------------------- contract

// TestComposerContractTypesMatchFrontend reconciles the composer view models
// with the TypeScript interfaces the page reads, in both directions - the same
// promise ticket 77 made for the approval card, extended to what this ticket
// adds. A renamed field on one side must go red here, not render "undefined"
// inside a permission indicator.
func TestComposerContractTypesMatchFrontend(t *testing.T) {
	root := panelRepoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, "frontend", "src", "lib", "panel.ts"))
	if err != nil {
		t.Fatalf("read frontend/src/lib/panel.ts: %v", err)
	}
	text := string(data)
	pairs := []struct {
		name    string
		goType  any
		tsIface string
	}{
		{"Snapshot", Snapshot{}, "PanelSnapshot"},
		{"ComposerState", ComposerState{}, "ComposerState"},
		{"ModeView", ModeView{}, "ComposerMode"},
		{"WorkspaceView", WorkspaceView{}, "ComposerWorkspace"},
		{"AttachmentRef", AttachmentRef{}, "ComposerAttachment"},
		{"ResultChunk", ResultChunk{}, "ResultChunkView"},
	}
	for _, p := range pairs {
		goKeys := jsonKeysOf(p.goType)
		tsKeys, ok := tsInterfaceKeys(text, p.tsIface)
		if !ok {
			t.Fatalf("frontend/src/lib/panel.ts declares no interface %s - the composer contract was renamed on one side only", p.tsIface)
		}
		if missing := subtract(goKeys, tsKeys); len(missing) > 0 {
			t.Errorf("Go %s emits %v that interface %s does not declare", p.name, missing, p.tsIface)
		}
		if extra := subtract(tsKeys, goKeys); len(extra) > 0 {
			t.Errorf("interface %s reads %v that Go %s never sends - those fields render as undefined", p.tsIface, extra, p.name)
		}
		t.Logf("%s <-> %s: %d JSON keys reconciled", p.name, p.tsIface, len(goKeys))
	}
}

// ------------------------------------------------- AC#1: mode is display-only

// composerModeWriteRe is door 2: every spelling of "the page changed the mode"
// this repository would recognise. ban #6's regex does not cover these, which is
// why they need their own gate (see the file header).
var composerModeWriteRe = regexp.MustCompile(
	`\bsetMode\b|\bmode\.set\b|\bpanel\.mode\.set\b|\bpermission\.mode\.set\b|` +
		`\bperm\.Store\.Set\b|\bPermissionMode\s*=|\bmode\s*=\s*["']auto_approve["']`)

// TestPlantedComposerModeWriteGoesRed is AC#1's positive nail. The fixture is a
// composer component written the WRONG way on purpose, and the assertion is that
// the instruments catch it - not that the real tree happens to be quiet.
func TestPlantedComposerModeWriteGoesRed(t *testing.T) {
	dir := t.TempDir()
	// Plant 1: the banned decision verb used to change the composer's own mode.
	// This is the shape ban #6 exists for (D33/F2: allow decisions are
	// native-side only, and a mode change IS an allow decision with no bound).
	plantDecision := `export function ComposerChip({ onChange }: { onChange: (m: string) => void }) {
  const bridge = window.wispBridge
  return (
    <button onClick={() => bridge?.postMessage(JSON.stringify({ method: "approval.decide", target: "mode", to: "auto_approve" }))}>
      mode
    </button>
  )
}
`
	// Plant 2: a mode write that never mentions the banned token.
	plantSetter := `export function setMode(next: string): void {
  post({ method: "panel.mode.set", to: next })
}
`
	cases := []struct {
		name    string
		content string
		want    *regexp.Regexp
	}{
		{"approval.decide planted in composer trips ban #6", plantDecision, panelDecisionIdentifierRe},
		{"panel.mode.set planted in composer trips the ticket 92 gate", plantSetter, composerModeWriteRe},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(dir, "composer.tsx")
			if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
				t.Fatalf("write plant: %v", err)
			}
			hits := reMatchesInFile(tc.want, path)
			if len(hits) == 0 {
				t.Fatalf("planted a composer that changes the permission mode from the renderer and NOTHING went red: pattern %q matched no line in the plant - the gate is decorative", tc.want)
			}
			t.Logf("red as required: %d line(s) matched, first=%s", len(hits), hits[0])
		})
	}

	// And the real tree must be clean under BOTH doors.
	root := panelRepoRoot(t)
	violations := scanComposerPermissionWrites(t, filepath.Join(root, "frontend"))
	if len(violations) > 0 {
		t.Errorf("frontend/ holds a composer-side permission write; only the native side may decide (PLAN.md:1588):%s", strings.Join(violations, "\n"))
	}
}

// TestModeViewHasNoWriteSurface is the Go-side half of AC#1: the view model the
// panel renders cannot be written through, and the panel package exposes no
// function that sets a mode. (The real setter, perm.Store.Set, lives in
// internal/perm and costs the R20/M4 L2 confirmation - there is no second one.)
func TestModeViewHasNoWriteSurface(t *testing.T) {
	mt := reflect.TypeOf(ModeView{})
	for i := 0; i < mt.NumMethod(); i++ {
		t.Errorf("ModeView has method %q - a mode view must carry no behaviour at all", mt.Method(i).Name)
	}
	src, err := os.ReadFile(filepath.Join(panelRepoRoot(t), "internal", "panel", "composer.go"))
	if err != nil {
		t.Fatalf("read composer.go: %v", err)
	}
	if hits := modeSettersIn(string(src)); len(hits) > 0 {
		t.Errorf("internal/panel exposes a mode setter (%v) - the panel may only display and request", hits)
	}
	if !strings.Contains(string(src), "ModeRequest") {
		t.Fatal("composer.go declares no ModeRequest: AC#1's shape is display + request, and a request type that does not exist becomes a direct write")
	}
}

// TestModeRequestResolvesThroughRisksOwnParser pins that a request cannot smuggle
// an unknown spelling into a default mode.
func TestModeRequestResolvesThroughRisksOwnParser(t *testing.T) {
	for _, s := range risk.ModeNames() {
		m, err := ModeRequest{To: s}.Parse()
		if err != nil {
			t.Errorf("ModeRequest{To:%q}.Parse() = %v, want the mode risk itself names", s, err)
			continue
		}
		if m.String() != s {
			t.Errorf("Parse(%q).String() = %q", s, m.String())
		}
	}
	if _, err := (ModeRequest{To: "yolo"}).Parse(); err == nil {
		t.Error("ModeRequest{To:\"yolo\"} parsed without error - an unknown mode must be refused, not defaulted")
	}
}

// --------------------------------------------------- AC#4: state survives close

func TestComposerStateSurvivesPanelCloseAndReopen(t *testing.T) {
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	state := NewComposerState(risk.ModeAskHighRisk, WorkspaceView{
		Set: true, Spelling: `~/notes`, Canonical: `D:\work\notes`, Rewritten: true,
	}, []AttachmentRef{{ID: "att-1", Name: "shot.png", MIME: "image/png", Kind: "image",
		SizeBytes: 1234, Artifact: "attachment-0011223344556677.png", Stored: true}}, MaxAttachmentBytes)
	if state.Workspace.Rewritten != true {
		t.Fatal("fixture premise broken")
	}
	first := NewSnapshot(nil, nil, state, now)
	if first.Composer.Mode.Current != risk.ModeAskHighRisk.String() {
		t.Fatalf("snapshot mode = %q, want %q", first.Composer.Mode.Current, risk.ModeAskHighRisk)
	}

	wire, err := json.Marshal(first)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// The panel closes (ball awake, panel never expanded) and reopens: every
	// value the page held is gone, and one push from Go restores it.
	var after Snapshot
	if err := json.Unmarshal(wire, &after); err != nil {
		t.Fatalf("unmarshal after reopen: %v", err)
	}
	if !reflect.DeepEqual(first, after) {
		t.Errorf("state after reopening the panel differs from the state before it:\n before %+v\n after  %+v", first, after)
	}
	if after.Composer.Workspace.Canonical != `D:\work\notes` || !after.Composer.Workspace.Rewritten {
		t.Errorf("workspace did not survive the reopen: %+v - AC#3's audit and AC#4's recovery read the same account", after.Composer.Workspace)
	}
	t.Logf("composer state round-tripped through %d bytes: mode=%q workspace=%q attachments=%d",
		len(wire), after.Composer.Mode.Current, after.Composer.Workspace.Canonical, len(after.Composer.Attachments))
}

// TestUnknownModeNeverRendersAsASafeOne is the fail-closed half of AC#4: when
// native cannot read the mode, the panel must show that, not the default.
func TestUnknownModeNeverRendersAsASafeOne(t *testing.T) {
	s := NewSnapshot(nil, nil, ComposerState{Workspace: WorkspaceView{
		Reason: "未选择工作区"},
	}, time.Now())
	if s.Composer.Mode.Current != "unknown" {
		t.Errorf("an unreadable mode rendered as %q, want \"unknown\"", s.Composer.Mode.Current)
	}
	if s.Composer.Mode.Current == risk.DefaultMode().String() {
		t.Error("the unknown mode string equals the default mode: a failed read would look like a chosen档")
	}
	if len(s.Composer.AttachmentMIMEs) != len(AcceptedMIMETypes()) {
		t.Errorf("a default snapshot hides the accepted attachment types: %+v", s.Composer.AttachmentMIMEs)
	}
}

// ------------------------------------------------- AC#7: no git switching, ever

// gitSwitchCapabilityRe names the cut feature's entry points. Owner cut git
// branch/repository switching from the product (R20/M5), so this is a checkable
// negative rather than a sentence in a ticket.
var gitSwitchCapabilityRe = regexp.MustCompile(
	`\bgit\s+checkout\b|\bgit\s+switch\b|\bswitchBranch\b|\bcheckoutBranch\b|\bchangeRepo(?:sitory)?\b|` +
		`\brepoPicker\b|\bbranchSelect(or)?\b|\bworktree\b|\bgit\.branch\b|\bgit\.repo\b|\bvcs\.switch\b`)

// TestNoGitSwitchCapabilityInThePanelSurface greps the two trees this ticket
// owns. Production files only (see the header of this test file for why _test.go
// is skipped: the patterns live in this file, and a self-match would make the
// instrument unreadable rather than wrong).
func TestNoGitSwitchCapabilityInThePanelSurface(t *testing.T) {
	root := panelRepoRoot(t)
	var violations []string
	scan := func(dir string, accept func(string) bool) {
		err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				switch d.Name() {
				case "node_modules", ".git", "dist", "fixtures":
					return fs.SkipDir
				}
				return nil
			}
			if !accept(p) {
				return nil
			}
			if hits := reMatchesInFile(gitSwitchCapabilityRe, p); len(hits) > 0 {
				rel, _ := filepath.Rel(root, p)
				violations = append(violations, rel+": "+hits[0])
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}
	scan(filepath.Join(root, "frontend"), func(p string) bool {
		if strings.HasSuffix(p, "_test.go") || strings.HasSuffix(p, ".test.ts") || strings.HasSuffix(p, ".test.tsx") {
			return false
		}
		return strings.HasSuffix(p, ".ts") || strings.HasSuffix(p, ".tsx") ||
			strings.HasSuffix(p, ".css") || strings.HasSuffix(p, ".html")
	})
	scan(filepath.Join(root, "internal", "panel"), func(p string) bool {
		return strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go")
	})
	if len(violations) > 0 {
		t.Errorf("a git branch/repository switcher exists in the panel surface; owner cut it from the product (票 92 AC#7):%s",
			strings.Join(violations, "\n"))
	}
	t.Logf("AC#7 negative criterion held: no git switch entry point in frontend/ or internal/panel/ (production files)")
}

// --------------------------------------------------------------- walk helpers

// modeSettersIn returns the lines of a Go source text that declare a mode
// setter in the panel package - the shape AC#1 forbids on either side of the
// bridge.
var modeSetterRe = regexp.MustCompile(`func\s+(Set|Write|Apply|Change)Mode\b|\bPermissionMode\s*=|perm\.Store\.Set\(`)

func modeSettersIn(src string) []string {
	var out []string
	for i, line := range strings.Split(src, "\n") {
		if modeSetterRe.MatchString(line) {
			out = append(out, "line "+strconv.Itoa(i+1)+": "+strings.TrimSpace(line))
		}
	}
	return out
}

func reMatchesInFile(re *regexp.Regexp, path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{"unreadable: " + err.Error()}
	}
	var out []string
	for i, line := range strings.Split(string(data), "\n") {
		if re.MatchString(line) {
			out = append(out, filepath.Base(path)+":"+strconv.Itoa(i+1)+": "+strings.TrimSpace(line))
		}
	}
	return out
}

// scanComposerPermissionWrites applies BOTH doors to one tree and returns the
// offending lines (empty = clean).
func scanComposerPermissionWrites(t *testing.T, root string) []string {
	t.Helper()
	var out []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".git", "dist":
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".ts") && !strings.HasSuffix(p, ".tsx") && !strings.HasSuffix(p, ".css") {
			return nil
		}
		for _, re := range []*regexp.Regexp{panelDecisionIdentifierRe, composerModeWriteRe} {
			if hits := reMatchesInFile(re, p); len(hits) > 0 {
				out = append(out, hits...)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(out)
	return out
}
