# 313-a1 — Go 侧只读现量（票 313：默认供应商＝stepfun、模型与密钥全在界面配）

**腿**：313-a1（只读现量腿，⛔ 产码；本件是本腿唯一写面）
**现量时刻**：2026-11-11 21:2x +08 · 分支 `dev` · HEAD `14b73de13`
**票面**：`.scratch/wisp/issues/313-default-provider-stepfun-and-model-config-stays-in-the-ui.md`（49 行全读；下文「票面 ：N」指该票行号，另注者除外）
**前勘**：票面 ：31 记 2026-10-13 首轮 313-a1（HEAD `0e5253fe`，四表在 `.scratch/wisp/probes/313/a1/`）；本件按「改前重跑」在现 HEAD 重现量，并吸收此后落地的 313-r1/r2/r3/v1（票面 ：32-：37、：40-：47）。
**行号纪律**：本件所有 file:line 均为本次现读（read/grep 现跑），非转抄。

---

## 逐格判定

| 格 | 判定（Go 侧） | 缺的部分归谁 |
|---|---|---|
| AC#1 默认选中 stepfun | **部分满足**（机制齐、Go 不持常量——这是 A877 裁定③的刻意形状，非缺口；页面消费物已由 313-page 落地） | 页面侧归前端会话（`dsh/feat/frontend-p0-v2` 树，本编队零写面）；真机端到端归机主签收 |
| AC#2 两个模型都能配 | **已满足**（Go 半；313-v1「成立」判语在案，本腿现读复核一致） | 其余三个文本角色（memory_extract/handoff/summarize）无面板字段——票面 AC#2 点名的是 chat，严格读不算 313 缺口；若产品要扩，扩面在 Go（本编队可动）＋页面 |
| AC#3 界面里录 key | **已满足**（机制＋三条守则的测试判据全在，逐一具名） | 无 Go 侧缺口 |

---

## AC#1 默认选中——Go 侧有没有「stepfun 是默认供应商」的落点

**有预设、无默认选中位；默认选中的承载物今天按裁定③放在页面，Go 送旗。**

1. **预设表里有 stepfun**：`internal/config/defaults.go:22-35` `providerPresets`（12 家），`:34` 逐字
   `"stepfun":     {Protocol: ProtocolOpenAIChat, BaseURL: "https://api.stepfun.com/v1"}`（即票面 ：10 那行，现读未漂）。`PresetNames()`（`defaults.go:39-46`，按字母序）就是面板供应商选择器的数据源。
2. **「默认选中」在 Go 侧今天零承载物——这是刻意的，不是没做**：
   - 全仓 `*.go` grep `stepfun`（79 处）：生产码里 stepfun 只出现在①预设表、②注释/测试。**没有任何「provider 未配时回落 stepfun」的运行时分支**；role 未配的语义是「该功能自降级/回退」（`internal/config/catalog.go:99-100` 注释逐字），不是「回落默认供应商」。
   - 默认标签层（D36 规则 3，`defaults.go:54-62`）也不设默认 provider：`default:"..."` tag 是唯一默认源、TOML 只覆盖，而 schema 里没有一枚 tag 指向 stepfun。
3. **⛔ 硬编码写死的判据（现读到的守卫）**：
   - Go 读面送「用户改过没有」的旗：`cmd/wisp/panel_config_store.go:112-113`
     `Configured: len(cfg.LLM.Providers) > 0 || strings.TrimSpace(cfg.LLM.Roles.Chat.Provider) != ""` ＋ `ConfiguredKnown: true`——用户建过任何 provider 行或配过 chat 角色，页面预选就必须让位。
   - 该旗的契约注释在 `internal/panel/config_handlers.go:175-177` 逐字："It is the \"改过以用户为准\" flag: the page default-highlights stepfun only while this is false (A877 ruling 3)"——即 A877 裁定③：**默认高亮常量归页面侧持有（`OWNER_DEFAULT_PROVIDER="stepfun"`），且必须核名字真在 Go 预设清单里**（落地记录＝票面 ：48 313-page：两道门＋`useState` 仅初始值、无 effect 强加、不点添加不写盘）。
   - 判据成立 ⇒「不是硬编码写死」＝①运行时零 stepfun 特判分支；②预选只在 `configuredKnown && !configured` 时发生；③用户配置一经存在即以用户为准。
4. **成对规则仍是硬门**：`internal/config/catalog.go:99-106`（provider/model「set both or neither」）、`:108-115`（provider 必须在 catalog、model 必须存在于该 provider 下）——「把 stepfun 设为默认」在任何层都不能只写半枚（票面 ：11 现读未漂）。

**小结**：Go 侧该给的两样（预设清单＋`configured`/`configuredKnown` 旗）都在载荷里（键名见 AC#2 读面）；"stepfun" 三个字母本身归页面常量。Go 侧无剩余缺口。

---

## AC#2 两个模型都能配——设置 schema/API 今天支持到什么程度

**读面与写面都齐，键名逐一现读如下（313-v1「Go 半成立、三方键集全等」判语与本腿现读一致）。**

### 读面（`config.get` 的结构载荷）

- 载荷是 `config.get` 回执里**散文之后的一行 JSON**，行首标记 `wisp-config-view:`（`internal/panel/config_handlers.go:597` 常量、`:588-592` 契约注释、`:381` 与 248 散文同行返回、`:611-679` `renderSettingsPayload` 渲染本体）。
- 顶层键（冻结钉在 `internal/panel/config_313_test.go:321`）：`readable` `credential` `presetsKnown` `presets` `configuredKnown` `configured` `chatProvider` `chatModel` `chatModelKnown` `providers` `voice`。
  生产写入点：`presets`＝`config.PresetNames()`、`chatProvider`、`configured` 两枚旗在 `cmd/wisp/panel_config_store.go:110-113`。
- provider 行键（`config_313_test.go:348` 冻结）：`name` `apiKeyRef`（引用、非值）`hasBaseURL` `baseURL` `protocol` `billing` `modelCount` `models`（**模型名清单**，票 248 只有计数、313 补的名字）`credentialSet` `credentialError`。生产写入点 `panel_config_store.go:135-152`（`:146-152` 收集并排序模型名）。
- 语音块键（`config_313_test.go:361` 冻结）：`cloudASRChain` `cloudTTSChain` `realtimeEnabled` `realtimeProvider` `realtimeModel` `aecEnabled` `inputDevice` `sampleRate`（后三枚＝2026-11-11 本机音频旁支，读者门控、nil 即不上线）。生产写入点 `panel_config_store.go:114-123`。

### 写面（两表并查；`ApplySetting` 分派在 `cmd/wisp/panel_config_store.go:226-392`）

- **票 248 第一张表**（`internal/panel/config_handlers.go:91-94` 名册，7 枚）：`provider_base_url` `provider_api_key_ref` `model_context_window` `model_price_in` `model_price_out` `role_chat_model` `provider_credential`（`:59-69` 常量；`provider_credential` 是只写枚，值不走它，见 AC#3）。
- **票 313 第二张表**（`internal/panel/config_write_313.go:114-128` 名册，13 枚）：
  - `voice_cloud_asr_chain` / `voice_cloud_tts_chain` → `voice.cloud_asr_chain` / `voice.cloud_tts_chain`（值＝逗号分隔 `provider/model` 链；`panel_config_store.go:493-500` `splitChainValue` 按逗号切；多元素链有钉 `Test313MultiElementCommaChainKeepsEveryElement`，票面 ：47 313-r3）
  - `voice_realtime_provider` → `voice.realtime.provider`（setter `internal/config/settings_313.go:119-126`）；`voice_realtime_model` → `voice.realtime.model`（`:132-137`）
  - `provider_row_create` → 显式建 `[llm.providers.<name>]` 行（**只收预设名、行已存在即拒**，`settings_313.go:186-190` 附近拒绝句「已有行请用逐字段写口」；A877 裁定②的显式建行形状）
  - `provider_billing` / `provider_row_create_custom`（票 320）/ `model_row_create` / `model_row_delete` / `provider_row_delete` / `role_chat_pair`（chat 角色 provider+model **原子成对写**，A881 裁定①②）
  - 本机音频旁支：`voice_aec_enabled` / `audio_input_device`（`:108-109`）
- 并查 accessor＝`config_write_313.go:176-183`（union 单一入口；另有 187/65/wave 三张旁支名册）。
- **端到端实证**：`cmd/wisp/panel_config_313_test.go:143-159`——空配置起六发写穿真实门（`provider_row_create` stepfun → 两条云链 `stepfun/<model>` → realtime 两枚 → `role_chat_pair`），落盘 config.toml 逐字含 `[llm.providers.stepfun]`、`[llm.roles.chat]`、`cloud_asr_chain`、`cloud_tts_chain`、`voice.realtime`；回执载荷各键现读核（`:175-193`）。
- 语音链的 stepfun API 形：链元素就是 `stepfun/<model>` 字符串（`internal/speech/cloudendpoint.go:88` 注释以 stepfun BaseURL 为例），Go 侧无形状障碍。

### schema 侧（什么能配进 config.toml）

- `llm.providers.<name>`：`internal/config/schema.go:401-409`（`Provider`：protocol/base_url/**api_key_ref**；`:396-400` 注释逐字列 12 家预设含 stepfun，空字段继承预设默认）。
- `llm.roles.*` 四角色＋成对/可知性校验：`internal/config/catalog.go:89-119`。
- `voice.realtime.{enabled,provider,model,api_key_ref}`：`schema.go:252-258`＋`validateRealtime`（`catalog.go:121-139`：enabled=true 时必须成对且可知）；`voice.realtime.api_key_ref` 在档位表 `internal/config/tiers.go:81` 是 TierReload。
- 255 读者名册（设置读行的证据行）：[llm] 行在 `cmd/wisp/config_readers_255.go:126`（锚 token＝`panel_config_store.go:98` 的 `ChatModel:` 行）、[audio] 行在 `:246`（S 腿重裁后落点）。

### 如实记的两枚边角（都不翻本格判定）

1. **`voice.realtime.enabled` 没有写字段**：读面有 `realtimeEnabled`，两张名册（7＋13 枚，逐字现读）里没有 realtime_enabled 的开关位。用户想开实时头今天仍需手写 config.toml。票面 AC#2 只要求「能配 provider＋model」，故不判缺；要补的话落点在 Go（第二张表加字段＋setter）＋页面，本编队可动 Go 半。
2. **其余三个文本角色**（`llm.roles.memory_extract/handoff/summarize`，`catalog.go:91-96`）在 `internal/panel` 零字段（grep `memory_extract|handoff|summarize` 于 internal/panel 零命中）：票面 AC#2 原文点名 `llm.roles.chat` 等——严格按票面只判 chat，chat 已齐；三角色若要上界面是新需求，Go 侧可动。

**小结**：chat 角色对＋ASR/TTS/实时三段的读面与写面今天全通（313-r1/r2/r3 落地、v1 验收、r3 补钉，票面 ：32-：37、：40-：47）；Go 侧无本票射程内的缺口。

---

## AC#3 界面里录 key——`api_key_ref`（DPAPI blob）机制今天落到哪一步

**机制链完整：录入（面板→DPAPI blob→引用落盘）→读取（引用→存在性，永不读值）→三条守则各有测试判据在守。**

1. **存储层**（`internal/secret/`，C28）：blob 每引用一枚存 `<data>\secrets\` 下；Windows 用 DPAPI `CryptProtectData` CurrentUser；包职责注释逐字含「last-4-only logging」「便携模式回退 env: 永不明文」（`internal/secret/doc.go:1-22`）。`store.go`：`NewRef`（随机 id，:47）→`Store`（:60）→`Exists`（:122）→`Resolve`（:139，内建 last-4 脱敏 :148-163）。
2. **schema 层**（D36 规则 5）：`api_key_ref` 只收 `"dpapi:<blob-id>"` 或 `"env:NAME"`，明文结构性放不进——provider 位 `internal/config/schema.go:406-409`、realtime 位 `:258`；加载闸二次验形（`internal/config/validate.go:261-271` provider、`:299-301` realtime；测试 `validate_test.go:293-308`）；面板写入口再验一次（`panel_config_store.go:229-238`，拒绝句点名「只接受 dpapi: 或 env:」）。
3. **录入路径（值只走一条专用窄门）**：`credentialValue` 是全仓唯一携带凭据值的 JSON 键（`panel_config_store.go:53` 常量），**只声明在写专用封套** `credentialWriteRequest`（`:59-66`），共享封套 `panel.ComposerRequest` 不含它——值从线上读下即交 `secret.Store`（`:460-463`：先铸随机 ref、再存值、最后写引用进配置），回执只回 `{ref, recorded, provider}`（`:473`），回执注脚逐字「值本身不在配置里，也不在任何回话里」（`:474`）。`ApplySetting` 若收到 `provider_credential` 字段会点名拒绝（`:387-389`「凭据值不走这条写入路径…必须走 StoreCredential」）。
4. **读取路径（值永不回面板）**：`ReadSettings` 只问存在性——注释逐字 "Existence is asked, content is not"（`panel_config_store.go:85-87`；实现 `:153-166` `refRecorded`：dpapi 走 `Exists`、env: 查进程环境 `:195-200`）；聚合答案只是计数（`:175-184`：`no_ref_declared/all_recorded/partly_missing`，缺哪家不点名）。页面字段契约：连 last-4 都不上页面（`internal/panel/config_handlers.go:140-142` 逐字：last-4 是 CLI 终端的工具，不是页面字段）。
5. **三条守则的判据（grep＋现读测试件）**：
   - **不回显**（回执/载荷面）：哨兵 `canary248`＝`cmd/wisp/panel_config_248_test.go:39`，走真实写入路径后对六面扫值（回执/audit/slog sink/快照/台账行/data-root，`:186-200`）零命中，且有正控证明尺不瞎（`:218-259` `TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere`）；路由面回执含 ref、不含值、含「只写不回显」话术（`internal/panel/config_route_248_test.go:348-377`）；常备镜像尺 `internal/panel/inbound_raw_leak_35r7_test.go:49-56`（canary 逐字镜像＋通用厂商键形状规则）对回执/审计双扫；载荷面 `renderSettingsPayload` 注释钉「no credential value and no partially-redacted value」（`config_handlers.go:611-616`）＋键集冻结且**零值承载键**（`config_313_test.go:387-397`：`credentialvalue/apikey/apikeyvalue/value/secret/token/key` 全是禁区词；`apiKeyRef` 是 schema 写的引用、与任何值键不同名）。
   - **不进日志**：last-4 契约三件钉——`internal/secret/store_test.go:257-299`（日志路径只准 `RedactSecret` 形式＋`TestRedactSecret` 钉函数本体，`internal/secret/redact.go:7-13`）、`internal/observe/logging_test.go:30`（日志属性面）、`cmd/wisp/secret_test.go:203` 与 `:238`（CLI 输出面：整值零命中、last-4 形式在）。
   - **不进模型可读文本**：结构钉＝共享封套不得声明 `credentialValue`（`panel_config_248_test.go:344-347`「写专用封套必须保有、共享封套必须没有」＋`:379-380` 红句）——值进不了 agent 侧能读到的任何类型；行为钉＝值路径的六面扫描含快照的 results 通道与台账行（即模型可读面，`:186-200`）；注入面钉＝把 key 形状塞进消息路由的 text 字段（共享封塞 credentialValue／未列 method／伪造 source／无 handler 五个子测，`inbound_raw_leak_35r7_test.go:115-203`）路由产物必须形状干净，正控 `:208-222`；上游先例＝票 63「key 从不进 argv/日志/对话」家族（票面 ：19 认引）。
6. **值真正只存一份**：`TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef`（`panel_config_248_test.go:263-294`）——ref 能解回原值、config.toml 只有 `api_key_ref = 'dpapi:…'`、secrets 目录恰好一枚 blob。
7. 313-v1 验收腿已亲跑「密钥面 canary 四表面零命中」（票面 ：37、：42），本腿为静态现读，未重跑活体。

**小结**：三条守则今天各有具名测试判据在守，机制无缺口。

---

## 一句总判

**Go 侧今天三格不缺产码：AC#1 部分满足（Go 该给的预设清单＋`configured` 旗已上载荷，「stepfun 默认」常量按 A877③ 刻意归页面且 313-page 已落地）、AC#2/AC#3 已满足（写读两表 7＋13 枚字段＋成对硬门＋DPAPI 窄门与三条守则判据全在）；剩下归页面侧与真机签收的是——313-page 产物在 `dsh/feat/frontend-p0-v2` 树（本编队零写面）、真机双击端到端归机主，Go 侧可选残山＝`voice.realtime.enabled` 开关位与三个非 chat 文本角色的面板字段（票面未要求，属新需求）。**

## 本腿门禁自证

- 只读腿：未跑测试、未改任何 `.go`/`frontend/**`/`cmd/wisp/**`；写面仅本件。
- 现读清单（HEAD `14b73de13`）：票面全文；`internal/config/{defaults,catalog,schema,settings_313,tiers,validate,validate_test}.go`；`cmd/wisp/{panel_config_store,config_readers_255}.go`；`internal/panel/{config_handlers,config_write_313,config_313_test,config_route_248_test,inbound_raw_leak_35r7_test}.go`；`cmd/wisp/{panel_config_248_test,panel_config_313_test,secret_test}.go`；`internal/secret/{doc,redact,store,store_test}.go`；`internal/observe/logging_test.go`；`internal/speech/cloudendpoint.go`。
- grep 面＝`stepfun` / `canary` / `RedactSecret` / `credentialValue` / `api_key_ref|ParseRef` / `Field[A-Za-z0-9]+ = "` / `memory_extract|handoff|summarize` / `renderSettingsPayload|SettingsPayloadMarker`（均含 `*_test.go`）。
- git：`add` 与 `commit` 两处显式 pathspec 仅本件；⛔ push；⛔ 动 thresholds/golden/PLAN.md/docs/specs。
