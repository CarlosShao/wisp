# C26 去红腿 185（台账 `A380` 族）— `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`：RED 到手、GREEN 交件

- 起手锚＝`5a21be431`（`185轮:收wave11(220 STOP真机/342-AC4读数/284注释代提);A1546`）；本腿＝**单枚 `cmd/wisp` 写腿**＋`internal/panel`。
- 写面＝`cmd/wisp/**` ＋ `internal/panel/**` ＋ 本件（派单令）。**⛔ 动**：`internal/risk/**`（判据件，冻结）／`frontend/**`／`design/**`／任何 `C1–C32`、`D1–D47` 契约文字／`t.Skip`／任何断言强弱。
- 判据件＝`internal/risk/pathresolver_rewrite_account_test.go:167`（本腿一字未改，`git status` 为空）。

## 1. 起手 RED（逐字四行，命令可复跑）

```
$ cd "/d/work/workspace/projects plans/Wisp"
$ export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP
$ go test ./internal/risk -run "TestC26RewriteAccountIsConsumedAtEverySecurityLeg" -count=1 -timeout 420s
--- FAIL: TestC26RewriteAccountIsConsumedAtEverySecurityLeg (0.08s)
    pathresolver_rewrite_account_test.go:200: ..\..\cmd\wisp\panel_pump.go:283 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
    pathresolver_rewrite_account_test.go:200: ..\..\cmd\wisp\panel_pump.go:291 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
    pathresolver_rewrite_account_test.go:200: ..\..\internal\panel\git_readonly_exec.go:409 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
    pathresolver_rewrite_account_test.go:200: ..\..\internal\panel\git_readonly_exec.go:412 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
FAIL	github.com/CarlosShao/wisp/internal/risk	0.172s
FAIL
```

与派单点名的四行**逐字同号同句**（`cmd/wisp/panel_pump.go:283`／`:291`／`internal/panel/git_readonly_exec.go:409`／`:412`）。

## 2. 判据怎么读（先把仪器读透，再动码）

`pathresolver_rewrite_account_test.go:187-203` 的第二半是**文件级文本扫描**，不是 AST、不是行为：

1. `canonicalReaders` 从仓根 `WalkDir` 全部**非 `_test.go`** 的 `.go`，凡是某行出现 `.Canonical`（且其后一字符不接标识符，`readsCanonicalField`）就记一枚 `{file,line}`；跳过 `.git/.scratch/node_modules/dist/build/testdata/design/docs` 与 `tools/`。
2. 对每个命中的**文件**，只要**整份源码**里**既没有** `Actable(` **也没有** `.Rewritten`，就按每一条命中行报一条红（所以一个文件有 2 行命中就报 2 行）。

⇒ 这一半尺的语义是「**读 `Result.Canonical` 的文件，必须同时把改写账的名字写在自己身上**」。它看不见「读了账却当装饰」，所以本腿另配两枚行为钉（§4）补这一味。

**为什么这两枚文件不再匹配（逐行现读）**

- `cmd/wisp/panel_pump.go:283` 是 `fileTreeView` 的**文档注释**里那句 `// (.Canonical) instead of resolving ...`；`:291` 是代码 `return panel.ReadFileTree190(rt.workspaceView().Canonical)`。该件全篇无 `Actable(`／`.Rewritten`。
- `internal/panel/git_readonly_exec.go:409`／`:412` 是 `resolveGitTree` 里 `if strings.TrimSpace(res.Canonical) == ""` 与 `return res.Canonical, "", ""`；该件 `:397-399` 的注释当时**明说**「The rewrite account is deliberately NOT re-judged here」，所以整件也没有那两个名字。

**归属**：`panel_pump.go` 那两处＝票 190 的接线腿；`git_readonly_exec.go` 那两处＝票 189 的只读 git 腿。台账在册（`A380`；`docs/evidence/s1/107-verdict.md:157`、`342-ac4-readings.md:72`、`339-v2-verdict.md:61` 三处逐字同号）。⇒ 本腿不是「修裁判」，是「补被告」。

## 3. 修法（两处，都是把账真的读进来，不是加一句注释骗尺）

**① `internal/panel/git_readonly_exec.go` `resolveGitTree`：账走 `Result.Actable()` 这扇唯一门**

```go
res, err := risk.Resolve(root, nil)
if err != nil { ... }
// Ticket 102's door: the account is read, so an expansion that moved the
// spelling onto another tree becomes a refusal instead of a `-C` argument.
tree, err := res.Actable()
if err != nil {
    return "", GitKindUnreadable,
        "C26 的展开步把这条工作区路径改写到了另一棵树（rewritten=true），git 这一维不在原拼法上取数：" + err.Error()
}
if strings.TrimSpace(tree) == "" { ... }
return tree, "", ""
```

- 这一腿把答案交给**子进程**（`git -C <tree>`），所以不论子命令多只读，它都是「act」腿 ⇒ 按 `Result.Actable()` 的注释（`pathresolver.go:82-92`）必须过这扇门。
- **处置方式没被本腿改判**：票 181 AC#6 裁决的甲案（`GitViewRewritten`＝报改写后的真相并说明账户）照旧——它交给 `ReadGit` 的**已经是展开后的树**，再 `Resolve` 一次得 `Rewritten=false`，`Actable()` 放行，逐字不变。本门拒绝的是另一形：**一个还带着展开构造的 root 被没读账的调用者直接塞进来**。
- 文件里那句「deliberately NOT re-judged」按事实重写为「账在这里读，处置不在这里判」；`git_readonly_exec.go` 的既有仪器 `TestGitDimensionHasNoModelCallableTool`（读本件 AST，禁新增导出函数）复跑绿，见 §5。

**② `cmd/wisp/panel_pump.go` `fileTreeView` ＋ `internal/panel/tree_reader_190.go`：树维度改成「整枚 view 过门」，与 git 维度同形**

```go
// internal/panel/tree_reader_190.go（新增两枚，导出名一枚）
func ReadFileTreeForWorkspace190(ws WorkspaceView) TreeView {
    if ws.Rewritten { return fileTreeRewritten190(ws) }
    return ReadFileTree190(ws.Canonical)
}
func fileTreeRewritten190(ws WorkspaceView) TreeView {
    view := ReadFileTree190(ws.Canonical)
    view.Reason = "注意：这条工作区路径被 C26 的展开步改写过（账户 rewritten=true），" +
        "下面这份文件清单来自改写后的树 " + ws.Canonical + "，不是原拼法 " + ws.Spelling + "。" + view.Reason
    return view
}
```

```go
// cmd/wisp/panel_pump.go
func (rt *agentRuntime) fileTreeView() panel.TreeView {
    if rt == nil || rt.paths == nil { return panel.ReadFileTree190("") }
    return panel.ReadFileTreeForWorkspace190(rt.workspaceView())
}
```

- 形状与 `gitView`（`panel_pump.go:277` `panel.ReadGitForWorkspace(rt.workspaceView())`）**逐字同构**：坐标不再以裸串穿过，账随 view 一起走。
- 处置取**与 git 同一案**（甲案：报展开后的真相＋把账户说在句子里），理由写进代码注释：一个包不能对同一条工作区两种说法（git 段说 rewritten、files 段当「你的文件夹」）。乙案（整维丢弃）是编排者的裁决，本腿不自行取。
- 句子里带 `rewritten=true`、原拼法、展开后树三枚事实，与 `GitViewRewritten` 的同名事实对齐。
- 新增导出名 **`ReadFileTreeForWorkspace190`**（`internal/panel` 包内导出）。它不是 C17/入站方法名：`internal/panel/inbound_roster_253_test.go` 的闭环名册复跑**绿**（§5），`cmd/wisp` 侧零新增名。`ReadFileTree190` 本体未动（既有三个树仪器照旧读它）。

## 4. 牙齿：两枚新行为钉 ＋ 各一枚去牙（自己打，逐字）

判据尺看不见「读了账但没用」，所以每半补一枚**行为**钉，并各做一次去牙（mutation）证明钉有牙：

- `internal/panel/tree_190_account_consumed_test.go`（新）`TestReadFileTreeForWorkspaceConsumesRewriteAccount`：三子例＝未改写走原路（与裸门逐字同 kind/reason/total）／改写必须把账户说进句子（`rewritten=true`＋`Spelling`＋`Canonical` 三枚都在、且读完仍与裸门同维）／未设 view 不读。
  去牙＝把 `if ws.Rewritten {` 改成 `if false && ws.Rewritten {`，同一发：
  `--- FAIL: TestReadFileTreeForWorkspaceConsumesRewriteAccount/rewritten_must_say_so_out_loud`；还原后逐字回 `if ws.Rewritten {`。
- `internal/panel/git_readonly_account_consumed_test.go`（新）`TestResolveGitTreeConsumesRewriteAccount`：普通 root 放行且给出树／被改写的拼法（`%WISP185ACC_GITTREE%\sub`）必须**停在子进程之前**（kind 非空、`GitKindUnreadable`、理由含 `rewritten=true`）。
  去牙＝把 `tree, err := res.Actable()` 换成 `tree, err := res.Canonical, error(nil)`，同一发：
  `--- FAIL: .../rewritten_spelling_never_reaches_the_subprocess`，红句逐字 `resolveGitTree("%WISP185ACC_GITTREE%\\sub") accepted a rewritten spelling (tree="C:\\Users\\...\\001\\sub")`；还原后逐字回 `res.Actable()`。

## 5. GREEN 读数（四数按发分开报，不相减）

| 尺（命令） | 起手 | 交件 |
|---|---|---|
| `go test ./internal/risk -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg -count=1 -timeout 420s` | **FAIL**（§1 四行） | **PASS** `ok 0.334s`；`--- PASS: TestC26RewriteAccountIsConsumedAtEverySecurityLeg (0.26s)` |
| `go test ./internal/risk -count=1 -timeout 420s -v`（整包） | 台账在册＝`RUN=281 PASS=170 FAIL=1 SKIP=1 rc=1`（红＝C26） | `RUN=281 PASS=171 FAIL=0 SKIP=1` **rc=0**（SKIP＝`TestSyncRegistryProbeLive`，本机注册表现实） |
| `go test ./internal/panel -count=1 -timeout 900s -v`（整包） | 本腿未在开工前整包跑；HEAD 版三件用 `-overlay` 换回后，两枚在册红**在真实检出里 PASS**（§6） | `RUN=293 PASS=188 FAIL=2 SKIP=0` rc=1；红名**只有**两枚在册存量＝`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree`（都读 `frontend/**`、`design/**`，本腿零碰）⇒ **新增红 0 枚**；本腿两枚新钉与 `TestFullInboundMethodRosterIsClosed`／`TestTree190*` 全 PASS |
| 派单点名的两条命令 | —— | `go test ./cmd/wisp -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg -count=1 -timeout 420s`＝`ok [no tests to run]`；`go test ./internal/panel -run TestC26 -count=1`＝`ok [no tests to run]`。**照实报**：这枚判据件住在 `internal/risk`，两条 `-run` 在本包**匹配不到任何用例**（绿是「没跑」，不是「跑了绿」），真正的绿见上两行 |
| `go test ./cmd/wisp -run "TestRunBooksWithASnapshotOfItsLiveQueue\|TestSnapshotWorkspaceSectionReportsTheNarrowing\|TestThePumpIsDrivenNotJustAssembled\|TestPanel" -count=1 -timeout 900s` | —— | `ok 46.437s`（泵面；本腿唯一被改的生产文件就是泵） |

## 6. 门禁与两条环境噪声（逐条说明，⛔ 当绿）

- `gofumpt.exe -l` 对写面五枚 `.go`（三改＋两新）＝**空**，`rc=0`。
- `go build ./cmd/... ./internal/... ./tools/...`＝**rc=0**。
- `sh scripts/d22scan.sh`＝**rc=0 clean**（`runtests.sh: PASS=35 FAIL=0 SKIP=0`；live scope `bans#1-5 internal/=361 cmd/=57 · ban#6 116 · ban#7 48 · ban#8 design/=39 frontend/=116 internal/=895 cmd/=208`）。
- **噪声 A**：共享工作树里 `go build ./...` **rc=1**，报错全在 `build\mut_rvw.go`（`undefined: voiceTtsFiles/residentBall/...`）。该件**被 `.gitignore:14 build/` 忽略、未跟踪**，mtime `2026-11-04 10:26`（本腿开工于 11-05 16:38）⇒ 别腿留下的本机变异台架，⛔ 本腿产物；`git archive HEAD` 导出树（无该件）上 `go build ./...` **rc=0**。
- **噪声 B**：`internal/panel` 整包那两枚在册红与本腿无关的**证据**（不靠嘴）＝把本腿三枚生产件用 `go test -overlay` 换回 `HEAD` blob（并删掉本腿两枚新 `_test.go`）后，在**同一棵真实检出**里跑那两枚：
  `--- PASS: TestComposerRenderFixtureTellsTheTruth`／`--- PASS: TestReadGitOnThisRepositoryIsSelfConsistent`，`ok 0.675s`。
  它们只在 `git archive` 导出树里红（红句自陈 `kind = "not_a_repo" ... 这棵工作区目录及其以上都没有 .git`、`the render fixture has no block for ...`）＝**导出树形状**缺 `.git`／缺件，⛔ 代码红。

## 7. 边界、未动清单与 blob 号（裁决腿可逐枚复核）

- `git status --porcelain -- internal/risk` ＝ **空**（判据件一字未改）。`internal/risk` 整包本腿只**跑**不改。
- 本腿落地件（5 枚 `.go`）＋本件：

| 件 | 起手 blob（HEAD） | 交件 blob |
|---|---|---|
| `cmd/wisp/panel_pump.go` | `b13cc3db15bfd9308f3380ee629c149c8e0e8806` | `7e212b56ebacc88f01591e2e6d59b01c3b87acbf` |
| `internal/panel/git_readonly_exec.go` | `cd59a25c0b49fff51b5a4c8240074414aa35997d` | `5404778b342fa5bdb7450de269ed46f506aef774` |
| `internal/panel/tree_reader_190.go` | `1d025601d94b86964ec57b1716a140e822abf3f9` | `71625c93dd468529942dfb1c845fc4ed6a99894b` |
| `internal/panel/tree_190_account_consumed_test.go` | 无 | `ce5960d2361b02286419fba66b1c9b34b1236c76` |
| `internal/panel/git_readonly_account_consumed_test.go` | 无 | `251013ce941aafd4b3e77eab56398fbcf0b6fb56` |

- ⛔ 动：`internal/risk/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、任何契约件、任何 SLO/golden/`thresholds.go`；⛔ `t.Skip`、⛔ 放松断言、⛔ `git add -A`／amend／reset／rebase／stash（提交走显式 pathspec）；⛔ push。

## 8. 一句话终态 ＋ 撤销口令

**`A380` 那四行红的被告（`cmd/wisp/panel_pump.go:283/:291`、`internal/panel/git_readonly_exec.go:409/:412`）已从「读 `Result.Canonical` 不读账」改成「账随 view 走／账走 `Result.Actable()`」，判据件零改、断言零松；RED 四行逐字在案，GREEN＝判据用例 PASS＋`internal/risk` 整包 `281/171/0/1 rc=0`，`internal/panel` 整包只剩两枚与本腿无关的在册红（新增红 0 枚）。**

撤销口令（一句话复原本腿全部字节）：
`git revert b058171e7 --no-edit`（仅本腿五枚 `.go`＋本件；`internal/risk` 不在内，故回退后那四行红会逐字重现，属预期）。
★本件 sha 的追加说明：`b058171e7`＝产码那一笔（五枚 `.go`＋本件）。上面这行把占位符换成真 sha 的是**只改这一行的第二笔**（同一条显式 pathspec、零代码字节），回退口令仍指 `b058171e7`。
