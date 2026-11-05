# 339-r2 落地读数：票 339 第 7 格（`AC#2`）两平台去牙反控已可达＝补世界，不是藏红

- 落地腿：`339-r2`（实现者；本件只给读数，**判定归非实现者**，本腿一格不翻）
- 起手锚 `557096c82`；读数期间共享工作树只被本腿改三件（`git diff --name-only` 两条 + 一条新增，见第 3 节）
- 写面：`internal/risk/anchor_miss_reason_334_test.go`（票面点名的射程）＋ `internal/risk/anchor_miss_reason_334_windows_test.go`（**新增，唯一的带档文件**）＋ `internal/risk/anchor_miss_reason_339_other_test.go`（票 339 自己那枚 POSIX 件，`339-r1` 建立）＋ `.scratch/wisp/probes/339/r2/**` ＋ 本件 ＋ 票尾一条只追加 bullet
- 零 `.go` 产码字节、零新导出名、⛔ `cmd/wisp`、⛔ `frontend/**`、⛔ 动票 334 一字、⛔ 翻框、⛔ push

## 0. 起手第一笔：阻塞链现跑（347 已清，不是转述）

| 事实 | 读数 | 尺 |
|---|---|---|
| 票 347 已结案 | 票面件现名 `347-...-done.md`（`Status: DONE-带残余`）；机主 2026-10-31 14:1x 追认（台账 `A1311①`）、2026-11-05 15:0x 再次「按推荐」（`A1541` 第 2 条） | `ls .scratch/wisp/issues/ \| grep 347`；`A1541` 末行逐字"下波派 51-B／**339（依赖347）**／R9" |
| 解析器真落地 | `internal/risk/pathresolver_other.go:41-54` 现读 `resolveHandle` = `filepath.EvalSymlinks` ⇒ 成功分支返回真路径；`reparseComponents` = 逐组件 `Lstat` 拒 symlink | `sed -n 1,55p internal/risk/pathresolver_other.go` |
| 所以 A 档在 POSIX 有世界 | `formsOf` 的 `f.root` 赋值点唯一（`pathresolver.go` 的 `resolveHandle` 成功分支）；`uncertainAnchorMiss`（`blacklist.go:377`）首行 `if cand.root == "" \|\| anchor.root == ""` 现在两侧都非空 | 两平台运行读数见第 4 节；POSIX 侧红句里逐字出现"its own walk stopped at …" |

## 1. 票 339 现况：剩的恰是 `AC#2`

`grep -cE "^- \[ \] \*\*AC#"` = **1**、`grep -cE "^- \[x\] \*\*AC#"` = **6**（起手与收尾各跑一次，同数）。

那一格（票面逐字）要的是：**去牙反控＝把回落支整支摘掉 ⇒ 修好后的仪器必须在两平台都点名红**。
本腿交的正是这一维；判定仍归非实现者。

## 2. 改前：仪器只在一个平台有"世界"（RED-before 逐字）

改前 `internal/risk/anchor_miss_reason_334_test.go:1` = `//go:build windows` ⇒ Linux 上那三枚用例**不编译、不跑**（名册里没有名字）。

只摘 tag、不补世界（其它字节一字不动）⇒ Linux 当场红 5 句：

```
GO_TEST_RC=1  RUN=6 PASS=1 FAIL=2 SKIP=0
anchor_miss_reason_334_test.go:81:  fixture guard: the fail-closed shape must still be A here, got class=B allow=false needL2=true reason="B-list sensitive path; single-file override requires one L2 confirm: id_*"
anchor_miss_reason_334_test.go:106: fixture guard: want A from the fail-closed clause, got class=B ...
anchor_miss_reason_334_test.go:143: class: got B, want A (the fail-closed tier is unchanged by AC#1)
anchor_miss_reason_334_test.go:147: Gate: got class=B ... want class=A needL2=false allow=false (non-overridable)
anchor_miss_reason_334_test.go:151: an A deny from the fail-closed clause was unlocked by an override key: class=B allow=true ...
```

- 日志：`.scratch/wisp/probes/339/r2/logs/linux-red-tag-only.log`（3308 字节）
- 台件：仓外导出树（`git archive HEAD`）、`golang:1.27`、`CGO_ENABLED=0`、`GOPROXY=off`、挂宿主 module 缓存只读、`TMPDIR=/tmp/gt`

这就是票 339 报的那枚病：`pathForms.partial` 的来源在 POSIX 与 Windows **不是同一件事**，而仪器把 Windows 那一件写死在用例里。

## 3. 改后：同一枚仪器，两平台都有世界

拆法只有一处：把"怎么让一个拼写变成不可证"（`pathForms.partial` 的来源）按平台供给，叫 `partialCandidate334`。

| 文件 | 行 | 内容 |
|---|---|---|
| `anchor_miss_reason_334_test.go` | `:78` `:103` `:141` `:157` | 三枚用例与那枚 class 用例改叫 maker；**`//go:build windows` 摘掉**（`numstat` 29/31，全是注释重写 + 取候选那几行） |
| `anchor_miss_reason_334_windows_test.go` | `:27-38`（新增 38 行） | Windows 半：`hostile334`（U+001B/U+202E，OS 拒绝解析）+ **前提断言**（`errors.Is(err, os.ErrNotExist)` 就必须 fatal，否则 partial 会静默变 false） |
| `anchor_miss_reason_339_other_test.go` | `:60-69`（+11 行，纯追加） | POSIX 半：**委派**给 347 落地的 `posixPartialCandidate347`（悬空符号链接），⛔ 重造一遍世界 |

其余（假主目录世界 + 装有/缺位叶子清单）仍旧只走 `anchor_world_334_test.go` 的共享 builder ⇒ 两平台不会漂成两份期望（347 的合件条件保名 `candidateOffAnchorChain334`，本腿零触碰）。

**没有放宽任何断言**（可机读的尺）：

- `git diff -U0 -- internal/risk/ | grep '^-' | grep -E 'want A|ClassA|t\.Skip'` => **0 命中**（那三句承重断言逐字在场）
- 用例数：Windows 名册 `22` 枚改动前后一字不变（不是删、不是加）；Linux `19=>24`（多出来的正是仪器自己）
- 新增 `t.Skip`：**0**；新增 `-short`/`testing.Short` 分支：**0**
- 三枚 A 档句子的 `d.Class != ClassA` / `d.NeedL2 || d.Allow` / override-key 三句**逐字未动**

## 4. 颜色（四数逐枚点名，⛔ 混用 ⛔ 相减）

名册（两平台同一发）：`-run 'TestAnchorMiss|TestAnchorFallbackTeeth'`，`-count=1 -timeout 30m -v`。
台面：Linux = `golang:1.27` 容器 + 仓外导出树；Windows = 本机 `go1.27.1` + 仓外导出树（同锚）。

| 台面 | RUN | PASS | FAIL | SKIP | rc |
|---|---|---|---|---|---|
| Linux 改前（`linux-before`） | 19 | 3 | 0 | 0 | 0 |
| Linux 改后（`linux-after`） | 24 | 5 | 0 | 0 | 0 |
| Windows 改前（`win-before`） | 22 | 3 | 0 | 0 | 0 |
| Windows 改后（`win-after`） | 22 | 3 | 0 | 0 | 0 |

- Linux 那 +5 = 2 枚顶层 + 3 枚子用例（`absent-aws-leaf-accuses-aws-and-says-so`／`installing-aws-moves-the-named-anchor-to-kube`／`proven-match-keeps-the-bare-literal`）＝`339-r1` 甲形在 Linux 上让出的那 5 枚名字，现在回来了。
- 日志：`logs/linux-base-miss-suite.log`／`logs/linux-green-miss-suite.log`／`logs/win-base-miss-suite.log`／`logs/win-green-miss-suite.log`（逐枚非 0 字节）。
- **收尾尺（从已提交树重打，不是工作树）**：`git archive dfd6f624b` 重建两平台台件 => Linux `24/5/0/0 rc=0`、Windows `22/3/0/0 rc=0`（`logs/linux-final-miss-suite.log`／`logs/win-final-miss-suite.log`）＝读数属于提交字节。
- POSIX 侧红句真长得出来（不是空跑）：`logs/linux-green-miss-suite.log` 里 `anchor_miss_reason_334_test.go:97: honest reject: A-list sensitive path (non-overridable): ~/.aws/credentials [fail-closed, NOT a match on this rule]: the candidate "/tmp/gt/.../profile/work/id_ed25519" is not proven handle-expanded (its own walk stopped at "/tmp/gt/.../profile/work") and the anchor ".../.aws/credentials" is not installed here (its expansion stopped at "/tmp/gt/.../profile") ...`

### 整包 `./internal/risk/`（同锚、同台面，`-v -count=1`）

| 台面 | RUN | PASS | FAIL | SKIP | rc | 红名册 |
|---|---|---|---|---|---|---|
| Linux 改前 | 250 | 150 | 1 | 0 | 1 | `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（存量 `A380`） |
| Linux 改后 | 255 | 152 | 1 | 0 | 1 | 同上（逐名同） |
| Windows 改前（导出树） | 281 | 166 | 5 | 1 | 1 | C26 + 4 枚 `TestSync*Ancestor*/116`（导出树台面） |
| Windows 改后（导出树） | 281 | 166 | 5 | 1 | 1 | **逐名 IDENTICAL**（`diff` 空） |

- Windows 那 4 枚 `116` 红与 1 枚 `TestSyncRegistryProbeLive` SKIP 在**改前改后同数同名**＝导出树台面（仓外无 `.git`）的形状，不是本腿；共享工作树里的权威色见第 6 节（只有 C26 一枚红）。
- 日志：`logs/linux-fullpkg-{before,after}.log`、`logs/win-fullpkg-{before,after}.log`。

## 5. `AC#2` 的去牙反控（两平台成对，红句逐字）

突变＝把 `internal/risk/blacklist.go` 的 `aRelates` 尾行回落支**整支摘掉**（同 `334-v1` 的 M-V1 形）：内容锚 `\treturn uncertainAnchorMiss(cand, anchor)` 命中 **1**，替换成 `\treturn false`；突变只打在仓外导出树的物理字节上，逐发 `git show HEAD:...` 还原并 `hash-object` 自证 `restored=44f6eb9b508eeaa0943cdeb9123bd19ad993bc78 match=YES`（两平台各一次）。

| 台面 | RUN | PASS | FAIL | SKIP | rc | 红名册 |
|---|---|---|---|---|---|---|
| Linux **改前**树 | 19 | 0 | 3 | 0 | 1 | 336 牙 + 339 前提孪生两枚；**票 334 仪器的两枚名字不在名册里**（摘牙对它＝零变化） |
| Linux **改后**树 | 24 | 0 | 5 | 0 | 1 | + `TestAnchorMissReasonNamesCandidateAndMissingAnchor` + `TestAnchorMissKeepsNonOverridableClassAndNeedL2` |
| Windows 改后树 | 22 | 0 | 3 | 0 | 1 | 334 两枚 + 336 |

Linux 改后树的红句（逐字，`logs/linux-deteeth.log`）：

```
--- FAIL: TestAnchorMissReasonNamesCandidateAndMissingAnchor
--- FAIL: TestAnchorMissKeepsNonOverridableClassAndNeedL2
    anchor_miss_reason_334_test.go:81:  fixture guard: the fail-closed shape must still be A here, got class=B allow=false needL2=true reason="B-list sensitive path; single-file override requires one L2 confirm: id_*"
    anchor_miss_reason_334_test.go:106: fixture guard: want A from the fail-closed clause, got class=B ...
    anchor_miss_reason_334_test.go:143: class: got B, want A (the fail-closed tier is unchanged by AC#1)
    anchor_miss_reason_334_test.go:147: Gate: got class=B ... want class=A needL2=false allow=false (non-overridable)
    anchor_miss_reason_334_test.go:151: an A deny from the fail-closed clause was unlocked by an override key: class=B allow=true ...
```

=> 票面那句"修好后的仪器必须在**两平台**都点名红"现在成立，且**是同一枚仪器**（同一批用例名），不是"Linux 上另养一枚更弱的孪生"。
=> 作差逐名归位：Linux 改前 3 红里没有票 334 的名字，改后 5 红里有两枚＝"补世界"补回的就是这两枚。
=> 正控：同树未突变 ⇒ 两平台全绿（第 4 节四数）。

## 6. 门禁四数（共享工作树，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 门 | 命令 | rc |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **0** |
| test | `go test ./internal/risk/ -count=1 -timeout 240s` | **1**（存量 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`，见下） |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l <写面三件>` | **0**，输出 **0 行** |
| d22scan | `sh scripts/d22scan.sh` | **0**（clean；ban#8 internal 885 枚、cmd 207 枚） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **1**（既有 roster 漂移，见下） |
| vet | `go vet ./internal/risk/` | **0** |

两枚非 0 的归因（⛔ 拿"反正不是我的"顶）：

1. **`TestC26RewriteAccount...`**：在 **HEAD 导出树**（未含本腿任何字节）里复跑同样红，红句点名 `cmd/wisp/panel_pump.go:283/:291` 与 `internal/panel/git_readonly_exec.go:409/:412`（"reads Result.Canonical without reading the rewrite account"）＝347 结案残余第 1 条（账 `A380`，panel 禁写面）。⇒ 存量，非本腿。
2. **长度门 rc=1**：`check-path-length-budget.sh:328`（**HEAD 里就是这一行**）点名的 roster 条目是 `179-...-sentence.md`，而仓里被跟踪的件已是 `179-...-sentence-done.md` ⇒ 票 179 改名 `-done` 之后 roster 没跟上＝既有漂移，与本腿三件（最长 38 字符名）无关。

## 7. 边界／未做（具名）

- ⛔ 翻任何框（票面至今 `未勾 1／已勾 6`，尺在 `grep -cE`）；⛔ 改票名／`Status:`；⛔ 加 `-done`。
- ⛔ 动票 334 的判据文字（本腿只动该 .go 件里的注释 + tag + 取候选那几行，票 334 的 issues 件零字节）。
- ⛔ 动产码：`blacklist.go` / `pathresolver*.go` 在共享工作树里 `git diff --name-only` **零命中**；去牙突变只在仓外树、逐发还原。
- ⛔ `./cmd/wisp` 整包测试、⛔ `go test ./...`（本腿零次）。
- ⛔ 新导出名、⛔ `thresholds.go`/golden/`allowlist.txt`/`go.mod`/`go.sum`/`.gitignore`/`docs/specs/**`/`docs/PLAN.md`/`frontend/**`/`design/**`。
- ⚠ 本腿把 probes 落在 `.scratch/wisp/probes/339/r2/**`（票面点名的是 `r1`＝上一腿的目录）；这是命名空间、不是扩面。
- ⚠ 一枚已知的注释漂移（未动、具名上报）：`anchor_world_334_test.go:45-48` 那句把 `hostile334` 的读者写成"336 牙 + 334 件"，现在 334 件的读者换成了它的 `_windows` 半（同族、仍带档，`var _ = hostile334` 的用途不变）。
