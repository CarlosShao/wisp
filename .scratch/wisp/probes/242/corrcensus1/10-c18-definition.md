# 242-corrcensus-1 · ① C18 契约与 SPEC 侧 correlation 定义（逐字抄）

锚：HEAD `8dbae00`（8dbae00c4948939c931f1162fe9d5d13361d8ce1，2026-10-08 19:18 +0800，分支 dev）。起手锚见 `00-anchor.md`。
取数命令与 rc（全部现取，无一接 `2>/dev/null`）：
- `grep -niE 'correlation' docs/PLAN.md` rc=0（命中 10 行）
- `sed -n '1355,1375p' docs/PLAN.md` rc=0
- `git grep -niE 'correlation' HEAD -- 'docs/specs/*.md'` rc=0（命中 11 行）
本件为**只读**产出：PLAN.md / docs/specs/** 一字未动（射程判断，非内容引用改动）。

---

## 1. 权威定义：C18（docs/PLAN.md:1368，整行逐字）

> | **C18** | **`ApprovalQueue`** | 全局 FIFO；每项含 `correlationId` / 所属任务标识 / 工具名 / **完整参数** / 风险级 / 请求理由。**队头单显 + 深度计数**；**超时 = 300s，一律判拒绝**（第四轮补数值，原文无数值 → agent 会自己拍一个或写成无限等待，见 §16.9 第 4 条）；**超时前 30s 醒目提示**；拒绝后**任务 root ctx 不取消，可一键重放**；宿主不可达时 **fail-closed**（学 DSH 的 `'unavailable'`）；**「允许」决策只接受原生侧来源（D33/F2，面板来源直接拒绝）** | D31, **D32/F2** |

照抄语义要点（不改写）：
- `correlationId` 是**队列每项的一个字段**，与"所属任务标识"是**并列的两个字段**；该行**没有**说两者必须相等。
- C18 行内不出现"路由"字样；"回复按 correlationId 路由"的语义在相邻条文里（下节）。

## 2. 相关权威条文（逐字；行号与短语同一发现）

- docs/PLAN.md:1134-1135（D31 ①「全局审批 FIFO 队列」块内）：
  >    每项携带 correlationId + 所属任务标识
  >    回复按 ID 路由 → 用户的点击不可能落到别的请求上
- docs/PLAN.md:1367（C17 `PanelBridge` 行，节选）：
  > **回复必须按 correlationId 路由**；前端**必须无状态**（WebView 销毁后一切从 Go 侧重读）
- docs/PLAN.md:1377（C27 `PanelManager` 行，节选）：
  > 多任务共用，按 correlationId 分区渲染
- docs/PLAN.md:1430（S7 切片行，节选该分句）：
  > **C18 审批队列**：并发任务各自请求确认时**按 correlationId 正确路由**、队头单显 + 深度计数、**超时判拒绝**；
- docs/PLAN.md:2406（§F2 面板 XSS 第三层威胁句）：
  > 一次 XSS 即可执行 `PanelBridge.invoke('approval.decide', {correlationId, allow: true})`
- docs/PLAN.md:2703（D35 `tool_call` 表行，节选）：
  > `tool_call` | `id` PK, `task_id` FK, `seq`, `tool`, `args_json`, **`risk_level`**, **`decision`**, `decided_at`, `started_at`, `ended_at`, `outcome`, `error_class`, **`correlation_id`**, `grant_id` | `task_id,seq` · `correlation_id` | **30 天** | **取证记录**：…
- docs/PLAN.md:2971（D42 可靠性缺口登记 ③）：
  > ③ correlationId 路由 + 事件推送的背压与合并（D38(d)）
- docs/PLAN.md:3025（§关键时序第 3 条，节选）：
  > A 被原生侧批准 → 按 `correlationId` 路由回 A → 出队 → 面板（C27 单窗口）切到 B 的请求 →
- docs/PLAN.md:3485（早期 UI「工具调用（展开）」行，节选）：
  > 结果摘要 + 耗时 + `correlationId`（供 D31 取证）

## 3. SPEC 侧（`git grep -niE 'correlation' HEAD -- 'docs/specs/*.md'` rc=0）

- docs/specs/SPEC-06-security-gatekeeping.md:102（§7 C18 ApprovalQueue 与 D31 并发审批）：
  > - 全局 FIFO；每项含 `correlationId`/所属任务/工具名/**完整参数**/风险级/请求理由。
- docs/specs/SPEC-06-security-gatekeeping.md:106：
  > - 回复按 correlationId 路由（点击不可能落到别的请求上）。
- docs/specs/SPEC-06-security-gatekeeping.md:151（§11 测试决策）：
  > - 审批队列：并发两任务同时 L2 → correlationId 路由正确、队头单显、超时判拒绝、重放可续。
- docs/specs/SPEC-02-data-storage.md:85（DDL 行内注释）：
  >   correlation_id TEXT NOT NULL,          -- C18 路由与取证
- docs/specs/SPEC-02-data-storage.md:89：
  > CREATE INDEX idx_tool_call_corr ON tool_call(correlation_id);
- docs/specs/SPEC-08-ui-ball-panel.md:154：`- 多任务共用单窗口，按 correlationId 分区渲染。`
- docs/specs/SPEC-08-ui-ball-panel.md:159：`correlationId 路由；每方法标注 capability 与是否需原生侧授权；未列出方法名 → 拒绝并记日志。`
- docs/specs/SPEC-08-ui-ball-panel.md:193（工具调用 chip 行节选）：`展开 = 着色 JSON + 耗时 + correlationId`
- docs/specs/SPEC-08-ui-ball-panel.md:235：`- C17：correlationId 路由测试（并发 10 个请求乱序回复）、方法白名单拒绝、resync 后前端状态与`
- docs/specs/SPEC-12-roadmap-governance.md:23（S7 行节选）：`C18：并发确认按 correlationId 路由、队头单显+深度计数、超时判拒绝；`

射程结论：docs/specs 里**没有**比 C18 表更细的"correlation 是什么"的定义段；上列即全部命中（rc=0，射程＝`docs/specs/*.md`）。

## 4. 产码注释里的实现口径（非契约；原句逐字，供裁决者对照）

- internal/agent/approval/queue.go:36-41（`names` 字段注释）：
  > // names are the OTHER strings this one card can legitimately be addressed
  > // by: the correlation id as it arrived before the queue re-stamped it, and
  > // the task id the host keys its own bookkeeping on (internal/tools/bridge.go routes
  > // cancels by orDefault(CorrelationID, TaskID)). Ticket 87: a card that is
  > // on screen is addressable by whatever the display was built from, and a
  > // reply carrying one of those names is not a reply for a different request.
- internal/agent/approval/gate.go:516-522（L2 `PendingApproval` 内）：
  > 	corr := it.Corr
  > 	// D31 bookkeeping keys on the INCOMING correlation id, because that is the
  > 	// value the bridge's cancel bus looks a started call up by; it.Corr is the
  > 	// card's address in the queue, and the two are not the same string once the
  > 	// queue has re-issued one. Captured before d is re-stamped below.
  > 	incoming := orDefaultText(d.CorrelationID, d.TaskID)
  > 	d.CorrelationID = corr
- internal/agent/approval/window_read.go:35-40（`L1Window` 头注释）：
  > // L1Window is one confirmation window an in-process observer sees as waiting for
  > // a veto right now. CorrelationID is the id the window is keyed by in g.windows
  > // (which is the task id whenever the decision carried no correlation of its own,
  > // per orDefaultText at gate.go:275); TaskID is the task the call belongs to, so a
  > // reader that pairs rows by task id can join without guessing at the correlation.
  （注释内 `gate.go:275` 与 HEAD 现量 281 差 6 行：注释行号已漂，本程只登记。）
- internal/panel/subagent_roster_197.go:56-58（`ParentTaskID` 字段注释）：
  > // The entity leg files the deriving call's correlation id, and the agent
  > // loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),
  > // so on the production path this reads as the parent's task id.
  （注释内 `loop.go:647` 与 HEAD 现量 654 差 7 行：注释行号已漂，本程只登记。）
- internal/tools/task.go:702-705（`task.cancel` caller 注释）：
  > // The caller's identity is the host's, not an argument: CorrelationID is what
  > // the bridge stamps from the call it dispatched, so a model cannot name whose
  > // child it is by editing parameters. Same source task.spawn uses to fill in
  > // ParentTaskID, which is what makes the two sides agree by construction.
- internal/tools/cancel.go:53-54（访问器注释）：
  > // CorrelationID returns the id of the call the running tool belongs to, or ""
  > // when the tool was invoked outside the bridge.

## 5. 一句话定义（本程读数，非裁定）

C18 表义：`correlationId` ＝**审批队列项的字段与答复路由键**，与"所属任务标识"并列且未要求同值；SPEC-06:106 补"回复按它路由（点击不可能落到别的请求上）"；实现注释（queue.go:36-41、gate.go:516-522）进一步把它区分为"队列入队前的 incoming 值"与"队列重发后的卡片地址 it.Corr"两个可能不同的串。
