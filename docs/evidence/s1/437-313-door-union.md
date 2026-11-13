# 437 — the 313 door's union report vs the two rosters

- **Ticket**: `.scratch/wisp/issues/437-the-313-door-reports-seven-writable-fields-the-two-rosters-do-not-name.md`
- **Wave**: 2026-11-13, write leg. Write face for this leg = `internal/panel/**`.
- **Verdict**: **the door is honest; the assertion is stale.** All seven differing names have a
  real write branch in the leg. **No `internal/panel` change is correct**, and every
  `internal/panel`-only repair I could build turns a green suite red (M2/M3/M6, §5).
  The fix belongs on the assertion side (`cmd/wisp/panel_config_313_test.go`), which is
  **outside this wave's write face** ⇒ reported for a later wave, not changed here.
- **Reading**: 27 reported vs 20 in the two rosters = **7**, not the ticket's "26 / 6".

## 1. Reproduction (verbatim, via the carrier)

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test313WritableUnionEqualsTheTwoRosters$'"
```

```
=== RUN   Test313WritableUnionEqualsTheTwoRosters
    panel_config_313_test.go:252: the door reports the union [audio_input_device ball_click_through ball_opacity_idle ball_size ball_style model_context_window model_price_in model_price_out model_row_create model_row_delete panel_font_size panel_width provider_api_key_ref provider_base_url provider_billing provider_credential provider_row_create provider_row_create_custom provider_row_delete role_chat_model role_chat_pair role_chat_thinking voice_aec_enabled voice_cloud_asr_chain voice_cloud_tts_chain voice_realtime_model voice_realtime_provider]; the two rosters say [audio_input_device model_context_window model_price_in model_price_out model_row_create model_row_delete provider_api_key_ref provider_base_url provider_billing provider_credential provider_row_create provider_row_create_custom provider_row_delete role_chat_model role_chat_pair voice_aec_enabled voice_cloud_asr_chain voice_cloud_tts_chain voice_realtime_model voice_realtime_provider] (A881: union == table1 + table2)
--- FAIL: Test313WritableUnionEqualsTheTwoRosters (0.00s)
FAIL
```

0.00s ⇒ stable, load-independent.

## 2. My own re-count (the ticket's number is wrong)

Counted mechanically by parsing the red line out of the run log, not by eye:

| side | count |
|---|---|
| the door reports | **27** |
| the two rosters (`WritableFields()` + `AdditionalWritableFields()`) | **20** |
| **difference** | **7** |
| roster-only (names the tables have and the door omits) | **0** |

The ticket says 26 in the headline and body, then lists seven names and flags its own
arithmetic as uncertain. **The list was right and the count was wrong**: 27 − 20 = 7, and the
door list carries no duplicates (checked by `Group-Object`). The ticket's "6" is the error.

The difference, by name (all seven door-only, sorted):

1. `ball_click_through`
2. `ball_opacity_idle`
3. `ball_size`
4. `ball_style`
5. `panel_font_size`
6. `panel_width`
7. `role_chat_thinking`

## 3. Arithmetic closure — where the seven come from

The door's report is not a two-roster union; it is a **five-roster** union
(`acceptedWritableFields`, `internal/panel/config_write_313.go:176-183`), while the assertion
reads only the first two:

| roster | accessor | size | declared at |
|---|---|---|---|
| 248 | `WritableFields()` | 7 | `config_handlers.go:91-94` |
| 313 | `AdditionalWritableFields()` | 13 | `config_write_313.go:114-128` |
| 187 | `WritableFields187()` | 1 | `config_write_187.go:45-47` |
| 65 | `WritableFields65()` | 2 | `config_write_65.go:55-58` |
| wave | `WritableFieldsWave()` | 4 | `config_write_wave.go:44-49` |

`7 + 13 = 20` (what the assertion wants); `7 + 13 + 1 + 2 + 4 = 27` (what the door reports).
The seven-name difference is **exactly the three rosters that landed after ticket 313**:
1 (187) + 2 (65) + 4 (wave) = 7.

The door's own gate asks all five: `knownWritableField` (`config_handlers.go:465-469`).

## 4. Per-field: does the leg have a real write branch?

Judged against the dispatcher's write switch, `configManager.ApplySetting`
(`cmd/wisp/panel_config_store.go:211`), whose default branch refuses
(`:390-391`, `字段 %q 没有对应的写入腿，未落盘`). Every one of the seven has a real `case`
that calls a real config setter — none falls to the default:

| field | leg branch (`file:line`) | calls | config setter (`file:line`) |
|---|---|---|---|
| `ball_click_through` | `cmd/wisp/panel_config_store.go:354` | `SetBallClickThrough(on)` `:364` | `internal/config/settings_appearance_wave.go:69` |
| `ball_opacity_idle` | `cmd/wisp/panel_config_store.go:340` | `SetBallOpacityIdle(v)` `:353` | `internal/config/settings_appearance_wave.go:56` |
| `ball_size` | `cmd/wisp/panel_config_store.go:326` | `SetBallSize(v)` `:339` | `internal/config/settings_413.go:45` |
| `ball_style` | `cmd/wisp/panel_config_store.go:321` | `SetBallStyle(w.Value)` `:325` | `internal/config/settings_65.go:46` |
| `panel_font_size` | `cmd/wisp/panel_config_store.go:377` | `SetPanelFontSize(v)` `:386` | `internal/config/settings_appearance_wave.go:95` |
| `panel_width` | `cmd/wisp/panel_config_store.go:365` | `SetPanelWidth(v)` `:376` | `internal/config/settings_appearance_wave.go:80` |
| `role_chat_thinking` | `cmd/wisp/panel_config_store.go:262` | `SetRoleChatThinking(w.Value)` `:267` | `internal/config/settings_187.go:104` |

For completeness I mapped **all 27** reported names to a leg case mechanically:
**27 reported names / 27 `case panel.Field…` branches / 0 without a branch.** The door lists
nothing it cannot write.

### The seven are not merely listed — they are accepted end to end

The test's own half 2 (`panel_config_313_test.go:255-266`) sends every union member through
`disp.Handle` and demands it reach the leg **exactly once**. I extended that loop to all
27 by overlay (M4, §5): **PASS**, 0 refused, every field reaching the leg once. So the seven
are not "reported but unwritable" — they clear the door, pass the selector/value checks, and
land on a real setter.

## 5. Mutation controls (all run; overlays live outside the repo, `D:\tmp\wisp\20261113\437-313\`)

Runtime-behaviour check ⇒ `go test -overlay` with copies held **outside** the working tree is a
valid instrument here (contrast ledger `A1804`, where a disk-reading static check made an
overlay useless). The repo was never written to for any of these.

| id | mutation | instrument | reading |
|---|---|---|---|
| **M1** | assertion `want` also appends 187 + 65 + wave | carrier, this pin | **PASS** — the report and a five-roster `want` agree exactly |
| **M2** | ticket's literal 甲: add the seven names to `fields313` | carrier, this pin | **FAIL** — report grows to 34 with seven names **doubled**; the door's union is five rosters, so adding to one roster duplicates |
| **M2** | same overlay | `go test ./internal/panel/` | **FAIL** — 3 disjointness pins red: `TestWaveRosterJoinsTheUnionAndStaysDisjoint` (`appearance_wave_test.go:54`), `Test65RosterJoinsTheUnionWithoutMovingOlderPins` (`ball_style_65_test.go:51`), `Test187RosterIsDisjointAndTheUnionAsksAllThree` (`config_write_187_test.go:116`) |
| **M3** | shrink the report to the two rosters (`acceptedWritableFields` drops 187/65/wave) | carrier, this pin | **PASS** (the assertion is satisfied by lying) |
| **M3** | same overlay | `go test ./internal/panel/` | **FAIL** — 4 pins red, all `acceptedWritableFields() lacks …`: wave `:67`, 422 `:54`, 65 `:55`, 187 `:134` |
| **M4** | five-roster assertion **and** half-2 loop over all 27 | carrier, `-test.run '^Test313'` | **PASS** — all three Test313 cases green |
| **M5** | five-roster assertion **+** report truncated by one roster | carrier, this pin | **FAIL** — `panel_config_313_test.go:255` names exactly the four missing wave fields ⇒ the proposed shape keeps its tooth |
| **M6** | collapse 187+65+wave into `AdditionalWritableFields()` (the one `internal/panel`-only way to make `W()+Add()` equal the report) | carrier, this pin | **FAIL** — same duplication as M2 |
| **M6** | same overlay | `go test ./internal/panel/` | **FAIL** — same 3 disjointness pins |

**The M2/M3 pair is the decisive one.** M2 is the ticket's tempting 甲 fix: it does not merely
fail to help, it makes the report *dishonest in a new way* (duplicate names). M3 shows the
reverse: silencing the assertion by shrinking the report is exactly the "报告撒谎" the ticket
forbids, and `internal/panel`'s own four pins catch it. **There is no correct edit inside
`internal/panel`.** The disagreement is not in the tables and not in the door; it is in the
assertion's `want` expression, which was written when only two rosters existed.

## 6. `role_chat_thinking` — the separate question

**It arrived with a different change, and it is fully legitimate.** Provenance:

- `git log -S 'role_chat_thinking'` ⇒ born in **`007de34f0` (2026-11-07)**, ticket 187's
  landing leg: "config.set 第三张名册 role_chat_thinking(248/313钉子原样…)". Its leg branch
  came in the follow-up **`7817ad09a` (2026-11-07)**: "ApplySetting 补末行
  case FieldRoleChatThinking→SetRoleChatThinking" — that commit's own message records that
  before it the field "走 default 被如实拒", i.e. it was **fail-closed until the branch landed**.
- The other six arrived 2026-11-10 (`587e0e22e`, ball_size) and 2026-11-11 (`7f3177f93`,
  the appearance wave's four; `ball_style` in `30d8f2e38`/65 family). So `role_chat_thinking`
  is **three days older and from a different ticket** than the ball/panel six — the ticket's
  suspicion is correct, and it is the *only* member of the difference that is not ball/panel.
- It is not a stray: `fields187` is a deliberate one-name roster
  (`config_write_187.go:45-47`), it carries its own measured tier (`fieldTier187`,
  `:59-61` → `EffectiveNextTask`, corrected at `config_handlers.go:451-453`), and
  `config_write_187_test.go:119-135` pins that the door asks it and that the refusal text
  carries it.

**Conclusion**: `role_chat_thinking` is the same class as the other six — table complete,
door honest, **assertion stale** — but it is a *different family and a different ticket*,
so it must be judged on its own (done above) rather than lumped in as a ball/panel extra.

## 7. Judgement, per field

The ticket offers 甲 (table incomplete ⇒ add to the table) and 乙 (door over-reports ⇒ fix the
door). **Both framings assume the disagreement is inside `internal/panel`. It is not.** The
third possibility is the true one: **the door is right, both tables are right, and the
assertion's `want` is the stale side.**

| field | door has a real branch? | judged side | why |
|---|---|---|---|
| `ball_click_through` | yes `panel_config_store.go:354` | **assertion stale** (neither 甲 nor 乙) | real `case` → `SetBallClickThrough`; in the wave roster by design; adding it to `fields313` duplicates it (M2) and reds the wave disjointness pin |
| `ball_opacity_idle` | yes `:340` | **assertion stale** | same as above |
| `ball_size` | yes `:326` | **assertion stale** | 422 grew the 65 roster on purpose (`config_write_65.go:44-51` names the landing); its own pin asserts `acceptedWritableFields()` carries it (`ball_size_422_test.go:53-54`) |
| `ball_style` | yes `:321` | **assertion stale** | 65 roster, tier measured at `:398-407`; `ball_style_65_test.go:54-55` asserts the report carries it |
| `panel_font_size` | yes `:377` | **assertion stale** | wave roster; `config_write_wave.go` header records the per-create reader |
| `panel_width` | yes `:365` | **assertion stale** | same as above |
| `role_chat_thinking` | yes `:262` | **assertion stale** (separate family) | 187 roster, own tier correction, `config_write_187_test.go:133-134` asserts the report carries it |

Not one field is 乙. The door never lists a name it cannot write: **27/27 have a branch**, and
M4 shows all 27 are accepted once end to end.

## 8. Why I did not change anything, and what the fix is

**Every `internal/panel`-only repair is measurably wrong** (§5): M2 duplicates the names, M6
duplicates them by a different route, M3 lies in the other direction. The only correct edit is
the assertion, which lives in `cmd/wisp/panel_config_313_test.go` — **another leg owns that
package exclusively this wave**, and the dispatch says: if the assertion side must change, stop
and report.

**Recommended fix for the later wave** (assertion side only, ~3 lines at
`cmd/wisp/panel_config_313_test.go:249`):

```go
want := append(append([]string{}, panel.WritableFields()...), panel.AdditionalWritableFields()...)
want = append(want, panel.WritableFields187()...)
want = append(want, panel.WritableFields65()...)
want = append(want, panel.WritableFieldsWave()...)
sort.Strings(want)
```

Measured: green (M1/M4), and still has teeth — truncating the report by one roster reds it and
names exactly the missing fields (M5). Half 2 then also covers all 27, which is a *stronger*
pin than today (M4 passes, so all 27 really do reach the leg once).

⚠ **Do not** take the shortcut of adding the seven names to `fields313`: M2 is that exact
experiment and it is red in two independent ways.

⚠ **Do not** delete the three later roster joins to satisfy the assertion: M3 is that
experiment and four `internal/panel` pins red.

The assertion's comment (`:219-224`) should be updated in the same edit: "the union the door
REPORTS is exactly `WritableFields()` + `AdditionalWritableFields()`" is the stale sentence —
the door reports a **five**-roster union, and `acceptedWritableFields` is the single accessor
that computes it (A881 条件ⓑ). Reading `acceptedWritableFields()` directly would be the
most faithful form, except it is unexported in `internal/panel`; hence the explicit
five-roster append above.

## 9. Gates (AC#3), read on this HEAD with **no** repo change

| gate | command | reading |
|---|---|---|
| build | `go build ./cmd/wisp/` | exit **0**, no output |
| vet | `go vet ./cmd/wisp/` | exit **0**, no output |
| format | `D:\work\base\gopath\bin\gofumpt.exe -l internal/panel/ cmd/wisp/` | **empty**, exit 0 |
| D22 | `D:\work\soft\Git\bin\bash.exe -c "scripts/d22scan.sh"` | `d22scan: clean - no D22 ban violations`; 480 production Go files examined; `PASS=45 FAIL=0 SKIP=0` |

Note on gate 3: `gofumpt` is **not on PATH** on this host; it is at
`D:\work\base\gopath\bin\gofumpt.exe`. The first attempt to run it bare failed with
"not recognized" — that is a setup miss, not a clean reading, so it was re-run by full path.

`build/mut_rvw.go` untouched and still present (checked).

## 10. Boundary (AC#2)

`git diff --name-only` for this leg: **no production file at all** — no `internal/panel/**`
change (correctly: none is warranted), and no `cmd/wisp/**` change. The only files this leg
adds are this evidence file and one append-only Progress-log line in the ticket. Nothing under
`docs/PLAN.md`, `docs/specs/**`, `internal/observe/thresholds.go`, golden files,
`tools/d22scan/allowlist.txt`, `frontend/**` or `design/**` was touched. No AC box flipped.
No push.
