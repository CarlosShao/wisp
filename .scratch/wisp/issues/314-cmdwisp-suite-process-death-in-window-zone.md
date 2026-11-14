# 票 314 — `cmd/wisp` 测试二进制在真窗测试区**整进程死亡**：4/4 连发复现（16:54–17:17 窗口）；同日早些时候同形 4/4 全活（16:22–16:50）——未定性，归因候选已列

**立票**：2026-10-13 17:1x 编排者（来路＝亲跑票 311 `Q-311-r1b-1` 的 q1＋顺手复跑 q2/q3/q4；票 311 那格只判"守卫一声不响＋四数行在场"，死亡**出那格射程**、另立本票）
**性质**：★**事故票（测试台面稳性）**。与票 310（死掉的二进制留下临时目录）同因不同面：310 管死后残留，本票管**死本身**。⚠ **未定性为 299 回归**——归因候选 4 枚见 §现量末；⛔ 任何结论先过 `AC#0/#1`。

**为什么现在立**：① `ci.yml:621`（`run: bash scripts/wisp-cli-tests.sh`）那一跳就是默认形 ⇒ **下一次推送，CI 会看见这个死亡**；② 4/4 连发复现，且 q1 与 q4 死在**同一枚测试**（`TestAC14GoSideEvalPushReachesThePage`）。

## 现量（逐发；原始日志根 `C:\Users\swq\AppData\Local\wisp-test-gotmp\orch311-q-20261013\logs\`；⛔ 引用前重跑）

| 发 | 调用形 | 读数 | 死点（最后一行 → ⛔ 无下一行） |
|---|---|---|---|
| q1 | `wisp-cli-tests.sh` 默认形（单包＋CI 那份 `-skip`） | 死于 **185.422s**；`runtests` 四数＝`=== RUN=176 --- PASS=126 --- FAIL=0 --- SKIP=0`（**无任何 `--- FAIL` 行**） | `=== RUN TestAC14GoSideEvalPushReachesThePage` → 窗建起（`MARK299 after CreateWindowExW hwnd=0x17950e52`＋数条 wndproc）→ 死 |
| q2 | 同上复跑 | 死于 **187.924s**；四数 `RUN=172 PASS=120 FAIL=2 SKIP=0`（两枚具名红＝`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 2.71s／`TestAC4FocusReturnToPriorWindowGap33r5` 0.35s） | `TestAC13BringUpSurvivesAReusedThreadQuit` → `after CreateWindowExW hwnd=0x2cd0ed4`＋`msg=0x86` → 死 |
| q3 | **直发**（`go test -count=1 -v ./cmd/wisp/ -skip …`；⛔ 走 wrapper） | 死于 **211.592s**；我 grep＝`RUN=173 PASS=170`＋1 枚具名红（`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` 41.24s） | `TestAC13BringUpSurvivesAReusedThreadQuit` → `after CreateWindowExW hwnd=0xeeb0382` → 死（更早一帧） |
| q4 | **两包**（`./cmd/wisp/ ./internal/audio/`，无 `-skip`＝与 16:22 活的那 4 发**同形**） | 死于 **173.961s**（cmd/wisp）；`ok internal/audio 16.162s` 照常 | `TestAC14GoSideEvalPushReachesThePage` →（同 q1） |
| **对照**：16:22–16:50 的 4 发 | 同 q4 形 | **4/4 全活**（cmd/wisp 各 397–407s；恰两枚真窗红＝软形态） | — |
| q5 | 形②（`WISP_TEST_GOTMP_ROOT` 指同卷换目录）＋**user-data 夹先挪走**（fresh profile，跑中重生成） | 死于 **178.028s**；守卫**静**（GUARD=0／INCONCLUSIVE=0）、**换目录被认领**（`GOTMPDIR=…/orch311-q-20261013/root2/cli-20261013-171826-1438` 逐字）、四数行在（`RUN=172 PASS=121 FAIL=1 SKIP=0`；1 枚具名红＝`TestTicket223ModeLoosening…` 2.27s） | `TestAC13BringUpSurvivesAReusedThreadQuit` →（同 q2/q3） |

- ★**调用形与共享夹的嫌疑全排除**：q3 干掉 wrapper（直发照样死）、q4 干掉"单包/`-skip`"差异（同形照样死）、**q5 干掉"共享 user-data 夹状态"主嫌**（夹先挪走、跑中重生成，照样死）⇒ 当前候选（按嫌疑排序）：
  (a) ★**桌面/会话状态**〔现最强〕：**边界证据＝存活窗口 16:22–16:50 与机主在旁/交互重叠；死亡窗口 5/5 起于他 16:50 最后一条消息之后（此后无交互）**——候选机制＝锁屏/空闲桌面对真窗测试的影响（⛔ 未验；下一发便宜实验＝他人在时复跑同形一发）；
  (b) 299 宿主（`panel_control_windows.go`）／新前端页面字节在真窗区的负载（⛔ 未验；死点只落在 2 枚测试＝更像"这两枚附近"而非全测试面）；
  (c) 随机性（5 死 4 活若 p=0.5 是 1/512 ⇒ ⛔ 当纯巧合，但也未到定性线）。
- ⚠ **q5 副产物（正向）**：该发同时是票 311 形②的闭档读数——守卫静＋换目录被认领＋四数行在 ⇒ 形② 判据成立。
- 死法共同特征：**无 Go panic 输出、无 `--- FAIL`（q1/q4 一枚都⛔ 有）、无 `exit status` 行**——整进程消失；死点全落在 `CreateWindowExW` 返回到首几条 wndproc 之间。
- 前例与对照：`A773`（10-09）**同脚本跑完**（`PASS=260/FAIL=4`）；`pending-and-issues.md:14768` 那族"真窗族不是确定性基线"是**软形态**（红绿漂），本票是**硬形态**（整进程走）。
- 关联票：票 310（死的后果＝临时目录残留；本票是死本身）；票 303（页面→Go 回执断口，已 `-done`）。

## 要建什么（草案；落地腿定最终形）

- [ ] **AC#0 死法取证**：`go test -c` 编出 cmd/wisp 测试二进制 → **直跑它**（`-test.v -test.skip='…'`），拿**进程退出码**（native 异常码 vs `os.Exit(n)` vs `-test.timeout` panic）＋stderr 全量；⛔ 只看 wrapper 汇总。
- [ ] **AC#1 最小重现＋归因 (a)**：q5 A/B 式（fresh profile ⇄ 旧 profile 各 ≥2 发）；判"(a) 共享夹状态"是否充分；若成立 ⇒ 修法候选＝夹具改用每发独立 profile（先查它为什么是固定名：`:61` 处注释逐字）。
- [ ] **AC#2 归因 (b)/(d)**：安静桌面（⛔ 人不用机器／⛔ 挂窗工具退场）下 1-2 发；以及 `main` 对照树 1-2 发（⛔ 动共享工作树；导出树跑）。⚠ 涉及"机主别用机器"的窗口＝要编排者约，⛔ 腿自约。
- [ ] **AC#3 CI 暴露与推送前门**：现读 `ci.yml:601-621` 一带确认哪一跳跑它＋**推送前**按本票形状重跑并把颜色写进推送前检查（⛔ 没跑过不许推）。⚠ 首推会让 CI 见生死——先在本机把死法定性再推。
- [ ] **AC#4 处置**：修（谁）／登记（多长）；⛔ 为变绿放宽任何断言、⛔ 动阈值／golden。

## 排程与边界

- 本票的"跑"（`AC#1/#2`）要 `cmd/wisp` 整包色＋静桌面 ⇒ ⛔ 与 `294-r1`·任何在飞整包腿同波。
- 写面（未来落地腿）：`probes/314/**`＋必要时 `cmd/wisp/*_test.go`（夹具 profile 名一处）；⛔ 本票⛔ 动 `.github/**` 一字。
- **不阻塞机主那条线**（他的产品面正常；死的是测试台面）⇒ ⛔ 拿它去打扰他；推送计划受影响（见 `AC#3`）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）

- [2026-10-13 17:22:56 +08] agent=编排者 did=立票即连发 5 枚（q1–q5）**全部整进程死亡**（174.0–211.6s；死点仅 2 枚测试：3×`TestAC13BringUpSurvivesAReusedThreadQuit`＋2×`TestAC14GoSideEvalPushReachesThePage`）；排除面＝wrapper（q3 直发照死）／单包与 `-skip`（q4 两包无 skip 照死）／共享 user-data 夹（q5 挪夹重生成照死）；对照＝16:22–16:50 同形 4/4 全活 ⇒ 边界＝机主 16:50 最后交互前后（候选 (a) 桌面会话状态，现最强、未验）next=机主在场时复跑同形一发（便宜实验）；`AC#0`（直跑测试二进制取退出码）随下一枚能占 `cmd/wisp` 色的腿走；原始日志 `orch311-q-20261013/logs/{qshot-default,qshot-default-2,qshot-direct,q4-two-package,qshot-shape2}.txt`（⛔ 入仓，同根另存）

- [2026-10-13 21:23 +08] agent=编排者 did=★"机主在场"对照数据入档：`294-r1b` 的 prescribed-form 整包 **5/5 全死**（157.2–174.4s；死点 4 枚具名＝`TestAC14GoSideEvalPushReachesThePage`／`TestAC13BringUpSurvivesAReusedThreadQuit`×2／`TestPanelThreadIsSTAAndExitsCleanly`；覆盖边界＝`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`；r1b 件 `probes/294/r1b/**`），且机主 20:4x–21:1x **全程在机器旁** ⇒ "桌面会话状态"候选不因在场消解 next=等 314 定/修后按 AC#5 重跑同形 2+2

- [2026-10-15 16:2x +08] agent=编排者 did=**★本票的"CI 暴露"来了（AC#3 预告的那一发）**：主 CI 的 `test-windows/cmd/wisp` 步**连续两发跑到窗口区整进程死亡**，与旧台账同签名——`MARK299 before/after CreateWindowExW`＋`MARK299 wndproc …` 后直接 `FAIL … <秒>s`、**该枚用例无 verdict**。逐发读数：
  - run **`38690995413`**（sha `0d1a0d3e`）：四数 `RUN=182 PASS=130 FAIL=2 SKIP=0`；**死在 `TestAC13BringUpSurvivesAReusedThreadQuit`**（起跑后仅 7 行 MARK299，然后 `228.091s`＋`FAIL`）——与旧账"3× AC13BringUp"同枚。
  - run **`38687919508`**（sha `1de62995`）：四数 `RUN=193 PASS=142 FAIL=2 SKIP=0`；死在 `TestLegacySubShapeOneHookDiesInAReentryLoop`（`panel_transport_35r2` 族）之后、下一枚无 verdict。
  - **对照三发（同一步、同 skip 名单、同 runner 标签 `windows-latest`，sha 早于两枚死发的树）全部跑完**：`38582016056`＝`414/301/2/1`、`38673468762`＝`415/302/2/1`、`38676514264`＝`415/301/3/1`（收尾到 `TestZZ315A1*`）。
  - **本机对照**（编排者现跑，`1de62995` 仓外快照、同 skip 名单）＝**`RUN=413`、跑到 `TestZZ315A1*`、`FAIL … 338.989s`**（19 枚红全属本机已知噪声族）⇒ **树面在本机能跑完整套**；⇒ 判**仍是本票那枚未定性的整进程死亡在 CI 上的实例**（不是"少跑/分母缩小"）。
- ⚠ **给本票 AC#3 的两条现读**：① `ci.yml:601-621` 那一跳今天确实"会看见生死"（连续两发死在窗口区）；② 因此**该步的四数不是覆盖面分母**——编排者已立新口径（`A976`）：读 `cmd/wisp` 四数前**先证该趟跑到终点锚 `TestZZ315A1*`**，没跑到就只记"该趟早停、四数不作分母"。⚠ 两枚死发的树都含 `319-r2` 新加的两枚 319 测试，但**本机同树 413 枚跑完**＋死点两发各不相同（AC13BringUp ／ panel_transport 族）⇒ 编排者**暂无证据**把这两发归给 319 新用例；待下一批推送后再取 1–2 发读色（若连发都在窗口区死，再考虑"窗口区用例枚数增长"这条候选）。
- ⛔ 本笔未跑任何本票 AC 的"整包重跑"（`cmd/wisp` 色本刻被 `319-r3/r4` 的腿占；且本票 AC#1/#2 要求不在编队跑动时采读数）；⛔ 未动 `.github/**`；⛔ 未碰任何测试件。

- [2026-11-08 15:3x +08] agent=408落地腿 did=★分界立票：2026-11-08 整包超时族（329 腿 1802s panic timeout／220 腿三发被掐）已定性另一根因＝所有进程共享一个 WebView2 UDF（orbDataPath 固定 `%LocalAppData%\wisp\orb`，多进程自挤 ⇒ S_OK+NULL env 崩 ⇒ recover 吞 ⇒ STA 死 ⇒ ball.New 无限等），修法（UDF 按进程唯一化）与读数全部落在新票 **408**（含残余共享态清单，面板测试件固定名 profile 目录是下一嫌疑）——本票 314 的"整进程无 panic 死亡"形状与那族**不同根**，AC#0（直跑测试二进制取退出码）继续待认领 next=408 落地读数出来后，若整包仍早死则比对死点是否落回本票签名

- [2026-11-12 14:5x +08] agent=编排者（非实现者；实跑取证，⛔ 不翻任何 AC 框）did=**★根因已由票 408 修掉，本票"进程消失"死法不再复现**。①现读：票 408 已 -done，落地产码 `e7f420a6c`「WebView2 UDF 按进程唯一化——`orbDataPath` 从共享 `%LocalAppData%/wisp/orb` 改为 `<base>/wisp/orb/<pid>-<8hex>`（sync.Once 进程内稳定）」，其自述理由逐字＝「WebView2 多实例官方形状是一进程一 UDF，共享夹使 `go test ./cmd/wisp` 整包自撞（S_OK+NULL env 崩/recover 吞/STA 死/ball.New 无限等=整包超时）」——与本票 AC#1 要判的"共享夹状态"是同一根因。②**实跑**（无活 wisp.exe、sherpa PATH、unset GOTMPDIR/TMPDIR）：`go test ./cmd/wisp -count=1 -timeout 30m` 跑到 **40 分钟仍在推进**，`wisp.test.exe` CPU 持续增长（15s delta 0.59）、**零 `--- FAIL`、零 panic、进程没消失**。③**新诊断**：整包慢的主因是**包内变异测试自起嵌套 `go test` 子进程**（实测 PID 33564 `-test.run=^Test219AC2RevokeRoundTripByteForByte$ -test.timeout=20m`，父 PID 28716 `go test ... -overlay=...overlay-m1.json`），即"整包时长≈各变异重投之和"，与 408 修前的 UDF 自撞挂死是两种现象。④⇒ **AC#0/AC#1 的前提已变**：应按"根因已修、只需复验＋登记"处置；AC#2 的"安静桌面"那一发、AC#3 的推送前门仍待办。⑤我 kill 了自己起的整包 job 与残留进程（0 枚残留、CPU 17%），未删文件、未动产码。**本节不翻任何框。**

- [2026-11-13 21:25:43 +08] agent=314 测量/定性腿（⛔ 非写腿）did=**★在 HEAD `397cd9e9d`（派单写的 `02939b8a2` 是其父提交，差 1 枚）走 carrier 实测 6 发**，证据件＝`docs/evidence/s1/314-window-zone-death.md`（新），读数根＝`.scratch/wisp/probes/314/**`。**逐发**：①**run1 全量**＝**卡死**（不是死亡）：卡在第 11 枚 `Test219AC2RevokeRoundTripByteForByte`，进程活着、**CPU 30 秒增量 0.0156 s**、日志 8m31s 零增长，我 kill；②**run2 全量**（自加 `-test.timeout=10m`）＝**rc=2、604s**、四数 `RUN=380 PASS=295 FAIL=2 SKIP=0`、终点锚 `TestZZ315A1` **0 命中**、卡点**具名**＝`TestRunBooksWithASnapshotOfItsLiveQueue`（`panel_pump_test.go:177` 的 `<-done`，等一张永不回答的 L2 卡片；同 dump 里还有一枚 3 分钟前的 `lifecycle_43` owner goroutine 泄漏）；③**死点 5 枚单跑 3 发**（iso1/2/3）＝**3/3 全 PASS、rc=0**、9/18/13s；④**`Test219AC2` 单跑**＝5s PASS ⇒ **卡死是顺序依赖**（单跑活、整包卡）。**桌面状态**（候选 (a) 代理变量）：每发前后 `input_desktop=Default`（**未锁屏**）、session 1 `LogonUI`=0（2/3/4 是断开 RDP 会话、各 1 枚，与本机桌面无关）、前台＝微信、CPU 67–100% ⇒ **(a) 的机制本轮不在场、既不能坐实也不能证伪**；**`WISP_*` 真窗开关全仓零命中**。**★rc=139 判定**：**是同一类现象**——产品码 `panel_host_windows.go:581-586` 逐字把 `rc=139` 等同「take the whole process down」，`33-panel-host-c27-r7.md:99` 记「SIGSEGV，rc=139 ⇒ 整枚测试二进制没掉」，算术 `139=128+11=SIGSEGV`、Windows 未处理 `0xC0000005` 经 MSYS2 翻译成 SIGSEGV；**但本轮⛔ 没亲眼见 139**（两发全量都是「卡」不是「死」），也⛔ 不能证明票面 2026-10-13 那 5 发当时的退出码就是 139。**裁定建议**：定性尚不清且「死法」已扩成两种形态（无声消失／卡到超时 panic `rc=2`）⇒ **继续测量**，缺 ①同桌面连续 ≥10 发全量的退出码＋四数＋终点锚＋卡点栈 ②死亡真发生那一发的退出码原值 ③**锁屏桌面**下 1–2 发（要机主按 Win+L）④跑到终点锚的读数；最小修⛔ 不建议在定性前动 `cmd/wisp/**` 夹具，先按 `314-ac3 §3` 给 `scripts/portable-tests.sh` cli 分支加终点锚断言（脚本面）。门禁 `sh scripts/d22scan.sh`＝**clean**。⛔ 未翻任何 AC 框、⛔ 未 push、⛔ 零产码改动。next=①②④ 由腿再跑一轮（每发 6–20 分钟）；③ 等机主锁屏一句话。

- [2026-11-13 21:46:04 +08] agent=314 只读诊断腿（⛔ 非写腿、⛔ 未翻任何 AC 框）did=**五问取证落 `docs/evidence/s1/314-hang-rootcause.md`（新）＋追加 `.scratch/wisp/probes/314/logs/repro-run1-subset.log`**。**Q1**：`done` 声明 `panel_pump_test.go:155`，关闭点 `panel_pump_test.go:159`（`defer close(done)`，goroutine 起于 :158），关闭条件＝`executeOn145`（:163）返回；永不关的分支列全（`queue.go:159` 的 `q.mu`／`queue.go:359` 的 `mintGrant`／`run.go:2739` `consoleApprovalUI.Prompt` 写 stdout 阻塞／`gate.go:663-703` select 四支）。★**诚实边界**：本例 `ctx` 有 45s 上界（`panel_pump_test.go:112`）、队列 31s（:136），故 `panel_pump` 的 `<-done`（:177）是**有界等待**，不是无限。**Q2**：链＝`gate.go:628`入队→`gate.go:657`上屏（`run.go:2739`）→**没人答**（`run.go:1458` 只在 `s.reply!=nil` 时挂监听，而 `run_test.go:89` 的 fixture 从不设 `reply`）→`gate.go:692` 到期自动拒→`panel_pump_test.go:159` 关。**★顺序依赖候选机制判「无据」**：gate/queue/ui 逐 runtime 私建（`run.go:1052-1054`），`approval` 包无包级队列，「别的用例留一张未答复卡占住队列」结构上不成立；真共享只有 `observe.Default`、`risk.Logf`（`blacklist.go:54`，每次装配被 `run.go:1050` 覆写）等，均不阻塞。**Q3**：`lifecycle_43` owner ＝**后果**——`lifecycle_43_windows_test.go:133` `Spawn(..., nil, …)`，`goroutine.go:288-292` root=nil ⇒ ctx=Background 永不 Done，fixture 的 `t.Cleanup(cancel)`（:100-101）cancel 不到它；dump（`run2.log:3147-3153`）显示它停在 select、不持锁，故不可能阻塞别人。**Q4**：卡点漂移**不是**单调共享状态——本腿复跑 run1 前 11 枚 **11/11 PASS、25s**（`logs/repro-run1-subset.log`），run2 的「卡点」是被 10m 预算在 29s 处砍到（`run2.log:3123,3125`），离 31s 到期只差约 2s；真机制＝**竞态 ＋ 嵌套变异用例抬长整包**（`run2.log:241-242` `Test219AC2Mutation…` 39.70s，`approval_always_219ac2_test.go:379-382` 自起 `go test`）。**Q5**：(a) 隔离＝中代价、只清 goroutine 累积；(b) 裸 `<-ch` 加有界超时＝小代价、直接消灭「静默卡死」形态；(c) cli 分支终点锚断言＝最小代价、只让 CI 响亮变红。**推荐：先 (c)（脚本面、零语义风险）再 (b)（测试面、把裸接收全部加上界），(a) 排后**。门禁 `sh scripts/d22scan.sh`＝**clean（rc=0）**。⛔ 零代码/测试/脚本改动、⛔ 未跑全量（仅 11 枚子集）、⛔ 未 push。next=把 (c) 的终点锚断言交给 `scripts/portable-tests.sh` 写腿；(b) 的裸 `<-ch` 清单（`panel_pump_test.go:177`、`approval_always_219ac2_test.go:242-243`）交 `cmd/wisp/**` 写腿。

- [2026-11-13 23:14:23 +08] agent=314b 写腿（cmd/wisp 面）did=**建议 (b) 落地**：两处裸接收改有界等待（`panel_pump_test.go:177` `<-done` → `waitDoneBounded314` 46s＝45s Execute ctx `:155`＋1s 边距；`approval_always_219ac2_test.go` 三枚用例的 `<-w.done`/`<-w.reply` → `waitRoundTripBounded314` 46s＝40s L2 窗 `approval_reply_201_test.go:154`＋6s 边距，reply 与 done 各一次有界接收共用同一 deadline）；超时 Fatalf 四要素全点名（等什么/等了多久/上界出处行号/当时可判读数——pump 报 pending 张数×2＋实显卡数＋发布数，219 报 done 闭没闭＋WIDEN-APPLIED 是否入账）。★正控：`-overlay` 变异（上界 50ms＋219 注入 reply 不关，变异件全在 `D:\tmp\wisp\20261113\314b\`，仓零触碰零还原）⇒ 两枚都红且失败语逐字点名（读数在证据件）；还原（nothing to restore）⇒ 真 carrier 定向 4 枚 **RUN=4 PASS=4 FAIL=0 SKIP=0**、rc=0。顺带同轮干了新红的 16r1 premortem 定因（判 (a)，修 `voice_e2e_demo_16r1_windows_test.go` 的 5s 墙钟猜想→mtime 前提）——全文见 `docs/evidence/s1/314-bounded-waits-and-16r1.md`（新）与票 111 progress log 同刻那行。门禁 build/vet/winlive-vet/d22scan/gofumpt 全绿。⛔ 未翻 AC 框、未 push、零产码/脚本/internal/models 改动。next=AC#0/AC#1 复验（整包色）仍待不被 portable 腿占用的窗口。

- [2026-11-14 09:18:15 +08] agent=314c 写腿（脚本面；⛔ 非实现者）did=**建议 (c) 落地**：`scripts/portable-tests.sh` cli 分支加终点锚断言（+131 行／0 删）。跑前 `go test -list '.*' ./cmd/wisp/` 数出分母 **619** 枚、取最后枚举名 `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` 当锚；跑后必须出现该枚的 `--- PASS|FAIL|SKIP:` verdict 行，缺则响亮红并点名**跑了几枚／分母几枚／缺的锚／最后一行**，且 rc=124/137、128+n 信号（含 139）、`panic: test timed out` 归类为「被砍」而非「测试失败」。★正控真造早停：`-test.run` 只跑 1 枚 ⇒ carrier 自己打 `ok` 行、GUARD B 也放行，但断言红 rc=1（`RUN=1` vs 分母 619、锚缺失、末行 `ok ... 3.00s`）；拆掉断言同一发 rc=0 全绿。★负控本机两发真实全量**都给不出绿**（第一发 rc=139 死在 `Test145ArgsJSONComesFromTheBookedForensicsRow`、RUN=418/370s；第二发跑到 RUN=795 撞 90m `-test.timeout` panic rc=2、36149s 含挂起）——两发都被新断言当场点名，其中 panic 那发是最终码的真实读数（`KILLED, NOT FAILED (panic: test timed out)`）；另做合成完成流证明锚出现时断言静默 rc=0。门禁三条终态：`bash -n` rc=0／`sh scripts/d22scan.sh` clean／`sh scripts/check-path-length-budget.sh` GREEN。证据件 `docs/evidence/s1/314-end-anchor.md`（新）。⛔ 未翻 AC 框、⛔ 未 push、零 Go/阈值/golden/ledger/carrier 改动。next=本机全量绿仍欠（产品面 ticket 314 本人），交下一发 CI 第 7 步四数＋锚判色。

## Progress log（编排者 2026-11-14 09:5x）

**独占干净全量（单发、无竞争者）：`rc=1`、1442s、PASS=786 / FAIL=1 / SKIP=12，跑到终点锚 `TestZZ315A1`。**
- **没有 `rc=139`、没有挂起、没有 `panic: test timed out`** ⇒ 本票的「整进程死亡」在单发条件下**不复现**。
- 此前两处异常（`rc=139` @370s、90 分钟挂起）都发生在**两发整包并发**的窗口内（编排者自己的测量与一条腿的全量相撞）⇒ 机制是**并发**，与诊断腿 `314-hang-rootcause.md` 的「竞态＋嵌套变异抬长整包」一致。
- 已落地两闸：①终点锚（`832e04489`，`docs/evidence/s1/314-end-anchor.md`）使「没跑完」响亮变红并点名；②并发纪律入 goal：整包同一时刻只许一发。
- 唯一红 `TestAC4FocusReturnToPriorWindowGap33r5`＝真窗前台时序已知红，与本票无关；12 枚 SKIP 全部具名带 `文件:行`。
- **仍未验**：候选 (a) 锁屏桌面（编排者自锁失败，`LogonUI` 前后为 0），待机主 `Win+L` 后补跑。
- 证据件：`docs/evidence/s1/314-clean-single-run.md`。
- [2026-11-14 10:54:04 +08] agent=314 终点锚审计腿（非实现者；⛔ 未翻任何 AC 框）did=独立复核 `832e04489` 的 cli 终点锚。★甲通过：临时 `-test.run` 只跑 1 枚 ⇒ carrier 自己打 `ok ... 3.00s`、GUARD B 打 `ok (own line)`，断言仍红 `rc=1` 并逐字点名 `=== RUN=1`／`roster=619`／缺锚 `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`／末行 `ok ... 3.00s`；同一发拆掉断言 `rc=0`。★乙通过（小全集，本审计未跑整包）：只跑最后两枚 ⇒ `RUN=2 PASS=2`、锚 `--- PASS` 命中、`rc=0`、无锚判决段；合成矩阵证明 `=== RUN` 单行 NOMATCH、只有 `--- PASS|FAIL|SKIP` 算到达，stub 只打锚的 RUN 仍 `rc=139` 红。★rc 分类通过：124/137/139/143 与 `panic: test timed out` 均 `KILLED, NOT FAILED`；rc=1 无签名 → `EARLY STOP / HANG`（四数保留 FAIL=1）。★未放宽：脚本 +131/0，运行/rc/锚判定无 `|| true`，唯一 rc 写是 0→1；临时改动已整文件备份还原（SHA256 与动前相同、六条 grep 全 0、`git diff` 空）。证据件 `docs/evidence/s1/314-end-anchor-audit.md`（新），读数根 `.scratch/wisp/probes/314-anchor-audit/**`。⛔ 未跑整包、⛔ 未 push、⛔ 未动产码/测试/阈值/golden。next=编排者决定是否据此收 314 的锚判色。

- [2026-11-14 11:23:19 +08] agent=314 测量腿（⛔ 非写腿、⛔ 未翻任何 AC 框）did=**AC#1 归因＋AC#2 A/B 串行实测**，证据件 `docs/evidence/s1/314-profile-ab.md`（新），读数根 `.scratch/wisp/probes/314-profile/**`。★**逐字读 `:61`**：旧版（`68275983e^`）`:61` = `dataPath := filepath.Join(os.TempDir(), "wisp-33r5-panel-profile")`，其上方四行注释逐字讲的是「生产形状／夹在系统 temp 下／cleanup 已注册」——**没有一个字解释"为什么固定名"**；固定名来自 33-r5 首落 `13acad460`，其 commit body 对 profile/user-data 零命中 ⇒ **票面说的"此处注释解释原因"不成立，那里没有原因**。唯一写"故意不删"的是 `cmd/wisp/panel_host_windows_test.go:604-608`（`t.TempDir` 的 RemoveAll 会与浏览器 profile-DB 退场赛跑），那是"为什么不删"、仍不是"为什么固定名"。★**因果链（file:line）**：夹具夹＝WebView2 的 UDF —— `cmd/wisp/panel_host_windows.go:625`（`newPanelControl(m.windowOptions(), m.dataPath)`）→ `cmd/wisp/panel_control_windows.go:355`（`browser.DataPath = dataPath`）→ `:433`（`browser.Embed(hwnd)`）→ `forks/go-webview2/pkg/edge/chromium.go:95`（`createCoreWebView2EnvironmentWithOptions(nil, <dataPath>, …)`）。★**A/B（每发独占、全程串行，15 发）**：A1 固定/旧夹 ×3（rc=0／4.80、4.53、4.60s／4 PASS，同夹文件 175→178→179→182 逐发累积＝真复用）；A2 固定名+清理 ×2（rc=0／4.81、4.79s；`RemoveAll` 只半删 175→53→20）；B 每发独立 ×4（rc=0／8.30、5.49、5.61、5.85s）；紧循环（上一发退出即起下一发）A1t ×3 6.13/8.22/6.57s、Bt ×3 6.20/7.84/8.72s —— **rc 全 0、verdict 全 PASS、panic 0、`boot failed` 0、`WebView2 Embed failed` 0、native 异常 0**，每发 competitor `wisp.test`=0、每发建窗 4 枚。★**同夹并存取证**（`udf-share.ps1`，仍一次一个测试进程）：固定名下**同一 UDF 同一 tick 出现两个顶层浏览器 host** —— parent=已退出的上一发 pid 49456 与 parent=当前发 pid 344，`tick=0..3`（约 1.2s）连续共存，`wisp-33r7-queued-close` tick 6 同形；两发 rc=0（warm 4.89s／sample 8.41s）。对照 B 组 `coexist=none`、`maxHostsOnOneUdf=1`。⇒ **"固定名⇒两发互踩同一夹"＝坐实；"共享夹状态足以解释 139/挂起"＝不充分** —— 缺的正是 408 记录的"两进程**同时**调 `CreateEnvironment` 争同一夹"（`internal/ball/orb_windows.go:436`）或"新 environment **JOIN** 濒死浏览器进程"（`:447-458`）那一瞬，它只在真并发下出现，是本腿边界外。★**时间线反证**：r3 守卫 `da5d559b7`（11-08 20:13，S_OK+nil env ⇒ `bootFailed` ⇒ `cmd/wisp/panel_control_windows.go:433-436` 具名红）＋三处 de-share `68275983e`（11-11 19:22）都早于 11-14 那枚 139 ⇒ 今天这条链的预期形状是具名红，不是整进程 139；要把 139 挂回它需要"崩溃在守卫射程外"的证据，本轮零证据、不主张。★**残余面**：de-share 只动三处，默认档仍有两枚**真建窗**用例共享固定夹 —— `cmd/wisp/panel_host_windows_test.go:609`（`wisp-33r1-panel-profile`，`TestPanelHostRealWindowHopAndLifecycle`）与 `cmd/wisp/panel_assets_299_windows_test.go:161`（`wisp-299r1-panel-profile`，`TestPanelPageFetchesItsOwnAssets299`）；`panel_holder_319_windows_test.go:53` 自述"无窗、当名字用"（`:45-50`），`panel_unbound_253r6_windows_test.go:73` 未判 ⇒ **并发时"互踩"通道未被关死**。门禁 `sh scripts/d22scan.sh`＝**clean（rc=0，85s）**。⛔ 未跑整包（四枚定向，最长 8.72s）、⛔ 未做真并发、⛔ 零产码/既有测试/阈值/golden/契约/`scripts/**` 改动（A 组走 `-overlay`，`cmd/wisp` 那份测试件一字未动、`git status --porcelain cmd/wisp/` 空读数）、⛔ 未 push。三枚测试二进制 62.5MB×3 只在盘上、**不入 commit**（证据件 §6 可现场重编）。next=若编排者要在真并发下闭合这条链，须在独占机器的两发窗口里抓"同时 CreateEnvironment"那一瞬（本腿不碰）；`panel_unbound_253r6:73` 建不建窗待判。
- [2026-11-14 11:54:58 +08] agent=314d 写腿（实现者；⛔ 审计者；⛔ 未翻任何 AC 框）did=**三处固定名面板 WebView2 user-data 夹改每发独立**（`panel_host_windows_test.go:706`／`:1027`（同文件第二处）／`panel_assets_299_windows_test.go:164`），提交 **`26e9e9e73`**（只 commit ⛔ push）。★**形状照 `68275983e`**：新 helper `panelProfileLeaf(prefix)` ＝ 先清同族旧夹再交 `uniqueProfileLeaf(prefix)`（复用 `panel_resident_windows_test.go:61`，⛔ 未另造命名法），三处**都不加 `t.Cleanup(RemoveAll)`**（`panel_host_windows_test.go:701-705` 的原文理由：自动 RemoveAll 与浏览器 profile-DB 退场赛跑会撞锁文件），父目录仍 `os.TempDir()`；"不删"改由**启动时清 >1h 同族旧夹**付（`sweepStalePanelProfiles`：每进程每 prefix 只跑一次，只删名字**恰为** `prefix-<pid>-<nanotime>` 且 mtime >1h 的目录，错误全忽略）。★**并发正控**（`probes/314-deshare/concurrent.ps1`，两发背靠背同时起、采样器每 250ms 读所有 `msedgewebview2` 的 `--user-data-dir`）：**两发各建一窗**时 `deshare-299-r1/r2`、`deshare-ac4-r1` **两发全 rc=0**、`maxTopLevelHostsOnOneUdf=1`、**同夹共存 0 tick**、两发 UDF 路径不同（逐条路径在 `results-concurrent.jsonl`）。★**反证**（同形状，三处用 `-overlay` 临时改回固定名重编，`overlay/` 三件入仓，编译 rc=0）：`fixed-ac4-r2`（A.pid=25976／B.pid=5880）**同一 UDF 连续 4 tick 两个顶层 host、parents 恰是两发 pid**、两发自己也把 UDF 读成同一条；`fixed-ac4-r1` 2 tick、`fixed-299-r1` 1 tick —— **但四发 rc=0，这一形我造不出致命红（逐字记）**。★★**诚实边界**：派单字面的"两发各跑两枚真建窗用例、两发都 rc=0"我**拿不到** —— 该形 **4/4 对里有 1 发原生消失**（`0xC0000005`×3、`0xC000041D`×1；无 `panic:`、无 `--- FAIL`、死在 `=== RUN <第二枚>` 之后），且**固定名构建同形也死一发（同一退出码）** ⇒ 该死亡不随夹名改变，既不是我引入的、也不是 de-share 治得了的，我不归因（派单 ⛔ 再归因）；另 `lifecycle` 自带 D32 冷启 1500ms 预算，并发下实测 1607.9／1641.1／1864.3ms **必红**（单发 1147／1322ms），`AC4` 的前台尺读机器级前台窗也会红（`deshare-ac4-r2` 两发红在同一句、读同一个别人的窗）—— 两者都**不是夹读数**，我没放宽任何断言。★**清理正控**（`sweep-check.ps1`）：种 5 枚死物（旧同族叶／新同族叶＝并发替身／别的 prefix／旧固定名／形近名）→ 结果 `oldLeaf` 删、其余 4 枚全留、`pass=true`。★门禁：`go vet -tags winlive ./cmd/wisp` rc=0；`gofumpt -l cmd/wisp` 空（用 `$(go env GOPATH)/bin/gofumpt.exe` v0.12.0）；`sh scripts/d22scan.sh` **clean rc=0**；真跑定向四数：单发 299＝RUN1/PASS1/FAIL0/SKIP0、单发 lifecycle＝1/1/0/0、单发 AC4＝1/1/0/0、单发两枚＝RUN2/PASS2/FAIL0/SKIP0（4.5–5.8s）。★**还原**：反证走 `-overlay`，两枚跟踪件 **SHA256 改动后前后同串**（`C080903D…88D7E`／`42735AA5…2760B`）、`git diff` 仍只有这 2 件（+113/−5）、无第 3 件；未跑整包、⛔ 零产码／`scripts/**`／`internal/**`／`frontend/**`／阈值／golden／契约改动。证据件 `docs/evidence/s1/314-deshare-fixed-profiles.md`（新），读数根 `probes/314-deshare/**`。next=①这枚"两发各两窗必死一发"是否另立票（我只有读数、无归因）②`lifecycle` 的 D32 冷启预算要不要在并发控制里排除（阈值口径、须人工批，⛔ 我没动）

## 编排者裁定（2026-11-14 12:0x）：de-share 三处落地、缺口补派、两问已拍

### 一、已落地（`26e9e9e73`＋`237f41fcb`＋`edccb3449`，已推）
三处固定名 WebView2 user-data 夹 → `panelProfileLeaf(prefix)`（唯一叶＝pid＋纳秒），照 `68275983e` 形状；不加 `Cleanup(RemoveAll)`（原理由保留），改由「启动时清 >1h 同族旧夹」兜。
- **并发正控**：两发背靠背（|skew|≤83ms）⇒ 两发 `rc=0`、`maxTopLevelHostsOnOneUdf=1`、**同夹共存 0 tick**、两发 UDF 路径不同。
- **反证（overlay 改回固定名）**：`fixed-ac4-r2` **同一 UDF 连续 4 tick 两个顶层 host，parents 恰是两发 pid**；另两对 2 tick／1 tick。
- **实现者如实记**：固定名侧**四发 rc=0、无 panic／boot failed／Embed failed ⇒ 这一形造不出致命红**。

### 二、⚠ 我查出 de-share **没做全**：还剩 5 处固定名夹
`panel_unbound_253r6_windows_test.go:73`（tag **只有 `windows`，进默认跑**）、`panel_host_windows_budget_winlive_test.go:57`、`panel_host_windows_live_test.go:37`、`panel_reshow_255r8_winlive_test.go:125`、`panel_transport_live_35v2_windows_test.go:186`（后四处 `windows && winlive`）。
⇒ **已派 `w314-deshare-c` 补完**（复用既有 `panelProfileLeaf`；若某处故意要固定名，须写理由＋`文件:行`，⛔ 无理由照改）。

### 三、两问拍板（编排者自己拍，⛔ 摆机主）
1. **「两发各跑两枚真建窗用例 ⇒ 4 对里 1 发原生消失（`0xC0000005`×3、`0xC000041D`×1，无 `panic:`、无 `--- FAIL`）」要不要另立票？**
   ⇒ **不另立票，就记在本票**。理由：这正是本票的主题（窗口区整进程死亡），且**固定名构建同形也死一发（同一退出码）** ⇒ 与夹名无关 ⇒ 它是本票的**现成新读数**，不是另一件事。**该腿正确地没有归因**，我也不归因。
2. **`lifecycle` 并发下冷启 1607.9／1641.1／1864.3ms > D32 的 1500ms 预算，要不要在并发控制里排除？**
   ⇒ **属阈值口径（D32），本编队动不得**；且**不做排除**——单发读数是 1147／1322ms（在预算内），并发才超 ⇒ 正确处置是**「并发跑不许用作 D32 预算判据」**写进纪律，⛔ 不是改阈值。**若机主要动 D32 数字，那才是他批的事**（已列进待机主项）。

### 四、审计安排
⇒ 审计腿**押后到 `w314-deshare-c` 交付后**再派：此刻 cmd/wisp 正被写腿改，并发跑读数会脏（我已因并发污染作废过两条全量读数，⛔ 不再犯）。

## Progress log（写腿 `w314-deshare-c`，2026-11-14 12:2x）

- [2026-11-14 12:24:00 +08] agent=314e 写腿 `w314-deshare-c`（实现者；⛔ 审计者；⛔ 未翻任何 AC 框）did=**补完剩 5 处固定名面板 WebView2 user-data 夹的 de-share**，提交 **`b1e1a738f`**（只 commit ⛔ push；其代码提交在 `48457eafe` 之上，落地后数秒并行的票 44 腿提交了 `6f23ce632`，**不是本腿的**）。★**5 处全部照用既有 `panelProfileLeaf(prefix)`，零新 helper**：`panel_unbound_253r6_windows_test.go:71`（**默认档**）、`panel_host_windows_budget_winlive_test.go:55`、`panel_host_windows_live_test.go:36`、`panel_reshow_255r8_winlive_test.go:125`、`panel_transport_live_35v2_windows_test.go:185`；import 只在"被替换行是最后使用点"处收窄（`os`/`path/filepath`），`panel_reshow` 的 seed 写入仍用 `os`/`filepath` 故未动。★**逐处核过"是否故意固定名"⇒ 没有**：4 处上下文注释与夹名无关；唯一相关的是 `panel_reshow_255r8_winlive_test.go:118-123`，它论证的是"根目录别用 `t.TempDir()`"（RemoveAll 与 Edge profile 落盘赛跑），`panelProfileLeaf` 保留 `os.TempDir()` 根 ⇒ **不构成保留固定名的理由**；派单要求的"理由＋文件:行"逐处写在证据件 §1.3。★★**默认档正控**（`probes/314-deshare2/concurrent.ps1`，两发背靠背同时起，采样器每 250ms 读 `msedgewebview2` 的 `--user-data-dir`）：两轮 `deshare-253r6-r1/r2`，skew **−53.2／−77.3 ms**，**两发都 rc=0、各 1 PASS**，`maxTopLevelHostsOnOneUdf=1`、**同夹共存 0 tick**、两发 UDF 路径**不同**且都含 pid＋纳秒叶（逐字路径在 `results-concurrent.jsonl`）。★**反证**：`-overlay` 把该处临时改回固定名重编（5 份副本字节级等于 `48457eafe:cmd/wisp/<file>`，编译 rc=0；跟踪文件一字未动），同形 3 轮里 **2 轮复现同一 UDF 两个顶层 host**（各 1 tick，parents 恰是两发 pid，逐字：`tick=0 udf=…\wisp-253r6-panel-profile\EBWebView topLevelHosts=2 parents=53088,18276 browsers=65520,41256`）；第 3 轮 skew 225.6ms 更大、没抓到共存那一瞬，**如实记"这一轮造不出"**，不把 2/3 说成 3/3；固定名侧两发**都 rc=0**，仍未造出致命红。★门禁：`go vet -tags winlive ./cmd/wisp` **rc=0**；`gofumpt -l cmd/wisp` **空**；`sh scripts/d22scan.sh` **clean rc=0**（首次 12:12 曾红在 `internal/observe/quota_44.go:223` 的 phantom-citation，12:15 该文件由并行票 44 腿落盘后重跑即 clean ⇒ 与本腿无关）；定向五枚最终同轮 **RUN=1/2、PASS=全部、FAIL=0、SKIP=0**（默认 253r6／2s 退出／D32 预算对／255r8 reshow／35v2 transport）。★**负载敏感红如实留档**：本机 CPU 91–100%，2s 退出条款 3 轮 A/B 里 deshare 红 2 绿 1、改前固定名构建红 1 绿 2；D32 预算负载尖峰时 hot 207.7ms 红、同一构建随后 142.1／44.0／84.1ms 绿，且改前构建也读到 cold 1352.7ms 贴近上限；35v2 teardown round2 两个构建同红（tree=8）——**均用改前二进制同形对照证明与夹名无关；⛔ 未放宽任何断言、⛔ 未加 retry、⛔ 未删断言、⛔ 未动阈值/golden/契约**。★**还原核**：反证走 overlay ⇒ 5 个跟踪件 SHA256 同串（`C9A14528…40551`／`8174EC64…C21D`／`71C9F788…94A6D`（略写见证据件）／`96A38B69…33086`／`CA478141…4D7D8`），跟踪测试件里旧固定名 `filepath.Join(os.TempDir(), …)` 现量 **0 条**，`git diff` 只有这 5 件（+5/−11）。证据件 `docs/evidence/s1/314-deshare-remaining-five.md`（新），读数根 `.scratch/wisp/probes/314-deshare2/**`。⛔ 未跑整包、⛔ 零产码／`scripts/**`／`internal/**`／`frontend/**`／阈值／golden／契约改动。next=审计腿可在 cmd/wisp 空窗后开跑（本腿已交付）；另：本腿只复跑"两发各一窗"的正控，**没有**重跑上一腿"两发各两窗必死一发"的形状（不属本腿派单）。

- [2026-11-14 13:2x +08] agent=today-four-fixes-audit（非实现者，审计腿）did=审计 26e9e9e73+b1e1a738f：8/8 目标调用点已走 panelProfileLeaf，tracked 测试只 +5/-11 路径/import、零断言改动；A3 把默认档 :71 改回固定名（-c 编译 rc=0）后 Test253r6 PASS=造不出红（仓内无唯一性断言，ps1 未接门禁）；D 全仓仍有 2 处固定名 WebView2 DataPath（319 live :75 真建窗、319 holder :53 无窗）；结论=部分；AC#4 不能翻（处置未清、AC#2 静桌面缺）。证据 docs/evidence/s1/today-four-fixes-audit.md

## 非实现者审计回执（2026-11-14 13:2x，`b70ea907f`）：四笔 **成立/成立/成立/部分**

审计件：`docs/evidence/s1/today-four-fixes-audit.md`（原始读数 `.scratch/wisp/probes/audit/logs/**`）。**A 类突变全在 `C:\Temp\wisp-audit-20261114` 导出树做、能编译、跑完还原。**

### 本票（314）相关：**停在「部分」，AC 不翻**
审计给出三条硬读数：
1. **A3：这一形造不出红** —— 把默认档 `panel_unbound_253r6_windows_test.go:71` 改回固定名，`go test -c` rc=0、`Test253r6…` **PASS**。**仓里没有任何断言读「叶名唯一性」**，唯一控制是 `.scratch/wisp/probes/314-deshare/concurrent.ps1`，而它**没有接进任何门禁**。
2. **D③：固定名未清零** —— 8/8 目标调用点已每发独立，但全仓仍有 2 处固定名 WebView2 DataPath：`panel_death_replay_319_winlive_test.go:75`（**真建窗**）、`panel_holder_319_windows_test.go:53`（**无窗**，其注释自述故意共享）；另有 spike 默认名 1 处。
3. **AC#2 静桌面读数仍缺**（机主一直在用机器 ⇒ 那一发还没跑）。

⇒ **处置（我拍）**：
- 314 的 AC#4 **不翻**，理由写清（上三条）；撤销口令＝「314 AC4 复议」。
- **下一步两件**（待 `cmd/wisp` 写窗空出后派）：① de-share `panel_death_replay_319_winlive_test.go:75` 那处**真建窗**；② 在仓内加一枚**源码形状尺**（Go 测试，扫 `cmd/wisp/*_test.go` 里 `filepath.Join(os.TempDir(), "…-panel-profile")` 形状，命中即红，每处豁免须在测试内具名带理由）—— **这样 A3 那种回退才会红**，而不是靠一个没人跑的 ps1。⛔ 不走「改 `ci.yml`」那条路，Go 测试天然进 CI。
`panel_holder_319_windows_test.go:53` 按它自己的注释**保留**（无窗、故意共享），但要在形状尺的豁免表里具名。

### 另三笔的审计结论（供票 44／197 用）
- `6f23ce632` 与 `756ebbcfa`+`f8e94bbd8`（票 44 配额锁）：**成立**。A1 去掉 `Spend` 锁 ⇒ `-race` rc=1、**25 条 DATA RACE ＋ 1 条 fatal**；A1b 去掉 `SeedFrom` 锁 ⇒ **8 条 DATA RACE**。B 类无放宽。
- `d8e2f4dc8`+`0c9ff7c0c`（票 197 发布窗口）：**成立**，但**具名弱点**：把一处 `PublishSubagentInFlight` 改回裸发布时，**形状尺红而症状尺仍 PASS** ⇒ 症状尺对回退没有确定性牙；另外 `Record`/`PublishSubagent`/`AttachCancel` 仍是公开 API，形状尺只扫 `subagent_197.go`，**新文件可绕开**。

### 一处小账（不另立票）
`panel_host_windows_test.go:706` 在证据件里写作 `:711`（行号漂移，非功能问题）。

- [2026-11-14 13:49:06 +08] agent=314f 写腿（实现者；⛔ 审计者；⛔ 未翻任何 AC 框）did=**上一节裁定的两件都落盘**（代码、证据件与本行同在本提交；⛔ 未 push）。①`panel_death_replay_319_winlive_test.go` 那处**真建窗**的固定名夹改成 `panelProfileLeaf("wisp-319r2-live-profile")`（helper 既有 `panel_host_windows_test.go:70`，零新 helper；`path/filepath` 是该文件最后使用点故同步收窄 import）；验证＝`GOFLAGS='-tags=winlive' go test -c` **rc=0**（63,072,718 B）；**未跑那枚真窗用例**——本刻桌面被另一条腿占着，不抢别人的窗口换读数。②新件 `cmd/wisp/panel_profile_leaf_shape_windows_test.go`（`//go:build windows`、纯 go/parser、**不建窗**、0.82 s）：扫 `cmd/wisp/*.go` **273** 件 ＋ 具名的 `scripts/spike/webview2-latency/main.go` ＝ **274** 件；两条规则——**NAME**（`Join` 里 `os.TempDir()`／其局部或包级别名 ＋ 含 `profile` 的字符串字面量叶，参数顺序/跨行/`+` 串接都算，不看流向）与 **SINK**（任意叶一旦进 `NewPanelManager`／`newPanelControl`／`DataPath` 字段赋值／`DataPath:` 键，直传或经一次局部绑定）⇒ 命中即红，红句逐字点名 `文件:行` 并写明「这个形状会让两发并发互踩同一个 WebView2 user-data 夹」。豁免 **2 条**，全在测试件内部具名：`panel_holder_319_windows_test.go`（无窗、故意共享）与 spike 的 `-datapath` 默认值；每条带**从源码逐字抄下的理由句**，并每次机器复核「理由句仍在那个文件里」「豁免仍覆盖 ≥1 处活命中（过期即红）」「`noWindow` 的文件不得出现 `bringUp`/`newPanelControl`/`startPanelHolder319` 调用」。★**正控（派单必须，已做）**：把审计 A3 同一处 `panel_unbound_253r6_windows_test.go:71` 改回 `filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")`（加回 `os`＋`path/filepath` 两行 import ⇒ 该行落 `:73`），突变体先编译再跑 ⇒ **rc=1**、红句逐字点名 `cmd/wisp/panel_unbound_253r6_windows_test.go:73`、读数 274 源文件/3 命中/2 豁免/**1 红**；还原＝突变前整文件备份原字节拷回，**SHA256 同串 `C9A14528D6E64CE98E13EDEF5F1064585B975234DF22BD384E78ADCAFB240551`**（与审计 §0.1 同串）、`git status`/`git diff` 该文件空、再跑 **rc=0 绿**（274/2 命中/2 豁免/0 红）。尺子另带**常驻合成正控 14 例**（9 必红＋5 必静）每跑一次，机制本身不会悄悄变瞎。★**一处设计取舍（防「把豁免表写宽」）**：初版规则「任何 `os.TempDir()` ＋字面量」第一次跑就误中 `panel_models_320_test.go:580` 的 `planted-key-320-temp.txt`——那是票 320 泄漏检测的**阳性控制文件**、不是 WebView2 夹 ⇒ 收紧到票面主体（NAME/SINK 两规则），320 那处**天然静默、零豁免**，没有为它加一条宽豁免。门禁终态：`go vet -tags winlive ./cmd/wisp` **rc=0**；`gofumpt -l cmd/wisp` **空**；`sh scripts/d22scan.sh` **clean rc=0**（`ban #8 cmd/=280`）；`sh scripts/check-path-length-budget.sh` **GREEN**（tracked=19,856；over-budget=69 全在名册；not in roster=0；wall interval 内 0）。证据件 `docs/evidence/s1/314-shape-ruler.md`（新），读数根 `.scratch/wisp/probes/314-ruler/**`（突变红日志、还原前后哈希、d22/路径预算日志；winlive 测试二进制被 `*.exe` 忽略、不入库）。⛔ 未改产码/`scripts/**`/`internal/**`/`frontend/**`/`.github/**`/阈值/golden/契约；⛔ 未跑整包；⛔ 未 push。next=交非实现者审计这枚尺子（重点攻三条射程外：固定名挪进 helper 形参、两个半截字面量拼叶、非 `os.TempDir()` 的根）；314 AC#4 的「清零」仍差 holder 那处无窗共享（本件按裁定保留并具名豁免）与 AC#2 的静桌面读数。

---

## 非实现者审计（源码形状尺）：3 条射程外攻击 ＋ 2 件复核（2026-11-14 14:2x）

审计件 **`docs/evidence/s1/314-ruler-audit.md`（新）**，读数根 `.scratch/wisp/probes/ruler-audit/**`（13 份日志、4 枚探针源、2 份原件备份）。
审计对象＝`cmd/wisp/panel_profile_leaf_shape_windows_test.go`（SHA256 `7D0CB395…`，blob `cb308abd`，提交 `0b28c7d98`，此后零改动）。
树＝`git archive` 的 `50285dff1` 导出树 `C:\Temp\wisp-ruler-audit-20261114`（审计期间仓 HEAD 前进到 `069f97cac`，是别的腿的票 245，未动这把尺）。⛔ 仓内零改动（写面只有本件、审计件、`probes/ruler-audit/**`）；⛔ 未跑整包（只 `-test.run '^Test314'`）；⛔ 未 push。**本件不翻任何 AC 框。**

**五节结论：攻击 1 抓住／攻击 2 漏了／攻击 3 漏了／豁免成立且过期即红／14 例计数对得上（9＋5）。**

- **基线**：导出树现跑 `274 源文件 / 2 命中 / 2 具名豁免 / 0 红`、rc=0；14 个正控子例全 PASS（`-test.v` 逐条数：9 必红＋5 必静）。
- **攻击 1（helper 返回值 → DataPath）抓住，但不是数据流**：命中的是 `:375` `strings.Contains(strings.ToLower(fixed.leaf), "profile")`——叶里写着 `profile` 就红（探针 A `:19` 红句逐字点名）。**同一形状把叶改成不含 `profile` 即静默**（A1b）；头注释自交底的「根进形参」（A1c）静默；另测出**第四条**它没交底的静默：包级 `DataPath:` 键 / 包级固定 Join 变量（SINK 只走 `fn.Body` `:353-365`，包级只走 NAME `:344-348`；包级 Join 叶改含 `profile` 立刻红，证明护栏是拼写不是作用域）。
- **攻击 2（半截字面量）全漏**：`"wisp-33r1"+"-panel-profile"`、整条 `+` 链、`strings.Join`、Join 外再加尾巴、拼完先绑局部——五形全静默。因 `:455` 只认单枚 `*ast.BasicLit`，`:469-478` 的 `+` 分支要求一边字面量一边 temp root（`:473`/`:476`），`:429` 不递归实参子表达式。
- **攻击 3（非 os.TempDir 根）全漏**：`os.Getenv("TEMP")`、`os.UserHomeDir()`、绝对路径字面量三形静默——`panelProfileIsTempRoot`（`:559-574`）只认 `os.TempDir()`/其别名/括号；`A3a/A3b` 的叶里**写着 `profile`** 也救不回来。
- **复核 4（豁免表）成立**：2 条豁免的理由句逐字在源码（holder `:46`/`:48`、spike `:307`），各覆盖 1 处真命中（现跑 2 命中＝2 豁免，无空转）；身份是 `(file, leaf)` 精确匹配、无 glob。四发突变（末一条是加测）：删理由句→`:190` 红；改叶名→`:190` 过期红＋新叶作固定名红；改 spike 的 `flag.String("datapath"`→`:190` 红；加测「在 holder 里调 `bringUp`」→`:190` noWindow 红。（一处口径：这些红下汇总行仍印 `0 处红`，因该计数器只数未豁免的形状命中——别拿它当绿证。）
- **复核 5（14 例）计数对得上**：用例表 `:643-841` 逐条数出 9 必有命中＋5 必静，现跑 14 `=== RUN`/14 `PASS`；但**四条 SINK 臂只覆盖两条**（`NewPanelManager` 直传/局部绑定），`newPanelControl`、`DataPath` 字段赋值、函数体内 `DataPath:` 键**零常驻正控**（探针 C3/C4 证明后两臂当前是好的）。
- **AC#4 判读**：**不能单独翻**。它把「A3 回退无牙」变成默认套件里的确定性红（且 `ci.yml:1165` → `wisp-cli-tests.sh` → `portable-tests.sh --scope=cli` 链无 `-test.run` 过滤，确认在 windows 腿分母里）；但残余两处固定名是**登记**（holder `:53` 无窗共享、spike `:307` 手工工具）不是清零，换根/拆字面量即全盲，AC#2 静桌面读数仍缺。
- 撤销口令＝「314 形状尺审计不采信」：删本节的结论行与 `docs/evidence/s1/314-ruler-audit.md`，读数根可原样留着复跑。
## Progress log（写腿 `w314-ruler2`，2026-11-14 14:5x）：审计 D1–D5 全部修好；A1b/A1c 具名仍漏

- [2026-11-14 14:5x +08] agent=w314-ruler2 写腿（实现者；⛔ 审计者；⛔ 未翻任何 AC 框）did=**把形状尺补到它自己声称的样子**，提交 **`c5be2d703`**（只 commit ⛔ push）。★**D1**：`panelProfileFold` 折叠 `+` 链／`strings.Join` 字面量切片／`Join(...)`＋字面尾巴，叶在折叠值上判；只一半常量的叶（`"wisp-panel-"+pid`）判静默（收紧，不是放宽）。★**D2**：根扩到 `os.Getenv("TEMP"|"TMP")`、`os.UserHomeDir()`（含 `home, err :=` 二值赋值）、盘符绝对字面量；孤立绝对字面量只走 SINK（A3c 的正形就是交给 `NewPanelManager`），所以测试数据里只是拼着 `profile` 的 Windows 路径不冤枉。★**D3**：包级声明同时走 NAME 与 SINK，包级固定 `Join` 变量 seed 进每个函数体。★**D4**：`newPanelControl`／`DataPath` 字段赋值／函数体内 `DataPath:` 键各加 1 条常驻必红控制。★**D5**：汇总行印「形状命中 N／其中已豁免 M／实红 K」＋「豁免自检红 F」，`K>0` 或 `F>0` 一律写【红】。★★**仓外可编译突变体**（`C:\Temp\wisp-ruler2-20261114`，旧尺/新尺两棵树放**同一枚**探针，`go test -c` 都 rc=0）：旧尺 275 文件/2 命中/**0 红 PASS**（12 形全静默），新尺 275 文件/14 命中/**12 红 FAIL**，红句逐字点名探针 `:36/:40/:44/:48/:52/:57/:61/:69/:73/:77/:84/:86`。★★**常驻控制有牙**：仓外 `teeth` 树四点自残——从 sink 集合删 `newPanelControl`／让 selector-LHS 字段臂永不触发／让 `DataPath:` 键臂空转／汇总退回旧口径——分别红在 `the_newPanelControl_sink_arm`／`the_DataPath_field-assignment_sink_arm`／`the_DataPath_key_sink_arm_inside_a_function_body`＋`a_package-level_DataPath_key_is_a_sink_where_it_is_declared`／`the_summary_line_cannot_read_green_while_anything_is_red`；每次还原后 `RESTORED_SHA == REPO_SHA == 6E54A0B3…` 且跑绿。★**常驻正控 14→36 例＋1 条汇总控制，全 PASS**；EXEMPTIONS **原样 2 条未动**。★门禁：`go vet -tags winlive ./cmd/wisp` **rc=0**、`gofumpt -l cmd/wisp` **空**、`sh scripts/d22scan.sh` **clean rc=0**（cmd/ 280 文件）、path-length **VERDICT GREEN**。★**仍漏（已具名写进测试件头注释）**：A1b（helper 返回非 profile 叶、值到下沉点）与 A1c（根在 helper 形参）——都要值跨函数边界，要跨过程分析，是另一件仪器；四根之外的根、一半常量的叶刻意不追（理由在头注释）。★**未跑整包 `./cmd/wisp`**、未 push、未改产码/其他测试/`scripts/**`/`internal/**`/规格/阈值/golden。证据＝`docs/evidence/s1/314-shape-ruler.md` §7，读数根＝`.scratch/wisp/probes/314-ruler2/`。撤销口令＝「314 形状尺补洞不采信」：`git revert c5be2d703`，并删本节这一行。

---

## 非实现者审计（第二轮，补洞后的形状尺）：D4/D5 真修好；D1/D2/D3 只修了示例（2026-11-14 15:2x）

审计件 **`docs/evidence/s1/314-ruler2-audit.md`（新）**，读数根 `.scratch/wisp/probes/ruler2-audit/**`（9 份日志、4 枚探针源、旧/HEAD 尺备份）。
审计对象＝`cmd/wisp/panel_profile_leaf_shape_windows_test.go` @ HEAD `a46352fd3`（blob `627c2fcd`，SHA256 `6E54A0B3…`＝`c5be2d703` 的交付）。
树＝`git archive HEAD` 导出 `C:\Temp\wisp-ruler2-audit-20261114`；⛔ 仓内零跟踪件改动；⛔ 未跑整包（只 `-test.run '^Test314'`）；⛔ 未 push。
★审计期间工作树那枚尺被另一个腿改成未提交的 442 实现（1630→2362 行，SHA 变两次），本件只认 `git archive HEAD` 的已提交版本。

**六节结论：1 漏／2 漏／3 漏／4 抓／5 抓／6 抓。**

- **基线**：274 文件/2 命中/2 豁免/实红 0/自检红 0，rc=0；`-test.v` 36 个子例＝35 形状＋1 汇总（证据 §7.4/§7.7 的「36 形状＋1 汇总」多算一条）。
- **攻击 1（D1 折叠）漏**：`"a"+suf`（suf 是同文件/别处 const，全常量）、`fmt.Sprintf`、`strings.Replace`、切片变量 `strings.Join`、`append(...)...` 展开——全静默（探针 `:24/:28/:33/:38/:43/:47`）；三参 Join `:51`、嵌套 Join `:56`、Join+Join `:60` 红。因 `panelProfileFold` 的 Ident 分支（`:637-641`）不做常量求值，CallExpr 只认三种调用（`:655-666`），`strings.Join` 只吃 `CompositeLit`（`:699-702`），`...` 落兜底 opaque（`:666`）。
- **攻击 2（D2 根）漏**：`os.Getenv(变量)`（探针 `:71`）、`filepath.Dir(os.Executable())`（`:76`）、根先经 `+` 拼一段再绑局部（`:104`）静默；`path.Join`（`:94`，被当 Join）与 `root := os.Getenv("TEMP")`（`:99`）红。`Getenv` 只认字面量（`:878-880`），根源只认三调用（`:848-869`），别名只登记直接根源（`:806-808`）。`TMPDIR`／`t.TempDir`／`MkdirTemp` 静默是具名射程/每发独立，不算洞。
- **攻击 3（D3 包级）漏**：包级固定 Join 在**另一文件**（声明 `d3a:15` → 使用 `d3b:9`）、`init()` 里赋值（`d3a:20`→`:35`）、`sync.Once` 里赋值（`:29`→`:40`）——三形全静默；同文件 `var = Join(...)` 对照（`d3c:14`→`:17`）红。seed 逐文件算（`:401/:405/:445-448`），`panelProfilePackageFixedJoins` 只收 `ValueSpec`（`:545-553`），`AssignStmt` 只 `stored` 不 `add`（`:497-500`）。⇒ 修的是「同文件值规格」，不是包作用域。
- **攻击 4（D4 正控）抓**：仓外 `teeth` 树逐臂打断——T1 删 `newPanelControl`→只红 `the_newPanelControl_sink_arm`；T2 字段臂永不触发→只红 `the_DataPath_field-assignment_sink_arm`；T3 键臂空转→红 `the_DataPath_key_sink_arm_inside_a_function_body`＋`a_package-level_DataPath_key_is_a_sink_where_it_is_declared`（同一臂，2 控制）。无错配；还原后 SHA 同串、跑绿。
- **攻击 5（D5 汇总）抓（两枚被指读数）**：K>0（探针 `:11`）→`【红】…实红 1 处`；F>0（删 holder `:46` 理由句）→`【红】…实红 0、自检红 1`。残余：正控表红时扫描测试的汇总仍印 `【绿】`（`a4-teeth.log`），因 `clean` 只看本测试两计数（`:247-256`）；「任何一处红」措辞大于射程，实际数字没回退。
- **攻击 6（36 子例）抓**：14 例与旧尺 `cb308abd` 逐块归一化比较**全 SAME**；断言从 `Contains(code, leaf)` 变 `raw != "" && Contains(code, raw)`＋显式 `lineNeedle`，**更强**不是更弱；总数 35 形状＋1 汇总＝36。
- **AC#4 判读不变**：本尺仍不能单独翻 AC#4（残余 holder/spike 是登记不是清零；换根/拆字面量仍可绕；AC#2 静桌面读数仍缺）。
- 撤销口令＝「314 形状尺二轮审计不采信」：删本节与 `docs/evidence/s1/314-ruler2-audit.md`，读数根可原样留着复跑。
