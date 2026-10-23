# 票 348 续跑腿 `348-r2` — 产码接缝落地＋Windows `AC#2` 成对三件（实现者交件）

**身份**：写腿（实现者）。⛔ 翻框、⛔ 改票名／`Status:`／`-done`、⛔ 改台账、⛔ push。
**起手锚**＝`51bb042b`（我第一笔 `git rev-parse HEAD` 现认；派单存档提交 `76ae3bbb` 之后编排者又落了 `51bb042b`＝台账＋票 164 结案，⛔ 碰）。
**写面**＝`internal/winsec/**` ＋ 本件。工作树里 `internal/winsec` 起手 `git status --porcelain` ＝ **0 条**（无别人未提交改动）；`docs/evidence/s1/362-r1-impl.md` 的 `M` 是**另一条腿**的，⛔ 我一只手没碰。

## 0. 机主批准（★逐字）

> 「好的，这三件事我都允了，同意，你可以继续了」

＝派单存档 `.scratch/wisp/probes/orch/20261023/txt/dispatch-three-legs-after-approval.md` 第 ① 节所引、批准范围第 ① 项即为本腿：「票 348 加"删除那一跳"的产码接缝（**第二道人工批准已获**）」。★这只解冻 **`internal/winsec` 产码里那一处最小面**；⛔ 借此重构、⛔ 动 `internal/{tools,config,session,memory,risk,panel,agent}/**`、⛔ 动 `cmd/**`／`scripts/**`／`.github/**`／`tools/d22scan/**`／`go.mod`／`go.sum`／`docs/{PLAN.md,specs/**,reports/**}`／`frontend/**`／`design/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件。核：本腿改动面 ＝ `internal/winsec/winsec.go`（`git diff --numstat` ＝ `23 2`）＋ `internal/winsec/deleteseam_348r2_windows_test.go`（新）＋ 本件。

## 1. 接缝加了什么（逐行）

**`internal/winsec/winsec.go`**（`git diff --numstat` ＝ `23 2`＝+23/−2）：

| 行 | 内容 |
|---|---|
| 新增 `var removeRefusedLeaf = unlinkPinnedLeaf` （1 行代码＋注释块） | 承接「删除那一跳」。签名＝`unlinkPinnedLeaf` 原有的 `(leafPin, string, fs.FileInfo)`；默认值＝**落地自己的删除**；⛔ 新导出名、⛔ 新参数、⛔ 改任何既有函数形状、⛔ 改 `privateFile` 返回语义 |
| `:186` 密封失败支 | `unlinkPinnedLeaf(parent, filepath.Base(resolved), pinned)` → **`removeRefusedLeaf(...)`** |
| `:191` 写失败支 | 同上 |

★**为什么是 `var` 而不是加参数**（同 `applyDescriptor` 的理）：测试替换它、**在那一个瞬间动手**、然后把原件**原样调回来** ⇒ 读数量的是**这条路由的删除**，⛔ 是"测试自己删文件"。★反向判据读同一枚 `var` 的另一个方向。

**`internal/winsec/deleteseam_348r2_windows_test.go`**（`wc -l`＝**151 行**，`//go:build windows`）：一枚用例，把接缝换成「落祖先翻转（`reparseAt` 现证）→ 调本树原件 → 复原」，跑 400 发，门槛＝`deleted == 0` ＋ `seamCalls == attempts` ＋ **`hopLanded == seamCalls`**（liveness 地板：零只有在"动手那一瞬间祖先真的站在链接上"时才算证据；`AC#2` 逐字「⛔ 拿 0 命中当凭据」）。

## 2. ★Windows ① 加缝前 vs 加缝后

**台件** ＝ 仓外导出树（`git archive HEAD` ＋ 我那两枚文件；单差异＝`internal/winsec/winsec.go`，见每格）。`unset GOTMPDIR`、`GOFLAGS=`、CGO 未设 0。

### 2.1 加缝**前**（T1 ＝ `winsec.go` 与 git 基线 `fde1df02` **逐字节相同**，唯 `wrapPath` 的 `%v`→`%w` 是后来另一票的行；**无接缝**）

**★这一格推翻派单的预期，逐字报**：派单写「加缝之前…（**应仍取不到**）」。**取得到** —— 只是**不稳**，原因不在世界不可达、在"删除支几乎不跑"：

| 仪器 | 3 发读数 | 说明 |
|---|---|---|
| `348-r1` 那一形（赛跑，**不注入密封失败**） | `deleted=0 / 0 / 0`；`truncated=11 / 7 / 11`；`provedReparse=71 / 73 / 68`；`formBLegRefusals=0` | ★这就是 `348-r1` 报的「①git 基线取不到」。但同一发 `truncated>0` 证窗口到达过；`deleted=0` 的真因＝**密封没失败⇒删除支没跑**（不是删除安全） |
| 同一棵树 ＋ 用**本包既有的 `applyDescriptor` 缝**注入密封失败（⇒ 删除支 129/138/129 发真跑） | **`deleted=10 / 8 / 8`**（of 400）；`provedReparse=75 / 82 / 78`；`truncated=0 / 0 / 0` | ★**①在 git 基线字节上是取得到的**（≈2%/发、占跑到删除支的那 6–8%），但**随调度器摆动**。⇒ 接缝买的不是"可达"，是"**确定**" |

### 2.2 加缝**后**（T2 ＝ 落地 ＋ 接缝 ＋ **路由回到修前形**：⛔ 钉 ＋ 删除体＝`os.Remove(filepath.Join(p.dir, base))`＝原 `os.Remove(resolved)`）

| 树 | `deleted` ×3 | `hopLanded` ×3 | `seamCalls` | 判 |
|---|---|---|---|---|
| **T2 修前形＋接缝** | **`400 / 400 / 400`** | `400 / 400 / 400` | 400 | **红**（①>0，确定性，3/3） |
| T3 落地＋接缝 | `0 / 0 / 0` | `400 / 400 / 400` | 400 | 绿（②=0，③同发在场） |
| T4 落地＋接缝，**只**把删除体换回修前 | **`400 / 400 / 400`** | `400 / 400 / 400` | 400 | **红** |

**⇒ ①`400/400`（加缝后，3 发逐字同号）⇄ ②`0/400`，③`hopLanded=400` 同一发现证祖先站在链接态。成对三件在 Windows 上齐了，且是确定性的。**

## 3. ★T4＝承重归因（本趟最硬的新东西，⛔ 派单要求的，但不说清就会读错）

T4 **保留钉**、**只**把删除体换成修前那一句 ⇒ `deleted=400/400`。⇒ ★在**删除那一跳的那一个瞬间**：
- `pinnedParent` 的钉**不**阻止这次删除被搬走（`deleteseam` 的翻转把 `sub` 换成**同一父目录内**的名字，测到 `err=<nil>`）；
- `unlinkPinnedLeaf` 的**身份复验也不**救它（身份复验在**被换掉的那个体**里；换成按名体就没这一步）；
- ⇒ **真正承重的是"fd-相对"这一半**：删除相对**被钉住的那个对象**，所以名字（`sub`）被人翻走以后，删除仍落在原对象里，落不到外来树。

★同时更正 `348-r1` 的成因 (a) 的**射程**（我复跑+扩测，`logs/03-*`，本机 1 发）：祖先 `rename` 的真实三态是**四态**——

| 态 | 读数 |
|---|---|
| H0 空手 | `rename(data)=<nil>` |
| H1a 叶子按名开着 | `rename(data)=Access is denied` |
| H1b 叶子按名开着，rename 叶子自己的目录 | `Access is denied` |
| **H1c 刚 `Close` 后**（修前路由唯一空手的瞬间） | **`rename(data)=<nil>`** |
| H2a 钉在 `sub` ＋ 叶子经 `Root` 开着 | `rename(data)=Access is denied` |
| H2b／H2d 钉在 `sub`，叶子开／关，把 `sub` 挪到**别的父目录** | `in use by another process` ×2 |
| **H2e 钉在 `sub`，叶子关着，把 `sub` rename 到**同一父目录内** | **`<nil>`** ★这就是接缝仪表能落世界的原因** |

⚠**两处我没解决的、照实标**：①`H2c`（钉开着、叶子关着、rename **祖祖**目录 `data`）＝`Access is denied`，而 `H2f`（同一姿势，在 H2e 那一轮之后）＝`<nil>` ⇒ **这两发互相矛盾、我⛔ 归因、也⛔ 拿它当判据**（本票的路不用它）；②H2e 这一条**改写了**我在 `winsec.go` 与测试头注释里的第一版措辞（原写"rename 任何祖先都被拒"）⇒ 已按测量改成"**挪出父目录**被拒、**父目录内换名**不拒"，两处注释逐字见产码。

## 4. POSIX 三读数（★**复用**，⛔ 我重跑；标清来路）

★本腿**没有**重跑 POSIX 容器台件（派单允许「复用或重跑、要标清」）。逐条来路＝既有已入库交件，**不是**本腿现跑：

| 件 | 读数 | 谁量的／哪一发 |
|---|---|---|
| ①修前 >0 | 确定性红（A3 那枚：`unlink err=<nil>` ＋ 见证 `no such file`）＋ 竞态 `deleted=76/1500`、`62/1500` | `348-r1`（派单第 ① 节逐字引）；`342-v2` 独立复量 `86/71/73/73` of 1500 |
| ②修后 0 | `0/1500` | `348-r1`；`342-v2` 独立复量 `0` ×4 |
| ③同发证窗口到达链接态 | `provedLinkCycles=267/256`、`pinnedWalkRefusals=50/54`、`pinnedParentBeforeOpen=29/39` | `348-r1`；`342-v2` 独立复量 `provedLinkCycles=159/152/147/156`、`pinnedParentBeforeOpen=22/26/25/15`、`pinnedWalkRefusals=28/32/43/26` |

★本腿**自己**为 POSIX 做的只有**编译门**：`GOOS=linux` `go build`＋`go vet` 各 rc=0、`GOOS=darwin` 同（接缝是**无平台标签**的 `winsec.go` 里一行，三枚 `unlinkPinnedLeaf` 签名一致 ⇒ 编译证接缝没把 POSIX 打坏，⛔ 证行为）。⇒ **⛔ 读成"本腿复跑了 POSIX"**。

## 5. ★两平台各一发（⛔ 一发代两发）

- **Windows**：本腿**新跑**（§2.2 那 3 发＋§2.1 那 2 组），台件＝仓外导出树，逐个日志在 `D:\tmp\wisp\20261023\348r2\logs\`（`10-*`…`40-*`）。
- **POSIX**：**复用** §4（`348-r1` ＋ `342-v2` 两枚独立台件），⛔ 本腿新跑。
⇒ ★两平台**各有自己的一发**，但**来路不同**（一新一复用），任何"本腿两平台都现跑"的读法都不成立。

## 6. 反向判据 ＋ 还原核

**反向判据（同一发、只差删除体）**：T4 ⇒ 必须红，逐名红行（逐字，3 发同号）：

```
--- FAIL: TestAC2SeamWindowsTheRefusalPathRemovalCannotBeMovedByAnAncestorFlip348r2
deleteseam_348r2_windows_test.go:118: AC#2' first damage verbatim: witness ABSENT after attempt 0
  (GetFileAttributesEx …\someone-elses-tree-w\leaf-342r1w.bin: The system cannot find the file specified.)
deleteseam_348r2_windows_test.go:132: AC#2 RED: 400 of 400 attempts found the foreign witness GONE
  through the refusal-path removal (first: witness ABSENT after attempt 0 (…)), so that removal was
  moved by an ancestor flip and deleted a name in somebody else's tree
```

★**为什么这枚红行咬得住**：同一发里 `hopLanded=400`（世界到达）、`seamCalls=400`（仪表确在那一跳上）、`deleted=400`（见证被删）——三个数**同发**，⛔ 三段拼接。T2（修前整条路由）同号红；T3（落地）同号绿。

**还原核**：★核 hash 前先答「谁算的／什么算法／几位」——
- **`git hash-object`** ＝ **git** 算的、**SHA-1**、**40 位十六进制**；
- **`sha256sum`** ＝ **sha256sum** 算的、**SHA-256**、**64 位十六进制**。**⛔ 两者混比**（`ruler-scope-checklist` 第 9 例）。

| 文件 | SHA-1（`git hash-object`，40 位） | SHA-256（`sha256sum`，64 位） |
|---|---|---|
| `internal/winsec/winsec.go` | `fc4e0d696ac482951902a0b82dcf9429284bc33b` | `6bed1b5bfae9bf37f7b5312e7040a290a874246406f5eb4756a1b42e75a6e52e` |
| `internal/winsec/deleteseam_348r2_windows_test.go` | `3272fdace3b8ffb21ea77768c62332bc9f7d449b` | `41d1121d3dd41759afc1b1cbdafc4d9c50999d9ff55266130f20ab5e11eeb88f` |

★**还原的实质**＝**工作树从来没有被突变过**：全部修前形/反向形都打在**仓外**三棵导出树（T1/T2/T4）上，⛔ `-overlay` 之外一次工作树改动都没有，⇒ 无"突变窗口"、窗口期零提交。上表两枚 hash 在全部变异运行**之后**复算，与运行**前**同值（同算法同位数逐位比）。`git status --porcelain -- internal/winsec` 收件时＝ `M winsec.go` ＋ `?? deleteseam_348r2_windows_test.go` 两条，无第三枚。★**注**：`winsec.go` 的 hash 与 `HEAD:internal/winsec/winsec.go`＝`757b16a8…` 不同是**应然**（接缝就是本腿的改动），⛔ 读成"没落地"。

## 7. 格式：甲乙两尺都报

| 尺 | 对象 | 读数 | 判 |
|---|---|---|---|
| **甲**（判据） | 仓外导出树 T3（＝`git archive HEAD` ＋ 我那两枚文件 ⇒ **与提交后仓外导出等价**）跑 `/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec` | **0 行**，rc=0 | **GREEN** |
| **乙**（只诊断） | 共享**工作树**跑同一命令 | **0 行**，rc=0 | 诊断，⛔ 判据 |

归属判据用 `[\\/]`（Windows 上 `gofumpt` 输出反斜杠，`A1140⑥`）：甲那次 `-match '[\\/]'` 命中 **0**（＝本来就没列任何文件，⛔ 是"列了但被过滤掉"）。`gofumpt` 一律**绝对路径** `/d/work/base/gopath/bin/gofumpt.exe`（`Test-Path` ＝ True）。`tr -cd "\r" | wc -c` 两枚 `.go` 各＝**0**（LF 干净，⛔ CRLF）。

## 8. 门禁 rc（串行；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `go build ./...` | **0** | — |
| vet（本机） | `go vet ./internal/winsec/` | **0** | — |
| vet（POSIX 编译） | `GOOS=linux` / `GOOS=darwin` 各 `go build`＋`go vet` | **0** / **0** | 四发全 0 |
| gofumpt | 甲／乙见 §7 | **0** / **0** | 0 行 |
| pathlen | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`、`tracked paths=17609`、`over-budget=61`、`not in roster=0` |
| **d22scan** | `sh scripts/d22scan.sh` | **1** | ★**红在写面之外**：正控模块 `PASS=33 FAIL=2`（`TestScannerSelfScanOfRealRepoIsGreen`／`TestRealRepoLedgerIsHonest`，后者 ⛔ 打原因）；真扫 1 条 finding＝`internal/memory/schema.go:26: [phantom-citation] … "docs/evidence/s1/196-r1-impl.md" which does not exist on disk`（**票 196 的腿**在飞）。★**`grep -c winsec` 该日志＝0** ⇒ 本腿两枚文件 **0 条 finding**。⛔ 当绿、⛔ 当我的红 |
| 整包 `./internal/winsec/` | `go test -count=1 -v`（3 发） | **1** | 四数逐发：`RUN=123 PASS=121 FAIL=2 SKIP=0`／`123/120/3/0`／`123/120/3/0`（和＝RUN，逐发对得上）。红名逐枚：`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（2/3 发，子测试 `junction`／`symlink-dir`）、`TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`（1/3 发）。★**全部是票 342 的 `truncated` 残留**（`truncated=1…3` of 400，正是 `pinnedleaf_342_windows.go` 与 `winsec.go:118-140` 自陈"窗口未关"的那一味），**⛔ 票 348 的 `deleted`**——同一批 3 发里 `deleted=0` 逐发（`ancestorunlink_348_windows_test.go:94/:112`）。⇒ **不是本票新造的红**（`AC#5` 的归因口径：安静复量 3 发后归因）。本腿新钉 3 发全绿 |

## 9. 不声称的（★照实）

1. ⛔ 声称"Windows ① 只能靠接缝拿到"——§2.1 测到它在 git 基线字节上**本来就取得到**（`10/8/8` of 400），接缝买的是**确定性**（`400/400`）与**可归因**（T4）。
2. ⛔ 声称 `deleteseam` 那一枚覆盖 `AC#2` 的 ③**在落地侧**——落地侧的零是"世界被进入且见证没被删"（`hopLanded=400`），**不**再借票 342 那枚赛跑仪器的 `provedReparse` 当凭据（两处都有、⛔ 互借）。
3. ⛔ 声称 POSIX 是本腿现跑（§4 逐条标了来路）。
4. ⛔ 声称 d22scan 绿（§8 红在别人写面，我给了 `winsec` 命中 0 的读数）。
5. ⛔ 声称 `H2c`／`H2f` 那一对矛盾已归因（§3 照实标"⛔ 归因、⛔ 当判据"）。
6. ⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`／台账、⛔ push、⛔ `cmd/wisp` 整包、⛔ `go test ./...`。
7. ⛔ 声称 `AC#2` 已满足 —— **判语归非实现者**。本件只交读数：①`400/400`（T2，修前形＋接缝）／②`0/400`（T3，落地）／③`hopLanded=400` 同发（两树皆是）／④反向判据逐名红行（§6）／⑤两平台来路不同（§5）。

## 10. ★尺射程自纠（本腿自己犯的第 1 处，`A1127` 那一族）

★**我第一版把本件报成"102 行"、把测试件报成"140 行"——两个都是错数。**
- **错在哪**：我用 `Get-Content | Measure-Object -Line`。这把尺**丢空行**：本件 `102`（非空）＋`42`（空行）＝`144`；测试件同理 `151`。
- **§ 落笔前五问第 ④ 问没问**："这个数有没有和总数对上" —— 我一枚都没核。
- **抓出来的地方**：`git diff --cached --numstat` 报 `144 0` ⇄ 我写的 `102`，两个数对不上 ⇒ 我回头用 `wc -l`（**git 与 wc 同口径：按 `\n` 数，含空行**）重取：**本件 144 行**、测试件 **151 行**。
- **射程**：本件帽 ≤160 ⇒ **144 ≤ 160，未破帽**，⛔ 需要拆 `-partN`。
- ★**未更正的残留**：提交 `5ffa553e3` 的**提交正文里**仍写着"102 行 ≤160 帽"与"140 行"（⛔ amend／⛔ 改写已落的历史）⇒ **以本件 §1/§10 的 `wc -l` 读数为准**，提交正文那两个数按本行更正。

