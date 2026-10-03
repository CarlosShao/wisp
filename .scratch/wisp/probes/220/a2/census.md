# 票 220 只读普查 a2 — 「`L1Windows` 这一跳今天断在哪、下一跳要动哪几枚文件」

## §0 起手（现跑）

- 采集起点时刻：`2026-10-03T10:33:43+08:00`（`date -Iseconds` 现跑）
- 起手锚点：`git log -1 --format=%H` = `e16e303773f59be81d2f02eee3cc638870aa1630`（分支 `dev`）
- **采集期间锚点前移**：第二次现跑＝`1289d8bbe0b90c2f58854c2b58e21d9f518ec62d`（`10:38:41 +0800`，`docs(probes/252-r2)` 骨架），
  本文所有行号在**两枚锚点之间被复量过一轮**，复量结果逐格写在 §8。
- 票面：`.scratch/wisp/issues/220-roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops.md`（全文读过）
- 前一枚同票普查：`.scratch/wisp/probes/220/c1/census.md`（锚 `6233dedc`，2026-09-29，**在 220-r1 之前**，其结论已被两枚写腿部分追平）；
  写腿交件 `.scratch/wisp/probes/220/r1/impl.md` §8（本腿量的就是它留下的那一跳）。
- 性质：**只读普查**。零编译零测试零写码；⛔ 本腿没跑过任何 Go 命令（闸门见 §7）。
- `frontend/**`／`design/**`：未读、未引、结论不转述其内容。所有 grep 均写显式搜索根
  （`cmd internal tools docs scripts`），无一根用 `.`。
- 本腿唯一写面＝本文件；骨架第一轮 commit＝`3bc10b65`。

## §1 脏面名册（起手 `git status --porcelain` 现抄）

起手时 `cmd/wisp` 确有一枚未提交改动，**采集期间它被 255-r2 自己提交了**，两枚读数都留在下面。

| 时刻 | 现跑读数 |
|---|---|
| 起手 `10:33:43` | 全仓脏名册里 Go 产码只有两枚：`M cmd/wisp/config_reload.go` ＋ `?? cmd/wisp/restart_tier_keys_255r2_test.go`；其余是 `.scratch/**` 证据件与 `design/**` 的 31 条删除／改动 |
| 复量 `10:39+` | `git status --porcelain -- cmd/wisp` ＝ **空**；两枚已由 `67ab595d`（`10:37:18 +0800`，`255-r2（编排者代提·死腿尾部那一格落定）`）入库，该 commit 的 stat 恰含 `cmd/wisp/config_reload.go`（＋8/−5 行段）与 `cmd/wisp/restart_tier_keys_255r2_test.go`（＋136） |

⇒ **本腿读 `cmd/wisp` 时唯一需要标注的中间态窗口＝`10:33:43`–`10:37:18`**，那几分钟里 `config_reload.go` 是 255-r2 的工作树版本。
本腿对 `cmd/wisp` 的全部引用落在 `run.go`／`panel_pump.go`／`panel_config_store.go`／`resident_*.go`／`panel_host_windows.go`
／`panel_resident_windows.go` 与若干 `_test.go`，**没有一枚落在 `config_reload.go`** ⇒ 行号不受那个窗口影响
（`run.go` 在该窗口内零改动，由上面第二行读数直接证）。

仍在飞的写腿（编排者令里点名的三枚，本腿一律未碰其写面）：`255-r2`（cmd/wisp＋internal/config）、
`171-r3`（门钉尺）、`252-r2`（internal/tools；其骨架件 `1289d8bb` 就落在这几分钟）。

## §2 待验断言（编排者令里写下的每一句，全部先当假的）

| # | 断言 | 复量结论 |
|---|---|---|
| F1 | `ce693ce4` 交回 `window_read.go`（`L1Window` 三枚 string／零方法＋`(*Gate).LiveL1Windows()`） | 成立，§8 F1 |
| F2 | `480b970d` 交回 `L1WindowWait`＋`PumpSources.L1Windows` 那一枚 reader 位 | 成立，§8 F2 |
| F3 | 「`cmd/wisp` 里没人给 `PumpSources.L1Windows` 供数据」 | 成立且**偏轻**，§8 F3 |
| F4 | 「L1 窗口的开／关缺发布触发」 | **半不成立 ⇒ 具名推翻**，§8 F4 ＋ §4 |
| F5 | 「`git_test.go` 一带有 `want-4` 那族锚，只认 `bridge.go` 里恰有四枚 `panel.*`」 | 成立但要补一条边界，§8 F5 ＋ §6 |
| F6 | 「票 256：`approval.New` 建得太早」会不会同样卡住 A 这一跳 | 断言成立、**与本跳无因果**，§8 F6 ＋ §5 |
| F7 | 「`cmd/wisp` 此刻是脏面（255-r2 未提交）」 | **已被时间推翻**（非谎，是过期），§1 |
| F8 | 票面「现量」那张表的行号 | 已漂移（票自己写了「别信这里的行号」），逐行对照见 §8 F8 |
| F9 | 「基线那 4 枚 FAIL 根因＝别腿在 `design/**` 的工作树删除」 | **只有一半成立**，§8 F9 |
| F10 | 「`markStarted` 之后的调用不会被 `LiveL1Windows` 报成在等」（`window_read.go:59-62` 的自述） | **代码里没有这一条**，§8 F10 ＋ §4 |

## §3 问题 A：供给链逐跳

### 3.1 `PumpSources` 今天有哪几枚 reader 位、生产里各自由谁填

全仓 `panel.NewSnapshotPump(` 的**生产调用者＝恰一枚**＝`cmd/wisp/run.go:699`
（本腿现跑：搜索根 `cmd internal tools`，排除 `_test.go`，命中 2 行＝那一枚调用＋`internal/panel/pump.go:229` 的定义本身；
测试侧 29 枚）。逐枚 reader 位（声明行＝`internal/panel/pump.go`）：

| reader 位 | 声明 | 生产填充点（`run.go` 字面量） | 填进去的那枚函数的定义处 |
|---|---|---|---|
| `Verdicts` | pump.go:140 | run.go:700 | cmd/wisp/panel_pump.go:58（`:62` 读 `rt.gate.Queue().LiveApprovals()`） |
| `Mode` | pump.go:142 | run.go:701 | internal/perm/store.go:155 |
| `Workspace` | pump.go:144 | run.go:702 | cmd/wisp/panel_pump.go:83 |
| `Git` | pump.go:150 | run.go:703 | cmd/wisp/panel_pump.go:228 |
| `Model` | pump.go:157 | run.go:704 | cmd/wisp/panel_pump.go:242 |
| `Credential` | pump.go:163 | run.go:710 | cmd/wisp/panel_config_store.go:285 |
| `Results` | pump.go:165 | run.go:711 | internal/panel/pump.go:559 |
| `Instructions` | pump.go:177 | run.go:718 | cmd/wisp/panel_pump.go:108 |
| `Tasks` | pump.go:190 | run.go:724 | cmd/wisp/panel_pump.go:146 |
| **`L1Windows`** | **pump.go:200** | **无——字面量里根本没有这一行**（run.go:699-726 逐枚读回） | 无 |
| `AttachmentMax` | pump.go:205 | 未填（缺省走 pump.go:266-268 的 `MaxAttachmentBytes`，故意） | — |
| `Now` | pump.go:209 | 未填（缺省 `time.Now`，pump.go:299-302，故意） | — |
| `Out` | pump.go:214 | run.go:725 | cmd/wisp/panel_pump.go:319（落台账，⛔ 不是页面） |

⇒ `L1Windows` 缺的是**「没人填」，不是「没有可填的源」**：源今天已经在树上并且是可用形状
（`internal/agent/approval/window_read.go:63` 的 `(*Gate).LiveL1Windows()`），
消费端也已经在树上（`internal/panel/pump.go:323-327` ＋ `internal/panel/subagent_roster_197.go:216-221`）。
断的是中间那一枚**投递**。

### 3.2 三问「谁产出／谁投递／谁落盘」（⛔ 本腿不数以枚数论的显式调用者）

- **产出**＝`internal/agent/approval`：`gate.go:91` 的 `windows map[string]*window`（`window` 结构 `gate.go:99-109`，
  其中 `taskID` 在 `:106`＝ce693ce4 补的那一枚纯数据字段）→ 由 `openWindow`（gate.go:352-360）写、`closeWindow`（:362-366）删 →
  读出＝`window_read.go:63-78`。
- **投递**＝**今天没有人**。本腿现跑（搜索根 `cmd internal tools docs scripts`，字面 `LiveL1Windows|L1WindowWait|L1Windows`）：
  `cmd/` 下**零命中**；`internal/panel/pump.go:115-200,323-327` 是载体与消费；
  `internal/agent/approval/ticket220_l1_window_read_test.go` 与 `internal/panel/subagent_blocked_220_test.go:73-81,151-187` 是测试；
  其余命中全在注释／`.scratch/**` 归档。⇒ **`LiveL1Windows()` 的产码调用者＝零枚**（同 A483 那一族「载具在、没人喂」）。
- **落盘**＝`internal/panel`：`pump.go:327` 把 `waits` 递给 `taskRosterSectionFrom`（定义 `subagent_roster_197.go:198`，
  导出那枚 `TaskRosterSectionFrom` 在 `:187` 以 nil 委托），
  该函数在 `:214` 先把卡的 corr 登记、`:216-221` 把窗口的 corr 与 taskID **两枚都**登记（`indexWaitingKey` 定义在 `:262`，
  认队列改写后缀的 `splitClashSuffix` 在 `:277`），最后一格在 **`:236`**
  （`BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`）。

### 3.3 从 `LiveL1Windows()` 到面板快照的每一跳，三档判定

| 跳 | 内容 | 出处 | 判定 |
|---|---|---|---|
| 1 | `g.windows` 里那枚窗口 | gate.go:91 ／开 :283、:352-360／关 :286、:362-366 | 已有 |
| 2 | 只读枚举口 | window_read.go:63-78 | 已有 |
| 3 | 把 `[]approval.L1Window` 映射成 `[]panel.L1WindowWait` 的那枚宿主 reader（同形先例＝panel_pump.go:58-76） | **cmd/wisp 里不存在** | **缺一整块** |
| 4 | `PumpSources.L1Windows:` 那一行赋值 | run.go:699-726 字面量 | **缺一行赋值** |
| 5 | 泵读 reader（nil ⇒ 不读，绝不读成「没有窗口在等」） | pump.go:323-326 | 已有 |
| 6 | 递给 join | pump.go:327 | 已有 |
| 7 | 窗口两枚 id 进 `waiting` | subagent_roster_197.go:216-221 | 已有 |
| 8 | 行上那一格 | subagent_roster_197.go:236（字段 `:132`，JSON 键 `blockedOnApproval` 不变） | 已有 |
| 9 | 「开／关时重发一发快照」的触发 | 见 §4：开＝控制台形状已有；关＝晚一拍；常驻形状＝**缺一整块** | 分形状而定 |
| 10 | 出口 | pump.go:348-368 → run.go:725 → panel_pump.go:319（**只到台账**） | 已有；⚠ **Go→页面那一跳整块缺**（见 3.4） |

### 3.4 还差的那一跳不止「reader 没接」

编排者那句「面板上那一行今天还看不见在等批准」的原因比 A 这一跳更深一层，本腿现跑为证：
全仓**没有任何 Go→页面推送口**。搜索根 `cmd internal`、字面 `PostMessage|EvaluateScript|WebMessage|PushSnapshot`，
排除 `_test.go` ⇒ 生产侧**零命中**（唯一命中是 `cmd/wisp/testdata/esclistener/main.go:66` 那枚测试小工具用
`PostMessageW` 发 `WM_CLOSE`，与本跳无关）。快照的 `Out` today 只有台账（panel_pump.go:319，
其 :290-318 那段把头注写死了「ledger 装不下整包字节，只书摘要＋sha256」）。
⇒ 落地腿把 3-4 两跳转完之后，**新增的是台账里那一发包裹的字段为真**，界面上那一行仍要等票 33/35 的 transport。
票 220 AC#5 末尾要求具名写清的「这一格修完之后界面上还差哪一跳」＝**的就是这一跳**，它和 `L1Windows` 无关、也不该由本票背。

## §4 问题 B：发布触发

### 4.1 仓里既有的三枚触发口（全在 `cmd/wisp`，这就是「可复用的那一形」）

| 触发 | 发出点 | 挂载点 | 覆盖形状 |
|---|---|---|---|
| T1 卡片上屏 | `consoleApprovalUI.Prompt` 末段 run.go:1386-1388 | run.go:735 `rt.ui.publish = rt.publishPanelSnapshot`（`if rt.ui != nil` 才挂，run.go:734） | **仅控制台形状** |
| T2 卡片消失／交接执行 | `consoleApprovalUI.Update` 里只认 `EventDismissed`／`EventStarted` 两枚 kind：run.go:1414-1418 | 同 T1 | 仅控制台形状 |
| T3 环路事件 | `consoleSink.Publish`：`EvToolStart`／`EvToolEnd`／`EvStopped`… 五臂置 `changed`（run.go:1275-1296），末尾 run.go:1297-1299 发 | run.go:998（`execute()` 里建 sink） | **两枚形状都有**（常驻腿也走 `execute`：resident_task_source_windows.go:433） |

「谁调 T1」＝闸门自己：`gate.go:294`（L1 路）与 `gate.go:536`（L2 路）各调一次 `g.ui.Prompt`。
「谁调 T2」＝`gate.go:322`（否决后 `EventDismissed`）、`gate.go:340`（到期后 `EventStarted`）、
`gate.go:554`／`:573`（L2 两条结束臂）。⚠ `EventTick` 今天**零发射者**（只有声明 `internal/agent/approval/ui.go:64`），
`internal/panel/pump.go:50-55` 那段注的就是这件事——想拿倒计时当触发，没有那个事件源。

### 4.2 「开」这一侧

- **控制台形状：触发已存在，且时机是对的。** `openWindow`（gate.go:283）**先于** `g.ui.Prompt`（gate.go:294）
  ⇒ T1 那一发 publish 读 `g.windows` 时窗口已经在表里。旁证不是注释而是测试自己的话：
  `cmd/wisp/panel_pump_test.go:198-206` 写「the FIRST was booked by consoleApprovalUI.Prompt, i.e. by the run itself」，
  并要求台账里 `depth=1` 的记录 ≥2 枚。
- **常驻（GUI／球）形状：触发真的没有。** `assembleRuntime` 在注入闸门那支里把 `rt.ui` 置 nil（run.go:607，
  理由见 run.go:583-590 那段），于是 T1／T2 都无处可挂（run.go:734 的 `if rt.ui != nil` 直接跳过）；
  而常驻那枚 UI 是 `ballCardUI`（resident_approval_windows.go:378，`Prompt` :428／`Update` :465），
  **结构体里根本没有 publish 字段**，两枚方法里一声 publish 都不发。
  ⇒ 常驻进程里一枚 L1 窗口正在等，快照**只在 `EvToolEnd` 那一拍**才被重建（那时窗口已关）——「在等批准」那一维在这一形状下**永远不会被读到**。
- ⛔ 结论：**「开窗口没有事件源」这句话不成立**（源就是 gate.go:294 的 `ui.Prompt`，控制台形状已经在用它）；
  成立的是「**常驻形状没有把那枚源接到泵上**」。

### 4.3 「关」这一侧：不是缺触发，是**顺序晚了一拍**

三条关法（都从 `PendingWindow` 返回，返回才触发 `defer g.closeWindow(w.corr)`，gate.go:286）：

1. 否决：`gate.go:322` 先发 `EventDismissed`（＝T2 publish），`gate.go:323` 才 return
   ⇒ **那一发快照里这枚窗口仍被列为在等**。
2. 到期：`gate.go:330` `markStarted`、`gate.go:340` 发 `EventStarted`（＝T2 publish），`gate.go:344` 才 return
   ⇒ 同上；且这一臂更糟，见 §8 F10。
3. 任务取消：`gate.go:346-347` 直接 return，**不发任何 `Event`** ⇒ 这一臂确实**既无 T2 也无 publish**。

补齐者是 T3：闸门返回后调用方一路走到 `internal/agent/loop.go:716-719` publish `EvToolEnd`
（否决／拒绝都经 `internal/tools/bridge.go:443-449` 的 `AnswerVeto`／`default` 两臂折成一次失败的 tool 结果，
不绕过那一发）⇒ 那时 `PendingWindow` 已返回、defer 已跑、窗口已从 `g.windows` 删掉（`delete` 只在 gate.go:364 一处）。
⇒ **「关缺触发」这句要改成**：关有触发、但**当场那一发读到的是旧值**，真实状态要等 `EvToolEnd` 那一拍；
只有 `ctx.Done()` 那一臂真无 T2（而它仍会被 T3 兜住）。

### 4.4 顺带量到的一条反向风险（与 §8 F10 同源）

`markStarted`（gate.go:371-385）**只写 `g.running`，不动 `g.windows`**；
`LiveL1Windows`（window_read.go:63-78）唯一的过滤是 `w == nil`（:70-72）。
⇒ 从 `markStarted` 到函数返回之间那一小段（含 `gate.go:340` 那一发 T2 publish），
一枚**已经开始执行**的调用仍会被报成「在等批准」——正是 `window_read.go:59-62` 声称「NOT listed」的那一形。
本改动面**不需要**动这段（要动就动判定分支，落进票 220 禁区），但落地腿若要当场补「关」的触发，
**必须同时**把这一形想清楚，否则会把「需要批准吃掉了还在干活」的镜像谎从窗口这一侧搬进来
（票面「外部对照」那条约束管的就是这一维）。

## §5 问题 C：装配根现状 ＋ 问题 D：最小改动面

### 5.1 三者建点与传参顺序（逐枚 file:line）

| 物件 | 控制台形状（`wisp run`） | 常驻形状（GUI／球） |
|---|---|---|
| `approval.Gate` | run.go:612 `approval.New(approval.Options{…})`（给了 `Window`/`ApprovalTimeout`/`Grants`：:615/:616/:618）；或由注入走 run.go:606 | resident_approval_windows.go:109-113 `approval.New`，**只给 `UI`/`Channels`/`Logf` 三项**；建点在 resident_windows.go:123 `newResidentApproval()` |
| 快照泵 | run.go:699（`assembleRuntime` 内，全仓唯一生产建点） | **同一枚**：常驻任务源经 resident_task_source_windows.go:265 `assembleRuntime(spec)` 拿到泵（spec 里 `gate: ra.gate, ui: ra.ui`＝:256-257） |
| `NewPanelManager` | 不经这里 | panel_host_windows.go:175（定义）← panel_resident_windows.go:204（唯一非测试调用）← `newResidentPanelManager` 定义在同文件 :183 ← resident_windows.go:142 |
| 先后 | Gate（:606／:612）**先于**泵（:699） | Gate（resident_windows.go:123）**先于**PanelManager（:142）**先于**泵（任务源启动时 :265） |
| 泵→宿主 | 无（`Out`＝台账 panel_pump.go:319） | 无：`NewPanelManager(disp, assets, dataPath)` 三参数里**没有泵、也没有 publish 回调**（`handOverPump` panel_resident_windows.go:268 那枚「pump」是 WebView2 消息泵，与快照泵同名不同物，见 :270-273 注） |

### 5.2 票 256 那枚「建得太早」会不会卡住 A 这一跳

**不会。** 两条都成立：
1. 两枚形状里 `*Gate` 都在泵建好**之前**就存在（见 5.1 的先后列），泵要读的正是那枚 `rt.gate`；
2. 更关键的是形状：`PumpSources` 每一枚位都是**方法值／闭包**（run.go:700-725 全是 `rt.xxx`），
   取值发生在 publish 那一刻而不是装配那一刻——同一条后绑定今天已经让 `Verdicts`（panel_pump.go:59 的 `rt.gate == nil` 卫哨）
   在注入形状里正常工作。所以「先有鸡还是先有蛋」在这一跳不存在。

256 的「早」咬的是另一维：常驻那枚门建时没拿 `Options.Grants`／`Window`／`ApprovalTimeout`
（resident_approval_windows.go:109-113 三项，票 256 §1 :8 逐字），
⇒ 后果是 L1 长度取默认 `3s`（queue.go:116 `DefaultL1Window`，clamp 在 gate.go:143-150）且「本会话内允许」无记账（run.go:592-599 那段把这条限制写在源里）。
它与「`L1Windows` 这一跳接不接得上」**正交**；唯一间接影响：窗口越短，「等批准」那一格的可读时间窗越短（3s 默认 vs `[risk]` 配置值）。
⛔ 本腿不动 256 的射程；并记母票 248 AC#10 已定案 ⓑ「不移动 `approval.New`」（票 256 §5 :51，账 `A534`）。

### 5.3 最小改动面（⛔ 不选形、不写码，只给面）

| # | 文件 | 面 | 会不会新开导出符号 |
|---|---|---|---|
| 1 | `cmd/wisp/panel_pump.go` | 新增一枚**非导出**方法 `(*agentRuntime).liveL1Windows() []panel.L1WindowWait`，形状照本文件 :58-76（含 `rt == nil || rt.gate == nil` 卫哨与逐条映射）。产码 8-12 行，注释按本文件惯例另计 | 否（`panel.L1WindowWait` 与 `Gate.LiveL1Windows` 都已在树上） |
| 2 | `cmd/wisp/run.go` | 字面量里加**一行** `L1Windows: rt.liveL1Windows,`（位置＝`Tasks:`（:724）之后、`Out:`（:725）之前，与「carrier 在⇒必须有人喂」那条定式同形） | 否 |
| 3 | 发布触发（按 §4 决定要不要动） | 控制台形状：**可以不动**——开＝T1 已有（§4.2），关＝T3 晚一拍兜住（§4.3）。若要「关也当场」，cmd/wisp 这侧没有位可加（`closeWindow` 的 defer 在 `internal/agent/approval` 包里），必须动 gate.go 的返回时机＝判定分支形状 ⇒ 落进票 220 禁区，需编排者另裁 | 否（但要动就动 approval 包，见 §6 尺 #12） |
| 4 | 常驻形状要出声（可选，另判） | `cmd/wisp/resident_approval_windows.go`：`ballCardUI`（:378）加一枚 publish 钩子字段、`Prompt`（:428）／`Update`（:465）末尾各发一次、装配后挂上（同 run.go:735 那一形）。三处、10-15 行 | 否（包内非导出） |

**契约面点名（编排者要求「若需要新增 `panel.*` 方法名或导出 API 必须具名点出」）**：
**这一跳不需要任何新增 `panel.*` 入站方法名，也不需要任何新增 `internal/panel` 导出 API。**
⇒ 不触 C17、不动 `internal/panel/bridge.go`、不动 `knownComposerMethod`、不动 `composer_dispatch.go`。
票 220 AC#2 甲形那枚导出（`Gate.LiveL1Windows`）已由编排者具名落账（ce693ce4 的 commit message 与 480b970d 均写 `A562`；
本腿只复核了代码事实，⛔ 未去 `docs/reports/pending-and-issues.md` 复量 A562 那一行原文 ⇒ 见 §6 量不到 #7）。

**`want-4` 那两枚锚会不会被打红**：不会。它们只读 `internal/panel/bridge.go` 一枚文件的带引号字面
（`internal/panel/git_test.go:381-389` 与 `:503-518`），而上面 1-4 号面**一枚都不碰 `bridge.go`**。
⚠ 反过来必须点名的是另一族：`cmd/wisp/subagent_blocked_197_test.go:190-192`
（「root 行没有卡命名它就不许 blocked」）——reader 一供上，`ui.Prompt` 那次取样里若同时有 root 的 L1 窗口开着，
那一格会变真、这一枚就红。本腿**不跑测试**，该格实到读数＝量不到（§6 量不到 #1）；
⛔ 落地腿不许为此放宽那枚断言，只能种样本避开 L1 路或把「此刻无窗口」量成前提。

## §6 问题 E：邻居隔离规矩（词面型／计数型负向尺逐枚）＋ 量不到的格子

### 6.1 相邻尺名册

| # | 尺 | 它数什么 | 落地腿怎么隔离才不误响 |
|---|---|---|---|
| 1 | `internal/panel/git_test.go:362`（门 3 在 `:381-389`；抽源函数 `:407`，正则 `panelMethodRe` `:394`） | `bridge.go` 里带引号的 `panel.` 字面去重后**恰等于**四枚常量 | 面 1-4 不写 `bridge.go` ⇒ 天然不响。**⛔ 别顺手在 `bridge.go` 加任何带点字面** |
| 2 | `internal/panel/git_test.go:484` 的 Door 3（`:503-518`） | 上那枚的正控：种第五枚必须数到 5、真实 `bridge.go` 必须数到 4 | 同上；这枚提醒「四」是**分母**，不是余量 |
| 3 | `internal/panel/inbound_roster_253_test.go:54-66` ＋ `:69-76` | 闭集六枚（4 枚带 `panel.` 前缀＋2 枚 `config.`）与两枚计数常量必须同动 | 只走 reader 面就不响；真要新增入站名，**名单与 `wantFullInboundRosterSize253`／`wantNonPanelPrefixedInbound253` 必须同批动**并先落台账批准 |
| 4 | 同文件 `TestFullInboundMethodRosterIsClosed`（`:419` 起）五枚分母：`:434` 声明／`:460` 运行 guard／`:469` guard case／`:493` 路由 case／`:506-517` **整个 `internal/panel` 生产件里被 guard 答到的带点字面** | 「恰四枚」那族看不见的名字这枚看得见 | 面 1-4 全在 `cmd/wisp` ⇒ 不在分母里；⛔ 别在 `internal/panel` 新增包级带点字符串常量 |
| 5 | `internal/panel/pump_test.go:123-124` ＋ `:291-292` | 快照**顶层**字节键集必须恰为 `composer,generatedAt,pending,results` | `L1Windows` 只进既有布尔那一格（`tasks.rows[].blockedOnApproval`）；**不许新开顶层键**（那是票 145 AC#3／Q-51 射程） |
| 6 | `internal/panel/approval_test.go:107`／`internal/panel/composer_test.go:50` | 双向对照：Go 侧 JSON 键集 vs 一枚仓外快照件的接口（路径字面就写在这两行里） | 这两枚**今天已红**（见 §8 F9），且它们一行 `design/` 都不读 ⇒ 落地腿**别把它们算成本跳的新红**，也⛔ 不许为了变绿去动那枚仓外件（两层禁令） |
| 7 | `cmd/wisp/panel_host_gate_test.go:30-31` | **按文件名寻址**：只 `parser.ParseFile("panel_host_windows.go")`（:33），数它的 import 面（`:37-44` 禁 `net*`） | 精确后果（本腿现读得出，⛔ 不是推断）：把宿主**改名／搬走** ⇒ `ParseFile` 失败 ⇒ `t.Fatalf`（`:35`）＝**打红**；造成**假绿**的那一动是**另开一枚宿主文件**（例如把「推快照」写进新文件），那枚文件根本不在它的分母里。⇒ 落地腿若给宿主加推快照那一臂，**必须仍写在 `panel_host_windows.go` 内**，或同批把分母改成多文件 |
| 8 | `cmd/wisp/subagent_carrier_197_test.go:430`（抽源函数 `:454` 起） | `internal/`＋`cmd/` 下非测试件里等于该枚前缀字面的包级字符串字面必须**恰 1 枚** | 新 reader 别手敲前缀，用 `panel.SubagentStreamKey`（`internal/panel/pump.go:504`） |
| 9 | `cmd/wisp/panel_pump_test.go:413`（断言在 `:426-428`／`:432-434`） | `snapshotCount()` 在 `rtHook` 那一刻必须 **0**，跑完必须 ≥1 | **装配期不许 publish**：别在 run.go:699 之后补一发「装配好先发一张」，否则这枚必红 |
| 10 | `cmd/wisp/panel_pump_test.go:189-211` | 展示的卡必须 1 枚、取样 `pending` 必须 1 条、台账里 `depth=1` 且含该 corr 的记录必须 ≥2 枚 | 面 1-2 不新增卡、不改 pending ⇒ 不响；这一枚也是 §4.2「开＝已有触发」的旁证 |
| 11 | `cmd/wisp/subagent_blocked_197_test.go:190-192`＋`:223-226` | 一行不许 section-wide 点亮；卡走后那一格必须回假 | **本跳最近的一枚负向尺**，见 §5.3 末尾 ⛔ 那条（不许放宽，只能改前提） |
| 12 | `internal/agent/approval/ticket220_l1_window_read_test.go:338-376`（反射钉）＋`:380-391` | `L1Window` 必须 0 方法／字段全 string／字段名恰 `{CorrelationID,TaskID,Tool}`，且 `*Gate` 上除 `LiveL1Windows` 外无人交出该类型 | 落地腿**不要**给 `L1Window` 加方法、也不要另写第二枚交出它的 getter（那会同时顶红 :361-373 那半） |
| 13 | `internal/agent/approval/ticket220_l1_window_read_test.go:296-337` | 「读一把不动闸门状态」八维对照（`windows` 键集合＋指针身份／`running`／`admitted`／队列 pending／grants 计数／UI `Update` 次数／连读两次同序） | 新 reader 只做**复制**（照 `panel_pump.go:58-76` 那枚 `liveVerdicts`），⛔ 不许在 reader 里顺手 `Veto`／清窗口 |
| 14 | `internal/panel/inbound_roster_253_test.go:351-365`（`productionPanelFiles`） | 该尺**跳过 `_test.go`**，所以往 `internal/panel` 加测试件不动任何分母 | 若落地腿要在 `internal/panel` 加正控用例，安全；⛔ 别加非测试件 |

### 6.2 量不到的格子（⛔ 不猜）

| # | 格子 | 为什么量不到 | 谁能量 |
|---|---|---|---|
| 1 | 接上面孔 1-2 后 `cmd/wisp` 整包是否仍绿；尤其尺 #11 会不会红 | 本腿被令⛔ 禁跑任何 Go 命令，且此刻 3 枚写腿在飞 | 落地腿自己跑 `go test ./cmd/wisp` 的终态名册 |
| 2 | 默认权限档下 `carrier197`／`pump145` 那些 mockllm 场景究竟会不会开出 L1 窗口（决定 #1 是「必红」还是「碰不到」） | 要跑，或要逐行核 D34 权威表每枚工具的等级；本腿只到「表在 `docs/PLAN.md` D34 一节（`AGENTS.md` 记 `PLAN.md:2514`）」这一层 | 落地腿 |
| 3 | 尺 #11 那枚卡被 cancel 之后到 `EvToolEnd` 之间**是否真的**发出一发快照 | 机制成立（loop.go:716-719 → run.go:1278/1297），但「该次 run 一定走到 dispatch 返回」要跑才证 ⇒ **半量不到** | 落地腿 |
| 4 | §8 F9 那两枚 `*MatchFrontend*` 尺今天红的**逐字报错行**与真实差集 | 要跑才知道；本腿只否证了「design 删除是根因」这一句 | 编排者派腿归族 |
| 5 | `internal/observe/thresholds.go` 与 SLO 带：新增一枚 reader 会不会动任何一枚资源带 | 本腿不改阈值也不跑仪器 ⇒ 未量 | 落地腿（⛔ 阈值一字节都不许动） |
| 6 | 票 256 的 `Confirming` 那一维会不会被面 4（常驻 UI 加 publish 钩子）碰到 | 那是 256 AC#0 的射程，本腿只读了票面与 §5.1 的建点 | 256 |
| 7 | `A562` 那一行台账原文（甲形的具名批准）到底写了什么边界 | 本腿只到 `.scratch` 归档与两枚 commit message 的转述层，未去 `docs/reports/pending-and-issues.md` 复量原文 | 编排者（他落的账） |
| 8 | 界面侧那一格今天能否渲染这一枚布尔 | `frontend/**` 两层禁令 ⇒ 未读未引 | 页面侧那枚 agent（票 220 §7 已划界） |

## §7 本腿闸门与自查

- ⛔ 未跑 `go test`／`go build`／`go vet`／`go list`／`go run`／`wisp`；未跑 `staticcheck`；未跑 `d22scan`。
  凡「要跑才知道」的格子一律进 §6.2，未猜任何一枚。
- ⛔ 未写、未还原、未 commit 任何别人的文件；起手脏面那枚 `cmd/wisp/config_reload.go` 本腿只**读**过一次
  （且未在本文件引用其内容或行号）。
- ⛔ 未碰票面任何一枚复选框（勾选状态一律留给编排者）；未改 `docs/specs/**`／`docs/PLAN.md`／`internal/observe/thresholds.go`／golden／`thresholds` 一字。
- 不 push；commit 带显式 pathspec（`-- .scratch/wisp/probes/220/a2/census.md`）；
  未用 `git add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；未在仓库内建 worktree。
- 搜索根一律显式（`cmd internal tools docs scripts`），⛔ 未用 `.` 作根；`frontend/**`／`design/**` 零读零引。
- 本文件零占位标记、零未交空格子；尺的字面文本未搬进被扫面（`.scratch/**` 不被任何尺扫）。

## §8 问题 F：推翻清单（逐条复量）

**F1 成立。** `ce693ce4`（`2026-10-03 10:26:19 +0800`）stat＝`gate.go +14`／`window_read.go +78`／
`ticket220_l1_window_read_test.go +400`；`git merge-base --is-ancestor ce693ce4 HEAD` 现跑＝真。
`L1Window` 确为三枚 string、零方法（`internal/agent/approval/window_read.go:43-47`），
`(*Gate).LiveL1Windows()` 在 `:63-78`，`window` 结构那枚纯数据字段 `taskID` 在 `internal/agent/approval/gate.go:106`。

**F2 成立。** `480b970d`（`10:27:02 +0800`）：`L1WindowWait`＝两枚 string、零方法（`internal/panel/pump.go:127-130`）；
`PumpSources.L1Windows` reader 位在 `:200`（规矩 :196-199「nil＝这枚宿主看不见这一维，不等于没有窗口在等」）；
消费点 `:323-327`；join 落点 `internal/panel/subagent_roster_197.go:218-221` ＋ `:236`；
导出签名 `TaskRosterSectionFrom(state, pending)` 仍在 `:187`（未改），干活的是非导出那枚 `:198`。

**F3 成立，且比那句话更重。** `cmd/` 下对 `LiveL1Windows|L1WindowWait|L1Windows` **零命中**（本腿现跑，
搜索根 `cmd internal tools docs scripts`）⇒ 不只「没人给 `L1Windows` 供数据」，
而是 `(*Gate).LiveL1Windows()` 今天**产码调用者零枚**（唯一非测试命中是 `internal/panel/pump.go:117` 的一行注释）。
这正是 A408／A483 那一族「载具在、没人喂」的第三次命中，⛔ 数调用者枚数在这里给出的结论与数「三问」一致：
缺的人＝投递者（§3.2）。

**F4 半不成立 ⇒ 具名推翻。** 编排者令里写「L1 窗口的『开／关』缺发布触发」，§4 的现读是：
- 「开」在**控制台形状**已有触发且时机正确：`gate.go:283` 的 `openWindow` 先于 `gate.go:294` 的 `ui.Prompt`，
  后者末尾 `run.go:1386-1388` 就发；旁证是测试自己的话（`cmd/wisp/panel_pump_test.go:207-208`）。
- 「开」**真的没有**的那一枚是**常驻形状**：`run.go:607` 把 `rt.ui` 置 nil ⇒ `run.go:734` 不挂 T1/T2；
  `ballCardUI`（`cmd/wisp/resident_approval_windows.go:378`，`Prompt` `:428`／`Update` `:465`）没有 publish 字段。
- 「关」不是缺触发，是**当场那一发读到旧值**：`gate.go:322`／`:340` 的 `ui.Update` 都早于 `PendingWindow` 返回，
  而 `closeWindow` 是返回时才跑的 defer（`gate.go:286` ⇒ `:362-366`，`delete` 只在 `:364`）；
  补上真实状态的是 `EvToolEnd` 那一拍（`internal/agent/loop.go:716-719` → `run.go:1278`/`:1297`）。
  唯独 `ctx.Done()` 那一臂（`gate.go:346-347`）既不发 `Event` 也不 publish。
⇒ **若照编排者那句话派腿，会多造一枚「开触发」并把真缺的常驻形状漏掉。**

**F5 成立，补一条边界。** `want-4` 那族今天两枚：`internal/panel/git_test.go:381-389`（真实比对）与
`:503-518`（种第五枚⇒5、真实⇒4 的正控）。它数的**只是 `internal/panel/bridge.go` 一枚文件**里的带引号 `panel.` 字面
（`panelMethodRe` `:394`、抽源函数 `:407`）。⚠ 所以「恰有四枚 `panel.*`」在 `bridge.go` 上成立；
写成「全仓入站方法名＝四枚」就**不成立**：闭集今天 6 枚（4 枚 `panel.` 前缀＋2 枚 `config.`，
`internal/panel/inbound_roster_253_test.go:69-76`），且 253 那枚尺存在的全部理由就是前一枚尺看不见无前缀的名字
（同文件 `:437` 逐字）。

**F6 断言成立、与本跳无因果。** 票 256（`.scratch/wisp/issues/256-resident-leg-cannot-read-those-risk-config-keys-…md` §1 :8）
的「早」指 `cmd/wisp/resident_approval_windows.go:109-113` 那枚 `approval.New` 只带 `UI`/`Channels`/`Logf`。
泵要的 `*Gate` 在两枚形状里都**先于泵**存在（§5.1 表），且 `PumpSources` 全是方法值（后绑定）⇒
**同一枚「建得太早」的顺序问题不会卡住 A 这一跳**。256 咬的是窗口长度与授权记账
（默认 `3s`＝`internal/agent/approval/queue.go:116`，clamp `gate.go:143-150`；限制自述 `run.go:592-599`），
并且母票 248 AC#10 已定案 ⓑ「不移动 `approval.New`」（票 256 §5 :51，账 `A534`）。

**F7 已被时间推翻（过期，非谎）。** 起手 `10:33:43` 时 `cmd/wisp/config_reload.go` 确为 ` M`；
`10:37:18` 由 255-r2（编排者代提）自行入库 `67ab595d`，本腿 `10:39+` 复量 `git status --porcelain -- cmd/wisp`＝空。
⇒ 编排者写那句时是对的；本腿交件时该事实已过期。本腿未写、未还原那枚文件，且未引用其内容。

**F8 票面行号漂移（票自己写了「别信这里的行号」，登记不复算为推翻）；但代码注释里的过期行号是 220-r1 落下的新瑕。**

| 票面／注释写 | 今天真身 |
|---|---|
| 票面「开／关 `gate.go:302-316`」 | `gate.go:352-366`（开 `:352-360`／关 `:362-366`；调用点 `:283`／`:286`） |
| 票面「`:377` 只有 `Veto` 按名查」 | `Veto` 在 `gate.go:421`，按名查在 `:427-428` |
| 票面「`orDefaultText` 在 `gate.go:247`／`:471`」 | 今天 `gate.go:281`／`:521` |
| **`window_read.go:38`「per orDefaultText at gate.go:275」** | 真身 `gate.go:281`（差 6 行） |
| **`internal/panel/subagent_roster_197.go:180`「gate.go:275, :471」** | 真身 `:281`／`:521` |
| 票面「`LiveApprovals` `pending_read.go:106-121`」 | 函数起 `pending_read.go:104`（差 2 行，语义同） |
| 票面「`queue.go:152-154`」 | **逐行复认成立**（`if _, clash := q.byID[corr]; clash` 在 `:152`，改写在 `:153`） |
| 票面「`subagent_roster_197.go:190-210`／`:127`」 | join 今天 `:213-221`、那一格 `:236`、字段 `:132`（480b970d 加注释后的位移） |
| `internal/panel/pump.go` 头注「泵读 `p.src.Verdicts()`」引 `:196-204`／`:263`（票面亦引） | 今天 `:239-241`／`:310-327` |
| `cmd/wisp/subagent_blocked_197_test.go:22` 引「run.go:1027-1032 说在源里」 | 真身 `run.go:1386-1388` |

⇒ 建议：`window_read.go:38` 与 `subagent_roster_197.go:180` 那两处注释行号属于**已入库的错引**，
落地腿或收口腿顺手改准即可（纯注释、不动判定），⛔ 别让它活到下一枚引用它的腿。

**F9 半成立 ⇒ 拆成两族。** 220-r1 的 commit message 与 `impl.md` §1 都写「基线那 4 枚 FAIL 根因＝同机别腿在
`design/**` 的工作树删除」。本腿现读：
- `TestC21DesignTokensFourWayAgree`：`internal/panel/tokens_fourway_test.go:50` 那枚路径常量正是 `design/assets/tokens.css`，
  而该文件今天为 ` D`（§1 的起手名册）⇒ **成立**。
- `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`：在 `internal/panel/frontend_hygiene_test.go`，
  该文件 `:217` 提到 `design/assets/tokens.css` ⇒ **大概率同因**（⛔ 本腿未逐行确认它的读取面，标 §6.2 未量尽）。
- `TestComposerContractTypesMatchFrontend`（`internal/panel/composer_test.go:50`）与
  `TestApprovalCardViewJSONKeysMatchFrontendTypes`（`internal/panel/approval_test.go:107`）：
  两枚都只打开那**一枚仓外快照件**（路径字面在上面两处）与 Go 侧，**两枚测试文件里 `design` 零命中**，
  且今天 `git status` 里 `frontend/**` **零改动** ⇒ 这两枚红**不可能**由 `design/**` 删除造成。
  另一枚根因**已经在树上写着**：`internal/panel/pump.go:31-38` 逐字说明票 248 让 Go 多出两枚 composer 段内的键、
  页面侧接口未声明、「this half cannot go self-green」（并指 `docs/evidence/s1/248-settings-write-path-r1.md`）。
⇒ 建议编排者把这 4 枚红**归两族**（令牌族＝被工作树删除挡住；契约对照族＝248 留下的自绿不能，需人拍板），
别当一枚、更别让下一枚跑 `internal/panel` 的腿拿它当「基线同一集合」再抄一遍。

**F10 新增（代码与自述不符）。** `internal/agent/approval/window_read.go:59-62` 写：
「past its deadline is NOT listed: markStarted has already moved that correlation into g.running」。
现读两条都不成立：
1. `markStarted`（`gate.go:371-385`）只写 `g.running`／`g.order`，**不删 `g.windows`**；全仓唯一删那枚 map 的地方是
   `closeWindow`（`gate.go:364`），而它只在 `PendingWindow` 返回时经 defer 跑（`gate.go:286`）。
2. `LiveL1Windows`（`window_read.go:69-74`）唯一的过滤是 `w == nil`。
⇒ 从 `gate.go:330` 到 `gate.go:344` 之间（含 `gate.go:340` 那一发 `EventStarted` publish），
一枚**已在执行**的调用会被枚举成「在等」。配套的钉也拦不住：
`ticket220_l1_window_read_test.go:380-391` 那枚「started 不报」是先 `markStarted`、**从未开过窗口**
（`g.windows` 本来就是空）⇒ 它对真实时序没有鉴别力（正控缺位，属票 161／171 那一族）。
⛔ 本腿不修、不派形；但**落地腿若照 §4.3 去补「关」的当场触发，会第一次真读到这条不符**，故在此点名。
