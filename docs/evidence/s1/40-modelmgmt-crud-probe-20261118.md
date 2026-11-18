# 票 40 AC#3 前两半 — 模型管理 CRUD 真落 config.toml ＋ 探针按钮真调 provider_health

腿：agent-136（写腿，GLM 5.3 Flash 路由）。日期：2026-11-18。
派单：`.scratch/wisp/probes/orch/20261118/txt/40-modelmgmt-a-dispatch.txt`（编排者已随 `005f545fa` 入库）。

## 交付

### 前提盘点（现读，不按票面旧行号）

CRUD 的**写腿**（row create/delete、context_window、price、audio_in/out、base_url、billing、role pair）是票 320/313/187 交付的现成面：面板门 `internal/panel/config_handlers.go`（白名单＋refused-drift）→ cmd/wisp `panel_config_store.go ApplySetting` → `internal/config` 各 setter（受守写路径 mergeWrite，`writeguard.go`）。**缺的三样**＝①模型行 CRUD 的"读回逐字段相等"牙（320 测试只断言 contains）②探针按钮（页面 `测试连接` 是诚实桩 `COPY.testConnectionUnavailable`，provider_health 机器在库但面板零路由）③endpoint 类写与探针收据的档位语义钉。

### 新路由：`config.model.probe`（写面 internal/panel ＋ cmd/wisp）

- `internal/panel/bridge.go`：名册第 21 名 `MethodModelProbe = "config.model.probe"`（无 `panel.` 前缀、无审批词表、不走 `models.*` 命名空间——config.models.list 先例）＋ `knownComposerMethod` 守卫行。
- `internal/panel/composer_dispatch.go`：`ModelProbeHandler` 接口＋`ModelProbes` socket＋派发 case（门数 20→21）。
- `internal/panel/probe_40.go`：`ModelProbePayloadMarker = "wisp-model-probe:"`、`ModelProbeResult/Outcome/Mismatch` 线形状、`RenderModelProbeReceipt/Payload`（回执渲染单源：p50 延迟、逐能力 声明/实测 配对词、mismatch 冻结句 Label()＋detail 原文、provider_health 落库半句）。**回执零档位句**——探测只测不存，档位词（需要重启/立即生效）出现在探测回执上就是反向错误档谎言。
- `cmd/wisp/panel_probe_40.go`：真机器——`llm.NewResolver(cfg, h.secrets)`（键从凭据存储走，不可解析的 ref 按名字拒绝，⛔无键裸测）→ `llm.BuildEndpointProvider(eps[0], chainOptions(cfg))` → `llm.RunProbeSuite`（Sink＝`storeHealthSink{store}`，判判决真写 provider_health，ticket 09 sink）→ 回执＋载荷行。存储按次开关（cost/privacy handler 同形）。结构失败（Unconfigured/provider 构造/存储）走 Handle 的 err 路＝RefusedEnvelopeForUser 原文；套件跑不完走载荷 `error` 字段原文。
- `cmd/wisp/panel_inbound.go`：装配 `newModelProbeHandler40(dir, mgr, secrets, auditf, actor)`（与设置腿同一 Manager、同一凭据存储——测的就是页面显示的那一行）＋ `ModelProbes: probes`。
- 名册尺：`inbound_roster_253_test.go` 名册＋`wantFullInboundRosterSize253=21`＋`wantNonPanelPrefixedInbound253=17`＋两处 plant guardLine 逐字副本同步。

### 前端（frontend/src/** ＋ scripts/render-*）

- `panel.ts`：`requestModelProbe(provider, model)`（route literal 唯一落点门文件）。
- `config-view.ts`：`MODEL_PROBE_MARKER`＋`ModelProbeResultView/OutcomeView/MismatchView`＋`parseModelProbeReply`（缺 identity 字段/非数字延迟 ⇒ no-payload 显示 prose，⛔画零冒充实测；payloadError 具名态）。
- `config-actions.ts`：`performModelProbe` → `ModelProbeRun`（no-host＝「未连接到本机服务，探测没有发生」；拒绝 prose 逐字上屏）。
- `model-settings.tsx`：删供应商行诚实桩按钮「测试连接」（`COPY.testConnectionUnavailable` 随之删除）；模型目录每行加「测试」按钮 → 真探针；`ModelProbeResultCard`（导出纯件，仪器可静态挂载）：p50 数字、逐能力词（声明可用/未声明/实测可用/实测不可用/未实测——全词，零 glyph）、失败 detail 原文、mismatch 冻结句、`判定已写入 provider_health；探测不改配置，不落 config.toml`。探测中态写明可能要几十秒。
- `fixtures/settings-v2-preview.tsx`：四处挂载穿 `performModelProbe` 真路径（与产品同门）。
- `render-model-settings.tsx`：新增 tooth 7 共 12 枚（真载荷逐字、拒句转义逐字、无档位词、桩句消失、route literal 只在 panel.ts、页面不判第二遍）。

## 牙与读数（cmd/wisp，载具 `scripts/wisp-cli-staged-tests.sh`）

| 牙 | 读数 |
|---|---|
| `Test40ModelMgmtCRUDRoundTripFromDiskAndBack` | PASS 0.07s——增（CreateModelRow）→**fresh LoadFile 读盘逐字段**（display/context_window/max_output/thinking_levels/billing/price/quota/enabled/capabilities 九维，billing=pay-per-token＋enabled=true 两枚 default 标签）；改（context_window 64000、audio_in true）→ 读盘逐字段相等＋字段隔离；读回脸 config.get payload 列行＋caps 位；删 → 读盘无此条且 seed 行仍在、读脸不再列 |
| `Test40ModelMgmtEndpointWriteSaysRestartVerbatim` | PASS 0.02s——`provider_base_url`（endpoint 类）回执逐字含 39 AC#1 全句「这一项要重启进程并重新运行才算用上（同一次运行里模型通路不会重建）。需要重启：本次运行不会读取新值。」，零 立即生效/保存即生效 |
| `Test40ModelProbeMeasuresTheRealProviderAndWritesProviderHealth` | PASS 3.83s——活 mockllm：载荷 provider/model 逐字、written=true、latencyMS 数值、逐能力 attempted+ok；回执含 实测 p50；回执零 需要重启/立即生效/设置已写入；provider_health 真行（fc=true） |
| `Test40ModelProbeReportsTheMismatchAndTheRawDetail` | PASS 2.63s——server thinking=broken：payload mismatches 含 thinking＋label/detail 非空；回执含 能力实测不符＋声明可用/实测不可用配对 |
| `Test40ModelProbeFailureCarriesTheRealErrorNeverOK` | PASS 59.62s——死端点（127.0.0.1:1）：全能力 attempted+fail、detail 非空且逐一致（真原文）、回执逐字回显 detail＋实测不可用（细节：；allOK=false；⛔捏造可用 |
| `Test40ModelProbeUnconfiguredRefIsNotSilentlyKeyless` | PASS 2.22s——dpapi:absent：拒绝句带 Unconfigured 原文、无载荷行、零捏造判定 |
| `Test40ModelProbeRefusesByNameAndWiringNail` | PASS 2.28s——未知 provider/model 按名拒绝（没有服务商/没有模型）；源级钉：panel_inbound.go 的 `newModelProbeHandler40(dir, mgr, secrets, auditf, actor)`＋`ModelProbes: probes`、panel_probe_40.go 的 `llm.RunProbeSuite(`＋`storeHealthSink{store: store}`＋`llm.NewResolver(cfg, h.secrets)`——摘任一线即红 |

### 变异证明（仓内原地＋pristine 备份 `.scratch/wisp/probes/orch/20261118/40-modelmgmt/`，测毕逐字节还原）

- **M1 只改内存不落盘⇒红**：`internal/config/writeguard.go` mergeWrite 的 `SaveFile` 调用套 `if false {}` → `Test40ModelMgmtCRUDRoundTripFromDiskAndBack` 红：「the created entry is not on disk」。还原后绿。
- **M2 写盘不回读⇒红**：`cmd/wisp/panel_config_store.go` ReadSettings 的 `row.Models = nil` → 同测试红：「the read face never lists the disk row model-under-test-40」（payload models:null 被点名）。还原后绿。
- **M3 空转返回OK⇒红**：`cmd/wisp/panel_probe_40.go` 载荷 `OK: true` 硬编 → 死端点牙红（三能力「measured ok:true want attempted+fail」＋「receipt claims a capability worked」）。还原后绿。

## 门禁读数（最终数以回话为准）

- `sh scripts/d22scan.sh` → **clean**（末行 `d22scan.sh: OK - control rc=0, scan rc=0`；ban#8 frontend/ 184 text files）。
- `go build ./...` **rc=0**。
- `go test ./internal/panel` **ok 27.163s**（253 名册尺＋page 双向门含新名）；`go test ./internal/config` ok 2.852s。
- 载具本包全量第一发：**2 红**——①`TestTicket255RosterStillMatchesTheActualReadSites`＝**本腿自己的新读口被尺逮住**（`cmd/wisp/panel_probe_40.go` 读 `[llm]`，255 名册未重裁）：按 441 同款先例把该文件登进 `config_readers_255.go` 的 sectionReadSites[LLM] 行＋证据句（live Manager at use time，与 panel_models_320 同类），定向复跑 `^TestTicket255Roster` **3/3 PASS rc=0**；②`TestAC4FocusReturnToPriorWindowGap33r5`＝**既有红，非本腿**：441 腿同日更早的全量（`.scratch/wisp/probes/orch/20261118/441-vehicle-full-run.txt`）同枚已红，且其 commit-msg 收据点名它来自 7f3177f93/43bf68b51 自己的收据（planted-thread WM_CLOSE 活窗面）；本腿四处生产改动（probe handler/dispatch/名册/读口行）无一涉及窗线程。修后收口第二发全量读数见回话。
- 前端：`typecheck` rc=0／`lint` rc=0（45 warnings 全存量，0 error；新代码零命中）／`build` rc=0／`render:model-settings` rc=0（12 枚新探针牙＋存量 6 组全绿）／`render:security40`、`render:privacy40`、`render:nav`、`render:composer`、`render:cfggolden` 全 rc=0。
- dist 产物 git-ignored，未入库；⛔ 未 push；未动 thresholds/golden/allowlist/PLAN/specs/D43/design。

## 诚实边界

- `测试密钥` 桩未动（不在本单射程；census B 类诚实桩保留）。
- quota 范围/capability 过滤下拉/chain 拖拽序：本单不做（派单明文划出）。
- 桩删除后供应商行不再有「测试连接」：探针按票 40 文本落**每模型一行**（provider_health 主键＝provider/model 对），供应商级按钮没有可测模型时不发明选择。
- 变异三枚在仓内原地做（pristine 备份在 probes 目录，git diff 复核零残留）——非仓外拷贝，形态与 256/263 腿先例一致。
- 探针死端点用例单发 ~60s（套件自身重试/退避，非本腿新增机制）；全量载具读数含它。
- `frontend/fixtures/composer-states.html` 在飞件（他人 M）：未触碰。
