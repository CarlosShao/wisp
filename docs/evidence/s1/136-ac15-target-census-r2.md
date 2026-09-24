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

---

## 2. What changed under it

Section 2 anchor: `git rev-parse HEAD` = `5889559b767586ab73dd859f0a5f00e208612f49`
(moved past section 1's `a1fd5bf` because other agents landed docs in between; `5889559` verified with
`git cat-file -t`, which returns `commit`. Every sha cited in this file was checked the same way.)

Range asked for: `git log --follow -p ca2c55e..HEAD -- internal/observe/sampler_settle_coverage_136_test.go`.
`ca2c55e` = `docs(A179,A180): 128 AC#4 收＋AC#5 追加；136 预检四裁定入库`, verified `git cat-file -t` = `commit`.

Three commits touch this file in the range, and one more touches the package:

| sha | subject (trimmed) | numstat for this file | touched this leg? |
|---|---|---|---|
| `aef82f5` | feat(136 AC#14): settle 报告自陈门行 + 六枚探针 | +30 / -2 | no - it added the *sibling* file `sampler_settle_gate_136_test.go` (0 -> new) and the `Verdicts`/gate machinery; in this file it edited the header comment of the three *other* legs |
| `52191ce` | feat(136 AC#14 Gate): flip settleCoverageRowGates true + rewrite 2 coverage legs | +52 / -4 | yes - one hunk lands **inside** `TestCheckSettleHalfTheReadsFailedReportsItsLoss`, but only below the guards |
| `f9bc512` | docs(136 AC#14b 收口#1): 改写四处过期注释 | +7 / -7 | comments only |
| `c28d4e8` | docs(136 收口#2): sampler.go:588 第 5 处过期注释 | not in this file | comments in `sampler.go` only |

Proof the leg's own skeleton is untouched, not an inference from the diff being small:

```
git show ca2c55e:internal/observe/sampler_settle_coverage_136_test.go | sed -n '176,206p' > /tmp/old_half.txt
sed -n '201,231p' internal/observe/sampler_settle_coverage_136_test.go > /tmp/new_half.txt
diff /tmp/old_half.txt /tmp/new_half.txt      # empty
```

That 31-line span is `func TestCheckSettleHalfTheReadsFailedReportsItsLoss` through the `kept < 2`
guard: declaration, fixture, `ReleaseMemory()`, the `CheckSettle` call with `100ms/10ms`, `tree.reads < 4`,
`kept < 2`. It is byte-for-byte identical and simply shifted 25 lines down. `sampler.go`'s settle read
loop is likewise untouched in the range: `git diff ca2c55e..HEAD -- internal/observe/sampler.go` yields
exactly two hunks, `@@ -452,7 +452,19 @@` (the `SettleReport` doc + `Verdicts` field) and
`@@ -514,9 +526,72 @@` (the fold call at what is now `:537`, then `settleCoverageMetric`,
`settleCoverageRowGates`, `buildSettleVerdicts`, `foldSettlePass` appended after the function). Nothing in
`491-526` changed. **The mechanism that produces the flake is therefore the same machine that produced it
when AC#15's face was written; only its address moved.**

### Does the "half the reads failed" leg now depend on the gate row?

It depends on the gate row for its *assertion*, not for its *flake*. Precisely:

- Old form (`ca2c55e:208-210`, the same address the ticket's AC#14 cell names):
  `if !rep.Pass { t.Fatalf("disclosure leg, not a verdict leg: this window
  still passes, report=%+v", rep) }`. This required `rep.Pass == true`, i.e. it required
  `backInTime && memOK && releaseOK` **and** (post-`aef82f5`) whatever the fold let through. Had that line
  survived the gate flip it would be **deterministically red**, because at `settleCoverageRowGates = true`
  a half-covered window's pass bit is now vetoed by its own row. So the flip and the rewrite are one
  indivisible change; that is what `52191ce` did.
- New form (`HEAD:244-263`): looks up `rep.Verdicts[i].Metric == "sampling"` and asserts
  `coverage != nil`, `!coverage.Pass`, and that the red row prints `sample_errors`. `buildSettleVerdicts`
  sets `covered := rep.SampleErrors == 0 && len(rep.Samples) > 0` (`sampler.go:565`), so with
  `SampleErrors >= 1` the row is `Pass: false` **for either value of the constant**. The last block
  (`:259-263`, `if v.Gate && !v.Pass && rep.Pass`) is the only gate-sensitive line, and it is an
  invariant that is vacuous at `false` and enforced at `true`.

So: behaviour at `settleCoverageRowGates = true` vs. at `false` differs for exactly one assertion in this
leg, and that assertion is an invariant, not a count. **Nothing in the leg changed its timing exposure.**
Two secondary consequences do matter for AC#15:

1. The old `!rep.Pass` line was itself weakly timing-coupled (`BackWithinCapMS >= 0` needs a *kept* read).
   The replacement `coverage.Pass` needs only one *dropped* read, which the alternating fixture produces
   on its very first call. The rewrite therefore *narrowed* the leg's timing surface rather than widening
   it, and the sole remaining timing-coupled assertions in it are the group-A guards listed in section 1.
2. Exposure per whole-package shot roughly doubled since the 1-in-240 measurement, because `aef82f5` added
   `sampler_settle_gate_136_test.go`, which does not exist at `ca2c55e`
   (`git ls-tree ca2c55e -- ...gate_136_test.go` -> empty) and calls `gateSUT` 7 times
   (`:155`, `:199`, `:237`, `:238`, and 3 iterations of the `scripts` loop at `:313-322`), each one a
   200 ms ticker window (`:109`). Counting helper *call sites* rather than definitions, the package's
   ticker-window budget per full-package run goes from about 1560 ms at `ca2c55e` to about 2960 ms at
   `5889559`, of which the AC#15 target itself contributes 100 ms in both. Consequence for section 5: the
   historical `1/240` was a per-package-run rate on a package that spent half as much wall time inside
   sampling windows. Reusing it as a prediction for the current tree would understate the rate; the
   denominator also now has 6 more cases in it, so the roster diff in section 5 is mandatory.

The arithmetic behind those two numbers, spelled out because the obvious instrument gets it wrong twice:

| file | windows per full-package shot | each | ms |
|---|---|---|---|
| `sampler_test.go` | 11 `SampleState`/`CheckSettle` calls (`:95 :134 :153 :168 :182 :203 :234 :273 :286 :307 :323`) | 30-120 | 750 |
| `sampler_settle_coverage_136_test.go` | 4 (3 via `settleSUT` at `:124 :292 :342`, 1 inline at `:208`) | 100 | 400 |
| `sampler_settle_zerosample_136_test.go` | 2 (`:54 :137`) | 100 | 200 |
| `sampler_zerosample_136_test.go` | 2 (`:52 :146`) | 50, 60 | 110 |
| `observer_cost_test.go` | 1 (`:75`) | 60 | 60 |
| `sampler_faketree_guard_136_test.go` | 1 (`:101`) | 40 | 40 |
| `sampler_settle_gate_136_test.go` | 7 (direct `:155 :199 :237 :238` plus 3 iterations of the loop calling `:322`) | 200 | 1400 |

`ca2c55e` = everything except the last row = **1560 ms**. `HEAD` = all rows = **2960 ms**. Two traps in
getting there mechanically, both hit while producing this table: (i) `grep -c "settleSUT(t,"` and the
helper's own definition line both match the `CheckSettle(context...` text, so a naive sum double-counts
each helper once (it inflates coverage by 100 ms and the gate file by 200 ms); (ii) `grep -c` counts
*lines*, not *executions*, so the 3-iteration `scripts` loop at `:313-322` reads as one window instead of
three. The naive pipeline yields 1660 / 2860; the corrected per-execution tally is 1560 / 2960. Quote the
table, not the pipeline, and note that only the ratio matters for the conclusion.

---

## 3. The two dimensions, measured separately

Section 3 anchor: `git rev-parse HEAD` = `5889559b767586ab73dd859f0a5f00e208612f49`.

The mistake on record tonight was writing one composite number ("4 sites in 2 `.go` files", real answer
4 sites in 3 files) as if a single command could produce both factors. It cannot: a site count and a file
count come from two different pipes. Both are therefore run twice below, once per definition of "family".
Quoted regexes are single-quoted in the shell so `$` and `[` stay literal.

### 3.1 The broad family (every shortfall guard in `internal/observe/`, section 1 lists them all)

Sites:

```
grep -rn -E "^[[:space:]]*if (reads|tree\.reads|kept|lost|len\(rep\.Samples\)|len\(none\.Samples\))( < | == 0)" \
  --include=*_test.go internal/observe/ | wc -l
# 15
```

Files:

```
grep -rn -E "^[[:space:]]*if (reads|tree\.reads|kept|lost|len\(rep\.Samples\)|len\(none\.Samples\))( < | == 0)" \
  --include=*_test.go internal/observe/ | cut -d: -f1 | sort -u | wc -l
# 5
```

The file dimension is not derivable from the site dimension, and `uniq -c` on the same stream is what
shows the skew (6 + 3 + 2 + 2 + 2 = 15 over 5 files):

```
internal/observe/sampler_settle_coverage_136_test.go        6
internal/observe/sampler_settle_gate_136_test.go            3
internal/observe/sampler_settle_zerosample_136_test.go      2
internal/observe/sampler_test.go                            2
internal/observe/sampler_zerosample_136_test.go             2
```

So the composite reading of the broad family is **15 sites in 5 files**, and the two numbers must be
quoted with their definitions attached, never as "15x5".

### 3.2 The narrow family (the one AC#15 is about: a guard that can fire while the window still took
at least one read, i.e. a *partial* window is enough)

Sites:

```
grep -rn -E "^[[:space:]]*if (tree\.reads < [2-9]|reads < [2-9]|kept < [2-9]|lost < [2-9]|kept < 1 [|][|]|len\(rep\.Samples\) < 3)" \
  --include=*_test.go internal/observe/ | wc -l
# 8
```

Files:

```
grep -rln -E "^[[:space:]]*if (tree\.reads < [2-9]|reads < [2-9]|kept < [2-9]|lost < [2-9]|kept < 1 [|][|]|len\(rep\.Samples\) < 3)" \
  --include=*_test.go internal/observe/ | wc -l
# 3
```

**8 sites in 3 files**, and the 8 are exactly section 1's group A rows 1-8. The third file beyond the two
`136` settle files is `sampler_test.go:99` - which is why "2 `.go` files" was never going to be right for
any census of this family that looked at only the two ticket-named files.

Two honesty notes about this instrument, because it is regex-shaped and a regex has no idea what the
fixture does:

- `sampler_settle_gate_136_test.go:201` is in the narrow family only via the `kept < 1` disjunct: with the
  alternating tree the first read is always a failure, so `reads == 1` already makes it red. A pure
  threshold-on-`reads` regex misses it; a pure "message contains reads" regex double-counts the shadowed
  `:216`/`:221` pair as independent flake surfaces. The hand classification in section 1 is what carries
  the meaning; the two commands above only prove the counts are reproducible.
- Guards whose predicate goes red on *surplus* rather than shortfall
  (`sampler_zerosample_136_test.go:60`, `sampler_settle_zerosample_136_test.go:65`, both
  `len(rep.Samples) != 0`) are deliberately outside both families. Including them would raise the site
  count to 17 and the file count stays 5, which is exactly the kind of silent composite drift this
  section exists to prevent; anyone re-running the census should say which predicate they meant.

---

## 4. Mechanism candidates, ranked, with the code that makes each possible

Section 4 anchor: `git rev-parse HEAD` = `77eac9f` (my own section 3 commit; HEAD moves under this run
because other agents are committing docs, which is why section 7 re-proves `internal/observe/` is
byte-unchanged since section 1's `a1fd5bf` rather than trusting any single anchor).

First, what the test does **not** do, because the answer decides the fix: it does not sleep, it does not
poll, it does not use a channel with a timeout, and it does not call a sampler with a fixed interval it
controls. It opens one real `time.Ticker` window and then asserts a count at the end of it. Everything
below follows from that single fact.

**C1 (primary). The settle window's read count is produced purely by ticker delivery, with no leading
read and no lower-bound criterion.**
Code: `internal/observe/sampler.go:492-500` -
`deadline := startAt.Add(within)`, `t := time.NewTicker(interval)`, `for { select { case <-ctx.Done(): ...
case <-t.C:` and the read at `:501` is reachable *only* from a tick. Contrast `SampleState`, which reads
once before the loop (`:269`) and once to close it (`:306`), so it can never report a zero-read window.
For this leg (`:208`, `within=100ms`, `interval=10ms`) the realised count is about 10 and the assertion's
floor at `:213` is 4.
What has to be true on the machine: the runtime must deliver fewer than 4 ticks of a 10ms ticker inside a
100ms monotonic budget - either each wake-up running about 25ms late, or one wake-up roughly 70ms late
(the loop is serial: read, check deadline, block again). Any of a GC stop-the-world, a long syscall on the
timer goroutine's path, or the OS not giving the process a thread for that span does it. It is not
"load" as an uncaused cause; it is that the leg asserts a *count* whose only guarantee is a *duration*.

**C2 (the actual AC#15 ② root cause, stated as the missing criterion).**
Code: the pair that never meet - `sampler.go:493` (`time.NewTicker(interval)`) with `sampler.go:523`
(`if time.Now().After(deadline) { break }`) on one side, and
`sampler_settle_coverage_136_test.go:213` (`if tree.reads < 4 {`) on the other.
Between the window and the assertion there is no statement of the form "this window is allowed to be
judged once N reads have been observed, and must be declared broken if N have not arrived before
deadline T". Nothing on the machine has to be true for this to be the defect: the code as written would
still be missing it under a scheduler that behaved perfectly, because the guarantee `100ms / 10ms => at
least 4 reads` is not something the program establishes anywhere. Fix shape is dictated by AC#15 ③: the
test must wait *on a criterion* (poll `tree.reads` until it reaches the floor, bounded by a monotonic
`observe.Timeout`, `internal/observe/clock.go:25-52`), rather than wait a fixed time and hope.

**C3 (amplifier, not a cause). Windows timer granularity feeding a 10ms ticker.**
Code: same `:493`. If the effective system tick is 15.6ms rather than 1ms, the expected read count for a
100ms budget drops from about 10 to about 6, which is 2 shots from the floor of 4 instead of 6.
What has to be true: the platform's timer resolution being coarse for that process. This cannot be the
whole story - a systematically coarse clock would make the leg red far more often than 1 in 240 - but it
sets how much headroom C1 has. It is also the reason "just raise the interval" is not a fix: it moves the
mean, it does not add the missing criterion.

**C4 (plausible trigger for the single observed shot). A stop-the-world pause inside the window.**
Code: the window is entered right after `coverage:207` `ReleaseMemory()` and the leg's own allocations;
`Sample.At` is stamped with `time.Now()` per read (`sampler.go:512`) and the package's other cases keep
the heap moving. A GC pause longer than about 70ms freezes tick delivery for exactly the length C1 needs.
What has to be true: one long enough STW inside that 100ms. Not reachable by reading code alone, and
section 5 is what decides it: `-count=N` in one process makes GC pressure per shot much higher than a
single whole-package shot, so the two denominators should separate this candidate if it is real.

**C5 (checked and excluded by code, list them so nobody re-litigates them).**
- The pre-window release is not a stall: `sampler.go:416-419` shows `ReleaseMemory()` only does
  `freeOSMemoryCalls.Add(1)` plus `debugFreeOSMemory()`, and the leg replaces that with a no-op at
  `coverage:205-206` (`debugFreeOSMemory = func() {}`, restored by `t.Cleanup`).
- Cancellation is not a cause: the `<-ctx.Done()` branch (`sampler.go:497-498`) cannot fire, the leg
  passes `context.Background()` (`:208`).
- Parallel interference is not a cause: `grep -rc "t.Parallel()" --include=*_test.go internal/observe/`
  sums to 0, so cases cannot interleave; note in passing that the shared-global dance on
  `debugFreeOSMemory` (`:205-206`, `:96-97`, `gate:105-107`) is safe *only* because of that, so any future
  `t.Parallel()` in this package turns into a different flake. This census changes nothing, just records
  the dependency.
- A leaked busy goroutine from another case is not a cause: no `TestMain` in the package
  (`grep -rn "func TestMain" --include=*.go internal/observe/` -> 0), and the only goroutines the tests
  start are `logging_test.go:273-279` behind a `sync.WaitGroup` and the short-lived workers of
  `goroutine_test.go`, which is itself the file that documents "no sleep is used to paper over the
  window" (`goroutine_test.go:24`).
- Parity is not a cause: `:227` `if kept > lost+1 || lost > kept+1` cannot fire for `alternatingTree`
  (`:278-284`) at any read count, so the half-and-half property is fixture-guaranteed and only the
  *count* floors are live.

---

## 5. What the real measurement program must run (this census ran none of it)

Section 5 anchor: `git rev-parse HEAD` = `10612382fa6dae71ddd0650d8ca7d09a71e47bb0` (this census's own
section 4 commit). `internal/observe/` is byte-identical to section 1's anchor `a1fd5bf` throughout
sections 3-5; proof command and re-read are in section 7.
Preconditions the runner must satisfy and record before shot 1, not afterwards: the self-hosted runner is
idle (`gh run list --limit 20` shows nothing in flight on this machine, and `slo-full` in particular is
not running), `git rev-parse HEAD` recorded, `git diff <that anchor>..HEAD -- internal/observe/` empty at
the end of every batch, and one pristine shot whose full `=== RUN` roster is saved as the golden roster.

### 5.1 Denominator A - whole package, one process per shot

```
mkdir -p /tmp/ac15/A
for i in $(seq 1 1440); do
  go test ./internal/observe/ -v -count=1 > /tmp/ac15/A/shot-$(printf %04d $i).log 2>&1
done
```

Rate A = (shots whose log contains the red sentence) / 1440, per shot, with:

```
grep -l "precondition broken: only [0-9]* reads taken" /tmp/ac15/A/shot-*.log | wc -l
```

One shot is one package process, so this is the denominator the historical `1 red in 240 whole-package
-v runs` was measured on and it is the only one comparable to it.

### 5.2 Denominator B - same process, repeats

```
go test ./internal/observe/ -run '^TestCheckSettleHalfTheReadsFailedReportsItsLoss$' -v -count=1200 \
  > /tmp/ac15/B-target.log 2>&1
go test ./internal/observe/ -run '^TestCheckSettle' -v -count=120 > /tmp/ac15/B-family.log 2>&1
```

Rate B = (occurrences of the red sentence) / (occurrences of `--- PASS` plus `--- FAIL` for that name),
read out of the same log:

```
grep -c -- "--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss" /tmp/ac15/B-target.log
grep -c -- "--- PASS: TestCheckSettleHalfTheReadsFailedReportsItsLoss" /tmp/ac15/B-target.log
```

**A and B are two rates and must be written as two numbers. They are never added, never averaged, and
never reported as one percentage** (AC#11 already pinned this, and the reason is physical here: B repeats
the same case in one process, so its GC heap, timer heap and page-cache state at shot 500 are not the
state of a fresh process at shot 1 - see C4 in section 4. B is the *mechanism* probe, A is the *rate*).
The 100 ms window means B-target's 1200 shots cost about 2 minutes of sampling, which is why B is where
the count is cheap; A is where the number is honest.

### 5.3 How many shots is defensible for a 1-in-240 event

`1/240` is a sighting, not a rate: with 1 hit in 240 shots the 95% interval is roughly 0.01% to 2.3%, so
it is compatible with anything from 1/10000 to 1/44. Two different goals, two different n:

- To **estimate** the rate to within a factor of about 2 you need at least 5 expected hits, i.e.
  `n = 1200` shots of A, and the run must report the hit count with its interval, not a bare percentage.
- To **prove the fix** in the sense of "no longer visible at the historical rate", use the rule of three:
  after `k = 0` hits the 95% upper bound is `3/n`. `3/720 = 1/240` is exactly the historical value, so 720
  buys nothing. The bound only clears the historical rate from below at `n > 720`; run **`n = 1440`** for
  an upper bound of `0.21%`, i.e. strictly tighter than 1/240, and state it as a bound rather than as "0
  hits so it is fixed".
- AC#15 ④'s `>= 30 shots with 0 hits` is kept as the **entry ticket, not the verdict**: 30 clean shots only
  bound p below 10%, which cannot separate "fixed" from "1/240 still sitting there". A submission that
  brings 30 shots satisfies ④ literally and still fails ① and ② in substance, so the acceptance table
  should demand ④'s 30 (both denominators, each batch's per-shot readings left in the evidence, named
  individually) plus the 1440-shot A bound above.

### 5.4 The roster set-difference, and why it is not optional

A case that *panics* takes the test binary down with it, so every case scheduled after it emits no `RUN`
and no `PASS` line at all: a shrunken denominator looks identical to a clean 0 hits. This package has
measured experience with it: `internal/observe/sampler_test.go:330-335` records the pre-guard shape as
"an empty `rep.Samples` panicked here, and the panic killed the test binary, so every case scheduled after
this one stopped reporting at all (measured: 52 RUN / 47 PASS / 5 FAIL + panic, 4 cases swallowed)".
So per shot, and against the golden roster:

```
grep -aoE '^=== RUN  +[^ ]+' $log | awk '{print $3}' | sort -u > /tmp/ac15/run.txt
grep -aoE '^ *--- (PASS|FAIL): [^ ]+' $log | awk '{print $2, $3}' | sort -u > /tmp/ac15/res.txt
comm -23 /tmp/ac15/run.txt /tmp/ac15/reported.txt   # ran but never reported: must be empty
comm -13 /tmp/ac15/run.txt /tmp/ac15/golden-run.txt # in golden but never ran: must be empty
grep -c -- '--- SKIP' $log                          # must be 0 (this repo: skip launders "untested" into "passed")
```

Both directions, every batch, plus the count of `--- FAIL` names compared against the historical red list
so "0 hits" cannot be smuggled in as "the target case silently stopped existing".

On "does any test in that package still index a slice without a length guard": yes, and not in the
`Samples` path.

- `internal/observe/sampler_test.go:340` `if rep.Samples[0].RuntimeGoroutines < 1 {` is guarded two lines
  above at `:336` `if len(rep.Samples) == 0 {` - that is AC#8's fix, in place and correct.
- `internal/observe/earlylog_130_test.go:130-135` (`recs := sink.snapshot()`, then `recs[1]`, `recs[0]`,
  `recs[2]`) and `:232` (`got := sink.snapshot()[0]`) have **no length guard of their own** on
  `recs`/`got`. They are safe today only transitively, via `flushed != 3` at `:123`, `len(got) != 3` at
  `:127` and `flushed != 1` at `:229`, plus the fact that `msgs()` (`:57-59`) is built by ranging over
  `snapshot()`. If `drain`'s contract ever diverges from the sink's contents these become panics, and they
  sit in the same binary as the settle family, so they would eat the rest of the roster. Not a bug today;
  described, not fixed, per this census's rules.
- The one that matters for a mutation-based fix check:
  `internal/observe/sampler_settle_gate_136_test.go:262`
  `never := buildSettleVerdicts(SettleReport{})[0]` indexes the builder's result with no length guard. Any
  mutation that makes `buildSettleVerdicts` return an empty slice turns that into a panic, which will show
  up as a *shrunken roster*, not as a readable red. Whoever runs the fix must either expect that or read
  the roster diff before concluding anything from those shots.

### 5.5 Also required, because section 2 changed the exposure

Re-run 5.1's A batch once on the pre-fix anchor (the commit before the fix lands) and once on the post-fix
HEAD, both at `n = 1440`, and report the two bounds side by side. Comparing a post-fix number against the
historical `1/240` alone is not a comparison: as recorded in section 2, the package's ticker-window budget
per whole-package shot has roughly doubled since those 240 shots were counted (about 1560 ms then, about
2960 ms now), so the pre-fix rate on the current tree is the only admissible baseline.




