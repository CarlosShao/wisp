# 票 283 — `corr` 落地后**两处新命名的盲区**：回退形无仪器／身份改读无整链尺

**立票**：2026-10-08 20:2x 编排者（机主令「及时补票」；来路＝非实现者腿 `282-v1b` 的判语＋台账 `A749`）
**性质**：补尺或登记"不测"，都行；⛔ 但**不许默认不管**。
**Status:** **done**（2026-10-09 09:4x 编排者结案：三格全勾，凭据＝编排者自己另下的两发反例（种 A／种 C 各打红并逐字还原，hash 回基线）＋`git show --stat 039ec93c`＝+227/−0；非引腿的自述。残余两格具名转立：seed B `:206` 那条上游 L2 通道未追到底、`loop.go` 的 `CorrelationID: callCorr(...)` 那一发突变未跑 ⇒ 归 **`282` AC#3** 与票 242 一族，见台账 `A752`）

## 两处盲区（引用前先重跑）

1. **回退形无仪器**：把 `internal/tools/cancel.go:69` 那支的返回值种成 `h.corr` ⇒ 相关三枚用例**全绿**（`282-v1b` 实测；⚠ 派发路径全走 corr==taskID 的构造 ⇒ 分不出）。
2. **身份改读无整链尺**：`internal/tools/subagent_197.go:262`＋`internal/tools/task.go:708` 改读 `TaskID(ctx)` **前后**，同一形状三枚用例（197 深度／221 父停子／221 兄弟自停）**双双全绿** ⇒ 套件分不出改读对不对；缺"真 loop（per-call corr）→bridge"整链尺。

## 要建什么

- [x] **AC#1 查重**：立任何新尺前先跑三处查重（票池／`issues/README.md` 归口句／台账），⛔ 别与既有尺重复。**（2026-10-09 09:4x 编排者复跑核过：`grep -c "归口" issues/README.md`=0（rc=1）；台账里 `整链尺|回退形无仪器` 全仓 1 命中＝`A749` 那条**缺口登记**、非既有尺；`git grep -l "TaskID(ctx)" -- internal/tools/*_test.go` 只回本尺新文件＝无同形尺可重复。腿的三处名册见 `probes/283/r1/10-positive-controls.md` §AC#1）**
- [x] **AC#2 两处各给"补尺 / 登记不测"二选一**：选补尺 ⇒ 最小形状＋正控（种坏必红）；选登记 ⇒ 写明"为什么不测"＋把该盲区写进对应票面/注释的**具名节点**（⛔ 不许只留本票）。**（两处皆选补尺，落点＝新文件 `internal/tools/ticket283_corr_identity_rulers_test.go`（227 行）。正控由编排者自己另下两发、非引腿的读数：种 A `cancel.go:69 return h.taskID→h.corr` ⇒ rc=1，红句 `…_test.go:40 TaskID = "corr-283", want the task id` ＋ `:194 ParentTaskID="…#call-283-spawn"` ＋ `:206`；种 C `task.go:708 TaskID(ctx)→CorrelationID(ctx)` ⇒ rc=1，单元尺 PASS、整链尺 FAIL 于 `:206/:209`「拒绝停止：…不是调用者 …#call-283-cancel 的孩子」。两发均逐字还原：`cancel.go 87ac0624`／`task.go 3ec8d498`／`subagent_197.go 8c9266a7` 三枚 hash 回基线、`git status --porcelain -- internal/tools/` 空、复跑 rc=0／2 PASS。读数＝`probes/orch-283r1/{00-baseline,01-seedA,02-seedC,03-restored}.txt`）**
- [x] **AC#3 不许放宽**：两处任何动作⛔ 不许改松既有断言。**（`git show --stat 039ec93c`＝`1 file changed, 227 insertions(+)`、零删除；编排者另把自己两发突变只动产码返回表达式、未碰任何 `_test.go` 断言；三枚冻结件零碰；本程未跑整包、只按名靶向 `-run 'Test283'`）**

## 禁区

⛔ 零翻框；⛔ 三枚冻结件一字不动；⛔ 零 push。
