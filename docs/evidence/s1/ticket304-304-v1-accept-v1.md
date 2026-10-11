# Ticket 304 — non-implementer adjudication (leg `304-v1b`), five cells `AC#1`/`AC#2`/`AC#3`/`AC#4`/`AC#4b`

Judge: leg `304-v1b` (continues dead leg `304-v1`; inherits on-disk leftovers only). Date: 2026-10-11 evening.
Anchor for every table header above: `git rev-parse --short HEAD` = `3b8066a2` (own reading; leg anchor commit `83b1aea9`).
Nothing here flips a ticket box; the orchestrator reruns these numbers first.

| cell | verdict | basis (my own re-runs) | evidence |
|---|---|---|---|
| `AC#1` | HOLDS | per-commit ruler over the four `304-r1` commits (`e875f00f`/`bd123485`/`27f7b197`/`04f2ecbc`): both spawn points in one batch (`mockllm.go:85`, `mockllm_integ_test.go:87`); reuse of `internal/proc.JobScope` only (`KILL_ON_JOB_CLOSE`), no second job wrapper; five hard constraints each = 0 forbidden diffs; ticket bullet 1+/0- | `probes/304/v1b/10-ac1.md`, `logs/ac1-*.txt` |
| `AC#2` | HOLDS | 5 own shots in out-of-repo trees (6-piece archive, `go_files=682`): FWD (b)/(c) openaichat + FWD/REV helper (c) + REV (b) openaichat; every shot carries world-construction evidence (child PID/cmdline while parent alive; `panic: test timed out after 3s`=1 or `kill_testbinary ... parent_alive_at_kill=YES`); forward post-counts 0; reverse residual count 1 (named pid, parent=GONE) then stopped by three guards, final 0. 2 own void shots named | `probes/304/v1b/15-ac2-shots.md`, `logs/ac2-*` |
| `AC#3` | HOLDS | own re-runs: `d22scan_rc=0`; `pathlen_rc=0` (over-budget 57 = roster 57, worst full path 228); census totals byte-equal to frozen baseline (`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`); format worktree roster 0 / HEAD-blob roster 0, blobs byte-identical. `304-r1` 0-byte trio = empty shells (readings live in `g-build.txt build_rc=0` / `g-vet.txt vet_rc=0` / `g-summary.txt`) | `probes/304/v1b/20-gates-ac3.md`, `logs/g-own-*` |
| `AC#4` | HOLDS with one named condition | canonical `probes/orch/gate-roster-304.md`: corrected process-name line (mockllm.exe required, wisp.test.exe added) + three questions with parent-PID demoted to corroborating (PID-reuse case written); my live run of the canonical four names: all 0, `gate_rc=0`, attribution dump empty. Adding wisp.test.exe = NOT overreach (ticket defers `*.test.exe` to the AC#0 roster; that roster's loose ruler already measures the family). Condition: literal `wisp.test.exe` under-covers sibling `<pkg>.test.exe` names. Five historical copies left unchanged = satisfies AC#4 (snapshots; ticket asks for corrected text + live run, both exist) | `probes/304/v1b/30-cells.md`, `logs/g-own-gatelive.txt`, `logs/g-own-templates.txt` |
| `AC#4b` | HOLDS | canonical S2 covers "every non-zero count points to a tree, else 归属不明", "no bare counts into final anchors", "no stop-debts from bare counts", plus the in-flight-window caveat — matching the ticket text fully. My live run had zero non-zero counts (stated) | `probes/orch/gate-roster-304.md` S2, `logs/g-own-gatelive.txt` |

Mandated controls handled: the orchestrator's two no-world readings (their file S1) are void and independently named;
the three 0-byte files in `304/r1/logs` are resolved as shells with rc-bearing companions in the same cell.
Bonus (named in `30-cells.md`): `Q-304-r1-2` = true by construction, locally provable but not shot by me;
`Q-304-a1-B` / `Q-304-r1-4` = 〔需 CI 色〕, locally only "no Assign failure observed in 5 shots" can be said.
Debts: see `probes/304/v1b/30-cells.md` named list (incl. my stray `-watchkill.txt` left in the orchestrator txt dir per create-only rule).
