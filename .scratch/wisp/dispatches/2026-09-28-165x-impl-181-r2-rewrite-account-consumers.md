# 派单：`181-r2`（写码腿·小面）——三处新消费者没读"改写账户"：全仓那把尺已经为此变红，修完必须转绿

- 派单时刻：`2026-09-28 16:5x`｜编排者锚点 HEAD＝**起手自取**（`git log -1 --format='%h'`，派单时盘上是 `e9c216b8`；⚠ **别照抄我的号，本轮已第 5 次锚点漂**）
- 工单：`.scratch/wisp/issues/181-…-git…md`（票 181 的 AC#6 欠款）＋根因规矩出处＝**票 102**（C26 展开会把路径**改写到另一棵树**并照样报成功）
- 性质：**缺陷修复腿（小而准）**｜⚠ **AC 框一枚不许勾**

## 0. 为什么现在派（一句话）

`33-r1` 交件时具名报了一枚**不在它名册里的红**（它没碰 `internal/risk/**`，判其对错交回给我），我本程现跑复现并归因完成：

```
go test -count=1 -run TestC26RewriteAccountIsConsumedAtEverySecurityLeg ./internal/risk/
--- FAIL: TestC26RewriteAccountIsConsumedAtEverySecurityLeg
  pathresolver_rewrite_account_test.go:200: ..\cmd\wisp\panel_pump.go:93  reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
  pathresolver_rewrite_account_test.go:200: ..\cmd\wisp\panel_pump.go:104 reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
  pathresolver_rewrite_account_test.go:200: ..\internal\panel\git.go:132  reads Result.Canonical without reading the rewrite account (Actable/Rewritten) - ticket 102's hole re-opened by a new consumer
```

⇒ **这不是那把尺坏了，是它今天抓到活东西**：`181-r1` 新增的三处消费者（面板快照那两行＋git 读面那一行）**只取 `Canonical`、不看"这条路径有没有被 C26 改写过"**。票 102 的原始洞就是"改写之后仍按原树行事"，而这枚判据存在的意义正是**"每个新消费者都必须显式消费改写账户"**。⚠ 我上一轮收 `181-r1` 时跑了 `panel`／`config`／`tools`／`d22scan`／`gate-clauses`，**唯独没跑 `./internal/risk/`** ⇒ 按包门禁结构性看不见这条全仓不变式，**这一格漏记在我账上**（台账 `A380`）。

## 1. 写面（只这些，逐枚具名）

✅ `internal/panel/git.go`｜✅ `cmd/wisp/panel_pump.go`｜✅ 这两枚各自的判据件（`internal/panel/git_test.go`／新建 `cmd/wisp/` 名下件——⚠ 若你判"必须新建 `cmd/wisp/*_test.go`"，先读 §4 的票 98 那一格再决定）｜✅ 新建 `docs/evidence/s1/181-rewrite-account-r2.md`｜✅ 新建 `.scratch/wisp/probes/181/r2/**`｜✅ 追加票 181 一段 Progress log（**只追加、不改原句**）。
⛔ **`internal/risk/**` 一字不许动**——那枚判据（`pathresolver_rewrite_account_test.go:200`）是**裁判**，不是被告。⚠⚠ **最诱人的烂修法就是"把三行加进白名单／把扫描范围收窄"**，那等于把票 102 的洞焊死成规矩：**判"必须改那枚判据才能绿"＝停手上报**。
⛔ `internal/panel/composer_dispatch.go`／`composer_dispatch_test.go`（`33-r1` 刚落地的入向那一跳）不动；⛔ `bridge.go`（白名单起手终态都要现量 `grep -cE '=\s*"panel\.'` ＝ **4**）；⛔ `docs/PLAN.md`、`docs/specs/**`、`go.mod`／`go.sum`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量、`internal/panel/tokens_fourway_test.go`／`composer_test.go`／`l2_grant_boundary_test.go` **一字不动**。
⛔ `frontend/**`／`design/**` 不读不写不引不转述（别家地界）。
⛔ 零删除命令；不许跑 `probes/161/r6/flip-declaration.sh`；不许改 `scripts/d22scan.sh`／`probes/154/gate-clauses.sh`；只 commit **不 push**；显式 pathspec；禁 `add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`／`switch`／`merge`。
⚠ **写面闸门**：起手跑 `git status --porcelain -- internal/ cmd/` 把名册原样贴进表 §1；**终态必须等于起手名册（逐枚具名差集为空）**——⚠ **不是"必须为空"**（同树有别家程在飞是常态）。
⚠ **变异载具的两枚仪器坑（`33-r1` 用 39 枚调用换来的，别再交一遍学费）**：① 载具脚本里那个 `ROOT` 我上次数错了层数（4 层 vs 5 层），**第一发什么都没变异、而还原检查还打印 "IDENTICAL"——因为两枚哈希都是空**；⇒ **载具必须先自证"变异真的落到了那个字节"（比对变异前后哈希，且**两枚都必须非空**）**；② `grep '^^- [ ]'` 里的 `[ ]` 是**字符类**，永远数不出未勾框，要写 `'^- \[ \]'`。

## 2. 必答四格（先测→再写→再提交；引代码要 `file:line` 逐字）

**AC#α 先把"改写账户到底是什么"现读清楚**（别按名字猜）：`internal/risk` 里那个返回结构（`Result`）有哪几枚字段表达"被改写过／这条答案可用"（尺：`grep -rn "Actable\|Rewritten" --include=*.go internal/risk/ | grep -v _test.go` 现跑）；**既有消费者是怎么消费它的**——挑两枚**已经合规**的最近邻逐字抄它的形状（尺：`grep -rn "\.Rewritten\b\|\.Actable\b" --include=*.go internal/ cmd/ | grep -v _test.go`，取非本程写面里的那几枚）。
**AC#β 三处逐枚修**：`cmd/wisp/panel_pump.go:93`／`:104`、`internal/panel/git.go:132`。每枚都要答**"被改写过的时候，面板这一格该显示什么"**——⚠ **三档候选**：① 显示改写后的真值＋在 `reason` 里明说"这条路径被 C26 改写过"；② 整维降级成 `unreadable`＋一句原因；③ 照旧显示 `Canonical` 不加说明（＝今天的病）。**先给三档各自的后果，再由我裁**（不许自选，这是票 102 那条口径的延伸）。
**AC#γ 判据要长在"新消费者还会再犯"那一形上**，不是只钉这三行：⚠ 现读那枚全仓判据（`internal/risk/pathresolver_rewrite_account_test.go`）的**扫描形状**（它怎么找消费者、按什么算"读了 Canonical"），然后回答：**"我修完之后，如果第四枚消费者再犯，它会红吗？"**——给一发可重跑的尺（例：临时造一枚读 `Canonical` 不读账户的假消费者，看那枚判据响不响；**假消费者只准建在 `-overlay` 或仓外，⛔ 不许落进跟踪件**）。
**AC#δ 别把这条规矩渗到不该渗的地方**：`git.go` 那一行取的是 `WorkspaceView.Canonical`（**工作区视图**，不是某次 `fs.read` 的解析结果）——⚠ 现量答：**工作区视图这条路上到底有没有"改写"这回事**（有 ⇒ 修；没有 ⇒ 那这枚判据为什么把它算成消费者？具名报，别硬修，也别硬绕）。

## 3. 门禁（终态读数取在你自己最后一枚 commit 之后；⚠ 按 `A363` 不充当结案凭据）

⚠ **本票的结案凭据之一就是那枚红转绿**：`go test -count=1 ./internal/risk/`（**必须单包跑**，四包并发假红＝`A359`）⇒ 目标 **`ok`**；逐包 `go test -count=1 ./internal/panel/ ./cmd/wisp/ ./internal/config/ ./internal/tools/`（⚠ `./cmd/wisp/` 本机缺 DLL 会 `0xc0000135`＝票 98 ⇒ 要跑必须带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，**两形读数都要贴**：带的和不带的）；`sh scripts/d22scan.sh`（⚠ **基线 `ban #8 internal/`＝438**，不是 433/436；`33-r1` 落了两枚件）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（在册只 `G6neg`）；`gofumpt -l` 名下件空。
⚠ **`./internal/panel/` 今天有 3 枚在册红**（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）⇒ 照实记、不许当绿、不许修（前两枚归别家地界）。

## 4. 表骨架（第 ≤5 枚调用内落盘并 commit 第一枚）

① 起手锚＋起手脏件名册 ② AC#α 字段语义与两枚已合规最近邻（逐字） ③ AC#β 三处逐枚：改前三行／改后三行／三档后果 ④ AC#γ "第四枚再犯会不会红"的尺与读数 ⑤ AC#δ 工作区视图那条路有没有改写这回事 ⑥ **那枚红转绿的逐字读数**（改前 `--- FAIL` 三行＋改后 `ok`） ⑦ 门禁全套（含 PATH 两形） ⑧ 变异自证（含"两枚哈希非空"那条自证） ⑨ 本程没测的（逐名） ⑩ 被拒调用＋零删除自证＋工具调用终值 ⑪ next＝还有谁在消费 `Canonical` 没被这把尺覆盖到

## 5. 硬顶与计数

工具调用**硬顶 30**；**每答完一格自报累计**；**第 22 枚起不许开新探索**；判"必须动 `internal/risk/**` 那枚判据才能绿"＝**停手上报**（那是裁判）。

## 6. 结束消息回我六项

① 三处逐枚改前／改后（`file:line`＋三行 diff 摘要）；② 我裁那一格要看的**三档后果**（你倾向哪档、为什么）；③ AC#γ 那发尺的读数（第四枚假消费者响不响）；④ AC#δ 的结论（工作区视图有没有"改写"这回事，那枚判据算不算误报）；⑤ **`./internal/risk/` 从红转绿的逐字**＋门禁全套＋`bridge.go` 枚数两向＝4；⑥ 起手/终态脏件名册差集＋被拒调用＋零删除自证＋工具调用终值＋你没测的逐名。
