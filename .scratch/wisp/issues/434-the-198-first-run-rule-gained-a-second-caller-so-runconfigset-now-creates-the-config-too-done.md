# 434 — `TestTicket198FirstRunCallerIsTheRunEntryOnly` 真红：`ensureFirstRunConfig` 多了一个调用者 `runConfigSet`，首跑创建从此有两个入口

- Status: **done**（写面＝`cmd/wisp/run.go` 的 `runConfigSet` **那一处**，或 `cmd/wisp/firstrun_198_test.go` 的断言——⚠ **先只读取证，判清是"产品多了入口"还是"断言该承认这个入口"，再定改法**）
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

- [x] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴逐字红句与耗时：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestTicket198FirstRunCallerIsTheRunEntryOnly$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **抄出票 198 的 AC J1 原文**（逐字，⛔ 转述）。
  3. `git log --oneline -10 -- 'cmd/wisp/run.go'`，找**哪个提交给 `runConfigSet` 加了那个调用**，贴出那个提交的说明。
  4. 读 `runConfigSet` 的完整函数体，写清它**为什么**需要 `ensureFirstRunConfig`（或者为什么不需要）。
  5. 判甲／乙，给一句话结论＋`文件:行`＋AC J1 原文引用。

- [x] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲**：从 `runConfigSet` 去掉那个调用。⚠ **必须确认 `config set` 在没有配置文件时的行为仍然合理**（贴出改后的行为）。
  - **乙**：改断言承认两个调用者，⛔ 不许放宽成"只要不含常驻腿就算过"（那会丢掉"恰好两个"这个约束）。
  - **必须配突变控制真跑**：甲 ⇒ 把调用加回去必须红；乙 ⇒ 把 `runConfigSet` 从允许名单去掉必须红。
  - ⚠ **仪器**：**这枚读磁盘源码**，`-overlay` **无效**（台账 `A1804`）。必须用**仓外磁盘改写**，并说明你怎么确认仪器有效（例如同一二进制读两棵不同树给出两种判定）。

- [x] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [x] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [x] **AC#4 与票 174 的关系写清**：票 174 动过 `run.go` 的 `[fs]` 授权面（已回退，`7dcb145ff`）。**确认本票的改动与那处无关**，⛔ 顺手碰 `[fs]`。

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

## Progress log

- **2026-11-12 23:0x 写腿 `434-firstrun198`：判乙，断言侧修完（AC 框一枚未翻、未 push）**
  - **AC#0 只读**：carrier 复现红句逐字 `AC J1 RED: ensureFirstRunConfig is called by [runConfigSet65 runTextTask], want exactly [runTextTask]`、**0.44s**。AC J1 原文逐字（票 198 `:13` ＋台账 `A514:10469`）＝**"首建只挂在 `wisp run` 这一条「用户明确发起」的入口；⛔ 不许进 `assembleRuntime`；⛔ 不许改那两枚钉"**，理由具名是"常驻腿直达装配根"。
  - **加调用的提交＝`ef00174ac`**（2026-11-07 20:29，票 65 CLI 切球路径），说明里那句逐字：**"缺文件set走ensureFirstRunConfig唯一首建口,get只答默认不建文件"**。它是新入口**出生自带**的一跳（该提交新增 `config_cmd_65.go` 236 行），不是某腿顺手多写。
  - **判乙**：J1 原文的禁止对象是"进 `assembleRuntime`"（理由＝常驻腿直达），`runConfigSet65`（`cmd/wisp/config_cmd_65.go:173`）只从 `main.go:149-157` 的独立 `case "config"` 来、全程不碰 `assembleRuntime` ⇒ **不违反原文**；被违反的是断言自己硬化的"计数＝1"。甲会红掉票 65 自己的钉子 `config_cmd_65_test.go:84-86`（逐字 "set on missing file must first-run-create"）⇒ 甲是拆别人票已验收的判据。
  - **产品语义判定**：`wisp config set` 在无配置文件时**该创建**。依据＝①`config.NewManager` 第一句就 `LoadFile`，缺文件即错（`manager.go:111-123`），不建就恒不可用；②创建能力/形状已由票 198 定死（`SaveFile(NewDefaults())`，单一默认源，禁第二真相源），复用 `ensureFirstRunConfig` 是唯一合规写法；③主流形状同向（`git config`／`kubectl config set` 都建文件，正是票 65 自陈的选形依据 `config_cmd_65.go:12-18`）；④`get` 读不建／`set` 写建的分界是票 65 明写已验收的设计（`config_cmd_65.go:161-163`）。
  - **修法（只改一处＝断言）**：`cmd/wisp/firstrun_198_test.go` 的 `want` 改成精确集合 `[]string{"runConfigSet65","runTextTask"}` ＋ `slices.Sort`（读数来自 map，两个调用者后无序比较＝抛硬币）＋ `slices.Equal`。**⛔ 没放宽成"不含常驻腿就算过"**：多一个红、少一个也红。
  - **突变控制（仓外磁盘改写，⛔ 非 overlay）**：台件 `D:\tmp\wisp\20261112\434-firstrun198\` 三棵 269 枚 `.go` 全拷贝树，**同一个二进制**（`sha256 89A527C3…D183F`，三发前后逐字相同）只换 `cd`：`tree-pristine` ⇒ **PASS 0.31s**；`tree-no-set-call`（删 `config_cmd_65.go:173-176`）⇒ **FAIL**，红句 `…called by [runTextTask], want exactly [runConfigSet65 runTextTask]`；`tree-assemble`（种进 `assembleRuntime`）⇒ **FAIL**，红句 `…called by [assembleRuntime runConfigSet65 runTextTask]`。
  - **仪器有效性怎么确认（三把尺）**：①**同一二进制两种判定**（唯一变量是读哪棵树；若读编译期烘焙内容三发必同色）；②红句里的**集合逐字跟着树变**；③红句**行号从 `:281` 漂到 `:305`**（因为我同文件加了注释行）⇒ 读的是盘上当前那份文件，不是缓存。⚠ 先踩一坑记我：PowerShell 直接 `& $exe -test.run …` 会把旗标吞掉、三发全 rc=2 输出 `-help` 文本；改经 `D:\work\soft\Git\bin\bash.exe -c` 才拿到真读数。
  - **越界**：本腿写面＝`cmd/wisp/firstrun_198_test.go` 一枚（＋本票 log ＋证据件）。`cmd/wisp/run.go` **零改动**（`git diff --name-only -- cmd/wisp/run.go` 空）⇒ 与票 174 的 `[fs]` 授权面（已回退 `7dcb145ff`）**无关**。禁止面（`PLAN.md`／`specs`／`thresholds.go`／golden／`allowlist.txt`／`frontend`／`design`）零命中。⚠ `git diff --name-only` 的全量输出里另有 12 枚**别人的在飞腿**改动（票 426–433 票面、291 探针、247 测试件、前端 fixture、346 派单件），本腿未碰、commit 用显式 pathspec 只提三枚。
  - **门禁四读数**：build `rc=0`／vet `rc=0`／gofumpt `list-count=0`／`scripts/d22scan.sh` **clean**（正控 `PASS=45 FAIL=0 SKIP=0`，分母 `bans #1-5 internal/=403 cmd/=77`、`ban #8 internal/=1004 cmd/=276`）。`build/mut_rvw.go` 未删未碰。
  - **附带复跑**：`^TestTicket198` 整族 **10/10 PASS**（含 `TestTicket198R2J1TheAssemblyRootStillCreatesNothing` 0.03s）／`^TestConfig`（票 65 族）**6/6 PASS**（含逐字要求"set 建件"的那枚）⇒ 乙没拆票 65 的判据。
  - **残余三条**（详证据件 §6）：①测试名 `…IsTheRunEntryOnly` 名不副实但**不改名**（≥12 份证据件／台账按名引用，改名会让历史读数失联），已在函数头注释写清读法；②台账 `:10953` 那句"唯一产码调用点仍是 `run.go:245`"**自本票起过期**，但它描述的真缺口（双击链不建件）没变、本票没动，归票 256 族；③"将来加第三条入口怎么办"＝必须显式改这枚断言，那是有意的摩擦。
  - **证据件**：`docs/evidence/s1/434-firstrun-198-second-caller.md`。**撤销口令：「434 回甲」**。
- [2026-11-13 09:1x +08] **编排者：收件＋三把尺全过＋AC#0–AC#4 五格翻勾＋`-done`**。腿交件 `9661bc6f4`（3 files／+299−5；证据件 `docs/evidence/s1/434-firstrun-198-second-caller.md` 19335B）。**三把尺**：①commit 在；②件 19335B 在盘；③**我自己复跑 PASS 0.28s**。**越界**：只改 `cmd/wisp/firstrun_198_test.go` 一枚，**`cmd/wisp/run.go` 零改动**（与票 174 的 `[fs]` 面无关）。
- **★判定＝乙（断言该承认这个入口），只改断言。** 一句话理由：**AC J1 原文的禁止对象是"进 `assembleRuntime`"、理由具名是"常驻腿直达装配根"**，而 `runConfigSet65`（`cmd/wisp/config_cmd_65.go:173`）只从 `main.go:149-157` 的独立 `case "config"` 来、**全程不碰 `assembleRuntime`** ⇒ 被违反的是**断言自己硬化的"计数＝1"**，⛔ 那条规矩。
- **★AC J1 原文逐字（腿从 `198-…-done.md:13` 抄出）**：`首建只挂在 wisp run 这一条"用户明确发起"的入口；⛔ 不许进 assembleRuntime；⛔ 不许改那两枚钉`。裁定理由（台账 `A514`／`:10469`）逐字：`assembleRuntime 被常驻腿共用 ⇒ 首建放进 assembleRuntime 今天就红 3 处`。
- **★哪个提交加的调用**：`ef00174ac`（2026-11-07 20:29:16，票 65 CLI 切球路径），说明逐字 `缺文件set走ensureFirstRunConfig唯一首建口,get只答默认不建文件`。**是新入口出生自带的一跳**（该提交新增 `cmd/wisp/config_cmd_65.go`，`--numstat` 逐字 `236 0`），⛔ 顺手多写一行。
- **★甲不可取（腿的证据）**：去掉那跳会让 `config set` 在干净机器上**必退码 1**（`config.NewManager` 第一句 `LoadFile`，`manager.go:111-123`），并**红掉票 65 自己的钉子** `config_cmd_65_test.go:84-86`（逐字 `set on missing file must first-run-create config.toml`）——**那是拆别人票已验收的判据**。
- **★`config set` 无配置文件时该创建（四条依据）**：①不建就恒不可用（`manager.go:111-123` 注释逐字 `first-run file creation is the caller's concern`）；②创建能力/形状**已由票 198 定死并验收**（`SaveFile(cfgPath, NewDefaults())`，默认值单一来源，禁第二真相源）；③**主流形状同向**：`git config user.name`／`kubectl config set` 都建文件（正是票 65 自陈的选形依据 `config_cmd_65.go:12-18`）；④**读不建／写建的分界是已验收的设计**（`get` 缺文件答默认值 rc 0，`:131-156`；`set` 建，`:158-176`）。
- **★突变控制＋仪器有效性（腿做得硬）**：**同一个二进制**（`sha256 89A527C3…D183F`，三发逐字相同），只换读哪棵仓外树（各 269 枚 `.go` 全拷贝）：`tree-pristine` ⇒ **PASS 0.31s**；`tree-no-set-call`（删 `config_cmd_65.go:173-176`）⇒ **FAIL**（红句 `…called by [runTextTask], want exactly [runConfigSet65 runTextTask]`）；`tree-assemble`（种进 `assembleRuntime`）⇒ **FAIL**（`…called by [assembleRuntime runConfigSet65 runTextTask]`）。**仪器有效性三把尺**：①同一二进制两种判定（唯一变量是读哪棵树）；②红句集合逐字跟着树变；③**红句行号 `:281`→`:305`**（同文件加注释行所致）⇒ **读的是盘上当前文件，⛔ 缓存**。
- **★腿踩坑记（值得记）**：PowerShell 直接 `& $exe -test.run …` 会**把旗标吞掉**、三发全 rc=2 输出 `-help` 文本；改经 `D:\work\soft\Git\bin\bash.exe -c` 才拿到真读数。
- **附带复跑**：`^TestTicket198` **10/10 PASS**、`^TestConfig`（票 65 族）**6/6 PASS** ⇒ **乙没拆票 65 的判据**。
- **★腿留的残余三条（值得后续票）**：①测试名 `…IsTheRunEntryOnly` **名不副实但⛔ 改名**（≥12 份证据件按名引用，改名会让历史读数失联），已在函数头注释写清读法；②台账 `:10953` 那句"唯一产码调用点仍是 `run.go:245`"**自本票起过期**（但真缺口"双击链不建件"没变、本票没动，归票 256 族）；③将来加第三条入口**必须显式改这枚断言**——那是有意的摩擦。

