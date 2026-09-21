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

	// Door 2 forgives prose and nothing else: a mode write wearing a comment's
	// clothes must still be caught, or the comment filter would be a hole in the
	// gate rather than a filter on it. (This file's own author learned the rule
	// the cheap way: a header line saying "there is no setMode-shaped call here"
	// was the first real-tree hit door 2 ever produced.)
	t.Run("door 2 ignores comment prose but not code behind a comment opener", func(t *testing.T) {
		prose := filepath.Join(dir, "prose.tsx")
		proseBody := "/**\n * 不要在这里 " + "setMode" + " - 档位由原生侧决定\n */\nexport const ok = 1 // " + "setMode" + "\n"
		if err := os.WriteFile(prose, []byte(proseBody), 0o600); err != nil {
			t.Fatalf("write prose plant: %v", err)
		}
		if hits := modeWriteMatchesInCodeFile(prose); len(hits) > 0 {
			t.Errorf("comment prose tripped the gate (%v) - the filter is missing the prose it claims to allow", hits)
		}
		hidden := filepath.Join(dir, "hidden.tsx")
		if err := os.WriteFile(hidden, []byte("/* just a note */ set"+"Mode(\"auto_approve\")\n"), 0o600); err != nil {
			t.Fatalf("write hidden plant: %v", err)
		}
		if hits := modeWriteMatchesInCodeFile(hidden); len(hits) == 0 {
			t.Error("a mode write placed after a comment opener escaped door 2 - the comment filter is a hole")
		}
	})
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
	}, []AttachmentRef{{
		ID: "att-1", Name: "shot.png", MIME: "image/png", Kind: "image",
		SizeBytes: 1234, Artifact: "attachment-0011223344556677.png", Stored: true,
	}}, MaxAttachmentBytes)
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
	s := NewSnapshot(nil, nil, ComposerState{
		Workspace: WorkspaceView{
			Reason: "未选择工作区",
		},
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
		// AC#7's file set is deliberately NOT widened to .js/.mjs: the acceptance
		// round's R-92-1 is about the mode-write door, and measuring this very
		// predicate with rendererSourceFile(p) went red on
		// frontend/scripts/vendor-shadcn.mjs:2, a build script whose comment only
		// mentions "git checkout". Capability entry points are pinned structurally
		// by TestTheRendererHoldsExactlyOneDoorToTheHost, not by more prose surface.
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

// rendererSourceFile is the file set the renderer can actually load: the text
// types tools/d22scan walks under frontend/, which is wider than the trio
// (.ts/.tsx/.css) door 2 used to open. R-92-1's first shape is a route spelled in
// a plain .js file, a file the old filter never read at all.
//
// frontend/dist is excluded on purpose, not by omission: it is generated output,
// so a gate over it would flip red and green with whether somebody re-ran
// npm run build rather than with what the source says (this ticket has not
// rebuilt it - see its residual 3). scripts/d22scan.sh already walks dist for
// ban #6, and no bundle can hold a call site its own source tree did not
// already contain, so the structural nails below reach dist without reading
// generated text.
func rendererSourceFile(p string) bool {
	switch filepath.Ext(p) {
	case ".ts", ".tsx", ".js", ".mjs", ".jsx", ".css", ".html":
		return true
	}
	return false
}

// hostBridgeCallRe matches a CALL into a postMessage-style host bridge. The
// leading dot is required, so panel.ts's own interface declaration
// "postMessage(message: string): void;" is not counted as a call site.
var hostBridgeCallRe = regexp.MustCompile(`\.\s*postMessage\s*\(`)

// dottedJoinRe matches a name assembled at runtime out of pieces
// (["panel","mode","set"] then parts.join(".")) - the shape no literal regex can
// see however wide its extension set gets, so it is nailed by its own syntax.
var dottedJoinRe = regexp.MustCompile(`\.join\(\s*["']\.["']\s*\)`)

// computedSendRequestRe matches sendRequest( whose first argument is not a
// string literal: the remaining way to route the panel's one envelope towards a
// method the source never names.
var computedSendRequestRe = regexp.MustCompile("sendRequest\\(\\s*(?:[A-Za-z_$][\\w$]*|`|\\[)")

// routeLiteralRe pulls every "panel.<something>" string literal out of a line.
var routeLiteralRe = regexp.MustCompile(`"panel\.[A-Za-z.]+"`)

// rendererDoorReport is one run of the structural scan.
type rendererDoorReport struct {
	callSites []string // every host call site in the tree
	outside   []string // ... of which: not inside src/lib/panel.ts
	assembly  []string // route names built at runtime
	computed  []string // sendRequest called with a non-literal route
	unknown   []string // a "panel.*" literal the Go side does not answer
	files     int
}

// composerRouteLiterals is the closed vocabulary the renderer may name. Growing
// it is a two-sided change on purpose: adding a route has to touch the Go side
// that answers it for this set to accept it.
func composerRouteLiterals() map[string]bool {
	return map[string]bool{
		`"` + MethodModeRequest + `"`:      true,
		`"` + MethodWorkspaceRequest + `"`: true,
		`"` + MethodAttachmentAdd + `"`:    true,
		`"` + MethodMessageSend + `"`:      true,
		// The approval card's route, whose exact spelling is pinned by
		// TestFrontendComposerRequestsMatchTheEnvelope.
		`"panel.approval.request"`: true,
	}
}

// scanRendererHostDoors runs the structural nails over one frontend/src tree. It
// takes a directory rather than the repo root so the planted-shape case can point
// the same instrument at a tree that is knowingly wrong.
func scanRendererHostDoors(t *testing.T, srcDir string) rendererDoorReport {
	t.Helper()
	allowed := composerRouteLiterals()
	rep := rendererDoorReport{}
	err := filepath.WalkDir(srcDir, func(p string, d fs.DirEntry, err error) error {
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
		if !rendererSourceFile(p) {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rep.files++
		rel := filepath.ToSlash(p)
		if i := strings.Index(rel, "/frontend/"); i >= 0 {
			rel = rel[i+len("/frontend/"):]
		} else {
			rel = filepath.Base(p)
		}
		inPanelLib := strings.HasSuffix(rel, "src/lib/panel.ts")
		for i, raw := range strings.Split(string(data), "\n") {
			line := codeOnly(raw)
			if line == "" {
				continue
			}
			at := rel + ":" + strconv.Itoa(i+1) + ": " + line
			if hostBridgeCallRe.MatchString(line) {
				rep.callSites = append(rep.callSites, at)
				if !inPanelLib {
					rep.outside = append(rep.outside, at)
				}
			}
			if dottedJoinRe.MatchString(line) {
				rep.assembly = append(rep.assembly, at)
			}
			if strings.Contains(line, "function sendRequest") {
				continue
			}
			if computedSendRequestRe.MatchString(line) {
				rep.computed = append(rep.computed, at)
			}
			for _, lit := range routeLiteralRe.FindAllString(line, -1) {
				if !allowed[lit] {
					rep.unknown = append(rep.unknown, at+" names "+lit)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", srcDir, err)
	}
	return rep
}

// TestTheRendererHoldsExactlyOneDoorToTheHost is AC#1 nailed structurally instead
// of lexically, and it is the answer to R-92-1: door 2 was a literal scan over
// three extensions and the "two postMessage call sites" nail counted inside
// panel.ts alone, so a .js file sending panel.mode.set and a route assembled at
// runtime both walked through. What is pinned here instead is the reachability
// set - which file may talk to the host at all, and which words it may say - so
// neither shape survives regardless of how the method name is spelled.
//
//	(i)   every .postMessage( call site under frontend/src lives in lib/panel.ts,
//	      and there are exactly two of them;
//	(ii)  no dotted name is assembled at runtime anywhere in the renderer;
//	(iii) sendRequest is never called with a computed route;
//	(iv)  every "panel.*" literal is one the Go side answers.
func TestTheRendererHoldsExactlyOneDoorToTheHost(t *testing.T) {
	root := panelRepoRoot(t)
	rep := scanRendererHostDoors(t, filepath.Join(root, "frontend", "src"))
	if len(rep.outside) != 0 {
		t.Errorf("the renderer reaches the host from outside src/lib/panel.ts; there is no second channel in this design"+
			" (ticket 92 AC#1, R-92-1):%s", "\n  "+strings.Join(rep.outside, "\n  "))
	}
	if len(rep.callSites) != 2 {
		t.Errorf("host call sites in frontend/src = %d, want 2 (the approval request and sendRequest):%s",
			len(rep.callSites), "\n  "+strings.Join(rep.callSites, "\n  "))
	}
	if len(rep.assembly) != 0 {
		t.Errorf("a dotted route name is assembled at runtime, which is a method name the source never states and no"+
			" literal gate can judge:%s", "\n  "+strings.Join(rep.assembly, "\n  "))
	}
	if len(rep.computed) != 0 {
		t.Errorf("sendRequest is called with a computed route; the envelope's method must be a literal so the"+
			" vocabulary stays closed:%s", "\n  "+strings.Join(rep.computed, "\n  "))
	}
	if len(rep.unknown) != 0 {
		t.Errorf("the renderer names a route the Go side does not answer:%s", "\n  "+strings.Join(rep.unknown, "\n  "))
	}
	t.Logf("renderer door held: %d files scanned, %d host call sites (all in src/lib/panel.ts), %d panel.* route literals",
		rep.files, len(rep.callSites), len(composerRouteLiterals()))
}

// TestPlantedRendererDoorShapesGoRed is AC#1's positive nail for the two shapes
// the acceptance round measured as GREEN under the old instruments: a
// panel.mode.set sent from a plain .js file, and the same route assembled at
// runtime from an array. Both are planted into a copy of the real tree, so the
// run also shows the clean half still passes.
func TestPlantedRendererDoorShapesGoRed(t *testing.T) {
	root := panelRepoRoot(t)
	real := filepath.Join(root, "frontend", "src")
	dir := t.TempDir()
	// Carry the honest half along, otherwise "outside panel.ts" would be red for
	// the wrong reason (no panel.ts at all).
	os.MkdirAll(filepath.Join(dir, "frontend", "src", "lib"), 0o755)
	data, err := os.ReadFile(filepath.Join(real, "lib", "panel.ts"))
	if err != nil {
		t.Fatalf("read src/lib/panel.ts: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "frontend", "src", "lib", "panel.ts"), data, 0o644); err != nil {
		t.Fatalf("plant panel.ts: %v", err)
	}
	// Shape A: the extension door 2 never opened, spelling a banned route.
	plantA := `window.chrome.webview.postMessage(JSON.stringify({ method: "panel.mode.set", to: "auto_approve" }));`
	// Shape B: no banned token anywhere in the text; the name is built at runtime.
	plantB := "const parts = [\"panel\", \"mode\", \"set\"];\n" +
		"const route = parts.join(\".\");\n" +
		"window.chrome?.webview.postMessage(JSON.stringify({ method: route }));\n"
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, "frontend", "src", name), []byte(body+"\n"), 0o644); err != nil {
			t.Fatalf("plant %s: %v", name, err)
		}
	}
	write("ac92-plant-a.js", plantA)
	write("ac92-plant-c.ts", plantB)

	rep := scanRendererHostDoors(t, filepath.Join(dir, "frontend", "src"))
	if len(rep.outside) < 2 {
		t.Errorf("shape A and shape B each add a host call site outside src/lib/panel.ts; the scan found %d:%s",
			len(rep.outside), "\n  "+strings.Join(rep.outside, "\n  "))
	}
	if len(rep.assembly) != 1 {
		t.Errorf("shape B's runtime-assembled route = %d hits, want 1:%s",
			len(rep.assembly), "\n  "+strings.Join(rep.assembly, "\n  "))
	}
	if len(rep.unknown) != 1 || !strings.Contains(rep.unknown[0], "panel.mode.set") {
		t.Errorf("shape A's route literal should be named once, got %v", rep.unknown)
	}
	// The widened door 2 must also reach a .js file now: this is the half R-92-1
	// measured as invisible, and it is fixed by extension, not by loosening.
	hits := scanComposerPermissionWrites(t, filepath.Join(dir, "frontend"))
	if len(hits) == 0 {
		t.Error("door 2 read neither plant: its extension set is still narrower than what the renderer loads")
	} else if !strings.Contains(strings.Join(hits, "\n"), "ac92-plant-a.js") {
		t.Errorf("door 2 must name the planted .js file, got %v", hits)
	}
	t.Logf("red as required: %d outside call sites, %d assembled routes, %d unknown literals, door 2 named %d line(s)",
		len(rep.outside), len(rep.assembly), len(rep.unknown), len(hits))
}

// scanComposerPermissionWrites applies BOTH doors to one tree and returns the
// offending lines (empty = clean).
//
// The two doors are applied differently on purpose:
//
//	door 1 (ban #6) is applied to the raw text, exactly as tools/d22scan walks
//	   frontend/, because that is the instrument whose verdict AC#6 reports and
//	   it does not care whether a banned identifier sat in a comment;
//	door 2 (this ticket's mode-write gate) is applied to code only. It is a
//	   naming gate over prose-heavy files, and a component that documents "do
//	   not call the setter" in its header would otherwise read as a violation.
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
		if !rendererSourceFile(p) {
			return nil
		}
		if hits := reMatchesInFile(panelDecisionIdentifierRe, p); len(hits) > 0 {
			out = append(out, hits...)
		}
		if hits := modeWriteMatchesInCodeFile(p); len(hits) > 0 {
			out = append(out, hits...)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(out)
	return out
}

// modeWriteMatchesInCodeFile applies door 2 to the code half of a file.
func modeWriteMatchesInCodeFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{"unreadable: " + err.Error()}
	}
	var out []string
	for i, line := range strings.Split(string(data), "\n") {
		code := codeOnly(line)
		if composerModeWriteRe.MatchString(code) {
			out = append(out, filepath.Base(path)+":"+strconv.Itoa(i+1)+": "+strings.TrimSpace(code))
		}
	}
	return out
}

// codeOnly strips the comment half of one TypeScript/TSX line: a trailing //
// comment, and any line that is nothing but comment syntax (// ... or the *
// continuation of a block comment). A line OPENING with /* is still scanned,
// because "/* whatever */ setMode(...)" is code wearing a comment's clothes, and
// that is the one shape a stripping filter must not forgive.
func codeOnly(line string) string {
	t := strings.TrimSpace(line)
	if strings.HasPrefix(t, "//") || (strings.HasPrefix(t, "*") && !strings.HasPrefix(t, "**")) ||
		strings.HasPrefix(t, "*/") {
		return ""
	}
	if i := strings.Index(line, "//"); i >= 0 {
		// Do not cut inside a string literal - the shapes that matter here are
		// URLs and method names, both of which can carry slashes.
		if !strings.Contains(strings.TrimSpace(line[:i]), `"`) &&
			!strings.Contains(strings.TrimSpace(line[:i]), "'") {
			line = line[:i]
		}
	}
	return line
}

// TestComposerRenderFixtureTellsTheTruth reads the committed render evidence
// (frontend/fixtures/composer-states.html, produced by `npm run render:composer`
// over the real React component) and checks the three states it claims. This is
// the non-interactive render check: it proves the values the composer is given
// end up on screen, including the two states where lying is easiest - no
// snapshot yet, and an attachment the native side refused.
//
// It is not a screenshot. The differential-screenshot pass for owner's sign-off
// is still open (ticket 92 AC#6's residual), and this test says nothing about
// pixels, only about painted text.
func TestComposerRenderFixtureTellsTheTruth(t *testing.T) {
	path := filepath.Join(panelRepoRoot(t), "frontend", "fixtures", "composer-states.html")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the render fixture: %v - regenerate it with `npm run render:composer`", err)
	}
	blocks := map[string]string{}
	for _, chunk := range strings.Split(string(data), "<!-- ") {
		if chunk == "" {
			continue
		}
		name, body, found := strings.Cut(chunk, " -->\n")
		if !found {
			continue
		}
		blocks[strings.TrimSpace(name)] = body
	}
	for _, name := range []string{
		"no host attached", "workspace chosen, video stored",
		"unsupported attachment told to the user",
	} {
		if _, ok := blocks[name]; !ok {
			t.Errorf("the render fixture has no block for %q - the harness and this test disagree "+
				"about which states are evidenced", name)
		}
	}
	if b := blocks["no host attached"]; b != "" {
		for _, want := range []string{"档位未知", "尚未收到原生侧的状态快照"} {
			if !strings.Contains(b, want) {
				t.Errorf("with no snapshot at all the row does not say so (%q missing):\n%s", want, b)
			}
		}
		if strings.Contains(b, "全自动") {
			t.Errorf("an unknown mode rendered as a named档:\n%s", b)
		}
	}
	if b := blocks["workspace chosen, video stored"]; b != "" {
		for _, want := range []string{`D:\work\Wisp\notes`, "clip.mp4", "video/mp4", "需原生 L2 强确认"} {
			if !strings.Contains(b, want) {
				t.Errorf("the chosen-workspace state is missing %q:\n%s", want, b)
			}
		}
	}
	if b := blocks["unsupported attachment told to the user"]; b != "" {
		if !strings.Contains(b, "不受支持") || !strings.Contains(b, "MZ") {
			t.Errorf("a refused attachment is not told to the user verbatim:\n%s", b)
		}
		if strings.Contains(b, "已存入附件目录") {
			t.Errorf("a refused attachment ALSO renders as stored - the user cannot tell what happened:\n%s", b)
		}
	}
	t.Logf("render fixture verified across %d painted states", len(blocks))
}
