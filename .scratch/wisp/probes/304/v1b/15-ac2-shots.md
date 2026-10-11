# 304-v1b 15 — my own AC#2 paired readings (both spawn points, world-construction evidence in every shot)

Judge leg 304-v1b. Anchor `3b8066a2`. Trees: T1 = forward export tree `/d/tmp/wisp-304v1b-fwd-UtGhQF`
(682 .go files, guard present), T2 = reverse export tree `/d/tmp/wisp-304v1b-rev-qCcWkx` (682 .go files,
two spawn files restored to `b9f37892` pre-guard form: `mockllm.go:79` and `mockllm_integ_test.go:80` are
plain `proc.cmd.Start()`, `grep -c StartChildGuarded` = 0 0; setup log `logs/ac2-rev-setup.txt`).
Six-piece export template used for both (`go.mod go.sum internal cmd tools scripts`), non-empty counted first.
Attribution in every roster = named path in the command line (test-name TempDir + tree golden dir); never counts, never parent-PID existence.
Only processes I started were touched; `msedgewebview2.exe` never named.

## Shots (each with its own "the bad world was really constructed" line)

1. FWD (b) openaichat spawn point, T1: `logs/ac2-fwd-b-head.txt` + `-body.txt`.
   World: `hit pid=5808 ppid=18332 parent=ac2t.exe` watch_seen=1 at poll 2; body line 4 = child serving with
   `golden dir D:\tmp\wisp-304v1b-fwd-UtGhQF\internal\llm\testdata\golden`. Bad world: `panic: test timed out after 3s` count=1
   (stacks in `mockllm_integ_test.go:188/249/217`), `go_test_rc=2`. Post roster `count_tree=0` => job reaped the child although the panic skipped every cleanup.
2. FWD (c) openaichat spawn point, T1: `logs/ac2-fwd-c-head.txt`.
   World: child `pid=28692`, watcher saw it while parent alive. Then `kill_testbinary pid=27132 parent_alive_at_kill=YES kill_targets=1`;
   body stops mid-subtests (no PASS). Post roster `count_tree=0` => OS killed the member with the test binary (AC#2's extra sentence satisfied: not "test finished so it looks clean").
3. REV (b) openaichat spawn point, T2 (guard removed): `logs/ac2-rev-b-head.txt` + `logs/ac2-rev-kill.txt`.
   World: child `pid=24372 ppid=32404 ac2t.exe` seen at poll 2; `panic: test timed out after 3s` count=1; go_test_rc=2.
   Post roster = `hit pid=24372 ppid=32404 parent=GONE ... count_tree=1` (SAME pid the watcher saw). => residual comes back.
   Stopped by me with three guards: `killed pid=24372 guards=name=mockllm.exe+tree-path+parent-gone`, `count_tree_after_kill=0`.
4. FWD (c) helper spawn point (`adaptertest.StartMockllm`, test `TestFailoverAfterPrimaryExhaustion`), T1: `logs/ac2-fwd-helper3-c-watchkill.txt`.
   World: `world=CONSTRUCTED child_pids=31768 t=0.9`; `kill_testbinary pid=11408 parent_alive_at_kill=YES`; go_test_rc=127 (killed mid-run).
   Post `count_tree=0`. Single-process watcher `wk.ps1` (polls 40 ms, kills within the same PS process).
5. REV (c) helper spawn point, T2 (guard removed): `logs/ac2-rev-helper-c-watchkill.txt`.
   World: `world=CONSTRUCTED child_pids=1852 t=1.58`; `kill_testbinary pid=18428 parent_alive_at_kill=YES`; go_test_rc=127.
   Post: `hit pid=1852 ppid=18428 parent=GONE count_tree=1` => residual; stopped by guards; `count_tree_after_kill=0`.

## Void shots (mine, named so nobody counts them)

- `logs/ac2-fwd-helper-c-head.txt` (2026-10-11 21:43): `watch_child_seen=0`, `go_test_rc=0`, PASS — test finished before the watcher caught the child => VOID (no world).
- `logs/ac2-fwd-helper2-c-head.txt` (21:44): watcher saw a child but `kill_targets=0`, test PASSed 2.58s => VOID for shape (c) (kill lottery lost; that shot says nothing).
- One inline attempt (21:45) was destroyed by my own shell chaining (whole chain backgrounded, vars empty); the T1 side ran unwatched => VOID; it also left a stray 587-byte `-watchkill.txt` in the orchestrator txt dir (named in the debt list; not deleted per fleet rule).

## The two controls the orchestrator buried (handled, not cited as evidence)

- Their two readings that "never started a world" (`verify-304r1-309r3.md` S1 rows 2 and 4) are void in their own file;
  none of my five shots above relies on them. My own rule applied: any shot with `watch_child_seen=0` is VOID (cases 4th bullet above).
- The three 0-byte files in `304/r1/logs` are resolved in `20-gates-ac3.md` (rc-bearing companions exist).

## Verdict

`AC#2` = **HOLDS** on my own paired re-runs: guard present => child reaped under both death shapes on both spawn points;
guard removed (byte-restored pre-guard files in a paired out-of-repo tree) => residual measured back with named PID and cmdline, then stopped by guards.
