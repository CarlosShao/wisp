# 59 — P15 AEC spike: license, binding, echo quality, CPU, false-trigger thresholds

**Status:** **DONE-带残余**（2026-10-31 编排者收口：59-v2 有条件 PASS、两条件读数＝重开首项；残余见末节）
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 13-audio-capture
**Blocks:** 26 (Path C barge-in part), 28 (Path C full duplex), 32 (scenario ④ gate), 60 (C32)
**Parallel slots:** ≤2 sub-agents (A: license/binding/build integration; B: measurement rig +
threshold report)
**Spec refs:** D47 (§16.12), P15, C32 gate condition, SPEC-01 whitelist entry, D32 16.3.3 D47 note

## What to build
The P15 feasibility spike that gates ALL Path C full-duplex work: verify webrtc-audio-processing
(AEC3) license + Go binding viability, measure echo-cancellation quality with our self-rendered
TTS PCM as reference, measure sustained CPU, and set the barge-in false-trigger thresholds —
including the external-audio-source case.

## Key constraints
- **License first**: webrtc-audio-processing is BSD-3 — verify the exact build/distribution we
  fetch (deps.toml entry + fetch-deps, SPEC-11 §2.1); any license surprise → Path C barge-in
  degrades to hotkey/click interruption, whitelist entry retracted (D47 rule; do NOT revive
  Path T full duplex).
- Reference signal = OUR OWN rendered TTS PCM (known exactly) + WASAPI render position —
  measure alignment quality and residual echo; clock drift concern is "an order of magnitude
  smaller" than generic AEC (D47 argument) — VERIFY, don't assume.
- **Cross-endpoint case**: capture and render on different devices (Bluetooth speaker + wired
  mic) — drift returns; measure and set expectations/thresholds.
- **External-audio case (missed by everyone, mandatory)**: user plays music/video while chatting
  — third-party sound is NOT in our reference; measure barge-in false-trigger rate (VAD on
  AEC-cleaned signal) with music/speech/video fixtures; threshold = external audio must not
  trigger barge-in (target: ≤1 false barge-in per 10min playback under normal room volume);
  if unachievable, Path C keeps barge-in but requires voice-activity evidence (e.g., direction
  of AEC residual), record honest limits.
- Sustained CPU: AEC+VAD processing during playback ≤ single-core 5% PROVISIONAL (SLO table
  row) — measure and backfill the acceptance value (it's a continuous load, not transient).
- 400ms gate sanity: from speech onset detection (AEC-cleaned VAD) to TTS-stop call — the
  ≤400ms budget includes detection (~100ms-class) + stop latency; measure the chain, not just VAD.
- Deliverable: `docs/PRECHECK.md` P15 conclusions + threshold table + go/no-go for Path C
  barge-in and C32 gate condition #1.

## Out of scope
- Product AEC integration wiring (13 does the module, 26/28 consume); C32 provider work (60);
  conversation UX.

## Acceptance criteria
- [ ] License verdict recorded (exact source/build + license file committed to evidence).
- [ ] Go binding builds + links on the pinned toolchain (mingw-w64 per BUILD.md); cgo smoke for
      AEC3 passes; cross-compile status noted (mingw Linux-cross likely NOT viable for this lib —
      record; native Windows build path unaffected).
- [ ] Echo-quality rig: loop-back fixture (render TTS wav via WASAPI + capture via mic/loopback)
      → residual echo measured (ERLE-class metric or perceptual check), self-transcription
      suppression proven (zero ASR text from playback).
- [ ] Cross-endpoint drift data recorded; thresholds or restrictions documented.
- [ ] External-audio false-trigger table (music/speech/video × volumes) with recommended
      barge-in threshold config.
- [ ] Sustained CPU number backfilled into SPEC-10 §3.1 Path C row (replacing the 5% provisional).
- [ ] PRECHECK P15 = PASS/NO-GO recorded; if NO-GO → tickets 26/28 Path-C parts + 60
      auto-defer with DEFERRED registry entries (no silent disappearance).

## Progress log (append-only, newest last)

## ★★★终裁落地（2026-10-31 15:4x 编排者，凭非实现者 `59-v2`，件 `docs/evidence/s1/59-v2-verdict.md` 139 行、`3ffabfd95`）

★**PRECHECK P15 ＝ 有条件 PASS**（既非无条件 PASS、也非 NO-GO）——★**59-v2 全链自建自跑**（export 树 bc2c278bd＋junction 穿透、两 tarball sha256 双中、380 obj 逐字同、rig 三跑）。
★★★**最重发现＝AEC3 冷启动收敛双峰**：★fresh 臂 4 发 live 采捕 **30.7/29.0/8.9/8.9dB**（**n=4 仅 1/4 跨 ≥30dB 线**、钉死臂 0/1＝24.5dB）——**aec3-r1 的「两臂跨线」是 n=1 读数、支柱③未稳定翻**；塌峰态 ASR 转出两音节（「为什么」/「我是」）⇒ **单 VAD/单 ASR 面都不可靠**。★支柱①②真翻（引擎真在＋license 全核＋逐字转写 5/5 消失）。
★★**PASS 条件（重开即首格验收、非再一轮 spike）**：**(a) 长窗稳态 ERLE**（≥30s 连续采捕＝塌峰是否消失的唯一裁决读数）＋**(b) 组合 barge-in 门产品级读数**（VAD∧能量∧ASR 文本门两态误触发率）——**任一不达 ⇒ 回 NO-GO**。
★**悬案三件定性**：①VAD 单起 1/3.4s **不是 PASS 阻塞**（触发面权威定义在 **SPEC-04:137＋D47 §16.12**、不在 D42/D25）——但塌峰态证明单面不可靠 ⇒ **升为条件 (b)**；②长窗稳态＝**阻塞无条件 PASS 的唯一硬缺口＝条件 (a)**；③真跨端点原样＝残余不阻塞。
★**abseil 边界解读成立**（release wrap source_hash/patch_hash 与盘上件逐字吻合）；★license 8/8 逐字节同；★deps.toml diff 恰两节。
★**残余**：两条件读数＋真跨端点＋SPEC-10:47 AEC 终值回填（票 13/机主）。★**三条 DEFERRED＝有条件重开**（重开≠立即实现、首项＝两条件读数＋组合门设计；**登记行文字改动摆机主**）。
★**59-v2 流程缺陷自陈两笔**：派单存档时点错位（编排者已补、commit e213c2d7a 在案）／pwsh 直传 `-test.run` 拆旗标（cmd /c 留痕）。
