# 427 — `Test12r2SeamCallSitesAreTheMeasurement` 恒红：AST 检查仍在 `openScope` 里找 `markSettling`，而它已重构进 `scopeOptions()` 的 `Hooks.Logf` 闭包

- Status: **ready-for-agent**（写面＝`cmd/wisp/settle_trace_12r2_windows_test.go` 的**断言侧唯一一处**；⚠ 本票**不碰**任何产码、不碰 `cmd/wisp/resident_voice_engines_windows.go`、不碰 `settle_trace_12r2.go`）
- 来源：round 88 整包复跑（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`，`PASS=593 / FAIL=14`）
- 关联：票 12（AC#9 的 T0/T1/T2 三戳）· 票 382WR2（同 `resident_voice_wire` 面）· 台账 `A1801`

## 这是什么（人话）

有一枚测试要保证：**"进 Settling 的那一刻"这个时间戳，必须由 scope 自己宣布的那条审计行来打**。
它不跑产品，而是**读源码**（AST 静态检查）：去 `openScope` 这个函数里找两个东西——
① 有没有调用 `markSettling`（打 T1 戳）② 有没有那个 `"session: SCOPE-SETTLING"` 字面量。

问题：后来产码**重构**了，这两样东西**搬到了 `scopeOptions()` 的 `Hooks.Logf` 闭包里**。
测试还在老地方（`openScope`）找 ⇒ **找不到 ⇒ 必红**。
**产品没坏，是检查的位置过期了。**

## 现量（可逐字复跑）

**逐字红句**（`cmd/wisp/settle_trace_12r2_windows_test.go:193`）：

> `12r2: openScope no longer stamps T1 on the scope's SCOPE-SETTLING announcement, so the settle entry leg of AC#9 item 1 is unmeasurable`

**耗时 0.00s**（纯静态检查，不跑产品）⇒ **稳定，与负载/时长/机器无关**。

**同族另一枚是绿的**（旁证：这套机制本身活着）：`Test12r2SettleWalkFiveSamplesEndToEnd` **PASS 22.66s**。

## 机制链（逐行可核）

**测试侧**（`settle_trace_12r2_windows_test.go`）：
- `:129` `parser.ParseFile(fset, "run.go", nil, 0)`
- `:133` `parser.ParseFile(fset, "resident_voice_engines_windows.go", nil, 0)`
- `:127` `func Test12r2SeamCallSitesAreTheMeasurement`
- T1 检查（`:170`–`:195` 区间）：找 `markSettling` 的 `SelectorExpr`（`:176`）与 `"session: SCOPE-SETTLING"` 的 `BasicLit`（`:184`），二者缺一即 `:193` Fatal。
- ⚠ **关键**：它用 `ast.Inspect` 扫的是**某个函数体**——`openScope` 的函数体（`:545` 起）。**去 `:160`–`:175` 附近读清它到底把范围锚在哪个 `FuncDecl` 上**，这决定了修法。

**产码侧（现在的真实位置）**：
- `openScope`：`cmd/wisp/resident_voice_engines_windows.go:545` `func (r *voiceResidency) openScope(reason string) (*session.Scope, error)`
  - **里面既没有 `markSettling` 也没有 `SCOPE-SETTLING`**（现读确认）。
- **两样东西实际在** `cmd/wisp/resident_voice_engines_windows.go:477-487` 的 `scopeOptions()` 里：
  ```go
  477  func (r *voiceResidency) scopeOptions() session.ScopeOptions {
  478      hooks := session.Hooks{
  479          Logf: func(format string, args ...any) {
  480              // 票 12 AC#9 的 T1：scope 自己宣布进 Settling（SCOPE-SETTLING 是它的
  481              // 审计行首，D43#31 的 session 半边）的那一刻。锚在行首前缀上，行文
  482              // 漂移只会让戳停掉、被回落测试当红抓住，不会给错数。
  483              if strings.HasPrefix(format, "session: SCOPE-SETTLING") {
  484                  r.settle.markSettling()
  485              }
  486              slog.Info("voice[scope]: " + fmt.Sprintf(format, args...))
  487          },
  ```
- ⇒ **`markSettling` 的唯一调用点就是 `:484`**（`git grep -n 'markSettling()'` 只回两处：定义 `settle_trace_12r2.go:78` 与这一处）。

## 要做的（按序）

- [ ] **AC#0 先复现＋读清扫描范围（⛔ 不许直接改）**：
  1. 用 carrier 跑一次，贴逐字红句：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test12r2SeamCallSitesAreTheMeasurement$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **逐字读 `settle_trace_12r2_windows_test.go` 的 `:160`–`:200`**，写清：T1 检查把 `ast.Inspect` 的范围**锚在哪个 `FuncDecl`** 上？是 `openScope` 还是别处？**这决定修法，⛔ 不许猜。**

- [ ] **AC#1 把检查锚点移到产码真实所在（⛔ 不许删检查、不许放宽成"整个文件里出现过就算"）**：
  - ⛔ **禁止**改成"在整个 `resident_voice_engines_windows.go` 里 grep 到 `markSettling` 就算过"——那会把"**必须在 `SCOPE-SETTLING` 前缀判断的同一分支里**打戳"这个真正的承诺**丢掉**，等于把检查废掉。
  - 正确改法（二选一，你按"最符合主流 harness"自己拍并写理由）：
    - **甲**：把 T1 的扫描范围从 `openScope` 的 `FuncDecl` 改成 **`scopeOptions` 的 `FuncDecl`**，判据不变（`markSettling` 与 `"session: SCOPE-SETTLING"` 都要在同一函数体里出现）。
    - **乙**：更强——要求 `markSettling` 的调用**在 `if strings.HasPrefix(format, "session: SCOPE-SETTLING")` 这个 `IfStmt` 的 `Body` 内**（即真的钉住"打戳发生在这个前缀判断的分支里"，而不只是"同一函数里两样都在"）。
  - **推荐乙**（更强、更贴近产码注释里那句"锚在行首前缀上"），但**你自己拍并写清理由**。

- [ ] **AC#2 突变控制实测（负向必响，这一发必须真跑）**：
  用 `go test -overlay`（拷贝放**仓外**，⛔ 不改仓内产码）做**两发**：
  - **发①**：把 `:484` 的 `r.settle.markSettling()` **删掉** ⇒ 新断言**必须红**。
  - **发②**：把 `markSettling()` 从 `if` 分支里**挪到 `if` 外面**（同函数内、但脱离前缀判断）⇒ **若你选乙，必须红**；若你选甲，此发可绿（这正是甲乙的差别，写清即可）。
  两发读数都贴进 Progress。⛔ 不许只写在注释里说"删了就会红"。

- [ ] **AC#3 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／**任何产码 `.go`（除本测试文件）** ⇒ **直接退回**。

- [ ] **AC#4 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · `tools/d22scan` 不新增违规。四条读数全贴。
  ⚠ 已知残留：`build/mut_rvw.go` 是**未跟踪、被 `.gitignore` 吃掉的残留件**，会让 `go build ./...` 红。**⛔ 不是你的写面，⛔ 不许删**；用 `go build ./cmd/wisp/`。

- [ ] **AC#5 与票 12 的关系写清**：本票只修这一枚的检查锚点，⛔ 不翻票 12 任何框，⛔ 不宣称 AC#9 的 T0/T1/T2 三戳全绿（T2 那半是 `:196-199`，另计）。

## 本票**不**解决

- 不碰 `Test12r2SettleWalkFiveSamplesEndToEnd`（它已 PASS 22.66s）。
- 不碰产码里 `markSettling` 的调用位置（那是既定实现，改它＝改产品行为，需单独立票）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 写面：**仅** `cmd/wisp/settle_trace_12r2_windows_test.go`。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push（推送由编排者做）。
- ⛔ 翻任何 AC 框（翻框由编排者做，非实现者）。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\427-<腿码>\`。
- 跑测试**只经 carrier**，⛔ 裸 `go test ./cmd/wisp`（DLL 遮蔽会给出假读数）。
