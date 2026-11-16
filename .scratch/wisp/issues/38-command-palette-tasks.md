# 38 — Command palette (cmdk) + tasks view: text input channel, force rules, discovery

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: palette + text entry; B: tasks list + tools directory)
**Spec refs:** SPEC-08 §5.2/§7, D9 (four roles), D11③ force rules, D34 task.list/cancel, S5

## What to build
The command palette as D9's four-role surface: searchable text entry fallback, L2 confirm
reliable channel, plugin/tool discovery with risk visibility, and the concurrent-tasks list
view (task.list/cancel tools backing it).

## Key constraints
- Palette = cmdk: fuzzy commands (open panel pages, mute, settings, tools) + free-text prompt
  entry (keyboard-first, ⏎ submits to agent as text task — the D9 noise-environment fallback);
  reachable by hotkey + tray + ball click; Chinese IME composition handled correctly.
- Tool discovery: lists resident + third-party tools with name/description/risk/capability
  ("我到底装了什么" transparency role); invoking a tool from palette goes through the SAME host
  bridge (no bypass).
- Tasks view: task cards (id, state, started, cost-to-date, current step) + cancel button →
  `task.cancel` (L1); blocked tasks show "waiting on what" (D31); queued-behind-pathlock tasks
  render their blocker (D43 #20/#39 groundwork — real multi-task at 47).
- force_tool/force_chat TOML rules (D11③): palette offers per-conversation override chips +
  config file support; overrides surface in agent context assembly.
- All palette state via bridge resync only (stateless); palette → result handoff switches to
  the result view (36) with streaming continuation.
- Cancel semantics: cancelling from palette → same ctx cancellation + applied-steps report
  (20) if mid-execution.

## Out of scope
- Multi-task concurrency engine (47); voice path; config editing (39).

## Acceptance criteria
- [ ] Text prompt via palette → streams into result view; identical behavior to CLI path.
- [ ] IME: Chinese composition (pinyin → 选字 → commit) submits final text only (e2e with
      Windows IME on runner or manual evidence).
- [ ] Discovery: installed tool set listed with correct risk pills; launch routes through gate.
- [ ] Tasks view: live task state transitions reflected; cancel mid-tool → applied-steps report.
- [ ] force rules: per-task override chip + TOML rule both steer agent mode (prompt snapshot).
- [ ] Keyboard-only: full palette flow reachable without mouse (tab/arrow/enter audit).

## Progress log (append-only, newest last)

> **2026-09-26 进度追记（owner 把前端交 ZCode 直管后的可视壳进度，只追加不勾框）**：
> 本票的可视面已按 beautiful-ui 主题收编第四代重画，演示形态在 `?harness=1`
> 应用全观（左栏三段式/会话更改卡/右栏审查·终端·浏览器）与 `?harness=2` 陈列室。
> **本票 AC 的真数据 e2e 格仍等票 145/35/33 的字段与宿主，一格未勾、不由本程
> 勾**；可视壳的 commit 链见票 77 进度追记（`424ac84`…`5e23d99`）与
> `docs/reports/frontend-session-log-zcode.md` §1–§7。

> **2026-11-16 23:0x 腿（agent-99，GLM 舰队五票顺序跑第 3 票；铁律 4 前端归本编队）**：
> 命令面板从死壳接成三面真面（`frontend/src/App.tsx` + `palette-screen.tsx`）：
> ① 自由文本行 → `panel.message.send`（composer 同一封套，`lib/panel.ts:748`；
>   生产读者 `cmd/wisp/panel_message_453.go:102` → relay → admitTask → agent.Loop），
>   结果照旧走 `snapshot.results` 流进 36 的结果面——无第二通道、无旁路；
> ② 命令名册（票 213 交付 3）第一次有 UI 读者：`snapshot.commands`
>   （读面 `internal/panel/commands_roster_213.go:45-51`，装配 `pump.go:466-476`）
>   按查询过滤展示 name/summary/approval/door/args，行只读不执行（名册只承载
>   不执行，`commands_roster_213.go:21-26`）；
> ③ 中文 IME 组合守卫（AC#2）：compositionstart/end + `isComposing` 双保险，
>   组合期 Enter 交给输入法，只提交选字后的终文；输入框挂载即聚焦（AC#6）。
> **未接线 dormancy 全部具名（面板 notes 区，保持未接线、不改"即时"）**：
> 命令执行（C17 名册 `internal/panel/bridge.go:41-107` 八法之外无执行通路）、
> 取消任务（无 task.cancel 路由；要接须 Go 侧在 `internal/panel/composer_dispatch.go`
> 派发表加处理器——注意 `internal/tools` 票 188 task-axis 扫描会红 task.* 路由名，
> 归 Go 腿）、模式强制 force_tool/force_chat（**DEFERRED(D11-3)**，
> `internal/agent/control.go:13-15`：规则表未实现，快照无读取面、C17 无写入路由）、
> 命令名册无 reader 形（快照无 commands 键 → 具名 `PumpSources.Commands`，
> `internal/panel/pump.go:221`）。
> 任务视图半格：活状态转移走任务监控浮层既有真接（`App.tsx:103` monitorAgents ←
> `snapshot.tasks`，读面 `internal/panel/subagent_roster_197.go:111-183`，等批准 join、
> 流截断三元组照线显示）；**卡片三字段 started/cost-to-date/current step 线上无此字段**
> （`TaskRowView:111-142` 不带），零造值，要加字段归 Go 腿（申请：TaskRowView 增
> startedAt/costToDate/currentStep 三字段 + 名册 reader 填充，`subagent_roster_197.go`）。
> 新牙 `npm run render:palette38`（send 行逐字回显、名册四事实真数据往返、只读断言、
> IME/聚焦/方向键源级断言、dormancy 具名断言、started/cost/step 零造值、双线形 App 不炸）
> OK 0 失败；build/typecheck/lint rc=0；render 仪器回归（nav/approval37/history36/
> shell77/receiver/stream/composer）全绿；`render:l2` 裸跑 rc=2 是其自带的"需 JSON 参数"
> usage 退出（既有行为，与本腿无关）。d22scan 对 frontend/ 零命中；仅有的两枚红
> （`cmd/wisp/approval_reply.go:379/464`）在本波 cmd/wisp 另一腿的未提交 +186 行里，
> 非本腿写面，已如实转报。AC 勾框一律留给非实现者。
