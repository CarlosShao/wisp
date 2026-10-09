# 247 v1 — adversarial acceptance leg (skeleton, committed first)

Leg = `247-v1`, non-implementer adjudicator. Implementer = `247-r1`.
Write surface of this leg = ONLY new `.md` files under `.scratch/wisp/probes/247/v1/`.
No production code, no test file, no AC box, no ledger, no HANDOVER was or will be touched.

- Clock at start (from `date` stdout): `2026-10-09 13:25:48 +0800`
- HEAD at start: `b04ec963` (branch `dev`, ahead 764, zero push by this leg)
- Batch under review (implementer commits): `5d407e77` `b1bdc61f` `cdaf5953` `ada5c563`
  `fb846b69` `618d833b` `303670a5` `b04ec963`, baseline `2304ca12`.
- Ticket verbatim = `.scratch/wisp/issues/247-the-capture-stack-has-zero-importers-so-no-real-microphone-level-ever-reaches-the-ball.md`
  (read in full by this leg, 94 lines, not from the dispatch summary).

## Rulers this leg will run itself (nothing copied from 247-r1's readings)

| # | ruler | where the reading lands |
|---|---|---|
| R1 | `GOFLAGS= go build ./...` | `10-gates.md` |
| R2 | gofumpt: whole-repo denominator + per-file for the 10 roster files | `10-gates.md` |
| R3 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/audio/ ./internal/ball/ -count=1` (and `-v` per-name roster) | `10-gates.md` |
| R4 | `./tools/d22scan/d22scan.exe` (and `sh scripts/d22scan.sh`) | `10-gates.md` |
| R5 | AC#1 importer ruler (ticket 现量 1 command, verbatim) | `20-ac-verdicts.md` |
| R6 | AC#2: level ruler `internal/audio/level.go` reading + a minimal DC-offset probe of my own | `30-ac2-level-ruler.md` |
| R7 | AC#3: every new cross-package call signature with `file:line` | `20-ac-verdicts.md` |
| R8 | AC#4: default values + read sites + the two assembly tests re-run | `20-ac-verdicts.md` |
| R9 | AC#5: bare `go func(` count, owner, `observe` roster, `git diff` on `internal/proc` | `20-ac-verdicts.md` |
| R10 | AC#6: the three device-failure subtests re-run by name | `20-ac-verdicts.md` |
| R11 | AC#7: `git diff --name-only --no-renames` over the batch vs forbidden paths | `20-ac-verdicts.md` |
| R12 | AC#10: `PrototypeVisualsEnabled` + `EnablePrototypeVisuals` caller ruler | `20-ac-verdicts.md` |
| R13 | Q2: ticket-255 roster file diff + `hotClaim*` constants + `-run '255'` final state | `40-cross-ticket-255.md` |
| R14 | Q3: seam legality — every AC#6 red sentence checked for "injected failure claimed as real device" | `50-seams-and-defaults.md` |
| R15 | Q4: `config.NewDefaults()` fallback shape vs ticket 198 / `A754` | `50-seams-and-defaults.md` |
| R16 | Baseline-diff: per-name red roster, new-red count | `10-gates.md` |
| R17 | `TestC21TableColourRowsMatchTokensCSS`: `git cat-file -e HEAD:design/assets/tokens.css` + does the test read from disk | `10-gates.md` |

## Files this leg will author

- `00-skeleton.md` (this file)
- `10-gates.md` (four hard rulers + red-roster diff, each with its own `rc=N` line)
- `20-ac-verdicts.md` (nine cells AC#1..AC#8, AC#10, each one line: 成立 / 不成立 / 部分成立 + evidence)
- `30-ac2-level-ruler.md` (question 1: is the 0.37 floor a defect of the level ruler itself)
- `40-cross-ticket-255.md` (question 2: roster re-adjudication)
- `50-seams-and-defaults.md` (questions 3 and 4)
