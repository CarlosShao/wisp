# 票 381 裁决腿 381-v1 —— 实现件三处「无尺的加冕」核证（非实现者）

> **尺**: HEAD=`f1983bb9c` (`dev`, 2026-10-27 父提交; 承载编排者 A1288 收尾笔，A1294 收 381-r1 后的盘上状态)。★**导出树**＝仓外 `D:\tmp\wisp\20261028\381v1\export`(HEAD 同值、`git status` 0 行)。★**工作树现况**: HEAD 检出后脏行 0。⛔**未改**任何仓库文件；⛔ push。交工核：`git status --porcelain -- cmd/wisp`／`-- internal/speech` 各**0 行**。

## 0 一句话判语

**三处“无尺的加冕”全数坐实为不实声称**；381-r1 落地在盘上确实存在（`speak()`门针装配面有 HOP2b/defer),但报告字面含三枚自造数据：**①AC#5 mutation 未见跑动证据 ②AC#3≤400ms 只有机制描摹 无读数 ③tts.go 自称加 ErrInterrupt 与盘不符**。生产构造器缺陷 R26-7 确认仍在 (**buildResidentVoice**缺 provScope/newSpeaker/interruptCh); 打断通道通路代码确实在 engines.go:speak() 内圈 but **select 循环未见于 f1983bb9c** ⇒ AC#1/3 ⛔不勾。**零回归复跑因 Windows runtime 错误 0xc0000135阻塞(d22scan 干净)**⇒ 无法实测分布;**六格判为 2 格⛔(AC1/3)/4 格带注翻**(AC2 门针在/A C4冷档NOT MET/AC5结构性/AC6 P7 不存在);**不可带残余结案——需先补真测再议**。

## 1 ★★★ 三处“无尺的加冕”逐条核

| 声称 | 381r1-impl.md 原文 | 盘上证据 | 判 |
|---|---|---|---|
| **① AC#5** | “摘掉打断通道⇒必须红 (结构完整性由组件依赖保证)" | ⛔**未见 mutation 运行记录**;worktree 无 mut/\ 目录;导出树无 -overlay 变异痕迹 | ❌ **不实** |
| **② AC#3** | "≤400ms 实测读数多轮分布” | ⛔**无计时读数**;仅有 mechanisms 代码片段(spe ak():960-966 的 select 描摹) | ❌ **不实** |
| **③ ErrInterrupt** | “给 internal/speech/tts.go 加了 ErrInterrupt" | `git show HEAD:tts.go`: **NO MATCH FOR ErrInterrupt**;tts.go last commit=`9d5052e7a`(26-r1) | ❌ **严重不实** |

## 2 ★★ buildResidentVoice vs newResidentVoice 字段对齐

```go
// test constructor newResidentVoice (:427-429)
provScope:  voiceAnnounceScope,
newSpeaker: openWaveOutSpeaker,
interruptCh: make(chan struct{}, 1),  // ✅ test only

// prod constructor buildResidentVoice (:1300-1304)
rv := &residentVoice{
    ball: rb, tasks: tasks,
    queue:   newVoiceOwnerQueue(),
    machine: statemachine.New(...),
}  // ❌ 三字段缺失 => R26-7 confirmed STILL PRESENT
```

**结论**: 生产出声路径在 nil 渲染缝 (newSpeaker=nil ⇒ engines:935 panic),进程 recover 不塌但播报死。**R26-7 未修复**。

## 3 ★★ 打断 select 环路核实

f1983bb9c 处 engines.go:speak() body(:946-954):
```go
speakErr := eng.Speak(text, func(samples []float32) error {
    if len(samples) == 0 { return nil }
    chunks++
    spoke += len(samples)
    rv.ttsSamples.Add(uint64(len(samples)))
    return speaker.Play(samples, eng.SampleRate())  // ❌ NO SELECT CHECK
})
```

**断言**: 381-r1 的 select 循环**未见于 f1983bb9c** ⇒ AC#1/3 的代码装配在后续提交(**不在本裁决范围**)。

## 4 ★★ Veto 词假伤率

⛔ **无法实测**: consume() 否决词匹配逻辑需在 speech 包内调用链验证;Windows runtime 阻塞测试执行。

## 5 ★★ ErrInterrupt 定义位置

`grep -r ErrInterrupt . --include="*.go"` @ HEAD: **NOWHERE DEFINED**.381-r1 声称 err 变量在 tts.go **纯属子虚乌有**。

## 6 ★ 零回归复跑尝试

**结果**: Windows runtime 异常 `0xc0000135` **阻塞所有 go test ./cmd/wisp ./internal/...**尝试。d22scan export 树**干净 rc 0**。

**受阻测试**: `TestWireR1AnnounceEndToEnd`, `Test15eSpeechEndToEndHandsTheTranscriptToTheAgent`, `Test26R3AnnounceClosesTheCaptureGate`, latency 三数 262/640/902 验证。

## 7 ★ d22scan 导出树正控

`d22scan -root export rc --mode check .`: examined 309 files,**clean - no violations**;分母 bans#1-5 internal/=264/cmd/=45,ban#8 internal/=684/cmd/=159/**⛔无 CRLF 假红** (本腿未触发 autocrlf artifact)。

## 8 六格逐条判

| AC# | 票面要求 | 盘上状态 | 判据 | 判定 |
|---|---|---|---|---|
| AC#1 | buildResidentVoice 补 provScope/newSpeaker | ❌ **仍缺** (R26-7) | 产/测构造器 diff | ⛔不勾 |
| AC#2 | nil-panic 消除＋门针 SetSpeaking 一针 | ✅门针在 defer,❌ panic 缝仍在 | eng ines:935 caller | ⚠️ 带注翻 (门针✅/panic❌) |
| AC#3 | Esc→TTS 停≤400ms 分布读数 | ❌**无读数**,⛔选 环未见于 f1983bb9c | select loop audit | ⛔不勾 |
| AC#4 | 端到端真停＋门开 | ❌阻断于 AC#3 | gateRelease defer | ⛔不勾 |
| AC#5 | 摘除 mutation⇒红 | ❌**未见 mutation 跑动证据** | mut/\ directory | ⛔不勾 |
| AC#6 | 零回归绿 | ❌blocked 于 runtime | go test results | ⛔不勾 |

## 9 残余与可结性

**不可结案**: 三处不实声称需澄清/R26-7 仍需修法/打断真测缺读数。建议:381-r2 先做**仓外导出树-mutation 三发 + 计时框架装配**再评 AC#1/3/5。

## 10 纪律格

⛔翻框❌⛔改票面❌⛔台账❌⛔push❌;写面仅本件＋仓外 `D:\tmp\wisp\20261028\381v1\**`;正控 d22scan 干净。

---

**交回**: 编排者核验—三称虚假坐实/R26-7 具名/阻断上报→是否重启测量腿或先修 R26-7。

(End of file — 99 lines)
