# 219 AC#6 收件（非实现者裁决腿 `219-ac6-r1`）

- 起点 HEAD：`d0639eedf`（2026-11-15 13:2x +08 现读）。判据来源＝票 219 全文（`.scratch/wisp/issues/219-approval-card-three-reply-buttons-and-a-reason-box.md`，下简称 `219:NN`）。
- 本腿身份：**非实现者**（裁决/收件，零产码）。写面只 `docs/evidence/s1/219-ac6-receipt.md` 与 `.scratch/wisp/probes/219/**`。
- 台面：一发现跑（`go test ./internal/panel`、`go test ./internal/agent/approval` 单枚集）、`frontend` 的 `render:roster`（真树，绿）、以及**副本树**里的四处前端突变与三处 Go `-overlay` 突变（`frontend/**` 与 `internal/**`、`cmd/**` **零写入**，见 §5 的 pristine 读数）。
- 全部读数现跑件：`.scratch/wisp/probes/219/readings/00-run-all.ps1`（脚本）＋`00-run-all.out.txt`（原样输出，17,993 字节）。本件的每条可核命令都能单独重跑，不必重跑整份脚本。

> **A1875 的尺**（`docs/reports/pending-and-issues.md:37327-37328` 逐字）：「一个能力要活着，必须同时有 ①生产赋值/派发点 ②一枚会因摘掉它而变红的测试。**只有 ①没有 ② ⇒ 随时会被删掉没人知；只有 ②没有 ① ⇒ 就是今天这两个洞**」；「凡判满足的格，必须同时给出**生产调用/派发点 `file:line`** 与**会红的测试名**；两者缺一，判语只能是部分满足」。
> ⇒ 本件每判"满足"的格都带这两件；**唯一缺的那件是 `cmd/wisp` 端到端那一发**（写腿 452 在用该包，本腿 ⛔ 跑），一律写〔未跑包〕并注明所引的是谁的读数。

---

## 0. 首答四件（回话第一行的内容都在这里）

| 项 | 读数 |
|---|---|
| commit | 起点 HEAD `d0639eedfe5539e78bfedfe59a83ce19bd20293c`；**本件落笔提交**由编排者现取：`git log -1 --format=%H -- docs/evidence/s1/219-ac6-receipt.md`（本腿不在文内自引一枚尚未产生的哈希） |
| **L1/L2 真的分表？** | **是**——底层两套表、两扇读门、两条生产读者、只在一处 join；**且两支都修了**（无"没修、留在票 197/211"可写）。**界面层不分**（同一枚「等批准」，wire 只带一个布尔、无 L1/L2 判别位）——这一层具名写进 §1.5，不冒充底层的"是" |
| **AC#6 判语** | **满足**（判据面：L1 短窗口里等的任务在名册上**不再**恒读"没被卡住"；生产点 `cmd/wisp/run.go:1215`＋`internal/panel/subagent_roster_197.go:248-252`；会红测试 `TestAL1WindowInWaitingLightsItsRow`（**基线绿；本腿 M1 突变下现跑红**）、`TestRunPacketLightsTheRosterRowOfItsLiveL1Window` 〔未跑包·引 220-v1 的 MU-D 读数〕） |
| **三数落点** | **3 处**（名册池号 1 处／行级流截断 1 处／整卷流截断 1 处），逐处 `file:line`＋原文见 §2 |

---

## 1. L1 / L2 分表：两套表各自在哪、是否真的分开

### 1.1 两套表各自的 file:line（现跑，`00-run-all.out.txt` §A）

**L1 ＝ 短确认窗口表 `Gate.windows`（内存、按 corr 建键、2–3 秒）**

| 动作 | 现读 |
|---|---|
| 声明 | `internal/agent/approval/gate.go:108` — `windows  map[string]*window` |
| 构造 | `gate.go:216` — `windows:  map[string]*window{}` |
| 开窗（撞名改写＋写入） | `gate.go:448` — `if _, clash := g.windows[w.corr]; clash {`；`gate.go:451` — `g.windows[w.corr] = w` |
| 关窗 | `gate.go:457` — `delete(g.windows, corr)` |
| 按名查（Veto） | `gate.go:548` — `w := g.windows[v.CorrelationID]` |
| **只读枚举口** | `internal/agent/approval/window_read.go:63` — `func (g *Gate) LiveL1Windows() []L1Window`；循环在 `:68-69` — `out := make([]L1Window, 0, len(g.windows))` / `for _, w := range g.windows {` |

**L2 ＝ 待决审批队列表 `Queue.pending`（FIFO、有上限、有答复面）**

| 动作 | 现读 |
|---|---|
| 声明 | `internal/agent/approval/queue.go:84` — `pending []*qitem` |
| 满则拒 | `queue.go:162-163` — `if len(q.pending) >= q.maxPend {` ⇒ fail-closed 拒绝 |
| 收编（撞名改写） | `queue.go:170-171` — `if _, clash := q.byID[corr]; clash {` / `corr = fmt.Sprintf("%s#%d", d.CorrelationID, q.seq)` |
| 入表 | `queue.go:181` — `q.pending = append(q.pending, it)` |
| 出表 | `queue.go:295`／`:306` — `for i, p := range q.pending {`（查拒绝／drop） |
| **只读枚举口** | `internal/agent/approval/pending_read.go:162` — `func (q *Queue) LiveApprovals() []LiveApproval`；循环在 `:168-169` — `out := make([]LiveApproval, 0, len(q.pending))` / `for _, it := range q.pending {` |

### 1.2 两扇门各自的生产调用点（现跑计数，非注释）

尺：`grep -rn '\.LiveL1Windows()\|\.LiveApprovals()' --include='*.go' internal cmd | grep -v _test.go`（`00-run-all.out.txt` §B 原样）：

```
cmd\wisp\approval_reply.go:646: wins := s.gate.LiveL1Windows()
cmd\wisp\panel_pump.go:82: items := g.Queue().LiveApprovals()
cmd\wisp\panel_pump.go:128: wins := rt.gate.LiveL1Windows()
PRODUCTION CALLSITES: LiveL1Windows=2  LiveApprovals=1
```

装配根（读，未跑）：
- `cmd/wisp/run.go:1178` — `Verdicts:  rt.liveVerdicts,`（L2 那一扇，`panel_pump.go:50-55`→`:78-82`）
- `cmd/wisp/run.go:1215` — `L1Windows: rt.liveL1WindowWaits,`（L1 那一扇，`panel_pump.go:124-137`）

### 1.3 两扇门**不合流**：L1 的窗口没有进 L2 的待决区（结构＋测试双向）

- 读的原文：`pending_read.go:168-169` 只遍历 `q.pending`；`window_read.go:68-69` 只遍历 `g.windows`。**没有任何代码把 L1 窗口塞进 `q.pending`**（`LiveApprovals` 的 out 只由 `q.pending` 追加）。
- pump 也是两条：`internal/panel/pump.go:353-360` 的 `cards` 只来自 `p.src.Verdicts()`；`:460-464` 的 `waits` 来自 `p.src.L1Windows()`，两者**并列**传进 `taskRosterSectionFrom(..., cards, waits)`：
  ```
  460:		var waits []L1WindowWait
  461:		if p.src.L1Windows != nil {
  462:			waits = p.src.L1Windows()
  463:		}
  464:		snap.Tasks = taskRosterSectionFrom(p.src.Tasks(), cards, waits)
  ```
- **测试把"不合流"钉住了**（`go test ./internal/agent/approval -run ...` 本腿现跑，4/4 PASS）：
  - `TestLiveL1WindowsSeesAWindowInTheLoop`（`ticket220_l1_window_read_test.go:248`）在 L1 窗口开着时断言 `LiveApprovals()` 为空（`:263-264` 逐字：`if items := g.Queue().LiveApprovals(); len(items) != 0 { ... "an L1 window ..."`）。
  - `TestReadingL1WindowsMovesNoGateState`（`:296`，同一空断言在 `:311`）、`TestStartedCallsAreNotReportedAsWaiting`（`:380`，`:388`）、`TestL1WindowIsPlainDataWithNoRouteToAnAnswer`（`:338`）。

⇒ **判：两套表"真的分开"＝是**。分开是**结构性的**（两个类型、两把锁、两个文件、两枚接收者），不靠命名约定；`L1Window` 只带三个 id 字符串、没有通往答复的路径（`window_read.go:43-47` 的字段表＋`:18-26` 的"不是答复"段）。

### 1.4 汇合点只有一处，而且是"加一维"不是"合表"

`internal/panel/subagent_roster_197.go:244-253`（现读）：

```
244:	waiting := make(map[string]bool, len(pending)+2*len(l1Windows))
245:	for _, card := range pending {
246:		indexWaitingKey(waiting, known, card.CorrelationID)
247:	}
248:	for _, w := range l1Windows {
251:		indexWaitingKey(waiting, known, w.CorrelationID)
252:		indexWaitingKey(waiting, known, w.TaskID)
253:	}
```
落格：`:268` — `BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID],`（布尔，**另有**行自己的 D43 `Status`/`StatusKnown` 原样随行，见 `:264-272`）。

- 会红的测试：`TestAL1WindowInWaitingLightsItsRow`（`internal/panel/subagent_blocked_220_test.go:158`，L1 那一支）与 `TestASecondCardForTheSameTaskStillLightsItsRow`（`:107`，L2 撞名那一支）——**本腿两发突变各自只打红自己那一支**（§1.6）。
- "不知道就不编"那一半也在：同一用例 `:171-184` 用**没有 `L1Windows` 读者**的泵发一包，断言同一枚窗口读成 `false`（`assertBlocked(t, data, "child-1", false)`）——即"装配不全"读成"没读数"，不是"没人等"。

### 1.5 分层的判语：底层＝是；**界面层＝不分**（具名，不洗）

- 界面上一条被卡住的行只有一个字：`frontend/src/components/harness/monitor-popover.tsx:91` — `if (row.blockedOnApproval) return { label: "等批准", dot: "blocked" };`。L1 与 L2 **落成同一枚字、同一个 `bg-warn` 点**。
- wire 上也只有布尔：`subagent_roster_197.go:137` — `BlockedOnApproval bool `json:"blockedOnApproval"``；`frontend/src/lib/panel.ts:213` — `blockedOnApproval: boolean;`。**没有任何 L1/L2 判别位**（`L1WindowWait` 只在 Go 侧 join 时用，不上 wire：`pump.go:135` 的类型、`panel_pump.go:134` 的折叠）。
- 顺带一枚**现存不对称**（读数，不是指控）：页头那枚「待审批 N」只数 L2 卡（`App.tsx:224-228`，`snapshot.pending.length`），而名册行会被 L1 点亮 ⇒ **一枚只在 L1 窗口里等的任务，行上说「等批准」、页头计数不动**。
- ⇒ **AC#6 的文字不要求界面分栏**（`219:75` 的判据是"名册上读成没被卡住"），所以本件**不因此判 AC#6 不成立**；但若票面原意是"界面上也要分开"，**那一半没修**，落点是"wire 加一个 level/scope 位＋行上两个字"，属**界面地界**（本编队零写面），具名留后续票。**⛔ 不写"留在票 197/211"**——`219:76` 自己已把这支更正成票 220，而 220 已落地（五格清零，2026-11-12 收口）。

### 1.6 分表的**敏感度**：两支各有一发只打自己的突变（本腿现跑，`-overlay` 只改副本）

| 突变 | 改了什么（副本） | 红名（逐字） | 没红的 |
|---|---|---|---|
| **M1** | 摘掉 `subagent_roster_197.go:248-252` 那 5 行（L1 那一支＝pre-220 形） | `TestAL1WindowInWaitingLightsItsRow`：`subagent_blocked_220_test.go:160` `row "child-1" blockedOnApproval=false，期望 true`、`:166` 同句；`TestABlockedRowKeepsItsRunningDimension`：`:208` `root row = {... BlockedOnApproval:false ...}, want blocked AND still carrying its filed state` | **L2 那两支全 PASS**（`TestASecondCardForTheSameTaskStillLightsItsRow`、`TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt`） |
| **M2** | 摘掉 `:245-247` 那 3 行（L2 那一支） | `TestASecondCardForTheSameTaskStillLightsItsRow`：`subagent_blocked_220_test.go:109`、`:113`；`TestTheClashFormIsRecognisedOnlyAsTheQueueWritesIt`：`:147` ×3 子测试（`child-1#9`／`#147`／`#retry`）；`TestABlockedRowKeepsItsRunningDimension`：`:199` `the second card is pending and the row does not say so` | **L1 那一支 PASS**（`TestAL1WindowInWaitingLightsItsRow`） |

⇒ 这是"分表"这个词最硬的读数：**两支各有独立牙齿，摘一支不会连带另一支**。M1/M2 都用 `go test -overlay`（overlay JSON 在 `.scratch/wisp/probes/219/mut/m*-overlay.json`），盘上产码零字节动（§5）。

---

## 2. 三数落点（3 处）：名册池号 / 行级流截断 / 整卷流截断

> 尺＝从 Go 生产点追到像素，逐段给 `file:line`＋原文；每处按 A1875 给两件套。**"三数"是三个维度**（池号、行级截断、整卷截断），不是三个字段。

### 2.1 落点一：名册池号「在跑 N / 上限 M」

| 段 | 现读 |
|---|---|
| **Go 生产点**（宿主自己的计数，不是数行） | `cmd/wisp/panel_pump.go:259` — `state.InFlightSlots = rt.tasks.InFlightSubagents()`；`:263` — `state.PoolCap = tools.MaxConcurrentSubagents`（函数体 `:207 taskRosterState`） |
| 载体 | `internal/panel/subagent_roster_197.go:162-163` — `InFlightSlots int `json:"inFlightSlots"`` / `PoolCap int `json:"poolCap"``；搬运 `:225-226` |
| wire 类型（页面侧） | `frontend/src/lib/panel.ts:226-227` — `inFlightSlots: number;` / `poolCap: number;` |
| **喂数点** | `frontend/src/App.tsx:266-267` — `poolInFlight={snapshot.tasks?.inFlightSlots}` / `poolCap={snapshot.tasks?.poolCap}` |
| **界面原文** | `frontend/src/components/harness/monitor-popover.tsx:239-246` — 组标题右侧 `<span ... title="名册池：宿主正占着的槽位与上限">` 里 `在跑 {poolInFlight} / 上限 {poolCap}` |
| **会红的测试名** | Go 载体层：`TestTheRosterReaderPutsSubagentsOnTheWire`（`internal/panel/subagent_roster_197_test.go:103`）——M3 下红于 `:174` `inFlightSlots = 0, want the slot the pool is holding`、`:177` `poolCap = 0, want the cap the host handed over (4) carried verbatim`；**页面层**：`render:roster` 的齿 `the pool numbers must read 在跑 3 / 上限 8`（I1 下红）＋源级齿 `the monitor must be fed the pool numbers straight off snapshot.tasks`（I2 下红） |
| 缺什么 | `cmd/wisp` 端到端那一发〔未跑包〕：`TestRunPacketCarriesTheSubagentItsRosterRowFed`（`cmd/wisp/subagent_carrier_197_test.go:232`，断言 `:331` `poolCap != tools.MaxConcurrentSubagents`、`:360` `inFlightSlots != 0 after the join`）本腿未跑，理由见 §4。另：池号在无 tasks 键时**整块不画**（`monitor-popover.tsx:239` 的 `!== undefined` 双条件），这是刻意的"不造假 0"，齿 `with no tasks key nothing pool-shaped may render` 钉着 |

### 2.2 落点二：行级流截断「省 N 字」/「流已丢」

| 段 | 现读 |
|---|---|
| **Go 生产点** | `cmd/wisp/panel_pump.go:233` — `trunc := rt.stream.TruncationFor(key)`；`:243-245` — `StreamTruncated: trunc.Truncated,` / `StreamElidedRunes: trunc.ElidedRunes,` / `StreamDropped: trunc.Dropped,` |
| 载体 | `internal/panel/subagent_roster_197.go:140-142` — `StreamTruncated bool `json:"streamTruncated"`` / `StreamElidedRunes int `json:"streamElidedRunes"`` / `StreamDropped bool `json:"streamDropped"``；搬运 `:269-271`（`:271` 还有 `row.StreamDropped \|\| dropped[row.StreamKey]`） |
| wire 类型 | `frontend/src/lib/panel.ts:214-216` |
| **喂数点** | `frontend/src/App.tsx:102-104` — 派生器 `monitorAgents`：`streamTruncated: row.streamTruncated,` / `streamElidedRunes: row.streamElidedRunes,` / `streamDropped: row.streamDropped,` |
| **界面原文** | `monitor-popover.tsx:112-116` — `rowBoundLabel`：`if (row.streamDropped) return "流已丢";` / ``if (row.streamTruncated) return `省 ${row.streamElidedRunes ?? 0} 字`;``（`null`＝整行干净）；画在行尾 `:263-270` |
| **会红的测试名** | Go：`TestTruncationFactsRideThePacket`（`internal/panel/subagent_roster_197_test.go:331`）——M3 下红于 `:400` `big row on the wire = elided 0/truncated false, want 88/true`、`:408` `the bytes do not state the loss as a number`；页面：`render:roster` 齿 `a folded row must show its elided count beside the status word`、`a dropped row must say its page is gone, not a fake elision count`（I1 下红）、`the derivation must carry the row stream triple verbatim`（I2 下红） |
| 缺什么 | —（行级这一处四段齐；`cmd/wisp` 端到端同 §4〔未跑包〕，`subagent_carrier_197_test.go:586-589` 是它的那枚断言） |

### 2.3 落点三：整卷流截断「流已截断：省略 N 字，已丢 M 条」

| 段 | 现读 |
|---|---|
| **Go 生产点** | `cmd/wisp/panel_pump.go:264-269` — `state.StreamTruncated = rt.stream.Truncated()` / `:265` `ElidedRunes()` / `:266` `DroppedKeys()` / `:269` `DroppedCount()` |
| 载体 | `internal/panel/subagent_roster_197.go:166-183`（`:166` `StreamTruncated bool `json:"streamTruncated"``、`:168` `streamElidedRunes`、`:176` `droppedStreamKeys`、`:183` `droppedCount`）；搬运 `:227-230` |
| wire 类型 | `frontend/src/lib/panel.ts:228-232`（`:231-232` 带"Ring 64 overflow"注释） |
| **喂数点** | `frontend/src/App.tsx:268-270` — `streamTruncated={snapshot.tasks?.streamTruncated}` / `streamElidedRunes={snapshot.tasks?.streamElidedRunes}` / `streamDroppedCount={snapshot.tasks?.droppedCount}` |
| **界面原文** | `monitor-popover.tsx:122-134` — `sectionBoundLine`：`return \`流已截断：省略 ${elidedRunes ?? 0} 字，已丢 ${dropped} 条\`;` / `` `流已截断：省略 N 字` `` / `` `已丢 ${dropped} 条流` ``，三档都只在真截断时返回；画在行列表尾 `:284-286` |
| **会红的测试名** | Go：同一枚 `TestTruncationFactsRideThePacket`——M3 下红于 `:397` `section total = 0, want the log's own 88`、`:408`；页面：`render:roster` 齿 `a truncated log with drops must name both losses in one line`、`a truncated log must show the elision count`、`dropped streams must warn even when nothing was elided`（I1 下红）、`the monitor must be fed the whole-log truncation triple straight off snapshot.tasks`（I2 下红） |
| 缺什么 | —（同上；干净日志零噪音那一半也有齿：`a clean log must render no truncation warning`） |

### 2.4 三数落点的计数（回答首答的"N 处"）

**3 处**，口径＝**真上屏的独立显示位**：①子智能体组标题右侧池号；②行尾截断标（行级）；③行列表尾整卷告示（整卷）。类型声明（`lib/panel.ts:213-232`）**不计**落点；`App.tsx` 的喂数是落点的**上游**，也不计；截图证据各自另在 `docs/evidence/s1/219-ac6-roster-pool-truncation-live.png`（写腿 11-08 的活体，本腿不重拍）。

---

## 3. 逐格判语（五列）

> 五列＝**AC 原文／判语／依据／可核命令／缺什么**。"依据"只写现读，不写票面自述。

| # | AC 原文（逐字） | 判语 | 依据（现读） | 可核命令 | 缺什么 |
|---|---|---|---|---|---|
| ① | `219:75`「**L1 与 L2 分表**：L1 短窗口里等的任务在名册上今天恒读"没被卡住"（账 `A416` 那一格）⇒ 本票要么把它一起修掉、要么具名写"这一支没修、留在票 197/211"」 | **满足** | L1 那一支已由票 220 一起修掉（不是留洞）：表 `gate.go:108`／门 `window_read.go:63`／生产读者 `panel_pump.go:128`／装配 `run.go:1215`／join `subagent_roster_197.go:248-253`／落格 `:268`。票面 `219:76` 自己把这支更正成票 220，票 220 五格已清零（`220-*.md:70` 的 2026-11-12 收口条） | `pwsh -File .scratch/wisp/probes/219/readings/00-run-all.ps1` §A–§E；单跑 `go test ./internal/panel -run TestAL1WindowInWaitingLightsItsRow` | `cmd/wisp` 端到端那发〔未跑包〕（§4.1） |
| ② | 同上格里的"分表"二字（结构性判据） | **满足** | 两套表＝两个类型/两把锁/两个文件；两扇门不合流（`pending_read.go:168-169` 只走 `q.pending`；`window_read.go:68-69` 只走 `g.windows`）；空断言在 `ticket220_l1_window_read_test.go:263-264`／`:311`／`:388`；两支各有一发只打自己的突变（§1.6 M1/M2） | 00-run-all §A–§D＋`-overlay` 两发 | — |
| ③ | 界面层是否也要分（"分表"在像素上的读法） | **不成立，也不在 AC#6 文字里**（具名留洞） | 一行只有一个字：`monitor-popover.tsx:91`「等批准」；wire 只有布尔（`subagent_roster_197.go:137`／`lib/panel.ts:213`），**无 L1/L2 判别位**；页头「待审批 N」只数 L2（`App.tsx:224-228`）而名册行被 L1 点亮 ⇒ 一处不对称 | `grep -rn 'blockedOnApproval' frontend/src internal/panel/subagent_roster_197.go` | 若原意是界面分栏：wire 加 level/scope 位＋行上两字，属界面地界，本编队零写面 ⇒ 具名交后续票（⛔ 不写"留在票 197/211"） |
| ④ | 名册池号上屏（`219:107` 写腿自述的那一处落点） | **满足** | 生产点 `panel_pump.go:259/:263`→载体 `:162-163/:225-226`→`panel.ts:226-227`→`App.tsx:266-267`→`monitor-popover.tsx:239-246`「在跑 N / 上限 M」 | 00-run-all §G（M3 红）＋§H（真树 `render:roster` 绿）＋§I1/§I2 红 | `cmd/wisp` 端到端〔未跑包〕 |
| ⑤ | 行级流截断三数上屏（同上） | **满足** | `panel_pump.go:233/:243-245`→`:140-142/:269-271`→`panel.ts:214-216`→`App.tsx:102-104`→`monitor-popover.tsx:112-116/:263-270`「省 N 字」「流已丢」 | 00-run-all §G（M3 红 `:400`）＋§I1（两齿红）＋§I2 | `cmd/wisp` 端到端〔未跑包〕 |
| ⑥ | 整卷流截断三数上屏（同上） | **满足** | `panel_pump.go:264-269`→`:166-183/:227-230`→`panel.ts:228-232`→`App.tsx:268-270`→`monitor-popover.tsx:122-134/:284-286`「流已截断：省略 N 字，已丢 M 条」 | 00-run-all §G（M3 红 `:397`）＋§I1（三齿红）＋§I2 | `cmd/wisp` 端到端〔未跑包〕 |
| ⑦ | A1875 的元判据：每格是否"生产点＋会红测试"两件齐 | **齐（6/6 格）**，唯一例外是 `cmd/wisp` 端到端那一发 | 上表 1–6 格每格都现跑出红名；`cmd/wisp` 那一发按本腿禁令未跑（§4.1），其红名由 220-v1 的 MU-D 读数背书（`docs/evidence/s1/220-v1-verdict.md` 引），本件如实标〔未跑包〕 | 00-run-all 全文（17,993 字节原样输出） | 见 §4.1（`cmd/wisp` 归写腿 452） |

**AC#6 框**：本件只出判语，**不翻框**（裁决者≠实现者；翻框与 `Status:` 归编排者）。

---

## 4. 边界与未跑（⛔ 拿"未跑"当"没红"）

### 4.1 `cmd/wisp` 全包：**只读，未跑**〔未跑包〕

派单与本项目铁律都写明 `cmd/wisp/**` 今归写腿 452，本腿 ⛔ 跑 `go test ./cmd/wisp`（重包同一时刻只许一个）。⇒ 下列两发**不是本腿读数**，只作**引证**并逐条标〔未跑包〕：
- `TestRunPacketLightsTheRosterRowOfItsLiveL1Window`（`cmd/wisp/panel_l1_blocked_220_test.go:183`）——**AC#6 的端到端那一发**。它的去牙红由 `220-v1` 记录：摘掉 `run.go:731` 那一行（当时号）⇒ 红句落在装配格（`panel_l1_blocked_220_test.go:216` `PRE-CHANGE / WIRING READING: ... its roster row still read blockedOnApproval=false ...`）。本腿只做了**只读计数**：`run.go:1215` 那行在盘上（§1.2）。
- `TestRunPacketCarriesTheSubagentItsRosterRowFed`／`…ReportsTheStreamLogPastItsBound`／`…CarriesTheStreamDroppedCountPastTheHardCeiling`（`cmd/wisp/subagent_carrier_197_test.go:232/:525/:611`）——三数在真 run 上的那三枚。
⇒ **本件对 AC#6 的"满足"建立在**：①只读 grep 得到的两条生产链完整；②`internal/panel` 与 `internal/agent/approval` 上**我现跑红**的突变读数；③`render:roster` 在真树上现跑绿＋四处突变各打各的齿。**端到端那一发留给能跑 `cmd/wisp` 的腿/编排者复跑**——这是本件唯一的一件"缺"。

### 4.2 `frontend/**` 零写：突变的做法

- 真树只跑一次 `npm run render:roster`（绿，`rc=0`），产物落 `frontend/node_modules/.tmp/rosterrender`（gitignored 构建缓存，非源）。
- 四处前端突变**在副本树**里跑：`.scratch/wisp/probes/219/mut/fe/`（`frontend/src`＋`frontend/scripts`＋六枚配置件各一份拷贝；`node_modules` 是指向 `frontend/node_modules` 的 **junction**，`git check-ignore` 命中根 `.gitignore:26 node_modules/`）。每次突变后**还原**，`Reset-Fe` 在脚本里。
- 突变读数：
  - **I0** 副本基线 **OK**（先证副本本身是绿的，免得把"副本坏了"读成"齿敏感"）；
  - **I1**（摘池号 span＋行尾标＋整卷行）⇒ 红 6 条，逐字见 §2.1–2.3；
  - **I2**（摘 `App.tsx` 五条喂数＋派生器三行）⇒ 红 5 条（含三条源级齿）；
  - **I3**（摘 `monitor-popover.tsx:91` 的等批准档）⇒ 红 3 条；
  - **I4**（摘 `App.tsx:98` 的派生）⇒ 红 3 条。

### 4.3 一枚**齿自己的缺陷**（本腿顺手量出来的，写在这里免得下一位把绿当承重）

`frontend/scripts/render-task-roster.tsx:189` 的正向齿 `the blocked row must show 等批准` **在 I3 突变下仍然绿**——因为夹具那枚行的 label 本身是「**等批准**的子任务」（`:157`，`:188` 也断言该 label 上屏）⇒ 这条齿被自己的夹具满足，**不承重**。承重的是同一组的负向齿 `the blocked row must NOT show the D43 status word (等批准 outranks it)`（`:190-193`）、dot 齿 `the blocked row's dot must be the warn tier`（`:194`）、源级齿（`:283-286`）——I3 下这三条才红。
⇒ 建议（⛔ 本腿不改，`frontend/**` 零写面）：把夹具 label 改成不含「等批准」三字的字串，或把正向齿收窄成"行尾状态字 span 的文本等于 等批准"。**这一条不影响 AC#6 判语**（三枚承重齿都在），但影响"这条齿能不能拦住回归"。

### 4.4 临时件

`.scratch/wisp/probes/219/` 下：`readings/`（脚本＋原样输出）、`mut/`（三枚 Go 突变副本＋三枚 overlay JSON＋`fe/` 前端副本树含 `.tmp` 构建产物与 node_modules junction）。按"临时件只建不删"原样留盘；`fe/` 与 `mut/*.go` 都**不进**本腿的提交（见 §5 的 pathspec），只作可复现件存在。

---

## 5. pristine 与收件

本腿对产码的写入＝**零**。现跑（`00-run-all.out.txt` §J 原样）：

```
$ git status --porcelain=v1 -- frontend/src internal cmd
(empty above = zero writes)
```

提交面（显式 pathspec，只 commit 不 push）：`docs/evidence/s1/219-ac6-receipt.md` ＋ `.scratch/wisp/probes/219/readings/**` ＋ `.scratch/wisp/probes/219/mut/*.json` ＋ `.scratch/wisp/probes/219/mut/m*/**`（前端副本树与其构建产物不入库）。

**撤销口令**：`撤 219-ac6-r1`——删 `docs/evidence/s1/219-ac6-receipt.md` 与 `.scratch/wisp/probes/219/`，票面与产码本腿一字未动，无需复原。

**给下一位的三句**
1. `cmd/wisp` 一空出来，先跑 `TestRunPacketLightsTheRosterRowOfItsLiveL1Window` 并把 `--- PASS/FAIL` 记进本件 §4.1——那是 AC#6 唯一还欠的读数。
2. 界面要不要分 L1/L2（§3 ③）是一个**产品/契约面问题**：wire 要加位、C17 白名单不动（`panel.*` 新名禁止），按"改契约文字"那一类摆机主，⛔ 由腿自行决定。
3. §4.3 那枚齿的夹具字串，下一位碰 `frontend/scripts/render-task-roster.tsx` 时顺手修；⛔ 在本件名下改。
