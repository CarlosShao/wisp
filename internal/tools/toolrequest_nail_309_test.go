package tools

// 票 309 AC#2 —— "将来任一 caller 手搭 agent.ToolRequest 递进空 TaskID" 这件事要变成
// "有人新写构造点就会红"的**枚数钉**，⛔ 依赖任何未来接线。落点（309-r2 按编排者 §4 收紧后）＝
// **全 roots 之内生产构造点合计 == 1 枚，且具名那枚落在 internal/agent/loop.go**
// —— 旧钉问的是"internal/agent 之外 == 0"，那是同一件事的另一半，但它把 internal/agent 的
// 24 枚非测试源整包跳过（重试包装器最自然的落点就在那包里），所以"往里再加一条"⛔ 响。
//
// 两刀（判据⛔ 是"我记得的规矩"，是下面现跑的尺）：
//
//	(1) 射程：roots 从 cmd＋internal 扩到 cmd＋internal＋tools＋scripts；排除支从
//	    strings.Contains(rel,"internal/agent/") 换成"合计 == 1 ＋具名那枚在 loop.go"。
//	(2) matcher：旧形只认 *ast.CompositeLit 且 Type 是 agent.ToolRequest 那枚 SelectorExpr，
//	    于是另外两形今天绿（309-v1 实测各 rc=0＝假阴性）。现在认三形：字面量（限定名＋同包裸名）／
//	    可解析到 agent.ToolRequest 的 var 声明（声明＋逐字段赋值那一形，没有复合字面量）／
//	    同文件类型别名指向同一限定名之后再写字面量。
//
// roots 那两把尺（现量 HEAD 202ed4a5，2026-10-11 14:3x +08；R359-1 读数在
// .scratch/wisp/probes/309/r2/logs/roots-rulers.log）：
//   ① module 路径前缀尺＝`find . -name go.mod` 里三枚 nested 模块逐字
//     github.com/CarlosShao/wisp/tools/d22scan／.../tools/mockllm／.../scripts/spike
//     ⇒ 全落在 github.com/CarlosShao/wisp/ 那一截之内。Go 的 internal/ 可见性按 import 路径
//     前缀判定 ⇒ "独立模块进不来 internal/"⛔ 成立，tools 与 scripts 同判 ⇒ **两枚都进 roots**。
//   ② 活例尺＝`grep -rn "CarlosShao/wisp/internal" tools scripts --include=*.go`＝1 命中＝
//     tools/signmodels/main.go:34 逐字 "github.com/CarlosShao/wisp/internal/models"
//     ⇒ nested 模块伸手进 internal/ 在本仓已是事实。scripts 自己今天 0 命中，但①已经放行它，
//     ⛔ 另造一条"看谁此刻还没伸手"的规矩（那正是本票在防的那类"接线事实当结构保证"）。
//
// 分母（⛔ 一枚空名册当"扫过了"；先例逐字＝task_state_188_test.go 那句
// "a scan whose denominator is zero reports clean for a tree it never read"）：旧钉现量 245
// ＝cmd＋internal 的 269 枚非测试源减 internal/agent 的 24；这枚新钉⛔ 再排任何目录，
// 现量＝293 枚（cmd 40／internal 229／tools 13／scripts 11，第三把尺＝find <root> -name '*.go'
// -not -name '*_test.go' | wc -l，与钉自打印逐枚对得上，读数 R359-7），每枚 root 还各自断言
// "至少解析到一枚源"（root 名册为空那一支直接 t.Fatalf，⛔ 报绿）。
//
// 词面尺的现量（★这族数⛔ 是常数，引之前先重跑；本腿改完 bridge.go 与本件之后重测＝HEAD
// 1da192e9，2026-10-11 14:4x +08，工作树面。派单锚 202ed4a5 那一发＝69／66／3，见票 309
// 「编排者复跑＋裁定」§1 末行）：
//   `grep -rn 'ToolRequest{' internal cmd --include=*.go`＝总 71 行／文件名以 _test.go 结尾
//   69 行（**含本件自己的注释与测试源码**）／非测试 2 行；那 2 行＝internal/agent/loop.go
//   的真构造点＋bridge.go 里引用这把尺的那行注释自己，剥掉行首 // 之后＝1。
//   旧钉头抄的 60／59／1 是同尺在 10-11 11:5x 锚上的数；309-r1 那两枚件把它抬到 69／66／3
//   （编排者 14:1x 复跑），本腿再加 2：＋3（本件，词面尺段与控制码里的字面量）−1（bridge.go
//   少引一枚字面量），逐枚对账＝本件 7→10、bridge.go 2→1（尺＝git grep -c 同一 pattern 打在
//   202ed4a5 ⇄ grep -c 打在工作树，读数 R359-6）。⛔ 谁再把任何一组当常数抄下去。
//   tools＋scripts 两枚 root 今天 0 命中，所以把 roots 加进这把词面尺⛔ 动它（仍 71）。
//   ★还有一枚本腿自抓的口径洞：`| grep -v _test.go` 滤的是**整行内容**、⛔ 文件名字段 ⇒
//   一枚同时带 `ToolRequest{` 与 `_test.go` 两个字样的注释行会被这把尺自己吞掉（这腿先就被它
//   骗过一次，读出"非测试 1 枚"的假数；读数 R359-6 两把口径并排）。⇒ 结论照旧（生产构造点
//   1 枚），但那句"带尺重述"要按这两把口径读，⛔ 只抄一个数。
//
// ★这把尺**仍⛔ 是全覆盖**。下面这些写法今天抓不住，具名列出，⛔ 抄成"全覆盖"：
//   · 别名再套一层落到**别的包**（pkg q 里 type Y = agent.ToolRequest，第三方写 q.Y{...}）＝漏；
//   · 函数体内声明的别名（本 matcher 只读文件级 type 声明）＝漏；
//   · 改名 import（import ag "…/internal/agent" 后 ag.ToolRequest{...}）与 dot-import ＝漏；
//   · new(agent.ToolRequest)／make([]agent.ToolRequest, n)／反射或解码构造（json.Unmarshal 进结构体）＝漏；
//   · 定义型 type X agent.ToolRequest（⛔ 别名）再 agent.ToolRequest(x) 那枚转换＝漏；
//   · 省了类型的嵌套元素 []agent.ToolRequest{{TaskID: ""}} 的内层 {...}＝漏（外层那枚是 ArrayType，
//     同样⛔ 计）⇒ 这一形实测两向都计 0，见 .scratch/wisp/probes/309/r2/logs/nail-rig-5-elided.log；
//   · 裸名那一支按 **package 子句**放行⛔ 按目录放行 ⇒ 别处真有一枚叫 agent 的包时，它自己的
//     ToolRequest{} 会进分子（今天尺＝grep -rl '^package agent$' cmd internal tools scripts
//     全部落在 internal/agent/ 内，读数 R359-2）。
//   ⇒ 判据是"恒真形必须响"，⛔ 类型级完美（那要引 go/types，本票⛔ 装那个复杂度）。
//
// "加一枚即红"的六发 Rig（⛔ 一次跑完，每发当场还原 bridge.go，cmp＋sha256 留证；脚本与读数
// 都在 .scratch/wisp/probes/309/r2/logs/ 里，脚本＝nail-rigs.txt 与 nail-rigs-tree.txt）：
//   RIG-1 限定名字面量→红 rc=1｜RIG-2 声明＋逐字段赋值→红 rc=1｜RIG-3 同文件别名→红 rc=1
//   RIG-4b **internal/agent 里**再加第二条（同包裸名，导出树里打）→红 rc=1，名册具名那枚新件
//   RIG-5 省类型的嵌套元素→绿 rc=0（＝上面具名的漏形之一）｜RIG-6 别名落别的包再写 q.Mid{}
//   →绿 rc=0（＝漏形之二）。三形⛔ 响⛔ 算这一格交；两发绿⛔ 算"这枚钉没用"——它们钉的是
//   "本 matcher 的射程边界"，边界内的三形各有一发红。分母在每发里都打印且⛔ 漂（293/294/295
//   的差＝Rig 自己往树里加的那枚文件），红键的是分子。
//
//
// 范本⛔ 去 cmd/wisp 找，本包就有两枚现成的（paths_twocontainments_252_r2_test.go 的
// ReadDir+排除 _test.go+parser.ParseFile；task_state_188_test.go 的 WalkDir+把 examined
// 当返回值）。本格⛔ 新造一条生产可达路径：只读源码、不 import、不调用。

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// nail309Point is one production construction point the matcher found, named by file,
// line and the shape it was written in. Naming the shape is what lets the head comment
// claim exactly the coverage it has - and nothing more.
type nail309Point struct {
	file  string
	line  int
	shape string
}

// toolRequestNamesIn309 returns the identifiers that denote agent.ToolRequest inside this
// file: the bare name ToolRequest, but ONLY when the file's own package clause says
// "agent", plus every file-level type alias that resolves to it. Gating the bare name on
// the package clause is deliberate - package tools has callers that write their own
// unrelated identifiers, and a foreign type that happens to be called ToolRequest must
// not get into the numerator.
func toolRequestNamesIn309(file *ast.File) map[string]bool {
	names := map[string]bool{}
	pkgIsAgent := file.Name != nil && file.Name.Name == "agent"
	if pkgIsAgent {
		names["ToolRequest"] = true
	}
	// Aliases may chain within one file, so resolve to a fixpoint. Only `type X = T`
	// (ts.Assign is a valid position) is an alias; `type X agent.ToolRequest` creates a
	// NEW type that reaches agent.ToolRequest only through a conversion, and that shape is
	// named in the head comment as a residual blind spot rather than silently counted.
	for changed := true; changed; {
		changed = false
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Assign == token.NoPos || ts.Name == nil {
					continue
				}
				if denotesToolRequest309(ts.Type, names) && !names[ts.Name.Name] {
					names[ts.Name.Name] = true
					changed = true
				}
			}
		}
	}
	return names
}

// denotesToolRequest309 reports whether expr spells agent.ToolRequest - either the
// qualified selector or a name this file resolved to it (bare same-package name or alias).
func denotesToolRequest309(expr ast.Expr, names map[string]bool) bool {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		pkg, ok := e.X.(*ast.Ident)
		return ok && pkg.Name == "agent" && e.Sel.Name == "ToolRequest"
	case *ast.Ident:
		return names[e.Name]
	default:
		return false
	}
}

// literalShape309 labels which literal shape was counted, so the failure message and the
// head comment agree about what the matcher actually recognised.
func literalShape309(expr ast.Expr, pkgIsAgent bool) string {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		return "qualified-literal"
	case *ast.Ident:
		if pkgIsAgent && t.Name == "ToolRequest" {
			return "same-package-bare-literal"
		}
		return "alias-literal"
	default:
		return "literal"
	}
}

// findToolRequestPoints309 returns the construction points in one parsed file. Two shapes
// count: a composite literal whose type denotes agent.ToolRequest, and a var declaration
// that gives a variable that type with NO initializer - the "var req agent.ToolRequest"
// followed by "req.TaskID = ”" shape, which creates a zero-valued struct (empty TaskID
// included) without ever writing a literal. `var req = agent.ToolRequest{...}` is counted
// once, by its literal, so this number means "construction points", not "syntax nodes".
func findToolRequestPoints309(fset *token.FileSet, file *ast.File) []nail309Point {
	names := toolRequestNamesIn309(file)
	pkgIsAgent := file.Name != nil && file.Name.Name == "agent"
	var out []nail309Point
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CompositeLit:
			if denotesToolRequest309(n.Type, names) {
				pos := fset.Position(n.Pos())
				out = append(out, nail309Point{
					file:  pos.Filename,
					line:  pos.Line,
					shape: literalShape309(n.Type, pkgIsAgent),
				})
			}
		case *ast.GenDecl:
			if n.Tok != token.VAR {
				return true
			}
			for _, spec := range n.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || len(vs.Values) > 0 || vs.Type == nil {
					continue
				}
				if denotesToolRequest309(vs.Type, names) {
					pos := fset.Position(vs.Pos())
					out = append(out, nail309Point{
						file:  pos.Filename,
						line:  pos.Line,
						shape: "var-decl-no-initializer",
					})
				}
			}
		}
		return true
	})
	return out
}

// moduleRel309 turns a path the walk produced (../../internal/agent/loop.go) into a
// module-relative slash path so the named-file assertion holds whether the test runs from
// the worktree or from an out-of-repo export tree.
func moduleRel309(p string) string {
	rel := filepath.ToSlash(p)
	for strings.HasPrefix(rel, "../") {
		rel = rel[3:]
	}
	return rel
}

// scanToolRequestNail309 walks every non-test Go source under the given roots and returns
// the construction points it found, the denominator (how many files it actually parsed)
// and the per-root split of that denominator. The per-root split is what stops one root
// that quietly stopped existing - or quietly holds nothing - from being read as
// "scanned and clean".
func scanToolRequestNail309(t *testing.T, roots []string) (points []nail309Point, examined int, perRoot map[string]int) {
	t.Helper()
	perRoot = map[string]int{}
	for _, root := range roots {
		before := examined
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			fset := token.NewFileSet()
			file, perr := parser.ParseFile(fset, path, nil, 0)
			if perr != nil {
				t.Fatalf("parse %s: %v", path, perr)
			}
			examined++
			points = append(points, findToolRequestPoints309(fset, file)...)
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
		perRoot[moduleRel309(root)] = examined - before
		if examined == before {
			t.Fatalf("root %q 一枚非测试源也没解析：空名册等于扫过了，这一格没交", root)
		}
	}
	return points, examined, perRoot
}

// roster309 renders the found points as a sorted, named list - the message a reader needs
// in order to re-run the ruler, not a count on its own.
func roster309(points []nail309Point) string {
	names := make([]string, 0, len(points))
	for _, p := range points {
		names = append(names, fmt.Sprintf("%s:%d(%s)", moduleRel309(p.file), p.line, p.shape))
	}
	sort.Strings(names)
	if len(names) == 0 {
		return "(空名册)"
	}
	return strings.Join(names, " | ")
}

func perRootString309(perRoot map[string]int) string {
	keys := make([]string, 0, len(perRoot))
	for k := range perRoot {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s=%d", k, perRoot[k]))
	}
	return strings.Join(out, " ")
}

func Test309ToolRequestProductionConstructionNail(t *testing.T) {
	// go test runs with the package dir (internal/tools) as cwd; the module root is two
	// levels up. Scanned roster is named in every failure so the next reader can re-run.
	roots := []string{
		filepath.Join("..", "..", "cmd"),
		filepath.Join("..", "..", "internal"),
		filepath.Join("..", "..", "tools"),
		filepath.Join("..", "..", "scripts"),
	}

	points, examined, perRoot := scanToolRequestNail309(t, roots)

	// The nail must carry its own denominator. 200 is a FLOOR, not a pin: today's reading
	// is far above it, and the floor only rings when the scan stops reading the tree (a
	// narrowed root list, a moved directory, an archive that exported nothing).
	if examined < 200 {
		t.Fatalf("枚数钉分母不可信：只解析了 %d 枚非测试生产源（roots=%q，每 root=%s）。"+
			"分母接近零＝这枚扫根本没读到树，报绿是假绿",
			examined, roots, perRootString309(perRoot))
	}

	// The pin: exactly ONE production construction point across every root, and it is
	// named. Before 309-r2 this was "0 outside internal/agent" with the whole agent
	// package skipped, so a second point added inside that package stayed green.
	const (
		wantTotal      = 1
		wantOwner      = "internal/agent/loop.go"
		ownerRationale = "the one newTaskID() dispatch, id always 36 chars, never empty"
	)
	single := len(points) == wantTotal && strings.HasSuffix(moduleRel309(points[0].file), wantOwner)
	if !single {
		t.Errorf("枚数钉响：roots=%q 里数到 %d 枚 agent.ToolRequest 生产构造点，票 309 钉的是 %d 枚且那枚落在 %s（%s）。"+
			"名册=%s。分母=%d 枚非测试源（每 root=%s）。"+
			"新写一条构造点的人：要么让它接上非空 id，要么把它放进 309 AC#1 那枚仪器覆盖的入口——"+
			"空 TaskID 走到 assessorFor 会静默摘掉 R4 探测器，见 assessorfor_taint_instrument_309_test.go 的 ②",
			roots, len(points), wantTotal, wantOwner, ownerRationale,
			roster309(points), examined, perRootString309(perRoot))
	}

	// Controls, both directions: the matcher must count each shape it claims to cover, and
	// must NOT count a foreign type that merely shares the name. In-memory only - nothing
	// touches disk, no production path is added. If the ruler went blind (or turned into a
	// catch-all) it rings here instead of reporting a silent green.
	controls := []struct {
		name string
		want int
		src  string
	}{
		{"shape 1 qualified literal", 1, "package p\n" +
			"import \"github.com/CarlosShao/wisp/internal/agent\"\n" +
			"func f() { _ = agent.ToolRequest{TaskID: \"\"} }\n"},
		{"shape 2 var decl plus field assignment", 1, "package p\n" +
			"import \"github.com/CarlosShao/wisp/internal/agent\"\n" +
			"var rig agent.ToolRequest\n" +
			"func f() { rig.TaskID = \"\"; _ = rig }\n"},
		{"shape 3 type alias then literal", 1, "package p\n" +
			"import \"github.com/CarlosShao/wisp/internal/agent\"\n" +
			"type T = agent.ToolRequest\n" +
			"func f() { _ = T{TaskID: \"\"} }\n"},
		{"shape 1b same-package bare literal", 1, "package agent\n" +
			"func f() { req := ToolRequest{TaskID: \"\"}; _ = req }\n"},
		{"guard: bare literal in a NON-agent package", 0, "package p\n" +
			"type ToolRequest struct{ TaskID string }\n" +
			"func f() { _ = ToolRequest{TaskID: \"\"} }\n"},
		{"no double count: var initialised by a literal", 1, "package p\n" +
			"import \"github.com/CarlosShao/wisp/internal/agent\"\n" +
			"var rig = agent.ToolRequest{TaskID: \"\"}\n"},
	}
	for _, c := range controls {
		fset := token.NewFileSet()
		file, perr := parser.ParseFile(fset, "ctrl.go", c.src, 0)
		if perr != nil {
			t.Fatalf("control %q: parse: %v", c.name, perr)
		}
		if got := len(findToolRequestPoints309(fset, file)); got != c.want {
			t.Fatalf("control %q: matcher counted %d, want %d — the ruler is blind or a catch-all, "+
				"so its green over the scanned tree means nothing", c.name, got, c.want)
		}
	}

	t.Logf("AC#2 枚数钉：roots=%q 里 agent.ToolRequest 生产构造点=%d（==1，那枚＝%s），分母=%d 枚非测试源（每 root=%s）",
		roots, len(points), wantOwner, examined, perRootString309(perRoot))
}
