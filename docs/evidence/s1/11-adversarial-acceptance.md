# Ticket 11 — adversarial acceptance (LLM REST adapters)

**Acceptor:** orchestrator (not an implementer of this ticket). **Date:** 2026-09-20.
**Commits under review:** `35d4c74`, `94071ff`, `dff7abd`, `43f9a53`, `f87c406` (plus the earlier
adapter work recorded at log `07:10Z`/`07:55Z`).
**Method:** every row below rests on something I ran or read myself. Where the implementer's log
made the claim, I say so and mark it as not independently derived.

## Verdict table (1:1 with the ticket's AC boxes, per README rule 6)

| AC | Verdict | Evidence I produced myself |
|---|---|---|
| 1 — both adapters pass the **same** golden/fault suite as OpenAI Chat | **PASS** | This was the shortcut I explicitly predicted and warned about ("write a weaker second harness and call it the same one"). **My suspicion was disproved by the code.** One shared table exists: `adaptertest/harness.go:174 func CaseTable()`, iterated by `RunAll` at `harness.go:479`. All three adapter packages call it — `anthropic/suite_test.go:36` and `openairesponses/suite_test.go` via `adaptertest.RunAll(t, unit())`, `openaichat/harness_golden_test.go:70` directly. Two anti-skip guards, both real: `AssertCoverage` (a skipped row or orphan fixture fails) and `harness.go:277-279` asserting `len(CaseTable()) == len(AllScenarios())`. I counted **13 rows** in `CaseTable` and ran `TestScenarioCoverage` → **PASS**. All three packages green at `-count=2`. |
| 2 — Anthropic cache-breakpoint: stable cached prefix across turns, cache-read usage populated | **PASS (log-derived, structurally spot-checked)** | `anthropic/cache_test.go` reads the request body from the **live** mockllm via `/__control/last_request` rather than a locally built struct, asserts exactly one `cache_control` on the last system block, none in the conversation suffix or tools, byte-identical prefix across two turns, and `cache_read_input_tokens → Usage.CachedTokens=8`. I verified the body-capture mechanism exists and is recorded before the first response byte (`recordedRequest` in `tools/mockllm/server.go`, populated on the routes that call `recordRequest`). I did **not** re-run this test in isolation; it passed in my `-count=2 ./internal/llm/...` sweep. |
| 3 — `ReasoningDelta` surfaces on both adapters | **PASS (same caveat)** | Covered by the shared table (`RunAll` rows for thinking responses) plus the dialect fixtures; `TestC6EventsIdenticalAcrossAdapters` compares a canonical event trace across all three, so a missing delta on one adapter breaks cross-adapter equality rather than passing silently. Green in my sweep. |
| 4 — fallback chain: `fail_next(3)` on primary → fallback → completes; double failure → `Error(provider)`, ctx preserved, resumable | **PASS** | `internal/llm/fallback_test.go` is written against the exact trap: its header comment states "a chain that never calls the primary still passes a naive fallback test", and every case asserts **the primary server's own request counter** (3 = initial + 2 retries), so a chain that skips the primary fails. The fixture stands up two live mockllm instances with the primary on the Anthropic dialect. Double failure asserts `Error(provider)`, `ctx.Err()==nil`, and that the same `*Request` re-streams after the fault clears. |
| 5 — §14.2 failure matrix: one test per row asserting the emitted user-visible event/state | **PASS (same caveat)** | `matrix_14_2_test.go`, one test per §14.2 row, including the row-2 requirement that exhaustion produces an explicit classified `Error` and never silence, and row-3's 401/403 → `Unconfigured` + never retried. Row 1's ">5s notice" is asserted on an **injected monotonic clock**, not on wall-clock deltas — which is the shape this repo's forbidden-pattern rule demands. Green in my sweep. |
| 6 — probe suite: capability emulation → `provider_health` ✓/✗ + declared-vs-measured mismatch | **PARTIAL — the proven half is real; the missing half is transferred, not failed** | See the two sections below. |
| 7 — token bucket: rpm=10 → 11th request within 60s waits locally, no server 429 | **PASS (log-derived)** | `94071ff`'s test asserts the 11th request never leaves the process, which is the only assertion that distinguishes local pacing from a server-side 429 round-trip. Green in my sweep. |

## AC#6 — what I verified myself, and where its edge actually is

**The measuring half is proven, and I proved it independently rather than trusting the log.**

1. **The judge was not bribed.** `tools/mockllm/capability_test.go` — the file containing the three
   tests that were red when I took over — has **exactly one commit in its entire history** (`dff7abd`,
   its addition). Nothing edited it afterwards, and the working tree carried no modification to it.
   So the three tests went green on **implementation changes alone**. This is the single most
   important check on this ticket and it is now closed with git as the witness.
2. **The probe measures instead of restating the catalog.** My own mutation: I replaced
   `rep.Flags.Set(capabilityFlagKey(c), res.OK)` at `probe_health.go:255` with
   `declaredCapability(o.Declared, c)` — i.e. the probe reports what config claims — and
   `TestProbeSuiteMeasuresBrokenFC`, `TestProbeSuiteMeasuresBrokenVision` and
   `TestProbeSuiteBrokenOnAllThreeDialects` all went **RED**. Reverted; `grep -c MUTATION` on that
   file is now 0.
3. **Non-vacuity is structural, not asserted once.** Each case pins three independent facts: the live
   mockllm is provably in the mode (`/__control/state` read back), requests really left
   (the server's own per-route counter, 0-vs-3), and the SQLite `provider_health.probe_json` row says
   the opposite of config when the wire does. Nine `TestProbeSuite*` cases exist, including the two
   mirror directions (`MeasuresBrokenFC`/`MeasuresBrokenVision`), an honest-provider zero-mismatch
   case, an **honest-negative** case (declares ✗ and measures ✗ → must NOT emit a mismatch, because
   the event is about contradiction, not failure), `AudioStaysUnprobed` (records no verdict rather
   than a false one), and `RefusesSilentRuns` (no sink or no notice → error, not a no-op).
4. **A real hole the implementer found on its own:** mockllm's Anthropic dialect ignored
   `tool_choice {"type":"any"}`, which is what the adapter actually sends for a required tool call —
   so an fc probe could never have measured a *capable* Anthropic provider. Only the broken direction
   was detectable. `TestProbeSuiteCapableOnAllThreeDialects` now pins both directions on all three
   wire shapes. I confirmed that test exists.

**The half that is genuinely missing, and why the box stays unticked.** Nothing in the shipped binary
calls `RunProbeSuite`: the composition root is `cmd/wisp`, and the probe's write path is a `HealthSink`
interface rather than a direct dependency on the storage layer (`internal/llm` still does not import
it). So on a real machine the 「声明 ✓ / 实测 ✗」 event can never fire. Secondary gap: the thinking
capability reuses ticket 09's check, which accepts a plain text answer and therefore **cannot detect a
broken thinker**.

Both are registered as **A11** in `docs/reports/pending-and-issues.md` and handed to **ticket 12**
(the S1 end-to-end gate) with completion criteria. The implementer's own breakpoint said
"next = user wires RunProbeSuite"; **that wording is void** — this owner does not write code, and the
item belongs to a ticket.

## Gates I re-ran myself

```
cd tools/mockllm && gofmt -l . && go vet ./... && go test -count=2 ./...   → ok 0.214s
go vet ./internal/llm/...                                                  → clean
go test -count=2 ./internal/llm/...  → ok internal/llm  ok anthropic 29.663s
                                      ok golden 0.078s  ok openaichat 46.628s
                                      ok openairesponses 38.088s
go test -run TestScenarioCoverage ./internal/llm/anthropic/                → PASS
```

## Disposition

**Clean to close with AC#6 formally transferred.** No BLOCKER, no MAJOR. Two notes worth carrying
forward, neither of which is a defect in this ticket's code:

- **A pattern, now twice observed.** Ticket 63's AC#6 (A8) and ticket 11's AC#6 (A11) failed the same
  way: *proven inside a test harness, never called from the composition root.* Both were handed to
  ticket 12. When ticket 12 is planned, treat "wires the thing that already exists" as its primary
  job rather than an integration afterthought — that is where this project's remaining gaps live.
- AC#2/#3/#5/#7 rest on the shared suite passing plus my structural reading, not on per-test
  re-runs. I say so plainly instead of implying I re-derived each one.
