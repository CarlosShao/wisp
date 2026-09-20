# 19 — Provenance (C25): taint marking, ≥8-char leak matching, sync-dir detection

**Status:** review
**Claimed by:** agent-ticket19-fix
**Last update:** 2026-09-20T04:35Z
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
- [2026-09-20T03:50Z] agent=agent-ticket19-fix did=accepted FAIL verdict, starting must-fix set B-1/B-2/M-1/M-3/M-6 next=red-team-tests-then-gates
- [2026-09-20T04:35Z] agent=agent-ticket19-fix did=per-defect disposition of the FAIL verdict — **B-1 FIXED** (syncdirs.go `resolveTarget`/`deepestExistingAncestor`/`splitPathComponents`: 目标不存在时先用 C26 Resolve 锚定「最近已存在祖先」再接词汇尾巴，只有连祖先都无法核实才落 sync-suspect；未采用判定书 §4 的一行式「!res.Resolved 即从严」，因为 fs.write 目标天生 Resolved=false，那样会把每次正常新建写入推去 L2)：钉 `TestSyncRedTeamNewFileSpellingsDeniedAndLandNoBytes`(4 subtest，8.3 别名由 GetShortPathNameW 取回 ONEDRI~1，非手抄) + `TestSyncRedTeamRealOneDrive`(真机 `\?\C:\Users\swq\OneDrive\...` → sync=true root=OneDrive/env，R4 hit，门控写盘后 0 字节落地) + 反向护栏 `TestSyncRedTeamGuardPlainNewFileWriteAllowed`/`TestSyncNormalNewFileWriteNotFlagged`(普通新文件写不判同步且字节确实落地) + `TestSyncRedTeamUnverifiableChainFailClosed`(Q盘/Foreign UNC 整链不存在→从严)。变异复核：把 resolveTarget 改回旧词汇回退后上述红队用例全红(EXPLOITABLE)，改回修复后全绿。 **B-2 FIXED** (provenance.go Inspect：channelTable 只贴 Channel 标签，其余参数一律 collectStrings 以 ChUnknown 追加扫描，写载荷仍受同步门约束；collectStrings 另收 []byte)：钉 `TestNamedChannelParamEscape`(10 例含 url/q/title/payload/嵌套 body/数组/[]byte/深层 map + 冻结缝复测)。 **M-1 FIXED** (provenance.go `scopeMarks`：scope 未注册且引擎任意处存在 taint → Hit{ChUnknown, SrcUnboundScope} 从严；`TestInspectUnknownScopeIsEmptyStore` 已改为断此方向，并覆盖 CheckText 与 Detector 缝)。 **M-2 FIXED** (探针链首加 `%OneDrive%/%OneDriveConsumer%/%OneDrivePersonal%/%OneDriveCommercial%/%OneDrivePublic%` → Source="env"；注册表探针逻辑下沉为可注入的 `registryProbeFor(registryValueSource)`)：钉 `TestSyncRegistryProbeTable`(4 subtest) + `TestSyncEnvConfiguredRoots` + `TestSyncRegistryProbeLive`(真机 0 条 → t.Skip 并写明原因) + 真机证据现记 source=env。 **M-3 FIXED** (syncdirs.go `gradeConfirmed` + `finalize`：只有 env/registry/config 级且 C26-canonical 的根才撤销 profile 兜底，default/fixture/options/未证实根只增加根；顺带吃掉 **N-2**——`canonical` 字段从此被消费)：钉 `TestSyncFallbackNotDisarmableByWeakRoot` + `TestSyncUnverifiedRootKeepsFallback`，`TestSyncFixtureFallbackAndMatch` 改断新方向。 **M-4 FIXED** (MaxParamDepth 默认 8，超深=截断 → Hit{SrcUnscannedNesting} + 日志)：钉 `TestDeepNestingFailsClosedNotSilent`。 **M-5 FIXED** (删死字段 winStrs + 不可达碰撞探邻居分支 + hashWindowPairs；contains 不再每候选 string(norm) 全量拷贝)：同进程 old/new 对照实测——单满额源常驻 **10.1→1.8 MiB**、五源 **50.5→8.9 MiB**；建索引 **38.3→22.9 ms、10.6→2.8 MiB/op、233014→5 allocs/op**；响应路径扫描 40k 参数×20 源 **55.8~58.9ms vs old 57.4~57.8ms（换序复跑，差异在噪声内，未宣称改进）**；Mark 满额源 27.8ms/4.3MiB；新增可复跑基准 `BenchmarkMarkFullSource`/`BenchmarkScanOverTenSources`，口径写入 docs/PRECHECK.md「C25 片段索引的 D32 预算」。 **M-6 FIXED** (normalizeTaint 丢弃 U+00AD/200B/200C/200D/2060/FEFF + unicode.Mn)：钉 `TestNormalizeTaintDropsIgnorableChars` + `TestNormalizationEndToEnd`(含 Mark→Inspect 全链路隐形字符/全角/空白变形，顺带补 **N-4**)。 **N-1 FIXED**(删除死字段 ProvOptions.AllowedDirs，规则明确为「profile 下任意路径，严于 allowed_dirs」，PRECHECK 同步)。 **N-3 FIXED**(死分支删除；`TestFragmentHashCollisionCannotFakeHit` 改为注入 ghost 哈希制造真碰撞)。 **N-5 FIXED**(suspect 兜底改组件边界 isUnder，`TestSyncSuspectFallbackIsComponentBounded` 钉 profileevil 不受牵连)。 **N-6 FIXED**(票 26 补 TTS 必须过 CheckText(ChTTS) 的 AC 行；票 55 补 macOS 云盘探测判据 AC 行，syncdirs_other.go 的 DEFERRED 头同步点名)。 **N-7 更正**：上一轮日志「21 new tests」实为 30 个测试函数；本轮修复后 68 个顶层测试 + 38 个 subtest，全绿 1 skip(仅真机注册表探针)。 未改动冻结面：assessor.go / pathresolver.go / docs/PLAN.md / docs/specs/* 零 diff（git diff --stat 空）。gates: vet CLEAN, -count=2 ok, -race ok, d22scan clean（证据 docs/evidence/s3/19-adversarial-fix-round.txt）。 next=交独立复验（Status 回到 review；本代理不自证通过）
