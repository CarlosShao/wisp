# 26 — TTS output: sherpa matcha engine, serial slot, Speaking pipeline, P7 gate

**Status:** **DONE-带残余（2026-10-27 20:2x 编排者收口：4 格已翻（AC2 带两注/AC4 带冷档注）；AC1/AC3/AC7 具名归口＋R26-7 新抓生产缺陷＝7 条残余见末节；承接票由编排者下一轮开）** ~~ready-for-agent~~
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 15-speech-engines-cer-harness
**Parallel slots:** ≤2 sub-agents (A: TTS engine + playback; B: half-duplex integration +
latency/P7 harness)
**Spec refs:** SPEC-04 §4, §8, D5, D16 half-duplex, D32 16.3.3/16.3.4, P7, S4

## What to build
`TtsEngine` (C9) on sherpa matcha/vits zh (int8): synthesize-and-play for short results, serial
engine-slot residency with ASR, half-duplex mic gating during playback, interrupt handling, and
the P7 subjective-quality gate harness.

## Key constraints
- TTS slot via the 15 engine-slot mutex: **Spike backfill (T02): serial switch measured P50 4726ms > 1.6s budget → D32 16.3.3 preset
      degradation ACTIVE: TTS RESIDENT + ASR on-demand is the default Path T policy (S4 may
      re-evaluate serial after tuning).** peak = max(ASR,TTS) not sum; unload via DisposalScope;
  measured switch latency vs "first-token→TTS-first-frame ≤800ms (loaded) / ≤1600ms (load)".
- Half-duplex (D16, **Path T only per D47**): during playback, mic capture closed AND KWS
  inference suspended (except veto-word channel rule in 41); playback end → Speaking→Warm
  transition (mic stays closed, default policy). **Path C (Conversation) is the D47 exception:
  mic STAYS OPEN with AEC (ticket 59/P15-gated), user speech mid-playback = barge-in → stop TTS
  within ≤400ms and never transcribe our own playback audio.**
- Interrupt channels during Speaking: hotkey / ball click / veto word (when KWS running) →
  stop synthesis+playback, release audio out, → Listening (D43 #27).
- Streaming: speak first sentence as soon as available (sentence-chunk synthesis) to meet the
  800ms budget; long text speaks summary only (routing at 30).
- Audio out: WASAPI render; volume follows system; no persistence of synthesized audio.
- P7 gate harness: fixed 20-sentence Chinese set; double-blind human scoring ≥7/10 recorded in
  evidence; below → fallback decision point (cloud TTS provider slot stays RESERVED, not
  implemented).
- Failure: TTS engine load fail → `model`/`resource` error class; playback device fail →
  `audio_device`; text-only fallback path (notify) still delivers result.

## Out of scope
- Result routing tiers (30); Conversation-mode mic policy internals (28, Path C owner); punctuation
  (27); AEC implementation (13 + spike 59 — this ticket consumes its gate for Path C barge-in).

## Acceptance criteria
- [ ] Speak path: short result → TTS audible, ball Speaking with success breathing; interrupt
      mid-playback → instant stop + Listening.
- [x] Half-duplex (Path T): playback injects no ASR transcripts (self-echo fixture); capture gate
      closed during playback (13 counter asserts zero frames).
- [ ] **Path C barge-in (D47, blocked by 59): with AEC on, speak mid-playback → TTS stops ≤400ms →
      Listening; loop back our own TTS audio into the mic input → zero ASR transcripts (no
      self-excitation); P15 numbers recorded (CPU sustained, false-trigger incl. external audio).**
- [x] Latency: first-audio ≤800ms warm / ≤1600ms cold-load measured (10 runs, P50/P95 recorded).
- [x] Serial slot: ASR and TTS never co-resident (sampler peak = max, asserted); unload frees.
- [x] P7: score sheet committed ≥7/10 or blocked-decision logged for cloud-TTS fallback.
- [ ] C25 wiring (registered by ticket 19's DEFERRED(C25-loop-wiring), adversarial report
      N-6): every announced string goes through `risk.Provenance.CheckText(scope,
      risk.ChTTS, text)` before playback; a hit means the announce is replaced by the L2
      confirmation (source named) and never spoken aloud. Test: the TTS channel assertion in
      `internal/risk/provenance_test.go` (TestFourChannelExfilSuite) plus a 26-side gate test.

## Progress log (append-only, newest last)

## ★★★终裁判语（2026-10-27 20:1x 编排者，凭非实现者 `26-v1`，件 `docs/evidence/s1/26-v1-verdict.md` **83 行**、`17c815505`）

★**七格 4 翻（`AC2` 新入列）／3 格具名残余** ⇒ ★**能带残余结案**（★先例：`15-done` 3 格未勾、`347-done` `AC#0` 未勾）。

### ★★★① `AC2` ✅ 翻（**带两注、⛔ 静默吞票面原文**）
★**13-counter 半句以门态读数等价收口** —— ★票面"13 counter zero frames"真身 ＝ `internal/audio/gate_test.go:47`（★席位在 13、**设备级零帧归生产装配态**）
★**针在盘上**：★`:912 SetSpeaking(true)`（HOP 2b、换族后首样本前）＋★`:903 false`（`gateRelease` defer、★**六出口恰一次**：闸拦/无引擎×2/无设备/合成失败/drain 失败/成功）
★**self-echo 夹具复绿**：★`GATE-26R3 closed=1 reopened=1 mid_speech_transcripts=0 spoken_samples=89389`／★`ECHO-26R3 frames_fed=6 transcripts_after=1`／★`REOPEN-26R3 finals=2`
★★★**`tri-v1` 那条「ring 重放带播放期回声（机理未测）」现判**结构性消除****：★**门闭 ⇒ inner 源 `Stop()`（`gate.go:12-16`＋`gate_test.go:47` 钉死帧计数冻结）⇒ 播放期零新帧可进 ring ⇒ 下个 epoch `Accept(ring)` 无回声可重放**
★**注 2**：★**夹具喂的是静音帧、真回声自激面随 Path C 归 59**

### ★★② `AC4` ✅ 翻（★**冷档 NOT MET 必带注**）
★**冷档本轮实测 `4.0550224s`**（★wire `4.29`／tri `4.82`）⇒ ★**对 1600ms 全线 `2.5–3.8×`**（`R26-4`：归换族/驻留策略票）

### ★★③ `AC5`/`AC6` ✅ 翻
★`AC5`：★tri `0.41<0.75`、本轮 speech 套件含该案全绿；★`AC6`：★`R5` 第二支已入账、★**P7 表不存在（`R26-6` 不挡结案）**

### ★★④ 三格不勾 ＝ 具名归口
- ★**`AC1` 不勾 ＝ `R26-1`**：★**生产打断通道 0**（`EvInterrupt` 仍只 `:940/:970` 两条失败分支、≤400ms 无读数）
- ★**`AC3` 不勾 ＝ 归票 59**（★`R26-3`：59 现跑 `ready-for-agent` 未 done）
- ★★**`AC7` 不勾 ＝ `R26-5`**：★**调用点半边已收**（★`speak():864` ＋ wire-M1 牙 ＋ `TestWireR1AnnounceGateRefusesTaintedText` 绿）；★★**但生产 `.Mark(` 全仓仍只 `panel_assets.go:245` 一处、`voice-announce` 的 Mark 唯一写者仍是测试**、★**命中后只 `noteProblem` 无 L2 确认**

### ★★★⑤ `R26-7`（★`26-v1` 新抓的生产缺陷）
> ★★**`buildResidentVoice` 补了 queue 但**仍缺 `provScope`/`newSpeaker` 两字段**（★只在测试构造器 `:427-428`；★核过 wire-r1 原版 `d803344ae` 与 26r3 终版 `06609feaf` 两版 blob 都如此）
> ★**`provScope` 有 `announce():798-801` 空串回落兜住**；★★**`newSpeaker=nil` ⇒ 生产出声 HOP 3/4（`:935`）**`nil` 函数值调用 panic**** —— ★**registry `recover` 边界兜进程不塌、`defer` 仍放门、但**该次播报死****
> ★**修法 ＝ 生产字面量补两字段或两构造器合一**（★**建议与 `R26-1` 同一张承接票、同在 `cmd/wisp` 装配面**）
⇒ ★★★**这是"测试构造器 vs 生产构造器是两条路"的**第三座碑****（★**`26-r3` 修了 queue、`26-v1` 又抓两字段**）

### ★★⑥ 反向判据（★`26-v1` 自造、导出树＋`-overlay`）
★**M1（删 `SetSpeaking(true)`）⇒ `rc 1`**：★逐字 `…(closed=0 reopened=1): SetSpeaking(true) is not on the path, so the D16 half-duplex close never happens and the microphone would hear the announcement`
★★**M2（删 `SetSpeaking(false)`）⇒ `rc 1`**：★**逐字 `…(closed=1 reopened=0 speaking=true): SetSpeaking(false) is not on the path; capture stays closed - stuck in Speaking - and no later utterance can ever be captured again`**
★**还原核（时机＝M2 后）**：★无 overlay 重编 **`Test26R3` 3/3 PASS**；★导出树 status 0 行/`HEAD` 同值；★`MUTANT 26r3` 4 处全在测试件；★**工作树零脏**
★**⚠ 口径**：★**M1 首发红是 cwd 错**（★`voice15eRepoFile` 从 cwd 找仓根）—— ★**改 cwd=仓根才得真红**（★**清单第 3 例**）

### ★★⑦ 正控／零回归／门禁
★**正控**：★**120 帧静音 ⇒ `SILENCE-26R3 finals=0 announce_asked=0 gate_closed=0 gate_reopened=0`（门针零触发）**
★**零回归**：★**`TestWireR1|Test15e|Test26R3` ＝ `16/16 PASS`**、★**`262/640/902` 逐位同值**、★**`drained=169/181`（固有抖动族）**；★**`internal/speech` staged 带 punc ＝ `67/0/0`**；★session/audio/vet 全绿；★**`cmd/wisp` 定向 `rc 0`**
★★**`d22scan`（导出树）**：★**干净 `rc 0`**（分母 **264/45/684/159**、与 `26r3` 终版逐位同）；★**正控 A（`U+2713`）/B（裸 `go func`）各 `rc 1`、删后 `rc 0` 分母回落 ＝ 可逆**（★**第 21 处定式**）

### ★★★⑧ 残余清单（★Progress log 追加、逐条）
★`R26-1` 打断通道（**承接票待开**）｜★`R26-3` **归 59**｜★`R26-4` 冷档 NOT MET（归换族/驻留策略票）｜★`R26-5` 毒源＋L2 确认（**承接票待开**）｜★`R26-6` P7（**`R5` 已收口**）｜★★**`R26-7` 生产构造器补 `provScope`/`newSpeaker`（建议与 `R26-1` 同一张承接票）**
★**`R26-2` 收口**（两注在案）
★**⚠ `tri-v1` 要求 "AC7/AC1 各开承接票" —— ★**工单池现跑 0 张 R26 票**、★**开票是编排者动作**（★**终裁腿未越权代开、已具名上报**）⇒ ★**编排者接手：下一轮开两张承接票**

### ★★★⑨ 结案判语
★**能带残余结案**（★**`tri-v1` 最小可结条件两件中第一件——`SetSpeaking` 一针＋self-echo 夹具——已做完且 `26-v1` 全链复验**；★**第二件"各开承接票"＝ 编排者动作、本轮执行**）

> ★**R26-1「打断」半格收口（2026-10-31 记）**：★**票 381（381-r2/381-v3）落地 ≤400ms MET（三跑 15 轮全 ≤59ms、真设备真 matcha、生产 TriggerInterrupt 通路、D43#27）** ⇒ **`R26-1` 打断半格具备收口读数并记收口**（★**R26-7 生产构造器同票收口**；★**残余：球单击未装配建议独立票／Drain 早回归 26 侧核**）。凭据＝`docs/evidence/s1/381-v3-verdict.md`。
