# 275-r1 — impl log (WRITE leg, shape 丁 on `.scratch/wisp/probes/161/r5/attrib.sh`)

Leg: `275-r1`. Role: WRITE leg + first-hand reader of my own before/after.
Not the researcher, not the reviewer of my own work. Every reading below must be
reproducible from a file saved under `.scratch/wisp/probes/275/r1/logs/**` or
仓外 `D:/tmp/wisp275r1/**`.

Authorization: named unfreeze `A663` (docs/reports/pending-and-issues.md:13133 and :13149).
Revoke phrase: 「撤 A663 解冻」 = `git show HEAD:.scratch/wisp/probes/161/r5/attrib.sh` restored.

Authorized change (verbatim from A663 :13149):
  文件 = `.scratch/wisp/probes/161/r5/attrib.sh` (ticket 161 tracked instrument)
  行   = `:338-344` 那一支, 只准在 `exit 2` **之前** 加 "逐枚打印 `A_OUT`"
  边界 = 不改退码语义 / 不缩分母 / 不排除 `.scratch/**` / 不动 §3 两枚台件 / 不碰 CI 步骤与 `if: !cancelled()` 守卫

---

## §0 起手锚点 (step A, captured verbatim BEFORE any edit and before any long-running command)

```
$ date
Wed Oct  7 12:28:54 CST 2026

$ git rev-parse --short HEAD
663a176d

$ git status --porcelain -- cmd internal scripts tools .github docs frontend
(empty — clean)

$ wc -l .scratch/wisp/probes/161/r5/attrib.sh
440 .scratch/wisp/probes/161/r5/attrib.sh

$ md5sum .scratch/wisp/probes/161/r5/attrib.sh
a8dd508846042ae6b7ca728f58613909 *.scratch/wisp/probes/161/r5/attrib.sh
```

Environment facts captured with the same anchor (needed to make later readings reproducible):

```
$ go env GOPATH
D:\work\base\gopath

$ D:\work\base\gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)          (ver_rc=0)

$ git config --get core.autocrlf
true

$ git ls-files '*.go' | wc -l
935
```

Ruler note for this machine: gofumpt is NOT on PATH (`which gofumpt` -> rc=1); the
instrument resolves it via `go env GOPATH`/bin/gofumpt.exe (attrib.sh:137-140), i.e.
`D:\work\base\gopath/bin/gofumpt.exe`. This matches the task's stated location.

The instrument's `:338-344` branch (read in place, not from the ticket paraphrase):
  :340 `A_RC=0`
  :341 `A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?`
  :342 `if [ "$A_RC" -gt 1 ]; then`
  :343 `echo "attrib.sh: (A) gofumpt exited $A_RC on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler" >&2`
  :344 `exit 2`
`:338` is the separate `TRACKED_GO -gt 0` guard (fires on 0 files, no A_OUT yet);
the A663 "print A_OUT before exit 2" target is the `:342-345` `A_RC -gt 1` branch.

PORCELAIN ANCHOR (final self-proof compares porcelain verbatim to this):
  `git status --porcelain -- cmd internal scripts tools .github docs frontend` = EMPTY at 663a176d.

---

## §1 未改读数 — BEFORE the fix (step C, unchanged instrument)

`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`
raw: `logs/before-stdout.txt` / `logs/before-stderr.txt`

- rc = **2**
- stdout = **1 line, 134 bytes**: `== (A) tracked tree: git ls-files -z '*.go' | xargs -0 D:\work\base\gopath/bin/gofumpt.exe -l   (tracked .go files handed to it: 935)`
- stderr = **1 line**: `attrib.sh: (A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler`
- CR bytes in stdout/stderr = 0 (LF).

Matches the task baseline (rc=2 / stdout 1 line / stderr "not readable" line; the
935-file gofumpt run exits 123 because of probes/185/c1/mut/fs_broken.go).

## §2 ★授权自审 — does printing A_OUT give a NAMED roster in the rc=2 case? (step D)

Replicated ruler_a line 341 verbatim, capturing stdout+stderr merged exactly as
`2>&1` does inside the instrument's command substitution:
`A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1)`
raw: `logs/aout-merged-before.txt` (also 仓外 `D:/tmp/wisp275r1/aout-merged.txt`)

- **A_RC = 123** (xargs -> 123 because gofumpt exits non-zero on the parse error)
- **merged A_OUT = 33 lines** = **1 diagnostic row**
  (`.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: imports must appear before other declarations`)
  + **32 named bare `.go` file rows** (distinct named files via norm = 32).

★ Verdict: gofumpt -l DOES emit a per-file roster on stdout even while it hits an
unparseable file and exits non-zero. Printing A_OUT in the exit-2 branch therefore
produces a NAMED roster (the broken file is named by its own diagnostic row, and the
32 tracked-dirty files — including the two planted controls `probes/241/v1/posctl/
badly_formatted.go` and `probes/259/r1/gofumpt-negctl/probe.go` — are named).
**D did NOT stop me; the authorized shape 丁 is effective; I proceeded to code.**
No widening was needed.

## §3 落地 — the edit (shape 丁)

Confined to the `A_RC -gt 1` branch of `ruler_a`, inserted between the existing
stderr `echo` (:343) and `exit 2` (:344): a `while read` loop that prints each
non-empty `A_OUT` row to stdout under a labelled header (`A-ROSTER\t<row>`),
then still `exit 2`. Full diff lives in commit `febca8f9`.

- `sh -n` syntax check: rc=0
- wc -l: **440 -> 459** (+19); md5: **a8dd508846042ae6b7ca728f58613909 -> 77ee11a8d0cc8dd8248c6a9b8e164a78**

A663 boundary held, each one re-checked:
- exit code stays **2** (no change to `exit 2`), guard `A_RC -gt 1` unchanged.
- denominator unchanged: still every tracked `*.go` (`git ls-files -z '*.go'`), 935 handed to gofumpt (banner confirms).
- `.scratch/**` NOT excluded; not shrunk.
- the two planted controls (241 posctl / 259 negctl) NOT touched (only READ by gofumpt; they appear as roster names, never edited).
- `.github/workflows/ci.yml` and `if: ${{ !cancelled() }}` NOT touched (porcelain group `.github` clean).

## §4 改后读数 — AFTER the fix (step E, SAME command, broken file still present)

`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`
raw: `logs/after-stdout.txt` / `logs/after-stderr.txt`

- rc = **2 (unchanged)**
- stdout = **35 lines, 2276 bytes** = 1 banner + 1 roster header + **33 `A-ROSTER` rows**
  (32 named `.go` files + 1 `fs_broken.go:4:1` diagnostic) — i.e. the red now names files.
- stderr = **1 line, 137 bytes (unchanged)**: the same "not readable by this ruler" line.
- CR bytes in stdout = 0 (LF).

## §5 未把门做瞎 — rc=1 shape roster unchanged (gate-not-blinded proof)

rc=1 shape produced WITHOUT touching the shared tree or shrinking the committed
denominator: a **temp index 仓外** (`D:/tmp/wisp275r1/tmpidx`) that drops ONLY
`fs_broken.go` (the deliberately-broken sample -> gofumpt then reads all 934 ->
`A_RC<=1` -> normal roster path). Disk bytes untouched; `git rm --cached` was run
ONLY against the temp index, never the shared one.
raw: `logs/before-rc1-stdout.txt` + `logs/after-rc1-stdout.txt` (stderr captured too).

- **before edit**: rc=1, `(A) tracked-tree  lines=32 files=32`, 32 `TRACKED-DIRTY` rows.
- **after edit**: rc=1, `lines=32 files=32`, 32 `TRACKED-DIRTY` rows.
- `diff before-rc1 after-rc1` on **stdout AND stderr: IDENTICAL**.
- `A-ROSTER` count in the rc=1 shape = **0** (my block only runs when `A_RC>1`,
  which never happens in the rc=1 shape) -> the normal red path is byte-for-byte
  what it was.

EOL context (state it so the number is not read as 19): this runs against the
**current working tree**, where the 13 A663-listed files are stale **CRLF** on disk
(`core.autocrlf=true` while `*.go text eol=lf`), so gofumpt flags them ->
**lines=32 files=32** (the "working-tree shape" the dispatch names). The CI-equivalent
**LF** shape (fresh `git archive`/checkout at eol=lf) prints **19**; I did **not**
reproduce 19 here — I deliberately used the working-tree-CRLF context, which is why
the count is 32 not 19. (CRLF judged with `tr -cd '\r' | wc -c`, not `grep -c $'\r'`.)

## §6 门禁 — gates (all run; positive controls shown)

- `sh scripts/d22scan.sh` -> **rc=0**. Its **positive control runs FIRST** and passes
  (`runtests.sh -C tools/d22scan ./...`, forced `-count=1`):
  `TestScanDetectsAllSeededViolations`, `TestScanCleanRepoIsGreen`, `TestAllowlistSuppressesOnlyListedPaths`,
  `TestCheckRootRejectsBlindRoots` (+4 subcases), `TestScanAloneIsNotAFalsifier`, `TestCheckRootAcceptsRealRepo`,
  `TestScannerSelfScanOfRealRepoIsGreen`, `TestEmojiBanCoversGoSourcesNotJustDesign`, `TestBan8MathBandAndRemainingGaps` all PASS;
  then `d22scan: clean - no D22 ban violations`. (d22scan is its OWN module
  `tools/d22scan`, not the root module — no root `go build`/`go test` was run.)
  raw: `logs/d22scan.txt` / `logs/d22scan.err`.
- Ticket 161 self-check carrier for THIS script = `attrib.sh --self-test` (no separate
  wrapper exists under scripts/ or .scratch/wisp/probes/161/). On my **edited** file ->
  **rc=0 SELF-TEST GREEN**, `cases=8 failures=0` (5 RING + 3 SILENT). raw: `logs/selftest-green.txt`.
- Carrier's positive control still reddens: copied my edited attrib.sh to 仓外
  (`D:/tmp/wisp275r1/attrib-mut.sh`), blinded the classifier's UNATTRIBUTABLE branch
  (`CL_RC=1` -> `0`, the 4-space line; the shared file was NOT edited), ran with
  `ATTRIB_ROOT=<repo>` -> **rc=1 SELF-TEST RED, failures=4** (the 4 unattributable
  RING cases, while tracked + silent cases stay correct). raw: `logs/selftest-RED-proof.txt`.
  Shared attrib.sh md5 re-checked after this = 77ee11a8 (untouched by the 仓外 proof).

NOTE on filenames: the committed `.gitignore:8` has a repo-wide `*.out` rule, so my
first evidence commit silently skipped the three `.out` gate logs. They were renamed
to `.txt` (not ignored) and re-committed; the raw gate output is now under
`logs/d22scan.txt`, `logs/selftest-green.txt`, `logs/selftest-RED-proof.txt`.

## §7 最终自证 — final self-proof

- date (end) = `Wed Oct  7 12:38:03 CST 2026`; HEAD = `febca8f9`
- `git status --porcelain -- cmd internal scripts tools .github docs frontend` = **EMPTY**
  -> **verbatim equal to the §0 anchor** (I never touched those groups; my work is all under `.scratch/**`).
- my commits (`git log --oneline`), pathspec-limited to my own files:
  - `4dc6cab6` 275-r1: land §0 pre-work anchor — `1 file changed, 67 insertions(+)` (`.scratch/wisp/probes/275/r1/impl.md`)
  - `febca8f9` 275-r1: shape 丁 ... rc=2 exit now prints the per-file roster — `1 file changed, 19 insertions(+)` (`.scratch/wisp/probes/161/r5/attrib.sh`)
- attrib.sh md5: before `a8dd508846042ae6b7ca728f58613909` -> after `77ee11a8d0cc8dd8248c6a9b8e164a78`; wc `440 -> 459`.
- (shared-tree note: leg `274-v1` landed `5da17c58` on the branch between my §0 anchor `663a176d`
  and my first commit; it is not mine and my pathspecs never included its file.)

## §8 未做完的格 — cells I did NOT measure (naming what was not done)

- **AC#3 (零误伤)**: not done. Needs a fresh syntactically-valid-but-unformatted
  tracked sample (仓外 / `-overlay`) proving a real regression is still SEEN after
  丁. That is a non-implementer / acceptance cell for AFTER landing; not mine to fake.
- **AC#4 (恒真性进攻 on the 丁 print)**: not run as a dedicated attack. I proved the
  ticket-161 carrier self-test still reddens, but I did NOT plant a mutation whose
  only effect is "stop the roster printing" and show a named criterion goes red.
  Note for the acceptance leg: the 丁 roster print has NO assertion yet — it prints;
  whether it prints is not itself gated.
- **CI-side "the guarded step really ran and its color"**: needs a push (owner said
  no push) -> [待推送取数]. I did NOT substitute a local reading for a CI reading.
- **The CI-equivalent LF shape (19)**: not reproduced (context = working-tree CRLF
  -> 32). Stated in §5.
- **gofumpt version**: pinned here at `v0.12.0 (go1.27.1)`; the roster/exit-code
  behaviour I measured is that version's.
- I did NOT run `go test`/`go build` on the root module (another leg holds that window).
