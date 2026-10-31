# 票 389 — Path-C barge-in 组合门＋动态门控 终裁判语（腿 389-v1，非实现者）

日期 2026-10-31 · 腿 `389-v1`（终裁腿；写面仅本件）· 收件 `651067412`（389-r1：4 枚 1743 行）·
起手 HEAD `b109a3010`（编排者 A1327 收件提交，实现提交＝HEAD~1；`internal/audio` 收件件与
HEAD diff=0、work 树零脏——本腿读码即读收件件内容）。接续链：59-v3 §0 P15 无条件 PASS＋§6
「26-PathC 立即重开」→ 389-r1 落地 → 本腿终裁。派单存档 `389-v1-dispatch.txt` 在案。
实验全部在导出树 `D:\tmp\wisp\20261031\389v1\export`（git archive `651067412`＋junction
third_party）；开工现核在飞三腿写面 `cmd/wisp`／`internal/speech`／`tools/mockllm` tracked
零脏——零触碰、零归因。

## 0. 判语：AC#1–AC#5 全部翻绿；**可带残余结案**——残余恰是票面自己划的 AC#6 集成边界

一句话：产品形态读码核四形状全过；两态读数按 389r1 §7 处方**逐字复跑 rc=0**，四格矩阵＋正控
逐格复现、产出 JSON 与 389-r1 在盘件**逐字节相同**；反向双闭环独立复刻（两路突变 rc=1＋还原
核 rc=0 逐字节还原）；四数 94/83/0/11 对 59-w 基线逐名作差 +11 全本腿、removed=0、基线 83 枚
颜色零漂移；两尺＋d22scan 实扫 clean＋正控 35/0 全绿。残余＝AC#6（400ms 停播链／产品麦面零
转写／−30 裕量重测）——59-v3 §7 已裁死归票 26 集成面，本腿确认 389-r1 未越界、未假报。

## 1. 本腿自己跑了什么（复核凭据，非转抄）

- **读数逐字复跑**（389r1 §7 处方逐字：`WISP_AEC_SPIKE=1`＋sherpa PATH 前缀＝A1307 在册形状；
  落点导出树）：rc=0 PASS。四格：−45 静态链 good 2＝**5.41**／collapse 2＝**11.50**（59-w §4
  基线逐字）；产品门动态开 **0.00/0.00**；−30 轴开/关全 0。正控：control-26（−26.70）两轴过、
  control-36（−36.58）−45 过／−30 被吃——389r1 §3 表逐格对上。
- **产出件逐字节同**：`bargein-readings.json` 与 389-r1 在盘件 byte-identical（复跑前先快照到
  `389v1\out\bargein-readings.pre-v1.json`，终态与原件相同，覆盖可还原）。
- **反向闭环独立复刻**：①矩阵内 −45/关＝11.50 基线复现；②389-r1 突变件内容（收紧接线
  `g.wins[w]`→`false`）经重键 overlay 在导出树复跑 **rc=1**，断言逐字
  `dynamic gate is a no-op on the collapse state: on=2 >= off=2`；③**本腿自产突变件＋自写
  overlay** 独立复刻同 rc=1 同断言；④还原核（摘 overlay）**rc=0**、盘上读数件逐字节还原。
  突变后动态开读数与关逐格同（5.41/11.50）——读数差异的全部来源就是被切断的那段收紧接线。
- **四数**（导出树、文档化跑法）：**94/83/0/11**（rc=0）。对 59-w 盘上基线日志（83/73/0/10）
  逐名作差：**added 11 枚全本腿**（audio 包 10 枚 bargein 单测＋读数驱动 1 枚按设计 SKIP）、
  **removed=0**、**基线 83 枚颜色零漂移**（逐名 PASS/SKIP 状态比对 0 差）。
- **两尺**：gofumpt `-l internal/audio`＝0 行；vet 三形状 rc=0（无 tag／`-tags webrtcap`／
  `GOOS=linux`）；build 两形状 rc=0。
- **d22scan**（Git Bash，导出树脚本自定根）：实扫 **clean** rc=0；正控 selftest **PASS=35
  FAIL=0 SKIP=0（RUN=77）**逐字同。分母实测 bans#1-5 internal/=**286**、cmd/=**51**、
  ban#8 internal/=**726**（§4 注记）、ban#6 frontend/=116、ban#7 internal/tools/=25。

## 2. AC#1 产品形态读码核（bargein.go 540 行＋两测试件全文）：**四形状全过**

1. **静态三档＋三件套动态门控形状**：`Evaluate`＝VAD（调用方供给）∧能量∧ASR 非空转写
   （bargein.go:513-515）；三件套＝`CollapseDetector`（中位 −6／绝对 15／持续 3、inactive 窗
   跳过不破 streak、一次每 streak、在持查询 `Sustained`，:216-253）＋`FarCorrelationDetector`
   （floor 0.8；②③合流一量——产品 far 真实可得，与 59wr1 §5.5 换算判语一致，:274-357）。
   常量与 59wr1 §5 同值零发明。
2. **−45 标称＋持窗收紧 −30＝59-v3 两常量产品化**：`BargeinNominalEnergyDbfs=-45`／
   `BargeinTightenedEnergyDbfs=-30`（:80-83）；收紧查询在 episode 起点窗（:505-512）、报警
   解除（streak 断）回标称。**收紧是裕量交换不是抑制开关**：只动阈值不旁路判决，−26.70 真话
   在收紧档仍过（单测 `TestBargeinGateDynamicTighten` 在 untrusted 窗内钉死＋读数正控实测），
   被吃的塌峰伪迹带实测 −32.24/−35.59（读数件逐枚在案）。
3. **双接缝同达唯一决策路径、无预 OR**：`Feed`（:445-486，帧级、内部 lag 环）与 `ObserveWindow`
   （:528-532，窗级原值）都只写 `g.wins`；判决只发生在 `Evaluate` 读 `wins[w]`。读数驱动以
   **原始窗值**进门（ERLE dB 经 10^(f/10) 幂重建进 detector、相关原值进 floor 判定）——判据
   在产品门内真实行使；突变闭环反向证明：切断门内收紧接线读数即翻红。
4. **零 goroutine、零 internal/speech 依赖**：全文唯一 import＝`math`；`BargeinEpisode` 纯数据
   接缝（:375-380）把 VAD/ASR 留在调用方。SKIP 11 枚＝59-w 基线 10 枚＋新读数测试（无
   `WISP_AEC_SPIKE` 按设计 SKIP），零丢失。

## 3. 离线复用合法性裁（派单④）：**成立，效力边界如实划**

- **合法性**：两态读数的对象是「产品形态门在 59-w 材料集上的误触发率」——门判决是
  （episode RMS，起点窗 trust 态，转写）的纯函数；59-w 的 episode 全集（含 ASR 转写）与逐窗
  判据原值都在盘上工件里，重放与实测是同一材料集。**校准锚**：静态链格（−45/关）与 59-w rig
  实测逐格同值、误触发 2/2 逐枚同名同电平（「你能」−38.85／「哎我」−32.32／「我们」−35.59／
  「听我」−32.24）——材料切分、episode 集、静态判决三方对齐；动态开格与关格的差异因此**只**
  来自收紧判决本身（突变闭环再反向钉死）。零 rig 零 ASR 重跑不损这两格读数的效力。
- **event 窗轴坑已修核**：材料窗段取 `probe.json` 的 events（与 59-w 门材料同源切分）；
  autocorr.json 只取逐窗 selfCorr 值（>0 才有效）、**不取其事件表**——代码逐行核过
  （bargeinEventSpan／bargeinRunWindows）。
- **效力边界（不能由此宣称的）**：①产品在线 `Feed` 路径能否算出同值判据（帧级流 vs 窗级重放
  的数值等价）＝票 26 集成面（Feed 等价形状已有单测钉住：`TestBargeinGateLagShift`／
  `TestBargeinGateWindowSeam`）；②59-w 材料集外的新材料；③400ms 停播与产品麦面零转写
  （AC#6）。误触发样本 n=2+2 如实记账。

## 4. d22scan 分母注记（诚实差，不涉判）

- 派单四数 286/51/726/**174** 的末位：本腿实测 **176**（389-r1 自报同为 176）。174 是 59-w 锚
  的旧值；174→176 的 +2 属 59-w 锚后其它波的 cmd/ 提交——389-r1 pathspec（commit stat 4 枚）
  **零 cmd 写面**，归因不沾。前三位 286/51/726 与申报逐字同。
- ban#8 design/ 本腿量 **14**（389-r1 自报 39）：59-v3 §8 既有定性——导出净树 vs 真工作树
  未跟踪 design 件；非差错。

## 5. AC247 偶发红归因核：**维持既有归因（设备面负载抖动），与 389-r1 无因果**

- 台账 `pending-and-issues.md:29266`（A1327 收件时照录）：`TestAC247LiveMicrophoneLevelsReach
  TheBallSeam` 偶发 1 红＝单跑/复跑 PASS，负载抖动形状、12-r2/384 先例同类。
- 本核加两条独立证据：①389-r1 写面（commit stat）零 cmd/wisp 文件——AC247 是 cmd/wisp 真麦
  活体测试，编译面无交集；②本腿 internal/audio 全套四数 FAIL=0——抖动不在 audio 面。
  不新开调查。

## 6. 六格逐条判

| 格 | 判 | 依据 |
|---|---|---|
| AC#1 组合门产品形态 | **能翻（绿）** | §2 四形状全过（全文读码＋突变反向证明） |
| AC#2 两态读数 | **能翻（绿）** | §1 逐字复跑 rc=0、四格逐格复现、JSON 逐字节同 |
| AC#3 反向判据 | **能翻（绿）** | 开关关=基线复现＋双突变 rc=1＋还原核 rc=0 逐字节还原 |
| AC#4 零回归 | **能翻（绿）** | 94/83/0/11、+11 全本腿、removed=0、基线漂移 0 |
| AC#5 纪律格 | **能翻（绿）** | pathspec 4 枚零越界；thresholds/golden/AEC3 源/deps.toml 未动；⛔ push 遵守；两尺＋d22scan 全绿 |
| AC#6 集成边界 | **⛔ 本面不可翻（票 26 集成面）** | 400ms 停播／产品麦面零转写／−30 裕量重测——389-r1 §8 未越界未假报，59-v3 §7 判语维持 |

**可否带残余结案：可以。** 残余定性：①AC#6 集成边界三件真数在票 26 集成面出＋−30/−36 裕量
2–6dB 重测加严——票 26 已开，标定集已备好（`bargein-readings.json` 的 ThresholdDbfs/untrusted
逐枚在案，−26.70/−36.58/−32.24/−35.59 四电平就是「收紧档裕量交换」的实测标定）；②误触发
小样本 n=2+2——同集重放效力成立（§3），泛化靠集成面复验；③处方可移植性注记：389-r1 的
overlay.json 用仓内绝对路径键，导出树复跑须重键（本腿已备 overlay-theirs-rekeyed.json／
overlay-v1.json 凭据）——非缺陷，是「仓外复刻」的工程注记。

## 7. 处置建议

1. **389-r1 按本判语结案**（AC#1–#5 翻绿带残余；残余全数归票 26 集成面，无 389-r1 欠账）。
2. 票 26 集成面验收直接复用本门读数面与四电平标定集；Feed 在线路径与窗级重放的数值等价性
   在集成面首格顺手核一枚即可。
3. AC247 偶发红维持台账归因；编排者收 390-r1 时若再现，照 12-r2/384 先例单跑/复跑处置。

## 8. 写面与凭据

- 仓内唯一写面＝本件（commit pathspec 单枚）。⛔ 票面/台账/specs/PLAN 零触碰；⛔ push；
  在飞 390-r1/387-r1/386-r1 写面零触碰、零归因。
- 仓外凭据 `D:\tmp\wisp\20261031\389v1\**`（只建不删）：`export\`（git archive `651067412`＋
  junction third_party）＋`export.tar`；`logs\readings-v1.log`（rc=0）、`logs\mutant-theirs-
  rekeyed.log`（rc=1）、`logs\mutant-v1.log`（rc=1）、`logs\restore-v1.log`（rc=0）、
  `logs\four-numbers-v1.log`（94/83/0/11）、`logs\d22scan-export-v1.log`（clean＋正控 35/0）；
  `out\bargein-readings.pre-v1.json`（复跑前快照）；`mutant-v1\`＋`mutant-theirs\`＋两份
  overlay（独立突变闭环凭据）。
