# 票 41 终裁腿 `41-v1`（非实现者裁腿裁决表）— KWS 接线 41a2-r1（`3d41c1791`）的六格终裁

> **尺**：开工 `HEAD=cf35a3c39`（`dev`）。★**复跑读数含 381-r2（第 21 处定式）**：381-r2＝`77c2cabfb`/`c4f4799f4` 已在 dev、其打断/Abort/token 案全部随跑。全部实验走仓外导出树 `D:\tmp\wisp\20261031\41v1\export`（`git clone --shared`，收工 HEAD=`8e2a411bd`＝比开工基线多一笔台账 `pending-and-issues.md` +17 行、代码面 0 差；`3d41c1791` 与 `77c2cabfb` 均在祖先链；`git status` 净），`third_party`/`build` 以 junction 补 staging。工作树在飞脏面＝`cmd/wisp/resident_windows.go`＋`run.go`（43-r1，⛔碰⛔归因）。现跑 2026-10-31 13:02–13:3x +08。★读数逐条现跑，⛔ 转述。
> **写面**：本文件＋仓外 `D:\tmp\wisp\20261031\41v1\**`（bin/mut/overlay/logs/export）。零 `.go`、零翻框、⛔ push。
> **staging 形状（复跑用）**：15e 契约原样——`go test -c -o bin/wisp41v1.test.exe`＋三 DLL 拷入 bin（**DLL 贴身**：`System32\onnxruntime.dll`=1.17.260311 是 Windows 自带、先于 PATH，PATH 注入被它遮蔽）＋全套 `WISP_KWS_DIR/WISP_SILERO_MODEL/WISP_SPIKE_MODELS/WISP_PUNC_MODEL/WISP_TTS_DIR`＋CWD=export 根（manifest 走 junction）。

## 0. 一句话

**★六格今天一格都翻不了整格（每格有缺角），但按 15/28/42 先例★可带残余结案**：实现面（KwsEngine→run() 帧管道→D43#5/#7/#11→warmFeed 真门→Facts.KwsLoaded）全部复验为真，41a2-r1 六格自评与我复判**逐格相符**（诚实、无纸糊）；七条残余全部具名、各有归属、无一被本树代码阻塞。

## 1. ★★★ 唤醒端到端复跑（staged 二进制，五案 5/5 rc=0）

| 案 | rc | 读数（逐字要点） |
|---|---|---|
| WakeEndToEnd（16.78s） | **0** | 0.wav 负控：kws_fed 前进、**0 hit**、机器守 Armed、handoffs=0；4.wav：kws_hits=1、`Sleeping+kws.enabled→Armed at_sample=0`（#5）→`Armed+kws.wake-word→Listening at_sample=115200`（#7）→`Listening+asr.vad-stop→Thinking at_sample=163840`（#11，语音 2816ms）、asr_loaded、handoffs=1、handoff_exit=0、**partials=11 finals=1**、ball 收 Armed→Listening |
| WarmGateFeedsKwsNotTheAsr（15.36s） | **0** | Warm 态 40 帧静音 ⇒ **kws_fed 恰 +20480、asr_fed +0（冻结）**；静音 0 额外命中；Warm 唤醒 ⇒ hits=2、warm_wake=1、门重开（warmFeed=false）、推 8 帧 asr_fed 前进、**loads_after_warm=0**、acquire_log=[3.15s 0s]（第二位 0s＝Warm 免装载，B4） |
| GuardRowsWalkWithTheProductionFact | **0** | #8/#10/#32 双分支全走：生产 fact getter、真实手势入口 `onMuteGesture`→owner 队列（D38b）；kws-not-loaded 分支同样走通 |
| ProductionTableRejectsTheDemoWavs（3.89s） | **0** | 生产拼法引擎×demo 7 夹具 **0 hit**（管线级控制案，防呆在案） |
| OptInGateResolvesTheManifestSpelling | **0** | off ⇒ Kws 空＋note「enabled=false」；on ⇒ manifest 五件逐名解出（epoch-12 encoder 尾名钉） |

三跳 AtSample=0/115200/163840 与 `41a2r1-impl.md` §2 **逐位同**；acquire_log 首位 3.15s（实现件 2.94s）＝装载钟时长、如实报。

## 2. ★★★ mutation 两发复刻（仓外 `-overlay`，树零突变）＋还原核

| 发 | 变异（按 HEAD 形状） | 面 | rc | 红句（逐字） |
|---|---|---|---|---|
| M-a | `kwsFeed` 首行插 `if true { return false }` | WakeEndToEnd | **1** | `the engine was not fed in Armed (kws_fed=0 before, 0 after): the Armed gate is closed` |
| M-b | `resident_voice_windows.go:933` 门改 `if false && rv.warmFeed.Load() && rv.kwsFeed(samples)` | WarmGate | **1** | `kws_fed moved 0, want exactly 20480: the KWS is not the Warm ears` |

**还原核四件**：①净二进制重跑两案 **rc=0**（21.17s/14.28s）；②sha256 台账：树件 `DD2F88D6…E8FB86D`＝开工值＝**工作树同件**（mut-a `206DFA2C…`、mut-b `54510722…`）；③`git diff --no-index`：mut-a 恰 1 insertion、mut-b 恰 1 ins/1 del（只差变异块）；④export 树 `git status` 净。⇒ 工作树从未进突变窗口。

## 3. ★★ 零回归合跑（`-run 'Test15e|TestWireR1|Test26R3|Test41a2r1|Test381r2'`）

**26 PASS／0 FAIL／0 SKIP，rc=0**（15e 6＋WireR1 7＋26R3 3＋本腿 5＋**381-r2 5**——`RealDevicePlayingPhaseAbortWithin400ms` 真跑 35.52s 非 skip）。★归因：合跑在 HEAD（含 381-r2＋102r2 修复）全绿 ⇒ **两腿无交叉红，本腿接线与 381-r2 的打断/Abort/token 变更共存无犯**。
- 四数：`openDelay=262ms chunk=640ms ballLatency=902ms` ⇒ **262/640/902 逐位同**。
- 15e 计数：`drained=181/181 partials=5 finals=1 handed=1 handoffs=1 asrLoaded=true problems=[]` ⇒ **本跑零抖动**（首跑抖动案 WarmIdleDisposesOnTheOwner 本跑绿；抖动只报不判）。
- 26R3：3 PASS（announce 关捕获门语义未动）。

## 4. ★ d22scan 导出树＋正控＋两尺＋vet

- `sh scripts/d22scan.sh`（export 树）：**rc=0 clean**，`bans#1-5 cmd/=45`、`ban#8 cmd/=162`（实现件记 161＝其新测试件；本跑 +1＝381-r2 的 interrupt 测试件，账对得上）。ban#8 射程＝design/frontend/internal/cmd，**未扫 docs/**。
- **正控有牙**：植入裸 `go func(` ⇒ rc=1，红句逐字 `[bare-goroutine] bare `go func(` is banned (D22/D38b)`；删后复跑 rc=0 clean。
- gofmt/gofumpt 对本腿三件（HEAD 形状）单点**空**；`go vet ./cmd/wisp` rc=0 空。

## 5. ★★★ 六格逐条判（能翻/⛔＋缺什么）

| # | 判 | 翻格前置（缺什么） |
|---|---|---|
| AC1 | **⛔ 半格成立** | sampler RSS/CPU 实测＋唤醒 ≤300ms 实测＝票 42 watchdog/sampler 面；样本钟原料已备（三跳 AtSample） |
| AC2 | **⛔ 半格成立**（`hit on fixture audio` 已验） | `reload without restart`：keywords 键 config 层有热更分类管线（manager.go:395/:410、tiers.go:67）但 **KWS 路径零消费**（无词表→keywords 文件生成器）；thresholds/veto_words 扫描零命中（config_readers_255.go:237/:238）＝config 票/B 腿 |
| AC3 | **⛔ 结构就位、e2e 无从做** | `kwsVetoWords` 空表（internal/speech/kws.go:83）⇒ 本树不可能真 veto 命中（B 腿改表即激活）；「KWS absent→UI 明示」＝Confirming 通道票 |
| AC4 | **⛔ 不覆盖** | P14 阈值表＋推荐值未落（阈值仍包内常量 Score=3.0/Threshold=0.1）；本腿负控读数是原料 |
| AC5 | **⛔ 半格成立**（机制＋guard 行双分支已验） | 真引擎管线级「Muted 里 kws_fed 冻结」断言缺（GuardRows 用 stub 引擎；门语义在 kwsFeed :1028 静态可证：Muted 非 Armed 且 warmFeed=false ⇒ 拒喂） |
| AC6 | **⛔ 半格成立**（plan 级 opt-in 门已验） | 设备句柄探针（采集面）＋Sleeping 预算实测（42 面） |

## 6. 残余 7 条逐条定性归口（对照 41a2r1-impl.md §9）

1. **AC1 sampler 半格** → **票 42 watchdog/sampler 面**；不阻结案。
2. **AC2 词表热更** → **config 票**（255 审计既有债）＋**B 腿**（词表→keywords 文件生成器）；不阻。
3. **AC3 veto e2e** → **B 腿**（改表即激活）；UI 半句 → **Confirming 通道票**；不阻。
4. **唤醒句话轮过滤**（唤醒词转写会交给文本路）→ 票面/D43 未定义 ⇒ **先定义再接 `consume()` 一处**（B 腿/P14 或新小票）；实现腿未自创过滤＝**正确形状**；不阻。
5. **#32 内部定时器 nil facts** → `internal/statemachine` 冻结面开 facts 钩＝**改冻结面、须人工批准**（编排者摆）；现行 **#5 重武装自愈**已 e2e 验证兜住；不阻。
6. **KWS 字节不走 hand-off walk／models ensure 无 KWS 入口** → **models 票**（D26 面，36MB opt-in 同 VAD 先例）；不阻。
7. **Speaking 双门**（发声期 KWS 靠 warmFeed/状态双门停喂、26r3 捕获门是第一道）→ **观察项**；若要第二道显式断言归 26 家族后续；不阻。

## 7. ★ 可否带残余结案

**★可**（先例在案：15 收 3 格未勾、28 DONE-带残余 1 格、42 五格全翻带残余——带未勾格结案有谱）。判语：**票 41 的 A2 接线收；六格按各自缺角保持未勾，由 B 腿/42/config/models 后续翻格**；R28-2 的注销（warmFeed 升真喂食门已落地并复验）随本判归编排者记账。41a2-r1 自评与复判逐格相符：AC1/AC2/AC5/AC6 半格、AC3 结构就位、AC4 不覆盖——无冒称、无纸糊。

## 8. staging 注记（复跑三红归因，零代码红）

首跑三红**全是 staging 缺口、与 41a2-r1/381-r2 均无关**：①`WISP_SPIKE_MODELS` 未设（voice15eInstallAsrFiles 硬依赖 env、不走路旁回退）×2；②`System32` 自带 onnxruntime 1.17 遮蔽 PATH ⇒ 正解＝DLL 贴测试二进制（PATH 注入无效）×1；③VAD 缺失令 buildResidentVoice 早退空 plan（KwsNote=""）×1。复跑读数全部以 staged 二进制取。
