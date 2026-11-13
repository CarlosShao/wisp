# 437b — the 313 union assertion now wants the same five rosters the door reports

- **Ticket**: `.scratch/wisp/issues/437-the-313-door-reports-seven-writable-fields-the-two-rosters-do-not-name.md`
- **Wave**: 2026-11-13, write leg. Write face for this leg = **`cmd/wisp/**`** (exclusive).
- **Dispatch**: `.scratch/wisp/probes/orch/20261113/txt/437b-313assert-dispatch.txt`
- **Predecessor**: `docs/evidence/s1/437-313-door-union.md` (read-only diagnosis, commit `018b3112d`).
  Its verdict is **taken as given and not re-derived**: the door is honest, the assertion is stale,
  and no `internal/panel` change is correct.
- **Change**: one construction, three added lines, in one file. `git diff --name-only` over the
  whole tree still shows other legs' pre-existing dirt; this leg's own delta is exactly
  `cmd/wisp/panel_config_313_test.go` (§7).

## 1. The five rosters, by their real function names

The door's report is computed by `acceptedWritableFields` (`internal/panel/config_write_313.go:176-183`),
which unions **five** accessors in this order:

| # | real function name | file:line | size |
|---|---|---|---|
| 1 | `WritableFields()` | `internal/panel/config_handlers.go:100` | 7 |
| 2 | `AdditionalWritableFields()` | `internal/panel/config_write_313.go:154` | 13 |
| 3 | `WritableFields187()` | `internal/panel/config_write_187.go:66` | 1 |
| 4 | `WritableFields65()` | `internal/panel/config_write_65.go:69` | 2 |
| 5 | `WritableFieldsWave()` | `internal/panel/config_write_wave.go:62` | 4 |

`7 + 13 + 1 + 2 + 4 = 27` — which is exactly the count the door reports (§2).

The dispatch said `config_write_313.go:176-183` names the rosters. It does, but only as **call
sites**; the accessor definitions live in the five files above, so the names were read off the
definitions rather than copied from the call sites.

## 2. Before and after — the expected-set construction

**Before** (two rosters; the stale side):

```go
want := append(append([]string{}, panel.WritableFields()...), panel.AdditionalWritableFields()...)
sort.Strings(want)
```

**After** (the same five rosters, in the door's own order):

```go
want := append(append([]string{}, panel.WritableFields()...), panel.AdditionalWritableFields()...)
want = append(want, panel.WritableFields187()...)
want = append(want, panel.WritableFields65()...)
want = append(want, panel.WritableFieldsWave()...)
sort.Strings(want)
```

Three things were deliberately **not** done:

- the comparison is still `!reflect.DeepEqual(reported, want)` — an **exact-set** comparison, not
  `Contains` and not a subset check;
- `sort.Strings(want)` is untouched, and `reported` is still compared in the door's emitted order
  (the door sorts its own report, so both sides stay sorted);
- half 2 and every other test are untouched. The only edit is the three `append` lines.

The failure sentence was left verbatim (`the two rosters say …`), because the dispatch bounds this
leg to *the construction of the expected set*. The wording is now slightly stale, but it is not
part of the comparison and changing it would widen the diff beyond the one point this leg owns.
That is a note for a later leg, not a change made here.

## 3. Whole-family rerun (all three `^Test313` cases, not just the pin)

Carrier only: `D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test313'"`.

The family is exactly three tests (`Test313EmptyConfigToStepfunRowChainAndPair`,
`Test313WritableUnionEqualsTheTwoRosters`, `Test313MultiElementCommaChainKeepsEveryElement`).

| run | PASS | FAIL |
|---|---|---|
| before the change | 2 | **1** (`Test313WritableUnionEqualsTheTwoRosters`, 0.00s) |
| after the change | **3** | **0** |

Before, verbatim:

```
--- FAIL: Test313WritableUnionEqualsTheTwoRosters (0.00s)
    panel_config_313_test.go:252: the door reports the union [audio_input_device ball_click_through
    … role_chat_thinking …]; the two rosters say [audio_input_device …] (A881: union == table1 + table2)
```

After:

```
=== RUN   Test313EmptyConfigToStepfunRowChainAndPair
--- PASS: Test313EmptyConfigToStepfunRowChainAndPair (0.04s)
=== RUN   Test313WritableUnionEqualsTheTwoRosters
--- PASS: Test313WritableUnionEqualsTheTwoRosters (0.00s)
=== RUN   Test313MultiElementCommaChainKeepsEveryElement
--- PASS: Test313MultiElementCommaChainKeepsEveryElement (0.03s)
PASS
```

## 4. Half 2 now verifies all twenty-seven

Half 2 (`panel_config_313_test.go:258-269`) iterates `want` and demands each member reach the leg
exactly once. Because `want` is now the five-roster union, half 2 covers **27** instead of 20.

Measured, not inferred — an out-of-repo copy added one `t.Logf("PROBE437B half2 covers %d union members", len(want))`
and the run read:

```
=== RUN   Test313WritableUnionEqualsTheTwoRosters
    panel_config_313_test.go:259: PROBE437B half2 covers 27 union members
--- PASS: Test313WritableUnionEqualsTheTwoRosters (0.00s)
```

`PASS` with zero per-member errors means all 27 are accepted by the door and each reaches the leg
once. So the widening does not buy a green by skipping seven fields — it pins seven more.

## 5. Mutation controls (both mandatory directions, both run)

Runtime-behaviour check ⇒ `go test -overlay` with the copies held **outside** the working tree is a
valid instrument. Every copy below lives in `D:\tmp\wisp\20261113\437b-313assert\`; the repository
was never mutated for any of them, and `GOFLAGS` was cleared before the restore run.

| id | mutation | reading |
|---|---|---|
| **M1** (mandatory) | `want` put back to the two-roster construction | **FAIL** — the assertion goes red again, at `panel_config_313_test.go:253`, naming all seven extras on the door side |
| **restore** | the tree's five-roster version, no overlay | **PASS** — 3/3 family green |
| **M5** (extra teeth) | five-roster `want`, but the door's report truncated by dropping the wave join from `acceptedWritableFields` | **FAIL** — names exactly the four wave fields (`ball_click_through`, `ball_opacity_idle`, `panel_font_size`, `panel_width`) as missing |

M1 is the required control: the old two-roster expectation is red, so the change is what makes the
pin pass, not an unrelated drift. M5 shows the comparison kept its tooth in the other direction —
an exact-set comparison still catches a report that is one roster short, which a subset or
`Contains` check would have silently accepted.

Note the overlay is demonstrably live rather than a no-op: M1 and M5 both changed the verdict, and
each produced the specific message the mutation predicts.

## 6. Gates (AC#3)

| gate | command | reading |
|---|---|---|
| build | `go build ./cmd/wisp/` | exit **0**, no output |
| vet | `go vet ./cmd/wisp/` | exit **0**, no output |
| format | `D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp/panel_config_313_test.go` | **empty** (0 lines), exit 0 |
| D22 | `D:\work\soft\Git\bin\bash.exe -c "scripts/d22scan.sh"` | `d22scan: clean - no D22 ban violations`, exit 0 |

`build/mut_rvw.go` was not deleted and is still present. `go build ./cmd/wisp/` was used, not
`go build ./...`.

## 7. Boundary (AC#2)

`git diff --name-only` on this tree:

```
.scratch/wisp/issues/426-…md            (deleted, other leg)
.scratch/wisp/issues/427-…md            (deleted, other leg)
.scratch/wisp/issues/428-…md            (deleted, other leg)
.scratch/wisp/issues/429-…md            (deleted, other leg)
.scratch/wisp/issues/430-…md            (deleted, other leg)
.scratch/wisp/issues/431-…md            (deleted, other leg)
.scratch/wisp/issues/432-…md            (deleted, other leg)
.scratch/wisp/issues/433-…md            (deleted, other leg)
.scratch/wisp/issues/434-…md            (deleted, other leg)
.scratch/wisp/issues/435-…md            (deleted, other leg)
.scratch/wisp/issues/436-…md            (deleted, other leg)
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go   (other leg)
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt              (other leg)
cmd/wisp/panel_config_313_test.go       <-- this leg's only code change
frontend/fixtures/composer-states.html  (other leg)
```

The other entries are pre-existing dirt from other legs, not this leg's writes — proven by
modification time, since this session began at **2026-11-13 09:43**:

| file | last write |
|---|---|
| `cmd/wisp/panel_config_313_test.go` | **2026-11-13 09:44:40** (this leg) |
| `frontend/fixtures/composer-states.html` | 2026-11-11 16:59:50 |
| `.scratch/wisp/probes/291/mic/…_windows_test.go` | 2026-11-09 22:01:05 |
| `.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt` | 2026-11-12 13:39:05 |

`git diff --name-only -- internal/panel/` is **empty**: the panel package is untouched, which is
the point of the predecessor's measurement (M2/M6 duplicate names and red three disjointness pins;
M3 makes the refusal text omit an accepted name and reds four pins).

Nothing under `docs/PLAN.md`, `docs/specs/**`, `internal/observe/thresholds.go`, golden files,
`tools/d22scan/allowlist.txt`, `frontend/**` or `design/**` was touched by this leg. No AC box
flipped. No push. `git add`/`git commit` both used an explicit pathspec.
