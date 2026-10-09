# 247 v1 — the four hard gates, self-run, plus the per-name red-roster diff

All commands run by `247-v1` on this machine, 2026-10-09 13:2x–13:5x +08, tree at `b04ec963`
plus this leg's evidence-only commits. Ticket text for this cell (`247…md:31`) verbatim:
"⛔ 框归编排者，产码腿与验收腿一枚都不许碰 … 逐名照抄终态".

## 1. Gate lines (one `rc=N` each)

| # | command (verbatim as run) | rc | terminal reading (this leg's own) |
|---|---|---|---|
| G1 | `GOFLAGS= go build ./...` | `rc_build=0` | zero output (13:33:49) |
| G2a | `"$(go env GOPATH)/bin/gofumpt" -l <the batch's 10 files, named individually>` | `rc_gofumpt_roster10=0` | **empty output** (files listed in §3) |
| G2b | `"$(go env GOPATH)/bin/gofumpt" -l internal/audio cmd/wisp` | `rc_gofumpt_dirs=0` | 3 names, **none of them in this batch**: `cmd\wisp\models.go`, `cmd\wisp\panel_inbound_guards_35r3_test.go`, `cmd\wisp\panel_transport_35r2_test.go` |
| G2c | `"$(go env GOPATH)/bin/gofumpt" -l .` (whole repo, denominator only) | 0 | 43 lines, dominated by other legs' scratch `.go` files under `.scratch/wisp/probes/**` (e.g. `probes\161\r1\runs\b8-probe-bands\internal\a\p_b8.go`, `probes\163\a1\main.go`, `probes\212\v1\mut\main-noq9.go`) — this leg judged **nothing** from the whole-repo form, per the dispatch's warning that its denominator contains other people's files |
| G3 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/audio/ ./internal/ball/ -count=1` | `rc=1` | `--- FAIL` ×6 (5 in `cmd/wisp`, 1 in `internal/ball`); `FAIL github.com/CarlosShao/wisp/cmd/wisp 448.159s` / `ok github.com/CarlosShao/wisp/internal/audio 16.184s` / `FAIL github.com/CarlosShao/wisp/internal/ball 0.187s` |
| G3v | the same command with `-v` | `rc=1` | `PASS=374 FAIL=6 SKIP=3` (this leg's own `grep -c` over `--- PASS` / `--- FAIL` / `--- SKIP`; 501 `=== RUN` lines) |
| G4a | `./tools/d22scan/d22scan.exe` | `rc_d22scan_exe=0` | `d22scan: clean - no D22 ban violations; … bans #1-5 internal/=229, bans #1-5 cmd/=39 …` and `examined 268 production Go files under internal/ and cmd/` |
| G4b | `sh scripts/d22scan.sh` (the CI wrapper: positive control first, then `go run . -root`) | `rc_d22scan_sh=0` | positive control ran and passed first (`--- PASS: TestScanDetectsAllSeededViolations (0.04s)`, `--- PASS: TestScanCleanRepoIsGreen (0.03s)`, `--- PASS: TestAllowlistSuppressesOnlyListedPaths …`), then the same clean verdict |
| G5 | `PATH=… WISP_LIVE_MIC=1 go test ./cmd/wisp/ -run TestAC247LiveMicrophoneLevelsReachTheBallSeam -count=1 -v -timeout 180s` | `rc_live=1` | real device, real thread, **red by design** — readings in `30-ac2-level-ruler.md` |

Two rulers about the instruments themselves, because "the gate said clean" is only worth what the
gate can see:

- `go run ./tools/d22scan` is indeed not the way to run it (ticket's ⛔ note): `tools/d22scan` is its
  own module (`tools/d22scan/go.mod`), so G4a/G4b are the only two forms this leg used. The shipped
  `d22scan.exe` is dated `2026-10-04 09:10:01` while `tools/d22scan/main.go` is dated
  `2026-10-04 10:54:10`, i.e. **the exe predates its own source by ~1h45m**; that is why this leg ran
  G4b too (it recompiles from source through `go run .`) and got the same verdict, so the clean
  reading does not rest on a stale binary.
- Without the DLL harness the same package is an *environment* red, not a code red:
  `go test ./cmd/wisp/ -count=1 -run TestNoSuchTest247v1` => `exit status 0xc0000135` with no
  `--- FAIL` line at all (the test binary never starts). Recorded so nobody reads a 0xc0000135 as
  something this batch caused.

## 2. Per-name red roster, and the diff against the baseline (criterion = names, not counts)

Baseline as given in the dispatch (orchestrator, same tree, `-v`, 12:0x–12:2x): 5 names, all the
real-window family. This leg's own `-v` run:

```
--- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.05s)
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (5.11s)
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.01s)
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.02s)
--- FAIL: TestAC14GoSideEvalPushReachesThePage (20.01s)
--- FAIL: TestC21TableColourRowsMatchTokensCSS (0.00s)      [internal/ball]
```

- `cmd/wisp`: same 5 names, in the same shape (three at 20.01s = the window-wait timeout form, two at
  ~5s). **Set difference against the baseline = 0 new names.**
- The count itself is not the criterion and this leg did not use it: the dispatch says the family
  drifts between 4 and 5 on this tree; the non-`-v` run and the `-v` run here both gave 5.
- This leg added **no `-skip`** and relaxed no assertion (it changed no test file at all — write
  surface = `probes/247/v1/*.md`). The 3 skips in this leg's run are:
  `TestPanelHostLatencyPercentilesAC2` (pre-existing), `TestAC247LiveMicrophoneLevelsReachTheBallSeam`
  (env-gated at `resident_audio_247_live_windows_test.go:129-131`, same precedent as)
  `TestLiveWasapiSmoke` (pre-existing env gate).

### The sixth red: worktree reading, **not** a defect of the batch

This leg judged the implementer's claim independently rather than accepting it:

```
git cat-file -e HEAD:design/assets/tokens.css          => rc=0 (the file IS in HEAD)
git status --porcelain -- design/assets/tokens.css     =>  D design/assets/tokens.css
                                                        (plus: "warning: could not open directory 'design/assets/'")
git status --porcelain -- design | grep -c "^ D"       => 16   (somebody else's uncommitted deletions)
ls design/assets/tokens.css                            => No such file or directory
```

And the case really does read from disk, which is what makes a HEAD-only check insufficient. Note the
two different line numbers, because the red line is **not** the read line: the report prints
`tokens_table_test.go:1468` (` TestC21TableColourRowsMatchTokensCSS`'s call
`	dark, light := c21ParseTokensCSS(t, root)`) only because `c21ParseTokensCSS` declares `t.Helper()`
(`:1382`); the actual disk read is `internal/ball/tokens_table_test.go:1384`
`	data, err := os.ReadFile(path)` with `path` built from `:1364`
`const c21TokensCSSPath = "design/assets/tokens.css"`, and the fatal text is assembled at `:1386`
`		t.Fatalf("read %s: %v - the CSS leg of this check must never skip", c21TokensCSSPath, err)`.
This leg's own run printed it verbatim:

```
    tokens_table_test.go:1468: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified. - the CSS leg of this check must never skip
```

Attribution ruler: the batch's file roster (§3) contains no `internal/ball/**` and no `design/**`
entry, and `git diff --name-only --no-renames 39d17f0b..b04ec963 -- internal/ball design` => empty.
⇒ **judgement: this is a reading of the working tree, not a defect of the thing under review**;
consistent with the same root cause the ledger already books for `internal/panel`'s four-way token
check. CI (checking out HEAD) cannot reproduce it. This leg restored nothing and committed nothing
under `design/**`.

## 3. The batch's own file roster (what G2a was pointed at)

`git diff --name-only --no-renames 39d17f0b..b04ec963` => 13 names:

```
.scratch/wisp/issues/247-…-ball.md            (Progress-log line only, appended by 247-r1)
.scratch/wisp/probes/247/r1/00-baseline.md
.scratch/wisp/probes/247/r1/10-gates.md
.scratch/wisp/probes/247/r1/20-ac-readings.md
cmd/wisp/config_readers_255.go
cmd/wisp/resident_audio_247_live_windows_test.go
cmd/wisp/resident_audio_247_windows_test.go
cmd/wisp/resident_audio_windows.go
cmd/wisp/resident_windows.go
internal/audio/capturelevel_windows_test.go
internal/audio/captureopt.go
internal/audio/captureopt_test.go
internal/audio/wasapimic_windows.go
internal/audio/wavinjector.go
```

(10 Go files; the 3 remaining `cmd/wisp` files G2b lists — `models.go`, `panel_inbound_guards_35r3_test.go`,
`panel_transport_35r2_test.go` — belong to other legs: `5e8748b3` (212-r1＋258-r1), `7f9d6e40` (35-r3),
`3a343bc7` (35-r5), by `git log -1 --format="%h %s" -- <file>`.)

## 4. Worktree hygiene at hand-in (this leg's own)

`git status --porcelain -- internal cmd docs .scratch/wisp/issues` => **empty** (this leg touched no
tracked file); the only writes are the six `.md` under `.scratch/wisp/probes/247/v1/`. `.gitignore`
still ` M`, `design/**` still 16 ` D` + 15 ` M`, `build/wisp.exe` untouched — all somebody else's,
nothing restored, nothing committed by this leg; zero push.
