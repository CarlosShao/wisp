# ci223-1：`TestTicket223HandEditedFsLooseningCostsAnL2Card` 基线 PASS 到 CI 新发 FAIL 的归因（只读取证）

> 本腿零 `go test` / 零 `go build` / 零 `go vet` / 零 exe 执行（题面硬边界 #1，盘上有两枚写码腿在跑毫秒级计时）。
> 全部结论只来自：读文件 + 读盘上已有的 CI 日志 + `git log/diff/show`。
> §3 与 §5 里凡"尺"是 `go test` 的那几把，**本腿没有跑、也不许跑**，逐条标了〔未自跑〕；其余每一把我都自己跑过。

## 0.-1 引用锚（先读这一条，否则会拿到过期行号）

**本文所有 Go 源码行号一律钉在 `8ae4c23e`**（＝新发那发 CI 真正编译的 commit），钉法：`git show 8ae4c23e:<file> | grep -n <串>`。
原因：本腿起手锚是 `6fe300c5`（08:41:33 +08），写到一半 HEAD 漂到了 `0d87a681`（09:03:15 +08），`cmd/wisp/run.go` 被别的腿改过，**我引的那几行整体下移了 17 行**（例：卡的 `Fprintf` 在 `8ae4c23e` 是 `run.go:1332`，在 `0d87a681` 是 `1349`）。不钉锚就读行号，下一枚腿会读到错的地方。
两枚 CI 锚点上这张卡的打印语句各自在哪：`0589fd9c` = `cmd/wisp/run.go:1146`，`8ae4c23e` = `cmd/wisp/run.go:1332`。
CI 日志件（`clean_*` / `raw_*` / `jobs_*`）的行号是**盘上文件自己的行号**，与 commit 无关，可直接用。

---

## 0.-2 边界自报（一处技术性越界，具名，请编排者裁）

题面硬边界 #2 写的是 `frontend/**` 与 `design/**` **零读取、零引用、零写入**。
本腿自查：为了找那句渲染语句的出处，我跑过一发 `grep -rn "确认 L2" --include=*.go .`（全仓递归）。
- 后果实测：`find frontend -name '*.go'` 现量只有 **1 枚** `frontend/embed.go`，`grep --include=*.go .` 会去打开它做字节匹配 ⇒ **这一发在"零读取"上是技术性越界**，本腿不辩解。
- 材料性：`grep -rn "确认 L2" --include=*.go frontend` 现量 **0 命中**；那发的输出里**没有任何一行来自 `frontend/**`**，本腿全文也**没有引用**过它的任何内容。`design/**` 侧 `find design -name '*.go'` = **0 枚**，那发递归 grep 因此**一个 `design` 字节都没打开过**。
- 本腿改用过的正确形状（后续腿请照这个）：把搜索根限定在 `cmd internal tools`，而不是仓库根。
- 写入面：`git show --stat` 自证本腿那一枚 commit（见交付回执）**只含本腿那 9 枚路径**，没有 `frontend/**`、没有 `design/**`。
  ⚠ 本腿写面时盘上另有他人的未提交件（`M cmd/wisp/panel_config_248_test.go`、`D design/assets/base.css`、`M .gitignore`、`.scratch/wisp/probes/161/**` 等），本腿**一字未动、未 stage**。

---


## §0 起手锚与写面声明

### 0.1 锚（自己取的，没用题面任何 sha 当锚）
- `git log -1 --format='%H %ad' --date=iso` = `6fe300c5035a60586e681f586a4ead5e2b56137b 2026-10-02 08:41:33 +0800`
- `git rev-parse --abbrev-ref HEAD` = `dev`
- `date` = `Fri Oct  2 08:46:01 CST 2026`
- **写面过程中复量**（`09:03:25 +08`）：HEAD 已漂到 `0d87a681b050d2557f4ae37a25799bf28633b168`（`2026-10-02 09:03:15 +0800`）⇒ 见 0.-1。**这正是本仓"未推枚数会漂"那一坑在读码侧的形态：行号也会漂。**

### 0.2 `git status --porcelain cmd internal` 读数（取数时刻 `2026-10-02 08:46 +08`）
```
 M internal/panel/bridge.go
 M internal/panel/composer_dispatch.go
 M internal/panel/composer_dispatch_test.go
 M internal/panel/l2_grant_boundary_test.go
?? internal/ball/sta_release_windows_test.go
?? internal/config/settings.go
?? internal/panel/config_handlers.go
```
= 7 行。全是别人的未提交件（在飞腿 `internal/panel` / `internal/ball`），本腿一字未动。
`design/**` 本腿**零读取**（现量 `find design -name '*.go'` = 0，故那发递归 grep 一个 design 字节都没打开）；`frontend/**` 有 **1 枚 `frontend/embed.go`** 被一发递归 grep 打开做过字节匹配，**0 命中、零引用**——这一条记在 **0.-2**，不在这里算成"零读取"。
写完时复量 `git status --porcelain cmd` = **空**（0 行，`09:03:25 +08`），因为那 7 行里的 `cmd` 侧本来就没有脏件，而 `internal` 侧已被别的腿提交进去了——**这条也归"会漂"，别当常量读**。
本腿唯一写面＝`.scratch/wisp/probes/orchestrator/ci223-1/` 这一个目录。产出件（只建不删）：
- `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\orchestrator\ci223-1\attrib.md`（本文）
- `...\ci223-1\raw_range_roster.txt`（`git log --oneline 0589fd9c..8ae4c23e`，231 行）
- `...\ci223-1\raw_commits_touching_cmdwisp.txt`（逐枚 numstat 筛出的 29 枚）
- `...\ci223-1\cmdwisp_base_names.txt` / `cmdwisp_new_names.txt`（两发 cmd/wisp 步顶层名册，191 / 241）
- `...\ci223-1\before_ours_base.txt` / `before_ours_new.txt`（本枚之前那 13 枚的名册）
- `...\ci223-1\order_clean_base.txt.txt` / `order_clean_new.txt.txt`（整份日志的**首次出现序**顶层名册，1769 / 886 枚；文件名里那两遍 `.txt` 是本腿拼 `order_$f.txt` 时的手滑，**件是对的、名是丑的**，临时件只建不删所以就这么留着）

### 0.3 复认题面给的两条盘上原始读数（复认成功，但行号要纠）
题面"`clean_new.txt:2170` 起是新发日志、`:3022` 起（同目录 `clean_base.txt`）是基线日志"——**对，但都不是判决行**：
- 新发：`=== RUN` 在 `clean_new.txt:2170`；`--- FAIL: … (2.43s)` 在 **`clean_new.txt:2185`**。
- 基线：`=== RUN` 在 **`clean_base.txt:3022`**；`--- PASS: … (2.13s)` 在 **`clean_base.txt:3033`**。
照题面那两行去读会读到 `=== RUN`，量不到判决。

题面"报错的两行是 `:303` 与 `:306`，`'the re-confirmation card is %s, want L2'`（`:297`）没响"——**复认成功，且比题面更强**：本腿去**原始 CI 日志**（不是切片件）数过本枚的全部 `config_reload_223_test.go:NNN`，只有 `:303` 与 `:306` 各一次（§4.1 第 3 环）。该报而未报的那几行——`:297`（want L2）、`:333`（want exactly 1 card）、`:349`、`:354`（后来那两句操作员句）——都没响，后三条是本题的钥匙。

### 0.4 本腿推翻的题面断言（具名三条，全部现量复认）

**推翻 T1 —— "这枚是差集里唯一'真回归嫌疑'的一枚"：不成立，按题面自己的口径（"测试文件在区间内没变"）差集里有 2 枚。**

| 新增红 | 测试文件（全路径） | 该文件在 `0589fd9c..8ae4c23e` 被改几枚 | 基线存不存在 |
|---|---|---|---|
| `TestTicket223HandEditedFsLooseningCostsAnL2Card` | `cmd/wisp/config_reload_223_test.go` | **0** | 存在 |
| `TestResolvePerCallBudget` | `internal/risk/pathresolver_budget_norace_test.go` | **0** | 存在 |
| `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking` | `cmd/wisp/ticket224_assembly_test.go` | 1 | **不存在**（区间新写） |
| `TestAC246ResidentPipelineAsksThroughTheOneGate` | `cmd/wisp/resident_task_source_246_windows_test.go` | 1 | **不存在** |
| `TestPanelHostRealWindowHopAndLifecycle` | `cmd/wisp/panel_host_windows_test.go` | 7 | **不存在** |

取数命令与时刻（`2026-10-02 08:5x +08`）：`git log --oneline 0589fd9c..8ae4c23e -- <file> | wc -l` 加 `git grep -c "func <名>(" 0589fd9c -- <file>`。
`TestResolvePerCallBudget` 两发读数：`clean_base.txt:4539 --- PASS (1.81s)`、`clean_base.txt:7478 --- PASS (1.97s)`、`clean_new.txt:4045 --- FAIL (2.51s)`。
⇒ **题面说"唯一"，实测两枚，而且第二枚是那两枚里更能说明问题的**（§3(c) 第 3 条、§4.2）。

**推翻 T2 —— `0589fd9c..8ae4c23e` 被当成小增量叙述：它是 231 枚 commit 的一次整段推送。**
`git rev-list --count 0589fd9c..8ae4c23e` = **231**；`cmd/wisp` 单包顶层名册从 **191 涨到 241**（§1.3）。
任何"区间里某枚 commit 改坏了"的归因在这个规模上必须逐枚 numstat——题面硬边界 #2 已料到，本腿照做（§2）。

**推翻 T3 —— "这是'CI 那台机器上才红'的一枚"：形不对，两发都是 CI，而且是两台不同 VM。**
- 基线 `test-windows` runner = `GitHub Actions 1000001150`，新发 = `GitHub Actions 1000001157`
  （出处：`D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\orchestrator\ci-delta-1\jobs_base.txt` / `jobs_new.txt`，本腿 grep 复认）。
- ⇒ 准确的形不是"CI 才红"，是"**换了那台 runner VM 就红**"。这不是修辞：它把归因轴从题面 §3(a) 的"CI vs 本机"挪到了"VM 抽签"，而后者便宜得多（§5.5 甚至可能不用推码）。
- 题面"我本机在 HEAD 上整包复跑 `cmd/wisp` 是全绿"本腿**独立复认**：`…\orchestrator\r9-head-fullpack.txt` 顶层 `--- PASS:` **160** / `--- FAIL:` **0** / `--- SKIP:` **0**，本枚在 **`:154 --- PASS (2.18s)`**。
  ⚠ 但该文件**没记 SHA**，本腿只能证它是 `2026-10-01 23:59 +08` 的一次本机 `./cmd/wisp/` 整包。它与 `8ae4c23e` 是否同码**未定性**（§6 C5）。

---

## §1 差集复认（自己从 `ci-delta-1` 逐名抽，没照抄 `delta.md`）

### 1.1 口径（不写口径这些数不可比）
- **判红只认行首 `--- FAIL: Test`。** `t.Logf` 也带 `file.go:NN:` 前缀——本腿当场撞到一例：`clean_new.txt:4048-4050` 三行 `pathresolver_expansion_test.go:80/83/86:` 看着像红，其实是 `TestC26ExpansionMustNotRewriteOntoAnotherTree` 的 Logf，而那枚用例是 PASS。
- `clean_base.txt` / `clean_new.txt` 是**整份 CI 日志**（跨步、跨包），**不是同步切片**：顶层唯一 `=== RUN` 名分别 **1769 / 886** 枚。⇒ **跨这两个文件的绝对行号与"第几枚"不可比。**
- 比跑序必须先按 `##[group]Run` 边界把同一步切出来（本腿用：基线 `clean_base.txt:2790,3021`、新发 `clean_new.txt:1932,2184`）。
- 耗时取 `--- (PASS|FAIL): <名> (<耗时>)` 括号原值。

### 1.2 本枚两发的逐名读数

| | 全路径:行 | 判决行 | 耗时 |
|---|---|---|---|
| 基线 | `…\ci-delta-1\clean_base.txt:3022` | `=== RUN   TestTicket223HandEditedFsLooseningCostsAnL2Card` | — |
| 基线 | `…\ci-delta-1\clean_base.txt:3033` | `--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.13s)` | **2.13s** |
| 新发 | `…\ci-delta-1\clean_new.txt:2170` | `=== RUN   TestTicket223HandEditedFsLooseningCostsAnL2Card` | — |
| 新发 | `…\ci-delta-1\clean_new.txt:2185` | `--- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.43s)` | **2.43s** |

新发那两条错误行（原文，`NNN` 即 `t.Errorf` 报的行）：
- `config_reload_223_test.go:303`：`the console did not render the reload card: wisp run: 答复监听已接入……`
- `config_reload_223_test.go:306`：`the card does not name the key it is asking about: wisp run: 答复监听已接入……`
两条 `%s` 打出来的 `shown` **只有同样两行 boot 横幅**（`答复监听已接入` + `配置热加载已接管`），没有卡。⇒ **题面"卡拿到了、级别也对、只有控制台那两行没出来"复认成功**（但"没出来"这个词不准，见 §4.1 结论）。

### 1.3 整条差集（本腿自数）
基线顶层 FAIL **24** 枚 / 新发顶层 FAIL **25** 枚（同一把尺 `grep -c "^--- FAIL:"`，取数 `2026-10-02 08:5x +08`）。
- **新发独有 5 枚** ＝ §0.4 T1 那张表：4 枚在 `cmd/wisp`、1 枚在 `internal/risk`。
- **基线独有 4 枚**：`TestC21DesignTokensFourWayAgree`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`——全在 `test-core` 的 `internal/panel`，`clean_new.txt` 里**连 `--- PASS` 都没有**（分母整块没了，不是修好；同步步的 `internal/panel` 包级 `FAIL` 行一并消失可佐证）。
  ⚠ 本腿按硬边界 #2 **不读 `frontend/**`**，这 4 枚是否算缺陷归编排者按 `delta.md` 已登记的 `frontend/dist` 只有占位那条走，与本枚无关。
- **两发都红的老红（不在差集里）**：`TestL1VetoNeedsAChannelTheHostReallyWired`（24.04→22.87）、`TestTicket223PermissionDeniedSitsInItsOwnSentence`（1.12→1.40）、`TestRunPacketCarriesTheLoadedInstructionFiles`（1.11→1.25）、`TestTicket101*` 三枚、`TestComposedGateBlocksAWriteForTwoSeconds`（301.12→301.70）、`internal/risk` 那 7 枚。
  ⇒ 本枚同文件的同族 `TestTicket223PermissionDeniedSitsInItsOwnSentence` **两发都红**：这个包的红底本来就高，本枚不是孤例。
- **cmd/wisp 分母涨了多少**：新发该步步尾 `clean_new.txt:3437` 逐字 `=== RUN=241  --- PASS=147  --- FAIL=11  --- SKIP=2`；基线同步步 `=== RUN=191`。包级耗时 `424.719s → 579.297s`（**+36.4%**）。

### 1.4 本腿自破的一把坏尺（先怀疑尺，没怀疑数据）
第一把我做了"比较本枚在两发里各是第几枚跑的"：读出基线 **267**、新发 **115**，差 152。
**这个读数不可用**，是尺坏了不是数据怪：`clean_base.txt`（1 201 637 B）与 `clean_new.txt`（699 575 B）是整份日志、含多包同步，唯一顶层名 1769 / 886——**两把尺量的不是同一个分母**。改用 1.1 第 3 条的按步切法才拿到可用读数（1.5）。
⇒ **给下一枚腿：任何拿这两个 `clean_*` 做"位置/序号"比较的尺都是坏尺，只能做步内比较。**

### 1.5 有效的那把尺：本枚**之前**跑了什么（逐名逐序比对）
切法见 1.1。结果：
- 顶层 `=== RUN` 计数：基线 **13** 枚、新发 **13** 枚。
- 逐名**有序** `diff`（本腿实跑 `diff <(sed -n '2790,3021p' clean_base.txt | grep '^=== RUN   Test') <(sed -n '1932,2184p' clean_new.txt | grep '^=== RUN   Test')`）：输出只有一条 `13a14 > 本枚自己`（新发那把右边界把本枚的 `=== RUN` 圈了进去）⇒ **前 13 枚逐名逐序完全相同，零行差。**
- 前 13 枚里**没有一枚是区间新增的**：`comm -13 cmdwisp_base_names.txt before_ours_new.txt` = **空**。
- 泄漏压相同：本枚之前 `goroutine outside the D38 roster` WARN 基线 **12** 行 / 新发 **12** 行；整步 84 / 98（多的 14 行全在本枚**之后**）。
- 包内无并发：`grep -rn "t.Parallel()" cmd/wisp/ --include=*.go` = **0 命中**；该步命令行（`logs_new_all.txt` 里 `runtests.sh: go test exited 1 - packages=[./cmd/wisp/ -count=1 -skip ^(...)`）不带 `-parallel`。

⇒ **本枚的上游在两发里逐字同形**（同 13 枚、同名、同序、同泄漏量、同串行模型）。这一条把 §3(b) 的"新增用例污染了本枚上游"整支钉死。

---

## §2 那枚区间到底动了什么（231 枚逐枚 numstat，未用 range diff 归因）

### 2.1 名册
- `git rev-list --count 0589fd9c..8ae4c23e` = **231**（首 `8ae4c23e`、尾 `0bc0f249`）。全量落 `...\ci223-1\raw_range_roster.txt`。
- **逐枚** `git show --numstat --format='' <c> -- 'cmd/wisp/**'`，取非空者 ⇒ **29 枚碰了 `cmd/wisp`**。清单（短 sha / 时间 / 标题 / 文件数）落 `...\ci223-1\raw_commits_touching_cmdwisp.txt`。

### 2.2 只挑"能碰到审批卡渲染到控制台这条路径"的产码枚
行号钉在 `8ae4c23e`。区间史尺：`git log --oneline 0589fd9c..8ae4c23e -- <file>`。

| 产码件 | 在这条路径上的位置（`8ae4c23e` 行号） | 区间内被改 | 改的是不是渲染 |
|---|---|---|---|
| `cmd/wisp/config_reload.go` | 卡文本生产者 `:246 PendingApproval` / `:250 Reason: fmt.Sprintf(...)`；**本枚后来读到的那两行** `:195-196`、`:201-202` | **0 枚** | — |
| `cmd/wisp/approval_always.go` | `bookWaitingState`（`:171-184`），正好卡在 record 与印之间 | **0 枚**（`-S bookWaitingState` 亦 0 命中） | — |
| `cmd/wisp/run.go` | **渲染本体**：`Prompt` `:1312`、`record` `:1322`、`bookWaitingState` `:1329`、**卡字 `Fprintf` `:1332`** | **4 枚**：`3b78c246` `bb37fac2` `3edc11d4` `2071f59e` | **见下** |
| `cmd/wisp/approval_reply.go` | `nativeCards.pending` `:158`（`awaitCard` 读的就是它） | 3 枚（`3b78c246` `bb37fac2` `3edc11d4`） | 只动答复路由与 help 文案 |
| `internal/agent/approval/gate.go` | `PendingApproval` `:496` → **调用点 `:530 g.ui.Prompt(ctx, p)`** | 2 枚（`3b78c246` `2a0175ba`） | **`ui.Prompt`/`promptFor`/`q.push` 调用点零改** |

`gate.go` 调用点那条的尺（本腿实跑）：
`git diff 0589fd9c 8ae4c23e -- internal/agent/approval/gate.go | grep -nE "^[+-].*(ui\.Prompt|promptFor|g\.q\.push)"` = **空输出**。区间 diff 全是 `Grants`/`GrantRecorder` 的**增补**（票 224 会话授权记账），不碰展示。

**唯一一枚真碰到渲染本体的改动**（`run.go`，`bb37fac2`+`3edc11d4` 这对抗）：
```
-	rt.ui = &consoleApprovalUI{out: s.stdout, run: rt, live: rt.liveCards}
+		rt.ui = &consoleApprovalUI{out: s.stdout, run: rt, live: rt.liveCards}
```
减号在基线位置，加号在 `8ae4c23e` 的 **`cmd/wisp/run.go:585`**，**参数字面一字未改**。它被搬进了票 246 AC#7「一个进程只许一枚 approval.Gate」的两档姿态：`if s.gate != nil`（`:574`）那一支把 `rt.ui` 置 nil（`:581`），本枚走的是 `else`（`:584-585`）。
⇒ 对本枚**行为等价**，本腿用盘上事实钉它：`cmd/wisp/approval_reply_201_test.go` 的 harness `runSpec` 字面量（`:167` 起，`return runTextTask(runSpec{`）**没有 `gate:` 字段**（`grep -n "gate:"` 该文件 = 0 命中）⇒ `s.gate == nil` ⇒ 走 `else` ⇒ 构造与基线逐字相同。

### 2.3 §2 层结论
29 枚碰 `cmd/wisp` 的 commit 里，**没有一枚改过"卡的文字怎么进 stdout"这条语句**：`Prompt` 函数体无 diff、`config_reload.go` 0 改、`bookWaitingState` 0 改、`gate.go` 的 `ui.Prompt` 调用点 0 改。
而且"`record` 早于 `Fprintf`"这个**可竞争顺序在基线锚点上就已经是这样**：
尺（本腿实跑）`git show 0589fd9c:cmd/wisp/run.go | sed -n '/func (u \*consoleApprovalUI) Prompt/,/^}/p' | grep -nE "record|bookWaitingState|Fprintf"` ⇒ 体内第 **11** 条 `u.live.record(p)`、第 **18** 条 `u.run.bookWaitingState("ui-prompt")`、第 **21** 条 `fmt.Fprintf(u.out, "\n[确认 %s %s] %s\n", p.Level, p.Tool, p.Reason)`。
⇒ **产码侧那个"卡先记账、后印字"的窗口，区间前后同形。§3(d) 真回归在读码层被证伪。**

---

## §3 机制候选，逐个配"能证伪它的尺"

可能性由高到低：**3(c) ≈ 3.0 竞争窗口（同一件事的两面）＞ 3(a) CI 环境口径 ＞ 3(b) 顺序依赖（已判死其主形）＞ 3(d) 真回归（已证伪）**。

### 3.0 先把"窗口"说清楚（后面三支都挂在它上面）
本枚的断言形状是**一次读、不等**（`cmd/wisp/config_reload_223_test.go`，`8ae4c23e` 与当前 HEAD 逐字节相同，本腿 `git diff --stat 8ae4c23e -- <该文件>` = 空）：
`:295 card := r.awaitCard(...)` → `:301 shown := r.h.out.String()` → `:302`/`:305` 两次 `strings.Contains(shown, …)`。
`awaitCard`（`:161-180`）轮的是 `r.rt.liveCards.pending()`（`cmd/wisp/approval_reply.go:158`），tick 间隔 `:165 time.NewTicker(20 * time.Millisecond)`；字是 `run.go:1332` 才写的。
**同一个 goroutine（配置轮询腿）的 program order**：
`config_reload.go:143 case <-t.C` → `:152 reloadOnce()` → `:153 mgr.CheckAndReload()` → …… `:246 rt.gate.PendingApproval(...)` → `gate.go:530 g.ui.Prompt(ctx, p)` → `run.go:1322 u.live.record(p)` → `run.go:1329 u.run.bookWaitingState("ui-prompt")`（内含 `auditf` `run.go:899`，其中 `:901` 一次 stderr 写 + `:902` 一次 sink handler）→ `run.go:1332 Fprintf(u.out, "[确认 …")`。
而 `reportReload` 是在 `CheckAndReload()` **返回之后**才被 `reloadOnce` 调（`config_reload.go:161`）。
⇒ **`record` 与卡的 `Fprintf` 之间隔着一次日志落盘；`awaitCard` 看不见那次落盘，它只看得见 ledger。**
⇒ **输的条件**：那一下落盘＋调度 > ~20ms。**赢的条件**：< 20ms。没有产品语义参与。

### 3(a) CI 与本机口径差异（缺 env / `%APPDATA%` 形状 / 中文编码换行）
- **支持**：本机与 CI 的 temp 根确实不同形——本机 `r9-head-fullpack.txt` 里是 `C:\Users\swq\AppData\Local\Temp\TestTicket223HandEditedFsLooseningCostsAnL2Card3203839337\002\logs`；两发 CI 都是 `C:\Users\RUNNER~1\AppData\Local\Temp\…`（`clean_base.txt:3024`、`clean_new.txt:2172`）。
- **反对（三条独立）**：
  1. **轴不对**：基线与新发**都是 CI**，`%APPDATA%` 形状**两发逐字相同**（都是 `RUNNER~1` 8.3 短名形）。两发同形的条件解释不了只有新发红。
  2. **中文编码 / `strings.Contains` 不命中：证伪。** 两把 needle（`[确认 L2 config.reload]`、`fs.allowed_dirs`）是**编译进测试二进制的 Go 字面量**，比较是字节级；`config_reload_223_test.go` 区间内 **0 枚 commit**（§0.4 表），源码字节没动。且同一份日志里中文照样印得出来（`clean_new.txt:2180-2181` 那两行 `wisp run: 答复监听已接入…`／`配置热加载已接管…` 完整可读）。
  3. **环境变量：证伪到具体语句。** 渲染语句 `run.go:1332` 的 `fmt.Fprintf(u.out, "\n[确认 %s %s] %s\n", p.Level, p.Tool, p.Reason)` 不读 env；`Reason` 由 `config_reload.go:250` `fmt.Sprintf("config.toml 里手改的这一条会放宽 [%s]：%v。…", section, keys)` 拼，也不读 env。
     ⚠ 这一条是**读码**结论（"我没在这条路上看到 `os.Getenv`"），不是"跑过证明没有"。
  4. **`%APPDATA%` 形状与红负相关**：本机是 `C:\Users\swq\…` 那一形，本机绿；CI 是 `RUNNER~1` 那一形，CI 一发绿一发红。⇒ 形状本身不是因。
- **能一刀判死它的尺**：〔未自跑，本腿禁 go test〕§5.1 尺 A（单发）。
  若同一台机器上单发稳定绿 ⇒ "CI 有系统性口径差异"整支判死（同一枚单跑绿 ⇒ 差异不在环境、在时序）。
  更强的判死已经由 §4.1 拿到（盘上日志自证渲染发生过），**本腿不需要为这一支再跑任何东西**。

### 3(b) 顺序依赖 / 并发（先跑的用例留下状态）
- **支持**：本仓确有同族先例（题面给的"一枚起真窗/锁线程的钉自己就是毒源，单跑永远绿、整包才红"）；新发整包多 50 枚、慢 36.4%（§1.3）；本包内确实有 `mockllm-stdout-reader` 这类**前序用例留下的活 goroutine**（本枚之前 12 行 WARN）。
- **反对（四条，都是本腿现量，见 §1.5）**：
  1. 本枚之前 13 枚，两发**逐名逐序零行差**。
  2. 新增那 50 枚里**没有一枚**排在本枚之前（`comm -13 cmdwisp_base_names.txt before_ours_new.txt` = 空）。
  3. 泄漏压相同：12 行 vs 12 行。
  4. `t.Parallel()` 0 命中 + CI 命令行无 `-parallel` ⇒ 顶层严格串行，"别的用例同时抢 CPU"在本包内不成立。
- **边界说死（不许把没证到的写成已证）**：以上判死的是"**区间新增用例污染了本枚上游**"这一支，**没有**判死"**同一次 run 里前 13 枚留下的 OS 级残债**（Defender 扫 temp、句柄数、SQLite 文件堆积）"。那需要 §5.1 单发 vs §5.3 整包在同一台机器上对比才分得开（§6 C2）。
- **能一刀判死它的尺**：〔未自跑〕
  ```
  cd "D:/work/workspace/projects plans/Wisp" && \
  PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v \
    -run '^(TestTicket223RunArmsTheReloadTick|TestTicket223HandEditedFsLooseningCostsAnL2Card)$' ./cmd/wisp/
  ```
  期望读数：两枚都 `--- PASS`。若这一发**红**而单发（§5.1）绿 ⇒ 毒源就是紧邻的前一枚 `TestTicket223RunArmsTheReloadTick`（两发都 PASS，但新发慢了 0.15s：2.14→2.29），本支复活为真。

### 3(c) 负载 / VM 抖动（**本腿投这一支票最多**）
- **支持（四条，含一把跨包独立佐证）**：
  1. **两台不同 VM**：`1000001150` vs `1000001157`（§0.4 T3）。
  2. **同族集体变慢**：本枚之前那 13 枚里 6 枚基线 ~1.0s 的涨到 ~1.4-1.8s（`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` 1.22→1.78、`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` 1.03→1.57、`TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused` 1.01→1.47、`TestReplyListenerAllowsAnL2CardFromTheNativeSide` 1.02→1.49、`TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel` 1.00→1.43、`TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` 1.03→1.74）；本枚 2.13→2.43（+14%）。
  3. **独立佐证（最硬一条）**：`internal/risk/pathresolver_budget_norace_test.go` 的 `TestResolvePerCallBudget` 是一枚**纯挂钟性能钉**——`:12 const resolveBudget = time.Millisecond`、`:36 if time.Duration(ns) > resolveBudget { t.Fatalf(...) }`，判的是 `testing.Benchmark` 的 `NsPerOp()`。它的**测试文件区间内 0 改**，且 `internal/risk/` **整包区间内 0 枚 commit**（尺：`git log --oneline 0589fd9c..8ae4c23e -- internal/risk/` = **空输出**）。它还是从 PASS 翻成 FAIL。
     ⇒ **一枚"只可能因机器慢而红"的用例，与我们的用例在同一次 run 里一起翻红，而它那条路径上连代码都没变。** 单靠"本枚这处有竞争窗口"解释不了它。
  4. 包级耗时同向：`424.719s → 579.297s`（+36.4%）。
- **反对（本腿不许把没证到的写成已证）**：
  - **不是均匀变慢**：`TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` 反而快了（1.96→1.34），`TestL1VetoNeedsAChannelTheHostReallyWired` 也快（24.04→22.87）。⇒ "整台 VM 慢 35%"这个说法被这两枚**直接反驳**；更像"那 6 枚共同跨过的 1s tick 边界被更多次跨过"（与 `reloadCaseBudget`/1s 轮询形状一致）。⇒ 3(c) 只能写成"**VM 抖动更大**"，不能写"VM 更慢"。
  - 新发多出的 50 枚（含真起 WebView2 窗的 `TestPanelHost*`、`TestAC13*`）排在**本枚之后**（§1.5），所以"包变长把本枚挤慢了"这一条**不成立**。
- **"它为什么不能解释成 C18 那 300s 常量"（题面点名要的判断）**：
  `internal/config/schema.go:450` `ConfirmTimeoutSec int \`toml:"confirm_timeout_sec" default:"300"\`` 是本仓那族 `300.0x s` FAIL 的来源。现量（尺：`grep -cE '^--- FAIL: Test.*\(30[0-9]\.'`）：**基线 1 枚 / 新发 1 枚**，同一枚 `TestComposedGateBlocksAWriteForTwoSeconds`（301.12s / 301.70s），**两发都红 ⇒ 不在差集里**。本枚耗时 **2.43s**，与该常量差两个数量级；本枚那 2 条是 `t.Errorf`（`:303`/`:306`），不是任何超时的 `t.Fatal`。
  ⇒ **"300s 审批超时"这一支与本枚无关，判死。**
- **能坐实/判死它的尺**：〔未自跑〕§5.2（`-count=25` 拿频率）＋ §5.4（那枚挂钟钉同机重复采样）。
  两枚一起按次翻 ⇒ 3(c) "同一只手"坐实；只有本枚翻、那枚 25/25 稳 ⇒ 两支解耦，本腿 §4.2 作废。

### 3(d) 真回归：区间里某枚产码把控制台渲染改坏了
- **支持**：只有形状上那一条——`run.go` 区间被 4 枚碰过，其中 `bb37fac2`+`3edc11d4` 确实**搬动过** `consoleApprovalUI` 的构造（§2.2）。
- **反对（四条，前两条已足）**：
  1. 搬动后的构造**参数字面逐字未改**，且本枚走的就是那一支（`approval_reply_201_test.go:167` 的 `runSpec` 无 `gate:` ⇒ `run.go:584-585` `else`）。
  2. **渲染本体零改动**：`Prompt` 体无 diff；`config_reload.go` 0 改；`bookWaitingState` 0 改；`gate.go` 的 `ui.Prompt` 调用点 0 改；`record 早于 Fprintf` 在 `0589fd9c` 上就是现状（§2.3 那把尺）。
  3. **本枚自己的日志已经把"卡没印出来"证伪了**（§4.1，全链只读盘，不依赖任何未跑的尺）。
  4. 若真坏了，同文件同族的 `TestTicket223PermissionDeniedSitsInItsOwnSentence` 该**从绿转红**；它其实**两发都红**（1.12s/1.40s），不在差集里。
- **能判死它的尺**：**已由 §4.1 判死**（读盘）。补一把给下一枚腿兜底：〔未自跑〕§5.1，期望 `--- PASS`；若同机**稳定红**，本腿 §4.1 的链有洞，立刻升回 (d)。

---

## §4 读码层能到哪儿（把话说死的边界）

### 4.1 〔已证：只读盘上 CI 日志，不跑任何东西〕控制台**确实印了那张卡**，是本枚读早了
闭合链，每环都带全路径行号（行号钉 `8ae4c23e`）：

1. **同一个 writer 对象**。`cmd/wisp/approval_reply_201_test.go:169` 把 `h.out` 交给 `runSpec.stdout`；`cmd/wisp/run.go:370` `rt := &agentRuntime{spec: s, stdout: s.stdout, …}`；`cmd/wisp/run.go:585` `rt.ui = &consoleApprovalUI{out: s.stdout, run: rt, live: rt.liveCards}`。
   ⇒ 卡的文字与热加载的句子写的是**同一枚 `*syncWriter`**（`approval_reply_201_test.go:64-79`：`Write` 与 `String` 共用一把 `sync.Mutex` 加一个 `bytes.Buffer`，**只追加、不截断**）。
2. **同一个 goroutine，卡的 `Fprintf` 严格在前**。卡字在 `run.go:1332`（`CheckAndReload` 调用栈内，`config_reload.go:153`）；`放宽已经过 L2 重新确认` 在 `config_reload.go:195-196`，由 `reportReload`（`:167`）打印，而 `reportReload` 是 `reloadOnce` 在 `CheckAndReload()` **返回之后**才调的（`:161`）。两者都在 ticker 腿上（`:143 case <-t.C`）。
   ⇒ 单 goroutine 的 program order 保证：**卡字先进缓冲，那句后进。**
3. **本枚在同一次红的那一发里读到了"后进"的那句**。原始 CI 日志（不是切片件）里本枚的全部 `config_reload_223_test.go:NNN` 只有两个：
   尺（本腿实跑）`grep -a -oE "config_reload_223_test.go:[0-9]+" …\ci-delta-1\logs_new_all.txt | sort | uniq -c` ⇒ `1 …:303` / `1 …:306`；`raw_tw_new.txt` 同一读数。
   ⇒ `:333`（want exactly 1 card）、`:349`（`放宽已经过 L2 重新确认`）、`:354`（`仍按启动时建好的 C26 名单`）**一行都没响 = 三条都过了**。`t.Errorf` 不中断执行，所以这三条确实执行到了。
4. **⇒ 结论**：到 `:349` 那一刻缓冲里已有 `config_reload.go:196` 那句；由 1＋2，卡字**必然也已经在里面**。而 `:301` 那次读没看到它。
   ⇒ **不存在"控制台没渲染"这回事；是断言读早了。** `[确认 L2 config.reload]` 与 `fs.allowed_dirs` 两把 `strings.Contains` 挂在**同一次**未加界的 `:301` 单读上，所以它们**必然一起响**——正是日志的样子。`:297` 读的是 ledger 不是控制台，所以它不响。
   **题面那句"只有控制台渲染那两行没出来"，准确写法是"只有那一次未加界的 stdout 单读没赶上"。**
   ⇒ 原始日志里只报了 303/306，⇒ 该文件其余 **56** 条候选报文行**都没响**：`:297`(L2 级别)、`:330`(Config() 冻住没答)、`:333`(卡片数≠1)、`:339`/`:343`(ledger 记账)、`:349`(放宽已经过 L2)、`:354`(C26 名单)、`:445`(同形状的另一枚卡) 等。
   尺（本腿实跑，`8ae4c23e`＝当前 HEAD 逐字节同）：`grep -cE "t\.(Errorf|Fatalf|Fatal)\(" cmd/wisp/config_reload_223_test.go` ⇒ 全文共 **58** 条候选报文行，逐条比对原始日志后只有 2 条命中。
5. **这不是新知识，是本文件自己登记过的形状**。同文件 `:129-139` 的注释逐字写着 `awaitStdout` 存在的理由是"reload 路径先写审计、后写操作员那句"，并给了实测频率 `:134` "measured 1/25 here, 1/5 on the v1 leg, 1/15 by the orchestrator"；`:478` 记的是 AC#5 那处正是改用 `awaitStdout` 收口的。
   ⇒ **`:301` 是同一批改造的漏网之鱼。** 文件里现在还剩这些 `h.out.String()` 出现行（尺：`grep -n 'h\.out\.String()' cmd/wisp/config_reload_223_test.go | cut -d: -f1 | tr '\n' ','`，实跑读数 `65,147,176,206,253,254,274,276,301,348,350,353,378,380,421,422,443,`）：其中 `65/147/176/206/254/276/350/380/422` 在 helper 内部或只是失败信息的转储；真正的"**awaitCard/awaitAudit 之后一次读、不等**"是 **`:301`** 与 **`:443`** 两处（`:443` 是 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`，`:442` 刚 `awaitCard`、`:443` 立刻单读、`:445` 就 `strings.Contains(shown, "risk.permission_mode")` —— **与本枚完全同一形状**，只是这一发没轮到它）。另有一批"审计之后单读操作员句"的次级同形：`:253`、`:274`、`:348`、`:353`、`:378`、`:421`。本腿**不修**（边界 #5）。

### 4.2 〔已证：同上，只读盘〕`internal/risk` 那枚同时翻红与本枚同向
`TestResolvePerCallBudget` 是纯 `NsPerOp` 挂钟钉（`internal/risk/pathresolver_budget_norace_test.go:12/:36`），测试文件与 `internal/risk/` **整包区间内 0 改**，却 PASS→FAIL。
⇒ **已证的是**：一次 run 里可以同时把"一枚有竞争窗口的功能钉"和"一枚纯速度钉"一起翻红，**无需任何代码变化**。
⇒ **已证的边界止于此**：这不等于证明本枚就是被同一只手翻的（那是 3(c)，需要 §5.2/§5.4 的频率）。

### 4.3 〔读码推的，未证〕窗口被什么撑开
- `record → bookWaitingState → auditf` 里那一次 `auditf`（`run.go:899-902`）＝ 一次 `fmt.Fprint(rt.stderr, …)` **加**一次 `rt.spec.sink.logger().Info(...)`，而 `cmd/wisp/logsink.go:122-126` 走的是 `slog.New(s.pipeline.Handler())`。**本腿没读进 `internal/observe` 的 pipeline 内部**，没确认它是同步落盘还是带 channel 的批处理。若带缓冲，CI 上被撑过 20ms 的具体机制就在那儿——**这一步是推测**。
- "新 VM 抖动更大"而非"更慢"：由 3(c) 的反例给出方向，但本腿没有任何 CPU 数、没有 runner VM 规格、没有同机重复采样。它只是"对得上数据的假设"。

### 4.4 〔必须真跑才答得出〕三件事，读码给不出
1. **单发会不会绿**（§5.1）——"竞争窗口 vs 别的"的分水岭。
2. **红频是多少**（§5.2）——只有拿到频率才能写"flaky"，才能与 `:134` 那句 1/25、1/5、1/15 并账。
3. **整包跑序那一发会不会因新增 50 枚而稳定复现**（§5.3）——若整包稳定红而单发稳定绿，3(b) 的"OS 级残债"版本就活了，与 3(c) 分家。

### 4.5 未定性（明写，不糊）
**"控制台那两行为什么在 CI 上这一次没印出来"——本腿不给机制级的因。**
本腿证得到的是**反命题已死**："卡没渲染 / 渲染路径被区间内某枚产码改坏"已由 §4.1 判死（盘上日志自证，不需再跑）。
本腿证不到的是**正命题的因**：窗口被撑开（落盘/调度 > 20ms）、还是同机 OS 残债、还是 VM 抖动。
**要区分这一支需要跑哪一发**：§5.1 单发连跑 3 次 ＋ §5.2 `-count=25`。红了给出频率 ⇒ 窗口/抖动坐实；三发全绿且 `-count=25` 全绿 ⇒ §4.1 的链有洞，回到 (d)。

---

## §5 下一枚复现腿的确切派单料

### 5.0 先验 PATH（不做这步，拿到的"绿"是假的）
`cmd/wisp` 测试二进制链 sherpa-onnx；**缺 DLL 时 `go test` 直接 `exit status 0xc0000135`，日志里一个 `--- FAIL` 都没有**（`cmd/wisp/config_reload_223_test.go:16-17` 写了这条；`scripts/wisp-cli-tests.sh:7`、`:13-20` 写了同一事实的四个读数；CI 侧同一事实是 `.github/workflows/ci.yml:455-465`）。⇒ **那读数是"根本没跑"，不是绿。**
- DLL 实际落位（本腿现量 `find third_party/sherpa-onnx -maxdepth 2`）：**`third_party/sherpa-onnx/` 目录本身**，内含 `onnxruntime.dll`、`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`、`.cache-manifest.json`。**没有 `bin/` 子目录**（`ls third_party/sherpa-onnx/bin` = 不存在）。`scripts/wisp-cli-tests.sh:109` 的 `export PATH="$dll_dir:$PATH"` 就是这个形状。
- 跑之前先自证 DLL 到位（不跑测试，只数名；本腿实跑过，读数 **3**）：
  ```
  cd "D:/work/workspace/projects plans/Wisp" && ls third_party/sherpa-onnx/*.dll | wc -l
  ```
  期望 **3**。不是 3 ⇒ 先 `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/fetch-deps.ps1`（或 `scripts/build.ps1 -Env dev`，即 CI 的 `ci.yml:453` 那一步）把 DLL 请回来，别急着跑。

### 5.1 尺 A：单发（先跑这把，它是分水岭）
```
cd "D:/work/workspace/projects plans/Wisp" && \
PATH="$PWD/third_party/sherpa-onnx:$PATH" \
go test -count=1 -v -run '^TestTicket223HandEditedFsLooseningCostsAnL2Card$' ./cmd/wisp/ 2>&1 \
| tee .scratch/wisp/probes/<leg>/A-alone.txt
```
**期望读数**：`--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card (2.0s-2.6s)`（参照：本机 2.18s、CI 基线 2.13s、CI 新发 2.43s）。
- **绿** ⇒ §3(d) 真回归**当场判死**，结论收在 §4.5：竞争窗口 ＋ 跑序/机器。下一发走 5.2。
- **红，且红的就是 `:303`/`:306`** ⇒ 本腿 §4.1 那条链**有洞**，判死作废，立刻升真回归候选。
- **红，但红在别的行** ⇒ 与本腿读数不同形，先把口径报回来再往下跑。
建议**连跑 3 次**（同机同码），把 3 个耗时都记下来。

### 5.2 尺 B：拿频率（判"是不是 flaky"唯一把那把尺）
```
cd "D:/work/workspace/projects plans/Wisp" && \
PATH="$PWD/third_party/sherpa-onnx:$PATH" \
go test -count=25 -v -run '^TestTicket223HandEditedFsLooseningCostsAnL2Card$' ./cmd/wisp/ 2>&1 \
| tee .scratch/wisp/probes/<leg>/B-count25.txt
```
数（这两行就是数）：
```
grep -c '^--- PASS: TestTicket223HandEditedFsLooseningCostsAnL2Card' .scratch/wisp/probes/<leg>/B-count25.txt
grep -c '^--- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card' .scratch/wisp/probes/<leg>/B-count25.txt
```
**期望读数**：两数之和 **= 25**（本腿按 §1.1 口径写死：只认行首）。
- **FAIL ≥ 1** ⇒ **flaky 坐实**（同码同环境按次翻，产品语义未变）。把 n/25 报进台账，与 `config_reload_223_test.go:134` 那句"1/25 here, 1/5 on the v1 leg, 1/15 by the orchestrator"**并账**。
- **FAIL = 0（25/25）** ⇒ 本机复现不了，3(c) 未坐实；走 5.3，并在 CI 侧优先走 **5.5**（那比任何 `go test` 都便宜）。

### 5.3 尺 C：整包跑序那一发（只有与 5.1/5.2 对比才有意义）
```
cd "D:/work/workspace/projects plans/Wisp" && \
PATH="$PWD/third_party/sherpa-onnx:$PATH" \
go test -count=1 -v ./cmd/wisp/ 2>&1 | tee .scratch/wisp/probes/<leg>/C-fullpkg.txt
```
**本腿能给的基线就是编排者那发本机整包**：`r9-head-fullpack.txt` 现量顶层 `--- PASS:` **160** / `--- FAIL:` **0** / `--- SKIP:` **0**，本枚在 `:154` PASS (2.18s)。
- 若当前 HEAD 上仍是 `160 / 0 / 0` ⇒ "本机整包不可复现"复认，本枚的红是 **CI 侧**属性。
- 若冒出红 ⇒ 先记**是不是本枚**，并连口径一起报（哪个包、哪一发、`-count` 几）。
⚠ **分母核对，别拿 241 当 160**：CI 的 `=== RUN=241 / PASS=147 / FAIL=11 / SKIP=2`（`clean_new.txt:3437`）里 241 **含子测试**；顶层**结果**数 147+11+2 = **160**，与本机那 160 是同一把尺。跨这两把报数必须写明用的哪一把。

### 5.4 尺 D：把两枚"未变而翻红"钉到同一只手上（顺手，一发）
```
cd "D:/work/workspace/projects plans/Wisp" && \
go test -count=25 -v -run '^TestResolvePerCallBudget$' ./internal/risk/ 2>&1 \
| tee .scratch/wisp/probes/<leg>/D-riskbudget.txt
```
（`internal/risk` 不链语音，理论上不需要 sherpa PATH；**但本腿没有验过这一条**，下一枚腿先加 PATH 再跑，别拿本腿没验的断言去省那一步。）
**期望读数**：这台机器若和 CI 那台一样抖 ⇒ `FAIL ≥ 1`。
- **两枚都按次翻** ⇒ 3(c) 坐实，本枚从"真回归嫌疑"降为"CI 抖动受害者"，两包一起报。
- **只有本枚翻、那枚 25/25 稳** ⇒ 两支解耦，本腿 §4.2 那条"同向"作废（§6 C1 的复活条件）。

### 5.5 CI 侧那一发（**最便宜，建议排第一**）
盘上已经躺着枚**同码二次采样**：`…\ci-delta-1\delta.md` §0 自己记的 `run 36940536372`——`workflowName=ci`、`event=schedule`、`headSha=8ae4c23e`、`status=completed`、`conclusion=failure`。
⇒ **判"CI 上是不是按次翻"可能根本不用推新码**：拉那枚 run 的日志，按 §5.0-5.1 的口径数本枚与 `TestResolvePerCallBudget` 两枚的名字。
⚠ 本腿**没有**拉它：硬边界只许"读文件 + 读盘上已有的 CI 日志 + git"，而 `ci-delta-1/` 里没有这枚 run 的日志件（`ls` 现量确认，见 §0.2）。**这是本腿交的第一优先待办。**

### 5.6 派单必须带上的三句口径（本腿现场踩过）
1. **判红只认行首 `--- FAIL: Test`。** `t.Logf` 也带 `file.go:NN:` 前缀：`clean_new.txt:4048-4050` 那三行 `pathresolver_expansion_test.go:80/83/86:` 是 `TestC26ExpansionMustNotRewriteOntoAnotherTree` 的 Logf 而那枚 PASS——照"有 `:行号` 就是红"的尺会把绿读成红。
2. **`clean_base.txt` 与 `clean_new.txt` 是整份日志、不是同步切片**（唯一顶层名 1769 / 886），跨文件的行号/序号不可比（§1.4 本腿就是在这儿摔的）。比跑序先按 `##[group]Run` 切。
3. **`0xc0000135` 不是红、是没跑**（§5.0）。
4. **行号要钉锚**：`cmd/wisp/run.go` 的行号在 `8ae4c23e` 与 `0d87a681` 之间差 17（0.-1）。引源码行号时先说清钉在哪。

---

## §6 判不动的地方（具名，每条：为什么判不动 ＋ 谁能判）

**C1 —— 窗口到底被什么撑开（stderr 写 / sink 落盘 / 调度）。**
本腿读到了 `run.go:1322 → :1329 → :1332` 这条链和 `logsink.go:122-126` 的 `s.pipeline.Handler()`，但**没读进 `internal/observe` 的 pipeline 内部**确认它是同步 `Write` 还是带 channel 的批处理，所以"CI 上那一下 > 20ms"具体是哪一次系统调用判不动。
谁能判：能读 `internal/observe` pipeline 的读码腿（本腿可以读，只是没排到）；或由 §5.4 那发把它降级成无关。

**C2 —— 新增 50 枚用例的**进程级**成本（二进制更大、包初始化更多）会不会落到本枚头上。**
`t.Parallel()` 0 命中保证了同包顶层用例不同时抢 CPU，但**同一枚二进制的 loader/init 是进程级**的，本腿无法从日志把它切出来。
谁能判：§5.3 整包发 与 §5.1 单发**在同一台机器上对比**；或测 `go test -c` 产物的 init 耗时（要跑 exe，本腿禁）。

**C3 —— `TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole` 为什么反而快了（1.96s→1.34s）。**
这条反例是 3(c) 从"VM 更慢"降级为"VM 更抖"的唯一支点，本腿只知道它快了、不知道为什么。若它是**固定 sleep 被跳过**而不是抖动，3(c) 的解释力还要再降。
谁能判：读 `cmd/wisp/approval_reply_201_test.go` 里那枚的等待结构 ＋ 同机重复采样（§5.2/§5.4 频率法）。本腿只报了数，没判因。

**C4 —— 本枚是不是历史上就翻过（同码不同 run 的红/绿史）。**
本腿只有两发 CI（`0589fd9c`、`8ae4c23e`）加一枚没记 SHA 的本机 fullpack。`delta.md §0` 明确记了第三发 `run 36940536372`（schedule、同 `8ae4c23e`、已跑完、也红），**其日志不在盘上**，本腿按边界不能去拉。
谁能判：编排者（他有 `gh`，本腿没有也不该用）。**这是最便宜的一手**，见 §5.5。

**C5 —— 本机那发全绿到底锚在哪枚 commit。**
`r9-head-fullpack.txt` 只有测试输出、**文件里没有 SHA**。题面说"在 HEAD 上"，本腿能确认的只有"2026-10-01 23:59 +08 的一次本机 `./cmd/wisp/` 整包，160/0/0"。它与 `8ae4c23e` 是否同码判不动；而当前 HEAD 已是 `0d87a681`（比 `8ae4c23e` 新 **21** 枚（尺：`git rev-list --count 8ae4c23e..HEAD` = 21，取数 `2026-10-02 09:0x +08`））。
谁能判：跑那发的腿自己补记锚；或编排者在 §5.3 那发上现量重取。

**C6 —— `internal/panel` 那 4 枚"消失的红"与本枚有没有关系。**
本腿按硬边界 #2 **不读 `frontend/**`**，也没读那 4 枚的判据来源；`delta.md` 已把它们归到 `frontend/dist` 只有占位文件撞 `tools/d22scan/runtests.sh:98` 的"跳过即红"。本腿**不复核这一条**，也不据此下任何判。
谁能判：面板/前端侧的腿，或编排者按已排的那条"摘回默认档"决定处理。

**C7 —— 该不该修（把 `:301` 那次未加界的单读改成 `awaitStdout`）。**
本腿**判不动且明确不判**：本仓规矩是"改契约＝人工批准""不许为了变绿放宽任何断言"（AGENTS.md §1.1）。把 `cmd/wisp/config_reload_223_test.go:301-307` 改成 `awaitStdout` 是**给断言加等待**——同文件 `:137-138` 那句"它不是 assertion、也不是 loosening：句子永远不来照样红"是**该件作者的主张，不是本腿的裁定**。
谁能判：票 223 的 owner 与编排者裁；裁下来才轮到写码腿动那枚文件。

**C8 —— 本腿所有 `go test` 尺的真实读数。**
硬边界 #1 禁本腿跑任何测试/构建 ⇒ §3 的 (a)/(b)/(c)/(d) 四把"能判死它的尺"与 §5 全部尺**都是未自跑的**，逐条标了〔未自跑〕。
本腿自跑过的只有：§0/§1/§2 的 grep / sed / git / diff / comm / find 读数，与 §4.1 那条**只读盘上日志**的推理链。
谁能判：下一枚复现腿（§5）。

**C9 —— `TestTicket223PermissionDeniedSitsInItsOwnSentence` 那枚两发都红的同族用例到底为什么红。**
本腿只证了它"两发都红 ⇒ 不在差集、不能用来支持真回归"，**没查它的红因**（那属另一枚票的账，且与本枚结论正交）。
谁能判：票 223 的 owner；或另派一枚只读腿专查它。

---

## 附：本腿对"CI 才红"这一形解释到什么程度算数

一句话：**解释到"反命题已死 ＋ 候选收敛到时序"这一层，算数；到不了"因"那一层。**
- **已死**（盘上日志自证，§4.1，不需再跑一发）：卡没渲染 / 渲染路径被区间内某枚产码改坏。这条直接取消了"真回归嫌疑"这个标签。
- **收敛**（§3）：剩下只有"那次未加界的 stdout 单读没赶上"。它由两条已存在的事实供养——同文件 `:129-139` 早已登记同一形状并给了 1/25、1/5、1/15 的实测频率；以及一枚**代码零改动却同发翻红的纯挂钟钉**（§4.2）。
- **不算数的那半**（§4.5 / §6 C1）：窗口**被什么撑开**、以及**它在本机与 CI 上各自的翻红频率**。那半必须跑，§5 给了确切尺。

对推送的处置（本腿只是取证腿，**不是决定者**）：本枚**不构成**"新代码把产品改坏了"的证据，也**不构成**"可以划成绿"的证据。要把差集里这枚与 `TestResolvePerCallBudget` 一起登记成"CI 抖动"，至少需要 §5.2 或 §5.5 **任意一发**的读数；在那之前它应按 **"未定性 ＋ 已排除真回归"** 入账，而不是按"回归待修"入账。
