# 票 59 — P15 AEC spike 裁决读数（腿 59-r1）

日期 2026-10-31 · 腿 `59-r1`（写腿 AEC spike）· 起手锚 `dcf6e910`（dev，非 detached；开工三核过：
索引无 staged、`internal/` 零脏、cmd/wisp 的脏件属 41a2-r1 未碰）
写面＝`internal/audio/aecspike/**`（新 spike 子包，9 文件）＋ `docs/evidence/s1/59r1-*.md`（本件＋
license 判定件）。⛔ cmd/wisp、⛔ internal/speech 语义（只 import 读）、⛔ 票面/台账/specs/PLAN、
⛔ push、⛔ 下载。

## 0. 一句话结论

**PRECHECK P15 = NO-GO（本机证据范围内）**：票面点名的 webrtc AEC3 无 fetch 不可评估（⛔ 下载边界），
零下载候选（Windows 内置语音 AEC）在本机**不存在**（DMO 0x80040154＋APO 清单无语音 AEC），而 rig 实测
证明自回声在采集通路**真实存在且被 ASR 逐字转写**、单增益消参考不足以抑制 ⇒ Path C 全双工的
barge-in 在没有真 AEC 引擎时必然自触发。按票面格 7：**26/28 的 Path-C 部分与 60 建议随 NO-GO 自动
defer（登记 DEFERRED，摆编排者/机主）**；本判语随两个可翻转条件失效（见 §8）。

## 1. 环境与工具链（钉定核）

- Go **1.27.1**（toolchain 指令生效）；C 工具链＝**MSYS2 mingw-w64 GCC 16.2.0 (Rev3, MSYS2 project)**
  于 `E:\work\base\msys64\mingw64\bin\gcc.exe`（PATH 解析）——与 `docs/BUILD.md` §1 钉定版逐字同。
- `GOFLAGS=` 空、`GOTMPDIR` 未设、CGO 默认启用、GOPROXY=goproxy.cn（既有）。
- 音频面：Realtek Audio（默认渲染＋默认采集）＋Intel/NVIDIA HDMI 渲染；OS=Windows 10 Pro build
  26100（24H2）。主控台会话当时**锁屏**（LogonUI 在场）——音频栈照常工作。

## 2. 格 1 — License verdict → 判语在 `59r1-license-verdict.md`

乙（AEC3）＝本腿 NO-GO（无 fetch 不可核，⛔ 下载边界）；甲（Windows 内置）＝本机不存在
（0x80040154，APO 清单无语音 AEC）；丙＝无 license 面。**「license 文件入 evidence」一格在本腿
不可满足**（无 fetch 就没有真实 license 文件；凭记忆抄＝伪造，⛔），已如实具名。

## 3. 格 2 — Go binding 编译链接 ＋ cgo smoke（钉定工具链）→ **可评部分全过**

- **cgo smoke 一发实跑**：`aecsmoke.c`（Voice Capture DMO 驱动，C 约 450 行）由钉定 mingw gcc
  经 cgo 编译＋链接（仅 `-lole32`，全部 GUID/属性键本地钉定）→ `go build`/`go vet` rc=0 →
  运行时真调 `CoCreateInstance(CLSID_CWMAudioAEC)` → **hr=0x80040154**（类未注册＝候选甲缺席的
  实锤读数，不是编译失败）。cgo 面本身（编译/链接/运行/COM 调用）**全通**。
- 仓内既有 sherpa-onnx 绑定（go.mod v1.13.8，模块内自带 mingw 库）同样在钉定工具链上编译链接
  运行成功：silero VAD、streaming Paraformer ASR、matcha TTS 三引擎本腿全部真载入真推理
  （ASR 远端转写「你好我是威斯普今天天气」＝TTS/VAD/ASR 全链路活）。
- ⚠ 运行时 DLL 形状（给后续腿的坑）：System32 里驻有 **onnxruntime 1.17.1**（Windows 自带，
  加载器先于 PATH 命中）⇒ `go test` 直跑会在 `NewMatchaTts` 崩（0xc0000005）。本腿用
  **`go test -c` 出 exe＋钉定三 DLL 拷到 exe 旁** 的形状（127/131 票先例）解决。
  ⚠ **顺带上报**：`internal/speech/tts.go` 的 `NewMatchaTts` **没有** `checkOnnxruntime` 守卫
  （`NewDetector`/`NewStreamingAsr` 都有）⇒ ORT 被遮蔽时 TTS 直接崩而不是拒绝。⛔ speech 语义
  未动，此味归票 13/26。
- 交叉编译：`GOOS=linux go vet ./internal/audio/...` rc=0（本包纯 Go 面干净）；AEC3 的
  mingw-cross 可行性因候选不可 fetch **无法现核**，票面假设照抄为待验。

## 4. 格 3 — 回环 rig（WASAPI render TTS＋capture loopback）→ **真读数到手**

rig＝同一钉定线程上开共享模式 render＋同端点 `AUDCLNT_STREAMFLAGS_LOOPBACK` 采集，播 matcha TTS
「你好，我是威斯普，今天天气不错。」（22.05k→16k→48k 混音率）：

- **对齐**：lag=493 样本（30.8ms，共享引擎 render→loopback 通路延迟，物理合理）；
  归一化相关 **strength=0.989 / corr=0.978**（near 就是 far 的回声，判据成立）。
- **ERLE（单增益消参考，spike 的 AEC-cleaned 替身）＝16.8dB**。
- **自转写抑制＝不成立**：对回环采集信号跑产品同源 streaming ASR——
  raw＝「你好我是威斯普今天天气」（回声被逐字转写＝Path C 无 AEC 时 barge-in 必自触发）；
  16.8dB 扣除后的 residual 仍被逐字转写 ⇒ **仅消参考不够，必须上真 AEC 引擎（≥30dB 级）**。
- ⚠ 仪器自曝的本机怪癖（全部留痕在日志）：①回环 tap 把 int16 混音样本以**位型**塞进 float32
  （值＝样本×2^-149，全 denormal）——rig 检测该特征后按 2^134 重标定；②Go `syscall.SyscallN`
  传不了 XMM 寄存器 ⇒ 整型路径调 `ISimpleAudioVolume::SetMasterVolume` 曾把本进程会话音量写成
  ≈0 且被 Windows **持久化**（正是早先"回环全零"的根因）——rig 现以 cgo C 调真 float ABI 复位
  1.0；③System32 ORT 遮蔽（§3）。

## 5. 格 4 — 跨端点漂移 → 同端点读数到手；真跨端点本机不可测（如实具名）

- 同端点（D47 的主论证场景：参考＝自渲染 PCM）：渲染/采集同一硬件时钟，captured−pushed 实测
  +885～+2885ppm，且随采集尾巴波动 ⇒ **同端点无设备时钟漂移，残差＝引擎缓冲/采集尾部**；
  D47「小一个量级」的论证方向被支持（但见下条的真门槛）。
- **合成漂移敏感度**（对 far 叠加 0/25/50/100/200/400ppm 后再对齐扣除）：单增益扣除在 0ppm 干净、
  **200ppm 即崩**（ERLE<10dB，单元测试钉住）⇒ 真跨端点（蓝牙音箱＋有线麦）若无重采样/漂移跟踪，
  参考扣除类方案不可用。
- **真跨端点（蓝牙音箱＋有线麦）本机测不了**：只有一枚采集端点（Realtek）、无蓝牙渲染设备。
  owner 可跑命令＝`PATH="$PWD/third_party/sherpa-onnx:$PATH" WISP_AEC_SPIKE=1 WISP_AEC_SPIKE_ROOT=\
  <repo> go test ./internal/audio/aecspike/ -run TestSpikeSelfEchoPathRig -v`（在 BT 渲染为默认的
  机器上）——rig 代码已就绪。

## 6. 格 5 — 外部音源误触发表（9 格，silero VAD＝产品同源，合成 fixture 如实标注）

far 参考＝TTS 语音；采集＝端点回环（数字回环、无房间混响 ⇒ 全表为真实现场的**下界**）。
触发率单位＝次/10min（8s fixture 折算）；corr＝对齐后 near-vs-TTS参考 相关均值。

| fixture×音量 | raw 触发 | 残差触发 | corr |
|---|---|---|---|
| speech-tts 0.2/0.5/0.8 | 1/1/1（176.5×3） | 1/1/1 | **0.975/0.978/0.979** |
| music-synth 0.2/0.5/0.8 | 1/0/0 | 0/0/0 | −0.002 |
| video-synth 0.2/0.5/0.8 | 1/1/3（75–225） | 2/2/1 | 0.005 |

读数：
1. **自回声在三档音量全部触发 barge-in**，且扣除后仍触发（§4）——无 AEC 的 Path C＝「播报中
   自己打断自己」。
2. **外部源与参考的 corr≤0.014，自回声 corr≈0.98** ⇒ **corr 判别器干净分离，建议阈值≈0.5**：
   barge-in 门＝VAD ∧ corr(near, far_ref)<0.5（票面「direction of AEC residual」的可落地方）。
3. 合成音乐对 silero 基本**不触发**（0–1 次/8s，边界抖动）；**含语音的 video 触发且随音量上升**
   （0.8 档 3 次/8s）⇒ 「外部音源必须不触发」的 ≤1 次/10min 目标在音乐上达标、在**含人声的
   外部内容上不达标**——该结论为数字回线下界，真房间混响会更糟。
4. fixtures＝合成（chord+kick 音乐／含 TTS 摘段的 video）＋**真 matcha TTS**（speech 行）——
   合成部分如实标注，非真实媒体文件。

## 7. 格 6 — Sustained CPU（数留 evidence；SPEC-10 回填摆机主）

- 消参考＋对齐链（纯 Go DSP）：**≈0.00% 单核**（3.4s 音频毫秒级算完）。
- silero VAD（产品同源、含两端点状态机）：两遍过 6.8s 音频合计 **53.1% 单核**（单遍 ≈26%），
  窗口已剔除模型加载。**显著高于 SPEC-10 §3.1 Path C 行的 5% 暂定值** ⇒ 5% 那格不被本 rig 的
  VAD 读数支持。
- AEC 引擎本体 CPU：候选甲不在本机、乙不可 fetch ⇒ **本机测不到 AEC 引擎的真实 CPU**——Path C
  行的最终回填缺「真 AEC 引擎＋产品集成」的实测，**回填动作摆机主/票 13**（⛔ 本腿不改 SPEC-10
  一字）。

## 8. 格 7 — PRECHECK P15 = **NO-GO**（判语＋依据＋翻转条件）

**判语**：在本机证据范围内，Path C 全双工（播报中说话 ≤400ms 停播）**当前不可开工**：

1. AEC 引擎三候选全灭（§/license 件：乙要 fetch=⛔、甲本机不存在=探针实锤、丙=砍功能）；
2. rig 实测证明**无 AEC 的 Path C 必然自触发**（自回声 corr 0.978、ASR 逐字转写、三档音量全触发）；
3. 400ms 预算内的检测链（VAD-onset→stop）本 rig 已证明可跑（VAD 在回声上 32ms 窗粒度出-start，
   检测延迟受 MinSpeech 300ms 支配），但「停播」建立在 barge-in 判对的前提上——前提不成立。
4. **自动 defer 建议（按票面原文）**：26 的 Path-C barge-in 部分、28 的 Path C 全双工、60 的 C32
   realtime 门 —— 登记 DEFERRED（五字段），⛔ 静默消失。
5. **翻转条件（满足任一即翻 PASS-候选，需重跑本 spike 的质量/CPU/误触发三件）**：
   (a) 机主批准 fetch webrtc-audio-processing（BSD-3 上游声明待现核）→ 票 13 落 deps.toml＋
   license 文件 → 重测；或 (b) 在 OS 语音 AEC 在场的机器上（本机不在场）用 rig 实测 ERLE≥30dB 级
   且 VAD-CPU 有产品级数。

### 三形代价表（⛔ 本腿不选形，按 337/340 先例摆料）

| 形 | 内容 | 改动面 | 代价/风险 |
|---|---|---|---|
| 甲 | Windows 内置语音 AEC（OS 组件） | 票 13 的模块＋comms 模式采集 | 零下载零 license 面；**本机不存在**（DMO 缺席、APO=Realtek 厂商件）；每机可用性不一，须运行时探测＋降级 |
| 乙 | webrtc-audio-processing AEC3（票面点名） | deps.toml＋fetch-deps＋cgo 绑定封装＋票 13 模块 | BSD-3（上游声明）；要 fetch＝人工批准；跨编译存疑；质量/CPU 业界最优级 |
| 丙 | 不上 AEC：Path C barge-in 退化为热键/单击 | D47 白名单条目撤回（**改契约＝人工批准**）＋26/28/60 defer | 零新增依赖；代价＝Path C 头牌功能（播报中说话）不存在 |

## 9. 交付物与门禁（四数／两尺／d22scan 导出树＋正控）

- **写面**：`internal/audio/aecspike/`＝doc.go、aecspike.go、fixtures.go、aecspike_test.go、
  rig_windows.go、rig_windows_test.go、dmo_cgo_windows.go、dmo_nocgo_windows.go、aecsmoke.c/.h。
- **四数**：`go test ./internal/audio/... -count=1 -v` ⇒ **RUN=52 PASS=39 FAIL=0 SKIP=1**
  （SKIP＝既有 `TestLiveWasapiSmoke` live 硬件门，非新增；日志 `D:\tmp\wisp\20261028\59r1\four-numbers.log`）。
- **两尺**：gofumpt（`D:\work\base\gopath\bin\gofumpt.exe`）`-l internal/audio/`＝空；
  `go vet` windows＝rc=0、**GOOS=linux**＝rc=0（真类型面）。
- **d22scan**（Git Bash 跑 `scripts/d22scan.sh`）：正控 **PASS=35 FAIL=0 SKIP=0（RUN=77）**、
  实扫 **clean**，导出树＝bans#1-5 internal/=280、cmd/=48、ban#6 frontend/=116、
  ban#7 internal/tools/=25、ban#8 design/=39、frontend/=116、internal/=715、cmd/=169，rc=0。
- ⚠ **路径长度门 RED＝既有 offender**（票 233 工单文件名 105 字符，超 100 帽；本腿写面最长相对
  45 字符，不涉名册）——如实记账，⛔ 越界改票名。
- 临时件/artifacts（仓外 `D:\tmp\wisp\20261028\59r1\`）：inventory.json、selfecho.json、
  selfecho-{far,near,residual}-16k.wav、false-trigger.json、gated-*.log、four-numbers.log、
  aecspike.test.exe＋三 DLL。音频 artifact＝合成 fixture＋自渲染 TTS，无用户语音（D16）。

## 10. 复跑处方（逐字）

```
cd <repo>
go test -c -o D:\tmp\wisp\20261028\59r1\aecspike.test.exe ./internal/audio/aecspike/
copy third_party\sherpa-onnx\{onnxruntime.dll,sherpa-onnx-c-api.dll,sherpa-onnx-cxx-api.dll} D:\tmp\wisp\20261028\59r1\
cd D:\tmp\wisp\20261028\59r1
set WISP_AEC_SPIKE=1
set WISP_AEC_SPIKE_ROOT=<repo>
aecspike.test.exe -test.run TestSpike -test.v -test.timeout 12m
```

（ ungated 单测在 CI/任何机器：`go test ./internal/audio/aecspike/`，无设备自动 SKIP。）
