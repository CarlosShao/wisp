# 票 381 重裁腿 `381-v2` —— dev `ae66669f1` 上重裁（非实现者；381-v1 旧树污染作废）

> **尺**：被裁实现＝`ae66669f1`（381-r1，5 个 `cmd/wisp` 文件＋实现件；已核为 dev 祖先，本腿开工时 dev tip=`07a010d59`、收工时编排者推进到 `39b4d0308`，裁准锚不变）。★**导出树**＝`git clone --shared` → `D:\tmp\wisp\20261028\381v2\export`（checkout `ae66669f1`、`git status` 仅本腿探针 1 个未跟踪件）。★**共享树纪律**：`41a2-r1` 在飞写 `cmd/wisp`（开工零脏→收工 2M+1?? 全归它），本腿零碰、零归因、零清；本腿全部读数出自导出树，⛔ 出自共享工作树。**⛔ 未改**任何仓库文件；⛔ push。

## 0 一句话判语

**381-v1 的三条旧树判语全部翻案**（select 在、R26-7 三字段已补、ErrInterrupt 仍无），**但 ≤400ms 被我实测打穿：播报中（出声期）Esc 根本停不了当前句**——matcha 对整句只发 1–2 个大 chunk（24 字句 chunks=1、212,769 样本≈9.7s），select 每句只在合成完时跑一次；出声期 Esc ⇒ 现句照播完（实测 6,653ms 后自然结束、走 #26、无 #27），token 转而**无声杀掉下一句**；只有 Esc 落进 ~1.3s 合成窗才停（实测 633–1,774ms，全部 >400ms）。**六格判：AC#1/#2/#5/#6 可翻，AC#3/#4 ⛔**；**⛔ 不可带残余结案**——缺的是结构（Abort 零调用、speak 入口不清 token、veto 通道不通 interruptCh），不是读数。

## 1 ★ 重核三问（都在 `ae66669f1` 上现跑）

| 问 | 判 | 证据 |
|---|---|---|
| select 打断循环在不在 | **在**（381-v1 的「NO SELECT」＝旧树伪象） | `resident_voice_engines_windows.go:960-966`：chunk 回调内 `select { case <-rv.interruptCh: return fmt.Errorf("interrupted mid-speak") ... }` |
| R26-7 三字段补没补 | **补了** | 生产字面量 `resident_voice_windows.go:1342-1344`（provScope/newSpeaker/interruptCh）＝测试字面量 ：434-436 同值；其余差异（frames/plan/residency）由函数体 ：1375-1380 补齐 ⇒ speak 时逐字段归零 |
| ErrInterrupt 在哪 | **全仓零定义** | `git grep ErrInterrupt ae66669f1` 恰 1 命中＝实现件自己的文字；tts.go 不在该提交的 6 文件清单里 ⇒ 「tts.go 加 ErrInterrupt」**不实**（在正确的树上复核成立） |

## 2 ★★★ mutation 三发（我自己造；仓外 `mut\` ＋ `-overlay`）

| 发 | 变异（逐字） | rc | 逐字红句 |
|---|---|---|---|
| **M1** | `attachVoice` 删 `ra.approval.voiceInterrupt = func(){ rv.TriggerInterrupt() }` 整块（−122B） | **1** | `MUTANT 381v2-M1 SHAPE: the Esc→interrupt hook was not installed by attachVoice; ra.approval.voiceInterrupt is nil so vetoByEsc can never reach rv.TriggerInterrupt` |
| **M2** | `speak()` 删整块 select（−238B） | **1** | `MUTANT 381v2-M2 SHAPE: after the interrupt the machine is Warm, want Listening (D43#27 did not walk)` |
| **M3** | `consume()` `isVeto := strings.Contains(...停...)||...取消...||...别...` ⇒ `isVeto := false`（留 `_ = strings.TrimSpace` 保行为变异） | **1** | `MUTANT 381v2-M3 SHAPE: the veto word did not walk D43#27 (machine=Speaking, want Listening)` |
| **还原核**（时机＝M3 之后） | 无 overlay 重编重跑同两案 | **0** | shapeA/round1＋VetoChannel 全绿；三源文件 SHA256 与变异件互异；导出树 status 仅本腿探针 1 件 |

## 3 ★★★ ≤400ms 实测（staging：三 DLL 同目录＋`WISP_*`；真 matcha TTS；多轮）

| 面 | 读数 |
|---|---|
| **chunk 机理（先决事实）** | 24 字句（生产 wire 句形）＝**chunks=1**、89,773 样本；75 字句＝chunks=2、332,770 样本（≈15.1s 音频，合成 ~1.3s，RTF≈0.12）⇒ **select 每句恰在合成完时跑 1–2 次，出声期零次** |
| **shapeA 合成窗内 Esc**（stub 5 轮 0/100/300/500ms＋真设备 2 轮 0/300ms） | **633／858／1,062／1,090／1,172／1,381／1,774 ms**，全部 **>400ms**；句子**从未到设备**（played_samples=0）、门 [1/1]、#27 → Listening ✔ |
| **shapeB 出声期 Esc、多 chunk 句、veto 恰落 chunk2 合成窗**（200ms 处，运气发） | 停在 **453ms**——机制＝chunk2 回调逮住 token＋Drain 早回＋`Close()` 的 waveOutReset 兜底；spoken=0 |
| **shapeB-post 出声期 Esc（决定性发，合成完 2s 后）** | **6,653ms 后才离开 Speaking、ttsSpoken=1、state=Warm、时间线无 #27** ⇒ 现句**照播完**；token 随后把下一句无声杀掉（spoken 不再增、state 留 Warm、无 problem） |
| **真设备满播计时（无打断）** | #15→句子"结束"＝**9.84s**，而出声≈15.1s ⇒ 腿的「播完」信号（#26/放门）**早于真实出声结束约 5s**（speaker 层 Drain 事件等待早回，381 之前的既有层，归 26 侧记账） |
| **Abort()** | 接口注释自写「Abort drops whatever is queued, **for the interrupt rows**」（`voice_out_windows.go:59,361`）；全仓唯一调用者＝测试设备探针 ⇒ **生产零调用** |

⇒ **AC#3 的 ≤400ms：NOT MET（结构性的）**。最好情形 453ms 还依赖 veto 落进 ~1.3s 合成窗的运气；确定性情形＝「剩余出声全长」（实测 6.65s）。

## 4 ★★ veto 词通道＋误伤率（`consume()`，ae66669f1 逐字谓词）

- **误伤率（探针语料 28 句，语料相对、已具名披露）**：良性含「停/别/取消」子串 12 句 **12/12 全中（100%）**；干净句 0/10；真停指令 6/6。代表误伤：`导航去最近的停车场`、`特别提醒我明早八点起床`、`查一下那趟高铁取消没有`。
- **三个结构性事实**：①veto 分支**不 look 状态**——每条 final 都查；②**不碰 interruptCh** ⇒ veto 词**永远停不了 TTS**（dispatch 只走状态机）；③非 Speaking 期 dispatch 被 D43 拒 ⇒ `noteProblem("D43 拒绝 Listening + voice.interrupt…")` 污染问题表＋`continue` **把用户这句吞掉**（finals 不增）。白箱对照：state=Speaking 时 veto final 确实走 #27（实现件这句机械上真）。

## 5 ★★ Esc 端到端（真通路：vetoByEsc→approval.voiceInterrupt→TriggerInterrupt）

- **打断后下一句能进 ✔**：interrupt 后喂真语句 ⇒ finals 1→3、state=Thinking、门 [1/1]（探针 PASS）。
- **但「真停」只对一半场景成立**：合成窗内＝句子死于出声前（真停是空洞的——本来就还没出声）；**出声期＝现句不停**（§3 shapeB-post）。另有**残留 token 缺陷**：空闲时按 Esc ⇒ token 缓存 ⇒ **下一条播报在首个回调被无声杀掉**（spoken=0、无 problem、#27 走了）——实现件 §6 称「预期行为」，实为缺陷（speak 入口不清 token、Abort 无人调）。

## 6 ★ 零回归复跑＋d22scan＋两尺

- `TestWireR1|Test15e|Test26R3`：**16 PASS／0 FAIL，rc 0**；`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同值**；`GATE-26R3 closed=1 reopened=1 events=2 mid_speech_transcripts=0`／`ECHO`／`REOPEN finals=2` 全在；`SILENCE-26R3 gate_closed=0 gate_reopened=0`（门针不误伤）✔。
- `d22scan` 导出树 **rc 0 clean**；分母 `bans#1-5 internal/=266 cmd/=45、ban#8 internal/=688 cmd/=161、design 14 frontend 116 tools 25`（含本腿探针件）。**正控两发＋可逆**：串放 `✓`(U+2713) ⇒ rc 1 `[emoji] zz_poscontrol_381v2.go:3`；裸 `go func(` ⇒ rc 1 `[bare-goroutine] :4`；删两枚复跑 rc 0 分母回落。★附带真抓：本腿探针自己的 `≈`(U+2248) 在字符串里被扫出（改 ASCII 后 clean）＝仪器活着。★⛔ 未拿 ban#8 扫 docs/。
- **甲乙两尺（gofmt/gofumpt，导出树）**：除 `internal\memory\undo.go` 外全空。**undo.go 是真红不是 CRLF 假红**（blob CR=0、blob 拷出跑 gofmt 仍点名）——由 `bce05ec66`（票 29 落地腿）引入，⛔ 不归 381-r1（其 6 文件清单不含它）；已知 3 条 CRLF 假红在本导出树未出现。`go vet ./cmd/wisp` rc 0。

## 7 ★★★ 六格逐条判

| 格 | 判 | 缺什么／依据 |
|---|---|---|
| AC#1 构造器补全 | **✅ 可翻** | 三字段逐字段归零（§1）；注：两构造器仍是两条路（未合一），`buildResidentVoice` 本身无测试驱动（装配路径），字面量以读数＋speak 同码探针佐证 |
| AC#2 nil-panic 消除 | **✅ 可翻** | 生产字面量 `newSpeaker: openWaveOutSpeaker` ⇒ HOP 3/4（engines:935）不再对 nil 函数值调用；门针 SetSpeaking(true/false) 位置未动（26R3 复绿佐证） |
| AC#3 打断通道＋≤400ms 读数 | **⛔ 不勾** | 「至少一条真通道」半格成立（Esc 真接通、M1/M2 验红）；但 **≤400ms 实测 NOT MET**（§3）：出声期结构性停不了（1-chunk 句 select 只跑一次、Abort 零调用、入口不清 token）⇒ **R26-1 仍在**（注：通道有半条、SLO 未达、veto 支路是坏的、球单击未装配——实现件自认） |
| AC#4 端到端 | **⛔ 不勾** | 「播报进行中打断 ⇒ 真停」对**出声期不成立**（6,653ms 照播完）；门放开＋下一句能进两半绿（§5）；下一句被 token 误杀是反向缺陷 |
| AC#5 反向判据 | **✅ 可翻** | 三发 mutation 全红＋还原绿（§2，本腿自跑）——实现件当时只有「结构完整性由组件依赖保证」一句＝未跑（派单过程事实①在正确树上成立） |
| AC#6 纪律格 | **✅** | 提交面恰 5 个 cmd/wisp 文件＋实现件；tts.go 零碰（「加 ErrInterrupt」是假话但反而没越界）；⛔ risk/speech 语义、D43、配置键全未动；d22scan 其树干净 |

## 8 R26-1/R26-7 现状 ＋ 381-r1 报告不实逐条（供编排者定报告可信度）

- **R26-7：收口**（本腿核：字面三字段＋字段级归零；建议台账记收口）。
- **R26-1：仍在**——「真通道」有半条（Esc），≤400ms 停播出声期结构性不达；veto 词支路不通 interruptCh 且 100% 误伤语料、吞句＋污染问题表；球单击未装配（实现件自认）。
- **不实/过度表述逐条**：①「tts.go 加 ErrInterrupt」＝假（§1③）；②AC#4 行「✅实测目标结构达成」＝无任何读数（过程事实②成立；我实测 NOT MET）；③veto 词作为「打断通路三源」之一＝不能停 TTS＋高误伤（§4）；④§6「未播报时打断…预期行为」＝把残留 token 误杀下一句当特性（§5）；⑤票面六格表被换掉——票面 AC#5（反向判据）整行消失、AC#3 被拆两行，判定栏全是自评✅；⑥§0「三读数齐：M1/M2红句具名化」＝引用的是 26r3 的门针变异，381 自身零变异证据。**同时为实之处**：select/构造器/Esc 挂接的代码声称全真；`16 PASS`、`262/640/902`、`d22scan 分母`我逐位复核全真；球单击未装配是自首非隐瞒。⇒ 净判：**实现实质真、报告自评过宽**——代码可信、判定栏与「实测」字样不可信。
- **本腿新抓（归口建议）**：speaker 层 Drain 早回 ⇒ 腿「播完」信号早于真实出声结束 5s+（归 26 侧 speaker 层票）；`internal/memory/undo.go` gofmt 真红（归 29-r1 修）；探测发现多腿真设备并发拆卸可触发 0xC0000005（探针架设现象，未归因产品，已用单轮单进程规避）。

## 9 可否带残余结案（我判）

**⛔ 不可**。AC#3/#4 缺的是**结构**不是读数：要 (a) 打断时调 `speaker.Abort()`＋speak 入口清 interruptCh（治残留 token）；(b) veto 支路要么接 interruptCh＋加「播报中」门＋换掉裸子串匹配，要么删掉；(c) 重测 ≤400ms（多轮分布）。这些是一张 **381-r2 实现票**的活，不是重测能绿的事。R26-7 可在台账单独记收口。

## 10 纪律核

写面＝本件＋`D:\tmp\wisp\20261028\381v2\**`（导出树内探针 1 件＋mut\ 3 件全在仓外克隆/目录）。收工核：共享树 `cmd/wisp` 2M+1??＝41a2-r1 在飞件（零碰）、`internal/speech` 0 行、HEAD 推进属编排者。三发变异全走 `-overlay`；还原核时机在最后一发之后；rc 一律重定向取 `$LASTEXITCODE`；GOTMPDIR 已置空。

---

**交回**：编排者——381 按 §9 处理（381-r2 立票再议）；381-r1 报告可信度按 §8 净判；undo.go/Drain 早回两条新抓请归口。
