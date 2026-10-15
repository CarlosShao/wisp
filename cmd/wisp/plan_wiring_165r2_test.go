package main

// 票 165-r2（装配半格）· 尺 a：调用点枚数尺。
//
// 这一族本仓有名字："写了没接线"。一个没人调用的导出函数照样 `go build` 干净、
// `go vet` 干净，所以 rc=0 证明不了接线。165-r1 交完时编排者现量
// "grep -rn BuiltinPlanEntries --include=*.go internal cmd | grep -v _test"
// 只命中定义自己（internal/tools/plan_present.go:134 注释＋:138 func），零枚生产
// 调用者。产码那一处由本腿补在 cmd/wisp/run.go，形状照同一循环里已接的先例
// BuiltinSubagentEntries；这枚尺量的就是"那一处在不在、形状对不对"：
//
//	射程＝internal/ 与 cmd/ 的全部 .go，剥掉 _test.go（与上面那把 grep 尺同射程）；
//	判据＝tools.BuiltinPlanEntries 作为*调用表达式*出现 1 处，且那一处必须是某个
//	for range 的 RangeExpr、循环体里有 Register 调用——"读了名册却没注册进
//	registry" 在这一族里同样算没接，所以形状也进判据，不只枚数。
//	阳性对照＝同一把尺先数已接的先例 BuiltinSubagentEntries＝1 枚且形状同形；
//	先例若也数出 0，说明尺瞎了（或 AST 读法变了），plan 那枚的读数就不作数。
//
// ⛔ 拿 internal/tools 测试自己 append 的那份名册（plan_present_165_test.go:76）
// 当"已接"凭据：那是测试拼的。装配可达那一半（真 assembleRuntime＋真 bridge＋
// 真确认卡）在同批的 plan_assembled_165r2_test.go，两枚件都受"每件 ≤260 行"约束。

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const (
	planToolName165r2     = "plan.present"
	planBuilder165r2      = "BuiltinPlanEntries"
	subagentBuilder165r2  = "BuiltinSubagentEntries"
	planAssemblyFile165r2 = "cmd/wisp/run.go"
)

// repoRoot165r2 is the hop from where a go test binary for cmd/wisp actually
// starts (the package directory) to the repository root the roster names files
// against.
var repoRoot165r2 = filepath.Join("..", "..")

type callSite165r2 struct {
	file      string
	line      int
	inRange   bool // the call is the RangeExpr of a for loop
	registers bool // that loop's body calls Register
}

func (c callSite165r2) String() string {
	return fmt.Sprintf("%s:%d range=%v register=%v", c.file, c.line, c.inRange, c.registers)
}

// callSites165r2 walks the production sources under roots and returns every
// `tools.<fn>(...)` call expression it finds, sorted by file:line. Paths are
// reported relative to the repository root, because a go test binary starts in
// the package directory (cmd/wisp), not the repo root: 本腿第一次跑就在那一枚
// 相对根上拿到 `GetFileAttributesEx internal: The system cannot find the file
// specified.`，那是台面故障、不是判据（已具名落在 probes/165/r2 的台面一节）。
func callSites165r2(t *testing.T, roots []string, fn string) []callSite165r2 {
	t.Helper()
	var out []callSite165r2
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			t.Fatalf("台面故障（不是判据）：扫描根 %q 读不到：%v", root, err)
		}
		err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			slash := filepath.ToSlash(p)
			if d.IsDir() || !strings.HasSuffix(slash, ".go") || strings.HasSuffix(slash, "_test.go") ||
				strings.Contains(slash, "/testdata/") {
				return nil
			}
			out = append(out, fileCallSites165r2(t, p, repoRel165r2(t, p), fn)...)
			return nil
		})
		if err != nil {
			t.Fatalf("台面故障（不是判据）：扫描 %s 失败：%v", root, err)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].file != out[j].file {
			return out[i].file < out[j].file
		}
		return out[i].line < out[j].line
	})
	return out
}

// repoRel165r2 renders one scanned path the way the roster names it: relative to
// the repository root, forward-slashed, so `cmd/wisp/run.go` compares equal no
// matter which volume or separator this bench runs on.
func repoRel165r2(t *testing.T, p string) string {
	t.Helper()
	rel, err := filepath.Rel(repoRoot165r2, p)
	if err != nil {
		t.Fatalf("台面故障（不是判据）：%s 相对仓库根 %s 算不出来：%v", p, repoRoot165r2, err)
	}
	return filepath.ToSlash(rel)
}

// fileCallSites165r2 parses one production file (path, resolved from the test's
// own working directory) and reports that builder's call sites under rel (the
// name the roster answers with), annotated with whether each sits in the
// "for range + Register" shape.
func fileCallSites165r2(t *testing.T, path, rel, fn string) []callSite165r2 {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("台面故障（不是判据）：%s 解析失败：%v", filepath.ToSlash(path), err)
	}
	// Pass one: which call expressions are a for-range's source, and does that
	// loop body register anything.
	inLoop := map[*ast.CallExpr]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		rs, ok := n.(*ast.RangeStmt)
		if !ok {
			return true
		}
		if call, ok := rs.X.(*ast.CallExpr); ok {
			inLoop[call] = loopRegisters165r2(rs.Body)
		}
		return true
	})
	// Pass two: the builder calls themselves.
	var out []callSite165r2
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != fn {
			return true
		}
		if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "tools" {
			return true
		}
		registers, ranged := inLoop[call]
		out = append(out, callSite165r2{
			file: rel, line: fset.Position(call.Pos()).Line,
			inRange: ranged, registers: registers && ranged,
		})
		return true
	})
	return out
}

func loopRegisters165r2(body *ast.BlockStmt) bool {
	if body == nil {
		return false
	}
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Register" {
			found = true
		}
		return true
	})
	return found
}

func Test165r2PlanWiringCallSiteCensus(t *testing.T) {
	roots := []string{filepath.Join(repoRoot165r2, "internal"), filepath.Join(repoRoot165r2, "cmd")}
	precedent := callSites165r2(t, roots, subagentBuilder165r2)
	if len(precedent) != 1 {
		t.Fatalf("阳性对照失败：先例 %s 的产码调用点数出 %d 枚（要 1 枚，run.go 的 subagent 块）"+
			"⇒ 这把尺此刻读不到任何装配，plan 那枚的读数不作数：%v",
			subagentBuilder165r2, len(precedent), precedent)
	}
	if !precedent[0].inRange || !precedent[0].registers {
		t.Fatalf("阳性对照失败：先例那枚 %s 不是「for range 该调用 + 循环体调 Register」那一形"+
			"⇒ 形状判据本身要重做", precedent[0])
	}
	t.Logf("阳性对照成立：%s 是已接形状（range＋Register）⇒ 尺读得出「接上了」这一形", precedent[0])

	sites := callSites165r2(t, roots, planBuilder165r2)
	if len(sites) != 1 {
		t.Errorf("接线 RED：tools.%s 在 internal/ 与 cmd/ 的产码里调用点 %d 枚，要 1 枚"+
			"（%s 的装配那一处）；命中清单＝%v（0 枚＝写了没接线那一族，"+
			"plan.present 在生产装配里根本够不到）",
			planBuilder165r2, len(sites), planAssemblyFile165r2, sites)
		return
	}
	s := sites[0]
	if s.file != planAssemblyFile165r2 {
		t.Errorf("接线形状不合规：调用点落在 %s，装配面只认 %s 那一处", s.file, planAssemblyFile165r2)
	}
	if !s.inRange || !s.registers {
		t.Errorf("接线形状不合规：%s 不是「for range 该调用 + 循环体调 Register」那一形"+
			"（先例 subagent 块是）：读了名册却没把它的工具注册进 registry，等于没接", s)
	}
	t.Logf("接线读数：%s ⇒ 产码调用点 1 枚，形状与先例同形", s)
}
