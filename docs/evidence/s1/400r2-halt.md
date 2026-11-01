# 400-r2 停手件：wire 携 DroppedCount 需先动 internal/panel wire 类型面（未定义即停②）

- 腿＝`400-r2`（票 400 残余①执行腿、写腿·实现者）；开工锚＝`5bc564081`（dev，恰为派单锚点；索引零已跟踪改动、写面 `cmd/wisp` 无他人脏）。
- 派单使命＝把 400-r1 落的 `DroppedCount()` 真挂到面板 wire（写面限 `cmd/wisp/panel_pump.go`）；三条硬边界之一＝⛔ `internal/panel/**`（400-r1 刚落、⛔ 碰）。
- **停手结论＝未定义即停②命中（派单原文：「wire 必须动 panel 契约 ⇒ 停手——那违反『只是接线』」）**：wire 携 DroppedCount **必须先动 `internal/panel/subagent_roster_197.go` 的 wire 类型面三处**，外加 TS 名册跟键；`cmd/wisp/panel_pump.go:250` 一行只是四分之一，且单独动它**没有左值可赋、连编译都过不了**。本腿零产码改动停手。

## 断点链路（四处现读，HEAD=`5bc564081`）

1. **数在源头已就位（可读）**：`StreamLog.DroppedCount() int`（`internal/panel/pump.go:750-757`，`droppedTotal` 总数 accessor，注释自述「including the ones the naming ledger no longer holds room for」）——400-r1 落的就是这一半＋Ring 64（`streamDroppedNamesCap`，pump.go:520）。
2. **装配点（与本票现读对得上，但无处可挂）**：`cmd/wisp/panel_pump.go:250`＝`state.DroppedStreamKeys = rt.stream.DroppedKeys()`（与 400-v1-verdict.md:52 现读逐字一致）。但 `state` 的类型 `panel.TaskRosterState`（`internal/panel/subagent_roster_197.go:94-101`）字段集＝`Rows / InFlightSlots / PoolCap / StreamTruncated / StreamElidedRunes / DroppedStreamKeys`——**无 `DroppedCount` 字段**。`state.DroppedCount = rt.stream.DroppedCount()` 没有左值，单改装配点编译不过。
3. **wire 类型面无位置**：`panel.TaskRosterSection`（同文件 :150-172）wire JSON 键集＝`rows / inFlightSlots / poolCap / streamTruncated / streamElidedRunes / droppedStreamKeys` 共 6 键（与 400-v1 手核的「TaskRosterSection 6」一致）——**无 `droppedCount` 键**；render `taskRosterSectionFrom`（:210-223）逐字段拷贝 state→sect，无 DroppedCount 拷贝行可挂。
4. **TS 名册同形状**：`frontend/src/lib/panel.ts:221-227` `TaskRosterSection` 同 6 键、无 `droppedCount`。Go 侧加键则 TS 名册必须跟——正是 400-r1/400-v1 刚治的「Go 侧恒发、TS 名册未跟」同病；该树归机主外部 agent（本编队 ⛔ 写面），且 400-v1-verdict.md:17 追认附条件③逐字「**此后此类跳点一律先转界面支或先取机主预批，不许再以『具名判断』自我越权**」。

## 根因：两份在案文书的措辞与现读错位

- `400r1-impl.md:46` follow-up① 逐字「**`TaskRosterState`/wire 若要携带 `DroppedCount` 需动 cmd/wisp 的装配点（panel_pump.go:250），本腿写面不含 cmd/wisp，未动**」；`400-v1-verdict.md:52` 残余①同向（「wire 携 `DroppedCount` 需动 cmd/wisp 装配点 **panel_pump.go:250**」）。
- 两条措辞都把「需动」只落在 panel_pump.go:250，暗示 wire 类型已就位；现读坐实 400-r1 在 wire 类型面（TaskRosterState／TaskRosterSection／render 三处）**一处未动**，落的只是 StreamLog 侧 accessor。派单据此措辞把写面压在 cmd/wisp——写面与病灶错位。
- 装配点本身（:250 的内容）与 400-v1 现读**一致**；对不上的是「只动装配点即可挂 wire」这个前提，故本腿按未定义即停②停手具名，而非①。

## 未走的路（为什么不就地变通）

- **单方改 internal/panel**：撞硬边界 1 的路径禁区（⛔ `internal/panel/**`），且 wire 加键属契约面扩展——「wire 只是接线」的前提不成立。
- **只加 Go 侧不加 TS**：仪器（对账尺）只对顶层键、嵌套类型无仪器覆盖 ⇒ 不会红，但那是**刻意制造下一轮「Go 恒发、TS 未跟」病**（400-r1 AC#1 治的正是它），且打破 400-v1 手核在案的「TaskRosterSection 6 键一字不差」对齐状态。
- **污染现有键**（如把总数编进 `droppedStreamKeys`）：动契约语义，比加键更糟。
- 三条路都违反派单明文边界，故停手上报，不自行假设。

## 建议的重派形状（供编排者立票）

1. `internal/panel/subagent_roster_197.go` 三处：`TaskRosterState` 加 `DroppedCount int`；`TaskRosterSection` 加 `DroppedCount int `json:"droppedCount"``（含注释沿用 DroppedStreamKeys 的「never unread」口径）；`taskRosterSectionFrom` 加拷贝行。
2. `cmd/wisp/panel_pump.go:250` 后加 `state.DroppedCount = rt.stream.DroppedCount()` 一行。
3. `frontend/src/lib/panel.ts` `TaskRosterSection` 加 `droppedCount: number`——**该枚在机主外部 agent 域，须按 400-v1 附条件③走机主预批或由编排者转界面支**。
4. `cmd/wisp/subagent_carrier_197_test.go` 同族断言（:561 现形状 `sect.DroppedStreamKeys` vs `stream.DroppedKeys()` 可照抄为 DroppedCount vs `stream.DroppedCount()`）＋flood 反向（摘装配点行⇒面板侧读数零）。
5. 重派写面须如实含 `internal/panel/**`（并处理与 ⛔ 禁区的关系：要么编排者解禁该一枚文件，要么把①③划归后续腿）。

## 本腿读数与纪律格

- 开工三核：HEAD=`5bc564081`（dev）＝派单锚点；`git status --porcelain --untracked-files=no`＝空（零已跟踪改动）；写面 `cmd/wisp` 无他人脏（397-r4 已交、cmd/wisp 净）。
- 本腿**零产码改动**（`.go` 零枚）：`go build ./cmd/wisp ./internal/panel` rc=0（树健康、现读非环境伪）；四数／两尺／d22scan 无对象（零 `.go` 改动 ⇒ 树与 HEAD 逐字节同，零回归恒真；gofumpt/d22scan 的扫描对象不存在）。停手时树状态存档 `D:\tmp\wisp\20261101\400r2\logs\01-halt-tracked-status.txt`、`02-build-rc.txt`。
- 本腿写面＝本文档一枚＋仓外 `D:\tmp\wisp\20261101\400r2\**`；⛔ 未动 `internal/panel/**`、`frontend/**`、thresholds/golden、testdata、票面、台账、specs、PLAN；⛔ push；commit 显式 pathspec 仅本文档。
- 未定义即停检查：派单三条停手条件中②命中（①装配点现读一致；③链路未断，断的是写面边界），全程未做未授权假设，停手时点在开工三核＋四份先读件＋四处断点现读全部完成后。
