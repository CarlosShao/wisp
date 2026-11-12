# 427 — 12r2 的 AST 检查锚点过期：从 `openScope` 移到 `scopeOptions`，并收紧到前缀分支

- 票面：`.scratch/wisp/issues/427-the-12r2-ast-check-still-looks-for-marksettling-inside-openscope-but-it-moved-into-scopeoptions-hooks-logf.md`
- 写面：**仅** `cmd/wisp/settle_trace_12r2_windows_test.go`（产码零字节）
- 实现者：写腿 427-12r2（裁决者≠实现者，翻框归编排者）

---

## §1 AC#0 — 复现（逐字红句）＋ 扫描范围结论

### 1.1 复现（改前，carrier 真跑）

```
$ D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test12r2SeamCallSitesAreTheMeasurement$'"
=== RUN   Test12r2SeamCallSitesAreTheMeasurement
    settle_trace_12r2_windows_test.go:193: 12r2: openScope no longer stamps T1 on the scope's SCOPE-SETTLING announcement, so the settle entry leg of AC#9 item 1 is unmeasurable
--- FAIL: Test12r2SeamCallSitesAreTheMeasurement (0.01s)
FAIL
[exit code: 1]
```

耗时 **0.01s**（纯静态检查，不跑产品）⇒ 稳定，与负载/时长无关。

### 1.2 ★扫描范围结论：T1 检查把 `ast.Inspect` 锚在哪个 `FuncDecl` 上

**逐字读改前 `:167`–`:191`**，答案＝**`openScope`**：

```go
167: 	ast.Inspect(enginesAst, func(n ast.Node) bool {
168: 		fn, ok := n.(*ast.FuncDecl)
169: 		if !ok || fn.Name.Name != "openScope" || fn.Body == nil {   // <== 锚点在这里
170: 			return true
171: 		}
172: 		ast.Inspect(fn.Body, func(m ast.Node) bool {
```

- 外层 `ast.Inspect(enginesAst, ...)` 走整棵 `resident_voice_engines_windows.go` 的 AST，
  但**一旦命中名字等于 `"openScope"` 的 `FuncDecl` 就 `return false`**（`:190`），停止外层遍历；
  真正的判据 `:172`–`:189` 全部跑在 **`fn.Body`＝`openScope` 的函数体**里。
- 两枚判据：`markSettling` 的 `SelectorExpr`（`:176`）与 `"session: SCOPE-SETTLING"` 的 `BasicLit`（`:184`）。
- ⇒ 判据面＝**`openScope` 一个函数体**，不是整个文件。**这决定了修法**：不是放宽成全文 grep，
  而是把锚点搬到产码真实所在的那个 `FuncDecl`。

### 1.3 产码侧真实位置（现读，非猜）

```
cmd/wisp/resident_voice_engines_windows.go
477: func (r *voiceResidency) scopeOptions() session.ScopeOptions {
479:     Logf: func(format string, args ...any) {
483:         if strings.HasPrefix(format, "session: SCOPE-SETTLING") {
484:             r.settle.markSettling()        <== markSettling 的唯一调用点
485:         }
508:             r.settle.markFellBack()        <== T2
545: func (r *voiceResidency) openScope(reason string) (*session.Scope, error) {
546:     return session.NewScope(r.scopeOptions())   <== 里面既无 markSettling 也无 SCOPE-SETTLING
```

`git grep -n 'markSettling()'` 只回两处：定义 `settle_trace_12r2.go:78` 与调用 `resident_voice_engines_windows.go:484`。
⇒ **产品没坏，是检查的位置过期**（同族 `Test12r2SettleWalkFiveSamplesEndToEnd` PASS 22.66s 是旁证）。

---

## §2 AC#1 — 选形：**乙**（锚到前缀分支），并说明为什么不是甲

### 2.1 落地的判据（改后）

1. 锚点 `FuncDecl`：`openScope` ⇒ **`scopeOptions`**（T1 与 T2 一并搬，两枚戳现在都在这一个函数里）。
2. T1 判据从"函数体里出现过"收紧为 **"`markSettling` 在 `HasPrefix(..., "session: SCOPE-SETTLING")` 这个 `IfStmt` 的 `Body` 内"**：
   - `testsSettlingPrefix(cond)`：条件里必须有 `HasPrefix` 调用，且 `"session: SCOPE-SETTLING"` 是**该调用的实参**；
     接受复合条件（前缀判断 AND 上别的守卫仍算测前缀），但不认"别的分支里恰好也有这个字面量"。
   - 只扫 `ifStmt.Body`：**条件本身与 `else` 臂都不算**——那是两处"看起来锚着、其实没锚"的位置。
3. 新增一条独立的 `!settlingBranch` 失败：**分支本身消失**与**戳从分支里掉出去**是两种不同的回归，
   红句分开点名（前者＝锚没了，后者＝戳没锚在分支里）。
4. `settle12r2SettlingPrefix` 常量单点持有那个字面量，判据与红句不可能各自漂移。
5. T0（`run.go` 的 `consoleSink.Publish`）与 T2（`markFellBack`）判据语义不变，T2 只换锚点函数。

### 2.2 为什么选乙（一句话）

甲只能证明"同一个函数里两样都在"，而 AC#9① 买的承诺是**打戳被前缀判断门控**——
把 `markSettling` 提到 `if` 外面，产品会**在每一条 scope 审计行上都打 T1**（错读数，不是搬家的读数），
而甲对这一形**照样绿**；乙是唯一能抓住"戳在、锚掉了"这一形的形。§3.3 的实测把这条差别量成了读数。

---

## §3 AC#2 — 突变控制：两发**真跑**（仓外拷贝，⛔ 未用 `-overlay`）

### 3.1 为什么这一发不能走 `-overlay`（本程实测的仪器事实，⛔ 藏）

本检查在**运行时从盘上读源码**（`parser.ParseFile(fset, "resident_voice_engines_windows.go", nil, 0)`，
相对 cwd＝`cmd/wisp`），而 `-overlay` 只换**编译期**输入。实测坐实（与 `120-fd-anchor.md` §3 同族结论）：

```
# 用 overlay 把「删掉 markSettling」的变异喂给编译器，产出的测试二进制跑在**未变异的盘上源码**旁：
$ go test -c -overlay=overlay-del.json -o ov/wisp.test.exe ./cmd/wisp/     # build rc=0
$ (cd cmd/wisp && ov/wisp.test.exe -test.run '^Test12r2SeamCallSitesAreTheMeasurement$')
--- PASS: Test12r2SeamCallSitesAreTheMeasurement (0.01s)
PASS
rc=0
```

⇒ **overlay 在这枚检查上是假绿仪器**：变异进了二进制、没进被读的字节，读数恒绿。
故本程按票面"拷贝放仓外"的字面做**盘上改写**，工作树全程零突变：

- 基线树 `D:\tmp\wisp\20261112\427-12r2\mut1\`：`run.go` + `resident_voice_engines_windows.go` + 测试二进制 + 三枚 DLL（同目录，loader 第一搜索位）。
- **同一枚测试二进制**（`sha256 917345372134F1BF…A012F2`，三棵树逐字节相同）配**不同的盘上文本** ⇒ 不同判决。
  这是本程读数"不是编译期差异"的硬证据。

### 3.2 发①：删掉 `:484` 的 `r.settle.markSettling()`

变异：删整行（分支保留、变空），`diff` 只此一行，文件 1311 → 1310 行。

```
=== RUN   Test12r2SeamCallSitesAreTheMeasurement
    settle_trace_12r2_windows_test.go:217: 12r2: scopeOptions' "session: SCOPE-SETTLING" branch no longer stamps T1 (markSettling is not inside that branch's body), so the settle entry leg of AC#9 item 1 is unmeasurable
--- FAIL: Test12r2SeamCallSitesAreTheMeasurement (0.02s)
FAIL
rc=1
```

⇒ **必红，实测红，rc=1。**

### 3.3 发②：把 `markSettling()` 从 `if` 分支里挪到 `if` 外面（同函数内、脱离前缀判断）

变异后该闭包逐字：

```go
		Logf: func(format string, args ...any) {
			// …注释不动…
			r.settle.markSettling()                                  // <== 提到 if 外面
			if strings.HasPrefix(format, "session: SCOPE-SETTLING") {
			}
			slog.Info("voice[scope]: " + fmt.Sprintf(format, args...))
		},
```

```
=== RUN   Test12r2SeamCallSitesAreTheMeasurement
    settle_trace_12r2_windows_test.go:217: 12r2: scopeOptions' "session: SCOPE-SETTLING" branch no longer stamps T1 (markSettling is not inside that branch's body), so the settle entry leg of AC#9 item 1 is unmeasurable
--- FAIL: Test12r2SeamCallSitesAreTheMeasurement (0.01s)
FAIL
rc=1
```

⇒ **选乙 ⇒ 必红，实测红，rc=1。**

### 3.4 ★甲乙差别：量出来的，不是嘴上说的

仓外纯 `go/ast` 探针 `D:\tmp\wisp\20261112\427-12r2\predicate\main.go`（零包依赖，只读路径实参），
把甲的判据（"`scopeOptions` 函数体里两样都在"）与乙的判据（"`markSettling` 在前缀分支体内"）
同时喂给三棵树：

```
mut1（基线）        A(jia): markSettling=true  literal=true  => GREEN | B(yi): branch=true  stamp-in-branch=true  => GREEN
mut1-del（删戳）    A(jia): markSettling=false literal=true  => RED   | B(yi): branch=true  stamp-in-branch=false => RED
mut2-hoist（提出去）A(jia): markSettling=true  literal=true  => GREEN | B(yi): branch=true  stamp-in-branch=false => RED
```

**逐格读**：发① 甲也红（戳真没了）；**发② 甲绿、乙红**——这正是票面 :75 预言的差别，现量坐实。
甲在这一格放过的是"T1 会对每一条 scope 审计行都打戳"这个错读数，所以乙更强、也是本程选的形。

### 3.5 还原核

基线树 `mut1\` 自始至终是产码原样（`git show HEAD:` 逐字节），改后重跑：

```
--- PASS: Test12r2SeamCallSitesAreTheMeasurement (0.01s)
PASS
rc=0
```

主仓 `git status --porcelain -- cmd/wisp` 全程只有本腿那一枚测试文件（另有一枚他腿 426 的脏改，见 §4）。

---

## §4 AC#3 — 越界检查（`git diff --name-only` 逐行）

```
$ git diff --name-only
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt
cmd/wisp/resident_voice_wire_r1_windows_test.go
cmd/wisp/settle_trace_12r2_windows_test.go
frontend/fixtures/composer-states.html
```

**逐行归属**（⛔ 藏他腿的活）：

| 行 | 是谁的 | 判据 |
|---|---|---|
| `cmd/wisp/settle_trace_12r2_windows_test.go` | **本腿（唯一写面）** | 本程唯一编辑 |
| `cmd/wisp/resident_voice_wire_r1_windows_test.go` | **他腿 426**（`TestWireR1WarmIdleDisposesOnTheOwner` 的 mic-close 计数） | mtime `20:09:04`，本程起手 `20:03`；diff 内容是 owner-queue barrier 改写，与本票无关；本程 ⛔ 未触碰、⛔ 未 commit |
| `.scratch/wisp/probes/291/…`、`.scratch/wisp/probes/orch/…`、`frontend/fixtures/composer-states.html` | 起手前既存脏改 | 本程 `20:03` 首次 `git status` 即在，本程零触碰 |

- 本腿**产码零字节**：`resident_voice_engines_windows.go`／`settle_trace_12r2.go` 未改（两枚戳位置是既定实现，见票面「本票不解决」）。
- 禁面全部未动：`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`。
- `build/mut_rvw.go` 未删（未跟踪残留件，非本腿写面）。
- commit 用**显式 pathspec**（只 `cmd/wisp/settle_trace_12r2_windows_test.go` ＋ 本证据件），⛔ push。

---

## §5 AC#4 — 门禁四读数（本程现跑）

| 尺 | 命令 | 读数 |
|---|---|---|
| 构建 | `go build ./cmd/wisp/` | **rc=0**（按票面用 `./cmd/wisp/`，⛔ 不用 `./...`，避开 `build/mut_rvw.go` 残留） |
| 静态检查 | `go vet ./cmd/wisp/` | **rc=0** |
| 格式（两尺） | `gofmt -l cmd/wisp/settle_trace_12r2_windows_test.go`<br>`D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp/settle_trace_12r2_windows_test.go` | **两尺均零输出**（rc=0） |
| D22 禁令 | 自建 `D:\tmp\…\d22\d22scan.exe`（`tools/d22scan` 是独立 module，`go build .`）跑主仓 | **clean，rc=0，零新增违规**：bans#1-5 internal/=403、cmd/=77；ban#6 frontend/=123；ban#7 internal/tools/=52；ban#10 internal/risk/=21；ban#8 design/=39、frontend/=123、internal/=1004、cmd/=276 |

> `tools/d22scan` 在 `go build ./tools/d22scan/` 下报
> `main module does not contain package …/tools/d22scan`（它是自带 `go.mod` 的独立 module），
> 故在其目录内 `go build .`，与既有各腿口径一致。

**载波口径**：测试一律经 `D:\work\soft\Git\bin\bash.exe` ＋ `scripts/wisp-cli-staged-tests.sh`；
本机 PATH 上的 `bash` 是 WSL 转发器（一枚测试都不跑），本程未用。

---

## §6 AC#5 — 与票 12 的关系（⛔ 不翻任何框）

- 本票**只修这一枚的检查锚点**：`Test12r2SeamCallSitesAreTheMeasurement` 从红转绿（`PASS 0.01s`）。
- ⛔ **不翻票 12 任何框**；⛔ 不宣称 AC#9 的 T0/T1/T2 三戳全绿——T2 那半（`:196-199` 原判据，现 `:220`）
  本程只换锚点函数、语义未变，**它是绿是红由裁决者按票 12 口径另判**。
- 同族 `Test12r2SettleWalkFiveSamplesEndToEnd` 本程未跑（票面明示不碰；它已 PASS 22.66s 是"机制活着"的旁证）。
- 产码 `markSettling` 的调用位置＝既定实现，⛔ 本程未动；要改它＝改产品行为，需单独立票。

---

## §7 未做／具名残余（⛔ 藏）

1. **本程未跑整包 `cmd/wisp`**：本机整包有他腿在飞（426 的脏改就在同包），且票面只要求这一枚。
   全包读数不归本腿声明。
2. **`-overlay` 仪器坑只在本检查上实测**（§3.1）：本程没有普查其他"运行时读盘"的检查是否也被 `-overlay` 假绿。
   已知同族先例＝`120-fd-anchor.md` §3 与 `102r2-r1-verdict.md` §方法注记；这是方法层的公共债，不归本票。
3. **`testsSettlingPrefix` 只认 `HasPrefix` 这个名字**：若产码换成 `strings.Contains` 或自带前缀函数，
   检查会红（＝要求人来看一眼锚是否还在），这是**故意的**方向：宁可红，不可假绿。
4. **`run.go` 那半（T0）本程零改动**，其锚点是否也有过期风险，本程未审计。

---

## §8 停手条件自核

1. **未定义即停？** 未触发：票面已给定甲/乙两形并授权实现者自拍；本程选乙并写明理由（§2.2），无自行发明判据。
2. **要不要改产码？** 不要：两枚戳位置是既定实现，票面明列「本票不解决」。
3. **要不要翻框？** 不要：翻框归编排者，本程只追加 Progress log。
