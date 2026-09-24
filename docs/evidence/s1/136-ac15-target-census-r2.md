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


