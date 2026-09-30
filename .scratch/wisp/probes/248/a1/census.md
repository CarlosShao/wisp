# 票 248 · 只读普查件 `248-a1` — 「设置」那条路的三问 + C17 文档漂移一问

> 腿：`248-a1`（只读普查） · 票面：`.scratch/wisp/issues/248-...md` AC#0 · 归口：面板
> 起手锚点：`2b4e2602`（分支 dev）· 起手钟：`2026-09-30 23:41:16 +08`（`date` stdout）
> 写点只有本文件一处；⛔ 本腿未改任何产码／测试／配置／`go.mod`／`go.sum`／票面。
> ⛔ 本件不含任何凭据值，只写字段名与变量名。
> ⚠ 本件里所有对 `cmd/wisp` 的引用一律是〔起手读数，可能被票 33-r1 改动〕。

状态：**①-④ 已填；⑤（42 把尺读数）／⑥（12 条自我推翻）／⑦（J1-J9，甲／乙／不做）先于结论表落盘。**
本腿只写本文件一处，未碰产码／测试／配置／`go.mod`／`go.sum`／票面的 `- [ ]` 框；⛔ 未跑 `go build`／`go vet`／`go test`／任何 `./...`。

## ① 要新增哪几枚方法名：参数形状 / 错误形状 / 生效级别 / 重启档 / posture 敏感字段

### ①.1 最小集＝两枚（读一枚、写一枚），凭据若要过境是第三枚

候选名（⛔ **本腿不命名，命名是 `A484` 的射程**，见⑦ J6）：`panel.config.read` + `panel.config.write`；
若⑦ J4 选甲，再加一枚**只写不回显**的凭据枚（如 `panel.secret.set`）。三枚都要落进同一处四件套：

| 位置 | 现读 | 扩一枚要动什么 |
|---|---|---|
| 常量声明 | `internal/panel/bridge.go:41-46`（四枚 const） | 加常量；⚠ `bridge.go:32-40` 那段注释自陈「名册改一枚 ⇒ `composer_test.go:502` 那枚钉会红」 |
| 守卫 | `bridge.go:104-110` `knownComposerMethod` 的 case | **动这一行会撞冻结件**（⑦ J1：`l2_grant_boundary_test.go:2051` 逐字钉着那行 case） |
| 封套 | `bridge.go:69-79` `ComposerRequest`（所有已答方法**共享同一枚封套**） | 新参数字段只能进这枚共享封套 ⇒⑦ J4 |
| 派发 | `composer_dispatch.go:137-162` switch + `:67-88` 四枚 handler 接口 + `:97-111` 结构体字段 | 新增接口 + 新增字段 + 新增 case；`:27-35` 写明 default 支的存在理由就是「bridge 扩了而这文件没被教」 |

**参数形状**（照既有四枚的宽度写，不新造通道）：读侧无参（或可选 `section` 过滤），错误/回执按
`RefusedEnvelopeForUser`（`bridge.go:115-125`）那句「被拒绝 [method requestId]」形；写侧一枚键路径 + 一枚值。
现成的键路径写法有唯一口径：`writeguard.go:69-72`（`fs.allowed_dirs` / `risk.permission_mode`）、
`:199`（map 元素按 `[llm.providers.<name>]` / `[plugins.<id>]` 具名）——⚠ 新写者必须**沿用这套拼法**，
否则同一枚键会在日志与 D36 裁决里被叫成两个名字（`writeguard.go:66-68` 与 `unwired.go:44-48` 两处都专为这件事写）。

**错误形状**（既有词汇，能直接复用，别造第二套）：
- 封套不合格 → `ErrComposerRequest`（`bridge.go:49`，每条都带 reason 与 requestId，`:84-102`）；
- 在册而本机没接处理器 → `ErrNoHandlerAttached`（`composer_dispatch.go:62`、`:182-187`，AC#1 要的那声响）；
- 在册而派发表没列出 → `ErrRosterMismatch`（`composer_dispatch.go:171-178`）；
- 写腿未接 → 既有 `ErrNoModeWriter` 形（`composer_handlers.go:60`）；
- **放宽而无确认腿** → `ErrNoL2Confirm`（`composer_handlers.go:57`、`:133-138`）；
- 当前值不可读所以无法判方向 → `ErrNoCurrentMode`（`composer_handlers.go:66`，注释写明「报一个它其实没在跑的宽档」会把闸门整个关掉）；
- 跨边界错误分类 → `observe.ClassConfig` + `observe.Wrap`（`permmode.go:72,83-85`、`writeguard.go:120-122,129-130`），这是 D37 的形状。

### ①.2 读侧要能答的字段：`internal/config/schema.go` 里**今天真存在**的模型相关字段（逐枚带行号与默认值）

> D36 规则 3（`schema.go:8-10`）：**默认值只活在 `default:"…"` tag 里**，TOML 文件不是第二个默认源。
> ⇒ 下表「无 tag」的字段，其零值是**「没填／未知」**，不是默认值；读侧必须像 `composer.go:256` 的
> `modelKnown` 那样**另给一枚 known 位**，否则「未配」与「配了个 0」在页面上长成一个样。

**`[llm]` 段**（`schema.go:406-420`）
| 字段（行） | tag/默认 |
|---|---|
| `LLMSection.TextChain :411` | 无默认（有序兜底链，元素 `provider/model`） |
| `LLMSection.Roles :413` → `Role{Provider :286, Model :287, Temperature :289 default:"0.7", ThinkingIntensity :291 default:"off", MaxOutputTokens :293}` | 四枚固定角色 `Roles{chat,memory_extract,handoff,summarize} :297-302` |
| `LLMSection.TimeoutMS :415 default:"60000"` | — |
| `RetryConfig{Max :306 default:"3", BackoffMS :307 default:"1000"}` | — |
| `LLMSection.Providers :419`（`map[string]Provider`） | 键＝服务商名 |

**`Provider`**（`schema.go:379-401`）
| 字段（行） | tag/默认 |
|---|---|
| `Protocol :381` | 枚举 `openai-chat\|openai-responses\|anthropic`（`:76-78`）；内置预设名可留空继承默认（`:374-378`） |
| **`BaseURL :383`** | 空＝预设默认 ⇒ 「服务商地址」那格 |
| **`APIKeyRef :387`** | 只允许 `dpapi:<id>` 或 `env:NAME`；**明文 key 在这枚结构体里结构性不存在**（`:11-17` D36 规则 5，`boundary_test.go` 钉着）；空＝无钥服务商 |
| `Billing :391 default:"pay-per-token"` / `PlanCreditTotalMicro :392` | 套餐总额＝手录，余量永远标「估算」 |
| `Compat{Loose :315 default:"false", AllowMissingUsage :316 default:"false", ExtraHeaders :317}` | — |
| `RPM :397` / `TPM :398` | 无 tag ⇒ 0＝不限流 |
| `Models map[string]ModelSpec :400` | 模型目录 |

**`ModelSpec`**（`schema.go:348-372`）
| 字段（行） | tag/默认 |
|---|---|
| `Display :350` | 空＝用模型 id |
| `Capabilities :352`（`Capabilities{text,vision,audio_in,audio_out,realtime,thinking,fc,stream} :324-333`） | 新模型**全 false**（`:346-347` 注释），且**声明必须经 probe 实测**（`SPEC-03:63-64`） |
| **`ContextWindow :354`** | 无 tag ⇒ `0 = unknown`（原文） |
| `MaxOutputTokens :356` | 0＝服务商默认 |
| `ThinkingLevels :359` | 枚举逐元素校验（`:83-88,98`） |
| `Billing :361 default:"pay-per-token"` | — |
| **`Price :363`** → `Price{In,Out,Cached,AudioIn,AudioOut} :337-343` | 无 tag ⇒ 0＝免费或未填；单位＝**每 1M token 的 micro-USD**（`:335-336`，C23/ticket 44 消费） |
| `QuotaDailyMicro :367` / `QuotaMonthlyMicro :368` | 0＝不设限；硬停 100% 走兜底链（`:365-366`） |
| `Enabled :371 default:"true"` | — |

**另两枚「模型相关」但不在 `[llm]` 的**：
`RealtimeConfig{Enabled :233 default:"false", Provider :234, Model :235, APIKeyRef :236, BaseURL :237}`（C32／D47 的云端 realtime 大脑，**也带一枚 `api_key_ref`**）；
`AgentSection.TokenBudget :433 default:"200000"` 与 `MaxRounds :431 default:"50"`（上下文预算消费方，D15）。
以及 `[risk] permission_mode :470 default:"ask_every_step"`——票面 §5 那张表里 owner 要的「权限档位」。

### ①.3 写侧每一枚字段的生效级别（三档逐字引，不自己发明）

三档的**代码定义**逐字（`internal/config/schema.go:33-44`）：
> `TierHot` "applies immediately: the new value is visible through Store.Config() as soon as the reload merge completes."
> `TierReload` "applies immediately AND emits a reload event so the owning subsystem can reload (e.g. voice model swap triggers model load/unload)."
> `TierRestart` "is NOT applied at runtime: the old value stays active until process restart. The reload report carries the pending section so the app can surface \"restart required\"."

三档的**文档口径**逐字（`SPEC-03-config-secrets-envs.md`）：
- `:75-76`「可热加载：`[llm]` `[agent]` `[session]` `[ball]` `[hotkey]` `[panel]` `[cost]` `[memory]` `[observe]` `[models]`（阈值/开关类）」
- `:77`「reload 级：`[voice]` 换模型（走重载子系统路径，触发模型加载/卸载）」
- `:78`「restart 级：`[app]`（theme 除外）」
- `:79`「🔒 安全节放宽：**弹 L2 重新确认卡（原生侧），确认后生效并写日志；拒绝则保留旧值**」
- `:32` 行尾逐字：`[llm]` ＝「**hot（api_key_ref 变更须重新解密；链/角色变更 hot）**」
- `:81-85` §4.3 逐字：「**GUI 改配置 = 写回 TOML → 同一检测路径生效（单一真相源始终是文件）**」
代码分档处：`manager.go:243-300`（四个 `planLocked` + `planApp :338` + `planVoice :361` + 其余 12 段整体 hot，`llm` 在 `:281`）。

### ①.4 哪些字段「写下去要重启才生效」——**实践面比那张表宽，这是本腿最该被读的一段**

〔起手读数，可能被 33-r1 改动〕
1. 表上写的 restart 档＝`[app]`（theme 除外）。生产里那句话在 `cmd/wisp/config_reload.go:97-99`：
   常量 `hotReloadDisabledPanelInbound`（「本宿主没有接 config.toml 轮询……要生效请重启进程并用 `wisp run`」），
   它在 `cmd/wisp/panel_inbound.go:212` 每启动念一次；`reportRestartPending` 在 `config_reload.go:278-294`，
   涉及键名册 `restartTierKeys` 在 `:299-301`（`app.language` / `app.autostart` / `app.single_instance`）。
2. **`[fs]` 放宽即使 L2 批了，本次运行也不参与判定**：`config_reload.go:198-206` 那句「新放宽的目录要重启进程才参与判定（原因：canonicalizer 只在装配时构造一次）」，
   同文件 `:55-61` 把它列为「WHAT THIS FILE DOES NOT FIX（具名，不藏）」。
3. ⚠ **`[llm]` 名义 hot，但正在跑的 provider 不会因此重建**：endpoint/provider 只在装配根构造一次——
   `run.go:410`（`res.ResolveRole(llm.RoleChat)`）、`:415` `rt.endpoint = ep`、`:416-421` `BuildEndpointProvider` → `rt.provs`；
   全仓 `grep` 只有 `cmd/wisp/providers.go:175` 另一处构造，**没有任何一处从 reload 重建 `rt.provs`**（R37，23:53:03）。
   而 `Manager.OnReload`（`manager.go:52-54`，`:198` 只对 **reload 档**触发）在 `cmd/wisp` 里**零个生产赋值点**（R40：只有 `cmd/balldebug/main.go:244`）；
   `[llm]` 是 hot 档，本来也不会走 `OnReload`（`internal/ball/hotkey_reload.go:26-30` 把这条限制写得很清楚）。
   ⇒ **对票 248 AC#4 的后果**：同一次 `wisp run` 里「保存 → 下一次任务真用上刚配的模型」这一发，按现读**不通**；
   要么新增「从热重载重建 endpoint/provider」那一跳（那是另一票），要么 AC#4 得按重启来判。**这条我判不动，交裁（⑦ 追加问，写在下面 J9）。**
4. 实时语音那一族（`voice.realtime.*`）按表是 reload 档，本进程的音频管线重挂路径不在本票射程。

### ①.5 写了会立刻改变安全 posture 的字段（最要紧那一族）

| 字段 | 为什么敏感 | 现成的闸门与出处 |
|---|---|---|
| **`risk.permission_mode`**（放宽＝往 `auto_approve` 走） | 直接移除本该询问的确认 | `manager.go:437-447`（mode 排序即宽严排序，高＝松）；`config_reload.go:221-267` `confirmLockedLoosening` 一张 L2 卡；`run.go:805-827` `confirmModeSwitch`；面板侧唯一合法入口＝`panel.mode.request` + `composer_handlers.go:133-138`（无确认腿即 `ErrNoL2Confirm`） |
| `risk.l1_window_sec` 调大 | 自动放行窗口变长 | `manager.go:432-436`；schema `:451-453` 注释自陈「raising it auto-approves more, so an increase is a loosening change」 |
| `risk.shell_enabled` / `risk.allow_shell_string` / `risk.shell_allowlist` / `risk.blacklist_overrides` | 「免审」类 | `manager.go:422-431` 判方向；⚠ **今天这四枚写了根本进不来**：`unwired.go:60-97` 在**加载时**以「写了但不起作用」拒绝（`validateUnwired`，由 `validate.go:24,32` ← `loader.go:123` 触发）⇒ 设置页若给它们摆控件，第一发点击就会把 config.toml 变成读不了的那一枚（见下一条） |
| `fs.allowed_dirs` 增项 / `fs.delete_enabled` / `fs.reparse_point_exceptions` | 扩大可读写根 | `schema.go:473-482`；写腿已存在：`allowdirs.go:62,109`，其头段 `:16-24` 写明「这里的写放宽的锁定段，**调用者欠那张 L2 卡**」 |
| `net.allowlist` 增项 / `net.block_private_ranges=false` / `net.proxy.mode` 离开 `none` | 新开外泄通道 | `schema.go:493-502`、`ProxyConfig :485-491`；后者两枚同样在 `unwired.go:86-96` 被拒 |
| `plugins.*` 全族 | 授予宿主 API 与能力 | `schema.go:551-576`；`unwired.go:138-146` 标注引擎未建（票 50/51） |
| **硬编码三枚：绝不可写** | `privacy.keep_transcript`/`keep_audio`（`schema.go:514-517`：写 true＝load error）、`models.verify_signature`（`:586-588`）、`audio.half_duplex`（`:273-275`） | 读侧可显示，**写侧必须结构性拒绝**；`SPEC-03:37,42,31` 同口径 |
| ⚠ **写路径今天没有校验闸门** | `SaveFile`（`loader.go:238-249`）只做 deepCopy + 强制 `schema_version` + `MarshalCanonical` + `atomicWrite`，**不跑 `validate()`**；`validate()` 只在 `LoadFile`（`:123`）与迁移（`migrate.go:77`）里跑 | ⇒ 一枚非法值（枚举拼错、链元素指向不存在的 `provider/model`、`..` 之外的怪路径）会**先落盘**，下一次读才炸——`config_reload.go:357-360` 的 `cause=invalid` 那一句，或下一次启动直接 Unconfigured（`SPEC-03:70-71,90`）。**任何配置写方法都必须写前 validate**，否则设置页一次点击就有一次把下次启动变砖的能力。这是我给 248-r1 的最硬一条判据。 |

### ①.6 凭据绝不通过 `panel.message.send` 传：理由 + 机读判据形状

**理由（四条，都指着行）**
1. 消息文本会进**模型请求**，也就进第三方服务商日志：`ComposerRequest.Text :77` 是 composer 的正文，
   `AgentSection.TokenBudget :433`／D15 的上下文预算把它计进去；一枚 key 进 prompt ＝交给对端保管。
2. 消息文本会进**出向快照**：`ResultChunk{Text :97}` 就是 `Snapshot.results`（`composer.go:59`）里回给页面的那一族键，`panel_pump.go:350-360` 还会把摘要与 sha256 记进持久 ledger。
3. **封套是共享的**：`bridge.go:69-79` 一枚 `ComposerRequest` 服务所有已答方法 ⇒ 一旦允许某枚带值键，
   每条已答路由都能带它（冻结件 facet 2 只看**键名是不是裁决形**，`grantRouteWords :192-195` 不含 value/secret，**拦不住值**）。
4. D30 间接提示注入：进了上下文的文本是模型可引用的，key 会落进 R4 污染面（`unwired.go:81` 那条「a file on disk cannot have clicked anything」是同一种诚实）。

**机读判据形状（可裁、非词面）**
- 判据一律**问能力**（AC#1 原话：不许做成扫注释词面）：
  「今天没有任何一枚已答方法能把凭据送进宿主」为**真**；落地后必须是
  「**恰有一枚**方法可送、它只写不回显、且 `panel.message.send` 送值必被响亮拒绝」。
- 负向尺：种 `{"method":"panel.message.send","text":"<哨兵串>"}` ⇒ 必拒，错误里点名 requestId 与「凭据不走对话通道」。
- ⛔ 单有负向尺不算交件（本仓规矩：负向尺必配「种 X 必响」的正控）：
  ① 同一枚哨兵串走**合法**写路径 ⇒ 处理器仍被调用，且快照／`[audit]` 行／持久日志三处 grep **0 命中**；
  ② 把哨兵串种进一枚既有通道的**非值**字段（如 `panel.mode.request` 的 `to`）⇒ 尺必须**不**报响，
  证明它看的是「哪枚键能带值」，而不是「整枚封套里有没有这串字符」；
  ③ 把带值键从写处理器上删掉 ⇒ 该判据必须变红（证明判据挂在真结构上，不是挂在一句注释上）。
- 机读形状本身建议按**反射读 bindable 键名册**（冻结件已有这件工具：`bindableKey`，`l2_grant_boundary_test.go:252-262`；
  它 facet 2 的做法就是「从本包自己的 JSON decode 调用点找封套，递归穿过 embedded/nested/pointer/slice」），
  只是词表要从裁决词表换成**值类词表**——⚠ 那是新建仪器，不改冻结件，别顺手改它（⑦ J1/J4）。

## ② 写路径落在谁手里

### ②.1 唯一写者与唯一持有者（现读）

- **config 的唯一写者＝`config.Manager`，而且它今天只有两枚导出写方法。**
  「一个 Manager、一份真相」那句逐字在 `cmd/wisp/run.go:381-390`〔起手读数〕：
  > "Manager rather than LoadFile (ticket 101): the permission mode is the one preference R20/M3 lets outlive a session, and it persists THROUGH the config file, so the assembly root needs the object that can write it back atomically (Manager.SetPermissionMode). **Loading the file twice over - once read-only here, once writable elsewhere - is exactly the split state SPEC-03 §3.1 exists to prevent, so there is one Manager and one truth.**"
  受保护写原语：`(*Manager).mergeWrite(ownedKey, setOn)`（`writeguard.go:111`）——它**先重读文件**、只替换自己那一枚键、
  写出后报告「真正改了的键路径」与「文件里有而本进程没有的手改」（`:145-186`），失败即回滚内存（`allowdirs.go:142-145`、`permmode.go:79-85`），
  并且只在「落盘字节＝内存字节」时才认领新 mtime+size（`:105-110,157-163`）。
  现有导出写者：`SetPermissionMode`（`permmode.go:70`）、`AddAllowedDir`/`SetAllowedDirs`（`allowdirs.go:62,109`）。
  Manager 的其余导出面只有读与重载（R10：`NewManager`/`Config`/`Resolved`/`CheckAndReload`）⇒ **今天没有通用 `Set`**。
- **凭据的唯一持有者＝`secret.Store`（DPAPI）**：`secret.NewStore` 在 `cmd/wisp/run.go:371` 是唯一装配点〔起手读数〕；
  CLI 侧 `cmd/wisp/secret.go:296-298` 另开一次同一构造（`openStore`），头段 `:31-34` 写明「不需要 mutex、last-writer-wins 是票 06 既有语义」；
  DPAPI 那一发在 `internal/secret/dpapi_windows.go:24`（`CryptProtectData`，`CRYPTPROTECT_UI_FORBIDDEN`）；
  `Store.Store/Resolve/Exists/Delete/Blobs`（R11）＋ `RedactSecret`（末 4 之外全 `*`，`redact.go:6-14`）是全部出口。
  AC#3 要的「只走一枚凭据存储」在现读上**已经成立**：`schema.go:11-17` 让明文 key 在结构体层面**结构性不存在**（`boundary_test.go` 钉），
  要新造第二套存储就得先动那枚钉着的边界。

### ②.2 面板侧只许「发起请求」要不要新依赖边（尺的读数）

- 起手名册（`go list -deps ./internal/panel`，R13/R14）：panel 的**直连**内部边只有 `risk`/`projctx`/`streamkey`（＋ embed 用的 `frontend`，⛔ 两层禁令，本腿不读不结论）；
  `secret`/`winsec`/`observe` 是**传递**边（`observe → secret`，R15）。**panel 今天不 import `internal/config`。**
- 两种接法的代价，都量过：
  - **甲：panel 直连 config** ⇒ 闭包多 **恰 1 枚**（`github.com/CarlosShao/wisp/internal/config`），wisp-internal 口径与全量口径同值（R16：23:44:00 与 23:47:34 两发都是 1）。
  - **乙：装配根注入接口，panel 不 import config** ⇒ 新依赖边 **0 枚**。仓内有现成先例：
    `composer_handlers.go:68-78` 的 `ModeWriter` 本地声明、注释逐字写着「Declared locally, **with no import of internal/perm**, so this package's dependency face does not move for one ticket's gate」，
    而 `*perm.Store` 正好满足它（`internal/perm/store.go:49-50`）。
- 无论甲乙，**写这个动作的归属都不变**：`config.Manager` 仍是唯一写者、`secret.Store` 仍是唯一凭据持有者；
  面板侧只到「发起请求 + 被响亮拒绝」为止（`composer_dispatch.go:97-111` 那些 nil 字段就是这个形状的样板）。

### ②.3 两处既有裁定——⚠ **下面两条推广是本腿做的，不是规矩本身**

1. 票 197 那条「装配根是唯一的接缝」原文在 `cmd/wisp/panel_pump.go:138`：
   "the composition root is where **A406** said the injection belongs, and it is also the only place that knows both halves"。
   **我的推广**：把「两半都在装配根才知道」换成「面板要的读面（config/secret）也只能在装配根接」——A406 裁的是 subagent 流键那一枚具体注入点，**没有裁过 config 读面**。⇒ 若要按这条拒掉甲接法，请把它当**新的一裁**登记，别当既有规矩引用。
2. 票 238 那条「正向依赖边一律不开、改注入」：本腿**没有**在 `cmd/wisp`／`internal/**` 找到它的成文落点（票 238 的标题是 first-run 建数据根那件事：`.scratch/wisp/issues/238-first-run-...md`）。
   **我的推广**：把它读成「panel 不新增 `config` 直连边」⇒ 只支持乙接法。这一条现读只支持到「本包确实一直刻意避开 `internal/perm` 的 import」这个程度（R14 + `composer_handlers.go:70-72`）。⇒ **具名等你确认**：这条推广到底存不存在、归谁写。
3. 第三枚是我自己给的：**方向不对称可以推广到配置段**（把 `ModeWriteHandler` 的「变宽要腿、变窄永远不拦」`composer_handlers.go:16-19,133-138` 推广到所有 🔒 段）——
   这条在票面没覆盖，写进⑦ J5 等你裁，本腿不按它写结论表。

## ③ 读侧怎么回答「配了没给／配好了」而不泄露值

### ③.1 四态判据（可机读，且每一态都有现成的源）

| 态 | 机读判据（不碰值） | 现读源 |
|---|---|---|
| **没配** | `api_key_ref` 这一枚键的字符串长度为 0 ⇒ `unset` | `schema.go:387`（空＝无钥服务商）、`:236` |
| **配了没给**（写了 ref，凭据不在） | `dpapi:` 族：`Store.Exists(ref)` 为 false；`env:` 族：`os.LookupEnv(name)` 为空 | `store.go:129-145`（**`Exists` 只答 `dpapi:`，`env:` 直接报错** ⇒ 两族必须分开判，不能让一枚尺去问另一族）；`refs.go:26-44 ParseRef` 给出 `kind/value` |
| **配好了**（解析得开） | `Store.Resolve(ref)` 无错，或 `llm.Resolver` 走到 `BuildEndpointProvider` 不返回 Unconfigured | `resolver.go:141-144`（解不开＝Unconfigured 文案）、`:139`（有 ref 无 resolver＝Unconfigured） |
| **不可判**（读不了文件／没权限） | 上面三条**任一返回 error 而非 false** ⇒ 必须自成一态 | 与 `describeReloadFailure` 四支不同形同纪律：`config_reload.go:315-369`（缺失/语法/权限/schema 四支各一句，票 223 AC#4 禁止合并） |
| **配了但目录里没这个模型** | 链元素引用的 `provider/model` 不存在 ⇒ 校验点名失败 | `schema.go:409-411`、`SPEC-03:59-60` |

**绝不进快照的**：任何 `Resolve` 的返回值。读侧只能出**枚举名 + 计数 + ref 名**（ref 不是秘密：`dpapi:<blob-id>` 是文件名，`refs.go:70-79` 自陈「随机 id 所以 ref 不泄漏它是哪家的」；`resolver.go:154 secretRefShape` 也是把 ref 写进错误文案的既有先例）。
⛔ 也不许出「末 4 位」：`RedactSecret`（`redact.go:6-14`）是**CLI 给人核对用的终端输出**，不是页面字段；AC#3 那句「不许新增任何把值回显给页面的方法」把 `api_key_ref` 之外的一切都算回显。

### ③.2 「宁缺毋造」那条规矩写在哪 + 出向快照加那一维要动哪几枚文件

规矩的**成文位置**（R07，23:42:55）：`internal/panel/pump.go:25-46`（"NO NEW KEYS"，写明加第五枚键会同时红哪四枚钉）与
`internal/panel/composer.go:44-56`（四枚钉名册：`composer_test.go:48`、`approval_test.go:105`、`pump_test.go:111-124`、`:270-276`，加一句「a key without that reader is the constant this file was moved here to stop」）；
票面口径在 `docs/evidence/s1/145-snapshot-field-census-r1.md:63`（「一枚字段若只能靠常量、靠 0、靠把另一行的文案挪过来才'通'，本件一律判 **无源**」）与 `145-snapshot-fields-landed-r1.md:230-236`。
> ⚠ 上面那四个「钉」的行号是**`composer.go:45-49` 自己写的**，不是我先量出来的。我随后独立复量的结果（23:58:46）：
> `composer_test.go:48` ＝ `TestComposerContractTypesMatchFrontend` ✓；`approval_test.go` 里读 `frontend/src/lib/panel.ts` 那行在 **`:107`**（注释说的是 :105，函数体偏移两行）；
> `pump_test.go` 那句 "the four keys are the contract" 在 **`:109`**（注释说的是 :111-124）；另 `composer_test.go:502`/`:533` 两枚 ✓、`git_test.go:362` ✓。
> ⇒ 名册可信，行号有 ±2 的漂移，248-r1 动手时**按函数名找，别按我抄的行号找**。

**现量更正**（R08；这条与票面现量 #5 不同，要按本件为准）：`Snapshot` **不是四枚键**，而是
`pending`/`results`/`composer`/`generatedAt` 四枚恒发（`composer.go:58-61`）+ `instructions,omitempty`（`:74`，票 200）+ `tasks,omitempty`（`:91`，票 197）＝**结构层 6 枚**；
票 145 AC#3 已裁第五枚键「不在那片切片」（`pump.go:27-29`，Q-51 未定案）。

要加「凭据是否已录入／配置是否可读」那一维，动哪几枚文件（**不新增通道，两接法都列**）：
- **甲（我倾向，但仍交裁）**：加进 `composer` 段——`internal/panel/composer.go:235-256` 的 `ComposerState` 已经有
  「一枚事实 + 一枚 known 位」的先例（`CurrentModel :255` / `ModelKnown :256`），照它的形加
  `credentialState`（枚举）/ `configReadable`（枚举）两枚键 ⇒ 顶层键集不动，`Snapshot` 的 4 枚恒发不动。
  ⚠ 代价：**同样动双向核对**（`composer_test.go:48`、`approval_test.go:105` 是对着页面声明比键集的），所以仍要过⑦ J3。
- 顶层加第五键：`composer.go:57-92` + `cmd/wisp/panel_pump.go` 的构造点 + `cmd/wisp/run.go:673` 的装配处加 reader〔起手读数〕。
- **两接法都跑不掉的那一步**：reader 必须在装配根递给泵，因为
  `NewSnapshot`（`composer.go:106-131`）对没读的字段一律填「unreadable/unknown」那一支（`:115-126` 那三处「绝不渲染成听起来最安全的那个」），
  没有 reader 的键就是 `composer.go:52-55` 点名要停掉的那枚常量。
- 哨兵正控落点（AC#2 要常驻用例，不是表里一次性读数）：写路径 → `bookPanelSnapshot`（`panel_pump.go:319-325`）与
  `panelSnapshotSummary`（`:331-362`，它已经自带 440 字符钳 `summaryClamp :367`，理由就是 ledger 单串上限 512）
  → grep 面要覆盖：`[audit]` 族、持久日志 sink（`installLogSink`，`cmd/wisp/secret.go:226` 是同族调用点）、快照字节。

## ④ C17 那张方法表在文档里有几处拷贝、漂移会出现在哪几行、有没有仪器看得见

> ⛔ 本节**只读只引行号**，`docs/PLAN.md`／`docs/specs/**` 一字未改。

**「拷贝」的真身只有一处是表**；其余 C17 出现都是**要求**或**追溯链接**，不枚方法名。

| # | 位置 | 里面有什么 |
|---|---|---|
| 1 | `docs/specs/SPEC-08-ui-ball-panel.md:156-175`（§5.2「C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】」） | **唯一一张名册表**：`panel.resync`、`tasks.list`/`task.detail`、`history.query`/`transcript.get`、`approval.current`/`approval.queue`、`approval.decide`（「allow」拒一切面板来源、仅 reject 可面板发起）、**`config.get`/`config.set`（旁注「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」）**、`grants.list`/`grants.revoke`、`privacy.purge`/`privacy.export`、`cost.summary`、`models.list`/`models.delete`、`diagnostics.export` + 事件推送五枚 |
| 2 | `docs/PLAN.md:2434`（F2 配套） | 只写要求：必须给方法白名单、逐枚标 capability 与是否需原生二次授权、未列出→拒绝并记日志。**无名册** |
| 3 | `docs/PLAN.md:2968-2974`（D39「C17 补四项」） | 四项要求（白名单／capability+二次授权／correlationId 背压合并／`panel.resync` 强制无状态）。除 `panel.resync` 外**无名册** |
| 4 | `docs/PLAN.md:1367`（C17 契约行） | 通道定义（`invoke(method,args)` + Go→前端事件、correlationId 路由、前端必须无状态）。**无名册** |
| 5 | `docs/PLAN.md:1747` | 交付物索引 `docs/contracts/C1..C31.md`——**该目录不存在**（AGENTS.md §4 已警告）⇒ 那张表本该有的第 4 份拷贝**从未落地** |
| 6 | `docs/PLAN.md:1808` | 自陈「C6/C7/C11/**C17**/C19 只有名字没有内容 ⇒ 补规格（D39）」 |
| 7 | `docs/PLAN.md:2164` | D22「禁改 D1-D25」三处早于 C17-C23 的更正记录 |
| 8 | `docs/specs/SPEC-08:3`、`SPEC-01:59`、`SPEC-10:3,18`、`SPEC-12:48` | 追溯链接与测试口径（`SPEC-10:18`：面板＝C17 接缝，测「路由/白名单/resync/审批拒绝面板来源」）；`SPEC-12:48`＝**未定案项**「C17 方法白名单定稿（S5 切片卡批准）」 |
| 9 | `docs/specs/SPEC-03-config-secrets-envs.md:34`（🔒 `[risk]` 行） | 段级名册，不是方法表——但**缺 `permission_mode`**（R21：SPEC-03 全文 0 命中，代码 `schema.go:470` 有）⇒ schema 层已经在漂 |

**代码侧的同一件事散在哪**（这才是漂移面）：`bridge.go:42-45`（名册真身）→ `bridge.go:104-110`（守卫）→
`composer_dispatch.go:137-158`（派发表）→ `composer_test.go:409-419 composerRouteLiterals()`（渲染侧允许的字面量集，含 `"panel.approval.request"`）→
`l2_grant_boundary_test.go:1496 / :1635 / :1983 / :2043 / :2051 / :2139`（冻结件的尺寸与拼写钉）→
`git_test.go:381-390 / :514 / :517`（票 181 的反漂移钉，也写死四枚）。**共 6 个文件里 11 处名册/尺寸表述。**

**代码扩了而文档没扩，漂移会长在哪几行**：
- `SPEC-08:168` 那行（`config.get` / `config.set` 已在文档、不在代码）——**方向与常规相反**：文档先行、代码未跟上。
  若 248 落地时选的名字**不等于**这两枚（例如 `panel.config.read/write`），那行就变成**两处不一致的真名册**；
  若名字恰好等于，那行才从「提案」变成「已实现」（它头上还挂着 `SPEC-12:48` 那句「S5 定稿走契约批准」）。
- `SPEC-08:167`（`approval.decide` 面板侧仅 `reject`）——文档用**未加 `panel.` 前缀**的名字拼裁决路由；
  代码侧这枚名**必须永远不被答**（`l2_grant_boundary_test.go:1248-1267` 把它列进「不许被答」候选）。这条不是新漂移，是**既有**术语形状不一致（文档 `x.y`，代码 `panel.x.y`）。
- `SPEC-08:163` 的 `panel.resync` 是文档表里唯一带 `panel.` 前缀的名字（另在 `:150` 出现在正文）——它**不在代码名册**
  （Go→页面方向根本没有通道，`pump.go:15-23` 自陈）；`SPEC-08:238` 还把 `approval.decide({allow:true})` 当成预期行为写着，
  而那一发今天是由**冻结件**保证的，不是由那张表保证的。
- 表外那些（`tasks.list`/`history.query`/`grants.*`/`privacy.*`/`cost.summary`/`models.list`/`diagnostics.export`，`:169-173`）**一枚都不在代码里**；
  票 181 与票 187 已各自把「模型/thinking 档」「worktree 切换」这类名册缺口登过票（`.scratch/wisp/issues/187-...md:24`、`181-...md:20`）。

**有没有仪器能看见（现量＝没有，你预计成立；但这条尺我第一次读错了，见⑥ 第 11 条）**：
- `grep -rn "SPEC-08" --include=*.go .`（去 `.scratch`）**82 行命中**（23:55:34），**但没有一枚 Go 文件把 `SPEC-08` 从盘上读进来**：
  同一批命中里 `grep ReadFile|os.|filepath` **0 命中**（23:55:42），82 行全是注释／flag 文案（如 `internal/session/grants.go:40`、`cmd/balldebug/main.go:50`）。
  ⇒ 「文档那张表是**给人读的引用**，不是任何尺的预言机」这一条成立；我先前写的「0 命中」是错读，已按这次读数改。
- 全仓唯一一枚「拿文档当预言机」的尺是 `git_test.go:365 os.ReadFile(docs/PLAN.md)`，它核的是 **D34 工具表**里不许出现 `git.*` 行（`:359-363` 三条门），**与 C17 名册无关**。
- `tools/d22scan/main.go:5-40` 八条禁令里没有任何一条看方法名册（R28）。
- 反向倒是有三枚**很硬**的尺（都只钉代码↔代码）：`git_test.go:517`（`len(real) != 4`）、
  `composer_test.go:502-531`（渲染侧每一枚 `panel.*` 字面量必须是 Go 答的那枚）、
  `l2_grant_boundary_test.go:1496`（守卫在 24 枚网格名里答的必须**恰好是那四枚**）。
  ⇒ **实际后果**：代码加名 ⇒ **红的是钉着「四」的那几枚尺**，而文档那张表**一声不响**。
  漂移是看不见的，**不漂移反倒会被看见**（这条判断是结构推理，我没跑测试，见⑥ 第 2 条）。
- ⚠ 一句必须说的：**「名字对得上文档」并不是安全判据**。冻结件 facet 1（`:1248-1267`）与 `carriesGrantWord`（`:242-250`）
  管的是「哪一枚被答」的形状（裁决词不进出向名册），那比文档同步更硬。任何给 248 加的新尺都得跑在**同一枚运行守卫**
  （`knownComposerMethod`）上，而不是跑在一张表上——`:1435-1455` 那段「sweep 的 ask 与 accumulate 被短路 ⇒ `hits=[]` 读起来像干净」是这条纪律的现成教训。

**⑦ 追加问（写这张表时新撞上的一条，票面没覆盖）**：

> **J9** ①.4 第 3 条那个「`[llm]` 名义 hot、实则本进程的 endpoint 不重建」与 **AC#4**（保存 → 下一次任务真用上刚配的模型）正面冲突。
> 现量：`run.go:410-421` 只在装配时构造 endpoint/provider（R37）；`Manager.OnReload` 在 `cmd/wisp` 零生产赋值点（R40），且 `[llm]` 是 hot 档本来不走 `OnReload`（`internal/ball/hotkey_reload.go:26-30`）。
> - 甲：本票范围内加「热重载后重建 provider/endpoint」那一跳（会动 `cmd/wisp` 与 `internal/llm`，比票面写的范围大）。
> - 乙：AC#4 按「重启后真用上」判，并让面板明确念出「这次改动要重启才用上」（复用 `reportRestartPending` 那一形，`config_reload.go:278-294`）。
> - 不做：先承认 AC#4 这格在重建那一跳落地前**永远勾不上**（与票面「票 33 宿主没起来之前不许勾」同形，但这是另一枚理由）。
> 我一条都不选：这条决定 248 的收口口径，不是只读腿的活。

## ⑤ 我跑了哪些尺、每条真实读数

> 纪律：本腿⛔ 未跑 `go build`／`go vet`／`go test`／任何 `./...`；只跑过 `grep`／`find`／`ls`／`sed -n`／
> `git status`／`git log`／指名单包的 `go list -deps`。所有当下状态断言与取数同发，带 `date` 钟点。
> 共 **42 把尺**（R01-R42；R34 那把我跑空后重跑，R42 那把我读错后重测）。

| 号 | 尺（命令或读法） | 真实读数 | 钟点 |
|---|---|---|---|
| R01 | 编排者原尺：`grep -rn "panel\.[a-z]*\.[a-z]*" --include=*.go internal/panel/ \| grep -v _test \| grep -o … \| sort -u` | **恰四枚**：`panel.attachment.add`／`panel.message.send`／`panel.mode.request`／`panel.workspace.request` | 23:42:32 |
| R02 | 同一把尺扩到 `internal/ cmd/ tools/`（非测试） | 同样**四枚**，一枚不多 | 23:42:32 |
| R03 | 扩到测试全量计数（`uniq -c`） | 去重 **23 枚**；最高 `panel.approval.request` 17 次、`panel.review.allow` 11、`panel.mode.set` 7 | 23:42:32 |
| R04 | `grep -rniE "panel\.(config\|secret\|credential\|settings)[a-z._]*" --include=*.go .` | **零命中** ⇒ 票面现量 #1 成立 | 23:42:32 |
| R05 | 每枚字面量的 `file:line` 归属 | 生产码只有 `internal/panel/bridge.go:42-45`（声明）、`knownComposerMethod`（派发）、`composer_handlers.go:85,109`、`cmd/wisp/panel_inbound.go:196`；其余 19 枚**全在 `*_test.go`**，其中绝大多数在冻结件 `l2_grant_boundary_test.go`（植入负控） | 23:42:32 |
| R06 | `sed -n '95,135p'`／`'490,560p'` `internal/panel/composer_test.go` | `panel.mode.set` 是**被禁形状**（`composerModeWriteRe`），只出现在植入件里；`TestTheRendererHoldsExactlyOneDoorToTheHost` 要求「渲染侧每一枚 `panel.*` 字面量都是 Go 侧答的那枚」 | 23:42:47 |
| R07 | `grep -rn "宁缺毋造" --include=*.go --include=*.md .` | 16 处，**全在 `.scratch/**` 的票与派单、`docs/evidence/**`**；Go 侧成文的同规矩在 `internal/panel/pump.go:25-46`（"NO NEW KEYS"）与 `composer.go:44-56` | 23:42:55 |
| R08 | `grep -n "json:\"" internal/panel/composer.go` | `Snapshot` 结构体 **6 枚带 tag 字段**：`pending`/`results`/`composer`/`generatedAt` 恒发 + `instructions,omitempty`(:74) + `tasks,omitempty`(:91)；`ComposerState`（`:235-256`）**9 枚键**：mode/workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError/git/currentModel/modelKnown | 23:43:08 |
| R09 | `sed -n '340,430p' cmd/wisp/run.go` | 单写者那句在 **`run.go:381-390`**；`secret.NewStore` 唯一装配点在 `run.go:371` | 23:43:08 |
| R10 | `grep -nE "^func (m \*Manager)…"` `internal/config/manager.go` | Manager 导出面只有读与重载：`NewManager`/`Config`/`Resolved`/`CheckAndReload`（+`plan*` 私有），**没有任何通用 Set** | 23:43:15 |
| R11 | `grep -nE "^func (s \*Store)…"` 四枚文件 + `sed -n '18,30p' internal/secret/dpapi_windows.go` | `Store.Store/Resolve/Exists/Delete/Dir/Blobs`；`protect` 在 `:18-28`，`CryptProtectData` 那行＝**`:24`**（票面引用成立） | 23:43:15 |
| R12 | 写面检索：`grep -rn "^func " writeguard.go permmode.go` + `grep -rn "SetPermissionMode" \| grep -v _test` | 受保护的写原语只有 `(*Manager).mergeWrite(ownedKey, setOn)`（`writeguard.go:111`）；其上只有**两枚导出写者**：`SetPermissionMode`（`permmode.go:70`）、`AddAllowedDir`/`SetAllowedDirs`（`allowdirs.go:62,109`）。生产调用者：`internal/perm/store.go:208` | 23:43:21 |
| R13 | `go list -deps ./internal/panel`（指名单包） | wisp-internal 起手名册 **8 枚**：`frontend`(embed)/`observe`/`panel`/`projctx`/`risk`/`secret`/`streamkey`/`winsec` | 23:43:36 |
| R14 | `grep -rn "CarlosShao/wisp/internal/" internal/panel/*.go \| grep -v _test` | 直连边只有三枚：`risk`(5 文件)、`projctx`(2)、`streamkey`(1)；**panel 不直连 `config`，也不直连 `secret`** | 23:43:46 |
| R15 | 谁把 `internal/secret` 拉进 panel 的闭包 | `internal/observe -> internal/secret`（redact 那侧），是**传递**边不是面板的直连边 | 23:43:52 |
| R16 | 提议接法的依赖增量：`comm -13` 两份 `-deps` 读数（wisp 口径 + 全量口径） | `panel -> config` 多出 **恰好 1 枚**（两种口径同值：只有 `internal/config` 本身；第三方零增量） | 23:44:00 / 23:47:34 |
| R17 | `grep -rn "C17" docs/PLAN.md docs/specs/` | **13 行 / 5 个文件**：`PLAN.md:1367,1747,1808,2164,2434,2968`；`SPEC-01:59`；`SPEC-08:3,156,235`；`SPEC-10:3,18`；`SPEC-12:48` | 23:43:36 |
| R18 | `sed -n '150,190p' docs/specs/SPEC-08-ui-ball-panel.md`（§5.2 那张表） | 文档表列 **20 余名**（`panel.resync`/`tasks.list`/`history.query`/`approval.decide`/**`config.get` / `config.set`**/`grants.*`/`privacy.*`/`cost.summary`/`models.list`/`diagnostics.export` + 五枚事件推送），标题自署**【SPEC 提案，S5 定稿走契约批准】** | 23:45:06 |
| R19 | `sed -n '2960,2985p'`／`'2430,2440p'` `docs/PLAN.md` | PLAN 侧对 C17 只提**四项要求**（白名单/capability+是否二次授权/correlationId 背压/`panel.resync` 无状态），**全篇未枚一枚方法名** | 23:45:06 |
| R20 | `sed -n '20,95p' docs/specs/SPEC-03-config-secrets-envs.md` | §3 表 `:24-43`（逐段生效级别）；§4.2 `:73-79`（hot/reload/restart/🔒 四支）；§4.3 `:81-85`：**"GUI 改配置 = 写回 TOML → 同一检测路径生效（单一真相源始终是文件）"** | 23:46:10 |
| R21 | `grep -n "permission_mode" SPEC-03` | **0 命中** ⇒ 文档 §3 的 `[risk]` 行（`:34`）里没有 `permission_mode`，代码 `schema.go:470` 有 ⇒ 文档↔代码**已在漂** | 23:46:24 |
| R22 | `sed -n '243,300p' internal/config/manager.go` | 分档实现：`planLocked` 四枚锁定段 → L2；`planApp`＝restart（theme 除外）；`planVoice`＝reload；其余 12 段整体 hot（`llm` 在 `:281`） | 23:46:33 |
| R23 | 冻结件 `internal/panel/l2_grant_boundary_test.go` 五处读段（头段 1-60、1240-1290、1360-1395、1408-1510、2021-2160、2290-2345） | 四条 facet + 植入件族；`git log` 该文件工作树**干净**、2380 行（23:48:08 同发） | 23:44:11-23:49:36 |
| R24 | `grep -n "grantRouteWords = " -A 20` + `carriesGrantWord` 实现 | 裁决词表 11 枚：`approve/approval/grant/allow/permit/ratify/authorize/authorised/decide/decision/verdict`；归一化**抹掉 `._-` 与空格再子串匹配**（`:242-250`） | 23:44:23 |
| R25 | 尺寸／拼写硬钉位置检索（`wantAnswered`、`!= 4`、`guardAnchor`、`fieldAnchor`） | `:1469 wantAnsweredInGrid=4`、`:1496 wantAnswered{四枚常量}`、`:1635`/`:1983` 逐枚循环、`:2043 fieldAnchor="\tText string \`json:\"text,omitempty\"\`"`、`:2051 guardAnchor="case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:"`、`:2139 got != 4` | 23:44:30 / 23:48:56 / 23:49:15 |
| R26 | `sed -n '350,400p'`／`'505,525p'` `internal/panel/git_test.go` | 票 181 的反漂移钉 AC#3 门 3：`wantMethods` = 四枚常量（`:381-390`，排序全等比较）；同一函数里还有正控（植入第五枚 `panel.worktree.switch` 必数到 5，`:514`）与实尺（`len(real) != 4`，`:517`） | 23:45:19 / 23:48:35 |
| R27 | 「有没有仪器把代码方法名核文档那张表」：`grep -rln "SPEC-08\|docs/PLAN.md" --include=*.go` + `grep -rn "docs/specs\|ReadFile" ` + `grep -rniE "config\.get\|config\.set\|models\.list" --include=*.go` | **没有。** 全仓只有一枚仪器把文档当预言机：`git_test.go:365 os.ReadFile(docs/PLAN.md)`，核的是 **D34 工具表**里不许有 `git.*` 行；Go 侧提到文档方法名的只有 `internal/session/grants.go:40` 一句注释 | 23:45:19 / 23:45:35 |
| R28 | `grep -nE "1 bare-goroutine…8 emoji" tools/d22scan/main.go` | 八条禁令射程内**没有任何一条**看 C17 方法名册；#3 plaintext-key 只看「key/token/secret 命名的标识符被赋字符串字面量」 | 23:45:35 |
| R29 | `grep -n "confirmModeSwitch\|rt.pump = " cmd/wisp/run.go` + 函数体 | `confirmModeSwitch` 在 `run.go:805-827`（一张 C18 卡、`host:mode-switch`、非 Allow 即错）；只有 `wisp run` 装它（`:626`）；泵装配在 `run.go:673` | 23:46:06 |
| R30 | 凭据读侧原语语义：`sed` `store.go:120-175`／`refs.go:1-80`／`configrefs.go`／`redact.go` | `Exists` **只答 `dpapi:`**，`env:` 直接错（`store.go:129-145`）；`ConfigRefs(configPath) map[字段]ref`（`configrefs.go:23`）；`RefFieldNames(configPath, ref)`（`:62`）；`RedactSecret`＝除末 4 字符全 `*`（`redact.go:6-14`）；`NewRef()` 随机 hex（`refs.go:72`） | 23:46:47 |
| R31 | `grep -nE "func \|Err \|api_key" internal/llm/resolver.go` | 解不开 key 的两个响亮出口：`:139`（有 ref 无 resolver，Unconfigured）、`:141-144`（`Keys.Resolve` 失败 → Unconfigured）；错误文案带 ref **形状**不带值（`secretRefShape`，`:154`） | 23:46:54 |
| R32 | `grep -n "func routeLabelsInFunc" -A 40` 冻结件 | 枚举行走的是**函数内每个 `ast.SwitchStmt` 的每条 `CaseClause`**（`:887-913`）⇒ 新增独立 case 行也进名册，`:2139` 那枚「4」照样变 5 | 23:49:51 |
| R33 | `git status --short` + `git rev-parse --abbrev-ref HEAD` + `git log --oneline -3` | 起手锚点 `2b4e2602`（dev），工作树有大量他人未提交改动（`design/**` 删除等）⇒ 我全程只用显式 pathspec | 23:41:16 |
| R34 | ⚠ 一把跑空的尺：`grep docs/specs/SPEC-03-config-secrets.md` | `No such file or directory`（实名 `SPEC-03-config-secrets-envs.md`）⇒ 23:46:06 那次读空，23:46:10 改用通配重跑才有 R20 的读数 | 23:46:06 |
| R35 | 分向谓词：`grep -n "func riskDirection" -A 26 internal/config/manager.go` | `manager.go:421-447`：`shell_enabled`/`allow_shell_string` false→true＝松（`:422-431`）、`l1_window_sec` 调大＝松（`:432-436`）、`permission_mode` 按 rank 比较，**不可解析的值两侧都解析成最严 ⇒ 只可能被记成收紧或中性**（`:437-447`，注释写得很清楚） | 23:52:20 |
| R36 | 读写两侧的实现：`sed` `loader.go:64-90`／`loader.go:238-249` | `readConfigFile` 会走迁移与校验；`SaveFile` 只做 deepCopy＋强制 `schema_version`＋`MarshalCanonical`＋`atomicWrite` | 23:52:20 |
| R37 | 端点重建检索：`grep -rn "BuildEndpointProvider\|ResolveRole\|rt.endpoint\|rt.provs" cmd/wisp/ \| grep -v _test` | 构造点只有 `run.go:410/415/416/421`（装配根）与 `providers.go:175`；消费点 `run.go:960-962`；**没有任何一处从 reload 重建 `rt.provs`** | 23:53:03 |
| R38 | 两条既有裁定的落点检索：`ls .scratch/wisp/issues/ \| grep -E "^238\|^197"` | 票 197 的成文句在 `cmd/wisp/panel_pump.go:138`（引 `A406`）；票 238 的文件名是 first-run 建数据根那件事（`238-first-run-on-a-fresh-machine-...md`），**我没找到「正向依赖边一律不开」的成文落点** ⇒②.3 第 2 条按推广写 | 23:53:03 |
| R39 | reload 消费者：`grep -rn "rep.Reload\|OnReload" cmd/wisp/ internal/ \| grep -v _test` | `cmd/wisp` 只在 `config_reload.go:169` **报** `rep.Reload`，无人接；`internal/ball/hotkey_reload.go:26-30,107-120` 是唯一的 OnReload 形 | 23:53:11 |
| R40 | `grep -rn "OnReload\\s*= " --include=*.go .`（去测试与 `.scratch`） | `Manager.OnReload` 在 **`cmd/wisp` 零个生产赋值点**；全仓只有 `cmd/balldebug/main.go:244` 一枚 | 23:53:19 / 23:53:23 |
| R41 | 校验入口归属：`grep -rn "validateUnwired\|func validate(" \| grep -v _test` ＋ 调用点筛 | `validate()`（`validate.go:24`，内含 `validateUnwired` `:32`）**只有两个调用点**：`loader.go:123`（LoadFile）与 `migrate.go:77` ⇒ **受保护写路径 `SaveFile`/`mergeWrite` 今天没有写前校验** | 23:52:29 |
| R42 | 复查我自己读错的那把尺：`grep -rn "SPEC-08" --include=*.go`（去 `.scratch`）＋ 同批命中里再筛 `ReadFile\|filepath` | 命中 **82 行**（不是 0），但把文档当数据读的 Go 文件 **0 枚**；82 行全是注释／flag 文案 ⇒ 结论不变、我原来的证据写错了 | 23:55:34 / 23:55:42 |

共 **42 把尺**（R01-R42，含 R34 那把我跑空了又重跑的、与 R42 那把我读错又改对的）。

## ⑥ 我可能写错的条目（对抗我自己）

1. **我最先要推翻我自己的一条**：写这节前我曾把「白名单扩一枚 ⇒ 冻结件红」的证据记成 `l2_grant_boundary_test.go:2144-2153`
   一段写着 `want 4: the whitelist changed size` 的文本。
   23:48:28 与 23:48:35 两次复查：`grep -rn "changed size" --include=*.go .` 与 `grep -rn "did not move with it" .` **全仓零命中**。那段文字不属于本仓任何文件。
   ⇒ 本件只承认查得的三枚真钉：`:2051 guardAnchor`（精确拼写）、`:2139 got != 4`（植入件的实尺快照）、`git_test.go:385/:514/:517`。任何后来人按我原来那句去找 `:2144` 会找不到。
2. **我一个测试都没跑**（票面禁令：246-r2 正在跑 `cmd/wisp`）。所以「加第五枚会让 `:2139` 红」「facet 2 不拦 config 族键名」这类结论**都是从码面推出来的结构判断，不是读数**。裁决者应把它们当预演，不该当已验。
3. R08 数的是**结构体 tag**，不是**发射态**。票面现量 #5 说「快照今天只有四键」，`pump.go:25-46`／`composer.go:44-56` 也说四键，但 `Snapshot` 有 6 枚带 tag 字段，其中两枚带 `omitempty`。⇒ 真实一发到底送 4 枚还是 6 枚，只有跑泵才知道；本件只报「4 恒发 + 2 视装配」，**没跑过就不写死**。
4. owner 给的 `run.go:356` 与我读到的 `:381-390` 不是同一处 ⇒ 行号已被 246-r2 落地推移。我引的是〔起手读数，可能被 33-r1 改动〕，**别按 356 去对质**。
5. 我把「SPEC-08 §5.2 有 `config.get`/`config.set` 而代码没有」写成「文档↔代码漂移」。可能判错：那张表标题自署**【SPEC 提案，S5 定稿走契约批准】**（`SPEC-08:156`），它可能是一张**尚未定稿的提案**，不是「批了没实现」。两种读法对 248-r1 的命名自由度影响不同 ⇒ 交裁（见⑦ J8）。
6. facet (2)（信封不得绑定裁决形键名）我只读了头段、`carriesGrantWord`、`grantRouteWords` 与调用点，**没逐行读 `scanGrantBoundary` 的键名收集实现**。如果它对「值／secret 形键名」还有额外约束，我①里「config/secret 族名安全」的结论就只覆盖了裁决词表那一半。
7. 冻结件行号可能在我只读期间被别的腿改动。同发核过的那一发是 23:48:08：`git status --short` 对 `internal/panel/l2_grant_boundary_test.go` 与 `bridge.go` **干净**、`wc -l`=2380。此后引用不带同发读数的一律按「某时刻读数」读。
8. R16「panel→config 只多 1 枚」是**新增直连边**的代价；如果实现腿改成从装配根（`cmd/wisp`）把 config 读写腿以接口注进 panel（像 `composer_handlers.go:73-78` 的 `ModeWriter` 那样**本地声明接口、不 import**），代价＝**0 枚新依赖边**。我把两种接法都列进②，没替实现腿选。
9. ③里「消息文本进上下文预算与日志」我只给了**判据形状**，没实测 logsink 的落盘目录与轮转后可 grep 面（`installLogSink` 我在 `secret.go:226` 读到调用，没读到实现细节）。哨兵正控真正要做到的覆盖面，得由跑的人定。
10. 我把 `panel.mode.set` 读成「被禁的渲染侧拼写」（R06）。它是**植入件**里的字符串，不是名册成员；如果实现腿误以为它是既有能力去「复用」，会把 `composer_test.go:111-120` 那枚正控钉踩响。本件按 R05 的归属读数写它＝只存在于测试。
11. **我在 ④ 的第一版里把一把尺的证据写错了**：我写成 `grep -rn "SPEC-08" --include=*.go` 命中 0。23:55:34 重跑（R42）＝**82 行命中**（如 `internal/session/grants.go:40`、`cmd/balldebug/main.go:50`），
    只是其中**把文档从盘上读进来的文件是 0 枚**（同批筛 `ReadFile|filepath` 无命中）。⇒ ④ 的结论（没有任何尺核「代码名册 ⊆ 文档表」）**不变**，
    但证据换了；后来人若按我原来那句「0 命中」去复核会当场推翻我。这条按本仓规矩留在⑥，不是留在正文里当没发生。
12. ①.5 表末那行「写路径今天没有校验闸门」是我从 `SaveFile` 函数体（R36/R41）读出来的**结构性判断**，
    我**没有**种过任何一枚非法值去验它真的能落盘（本腿禁跑测试，也不该去动生产配置）。
    ⇒ 这条最该由 248-r1 或裁决者用一发实测转成读数；如果它错了（比如 `mergeWrite` 上游另有校验），①.5 那行的判据强度要降级。

## ⑦ 判不动的地方（逐条甲／乙／不做 + 现量）

> 以下**九条**一条都没动，全部交编排者裁。**J1 是这票真正的闸门**；
> **J9**（热加载与 AC#4 冲突）写在 ④ 之后（①.4 第 3 条指它），内容与其他九条同级，别漏。

**J1 冻结件与 AC#1 正面相撞：白名单每扩一枚，`l2_grant_boundary_test.go` 必红。**
现量：①`:2051` 要求 pristine `bridge.go` 里**逐字**含 `case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:`，不含即 `t.Fatalf("plant B has no anchor")`；②`:2139` 在「复制真实包 + 植入」的快照上断 `len(pkg.answered) != 4` 即红，而 R32 证明 `routeLabelsInFunc` 走**全部** case 子句；③`:2043` 还逐字钉着信封里那行 `Text string`。该文件是 AGENTS.md §1.5／票 248 AC#5 的**三枚冻结件之一，实现腿禁改**。
- 甲：把「冻结件里那两枚**尺寸/拼写**钉」按人工批准改成**从 guard 派生**（`len(answered)` 与常量名册一致，不再写死 4），与 `A484` 同批登记；`AGENTS.md §1.2` 的「不许为变绿放宽断言」在这条形不成文——它是尺寸快照不是安全断言，但这话不归我说。
- 乙：白名单**不扩进 `knownComposerMethod`**，新能力另开一条管道——但 AC#1 明确要求「进 `bridge.go` 白名单与 `ParseComposerRequest`」，且另开门正是该冻结件 facet 1 的「second switch/if chain」要判红的形。
- 不做：本票 248 的 AC#1 在冻结件重排之前不可落地；先做②③（读侧那一维）也不通，因为没有入口。
⇒ **要我拍**：甲/乙/不做，选一个。

**J2 票 181 的反漂移钉同样写死四枚（非冻结，但那是别人 AC#3 的判据）。**
现量：`internal/panel/git_test.go:381-390` 排序全等比较 + `:517` `len(real) != 4` + `:514` 植入第五枚必数到 5 的正控。
- 甲：248-r1 在同一枚 commit 里把 `wantMethods` 与那句「this ticket may not extend」的报错文案一并改写，并在证据里写明是票 248 动的。
- 乙：把票 181 那条判据改成「读侧展示不新增方法」的性质判断（能力形），不再数枚数——AC#1 的判据形状本来就写着「判据一律问能力，不许做成扫注释词面」。
- 不做：让 248-r1 去撞一枚写着「本票不许扩」的尺。
⇒ 要个口径；我倾向**乙**，但那是票 181 的判据，不是我的。

**J3 快照那一维落在哪：AC#2 要第五枚键，`pump.go:25-46` 写「NO NEW KEYS」，票 145 AC#3 已裁「不在这片切片」，Q-51（谁准动 `frontend/src/lib/panel.ts`）未定案。**
现量：R08（4 恒发 + 2 可选；`ComposerState` 另有 9 枚键）；四枚钉住键集的件列在 `composer.go:44-56`（`composer_test.go:48`、`approval_test.go:105`、`pump_test.go:111-124`、`:270-276`）。
- 甲：这一维**加进 `composer` 段**（与 `currentModel`/`modelKnown` 同形，`:255-256` 已是「一枚事实 + 一枚 known 位」的先例），不动 `Snapshot` 顶层。
- 乙：加顶层第五键，按 145 的路子在**同一 commit** 里动那四枚钉（其中两枚是对 `frontend/src/lib/panel.ts` 的双向核对 ⇒ 触发 Q-51）。
- 不做：这一维先不进快照，只走**读方法**的回执（面板点开设置时才问），AC#2 那格因此勾不上。
⇒ 三条都改得到东西，我没资格选。

**J4 凭据值过境的信封形状：AC#0 只禁了「塞进 `panel.message.send`」，但没有仪器拦得住。**
现量：所有已答方法解进**同一枚** `ComposerRequest`（`bridge.go:69-79`），facet 2 的裁决词表（R24）**不含** `value`/`secret`/`config` 等名 ⇒ 一枚带值键一旦进信封，`panel.message.send` 那条路**同样能带值**且解析通过；今天没有任何尺按路由分键。
- 甲：值**不进信封**——写方法只带 `api_key_ref`（`dpapi:<id>`），值另走一次「只写不回显」的一击（形状待定，可能必须新增一枚 `panel.secret.set`，其处理器**立即**交 `Store.Store` 且不留字段）。
- 乙：进共享信封但**按路由拒键**：新增一枚机读判据＝种「`panel.message.send` 带 value」必被拒，且给这判据配正控。
- 不做：靠票 248 §5 那句「界面侧不要把密钥塞进聊天消息」当约束（那是文案，不是闸门）。
⇒ 要个形状；这一条还直接决定 AC#3「不许新增任何把值回显给页面的方法」怎么被核。

**J5 `permission_mode` 的写：面板侧改档今天只有 `panel.mode.request`，而 `wisp panel-inbound` 的 `Confirm` 是**刻意 nil**。**
现量：`panel_inbound.go:213-231`（`Confirm: nil` + 头段 `:32-39` 的理由）、`composer_handlers.go:133-138`（变宽即拒 `ErrNoL2Confirm`）、`run.go:626,805-827`（只有 `wisp run` 装了那张卡）。
- 甲：配置写方法**显式拒写 `risk.*` 全族**（连收紧也拒），档位一律走既有 `panel.mode.request`，一行代码都不给「从设置页放宽」留门。
- 乙：允许写**收紧向**、放宽向交给 `Confirm` 腿；宿主没腿即 `ErrNoL2Confirm`（把 `ModeWriteHandler` 的方向不对称推广到配置段——⚠ 这是我的推广，见②）。
- 不做：让设置页直接写 `risk.permission_mode`（那正面撞 `PLAN.md` D33/D36 的「放宽不得静默生效」与 AGENTS §1.2「面板侧来源的 L2 允许」）。
- 现量补充：票面 §5 第 1 条把「权限档位」列进了那页表单，**这张表长什么样不归本编队**（AC#5），但「档位能不能从这页写」归本票裁。

**J6 一枚通用 `config.set(section, key, value)` 还是按字段数枚具名方法。**
现量：`A484` 的边界＝「只新增配置读/写那一族方法名」；写面既有形状是**按键**的 `mergeWrite(ownedKey, …)`（`writeguard.go:111`，`ownedKey` 是逐枚键路径，`llm.providers.<name>` 在 `:199` 有具名形）；生效级别在代码里是**按段**判的（R22）。
- 甲：读一枚 + 写一枚（键路径 + 值），级别与方向由 config 层判（复用 `plan*`/`riskDirection` 那套），白名单只多两枚名。
- 乙：按字段族数枚（provider 基址/模型名/上下文窗口/价格/`api_key_ref` 各一枚），名册大，但每枚参数字面窄、错误形状具体。
- 不做：本腿继续猜。
⇒ 这决定 `A484` 那张批准记录的射程要不要重写（它写的是「那一族」，不是枚数）。

**J7 重启档与「写完没生效」的回执形状。**
现量：restart 段今天只有 `[app]`（theme 除外）＝`SPEC-03:78` + `manager.go:338 planApp`；那句「本宿主没有接 config.toml 轮询…要生效请重启进程并用 `wisp run`」是 `config_reload.go:97-99` 的 `hotReloadDisabledPanelInbound`，在 `panel_inbound.go:212` 每启动念一次；`Report` 有 `Restart []string`（`manager.go:92-94`）但**没有任何出口把它交给面板**（H10 未建）。
- 甲：写方法的**返回值**里就带 tier 判定（`applied`/`needs-restart`/`denied-loosening`），不新增回显通道。
- 乙：只回执「已写入」，把「要不要重启」全交给下一次 `panel.resync`（依赖票 33/35）。
- 不做：让页面自己按字段名猜哪枚要重启（那是第二套真相）。
⇒ 顺带一问：设置页写完之后**谁来念那句 hotReload 的话**——`panel-inbound` 宿主现在只在 stderr 念，页面上什么都没有（`panel_inbound.go:44-46` 自陈 H10 开着）。

**J8 文档那三张表谁权威、代码扩名后要不要跟着动（本腿一字未改）。**
现量：`docs/specs/README.md` 的优先级（PLAN 定案最高）＋ R17 的 13 处 C17 ＋ R18 那张自署「提案」的 §5.2 表 ＋ R21 的 SPEC-03 缺 `permission_mode`。
- 甲：以 §5.2 那张表为**待批提案**，248 落地时**只**在该表加两行并标批准号（属 `docs/specs/**`，实现腿禁写 ⇒ 谁写？）。
- 乙：承认文档那张表已过期，另开一票做文档同步，248 只动代码。
- 不做：什么都不动 ⇒ ④ 的漂移面再加一层（代码 6 枚 / 文档 20+ 枚 / 提案未定稿三种名册并存）。
- ⛔ 我不改 `PLAN.md`/`docs/specs/**` 一字；这条只要一个裁。
