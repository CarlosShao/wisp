# Wisp Ticket Index

Source: `docs/specs/SPEC-*.md` (authoritative for all numbers/rules cited here).
Slice order: S0→S1→{S2∥S3}→S4→S5→S6→S7→S8 (SPEC-12 §2).
All tickets are **vertical tracer bullets**; work the **frontier** (all blockers done).

## Status lifecycle (anti-interruption rules — MANDATORY)

| Status | Meaning |
|---|---|
| `ready-for-agent` | Frontier ticket, unclaimed |
| `in-progress` | Claimed; must have ≥1 Progress-log entry |
| `blocked` | Waiting on external decision (ticket notes which) |
| `review` | Work done, awaiting adversarial acceptance (SPEC-10 §8) |
| `done` | All boxes checked; **rename file with `-done` suffix** + title `(DONE ✅)` |

**Rules (prevent abandoned/in-progress tech debt):**
1. Before writing any code: set `Status: in-progress`, fill `Claimed by`, `Last update`,
   append a Progress-log line → **commit + push (both remotes)** in the same step.
2. Append a Progress-log line after every meaningful unit (`- [UTC ts] agent=<id> did=… next=…`),
   commit+push. Sub-agent lines tagged with their id. **Max 2 sub-agents per ticket**;
   recommended overall fleet width is also 2.
3. An `in-progress` ticket with no log update for >24h: next agent must verify on-disk state
   against the log before continuing; never reset checked boxes; append `resumed` line.
4. Completion: check all acceptance boxes → `Status: done` → rename file `NN-slug.md` →
   `NN-slug-done.md` → update this index → commit+push.
5. If a decision is missing, set `Status: blocked` with the open question in the log
   (D22 闸门③: undefined = stop and ask, never assume).

## Dependency graph (blockers in parentheses)

- **S0**: 01 build-chain ✅done · 02 s0-spike ✅done (01)
- **S1**: 03 skeleton ✅done (01) · 04 sqlite-core ✅done (03) · 05 config-model ✅done (03) · 06 secretstore-envs ✅done (03) ·
  07 ball-state-machine ✅done (03) · 08 observability-slo ✅done (03,07) · 09 llm-provider+mockllm ✅done (03) ·
  10 agent-loop-core (05,09) · 11 llm-adapters-rest (09) · 12 cli-text-path=S1 gate (04,07,08,10)
- **S2**: 13 audio-capture ✅done (03) · 14 model-distribution ✅done (03) · 15 speech-engines+CER (02,13,14) ·
  16 s2-acceptance (12,15)
- **S3**: 17 risk-assessor C19 (03) · 18 path-resolver C26 (03) · 19 provenance C25 (17) ·
  20 host-bridge+fs-tools (17,18) · 21 approval-gates-minimal (17) · 22 web-tools+D30 (19,20) ·
  23 system/window/input-tools (20) · 24 doc/search-tools (20) · 25 s3-acceptance (21,22,23,24)
- **S4**: 26 tts-output (15) · 27 punctuation (15) · 28 session-scope Warm (07,15) ·
  29 memory-l1/l2 (04,10,28) · 30 result-routing D10 (10,26) · 31 reminders (04,30) ·
  32 s4-acceptance scenarios ③④ (27,29,30,31) — **26/28/32 的 Path C 部分（D47 全双工/barge-in）
  另被 59 阻塞**
- **S5**: 33 panel-host C27 (07,12) · 34 frontend-scaffold (01) · 35 panel-bridge C17 (33,34) ·
  36 result-history-panel (30,35) · 37 approval-ui (21,35) · 38 palette+tasks (35) ·
  39 config-editor-gui (05,35) · 40 security/privacy/cost-pages (29,35,44)
- **S6**: 41 kws-wake-word (15,28) · 42 watchdog (08,28) · 43 power-lifecycle (28) ·
  44 costmeter C23 (04,10) · 45 diagnostics+guards (42) · 46 s6-acceptance (41–45)
- **S7**: 47 task-scheduler+pathlock (10,21) · 48 approval-queue-full (37,47) ·
  49 session-grants D45-2 (04,48) · 50 tier1-manifest-plugins (20) · 51 tier2-goja (17,20) ·
  52 d46-command-plugins (50) · 53 larkcli-plugin-e2e (48,51,52) · 54 s7-acceptance (47–53)
- **S8 (deferred, not ready-for-agent)**: 55 macos · 56 signing+dist+naming · 57 plugin-sdk+registry · 58 i18n
- **D47 追加（2026-09-19，编号晚于 58 但入图拓扑早于被阻塞票）**：
  **59 p15-aec-spike (13) → 阻塞 26/28/32 的 Path C 部分与 60**；
  **60 c32-realtime-engine (16,28,44,59；门控双条件：59 通过 + 用户有 Key，缺一自动推迟登记)**。
  ⚠ 编号惯例说明：后补票用更大编号以保持文件名稳定，忽略编号与拓扑序的差异，以 Blocked by 为准。
- **LLM 接入补充（2026-09-19 用户批准）**：
  **61 cloud-voice-providers-C9 (16,11,05)** —— 云级联 ASR/TTS（语音模型下拉的落地层，
  候选 StepAudio 2.5 ASR / StepAudio 3 TTS）；配套修订：票 05/09/11/40/44/60 均已注入
  目录/能力位/计费模式/配额三层/兜底链/角色化默认/probe/限流要求；
  SPEC-02（provider_health v2 表）、SPEC-03 §3.1、SPEC-05 §3.3a、SPEC-12 登记表（云端语音
  DEFERRED→已采纳·提前）同步。

⚠ **Safety-incomplete period**: tickets 21–32 run with `TaskScheduler` single-task only;
multi-task concurrency unlocks only at 47 (SPEC-12 §2).

## Hard global constraints (apply to EVERY ticket)

- Never modify D1–D46 / C1–C31 / R1–R9 / D43 transition table (contract change = human approval).
- Forbidden patterns (auto-checked in CI): bare `go func(` without owner/recover; `filepath.Clean|Abs`
  for fs decisions outside `risk.PathResolver`; plaintext API keys; wall-clock-difference timeouts;
  hashes fetched from mirrors; panel-sourced L2 "allow"; host-internal artifact writes implemented
  as gated Tools; any emoji in UI code (U+2190–U+2BFF, U+1F300–U+1FAFF, U+FE0F).
- `DEFERRED(D-xx)` code markers must map 1:1 to SPEC-12 §5 registry (bidirectional check).
- Tests inject at seams only: C8 AudioSource (wav), C5 LlmProvider (golden SSE), C17 PanelBridge,
  CLI `wisp run`. No mock-instead-of-real to fake completion; never weaken SLO thresholds.
- Each slice completion requires: gap audit + pre-mortem + adversarial acceptance by a DIFFERENT
  agent (SPEC-10 §7–8).
