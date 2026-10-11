# `309-r3` 门禁（四道各自落一行 `rc=N`，件⛔ 0 字节）

每把尺的件都在 `logs/`，末行都是 `rc=` 本身（⛔ `cmd | head; echo $?` 那种取错退码的写法：每一发都是"先落文件、再取退码"）。
台面：**共享工作树**（`internal cmd tools scripts` 全程 0 行脏，见 `20-teeth.md` 末节与 `logs/rig-shapes.log` 末节；本腿⛔ 在树里留过一枚变异件——六发全打在仓外导出树）。

## 1. 四道门禁（第一轮＝两件代码笔落完之后、`10`／`20` 两件证据件落盘之时）

| 门禁 | 尺（逐字） | rc | 件（字节） |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **`rc_build=0`** | `logs/gate-build-mid.log`（输出 0 字节＝无诊断 ＋ 一行 `rc_build=0`） |
| vet | `GOFLAGS= go vet ./internal/tools/` | **`rc_vet=0`** | `logs/gate-vet.log`（同上，含 `rc=0` 行） |
| d22scan | `sh scripts/d22scan.sh` | **`rc_d22scan=0`** | `logs/gate-d22scan.log` 末两行逐字含 live 那句 `d22scan: clean - no D22 ban violations` ＋ `rc_d22scan=0` |
| path-length | `bash scripts/check-path-length-budget.sh --with-self-test` | **`rc_pathlen=0`** | `logs/gate-pathlen.log`（`VERDICT GREEN - every over-budget tracked path is rostered by name …`） |

★ban #8 那一族字形（`⛔ ≈ ✓ ≤` 与箭头/表情）本腿一次都没写进 **Go 字符串字面量**：新增的字符串只有 shape 标签（ASCII）、控制码源串（ASCII）与一句中文错误文案。
d22scan 两轮都 `rc=0` 就是这一条的读数（第二轮见下面第 3 节）。⛔ 动 `tools/d22scan/allowlist.txt`（尺＝`git status --porcelain -- tools scripts internal/observe`＝**0 行**，第 4 节）。

## 2. 格式名册两把并排（尺＝`logs/gate-fmt.log`）

| 尺 | 射程 | 工作树面 | HEAD blob 面（`git show HEAD:<f>` 逐枚导进仓外临时目录后量） |
|---|---|---|---|
| `gofmt -l internal/tools` | `internal/tools/**`（78 枚 tracked `.go`，尺＝`git ls-files internal/tools \| wc -l`） | **0 枚**（empty above＝0 枚） | **0 枚** |
| `gofumpt -l internal/tools` | 同上（`gofumpt` 在场＝`$(go env GOPATH)/bin/gofumpt`，⛔ 免检） | **0 枚** | **0 枚** |
| `gofmt -l .`（整仓，只作旁证） | 整仓 | **37 枚**——起手就在册、⛔ 归本腿（本腿格式面新增＝0：两把 scoped 尺都 0⇄0） | — |

⇒ 本腿新增格式违例＝**0 枚**（判据＝同一把 scoped 尺在工作树⇄HEAD blob 两面都 0）。
HEAD 面此刻＝本腿 (B) 那笔 `a1f84d39`；两面的 `internal/tools` 内容一致，因为本腿交回前⛔ 再动过产码。

## 3. 第二轮（写完自家证据件之后重跑，防"证据件自己把门禁打红"）

| 尺 | rc | 件 |
|---|---|---|
| `sh scripts/d22scan.sh`（此时 `00`／`10`／`20` 三枚 `.md` 已在盘上） | **`rc_d22scan_with_my_files=0`** | `logs/gate-d22scan-with-my-files.log` |
| `bash scripts/check-path-length-budget.sh --with-self-test` | **`rc_pathlen_with_my_files=0`** | `logs/gate-pathlen-with-my-files.log`；承重行逐字＝`denominator: tracked paths=10828  over-budget=57  covered by roster=57  not in roster=0`、`bands: over the hat=57 … roster entries=57`、`VERDICT GREEN` |

⚠ `tracked paths` 从 `A851` 那次的 10779 漂到 **10828**＝本波其它腿与本案的件一起进仓的分母漂移，⛔ 归因给本腿也⛔ 洗掉；判据只看 `not in roster=0` 与 `VERDICT GREEN`。

## 4. 越界名册（写面⛔ 出派单 §2.3）

| 尺 | 读数 |
|---|---|
| `git diff --name-only 3974ab4d HEAD`（代码两笔＋证据笔累加） | 名册＝本腿两枚产码／测试件（`internal/tools/bridge.go`、`internal/tools/toolrequest_nail_309_test.go`）＋`.scratch/wisp/probes/309/r3/**`＋（末笔）票 309 一件 |
| `git status --porcelain -- tools scripts internal/observe internal/risk` | **0 行**（empty above＝0 枚）⇒ allowlist／`thresholds.go`／`internal/risk/**`／`scripts/**` 一枚⛔ 动 |
| `find .scratch/wisp/probes/309/r3 -name '*.go' \| wc -l` | **0**（⛔ 在 `probes/**` 下建 `.go`；台件＝`.txt`、读数＝`.log`） |
| 三枚冻结件（`internal/panel/tokens_fourway_test.go`·`internal/panel/l2_grant_boundary_test.go`·`internal/perm/ticket90_persist_test.go`） | 本腿⛔ 打开过；尺＝上面那发 `git status --porcelain` 覆盖⛔ 到 `internal/panel`、`internal/perm`，改由名册尺核：本腿两笔代码笔的 `git show --name-only` 各＝**1 枚文件**（见 `logs/commit-code.log`、`logs/commit-bridge.log`） |

## 5. ⛔ 放宽断言／⛔ 造 SKIP 的自证

- 名册尺＝`logs/roster-control.log`：`SKIP=0`⇄`SKIP=0`（两枚计数件都空＝empty above＝0 枚）、`newly_red_count=0`、顶层 211⇄211、子测试 68⇄68。
- 本腿⛔ 改任何 `want*` 常量：`wantTotal=1`／`wantOwner="internal/agent/loop.go"`／分母地板 `200` 三枚逐字未动（尺＝`git diff 3974ab4d a1f84d39 -- internal/tools/toolrequest_nail_309_test.go` 里 `grep -E '^[+-].*(wantTotal|wantOwner|examined < )'` ＝**空名册，0 枚**；读数一并落在 `logs/`，见 `30-gates` 提交里的 `gate-assertions.log`）。
