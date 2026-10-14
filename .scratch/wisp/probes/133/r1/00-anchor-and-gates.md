# 133-r1（落地腿）：给"分发跳没人走"补新尺 —— AC#1 五发×门开/门关、AC#3 自证腿、AC#4 静默缩短、AC#5 门禁、AC#6 CI

- **锚**：开工 `git rev-parse --short HEAD` = `a7b88997`（与派发锚逐字同，零漂移）。仪器 commit `da389ca5`（本腿新增第二枚仪器文件＋同尺增补）；**⛔ 全程只 commit 未 push**。
- **身份**：实现者。⛔ 未翻任何框、⛔ 未动票面 `Status:`/票名；判语归非实现者。
- **可交付物**：`cmd/wisp/dispatch_hop_selfproof_133_test.go`（新文件：AC#3 自证腿＋AC#4 判据）＋ `cmd/wisp/leg_dispatch_gate_133_test.go`（同尺增补：下限 N=12、注册表两半成对、注册表↔census 双向差、逐枚可区分）。两枚文件都在 CI 既有的 `cmd/wisp` step 里跑（§5）。
- 本目录 logs（7 枚上限内的 6 枚）是**逐字摘录**；全量原档在仓外 `D:/tmp/wisp133r1/logs/`（只建不删），摘录里每段都写明原档名。

## 0. 环境、分母与宿主噪声（先登记，后面所有读数都在这条底线上比）

- `GOTMPDIR='C:\tmp\wisp133r1gotmp'`（派单逐字形，先自建）；`PATH` 前置 `third_party/sherpa-onnx`（MSYS 形）；`go1.27.1 windows/amd64`。
- **宿主底噪 7 枚**（pristine 树、零变异、`-count=1 -v` 全包，原档 `01-pristine-7reds-rerun.log`）：`TestAC2SealNoticeLandsInTheRunLegLogFile`、`TestAC4FocusReturnToPriorWindowGap33r5`、`TestAC14GoSideEvalPushReachesThePage`、`Test258BridgeMutationNoSrcKeepsOldBinding`、`Test258OccupiedCombinationNamesTheNewValue`、`Test258V1ProbeSummonEditRebindsLiveBall`、`TestAC246DevLegIgnoresTheTestTaskInjection`。
  - **归因①（根/ACL 系 2 枚）**：把 `GOTMPDIR` 换成 `%LOCALAPPDATA%\wisp-test-gotmp\cli-133r1-attr` 后这两枚转绿（原档 `02-localappdata-root-7reds.log`；`TestAC2SealNotice...` 的红是 `records naming "seal cleared principals..." = 2, want 1`，temp 根继承 ACE 不同 ⇒ 多一条 seal 记录）⇒ **票 311 同形**（"gate 的颜色是根的函数"），与派单指定根绑定，**不是本腿引入**。
  - **归因②（热键/真窗系 5 枚）**：日志行 `hotkey ... err="Hot key is already registered."`（本机 `wisp.exe`＋微信等占着 `Ctrl+Alt+M/P`）＋ WebView2 真窗 ⇒ 宿主环境分母；CI runner 无这些占位者。
- **314 死区**：本腿起过的 **9 趟全包进程里 2 趟整进程死亡**（无判决断在面板/WebView2 段：`n3-10` 18:31:10→18:36:35 = 317.8s、RUN=186 停在 `TestShapeA3SecondPagePostStillDelivers`；CI 形 step 第 1 次 19:35:23→19:41:12 = 306.2s 停在 `TestAC13BringUpSurvivesAReusedThreadQuit`），另 1 趟撞 Go 默认 10m `-timeout`（`panic: test timed out after 10m0s`，603.1s）；其余 6 趟走完（3 趟 306/226/7/3＋2 趟 AC#5＋1 趟 CI 形）。逐趟窗口表在 `logs/04-ac5-gates.txt` 末节；tasklist 自证见同文件头部。⇒ 本腿不把任何单趟四数当判语，主证是红名集合。

## 1. AC#1 五发×两份读数（门开＝131 的门在跑；门关＝`-run` 只跑本尺，不动 131 的文件）

逐发红名/落地行/还原证的逐字摘录：`logs/01-ac1-five-shots.txt`；X14 两拍＋no-doors：`logs/02-...txt`。

| 发 | 变异（落地行，sha256 见 01 原档） | 门关：本尺 | 门开：两门 | 131 单跑 | 备注 |
|---|---|---|---|---|---|
| N-3 | 摘掉 `main.go` 里 `os.Exit(cmdModels(...))` 调用者（`attachParentConsole()` 留在 case 里） | **rc=1 红名 `TestAC1AC2DispatchHopGate133`**；红点名 `cmdModels (models.go:104)`、`modelStore.handOffModel`、`modelsEnsure` 三枚孤儿＋nail-hop | rc=1 双红 | rc=1（红） | 三态见下 |
| X4 | `if len(args) > 0 && args[0] == "--diag"` 早退腿＋`diag131x.go` | **rc=1 红名本尺**，红名点到 `leg "--diag" (main.go:88)` | rc=1 双红 | rc=1（红） | |
| X8 | `case sloCmdName131:`（1 标签换 const） | **rc=1 红名本尺**：非字面量标签＋裁决句孤儿＋usage 差 | rc=1 双红 | rc=1（红） | 不是"整包坏掉才响"：账本仍出全 12 行为底、只这一行变合成腿 |
| X12 | `var sinkAlias131 = installLogSink`＋`case "fake131"` | **rc=1 红名本尺** | rc=1 双红 | rc=1（红） | 别名已被本尺走图接住（131 的 aliases 是它的事，本腿自查） |
| X14 第一拍 | `sfx131.go`：结构体字段初值 `holder131.open=installLogSink`，实二进制真活 | **rc=1 红名本尺**：`leg "sfx131" ... reaches installLogSink on this path: dispatch -> holder131.open -> installLogSink` | rc=1（131 门此时 FAIL——本发种子里带 slog） | 带 slog 版：FAIL `records=true -> RED records...` | 票面原形（无 slog）另跑：**x14c** ⇒ 本尺红、**131 单跑 PASS**（见下） |
| X14 第二拍 | 字段初值与那次调用删掉（腿仍在） | **rc=1 红名本尺**：`leg "sfx131" ... covered by nothing` | rc=1 | 带 slog 版：FAIL；**x14b2c 原形：131 PASS** | 不许用 X12 已红抵这发：本读数独立 |

**X14 原形的独占目击（票面 §2.0.1 那句"只有 X14 两拍是本尺独占目击"的今日复核）**：x14c/x14b2c（无 slog 版）下 131 的账本逐字 `leg sfx131 main.go:89 install=false records=false ruled=false nails=- -> no records` ＋ `--- PASS`；同一刻本尺门关着 rc=1 红。⇒ 这一形今天仍然**只有本尺看得见**。

**N-3 三态（票面字面要求＋09-23 §2.0.1 裁定的偏差登记）**：
1. **不加尺**（`-skip '^TestAC1AC2DispatchHopGate133$'`，全包）：rc=1，红名只有 131 的门（`TestAC4EveryLegIsNailedOrRuled`）＋宿主 3 枚 ⇒ 字面要求"全绿"今天**不可能成立**（131 已接住这一跳），原档 `n3-10-no133instr.log`；**两门都不跑**（`-skip '^(本尺|131门)$'`）⇒ 306/226/**7/3**，7 枚全是宿主底噪、无一门一提及 ⇒ **"没人沿边走"这一形照旧**（原档 `n3-13-nodoors.log`）。
2. **装尺**：见上表（门关 rc=1、红名本尺）。
3. **还原复绿**：`main.go` sha256 回 `24efbffb...`（`identical=True`）、`git status --porcelain cmd/wisp/` 只剩本腿两枚仪器文件；复跑四门 rc=0 全 PASS（原档 `04-restore-green.log`）。
- **方法值那一支（`v.M(...)`）未造**（包内没有可借形状）⇒ 按票面口径**不算已测**；**跨包那一支出 `cmd/wisp/**` 地界，本票不接**（逐字照票面）。

## 2. AC#3 自证腿（变异自证，不是注释声明）：`logs/03-ac3-ac4-mutations.txt`

新文件把每枚读者喂一份**已知坏形状的夹具包**（`t.TempDir()` 里的迷你 cmd/wisp，`loadPackage133` 原样读它），并断言它必须响；夹具孪生腿证明 distinctness 判据能答"否"。两次**独立破坏**读数：

- **m1（R-121-2 同形：成员判定恒真）**：`isRunnableCase133` 顶部插 `return true`（1 行）⇒ **`TestAC3SelfProof133TheDispatchReadersCanAnswerNo` 红**：夹具里的方法 `TestPhantomMethod133` 被洗成覆盖（`covered=test TestPhantomMethod133`），自证腿当场点破"方法不是用例"。⚠ 同一发门自己**PASS** ⇒ 没有这枚自证腿，这一形在这一棵树上就是哑的。
- **m2（Y3 同形：判定直接取期望）**：`runRosterReds135` 里 `if roster[name]` 改 `if true`（1 行，对账恒真）⇒ **自证腿红**（"a roster WITHOUT the credited case has to red the row..."），门本身仍 PASS。
两发都出仓外备份、逐发还原（sha256 identical 逐字在 03）。**"摘掉自证腿 ⇒ 其余用例不会全绿"**：m1/m2 两读数即为反证（门 PASS、只有自证腿红）。

## 3. AC#4 静默缩短要能响：三枚清单、双向差、逐枚可区分、Y2/Y3

- **① 双向差集为空**：`censusVsUsage133`（门内既有）把"`func main` 分支现算的腿清单"与"给人看的 `usage` 文字块"**双向**对差，任一方向非空即红。今日 12 legs ↔ usage 双向空；**y1drop**（删整枚 `panel-inbound` 腿，Y1 形）⇒ 双向同时咬（见 ③）。注册表这半边本腿新增：每条 claim 的 leg 必须仍在 census（改名/删腿会留一枚指向不存在命令的 claim，Y1 形），健康树零红。
- **② 逐枚可区分（Y2/Y3 判据）**：门新增 `distinctnessReds133`（本文件同尺调用）：两枚腿不许穿同一行签名（`installs|handoff|covered|entries|aliases`）、覆盖列不许整表一个值（恒真即红）。今日读数：**12 legs / 12 distinct rows / kinds=[nail ruling test]**。另新增 **nail 两半成对**：claim 点名的用例必须驱动 claim 点名的 entry（**Y2 实体**）。
  - **m3 改前（本腿自己量，前两程都没造过这一形）**：把 models↔secret 两枚 claim 的 `test` 名对调（2 行、腿名/注释/entry/count 全保留）⇒ **门 rc=0 全绿、两门皆 PASS**（`logs/03` m3 段）——正是 Y2"数量对、归属错"。
  - **m3b 改后**：同一发 ⇒ **门 rc=1**、两枚红逐字点名 `leg "models"'s nail pairs the case "TestAC3Secret..." with the entry "cmdModels", and that case never calls "cmdModels"`（secret 对称一枚）。⇒ Y2 形自此被接住。
  - **no-args 那枚的例外**：它的用例是**进程形**（真启动二进制，不 in-process 调 `runResident`）⇒ 注册表加显式 `process: true`，账本逐字披露 `(process claim, name-only: ...)`，**不 waived 成绿、不假装成对已查**。这不是漏，是把"查不了"写在行上（131 自己的账本对同一枚钉写 `subprocess:` 同义）。
- **③ 下限 N＝12（写死并给实测值）＋"全瞎必须红着说看不见"**：`minLegs133` 由 11 抬到 **12**（今天实测：10 枚字面标签〔run providers doctor secret models slo panel-assets panel-inbound version help〕＋no-args＋default）；谓词抽成 `floorReds133`，**门与自证腿用同一个数**。
  - **y1drop**（删一枚腿 ⇒ 11）：门 rc=1，逐字 `the leg census enumerated 11 legs (12 is the floor written down next to minLegs133)` ＋ usage 方向"documents command \"panel-inbound\", which func main does not dispatch"。
  - **y2blind**（`args := os.Args[1:]` ⇒ `args := os.Args`，解析器整个瞎）：门 rc=1，**红着说** `func main binds no local to os.Args, so this gate cannot tell a dispatch from a call` ＋ 0 legs 撞下限 ⇒ **不许返回空清单当"没有腿要测"**。
  - ⚠ **票面 AC#4 第三枚清单（"131 那扇枚举门的腿清单"）的落点说明**：131 的门（`cmd/wisp/leg_sink_gate_131_test.go`）本腿**一字未改**——派单明写"修复须动 131 的门 ⇒ 停手上报"，且本腿的地界是新文件。对"同一枚 census 无下限"这一要求，本尺以**自己的** `floorReds133`（同一个 `main.go` 现算的 12 条腿）＋全瞎红覆盖（③上面两条读数），**不借 131 的门**。131 自家门内那三行（`len(legs)==0` 的 Fatal 有、`>=N` 下限无）留给它的 owner，本腿只登记事实。

## 4. AC#5 门禁：`logs/04-ac5-gates.txt`（四数只从 `-v` 量；计数口径＝仓内 `tools/d22scan/runtests.sh` 那四把 grep）

- **字面命令 `go test -count=2 -v ./cmd/wisp/`**：603.108s 处 `panic: test timed out after 10m0s`（Go 默认 `-timeout`；本机单趟就要 5–10 分钟）⇒ 字面跑不完，按 314 如实报；同一断言集换 `-timeout 40m` 跑完（1765.788s）：**RUN=824 PASS=575 FAIL=23 SKIP=6**（两趟），红名＝7 枚宿主底噪 ×2 ＋ 9 枚负载诱发面板/STA 超时（15s/17.95s 那种，R-133-12 同族）——本机四数不是稳定量，主证是红名集合，低负载同树单趟仍 **306/226/7/3**（7 枚全是宿主底噪）。
- **CI 形 step `bash scripts/wisp-cli-tests.sh`**（run Id 链见 §5）：第 2 次跑完读 **RUN=282 PASS=217 FAIL=4 SKIP=1**；4 枚红＝宿主底噪 3（Seal notice/Focus return/GoSideEval）＋ `TestAC1ResidentLegInstallsItsLogListenerOnDisk`（19.64s，R-133-12 负载族——低负载同树它是绿的）。**本腿四枚用例在该 step 的分母里全 PASS**（`TestAC1AC2DispatchHopGate133`、`TestAC4EveryLegIsNailedOrRuled`、`TestAC3SelfProof133...`、`TestAC4Distinctness133...`）。
- **`gofmt -l cmd/wisp/`**：只剩 `cmd/wisp/models.go`（gofmt 与 gofumpt 都点它）；**归因**：该文件与 HEAD blob 只差行尾（盘上 CRLF、blob LF，`cmp` 第 13 字节即差）⇒ 工作树 CRLF 伪影、`git status` 看不见；同日前一版快照（`git archive da389ca5^`）gofmt 干净、gofumpt 只点 `panel_transport_35r2_test.go`（commit `18d43f8e` 自述的 "HEAD 真格式债"）。**两枚都不在本腿 pathspec 里**；本腿两枚文件逐枚干净。`gofumpt` 版本 `v0.12.0 (go1.27.1)`（CI 装的是 `@latest`，跨版本可比性仅此一注）；`go vet ./cmd/wisp/` **rc=0**。
- **`sh scripts/d22scan.sh` rc=0 纯净**；各 scope（含本腿两枚文件）：`bans#1-5 internal/=237 cmd/=40 · #6 frontend/=115 · #7 internal/tools/=24 · #8 design/=39 frontend/=115 internal/=556 cmd/=128`。**同日控制组**（`da389ca5^` 的 `git archive` 快照）：除 `#8 cmd/` 127→128（恰为本腿新增的 1 枚 `_test.go`）外逐数相同 ⇒ **不降**有同日基线（控制组非 git 仓、gitignore 不生效，`#8 design/` 30 vs 39 是姿态差，已登记）。
- **票 123 那批 CLI 用例（"审批超时 1/300 秒未确认 ⇒ 一律判拒绝"）本腿零触碰、零放宽**（04 的逐名字段里它既没被 -skip、也没被改阈值）。

## 5. AC#6 仪器自己进 CI：`logs/05-ac6-ci.txt`

新文件是 `cmd/wisp` 的 `_test.go` ⇒ 随既有 step 跑：job **`test-windows`**（`.github/workflows/ci.yml:533`）里的 step **"cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"**（ci.yml:601）→ `bash scripts/wisp-cli-tests.sh`（ci.yml:621）→ `portable-tests.sh --scope=cli` → `scope=(./cmd/wisp/)`＋`cli_pin=github.com/CarlosShao/wisp/cmd/wisp`。**无需改 `.github/**`**（若需改则按派单停手上报）。真 run 读数（`gh`，本腿 2026-10-14）：**run `38554642217` / sha `72db0f39`，job `test-windows` databaseId `115820641622`，step #7 同名**，该 step 自己的行 `four numbers: === RUN=410 --- PASS=297 --- FAIL=2 --- SKIP=1`（红是那棵树的状态，非本腿）；另一条 `38551683523` 同 step 亦跑过。本腿自跑该脚本两次：一次 314 死（306.2s）、一次跑完（282/217/4/1）——**本腿四枚用例两次都在分母里、跑完那次全 PASS**（§4）。

## 6. 没做的、偏离与停手上报项（逐条具名）

1. **方法值 `v.M(...)` 未造**；**跨包 install 支**出地界未接（票面原文照抄）。
2. **AC#4 第三枚清单的"131 门内下限"未在 131 的文件里加**（停手上报规则的直接后果，§3③ 已写清本尺自己的替代覆盖）。
3. **AC#1 三态里"不加尺全绿"字面不可满足**：与 09-23 §2.0.1 同结论（131 已接住 N-3/X4/X8/X12），本腿以 no-doors 读数（§1）给出"两门都摘仍零门提及"的第一手复核，X14 两拍由 x14c/x14b2c 承担独占目击。**读数与票面 09-23 记录不同的一处**：那一次 X14 种子账本读 `install=false records=false`（无 slog），与本腿 x14c/x14b2c 逐字同；本腿主发 x14/x14b2 的种子里带 slog ⇒ 131 也红（不同种子，已并列登记，不算分歧）。
4. **AC#2 本腿未动**（已由编排者翻勾）；**票面框一枚未翻**。
5. 全包在 `C:\tmp` 根下带 7 枚宿主底噪（§0），AC#5 四数**不是 0/0**，逐枚归因已给；"不降"本腿另取**同日控制组**（`da389ca5^` 的 `git archive` 快照，除 `#8 cmd/` +1〔本腿新文件〕外逐数相同），但仍比不出"进仓历史基线"（`R-133-7` 口径照旧）。
