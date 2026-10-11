# 308-r2 · 30-gates（四把门禁 rc ＋ 格式名册两把并排 ＋ 整包名册作差）

所有件自落 `rc=` 行（`logs/*.rc.txt`）。台面＝工作树（改后）＋**派发锚 `202ed4a5` 的仓外导出树**（改前那一发⛔ 在共享工作树跑）。

## 1. 四把门禁

| 门禁 | 命令 | rc（最终内容） | 件 | 备注 |
|---|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | **0** | `logs/gate-build-final.txt`（stdout **0 字节**＝绿时本就空）＋早一发 `logs/gate-build.txt` rc=0 | ⛔ 跑 `cmd/wisp` 那包（sherpa DLL 前置归编排者）；本发是 `go build`，⛔ 跑测试 |
| vet | `GOFLAGS= go vet ./internal/risk/` | **0** | `logs/gate-vet-final.txt`（0 字节） | |
| D22 扫描 | `sh scripts/d22scan.sh` | **0 ⇒ 1 ⇒ 1**（三发，见 §2） | `logs/gate-d22scan.txt`／`gate-d22scan-run3.txt`／`scan-worktree-direct.txt` | 红＝**别腿在飞态**，⛔ 我账；我这面另有一发隔离扫描 rc=0 |
| 路径长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `logs/gate-pathlen.txt` | `VERDICT GREEN`、tracked paths=10457、over-hat=**57** ⇄ roster=**57**、not in roster=**0**、worst full path=**228**、longest rel=184 ⇒ 与票 §规矩 登记的 11:0x 读数同值，⛔ 动别人名册（我的件名全短） |

## 2. ★d22scan 三发与归属（这格⛔ 我盖章"过"，只报读数）

| 发 | 时点 | rc | 红在哪 |
|---|---|---|---|
| run1 | 句 1–3 落地后 | **0** | `clean - no D22 ban violations`（live scope：bans #1-5 internal/=229、cmd/=39、ban #7 internal/tools/=23、ban #8 internal/=531 Go files） |
| run2 | 句 3 收窄后 | **1** | 顶层 `PASS=33 FAIL=2`；两枚 finding＝`internal/tools/bridge.go:887`／`:966`，ban **#9 phantom-citation**（引 `cmd/internal/tools/scripts` 与 `.scratch/wisp/probes/309/r2/30-comments.md`，盘上⛔ 存） |
| run3 | 同内容复跑 | **1** | 同一枚两枚 finding，行号漂到 `:892`／`:971` |

归属尺（三发都⛔ 指向我）：
- `git status --porcelain -- internal/tools`＝**2 行**：` M internal/tools/bridge.go`、` M internal/tools/toolrequest_nail_309_test.go` ⇒ **309-r2 正在写**（起手时该射程 0 行，见 `00-anchor.md` §0）。
- `git show 202ed4a5:internal/tools/bridge.go | grep -cE 'cmd/internal/tools/scripts|probes/309/r2/30-comments.md'`＝**0** ⇒ 那两枚引用是**在飞态内容**，⛔ 是跟踪内容。
- run2⇄run3 同一 finding 的行号从 887/966 漂到 892/971 ⇒ 那一程别腿仍在插行。
- ⚠ 脚本 `set -eu`：步骤 1（`runtests.sh` 正控制，含读实时盘的 `TestScannerSelfScanOfRealRepoIsGreen`／`TestRealRepoLedgerIsHonest`）红 ⇒ **步骤 2（真扫描）根本⛔ 跑到**。所以 run2/run3 的 rc=1 ⛔ 回答"我这面干不干净"。

我这面的独立读数（两把）：
- 直接跑步骤 2 于工作树：`cd tools/d22scan && GOFLAGS= go run . -root <repo>` ⇒ rc=**1**、2 枚 finding **全部**在 `internal/tools/bridge.go`、`grep -c 'internal[/\\]risk'` 命中＝**0**（件 `logs/scan-worktree-direct.txt`）。
- 隔离台面：`git archive 202ed4a5` 导出树（995 枚 `.go`，带 `go.mod`/`go.sum`）＋**只植入我这两枚最终件** ⇒ `go run . -root <tree>` rc=**0**、`clean - no D22 ban violations`（件 `logs/scan-mine-tree-final.txt`；scope 分母与工作树同值：bans #1-5 internal/=229、cmd/=39、ban #7 internal/tools/=23）。⚠ 具名局限＝导出树里⛔ `.git`，扫描器自己打了一行 `gitignore rules NOT APPLIED` 警告 ⇒ 这一发只回答"我的注释⛔ 触发任何 ban"，⛔ 替代 run1 那发全门通过。

⇒ 交回口径：**run1 rc=0 是我这三句落地的门；run2/run3 rc=1 是 309-r2 的在飞态，具名归它，⛔ 我消也⛔ 我盖。**

## 3. 格式名册（gofmt ⇄ gofumpt 各两把并排，每把写射程；先数 CR 再喊脏）

工具：`gofmt`＝`/d/work/base/go/bin/gofmt`；`gofumpt`＝`$(go env GOPATH)/bin/gofumpt.exe`（**PATH 上没有**，具名版本 `v0.12.0 (go1.27.1)`）。件＝`logs/format-rosters.txt`。

| 台面／射程 | gofmt -l | gofumpt -l | 逐枚成色 |
|---|---|---|---|
| 工作树，`internal/risk` | **1** | **1** | `internal/risk/provenance.go` CR=**1115**、`git check-attr`＝`text: set / eol: lf` ⇒ **CRLF 假枚**（与票面 A839/A847 记的同一枚），⛔ 我碰它、`git status` 空 |
| 工作树，`internal cmd tools` | **4** | **5** | 4 枚＝`internal/agent/approval/pending_read.go`(CR=131)／`internal/agent/tools.go`(CR=225)／`internal/risk/provenance.go`(CR=1115)／`cmd/wisp/models.go`(CR=334) ⇒ 全 CR>0 ＝**换行假枚**；gofumpt 第 5 枚＝`cmd/wisp/panel_transport_35r2_test.go`（CR=**0**，⛔ 假枚） |
| 派发锚导出树，`internal/risk` | 0 | 0 | 同一批文件在 LF 台面下**全绿** ⇒ 上面那 4 枚确认为假枚 |
| 派发锚导出树，`internal cmd tools` | 0 | **1** | 只有 `cmd/wisp/panel_transport_35r2_test.go` ⇒ 它在 **HEAD 里就真脏**（gofumpt 意义），`git status` 空＝⛔ 在飞态、⛔ 本票账 ⇒ 具名留给编排者（不在我射程、我⛔ 改） |
| **我的两枚文件** | **0** | **0** | `gofmt -l <两枚>`＝0 行、`gofumpt -l <两枚>`＝0 行、CR＝0 ⇄ 0（件 `logs/strip-result-final2.txt`） |

⚠ 另记一笔（⛔ 我造、⛔ 我删，只登记）：`git status --porcelain -- .scratch/wisp/probes | grep '\.go$'`＝**3 枚未跟踪 `.go`**，落在 `probes/224/v2/posctrl.go`、`probes/33/r5/block_ac13_ac14.go`、`probes/33/r5/logs/prev-test-copy.go` ⇒ **别家腿**的，本腿在 `probes/308/r2/**` 里新建的 `.go`＝**0 枚**（`find` 计数 0）。任何把 `.scratch` 扫进去的宽域格式尺都会先撞上这三枚，⛔ 算我挪的分母。

## 4. 整包名册作差（改前⇄改后，逐名）

改前＝`git archive 202ed4a5` 导出树一发（`logs/pre-head-full.txt`，rc=0）；改后＝工作树最终内容一发（`logs/post-worktree-final-full.txt`，rc=0）。别名台面（`TMP=RISKAT~1\RUNNER~1`，cmd.exe 绕开 MSYS 改写）另有一发 `logs/pre-head-full-alias.txt` 作对照。

| 尺（具名口径） | 改前 | 改后 | 别名台面（改前） |
|---|---|---|---|
| `GOFLAGS= go test ./internal/risk/ -count=1 -v` rc | 0 | 0 | 0 |
| `^=== RUN` | 201 | 201 | 201 |
| 顶层 `^--- PASS`（一把尺） | **132** | **132** | 132 |
| 含子测试的 `--- PASS`（另一把尺，`^[[:space:]]*--- PASS`） | **200** | **200** | 200 |
| 其中子层 `^    --- PASS` | 68 | 68 | 68 |
| 顶层名册枚数（PASS＋FAIL＋SKIP） | 133 | 133 | 133 |
| `--- SKIP` / `--- FAIL` | 1 / 0 | 1 / 0 | 1 / 0 |
| `panic: runtime error`（词面） | 0 | 0 | 0 |
| **名＋色作差**（剥 `(0.01s)` 后排序 `diff`） | — | — | **空，diff rc=0**（件 `logs/roster-pre.txt`／`logs/roster-post-final.txt`） |

⇒ ⛔ 新增 FAIL、⛔ 新增 SKIP、⛔ 少名、⛔ 变色；唯一色变＝⛔ 有（三处计时数字漂移⛔ 是色）。哨兵 A/B 在两发里都＝0（⇒ 句 2 那句"到达数＝0"仍立）。⚠ 词面尺具名：`grep -ciE 'panic|slice bounds'`＝**4** 但那 4 行全是用例名里的 `Panic` 字样 ⇒ 承重尺＝`grep -c 'panic: runtime error'`＝0。
