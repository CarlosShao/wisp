# 247 r1 — AC#8 gates, rosters, edge rulers

Leg `247-r1`. Own-file roster = `git diff --name-only 39d17f0b HEAD` (the commit before this
leg's first): 10 files, listed at the bottom. All readings below are this leg's own runs on
this machine, 2026-10-09 12:2x–13:0x +08.

## 1. The five gate lines (ticket AC#8, verbatim commands)

| # | command | rc | terminal reading |
|---|---|---|---|
| 1 | `GOFLAGS= go build ./...` | `rc=0` | zero output (success) — run twice: after the first wiring pass and again after the final commit `ada5c563` |
| 2 | `"$(go env GOPATH)/bin/gofumpt" -l internal/audio cmd/wisp` (bare `gofumpt` is not on this machine's PATH; `rc=127` on the first try, then the same command through `$(go env GOPATH)/bin`, where `gofumpt.exe` sits next to `staticcheck.exe`) | `rc=0` | lists 3 files, NONE of them mine: `cmd/wisp/models.go`, `cmd/wisp/panel_inbound_guards_35r3_test.go` (ticket 288's, untouched by order), `cmd/wisp/panel_transport_35r2_test.go`. My 10 files listed individually answer empty: `gofumpt -l internal/audio/captureopt.go internal/audio/captureopt_test.go internal/audio/capturelevel_windows_test.go internal/audio/wasapimic_windows.go internal/audio/wavinjector.go cmd/wisp/resident_audio_windows.go cmd/wisp/resident_audio_247_windows_test.go cmd/wisp/resident_audio_247_live_windows_test.go cmd/wisp/resident_windows.go cmd/wisp/config_readers_255.go` => no output. No `-w` over the repo: `-w` was pointed at exactly three files this leg authored (`captureopt_test.go`, `resident_audio_windows.go`, `config_readers_255.go`) |
| 3 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/audio/ ./internal/ball/ -count=1` | `rc=1` | run twice on this leg's tree. 12:4x (before the roster re-adjudication): `FAIL cmd/wisp 449.464s` / `ok internal/audio 16.119s` / `FAIL internal/ball 0.186s`, and its 6th red was `TestTicket255RosterStillMatchesTheActualReadSites` — **caused by this leg**, fixed in `ada5c563`. Final tree: the same command with `-v` appended, `rc=1`, `PASS=374 FAIL=6 SKIP=3`, name roster in §2 (the red names are identical in both forms; only the 255 one disappeared) |
| 3v | same command with `-v`, final tree at `fb846b69` | `rc=1` | `PASS=374 FAIL=6 SKIP=3` — name roster in §2; `FAIL github.com/CarlosShao/wisp/cmd/wisp 447.762s` / `ok github.com/CarlosShao/wisp/internal/audio 16.518s` / `FAIL github.com/CarlosShao/wisp/internal/ball 0.178s` |
| 4 | `./tools/d22scan/d22scan.exe` | `rc=0` | `clean - no D22 ban violations`, live scope work `bans #1-5 internal/=229, cmd/=39`, `ban #8 internal/=524 Go files, cmd/=116`, `examined 268 production Go files` — i.e. this leg's production files are inside the examined count and no bare `go func(` was added. Run once mid-leg and once on the final tree, same answer |

Gate 1 and 4 were also re-run after the last commit (`final-build rc=0`, `final-d22 rc=0` in
the same command as the `-v` run above).

The `-v` run's raw log is 330 KB and is deliberately **not** in the repo (MB-scale rule);
only the extracted lines below are.

## 2. Red roster after the change, by name, and the diff against the orchestrator's baseline

Baseline given in the dispatch (orchestrator, same tree, `-v` form, 2026-10-09 12:0x–12:2x):
`TestPanelHostRealWindowHopAndLifecycle`, `TestAC4FocusReturnToPriorWindowGap33r5`,
`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`, `TestAC14AwaitedBindingReplyReachesThePage`,
`TestAC14GoSideEvalPushReachesThePage` (5 names, all the real-window family).

After this leg (final `-v` run):

```
--- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.05s)
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (5.09s)
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.02s)
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.02s)
--- FAIL: TestAC14GoSideEvalPushReachesThePage (20.02s)
--- FAIL: TestC21TableColourRowsMatchTokensCSS (0.00s)          [internal/ball]
```

Set difference against the baseline = **0 new names in cmd/wisp**.

The sixth red is in `internal/ball` and its cause is a worktree file somebody else deleted:

```
tokens_table_test.go:1468: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css:
The system cannot find the path specified. - the CSS leg of this check must never skip
```

Ruler that it is not this leg's: `git cat-file -e HEAD:design/assets/tokens.css` => the file
**exists in HEAD**, so the absence is only in the working tree (the dispatch's §6/§7 names the
same 16 uncommitted `design/**` deletions and forbids restoring them); and this leg's own file
roster (§5) contains no `internal/ball` and no `design/**` file at all. Same root cause the
dispatch already booked for `internal/panel`'s `TestC21DesignTokensFourWayAgree`.

Skips after the change (3, none added by `-skip`):

```
--- SKIP: TestPanelHostLatencyPercentilesAC2           (pre-existing)
--- SKIP: TestAC247LiveMicrophoneLevelsReachTheBallSeam (this leg's AC#2 harness: needs WISP_LIVE_MIC=1, the precedent TestLiveWasapiSmoke uses)
--- SKIP: TestLiveWasapiSmoke                          (pre-existing, same env gate)
```

This leg added no `-skip` and relaxed no assertion.

## 3. AC#1 ruler (the ticket's own command, both shapes)

```
grep -rl "CarlosShao/wisp/internal/audio" --include=*.go . | grep -v "/internal/audio/" | grep -vc _test
```

- before (12:2x, file `00-baseline.md`): `0`
- after (12:3x, same tree, no other leg in between): `1`
- the one name: `./cmd/wisp/resident_audio_windows.go` — the resident leg itself, not another
  debug cmd (AC#1's parenthetical).

## 4. AC#3 / AC#0 edge rulers (capability, not word-level)

```
GOOS=windows go list -deps ./internal/audio | grep -c "internal/ball\|internal/statemachine"   => 0
GOOS=windows go list -deps ./cmd/wisp        | grep "internal/audio"                             => github.com/CarlosShao/wisp/internal/audio
```

So: the audio package still imports neither the ball nor the state machine, and the single new
package-level edge is `cmd/wisp -> internal/audio`.

The new cross-package signatures this leg added, complete:

| new call | signature | crosses into internal/ball? |
|---|---|---|
| `audio.WithLevelSink` | `func(fn func(float32)) CaptureOption` | no |
| `audio.WithCaptureRegistry` | `func(r *observe.Registry) CaptureOption` | no |
| `audio.SpawnCapture` (pre-existing seam, first production caller) | `func(*observe.Registry, func(ctx context.Context)) *observe.Handle` | no |
| `residentBall.setAudioLevel` -> `Ball.SetAudioLevel` | `func(level float32)` | yes, one scalar |

No `[]byte`, no `[]int16`, no `samples` parameter appears in any call that crosses into
`internal/ball`. The C8 seam's own `AudioSource.Start(ctx, chan<- []byte)` is ticket 13's
pre-existing exported contract, called by cmd/wisp against internal/audio; frames stop at this
process's bounded channel and never reach the ball. Reading AC#3's second sentence as covering
that pre-existing seam would forbid the落点 the orchestrator ruled in AC#0, so it is judged
against the ball boundary, as its first sentence and the 更正 section (`liquid_windows.go:42`)
state. Flagged as an interpretation in the leg report.

## 5. This leg's own file roster

```
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

`git diff --name-only 39d17f0b HEAD | grep -c "<AC#7 forbidden paths>"` => `0`.

## 6. Commits (this leg, zero push)

```
2304ca12 pre-change baseline evidence
5d407e77 production wiring (cmd/wisp -> internal/audio, level sink, gate, step-4 registrant)
b1bdc61f capture-seam tests (internal/audio)
cdaf5953 assembly tests (AC#4/AC#5/AC#6) + the WISP_LIVE_MIC AC#2 harness
ada5c563 ticket-255 config-reader roster re-adjudication (forced by its own instrument)
fb846b69 AC#10 reading test (PrototypeVisualsEnabled()=false at the shipped default)
(+ this evidence file and the ticket Progress-log line, one further commit)
```

Zero push. Nothing was deleted from the repo; the two raw long logs this leg produced
(`/tmp/247-gate3-v.log`, `/tmp/247-final-v.log`, ~330 KB each) are left outside the repo and not
committed, per the MB-scale rule, and only their extracted lines are in here.
