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
