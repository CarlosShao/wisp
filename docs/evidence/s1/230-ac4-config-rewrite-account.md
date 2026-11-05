# C26 rewrite account x `risk.blacklist_overrides` — internal/config de-red (2026-11-05)

- Leg: `internal/config` de-red + wording cleanup. Dispatch anchor `0fa1527f9`; the worktree HEAD when this leg started was `0af3b17aa`, and other legs landed `23cfd5725`, `ad4aa9b90`, `2addacf7d` while this one ran (the last is the composition-root arming, see §5), so this leg's commits interleave with theirs (`git show --stat` on each of this leg's commits lists only this leg's files). Branch `dev`. Owner dispatch: ticket 230 AC#4 (`A1549`) meets the ticket 102 AC#5 static criterion.
- Write face actually touched: `internal/config/overrides.go`, `internal/config/overrides_230_test.go`, `internal/config/overrides_c26_rewrite_test.go` (new), `internal/config/unwired.go`, this file.
- Not touched: `internal/risk/**` (frozen), `internal/tools/**`, `cmd/wisp/**`, `docs/reports/**`, the ticket pool.

## 1. RED before (this leg's own run, worktree HEAD `0af3b17aa`, dispatch anchor `0fa1527f9`)

```
go test ./internal/risk -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg -count=1
--- FAIL: TestC26RewriteAccountIsConsumedAtEverySecurityLeg (0.15s)
    pathresolver_rewrite_account_test.go:200: ..\..\internal\config\overrides.go:78 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
FAIL	github.com/CarlosShao/wisp/internal/risk	0.215s
```

That is the C26 shape exactly: `overrides.go` did `path := res.Canonical` straight after `risk.Resolve`, and that spelling was then **published** as the exemption key the gate matches.

## 2. What changed

### 2a. `overrides.go` — consume the account, never key on a rewritten spelling

- The tier check still runs on the **expanded** form, on purpose (SPEC-06 §4: expand, then classify). That is what makes `~/.ssh/id_rsa` reach `ClassA` instead of matching the literal B-tier `id_*` name rule; the A-tier branches keep their tier reason and now also append the expansion (`C26 also expanded this spelling: home; ...`) when one happened.
- The **B tier** is the only branch whose answer becomes a published key, so that is where the account binds: `path, err := res.Actable()`. A rewritten spelling is a load error — `config.toml: risk.blacklist_overrides[i] %q cannot be a standing exemption: <ErrRewrittenPath: spelling expands to canonical (kind)> ... Write the entry as the literal path, or remove it`. `c.Risk.BlacklistOverrides` never receives the expanded tree.
- A rewritten entry that also matches no B rule reports both facts (no-B-rule + the expansion). Handle-resolution substitutions (8.3 short names, junction targets, mixed separators, redundant `.` segments) are **not** rewrites, so normal canonicalization is unchanged.
- Header duties rewritten from three to four (CONSUME THE REWRITE ACCOUNT first), and the stale "STILL UNDONE: tools.Options.Confirmations ... Gate's bOverrides ... comparison fold unexported" paragraph replaced: `risk.OverrideKey` and `risk.GateStanding` now exist, `tools.Options.PersistedExemptions` is the provider, and the composition-root arming landed in `2addacf7d` during this leg (§5).

### 2b. `unwired.go:125` — the disposition named symbols that do not exist

Before:

> the list is the persisted exemption set tools.Options.Confirmations will key for risk.Gate's bOverrides map

After:

> the list is the persisted exemption set tools.Options.PersistedExemptions returns to risk.GateStanding as its third argument (the second is the runtime L2 confirmations), armed by cmd/wisp (armPersistedExemptions, keyed with risk.OverrideKey) and re-armed after every reload

Verified symbols: `tools.Options.Confirmations` is the **runtime L2** confirmation flow (`internal/tools/bridge.go:103-108`), a different claim about the same file; `tools.Options.PersistedExemptions` (`bridge.go:109-122`) is the config list's provider, read at `internal/tools/mode.go:99,106` and passed to `risk.GateStanding` (`internal/risk/blacklist.go:105`, `func GateStanding(canonical string, l2, persisted map[string]bool)`). The disposition names those and the arming that landed in `cmd/wisp` (`armPersistedExemptions` + the `reloadOnce` re-arm, §5).

### 2c. Tests

- `overrides_230_test.go` positive half: the second accepted entry is now a **literal** second `.env` spelled with a redundant `.` segment (`later/./.env`), not `~/later/.env`. The reading is unchanged — the Config carries the C26-canonical spelling, not the raw string — but the entry no longer depends on expansion, because expansion is now the refusal the new test pins. (The A-tier test is untouched: `~/.ssh/id_rsa` still reports the tier verdict.)
- `overrides_c26_rewrite_test.go` (new) `TestBlacklistOverridesRefuseARewrittenEntry`: `%VAR%\project\.env` (kind `env`) and `~/later/.env` (kind `home`) both fail the load with an error that (a) names `risk.blacklist_overrides`, (b) is the rewrite refusal, not a tier verdict, (c) reports `expands to` + the kind, (d) tells the operator to spell the literal path, and (e) does not fall back to the old "does nothing" sentence.

## 3. Gates (this leg's own runs, `unset GOTMPDIR TMPDIR TEMP TMP` before each)

| gate | command | reading |
|---|---|---|
| C26 static criterion | `go test ./internal/risk -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg -count=1` | `ok github.com/CarlosShao/wisp/internal/risk 0.203s`, `rc=0` (before: the FAIL above) |
| full risk package | `go test ./internal/risk -count=1 -timeout 1200s -v` | `ok github.com/CarlosShao/wisp/internal/risk 16.575s`, `rc=0`; `--- PASS` 288, `--- FAIL` **0**, `--- SKIP` **1**: `TestSyncRegistryProbeLive` (`syncdirs_windows_test.go:106: no registry-grade sync record on this machine (HKCU Accounts without UserFolder is the documented reality here)`) — environment, not this change |
| full config package | `go test ./internal/config -count=1 -timeout 600s` | `ok github.com/CarlosShao/wisp/internal/config 2.837s`, `rc=0` |
| focused config tests | `go test ./internal/config -run "TestBlacklistOverrides\|TestATierStaysDenied\|TestUnwired\|TestEveryLockedSectionKeyIsAccountedFor" -count=1 -v` | PASS, including the new rewrite-refusal test (2 subtests) |
| cross-check, landed arming | `go test ./cmd/wisp -run Test230AC4 -count=1 -timeout 420s` | `ok github.com/CarlosShao/wisp/cmd/wisp 5.161s`, `rc=0` (read-only; `cmd/wisp/**` not edited) |
| build | `go build $(go list ./... \| grep -v "/build$")` (41 packages) + `go build ./cmd/wisp` | both `rc=0`; literal `go build ./...` `rc=1` with the sole diagnostic `# github.com/CarlosShao/wisp/build` (git-ignored stray, see §4) |
| format | `/d/work/base/gopath/bin/gofumpt.exe -l internal/config` | no output, `rc=0` |
| D22 bans | `sh scripts/d22scan.sh` | `rc=0`, clean, see §4 |

## 4. Closing readings

- `go build ./internal/... ./tools/...` — `rc=0`; `go build ./cmd/wisp` — `rc=0`; `go build $(go list ./... | grep -v "/build$")` (41 packages) — `rc=0`.
- `go build -overlay <HEAD-content-overlay> ./cmd/...` — `rc=0` (first measurement, while the concurrent wiring leg's files were still half-written). The overlay restores `cmd/wisp/run.go` and `cmd/wisp/config_reload.go` to `HEAD:0fa1527f9` and replaces the two untracked files with an empty `package main`, so this reads "the app as HEAD commits it".
- Literal `go build ./...` on the shared worktree — `rc=1`, **not this leg's**, with exactly one failing package: `github.com/CarlosShao/wisp/build`. The stray `build/mut_rvw.go` is `package main` but references `cmd/wisp` symbols (`voiceTtsFiles`, `residentBall`, ...), so `go build ./...` treats the stray `build/` directory as a package. The whole directory is git-ignored (`.gitignore:14: build/`), so it is not in HEAD and does not even show as `??`; it is a mutation-test artifact and repo rule 8 forbids deleting temp files, so it is reported, not touched. Excluding only that directory, all 41 module packages build (`rc=0`).
- During the FIRST measurement of this leg, `github.com/CarlosShao/wisp/cmd/wisp` also failed: untracked, half-written `cmd/wisp/blacklist_overrides_230.go:40: "sync/atomic" imported and not used`, while another leg was editing that face (`M cmd/wisp/run.go`, `M cmd/wisp/config_reload.go`). That leg landed its work as `2addacf7d` during this leg (§5); the closing readings are `go build ./cmd/wisp rc=0` and `go test ./cmd/wisp -run Test230AC4 rc=0` (5.161s). `cmd/wisp/**` stayed outside this leg's write face throughout.
- `sh scripts/d22scan.sh` — `rc=0` on the committed tree, final line: `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=363, bans #1-5 cmd/=58, ban #6 frontend/=116, ban #7 internal/tools/=48, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=904, ban #8 cmd/=211` (904 because this leg added one `internal/config` test file; the shared tree's cmd/ count moved for the other leg's untracked file).

## 5. Wiring state (another leg landed it mid-leg)

`2addacf7d` — not this leg, write face `cmd/wisp/**` — closed the composition-root half while this leg ran: `cmd/wisp/run.go` arms `tools.Options.PersistedExemptions: rt.persistedExemptions` before the bridge is built (`run.go:984`, `run.go:1016`), `cmd/wisp/blacklist_overrides_230.go` folds each entry with `risk.OverrideKey` (`exemptionKeySet`), and `cmd/wisp/config_reload.go:192` re-arms from `reloadOnce`. So the list this leg publishes now has a real consumer: an accepted entry can come back from `risk.GateStanding` as `Allow` with the `source=config:risk.blacklist_overrides` log line.

This leg's share stays what it was: only literal, C26-canonical, B-tier spellings can enter the list, so the armed map cannot contain a tree the config never named.

Not verified by this leg (outside its face): the wiring leg's own end-to-end/reload evidence; its commit message records `go build ./cmd/... ./internal/... ./tools/... rc=0`, `go test ./cmd/wisp -run Test230AC4 rc=0 PASS 2`, `gofumpt -l cmd/wisp` empty, `d22 rc=0`. This leg only re-ran `Test230AC4` as a cross-check (5.161s, rc=0) and did not redo their mutation readings.

Concurrency note: `cmd/wisp/blacklist_overrides_230.go` was untracked and half-written when this leg measured its first build; it is now committed by `2addacf7d`. This leg never edited `cmd/wisp/**`; the only files it committed are listed at the top.
