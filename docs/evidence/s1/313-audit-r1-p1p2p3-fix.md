# 313 审计 R1/P1/P2/P3 修正（写腿，2026-11-16）

> 派单：`.scratch/wisp/probes/orch/20261116/txt/313-audit-r1-p1p2p3-dispatch.txt`。
> 审计判「成立 5／不成立 4」，本单修 4 条：P3 能力位可达（读面＋写面＋前端）、R1 两卡文案条件式、
> P1 255 名册补删装配点必红的牙、P2 SLO 附录 E 补取样态。⛔ 反驳审计、⛔ 只改文案不改机制。

## 0. 根因（审计已复读源码，本腿照单执行，未反驳）

写面只查 provider/model 存在（`internal/config/settings_313.go` requireChainElements →
`internal/config/catalog.go` checkChainElement，无能力位检查）；装配面查能力位
（`internal/speech/cloudcatalog.go:97/:101` lacks capabilities.audio_in/audio_out）；
能力位默认全 false（`internal/config/schema.go` ModelSpec probe-first）；读面
`cmd/wisp/panel_config_store.go` VoiceSettingsView 无能力位、写面也无。
⇒ 页面写「已接入·重启生效」，但新用户无处声明 audio_in/audio_out ⇒ boot 时链必被拒＝反向假话。

## 1.【P3】能力位可达 —— 三面各一个落点

- **写面（config setter）**：`internal/config/settings_313.go:207` `SetModelAudioIn` /
  `:218` `SetModelAudioOut`（共同底座 `:220` setModelAudioBit）。每枚一叶一键，走
  `writeOneKey`（票 248 受保护 setter 同一条路，⛔ 未另造新路；非密钥，不进 DPAPI），
  键路径 `llm.providers.<p>.models.<m>.capabilities.audio_in|audio_out`，落盘前
  `requireCatalogEntry`（拒因 2：行不存在）。能力位**不**加进链写面的拒因——链保存仍只查
  存在性，装配结果由读面＋前端如实呈现（派单判据就是这么定的）。
- **写面（面板名册）**：`internal/panel/config_write_313.go:111-112` `FieldModelAudioIn`
  =`model_audio_in` / `FieldModelAudioOut`=`model_audio_out`（313 名册一族，needsProvider+
  needsModel+takesValue 三列都进）；接线 `cmd/wisp/panel_config_store.go:340` ApplySetting
  两个 case（布尔形状话在本腿答，行存在尺在 config 侧——与 price 叶同一分工）。
- **读面**：`internal/panel/config_handlers.go:291` `ModelCapabilityBits`（AudioIn/AudioOut）、
  `ProviderSettings.ModelCapabilities`（按 model id 的 map，nil＝腿没答）、wire
  `providerPayload.modelCapabilities`（`omitempty`：没答＝键不上线，与语音 payload 指针门同
  一类）；填充 `cmd/wisp/panel_config_store.go` ReadSettings（每个模型两位，从同一批 specs 读）。
- **前端（模型行两枚开关）**：`frontend/src/components/model-settings.tsx` 模型目录每行
  两枚 Switch（audio in / audio out），走 `model_audio_in` / `model_audio_out` 一写一位；
  能力位快照没答＝「等快照」标签，不是猜的 false。类型
  `frontend/src/lib/config-view.ts` ConfigViewProvider.modelCapabilities。
- **payload 键 pin 更新（非放宽）**：`internal/panel/config_313_test.go`
  Test313PayloadKeySetIsFrozenAndCarriesNoValueBearingKey 的 wantRow 键集**仍是逐字节
  exact-match**，只是按派单明令的读面新增把 `modelCapabilities` 一枚键加进去（测试内注释
  写明这是点名新增的 re-adjudication，不是松断言）；value-bearing 键扫描照旧跑在增后的键集上。

## 2.【行为判据】默认 false ⇒ 装配拒绝；打开 ⇒ 装配成功（改位 ⇒ 装配结果跟着变）

`cmd/wisp/panel_modelcaps_313_windows_test.go`（windows 载具面，走真 settings 腿）：

- `Test313CapabilityBitFlipsAssemblyOutcome`（:73）：
  1. 种子 config（无密钥 provider + 一条 probe-first 全 false 模型 + 两条链指着它）⇒
     `speech.ResolveCloudChain` 按名拒绝，红句原文 `lacks capabilities.audio_in`
     （TTS 同款 `lacks capabilities.audio_out`）；
  2. 读面回答默认 false/false（真答案，不是缺键）；
  3. 经 `configStore.ApplySetting(model_audio_in=true)` 走票 248 门落盘，回执点名
     `llm.providers.stepfun.models.<id>.capabilities.audio_in`、档位 restart；
     reload 后**同一条链装配成功**（端点 head＝链元素）；audio_out 方向仍拒（不夹带）；
  4. 读面跟着翻（AudioIn true / AudioOut false）；
  5. TTS 孪生半：`model_audio_out=true` ⇒ 朗读链装配成功。
- `Test313CapabilityWritesRefuseByNameAndMoveNoByte`：点名不存在模型 ⇒ 拒因 2「行不存在」；
  非布尔词 ⇒ 形状话；两种拒绝 config.toml 字节不变。

## 3.【R1】两卡文案条件式＋点名门槛（frontend/src/components/voice-settings.tsx）

- **听写（ASR）卡 description 原文**：
  「已接入：常驻语音启动时按这条链的顺序装配云端听写，链尽回退本地引擎。门槛：链里每个模型都要在能力位里勾上 audio in（「模型」页的模型行），保存后重启生效。」
- **朗读（TTS）卡 description 原文**：
  「已接入：常驻语音启动时按这条链的顺序装配云端朗读，链尽回退本地引擎。门槛：链里每个模型都要在能力位里勾上 audio out（「模型」页的模型行），保存后重启生效。」
- **页头 lede** 同步点名门槛（「……门槛：链里每个模型都要在能力位里勾上对应的 audio in / audio out（「模型」页的模型行），保存后重启生效……」）。
- **卡顶标签不再是静态「已接入·重启生效」**：按已存链对读面判——能装配＝「已接入·重启生效」
  (ok)；缺位＝「装配不了·缺能力位」(warn) 并点名缺哪一位（`stepfun/x 缺
  capabilities.audio_in`，与 boot 拒绝句同一词汇）；快照没答＝「已接入·装配未判（等快照）」；
  链为空＝「未设置」。判定函数 `chainAssembly`（本文件），判定词表与 cloudcatalog.go 同源。
- **保存回执分流**（「这链现在能不能装配」）：能＝`已接入：保存后重启生效`(ok)；
  不能＝error 醒目「已写入配置文件，但这条链现在装配不了：<点名缺哪一位>。在「模型」页把缺的那
  一位打开并保存，再回来存这条链，重启生效。」；判不了＝warn「已保存，重启后生效；这条链现在能
  不能装配判不了（能力位快照还没读到）。」

## 4.【P1】255 名册补删装配点必红的牙

`cmd/wisp/config_receipt_255_test.go:970` `TestTicket255CloudChainAssemblyCallPointIsAlive`
——**按真调用图判，不断判决字符串**（go/parser 走 cmd/wisp 生产源码，非测试文件）：

1. 生产 cmd/wisp 必须存在至少一处 `speech.ResolveCloudChain` 调用（装配点本身）；
2. 持有该调用的函数必须仍被至少一处其他生产函数调用。

删 `resolveVoiceCloudChains` 里的 ResolveCloudChain 调用 ⇒ 尺 1 红句
「拆掉装配点会红：production cmd/wisp no longer calls speech.ResolveCloudChain anywhere…」；
删 `buildResidentVoice` 里对 resolveVoiceCloudChains 的调用 ⇒ 尺 2 红句
「…nothing in production cmd/wisp calls <fn> any more - the assembly point is orphaned…」。
行为半边本来就有（`voice_cloudchain_313_windows_test.go` 的
Test313CloudVoiceChainAssemblesIntoTheResidentLoaders / ...FollowsTheConfig 驱动
buildResidentVoice 端到端）；本牙是名册自己的跨平台半边，两半拼齐「删装配点必红」。

## 5.【P2】SLO 附录 E 取样态（叠加不改写，仅附录）

`docs/SLO.md` 附录 E 末尾补三句：① 本发是负载态（审计实测 CPU 57%＋Weixin /
magpie-windows-amd64 / DeepSeek Harness 在跑，非附录 A 的安静窗口口径）；② 同机同态对照
（33-r4 安静态 815.826ms vs 本发负载态 1693.016ms＝2.07 倍差，同机不同态只记倍差不判预算）；
③ ⛔ 本发不得与安静态 815.826 合并成一个冷启 P95（合池须同态重测后由验收腿明令）。阈值零改动。

## 6. 构建与测试

- `npm --prefix frontend run typecheck` rc=0；`npm --prefix frontend run build` rc=0
  （node 24.21.0 @ D:\work\server\node24）；构建产物按硬约束删回 `frontend/dist/.gitkeep`。
- `go vet` ./internal/config ./internal/panel ./cmd/wisp rc=0；
  `go test ./internal/config ./internal/panel` rc=0（313 面板八测全过）。
- `go test ./cmd/wisp` 只经载具 `scripts/wisp-cli-staged-tests.sh` 一次一发（bash.exe
  D:\work\soft\Git\bin\bash.exe，跑前 Get-Process wisp.test＝0），结果见 §7。
- ⚠ Wisp 球 PID 49348 在跑：TestAC4Focus/TestAC246 形的环境红按环境报，不算本腿弄坏。

## 7. 载具整包读数（A1882 口径）

- **载具**：`scripts/wisp-cli-staged-tests.sh`（经 `D:\work\soft\Git\bin\bash.exe`，一次一发；
  两次跑前 `Get-Process wisp.test` 计数均为 0）。日志留档
  `.scratch/wisp/probes/orch/20261116/txt/313-audit-r1-vehicle-run.log`。
- **如实记一笔**：第一次载具跑是 setup failure（`go test -c` 编译错：新测试一个未使用的局部
  变量，未跑任何测试），修正后重跑；下面读数是第二次（成功跑到包结果行）。
- **包结果**：`FAIL github.com/CarlosShao/wisp/cmd/wisp 1148.00s`，rc=1；**全部 3 枚红都是
  派单点名的环境红**（Wisp 球 PID 49348 在跑，别当本腿弄坏）：
  `TestAC4FocusReturnToPriorWindowGap33r5`、`TestAC4PriorFocusSurvivesARefusedPanelSample`
  （AC4 同族，前台 HWND 被常驻球顶掉）、`TestAC246DevLegIgnoresTheTestTaskInjection`
  （stderr「another instance is running in this session; activated it; exiting」＝单实例碰撞）。
- **终点锚（A1882）**：`TestZZ315A1` 三枚全命中且全 PASS（RootLoop / SiblingChild /
  ProjectInstructionsSeed）；**最后一条 === RUN**＝`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`，
  日志以包结果行收尾（8032 行，无截断、无无声断掉）。
- **新测试位置**：在截断点（包尾）**之前**：`TestTicket255CloudChainAssemblyCallPointIsAlive`
  PASS（log:1282）、`Test313CapabilityBitFlipsAssemblyOutcome` PASS（log:3149）、
  `Test313CapabilityWritesRefuseByNameAndMoveNoByte` PASS（log:3153）。
- SKIP 均为载具名册的有名跳过（fixtures/活麦克风/活窗口族），FIXTURE-GAP＝0（无名册缺口）。
