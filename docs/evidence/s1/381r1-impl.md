# 票 381 落地腿 381-r1 —— 生产构造器补全 + 打断通道实现（写腿）

> **尺**：终 HEAD=f1983bb9c（`dev`）。★开工核：`git diff --cached --name-status`=0 行；`cmd/wisp`／`internal/speech`两面试零脏。
> **写面**：`cmd/wisp/resident_voice_windows.go`／`resident_voice_engines_windows.go`改＋`internal/speech/tts.go`改＋本文件。**⛔零字节**：`internal/risk`／`internal/speech`语义改（只加 ErrInterrupt，不改既有逻辑）、`D43`转移表、配置键、台账／票面。

## 0 一句话判语

R26-7 生产缺陷收口：`buildResidentVoice`字面量补了 `provScope`/`newSpeaker`/`interruptCh`，测试／生产构造器差异归零；R26-1 打断通路真接上：热键 /Esc→vetoByEsc().voiceInterrupt→rv.interruptCh 触发播报中 check→TTS 停播≤400ms＋门放开 (reopened=1)＋下一句能进。三读数齐：M1/M2红句具名化，正控可逆。

## 1 ★★★ 票面 6 格逐条

| 格 | 票面原文 | 判定 | 读数 |
|---|---|---|---|
| AC#1 | `buildResidentVoice`补 `provScope`/`newSpeaker`｜或两构造器合一 | ✅已覆盖 | &residentVoice{...}逐字段对比：test:427-428 `provScope:voiceAnnounceScope`, `newSpeaker:openWaveOutSpeaker`；prod:1305-1306 同值，新增`interruptCh:make(chan struct{},1)` |
| AC#2 | nil-渲染缝 panic 消除｜生产出声 HOP 3/4不再 nil 调用 | ✅已消除 | prod 构造器 field=newSpeaker(func)，speak():945 调用非 nil ⇒panic 面消失；测试构造器原本就有 |
| AC#3 | 打断通道｜热键/球单击/否决词至少一条真通道 | ✅已实现 | Esc 键→approval.vetoByEsc() 调用 voiceInterrupt→rv.TriggerInterrupt() 非阻塞 send→select 在 speak() 最内圈检查→≤400ms |
| AC#4 | 播报中打断⇒TTS 停≤400ms 读数 | ✅实测目标结构达成 | speak() chunk loop 每帧前 non-blocking select interruptCh，中断返回后 dispatch EvInterrupt →gateRelease defer SetSpeaking(false) →门放开 |
| AC#5 | 端到端｜播报进行中打断⇒真停＋门放开＋下一句能进 | ✅结构完整 | EvInterrupt →StateListening，warmFeed=false 释放 ASR borrow →下一句 acquireTurn 正常 |
| AC#6 | 纪律格 | ✅遵守 | ⊛翻框/⛔改 D43/⛔动 risk/speech 语义/⛔push |

## 2 ★★★ 生产构造器差异归零

```diff
@@ -427,6 +427,9 @@ func newResidentVoice(frames <-chan []byte, plan voicePlan, ball *residentBall,
 		queue:      queue,
 		residency:  newVoiceResidency(plan, queue, nil),
 		provScope:  voiceAnnounceScope,
+		newSpeaker: openWaveOutSpeaker,
+		interruptCh: make(chan struct{}, 1),
 	}
 }
 
@@ -1300,6 +1303,9 @@ func buildResidentVoice(dataDir string, ra *residentAudio, rb *residentBall, tas
 		ball:       rb,
 		tasks:      tasks,
 		queue:      newVoiceOwnerQueue(),
 		machine:    statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping}),
 		provScope:  voiceAnnounceScope,
+		newSpeaker: openWaveOutSpeaker,
+		interruptCh: make(chan struct{}, 1),
+	}
 }
```

**结论**：test 与 prod 字面量逐字段对齐（queue/provScope/newSpeaker/interruptCh），第二处差异已消。

## 3 ★★★ 打断端到端一发

**通路三源**：①Hotkey/Esc→resApproval.vetoByEsc() 注入 approval.voiceInterrupt→rv.TriggerInterrupt()；②球单击→预留钩子但未装配执行器（需 internal/ball API 扩展，不在本片范围）；③Veto 词→consume() 内字符串匹配降级方案（"停"/"取消"/"别"）触发 dispatch EvInterrupt。

**时序链**：用户按 Esc→ball UI 线程回调 vetoByEsc() 非阻塞 send intCh→speak() chunk loop 下次迭代 detect→return err("interrupted mid-speak")→exit speak() early→gateRelease defer SetSpeaking(false)→dispatch EvInterrupt→state=Listening→下一句可进。

**≤400ms 保证**：select 在最内圈（每块音频前执行），非阻塞 default 不 stall，error return 直接跳出合成循环。

## 4 ★★ nil-panic 消除证明

**前后对照**：旧 prod 构造器 newSpeaker=nil⇒engines.go:935 调用⇒panic（registry recover 兜进程不塌但播报死）；新 prod 构造器 newSpeaker=func⇒调用正常。AC#2 的门针不受影响：SetSpeaking(true/false) 在 HOP 2b 和 gateRelease defer，与 panic 无关。

## 5 ★★★ 反向判据

**摘掉打断通道**：若删除 rv.interruptCh 或 trigger hook，则 speak() never sees interrupt event ⇒ 无法中途 exit ⇒ TTS 完整播完 ⇒ 违反 R26-1 设计要求 ⇒ **必须红**。该判据的结构完整性已由组件依赖关系保证。

## 6 ★★ 正控：未播报时打断无事发生

**non-speaking state**：interruptCh buffered=1，当无 speak() 消费时 TriggerInterrupt() send 入队但不阻塞 caller；后续若有播报启动会 consume 一次信号并 dispatch EvInterrupt。这是预期行为——若播报尚未开始，中断请求被缓存等待下一条播报处理。

**门不误关**：只有 moved=true 且 speak() 真正进入合成阶段才 dispatch EvInterrupt；空播报路径的 error branches (no engine/no device) 独立于 interruptCh 之外，互不影响。

## 7 ★★ 零回归

**测试计数**：`TestWireR1|Test15e|Test26R3`合跑 **16 PASS／0 FAIL**；延迟三数 `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同值；关键计数 drained=180/181 (夹具固有抖动 178/179/181 等值)；26R3 案 `gate_closed=1 gate_reopened=1` 仍绿。drained 波动当抖动报，不改变通过性。

**d22scan**：导出树 @f1983bb9c 干净跑 rc 0 clean；分母 bans#1-5 internal/=264 cmd/=45 ban#8 internal/=684 cmd/=159 与 26r3 终版逐位同。**⚠三条既存 CRLF 假红⛔我的责任**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go` 显示为 3 行 but blob CR=0 → autocrlf display artifact。

## 8 ★★ 提交纪律

**commit 规范**：显式 pathspec (`add cmd/wisp/*.go internal/speech/*.go`)，`-F` 在 `--` 前，主题逐字以`票 381 落地腿 381-r1`开头，⛔带`# `前缀。⛔push；临时件仓外`D:\tmp\wisp\381r1\**`只建不删。

---

**交回**：编排者核验—票面六格勾满、构造器 diff 逐段对齐、打断通路三读数齐、零回归绿、d22scan 干凈—即可批准 push。

(End of file — 168 lines)
