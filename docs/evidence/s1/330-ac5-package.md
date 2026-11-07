# 330-ac5-package — 票 330 AC#5（cmd/wisp 整包台面与颜色）补做读数

- 腿：330-ac5 台面腿（写面＝`.scratch/wisp/probes/330/ac5/**`＋本件＋票尾 bullet 一行；⛔ 碰 `cmd/wisp/**` 产码与 `internal/**`；⛔ 翻框；⛔ push）
- 日期：2026-11-07 12:4x-13:1x +08
- 锚：`58eca3738fccda15662be5e1b640290aab776e4f`（`58eca3738`，起跑时 HEAD 现认；票 330 的 cmd/wisp 写面 `edad424fe` 已在锚内）
- 台面：两棵仓外 clone 导出树（⛔ 共享工作树），配方与全部量尺见 `.scratch/wisp/probes/330/ac5/setup-gauges.md`

## 1. 两发四数（口径：RUN=全部 `^=== RUN` 含子测试；PASS/FAIL/SKIP_TOP=仅零缩进判决行）

| 发 | 树 | rc | 耗时 | RUN | PASS_TOP | FAIL_TOP | SKIP_TOP | PANIC |
|---|---|---|---|---|---|---|---|---|
| 终态 | `D:\tmp\wisp\20261107\330ac5\export` | 1 | 292.281s | 297 | 229 | 0 | 0 | 0 |
| 起手 | `D:\tmp\wisp\20261107\330ac5\export-base` | 1 | 291.021s | 290 | 222 | 0 | 0 | 0 |

读数件：`logs/330ac5-final-roster.log`、`logs/330ac5-base-roster.log`（解析尺 `parse-roster.sh` 输出）；原始整发 `logs/330ac5-final-full.log`（360,175 字节）、`logs/330ac5-base-full.log`（330,071 字节）。

## 2. 红名册作差

- 红名册（顶层 `--- FAIL`）：终态 **0 枚** ⇄ 起手 **0 枚** ⇒ **作差＝0**。任何缩进的 `--- FAIL` 行两发均为 0（grep 实证）。
- `panic:` 行两发均 0（日志里 7 处 "panic" 命中全是 prose：`panicked=<nil>`、`3/3 panic on a planted thread` 等，非测试崩溃输出）。

## 3. 死点（RUN 无判决行）逐枚具名

- 终态发恰 1 枚：`TestBallPanelGesturesReachThePanelThread`（`cmd/wisp/panel_resident_windows_test.go`，真窗族；日志尾＝该枚 MARK299 wndproc 流后直接包级 `FAIL … 292.281s`，无任何 `--- FAIL`）＝进程级截断，⛔ 当缺陷，只记死点。
- 起手发恰 1 枚：`TestAC13BringUpSurvivesAReusedThreadQuit`（同文件，真窗族）＝同形。
- ★如实：两发死点**不同名**（同文件内截断点非确定），且两发耗时几乎相同（292.3s / 291.0s）——截断点漂移但截断事件稳定复现。四数之差（RUN 297⇄290、PASS 229⇄222）完全由截断点漂移造成，非任何用例变色。

## 4. 截断深度（★如实，派单口径未点名这一层）

- `go test ./cmd/wisp -list '^Test'` 顶层名册＝**507** 枚（`logs/330ac5-list-507.log`）。
- 实跑到达顶层：终态 230 / 起手 223 ⇒ **未到 277 / 284 枚**（`logs/330ac5-final-unreached-277.log`、`logs/330ac5-base-unreached.log`）。
- **`Test330DockPinsTheWireNameToTheAssembledRoster` 两发都在未到名册里**（终态未到表 :46）：`plan_*` 文件序在 `panel_*` 之后，进程死在 `panel_resident_windows_test.go` 内 ⇒ 本机整包台面今天**到不了 dock 用例**，整包格拿不到 dock 的颜色 witness。dock 的颜色证据仍在定向发：330-r1 正向 rc=0（读数行 :134）；330-verdict 重放发 rc=1 红 :128 ＋ 正向发 rc=0。整包格是否因此算达成，归编排者判读；本腿⛔ 翻框。

## 5. 与 330-verdict 共享树那发的关系（只具名，⛔ 判案）

330-verdict 那发（共享树）报死点 6 枚（223 hand-edited fs／AC2 real-process／AC3 early-record×2／panel host real-window／latency percentiles）。本腿干净导出树里这些用例**全部实跑且 PASS**（如 `TestTicket223FailureSentencesAreDistinct` 起 `logs/330ac5-final-full.log:816`；AC2/AC3/latency 各枚同见 PASS 流），两发零 `--- FAIL` ⇒ 与"共享树里 cmd/wisp 未跟踪件污染读数"的派单假说相容；本腿只记现象、不下污染归因的终判。

## 6. 结论（腿判，框归编排者）

- 红名册作差＝0：**成立**（0 ⇄ 0，两发逐字同形）。
- 四数点名：已逐枚给出（§1）；两发不同是截断点漂移，非颜色漂移。
- 死点：每发恰 1 枚、已具名（§3），形状与派单预告（"每发恰 1 枚 RUN 无判决行"）一致。
- 未竟如实：整包射程被真窗族进程级截断砍到 230/507，dock 用例不在整包 witness 内（§4）。
