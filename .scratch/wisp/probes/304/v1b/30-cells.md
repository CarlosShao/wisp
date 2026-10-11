# 304-v1b 30 — five per-cell verdicts, bonus answers, named debts

Judge leg 304-v1b, anchor `3b8066a2`. Verdicts rest on MY OWN re-runs (paths under `probes/304/v1b/**`);
the dead leg's readings are cited only where marked "from a dead leg's logs" (its `logs/a00..03`). No box flipped, no `Status:` touched, no production byte changed.

## Verdicts

- `AC#1` = **HOLDS** — per-commit ruler on the four `304-r1` commits; both spawn points changed in `e875f00f`
  (`mockllm.go:85`, `mockllm_integ_test.go:87`); reuse of `internal/proc.JobScope` (`OpenJobScope`/`StartInJob`/`Close`),
  no second job wrapper (`CreateJobObject` only in `internal/proc`); five hard constraints each checked (golden/observe/ci/Q-57 =
  0-line diffs; ticket bullet = 1 add / 0 del). Evidence: `10-ac1.md` + `logs/ac1-*.txt`.
- `AC#2` = **HOLDS** — five own shots (both spawn points, forward+reverse, every shot with its world-construction line:
  `watch_child_seen=1` + named PID/cmdline, or `world=CONSTRUCTED`; forward-after-kill counts 0; reverse-after-panic count 1 with
  parent=GONE, then stopped by three guards). Two of my shots were void and are NAMED as such. Evidence: `15-ac2-shots.md` + `logs/ac2-*`.
- `AC#3` = **HOLDS** — my own re-runs: `d22scan_rc=0`, `pathlen_rc=0` (over-budget 57 = roster 57, worst 228),
  census totals byte-equal to the frozen baseline, format two-parallel rosters 0/0 + blobs byte-identical.
  The three 0-byte files in `304/r1/logs` are empty shells: the same cell carries `build_rc=0`/`vet_rc=0` in `g-build.txt`/`g-vet.txt`.
  Evidence: `20-gates-ac3.md` + `logs/g-own-*.txt`.
- `AC#4` = **HOLDS, with one named condition** — canonical `probes/orch/gate-roster-304.md` carries the corrected line
  (mockllm.exe required + wisp.test.exe added) and the three-question rule with parent-PID demoted to corroborating (PID-reuse case stated).
  My own live run of the canonical four names: all 0, `gate_rc=0` (`logs/g-own-gatelive.txt`), attribution dump empty.
  Ruling on the orchestrator's question "(1) is adding wisp.test.exe overreach?": NOT overreach — the ticket text itself contemplates
  "另数 *.test.exe" deferred to the AC#0 ① roster, and that roster's own loose ruler already measures the `*.test.exe` family
  (from the dead leg's logs: `probes/304/a1/10-residual.md:10`). CONDITION: the fixed literal `wisp.test.exe` under-covers sibling
  `<pkg>.test.exe` names; either count the family or record the narrow scope as accepted debt.
  Ruling on "five historical template copies left unchanged + separate canonical": satisfies AC#4 — the ticket's judgement asks for the
  corrected template text plus a live run, both exist; the five copies are historical snapshots (three of them untracked, so editing them
  would rewrite history, and one is an append-only ledger file). Named caveat: those five still carry the old two-name line.
- `AC#4b` = **HOLDS** — canonical S2 covers verbatim what the ticket asks: every non-zero count points to a tree (root name + root PID),
  else "归属不明"; no bare counts into the final anchor; no "stop" debts generated from bare counts; in-flight-window caveat included.
  My live run had zero non-zero counts, so the pointer rule could not be exercised live today (stated, not hidden).

## Bonus answers (requested shape: can it be proven locally / name it)

- `Q-304-r1-2` (`go build -o` build child in no job): TRUE by construction (guard wraps only the exe via `StartInJob`;
  the build child is a plain `exec.Command(...).CombinedOutput()` at `mockllm.go:65` / `mockllm_integ_test.go:66`).
  Locally provable in principle (kill the test binary mid-build, roster `go.exe`/`compile.exe` carrying the test-name TempDir),
  NOT shot by me — named 〔本机可证、本腿未打〕.
- `Q-304-a1-B` / `Q-304-r1-4` (JobScope reach on a hosted runner; `StartInJob` Assign failure path): needs CI color
  〔需 CI 色〕. What I CAN state locally: my five shots show no Assign failure on this bench (no `t.Fatalf("start mockllm")` anywhere);
  the failure path itself is `childjob_windows.go:39-41` -> error -> `t.Fatalf` in both spawn points = hard red, not a silent leak (code read, not a measurement).

## Named debts (carried, not closed by me)

1. `Q-304-r1-1` ubuntu branch still leaks (explicit no-op, self-named in `childjob_other.go`).
2. `Q-304-r1-2` local mid-build shot not run by me (answer above).
3. `Q-304-r1-3` census "before" shot remains r1's own (I confirmed after-side; r1's before/after delta logged `IDENTICAL`).
4. `Q-304-r1-4` / `Q-304-a1-B` CI color for runner job scope.
5. `Q-304-r1-5` the other 16 helper-caller test files not re-run individually; `Q-304-r1-6` other `-timeout` tiers.
6. `-G` magpie/webview counts, `-H` per-dir attribution of the 23 `Test*` dirs, `-I` direct `os.StartProcess` roster: not measured by me (out of my scope).
7. My own two void shots (named in `15-ac2-shots.md`) and one shell-chaining accident that left a 587-byte stray `-watchkill.txt`
   in `probes/orch/20261011/txt/` (left in place per create-only rule; needs orchestrator deletion).
8. The AC#4 five historical copies remain stale by design; the wisp.test.exe literal-vs-family gap is the AC#4 condition.
