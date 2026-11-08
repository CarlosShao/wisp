# 145-snapshot-fields — wave 写腿：快照新增顶层 `run` 工作态段（reader-gated，213/187/65 同范式）

- 日期：2026-11-07 20:3x–21:0x（ Asia/Shanghai）
- 腿：票 145 本波写腿；写面 = `internal/panel/**` ＋新测试件＋本证据件＋票面 Progress log 一行
- 起手锚：`7202103f3`（HEAD，wave74 已收）；起手 `git status --porcelain -- internal/panel/` ＝**空**（无人在飞）
- 收件形式：只 commit、不 push；每次 git add/commit 均带显式 pathspec

## 0. 起手三核读数（票面锚过期，现量为准）

1. **票面标题「恰四字段」「十一个无输入」两枚锚均已过期**（与票面自己 10-03/10-07 的两节更正一致，本程逐枚复跑）：
   现量 `Snapshot`（`internal/panel/composer.go:58` 起）＝**4 基线字段＋8 枚 reader-gated 段/键**
   （`instructions`=200、`tasks`=197、`commands`=213、`modelSwitch`=187、`ballStyle`=65、
   `git`/`files`/`currentModel`+`credentialState` 等=181/190/145-r3/248）；页面侧 `UnfedScreen` 全树 0 命中、
   屏名册 9 枚（`panel-views.ts:26-35`）——与 `145-f1`（`A657`）、`145-a2` 的读数一致。
2. **十四行表真身**＝`docs/PLAN.md:3479-3494`（票面引 `:3473-3488` 已漂 6 行）。
3. **`git status --porcelain`**：`internal/panel/` 空；`cmd/wisp/` 只有他人两枚 323 未跟踪测试件（本程零触碰）；
   `internal/ball/`、`internal/config/` 空。无红/绿来自在飞件。

## 1. 本程的关键判断（铁律 ⓪⓪⓪，逐字段「有真源/无真源」裁定，每条可复跑）

判据＝AC#1 ⓑ「一个字段若只能靠填常量才通，不许进落地集」＋AC#6「哪枚用例断言值来自真源」。
本程把「有真源」收紧为：**值已由生产码写在某个今天存在的载体上**（事件值／既有导出读法），
holder 只是搬运。逐字段裁定表：

### 1a. 落地（有事件真源，三枚新牙全绿）

| 字段 | 真源（file:line，2026-11-07 现量） | 判据 |
|---|---|---|
| `run.tools[]`＝`{callId,name,outcome,resultText}` | `EvToolEnd` 三枚发布点：`internal/agent/loop.go:845`（正常结束）、`:988`（截断）、`:1023`（failOpenCalls 补判）；四值词表＝`internal/agent/journal.go:40-43`（success/error/cancelled/truncated），逐字携带不自造 | 行 4/5 的「条目＋终态」半格 |
| `run.reasoning` | `EvReasoningDelta`：`internal/agent/loop.go:645`（`l.publish(Event{Kind: EvReasoningDelta, ..., Text: ev.Text})`）；`cmd/wisp/run.go:2268-2273` 的 console sink 注释逐字点名「the reasoning field is ticket 145 row 3, which is a key the snapshot does not have」 | 行 3 的文本半格 |
| `run.error`＝`{class,detail,text}` | `EvError` 携带分类错误的发布点：`internal/agent/loop.go:1047`（failWith，`Err: err`）；`observe.Error` 字段＝`internal/observe/errors.go:174-180`（`Class`/`Detail` log-safe） | 行 10 的 `error_class`＋技术细节半格 |
| `run.stuck`＝`{text,repeatLevel}` | `EvStuck`：`internal/agent/loop.go:576`（C22 阶梯，`Text: rem.Text, RepeatLevel: rem.Level`）＋`:1038`（brakeStuck） | 行 11 的「重复了什么」半格 |

holder＝`internal/panel/run_feed_145.go` 的 `RunFeed`（`sync.Mutex`，StreamLog 同款锁纪律；
方法收**纯值**不收 `agent.Event`——`internal/panel` 不 import `internal/agent`，无新包边）。
两枚界：`runToolRowsKept=64`（逐出计数进 `toolsDropped`，不静默）、`runReasoningKeepRunes=4096`
（尾窗＋省略计数在 marker 里）。行 4/5 的**运行中**态与**耗时**列、行 3 的耗时列见 1b。

### 1b. ⛔ 登记（无真源／跨 seam／归别票，本程零字落地）

| 维度（十四行表行号） | 裁定 | 依据（现量） |
|---|---|---|
| 工具行「运行中」态 | 无写者，不落 | `EvToolStart`（`internal/agent/sink.go:25`）生产码**零发布点**：`grep -rn EvToolStart --include=*.go internal cmd` 剔 test 仅命中 kind 声明＋run.go:2275 的 handler |
| 工具行参数摘要／风险徽标／耗时 | 无事件源，不落 | `agent.Event`（`internal/agent/sink.go:51-66`）不携带 args/risk/时长；`panel_pump.go:256 panelArgs` 的摘要在 cmd/wisp 事件之外，泵不可达 |
| `usage`／`cost`（行 2 token 计数、行 12 成本） | 跨 seam，不落 | 累加器 `cost := Cost{...}` 是 `internal/agent/loop.go:416` 的 Execute 局部，`:1061-1063` 随 `res` 出口；进程内无活 holder、无 per-turn usage 事件；`EvDone` 的 `TokensIn/TokensOut`（`loop.go:1095`）只是无 cached/金额/状态词的终末回声，单带它=半回声假行 |
| 终结态词＋取消原因（行 13） | 无事件源，不落 | `EvDone.Text = summary = res.Text || msg`（`loop.go:1071-1074`）：取消时带的是**部分回复文本**不是取消原因，D43 终结态词只进 `res.Status` 不上事件；按它渲染「已取消：…」＝造假 |
| `phase`/elapsed（行 1）、`humanText`（行 10）、`remainingMs`（行 8）、`reasoningMs` | 纯无源，不落 | 复认 `145-c2`/`145-a2` 两程同判：映射表不存在、clock.After 不给余量 |
| `pending[].position`、`context.occupancy`（行 6/12 邻域） | 归票 167-r2，不落 | 票面 09-28 明令「两条腿不许各加各的 pending[] 键」；`A699` |
| `composer.git` 分支/工作树清单 | 归票 186，不落 | `145-a2` C2 切面 |
| `view`「当前屏」（行 7 邻域） | AC#3 裁决维持，不落 | Go 侧不知九枚屏名；ⓑ 重开 C17 白名单＝未定义即停 |

## 2. 落的形状（reader-gated，⛔ 第四范式零新增）

- `internal/panel/composer.go:160`：`Run *RunSection json:"run,omitempty"`（顶层指针段＋omitempty；
  注释写明 Q-51 未声明红＝213/187/65 同一既定程序）。
- `internal/panel/pump.go:250`：`PumpSources.Run func() *RunSection`；`:458-464` 填段（nil reader ⇒ 键缺席，
  四键字节钉 `pump_test.go:123`／`:291`／`subagent_roster_197_test.go:214` 不触发）。
- `internal/panel/run_feed_145.go`：`RunSection`/`ToolRunRow`/`RunErrorRow`/`RunStuckRow` 载体＋`RunFeed` holder；
  空 feed 渲染 `"tools":[]`（commands/modelSwitch 同款，永不 null）。
- **留给 cmd/wisp 面的两跳（本程零触碰，finder 不用行号）**：①喂 holder——console sink 的事件 switch
  （`cmd/wisp/run.go` 的 `func (c consoleSink) Publish`，EvToolEnd/EvReasoningDelta/EvError/EvStuck 四个 case
  全已存在，各加一行 `c.runFeed.Xxx(...)`）；②接 reader——装配根建 `PumpSources` 处
  （`grep -n "PumpSources{" cmd/wisp/*.go`）加 `Run: runFeed.Section`。
  在这两跳落地前，生产包不带 `run` 键（钉保持绿的形状），测试经 `PumpSources.Run` 直接验线。

## 3. 用例（三枚，`internal/panel/run_feed_145_test.go`）

1. `TestRunSectionTravelsOnlyFromItsReader`（`:27`）——reader-gated 断言＋**突变牙**：摘掉 reader，
   `"run"` 键必须缺席（字节级）；空 feed 渲染 `"tools":[]`。
2. `TestRunFeedCarriesOnlyWhatEventsWrote`（`:55`）——**AC#6 答句**：按 loop.go 发布点的原值喂脚本
   （两枚工具行＋分类错误＋stuck＋三枚 reasoning delta），逐字段断言与脚本逐字一致、顺序保持、
   字节含 `"outcome":"success"` 等；notice 形错误（无 class/detail/text）**不写行**（不造假分类）。
3. `TestRunFeedBoundsAreHonest`（`:122`）——两枚界的诚实规则：逐出计数、最老先逐、
   尾窗保留 tail＋省略计数在 marker。

## 4. 门禁读数（环境：unset GOTMPDIR TMPDIR TEMP TMP；PATH+=build、third_party/sherpa-onnx；⛔ 未设 CGO_ENABLED=0）

| 门 | 读数 |
|---|---|
| `gofumpt -l`（四枚点名件＋全 `internal/panel/`） | **空** |
| `go vet ./internal/panel` | **rc=0** |
| `go build ./internal/panel` | **rc=0** |
| 新三牙单跑 | 3/3 PASS 0.067s |
| `go test ./internal/panel -count=1`（改后） | rc=1，RUN 310 / PASS 204 / **FAIL 3** |
| **改前基线**（`-overlay` 指 HEAD 版 composer.go/pump.go，仓外副本 `.scratch/wisp/probes/145/wave-leg/`） | rc=1，红名册**逐名同册 3 枚**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`、`TestComposerContractTypesMatchFrontend`、`TestC21DesignTokensFourWayAgree` |
| 红差量 | **零枚新红**；两枚前端声明尺本就红（commands/modelSwitch/ballStyle 未声明，213/187/65 既定程序），本程只让点名列表多出 `run` 一词；C21 尺红在 ball tokens，非本面 |
| `sh scripts/d22scan.sh` | 建证据件前 1 finding＝ban #9 phantom-citation（run_feed_145.go:25 引本件、件未建）；**本件落盘后复跑 rc=0 clean**（最终读数见 Progress log） |
| LF | `file` 四枚件无 CRLF |

## 5. TS 侧「应当长这样」逐键清单（按 `A655` 口径，由 owner/前端自己落；本程 `frontend/**` 零写面）

```ts
// frontend/src/lib/panel.ts 的 PanelSnapshot 增补（与 Go 侧逐键对齐）：
run?: {
  tools: { callId: string; name: string; outcome: string; resultText?: string }[];
  toolsDropped?: number;
  reasoning?: string;
  error?: { class: string; detail?: string; text?: string };
  stuck?: { text: string; repeatLevel: number };
} | undefined;
```

## 6. 离「十四态能真画」还差什么（本程现量口径）

1. **喂线两跳**（§2 末）落在 cmd/wisp——本波 wave75 的面。
2. **出向最后一跳**（Go→页推＋页面接收器）＝票 35，`A657` 判语维持。
3. **页面声明**＝Q-51（两枚双向尺的解冻钥匙）。
4. usage/cost/终结态/计时类见 §1b——缺的是**源**（loop 局部 holder／事件字段），不是字段；先解源再谈。

## 7. 11-08 19:1x　只读判定程 wave3（零产码；起手三核＋「还剩什么」终局判定；锚 `38628bee2`）

**起手三核**：①派单前提「Snapshot 恰四字段」**过期**——现量顶层 **11 枚直接字段**（`internal/panel/composer.go:57-174`：pending/results/composer/generatedAt/instructions/tasks/commands/modelSwitch/ballStyle/run/plugins），`PumpSources` **21 面**（`pump.go:146-299`）且**全部在装配根接线**（`cmd/wisp/run.go:998-1127`，含 `Run: rt.runFeed.Section` `:1059`）；②引文＝census r1 §4.2 甲组 12 行现读真字段；③门禁在 `git clone -c core.longpaths=true --shared` 导出树跑（主树 `frontend/**` 三枚 M＝前端腿在飞）。克隆一次撞 Windows MAX_PATH（一枚长文件名 tracked 件 checkout 失败），失败件留 `.scratch/wisp/probes/145/wave3-clone/` 不删，重克隆 `wave3-clone-lp` 成功（HEAD 同锚、树净）。

**AC#1 甲组 12 行逐行现量判定**（权威名册＝`145-snapshot-field-census-r1.md` §4.2）：

| 行 | 字段 | 现量判定 |
|---|---|---|
| 3 | reasoning 文本 | **已落**＝`run` 段（`run_feed_145.go`），装配根 `:1059`＋sink `:2366` 已接（wave-leg2 闭账） |
| 11 | stuck | **已落**（`stuck` 键，sink `:2380`） |
| 4 | 工具调用 | **部分落**（callId/name/outcome，sink `:2374`）；「运行中」态未落（见下 5 族之 3） |
| 5 | 展开半 | **部分落**（callId/resultText）；ArgsJSON 未落（5 族之 4） |
| 10 | 失败 | **部分落**（class/detail/text，sink `:2390`）；Retryable 未落（5 族之 5） |
| 7 | pending 十字段 | **早已通**（零 Go 改动行） |
| 14 | 注入来源 | **早已通**（reason 链） |
| 8 | L1 窗口 | **按裁定形状落**＝220 的 `L1Windows` id-join（`run.go:1036`）；remainingMs 在 `run_feed_145.go` 头注释具名不落 |
| 6 | Depth | **裁走**＝167 射程（position/occupancy）＋前端可算（r3 拒收维持） |
| 1/9 | 无源／原生侧 | 不在本组（census 原判维持） |

**AC#2b 两维＝票 187 已落，⛔ 不重复做**：①模型清单＝`ModelSwitch.Models`（**enabled-only**：`internal/config/settings_187.go:202-205` ∩ 装配根 `run.go:1071-1084`）；②档位集＝`ThinkingLevels`＝当前模型自己声明的表（`settings_187.go:200`）。**两枚半枚未落且都卡禁面**：`Display`（源在 `schema.go:372`，`ChatSwitchOptions` 不导出，加＝动 internal/config）；「各家适配器实际接受的词表」（三张包内私有 `var` 无导出面，r2 R-3② 维持）。派单点名的 313 腿 `SettingsView.Presets/configured` 现读在 `config_handlers.go:151-185`——那是**入向读答路**（页面问、Go 答），不是快照推送，与 AC#2b① 不同通路、不构成部分落地。

**剩余未落的有源字段＝5 族，逐族具名锁，全部卡在本波写面（internal/panel/**＋run.go PumpSources 字面量一处）之外 ⇒ 落地集＝空、零产码＝宁缺毋造**：

1. **usage/cost 活值**（行 2/12）：`publishUsage` 只写 `Result`（`loop.go:1115-1119`）；`agent.Event` 无 Usage/CostMicros 字段（`sink.go:51-66`），`EvDone` 只带 TokensIn/TokensOut 半回声（`:1093-1096`）；`Guard` 每任务局部（`loop.go:415`），cmd/wisp 不持有 ⇒ 装配根无活读口可接。锁＝internal/agent 事件面＋`consoleSink.Publish` 喂线（`run.go:2346-2404`，字面量之外）；§2.13 Micros/Currency 成对裁定仍未解。
2. **终结态**（行 13）：`Event.Stop` 存在（`sink.go:59`）但唯一填点＝EvError（`loop.go:553`），`EvDone` 不带 Stop；census §2.14「走 Event.Stop 还是 Result.Status」**未裁**。锁＝internal/agent＋裁定＋喂线。
3. **工具「运行中」态**（行 4）：`EvToolStart` kind＋console handler 在（`sink.go:25`、`run.go:2367`）但**生产码发布点仍 0**（本程 grep 复核）。锁＝`loop.stream` 发布点（internal/agent）。
4. **ArgsJSON/Level**（行 4/5 尾）：事件不带；store 读面在（`ToolCallByID`/`ListToolCallsByTask`）但接线点＝sink `ToolEnded` 行（`run.go:2374`）。锁＝喂线。
5. **Retryable**（行 10 尾）：`observe.Error.Retryable()` 在，喂线＝`run.go:2390`。锁＝喂线。

⇒ **给下一程派单的一句**：每族的落地点都在 `cmd/wisp/run.go` 的 **consoleSink.Publish 喂线行**（或 internal/agent 发布点／internal/config 导出面）——派单必须把这些行开进写面，否则加了就是 AC#6 答不出的装饰品。

**门禁读数**（导出树 HEAD `38628bee2`；环境卫生 unset GOTMPDIR/TMPDIR/TEMP/TMP＋PATH 加 build/sherpa-onnx，未设 CGO_ENABLED）：`/d/work/base/gopath/bin/gofumpt.exe -l internal/panel cmd/wisp` **空**；`go vet ./internal/panel ./cmd/wisp` **rc0**；build **rc0**；`go test ./internal/panel -count=1`＝===RUN **313**／---PASS **208**／---FAIL **2**，红册＝`TestC21DesignTokensFourWayAgree`（**在册**，球 token 文档族）＋`TestComposerRenderFixtureTellsTheTruth`（**环境性红**：生成件 `frontend/fixtures/composer-states.html` 不入 git、导出树缺席；自主树拷入克隆后单跑 **PASS**「verified across 3 painted states」）。★**145 在册三红现量只剩一枚**：两枚契约尺（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`）在 HEAD **已绿**＝TS 侧声明已入库（Q-51 那半张的解冻已有人落）。`sh scripts/d22scan.sh` **rc0 clean**（451 生产 Go 文件）。

**Git**：零产码；只 commit、未 push；显式 pathspec 逐枚点名（本件＋票面 Progress log 一行）；禁面（frontend/**、internal/ball/**、internal/config/**、cmd/wisp 除装配点外、thresholds/golden、docs/PLAN.md、docs/specs/**、`tokens_fourway_test.go`）零字节；票面 AC 勾与 Status 零动。
