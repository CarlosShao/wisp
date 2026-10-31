# 41 — KWS wake word: Armed state, keywords+veto words, mute, opt-in privacy

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 15-speech-engines-cer-harness, 28-session-scope-warm
**Parallel slots:** ≤2 sub-agents (A: KWS engine + Armed lifecycle; B: veto words + L1 veto
integration + thresholds P14)
**Spec refs:** SPEC-04 §5, D2/D5/D16, B1, D43 #5,7–10, P14, S6

## What to build
`kws-zipformer-wenetspeech-3.3M` wake word engine (opt-in): Armed state lifecycle with ≤90/110MB
budget, keywords-file customization (no training), veto words sharing the model for the L1
block-window channel, one-key mute, and the P14 multi-keyword/threshold verification.

## Key constraints
- Opt-in only (`[voice] wake_word.enabled=false` default); enabling loads model (Armed);
  disabling unloads (Sleeping cap); mic stays OFF until wake or explicit listening (D16①).
- Armed budget: tree RSS ≤90MB (path Y) / ≤110MB (path X) per spike verdict; goroutine ≤7;
  audio buffers never persisted. `kws-infer` named goroutine.
- Keywords file: user-editable text list + per-word thresholds (`[voice] wake_word`); wake hit
  → Listening (KWS paused, resumes per state rules).
- **Veto words** (取消/停下/别) registered in the SAME model with per-word thresholds: hit during
  Confirming → veto (the B1 voice-cancel channel, ~100ms-class, no ASR load). P14: measure
  false-trigger rate of veto words under TTS/environment audio; too high → raise thresholds or
  disable channel with explicit UI note (decision point recorded).
- KWS-not-loaded honesty: Confirming UI shows 「语音取消不可用」 whenever KWS isn't running
  (short-cut-key sessions) — asserted e2e.
- Mute: global mute key → Muted (KWS inference stopped, model retained per #8); unmute
  restores; Armed visual per SPEC-08 §2.1.
- Wake latency: wake-word-end → Listening ≤300ms (D32) measured with fixture audio.
- **D47 interplay: KWS belongs to the Armed→Listening path (Path T semantics). In Conversation
  (Path C) the mic is dedicated to ASR+AEC — KWS stays OFF there (not "paused per state" but
  structurally unused); veto words remain the Confirming-channel mechanism for Path T only.**
- Watchdog contract honored: Armed over-limit NEVER unloads KWS (42 rule; integration test).

## Out of scope
- AEC/barge-in (deferred); fast-path voice veto without KWS (deferred B1).

## Acceptance criteria
- [ ] Armed RSS/CPU gate measured (sampler) within D32 row; wake latency ≤300ms.
- [ ] Keywords: custom keyword list change → reload without restart; hit on fixture audio.
- [ ] Veto: KWS running + Confirming + veto-word wav → L1 call cancelled; KWS absent → UI
      message asserted.
- [ ] P14 false-trigger report (thresholds table + recommended values) committed.
- [ ] Mute cycle Armed→Muted→Armed with inference stop/start asserted (no model reload).
- [ ] Opt-in off → zero mic opens (device-handle probe) + Sleeping budget.

## Progress log (append-only, newest last)

## ★★★终裁落地（2026-10-31 13:4x 编排者，凭非实现者 `41-v1`，件 `docs/evidence/s1/41-v1-verdict.md` 72 行、`e926531d5`）

★**六格各有缺角、按缺角保持未勾 ⇒ ★可带残余结案**（先例 15〔收 3 格未勾〕／28〔带残余 1 格〕／42-done）。
★**六格缺角逐格**：AC1 半格（缺 42 面 sampler 实测＋≤300ms 实测）；AC2 半格（hit on fixture audio 已验；**reload 缺＝`keywords` 键 config 层有热更分类但 KWS 路径零消费、thresholds/veto_words 255 审计零命中→config 票/B 腿**）；AC3 结构就位 **e2e 无从做（`kwsVetoWords` 空表 `kws.go:83`→B 腿改表即激活**；UI 半句→Confirming 票）；AC4 不覆盖（P14 报告未落、负控是原料）；AC5 半格（机制＋guard 双分支已验；缺真引擎管线级「Muted 里 kws_fed 冻结」断言）；AC6 半格（plan 级 opt-in 门已验；缺设备句柄探针＋Sleeping 预算→采集/42 面）。
★**41a2-r1 自评与复判逐格相符、无纸糊**；★**复跑三红全 staging 缺口＝零代码红**（WISP_SPIKE_MODELS 未设×2／System32 onnxruntime 1.17 遮蔽〔正解＝DLL 贴测试二进制、PATH 注入无效〕／VAD 缺失 plan 早退）。
★**残余 7 条归口**：①sampler→票 42 ②词表热更→config 票+B 腿 ③veto e2e→B 腿+Confirming 票 ④唤醒句话轮过滤→**票面/D43 未定义、先定义再接 consume() 一处（未自创＝正确形状）** ⑤**#32 nil facts→statemachine 冻结面钩＝须人工批准（摆机主）**、现行 **#5 重武装自愈已 e2e 验证兜住** ⑥models ensure 无 KWS 入口→models 票 ⑦Speaking 双门→26 家族观察项。
★★★**R28-2 注销随本判**：★**`warmFeed` 真喂食门已落地并 e2e 验证**（Warm 静音帧 ⇒ `kws_fed`+20480/`asr_fed`+0；命中 ⇒ 门重开 `loads_after_warm=0`）——**「mic 是闩不是门」问题正式注销**。
