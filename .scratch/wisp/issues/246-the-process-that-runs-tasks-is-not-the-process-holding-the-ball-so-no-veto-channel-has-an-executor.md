# 246 — 反向缺口：**会跑任务的那条腿不在带着球的常驻进程里** ⇒ 审批卡片进不来，D43 承诺的四条否决通道（单击球／Esc／语音否决词／面板拒绝）**在生产进程里一枚都没有执行者**

- Status: **已立，未派**（09-30 18:2x，编排者立；账 `A480`）。触发＝验收腿 `245-v1` 交回的一条事实，我自己现跑复认过。
- 与票 228 的关系：票 228 修的是**半边**（"球与托盘不在跑任务的那个进程里"⇒ 把球搬进常驻进程，已落）。**这一枚是另半边**：常驻进程今天仍然**不跑任务、没有审批门、没有面板宿主** ⇒ 球进来了但**没有东西可否决**。两枚合起来才是那句产品事实："**双击图标起来的那个进程，能干活、能被打断**"。

## 现量（09-30 18:23 编排者自己跑，别信行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| 常驻那条腿自己逐字承认没有管线 | `cmd/wisp/resident_ball_windows.go:22` 逐字「pipeline, no microphone, no approval gate and no panel host today, so this…」；`:126` 逐字「gestures", "recorded only: this leg has no task pipeline, no microphone and no approval gate」；`:163` 逐字「…so the gesture has no executor here」 | `grep -n "Confirming\|approval" cmd/wisp/resident_ball_windows.go` 我现跑 |
| "借 Esc"那一支在生产里零调用者 | `internal/ball/ball_windows.go:873 TakeEscForCancel`／`:895 ReleaseEscAfterSession` 的非测试消费者**只有旁支程序** `cmd/balldebug/main.go:635`／`:638` | `grep -rn "TakeEscForCancel\|ReleaseEscAfterSession" --include=*.go internal cmd \| grep -v _test` 我现跑 |
| 契约要求的是**四条**否决通道 | `docs/PLAN.md:3082`（**D43 转移表＝C12 冻结**）第 22 行逐字「否决（单击球 / `Esc` / KWS 否决词 / 面板拒绝，B1）」；`internal/agent/approval/approval.go:90` 逐字 `ChannelEsc: "按 Esc 键"` | 票 245 已逐字读过，本票**不改那两枚文件一字** |

## 前段（派单第一枚只能是普查，不许直接开写）

- [x] **AC#0 先把"卡片怎么进常驻进程"的两形代价摆开（不许直接开写）**：甲＝常驻腿**自己起一条 loop／审批门**（新依赖边？谁装配？）／乙＝**装配根 `cmd/wisp` 注入**（与台账「票 197 段：装配根是唯一的接缝」一致）。完成判据＝两形各带"要动哪几枚文件＋新增哪几条依赖边＋会不会破 `internal/proc` 的 D38(e) 十步顺序＋与票 228 后续片（`config.toml` 未接、托盘「退出」无执行者）谁先谁后"，由编排者裁后再派落地腿。⛔ **普查腿不许改任何产码**；写点只准落在 `.scratch/wisp/probes/246/a1/` 与本票面。

## 判据（前段裁完之后才许补，⛔ 现在不填）



## 编排者裁定 — 09-30 **18:40**（账 `A481`）：`246-a1` 两形代价已交，**我裁走乙（装配根注入）**；AC#0 翻勾，判据从"待补"改成本票下面的六格

**盘上实证**：表 `.scratch/wisp/probes/246/a1/census.md`＝**206 行／31,880 字节、真占位符 0 枚**；四枚 commit 逐枚 `git log -1` 复认（「f9cbe725 骨架 18:29 → 63508c50 填满 18:38 → 56c6c46f 票面一行 18:38 → 2e67b24f 起手名册 18:39」）；它零改产码（`git status --porcelain internal cmd docs/PLAN.md`＝**0 枚**）、票面 AC 框一枚未动。
**我另跑的三把复认尺（乙形的三条承重理由我逐条自己读过，不是信它）**：① `cmd/wisp/approval_reply.go:67` 与 `cmd/wisp/run.go:47` **今天就在产码里 import `internal/agent/approval`** ⇒ 乙形**确实零新增包级边**；② `cmd/wisp/run.go:336 assembleRuntime`／`:782 (*agentRuntime).close` 在场，而常驻那条腿 `cmd/wisp/resident_windows.go:118` 只做了 `rb := startResidentBall(rt.Registry)` ⇒ 甲形要复刻 `assembleRuntime` ＝**在同一棵树里造第二个审批门真相源**；③ `internal/proc/boot_windows.go:151-161` 现读逐字 `hooks := ShutdownHooks{}` 之后**只填 `CloseJob`** ⇒ D38(e) 十步里第 1–7 步今天**零生产者**（它这条判语成立）。

## 灰区那一枚，我自己裁（⛔ 不摆给 owner，因为代码自己的注释已经答了）

`246-a1` 把"给 `proc.Shutdown` 补钩子注册入口算不算契约面"判不动上交。我裁：**不算**。凭据是那段注释**逐字**写着「Later tickets register their hooks on the same sequence - **the order itself is frozen and audited**」（`internal/proc/boot_windows.go:148-150`）⇒ 冻结的是**顺序**，注册入口是它**预先写给后续票的**。本条即批准记录；⛔ 十步顺序一字不许动，只准往里挂钩子。

## 判据（本票现行射程＝最短链：让一张真卡片进得来、Esc 真能否决一次、退出时不撒谎）

**AC#0 已勾（凭据见上面那一格与本节）**：（凭据＝上面三把我自己跑的尺＋表 §甲/§乙；裁定＝**乙**）。
- [x] **AC#1 装配根注入一枚真审批门到常驻腿**：由 `cmd/wisp`（装配根）把 gate 递进常驻那条腿，⛔ **不开任何新包级依赖边**（尺＝`go list -deps`／import 差集落地前后**逐名相同**）。完成判据＝差集现跑读数进表。
- [x] **AC#2 一张真卡片能在常驻进程里挂起来**：走 `gate.go:188/246/415`＋`replies.go:252/455` 那一族既有机制造**真**待决项（⛔ 不许 mock 代替真的来报完成），状态机真进 `Confirming`。**可见证据今天只到"日志＋状态位"那一层**：⛔ 面板开不出来是既成事实（`approval_always.go:165` 逐字 this binary links no WebView2 host），所以这一格**不许写成"用户看得见卡片"**。
- [x] **AC#3 `Confirming` 期间 Esc 真能否决一次、完事归还**：复用票 245 那套"第二个进程观察 Esc"的台件（⚠ 与票 245 AC#9 有先后：**台件先进仓，本格才守得住回归**）；完成判据＝借到／否决生效／归还后另一进程收得到 Esc 三形读数都在。
- [x] **AC#4 卡片挂着时收到退出信号不许撒谎**：D38(e) 顺序一字不动，把"取消任务根"做成**真注册的钩子**（现在只填 `CloseJob`＝第 1–7 步零生产者），并保证待决卡片在退出路径上被判**拒绝＋留审计**；完成判据＝一发真机：卡片挂起 → 发退出 → `StepRecord` 里那一步**不是 skipped**，且 `problems()` 不冒出"有未注册产物"的假账。
- [x] **AC#5 起手第一发必须先归因那枚红**：`cmd/wisp` 默认层今天有**一枚未归因的红**（用例名丢失，账 `A480` ④）⇒ 本腿开工前跑默认层、把红名取出来、判"起手即在 vs 本票造成"，读数写进表 §门禁；⛔ 不许当已知常红略过、也不许顺手修别人的东西。
- [x] **AC#6 四条否决通道的其余三条不许顺手做**：单击球与 KWS 否决词、面板拒绝本票**不实现**（面板属前端会话、由 owner 自己带话；KWS 属唤醒词那一段）⇒ 完成判据＝票面「Progress log」具名写"本票只落 `Esc` 一条通道，另三条各有归口"，⛔ 不许留成"看起来四条都通了"。

**归口（`246-a1` 的 D-1..D-5，我一枚没动）**：proc 七步空转＝本票 AC#4；缺注册入口＝本票 AC#4 的前置（灰区已裁）；过期注释两处（含 `ball_windows.go:809-812`）归票 228 后续片与票 245；那枚红归因＝本票 AC#5；`ball.New` 失败路径清理归票 228。
**串行**：本票动 `cmd/wisp`＋`internal/proc` ⇒ ⛔ 与票 228 后续片、票 245 AC#6..AC#9 全部串行；派单必写"此刻哪几枚包有别的写腿"。

## 禁区

- ⛔ 不动 `docs/PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`allowlist.txt`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；D43 转移表那句"四条否决通道"**本票只能去实现它，不能去改它**。
- ⛔ 不新造第五枚否决通道、不改冻结契约 C12／C18 的形状；`frontend/**`／`design/**` 零读零写零转述。
- ⛔ 面板侧要的东西**只写进本票面与台账**，由 owner 自己带给他用的那枚前端 agent；本编队永不调用跨会话工具去联系任何会话。
- ⚠ **下一枚动 `cmd/wisp` 的腿必须先归因一枚已知的红**：`245-v1` 自报该包默认层「1 红 2 绿」但用例名丢了（`A480` ④）⇒ 起手先跑默认层把红名取出来、判断是不是本票造成的，⛔ 不许当已知常红略过。
- ⛔ **不新增依赖边要先经编排者裁**（本仓有"两枚正向边一律不开、改注入"的先例：票 238 第 1 刀）。
- git：只 commit 不 push；显式 pathspec；禁 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件。

## 编排者收表 — 09-30 **19:51**（账 `A482`，**这一枚才是真落地的 A482**）：非实现者「246-v1」交回 ⇒ AC#0..AC#6 全勾，⛔ **不加 `-done`**（另补两格未勾）

裁决表 `docs/evidence/s1/246-veto-channel-in-the-resident-process-v1.md`＝**167 行／19,471 字节、占位符 0**；两枚 commit 复认（「14676b1c 骨架 19:38 → 71d67b00 终表 19:50」）；终态名册＝起手名册（`git status --porcelain internal cmd` 现跑 **0 枚**，突变全还原）。
**凭据逐格**（表 §1，全部它自己现跑）：AC#1 零新增边（`go list -deps ./cmd/wisp`＝276、`./internal/proc`＝147，本票 5 枚内部 import 在 base「ce1ade1f」就已是 `cmd/wisp` 直依）；AC#2 真门真卡（真 `approval.New`＋真 `ballCardUI`＋winlive 断 `WaitingState()==Confirming`；**无 mock 代真、无"看得见卡片"主张**）；AC#3 三形齐全（`observer keydown 1→0→1`、`ANSWER-VETO … channel=esc`、**归还由第二个进程真收到**，且台件 `cmd/wisp/testdata/esclistener/main.go` **已进仓**＝票 245 AC#9 那条前置一并销掉）；AC#4 退出不撒谎（step3 `Skipped==false`、待决卡真被拒、`shutdown.go` 的 `RunShutdownSequence` **一字未动**＋顺序审计 PASS）；AC#6 其余三条通道没被顺手做。
**⛔ 一句要认的**：AC#5 本格我判勾，但**归因不是写腿做的**——写腿撞帽前根本没跑整包。是这枚验收腿把 `A480` ④ 那枚"丢了用例名的红"名字找回来：**`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`**（`cmd/wisp` 默认层整包序 3 发 1 红、隔离 `-count=5` 5/5 绿、与 246 兄弟用例配对 `-count=6` 6/6 绿 ⇒ 起手即在的 boot 时序竞态，落点在 `resident_windows.go:62→106` 这段本票没碰的代码）。**它同时判写腿那句"起手 2/2 绿＝无红可归因"是过度主张**，我认。
**第 6 面（我自己代提那半发增量）被定向突变验过不是装饰**：MUT-6a 中和 `if !rb.cancelHosted` 守卫 ⇒ 红句逐字「bindBallHost loaded the cancel channel with no executor behind it (ball up = true)」＋`--- FAIL: TestAC246ChannelNeedsBothWindowAndExecutor`。**记它一功：这发增量从没被写腿完整验证过，是这枚腿第一次替我兜住。**

- [x] **AC#7（新补）常驻进程今天还没有"举卡的来源"**：`newResidentApproval()` 有产码调用者（`resident_windows.go:122`，双击那个二进制里真装配），但 **`askConfirmation`／`AskOnTaskRoot` 的产码调用者＝0**（只有用例）⇒ 门在、卡进得来，**还没有任务会去举它**。完成判据＝常驻那条腿真起一条任务管线（或经装配根注入一条最小任务源），一发真机走通"起任务 → 举卡 → Esc 否决 → 任务被取消"四步；⛔ 不许用测试构造的任务源冒充。**这是票 246 标题那句"永久之家"真正落地的地方，属后端主链，排队列最前。**
- [ ] **AC#8（新补）那枚 boot 时序 flake 要有归属**：修法在常驻启动路径（候选＝把 `signal.Notify(bootExit)` 提到 `installLogSink` 之前，消掉 `[sink→Notify]` 那段窗口），**地界属票 228**（用例名就是它的 `TestAC228…`）⇒ 本格只做"归口＋判据"，⛔ 不在本票改；判据＝整包序 `-count=10` 该用例**零枚** `0xc000013a`＋红，且不得靠放宽断言或加 `sleep` 达成。

**编排者收件（2026-10-01 09:3x，`246-v2` 判语到手 ⇒ AC#7 翻勾）**

裁决表 `docs/evidence/s1/246-resident-task-source-v2.md`＝**237 行／38,408 字节**；它死于连接中断（09:29 通知），但 §2 判语／§5 十一条／§6 七笔更正**全在盘上没来得及提交**，由我代提为 `6b8327d4`（95 增／4 删，删的只有它自己那三行占位标题）。⛔ **§7 收尾三把尺与 §8 两节至今为空，我不代填**（填了就把"谁做的判"洗混）；我自己现量的卫生三把尺入 `A488`：`git status --porcelain -- cmd internal`＝**0 行**、`git diff --numstat -- cmd internal`＝**0 行**、`grep -rn "MUT-" --include=*.go cmd internal` 在**产码**里＝**0 处**（测试里那 15 处是早先票号的历史指认，本腿一枚没新增），六枚被改锚点逐枚 `grep` 复现原句（`interactiveStdin()`／`if s.gate != nil {`／`baseCtx := parent`／`src := startResidentTaskSource`／地标句子 2 处／谓词 1 处）。
**翻勾凭据（它的自取读数，⛔ 措辞必须带栏头）**：〔接缝注入〕四步 6 发真机全过（`-count=1` 一发＋`-count=2` 两发×两枚用例，0 skip、卡片编号每发不同）；〔真模型〕那一栏**本机为空**（`%APPDATA%\wisp` 无 `config.toml`、无 secrets，与编排者 22:5x、`246-r2` 23:2x 三枚读数同形）。判语正文第 3 截那枚可操作问句我认：**把产码里那枚消费者删掉，四步读数会不会死？——MUT-4 会让两枚真机用例立刻死**，所以"测试构造的任务源冒充"没有触犯。零新依赖边 276→276 逐名相同；重锚的钉消费者 3 枚→6 枚、MUT-5 三枚点名红。

**它交回十一条判不动，逐条裁（K1–K11）**：

- **K1 Esc 的语义对象**：不动 `D43`／`C12` 一字，票面那句"任务被取消"按**两形都量过**读（D43 第 22 行＝取消那一发调用；D38(e) 第 3 步＝取消在跑的任务）。冻结文字没被碰 ⇒ **不构成契约变更、不摆 owner**。
- **K2 真控制台那一支的台件**：**今天不立**。`grep -rn "runConsoleLoop" --include=*_test.go cmd/`＝**0 处**（只有反向覆盖 MUT-1＋写腿手测）⇒ 登记为欠账，**排在票 244 定案之后**——244 一落地那一支可能整个没有用户，先造台件＝可能白造。
- **K3 两枚零真消费者接缝**（`AskOnTaskRoot` 非测试调用者 0／`askConfirmation` 唯一产码调用者就是它 ⇒ 传递不可达）：**保留、不删**。具名一句给下一程：`resident_approval_windows.go:252-256` 那段"它的调用者只有本包用例"**今天仍为真，不许当缺陷清掉**；要"宿主发起的卡片"那一形的那程再走它，走之前不许把"接缝存在"读成"接上了"。
- **K4 两处被本票改假的自述句**：`resident_ball_windows.go:21`（注释）／`:160`（**运行期真打进台账的 slog 字段**）／`:203`（`ballGestureWhy` 常量）逐字"this leg has no task pipeline"，`main.go:61` 注释仍"empty event loop"。⇒ **当场立格，不留残言**：归口票 228 后续片新 AC（那一腿的球自述要在常驻腿上补），本轮由我在票 228 面落格。
- **K5（最重一笔）`-H=windowsgui` 一落地，AC#7 那条真任务源会静默关掉而仪器不响**：三把凭据我复认成立——① `objdump -p build/wisp.exe`＝`Subsystem 00000003 (Windows CUI)` ⇒ **今天**双击进程有真控制台缓冲、`interactiveStdin()` 返回 `os.Stdin`；② `cmd/wisp/console_windows.go:33-43` 的 `attachParentConsole` 只会 `AttachConsole(ATTACH_PARENT_PROCESS)`，Explorer 拉起的 GUI 进程没有父控制台；③ `SPEC-11:50` 逐字要求"无参＝GUI（`-H=windowsgui`）"。⇒ **已落成票 244 的新硬 AC（见 `A488`）**：切那一刀之前必须先答"GUI 进程的任务源从哪来"，且必须有一枚尺在"常驻腿零任务生产者"这一形上**响**（带正控）。
- **K6 注入 gate 的两处真实代价**（`Options.Grants` 为 nil ⇒「本会话内允许」放行但不落盘并写 `GRANT-DROPPED`；吃不到 `[risk] confirm_timeout_sec`——窗口那项即便接了也钳在 3s，只有超时是真差异）：**不现在改 `approval.New` 的位置**（那等于改 AC#1 裁过的乙形次序）。⇒ 落成票 248 的新 AC 一格，由"配置进得来"那票一并收。
- **K7 D38(e) 第 1／第 7 步钩子**：AC#4 原文只写第 3 步 ⇒ **我追认，并具名"这是我的推广"**：十步顺序一字未动（只是往 `RegisterShutdownHook` 挂函数），真机断言 `10 steps, 0 failed`＋boot 报告含 `1:scheduler-close`／`7:flush-logs-close-db`，6 发都过。授权记录＝`A488`。
- **K8 新增的 `task <文本>` 控制台动词**：**认可**。尺＝`grep -n "控制台" docs/PLAN.md` 只命中 `:146` 那枚 Porcupine 厂商控制台、`grep -n "命令表\|动词表" docs/PLAN.md docs/specs/*.md`＝**零命中** ⇒ 冻结表里没有"CLI 动词清单"这一张，加一枚动词不触碰契约面；且裸行一律交现成 `replySurface.handle`、未知指令打印拒绝＋复读帮助（`:376-387`），不会把打错的答复变成一次真工具调用。
- **K9 一次只接一发任务**（`:411-414` 第二行**响亮拒绝**）：认可为设计选择；并发属票 197／211 那一族，不在本票。
- **K10 `winlive` 在 CI 无 tag 档**：第四次具名（票 62 AC#9／票 245／票 246-v1／本次）。真做要单开 workflow＝动票 134 的 C+B 形状（契约级）⇒ 登记成 `Q-75` 三栏（甲／乙／不做），**默认不做、本轮不摆 owner**（他要的是核心功能先过）。
- **K11 `internal/ball` 那枚红的归属**：它判"由工作树里别人未提交的删除造成、非本票、非 flake（`-count=2` 2/2 复现）"，但它引用的凭据在**它没写的 §7** 里 ⇒ 我 09:3x 自己复跑 `./internal/ball ./cmd/wisp` 定案，读数进 `A488`；⛔ 那个地界属别人，我不还原、不提交、不删。

**AC#8 保持未勾、保持归票 228**（它复认归口对：用例名是票 228 的名字，落点那段启动窗口本票五枚提交一行未动；它**没跑** `-count=10`，起手那一发整包 0 红既不构成复现也不构成推翻）。

**未列进判据、但具名登记的两条**：① 验收腿判不动"**本票两枚新真实进程／球窗用例会不会抬高 flake 发生率**"（方向可推，3 样本不足定论）；② `winlive` 这一族在 CI **没有 `-tags winlive` 档** ⇒ 上面那些"真机读数"目前**只有本机有**（与票 62 AC#9、票 245 交回的同形，第三次出现，仍无人立尺）。
**另两处过期指认（它顺手发现，我一枚没改）**：`approval_always.go:165` 那句"球只由 `cmd/balldebug` 建造"在 228／246 之后**已过期**；`PLAN.md:73`（D2"常驻是不是任务管线之家"）本票三程都**没逐字核过** ⇒ 归票 228 后续片与本票 AC#7 的前置阅读。
**本票状态**：勾 8／未勾 1（只剩 AC#8，地界在票 228）⇒ **不加 `-done`**。

## 排程与串行



- 排在**队列最前**（先于 `197-r3`）：owner 09-30 原话「你可千万别忘了要紧事儿，**你得赶紧干后端的活儿哈**」——本票是后端主链（"能干活＋能被打断"），票 197／224 是仪器精度。
- ⛔ 与票 228 后续片、票 245 的 AC#6..AC#9 同撞 `cmd/wisp`／`internal/ball` ⇒ **一律串行**，跑突变的腿绝不并发；派单里写死"此刻哪几枚包有别的写腿"。
- 测量坑：跑真机热键／桌面用例前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 与 `wisp.exe` 计数必须为 0；`cmd/wisp` 缺 sherpa PATH 会 `0xc0000135` 且**无 `--- FAIL`**＝根本没跑。

## Progress log (append-only, newest last)
- [2026-09-30 18:2x +08] agent=246-a1 did=前段只读普查交件 `.scratch/wisp/probes/246/a1/census.md`(206 行/31,880 字节, 占位符 0, 骨架先落 f9cbe725→填满 63508c50)。三条现读凭据独立复认成立; 两形代价表各 ①-⑥; 依赖边差集为空(甲乙合规做都 0 条新增包级边, cmd/wisp 已 import agent/approval/ball/proc); D38(e) step3 cancel-task-roots 是 proc 侧两形共用缺口(boot_windows.go:151-161 只填 CloseJob, step1-7 零生产者)。建议乙形(装配根注入, 守"唯一接缝"、不造第二真相源), 最短链『卡片进得来+Esc 真能否决一次』可只注入裸 gate(不需 loop/provider/config/麦)。硬契约改动≈0(D43/C12 一字未动), 需告知功能 3-4 枚, proc hook 注册入口是否算契约面待 owner 裁。go build/vet rc=0; 未跑 go test; A480④ 未归因红留落地腿。零产码、未碰任何 AC 框。next=编排者裁甲/乙后派落地腿

## 编排者裁定（2026-10-08 18:5x；凭据＝只读普查腿 `246-raisercensus-2`；⛔ 票面原句一字未改、⛔ 不撤任何勾）

- **交件**＝`bb1d5ef8`（起手锚）／`eb2a217a`（普查件 `01-census.md`），腿名 `246-raisercensus-2`。它接着上一枚（被我含糊的停手条件卡住的那条）把**第 1/2 步**跑完；全程零 Go 命令。
- **逐跳（锚→枚数→档位）**：`main.go:66`→`resident_windows.go:260` 起任务源 1 枚〔已接〕；`rts:278` `assembleRuntime(spec:269/270)`〔已接〕；`run.go:606` 注门＋`:748` 交 bridge〔已接〕；控制台 `task` 动词 `rts:387`→`submitTask:430`〔已接〕；`bridge:352`→`:443`/`:458` 两读点〔已接；条件＝有效级 L1 无会话授权／L2 才问门〕；`gate.go:294`/`:536` 两处 `ui.Prompt`〔已接，我复跑：全包恰 2 枚〕；`ballCardUI.Prompt:864`〔已接〕；**`AskOnTaskRoot:697`→0 枚〔建了但没接〕**；**`Gate.Replay` 零调用者〔建了但没接〕**。Esc 链与 D38 第 3 步钩子在场。
- ★**裁定（那一问：起管线算不算满足 `AC#7`）**：**算**——`AC#7` 的完成判据是"起一条任务管线（或经装配根注入）"，该半今天**静态成立**（无控制台则打 claim 返回 nil，⛔ 不冒充）；**"真举出一张卡"不是这一格**，它要的是**行为凭据**（任务真跑到需要批准那一步 ⇒ 门调 `ui.Prompt`），今天零凭据（winlive 被抑制、〔真模型〕空）。⇒ ⛔ **不撤勾**；把"真举卡行为凭据"**另立为一条待人派的真机读数**（归"真控制台/真窗"那一族，⛔ 不是契约问题、⛔ 不需要 owner 拍板），已登记进台账当下一波候选。
- ⚠ **两处引用修正（都在我或旧票面上）**：①我 `A725` 引 `resident_windows.go:261`，实测 `src := startResidentTaskSource(rt, ra)` 落 **`:260`**（今日第三次"行号与短语错位"，⛔ 不回改、就地打旧）；②**票面旧引 `:122` 已过期**，现读作 `WithConfig:132`（它现取）。
- ⛔ 本轮**不派码**：`AskOnTaskRoot`／`Replay` 两条"建了但没接"只登记（是死路还是备用路＝待裁，不在本格）。
