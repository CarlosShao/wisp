# 票 26 Path-C 集成面 — 落地腿 26p-r1 实现件（三个集成真数＋接线＋反向）

日期 2026-10-31 · 腿 `26p-r1`（写、cmd/wisp 独占＋internal/audio 接线面＋本件）·
起手 HEAD `41295d4d5`（dev，开工三核：锚点在／索引净／写面无他人脏）·
派单存档 `.scratch/wisp/probes/orch/20261031/txt/26p-r1-dispatch.txt` 在案。

## 0. 一句话

组合门（389）＋AEC3 引擎（59-w/aec3-r1 的绑定与成品库）接进 resident 语音路径的播报窗，
**三个集成真数全部在 staged 全链（真 waveOut＋真 matcha＋真 C8 帧＋真 AEC3＋真组合门＋真
silero＋生产 TriggerInterrupt→Abort 停播链）上出数**：(a) 播报中发声 ⇒ barge-in 触发停播
**99/56/67/65/40ms 五轮全 ≤400ms**；(b) 产品麦面零转写＋零自触发（含反向 M2 摘 AEC3 ⇒
**逐字转写复现**「我已经完成明天下午三点的会议改到四点请准时参加谢谢」＋还原核）；首格
**Feed 在线 vs 窗级重放 6 窗零漂移**；(c) -30 裕量重测加严读数（含反向 M1 摘组合门 ⇒
**回声自触发复现** -36.74dBFS ＋还原核）。

## 1. 接线（写面：cmd/wisp/**＋internal/audio/gate.go＋internal/audio/align.go）

- **`cmd/wisp/resident_voice_bargein_windows.go`（新）＝ `voiceBargeinMonitor`**：一个公告窗的
  barge-in 判决器。
  - far 参考＝公告本身：speak() 的 chunk 回调在进设备前把合成块交给 `speakFar`（流式重采样
    22050→16k，绝对坐标不漂移）——D47 的「回声参考是我们自己渲染的 PCM」。
  - 采集帧＝真 C8 seam：arm 时把 `audio.HalfDuplexGate.SetConversation(true)`（D47 会话窗，
    见下）打开，monitor 的 drain goroutine（C11 作用域、冻结临时名册名 `tts-infer`）成为公告
    窗内 `rv.frames` 的唯一读者——消费 goroutine 恰在同一括号内身处 speak()，单读者纪律由
    同一括号成立。
  - 判决链＝产品形状：AEC3 render-then-capture 10ms 锁步（时延提示在线测得——
    `internal/audio/align.go` 包络相关对齐器，rig 仪器 `AlignByEnvelopeCorrelation` 的产品形态，
    跨件对拍测试逐值同）→ `audio.BargeinGate.Feed`（在线）→ 自建 silero 会话分段**残余** →
    episode 起点走**真 Evaluate**。D32:2315 形状：播放期挂起 ASR、400ms 门禁由 VAD 而非 ASR
    满足——episode 不带 ASR 文本，Evaluate 仍是唯一阈值权威。
  - 过档 ⇒ 触发一次 `TriggerInterrupt`（Esc/球单击同一条打断通道），381-r2 的 abort hop 真停
    设备。
  - **退化诚实**：无 `-tags webrtcap` 绑定 ⇒ NewApm ErrUnavailable ⇒ monitor 不武装，公告窗
    保持今天的精确 D16 形状（PathT 播报闭麦）；无 VAD/无作用域/无触发器同理，reason 落读数。
- **`internal/audio/gate.go`：`HalfDuplexGate.SetConversation(on bool)`**＝D47 Path-C 会话窗：
  ON 时 PathT 门在 SetSpeaking(true) 期间不关采集（「echo cancellation is the AEC layer's job」
  ——gate.go 自己的 PathC 注释从第一wave就写着这句）；OFF（默认）＝逐字节回到 26p 前形状；
  muted 仍最高优先。单测 `TestGateConversationWindowOnPathT` 六跳钉死。
- **`internal/audio/align.go`（新）**：对齐器产品形态（rig 的包络相关时延测量），常数与算法同
  rig；`internal/audio/aecspike/align_cross_26p_test.go` 钉两实现逐值同（rig 不进产品二进制）。

## 2. 公告窗工作点＝收紧档 -30（集成面首个设计输入，实测钉住）

首次 staged 全链跑出一个 rig 给不出的形状：**健康匹配回声抵消下两路动态判据全读「可信」**
（far-corr 0.78-0.9 ⇒ 回声即参考；ERLE ~23dB ⇒ 无塌峰报警），loose 标称 -45 生效，而**引擎残余
地板（实测 -42.11dBFS）在它之上 ⇒ 回声自触发**（读数在案：episode rms -42.11 pass=true、
触发 1 次）。59-v3 §3 的「不自激」验证点本来就是 **tier-3@-30 字面达标**，且收紧是裕量交换
不是抑制开关（control-26 真话在收紧档仍过）——⇒ 公告窗的门以收紧档武装
（`EnergyThresholdDbfs: BargeinTightenedEnergyDbfs`），动态三件套照跑照记（读数面仍然活着）。
M1 反向＝回 loose 轴＋dynamicOff＝rig 的 11.5/10min 基线形状。残余的**连续性**还逼出了第二件
产品形状：AEC 残余永不给 endpointer 700ms 静音 ⇒ 首个残余 episode 永不闭合、说话者 onset 永不
产生 START（实测轮 3 判了一个窗 0 后全程静默）——⇒ monitor 分段器加**噪声门**（滚动 2s 最小值
地板×2＋512ms 保持，常数与发现引注进码）：平地板被门掉、有语音动力学的东西（说话者／M2 旁路
下的裸回声）通过，silero+Endpointer 在门后照跑（D32 的 VAD+AEC 形状保住）。

## 3. 三个集成真数（staged 全链读数；工件 `D:\tmp\wisp\20261031\26pr1\**`）

**(a) 播报中发声 ⇒ barge-in 触发停播 ≤400ms — `Test26PBargeinStopsWithin400ms` PASS**
五轮，真设备＋真 matcha 公告＋公告自身回声（匹配耦合 0.2）经 C8 注入＋-20dBFS 真近麦说话者
（真 matcha 渲染）+1.2s 注入：

| 轮 | episode rms (dBFS) | 阈 | VAD起点→停播 |
|---|---|---|---|
| 1 | -21.88 | -30 | **99ms** |
| 2 | -26.63 | -30 | **56ms** |
| 3 | -23.07 | -30 | **67ms** |
| 4 | -26.77 | -30 | **65ms** |
| 5 | -20.56 | -30 | **40ms** |

全部 ≤400ms（MET）。分解的诚实账：onset→VAD-start 一段是 MinSpeech 300ms 设计常数＋silero
自身时延（59-r1 §8「预算紧、设计空间在」的形状）；monitor 冷启动对齐预热（~0.6s，在线测
lag/强度随读数在案）由说话者 +1.2s 注入避开——稳态插话时延是上表读的量。

**首格 Feed 在线 vs 窗级重放 — `Test26PFeedOnlineMatchesWindowReplay` PASS**
回声-only 轮（无触发 ⇒ 6 窗全闭）：同一 (far,near,residual) 流，在线 `Feed` 与离线
`ObserveWindow` 的 per-window untrusted 判决**6/6 逐位同、零漂移**（389-v1 §3 边界①在此面闭合；
lag=644、strength=0.89）。工件 `26p-feed-equivalence.json`。

**(b) 产品麦面零转写（回声不触发）— `Test26PEchoNeverTriggersAndNeverTranscribes` PASS**
公告自身回声整窗注入（播报音频进麦克风回路＝S4 字面的注入形状）：episodes=1（残余 -42.88
被收紧档吃）、**triggers=0、finals Δ=0、handed Δ=0**（转写面结构性与实测双成立——公告窗帧由
monitor 独占，产品转写路收不到），gate 2/2，问题表空。工件 `26p-b-echo-zero.json`。

**(c) -30 裕量重测加严读数 — `Test26PMarginRetestAndGateReverse` PASS**（两发读数在案）
- **L1（control-26 -26.70 过链）**：post-AEC **-29.42（margin +0.58dB，过）** 与 **-38.24
  （margin -8.24dB，被吃）** 两发——**AEC 自身双讲抑制对控制电平的损伤（0-8dB，随引擎状态摆）
  是主导项，offline 标定的 3.3dB 裕量不保证过链**——这就是 (c) 要的加严读数，也是 28 全双工
  集成时的设计输入（rig 的 -26.70/-36.58 是混合电平，从未过引擎）。
- **L2（quiet speech -36.58）**：post-AEC -41.69，收紧档吃（如实代价）；离线标称轴 re-judge
  过（loose 档的可达性仍在）——「收紧是裕量交换不是抑制开关」在集成面复现。
- **L3（失配回声伪迹族，耦合 0.25）**：残余 -36.38，距收紧档 **6.38dB**——比 59-w rig 塌峰族的
  2.2-5.6dB 更宽（rig 的塌峰是引擎失敛，本面的失配是温和扰动）。工件 `26p-c-margins.json`。

## 4. 反向判据（三读数＋还原核）

- **M1 摘组合门（dynamicOff＝回 loose 轴 -45）⇒ (a) 塌（回声自触发复现）**：同一失配回声材料，
  loose 链在残余 -36.74dBFS 上**自触发**（announcement aborted，无说话者注入）——11.5/10min
  族在集成面复现；**还原核**：dynamic 回 ON 同材料 triggers=0 held ✓。
- **M2 摘 AEC3（aecBypass＝判据吃原始采捕）⇒ (b) 逐字转写复现**：回声 raw 进判决面 ⇒ episodes
  打开（判据看见了公告音频）＋同一材料过**真 ASR 链**（silero+流式 Paraformer）逐字转写
  **「我已经完成明天下午三点的会议改到四点请准时参加谢谢」**（59-r1 的 5/5 形状在本面复现）；
  **还原核**：AEC 回来同材料 episodes=0、triggers=0（(b) 面复原，成对闭合）。

## 5. 零回归

- **六族 staged 合跑**（no-tag staged 二进制＋A1326 staging 全套：三 DLL 同目录＋WISP_* env＋
  cwd=包目录）：**34 PASS / 1 FAIL / 0 SKIP**；唯一红 `Test384ResidentWatchdogEndToEnd`＝合跑
  负载时窗——**单跑 PASS（1.15s）**，照 384-r1/390-r1（A1329）既有归因：HEAD 既有形状、红名
  集合零新增。真设备案全过：`Test381r2RealDevicePlayingPhaseAbortWithin400ms`（34.57s）、
  `Test390r1BallClickRealDeviceWithin400ms`（33.93s，CLICK-27 在内）。
- **`262/640/902` 逐位同**：`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`（15e 端到端）。
- `go test ./internal/audio`（无 tag 全包＋`-tags webrtcap` 全包两形）ok；aecspike 跨件对拍 PASS。
- 两尺：gofumpt `-l`（两树全部触件）＝0；vet（无 tag／`-tags webrtcap` 两形状）rc=0；build 两
  形状 rc=0。
- d22scan：**clean**（分母 bans#1-5 internal/=287、cmd/=52、ban#8 internal/=729、cmd/=180——
  增量＝本腿新件）；`-self-test` **37/37 方向检查过（20 ring/17 silent）**。
- ⛔ 越界核：internal/{risk,session,speech} 零触碰（monitor 对 speech 只用公共面：NewDetector/
  NewStreamingAsr/NewTranscriber/MatchaTts）；D43 表/statemachine/thresholds/golden 零触碰；
  deps.toml/AEC3 源零触碰；testdata 仓根零建；⛔ push。

## 6. 写面申报

`cmd/wisp/resident_voice_bargein_windows.go`（新）·`cmd/wisp/resident_voice_bargein_26p_windows_test.go`（新）·
`cmd/wisp/resident_voice_bargein_26p_echo_windows_test.go`（新）·`cmd/wisp/resident_voice_windows.go`（M：
bargein 字段＋构造＋scope/conv 接线）·`cmd/wisp/resident_voice_engines_windows.go`（M：arm/far
tap/release 三插点）·`internal/audio/gate.go`（M：SetConversation）·`internal/audio/gate_test.go`（M：
会话窗单测）·`internal/audio/align.go`＋`align_test.go`（新）·`internal/audio/aecspike/align_cross_26p_test.go`（新）·
本件。仓外凭据 `D:\tmp\wisp\20261031\26pr1\**`（bin/out/logs/cases，只建不删）。

## 7. 残余与边界（如实记账）

1. **AEC 双讲损伤的摆动**（L1 两发 -29.4/-38.2）：引擎默认参数未调（绑定按 59 系 pin 的默认
   形状）；28 全双工集成若把公告窗常开为产品姿势，损伤调参是它的活——本面只出读数。
2. **-20dBFS 是 (a) 的说话者电平**（近麦正常说话）；-26.70 控制的过链裕量在 L1 如实报
   （可正可负）。若把「收紧档必须吃 -26.7」当验收，须先调引擎，不是本面单方面能钉的。
3. **monitor 冷启动预热**（对齐 ~0.6s）由 +1.2s 说话者避开；跨公告保温是 28 的集成深度。
4. 公告窗工作点 -30 是**monitor 构造面的选择**（59-v3 §3 验证点）；389 门本体零改——若把它
   反写进门常量属 389 面。
