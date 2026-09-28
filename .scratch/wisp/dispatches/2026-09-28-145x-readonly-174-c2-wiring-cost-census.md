# 派单：`174-c2`（只读代价普查·零产码）——"把路径判定者接到生产"这一跳要动哪几层、会把谁顶出去

- 派单时刻：`2026-09-28 14:5x`｜编排者锚点 HEAD＝`03c01d71`
- 工单：`.scratch/wisp/issues/174-task-output-points-at-a-file-the-model-cannot-read-by-default-spill-artifacts-are-outside-fs-allowed-dirs.md`
- 只答三格：**AC#2b（接线代价表·不落地）**＋**AC#2d（那一支的形状确认）**＋**AC#3（三条禁区自证，只读能答的部分）**
- 性质：**只读** ⇒ **AC 框一枚不许勾**、**产码零字节不动**、**不落地、不选方案**。

## 0. 一句话背景
票 174 的主缺陷（回执指着一条模型默认读不到的路径）已由 `174-r1` 把**文案**修成说实话；剩下的是"**把路径判定者真接到生产**"那一跳——`cmd/wisp/run.go` 里 `TaskDeps{}` 的 `Paths` 字段今天**没赋值**（尺：`grep -n "TaskDeps{" cmd/wisp/run.go`，⚠ 号会漂、以你现跑为准）。⇒ 本程**只交代价表**，接不接、什么时候接归编排者裁。

## 1. 写面（只这些）
✅ 新建 `docs/evidence/s1/174-wiring-cost-census-c2.md`（骨架先落盘）｜✅ 新建 `.scratch/wisp/probes/174/c2/**`｜✅ 追加票 174 一段 Progress log（只追加、不改原句、不勾框）。
⛔ `internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量 **一字不动**；起手与终态两条 `git status --porcelain -- internal/ cmd/` 都必须为空，不为空**停下报我**。
⛔ 三条禁区一票否决（照票 174 AC#3 原文，**不许自己判"这次例外"**）：(i) 不许把宿主内部 artifacts 实现成**受门控的 Tool**；(ii) 不许在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做决策；(iii) 不许把 `artifacts` 目录**塞进 `allowed_dirs` 默认值**（⚠ `allowed_dirs` 是**判级输入**、不是执行时硬边界——这条已定案，别顺手改成硬边界）。
⛔ `frontend/**`／`design/**` 不读不写不引不转述；不属于你的脏件不动不提交不评论。
⛔ 不许跑任何删除命令；不许跑 `probes/161/r6/flip-declaration.sh`；只 commit 不 push、显式 pathspec、禁 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`。
⚠ 重跑既有台件前先 `grep -n "WriteFile\|OpenFile" <那枚台件>` 看写出路径（就地覆盖型 sink 会洗掉别家票逐行引用的读数＝`A367`）。

## 2. 必答（每条都要现跑，不许推理）

**Q1 接线那一跳的完整形状**：从 `cmd/wisp/run.go` 的 `TaskDeps{}` 到 `internal/tools` 真用它的那处，逐环列：谁构造、谁持有、谁调用（每环一条 `grep -n` 尺＋文件:号）。⇒ 明确"要动几枚文件、每动一处哪几枚既有判据会红"（**逐枚点名**，不许只报枚数）。
**Q2 顶出去效应（这一问最值钱）**：接上之后，**哪些今天能读到的路径会变读不到、哪些今天被拒的会变放行**？现读那枚判定者的语义（`grep -n "func .*ResolveWorkspace\|allowed_dirs" internal/risk/pathresolver.go internal/tools/paths.go`），并给一发**只读台件**（放 `probes/174/c2/`，走 `-overlay`，**别动跟踪件**）量出"接上前后"两发对照。⚠ 收益与退让**必须一起报**（这是这仓写死的老规矩）。
**Q3 与票 183/185 那根管子的关系**：接上之后，`task.output` 回执里那条产物路径**是否就自动读得回**？（⚠ 别假设——票 183 的教训是"包级全绿≠端到端通"。）给一条现跑或明确写"本程未验"。
**AC#2d**：本格原文说它"不是 AC#2 的缺陷形"，而是 `174-v1` 现跑发现的**零覆盖那一支**。⇒ 现量它今天到底有没有判据钉着（尺：`grep -rn "<那一支的关键行为>" --include=*_test.go internal/ | head`），**没有就明说没有**，不许拿"文案已改"充当。
**AC#3 可只读自证的部分**：三条禁区逐条给"我这程没碰"的证据（`git status` ＋ `git diff --numstat` 删除列逐枚 0）。

## 3. 门禁（按 `A363` 薄规矩不充当结案凭据）
`sh scripts/d22scan.sh`（基线 `ban #8 internal/` **examined=433**）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册**（在册只有 `G6neg`）；`go test -count=1 ./internal/tools/`（单跑）。⚠ **不要改 `gate-clauses.sh` 那把尺**。终态读数取在你自己最后一枚 commit 之后。

## 4. 表骨架（第 ≤5 枚调用内落盘并 commit）
① 起手锚＋写面闸门 ② Q1 接线逐环 ③ Q2 顶出去两发对照 ④ Q3 与 183/185 的关系 ⑤ AC#2d 有无判据 ⑥ AC#3 三条禁区自证 ⑦ 本程没测什么（逐名） ⑧ 门禁终态 ⑨ 被拒调用＋零删除自证＋工具调用终值 ⑩ next＝落地腿派之前还缺什么（含"要不要人先批准哪一格"）

## 5. 硬顶
工具调用**硬顶 35 枚**，**到第 25 枚停止新探索**；每答完一问 commit 一枚；判"必须改产码才能答完"＝停手上报。

## 6. 结束消息回我六项
① 接线逐环（动几枚文件、哪几枚判据会红，逐枚名）；② 顶出去那一对读数（接上前／后）；③ Q3 的结论或"未验"；④ AC#2d 有没有判据钉着；⑤ 门禁四数＋红腿名册＋两条 `git status` 为空；⑥ 工具调用终值＋被拒调用＋有没有跑删除命令＋next。
