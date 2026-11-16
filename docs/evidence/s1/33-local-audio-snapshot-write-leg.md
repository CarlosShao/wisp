# 33 本机音频卡接真（快照当前值＋写入路由＋前端三行）— 本腿交付表

日期：2026-11-16。派单：`.scratch/wisp/probes/orch/20261116/txt/33-audio-snapshot-dispatch.txt`。
起手锚（HEAD）：`e49962d55`。写面：`internal/panel/**`（一处拒绝文案）、`frontend/src/**`（本机音频卡三行）＋本证据件。`cmd/wisp/**` 与 `internal/config/**` **零写面**——查证后发现主目标的 Go 侧两半（快照当前值＋写入路由）**已由 2026-11-11 的 leg R 落地**（`d64d48e03`），本腿只做接真核验＋前端接线，未重造任何宿主或写入路径。

## 1. 主目标逐项现量（起手已真，本腿核验）

- **①出向快照带当前值**：`internal/panel/config_handlers.go` 的 `voicePayload`（`aecEnabled`/`inputDevice`/`sampleRate`/`aecBuildAvailable`，reader-gated 指针，nil＝键不上线）；生产读口＝`cmd/wisp/panel_config_store.go:103-137` 的 `ReadSettings`（每次现读 `mgr.Config()`，`aecBuildAvailable` 来自 `aec.BindingBuilt()` 编译事实）。**不是"只出已接真标签不出值"**——值本体（如 `sampleRate:16000`、`inputDevice:"Blue Yeti"`）上线，且被字节钉住（`internal/panel/config_313_test.go:390` 的 `wantVoiceJSON`）。
- **②写入走票 248 既有 `config.set` 受保护 setter**：`config_write_313.go` 名册两枚平面字段 `voice_aec_enabled`/`audio_input_device` → 门（`config_handlers.go write()`）→ 腿（`cmd/wisp/panel_config_store.go:291-309`）→ `internal/config/settings_local_audio.go` 的 `SetVoiceAecEnabled`/`SetAudioInputDevice`（`writeOneKey` 受守写路径，落盘前 validate）。**零新路**。
- **③ sample_rate 冲突如实拒绝＋具名原因**：`audio.sample_rate` **没有也不设**写入字段（`Test313SampleRateStaysAGuardKeyNotAWriteField` 钉住）；任何对它的 `config.set`（点路径 `audio.sample_rate` 或平面名 `audio_sample_rate`）都被门按 `ErrConfigFieldRefused` 拒绝、不落盘。本腿补的一处：点路径 `audio.` 族的拒绝句原来套用"安全锁定族"的通用话（对 audio 是谎话——`audio_input_device` 是可写的，sample_rate 是守卫不是安全键），改为具名真话（`internal/panel/config_handlers.go`，+8 行）：点出可写的平面字段名与 16000 加载期守卫的拒绝单源。
- **前端三行**：见 §2。

**本腿 Go 侧测试读数（载具一次一发，`scripts/wisp-cli-staged-tests.sh`，起手 `Get-Process wisp.test`＝0 枚，fixtures 全在＝全分母零 skip）**：

```
--- PASS: TestConfigLocalAudioSnapshotKeysRideTheLiveReader (0.01s)   ← 快照键 ride 现读（写配置文件后下一发读数即变）
--- PASS: TestConfigLocalAudioWriteFieldsRoundTripAndRefuse (0.03s)   ← 两写入字段往返＋坏值按名拒绝
--- PASS: TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo (0.05s)    ← 能力位与用户开关两枚独立键
ok  github.com/CarlosShao/wisp/cmd/wisp 7.00s   （15 枚选中全 PASS／0 FAIL／0 SKIP；vehicle rc=0）
```

`internal/panel` 定向（裸跑，允许）：`go test ./internal/panel/ -run 'TestAC1FieldTableCarriesNoApprovalVocabulary|TestConfigRoute|Test452|Test313|TestConfig248' -count=1` ＝ `ok 0.185s` rc=0（含 `Test313SampleRateStaysAGuardKeyNotAWriteField`、`Test313AdditionalFieldNamesCarryNoApprovalWordAndNoLockedFamily`——本腿文案改动未破任何钉）。`go vet ./internal/panel/` rc=0；`gofmt -l` 空（本机 `gofumpt` 不在 PATH，以 `gofmt` 代，如实记）。

## 2. 前端三行改成什么（原文对照）

写面四枚文件：`frontend/src/lib/config-view.ts`、`frontend/src/components/voice-settings.tsx`、`frontend/src/fixtures/settings-v2-preview.tsx`、`frontend/src/components/harness/settings-page.tsx`（注释）。

- **回声消除（voice.aec.enabled）**：原文「后端已接：喂给打断监测的回声消除开关（重启后生效）」＋`StateTag "等快照"` → 改为显示**真当前值**（`当前：开/关`，读 `view.voice.aecEnabled`）＋**Switch 草稿态＋保存按钮**（写 `config.set` 字段 `voice_aec_enabled`，值 `true/false`）；`aecBuildAvailable === false` 时开关置灰并明说「这枚二进制没有编译回声消除（缺 webrtcap 构建标签），开关只能保存配置意图，当前不会生效」（票 452 AC#2 的两维不合并）。null（腿没答）仍逐行 `等快照`，指针闸保留。
- **输入设备（audio.input_device）**：原文「后端已接：每次开麦与热插拔重开都按这个名字解析设备」＋`等快照` → 显示**真当前值**（`default`/空串渲染为 `default（跟随系统默认）`）＋**文本输入＋保存**（写字段 `audio_input_device`；空值前端先拒：「输入设备名不能为空；跟随系统默认请填 default」，与门行为一致），解析失败时开麦侧按名拒绝的那句话保留。
- **采样率（audio.sample_rate）**：原文「后端守卫：采集链只跑 16000，配置里写别的值会被拒」＋`等快照` → 显示**真当前值**（`当前：16000 Hz`）＋`守卫·只读` 标签，**不提供控件**——写入面刻意没有这一字段（守卫键），页面不得画一枚改不动的滑杆（票 83 那味）。
- `ConfigViewVoice` 类型补四键（`aecEnabled`/`aecBuildAvailable`/`inputDevice`/`sampleRate`，null＝未答）；两份夹具（`settings-v2-preview.tsx`）补同键：未配置夹具＝四枚 null（等待态）、已配置夹具＝`aecEnabled:false, aecBuildAvailable:false, inputDevice:"default", sampleRate:16000`（shipped 构建的事实形状）；页头导语与三处过期注释同步改真。

**门禁读数**（node 在 `D:\work\server\node24`，v24.21.0／npm 11.19.0）：

```
npm --prefix frontend run typecheck  rc=0
npm --prefix frontend run build      rc=0   （dist/index.html 1.04 kB + assets/index-Cpu0tTrq.css 65.27 kB + assets/index-DZQOdBdK.js 643.77 kB）
sh scripts/d22scan.sh                rc=0 clean（ban #8 frontend/=123 含本腿 TSX，零 emoji 违规）
```

## 3. AC#10 五枚 [panel] 字段生产者读数（载具同发，8 枚全 PASS）

`TestAC10*` 8 枚 `--- PASS`（0.00s-0.01s 各发）。五枚字段归位（schema `internal/config/schema.go` `PanelSection`）：

| 字段 | 生产者 | 读数 |
|---|---|---|
| `width` | 票 255 `panelGeometrySource`→`windowOptions()` 交唯一建窗点 | 在产（本腿未动） |
| `height` | 同上 | 在产（本腿未动） |
| `enabled` | 33-r3 `withPanelEnabled(panelEnabledSource(dataDir))`，闸在 `bringUp` 顶部、runtime 探测之前 | `TestAC10DisabledPanelRefusesBeforeTheRuntimeProbe`＋`...EnabledSourceReadsTheRealConfigFile` PASS |
| `keep_alive_in_session` | 33-r3 `withPanelKeepAlive(...)`，`Hide` 末 `teardownAfterHideUnlessKept()` 真拆窗 | `TestAC10KeepAliveFalseTearsTheWindowDownOnHide`（3 子测）PASS |
| `scale` | **无，具名**：`forks/go-webview2` 的 `PutZoomFactor` 只有原始 `ComProc` vtbl 槽、无 Go 包装、controller 未导出——落地必须改 `forks/**`（本腿禁写） | `TestAC10ScaleHasNoProducerAndTheReasonIsNamed` PASS；摘掉障碍（fork 补包装）该钉会自动转红逼接线 |

AC#10 框**不翻**（勾归非实现者）。

## 4. AC#12 读数（构建后 Go 侧能力读数，跑完已删回 `.gitkeep`）

`TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` PASS，逐字：

```
panel_host_gate_test.go:118: AC#12 reading (head e49962d55): shape=page-bundle built=true
entry-bytes=1044 entry-ctype="text/html; charset=utf-8" entry-err=<nil> refs=2 check-err=<nil>
manifest-entries=4 manifest-err=<nil> | git-metadata=true tracked=1 tracked-beyond-anchor=0
tracked-has-entry=false ignored-or-untracked=2
```

即：`frontend.Dist()`（embed 后）4 条目（`.gitkeep`＋`index.html`＋2 枚带哈希 assets）、入口真存在且 `Check()` 2/2 引用可解析、shape 分得开 anchor-only 与 page-bundle。**跑完已删回**：`frontend/dist` 现只剩 `.gitkeep`，`git status --porcelain -- frontend/dist`＝空；**本腿未提交任何 dist 字节**。

## 5. AC#13 读数（冷启动两发次序＋真页面断言，两枚全 PASS）

```
TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe PASS (0.00s)
  33-r10 world: embed resolves 1044 entry byte(s), 1 element id(s) [root]
  read: 2 handover act(s), first=doc:"<!doctype html>…wispProbeRT…" last=nav:"https://wisp.panel/index.html"
  read: last act navigates to the entry URL and that URL serves 1 of 1 entry id(s), lookable=true,
        bytes-equal-to-resolved-entry=true
TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe PASS (2.24s)  ← 真窗
  probes from the resolved entry (1044 bytes): 1 id(s) [root]
  447-ac13-handover-settled on probe attempt 10: "https://wisp.panel/index.html|complete|Wisp"
  page answer: "1" - 1 of 1 probe id(s) present in the live document
```

最终文档＝入口页本体（`bytes-equal-to-resolved-entry=true`），不是探测页；真窗用例的凭据是**页面自己报**"1/1 个探针 id 在活文档里"，不是"SetHtml 被调用过"。次序修复（`coldStartPageHandover`，先探测后 `serveEntry`）为 33-r10 在册产码，本腿零改动。

## 6. AC#14 只读结论（⛔ 本轮未动消息环一字）

**"Go→页面由谁投递"今天的答案（现量行号，本树核）**：专用面板 STA 线程泵库的 `Run()`＝`cmd/wisp/panel_resident_windows.go:437`；`Dispatch` 生产调用点＝`:498`（投递闭包上线程）＋`panel_host_windows.go:713`（`w.Dispatch(func(){ m.deliverPanelPush(w) })`，票 329 快照推送）＋`panel_control_windows.go:807/814/818/822`（RPC 回执 resolve/reject，全部 `Dispatch`→泵→`Eval` 注 JS）；直接 `Eval`＝`panel_host_windows.go:772`。**即票面三形中的 ⓐ（换库 `Run()` 形）已于 2026-10-01 裁定 P1 并落地为现状**；`dispatchq` 不再"只涨不取"。

**本腿两枚真窗凭据（载具同发，页面自己报话，非"Dispatch 被调过/Eval 没报错"）**：

```
TestAC14AwaitedBindingReplyReachesThePage PASS (1.78s)
  nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED"（3/3 真请求的绑定回话真到页面）
TestAC14GoSideEvalPushReachesThePage PASS (2.19s)
  nail 2 (Eval push hop), page's own words: title="PUSHED-33R5-OK"（主动推送与回执两维分开，各自有钉）
```

**三种形取舍表（给下一波裁决；本腿零选型执行）**：

| 形 | 内容 | 现状 | 代价 | 取舍建议（腿评，非裁定） |
|---|---|---|---|---|
| ⓐ | 宿主消息环换成库 `Run()` 形（专用 STA 线程） | **已落地＝现状**（裁定 P1＋`panel_resident_windows.go:437`；线程名 `panel-sta` 不入 D38 名册＝测试在册核过 `TestPanelThreadIsSTAAndExitsCleanly`/`TestPanelThreadNameIsNotInResidentRoster`） | 已付：一条常驻线程＋退出链（pump 返回即线程收尾） | 维持；撤销口令「33 面板线程改回投 ui-sta」仍在册，撤销＝重开 33-a2 ⑦-A 那一族问题 |
| ⓑ | Go→页推送一律改走 `Eval` 直塞 | **部分在用**：主动推送可用 Eval（`:772`、测试钉 PASS），但**绑定回话救不了**（回执走 `Dispatch` 队列，R27 已量） | 若全改 Eval：绕过 correlationId 路由的队列纪律，C17 那半句要重写 | 不取；维持"回执走队列＋推送可 Eval"的分工，两枚钉已分开钉住两维 |
| ⓒ | 换依赖／fork 加可注入 drain | **不再需要**：`Run()` 在泵，`dispatchq` 有取出者；`WebView` 接口名册里无可导出 drain 的缺口已被 ⓐ 消掉 | 最大（fork `go.sum` 边、票 244 零新增依赖族尺） | 不取；除非 ⓐ 被撤销 |

**剩余缺口（如实）**：AC#14 框仍未勾——勾归非实现者；"出向主动推快照"那一跳票 329 已接（`deliverPanelPush`），本腿读数只补"两枚真窗钉在当前 HEAD 仍绿"这一维。

## 7. 硬约束自证

- ⛔ `thresholds.go`/golden/`allowlist.txt`/`docs/PLAN.md`/`docs/specs/**`/D43/`design/**`：零改动（`git diff --stat` 只有 §2 四枚前端文件＋`internal/panel/config_handlers.go`＋本件）。
- ⛔ 放宽断言：零（Go 侧仅新增一条拒绝文案分支，未动任何既有断言；前端零测试删改）。
- ⛔ push：本腿只 commit。
- 新协程：零枚新增。
- ⛔ 反向动 C17 白名单：`internal/panel/bridge.go` 零字节改动。
- ⛔ 提交 dist 字节：`git ls-files frontend/dist`＝1（`.gitkeep`），构建产物已删回。
- 环境注意：机器上有 Wisp 球 PID 49348 在跑；`TestAC4Focus…`/`TestAC246…` 不在本腿 -test.run 名册内，未受其染。
