# 票 289 / 腿 289-v1 — AC#4 判语：门禁逐把自跑（每把自己落 rc）＋越界名册

**判语：成立。**（四把门禁全过；产码笔名册只含那 4 枚件；三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／`testdata/**` 相对 HEAD 零字节；⚠ 树**不纯净**这一条具名写明，没有据它判绿。）

每把尺的 rc 都由**被量命令自己的 `$?`** 在下一条语句取（⛔ 不跨管道取，避免取到管道尾的码）。原始输出＝本目录 `raw-ac4-gates.md`（§1–§9 逐节，每节一行 `rc=`）与 `raw-d22scan.md`／`raw-d22scan-badroute.md`。

| 尺 | 命令逐字 | 读数 | rc |
|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | 空输出 | **rc=0** |
| gofmt | `gofmt -l internal/panel/subagent_roster_197.go internal/panel/subagent_roster_197_test.go internal/panel/subagent_blocked_220_test.go cmd/wisp/subagent_carrier_197_test.go` | 空（4 枚件无一被列） | **rc=0** |
| gofumpt（裸 PATH） | `gofumpt -l internal/panel/subagent_roster_197.go` | `/usr/bin/bash: line 1: gofumpt: command not found` | **rc=127**（⚠ 具名不掩盖：派单说的"不在 PATH"本腿现跑复现） |
| gofumpt（绝对路径） | `"$(go env GOPATH)/bin/gofumpt.exe" -l <同那 4 枚件>` | 空 | **rc=0** |
| d22scan | `sh scripts/d22scan.sh` | 末行逐字 `d22scan: clean - no D22 ban violations`；ban #8 射程逐字 `internal/ 524 Go files, comments and _test.go included`＋`cmd/ 116 Go files, comments and _test.go included` ⇒ 本票那 4 句注释**在射程内且未触雷** | **rc=0** |
| d22scan 的反面（派单的警告） | `GOFLAGS= go run ./tools/d22scan` | 逐字 `main module (github.com/CarlosShao/wisp) does not contain package github.com/CarlosShao/wisp/tools/d22scan` | **rc=1**（⛔ 这一发本来就该失败，`tools/d22scan` 是独立模块；正解是上面那一发 `sh scripts/d22scan.sh`） |

⚠ **树不纯净，具名写明**：起手 `git status --porcelain | wc -l` ＝ **770** 行在飞，其中 ` D design/` 计数＝**16**（别人删的，含 `design/assets/tokens.css`），另有 `.gitignore` 与多枚别人的 `.scratch/**` 探针件。本腿**未还原、未提交、不据它们判绿**；上表的 `rc` 是在这棵脏树上量到的（d22scan 在脏树上仍 clean，与 r1 同读）。本腿对 `frontend/**`／`design/**` **零读零写**——`design/assets/tokens.css` 的存在性只走 git 对象层：
`git cat-file -e HEAD:design/assets/tokens.css` → **rc=0（HEAD 里那枚文件存在）**；`git status --porcelain -- design/assets/tokens.css` → ` D design/assets/tokens.css`（工作树被删、未入库）。

## 1. 产码笔名册尺（⛔ 只含那 4 枚件）

`git show --name-only --format="" 73c30dfe` 原样（4 行）：
`cmd/wisp/subagent_carrier_197_test.go`／`internal/panel/subagent_blocked_220_test.go`／`internal/panel/subagent_roster_197.go`／`internal/panel/subagent_roster_197_test.go`。

三笔逐笔名册（`git show --name-only --format="" <c>`，`898d1f89`＝3 枚 probes 件、`913161ae`＝9 枚 probes 件）合起来过一把禁列尺：

```
git diff-tree --no-commit-id --name-only -r 898d1f89 73c30dfe 913161ae \
  | grep -E '^(docs/|frontend/|design/|testdata/)|tokens_fourway_test|l2_grant_boundary|ticket90_persist|thresholds\.go|allowlist\.txt|internal/agent/loop\.go'
```
→ 输出 **0 行**，`rc(docs-forbidden-in-r1-commits)=1`（＝零命中）。

⚠ 一条容易读错的形状要具名：`git diff --name-only fd692f7e..913161ae` 会把**别人的笔**一起带出来（`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、票 290/291 的票面、`probes/290|291|orch` 等），因为共享工作树在 r1 的三笔之间被别的腿推进过（`git rev-parse --short 898d1f89^` ＝ `fd692f7e`，而 `3adabba0` 那一笔台账是别人的）。⇒ **"写面只含那 4 枚件"这一条只能用逐笔 `git show`／`git diff-tree` 量，不能用区间的 `git diff` 量**；用区间量会得到"r1 写了 docs"的错误结论，本腿两种都跑了，逐笔尺干净。

## 2. 禁区零字节尺（工作树相对当前 HEAD）

```
git diff HEAD --numstat -- internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go \
  internal/perm/ticket90_persist_test.go internal/observe/thresholds.go tools/d22scan/allowlist.txt
```
→ 输出 **0 行**，`rc=0`。同理 `git diff HEAD --name-only -- internal/agent/loop.go … docs/PLAN.md docs/specs` ＝ 0 行（`raw-ac4-gates.md` §6，`rc(forbidden-count)=0`），
`git diff HEAD --name-only -- '*/testdata/*' testdata '**/golden/**'` ＝ 0 行（§7，`rc(testdata-count)=0`）
⇒ 三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／`testdata/**` **零字节**，`internal/agent/loop.go` 的 `:367`／`:594` 两句实话没被动（那枚文件整体零变化）。

## 3. 顺手量到的一条与本票无关的红因（具名，⛔ 不算本票账）

`TestC21DesignTokensFourWayAgree` 在**本腿自己**那一发改后全量跑里仍红，失败行逐字（`raw-v1-post.md:464`）：

```
    tokens_fourway_test.go:441: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified.
```

⇒ 它的红**只由别人在工作树里删掉的 `design/assets/tokens.css` 造成**（HEAD 里该文件存在，见上），⛔ 不许读成票 289 造成的红，本腿也不还原它。

## 4. Git 纪律自证（本腿）

零 push（机主从未授权）；两笔 commit 均带**显式 pathspec**（`git add -- <逐条路径>` → `git commit -- <逐条路径>`）；⛔ 无 `git add -A`／`git add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内零删文件、零 worktree；临时件只建不删。
⚠ 一处小形状要具名：起手第 1 笔我直接 `git commit -- <新件路径>` 失败（`pathspec ... did not match any file(s) known to git`，因为件还没进索引），改为 `git add -- <同一条显式路径>` 后再 commit 成功——这不是违规，只是"新建件必须显式 add 一次"，记下来免得下一位再撞。
