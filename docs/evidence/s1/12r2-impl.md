# 票 12 落地腿 `12-r2` —— AC#9（registry A25）「回复后 3s 回落有实测记录」（写腿）

> **尺**：开工 `HEAD=07f6c29bc`（`dev`，三核过：索引零已暂存、`cmd/wisp` 无他人脏）；正式读数现跑
> 2026-10-31 17:5x–18:1x +08，staged 导出树锚 `08ffbd72d`（编排者当日只追加 `.scratch` 派单存档，
> 代码面与开工核同树）。★读数逐条现跑，⛔ 转述。
> **写面（6 枚，逐名）**：新 2 枚（`cmd/wisp/settle_trace_12r2.go`＝「最后 token 时间」载体＋样本件、
> `cmd/wisp/settle_trace_12r2_windows_test.go`＝≥5 样本实测案＋反向判据常驻牙）＋改 3 枚
> （`cmd/wisp/run.go` +14/−1：`agentRuntime.settle` 字段、`consoleSink.settle` 字段、`Publish` 的
> EvTextDelta 分支 1 行 T0 戳、`Sink:` 字面量 1 处；`cmd/wisp/resident_voice_engines_windows.go`
> +25/−2：residency 的 `settle`/`ttsLoader` 两字段、TTS 工厂 1 行替身缝、`openScope` 的 Logf/
> FreeOSMemory 两个既有闭包各加 T1/T2 戳；`cmd/wisp/resident_windows.go` **+8/−0 纯插入**＝30-w-r1
> presenter 装配段之后的兄弟挂接段，`git diff` 可核其行零触碰）＋本文件。
> **零字节**：`internal/statemachine/**`（D43 表/events.go/timeouts.go 一字未动）、`internal/session/**`
> （Settle/EnterWarm/3s 静置窗零字节，只观测）、`internal/{risk,speech,panel,observe,agent,proc}/**`、
> thresholds/golden、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、
> **`docs/SLO.md`**（见 §2 落点判定）、票面、台账、`docs/PLAN.md`、`docs/specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\12r2\{bin,logs,mut,tmp,export,head.tar,d22scan.exe}`。
> 突变在仓外 `export` 导出树（`git archive HEAD`＋overlay 本腿 5 件后变异），共享工作树零突变窗口。
> **staging**：`go test -c -o bin\wisp.test.exe ./cmd/wisp`（主仓树 build）＋三 DLL 同目录
> （`third_party\sherpa-onnx\{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll`）＋
> `WISP_SILERO_MODEL/WISP_TTS_DIR/WISP_SPIKE_MODELS/WISP_KWS_DIR/WISP_PUNC_MODEL/WISP_15E_TMP`
> ＋binary 直跑 cwd=包目录（381-v3 清单同套）。

## 0. 一句话

AC#9 的「最后 token 时间」载体落在 resident/session 侧三条**已经存在**的缝上（流式 sink 的 delta 分支、
scope 自己的 SCOPE-SETTLING 审计行、dispose 级联的 FreeOSMemory 步），球那台机器的 Sink 面与
statemachine 冻结面一个字节未碰（机器效果投递面归票 273）；产线走位实测 ≥5 样本全报：
**token→Settling 1201–1206ms（≤3000ms 判据全过）、Settling→回落 3000–3003ms（冻结 3s 静置窗真值）**；
摘任何一枚戳 ⇒ 实测案红（两发仓外变异＋还原核）；六族＋LATENCY 262/640/902 逐位同值零回归。

## 1. ★ 必给读数一：AC#9 字面＋PLAN 现核（第 20 处）

- **票面 `:150-157` 逐字**：「（2026-09-20 编排者新加，registry A25）回复后 3s 回落有实测记录：PLAN `:1424`
  的 S1 第一判据原文是「打字→得回复→**3s 回落**」，逐项对账后确认这是**四项里唯一从未被量过的一条**
  （`wisp slo` 量空闲态、`balldebug -diff` 量逐态驻留，都不量"最后一个 token 之后多久进
  `Settling`/`Sleeping`"）。完成判据：①从"最后一个 token 送达"起算到进入 `Settling` 的耗时 **≤3000ms**，
  **必须用单调钟**（`observe.Timeout` 那条路；墙钟差判超时是 D22 禁令 #4）；②**≥5 个样本全报**，不许用
  平均抹掉坏尾部、不许重测到运气好的那次，超线就录 FAIL；③机器可判的形态优先：若 D43 的回落行
  （`EvAfterReply`/notify→Warm 那条）没有暴露"最后 token 时间"，**加一个事件钩子/接缝就是代码改动**，
  写进本票而不是绕过去用肉眼数秒表。」
- **PLAN 现核（第 20 处）**：票面引 `:1424` 是行号漂移（12-r1 已核），S1 行在 **`PLAN.md:1430`**，逐字：
  「打字→得回复→3s 回落；空闲 private RSS ≤25MB 实测通过；`wisp run "..."` CLI 可用；C21 设计 token 落地」。
  **判定：验收语义与「3s 回落」对得上，不停手**——票 12 自己的 What-to-build 把它钉死为
  「Settle: result presentation done → 3s → DisposalScope session teardown (incl. FreeOSMemory) →
  RSS back to idle cap ≤10s」：**「3s」＝回复呈现完之后的那扇 3s 静置窗（D43#32 的 session 半边，
  `session.SettleStaticDefault=3s` 冻结默认），「回落」＝窗满后的 dispose 级联**；AC#9① 的 ≤3000ms
  钉的是「最后 token → 进 Settling」这条**进窗腿**（管道时延上限）。两处出处互证，无未定义。
- **现成缝（12-r1 预筛复核，逐条对上）**：①`internal/statemachine/events.go` 42 枚事件无任何时间载体
  （无 EvAfterReply/LastToken/tokenAt——本腿现跑 grep 复核同零）；②CLI `run.go` 跑完即退、无 Settling
  计时；③活着的回落走位在 resident/session 侧：`voiceResidency.openScope` 把 `session.Scope` 的
  Clock/Hooks/WarmTimeout（=[session].warm_timeout_sec）接进产线，Warm→`Settle()`（SCOPE-SETTLING
  审计行）→3s 静置→`onSettleExpired`（D43#32 触发器）→dispose 级联（FreeOSMemory 第 3 步）。

## 2. ★ 记录落点判定（派单第 3 条：「按 AC#9 字面判，拿不准停手问」）

**判定：实测记录落 `docs/evidence/s1/`（本文件）；`docs/SLO.md` 一字不动。**
理由（可核，不是自创）：AC#9 的字面只要求「有实测记录」，**没有像 AC#2 那样点名
「PASS recorded into docs/SLO.md appendix」**——同票两框一处点名附录、一处不点名，这个对比就是
票面自己的落点答案；往 SLO.md 加附录行才是自创章程位。落点因此为证据表＋票面引用，由对抗验收腿/
编排者按票面翻框。

## 3. ★ 接缝（最小侵入证明）：三条已有缝上的三个戳，零新钩子

| 戳 | 挂点（全部既有缝） | diff |
|---|---|---|
| **T0 最后一个 token 送达** | `consoleSink.Publish` 的 `EvTextDelta` 分支（run.go:1700 区，票 12 原生 sink、非 30-w-r1 行）——LLM 流式回复交付给消费面的最后一跳 | `consoleSink`/`agentRuntime` 各 +1 字段＋`Sink:` 字面量 1 处＋分支内 **1 行** `c.settle.markLastToken()` |
| **T1 进 Settling** | `openScope` 的 **Logf 闭包**（engines 文件既有接线）——scope 自己宣布 `session: SCOPE-SETTLING`（`internal/session/scope.go:596` 的审计行首，行号漂移只会让戳停掉、被测试当红抓住） | 闭包内 **1 个前缀判据 + 1 行** `markSettling()` |
| **T2 回落完成** | `openScope` 的 **FreeOSMemory 闭包**（engines 文件既有接线，级联第 3 步）——`onSettleExpired`＝D43#32 触发器 | 闭包内 **1 行** `markFellBack()`，置于真释放之前 |
| 载体本体 | `settle_trace_12r2.go`：三枚 `atomic.Int64` 单调钟戳（`time.Since` 底座，**测量不是超时**，D22#4 正面）＋pending/sample 相干守卫（无回复不成样本、family 换族不成样本、窗被 #33 唤醒取消则被下一次 Settling 覆盖）＋有界 64 样本 | 新文件，纯 Go 跨平台 |
| 注入半边 | `voiceResidency.ttsLoader`（镜像既有 `asrLoader` 的替身缝：nil=产线 matcha 工厂）＋`resident_windows.go` 兄弟挂接段 1 行（`src.run.settle = rv.residency.settle`） | +25/−2、+8/−0 |
| **⛔ 两头不碰的边界** | 球那台机器的 Warm 驻留走冻结 `DefaultTimeouts`（90s/#31）——**机器 Sink 的效果投递面归票 273**（273 AC#2 三形未选形、384-r1 书面「273 落真 Sink 时本缝零改动」），本腿不给生产机器带 Sink；`internal/session` 的 Settle/窗一个字节不动 | grep 可核 |

## 4. ★★★ 必给读数二：≥5 样本 3s 回落分布（staged 正式跑，单调钟、逐样本毫秒）

样本＝产线全链：`newRunFixture`（mockllm＝C5 缝替身）流式回复（T0 戳在产线 sink）→
`rt.voiceAnnounce=rv.announce`（wire-r1 生产字段）→ 产线 `speak()`（TTS＝C9 三方法桩、设备＝
wire-r1 speaker 桩，使被测腿只量回落不量合成）→ 产线 `#15/#26` 派发 → `enterWarm("announce-end")`
→ **`voicePlan.WarmTimeout`=1.2s**（=[session].warm_timeout_sec 的读数位，让样本秒级完成）→
scope `SCOPE-SETTLING`（T1）→ **冻结 3s 静置窗** → dispose 级联（T2）。机器桥＝测试侧按冻结表
自己的行走位（#4/#11，26r3/381r2 同款先例）；三枚戳全在产线缝上。

**staged 导出树正式跑（log `staged-measure.log`，rc=0，24.63s，5/5 全报）**：

| 样本 | token→Settling | Settling→回落 | token→回落 |
|---|---|---|---|
| 1 | 1206ms | 3000ms | 4206ms |
| 2 | 1202ms | 3003ms | 4206ms |
| 3 | 1201ms | 3000ms | 4201ms |
| 4 | 1202ms | 3000ms | 4202ms |
| 5 | 1202ms | 3000ms | 4203ms |

- **判据①**：5/5 ≤3000ms（分布 1201–1206ms；管道增量＝announce 尾＋执行尾，个位数 ms）。
- **「3s 回落」本体**：Settling→回落 3000–3003ms＝冻结 `SettleStaticDefault` 真值（计时器不早火、
  经 owner 队列到戳），**回落真的发生、真的 3s**——这就是 AC#9 标题要的那条实测记录。
- **还原核兼复现跑**（log `restored.log`）：5/5 同分布（1201–1203 / 3000–3001 / 4201–4204ms）。
- 机器走位逐行在案（log）：`SCOPE-WARM reason=announce-end warm_timeout=1.2s` → `SCOPE-SETTLING
  reason=warm.idle static=3s` → `SCOPE-DISPOSED reason=settling.expired cascade=…`，且
  `D43#15 首个纯文本 token 让球走 Speaking`、`D43#26 Speaking 到 Warm` 两行产线派发同 log 可引。

**诚实边界（⛔ 不藏）**：生产默认暖窗＝90s（scope 走 `[session].warm_timeout_sec` 默认、球机器走冻结
`DefaultTimeouts` 90s）⇒ 生产默认下同一条进窗腿 ≈90s＋管道，**AC#9① 的 ≤3000ms 在默认配置下按字面
必超**——该判据在本记录里的可判形态＝**staging 配置暖窗（1.2s，具名）下的进窗腿**；3s 静置窗与
dispose 级联不受该配置影响（冻结默认）。这与票面「Settle: result presentation done → 3s → teardown」
的 S1 语义（文本路无语音引擎、无保热窗）一致；语音路 90s 保热是票 28/D43#26/#31 的引擎热窗设计。
进窗腿在生产默认下的形状由本缝今后在任何配置下都可实测（摘缝判据见 §5）。

## 5. ★★★ 必给读数三：反向判据（摘接缝 ⇒ 不可测 ⇒ 红；三读数＋还原核）

全部在仓外导出树 `-overlay` 后变异（共享工作树零突变）：

- **基线读数**：§4 的 staged 正式跑 rc=0（5/5 全报）＝**读数一（绿）**。
- **M1 摘 T1/T2 戳**（export 副本删 `markSettling`/`markFellBack` 两行，vet rc=0、编译 rc=0）⇒
  `Test12r2SettleWalkFiveSamplesEndToEnd` **rc=1**，红句逐字 `timed out waiting for the reply's
  fall-back to land in the trace (T2)`（log `mut1.log`）＝**读数二（红）**：窗照常走完（`SCOPE-DISPOSED`
  仍在 log）但 **3s 回落不可测**——正是 AC#9 要防的「行为在、记录无处落」。
- **M2 摘 T0 戳**（export 副本 `Publish` 的 `markLastToken` 行换 `_ = c`，vet rc=0、编译 rc=0）⇒
  同案 **rc=1**，红句同形（log `mut2.log`）＝**读数三（红）**：无「最后 token 时间」⇒ 无样本。
- **还原核**：export 副本按主仓原样还原（`markLastToken`=1、`markSettling|markFellBack`=2 现跑计数）
  ⇒ 重建重跑 **rc=0、5/5 同分布**（log `restored.log`）；主仓 `git status -- cmd internal` 全程恰本腿
  5 件、两枚生产文件与 `git show HEAD:` 对回无变异。
- **常驻牙**：`Test12r2SeamCallSitesAreTheMeasurement`（AST 钉三枚戳调用点＋SCOPE-SETTLING 行首锚，
  357 同款形状）随套常跑——不做变异也红。

## 6. ★★ 必给读数四：零回归（六族＋262/640/902）

| 面 | 结果 | 出处 |
|---|---|---|
| 15e 6 案 | 6/6 PASS；**`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同值** | log `zeroreg-sixfamilies.log` |
| 26R3 3 案 | 3/3 PASS | 同上 |
| 41a2r1 5 案 | 5/5 PASS（含 24.2s/18.4s 两枚 KWS 真装载案） | 同上 |
| 381r2 5 案 | 5/5 PASS（含 37.6s 真设备案） | 同上 |
| Test43 9 案 | 9/9 PASS | 同上 |
| Test384 3 案 | **单发 3/3 PASS**（E2E 1.39s） | log `zeroreg-384-single.log` |
| 邻接包 | `internal/{session,statemachine,observe,config,watchdog,proc}` 全 ok（rc=0） | 门禁段 |
| **抖动申报（非零回归差）** | 六族合跑首遍 `Test384ResidentWatchdogEndToEnd` 红 1 枚：`never observed the tree back under the Armed cap … within 3s`——**纯负载敏感读数**（串跑前刚跑完 26R3 29.6s＋381r2 37.6s 两枚真设备重案），**单发复跑绿**，与 384-r1 记录的族串跑装载抖动同族（该案不经本腿任何写面：本缝只加原子戳与 ≤64 样本切片，零 goroutine、零定时器） | 两 log 并报 |

## 7. ★ 四数＋两尺＋d22scan 导出树＋正控

- **甲 `gofmt -l cmd/wisp/`**：空（rc=0）。
- **乙 gofumpt**（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe -l`，只喂本腿 5 枚 `.go`）：空（rc=0）。
- **`go vet ./cmd/wisp/`**：rc=0（含测试文件）。
- **d22scan**（独立 build→`D:\tmp\...\12r2\d22scan.exe`）：主仓（含本腿件）**rc=0 clean**（bans#1-5
  internal/=286、cmd/=51；ban#8 internal/=726、cmd/=176）；**导出树＋overlay 同 clean rc=0**。
  **正控**：种真 U+2264 进 `cmd/wisp/zz_d22_positive_control_12r2.go` ⇒ **rc=1 且逐字点名**
  `cmd/wisp/zz_d22_positive_control_12r2.go:5: [emoji] ban #8 glyph in scope cmd/ is banned (D23):
  non-comment text, string literals included; comments are exempt per Q-46(c)`；种后即删、复扫 rc=0。
  ⚠ **仪器在本腿自己身上咬了一口（照录）**：本腿测试串初稿含 U+2460（①）与 U+2264（≤）各一枚，
  首轮扫描 rc=1 逐字点名 ⇒ 已改文字，两树复扫 clean——尺是活的，且对写腿自己也咬。
- **红名册**：本腿按族尺跑（全包直跑在本机撞 33 族真窗口案，43-r1/384-r1 同口径不作读数）；
  上面六族＋邻接包＝起手零新增红（384 串跑抖动单发绿，见 §6）。

## 8. 未做／具名残余（⛔ 藏）

1. **球机器半边（#31/#32 的 machine walk）不在样本里**：其 Warm 驻留是冻结 `DefaultTimeouts` 90s，
   且机器 Sink 效果投递面归票 273（AC#2 三形未选形）——本腿零发明、零触碰；本记录测的是**同一张
   D43 行的 session 半边**（`Settle`→3s 静置→`onSettleExpired` 级联），球机器那半边的 Sink 落地后
   本缝零改动即可把 T1/T2 也钉到机器转移上。
2. **文本路今天仍无回落走位**（12-r1 预筛发现维持）：`submitTask`→`execute` 不派发任何 D43 行、
   CLI 跑完即退——本腿交付的是载体＋活走位实测，把文本任务接进状态机是行为改动、超出接缝派单，
   留给票面/编排者裁定。
3. **AC#9① 的生产默认形状**（§4 诚实边界）：≤3000ms 判据按 staging 暖窗 1.2s 判；默认 90s 下该腿
   ≈90s＋管道。判据文本、配置位、两套形状都已具名，翻框裁量归编排者。
4. **273 与 30-w-r1 边界**：本腿未给任何生产 `statemachine.New` 带 Sink（grep `statemachine.New(` 产码
   构造点仍 2 枚均无 Sink）；`resident_windows.go` 只 +8 纯插入，30-w-r1/43-r1/384-r1 挂接行零触碰
   （diff 可核）。

## 9. 停手条件自核（派单三条「未定义即停」）

1. **PLAN 语义对不上？**——现核 `PLAN.md:1430` 与票 12 What-to-build 互证，语义对得上（§1），未触发。
2. **接缝必须动 statemachine？**——没有：三条缝全在 cmd/wisp 既有代码上（sink 分支、Logf 闭包、
   FreeOSMemory 闭包），`internal/statemachine` 零字节，未触发。
3. **与 30-w-r1 presenter 装配正面冲突？**——没有：挂接段是其兄弟段（+8 纯插入），presenter 装配行
   零触碰；样本链走的是 legacy announce 路（presenter=nil 的 CLI 同形分支），未占其 Dispatch 调用点，
   未触发。
