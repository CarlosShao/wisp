# 票 393 终裁腿 `393-v1` —— AEC3 绑定层 flush/reset API（非实现者裁决表）

日期 2026-11-01 · 裁决对象 `35084d0b0`（6 枚 +718/−0，真父 `a7d9f7f06` 逐枚 diff 坐实
＝申报写面恰好 6 枚、越界面 `aec.go`/`bargein.go`/`cmd/wisp`/`deps.toml`/PLAN/thresholds
**0 字节**）· 材料现核：票面 5 格＋`393r2-flush-reset.md`＋`392r1-session-integration.md`
＋钉定树（`D:\tmp\wisp\20261031\aec3r1\webrtc-audio-processing-2.1\webrtc\`）逐面现读
· 本腿全部实验在导出树 `D:\tmp\wisp\20261101\393v1\`（`git archive 35084d0b0`＋补 3 枚
sherpa DLL，A802 尺过）——工作树零触碰、⛔ push。

## 0. 一句话

**五格全过、可收口**；停手裁定（reset 半边按公开面落地、flush 半边具名不伪造）
经钉定树逐行现核成立；AC#2/AC#3/AC#4 本腿独立复跑**逐字节/逐位复现**（工件与
393-r2 在盘件 diff rc=0）；主新发现的两引擎形**可观测差坐实**、但其「生产吃全带 vs
tuned 吃分带」的**机制归属未证**（本腿读出两条实现未计入的懒同步路径）——处置建议
见 §6，保温修复票收窄为 resampler/detector 半成立。

## 1. 半停手裁定核（钉定树逐行现读，★最关键格）

1. `EchoControl`（`api/audio/echo_control.h:23-62`）：方法集 AnalyzeRender/
   AnalyzeCapture/ProcessCapture×2/GetMetrics/SetAudioBufferDelay/
   SetCaptureOutputUsage/ActiveProcessing——**无 reset、无 flush**。✓
2. `AudioProcessing::Initialize()`（`api/audio/audio_processing.h:497-507`）逐字：
   「Initializes internal states, while retaining all user settings. This should be
   called before beginning to process a new audio stream.」＝**真 reset 公开面**。✓
3. `InitializeLocked`（`audio_processing_impl.cc:516-588`）：首行
   `UpdateActiveSubmoduleStates()`、全部 AudioBuffer 重分配、`AllocateRenderQueue()`
   （:572）、逐子模块重建（:574-583）；`InitializeEchoController`（:1880-1916）经
   `echo_control_factory_->Create`（:1887-1889）重造——**tuned 会话重造同参数引擎**
   （factory 自持 config_，shim `TunedEchoControlFactory` 形）。✓
4. flush 半边无公开面：`kBufferFlush`＋`delay_controller_->Reset(true)` 全在
   `block_processor.cc:136-144`（render overrun 内部路径；:122-124 亦内部）；
   `kPlayoutAudioDeviceChange` 只进 render 队列（:821-822）、到不了 echo controller；
   `kPlayoutVolumeChange` 进双队列（:829-837）→ `echo_path_gain_change`
   （:1308-1311）→ `EchoPathVariability`（block_processor.cc:132-134）＝**加速重
   适配标志、非清态**。✓ 具名成立，伪造 Flush 判定正确。
5. **判语**：「reset 半边按铁律落地（APM Initialize＝公开面、语义正对保温窗间）、
   flush 半边具名不存在不伪造」**裁定对**。精度注记：reset 面在 `AudioProcessing`
   而非 `EchoControl`——shim 持 APM 句柄，落点正确；停手条款①对 flush 成立、对
   reset 不成立，与实现判语一致。

## 2. AC#1 形状核（对照 `wap_new_tuned`）

- `wap_reset`（wap_shim.h:37-44/cpp:138-145）＝`apm->Initialize()`，错误约定
  （note_error＋返回 rc）与 process 函数同形；Go `Reset()` 与 `ProcessCapture`
  同错误包装形；stub `Reset()` 镜像方法集；doc 具名 flush 半边不存在并指向证据件。✓
- `wap_new_tuned`/`aec.go` Tuning 零触碰（0 字节 diff）；**AEC3 源树零触碰**：钉定树
  全树 2026-10-31 12:00 后改动文件＝**0**、`libwebrtc_audio_processing.a` 未重建。✓

## 3. AC#2 复跑（导出树、`WISP_AEC393R2=1` 门、两发）

本腿自建（CGO 旗标＝build.py COMMON 全集＋`-include cstdint`＋c++17，build rc=0），
两发 rc=0：**run1/run2 读数逐位同（唯一差＝两行 PASS 时长）**；两枚工件 JSON 与
393-r2 在盘件 **diff rc=0（逐字节同）**。
- **reset==rebuild-control 逐位：4/4 发、帧差 0**＝carry 状态被完全清除。✓
- reset 带 −48.30..−46.94，**全部 ≤−39**（比 −24 近端回声低 15dB＝第二窗重新会判）。✓
- fresh 带 −48.68..−47.08（摆动 1.61）。门核：无 tag 四数里该测试**不编译**（SKIP=11
  集不含它）、门下两发全跑——门行为与申报一致。✓

## 4. AC#3 反向复刻＋还原核

- carry 臂均值 −39.73 vs reset −47.60＝**退化 7.87dB**——392 根因在绑定层隔离复现，
  Reset 移除之。✓
- 还原核帧差 0（仪器链全程保持）；tuned 会话 carry→Reset→刺激 与 fresh tuned
  **逐位一致**（注入 factory 被 Initialize 保留）＝调参面正面读数；旋钮同形效应 0
  如实记（合成刺激 enr 未上线）。✓

## 5. AC#4 零回归（口径差如实记＋判据落点核）

本腿现跑（导出树、`PATH=$PWD/third_party/sherpa-onnx` MSYS 形、-timeout 30m）：
- **四数无 tag：rc=0、`-v` 口径 RUN=98 FAIL=0 SKIP=11、顶层 PASS=75**——与
  `392r1-session-integration.md:108` 自记口径「RUN=98 FAIL=0 SKIP=11」**逐位同**；
  SKIP 集＝TestLiveWasapiSmoke＋aecspike 十枚硬件/工件门（本腿日志中一枚门名被
  密掩占位符替换，枚数与集形吻合）。✓
- **派单基线 94/83/0/11 口径差**：其源＝A1326 时代 389-v1 终裁行（台账 :29285，
  10-31）的旧测试总体口径；与本次两口径（98/0/11、顶层 86/75/11/0）均不对齐——
  实现已如实记为口径差；**判据落点（FAIL=0＋rc=0＋SKIP 集同形＋392-r1 `-v` 口径
  逐位复现）本腿独立复现成立，零回归判成立**。✓
- gofumpt `-l internal/audio/`＝0 脏；build 无 tag rc=0、vet rc=0、GOOS=linux vet rc=0。
- d22scan：selftest **RUN=77 PASS=35 FAIL=0 SKIP=0**；实扫 clean，分母
  internal/=287、cmd/=52（ban#1-5）、730/182（ban#8）逐枚命中；ban#8 未扩 docs。✓
- cmd/wisp 的 392 K=4 重跑不在本腿写面（调参面权威读数仍 392-r1）；零触碰＋tuned
  过 Reset 逐位存活已覆盖「调参面不退化」的本次判据。

## 6. ★主新发现核与处置建议（给编排者）

**可观测差＝坐实**：fresh vs reset 帧差 0.0411、读数带差 ≤1.7dB、`default-fresh`
唯一异形（本腿复跑逐字节复现；default-reset ≡ tuned-reset ≡ tuned-fresh 由
dTuned=0/dKnobs=0 坐实）。
**机制归属＝未证（本腿反证在案）**：证据件 §4 的「ApplyConfig 留 stale
`submodule_states_` ⇒ 分带门不开 ⇒ 生产吃全带（持续性）」依赖四条源码事实（本腿
逐条核真：ApplyConfig :709-711 只局部重建、InitializeLocked :517/:591 跑 Update、
门组成＝纯 `echo_controller_enabled_` 位 :378-383、factory 预标 :476-478）——但
`MaybeInitializeCapture`（:878）在**每次** float ProcessStream 都先跑
`UpdateActiveSubmoduleStates()`，stale（false→true）即触发全量
`InitializeLocked`（:889-897）；任何格式变更 reinit（`MaybeInitializeRender`
:501-513）同样同步。⇒ **stale 窗最迟在首个处理帧前关闭**，「持续全带」叙事未证；
两形差的真机制（首帧内一次性 reinit 侧效应：render 历史丢弃/首帧 delay 状态重置）
待仪器化副树探针坐实，本腿未越界造探针。
**统一器事实**：default+Reset ≡ factory 形逐位（已证）⇒ 生产会话自首次 Reset 起
即落 Initialize 形；「fresh 形」只存在于 Create 后、首次 Reset 前的首窗，量级
0.0411/≤1.7dB＝启动暂态尺寸。
**处置建议**：① 维持本腿不动生产路径（对：改 `wap_new` 形＝移位 392 基线数）；
② 后续票引用 §4 时把归属降级为「两形实存、全带/分带标签未证」；③ 若要统一引擎形
或改 `wap_new` 为 factory/Initialize 形＝**先跑 392 K=4 对照＋先用仓外仪器化探针
坐实归属再动**（与实现落点同向、加了机制坐实前置）；④ 保温修复票无需动生产路径：
窗间调 Reset 即自动统一（第二窗起全落 Initialize 形）。

## 7. 五格逐条判＋收口＋残余

- **AC#1 ✓**（§2）／**AC#2 ✓**（§3）／**AC#3 ✓**（§4）／**AC#4 ✓**（§5，口径差
  如实在案）／**AC#5 纪律 ✓**（真父 diff＝恰好申报 6 枚 +718/−0；越界面 0 字节；
  票面/台账/D43/thresholds/golden 零触碰；无 push）。
- **可收口：是**（第 68 票候选；AC 框勾选与票面 Status 归编排者收口动作，本腿⛔改）。
- **残余四件**：① flush 半边公开面不存在（具名；未来 webrtc 升级带来公开 flush 时
  照 `wap_reset` 形加一法即可）②两引擎形机制归属未证（§6，仪器化探针待需要时做）
  ③旋钮同形效应 0＝合成材料局限，权威读数仍 392-r1 cmd/wisp K=4 ④392 K=4 重跑
  不在任何已跑写面（验收/后续腿可跑，非收口阻塞）。
- **保温修复票收窄＝成立**：APM 半可原地清干净（reset 后 −47.60＝cancelled 深度），
  问题收窄为 resampler/detector 半（392-r1 :36-39 具名残余）。
- 393-r1 死亡核：工件/证据件 0 命中（`393r1-*` glob 零）；派单存档件按规矩在盘、
  无 393-r1 工作 commit——「未留痕」在工件面成立。

## 附：本腿自跑尺（全部可重放）

`D:\tmp\wisp\20261101\393v1\`：tree＝`git archive 35084d0b0` 导出＋3 枚 sherpa DLL；
`aec393v1.test.exe`（webrtcap 两形状 build/vet/linux vet 全 rc=0）；logs/393v1-run{1,2}.log
（门下两发）；out/run{1,2}/*（工件，与 `D:\tmp\wisp\20261101\393r2\out\*` diff rc=0）；
logs/four-numbers.log（98/0/11/75）；logs/d22scan-*.log；logs/gofumpt.out（0 行）。
CGO_CXXFLAGS 用 `D:/` 盘符形（MSYS `/d/` 形会被 cgo 当相对路径——本腿踩坑具名，
aec3r1 处方即 Windows 形）。
