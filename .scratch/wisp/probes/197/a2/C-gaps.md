# 197-a2 §C — 落地缺料清单（逐块；⛔ 本腿只列料，不派活不拍板）

> 落点列＝真需要动哪棵树；⚠ 标「要 `cmd/wisp`」的块，其写面正排给 `294-r1`／`313-r1a`／`186-r1`（`git log -1` 派单存档 `8e4f6f83` 逐字），派单时须串行。
> 依赖列＝是否要等票 201／219 族，或独立。

## C1｜五格的非实现者裁决料（AC#0／AC#1／AC#2／AC#3／AC#6）

- **缺什么**：一枚 `197-v2`（非实现者）裁决表，判这五格「盘上形状是否＝票面声称」。今天 `docs/evidence/s1/` 下**只有** `197-ac5-selfapproval-v1.md`（只吃 AC#5）与 `197-blocked-carrier-r4.md`（只吃 AC#4）。
- **落点**：`docs/evidence/s1/197-*.md`（零 Go 产码）。
- **依赖**：**独立**（被验对象 b9fa815b／335b8d2b／94ea50f7／`run.go:780` 全在盘上）。
- ⚠ **要 `cmd/wisp` 的读数面**：五格要跑 `go test ./cmd/wisp ./internal/...` ⇒ 与 `294-r1`／`313-r1a` 的写面**互斥**，须排在其后（先例：`197-r1` §⑤ 让出门禁窗口给验收腿）。
- 尺：`ls docs/evidence/s1/ \| grep 197`。

## C2｜AC#6 链上那枚 goroutine 名（票 229）

- **缺什么**：收尾协程名 `subagent-finish-<id>` **不在 D38 名册**（`grep -c "subagent" internal/observe/goroutine.go`＝0）⇒ 每次 spawn 打一行「疑似泄漏」WARN，且抬高 SLO 协程总数。
- **落点**：`internal/tools/subagent_197.go:383-384`（spawn 点）——**甲形**＝把名字换成名册内已有形状（每任务前缀族 `agent-task-` 最贴）＝**可在 `internal/` 内直接做**；**乙形**＝往 D38b 表加 `subagent-finish-`＝**动冻结表＝人工批准**（须先落 `A##`）。⛔ 第三支「把 WARN 调轻／删」＝削弱仪器，票面明禁。
- **要 `cmd/wisp` 吗**：**不需要**（写面全在 `internal/tools/`）。⚠ 但票 229 自带 `AC#6 整包终态`（`go test ./cmd/wisp ./internal/...`）⇒ 验证面仍要碰 Go 测试面。
- **依赖**：**独立**（票 229 自有六格，`0 勾`）。
- 尺：`sed -n '383,384p' internal/tools/subagent_197.go`；`grep -c "subagent" internal/observe/goroutine.go`。

## C3｜AC#4 后半那格（票 220）

- **缺什么**：`blockedOnApproval` **只报 L2**（`LiveApprovals` 只走 `q.pending`）＋同一任务第二张卡因 join 用 `row.TaskID` 而卡键是 `corr` 静默读成「没卡住」；L1 短窗口住 `Gate.windows`、**全仓无枚举口**。
- **落点**：①L1 只读枚举口＝`internal/agent/approval/gate.go`（**新增导出方法** ⇒ 撞票 197 禁区「不新增方法名」，c2 §⑥ 已裁＝把它**拆进票 194 的批准面**，票 220 只做「把已有数据读对」）；②join 键改 corr／taskID 双查＝`cmd/wisp/panel_pump.go` ⇒ **要 `cmd/wisp`**。
- **依赖**：**独立**（票 220 自票、`0 勾／5 未勾`）。票面逐字「与票 219 的关系＝219 `AC#6` 早就指过这一族」⇒ **同一族但不互等**。
- 尺：`head -6 .scratch/wisp/issues/220-*.md`；`grep -c "^- \[x\]" .scratch/wisp/issues/220-*.md`。

## C4｜AC#5 那三条尺的修法（`197-v1` 的 C1／C2／C3）

- **缺什么**：C1 第七枚载具 `spawnText` 恒空；C2 `:254` 两值集合分不清「门拦的」与「权限拦的」；C3 脚注 `M1` 指错落点。
- **落点**：`cmd/wisp/subagent_selfapproval_197_test.go`（`:156`／`:315`／`:248-261`／文件末脚注）⇒ **要 `cmd/wisp`**。
- **依赖**：**独立**；⚠ 动被测文件必须同批把 **P1／M1 真形／D1** 三发当正控重跑（票面第 56 行逐字），台件 `.scratch/wisp/probes/197/v1/` 可直接复用。
- **C4 那枚**：跨卡借证形 → **票 242 已 `-done`**（`242-…both-have-zero-rulers-done.md`）。
- 尺：`grep -n "spawnText" cmd/wisp/subagent_selfapproval_197_test.go`（现只 `:156`／`:315`，零赋值）。

## C5｜界面那一跳（⛔ 跨边界，归 owner）

- **缺什么**：`frontend/src/lib/panel.ts:129-137` 的 `PanelSnapshot` 缺 `tasks?` 与 `instructions?` 两枚 ⇒ `internal/panel/composer_test.go:48` 双向对账**今天为红**。
- **落点**：`frontend/**`（**本编队零写面、不读不引**）⇒ 由 owner 自己带给前端那枚 agent。
- **依赖**：**独立**；⛔ 不许为过尺让 Go 少发键、不许改尺或加豁免（票面第 91 行逐字）。
- 尺：`sed -n '122,140p' frontend/src/lib/panel.ts`；`grep -n "does not declare" internal/panel/*_test.go`。

## C6｜依赖关系速查（票 201／219 族 vs 独立）

| 块 | 是否要 `cmd/wisp` | 依赖 201／219 族？ |
|---|---|---|
| C1 五格裁决 | 读数面要（写面不要） | **独立** |
| C2 票 229 名字 | 否（写面在 `internal/tools/`） | **独立** |
| C3 票 220 | **是**（`panel_pump.go`） | **独立**（同族不等） |
| C4 `197-r3` 尺修 | **是** | **独立** |
| C5 界面跳 | ⛔ 归 owner | **独立** |

- 票 201 现 2 勾／5 未勾（`部分结案`），票 219 现 0 勾／6 未勾（`接单＝做`）⇒ 两枚都**未 -done**；但 `197-r1` 那枚正控已用现成 `Replies.Allow` 走通 ⇒ ticket 197 的**裁决料不欠 201／219**。
