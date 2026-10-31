# 票 389 落地腿 389-r1 — Path-C barge-in 组合门＋三件套动态门控产品形态

日期 2026-10-31 · 腿 `389-r1`（写腿＝实现者；写面 `internal/audio/**` ＋ 本件）·
起手锚 `08ffbd72d`（dev；开工三核过：①HEAD 在 dev ＝ `08ffbd72d` ②索引干净 ③写面
`internal/audio` 零脏、`docs/evidence/s1` 仅他人未跟踪 381/382 判语不碰）。
接续链：`59-v3-verdict.md` §0 PRECHECK P15 无条件 PASS＋§6「26-PathC 立即重开进实现」
⇒ 本腿＝重开首项。派单存档 `389-r1-dispatch.txt` 在案。

## 0. 一句话结论

**组合门产品形态已落 `internal/audio` 并在 59-w 盘上工件上出齐三读数**：产品门
（`BargeinGate`：VAD∧能量∧ASR 静态三档＋三件套动态门控）在 −45 宽松轴上把塌峰态
误触发从 **11.5/10min 压到 0.00/10min**、好峰态从 5.4 压到 0（两态同时归零，好峰水平
或更低 ✓）；静态 −30 轴两态保持 0（零回归 ✓）；正控 −26.70dBFS 真语音两轴全过、
−36.58dBFS 安静语音在 −45 轴**从静态 −30 门的「被吃」变「可达」**。反向判据双闭环：
产品开关关 ⇒ 塌峰回到 11.5 基线；仓外 `-overlay` 突变（切断收紧接线）⇒ 动态开读数与
关逐格相同 ⇒ 断言抓到 no-op 红 ⇒ 摘 overlay 还原核回到 0/0 绿。四数 94/83/0/11
（对 59-w 基线 83/73/0/10：+11/+10 全为本腿 11 枚新增、removed=0、FAIL=0）、两尺
clean、d22scan 实扫 clean＋正控 PASS=35。

## 1. 环境与复用（⛔ 重造边界全部遵守）

- 引擎/构建产物：`D:\tmp\wisp\20261031\aec3r1\build\libwebrtc_audio_processing.a`
  （aec3-r1 钉定成品；零重建、零源树改动、零构建参数改动）；`deps.toml` 未动一字。
- **材料＝59-w 盘上工件离线复用，零 rig、零实时采捕、零 ASR 重跑**（派单材料复用
  边界；59-v3 §7「rig 已饱和、不再需要独立测量腿」逐字遵守）：episode 全集（含
  ASR 转写与 RMS）来自 `59wr1/out/gates-probe.json`（55 材料）＋`gates.json`（正控
  2 材料）；窗级判据值来自各 run 的 `-curve.json`（逐窗 ERLE）＋`autocorr.json`
  （TTS 臂免参考 selfCorr，>0 才有效——探针在早期窗与 inactive 窗写 0）＋
  `probe.json`（music 臂 corr_vs_far，真 far 逐位可重建；TTS 臂的 corr_vs_far 不可
  用，corr[0]=0.126，与 59-w §5.6 判别器选择一致）。事件材料窗段取 `probe.json`
  的 events（**与门材料切分同源**；autocorr 自己的事件列表分段规则不同，只可取
  逐窗值、不可取其事件表——本腿踩过这个错位并以 probe 段修正，见 §7 处方注）。
- 音频错峰：本腿全程零音频设备使用（纯离线 JSON 读数＋纯数学单测），无进程表
  冲突面；四数跑法照 59-w §6 文档化形状（PATH 前缀 third_party/sherpa-onnx，裸跑
  0xc0000135 为 59-r1 时代既有形状）。

## 2. 产品形态（AC#1：组合门＋动态门控落 internal/audio）

交付件 `internal/audio/bargein.go`（无 build tag、纯 DSP、纯样本计数、无 goroutine）：

- **静态三档**（59-w rig 同形）：episode 级判决＝能量门 ∧ ASR 文本门（非空转写；
  VAD 起点由调用方供给——本包不分割音频、不依赖 internal/speech，`BargeinEpisode`
  是数据接缝）。
- **三件套判据产品化**（59wr1 §5 判据形状、常量同值零发明）：
  ①`CollapseDetector` 滑窗 ERLE 骤降报警（中位 −6dB／绝对 15dB／持续 3 窗、一次每
  streak、inactive 窗跳过不破 streak）＋在持查询 `Sustained()`；②③`FarCorrelationDetector`
  对 far 参考的滑窗相关（floor 0.8，music 地面真值标定自 59-w；**②③在产品形态
  合流为一个量**——产品 far 真实可得，rig 时代 TTS 臂的免参考自相关只是 far 不可
  重建时的替代仪器，59wr1 §5.5 自己写明「对产品形态换算＝对远端参考的滑窗相关」，
  窗级相关值同时就是「对齐强度」，60b 头部 0.65 整发可见的判据由同一量承担）。
- **动态门控语义**（59-v3 §3 两常量的产品化正解）：判据任一在持 ⇒ 窗标 untrusted ⇒
  该窗内 episode 的能量门**收紧到 −30dBFS**（59-v3 §3 的常量——「静态 −30 不是
  产品门」的产品化读法＝−30 从静态门变成动态收紧档）；报警解除（streak 断）⇒ 回
  标称 −45（放行）。收紧查询在 episode **起点窗**（barge-in 决策时刻）。
- **double-talk 不误伤**（59-w §3「这正是产品必须存活的行为」）：真人插话压低 ERLE
  触发报警 ⇒ 收紧到 −30，实测真人带 −26.70dBFS（control-26）仍过；被收紧吃掉的
  塌峰伪迹带实测 ≤−32.1dBFS——收紧只是在报警窗内把「贴阈伪迹带」与「真人带」
  的裕量交换方向，不是抑制开关。
- 两条输入接缝同达一个决策路径：`Feed(far, near, residual)`（原始帧，内部做 lag 环
  形移位，照 59-w longAecPass 的 echo-source-time 形状）与 `ObserveWindow(w, corr)`
  （窗级——引擎自报监督数或离线重放用），`Evaluate(episode)` 唯一判决。
- `BargeinOptions` 全字段显式（`NewBargeinGate` 零值钉实测常量）；`DynamicGating`
  开关是产品形态的一部分（关＝纯静态三档），不是测试后门。

## 3. 两态读数（AC#2：对照 59-w 的 5.4/11.5 基线）

读数驱动＝`TestSpike389r1BargeinReadings`（`internal/audio/aecspike/
bargein_readings_windows_test.go`）：每材料的窗级原值经 `ObserveWindow` 折进产品门
（**判据在产品内行使**——ERLE dB 反推窗 stat 进 CollapseDetector、相关原值进 floor
判定，无预 OR 短路），episode 经 `Evaluate` 真实决策路径。材料分母逐字对齐 59-w §4
表：good 9 份 222.0s、collapse 44 份 104.4s（43 事件窗＋59v2 存档；state 字面
"collapsed" 的 cold-tts-2 按表算术不在分母——59-v3 §1 复算同值，本腿如实记录这个
表格字面与算术的差并以数字为准）。

| 轴 | 动态门控 | good /10min | collapse /10min |
|---|---|---|---|
| −45（宽松档） | **关**（静态链＝59-w 形） | 2 事件＝**5.41** | 2 事件＝**11.50**（基线逐格复现） |
| −45（宽松档） | **开**（产品门） | **0.00** | **0.00** |
| −30（收紧档） | 关 | 0 | 0 |
| −30（收紧档） | 开 | 0 | 0 |

- **AC#2 达成**：动态门控把塌峰态 11.50 → 0.00，等于好峰水平（0 ≤ 0）或更低——
  两态同时归零。逐枚归因：塌峰 2 枚误触发（−32.24「听我」@120s 发 event10、−35.59
  「我们」@event5）在事件窗内均被判 untrusted（窗级 selfCorr 0.339–0.66 <0.8）⇒
  收紧 −30 吃掉；好峰 2 枚（−38.85「你能」、−32.32「哎我」@steady）所在窗 trusted
  与否不改变结果（−45 轴上它们原本就是仅有的两枚，动态开后归零）。
- **正控**（真语音电平，`gates.json` 实测 episode RMS）：control-26（−26.70dBFS，
  转写「天气怎么样帮我看一下日历」）两轴全过——门不吃真人；control-36（−36.58dBFS）
  在 −45 轴**过**（静态 −30 门会吃的安静语音被宽松档＋动态门控救回，59-v3 §3 的
  「−36 安静语音被吃」代价在产品形态下消失）、在 −30 轴被收紧档吃（报警窗内安静
  语音的代价仍按 59-v3 记账，是收紧档的接受面）。
- **−30 轴零回归**：两态在动态开/关下都保持 0（对 59-w 静态 −30 读数无劣化）。

## 4. 反向判据（AC#3：摘动态门控 ⇒ 塌峰回升 ⇒ 红；三读数＋还原核）

两形状闭环，都在 §3 同一材料/同一决策路径上：

1. **产品开关反向**（`DynamicGating=false`）：−45 轴塌峰回到 **2 事件＝11.50/10min**
   （59-w 静态基线复现）、good 回 5.41——摘门控 ⇒ 回升，且读数测试以关系断言钉住
   （collapse(on) < collapse(off)、collapse(on) ≤ good(on)、good(on) ≤ good(off)、
   −30 轴全 0），无 golden 绝对值。
2. **仓外突变核**（`-overlay`，`D:\tmp\wisp\20261031\389r1\mutant\bargein.go`）：把
   `Evaluate` 的收紧接线 `g.wins[w]` 突变为 `false`（三件套照跑、结果不进决策——
   区别于开关关的「静态档」语义）⇒ 动态开读数与关**逐格相同**（good 2/collapse 2）
   ⇒ 断言「dynamic gate is a no-op」红（rc=1）。证明读数差异的全部来源就是被突变的
   那段决策接线。
3. **还原核**：摘 overlay 重跑 ⇒ 四格回 0/0/5.41/11.50 原读数、PASS（rc=0）；盘上
   `bargein-readings.json` 被还原核重写为还原后读数。

## 5. 零回归与门面（AC#4/AC#5）

- **四数**（`./internal/audio/...` 全集，文档化跑法＝PATH 前缀 third_party/sherpa-onnx；
  裸跑 0xc0000135 为既有形状；日志 `D:\tmp\wisp\20261031\389r1\four-numbers.log`）：
  `-count=1 -v` ⇒ **RUN=94 PASS=83 FAIL=0 SKIP=11**（rc=0）。对 59-w 基线 83/73/0/10
  逐名作差：**added 11 枚全为本腿**（`bargein_test.go` 10 枚：TestBargeinWindowStat
  ERLEdB／TestSlidingERLEFoldsWindows／TestCollapseDetectorShapes／TestCollapseDetector
  Sustained／TestFarCorrelationDetector／TestBargeinGateStaticTiers／TestBargeinGate
  DynamicTighten／TestBargeinGateDynamicDisabled／TestBargeinGateLagShift／TestBargein
  GateWindowSeam，全 PASS；`bargein_readings_windows_test.go` 1 枚 TestSpike389r1
  BargeinReadings，无 WISP_AEC_SPIKE 时按设计 SKIP）；**removed=0**；59-w 的 10 枚
  SKIP 集合零丢失；基线 83 枚枚枚同名同态零变红。
- **两尺**：gofumpt（`D:\work\base\gopath\bin\gofumpt.exe`）`-l internal/audio/`＝空；
  `go vet` 三形状 rc=0（无 tag ＋ `-tags webrtcap` ＋ `GOOS=linux`）。
- **build 两形状** rc=0（无 tag ＋ webrtcap 的 `./internal/audio/...`）。
- **d22scan**（Git Bash `scripts/d22scan.sh`，日志 `D:\tmp\wisp\20261031\389r1\d22scan.log`）：
  实扫 **clean** rc=0；正控 selftest **PASS=35 FAIL=0 SKIP=0（RUN=77）**。导出树
  bans#1-5 internal/=**286**、cmd/=**51**、ban#6 frontend/=116、ban#7 internal/tools/=25、
  ban#8 design/=39、frontend/=116、internal/=**726**、cmd/=**176**。对 59-w 的
  285/49/…/723/174：**internal +1＝本腿 `bargein.go`、ban#8 internal +3＝本腿三件**
  （bargein.go＋bargein_test.go＋bargein_readings_windows_test.go）；**cmd +1 与本腿
  无关**（本腿零 cmd 写面，差属在飞 12-r2——59-w §6 对同类差的归因先例）。
- 路径长度门：本腿最长相对路径
  `internal/audio/aecspike/bargein_readings_windows_test.go`（52 字符＜100 帽）。

## 6. 写面申报（逐枚）

仓内（commit pathspec 逐名列出）：

1. `internal/audio/bargein.go` —— 产品模块：三件套判据（SlidingERLE／CollapseDetector／
   FarCorrelationDetector）＋组合门（BargeinEpisode／BargeinOptions／BargeinGate：
   Feed 帧接缝＋ObserveWindow 窗级接缝＋Evaluate 判决＋动态收紧）＋实测常量组。
   无 build tag（CI 全平台可算）。
2. `internal/audio/bargein_test.go` —— 上件的 CI 单测 10 枚（判据形状、chunk 不变性、
   收紧语义、开关反向、正控常量 −26.70/−36.58、double-talk 不误伤、lag 环、窗级接缝）。
3. `internal/audio/aecspike/bargein_readings_windows_test.go` —— 读数驱动 1 枚
   （windows 门控测试）：59-w 盘上工件离线重放 ⇒ 四格矩阵＋正控＋关系断言（先落盘
   后断言，失败时读数在盘上）。
4. `docs/evidence/s1/389r1-bargein-gate.md` —— 本件。

仓外（`D:\tmp\wisp\20261031\389r1\**`，只建不删）：`out\bargein-readings.json`
（四格矩阵＋逐枚 tier-3 事件＋正控＋逐 run untrusted 窗计数）、`four-numbers.log`、
`d22scan.log`、`mutant\bargein.go`＋`overlay.json`（突变核工件）。零音频 artifact
（本腿无任何设备采捕/渲染）。

⛔ 未碰：`cmd/wisp/**`（12-r2 在飞）、`internal/{risk,session,speech}/**` 语义
（零 import、零调用——episode 数据接缝把 VAD/ASR 留在调用方）、AEC3 源树/构建
参数、deps.toml、59-w 的 longrun 4 件（`longrun*.go` 零改动）、票面/台账/specs/
PLAN/testdata、thresholds/golden。

## 7. 复跑处方（逐字）

读数（免 rig、免模型；须有 59-w 工件在盘，缺则 SKIP）：

```
cd <repo>
set WISP_AEC_SPIKE=1
set PATH=<repo>\third_party\sherpa-onnx;%PATH%
go test ./internal/audio/aecspike/ -run TestSpike389r1BargeinReadings -count=1 -v
```

（读数落盘 `D:\tmp\wisp\20261031\389r1\out\bargein-readings.json`。材料窗轴取
`probe.json` 的 events——autocorr.json 的事件表分段规则不同，只可取其逐窗
selfCorr 值；窗有效性＝selfCorr>0，探针在早期窗/inactive 窗写 0。）

突变核（收紧接线切断 ⇒ 预期红）：

```
go test -overlay D:\tmp\wisp\20261031\389r1\overlay.json ./internal/audio/aecspike/ -run TestSpike389r1BargeinReadings -count=1
```

（还原核＝同命令摘 `-overlay` ⇒ 预期 PASS。）

四数／两尺／d22scan：§5 逐字（四数日志与 d22scan 日志在临时根）。

## 8. AC#6 集成边界（59-v3 判语逐字遵守）

**400ms 停播链／产品麦面零转写／−30 裕量重测加严＝票 26 集成面出真数**——本腿零
rig、零测量腿（59-v3 §7「rig 可答的问题已饱和」）；本腿把 −30 钉为动态收紧档常量、
把裕量语义写进门判决读数面（`ThresholdDbfs` 逐 episode 在案），集成面复验时直接读。
三件套判据产品化**零 AEC3 源改动**（未定义即停②不触发）；零 cmd/wisp 触碰
（③不触发）；动态门控在塌峰态**拉平达成**（①不触发——无差值可报）。

## 9. 设计输入对照（两头钉死 → 门内落点）

- 59wr1 §5 三件套判据 → §2 判据①②③（常量 6/15/3 与 0.8 同值零发明）。
- 59-v3 §3 两常量 → 标称 −45（宽松档：−36 可达）＋收紧档 −30（报警窗内）；
  「静态 −30 非产品门」的产品化＝−30 从静态门变动态收紧档。
- 59-w §4 基线 5.4/11.5 → §3 表逐格复现（5.41/11.50）后压到 0/0。
- 59-w §3「产品必须存活 double-talk」→ §2 收紧语义＋control-26 两轴全过实测。
