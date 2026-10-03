# 票 145 · AC#2b 前置普查（腿 `145-a2`，只读）

- 起手：2026-10-03T10:33:46+08:00；锚点 `e16e303773f59be81d2f02eee3cc638870aa1630`（branch `dev`）。
- 本腿禁跑任何 Go 命令（3 枚写腿在飞：255-r2 / 171-r3 / 252-r2）：结论只来自读文件＋grep；跑才知道的格子逐枚登记在末节「量不到」。
- 票面 AC#2b 原文（本腿复量，非转述）：`.scratch/wisp/issues/145-panel-snapshot-has-four-fields-so-eleven-of-the-fourteen-ui-states-have-no-input-to-render-grow-the-go-side-carrier.md`
  - 追加节标题 :91；现量 1~5 :95-101；勾框与判据 :103-104。
  - 两维逐字：①当前可选模型清单（provider＋model id＋`Display`，来源＝配置目录 ∩ `enabled=true`；不许把未启用的也列进去）；②当前模型的思考档位集（`Capabilities.Thinking` 决定「有没有这一维」＋该适配器实际接受的档位词表决定「有哪几档」，两问分开答）。
  - 判据逐字：拿一份两模型／两档位的假配置树喂进去，快照里读到的清单与那棵树**逐字一致**，不许用「字段非空」充当判据。
  - 只做显示；任何写回＝`bridge.go:42-45` 四枚 `panel.*` 之外新增方法＝C17 契约变更，已并入 `Q-64`（默认不做）（:103 末段）。
- `frontend/**`／`design/**` 既不读也不引；TS 对齐那一格登记为射程外（见 D 末与「量不到」）。

---

## A 清单的源（provider ∩ model 目录逐枚指认）

### A1 配置侧（目录本体）

| 东西 | 位置 | 谁建／真值来源 |
|---|---|---|
| `Config.LLM`（`toml:"llm"`） | `internal/config/schema.go:117`（Config struct 枚定于 :107-134，**顶层无 `providers` 字段**） | `config.toml` 经 `config.LoadFile`／`config.NewManager` 解码 |
| `LLMSection.Providers map[string]Provider` | `internal/config/schema.go:419` | 同上 ⇒ 真路径＝`[llm.providers.<名>]` |
| `Provider.Models map[string]ModelSpec` | `internal/config/schema.go:400` | ⇒ `[llm.providers.<名>].models.<id>` |
| `ModelSpec`（注释逐字 `ModelSpec is llm.providers.<name>.models.<model-id>: one catalog entry` :345） | `internal/config/schema.go:348-372`：`Display` :350（注释：empty = model id）、`Capabilities` :352、`ThinkingLevels` :359、`Enabled` :371（tag `default:"true"`，注释 :369-370 「removes the model from discovery/selection」） | |
| `Capabilities`（含 `Thinking bool` :330） | `internal/config/schema.go:324-333`；头注 :320-323 逐字：声明可由人填、必须由探测核实（票 11），探测结果住 SQLite `provider_health`，这里刻意不留字段 | |

- `[providers.x]` 坑：解码器开 `DisallowUnknownFields`（`internal/config/parse.go:72`；库＝`github.com/pelletier/go-toml/v2 v2.2.4`，`go.mod:9`＋`parse.go:14` import），Config 顶层 struct 没有 `providers` 字段 ⇒ 静态判：未知键、加载即拒。**成立**（运行时错误文案量不到，见末节）。
- ⚠⚠ **`enabled` 默认方向（本腿复量＝推翻票面注释的宣称，读码实证）**：
  1. `applyDefaults` 只走 struct 字段、Map 一律留 nil（`internal/config/defaults.go:64-88`，Map 分支 :77-78 逐字 "leave nil (see NewDefaults)"）——**不会进入 `Providers`/`Models` 的 map 值**，`:371` 的 `default:"true"` tag 对 TOML 解码出来的条目**不生效**；
  2. `normalizeConfig`（`defaults.go:170-176`）只做「空 map／空切片→nil」，不填任何默认；
  3. ⇒ 手写 `config.toml` 里**没写 `enabled = true` 的目录条目解码后 `Enabled=false`**，与 `schema.go:346-347` 注释「New models default to ... enabled=true」**相反**（该注释只在 `ImportDiscovered` 路径兑现：`internal/llm/discover.go:116-120` 显式 `Enabled: true`、`Display: m.ID`、capabilities 全 false）。
- 目录校验只查枚举不查 enabled：`internal/config/validate.go:153-183`（provider 侧 ref/billing/protocol）＋ `validateModelSpec` :196-209（billing、`thinking_levels` 元素必须 ∈ off|low|medium|high，词表 `schema.go:88` `thinkingIntensity`）。

### A2 llm 侧（谁消费目录）

| 东西 | 位置 | 现状 |
|---|---|---|
| `Resolver` 直接持有 `Providers map[string]config.Provider` | `internal/llm/resolver.go:83`；构造 `NewResolver` :90-100 | 生产构造：`cmd/wisp/run.go:435`（`res := llm.NewResolver(cfg, st)`）→ `:436 ResolveRole(llm.RoleChat)` → `:441 rt.endpoint = ep` |
| `resolveEndpoint` 取一对 provider/model | `internal/llm/resolver.go:113-135` | **不检查 `spec.Enabled`**（:119 取到 spec、:124-135 组 Endpoint，全程无 enabled 判断） |
| 唯一「列目录 id」现成导出件 `DiscoveredModels` | `internal/llm/resolver.go:286-295`（sorted ids，单 provider） | **非测试生产调用者＝0**（grep `DiscoveredModels(` 过 cmd/internal/tools，除定义外零命中）；不带 Display、不过滤 Enabled |
| `Endpoint.Model` | `internal/llm/resolver.go:46`（struct :44-58） | 现量确认；已进快照（`composer.currentModel`，见 B） |
| ⚠ 全仓**读 `ModelSpec.Enabled` 的生产代码＝零行** | grep `\.Enabled` 非测试命中只有 `cmd/wisp/logsink.go:205,210`（slog 的同名字段） | 「∩ enabled=true」这一维今天**没有执法者也没有读取者**，票面 AC#2b 的来源式在产码里落不了地（见 C） |

### A3 探测核实位（档位「有没有」的第二问）

- 写侧（实测）：`internal/llm/probe_health.go`（头注 :14-22：RunProbeSuite 的 verdict 只来自真实请求，声明位只用来命名 mismatch）；生产入口＝CLI `wisp providers probe <p>/<m>`（`cmd/wisp/providers.go:152-220`，Declared 取 :193，落库 Sink＝`storeHealthSink` 定义于 `cmd/wisp/run.go:1168-1172`）。
- 读侧：`internal/memory/dao_providerhealth.go` —— `ProviderHealthRow` :167、`ListProviderHealth` :175；**非测试调用者＝0**（grep 全仓非测试只命中定义＋`internal/memory/schema.go:20` 注释）。⇒ 核实位「写得出、没人读回」。
- 发现式录入把 capabilities 全填 false（`internal/llm/discover.go:116-119`），且 `ImportDiscovered` 非测试调用者＝0（只有 `catalog_test.go`／`mockllm_integ_test.go`）。

### A4 现成的读面指认（票面点 5 的复量）

- `cmd/wisp/models.go`：`cmdModels` 在 **:102**（票面 :104 漂 2 行）；其 `modelsList`（:201-225 区）读 `store.manifest.Models`（:216），注释逐字 "not about this boot's [models] section"（:206-208 一带）⇒ **签验下载清单（C29 manifest），不是 llm 目录**。票面「复用它的取数」**复不出来**（r2/R-3① 的断言，本腿复量成立）。
- `cmd/wisp/providers.go`：`discover`＝打远端 `/v1/models`（:116-150，`llm.DiscoverModels` :136），`probe`＝实测（:152+）⇒ 都不是「本地目录 ∩ enabled」的列法。
- `cmd/wisp/panel_config_store.go:88-133`（`ReadSettings`）：已在遍历 `cfg.LLM.Providers`（:99）＋数 `p.Models`（:112-114，`row.ModelCount++`），**但只出计数**，且**无 enabled 过滤**；出口＝`config.get` 路由渲染的人读文本（`internal/panel/config_handlers.go:446-462` `renderSettingsView`），不是结构化快照。⇒ 它是「取数形状可参照」的腿，不是可直接复用的枚举 API。

---

## B 快照的路（构造点枚数／归属＋每一跳三档）

### B1 生产构造点（现量，非测试；票面 AC#2 点名的格）

`PanelSnapshot{` 字面量＝0 枚（Go 侧类型叫 `Snapshot`，编排者 :155 自记的坏尺复量仍成立）。真身 `Snapshot{` 非测试 6 枚，逐枚归属：

1. `internal/panel/composer.go:127` —— `NewSnapshot`（:106-133）的 return，**唯一整包构造器**；填 `Pending/Results/Composer/GeneratedAt` 4 枚，带 mode 空→unknown（:115-118）、MIME nil→全集（:119-121）、Git.Kind 空→NotProbed（:122-126）三道守卫；**不碰** Instructions/Tasks。
2. `internal/panel/pump.go:236` —— `Snapshot()` 的 nil-pump 零值支。
3. `internal/panel/pump.go:350` —— `Publish()`（:348-365）的 nil-pump 错误支。
4. `cmd/wisp/panel_pump.go:385` —— `lastPanelSnapshot()` 的 nil 零值支。
5. `internal/perm/store.go:247` —— 同名不同物（perm 的 Snapshot）。
6. `internal/ball/liquid.go:276` —— 同名不同物（motionSnapshot）。

**真·生产装配链唯一装配根（现锚行号）**：
`cmd/wisp/run.go:699-726`（`panel.NewSnapshotPump(panel.PumpSources{...})`，`NewSnapshotPump(` 非测试调用者＝1，复量与 r2 口径一致）→ `internal/panel/pump.go:234 Snapshot()` → `:303 NewSnapshot(...)` → `composer.go:127`。
⚠ r3 台件引的 `run.go:308/:442`、`:365` 已漂移到现锚 `:441`、`:544`、`:699`（改判行号，非改判事实）。

### B2 装配根现接的读口（逐枚点名）

`run.go:700-725` 字面量接 10 根线：`Verdicts:700 / Mode:701 / Workspace:702 / Git:703 / Model:704 / Credential:710 / Results:711 / Instructions:718 / Tasks:724 / Out:725`。
`PumpSources` 共 13 枚字段（`pump.go:138-224`），未接 3 枚：`L1Windows`（:200，nil＝诚实无枚举）、`AttachmentMax`（:205，零→常量）、`Now`（:209，nil→time.Now）。

### B3 谁产出／谁投递／谁落盘（判「接上了没有」的口径）

- 产出：`pump.go:234-331 Snapshot()`（composer 段逐读口填充 :256-298：Git :275-279、Model :280-286、Credential :287-298；顶层指针 Instructions :304-309、Tasks :310-329）。
- 投递：`pump.go:348 Publish()` → `json.Marshal` → `src.Out`；`Out` 生产实参＝`rt.bookPanelSnapshot`（`run.go:725`）。
- 落盘：`cmd/wisp/panel_pump.go:319-325` —— 进程内留 `rt.lastSnap/lastSnapBytes` ＋ ledger 记**有界摘要行**（:290-313 自述：512-rune 界装不下整包，最小样本 552B/518r）。
- 触发者：`publishPanelSnapshot`（`panel_pump.go:401-409`），生产调用＝`run.go:735`（`rt.ui.publish`）＋`run.go:998`（consoleSink）。
- ⚠ 最后一公里：`run.go:691-697` 注释说「this tree carries no WebView2 host」**过期**——`cmd/wisp/panel_host_windows.go` 存在（`//go:build windows` :1，C27/票 33 宿主），但该文件 grep `Snapshot` **零命中** ⇒ 准确形状＝「宿主存在、不消费快照」。包的真实落点＝进程内＋摘要账，页面侧今天本来就拿不到整包。

### B4 「清单到面板」逐跳三档（〔已有〕／〔缺一行赋值〕／〔缺一整块〕）

维度① 可选模型清单：
1. 目录数据在对象里：`rt.cfg`/`rt.mgr`（`run.go:424-425`）＝ `cfg.LLM.Providers` 全导出可遍历 ——〔已有〕。
2. 「provider＋id＋Display ∩ enabled」枚举件：不存在（`DiscoveredModels` 不带 Display/不过滤且零调用者，A2）——〔缺一整块〕。
3. 泵读口 `PumpSources.Models`＋`Snapshot()` 填块＋`ComposerState` 字段（数组＋known 伴生位）——〔缺一整块〕。
4. 装配根字面量多接一根 `Models:` ——〔缺一行赋值〕（函数体属 3 那块）。
5. 投递／落盘：Out 链现成 ——〔已有〕。
6. 页面侧渲染：本腿射程外（Q-51；宿主不消费快照，见 B3）。

维度② 思考档位（两问分开）：
1. 「有没有」——声明位 `Capabilities.Thinking` ——〔已有〕；实测核实位——〔缺一整块〕（DAO 读侧零生产调用者，A3）。
2. 「有哪几档」——声明枚举 `ModelSpec.ThinkingLevels` ——〔已有〕（枚举被 :201-206 校验；⚠ 未写键解码为 nil→normalizeZero 后仍是 nil，与「声明为空」塌成同形）；适配器接受词表——〔缺一整块〕：三张表全为**包内私有 var、零导出访问器**（`internal/llm/openairesponses/request.go:63 effortLevels`＝low/medium/high 三档，:128-135 消费；`internal/llm/anthropic/request.go:92 thinkingBudgets`＝三档，:113 消费；`internal/llm/openaichat/adapter.go:14-17` 逐字 "intentionally NOT mapped ... 严格端点会 400"＝一档不映射）。r2/R-3② 的断言**复量成立**。
3. 泵读口／填块／字段／装配：同①各〔缺一整块〕。

---

## C 宁缺毋造（哪几枚子字段真有源／哪几枚只能留空／留空会不会画成假话）

**今天真有源（可直接带值）**：
- provider 名（`Providers` map 键）、model id（`Models` map 键）、`Display`（空→用 id，`schema.go:349` 注释是 Go 侧现成契约句）、`ThinkingLevels` 声明枚举（校验过的词表元素）。
- `Capabilities.Thinking` **声明位**（但语义＝「人没改之前不算数」：探测前全 false 是录入默认，`discover.go:116-119`）。

**今天无源、只能留空／必须带状态位**：
1. 「∩ enabled=true」过滤语义本身：产码零读 `Enabled`（A2），且解码默认方向＝**缺键 false**（A1 ⚠⚠）。照票面直接做交集，对手写未带 `enabled` 键的配置会交出**空清单**，而实际请求链照跑不误（resolver 不看它）——这不是「留空」，是**把目录里的模型抹成不存在**，属假话形状。
2. 思考实测核实：provider_health 读侧零生产调用者（A3）⇒ 「探测过没有」今天无从可读；票面 R-3③ 逐字「不许把『未探测』画成『不支持』」。
3. 适配器接受档位词表：私有 var、无导出访问器（B4②-2）⇒ AC#2b 第二问的「有哪几档（实际侧）」**缺的是源，不是字段**（r2 原话，复量成立）。

**留空会不会显示成假话（Go 侧证据，不读前端）**：本包的反假话先例是成文的——空 mode 塌 unknown（`composer.go:115-118`）、Git 未读保 NotProbed（:122-126）、Credential 无读口报 Unknown 而非「未录入」（:287-298 注释逐字「nobody asked ≠ nothing recorded」）、模型无读口 `modelKnown=false`（`pump.go:280-286`）、票 200 教训「list＋omitempty 把四种真相塌成一枚缺键」（`composer.go:69-74`）。⇒ 新加 `models`/`efforts` 若**不带 known/状态伴生位**，「目录读不到」「目录真空」「全被 enabled 关了」三种真相将同形——正是本仓已命名过的「文案在、控件不在」族（ledger A408 引用处 `composer.go:73`）。落地面必须按先例每维配位，否则宁缺毋造只完成一半。

---

## D 改动面最小形状（只给面，不选形、不写码）

维度①（可完全落进已开三件＋run.go 一行＋用例家）：
| 文件 | 动什么 | 量级 |
|---|---|---|
| `internal/panel/composer.go` | ＋1 枚行型 view struct（Provider/Model/Display 三枚导出字段）；`ComposerState` ＋`models []…` ＋`modelsKnown bool` | ~15-25 行（含守卫注释） |
| `internal/panel/pump.go` | `PumpSources` ＋`Models func() …` 读口；`Snapshot()` 在 :286 后加填块（reader 存在才信） | ~8-15 行 |
| `cmd/wisp/panel_pump.go` | ＋1 枚 reader 函数：遍历 `rt.cfg`/`rt.settings.mgr.Config()` 的 `LLM.Providers`→`Models`（字段全导出，**无需在 internal/llm、internal/config 新开任何导出 API**）；需先裁 enabled 语义（C-1） | ~15-25 行 |
| `cmd/wisp/run.go` | :699-726 字面量 ＋1 行 `Models: rt.…,` | 1 行（⚠ run.go 不在三件具名内，需编排者再具名放行；先例＝r3 :165） |
| 用例之家 | `internal/panel/composer_test.go` ＋1 枚 `TestXxx`（假树两模型两档位逐字对，票面 :103 判据形；先例＝r3 :163/:167，其越界裁量未销账） | 测试件，非产码 |

维度②：
- 若只落**声明面**（ThinkingLevels 枚举＋Thinking 声明位＋known 伴生）：同上三件＋run.go 一行，无新增导出 API。
- 若按票面第二问落「**适配器实际接受**」交集：必须动 `internal/llm/openairesponses/request.go`、`internal/llm/anthropic/request.go`、`internal/llm/openaichat/adapter.go` 各加 1 枚导出访问器（每枚 ~3-6 行）——**这三枚文件不在票 145 任何已具名写面内＝新开面，需编排者先落批准记录**（本腿不代裁）。
- 实测核实位（provider_health 读回）要另开 `cmd/wisp` 读腿＋`internal/memory` 调用，整块新面，本格不建议。

`panel.*`／`bridge.go`：**纯显示不需要任何新方法名**，`bridge.go:42-45` 四枚常量零触碰 ⇒ `want-4` 锚（git_test.go，见 E#1）不动。只有「面板改档位／换模型」的**写回**才需要第五枚门＝C17 契约变更，票面已并入 `Q-64` 默认不做（:103）；碰它之前要编排者落批准记录。
TS 页面侧对齐那一跳（`PanelSnapshot`/`ComposerState` interface 增键）：需前端会话自取，**本腿射程外**（两层禁令；双向钉 `composer_test.go:49` 要求两侧同 commit 挪动，属 `Q-51` 族）。

---

## E 邻居尺名册（负向尺逐枚＋隔离一句＋冻结范围核查）

词面型／计数型负向尺（文件:行＋数什么；本锚复量）：

1. `internal/panel/git_test.go:362-390` `TestGitDimensionHasNoModelCallableTool` 三门：D34 表 git 行、`internal/tools` 的 `git.*` 名、**bridge.go 引号内 `panel.*` 字面量恰四枚**（`wantMethods` :381-383，比对 :385；正则 `panelMethodRe` :394；提取器 `whitelistMethodsFromSource` :407-423；正控 :504-519，其中 :517 现量真 bridge＝4）。隔离：不往 bridge.go 加任何 `"panel.*"` 字面量；显示维字段名不带 `panel.` 前缀形状。
2. `internal/panel/composer_test.go:405-417` `composerRouteLiterals` 闭集（4 枚 request＋`panel.approval.request`）＋ `:502` `TestTheRendererHoldsExactlyOneDoorToTheHost`（扫 frontend/src：postMessage 恰两枚、无拼装路由、渲染器字面量必须被 Go 答）。隔离：只加 Go 侧输出字段＝不响；加任何 inbound 名才响。
3. `internal/panel/composer_test.go:49` `TestComposerContractTypesMatchFrontend`：Snapshot/ComposerState/ModeView/WorkspaceView/AttachmentRef/ResultChunk 六对↔`panel.ts` 同名 interface **双向 JSON 键差集**。⚠ 加 composer 嵌套键**必响**（在册三红之一，「更红＝预期」先例见票面 r3 :164 口径）；不是隔离对象，是**必须随行的页侧声明**（Q-51，前端会话）。
4. `internal/panel/approval_test.go:105` `TestApprovalCardViewJSONKeysMatchFrontendTypes`：含 `Snapshot`↔`PanelSnapshot` 顶层键差集（pairs :114-118）。隔离：新键嵌进 composer 段、顶层键集不动 ⇒ 保持绿。
5. `internal/panel/pump_test.go:123-125` 与 `:291-293`：出口字节顶层键**恰四** `composer,generatedAt,pending,results`。⚠ 票面记的 `:124/:276` 已漂（现锚＝:123/:291，r2 台件行号过期）。隔离：不加顶层第五枚键 ⇒ 不响、也不触发那两行的预解冻条件。
6. `internal/panel/subagent_roster_197_test.go:214`：**第三枚同尺**（同四字面串）。票面「两枚字节钉」的说法**少计一枚**——顶层加键要同时挪三处，别只数两枚。隔离同上：嵌套。
7. `internal/panel/inbound_roster_253_test.go:419` `TestFullInboundMethodRosterIsClosed`（＋正控 :597/:696/:769/:809）：按结构闭枚全部 inbound 方法名（头注 :6-11 自述：git/composer 两尺硬编码 `panel.` 前缀，看不见 `config.get/config.set`＝bridge.go:66-67，本尺补盲）。隔离：不新增 inbound 方法名／`Method*` 常量 ⇒ 不响。
8. `internal/panel/l2_grant_boundary_test.go`：`:1229`（answered roster 无 approval-shaped 名）、`:1547`（无 inbound 封套能 bind verdict）、`:1595`（grant 词表对真封套不满足）、`:1959`、`:2168`（种桩：给**入站封套**加 `outcome` 字段即红——「InASnapshot」的 snapshot＝TempDir 源码副本，不是面板快照）、`:1807`（AST/reflect/encoding-json 三方同规则）。词表 :185-193＝allow*/permit/verdict/outcome/bypass/override/approve/approval/grant。隔离：显示侧新字段名**避开这些词根**（`models/efforts/modelsKnown/effortsKnown` 不碰）；不新增入站类型。
9. `internal/panel/frontend_hygiene_test.go`（:169/:196/:222/:246/:281）：只 walk `frontend/src`（:76-92）。Go 侧加字段零触碰；:196 在册红。
10. `internal/panel/tokens_fourway_test.go:439`：四方颜色对账（design/frontend/ball/tokens.go/c21 md），不扫 composer 结构；在册红。
11. `sh scripts/d22scan.sh`（tools/d22scan）：bans #1-8 扫 `design/ internal/ cmd/`（豁免注释、不豁免字符串；`frontend_hygiene_test.go:26-31` 自述 frontend/ 不在其 declared scope）。新读口不起 goroutine、不做路径决策（只读 map 键）⇒ 不触 #1/#2；字段**永不携带** `APIKeyRef` 值或任何密钥形态（#3）；文案字形避开仪器带谱（AGENTS.md §1.2：对勾带/数学带在仪器射程、字符串不豁免；箭头带是刻意留的空隙但建议仍用 ASCII）。

**改动面 vs 仍冻结的尺（只答会不会碰到，不代裁解冻）**：D 的最小形状＝`composer.go`＋`pump.go`＋`panel_pump.go`（票面可考的放开先例：r1 路径并集含 `internal/panel/{composer,pump}.go` :69；r3 实写含 panel_pump/currentModel :162）＋`run.go` 一行（需再具名）＋用例家 `composer_test.go`（r3 先例，越界裁量未销 :167）。**不碰** `tokens_fourway`／`l2_grant`／`frontend_hygiene`／`inbound_roster_253`／`pump_test.go:123/:291`／`approval_test.go:105`／`bridge.go:42-45`。唯一「会红」的在册尺＝#3（本属随行项，非误响）。若维度②选适配器交集，新开 internal/llm 三枚文件＝**写面扩张，不是尺冻结问题**。
并发冲突提示（归属点名）：`255-r2` 正写 `cmd/wisp`＋`internal/config`（`config_readers_255.go` 一族）——与本腿指认的读面同包不同文件；`252-r2` 在 `internal/tools`，D 面不碰。落写时 run.go/panel_pump.go 的行号必再复量。

---

## F 推翻清单（逐条断言→现量）

1. 「AC#2b 09-28 追加、至今未做」：**成立**（本锚现量：ComposerState 11 枚无 models/efforts；r3 落的 `currentModel/modelKnown` 不是该两维；票面四框未勾）。
2. 「两维＝可选模型清单＋思考档位」：**成立**，原文 :103；转述漏了判据形（假树逐字一致、禁「字段非空」）与「两问分开答」。
3. 「来源＝配置目录 ∩ enabled=true」：**纸面成立、执法不成立**——产码读 `ModelSpec.Enabled`＝0 行；缺键解码默认 **false**（与 `schema.go:346-347/371` 的注释/tag 宣称相反，机制在 `defaults.go:77-78`＋`normalizeConfig` 不填默认）。具名推翻对象＝schema.go 注释宣称；AC#2b 落地时该式照写会产出假空清单。
4. 「真路径 `[llm.providers.<名>]`、`[providers.x]` 撞 DisallowUnknownFields」：**静态成立**（Config struct 无顶层 providers；parse.go:72）。运行时错误文案＝量不到。
5. 「`want-4` 锚在 `internal/panel/git_test.go` 一带」：**成立**（:362-390/:407/:514-519）。⚠ 精确度更正：数的是 bridge.go **引号 `panel.*` 字面量**；同文件另有 `config.get/config.set`（:66-67，无前缀、不入该四枚、但被 inbound_roster_253 盯）——把 want-4 读成「bridge.go 只有四个方法」会漏两枚。
6. 「Go 侧 `Snapshot` 字段太少」：**过期**——现量 Snapshot 直接字段 6（`composer.go:57-92`）、ComposerState 11（:235-270）；票面 :176（c2）已更正，本腿复量一致。
7. 「票 145 局部解冻只放开 composer.go/pump.go/panel_pump.go、只到新增字段」：**转述非原文**。票面可考记录三条并不同形：r1＝`internal/panel/{composer,pump}.go`（:69），r3 编排者＝run.go 具名两行、射程「只多接一根读口」（:155），r3 实写＝三件＋run.go 两行＋composer_test.go 一枚（:162-167）。「三枚冻结测试件」在禁面清单（:149/:170）**未具名**——把它们对应到 tokens_fourway/l2_grant/frontend_hygiene 是转述猜的，文件存在、内容如上，但票面没写这三枚文件名。
8. 「复用 `cmdModels` 取数」（票面现量点 5）：**推翻复量成立**——modelsList 读签验 manifest（`cmd/wisp/models.go:200-225`），非 llm 目录；最接近的现成枚举＝`resolver.DiscoveredModels`（:286-295），缺 Display/缺过滤/零生产调用者。
9. 「档位词表＝私有 var 无导出」：**复量成立**（E/B4 三处行号实钉）。
10. r3 台件行号（run.go :308/:365/:442、composer.go :44-49/:79 族）：**漂移**——现锚 :441/:544/:699-726；r2 记的 pump_test 预解冻线现锚＝:123/:291；票面现量点 5 的 cmdModels :104 现＝:102。
11. 「run.go:691-697 本树无 WebView2 host」：**措辞过期**——panel_host_windows.go 存在但不消费快照（grep `Snapshot` 零命中）；「快照到不了页面」结论不变、理由换了。
12. 「票面十四行表没要求过模型清单」（:104 的尺）：本腿不复量 PLAN.md 表体（docs 读得到但那是票面已钉的现量，转述与原文无冲突）。

---

## 量不到的格子（全因本腿禁跑 Go；逐枚点名）

1. 在册三红（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）在 `e16e3037` 的**当前颜色**——三枚函数在盘（composer_test.go:48／frontend_hygiene_test.go:196／tokens_fourway_test.go:439）、名与册一致，跑色量不到。
2. `pump_test.go` 四键钉与 `subagent_roster_197_test.go:214` 当前是否绿——同上，仅静态确认断言行在。
3. `[providers.x]` 被拒的**实际错误文案／退出码**（parse.go:72 静态指向，运行时形态未验）。
4. 「缺 `enabled` 键→解码 false」的**行为实证**：仓内未找到同款用例（grep internal/config/internal/llm 测试无 omitted-enabled 断言；catalog_test/fallback_test 都显式写 `Enabled:true`）；本腿判读来自 applyDefaults 的 Map 盲区＋normalizeConfig 只做空→nil 两枚独立静态证据。写腿动刀前建议跑一枚三行解码探针坐实。
5. 新字段落地后的**出口字节数／ledger 512-rune 界**（panel_pump.go:290-313 的自测量不覆盖新维）。
6. TS 侧 `panel.ts` 当前键集与两维应增键清单的差量——两层禁令，需前端会话自取。
7. `wisp providers probe` 落库后 provider_health 行的真实形态（DAO 读回零调用者，写侧只有 CLI 触发路径；本腿不能跑）。
