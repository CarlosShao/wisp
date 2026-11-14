// ticket 448 -- 446 的 errno 判定必须有枚「留下来的」牙。
//
// 446（be0f599f3）把判定从比较英文句子改成 errors.Is(e, syscall.Errno(0))，函数
// 叫 errnoPresent，两处调用点在 system_backend_windows.go。它「能变红」的证明是
// 一枚 -overlay 挂进包的临时探针，跑完就没了。审计 fc580501b 现跑量到的后果：
// internal/tools 里零测试提到 errnoPresent ⇒ 把判定改回句子比较，
// `go test ./internal/tools` 照样全绿。
//
// 本文件是那枚留下来的牙。它读的是**源码形状**而不是运行时返回值，所以任何 GOOS
// 都跑（不需要 Windows，也就碰不到语言/区域设置），钉三件事：
//
//	闸1  errnoPresent 体内不许出现字符串字面量、不许调 .Error()，且必须有
//	     errors.Is(<...>, syscall.Errno(...)) 形状的码判据 —— 这是「判定退回句子
//	     比较」（M1）那一条。
//	闸2  凡是把 Proc.Call 第三返回值接进命名变量、并**拿它作决定**（比较、if/for
//	     条件、switch tag）的调用点，那个决定必须由 errnoPresent(该变量) 来做 ——
//	     这是「调用点绕过判定」（M2）那一条，包括改回内联句子比较，也包括内联重写
//	     成 errors.Is（那时 errnoPresent 本身还是对的，只是没人用，行为测试看不出
//	     来）。只把 errno 印进消息、从不判它存不存在的调用点放行，理由见
//	     errnoCallSiteAudit448。
//	闸3  整个文件里错误文本（.Error()）不许被拿来作决定：不作为比较的操作数、不
//	     交给 strings.* 处理、不当 switch 的 tag。
//
// 三条闸都两头钉：Test448TheStaticChecksGoRedOnTheOldShapes 把同一批检查函数跑在
// 446 换掉的那两段源码字面量上，必须报红；再跑一段修好后的样本，必须静默。所以
// 「一处都没匹配到」不会读成绿。
//
// 为什么落在 internal/tools/*_test.go 而不是 tools/d22scan 的静态闸（票 448 要求
// 二选一并写理由）：
//   - 这条不变量是包内的。errnoPresent 是 package tools 的未导出辅助函数，出这个
//     文件没有调用者；d22scan 管的是全仓横切的禁止形状，不是某个包的某个辅助函数。
//   - d22scan 加一条 ban 是动契约。那个工具的 ban 清单是 PLAN.md D22 冻结清单的
//     转写（文件头逐条列 1..10，ban #9 的注释写明「加一棵树」都属业主批准的 scope
//     变更）⇒ 为一个符号加 ban #11 是把仪器放到契约前面，会把一次人工批准拖进一张
//     小票，还会在工具里永久留下一条「对象可以消失」的规则。
//   - 覆盖面。internal/tools 在 scripts/portable-tests.sh 的 core scope 里，本文件
//     因此在**两个** CI job 里跑（test-core 的 ubuntu 与 test-windows）；d22scan 只
//     在 lint job。守护这个文件的闸，应当在编译这个文件的那两个 job 里读数。
//   - 牙的成本。d22scan 新 ban 要配 selftestsamples.go 的种子样本＋scan_test.go 的
//     钉子；这里变红的证据直接来自票面点名的两次变异（M1/M2），加上下面的种子字面量。
//
// 反过来，若命题是「全仓任何文件都不许比较 .Error() 字符串」，那是横切形状，该进
// d22scan；本票没有主张那个命题。
package tools

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// errnoGateFile448 是这三条闸钉的唯一产码文件。
const errnoGateFile448 = "system_backend_windows.go"

// errnoGateMinSites448 是 446 落地时消费 errno 的调用点数量（EnumWindows 与
// SendInput）。把它变成下限，是为了让「扫到 0 处」也红：一条点数比事实少的规则
// 不许报绿（仓里 d22scan 对空 scope 用同一条纪律）。
const errnoGateMinSites448 = 2

// errnoGateSite448 是一处「Proc.Call 的 errno 接进命名变量」的现场。
type errnoGateSite448 struct {
	fn   string
	name string
	node ast.Node
}

func errnoGateSource448(t *testing.T) []byte {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) 失败：闸门定位不到自己，就读不到被测源码")
	}
	path := filepath.Join(filepath.Dir(thisFile), errnoGateFile448)
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("闸门读不到 %s: %v -- 判据就在这个文件里，读不到不能算绿", path, err)
	}
	return src
}

func errnoGateParse448(t *testing.T, name string, src []byte) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("闸门解析 %s 失败: %v -- 源码读不懂只能算红，不能算绿", name, err)
	}
	return fset, f
}

func errnoGateAt448(fset *token.FileSet, n ast.Node) string {
	if fset == nil || n == nil {
		return "<unknown>"
	}
	p := fset.Position(n.Pos())
	return fmt.Sprintf("%s:%d", p.Filename, p.Line)
}

// errnoPredicateViolations448 列出 errnoPresent 不再是「拿码判」的每一条理由。
// 生产测试与种子反证共用它，所以「检查器不再匹配任何东西」会被种子反证抓住，
// 而不是以静默收场。
func errnoPredicateViolations448(fset *token.FileSet, file *ast.File) []string {
	var out []string
	defs := make([]*ast.FuncDecl, 0, 1)
	for _, d := range file.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || fd.Name == nil || fd.Name.Name != "errnoPresent" {
			continue
		}
		defs = append(defs, fd)
	}
	switch len(defs) {
	case 0:
		return []string{"errnoPresent 不见了：两处调用点再没有共同判定可走（被删、被改名、被内联都算）"}
	case 1:
	default:
		out = append(out, fmt.Sprintf("errnoPresent 有 %d 个定义（第一个在 %s）：判定只能有一个，多个就说不清哪一处在两处生效",
			len(defs), errnoGateAt448(fset, defs[0])))
	}
	def := defs[0]
	if def.Body == nil {
		return append(out, errnoGateAt448(fset, def)+": errnoPresent 没有函数体")
	}
	ast.Inspect(def.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				out = append(out, fmt.Sprintf("%s: errnoPresent 体内出现字符串字面量 %s -- 判定又回到比较句子了，判据只能看码",
					errnoGateAt448(fset, x), x.Value))
			}
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Error" && len(x.Args) == 0 {
				out = append(out, fmt.Sprintf("%s: errnoPresent 体内调了 .Error() -- 错误文本不是判据（换语言、换构建就变）",
					errnoGateAt448(fset, x)))
			}
		}
		return true
	})
	if !errnoPredicateSeesCode448(def) {
		out = append(out, fmt.Sprintf("%s: errnoPresent 体内没有 errors.Is(<...>, syscall.Errno(...)) 形状的码判据 -- 这正是 446 修好的那半",
			errnoGateAt448(fset, def)))
	}
	return out
}

// errnoPredicateSeesCode448 判 errnoPresent 是否还按码判：存在一次
// errors.Is(<...>, <...syscall.Errno...>) 调用。
func errnoPredicateSeesCode448(def *ast.FuncDecl) bool {
	sees := false
	ast.Inspect(def.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Is" {
			return true
		}
		if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "errors" {
			return true
		}
		for _, a := range call.Args {
			ast.Inspect(a, func(m ast.Node) bool {
				x, ok := m.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if id, ok := x.X.(*ast.Ident); ok && id.Name == "syscall" && x.Sel.Name == "Errno" {
					sees = true
				}
				return true
			})
		}
		return true
	})
	return sees
}

// errnoCallSites448 收集消费 Proc.Call errno 的现场：*.Call(...) 第三返回值接进
// 命名变量（接进 _ 的丢弃写法不算消费，也不该被这条闸要求什么）。
func errnoCallSites448(fset *token.FileSet, file *ast.File) []errnoGateSite448 {
	var sites []errnoGateSite448
	for _, d := range file.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			as, ok := n.(*ast.AssignStmt)
			if !ok || len(as.Lhs) != 3 || len(as.Rhs) != 1 {
				return true
			}
			call, ok := as.Rhs[0].(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Call" {
				return true
			}
			id, ok := as.Lhs[2].(*ast.Ident)
			if !ok || id.Name == "_" {
				return true
			}
			sites = append(sites, errnoGateSite448{fn: fd.Name.Name, name: id.Name, node: id})
			return true
		})
	}
	return sites
}

func errnoFuncByName448(file *ast.File, name string) *ast.FuncDecl {
	for _, d := range file.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if ok && fd.Recv == nil && fd.Name != nil && fd.Name.Name == name {
			return fd
		}
	}
	return nil
}

// errnoDecisionExprs448 收一个函数体里所有「作决定」的表达式：每一处比较、每个
// if/for 的条件、每个 switch 的 tag。errno 只有在这些位置上才叫「被判过存在性」。
// 收全部比较（而不只收 if 的条件）是为了堵住「先算成布尔再判」那类搬移：
// has := e != nil && ... 里的 e != nil 本身就是一处比较。
func errnoDecisionExprs448(fd *ast.FuncDecl) []ast.Expr {
	var out []ast.Expr
	if fd == nil || fd.Body == nil {
		return out
	}
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BinaryExpr:
			if errnoIsComparisonOp448(x.Op) {
				out = append(out, x)
			}
		case *ast.IfStmt:
			if x.Cond != nil {
				out = append(out, x.Cond)
			}
		case *ast.ForStmt:
			if x.Cond != nil {
				out = append(out, x.Cond)
			}
		case *ast.SwitchStmt:
			if x.Tag != nil {
				out = append(out, x.Tag)
			}
		}
		return true
	})
	return out
}

func errnoExprMentions448(expr ast.Expr, name string) bool {
	hit := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name == name {
			hit = true
		}
		return true
	})
	return hit
}

// errnoExprCallsPresent448 判这段表达式里有没有 errnoPresent(<name>)。
func errnoExprCallsPresent448(expr ast.Expr, name string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		id, ok := call.Fun.(*ast.Ident)
		if !ok || id.Name != "errnoPresent" {
			return true
		}
		for _, a := range call.Args {
			if arg, ok := a.(*ast.Ident); ok && arg.Name == name {
				found = true
			}
		}
		return true
	})
	return found
}

// errnoCallSiteAudit448 是闸2 的本体。规则按「errno 有没有被拿来作决定」判，不是
// 按「函数里出现没出现 errnoPresent」判：
//
//   - 一处调用点把 Proc.Call 的 errno 接进命名变量，并在任何决定位置用到它（比较、
//     if/for 条件、switch tag），那这个决定必须由 errnoPresent(那个变量) 来做；
//   - 只把 errno 印进消息、从不判它存不存在的调用点放行。仓里现有一处正这样：
//     system_backend_windows.go 的 LookupPrivilegeValueW，它拿 r == 0 判失败、把
//     errno 交给 %v 打印 —— 那里没有「errno 到底有没有」这个问题可绕，硬要求它走
//     errnoPresent 是把闸门开成误报。
//
// 回传「被判过的调用点数」，让调用方能对分母做下限断言（扫到 0 处不许报绿）。
func errnoCallSiteAudit448(fset *token.FileSet, file *ast.File) (decided int, violations []string) {
	sites := errnoCallSites448(fset, file)
	cache := make(map[string][]ast.Expr, 4)
	for _, s := range sites {
		decisions, ok := cache[s.fn]
		if !ok {
			decisions = errnoDecisionExprs448(errnoFuncByName448(file, s.fn))
			cache[s.fn] = decisions
		}
		used, bypassed := false, false
		for _, d := range decisions {
			if !errnoExprMentions448(d, s.name) {
				continue
			}
			used = true
			if !errnoExprCallsPresent448(d, s.name) {
				bypassed = true
			}
		}
		if used {
			decided++
		}
		if bypassed {
			violations = append(violations, fmt.Sprintf("%s: 函数 %s 把 Proc.Call 的 errno 接进 %q 之后拿它作了决定，却没有交给 errnoPresent -- 调用点绕过了 446 的判定（判定函数本身可能还是对的，只是没人用它）",
				errnoGateAt448(fset, s.node), s.fn, s.name))
		}
	}
	return decided, violations
}

// errnoTextDecisionViolations448 报出「拿错误文本作决定」的每一处。
func errnoTextDecisionViolations448(fset *token.FileSet, file *ast.File) []string {
	var out []string
	var stack []ast.Node
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		if call, ok := n.(*ast.CallExpr); ok && errnoIsErrorTextCall448(call) {
			if why := errnoDecidesOnText448(stack, call); why != "" {
				out = append(out, fmt.Sprintf("%s: 错误文本被拿来%s -- 判据不是句子（446）",
					errnoGateAt448(fset, call), why))
			}
		}
		stack = append(stack, n)
		return true
	})
	return out
}

func errnoIsErrorTextCall448(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Error" && len(call.Args) == 0
}

// errnoDecidesOnText448 沿父链判这个 .Error() 的结果是否落进了「作决定」的位置。
// 遇到语句边界就停（那里的文本只是被打印，不是被判）。
func errnoDecidesOnText448(stack []ast.Node, expr ast.Expr) string {
	var child ast.Node = expr
	for i := len(stack) - 1; i >= 0; i-- {
		switch p := stack[i].(type) {
		case *ast.BinaryExpr:
			if errnoIsComparisonOp448(p.Op) {
				return "做比较"
			}
		case *ast.CallExpr:
			if sel, ok := p.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "strings" {
					return "交给 strings." + sel.Sel.Name + " 处理"
				}
			}
		case *ast.SwitchStmt:
			if p.Tag == child {
				return "当 switch 的 tag 用"
			}
		case ast.Stmt:
			return ""
		}
		child = stack[i]
	}
	return ""
}

func errnoIsComparisonOp448(op token.Token) bool {
	switch op {
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
		return true
	}
	return false
}

// Test448ErrnoPresentDecidesByCodeNotText 是 M1（判定退回句子比较）那条闸。
func Test448ErrnoPresentDecidesByCodeNotText(t *testing.T) {
	fset, file := errnoGateParse448(t, errnoGateFile448, errnoGateSource448(t))
	if v := errnoPredicateViolations448(fset, file); len(v) != 0 {
		t.Fatalf("446 的判定退回「比较句子」了：\n  %s", strings.Join(v, "\n  "))
	}
}

// Test448ErrnoCallSitesGoThroughErrnoPresent 是 M2（调用点绕过判定）那条闸。
func Test448ErrnoCallSitesGoThroughErrnoPresent(t *testing.T) {
	fset, file := errnoGateParse448(t, errnoGateFile448, errnoGateSource448(t))
	decided, bad := errnoCallSiteAudit448(fset, file)
	if decided < errnoGateMinSites448 {
		t.Fatalf("只有 %d 处调用点拿 errno 作了决定，至少应有 %d 处（EnumWindows 与 SendInput）：分母比事实少，这条闸就是空的",
			decided, errnoGateMinSites448)
	}
	if len(bad) != 0 {
		t.Fatalf("调用点绕过 errnoPresent 了：\n  %s", strings.Join(bad, "\n  "))
	}
	t.Logf("448 闸2：%d 处调用点拿 errno 作了决定，全部走 errnoPresent", decided)
}

// Test448NoErrorTextIsDecidedOnInTheBackend 是兜底：句子比较换了写法（strings.*、
// switch、拼过再比）也跑不掉。
func Test448NoErrorTextIsDecidedOnInTheBackend(t *testing.T) {
	fset, file := errnoGateParse448(t, errnoGateFile448, errnoGateSource448(t))
	if v := errnoTextDecisionViolations448(fset, file); len(v) != 0 {
		t.Fatalf("这个文件又在拿错误文本作决定了：\n  %s", strings.Join(v, "\n  "))
	}
}

// 下面三段是 446 换掉的形状与它修好后的形状，都是**字面量**而不是盘上的文件：
// 种子反证因此任何机器都跑，也不会因为生产文件恰好干净而被满足。
const errnoSeededOldPredicate448 = `package tools

func errnoPresent(e error) bool {
	return e != nil && e.Error() != "The operation completed successfully."
}
`

const errnoSeededOldCallSite448 = `package tools

func sendInputs() error {
	n, _, e := procSendInput.Call(0)
	if int(n) != 0 {
		if e != nil && e.Error() != "The operation completed successfully." {
			return e
		}
	}
	return nil
}
`

const errnoSeededGood448 = `package tools

func errnoPresent(e error) bool {
	return e != nil && !errors.Is(e, syscall.Errno(0))
}

func sendInputs() error {
	n, _, e := procSendInput.Call(0)
	if int(n) != 0 && errnoPresent(e) {
		return e
	}
	return nil
}
`

// Test448TheStaticChecksGoRedOnTheOldShapes 是三条静态闸自己的牙：同一批检查
// 函数在 446 修复前的两段源码上必须报红，在修好后的样本上必须静默。
func Test448TheStaticChecksGoRedOnTheOldShapes(t *testing.T) {
	t.Run("446 修复前的判定体必须报红", func(t *testing.T) {
		fset, file := errnoGateParse448(t, "seeded-old-predicate.go", []byte(errnoSeededOldPredicate448))
		v := errnoPredicateViolations448(fset, file)
		if len(v) == 0 {
			t.Fatal("闸1 在 446 修复前的判定体上一处都没报 -- 这个检查已经不会红了")
		}
		t.Logf("闸1 命中 %d 处:\n  %s", len(v), strings.Join(v, "\n  "))
		if tv := errnoTextDecisionViolations448(fset, file); len(tv) == 0 {
			t.Fatal("闸3 在 446 修复前的判定体上一处都没报 -- 兜底检查已经不会红了")
		}
	})

	t.Run("446 修复前的调用点必须报红", func(t *testing.T) {
		fset, file := errnoGateParse448(t, "seeded-old-callsite.go", []byte(errnoSeededOldCallSite448))
		decided, v := errnoCallSiteAudit448(fset, file)
		if decided != 1 {
			t.Fatalf("种子样本里应有 1 处调用点拿 errno 作了决定，实得 %d -- 扫描器已经不认识这个形状了", decided)
		}
		if len(v) == 0 {
			t.Fatal("闸2 在 446 修复前的调用点上一处都没报 -- 这个检查已经不会红了")
		}
		t.Logf("闸2 命中 %d 处:\n  %s", len(v), strings.Join(v, "\n  "))
	})

	t.Run("446 修复后的样本必须静默", func(t *testing.T) {
		fset, file := errnoGateParse448(t, "seeded-good.go", []byte(errnoSeededGood448))
		if v := errnoPredicateViolations448(fset, file); len(v) != 0 {
			t.Fatalf("闸1 在干净样本上误报：\n  %s", strings.Join(v, "\n  "))
		}
		if decided, v := errnoCallSiteAudit448(fset, file); len(v) != 0 || decided != 1 {
			t.Fatalf("闸2 在干净样本上误报（decided=%d）：\n  %s", decided, strings.Join(v, "\n  "))
		}
		if v := errnoTextDecisionViolations448(fset, file); len(v) != 0 {
			t.Fatalf("闸3 在干净样本上误报：\n  %s", strings.Join(v, "\n  "))
		}
	})
}
