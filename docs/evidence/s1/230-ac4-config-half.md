# 230 AC#4 — B-tier override: internal/config half (landed) and the wiring seam (next wave)

- Ticket: `.scratch/wisp/issues/230-four-cells-left-unfinished-inside-closed-tickets.md`, AC#4.
- Owner ruling: **B** — open the B-tier override as a user-customizable path (2026-11-05, ledger `A1549`, ticket tail line 63).
- Leg: `230-ac4-config` (implementation leg, `internal/config` face only). Start anchor `02823947a` (the ruling commit), branch `dev`.
- Write face actually touched: `internal/config/overrides.go` (new), `internal/config/overrides_230_test.go` (new), `internal/config/validate.go`, `internal/config/unwired.go`, `internal/config/schema.go`, `internal/config/unwired_test.go`, this file.
- Not touched: `internal/risk/**` (frozen), `internal/tools/**`, `cmd/wisp/**`, `PLAN.md`, `docs/specs/**`, `frontend/**`, `design/**`, the ticket file itself (frame and boxes belong to the orchestrator and a non-implementer).

## 1. What actually changed, in one paragraph

`risk.blacklist_overrides` stopped being a registered lying key and became a validated setting. The guard row and its `lockedKeyDisposition: unwired:` entry are gone; in their place `validateBlacklistOverrides` runs in the load pipeline: every entry goes through the C26 resolver (`risk.Resolve`, with the same `[fs].reparse_point_exceptions` the tool path uses), the Config carries the resulting canonical spelling, an entry that classifies `ClassA` is a load error (A tier is absolutely denied and never overridable), and an entry that matches no B-tier rule is a load error too (an exemption that could never fire is the same lie in a new coat). The A/B verdict itself was not touched: it still comes from the frozen `internal/risk` tables, and `Validate` calls the classifier rather than restating the rules.

## 2. RED before GREEN (both readings are this leg's own runs)

RED — test file added, production code untouched (`02823947a`):

```
go test ./internal/config -run "TestBlacklistOverrides|TestATierStaysDenied" -count=1 -timeout 240s -v
rc=1  RUN=13  top-level FAIL=2
  --- FAIL: TestBlacklistOverridesParseAsAnExemptionList
  --- FAIL: TestBlacklistOverridesRefuseATier   (9 subtests, one per A-tier example)
error text observed: "config.toml: risk.blacklist_overrides is written but does nothing: ..."
```

That is the ticket's premise, measured: a user who wrote the key got an error whose reason was "nothing reads this value", not a tier verdict.

GREEN — after the change:

```
go test ./internal/config -run "TestBlacklistOverrides|TestATierStaysDenied|TestUnwired" -count=1 -timeout 240s -v
rc=0  RUN=...  --- PASS: 9 top-level
  TestBlacklistOverridesParseAsAnExemptionList      PASS
  TestBlacklistOverridesRefuseATier                 PASS (9 subtests)
  TestBlacklistOverridesRefuseAnInertEntry          PASS
  TestATierStaysDeniedEvenWhenItReachesTheOverrideSet PASS
  TestUnwiredSecurityKeysFailLoudly                 PASS (5 rows now, was 6)
  TestUnwiredGuardLeavesHonestConfigsAlone          PASS
  TestUnwiredGuardFiresBeforeAnyFileWrite           PASS
  TestUnwiredGuardIsScopedNotABigStick              PASS
  TestUnwiredKeysStillRoundTripAtTheByteLevel       PASS
```

What each new test pins, and the failure it would catch:

| test | reading it pins | the shape that goes red |
|---|---|---|
| `TestBlacklistOverridesParseAsAnExemptionList` | two valid B entries (a real `.env`, a `~/later/.env`) load, and the Config carries exactly the `risk.Resolve` canonical spelling of each; every surviving entry classifies `ClassB` | a loader that keeps the operator's raw spelling (the exemption may then silently never match), or one that drops/duplicates entries |
| `TestBlacklistOverridesRefuseATier` | nine A-tier examples are refused by name: `~/.ssh/id_rsa` (both spellings), `.git-credentials`, `.aws/credentials`, `.kube/config`, a Chromium `Login Data` store, DPAPI master keys, the Windows credential vault, `%APPDATA%\wisp\config.toml` | a tier check that runs on the literal string before expansion: `~/.ssh/id_rsa` would match the B-tier `id_*` name rule and slip through |
| `TestBlacklistOverridesRefuseAnInertEntry` | an entry that matches no B rule (`notes.txt`) is refused, and the refusal is the tier verdict, not the old "does nothing" sentence | a validator that accepts any absolute path, leaving a standing "exemption" that exempts nothing |
| `TestATierStaysDeniedEvenWhenItReachesTheOverrideSet` | even with the A path keyed into `risk.Gate`'s override map under both its canonical and lowercased spelling, the verdict is `ClassA`, no `Allow`, no `NeedL2`, and an A-tier security line is emitted | a wiring wave that ever lets the persisted set sit in front of the tier check |

The A-tier examples are spelled through fake anchors (`HOME`/`USERPROFILE`/`APPDATA`/`LOCALAPPDATA` all under one temp root), so the test measures the real tables without touching the machine's profile.

## 3. What this half does and does not give the AC

Satisfied now:

- the key is a **real validated setting** (no loud-failure registration left; `TestEveryLockedSectionKeyIsAccountedFor` still passes on the new disposition);
- a user can **persist** the list (schema unchanged, `SaveFile` round-trip still green) and it survives load as a canonical exemption list;
- **A-tier is not overridable** is enforced at load for every example the ticket names, and re-proved at the gate for a set built by other means;
- the "lying key" defect that AC#4 inherits from ticket 80 is gone: the error a user sees now is the tier verdict, not "nothing reads this".

Not satisfied by this half (needs the wiring wave below, by construction):

- no consumer yet: `tools.Options.Confirmations` is still nil in production, so `risk.Gate` still receives an empty override map and every B-tier path still answers "ask first";
- therefore **no overridden use happens yet**, which is also why the "every overridden use emits a security log line" half of the AC cannot be accepted in this wave (the line exists in `internal/risk/blacklist.go:87`, and `risk.Logf` is armed in production at `cmd/wisp/run.go:816` and `internal/agent/approval/gate.go:187`; what is missing is the use, and the sentence's honesty for a config-sourced use, see 4a).

## 4. The seam, exactly — what the wiring wave has to change

This leg STOPPED here on purpose: the remaining work is in `internal/risk` (frozen) and in `internal/tools` / `cmd/wisp` (other legs this wave). Three edits, in dependency order.

### 4a. `internal/risk` (frozen: needs its own wave/approval) — two additions, no classification change

1. **Export the comparison fold.** `overrideApplies` (`internal/risk/blacklist.go:116-127`) compares `bOverrides[cand.real]`, and `cand.real` is `normPath` output; `normPath` is unexported (`internal/risk/blacklist.go:134`). `internal/risk/pathshape_portable_test.go:172-174` already says the map is "keyed on normPath output" and names `risk.blacklist_overrides` as a future feeder. The composition root cannot key the map without this fold, and a second copy of it in `internal/tools`/`cmd/wisp` would be a second implementation of a security-comparison key. Minimal shape:

   ```go
   // OverrideKey folds a C26-canonical spelling into the key space bOverrides
   // uses (the same comparison fold Gate applies to its argument). The
   // composition root keys the persisted exemption set with this; there must
   // not be a second implementation of the fold.
   func OverrideKey(canonical string) string { return normPath(canonical) }
   ```

2. **Let the existing override allow line name its source.** `internal/risk/blacklist.go:87` prints `B-list OVERRIDE (single-file, L2 confirmed)`. For a persisted exemption nobody clicked anything, so the sentence would be false exactly where the AC wants evidence. Shape that keeps today's signature (so `internal/tools/mode.go:98` and the frozen tests keep compiling):

   ```go
   func GateStanding(canonical string, l2, persisted map[string]bool) PathDecision
   func Gate(canonical string, bOverrides map[string]bool) PathDecision {
       return GateStanding(canonical, bOverrides, nil)
   }
   ```

   with one line per allowed use, `source=l2-confirm` or `source=config:risk.blacklist_overrides`, and a `PathDecision.Reason` for the config case that does not claim an L2 confirm (that sentence reaches the approval card and the model-facing text; ticket 334 is the precedent that a dishonest reason is a defect, not cosmetics). Constraint to respect: `internal/risk/allow_substring_pin_216_test.go:47-49` pins that the `overrideApplies` window contains `bOverrides[`, and `:115` lists `return bOverrides[cand.real]` as a lookalike that must keep counting zero. Keep that lookup in place, or move the pin deliberately with a stated reason.

### 4b. `internal/tools` — consume the persisted set without merging the two sources

`Options.Confirmations` (`internal/tools/bridge.go:103-108`, fanned out at `:196-199`) is the documented seam for `Gate`'s `bOverrides` map, and `internal/tools/mode.go:91` (`overrides := b.confirmations()`) / `:98` (`risk.Gate(c, overrides)`) is the single production call site. Do **not** union the config list into `Confirmations`: that map is "B-tier single-file overrides on record" (runtime L2), and merging loses the source tag the log line in 4a2 needs. Add a second provider (e.g. `Options.PersistedExemptions func() map[string]bool`), key it with `risk.OverrideKey(entry)`, and call `GateStanding(c, l2, persisted)`.

### 4c. `cmd/wisp` — arm the closure from the loaded config

The closure returns `cfg.Risk.BlacklistOverrides` folded with `risk.OverrideKey`, re-armed on reload. The reload path already puts a locked-section loosening behind the L2 hook (`internal/config/manager.go:334-356`; `riskDirection` names `risk.blacklist_overrides` as a loosen key at `:487`), so an exemption added by hand while the process runs still owes one click, while the startup load is the operator's own standing instruction. Because the list is already C26-canonical (this wave), arming is a pure fold and can cost startup/reload time, not per-call time.

### 4d. Acceptance criteria for the wiring wave (the rest of AC#4)

1. with the persisted set armed, a B-tier canonical path is `Allow: true`;
2. the same path with the set empty is still `NeedL2: true` (fail-closed default unchanged);
3. an A-tier path placed in the persisted set by construction is still `ClassA`, `Allow: false`, `NeedL2: false`;
4. every persisted-exemption use emits exactly one security line naming path, rule and `source=config:...`; the L2 path still says `source=l2-confirm`;
5. a reload that adds or removes an exemption is reflected by the closure (after the L2 hook for the loosening direction).

## 5. Gates (this leg, shared worktree, `unset GOTMPDIR TMPDIR TEMP TMP` before every run)

| gate | command | reading |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | `rc=0` |
| package tests | `go test ./internal/config -count=1 -timeout 240s` | `rc=0`, `ok github.com/CarlosShao/wisp/internal/config 2.837s` |
| format | `/d/work/base/gopath/bin/gofumpt.exe -l internal/config` | output empty, `rc=0` (first pass listed `internal/config/overrides_230_test.go`; `-w` on the two new files, re-check empty) |
| D22 bans | `sh scripts/d22scan.sh` | see the closing reading below |

D22 first run was **not** green, and the finding was this leg's own: `internal/config/overrides.go:43: [phantom-citation] comment cites repo path "docs/evidence/s1/230-ac4-config-half.md" which does not exist on disk (ticket 212 ban #9)`. It failed two tests (`TestScannerSelfScanOfRealRepoIsGreen`, `TestRealRepoLedgerIsHonest`) for one finding. The citation in the code comment now points at this file, which exists, and the re-run reading is recorded as the closing line of this section.

Closing D22 reading: `sh scripts/d22scan.sh` `rc=0`, last line `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=363, bans #1-5 cmd/=57, ban #6 frontend/=116, ban #7 internal/tools/=48, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=901, ban #8 cmd/=209` (420 production Go files examined). One denominator moved between the two green runs - `ban #8 cmd/` 208 then 209 - and it is not this leg's: the working tree gained a `cmd/wisp` test file from another leg running concurrently, while every other denominator and both `internal/` counts stayed identical. Raw run logs are kept outside the repo, next to this leg's other readings, in `D:\tmp\wisp\20261105\230ac4\` (`05-d22scan.log` red, `06-d22scan-after-doc.log` and `12-gate-d22scan-final.log` green).

### 5b. POSIX readings (the CI runs this package on two platforms)

The four new tests and the five guard tests were run in `golang:1.27` (linux/amd64, module cache mounted, `CGO_ENABLED=1`, `GOOS=linux` printed by the container): `PASS`, `ok github.com/CarlosShao/wisp/internal/config 0.107s`. The `%APPDATA%` A-tier case is spelled with `filepath.Separator` on purpose - on POSIX a backslash is an ordinary file-name character, so a literal Windows spelling there names one long file name and is legitimately not the config file.

Full package on Linux: **one** red, `TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed` (`writeguard_unnailed_341_test.go:240: mtime did not move on the write branch`). It is not this leg's: the same package on a pristine `git archive HEAD` tree (this change absent) in the same container reports the identical single red, and the same package on Windows is green (`rc=0`). It is the known mtime-granularity family on the container filesystem, not a verdict about this change.

## 6. Open items for the orchestrator (not decided here)

1. **Ticket frame**: AC#4's box was not flipped, and the ticket's own sentence says only a non-implementer may judge. This half is the `internal/config` half; the AC is not satisfied until 4a-4c land.
2. **The `unwired.go` disposition wording** now says "consumed: validated and canonicalized at load ... the list is the persisted exemption set `tools.Options.Confirmations` will key ... that arming is the remaining wiring leg". If the wiring wave chooses different names (4a2's `GateStanding`, 4b's `PersistedExemptions`), that sentence must move with it - a disposition that names a symbol nobody added is the next stale comment.
3. **Governance reading of the ruling**: with this change a hand-written line in `config.toml` is a standing exemption that takes effect at the next load without a click. That is the owner's explicit call (`A1549`: "完全不给放行太严格，应让用户自定义"), and D36 rule 1 keeps a runtime *change* behind the L2 hook; a future reader who finds "must be clicked by a human" in older tickets should read it as being about changes, not about the file's own contents.
4. **`internal/risk` is frozen**, so 4a cannot be done by a leg whose write face excludes it; the wiring wave needs that face granted, or the ticket needs an explicit risk-side sub-ticket.
