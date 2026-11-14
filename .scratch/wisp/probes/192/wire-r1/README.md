# probes/192/wire-r1 - ticket 192's wiring leg

The three `cmd/wisp` steps that make the landed preview stack reachable, and the
readings each one has to produce. Everything here re-runs on this machine; none
of it needs a window (the filter's lookup and the composer dispatch are both
plain functions).

## Files

| file | what it is |
|---|---|
| `01-callsites-before-after.sh` | the call-site ruler: `QuoteForModel`, `panel.NewPreviewResolver` and the preview branch in `cmd/wisp`, before (`HEAD`) and after (worktree, `--untracked` so the new production file is counted) |
| `logs/01-callsites.txt` | that script's output (0 -> 1, 0 -> 1, 0 -> 4) |
| `logs/02-d22scan.txt` | `tools/d22scan/d22scan.exe -root .` |
| `logs/03-vet-gofmt.txt` | `go vet ./cmd/wisp/` + `gofmt -l` + `gofumpt -l` |
| `logs/10-cmdwisp-full.txt` | the whole `./cmd/wisp` suite through `scripts/wisp-cli-staged-tests.sh` (the carrier; bare `go test` fails on this box with `0xc0000135`), 630 cases in 39m15s |
| `logs/11-t219-alone.txt` | the one red the dispatch did not name, re-measured alone (PASS 1.70s after 55.78s in the loaded run) |

The step-1/step-2/step-3 readings themselves are `t.Logf` lines inside
`cmd/wisp/panel_preview_192_windows_test.go`; they are in
`logs/10-cmdwisp-full.txt` (grep `step 1`/`step 3`, and the six
`--- PASS: Test192...` lines) and in `docs/evidence/s1/192-preview-wiring-r1.md` §2.

## The one trap

`git grep` skips untracked files. The whole point of this leg is a NEW
production file, so a reading taken without `--untracked` answers 0 both before
and after and looks like "nothing changed". The script passes it, and the
evidence quotes the line (`cmd/wisp/panel_preview_192.go:294`).
