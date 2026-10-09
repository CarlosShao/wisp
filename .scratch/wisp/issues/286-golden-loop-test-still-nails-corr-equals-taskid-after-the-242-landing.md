# 票 286 — `internal/agent` 里有一枚**起手即在的红**：`TestGoldenSingleToolCall` 还钉着"corr 等于 taskID"，而票 242 落地的形状已把它铸成 `taskID#call_e1`

**立票**：2026-10-09 11:2x 编排者（机主令「及时补票」；来路＝写腿 `285-r2` 具名报回"我未裁定、未放宽、未顺手改"＋编排者自己现跑复认）
**性质**：**这不是新 bug，是我 10-08 落地那跳留下的未收口账**——产码按契约改了形状，**同一形为的测试期望没人跟着搬**，于是树里从今天起一直躺着一条红。⛔ 本票不新造功能、⛔ 不动 golden。

**Status:** **已裁完，待写腿**（2026-10-09 11:3x 编排者：`AC#2` 权威判定已下＝搬测试期望、不回退产码；`AC#1`/`AC#3` 派只读归因腿，`AC#4`/`AC#5` 按在其后）。本格未全勾 ⇒ ⛔ 不加 `-done`。
## 现量（引用前先重跑；行号与被引短语同一次取数，尺一律 `git show HEAD:<path>`）

- 编排者现跑（同一条命令、`-count=1`）：`go test ./internal/agent/ ./internal/tools/ -count=1 -v` ⇒ **rc=1**／`internal/agent` **1 FAIL**（`--- FAIL: TestGoldenSingleToolCall`）／`internal/tools` `ok 13.6s`。
- 红句逐字（我的运行输出）：`loop_golden_test.go:70: call identity = {TaskID:03a55baf-… CorrelationID:03a55baf-…#call_e1 CallID:call_e1 Name:echo …}, want task 03a55baf-…`。
- 断言原文（`HEAD:internal/agent/loop_golden_test.go:70`）：`if call.Req.TaskID != res.TaskID || call.Req.CorrelationID != res.TaskID {` ⇒ 它要求**请求侧的 corr 逐字等于 taskID**。
- 产码现形（`HEAD:internal/agent/loop.go:603` `callCorr(taskID, callID string, index int)`，唯一调用点 `:676` `TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i)`）⇒ 带 call id 时返回 `taskID + "#" + callID`，**永不等 taskID**（除非 callID 为空走 `%s#call-%d`，那也带后缀）。
- ★**同文件里还躺着一枚绿的反形**：`HEAD:internal/agent/loop_golden_test.go:341` 逐字 `if row.Tool != "echo" || row.CorrelationID != res.TaskID {` ——它钉的是 **journal 行**那一侧，而 journal 侧至今仍是 taskID（`loop.go:369` `newTaskJournal(l.opt.Journal, taskID, taskID)`，第二实参经 `journal.go:59` 落到 `:81 CorrelationID: t.corrID`）。⇒ **同一枚文件里两枚同名"correlation id"的期望现在互相矛盾：请求侧那条被真值打红、账本侧那条靠残余保持绿。**这笔残余不是我今天才发现：票 282 `AC#4`(a) 已裁"真存在／判留／归票 230 那一族"（台账 `A758`）。
- ⚠ 射程（⛔ 别把这行当全集）：本票只处理**请求侧那一枚红**；`285-r2` 报的是"在锚 `ce18b3d6` 上起手即红"，我没有在旧 commit 上重跑过（⛔ 共享工作树里禁 `checkout`/`stash`/worktree），所以**"哪一笔落地让它变红"属 §要建什么 的 AC#1，不在现量里**。

## 要建什么

- [ ] **AC#1 归因到那一跳**（只读腿，⛔ 不许 checkout）：用对象层两态对拉判出"请求侧 corr 何时开始带 `#call-` 后缀"——尺＝`git show <ref>:internal/agent/loop.go`（`callCorr` 在不在／`:676` 那行形状）× `git show <ref>:internal/agent/loop_golden_test.go:70`（期望在不在），逐枚给出候选 commit 名册与**边界那一枚**；同时答："那枚期望自它诞生起钉的是什么"（⛔ 不许按票号猜，读代码）。⚠ 这一格交的是**归因**，⛔ 不交修法。
- [ ] **AC#2 谁是权威（本格归编排者裁，⛔ 实现腿与验收腿都不许自裁）**：`loop_golden_test.go:70` 那枚等值期望，与 `docs/PLAN.md:1368` 的 C18 句（"corr 坐在 task id **旁边**、⛔ 不要求等于它"）＋台账 `A745` 我给落地腿的派单原话（"同一任务并发两问 ⇒ 两枚 corr **不同**且各自可路由"）冲突 ⇒ 判"该搬测试期望"还是"该回退产码"。先例已盘上：票 282 `AC#2` 我裁过"那枚等值钉当年钉的是**这个 bug 本身**"（`A758`）。⛔ **无论判哪一支，"把断言删掉／改成永远通过"都不是选项。**
- [ ] **AC#3 全仓同类旧钉名册（按语义变体扫，⛔ 只搜一种字样＝本仓抓过的假话形）**：把"corr 等于 taskID"这一形的所有写法扫出来——至少含 `CorrelationID != res.TaskID`／`== taskID`／`corr == t`／位置参数 `taskID, taskID`／`CorrelationID: taskID` 诸形，逐枚给"该搬／该留（并写清为什么这一侧允许相等）／待人裁"。⚠ 名册里**必须并排 `:341` 那一枚**（它今天靠 journal 残余而绿，是本票最容易被"顺手一起改"带跑的一枚）。
- [ ] **AC#4 搬完之后要有仪器证明它没被搬松**（由写腿交、非实现者裁）：① 改前整包红名册（具名 `TestGoldenSingleToolCall`）；② 改后 `go test ./internal/agent/ -count=1` **rc=0**、且**不许新增 `t.Skip`**、PASS 枚数只增不减；③ **反形自证**：把产码 `callCorr` 换回"恒返回 taskID"那一形 ⇒ 新期望**必须红**（证明搬过的钉仍然咬得住，不是一条被抹平的等值句）；四件套齐（种前 hash／`sed -n` 复量／红句 `文件:行`／还原 hash 等值＋`git status --porcelain -- internal cmd` 空）。
- [ ] **AC#5 越界与门禁**：⛔ 零触碰 `internal/agent/testdata/golden/**`（尺＝`git show --stat` 名册里不许出现 testdata 字节）、`docs/PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／`allowlist.txt` 一字节不动；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）一字不动；`frontend/**`／`design/**` 零读零写零转述；`sh scripts/d22scan.sh` 纯净树 rc=0；`gofmt -l`／gofumpt 对动过的件空输出。

## 禁区

- ⛔ **不许为变绿放宽断言**：搬之后的期望必须**至少和原来一样严**（要钉的是"corr ＝ `taskID` 前缀 ＋ 该 call 自己的那一枚 id"这一形，⛔ 不是"corr 非空"）。
- ⛔ 不许动 journal 那一侧（`:341`／`loop.go:369`／`journal.go:81`）来"让两枚一致"——那是票 282 `AC#4`(a) 已裁**留**的残余，动它要先在票 230 那一族里挂号，⛔ 不在本票射程。
- ⛔ 不许新造 D43 之外的状态词；不许把 `SessionID`／grant 塞进 `tools.Decision`／`agent.ToolRequest`（票 224／242 撞过同一枚反射钉）。
- ⛔ 票 284 那一格（`gate.go` 注释）**不搭本票的车**：本票不动 `gate.go`，别顺手捎带。
- ⛔ 零 push；commit 必带显式 pathspec。

## AC#2 编排者裁决（**2026-10-09 11:3x**，本格归我裁，写腿与验收腿⛔不许自裁）

- **裁＝搬测试期望，⛔ 不回退产码**。凭据三条（同一次现读）：① `docs/PLAN.md:1368` 的 C18 句把 corr 定义成"坐在 task id **旁边**、⛔ 不要求等于它"；② 台账 `A745` 我给落地腿的派单原话＝"同一任务并发两问 ⇒ 两枚 corr **不同**且各自可路由"；③ 票 282 `AC#2` 我已裁过同形先例——那枚等值钉当年钉的是**这个 bug 本身**（corr 逐字＝taskID），保留原字面＝造一格"只能靠回退 bug 才能满足"的死格（`A758`）。⇒ `loop_golden_test.go:70` 的期望属"钉住旧缺陷"那一族，**该搬**。
- **搬成的形状（至少和原来一样严，⛔ 不许删、不许改成"非空即可"**）：请求侧那枚 corr 要同时钉 ① 以 `res.TaskID` 为**前缀**、② 后缀＝该 call 自己的 id（golden 回放里逐字 `call_e1`）、③ **不等于 `res.TaskID`**；再加一枚**反形自证**（把 `callCorr` 换回"恒返回 taskID" ⇒ 新期望必须红）。
- ⛔ **同文件 `:341` 那一枚不在本票射程**（journal 侧的 `row.CorrelationID == res.TaskID` 今天靠 `loop.go:369 newTaskJournal(l.opt.Journal, taskID, taskID)` 这枚残余而绿；那笔 10-09 已由票 282 `AC#4`(a) 裁"留"、归票 230 那一族，`A758`）。⛔ 任何程不许为了让两枚"看起来一致"去动 `:341`／`journal.go:81`——那会**悄悄改账本侧的形状**。
- **排程**：AC#1（归因）＋AC#3（全仓同类旧钉名册）＝只读腿可办，⛔ 不许带产码改动；AC#4（搬＋反形自证）＝写腿，按在 AC#1/AC#3 之后（先有名册再动手，否则会漏改第二枚同形）。AC#5 门禁随写腿那一笔同批交。
