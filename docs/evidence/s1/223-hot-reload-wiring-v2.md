# 票 223 非实现者对抗验收 `223-v2` —— AC#5 终态那一发 ＋ `223-r2` 修法的牙

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 被裁的实现件：`docs/evidence/s1/223-hot-reload-wiring-r2.md`（写腿 `223-r2`，收码 commit `4d925c2d` 一系）
- 上一枚验收腿：`docs/evidence/s1/223-hot-reload-wiring-v1.md`（判 AC#5 不成立、AC#4 两句归反）
- 更早的实现件：`docs/evidence/s1/223-hot-reload-wiring-r1.md`
- 本腿起手锚点：`git rev-parse --short HEAD` = **`42911ce5`**（读数全部现跑于骨架锚点 `21552b44`，两者之间零产码变更），分支 `dev`
- 本腿性质：**只裁、只跑台件、不产码、不 commit 任何源码改动**（突变＝一进一出，逐字还原＋SHA256 自证）
- 台件目录：`.scratch/wisp/probes/223/v2/`（原始输出全量落盘，绝不接 `| head`／`| tail`）
- ⛔ 零读零写零转述：`frontend/**`、`design/**`
- ⛔ 不碰：`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件、`docs/reports/**`、任何工单票面
- ⛔ 不碰别人的脏改动：`.gitignore`、`design/**`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`

## 入库清单（本腿会 commit 的路径，逐枚具名）

| 路径 | 是什么 |
|---|---|
| `docs/evidence/s1/223-hot-reload-wiring-v2.md` | 本表 |
| `.scratch/wisp/probes/223/v2/` | 全部原始读数与台件（含 `full-v.txt`、两发突变日志、J1 三形台件） |

---

## 起手读数

- 起手 `date` = **Tue Sep 29 17:00:07 CST 2026**；起手锚点 `git rev-parse --short HEAD` = **`42911ce5`**（⛔ 不是派单里写的任何号；派单没给号，本腿自取）。
- 骨架提交＝**`21552b44`**，本腿所有读数都跑在 `21552b44` 这枚锚点上（＝`42911ce5`＋编排者的 `cfc4df9a` 台账件＋本腿骨架件）。
- ⚠ **产码归因自证（现跑）**：`git diff --stat 42911ce5..HEAD -- cmd internal '*.go'` ＝ **空** ⇒ 起手到现在**没有任何产码变动**，那两枚中间 commit 只动 `docs/reports/**` 与本腿的表。⇒ 本腿的读数就是 `223-r2` 交付的码（`5a755c3c` 一系，收码 commit `4d925c2d`）。
- 在飞检查：`ps -W | grep -i -E "go\.exe|test\.exe|go-build"` 起手＝**空**（17:01:24 现跑）；`git status --porcelain -- cmd/wisp internal/config` 起手＝**空**；`git diff --cached --name-only` 起手＝**空**（索引干净）。
- 别人地界的脏改动本腿一律不碰（起手 `git status --porcelain` 现读到名字、未读内容）：`.gitignore`、`design/**`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`、`docs/reports/**`、`*-done.md` 票面（＝另一枚只读腿 `done-fix-1` 的地界，本腿零接触）。

| 尺（本腿现跑） | 读数 |
|---|---|
| `grep -rn "CheckAndReload()" --include=*.go cmd/ internal/ \| grep -v _test.go` | **2 枚生产调用者**＝`cmd/balldebug/main.go:243` ＋ `cmd/wisp/config_reload.go:153`；`internal/ball/hotkey_reload.go:21` 是注释；定义在 `internal/config/manager.go:142` |
| `grep -rn "ConfirmLocked = \|OnRestartPending = \|OnReload = " --include=*.go cmd/ internal/ \| grep -v _test.go` | 三枚＝`cmd/wisp/config_reload.go:114`（`ConfirmLocked`）／`:115`（`OnRestartPending`）／**`OnReload` 只有 `cmd/balldebug/main.go:244`** ⇒ **`wisp run` 的 reload 档今天仍无生产赋值点**（AC#2 半格复认） |
| `grep -rn "本次运行不会生效" --include=*.go .` | **全仓 3 处**＝产码 **1 处**（`cmd/wisp/config_reload.go:289` 那句 `Fprintf`）＋测试 2 处（`config_reload_223_test.go:481` 的 needle、`:490` 的既有断言）。⇒ `awaitStdout` 等的那枚 needle **在产品 stdout 上只有一个来源**，不存在"别的话撞上了这句"的通道（四问第 ① 问的依据，详见 ② 节） |

## ① AC#5 缺的那一发：带 `-v` 的终态整包读数 —— **成立（这一发的红名册里已经没有本票的用例）**

- 尺（逐字）：先 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，再 `go test ./cmd/wisp ./internal/... -count=1 -v -timeout 30m`；台件＝`.scratch/wisp/probes/223/v2/run-full-v.sh`（脚本原文入库，`date` 与锚点都写在原始件首尾）。
- **起手 `date`＝`Tue Sep 29 17:01:29 CST 2026`／终态 `date`＝`Tue Sep 29 17:03:53 CST 2026`**（跑前 17:01:24 现查 `ps -W`＝空，全程只有这一发门在飞）。
- 原始输出**全量**落盘＝`.scratch/wisp/probes/223/v2/full-v.txt`（**7,097 行／844,992 字节**，⛔ 未接 `| head`／`| tail`），`GATE_EXIT=1`（第 7095 行）。
- **包级读数（26 枚包＝24 `ok`／2 `FAIL`）**：两枚红＝`internal/ball 0.304s`、`internal/panel 4.327s`（**都是别人地界的历史红**）。⚠ **本票相关的两包都在绿侧**：`ok cmd/wisp 135.353s`、`ok internal/config 2.569s`；`internal/risk 9.095s` 也是 `ok`。
  逐包原文在 `full-v.txt` 内可数：`ok` 行 24 枚、`FAIL<tab>包` 行 2 枚（尺＝`grep -E "^(ok|FAIL)[[:space:]]" full-v.txt`）。
- **口径（票面第 5 行定死＝比"用例名集合"、不比枚数；本腿把两套都交）**：
  顶层结果行 `--- PASS` **1,161**、`--- FAIL` **5**、`--- SKIP` 顶层 **7**；缩进层 `--- PASS` **573**、缩进层 `--- FAIL` **0**；顶层 `=== RUN` **1,746**。
  ⚠ **与 `223-v1` 那发的同名集合对拉**（脚本 `.scratch/wisp/probes/223/v2/make-rosters.sh`，产物 `v1-name-roster-top.txt` 1,166 名／`v2-name-roster-top.txt` 1,167 名）：
  - **只在 v2 出现的顶层名＝1 枚**＝`TestTicket223R2FailureSentenceRouting`（＝r2 新增的那枚常驻路由用例，8 枚子例）。
  - **只在 v1 出现的顶层名＝0 枚**＝**没有任何用例从名册上消失**。
  - **状态翻动的名＝2 枚**：`TestTicket223RestartTierSaysItWillNotApply` **FAIL→PASS**、`TestResolvePerCallBudget` **FAIL→PASS**。
  - 两发的 `=== RUN` 差（1,737→1,746＝＋9）与名集合差**完全自洽**＝＋1 枚顶层名＋8 枚子例，⛔ 不是有测试在两次之间被删或被加。
- **逐名红册（5 枚，全部为在册历史红＝别人地界，本票不新增、本腿不顺手修、本腿不读 `frontend/**`／`design/**` 去解释它们）**：

| 用例名 | 包 |
|---|---|
| `TestC21TableColourRowsMatchTokensCSS` | internal/ball |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | internal/panel |
| `TestComposerContractTypesMatchFrontend` | internal/panel |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | internal/panel |
| `TestC21DesignTokensFourWayAgree` | internal/panel |

  ⇒ **与派单给的"在册 5 枚"逐字相同、零多零少**。⛔ **本票那一枚红（`TestTicket223RestartTierSaysItWillNotApply`）在这一发里没有再红**（`full-v.txt:206` ＝ `--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.14s)`）⇒ 按派单那句判据：**AC#5 不成立的那个点已消除**，且这是**第一发跑在"红已堵上之后"的带 `-v` 终态**（`223-r2` 表"没做完"第 1 条与票面"仍缺的那一发"要的就是它）。
- **逐名绿册（＝本腿交付的另一半，`223-v1` 欠的那半）**：1,161 行顶层绿名全量落盘＝`.scratch/wisp/probes/223/v2/v2-green-roster-top.txt`，573 行子例绿名＝`v2-green-roster-sub.txt`，两发合起来的逐枚结果行（含状态）＝`v2-all-result-lines.txt`。⚠ **1,161 枚名字不抄进本表**（那会把表撑成不可读）；名册以盘上原件为准，本表只抄**本票射程内的 13 枚顶层＋12 枚子例**（下表，逐名，`full-v.txt` 行号可回溯）：

| 本票用例（顶层／子例） | 这一发的读数 |
|---|---|
| `TestTicket223RunArmsTheReloadTick` | PASS 2.21s |
| `TestTicket223HandEditedFsLooseningCostsAnL2Card` | PASS 2.30s |
| `TestTicket223RefusedLooseningKeepsOldValues` | PASS 5.34s |
| `TestTicket223TighteningRaisesNoCard` | PASS 2.19s |
| `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` | PASS 2.17s |
| **`TestTicket223RestartTierSaysItWillNotApply`** | **PASS 2.14s**（＝v1 那发的 FAIL，本腿复量见下） |
| `TestTicket223FailureSentencesAreDistinct` | PASS 8.61s（4 子例全绿：缺失 2.12／语法错 2.25／声明了版本但坏在后面 2.08／schema未知键 2.16） |
| `TestTicket223PanelInboundSaysHotReloadIsDisabled` | PASS 0.01s |
| `TestTicket223PermissionDeniedSitsInItsOwnSentence` | PASS 2.29s（真 ACL 这台上跑通） |
| **`TestTicket223R2FailureSentenceRouting`** | **PASS 0.53s**（8 子例逐名绿：C1 0.07／E1 0.07／G1 0.06／L1 0.07／J1 0.07／A1 0.07／A2 0.07／解析得开未知键 0.07） |
| `TestConfirmLockedRunsOutsideTheManagerLock` | PASS 0.03s |
| `TestCheckAndReloadAsksOncePerFileChange` | PASS 0.03s |

- **`TestResolvePerCallBudget`：这一发是绿的**（`full-v.txt:3745` ＝ `--- PASS (1.47s)`，所在包 `ok internal/risk 9.095s`）⇒ 派单那条"若红＝争用型假红、安静后 `-count=3` 复量"的处置**没有触发**（没有红可复量）。⛔ `thresholds.go`／golden 本腿全程零接触（连文件都没打开过）。
- **顺带一枚跨票读数（登记不处置，归票 226／227 的账）**：`223-v1` 说过票 226 那枚唯一走真宿主的用例靠"整条 run 跑不满 1s tick"过关（它量到 1.11／0.91／0.93s）。**这一发它是 `--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (4.82s)`**（`full-v.txt:18`）＝在同一条 1s tick 下**跑满了约 4 发 tick 仍然绿**。⛔ 本腿不据此翻 v1 的结论（那是单包 3 连发 vs 整包并发两种条件），只把这一行读数交回：**"时序侥幸"这一解释在本条件下没有复现**，归票 227 的射程裁。
- **"拿掉实现必红"这一格本腿另有独立读数**：`cmd/wisp/config_reload_223_test.go` 全文不 assign 三枚钩子、也不自调 `CheckAndReload`（`223-v1` 已证、本腿现读复认同一形状），所以这些绿是**被观测到的通电**，不是字段在场。②节的两发突变（正控＋拿掉 loader 新支）是本腿自己重做的、更硬的那一版。

## ② 攻 `223-r2` 那发正控与新轮询助手的牙

〔待填：正控独立重做（原始红句＋前后 SHA256）、`awaitStdout` 四问逐答带 `file:line`、AC#4 拿掉实现那发的红形名册 vs 实现腿自述〕

## ③ J1 那一句裁定的实测（三形）

〔待填：(a) 未来版＋语法坏／(b) 未来版＋解析得开／(c) 未来版＋未知键 的实际句子原文，裁定坐实或推翻〕

## ④ AC#1..AC#7 逐格重裁

〔待填：七格各一句＋依据行；重点核三笔＝`config_reload.go` 11/6 纯注释归真是否零行为、新增 124 行用例文件是否只测归句、AC#2 reload 档半格归属〕

## 没做完／判不了（具名清单）

〔待填〕

## 门禁与尺：时刻表

〔待填〕

## 交回编排者的六节（大白话）

〔待填〕
