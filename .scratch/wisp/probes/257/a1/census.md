# 257-a1 普查件 — AC#0 三形代价门（只读，未选形）

> 腿：`257-a1`；票：257。
> 锚点：进场 `16:37:55` `d3709381` → 量程中 `17:01:55` `ed769600` → `17:19:22` `41fa10a7` → 落件 `17:28:07` `41fa10a7`（共享树上别腿在落 commit；本腿零产码写入、零 Go 命令）。
> ⛔ 本件只量代价，不裁形。所有行号现读于上列锚点之间的树；禁 Go 命令的格子已在 §5 具名。

## 0. 票面四环现量复核

| # | 票面断言 | 复核 | 现量 |
|---|---|---|---|
| 1 | `defaults.go:77-78` 对 map「leave nil」；`schema.go:419` 无 default 标签 | **成立（行号全对）** | `defaults.go:77-78` 逐字：`case reflect.Map:` / `// leave nil (see NewDefaults)`；注释 `defaults.go:56-57` 还给了理由（go-toml merges into pre-existing maps ⇒ 预填会泄漏幻影条目）。`schema.go:419`＝`Providers map[string]Provider \`toml:"providers"\``，确无 `default` 标签 |
| 2 | 写侧要求行已存在：`settings.go:310-314`／`:294-308` 是**拒**不是建行 | **成立（行号对）** | `:310-314`＝`unknownProviderErr`（"a settings write addresses an existing entry, it does not create one"）；`:294-308`＝`requireCatalogEntry`（"refusing to invent one from a settings write"）。两把都被 `settings.go:60-76`/`:81-97` 的 apply+check 双路引用；模型级三枚（`:103-126`/`:132-139`/`:141-171`）走 `requireCatalogEntry` |
| 3 | `role_chat_model` 撞 `catalog.go:97-101`：model 有值 provider 空 ⇒ 拒 | **成立（行号对）** | `validateRoles`（`catalog.go:89-119`）的 switch：`role.Provider == ""` 分支在 `:101-103` 逐字「model set without provider (set both or neither)」。⚠ 精确形状：`SetRoleChatModel`（`settings.go:177-183`）**只写 model、不碰 provider**，而干净机首份配置里 `llm.roles.chat.provider` 与 `model` 都是空串 ⇒ 任何非空 model 都落进"model 有、provider 无"＝`:101-103` 拒。**不是** `:97-101` 整块——那是 case 头+注释行，真正的拒句在 `:102-103`（差 1 行，判读不变） |
| 4 | 名册七枚里没有能创建 provider 条目的字段 | **成立** | `config_handlers.go:57-69` 七枚常量 + `:90-93` `allWritableFields`：`provider_base_url`/`provider_api_key_ref`/`model_context_window`/`model_price_in`/`model_price_out`/`role_chat_model`/`provider_credential`。全部是"行已存在"前提下的叶子写入；无一枚能造 `llm.providers.<名>` 行 |
| 5 | 净结果＝干净机器七枚全写不进去 | **静态链成立**（真跑被禁，§5 N1） | 逐枚过门：provider/model 级六枚在 `check`（unknownProviderErr/requireCatalogEntry）就拒；`role_chat_model` 在写后 `validate()`（`settings.go:230`）撞 `catalog.go:101-103` 拒且文件回滚。七枚各有自己那道门，全部关着 |

补充现量（票面没写但对选形要紧）：
- **零provider 不是"空 map"是 nil**：`normalizeZero`（`defaults.go:148-151`）还把"文件里写了空 `[llm.providers]`"归零回 nil ⇒ 即使手加一节空表头，加载后读数与"没有"同形（round-trip 钉住，`defaults.go:135-138` 注释）。
- **写侧双层都关**：`writeOneKey` 的 `check` 拒（provider/model 六枚）之外，apply 也返回 false（`:62-69` `p, ok := base.LLM.Providers[provider]; if !ok { return false }`）⇒ 即使 check 被绕过，apply=false 走 `:224-227` 的"does not exist in the file"拒绝。要"建行"两处都得动，不是一个 if。
- **17:01 那枚钉先记下**：`settings_248_test.go:130-133` `TestAC7InvalidValueLeavesTheFileByteIdentical` 有一案逐字 `{"a provider the config does not declare", ... SetProviderBaseURL("inventco", ...)}` 断言**必须拒且文件逐字节不变**——ⓐ/ⓑ 改"行不存在则建"必然打红它，见 §1/§2 的钉清单。

## 1. ⓐ 空壳行（首建默认建一条可用空壳 provider 行，写侧就能改它）

**动的面（文件:行）**：
1. `internal/config/defaults.go:77-78`——`reflect.Map` 分支要么整体放行（给全部 map 发默认＝改变 NewDefaults 全局语义，撞 `defaults.go:56-57` 注释自述的幻影条目理由），要么加一条"仅 Providers"特判（`if fv.Type() == reflect.TypeOf(map[string]Provider(nil))`）。特判形＝默认表里第一次出现"非 tag 来源的值"，直接顶红下面两枚钉。
2. `cmd/wisp/firstrun_198r2_test.go:221` `TestTicket198R2AC2CreatedFileCarriesNothingButTheSchemaDefaultTags`——首建文件必须逐字节等于"`default` 标签渲染出来的文件"（`render198r2`），且 8 类 kind 正控种值必红（`:245-258`）。塞空壳行后 `render198r2(t, expect)` 与盘上文件不等 ⇒ **必红**。
3. `cmd/wisp/firstrun_198r2_test.go:274` `TestTicket198R2AC2ExportedDefaultTableStillRendersAsItsTags`——`NewDefaults()` 对 tag 渲染同理必红。
4. `cmd/wisp/firstrun_198_test.go:192` `TestTicket198AC1CreatedFileHoldsEveryStaticSectionAndNoInventedTables`——`:220-224` 逐字禁 `[llm.providers` 出现在首建文件（"an invented entry, not a default"）⇒ **必红**。
5. `internal/config/settings_248_test.go:130-133`——若"空壳行"同时松写侧"行不存在则建"，此案必红；若只建空壳行不松写侧，此案不红（写侧行为没变）。
6. schema 校验谁兜（票面问的那句）：空壳行若走**预设名**（如 `deepseek`），`applyPresets`（`loader.go:188-204`）加载时补 protocol/base_url ⇒ 校验全绿，这条是现成的免校验通道；若走**非预设名**（`inventco`），`validate.go:171-177` 逐字要求"非预设 provider 必须显式给 protocol"⇒ **首建文件自己就加载不过**（`NewManager`→`LoadFile`→`validate` 在 `manager.go:102`，文件写得出、进程起不来）。所以"可用的空壳行"实际只有两个子形：ⓐ-甲 塞预设名行（protocol/base_url 由 preset 继承，机主只欠 key 和 model）；ⓐ-乙 塞非预设名+显式 protocol（defaults 表要带 protocol 值＝更像"发明默认"）。
7. `role_chat_model` 的连锁：空壳行**不会**自动解锁它——`catalog.go:112` `modelExists` 仍要 provider 行下有 `models.<id>` 条目；role 指向该 provider 但 model 空则走 `:104-106` 拒。ⓐ 单独成立时，七枚里解锁的是 `provider_base_url`/`provider_api_key_ref`（+`provider_credential` 的 StoreCredential 链）＝**3/7**；`model_*` 三枚仍撞 `requireCatalogEntry`，`role_chat_model` 仍撞 catalog——除非再给空壳行预置一条空 model 行（那会同时顶红 `firstrun_198_test.go:220` 的 `[llm.providers` 禁令并叠加 §1-4 的红）。

**owner 视角**：设置页"配了没配"读数（`panel_config_store.go:96-131`）出现一行真 provider（如 deepseek）＝"服务商 1 家；凭据状态 no_ref_declared"；base_url/api_key_ref 两枚可写并回执"已写入，重启生效"；但目录模型 0 个、聊天模型写不进（拒句"no model "" in its catalog"）——机主要看见"配了一半"的诚实中间态。

## 2. ⓑ 名册加 provider 字段（＝C17 面，⛔ 具名上报不决定）

**C17 白名单现量**：`internal/panel/bridge.go:41-68` 共 **6 枚**方法名——`panel.mode.request`/`panel.workspace.request`/`panel.attachment.add`/`panel.message.send`（4 枚 `panel.*` 前缀）＋ `config.get`/`config.set`（刻意无前缀，注释 `:50-66` 给了三条理由，含与冻结审批仪器 `panel*` 命名空间的刻意隔离）。**新增方法名的步数**（若走"新方法"形；走"复用 config.set+新字段"形则只动字段表，不走这六步）：
1. `bridge.go:41-68` 加常量 → 2. `bridge.go:146-152` `knownComposerMethod` case → 3. `composer_test.go:409-419` `composerRouteLiterals()` 加一枚 → 4. `composer_dispatch.go:197-202` 一带加路由 → 5. 前端 `frontend/src/lib/panel.ts` 发同名信封（⛔ 本腿禁入 frontend，未量前端改动面）→ 6. `git_test.go:381-389` 的四枚 `panel.*` 断言：`whitelistMethodsFromSource` 只提 `"panel\.[a-z0-9_.-]+"` 引号字面（`:394`）——**`config.*` 形不加前缀就不进这把尺**（台账 1002 那笔"config.get/set 掉在两把名册尺外"同源）；但若新方法带 `panel.` 前缀，`TestGitDimensionHasNoModelCallableTool` 的 `len(real) != 4`（`git_test.go:517`）**必红**。
- 另一枚枚数门：`l2_grant_boundary_test.go:1293-1301`（answered 集 vs `Method*` 常量 1:1）——只要求"answered 的必有常量"，加常量本身不顶红；grant 词黑名单（`:1259-1266`）要求新名不带 approve/grant/allow/decide 族词。
- `config_handlers.go:259-263` `TestAC1FieldTableCarriesNoApprovalVocabulary` 的同款词禁也罩着字段名与路由名。

**两界对账门现状（今天已经红的枚数与键）**：
- `TestComposerContractTypesMatchFrontend`（`composer_test.go:48`）对 6 对结构体（`composer_test.go:55-66`：Snapshot/ComposerState/ModeView/WorkspaceView/AttachmentRef/ResultChunk）双向逐键对 `frontend/src/lib/panel.ts`。
- `TestApprovalCardViewJSONKeysMatchFrontendTypes`（`approval_test.go:105`）对 3 对（ApprovalCardView/ResultChunk/Snapshot）。
- **台账现量**（evidence `145-snapshot-growth-r2.md` §表：`internal/panel` 红三枚＝这两枚里只红 `TestComposerContractTypesMatchFrontend`＋`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`＋`TestC21DesignTokensFourWayAgree`；`145-snapshot-growth-r3.md:112` 载 composer 对账门"因本程两枚键更红"）。**静态推其今天红在 `ComposerState` 键组**：Go 侧 `composer.go:235-270` 逐字键＝`mode/workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError/git/currentModel/modelKnown/credentialState/credentialKnown`（11 枚）；TS 侧 `panel.ts:112-121` `ComposerState`＝`mode/workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError`（6 枚）⇒ **Go 多 5 键：`git`、`currentModel`、`modelKnown`、`credentialState`、`credentialKnown`**（每键双向都会报：Go 发→TS 未声明，方向 `subtract(goKeys, tsKeys)`，`composer_test.go:73-75`）。同样 `Snapshot`：Go 有 `instructions`/`tasks` 两枚 omitempty（`composer.go:74`/`:91`），TS `panel.ts:129-139` `PanelSnapshot` 只有 4 键 ⇒ 再多 2 键（两注释自陈"reconciliation now name them"）。**即：今天这枚红句点名 ≈7 键（ComposerState 5 + Snapshot 2），两处 `t.Errorf`（missing+extra）都可能响。**⚠ 此为静态推算（禁 Go ⇒ 量不到活跑句），§5 N2。
- **ⓑ 关键判读**：C17 面的"字段"若落在 `ComposerRequest`（信封，`bridge.go:117-120`）或 `SettingsView`（读侧，`config_handlers.go:142-156`）——这两枚**不在**上述任何对账对里（对账只罩 6+3 对视图模型）⇒ 新增枚**不顶红两界对账门**；但会顶 `cmd/wisp/panel_config_248_test.go:369-382` `TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope`——它断言信封上新键**只有** configField/configProvider/configModel/configValue 四枚 ⇒ 加第 5 键**必红**（此钉改不改＝契约决定，具名上报）。

**名册登记点**（走"config.set+新字段"形的最短路径）：
1. `internal/panel/config_handlers.go:57-69` 加 `FieldProviderCreate = "provider_create"`（或类似）常量；`:74-94` 三张表加行（`fieldNeedsProvider` 不加——它自己就是 provider 参数；`fieldTakesValue` 看语义）；`:90-93` `allWritableFields` 加枚。
2. `config_route_248_test.go:269-279` `want` 切片**逐字列 7 枚并 `reflect.DeepEqual`** ⇒ 加第 8 枚**必红**，改它＝动测试断言（许可范围＝写面腿 257 本票，但按"不许为变绿放宽断言"纪律要写清"枚数 7→8 是契约变更"）。
3. `cmd/wisp/panel_config_store.go:182-222` `ApplySetting` switch 加 case（新 setter 的接线）＋ `config` 包新 setter（`settings.go` 家族，"行不存在则建"形——撞 §1-5 那枚 `settings_248_test.go:130` 钉与 `unknownProviderErr` 语义反转）。
4. 新 setter 的"建行"还要过 `writeOneKey` 的 validate 门（`settings.go:230`）：非预设名必须带 protocol 一起落（一枚键一写的 AC#11 第二条注释 `settings.go:25-28` 明说没有多键原子写）⇒ **"建 provider"天生是多键写**（名+protocol 至少两键），要么发明多键原子写（撞 `settings.go:27-28`"inventing one is ticket 226's territory and needs its own ruling"），要么一条请求连发两枚守卫写（中途失败会留下半行＝新失败形状）。

**composerRouteLiterals() 要不要动**：`internal/panel/composer_test.go:409`（读作 `composer_test.go:409-419`）——**只在"新增 panel.* 方法"形下要动**（`routeLiteralRe`＝`"panel\.[A-Za-z.]+"`，`composer_test.go:394`，扫 `frontend/src` 全树：页面里出现新 `"panel.*"` 字面而 Go 不答 ⇒ `rep.unknown` 红，`composer_test.go:521-523`）；`config.*` 无前缀字面**不在**该 regex 射程 ⇒ 走字段表形不用动它。⚠ 台账 1002 已把"`config.get/set` 不带 `panel.` 前缀掉在两把名册尺外"记为**已知洞**——选形时这是现成邻证。
**owner 视角**：设置页多一枚"添加服务商"控件；成功回执要么"已创建 llm.providers.<名>，需重启"要么列出两枚落盘键路径；失败形状变多（protocol 欠/半行/名字撞预设）——AC#2"三种拒因各配一句"的第三种"校验不过"会在这形里最频繁出现。

## 3. ⓒ 指引手加 `[llm.providers.x]`（读 firstrun.go:92-115，只读不改）

**先纠一处票面引用形**：票面 AC#0-ⓒ 写"`[providers.x]`"——schema 真路径是 **`[llm.providers.<名>]`**（`Config.LLM`＝`toml:"llm"`（schema.go:116），`LLMSection.Providers`＝`toml:"providers"`（schema.go:419）；顶格写 `[providers.x]` 是**另一个键路径**）。后果链：`decodeStrict`（parse.go:62-88）`DisallowUnknownFields` ⇒ 顶层 `[providers]` 是未知键 ⇒ `formatDecodeError`（parse.go:139-149）报 `unknown key "providers"` ⇒ **LoadFile 直接失败**、`NewManager` 起不来 ⇒ 不止不解锁七枚，**面板整条装配死**（`panel_inbound.go:230-233` "配置未就绪"返错）。现成邻证：`cmd/wisp/firstrun.go:116-118` 的回执第二句写的正是 `[llm.providers.<名>]` 全路径；`firstrun_198r2_test.go:381` 也把 `[llm.providers.` 列为回执必名串。⇒ 若选 ⓒ，指引文案必须写全路径（或修票面那半句），一字之差是"面板死"与"面板活"的差。

**静态链：手加一段后重启认不认**（认，前提是写对三样）：
1. 路径对（`[llm.providers.deepseek]`）＋ `schema_version` 不动 ⇒ `readConfigFile`（loader.go:64-127）管线全过：`peekSchemaVersion`（ver=2 不走迁移）→ `decodeStrict`（`[llm.providers.x]` 是 schema 内合法动态表）→ `normalizeConfig`（非空 map 不归零）→ `applyPresets`（**预设名**自动补 protocol/base_url，loader.go:188-204；**非预设名**必须手写 `protocol`，validate.go:171-177 否则加载拒）→ `validate`。
2. `NewManager`（manager.go:101-111）全量重读 ⇒ **七枚解锁数**：provider 级 `provider_base_url`/`provider_api_key_ref`/`provider_credential` **3 枚**解锁（行已存在）；`model_*` 三枚仍锁（要 `[llm.providers.x.models.<id>]` 子表——**手加时一并写上 model 子表则再解锁 3 枚＝6/7**）；`role_chat_model` 仍锁——`validateCatalog` 的 `modelExists` 只查 provider 行下有没有那个 model 条目（catalog.go:145-148），**且 `role_chat_model` 写侧只写 model 不写 provider**（settings.go:177-183）⇒ 还要求手加时把 `[llm.roles.chat]` 的 `provider` 也点名，才能过 `validateRoles`（catalog.go:108-115）⇒ **七枚全解锁的手加形状＝provider 行＋至少一条 model 子表＋roles.chat.provider 点名**。这正是 firstrun 回执第二句教的顺序（`[llm.providers.<名>]` 补 api_key_ref/base_url/models.<id>，再在 text_chain 或 roles.chat 点名，firstrun.go:116-118）。
3. **回执是现成邻居**：`cmd/wisp/firstrun.go:92-118` 三段回执（建件句 :92-95＋key 句 :111-115＋模型句 :116-118）已在教手改 config.toml；ⓒ 的增量实际只是"把模型句往前挪/加重"（比如回执里明说"在设置页能写的只有这几样，其余要手改文件"）。**动的文件＝`cmd/wisp/firstrun.go` 一处（+可能文案钉）**：
   - `firstrun_198r2_test.go:363-397` `TestTicket198R2AC4ReceiptNamesTheRealEntryPoints`：回执必名 7 串（`wisp secret set`/`api_key_ref`/`dpapi:`/`env:`/`[llm.providers.`/`text_chain`/`roles.chat`）——ⓒ 改文案不能丢这 7 串，且命令名扫描器（`:386-395`）要求回执里每个 `wisp <词>` 都真有 `cmd<词>` 函数。
   - `firstrun_198_test.go:237` `TestTicket198FirstRunReceiptStaysOutOfTheFourCauseSentences`：回执不得混入 cause=missing 那句的标记（"继续用内存里的旧配置"）。
   - 第二次运行不重复回执的钉（`firstrun_198r2_test.go:399-405`）。
4. **ⓒ 的诚实边界**：firstrun.go:96-110 注释自陈"模型目录没有任何 CLI 写入者……模型这一格今天只有改文件一条路，文案照实说"——ⓒ 与这句现量完全一致，**不需要**新写入者；代价集中在"教得够不够清楚"，不在权限面。
5. ⓒ 解锁不了的那一格：机主照指引改文件时若进程在跑，`wisp run` 有 1s 热加载（config_reload.go:101-153→run.go:813），**改完不用重启**（写错则 HOT-RELOAD not-applied，内存配置保持）；但设置页写入腿的回执仍说"重启生效"（panel_config_store.go:26-33：模型通路装配期建一次、reload 不重建）——ⓒ 教用户手改+热加载与设置页"要重启"两句**并存且都对**，回执措辞要分清"手改文件（热加载认）/面板写（要重启）"两条通道，否则 AC#2"不许折成一句"会在这形里复发。

**owner 视角**：首建回执多一段"去哪儿配"；设置页在干净机上仍七枚灰/拒，但每枚拒句应指向"文件里加 `[llm.providers.<名>]`"（拒句今天已如此——`unknownProviderErr` 的 "addresses an existing entry, it does not create one"）⇒ owner 看见的诚实台词＝"这页改不了服务商的存在性，去文件里加一节"。**代价最小、产码面最窄（只动 cmd/wisp 文案+其钉），但"机主要点开设置自己配模型"的原始诉求只是被绕开、没被兑现**——这句判读是否算"我替编排者裁了"：不算，这是把机主体验对照票面背景句量出来，选形仍归编排者。

## 4. 我可能判错的条目

- **N1（§0-5）**："七枚全写不进去"我是静态链推的；票 257 立票时 e2e 腿是真跑量过（A531/H11）。若常驻面板链在干净机上还有我漏看的建行/预置路径，此结论错。我已查 `panel_inbound.go` 全文与 `newComposerDispatchChain`：链上只有 `config.NewManager`（读），无 ensure/建行调用 ⇒ 静态链很硬，但禁 Go ⇒ 未能复跑。
- **N2（§2 两界对账红键集）**："红句点名 ≈7 键"是我拿 `jsonKeysOf` 反射算法（approval_test.go:174-187，只收带 json tag 的导出字段、嵌套 struct 不下钻——`ComposerState` 的 `mode`/`workspace` 是**嵌套 struct 键本身**，不下钻比对其内部）静态拼的。`tsInterfaceKeys` 的 regex（`composer_test.go` 文件内 `tsBlockRe`/`keyRe`）对 `panel.ts` 的解析同理静态拼。真跑可能多红/少红一两枚键名；**"两门今天已红"这一定性有台账三处独立佐证（145-growth-r2/r3 的 FAIL 名单），定性稳、键集清单是推算**。
- **N3（§1-6 预设名免校验通道）**："ⓐ-甲 塞预设名＝校验全绿"依赖 `applyPresets` 在 validate **之前**跑（loader.go:121-122 顺序）——顺序我读过，但"塞了预设名的空壳行能过 `validateAPIKeyRefs` 的 protocol 分支"是按 `validate.go:165-171` 的 `if p.Protocol != "" ... else if isPreset` 静态推的，没真跑过这份首建文件。
- **N4（§2 C17 步数）**："新方法名走 6 步"是全 grep 出来的编辑点清单，`internal/panel/l2_grant_boundary_test.go` 那把 AST 尺（routeNamePool 从全包字面收池）可能还有我没点名的断言会被新字面顶红——池尺的判定是"包里写了但 guard 不答才红"（l2_grant_boundary_test.go:1132），新方法若真在 guard 里答了就不红，但**写法顺序**（先加常量后加 case 的中间 commit）会让它红——本腿没量 CI 分步语义。

## 5. 判不动／量不到

- **N1 真跑**：AC#0-ⓐ/ⓑ/ⓒ 任何一格的"改后它真的绿/红"都需要 `go test`——⛔ 票面禁令（同机 255-r1）。全部格子以静态链交付。
- **N2 活跑红句**：`TestComposerContractTypesMatchFrontend`/`TestApprovalCardViewJSONKeysMatchFrontendTypes` 今天的**逐字红句**量不到（台账 145-growth-r2 载当时红句只点名 `git` 一枚——那是 09-26 读数，其后 197/200/248 各程落了 `tasks`/`instructions`/`credentialState` 等键，台账 145-snapshot-growth-r3:112 记"因本程两枚键更红"⇒ 红句键集随程增长，我只给静态推算集）。
- **frontend 改动面**：⛔ 票面禁 `frontend/**` ⇒ ⓑ 若选"新方法"形，页面侧要加哪几行量不到（只量到 `panel.ts` 是对账门的读取对象、`WispHostBridge`/`postMessage` 两处 call site 钉在 `bridge_test.go:157-164`）。
- **CI 分步语义**（N4）：不跑 CI 量不到。
- **`settings.go` 新 setter 的"半行"失败形状**：静态推（两枚守卫写中途失败）没有实测粒度。

## 6. 我推翻编排者哪一句

**没有推翻**。票面五行断言全部复核成立（§0）。两处**精度修正**（非推翻）：
1. `catalog.go:97-101` 那道拒的真身在 `:101-103`（case 头在 :97-100，拒句 :102-103）——判读不变，行号差 1。
2. 票面 ⓒ 的"`[providers.x]`"拼法若照抄进回执会把面板整条链打死（§3 首段）——这是票面文字的引用形问题，不是断言问题。
