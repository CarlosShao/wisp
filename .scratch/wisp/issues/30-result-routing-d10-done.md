# 30 — Result routing (D10): four-tier dispatch, artifacts, summaries, badge fallback

**Status:** **DONE-带残余**（2026-10-31 编排者收口：六格全翻；残余=mockllm 400 字节截断/D10 量词留痕/cmd-wisp 装配缝） ~~ready-for-agen~~
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 10-agent-loop-core, 26-tts-output
**Parallel slots:** ≤2 sub-agents (A: tier classifier + artifacts writer; B: channel emitters
+ badge/notify fallback)
**Spec refs:** SPEC-05 §7, D10 (thresholds final), D34 note②, D42#8, S4

## What to build
The output dispatcher that decides where a result goes: TTS-only, notify+panel, or file-drop —
with same-response ≤60-char spoken summaries, artifacts writes, clipboard tier-4 rule, and the
Focus-Assist-proof badge update.

## Key constraints
- Tier rules (first match wins, structured markers beat length):
  1. code block / table / ≥3-item list / ≥2 links / file-path list → artifacts file + panel +
     speak ≤60-char summary + path;
  2. >400 hanzi (or >800 chars) → file + panel + summary+path;
  3. 61–400 hanzi → notify + panel + speak first sentence ≤60;
  4. ≤60 hanzi & unstructured → TTS + ball + clipboard write.
- Summary comes from the SAME LLM response (prompt ⑥ section instructs long results to carry a
  ≤60-char spoken summary) — no extra LLM request ever (test asserts single completion).
- Artifacts writes = host-internal (no gating), data-dir only, quota/LRU via 04, path logged,
  panel-visible/deletable later (36); tier-4 clipboard ONLY (never clipboard on long results).
- Tier 3 must ALWAYS also update ball badge + task list entry (Focus Assist may swallow toast —
  D42#8: unconditional dual-channel, no API probing).
- Structured-marker detector: markdown parsing for code fences/tables/lists/links/paths —
  unit-tested with the "40-char table" case (structured wins over short).
- Streaming interplay: tiers decided on final text; Speaking begins with summary immediately
  (sentence-chunk from 26); panel stream continues if open.

## Out of scope
- Panel rendering surface (36); TTS engine internals (26); notification Focus-Assist probing
  (45 owns the "no API" conclusion).

## Acceptance criteria
- [x] Tier matrix tests: crafted outputs for every tier + boundary values (60/61, 400/401, 800)
      + structured-vs-length precedence cases.
- [x] Single-request rule: mockllm request count == 1 per result (summary present in output).
- [x] Artifacts: file written under artifacts dir, path in logs + result payload; quota respected.
- [x] Clipboard: tier 4 writes; tiers 1–3 never touch clipboard.
- [x] Tier 3 dual-channel: toast sent AND badge/task-list updated (event trace), toast-blocked
      fixture → badge still updates.
- [x] E2E: four results through the real pipeline; ball/panel/notify/audio observed per tier.

## Progress log (append-only, newest last)

## ★★★终裁落地（2026-10-31 12:1x 编排者，凭非实现者 `30-v1`，件 `docs/evidence/s1/30-v1-verdict.md` 71 行、`1f16b9a51`）

★**六格全翻、可带残余结案**（★`30-v1` 九项全自跑：present 19/0/0＋`-race`；**三发 mutation 各 rc=1 红句正中判据**；**AC2 的 mockllm 400 字节截断绕行判合法**——`chat.go:159-162` 逐字核到截断、402 汉字物理回不来，`scriptedC5` 经 `agent.Options.Provider` 注入＝AGENTS 1.3 具名缝 C5、四案同回路同 router、`present` 零 provider 引用＝**结构性无二次请求**）。
★**残余三条（随结案登记）**：①**mockllm 回声 400 字节字节级截断**（`tools/mockllm/chat.go:160`、多字节拦腰切）⇒ **归 tools/mockllm 所有者立票**；②D10 两量词解释已三处留痕（30-r1 §8⑤/测试名/裁决表 §5）；③**`cmd/wisp` 装配缝**（`present.New(Channels{...})`＋`loop.Run` 后对 completed `Dispatch`；Speak→`rt.voiceAnnounce` 闭包不重复设闸；Clipboard 适配器先 `CheckText`）⇒ **归接线腿（41a2-r1/381-r2 之后波次）**。
