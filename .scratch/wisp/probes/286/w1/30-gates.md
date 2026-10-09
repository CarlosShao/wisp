# 286-w1 / 30 门禁读数（票 286 AC#5）

## 门禁四数

| 门 | 命令 | rc | 输出 |
|---|---|---|---|
| 1 | `gofmt -l internal/agent/loop_golden_test.go` | `0` | **空**（无待格式化件） |
| 2 | `$(go env GOPATH)/bin/gofumpt -l internal/agent/loop_golden_test.go` | `0` | **空**；版本 `gofumpt v0.12.0 (go1.27.1)` |
| 3 | `sh scripts/d22scan.sh` | `0` | `d22scan: clean - no D22 ban violations` |
| 4 | `sh scripts/check-path-length-budget.sh` | `0` | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |

d22scan 分母原文（关键几行）：

```
d22scan.sh: scan of /d/work/workspace/projects plans/Wisp
d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]
d22scan: examined 266 production Go files under internal/ and cmd/
d22scan: scope bans #1-5 internal/      examined 228 production Go files
d22scan: scope bans #1-5 cmd/           examined  38 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  23 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 521 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined 113 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: ... ban #8 emoji coverage: ... internal/ 521 Go files, comments and _test.go included; cmd/ 113 ...
```

⇒ 我改的那一枚 `_test.go` 在 ban #8 的 `internal/` 射程内（521 枚，注释与 `_test.go` 计入）且**判 clean**：
新写的注释与 `t.Errorf` 文案零 emoji。

path-length 分母原文（关键几行）：

```
check-path-length-budget.sh: hat: rule 9 name cap=100 + issues dir prefix=21 -> relative hat=121 ; worst checkout prefix=44 -> full-path budget=165
check-path-length-budget.sh: denominator read: 8485 tracked paths
check-path-length-budget.sh: denominator: tracked paths=8485  over-budget=57  covered by roster=57  not in roster=0
check-path-length-budget.sh: bands: over the hat=57 ... roster entries=57
```

⚠ 这一发跑在证据件与票面追加**尚未 commit** 的时刻（分母 8485 未含我新建的四枚 `.md`）。
新件相对长度实测见下，全部远低于 hat=121、末段名 <100 ⇒ 结构上不可能越门；
commit 后再跑一发把带上新件的分母与 rc 补在本件文末（`40-post-commit-recheck`）。

```
.scratch/wisp/probes/286/w1/00-anchor.md      = 42 chars
.scratch/wisp/probes/286/w1/10-readings.md    = 45
.scratch/wisp/probes/286/w1/20-mutation.md    = 45
.scratch/wisp/probes/286/w1/30-gates.md       = 44
.scratch/wisp/probes/286/w1/logs/before.md    = 44   （logs/before.md, after.md, mutant-pre-hash.md,
                                                       mutant-post-hash.md, mutant-run.md, restored-run.md 同量级）
```

## `git show --stat` 名册（本腿的 commit，⛔ 不含并行别家的）

起手锚 `332e2828`。本腿名册逐笔：

```
--- 77015dbd  286-w1 起手锚
 .scratch/wisp/probes/286/w1/00-anchor.md | 99 ++++++++++++++++++++++++++++++++
 1 file changed, 99 insertions(+)

--- 9fe2f39c  286-w1 搬票 242 之后未跟搬的测试期望
 internal/agent/loop_golden_test.go | 15 +++++++++++++--
 1 file changed, 13 insertions(+), 2 deletions(-)
```

⇒ 两笔名册合起来只有 **1 枚产测文件 + 1 枚证据件**：
⛔ 无 `internal/agent/testdata/**`、⛔ 无三枚冻结件（`internal/panel/tokens_fourway_test.go`／
`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）、
⛔ 无 `frontend/**`、⛔ 无 `design/**`、⛔ 无 `docs/PLAN.md`／`docs/specs/**`／`docs/SLO.md`／
`internal/observe/thresholds.go`／`allowlist.txt`。

## 硬尺（对禁区逐条量，不靠名册目测）

```
git diff --stat 332e2828..HEAD -- internal/agent/testdata internal/agent/loop.go internal/agent/journal.go
（空输出）
PROD_DIFF_LINES=0
```

⇒ golden 语料零字节、产码 `loop.go`（含 `:369` `newTaskJournal(..., taskID, taskID)` 与 `:603` `callCorr`）
与 `journal.go`（含 `:59`／`:81`）在"锚→我最后一笔"整段区间内**零改动**。
（这一把尺同时把并行别家 commit 一并计入 ⇒ 更强：连别人都没动过它们。）

⚠ 我自造的一把假尺要具名登记：第一发越界 grep 用了 `...|golden|...` 这一支，
它匹配到了**文件名** `internal/agent/loop_golden_test.go`（本票目标件，不是 golden 语料目录）。
该支的判据无效，以 `internal/agent/testdata/` 目录尺与上面 `PROD_DIFF_LINES=0` 为准。
（同一形病＝票 141/242 那条"全称式句子必须并排穷举尺"，我这次是尺自己误命中，不是漏计。）

## 并行在飞字节（不是我的，登记以免被当成我的）

`git log --oneline 332e2828..HEAD` 里除我两笔外还有 `1e1017a5`／`22812cda`（腿 `35-a2`，票 35 只读普查）
与 `b02c41d8`（腿 `247-a4`，票 247 只读普查），三笔名册均为 `.scratch/wisp/probes/**.md`，零产码字节。
`22812cda` 那笔正文具名登记了"锚后在飞字节 `internal/agent/loop_golden_test.go`"＝本腿所写；
它只读 HEAD 对象层、显式 pathspec 未带走它 ⇒ 两腿无覆盖冲突。
