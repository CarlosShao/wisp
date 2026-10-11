package tools

// 票 309 AC#2 —— "将来任一 caller 手搭 agent.ToolRequest{ 递进空 TaskID" 这件事要变成
// "有人新写构造点就会红"的**枚数钉**，⛔ 依赖任何未来接线。落点＝ internal/agent 之外的
// 生产构造点枚数 == 0（分母⛔ 零：a scan whose denominator is zero reports clean for a
// tree it never read —— 抄自 task_state_188_test.go 的既有规矩，本枚自带分母）。
//
// 现量基线（两把独立尺，写在失败句里让下一位能重跑）：
//   枚数尺＝`grep -rn 'ToolRequest{' internal cmd --include=*.go`：全仓 60 枚，_test.go 内
//     59，非测试仅 1 枚＝internal/agent/loop.go:675 `req := ToolRequest{`（裸名，同包，
//     本枚按限定名 agent.ToolRequest 匹配本就⛔ 命中它，且它落在 internal/agent 排除项里）。
//   门尺＝票 154 .scratch/wisp/probes/154/gate-clauses.sh 的 G1 那格登记 want_n 0，今日对得上。
// ⇒ internal/agent 之外的生产构造点＝**0**。
//
// 范本⛔ 去 cmd/wisp 找，本包就有两枚现成的（paths_twocontainments_252_r2_test.go 的
// ReadDir+排除 _test.go+parser.ParseFile；task_state_188_test.go 的 WalkDir+把 examined
// 当返回值）。本格⛔ 新造一条生产可达路径：只读源码、不 import、不调用。

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// countToolRequestLits parses one file and returns how many `agent.ToolRequest{...}`
// composite literals it holds. Matching is by the qualified type selector, never by
// a helper name or a string grep, so a renamed variable keeps this ruler honest while
// a genuinely new construction point cannot hide.
func countToolRequestLits(fset *token.FileSet, file *ast.File) int {
	n := 0
	ast.Inspect(file, func(node ast.Node) bool {
		cl, ok := node.(*ast.CompositeLit)
		if !ok {
			return true
		}
		sel, ok := cl.Type.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pkg, ok := sel.X.(*ast.Ident)
		if ok && pkg.Name == "agent" && sel.Sel.Name == "ToolRequest" {
			n++
		}
		return true
	})
	return n
}

// scanToolRequestNail309 walks the non-test Go sources under the given roots (skipping
// every file under internal/agent, which owns the one trusted newTaskID() construction),
// and returns the total number of `agent.ToolRequest{` production construction points
// together with the number of files it actually parsed — the denominator.
func scanToolRequestNail309(t *testing.T, roots []string) (constructions, examined int) {
	t.Helper()
	for _, root := range roots {
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
			rel := filepath.ToSlash(path)
			if strings.Contains(rel, "internal/agent/") { // exclude the trusted point's own package
				return nil
			}
			fset := token.NewFileSet()
			file, perr := parser.ParseFile(fset, path, nil, 0)
			if perr != nil {
				t.Fatalf("parse %s: %v", path, perr)
			}
			examined++
			constructions += countToolRequestLits(fset, file)
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}
	return constructions, examined
}

func Test309ToolRequestProductionConstructionNail(t *testing.T) {
	// go test runs with the package dir (internal/tools) as cwd; the module root is two
	// levels up. Scanned roster is named in every failure so the next reader can re-run.
	roots := []string{
		filepath.Join("..", "..", "cmd"),
		filepath.Join("..", "..", "internal"),
	}

	constructions, examined := scanToolRequestNail309(t, roots)

	// The nail must carry its own denominator: a scan that read nothing reports a false green.
	if examined < 10 {
		t.Fatalf("枚数钉分母不可信：只解析了 %d 枚非测试生产源（roots=%q，"+
			"排除 _test.go 与 internal/agent/**）。分母接近零＝这枚扫根本没读到树，报绿是假绿",
			examined, roots)
	}

	if constructions != 0 {
		t.Errorf("枚数钉响：internal/agent 之外发现了 %d 枚 `agent.ToolRequest{` 生产构造点 "+
			"(分母=%d 枚解析源，roots=%q)。票 309 的保证＝今天非空只靠 loop.go 那枚 newTaskID() "+
			"构造点（接线事实）；若有人新写一条生产构造点，它可能递进空 TaskID 让 assessorFor 静默摘掉 "+
			"R4 探测器（见 assessorfor_taint_instrument_309_test.go 的 ②）。要么给它接上非空 id，"+
			"要么把它放进被 309 AC#1 仪器覆盖的入口",
			constructions, examined, roots)
	}

	// Positive control: prove the matcher is not blind. Parse an in-memory source holding
	// one qualified construction point and assert the counter sees exactly one. (This does
	// not touch disk and does not add a production path; it only shows the ruler bites.)
	fset := token.NewFileSet()
	ctrl, perr := parser.ParseFile(fset, "ctrl.go", `package p
import "github.com/CarlosShao/wisp/internal/agent"
func f() { _ = agent.ToolRequest{TaskID: ""} }
`, 0)
	if perr != nil {
		t.Fatalf("positive control parse: %v", perr)
	}
	if got := countToolRequestLits(fset, ctrl); got != 1 {
		t.Fatalf("positive control: matcher counted %d, want 1 — the ruler is blind, "+
			"its green means nothing", got)
	}

	t.Logf("AC#2 枚数钉：internal/agent 之外 agent.ToolRequest 生产构造点=%d（==0），分母=%d 枚非测试源",
		constructions, examined)
}
