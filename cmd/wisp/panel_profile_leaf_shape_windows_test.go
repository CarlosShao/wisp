//go:build windows

package main

// Ticket 314 AC#4's source-shape ruler: a fixed-name panel WebView2 user-data
// leaf anywhere in this package's sources is red, so the de-sharing cannot be
// undone silently.
//
// WHY A SOURCE RULER AND NOT A BEHAVIOURAL ONE. The 2026-11-14 non-implementer
// audit (docs/evidence/s1/today-four-fixes-audit.md section 4) put the default
// tier's panel_unbound_253r6_windows_test.go:71 back to
// filepath.Join(os.TempDir(), "wisp-253r6-panel-profile"): the mutant compiled
// (go test -c rc=0) and Test253r6TheSilentDropBranchIsLoud PASSED. Nothing in
// this repo asserted that a profile leaf was unique; the only control was
// .scratch/wisp/probes/314-deshare/concurrent.ps1, which no gate calls. A
// behavioural ruler cannot see this class anyway: a shared folder is only a
// defect when two runs land at the same moment (the same audit's A/B readings
// were rc=0 on BOTH sides of the fixed-name shape). So the shape itself is the
// subject - "is the leaf a fixed name" is a fact about the source, it is the one
// fact the next revert changes, and this ruler reads it in under a second
// without creating a window.
//
// WHAT IT READS. Every cmd/wisp/*.go - production and test, because a literal
// can move one file over - plus the one out-of-package source that can point a
// real WebView2 at a reusable folder: scripts/spike/webview2-latency/main.go.
// It parses with go/parser, NOT with a string search, for three reasons:
// comments in this package discuss os.TempDir() constantly and a text ruler
// would convict prose; a parse yields the exact line the red sentence must name;
// and the leaf can be judged as a literal ARGUMENT rather than as any occurrence
// of a name. A source that cannot be read or parsed is FATAL, never skipped - a
// silently skipped file is how a shape ruler goes blind, and an empty scan set
// fatals for the same reason.
//
// THE SHAPES IT CALLS FIXED - two rules, unioned (panelProfileHitsInSource), and
// both judge the leaf on the value AFTER constant folding (panelProfileFold):
// `"wisp-33r1" + "-panel-profile"`, a strings.Join over a literal slice, and a
// Join whose result gets a literal suffix are folded first, so a leaf assembled
// from pieces is judged as the name it spells rather than as the token it was
// written as. A leaf that is only PART constant and part runtime value is not a
// fixed name and stays silent on purpose (a constant prefix plus a pid is
// per-run).
//   - NAME: a join carrying a temp root together with a leaf that spells a
//     profile folder ("...profile"), in either argument order, across wrapped
//     lines, through a local or package-level alias of the root (`root :=
//     os.TempDir()` then Join(root, "...")), as `os.TempDir() + "/literal"`
//     concatenation, or as a leaf folded out of literal pieces. This is the
//     audit's exact shape, caught wherever it appears - sink or no sink - which
//     is what keeps a leaf returned from a helper (the 319 holder fixture)
//     inside the ruler;
//   - SINK: the same join with ANY fixed leaf once the value reaches a WebView2
//     data-path sink - NewPanelManager, newPanelControl, a DataPath field, a
//     DataPath key in a struct literal - directly, through a local binding,
//     through a package-level binding, or at package level itself (a `var` whose
//     initialiser is a struct literal with a DataPath: key IS a sink where it is
//     declared). Without it, a fixed name that avoids the word "profile" would
//     be invisible.
// A temp root is not only os.TempDir(): os.Getenv("TEMP"), os.Getenv("TMP"),
// os.UserHomeDir() (also through `home, err := os.UserHomeDir()`), and a
// drive-letter absolute literal (`C:\Temp\...`, which carries the whole path and
// needs no root expression) all count. The root and the leaf still have to meet
// in one expression in the scanned file - that is boundary 1 below. A lone
// absolute literal is a path rather than a join, so the NAME rule leaves it to
// the SINK rule: it is a shape where it reaches a data-path sink (which is the
// audit's A3c), and a quoted Windows path that merely spells "profile" in test
// data is not accused.
// A fixed temp leaf for another purpose stays silent on purpose: ticket 320's
// planted-key control file in os.TempDir() is a regular file and a detector
// control, not a browser profile folder, and accusing it would be wrong.
//
// WHAT IT STILL CANNOT SEE - named, not hidden, the same honesty the 197 publish
// ruler's audit asked for (today-four-fixes-audit section 3 D):
//  1. a helper that joins a literal with a temp root behind a parameter, and a
//     helper whose fixed join (with a non-profile leaf) is handed to a sink by
//     its caller: the root and the literal never meet in one expression in the
//     scanned file, and this ruler does not follow values across a function
//     boundary. Closing either needs interprocedural analysis, a different
//     instrument (the 2026-11-14 audit's A1b/A1c);
//  2. a root that is none of the four above: os.UserCacheDir(), a UNC literal,
//     %TEMP% re-spelled through another variable;
//  3. a leaf that looks fixed but ends in a runtime value (`Join(root, "wisp-")
//     + pid`): the constant text is mixed with a runtime value, so this ruler
//     treats the name as per-run. That is a deliberate narrowing, not an
//     oversight - a name that is only half constant is not the shared-folder
//     defect.
// Moving a fixed name into one of those three is a new decision to argue on the
// ticket, not an edit this ruler silently blesses.
//
// THE 2026-11-14 AUDIT'S SILENT SHAPES, ONE BY ONE (docs/evidence/s1/
// 314-ruler-audit.md). Every FIXED line below has a resident must-red control in
// Test314PanelProfileShapeRulerRingsOnTheFixedNameShape:
//   A2a "wisp-33r1"+"-panel-profile" inside one Join   FIXED  (fold)
//   A2b the whole path by + concatenation              FIXED  (fold)
//   A2c leaf by strings.Join                           FIXED  (fold)
//   A2d Join(...) + a literal suffix                   FIXED  (fold)
//   A2e folded value bound to a local, then to a sink  FIXED  (fold + SINK binding)
//   A3a root os.Getenv("TEMP")                         FIXED  (root rule; TEMP and TMP)
//   A3b root from os.UserHomeDir()                     FIXED  (root rule + 2-value alias)
//   A3c root is a drive-letter absolute literal        FIXED  (root rule)
//   X1  package-level DataPath: key                    FIXED  (SINK over package decls)
//   X2  package-level fixed Join given to a sink       FIXED  (package bindings seeded)
//   A1b helper return with a plain leaf, value to sink STILL OUT (function boundary)
//   A1c root behind a helper parameter                 STILL OUT (root and literal never meet)
//
// EXEMPTIONS (panelProfileLeafExemptions) are named HERE, in the test - never in
// a workflow file or a scanner allowlist - and each is checked three ways:
//   - it must still cover at least one live hit (a stale exemption is red, so
//     the table cannot accumulate blank cheques);
//   - every string in its quotes must still appear verbatim in that file (the
//     reason is read from the source, so deleting the justifying comment kills
//     the exemption instead of leaving the claim behind);
//   - an exemption carrying noWindow reds if the file starts calling the
//     window-creation path (bringUp / newPanelControl / startPanelHolder319):
//     "nothing here creates a window" is a measured claim, not a promise.
//
// Ticket 314's audit ruling: panel_holder_319_windows_test.go keeps its fixed
// name (headless tier, deliberately shared) and this table is where that is
// booked. Every other source must go through panelProfileLeaf.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// panelProfileSpikeSource is the one source outside cmd/wisp that can hand a
// real WebView2 a reusable user-data folder. It is NAMED, not globbed: a file
// that moved must turn this ruler red (the scan set and the exemption table have
// to move in the same commit), never shrink the scan quietly.
const panelProfileSpikeSource = "../../scripts/spike/webview2-latency/main.go"

// panelProfileSpikeDisplay is the repo-relative spelling the red sentence uses
// for the source above (the test runs with cmd/wisp as its working directory).
const panelProfileSpikeDisplay = "scripts/spike/webview2-latency/main.go"

// panelProfileWindowCalls are the call names that put a real window on screen.
// An exemption claiming "this file creates no window" is re-checked against
// these on every run.
var panelProfileWindowCalls = []string{"bringUp", "newPanelControl", "startPanelHolder319"}

// panelProfileShapeHit is one fixed-name leaf found in one source line.
type panelProfileShapeHit struct {
	display string // repo-relative path, as the red sentence must name it
	line    int    // 1-based line of the offending expression
	leaf    string // the fixed leaf, after constant folding, unquoted
	raw     string // the literal token that leaf starts with (quoted, as written)
	code    string // that source line, trimmed, so the red shows the shape
}

// panelProfileLeafExemption books one deliberate shared folder. file+leaf are
// the identity (not a line number: lines drift, a file and its exact literal do
// not), quotes are the justification read out of that same file, and noWindow
// says the file has no window at all.
type panelProfileLeafExemption struct {
	file     string
	leaf     string
	reason   string
	quotes   []string
	noWindow bool
}

// panelProfileLeafExemptions is the complete list. Two entries, one reason each,
// every one readable off the source it points at. Widening this table is how the
// ruler would be defeated, so each entry has to survive the three checks in the
// header; there is deliberately no prefix/glob form.
var panelProfileLeafExemptions = []panelProfileLeafExemption{
	{
		file: "cmd/wisp/panel_holder_319_windows_test.go",
		leaf: "wisp-319r1-holder-profile",
		reason: "本文件一次窗都不建（headless tier 从不向 bringUp 要窗口），所以那个夹只是一个名字、" +
			"从来不是 WebView2 user-data 夹；本文件内两个夹具共用它是故意的。",
		quotes: []string{
			"Nothing here creates a window (the headless half never asks bringUp), so",
			"shared here on purpose: two fixtures in one run would",
		},
		noWindow: true,
	},
	{
		file: "scripts/spike/webview2-latency/main.go",
		leaf: "wisp-spike-webview2-profile",
		reason: "手工逐发跑的延迟测量工具（P11/D32 的 cold/hot 基数）的 -datapath 默认值：共用同一个夹正是它要测的 warm profile，" +
			"每一次调用都由人显式串行发起、且 -datapath 可覆盖；测试套件从不启动它。",
		quotes: []string{
			`flag.String("datapath"`,
			`"WebView2 user data folder"`,
		},
	},
}

// panelProfileShapeSource is one scanned file: its display path, its bytes, and
// the parse the census walks.
type panelProfileShapeSource struct {
	display string
	path    string
	src     []byte
	file    *ast.File
	fset    *token.FileSet
}

// Test314FixedPanelProfileLeafShapeIsRed is the ruler. A fixed leaf that is not
// in the exemption table reds with the file, the line and the reason the shape
// matters; a source it cannot read or parse fatals rather than being skipped.
func Test314FixedPanelProfileLeafShapeIsRed(t *testing.T) {
	sources := panelProfileShapeSources(t)
	var hits []panelProfileShapeHit
	for _, src := range sources {
		hits = append(hits, panelProfileHitsInSource(t, src)...)
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].display != hits[j].display {
			return hits[i].display < hits[j].display
		}
		return hits[i].line < hits[j].line
	})

	exempted := 0
	for _, hit := range hits {
		if ex, ok := panelProfileLeafExemptionFor(hit); ok {
			exempted++
			t.Logf("豁免命中：%s:%d 叶 %q -> %s", hit.display, hit.line, hit.leaf, ex.reason)
			continue
		}
		panelProfileFixedNameRed(t, hit)
	}

	faults := panelProfileCheckExemptions(t, hits, sources)

	// The summary carries all three counts, and a run with any red in it must not
	// be readable as green: the audit saw the old line print "0 处红" while this
	// test was failing on a stale exemption reason.
	line, _ := panelProfileSummaryLine(len(sources), len(hits), exempted, faults)
	t.Logf("%s", line)
}

// panelProfileSummaryLine renders the one line a reader takes for the verdict. It
// prints the three numbers (shape hits / exempted / actually red) plus the
// exemption self-check reds, and it says 红 whenever ANY of them is non-zero. The
// old line counted only unexempted shape hits, so a stale exemption could fail
// this test while the summary still read "0 处红" - the audit's misleading-green
// reading. The verdict word is decided here, once, so no caller can phrase it.
func panelProfileSummaryLine(sources, hits, exempted, faults int) (string, bool) {
	red := hits - exempted
	clean := red == 0 && faults == 0
	verdict := "红"
	if clean {
		verdict = "绿"
	}
	return fmt.Sprintf("314 形状尺【%s】：扫了 %d 个源文件，形状命中 %d 处；其中已豁免 %d 处、实红 %d 处，豁免自检红 %d 处 —— 只有「实红 0 且自检红 0」才是绿，任何一处红都不能当绿证",
		verdict, sources, hits, exempted, red, faults), clean
}

// panelProfileFixedNameRed is the one sentence this ruler exists to say. It
// names the file and line, prints the shape, says what breaks when two runs
// share the folder, and quotes the audit reading that proved nothing else in the
// repo can stand in for it.
func panelProfileFixedNameRed(t *testing.T, hit panelProfileShapeHit) {
	t.Helper()
	// A folded leaf is not written on one line, so the reader is told where the
	// quoted name came from instead of being left to wonder why the line below
	// does not contain it.
	foldNote := ""
	if written, err := strconv.Unquote(hit.raw); err == nil && written != hit.leaf {
		foldNote = "\n（这枚叶由多段字面量折叠得到：上面那行是它的拼法，红句写的是折叠后的值）"
	}
	t.Errorf("%s:%d 命中固定名面板 profile 形状（叶 %q）：\n    %s\n"+
		"这个形状会让两发并发互踩同一个 WebView2 user-data 夹：一个夹就是一个浏览器 profile，"+
		"两个测试进程在同一条夹上各起一套顶层 browser host —— 先到的 CreateEnvironment 争同一夹会回 S_OK+NULL environment，"+
		"后到的 environment 会 JOIN 一个正在退场的浏览器进程（internal/ball/orb_windows.go:432-458 的两条守卫）；"+
		"票 314 的测量把同一夹同 tick 两个顶层 host 抓成过约 1.2 秒的共存读数。%s\n"+
		"ticket 314 的非实现者审计量过这形的牙是零：把默认档 panel_unbound_253r6_windows_test.go:71 改回固定名，"+
		"go test -c rc=0、Test253r6 PASS —— 仓里没有断言读叶名唯一性，所以只有这把源码形状尺能让这种回退红。\n"+
		"改法＝panelProfileLeaf(%q)（cmd/wisp/panel_host_windows_test.go 的既有 helper，叶＝prefix-pid-纳秒）；"+
		"确有理由共用同一个夹的，在 panelProfileLeafExemptions 里按「文件＋叶名」具名登记，并把源码里那句理由逐字抄进 quotes。",
		hit.display, hit.line, hit.leaf, hit.code, foldNote, hit.leaf)
}

// panelProfileCheckExemptions runs the three checks every exemption owes: it is
// live, its quoted reason is still on disk, and a noWindow claim is still true.
// It answers how many of those checks went red, so the summary line can count
// them: a stale exemption is a red run even when no shape hit is unexempted.
func panelProfileCheckExemptions(t *testing.T, hits []panelProfileShapeHit, sources []panelProfileShapeSource) int {
	t.Helper()
	faults := 0
	byDisplay := make(map[string]panelProfileShapeSource, len(sources))
	for _, src := range sources {
		byDisplay[src.display] = src
	}
	for _, ex := range panelProfileLeafExemptions {
		covered := 0
		for _, hit := range hits {
			if hit.display == ex.file && hit.leaf == ex.leaf {
				covered++
			}
		}
		if covered == 0 {
			faults++
			t.Errorf("豁免表过期：%s 的固定名叶 %q 在这次扫描里覆盖 0 处。留着它就是把一张空白支票挂在尺子上："+
				"删掉这条，或者把叶名改回它真正豁免的那一处。", ex.file, ex.leaf)
		}
		src, ok := byDisplay[ex.file]
		if !ok {
			faults++
			t.Errorf("豁免表点名了扫描集之外的源文件 %s：扫描集与豁免表必须同一次改动（否则这条豁免只活在纸上）", ex.file)
			continue
		}
		for _, quote := range ex.quotes {
			if !strings.Contains(string(src.src), quote) {
				faults++
				t.Errorf("豁免 %s 的理由句在源码里找不到了：%q\n"+
					"这条豁免的正当性是从那段注释读出来的；注释改掉后它还留着，就成了一张没有出处的通行证。"+
					"重读源码、逐字重抄一句，或删掉这条豁免。", ex.file, quote)
			}
		}
		if !ex.noWindow {
			continue
		}
		if name, line, ok := panelProfileWindowCallIn(src); ok {
			faults++
			t.Errorf("豁免 %s 自述「这里不建窗」，但 %s:%d 调用了 %s：这句话现在是假的，"+
				"那处固定名夹也就重新变成两发互踩的碰撞面。要么把固定名改成 panelProfileLeaf，要么删掉这条 noWindow 豁免并重新论证。",
				ex.file, ex.file, line, name)
		}
	}
	return faults
}

// panelProfileLeafExemptionFor answers the exemption covering hit, if any. The
// identity is (file, leaf): a changed literal is a new shape and gets no ride.
func panelProfileLeafExemptionFor(hit panelProfileShapeHit) (panelProfileLeafExemption, bool) {
	for _, ex := range panelProfileLeafExemptions {
		if ex.file == hit.display && ex.leaf == hit.leaf {
			return ex, true
		}
	}
	return panelProfileLeafExemption{}, false
}

// panelProfileShapeSources builds the scan set and parses it. Read/parse
// failures are fatal: this ruler's subject is source text, and a subject it
// could not read is not an exempt subject.
func panelProfileShapeSources(t *testing.T) []panelProfileShapeSource {
	t.Helper()
	names, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("列 cmd/wisp/*.go 失败：%v —— 这把尺连主体都找不出来", err)
	}
	sort.Strings(names)
	if len(names) == 0 {
		t.Fatalf("扫描集为空：cmd/wisp 下没有 *.go。空扫描集不许印绿（盲尺不是干净的尺）")
	}
	out := make([]panelProfileShapeSource, 0, len(names)+1)
	for _, name := range names {
		out = append(out, panelProfileShapeSourceWithSource(t, "cmd/wisp/"+name, name))
	}
	out = append(out, panelProfileShapeSourceWithSource(t, panelProfileSpikeDisplay, panelProfileSpikeSource))
	return out
}

// panelProfileShapeSourceWithSource reads and parses one named source.
func panelProfileShapeSourceWithSource(t *testing.T, display, path string) panelProfileShapeSource {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读不到 %s（%s）：形状尺的落点是这棵树里的源文件；读不到就是尺子瞎了，跳过不是豁免", display, path)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("%s 解析失败：%v —— 一把读不了源文件的尺子等于没有（红要红在形状上，不能红在解析器上还假装扫过）", display, err)
	}
	return panelProfileShapeSource{display: display, path: path, src: src, file: file, fset: fset}
}

// panelProfileHitsInSource is the census of one parsed source. It applies two
// independent rules and unions what they find:
//
//   - the NAME rule: a join carrying a temp root and a leaf that calls itself a
//     profile folder ("...profile") is the audit's shape wherever it appears,
//     sink or no sink. This is what keeps a leaf that is returned from a helper
//     (the 319 holder fixture) inside the ruler instead of outside it;
//   - the SINK rule: a join carrying a temp root and ANY fixed leaf once its
//     value reaches a WebView2 data-path sink (NewPanelManager / newPanelControl,
//     a DataPath field, a DataPath key in a struct literal), directly, through a
//     local binding, through a package-level binding, or at package level. It
//     keeps a fixed name that avoids the word "profile" from being invisible
//     just because of its spelling.
//
// Both rules read the leaf after constant folding (panelProfileFold), so a leaf
// assembled from literals is judged as the name it spells.
//
// A fixed temp leaf for another purpose is silent on purpose. Ticket 320's
// planted-key control file in os.TempDir() is the measured one: it is a regular
// file and a detector control, not a browser profile folder, and a red that
// accused it would be a false accusation.
func panelProfileHitsInSource(t *testing.T, src panelProfileShapeSource) []panelProfileShapeHit {
	t.Helper()
	lines := strings.Split(string(src.src), "\n")
	packageRoots := panelProfilePackageTempRoots(src.file)
	packageJoins := panelProfilePackageFixedJoins(t, src.display, src.file, packageRoots)
	var hits []panelProfileShapeHit
	seen := map[string]bool{}
	// The line is the LEAF literal's line, not the call's: on a wrapped Join the
	// literal is the thing a reader has to change, and this keeps the printed
	// source line containing the leaf it names.
	add := func(fixed panelProfileFixedJoin) {
		line := src.fset.Position(fixed.pos).Line
		key := strconv.Itoa(line) + "|" + fixed.leaf
		if seen[key] {
			return
		}
		seen[key] = true
		code := ""
		if line >= 1 && line <= len(lines) {
			code = strings.TrimSpace(lines[line-1])
		}
		hits = append(hits, panelProfileShapeHit{display: src.display, line: line, leaf: fixed.leaf, raw: fixed.raw, code: code})
	}
	// Package-level declarations carry BOTH rules, not just NAME: a package-level
	// `var` initialiser can hold the fixed join (NAME) and can be the sink itself
	// (a DataPath: key in a struct literal). Until the 2026-11-14 audit this walk
	// was NAME-only, so a package-level sink was invisible (X1).
	for _, decl := range src.file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			panelProfileNameRule(t, gen, src.display, packageRoots, add)
			panelProfileSinkRule(t, gen, src.display, packageRoots, packageJoins, add)
		}
	}
	// Function bodies: each one gets its own alias scope, so one function's
	// `tempRoot` cannot make another function's unrelated join look like a temp
	// root. Closures inside a body share that body's scope (over-approximate:
	// an extra red is possible, a missed one is not). The package-level fixed
	// joins are seeded into every body, so a package-level binding handed to a
	// sink is caught (X2).
	for _, decl := range src.file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		roots := panelProfileCopyRoots(packageRoots)
		panelProfileTempRoots(fn.Body, roots)
		panelProfileNameRule(t, fn.Body, src.display, roots, add)
		panelProfileSinkRule(t, fn.Body, src.display, roots, packageJoins, add)
	}
	return hits
}

// panelProfileNameRule is the shape the ticket is named after: a temp root
// joined with a leaf that spells a profile folder.
func panelProfileNameRule(t *testing.T, node ast.Node, display string, roots map[string]bool, add func(panelProfileFixedJoin)) {
	t.Helper()
	ast.Inspect(node, func(n ast.Node) bool {
		expr, ok := n.(ast.Expr)
		if !ok {
			return true
		}
		fixed, ok := panelProfileFixedJoinIn(t, display, expr, roots)
		if ok && !fixed.bare && strings.Contains(strings.ToLower(fixed.leaf), "profile") {
			add(fixed)
		}
		return true
	})
}

// panelProfileSinkRule catches a fixed join with ANY leaf once its value is
// handed to a WebView2 data-path sink - directly, through a local binding, or
// through a DataPath field/key. seeded carries the package-level fixed joins, so
// a package-level boundary is no longer a hiding place (audit X2). A helper
// return is still a boundary this rule cannot cross (header, boundary 1).
func panelProfileSinkRule(t *testing.T, node ast.Node, display string, roots map[string]bool, seeded map[string]panelProfileFixedJoin, add func(panelProfileFixedJoin)) {
	t.Helper()
	stored := make(map[string]panelProfileFixedJoin, len(seeded)+4)
	for name, fixed := range seeded {
		stored[name] = fixed
	}
	hits := func(expr ast.Expr) {
		if fixed, ok := panelProfileFixedJoinIn(t, display, expr, roots); ok {
			add(fixed)
		}
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.AssignStmt:
			if len(expr.Lhs) != 1 || len(expr.Rhs) != 1 {
				return true
			}
			name, isIdent := expr.Lhs[0].(*ast.Ident)
			fixed, isFixed := panelProfileFixedJoinIn(t, display, expr.Rhs[0], roots)
			if !isFixed {
				return true
			}
			if isIdent {
				stored[name.Name] = fixed
				return true
			}
			if sel, isSel := expr.Lhs[0].(*ast.SelectorExpr); isSel && sel.Sel != nil && panelProfileSinkKey(sel.Sel.Name) {
				add(fixed)
			}
		case *ast.ValueSpec:
			if len(expr.Names) == 1 && len(expr.Values) == 1 {
				if fixed, ok := panelProfileFixedJoinIn(t, display, expr.Values[0], roots); ok {
					stored[expr.Names[0].Name] = fixed
				}
			}
		case *ast.CallExpr:
			if !panelProfileSinkCall(panelProfileCalleeName(expr)) {
				return true
			}
			for _, arg := range expr.Args {
				if name, isIdent := arg.(*ast.Ident); isIdent {
					if fixed, ok := stored[name.Name]; ok {
						add(fixed)
					}
					continue
				}
				hits(arg)
			}
		case *ast.KeyValueExpr:
			key, isIdent := expr.Key.(*ast.Ident)
			if !isIdent || !panelProfileSinkKey(key.Name) {
				return true
			}
			hits(expr.Value)
		}
		return true
	})
}

// panelProfilePackageFixedJoins collects the package-level names bound exactly to
// a fixed join, so a function that hands one to a sink is caught (audit X2). A
// value computed at run time is not a fixed name and is not collected.
func panelProfilePackageFixedJoins(t *testing.T, display string, file *ast.File, roots map[string]bool) map[string]panelProfileFixedJoin {
	t.Helper()
	out := map[string]panelProfileFixedJoin{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok || len(value.Names) != 1 || len(value.Values) != 1 {
				continue
			}
			if fixed, ok := panelProfileFixedJoinIn(t, display, value.Values[0], roots); ok {
				out[value.Names[0].Name] = fixed
			}
		}
	}
	return out
}

// panelProfileFixedJoinIn answers the fixed leaf an expression carries when a
// temp root and constant text meet in it, after constant folding:
// Join(...root..., "leaf"), `root + "/leaf"`, a leaf assembled from literal
// pieces, and a drive-letter absolute literal (which is the whole path).
func panelProfileFixedJoinIn(t *testing.T, display string, expr ast.Expr, roots map[string]bool) (panelProfileFixedJoin, bool) {
	t.Helper()
	folded := panelProfileFold(t, display, expr, roots)
	if !folded.rooted || folded.leaf == "" || folded.opaque {
		return panelProfileFixedJoin{}, false
	}
	return panelProfileFixedJoin{
		leaf: folded.leaf,
		pos:  folded.pos,
		raw:  folded.raw,
		bare: panelProfileBareAbsolute(t, display, expr),
	}, true
}

// panelProfileBareAbsolute answers whether expr is nothing but one drive-letter
// absolute literal. Such a literal is a path, not a JOIN: it is only a defect
// where it reaches a sink (A3c is handed straight to NewPanelManager), so the
// NAME rule leaves it to the SINK rule. Without this split the ruler would
// convict any quoted Windows path that happens to spell "profile" - test data,
// documentation strings - which is the spelling self-certification the audit
// criticised, not a data path.
func panelProfileBareAbsolute(t *testing.T, display string, expr ast.Expr) bool {
	lit, ok := panelProfileUnparen(expr).(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	return panelProfileDriveAbsolute(panelProfileUnquote(t, display, lit.Value))
}

// panelProfileFixedJoin is one fixed leaf: what it names, where it is written,
// the literal token it starts with, and whether it is a lone absolute path. A
// folded leaf is not spelled verbatim on any one line, so the red sentence and
// the controls check the named line against that token instead.
type panelProfileFixedJoin struct {
	leaf string
	pos  token.Pos
	raw  string
	bare bool
}

// panelProfileFolded is one expression read as a piece of a data path: whether a
// temp root lives in it, the constant text it contributes with roots elided, and
// the literal token that text starts with.
type panelProfileFolded struct {
	rooted bool
	opaque bool
	leaf   string
	raw    string
	pos    token.Pos
}

// panelProfileJoinSeparator is the separator Join puts between its arguments.
// Keeping arguments apart (instead of pasting them together) stops two path
// SEGMENTS from spelling "profile" across the boundary between them.
var panelProfileJoinSeparator = string(filepath.Separator)

// panelProfileFold reads expr as constant text plus temp roots, so the leaf is
// judged on the folded value rather than on one literal token: a `+` chain, a
// strings.Join over a literal slice, and a Join over literal pieces all fold. A
// value that is only part constant is opaque, and an opaque expression is NOT a
// fixed name - that is what keeps a constant prefix plus a pid silent.
func panelProfileFold(t *testing.T, display string, expr ast.Expr, roots map[string]bool) panelProfileFolded {
	t.Helper()
	switch node := panelProfileUnparen(expr).(type) {
	case *ast.BasicLit:
		if node.Kind != token.STRING {
			return panelProfileFolded{opaque: true, pos: node.Pos()}
		}
		value := panelProfileUnquote(t, display, node.Value)
		return panelProfileFolded{
			rooted: panelProfileDriveAbsolute(value),
			leaf:   value,
			raw:    node.Value,
			pos:    node.Pos(),
		}
	case *ast.Ident:
		if roots[node.Name] {
			return panelProfileFolded{rooted: true, pos: node.Pos()}
		}
		return panelProfileFolded{opaque: true, pos: node.Pos()}
	case *ast.BinaryExpr:
		if node.Op != token.ADD {
			return panelProfileFolded{opaque: true, pos: node.Pos()}
		}
		left := panelProfileFold(t, display, node.X, roots)
		right := panelProfileFold(t, display, node.Y, roots)
		return panelProfileFolded{
			rooted: left.rooted || right.rooted,
			opaque: left.opaque || right.opaque,
			leaf:   left.leaf + right.leaf,
			raw:    panelProfileFirstRaw(left, right),
			pos:    panelProfileFirstPos(left, right, node.Pos()),
		}
	case *ast.CallExpr:
		if panelProfileIsTempRoot(node, roots) {
			return panelProfileFolded{rooted: true, pos: node.Pos()}
		}
		if panelProfileQualifiedCallee(node) == "strings.Join" {
			return panelProfileFoldStringsJoin(t, display, node, roots)
		}
		if panelProfileCalleeName(node) == "Join" {
			return panelProfileFoldJoin(t, display, node, roots)
		}
	}
	return panelProfileFolded{opaque: true, pos: expr.Pos()}
}

// panelProfileFoldJoin folds one Join over its arguments: literal pieces keep
// their order, roots contribute no text, and an argument the ruler cannot read
// makes the whole join opaque.
func panelProfileFoldJoin(t *testing.T, display string, node *ast.CallExpr, roots map[string]bool) panelProfileFolded {
	t.Helper()
	out := panelProfileFolded{pos: node.Pos()}
	var pieces []string
	for _, arg := range node.Args {
		part := panelProfileFold(t, display, arg, roots)
		out.rooted = out.rooted || part.rooted
		out.opaque = out.opaque || part.opaque
		if part.leaf == "" {
			continue
		}
		pieces = append(pieces, part.leaf)
		if out.raw == "" {
			out.raw, out.pos = part.raw, part.pos
		}
	}
	out.leaf = strings.Join(pieces, panelProfileJoinSeparator)
	return out
}

// panelProfileFoldStringsJoin folds strings.Join over a literal slice: the one
// shape the audit assembled with no Join of the path kind at all (A2c).
func panelProfileFoldStringsJoin(t *testing.T, display string, node *ast.CallExpr, roots map[string]bool) panelProfileFolded {
	t.Helper()
	if len(node.Args) != 2 {
		return panelProfileFolded{opaque: true, pos: node.Pos()}
	}
	slice, ok := panelProfileUnparen(node.Args[0]).(*ast.CompositeLit)
	if !ok {
		return panelProfileFolded{opaque: true, pos: node.Pos()}
	}
	sep := panelProfileFold(t, display, node.Args[1], roots)
	if sep.rooted || sep.opaque {
		return panelProfileFolded{opaque: true, pos: node.Pos()}
	}
	out := panelProfileFolded{pos: node.Pos()}
	pieces := make([]string, 0, len(slice.Elts))
	for _, elt := range slice.Elts {
		part := panelProfileFold(t, display, elt, roots)
		if part.rooted || part.opaque {
			return panelProfileFolded{opaque: true, pos: node.Pos()}
		}
		pieces = append(pieces, part.leaf)
		if out.raw == "" && part.raw != "" {
			out.raw, out.pos = part.raw, part.pos
		}
	}
	out.leaf = strings.Join(pieces, sep.leaf)
	return out
}

// panelProfileFirstRaw answers the raw token of the first literal the fold kept.
func panelProfileFirstRaw(a, b panelProfileFolded) string {
	if a.raw != "" {
		return a.raw
	}
	return b.raw
}

// panelProfileFirstPos answers the position of the first literal the fold kept,
// falling back to the expression's own position.
func panelProfileFirstPos(a, b panelProfileFolded, fallback token.Pos) token.Pos {
	if a.raw != "" {
		return a.pos
	}
	if b.raw != "" {
		return b.pos
	}
	return fallback
}

// panelProfileUnparen strips any number of parentheses.
func panelProfileUnparen(expr ast.Expr) ast.Expr {
	for {
		paren, ok := expr.(*ast.ParenExpr)
		if !ok {
			return expr
		}
		expr = paren.X
	}
}

// panelProfileSinkCall answers whether a call name takes a WebView2 user-data
// folder as an argument in this package.
func panelProfileSinkCall(callee string) bool {
	switch callee {
	case "NewPanelManager", "newPanelControl":
		return true
	}
	return false
}

// panelProfileSinkKey answers whether a field or struct key names the data path.
func panelProfileSinkKey(name string) bool {
	switch name {
	case "DataPath", "dataPath":
		return true
	}
	return false
}

// panelProfilePackageTempRoots collects the package-level names bound to a temp
// root (os.TempDir(), os.Getenv("TEMP"|"TMP"), os.UserHomeDir()). Function
// bodies are deliberately NOT walked here: their aliases belong to their own
// scope (see panelProfileHitsInSource).
func panelProfilePackageTempRoots(file *ast.File) map[string]bool {
	roots := map[string]bool{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		panelProfileTempRoots(gen, roots)
	}
	return roots
}

// panelProfileTempRoots adds to roots every name a temp root is bound to within
// node (a block, a GenDecl, or a whole function body).
func panelProfileTempRoots(node ast.Node, roots map[string]bool) map[string]bool {
	if roots == nil {
		roots = map[string]bool{}
	}
	record := func(lhs, rhs []ast.Expr) {
		if len(lhs) == 0 || len(rhs) != 1 {
			return
		}
		name, ok := lhs[0].(*ast.Ident)
		if !ok {
			return
		}
		// More than one name on the left is how a shared root is spelled when the
		// call also returns an error (`home, err := os.UserHomeDir()`): the first
		// result is the path, so that is the name bound to the root.
		if panelProfileIsRootSource(rhs[0], roots) {
			roots[name.Name] = true
		}
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.AssignStmt:
			record(expr.Lhs, expr.Rhs)
		case *ast.ValueSpec:
			lhs := make([]ast.Expr, len(expr.Names))
			for i, name := range expr.Names {
				lhs[i] = name
			}
			record(lhs, expr.Values)
		}
		return true
	})
	return roots
}

// panelProfileIsTempRoot reports whether expr is a temp root: os.TempDir(),
// os.Getenv("TEMP"|"TMP"), os.UserHomeDir(), a name bound to one of those, or a
// drive-letter absolute literal, which carries the whole path and needs no root
// expression (the audit's A3c).
func panelProfileIsTempRoot(expr ast.Expr, alias map[string]bool) bool {
	expr = panelProfileUnparen(expr)
	if panelProfileIsRootSource(expr, alias) {
		return true
	}
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	value, err := strconv.Unquote(lit.Value)
	return err == nil && panelProfileDriveAbsolute(value)
}

// panelProfileIsRootSource answers whether expr names a shared root the way a
// binding can carry it: one of the three os calls, or a name already bound to
// one. An absolute literal is deliberately NOT a source here: it is a root where
// it is written, and promoting every absolute base directory to "temp root"
// through a variable would convict hard-coded asset paths that are not profiles.
func panelProfileIsRootSource(expr ast.Expr, alias map[string]bool) bool {
	switch node := panelProfileUnparen(expr).(type) {
	case *ast.Ident:
		return alias[node.Name]
	case *ast.CallExpr:
		sel, ok := node.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil {
			return false
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok || pkg.Name != "os" {
			return false
		}
		switch sel.Sel.Name {
		case "TempDir", "UserHomeDir":
			return true
		case "Getenv":
			return panelProfileGetenvIsRoot(node)
		}
	}
	return false
}

// panelProfileGetenvIsRoot answers whether a Getenv call names a shared temp
// root. Only TEMP and TMP: those are the two names %TEMP% is spelled as, and a
// Getenv of anything else is not a root this ruler claims.
func panelProfileGetenvIsRoot(call *ast.CallExpr) bool {
	if len(call.Args) != 1 {
		return false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return false
	}
	switch strings.ToUpper(value) {
	case "TEMP", "TMP":
		return true
	}
	return false
}

// panelProfileDriveAbsolute answers whether a literal value is an absolute
// Windows path with a drive letter. A bare "C:" is drive-RELATIVE, not absolute,
// so it does not count.
func panelProfileDriveAbsolute(value string) bool {
	if len(value) < 3 || value[1] != ':' {
		return false
	}
	first := value[0]
	if (first < 'A' || first > 'Z') && (first < 'a' || first > 'z') {
		return false
	}
	return value[2] == '\\' || value[2] == '/'
}

// panelProfileCopyRoots copies a root set so one function's own aliases cannot
// leak into the next function.
func panelProfileCopyRoots(roots map[string]bool) map[string]bool {
	out := make(map[string]bool, len(roots))
	for name := range roots {
		out[name] = true
	}
	return out
}

// panelProfileQualifiedCallee answers "pkg.Name" for a selector call and "" for
// anything else. It is how strings.Join is told apart from filepath.Join.
func panelProfileQualifiedCallee(call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil {
		return ""
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return ""
	}
	return pkg.Name + "." + sel.Sel.Name
}

// panelProfileCalleeName answers the final name a call goes through, for both
// name(...) and recv.name(...). Anything else is "not a named call".
func panelProfileCalleeName(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name
	case *ast.SelectorExpr:
		if fun.Sel != nil {
			return fun.Sel.Name
		}
	}
	return ""
}

// panelProfileWindowCallIn answers the first window-creation call the source
// makes, if any. It walks the AST, so a comment claiming "never asks bringUp"
// cannot satisfy or violate it.
func panelProfileWindowCallIn(src panelProfileShapeSource) (string, int, bool) {
	var name string
	var line int
	found := false
	ast.Inspect(src.file, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		callee := panelProfileCalleeName(call)
		for _, want := range panelProfileWindowCalls {
			if callee == want {
				name, line, found = callee, src.fset.Position(call.Pos()).Line, true
				return false
			}
		}
		return true
	})
	return name, line, found
}

// panelProfileUnquote returns the leaf literal's value, fataling on anything
// that is not a plain string literal: a shape this ruler cannot name exactly is
// not a shape it will report approximately.
func panelProfileUnquote(t *testing.T, display, literal string) string {
	t.Helper()
	value, err := strconv.Unquote(literal)
	if err != nil {
		t.Fatalf("%s 里的字符串字面量 %s 读不出来：%v", display, literal, err)
	}
	return value
}

// Test314PanelProfileShapeRulerRingsOnTheFixedNameShape is the ruler's own
// positive control, in the repo and run on every pass: the exact mutant the
// audit ran must ring, the two launder shapes (alias, concatenation) must ring,
// and the shapes that are NOT the defect (per-run leaf, helper call, prose,
// quoted text, a test-local TempDir) must stay silent. Without this, a matcher
// that quietly stopped matching would print a clean scan forever while the
// exemption table still "covered" everything - the blind-guard the audit's 197
// reading asked for.
func Test314PanelProfileShapeRulerRingsOnTheFixedNameShape(t *testing.T) {
	cases := []struct {
		name       string
		source     string
		want       []string // expected fixed leaves; nil means the shape must stay silent
		lineNeedle string   // text on the line the red must name; needed when the leaf is folded
	}{
		{
			name: "the audit's revert, default tier verbatim",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() string {
	return filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")
}
`,
			want: []string{"wisp-253r6-panel-profile"},
		},
		{
			name: "the literal first, the root second",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() string {
	return filepath.Join("wisp-arg-order-profile", os.TempDir())
}
`,
			want: []string{"wisp-arg-order-profile"},
		},
		{
			name: "wrapped across lines",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() string {
	return filepath.Join(
		os.TempDir(),
		"wisp-wrapped-profile",
	)
}
`,
			want: []string{"wisp-wrapped-profile"},
		},
		{
			name: "local alias of os.TempDir",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() string {
	root := os.TempDir()
	return filepath.Join(root, "wisp-alias-laundered-profile")
}
`,
			want: []string{"wisp-alias-laundered-profile"},
		},
		{
			name: "package-level alias of os.TempDir",
			source: `package main

import (
	"os"
	"path/filepath"
)

var panelTempRoot = os.TempDir()

func mutant() string {
	return filepath.Join(panelTempRoot, "wisp-package-alias-profile")
}
`,
			want: []string{"wisp-package-alias-profile"},
		},
		{
			name: "concatenated instead of joined",
			source: `package main

import "os"

func mutant() string {
	return os.TempDir() + "/wisp-concatenated-profile"
}
`,
			want: []string{"/wisp-concatenated-profile"},
		},
		{
			name: "a profile-shaped leaf behind a helper return (the 319 holder shape)",
			source: `package main

import (
	"os"
	"path/filepath"
)

func panel319DataPath() string {
	return filepath.Join(os.TempDir(), "wisp-319r1-holder-profile")
}
`,
			want: []string{"wisp-319r1-holder-profile"},
		},
		{
			name: "a non-profile leaf handed to the panel sink",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-quicklook"))
}
`,
			want: []string{"wisp-quicklook"},
		},
		{
			name: "a local binding handed to the panel sink",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() *PanelManager {
	dataPath := filepath.Join(os.TempDir(), "wisp-bound")
	return NewPanelManager(nil, nil, dataPath)
}
`,
			want: []string{"wisp-bound"},
		},
		{
			name: "a non-profile fixed temp file for another purpose stays silent",
			source: `package main

import (
	"os"
	"path/filepath"
)

func fine() {
	_ = os.ReadFile(filepath.Join(os.TempDir(), "planted-key-320-temp.txt"))
}
`,
		},
		{
			name: "per-run leaf through fmt.Sprintf",
			source: `package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func fine() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d-%d", "wisp-253r6-panel-profile", os.Getpid(), 7))
}
`,
		},
		{
			name: "the panelProfileLeaf helper",
			source: `package main

func fine() string {
	return panelProfileLeaf("wisp-253r6-panel-profile")
}
`,
		},
		{
			name: "a test-local temp dir is per run by construction",
			source: `package main

import "path/filepath"

type dirs struct{}

func (dirs) TempDir() string { return "" }

func fine(d dirs) string {
	return filepath.Join(d.TempDir(), "wisp-not-a-shape")
}
`,
		},
		{
			name: "prose in a comment and text in a quoted string",
			source: "package main\n\n" +
				"// Old shape: filepath.Join(os.TempDir(), \"wisp-prose\") must not come back.\n" +
				"var panelShapeDoc = \"filepath.Join(os.TempDir(), \\\"wisp-in-a-string\\\")\"\n\n" +
				"func fine() {}\n",
		},
		// --- 2026-11-14 audit attack 2 (D1): a leaf assembled from pieces -------
		// The leaf is judged on the folded value, so each of these five is the
		// same shape as the single literal above it.
		{
			name: "two literal halves inside one Join",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-33r1"+"-panel-profile"))
}
`,
			want:       []string{"wisp-33r1-panel-profile"},
			lineNeedle: `"wisp-33r1"+"-panel-profile"`,
		},
		{
			name: "the whole path by concatenation",
			source: `package main

import "os"

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, os.TempDir()+"/wisp-33r1"+"-panel-profile")
}
`,
			want:       []string{"/wisp-33r1-panel-profile"},
			lineNeedle: `os.TempDir()+"/wisp-33r1"+"-panel-profile"`,
		},
		{
			name: "the leaf by strings.Join",
			source: `package main

import (
	"os"
	"path/filepath"
	"strings"
)

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		strings.Join([]string{"wisp-33r1", "-panel-profile"}, "")))
}
`,
			want:       []string{"wisp-33r1-panel-profile"},
			lineNeedle: `strings.Join([]string{"wisp-33r1", "-panel-profile"}, "")`,
		},
		{
			name: "a literal suffix outside the Join",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-33r1")+"-panel-profile")
}
`,
			want:       []string{"wisp-33r1-panel-profile"},
			lineNeedle: `filepath.Join(os.TempDir(), "wisp-33r1")+"-panel-profile"`,
		},
		{
			name: "a concatenated leaf bound to a local, then to the sink",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() *PanelManager {
	p := filepath.Join(os.TempDir(), "wisp-33r1") + "-panel-profile"
	return NewPanelManager(nil, nil, p)
}
`,
			want:       []string{"wisp-33r1-panel-profile"},
			lineNeedle: `p := filepath.Join(os.TempDir(), "wisp-33r1") + "-panel-profile"`,
		},
		{
			name: "a folded leaf with no temp root stays silent",
			source: `package main

import "path/filepath"

func fine() string {
	return filepath.Join("data", "wisp-"+"33r1-profile")
}
`,
		},
		{
			name: "a constant prefix plus a runtime suffix is not a fixed name",
			source: `package main

import (
	"os"
	"path/filepath"
	"strconv"
)

func fine() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-panel-"+strconv.Itoa(os.Getpid())))
}
`,
		},
		// --- 2026-11-14 audit attack 3 (D2): roots that are not os.TempDir() ---
		{
			name: "root os.Getenv(TEMP)",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.Getenv("TEMP"), "wisp-audit3a-panel-profile"))
}
`,
			want:       []string{"wisp-audit3a-panel-profile"},
			lineNeedle: `os.Getenv("TEMP")`,
		},
		{
			name: "root os.Getenv(TMP)",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.Getenv("TMP"), "wisp-audit3a2-panel-profile"))
}
`,
			want:       []string{"wisp-audit3a2-panel-profile"},
			lineNeedle: `os.Getenv("TMP")`,
		},
		{
			name: "a Getenv that is not TEMP or TMP is not a root",
			source: `package main

import (
	"os"
	"path/filepath"
)

func fine() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.Getenv("TMPDIR"), "wisp-quicklook"))
}
`,
		},
		{
			name: "root os.UserHomeDir through its two results",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "wisp-audit3b-panel-profile")
}
`,
			want:       []string{"wisp-audit3b-panel-profile"},
			lineNeedle: `filepath.Join(home, "wisp-audit3b-panel-profile")`,
		},
		{
			name: "a cache directory is not one of the roots this ruler claims",
			source: `package main

import (
	"os"
	"path/filepath"
)

func fine() *PanelManager {
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil
	}
	return NewPanelManager(nil, nil, filepath.Join(cache, "wisp-quicklook"))
}
`,
		},
		{
			name: "root is a drive-letter absolute literal",
			source: `package main

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, "C:\\Temp\\wisp-audit3c-panel-profile")
}
`,
			want:       []string{`C:\Temp\wisp-audit3c-panel-profile`},
			lineNeedle: `"C:\\Temp\\wisp-audit3c-panel-profile"`,
		},
		{
			name: "a drive-letter absolute literal joined with a leaf",
			source: `package main

import "path/filepath"

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join("C:\\Temp", "wisp-abs-joined-panel-profile"))
}
`,
			want:       []string{`C:\Temp\wisp-abs-joined-panel-profile`},
			lineNeedle: `filepath.Join("C:\\Temp", "wisp-abs-joined-panel-profile")`,
		},
		{
			name: "a relative literal root is not a root",
			source: `package main

import "path/filepath"

func fine() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join("temp-dir", "wisp-quicklook"))
}
`,
		},
		// --- 2026-11-14 audit X1/X2 (D3): package-level declarations ----------
		{
			name: "a package-level DataPath key is a sink where it is declared",
			source: `package main

import (
	"os"
	"path/filepath"
)

type auditSink struct {
	DataPath string
}

var auditPackageSink = auditSink{DataPath: filepath.Join(os.TempDir(), "wisp-auditpkg-sink")}
`,
			want:       []string{"wisp-auditpkg-sink"},
			lineNeedle: `var auditPackageSink = auditSink{DataPath: filepath.Join(os.TempDir(), "wisp-auditpkg-sink")}`,
		},
		{
			name: "a package-level fixed join handed to a sink in a function",
			source: `package main

import (
	"os"
	"path/filepath"
)

var auditPackageJoin = filepath.Join(os.TempDir(), "wisp-auditpkg-var")

func mutant() *PanelManager {
	return NewPanelManager(nil, nil, auditPackageJoin)
}
`,
			want:       []string{"wisp-auditpkg-var"},
			lineNeedle: `var auditPackageJoin = filepath.Join(os.TempDir(), "wisp-auditpkg-var")`,
		},
		{
			name: "a package-level fixed join that never reaches a sink stays silent",
			source: `package main

import (
	"os"
	"path/filepath"
)

var auditPackagePlain = filepath.Join(os.TempDir(), "wisp-auditpkg-plain")

func fine() string {
	return auditPackagePlain
}
`,
		},
		// --- D4: the three SINK arms that had no resident control -------------
		// The audit found these arms working on a one-off probe but uncovered by
		// this table, so a later break in them would have been all green.
		{
			name: "the newPanelControl sink arm",
			source: `package main

import (
	"os"
	"path/filepath"
)

func mutant() {
	_, _ = newPanelControl(webview2.WindowOptions{}, filepath.Join(os.TempDir(), "wisp-control-arm"))
}
`,
			want:       []string{"wisp-control-arm"},
			lineNeedle: `_, _ = newPanelControl(webview2.WindowOptions{}, filepath.Join(os.TempDir(), "wisp-control-arm"))`,
		},
		{
			name: "the DataPath field-assignment sink arm",
			source: `package main

import (
	"os"
	"path/filepath"
)

type auditConfig struct {
	DataPath string
}

func mutant(c *auditConfig) {
	c.DataPath = filepath.Join(os.TempDir(), "wisp-field-arm")
}
`,
			want:       []string{"wisp-field-arm"},
			lineNeedle: `c.DataPath = filepath.Join(os.TempDir(), "wisp-field-arm")`,
		},
		{
			name: "the DataPath key sink arm inside a function body",
			source: `package main

import (
	"os"
	"path/filepath"
)

type auditConfig struct {
	DataPath string
}

func mutant() auditConfig {
	return auditConfig{DataPath: filepath.Join(os.TempDir(), "wisp-key-arm")}
}
`,
			want:       []string{"wisp-key-arm"},
			lineNeedle: `return auditConfig{DataPath: filepath.Join(os.TempDir(), "wisp-key-arm")}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hits := panelProfileHitsInText(t, tc.name+".go", []byte(tc.source))
			got := make([]string, len(hits))
			for i, hit := range hits {
				got[i] = hit.leaf
			}
			sort.Strings(got)
			want := append([]string(nil), tc.want...)
			sort.Strings(want)
			if strings.Join(got, "|") != strings.Join(want, "|") {
				t.Fatalf("这把尺对这个形状的读数是 %q，应为 %q（叶）——正控与尺子必须同生同死", got, want)
			}
			for _, hit := range hits {
				if hit.line < 1 || hit.raw == "" || !strings.Contains(hit.code, hit.raw) {
					t.Errorf("命中 %q 的行号与源码行不对（line=%d raw=%q code=%q）：红句必须点名叫对文件与行，并印出叶的起始字面量",
						hit.leaf, hit.line, hit.raw, hit.code)
				}
			}
			if len(want) == 0 {
				return
			}
			// The named line is the first literal's line: the deployed literal itself
			// for a single-literal shape, an explicit needle for a folded one (a
			// folded value is not spelled on any one line, so it cannot be searched
			// for as a quoted string).
			needle := tc.lineNeedle
			if needle == "" {
				needle = `"` + tc.want[0] + `"`
			}
			named := hits[0].line
			if lineOf := panelProfileLineOf(t, tc.source, needle); named != lineOf {
				t.Errorf("%s 的点名行错了：红句写 :%d，叶字面量其实在 :%d", tc.name, named, lineOf)
			}
		})
	}

	// D5: the summary line itself is a control. The audit's M1 mutation (delete
	// the holder's justifying sentence) failed this test while the old summary
	// still printed "0 处红", which reads as green to the next person. The verdict
	// is now decided by one function, and this is what pins it.
	t.Run("the summary line cannot read green while anything is red", func(t *testing.T) {
		summaryCases := []struct {
			name                            string
			sources, hits, exempted, faults int
			clean                           bool
		}{
			{name: "clean run", sources: 274, hits: 2, exempted: 2, clean: true},
			{name: "one unexempted shape hit", sources: 274, hits: 3, exempted: 2, clean: false},
			{name: "only a stale exemption reason", sources: 274, hits: 2, exempted: 2, faults: 1, clean: false},
		}
		for _, sc := range summaryCases {
			line, clean := panelProfileSummaryLine(sc.sources, sc.hits, sc.exempted, sc.faults)
			if clean != sc.clean {
				t.Fatalf("%s：汇总行判定 clean=%v，应为 %v（%s）", sc.name, clean, sc.clean, line)
			}
			if !sc.clean && (!strings.Contains(line, "【红】") || strings.Contains(line, "【绿】")) {
				t.Errorf("%s：有红时汇总行必须写【红】、不得出现【绿】，实际：%s", sc.name, line)
			}
			if !strings.Contains(line, fmt.Sprintf("实红 %d 处", sc.hits-sc.exempted)) {
				t.Errorf("%s：汇总行必须印出实红数字，实际：%s", sc.name, line)
			}
			if sc.sources == 274 && !strings.Contains(line, "形状命中 2 处") && sc.hits == 2 {
				t.Errorf("%s：汇总行必须印出形状命中总数，实际：%s", sc.name, line)
			}
		}
	})
}

// panelProfileHitsInText runs the same census over an in-memory sample, so the
// matcher can be exercised without touching the tree.
func panelProfileHitsInText(t *testing.T, display string, src []byte) []panelProfileShapeHit {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, display, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("%s 的合成样本解析失败（正控自己坏了）：%v", display, err)
	}
	return panelProfileHitsInSource(t, panelProfileShapeSource{display: display, path: display, src: src, file: file, fset: fset})
}

// panelProfileLineOf answers the 1-based line of needle in a synthetic sample.
func panelProfileLineOf(t *testing.T, src, needle string) int {
	t.Helper()
	idx := strings.Index(src, needle)
	if idx < 0 {
		t.Fatalf("合成样本里找不到 %q（正控自己坏了）", needle)
	}
	return strings.Count(src[:idx], "\n") + 1
}
