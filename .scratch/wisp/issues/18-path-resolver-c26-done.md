PathResolver (C26) + A/B sensitive-path blacklists (DONE ✅)

**Status:** done
**Claimed by:** orchestrator -> sub-agent T18-impl
**Last update:** 2026-09-20T00:10:30Z
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: resolver pipeline; B: blacklists + red-team test suite on
Windows runner)
**Spec refs:** SPEC-06 §4, D33/F1 (CRITICAL), C26, D22 path ban

## What to build
The ONLY path normalization entry point in the codebase: expand → absolutize → Clean →
handle-based real path (`GetFinalPathNameByHandle`) → reject reparse points by default →
expand 8.3 short names → normalize UNC; plus the A/B sensitive-path blacklists and the static
ban on raw `filepath.Clean|Abs` for fs decisions.

## Key constraints
- Pipeline order fixed (SPEC-06 §4): env/~ expansion → absolute → Clean → open-handle real path
  (VOLUME_NAME_DOS) → `FILE_ATTRIBUTE_REPARSE_POINT` detection → default DENY junction/symlink
  ("parse-then-continue" is forbidden — it misses source-outside/target-inside edges);
  `reparse_point_exceptions` allows per-explicit-path bypass (TOML, user-set).
- macOS path (for later): realpath + lstat symlink reject — interface stub only (DEFERRED).
- A blacklist (absolute, non-overridable, read AND write denied): `~/.git-credentials`,
  `.git/config`, `%APPDATA%\wisp\config.toml`, `~/.ssh/**`, `~/.aws/credentials`,
  `~/.kube/config`, browser credential stores (Login Data/Cookies/Web Data/Local State),
  DPAPI master keys (`%APPDATA%\Microsoft\Protect\**`), `%LOCALAPPDATA%\Microsoft\Credentials\**`.
- B blacklist (default deny; per-file override = one L2 confirm + log): `.env*`, `*.pem`,
  `*.p12|*.pfx`, `id_*`, `secrets.*`, `*credentials*.json`.
- Resolution results feed R2/R3 (17): allowlist membership + blacklist class.
- Static enforcement: linter/grep rule failing CI on `filepath.Clean|filepath.Abs` outside the
  risk package's sanctioned call (D22 ban) — add to CI in this ticket.
- All path red-team tests must create REAL junctions/short names via `mklink /J`/fsutil on the
  Windows runner — string-mock tests are invalid.

## Out of scope
- Taint marking (19); fs tool enforcement call-sites (20).

## Acceptance criteria
- [x] Red-team four-bypass suite on real Windows: junction, 8.3 short name, UNC, `\\?\` prefix —
      each pointing at an A-list file → all DENIED (plus case/dot-mix variants).
      （裁决：adversarial 报告 §2 PASS——真 mklink /J + GetShortPathNameW + UNC×3 拼写 + `\?\` + evil-twin 兄弟 junction 全拒）
- [x] Reparse default-deny + explicit exception path works (allowlisted junction passes; log).
      （裁决：报告 §4 PASS——按具体路径、大小写不敏感、无前缀泄漏）
- [x] A-list: all entries denied read+write; no grant/override can unlock (unit + integration).
      （裁决：报告 §3 PASS——A 档 9 类锚点全 deny 且 Gate 无 override）
- [x] B-list: default deny; override path emits L2 request + log entry.
      （裁决：报告 §3 PASS——B 档 6 规则默认 L2 + 单文件豁免出审计日志）
- [x] CI static ban active: seeded violation fails the build.
      （裁决：报告 §5 PASS——tools/d22scan 全仓 clean，allowlist 1 行豁免经裁定）
- [x] Resolver idempotence + perf: ≤1ms per call on warm handle cache (budget for ≤50 calls/task).
      （幂等：`TestResolveIdempotent`；性能：2026-09-20 补 bench 实测 **0.397–0.419 ms/op**，
      预算 1ms 余量 ~2.4x，50 次/任务 ≈ 21ms。重跑判据见下方 Progress log 末行）

## Progress log (append-only, newest last)
- [2026-09-20T00:10:30Z] agent=orchestrator claimed=T18-impl did=dispatched (50-min deadline window; hard stop 08:40 local, clean-unit boundary only) next=work
- [2026-09-20T00:31:00Z] agent=T18-impl did=C26-pathresolver+A/B-blacklists+red-team-suite(10 tests, real mklink/J junctions+8.3 via GetShortPathName+UNC+\?\ all deny A-list; reparse exceptions pass) next=none-resolver-done; R2/R3 feed-in (17) uses Resolve+Classify+Gate; T17 parallel files untouched (their TestRuleTaintR4/vet errors pre-existing)
- [2026-09-20T00:27:30Z] agent=orchestrator did=adversarial PASS (orchestrator-executed; reports docs/evidence/s1/18-adversarial-acceptance.md; full-package tests + race green after both agents merged) next=ticket DONE
- [2026-09-20T02:35Z] agent=orchestrator( Carlos ) did=AC boxes reconciled against the acceptance report (all 6 PASS; report §1–§6) + **closed the DEFERRED perf AC**: added `internal/risk/pathresolver_bench_test.go` (BenchmarkResolveWarm) and `internal/risk/pathresolver_budget_norace_test.go` (TestResolvePerCallBudget, build-tagged !race so -race doesn't fake a breach) — measured 0.397 ms/op (test, 2000 samples) / 0.419 ms/op 3688 B/op 53 allocs (bench -benchtime=2000x) vs 1ms budget. Also repaired the garbled evidence path in the line above next=none-ticket-18-fully-booked
  - 重跑判据（任何人可验证，勿信本行文字）：`export PATH="/d/work/base/go/bin:/e/work/base/msys64/mingw64/bin:$PATH"; export GOPATH=/d/work/base/gopath; go test ./internal/risk/ -run TestResolvePerCallBudget -bench BenchmarkResolveWarm -benchtime=2000x -count=1 -v`
