
## 5. 尺读数与未跑清单

> 每把尺给四件：**正则／抽取规则＋根＋读数＋时刻**。全部 `grep`/`sed` 现量，⛔ 零 `go` 命令。
> 起手段的 `git` 读数（`c6cf66e6`／`dev`／五根干净）在 §0，不在此重复。

| 尺 | 规则逐字 | 根 | 读数 | 时刻 |
|---|---|---|---|---|
| R-① | `grep -rniE "occupanc\|contextused\|usedtokens\|windowused\|prompttokens" --include=*.go cmd/wisp internal/panel internal/agent \| grep -v "_test.go"` | 快照两包＋环路包 | **0 行**（占用命名零枚） | 11:25:29 |
| R-①b | 同正则，根放宽到 `cmd internal tools` | 三根 | **19 行**，逐枚过目后全在别的地界：`internal/llm/openaichat/wire.go`（provider 报文 `prompt_tokens`）、`internal/projctx/projctx.go:121/:141/:431`、`internal/tools/subagent_197.go:362` 与 `internal/winsec/resolve.go:108`（两枚注释里的英文 occupancy）、`tools/mockllm/*` | 11:16x |
| R-② | `grep -rn "^func (l \*Loop) [A-Z]" internal/agent/*.go` | `internal/agent` | **8 枚**：`loop.go:260 Budgets`／`:263 History`／`:272 Reset`／`:282 Steer`／`:328 RunAsync`／`:339 Run`＋`prompt.go:280 AttachProjectInstructions`／`:288 ProjectInstructionManifest`。**无占用形方法** | 11:10:59 |
| R-③ | `grep -rn '^func (l \*Loop) [A-Z][A-Za-z0-9]*(.*taskID' internal/agent`（＝`.scratch/wisp/probes/154/gate-clauses.sh:446-447` 那枚 G3 腿的正则） | `internal/agent` | **0 枚**（G3 今日安静，与尺自设的 `want_n 0` 同形） | 11:25:29 |
| R-④ | `grep -rn "loop.Budgets()" --include=*.go cmd \| grep -v _test` | `cmd` | **3 枚**：`cmd/wisp/run.go:1056`／`:1062`／`:1109`；**取 `.ContextWindow` 的＝0 枚**（三枚分别只取 `.PromptTotal`／`.PromptTotal`／整包递 `agent.NewSpiller`） | 11:25:29 |
| R-⑤ | `grep -rn "\.History()" --include=*.go cmd internal tools \| grep -v _test \| grep -v "internal/agent/"` | 三根、排包内 | **0 枚**（`Loop.History()` 生产零包外调用者） | 11:25:29 |
| R-⑥ | `grep -rc "json:\"" --include=*.go internal/agent/approval` 求和；同尺加 `\| grep -v "_test.go"` | `internal/agent/approval` | **两口径同为 0**：含 17 枚测试件＝0，排测试件（9 枚产码件）＝0 | 11:25:29 |
| R-⑦ | 读两把尺的型别名册：`sed -n '55,70p' internal/panel/composer_test.go`／`sed -n '113,122p' internal/panel/approval_test.go` | `internal/panel` 测试件 | A＝**6 对**（`Snapshot`/`ComposerState`/`ModeView`/`WorkspaceView`/`AttachmentRef`/`ResultChunk`）；B＝**3 对**（`ApprovalCardView`/`ResultChunk`/`Snapshot`）；交集＝`Snapshot`＋`ResultChunk`；**`TaskRosterSection`／`TaskRowView`／`InstructionsSection`／`GitView` 两把尺都不在册** | 起手逐行过目时（约 11:0x） |
| R-⑧ | `sed -n '235,270p' internal/panel/composer.go \| grep -cE "^\s+[A-Z][A-Za-z0-9]+ +"` 对 `… \| grep -c "json:\""`（同法跑 `Snapshot` 57-92、`ApprovalCardView` 39-60） | `internal/panel` | `ComposerState` **11 字段／11 tag**；`Snapshot` **6／6**；`ApprovalCardView` **10／10** ⇒ 快照三型**零枚**无 tag 导出字段 | 11:27:23 |
| R-⑨ | `grep -rn "composer,generatedAt,pending,results" --include=*.go . \| grep -v "\.scratch"` | 全仓（排探针） | **3 枚**：`internal/panel/pump_test.go:123`、`internal/panel/pump_test.go:291`、`internal/panel/subagent_roster_197_test.go:214` | 11:25:29 |
| R-⑩ | `grep -rn "NewSnapshotPump" --include=*.go cmd internal tools \| grep -v _test \| grep -v "func NewSnapshotPump"` | 三根 | **1 枚**＝`cmd/wisp/run.go:699`（产码装配点全仓仅此一处） | 约 11:0x |
| R-⑪ | `grep -rn "L1Windows:" --include=*.go cmd internal \| grep -v _test` | 两根 | **0 枚**（唯一命中在测试件 `internal/panel/subagent_blocked_220_test.go:81`） | 11:25:29 |
| R-⑫ | `grep -rn "agent.Loop" --include=*.go cmd \| grep -v _test` | `cmd` | **1 行，且是注释**（`cmd/wisp/run.go:21`）⇒ `agentRuntime` 里 `*agent.Loop` 类型字段 **0 枚** | 11:25:29 |
| R-⑬ | `grep -rn "rt.execute(\|run.execute(" --include=*.go cmd \| grep -v _test` | `cmd` | **2 枚**：`cmd/wisp/run.go:262`、`cmd/wisp/resident_task_source_windows.go:462` | 11:27:23 |
| R-⑭ | `grep -rn "\.TokensIn" --include=*.go cmd internal tools \| grep -v _test` | 三根 | **2 枚**，都在 `internal/memory/dao_misc.go:158/:179`（SQLite 日成本行，与快照无关）⇒ **`cmd/wisp` 对 `Event.TokensIn` 的产码消费＝0 枚** | 11:25:29 |
| R-⑮ | `grep -rn "MaxLoggedString = " internal/observe/redact.go` | `internal/observe` | `:37` 逐字 `	MaxLoggedString = 512`（出口 `:141` 逐字 `	return truncate(s, MaxLoggedString)`）⇒ 占用若想在账本里带数字，仍受这枚 512 约束（本腿未动它一字） | 11:25:29 |

### 5.1 未跑清单（⛔ 硬闸禁跑 `go`，全部标〔预测〕并具名归口给编排者）

| 格 | 只有跑包才能答的问题 | 本腿给的〔预测〕＋理由 | 归口 |
|---|---|---|---|
| U-1 | `go test ./internal/panel/` 今天几枚红、名册逐名是什么 | 〔预测〕**≥4 枚**：尺 A＋尺 B（`docs/reports/HANDOVER.md:485-486` 09-28 名册在册）＋`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`＋`TestC21DesignTokensFourWayAgree`（`internal/panel/tokens_fourway_test.go:439`，红因在 `design/**` 工作树，不属本编队）。⛔ 这不是读数 | 编排者跑 |
| U-2 | 新增一枚带 tag 的 `ComposerState` 字段后，尺 A 的红句多列几枚键 | 〔预测〕多 **1 枚**（该字段名），且**只有尺 A 红、尺 B 不红**（R-⑦ 的交集里没有 `ComposerState`） | 编排者跑 |
| U-3 | 新增一枚 `Snapshot` 顶层字段（值类型、无 omitempty）后三枚四键钉红几枚 | 〔预测〕**3 枚全红**（R-⑨ 三枚 fixture 都不设该 reader，但值类型节无法缺席）；若做 pointer＋omitempty 且 reader 缺席 ⇒ 〔预测〕**0 枚红** | 编排者跑 |
| U-4 | `go test ./cmd/wisp/` 是否整包无实测读数 | 不可判——票面 §9 第 5 条逐字记「**`cmd/wisp` 整包至今没有实测读数**」，本腿按硬闸未跑。**这一格直接决定 167-r2 现在能不能起飞**（票面把它排在 `267-r1` 之后） | 编排者 |
| U-5 | `gofumpt -l` 与 `scripts/d22scan.sh` 的基线数 | 不可判（禁跑）。⚠ 具名提醒落地腿：`AGENTS.md §1.2` 记录的 `emojiRe` 射程含 `U+2200–U+22FF` 与 `U+2600–U+27BF`、**注释豁免、字符串不豁免** ⇒ 占用那一节的中文 reason 里写 `≤`(U+2264) 或 `✓`(U+2713) 会被仪器抓，写 `→` 不会（箭头带不在射程内） | 编排者／落地腿 |
| U-6 | 两把尺的红腿名册与本腿交付的差集 | 不可判（本腿零改动，差集**应**为空集，但**没跑就没有名册**） | 编排者 |

## 6. 判不动／量不到（具名，⛔ 不写"应该没问题"）

1. **尺 A／尺 B 今天的实际颜色与红句里已列的键名册**——见 §5.1 U-1。本件所有"会红一发"的判语都是**射程**判语，不是颜色判语。
2. **`frontend/src/lib/panel.ts` 到底声明了哪些键**——⛔ 硬闸 2 禁读 `frontend/**`。⇒ 直接后果：§4.1 末段那批"新键是否已被页面声明"我一律判不了，只能引 `internal/panel/composer.go:87-90` 与 `internal/panel/pump.go:31-38` 的**代码注释自述**（注释不属权威来源）。这一格是**禁令造成的量不到**，不是能力问题。
3. **票 145 那一寸能不能盖住"新增一枚带 JSON 键的字段"**——派单要本件"逐条判并具名写越界与否"，本件按字面判了（§3.2），但**权威边界文本**在 `docs/reports/pending-and-issues.md` 的 `A273`／`A389`／`A560`／`A590`／`A610` 诸条里，⛔ 硬闸 3 不许本腿读那枚台账 ⇒ 本件的"寸内"判语**降级为按票面 AC#7 与 `.scratch/wisp/dispatches/2026-09-26-093x-thaw-panel-for-145.md` 字面读**的结果，编排者须用台账原文复核。
4. **子代理那一发的上下文用量到底能不能读**——量到的是"读不到"（`internal/tools/subagent_197.go:319`／`:344` 两枚局部变量，`internal/tools` 无字段握它），但**"要不要为一发子代理单独报占用"是产品问题**，本腿不判。
5. **占用该报"当前历史"还是"跨轮累加"**——票面 §9 第 1 条与 `167-c2` §4② 都把这格列为"不答则整枚不做"的前置裁夺；本件只把两轴的单位差异钉死（§2.1／§2.3），⛔ 不选形。归编排者或机主。
6. **resident 腿的快照泵到底存不存在**——本腿现量产码装配点只有 `cmd/wisp/run.go:699` 一枚（R-⑩），与 `167-c2` §5 第 3 条"常驻面板今天既不发快照也不显示快照内容"同形；但常驻腿会不会**复用**同一枚 `rt.pump`（同一 `agentRuntime` 实例）本腿未追到底 ⇒ 若落地腿要靠 `execute()` 往 runtime 写占用，**常驻腿两发任务之间的残留读数**是量不到的一格（`cmd/wisp/resident_task_source_windows.go:440` 的单槽规则与 `:474` 的 `src.running = false` 在册，占用读数的清理点未追）。
7. **`rt.loopOpt` 的并发安全**——`cmd/wisp/run.go:1022` 那行写入**没有任何互斥**（对比 `rt.instrLoader` 有 `instrMu`、`rt.seenTasks` 有 `taskMu`、`rt.lastSnap` 有 `snapMu`）。本腿只量到"无锁"这一形，**没跑 `-race`** ⇒ 它今天是否已是活体竞争，属量不到（若落地腿加同型字段，请带锁，别复制这一形）。
8. **`gate-clauses.sh` 的聚合退码与红腿名册**——禁跑（那枚脚本会跑 `go`）。⛔ 本件只复量了其中 G3 那一支的**正则命中数**（R-③），没有复量它今天的 `腿数`／`BAD 名册`；派单提到的"在册只 `G6neg`"是过期读数，本件不复用。

## 7. 我写错的读数（自我对抗，真改）

1. **把 `167-c2` 的"14 枚 JSON 字段"当既有事实抄了下来。** 首版 §1.1 跳⑨ 写 `ComposerState` 有 14 枚 JSON 键；现量（R-⑧，11:27:23）＝**11 枚字段／11 枚 tag**。已就地改为 11 并标注推翻。**这是我抄前趟计数未复量的直接后果**——同一个错 `167-c2` §6.1 刚为 `167-a1` 记过一笔。
2. **`internal/agent/compress.go` 的累加行写错。** 首版引 `compress.go:112-115 累加 ApproxTokensOf(m.Content)`；现量累加那句在 `:113` 逐字 `		n += ApproxTokensOf(m.Content)`，`:112` 是 `	for _, m := range hist {`。已改。
3. **`internal/agent/budgets.go` 的 `ApproxTokens` 注释行号写错两格并漏了字。** 首版在 §2.1 N8 引 `:126-128`、在 §2.3 引 `:125-128` 却把第三行行首的 `4)` 归并到上一行。现量注释起 `:125`、止 `:128`，第三行逐字是 `// 4) so the loop's budget math and the seam's pre-flight estimate cannot`。两处已统一到 `:125-128` 并补回漏字。
4. **`internal/panel/composer.go` 的 `ModeView.Current` 行号差一行。** 首版引 `:143`（那是注释行），现量逐字 `	Current string \`json:"current"\`` 在 `:144`。已改。
5. **`internal/panel/git.go` 的 `GitSwitchBlockedReason` 引成了注释段。** 首版写 `git.go:72-80 const GitSwitchBlockedReason = ...`；现量常量在 `git.go:82`（`:72` 起是注释）。已改。
6. **`internal/panel/l2_grant_boundary_test.go` 的 `reflectionTopKeys` 行号。** 首版引 `:1755`（那是它体内的 `NumField` 循环），现量函数起 `:1753`。已改。
7. **`internal/tools/task.go` 的 fail-closed 句行号。** 首版引 `:159`，现量该逐字句在 `:160`（`:158` 是函数签名、`:159` 是 `if o.State == ""`）。已改。
8. **`cmd/wisp/config_receipt_255_test.go 一族`。** 首版用"一族"这种没有单一 `file:line` 的写法交付 B5 的读者，属**不合格引用**；现量唯一读者是 `internal/panel/config_route_248_test.go:323`，已换成逐字行。
9. **对 `167-c2` 的行号过度乐观。** 起手打算沿用 c2 的 `loop.go:253/256`；本件起手即现量到它们已移到 `:260/:263`（差 7 行），并在 §1.0 显式登记。**若当时抄了，§1 的名册里会有两枚指向错字符。**
10. **一处判语在写完才想到反例**：§3.2 表庚原写作"零新键的复用"——严格讲它**不是复用，是把别人的维改了**（`StatusKnown` 的语义是 D43 状态名，不是忙闲）。表庚的判语已收紧为"唯一零新键的形状，而它是被明令禁手的那一枚"，§3.2 末句同步改过。**保留这一格而不是删掉，是因为它记录了"便宜路看起来存在"的错觉从哪儿来。**
11. **`internal/panel/pump.go:45-47` 的自述被我先当成读数、后才发现它过期。** 该注释逐字说四枚钉里"the two byte-level key-set nails in pump_test.go (:111-124, :270-276)"；现量第三枚同名钉在 `internal/panel/subagent_roster_197_test.go:214`（R-⑨＝3 枚），且注释点名的两个行段今天实际是 `:111-125`／`:283-293`。⇒ 我在 §4.2 已把它当"注释自述、非读数"引用并具名标差；**记在这里是因为我第一版差点直接引它当名册**。

## 8. 交件判语三行

1. **落地腿最少要动哪几处，占用条才能有别于"永远未知"（四处，缺一即恒空；⛔ 一处都不必新增 `Loop` 方法，G3 射程已复量 0 命中 R-③）**：
   ① `cmd/wisp/run.go:266-377` 的 `agentRuntime` **加一枚握读数的字段＋一把锁**（照 `:361-362 instrMu`／`instrLoader` 那对形，⛔ 别照 `:292-293 loopOpt` 那对无锁形）；
   ② `cmd/wisp/run.go:1023` 造出 `loop` 之后、`:1099` 起跑之前，在 `execute()` 里**交一次手**（照 `:1070` 那枚 `rt.setInstructionLoader(instrLoader)` 的时序——那里已有"泵先装配、读口后落地"的既有先例，见 `cmd/wisp/panel_pump.go:90-92` 自述 `// execute(), after the loop exists and before its first turn.`）；
   ③ `cmd/wisp/panel_pump.go`（145 寸内文件）**新加一枚 reader 方法**：分子由 `internal/agent/loop.go:263 History()` ＋ `internal/agent/budgets.go:150 ApproxTokensOfMessage` 在装配根复算，分母取 `cmd/wisp/run.go:302 rt.endpoint.ContextWindow`（它的 `0` 有 `internal/config/schema.go:353` 逐字具名"unknown"）；⛔ **不许取 `Budgets().ContextWindow` 当分母**——`internal/agent/budgets.go:85-88` 把它折成恒非零，"未知"那一维在那一步已经死了；并在 `cmd/wisp/run.go:699` 那**全仓唯一一枚**产码装配点连线（R-⑩）；
   ④ `internal/panel/pump.go:138-215` 的 `PumpSources` **加一枚 reader 槽** ＋ `internal/panel/composer.go:235-270` 的 `ComposerState` **加带 json tag 的标量字段**，"未知"按 §3.1 的 A／B 两族任一模子长。
   ⇒ 这四处里**第 0 跳是①＋②**；只做③④＝读到恒空，正是 `167-c2` §6.1 补的那一步，本件把它量成了具体两行。⚠ 另有一处**不必动但必须写进判据**的坑：§4.3 的 W④——占用字段落进 `composer.go` 时，注释里**不许出现** `SetMode`／`PermissionMode =`／`perm.Store.Set(` 三串字面（`internal/panel/composer_test.go:171` 那发逐行扫全文，注释也扫）。
2. **哪些一寸属越界、要停手上报（五条，逐条具名）**：
   (a) **新建 `internal/panel/occupancy.go` 之类的新文件／新导出类型**——票 145 那一寸的字面名册只有 `composer.go`／`pump.go`／`panel_pump.go` 三枚**已存在**文件，批准文字逐字是"只到新增字段为止"；先例说明新文件要靠**自己那份派单的具名授权行**（`.scratch/wisp/dispatches/2026-09-28-150x-impl-181-r1-go-side-git-read-surface.md` §1 逐字 `✅ 新建：internal/panel/git.go`（或同层你判更合适的名字）），不是靠 145。⇒ **要编排者在派单里补一行；不补就停手。**
   (b) **改 `internal/panel/instructions_200.go:158-194` 借 `projctx.Bundle` 那两个现成数**——文件不在名册里＝越界，且语义错（§2.1 N7：说明子预算 ≠ 对话窗口）。⇒ 即使后来裁"再开一枚项目说明预算小格"，也要新授权，且必须与占用**分键**。
   (c) **动 `internal/panel/subagent_roster_197.go` 或 `internal/tools/task.go` 的 `MarkRoot`／`StateAnswer` 那一维**——票面 §9 第 1 条末行逐字「⛔ 不许顺手把那一维改了」（票 196／`A394`）。⇒ 双越界（文件＋契约维），**必停**。
   (d) **给 `internal/agent.Loop` 新加任何导出方法**——`docs/reports/HANDOVER.md` 18:5x 那条在册裁定逐字是"裁定＝入口放 tools 侧，非放 Loop 不可就停手上报"；⚠ 本件同时量到 G3 只按字面 `taskID` 抓（R-③＝0 枚）、且 HANDOVER 自己写过"G3 安静≠合规（改名即失灵）" ⇒ **"G3 没响"绝不可以当"合规"的证据**。
   (e)（附带，本件的边界自证）**序号／停止／草稿／崩溃自救四枚一律未动笔**：`internal/panel/approval.go:39-59` 的十枚卡面键、`internal/panel/bridge.go:42-45` 的四枚入向名册、`cmd/wisp/panel_inbound.go`、`internal/observe/**`、`cmd/wisp/doctor.go` 本件**只读不判**，即便 §3／§4 顺手量到了它们的名字。
3. **哪些归机主裁（三枚，⛔ 编排者不代裁）**：
   (a) **分子语义**＝"当前这轮历史吃多少"（估算 token）还是"这次任务一共花了多少"（计费 token、跨轮累加）——两轴单位不同、不能相除（§2.3）；票面 §9 与 `167-c2` §4② 都把这格列为"不答则整枚不做"。
   (b) **覆盖范围**＝占用条只报根任务，还是也要报子代理那一发；后者今天**结构上读不到**（§1.0 末段），要做就是新增一跳交接口＝扩大写面，属人工批准面。
   (c) **"页面侧同批声明"那一跳谁写**——占用一旦新增 JSON 键就让 §4.1 那把尺点名，转绿在 `frontend/**`，而票面 Status 逐字 `skipped=frontend(owner-delegated)`、硬约束 3 禁本编队写 `frontend/**` ⇒ 与 `Q-51`／`Q-76` 同族摆给机主。⚠ 本件**不主张**"为了不发红就把字段做成无 tag"——§4.4 已量明那一形会让两把尺同时哑火而线上多一枚幽灵键，那是把仪器骗过去、不是把契约对上。
