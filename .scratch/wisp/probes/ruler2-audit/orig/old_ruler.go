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
// THE SHAPES IT CALLS FIXED - two rules, unioned (panelProfileHitsInSource):
//   - NAME: a Join carrying os.TempDir() together with a string-literal leaf
//     that spells a profile folder ("...profile"), in either argument order,
//     across wrapped lines, through a local or package-level alias of
//     os.TempDir() (`root := os.TempDir()` then Join(root, "...")), or as
//     `os.TempDir() + "/literal"` concatenation. This is the audit's exact
//     shape, caught wherever it appears - sink or no sink - which is what keeps
//     a leaf returned from a helper (the 319 holder fixture) inside the ruler;
//   - SINK: the same join with ANY literal leaf once the value reaches a WebView2
//     data-path sink - NewPanelManager, newPanelControl, a DataPath field, a
//     DataPath key in a struct literal - directly or through a local binding.
//     Without it, a fixed name that avoids the word "profile" would be invisible.
// A fixed temp leaf for another purpose stays silent on purpose: ticket 320's
// planted-key control file in os.TempDir() is a regular file and a detector
// control, not a browser profile folder, and accusing it would be wrong.
//
// WHAT IT STILL CANNOT SEE - named, not hidden, the same honesty the 197 publish
// ruler's audit asked for (today-four-fixes-audit section 3 D):
//  1. a helper that joins a literal with os.TempDir() behind a parameter
//     (fixedProfile("lit") with the root spelled inside the helper): the root
//     and the literal never meet in one expression, and the SINK rule stops at a
//     function boundary too (that is how a fixed name avoids both rules);
//  2. a leaf assembled from literals ("wisp-" + "319r1-..."), or faked with a
//     placeholder that is not a per-run pid/nanotime;
//  3. a root that is not os.TempDir() at all (os.UserCacheDir, %TEMP% through
//     os.Getenv, an absolute literal). Those are different shapes; the ticket's
//     ruling names os.TempDir().
// Moving a fixed name into one of those three is a new decision to argue on the
// ticket, not an edit this ruler silently blesses.
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
	leaf    string // the fixed leaf literal, unquoted
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

	panelProfileCheckExemptions(t, hits, sources)

	t.Logf("314 形状尺：扫了 %d 个源文件，%d 处固定名命中（%d 处具名豁免、%d 处红）",
		len(sources), len(hits), exempted, len(hits)-exempted)
}

// panelProfileFixedNameRed is the one sentence this ruler exists to say. It
// names the file and line, prints the shape, says what breaks when two runs
// share the folder, and quotes the audit reading that proved nothing else in the
// repo can stand in for it.
func panelProfileFixedNameRed(t *testing.T, hit panelProfileShapeHit) {
	t.Helper()
	t.Errorf("%s:%d 命中固定名面板 profile 形状（叶 %q）：\n    %s\n"+
		"这个形状会让两发并发互踩同一个 WebView2 user-data 夹：一个夹就是一个浏览器 profile，"+
		"两个测试进程在同一条夹上各起一套顶层 browser host —— 先到的 CreateEnvironment 争同一夹会回 S_OK+NULL environment，"+
		"后到的 environment 会 JOIN 一个正在退场的浏览器进程（internal/ball/orb_windows.go:432-458 的两条守卫）；"+
		"票 314 的测量把同一夹同 tick 两个顶层 host 抓成过约 1.2 秒的共存读数。\n"+
		"ticket 314 的非实现者审计量过这形的牙是零：把默认档 panel_unbound_253r6_windows_test.go:71 改回固定名，"+
		"go test -c rc=0、Test253r6 PASS —— 仓里没有断言读叶名唯一性，所以只有这把源码形状尺能让这种回退红。\n"+
		"改法＝panelProfileLeaf(%q)（cmd/wisp/panel_host_windows_test.go 的既有 helper，叶＝prefix-pid-纳秒）；"+
		"确有理由共用同一个夹的，在 panelProfileLeafExemptions 里按「文件＋叶名」具名登记，并把源码里那句理由逐字抄进 quotes。",
		hit.display, hit.line, hit.leaf, hit.code, hit.leaf)
}

// panelProfileCheckExemptions runs the three checks every exemption owes: it is
// live, its quoted reason is still on disk, and a noWindow claim is still true.
func panelProfileCheckExemptions(t *testing.T, hits []panelProfileShapeHit, sources []panelProfileShapeSource) {
	t.Helper()
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
			t.Errorf("豁免表过期：%s 的固定名叶 %q 在这次扫描里覆盖 0 处。留着它就是把一张空白支票挂在尺子上："+
				"删掉这条，或者把叶名改回它真正豁免的那一处。", ex.file, ex.leaf)
		}
		src, ok := byDisplay[ex.file]
		if !ok {
			t.Errorf("豁免表点名了扫描集之外的源文件 %s：扫描集与豁免表必须同一次改动（否则这条豁免只活在纸上）", ex.file)
			continue
		}
		for _, quote := range ex.quotes {
			if !strings.Contains(string(src.src), quote) {
				t.Errorf("豁免 %s 的理由句在源码里找不到了：%q\n"+
					"这条豁免的正当性是从那段注释读出来的；注释改掉后它还留着，就成了一张没有出处的通行证。"+
					"重读源码、逐字重抄一句，或删掉这条豁免。", ex.file, quote)
			}
		}
		if !ex.noWindow {
			continue
		}
		if name, line, ok := panelProfileWindowCallIn(src); ok {
			t.Errorf("豁免 %s 自述「这里不建窗」，但 %s:%d 调用了 %s：这句话现在是假的，"+
				"那处固定名夹也就重新变成两发互踩的碰撞面。要么把固定名改成 panelProfileLeaf，要么删掉这条 noWindow 豁免并重新论证。",
				ex.file, ex.file, line, name)
		}
	}
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
//   - the NAME rule: a join carrying os.TempDir() and a leaf that calls itself a
//     profile folder ("...profile") is the audit's shape wherever it appears,
//     sink or no sink. This is what keeps a leaf that is returned from a helper
//     (the 319 holder fixture) inside the ruler instead of outside it;
//   - the SINK rule: a join carrying os.TempDir() and ANY literal leaf once its
//     value reaches a WebView2 data-path sink (NewPanelManager / newPanelControl,
//     a DataPath field, a DataPath key in a struct literal), directly or through
//     a local binding. It keeps a fixed name that avoids the word "profile" from
//     being invisible just because of its spelling.
//
// A fixed temp leaf for another purpose is silent on purpose. Ticket 320's
// planted-key control file in os.TempDir() is the measured one: it is a regular
// file and a detector control, not a browser profile folder, and a red that
// accused it would be a false accusation.
func panelProfileHitsInSource(t *testing.T, src panelProfileShapeSource) []panelProfileShapeHit {
	t.Helper()
	lines := strings.Split(string(src.src), "\n")
	packageRoots := panelProfilePackageTempRoots(src.file)
	var hits []panelProfileShapeHit
	seen := map[string]bool{}
	// The line is the LEAF literal's line, not the call's: on a wrapped Join the
	// literal is the thing a reader has to change, and this keeps the printed
	// source line containing the leaf it names.
	add := func(pos token.Pos, leaf string) {
		line := src.fset.Position(pos).Line
		key := strconv.Itoa(line) + "|" + leaf
		if seen[key] {
			return
		}
		seen[key] = true
		code := ""
		if line >= 1 && line <= len(lines) {
			code = strings.TrimSpace(lines[line-1])
		}
		hits = append(hits, panelProfileShapeHit{display: src.display, line: line, leaf: leaf, code: code})
	}
	// Package-level declarations: the NAME rule only, because a package-level
	// value has no sink in sight where it is declared.
	for _, decl := range src.file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			panelProfileNameRule(t, gen, src.display, packageRoots, add)
		}
	}
	// Function bodies: each one gets its own alias scope, so one function's
	// `tempRoot` cannot make another function's unrelated join look like a temp
	// root. Closures inside a body share that body's scope (over-approximate:
	// an extra red is possible, a missed one is not).
	for _, decl := range src.file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		roots := make(map[string]bool, len(packageRoots))
		for name := range packageRoots {
			roots[name] = true
		}
		panelProfileTempRoots(fn.Body, roots)
		panelProfileNameRule(t, fn.Body, src.display, roots, add)
		panelProfileSinkRule(t, fn.Body, src.display, roots, add)
	}
	return hits
}

// panelProfileNameRule is the shape the ticket is named after: os.TempDir()
// joined with a leaf that spells a profile folder.
func panelProfileNameRule(t *testing.T, node ast.Node, display string, roots map[string]bool, add func(token.Pos, string)) {
	t.Helper()
	ast.Inspect(node, func(n ast.Node) bool {
		fixed, ok := panelProfileFixedJoinIn(t, display, n, roots)
		if ok && strings.Contains(strings.ToLower(fixed.leaf), "profile") {
			add(fixed.pos, fixed.leaf)
		}
		return true
	})
}

// panelProfileSinkRule catches a fixed join with ANY leaf once its value is
// handed to a WebView2 data-path sink - directly, through a local binding, or
// through a DataPath field/key. A helper return is the boundary this rule cannot
// cross (header, escape surface 1).
func panelProfileSinkRule(t *testing.T, node ast.Node, display string, roots map[string]bool, add func(token.Pos, string)) {
	t.Helper()
	stored := map[string]panelProfileFixedJoin{}
	hits := func(expr ast.Expr) {
		if fixed, ok := panelProfileFixedJoinIn(t, display, expr, roots); ok {
			add(fixed.pos, fixed.leaf)
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
				add(fixed.pos, fixed.leaf)
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
						add(fixed.pos, fixed.leaf)
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

// panelProfileFixedJoin answers the fixed leaf a join carries when the temp root
// and a string literal meet in one expression: Join(...os.TempDir()..., "leaf"),
// the same through an alias of os.TempDir(), and `os.TempDir() + "/leaf"`.
func panelProfileFixedJoinIn(t *testing.T, display string, expr ast.Node, roots map[string]bool) (panelProfileFixedJoin, bool) {
	t.Helper()
	switch node := expr.(type) {
	case *ast.CallExpr:
		if panelProfileCalleeName(node) != "Join" {
			return panelProfileFixedJoin{}, false
		}
		var fixed panelProfileFixedJoin
		rooted := false
		for _, arg := range node.Args {
			if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				if fixed.leaf == "" {
					fixed.leaf = panelProfileUnquote(t, display, lit.Value)
					fixed.pos = lit.Pos()
				}
				continue
			}
			if panelProfileIsTempRoot(arg, roots) {
				rooted = true
			}
		}
		if rooted && fixed.leaf != "" {
			return fixed, true
		}
	case *ast.BinaryExpr:
		if node.Op != token.ADD {
			return panelProfileFixedJoin{}, false
		}
		if lit, ok := node.Y.(*ast.BasicLit); ok && lit.Kind == token.STRING && panelProfileIsTempRoot(node.X, roots) {
			return panelProfileFixedJoin{leaf: panelProfileUnquote(t, display, lit.Value), pos: lit.Pos()}, true
		}
		if lit, ok := node.X.(*ast.BasicLit); ok && lit.Kind == token.STRING && panelProfileIsTempRoot(node.Y, roots) {
			return panelProfileFixedJoin{leaf: panelProfileUnquote(t, display, lit.Value), pos: lit.Pos()}, true
		}
	}
	return panelProfileFixedJoin{}, false
}

// panelProfileFixedJoin is one fixed leaf: what it names and where it is written.
type panelProfileFixedJoin struct {
	leaf string
	pos  token.Pos
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

// panelProfilePackageTempRoots collects the package-level names bound exactly to
// os.TempDir(). Function bodies are deliberately NOT walked here: their aliases
// belong to their own scope (see panelProfileHitsInSource).
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

// panelProfileTempRoots adds to roots every name assigned exactly os.TempDir()
// within node (a block, a GenDecl, or a whole function body).
func panelProfileTempRoots(node ast.Node, roots map[string]bool) map[string]bool {
	if roots == nil {
		roots = map[string]bool{}
	}
	record := func(lhs, rhs []ast.Expr) {
		if len(lhs) != 1 || len(rhs) != 1 {
			return
		}
		name, ok := lhs[0].(*ast.Ident)
		if !ok {
			return
		}
		if panelProfileIsTempRoot(rhs[0], roots) {
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

// panelProfileIsTempRoot reports whether expr spells os.TempDir() (or a name
// bound to it) through any number of parentheses.
func panelProfileIsTempRoot(expr ast.Expr, alias map[string]bool) bool {
	switch e := expr.(type) {
	case *ast.CallExpr:
		sel, ok := e.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil || sel.Sel.Name != "TempDir" {
			return false
		}
		pkg, ok := sel.X.(*ast.Ident)
		return ok && pkg.Name == "os"
	case *ast.Ident:
		return alias[e.Name]
	case *ast.ParenExpr:
		return panelProfileIsTempRoot(e.X, alias)
	}
	return false
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
		name   string
		source string
		want   []string // expected fixed leaves; nil means the shape must stay silent
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
				if hit.line < 1 || !strings.Contains(hit.code, hit.leaf) {
					t.Errorf("命中 %q 的行号与源码行不对（line=%d code=%q）：红句必须点名叫对文件与行",
						hit.leaf, hit.line, hit.code)
				}
			}
			if len(want) == 0 {
				return
			}
			// The first case is the audit's own mutant; its line must be the
			// literal's line, which is the needle's line in the synthetic file.
			needle := `"` + tc.want[0] + `"`
			if got := hits[0].line; got != panelProfileLineOf(t, tc.source, needle) {
				t.Errorf("%s 的点名行错了：红句写 :%d，叶字面量其实在 :%d", tc.name, got, panelProfileLineOf(t, tc.source, needle))
			}
		})
	}
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
