# `.scratch/wisp/probes/161/r1` — D22 gate blind-spot bench (ticket 161 AC#1)

Everything here is a **measurement bench**. It writes zero production bytes:
`tools/d22scan/**` is only ever read and compiled.

## What it does

`bench.sh` builds, for every numbered ban in `tools/d22scan/main.go`'s `// Bans`
roster, a **known-violation** fake repository root and a **known-clean** one of
the same shape, runs the scanner over it with `-root`, and records whether that
ban rang. The clean twin is load-bearing: a sample that rings proves nothing
unless the same-shaped clean sample stays silent.

```
sh bench.sh all          # build the binary, run every cell, print logs/summary.tsv
sh bench.sh cells        # the cell roster
sh bench.sh run <cell>   # one cell: build + scan -> logs/<cell>.scan.log
sh bench.sh gorun <cell> # the same cell through `go run .` (binary-freshness cross-check)
sh bench.sh tool         # rebuild bin/d22scan-probe161.exe from tools/d22scan
```

## Layout

| path | what | tracked? |
|---|---|---|
| `bench.sh` | the whole bench: skeleton + 27 payloads + driver | yes |
| `.gitignore` | keeps `runs/` and `bin/` out of the repo | yes |
| `runs/<cell>/` | generated fake repo roots (go.mod, allowlist.txt, internal/, cmd/, frontend/, design/) | **no** |
| `bin/` | the scanner binary built from `tools/d22scan` to run 27 cells fast | no |
| `logs/<cell>.scan.log` | verbatim command + full scanner output + `# rc=N` | yes |
| `logs/summary.tsv` | one row per cell: ban, expectation, rc, ring count, finding lines | yes |
| `logs/pre-baseline-*` `logs/post-*` | the §6 gate readings this ticket's evidence file quotes | yes |

## Why the fixtures are whole fake roots

`main()` -> `checkRoot()` refuses to print a verdict for a tree that is not a
wisp root (it wants `go.mod`, `tools/d22scan/allowlist.txt`, `internal/`, `cmd/`
and >= 10 production `.go` files), and `verdict()` exits **2** if any declared
live scope walked zero files. So a sample must satisfy both or the run answers
"empty instrument" instead of the question. `skeleton-clean` is the control that
proves the layout itself is silent, so every other silence is attributable to
the payload and not to the scaffolding.

Because a fixture root sits **inside** this repository, `runGitIndex()` detects
`rev-parse --show-prefix` != "" and applies **no** ignore rule at all, printing
`gitignore rules NOT APPLIED ...` on every cell run. That is the loud direction:
nothing in a fixture can be silently skipped.

## Reading the table

`ring_count` counts lines shaped like `path:N: [ban] ...`. `said_ban` is whether
the ban this cell is about is among them. For probe cells the expectation column
deliberately says `unknown`: a probe measures, it does not assert.
