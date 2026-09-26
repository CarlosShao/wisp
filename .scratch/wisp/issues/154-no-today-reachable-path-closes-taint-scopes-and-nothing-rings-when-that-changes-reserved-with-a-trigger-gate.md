# 154（**这一枚作废，勿派**）— 前提已被否证

- Status: 作废（作废原因见下；已改为另一枚文件名，见最后一行）
- 作废时刻：2026-09-26 10:3x，编排者自己作废自己写的这一版。

## 为什么作废

本文件第一版的前提是："**今天没有任何生产路径会关污点 scope，因为 `cmd/wisp/run.go:513` 写着 `rt.bridge.OpenTask("")`**"。
⇒ **那行代码不存在，且从未存在过**：`git log -S 'OpenTask("")' --all -- cmd/wisp/run.go` 零枚 commit；
`git grep -n OpenTask 86b0161 -- internal cmd | grep -v _test` 的唯一命中是 `internal/tools/bridge.go:559` 的 `b.OpenTask(dec.TaskID)`，**带的是真 task id**。
⇒ 真形状是反过来的：**"每轮一枚 scope"是可达的**，而票 151 的修法（`5d46f24`）已经把**环路那一枚 id** 关掉了；剩下的账是"**只关了一形**"，不是"没人关"。

## 那份假账的来路（要留痕，别让下一位以为这是个笔误）

那句"空串"与三枚配套的空号（`RunTextTask`／`unboundScopes()`／两枚 sha）**来自一枚仍在运行的验收程的中途输出**，被我当成终判写进了台账与这枚票。
⇒ 全案（含逐条否证命令）＝台账 `docs/reports/pending-and-issues.md` 的 **`A280`**；三处未提交的假账已在提交前撤销。

## 顶替它的那一枚

`.scratch/wisp/issues/154-the-close-only-covers-the-loop-task-id-so-host-supplied-task-ids-still-have-no-owner-and-concurrent-tasks-have-zero-readings-reserved-with-a-trigger-gate.md`
（同号码、五条前提全部换成编排者本轮现跑的命令；本文件按"临时件只建不删"的规矩**留在原地不删**，只挂这块牌子。）
