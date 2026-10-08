# gate-snapshot-1 · 01 门禁快照（只读，零修改）

- 快照时刻：2026-10-08 18:56–18:58 +0800
- 引数锚 HEAD：`50b34971d592c6b9e161181dd01f19afda0992fc`（`50b34971`，2026-10-08 18:57，本腿起手锚 commit；其父＝`0f7d4c5`）
- 声明：本轮只跑读命令；⛔ 未跑 go build／go test／go vet；⛔ 未改任何既有文件、未碰他腿脏改动；唯一写面＝本目录 `*.md`。每把均自检 rc；分母随锚点，**引用时报 HEAD＝50b34971**。

## 1 d22scan（独立 module，经 `scripts/d22scan.sh` 走）

命令原文：`sh scripts/d22scan.sh`（输出落 `/tmp/gs1-d22scan.log`，仓库外；临时件只建不删）
→ **rc=0**，日志 251 行。末 11 行（原文）：

```
d22scan.sh: scan of /d/work/workspace/projects plans/Wisp
d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (2 path(s))
d22scan: examined 266 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope bans #1-5 internal/      examined 228 production Go files
d22scan: scope bans #1-5 cmd/           examined  38 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  23 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 516 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined 113 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=516, ban #8 cmd/=113; ban #8 emoji coverage: design/ 39 text files; frontend/ 85 text files; internal/ 516 Go files, comments and _test.go included; cmd/ 113 Go files, comments and _test.go included
```

（上引为末 11 行；`bans #1-5` **不含** `_test.go`，`ban #8` **才含**。日志前段另有自带 runtests：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`。）

## 2 gofmt

命令原文：`gofmt -l cmd internal`（⛔ 未加 `-w`、本腿不修）
→ **rc=0**，命中 **7 枚**。逐枚尺＝`git log -1 --format='%h %ad' --date=format:'%m-%d %H:%M' -- <file>` ＋ `git status --porcelain -- <file>`：

| 命中文件 | 最后落盘 | 工作树 |
|---|---|---|
| `cmd/wisp/models.go` | `5e8748b3` 10-03 21:07 | 干净 |
| `cmd/wisp/panel_inbound_guards_35r3_test.go` | `7f9d6e40` 10-07 21:25 | 干净 |
| `cmd/wisp/panel_transport_35r2_test.go` | `3a343bc7` 10-08 09:37 | 干净 |
| `internal/agent/approval/pending_read.go` | `d03d166f` 10-04 10:37 | 干净 |
| `internal/agent/tools.go` | `5413f46d` 10-04 09:19 | 干净 |
| `internal/risk/provenance.go` | `5e8748b3` 10-03 21:07 | 干净 |
| `internal/tools/bridge.go` | `5413f46d` 10-04 09:19 | 干净 |

判定：**6 枚＝既有**（最后落盘早于今天：10-03／10-04／10-07）；**1 枚＝今天落盘**（`3a343bc7` 10-08 09:37，已提交、工作树干净）。7 枚均非"今天新写入未提交"形。

## 3 两道常量/名册尺

① 命令原文：`grep -c 'answered' internal/panel/l2_grant_boundary_test.go` → **rc=0，计数=62**。
（仅**射程计数，非内容引用**；该测试文件此刻在 259-r4 的写射程上，数值随其进度漂移；本腿未改其一字。）
② 命令原文：`git grep -c 'go func(' HEAD -- 'internal/*.go' 'cmd/*.go' ':!*_test.go'` → **rc=0，词面命中总数=1**，唯一位点 `cmd/wisp/testdata/esclistener/main.go`（计 1）。
⚠ 这只是**词面命中数，不是判定**（ban 判的是"裸 `go func(` 无 owner/recover"）——**这不是判定，裁决要人看**。

## 4 三面空否

命令原文：`git status --porcelain -- cmd/wisp internal/panel internal/tools internal/agent/approval` → **rc=0，输出全空（0 行）**。
（截至 18:58，四路径未见未提交改动 ⇒ 盘面未显示 259-r4 的落笔，属时点差可能；本腿不解读。）

## 5 未推计数

命令原文：`git rev-list --count origin/dev..HEAD` → **rc=0，count=427**（自取时刻＝2026-10-08 18:58:00 +0800；含本腿 `50b34971`；本地 `origin/dev` 引用、未 fetch，**会漂**）。
`git log -1 --format='%H %ad'` → `50b34971d592c6b9e161181dd01f19afda0992fc 2026-10-08 18:57`。

## 附：冻结件与写面（本腿口径）

- 本腿**一字未动**任何文件；对 `internal/panel/l2_grant_boundary_test.go` 只做 `grep -c` 计数＝**射程判断，非内容引用**。
- 台账里"三枚冻结件"的枚枚举在本轮射程搜索结果中未见定义 ⇒ **不猜名**；对 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt` 等冻结面零接触。`frontend/**`／`design/**` 不读不引。
- 本腿 commit：`50b34971`（00-anchor.md）＋本文件一笔（paths 显式、零 push）。

## 附：与派单转述的冲突

无。五把均按盘上原文跑通；唯一补充＝(4) 全空为时点读数、(3)① 计数随 259-r4 并发漂移，均已注明。
