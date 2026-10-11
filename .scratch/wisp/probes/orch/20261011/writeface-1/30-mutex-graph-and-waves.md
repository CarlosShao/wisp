# writeface-1 — 30 互斥图与波次

判据口径：**两两配对**，"⛔ 并发"必须落在**一枚具体文件名**或**一枚具体 Go 面试图**上；⛔ 出现"都在 `cmd/wisp` 里所以互斥"这种没名字的判语。
写面名册取自 `10-per-leg-file-roster.md`（本腿自跑尺，HEAD `934b350a`；现 HEAD `47e0ff78`，⛔ 源码字节差＝尺 `git diff --stat 934b350a HEAD -- internal cmd` 输出**空**）。

## 一、逐对判定

| 对 | 写面交集（文件级） | 判 | 具名理由（⛔ 空口） |
|---|---|---|---|
| `296-r2` × `295-r1` | `{resident_windows.go}` × `{resident_audio_windows.go, resident_mute_290_windows_test.go, resident_audio_247_windows_test.go, resident_audio_247_live_windows_test.go}`＝**∅** | ⛔ **⛔ 文件级互斥不成立**（★顶回点） | 唯一真耦合＝**同取 `./cmd/wisp/` 那一枚包**：`296-r2` 的门禁逐字含 `GOFLAGS= go build ./...`，而整模块编译要把 `295-r1` 在共享工作树上未提交的 `resident_audio_windows.go`（`	started bool`→`	started atomic.Bool` 那一发）**一起**编；半成品在＝`296-r2` 假红。⇒ 这是**颜色归因**冲突，⛔ 文件名冲突 |
| `296-r2` × `294-r1` | `{resident_windows.go}` × `{resident_mute_294_test.go(新), …windows 测试件}`＝**∅** | ⚠ **软边**：写-读同一枚文件 | `294-r1 AC#1` 的钉**parse** 那枚文件（模板逐字 `parser.ParseFile(fset, filepath.Join(".", "resident_windows.go"), …)`）。`296-r2` 改的是 `hotkeyReloadSource296` 的文档块，**在 `runResident` 体⛔** ⇒ 走查计数⛔ 变；唯一被咬的形＝**行号**：`296-r2` 若改注释行数，`resident_windows.go` 的 `:342`／`:249` 等号全漂。票 294 裁定节点⛔ 已自定"派单一律内容锚" ⇒ 只要⛔ 出带行号的锚，两枚可同波 |
| `295-r1` × `294-r1` | **∅（有条件）** | ⛔ **⛔ 编排者那条硬串行⛔ 是文件必然**（★顶回点，详见 §四） | 条件＝派单钉死三条：① `294-r1` 写面⛔ 出它自己的新件；② `AC#3` 反形（`resident_audio_windows.go` 的 `	case gate.Muted():`／`	case gate.Open():` 两支 `true`→`false`）走 **`-overlay`**，⛔ 改树——该形在本仓已成型（`298-v1` 用 `-overlay` 载改前内容；`296-v1` 用仓外备份＋`sha256sum` 还原证明）；③ `AC#2` 的镜像常量与真窗用例⛔ 另起一枚 `_windows_test.go`，⛔ 加进 `resident_mute_290_windows_test.go`（那枚件 `295-r1` 必动 4 枚读点）。三条任一⛔ 钉 ⇒ 互斥**当场成立** ⇒ 那就维持串行 |
| `309-r1` × 三枚 `cmd/wisp` 腿 | `{internal/tools/bridge.go, internal/tools/ 1 枚新测试件}` × `cmd/wisp/**`＝**∅** | ⛔ **互斥不成立（可同波）** | 包也⛔ 相交（`go list ./...` 现量 35 枚包，二者并列存在）；`309-r1 AC#4` 要的是 `./internal/tools/` 色 ⇄ `go vet ./internal/tools/`。⚠ 一枚**唯一**真耦合还是 `go build ./...`（整模块）——同上一格那句。另：`309-r1 AC#2` 那枚枚数钉的**分母**会读盘（`os.ReadDir`/`WalkDir`），若 `294-r1` 正把新件躺在工作树里，枚数随台面变 ⇒ 派单要钉"分母射程＝`internal cmd` 的**非 `_test.go`**"（现量：`ToolRequest{` 全仓 tracked **252**／`internal cmd` **60**／产码 **1**＝`internal/agent/loop.go:675`） |
| `298 残余` × 全部 | **∅（落地半边⛔ 存在）** | 无边 | 尺①`git merge-base --is-ancestor 18d43f8e HEAD`＝YES；尺②两枚目标件 `HEAD:<path>` ⇄ `18d43f8e:<path>` blob 名逐字相同 ⇒ 见 `10-…` E 节 |

## 二、可并行分组（波次建议）

★先交一枚**台面更新**（它改了上面所有"在飞"前提）：起手时编排者说 `298-v1` 与 `308-a1` 在飞；本腿收尾现量 `git log --oneline 934b350a..HEAD`（⛔ 本腿的笔）＝**2 笔**：`7d163392`（`298-v1` 验收交件）＋`dd638f98`（`308-a1` 台面读数）⇒ **两枚都已交回**，且 `git status --porcelain -- internal cmd` 复跑仍 **0 行**。
⇒ 票 309「排程更新」那句"⚓ 按住 `309-r1`，因 `298-v1` 正在取验收色、`308-a1` 正在读 `internal/risk` 台面色"的**两条前提今天都已消失**（该句在票面上，⛔ 由本腿改；记为过期读数）。

- **波 A（现在即可同波派 2 枚）**：`296-r2` ＋ `309-r1`。
  理由＝上表两枚边都⛔ 成立（文件∅、包∅、⛔ 改语义）。
  ⚠ 波 A 里 `296-r2` 那一枚带 M-4（`-run 'Test296|Test258'` 两族起真球窗）⇒ 无桌面会话时把它挪到波 A 尾或编排者亲跑那一发，⛔ 让腿交 SKIP。
- **波 B**：`295-r1` 独占 `cmd/wisp` 写面＋`./internal/audio/` 色。
- **波 C**：`294-r1` → `294-v1`。
  **B/C 可并成一波**，代价＝派单必须逐字钉上面那三条（①②③）；编排者若⛔ 愿意钉那三条 ⇒ 维持 `B→C` 串行。本腿⛔ 断言"串行一定浪费机时"，本腿只断言一件事：它⛔ 是**文件必然**，而是"票 294 排程节那句旧话＋同包名册归因"的产物。
- 与编排者原队列的差：原队列 `296-r2 → 295-r1 → 294-r1 → 294-v1` 是**四拍全串**；按本件 ⇒ **`296-r2` ∥ `309-r1`（1 拍）→ `295-r1`（1 拍）→ `294-r1`→`294-v1`**，省 1 拍；若接受 ①②③ ⇒ 再省 1 拍（`295-r1` ∥ `294-r1`）。

## 三、`294-r1 AC#1` 那把尺的形状（雷#2 的答案）

编排者那句"⚓ 照抄 258 会造恒绿假钉（`call.Fun` 是 `*ast.SelectorExpr` 不是 `*ast.Ident`）"——**技术核〔已证〕；但那枚钉的"颜色"写反了**：

- 核（现量证实）：258 那把尺逐字 `		if id, ok := call.Fun.(*ast.Ident); !ok || id.Name != "startResidentBall" {`——它成立是因为 `startResidentBall` 是**包级函数**调用；生产那枚是方法调用 `	rb.attachMuteGate(raudio.toggleMute)` ⇒ `call.Fun` 是 `*ast.SelectorExpr`，纯照抄＝计数永远 0。
- 但"永远绿"⛔ 是这一形的颜色：票面 `AC#1` ⓐ 逐字要求次数 **== 1** ⇒ 计数恒 0 的钉会**永远红**（响亮），⛔ 静默放行。
- **真正恒绿的形**是 `AC#1` 的 ⓑ/ⓒ 那两半被写成"对找到的每一枚调用检查 X"（0 枚找到 ⇒ 循环体⛔ 执行 ⇒ ⛔ 判空集）＝票 309 范本注释逐字点名的病：*"a scan whose denominator is zero reports clean for a tree it never read"*（`internal/tools/task_state_188_test.go`）。⇒ **派单该钉的是"ⓐⓑⓒ 三件事都⛔ 共用同一枚计数，且计数==1 单独成断言"**。
- 省力读数：⛔ 非要新造 SelectorExpr 层——**同包现成范本**逐字 `		if sel, ok := d.Call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "stop" {`（`cmd/wisp/resident_ball_228_test.go`，即票 294 现量节点名的三枚 AST 先例之一）；尺＝`git grep -n "ast.SelectorExpr" HEAD -- cmd/wisp` ⇒ **41 命中／13 枚件**（含测试，射程 `cmd/wisp`）。
