# 433 — `TestAC246DevLegIgnoresTheTestTaskInjection`: the TempDir cleanup red

Ticket: `.scratch/wisp/issues/433-the-246-dev-leg-test-is-stably-red-on-its-own-tempdir-cleanup-failing-because-a-handle-outlives-the-child.md`
Write surface: `cmd/wisp/resident_task_source_246_windows_test.go` (harness side only, no product code)
Measured: 2026-11-12 22:25–22:46 +08, host `swq`, desktop otherwise idle, no parallel legs.

---

## §1 AC#0 — the readings (read-only, nothing written)

### §1.1 The complete failing filename

The ticket's log truncates at `EB…`. Rerun 1 supplies the whole line:

```
testing.go:1617: TempDir RemoveAll cleanup: unlinkat C:\Users\swq\AppData\Local\Temp\TestAC246DevLegIgnoresTheTestTaskInjection2168101192\002\wisp\orb\62760-905e797d\boot-1\EBWebView\lockfile: The process cannot access the file because it is being used by another process.
```

So the object is **`EBWebView\lockfile`** — the WebView2 user-data directory's lock file,
under the ball's per-boot leaf `orb\<pid>-<8hex>\boot-1\`, inside the throwaway `APPDATA`
that the case hands the child as `LOCALAPPDATA`.

### §1.2 Rerun counts and durations — the ticket's "stably red" is **wrong**

| # | command | result | duration |
|---|---|---|---|
| 1 | carrier, `-test.run '^TestAC246DevLegIgnoresTheTestTaskInjection$'` | **FAIL** | 23.31s |
| 2 | same, fresh stage dir | **PASS** | 13.42s |
| 3 | instrumented overlay (below) | **FAIL** | 11.24s |
| 4 | instrumented overlay | **FAIL** | 13.13s |
| 5 | instrumented overlay, graceful-exit mode | **FAIL** | — (cleanup) |
| 6 | instrumented overlay, RemoveAll probe | **PASS** | 15.75s |

**Correction to the ticket:** it is **not stable**. It is a **race that fails most of the
time but not always** — run 2 above passed with no change to the code. The ticket's own
"stably red, not occasional" (`:23`) and the dispatch's "稳定红" (`:8`) are both
contradicted by a measured green on an idle desktop. This matters for the fix: a
sometimes-green test is exactly what a 2s-vs-5.58s budget race looks like.

### §1.3 The moment of cleanup: had the child exited? — **YES, it had**

From the overlay instrument (`resident_task_source_246_windows_test.go:492` in the
overlay copy, never committed):

```
433 INSTR: shutdown returned in 54ms waited=true waitErr=exit status 1 childPID=38616 childAlive=false
433 INSTR: t+0 1 lockfile(s) under the temp root, 1 still held: [...\boot-1\EBWebView\lockfile]
```

`childAlive=false` is read by opening the pid with `PROCESS_QUERY_LIMITED_INFORMATION` and
comparing `GetExitCodeProcess` against `STILL_ACTIVE` — a direct process-table reading,
not an inference from `cmd.Wait`. **The child was gone and the handle was still held.**

### §1.4 Who actually holds it

```
433 INSTR: t+0 survivor 22060 <- parent 38616 at 22:38:17 :: msedgewebview2.exe --embedded-browser-webview=1 --webview-exe-name=wisp.exe --user-data-dir="...\boot-1\EBWebView"
433 INSTR: t+0 survivor 44872 <- parent 22060 :: --type=crashpad-handler --user-data-dir=...\boot-1\EBWebView
433 INSTR: t+0 survivor 63112 <- parent 22060 :: --type=gpu-process ...
433 INSTR: t+0 survivor 62968 <- parent 22060 :: --type=utility ...NetworkService
433 INSTR: t+0 survivor 18100 <- parent 22060 :: --type=utility ...StorageService
```

Five `msedgewebview2.exe` processes, **orphaned by the dead child** (their parent is
`22060`, the browser process; the browser's parent is the dead `wisp.exe`, pid 38616), all
pointing `--user-data-dir` at the very tree being removed. One of them holds
`EBWebView\lockfile` open, which is what `unlinkat` reports.

### §1.5 How long the handle outlives the child — the number that decides the fix

Measured against the predicate the framework actually uses (`os.RemoveAll` on the temp
root, not "can the lockfile be opened" — a delete-pending file answers
`ERROR_FILE_NOT_FOUND` forever and would have been a false reading):

```
433 INSTR: t+75ms   RemoveAll still fails: ...lockfile ... being used by another process
433 INSTR: t+1.095s RemoveAll still fails
433 INSTR: t+2.111s RemoveAll still fails
433 INSTR: t+3.131s RemoveAll still fails
433 INSTR: t+4.150s RemoveAll still fails
433 INSTR: t+5.168s RemoveAll still fails
433 INSTR: tree removable at t+5.576s after 55 tries
```

**The tree becomes removable 5.58s after the child exits.**

The Go framework's budget, read at `D:\work\base\go\src\testing\testing.go`:

- `:1615-1619` — the `t.TempDir` cleanup calls `removeAll(c.tempDir)` and reports failure via `c.Errorf`
- `:1649-1668` — `removeAll` retries only while `isWindowsRetryable(err)` **and** the elapsed time is under `arbitraryTimeout`
- `:1650` — `const arbitraryTimeout = 2 * time.Second`

**5.58s needed, 2s granted.** That gap is the entire red.

### §1.6 `leg.stop()` / `leg.exitedWithin()` — do they guarantee the child exited?

- `resident_sink_nail_127_windows_test.go:290-301` — `stop()` reaps if already done,
  else `Process.Kill()`, then waits up to 200 ticks (10s) for the `cmd.Wait` result.
- `resident_sink_nail_127_windows_test.go:265-272` — `exitedWithin` polls `reap` for N×50ms.
- `resident_sink_nail_127_windows_test.go:275-287` — `reap` takes the `cmd.Wait` result.

**Yes, they guarantee the *child* exited** — and the instrument confirms it (`childAlive=false`,
`waited=true`). **They do not, and cannot, guarantee the child's WebView2 subtree exited.**
Killing `wisp.exe` orphans the browser process rather than killing it, because the ball's
WebView2 tree is not in the Job Object that would take it down (`internal/proc/jobscope_windows.go:20,77`
sets `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, but that scope is the runtime's, and the browser
process is spawned by WebView2, not assigned to it).

### §1.7 The judgment: **乙 (WebView2 handle), not 甲, not 丙**

**乙 — a WebView2 handle outlives the child.** One sentence: the child is provably gone
(`childAlive=false`) while its orphaned `msedgewebview2.exe` still holds
`EBWebView\lockfile` for 5.58s, which is past the framework's 2s retry.

Why not the other two:

- **甲 (child not really exited)** — refuted by measurement: `childAlive=false` at t+0, and
  `stop()`/`exitedWithin()` do join `cmd.Wait` (`resident_sink_nail_127_windows_test.go:290-301`).
- **丙 (child cwd/temp points into the tree)** — refuted by reading: the child's `cmd.Dir`
  is `filepath.Dir(exe)` (`resident_task_source_246_windows_test.go:493`), i.e.
  `...\<test>\001` (the build dir), while the tree being removed is `...\<test>\002`
  (the `APPDATA`). The instrument printed both: `childCWD="...\001"`,
  `appData="...\002"`. Different subtrees; the cwd is not what blocks the unlink. (Note
  that `t.TempDir` cleanup would have removed `\001` too and it did not complain — only
  `\002` did.)
- The culprit path is **inside** the tree because `LOCALAPPDATA` is pointed at it
  (`:459`), and the ball roots its WebView2 profile at `LOCALAPPDATA` —
  `internal/ball/orb_windows.go:471-477` (`orbDataBase` reads `LOCALAPPDATA`) and
  `:529-535` (`orbBootDataDir` appends `boot-<n>`), with `:315` assigning
  `browser.DataPath = bootDir`. That is the chain that puts `EBWebView` under the
  managed directory.

### §1.8 This is **not** a desktop-focus case

Confirmed as the ticket requires. The case starts one child (`exec.Command(exe)` with no
args, `:492`) and never asks for foreground, focus, or a real window owner; it reads
stdout and a log sink. It is not `TestAC4Focus*`, and its red is reproducible with the
desktop idle. The two are unrelated.

---

## §2 AC#1 — the fix, and the mutation control

### §2.1 What was changed (one place)

`cmd/wisp/resident_task_source_246_windows_test.go`, the `appData` line of
`TestAC246DevLegIgnoresTheTestTaskInjection` only:

```go
-	appData := t.TempDir()
+	appData := managedDir433(t, "ac246-dev-leg-appdata")
```

plus the `managedDir433` helper (registered cleanup, bounded by **attempt count**, not a
wall-clock deadline — D42#9; same shape as `exitedWithin(attempts)` at
`resident_sink_nail_127_windows_test.go:265`).

The path itself is `uniqueProfileLeaf` (`panel_resident_windows_test.go:61`), the repo's
**existing** answer to this exact problem for the panel's WebView2 profile
(`panel_resident_windows_test.go:55-63`, and the reasoning at
`panel_reshow_255r8_winlive_test.go:118-125`). This is not a new invention; it is the
house rule the panel already follows, applied to the one resident case that was still
borrowing `t.TempDir`.

**Not the forbidden fix.** The forbidden shape is "replace the managed temp dir with a
hand-made one that is never cleaned". Here the removal is registered in `t.Cleanup`, so it
runs on **every** exit path including `t.Fatalf`; it is retried with a real budget; and if
it still fails it **fails the test loudly** with its own named message instead of
disappearing. Nothing is hidden — the tree is still swept every run. What changed is
*which* cleanup owns the wait, and how long it is allowed to wait.

### §2.2 Mutation control — really run

| step | command | reading |
|---|---|---|
| fix in place | carrier | **PASS 14.48s** |
| **mutate**: `appData := t.TempDir()` restored, helper left unused | carrier | **FAIL 9.82s** — `testing.go:1617: TempDir RemoveAll cleanup: unlinkat …\boot-1\EBWebView\lockfile: The process cannot access the file because it is being used by another process.` |
| **restore** | carrier | **PASS 8.30s** |
| restore, repeat | carrier | **PASS 13.48s** |

The mutation reproduces the **exact original signature**, byte for byte, including the same
file (`EBWebView\lockfile`). The fix and the defect are one line apart.

### §2.3 The whole AC246 family, with the fix in place

`-test.run '^TestAC246'` → **PASS**, 0 failures, 14 cases including the three shipped-leg
cases (7.37s / 8.01s / 8.70s / 15.18s). No sibling regression.

---

## §3 AC#2 — boundary check

`git diff --name-only` (mine is the only line I own; the rest pre-date this leg and are
another agent's work in the shared tree):

```
.scratch/wisp/issues/426-…md            (pre-existing, not mine)
.scratch/wisp/issues/427-…md            (pre-existing, not mine)
.scratch/wisp/issues/428-…md            (pre-existing, not mine)
.scratch/wisp/issues/429-…md            (pre-existing, not mine)
.scratch/wisp/issues/430-…md            (pre-existing, not mine)
.scratch/wisp/issues/431-…md            (pre-existing, not mine)
.scratch/wisp/issues/432-…md            (pre-existing, not mine)
.scratch/wisp/probes/291/mic/…_test.go  (pre-existing, not mine)
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt (pre-existing, not mine)
cmd/wisp/resident_task_source_246_windows_test.go          <- MINE
frontend/fixtures/composer-states.html  (pre-existing, not mine)
```

None of `docs/PLAN.md`, `docs/specs/**`, `internal/observe/thresholds.go`, golden files,
`tools/d22scan/allowlist.txt`, `frontend/**` or `design/**` was touched by this leg.
`build/mut_rvw.go` **still present** (verified `Test-Path` → `True`), not deleted.

---

## §4 AC#3 — the four gate readings

| gate | command | reading |
|---|---|---|
| build | `go build ./cmd/wisp/` | exit 0, no output |
| vet | `go vet ./cmd/wisp/` | exit 0, no output |
| format | `go run mvdan.cc/gofumpt@v0.7.0 -l cmd/wisp/resident_task_source_246_windows_test.go` | **zero output** (clean) |
| d22scan | `scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations`; `runtests.sh: OK - top-level: PASS=45 FAIL=0 SKIP=0`; examined 480 production Go files; exit 0 |

(`gofumpt` is not on PATH on this host; run from the module cache at `mvdan.cc/gofumpt@v0.7.0`.)

---

## §5 AC#4 — relationship to ticket 246

This ticket fixes **only** the cleanup defect of one case. It does **not** flip any AC box
of ticket 246, and it does **not** claim the `TestAC246*` family is green as a body of
evidence — §2.3's family run is a regression check for this change, not a re-acceptance of
246's AC#7. The product was not modified: `internal/ball` still roots the orb profile at
`LOCALAPPDATA`, which is correct for the shipped binary.

## §6 What this ticket does **not** fix

- **The product-side question is left open on purpose.** A killed resident process orphans
  its `msedgewebview2.exe` subtree for ~5.6s. For a user that is invisible (the next boot
  sweeps dead leaves via `orbSweepStaleDirs`, `internal/ball/orb_windows.go:544-569`), and
  changing it would mean assigning WebView2's browser process to the Job Object, which is
  a product change this ticket has no evidence to justify. Recorded here, not fixed here.
- `TestAC4FocusReturnToPriorWindowGap33r5` — untouched; it is a desktop-focus case per the
  ticket and needs its own scheduling ticket.
- `thresholds.go` / SLO thresholds — untouched (frozen).
