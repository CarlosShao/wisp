# 257-a2 普查件 — 凭据通道名册 / 写侧建行卡点 / 两条重读通道（只读，零 Go 命令）

> 腿：`257-a2`；票：257（AC#0 已裁形 ⓒ，撤销口令「257 改形 ⓐ」）。
> 前件：`.scratch/wisp/probes/257/a1/census.md`（92 行，commit `ed3fd270`）。
> ⛔ 本腿零产码写入、零 `go` 命令（同机三枚写腿在飞：`cmd/wisp`／`internal/config`／`internal/agent/approval`／`internal/tools`＋`internal/risk`）。
> ⛔ 凭据**值**零外泄：本件只写变量名／字段名／blob 名／env 名。
> ⛔ `frontend/**`／`design/**` 两层禁令：未读、未引用、未转述。

## 0. 锚

| 项 | 读数 |
|---|---|
| 进场时刻 | `2026-10-03T09:12:30+08:00` |
| 进场 HEAD | `3736f0dd`（2026-10-03 09:11:55 +0800） |
| 进场 porcelain 行数 | 398 |
| 量程中 | （取数中） |
| 落件 | （取数中） |

## 1. 凭据通道名册（谁解析 / 何时读 / 失败句子逐字）

**先立一件事**：全仓只有一处做加解密（`internal/secret/dpapi_windows.go:24` `CryptProtectData`／`:39` `CryptUnprotectData`；非 Windows 编译垫片 `internal/secret/protect_other.go:21-27`）。`grep CryptProtectData|CryptUnprotectData` 在 `cmd internal tools docs scripts` 下只命中这一包（另一处是 `internal/secret/doc.go:1` 的注释）。**所以下面每一条的"读凭据"最终都落到 `internal/secret`——差别只在谁递进一个 ref、什么时候递。**

### 1.0 引用形只有两种，且写死在解析器里

`internal/secret/refs.go:15-21` 两枚常量：`RefKindDPAPI = "dpapi"`／`RefKindEnv = "env"`（前缀 `dpapi:`／`env:`）。解析者＝`ParseRef`（`internal/secret/refs.go:29-46`）：

| 形 | 校验 | 不合法时的逐字句 |
|---|---|---|
| `dpapi:<blob 名>` | `ValidBlobID`（`refs.go:51-67`：1–128 字符 `[A-Za-z0-9._-]`，非 `.`/`..`） | `secret: invalid ref: %q: blob id must be 1-128 chars of [A-Za-z0-9._-] and not "."/".."`（`refs.go:34`） |
| `env:<环境变量名>` | 非空、≤128、不含空白/`=`/`:`（`refs.go:39`） | `secret: invalid ref: %q: env name must be 1-128 chars without whitespace, "=" or ":"`（`refs.go:40`） |
| 其它一切形（**含明文 key**） | 落到 `default` 分支（`refs.go:43-45`） | `secret: invalid ref: %q must start with "dpapi:" or "env:"`（`refs.go:44`） |

⇒ **"配置文件里的引用形"与"明文"不是两条通道**：明文在 `ParseRef` 就被拒，且 `internal/config` 的 schema 里根本没有承载它的字段（`internal/config/boundary_test.go:36-45` 那条 D36 rule 5 的尺点名 `APIKeyRef` 是唯一可带密钥材料的字段）。

### 1.1 通道甲：`dpapi:<blob 名>`（DPAPI blob 文件）

- **落点**：一 ref 一文件，`<dataDir>\secrets\<blob 名>`（`internal/secret/store.go:42` 拼 `secrets`，`store.go:172-179` `blobPath` 是唯一路径咽喉；目录由 `winsec.PrivateDirAll` 密封，`store.go:49`）。
- **谁写入**：只有两枚生产入口——① `wisp secret set <name>`（`cmd/wisp/secret.go:333-387` `runSet`；值来自隐藏控制台输入 `secret.go:588/592` 或 `--from-stdin` `secret.go:576`，**argv 里永不存在值的槽位**：全文件无 string 型 flag，`secret.go:314-315` 注释点名这条不变量）；② 面板凭据写入腿 `StoreCredential`（`cmd/wisp/panel_config_store.go:238-279`：`secret.NewRef()` 现制随机 blob 名 `panel_config_store.go:263` → `s.secrets.Store(ref, value)` `:264` → 再把引用写进配置 `:267`）。
- **谁读**：`Store.Resolve`（`internal/secret/store.go:89-123`），`case RefKindDPAPI` 在 `:104-119`：`os.ReadFile` → `unprotect`。
- **什么时候读**：**每次构建一个 Endpoint 时读一次**，不是每次请求。链路：`cmd/wisp/run.go:435` `llm.NewResolver(cfg, st)` → `:436` `res.ResolveRole(llm.RoleChat)` → `internal/llm/resolver.go:136` 取 `p.APIKeyRef` → `:141` `r.Keys.Resolve(ref)` → 结果塞进 `Endpoint.APIKey`（`resolver.go:146`，那行注释自陈 "resolved plaintext; NEVER log"）→ `cmd/wisp/run.go:442` `BuildEndpointProvider` → 适配器把它留在 `opts.APIKey` 里（`internal/llm/provider.go:167`）→ **此后每次 HTTP 请求是从内存里那份副本取头**，不再碰 store：`internal/llm/openaichat/adapter.go:88-89`（`Authorization: Bearer`）／`internal/llm/openairesponses/adapter.go:124-125`（同）／`internal/llm/anthropic/adapter.go:154-156`（`x-api-key`）。
- **失败时用户看到的句子**（逐字，从下往上叠）：
  - blob 文件不在：`secret: resolve %q: no blob file under %s`（`store.go:107`）
  - 读得出但解不开（非便携档）：`secret: resolve %q: secret: dpapi CryptUnprotectData: %w`（`store.go:117` 叠 `dpapi_windows.go:40`）
  - 便携档解不开：`ErrPortableDecrypt` 全文（`store.go:17-19`）逐字 = `secret: portable mode: DPAPI blob cannot be decrypted (created by another user or on another machine); store this key as an env: reference instead (P13: there is no plaintext fallback)`——**注意这一支今天在生产里到不了面**：`WithPortable` 只被 `cmd/wisp/secret.go:297`（`wisp secret` 那一腿）传；其余四枚生产 `secret.NewStore` 都不传（`cmd/wisp/run.go:392`、`cmd/wisp/panel_inbound.go:262`、`cmd/wisp/providers.go:92`、`cmd/wisp/firstrun.go:79`）。⇒ 这句是**静态现读**（`grep secret.NewStore` 的命中集就是全部装配点），不是量不到。
  - 非 Windows 平台：`secret: DPAPI is only available on Windows; use an env: ref on this platform`（`internal/secret/protect_other.go:18-19`，由 `store.go:75`/`:112` 透传）
  - `wisp secret get <name>` 直接把上面那串原样打到 stderr：`wisp secret get: %v`（`cmd/wisp/secret.go:419`）
  - 面板侧**不解密、只 stat**：`refRecorded` → `s.secrets.Exists(ref)`（`cmd/wisp/panel_config_store.go:166` → `internal/secret/store.go:129-145`），错误进 `row.CredentialError`（`:121`），页面句 = `无法核对：` + 那句（`internal/panel/config_handlers.go:476`）
  - 凭据字段走错门：`凭据值不走这条写入路径（它只写不回显，必须走 StoreCredential），字段 %q 未落盘`（`cmd/wisp/panel_config_store.go:219-220`）
  - 存得进 store、引用没写进配置（半状态句）：`密钥已进凭据存储，但 %q 的引用没写进配置：%w`（`panel_config_store.go:272`）

### 1.2 通道乙：`env:<环境变量名>`

- **谁读**：同一枚 `Store.Resolve` 的 `case RefKindEnv`（`internal/secret/store.go:95-103`）→ `os.LookupEnv(value)`（`store.go:96`）。**每次都是即时读进程环境**（没有缓存层），但调用时机仍受 §1.1 的"每次建 Endpoint 一次"约束。
- **失败句逐字**：`secret: resolve %q: environment variable %s is not set`（`store.go:98`）／`secret: resolve %q: environment variable %s is empty`（`store.go:101`）。两式都不返回兜底值（`store.go:86-87` 注释：missing or empty -> explicit error, never a default）。
- **面板读侧走的是另一把尺**，不是 store：`refRecorded` 对 `env:` 支自己 `os.Getenv(value)` 判非空（`cmd/wisp/panel_config_store.go:157-161`），注释给的理由是 `Exists` 故意拒绝 `env:` 形（`internal/secret/store.go:135-136` 逐字：`secret: exists %q: env refs live in the process environment; only dpapi: refs have blob files`）。⇒ **"这一页说 env: 未录入"与"wisp run 起不来"可能同时成立**，因为两把尺问的是同一件事的两个进程环境（§3 会用到这条）。
- **写入面**：`Store.Store` 明确拒绝 `env:` 形（`internal/secret/store.go:69-71`：`secret: store: %q: env refs live in the process environment; only dpapi: refs are persisted`）。⇒ **env: 这一支没有任何 Wisp 写入者**，值只能由系统环境提供——`cmd/wisp/firstrun.go:115` 的回执句正是这么写的（"值由系统环境提供"）。

### 1.3 通道丙：`config.toml` 里的引用形（只是"指名"，不是凭据本身）

- 两枚字段：`llm.providers.<名>.api_key_ref`（`internal/config/schema.go:387` `toml:"api_key_ref"`）与 `voice.realtime.api_key_ref`（`schema.go:236`）。
- **谁解析形状**：加载期 `validateAPIKeyRefs`（`internal/config/validate.go:152-191`，逐枚调 `secret.ParseRef`；`llm` 侧拒句 `config.toml: llm.providers.%s.api_key_ref: %v` 在 `:156-157`，`voice` 侧 `config.toml: voice.realtime.api_key_ref: %v` 在 `:186-187`）；写侧还有一道**更早**的面板门（`cmd/wisp/panel_config_store.go:187-192`，句：`引用形状不合法（只接受 dpapi:<id> 或 env:NAME）：字段 %q 未落盘`）。
- **什么时候读**：加载期只读"名字"；解引用在 `resolveRefs`（`internal/config/loader.go:208-233`）——**但这一条在生产里是死线**：

  ⚠ **现读硬事实（本腿与 `257-a1` 不同的一处）**：`resolveRefs` 只在 `res != nil` 时干活（`loader.go:210-212` 直接返回空 `Resolved`），而**每一枚生产 `config.NewManager` 都传 nil**：`cmd/wisp/run.go:409`、`cmd/wisp/panel_inbound.go:230`、`cmd/balldebug/main.go:231`（其余 `LoadFile(..., nil)`：`cmd/wisp/providers.go:98`、`cmd/wisp/models.go:184`）。`Manager.Resolved()`（`internal/config/manager.go:123-134`）**在非测试代码里零调用者**（`grep ProviderKeys|RealtimeKey|\.Resolved\(\)` 命中集＝定义＋`internal/config/loader_test.go`＋`cmd/wisp/secret_test.go`）。
  ⇒ 结论：**"配置层解引用凭据"这条通道今天存在但没接线**；真正读凭据的只有 §1.1/§1.2 那两条经由 `llm.KeyResolver`/`secret.Store` 的线。这不是缺陷上报（`loader.go:37-40` 与 `resolver.go:86` 都写明 res/Keys 可为 nil），但**回执文案不能把这条例为"入口"**——`firstrun.go` 今天也没列它，见 §3。
  ⇒ 副作用一条：`voice.realtime.api_key_ref` 今天**没有任何运行时消费者**（`Resolved.RealtimeKey` 零读者；`voice.realtime` 族在 `internal/config/tiers.go:65-70` 注册为 reload，但 `OnReload` 在 `cmd/wisp` 零赋值，见 §3）。

### 1.4 通道丁：别针（现读到的其它"能读到凭据"的地方）

| 入口 | file:line | 性质 |
|---|---|---|
| `wisp secret get <name> [--show]` | `cmd/wisp/secret.go:390-434`；默认掩码 `RedactSecret`（`internal/secret/redact.go:7-13`，只留后 4 枚字符），`--show` 是全仓**唯一被批准的明文出口**（`secret.go:428-432`：只写 stdout、不留尾换行、不记日志，并先打一句警告 `wisp secret get: --show writes the plaintext to stdout; do not redirect it to a file or a log`） | 人工终端 |
| `wisp secret list` | `cmd/wisp/secret.go:437-469` → `internal/secret/list.go:34+`（`BlobInfo` 只带 ID/Created/Modified/Size，`list.go:14-16` 注释自陈 "the files are not [read]"） | 只元数据 |
| `wisp secret unset <name> [--force]` | `cmd/wisp/secret.go:473-535`：先用 `secret.RefFieldNames`（`internal/secret/configrefs.go:62-75`）扫 config.toml 引用，被引用则拒绝并逐行列出点名的字段路径（`secret.go:503-507`） | 只读引用、不读值 |
| `wisp providers discover/probe` | `cmd/wisp/providers.go:124-133`（discover 直接 `st.Resolve`；拒句 `wisp providers: provider %q 的 api_key_ref 无法解析（Unconfigured）` 在 `:129`）／`:168-174`（probe 借同一枚 `llm.NewResolver`＋`ResolveChain`） | 每次命令运行解一次 |
| `cmd/llmrecord`（黄金语料录制，dev 工具） | `cmd/llmrecord/main.go:62-67` 直接 `os.Getenv(keyEnv)`，**不经 store**；空值句逐字 `environment variable %s is empty: live golden recording is deferred until an LLM API key is provided (H2); the mockllm + replayer harness needs no key and is already in place` | 第五形：绕过 `ParseRef` 的裸环境变量读取，只活在 `cmd/llmrecord` 里 |
| `tools/signmodels` | `tools/signmodels/main.go:117` `os.Getenv("WISP_MINISIGN_KEY")`（minisign 签名私钥**路径**，构建期工具，非运行时 LLM 凭据） | 构建期 |

⇒ **机主真走得通的凭据入口只有两枚**：`wisp secret set <blob 名>`（CLI）与设置页那枚 `provider_credential` 字段（面板）。`cmd/wisp/firstrun.go:111-118` 的回执今天只教了前者。

### 1.5 谁带 portable 旗（决定 §1.1 那句 P13 话到不到得了面）

- 带：`cmd/wisp/secret.go:296-298` `openStore()` → `secret.NewStore(c.dataDir, secret.WithPortable(c.portable))`，`portable` 来自 `resolveSecretLayout`→`proc.ApplyPortableOverride`（`secret.go:106-118`）。
- 不带（四枚全裸传）：`cmd/wisp/run.go:392`、`cmd/wisp/panel_inbound.go:262`、`cmd/wisp/providers.go:92`、`cmd/wisp/firstrun.go:79`。
- ⇒ 判语：**"便携模式下 blob 解不开"这句带指引的话，今天只有 `wisp secret` 一条腿会说**；`wisp run`／面板只会说 `secret: resolve %q: secret: dpapi CryptUnprotectData: %w` 那种裸密码学失败。SPEC-02 §6 的便携×DPAPI 待定案项（AGENTS.md §2 列了）在码里的形状就是这一处缺口。〔此条是"通道存在但旗没人举"的现读，不是量不到。〕

## 2. 写侧建行卡点与最小改动面

**只给改动面与代价，不选形**（选形归编排者；票 257 §8 已裁 ⓒ，本节是给"哪天真要建行"那张预留票备料）。

### 2.0 入向路由（先确认门在哪几枚函数上，逐枚现读）

页面 → `config.set`（方法名常量 `internal/panel/bridge.go:67`，白名单判定 `bridge.go:148`）→ `internal/panel/composer_dispatch.go:202-206`（把 `raw` 原始字节一并递给处理器，凭据值就靠这条路不走共享封套）→ `panel.ConfigWriteHandler.HandleConfigRequest`（`internal/panel/config_handlers.go:269`）→ `write`（`:307`）→ 六道名字门：字段非空 `:309-314`、锁定族 `:315-320`、在名册内 `:321-326`、要服务商的没指 `:327-332`、要模型的没指 `:333-338`、要值的给了空 `:349-354` → `Store.ApplySetting`（`:355-358`）→ **cmd/wisp 的实现** `configStore.ApplySetting`（`cmd/wisp/panel_config_store.go:173-230` 的 switch，逐枚映射到 `internal/config` 的 setter）／凭据支 `configStore.StoreCredential`（`panel_config_store.go:238-279`）。

装配点两枚：CLI 入向腿 `cmd/wisp/panel_inbound.go:266-267`、常驻面板腿 `cmd/wisp/panel_resident_windows.go:170`（两枚都经由 `newComposerDispatchChain`，`panel_inbound.go:228`）。

### 2.1 "行已存在"具体卡在哪一行：**双层，而且读的是盘不是内存**

| 层 | 位置 | 干什么 | 拒句（逐字） |
|---|---|---|---|
| 门① `check`（写前校验用的那道） | `internal/config/settings.go:71-73`（base_url）／`:92-94`（api_key_ref）／经 `requireCatalogEntry` 的 `:124`（context_window）与 `:169`（price.in/out）→ 真身在 `:294-308` | provider 行不在 ⇒ 直接返回错 | `config: no provider %q in this config; a settings write addresses an existing entry, it does not create one`（`settings.go:311-313`）；model 不在 ⇒ `config: llm.providers.%s has no model %q in its catalog; refusing to invent one from a settings write`（`:303-305`） |
| 门② `apply` 返回 false | `settings.go:63-66`／`:84-87`／`:111-118`（provider 与 model 两级判空）／`:150-156` | 就算绕过门①，落笔时仍然不建 | `config: <key> does not exist in the file this write is based on; nothing was written`（`:224-227`） |
| 门③ 写前 `validate` | `settings.go:230`（`validate(candidate)`），失败句 `:231-232` | 建行也过不了这一道，见 §2.2 | `config: refused to write <key>; the value would not survive validation, so config.toml is unchanged` |

**★ 本腿与 `257-a1` 的一处精度差**（不推翻其结论，改的是"哪本账"）：门①②③共用的 `base` 是**磁盘上那份文件**，不是这个 Manager 的内存。证据链：`writeOneKey` 第一件事是 `m.readCurrentFile()`（`settings.go:205`），`readCurrentFile`（`:274-289`）里 `stat` 通了就 `readConfigFile(m.path)`（`:283`）——`candidate` 由它而来（`:212`），只有**文件不存在**时才退回内存/默认表（`:213-218`）。
⇒ 现读推论（可今天证实）：**机主手加一节 `[llm.providers.<名>]` 之后，面板那七枚的写门立刻能过**，连那个"自己不 tick 的 Manager"也一样能过（§3.2）；卡住面板的只有"内存里的读数"（`ReadSettings` 用 `s.mgr.Config()`，`panel_config_store.go:92` ⇒ `internal/config/manager.go:115-119` 返回的是内存深拷贝）。**"写门看盘、读面看内存"是同一枚 leg 上的两把尺**——这正是 owner 那句"我填了为什么没生效"的反向形状（他更常问的是"我改了为什么看不到"）。

### 2.2 干净机上七枚各自的死因（逐枚过门，现读）

| 字段（`config_handlers.go:57-69`） | setter | 死在第几门 | 具体行 |
|---|---|---|---|
| `provider_base_url` | `SetProviderBaseURL` | 门①+门② | `settings.go:63-66`/`:71-73` |
| `provider_api_key_ref` | `SetProviderAPIKeyRef` | 门①+门②（另有一道更早的形状门：`panel_config_store.go:187-192`） | `:84-87`/`:92-94` |
| `model_context_window` | `SetModelContextWindow` | 门①(`requireCatalogEntry`)+门② | `:111-118`/`:124` |
| `model_price_in` / `model_price_out` | `setModelPriceLeaf` | 同上 | `:150-156`/`:169` |
| `role_chat_model` | `SetRoleChatModel` | **门③**——它的 `apply` 恒真（`:179-181`）、`check` 恒过（`:182`），所以它一路走到 `validate(candidate)`，撞 `validateCatalog`→`validateRoles` 的"model 有值 provider 空"分支 | `catalog.go:101-103` 逐字 `config.toml: llm.roles.chat.model set without llm.roles.chat.provider (set both or neither)` |
| `provider_credential` | `StoreCredential` | 门①+门②——它最后一步就是 `SetProviderAPIKeyRef`（`panel_config_store.go:267`），所以 blob 存成了、引用写不进（那枚半状态句 `:272`） | 同上 `:84-87`/`:92-94` |

⇒ 票面四环在本腿复认成立。**注册表为什么是 nil**：`applyDefaults` 的 `case reflect.Map:` 后一句 `// leave nil (see NewDefaults)`（`internal/config/defaults.go:77-78`），理由写在 `defaults.go:54-57`；并且就算文件里写一节空 `[llm.providers]`，`normalizeZero`（`defaults.go:148-152`）也把它归零回 nil（`// empty map -> nil, matching the absent-key decode state`）。首建文件＝`SaveFile(cfgPath, config.NewDefaults())`（`cmd/wisp/firstrun.go:82`），所以里面没有 provider 行。

## 2.3 要让它**能创建行**：最小改动面（逐枚，附代价）

1. **八枚分支**：四枚 setter 各两层——`settings.go:63-66`/`:71-73`、`:84-87`/`:92-94`、`:111-118`/`:124`、`:150-156`/`:169`（后两枚还要 `requireCatalogEntry` `:298-306` 的 model 半边）。"拒"改"建"要两层同改，不然门①放行、门②仍 `apply=false` ⇒ 走到 `:224-227` 那句"does not exist in the file"。
2. **nil map 会 panic（新踩的坑，票面/a1 都没点到）**：门②的写法是 `base.LLM.Providers[provider] = p`（`settings.go:68`/`:89`/`:121`/`:166`）。干净机上 `base.LLM.Providers` 恰是 **nil**（§2.2 末）——向 nil map 赋值是运行时 panic。⇒ 建行必须先补 map 初始化（`defaults.go:77-78` 那条"go-toml merges into pre-existing maps ⇒ 预填会泄漏幻影条目"的理由在此处反转：一旦写侧能建 map，"默认表里不许有值"与"落笔处要能建"就得由**谁负责 new**来决定）。model 级同理（`p.Models` 为 nil，`settings.go:115` 的 `p.Models[model]` 判读之前得先建）。
3. **门③要喂 protocol**：新建非预设名的行如果 `protocol` 为空，`validateAPIKeyRefs` 的 `else if` 分支直接拒（`internal/config/validate.go:171-177` 逐字 `config.toml: llm.providers.%s.protocol must be set explicitly for a non-preset provider (one of %s)`）。两条现成出口都有代价：**(甲) 只允许内置预设名**（`defaults.go:22-35` `providerPresets`＋`LookupPreset` `:48-52`，12 枚名册）——代价是名册外服务商永远建不了；**(乙) 一次写两枚键**（名＋protocol）——撞 `settings.go:25-28` 明写的"没有多键原子写，发明它是票 226 的地盘、要它自己的裁定"。**注**：`applyPresets`（`loader.go:188-204`）**不在写路径上**（它只在 `readConfigFile` 里，`loader.go:122`），所以"建了个预设名的空行"落盘时 protocol/base_url 仍是空串，要到下一次加载才被填满 ⇒ 回执 `writtenKeyPaths`（`settings.go:254-268`，它 diff 的是**文件前后**）会说"只改了 `llm.providers.<名>.base_url`"，而内存里那行的 protocol 是空、下一次启动后就不是空了。**这是一处会写进回执的说谎点，量得到形状、量不到具体串（禁 Go）。**
4. **既有 setter 无一有先例**：三枚更老的 setter 都是"换一枚叶子值/换一个 slice"，且**不经过 `writeOneKey`**（`permmode.go:70-79`、`allowdirs.go:139-142` 直接 `mergeWrite`）——全仓今天没有任何一处写路径往 map 里加过条目。⇒ "建行"是新形状，不是照抄邻行。
5. **面板面（C17 契约面，⛔ 本腿不决定）**：新字段要动 `config_handlers.go:57-69`（常量）＋ `:74-93` 三张表 ＋ `allWritableFields` `:90-93`；名册枚数被 `internal/panel/config_route_248_test.go:271-276` 以 `reflect.DeepEqual` **钉死为 7 枚**（现读：`want` 切片逐字列 `FieldModelContextWindow/FieldModelPriceIn/FieldModelPriceOut/FieldProviderAPIKeyRef/FieldProviderBaseURL/FieldProviderCredential/FieldRoleChatModel`）；新封套键被 `cmd/wisp/panel_config_248_test.go:369-371`（`TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope`，`:371` 逐字列 4 枚 `configField/configProvider/configModel/configValue`）钉住。⇒ 动这两枚＝动测试断言＝按 AGENTS.md §1.1 要人工批准。
6. **"行不存在必拒"这条断言本身**：`internal/config/settings_248_test.go:130-132`（用例名逐字 `"a provider the config does not declare"`，动作 `m.SetProviderBaseURL("inventco", "https://x.example.invalid")`）挂在 `TestAC7InvalidValueLeavesTheFileByteIdentical`（`:113`）下——**任何"行不存在则建"的形必打红它**（与 `257-a1` §0 补充末条同判，此处独立复认）。
7. **门②之后的内存/文件顺序没变**：`writeOneKey` 是"先内存（`:241`）→ 再文件（`:242`）→ 文件失败回滚内存（`:243-245`）"，建行也走这一条，不需要新发明；但 `mergeWrite` 的 base 是**文件**（`writeguard.go:140-142`），所以建行在文件侧是"往盘上那份 map 里加一枚"，与内存侧那份是两次独立的 `apply` 调用（同一闭包，`:241`/`:242`）——**中途失败的形状**：内存已建、文件没建（回滚路径覆盖），或 `writtenKeyPaths` 复读失败（`:261-267` 那句"写后无法复读文件"）。

**代价小结（不选形）**：真正的最小面是第 1、2、3 条三簇（同包、同一枚 `settings.go` 可做完，⛔ 不碰 C17、⛔ 不动枚数钉）；但**做完也只解锁 3/7**——`model_*` 三枚还要求 model 行（`requireCatalogEntry` 的 `:302-306`），`role_chat_model` 还要求 `roles.chat.provider` 被同时点名（`catalog.go:101-103`），而后两件事天然是多键写 ⇒ 撞第 3-(乙) 条。这一段与 `257-a1` §1-7 的"3/7"读数**独立复认一致**。

## 3. 两条通道的重读时机（热加载 vs 需重启）

（取数中）

## 4. 我可能写错的条目（自我对抗）

（取数中）

## 5. 量不到的地方（具名）

（取数中）

## 6. 交件判语

（取数中）
