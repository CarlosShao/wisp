# 01 — S0①: build chain, one-shot clean build, BUILD.md frozen

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
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
