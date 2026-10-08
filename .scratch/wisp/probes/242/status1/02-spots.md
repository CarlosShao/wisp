# 242-status-1 三处具名朗读（行号与被引短语同一发现取；时点 2026-10-08 19:1x，HEAD `fe2b6b16`）

## 1. 产码默认：`internal/agent/loop.go:654`

- 尺：`grep -nE 'TaskID: taskID|CorrelationID: taskID' internal/agent/loop.go`（rc=0，命中 `:654`）＋`sed -n '645,662p'`（rc=0）。
- 逐字（`:654`，位于 `ToolRequest{...}` 构造、`l.reg.Spawn("tool-exec-"...)` 之前）：
  `			TaskID: taskID, CorrelationID: taskID, CallID: p.call.ID,`
- 判：产码默认**仍**把两枚 id 塌成同一个值 ⇒ 真进程里造不出两枚不同 correlation 的活卡。A737（19:0x）已登记此遗留：逐字引同一行、**归口"票 242 一族"**（引票 `:14`／`:60`）、查重三处后判"不新立票"。

## 2. 载具：`cmd/wisp/subagent_selfapproval_197_test.go`（旧 `:109` 一带）

- 尺：`grep -nE 'TaskID:|CorrelationID:' cmd/wisp/subagent_selfapproval_197_test.go`（rc=0）＋`sed -n '100,118p'`（rc=0）＋`grep -n 'startChildWrite197(' ...`（rc=0）。
- 现状逐字：
  - `:101`（签名）：`func startChildWrite197(rt *agentRuntime, taskID, corrID, path string) *childWrite197 {`
  - `:116`（原 `:109` 那句的现位置）：`			TaskID: taskID, CorrelationID: corrID,`
  - `:396`：`	w1 := startChildWrite197(rt, childID, childID+"-corr-1", file1)`
  - `:463`：`	w2 := startChildWrite197(rt, childID, childID+"-corr-2", file2)`
- 判：**今天已拆成两个不同值**（`-corr-1`／`-corr-2`；旧 `:109` 逐字 `CorrelationID: taskID` 已不在盘上）＝259-r4 `77150dcc`、票 259 `AC#4` 翻勾（A737 亲跑 `go test . -run '197' -v` rc=0，日志逐字 `corr1=…-corr-1 corr2=…-corr-2`）。

## 3. 票面 `:14`／`:60` 两句在盘上今天成不成立（整行已读）

- **`:14`**（现量行，09-30）原文逐字（现量格＋尺格）：
  `⚠ **造不出"两张不同名"的卡**（现成的载具缺陷） | `cmd/wisp/subagent_selfapproval_197_test.go:109` 逐字 `TaskID: taskID, CorrelationID: taskID` ⇒ 同一次跑里两张卡**共用同一个 correlation id**`
  **我的判**：在**载具**层面**已不成立**（`:116` 已拆、`:396`／`:463` 两个不同 corr 值）；在**真进程产码**层面**仍成立**（`loop.go:654` 同形未改）。
- **`:60`**（10-03 "本票终态"段）原文逐字开头：
  `**载具前置仍未做**：`cmd/wisp/subagent_selfapproval_197_test.go:109` 现量逐字 `TaskID: taskID, CorrelationID: taskID,`（我 09:3x `grep -n` 复认）＝票面 AC#1 那句"两枚同时活着、correlation id 不同"的前提**在盘上不成立**`
  **我的判**：**已不成立**（载具前置**已做**：259-r4 落地、票 259 `AC#4` 翻勾、票改名 `259-grant-binding-identity-ruler-holes-done.md`；未勾 0／已勾 6 现量）。注：该句后段"写面被 `255-r2` 占着"是当时时点信息，不构成今日状态。
- 两句均为票面**旧时点文字、未回改**（票 `:50` 起 §5／`:59`–`:63` 本票终态段原样保留；A737 在其节内引它们并登记归口，未改一字）。

## 附：本程用到的台账节（行号同发现取，`docs/reports/pending-and-issues.md`）

- A559（`:11088` 节头）；A562（`:11111` 节头，票 259 `AC#1` 裁 ⓐ）；A709（`:13905`）；A712（`:13966`）；A714（`:14009`）；A715（`:14023`）；A717（`:14057`）；A734（`:14373`）；A735（`:14382`）；A736（`:14389`）；A737（`:14395`）；A739（`:14411`）。
