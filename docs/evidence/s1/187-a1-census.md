# 187-a1 — 票 187 只读普查：模型目录/选择器从哪取、写腿落在哪几行、C17 要不要动、与 190/188/197 怎么分账

- 腿：`187-a1`（只读普查·零产码）。工单＝票 187 `AC#1`／`AC#2`／`AC#4` 的落地前料刷新。
- 锚：`HEAD a34ad0a7f`（分支 `dev`），2026-11-04 16:1x–16:3x +08，本程现跑。
- 产出：本文件一枚（仓库内唯一写点）。临时件在 `D:\tmp\wisp\20261104\187a1\readings.md`（只建不删）。
- ⚠ 本程**零产码**、未改任何 `*.go`、未改票面、未动 `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／C17 白名单／`frontend/**`。
- 工作树读数：`git status --porcelain` **1322 行全部是 `??`（未跟踪）**，跟踪件零改动 ⇒ 本程的 `git add` 只可能带进自己这一枚 md。
- 标〔现跑〕＝本程跑过命令/测试；标〔读断言〕＝只读源码/断言文本、未执行；行数是**读后估**，落地腿要复量。

---

## §0 结论摘要（先给答案）

1. **模型目录／选择器不需要新读面，也不需要重造**（回答工单 ①）：
   - **可选项 id（含 Enabled 过滤）**：`(*llm.Resolver).DiscoveredModels(provider)` **已导出**——`internal/llm/resolver.go:287-305`，
     过滤 `Enabled==false`、按名排序、未知 provider 回 nil；判据在 `internal/llm/enabled_reach_261_test.go:176-183` 与 `enabled_gate_261_r1_test.go:50-66`。
     **生产调用者今天 0 枚**（`grep -rn DiscoveredModels` 只中 resolver.go 自身与那两枚测试）⇒ 是"**有导出 reader、没人接**"，**不是**"无导出 reader"。
   - **当前值**：`cmd/wisp/panel_config_store.go:96-112 ReadSettings` 已经在读 `cfg.LLM.Roles.Chat.{Provider,Model}`，
     并且 `:134-140` 已经枚举 `cfg.LLM.Providers[p].Models` 的全量 id，经 `panel.SettingsView.Providers[].Models`
     → `renderSettingsPayload`（`internal/panel/config_handlers.go:585-644`）→ 页侧 `config-view.ts:51-66` **已经在用**（钉在 `cmd/wisp/panel_config_313_test.go:183-184`）。
   - **每模型声明的思考档**：`config.ModelSpec.ThinkingLevels`（`internal/config/schema.go:362-364`）是**已导出结构体上的导出字段**，
     直接从 `cfg.LLM.Providers[p].Models[m].ThinkingLevels` 读，**不需要任何新导出函数**；今天缺的只是"没人把它读进 payload"。
   - **当前思考档**：`cfg.LLM.Roles.Chat.ThinkingIntensity`（`schema.go:296`），同样 1 行就能读。
   - ⇒ **结论：不新读面、不重造目录。** 票面 `现量` 表那句"缺的不是数据源"是对的；但 `cmd/wisp/panel_pump.go:299-301` 与
     `internal/panel/composer.go:286-290` 里"**目录没有导出 reader**"这句**过强**（导出的是 `DiscoveredModels`，只是没人接），落地腿应顺手把它改准。
2. **写腿也不需要新 C17 方法**（回答工单 ②）：**模型今天已经能从面板改**——
   `config.set` + `role_chat_model`（`internal/panel/config_handlers.go:64`，248 表）／`role_chat_pair`（`config_write_313.go:89`，313 名册）
   已经端到端可写（`cmd/wisp/panel_config_store.go:242-243`／`:276-277` → `internal/config/settings.go:224 SetRoleChatModel`／
   `settings_313.go:208 SetRoleChatPair`）。**唯一缺的是思考档那一枚字段**：在 `fields313` 里追加一行 `role_chat_thinking` 即可，
   **0 新方法、0 新顶层键、0 新 import 边**。
3. **但"能改"今天等于"改了没用"**（本程最有价值的一枚发现）：`[llm.roles.chat].thinking_intensity` **到不了线上**。
   `llm.ResolveRole` 把它抄进 `RoleSettings`（`resolver.go:202-206`），`execute()` 每任务重读一遍（`run.go:1472-1478`），
   但**只有 `MaxOutputTokens` 被交给 `agent.Config`**（`run.go:1534`）；`agent.Config`（`loop.go:110-160`）根本没有这个字段，
   `Assembler.ThinkingIntensity`（`prompt.go:136`）**全仓非测试码零赋值**（`NewAssembler` 只 copy 了 MaxOutputTokens，`loop.go:260-261`）。
   ⇒ 今天往 config 写一个思考档，`llm.Request.ThinkingIntensity` 恒为 `""`，**一次请求都不会变**。
   票面 `现量` 第 2 行"它今天真的被用（不是死字段）"只对到 `RoleSettings` 为止，**对到线上是假的**。
   修法有现成模板：票 370 给 `[llm.roles.chat].max_output_tokens` 做的就是这一跳（`internal/agent/maxoutput_370_test.go:1-60` 的注释把三种"值存在但不旅行"的形状写死了），**3 行 + 一发判据**。
4. **模型换"当场生效"要动的那一跳到今天为止不存在**：`rt.endpoint` 只在启动时赋值一次（`run.go:517`），
   `execute()` 读的是那个字段（`:1463 ep := rt.endpoint`），快照的 `composer.currentModel` 也读它（`panel_pump.go:304-309`）。
   ⇒ 面板写完模型、1s tick 把 `[llm]` 热应用进 Manager（`tiers.go:40` `"llm": TierHot`；tick 在 `config_reload.go:100-101`），
   **cfg 变了、线上模型没变、两处显示也都还是旧值**。这与回执说"要重启"（`panel_config_store.go:288`）自洽，但与票面要的"当场改"矛盾。
   落地建议见 §3.2 案型 **乙**（每任务重解 endpoint + 重建 provider，值 18-26 行）。
5. **判语**：本票可以**不新增 C17 方法、不新增顶层快照键、不改 `composer.go`** 地做完；
   真正的成本在①那一跳思考档接线（3 行，必做）、②一枚 config setter（约 30 行）、③读面补字段（约 30 行）、
   ④模型"当场换"（18-26 行，可裁），加上约 140-200 行判据与 3 处既有断言的**具名改钉**（§4.3）。

---

## §1 ① 模型目录／选择器：从哪取（复用 vs 新读面）

### 1.1 今天存在三条"能拿到目录"的路，只有两条是给面板的

| 路 | 位置 | 取到什么 | 生产调用者 | 能不能直接用 |
|---|---|---|---|---|
| **A. `llm.(*Resolver).DiscoveredModels`** | `internal/llm/resolver.go:287-305` | 某 provider 下 `Enabled==true` 的模型 id，排序 | **0**（只在 `internal/llm/enabled_*` 两枚测试里） | **能**，正是"选择器能选的那几枚"的定义 |
| **B. 设置腿 `ReadSettings`** | `cmd/wisp/panel_config_store.go:134-140` → `internal/panel/config_handlers.go:585-644` → 页 `config-view.ts:51-66` | **全量** id（含 `Enabled=false`）、provider 行身份 | 生产：`panel_inbound.go:302-306`／`panel_resident_windows.go:228-238` 两条装配都在 | **已在线上，但不筛 Enabled** |
| C. 拉取腿 `config.models.list` | `internal/panel/bridge.go:75` ＋ `cmd/wisp/panel_models_320.go` | provider 自己的 `/v1/models` 网络结果（**不是**本地目录） | 生产：同上 | 不该当选择器：AC#4 要的是"目录与角色表里真存在的项"，网络结果会把目录外的 id 拉进下拉 |

⇒ **判据落点**：选择器的 id 集合**必须**走 A（或 B 改成与 A 同一条过滤），**不许**走 C，**不许**在 `internal/panel` 里再写一份 `Enabled` 过滤（那就是第二份名册）。
最省且不重造：`ReadSettings` 里把 `for id := range p.Models` 换成 `llm.NewResolver(cfg, nil).DiscoveredModels(name)`——**`internal/config` 与 `internal/llm` 之间没有环**
（`internal/llm` 已 import `internal/config`，`resolver.go:10`）。

### 1.2 逐项：要什么、在哪、缺什么

| 选择器要的东西 | 今天在哪 | 缺什么 | 代价 |
|---|---|---|---|
| 可选模型 id（筛 disabled） | 见 A；B 里已有全量 | B 没筛 | 1 行（换成 `DiscoveredModels`） |
| 每模型声明的思考档 | `schema.go:362-364 ThinkingLevels` | 没人读进 payload | 读 1 行 + payload 1 键 |
| 当前模型 / 当前思考档 | `panel_config_store.go:98-103`（模型已有）；`schema.go:296`（思考档） | 思考档没读 | 1 行 |
| "没有可用模型"这一形 | `configured`／`providers:[]`（`config_handlers.go:585-644`） | 有 `providers:[]` 但没有"某 provider 有行、目录 0 模型"的显式文案位 | payload 1 键 + 1 句文案（宁缺毋造：空就是空，不许回退到"全部模型"） |
| 当前模型在运行时用的**到底**是不是这一枚 | `panel_pump.go:304-309` 读 `rt.endpoint.Model` | 见 §3.2：`rt.endpoint` 启动后不再变 | 见 §3.2 |

### 1.3 不需要的那几件（明确写出来，防落地腿顺手加）

- **不需要**在 `internal/llm` 新开 `catalog.go` 或 `ListModels()`：A 已经够，且 `Resolver.Providers` 本身就是导出字段（`resolver.go:83`）。
- **不需要** `config.PresetNames()` 那类新导出：已在（`internal/config/defaults.go:39`，设置腿 `:101` 已在用）。
- **不需要**把 `ThinkingLevels` 的枚举再校一遍：`validate.go:238-246` 已经在加载期校过每个元素；面板只负责**列**。
- **不需要**碰 `frontend/**`（`A383`：本编队零写面）；payload 新增键对 `config-view.ts` 是**容错**的（该文件 `:105-108` 逐字"unknown fields are ignored on purpose (the wire shape is the Go side's to grow)"）。

---

## §2 ② 落点：文件／约行数／C17／顶层键

### 2.1 落点总表（生产码，不含判据）

| # | 文件 | 改什么 | 约行数 | 备注 |
|---|---|---|---|---|
| P1 | `internal/agent/loop.go` | `Config` 加 `ThinkingIntensity string`；`New` 里 `asm.ThinkingIntensity = opt.Config.ThinkingIntensity` | **10-14** | 照 `:260-261`（票 370）的形状；**必做**，否则本票是"改了没用" |
| P2 | `cmd/wisp/run.go` | `agent.Config{...}` 里补 `ThinkingIntensity: roleSettings.ThinkingIntensity`（`:1528-1534` 那一块） | **3-5** | 与 `:1472` 的每任务重读天然咬合 ⇒ 下一任务生效 |
| P3 | `internal/config/settings.go` | `SetRoleThinkingIntensity(level string)`（`writeOneKey` 形状，照 `:224-230`）＋ 建议加 `requireModelDeclaresLevel`（model 的 `ThinkingLevels` 非空时必须含该档） | **16-22 ＋ 10-12** | 第二半可选；不加就是"面板筛、文件不筛"（见 §4 雷 6） |
| P4 | `internal/panel/config_write_313.go` | `FieldRoleChatThinking = "role_chat_thinking"` ＋ `fields313` 一行 ＋ `fields313TakesValue` 一行 ＋ 表头一句 | **8-10** | **不要**往 248 那本表加（`config_313_test.go:55` 钉死 248 恰好七枚） |
| P5 | `cmd/wisp/panel_config_store.go` | `ApplySetting` 加一分支（照 `:242-243`）；`ReadSettings` 加 `ChatThinking`／`ChatThinkingKnown`／每模型 `ThinkingLevels`／disabled 名单（`:96-159`） | **19-26** | 无新 import（`config` 已在） |
| P6 | `internal/panel/config_handlers.go` | `SettingsView` 加 3 字段、`ProviderSettings` 加 2 字段、`renderSettingsPayload` 加键（`:151-230`／`:585-644`／`:649-668`） | **26-36** | 纯追加；每字段都要"known"位（票 92/248 的既有口径） |
| P7 | `cmd/wisp/panel_config_store.go` | 按字段定档：思考＝`EffectiveNextTask`；模型＝见 §3.2 裁 | **6-8** | 现在是全党 `EffectiveRestart`（`:288`） |
| **必做合计** | 6 枚文件 | | **≈ 82-121 行** | |
| P8（可裁） | `cmd/wisp/run.go` | `execute()` 里每任务重解 endpoint ＋ 重建 provider ＋ 更新 `rt.endpoint`（`:1463`／`:1490-1494` 一块） | **18-26** | 见 §3.2 案型乙；会让 `quota_44r5_test.go:182` 那枚源码字面量钉变红（§4 雷 7） |
| **含 P8 合计** | 6-7 枚文件 | | **≈ 100-147 行** | |

判据（另计，**140-200 行**）：`internal/agent/` 照 `maxoutput_370_test.go` 的"Config → Assembler → Request → 线上 JSON body"四跳各一发；
`internal/config/` 一发 setter（枚举拒绝＋model 声明的档拒绝）；`internal/panel/` 一发路由（字段进名册、值不落 `risk.Mode`）；
`cmd/wisp/` 一发端到端（面板写 → 文件 → 下一任务真的用新档）。

### 2.2 C17：**0 新方法**（判语：不需要）

- 面板入向白名单在 `internal/panel/bridge.go:41-76`，封闭集是 `knownComposerMethod`（`:154-160`）**七枚**：
  `panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`／`config.get`／`config.set`／`config.models.list`。
- `panel.*` 那四枚被 `internal/panel/git_test.go:497-513` **逐字钉成恰好四枚**；本票不加方法 ⇒ 这枚钉**不红**。
- 本票走的是 `config.set` 这条**已经存在**的门，增长的是它的**字段名册**（数据），不是方法表。
  ⚠ 但 `config_handlers.go:16` 逐字写着"**Growing this table is a contract change, not a bug fix**" ⇒
  **加一行字段＝契约面变更（记账级，非 C17 级）**，要编排者落 `A##`、并由机主批"契约文字"那一步；
  票 313／320 已两次按这个形状做过（`config_write_313.go` 的独立名册 ＋ `A877` 裁定），有先例可循。
- 票面 AC#2 说"若要新增 C17 白名单条目，先具名报我"——**本程的结论是不需要新增**，故不触发那个停手点。

### 2.3 顶层快照键：**不加**（判语：不加，也无需加）

- 顶层键集今天被四处钉住：`Snapshot` 四枚（`composer.go:57-92` 的 `pending`／`results`／`composer`／`generatedAt`，另两枚 `instructions`／`tasks` 是 `omitempty` 指针），
  字节级两发 `internal/panel/pump_test.go:111-125` 与 `:263-293`（`want exactly the four`），两发两向对账 `composer_test.go:52-85`／`approval_test.go:105-134`。
- 本票**不需要**新顶层键：选择器与当前值都挂在 `config.get` 的 payload 上，快照那边一个字节不用动 ⇒ **不新增撞 `TestComposerContractTypesMatchFrontend` 的面**。
- **连 `composer` 节也不用动**（推荐）：`currentModel`／`modelKnown` 已在（`composer.go:286-295`）；
  想加 `currentThinking` 也行，但那会让 `ComposerState` 的对账多一枚"Go 发了、页侧没声明"的键——
  该测试**今天已经是红的**（`files`，见 §4.2），再加一枚不会变出新颜色，只会把红句加长，且页侧归 UI 腿，本编队写不了。
  ⇒ **推荐把"当前思考档"只放进 `config.get` payload**（页侧解析器已容错），快照保持不动。

---

## §3 ③ 与 190／188／197 的边界

### 3.1 逐票分账（谁碰哪几枚文件）

| 票 | 它的写面 | 与 187 的重叠 | 处置 |
|---|---|---|---|
| **190**（文件树，已落地 6 commit，`A1491`） | `internal/panel/composer.go`（`ComposerState.Files` `:285`）、`tree_types_190.go`／`tree_reader_190.go`、`cmd/wisp/panel_pump.go`（`fileTreeView` `:287-292`）、`cmd/wisp/run.go`（`PumpSources.Files` `:909`）、**页侧 `panel.ts` 的 `files` 键（未写，红的成因）** | **`composer.go` 那一枚键** | 187 **不碰 `composer.go`**（见 §2.3）⇒ 无重叠；`files` 那枚红归 190 的页侧，⛔ 187 不许"顺手补" |
| **188**（任务监控，普查已交 `188-a1`） | 它的余格 `AC#2-e` 生产者＝`cmd/wisp/run.go` `execute()` 收口处（`TaskBackfill{...State}`，`188-a1` §3.1 指 `:1666-1680`；本程未复跑该行号） | 同为 `cmd/wisp/run.go`，**不同区段**（187 在 `:1463-1534` 的装配段） | 可同批但**必须分两次提交**；`cmd/wisp` 是共享编译面 ⇒ 每波只许一枚写腿（`AGENTS.md` §时代索引③） |
| **197**（子代理三层） | `Snapshot.Tasks`（`composer.go:91`）、`internal/panel/subagent_roster_197.go`、`internal/panel/pump.go`、`cmd/wisp/run.go:930`、`cmd/wisp/panel_pump.go` 的 `taskRosterState` | **`composer.go`**（它已占）＋ `run.go`（不同区段） | 187 不碰 `composer.go`／`subagent_roster_197.go`／`pump.go` ⇒ 与 197 零文件重叠；只共用 `run.go` 一枚文件 |
| （对照）**186**（git，已落） | `bridge.go` 白名单**没有**增方法，`git.go`／`git_test.go` | 无 | 187 不碰 `bridge.go` |

### 3.2 票内那两条"当场改"路线怎么分（本程的选型建议）

**思考档（必做，便宜）**：P1+P2 三行把值接到线上，`execute()` 每任务重建 loop（`agent.New(loopOpts)` 在 `:1585`）⇒
**下一任务用新档**，回执该说 `EffectiveNextTask`（`config_handlers.go:504-505` 那句现成："这一项在下一次任务开始时生效。"）。零重启、零竞态窗口。

**模型（要裁：甲／乙）**：

| 案型 | 写什么 | 约行数 | 线上读数 | 代价/红 |
|---|---|---|---|---|
| **甲（保持重启）** | 什么都不做，回执就说重启（今天就是这样，`:288`） | 0 | 改了要重启；`composer.currentModel` 与 `config.get.chatModel` 会**长期分叉**（前者=启动模型，后者=文件里的新值） | 与票面"当场改"正面冲突；AC#5 的"显示=真用"只能靠一句"两处各报各的"糊过去 ⇒ **不建议** |
| **乙（下一任务换，推荐）** | `execute()` 内：`res := llm.NewResolver(cfg, rt.settings.secrets)` → `ResolveRole` 得 `ep2`；`BuildEndpointProvider(ep2, chainOptions(cfg))` 得 `prov2`；`rt.endpoint = ep2`；把 `rt.provs`／`rt.names` 换成 `{prov2}`／`{ep2.Provider+"/"+ep2.Model}`；回执改 `EffectiveNextTask` | **18-26** | 面板写完 → 1s tick 热应用（`tiers.go:40`）→ 下一任务起真用新模型；快照 `currentModel` 同步跟上（它读 `rt.endpoint`） | 会打红 `cmd/wisp/quota_44r5_test.go:182` 那枚**源码字面量**钉（它逐字要 `quotaDispatchChain(cfg, rt.provs, ...)`）⇒ 改钉要具名；另有一个**真实代价**：每任务重建 `llm.BuildEndpointProvider` 会把 `BucketLimiter`（`resolver.go:284`）重建，**RPM 节奏跨任务归零**，快速连发可能超 provider 的 RPM ⇒ 要加一枚按 `provider/model` 缓存（+6-10 行）或把 limiter 提到 runtime 上 |

**为什么推荐乙（按机主铁律：最有价值＋最主流 harness）**：
主流 harness 换模型是"下一条消息就用新的"，不是"重启程序"；而本票的**标题**就是"当场改"。
甲会让这一页最显眼的两个控件里，一个（思考档）能当场生效、另一个（模型）永远要重启——同一排控件两种语义，正是票 146 那族"显示与真用分账"的形状。
乙的代价是**一枚源码钉改判**＋**一个限流器缓存**，买到的是"两个控件都当场生效、快照与文件不再分叉"。
⚠ 乙**不**改变 `D15` 的语义：`BudgetsFor(window)` 在 `agent.New` 里按 `ep.ContextWindow` 冻结（`loop.go:232-237`），换窗口只在**新任务**上生效，**不会**让正在跑的那一轮超限（那一支的答案＝"等到轮次边界"，不用拒绝切换）。

---

## §4 ④ 雷区与已知红

### 4.1 会立刻咬人的八条

1. **思考档到不了线上**（本票头号）：`prompt.go:136` 那枚字段非测试码**零赋值**；`agent.Config` 没这个字段（`loop.go:110-160`）。
   ⇒ 只做写入口＝交一个"改了没反应"的控件。判据必须**读到线上 JSON body**，不能只读结构体（照 `maxoutput_370_test.go:57-60` 的第三跳）。
2. **`rt.endpoint` 启动后不再变**：唯一赋值 `run.go:517`，读点 `:1463` 与 `panel_pump.go:308`。
   ⇒ 任何"模型已切换"的说法，不碰这三处就是谎。⚠ `run.go:508-511` 的注释写着"re-resolved per task in execute() … exactly like rt.endpoint"——
   **这句话与 `:1463` 的事实相反**（`ep := rt.endpoint` 是读字段，不是重解）；落地腿应把它改准，别照注释理解。
3. **`openai-chat` 协议根本不发思考参数**：`internal/llm/openaichat/adapter.go:14-17` 逐字"intentionally NOT mapped… sending provider-private params to strict endpoints causes 400s"。
   ⇒ 本机默认 provider（stepfun/deepseek 那族，`config_write_313.go` 的用例就是 stepfun）走的是 `openai-chat` ⇒ **即使接完线，思考档在这些 provider 上仍然是"设置成功、线上不变"**。
   这条必须在页面上说出来（AC#4 的"只列声明支持的档"要对齐到"**这个模型 + 这个协议**是否真能发"）：建议 payload 里给 `thinkingSupported`（协议侧），页面据此禁用控件并给一句原因。
   这是本票**最容易翻车又最容易被忽略**的一格。
4. **`enabled` 的隐藏语义**：手写目录条目缺 `enabled` 键 ⇒ 解成 `false`（`internal/config/enabled_261_test.go:5-24`）；且**任一次 settings 写入都会把 `enabled = false` 物化进文件**（同文件 `:14-20`）。
   ⇒ 选择器一筛 `Enabled`，用户会看到"我什么都没干，模型全不见了"。处置：**筛掉但别静默**——payload 另给 `disabledModels`，页面显示成"已停用（在 config.toml 里 enabled = false）"。
5. **写侧不校 `Enabled`**：`internal/config/settings.go:356-370 requireCatalogEntry` 只问"这一条在不在"，不问"开没开"；
   `ResolveRole` 在 `resolver.go:124-131` 才拒。⇒ 面板写进一枚 disabled 模型 = 写成功、下一个任务**整轮 Unconfigured**（退出码 2 那一族）。
   ⇒ 要么选择器不给选（§4 雷 4），要么 setter 加一条 `spec.Enabled` 检查（推荐两个都做）。
6. **"role 的档 ⊆ 模型声明的 levels"没有任何一处校验**：`validate.go:232-246` 只校枚举元素、`:251-274` 只校 role 的枚举值。
   ⇒ 面板要列"该模型声明支持的档"（AC#4），但文件里可以写一个模型不支持的档而**加载期不报错**。
   建议：setter 里加 `requireModelDeclaresLevel`（`ThinkingLevels` 空＝"未声明"，按"不拦"处理，否则现有手写目录全都写不进去）。
7. **`quota_44r5_test.go:175-209` 是源码字面量钉**：它要求 `run.go` 里逐字存在
   `prov, quotaGate, ledger, servedKey := quotaDispatchChain(cfg, rt.provs, quotaKeys, rt.quotaLedger, sink.Publish)`。
   ⇒ §3.2 案型乙会打红它。**改它是改既有断言，要具名**（改法：把那行换成读新变量，钉的 `want` 同步改成新行——不是删断言）。
8. **`git_test.go:497-513` 钉死 `panel.*` 恰好四枚**：本票不加方法 ⇒ **不红**；一旦有人"顺手"加 `panel.model.request`，这枚钉先红（它是 C17 变更的**自动审计**，别绕过）。
9. **票面与注释里的行号已经漂了三处**（本程现读，落地腿别照抄旧号）：
   - 票面 `现量` 第 2 行引 `internal/agent/prompt.go:113-114`＋`:225` 作为"思考档真被用"的凭据——今天那两处分别是 `InstructionSource` 接口尾部与别处，
     真正相关的两行是 **`:135-136`（字段声明）与 `:254`（塞进 `llm.Request`）**；而"塞进"那一步今天**恒为空**（§0 第 3 条）。
   - `cmd/wisp/panel_pump.go:296` 说"run.go**:308** assigns it"——`rt.endpoint` 的赋值今天是 **`:517`**。
   - 票面 `现量` 第 1 行引的 `internal/config/schema.go:291`／`:359` 与今天一致（`:296` 是 `thinking_intensity` 那行、`:364` 是 `ThinkingLevels`），差别只在 ±5 行，属正常漂移。

### 4.2 已知红（本程〔现跑〕，逐名照实记，不修、不当绿）

`go test -count=1 ./internal/panel/ -run 'TestComposerContractTypesMatchFrontend|TestC21DesignTokensFourWayAgree|TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme'` → **rc=1，两枚 FAIL**：

1. **`TestComposerContractTypesMatchFrontend`**（页侧契约）：红句逐字
   `Go ComposerState emits [files] that interface ComposerState does not declare`（`composer_test.go:78`）。
   成因＝票 190 的 `files` 键 Go 侧已发（`composer.go:285`）、页侧 `frontend/src/lib/panel.ts` 未声明（那是 UI 腿的写面）。
   同次读数：`Snapshot <-> PanelSnapshot: 6 JSON keys reconciled` ⇒ 页侧**已**声明 `git`／`currentModel`／`modelKnown`／`credentialState`／`credentialKnown`／`instructions`／`tasks`。
2. **`TestC21DesignTokensFourWayAgree`**（tokens 域）：六行
   `docs/evidence/s1/c21-native-tokens.md:58-60/:103-105 the table names Palette.GlassRing / GlassHi / GlassHiSoft, but internal/ball/tokens.go holds no such colour literal in DarkPalette()/LightPalette()`，
   汇总行 `78 colour rows … 72 agreed`。**与 187 无关**，⛔ 不许顺手改 `internal/ball/tokens.go` 或那份证据表。
3. **`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` 今天绿**（没出现在 FAIL 里）⇒ `188-r2`／`190` 面里"三枚在册红"的第三枚是**历史状态**，别照着去找。

⇒ 本票若照 §2.3 不动 `composer.go`，**这两枚红的颜色与红句都不会因 187 变化**；落地腿的读数是"仍是这两枚"，不是"多了/少了"。

### 4.3 会因 187 而需要**具名改钉**的三处（不是放宽，是钉在干活）

| 钉 | 位置 | 为什么会动 | 改法 |
|---|---|---|---|
| payload 顶层键集 | `internal/panel/config_313_test.go:290`（`want` 11 键字面量） | `chatThinking`／`chatThinkingKnown`／`disabledModels` 之类追加 | 把新键写进那个 `want` 字面量（键集是契约，增长就要写进契约） |
| provider 行键集 | 同上 `:317`（10 键字面量） | 每模型 `ThinkingLevels` 要挂在行上 | 同上；⚠ 若把 `models: []string` 改成对象数组，会**破**页侧 `config-view.ts:63` 的 `models: string[] | null` ⇒ **改成一个新键**（如 `modelThinking`），别改既有键的类型 |
| 名册并集 | `cmd/wisp/panel_config_313_test.go:225-267` | 它读**活名册**求和 ⇒ 往 `fields313` 追加**自动覆盖，不需要改** | 无需改（这正是"追加到 313 名册而不是 248 表"的好处） |
| 配额装配源码钉 | `cmd/wisp/quota_44r5_test.go:182` | 仅当选 §3.2 案型乙 | §4 雷 7 |

### 4.4 工程面

1. **共享编译面**：`internal/panel/**` 与 `cmd/wisp/**` 每波只许一枚写腿；`go test ./cmd/wisp` 重包同一时刻只跑一个（`AGENTS.md` §时代索引③）。
   本票要动**两边**（P4/P6 在 panel、P1/P2/P5/P7/P8 在 agent/config/cmd）⇒ 排期上必须独占一波。
2. **`go test ./cmd/wisp` 是重包**：里面有真 run 的测试（超时百秒级），落地腿先跑定向 `-run`，别整包乱跑。
3. **`internal/panel` 今天的两枚红**（§4.2）在落地前后都要照实记；`d22scan` 基线（票面写 `examined=433`）与 `gate-clauses` 名册（在册 `G6neg`）**本程未跑**，落地腿自量。
4. **回执文案会变**：`config_handlers.go:507` 那句"这一项要重启进程并重新运行才算用上（**同一次运行里模型通路不会重建**）"在 §3.2 乙之后对 `[llm]` 就是**假话**，必须改成 `EffectiveNextTask` 那句（`:504-505`）。这是**改文案**，不是改契约文字，但要在证据件里点名。

---

## §5 没测到／边界（照实写，不当已证）

1. **没跑 `cmd/wisp` 任何测试**（含 §3.2 案型乙会撞的 `quota_44r5_test.go`）⇒ 那枚红是**读断言后的推断**，真值要落地腿现跑。
2. **没跑 `internal/agent/`／`internal/config/` 全量、`d22scan`、`gate-clauses`**；只跑了 §4.2 那三枚定向测试。
3. **行数是读后估**：P1-P7 的"约行数"按同类既有实现（`SetRoleChatModel` 30 行、票 370 的 3 行接线、`fields313` 一行）推的，未写码验证。
4. **"每任务重建 provider"的真实代价只做了结构推理**：`BucketLimiter` 重建会让 RPM 节奏归零（`resolver.go:284`），没有压测；缓存形状建议按 `provider/model` 做 map，未验。
5. **协议面只核了三枚适配器的注释与请求编码**（openai-chat 不发、anthropic／openai-responses 发），**没有**对真实 provider 发请求验证"stepfun 到底吃不吃 reasoning 参数"——那是票 320 probe 那一族的地界。
6. **页侧零读数**：`frontend/**` 本程只读了两枚类型文件（`panel.ts`／`config-view.ts`）用于确认"为什么红"与"新键是否容错"，未跑任何前端工具链。
7. **`190`／`188` 的行号引用**取自它们的票面与普查件（`docs/evidence/s1/188-a1-census.md` §3.1 的 `run.go:1666-1680`），本程**未逐行复跑**。

---

## §6 纪律面

- 仓库内写点：**只有本文件**。`git add`／`git commit` 都用**显式 pathspec**（`-- docs/evidence/s1/187-a1-census.md`），未 `push`；
  未 `-A`／`--amend`／`reset`／`stash`／`checkout`；仓内零删除；未建 worktree。
- 临时件：`D:\tmp\wisp\20261104\187a1\readings.md`（只建不删）。
- `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／`internal/panel/bridge.go:41-76`／票面／`frontend/**`：**一字未动**。
- 本程跑过的命令：`git rev-parse`／`git status --porcelain`（只读）、`go version`／`go env GOFLAGS`、
  一发定向 `go test ./internal/panel/ -run …`（只读，未生成仓内件）、`New-Item` 建临时目录（仓外）。
  **未跑** `d22scan`／`gate-clauses`／`cmd/wisp` 任何包／前端工具链。
