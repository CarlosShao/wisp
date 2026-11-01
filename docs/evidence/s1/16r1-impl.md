# 票 16 证据组装腿 `16-r1` —— S2 切片验收门·格①②③④取证表（写腿）

> **身份**：写腿（实现者侧证据组装）。**⛔ 翻框**：五格勾选归 16-v1（对抗验收）＋编排者；
> 格⑤（对抗验收）本腿自裁、未做。**尺**：开工 `HEAD=19d715c5a3cc`（`dev`，编排者 A1342
> 收 12-v3 解锁票 16 的那一笔），索引 0 staged；写面无他人脏（在飞脏件＝`382-w-r2` 的
> `cmd/wisp/resident_voice_{engines_,}windows.go`＋`393-r2` 的 `internal/audio/aec/*`＋
> `.scratch/wisp/issues/30-*.md`——**零触碰、零归因、不进读数**）。读数现跑
> 2026-11-01 09:0x–09:5x +08，全部实验在仓外导出树
> `D:\tmp\wisp\20261101\16r1\export`（`git archive HEAD` 解树＋junction
> `build`/`third_party` 指主仓资产＋三 DLL 贴身＋`WISP_*` 全套 env＋cwd=导出根
> ——41-v1/390-v1 同款 staging）＋**主仓共享树零突变**。产品码零字节改动（见 §7 申报）。
> **音频错峰**：每轮真设备跑前 `tasklist | grep -iE "wisp|balldebug"` 现核（demo 前
> `DESKTOP_CLEAN` 在案）；全程单进程串行。**跑法纪律**：测试二进制一律 Git Bash 直调
> （390-v1 §9 踩坑同款规避）、退出码直接取不经管道、`unset GOTMPDIR`。

## 0. 一句话

五格中的 ①②③④ 取证完成：**①** staged 真硬件一轮完整语音 E2E 全链绿
（真 KWS 唤醒→真流式 ASR→真 agent 回合（tool_call 落真库）→真 matcha 播报（真
waveOut）→SCOPE-WARM→3s 静置→SCOPE-DISPOSED，D43 五行全产线派发）；
**②** 延迟对照表进 `docs/SLO.md` 附录 D（票面明令位，只加不改）；
**③** 六场景 4 全执行＋2 具名边界（见 §4）；**④** CER 门以 CI 同款 runner 本地绿
＋接线在案＋「真 CI 工件待 push」具名残余（网络/推送物理受限，判据见 §5）；
零回归六族 28/28＋262/640/902 逐位同。**格⑤＝16-v1，待派。**

## 1. 停手条件自核（派单三条「未定义即停」）

1. **demo 全链跑不通？**——未触发：全链一轮绿（§2）。
2. **六场景清单与 SPEC-10 §7 对不上？**——未触发，映射成立（§3 表头）：票面六场景
   ＝§7 基础八项中 S2 相关三项（麦克风被占/磁盘满/模型下载中断）＋D42 十二项中 S2 行
   （#12 麦克风被独占/权限被拒——SPEC-09 §4 表明写「S2」归口）＋#2 音频热插拔，
   外加 C29 篡改面（SPEC-10 §6 #11）拆出的「安装后损坏」半。**逐项可追溯，无缺号。**
3. **SLO 附录格式无先例？**——未触发：附录 A/B/C 已立「叠加不改写」体例，附录 D 同形。

## 2. 格①——E2E voice demo（log＋state trace＋tool rows、真硬件）

**测试**：`cmd/wisp/voice_e2e_demo_16r1_windows_test.go`
`Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware`——**PASS 28.37s**（log
`D:\tmp\wisp\20261101\16r1\logs\16r1-demo2.log`；首轮 `16r1-demo.log` 28.91s 除一条本腿
错用 12-r2 文本路判据的断言外全同，照录不删）。真硬件事实：本机 Realtek 采集＋渲染端点
在案，`TestLiveWasapiSmoke`（`WISP_LIVE_MIC=1`）当日两跑 PASS（导出树 HEAD 构建：
25 帧/820ms，`16r1-livemic-head.log`）——采集与播放两端都是真设备；音频数据缝按 C8 只
注入 wav（AGENTS §1.3 纪律），sherpa 三引擎与 matcha 全为真模型真推理。

**全链读数（逐行引 log）**：

| 链节 | 读数（逐字） |
|---|---|
| 负控 | 0.wav（checkpoint 自带负样本）入 Armed：`kws_hits` 增量 0、机器守 Armed |
| 唤醒 | 4.wav：`kws_hits=1`；`DEMO16 D43 Sleeping+kws.enabled->Armed at_sample=0`（#5）→ `Armed+kws.wake-word->Listening at_sample=115200`（#7）→ `Listening+asr.vad-stop->Thinking at_sample=163840`（#11） |
| 交还 | `handoffs=1 handoff_exit=0`（真 manifest 验签走位）；`asr_loaded=true` |
| 转写 | `TRANSCRIPT ["嗯大家如果对我的研究感兴趣呢嗯"]`（真流式 Paraformer 对 4.wav 的转写）；`partials=11 finals=1` |
| agent | golden SSE（C5 缝）两轮：`fs.read`（allowlist 内笔记）＋纯文本收尾；`DEMO16 TOOL_ROW tool=fs.read level=L0 decision=allow outcome=success corr=…#call_16r1_demo`；`TASK_LOG id=… state=done query="嗯大家如果对我的研究感兴趣呢嗯"` |
| 播报 | 产线 `speak()` 全 hops：`announcement spoken runes=21 chunks=3 samples=88856 rate=22050 acquire=3.7150564s`（真 matcha 合成→真 waveOut 播放）；`SPEAK asked=1 spoken=1 aborted=0 gated=0` |
| 通知 | `NOTIFY 1 posted ("Wisp\0已经读完了，笔记里只有一行话，随时可以继续")`（D10） |
| 回落 | `D43 Thinking+llm.first-token->Speaking`（#15 纯文本臂）→ `Speaking+tts.speak-done->Warm`（#26）→ `SCOPE-WARM warm_timeout=1.2s` → `SCOPE-SETTLING static=3s` → `SCOPE-DISPOSED reason=settling.expired cascade=…` |
| 球走位 | `BALL_STATES [Armed Listening Thinking Speaking Warm]`——**五行全部产线派发，本腿零手工桥** |
| settle | `SETTLE token_to_settling=8984ms（含真 TTS 装载 3.7s＋播报 4.1s）settling_to_fallback=3039ms token_to_fallback=12024ms` |
| 延迟 | wake 路 utterance `openDelay=6.95s chunk=640ms finalLatency=3.52s speechMS=2816 closedBy=end seamWaitToFirstText=9.99s`（唤醒后才开句的形状；262/640/902 基线是 Test15e 单句口径，同日六族复跑逐位同值钉住，两者⛔互当） |

**诚实边界（照录）**：
- **Acting 不在球走位里**：今日 HEAD 的常驻语音机在 agent 工具执行期不接任何 D43 派发
  （announce 路径以 `Facts{HasToolCall:false, HasAnnouncement:true}` 走 #15 纯文本臂，
  `resident_voice_engines_windows.go:979-980` 现核）；工具证据落 `tool_call` 行＋审计，
  不落球态。票面链字的「Listening→Thinking→Acting→notify reply」中 Acting 一环在
  HEAD 语音路无可观测载体——照录，归 16-v1/编排者裁量是否算链缺口。
- **settle 的 token→Settling 腿**（8984ms）含真 TTS 装载＋播报本体，与 12-r2 桩 TTS
  文本路形状（~1.2s）不可比；3s 静置窗本体 3039ms 为冻结真值，两形状同钉。
- **音频错峰**下 demo 独占真扬声器约 4s（播放期采集门关闭——`gate_closed` 语义由 26r3
  套件守，本腿未重断言）。

## 3. 格③——六场景 pre-mortem（SPEC-10 §7；检测→行为→文案→D37 类）

| # | 场景 | 执行方式与读数 | D37 类 | 边界 |
|---|---|---|---|---|
| 1 | 麦克风被占（test process holds device） | seam 注入真实 HRESULT `AUDCLNT_E_DEVICE_INUSE` 过**真 open 路径**（产码映射）→ Start 响亮失败＋`stats.LastError` 记录＋设备名＋`inUseGuidance`：`TestOpenOccupiedAndPermissionDenied/occupied` PASS（`16r1-audio-premortem.log`）；正控＝同路径真设备可开（live smoke PASS） | `audio_device`（D37 表「麦克风被独占」逐字典型） | 独占持麦需外部独占流持有者（产码只开共享模式、无 COM 探针面可脚本化）；共享模式共存不构成占用 |
| 2 | 麦克风权限被拒 | **seam 半**：`hrEAccessDenied` 过真 open 路径→`privacyGuidance`＋设备名：`/permission_denied` PASS。**OS 半（真触发尝试）**： ConsentStore 两枚杠杆（主键 `microcontroler`、`NonPackaged`）翻转 Deny 后真探针（`WISP_16R1_MIC_DENY_PROBE=1` 门控，TestLiveWasapiSmoke 同款门形）**设备照常打开**——负结果照录（`16r1-mic-deny.log`/`-deny2.log`，恢复后 `Value=Allow`/键删除核过）；本机 consent store 无 per-app 条目、两杠杆均未触及 WASAPI 强制面 | `audio_device`＋引导系统设置（SPEC-09 §4 #12 逐字） | 真 OS 级触发剩余路径＝设置页人工翻或 HKLM 组策略（admin）；按「未定义即停」精神**不冒充**：seam 层映射为真产码执行，OS 层触发不可达已具名 |
| 3 | 设备热拔 mid-utterance | **seam 三连**（枚举缝注入＋真采集循环、真重采样、真有界推送——`hotplug_test.go` 头注的既立纪律）：`TestHotplugStaleHandleFailsLoudly`（流在采集中途失效＋重开被占→终态错误**点名设备名**＋in-use 指引＋帧停流不静默）、`TestHotplugReenumerateOnce`（D42#2 事件→恰一次重枚举＋新端点数据续流）、`TestHotplugReenumerateFailureNamesDevice`（重枚举失败→点名前设备）——三发 PASS。**状态机半**：`TestEveryLegalRowFires/#14_audio_device_lost` PASS（`Listening+EvAudioDeviceLost→Error`，副作用 `error.device-name`） | `audio_device`（D37「热拔出」典型） | OS 级真禁用设备需 admin（`pnputil`/`Disable-PnpDevice`）；本机 Realtek 内 construção无物理可拔件——admin=no 现核，边界具名 |
| 4 | 转写落盘时磁盘满 | **现读执行**：真写失败面＝`run.go:1461-1473`——`StartTaskLog`/`FinishTaskLog` 失败逐条 audit（`wisp run: task_log open/close failed: %v`）**不静默**、响应路径不受影响；存储层 WAL 防线＝启动 `wal_checkpoint(TRUNCATE)`＋`wal_autocheckpoint(1000)`（`open.go:295-319` 注释自名 "disk full defense"），`go test ./internal/memory` 当日整包 ok 22.742s；健康路径正控＝demo 的 `task_log state=done` 行 | `resource`（D37「磁盘满」→`WatchdogAlert`） | **200MB 写前预检＋用户警示＝D42#5 归票 45（`ready-for-agent`，0 AC）未建**——预演如实记「今日无该守卫」；真磁盘满夹具需小虚拟盘/配额（admin，票 45 AC#2 自名），本机不可达，未伪造 |
| 5 | 模型安装后损坏（sha256 fail on load） | **全执行**（`Test16r1PremortemModelSha256CorruptRefusesTheHandOff` PASS）：改装后 store 的 encoder 一个字节（替换法，⛔ 不让硬链接把翻转带回 third_party 钉件）→ 检测半 `VerifyInstalled` 逐字 `sha256 mismatch (pins 81a70226a893..., got ab9962e4c1a2...)` 类=`model`；走位半＝生产 hand-off **拒交未验证字节**：自愈重取（12.884s vs 基线 ~2s）后 store 恢复可验证——**用户永不拿到被篡改字节**；不可自愈分支（镜像篡改→`sha256 MISMATCH … staged file deleted` ClassModel）由 models 包 `TestAC21FreshInstallIsGuardedOnEnsuresOwnReturnPath` 钉住（`sha256 mismatch` 断言在案） | `model`（D37「哈希/签名不符」逐字典型） | 无（真触发：真文件真翻转真校验） |
| 6 | 下载中杀进程＋恢复 | **全执行（真镜像真进程杀）**：`wisp models ensure asr-streaming-paraformer-zh-en`（237MB 工件，真 URL）跑至 staging `encoder=165,462,184 B＋decoder=2,673,470/71,664,561 B` 中断态时 `Stop-Process -Force`（PID 18808）→ **partial 存活**（staging 目录原样）→ 二跑 `已交还（state=FirstRun，ticks=4454，末次进度 100%，9.11s）` RC=0 → `wisp models verify` 逐字 `与已验签清单一致（store=…，1.81s）` RC=0；**Range 续传协议**由 `TestResumeContinuesAtByteOffset` 服务端 Range 断言钉住（当日 PASS；真镜像侧无法观测服务端 Range 头——kill 容错与协议证据的分工照录） | `model`/`network`（中断传输 `transfer interrupted (partial kept for resume)` ClassNetwork；重验失败 ClassModel——downloader.go:542/:557 现核） | 无 |

**六场景清点**：4 项全执行（3/5/6 真触发＋1 的真 HRESULT 真路径＋2 的 seam 半）；
2 项具名不可达边界（2 的 OS 级、3 的 OS 级、4 的真满盘）——**全部有负结果读数或
代码级现核在案，无一项静默跳过**。

## 4. 格②——延迟表（落到 `docs/SLO.md` 附录 D）

`docs/SLO.md` 追加**附录 D**（票面明令编辑位）：D32 16.3.4 预算行 vs 实测 12 行
（唤醒 0 样本差／VAD 判停→全文 704ms PASS／标点开销 0.53ms PASS／组装与 LLM 首 token
两行未测如实记／TTS 首帧需加载分支 3.715s **FAIL 照录**／冷会话首用 5.51s 超 P50 含首装
照录／Warm P50 deferred（票面自带分支）／打断 41–46ms 引 381-r2/390／3s 回落
3000–3039ms PASS／settle 261–263ms PASS／Sleeping 复跑 skeleton 口径引用 C.5）。
**正文与 A/B/C 零改写**；预算行一字未动。

## 5. 格④——CER 门（CI 工件判定）

**现核事实链**：①CI 接线在案——`.github/workflows/ci.yml:732` `run: bash
tools/d22scan/runtests.sh ./internal/speech/cer/`（Ubuntu 专步）＋
`scripts/portable-tests.sh:236/:295/:353` cer 入核心清单（2026-10-26 票 15 AC#2 接线）；
②仓内存档 CI 日志（run-37158259050=10-03、37240161874=10-04、37249563077）**全部早于
10-26 接线，无一含 CER 步**——真 CI 工件今日不存在；③本腿⛔ push（推送只归编排者，
A813），网络虽通也无法触发 CI。

**判定（拿得准，不停手）**：以 **CI 同款 runner 的本地执行＋接线引用＋具名残余** 作
格④的记录，理由＝「真 CI 工件」的唯一生成路径（push）不在本腿权限内且无任何本地动作
可替代；门的本体（24/24 绿）机器可核：
- `bash tools/d22scan/runtests.sh ./internal/speech/cer/` ⇒ **`OK - packages=
  [./internal/speech/cer/] top-level: PASS=24 FAIL=0 SKIP=0`**（`16r1-cer-ci.log`，
  rc=0——与 ci.yml:732 逐字同命令）；
- `go test ./internal/speech/cer -count=1 -v` ⇒ **24 PASS/0 FAIL/0 SKIP**（`16r1-cer.log`）；
- `near_clean` mandatory：门集含 60+40 提交集的逐 wav 哈希与 perfect=0/empty=100% 红
  断言（15-v1 §6 判语照录）；`noisy_far` status：阈值口径 EXERCISED-not-MET（15-v1 R1
  具名残余——无 `StreamingAsr`→`cer.Set.Score` runner，6%/15% 无真 ASR 跑分）照录。
- **具名残余**：真 CI 工件＝下一次编排者 push 后该步的 run 页留存（届时此格可凭真工件
  翻勾）；本腿记录不冒充 CI 工件。

## 6. 零回归＋四数＋两尺＋d22scan

| 尺 | 读数 |
|---|---|
| **六族合跑**（staged，`Test43|Test15e|Test26R3|Test41a2r1|Test381r2`） | **rc=0，28/28 PASS，0 FAIL/0 SKIP**（`16r1-sixfamilies.log`）；**`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同**（`seamWaitToFirstText=4.966s finalLatency=4.096s speechMS=3392` 与 12-v3 逐字同） |
| d22scan | 导出树（HEAD＋本腿测试件）**rc=0 clean**；分母 `bans#1-5 internal/=287 cmd/=52、ban#8 internal/=730 cmd/=183`；**正控**：种真 U+2713 字符串 ⇒ rc=1 逐字点名 `zz_16r1_poscontrol.go:3 [emoji] …string literals included`，删后复跑 clean（可逆） |
| 甲 gofmt | 本腿 `.go` **0 行**（rc=0） |
| 乙 gofumpt（`D:\work\base\gopath\bin\gofumpt.exe`） | 本腿 `.go` **0 行**（rc=0） |
| go vet | `./cmd/wisp/` rc=0（含本腿测试件） |
| 邻接包 | `internal/memory` 整包 ok 22.742s；`internal/audio` 定向五发 rc=0；`internal/statemachine` TestEveryLegalRowFires rc=0；`internal/speech/cer` 24/24 |
| 协议差照录 | `internal/speech` 整包在导出树构建的二进制下 `TestDefaultKwsConfigResolution` 红（`kws_test.go:284` walk-up 锚＝编译期源码路径，导出树构建解析出 junction 字符串）——12-v3 §5 同款 cwd/构建位协议差；本腿对 `internal/speech` **零字节**，六族（本腿义务面）全绿 |

## 7. 写面申报（逐枚）＋格⑤

**主仓提交 pathspec（全部写面）**：
1. `A cmd/wisp/voice_e2e_demo_16r1_windows_test.go` —— 格① demo＋格③ 场景 2/5 的
   测试取证面（产码零字节；现有测试缝的全部消费，零新缝）；
2. `A docs/evidence/s1/16r1-impl.md` —— 本件；
3. `M docs/SLO.md` —— 仅追加附录 D（票 16 AC#2 明令位；`git diff` 只含追加行）。

**临时件（只建不删）**：`D:\tmp\wisp\20261101\16r1\{bin,dl,export,logs,tmp15e,d22scan.exe}**`
（logs＝demo/demo2/sixfamilies/cer/cer-ci/audio-premortem/d43-rows-v/livemic-head/
mic-deny/mic-deny2/resume-range/memory/speech-staged* 等）。

**⛔ 未碰**：`internal/{risk,session,speech,audio,statemachine}/**`、D43 表、
thresholds/golden、testdata 仓根、frontend/design、scripts/.github、票面、台账、specs、
PLAN、`docs/SLO.md` 正文与附录 A/B/C。**⛔ push。**

**格⑤（对抗验收）＝16-v1，本腿自裁未做**——票面第五格「Adversarial acceptance pass
by a different agent (stub scan + contract conformance)」待编排者派非实现者腿；
本件与其引用的全部 log 即其验收对象。

**残余具名**：①真 CI 工件待 push（§5）；②OS 级麦克风权限/热拔/真满盘三触发不可达
（§3 边界列）；③Acting 球态载体（§2 边界）；④TTS 首帧 ≤1600ms 需加载分支 FAIL（附录 D，
S0 §4 同源方向，Warm 形状归票 28）；⑤`internal/speech` 整包在导出树构建位的
kws 路径字符串红（协议差，§6）。
