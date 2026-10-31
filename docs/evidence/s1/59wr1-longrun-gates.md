# 票 59 重开首项测量腿 59-w-r1 — 长窗 rig 扩展与 (a)(b) 两条 PASS 条件读数

日期 2026-10-31 · 腿 `59-w-r1`（写腿＝测量装备＋测量；写面 `internal/audio/**` ＋ 本件）·
起手锚 `d49f0d6d5`（dev；开工三核过：①HEAD 在 dev ＝ `d49f0d6d5` ②索引有**他人**在飞暂存件
`.scratch/wisp/issues/28-session-scope-warm-done.md`（在飞 5 写 3 常态，本腿 commit 全程显式
pathspec）③写面 `internal/audio` 干净、`docs/evidence/s1` 仅他人未跟踪 381/382 判语不碰）。
接续：`59-v2-verdict.md` §9 两条 PASS 条件 ＋ `aec3r1-remeasure.md` §7 复跑处方。
**PASS/NO-GO 判语摆编排者/59-v2——本件只给读数与设计输入。**

## 0. 一句话结论

**rig 长窗能力已建成并跑了 4 发（60/60/60/120s）**：长窗下「塌峰」不再以 59-v2 那种「整窗
8.9dB 不收敛」的形态出现；曲线呈「干净稳态平台＋瞬态事件＋秒级恢复」，事件**以采捕侧外源
能量（endpoint mix 污染）归因为主**——music 发用确定性真参考把「rig 时延阶跳」假设直接否定
（时延全程恒 493 样本），但 music w24 形（采捕纯回声、32/32 活动帧、ERLE 骤降 8.2 后 2 窗
恢复）证明**引擎侧再收敛扰动也真实存在，两类机制并存**。干净段稳态平台实测在场：60a 连续
32 窗（=32s）26.1–30.9dB（中位 27.4）；30–31dB 窗 8 枚在场但不连续，**「≥30dB」按平台中位
读未达、按窗在场读达——两层读数并报，裁决摆编排者**。组合门 (b)：**−45dBFS 能量阈下塌峰态 tier-3 仍误触发
11.5/10min（好峰 5.4/10min），门在塌峰态兜不住**；−30dBFS 阈下本样本两态归零但代价是
−36dBFS 安静语音不可达（正控 −36 那发被能量门吃掉）。塌峰表征：亚样本相位假设被数据否定
（好峰/塌峰 fractional lag 无分离）；最小复现条件＝「测量窗内外源会话在放音」；在跑塌峰
检测判据三件套（滑窗 ERLE 骤降报警／采捕自相关@循环周期／对齐强度）已实现并全部分析过。

## 1. 环境与复用（⛔ 重造边界全部遵守）

- 引擎/构建产物：`D:\tmp\wisp\20261031\aec3r1\build\libwebrtc_audio_processing.a`（aec3-r1
  钉定成品，零重建、零源树改动、零构建参数改动）；`deps.toml` 两节未动一字。
- 绑定层：`internal/audio/aec/`（`-tags webrtcap`，CGO 面照 `aec3r1-remeasure.md` §7）。
- rig：`internal/audio/aecspike/`（59-r1/aec3-r1 两代全部仪器内建复用：SyscallN XMM、24H2
  denormal 重标定、PATH 遮蔽、音量轴、音钟推导停止条件）。
- 音频错峰：每次 rig 前现核进程表（go/aecspike/wisp 测试进程零在场才开跑；本腿全部 4 次
  rig 启动前各核一次）。
- 采捕内容两种：matcha TTS「你好，我是威斯普，今天天气不错。」＋确定性合成 music
  （`SynthMusic`，8s 无缝循环）；音量 0.5（=基线协议）；endpoint＝console 默认 render
  （48kHz/2ch/float32 mix，denormal 重标定路径在工作）。

## 2. rig 长窗扩展（交付件申报见 §7）

新增 `RunLoopbackLong`（`longrun_windows.go`）：

- **≥30s 连续 render→capture**：far fixture 无缝循环，音钟推导目标帧数
  （`Seconds × mixRate`）双到即停——**⛔ 墙钟超时实现**（停止条件只有音频时钟；设备事件等待
  与基线 rig 同形）； runaway/idle 守卫照搬基线。
- **分段落盘防内存爆**：5s/段 WAV（mix-rate mono int16，480KB/段）边采边写；60s 发内存
  峰值＝单段＋夹具；120s 发＝24 段共 11.5MB 落盘。AEC 输出同形落盘（16k，5s/段）。
- 离线过带（`longAecPass`）：**每发冷启动全新 AEC3 引擎**（产品形状）；流式 resampler
  （`audio.NewResampler`，分块喂与整喂逐位相同——包内既有断言）分段喂；
  far 参考由 `FarClock` 相位连续重建（fix-point 3×2^32@48k→16k，无漂移）。
- **在跑塌峰检测判据**（`longrun.go`，全部纯音频时钟量、因果可算）：
  ① `ERLEAccum` 滑窗 ERLE（32ms 帧×1s 窗，far-active 才计，窗闭合即出数）；
  ② `CollapseDetector` 骤降报警（基线中位 −6dB 或绝对 <15dB，持续 3 窗，一次/streak；
  阈值组由 59-v2 双峰 8.9 vs 29–30.7 的分离度定）；
  ③ `FractionalLag` 亚样本时延（抛物线内插，表征用）；
  ④ 采捕自相关@循环周期（probe 仪器，免参考污染检测）。
- 纯数学层（`longrun.go`/`longrun_test.go`）无 build tag，CI 全平台可算（8 枚单测见 §6）。

## 3. (a) 长窗 ERLE 全程曲线读数

4 发（TTS×3 ＋ music×1；音量 0.5；每发冷启动引擎；lag/强度对齐在头部 3.4s＝基线协议形）：

| 发 | far | 时长 | lag(str) | 起步段[0,5s) | 稳态段[末30s] | 全程 | 头窗[0,3.4s) | 收敛窗 | 报警 | 事件剖析 |
|---|---|---|---|---|---|---|---|---|---|---|
| long-tts-60a | tts | 60s | 493(0.988) | **27.4** | 18.5 | 20.3 | **26.0** | 0 | 1 | 干净平台 w13–44；w45–48 主事件＋零星小坎 |
| long-tts-60b | tts | 60s | 493(0.650) | 6.5 | 11.8 | 10.6 | 5.3 | 3 | 2 | **头部即污染**（强度 0.65）整发低迷 |
| long-music-60 | music | 60s | 493(0.954) | 23.0 | 17.5 | 16.2 | 24.0 | 0 | 2 | 8 事件，事件间平台 25.9–29.3 |
| long-tts-120 | tts | 120s | 493(0.981) | 22.2 | 11.5 | 11.5 | 22.4 | 0 | 6 | 42 事件/120s＝21/10min 污染率 |

（稳态段读数低是因为「末 30s」机械切片把事件窗算进去了——干净段统计才是稳态平台：）

**干净段（≥25dB）统计**：60a 40/58 窗达标、**最长连续 32 窗（w13–44＝32s）26.1–30.9、
中位 27.4**；music 21/58 窗、最长连续 8 窗中位 28.1；120s 发 18/117 窗（事件密集）；60b
10/58（头部污染毁整发）。

**事件归因（load-bearing，三仪器互证）**：

1. **music 真参考逐窗时延轨迹**（`SynthMusic` 确定性⇒far 逐位可重建；probe 自证 corr[0]=
   0.984）：**lag 全程恒 493（60 窗无一阶跳）**，而 far-strength 在多数 ERLE 事件窗同步塌
   （w5: 0.999→0.615/ERLE 8.2；w25: 0.669；w49: 0.791）。⇒ 采捕时延从未移动（**rig 时延阶跳
   假设被否定**），这些窗的采捕内容里出现了与 far 不相关的外源能量（**endpoint mix 污染
   实锤**）。关联是方向性的、非阈值干净：每个 ERLE≤12 的重事件窗 strength ≤0.67；但
   w49（strength 0.791）ERLE 仍有 23——外源音量小/窗内占比低时 ERLE 未必塌。
2. **引擎侧事件真实存在（反例一枚，如实报）**：music **w24：far-strength 0.953、32/32
   far-active 帧（估计方差排除）、ERLE 8.2、w26 已回 22.4/w27 回 27.1**——采捕保持纯回声、
   时延恒定，引擎自行丢失消除约 2s 后恢复。机制（前期轻度污染扰动滤波器后的再收敛，或
   AEC3 内部时延估计在 8s 周期内容上跳拷贝）本腿无法二分；**产品的 far（真 TTS）无这种
   8s 周期性，该形状是否随内容消失未知**。这一枚是 §5「在跑检测判据必须进产品」的直接
   证据。
3. **采捕自相关@循环周期**（免参考，TTS 臂唯一可用判别器）：四发合计 75 个 ERLE<20 窗中
   25 个 selfCorr 0.33–0.79（外源能量明确在场），50 个 0.81–0.98；music 地面真值标定：
   far-strength 0.92 的**轻度污染即足以把 ERLE 打到 15**（w12–13: strength 0.924/0.932 →
   ERLE 15/15）⇒ selfCorr 0.8–0.95 的段多数是轻度污染形状，0.8 阈偏松（如实记）。
4. **事件-边界距离**：事件起点到 5s 分段落盘边界的距离均匀分布（0/1/2s 各占 ~1/3）⇒
   **分段写盘 flush 造成毛刺的假设被否定**（冷跑无 flush 也见到冷启动形低窗）。

**对 PASS 条件 (a) 的回答**（判语摆编排者）：

- **塌峰形态在长窗下消失**：59-v2 的「整窗 8.9dB 不收敛」形态在 4 发长窗里零出现；取而代之
  是「瞬态事件（1–4s，ERLE 3–12dB）＋2–6s 恢复」，且逐发归因到采捕侧外源能量。
- **稳态平台在场**：干净段 30s+ 连续 26.1–30.9dB（60a）/ 25.9–29.3（music）；冷启动「起步段
  <稳态段」的担心不成立——干净发的起步段（26.0/24.0/22.4）与稳态段同带，起步即收敛
  （convWin=0）。
- **「≥30dB 级」两层读数**：按**平台中位**读＝27.4/28.1，未达 30；按**窗在场**读＝30–31dB
  窗在长窗上真实出现（60a 8 窗 ≥30：31.3/30.3/30.5/30.8/30.9/30.7/30.5/30.6），但非连续
  平台。两层并报，裁决摆编排者。
- ⚠ **污染是 loopback 测量环境的属性**：loopback 端点捕获该端点**所有** render 会话（桌面
  其它 app 的系统音/提示音全进采捕）；产品麦克风不这样。但它等效演示了 double-talk 形状
  （非回声内容进采捕→ERLE 暂跌→引擎恢复）——这正是产品必须存活的行为，且事件后恢复
  时间 2–6s 是「引擎在真 mix 里的再收敛速度」的第一手读数。

## 4. (b) 组合 barge-in 门产品级读数

门链＝产品同源：AEC 输出 → `speech.Detector`（silero＋Endpointer，0.5s 分块喂）→ utterance
episode（关闭事件自带音频）→ 三档：**VAD 单独 / VAD∧能量(−45dBFS 阈) / VAD∧能量∧ASR 文本
（非空转写）**。材料全部来自本腿盘上 AEC 输出段；另有正控两发（真 TTS 语音「今天天气怎么
样，帮我看一下日历。」以 −26/−36dBFS RMS 绝对电平加注到好峰材料上）。

**合并率（per state）**：

| 态 | 材料 | 总时长 | VAD 起点（tier-1） | tier-3 误触发 @−45 | tier-3 @−30 |
|---|---|---|---|---|---|
| good（长窗稳态臂×4＋冷跑好峰×5） | 9 | 222s | 6 ＝ 16.2/10min | **2 ＝ 5.4/10min** | **0** |
| collapse（长窗事件臂＋冷塌＋59v2 存档） | 44 | 104.4s | 6 ＝ 34.5/10min | **2 ＝ 11.5/10min** | **0** |
| control（正控语音 −26/−36dBFS） | 2 | 6.8s | 2 | −26 发全过、−36 发 @−30 被能量门吃掉 | 见左 |

**误触发样本（全文在 gates-probe.json）**：good 态 120s 臂 ep「你能」(−38.9dBFS)、「哎我」
(−32.3)；collapse 态 120s 事件臂「我们」(−35.6)、「听我」(−32.2)；**整流 ASR 在 collapse
态转出 3+ 字短语（「哎大龙」「九十五」）⇒ 按文本长度设阈救不了塌峰态**；good 态整流也有
伪迹（60a 稳态 48s 整流＝「哦不哦不不不不哦」，但 endpoint 化后仅开 1 个 episode、ASR「我」
rms −49.6 未过能量门）。

**对 PASS 条件 (b) 的回答**（判语摆编排者）：

- **−45dBFS 阈：门在塌峰态兜不住**（tier-3 仍 11.5/10min，是 good 态的 ~2 倍）；且 good 态
  也非零（5.4/10min）——59-v2 未定义即停 ③「组合门在好峰态也高误触发」的弱化变体成立，
  如实报。
- **−30dBFS 阈：本样本集两态 tier-3 归零**，且 −26dBFS 正控仍全过——**「能量门＋ASR 文本门
  在 −30 阈上把回声侧误触发清零同时保住正常说话电平」在本样本上成立**；但代价实测：
  −36dBFS 级安静语音被吃（正控 −36 发 @−30 归零），且 collapse episode RMS −32.2/−35.6 距
  阈值只有 2–6dB 裕量。⇒ 设计输入：静态能量阈不足以独扛塌峰态，必须配在跑塌峰检测
  （§5 三件套）做动态门控；−30 阈的 2–6dB 裕量在票 26 落地时必须重测加严。
- 能量阈的原始分布（episode RMS by state）全在 `gates-probe.json`，阈值可任意重切。

## 5. 塌峰表征（为何 2/4 发塌 22dB——重开首项设计输入）

1. **亚样本相位假设死掉**：本会话 6 发冷跑 fractional lag＝好峰 +0.27..+0.34、塌峰 +0.27，
   **无分离**（`cold-tts-*` 表如下）。59-v2 §4「落哪个峰由采捕实例的亚样本相位决定」在本腿
   仪器下不成立。
2. **本会话冷跑矩阵**（3.4s/3.9s，AEC 全程 ERLE＝aec3-r1 同尺）：

| 发 | far 头部 | 全程 ERLE | fractional lag | 头对齐强度 | 判定 |
|---|---|---|---|---|---|
| cold-tts-1 | TTS 即时 onset | 19.8 | +0.28 | 0.935 | good |
| cold-tts-2 | TTS 即时 onset | **12.4** | +0.27 | **0.814** | collapsed（头部污染形） |
| cold-ttsgap-1 | 0.5s 前导静音 | 17.6 | +0.29 | 0.745 | good（强度同样低——污染） |
| cold-ttsgap-2 | 0.5s 前导静音 | 18.8 | +0.27 | 0.983 | good |
| cold-music-1 | music | 25.3 | +0.31 | 0.977 | good |
| cold-music-2 | music | 26.3 | +0.34 | 0.974 | good |

3. **头部内容敏感性**：前导静音两发（17.6/18.8）与即时 onset（19.8/12.4）无系统差；
   music 头两发全好（25.3/26.3）——**「采捕头部内容敏感」假设未获支持**（n 太小，如实记）。
4. **最小复现条件**：「测量窗内 endpoint mix 有外源会话在放音」。实测桌面事件率＝长窗
   75 事件/300s（其中 ≥25 个重污染实锤、music w24 一枚判引擎侧）≈ **15 事件/10min 量级**；
   3.4s 冷窗被击中的推算概率 ~8%/窗。⇒ 59-v2 时代的 2/4 塌峰率（50%）高于本推算——或当时
   桌面更吵、或另有未捕获机制；**本腿无法回溯验证他们的 4 发**（无逐窗数据），如实报差异，
   不下断言。
5. **在跑塌峰检测判据（产品可用雏形，票 26 设计输入）**：①滑窗 ERLE 骤降报警
   （`CollapseDetector`：1s 窗、中位 −6dB/绝对 15dB、持续 3 窗；本会话全部事件被曲线捕获、
   长窗报警工作）；②采捕自相关@循环周期（**免参考**污染检测：selfCorr<0.8＝外源能量在场；
   对产品形态换算＝对远端参考的滑窗相关）；③对齐强度跌落（60b 头部 0.65＝整发污染即刻
   可见）。三者皆纯音频时钟量、因果可算、无需未来样本。
6. ⚠ 本腿 rig 的 `MixFar` 只在内存、未落盘 ⇒ 离线 probe 对 TTS 臂无法重建 far
   （matcha TTS 跨进程**不逐位确定**，实证 probe corr[0]=0.126 vs music 0.984）——时延
   轨迹在 TTS 臂只能用免参考仪器。后续腿建议 rig 落盘 far（5.8MB/60s，一行改动）。

## 6. 四数／两尺／d22scan 导出树＋正控

- **四数（`./internal/audio/...` 全集，文档化跑法＝PATH 前缀 third_party/sherpa-onnx；
  裸跑 0xc0000135 为 59-r1 时代既有形状）**：`-count=1 -v` ⇒ **RUN=83 PASS=73 FAIL=0
  SKIP=10**（rc=0，日志 `D:\tmp\wisp\20261031\59wr1\four-numbers.log`）。**逐名归因**（对
  aec3-r1 基线日志同名作差）：新增 12 枚全为本腿——PASS 8（TestFractionalLag…、
  TestERLEAccum×2、TestCollapseDetector×3、TestGateTiers、TestFarClock）＋SKIP 4
  （TestSpikeAec3LongWindow、TestSpikeLongWindow{Probe,Autocorr,DelayTrack}——门控测试在
  无 tag 形状按设计 skip）；removed＝0；SKIP 全集＝预期门控/活体硬件集合（10 枚逐名
  见日志）。53 PASS/6 SKIP 基线零丢失零变红。
- **两尺**：gofumpt（`D:\work\base\gopath\bin\gofumpt.exe`）`-l internal/audio/`＝空；
  `go vet` 三形状 rc=0（无 tag ＋ `-tags webrtcap` ＋ `GOOS=linux`）；build 两形状 rc=0
  （无 tag ＋ webrtcap 的 `./internal/audio/...`）。
- **d22scan（Git Bash `scripts/d22scan.sh`）**：实扫 **clean** rc=0；正控 selftest
  **PASS=35 FAIL=0 SKIP=0（RUN=77）**（本腿代码注释只引用已存在证据件，phantom-citation
  零咬合）。导出树＝bans#1-5 internal/=**285**、cmd/=50、ban#6 frontend/=116、
  ban#7 internal/tools/=25、ban#8 design/=39、frontend/=116、internal/=**723**、cmd/=174。
  对 aec3-r1 的 283/49/…/719/171：**internal +2＝本腿 `longrun.go`+`longrun_windows.go`、
  ban#8 internal +4＝本腿四件**；cmd/design 的 +1/0 差属两腿间其它在飞腿（本腿零写面）。
- 路径长度门：本腿最长相对路径
  `internal/audio/aecspike/longrun_windows_test.go`（45 字符＜100 帽）。

## 7. 写面申报（逐枚）

仓内（commit pathspec 逐名列出）：

1. `internal/audio/aecspike/longrun.go` —— 长窗测量纯数学：FractionalLag／ERLEAccum（滑窗
   ERLE）／CollapseDetector（骤降报警）／GateTiers／FarClock／Rate10Min（无 build tag）。
2. `internal/audio/aecspike/longrun_test.go` —— 上件的 CI 单测 8 枚（全平台）。
3. `internal/audio/aecspike/longrun_windows.go` —— `RunLoopbackLong`：长窗 rig（循环 far、
   分段落盘、音钟停止条件、 runaway/idle 守卫；`RatioPpm` 语义变化已注释声明）。
4. `internal/audio/aecspike/longrun_windows_test.go` —— 门控驱动 4 枚：TestSpikeAec3Long
   Window（A/B/C 三阶段主驱动）、TestSpikeLongWindowProbe（离线相关判别＋门测材料重建）、
   TestSpikeLongWindowAutocorr（免参考污染判别）、TestSpikeLongWindowDelayTrack（逐窗时延
   轨迹；music 真参考交叉验证）。
5. `docs/evidence/s1/59wr1-longrun-gates.md` —— 本件。

仓外（`D:\tmp\wisp\20261031\59wr1\**`，只建不删）：`run\`（测试 exe×2＋三 DLL＋四个 .cmd
跑法存档）、`out\`（每发 `<name>.json`+`-curve.json`＋采捕分段 WAV＋AEC 输出分段 WAV；
`longrun-summary.json`、`gates.json`、`probe.json`、`gates-probe.json`、`autocorr.json`、
`delaytrack.json`）、`four-numbers.log`、`d22scan.log`。音频 artifact＝合成 fixture＋自渲染
TTS＋其回声，无用户语音（D16）；全部在 WISP_AEC_SPIKE=1 门下写入临时根。

⛔ 未碰：`cmd/wisp/**`（30-w-r1 在飞）、`internal/speech/**` 语义（只调用既有 API）、
AEC3 源树/构建参数、deps.toml、票面/台账/specs/PLAN、testdata。

## 8. 复跑处方（逐字）

构建（aec3-r1 §7 处方，本腿落点）：

```
cd <repo>
set CGO_CXXFLAGS=-ID:\tmp\wisp\20261031\aec3r1\webrtc-audio-processing-2.1\webrtc -ID:\tmp\wisp\20261031\aec3r1\abseil-cpp-20240722.0
set CGO_LDFLAGS=-LD:\tmp\wisp\20261031\aec3r1\build
go test -c -o D:\tmp\wisp\20261031\59wr1\run\aecspike.test.exe -tags webrtcap ./internal/audio/aecspike/
copy third_party\sherpa-onnx\{onnxruntime.dll,sherpa-onnx-c-api.dll,sherpa-onnx-cxx-api.dll} D:\tmp\wisp\20261031\59wr1\run\
```

主驱动（rig 前**现核进程表**：go/aecspike/wisp 测试进程零在场；约 6 分钟，其中 300s 为
实时采捕）：

```
cd D:\tmp\wisp\20261031\59wr1\run
set WISP_AEC_SPIKE=1
set WISP_AEC_SPIKE_ROOT=<repo>
set WISP_AEC_SPIKE_OUT=D:\tmp\wisp\20261031\59wr1\out
set PATH=<repo>\third_party\sherpa-onnx;%PATH%
aecspike.test.exe -test.run TestSpikeAec3LongWindow -test.v -test.timeout 40m
```

离线 probe 三连（免 rig；须先跑主驱动留件）：

```
aecspike.test.exe -test.run TestSpikeLongWindowProbe -test.v -test.timeout 20m
aecspike.test.exe -test.run TestSpikeLongWindowAutocorr -test.v -test.timeout 15m
aecspike.test.exe -test.run TestSpikeLongWindowDelayTrack -test.v -test.timeout 15m
```

（无 tag 四数：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./internal/audio/... -count=1 -v`。）
