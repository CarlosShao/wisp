# 275-v1 — 非实现者对抗验收（形丁落在 `.scratch/wisp/probes/161/r5/attrib.sh`）

Leg: `275-v1`. Role: adversarial acceptance (NON-implementer). Code under review = leg `275-r1`
(commits `4dc6cab6` → `febca8f9` → `d10697eb` → `31261b4d`, orchestrator's own landing commit `79fbcd57`).

Write surface: `.scratch/wisp/probes/275/v1/**` only. Scratch/mutations OUTSIDE the repo under
`D:/tmp/wisp275v1/`. Evidence extension = `.txt` (⛔ not `.out`, `.gitignore:8` is a repo-wide `*.out`).
I did not edit `attrib.sh`, did not flip any ticket box, did not touch `.github/**`, `scripts/**`,
`cmd/**`, `internal/**`, or the two other tickets' stage pieces.

---

## §0 起手锚点（在任何长跑命令之前落盘；本节即第一笔 commit 的内容）

Raw file: `anchor-raw.txt` (same bytes copied to `D:/tmp/wisp275v1/anchor-raw.txt`).

```
$ date                                  Wed Oct  7 12:45:32 CST 2026
$ git rev-parse --short HEAD            79fbcd57
$ md5sum …/161/r5/attrib.sh             77ee11a8d0cc8dd8248c6a9b8e164a78
$ wc -l …/161/r5/attrib.sh              459
$ sh -n …/161/r5/attrib.sh              sh_n_rc=0
$ git status --porcelain                753 lines  (full list -> D:/tmp/wisp275v1/porcelain-anchor.txt)
$ git status --porcelain -- cmd internal scripts tools .github docs frontend
                                        0 lines    (PORCELAIN ANCHOR = EMPTY for these groups)
$ git ls-files '*.go' | wc -l           935
$ git diff --numstat febca8f9^ febca8f9 -- …/161/r5/attrib.sh
                                        19  0  .scratch/wisp/probes/161/r5/attrib.sh   (19 added / 0 deleted)
```

Restore proof compares verbatim to: (a) `md5sum attrib.sh` = `77ee11a8d0cc8dd8248c6a9b8e164a78`,
(b) `git status --porcelain -- cmd internal scripts tools .github docs frontend` = **0 lines**,
(c) the 753-line full-repo porcelain file.

Authority read on disk (my own reading, not the ticket paraphrase):
- `attrib.sh:341` `A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?`
- `attrib.sh:342` `if [ "$A_RC" -gt 1 ]; then`
- `attrib.sh:343` stderr `echo` "…the tracked tree is not readable by this ruler"
- `attrib.sh:344-353` the added comment block (A663 / shape 丁)
- `attrib.sh:354` `if [ -n "$A_OUT" ]; then`
- `attrib.sh:355` banner `== (A) per-file roster gofumpt emitted before it hit an unparseable file …`
- `attrib.sh:356-361` `while IFS= read -r roster_line` + `printf '   A-ROSTER\t%s\n' "$roster_line"`
- `attrib.sh:363` `exit 2`

The landing is exactly the shape `A663` :13149 authorized (only print `A_OUT` before `exit 2`;
no exit-code, guard, or denominator change).

Note on the ticket vs disk: ticket §2 row 5 quoted `attrib.sh:173-178`; the ticket's own
correction block (`:31-36`) says the real branch is `:338-344`. On disk (459-line version) the
guard is at `:342` and `exit 2` now sits at `:363` after the insertion. **Disk wins** — I cite
line numbers from the current 459-line file.
