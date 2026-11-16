# 票 313 · 云端语音链装配腿（313-wire）交付证据 — 2026-11-16

**腿名**：agent-79（写腿，派单存档 `.scratch/wisp/probes/orch/20261116/txt/313-wire-voice-cards-dispatch.txt`）
**范围**：`cmd/wisp/**`（本波唯一 cmd/wisp 写腿）＋ `frontend/src/components/voice-settings.tsx`（＋仪器 `frontend/scripts/render-model-settings.tsx` 小改）＋ 本文件。`internal/speech/**` 零字节（只读）。

## 0. 现量复核（改前重跑）

- `grep -rn "NewCloudASRChain|NewCloudTTSChain" --include=*.go internal cmd`：除 `internal/speech/cloudchain.go` 定义与 `internal/speech/cloud61_test.go` 测试外 **全仓调用点 0** —— 云端链从未装配，四张卡的「占位·当前不生效」是真话。
- 设置写入路由已在树上（票 248/313-r1/r2/r3：`bridge.go:66-67` config.get/set、`panel_inbound.go:321` DPAPI、`panel_config_store.go:283-286` 两条链的 setter）⇒ 配置能写、没人拿配置去 new 链。

## 1. 生产调用点（AC#2 的运行那一半）

落点（本腿 diff 里的行号）：

| 职责 | 位置 |
|---|---|
| 链解析（boot 半）：`resolveVoiceCloudChains(dataDir, c)` —— `secret.NewStore`（DPAPI）作 KeyResolver，`speech.ResolveCloudChain` 解析 `voice.cloud_asr_chain` / `voice.cloud_tts_chain`；解不出＝具名 note＋该方向空表（绝不静默、绝不 boot 失败） | `cmd/wisp/voice_cloudchain_313_windows.go`（新件） |
| 链装配（loader 半）：`buildCloudAsrEngine` / `buildCloudTtsEngine` —— `speech.NewCloudASRChain` / `speech.NewCloudTTSChain`，本地 sherpa/matcha 作 terminal fallback（`voiceCloudAsrFamily` / `voiceCloudTtsFamily` 把链与本地回退的 Close 并成 slot 的一次 Close，堵住「链不 own fallback」的泄漏面） | 同上 |
| 装配进既有入口：`voiceResidency.loader` 的 EngineASR / EngineTTS 分支在**原有本地路径之前**加云端分支（空链走原路径，字节不变）；`voicePlan` 新增 `CloudASREps/CloudTTSEps` | `cmd/wisp/resident_voice_engines_windows.go` |
| boot 接线：`buildResidentVoice` 在 `resolveVoicePlan` 之后解析两链；verdict 前缀保持 `语音链已装配`，配置了链才追加 `；云端听写 stepfun/...；云端朗读 stepfun/...`；refusal 经 `rv.noteProblem` 具名 | `cmd/wisp/resident_voice_windows.go`（buildResidentVoice 内 resolveVoiceCloudChains 调用＝boot 生产调用点） |
| 协程注册表：`startResidentVoice` 在 `rv.start` 之前 `rv.residency.attachCloudRegistry(rt.Registry)` —— 云 ASR 的 per-utterance transport 协程经 `observe.Registry.Spawn` 起名（roster 一张图） | `cmd/wisp/resident_voice_windows.go` |

**禁区自证**：⛔ 新配置档位（四档不变，`voice.cloud_*_chain` 是既有 TierReload 键，`internal/config/tiers.go:83-84`）；⛔ 动 `thresholds.go`/golden/allowlist/PLAN/specs/D43/design；⛔ 新协程（全仓无裸 `go func`，transport 协程是 internal/speech 既有形状）；C17 白名单零动；`go.mod` 零动。

**降级语义（写死并钉住）**：链元素解不出（未知 pair / enabled=false / 缺 `capabilities.audio_in|audio_out` / api_key_ref 解不出）⇒ 该方向整链不装配，problems 具名元素与原因，boot verdict **不**点云端头，本地级联照常 —— 与 cascade 自己的 terminal-fallback 设计一致；verdict 与 loader 永不说两样话。

**已知边界（照实说，不藏着）**：
1. 装配吃的是 boot `config.LoadFile` 的一发读（与 voicePlan 每个字段同一根馈线）；`cloud_*_chain` 是 TierReload ⇒ **保存后重启生效**，这正是前端文案写的那句话。
2. QuotaGate / UsageSink 未接（级联 nil 闸＝「未配 quota」，nil sink＝「不计量」）：C23 的语音行记账归它自己的票，本腿不假装接了。
3. 面板建行（`CreateProviderRow`）产出的模型 capability 位全 false ⇒ `speech.ResolveCloudChain` 会拒绝它做链元素（SPEC-03 §3.1 的 probe-first 契约）——这不是本腿造的缺口，是目录声明位本身的形状；用户在 config.toml 里声明 `audio_in/audio_out=true`（或后续票给页面加 capability 位）后链即装配。降级形态已钉（见 §2 牙三）。

## 2. 行为的牙（`cmd/wisp/voice_cloudchain_313_windows_test.go`，三枚，载具现跑全绿）

| 测试 | 钉什么 |
|---|---|
| `Test313CloudVoiceChainAssemblesIntoTheResidentLoaders` | 配置 stepfun 行＋两链 ⇒ buildResidentVoice verdict 点名两链；residency.loader(EngineASR/EngineTTS) 真返回链族且 `Active()`＝配置元素逐字。**拆调用点红句**：删 `resolveVoiceCloudChains` 调用或 loader 云端分支 ⇒ `拆掉调用点会红：ASR 族里装的不是云端链（got %T）`（loader 落回本地工厂，该夹具形态下先报 `拆掉调用点会红（loader 落回本地工厂）：ASR loader refused...`） |
| `Test313CloudVoiceChainFollowsTheConfig` | 改配置元素 ⇒ 新一次生产 build 的链头跟着变（`stepfun/model-under-test-313-chain-asr-b`）——「改配置⇒装配跟着变」的另半 |
| `Test313UnresolvableCloudChainSaysSoAndKeepsTheLocalCascade` | 缺 capability 位的链元素 ⇒ verdict 不点云端听写、problems 具名 `云端听写链未装配（... element "stepfun/..." lacks capabilities.audio_in ...）`、boot 不倒 |

**载具现跑**（`D:\work\soft\Git\bin\bash.exe scripts/wisp-cli-staged-tests.sh -test.run 'Test313CloudVoiceChain|Test313UnresolvableCloudChain'`）：三枚 PASS，`ok github.com/CarlosShao/wisp/cmd/wisp 6.00s`，rc=0（2026-11-16 11:32）。

## 3. AC#1（默认供应商 stepfun，不写死）

- 常量复用 `model-settings.tsx` 的 `OWNER_DEFAULT_PROVIDER`（一处定义），语音页经 `voiceDefaultProvider(providers)` 使用：**门槛＝这行供应商真的存在**（链元素必须点名已建行的 provider/model，给不存在行的预选＝必然保存失败的坑）。
- 用法只落在「未设置字段的草稿初值」：链为空 ⇒ 第一段供应商预选 stepfun；realtime 未配置 ⇒ 草稿 stepfun（`view.voice.realtimeProvider || voiceDefaultProvider(providers)`）。全部 **useState 初始值，无 effect 强加，不点保存不写盘**；用户存过的值逐字显示、永不被默认盖写。
- 仪器牙（`render-model-settings.tsx` tooth 5，npm run render:model-settings rc=0）：fresh 页 3 枚 picker 预选 stepfun；源级钉 row-exists 门槛／initializer-only／saved-value-wins；VOICED 页 `stepfun / 示例模型-2` 逐字在（用户值赢）。
- 模型页的 AC#1 承载物（ProvidersCard 两道门预选）由 313-page 腿已落（票面 :48），本腿不重复。

## 4. 前端文案跟真相同步（唯一动 frontend 处）

**改了的**（`frontend/src/components/voice-settings.tsx`）：
- lede（原 `:81`）：「三张卡都是占位…当前不生效」→「云端听写链与云端朗读链已接入常驻语音的装配（保存后重启生效，链尽回退本地引擎）；实时语音的引擎尚未实现仍是占位；本机音频的后端读者已接真，面板先等快照。」
- 听写卡 description（原 `:144`）：「占位卡：…还没有接到听写运行路径，接入前改它不会生效。」→「已接入：常驻语音启动时按这条链的顺序装配云端听写，链尽回退本地引擎；保存后重启生效。」
- 朗读卡 description（原 `:154`）：同理 →「已接入：…装配云端朗读，链尽回退本地引擎；保存后重启生效。」
- 两张链卡的 StateTag：`占位·当前不生效`（warn）→ `已接入·重启生效`（ok）；链卡保存回执不再降级为 warn（真接了，`placeholderReceipt` 从链卡 save 移除）。

**没改的**（照 v2 的 B 级诚实Stub 原样保留）：实时语音卡（`占位·当前不生效` warn ＋ RealtimeEngine Reserved 的占位描述 ＋ warn 回执）；本机音频卡（等快照态）。⛔ 把假话换成另一种假话——没通的卡一个字没动。

**门禁**：`npm --prefix frontend run typecheck` rc=0；`npm run build` rc=0（产物 `index-DhVAo5BL.js`）；`npm run lint` 0 errors（36 warnings 皆既有类）；`npm run render:model-settings` rc=0（五牙全绿）；`tools/d22scan` clean（frontend/=123 文件零 emoji）；**`frontend/dist` 已删回 `.gitkeep`（构建字节零提交，票 77）**。

## 5. AC#3 复核（key 面）

- 值面走既有 `api_key_ref`（D36 规则 5）：`speech.ResolveCloudChain` 经 `secret.Store`（DPAPI）解析，明文只活在 `speech.CloudEndpoint.APIKey` 内部；本腿打印的每一个字符串（verdict/problems/slog）都是 ChainKey/baseURL/事件词汇形——`VoiceFailoverEvent.Detail` 与 `CloudEndpoint.String` 的契约本身禁 key。
- ⛔ 新造重复件：复跑现成哨兵 `TestAC2CredentialSentinelAppearsInNoArtifact`（`cmd/wisp/panel_config_248_test.go:154`），含在整包 run 内，读数见 §6。

## 6. 整包门禁读数

- `go build ./...` rc=0；`go vet ./cmd/wisp` rc=0；`gofmt` 空；`tools/d22scan` clean（ban #6/#8 frontend=123 文件）。
- **载具整包**（`scripts/wisp-cli-staged-tests.sh`，full-fixture 模式「voice fixture environment resolved」，2026-11-16 11:39–11:59，1202s，log＝`.scratch/wisp/probes/313/full-run-20261116.log`）：
  - `=== RUN` 726+，终点锚 `TestZZ315A1` **三案全命中且 PASS**（RootLoopFirstBridgeCall 0.27s / SiblingChildFirstBridgeCall 0.28s / ProjectInstructionsSeedCannotStrandTheFirstTask 0.23s）；最后一条 `=== RUN`＝`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`。
  - `TestAC2CredentialSentinelAppearsInNoArtifact` **PASS**（0.03s；单独复跑见下）。
  - FAIL 三枚，处置如下：
    1. `TestTicket255RosterStillMatchesTheActualReadSites` —— 名册警报（设计内要求 re-adjudicate）：①`[voice]` 的 reader 扫描点名了新文件 `voice_cloudchain_313_windows.go` 而 `sectionReadSites` 未列 ⇒ 已补（Voice 行）；②`voice.tts.speed` 行的 no-reader 前提「ResolveCloudChain 全仓非测试调用者 0 枚」被本腿正当推翻 ⇒ 按仓内 re-adjudication 惯例重裁为 **other-process**（boot 装配调用、随重启生效，仍非 consumed——本 host 无活读者），钉该前提的 `TestTicket255RosterReJudgedAfterTickets44And61` 同步重裁（钉改的是「前提」不是松断言：新红句要求点名新调用者 `voice_cloudchain_313_windows.go`，删装配点它照样红）；③顺带刷新本腿搬动的 `resident_voice_windows` 行号读数（1848→1968、1862→1982、1866/1867/1868→1986/1987/1988）。**重裁后 `TestTicket255*` 全族 22 枚载具现跑全绿（rc=0，20.00s）**。其余漂移行（run.go/panel_host/logsink 等，本腿未碰的文件）是 313-r2 已挂账的 255-r2 存量，只记 Logf 不红，未越界代裁。
    2. `TestAC4FocusReturnToPriorWindowGap33r5` —— **桌面态红**：panel Show 时前台被别的窗口占（红句自带读数 foreground 0x94e0122 vs panel hwnd 0x52705c6；机主桌面实时状态）。本腿零窗口代码；A1893 单跑复量同因。
    3. `TestAC246DevLegIgnoresTheTestTaskInjection` —— **桌面态红**：红句逐字「a dev Wisp is already running in this session… Stop it and re-run」——机主的 Wisp 球（PID 49348）正在跑，单实例互斥让该腿按设计交还激活。未杀机主进程（危险动作不做）；A1893 单跑复量同因。
  - 其余 SKIP 全部是测试自身的活机/WebView2 门（非载具名册跳过——本次是 full-fixture 模式）。
- **载具单跑读数**（新牙＋名册＋哨兵，rc=0，20.00s）：三枚 313 牙 PASS（6.64s/0.01s/0.01s）＋ `TestTicket255*` 22 枚 PASS ＋ `TestAC2CredentialSentinelAppearsInNoArtifact` PASS。

## 6.1 AC 勾选（终版）

- **AC#2**：☑（运行半＝本腿装配点＋三牙；读写面前腿已落）。
- **AC#1**：☑（模型页前腿＋语音卡草稿预选＋仪器 tooth5）。
- **AC#3**：☑（哨兵 PASS；key 只活在 endpoint 内部；页面半 313-page 已落）。
