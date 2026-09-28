# 174-AC#2d 判据第一枚（写码腿 `174-r2`·纯测试面）

> 派单：`.scratch/wisp/dispatches/2026-09-28-151x-impl-174-r2-ac2d-first-criterion.md`
> 工单：`.scratch/wisp/issues/174-...md`（**本程只碰 AC#2d 这一格**）
> 上游普查：`docs/evidence/s1/174-wiring-cost-census-c2.md` ＋ 台账 `A372`
> 骨架落盘时刻：`2026-09-28 14:58 +0800`（骨架在第 6 枚工具调用落盘，比派单 §4 的「≤5 枚」晚一枚：同批并发了三枚必读件（工单／`internal/tools/task.go:215-364`／`task_output_pointer_notice_test.go`），已登记偏差不改判）

## ① 起手锚 ＋ 写面闸门

- 起手 HEAD：`40aee084`（`Mon Sep 28 14:56:28 2026 +0800` · `docs(evidence): 33-a1 land the twelve-section skeleton for the panel inbound-hop survey (no product code touched)`）
- 派单声明的编排者锚点：`ee0ef8ec` 之后的实际 HEAD ⇒ 共享树里别人在推进，**按实际 HEAD 做**，偏差＝起手锚是 `40aee084` 而非 `ee0ef8ec` 的直接后继（登记，不改判）。
- 起手写面自证：`git status --porcelain -- internal/ cmd/` ＝**空**（零脏件）。
- 本程写面：只新增 `internal/tools/` 下的 `*_test.go`（文件名 `task_output_canonicalize_fail_174_test.go`）。
- 本程不接线：`cmd/wisp/run.go` 的 `TaskDeps{}` 补 `Paths` 是下一枚腿的事（同批有写腿在碰 `run.go`）。

（待填：终态 `git diff --numstat <起手锚>..HEAD -- internal/ cmd/` 非 `_test.go` 行必须为空）

## ② 三形判据与名字

（待填：新判据函数名 ＋ 正向／反向／不越界各断言什么）

## ③ 变异现量（摘哪处会红，逐字）

（待填：`-overlay` 变异 ⇒ 逐字红腿读数）

## ④ 本程没测什么

（待填）

## ⑤ 门禁终态

（待填：`go test -count=1 ./internal/tools/` ／ `sh scripts/d22scan.sh`（`ban #8 internal/` 枚数解释）／ `gate-clauses.sh` 红腿名册 ／ `gofumpt -l`）

## ⑥ 被拒调用 ＋ 零删除 ＋ 终值

（待填）

## ⑦ next

（待填：接线那一行还欠谁）
