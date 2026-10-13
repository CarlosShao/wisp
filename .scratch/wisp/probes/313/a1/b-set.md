# B. config.set 可写白名单 vs 313 缺口（只读普查腿 313-a1；HEAD 0e5253fe）

每行=结论＋尺＋件行号快照；行号=HEAD 现读，引用前重跑尺。

## B0 白名单逐枚（原文抄自 internal/panel/config_handlers.go:57-69）
```
FieldProviderBaseURL    = "provider_base_url"
FieldProviderAPIKeyRef  = "provider_api_key_ref"
FieldModelContextWindow = "model_context_window"
FieldModelPriceIn       = "model_price_in"
FieldModelPriceOut      = "model_price_out"
FieldRoleChatModel      = "role_chat_model"
FieldProviderCredential = "provider_credential"   // 只写：值不在共享封套上
```
- 选择器表（同件 :74-94）：要 provider={base_url, api_key_ref, context_window, price_in, price_out, credential}；要 model={context_window, price_in, price_out}；要值=除 credential 外 6 枚。
- 锁定族前缀（:108-134，逐字）：risk. / fs. / net. / plugins. / privacy. / models. / audio.（+ permission_mode、allowed_dirs 别名）⇒ 新字段名不得落这些前缀；`voice_*` 安全。
- 尺：`sed -n '57,94p;108,134p' internal/panel/config_handlers.go`

## B1 每枚的真路径（出处现读）
| 字段 | config.toml 真路径 | 出处 |
|---|---|---|
| provider_base_url | llm.providers.<p>.base_url | internal/config/settings.go:100-116 |
| provider_api_key_ref | llm.providers.<p>.api_key_ref | settings.go:121-137（先拒非 dpapi:/env: 形：cmd/wisp/panel_config_store.go:185-193；二道=validateAPIKeyRefs validate.go:189-196）|
| model_context_window | llm.providers.<p>.models.<m>.context_window | settings.go:143-166 |
| model_price_in/out | …models.<m>.price.in / price.out | settings.go:172-211 |
| role_chat_model | llm.roles.chat.model | settings.go:224-230（provider 那一半**没有写口**：requireChatProvider 拒，:235-243）|
| provider_credential | 值→DPAPI blob；引用落 llm.providers.<p>.api_key_ref | cmd/wisp/panel_config_store.go:238-279；值只走 `credentialValue` 键（:53、:59-66）|
- 写入纪律：一枚字段=一条 ownedKey、validate 前置于落盘（settings.go:25-28、:259-310；第 3 种拒因在 :288-294）。
- 尺：`sed -n '100,137p;143,166p;224,243p' internal/config/settings.go`

## B2 ★缺口点名：语音段今天零枚可写字段
- 313 AC#2 要"语音链（ASR/TTS/实时段）的 provider＋model 界面上改"——白名单 7 枚**全在 [llm]**，语音面一枚没有（B0 名单现读=证）。
- 目标键今天长这样（schema.go）：`voice.asr.provider/model`（:209-212，provider default="local-sherpa"）；`voice.tts.provider/voice/speed`（:215-219）；`voice.cloud_asr_chain`/`voice.cloud_tts_chain`（[]string，元素 "provider/model"；:259-265）；`voice.realtime.{enabled,provider,model,api_key_ref,base_url}`（:232-238）。
- **要新增哪几枚（形状建议＋为什么这个名）**：
  - `voice_cloud_asr_chain` / `voice_cloud_tts_chain`：值=链元素串（单元素如 "stepfun/<模型id>"，多元素逗号分隔）；名=键路径点→下划线（与 role_chat_model↔llm.roles.chat.model 同规）。校验=每元素过 checkChainElement（B3 逐字）。这是 audio3"仅 API"最贴的落点：链=云侧回退序列，元素就是 provider/model 对。
  - 实时段（对）：候选 α 一枚 `voice_realtime_pair`（值 "provider/model"，一次写 provider+model 两键）——撞"one key per call"纪律（settings.go:25-28 逐字"Each method owns exactly one dotted key path"；多键原子写"needs its own ruling"）⇒ 落地腿要先拿裁定。候选 β 两枚 `voice_realtime_provider`+`voice_realtime_model`——守 one-key-per-call，半对中间态过 validateRealtime（enabled=true 时缺一即拒，catalog.go:129-131）。
  - 可选 `voice_realtime_api_key_ref` / `voice_realtime_base_url`（audio3 若走 realtime 头自己的 key/base_url，schema.go:236-237）。
  - 变体（若机主指本地段选择器 asr/tts）：那两对**今天不被 catalog 校验**（validateCatalog 只查 chains+roles+realtime，catalog.go:24-44）⇒ 写口要自带新校验面，代价更大；且其语义是本地引擎选择器（default local-sherpa），不是 stepfun 落点。
- 尺：`grep -n 'CloudASRChain\|Realtime\|ASRConfig\|TTSConfig' internal/config/schema.go`

## B3 每枚必过的校验（现读逐字；票面 :95-113/:110/:113 已漂到 :99-115/:108/:112）
- catalog.go:99-100 都空 ⇒ 注释"unset: the owning feature falls back / disables itself"
- :101-103 只写 model ⇒ `"config.toml: %s.model set without %s.provider (set both or neither)"`
- :104-106 只写 provider ⇒ `"config.toml: %s.provider set without %s.model (set both or neither)"`
- :108-111 provider 不在 registry ⇒ `"config.toml: %s references unknown provider %q"`
- :112-115 model 不在该 provider ⇒ `"config.toml: %s references unknown model pair %q"`
- 链元素：checkChainElement（catalog.go:51-71）——:57-60 拒 local-sherpa（"the local sherpa cascade is the terminal voice fallback and must not appear in a cloud chain"）、:62-65 未知 provider、:66-69 未知 model；元素形状=第一个 '/'（:74-84）
- realtime：catalog.go:124-142——:129-131 enabled 时 provider/model 缺一即拒；:133-140 未知 provider/pair 拒；disabled 可带陈旧对
- 落盘前门：settings.go writeOneKey:288-294（同一份 validate；不过=第 3 种拒因、文件 0 字节动）
- 尺：`sed -n '51,71p;99,115p;124,142p' internal/config/catalog.go`

## B4 新字段落地要动的四处（现读点名）
1. internal/panel/config_handlers.go:57-94（字段表+选择器表）
2. cmd/wisp/panel_config_store.go:182-223 switch 新 case（不接=default 拒"字段 %q 没有对应的写入腿，未落盘" :221-222）
3. internal/config/settings.go 新 setter（各自一条 ownedKey）
4. ⛔ C17 方法白名单**不用动**（config.get/set 已在，bridge.go:66-67）；只有起新方法才碰它
- 尺：`sed -n '182,223p' cmd/wisp/panel_config_store.go`

## B5 "新增可写字段=契约变更"的规矩出处（本腿不写账，只点名）
- internal/panel/config_handlers.go:15-16 逐字："Growing this table is a contract change, not a bug fix."
- AGENTS.md 1.1：D1-D47 / C1-C32 契约变更须**人工批准**（SPEC-12 4.1 同旨）
- 账落点=docs/reports/pending-and-issues.md（A##/R##/Q## 只追加不删）
- 仅当起**新方法**才另碰 C17 白名单：frontend/src/lib/panel.ts:228-235（"reopens C17's method whitelist, a contract-level human-approval surface"）
- 尺：`sed -n '15,16p' internal/panel/config_handlers.go`

## B6 AC#2 的结构性缺口（比"缺字段"更重，落地腿必答）
- 今天链路**不能建行**：provider 行不存在⇒unknownProviderErr（settings.go:376-380 逐字"a settings write addresses an existing entry, it does not create one"）；catalog 行不存在⇒requireCatalogEntry（:356-370）；257 形ⓒ=页面不铸行（:36-39"the page does not mint a provider or a catalog row"）。
- ⇒"在界面上把 stepfun（audio3）配成"=先建 provider 行+audio 模型行（+文本 roles 对），今天一条都建不了。要么新增"建行/建节"写入能力（多键+建表，settings.go:25-28 点名"needs its own ruling"），要么维持手改文件——与机主"只在界面配"直接冲突。
- 尺：`sed -n '36,40p;356,380p' internal/config/settings.go`
