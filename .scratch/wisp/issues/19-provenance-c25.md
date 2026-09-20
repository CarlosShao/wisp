# 19 — Provenance (C25): taint marking, ≥8-char leak matching, sync-dir detection

**Status:** review
**Claimed by:** agent-ticket19-taint
**Last update:** 2026-09-20T03:05Z
**Blocked by:** 17-risk-assessor-c19
**Parallel slots:** ≤2 sub-agents (A: taint store + matcher; B: sync-dir detection P12 + four-
channel exfil tests)
**Spec refs:** SPEC-06 §5, D30①, D33/F4, C25, P12, 16.9#1

## What to build
Coarse-grained provenance tracking: sensitive-source taint marking, the six exfiltration
channels, normalized ≥8-char contiguous-fragment matching, and sync-directory detection so that
"read something sensitive then send it anywhere" upgrades to L2 with the source named.

## Key constraints
- Sensitive sources marked on tool output: `fs.read` (incl. out-of-allowlist attempts),
  `search.content` results, `clipboard.read`, `sysinfo` focused-window title, `web.fetch`
  responses, `doc.read` content, `screen.capture` images, user transcript (from 15).
- Exfil channels (all six): `web.search` query string, `notify` text+URL, TTS announce text
  (physical exfil), `clipboard.write`, `fs.write` into sync dirs, HTTP POST/URL body (net).
- Match: normalized (whitespace-stripped, case-folded, full/half-width unified) ≥8-char
  contiguous fragment of a tainted source appearing in any exfil parameter → R4 upgrade to L2;
  decision names the source ("包含来自 <源> 的内容"). Deliberately coarse (16.9#1): false-
  positive direction is safe; no token-level taint tracking (rejected — do not implement).
- Sync-dir detection (P12): registry + client-config probing for OneDrive/Dropbox/Nutstore/
  Google Drive known and configured locations; path-string matching alone is insufficient;
  undetectable → treat `[fs] allowed_dirs` entries under user-profile as sync-suspect (safe
  default). Conclusion recorded in PRECHECK; write into sync dir counts as exfil channel
  regardless of taint when content matches.
- Taint lifetime: bound to task/session scope (DisposalScope); tool_call rows log taint hits.
- Four-channel red-team tests: search-content marker → exfil via web.search query / notify /
  clipboard.write / sync-dir write → all upgrade L2 with named source.

## Out of scope
- Approval rendering (21/37); actual tool implementations (20/22).

## Acceptance criteria
- [x] Four-channel exfil suite: each channel upgrades to L2 with correct source attribution.
      `TestFourChannelExfilSuite` (web.search query / notify text+URL / clipboard.write /
      fs.write→sync-dir) + TTS & HTTP-body via `CheckText`; end-to-end through the real C19
      assessor in `TestR4EndToEndViaAssessor` (L2, [R4], SessionOverrideBlocked, reason
      「包含来自 doc.read …」).
- [x] Normalization tests: whitespace/case/full-width transforms still match; <8-char does not.
      `TestNormalizeTaintContractForm` + `TestFragmentMatch*` (exact 8/7 boundary, CJK, CJK
      full-width digits; hash-collision verification).
- [x] LLM-paraphrase tolerance test (documented limitation): `TestLLMParaphraseResidual`
      locks the miss as accepted residual (D30 layers 2/3/5; 16.9#1) — NOT a fix.
- [x] Sync-dir detection on a machine with OneDrive configured (self-hosted evidence) + fixture
      fallback; registry probe code reviewed against P12.
      `TestSyncDetectionOnThisMachine` on C:\Users\swq\OneDrive (evidence:
      docs/evidence/s3/19-provenance-c25-tests.txt); fixtures `TestSyncFixtureFallbackAndMatch`
      / `TestSyncDropboxHostDBConfig`; conclusion backfilled docs/PRECHECK.md §P12. Registry
      probe code awaits adversarial review (this ticket is Status: review).
- [x] Taint scoping: new session does not inherit old taints (DisposalScope test).
      `TestScopesNeverInherit` + `TestDisposalScopeClearsTaints` (real plugin.DisposalScope
      Defer→Dispose lifecycle) + `TestInspectUnknownScopeIsEmptyStore`.

## Progress log (append-only, newest last)
- [2026-09-20T02:20Z] agent=agent-ticket19-taint did=claimed next=engine+skeleton-tests
- [2026-09-20T02:55Z] agent=agent-ticket19-taint did=C25-engine (internal/risk/provenance.go taint store+channel Inspect+CheckText+Detector(scope) bound to frozen TaintHit seam; taintmatch.go normalize+>=8-rune window index w/ collision-verified hashing; syncdirs.go+_windows.go P12 registry/config/default/fixture probes + fail-closed suspect fallback; 21 new tests all green incl. four-channel suite through real C19 assessor, DisposalScope non-inheritance, threshold clamp, paraphrase residual; risk vet+count=2+-race+full d22scan clean; note: go build ./... currently broken in internal/agent/ (ticket-10 untracked WIP, left untouched)) next=PRECHECK-P12+evidence+status-review
- [2026-09-20T03:05Z] agent=agent-ticket19-taint did=P12-conclusion-backfilled(docs/PRECHECK.md §P12, incl. nutstore/google-drive encrypted-config residuals + suspect-fallback rule) + self-hosted evidence(docs/evidence/s3/19-provenance-c25-tests.txt: live OneDrive C:\Users\swq\OneDrive detected via default-location probe, registry Accounts key present without UserFolder) + DEFERRED(C25-loop-wiring) marker in provenance.go header (Mark/Inspect/OpenScope/CloseScope call-sites belong to tickets 10/20/21/22/26; R4 reason flows into existing tool_call DAO, Hit.Fragment never persisted) next=verifier (Status=review)
