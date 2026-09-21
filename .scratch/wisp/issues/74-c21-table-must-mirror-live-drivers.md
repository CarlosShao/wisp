# 74 — Make the C21 token table mirror the constants that actually drive the ball

**Status:** ready-for-agent
**Type:** defect-fix (contract fidelity, cosmetic surface but governance-loaded: C21 is a frozen contract)
**Blocks:** nothing · **Blocked by:** nothing — ticket 69 shipped the machine check that discovered this
**Packages:** `internal/ball` + `docs/evidence/s1/c21-native-tokens.md` only. Do **not** touch
`internal/tools`, `internal/observe`, `internal/secret`, `tools/d22scan`, `internal/llm` (other agents in flight).
**Spec refs:** C21 DesignTokens, registry **A24-D4**, ticket 69 report D7–D12

## Proven premise (measured by `4ca2c6c` / `9c29555`, do not re-derive)
`internal/ball/tokens_table_test.go` now cross-checks the C21 table against `tokens.go` both ways.
It emitted **29 named divergences** the check deliberately does not fail on. Three families:
1. **Code uses constants the table never mentions** — `hit.go` `RingMarginPx`/`ClickTolerancePx`,
   and **9 animation constants in `liquid.go`**.
2. **Semantic drift (the dangerous one)**: `tokens.go` declares `SwimLevelGain` / `SpinLevelGain`
   with **zero consumers**, while the real driver is `liquid.go`'s `Spin*RadPerS` group. So the table
   can document a number that is wired to nothing while the number that moves pixels on screen is
   undocumented. **This is the same shape as A33 ("声明✓ / 实测✗") and must be treated as the priority.**
3. **Table rows nobody consumes**: 18 `Palette` fields, `looks[*].Deep`, `Dur*Ms`, `CountdownFontPx`.

## What to build
Kill family 2 first, then family 1, then decide family 3 explicitly — do not silently do nothing
about any of them.
- **Family 2**: delete the zero-consumer `SwimLevelGain`/`SpinLevelGain` **in the same commit** that
  points the table at `liquid.go`'s real drivers (removing a documented token and adding its
  replacement in one step; never leave the table pointing at a dead constant). If you instead prove
  the dead constants are the intended future contract, write that proof in the commit message and
  make the *test* enforce it — do not leave both halves floating.
- **Family 1**: put the 11 unconsumed-but-live constants into the table with their real units.
- **Family 3**: for each of the 29→N remaining table-only rows, either (a) add the missing consumer
  assertion, or (b) list it in an explicit `documented-but-unconsumed` exemption table **with a
  one-line reason per row**. A blanket "not used yet" is not a reason.
- Tighten `tokens_table_test.go` only where you can name the exact bug class it newly catches: the
  geometry match is currently *containment*, not bijection (two rows can mask each other's number),
  and `FontFamily` has no value assertion. Fixing containment→bijection is in scope **only if you
  show a mutation that today passes and must fail**; otherwise record it as a gap.

## AC (1:1 verdict table required, one row per box)
- [ ] **AC#1** `SwimLevelGain`/`SpinLevelGain`: no longer both "declared" and "unconsumed" — either
  deleted-with-replacement-in-one-commit, or the test pins them to a named consumer. Show the grep
  that proves the consumer exists (`git grep -n`).
- [ ] **AC#2** Every constant that *changes pixels* under `internal/ball` (hit.go, liquid.go,
  renderer_windows.go) is either in the C21 table or in the exemption list with a reason; the
  machine check itself proves it (the report count goes to 0 unattributed).
- [ ] **AC#3** Mutation proof for the new/changed assertions: pick at least 2 rows, change the code
  value, show the specific test go red by name, grep-prove the mutation landed, restore, grep-prove
  it is gone (`git diff --quiet` on both files).
- [ ] **AC#4** Table↔`tokens.css` gap: either add the cross-check, or move this AC to a written
  hand-off (name the ticket). **"Deferred" without a ticket number is not an option** — the owner's
  standing rule is that postponed work must be registered with its completion criteria.
- [ ] **AC#5** Gates: `gofmt -l internal/ball` empty, `go vet ./internal/ball/...` rc=0,
  `go test ./internal/ball/... -count=2`, and `go test ./internal/ball/ -v` with **RUN count == 2×
  distinct test names and zero `SKIP` occurrences** in the log (SKIP counts as ok — that trap has
  burned this project twice).

## Explicit prohibitions for this ticket
- No `git commit --amend`, `reset`, `rebase`, `stash`, `checkout .` (A34). Use
  `git commit -F - -- <explicit paths>` with a **quoted** heredoc; an unquoted heredoc let the shell
  execute backticks in a commit message on `9bf05cb` (A31) and also swept a foreign staged rename
  into that commit.
- **Do not edit `docs/PLAN.md` or `docs/specs/*`.** C21's *contract text* is frozen; this ticket only
  touches the evidence-layer table and code. If you believe the contract itself is wrong, stop and
  report it as a D22 item instead of editing.
- Do not delete a guard in one commit and add its replacement in the next (AC#1 family 2 is the one
  place a deletion is allowed, and only when paired).
- No `t.Skip`, no lowered thresholds, no rewriting golden values to make a check pass.
- If an AC cannot be proven, **leave the box empty and say why**.
- First checkpoint commit **within your first 15 tool calls**; sync Status + boxes + a `next=`
  Progress log line on every commit (a ticket whose face lags HEAD gets its work redone by the next
  agent).
