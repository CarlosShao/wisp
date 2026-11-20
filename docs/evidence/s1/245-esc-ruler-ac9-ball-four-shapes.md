# 245 AC#9 evidence — the two-process Esc delivery ruler, standing in `internal/ball`

- Leg: `245-ac9-ruler` (write leg, not the adjudicator).
- Dispatch: `.scratch/wisp/probes/orch/20261120/txt/245-ac9-dispatch.txt` (committed by the orchestrator as `f62e312a9`).
- Ticket: `.scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md`, AC#9 (its turn was made on 2026-11-14 from `245-v1`'s second-process rig; ledger A480).
- Deliverables:
  - `internal/ball/esc_delivery_ruler_245_windows_test.go` (new, `//go:build windows`, default tier so the dispatch's gate command selects it).
  - This evidence file.
- Product bytes changed: **zero**. `git diff --quiet -- internal/ball/ball_windows.go internal/ball/hotkey_windows.go` rc=0 after both teeth; SHA-256 of `hotkey_windows.go` = `C66CE2E336A8402BA7389C272C1177B80408ED40C9B39D484E7719B7DE765C69` (same string the ticket already recorded) and `ball_windows.go` = `FA56C8AD2ADE15F77FEF925E531DAF5FCF80DCD31392297EEC1AA970F28FEDD8`, both re-verified after the final restore.

## 1. What the ruler is

`TestLive245EscDeliveryRulerFourShapes` runs the acceptance-side four shapes against the **ball itself**, on one real hidden ball with the production default bindings (bare-`Esc` cancel), plus the instrument's own two controls, with **no human driving anything**:

- The second real GUI process is this test binary exec'd in the classic helper shape (`exec.Command(os.Args[0], "-test.run=^TestHelperProcess$")` + env `WISP_245_ESC_DELIVERY_HELPER=watch|steal`, the same re-exec pattern `orb_datapath_windows_test.go` already uses). The helper (`TestHelperProcess`, this file only) creates a real top-level window (system class `STATIC`, subclassed), takes the foreground itself (bounded settle; SPI foreground-lock timeout zeroed and **restored**), waits for the trigger over stdin, injects one physical Esc through `keybd_event`, counts the `WM_KEYDOWN`/`WM_SYSKEYDOWN`/`WM_HOTKEY` it was delivered, prints one `READING` line. Single-threaded on a locked OS thread; the only clock is the Win32 quit timer (2600 ms) plus a 60 s process watchdog, so a dead harness can never read as a pass.
- Order on one timeline: E1 baseline (nothing of ours) -> E2 control (observer registers bare Esc itself) -> ball up -> shape2 (in-process roster) -> shape1 (idle delivery) -> shape3 (card waiting: borrow + veto) -> shape4 (card gone: return).
- Verdict discipline carried from AC#8: a "0" is first attributed through **the ball's own cancel gesture** (a leaked registration delivers `WM_HOTKEY` to OUR window and only ours — the discriminator a bookkeeping-only release cannot fool), then through a Win32 claim probe (NULL-window `RegisterHotKey`, 3 attempts, verbatim errno). ours = hard RED; foreign = `SKIP-LOUD` with attempts and errno; nobody = RED (instrument). Instrument unusable (no desktop / foreground refused) = `SKIP-LOUD` named with the numbers measured.
- Desktop-unavailable paths are named skips (`SKIP-LOUD (ticket 245 AC#9, NO INTERACTIVE DESKTOP ...)` on helper `create_window` failure, `... FOREGROUND REFUSED ...` when the observer never owned the foreground at injection time). `ball.New` failing **after** E1/E2 proved an interactive desktop is a RED (PRODUCT), not a skip. None of these fired on this machine today; every reading below was really taken.

## 2. The readings (gate command, three green runs today, last one verbatim)

Command: `go test ./internal/ball -count=1 -run "245|Esc" -v` — rc=0, `ok github.com/CarlosShao/wisp/internal/ball 15.579s`, 8/8 PASS, 0 FAIL, 0 SKIP. Full verbatim log: `.scratch/wisp/probes/245/ac9-ruler/post-restore-gate.txt` (also `final-gate.txt` for the run before the gofumpt re-wrap).

The ruler's final line (verbatim):

```
    esc_delivery_ruler_245_windows_test.go:843: AC#9 ruler four shapes: E1 baseline=1 E2 control=0/1 shape1_idle=1 shape2_roster=live3/cancel=standby shape3_borrow=observer0+veto1 shape4_return=1
--- PASS: TestLive245EscDeliveryRulerFourShapes (15.46s)
```

| shape | reading | evidence line |
|---|---|---|
| E1 instrument baseline (nothing of ours) | observer `keydown_esc=1` | `final-gate.txt` / `post-restore-gate.txt` |
| E2 instrument control (observer self-claims Esc) | `keydown_esc=0` and `wm_hotkey=1` | same |
| ① idle ball, other process receives Esc | observer `keydown_esc=1`, ball veto delta=0 | `:751 AC#9 shape1 (idle delivery): observer keydown_esc=1, ball veto delta=0` |
| ② steady state, no bare Esc bound | `HotkeyReport().Bindings()` live=3 (summon Ctrl+Alt+Q / mute Ctrl+Alt+M / panel Ctrl+Alt+P), cancel=standby, bare-held=false | `:724 AC#9 shape2 (steady state): Bindings() live=3 cancel=standby bare-held=false (summon/mute/panel)` |
| ③ card waiting, Esc borrowed | observer `keydown_esc=0` AND this ball's cancel gesture fired once (`veto delta=1`, real `WM_HOTKEY` through the ball's wndproc) | `:794 AC#9 shape3 (borrow + veto): observer keydown_esc=0, ball veto delta=1 (the cancel gesture fired)` |
| ④ card removed, Esc handed back | observer `keydown_esc=1` again, ball veto delta=0 | `:841 AC#9 shape4 (return): observer keydown_esc=1, ball veto delta=0` |

Shape3's in-process pre-reads on the same run: cancel line live with the production default pair `{Mods:16384 VK:27}` (`MOD_NOREPEAT`, `VK_ESCAPE`), `EscTakenOver()=true`, `Live()` holds 4; shape4's: `EscTakenOver()=false`, cancel standby, `idleRosterReds245` empty.

The other 7 selected tests (existing 245 nails, unchanged): `TestConfiguredCancelIdleBindingFollowsModifierRuling245`, `TestAC8IdleProbeSplitsForeignOwnerFromOurOwnSet245`, `TestAC6IdlePinJudgesTheBareKeyNotTheCount245`, `TestRefusedBorrowReleaseReturnsTheIdleRoster245`, `TestRefusedBorrowStaysNamedWhileTheCardWaits245`, `TestHeldBorrowReleaseStillReturnsStandby245`, `TestDefaultHotkeysIdlePassHoldsNoEsc` — all PASS, none loosened.

## 3. The teeth (mutation proof, both re-run against the final formatted bytes)

### Tooth A — the return step stops releasing => shape4 must red

Mutation (product `internal/ball/ball_windows.go`, `ReleaseEscAfterSession`): the `releaseEsc(b.hwnd)` call removed while the bookkeeping below it still runs (flags reset, report reset to standby) — the most deceptive form of "不 release": the report says returned, the key stays in Win32.

Verbatim red (`.scratch/wisp/probes/245/ac9-ruler/tooth-A-final.txt`, rc=1):

```
    esc_delivery_ruler_245_windows_test.go:819: ticket 245 AC#9 shape4 RED (PRODUCT): the return did not hand the bare Esc back: the observer received keydown_esc=0 (want 1) and this ball's cancel gesture fired after the release (1 time(s) on this shape) - that gesture can only come from a registration on this ball's own window, so the release dropped the bookkeeping but not the key. Report after the return:
        [{Name:summon ID:1 Binding:Ctrl+Alt+Q Status:live Acc:{Mods:16387 VK:81} Err:<nil> Note:} {Name:mute ID:2 Binding:Ctrl+Alt+M Status:live Acc:{Mods:16387 VK:77} Err:<nil> Note:} {Name:cancel ID:3 Binding:Esc Status:not bound while idle (cancel is borrowed only during Confirming) Acc:{Mods:0 VK:0} Err:<nil> Note:} {Name:panel ID:4 Binding:Ctrl+Alt+P Status:live Acc:{Mods:16387 VK:80} Err:<nil> Note:}]
--- FAIL: TestLive245EscDeliveryRulerFourShapes (16.46s)
```

Shapes 1-3 stayed green under this mutation (only the return is broken), and the catch came from the **delivery half** (observer 0 + our own gesture firing), not from a flag — exactly the zero-symptom shape a bookkeeping-only fake would need to survive.

### Tooth B — the steady state binds the bare Esc again => shape2 must red

Mutation (product `internal/ball/hotkey_windows.go`, `cancelHeldWhileIdle`): `return acc.Mods&realModifiers != 0` replaced by `return acc.VK != 0`, i.e. the pre-AC#6② shape where the idle pass registers the bare cancel binding.

Verbatim red (`.scratch/wisp/probes/245/ac9-ruler/tooth-B-final.txt`, rc=1):

```
    esc_delivery_ruler_245_windows_test.go:709: ticket 245 AC#9 shape2 RED: the idle ball's roster is not the standby one (want the three always-on slots and no bare key): [ticket 245 RED: the idle ball holds cancel = "Esc" as a bare key ({Mods:16384 VK:27}, no Ctrl/Alt/Shift/Win), and RegisterHotKey takes a bare key from the WHOLE desktop - every other application has lost it for as long as the ball is idle]
        [{Name:summon ID:1 Binding:Ctrl+Alt+Q Status:live Acc:{Mods:16387 VK:81} Err:<nil> Note:} {Name:mute ID:2 Binding:Ctrl+Alt+M Status:live Acc:{Mods:16387 VK:77} Err:<nil> Note:} {Name:cancel ID:3 Binding:Esc Status:live Acc:{Mods:16384 VK:27} Err:<nil> Note:} {Name:panel ID:4 Binding:Ctrl+Alt+P Status:live Acc:{Mods:16387 VK:80} Err:<nil> Note:}]
--- FAIL: TestLive245EscDeliveryRulerFourShapes (7.03s)
```

The red names the bare key through the shared AC#6② nail (`idleRosterReds245`), not a count.

### Restoration

After each tooth the mutated file was restored by editing the exact lines back; verified per file by SHA-256 (hashes in the header above, both equal to the pre-mutation values) and `git diff --quiet -- internal/ball/ball_windows.go internal/ball/hotkey_windows.go` rc=0, and the post-restore gate re-run green (§2). The intermediates are also on record: the first tooth runs (pre-gofumpt bytes) live in the session transcript; both final teeth ran on the committed bytes so the quoted line numbers match the file.

## 4. Gates

| gate | command | result |
|---|---|---|
| test | `go test ./internal/ball -count=1 -run "245|Esc" -v` | rc=0, 8/8 PASS, `ok ... 15.579s` (§2) |
| gofmt | `gofmt -l internal/ball` | empty |
| gofumpt | `go run mvdan.cc/gofumpt@v0.12.0 -l internal/ball` | empty (after `-w` on the new file only; no other file touched) |
| vet | `go vet ./internal/ball` | rc=0 |
| D22 | `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations ...` / last line `d22scan.sh: OK - control rc=0, scan rc=0` |

## 5. Boundaries

- Write surface: `internal/ball/**` (one new `_test.go`) + this evidence file. No `PLAN.md`, `docs/specs/**`, `thresholds.go`, golden, `allowlist.txt`, `cmd/wisp/**`, frontend, design.
- No existing assertion touched; the reused judgments (`idleRosterReds245`, the AC#8 ours/foreign split, `requireQuietBallDesktop`'s hard-red-on-pollution law) are the already-ruled ones, read, not redefined.
- Skips fired today: none. Named skip paths exist in the file for: no interactive desktop (`create_window` fatal), foreground refused at injection, foreign owner at baseline / shape1 / shape4 / during the control, borrow refused at shape3 — every one names the shape and carries the numbers it measured; the instrument can never read a failure as a pass.
- Not pushed. Not done: AC#4 remains blocked on ticket 228 (untouched, as the dispatch says).
