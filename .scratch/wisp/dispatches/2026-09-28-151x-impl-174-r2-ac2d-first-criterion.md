# 派单：`174-r2`（写码腿·**纯测试面**）——给"连规范化都没通过"那一支补上第一枚常驻判据（AC#2d）

- 派单时刻：`2026-09-28 15:1x`｜编排者锚点 HEAD＝**`ee0ef8ec`** 之后的实际 HEAD（共享树里别人在推进；**按实际 HEAD 做并登记偏差**）
- 工单：`.scratch/wisp/issues/174-task-output-points-at-a-file-the-model-cannot-read-by-default-spill-artifacts-are-outside-fs-allowed-dirs.md`（**只做 AC#2d 这一格**）
- 上游读数（已入库、编排者复跑同值）：`docs/evidence/s1/174-wiring-cost-census-c2.md` ＋ 台账 `A372`

## 0. 这一单干什么（一句话）
`174-c2` 现量：**"路径连规范化都没通过"那一支今天没有任何判据钉着**（尺：`grep -c '规范化' internal/tools/task_output_pointer_notice_test.go`＝**0**；`grep -rn "连规范化都没通过" --include=*_test.go internal/`＝**零命中**）⇒ 摘掉它今天全仓不红。**本程只补这一枚判据，不接线、不改文案、不动产码。**

## 1. 写面（只这些）
✅ 新建或追加跟踪测试件：`internal/tools/` 下（建议新文件 `internal/tools/task_output_canonicalize_fail_174_test.go`，命名带 `174`；⚠ **不要改 `task_output_pointer_notice_test.go` 里任何既有用例的断言**）
⛔ **产码零字节不动**：`internal/tools/task.go`、`internal/tools/bridge.go`、`cmd/wisp/run.go`、`internal/risk/**` 一律不许改（起手与终态 `git status --porcelain -- internal/ cmd/` 都必须为空——**你只加 `_test.go`**）
⛔ **不碰** `docs/PLAN.md`（含 `:2564` 那行）、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量；⛔ 两枚模板冻结钉（`task_output_pointer_notice_test.go` 的 `:299`／`:349`，其中一枚 md5 基线 `8bd433707c1fd9039ac2103da08e7d22`）**一字不改**
⛔ **本程不许接线**（`run.go` 那一行 `TaskDeps{}` 补 `Paths` 是下一枚腿，且同批有别的写腿在碰 `cmd/wisp/run.go`，撞了会互相吃读数）
⛔ `frontend/**`／`design/**` 不读不写不引不转述；不属于你的脏件不动不提交不评论
⛔ Git：只 commit 不 push；显式 pathspec；禁 `add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`／`switch`；**零删除命令**；不许跑 `probes/161/r6/flip-declaration.sh`
⛔ **AC 框一枚不许自己勾**（勾由编排者翻）

## 2. 判据要长什么样（三形，缺一形＝装饰）
1. **正向**：规范化**返回错误**那一形（`Canonicalize` 返错，不是"规范化成功但不在授权根内"——那一形已有判据）⇒ `task.output` 的回执必须走"按读不到处理"那一支，且**不许**是空成功。
2. **反向（这枚最值钱）**：**摘掉那一支的 fail-closed 处理，必须有一枚用例红**。⇒ 用 `-overlay` 造变异现量（**别动跟踪件**），把"未接线／规范化失败"那支换成空串或换成"当作可读"，读数逐字进表。⚠ 若摘掉它全仓不红＝你写的判据没牙，重写。
3. **不越界**：同一发要断言**回执没有把这条路径说成"能读"**（不许出现"全文见 … 读得回来"这类承诺），但**不许**顺手断言精确文案字节（那会撞上模板冻结钉）。

⚠ 落点提示（现跑取号，别照抄）：今天真机走的那支在 `internal/tools/task.go:320`（`grep -n "路径授权判定者未接线" internal/tools/task.go`），调用链是 `task.go:242-243` → `pointerNotice` → `:318-330`。

## 3. 门禁（终态取在你自己最后一枚 commit 之后）
`go test -count=1 ./internal/tools/`（基线 `ok 15.812s`）；`sh scripts/d22scan.sh`（基线 rc=0、`ban #8 internal/` **examined=433**，你新增跟踪测试件会变 **434**，具名解释）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（在册只 `G6neg`＝票 178；⚠ **不要改这把尺本身**，同批有程在跑名册差集）；`"$(go env GOPATH)/bin/gofumpt" -l <你的新件>` 必须为空。⚠ 同批另有 **3 枚程在飞**（`181-r1` 写 `internal/panel/**`＋可能碰 `cmd/wisp/run.go`；`180-c1＋182-c1`、`33-a1` 只读）⇒ **你的写面只有 `internal/tools/*_test.go`**，别去碰它们的文件。

## 4. 增量交付与硬顶
工具调用**硬顶 30 枚**，**到第 22 枚停止新探索**。第 ≤5 枚调用内把表骨架（`docs/evidence/s1/174-ac2d-criterion-r2.md`：① 起手锚＋写面闸门 ② 三形判据与名字 ③ 变异现量（摘哪处会红，逐字）④ 本程没测什么 ⑤ 门禁终态 ⑥ 被拒调用＋零删除＋终值 ⑦ next）落盘并 commit 第一枚。判"必须改产码才能让判据红"＝**停手上报**（那说明那一支今天根本不可达，是要登记的新事实，不是让你去改产码）。

## 5. 结束消息回我七项
① 新判据的函数名与三形各断言了什么；② 变异现量（摘掉哪一处 ⇒ 哪枚用例红，逐字读数）；③ **产码零改动自证**（`git diff --numstat <起手锚>..HEAD -- internal/ cmd/` 里非 `_test.go` 的行**必须为空**）；④ 门禁四数＋红腿名册＋`ban #8` 枚数解释；⑤ 两枚模板冻结钉 md5 复量未变；⑥ 被拒调用逐条＋有没有跑删除命令＋工具调用终值；⑦ next（含"接线那一行还欠谁"）。
