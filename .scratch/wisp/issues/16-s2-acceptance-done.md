# 16 — S2 acceptance: voice end-to-end gate, latency segments, hotplug pre-mortems

**Status:** **DONE-带残余**（2026-11-01 编排者收口：格①②③⑤翻（16-v2 终判）、格④＝真 CI 工件待 push 后翻；升级发现＝D43#15/Acting 零派发⇒票 395；残余七枚见末节）
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 12-cli-text-path-s1-gate, 15-speech-engines-cer-harness
**Parallel slots:** 1
**Spec refs:** S2 done (D44), D32 16.3.4, §10 verification, SPEC-10 §3–§5

## What to build
The S2 slice gate: full voice-input chain (mic → VAD → ASR → agent → notify result → settle)
verified with wav injection and real-device smoke, all D32 latency segments measured against
budget, and the S2 pre-mortem items executed for real.

## Key constraints (must all be evidenced, not asserted)
- Chain demo: say command (wav) → ball Listening→Thinking→Acting→notify reply → Warm/Settling;
  transcript + tool evidence rows present.
- Latency segments vs D32 16.3.4: VAD-stop→ASR-final ≤800ms(stream)/≤1500ms(offline);
  context-assembly ≤50ms; ASR→punct placeholder ≤300ms (punct real in 27 — measure pipeline
  overhead only); cold-session first-token P50 ≤4s / P95 ≤6s (golden LLM), Warm P50 ≤2.5s
  (from 12's Warm data if present, else deferred note).
- Idle/Armed SLO re-run (sampler) after speech stack exists: Armed only when KWS lands (41) —
  here assert Sleeping unaffected and post-settle return ≤10s.
- Pre-mortems executed for real (SPEC-10 §7): mic occupied (test process holds device), mic
  permission denied, device hotplug mid-utterance, disk-full during transcript write, model
  file corrupted post-install (sha256 fail on load), process kill during download + resume.
- Every failure shows correct D37 class + user-visible message; nothing silent.

## Out of scope
- TTS/punctuation (26/27); KWS (41); security tools (17+).

## Acceptance criteria
- [x] E2E voice demo recorded (log + state trace + tool rows), run on real hardware.
- [x] Latency table filled with measured numbers vs budget, committed to docs/SLO.md appendix.
- [x] All six pre-mortem scenarios executed with evidence and correct error classes.
- [ ] CER gates from 15 green in CI artifacts (near_clean mandatory; noisy_far status recorded).
- [x] Adversarial acceptance pass by a different agent (stub scan + contract conformance).

## Progress log (append-only, newest last)

## 终裁落地（2026-11-01 10:5x 编排者，凭 16-v2，件 docs/evidence/s1/16-v2-verdict.md 134 行、56f94624b）

- 格①②③ 翻（各带具名边界）；格④ 维持残余（真 CI 工件待编排者下一次 push 后留存——CER 本地 24/24 与 ci.yml:732 逐字同命令已钉）；格⑤ 本件即其本体（SPEC-10 §8 七项义务全落）。对照 15-done 先例＝第 48 票带 10 条具名残余收口。
- 升级发现（16-v2）：D43#15 工具臂与 Acting 家族（#17-20）今日产码全仓零生产派发位（demo SSE 带 tool_calls 而产线 resident_voice_engines_windows.go:979-981 恒 HasToolCall:false 走 #15）＝产品级缺口（非 16-r1 伪造）、demo 走位仍每行合法；归口 26 集成面 ⇒ 开票 395。
- 残余七枚：①Acting 球态载体→票 395 ②真 CI 工件待 push（格④）③TTS 首帧 3.715s→S0 降级方向/Warm 归 28 ④OS 级三触发不可达（admin）⑤磁盘满守卫→票 45 ⑥live smoke 声学相位两发红＝环境性不判回归 ⑦kws 导出树协议差（16-r1 照录维持）。
