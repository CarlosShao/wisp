# 票 397 落地腿 `397-r2` —— monitor 复活三处窗间清一体落地（写腿）

日期 2026-11-01 13:0x-14:0x +08 · 腿 `397-r2`（写、`cmd/wisp/**`＋`internal/audio/**` 如需＋本件；
编排者重授 cmd/wisp）· 起手 HEAD `d8cc6f785`（dev，开工三核过：①锚点等号在 ②索引脏区全
`.scratch/**` 他人临时件、写面零脏 ③进程表核无 rig 在飞）· 派单存档
`.scratch/wisp/probes/orch/20261101/txt/397-r2-dispatch.txt` 在案 · 397-r1 核清四事实（三清公开面
在、爆炸半径三处、树零漂移）免重跑，直接采纳。

## 0. 一句话

**三处窗间清一体落地进 monitor（`release()` keepalive 分支），staged 全链保温复活读数全愈：
第二暖窗 episodes=1、triggers=1@talker、engine_warm=true、判官全新实测（ready=9728），warm
剖面与冷 (a) 轮同带；反向四发两红两健康（no-Finalize⇒判死形、no-cooldown⇒hangover 瞬态形
定向复现；no-Reset/no-Flush 在本材料读健康＝跨面差异如实记）；`FINDING MOVED` 断言改钉
`WARM REVIVAL MOVED`（episodes=1/triggers=1@talker）＋工件壳修复（全导出、剖面在盘）；零回归
（26p 4/5＋(c) HEAD 对照树同形红、调参复跑绿、六族 34/1/0 单跑核绿、262/640/902 逐位同、
四数 98/75/0/11 逐位同、两尺/d22scan 全绿）；keepalive=false 默认路径逐字节不动；⛔ 三停手
条款全未触发。**

## 1. AC#1 三处窗间清一体落地（一腿一 commit、⛔ 拆开）

`cmd/wisp/resident_voice_bargein_windows.go`（唯一产品码面）：

1. **`release()` keepalive 分支**（替换原 Close/nil 形）：**APM 持有＋`Reset()`**（393 公开
   面；reset==rebuild 由 392-r2 矩阵与 393-v1 读数承着；`Reset` 拒绝时 `Close()+nil` 兜底＝
   下一 arm 冷重建，warn 在案）；**resampler `Flush()`**；**detector `Finalize(nil)`**（中断
   留下的开口话态在此关闭、边界事件弃置）**＋1.2s 静音冷段**（`voiceBargeinCoolWindows`=38
   整窗覆盖 ≥19200 样本、逐窗 `Accept(voiceBargeinSilence, nil)`、产生的边界事件一律弃置；
   整窗喂入使 detector 的 partial 缓冲在段末为空）。冷姿势分支 `closeEngines()` 原样。
2. **`arm()` :381-384 重建点消解**：新增 `apmTuning` 记录持有 APM 的构造调参——同调参→持有
   复用（release 已 Reset＝just-created 态）；异调参→冷重建（one-parameter-one-firing 规则
   原样保持，读数永不描述陈旧参数）；冷姿势 holds nothing、每窗照建（pre-392 形不变）。
   **judge 的 measure-then-arm＋~0.6s 盲头零动**（26p-r1 形状保持——消解只在构造层，lag 仍
   每窗全新实测，读数 ready=9728 与冷姿势同值）。
3. **`speakFar` 持有 resampler 的 rate 守卫**（持有复用的正确性前提）：held 的 from-rate ≠
   本公告 engineRate 时放下重建（`Rates()` 既有公开面）。
4. 注释同步（文件头 WHAT THE KEEPALIVE CHANGES、struct keepalive 段、arm/release 段）；
   `resident_voice_windows.go`（dispose 邻区）**零字节**——其 :751 注释在 Reset 形下语义仍真。

**⛔ 冷路径逐字节不动（现核）**：`keepalive: false`（原 :294）HEAD 逐字对照同；`closeEngines`
与 `dispose` 零 hunk（`git diff -U0` hunk 清单在案：改动只落文件头注释/常量块/struct 注释与
apmTuning 字段/arm 构造块/speakFar 守卫/release 注释与分支）。

**零新 API**：Reset（`aec_webrtc_windows.go:101`）/Flush（`resample.go:83`）/Finalize
（`vad.go:167`）全既有公开面；`internal/audio/**`、`internal/speech/**` 本腿 **0 字节**。

## 2. AC#2 保温复活读数（权威面＝cmd/wisp staged 全链；keepalive 显式开启）

**仪器**：A1326 staging 全套（webrtcap 二进制＋三 DLL 贴 binary＋WISP_* 五 env＋cwd=包目录
＋unset GOTMPDIR＋Git Bash runner；每跑前进程表核）。pristine 二进制
`wisp397r2.test.exe`（SHA256 `b77d27dc…eda51`）。发 `Test392KeepalivePostureFindingsAndColdDefault`。

**读数（run1，rc=0，31.30s）**：

| 窗 | 形 | 读数 |
|---|---|---|
| 冷 control | 无注入 | keepalive=false built=true warm=false（生产默认姿势钉住） |
| 冷 round 1/2 | (a) -20 talker | triggers=1、ready=9728、lag 482/513、built=true（逐窗重建） |
| 冷 round 3 | (b) echo-only | triggers=0、episodes=1 全被收紧档吃、res -41..-59（cancelled 深度） |
| 暖 round 4（W1） | (a) | triggers=1、built=true（冷轮刚收、首暖窗全新建） |
| **暖 round 5（W2）＝复活窗** | (a) | **episodes=1、triggers=1@talker、engine_warm=true、engine_built=false、ready=9728、lag=453、rms=-24.13@-30**；剖面 near=[-21.3,-24] res=[-28.9,-22.6]——与冷 (a) 轮同带（[-20.9,-23.6]/[-28.3,-22.9]、[-20.6,-22.6]/[-28.9,-23.2]） |
| 冷 restore（round 6） | (b) | triggers=0、built=true（切回冷姿势重建） |

**对照 392-r2 revival-cool（复刻仪器）**：episodes=1、triggers=1@talker、判官全新实测、残余
cancelled 深度——四元全同带。**392-r1 finding 形（第二暖窗判死）消除。** 断言改钉
`WARM REVIVAL MOVED`（episodes=1/triggers=1@talker）；还原核＝全突变跑完后 pristine 复读
**rc=0**。工件 `392-keepalive.json` 本腿全导出（见 §4 壳修复）。

## 3. AC#3 反向判据（仓外 overlay 四发＋还原核；三失败形定向）

突变件＝`release()` keepalive 分支逐字摘除一件（Compare-Object 逐行核：每件 diff 恰为目标
块）；overlay 构建独立二进制（SHA 互异且均 ≠ pristine）；发同一 pristine 测试。

| 发 | 摘除 | warm2 读数 | 判 |
|---|---|---|---|
| M3 | det `Finalize`＋冷段（det 裸保温） | **episodes=0、triggers=0**，红 | ✓ 定向：392-r1 finding 形回归（判官死）＝392-r2 revival-no-finalize 形 |
| M4 | 仅冷段（留 Finalize） | **episodes=1、triggers=0**，红 | ✓ 定向：hangover 假 start 吃掉 endpointer、talker 漏打断＝392-r2 revival（不冷）瞬态形 |
| M1 | APM `Reset`（裸携） | episodes=1、triggers=1、res=talker 带——**rc=0** | ✗ 失败形未复现（§6.1 跨面差异） |
| M2 | resampler `Flush` | episodes=1、triggers=1、res=[-36.2,-24.3] 健康带——rc=0 | ✗ 无失败形（与 392-r2 相位延续注记同向：卫生面无害） |

还原核：四发后 pristine 复读 rc=0；工作树文件 SHA256 `A07C1929…E59` 前后逐位同（overlay 零触树）。
三失败形中两形在 staged 链定向复现（detector 半的两件）；APM/resampler 半见 §6.1。

**踩坑账（复证在册处方）**：overlay JSON 的键必须**正斜杠**（`D:/work/...`）——反斜杠键被
静默忽略（首建四枚 SHA 全同 pristine 当场抓住；26p-k1 在盘 overlay 同形佐证）。

## 4. AC#4 零回归

- **26p 套件**（pristine webrtcap 二进制）：(a) keyed `Test26PBargeinStopsWithin400ms` PASS
  （BARGE-IN 36ms×2 全带内）；(b) `Test26PEchoNeverTriggersAndNeverTranscribes` PASS；M2
  `Test26PWithoutAECTheEchoTranscribesVerbatim` PASS（逐字转写＋还原核）；等价
  `Test26PFeedOnlineMatchesWindowReplay` PASS（6/6 零漂移）；(c)
  `Test26PMarginRetestAndGateReverse` 三发 1 绿 2 红（M1 反形红）——**纯净 HEAD 对照树
  （`git archive d8cc6f785` 导出、零 397 字节）同枚同 M1 形红**（其 L1 深压 -40.15）＝发动机
  摆动族 HEAD 既有形状、红名集合零新增（A1329 归因法；392-r1 自记「两跑一过一死」同族）。
- **392 调参族** `Test392AecTuningDoubleTalkSwing`：首发红（trigger-threshold-6 round 2 深压
  rms=-63.11 低于可读地板）→ **复跑 rc=0 全绿**（4 设置×4 发＋还原核）＋HEAD 对照 rc=0
  ——摆动一发归因闭合（K=4 设计所测现象本身；冷姿势码路径本腿零行为差）。
- **六族 staged 合跑**（no-tag 二进制＋三 DLL＋WISP_* env＋cwd=包目录，`-test.run
  "Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384"`）：**34 PASS / 1 FAIL /
  0 SKIP**——唯一红 `Test384ResidentWatchdogEndToEnd`（3.30s）＝合跑负载时窗（A1329/
  392-r1/26p-k1 既有形状），**单跑 rc=0**。红名集合零新增。
- **`262/640/902` 逐位同值**：`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`。
- **四数**（无 tag、`PATH=$PWD/third_party/sherpa-onnx`、`./internal/audio/...` `-count=1 -v`）：
  **rc=0，顶层口径 RUN=98 PASS=75 FAIL=0 SKIP=11——与 392-r1:108/392-r2/392-v3 自记逐位同**；
  SKIP 集＝TestLiveWasapiSmoke＋aecspike 十枚硬件/工件门（一枚门名密掩占位）逐名同形。
  （另 `go test ./...` 全仓扫 9 红全在写面外——panel/observe/risk/agent/winsec 的负载/环境
  形状；基线四数口径从不跑 `./...`，且两枚改动文件与那些包零编译域交集。）
- **keepalive=false 默认路径逐字节不动**：见 §1 现核。
- **两尺**：gofumpt `-l`（cmd/wisp＋internal/audio＋internal/speech）＝0；vet no-tag rc=0、
  `-tags webrtcap` rc=0（cmd/wisp 腿标准双形状）；GOOS=linux `internal/audio/...` rc=0
  （392-r2 同尺）。全仓 build（无 tag）rc=0。
- **d22scan**（Git Bash `scripts/d22scan.sh`）：实扫 **clean rc=0**；selftest **PASS=35 FAIL=0
  SKIP=0（RUN=77）**与基线逐位同；分母 bans#1-5 internal/=287、cmd/=52、ban#8 internal/=731
  （=392-r2 入树的 aecspike 件 +1、392-v3 在案——**本腿零新文件、零增量**）、cmd/=183。

**工件壳修复（申报为具名小扩展）**：`voice392Round` 字段未导出 ⇒ 392-r1 工件 round 对象全
`{}`（392-v3 §1 实锤）。本腿加 `voice392RoundOut`（json tag 全导出）＋转换器，两件 392 工件
（keepalive/aec-tuning）读数逐字段落盘；`finding` 键改 `revival`（语义随断言翻动，392-r1 原件
先备份 `D:\tmp\wisp\20261101\397r2\out-backup\`）。扩展落点在 :271-322 断言区外的 instrument
面，为 AC#2「残余深度」读数提供盘上支撑。

## 5. 写面申报（commit pathspec 逐枚）

仓内：
1. `cmd/wisp/resident_voice_bargein_windows.go`（M：三处窗间清＋arm 消解＋speakFar 守卫＋注释）
2. `cmd/wisp/resident_voice_bargein_392_windows_test.go`（M：断言改钉＋工件导出＋文档同步）
3. `docs/evidence/s1/397r2-monitor-revival.md`（本件）

仓外 `D:\tmp\wisp\20261101\397r2\**`（bin＝pristine＋四突变＋no-tag 二进制＋三 DLL；mut＝四
突变件＋四 overlay＋SHA；logs 全量；out-backup＝392-r1 工件备份；run.sh/build.sh；head-control
＝纯净对照树）只建不删。测试工件落 `D:\tmp\wisp\20261031\392r1\out\`（测试硬码根、原件先备份）。

⛔ 越界核：`resident_voice_windows.go` 零触碰；`internal/{audio,speech,risk,session}` 零触碰；
AEC3 源树零触碰；389 门本体（`internal/audio/bargein.go`）/thresholds/golden/D43 表零触碰；
deps.toml/下载/testdata 仓根/frontend/design/scripts/.github/票面/台账/specs/PLAN 零触碰；
⛔ push（推送归编排者）。

## 6. 残余与边界（如实记账）

1. **no-Reset / no-Flush 失败形未在 staged 链复现（跨面差异，具名）**：392-r2 复刻矩阵的
   revival-no-reset 臂（W2 triggers=0、残余 -42.78≈近端）在本 staged 材料/窗序上不复现——裸携
   APM 的 W2 照判（每帧 delay hint 重锚＋窗间声学几何同一，携带滤波器仍抵消）。392-v3 §1 已
   实锤 392-r1 对 APM 半的 staged 原读数是空壳无盘上支撑；本腿是 APM 半在 staged 链的第一份
   全导出读数（健康）。**实现保持 Reset 形**（392-v3 采纳路径①的统一器兼省构造裁定不依赖本
   发；复刻矩阵 19.24dB 承重读数仍为最敏感证据）；更广材料/几何上的 Reset 承重读数归默认翻
   转腿的全链复读。Flush 同理（复刻自身的 iso-det-fresh/相位延续注记本就判无害）。
2. **(c) M1 与调参深压＝发动机态摆动族**：两侧归因在案（HEAD 对照树同形红／复跑绿）；
   (c) 的摆动脆弱 26p-v2/392-r1 已在册，本腿零处置（验收腿/编排者裁量）。
3. **冷段是 owner 侧纪律**：silero session 无 RNN reset 面（`internal/speech` 无此 API），
   392-v3 §6 第 2 步的 internal/speech 票原样待开；「衰减多少够」仍是 1.2s 一点读数。
4. **staged 复活读数的材料边界**：单发＋还原核（同发二跑未做——26p 家族的摆动先例下，
   复活断言的多发稳定性读数属默认翻转腿的全链复读范围）。
5. 停止条款核：①arm() 消解与 26p-r1 形状冲突——**未触发**（judge 盲头原样，消解只在构造层）
   ②三清接入后 keepalive 显式开启读数仍判死——**未触发**（复活读数全愈）③必须动 keepalive
   默认——**未触发**（:294 逐字节不动、0 字节）。
