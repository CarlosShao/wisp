# 342-r3 交付件：`privateFile` 祖先那一支的两平台成对读数（实现腿＝只测量，产码 `.go` 零字节）

- 锚＝`7eb0e3f4b`（起手 `git rev-parse HEAD` 现认，2026-11-05 13:5x）。
- 写面＝本件 1 枚；`internal/winsec/**` **零字节改动**（收尾现量：`git status --porcelain internal/winsec` 空读数＝0 行）；两枚跑的就是盘上的仪器，**一字未改**。
- 结论＝那一条腿在 HEAD 上**已经落地并在两平台各自量到**，本腿不重落形，只把票 342 `AC#1` 要的**成对读数**（修前分子 > 0 ⇄ 修后 0 且 liveness 非零，两平台各一发以上）补齐，同时把两处口径缺陷具名报上去。

## 0. 先说清楚"为什么没有产码改动"

祖先那一支今天由两枚先前的落地腿关着，都在 HEAD 上：

| 腿 | 落了什么 | 在哪 |
|---|---|---|
| `342-r1` | 先钉父目录再动作：`pinnedParent` ＋ `openPinnedLeaf`（`openat`）＋ `unlinkPinnedLeaf`（`unlinkat`），父钉的 `close` 摆在路由末端 | `winsec.go:170-226`、`pinnedleaf_342_unix.go`、`pinnedleaf_342_windows.go` |
| `400-r1` | Windows 的 walk 改下探**调用方拼写**（`pinnedParent(path, resolved)`），封掉"C26 先读 reparse 判定、再读最终路径"两读之间换位让答案穿链成外来拼写那扇窗 | `winsec.go:194`、`pinnedleaf_342_windows.go` 的 `pinnedParent(caller, resolved)` |

所以本腿的可做面＝**读数**（`AC#1`）与门禁（`AC#4` 的本波射程）。票面 `AC#6` 写明七格由非实现者裁，本腿**一枚没翻**（起手 5 勾 / 2 未勾，收尾同）。

## 1. 台件：三棵仓外导出树，⛔ `-overlay`、⛔ 复用别人的树

全部 `git archive HEAD | tar -x` 到 `/d/tmp/wisp/20261105/342r3/`（仓外，⛔ 入库）：

| 树 | 是什么 | 怎么造 |
|---|---|---|
| `tree` | **修后**（HEAD 原样） | `git archive HEAD` |
| `treeOld` | **修前** | 只把 `winsec.go` 的 `privateFile` 体换回 `fde1df02` 的按名形（`os.OpenFile(resolved, …, O_TRUNC)` ＋ 拒绝路径 `os.Remove(resolved)`），其余字节全 HEAD |
| `treeOldNoUnlink` | **修前的隔离发** | `treeOld` 再把两处 `_ = os.Remove(resolved)` 摘掉（那是票 348 的动作），让"那次 open 的截断"与"拒绝路径的删除"分开计 |

- POSIX 档＝容器 `golang:1.27`、`CGO_ENABLED=0`、`GOPROXY=off`、`GOMODCACHE` 只读挂宿主、`TMPDIR`/`GOTMPDIR` **留容器内 `/tmp`**（⛔ 绑 Windows 卷——`342-v1` 那枚定式：bind mount 上 `os.SameFile` 失真，身份门会被台面摘掉）。
- Windows 档＝本机 `go1.27.1 windows/amd64`，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`。
- 逐发日志在 `/d/tmp/wisp/20261105/342r3/logs/`（`posix-head-1..3`、`posix-old-1..2`、`posix-oldnounlink-1..2`、`win-head-1..3`、`win-old-1..3`）。

## 2. `AC#1` 成对读数

### 2.1 POSIX（`TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`，1500 发/发）

| 树 | truncated | deleted | pinnedParentBeforeOpen | provedLinkCycles | 判 |
|---|---|---|---|---|---|
| 修前 `treeOld`（删支在） | 0 / 0 | **83 / 100** | 0 | 196 / 185 | RED |
| 修前隔离 `treeOldNoUnlink`（删支摘掉） | **6 / 9** | 0 / 0 | 0 | 275 / 205 | RED＝截断支的修前分子 > 0 |
| 修后 `tree`＝HEAD | **0 / 0 / 0** | 0 / 0 / 0 | 32 / 16 / 31 | 261 / 174 / 209 | PASS ×3 |

两条"修前"读数的关系是这张表最要紧的一句：预置路线在 POSIX 上**几乎总是先截断、再被拒绝路径的删除把名字拿走**，所以删支留着时 `truncated` 读 0 而 `deleted` 读 83/100；把删支摘掉的隔离发才让那次 `O_TRUNC` 单独露出来（6/9 发，首次损伤 body 逐字 `size 0, content ""`＝被截到空）。⇒ **"修前 truncated=0"在这条路上 ⛔ 能读成"没伤"**（票 342 `AC#1` 逐字要的就是这条口径）。

修后三发的 refusal 文本逐字（同一句，parent 面而非 leaf 面）：
`refused by the pinned parent before the open: winsec: path is not provably resolved, refusing to seal: …/data/sub: pinned object at sub is not what the name carries: the name carries a link, so the pinned object is not what the spelling reaches (winsec: cannot apply a private security descriptor)`。

### 2.2 Windows（`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`，两支形各 400 发/发）

| 树 | 形 | truncated | provedReparse | formBLegRefusals | 判 |
|---|---|---|---|---|---|
| 修前 `treeOld` | junction | **1 / 1 / 3** | 15 / 12 / 17 | 0 | RED |
| 修前 `treeOld` | symlink-dir | **8 / 8 / 14** | 58 / 58 / 66 | 0 | RED |
| 修后 HEAD | junction | 0 / 0 / 0 | 18 / 17 / 15 | 7 / 6 / 3 | PASS |
| 修后 HEAD | symlink-dir | 0 / 0 / 0 | 85 / 76 / 67 | 35 / 16 / 15 | PASS |

修前损伤 body 逐字＝`content "raced private payload 342r1 windows ancestor"`（44 字节产码载荷）＝真产码漏，⛔ 仪器自造。

## 3. 门禁（跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 门 | 命令 | 读数 |
|---|---|---|
| 构建 | `go build ./internal/winsec/` | rc=0 |
| 消费者构建 | `go build ./internal/...` | rc=0（`internal/risk`⇄`secret`⇄`config`⇄`agent`⇄`memory` 全在编译面内） |
| 测试 | `go test -count=1 -timeout 240s ./internal/winsec/` | `ok … 29.153s`（同波另两发 `ok … 23.554s`） |
| 格式化 | `gofumpt.exe -l internal/winsec/` | **空**、rc=0（本腿 0 枚 `.go` 改动） |
| D22 | `sh scripts/d22scan.sh` | rc=0 `clean`（ban #8 分母 internal=884 / cmd=207） |
| vet（本机） | `go vet ./internal/winsec/` | rc=0 |
| vet（交叉档） | `GOOS=linux go vet ./internal/winsec/` | rc=0 |
| vet（交叉档） | `GOOS=darwin go vet ./internal/winsec/` | rc=0 |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **rc=1 VERDICT RED**（⛔ 本腿造成，见第 4.3 节） |

`go build ./...` 在**共享工作树**上会撞 `build/mut_rvw.go:446: undefined: voiceBargeinMonitor` 等一串：`build/` 是 `.gitignore:14` 忽略掉的**未跟踪草稿目录**（`git check-ignore -v build/mut_rvw.go` 现验＝`.gitignore:14:build/`），不在跟踪面里；跟踪面用 `go build ./internal/...` 取，rc=0。⛔ 把这个读到红。

## 4. 本腿抓到的三处口径缺陷（⛔ 改，具名报上去）

### 4.1 `ancestoropen_342_unix_test.go:80-101` 的归属支**说反了**

盘上是 `landed := st.pinnedParentLeg > 0`（:80），`if !landed`（:90）判红，红句逐字写"…so this row was produced by bytes that **DO** carry the pinned route…"（:91-101）。可 `pinnedParentBeforeOpen == 0` 恰恰是"这些字节**没有**父钉"——真修前树上它读出来的是**反话**。本腿现量：`treeOldNoUnlink` 上 `PinnedParentBeforeOpen=0`、`truncated=6/9`，红句却断言这行"带父钉"。

同一文件的头注（:34-45）又写"an **old-side** row whose `pinnedParentBeforeOpen > 0` fails as 'not attributable'"——与代码的判据方向相反（代码在 `>0` 时**放行**并按"落地行"去 gate `truncated==0`）。⇒ 仪器今天在 HEAD 上是恒绿死支（HEAD 总有父钉），在真修前树上给的是**归属红而非损伤红**。本腿⛔ 改它：它是判据面，改它要裁决者点头。

### 4.2 `winsec.go:139-161` 的自陈已和同一函数的产码相反

`:149-154` 仍逐字写 "the ANCESTOR leg is still open on POSIX … This window is **not claimed closed** anywhere"，而 `:170-193` 落的正是父钉路由、`:194` 传的是调用方拼写、`:202` 的 open 在钉住的父 fd 上。`342-v3` 已裁"更新权在终裁之后"，本腿**不动它**；现在两平台读数都齐了（见第 2 节），该不该改成终态、怎么改，归裁决者/编排者。

### 4.3 长度门 RED 是**先前有**的、且不在本腿写面

`over-budget=65 / roster=64 / not in roster=1`，点名两行同一枚：
`.scratch/wisp/issues/179-loop-deciderisk-treats-a-declared-l0-as-unclassified-so-every-readonly-builtin-is-refused-on-the-real-cli-and-no-criterion-holds-that-sentence-done.md`（相对 175 字符、name 154、自托管 runner 上 219 > 165）。⇒ 是 `.scratch/wisp/issues/**` 的台账名与 `scripts` 名册漂移，归 `scripts`＋票面面；本腿 0 枚 `.go`／0 枚台账改动，碰不到它。

## 5. ⛔ 声称的（照旧具名，⛔ 圆场）

- Windows **叶子** 那一支（票 322 的残差）不在本读数射程内：本件只读**祖先换手**。
- `FILE_SHARE_DELETE` 顶名窗（`342-v3` 具名、归新票族）未动、未量。
- `cmd/wisp` 整包与 `go test ./...` 未跑（派单令）；`AC#4` 的"消费者五包成对＋`cmd/wisp` 那一发"仍归续体。
- 票 342 的框**一枚没翻**（归非实现者裁）；票名／`Status:`／`-done` 一字未动。
- 本腿没有新增任何产码或测试字节；本件是本腿唯一入库物。
