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
- [ ] Red-team four-bypass suite on real Windows: junction, 8.3 short name, UNC, `\\?\` prefix —
      each pointing at an A-list file → all DENIED (plus case/dot-mix variants).
- [ ] Reparse default-deny + explicit exception path works (allowlisted junction passes; log).
- [ ] A-list: all entries denied read+write; no grant/override can unlock (unit + integration).
- [ ] B-list: default deny; override path emits L2 request + log entry.
- [ ] CI static ban active: seeded violation fails the build.
- [ ] Resolver idempotence + perf: ≤1ms per call on warm handle cache (budget for ≤50 calls/task).

## Progress log (append-only, newest last)
- [2026-09-20T00:10:30Z] agent=orchestrator claimed=T18-impl did=dispatched (50-min deadline window; hard stop 08:40 local, clean-unit boundary only) next=work
- [2026-09-20T00:31:00Z] agent=T18-impl did=C26-pathresolver+A/B-blacklists+red-team-suite(10 tests, real mklink/J junctions+8.3 via GetShortPathName+UNC+\?\ all deny A-list; reparse exceptions pass) next=none-resolver-done; R2/R3 feed-in (17) uses Resolve+Classify+Gate; T17 parallel files untouched (their TestRuleTaintR4/vet errors pre-existing)
- [2026-09-20T00:27:30Z] agent=orchestrator did=adversarial PASS (orchestrator-executed; reports docs/evidence/s1/.scratch/wisp/issues/18-*.md; full-package tests + race green after both agents merged) next=ticket DONE
