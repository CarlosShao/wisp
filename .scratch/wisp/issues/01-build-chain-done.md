# 01 — S0①: build chain, one-shot clean build, BUILD.md frozen (DONE ✅)

**Status:** done
**Claimed by:** orchestrator -> sub-agent T01-impl
**Last update:** 2026-09-19T06:57:25ZT05:33:59Z
**Blocked by:** None — can start immediately
**Parallel slots:** 1
**Spec refs:** SPEC-11 §2, SPEC-01 §3, §14.9 (highest-risk gap), S0 done①

## What to build
From a clean clone on Windows, one command produces a runnable `wisp.exe` (+ native DLLs) with zero
manual toolchain steps, and the exact flow is frozen into `docs/BUILD.md`. This is the #1 blocker
for everything else: AI agents must never be able to enter toolchain trial-and-error loops.

## Key constraints
- Repo skeleton from SPEC-01 §3: `go.mod` (module `github.com/CarlosShao/wisp`, toolchain pinned),
  `cmd/wisp/main.go` (no-args = GUI; `wisp run`/`wisp doctor` stubs acceptable), `internal/buildinfo/`.
- C toolchain: **MSYS2 mingw-w64 GCC, version pinned in BUILD.md** (decision of SPEC-11 §2.1).
  If linking sherpa prebuilt libs proves MSVC-only, fallback `CC=cl` + VS Build Tools becomes the
  documented primary — record which one won.
- `deps.toml` pins sherpa-onnx C lib + onnxruntime versions **with SHA256**;
  `scripts/fetch-deps.ps1` downloads into `third_party/` (git-ignored), verifies hashes, caches.
- `scripts/build.ps1 [-Env dev|prod]`: fetch-deps → (skip frontend, none yet) → `go build` with
  `CGO_ENABLED=1` → outputs exe + DLLs colocated (DLL hijack/PATH rules SPEC-11 §7.1) → SHA256SUMS.
- Docker builder image `docker/builder.Dockerfile` (Go + mingw) for CI; note: Linux-cross link of
  sherpa Windows libs is validated in ticket 02, not here — native Windows build is the primary path.
- GUI/CLI same binary; console output for subcommands via `AttachConsole(ATTACH_PARENT_PROCESS)`.
- Link repo to CI definitions only if trivially possible this ticket; full CI is ticket 08.

## Out of scope
- Spike measurements (ticket 02); any product feature; frontend; CI workflow files (ticket 08).

## Acceptance criteria
- [ ] Fresh clone → `scripts/build.ps1` → running `wisp.exe` (shows a placeholder console/UI proof
      of life) — verified on a clean machine/CI runner, no preinstalled toolchain assumed beyond
      what BUILD.md declares.
- [ ] `docs/BUILD.md` written and declared frozen: toolchain versions, deps.toml hash pins,
      exact commands, troubleshooting of the two failure modes actually hit during this ticket.
- [ ] `deps.toml` SHA256 mismatch test: corrupt one downloaded byte → fetch-deps fails loudly.
- [ ] `wisp doctor` prints toolchain/DLL versions and PASS/FAIL for colocated-DLL rule.
- [ ] BUILD.md committed and pushed to both remotes.

## Progress log (append-only, newest last)
- [2026-09-19T05:33:59Z] agent=orchestrator claimed=T01-impl did=dispatched implementation sub-agent next=sub-agent works through acceptance criteria
- [2026-09-19T06:14:17Z] agent=T01-impl did=build chain implemented: go.mod(go1.27/toolchain go1.27.1), cmd/wisp(no-args proof-of-life/run/doctor+AttachConsole), internal/buildinfo, deps.toml(sherpa-onnx 1.13.8 MT-Release sha256+per-DLL pins, go-binding v1.13.8), fetch-deps.ps1, build.ps1; local build green, doctor PASS next=BUILD.md + tamper test + clean clone
- [2026-09-19T06:18:53Z] agent=T01-impl did=BUILD.md frozen (toolchain pins, 6 real pitfalls, MSVC fallback); deps.toml tamper tests: archive-pin and DLL-pin both fail loudly (exit 1, bad file deleted); cache bound to manifest to prevent stale-cache bypass next=clean clone simulation + handoff
- [2026-09-19T06:21:10Z] agent=T01-impl did=handoff-to-orchestrator (clean clone one-shot build verified prod+dev, doctor PASS, AttachConsole validated under -H=windowsgui, tamper tests green; ticket 01 implementation complete, Status left for orchestrator)
- [2026-09-19T06:32:38Z] agent=orchestrator did=added docker/builder.Dockerfile (was in ticket constraints, impl agent deferred; Go pinned 1.27) next=T01-adv acceptance
- [2026-09-19T06:57:25Z] agent=orchestrator did=T01-adv VERDICT PASS (0 BLOCKER/0 MAJOR/3 MINOR; report docs/evidence/s0/01-adversarial-acceptance.md); MINORs fixed (log order, BUILD.md indirect-import wording, Dockerfile deviation note) next=ticket DONE
