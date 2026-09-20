# 14 — Model distribution: C29 signed manifest, mirrors, resume, Downloading state (DONE ✅)

**Status:** done
**Claimed by:** orchestrator -> sub-agent T14-impl
**Last update:** 2026-09-19T23:09:31ZT15:11:33Z
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: manifest + minisign verify + downloader; B: state wiring +
local_override + mirror tooling/compose)
**Spec refs:** SPEC-04 §7, D26, C29, D33/F3, D42 model rows, P3/P5, S2 (P5 blocks S2)

## What to build
`models` module: signed model manifest, verified multi-mirror downloads with resume/retry/
progress/cancel, `Downloading` state wiring, local-dir escape hatch, and the compose
`model-mirror` service (normal/corrupt/bad-signature fixtures) used by tests and dev env.

## Key constraints
- `models/manifest.json` in-repo: per model `{id, purpose, urls[], sha256, size_bytes, license,
  quant}`; signed with **minisign**; public key hardcoded in `buildinfo` (C29). Verification is
  offline (manifest signature checked before any network use).
- Download flow: verify manifest signature → fetch bytes from mirror list (domestic mirrors
  first, official fallback; `[models] mirror` config) → verify file sha256 → install to
  `models\<id>\`. Mirror provides BYTES ONLY; hashes come exclusively from the signed in-repo
  manifest (F3: same-source hash = zero auth value).
- HTTP Range resume; retry with backoff; progress events (ball `Downloading` ring + percent);
  cancel cleans partial files to `staging\`.
- `[models] verify_signature=false` must ERROR, never disable (config already enforces; runtime
  double-checks).
- `[models] local_override{id:path}`: skip download, run integrity sha256 anyway (cheap), then use.
- P3: fill `license` field per model from license verification; any non-commercial finding = a
  blocked-decision log, model not shipped.
- Model set for this ticket: KWS zipformer 3.3M (bundled candidate — note: bundling decision
  recorded; default = download-on-first-use unless spike says bundle), silero VAD, streaming
  Paraformer int8, offline SenseVoice int8, CT-Punc int8, matcha-zh int8. Actual binaries pulled
  by 15; this ticket ships the pipeline + registry entries.
- Compose `model-mirror` (port 18081): serves good files, a corrupted file, a tampered manifest.

## Out of scope
- Engine loading/inference (15); GUI model management page (40); CPU-heavy CER (15/16).

## Acceptance criteria
- [x] Tamper tests: flip one model byte → sha256 fail → reject+delete; tampered manifest →
      signature fail before download; `verify_signature=false` in config → hard error.
- [x] Resume test: kill downloader mid-transfer → restart continues (byte-offset verified),
      completes, hash passes.
- [x] Mirror failover: primary mirror 404 → fallback used; progress events emitted throughout;
      cancel leaves no partials outside staging.
- [x] `Downloading` state walk on ball (enter/progress/exit) asserted via state log.
- [x] local_override: model loaded from local path with integrity check, no network calls.
- [x] manifest.json committed with license fields filled; P5 conclusion recorded in PRECHECK.
- note: AC#2 **PASS** — 见 docs/evidence/s2/14-adversarial-acceptance.md §"Addendum 裁决（2026-09-20，AC
      补裁）" AC#2 行：downloader_test.go:317 断暂存存活 2000B + **服务端实收 `Range: bytes=2000-`** +
      全文件哈希等值 + staging 清干净，`-count=2` 2×PASS（残余：真·冷启中断预演归 16/SPEC-10 §7）
- note: AC#3 **PASS** — 见同上报告 AC#3 行：:369 断主镜像确被试（`primary.Hits()!=0`）+ 备用确被用 +
      五阶段进度齐 + 装成内容哈希正确；:399 真卡在 PhaseDownloading 后 Cancel，`filepath.Walk` 整个
      DataDir 断无任何 staging 残件，`-count=2` 4×PASS
- note: AC#4 **PASS** — 见同上报告 AC#4 行：bridge_test.go:23 对 `bridge.StateLog()` 断序列全等
      `[Downloading FirstRun]`（SPEC-08 §3 行 #2/#37）+ `Ticks()!=0` + `LastPercent()==100` +
      #37 C29 副作用已发；:60 断 `[Downloading Error]`；:77 断非法进入被拒且态不变（残余：球窗口
      Downloading 环+百分比的渲染核对并入 H1），`-count=2` 6×PASS
- note: AC#5 **PASS** — 见同上报告 AC#5 行：downloader_test.go:451 注入 `panicTransport`
      （:486 任何出网直接 panic，非事后计数器）证"零网络"，且 :475 翻一个 bit 后同 mgr 再 Ensure
      必须失败＝override 不是校验旁路，`-count=2` 2×PASS

## Progress log (append-only, newest last)
- [2026-09-19T15:11:33Z] agent=orchestrator claimed=T14-impl did=dispatched next=sub-agent works through acceptance criteria
- [2026-09-19T15:54:22Z] agent=T14-impl did=internal/models pipeline green: minisign verify (hand-written Ed25519+Blake2b, format-exact), C29 manifest types+signed loader (offline verify before network), Manager (mirror chain, Range resume byte-offset, retry/backoff, failover, progress events, cancel-cleans-staging, sha256 reject+delete, tar.bz2 safe extraction), local_override offline integrity, DownloadingBridge (D43 #2/#37 walk); TDD tests all pass incl -race next=real manifest 6 entries + dev key + sign + spot-check downloads
- [2026-09-19T22:59:23Z] agent=T14-impl did=real manifest landed+signed (6 entries, all hashes from real downloads/local cross-check vs spike table; KWS/VAD/paraformer/SenseVoice=Apache-2.0, CT-Punc=Apache-2.0, matcha-zh-baker=NON-COMMERCIAL -> status=blocked-p3, Ensure refuses); dev minisign keypair E:\work\base\wisp-minisign (keyid a3c8794f3fd94fc5, never in repo), pub hardcoded in buildinfo, tools/signmodels + scripts/sign-models.ps1, models/manifest.json.minisig committed and verifies in TestRealManifestInRepoVerifies; compose.dev/test.yml + model-mirror fixtures (good/corrupt/missing); build.ps1 colocates manifest; REAL pipeline spot checks PASS (vad-silero 11.3s; punc 64MB archive download+extract+verify 16.5s, WISP_IT_REAL_MIRROR=1) next=full-repo gates + commit/push
- [2026-09-19T23:06:46Z] agent=T14-impl did=handoff-to-orchestrator all acceptance criteria met: gates green (go vet/test on my pkgs + repo-wide go test ./... PASS; build.ps1 -Env dev PASS incl manifest colocated into build/models; compose.dev/test config OK); unit commits 92f4861+bfcb230+5ab0339 pushed origin+cnb; NO uncommitted T14 files. Notes for orchestrator: (1) matcha-zh-baker P3 BLOCKED non-commercial (data-baker), status=blocked-p3 enforced in Ensure, TTS replacement decision pending [H5]; (2) dev minisign key E:\work\base\wisp-minisign (keyid a3c8794f3fd94fc5), prod rotation S8; (3) repo-wide vet has 1 issue in parallel agent's in-flight internal/observe/sampler_test.go:330 (not T14); (4) cmd/wisp doctor still labels buildinfo key "placeholder" - stale copy, one-line touch for next owner; (5) hf-mirror 401s on Xet-backed repos (punc) - manifest keeps official GitHub fallbacks, ghfast transport fallback automatic. next=T15 pulls binaries; acceptance states remain in-progress pending orchestrator
- [2026-09-19T23:09:31Z] agent=orchestrator did=T14-adv PASS (orchestrator-executed; report docs/evidence/s2/14-adversarial-acceptance.md; matcha P3 blocked ruling recorded) next=ticket DONE
- [2026-09-20T02:25Z] agent=agent-bookkeeping-1 did=AC boxes reconciled against docs/evidence/s2/14-adversarial-acceptance.md (6-row audit table from the deadline-window inline acceptance; only rows 2/3/4/5 map onto ACs): 2 checked, 4 left open (reasons above — resume/failover+progress+cancel/Downloading-walk/local_override were never adjudicated by name, though their tests exist on disk) next=none
- [2026-09-20T11:25Z] agent=ac-addendum-auditor did=补裁 AC#2/AC#3/AC#4/AC#5（四框原为空白裁决）：全部 **PASS** 并改勾。逐条实跑 `-count=2` 全绿（models 包 4+2+2+6 例），并读断言本体确认非空转——resume 的字节偏移是**服务端实收 `Range: bytes=2000-`**（downloader_test.go:346）而非内部自陈；failover 同时断"主镜像确被试"与"备用确被用"（:381-386）防假转移；cancel 用 `filepath.Walk` 扫全 DataDir（:434-446）；Downloading 走态对 `StateLog()` 断序列全等 `[Downloading FirstRun]`/`[Downloading Error]` 并含 `Ticks()!=0`+`LastPercent()==100`+D43 #37 C29 副作用（bridge_test.go:37-56）；local_override 的"零网络"用 `panicTransport`（任何出网直接 panic，:486）。两条残余明写不改判：真·冷启中断预演归票 16（SPEC-10 §7 要求真实触发）、球窗口 Downloading 渲染核对并入 H1。裁决写入 docs/evidence/s2/14-adversarial-acceptance.md §Addendum next=none
