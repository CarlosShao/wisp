# 136 AC#15 - target census r2 (READ-ONLY)

agent: `auditor-ticket136-ac15-census-r2`
date: 2026-09-24, ticket face: `.scratch/wisp/issues/136-the-zero-sample-fail-closed-guard-underneath-the-freshness-nail-has-no-nail-54-tests-stay-green-without-it.md` AC#15 (`:348-363`) plus its `>` correction block (`:365-370`)
role: census only. Nothing was built, run, formatted or edited except this file.
Constraint honoured: a resource-sampling CI run is live on this machine, so no `go build|test|vet`, no `gofmt`, no `wisp`, no `docker`. Every number below is either a line on disk (Read/Grep) or a `git` read.

---

## 1. Where the target actually is now

Anchor at the time of writing this section: `git rev-parse HEAD` = `a1fd5bf1795196df3913686b804664c787833700`.
`internal/observe/` is re-verified against this anchor at the end of the file (section 7).

AC#15's face text puts the guard at `sampler_settle_coverage_136_test.go:189`. That is stale, as the face
itself warns (`:288` "到时要现量重划，不许照抄今天的行号"). Under `ca2c55e` the same line sat at `:188-189`.
At the anchor above it sits at `:213-215`. The text of the red sentence is unchanged, verbatim.

### 1(a) the precondition / "not enough reads" guard of the flaky leg

`internal/observe/sampler_settle_coverage_136_test.go:213-215`

```
	if tree.reads < 4 {
		t.Fatalf("precondition broken: only %d reads taken, half-and-half needs a window to lose in", tree.reads)
	}
```

Owning case: `TestCheckSettleHalfTheReadsFailedReportsItsLoss`, declared at
`sampler_settle_coverage_136_test.go:201` (`func TestCheckSettleHalfTheReadsFailedReportsItsLoss(t *testing.T) {`).

Two sibling guards in the same case fire on the same physical condition and are therefore shadowed by
`:213`, never reachable as the first red:

- `:216-218` `if kept < 2 {` - with `alternatingTree` parity, `kept = floor(reads/2)`, so `kept < 2` implies
  `reads < 4`, which `:213` already caught three lines earlier.
- `:221-223` `if lost < 2 {` - same implication from the other side.
- `:227-229` `if kept > lost+1 || lost > kept+1 {` - **cannot fire at all** for this fixture: `alternatingTree`
  (`:278-284`) alternates on `t.reads%2`, so `|kept-lost| <= 1` by construction, for every read count.

### 1(b) the sampling loop that feeds it

The test does not sample by itself; it opens one settle window and then reads the seam's counter:

- call site `internal/observe/sampler_settle_coverage_136_test.go:208`
  `rep, err := s.CheckSettle(context.Background(), SLOSleeping, 100*time.Millisecond, 10*time.Millisecond, true, TreeMetrics{}, 100<<20)`
  i.e. `within = 100ms`, `interval = 10ms`.
- fixture `internal/observe/sampler_settle_coverage_136_test.go:276-284` `type alternatingTree struct{ reads int }`
  whose `ReadTree` increments `t.reads` on every call (`:279`).
- the loop proper `internal/observe/sampler.go:491-526`:
  `:491` `startAt := time.Now()`, `:492` `deadline := startAt.Add(within)`, `:493` `t := time.NewTicker(interval)`,
  `:495` `for {`, `:499` `case <-t.C:` (the loop blocks here; the read only happens on a tick),
  `:501` `m, err := s.tree.ReadTree()`, `:523-525` `if time.Now().After(deadline) { break }`.

Structural fact worth carrying into the fix design: **`CheckSettle` has no read before the first tick.**
The sibling `SampleState` reads once up front (`internal/observe/sampler.go:269`
`windowStart, err := s.tree.ReadTree()`, error returns at `:271`) and again to close the window
(`:306` `windowEnd, err := s.tree.ReadTree()`), so its read count is `2 + ticks`. `CheckSettle`'s count is
purely `ticks delivered inside the budget`. Expected value here is about 10; the assertion's floor is 4.

### 1(c) every other case in `internal/observe/` that can emit a "window not fully measured" red

Defined as: the case fails because the window returned fewer reads/samples than that leg needs.
Line numbers are guards; the red sentence is 1-2 lines below each. All verified verbatim on disk.

Group A, same family as the flake (needs a *nonzero* number of reads; a partial window is enough to fire):

| # | file:line | guard, verbatim | window |
|---|---|---|---|
| 1 | `sampler_settle_coverage_136_test.go:125` | `if reads < 3 {` | `settleSUT` `:99` 100ms/10ms |
| 2 | `sampler_settle_coverage_136_test.go:213` | `if tree.reads < 4 {` | **AC#15 target**, `:208` 100ms/10ms |
| 3 | `sampler_settle_coverage_136_test.go:216` | `if kept < 2 {` | same window, shadowed by #2 |
| 4 | `sampler_settle_coverage_136_test.go:221` | `if lost < 2 {` | same window, shadowed by #2 |
| 5 | `sampler_settle_coverage_136_test.go:293` | `if reads < 3 {` | `settleSUT` `:292` 100ms/10ms |
| 6 | `sampler_settle_coverage_136_test.go:343` | `if reads < 3 {` | `settleSUT` `:342` 100ms/10ms |
| 7 | `sampler_settle_gate_136_test.go:201` | kept/lost floor test, see block below | `gateSUT` `:109` 200ms/10ms; `kept==0` at reads==1 |
| 8 | `sampler_test.go:99` | `if len(rep.Samples) < 3 {` | `SampleState` `:95` interval 20ms, duration 120ms |

Group B, fires only on a *fully starved* window (needs 1-2 reads; the practical failure mode is "no tick at all"):

| # | file:line | guard, verbatim | window |
|---|---|---|---|
| 9 | `sampler_settle_gate_136_test.go:156` | reads floor + all-kept test, see block below | 200ms/10ms |
| 10 | `sampler_settle_gate_136_test.go:239` | reads floor + none-kept test, see block below | 200ms/10ms, two `gateSUT` calls `:237`/`:238` |
| 11 | `sampler_test.go:336` | `if len(rep.Samples) == 0 {` | `SampleState` `:323` interval 10ms, duration 30ms - the tightest window in the package |
| 12 | `sampler_zerosample_136_test.go:63` | `if reads == 0 {` | `SampleState` `:52` 10ms/50ms |
| 13 | `sampler_zerosample_136_test.go:150` | `if len(rep.Samples) == 0 {` | `SampleState` `:146` 10ms/60ms |
| 14 | `sampler_settle_zerosample_136_test.go:62` | `if reads == 0 {` | `CheckSettle` `:54` 100ms/20ms |
| 15 | `sampler_settle_zerosample_136_test.go:141` | `if len(rep.Samples) == 0 {` | `CheckSettle` `:137` 100ms/20ms |

Verbatim text of the three guards marked "see block below" (pipes break a markdown table cell, so they
are quoted here instead; each line is `grep -n`-able as-is):

```
internal/observe/sampler_settle_gate_136_test.go:156:	if reads < 1 || len(rep.Samples) != reads {
internal/observe/sampler_settle_gate_136_test.go:201:	if kept < 1 || lost < 1 {
internal/observe/sampler_settle_gate_136_test.go:239:	if reads < 1 || len(none.Samples) != 0 {
```

Excluded on purpose (they are *surplus* or *exact-shape* assertions, not shortfall ones, and a
short window cannot make them red): `sampler_zerosample_136_test.go:60` and
`sampler_settle_zerosample_136_test.go:65` (`len(rep.Samples) != 0`, red when too MANY land);
`sampler_settle_coverage_136_test.go:129`, `:296` (`!= 1`, satisfied by the clamped script trees for any
reads >= 1); `:177`, `:224`, `:346`; `sampler_settle_gate_136_test.go:204`, `:242`.

Adjacent claim worth checking rather than trusting: `sampler_settle_gate_136_test.go:34-37` states
"every precondition is stated as 'at least', never as an exact read count, so these legs cannot join the
'the window only took N reads' family AC#15 is chasing". Verified **half-true**: the *at least* form is
correct, but "cannot join the family" is not - rows #7, #9, #10 above are in that family with a lower
floor. The gate file's windows are 200ms at 10ms (2x the target's budget), so its exposure per shot is
smaller, not zero. This should be reported to the ticket as a wording defect in AC#14's file header,
not as a code defect.
