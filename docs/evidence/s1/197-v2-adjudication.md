# 197-v2 裁决件 — 票 197 五枚未勾格（AC#0／#1／#2／#3／#6）逐格裁决

> **本件是什么**：只读裁决腿 `197-v2` 对票 197 五枚未勾格的**非实现者裁决**（`SPEC-12 §4.3` #1/#3、D22 双角色）。
> **谁派的**：派单存档 `.scratch/wisp/probes/orch/20261115/txt/451-197-314-dispatch.txt` §腿 B（第 25–33 行）。
> **锚**：`HEAD = bf4adaa7b57156e8b79ac579763b381092fdee3b`（短 `bf4adaa7b`）；裁决腿起手 `2026-11-15 12:14:55 +08`，
> 落笔 `2026-11-15 12:24:43 +08`（本机 `Get-Date` 现读，与本件末节同一把尺）。
> **写面**：本件与 `.scratch/wisp/probes/197/v2/**`（含 5 发 overlay 突变件与逐发日志）。产码／票面／台账／阈值／契约**零字节改动**
> （终态尺见 §6.2）。`cmd/wisp/**` **只读**（本波归写腿 451 独占，见 §1.3）。
> ⚠ **本件不改票面任何 `- [ ]` 框**：勾框归编排者。本件只给判语与读数。

---

## 0. 逐格判语总表（五列，按派单要求）

| # | AC 原文逐字（票面 `.scratch/wisp/issues/197-…subagent-entity.md`） | 裁决 | 依据 `file:line` ＋原文 | 还缺什么料 |
|---|---|---|---|---|
| **AC#0**<br>（票面 `:36`） | 「**AC#0 先答"这一层归哪个切片"**：`PLAN.md:1544`／`:1447` 两行改判之后，子代理派生落在 S 几（S3？S5 Agent 内核？S7 并发？）——**落在哪一节决定它要不要现在就能开多任务**；答不出就摆 owner，不许我自己挑一节塞进去。」 | **部分** | 票面自答在 `.scratch/wisp/issues/197-…:21-24`：「落在 S5（Agent 内核）为家、S7（并发）为闸」「⇒ 所以本票**不需要"先开放多任务并发"这一步**：名册在场，缺的是**"子代理"这一层实体**」。<br>⚠ 现量三处对不上：① `docs/specs/SPEC-12-roadmap-governance.md:21` 逐字「\| **S5** \| 按需 WebView：命令/结果面板 + L2 确认卡 + 配置编辑器 \|」⇒ **S5 不是 Agent 内核**；② `grep -n "task\.spawn" docs/PLAN.md docs/specs/` **零命中** ⇒ 本层在 D34 权威表里**没有行**；③ 它同族的 `task.list / task.cancel`（`docs/PLAN.md:2570`）与 `task.output`（`:2571`）在表里的**切片列都写 `S7`**，而 `plan.present` 写 `S3`（`:2573`）。<br>「不用先开多任务」这一半**可核为真**：派生是「一次工具调用里起一个子环路」，不经 `TaskScheduler`（`internal/tools/subagent_197.go:286` 只取自己的池位），名册本就多行（`internal/tools/task.go:424`／`:455`）。 | **一枚落节更正**（不是产码）：<br>① 按 D34 的「切片」列读 ⇒ 本层与 `task.output`／`task.cancel` 同族 = **S7**；按「spec 为家」读 ⇒ **SPEC-05／SPEC-07**（`SPEC-05` 标题「Agent 核心：循环、LLM 抽象、上下文与记忆」）——**S5 两读都不成立**；<br>② `docs/PLAN.md:1452-1455` 逐字「这段期间**不得开放多任务并发**」今天一字未改，A399 那句改判是**台账上的总改判**，**没有一句落到 `task.spawn` 这一层**；读冻结文本的人仍会读成禁令 ⇒ 需要一句「本层不属 `TaskScheduler`，故不受 `:1453` 约束」落进票面或 `A##`；<br>③ 本格**没有任何尺**：删掉票面 §AC#0 不会有一枚测试红 ⇒ 落节更正须写成可复核的一行（`A##` 或票面），否则下一腿还会重新问一次。 |
| **AC#1**<br>（票面 `:37`） | 「**AC#1 实体层有生产调用者**：派生入口在非测试码里有真听众（⚠ 今天"零调用方有两种成因：差一根线／差一整层"＝`A382`/`A391`，本票起手就要分清是哪一种）。」 | **满足** | **①生产派发点**：`cmd/wisp/run.go:1418`「for _, e := range tools.BuiltinSubagentEntries(tools.SubagentDeps{」＋`:1428`「}) {」＋`:1429`「if err := reg.Register(e); err != nil {」——唯一装配根 `assembleRuntime`（`run.go:628`）里注册进 C4 目录；目录本身由 `internal/tools/bridge.go:278-279`「func (b *Bridge) Tools(...)」「entries := b.reg.List()」交给模型，派发在 `bridge.go:312`；常驻腿走**同一个** `assembleRuntime`（`cmd/wisp/resident_task_source_windows.go:297`「run, code := assembleRuntime(spec)」→ `:505`「code := src.run.execute(root.Ctx, text)」）。<br>**②摘掉它要红的测试**：`cmd/wisp/plan_wiring_165r2_test.go:178-187` 把本入口当**阳性对照**逐字钉住——`:178`「if len(precedent) != 1 {」、`:179`「阳性对照失败：先例 %s 的产码调用点数出 %d 枚（要 1 枚，run.go 的 subagent 块）」、`:183`「if !precedent[0].inRange \|\| !precedent[0].registers {」；另一枚 `cmd/wisp/subagent_carrier_197_test.go:232 TestRunPacketCarriesTheSubagentItsRosterRowFed` 走真装配的 `rt.bridge.Execute`（`:173`）派一发真子代理，入口不在册时它会 `:271`「the spawn call failed」。<br>**成因分类**（AC 点名的第二问）：票面 `:24` 自答「按 `A382`/`A391` 那把尺分类＝**差一整层，不是差一根线**」——现量与之一致：整层在三枚文件里（`internal/tools/subagent_197.go`／`internal/panel/subagent_roster_197.go`／`internal/streamkey/streamkey.go`）。 | 就本格原文：**无**。<br>⚠ 形态记账：②那两枚都在 `cmd/wisp`，本波该包归写腿 451 独占（§1.3），本腿**只读断言、未跑包**；两枚的断言是确定性的（AST 计数／真桥派发），不是"大概会红"。<br>⚠ 本格不覆盖：子代理**自身**的工具目录是否只剩父目录减 `task.spawn`（那是深度那条格，`subagent_197.go:603-629` 的读数在 `Test197ChildCannotDeriveSubagent`）。 |
| **AC#2**<br>（票面 `:38`） | 「**AC#2 名册与状态同源**：子代理状态**只能**是 D43 那一集；判据要能区分"填了真状态"与"字段恒空"（票 181 `AC#7` 那枚空转教训）。」 | **满足** | **①生产写入／读取点**：在跑态 `internal/tools/task.go:461-464`「r.recordLocked(taskID, TaskOutput{… State: statemachine.State(subagentStateRunning),」；收口态 `internal/tools/subagent_197.go:464-470`（Settling／Muted／Error 三支）＋`:475-478` 写回同一行；读侧唯一判决 `internal/tools/task.go:224-234`「func (o TaskOutput) StateAnswer()」，装配根在 `cmd/wisp/panel_pump.go:232`「status, statusReason := rec.Out.StateAnswer()」；载体对**不被承认的名字**只丢不译（`internal/panel/subagent_roster_197.go:273-282`）。词表只有 D43 那 20 枚：`internal/statemachine/states.go:11-31`＋`:39-49` `Valid`；两枚 197 产码里**没有任何状态字面量**（`grep 'State\s*(:\|=\s*"'` 零命中）。<br>**②摘掉它要红的测试（本腿现跑突变，见 §3）**：<br>• M1（在跑态置空）⇒ `Test197RowExistsBeforeFirstChildModelCall`、`Test197CancelIsPerRowAndNeverCascades`、`Test221ParentCancellationStillDoesNotCascade` **三枚红**，逐字红句「在跑的行 State = "", want "Thinking"」（`subagent_197_test.go:380`）；<br>• M2（收口态置空）⇒ `Test197SpawnPublishesIdentityRow` **红**，逐字「State = "", want "Settling"」（`:321`）＋「State "" 不在 D43 的 20 枚名字上」（`:324`）。<br>"真状态 vs 恒空"这一半另有 `internal/panel/subagent_roster_197_test.go:246 TestAStatusOutsideD43NeverReachesTheWire`（20 枚全过／legacy 词被丢／空白带理由三态）与 `cmd/wisp/subagent_carrier_197_test.go:301-307`、`:353-359`（同一行先后读 `Thinking`→`Settling`，常量字段过不了）。 | 就本格原文：**无**（两向突变都指名了红，且各带了绿着的对照：M1 下 `Test197SpawnPublishesIdentityRow` 仍绿、M2 下 `Test197RowExistsBeforeFirstChildModelCall` 仍绿 ⇒ 不是"一律红"的假牙）。<br>⚠ 本格**不覆盖**根任务行：根行状态今天仍是"空＋理由"（`internal/tools/task.go:471-479` 逐字「State is deliberately NOT written here」，A394 留给票 196）——那是**根**那一行，不是子代理这一格。 |
| **AC#3**<br>（票面 `:39`；判据在 `:32`） | 「**AC#3 每子代理一条流、不许合并**：上面第 2 层那发结构判据。」（第 2 层判据逐字在票面 `:32`：「**造 N＋1 个子代理，没有任何一条会被并进别人的**」） | **满足** | **①生产派发点**：子代理的正文只从两处进流，都用**孩子自己的**键——`internal/tools/subagent_197.go:577`「s.d.Stream.Append(SubagentStreamKey(e.TaskID), e.Text)」与 `:541`「t.d.Stream.Append(key, text)」（`key` 来自 `:535` `SubagentStreamKey(childID)`）；键的**唯一字面量**在 `internal/streamkey/streamkey.go:29`，两边都是别名（`subagent_197.go:54`、`internal/panel/pump.go:763`）；装配根按行铸键 `cmd/wisp/panel_pump.go:221-224`「key = panel.SubagentStreamKey(rec.TaskID)」。溢出语义已重裁：`internal/panel/pump.go:718`「DefaultStreamKeys = 32」、`:723`「StreamKeyHardCeilingMultiple = 2」、`:855-858`「no key is ever folded into another」、`:860-883` `enforceBoundLocked`（越界只把**每行自己**削成头尾，过硬顶才**整键停跟踪并具名**）、`:896-928` `clampLocked`（逐键算，不跨行）。<br>**②摘掉它要红的测试**：<br>• **本腿现跑** M3（把 `Append` 折成一枚键）⇒ `TestNPlusOneSubagentsEachGetTheirOwnStream` 红，逐字「rows = 5 for 4 subagents, want one row per agent」（`internal/panel/subagent_stream_197_test.go:70`）＋`TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther` 红「rows = 4, want 3」（`:184`）；<br>• 生产路径那一发（**未跑**，`cmd/wisp` 独占）：`cmd/wisp/subagent_carrier_197_test.go:525 TestRunPacketReportsTheStreamLogPastItsBound` 真派 **33 枚**（`panel.DefaultStreamKeys+1`，`:531`）并断言 `:578-581` 无两行共用键、`:592`「if len(seenKey) != panel.DefaultStreamKeys+1 {」（＝33 枚各有一条自己的键）；另一枚 `:232 TestRunPacketCarriesTheSubagentItsRosterRowFed` 在 `:308`／`:338` 用**独立算出**的键去找那一页（键被换成裸 taskID ⇒ `carrierPage` 当场 Fatal）。 | 就本格原文：**无**（结构判据两向都有牙：合并⇒红；且 `TestNPlusOne…` 同时断言"截断要说出来"`:106`、`elided=0` `:109` ⇒ 没靠偷偷丢内容换绿）。<br>⚠ 生产路径的 33 枚那一发本腿**未跑**（原因同上）；它比本腿跑的 `TestNPlusOne` 强在"真派 N+1"，弱在"每行正文只有 ~150 rune，读不到中段被削"（该文件 `:660-684` 自己写明了这两条边界）。<br>⚠ 本格**不覆盖界面**：点进去那一页归 owner 带给前端那支（`frontend/**` 零读零写，票面 `:33`／`:70`）。 |
| **AC#6**<br>（票面 `:59`） | 「**AC#6 取消语义**：父任务取消时子代理怎么收（`TaskRoster` 今天有没有 owner／recover 那条链要现读）；⚠ 裸 `go func(` 而无 owner／recover 是 `AGENTS.md` §1.2 硬禁。」 | **满足** | **①生产点（三处，一条语义一条）**：<br>(a) 父取消**不级联**：`internal/tools/subagent_197.go:340`「childCtx, cancelChild := context.WithCancel(context.WithoutCancel(ctx))」；父侧收口那句写给模型看 `:448`（逐字「子代理 %s 没有被级联取消，它仍在名册里，」）。<br>(b) 孩子**有停止出口且行与句柄同一次落地**：`:357-365`（准入钩子里 `:358` `PublishSubagentInFlight(...)`）与 `:387`（派生后同值再发一次）；原子发布体 `internal/tools/task.go:455-469`（`nil` 句柄不发布）；唯一停法 `task.go:569-586 Cancel`，生产调用者 `task.go:908`「stopped, why := t.d.Roster.Cancel(target)」，工具注册在 `cmd/wisp/run.go:826-835`（`task.go:723-729` 的 `BuiltinTaskEntries`）。<br>(c) owner／recover：收尾监视器走**名册命名**的注册表 `subagent_197.go:431-440`（`:432`「observe.Default.Spawn("agent-task-"+bg.ID, "tools", watchRoot, func(c context.Context) {」），recover 边界在 `internal/observe/goroutine.go`（D22 ban #1 的正形）；**产码里裸 `PublishSubagent` 今天 0 处**（`git grep` 非测试零命中），**裸 `go func(` 也无**（同一文件里只有 `observe.Default.Spawn` 这一处起协程）。<br>**②摘掉它要红的测试（本腿现跑突变，见 §3）**：<br>• M4（`WithoutCancel` 换回 `ctx`）⇒ `Test197CancelIsPerRowAndNeverCascades`＋`Test221ParentCancellationStillDoesNotCascade` **3 轮 ×2 枚全红**，逐字「父取消后孩子的 State = "Muted", want "Thinking"」（`subagent_197_test.go:770`）＋「停这一枚子代理没成功：…没有在跑（kind=subagent，state=Muted），没有可停的句柄」（`:775`）；<br>• M5（监视器改成裸 `go func`）⇒ `Test229FinishWatcherIsRosterNamedWhileChildRuns` 红，逐字「收尾协程不在册：CountByName("agent-task-…")=0」（`internal/tools/spawn_finish_watcher_229_test.go:70`）。 | 就本格原文：**无**（"父取消不级联"与"owner／recover 那条链"两个问句都各有现跑红）。<br>⚠ 本格**不覆盖**两件（票面 `:64-65` 已拆票，本腿不回收、也不拿它们拉低本格）：<br>① **用户**停一枚子代理今天仍无落点（`task.go:892-902` 对根行/无父行一律拒，票 181／220 地界）；<br>② 名册 `blockedOnApproval` 只报 L2（票 220）。<br>⚠ `:236` 那行「AC#6 停止出口已由票 221 落 `task.go:742`」的行号已漂：`task.go:742` 今天是 `taskCancelSchema` 的注释，`task.cancel` 名字在 `:749`、注册在 `:723-729`、执行在 `:811`、`Roster.Cancel` 调用在 `:908`（§5）。 |

---

## 0.1 四数与"缺料的格"

- **满足 4（AC#1／#2／#3／#6）· 部分 1（AC#0）· 不满足 0 · 裁决不了 0（总 5）**。
- **缺料的格具名＝AC#0**（唯一）：缺的是**一枚落节更正＋一句"`:1453` 不束本层"**，不是产码（票面 `:234` 那句"五枚未勾格里无一格欠新产码"**在本腿复核后成立**，见 §5.2）。
- 其余四格缺的是"裁决料"，本件即为该料；四格各自的**残余**都已具名指向别的票／别的层（AC#1 无、AC#2 指票 196 的根行、AC#3 指界面那支、AC#6 指票 181／220），**没有一条被本件吞掉**。

## 0.2 A1875 两件套逐格对账（"满足"的硬前提）

| 格 | ① 生产调用／派发点 `file:line` | ② 摘掉它会红的测试名 | ② 的形态 |
|---|---|---|---|
| AC#1 | `cmd/wisp/run.go:1418`＋`:1429`（`assembleRuntime` 里注册进 C4 目录）→ `internal/tools/bridge.go:279`／`:312` | `Test165r2PlanWiringCallSiteCensus`（`cmd/wisp/plan_wiring_165r2_test.go:175`，阳性对照逐字点名本入口）／`TestRunPacketCarriesTheSubagentItsRosterRowFed`（`subagent_carrier_197_test.go:232`） | 〔读断言／未跑包：cmd/wisp 归写腿 451 独占〕 |
| AC#2 | `internal/tools/task.go:461-464`（在跑）＋`internal/tools/subagent_197.go:464-478`（收口）＋`cmd/wisp/panel_pump.go:232`（读） | `Test197RowExistsBeforeFirstChildModelCall`／`Test197SpawnPublishesIdentityRow`（另两枚同族见 §0） | **〔本腿现跑·M1／M2〕** |
| AC#3 | `internal/tools/subagent_197.go:577`＋`:541`＋`cmd/wisp/panel_pump.go:223` | `TestNPlusOneSubagentsEachGetTheirOwnStream`（现跑）；生产路径 `TestRunPacketReportsTheStreamLogPastItsBound`（未跑） | **〔本腿现跑·M3〕**＋〔读断言〕 |
| AC#6 | `internal/tools/subagent_197.go:340`／`:358`／`:387`／`:432`＋`internal/tools/task.go:908`＋`cmd/wisp/run.go:826` | `Test197CancelIsPerRowAndNeverCascades`＋`Test221ParentCancellationStillDoesNotCascade`（3/3 轮）；`Test229FinishWatcherIsRosterNamedWhileChildRuns` | **〔本腿现跑·M4／M5〕** |
| AC#0 | ——（本格不是能力断言，是**落节断言**） | ——（删掉票面 §AC#0 不会有任何测试红） | ⇒ A1875 的②**结构性缺失** ⇒ 只能 **部分**（这是本格降级的机械理由，不是印象） |

---

## 1. 口径（先读这段，再看数字）

### 1.1 本腿跑过什么、没跑什么

- **跑过**（`internal/tools` 与 `internal/panel`，两包均非他腿写面）：
  1. 起手基线 11 枚定向用例（命令与逐枚读数见 §4，全绿）；
  2. 5 发 **overlay 突变**（§3，逐发日志落 `.scratch/wisp/probes/197/v2/<发名>/run.txt`）。
- **没跑**：`cmd/wisp` 任何测试／编译。**理由**：本波派单把 `cmd/wisp` 归写腿 451 独占（派单 `:17`），且工作树里该包正有他腿在飞（`git status --porcelain -- cmd` 显示 `M cmd/wisp/run.go` 与新增 `cmd/wisp/agent_profiles_451.go`）。
  ⇒ 凡引 `cmd/wisp` 测试的地方一律标〔读断言／未跑包〕，**不写成"已证红"**。
- **没跑**：`go test ./...` 整包、任何突变在 `cmd/wisp` 上的那一半。

### 1.2 overlay 突变的纪律（照票面 `:77-81` 那三条，一条不少）

1. **买到什么／买不到什么**：`go test -overlay` **只换编译期源**，盘上产码字节**一字未动**（§6.2 是终态尺）。它**买不到**"读写盘／走 AST 的测试"的可见性——本腿跑的 5 发**全部**是行为测试（真桥、真名册、真 StreamLog、真注册表），**没有一枚**靠读盘或 AST 判分；反之，`Test165r2PlanWiringCallSiteCensus`／`Test229SpawnNamesAreRosterShaped`／`TestSubagentStreamKeyHasOneMintSite` 这三枚**读盘/AST** 的尺本腿**没有**放进任何一发的 `-run` 集合，所以"overlay 对它们结构性失明"这件事在本件里**不构成任何读数**。
2. **锚点唯一性**：驱动器落地前逐枚断言锚点在该文件里**恰好命中 1 次**，否则 `throw`（`mutate.ps1` 的 `New-Mutant`）；5 发全部通过该断言。
3. **还原**：因为用的是 overlay，**没有需要还原的东西**——产码从未被写。每发只写三个临时件（突变副本／`overlay.json`／`run.txt`），全部落 `.scratch/wisp/probes/197/v2/`，**只建不删**。

### 1.3 与"只读"的边界

- 未 commit 任何产码；未 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；未 push；未建 worktree。
- 未碰 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 表／`frontend/**`／`design/**`；未动票面与台账。

---

## 2. 逐格详证

### 2.0 AC#0（部分）——为什么不是"满足"，也不是"不满足"

**票面要的东西只有一句**：答出"落在 S 几"，并因此决定"要不要现在就能开多任务"。

**已答出来的那一半（可核为真）**：
- 派生是**一次工具调用**里起一个子环路：`subagent_197.go:367`「child, err := agent.New(opt)」→ `:386` `child.RunAsync`；它**不碰** `TaskScheduler`（`subagent_197.go:286` 只 `TryAcquireSubagentSlot`，池上限 = 桥的 D38d 天花板 4，`:86`）。
- 名册本来就多行：`task.go:424`／`:455` 两枚发布器 + `cmd/wisp/panel_pump.go:248-255` 走 `Descendants` 枚举行。
⇒ 「本票不需要"先开放多任务并发"这一步」**成立**。

**对不上的那一半（三处现量）**：
1. `docs/specs/SPEC-12-roadmap-governance.md:21` 的 `S5` 逐字是「按需 WebView：命令/结果面板 + L2 确认卡 + 配置编辑器」——**不是 Agent 内核**。「S5（Agent 内核）」这个配对在仓里**没有出处**；`SPEC-05` 才是 Agent 核心（`docs/specs/SPEC-05-agent-core.md:1`「# SPEC-05 · Agent 核心：循环、LLM 抽象、上下文与记忆」），而 SPEC 编号与切片编号是**两套**。
2. `task.spawn` 在 `docs/PLAN.md` 与 `docs/specs/**` 里 **零命中**（`git grep -n "task\.spawn" -- docs/PLAN.md docs/specs/` 空）⇒ D34 的**切片列**给不出本层的落点。
3. 能借的只有同族：`docs/PLAN.md:2570`「\| `task.list` / `task.cancel` \| 查看/取消任务 \| L0 / L1 \| — \| **S7** \|」、`:2571` `task.output` 同列 `**S7**`；而 `plan.present` 是 `S3`（`:2573`）。⇒ 若按"切片"读，本层应写 **S7**，与票面的 S5 **不同**。

**为什么不是"不满足"**：AC 的实质问句（"要不要现在开并发"）答对了、且我复核为真；错的是**名字**，不是结论方向。
**为什么不是"裁决不了"**：这不是"未定义即停"里那六条待定案项之一，也不属"改契约文字"（`AGENTS.md` §1.1 没被触发——本格不需要改任何冻结文本，只需要一句落节与一句约束澄清）。按铁律 2／4，落节是技术判断，**agent 自己拍**，所以"裁决不了"会把该我拍的推回去。

**还缺什么料（三件，按可直接落笔的粒度）**：
1. **一句落节**：二选一并写明理由——(甲) 按 D34 切片列 ⇒「`task.spawn` 落 **S7**，与 `task.output`／`task.cancel` 同族」；(乙) 按 spec 为家 ⇒「落 **SPEC-05／SPEC-07**，本票不新增切片」。**两案都不写 S5。**
2. **一句约束澄清**：`docs/PLAN.md:1452-1455`「这段期间**不得开放多任务并发**（`TaskScheduler` 只允许单任务）」——需写明「本层不是 `TaskScheduler` 的行，故不受此句约束」。A399 的总改判在场，但**没有一句落到本层**，读冻结文本的人仍会拦。
3. **一条可复核载体**：本格今天没有任何尺（删票面那节不会红）。最低成本＝落 `A##` 一行或票面一行，**写明判据是"派生入口不经 `TaskScheduler`"**（该判据今天可核：`subagent_197.go:286`／`:86`）。

### 2.1 AC#1（满足）

- **生产派发链（四跳，逐跳可核）**：`run.go:628` 唯一装配根 `assembleRuntime` → `run.go:1418/:1428/:1429`（`BuiltinSubagentEntries` 进 registry）→ `bridge.go:279` `b.reg.List()`（模型看得见）→ `bridge.go:312` `Execute`（真派发）。常驻腿同根：`resident_task_source_windows.go:297` → `:505`。
- **AC 点名的分类问题**：不是"差一根线"，是"差一整层"——现量与本票 §AC#0（`:24`）一致。
- **为什么敢判满足**：②那枚 `Test165r2PlanWiringCallSiteCensus` 是**阳性对照**用法，不是顺带提到：`plan_wiring_165r2_test.go:177-186` 先数 `BuiltinSubagentEntries`（本入口）＝1 且形状 `range+Register`，不满足就 `t.Fatalf` 并明写"这把尺此刻读不到任何装配"。把人肉改一行（删 `:1418` 那行）⇒ 计数 0 ⇒ 该用例红。**注意它是读 AST 的尺**，对 `-overlay` 失明——所以本腿**没有**用 overlay 去碰它，也没把它的绿当"接线已证"。
- **残余**：无（本格原文之外的东西见右栏）。

### 2.2 AC#2（满足）

- **同源**：行的状态只有两个写者，都写在 `TaskRoster` 的同一行上（`recordLocked`，`task.go:294-316`）；判词只有一处（`StateAnswer`），装配根照抄（`panel_pump.go:232`），载体对不承认的名字**只丢不译**（`subagent_roster_197.go:273-282`）。
- **"真状态 vs 恒空"怎么被判住**：M1／M2 两向（在跑态、收口态）分别红，且**各带一枚绿对照**（§3 表）——这正是票 181 `AC#7` 要求的那种"空转也可判红"。
- **词表纪律**：两枚 197 产码零状态字面量；面板侧连 `statemachine` 都不 import（`subagent_roster_197.go` 只 import `sort`/`strings`，`:41-44`），有 `TestTaskState188AC3PanelHasNoWriteLeg`（`internal/tools/task_state_188_test.go:413`）钉住"面板不许写状态维"。

### 2.3 AC#3（满足）

- **键的单一铸造点**：`streamkey.go:29` 一处字面量，两边别名（`subagent_197.go:54`、`pump.go:763`），`TestSubagentStreamKeyHasOneMintSite`（`subagent_carrier_197_test.go:431`）钉住"非测试源里恰好一处"。
- **溢出不再合并**：`pump.go:855-883` 的注释与被测体是同一条（"no key is ever folded into another"），`clampLocked` 逐键算（`:896-928`），`DroppedKeys` 具名（`:874-881`）。
- **判据落点**：本腿现跑 M3 红（两层：4 枚代理 / 3 枚代理）；生产路径那发（33 枚真派生）**未跑**但断言在场且是确定性的（`:578-581`／`:592`）。
- **残余**：界面那一页（点进去看什么）本格不覆盖。

### 2.4 AC#6（满足）

- **语义三句**（都能在产码里指到行）：父取消不级联（`subagent_197.go:340`）；孩子可被**它的父**单独停（`task.go:908`，注册在 `run.go:826`）；孩子收不了自己、兄弟停不了（`task.go:892-902`，`Test221SubagentCannotStopSiblingOrItself` `task_cancel_221_legs_test.go:413`）。
- **"owner／recover 那条链"现读结论**：链在，且是**名册命名**的——`subagent_197.go:432` 用 `observe.Default.Spawn("agent-task-"+bg.ID, ...)`，不是裸 `go func`；票面 `:236` 记的 `subagent-finish-*` 旧名已被票 229 修掉（`Test229FinishWatcherIsRosterNamedWhileChildRuns` 现在断言 `agent-task-<id>` 在册、旧名不在 `Unknown` 里）。
- **M4 的 3/3 很重要**：这条语义的尺不是"偶发红"，三轮全红、且红句一致（`state = "Muted"`）⇒ 摘掉 `WithoutCancel` 是**确定性**红，不是竞态。

---

## 3. 五发突变逐发读数（overlay；产码零改）

驱动：`.scratch/wisp/probes/197/v2/mutate.ps1`（M5 的编译修正版 `mutate-m5.ps1`，原因见下）。逐发日志同目录 `<发名>/run.txt`。

| 发 | 目标格 | 突变（一个字面改动的形状） | 退出码 | **红的用例（逐名）** | 绿着的对照（证明不是一律红） |
|---|---|---|---|---|---|
| **M1** | AC#2 | `internal/tools/task.go:461-464` `PublishSubagentInFlight` 里 `State: statemachine.State(subagentStateRunning)` → `State: ""` | 1 | `Test197RowExistsBeforeFirstChildModelCall`、`Test197CancelIsPerRowAndNeverCascades`、`Test221ParentCancellationStillDoesNotCascade` | `Test197SpawnPublishesIdentityRow` **PASS** |
| **M2** | AC#2 | `internal/tools/subagent_197.go:464` `finalize` 里 `state := subagentStateSettled` → `statemachine.State("")` | 1 | `Test197SpawnPublishesIdentityRow` | `Test197RowExistsBeforeFirstChildModelCall` **PASS** |
| **M3** | AC#3 | `internal/panel/pump.go:792` `Append` 入口把 `key` 折成 `"subagent:merged"` | 1 | `TestNPlusOneSubagentsEachGetTheirOwnStream`、`TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther` | ——（本发没有设绿对照：这两枚本身就是"合并即红"的正控，且各自带自证输出：`rows = 5 for 4` 那一行直接把"多出来的一行"和"四条空行"一起打了出来） |
| **M4** | AC#6 | `internal/tools/subagent_197.go:340` `context.WithCancel(context.WithoutCancel(ctx))` → `context.WithCancel(ctx)` | 1 | `Test197CancelIsPerRowAndNeverCascades`、`Test221ParentCancellationStillDoesNotCascade`（**`-count=3`，3 轮 ×2 枚＝6 条 `--- FAIL`，无一例外**） | ——（同一轮里 `task.spawn` 的正常收口断言仍绿：`spawn: %v` 未触发） |
| **M5** | AC#6 | `internal/tools/subagent_197.go:432` `observe.Default.Spawn(...)` → 裸 `go func(...)(context.WithoutCancel(ctx))`（并加 `_ = watchRoot` 让编译过） | 1 | `Test229FinishWatcherIsRosterNamedWhileChildRuns` | ——（本发**没有**绿对照：Fatal 落在名册断言 `:70`，其后的「结论回父任务」那几行（`:94-107`）**没跑到**，本件不拿它们当对照） |

**逐字红句（引 M1／M2／M3／M4／M5 各一枚，原样照抄 `run.txt`）**

```
M1  subagent_197_test.go:380: 在跑的行 State = "", want "Thinking"
    subagent_197_test.go:770: 父取消后孩子的 State = "", want "Thinking"（级联会正好体现在这一维上）
    task_cancel_221_legs_test.go:519: 父取消把孩子也停了：state = "", want "Thinking"（级联会正好体现在这一维上）
M2  subagent_197_test.go:318: State 维不可读：宿主没有登记这一维（fail-closed：不替任务编一个状态）（state=""）
    subagent_197_test.go:321: State = "", want "Settling"（跑完的孩子落在 D43 的 Settling 上）
    subagent_197_test.go:324: State "" 不在 D43 的 20 枚名字上
M3  subagent_stream_197_test.go:70: rows = 5 for 4 subagents, want one row per agent -  {"subagent:merged" text="[alpha step1] [alpha ste<...>lta step1] [delta step2]" done=false} …
    subagent_stream_197_test.go:184: rows = 4, want 3 -  {"subagent:merged" text="R1 R2 R3 A1 A2 A3 B1 B2 B3 " done=false} …
M4  subagent_197_test.go:770: 父取消后孩子的 State = "Muted", want "Thinking"（级联会正好体现在这一维上）
    subagent_197_test.go:775: 停这一枚子代理没成功：任务 064428cf-… 没有在跑（kind=subagent，state=Muted），没有可停的句柄
    task_cancel_221_legs_test.go:526: 父任务事后停它被拒：没能停掉 fb9a609b-…：任务 fb9a609b-… 没有在跑（kind=subagent，state=Muted），没有可停的句柄（这一行没有被改动，别的孩子也没被牵连）
M5  spawn_finish_watcher_229_test.go:70: 收尾协程不在册：CountByName("agent-task-b5e7d13a-…")=0；此刻 observe.Default：Total=0 PerTask=0 Unknown=[]；旧名 "subagent-finish-b5e7d13a-…" 的 CountByName=0（AC#2 没落地，或名册判定被改）
```

⚠ **M5 的第一发是 `[build failed]`，本件不把它当读数**：`declared and not used: watchRoot`（突变体自己没编译过，不是判据）。修正后重跑取上面那一条；同目录的 `run.txt` 是**修正版那一发**的日志（首发日志被覆盖，其首行读数即上面这句编译错，本件照抄在此，没有第二份）。修正版驱动器 `mutate-m5.ps1` 与 `mutate.ps1` 的 M5 段同形，差别只有 `_ = watchRoot` 那一行。

---

## 4. 起手基线：11 枚定向用例全绿（本腿自己跑，非引他人）

```
go test ./internal/panel -run 'TestNPlusOneSubagentsEachGetTheirOwnStream|TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther|TestAStatusOutsideD43NeverReachesTheWire|TestTheRosterReaderPutsSubagentsOnTheWire' -count=1 -v
--- PASS: TestTheRosterReaderPutsSubagentsOnTheWire (0.17s)
--- PASS: TestAStatusOutsideD43NeverReachesTheWire (0.03s)
--- PASS: TestNPlusOneSubagentsEachGetTheirOwnStream (0.00s)
--- PASS: TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther (0.00s)
ok  	github.com/CarlosShao/wisp/internal/panel	0.586s

go test ./internal/tools -run 'Test197SpawnPublishesIdentityRow|Test197RowExistsBeforeFirstChildModelCall|Test197CancelIsPerRowAndNeverCascades|Test229FinishWatcherIsRosterNamedWhileChildRuns|Test221ParentCancellationStillDoesNotCascade|Test197SpawnPathPublishesRowAndStopHandleTogether|Test197SubagentPoolNeverExceedsBridgeCeiling' -count=1 -v
--- PASS: Test229FinishWatcherIsRosterNamedWhileChildRuns (0.00s)
--- PASS: Test197SpawnPublishesIdentityRow (0.00s)
--- PASS: Test197RowExistsBeforeFirstChildModelCall (0.00s)
--- PASS: Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)
--- PASS: Test197CancelIsPerRowAndNeverCascades (0.00s)
--- PASS: Test197SpawnPathPublishesRowAndStopHandleTogether (0.00s)
--- PASS: Test221ParentCancellationStillDoesNotCascade (0.00s)
ok  	github.com/CarlosShao/wisp/internal/tools	1.000s
```

⇒ 四格的 ② 不是"今天本来就红"（那会把"会红"当成已有绿）；基线是绿的，红只出现在突变下。**基线绿的这 11 枚里，没有任何一枚依赖本波他腿在飞的文件**（`internal/tools`／`internal/panel` 的产码 `git status` 干净，§6.2）。

---

## 5. 现量更正（票面与 a2 那行的行号／口径）

### 5.1 行号漂移（不影响判语，但引用必须改）

| 票面写法 | 今天（`bf4adaa7b`）现量 | 尺 |
|---|---|---|
| 「`PLAN.md:1544` ＝ REJECTED 多 Agent 协作」 | **`docs/PLAN.md:1550`** 逐字「\| **REJECTED** \| 多 Agent 协作 \| 草案 6.6 明确不做 \| — **不是债务** \| — \| 单 Agent，无编排（⚠ 注意与 D31 的**多任务并发**区分…）」 | `git grep -n "多 Agent 协作" -- docs/PLAN.md` |
| 「`PLAN.md:1447` ＝ S7 之前不得开放多任务并发」 | **`docs/PLAN.md:1453`** 逐字「S7 才补上审批队列与冲突检测。这段期间**不得开放多任务并发**（`TaskScheduler` 只允许单任务），」（上下文 `:1452-1455`）；`:1447` 今天是 C19／S3 那一行 | `git grep -n "不得开放多任务并发" -- docs/PLAN.md` |
| `:236`「AC#6 停止出口已由票 221 落 `task.go:742`」 | `task.go:749`＝`task.cancel` 的 `Name()`；注册 `:723-729`；`Execute` `:811`；`Roster.Cancel` 调用 `:908` | `git grep -n "task.cancel\|func (r \*TaskRoster) Cancel" -- internal/tools/task.go` |
| 票面 `:234`（a2）「实体 `cmd/wisp/run.go:780` 有生产调用者」 | **`cmd/wisp/run.go:1418`**（a2 的锚是 `8e4f6f83`，其后装配根长了 ~640 行） | `git grep -n "BuiltinSubagentEntries" -- cmd/wisp/run.go` |

### 5.2 复核 a2 那句「五枚未勾格里无一格欠新产码…缺的是非实现者裁决料」

- **结论：成立，附一处精确化。** 五格里四格（#1／#2／#3／#6）的产码与判据都已在树里，本件给的正是缺的那份裁决料；AC#0 **也不欠产码**，它欠的是**一句落节与一句约束澄清**（文字，不是码）。
- 同批复核（票面 §编排者裁定 `:244` 那句"生产码里裸 `.PublishSubagent(` 调用现为 0 处"）：`git grep -n "\.PublishSubagent(" -- "*.go"` 的**非测试命中全在 `.scratch/wisp/probes/**` 的历史突变体里，`internal/`／`cmd/` 零命中**；生产的两个发布点都走 `PublishSubagentInFlight`（`subagent_197.go:358`／`:387`）。**该句今天仍成立。**

### 5.3 本腿**没有**复核的既有判语（不重复判、也不背书）

AC#4（`- [x]`，证据 `docs/evidence/s1/197-blocked-carrier-r4.md`）与 AC#5（`- [x]`，`docs/evidence/s1/197-ac5-selfapproval-v1.md`）本腿**不重裁**；`197-v1` 那四条附条件（C1 `spawnText` 恒空等）今天是否已修，本腿**没有**重跑（不在本件射程，且属 `197-r3` 写腿的地界）。

---

## 6. 终态自证

### 6.1 本腿的写面（只这三类）

1. 本件 `docs/evidence/s1/197-v2-adjudication.md`；
2. `.scratch/wisp/probes/197/v2/**`（`mutate.ps1`／`mutate-m5.ps1`／5 发突变副本／`overlay.json`／`run.txt`）；
3. 一条 commit（显式 pathspec，只含上面两类）。

### 6.2 产码零改（终态尺）

```
git status --porcelain -- internal cmd tools
 M cmd/wisp/run.go                      ← 非本腿（票 451 写腿在飞）
?? cmd/wisp/agent_profiles_451.go       ← 非本腿（同上）
?? cmd/wisp/agent_profiles_451_test.go  ← 非本腿（同上）
?? tools/d22scan/check_blob_417.sh      ← 非本腿（来源未核，不属本件射程）
?? tools/jarvis-export/.verify.js       ← 非本腿（同上）
```
⇒ `internal/`／`internal/panel`／`internal/tools` **零改动**；`cmd/wisp` 的三处都是**票 451 写腿**的活，本腿只读（本件引用的 `run.go` 行号 826／1211／1418 在 HEAD 与工作树上同号——他腿那 10 行插在 `:2113` 之后，`git diff -U0` 只有 `@@ -2112,0 +2113,10 @@` 一个 hunk）。

### 6.3 门禁

- 本腿**没跑** `go vet`／`gofumpt`／`d22scan`：本腿不改产码，跑它们是替别人背书，且 `d22scan` 会把 `.scratch/**` 里的突变体扫进去（那些文件以 `.` 开头，go 工具链不编译它们，但 d22scan 的射程未必同宽）⇒ **明写不跑**，不拿"clean"当本件的读数。

---

## 7. 复核口令（下一位一条条重跑）

```powershell
# 1) 五格判语（本件 §0 那张表）
git show HEAD:docs/evidence/s1/197-v2-adjudication.md | Select-String -Pattern '^\| \*\*AC#'
# 2) 起手基线（§4）
go test ./internal/panel -run 'TestNPlusOneSubagentsEachGetTheirOwnStream|TestAStatusOutsideD43NeverReachesTheWire' -count=1 -v
go test ./internal/tools -run 'Test197RowExistsBeforeFirstChildModelCall|Test197SpawnPublishesIdentityRow|Test197CancelIsPerRowAndNeverCascades|Test229FinishWatcherIsRosterNamedWhileChildRuns' -count=1 -v
# 3) 五发突变（§3）——overlay，不改盘上字节
pwsh -NoProfile -File .scratch/wisp/probes/197/v2/mutate.ps1
pwsh -NoProfile -File .scratch/wisp/probes/197/v2/mutate-m5.ps1
# 4) 产码零改（§6.2）
git status --porcelain -- internal cmd tools
# 5) 撤销口令：本件只增不改——撤掉它＝
#    git revert --no-edit <本件 commit>   （只会删本件与 probes/197/v2 下本腿新增的文件）
```
