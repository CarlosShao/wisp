# 444 — seven race-window instruments that could turn RED on abstention

**Ticket**: `.scratch/wisp/issues/444-seven-more-race-window-instruments-can-red-on-abstention.md`
**Shape**: ticket 441's (`docs/evidence/s1/441-instrument-abstention.md`) — *"the window was REACHED"* may
abstain, *"the guardrail is still alive"* must stay true-RED, and the abstention must be a named line with
numbers, never a silent green.
**Base HEAD when this landed**: `6949484e1`. **Platforms**: keep-in-mind matrix — 4 of the 7 instruments are
`!windows`/`unix` (they run in `golang:1.27-bookworm`, image cached, module cache in the named volume
`wisp444gomod`), 3 are `windows`-tagged (they run natively, go1.27.1 windows/amd64).
**Write surface**: `internal/winsec/**` `*_test.go` only, this document, `.scratch/wisp/probes/444/**`, and an
append-only Progress bullet in the ticket. **Production bytes changed: 0** (§5).

---

## 0. What was actually wrong, in one paragraph

Ticket 441 established the distinction and fixed one instrument. The seven here are the same shape: each one
counted **two different things with one `if`** — *the racer's arrival* (did the window exist at all) and *the
route's landing* (did a `PrivateFile` call stand inside that window). Under load either count can legitimately
be zero: the route is a paced loop (`2ms` sleep per attempt) and the window is a few microseconds wide, so a
starved scheduler can make `landed` zero while the racer proves hundreds of flips, and a pinned handle can make
the racer's own `rename` fail outright (`Access is denied.`) so *neither* side arrives. Both zeros used to print
**RED**, i.e. the instrument cried wolf exactly when it had the least information. The fix, per 441: keep every
zero that the guardrail itself owns as a RED, turn the zero that only says *"the window was not sampled"* into a
greppable `ABSTAIN(444:<shape>:<why>)` line **carrying the numbers of that run**, and add an **un-raced control**
that reaches the window deterministically (one pin, one flip, one act) and asserts the guardrail there. The
control is what carries the row; the raced sample is evidence, not the threshold.

---

## 1. AC#1 — the red each instrument printed when the window was absent, read verbatim

Method: no assertion was touched to produce these. Each was read at a **shrunk denominator** (`WISP_444_ATTEMPTS`,
a knob of the run — `internal/winsec/attempts_444_test.go`, `n >= 0`, default = the caller's own constant) or, for
the seam, at `WISP_444_HOP_REFUSE=1` which parks a **non-empty** directory where the hop lands so `rename(2)` refuses
every hop — the same branch the platform takes on CI when the pinned handle turns the rename into *"being used by
another process"*. `GOMAXPROCS=1`.

| # | instrument (file:arm) | run that showed it | the sentence, verbatim |
|---|---|---|---|
| a | `openleaf_322_other_test.go:208` — `liveness == 0 \|\| flips == 0` | `WISP_444_ATTEMPTS=1`, 3/3 RED (log `05-base-unix-att1.log`) | `AC#1 RED instrument: across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 0 flips (0 errors, first ""), so the damaged count above is not evidence about a window that was never reached` |
| b | `ancestoropen_342_unix_test.go:96` — `!landed` → "not a landing reading" | `=1`, 3/3 RED (log 05) | `AC#1 RED (not a landing reading): pinnedParentBeforeOpen answered 0 attempts, and that numerator is the landing's own — the pre-landing route cannot raise it and the landing's liveness gate requires it — so this row cannot certify the landing's closed window either way. On pre-landing bytes these numbers are the pre-fix reading (342-v1's VPROD/VBASE measured deleted 65..86/1500 there; leg 342-r3 re-took deleted=0 of 1 on the by-name route), and reading any zero among them as 'the window was already closed' is the exact mistake ticket 342 AC#1 forbids (provedLinkCycles=0, attempts=1, flipErrs=0, first "")` |
| b′ | `ancestoropen_342_unix_test.go:116` — `provedLinkCycles == 0 \|\| pinnedParentLeg == 0` | same run, the arm the ticket names | `AC#1 RED instrument: the racer proved 0 link-state cycles and the pinned walk refused 0 attempts across 1 (flipErrs=0, first ""), so the zero above is not evidence about an open that was ever reachable` |
| c | `ancestorunlink_348_unix_test.go:154` — `provedLinkCycles == 0 \|\| liveness == 0` | `=1`, 3/3 RED (log 05) | `AC#2 RED instrument: the racer proved 0 link-state cycles and the pinned walk refused 0 attempts across 1 (flipErrs=0, first ""), so the zero above is not evidence about a removal that was ever reachable` |
| d | `ancestoropen_342_windows_test.go:57` — `provedReparse == 0 \|\| formRefusals == 0` | `=1`, 6/6 RED, both forms (log `01-base-342win-att1.log`) | `AC#1 RED instrument: reparseAt proved 0 flipped cycles and the landing's own legs answered 0 attempts across 1 (flipErrs=1, first "rename C:\\…\\data\\sub C:\\…\\data\\parked-real-sub: Access is denied.") - so the zeroes above are not evidence about a window that was ever reached` |
| e | `ancestorunlink_348_windows_test.go:112` — `provedReparse == 0 \|\| formRefusals == 0` | `=1`, 5/5 RED (log `03-base-348win-att1.log`) | `AC#2 RED instrument: reparseAt proved 1 flipped cycles and the landing's own legs answered 0 attempts across 1 (flipErrs=0, first "") - so the zero above is not evidence about a removal that was ever reachable` |
| f | `deleteseam_348r2_windows_test.go:131` — `landed == 0 && refused == 0` | HEAD source string (`git show HEAD:…` lines 130-131). This one has no run of its own at HEAD because HEAD has no denominator knob: the only way to reach `calls == 0` on HEAD bytes is a route that never calls the seam (mutation `g`), and under that mutation the *wiring* red at `:126` fires first with different text. After the fix the same state is `WISP_444_ATTEMPTS=0` and it prints the abstain line (log `10-postfix-348r2-att0.log`). | `AC#2 RED instrument: 0 seam calls produced neither a landed nor a refused flip, so the flip was never even asked and the deletion count below is unaccompanied by any statement about the ancestor` |
| g | `deleteseam_348r2_windows_test.go:159` — `landed != calls` | `WISP_444_ATTEMPTS=4 WISP_444_HOP_REFUSE=1`, 2/2 RED (log `04-base-348r2-hoprefuse.log`) | `AC#2 RED instrument: the ancestor flip stood at the removal instant on only 0 of 4 seam calls (refused=4, first refusal verbatim: "rename C:\\…\\data\\sub C:\\…\\data\\parked-real-sub: Access is denied."), so the deletion count is a reading about a world that was not entered rather than about the removal` |

Third finding, and it is the one the ticket did not list: **b′ and b are two different arms, and the one that
actually fires first is `:96`, not the `:116` the ticket names.** At `attempts=1` the `!landed` branch printed its
"not a landing reading" sentence and **`return`ed**, which skipped the damage assertions below it — so on the exact
run where the instrument had least to say it also said *less than nothing about damage*. That `return` is gone (§4).

---

## 2. AC#2 — the shape of the fix, and the ruler you can grep

### 2.1 The split, applied seven times

```
if <racer-side arrival == 0>            -> RED, unchanged bytes (this is the teeth that a deleted racer must hit)
if <route-side landing == 0>            -> ABSTAIN(name, numbers)   [new]
if <un-raced control violated>          -> RED (new, deterministic)
if <damage/deletion/narrowing > 0>      -> RED, unchanged bytes
```

"Racer-side" = the count the racer's own loop owns (`flips`, `provedLinkCycles`, `provedReparse`). If *that* is
zero the fixture machinery is gone, and no future run would ever be able to test anything — that must be red.
"Route-side" = the count that depends on the route being *scheduled inside* the racer's window
(`leafLegRefusals`, `pinnedParentLeg`, `formRefusals`, `seamCalls`). That one is a property of the scheduler, not of
the guardrail, so its zero abstains.

### 2.2 The greppable token

Every abstention prints exactly one line, first token `ABSTAIN(444:`:

| shape | token |
|---|---|
| 322 leaf leg | `ABSTAIN(444:322-leaf:leg-starved)` |
| 342 unix ancestor open | `ABSTAIN(444:342-unix-ancestor:route-missed-link)` |
| 348 unix ancestor removal | `ABSTAIN(444:348-unix-removal:route-missed-link)` |
| 342 windows ancestor open | `ABSTAIN(444:342-win-ancestor:route-missed-link)` |
| 348 windows ancestor removal | `ABSTAIN(444:348-win-removal:route-missed-link)` |
| 348-r2 seam, never reached | `ABSTAIN(444:348r2-seam:removal-never-reached)` |
| 348-r2 seam, hop refused | `ABSTAIN(444:348r2-seam-hop:flip-not-asked)` |

`grep -c 'ABSTAIN(444:' <log>` counts abstentions; `grep -c 'RED instrument' <log>` counts arrivals; a green run
with **0** of both is the silent-green failure this ticket exists to prevent, which is why the control's own
`RED instrument` sentence is *always* present in the source (it just does not fire). Each abstain line ends with the
same clause: *"the threshold that carries <the row> is the un-raced control above (…), not a green from absence"* —
so a reader who lands on the line cannot mistake it for a pass.

### 2.3 The un-raced control — one per instrument, and it is deterministic

Each control does what the racer does, but by hand and in a fixed order: **pin the parent → flip the ancestor →
perform the one act through the pin**, then read back *the foreign witness* and *the caller's own leaf*. Not racy ⇒
no scheduling dependence ⇒ this is the thing a mutation must break.

| # | control | what it asserts |
|---|---|---|
| a | `leafLegControl444` (`openleaf_322_other_test.go`) | plant a sibling symlink at the leaf name, call `openPinnedLeaf(pin, base, O_WRONLY\|O_CREATE\|O_TRUNC\|O_NOFOLLOW, 0o600)`; require `legLeafOpen322r1` refusal, and the witness's **body and mode unchanged** |
| b | `ancWindowControl444(t,"open")` (`ancestorworld_342_348_unix_test.go`) | `pinnedParent(name,name)`, `rename(sub,park)`, `symlink(outside,sub)`, prove `Lstat` is a link, `openPinnedLeaf(…\|oNoFollowPlatform)`; require `ownLeafAtPin == payload` (the act landed *under the pin*), `foreignBody`/`foreignMode` intact, `foreignThere` |
| c | `ancWindowControl444(t,"removal")` | same stage, then `removeRefusedLeaf(pin, base, pinned)`; require the foreign witness **still present and intact**, `actErr == ""`, and `ownLeafGone` (our own leaf, under the pin, is the one removed) |
| d | `winAncWindowControl444(t, tc.form, "open")` (windows) | `newWinAncWorld342r1(t, form)` **with the subtest's own form** — a junction subtest must not be carried by the symlink control and vice versa — then `w.toFlip()` (which proves `reparseAt`), then `openPinnedLeaf`; require `ownLeafAtPin == winAncPayload342r1`, witness intact |
| e | `winAncWindowControl444(t, symDir, "removal")` | same, `removeRefusedLeaf`; require witness intact + `ownLeafGone` |
| f/g | the same removal control, called **once at the top of the seam test, before `withInjectedSealFailure`** | it must run *outside* the injected failure, otherwise it can never reach the removal at all |

Honest limit, stated once and not papered over: `preOpenLeafGuard` is a plain function on POSIX (not a `var`), so
the control drives `openPinnedLeaf`/`removeRefusedLeaf` **directly, the way the route calls them** — it does not
re-drive `resolveString`/`privateFile`'s ordering. That ordering is pinned statically by
`routeguard_369_test.go` (ticket 369's precedent). The control therefore carries *the act's* identity, not *the
route's*; the raced reading carries the route.

---

## 3. AC#3 — one "guardrail removed" mutation per instrument, on a tree outside the repo

Tree: `D:\tmp\wisp\444\export` — `git clone --shared` of the repo at `6949484e1`, then the working-tree
`internal/winsec/*.go` copied in (identical for production files; the test files are this leg's). Mutations are
applied to **production** files in that tree only, one at a time. Restore = the same bytes back; the proof is the
SHA256 string comparing equal to the pristine value **and** the test green again.

| # | file | mutation | SHA256 pristine → mutant | the red it printed, verbatim (truncated to the leg it names) | restore |
|---|---|---|---|---|---|
| a | `winsec_other.go` | `const oNoFollowPlatform = syscall.O_NOFOLLOW` → `= 0` | `19F3C211…D3E2` → `9BA3F07B…4FA8` | `AC#1 RED instrument: the un-raced leaf-leg control did not get the refusal ticket 322 pays for. attempts=1 refused=false leg="" witness-body="" witness-mode=-rw-rw-rw- verbatim: <nil>` (and `AC#1 RED: 325 of 3000 foreign witnesses lost their own bytes or mtime … content is "", want its own bytes`) | `19F3C211…D3E2`, green |
| b | `pinnedleaf_342_unix.go` | `unix.Openat(int(p.fd()), base, …)` → `unix.Open(path, …)` (by-name) | `22D5D118…441A` → `5876B3F4…7A52` | `AC#1 RED instrument: the un-raced window control - the route's own pin, then the ancestor flip, then the act - did not keep this route's act inside the pinned object: … foreign-body="raced private payload 342r1 through a swapped ancestor" …` | `22D5D118…441A`, green |
| c | `pinnedleaf_342_unix.go` | `unix.Unlinkat(int(p.fd()), base, 0)` → `os.Remove(filepath.Join(p.name, base))` | `22D5D118…441A` → `264AA6CC…E31D` | `AC#2 RED instrument: the un-raced removal control … foreign-present=false …` **and** `AC#2 RED: 70 of 1500 attempts found the foreign witness GONE, so the refusal-path unlink in privateFile re-resolved the spelling through an ancestor swapped into a link and removed a name in somebody else's tree` | `22D5D118…441A`, green |
| d | `pinnedleaf_342_windows.go` | `p.root.OpenFile(base, flags, perm)` → `os.OpenFile(filepath.Join(p.dir, base), …)` | `1C15A7DD…17D8` → `07CC696D…07C5` | `AC#1 RED instrument: the un-raced window control … act-error="winsec: path is not provably resolved, refusing to seal: … refusing to open for a private write …"` — both forms (junction and symlink-dir) | `1C15A7DD…17D8`, green |
| e | `pinnedleaf_342_windows.go` | `p.root.Remove(base)` → `os.Remove(filepath.Join(p.dir, base))` | `1C15A7DD…17D8` → `C34BABE9…4766` | `AC#2 RED instrument: the un-raced removal control … foreign-present=false …` | `1C15A7DD…17D8`, green |
| f | same as (e), observed at the seam | — | — | `AC#2 RED: 400 of 400 attempts found the foreign witness GONE through the refusal-path removal` + the control's `RED instrument … neither kept the foreign name nor removed the pinned one` | green |
| g | `winsec.go` | the refusal-path removal call deleted from the seal-failure branch | `F0BBF9FF…55A1` → `8B7B3D50…167D` | `AC#2 RED instrument: the removal seam was reached 0 times in 400 attempts, so this reading is not about the refusal-path removal at all (an attempt that never got there cannot say anything about what that removal would have done)` — and note the `ABSTAIN(444:348r2-seam:removal-never-reached)` line printed **alongside it**: the abstention does not swallow the red | `F0BBF9FF…55A1`, green |

Row (g) is the load-bearing one for `:131`: deleting the removal makes `seamCalls == 0`, which is exactly the state
that used to print nothing but a false red. After the fix it prints the **wiring** red (row 126, unchanged bytes) and
*additionally* the abstain line — the reader gets both facts instead of a red that blamed the flip.

---

## 4. AC#4 — before/after assertion table (nothing deleted, nothing loosened)

"Kept" means the bytes are unchanged apart from the denominators now reading the run's own `attempts` variable
instead of the compile-time constant (so a shrunk run cannot print a number about a run it did not do).

| # | before | after |
|---|---|---|
| a | `if liveness == 0 \|\| flips == 0 { Errorf(...) }` | `flips == 0` → same Errorf (**kept**); `liveness == 0 && flips > 0` → `Logf ABSTAIN` (**split**); **new** `if !ctl.refused \|\| !ctl.witnessIntact { Errorf }`. `damaged > 0`, `narrowed > 0`, `landed > 0` unchanged. |
| b | `if !landed { Errorf("not a landing reading"); return }` | `Logf` the same sentence as a **reading** (attribution, not a verdict) and **no `return`** — so the damage arms below now run on exactly the runs where they used to be skipped (**strictly more gated**); `truncated/deleted > 0` unchanged; `provedLinkCycles == 0` **kept** red; `pinnedParentLeg == 0 && provedLinkCycles > 0` → `Logf ABSTAIN`; **new** control red. |
| c | `if provedLinkCycles == 0 \|\| liveness == 0 { Errorf }` | `provedLinkCycles == 0` **kept**; `liveness == 0 && provedLinkCycles > 0` → `Logf ABSTAIN`; `deleted > 0`, `narrowed > 0` unchanged; **new** control red. |
| d | `if provedReparse == 0 \|\| formRefusals == 0 { Errorf }` | `provedReparse == 0` **kept**; `formRefusals == 0 && provedReparse > 0` → `Logf ABSTAIN`; `truncated > 0` unchanged; **new** control red (per form). |
| e | same as (d) on `deleted` | same split; `deleted > 0` unchanged; **new** control red. |
| f | `if landed == 0 && refused == 0 { Errorf }` | `calls > 0` → same Errorf (**kept**); `calls == 0` → `Logf ABSTAIN` (`landed==0 && refused==0` ⟺ `calls==0`, since every seam call increments exactly one of the two) |
| g | `if landed != calls { Errorf }` | control satisfied → `Logf ABSTAIN`; control **not** satisfied → the same Errorf byte-for-byte (**kept**) |
| f/g | `calls != attempts` wiring red; `deleted > 0` red | **unchanged** (only the printed denominators follow `attempts`) |

Count: **0 assertions deleted, 0 relaxed.** 6 red arms kept verbatim, 6 split into (kept red + new abstain),
7 new deterministic control reds added, and 1 early `return` removed, which re-arms 2 damage assertions on the
starvation path.

---

## 5. AC#5 — gates

| gate | reading |
|---|---|
| `gofmt -l internal/winsec` | empty |
| `go vet ./internal/winsec` (windows/amd64, go1.27.1) | `rc=0` |
| `go vet ./internal/winsec` (linux/amd64, golang:1.27-bookworm) | empty output (rc=0) |
| `bash -n .scratch/wisp/probes/444/repro-444.sh` | `BASH-N-OK` |
| `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations` (482 production Go files examined) |
| production bytes changed | `0` — `git status --porcelain internal/winsec` lists only `*_test.go` (+ the new `attempts_444_test.go`) |
| `sh scripts/check-path-length-budget.sh` | **VERDICT RED — not mine.** The single out-of-roster path is `.scratch/wisp/issues/274-…-done.md` (tracked, added by an earlier leg; `git log -1 --` that path = `c0e48dbc4`). My own new paths measure 37 / 36 / 45 / 56 chars relative, all far under the 122 middle band. Reported rather than edited: that file is not on this leg's write surface. |

---

## 6. Boundary self-proof — what this leg proves, and what it does not

Proved here: (i) each of the seven zeros that used to red now has a named abstention carrying that run's numbers;
(ii) each instrument gained a deterministic un-raced control that reds when its guardrail is removed — verified by
seven mutations on an out-of-repo tree, each restored to the same SHA256 and green; (iii) no assertion was deleted
or relaxed, and one early `return` that skipped damage checks was removed.

**A residual I am not allowed to hide, with its measurements.** The *racer-side* zero stays RED by design (§2.1),
and a bench that shrinks the denominator far enough starves the racer itself, so a shrunk run is **not** a clean
demonstration of the abstention. Measured, so the reader can judge:

- `WISP_444_ATTEMPTS=20`, container, `-cpu=1`, `count=3` (log `repro-all-20-r0.log`): **5 abstain lines, 6 case-runs
  pass, 3 red — and the reds all name the racer** (`AC#1 RED instrument: the racer proved 0 link-state cycles across
  20 attempts (flipErrs=1, first "remove …")`). In the *same* log, 342-unix prints
  `ABSTAIN(444:342-unix-ancestor:route-missed-link) attempts=20 provedLinkCycles=2 pinnedParentBeforeOpen=0 …` and
  348-unix prints `ABSTAIN(444:348-unix-removal:route-missed-link) provedLinkCycles=1 pinnedWalkRefusals=0 …`. That
  side-by-side is the whole ticket in one screen: **racer zero → red that names the racer; route zero → named
  abstention.** Nothing was smoothed to make it look tidier.
- `WISP_444_ATTEMPTS=1`, `GOMAXPROCS=1`, `count=2` (log `repro-all-1-att1.log`): **6 of 6 case-runs red on the racer
  arm**, because at one attempt the route finishes in under a millisecond and the racer's goroutine is not even
  scheduled — `flipErrs=0, first ""`, `provedLinkCycles=0`. That is a bench artefact and the instrument says so by
  naming the racer; it is not the state CI is in.
- At the **committed** denominators the instruments were green in every run made for this leg: 400-attempt windows
  runs 4/4 (and at `attempts=1` on that platform 2 case-runs abstained while 4 red on the racer arm — log
  `09-postfix-win-att1.log`), 1500-attempt unix runs 5/5, 3000-attempt 322 runs 5/5.
- Under controlled contention (4 CPU burners + `GOMAXPROCS=1`, committed denominators, `count=3`, log
  `18-load-committed-unix.log`) all nine case-runs were green and the route-side counts were 322
  `leafLegRefusals=0,3,5` per 3000, 342u `pinnedParentBeforeOpen=13,11,11` per 1500, 348u `pinnedWalkRefusals=23,24,22`
  per 1500. **The `0` is the one that used to be a false red**: that run printed
  `ABSTAIN(444:322-leaf:leg-starved) attempts=3000 flips=130 flipErrs=0 leaf-leg-refusals=0 damaged=0 …` and passed.
  So the abstention is not cosmetic — the route-side count is 0.7–1% of attempts and reaches 0 on a loaded bench.

**Second residual.** The two Windows ancestor instruments depend on the bench having directory-symlink or junction
privilege; the junction subtest prints `Access is denied.` for its rename when the pinned handle is held, which is
why the control (not the racer) carries those rows. Neither instrument is skipped on a bench without the privilege:
they abstain (route side) or red (racer side) — never silently pass.

---

## 7. Re-run it yourself

```sh
# POSIX instruments (the three !windows / unix ones) - container, cached image:
docker run --rm -v "<repo>:/src" -v wisp444gomod:/go/pkg/mod -w /src \
    -e WISP_444_ATTEMPTS=1 -e COUNT=2 golang:1.27-bookworm \
    bash -c 'sh .scratch/wisp/probes/444/repro-444.sh'
# -> abstain-lines / controls-printed / red-arrival / red-guardrail / cases-pass, one line each

# Windows instruments (the three windows-tagged ones) - host, GOMAXPROCS=1:
go test ./internal/winsec -count=1 -v \
  -run '^(TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1|TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1|TestAC2SeamWindowsTheRefusalPathRemovalCannotBeMovedByAnAncestorFlip348r2)$'
# abstain the seam:   WISP_444_ATTEMPTS=0            -> 444:348r2-seam:removal-never-reached
# abstain the hop:    WISP_444_ATTEMPTS=4 WISP_444_HOP_REFUSE=1 -> 444:348r2-seam-hop:flip-not-asked
# abstain the routes: WISP_444_ATTEMPTS=1            -> 444:342-win-ancestor / 444:348-win-removal / 444:322-leaf
```

Logs for every reading above: `.scratch/wisp/probes/444/logs/` (`01`…`18`, base before the fix, post-fix,
mutations, restore-green, and the contention run).
