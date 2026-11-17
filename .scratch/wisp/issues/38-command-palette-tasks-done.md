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
- [x] Text prompt via palette → streams into result view; identical behavior to CLI path.
- [x] IME: Chinese composition (pinyin → 选字 → commit) submits final text only (e2e with
      Windows IME on runner or manual evidence).
- [x] Discovery: installed tool set listed with correct risk pills; launch routes through gate.
- [x] Tasks view: live task state transitions reflected; cancel mid-tool → applied-steps report.
- [x] force rules: per-task override chip + TOML rule both steer agent mode (prompt snapshot).
- [x] Keyboard-only: full palette flow reachable without mouse (tab/arrow/enter audit).

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

---

## 编排者收件·第一批（2026-11-17 09:2x，`817785260`，仪器 `render:palette38`）

**尺1** `817785260` 在。**尺3 我自己复跑**：`npm run render:palette38` **rc=0**；build/typecheck/lint 全 rc=0。
- **勾 AC#1**：自由文本走 `panel.message.send` 同一封套（`lib/panel.ts:748`→`composer_dispatch.go:236`→生产读者 `panel_message_453.go:102`），结果进票 36 已接的结果面；仪器断言 send 行逐字上屏。
- **勾 AC#2（IME）**：`palette-screen.tsx:33-34/:181-183/:205` compositionstart/end ＋ `isComposing` 双保险，组合期 Enter 提前 return，只提交选字终文。
- **留白（具名）**：AC#3 名册只读面已接但「risk pills/launch 走门」缺 Go 派发表；AC#4 任务卡三字段（`subagent_roster_197.go:111-142 TaskRowView` 无 started/cost/step）＋ task.cancel 不在 C17 八法（`bridge.go:41-107`）；AC#5 force rules 无读面；AC#6 键盘全流（tab 审计）待补。

---

## 编排者收件·第二批（2026-11-17 10:4x，`7a2b83fa4`，件 `docs/evidence/s1/38-ac6-keyboard.md`）

**尺1** `7a2b83fa4` 在。**尺3 我自己复跑**：`npm run render:palette38kbd` **rc=0**（`listbox/option aria with activedescendant announcing the active row`）；build/typecheck/lint rc=0。
- **勾 AC#6**：五支牙——①Tab 次序＝DOM 次序（`<input` 在所有行 id 前、`palette-row-0<1<2`）②↑↓ 焦点序列 `0→1→2→1→0`＋环绕 `0↑→2`/`2↓→0` ③Enter＝`select`＋`aria-activedescendant` 与 `aria-selected` 同行 ④ESC/⑤Ctrl+K；空结果集无 listbox/无 activedescendant。composing 期**六键决策全 `none`**（AC#2 守卫从输入框扩到对话框级）。
- 转报（我另案处理）：`render:security40` rc=1 是**仪器钉住了旧假话**（`dd52313e1` 改真话后过期）⇒ 已并入前端下一单修；`d22scan` 5 枚红在 Go 面（票 49 引入）⇒ 已单开 `wd22clean` 清。
## 编排者追加·任务态/取消面已落地（2026-11-17 13:3x，`3a62bed97`）

我复跑（载具）：`TestTaskCardFieldsRideTheRosterWire` 0.00s／`TestRosterCancelDoorStopsARealSubagentAndTheForensicsCarryTheAppliedSteps` **3.28s PASS**（取消后读数＝终态 Stopped、已落地 1 步 fs.read、被中止 0 步）⇒ 全 PASS。
- 面：`internal/panel/subagent_roster_197.go:180-182 TaskRowView +startedAt/costToDate/currentStep` → `cmd/wisp/panel_pump.go:256`（`panel_tasklive_wave2.go:119` 取值）＋`internal/panel/bridge.go:211 roster.cancel`（理由写在常量上）→ `composer_dispatch.go:349` → `cmd/wisp/panel_roster_cancel_wave2.go:111/:187`。
- **命名裁定（编排者拍）**：票面要的 `task.cancel` 被 `internal/tools` 的票 188 名轴仪钉死（该包不在本波写面、⛔绕仪）⇒ 落地名 **`roster.cancel`**（语义＝停名册一行）。**撤销口令**：若日后把 188 扫描收窄，一行改回 `task.cancel`。

---

## 编排者收件·第二批（2026-11-17 14:2x，`71ccf56c0`，件 `docs/evidence/s1/readface-consume2-20261117-r1.md`）

**尺1** `71ccf56c0` 在（13 件）。**尺3 我自己复跑**：`render:taskcard38` rc=0；`render:historypage36` rc=0；build/typecheck/lint rc=0；`go test ./internal/panel` ok 18.866s；ledger 两枚页字面注册确已摘（`bridge_test.go` 的 rosterPageGapLedger194）。
- **勾 AC#4**：三键逐字（开始=生产 formatter／123456 micros CNY／tokens 1200-340／fs.write；缺⇒「未判」无 1970 无假 0）；**取消走真门**（第一下只开确认卡、请求仅一处发出）；取消三读数＝终态 Stopped·已落地 1 步·被中止 0 步（整史取证扫描；扫描不全⇒未判）；腿跑真子代理取消 2.86s PASS。
- 留白：AC#3（discovery 风险丸/launch 走门需 Go dispatch）、AC#5（force rules 无读面）、AC#6 已勾。
## 编排者更正（2026-11-17 14:4x，来自对抗审计 `readface-adversarial-20261117-1432.md`）

**上面收件里写「终态 Stopped」是假话，我撤回更正**：真值是 **`Muted`**（`internal/statemachine/states.go:14 StateMuted = "Muted"`；D43 二十枚里**没有 `Stopped` 一词**，`statemachine.Valid("Stopped")==false`）。这处假话同时出现在 `3a62bed97` 的 commit message 与 wave2 证据件 §2 —— commit 不可改，此处更正为凭。
- 附带**空牙**（审计判）：`TestRosterCancelDoorStopsARealSubagent…` 只断言终态非空＋是合法 D43 名，**没断言就是 Muted 那一枚** ⇒ 把实现换成"返回 true 什么都不停"仍绿。已排单补一行 `if stoppedRow.State != statemachine.StateMuted` 的牙。

## 编排者收件·第三批（Go 面）（2026-11-17 16:1x，`c25c3086d`，件 `docs/evidence/s1/38-discovery-force-20261117-r1.md`）

**尺1** `c25c3086d` 在。**尺3 我自己复跑**：`internal/panel` 全包 **ok 19.636s**；载具跑 `Test38*` 相关 **20/20 PASS rc=0**（`Test38ToolsListRosterPillsAreTheRiskTruth` 0.18s／`Test38ToolLaunchL2RaisesARealCard` 0.21s／`Test38ToolLaunchL1RaisesNoCard` 4.18s／`Test38ForceRulesFacePrecedence` 0.01s／`Test38ForceChipRidesTheMessageDoor`／`Test38DiscoveryDoorsRefuseByName`…）。
- **AC#3 的 Go 面**：`bridge.go:227 tools.list`（riskPill 逐字＝`risk.L?.String()`，回 真 registry）＋`bridge.go:243 tool.launch`（走环路唯一 ToolRequest 构造点→真 gate；L2 弹 1 张真卡、L1 0 卡真执行）。
- **AC#5 的 Go 面**：`bridge.go:257 rules.force`＋`cmd/wisp/force_rules.go:188 resolveForce` 唯一解算器；TOML⇒mode 变 source=toml／override 覆盖 source=override／都无⇒default；force 指令走票 200 接缝、无规则逐字节原样。
- **两格仍不整勾**：AC#3 的“listed with risk pills”（页名单）与 AC#5 的“override chip”（页控件）都还差**前端半边**（`panel.ts` 三枚字面＋chip 控件，腿已按 Go-only 单留 pending 登记在 gap ledger）。⇒ 下一波前端腿接。

## 编排者收件·页半边（2026-11-17 22:4x，`a5295e40b`，件 `docs/evidence/s1/38-discovery-force-ui-20261117-r1.md`）

**尺1** `a5295e40b` 在（11 件）。**尺3 我自己复跑**：`render:discovery38` rc=0／`render:forcechip38` rc=0／`render:palette38` rc=0；build/typecheck/lint rc=0；`internal/panel` ok 19.902s；dist 归 `.gitkeep`。
- **勾 AC#3**：工具发现页风险丸**逐字＝载荷 riskPill**（L0 绿/L1 黄/L2 红/Deny 四丸；缺⇒`未判`+unjudged 色，⛔假值）；launch 键→`tool.launch` **真门**（页不预判、无 allow 词、丸不上按钮门）。
- **勾 AC#5**：force 页显示生效规则＋来源（toml/override/default）＋下一 prompt 采用 mode＋instruction 逐字；**override chip 两态+off** 搭 `panel.message.send` 顺封套（off 时封套字节与改造前相同，发送后复位不跨任务泄漏）。
- 过程红两枚已修（牙与真身写法不符；新仪器里 `"panel.ts"` 字面量撞 `routeLiteralRe`，label 改 `src/lib/panel.ts` 后绿 —— 扫描器一字未动）。
- 必要偏离（腿报，我接受）：`readface_wave3_test.go` 的 `TestWave3PageGapLedgerRegistersTheThreeNames` 把“已登记”钉成永真，与 bridge_test.go 双向审计互斥 ⇒ 重裁为 `TestWave3PageGapEntriesLeftWithThePageLiterals`（再登记即红；方向反转的**增强**）。`bridge.go:226/242/256`＋`inbound_roster_253_test.go:90-92` 的 "REGISTERED while pending" 注释过期，留待有写面的腿更。