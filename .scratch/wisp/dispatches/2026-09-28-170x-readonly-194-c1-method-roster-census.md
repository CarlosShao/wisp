# 派单：`194-c1`（只读代价普查·零产码）——规格那 18 枚方法逐枚核：有货／没货／撞禁区，再决定分几批落地

- 派单时刻：`2026-09-28 16:5x`｜编排者锚点 HEAD＝**起手自取**（`git log -1 --format='%h'`，派单时盘上是 `8ca2291f`；漂了登记即可）
- 工单：`.scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-code-to-the-spec.md`
- 性质：**只读普查** ⇒ **AC 框一枚不许勾**、**产码与 `go.mod` 零字节不动**、**不落地**。

## 0. 为什么派（一句话）

owner 09-28 16:5x 裁了 `Q-67`：「**方法名册，这个建议你按照当前规格补一下吧，别搞债务了**」⇒ 方向定了（**代码向 `docs/specs/SPEC-08-ui-ball-panel.md:163-174` 那张表对齐**）。但那张表里有几枚今天**连后端能力都还没有**，还有三枚**明令不许从面板发起**（票 194 AC#2）⇒ **不先摊代价就派写腿，一定会撞禁区或造出空壳方法**（"有名无实现"是这仓的老病：`panel.resync` 就是前例）。

## 1. 写面（只这些）

✅ 新建 `docs/evidence/s1/194-method-roster-census-c1.md`｜✅ 新建 `.scratch/wisp/probes/194/c1/**`（读数只往这里写）｜✅ 追加票 194 一段 Progress log（只追加、不改原句、不勾框）。
⛔ `internal/**`、`cmd/**`、`go.mod`／`go.sum`、`docs/PLAN.md`、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden **一字不动**；起手与终态都要跑 `git status --porcelain -- internal/ cmd/` 并把名册贴进表 §1。
⚠ **写面闸门（新规矩，`A374`）**：**终态必须等于起手名册（逐枚具名差集为空）**，⚠ **不是"必须为空"**——同树里有 `181-r2`、`185-r1` 两枚写腿在飞是常态，多出来的每一枚具名说"不是我的"，不许动、不许提交、不许还原。
⛔ `frontend/**`／`design/**` 不读不写不引不转述。⛔ 别家脏件不动不提交不评论。
⛔ 零删除命令；不许跑 `probes/161/r6/flip-declaration.sh`；不许改 `scripts/d22scan.sh`／`probes/154/gate-clauses.sh`；只 commit 不 push、显式 pathspec、禁 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`。
⚠ 重跑任何既有台件前先 `grep -n "WriteFile\|OpenFile" <那枚台件>` 看写出路径（`A367` 那枚坑）。

## 2. 必答五问（每条现跑；引规格与代码要**逐字引原句＋`file:line`**，不许转述）

**Q1 名册到底有几枚**：`sed -n '158,178p' docs/specs/SPEC-08-ui-ball-panel.md` 现读那张表，**逐枚列出方法名**（含"事件推送"那一行里的每一枚），并给出**你数出来的枚数**（⚠ 票 194 面写的是"18 枚"，那是我按眼睛数的、可能不准 ⇒ **以你现跑的为准并具名更正**）。代码那份名册现量：`grep -nE '=\s*"panel\.' internal/panel/bridge.go`（⚠ **别用 `grep -c 'panel\.'`**，那把尺起手是 5，因为 `bridge.go:35` 是含 `"panel.*"` 的注释——`33-r1` 报过，`A380`）。
**Q2 逐枚四问（本票主体，一张表 18±行）**：每枚答——① **今天有没有等价实现**（尺：`grep -rlF "<方法名>" --include=*.go internal/ cmd/ tools/ | grep -v _test.go` ＋ **按能力再找一遍同物异名**，例：`tasks.list` 可能已经以名册／快照的形式存在、只是不叫这个名字）；② 若有，落点在哪（`file:line`）；若没有，**最接近的现有地基**是哪一枚函数；③ **撞不撞票 194 AC#2 那三条禁区**（逐字引那三句：`SPEC-08:169` 的 allow 侧禁令／R20·票 92 的"权限输入口"／`SPEC-08:170` 的 `config.set` 安全节要 L2 重确认）；④ 需要新的门控／审计吗（对照 `l2_grant_boundary_test.go` 与 `composer_handlers.go` 今天的形状）。
**Q3 事件推送那一行单独答**：`task.delta`／`tool.chip`／`approval.request`／`ball.state`／`cost.tick` 这五枚是**Go→前端**方向，⚠ 出向管子今天真在跑（`cmd/wisp/panel_pump.go`，`33-a1` 现量）⇒ 逐枚答"这一枚今天有没有内容"（票 145 那张十六字段表是上游），**别把"管子有"说成"内容有"**。
**Q4 分批建议（不拍板）**：按"有货只差名字／有名字没货／撞禁区要先裁"三堆各给一枚清单＋每堆一枚最小写腿的写面范围（哪几枚文件），⚠ 并标出**哪些必须等真窗口那一环（票 33 的 H2/H3）才有生产调用者**——`33-r1` 已把六环建好但生产听众仍是零枚（`A380`），**补齐名册不等于按钮能点**，这条要逐枚说清，别让它变成又一枚"做了但没生效"。
**Q5 债务清点**：`AGENTS §2` 那条在册待定项「`C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿」今天定到哪一步（尺：`grep -rn "C24" docs/PLAN.md docs/specs/ | head`）；⚠ 本票只报现状**不替它定稿**。

## 3. 门禁（只读程也取数；⚠ 按 `A363` 不充当结案凭据）

`sh scripts/d22scan.sh`（基线 `ban #8 internal/`＝**438**；⚠ 若更高，那是 `181-r2`／`185-r1` 在飞新增件，**具名登记、不当漂移、不去动**）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册**（在册只 `G6neg`）；`go test -count=1 ./internal/panel/`（⚠ 今天 3 枚在册红：`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`，照实记、不许当绿、不许修）。⚠ **不要改那两把尺本身**。终态读数取在你自己最后一枚 commit 之后。

## 4. 表骨架（第 ≤5 枚调用内落盘并 commit 第一枚）

① 起手锚＋起手脏件名册 ② Q1 两份名册逐枚（含我票面枚数对不对的更正） ③ **Q2 逐枚四问主表** ④ Q3 事件推送那一行 ⑤ Q4 三堆分批建议 ⑥ Q5 债务清点 ⑦ 本程没测／判错的（逐名） ⑧ 门禁终态 ⑨ 被拒调用＋零删除自证＋工具调用终值 ⑩ next＝派写腿之前还缺哪几枚批准

## 5. 硬顶与计数

工具调用**硬顶 35**；**每答完一问自报一次累计**（写"本节末＝第 N 枚"）；**第 25 枚起不许开新探索**；表没落盘前不许继续取数；判"必须读 `frontend/**` 才能答完"＝**停手上报**。

## 6. 结束消息回我六项

① 名册真实枚数（我票面写的 18 对不对）；② Q2 主表三堆各几枚＋逐枚一句话；③ **撞禁区的那几枚逐名**（哪一枚撞哪一条）；④ Q3 那五枚事件推送"有管子没内容"的逐枚结论；⑤ Q4 分批建议（每堆写面范围＋哪几枚必须等真窗口）；⑥ 门禁四数＋红名册＋起手/终态名册差集＋工具调用终值＋有没有跑删除命令。
