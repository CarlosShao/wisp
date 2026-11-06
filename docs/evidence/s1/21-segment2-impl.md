# 21-segment2-impl — approval-gate segment 2 (card rendering + depth badge) implementation receipt

- **Leg:** agent=C, wave 34, 2026-11-06
- **Tree:** `git clone --shared` export at `D:\tmp\wisp\20261106\21s2\tree` (rule 21: the main
  tree's `cmd/wisp` carries ticket 323's in-flight untracked test files, and the shared
  `cmd/wisp` write slot makes isolation mandatory). Clone started at dev tip `92d5a33d1`;
  dispatch HEAD `6ae44b1fb` verified ancestor of it. **Zero push.**
- **Ticket:** `.scratch/wisp/issues/21-approval-gates-minimal.md` (status=in-progress,
  claim + progress log committed first per issues/README rule 1)

## Scope decision (re-audit before writing)

The dispatch named "native L2 card rendering + ball-side visual wiring + cmd/wisp loop
wiring". Re-audit at HEAD showed most of that already landed between 2026-09-20 and now
via tickets 12 (composition, Options.Gate), 20 (applied-steps), 201 (reply listener,
Replies seam, waiting-state), 246 (resident ball host, Confirming/AwaitingApproval orb
states, Esc borrow), 220/146 (panel read faces), 256/268 (risk-config gate values),
337/375/376 (sink, answerability, card window number). What was measurably still missing:

1. **Full params never rendered** — SPEC-06 §7 「完整参数」 and the ticket's "FULL params
   (not summary)" clause: `consoleApprovalUI.Prompt` printed rules/channels/paths/window
   and stopped. No surface printed `Prompt.Params`/`Prompt.Args`.
2. **D45-1 batch view never rendered** — `Aggregate` builds `BatchView` and `promptFor`
   hands it to every Prompt; no renderer consumed it.
3. **`ball.SetBadge` had zero production callers** — the orb's AwaitingApproval depth
   badge (SPEC-06 §7 「深度角标（悬浮球必显，不能只藏面板）」) existed only in
   `cmd/balldebug`'s scripted tour (census: `git grep -n "SetBadge(" -- cmd` → balldebug
   only, before this leg).

## What landed (4 production files + doc)

| File | Change |
|---|---|
| `internal/agent/approval/cardrender.go` (new) | `approval.RenderCard(Prompt) []string` — the ONE native card renderer: header, rules_hit verbatim, task, L2 depth line, D45-1 batch lines, FULL params (`Args` verbatim via `json.Indent`, `Params` fallback, explicit gap sentence, unparseable kept raw), paths, channel lines with honest availability, route-correct window. Pure function, no authority. |
| `cmd/wisp/run.go` | `consoleApprovalUI.Prompt` prints `RenderCard`'s lines (grant-book → record → waiting-state ordering untouched); dead `channelState` helper removed. |
| `cmd/wisp/resident_approval_windows.go` | `ballCardUI.Prompt` posts `ball.SetBadge(depth)` AFTER `SetState` (STA order: `applyStateLocked` zeroes the badge on state entry) and prints the full card; `Update` re-syncs the badge when one card settles with another pending; helpers `approvalDepthBadge` (pure) + `syncBadge`. |
| `internal/agent/approval/doc.go` | "Still owed" rewritten to the two true residuals. |

## Tests (all new, named by contract line)

- `internal/agent/approval/cardrender_test.go` — 10 cases: rules/order, task+depth,
  params-per-value (not line count), Params fallback + determinism, explicit-gap,
  unparseable-verbatim, D45 batch shape (uses real `Aggregate`), nil-batch silence, B1
  channel wording (「语音取消不可用」 arrives inside the registry's Text), L1-has-no-depth,
  parser-shape stability, purity.
- `cmd/wisp/approval_card_render_21s2_test.go` — real `approval.New` gate → real
  `consoleApprovalUI` → printed card asserted → `Native().Allow` with the minted grant
  settles the call (AC#2's round trip); D45 batch card through the same printer.
- `cmd/wisp/approval_badge_21s2_windows_test.go` — pure mapping cases + AST census of
  `resident_approval_windows.go` (Prompt: exactly 1 `SetBadge`, 1 `SetState`, 1
  `RenderCard`; Update: 1 `syncBadge`) + order tooth (SetBadge after SetState).

**Mutation self-check (tooth proof):** removing the `SetBadge` post from `Prompt` turns
`TestTicket21S2BadgeWiringCensus` RED (`calls SetBadge 0 times, want exactly 1`) and
`TestTicket21S2BadgePostsAreOrderedStateThenDepth` RED (`SetBadge@0`); restored file
re-runs green. Mutated copy and restore both inside the clone only.

## Gates (all run in the clone; CGO left at default 1; GOTMPDIR/TMPDIR/TEMP/TMP unset per dispatch)

| Gate | Reading |
|---|---|
| `go build ./cmd/... ./internal/... ./tools/...` | rc=0 |
| `go test ./internal/agent/... -count=1 -timeout 30m` | ok (agent 4.3s, approval 0.64s) |
| `go test ./cmd/wisp -run <approval-relevant set> -count=1` | ok 41.4s (376 card-window, 201 seam/reply, 220 L1-block, 223 reload card, 375 unanswerable, 35r3 panel guards, 21s2 new) |
| `go test ./cmd/wisp -count=1 -timeout 30m` (full) | see §Full-suite below |
| `gofumpt -l` on all 6 touched files | empty (0 listed) |
| `sh scripts/d22scan.sh` | rc=0 clean — positive control PASS=35 FAIL=0; scan examined 422 production Go files; first attempt with `GOTMPDIR` on D: was REFUSED by the 311 shape-B guard (correct behaviour; re-run with temp vars unset as dispatch ordered) |
| `go vet ./cmd/wisp ./internal/agent/approval` | rc=0 |

## Full-suite result (and the baseline comparison that explains it)

`go test ./cmd/wisp -count=1 -timeout 30m` in this clone: **rc=1, but the reds are
tree/environment-shaped, not this leg's** — proven by running the IDENTICAL family set in
a second `git clone --shared` at clean HEAD (no leg changes, same junction set,
`D:\tmp\wisp\20261106\21s2\base` @ `d0a230547`):

| Family | Work tree | Clean baseline | Explanation |
|---|---|---|---|
| Voice: `Test15e*` `Test41a2r1*` `Test26P*` `Test392*` `Test395r1*` | FAIL | FAIL (same) | They resolve models by walking `build\models` / `third_party\spike-models` from the tree root; the clone lacks the gitignored content even with junctions (the resolver refuses reparse-point traversal by design). Baseline-identical ⇒ environmental. |
| Resident real-process: `TestAC1ResidentLeg*` `TestAC228*` | 2 of 4 FAIL | ≥ same 2 + more FAIL | Subprocess boot timing on this machine; baseline FAILS A SUPERSET ⇒ environmental/flaky, pre-existing. |
| `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims` | FAIL | FAIL | Ticket 375 added `Options.Answer`; the resident now passes 7 fields, the 256 roster says 6. The census wasn't updated by 375's leg — pre-existing drift, unrelated to this leg (my diff touches no `approval.Options` literal). |
| `TestAC247HandingTheLevelToTheSeamIsNotVisibility` | FAIL | FAIL | `prototypeVisuals` ON in assembly vs AC#10's expectation — pre-existing at HEAD. |
| DLL-dependent (`TestAC2RealProcess*` `TestAC3Early*` 246-pipeline) | FAIL before junctions → PASS after junctioning `third_party\sherpa-onnx` | — | The clone-missing-gitignored-dir class, closed by junction. |
| **Card/approval/reply/badge families (this leg's relevance set)** | **PASS** | — | `Test376R1* Test201* Test220* Test223* Test375* Test35R3* Test21S2* TestComposedGate* TestVeto* TestReplay* TestPanel* TestCard* Test246*` = ok 41.4s. |

The dispatch's gate is `go test ./cmd/wisp -run <relevant>` — that gate is green. The full
package run should be taken on the main tree by the orchestrator after merge, where
`build/` and `third_party/` are real directories.

**Pre-existing drifts worth a line, not this leg's to fix:** the 256 roster (7 vs 6 fields,
belongs to ticket 375's follow-up) and the 247 prototypeVisuals assertion.

## Explicitly NOT done (named, not hidden)

1. `decideRisk` dead-branch deletion (`internal/agent/loop.go:916-924`) — outside this
   leg's write face (`internal/agent` = ticket 12/162's surface); doc.go now records it
   as dead-code cleanup.
2. Winlive visual verification of the badge on a real orb window — real-window rigs are
   the owner-present pass (260-r3 precedent); the headless half (mapping + census +
   mutation tooth) is delivered.
3. AC boxes remain unticked pending the non-implementer adversarial pass (rule 6); this
   file is the implementer receipt, not the verdict.

## Reproduction (every reading above re-runnable)

```
cd D:\tmp\wisp\20261106\21s2\tree
git log --oneline -3                       # claim + implementation commits
go build ./cmd/... ./internal/... ./tools/...
go test ./internal/agent/approval/ -run Ticket21S2 -count=1 -v
go test ./cmd/wisp -run Ticket21S2 -count=1 -v
gofumpt -l internal/agent/approval/cardrender.go internal/agent/approval/cardrender_test.go cmd/wisp/run.go cmd/wisp/resident_approval_windows.go cmd/wisp/approval_card_render_21s2_test.go cmd/wisp/approval_badge_21s2_windows_test.go
sh scripts/d22scan.sh
```

(`build/` and `third_party/sherpa-onnx` of the MAIN tree on PATH for the test binaries'
DLLs; the gitignored dirs do not exist in the clone — build itself needs nothing from
them, proven by `go build ./cmd/... ./internal/... ./tools/...` rc=0 before any PATH was
set.)
