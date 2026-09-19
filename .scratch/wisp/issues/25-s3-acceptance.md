# 25 — S3 acceptance: capability scenarios ①② + full security red-team suite

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 21-approval-gates-minimal, 22-web-tools-d30, 23-system-window-input-tools,
24-doc-search-tools
**Parallel slots:** ≤2 sub-agents (A: scenario scripts ①②; B: red-team suite execution)
**Spec refs:** §10 #5, #10, #11 (S3 rows), D44(b) matrix, SPEC-10 §6

## What to build
The S3 slice gate: scripted execution of user scenarios ① (system ops & file governance, 10
instructions) and ② (info & research, 5 instructions) with per-instruction evidence, plus the
complete security red-team suite green.

## Key constraints
- Scenario ① ten instructions (each: exact instruction text, tool-call trace, risk-level match
  against D34, screenshot/log evidence): open notepad; group PDFs in Downloads into one folder;
  volume 30%; close this window; find last week's contract; move file to recycle bin (assert
  trash=L1 not L2); lock screen; switch to browser; battery level; open Bluetooth settings page.
- Scenario ② five: check latest version of X; open link and summarize; read my screen error
  (assert screenshot entered LLM as ImagePart); read PDF chapter 3; search this error code.
- Red-team suite (all must pass; any failure = slice NOT done): out-of-allowlist access;
  undeclared capability; declared-L0 dangerous tool still judged L2 (C19); web.fetch private
  range deny; stopReason=length fails all tool calls; L2 without confirm never executes;
  net out-of-allowlist domain → L2; C26 four-bypass real junctions; A-list reads (git-credentials,
  .git/config, wisp config.toml) denied over any grant; C25 four exfil channels → L2 with source;
  C29 one-byte tamper reject; F2 forged panel-allow rejected; D45 batch boundaries; config
  loosening → L2 re-confirm.
- Single-task constraint still in force (TaskScheduler width=1) — concurrency tests deferred
  to 47/54 (recorded as explicit NOT-tested-here).
- Gap audit + pre-mortem (SPEC-10 §7) + adversarial acceptance (different agent) + scenario↔
  slice bidirectional matrix check.

## Out of scope
- Voice/TTS; panel UI; concurrency; session grants; KWS.

## Acceptance criteria
- [ ] 15 scenario instructions scripted and green with evidence artifacts committed (testdata/
      evidence/s3/).
- [ ] Full red-team suite green in CI (windows job for junction-real tests).
- [ ] Pre-mortem report + adversarial acceptance report committed (docs/evidence/s3/).
- [ ] DEFERRED cross-check: all markers ↔ registry consistent for touched code.
- [ ] RiskLevel conformance: no instruction's decision differs from D34 table (log diff = fail).

## Progress log (append-only, newest last)
