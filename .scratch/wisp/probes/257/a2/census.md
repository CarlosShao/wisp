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
| 进场 HEAD | `3736f0dd`（2026-10-03 09:11:55 +0800，`dev`） |
| 进场 porcelain 行数 | 398 |
| 量程中 | `09:32:46` HEAD `bb39a595` porcelain 420（同机别腿在连续落 commit：`145-c2`/`248-c3`/`255-c2`/`258-a2`/`228-a4` 均在此间落件） |
| 本腿自身四枚 commit | 骨架＋§0：`395b3364`／§1：`fbb27d4b`／§2：`21f85d95`／§3：`3759ce6b` |
| 落件 | `09:34:23` HEAD `598cc093` porcelain 420 |

⚠ 全部行号现读于上列锚点区间内的**共享工作树**（未 checkout、未 branch）；`internal/config`／`cmd/wisp` 正被 255/257/248 系写腿改，**本文的行号与"枚数"仅对这一区间负责**（§4-N1/§6 末条已具名）。§4/§5/§6 与 §1.6 属 09:28–09:34 间的补读与自我对抗，未回头改 §1–§3 的既有读数（只在 §3.2 就地打过一处自我纠正，见该节表格）。

## 1. 凭据通道名册（谁解析 / 何时读 / 失败句子逐字）

**通道名册（本问的答案，详表在下面各节）**：甲 `dpapi:<blob 名>`（§1.1，真读值）｜乙 `env:<环境变量名>`（§1.2，真读值）｜丙 `config.toml` 里的引用形（§1.3，只是"指名"；其解引用腿 `resolveRefs` 在生产里恒 `res=nil`）｜丁 别针五枚（§1.4：`wisp secret get --show` 是全仓唯一批准的明文出口、`list`、`unset`、`wisp providers`、`cmd/llmrecord` 的裸 env 读）｜戊 明文 `api_key` 迁移路（§1.6，**已建、今天没有生产调用者**）。**能真读到凭据值的只有甲、乙两条**，且都收敛到同一枚 `Store.Resolve`（`internal/secret/store.go:89-123`）。

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

### 1.6 通道戊：明文 `api_key` 的迁移路（**已建、今天没接**）——回执最容易被它绊到的一处

- **有一枚读明文的东西存在**：`secret.MigratePlaintext`（`internal/secret/migrate.go:86+`，词汇表在 `:21-34`：`PlaintextKey = "api_key"`／`RefKey = "api_key_ref"`／`BackupSuffix = ".bak-plaintext"`），把树里每一枚非空字符串 `api_key` 存成 DPAPI blob（**blob 名由 TOML 路径确定性推出**，`migrate.go:73-74` 注释），改写成 `api_key_ref = "dpapi:<id>"`，并把原件备份成 `config.toml.bak-plaintext`（`:76-78` 说明备份被 `winsec.PrivateFile` 收窄，因为**那个备份文件里就是明文**）。用户可见句在 `MigrationReport.Notice()`（`migrate.go:50-56`），逐字：`config security (D33): %d plaintext api_key field(s) in %s were migrated to DPAPI-protected references and removed from the file; the original was backed up to %s`。
- **⚠ 它在生产里没有调用者**（现读：`grep MigratePlaintext` 的非定义命中只有 `internal/config/migrate.go:85` 的一句注释、`internal/secret/configrefs.go:16` 的一句注释，以及 `internal/winsec/migrate_windows_test.go:82`/`:133` 两枚测试）。⇒ **D33 的这枚"明文→DPAPI"迁移今天是一台没通电的机器**；`internal/secret/configrefs.go` 那套引用索引（`ConfigRefs`/`RefFieldNames`）也只为 `wisp secret unset` 服务（`configrefs.go:18-20`）。
- **另一处确实会留下明文的地方，与它无关但更近**：schema v1→v2 的迁移备份。`internal/config/migrate.go:82-95` 逐字注释说得很白——备份是**迁移前配置的原文**，而"在还没走过 D33 的那台机器上，那是这里唯一还带着 plaintext api_key 的文件"，所以这枚备份也经 `winsec.PrivateFile` 密封（`:93`），并且它**是修复旧宽权限备份的路径**（`:89-92`）。
- **对 257 回执的直接影响（这是本腿要给落地腿的那句）**：`cmd/wisp/firstrun.go:114` 现在写的是"这份文件里没有写明文 key 的字段，要补的是那个名字"——**这句对 schema 成立**（`Provider` 结构体没有 `api_key` 字段，写 `api_key = "..."` 会被 `decodeStrict` 的 `DisallowUnknownFields`（`internal/config/parse.go:72`）打回 `unknown key "..." at line N`（`parse.go:145-148`），并由 `describeReloadFailure` 归到 `cause=unknown-key`，`config_reload.go:334-337`）。**但它在三件事上不完整**：① 旧机器上的 `<dataDir>\config.toml.bak-1`（`migrate.go:82` 的后缀形）里可能仍有明文，而**没有任何东西会把它迁走**（上面那条）；② 若机主把 key 直接写进 v2 的文件，他看到的不是"请改用 `wisp secret set`"而是 `unknown key` 那一类归因——运行中改的是 `cause=unknown-key` 那句（`config_reload.go:334-337`），**启动时改的是整条装配死**：`wisp run: 配置未就绪（Unconfigured）：%v`（`run.go:417`）套那句 `config.toml: unknown key ... at line N`；③ **带着明文 `api_key` 的 v1 文件今天根本迁不动**：v1 的 provider 表本来就在 `llm.providers.<n>` 下（现读：`migrateV1toV2` 取的是 `llm["providers"]`，`internal/config/migrate.go:127`），迁移只搬 `model`→`models.<id>`、`timeout`→`timeout_ms`、`default_provider`→`text_chain`/`roles.chat`（`:102-109` 自陈），**明文 `api_key` 不在被搬的名单里 ⇒ 原样留下**（`:111-112`"everything else passes through unchanged"）；迁完的 dry-check 走 `decodeStrict`（`:70`）⇒ `api_key` 对 v2 schema 是 unknown key ⇒ 拒句在 `:71-73`（`config.toml: migration to schema version %d produced an invalid config (%v); the file was left untouched`）。⇒ 也就是说 **D33 的"把明文收进 DPAPI"这一环，从 v1 迁移到首建文件，两条路都没有落地的执行者**；回执若承诺"写了明文我们会帮你搬走"，今天就是谎。（⚠ ③ 的整条链是**静态读码**得出的：`migrate.go` 的三步顺序 `:70`→`:77`→`:93`→`:96` 现读，但"这份 v1 文件真的停在原样"要跑起来才量得到 → §5 N5。）

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

**一句话结论（现读）**：全仓只有**一枚** tick 会重读盘，它挂在 `wisp run`／常驻任务腿自己的那枚 Manager 上；设置页写入腿用的**往往是另一枚 Manager**（不 tick）。写盘之后**没有任何事件**发给持有旧值的组件——一件都没有，连"给自己的 tick 报警"都被 `statOwnWrite` 主动压掉了。所以"手改能被看到"与"面板写要重启"两句**都真**，但它们不是同一件事的两种说法，而是**四件事**：① 盘被重写 ② 写它的那枚内存被更新 ③ 别的 Manager 的 1s 轮询把新值拉进自己内存 ④ 本进程真正干活的那条模型通路**只在装配时建一次**，谁都更新不了它。

### 3.1 谁在什么时候重读盘（唯一一枚）

| 步 | file:line | 要点／逐字 |
|---|---|---|
| 装配处接线 | `cmd/wisp/config_reload.go:114-115` | `rt.mgr.ConfirmLocked = rt.confirmLockedLoosening`／`rt.mgr.OnRestartPending = rt.reportRestartPending` |
| 起协程 | `config_reload.go:119` | `observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)`；形状自陈在 `:23-33`（"常驻 tick，不是文件通知"） |
| tick 体 | `config_reload.go:131-147` | `time.NewTicker(configReloadPollInterval)`，`configReloadPollInterval = time.Second`（`:88`） |
| 一次轮询 | `config_reload.go:152-162` | `rep, err := rt.mgr.CheckAndReload()` |
| **谁被接线** | `cmd/wisp/run.go:813`（`rt.startConfigReload()`，在 `assembleRuntime` 体内） | ⛔ 除此之外全仓无第二枚生产调用：`grep startConfigReload` 命中＝定义 `config_reload.go:105` ＋ `run.go:813` ＋测试 |
| 重读的是什么 | `internal/config/manager.go:149-160` | `os.Stat` → mtime+size 与 `seen*` 相同即 `return nil, nil`（`:156-159`）；不同才 `LoadFile(m.path, m.res)`（`:160`）——**整条管线重跑、含 `resolveRefs`**（`loader.go:46`），但因 §1.3 那处 `res==nil`，实际一枚凭据都不解 |
| 重读后落到哪 | `manager.go:169` `plan(fresh)` → `:192-196` `plan.commit()`＋`m.resolved = resolved` | `[llm]` 在 hot 名册（`manager.go:284`），改了就整段覆盖 `cur.LLM`；尺在 `tiers.go:30`（`"llm": "hot"`），名册与尺对不上会 `panic`（`manager.go:294-300`） |

⇒ **"手改配置文件能被热加载看到"这一句今天可证实**，但它的确切射程只是：`wisp run`（控制台）与常驻进程的**任务腿那枚 Manager** 的内存与 `Report`。**不**覆盖面板读写腿那枚（§3.2）。

### 3.2 一台机器上到底有几枚 Manager（owner 那句"我填了为什么没生效"的真身）

三枚生产装配点，各自独立持有一份 `cur`：

| 装配点 | Manager | 会不会重读盘 | 谁在用它 |
|---|---|---|---|
| `cmd/wisp/run.go:409` `config.NewManager(cfgPath, nil)` | `rt.mgr` | **会**（`run.go:813` 起了 tick） | `rt.settings`（`run.go:431`）→ 只被快照的凭据维度用（`run.go:710` `Credential: rt.settings.credentialStatus`）；`rt.permissionMode()`（`cmd/wisp/approval_always.go:151-156`，每次调用现读 `rt.mgr.Config()`）；`perm.Store.PermissionMode()`（`internal/perm/store.go:155-164`，经 `internal/tools/mode.go:35-47` 每次工具调用拉一次） |
| `cmd/wisp/panel_inbound.go:230`（同 `newComposerDispatchChain`，定义 `:228`） | 链上的 `mgr` | **不会**，且**明说自己不会**：`hotReloadDisabledPanelInbound` 逐字在 `config_reload.go:97-99`（"本宿主没有接 config.toml 轮询，这次运行期间手改配置不会生效……要生效请重启进程并用 wisp run"），由 `panel_inbound.go:218` 每启动播一次 | `config.set`／`config.get` 的全部写入与读数（`panel_inbound.go:266-267`） |
| `cmd/wisp/panel_resident_windows.go:170` → 同一枚 `newComposerDispatchChain` → 同一个 `panel_inbound.go:230` | 常驻面板链上的第二枚 | **不会**——`residentPanelHotReloadNote` 逐字在 `panel_resident_windows.go:162-165`（"panel host (resident): this leg does not tick config.toml either; the reload tick lives in `wisp run`…"），由 `:169` 播一次 | 同上 |

⇒ **常驻面板链那枚 Manager 永远没有 tick**（`panel_resident_windows.go:170` → `panel_inbound.go:230`，`residentPanelHotReloadNote` 已自陈）。而**同一进程里是否存在那枚会 tick 的 Manager，是有闸门的**：`assembleRuntime`（含 `run.go:813` 起 tick）在常驻进程里只由任务腿调（`resident_task_source_windows.go:265`），而它在 `:224-231` 有一道控制台闸——`interactiveStdin()`（实现 `cmd/wisp/approval_reply_stdin_windows.go:41-50`：stdin 句柄拿不到、或 `GetConsoleMode` 拒 ⇒ 返回 nil；非 Windows 垫片 `cmd/wisp/approval_reply_stdin_other.go:23` 恒 nil）且无测试注入文本时**直接 `return nil`**，逐字句：`wisp: %s：本进程仍然带球常驻、仍然会在退出时拒绝挂起的卡片，但没有任何东西会去举一张卡`。
⇒ 三种进程形状，回执必须分清（⚠ 第三行的"这台机器上到底是哪一种"属 §5 N4，只有跑起来量得到）：

| 形状 | 有几枚 Manager | 谁重读盘 | 后果 |
|---|---|---|---|
| `wisp run`（终端起，控制台在） | 1（`run.go:409`） | 那枚，1s tick | 手改 1 秒进内存；面板读数与写入腿**不在**这个进程里（`config.set` 只有 `wisp panel-inbound` 才走到） |
| `wisp panel-inbound`（一次性 CLI 缝） | 1（`panel_inbound.go:230`） | **没有** | 手改本次运行永不生效（`config_reload.go:97-99` 那句就是为此而播） |
| `wisp resident`（真日常面） | ≥1；**只有任务腿过了控制台闸才有第二枚** | 面板链那枚永不；第二枚存在才 tick | 面板那页的读数与写入腿走的是不 tick 的那枚；若这台机器上进的是无控制台形状，则**全进程没有任何东西重读盘**——手改与面板写在两个方向上都只能靠重启 |

- **今天可证（静态）**：面板 `config.get` 的读数（`ReadSettings` → `panel_config_store.go:92` `s.mgr.Config()` → `manager.go:115-119` 内存深拷贝）**只反映面板链那枚内存**，手改文件对它**在这枚进程活着的期间永远不可见**；而快照的凭据维度（`run.go:710`）走 `rt.mgr`，**1 秒内可见**——**前提是那个 runtime 被装配出来**（同上闸门）。两句同时成立的那个瞬间是否存在，见 §5 N4。可证的静态结论是：**两把尺读的可以是两枚不同步的内存**（例：手加一节 provider 后凭据维度可由 `no_ref_declared` 变为 `partly_missing`/`all_recorded`，而 `config.get` 那句仍是 `服务商 0 家；凭据状态 no_ref_declared`——句在 `internal/panel/config_handlers.go:451`）。
- **只有跑起来量得到**：这两句在真常驻进程里是否会被**同一秒**读到（取决于谁触发快照，§3.3 末行）。→ §5 N3。
- ⚠ 反过来的一条也要写进回执：**写门看盘、读面看内存**（§2.1 末）。所以"手加一节之后那七枚写不写得进"与"那一节显不显示在页面上"是两回事，前者立刻、后者在这枚进程里不。

### 3.3 写盘之后有没有事件发给持有旧值的组件：**逐枚否证**

| 候选事件 | 现读结果 |
|---|---|
| setter 直接通知谁 | **没有**。`writeOneKey`（`settings.go:199-249`）的出口只有 `check` 错／`apply` 假／`mergeWrite` 错或成；`grep "OnReload\|OnRestartPending"` 在非测试代码里的**赋值**只有 `config_reload.go:114-115` 两枚 ＋ `cmd/balldebug/main.go:244`（另一枚命令，不是 wisp run）；两枚钩子的**调用**只在 `manager.go:198-203`——即**只有 tick 会发**，写路径永远不发 |
| 写完让本 Manager 自己的 tick 再读一遍 | **主动压制**：`mergeWrite` 成功且文件与内存一致时 `m.statOwnWrite()`（`writeguard.go:175` → `allowdirs.go:157-163` adopt mtime+size）——写它的这枚 Manager 永远**不会因自己这次写而 reload**。只有文件里带着本机没读过的他人手改时才**不** adopt（`writeguard.go:178-186` 那支 `slog.Warn`），下一 tick 才读回来 |
| 让**别**的持有者知道 | **没有任何跨 Manager 机制**（无注册表、无 channel、无回调；三枚各自 `NewManager`）。唯一的"传播"就是别的 tick 各自 poll 到 mtime 变化 |
| 让装配期就快照走的组件知道 | 不可知，且代码自陈：`run.go:423-424` `cfg := mgr.Config()` 是深拷贝快照并 `rt.cfg = cfg`；模型通路由它一次建好——`run.go:435` `llm.NewResolver(cfg, st)`＋`:436` `ResolveRole`＋`:442` `BuildEndpointProvider`＋`:447` `rt.provs`；凭据也在同一步被烤进适配器（§1.1）。同一形状在 `config_reload.go:55-61` 写得很白：C26 canonicalizer 只在装配时建一次（`run.go:390`），所以**已批准的 [fs] 放宽也不改本次运行的路径判定**，并专门为此播一句（`:201-203`） |
| 让**页面**知道 | 页面只被**不相关的状态移动**推：`cmd/wisp/panel_pump.go:392-396` 注释逐字 "Called where a state the panel shows actually moved: a card opened, a card went away, a tool call started or ended, the stream closed."——**配置写成功不在这张单上**；`rt.pump.Publish()` 生产调用面只 `panel_pump.go:405` 一枚，函数值经 `run.go:735`／`run.go:998` 交出去。⇒ "写完不刷页面就看不到新数"是**今天可证的形状** |

### 3.4 两条通道今天各自说到用户嘴上的句子（逐字）

**通道甲（手改文件 → tick）**，`cmd/wisp/config_reload.go`：
- 装 tick 时（`:124-127`）：`wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。手改会按 D36 三档处理：可热加载段立即生效；[risk]/[fs]/[net]/[plugins] 的放宽要先答一张 L2 卡，不答按拒绝保留旧值；重启档的改动本次不生效，会另有一句告诉你为什么不生效。`
- 段被立即生效时（`:171`）：`wisp run: 配置热加载：这些段已立即生效（D36 立即档）：%v` ← **`[llm]` 的改动会落进这一句**
- reload 失败归因（`describeReloadFailure` `:315-370`）：`cause=missing`／`cause=permission`／`cause=syntax`／`cause=unknown-key`／`cause=migration`／`cause=invalid`／`cause=unclassified`（票 223 AC#4 的四分法今天已是七分）
- 不 tick 的宿主各播一句：`config_reload.go:97-99`（panel-inbound）／`panel_resident_windows.go:162-165`（常驻面板腿）

**通道乙（面板写 → 回执）**，`internal/panel/config_handlers.go`＋`cmd/wisp/panel_config_store.go`：
- 级别句（`tierSentence` `:428-441`）——写侧全部返回 `EffectiveRestart`（`panel_config_store.go:228`／`:275`），所以页面拿到的必然是 `:435` 那句：`这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。`
- 受理句（`renderSettingReceipt` `:381-396`）：`设置已写入。 实际落盘的键路径：<diff 出的键>。 这一项要重启…`；文件本来就是这个值时：`本次没有任何键路径发生变化（文件里已经是这个值）。`（`:385`）
- 凭据句（`renderCredentialReceipt` `:401-423`）：`密钥已录入（只写不回显：这一页任何时候都不会再显示它的值，只会显示已录入/未录入）。 服务商 <名>。 记录的是引用 <ref>。 …`
- 拒写三因（票 257 AC#2 要的就是这三句别折成一句）：**① 文件没建** → 面板链根本装配不出来，句在 `panel_inbound.go:148`（`wisp panel-inbound: 入向装配未完成（Unconfigured）：%v`）套 `panel_inbound.go:232`（`配置未就绪（%s）：%w`）套 `manager.go:102`→`loader.go:67`（`config.toml read`）——**这一因今天没有 Wisp 自己的中文句子**（那是票 198 首建腿的地盘）；**② 行不存在** → `settings.go:311-313` 或 `:303-305`，**英文**；**③ 校验不过** → `settings.go:231-232` 包着 `validate.go` 的英文句。

### 3.5 那么回执要说什么才不说谎（给 257-r1 备料；⛔ 不新建第二台回执机器）

**先点名已有那台**：三段式首启回执已在库里——`cmd/wisp/firstrun.go:92-95`（建件句）＋`:111-115`（key 句）＋`:116-118`（模型句），且**已点到真入口名**：`wisp secret set <blob 名>`、`--from-stdin`、`api_key_ref = "dpapi:<blob 名>"`、`env:<环境变量名>`、`[llm.providers.<名>]`、`text_chain`/`roles.chat`、`wisp providers discover`/`probe`。⇒ **本腿判读：257 的增量只应是这三段的措辞，不应出现第二个打印点。**

它今天缺的、且 §3 已用现读证实的五件事（措辞要点，逐条指回上面）：

1. **两条通道要分名，不能合成一句**（票 257 §8 边界①）。今天可证的分裂：手改 → `config_reload.go:171` 会说"已立即生效（D36 立即档）：[llm]"；面板写 → `config_handlers.go:435` 会说"要重启"。**两句都对，说的却不是同一样东西**：前者是那枚 tick 着的 Manager 的**内存**，后者是 `run.go:435-447` 那条**装配一次**的模型通路。
2. **"面板写要重启"的理由要说到底是两重的**：除"模型通路不重建"外，还有 §3.2 那枚**不 tick 的 Manager**。**但两重是否同时落在用户那台机器上，取决于 §3.2 那张表**：`wisp panel-inbound`（一次性进程）只有第一重；常驻若无控制台闸过不了，则**两重都不是"要不要重启"的问题，而是整个进程根本没有重读盘的东西**——那种形状下"重启"是**唯一**生效路径，对手改也成立。回执若照 §8 边界① 只写两格（手改=热加载／面板写=重启），在无 tick 形状下就是错的。
3. **反向的诚实也要给**：手改对**面板那一页的读数**同样不生效（§3.2/§2.1）。若回执只写"手改＝热加载认、面板写＝要重启"，仍是半谎——正确形状是"手改会被**运行中的 `wisp run`** 在 1 秒内读进它自己的内存；这一页上的读数不跟着它变；两处的通路都要重启才算用上"。现成邻居句可指：`config_reload.go:97-99`／`panel_resident_windows.go:162-165`。
4. **`[llm]` 是 hot 段这条规格话今天仍在码里**（`schema.go:403-405` 注释逐字 "entirely hot-tier (api_key_ref changes trigger re-resolution of the reference, never a restart)"; `tiers.go:30` `"llm": "hot"`）——它与 `panel_config_store.go:25-32` 那段自陈（"every accepted write answers restart required, even though the spec's table calls [llm] hot"）**并存且已知不一致**；而 §1.3 的"`res` 恒为 nil"更进一步：**api_key_ref 的重新解引用今天在生产里根本不发生**（无 resolver）或**只在新建 Endpoint 时发生**（`resolver.go:141`）。回执不应引用那句注释。⚠ 这属"注释与规格文字"层面，改它需人工批准，本腿只具名不提议。
5. **三因句子的语言不一致**：②③ 是英文（`settings.go:303-305`/`:311-313`/`:231-232`），①只有装配错的英文＋中文混合句（§3.4 末）；面板侧唯一成体系的中文拒句是名册/选择器那六道（`config_handlers.go:310-353`）。票 257 AC#2 要"三种拒因各配一句不同的话"——今天**只有第②与第③在结构上分得开**（门①/门② vs 门③），第①在面板链上是"根本没有这条腿"，与②③不在同一层。

## 4. 我可能写错的条目（自我对抗）

- **N1｜§1.3"`res` 恒为 nil、`resolveRefs` 是死线"**——我核的是三枚生产装配点（`run.go:409`／`panel_inbound.go:230`／`cmd/balldebug/main.go:231`）加两枚 `LoadFile(..., nil)`（`providers.go:98`／`models.go:184`），用的是 `grep config\.NewManager|config\.LoadFile\(`。⛔ 我**没**排掉"把 resolver 装进变量再传"的写法（那要靠读 `Manager.res` 的赋值点，只有 `manager.go:106` 一枚，我读过，但没做全仓 `res:` 字段写入扫描）。另：`internal/config` 正被 255/257 的写腿改，**我这句只在下面 §0 的锚点区间内成立**。
- **N2｜§1.1"凭据在装配时被烤进适配器，此后每次请求不再碰 store"**——依据是 `EndpointOptions.APIKey` 是 `string`（`internal/llm/provider.go:157-167`）、`NewProvider` 只按 protocol 取工厂（`provider.go:217-232`）、三枚适配器都从 `a.opts.APIKey` 取头。**风险点**：我没逐行读重试/限流包装层（`internal/llm/retry*.go`、`ratelimit.go`）确证它们不持 `KeyResolver`；缓解证据＝`grep KeyResolver` 非测试命中**全在** `internal/llm/resolver.go`（`:30/:32/:86/:90`）一处文件内 ⇒ 结构上不可能有第二处解引用者，但这是"以grep 代精读"，标出。
- **N3｜§2.1 标题句"写门读盘不读内存"**——**只在文件存在时成立**：`settings.go:213-218` 明确写了文件缺失时 `candidate` 退回 `m.cur`/`NewDefaults()`。而干净机上面板链在文件缺失时**根本装配不出来**（`panel_inbound.go:230`→`:232`），所以那半分支今天从面板这条路走不到；从 `wisp run` 的 `perm`/`allowdirs` 那两枚旧 setter 能走到，但它们不经 `writeOneKey`。**净判读不变，但别把它读成无条件**。
- **N4｜§2.3-2"向 nil map 赋值会 panic"**——这是对**假想改动**的代价，不是今天的缺陷：今天的代码在赋值前先 `if !ok { return false }`（`settings.go:63-66` 等），所以永不触发。证据链我核了两条（`applyDefaults` 的 map 分支留 nil，`defaults.go:77-78`；`copyValue` 对 nil map 直接 `return` 不分配，`defaults.go:213-216`），但**我没有真跑过**（禁 Go），所以"落地腿若不初始化 map 就会 panic"这一句是 Go 语义常识＋现读的两条 nil 证据，不是量到的栈。
- **N5｜§3.2 那张三形状表**——中间那列（常驻是否有第二枚 Manager）我核到了 `resident_task_source_windows.go:224-231` 的控制台闸与 `assembleRuntime` 的全仓仅有的两枚生产调用（`run.go:249`、`resident_task_source_windows.go:265`）。⚠ 但**"用户那台机器上进的是哪一格"我量不到**（要 `GetConsoleMode` 的真实答案）→ §5 N4。另：我**没有**审 `wisp resident` 的启动方式（快捷方式/计划任务/自启是否附控制台），那属票 33/228 的地盘。
- **N6｜§3.3"没有任何事件发给持有旧值的组件"**——我核的是两枚钩子的赋值面与调用面、`statOwnWrite`、以及 `SnapshotPump.Publish` 的生产调用面（`panel_pump.go:405` 一枚，函数值经 `run.go:735`/`:998` 交出）。⚠ 残余风险：`rt.ui.publish` 与 `consoleSink.publish` 的**触发点集合**我是靠 `panel_pump.go:392-396` 那句注释＋`run.go:1297-1298`/`:1329-1332` 的字段注释读出来的，**没有枚举行级**——若某条我漏看的路径在配置写后调了 publish，"页面不刷新"这半句要收回。
- **N7｜§1.6 的三处"不完整"**——① 与 ② 是纯现读（无调用者＋拒句都在码里）。③（v1 明文迁不动）是**静态链**：`migrate.go:127` 的取值路径＋`:70` 的 dry-check＋`parse.go:72` 的 strict 三处现读，但"这份文件真的被原样留下、退码与句子长什么样"只有跑起来才量得到（且 `MigratePlaintext` 无调用者这一条本身就是"这台机器有没有走过 D33"的历史问题）→ §5 N5。
- **N8｜枚数类读数**——`providerPresets` 我数了 12 枚（`defaults.go:23-34` 逐行）；C17 方法名 6 枚（`bridge.go:41-68`，`bridge.go:148` 的 case 列了同名 6 枚，与 `257-a1` §2 一致）；可写字段 7 枚（`config_handlers.go:57-69`/`:90-93`，被 `config_route_248_test.go:271-276` 钉死）。这三处都是**可复核的枚数**，不是推算。
- **N9｜我对 `257-a1` 的两处补强不是推翻**：a1 §0 四环全部复认成立（含它自己修正过的 `catalog.go:101-103`）；我在 §2.1 补的是"那道门读的是盘"，在 §2.3-2 补的是"nil map"，在 §3.2 补的是"面板链那枚 Manager 不 tick"。三处都**不改** a1 与票面 257 §8 的任何判语。

## 5. 量不到的地方（具名）

| # | 想量什么 | 为什么量不到 | 需要什么才能量到 |
|---|---|---|---|
| N1 | 票 257 AC#1 那一发：**干净机＋无 `config.toml` 起步、首建之后逐枚试写那七枚**的终态与逐字回执 | ⛔ 本腿禁跑任何 `go` 命令（三枚写腿在飞，同机会互相洗读数） | `go test`／真跑 `wisp run`＋`wisp panel-inbound` 喂七发封套（＝落地腿 257-r1 的 AC#1） |
| N2 | 门①②③各自**今天的完整拒句串**（含 `%q` 填上的实值）与页面渲染后的整段文字 | 同上；我只抄了**格式串**与拼装函数（`config_handlers.go:381-423`），没跑过一次渲染 | 一发真跑的 stdout/回执照 |
| N3 | "手改→1s 进内存"与"面板写→回执说重启"**在真进程里是否同秒可读**、以及 `hot=`/`restart=` 那行的实际内容 | 要跑：`config_reload.go:168` 那条审计行只在 tick 真命中时出 | e2e 腿（可参照 `.scratch/wisp/probes/e2e-panel-1/readiness.md` 的跳表法） |
| N4 | 用户那台常驻机**落在 §3.2 三形状的哪一格**（任务腿有没有过控制台闸） | 是运行期 `GetConsoleMode` 的事实＋启动方式，码里读不出 | 真机起一次 `wisp resident`，看有没有打出 `taskEntryDisabledClaim`／`taskEntryConsoleClaim` |
| N5 | §1.6-③ 那份"带明文 `api_key` 的 v1 文件"迁不动时**退给用户的整句与文件是否逐字节原样** | 静态链已给（`migrate.go:70-73`＋`parse.go:72`），但要跑一发才有串 | 种一份 v1＋明文键的文件，读退码与句 |
| N6 | 页面侧要什么、缺什么、有没有"重问一次"的按钮 | ⛔ `frontend/**`／`design/**` 两层禁令，未读未引 | 由 owner 自己带给他用的那枚前端 agent（本腿不代转、不猜页面里的字） |
| N7 | 便携档 `ErrPortableDecrypt` 会不会真的在用户的便携盘上出现（§1.1/§1.5） | 要跑＋真机 DPAPI；且 `WithPortable` 只在 `wisp secret` 一条腿上传 | 便携安装真跑一发 `wisp run`（另一枚待定案项：AGENTS.md §2 列的"便携×DPAPI"） |
| N8 | `d22scan` 等仪器对本次读码点是否报警 | ⛔ 不跑 `go vet`/工具 | CI（本腿只做读码） |

## 6. 交件判语

- **只读**：本腿全程只 `Read`/`Grep`/`git log|show|diff`/`sed -n`/`grep -n` 读码，**未写任何产码或测试文件**；写入面只有 `.scratch/wisp/probes/257/a2/**`（`census.md` ＋ `msg1..msg5.txt` 五枚临时件，**只建不删**）。
- **零 Go 命令**：未执行 `go test`／`go build`／`go vet`／`go run`／`go list` 中的任何一枚（第一硬规）。§5 那八格全部具名留在"量不到"，未以推测填空。
- **未 push**：只 commit；提交面逐枚带显式 pathspec，无 `add -A`/`.`/`-a`，无 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`restore`/`clean`/merge/worktree。
- **AC 框未碰**：`.scratch/wisp/issues/257-*.md` 只读；票面与四个 AC 框一字未改（`git status` 里那枚票文件不属本腿改动）。
- **冻结件未碰**：`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`／`.github/workflows/ci.yml` 一律未写。本文中引用它们的只有**射程判断**（§2.3 第 5/6 条引的是 248 系测试的断言形状，出处是**读到的测试文件**，属"读内部判射程"，非内容引用）。
- **`frontend/**`／`design/**` 两层禁令**：未读、未在证据件或回报正文中引用或转述（§5 N6 明确拒答页面侧）。
- **凭据零值外泄**：全文只出现**格式串**（含 `%q`/`%s` 占位）、字段名（`api_key_ref`/`api_key`）、blob 名占位（`<blob 名>`/`dpapi:<id>`）、env 名（`WISP_ENV`/`WISP_MINISIGN_KEY`/`USERNAME` 一类非密钥名）与包内符号名；**没有抄录任何一枚凭据值**（含测试文件里的假值，凡带值形状的串我一律未抄）。
- **三问均已答**（§1/§2/§3），§0 锚齐，骨架期的取数占位符已全部替换为本体读数（`grep` 全文占位形＝0 命中）。
- **对 `257-a1` 的复认**：★ 它那处推翻票面拼法的读数**复认成立**——真实路径是 `[llm.providers.<名>]`，顶格 `[providers.x]` 会撞 `DisallowUnknownFields`（现读：`internal/config/parse.go:72`，另 `:123` 是 `[plugins.<id>]` 条目那把同款尺；拒句形状 `unknown key %q at line %d` 在 `parse.go:145-148`；后果链＝`readConfigFile`→`decodeStrict` 失败→`NewManager` 失败→`panel_inbound.go:230-232`「配置未就绪」/`run.go:417` 退 2）。**未发现需要推翻 a1 或推翻票面 257 §8 判语的新证据**；a1 未覆盖、本腿补的三处在 §4-N9 列了。
