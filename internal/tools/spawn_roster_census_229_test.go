package tools

// Ticket 229, AC#4: the D38b roster is only worth freezing if every spawn name
// in the repo actually lands in it. CategoryUnknown is the leak alarm's唯一
// signal ("this goroutine is not accounted for"); each accidental off-roster
// name turns that signal into routine noise. This scan is the "尺 + 断" pair:
//
//   - walk every production .go under cmd/ internal/ tools/ frontend/ (the
//     module's first-party roots; test files are out of scope, they plant
//     fake names on purpose);
//   - a Spawn first argument that IS a string literal is classified against
//     the live roster: off-roster and not on the documented-debt list => red;
//   - a Spawn first argument that is NOT a literal can't be judged statically,
//     so the site itself must be on the reviewed list below; an unreviewed
//     site => red (someone added a spawn name and nobody wrote down where it
//     lands).
//
// A file that does not parse is counted and logged, not failed: this is a
// shared work tree where another leg may be mid-write, and this scan's target
// is the roster, not the parser. NewRoot/NewRootFrom first arguments are root
// ids (labels in panic records), never goroutine names; they are logged for
// the audit trail and deliberately not classified.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// spawn229OffRosterKnown is the documented-debt ledger: names that are
// deliberately NOT in the D38b roster, each with the in-code rationale that
// keeps the alarm honest (they are not accidents - each one is a site whose
// author recorded why the noise stays until an owner decides).
var spawn229OffRosterKnown = map[string]string{
	"panel-sta":               "cmd/wisp/panel_resident_windows.go:80-84 逐字写了：deliberately NOT one of observe.ResidentNames 的六枚常驻名（加名是 D38b 契约改）",
	"balldebug-hotkey-bridge": "cmd/balldebug/main.go:35-43 逐字写了：R16#5 禁止借名册名蒙过漏哨，噪声留着、归票 66",
	"balldebug-level-feeder":  "同上（票 66）",
	"mockllm-stdout-reader":   "internal/llm/adaptertest/mockllm.go:92；cmd/balldebug/main.go:42 登记为票 67 的 loose end",
}

// spawn229ReviewedSites is the manual-review ledger for Spawn sites whose
// first argument is not a string literal. Key = "<repo-relative path>|<expr
// text as written>"; value = the runtime name and the roster category it
// lands in. A site missing from here fails the scan on purpose.
var spawn229ReviewedSites = map[string]string{
	`internal/agent/loop.go|"agent-task-"+id`:                                            "per_task：任务根协程（loop.RunAsync）",
	`internal/agent/loop.go|"tool-exec-"+taskID`:                                         "per_task：每任务工具执行",
	`internal/plugin/disposal.go|name`:                                                   "参数；生产唯一调用点 internal/memory/retention.go:103 传 retention-job（temporary）",
	`cmd/balldebug/main.go|goroutineHotkeyBridge`:                                        "balldebug-hotkey-bridge：有意名册外，见 spawn229OffRosterKnown",
	`cmd/balldebug/main.go|goroutineLevelFeeder`:                                         "balldebug-level-feeder：同上",
	`cmd/wisp/panel_resident_windows.go|panelSTAName`:                                    "panel-sta：有意名册外，见 spawn229OffRosterKnown",
	`cmd/wisp/resident_task_source_windows.go|taskSourceConsoleGoroutine`:                "approval-waiter：per_task",
	`cmd/wisp/resident_task_source_windows.go|taskSourceGoroutinePrefix+fmt.Sprint(seq)`: "agent-task-resident-<n>：per_task 前缀族",
	`internal/tools/subagent_197.go|"agent-task-"+bg.ID`:                                 "per_task：票 229 修法后的子代理收尾协程（bg.ID 是孩子任务 id）",
}

// expr229 renders one AST expression back to source text, whitespace-folded so
// the ledger keys survive line moves and reformatting.
func expr229(fset *token.FileSet, src []byte, e ast.Expr) string {
	start := fset.Position(e.Pos()).Offset
	end := fset.Position(e.End()).Offset
	if start < 0 || end > len(src) || start >= end {
		return "?"
	}
	return strings.Join(strings.Fields(strings.ReplaceAll(string(src[start:end]), "\n", " ")), " ")
}

func Test229SpawnNamesAreRosterShaped(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// thisFile = <repo>/internal/tools/spawn_roster_census_229_test.go
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))

	fset := token.NewFileSet()
	var checked, skipped, spawnSites, rootSites, offRosterKnown, unreviewed int

	for _, root := range []string{"cmd", "internal", "tools", "frontend"} {
		base := filepath.Join(repoRoot, root)
		err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if name := d.Name(); name == "testdata" || strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			rel, rerr := filepath.Rel(repoRoot, path)
			if rerr != nil {
				return rerr
			}
			rel = filepath.ToSlash(rel)
			src, rerr := os.ReadFile(path)
			if rerr != nil {
				return rerr
			}
			file, perr := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
			if perr != nil {
				skipped++
				t.Logf("229 census: 解析跳过 %s（共享树里可能在写到一半）: %v", rel, perr)
				return nil
			}
			checked++
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok || len(call.Args) == 0 {
					return true
				}
				line := fset.Position(call.Pos()).Line
				switch sel.Sel.Name {
				case "NewRoot", "NewRootFrom":
					// Root ids are labels (panic records), not goroutine names:
					// ClassifyGoroutine is never called on them. Logged for audit.
					rootSites++
					t.Logf("229 census: %s:%d %s(%s) 根标签（不参与名册判定）", rel, line, sel.Sel.Name, expr229(fset, src, call.Args[0]))
					return true
				case "Spawn":
				default:
					return true
				}
				spawnSites++
				expr := expr229(fset, src, call.Args[0])
				if lit, isLit := call.Args[0].(*ast.BasicLit); isLit && lit.Kind == token.STRING {
					name, uerr := strconv.Unquote(lit.Value)
					if uerr != nil {
						t.Errorf("%s:%d: Spawn 首参字面量解不开: %v", rel, line, uerr)
						return true
					}
					if cat := observe.ClassifyGoroutine(name); cat != observe.CategoryUnknown {
						t.Logf("229 census: %s:%d Spawn(%q) -> %s", rel, line, name, cat)
						return true
					}
					if why, known := spawn229OffRosterKnown[name]; known {
						offRosterKnown++
						t.Logf("229 census: %s:%d Spawn(%q) 名册外但有登记：%s", rel, line, name, why)
						return true
					}
					unreviewed++
					t.Errorf("%s:%d: Spawn 名 %q 不在 D38b 名册：它每次 spawn 都会让"+
						"「goroutine outside the D38 roster (leak symptom)」响一次，泄漏报警就淹在噪声里。"+
						"要么把名字改成名册内的形状（每任务前缀族 agent-task-/tool-exec-/approval-waiter，"+
						"或 resident/on-demand/temporary 表里的名字），要么落一条批准记录后再动表；"+
						"确属有意的名册外站点，把名字与理由登记进 spawn229OffRosterKnown", rel, line, name)
					return true
				}
				if why, reviewed := spawn229ReviewedSites[rel+"|"+expr]; reviewed {
					t.Logf("229 census: %s:%d Spawn(%s) 非字面量、已复审：%s", rel, line, expr, why)
					return true
				}
				unreviewed++
				t.Errorf("%s:%d: Spawn 首参不是字符串字面量（%s）且没有复审记录："+
					"把运行期名字与它落的名册类别写进 spawn229ReviewedSites", rel, line, expr)
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", base, err)
		}
	}

	t.Logf("229 census 汇总：生产 .go %d 枚（解析跳过 %d 枚）、Spawn 站点 %d 个、"+
		"根标签 %d 个；名册外有登记 %d 枚、未复审 %d 枚",
		checked, skipped, spawnSites, rootSites, offRosterKnown, unreviewed)
	if checked < 50 {
		t.Fatalf("只扫到 %d 枚生产 .go，仓布局变了？", checked)
	}
	if spawnSites < 10 {
		t.Fatalf("只扫到 %d 个 Spawn 站点，扫描器坏了？", spawnSites)
	}
}
