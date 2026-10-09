# 票 289 / 腿 289-r1 — 门禁（AC#4）与越界自证

> 每枚门禁件自己落一行 `rc=N`。零字节的件＝那一格没交。
> 时钟（date stdout 原样，不手打）：`date -u -d "+8 hours" "+%Y-%m-%d %H:%M +08"` → `2026-10-09 14:06 +08`（起手）；
> 共享工作树的 HEAD 在本腿途中被别的腿推进：`34a01259` → `fd692f7e`（现量复量与全部门禁都在 `fd692f7e` 之上起跑）。

## 各把尺的读数

| 尺 | 命令逐字 | 读数 | rc |
|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | 空输出 | **rc=0** |
| gofmt | `gofmt -l internal/panel/subagent_roster_197.go internal/panel/subagent_roster_197_test.go internal/panel/subagent_blocked_220_test.go cmd/wisp/subagent_carrier_197_test.go` | 空输出（4 枚件无一被列） | **rc=0** |
| gofumpt | `"$PWD/$(go env GOPATH)/bin/gofumpt.exe" -l <同那 4 枚件>` | 空输出。⚠ 具名：`gofumpt` 不在 PATH（`command not found`，`rc=127`），落点在 `$(go env GOPATH)/bin/gofumpt.exe`，本腿用绝对路径跑的，⛔ 没有把那一格当成"跳过" | **rc=0** |
| d22scan | `sh scripts/d22scan.sh`（正控 `runtests.sh` ＋ `go run . -root "$root"` 两步，脚本内 `set -eu`，无 `|| true`） | `d22scan: clean - no D22 ban violations`；`examined 268 production Go files under internal/ and cmd/`；ban #8 覆盖 `internal/` 524 ＋ `cmd/` 116 枚 Go 文件（注释与 `_test.go` 计入）⇒ 本腿改的 4 句注释在射程内且未触雷 | **rc=0** |
| 改后 go test | 见下节 | 见下节 | 见下节 |

d22scan 全量输出＝本目录 `raw-d22scan.md`（末行 `rc(d22scan)=0`）。
⚠ 具名一句：本次**不是纯净树**（工作树里另有别人的在飞改动：`design/**` 若干删除、`.gitignore` 修改、别人的 probes 件），
本腿没有还原它们、也没有把它们算进自己的写面；在那样的树上 `sh scripts/d22scan.sh` 仍 `rc=0`。

## 越界自证（写面只有那 4 枚件）

- 尺＝`git diff HEAD --numstat -- <那 4 枚件>`（工作树相对当前 HEAD 的净变化，逐枚）：

```
4	4	cmd/wisp/subagent_carrier_197_test.go
2	2	internal/panel/subagent_blocked_220_test.go
3	3	internal/panel/subagent_roster_197.go
5	5	internal/panel/subagent_roster_197_test.go
```

- 尺＝`git diff --cached --numstat`（索引）＝**空** ⇒ 本腿没有 `git add` 过任何别人的件，
  提交用的是 `git commit -- <显式 pathspec>`，⛔ 无 `git add -A`／`git add .`。
- 禁区名册逐枚（`git diff HEAD --name-only -- ...` 逐字路径列表）：
  `internal/agent/loop.go`／`frontend`／`internal/panel/tokens_fourway_test.go`／
  `internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`／
  `internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／`docs/PLAN.md`／`docs/specs`／
  `cmd/wisp/panel_inbound_guards_35r3_test.go` ⇒ **输出 0 行**（见下面 §名册尺）。
- `design/**`：⚠ 这一族的删除/改动**先于本腿存在**——起手第一发 `git status --short`（在本腿写下任何一处改动之前）
  就已经列出 `D design/assets/base.css`／`D design/assets/icons.js`／`D design/assets/theme.js`／
  `D design/assets/tokens.css`／`D design/index.html`／`D design/screens/*.html` 与 `M design/doubao/**`。
  那是别人在飞的活：⛔ 本腿未还原、未提交、不据它判绿；本腿对 `frontend/**`／`design/**` **零读零写**
  （连结论都不引），`git show --stat` 名册尺（下一节）里也不会有它们。
- 零读零写自证：本腿读过的件名册＝票面本身、`internal/agent/loop.go`（只读注释与 callCorr）、
  `internal/panel/subagent_roster_197{,_test}.go`、`internal/panel/subagent_blocked_220_test.go`、
  `cmd/wisp/subagent_carrier_197_test.go`、`internal/tools/subagent_197.go`、`internal/tools/task.go`、
  `internal/panel/pump.go`、`internal/tools/ticket283_corr_identity_rulers_test.go`（第 5 枚件，只读未动）、
  `internal/tools/bridge.go`（grep 行）、`scripts/d22scan.sh`、`.gitignore`。**不含 frontend/design 任何一枚。**

## 改后 `go test`（AC#3 的后一半，读数和作差在 `30-test-comparison.md`）

- 命令＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -v ./internal/panel/ ./cmd/wisp/ -count=1`（CWD＝仓根），
  原始输出＝本目录 `raw-post.md`。
- **rc=1**（由被量命令自己的 `$?` 落，末行 `rc=1` 原样在 `raw-post.md` 里；⛔ 不是管道尾元素的退码）。
- 三数（顶层）：PASS=390 ／ FAIL=11 ／ SKIP=2；改前 389／12／2 ⇒ **新增红 0 枚**（名册逐名作差＝∅），
  另有 1 枚改前红在改后转绿（`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`，
  本腿靶向复跑 3 发均 rc=0 并把它具名归到串跑时序，⛔ 不记成本腿修的、⛔ 不据它判绿）——
  全过程与原文见 `30-test-comparison.md`。
- 两发之间 `git diff --name-only fd692f7e..HEAD | grep "\.go$"` ＝ 0 行 ⇒ 两发改的不是被测物。

## `git show --stat` 名册尺（产码那一笔只含那 4 枚件）

`git show --name-only --format= 73c30dfe | sort` 原样：

```
cmd/wisp/subagent_carrier_197_test.go
internal/panel/subagent_blocked_220_test.go
internal/panel/subagent_roster_197.go
internal/panel/subagent_roster_197_test.go
```

`git show --stat` 的合计行＝`4 files changed, 14 insertions(+), 14 deletions(-)`（⛔ 增删不等就是失败，这里相等）。
本票三笔 commit：`898d1f89`（证据件骨架）→ `73c30dfe`（产码 4 枚件，显式 pathspec）→ 第 3 笔＝证据件入库（显式 pathspec，本件所在那一笔）。
⛔ 零 push（机主从未授权）；⛔ 无 `git add -A`／`git add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；
⛔ 仓内零删文件、零 worktree。
