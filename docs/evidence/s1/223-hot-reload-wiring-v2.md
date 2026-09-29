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
- **逐名绿册（＝本腿交付的另一半，`223-v1` 欠的那半）**：1,161 行顶层绿名全量落盘＝`.scratch/wisp/probes/223/v2/v2-green-roster-top.txt`，573 行子例绿名＝`v2-green-roster-sub.txt`，两发合起来的逐枚结果行（含状态）＝`v2-all-result-lines.txt`。⚠ **1,161 枚名字不抄进本表**（那会把表撑成不可读）；名册以盘上原件为准，本表只抄**本票射程内的 12 枚顶层用例＋12 枚子例**（下表；分母尺＝`grep -rn "^func Test" --include=*_223*_test.go cmd/ internal/` 现跑＝**12 枚**（`config_reload_223_test.go` 8＋`config_reload_perm_223_windows_test.go` 1＋`config_sentences_223r2_test.go` 1＋`internal/config/manager_223_test.go` 2），逐名，`full-v.txt` 行号可回溯）：

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

## ② 攻 `223-r2` 那发正控与新轮询助手的牙 —— **正控有牙（本腿独立重做红）；`awaitStdout` 本身不吞失败，但它上面那半句拼接断言会吞一种（本腿实测出来）**

本腿一共做了**四发突变**（三发改产码＋一发 overlay 只读台件），全部逐字还原。**⛔ 本腿没有为了让任何一发变绿去放宽一条断言，也没有改一个字的产品文案留下不还原。**

### 2.1 正控（派单要的那一发，本腿自己重做，不转述 `223-r2` 的 1.4）

- 突变形状＝改**产码**：`cmd/wisp/config_reload.go` 的 `reportRestartPending` 里那句 stdout `Fprintf`（`:288-293`）整段注释掉；`:282`／`:284` 两行 stderr 审计**保留**。
- 尺（逐字）：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` → `go test ./cmd/wisp -count=1 -v -timeout 5m -run 'TestTicket223RestartTierSaysItWillNotApply'`。原始件 `.scratch/wisp/probes/223/v2/positive-control-v2.txt`（起手 **17:09:14**／终态 **17:10:00**，`EXIT=1`）。
- **读数＝红，`--- FAIL: TestTicket223RestartTierSaysItWillNotApply (42.18s)`**，红句逐字（原始件第 18 行）：

  > `config_reload_223_test.go:481: stdout never carried "本次运行不会生效" within 40s; full stdout:`

  贴出的 stdout 只有"答复监听已接入……"＋"配置热加载已接管……"两行；**同发 stderr 里三行审计全在**（第 25–27 行：`state=restart-pending sections=[app] effect=next-process-start tier=restart`／`RESTART-PENDING detail="…重启进程后生效" keys=[…]`／`state=applied hot=[] reload=[] restart=[app] locked=0`）。
  ⇒ **轮询非恒真**：句子拿掉必红，与 `223-r2` 自述的 42.18s／同一红句／同一行号**逐字对上**。⛔ 不是"20 连发全绿＝装饰"那种情况。
- 还原与哈希（`git cat-file blob <blob 14197622…> > cmd/wisp/config_reload.go`，⛔ 未用 `checkout`／`restore`／`reset`／`stash`／`clean`）：

| 时刻 | `certutil -hashfile cmd/wisp/config_reload.go SHA256` | 原件 |
|---|---|---|
| 突变 A 前基线 | `542f705ab7e3c7501ecd479c4096ffdda8a2192858ee7ce3557248a0d34f4670` | `hash-config_reload-baseline.txt` |
| 突变 A 后还原 | `542f705ab7e3c7501ecd479c4096ffdda8a2192858ee7ce3557248a0d34f4670` | `hash-config_reload-afterA.txt` |
| 突变 B 后还原 | `542f705ab7e3c7501ecd479c4096ffdda8a2192858ee7ce3557248a0d34f4670` | `hash-config_reload-afterB.txt` |
| 突变 D 后还原（收工态） | `542f705ab7e3c7501ecd479c4096ffdda8a2192858ee7ce3557248a0d34f4670` | `hash-config_reload-final.txt` |

  ⚠ 这枚基线哈希与 `223-r2` 表里记的那枚**逐字相同**（`542f705a…34f4670`）⇒ 本腿起手的工作区＝r2 交付态，四枚哈希之间零漂移；`git status --porcelain -- cmd/wisp internal/config` 收工＝**空**。

### 2.2 攻 `awaitStdout` 四问（逐问答，带 `file:line`；尺＝读码＋上面两发＋下面突变 B，不靠推断）

**问① 它只等"某句出现"，那句子说错（比如把重启档说成立即档）时它会不会照样绿？**
- **分两种"说错"，实测结论不同：**
  - **方向说反**（重启档被讲成立即档）⇒ **会红**，两道闸：`awaitStdout` 等不到 needle `本次运行不会生效`（`config_reload_223_test.go:481`→`:148-154`），且 `:495-497` 那枚负断言逐字钉着 `这些段已立即生效：[app]`。⇒ 这一形不靠轮询也能红，**没有洞**。
  - **说瘦**（仍含那五个字，但把"为什么／涉及哪几枚键"整段删掉）⇒ **实测照样绿**。突变 B＝把 `config_reload.go:288-293` 那句改成只剩 `"wisp run: 这些段的改动本次运行不会生效。"`（原因段、`涉及：%s` 段、"两件事都没发生"段全删），跑同一枚用例 ⇒ **`--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.24s)`、`ok cmd/wisp 2.272s`**（原始件 `.scratch/wisp/probes/223/v2/attack-b-stdout-content.txt`，17:11:00→17:11:06）。
  - **原因现读**：`:482-486` 那三枚内容 needle（`app.autostart`／`开机自启`／`重启进程后生效`）是在 **`why+out` 拼接串**上找的，而这三枚在 **stderr 那行审计里逐字全有**（产码 `config_reload.go:284-287` 的 `RESTART-PENDING detail="…开机自启注册、单实例锁、界面语言…重启进程后生效" keys=["app.language" "app.autostart" …]`）。⇒ 拼接把"stdout 那句话的内容完整性"整条让给了 stderr。
  - ⚠ **归责要说准**：这个拼接形状是 `223-r1` 交付时就有的（r2 表 1.2 第 3 条逐字写着"`:452` 的断言……全部逐字保留"，本腿现读复认它在 `:482-486`），**不是 r2 改出来的**；r2 改的只是 `out` 的取值方式。⇒ **产品事实面没问题**（现读 `:288-293` 那句确实带原因与涉及键），**这是验收用例的洞**：AC#7 的"用户能观测到**为什么**"在测试面上只被钉成"有那句话"，没被钉成"那句话里有为什么"。
- 还有一枚**前向**风险（本腿造不出当下反例，具名交给下一腿）：`awaitStdout` 是对**整条流**做 `strings.Contains`，不限"种下之后"。现在产品 stdout 只有那一处会打出这五个字（现跑 `grep -rn "本次运行不会生效" --include=*.go .` ⇒ 产码唯一命中 `config_reload.go:289`），所以今天没有假绿通道；**但如果哪天启动横幅（`config_reload.go:124-127`，实测原文见 2.1 原始件第 20 行，它现在写的是"重启档的改动**本次不生效**"，与 needle 差两个字）改成含 needle 的字样，用例会在 0ms 绿。**⇒ 下一腿该试的＝把断言窗口化（种之前先取 `out` 快照、要求 needle 出现在增量里），这与同文件 `plant`（`:71-94`）已有的"证明种下去了"纪律同形。

**问② 那枚 40s 期限会不会把真挂死伪装成"超时后才绿"，或反过来把真缺陷吞成 skip？**
- **不会伪装成绿**：期限到＝`t.Fatalf`（`:152-154`），只有 needle 真出现才 `return`（`:148-150`）；实测＝突变 A 那发 42.18s 红（不是 40s 过）。
- **不会吞成 skip**：helper 里没有 `t.Skip`，本票 12 枚顶层用例在 ①节那发整包里 **零 SKIP**（全发 7 枚 SKIP 全是历史在册的 wasapi／真下载／子进程助手，逐名见 `v2-redskip-roster-top.txt`）。
- **残余风险是反方向的，具名说清**：`awaitStdout` 只测"有没有"、**不测"多久"**——产品若要 8 秒才把这句打出来，用例照样绿。票面 AC#7 的字面判据没要求时延，所以本腿**不据此判红、也不加码**；要钉时延是另一格的事。
- 真死锁的形状反而比旧写法更早暴露：tick 侧死锁 ⇒ 40s 时在 helper 里红并 dump 双流，而不是等 `-timeout 30m`。

**问③ 它读的流是不是只有 stdout，stderr 那两行审计会不会被当成"已确认"？**
- **helper 自己只读 stdout**（`:147` `got := r.h.out.String()`），失败时才把 stderr 一并 dump（`:153-154`）⇒ **它没有把审计行当确认**。
- 但**用例整体**确实在 `:473` 先 `why := r.awaitAudit(t, "config: RESTART-PENDING detail=")` 把 stderr 那一行认成了"已到位"，再在 `:482-486` 用 `why+out` 找内容 needle ⇒ 见问① 的突变 B 实测。⇒ 精确说法：**`awaitStdout` 无辜，拼接断言有洞**；两者要分开记。

**问④ 期限常量与 `reloadCaseBudget` 是不是和文件里既有的 `awaitAudit` 同构？**
- **同构**（逐字对照 `awaitAudit :109-127` vs `awaitStdout :140-158`）：同一枚常量 `reloadCaseBudget = 40 * time.Second`（`:35`）、同 `time.NewTimer(reloadCaseBudget)`＋`time.NewTicker(20 * time.Millisecond)`、同 `select { case <-deadline.C: t.Fatalf(...) ; case <-tick.C: }`、同 `strings.Contains` 判定、同"返回整条流"的约定。同族的 `awaitCard :161-180`、`awaitLive :193-…` 也是这一形状。
- **差别只有两处，说出在哪**：① 读的流不同（`out` vs `err`）；② 失败消息多 dump 一条 stderr（`awaitAudit` 只 dump stderr）。⛔ 没有第三处差别，尤其**没有**"超时即算过"这种形状。
- 合规复认：ban #4 的尺射程＝`tools/d22scan` 的 `wallclockRe`（`\.Sub\(time\.Now\(\)\)`），helper 用 timer/ticker ⇒ 不沾；本腿收工再跑一次 `sh scripts/d22scan.sh` 自证（见时刻表，读数 **clean**）。

### 2.3 AC#4 那一发"拿掉实现必红"（本腿独立重做，逐名核对实现腿自述）

- 突变形状＝**整枚换回改前形状**：`git cat-file blob c774f8da:internal/config/loader.go > internal/config/loader.go`（`c774f8da`＝`223-r2` 自己的起手锚点，即 r2 新支不存在的那一发）。
- 尺：`go test ./cmd/wisp -count=1 -v -timeout 5m -run 'TestTicket223R2FailureSentenceRouting'`。原始件 `.scratch/wisp/probes/223/v2/teeth-loader-removed-v2.txt`（起手 **17:12:00**，`EXIT=1`）。
- **本腿这一发的红名册（逐名，8 枚子例全在原始件里）**：

| 子例名 | 本腿读数 | 实现腿 `223-r2` 自述 | 是否一致 |
|---|---|---|---|
| `声明当前版_注释以方括号开头_C1` | **FAIL** 0.07s | FAIL | 一致 |
| `声明当前版_CRLF_E1` | **FAIL** 0.07s | FAIL | 一致 |
| `声明当前版_无空格_G1` | **FAIL** 0.06s | FAIL | 一致 |
| `声明当前版_缩进版本行_L1` | **FAIL** 0.06s | FAIL | 一致 |
| `声明未来版_正文语法坏_J1` | **FAIL** 0.07s | FAIL | 一致 |
| `读不出版本_A1_基线不动` | PASS 0.07s | PASS（对照） | 一致 |
| `声明旧版_坏表头_A2_仍归迁移` | PASS 0.07s | PASS（对照） | 一致 |
| `解析得开_未知键_不抢语法错` | PASS 0.07s | PASS（对照） | 一致 |

  ⇒ **恰红 5 形＝C1／E1／G1／L1／J1，3 枚对照形纹丝不动**，与实现腿表 2.4 节自述**逐名全等**（本腿**没有**推翻它这一格；它的自述不多报也不少报）。
- 红句原文（5 枚同形，此处抄两枚，原始件第 30／33 行起）：

  > `config_sentences_223r2_test.go:111: shape 声明当前版_注释以方括号开头_C1 is booked "cause=invalid detail=\"config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置\"", want it to open with "cause=syntax" - the two sentences swapped again`

  ⇒ 用例把那句会说反的原文抄回来了＝**它测的是归句，不是判定**（见 ④ 节"新增 124 行用例文件有没有偷偷改判定"一格）。
- 还原与哈希（`git cat-file blob 42220c9c4dc4c090b3749ae2bc471d3be86bdfef > internal/config/loader.go`）：
  - 换回前基线（＝r2 交付态）`14251d9fa271e4eb93c7d58961fa49a494480b3fb71f2904242f9dabb28b0d30`
  - 换回后（c774f8da 形状）`f9c93512776ed94bde9068e29dce082b42d800bbbee3ee8b3865df7d85c624ca`
  - **再还原后** `14251d9fa271e4eb93c7d58961fa49a494480b3fb71f2904242f9dabb28b0d30` ⇒ **与基线逐字相同**（也与 `223-r2` 表里记的那枚相同），原件 `hash-loader-restored-v2.txt`。
- ⚠ **一发读数口径澄清**：`223-r2` 表 2.4 写"红的正是被修法接管的那 5 形"——本腿逐名复认，但**注意 J1 也在红列**，而 J1 归句是**编排者 09-29 16:5x 才裁的**（票面"我裁一句"）⇒ 这枚用例把"J1 归语法错"这件事**钉成了断言**。本腿 ③ 节对这一裁定的实测＝**支持**（见 3.1），所以本腿不判它越权；只把"用例钉住了待人裁点"这一形状登记在册，供裁点将来若被推翻时知道要动哪一行（`config_sentences_223r2_test.go:67-69`）。

### 2.4 额外两笔（本腿自己加的，不为难实现腿、只补 `223-v1` 具名没复现的那两发）

- **复现 m2（摘掉生产钩子）＝`223-v1` 推翻清单第 7 条里具名"本腿未复跑"的那一发，本腿跑了。** 突变＝注释掉 `cmd/wisp/config_reload.go:114` `rt.mgr.ConfirmLocked = rt.confirmLockedLoosening`（`:115` 与 `:119` 保留）。尺＝`go test ./cmd/wisp -count=1 -v -timeout 8m -run 'TestTicket223HandEditedFsLooseningCostsAnL2Card|TestTicket223ModeLooseningChangesTheRunningModeAfterAllow'`，原始件 `mutation-d-no-confirm-hook.txt`（17:14→17:15:5x）：
  **两枚全红**＝`--- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card (41.32s)`（红句 `config_reload_223_test.go:295: no "config.reload" card was displayed within 40s`）＋`--- FAIL: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (41.26s)`（红句 `:442` 同形）。
  ⇒ **AC#3 的正控在这条生产赋值路径上是有牙的**，而且本腿自己读到了"摘钩子之后不是静默变宽、是压根不弹卡"（fail-closed 现读 `internal/config/manager.go:172` 一带：`confirm := m.ConfirmLocked` 为 nil ⇒ 裁决不成 ⇒ 保留旧值）。⛔ 这一发**没有**触碰 `[fs]` 的任何真实放宽生效，最坏后果形状见本节末的三行定性。
  还原＝同一枚 blob `14197622…`，哈希＝`542f705a…34f4670`（上表第四行），`git status --porcelain -- cmd/wisp internal/config` **空**。
- **`cmd/wisp/config_reload.go` 那 11/6 到底删没删行为（派单点名要核的一格）**：`git diff c774f8da..HEAD -- cmd/wisp/config_reload.go` 现读＝**整段都在 `describeReloadFailure` 的迁移支注释里**（`-6` 行全是 `//` 开头的注释行，`+11` 行也全是 `//`），**函数体那一行 `return "cause=migration detail=\"" + …` 一字未动**。⇒ **"纯注释归真、零行为改动"这句为真**，本腿另用 ① 节整包（`ok cmd/wisp 135.353s`）与突变 D／A 两发反向钉住它没顺手改行为。

**涉及"权限／放宽／安全"字样的三行定性（本节 2.1／2.2／2.4 三发突变通用）**：
① **现象出现在哪**＝本机一个测试进程里，本腿亲手把产码的两行（产品句子／钩子赋值）临时注释掉之后、用例的读写时序与颜色；工作区在每一发之后都按 `git cat-file` 还原、哈希逐字相同，四枚文件收工零改动。
② **有没有本机被入侵的证据**＝**没有**。零外部输入、零网络、零权限变化、零真实 `config.toml` 被改（所有种文件都在 `t.TempDir()` 里）；`thresholds.go`／golden／`allowlist.txt`／三枚冻结件全程零接触（本腿连打开都没有打开过）。
③ **最坏后果是什么形状**＝突变 A／B 的最坏是"一句中文提示说瘦了而用例不报"（**文案与测试写法层面**，行为面仍是 fail-kept）；突变 D 的最坏是"放宽被静默拒绝、不弹卡"（**变严不变宽**，方向安全）；**没有一发把权限放宽过、也没有一发绕过过一张卡**。⛔ 本节不含任何安全事件定性。

## ③ J1 那一句裁定的实测（三形＋三形对照）—— **裁定前半被坐实；后半（"解析得开时版本更高那句才有资格说"）被本腿推翻：那一形今天根本不出口"版本更高"**

- 台件＝**`go test -overlay`**（⛔ 零源码改动、零产码迁就）：物理件 `.scratch/wisp/probes/223/v2/overlay/zz_v2probe_cmdwisp_test.go`＋`overlay.json`（虚拟路径 `cmd/wisp/zz_v2probe_cmdwisp_test.go`，`cmd/wisp` 目录里不存在这枚文件）。尺＝`go test -overlay .scratch/wisp/probes/223/v2/overlay/overlay.json ./cmd/wisp -count=1 -v -timeout 5m -run TestV2ProbeJ1Shapes`（先 export PATH），原始件 `.scratch/wisp/probes/223/v2/j1-three-shapes.txt`（起手 **17:13:31**／终态 **17:13:35**，`EXIT=0`）。
  ⚠ **口径说白**：这枚台件**只 `t.Logf` 读数、零判等**，它的 `PASS` 不代表实现绿，只代表"句子打出来了"。种子＝`config.SaveFile(path, config.NewDefaults())` → `config.NewManager` → 整枚换成对抗体 → **`mgr.CheckAndReload()`** → 错误原样交**生产分类器** `describeReloadFailure`（与 `reloadOnce` 逐字同路，`config_reload.go:153→155`）。文件都在 `t.TempDir()` 里，真实 `%APPDATA%\wisp` 零接触。
- **三形实测原文（逐字抄自原始件，⛔ 不是本腿改写）**：

**(a) 声明未来版＋正文语法坏**　body `"schema_version = 99\nbroken [[[\n"`
```
reload err    = config: config.toml parse: toml: expected character =
PRODUCTION LINE = cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"
```
**(b) 声明未来版＋正文解析得开（内含一份本程序不认的 `[app]` 键）**　body `"schema_version = 99\n\n[app]\nbogus_future_only_key = 1\n"`
```
reload err    = config: config.toml: schema_version 99 was written by a newer build (this build understands 2); upgrade Wisp or restore a backup
PRODUCTION LINE = cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"
```
**(c) 声明未来版＋顶层未知键**　body `"schema_version = 99\nthis_key_does_not_exist = 1\n"`
```
reload err    = config: config.toml: schema_version 99 was written by a newer build (this build understands 2); upgrade Wisp or restore a backup
PRODUCTION LINE = cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"
```
**三形对照（只为把 (b)(c) 的归因钉死，不是新增判据）**：
```
(d) 当前版＋顶层未知键 "schema_version = 2\nthis_key_does_not_exist = 1\n"
    err = config.toml: unknown key "this_key_does_not_exist" at line 2
    LINE = cause=unknown-key detail="config.toml 语法没问题，但里面有这份 schema 不认的键（拼错的键会被这样拒绝，而不是被忽略）。本次运行继续用内存里的旧配置"
(e) 当前版＋正文解析得开 "schema_version = 2\n"      -> err = <nil>，report = hot=[] reload=[] restart=[] locked=0（干净重载）
(f) 未来版＋全部键都合法 "schema_version = 99\n\n[ball]\nsize = 64\n"
    LINE = cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）…"（与 (b)(c) 同句）
```

### 3.1 逐条裁

- **(a) ⇒ 裁定前半坐实**。"版本比本程序还新**且正文解析不开**"这一形，第一句确实是**"读到了但解析不了"**，并且由常驻用例 `config_sentences_223r2_test.go:67-69`（J1 那一行，`wantCause="cause=syntax"`／`notCause="cause=invalid"`／`notFragment="语法没问题"`）**双向钉住**；本腿 ②2.3 那发把它打回改前形状时 J1 恰好转红 ⇒ 这颗钉有牙。⇒ **编排者这一条不用作废，r2 的现读裁与派单条形 1 一致、也被实测支持。**
- **(b) ⇒ 裁定后半被推翻（本腿明确写清正确说法）**。裁定原句（票面文末）＝「**"版本更高"那句只在"解析得开＋版本确实更高"时才有资格说**」。实测：解析得开的未来版文件（(b) 与 (f) 两形都试了，后者所有键都合法）**并没有说出"版本更高"**——它被 `describeReloadFailure` 归进 `cause=invalid`，那句中文逐字是「config.toml **语法没问题，但内容被校验拒绝（值不合法或引用解不开）**」。而现读 `internal/config/loader.go:106-110`：**"was written by a newer build (this build understands 2); upgrade Wisp or restore a backup" 这句是存在的、就在错误对象里**，只是 `cmd/wisp/config_reload.go:357-360` 那枚 `HasPrefix(d, "config.toml:")` 分支把它**吞进了"内容被校验拒绝"那一句**（匹配顺序：`config.toml parse` → `unknown key` → `cannot migrate` → `config.toml:`）。
  ⇒ **正确说法（替换裁定后半）**：*'版本更高'那句今天在**操作员面上零出口**——它只在配置层的错误文本里活着，生产分类器没有它自己的 `cause=`，所以任何"声明了更高 schema_version"的文件都被说成"内容被校验拒绝（值不合法或引用解不开）"，而那半句在这一形是**假的**（内容根本没走到 `decodeStrict`/`validate`，是版本闸拒的；对照 (e)＝解析得开的当前版才会真被校验）。操作员能执行的那句"升级 Wisp 或恢复备份"也一并丢了——讽刺的是**同一句中文在 `cause=migration` 那一支里是写全了的**（`config_reload.go:354-356`，但它只在 `ver < current` 时可达）。
  ⇒ 影响范围：(c) 同理——未来版＋未知键被版本闸抢在 `unknown-key` 那句之前，归了 `invalid`；对照 (d) 证明未知键那句本身是活的，所以这是**归句顺序**问题、不是分类器失灵。
- **归责要说准（⛔ 不许算到 r2 头上，本腿专门为这一条跑了对照发）**：`.scratch/wisp/probes/223/v2/j1-shapes-pre-r2-loader.txt`＝把 `internal/config/loader.go` 换回 **c774f8da**（r2 之前那一发）的形状后、**同一条 overlay 台件**的读数（17:20:14→17:20:18，`EXIT=0`；换回后哈希 `f9c93512…c624ca`、还原后 `14251d9f…b0d30` 与基线逐字相同）：
  - **(b)(c)(f) 三形改前改后逐字同一句**＝`cause=invalid`「config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）」⇒ **"版本更高那句零出口"这一形不是 r2 改出来的回归**，那一支 `ver > SchemaVersionCurrent`（`loader.go:106-110`）与那枚 `HasPrefix("config.toml:")` 分类分支（`config_reload.go:357-360`）r2 一字未动（② 节 2.4 的 11/6 全在注释里）。
  - **只有 (a) 变了**：改前 `cause=invalid`（那句"语法没问题"对着一个语法坏的文件＝假）→ 改后 `cause=syntax`（③ 节 3.1 的 (a) 原文）⇒ **r2 修的就是这一形，且实测确实换了句子。**
- **本格处置＝具名另立，不在本票翻勾里处理**：票面 AC#4 的字面判据只要求**四件各一句**（缺失／语法错／权限不够／热加载被禁用）——那四句本腿在 ① 节逐名读到、各有独立出口，所以 AC#4 仍**成立**；(b)(c) 这形是**第五件**，票面没写、派单没要求。⇒ 本腿**不判 AC#4 红**，只把这一形作为"AC#4 家族的新注＋建议另立票（或并入票 227 的归句矩阵）"登记，改法归写腿：给 `describeReloadFailure` 一枚自己的 `cause=newer-version` 分支（匹配串现成＝`was written by a newer build`），并在 `config_sentences_223r2_test.go` 的表里加一行；⛔ 本腿不动任何产码。

**三行定性（涉及"会说反／操作员被误导"字样）**：①现象出现在本机 `wisp run` 读一份"给更新版本写的 config.toml"时、写给操作员那一句中文的**归句选择**上，全部读数来自 `t.TempDir()` 里的对抗件；②没有本机被入侵的证据——零外部输入、零网络、零权限变化；③最坏后果形状＝**排查方向被指错**（该升级程序的人去查"值不合法"），当次重读仍然 fail-kept、内存续用旧配置。**不是权限放宽、不是绕卡、不是安全事件。**

## ④ AC#1..AC#7 逐格重裁（每格给依据行）

> 口径：以下"本腿现读"全部跑在锚点 `21552b44`（① 节）或其后的还原态（工作区与 HEAD 逐字节相同，见 2.1 的四枚哈希）。凡引 `223-v1` 的数字都标"二手"，本腿复跑过的标"本腿"。

| 格 | 本腿重裁 | 依据（行号＝本腿现读） |
|---|---|---|
| **AC#1 生产里真有人在轮询** | **成立**（注不变：这圈 tick 没被 join） | 生产调用者点数现跑＝**2**（`cmd/balldebug/main.go:243`＋**`cmd/wisp/config_reload.go:153`**）；链路 `cmd/wisp/run.go:620 rt.startConfigReload()` → `config_reload.go:105` → `:119 observe.Default.Spawn("watchdog","config",rt.reloadRoot,rt.configReloadTick)`（owner／recover 走 Registry，非裸 `go`）→ `:132 time.NewTicker(configReloadPollInterval)`，常量 `:88 = time.Second`；**不是墙钟差**＝重读与否由 `internal/config/manager.go:156` 的 `st.ModTime().Equal(m.seenMtime) && st.Size()==m.seenSize` 指纹相等决定，全仓 `manager.go` 里 `time.Now()/.Sub(` 零命中；收工 d22scan 本腿现跑 **clean**（`gates-v2.txt`）。注现读复认：`reloadHandle` 全仓**只有声明 `run.go:273` ＋赋值 `config_reload.go:119`、零读取点** ⇒ 未 join 那一笔仍在，归票 228（v1 说分母要从 1 变 2，本腿复认）。 |
| **AC#2 三档生效级别各有读数** | **成立带注**（立即档成立、重启档成立、**reload 档半格**） | 档位只有三枚（`internal/config/schema.go:33-44`＝`TierHot`／`TierReload`／`TierRestart`，⛔ 没人造第四档）；立即档有行为读数（本腿①发的 `TestTicket223RunArmsTheReloadTick` PASS 2.21s、`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` PASS 2.17s、`…TighteningRaisesNoCard` PASS 2.19s）；重启档有"明确告知"（见 AC#7）。**半格现跑**＝`grep -rn "OnReload = " --include=*.go cmd/ internal/ \| grep -v _test.go` 只有 `cmd/balldebug/main.go:244` ⇒ `internal/config/manager.go:198-199` 那一支在 `wisp run` 里恒不触发，`config_reload.go:168` 审计行的 `reload=` 列有值但**零消费者**。**该不该由本票闭合＝判不了，缺哪一句见下面"没做完"第 1 条**（⛔ 本腿不替编排者拍归属）。 |
| **AC#3 放宽必带 L2 复确认（D33 正控）** | **成立**（本腿有独立一发突变，＝v1 具名"没复跑"的那枚 m2） | 生产钩子赋值点唯一＝`config_reload.go:114`；本腿**把它注释掉后**跑两枚具名用例 ⇒ `--- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card (41.32s)`＋`--- FAIL: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (41.26s)`，红句逐字 `no "config.reload" card was displayed within 40s`（`:295`／`:442`；原件 `mutation-d-no-confirm-hook.txt`）⇒ **"拿掉复确认必须红"这一格本腿自己证到了，不再靠二手**。收紧不弹卡＝同格用例 PASS＋方向判定现读（`manager.go` 的 `plan` 只把 `len(loosen)>0` 送进 confirm 队列）；拒绝 fail-closed＝`…RefusedLooseningKeepsOldValues` PASS 5.34s。裁决位置＝`manager.go:172 confirm := m.ConfirmLocked`（锁外快照，见 AC#6）。 |
| **AC#4 不生效与读不到是两句话** | **成立带注**（注比 v1 那枚换了一处） | 四句各有出口且互不共用，本腿①发逐名绿：缺失／语法错／schema 拒绝（`TestTicket223FailureSentencesAreDistinct` 4 子例 PASS）、权限不够（`TestTicket223PermissionDeniedSitsInItsOwnSentence` PASS 2.29s，真 ACL）、热加载被禁用（`TestTicket223PanelInboundSaysHotReloadIsDisabled` PASS 0.01s，出口 `cmd/wisp/panel_inbound.go:211`）。**v1 抓到的"两句归反"已被 r2 修好且本腿逐名复认**（②2.3 那发：拿掉新支⇒恰红 C1/E1/G1/L1/J1）。**新的注＝③ 节实测出的第五形**：声明更高版＋解析得开（(b)(c)(f)）被归进 `cause=invalid`「内容被校验拒绝（值不合法或引用解不开）」，而它没被校验过、"升级 Wisp 或恢复备份"那句也被丢掉；这一形 v1 没测、票面四句没覆盖、r2 未碰（对照发证明改前改后同句）⇒ **本腿不判 AC#4 红，具名建议另立票／并入票 227 的归句矩阵**。 |
| **AC#5 整包终态读数** | **成立（本腿这一发＝票面要的那一发）** | 尺与逐名读数＝① 节：带 `-v`、全量 7,097 行未截断、起止 `date` 在案（17:01:29→17:03:53）；包级 **24 ok／2 FAIL**；红名册**按用例名集合**＝5 枚，与派单在册的 5 枚历史红**逐字相同、零新增**；本票 12 枚顶层用例逐名全绿，含 v1 那发的 `TestTicket223RestartTierSaysItWillNotApply`（本发 PASS 2.14s）；与 v1 那发逐名对拉＝只多 1 枚新用例名、消失 0 枚。⛔ 本票那枚红**没有在这一发里复发** ⇒ 派单那句"若仍红则 AC#5 直接不成立并退回 r2 的修法"**未触发**。**另有本腿自己的 10 连发复量**（原件 `flake-recheck-10.txt`，17:22:48→17:23:32）：**10/10 EXIT=0**，每发 2.017s～2.247s ⇒ 连同 r2 的 20 发与编排者的 20 发，同一把尺累计 **50 发零红**；再加 ②2.1 的正控（句子拿掉必红）⇒ "全绿不是恒真"这一条有读数。⚠ 统计口径本腿照写：累计连绿仍不等于"竞态已消失"，真正的堵口证据是**正控那一发的 42.18s 红**（helper 在等价缺陷下必响），这条本腿独立复现了。 |
| **AC#6 裁决在锁外** | **成立**（v1 已有专属突变读数，本腿复认形状、未重跑那发突变） | 现读 `internal/config/manager.go`：`confirm := m.ConfirmLocked`（`:172`，取快照）→ 解锁后才逐条裁决 → 回调 `:198-203`；`grep -rn applyLocked --include=*.go`（产码）本腿现跑＝**零命中**（旧形状已删）。专属牙＝`internal/config/manager_223_test.go` 的 `TestConfirmLockedRunsOutsideTheManagerLock`（本腿①发 PASS 0.03s；v1 二手读数＝挪回锁内时它 FAIL 6.02s、同发 `…HandEditedFsLooseningCostsAnL2Card` FAIL 41.12s）。⛔ 本腿**没有复跑锁内突变**（原因＝时间预算＋同一枚产码不得并发突变；不缺信息），具名列入"没做完"第 4 条。 |
| **AC#7 "重启后生效"要有出口** | **成立但带一条比 v1 更重的注** | 出口点数现跑＝1（`config_reload.go:115 rt.mgr.OnRestartPending = rt.reportRestartPending`），读取点 `manager.go:201-203`；内容三样齐（`:282-283` 审计带 `effect=next-process-start`、`:284-287` 带"为什么"与 `restartTierKeys`、`:288-293` 那句 stdout 带原因与涉及键）；不是"静默不生效"＝`planApp` 把 restart 键留在旧值并进 `rep.Restart`；独立用例＝`TestTicket223RestartTierSaysItWillNotApply`（没与 AC#2 合并，票面禁合满足）。⚠ **新注（本腿测出来的）**：**"为什么不生效"这半句没有被用例钉在 stdout 上**——②2.2 问① 的突变 B（那句被说瘦、只留 `本次运行不会生效`）实测 **PASS 2.24s**，因为 `:482-486` 的三枚内容 needle 是在 `why+out` 拼接串上找的，而 stderr 那行审计逐字全有。**产品事实对、验收用例弱**；补法（下一腿或写腿）＝把内容断言限定在 `out`，或窗口化到"种下之后的增量"。 |

### 4.1 派单点名要核的三笔（逐笔给读数）

1. **`cmd/wisp/config_reload.go` 那 11/6 是不是"纯注释归真、零行为改动"——是。** 尺＝`git diff c774f8da..HEAD -- cmd/wisp/config_reload.go` 现读（全文已抄进本腿的分析，输出 1 个 hunk）：删的 **6 行全部以 `//` 开头**（`describeReloadFailure` 迁移支那段 r1 写的边界自述），加的 **11 行也全部以 `//` 开头**；`return "cause=migration detail=\"" + …` 那三行产品文字与 `switch` 条件**一字未动**。⇒ 该文件的**行为面零改动**为真；本腿另有两发反向钉（突变 A／D 都在这个文件上，还原后整包 `cmd/wisp` 仍 `ok 135.353s`）。
2. **新增的 124 行用例文件是否只测归句、有没有偷偷改判定——只测归句，零改判定。** 尺＝通读 `cmd/wisp/config_sentences_223r2_test.go`（124 行，`git diff --numstat c774f8da..HEAD` 现读＝`124 0`，纯新增）：它做的事只有四件——种文件（`t.TempDir()`）→ `config.NewManager` → `mgr.CheckAndReload()` → 把错误交给**生产分类器** `describeReloadFailure`，然后断三件事（该开的 `cause=` 头、不许借的另一家 `cause=`、语法错支不许带「语法没问题」半句，`:110-121`）。⛔ **它没有**改任何判定阈值／权限／方向逻辑，**没有** assign 任何钩子，**没有**碰 `thresholds.go`／golden／`allowlist.txt`／三枚冻结件，**没有**跳过任何用例。归因同路证据现读：`reloadOnce` 的调用是 `config_reload.go:153→155`（`auditf("config: HOT-RELOAD state=not-applied %s", describeReloadFailure(err))`）⇒ 用例把同一个函数、同一枚错误对象拿在手里，只是同步跑、不走 tick（⚠ 这一条是**口径限制**不是缺陷：它证明不了"tick 那一刻也归对了句"，但 AC#4 那枚走真 tick 的 `TestTicket223FailureSentencesAreDistinct` 在本发同包内逐名绿，两枚互补）。
   ⚠ 唯一要登记的形状：`J1` 那一行（`:67-69`）把**编排者 09-29 16:5x 才裁的归句**钉成了常驻断言。本腿 ③ 节的实测**支持**这一裁定，所以不判它越权；但若将来裁点被推翻，要动的就是这三行＋`loader.go:90` 的半条件（r2 表"没做完"第 2 条早已具名写出改法）。
3. **AC#2 的 reload 档半格该不该由本票闭合——判不了，缺的是这一句定案（⛔ 本腿不拍）**：缺＝**「`TierReload` 的生产消费者（也就是 `Manager.OnReload` 的赋值点）归哪张票／哪个切片」**这一句归属定案。具体说，盘上今天同时成立的两件事实没有一句话把它们接上：
   - `internal/config/schema.go:36-38` 给 `TierReload` 的语义逐字是「applies immediately **AND emits a reload event so the owning subsystem can reload (e.g. voice model swap…)**」——它的定义**要求**一个 owning subsystem 去接那个事件；
   - 而接它的那枚钩子在 `wisp run` 里零赋值点（只有旁支 `cmd/balldebug/main.go:244` 一枚），`manager.go:198-199` 恒不触发。
   ⇒ 要人拍的原话形状＝"**给语音管线接 `OnReload` 这件事，是票 223 的射程内（那就本格未闭、223 不该 `-done`），还是另立票／归 S4·S5（那 223 按字面判据闭合、这笔账转出去）**"。票面七格的字面判据**没有一格**要求 reload 档读数（AC#2 只要"立即档变了"＋"重启档明确告知"），所以本腿既不能说它红，也不能替编排者说它不归 223。这正是 `AGENTS.md` §2"未定义即停"那一类；`223-v1` 的"缺归属定案、不是缺读数"这一句本腿**复认**，并把缺的那句写成了上面那行。

## 没做完／判不了（具名清单）

1. **`TierReload` 的消费端归哪张票——判不了，缺的是归属定案不是读数**（详见 ④4.1 第 3 笔那句要人拍的原话形状）。本腿不动票面、不自填。
2. **②2.2 问① 那枚"说瘦的句子照样绿"的洞，本腿只裁不修**（非实现者腿不产码）。补法具名＝把 `config_reload_223_test.go:482-486` 的内容 needle 从 `why+out` 改成只在 `out` 上找，或把 needle 断言窗口化到"种下之后的 stdout 增量"；⛔ **补它时必须同时保住 ②2.1 那发正控**（否则改出来的绿可能只是把 needle 挪到 stderr 上自己满足自己）。这一格属实现腿的活，本票若翻 `-done` 之前不补，就要在台账里具名记一笔"AC#7 的测试面只钉了有那句、没钉那句里有为什么"。
3. **③ 节那形（声明更高版＋解析得开 ⇒ 借 `cause=invalid` 那句、"升级 Wisp"零出口）本腿没写票**——改票面＝契约级动作。落点建议二选一由编排者拍：①并入票 227 的归句矩阵（它已经是"哪种坏文件说哪句话"的射程），②另立一票给 `describeReloadFailure` 加 `cause=newer-version` 分支＋路由表加一行。**本腿只交读数与两处行号（`loader.go:106-110`、`config_reload.go:357-360`）。**
4. **AC#6 的锁内突变本腿没复跑**（`223-v1` 那发的二手读数＝两枚具名红）。原因＝时间预算＋同一枚产码不得并发突变；**不缺任何信息**。若编排者要"本腿亲手证过 AC#6"，缺的就是这一发，尺与还原协议 v1 表里已写全，可直接照做。
5. **本腿没有重跑 20 连发**（那是 `223-r2` 与编排者各自的尺，两发都在盘上）。本腿只加了**自己那 10 连发**（`flake-recheck-10.txt`，10/10 EXIT=0）。累计 50 发零红的分母**来自三枚腿的三把尺**，不是本腿一发的读数——引用这句话时请带上这个口径。
6. **逐名绿册的"起手分母对拉"本腿做了对 `223-v1` 那一发的版本**（`make-rosters.sh`＋两份 name-roster，差集已抄进 ① 节），**没有**做对票 223 最起手（`223-r1` 之前）那一次的分母——那枚件 `probes/223/r1-start-green-roster.txt` 在盘上但**未提交**、且口径是"含子测试行 127 条"那把尺，与本腿的 `--- PASS` 名集合不同构。⇒ 若下一位要"从本票开跑前的基线逐名比"，缺的是**一枚与 `-v` 同构的起手名册件**，不是缺读数。
7. **`internal/panel` 4 枚＋`internal/ball` 1 枚历史红**：本腿未修、未读 `frontend/**`／`design/**` 一字（零读零引零转述）。d22scan 自己会报它扫了 `design/` 39 枚、`frontend/` 85 枚 text files——那是仪器的射程，本表只转引计数。
8. **`gofumpt` 不在本腿 PATH**（第一发 `GOFUMPT_EXIT=127`，见 `gates-v2.txt`）⇒ 本腿改用 `$(go env GOPATH)/bin/gofumpt.exe -l` 复检 r2 那四枚文件＝**零输出／EXIT=0**，已追加进同一枚原始件。这一条只是口径说明，不是缺陷。

## 门禁与尺：时刻表（本腿现跑，全部落盘）

| 时刻(+0800，`date` 现读) | 命令 | 落盘 | 结果 |
|---|---|---|---|
| 17:00:07 | 起手：`git rev-parse --short HEAD`＝`42911ce5`；`git status --porcelain -- cmd/wisp internal/config`＝空 | 本表起手节 | 锚点自取 |
| 17:00:4x | Write 骨架＋`git commit -F … -- <两枚 pathspec>` | 本表 | 骨架 `21552b44` |
| 17:01:24 | `ps -W \| grep -iE "go\.exe\|test\.exe"` | 本表 | **空**（无别人的门在飞） |
| **17:01:29→17:03:53** | **① 带 `-v` 整包**：`go test ./cmd/wisp ./internal/... -count=1 -v -timeout 30m`（先 export PATH） | `full-v.txt`（7,097 行／844,992 字节，未截断） | 24 ok／2 FAIL；红名册 5 枚历史红；`GATE_EXIT=1` |
| 17:05→17:06 | `make-rosters.sh`（逐名红册／绿册＋与 v1 名集合对拉） | `v2-green-roster-top.txt`／`v2-green-roster-sub.txt`／`v2-all-result-lines.txt`／`v2-redskip-roster-top.txt`／`v1-name-roster-top.txt`／`v2-name-roster-top.txt` | 只在 v2＝1 枚；消失＝0 枚 |
| 17:08:5x | `certutil -hashfile cmd/wisp/config_reload.go SHA256`（突变前基线） | `hash-config_reload-baseline.txt` | `542f705a…34f4670` |
| **17:09:14→17:10:00** | **② 正控**：产品 stdout `Fprintf` 整段注释掉 → `go test ./cmd/wisp -count=1 -v -run 'TestTicket223RestartTierSaysItWillNotApply'` | `positive-control-v2.txt` | **FAIL 42.18s**，红句 `:481 stdout never carried "本次运行不会生效" within 40s` |
| 17:10:3x | `git cat-file blob 14197622… > cmd/wisp/config_reload.go` ＋ `certutil` | `hash-config_reload-afterA.txt` | 与基线**逐字相同** |
| **17:11:00→17:11:06** | **② 攻牙突变 B**：stdout 那句只留 needle（原因／涉及／两件事都没发生 三段全删）→ 同一枚用例 | `attack-b-stdout-content.txt` | **PASS 2.24s**＝②2.2 问① 那枚洞的实测 |
| 17:11:3x | 还原＋`certutil` | `hash-config_reload-afterB.txt` | 与基线**逐字相同** |
| 17:11:5x | 现跑归因尺：`grep -rn "本次运行不会生效" --include=*.go .`／三枚钩子赋值点／`CheckAndReload()` 生产点数／`git diff c774f8da..HEAD -- cmd/wisp/config_reload.go` | 本表起手节＋④4.1 | 见对应节 |
| **17:12:00→17:12:1x** | **② AC#4 那发牙**：`git cat-file blob c774f8da:internal/config/loader.go > internal/config/loader.go` → `go test ./cmd/wisp -count=1 -v -run TestTicket223R2FailureSentenceRouting` | `teeth-loader-removed-v2.txt` | **恰红 C1/E1/G1/L1/J1、对照 3 形绿**，与实现腿自述逐名一致；`EXIT=1` |
| 17:12:3x | 还原 `git cat-file blob 42220c9c… > internal/config/loader.go` ＋ `certutil` | `hash-loader-restored-v2.txt` | `14251d9f…b0d30`＝基线**逐字相同** |
| 17:13:11→17:13:35 | **③ J1 三形**：`go test -overlay …/overlay.json ./cmd/wisp -count=1 -v -run TestV2ProbeJ1Shapes` | `j1-three-shapes.txt` | (a)`cause=syntax`／(b)(c)(f)`cause=invalid`／(d)`unknown-key`／(e) 干净重载 |
| **17:14:0x→17:15:5x** | **② 突变 D（复现 m2）**：注释掉 `config_reload.go:114` 钩子赋值 → 两枚 AC#3 用例 | `mutation-d-no-confirm-hook.txt` | **两枚全红 41.32s／41.26s**，红句 `no "config.reload" card was displayed within 40s` |
| 17:16:0x | 还原＋`certutil`（收工态） | `hash-config_reload-final.txt` | 与基线**逐字相同**；`git status --porcelain -- cmd/wisp internal/config` **空** |
| **17:20:14→17:20:18** | **③ 对照发**：loader 换回 `c774f8da` 形状后重跑同一枚 overlay 台件，再逐字还原 | `j1-shapes-pre-r2-loader.txt` | (b)(c)(f) 改前改后**同句**⇒不是 r2 回归；只有 (a) 换了句子 |
| 17:20:56→17:21:15 | 收工三把尺：`go build ./...`／`gofumpt -l`（四枚，见"没做完"第 8 条）／`sh scripts/d22scan.sh` | `gates-v2.txt` | BUILD_EXIT=0；GOFUMPT 零输出；**d22scan clean**（口径：bans #1-5 生产 248 枚，ban #8 internal 460＋cmd 63 枚含注释与 `_test.go`）；自检 `runtests.sh: OK PASS=34 FAIL=0 SKIP=0` |
| **17:22:48→17:23:32** | **① 补强**：本腿自己的 10 连发 `go test ./cmd/wisp -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply'` | `flake-recheck-10.txt` | **10/10 EXIT=0**（2.017～2.247s 逐发在案） |
| 收工 | `git status --porcelain -- cmd/wisp internal/config` | 本表 | **空**＝本腿零源码改动、零 commit 产码、未 push |

## 交回编排者的六节（大白话；面向不看技术的 owner）

**① 带 `-v` 的终态那一发（这是这张票欠了整天的那一发）**
17:01:29 起跑、17:03:53 收，7,097 行原始输出一个字节都没截（`full-v.txt`，844KB）。26 个包：**24 绿 2 红**。红的两包是 `internal/ball` 和 `internal/panel`，按**用例名**数出来正好 5 枚，全部是别人地界上早就在册的老红（面板 4 枚＋球 1 枚），**一枚新的都没有**。上一轮咬着这张票的那枚 `TestTicket223RestartTierSaysItWillNotApply`，这一发**绿了（2.14 秒）**；`internal/risk` 那枚计时用例这一发也绿了（不需要复量）。本票自己 12 枚顶层用例逐名全绿（表里都抄了），新加的那枚路由用例 8 个子例逐名绿。我还把这枚曾经间歇红的用例**自己连跑了 10 次，10 次全过**。跟上一轮验收腿那一发做"逐名对拉"：两次之间只多出一枚用例名（就是 r2 新加的那枚），**没有一枚用例从名册上消失**。**结论：AC#5 那格"不成立"的读数依据已经不存在了——这一格现在成立。**

**② 正控与那个新轮询助手有没有牙（本腿亲手打的，不是复述）**
把产品那句"本次运行不会生效"整段注释掉再跑：用例 **42.18 秒变红**，报的话逐字是"stdout 从来没出现过这句话（40 秒内）"，而同一次的后台账本里三行重启记录**都在**——**这就是"轮询不是装饰品"的硬证据**，上一枚腿量的 42.18 秒我一字不差复现了。改前改后的 SHA256 四枚**逐字相同**，工作区收工时零改动，我没提交任何源码。
新助手 `awaitStdout` 我逐行读了四个问题，答案是：**它自己不吞失败**（等不到就是红，没有"超时算过"，也没有 skip；形状与文件里既有的 `awaitAudit` 同构，只读 stdout 不读 stderr，用的不是墙钟时间差）。**但它上面那半句断言会吞一种**：我做了一发新突变——让产品照样说出"本次运行不会生效"这六个字、但把"为什么不生效、涉及哪几个设置"整段删掉，**用例仍然绿（2.24 秒）**。原因很具体：那三个内容检查是在"后台账本＋屏幕那行"拼起来的串上找的，而账本里那些字全有。⛔ 这不是 r2 改出来的（那半句拼接是第一版实现就有的），**产品说的话是对的、验收用例不够硬**；补法我写在"没做完"第 2 条，属写腿的活。

**③ J1 那形（"声明了更高的版本号，而且文件读不出内容"到底算哪句话）**
上一版实现把三形全测了一遍。**你的裁定前半被坐实**：文件既声明了未来的版本、正文又解析不开时，操作员听到的第一句确实是「**config.toml 读到了但解析不了：这一行不是合法 TOML 语法**」——你要的"先说读到了但解析不了"做到了，而且被常驻用例钉住（我把实现换回旧形状，它立刻红）。
**你的裁定后半被推翻**，以实测为准：你说"'版本更高'那句只在解析得开时才有资格说"——但实测**解析得开＋版本更高的文件（我试了三形：带未知键、带不认识的表、所有键都合法）根本不说"版本更高"**。它被归进了"config.toml **语法没问题，但内容被校验拒绝（值不合法或引用解不开）**"这一句——而那半句在这一形是**假的**（文件压根没走到校验，是版本闸把它挡在门外），而且配置层明明备好了那句"是更新的 Wisp 写的，请升级或恢复备份"，被分类器**丢掉了**。讽刺的是这句中文在"版本太旧"那一支里是写全了的，可那一支只在版本低于当前时才可达。**正确说法（供你在台账里作废原句、换成这句）**："'版本更高'那句今天在操作员面上**零出口**；要它出口，得给分类器加一枚自己的分支（`loader.go:106-110` 有现成的匹配串）。"**并且我把归责钉死了**：把实现换回 r2 之前的形状再跑一遍，这三形**改前改后逐字同一句**——所以这不是 r2 弄坏的新伤，是同一个"两句会说反"的老家族里 r2 没被要求处理的另一形。我没有替你把它判成 AC#4 红（票面四句话的四句各自有出口、都在绿册里），只是具名建议另立票或并入票 227。

**④ 这张票现在到底还差什么才能 `-done`**
按字面判据逐格重裁：**AC#1／AC#3／AC#5／AC#6 成立**，**AC#2／AC#4／AC#7 成立但各带一条注**。要翻 `-done`，我看还差的不是实现，是**四个决定**（都不该由我做）：
1. **AC#2 那半格归谁**——`TierReload` 的定义写着"要被那个子系统的拥有者接走"，可 `wisp run` 里没有一个人接。要不要 223 来接＝**需要你拍一句归属**（我写成了原话形状，见 ④4.1 第 3 笔）。这是七格里唯一一格**判不了**的，缺的是定案、不是读数。
2. **AC#7 那条新注要不要在翻勾前补**——现在的产品是对的，但用例钉不住"那句里到底有没有为什么"（我用突变实测出了绿）。要么现在补一发硬断言，要么在台账里具名记一笔"AC#7 的测试面只钉了有那句"。
3. **③ 那形要不要立项**——"给更新版本写的配置文件"这一形今天说一句假话（内容被校验拒绝），而它其实有现成的正确句子可说。
4. **AC#6 我沿用的是上一枚腿的突变读数**，本腿没有复跑锁内那一发（不缺信息，是时间预算＋同一枚文件不能两发同飞）。你要"本腿亲手证过 AC#6"，就差这一发，尺我照 v1 抄在时刻表下面了。
⛔ **我的判语不自动等于翻勾**，字面判据优先、翻勾按 `SPEC-12 §4.3` 由你做。另两笔不属于 223 但被它冲出来的账我**没有改任何票面**：226 那枚"没人轮询"的前提（本发它 4.82 秒跑满约 4 个 tick 仍绿，所以"靠时序侥幸"这个解释在这个条件下没复现，读数已交回、归票 227 裁）；228 那格"没 join 的协程"分母仍是 2 枚（`reloadHandle` 现跑零读取点）。

**⑤ 有没有为了变绿放宽任何断言／有没有碰到权限与安全**
**没有。**本腿全程零源码提交、四枚动过的文件每一发都按 `git cat-file` 逐字还原并核对哈希（五枚哈希两两逐字相同）、`git status --porcelain -- cmd/wisp internal/config` 收工＝空；`thresholds.go`／golden／`allowlist.txt`／三枚冻结件／`PLAN.md`／`docs/specs/**`／票面／`docs/reports/**` 全程零接触（工单票面是另一枚只读腿的地界，我一个字没写）；`frontend/**`／`design/**` 零读零引零转述。**没有一处在"权限／放宽"上被放宽过**：本腿三发产码突变的性质分别是"少说一句中文"、"把那句中文说瘦"、"摘掉弹卡的钩子"，最坏后果形状是**提示语不准**或**放宽被静默拒绝（变严，不是变宽）**，全程 fail-kept、没有任何一张卡被绕过、没有任何一条真实路径判定被改动。⛔ **本表不含任何安全事件定性，测试时序问题不被我写成安全事件。**

**⑥ 我没做完／判不了的（逐枚，别当成"都验过了"）**
`TierReload` 归属＝判不了，缺定案句；AC#6 的锁内突变＝本腿没复跑，缺那一发、不缺信息；从本票最起手基线做逐名绿册对拉＝缺一枚与本尺同构的起手名册件；"说瘦的句子照样绿"那枚洞＝我只裁不修，补法具名写了；"版本更高那句零出口"＝我没写票、只交行号；10 连发是本腿的加测，**50 发零红那个分母是三枚腿的三把尺加起来**，不是我一发的读数；`gofumpt` 不在 PATH 这件事我换了 `GOPATH/bin` 那枚复检并把两发读数都写进同一原始件（口径说明，不是缺陷）。
