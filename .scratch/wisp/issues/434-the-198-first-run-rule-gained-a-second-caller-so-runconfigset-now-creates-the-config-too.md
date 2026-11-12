# 434 — `TestTicket198FirstRunCallerIsTheRunEntryOnly` 真红：`ensureFirstRunConfig` 多了一个调用者 `runConfigSet`，首跑创建从此有两个入口

- Status: **ready-for-agent**（写面＝`cmd/wisp/run.go` 的 `runConfigSet` **那一处**，或 `cmd/wisp/firstrun_198_test.go` 的断言——⚠ **先只读取证，判清是"产品多了入口"还是"断言该承认这个入口"，再定改法**）
- 来源：round 89 整包（`D:\tmp\wisp\20261112\round89\full.log`，`PASS=583 / FAIL=11`）＋台账 `:32512` 存量真红
- 关联：票 198（AC J1：首跑创建的唯一调用者）· 票 174（`[fs]` 授权面，同族）· 台账 `A1806`

## 这是什么（人话）

有个规矩：**"第一次运行时创建配置文件"这件事，只许从"运行入口"这一条路发生**。
理由是：常驻腿（resident leg）会直接走到装配根，如果它也能触发首跑创建，
那么"首跑"就不再是"用户第一次运行"，而会变成"任何一条路走到装配根"——**语义被摊薄了**。

现在实测：**有两个调用者**——`runTextTask`（本该是唯一的）和 **`runConfigSet`**（多出来的）。
⇒ 规矩被打破了。

## 现量（逐字，可复跑）

**红句**（`cmd/wisp/firstrun_198_test.go:281`）：
> `AC J1 RED: ensureFirstRunConfig is called by [runConfigSet65 runTextTask], want exactly [runTextTask] - the resident leg reaches the assembly root directly, and first-run creation shared…`

**耗时 0.24s**（纯 AST 静态检查，不跑产品）⇒ **稳定，与负载无关。**

## 机制链（逐字，⛔ 凭印象）

- 测试：`cmd/wisp/firstrun_198_test.go:275-282`
  ```go
  for name, fi := range pkg.funcs {
      if fi.calls["ensureFirstRunConfig"] { callers = append(callers, name) }
  }
  if len(callers) != 1 || callers[0] != "runTextTask" {
      t.Errorf("AC J1 RED: ensureFirstRunConfig is called by %v, want exactly [runTextTask] - …")
  }
  ```
- **这是 AST 静态检查**（扫 `run.go` 的函数调用图），⛔ 运行时行为。
- ⇒ **仪器选择**：**`go test -overlay` 对这种检查无效**（它读磁盘源码，overlay 只换编译期输入）——见台账 `A1804` 那条分界线。**必须用仓外磁盘改写。**

## ⚠ 必须自己判清的核心岔路（⛔ 不许我替你定）

**"`runConfigSet` 该不该调 `ensureFirstRunConfig`"——两种答案，修法相反：**

- **甲：产品多了入口（真回归）** ⇒ `runConfigSet`（`wisp config set` 那条路）**不该**触发首跑创建，改产码去掉那个调用。
  **判据**：读票 198 的 AC J1 原文与 `docs/specs/**` 里关于"首跑"的定义。若规矩是"只有运行入口"，那 `config set` 确实不该创建。
  **⚠ 但要想清楚**：`wisp config set` 在**从没跑过 wisp** 的机器上，如果配置文件不存在，它**该失败还是该创建**？这是产品语义问题，⛔ 拍脑袋。
- **乙：断言该承认这个入口** ⇒ 若 `config set` 本来就需要一个配置文件才能工作、且产品设计允许它创建，那断言要改成 `[runConfigSet runTextTask]`。
  **判据**：AC J1 的原话到底约束的是什么——"唯一调用者"还是"常驻腿不得调用"？

**★腿必须先找到 AC J1 的原文**（票 198 票面 + 若有对应 spec），**然后按原文判**，⛔ 按"哪个改起来容易"判。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴逐字红句与耗时：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestTicket198FirstRunCallerIsTheRunEntryOnly$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **抄出票 198 的 AC J1 原文**（逐字，⛔ 转述）。
  3. `git log --oneline -10 -- 'cmd/wisp/run.go'`，找**哪个提交给 `runConfigSet` 加了那个调用**，贴出那个提交的说明。
  4. 读 `runConfigSet` 的完整函数体，写清它**为什么**需要 `ensureFirstRunConfig`（或者为什么不需要）。
  5. 判甲／乙，给一句话结论＋`文件:行`＋AC J1 原文引用。

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲**：从 `runConfigSet` 去掉那个调用。⚠ **必须确认 `config set` 在没有配置文件时的行为仍然合理**（贴出改后的行为）。
  - **乙**：改断言承认两个调用者，⛔ 不许放宽成"只要不含常驻腿就算过"（那会丢掉"恰好两个"这个约束）。
  - **必须配突变控制真跑**：甲 ⇒ 把调用加回去必须红；乙 ⇒ 把 `runConfigSet` 从允许名单去掉必须红。
  - ⚠ **仪器**：**这枚读磁盘源码**，`-overlay` **无效**（台账 `A1804`）。必须用**仓外磁盘改写**，并说明你怎么确认仪器有效（例如同一二进制读两棵不同树给出两种判定）。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 与票 174 的关系写清**：票 174 动过 `run.go` 的 `[fs]` 授权面（已回退，`7dcb145ff`）。**确认本票的改动与那处无关**，⛔ 顺手碰 `[fs]`。

## 本票**不**解决

- 不碰 `Test313*`／`TestLockedSuffix*`／`TestAC247*`（另计，会另派）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\434-<腿码>\`。
- **⚠ 本票与真窗族测试（`Test258*`／`TestAC4Focus*`）⛔ 并行**——它们要 `RegisterHotKey`，会互相污染（台账 `A1807`）。
