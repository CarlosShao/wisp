# gate-rerun-1 — ticket 234: 12 cells that only need a reading (mutation + gate rerun on current HEAD)

- Leg: `gate-rerun-1` (non-implementor rerun). Ticket: `.scratch/wisp/issues/234-the-twelve-cells-that-only-need-a-reading-mutation-and-gate-rerun-on-current-head.md`
- Anchor self-taken at start: `git rev-parse --short HEAD` = `015f6be1`, branch `dev`.
- Blocker check: ticket face Status requires ticket 222's acceptance leg to finish first.
  Evidence it finished: `d39f730e` (222-v1 shou-gong), `f23586fa` (222-v1 table + raw readings),
  `015f6be1` (ledger A447 collects 222-v1). Block released before this leg started.
- Scope: zero production code. `frontend/**` and `design/**` never read, never listed, never quoted.
- Order of work: 5 mutation cells -> 7 gate/regression cells -> AC#5 ledger relocation -> AC#6 bad yardstick.
  One commit per cell so that a mid-run death only loses the remaining cells.

## Baseline roster (AC#1)

- UNJUDGED

## Mutation cells (5)

### M-1 ticket 92 `:68` — AC#5 direction (iii), verifier marked [self-report only]

- UNJUDGED

### M-2 ticket 104 `:52` — AC#3 both directions

- UNJUDGED

### M-3 ticket 110 `:39` — AC#3 the new step must go red by itself

- UNJUDGED

### M-4 ticket 113 `:54` — AC#4 two shapes

- UNJUDGED

### M-5 ticket 115 `:58` — AC#4 two shapes

- UNJUDGED

## Gate / regression cells (7)

### G-1 ticket 92 `:73` — AC#6 three readings valid only against two old tree snapshots

- UNJUDGED

### G-2 ticket 97 `:55` — AC#5 per-package five numbers re-taken on current HEAD

- UNJUDGED

### G-3 ticket 104 `:57` — AC#4 four packages `-count=2 -v`

- UNJUDGED

### G-4 ticket 105 `:54` — AC#5 five readings + per-scope file counts

- UNJUDGED

### G-5 ticket 110 `:47` — AC#5 `bash -n` / `d22scan.sh` per scope not lower / new step ordering

- UNJUDGED

### G-6 ticket 113 `:56` — AC#5 in-container three packages four numbers + per-package vet + gofmt + `d22scan.sh`

- UNJUDGED

### G-7 ticket 115 `:64` — AC#6 per-package four numbers / gofmt / gofumpt.exe -l / go vet / d22scan

- UNJUDGED

## Old debts carried by 5 of the cells

- UNJUDGED

## Final roster comparison (AC#4)

- UNJUDGED

## AC#5 — 10 duplicate denominators relocated to still-open tickets

- UNJUDGED

## AC#6 — retire one bad yardstick (`NewGate(` zero hits)

- UNJUDGED

## Not finished / not judgeable

- (must not be left empty at the end of the leg)
