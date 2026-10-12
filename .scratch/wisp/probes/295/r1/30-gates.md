# 票 295 · 落地腿 `295-r1` · 第 3 笔：`AC#4` 门禁 ＋ `AC#5` 代价表

时刻 `2026-10-12 15:5x +08`。⛔ 走 `scripts/*.sh`、⛔ 走 `tools/d22scan/runtests.sh` 跑我的测试（`310-r2` 在飞态；本腿自查时 `git status --porcelain -- scripts tools/d22scan`＝空，仍按派单 §3 直接 `GOTMPDIR=… go test …`）。
专用根全族＝`$LOCALAPPDATA/wisp-test-gotmp/295r1-<YYYYMMDD-HHMMSS>-<pid>-<档名>`（**默认那一族**，逐字在每发件头）；⛔ 指到仓所在卷 `D:`（票 311 实测换卷改颜色）。
⚠ 派单 §3 的归档清单少一枚：`git archive … cmd/wisp internal tools scripts go.mod go.sum` 导出的树里 `cmd/wisp` **[setup failed]**（`internal/panel/assets.go:25` 引 `github.com/CarlosShao/wisp/frontend`，⛔ 在那六项内）⇒ 本腿加一枚 `frontend` 后重跑，逐字与两把计数见 `logs/10`（作废的那发）与 `logs/11`/`logs/12`；两棵成对树（改前⇄改后）都用同一套七项，差名册＝5 枚（我 4 枚 ＋ `scripts/winsec-tests.sh`，后者＝别腿 `3b857e53` 在两锚之间入库）。欠账 `Q-295-r1-archive-set`。

## 1. `AC#4` 逐条（每条带尺＋射程＋量的/算的＋`rc=`）

| # | 命令（逐字） | 读数 | rc |
|---|---|---|---|
| g1 | `GOFLAGS= go build ./...` | 全仓工作树，改后跑两遍（`logs/20`、`logs/41` 前）皆无声＝0 行输出 | `rc=0` |
| g2a | `gofumpt --version`／`which gofumpt`（裸名那把） | `which: no gofumpt in (…)`、`gofumpt: command not found` ⇒ **具名 `rc=127`**（⛔ 当"跳过"） | `rc=127` |
| g2b | `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` | 工作树口径 **2 枚**＝`cmd\wisp\models.go`、`cmd\wisp\panel_transport_35r2_test.go`；**改前同一把＝同 2 枚**（重建树 `cp -r cmd/wisp` ＋ 4 枚还原成 `git show HEAD:` 版 ⇒ `logs/20` 那一段）⇒ 本改动**净增 0 枚**；我那 4 枚单量＝**0 枚** | `rc=0` |
| g3 | `gofmt -l cmd/wisp` 改前／改后并排 | 改前・工作树字节口径＝**1 枚**（`models.go`）；改后・工作树口径＝**1 枚**（同 `models.go`）；改前・HEAD-blob／LF 导出树口径＝**0 枚** | 三把各 `rc=0` |
| g4 | 七发 `go test -count=1 -v ./cmd/wisp/ ./internal/audio/` | 见下面第 2 节（改前两发＝导出树；改后五发＝C 作废、F/G 成对导出树、D/E 工作树） | 各发 `GO_TEST_RC` 见件 |
| g5 | `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations`；分母逐枚在场（`internal/` 231、`cmd/` 39、ban #8 `cmd/` 119 枚含注释与 `_test.go`）；`Binary file … matches` 尺＝**0 命中**（无需 `-a` 重跑，两把的差因此为空，具名） | `rc=0` |
| g6 | `GOFLAGS= go build ./...`（漂移修复后再一遍）＋ focused 复绿发 | `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` ⇒ `--- PASS (0.00s)`；整发 `-run 'TestTicket255|TestAC247|TestAC290'`＝33 枚顶层 PASS＋1 SKIP、`ok cmd/wisp 8.766s` | `DRIFTFIX_RC=0` |

`models.go` 那一枚 g3/g2b 的假枚成因（量的）＝工作树字节里有 334 枚 CR（尺＝`tr -cd '\r' | wc -c` 对工作树量）而 HEAD blob 里 0 枚；`.gitattributes` 写的是 `*.go text eol=lf` ⇒ **本机 checkout 形状**，⛔ 我改的件、⛔ 拿两把的差去指控任何人。我那 4 枚件工作树里 CR＝0。

`g2b` 那 2 枚的归属（⛔ 我修，⛔ 放宽）：`panel_transport_35r2_test.go` 在 **HEAD 里就**被 gofumpt 判需要 `type (…)` 分组（尺＝`gofumpt -d` 那段 diff，导出的 LF 树里同样命中＝⛔ CRLF 假枚）；`models.go` 是上面那枚 CRLF。两枚都⛔ 在本票写面三件套内 ⇒ **原样报回**＝欠账 `Q-295-r1-gofumpt-two`。

## 2. `AC#4` 那七发与"新增红 0 枚"（名册与尺全文＝`logs/46-roster-all-runs.txt`）

计数尺（本仓口径）＝`grep -c -- '--- PASS'`（含子项）；顶层另一把同时打印。

| 发 | 台面 | PASS／FAIL／SKIP | 件 |
|---|---|---|---|
| A | 导出树 @`0615c6c9`（改前第 1 发） | 423／22／4 | `logs/11` |
| B | 同一棵树第 2 发 | 424／22／4 | `logs/12` |
| C | 工作树 @AC#1 第一版（漂移未销＝⛔ 改后凭据） | 449／**4**／2 | `logs/40` |
| F | 导出树 @`1c1d7df3`（改后成对第 1 发） | 424／22／3 | `logs/42` |
| G | 同一棵改后树第 2 发 | 425／21／3 | `logs/43` |
| D | 工作树 @`1c1d7df3`（改后第 1 发） | 450／3／2 | `logs/44` |
| E | 工作树第 2 发 | 449／4／2 | `logs/45` |

- 成对导出树（同台面）：改前 A∩B＝449 行、改后 F∩G＝448 行；**新增红（严尺＝改后 FAIL 名册里、改前⛔ 是 FAIL 的枚数，含 PASS→FAIL 与 SKIP→FAIL）＝`comm -13` ⇒ 0 枚**；反向翻绿 1 枚＝`TestPanelHostRealWindowHopAndLifecycle`（真窗族）。
- flake 三枚逐名具名（交集尺自动不收）：B-only `PASS TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`；F/G 之差 `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`（40.76s 审批超时）；D/E 之差 `TestTicket223HandEditedFsLooseningCostsAnL2Card`（`config_reload_223_test.go:406` 卡片没弹、答复走了控制台监听）。
- 工作树 D∩E 的红名册＝3 枚，逐名 `TestAC4FocusReturnToPriorWindowGap33r5`／`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`／`TestAC14GoSideEvalPushReachesThePage`＝台账 `docs/reports/pending-and-issues.md:14768` 已具名的"真窗族按 4↔5 会漂"那一族；其中 `TestAC13ColdStart…` 在四棵导出树里全 **SKIP**（无 `frontend/dist` ⇒ `entryIDProbes` 空即 `t.Skipf`，机制台账 `:10407` 具名）＝**台面差、⛔ 我的改动**。
- `internal/audio` 那包七发逐字 `ok github.com/CarlosShao/wisp/internal/audio 15.961s—16.519s`＝0 红 0 跳（派单 §1 那格 ② 裁"第二包买的是签名级回归"，这里给的就是那一发的色）。
- 编队尺随行：每发件头都带跑前／跑后的 `*.test` 家族读数；本机整段时间里有**别腿的 `wisp.test` 并发**（具名例：pid 22704＝`…\wisp-test-gotmp\orch-298ac4-20261012-142231\…\b001\wisp.test.exe`，编排者自己的根）⇒ 归属靠命令行路径＋根名，⛔ 记到我头上（定式见欠账 `Q-295-r1-fleet`）。

## 3. `AC#5`〔可选，排最后〕`-race` 代价表（⛔ 改 CI 一字，判"要不要请进门禁"归编排者＋机主）

现跑那一发（逐字）＝`GOTMPDIR="<专用根>" PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -race ./cmd/wisp/ -count=1 -timeout 8m`，件 `logs/60-ac5-race.txt`（360,724 字节）。

| 项 | 读数 |
|---|---|
| 墙钟 | 包级逐字 `FAIL github.com/CarlosShao/wisp/cmd/wisp 457.462s` ⇄ 同树非 race 两发 `377.287s`／`376.885s` ⇒ **＋约 80s（×1.21）**（尺＝逐字 `ok|FAIL` 行里的秒数，⛔ 算式冒充读数：80.2＝457.462−377.287） |
| 颜色 | **红**，`RACE_RC=1`；红名册 8 枚 |
| 恒红成因 | ① `WARNING: DATA RACE` **23 块**，栈顶全在 `bytes.(*Buffer).Write()` ⇄ `cmd/wisp.(*agentRuntime).auditf()` `run.go:933`（`config_reload.go:144` 的 tick 协程 ⇄ `internal/tools/bridge.go` 的 Bridge 写同一枚 buffer）；**`grep -c 'resident_audio' = 0`** ⇒ 与本票那枚 `started` ⛔ 相关；`race detected during execution of test` **4 枚**＝`TestRunPacketMarksTheRosterRowACardIsHolding`／`TestRunPacketCarriesTheSubagentItsRosterRowFed`／`TestRunPacketReportsTheStreamLogPastItsBound`／`Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`，这四枚在 D/E 非 race 两发里**全 PASS**（逐名尺）＝只因探测器而红；② 真窗族 3 枚（上面第 2 节）＋审批卡片那枚 flake ⇒ **今天请进门禁＝恒红**，且红因⛔ 在我这一格 |
| 反面一句（⛔ 结论、只报素材） | 本票那枚 `started` 的形状**在 race 探测器前不需要新用例来证明**（票面现量：`-race` 不在门禁名册 ⇒ 今天没有仪器看得见它）；这发唯一的用处是让"加它"这一步的代价可见。判归编排者＋机主 |
| 台件作废具名 | 15:37 那一发 `nohup … &`／后台任务那一形**子进程随工具 shell 退出被切**（`logs/60` 头 700 字节、跑后家族尺 0 枚、无 `go.exe`），同形还有 14:33 那发（`logs/12` 头两段作废说明）⇒ 两发都标作废、改前台重跑，⛔ 把"我没看到"读成"跑过了" |
