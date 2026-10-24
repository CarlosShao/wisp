# m4-survey ★ 把权限档位放宽时，主流 harness 怎么处理

> 调研腿 `m4-survey`（只读）。对象 = `D:\work\AI\open source\` 8 项：`cline`／`deepseek-harness`／`minimax-code`／`openchamber`／`pi`／`pi-upstream`／`Step-Code`／`_survey`。
> 只读纪律：`D:\work\AI\open source\**` 一字未写、未落探针；`oc.tar.gz` 未解压（§2.8）。引文逐字取自工作树，`文件:行` 为实测行号。

## ① `_survey` 里的现成结论（有就引，不重做）

1. **dsh 的授权梯度已量化**：4 个可见值（`workspace-write`/`danger-full-access` + 代码固定 `auto` + 派生 `custom`）、**2 个真实自由度**、**只有 1 种授权 `allowed-once`**——`_survey/dsh-agent-loop-and-context.md:490-524`。
2. **授权结果只有 4 值**，`unavailable` 的定义是"缺通道"且**开不了闸**——同 `:552-570`（引 `docs/subsystems/approval.md:21`：*"A missing, non-owning, throwing, or non-conforming answerer becomes `unavailable` rather than opening the gate"*）。
3. **dsh 的 5 层**：`tools/pre-execute` waterfall → 单调 guard（只许 deny）→ approval seam → 会话事件旋钮（`sandbox/mode`+`approval/policy`）→ 工具内升级——`_survey/cross-repo-tool-catalogue.md:445-469`。
4. **"无通道"的模型侧文案**：`tool "X" requires approval, but no approval channel is available`——同 `:516`。
5. ★**_survey 没答本格**：改档位**这一下**弹不弹、弹的条件、跨入口是否一致。下面 §2.1–2.6 都是本次新读的代码。

## ② 逐 repo 五子问题证据

### 2.1 deepseek-harness（dsh）★ 唯一"给最宽档位配确认"的 repo

- **Q1 ⇒ 弹，但只按"目标档位"弹，且只在 Web 客户端**。`packages/client/ui-permission-presets/src/client/index.ts:94-104`：
  `…(option.value === FULL_ACCESS_PRESET || option.value === AUTO_REVIEW_PRESET ? { confirmation: { title: …, acknowledgeLabel: t('confirm.acknowledge'), … } } : {})`；
  `FULL_ACCESS_PRESET = 'danger-full-access'`（`presentation.ts:4`）。**中间档 `workspace-write` 不在条件里**。
  另两处同规则：`PermissionSelect.tsx:129-138` `if (id === FULL_ACCESS || id === AUTO_REVIEW) { … setConfirmation(id); return }`；`PermissionRow.tsx:81-84` `if (id === FULL_ACCESS_PRESET) { setConfirmingFullAccess(true); return }`。
- **Q4 ⇒ 无方向性**：判据是"目标档位的 bundle 会不会拿掉提示"，不是变宽/变窄。`danger-full-access`（无沙箱无审批）与 `auto`（沙箱关 + LLM 复核）弹；`read-only → workspace-write` 这种变宽**不弹**。
- **Q2 ⇒ 放宽照旧生效**（确认门只在客户端，宿主不挡）：`packages/interaction/permission-presets/src/index.ts:264-273` 宿主 `/permission <preset>` 只有"名字存在吗"一道判断便 `this.apply(...)`；`apply` 落耐久事件 `session.append('permission/preset', …)`（`:403-418`）。⇒ 无 UI = **允许 + 记日志**。
- **Q3 ⇒ 两层**。改档位 0–1 道（仅 Web 两枚目标档）；每次调用 4 道，且**无通道即拒**：`packages/core/tools/src/index.ts:1731-1736` `const approval = this.ctx.get('approval')` / `if (approval === undefined) { return { decision: { kind: 'deny', reason: … requires approval (not yet supported) } } }`；`:1763` `… but no approval channel is available`；升级 `packages/sandbox/sandbox/src/escalation.ts:205` 直接 throw 同一句。
- **Q5 ⇒ 不一致，而且写在注释里**：`client/index.ts:7-8` *"the argued path (`/permission <preset>` still switches directly)"* ⇒ 手打 `/permission danger-full-access` 绕过弹窗；处理方式 = **不处理**。

### 2.2 cline ★ 把"总开关"删掉，改成按动作类别

- **Q1 ⇒ 完全不弹**。GUI 存盘是一条裸写：`apps/vscode/src/core/controller/state/updateAutoApprovalSettings.ts:11-37`（只做 version 合并后 `setGlobalState`，无确认）；CLI 侧 `/yolo on` 同样裸写：`apps/cli/src/utils/chat-commands.ts:385-397` → `apps/cli/src/runtime/run-interactive.ts:768-772` `setToolAutoApproveGlobally(enabled)`。
- **★设计方向**：`apps/vscode/src/hosts/vscode/vscode-to-file-migration.ts:264-266` 逐字 *"The "YOLO mode" and "auto-approve all" toggles were removed: both were blanket "run without asking" switches that bypassed the per-action auto-approval settings."*
- **Q3 ⇒ 档位换成 8 个动作类别**（同一文件 `:296-306`）：`readFiles`／`readFilesExternally`／`editFiles`／`editFilesExternally`／`executeSafeCommands`／`executeAllCommands`／`useBrowser`／`useMcp`。命令还分"安全命令 / 全部命令"两级。
- **Q2 ⇒ 拒（fail-closed）**：`apps/cli/src/runtime/interactive/approvals.ts:48-51` `if (refs.tuiToolApprover.current) { … } return { approved: false };`；hub 侧 `sdk/packages/core/src/hub/server/handlers/approval-handlers.ts:17-21` `if (state?.interactive === false) { return { approved: false, reason: "Tool approval requires an interactive session, but this session is non-interactive." } }`。
- **Q4 ⇒ 有一条方向性规则，但在 plan→act**：`apps/cli/src/runtime/interactive/mode.ts:39-40` `"…only call this after the user has explicitly approved the plan in a message sent AFTER you presented it … never treat the original task request as approval."`（变能力更强的一侧要人先点头；回到 plan 不需要）。
- **Q5 ⇒ 一致**：GUI RPC 与 CLI 写同一个 `autoApprovalSettings`（`updateSettingsCli.ts:62-80`），两条路都无门。

### 2.3 minimax-code ★ 三档与 Wisp 逐字同构

- **档位**：`packages/tui/src/tui/features/interaction/permission-mode-picker.ts:15-19` 逐字 `default: 'Confirm sensitive actions'`／`auto: 'Ask only when risk is high'`／`bypassPermissions: 'Run without confirmation'`（`MINIMAX_CODE_PERMISSION_MODES`，`packages/tui/src/application/permission-mode.ts:1`）。
- **Q1 ⇒ 不弹，回车即生效**：`permission-mode-picker.ts:55-58` `if (matchesKey(data, 'enter')) { const mode = MINIMAX_CODE_PERMISSION_MODES[this.selectedIndex]; if (mode) this.onSelect(mode); }`。唯一的"提示"是红色标记 `:117` `mode === 'bypassPermissions' ? chalk.hex(colors.error)('⚠') : ''`。
- **Q2 ⇒ 一律 settle 成 deny**：`packages/local-runtime/src/api/local-permission-approval-service.ts:212,229,241` 与 `:336-341` 的 abort 分支全部 `settle(record,'deny')`。**中间档缺分类器时降级为"问"**：`permission-core.ts:540-542` `if (mode === 'auto') return { interaction: 'ask', classifier: 'allowed', … }`，且无任何 allow 证据时 `:418-435` 落到 ask；无托管网关时 *"fail open to ask-user"*（`in-memory-cloud-gateway-client.ts:9`）。
- **Q3 ⇒ 两层**：档位（无门）+ 每次调用一个判定引擎（规则 / 安全分类器 / bypass-immune 硬拒集，`types.ts:8-18`、`engine.ts:437-442`）。
- **Q5 ⇒ 不一致**：ACP 客户端可直接 `setPermissionMode`，服务端无门（`packages/tui/src/acp/agent.ts:969-973`）。

### 2.4 openchamber ★ "档位不可用时把它从可选集合里去掉"

- **档位**：`packages/web/server/lib/permission-auto-accept/modes.js:1-18` `ask`／`safety`（安全网否决才问）／`auto`（全收）。
- **Q1 ⇒ 一键循环，不弹**：`packages/ui/src/components/chat/composer/ui/PermissionAutoAcceptButton.tsx:54-77`（`onClick={handlePermissionModeCycle}`）+ `packages/ui/src/stores/utils/permissionAutoAccept.ts:26-31` `nextPermissionMode`。
- **Q2 ⇒ 三种处理叠着用**：① 中间档跑不动就**从循环里消失**——`permissionAutoAccept.ts:22-31` `displayedPermissionMode`（`safety` 无分类器时按 `ask` 显示、循环跳过 `safety`）；② 分类器不给答案 ⇒ **挂起等人**：`runtime.js:239-242` `if (verdict?.action !== 'accept') return 'held';`，`modes.js:5-7` *"without an answer from Jev the request waits"*；③ 没有条目 ⇒ **默认 ask**：`permissionAutoAccept.ts:93,100`。
- **Q3 ⇒ 两层，且"改档位"本身就是一次授权事件**：`runtime.js:125-137` 改档后 `if (isAutoAnsweringMode(next)) await reconcilePending(...)`，`reconcilePending`（`:294-323`）把**已在队列里的请求**逐条 `replyOnce` → `{ decision: 'once' }`（`:244-249`）。⇒ 放宽会**追溯**清空挂起项。
- **Q5 ⇒ 一致**：策略在服务端（`policySnapshotFromWire`，`permissionAutoAccept.ts:43-56`），GUI / VS Code 桥 / 旧 on-off 客户端都读同一份；旧布尔值映射 `modes.js:20-24`。

### 2.5 Step-Code ★ 唯一把"无 UI 时怎么办"写进档位本体的 repo（且与本仓冲突最相关）

- **档位**：`packages/coding-agent/src/step/permissions.ts:37-70` 四个 preset，每项带 `mode` + **`nonInteractiveApproval`**：`ask`=confirm/deny、`read-only`=strict/deny、`bypass`=auto/**allow**、`autopilot`=auto/allow+autoResume。
- **Q1 ⇒ 不弹**。改档只有启动参数/持久化：`src/cli/args.ts:205-212`（`--approval-mode`）、`:525` 帮助文本 *"`--non-interactive-approval <mode>` Fallback without a UI: allow or deny"*；持久化 `settings-manager.ts:157`。
- **Q2 ⇒ 显式"降级到更严的档"**：`permissions.ts:518-521`
  `private effectiveState(cause) { if (this.state.mode !== "auto" || !(cause.refused || cause.unconfigured)) return this.state; return { ...this.state, mode: "confirm", nonInteractiveApproval: "deny" }; }`
  注释 `:505-517` 逐字：*"An unattended cause needs the mode downgraded rather than just the no-UI fallback consulted, because `auto` decides every call as `allow` and returns before the fallback is reached. Without this, `--non-interactive-approval deny` silently does nothing."*
  随后每调用一道 `:535-545`：无 UI 时普通确认**只有显式 `allow` 才放**，hazardous/分析不全**永远拒**并给出指名 flag 的理由（`:562-597`）。
- **★默认即放宽、但只在有 UI 时才成立**：`:252-258` *"Default to bypass … a run with no UI refuses the defaulted policy outright"*，对应 `unattendedCause` 的 `unconfigured = this.state.defaulted === true`（`:500-504`）。⇒ **隐式放宽无通道 ⇒ 拒；显式放宽无通道 ⇒ 照用 + 动作层 fail-closed**。
- **Q4 ⇒ 无方向性**，只有"显式 vs 缺省"这条轴。
- **Q5 ⇒ `ctx.hasUI` 一条线贯通**（pi 的 `hasUI`，`packages/coding-agent/src/core/extensions/runner.ts:578`）。

### 2.6 pi-upstream ★ 没有档位，只有"没 UI 就保持更严的那一侧"

- **Q1 ⇒ N/A**：全仓无 permission tier、无 `yolo`、无 `--dangerously-*`；`--approve/-a` 只是"本次信任项目本地文件"（`packages/coding-agent/src/cli/args.ts:229-231`），**这条 flag 本身即同意**，不再二次确认。
- **Q2 ⇒ 保持严的那一侧**：`packages/coding-agent/src/core/project-trust.ts:77-88`
  `switch (options.defaultProjectTrust ?? "ask") { … }` / `if (!options.projectTrustContext.hasUI) { return false; }`（`false` = 不信任）；`src/cli/project-trust.ts:19-21,31-34` 的 `select/confirm` 无 UI 时返回 `undefined/false`。
- **官方样例给出同一条判据**：`packages/coding-agent/examples/extensions/permission-gate.ts:20-24`
  `if (!ctx.hasUI) { // In non-interactive mode, block by default` / `return { block: true, reason: "Dangerous command blocked (no UI for confirmation)" }; }`
- **Q3 ⇒ 一层**（per-call `tool_call` 钩子）；**Q4 ⇒ 无**；**Q5 ⇒ 一致**（同一 `hasUI`）。

### 2.7 pi（`pi.js` + `pod_setup.sh`）⛔ 没有这个概念

GPU pod 管理器（vLLM 部署），全仓无 `permission/approval/autoApprove/mode/escalat` 任一概念的实现；`yolo` 只出现在 README 的局限说明里（`README.md:30` *"OpenAI endpoints exposed to the public internet (yolo)"*）。⇒ **这个 repo 里没有"权限档位"这回事**。

### 2.8 `oc.tar.gz`

顶层 41MB 压缩包，内容为 `openchamber-main/`（`tar -tzf` 首行）。**未解压**（遵只读纪律，也不解到源目录）；openchamber 的结论全部取自已展开的 `openchamber/` 树。

## ③ 对照表（repo × 五子问题）

| repo | Q1 放宽弹不弹／条件 | Q2 无确认通道时 | Q3 粒度（层数） | Q4 方向性规则 | Q5 跨入口 |
|---|---|---|---|---|---|
| **dsh** | 弹，**只对 `danger-full-access` + `auto` 两个目标**，且只在 Web 端 | 放宽照旧生效，落 `permission/preset` 事件（允许+记日志）；**每次调用**则 deny | 2 层（改档 0–1 道；调用 4 道） | 无（按目标档位，不按方向） | **不一致**：`/permission X` 绕过弹窗，注释承认 |
| **cline** | 不弹（GUI RPC 与 `/yolo on` 都是裸写） | per-call `{approved:false}`；非交互会话明文拒绝 | 2 层（8 个动作类别 + 每次调用） | 有，但在 plan→act：act 需人先批准 | 一致（同一 state，两条路都无门） |
| **minimax-code** | 不弹（回车即生效，仅红 ⚠ 标记） | 一律 settle deny；中间档缺分类器→降级为问 | 2 层（档位 + 判定引擎） | 无 | 不一致：ACP 可直接改档，服务端无门 |
| **openchamber** | 不弹（一键循环） | ①档位不可用就从循环里删掉 ②安全网无答案→挂起 ③缺省→ask | 2 层；**改档会追溯清空挂起队列** | 无 | 一致（策略在服务端，多客户端同一份） |
| **Step-Code** | 不弹（启动 flag / 持久化） | **显式降级 `auto→confirm`** + 每调用 fail-closed + 指名 flag 的理由 | 2 层（preset 含"无 UI 答案" + 判定函数） | 无方向性，有"显式 vs 缺省"轴 | 一致（`ctx.hasUI`） |
| **pi-upstream** | N/A（无档位；`--approve` 即同意） | **返回更严一侧**（`return false`／`block:true`） | 1 层（per-call 钩子） | 无 | 一致 |
| **pi** | ⛔ 无此概念 | ⛔ 无此概念 | ⛔ | ⛔ | ⛔ |
| **oc.tar.gz** | 未解压（= `openchamber-main/`，结论同 openchamber） | — | — | — | — |

## ④ 一句总判

**主流做法是：把"确认"挂在"目标档位有多危险"和"这一发动作"上，不挂在"变宽"这个方向上；档位本身照改不误，缺确认通道时在动作层 fail-closed（拒这一发，并给出可执行的下一步），或在档位层把这一档降级／从可选集合里删掉——6 个有实现的 repo 里，没有一个是"因为没有确认通道就拒绝放宽档位"的。**

## ⑤ 对本仓那条冲突的直接建议（只报"主流会怎么做"）

**本仓现状**：三档（每步都问／只问高危／全自动）；机主 2026-09-21 定案"只有切进全自动要确认"；昨天一条腿改成"任何放宽都要先有确认通道"，⇒ 放宽到中间档、无通道时被**拒绝**。

**主流三选一的倾向 + 依据**：

1. **倾向＝"允许 + 动作层 fail-closed"，不是"拒绝这个放宽"**。依据：dsh 宿主 `/permission <preset>` 无任何确认判断（`permission-presets/src/index.ts:264-273`），门只在 Web 端且只对两枚目标档（`ui-permission-presets/src/client/index.ts:94`）；Step-Code 的档位本体自己带"无 UI 时怎么办"（`permissions.ts:37-70`、`:518-521`），改档动作不被拒，被拒的是**需要确认的那一发调用**（`:535-545`）。
2. **"中间档不需要确认"这条，主流是支持的**：dsh 的中间档 `workspace-write` 完全不在确认条件里；minimax 的中间档 `auto`（"Ask only when risk is high"）在 picker 里回车即生效；openchamber 的 `safety` 同理。⇒ **机主"只有切进全自动要确认"的定案，方向上与 3/6 个 repo 一致**（dsh／minimax／openchamber 都不给中间档配确认）。
3. **若一定要为"无通道"设一条规则，主流选的是"降级"，不是"拒绝改档"**：Step-Code `effectiveState` 把 `auto` 降成 `confirm`（`permissions.ts:518-521`，注释写明"否则 `--non-interactive-approval deny` 静默失效"），openchamber 把跑不动的档从循环里删掉（`permissionAutoAccept.ts:22-31`），pi-upstream 无 UI 时返回更严的一侧（`project-trust.ts:86-88`）。三条都是**结果更严**，但都不是"拒绝用户的这个动作"。
4. **⚠ 与本仓定案相反的一处，具名报上**：**dsh 给第二枚目标档 `auto`（沙箱关闭 + LLM 复核）也配了确认**（`client/index.ts:94`、`PermissionSelect.tsx:132`、`PermissionRow.tsx:81`）。即"只对最宽那档弹"在 dsh 不成立——它是"对**拿掉提示**的档弹"，共两枚。若本仓将来加"全自动 + 复核"这种第四档，按 dsh 的样子它也该弹。
5. **落地建议（三步，最小改动）**：①把 gate 从"改档动作"移到"需要确认的那一发调用"——即中间档照改，调用层缺通道时按更严档执行（取 Step-Code `:518-521` 的形状）；②无通道时给模型/用户一句指名出路的文案（取 Step-Code `:562-597` 或 dsh `:1763`）；③把"改档"这一笔落成耐久事件（取 dsh `permission-presets/src/index.ts:416` 的 `session.append('permission/preset', …)`），**这样即使没弹确认也留有账**。
