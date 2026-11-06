# 33 — Panel host C27, implementation leg 33-r14 (AC#6 CSP + AC#8ⓐ clamp teeth)

**Leg:** 33-r14 (code `a1`..`r13` taken; codes checked against `.scratch/wisp/probes/33/` before naming)
**Start HEAD:** `aa7159b29` (dispatch anchor `ac7d14c7f` or current — current was `aa7159b29`, base moved to `e3abe0ebb` mid-leg by the parallel wave45 ticket-23 leg, internal/tools only, zero face overlap)
**Method:** rule 21 shared clone (`git clone --shared`) because `git status --porcelain -- cmd/wisp` was NOT clean (untracked `panel_resident_323_windows_test.go` / `panel_resident_assembly_323_test.go` from leg 323); patch landed on the main tree byte-identically (SHA256 verified, all 7 files).
**Prereqs:** 07 + 12 done (A1588 前置核查); 起手三核 done (ticket read fully, 406 lines; cited anchors re-measured at HEAD, drifts named below).

## What this leg landed (three pieces, all inside `cmd/wisp/**` + one test file each)

1. **AC#6 — CSP response header on every served document.** `panel_control_windows.go`:
   `panelCSPHeaderValue` = SPEC-06 §9's strict CSP verbatim (`docs/specs/SPEC-06-security-gatekeeping.md:126`,
   including `connect-src 'none'`); `panelCSPHeader` = the header line; `panelResponseHeaders(status, contentType)`
   = the one construction site, called from `respondWithEmbeddedAsset` for **200 and 404 alike**.
   The policy travels as a response HEADER on the one serving surface this tree has (the
   `AddWebResourceRequestedFilter` arm), not an injected `<meta>` a page could strip.
   Ruler: `panel_csp_33r14_windows_test.go` — two tests, window-free (they ask `panelResponseHeaders`,
   the filter's own construction; + one AST walk pinning the CSP literal to exactly ONE string literal
   in the file, comment-exempt like d22scan ban #8).
2. **AC#8ⓐ — the zero-execution summary-clamp branch gets teeth.** The branch
   (`cmd/wisp/panel_pump.go:422` at start HEAD; ticket text cites `:186-192`, drifted) still had no test
   reaching it at `aa7159b29` (measured: `panelSummary|summaryClamp` in `cmd/wisp/*_test.go` = 1 comment hit only).
   Because `Publish`'s `Out` takes the same function, a synchronous test could never observe the branch's
   product without taking the slow path too. Naming split (per the AC's own text): production entry renamed
   `bookPanelSnapshotBlocking` (behaviour byte-identical; `run.go:958` `Out:` follows the rename, so the
   production path is unchanged), plus new `tryBookPanelSnapshot` (TryLock; a caller that loses the race
   reports it and books nothing — the observable difference that makes the miss nameable).
   Ruler: `panel_clamp_33r14_test.go` — 5 tests, all through `panel.NewSnapshotPump(...).Publish()` with the
   production entry as `Out` (the exact value `run.go` wires), + the over-bound packet (3×120-rune ids,
   measured: 45-rune ids stayed under the bound) asserting the branch's own output shape (`pending=-`,
   `depth=3` kept, `sha256=` digest kept, line ≤ 440 runes) + the negative control (short packet keeps ids).
3. **Orchestrator's 2026-10-10 顺带 item — the close reading.** `panel_host_windows.go` `Hide()`: one
   stdout line on the same surface `panel_resident_windows.go:446` uses — `wisp: panel window hidden
   (focus returned to hwnd 0x…)` / `… (no prior window to return focus to)`. Observed firing live in the
   full-package run under `TestAC4FocusReturnToPriorWindowGap33r5`. Zero behaviour change (the Win32
   calls and the bookkeeping are untouched).

## Gates (all re-run on the MAIN TREE after the byte-identical copy, CGO_ENABLED=1, sherpa/build PATH, env vars unset)

| Gate | Reading |
|---|---|
| `gofmt -l cmd/wisp internal/panel` | **empty** (rc 0) |
| `gofumpt -l cmd/wisp internal/panel` (`/d/work/base/gopath/bin/gofumpt.exe`) | **empty** (rc 0) |
| `go vet ./internal/panel ./cmd/wisp` | **rc 0**, no output |
| `go test ./cmd/wisp -run '<7 new rulers>' -count=2` | **14/14 PASS** (each ruler twice) |
| `go test ./internal/panel -count=2` | 6 FAIL lines = the SAME 3 names per count round, **both with this diff and at pristine base** `e3abe0ebb` (control clone): `TestComposerContractTypesMatchFrontend` / `TestComposerRenderFixtureTellsTheTruth` / `TestC21DesignTokensFourWayAgree` ⇒ **pre-existing on the shared panel, zero new panel reds from this leg** (write face never touched `internal/panel/**`; dependency direction cmd→panel makes cross-contamination structurally impossible, and the control proves it empirically). ⚠ Note for the orchestrator: `TestComposerRenderFixtureTellsTheTruth` is NOT on the ticket's old known-red roster (which lists approval/composer/frontend_hygiene/tokens_fourway as the 4 at 10-01) — the roster needs re-baselining by a non-implementer. |
| `sh scripts/d22scan.sh` | **rc 0, clean** — `ban #8 internal/=915, cmd/=219` (file-count denominators) |
| Mutation checks (see below) | 3/3 RED with named teeth, restored byte-identical (md5) |

Full-package `go test ./cmd/wisp -count=2` (extra, beyond the dispatched gates): hit the default 10-minute
package timeout on the second count round with ~30 fails, **all reproduced at pristine base in named
controls** (`Test15e*` family: "cannot find build\models\manifest.json / third_party\spike-models\silero_vad.onnx
above the working directory" — untracked artifacts that exist in the main tree but not in a clone;
resident/AC228/AC1 family; `TestAC4FocusReturnToPriorWindowGap33r5` re-proven RED at base in isolation ⇒
the ticket's known load-sensitive focus red, unchanged). Named honestly: this leg did NOT produce a green
full cmd/wisp count=2 run in the clone, and the base cannot either under this environment; the dispatched
gates (internal/panel count=2 + vet + scan + fmt) are green modulo the three attributed panel reds.

## Mutation checks (载具: throwaway python scripts, zero deletion — files removed after restore; mutant never committed)

- **M1** `if len(line) > summaryClamp` → `if false && ...` ⇒ `TestAC8SnapshotClampBranchDropsIdsAndKeepsTheDigest33r14`
  **RED** with 2 named teeth (ids spelled out; line 496 > 440) while the short-packet negative control stayed
  **GREEN** ⇒ the ruler separates the branch from the function's own spelling. Restored, `md5sum -c` OK.
- **M2** CSP `connect-src 'none'` → `'self'` ⇒ `TestAC6CSPHeaderIsOnEveryServedResponse33r14` **RED**
  (policy-comparison tooth names the exact drift) + the one-build-site ruler **GREEN** (it counts the literal,
  not the policy). Restored, `md5sum -c` OK.
- **M3** miss response returns `""` instead of the CSP header ⇒ same ruler **RED** on the 404 tooth only.
  Restored, `md5sum -c` OK.

## What remains (named, per the dispatch's highest-value-slice mandate)

- **AC#7** — NOT implementable as written at this HEAD; two facts, both named not improvised:
  1. The letter's trigger ("本票落地 `Marshal()` 的第一枚**生产**调用者之后") **did not occur as the letter
     predicted**: production consumes the full bytes through `Publish()` (`run.go:958` `Out:` →
     `bookPanelSnapshotBlocking` retaining `lastSnapBytes`, and ticket 329's `pushPanelSnapshotToHost`
     pushing them at the host), not through a separate `Marshal()` call. The 13 `Marshal()` references at
     HEAD are all `_test.go` + definition. A `Marshal()` call added just to light the trigger would be a
     fake caller — the F-PUMP-4 sin this box exists to prevent.
  2. The named fix (hold the pump lock across `Snapshot()`) does NOT deliver the substance the box
     demands: coherence of the packet is NOT provided by `pump.mu` against verdicts/mode/results that live
     under EXTERNAL locks (the gate queue's, the store's, the stream log's); a state change between two
     source reads remains observable with or without it, so the red→green flip it requires is not reachable
     by that edit against real sources. Any real closure needs the owner to define WHICH two-instant
     coherence is required (queue depth vs mode vs results are three different external judges) — that is a
     contract-shape decision, not a leg-local edit.
- **AC#8ⓑ** — the `EvToolStart` consumer (`run.go:2166`, drifted from `:738`) still has **zero emitters**
  (`internal/agent/sink.go:25` definition + comment only, re-measured at HEAD). Deleting the case is a
  product-shape decision (the console tool line and the `changed=true` pulse are correct if an emitter ever
  lands); giving it teeth requires a fake event — F-PUMP-6 explicitly refuses that. Left named, not touched.
- **AC#1–AC#5, AC#9–AC#14** — unchanged by this leg (host already landed by r1–r13/255/299/329 legs; AC#14
  re-homed to ticket 329 by A-ruling 2026-10-16; white-screen to ticket 299). AC#10's [panel] producer
  landed by the 255 legs (schema.go:566-575 → panelWidthPx/panelHeightPx fallback + resident re-read at
  `panel_resident_windows.go:205` + reshow SetSize, with `panel_geometry_255_test.go`/`config_receipt_255_test.go`
  teeth in place). Ticket AC boxes untouched by this leg (non-implementer flips boxes).

## Honesty notes

- The `-count=2` full-package environment failures in the clone are attributed but NOT fixed here — the
  15e/resident family needs the untracked model artifacts, which is the main tree's shape, not a defect.
- The three `internal/panel` reds are pre-existing (base control), but their roster status is stale in the
  ticket (RenderFixture is not on any known-red list I can find) — flagged for the orchestrator, not
  silently accepted.
- `tryBookPanelSnapshot` has no production caller yet (by design: it is the AC#8ⓐ ruler's entry; the
  production path stays on the blocking entry). A merge/backpressure consumer adopting the try-shape is
  ticket 35 AC#5's territory.
- Zero writes to: `internal/panel/**`, `frontend/**`, `design/**`, `docs/reports/pending-and-issues.md`,
  ticket file (progress line lives in the commit message per this leg's dispatch write face),
  `docs/SLO.md`, `internal/observe/thresholds.go`, `go.mod`/`go.sum`.
- Zero GUI, zero `wisp.exe` run, zero `*winlive*` test, zero push.
