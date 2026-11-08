# 票 291 写腿 `291-mic` 读数台账 — AC#2 四形真麦 ＋ AC#3 落点表 ＋ AC#4 门禁

钟点：2026-11-08 16:49–17:02 +08（`date` 现跑）。树锚：主仓 `8b568eccd`（A1639）。
写面：`.scratch/wisp/probes/291/**` ＋ `docs/evidence/s1/291-level-dc.md` ＋ 票 291 Progress log 一行。**零产码**。
工作方式：`git clone --shared` 导出树 `/d/tmp/wisp/20261108/291-mic-tree`（临时件只建不删），全部采集与门禁都在导出树跑。

## 起手三核（现跑，⛔ 信转述）

1. **票面引文核**：`internal/audio/level.go` 求和段实际在 **`:96-107`**（`LevelOfSamples`，int64 平方累加在 `:100-104`，无去均值步）；「DC counts: RMS does not remove a DC offset」注释在 **`:37-39`**。票面写「`:107-118` 一带」偏 +2～+8 行，代码本体在、语义同。`FrameLevel` 实际 `:118-128`。
2. **引文出处核**：票面引的 `capturelevel_windows_test.go:21-24`（样值 12345）属实；`internal/ball/liquid.go:30` `SilenceLevelGate = 0.06` 属实；`docs/evidence/s1/241-audio-level-producer-v1.md:111` 一带「不除 DC…不是缺陷」属实。
3. **飞件核**：两把 tasklist 尺 16:5x 现读 0（`balldebug.exe`、`wisp.exe` 均无）；采集在共享克隆树跑，与主仓三腿互斥面（`frontend/**`／`internal/ball/**`）零接触。

## 台件真名与真麦状态

- `cmd/wisp/resident_audio_247_live_windows_test.go`（`WISP_LIVE_MIC=1` 门，三窗 harness：`bootAudioRuntime`/`writeAudioConfig`/`assembleCapture`/`requireNoCompetingAudioProcesses`/`livePhase`/`phaseRouter`）。
- `internal/audio/hotplug_test.go:524` `TestLiveWasapiSmoke`（同门）。
- 设备：`麦克风 (Realtek(R) Audio)` Status OK（`Get-PnpDevice -Class AudioEndpoint` 16:5x）。
- 冒烟：`TestLiveWasapiSmoke` **PASS**，25 帧/834ms，`capture="麦克风 (Realtek(R) Audio)" render="扬声器 (Realtek(R) Audio)"`，0 丢帧（16:55）。

## AC#2 四形读数（16:57:34–16:58:13，probe 台件 `resident_audio_291_probe_windows_test.go`，拷贝件在 `.scratch/wisp/probes/291/mic/`；probe 源只存在于克隆树，从未提交）

probe 四窗结构：安静 3s → 440Hz 已知幅度音调 6s（源样值幅度 12000，**源理论电平 12000/√2/32768 = 0.258950**）→ 说话窗 25s → 对照安静 3s；另做形④合成帧直过 `audio.LevelOfSamples`。

| 形 | 窗口 | 采到帧 | mean | max | levels/s |
|---|---|---|---|---|---|
| ① 完全安静（3s） | shape1-quiet | 93 | **0.003067** | **0.099908** | 31.00 |
| ③ 已知幅度音调外放（6s） | shape3-tone | 190 | **0.001321** | **0.002530** | 31.67 |
| 说话窗（25s，**无人说话**，⛔ 不是形②） | shape2-speech | 781 | 0.001490 | 0.012336 | 31.24 |
| 对照安静（3s） | control-quiet-after | 94 | **0.001305** | **0.001696** | 31.33 |
| ④ 合成 DC(12345)＋整周期正弦(6000) 帧 | LevelOfSamples 直采 | — | **0.398368** | — | — |
| ④' 合成纯 DC 12345 帧 | LevelOfSamples 直采 | — | **0.376740** | — | — |

- 形④闭式校验：`sqrt(12345² + 6000²/2)/32768 = 0.398367`，读数 0.398368，delta **4.369e-7**（int16 舍入量级）。
- 形④'：`12345/32768 = 0.376740`，与票面夹具值逐位一致。
- 边缘计数：`levels=1158 frames_sent=6 frames_dropped=1152`——有界信道满（本 harness 无识别器在抽另一条信道，247 live 台件同形，不属本票缺陷）；**levels 信道节奏 31/s 满速，读数路径完整**。

## 三个必须如实报的事实（按票面「如实」条款）

1. **形②采不到**：子代理身份无法触发人工交互（`ask_user` 被拒：「human interaction is unavailable while the calling agent is owned by another live agent」）。真人说话窗**没有采**，上表 25s 窗是无人窗，⛔ 不冒充形②。机主在场复采命令见文末。
2. **今日安静底噪 mean 0.003067，不是票面 2026-10-09 快照的 0.368462**：DC 底噪在今日这台机的设备状态里**不存在**（增益/驱动状态已变）。票面核心前提「真麦安静读 0.37」**今日不可复现**；`0.376740` 今天只在合成 DC 帧上成立（形④'）。「设备增益」与「尺不减 DC」两件事照票面禁令分开写：后者是定义事实（形④钉死），前者今日读数已翻面。
3. **形③分开不了（具名）**：音调源电平 0.258950，麦读 mean 0.001321 / max 0.002530——**低于安静窗①的 max 0.099908**，与对照窗 0.001305 分不开 ⇒ 扬声器→麦通路今日无效（扬声器音量/静音，或麦增益过低）。按票面判据「分开不了就具名说分开不了，⛔ 不许挑一个能过的门槛凑绿」。

## AC#2 判据的回答

「说话」与「安静」之间是否存在**既有的、不动用的判据**能分开——**形②缺席，判不了**。两条侧证照实给：
- 今日安静窗 mean 0.003067 < `SilenceLevelGate 0.06` ⇒ 球侧那道门**今日会正确关上**（10-09 的「永远开着」形状今日不复现）。
- 「尺不剥 DC」这条定义事实由形④钉死：读数精确含 DC 分量（0.398368 = 闭式），与设备状态无关。

## AC#3 落点表（⛔ 本票未落码，只出表）

| # | 文件:行（今日现读） | 最小改法 | 为什么 |
|---|---|---|---|
| 1 | `internal/audio/level.go:96-107`（求和段 `:100-104`） | 求和前先 int64 累加样本均值；第二遍累加 `(s-mean)²`；返回 `sqrt(方差)/LevelFullScale` | 去掉 DC 的唯一最小落点；两遍遍历保住「无分配、无状态、纯函数」的注释契约 |
| 2 | `internal/audio/level.go:17-39`（尺的注释块） | 重写端点段 `:22-26` 与 `DC counts` 段 `:37-39` | 纯 -32768 帧将从 1.0 变 **0.0**，端点语义搬家；注释是尺的定义面，必须与算术同步 |
| 3 | `internal/audio/level_test.go:72-87` `TestLevelFullScaleSquareEndpoints` | 改期望与注释（`all -32768 frame` 从「exactly MaxLevel」改「exactly MinLevel」；对称方波端点 32767/32768 不受影响保留） | 逐字钉着「all-minus-32768 → 恰好 1.0」（LevelOfSamples 与 FrameLevel 两处断言）⇒ **去 DC 后必红的头号钉**。改期望不是放宽断言，是尺的定义变了 |
| 4 | `internal/audio/level_test.go` 其余 9 枚 | 预判不红（输入全是零均值正弦或纯零帧，去 DC 对它们是恒等变换） | 名册：`TestLevelSilentFrameIsExactZero`/`TestLevelSineAgainstRootTwo`/`TestLevelIsLinearInAmplitude`/`TestLevelSineToleranceIsNamedAndBounded`/`TestLevelSameInputSameBits`/`TestFrameLevelRejectsNonSeamFrames`/`TestSeamCodecRoundTrip`/`TestLevelOverBoundedChannelFromWavInjector`/`TestLevelFrameIsOneSeamFrame`。⚠「改前后 `go test ./internal/audio/ -count=1 -v` 逐名红名册作差」按编排者 2026-10-09 14:2x 裁定欠在落地腿，本只读腿欠不了 |
| 5 | `internal/audio/capturelevel_windows_test.go:21-24` | 不改（期望值派生式 `:50 want := LevelOfSamples(...)`，全仓零枚 0.3767 字面 ⇒ 不会红） | 真代价＝夹具失去见证力：它见证的正是「含直流」这件事，去 DC 后恒定 DC 帧读 0，夹具空转 |
| 6 | `internal/ball/liquid.go:30` `SilenceLevelGate = 0.06` | **⛔ 不许动**（C21 冻结 token，机器钉 `internal/ball/tokens_table_test.go:552`） | 去 DC 后该 token 的含义搬家（今日读数 0.003 < 0.06，门会关），但数字一字不许动——含义搬家要在落地票里向 owner 摆明 |
| 7 | `cmd/wisp/resident_audio_247_live_windows_test.go:194`（2× 门槛） | **本表只登记不改**（归票 247） | 去 DC 后 loud/quiet 分离在真机上有望达成，但改那格的判据须连凭据交非实现者裁 |

## AC#4 门禁读数（全部现跑）

- 两把 tasklist 尺：`balldebug.exe`/`wisp.exe` 均 `INFO: No tasks are running`（16:5x）。
- 克隆树 `GOFLAGS= go build ./...`：**rc0**（17:0x；先把随依赖拷入的飞腿杂散件 `build/mut_rvw.go`、`build/probe61r2/main.go` 挪出树外后跑， genuine rc）。
- 克隆树 `sh scripts/d22scan.sh`：**clean**，rc0（`internal/=964, cmd/=242, frontend/=121, design/=14` 全查）。
- `/d/work/base/gopath/bin/gofumpt.exe -l` 对 probe 台件：**空**（rc0）。
- 主仓产码零改动：`git status --porcelain | grep -v '^??'` 读数见证据件（应只剩我写面的 3 处＋其他腿既有未收件）。
- `frontend/**`／`design/**` 零读零写；`PLAN.md`/specs/thresholds/golden 零碰；零 push。

## 机主在场复采形②的可核命令（编排者转交机主即可）

```bash
cd /d/tmp/wisp/20261108/291-mic-tree
unset GOTMPDIR TMPDIR TEMP TMP
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
export WISP_LIVE_MIC=1
/d/tmp/wisp/20261108/291probe.test -test.run TestAC291FourShapeLiveMicReadings -test.v -test.timeout 120s
```

流程：启动后约 10 秒音箱会嘟 440Hz（形③）；**听到嘟声立刻开始正常音量对麦连续说 25 秒**，说完停。全程约 50 秒。读数在 `291 reading:` 四行。⚠ 今日形③音调麦听不到——复采前建议先在系统音量里确认扬声器非静音、麦克风输入增益非 0，否则形③照样分开不了。
