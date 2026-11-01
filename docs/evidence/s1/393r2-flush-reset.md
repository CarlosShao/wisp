# 票 393 落地腿 `393-r2` —— AEC3 绑定层 flush/reset API（写腿）

日期 2026-11-01 · 腿 `393-r2`（写、`internal/audio/aec/**`＋本件）·
起手 HEAD `20392608e`（dev，开工三核过：①锚点在 ②索引净（脏区全 `.scratch/**`＋
`design/doubao`，写面零脏）③进程表现核——`wisp392v2.test` 在飞，音频错峰处置＝
合成读数（负载不敏感、两跑逐位复现）照跑、负载敏感的四数等 rig 退场后跑）
· 派单存档 `.scratch/wisp/probes/orch/20261031/txt/393-r1-dispatch.txt`（重派件：
393-r1 过夜静默死亡，本腿按同一票面接手）。

## 0. 一句话

**停手条款裁定：reset 半边有公开面、flush 半边没有**——钉定树
（webrtc-audio-processing 2.1，`D:\tmp\wisp\20261031\aec3r1\`）现核：
`EchoControl` 接口无任何 reset/flush；APM 公开面
`AudioProcessing::Initialize()`（`api/audio/audio_processing.h:497-507`）＝
「Initializes internal states, while retaining all user settings. This should be
called before beginning to process a new audio stream.」——这正是保温场景
（公告窗之间）的语义，落地为 `wap_reset`＋Go `Reset()`；**保滤波器、只清延迟
缓冲的 flush 无公开面＝具名**（引擎内部的 `delay_controller_->Reset`／
`kBufferFlush` 都不是公开 API；`kPlayoutVolumeChange` 只置
`echo_path_gain_change`＝「加速重适配」标志，不清滤波器不清缓冲）。不做伪造
的 Flush API。**K=4 合成读数（AC#2/AC#3）**：跨窗 carry 实测退化
**7.87dB**（392 根因在绑定层隔离复现），`Reset()` 后 **与「无历史 Reset 直读」
逐位一致（4/4 发，帧差 0）**、回声残余回到 cancelled 深度（-47.60 均值），
还原核逐位 0；tuned 引擎过 Reset 逐位存活。**额外发现（具名、未改）**：
钉定树 `ApplyConfig`-after-`Create` 路径留下 stale `submodule_states_`——
生产默认路径（`wap_new`）的引擎吃全带、Reset/factory 路径的引擎吃分带，
帧级差 0.0411、读数带差 ≤1.7dB（§4）。

## 1. 停手条款裁定（flush/reset 面有无，逐面现核）

对钉定树（deps.toml `[webrtc-audio-processing]` 2.1；源树
`D:\tmp\wisp\20261031\aec3r1\webrtc-audio-processing-2.1\webrtc\`）逐公开面核：

1. **`EchoControl`（`api/audio/echo_control.h:23-62`）**：方法集＝
   `AnalyzeRender` / `AnalyzeCapture` / `ProcessCapture(…, level_change)` /
   `GetMetrics` / `SetAudioBufferDelay` / `SetCaptureOutputUsage` /
   `ActiveProcessing`。**无 reset、无 flush。**
2. **`AudioProcessing`（`api/audio/audio_processing.h`）**：`Initialize()`（:507）
   ＝全状态重建、保留用户设置——**reset 面存在且文档语义正对保温场景**。
   `InitializeLocked()`（`audio_processing_impl.cc:516-588`）重分配全部
   AudioBuffer、`AllocateRenderQueue()`、逐子模块重建，其中
   `InitializeEchoController()`（:1880-1916）经 `echo_control_factory_` 重造——
   **tuned 会话重造同参数 tuned 引擎**（393-r2 读数 §3 逐位证实）。
3. **`RuntimeSetting`（:388-400）**：类型集
   `kCapturePreGain/kCaptureCompressionGain/kCaptureFixedPostGain/
   kPlayoutVolumeChange/kCustomRenderProcessingRuntimeSetting/
   kPlayoutAudioDeviceChange/kCapturePostGain/kCaptureOutputUsed`。
   逐型核（`audio_processing_impl.cc:817-835, 1007-1068`）：
   - `kPlayoutAudioDeviceChange` 在此版**只进 render 侧队列**、只转发
     `render_pre_processor`（本构建不存在）——**到不了 echo controller**；
   - `kPlayoutVolumeChange` 只更新 `capture_.playout_volume`，使下一帧
     `echo_path_gain_change=true`（:1308-1311）——该标志进
     `EchoPathVariability`（`block_processor.cc:132-145`），消费点是
     `refined_filter_update_gain`／subtractor／aec_state＝**「加速重适配＋
     残余估计更保守」，不是清态**。
4. **引擎内部 flush**：render overrun 时 `kBufferFlush`＋
   `delay_controller_->Reset(true)`（`block_processor.cc:136-144`）、
   capture 起点重置——全部内部，非公开面。

**裁定（铁律：最有价值；技术选型 agent 自拍）**：
- **reset 半边＝落地**：`wap_reset`→`apm->Initialize()`，Go 侧 `Reset()`。
  它正是 392-r1 根因链「APM 跨窗携带失准」所需之力（§3 读数证实）。
- **flush 半边（保滤波器、只清延迟缓冲）＝公开面不存在，具名**——派单停手
  条款①对 flush 半边成立；对 reset 半边不成立。不伪造 Flush API
  （`kPlayoutVolumeChange` 冒充 flush＝假音量假回声路变化，仪器不诚实；
  重命名 Initialize 为 Flush＝名不副实）。
- 停手条款②（shim 必须动 AEC3 源）**未触发**：AEC3 源树零触碰，
  shim 只增 `wap_reset` 一个声明/实现。

## 2. AC#1 落地面（对照 `wap_new_tuned` 形状）

| 件 | 改动 |
|---|---|
| `internal/audio/aec/wap_shim.h` | ＋`int wap_reset(void* session)` 声明（注释＝公开 Initialize 面、保用户设置、tuned factory 重造同引擎） |
| `internal/audio/aec/wap_shim.cpp` | ＋`wap_reset` 实现：`s->apm->Initialize()`，rc≠0 记 `wap_last_error`，返回 rc（与 process 函数同错误约定） |
| `internal/audio/aec/aec_webrtc_windows.go` | ＋`func (a *Apm) Reset() error`（doc 具名 flush 半边不存在、指向本件） |
| `internal/audio/aec/aec_stub.go` | ＋`Reset() error` stub（`ErrUnavailable`，方法集镜像不变） |
| `internal/audio/aec/aec393r2_reset_windows_test.go` | 新：AC#2/AC#3 测量件（门 `WISP_AEC393R2=1`，spike-gate 形状；工件写 `WISP_AEC393R2_OUT`） |

`aec.go`（Tuning＋Overrides）与 `NewApmTuned` 零触碰（392 调参面语义原样）。

## 3. AC#2 状态清洁读数＋AC#3 反向判据（K=4，四臂＋还原核）

**形状**（对照 392-r1 的 K=4：每设置 4 发、post-AEC episode RMS、min/max/摆动、
工件 JSON；392 的原测量在 `cmd/wisp`（⛔ 本腿写面），故在包内自足复刻同一
**测量形状**，刺激全合成确定性——跨臂比较是逐位而非摆动带）：
announcement（合成多音节调幅音，far 峰值 -10 dBFS）→ matched echo
coupling 0.2（近端回声 -24 dBFS，`set_stream_delay_ms(30)`）＋ talker
-26.70 dBFS（392-r1 L1 控制电平）@1.2s；echo-only 读数窗 0.8-1.15s
（冷启收敛后）、episode 读数窗 1.2-1.5s；每发 lag 抖动 2ms（480/512/544/576
样本——392-r1「远端 tap 与采集起点交错每公告抖动」的隔离读数）。

**carry 窗（392 根因复刻）**：render 喂满 3.4s、capture 2.0s 截断（窗末截断
留无配对渲染帧），同一 Apm 直接进下一窗。

**读数**（`-test.count=2` 两跑逐位同；工件 `393-aec-reset.json`、
`393-aec-reset-tuned.json`）：

| 臂 | echo_resid（4 发） | 摆动 | episode RMS 摆动 | 帧差 |
|---|---|---|---|---|
| **fresh**（冷对照，生产 `wap_new` 形） | -48.68 .. -47.08 | 1.61dB（＝lag 抖动敏感度） | 0.21dB | — |
| **carry**（截断→同 Apm 下窗，无 reset） | -40.35 .. -39.05 | 1.30dB | 1.02dB | vs fresh 0.047 |
| **reset**（截断→Reset→下窗） | -48.30 .. -46.94 | 1.36dB | 0.74dB | vs rebuild-control **0** |
| **rebuild-control**（Reset 无历史直读） | -48.30 .. -46.94 | 1.36dB | 0.74dB | vs fresh 0.0411 |

判定：
- **AC#2 状态清洁（精确定义）**：`reset == rebuild-control` **逐位
  （4/4 发，帧级 max diff = 0）**——Reset 后的会话与「从未携带任何东西、
  走过同一 face」的会话不可区分。carry 状态被完全清除。
- **域锚**：reset 臂 echo_resid 全部 ≤ -39（比 -24 近端回声低 15dB）＝
  「第二窗重新会判」形状成立——392 finding（episodes=0、残余骑近端电平
  -22~-25dBFS）在本面消失。
- **AC#3 反向判据**：carry 臂均值 -39.73 vs reset 臂 -47.60＝**退化
  7.87dB**——跨窗 carry 在绑定层隔离复现（392-r1 全链的「彻底失准」在
  binding-only 情形表现为 ~8dB 部分退化；全链更重＝resampler/det 半的
  叠加，那是 392-r1 的具名残余），且 Reset 移除之。
- **还原核**：全部臂跑完后全新会话复读 fresh[0]——帧差 0，仪器链全程保持。
- **tuned 存活**（392 调参面不退化的正面读数）：tuned 会话（enr 0.45）过
  carry→Reset→刺激，与 fresh tuned 会话**逐位一致**——注入的
  `EchoControlFactory` 被 `Initialize()` 保留并重造同参数引擎
  （`393-aec-reset-tuned.json`）。旋钮效应读数刻意取同引擎形对照
  （reset vs reset）：**本合成刺激上 0**（enr 阈值未上线——三旋钮都是
  状态机依赖的；调参面的权威读数仍是 392-r1 的 cmd/wisp K=4，本腿不能跑）。

**保温真解可行性读数（本票的交付问题）**：reset 后第二窗残余 -47.60 均值
＝cancelled 深度⇒**APM 半可以在原位清干净**，保温修复票把问题收窄为
resampler/detector 半（392-r1 §6.1 的具名残余）。

## 4. 发现（具名、未改）：钉定树的两引擎形（stale `submodule_states_`）

读数中 fresh vs reset 的帧级差 0.0411 追到根：

- `wap_new`＝`AudioProcessingBuilder().Create()`（构造：echo 关）**之后**
  `ApplyConfig(echo on)`——而 `ApplyConfig`（`audio_processing_impl.cc:668-743`）
  对 echo_canceller.enabled 变化**只局部重建**（:709-711 仅
  `InitializeEchoController`），**不跑 `UpdateActiveSubmoduleStates()`**；
- `submodule_states_` 于是停在构造时全 false → `ProcessStream` 的分带门
  （:1331-1335 `CaptureMultiBandSubModulesActive()`）不开 → **AEC3 吃全带**；
- tuned 路径（factory 在构造时预标 echo_controller_enabled，:477-478）与
  任何走过全量 `InitializeLocked()` 的会话（含本腿 Reset）→ 状态同步 →
  **AEC3 吃分带**。

两种引擎形读数带差 ≤1.7dB（本刺激）、帧差 0.0411（≈-27.7dBFS 单样本）。
探针核（同二进制）：`default-fresh` 是唯一的全带形；`default-reset ≡
tuned-reset ≡ tuned-fresh` 逐位。

**⛔ 本腿不改**：把 `wap_new` 改成 factory/Initialize 形＝改生产引擎数值
（392 基线带 -41.28..-28.12 是全带形上测的，392 调参面语义随之动）＝越界。
读数在案供验收腿/编排者裁量；若未来要统一引擎形，先跑 392 K=4 对照再动。

## 5. AC#4 零回归

- **build/vet 两形状** rc=0（webrtcap 面 CGO 旗标＝aec3r1 build.py 的
  COMMON define 全集＋`-include cstdint`＋c++17——aec3 头链进 shim 后
  必需；证据件配方里的 set 行只列了 -I，此处把全量记下：`-DWEBRTC_LIBRARY_IMPL
  -DWEBRTC_ENABLE_SYMBOL_EXPORT -D_WINSOCKAPI_ -DNDEBUG -DWEBRTC_WIN
  -D_WIN32 -D__STDC_FORMAT_MACROS=1 -DNOMINMAX -D_USE_MATH_DEFINES
  -DWEBRTC_ENABLE_AVX2 -DWEBRTC_APM_DEBUG_DUMP=0`）；GOOS=linux vet rc=0。
- **gofumpt**（`D:\work\base\gopath\bin\gofumpt.exe`）`-l internal/audio/`＝空。
- **四数**：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test
  ./internal/audio/... -count=1 -v`（无 tag；392-v2 rig 与六族退场后错峰跑）
  ＝ **rc=0，两包 `ok`（internal/audio 17.3s＋internal/audio/aecspike），FAIL=0，
  SKIP=11**（SKIP 集＝TestLiveWasapiSmoke＋aecspike 十枚硬件/工件门，与
  aec3r1/59wr1 记录的门集同形）。**-v 口径（`^=== RUN` 计数，Go 的子测试
  RUN 在列 0）＝RUN=98 FAIL=0 SKIP=11——与 392-r1 自记口径逐位同**；
  顶层终态口径＝RUN=86 PASS=75 SKIP=11 FAIL=0（12 枚为子测试终态）。
  派单基线写作 94/83/0/11——与本次两口径都对不齐（94≠98≠86、83≠75），
  如实记为口径差：判据落在 **FAIL=0、rc=0、SKIP 集、以及 392-r1 的 -v 口径
  逐位复现**上，零回归成立。
- **d22scan**（Git Bash `scripts/d22scan.sh`）：self-test **PASS=35 FAIL=0
  SKIP=0（RUN=77）**；实扫 clean（ban#1-5 internal/=287、cmd/=52；
  ban#8 internal/=730、cmd/=182＝392-r1 的 181＋本腿测试件 1）。首跑曾 2 红＝
  本件未落盘时的 [phantom-citation]（ticket 212 ban #9，aec3r1 同形状），
  本件落盘后复跑转绿。
- **392 调参面不退化**：`Tuning`/`NewApmTuned`/`wap_new_tuned` 零触碰；
  两形状 build/vet 过；tuned 引擎过 Reset 逐位存活（§3）；cmd/wisp 的 K=4
  重跑不在本腿写面（验收腿可跑）。

## 6. 写面申报（commit pathspec 逐枚）

仓内：
1. `internal/audio/aec/wap_shim.h`（M：wap_reset 声明）
2. `internal/audio/aec/wap_shim.cpp`（M：wap_reset 实现）
3. `internal/audio/aec/aec_webrtc_windows.go`（M：Reset 方法）
4. `internal/audio/aec/aec_stub.go`（M：Reset stub）
5. `internal/audio/aec/aec393r2_reset_windows_test.go`（新：测量件）
6. `docs/evidence/s1/393r2-flush-reset.md`（本件）

仓外 `D:\tmp\wisp\20261101\393r2\**`（aec.test.exe、build/vet/run 日志、
out\393-aec-reset*.json），只建不删。

⛔ 越界核：`cmd/wisp/**` 零触碰（382-w-r2、392-v2 在飞）；AEC3 源树零触碰
（面核只读）；`internal/audio/bargein.go`（389 门本体）零触碰；392 调参面
（`aec.go` Tuning／`wap_new_tuned`）零触碰；thresholds/golden/D43 表/票面/
台账/specs/PLAN/deps.toml/testdata 仓根/frontend/design/scripts/.github
零触碰；⛔ push（推送归编排者）。

## 7. 残余与边界（如实记账）

1. **flush 半边不存在（具名）**：保滤波器、只清延迟缓冲的公开面在钉定树
   没有。若未来 webrtc 升级带来公开 flush，绑定层照 `wap_reset` 形状加一个
   方法即可；在此之前保温修复票只有 reset 一档可用（§3 读数＝够用档）。
2. **两引擎形发现（§4）**：392 基线带的可比性注记——基线在全带形上测得；
   任何把生产路径改到 factory/Initialize 形的动作都会移位读数，先跑 392
   K=4 对照。处置摆编排者/后续票。
3. **旋钮同形效应 0＝合成材料局限**：三旋钮状态机依赖，本刺激未上线；
   调参面语义未动，权威读数仍是 392-r1 cmd/wisp K=4。
4. **测量件自带门与工件根**（`WISP_AEC393R2=1`、`WISP_AEC393R2_OUT`）：
   无 tag 四数永远 SKIP 不见它；门下两跑逐位复现，负载不敏感。
5. 393-r1 死亡未留任何盘上痕迹（`393r1-*` 证据件/glob 零命中），本腿为
   全新落笔；派单面的开工三核按重派件执行。
