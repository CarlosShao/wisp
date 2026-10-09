# 票 287 — 深度判定那一支（`internal/tools/subagent_197.go:267-271`）**今天既没有仪器也没有后果**：缺的是夹具，不是一枚种刀

**立票**：2026-10-09 11:5x 编排者（机主令「及时补票」；来路＝写腿 `282-r1` 顶回票 282 `AC#3` 那句"缺第 3／第 4 发读数"，腿件 `.scratch/wisp/probes/282/r1/20-ac3-before-after.md` 91 行＋编排者独立复跑）
**性质**：⚠ **这一支不是缺陷**。票 242／282 把"任务身份"从 `CorrelationID` 拆到 `TaskID(ctx)` 之后，深度判定那一跳**在生产链上走不到**（孩子的 `task.spawn` 在 `subagentToolProvider.Execute` 就被硬拒）。本票要的是**一枚能让它走到的夹具**＋钉住它，属于"纵深防御今天不可测"那一族。

## 现量（引用前先重跑；尺一律 `git show HEAD:<path>`）

- `internal/tools/subagent_197.go:262` 逐字 `parentID := TaskID(ctx)`（11:4x 编排者现读，还原后的形状）。
- **种坏 `:262` 会让身份那一跳红**（`ticket283_corr_identity_rulers_test.go:194` 指名红句带 `…#call-283-spawn` 对 task id 不等）；**而把同一枚访问器换成 `CorrelationID(ctx)` 时，深度那一支前后两发逐字同绿**：`--- PASS: Test197ChildCannotDeriveSubagent (0.00s)`（腿 `282-r1` 两发读数，命令 `go test ./internal/tools/ ./internal/agent/ -count=1 -v` 逐字相同、只有盘上那一行不同）。
- **三把尺证明它今天不可达**（腿现量，编排者未逐把复跑＝〔仅腿报〕，AC#1 要求复跑）：① 拒句「子代理不许再派子代理」在 `*_test.go` 命中 **0**；② `TaskKindSubagent` 的 5 枚测试命中**全部**是"断言孩子那行的 Kind"，没有一枚把**调用者**那行摆成 subagent；③ 生产链上孩子的 `task.spawn` 在 `subagentToolProvider.Execute`（`:573-580` 一带）被先行硬拒 ⇒ 根本到不了 `:267`。
- 相关：票 282 `AC#3` 已按"身份那一问"翻勾，本票承接它**没闭合的那一半**（票面 ✅ 注里具名转立）。深度上限这一维的契约出处＝`A399`／`A401`（09-28 定"并发 8／深度 1／子级永不自批"）。

## 要建什么

- [ ] **AC#1 先把三把尺自己复跑一遍**（不许照抄本票现量）：拒句在 `_test.go` 的命中枚数、`TaskKindSubagent` 的测试命中逐枚读法、以及"`:267` 到不到得了"的调用形状尺。⇒ 若任一把握**不复现**，本票前提当场作废并具名报回。
- [ ] **AC#2 造那枚夹具**：让**调用者那一行**的 `Kind == TaskKindSubagent`、派发经**真 loop**（因而 `corr != taskID`）、请求 `task.spawn` ——三件缺一都作不出读数（方法学前置＝票 282 那条：`TaskID != CorrelationID` 必须真在红句里同屏）。判据＝**新增用例今天必须是"能红"的**：把它种坏（把 `:267` 那一支的深度判定摘掉／改成恒通过）⇒ **指名用例必须红**，四件套齐（种前 hash／`sed -n` 复量／红句 `文件:行` 逐字／还原 hash 等值＋`git status --porcelain -- internal cmd` 空）。⛔ 不许用"字段非空"充当判据、⛔ 不许 `t.Skip`。
- [ ] **AC#3 顺手答一句"要不要让它走到"**（⛔ 不改产码）：孩子的 `task.spawn` 被 `subagentToolProvider` 硬拒 ＝ **设计如此**（深度 1）还是**挡住了一枚本该走到的检查**？两形各写"要动哪几枚文件＋会不会把深度上限放开"，交编排者裁。⚠ 这一格交的是**代价表**，⛔ 本票不许落地。
- [ ] **AC#4 边界与门禁**：写面只许 `internal/tools/**`（夹具＋用例）；⛔ 不碰 `internal/agent/loop.go` 的 `callCorr`、⛔ 不碰 `approval/gate.go`、⛔ 三枚冻结件一字不动（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；`GOFLAGS= go build ./...`、`go test ./internal/tools/ -count=1` 改前改后各一次（**PASS/FAIL/SKIP 三数并排＋红名册差集**）、`gofmt -l`／`gofumpt -l` 对动过的件空、`sh scripts/d22scan.sh` 纯净树 rc=0；`frontend/**`／`design/**` 零读零写零转述；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许为变绿放宽任何既有断言；⛔ 不许"顺手"把深度上限改成 2 或放开子级 `task.spawn`（那是契约面，须人工批准）；⛔ 不许把 AC#3 的代价表读成已裁。

**Status:** **未开工**（2026-10-09 11:5x 立票）。排程＝⛔ 按住等 Go 编译/测试面空（此刻 `286-w1` 独占），且**按在票 282 结案之后**——本票承接它未闭合的那一半，先有归属再开工。⛔ 零翻框、零 push。
