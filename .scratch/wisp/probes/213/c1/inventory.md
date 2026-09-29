# 213/c1 — 加号菜单候选命令的能力普查（工作树现读，无一次编译/执行）

尺：全部 `file:line` 出自 `D:\work\workspace\projects plans\Wisp` 的**当前工作树**（`internal/` `cmd/` `docs/`）。
本件不读 `.scratch/wisp/probes/**` 的任何归档件，故无〔归档快照〕引用。
本件不编译、不 `go test`/`go vet`/`wisp slo`（另一枚写腿正在动 `internal/agent/approval`、`internal/tools`、`cmd/wisp`）。
`frontend/**`、`design/**` 未读、不转述；菜单**分组**按票 213 第 0.1 节属界面那支，本件只管**成员 + 能不能执行**。

图例：〔已证〕=生产路径有调用方 / 〔建了但没接〕=有实现但生产零调用方 / 〔仅文档〕=只在 PLAN.md、spec 或工单文字里被点名

---

## 1. 每一枚候选命令的现状表

| 候选命令 | 能力现状 | 干活的东西 file:line | 不新增契约面能否挂上菜单 | 归属票 |
|---|---|---|---|---|
| **compact / 上下文压缩** | 能力〔已证〕，**但只有自动触发** | 压缩本体 `internal/agent/compress.go:174`（`(*Compressor).Compress`），阈值判定 `:119 Need`；唯一生产调用点 `internal/agent/loop.go:399`，由 `:397 l.comp.Need(hist)` 驱动 | **不能白拿**：落历史要经 `(*Loop).replaceHistory`（包内私有，`loop.go:401` 是调用点），包外没有入口 => 需一枚新的导出方法 | 213（AC#1 的"真执行者"：自动算真，手动不算） |
| **goal / 长期目标** | **零命中**（连文档都没有） | Go: `grep -rn --include=*.go -iE "goal|long_term|objective" internal cmd` 排测试 => **0 行**；`docs/PLAN.md` + `docs/specs/**` 里 `goal` => 0 命中 | 无从挂起：既无能力也无配置面。**最贵的一枚** | 213 的"初始集"点名，**无人 owning** => 要另立票 |
| **plan / spec（进·退方案模式）** | 〔仅文档〕 | `docs/PLAN.md:2568` D34 行 `plan.present`（**L2**，票 165 owner 批准新增，走甲路，复用 D43 的 #17/#23/#24/#25 四行）；Go 侧 `grep -rn --include=*.go "plan.present|planPresent" internal cmd` => **0 行** | 能挂"菜单里先不列"；一旦要列 => 先落 `plan.present` 工具 => 撞 `internal/tools/**`（写腿正在动） | 213（第 0 节交付 1 的"进/退计划模式"） |
| **model switch / 换模型** | 读〔已证〕、清单〔建了但没接〕、写**零** | 单值读：`cmd/wisp/panel_pump.go:242 (*agentRuntime).currentModel`（读 `rt.endpoint.Model`）-> 装配点 `cmd/wisp/run.go:477 Model: rt.currentModel` -> 快照键 `internal/panel/composer.go:255 CurrentModel` / `:256 ModelKnown`。清单：`internal/llm/discover.go:42 DiscoverModels` 有生产调用方（`cmd/wisp/providers.go:136`，但只作为 CLI `wisp providers discover`）；**`internal/llm/discover.go:102 ImportDiscovered` 生产零调用方**（只有 `internal/llm/catalog_test.go:212`、`internal/llm/openaichat/mockllm_integ_test.go:413`）。写：**没有** —— `internal/config` 唯一的导出 setter 是 `permmode.go:64 SetPermissionMode` | 读侧已通、不需新面；**换成要新 C17 方法名**（今天只有 4 枚入向方法）=> 要批准 | **187** |
| **thinking effort / 思考档位** | 〔已证〕一路进模型，但**面板读不到、也没有档位列可列** | 配置枚举 `internal/config/schema.go:291 thinking_intensity default:"off"`、词表 `:85-88`；逐枚校验 `internal/config/validate.go:222`；解析 `internal/llm/resolver.go:196`；带进请求 `internal/agent/prompt.go:136` -> `:254`；落点 `internal/llm/provider.go:75`。模型支持面 `schema.go:359 ThinkingLevels`（校验 `validate.go:201`）。**面板侧零**：`internal/panel/pump.go:124-130` 的 `Model` 读者注释逐字写着 "not the per-adapter effort vocabulary, neither of which is exported here" | 读侧要新快照键（打红 4 枚钉）、写侧要新 C17 方法 => 两样都要批准 | **187** |
| **permission preset / 权限档位** | **〔已证〕，全仓唯一走通的一枚"用户改设置"** | 枚举 `internal/risk/mode.go:74 ModeNames`、默认最严 `:68-70 DefaultMode`；写入口 `internal/perm/store.go:175 (*Store).Set`（带确认 + 审计）；面板处理器 `internal/panel/composer_handlers.go:111 HandleModeRequest`；生产装配 `cmd/wisp/run.go:447` 与 `cmd/wisp/panel_inbound.go:215`；入向方法 `internal/panel/bridge.go:42 MethodModeRequest = "panel.mode.request"`；路由 `internal/panel/composer_dispatch.go:131`；快照回读 `internal/panel/composer.go:236 Mode ModeView`；持久化 `internal/config/permmode.go:64` | **能，零新面**：名册里已有的第 1 枚方法。这枚就是"命令名册该长什么样"的现成样板 | 213 已含；本体归票 90 / R20 |
| **export / 导出会话** | 〔建了但没接〕**而且导的不是会话** | `internal/memory/privacy.go:198 (*Store).ExportPrivacy`（生产零调用方；只有 `internal/memory/privacy_test.go:93`、`artifacts_path_invariant_test.go:415`），按**隐私域**出 JSON。⚠ **"会话"这个实体不存在**：`internal/memory/schema.go:23-136` 的表是 profile/memory/task_log/tool_call/approval_grant/cost_daily/plugin_state/provider_health，**没有 session 或 conversation 表**；对话历史只在进程内 `internal/agent/loop.go:186 history []llm.Message`；`internal/session/doc.go:17` 逐字写 `DEFERRED(SessionScope): implemented by ticket 28` | **不能**：先要有"一次会话"的实体，再谈导出 => 前置是票 28 | 213 明写"能力还没落地的不许先进名册" |
| **skills / 技能** | **零命中**（比 215 起手说的"键在没人读"还低一格） | `grep -c Skill internal/config/schema.go` => **0**；`grep -rn --include=*.go -i skill internal cmd` 排测试 => **0 行**；`grep -rni skill docs/PLAN.md docs/specs/` => **0**；`grep -rn 技能 docs/` 只命中 `docs/evidence/s1/**` 里**别的 harness 塞进来的技能清单回显**（`111-ci-step-readings.md:153`、`147-offset-naming-r1.md:486` 等），与本产品无关 | 无从挂；要挂先补配置面 => **动 D36 那棵 section 树 = 配置模型契约面** => 人工批准 | **215b** |
| **MCP servers** | 〔仅接口位〕（与 D13 一致，不是缺漏） | `internal/tools/registry.go:22 KindMCP Kind = "mcp"`（只是一枚枚举值）；`:31` 错误串 `mcp=REJECTED per D13/16.9#7`；`:52-56 providerSlots` 里 `{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"}`。配置侧 `grep -rn -i mcp internal/config/` => **0 命中**（无 section）。全仓 MCP 命中合计 **4 行、全在上面那个文件** | 能挂一条**只回显"0 台，理由 D13"**的行，**不动契约**；不能挂"开关"，也不能列真机器 | **215**（票面交付 4 明写不接）+ D13 |
| **plugins / 插件** | 三层里两层〔建了但没接〕、装配层**根本没有** | 1) 配置面〔已证〕：`internal/config/schema.go:554 PluginEntry`（`:555 Enabled`、`:561 HostAPI`）、`:565-575 PluginsSection.Entries`（动态键，`internal/config/parse.go:93 buildPlugins`）；热重载 + 松紧方向审计 `internal/config/manager.go:185-190`、`:407-433 pluginsDirection`。2) **读者零**：`grep -rn --include=*.go .Plugins internal cmd` 排测试、排 `internal/config` 自己 => **0 命中**，即 `[plugins]` 整节今天只有 config 自己在读自己。3) 装配与引擎零：`internal/config/unwired.go:126-133 lockedKeyDisposition` 对 6 枚插件键逐枚写 `not built: plugin engine is ticket 50/51`；`internal/plugin/disposal.go` 是 **C11 DisposalScope**（生命周期原语，`:18-36`），不是插件引擎，其唯一生产 importer 是 `internal/memory/retention.go:10`。4) 清单存储〔建了但没接〕：`internal/memory/models.go:117 PluginState`（含 `Enabled bool`、`Hash`、`ExeHash`）+ `internal/memory/dao_misc.go:188/222/243/269` 四枚 DAO，**生产零调用方**（只有 `dao_test.go:434-455`），而 `memory.Open` 在生产里已经开着（`cmd/wisp/run.go:347`） | 列一份**插件清单**能零新面（走 CLI 腿，见第 5 节）；要在**快照里**列 => 新键 => 打红 4 枚钉；要**开关** => 现成锁段写路径已有（`manager.go:228 applyLocked`），但从界面点需要新 C17 方法 => 批准 | **215a** |
| **subagent spawn（`task.spawn`）** | 〔已证〕，**但它是模型可自调的工具，不是命令** | `internal/tools/subagent_197.go:180 Name() = "task.spawn"`、声明 `:202 TaskSpawnDecl`、入口 `:217 BuiltinSubagentEntries`，生产注册 `cmd/wisp/run.go:533`；子代理在自己的父目录里被藏掉 `:513 subagentHiddenTool` | **不该挂**：票 213 第 0.4 节与 AC#5 要"命令只有人能执行"，禁区逐字写"别把 `/compact` 做成模型可自调的工具"。界面展示另有归属 | **197** |
| **cancel / stop** | **半枚〔已证〕+ 半枚〔建了但没接〕—— 本表最便宜的赢面** | 〔已证〕控制词：`internal/agent/control.go:36 controlWords`（停/取消/重说/大声点/确认，一比一抄 D11）、`:51 MatchControl`、生产调用点 `internal/agent/loop.go:341`（**不经 LLM 往返**）、兜底执行 `loop.go:524 defaultControl` -> `:531 root.Cancel()`。⚠ `Options.Control` 这枚导出字段（`loop.go:158`）**生产没人设**：`cmd/wisp/run.go:692 opts := agent.Options{` 里没有 `Control:`，所以只有 loop 自己当前任务停得了，repeat/louder/confirm 一律 `Handled:false`（`loop.go:535`）。〔建了但没接〕按 id 停：`internal/tools/task.go:442 (*TaskRoster).Cancel(taskID)` **生产零调用方**（grep 只命中自身定义与 `subagent_197.go:39` 一句注释），而**停的把手已经存好** —— `subagent_197.go:330 AttachCancel(bg.ID, cancelChild)` / `:391 DetachCancel`，把手表 `task.go:415` `:429`。D34 的 `task.list` / `task.cancel` 两行（`docs/PLAN.md:2564`）今天**登记为 DEFERRED**（`docs/PLAN.md:1531`：名册在册、生产注册表**零枚**，"先补名册裁定再谈实现"） | **能，零新名**：停"别人家的任务"只需把已导出的 `TaskRoster.Cancel` 接到**已导出的** `Options.Control` 上，走的正是 D11 那条"命令只有人能执行"的路；**不要**把它做成 `task.cancel` 工具（要动 D34 名册 + DEFERRED 登记表） | 213（取消当前任务）+ **197**（每枚任务的界面停） |

---

---

## 2. "命令"这个词在仓里已经有什么（票 213 该复用谁）

**先给反向尺的读数（票 213 现量那一格的复算，命令照抄可跑）**：

```
grep -rn --include=*.go -iE "slashcommand|commandcatalog|commandregistry|parsecommand" internal cmd | grep -v _test.go | wc -l
=> 0

grep -rn --include=*.go -E '"\/[a-z][a-z-]*"' internal cmd | grep -v _test.go | wc -l
=> 4   （全是 URL 与协程名，无一枚是命令：
        internal/llm/anthropic/adapter.go:146 "/messages"
        internal/llm/discover.go:56          "/models"
        internal/llm/openairesponses/adapter.go:117 "/responses"
        internal/plugin/disposal.go:232      name+"/go"）

grep -rn --include=*.go -E "type [A-Za-z_]*Registry[A-Za-z_]* " internal cmd | grep -v _test.go | wc -l
=> 4   （internal/agent/approval/approval.go:141 ChannelRegistry
        internal/observe/goroutine.go:213 Registry
        internal/risk/syncdirs_windows.go:32 liveRegistry
        internal/tools/registry.go:79 Registry）

cmd/wisp 子命令派发的枚数：grep -c  'case "' cmd/wisp/main.go  => 10
usage 块里承诺的枚数：grep -cE '^  wisp [a-z-]+' cmd/wisp/main.go => 10   （两侧相等，由测试钉住）
```

**⇒ 零枚斜杠命令目录/注册表，票 213 的现量成立。** 但"名册"这一形仓里**已经有六处**，一枚都不该重造：

| # | 已有结构 | file:line | 对票 213 意味着 |
|---|---|---|---|
| 1 | **C17 入向方法名册**（今天就 4 枚） | `internal/panel/bridge.go:42-45` 四枚导出常量（`panel.mode.request` / `panel.workspace.request` / `panel.attachment.add` / `panel.message.send`）；白名单判定 `:104 knownComposerMethod`；解析器 `:84 ParseComposerRequest`；注释 `:32-40` 写明改名或加名会红在 `internal/panel/composer_test.go:502` | **这就是"命令名册"的权威层**。213 不该新造 `Command` 类型，而该把命令**登记成这枚名册的成员**；⚠ 加第 5 枚 = **新 C17 方法名 = 契约面 = 要批准**（票 213 禁区已自认这一条） |
| 2 | **入向派发表**（四扇门） | `internal/panel/composer_dispatch.go:90-104` 四枚处理器插座、`:113 Handle`、`:129 dispatch`；`ErrNoHandlerAttached:56`、`ErrRosterMismatch:169`；`:24-30` 逐字说明为什么**不许在这里写私有路由名** | 213 要的"一条执行通道"**已经在了**，且 default 分支正是为"名册先长、派发表后学"准备的（原文："for the day bridge.go gains a fifth method and this file has not been taught it"）。三枚插座是空的 |
| 3 | **控制词名册**：唯一"人名 -> 本地执行、不经模型"的现成注册表 | `internal/agent/control.go:17 ControlVerb`、`:36 controlWords` map、`:51 MatchControl`、`:61 ControlOutcome`、`:73 ControlHandler` | ⚠ **票 213 第 0.2 节选的"同一枚解析器同时供界面高亮与回车执行"，仓里已经有这一枚**（`MatchControl` 是那枚解析器，`loop.go:341` 是那个执行点）。第 0.4 节要的"命令只有人能执行"它本来就这个形状。**213 的名册应与它同源，而不是并列** |
| 4 | **工具注册表** | `internal/tools/registry.go:79 Registry`、`:95 Register`、撞名**响亮报错**在 `:101 duplicate tool name (C1 requires global uniqueness)` | 票 213 第 0.3 节选的"撞名要响亮报（Pi 那一形）"——**这仓已经这么写了**，照抄这条规矩即可 |
| 5 | **CLI 子命令名册与派发表的一致性钉** | `cmd/wisp/main.go:22 usage` 常量（10 枚条目）、`:81-121` 十枚 `case`、三枚 `WISP-LEG-COVERAGE-RULING` 文字裁决（`:65-80`）；尺 `censusVsUsage133`（`cmd/wisp/leg_dispatch_gate_133_test.go`） | 票 213 **AC#2「没有能力就进不了册」在这一层已有可运行先例**：名册与派发表相互核对，漏一侧就红 |
| 6 | **"写了不管用就响亮拒收"旗** | `internal/config/unwired.go:56 unwiredKeys`（6 枚，每枚带 `missing` 与 `lands`）、`:110 lockedKeyDisposition`（逐键写"谁在读"）、拒收点 `:137 validateUnwired`；完整性钉 `TestEveryLockedSectionKeyIsAccountedFor` | 票 215 AC#2（配置 true 但没人装配 => 判据必红）**在这一层已经实现了**。215 不许拆（票 215 禁区），213 可直接借这一形做"名册成员必须有真读者"的账 |

**⇒ 结论给票 213**：**复用第 1、2、3、6 四枚已存在的形状，不新增名类。**

⚠ **一条硬约束，本件实测到的**：`internal/panel/composer_dispatch.go:24-30` 记录的 AST 仪器（`internal/panel/l2_grant_boundary_test.go` 的 `poolJudgedByRealGuard`）会收集全包内**任何路由形状的名字**并问运行中的守卫；原文写着，一处守卫答得出而 case 表没点名的名字会被报成 "a second switch/if chain ... **exactly the shape a handler registry takes**"。也就是说：**另起一枚平行的命令注册表会自己打红一枚现成的钉。**

---

## 3. skills / MCP / plugins 三家对照：能列出来、能开关、能看出这次生效了没有

| 维度 | skills（技能） | MCP servers | plugins（插件） |
|---|---|---|---|
| **有 loader 吗** | **零命中**。尺：`grep -rn --include=*.go -i skill internal cmd` 排测试 => 0 行 | **零命中**（无连接、无子进程、无 provider）。尺：`grep -rn --include=*.go -iE "\bmcp\b" internal cmd` 排测试 => 4 行，全在 `internal/tools/registry.go:14/22/27/31` | **零 loader 与引擎**。`[plugins]` 的**读者**为零：`grep -rn --include=*.go "\.Plugins" internal cmd` 排测试、排 `internal/config` 自己 => **0 命中**。`internal/plugin/disposal.go` 是 C11 生命周期（`:18-36`），不是引擎；唯一生产 importer 是 `internal/memory/retention.go:10` |
| **有 config section 吗** | **零**。尺：`grep -c Skill internal/config/schema.go` => **0**（比"键在没人读"还低一格，票 215 的 09-28 更正成立）。且 `grep -rni skill docs/PLAN.md docs/specs/` => **0**，即**连文档都没命名它** => 加它 = 动 **D36** 的 section 树 = 契约面 | **零**。尺：`grep -rn -i mcp internal/config/` => 无输出 | **有，是全仓最完整的一节**：`internal/config/schema.go:554 PluginEntry`（`Enabled` / `Capabilities` / `NetAllowlist` / `HostAPI`）、`:565-575 PluginsSection`（动态 `[plugins.<id>]` 表；`internal/config/parse.go:93 buildPlugins` 解、`:172 marshalPlugins` 装回去） |
| **有 enable/disable 开关吗** | 无 | 无 | **有，但只到"文件"这一层**：`plugins.<id>.enabled`（`schema.go:555-558`）+ `plugins.tier2_enabled`；锁段松紧方向审计 `manager.go:185-190` + `:407-433 pluginsDirection`；放宽走 `manager.go:228 applyLocked`（D36 的 L2 再确认）。**答"能不能开关"= 能，但改的是 config.toml，不是点菜单**；从界面点需要新 C17 方法 => 批准 |
| **有"这次生效了没有"的逐轮回读吗** | 无 | **没有，也无可回读之物**。今天能给的只有"为什么没有"：`internal/tools/registry.go:52-56` 那枚 `{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"}`，形状是 `SlotErr`（`:37`，带 `Unwrap` `:47`，可 `errors.Is`） | **两枚半成品，都没接**。1) 磁盘账：`internal/config/unwired.go:126-133 lockedKeyDisposition` 对 6 枚插件键逐枚写 `not built: plugin engine is ticket 50/51` —— 它答的正是"**配置里写了但没人装配**"，而且是**响亮**的（`:137 validateUnwired` 直接拒收）。2) 装配账的**形状已存在、内容没接**：`internal/memory/models.go:117-124 PluginState`（带 `Enabled bool`、`Hash`、`ExeHash`）+ `internal/memory/dao_misc.go:188 UpsertPluginState` / `:222 PluginStateByID` / `:243 ListPluginStates` / `:269 DeletePluginState`，**四枚全生产零调用方** |

**这一栏最该抄的"生效了没有"回读样板（本仓自己的、已上线的）** —— 三枚，全都把"没有"和"读不到"写成两句话（票 216 要的就是这一形）：

1. `internal/panel/instructions_200.go:27-43`：**五态** `loaded` / `none_found` / `off` / `refused` / `not_run`，注释逐字写"两枚不同的事实不许共用同一个缺失键"；并且 `:63-65 Dropped bool` 的注释逐字是 **"面板必须能说出这一枚没生效，不许假装生效"** —— 这就是票 215 AC#2 要的那一句，Go 侧**已经有实现**。它的生产读者也在：`cmd/wisp/panel_pump.go:117 return rt.instrLoader.Last()`，快照键 `internal/panel/composer.go:66 Instructions *InstructionsSection`。
2. `internal/panel/composer.go:236 ModeView` + `:190 ModeUnknownView`：**读不到就写 `unknown`，不许渲染成最安全的那一档**。
3. `internal/panel/composer.go:255-256 CurrentModel` / `ModelKnown`：**没人读过就 `modelKnown=false`**，不许填空串。

**⇒ 一句话交付票 215**：skills 三问全零、MCP 三问全零（且按 D13 应当为零）、plugins 只有第二问（开关）在文件层成立；第一问有**两枚半成品**（`lockedKeyDisposition` 与 `ListPluginStates`），第三问有**一枚可直接照抄的样板**（`instructions_200.go`）。

---

## 4. MCP 的取舍：D13 说了什么 vs 今天建了什么

**尺**：`grep -n "D13" docs/PLAN.md` => 标题在 **`docs/PLAN.md:320 ### D13 — MCP 取舍`**（`AGENTS.md` 第 5 节索引同样指 `:320`）。

**定案原文（`docs/PLAN.md:321`）逐字**：**「现在不做 MCP client，但架构预留 `ToolProvider` 抽象点（推迟而非封死）。」**

- 三条理由在 `:323-328`：1) 所有能力一律经 host bridge，MCP server 是独立子进程、天然绕过那唯一收口点；2) MCP 的 tool 元数据里没有风险等级，host 无法像校验 Tier 1/2 清单那样强制校验，于是要么全按 L2 强确认、要么把 D4 交给第三方；3) 一个 Node server 即 +60-150MB，直接冲 D1 资源预算。
- `:329-337` 那张 `trait ToolProvider` 图里，**`[McpProvider]` 是带方括号的"未来可接入，不需改核心"**。
- `:340-342` 要求**诚实记录代价**：2026 年"不支持 MCP"是明显缺口，须在 README 正面说明取舍而不是回避。
- 契约层对应写法：`docs/PLAN.md:1354` C4 行 `**[Mcp] 仅接口位，不实现**`；`:1533` RESERVED 行 `— **无实现计划**`，后果栏逐字 = **"无法接入任何 MCP server"**；`:1778` "保留结论、替换理由"；`:3152` 第 16.9 节第 7 项**再次驳回**引入建议；`:3188` 自己承认 D13 的理由 2 已失效（只保留 1 与 3）。

**今天建了什么**：**正好只是那枚接口位，不多一分。** `internal/tools/registry.go:22 KindMCP`（枚举值本身）、`:20-21` 注释逐字 "Only Builtin has an implementation here"、`:27-31 ErrSlotNotLanded` 把理由**带在错误串里**（`mcp=REJECTED per D13/16.9#7`）、`:52-56 providerSlots` 里那一枚 `{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"}`。

**代码与 D13 的关系**：**没有偏差** —— 预留不等于实现，仓里也确实只预留了。**本件不提任何改 D13 的建议**（票 215 禁区与票 213 第 0.5 节同样禁止；D13 也在 `AGENTS.md` 第 1.1 节的 D1-D47 禁改清单里）。

**票 215 能不能"列出 MCP servers"而不动那份契约？——能，而且代价是零。**
- `[mcp]` section 不存在（`grep -rn -i mcp internal/config/` 无输出），所以任何合法配置下的 MCP 机器数**恒为 0**；而 `registry.go:52-56` 已经把这个"0"和它的**理由**做成了**可编程读出的数据**（一枚 `SlotErr{Kind, Ticket}`）。
- ⇒ 名册里那枚 MCP 条目的**诚实内容就是这条既有 `SlotErr` 的渲染**，正好落在票 213 第 0.5 节点名的"只回显、不动手"那一形（Step-Code 的 `/mcp` 也是只回显"配置了 vs 加载了"）。
- **三条界线（写腿要守）**：1) 不加 `[mcp]` 配置 section（= 动 D36 且实际推翻 D13）；2) 不实现 `Provider` 接口的 MCP 那一枚（= 动 D13）；3) 不新增 `Kind*` 枚举值或新工具名（= 动 C4 与 D34 面）。

---

## 5. 下一枚写腿的最省事落点

| 落点 file:line | 改什么 | 为什么不需要新 C17 方法名、也不需要新导出名 | 判定 |
|---|---|---|---|
| `cmd/wisp/run.go:692 opts := agent.Options{` | 把**已存在但没人设**的 `Control:` 装上（一枚 host handler），内部调 `rt.tasks.Cancel(taskID)`；顺带让 repeat / louder / confirm 得到真执行者，或**具名报"没落地"** | 插座与把手**全是现成导出面**：`Options.Control`（`internal/agent/loop.go:158`，注释逐字写 "nil = the loop's own default"）、`ControlHandler` / `ControlVerb` / `ControlOutcome`（`internal/agent/control.go:73` / `:18` / `:61`）、`TaskRoster.Cancel`（`internal/tools/task.go:442`，把手已由 `internal/tools/subagent_197.go:330 AttachCancel` 存好）。**零新名、零新方法名** | **最便宜**。一次接掉两枚〔建了但没接〕；直接喂票 213 的"取消当前任务"与票 197 的"每枚任务停得了" |
| `internal/panel/composer_dispatch.go:96-104`（`Workspace` / `Attachment` / `Message` 三枚空插座） | 把已实现的守卫接到空插座后面：workspace 走 `internal/panel/workspace.go:76 RequestWorkspaceSwitch`，attachment 走 `internal/panel/attachments.go:179 NewAttachmentBroker` + `:199 Ingest` | 接口**已声明**（`composer_dispatch.go:64` / `:73` / `:80` / `:86` 四枚 handler 接口），方法名**已存在**，入向方法名**已在名册里**（`bridge.go:43` / `:44` / `:45`）；`:24-30` 的注释逐字写明缺的是接线不是名字 | **便宜，零新名**。同一片文件正是票 213 的落点，先接它不会白改 |
| `internal/panel/pump.go:226 NewComposerState(mode, workspace, nil, maxAttachment)` | 把写死的 `nil` 换成真实附件读；上限一侧字段已有（`pump.go:162 AttachmentMax`，零值走常量 `attachments.go:47 MaxAttachmentBytes`） | **不加任何快照键**：四枚附件键早已存在（`internal/panel/composer.go:238-241`），`NewComposerState` 的形参位也已有（`composer.go:262` 第 3 位）=> 不触发票 213 现量点名的 4 枚钉（`composer_test.go:48`、`approval_test.go:105`、`pump_test.go:111-124`、`:270-276`）、不动 `panel.ts`（Q-51） | **便宜，零新名**（票 214 的落点） |
| `cmd/wisp/` 新增一枚**清单腿**（照 `wisp panel-inbound` 那一形，或扩 `wisp doctor`），枚举 `[plugins]` + `plugin_state` + MCP 那枚 `SlotErr` | 一份可枚举清单：名字 / 版本 / 开没开 / 为什么没生效 | 包是 `package main`，**导出名为零**；不动 C17、不动快照键 => 不打红 4 枚钉、不碰 Q-51。数据源全现成：`config.Plugins.Entries`（`internal/config/schema.go:575`）、`ListPluginStates`（`internal/memory/dao_misc.go:243`）、`SlotErr`（`internal/tools/registry.go:37`）、`lockedKeyDisposition`（`internal/config/unwired.go:110`） | **便宜，但有一枚前置钉要守**：新 `case` 必须同批改 `cmd/wisp/main.go:22 usage`，并被 `censusVsUsage133`（`cmd/wisp/leg_dispatch_gate_133_test.go`）核对；且按票 133 的规矩要么自带测试、要么写 `WISP-LEG-COVERAGE-RULING` 文字裁决（样板在 `main.go:65-80`） |
| `internal/agent/compress.go:174` 的手动触发入口（`/compact`） | 一枚"现在就压"的方法 | ⛔ **要停手上报**：落历史只有包内私有的 `(*Loop).replaceHistory`（唯一可见调用点 `internal/agent/loop.go:401`），包外无入口 => **必须新增一枚导出方法** | **⚠ 需人工批准（新导出名）**。另注意 `DEFERRED(D28-1)` 已登记在 `compress.go:25-27`（Warm-window hook 才是归宿），手动入口别和它抢 |
| `panel.model.request` 与 effort 档位请求（换模型、换思考档） | 让票 187 的"能改"落地 | ⛔ **三重都要批准**：1) 第 5 枚入向方法 = **新 C17 方法名**（`bridge.go:42-45` 是闭合集，加名红在 `composer_test.go:502`）；2) 快照里列清单或档位 = **新键**，打红上面那 4 枚钉，且 `panel.ts` 归 Q-51；3) 写侧只有 `Manager.SetPermissionMode`（`internal/config/permmode.go:64`）一枚导出 setter，模型与档位没有写路径 | **⚠ 需人工批准，且属票 187**，票 213 不许顺手做 |
| skills 的配置面（`[skills]` 或 `skills.<id>.enabled`） | 让"技能"第一次存在 | ⛔ **动 D36 那棵 section 树 = 配置模型契约面**；D36 在 `AGENTS.md` 第 1.1 节的 D1-D47 禁改清单里 | **⚠ 需人工批准（批准只落 `A##`、不改 `PLAN.md` 一字），属票 215b** |
| 把任何命令做成模型可自调的 C4 工具（`/compact` 做成工具、`task.cancel` 做成工具） | —— | ⛔ 票 213 禁区逐字："别把 `/compact` 做成模型可自调的工具"（与 D34 无 git 行同源）；`task.list` / `task.cancel` 也已**登记为 DEFERRED**（`docs/PLAN.md:1531`，"先补名册裁定再谈实现"），做成工具要同时动 D34 名册（`docs/PLAN.md:2564`）与 `SPEC-12` 第 5 节登记表（`AGENTS.md` 第 1.1 节要求 1:1 双向对得上） | **⛔ 不要走这一条**。本表第一行那枚 `Options.Control` 落点是它的零契约替身 |

**一句话给编排者**：**唯一一枚"零新名、零契约面"的落点是 `cmd/wisp/run.go:692` 那行 `agent.Options{}` 里补上 `Control:`** —— 它一次接掉两枚〔建了但没接〕（`Options.Control` 与 `TaskRoster.Cancel`），喂饱票 213 的"取消当前任务"，并天然满足票 213 第 0.4 节那条"命令只有人能执行"，因为它走的本来就是 D11 那条不经 LLM 的路。

**三票共用的前置（本件复核为真，未翻案）**：界面到 Go 那一跳今天只有命令行路 —— `internal/panel/composer_dispatch.go` 的注释逐字写它 "has NO production caller yet"，直到 `cmd/wisp/panel_inbound.go` 那枚 CLI 腿（`cmd/wisp/main.go:112-118 case "panel-inbound"`）把它接上，而该腿自己写着 "it is not the WebView2 receiver (H2/H3). The page's postMessage still has no listener in this tree"。⇒ 票 181 / 33 的入向线接通之前，213 / 214 / 215 的终态只能是"Go 侧可枚举 + 命令行可执行 + 界面看不到"。

---

## 6. 行号更正与漂移声明（**正文以本节为准**；本件写完后逐条复查所得）

### 6.1 两枚前提被本件推翻或收窄

1. **⚠ 派单给我的前提 (b) 措辞错了**：「四枚**附件配置键**已经存在」**不成立**。
   `grep -rn -i attachment internal/config/` ⇒ **零命中** —— `config.toml` 里没有任何附件键。
   真实形状（＝票 214 09-28 更正的原话）是四枚**快照 JSON 键**已存在：`internal/panel/composer.go:238-241`
   `attachments` / `acceptedAttachmentMimes` / `maxAttachmentBytes` / `attachmentError`；
   而"上限"不是配置，是 Go 常量 `internal/panel/attachments.go:47 MaxAttachmentBytes int64 = 64 << 20`，
   白名单也不是配置，是硬编码表 `attachmentTypes` 经 `:153 AcceptedMIMETypes()` 吐出。
   ⇒ 结论不变（不加键就不打红 4 枚钉），但**"配置里能调附件上限"这句话今天对 owner 不许说**。
2. **⚠ 派单给我的票号表（213/214/215/216/186/187/197/198）已经过时**：`ls .scratch/wisp/issues/` 现有
   **219 / 220 / 221** 三枚新票，且**"cancel 那格"的归属是 221，不是"213＋197"**。
   `docs/../issues/221-task-spawn-description-promises-task-cancel-that-is-not-registered.md:1` 已 owning
   本件在第 1 节判为"最便宜"的那枚〔建了但没接〕（`TaskRoster.Cancel` 生产零调用者，票 221 现量表第 3 行逐字同判）。
   ⇒ 本件第 1 节"cancel/stop"那一行的**归属票栏应读作：221（本体）+ 213（命令成员资格）+ 197（界面）**。
   ⚠ **但本件第 5 节那枚落点仍然是 221 没列的第三种形状**：221 只写了
   **甲形＝注册 `task.cancel` 工具**（碰 D34 权威表 + `PLAN.md:1531` 的 DEFERRED 登记 ⇒ 要先落 `A##`）与
   **乙形＝删掉 `subagent_197.go:189` 那句对模型的许诺**（纯代码）；
   本件提的"把 `Options.Control` 装上、内部调 `TaskRoster.Cancel`"**既不动 D34 也不动 DEFERRED 登记表**（它走 D11 那条"命令只有人能执行"的路，正合票 213 第 0.4 节），
   因此**要并入 221 当第三支、还是另立一枚，属编排者的决定，不是本件能定的**。
   并入时两条 221 的既有约束必须一起带过来：AC#2 的"停一名孩子后它自己的流要有终态、要留下被谁停的、名册状态只能用 D43 已有的名"；
   AC#3 的"子代理停兄弟／停自己两枚都必须被拒"；以及 221 禁区那句"票 201／220 的写面未空出之前不许派甲形（同撞 `internal/tools/task.go` + `subagent_197.go` + `cmd/wisp`）"。

### 6.2 工作树漂移：本件引用过的两处正在被别人改写

`git status --porcelain` 现读（只读命令，本件没动过任何源文件）：
`M cmd/wisp/main.go`、`M cmd/wisp/run.go`、`M internal/agent/approval/gate.go`、`M internal/agent/approval/queue.go`，
未跟踪新文件 `cmd/wisp/approval_reply.go` / `_stdin_other.go` / `_stdin_windows.go`。
`git diff --numstat`：**`cmd/wisp/main.go +21/-1`、`cmd/wisp/run.go +98/-4`** ⇒ 这两枚的行号**在我这次会话内就漂过**
（`opts := agent.Options{` 我起手读到 `:692`，收尾重读已是 `:742`）。
`internal/panel/composer.go`、`internal/tools/task.go`、`internal/tools/subagent_197.go`、`internal/config/**` 的 `git diff --numstat` **无输出 ⇒ 未被改动，行号可信**。
⇒ **引用 `cmd/wisp/**` 的落点一律要现跑 `grep -n` 复量再写**，别抄本件（也别抄工单）的行号。

### 6.3 逐条更正（正文写错的行号，一律以此表为准）

| 正文误引 | 现读正确值（工作树，`grep -n` 复量） |
|---|---|
| `loop.go:397 l.comp.Need(hist)` | `internal/agent/loop.go:393` |
| `loop.go:401 replaceHistory 调用点` | `internal/agent/loop.go:403`；定义在 `:1016` |
| `loop.go:186 history []llm.Message` | `internal/agent/loop.go:187` |
| `loop.go:504/524/531/535`（defaultControl 一族） | `:504 runControl`、`:524 defaultControl` 正确；`:531 root.Cancel()`、`:535` 未复量，引用前现跑 |
| `compress.go:25-27 DEFERRED(D28-1)` | `internal/agent/compress.go:27`（loop 侧那半句在 `internal/agent/loop.go:394`） |
| `PLAN.md:2568 plan.present` | `docs/PLAN.md:2567` |
| `session/doc.go:17` | `internal/session/doc.go:16` |
| `registry.go:20-21 "Only Builtin has an implementation here"` | `internal/tools/registry.go:14` |
| `registry.go:101 duplicate tool name` | `internal/tools/registry.go:103` |
| `registry.go:37 SlotErr` / `:47 Unwrap` | `internal/tools/registry.go:35` / `:48` |
| `composer_dispatch.go:64/73/80/86` 四枚 handler 接口 | `internal/panel/composer_dispatch.go:67` / `:74` / `:81` / `:86` |
| `composer_dispatch.go:90-104` 四枚插座 | 结构体在 `internal/panel/composer_dispatch.go:97`（字段 `:100-110` 一带） |
| `composer_dispatch.go:56 ErrNoHandlerAttached` / `:169 ErrRosterMismatch` | `:62` / `:176` |
| `composer_dispatch.go:113 Handle` / `:129 dispatch` / `:131` mode case | `:120` / `:137` / `:139` |
| `unwired.go:56 unwiredKeys` / `:110 lockedKeyDisposition` / `:137 validateUnwired` | `internal/config/unwired.go:60` / `:119` / `:100` |
| `unwired.go:126-133` 六枚插件行 | `internal/config/unwired.go:141-146` |
| `schema.go:555-558 Enabled` / `:565-575 PluginsSection` | `internal/config/schema.go:555`（注释）`:556`（字段）/ `:572-575` |
| `instructions_200.go:27-43` 五态 / `:63-65 Dropped` | `internal/panel/instructions_200.go:26-40`（`Loaded :27`、`NotRun :40`）/ `:59` |
| `composer.go:66 Instructions` | `internal/panel/composer.go:74`（`Tasks` 在 `:91`；`ModeView` 定义 `:141`、`ModeUnknownView :168`） |
| `composer.go:190 ModeUnknownView` / `:262` NewComposerState 形参位 | `:168` / `:260` |
| `cmd/wisp/run.go:447 rt.modeWrites` | **`:479`（收尾现读，会话内已漂）** |
| `cmd/wisp/run.go:347 memory.Open` / `:477 Model:` | **`:369` / `:509`（收尾现读）** |
| `cmd/wisp/run.go:375 / 395 / 533` 三处注册循环 | **`:397` / `:417` / `:565`（收尾现读）** |
| `cmd/wisp/run.go:692 opts := agent.Options{` | **`:742`（收尾现读）** |
| `cmd/wisp/main.go:112-118 case "panel-inbound"` | **`cmd/wisp/main.go:109`** |
| `cmd/wisp/main.go:81-121` switch / `:65-80` 三枚裁决注释 | **`:82` 起** / **`:68-80`** |
| `subagent_197.go:189`（票 221 的谎句）与 `task.go:595`（只注册 `task.output`） | 两处**已现读核实**：`internal/tools/subagent_197.go:189` 逐字「可以用 task.cancel 单独停它」；`internal/tools/task.go:595-597 BuiltinTaskEntries` 返回**一枚** `taskOutput` |

### 6.4 本件自己没做到的两件事（照票面口径具名报"没测"，不许留成空白）

1. **票 216 起手那枚反向尺本件没跑完**（不在我的射程内，交回给下一格）：
   `grep -rn --include=*.go -iE "sanitiz|control.?char|escape.*terminal|\\u001b" internal/ cmd/`。
   本件只顺手看到一枚相邻件：`internal/panel/attachments.go:165 NameGuard`（管名字合法、不管显示安全），与票 216 现量表一致。
2. **本件没有编译证据**。按派单硬约束没跑 `go test` / `go build` / `go vet` / `wisp slo`，
   因此全部结论是**源码层"有无生产调用方"的静态读数**，不是运行时行为读数；
   票 213 AC#1「每一条都有真执行者」这类判据最终仍需写腿在终态跑整包点名比对（221 的 AC#5 那形）。
