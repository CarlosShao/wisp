# 309-r1 · `AC#3` 文案一致性（`internal/tools/bridge.go`，只改行号与枚数）

判据：⛔ 改任何**契约声称**文字（"the TEXT loop already defers"／G1/G1b/G5 条款名／ticket 151 归属／"no test will tell you so"等一律原样保留）。只把漂了的行号换内容锚、把裸枚数换成带尺的现量。改动落在 `bridge.go:230/:252` 之外的两处 loop.go 引用与 "Today no production code…" 那句（`:230`、`:252` 属 `AC#3` 射程之外，⛔ 动，已确认未动）。

## 三处编辑（改前 ⇒ 改后）
1. **`bridge.go:861` 引 `loop.go:366`（＝U-5，本轮现读）**
   - 现读结论：**漂了**。票面/文档声称"TEXT loop already defers (…:366)"，但 `:366` 现落在 journal 块；真正 `defer revokeAdmission()`（AdmitTask 的 revoke）在 `loop.go:379`。
   - 改法：`…already defers (internal/agent/loop.go:366), so a` → `…already defers (in internal/agent/loop.go the "defer revokeAdmission()" line that runs right after the AdmitTask hook), so a`。**认代码本身**、⛔ 绝对行号；"already defers"这个契约声称未动。
2. **`bridge.go:871` 引 `loop.go:332/:321`（Run/RunAsync→newTaskID）**
   - 现读：**两枚都漂**。现量 `RunAsync` 起 `:328`（铸造 `:329` `id := newTaskID()`）、`Run` 起 `:339`（`:340` `return l.run(ctx, newTaskID(), input)`）。
   - 改法：换内容锚 `RunAsync's "id := newTaskID()" and Run's "return l.run(ctx, newTaskID(), input)", both in internal/agent/loop.go`。
3. **`bridge.go:877` "Today no production code dispatches on the bridge except loop.go, so being"**
   - 票面 AC#0/本票：该句是**会被一行 commit 缩短的现读句**，⛔ 当永久保证。
   - 改法：带尺重述为"as measured today (a present count, not a permanent promise) … Ruler: `grep -rn 'ToolRequest{' internal cmd --include=*.go` minus `_test.go` yields the single `internal/agent/loop.go` point; `Test309ToolRequestProductionConstructionNail` pins points outside internal/agent to 0. So being…"。枚数写明是"现量"；`so being` 与后句 `such a caller…` 衔接原样；G1/G1b/G5 那句⛔ 动。

## 尺（改后自证）
- AC#3 点名的两枚绝对号（`:861` 的 `loop.go:366`、`:871` 的 `loop.go:332/:321`）已全部换内容锚 ⇒ `grep -n 'loop\.go:[0-9]' internal/tools/bridge.go` 里**不再出现** `:366`/`:332`/`:321`。其余绝对号（`:113`/`:114` 的 `loop.go:563`/`696-700`、`:230` 的 `783-806`、`:252` 的 `692-694`、`:1084` 的 `loop.go:585`）均**非 AC#3 对象**、逐字保留、⛔ 动（实测这些仍在，正是"我没越射程"的证据）。
- 测试标识符未被换行折断：`grep -n 'Test309ToolRequestProductionConstructionNail\|Test309AssessorForEmptyTaskIDDropsTaintDetector' internal/tools/bridge.go` ⇒ 命中（`AC#2`/`AC#1` 名各出现，均整行）。
- 行为零改动：`git diff --numstat -- internal/tools/bridge.go`＝`32 3`（纯注释增删）；`go build ./...`/`go vet`/全量测试 rc=0（见 `40-gates-and-rosters.md`）。`AC#3` 闭合。
