# 调研交件 — MiniMax code `packages/agent-modules` 整棵树（腿 3）

> 只读调研子代理 `survey-mmx-mod2`，2026-09-28。派单：
> `.scratch/wisp/dispatches/2026-09-28-195x-readonly-survey-wave3-three-named-uncovered-surfaces.md` §1 腿 3。
> 参照树：`D:\work\AI\open source\minimax-code\packages\agent-modules`（**无 `.git` ⇒ 全文不引提交历史**）。
> 纪律遵守：未修改参照仓、未在参照仓建文件；本仓只写本文件；未读 `frontend/**` 与 `design/**`；未 commit。

## 0. 分母读数（先复核编排者的名册，推导式与口径）

| 口径 | 推导式 | 读数 |
|---|---|---|
| 模块枚数 | `ls packages/agent-modules \| wc -l` | **12**（与派单名册一致，无名册外目录） |
| `.ts` 总数 | `find . -name '*.ts' -type f \| wc -l`（在 `packages/agent-modules/` 下） | **147** |
| 非 `.ts` 文件 | `find . -type f ! -name '*.ts'` | **12**（每模块一枚 `package.json`，无 README、无测试文件、无 `.tsx`） |
| 名册 | `background-task context-manager conversation-contract cron goal mcp permission plugin-hooks runaway-guard session-report skills system-reminder` | 逐名对上，**无手误** |

每模块 `.ts` 枚数（同一 `find` 口径，逐目录）：

| 模块 | `.ts` | 备注 |
|---|---|---|
| `permission` | **47** | 占整棵树 32%，护栏里最重的一枚 |
| `goal` | 19 | |
| `mcp` | 13 | |
| `system-reminder` | 12 | |
| `cron` | 10 | |
| `runaway-guard` | 10 | |
| `context-manager` | 9 | |
| `plugin-hooks` | 9 | |
| `background-task` | 8 | |
| `conversation-contract` | 2 | 只有 2 枚，但见 §7（它是"契约"层） |
| `session-report` | 4 | **见 §6：推翻编排者"五家都没有日报"那条** |
| `skills` | 4 | |

---

# 1. `runaway-guard`（10 枚 `.ts`）—— 防跑飞：它到底数什么

**先给一句结论（重要，因为它同时纠正派单里两句前提）**：
这枚模块**不停任何东西**。它的全部权力是"往对话里塞一条假的 user 消息"，
原文逐字写着 `At most one Steer attempt per Turn; never rejects a tool or aborts a Turn.`
（`packages/agent-extension/src/runaway-guard.ts:67`，那是把本模块接进 loop 的适配器，**在 agent-modules 树外**，
我为了回答"停谁"这一问必须读它，越界之处登记在 §①）。
**"停"这一维我们已经有**：`internal/agent/guard.go:14-19` 的 C22 LoopGuard 三级梯度 `[3,5,8]`
在顶级会 `Stuck` 硬停并出显式文案（`internal/agent/loop.go:470-483` 调 `brakeStuck`，`BrakeRepeat`）。
⇒ **派单 brief 里"我们有审批门但没有它自己转圈这一维"这句不成立**，见 §②第 1 条。

## 1.1 它数什么：六枚信号种（口径＝枚举原文）

`packages/agent-modules/runaway-guard/src/contracts.ts:3-9` 逐字六枚：

| 信号 | 数什么 | 判据落在哪一行 |
|---|---|---|
| `exact_action_repeat` | **单条**工具调用的 (toolName + 规范化 arguments) 指纹连续累计 | `step-view.ts:153-171`（算键）＋ `signals.ts:25-33`（累计） |
| `exact_result_repeat` | 工具**结果**的整体指纹（`{toolName,content,details,isError}`）重复 | `step-view.ts:184-187`，`signals.ts:43-51` |
| `same_error_family` | 错误归一化后的"家族"连续同类 | `step-view.ts:189-206` ＋ `errorFamily()` 388-398 |
| `abab_action_cycle` | 批次序列 `A,B,A,B` 交替 | `signals.ts:199-217`（要 4 枚连续批次，`state.recentActionBatches` 滑窗 4：205-207） |
| `polling_repeat` | 后台任务轮询：同一 task 的 (status, next_offset) 游标不动 | `signals.ts:34-42` ＋ `step-view.ts:106-152` |
| `unchanged_progress_repeat` | 工具**自己申报**的进度证据（`stateChanged`/`artifactChanged`/`newFacts`）跨多次不变 | `signals.ts:61-66`、`observeProgress()` 71-156 |

**数值的两个门槛（这就是"阈值多少"的原文答案）**：
- **观测**门槛＝**第 2 次**就发 observation（`signals.ts:174` `if (prior < 2 && occurrences >= 2)`；进度维另写一份在 `signals.ts:121-123`）。
- **提醒**门槛＝`remindAfterOccurrences`，**必须 ≥3 的整数**（`guard.ts:31` 的 `integerAtLeast(..., 3, ...)`），
  不传＝**只观测不提醒**（`guard.ts:28-31`、`contracts.ts:141-143` "Omit for observation-only mode"）。
  **生产默认值＝3**：`packages/agent-extension/src/runaway-guard.ts:69` `options.remindAfterOccurrences ?? 3`。
- ABAB 与"结果重复"**永远不触发提醒**：`RunawayGuardReminderSignalKind` 白名单只有 4 枚（`contracts.ts:177-180`），
  `exact_result_repeat` 与 `abab_action_cycle` 只进遥测。

**提醒的优先级（同时命中只发一枚）**：`reminder.ts:10-15`
`unchanged_progress_repeat > same_error_family > exact_action_repeat > polling_repeat`。

## 1.2 判定之后停谁、给用户留什么话

| 问 | 答（带行） |
|---|---|
| 停单条工具调用？ | **不停**。`guard.observe()` 返回一枚 `RunawayGuardReminder{content, observation}`（`contracts.ts:161-164`），适配器做的事是 `event.agent.steer({role:'user', content, timestamp})`（`agent-extension/src/runaway-guard.ts:130`）。`steer` 的语义是往运行中的 loop 排一条待读消息（`packages/agent-core/src/pi-turn-runner/events.ts:200-202, 228-230`）。 |
| 停整轮？ | **不停**，且**每 Turn 至多一枚**：`state.reminderAttempted` 在发出之前就占坑，注释逐字 `Reserve before the adapter calls steer: a failed attempt must not retry.`（`reminder.ts:26-28`）——**发歪了也不补发**。 |
| 停整个会话？ | **不停**。全模块唯一的外部动作是 `onSignal` / `onReminder` / `onTurnSummary` 三个回调（`contracts.ts:143`、`runaway-guard.ts:79, 132, 151`），回调失败**一律吞掉**（`signals.ts:244-254`、`runaway-guard.ts:252-262`）。 |
| 那"硬停"在哪一层？ | **不在这枚模块里**。本模块自带的另外三枚刹车我们也有（见 1.4 对照）。它的定位写在自己的 description：`Observe deterministic convergence signals and inject one bounded strategy reminder.`（`runaway-guard.ts:144`） |
| 给用户留什么话 | **给用户：没有直接文案**，用户只看到模型行为变化。**给模型：逐字四套话术**（`reminder.ts:72-101` + `63-70`），例如 `same_error_family`：`[runaway guard] The same tool error family has now occurred ${occurrences} times in a row. Do not retry the same route unchanged. Diagnose the cause, change one controlled variable or switch route. Do not infer that the entire task has failed from this signal.`。**每一套都带一段"防污染"后缀**：`This is a temporary runtime reminder for the current Turn only, not a user preference or a durable rule; do not save this reminder or generalize it into Memory, Skills, or other persistent instruction files...`（`reminder.ts:59-60`）——这条正是我们 D30 间接提示注入防护的同族顾虑。 |
| 轮询那一枚的话术 | 包 `<system-reminder>` 标签，且**许了用户可见的承诺**：`you will be notified automatically and this conversation will resume when the background task completes`（`reminder.ts:66-69`）——这句话只有 `background-task` 真能唤醒会话才敢写，见 §3。 |

## 1.3 三条我们大概率想抄的机关（各带行）

1. **审批拒绝不算跑飞**（这直接关系我们门控的误伤）：被 `blockedBy === 'permission'` 挡掉的调用**从检测里整条剔除**，
   `step-view.ts:61-65` 建 `excludedCallIds`，`:83-84` 连结果与进度都不给它看，`:116-118` 逐字注释
   `Permission rejection stays excluded from the old detectors, but its trusted task identity must still reset the new polling continuity.`
   ⇒ 也就是说"被拒 5 次"不会把它逼到刹车，但**被拒的那一次仍会把轮询计数的连续性清掉**。
2. **不信任模型自报，只信宿主证据**：`RunawayGuardTrustedToolProvenance` 由宿主在真正派发时打点，
   注释逐字 `A bounded host fact captured at actual tool dispatch, never inferred from an assistant message.`（`contracts.ts:25-34`），
   校验只认 `builtin`/`captured-compatibility` 两种来源（`step-view.ts:337-351`），
   且"宿主事实可以否决一次轮询观测，但不能凭空造一次"（`step-view.ts:282-292` 注释）。
   适配器只对 `task_output` 这一枚工具打点（`agent-extension/src/runaway-guard.ts:160-182`）。
3. **指纹用 per-turn 随机密钥做 HMAC，且限长**：`randomBytes(32)` 每枚 Turn 一份（`state.ts:101`），
   `createHmac('sha256', secret)`（`fingerprint.ts:21-23`），默认预算 `16 * 1024` 字节（`fingerprint.ts:8`），
   序列化稳定化（对象键排序 `fingerprint.ts:97`、深度上限 32 `:10/:47`、键数上限 1024 `:94`、
   非纯对象/带原型/循环引用一律拒收 `:57-62`、`:71`、`:89`），**超预算不是截断而是整条不检测**，
   只累加 `fingerprintSkippedCount`（`step-view.ts:149`、`:160`、`:187`）并进了 Turn 摘要的 `projectionSkippedCount`
   （`state.ts:68`）。注释逐字 `raw keys never leave the Turn or reach observers`（`contracts.ts:37-39`）。
   ⇒ 这一条对我们**特别对口**：我们的 `guard.go:234-243` 是把 `toolName|args` 明文拼进 `lastSig`，
   参数里如果有密钥/文档正文，就在进程内存里留明文副本（只在内存；落盘的是 `tool_call` 行，见 1.4 末注）。

## 1.4 逐维对照我们现量（"他们没有 / 我们有"是双向的）

| 维度 | MiniMax（带行） | 我们（现量） | 落到我们身上的后果 |
|---|---|---|---|
| 有没有"转圈"这一维 | 有，独立模块 | **也有**：`internal/agent/guard.go:14`（C22 LoopGuard），`internal/config/schema.go:422-425` `agent.loop_guard.repeat_thresholds` 默认 `3,5,8` | 派单前提要改口径：**不是"没有"，是"只有整轮签名一枚"** |
| 检测粒度 | **单条调用**级（每枚 toolCall 各算各的键，`step-view.ts:79-171`） | **整轮**级：`turnSignature()` 把这一轮所有 `ToolCalls` 拼成一枚字符串比较（`guard.go:234-243`），`loop.go:470` 每轮调一次 | 我们**看不见**"一轮里 A 重复但整轮签名在变"这种转圈：模型只要每轮加一枚无关小调用就能绕过 `[3,5,8]`。这是最值钱的一条差 |
| A-B-A-B 交替 | 有（`signals.ts:199-217`） | **没有**：`guard.go:171-176` 签名一变即 `repeat=1` | 两个文件来回改同一处、改一次看一次，我们**永远不会提醒也不会停** |
| 错误家族 | 8 类正则归一化＋结构化码优先（`step-view.ts:388-409`，`:400-409` 表逐名：`timeout/rate_limit/network/auth/permission/not_found/invalid_argument/process_exit`），错误文本截 1024（`:396`），只取首段 text 块且 ≤4096（`:421-428`） | **没有**〔现读 0 枚：`grep -rniE "error.?family\|same_error" --include=*.go internal/` 无命中〕；我们有 D37 错误分类（`guard.go:271-277` `ErrorClassOfTurnError`）但**不参与刹车** | 同一类错误连续撞：他们第 3 次就换路提醒；我们只在"参数一字不变"时才计数，**换个参数就重置** |
| 进度证据（"改了没有"） | 工具自己申报 `stateChanged/artifactChanged/newFacts`（`contracts.ts:40-48`），有正向变化即**当预期结果**、不建"无变化"假说（`step-view.ts:135-137`、`:164-166`） | **没有**〔现读：`grep -rniE "stateChanged\|artifactChanged\|newFacts\|ProgressProjection" --include=*.go internal/`＝**8 命中**，逐名看**全部**是 `internal/audio/mmdevice_windows.go:177/196/231` 的 Windows 音频设备回调（`onDeviceStateChanged` 等），**与工具进度申报无关**〕 | 我们的"重复"只是形状重复；他们判的是**没有产生新事实**，这一维差一整层 |
| 轮询豁免 | `task_query`/`task_output` 登记成 `kind:'polling'`（`local-runtime-v2/src/service/turn-system/runaway-guard/tool-policy.ts:23-31`），轮询**不触发**动作重复提醒，只在游标不动 ≥3 时提醒；`rg/grep` 退出码 1 被 `isExpectedResult` 判为预期、不算错误串（`tool-policy.ts:34-62`） | **没有豁免表**〔现读：`internal/agent/guard.go` 全文无按工具名的白/黑名单；`grep -n "exempt" internal/agent/guard.go` 0 枚〕 | 我们的正常轮询会被误计成跑飞；他们为"哪些工具重复是合法的"专门开了一枚可注入 policy 表（`contracts.ts:98-108`） |
| 预算类刹车 | **本模块没有** | **我们有**：token 预算（`guard.go:148-152`）、50 轮下限（`:141-144`）、单工具超时（`:126`）、并发上限（`:129`） | 别把这张表读成"他们全更强"——**预算这一侧是我们更全** |
| 关掉的能力 | 配置 `enabled` 默认 true，可远端覆盖，读不到值＝不加覆盖而不是打开（`packages/config/src/runaway-guard-config.ts:10-21`）；宿主热杀开关 `isEnabled`（`agent-extension/src/runaway-guard.ts:46-47, 88-92`）；关闭时**清空 streaks**（`:89`） | 我们的 ladder 是**配置常量**，无运行时热开关〔现读：`grep -rn "LoopGuard" internal/config/`＝**5 命中**，全部在 `schema.go:422/423/437/438` 与 `boundary_test.go:140`；`grep -rniE "loop_guard.*enabled\|LoopGuardEnabled" --include=*.go internal/`＝**0 命中**〕 | 我们出误伤时只能改配置重启，不能一键关 |
| 豁免特定轮型 | `shouldRemind: (ctx) => ctx.turnIntent?.kind !== 'goal-verifier'`（`local-runtime-v2/src/service/turn-system/runaway-guard/extension.ts:21`）——目标校验子轮的重复调用**是它的本职**，不许提醒 | **没有轮意图概念**〔现读：`grep -rniE "turnIntent\|TurnIntent\|轮意图" --include=*.go internal/`＝**0 命中**〕 | 我们票 197 子代理三层落地后，**校验型子代理会天然撞我们的 `[3,5,8]`**，这条现在是隐形地雷 |
| 离线回放 | `replayRunawayGuardTrajectory` 跑**同一枚生产投影器**回放轨迹（`replay.ts:19-66`，注释 `never steers and has no persistence side effect`），并量化"信号之后又多烧了多少"（`replay.ts:68-101`：`postSignalStepCount`/`postSignalToolCallCount`/`postSignalProviderTokens`） | **没有**：`internal/agent/loop_golden_test.go` 是 golden SSE 而不是轨迹回放；不过原始数据我们有——`internal/agent/journal.go:20-26` 的 `tool_call` 行已带 `decision/outcome/error_class` | 他们能回答"这枚护栏上线后省了多少 token"，我们**答不出来**（对 D32 预算/SLO 是一条免费的证据链） |
| 全 fail-open | 检测/宿主事实/steer/观测器四处各自 `catch {}`（`agent-extension/src/runaway-guard.ts:97-99, 133-135`；`signals.ts:250-253`） | 我们刹车是**决定性**的（撞线就 Stuck），不存在 fail-open 分支〔现读 `loop.go:470-483`〕 | 注意**这是取向差不是缺口**：他们选"护栏永不伤害执行"，我们选"撞线必停"。改哪边都要过 D43 转移表，属契约级 |
| 遥测载荷 | 白名单投影：只出 `schema_version/policy_version('runaway-v1')/session_id/turn_id/agent_name/signal_kind/step_index/occurrences`，日志与上报各自 `try{}catch{}`（`local-runtime-v2/src/service/turn-system/runaway-guard/observation.ts:19-56`） | 我们出的是 `task_journal` 落库＋事件汇（`internal/agent/loop.go:472-476` `EvReminder`），**没有 policy_version 这一维**〔现读：`grep -rn "policy_version\|PolicyVersion" internal/` 0 枚〕 | 我们日后改 ladder 值，历史数据无法区分是哪一版策略产生的——回滚时无法归因 |

**跑在哪一层**：纯**进程内**、零 IO、零持久化，模块自己的注释逐字写着
`Turn-local detector and reminder policy. No Agent, lifecycle registration, host IO, Memory writes, or persistence belongs to this module.`
（`guard.ts:19-22`）。唯一系统依赖是 `node:crypto` 的 `randomBytes`/`createHmac`（`state.ts:1`、`fingerprint.ts:1`）
⇒ **落到 Win32＋Go 无障碍**：`crypto/hmac`＋`crypto/rand` 直换。

---

# 2. `permission`（47 枚 `.ts`，占全树 32%）—— 授权梯度长什么样

## 2.1 他们的梯度是**四枚正交轴**，不是一条梯子

| 轴 | 取值 | 行 |
|---|---|---|
| 判定结果 `PermissionBehavior` | `allow \| deny \| ask` | `types.ts:11` |
| 档位 `PermissionMode` | **六枚**：`default \| acceptEdits \| bypassPermissions \| auto \| dontAsk \| off` | `types.ts:13-19` |
| 规则来源 `PermissionRuleSource` | **三层**：`global \| agent \| session` | `types.ts:38`，落盘形状见 `types.ts:242-261`（`permission.json` 的 `allow/deny/ask/defaultMode`） |
| 内部策略 `AskForApproval` | 四枚：`on-request \| on-request-llm \| never \| deny`（档位到策略的**唯一映射表**） | `ask-policy.ts:26-54` |

**关键机关：档位不是布尔，是一枚四轴 profile。** `PermissionModeProfile`
（`permission-core.ts:112-121`）把"这一档允许做什么"拆成
`interaction: 'ask' | 'neverAsk' | 'headless-fail-closed'` ／
`classifier: 'allowed' | 'forbidden'` ／
`enforcement: 'strict' | 'off'` ／
`bypass: 'allowed' | 'forbidden'`，映射表逐档在 `permission-core.ts:523-549`。
最要紧的第三枚 `enforcement: 'off'` 只给 `mode:'off'`（`:532-539`）——**"权限层不参与阻断"与"永不询问"被刻意分开写**，
`bypassPermissions` 是 `neverAsk + enforcement:strict`（`:524-531`），也就是**仍然走完整判定，只是把 ask 折成 allow**。
另外 `interaction:'headless-fail-closed'` 这一支在没有审批通道时**直接判 deny**，
理由逐字 `'permission requires user confirmation but no approval channel is available'`（`permission-core.ts:422-430`）。
⇒ 对照我们：`internal/risk/mode.go:44-56` 的三枚 `Mode` 是**同一把梯子的三格**（`ModeAskEveryStep` 是零值，注释逐字
"an uninitialized or unrecognized setting lands on the strictest档, not on the loosest one"）——
**我们的"零值最严"比他们更安全**（他们 `mode ?? 'default'` 也算严，但 `off`/`bypassPermissions` 是两枚可持久化的档），
**但我们没有"无审批通道 ⇒ 拒绝"这一支**〔现读：`grep -rniE "headless|no approval channel" internal/agent/approval internal/risk --include=*.go` 无命中〕。

## 2.2 判定顺序（他们真正的资产是一张 14 级 precedence 表）

两条路径并存：`engine` 回滚路径（`engine.ts:252-513`，三步流：Reject→Allow→Fallback，注释在 `:182-201` 逐级编号 1a…3）
与 `core` 生产路径（`permission-core.ts:159-300`）。**core 路径的次序**（逐行）：

| 序 | 判据 | 结果 | 行 |
|---|---|---|---|
| 1 | 整工具 deny 规则 / MCP server deny | deny | `:162-164` |
| 2 | **工具检查器自己 deny** | deny，**且赢过更宽的 ask**（注释逐字 "otherwise adding a blanket ask can weaken a command-specific hard block"） | `:166-171` |
| 3 | 整工具 ask 规则 | ask（唯一例外：bash＋沙箱＋`autoAllowBashIfSandboxed`） | `:173-185` |
| 4 | MCP server ask | ask | `:187-194` |
| 5 | **检查器抛异常** | **ask（fail-closed）** | `:196-202` |
| 6 | 注册为"需要人交互"的工具（默认集合只有 `'question'`） | ask | `:204-215`、`context.ts:161` |
| 7 | 检查器 ask 且 `bypassImmune` | ask，带两枚旗标 | `:217-227` |
| 8 | **内容级 ask 规则**（对工具输入做 `includes()` 子串匹配） | ask | `:229-239`、`:339-350` |
| 9 | 检查器 ask（非 bypass 档） | ask | `:241-251` |
| 10 | bypass 档 | allow | `:253-259` |
| 11 | 整工具 allow 规则 / MCP allow | allow | `:261-269` |
| 12 | 检查器 allow | allow | `:271-277` |
| 13 | **没有注册检查器** | **allow（默认放行！）** | `:279-287` |
| 14 | 有检查器但没规则命中 | ask | `:289-299` |

**第 13 级是必须点名的取向差**：他们未知工具的兜底是**放行**，我们 R9 的兜底是**升 L2**
（`docs/specs/SPEC-06` 行 42 逐字 "R9 无判定器可用或任一判定器 panic → fail-closed 升 L2"；
`internal/risk/mode.go:26-28` 也把这写成红线）。
他们靠**第 12 级之前的 `decideUnknownToolPathCapabilities`**（`engine.ts:232`、`engine.ts:290`）兜住——
**未注册工具先按"参数里能不能提出路径意图"过一遍能力闸门**，提不出意图才落到 allow；
且 `apply_patch` 特例：提不出路径意图时**反而 ask**，注释逐字 "fail closed: ASK instead of letting the engine's
no-checker default-allow wave an opaque write through"（`tools/path-capability.ts:1049-1063`）。
多意图聚合是 **deny > bypass-immune ask > ordinary ask > allow**，注释解释为什么要**扫完全部**而不是取第一个非 allow
（`:1067-1099`，逐字提到"早期一枚普通 ASK 会遮蔽后面的 bypass-immune ASK，让 bypass 档把敏感写放过去"）。
⇒ 后果：**我们若日后把插件工具（D46 Tier-1 CLI 包装）接进判定器，"没有判定器的工具"这一格必须自己拍死是哪一档**，
他们的答案是"路径能力兜一层＋兜不住才放行＋写类工具例外"，不是简单的 allow 或 deny。

## 2.3 授权梯度（D45 那一维，他们做成了"候选作用域"）

`PermissionDecision.candidateScopes` ＋ `CandidateScope`／`CandidateScopeGroup`（`types.ts:151-195`）：
用户点"始终允许"时**不是写死一条规则**，而是**一次列出若干宽度供选择**：

- `kind: 'narrow' | 'byFirstWord' | 'byArgvPrefix2' | 'byDomain' | 'wholeTool'`（`types.ts:171`）
- **index 0 必须是 `narrow`**，向后兼容契约逐字写在 `types.ts:96-106`
  （`candidateScopes[0] ... MUST always be equivalent to the current ruleContents (the narrow default)`），
  更宽的档从 index ≥1 出现，客户端要显式回 `selectedScopeIndex` 才拿到宽档；
- 每个候选带 `labelKey`（i18n 键）＋ `labelParams`（`types.ts:173-176`），**宽度是 UI 可排序/可过滤的机器可读枚举**；
- 一组＝"一个待批主体"，复合命令里**每枚要问的子命令各一组**（`types.ts:179-195`）。

对照我们现量：
| 他们的档 | 我们 |
|---|---|
| 持久规则库三层（global/agent/session）＋落盘 `permission.json` | **只有 session 一档有设计、零一档有实现**：`docs/PLAN.md:2140-2152`（D45 三元组 (工具,路径模式,会话)，落 S7）；`docs/PLAN.md:1535` 逐字保留 **RESERVED**："无「永久允许此工具于某路径」；每次新会话都要重新授权一次" |
| `allow_session_grant` 决策值 | **词汇存在、生产者零枚**〔现读：`grep -rn "DecisionAllowGrant\|allow_session_grant" --include=*.go .` 非测试命中只有 4 处，全在 `internal/agent/journal.go:32`、`internal/memory/models.go:81/138`、`internal/memory/schema.go:76` 的**枚举/注释**里，`internal/agent/approval/` 无一处写入〕 |
| 宽度候选枚举（5 档） | **无枚举，只有 PLAN.md:2028 手写的三档**（"仅本次"／"本会话允许向此应用注入"／"本会话允许向所有应用注入"），且只针对 `input.type` 一枚工具；面板"安全"页可查可撤销也在那一条里写着 |
| "L2 永不进持久授权" | **一致**：`docs/specs/SPEC-06:121`、`docs/PLAN.md:1592`；他们侧的对应物是 `bypassImmune`（`types.ts:110-115`）与 `hardBlockedCategoryIsFinalDeny`（`classifier/dangerous-patterns.ts:190-203`） |

⇒ **这一枚是本轮对我们 D45 最直接可抄的东西**：它把"授权梯度"从**卡片上的三个按钮**升级成
**机器可读的宽度枚举＋index0 永远是窄档的兼容律**，于是"宽档"以后能被审计、被回滚、被按 kind 统计。
我们现在若加一档，只有改文案一条路。

## 2.4 风险的确定性两层表（可直接抄结构）

- **HARD_BLOCKED 15 类**（`classifier/dangerous-patterns.ts:130-188` 枚举逐名；枚数推导式＝
  `sed -n '130,158p' dangerous-patterns.ts | grep -c "^  | '"`＝**15**；同一份也逐名列进
  `HARD_BLOCKED_CATEGORIES` 集合 `:171-188`）：
  `system-secret · ssh-credential · sensitive-read · data-exfil · catastrophic-standalone · irrecoverable-delete ·
  disk-erase · storage-volume-delete · windows-secure-erase · windows-delete · ransomware-indicator ·
  fs-inode-direct · archive-remove-source · encoding-bypass · reverse-shell`
  ⇒ **只有其中 10 类是"最终拒绝"**：`hardBlockedCategoryIsFinalDeny()` 白名单在 `:190-203`。
  **剩下 5 类（`system-secret`/`sensitive-read`/`ssh-credential`/`encoding-bypass`/`reverse-shell`）是 bypass-immune ASK 而不是死拒**
  ——这跟我们 R3 的"敏感路径 A 档直接拒绝"比，**多了一个"敏感读可以问人"的中间态**。
- **SOFT_RISK 12 类**（`:792-809`；推导式＝`sed -n '792,810p' | grep -c "^  | '"`＝**12**）：
  `permission-change · network-listen · privilege-escalation · disk-tools ·
  interpreter-inline · powershell-encoded · sensitive-path · container-exec · scheduled-execution · content-clear ·
  remote-execution · windows-management`（＋注释说明每类为何是软）。软＝**踢给云端 LLM 判定**（`:789-790` 逐字
  `SOFT_RISK registry — kicks the fast-path into the cloud LLM gate`）。
- **敏感路径正则表**（`:217-251`）**已经含 Windows 目标**：`\.ssh/id_{rsa,dsa,ecdsa,ed25519}`（显式排除 `.pub`）、
  `.aws/credentials`、`/etc/{shadow,gshadow,sudoers}`、私钥词干＋`.key`、`.p12`、`.pfx`、`private*.pem`、
  macOS System.keychain、`/proc/{self,\d+}/environ`、**`Windows/System32/config/SAM` 与 `Windows/NTDS/`**
  （`:250-251`，注释逐字"accepting both separator styles in command text"）。
  另有 `HARD_BLOCKED_FS_READ_PATTERNS`（`:676`）与 `SECRET_READ_PATTERNS`（`:1467`）两张"读"侧表，
  以及 `matchSecretReadForm`（`:1529`）管 `printenv` 这类**不带路径的泄密形态**。
- ⚠ **他们自己写下的一处 Windows 缺口**（对我们极有参考价值）：`windows-delete` 这一类注释逐字
  `CLI-driven deletes on Windows skip the Recycle Bin entirely and are irrecoverable ... until that rewrite path is
  wired for Windows delete commands (del/rd/Remove-Item), the safest policy is HARD deny`（`:140-145`）。
  ⇒ 他们的 POSIX `rm` 有"改写成可回收删除"的通路（见 2.5），**Windows 三条删除命令还没接上，先用硬拒顶着**。

## 2.5 一件我们完全没有的能力：**权限层可以改写要执行的命令**

- 决策形状带 `rewrittenInput?: Record<string, unknown>`（`types.ts:108-109`），
  core 在 **6 个不同分支**上都把它往下传（`permission-core.ts:183, 237, 248, 257, 267, 275, 298`）。
- 改写类型目前是**一枚**：`ExecutionTransform = { type: 'recoverable-delete'; targets: string[] }`（`permission-core.ts:103`）。
- 理由文案逐字：`Allowed: rm command was rewritten to mavis-trash for recoverable deletion. Original command: ...`
  （`reason-format.ts:114`）；`DecisionReason` 里对应两枚变体 `rmRewrite` 与 `recoverableDeleteRewrite`（`types.ts:72-73`）。
- **跨层传递用的是 Symbol 键**：`windows-trash-execution.ts:6` 取
  `Symbol.for('@mavis/permission/windows-trash-execution')`，写入时 `Object.defineProperty(..., {configurable:false, writable:false})`（`:24-29`），
  读取端做全量类型复核（`:33-45`），注释逐字
  `Only literal trash targets cross this seam; the runtime owns the launcher, script path and environment.`（`:3-5`）
  ⇒ **一层塞给另一层的"改道凭证"，不是靠参数字段名约好的，是靠不可伪造的 Symbol＋冻结对象**。
- **跑在哪一层**：权限层进程内决策，真正执行由 local runtime 侧的 bash 实现拿 Symbol 里的目标去起 `mavis-trash.js`
  （子进程／脚本，见 `windows-native-delete.ts` 162 行与 `fs-permission.ts`；**我未逐行读完这两个文件的执行侧**，登记在 §①）。

## 2.6 "给用户留什么话"这一维他们做成了一等公民

- 每枚决策都要出**双语**理由：`ACTION_LABEL` `en{Allowed/Needs confirmation/Blocked}` /
  `zh{已允许/需要确认/已拒绝}`（`reason-format.ts:27-30`），表头注释逐字
  `When adding a new DecisionReason variant, add its template here too — missing keys surface as a TypeScript error.`（`:7-9`）
  ——**用类型系统钉住"新判据必须有文案"**，这条我们完全可以照抄成 Go 的穷尽 switch。
- 理由**分两个受众、两种语言策略**：给人看的那份本地化；**给模型看的那份前缀固定英文**，
  `formatBlockedToolReason` 注释逐字 `The leading [permission:denied source=<source>] tag is intentionally fixed English
  so the model (and log greps) can parse the origin regardless of the user's locale`（`:324-328`），
  实际串 `:330-338`。
- 拒绝来源分四枚 `ToolDenialSource = 'user' | 'rule' | 'safety' | 'safety-immune'`（`:306`），注释解释每一枚何时用
  （`:293-305`，其中 `safety-immune` 举的例子正是 **UNC／网络盘**＋`rm -rf /`＋`rm -rf ~`）。
  ⇒ 对用户"我拒的"和"规则拒的"和"墙挡的"是三种心智；模型也需要这个区分才知道该换路还是该放弃。
  **我们现在给模型的拒绝文本是哪一种形状**〔现读：`internal/agent/approval/` 与 `internal/risk/` 里
  `grep -rniE "denied source|permission:denied" --include=*.go` 无命中〕。

## 2.7 云端判定（层＝**服务端**，且有明确隐私代价）

- 端点：`POST {region-routed-host}/mavis/api/v1/permission/check`（`classifier/cloud-classify-client.ts:5-6, 51`），
  区域表 `cn → https://agent.minimax.cn`、`en → https://agent.minimax.io`（`:26-39`），
  与内容安全**共用同一网关**（`:23-24`）。
- 载荷字段逐名（`cloud-gateway.ts:10-21` 与 `http-cloud-gateway-client.ts:63-74`）：
  `tool_name · input（**序列化后的整条工具入参**） · platform · home_dir · workspace_root · mode:'auto' ·
  conversation_context · agent_id · session_id · daemon_prompt_version`。
  `conversation_context` 由本地渲染（`conversation-renderer.ts:73-99`）：
  最近 **3** 条用户消息（>500 字符则首尾各留 250，`:110-113`）＋最近 **5** 条整会话（每条截 200，`:34-36`、`:105-108`），
  工具调用压成 `tool:<name> args=<截断> result=<截断>`（`:121-128`）。注释逐字
  `the rendering algorithm must match what the server-side prompt expects byte-for-byte`（`:5-6`）。
- **裁决词表**：`allow | confirm | block` ＋客户端合成的第四枚 `timeout`（`cloud-gateway.ts:83-84`）。
- **fail-closed 到底怎么闭**：非 2xx／非 JSON／业务码非 0／verdict 不合法／网络错／abort **全部**折成 `kind:'timeout'`，
  注释逐字 `The daemon then surfaces an ask-user card so the user is never silently bypassed when the gateway is unavailable`
  （`http-cloud-gateway-client.ts:11-16`、`:165-172`）；单请求超时默认 **60_000ms**（`:31`、`:85` 用 `AbortSignal.timeout`）。
  还会区分"真超时"与"网络错"以便卡片说实话（`:139-149`）。
- **只在托管运行时启用**：`shouldUseCloudClassify()` 逐字 `return isManagedRuntime()`（`:61-63`），
  注释说明本地/离线开发守护进程保留本地路径（`:56-60`）。注入缝是 `CloudGatewayClient` 接口，
  测试给 `InMemoryCloudGatewayClient`，**离线/BYO-key 部署可以塞一枚恒返回 timeout 的空实现把流程逼回问人**
  （`cloud-gateway.ts:110-118`）。
- 我们侧现量：`docs/PLAN.md:1534` 把"Codex guardian 双轴评分（风险等级 × 用户授权度）"整条挂在 **RESERVED**，
  理由逐字 `L1/L2 每次 +1 LLM 调用且在响应路径上（+300–800ms），与 D18 延迟预算相左`，并写了重评触发条件
  `若实测发现 L1 误执行率不可接受，优先启用此项而非放宽门控`。
  ⇒ **对照结论**：他们把那一枚 LLM 判定**搬到了服务端**（因此不占我们最在乎的本地响应路径预算，
  但把"你的整条工具入参＋最近对话"送出了机器）。**这两条不能互相冒充对方**：
  我们若要抄"LLM 当二次判定"，`docs/PLAN.md:1534` 那条 RESERVED 的延迟论证在**本地跑**才成立；
  走云端就得另开一条隐私账（属 D4／D33 一级的契约变更，须人工批准）。

## 2.8 其余结构（本轮只点名，细节见 §①）

- 上下文是不可变对象：`Object.freeze` 深冻结＋`applyPermissionUpdate` 返回新对象（`context.ts:141-175`、`:190-204`、`:210-233`），
  三种更新 `addRules/replaceRules/removeRules`（`types.ts:215-236`），去重靠 `ruleValueEquals`（`context.ts:236-249`，
  **path 类规则要连 actions 集合一起比**，命令类只比 pattern）。
- 上下文里带**平台与 shell 家族**：`platform`＋`shellFamily: 'posix'|'cmd'|'powershell'|'unknown'`（`context.ts:63-67`，默认 `'unknown'` `:150`）
  ⇒ 我们侧对应物〔现读：`ls internal/risk | grep -c '_windows'`＝**8**（其中非测试仅 **2** 枚：
  `pathresolver_windows.go`、`syncdirs_windows.go`）；`grep -c '_other'`＝**4**（非测试 **2** 枚：
  `pathresolver_other.go`、`syncdirs_other.go`）〕，
  但**shell 家族没有作为判定输入的一等字段**〔现读：`grep -rniE "shellFamily|ShellFamily" --include=*.go internal/`＝**0 命中**〕。
- 检查器可注册/可注销（`engine.ts:165-178`），且**观测回调不许改动裁决**
  （`notifyCheckerDecision` 整体 `try{}catch{}`，注释 `Observability must never alter the permission verdict.`，`engine.ts:519-538`）。
  观测形状 `PermissionCheckerDecisionTrace`（`permission-core.ts:25-35`）带 `rewriteApplied: boolean`——**"改写有没有生效"是被单独钉住的一枚读数**。
- `CommandIntent` 是一套**跨 shell 的命令中间表示**：5 枚 kind（`filesystem`/`execute`/`network`/`script`/`opaque`）
  ＋ `PathValue` 带 `raw/resolved/dynamic/resolution` 四字段（`resolution: 'absolute'|'relative'|'home'|'environment'|'glob'|'unknown'`）
  （`permission-core.ts:46-90`）。**`opaque` 是一等公民**（承认"这条我解析不了"），
  且 `dynamic: true` 明确表示"环境变量/glob 还可能改变语义"。
  ⇒ 我们侧对应物〔现读：`internal/risk/rules_shell.go` 存在；`grep -n "opaque\|Opaque" internal/risk/*.go | grep -v _test` 无命中〕。

---

# 3. `cron`（10 枚）＋ `background-task`（8 枚）—— 定时与后台：怎么持久、重启怎么续、能不能列和取消

## 3.1 cron：配置落盘是 `.md` 带 frontmatter，调度器由宿主注入

| 问 | 答（带行） |
|---|---|
| 持久形态 | **每枚任务一个 `.md` 文件**：frontmatter 走 `CronFrontmatterSchema`（`cron/src/types.ts:83-97`），**prompt 在正文里不在 frontmatter**（`:76-77` 注释逐字 `The prompt field lives in the markdown body`）。字段逐名：`name · disabled · schedule · scheduleType('cron'\|'once') · runAtMs · deleteAfterRun · timezone · activeHours{start,end HH:MM} · session · delivery{channel,chatId}` |
| 会话落点 | `session.mode` 是三选一判别联合（`types.ts:26-39`）：`root`（回到主会话）／`sessionId`（指定会话）／`new`（每次开新会话，带 `keepSessions`：`undefined`→默认 **3**、`null`→**全留不归档**、`N≥1`→留最近 N 枚） |
| 断电/重启后怎么续 | ① 启动时 `CronRegistry.start()` **扫盘、逐枚重登记**（`registry.ts:101-141`）：`cronStore.listAll()` → 建 job → `nextRun = job.nextRun()`；② **重启前错过的档期不补跑**〔现读：`grep -rn "overdue\|compensat\|misfire\|missed" packages/agent-modules/cron/src packages/agent-tools/src/shared/cron-schedule-input.ts`＝**0 命中**；`lastRun` 在重启登记时一律写 `null`（`registry.ts:115`）〕——"每天 9:00"那条若 8:55 崩溃、9:30 拉起，**9:00 那一次不会补**；③ 启动还跑两把 fire-and-forget 清扫：会话保留清扫（`registry.ts:143-158`，注释逐字 `retention policy that was previously dependent on an in-memory map and therefore lost across daemon restarts`）＋孤儿会话清理（`:160-182`） |
| 恢复时可以拒绝执行 | `CronRegistryStartOptions.restoredTaskExecution: 'enabled' \| 'quarantined'`（`registry.ts:42-48`），注释逐字 `A quarantined clone still registers tasks for inspection and later user mutations`；实现落在 `:114`（quarantined ⇒ `enabled:false`）。⇒ **从别的机器/别的副本恢复来的定时任务可以先只登记不执行**，这一枚形状我们完全没有 |
| 定时器归谁 | **宿主注入**：`CronSchedulerPort.createJob/assertSchedulable`（`registry.ts:55-62`），没注入就直接抛 `CRON_SCHEDULER_MISSING`，错误文案逐字 `Runtime hosts must inject a scheduler instead of relying on agent-core timers.`（`:84-95`）；另有 `unrefTimers` 让短命嵌入宿主能在服务完一次性请求后退出（`:33-37`）。⇒ **模块自己不持有 setInterval**，这点对"能不能落到 Go 里"很干净 |
| 用户能不能列出来 | **能，两条路**：① 面板/HTTP：`listTasks(agentName)`（`registry.ts:196-205`）＋ `CronTaskResponse`（`types.ts:190-220`，含 `enabled/nextRun/lastRun/lastResult/lastError/status/activeHours/delivery`）；② **模型自己也能列**：`cron list / get / create / self / once / update / delete / resolve-model` 八条子命令写在工具说明里（`packages/agent-tools/src/desktop/builtin-defs.ts:989-996`，**越界读物**，登记 §①） |
| 用户能不能取消 | **能**：`cron delete`（同上 `:996`）、`disabled` 开关（updateConfig 路径 `registry.ts:245-258`，`update.disabled` 见 `types.ts:122`），以及**排队中的那次也能撤**：`cancelQueuedCron(agentName, cronName)` 返回取消枚数（`executor.ts:114-116` → `BusyQueue.removeCron`） |
| 忙的时候怎么办 | `BusyQueue`（`busy-queue.ts`）：**TTL 30 分钟**（`:4` `TTL_MS = 30 * 60 * 1000`），排空结果五枚枚举 `'sent' \| 'expired' \| 'removed' \| 'replaced' \| 'failed'`（`:6`）——**latest-wins 替换**与**过期**都是显式结局。排空触发点不是"只有会话结束"，注释逐字 `When the in-flight turn ends via error or abort ... otherwise it sits in the queue forever`（`executor.ts:77-81`），且 `session.error`/`session.abort` **按事件来源白名单过滤**（`:95-111`，逐字说清"session-title-service 发的 session.error 不许把 cron 队列排空"） |
| 跳过时的账 | `ExecuteResult.reason` 至少四枚：`'disabled' \| 'outside_active_hours' \| 'invalid_timezone' \| 'already_running'`（`executor.ts:153-183`），**时区配置非法是单独一枚 reason**（`:169-175`），不是悄悄跳过 |
| 交付到 IM | `delivery{channel, chatId}`（`types.ts:64-71`）＋ `report-delivery.ts` 211 行；两个旧字段 `report_to_root/report_to_main` 标 `@deprecated` 且注释逐字 `never writes to Root`（`types.ts:93-96`、`:110-113`） |
| 已退役任务的收尾 | `cleanupRetiredCronTasks`（`retired-cleanup.ts:42`，接口 `RetiredCronFileStore` `:19`） |

**⚠ 他们自己在工具说明里写下的一处倒退**：`cron create/update` 的 `activeHours` 参数描述逐字
`Legacy cron create/update: daily run window. Cron V2 does not support active-hour windows.`
（`builtin-defs.ts:824`）。⇒ **"只在某时段跑"这一维在他们的 V2 命令面上已经不被支持，但 schema 与 executor 检查都还在**
（`types.ts:91`、`executor.ts:160-176`）——这一处"规格还在、入口已关"的形状，正是我们票 164 AC#1 想终结的那类，可当反面样本。

## 3.2 background-task：状态机、活性模型与"投递账本"

枚数口径：`find packages/agent-modules/background-task -name '*.ts' | wc -l`＝8，合计 916 行（`wc -l` 同口径）。

| 机关 | 细节（带行） |
|---|---|
| 任务种类 | `bash \| subagent \| workflow \| custom`（`types.ts:1`）——**子代理本身就是一种后台任务** |
| 状态机 | 七枚：`queued \| running \| stopping \| succeeded \| failed \| canceled \| lost`（`types.ts:3-10`）。**`stopping` 与 `lost` 都是正式状态**，不是错误码 |
| `lost` 怎么来 | `TASK_LOST_ON_STARTUP`：重启后发现宿主进程不在了 ⇒ 落 `lost` ＋ `lastError.code` ＋发 `completed` 事件（`packages/local-runtime/src/background-task/startup-recovery.ts:88-106`，理由文案 `'Local runtime restarted before task completed'` `:89`）。**越界读物（local-runtime 包），登记 §①；但它是"断电后怎么办"的唯一真答案，值得抄** |
| 不许把年龄当证据 | 注释逐字 `Legacy rows have no owner identity. Only exact recovered Turn evidence permits settlement; age or an idle Session does not prove ownership loss.`（`startup-recovery.ts:147-149`）⇒ **只有"这一枚 Turn 确实没活下来"的精确证据才允许判 lost** |
| 投递账本 | 任务上有 `deliveredAt?`（`types.ts:82`），查询位 `undeliveredOnly`（`:108-110`，注释 `only return tasks whose completion has not yet been delivered`）。落库为列 `delivered_at_ms` 并**给它建了复合索引** `(owner_session_id, status, delivered_at_ms, ended_at_ms, task_id)`（`local-runtime/src/background-task/schema.ts:38-45`，migration version 20）。⇒ 这条就是 §1.2 里那句 `you will be notified automatically` 兑现的机制：**"跑完了但还没告诉会话"是一个可查、可索引的状态** |
| 取消语义 | `stop()` 先落 `stopping` 并发 `stop_requested`，再调 runner；**runner 停止失败也照样落 `canceled`**，但在事件 payload 里写 `stopFailed: true`（`manager.ts:120-147`）。已是终态直接返回不动（`:122`），已在 `stopping` 也直接返回（`:123`）——幂等 |
| 输出读取 | `TaskOutputStore {append, read(offset,limitBytes,stream), tail, finalize}`（`types.ts:62-67`），读结果带 `nextOffset/truncated/summary/outputRef`（`:54-60`）；流四枚 `stdout \| stderr \| transcript \| final_result`（`:37`）；outputRef 落点五枚 `file \| oss \| sandbox_job \| session_transcript \| memory`（`:26`） |
| 生命周期事件 | 六枚 `created \| started \| output_updated \| stop_requested \| status_changed \| completed`（`:117-123`），每条带 `sequence`（`manager.ts:207`），控制器可选 `watch?(query): AsyncIterable<TaskLifecycleEvent>`（`types.ts:203`） |
| 存储在哪一层 | 模块本身只有 in-memory 实现（`in-memory-store.ts` 172 行、`in-memory-output-store.ts` 129 行）；**真持久化在宿主**：SQLite 表 `local_runtime_background_tasks` ＋ `local_runtime_background_task_events`（`schema.ts:7-36`，migration 19），**整条记录另存一份 `record_json TEXT NOT NULL`**（`:15`，列取形状＋JSON 取全量） |
| 计量 | `TaskUsage {inputTokens, outputTokens, totalTokens, costUsd, metadata}`（`types.ts:18-24`）——后台任务烧了多少**是任务对象上的字段**，不是另算 |

## 3.3 对照我们现量

| 维度 | 他们 | 我们（现量） | 后果 |
|---|---|---|---|
| 有没有应用内调度器 | **有**（常驻 daemon ＋ cron 模块） | **刻意没有**：`docs/PLAN.md` D12 标题逐字「应用内零调度器，定时交给 OS」，处置为"CLI 入口＋Windows 任务计划/launchd/cron 调 `assistant run`"（`PLAN.md:306-320`），安全论据逐字 `一旦 Agent 能定时自触发，就出现 **无人值守执行 × D4 需要有人在场才能生效** 的冲突` | **这不是缺口，是取向**。但他们给了我们一条 D12 里没讨论过的逃生口：`restoredTaskExecution:'quarantined'`＋`activeHours`＋`deleteAfterRun`——"OS 定时唤起的这次运行"能不能带一份"只列不跑"的档期表，D12 那句话没有这一维 |
| 后台任务持久化 | SQLite 两枚表＋事件表，重启后逐枚对账（见 3.2） | **进程内一张表**：`internal/tools/task.go:165` 注释逐字 `TaskRoster is the process-local background task table`；结构 `byTask map[string]TaskOutput` ＋ `cancel map[string]func()`（`:171-189`）〔现读：`grep -rn "TaskRoster" --include=*.go internal cmd \| grep -v _test`＝**10 命中**，全在 `internal/tools/task.go` 与 `subagent_197.go`；`internal/memory/` 无对应表〕 | 崩溃/重启后后台任务全成孤儿，且没有 `lost` 这一枚状态去解释"它到底死没死" |
| 后台任务 kinds | 四枚（bash/subagent/workflow/custom） | **只有 subagent 一种有实现**：`internal/tools/subagent_197.go:20`（`a record with its own task id in the TaskRoster, its own agent.Loop`）；bash 类后台执行〔现读：`grep -rniE "background|后台" --include=*.go internal/tools/ | grep -v _test`＝**4 命中**，逐名是 `recycle_windows.go:32/45/140` 的"后台工具调用不许挂住"注释与 `task.go` 的表定义注释，**没有一枚是 bash 后台 runner**〕 | 我们说的"后台"其实是"子代理"，不是"任意长时间运行的操作" |
| 用户能否列出后台任务 | 能（`task_query` 工具＋`list(query)`＋`undeliveredOnly`） | **DEFERRED 且已自记**：`docs/PLAN.md:1531` 逐字 `**task.list / task.cancel 的实现**（D34 名册在册、生产注册表零枚）`，判据 `注册表里 task.list／task.cancel 各有一枚真实现并过契约测试`，后果列逐字 `Agent 侧查不到"谁被阻塞"，D31 那条"必须可见"只剩 UI 半条腿`；代码自证 `internal/tools/task.go:19-27`（`task.list - DEFERRED`／`task.cancel - DEFERRED`，只有 `task.output` 有实现） | **这条我们自己已记账并写明后果**，不算新发现；新发现是他们的 **`undeliveredOnly`＋`delivered_at_ms` 复合索引**这个形状——把"还没告诉用户"做成可查列，比"事件发了就忘"诚实 |
| 取消 | `stop()` 三段（stopping→runner→canceled），失败也标 `stopFailed` | 有取消路但是**单路径**：`TaskRoster.cancel` map ＋注释 `cancel is the ONE stop path a subagent has ... No second cancel mechanism exists for a roster row by design (197 §0)`（`internal/tools/task.go:174-177`）；另有 `internal/tools/cancel.go`（82 行，同 `wc -l` 口径） | 我们缺的不是取消能力，是**"取消失败"这件事的表达**：他们没有把 stop 失败谎报成 canceled，而是写进事件 payload |
| 忙时排队 | `BusyQueue` latest-wins＋30 分钟 TTL＋五种排空结局 | 会话级排队由 `internal/agent/scheduler`（C20 TaskScheduler＋PathLock）负责，**不是"定时任务撞上正在跑的会话"这一形**〔现读：`internal/agent/scheduler/doc.go:1-15`，职责逐字 `concurrent task management ... PathLock ... scheduler door`〕 | D12 落地时若真出现"9:00 的活儿撞上用户在聊"，我们没有 latest-wins/expired 这套结局枚举，现成一套可抄 |
| 可回收删除 | 权限层把 `rm` 改写成 trash（见 §2.5），**并自承 Windows 的 `del/rd/Remove-Item` 还没接**（`dangerous-patterns.ts:140-145`） | **我们反倒更完整**：`internal/tools/recycle_windows.go`（268 行）走 `SHFileOperationW` ＋ `FOF_ALLOWUNDO`，并**回读 `$Recycle.Bin\<SID>\$I*` 恢复记录证明真的进了回收站**（`:17-28` 注释逐字 `os.Remove is never reached from this file: an unlink is not a trash`） | 这一格**是我们领先**，且他们的注释正好承认这是他们的 Windows 缺口——写进差集时要反向标 |

---

# 4. `session-report`（4 枚 `.ts`，620 行）—— **没有推翻**编排者那句"五家都没有日报"

## 4.1 直接回答（先给结论，再给证据）

**结论：`session-report` 这枚模块不是"今天它替我做了哪些决定"的日报。它没有推翻 `A409` 第九节，
也没有推翻 `docs/reports/missing-features-2026-09-28-v2.md` 第九节那一行。**

它其实是**"会话诊断包／反馈工件收集器"**：把一棵会话子树目录里的文件逐个读出来汇成一枚 manifest，
交给上层去上传。三条硬证据（都在模块自己身上）：
1. `package.json` 的 description 逐字：`Host-independent local Session diagnostics and feedback artifact collection.`
   （`session-report/package.json:5`）——**diagnostics ＋ feedback artifact**，没有一个字提到"决定"或"摘要"。
2. 收集器的输出是**文件清单**，不是叙述：`SessionReportManifest = {schemaVersion:1, rootSessionId, sessionIds[], artifacts[]}`，
   每枚 artifact 只有 `{name, bytes, required, path|content}`（`src/contracts.ts:22-40`）。
3. **消费者是反馈上传链**〔越界读物，`grep -rln --include=*.ts "session-report|SessionReportCapability|session_report" packages | grep -v agent-modules/session-report`＝7 命中：
   `packages/agent-extension/src/session-report.ts`、
   `packages/local-runtime-v2/src/application/session/process-local-application.ts`、
   `packages/local-runtime-v2/src/service/session-system/{index,owner}.ts`、
   **`packages/tui/src/runtime/feedback/diagnostic-summary.ts`**、**`packages/tui/src/runtime/feedback/diagnostic-upload.ts`***
   ——落点在 `tui/runtime/feedback/`（用户点"提交反馈"那一条路），不是任何"日报"入口。
4. 反向佐证：全仓 `grep -rniE "daily|digest|recap|activity report|decision log|decision summary" packages --include=*.ts`
   的**非测试命中全部落在 `goal/` 的 sha256 digest 与 `goal/continuation.ts` 的 "audit" 提示词上**（见 4.3），
   **没有一枚"今日汇总"类的实现**。⇒ **A409 第九节维持原判断。**

## 4.2 但这枚模块里有三件我们该抄的机关（与日报无关，都带行）

| 机关 | 细节（带行） | 我们的现量 |
|---|---|---|
| **上报"我没能收到什么"** | 每枚失败都进 `skipped[{name, reason}]`，最后**把这份 skip 清单自己写成工件** `session-report-collection.json`（`session-report-service.ts:49-53`、`:93-101`）；单枚文件失败不影响整包（注释 `Collects readable artifacts independently; diagnostics must survive partial history loss.` `:43`）；**父目录坏了不许遮蔽子节点**（`:128-136` 注释逐字 `A missing parent directory/identity must not hide discoverable descendants`）；重复名另记 `duplicate_artifact`（`:90`）；`manifest.json` 缺失单独报 `manifest_missing`（`:150-159`） | 我们的同类形状是"宁缺毋造"那一类（票 145 快照扩字段），但**没有一枚"收集器自报漏收"的工件**〔现读：`grep -rniE "skipped" --include=*.go internal/ | grep -v _test`＝**50 命中**，抽读前 6 条全是注释／测试跳过语义（`audio/mmdevice_windows.go:47`、`ball/ball_windows.go:48`、`llm/adaptertest/harness.go:277`、`adaptertest/mockllm.go:50`、`:54`、`llm/openairesponses/stream.go:320`），**无一枚是"工件收集漏收清单"**；**50 条未逐名看完，故此条判据强度＝抽读**，登记在 §①〕 |
| **拒绝符号链接＋realpath 收敛** | 遍历会话目录时：**遇到 symlink 直接抛** `session_report_unsafe_symlink`（`:176-180`），文件还要 `realpath` 后验 `pathInside(root, resolvedPath)` 否则抛 `session_report_path_escape`（`:190-194`、`pathInside()` `:261-264`）；目录入口先 `lstat` 判"是目录且不是 symlink"否则 `session_report_history_missing`（`:247-255`）；`.DS_Store`／临时发布件／`llm-call`／`llm-context-inspector` 目录**白名单式忽略**（`:31-41`、`:228-238`） | 我们更严也更成契约：**在 `risk.PathResolver` 之外用 `filepath.Clean\|Abs` 做 FS 决策本身就是 D22 禁形**（AGENTS §1.2），且 junction/reparse 有专测〔现读：`internal/risk/pathresolver_junction_windows_test.go`、`syncdirs_ancestor_reparse_ticket105_windows_test.go` 存在〕。⇒ **这一格我们不欠他们**，但他们的"临时发布件名单"（`:37-41` 三枚正则，逐形钉 `.<pid>-<uuid>.tmp`、`history-catalog.json.<pid>.<n>.<hex>.tmp|restore`、`manifest.json.<pid>.tmp`）值得抄成我们的**忽略名单测试夹具** |
| **把"模型这一轮被喂了什么"当证据存** | `SessionLlmCallEnvelope` 是**"不含消息正文、不含密钥"的重建入参**（`contracts.ts:52-67`，注释逐字 `Secret-free, non-message inputs needed to reconstruct an agent LLM call.`）：`systemPrompt · tools[{name,description,parameters}] · model · provider · api · thinkingLevel · maxTokens · hostMaxOutputTokens · maxSerializedInputBytes · cacheRetention · outputSchema · outputRevisionInstruction`。存储侧：规范 JSON（`serializeCanonicalJson`）＋**sha256 指纹去重**（`llm-call-report-store.ts:47-53`）、上限 **4 MiB**（`:14`）、指纹表上限 **256**（`:15`）、原子替换（`:52`）、**压缩完成后把当前信封冻结在它描述的那枚快照旁边**（`:57-75`）、**Rewind 之后删除证据但注释说明"删的是上一轮的证据不是运行时配置"**（`:77-101`，逐字 `A committed Rewind invalidates the last call's evidence, not Runtime configuration.`）、按会话分车道串行（`:36` `lanes`、`inLane()`） | 我们有 prompt 落盘与 spill（`internal/agent/prompt.go`、`internal/agent/spill.go`），**但没有一枚"这一轮的运行时配置信封"的独立工件**〔现读：`grep -rniE "thinkingLevel|cacheRetention|maxSerializedInputBytes" --include=*.go internal/ | grep -v _test`＝**3 命中**，
全部是**模型能力枚举**（`internal/config/schema.go:357-359` 的 `ThinkingLevels []string \`toml:"thinking_levels"\``
与 `internal/config/validate.go:201` 的校验循环），**不是"这一轮实际用了哪一档"的落盘读数**；
`cacheRetention`／`maxSerializedInputBytes` 各自 **0 命中**〕。⇒ **这是本轮第二值钱的发现**（仅次于 runaway-guard 的六信号）：出问题时我们要回答"当时到底喂给模型什么配置"，今天只能靠代码反推 |

## 4.3 最接近"日报"的东西在 `goal` 里，而且形状不同

`packages/agent-modules/goal/src/continuation.ts` 把"完成审计"写成了**给模型的提示词模板**，不是给人的报表：
- `:50` `Completion audit:`、`:59` 逐字 `The audit must prove completion, not merely fail to find obvious remaining work.`
- `:63` `Blocked audit:`、`:67` 逐字 `a safety/policy refusal remains immediate. For every other blocker, treat the resumed run as
  a fresh blocked audit and call update_goal with mode "status" and status "blocked" only if the same condition
  repeats for at least three consecutive resumed goal turns.`
- `:90` `DEFAULT_GOAL_TERMINAL_AUDIT_TEMPLATE`、`:97` 恢复版模板、`:163-170` `renderTerminalAuditPrompt()`
- `:73` `Do not call update_goal unless the goal is complete, the explicit safety/policy refusal case applies, or the
  strict blocked audit above is satisfied. Do not mark a goal complete merely because you are stopping work.`

⇒ **性质差别要说清**：这是"**让模型自己写检查**"，不是"**把已经发生的决定摆给人看**"。
前者不可信（模型自报），后者要数据源（我们的 `tool_call` 表已有 `decision/outcome/error_class`，见 §1.4 末行）。
**owner 要的日报缺的是后者**，所以：**A409 第九节这一条不被推翻，但要多写一句**——
MiniMax 家至少有**四件可当日报数据源的现成形状**，别的家我本轮没查：
1. `permission` 的决策理由树 `DecisionReason` 14 变体（§2.1）＋ `PermissionCheckerDecisionTrace`（§2.2 末）；
2. `background-task` 的事件表 `local_runtime_background_task_events`（§3.2）与 `undeliveredOnly`（§3.2）；
3. `runaway-guard` 的 `RunawayGuardTurnSummary`（六枚信号各自 `episodeCount/maxOccurrences/firstStepIndex/lastStepIndex`，
   `runaway-guard/src/contracts.ts:120-136`）——**"今天它被拦了几次、各拦在哪一类"这件事在他们这里是每 Turn 一枚结构化摘要**；
4. `cron` 的每次执行历史 `<cronName>.sessions.json` ＋ `keepSessions` 保留策略（`cron/src/executor.ts:393`、`:427`、`:519`）。
   ⇒ **给编排者的实际后果**：日报这件事的难点也许不在"要不要做一张报表"，而在**"有没有一枚每 Turn 的结构化摘要"**。
   他们至少有 3 枚可当原料（runaway 的 `RunawayGuardTurnSummary` ＋ permission 的 `DecisionReason` 双语理由 ＋ 后台任务事件表）。
   **我们侧我这轮只量到"版本号那一维缺失"**（`policy_version`／`strategy_version` 全仓 0 命中，见 §1.4 末行与本节上方），
   **没有做"是否存在任何同类每轮摘要"的全仓反向取证** ⇒ 这一条请按弱版本引用：**"缺的是可归因的版本号与逐信号计数，
   不是断言我们一枚摘要都没有"**（详见 §①D 第 5 条）。

---

# 5. `system-reminder`（12 枚）＋ `context-manager`（9 枚）＋ `conversation-contract`（2 枚）—— 上下文与提醒的分层

## 5.1 他们把"提醒"做成独立模块，内部是**责任链＋每 provider 单独频控**

| 机关 | 细节（带行） | 层 |
|---|---|---|
| 责任链注册表 | `SystemReminderRegistry`：`append(name, fn)` 全局、`appendCritical(...)` 同链但**关掉全量提醒时仍跑**、`appendFor(frameworkType, ...)` 按框架、`appendCriticalFor(...)` 两轴叠加；`resolve()`＝defaults ＋ 该框架专属（`system-reminder/src/registry.ts:25-73`），`resolveCritical()` `:76-78` | 进程内 |
| 链的语义 | 注释逐字（`:1-11`）：provider 收到 `SystemReminderInput` 返回**字符串块或 undefined 表示跳过**，按登记顺序执行，输出汇成 `<system-reminder>`；`Different AgentFrameworkType can register different sets of providers` | 进程内 |
| **按 provider 的频控策略** | `ReminderPolicyEntry{name, frequency:{intervals[]}, thresholds:{cooldown_ms, todo_reminder_interval_turns}}`（`types.ts:409-422`）。语义写在 `:378-395`：`undefined`＝全跑（旧路径）、**`[]`＝fail-safe 全关**、非空＝只跑名字匹配的，且**`critical` 不被绕过**（逐字 `critical providers are NOT bypassed; they go through the same allowlist`） | 进程内（配置来自 `AgentConfig.system_reminders`） |
| **指数退避的提醒** | `EVOLUTION_CONFIG.MEMORY_REMINDER_INITIAL_INTERVAL = 10`、`MAX_INTERVAL = 40`（注释 `Max interval after exponential back-off (caps doubling)`），调度注释 `First fires at turn INITIAL_INTERVAL (10) / If ignored, gap doubles: 10 → 20 → 40 (capped at MAX) / Preserves backoff pace when agent writes memory`（`types.ts:426-431`、`evolution.ts:44-52`）；状态是按会话的 Map **且带容量上限** `boundMap(map, MAX_MAP_SIZE=200)`（`evolution.ts:21-29`） | 进程内，有状态 |
| 提醒有多少种 | `blocks.ts` 导出 **30 枚** 函数（推导式：`grep -c "^export function" system-reminder/src/blocks.ts`＝**30**，其中 29 枚是 `build*Block`/`build*Reminder` 形状，另 1 枚是 `formatPeersAggregated` `:25`）。可读名的有：`buildAgentContextBlock`/`buildSlimAgentContextBlock`（**同一内容两版浓淡**，`:89`/`:206`）、`buildPeersUpdateBlock`、`buildMemorySkillReminder`、`buildMemoryTopicsBlock`、`buildIdentityUpdateBlock`、`buildConfigUpdateBlock`、`buildBootstrapBlock`、`buildEvolutionReminderBlock`、`buildRelevantMemoryBlock`、`buildProactiveMemoryBlock`、`buildUserMemoryUpdateBlock`、`buildAgentMemoryUpdateBlock`、`buildMemorySummaryUpdateBlock`、`buildDailyMemoryUpdateBlock`、`buildPersonaMissingBlock`、`buildInboundMetaBlock`、`buildSkillsBlock`、`buildSkillEvolutionChannelsBlock`、`buildSkillSignalReportingBlock`、`buildActivePlanReminderBlock`、`buildAsyncAuditBlock`、`buildMediaOutputReminderBlock`、`buildTaskCompletionReminderBlock`、`buildBoardNudgeBlock`、`buildWorktreeReminderBlock`、`buildSecretEnvBlock` | 进程内 |
| **`buildSecretEnvBlock` 只报变量名不报值** | 输入是 `secretNames: readonly string[] | undefined`，三档语义写在注释（`types.ts:360-377`）：`undefined`→没接线，保持沉默；`[]`→库是空的，**逐字 `no point reminding the agent about nothing`** 也沉默；非空→渲染 `<secret-env>`，**且只在第 1 轮或名字集合变化时注入（按会话 hash）** | 进程内 |
| 会话/角色枚举自成一家 | `SessionType{Branch=0, Root=1}`（注释 `Numeric values are stable forever (persisted in sessions.session_type)`，`types.ts:21-32`）、`Role{User,Agent}` 与协议层 speaker role **明确区分**（`:34-41`）、`AgentRole{Worker, Orchestrator}`（`:43-47`）、`AgentFrameworkType{OpenCode(deprecated), PiAgent, Codex}`（`:49-55`） | 进程内 |

## 5.2 `context-manager`：预算是七个具名字段，估算**承认自己会错并锚定真值**

`DEFAULT_CONTEXT_MANAGER_SETTINGS`（`context-manager/src/settings.ts:3-11`）逐值：
`enabled:true · reserveTokens:16_384 · keepRecentTokens:20_000 · minMessagesToCompact:4 ·
contextWindowFallback:128_000 · safetyMarginTokens:2_048 · strategyVersion:'agent-modules-context-manager-v1'`

| 机关 | 细节（带行） |
|---|---|
| **触发线算法** | `computeCompactionTriggerAt()`：MiniMax-M3 在 512K/1M 模式走**产品定义的 90% 线**（`:25-30`），其余模型 `contextWindow − max(reserveTokens, perTurnMaxTokens + safetyMargin)`，下限 1（`:31-36`）；注释逐字 `Other models reserve the greater of the configured reserve and one turn's output budget plus the safety margin`（`:13-17`） |
| **供应商入参上限** | `resolveCompactionTokenBudget()`：`providerInputLimit = min(0.95×window, window−reserve, window−effectiveOutput−safety)`（`provider-budget.ts:26-33`，`PROVIDER_INPUT_RATIO = 0.95` 在 `:3`）；自动触发线再取 `min(providerInputLimit, window − min(2×reserve, ⌊window/4⌋))`（`:34-41`）——**主动压缩线早于硬上限线，且随窗口比例封顶** |
| **输出预算随上下文收缩** | `resolveDynamicMaxTokens()`：剩余 `window − estimated − safety`，但不低于 `outputFloor`（floor 自身封顶为 reserve，`:44-46`）（`:5-15`） |
| **CJK 估算不足是被点名的问题** | 文件头逐字（`token-estimator.ts:6-11`）：`pi-agent-core's estimateTokens / estimateContextTokens use Math.ceil(chars / 4), which severely under-counts CJK text — Chinese characters consume 1–2 BPE tokens each, so chars/4 (≈0.25 tokens/char) under-reports by 4–8×. Under-reporting moves the compaction trigger to the right, past the point at which Messages-compatible providers reject the request (e.g. minimax's MiniMax-M2.7: rejects when input > contextWindow - max_tokens)` |
| **估算错只允许发生在尾部窗口** | `:20-27` 三条：找最后一条**成功且带 usage** 的 assistant 消息 → **信它上报的 total 作为前缀和** → 只估它之后的消息；`getLastAssistantUsageInfo()` `:165`、`calculateContextTokens()` `:137` |
| 估算器是**注入缝** | `TokenEstimator` ＋ `EncodeFn`（`:71`、`:96`），注释 `A future count_tokens-API-backed estimator can drop in via constructor injection without touching manager logic`（`:13-18`）；默认 `BpeTokenEstimator` 用 `gpt-tokenizer` 的 `o200k_base`，**故意选偏高的方向**：逐字 `consistently over-estimates a touch (safe: earlier trigger), which is the direction we want when capacity miscalculation means a hard 4xx`（`:31-36`） |
| 结构开销与图像 | 每消息 4 token 的角色包装开销（`:38-41`）；视觉块按 **4_800** token 计费（`:43-45`、`isVisualContentBlock` `:133`）；超长字母数字连续串走上界（`hasOversizedAlphanumericRun` `:101`、`estimateTextTokensUpperBound` `:121`） |
| 压缩是**决策联合** | `ContextCompactionDecision = continue | skip(reason) | abort(reason) | replaceMessages(messages, metadata)`，metadata 含 `replacementId · strategyVersion · summary · compactedMessages · keptMessages · firstKeptIndex · tokensBefore/After · messagesBefore/After`（`types.ts:73-92`） |
| **摘要要带"usage 已失效"的计数** | `ContextCompactionSummaryMessage.keptMessageCount`，注释逐字 `Number of following messages retained from the old context, whose usage is stale. Persist this with the summary so a restored transcript has the same usage boundary. Absent on legacy summaries, which must use timestamp-based freshness instead.`（`types.ts:9-16`） |
| 四枚观测事件 | `ContextTriggerEvaluatedEvent`（含 `triggerAt/shouldCompact`）、`...SkippedEvent`（含 `reason`）、`...CommittedEvent`（前后两侧数）、`...FailedEvent`（含 `reason`）（`types.ts:94-127`），`ContextManagerObserver` 四者全可选（`:129-134`）；每枚带 `phase` 字段（分阶段归因） |
| 压缩前要落检查点 | `ContextManagerCheckpointOptions.beforeCompaction?: () => boolean | Promise<boolean>`，注释逐字 `Called only after compaction is required and a safe plan exists. False defers this attempt.`（`types.ts:146-149`） |
| 会话级串行＋计数来源 | `ContextManagerLock.withSessionLock(sessionId, action)`（`:60-62`）；计数结果带来源 `'remote_count_tokens' | 'local_estimate' | string`（`:28-31`，**来源是值的一部分**） |

## 5.3 `conversation-contract`：只有 2 枚文件，但它是三家里最像"冻结契约包"的一枚

- 口径：`find conversation-contract -name '*.ts'`＝2，`index.ts` 739 行＋`process-local.ts` 12 行。
- **`process-local.ts` 全文 12 行就是一句边界声明**，逐字两行：
  `/** Caller-controlled cancellation context; accepts no network identity or HTTP headers. */`（`process-local.ts:1`）
  ＋`/** Session frames consumed in-process; consumers must end the iterator early when needed. */`（`:5`）
  ⇒ **这一枚正面回答我们 `Q-49`（面板侧来源的 L2「允许」）**：他们不是"给面板加一枚来源字段"，
  而是**把"进程内入口"单独立成一份契约文件，并在类型文档里禁止它携带网络身份**。
  我们现量（逐名，别把这条读成"我们完全没有面板边界"）：
  `internal/panel/bridge.go` **存在**，`internal/panel/doc.go:2` 与 `:9` 把 C17 的职责写成
  `PanelBridge (C17) with method whitelist` ＋ `PanelBridge method whitelist + capability annotations + correlationId`
  （另见 `docs/specs/SPEC-08:156` 的 §5.2 与 `:235` 的追溯项）；
  面板侧出口在 `internal/agent/approval/ui.go`；`internal/agent/approval/gate.go:595` 注释里还有一句
  `a future PanelBridge`（说明那枚类型今天**还不叫这个名字**）。
  ⇒ **真正缺的是他们那种"入口层就不接受网络身份"的类型级声明**：我们的白名单是**方法名**层面的（能调哪些方法），
  他们这 12 行管的是**身份来源**层面的（这个入口不许从 HTTP header 取调用者）。两者不互相替代。
- **来源枚举 13 枚**（`index.ts:1-14`）：`api · cron · task · background-task · team · thread-goal · questionnaire ·
  communication · code_review · greeting · channel:wechat · channel:feishu · channel:telegram`
  ⇒ "这条消息是谁发起的"是**契约字段**（`ConversationSubmitInput.source`，`:357`），不是日志里的自由文本。
  我们侧〔现读：`grep -rniE "ConversationSource|MessageSource|MsgSource" --include=*.go internal/ | grep -v _test`＝**0 命中**；
  最接近的是 `internal/risk/provenance.go` 的 mark/origin，那是**内容**来源不是**触发**来源〕。
- **会话状态五枚** `idle | started | error | aborted | interrupted`（`:24-29`）与会话种类六枚
  `conversation | task | peek | channel | cron | unknown`（`:16-22`）**分成正交两轴**。
- **Turn 拒绝原因 12 枚，含一枚可展开的前缀**（`:319-331`）：
  `invalid-session · active-turn · compaction-active · session-deleting · session-mutating · priority-blocked ·
  ingress-conflict · policy:${string} · duplicate · invalid-input · active-model-override-unsupported · delivery-closed`。
  ⇒ **`policy:${string}` 这种"封闭枚举＋一条开放前缀"的写法**（TS 模板字面量类型）既机读又可扩，
  是我们 C 契约注释里可以借的形状（Go 侧等价物＝前缀常量＋校验函数，`internal/config/permmode.go` 一类已在做拼写校验）。
- 幂等与去重是一等字段：`dedupeKey · clientRequestId · requestedTurnId · expiresAt · allowQueue · queuePlacement:'front'`
  （`:355-366`，注释 `Internal control-plane continuation that must precede already queued ordinary input`）；
  steer 走 `producerId + idempotencyKey`（`:394-407`），有**投递前验收钩子**
  `ConversationSteerPreDelivery.accept({mode:'activated'|'steered', turnId})`（`:409-414`），
  结果分 `steered | duplicate | activated` 三态且**只有 activated 才持有 completion**
  （`:416-430`，注释逐字 `activated means the steer had no Turn to join and started one instead, so it owns a completion
  just like submit`）。
- **断点续读契约**：`ConversationResumeUserInput`（`:368-382`）注释逐字
  `Resumes a persisted user-input wait. This is a narrow internal admission seam: it never joins the ordinary Queue,
  and requestId is its durable replay identity.`，还带 `owner?: UserInputResumeOwner{kind:'thread-goal', goalId}`
  （`:377-378`、`:384-387`）。
  ⇒ 这条正面撞上 AGENTS §2 的第一枚待定项（`AwaitingApproval` 遇系统挂起、`PLAN.md:1530` DEFERRED，
  倾向"作废并判拒绝"）：**他们的做法是"续读走一条独立的窄入口，用 requestId 当重放身份"，
  而不是"恢复一个陈旧的批准请求"**——两条路都满足"不复活陈旧批准"，但多给了"挂起的用户输入可以被合法续上"这一格。
  **这不改变我们的定案，只补一条外部先例。**
- 错误字段五枚（`errorMessage/errorCode/errorSource/errorDetail/errorProviderId`，`:104-108`）——D37 跨边界传播的现成形状。
- ⚠ 一处**我们绝不能照抄**：`ConversationRunLocation.mode = 'current' | 'new-worktree' | 'existing-worktree'`
  （`:79-85`，还带 `resolvedBranch/parentRepoDir`）——**他们把 git worktree 当成会话的一等运行位置**。
  我们 AGENTS §1.4 逐字禁 `绝不在仓库目录内建 worktree 或 checkout`，这条要么明确不接、要么走 D17/D41 单独议。

## 5.4 三者分层对照我们的 D15／D39

| 层 | 他们 | 我们（现量） | 后果 |
|---|---|---|---|
| 提醒（provider 链） | 独立模块、30 枚 block 构造函数、按 provider 频控＋指数退避＋critical 档 | 只有一种提醒（C22 重复调用），注入在 `internal/agent/loop.go`〔现读：`grep -rn "EvReminder" --include=*.go internal/ \| grep -v _test`＝**3 命中**；全仓非测试 `reminder` 命中分布＝`guard.go` 9／`loop.go` 6／`sink.go` 2／`statemachine/table.go` 1（共 18）；**按提醒名的频控/退避表为 0 命中**：`grep -rniE "reminder.*(interval\|cooldown\|backoff)" --include=*.go internal/`＝**0**〕 | 我们要加"计划提醒/记忆提醒/身份提醒"这类**得一条条写代码**；他们是一条链＋一份配置 |
| 上下文预算 | 7 个具名字段＋两条线（主动线／供应商硬线）＋动态输出预算 | `internal/agent/budgets.go`（150 行）：按 128k 参考窗口**等比缩放**、`Scale` 夹在 `(0,1]`（`:16-18`、`:51-55`）；参考值逐名 `refPromptTotal=2300 · refSpillTokens=4000 · refSpillHead=500 · refSpillTail=200 · refHistoryTokens=12000 · KeepRawRounds=3 · MaxToolConcurrency=4`（`:22-47`），注释逐字 `D15 makes scaling by the configured model's context_window a MUST, not an optimization`（`:7-11`） | **两种哲学**：我们"随窗口等比缩小、大窗口不抬高冻结值"（D15 原文要求），他们"绝对余量＋按模型开特例（M3 的 90% 线）"。他们那枚 `strategyVersion` **我们没有**〔现读：`grep -rniE "strategyVersion|strategy_version" --include=*.go internal/`＝**0 命中**〕 |
| token 估算 | 宁高不宁低（早触发）＋**锚定 provider usage 前缀和**＋来源标注 | `ApproxTokens(s) = len(s) / 4`（`internal/agent/budgets.go:125-129`，注释逐字 `deliberately matches llm.Request.EstimateTokens (utf-8 bytes / 4) so the loop's budget math and the seam's pre-flight estimate cannot disagree; exact accounting comes from Usage events (C23)`）；`Compressor.TotalTokens` 逐消息累加估算（`internal/agent/compress.go:110-116`），**不锚定上一条真 usage** | ⚠ **本轮我认为最该开一枚票的技术风险**（下面的推导是**我本人算的**，不是他们的原文）：中文 UTF-8 每字符 3 字节 ⇒ `bytes/4 ≈ 0.75 token/汉字`，而 BPE 实际约 **1–2 token/汉字**（他们原文的量级），
  ⇒ **我们低估约 1.3–2.7 倍**；他们把 `chars/4`（≈0.25/字）低估 4–8 倍这件事**写进源码注释并点名会造成 hard 4xx**
  （`token-estimator.ts:6-11`）。Wisp 是中文优先产品，**每一条预算阈值都经过这枚函数**（`budgets.go:125-126` 自述），
  压缩触发（`compress.go:119-121` 比 `HistoryCompressTokens`）与 guard 的 token 预算（`guard.go:148-152`）都在射程内 |

---

# 6. `plugin-hooks`（9 枚 `.ts`，5 194 行）—— 插件能钩哪些点（对 D19／D46）

## 6.1 钩点名册（11 枚，逐字来自数组）

`PLUGIN_HOOK_EVENTS = ['SessionStart','SessionEnd','UserPromptSubmit','PreToolUse','PermissionRequest',
'PostToolUse','SubagentStart','SubagentStop','Stop','PreCompact','PostCompact']`
（`plugin-hooks/src/contracts.ts:2-14`）。`SessionStart` 还带来源六枚
`startup | resume | clear | compact | fork | plugin_activation`（`:24-30`）。

**对照我们**：`grep -rniE "PreToolUse|PostToolUse|UserPromptSubmit|SessionStart|SubagentStart|HookEvent" --include=*.go internal/ cmd/`
＝**0 命中**〔现读〕。全仓 `hook` 一词非测试命中 **76** 条，按文件分布是
`internal/proc/shutdown.go` 12、`internal/tools/bridge.go` 7、`cmd/wisp/run.go` 7、
`internal/tools/fs_write.go` 5、`internal/proc/boot_windows.go` 5、`internal/agent/compress.go` 5、
`internal/ball/hotkey_reload.go` 4 …（抽点统计）⇒ **全部是 Windows 键盘钩子／关机钩子／桥回调一类 OS 语义，
没有一枚是"插件事件钩子名册"**。⇒ 我们的 D19/D46 目前**没有"插件可参与决策"这一层**（这是现量结论，不是推测：上面两把尺都是 0/OS-钩子语义），
这既是安全优势（攻击面小）也是能力缺口（插件不能做合规预处理）。**要不要开这一层是契约级问题，不是票级问题。**

## 6.2 handler 是**子进程**，形状直接对应我们的 D46

`PluginHookCommandHandler`（`contracts.ts:32-55`）逐字段：
`kind:'command' · sourceFormat:'MINIMAX'|'CLAUDE'|'CODEX' · pluginName · pluginRoot ·
pluginDataDir（注释：`Host-owned writable state directory; the Plugin package root is immutable`）·
activationKey（`Changes whenever this Plugin capability is updated or re-activated`）· sourcePath · event ·
matcher · command · args（缺省即 shell 形式）· shell:'bash'|'powershell' ·
condition（注释 `Compatible tool-input predicate such as Bash(rm *)`）· timeoutMs ·
additionalContextLimit（**按 handler 的模型可见上下文上限；`0 disables spilling`**）· declarationOrder`

执行层（`runner.ts`）——**层＝子进程**，这是三家里给我们最现成的一份 Windows 取消作业：
- `spawn` 自 `node:child_process`（`:1`、`:262`），`detached: process.platform !== 'win32'`、`windowsHide: true`（`:266-267`）；
- **Windows 上的树杀是 `taskkill /pid <pid> /T /F`**（`runTaskkill()` `:520-529`，`error` 事件折成 false `:526`），
  POSIX 走 `process.kill(-pid,'SIGKILL')`（`:513`），**两条都只给 500 ms 排空预算**
  （`PROCESS_DRAIN_BUDGET_MS = 500` `:38`，用法 `:504`、`:510`）；
- **输入/输出硬闸**：`MAX_INPUT_BYTES = 1 MiB`（`:31`，超限就不写 stdin `:1702`）、
  `MAX_OUTPUT_BYTES = 64 KiB`（`:30`，溢出即判 `HOOK_INVALID_OUTPUT` 并杀进程 `:418-429`、`:289-292`）、
  `MAX_REASON_CHARS = 4_096`（`:32`）、`MAX_OUTPUT_VALUE_NODES = 10_000`（`:35`）；
- 事件级预算：`ORDINARY_EVENT_BUDGET_MS = 15_000`、**`SESSION_END_BUDGET_MS = 3_000`**（`:36-37`，收尾事件预算更少）；
- 结局码六枚 `HOOK_ABORTED | HOOK_INVALID_INPUT | HOOK_INVALID_OUTPUT | HOOK_PROCESS_ERROR | HOOK_PROCESS_EXITED | HOOK_TIMEOUT`
  （`:193-201`），每枚带 `{pluginName, sourcePath, event, declarationOrder}`（`:57-63`）
  ——**出错能定位到"哪个插件的哪份文件第几个声明"**；
- 观测器两轴 `onHandler({event,format,outcome,durationMs,processKilled})` 与 `onEvent({event,outcome,durationMs})`；
  outcome 词表 handler 轴是 `success|deny|ask|error|timeout|aborted`、事件轴是 `success|deny|ask|degraded|timeout|aborted`
  ——**事件轴多一枚 `degraded`**（`:232-245`）；
- `condition` 的解析是一条具名正则：`/^([A-Za-z0-9_.:/-]+)\(([^\r\n()]*)\)$/u`（`:1832`）；
- matcher 安全性单独成函数并导出：`isSafePluginHookMatcher`（`index.ts` 第 1 行导出，来自 `parser.ts` 417 行）。

## 6.3 最要紧的一段：**插件钩子能改哪些东西**（对我们 D19 是一张现成红线清单）

`PluginHookDecision`（`contracts.ts:155-191`）逐字段后果：

| 字段 | 后果（带行） |
|---|---|
| `decision: allow\|deny\|ask\|defer` ＋ `permissionDecision: allow\|deny\|abstain` | **插件可以替用户批准**，注释逐字 `an explicit allow authorizes the tool **without showing the product UI**`；缺省输出＝`abstain`（`:156-161`） |
| `permissionAutoApproval: 'ordinary_only'\|'any_prompt'` | 批准的范围；注释点明两家语义不同：`Compatible still evaluates deny/ask rules; Codex treats the Hook verdict as the final approval`（`:162-166`） |
| `updatedInput` | **改写工具入参**（`:171`） |
| `updatedResult` ＋ `updatedResultFormat` | **改写模型可见的工具结果**，注释逐字 `Replacement for the model-visible tool result; **audit/persistence keeps the original result**`（`:187-190`）——**"给模型看的"与"入账的"允许不一致，但必须显式声明格式归属** |
| `updatedPermissions[]` | **改权限规则**：`addRules/replaceRules/removeRules`（带 `behavior` 与 `destination`）、**`setMode`（档位词表含 `bypassPermissions` 与 `plan`）**、**`addDirectories/removeDirectories`（扩目录授权）**（`:137-153`、`:120-126`）；`destination` 四档 `session\|localSettings\|projectSettings\|userSettings`，注释 `cliArg is intentionally not host-persistable`（`:113-118`）；**宿主必须整条原子应用**（`:133-136`） |
| `interrupt` / `continue:false` / `stopReason` | **停掉正在跑的 agent**；两条语义分别写在 `:174`（`deny-only request to stop the active agent run`）与 `:178-179`（`Universal ... false stops the active agent run`） |
| `systemMessage` / `suppressOutput` / `postToolFeedback` | 给用户看一条／压掉输出／**给下一步模型一条反馈**（注释 `never an agent hard stop`，`:181-186`） |
| `terminalSequence` | **写终端转义序列**，但出口有白名单函数 `isAllowedPluginHookTerminalSequence`（字段 `:183-184`；`index.ts` 导出该函数） |
| `defer` | 把决定推后再问（`:177`），与 `decision:'defer'`（`:108`）配套 |

**他们自己钉的一道闸（对我们 Q-49 直接可用）**：`runner.ts:781` 逐字
`if (sourceFormat === 'MINIMAX' && decision.updatedPermissions !== undefined) return undefined;`
——**自家格式（MINIMAX）的 hook 不许带 `updatedPermissions`**，只有 `CLAUDE` 兼容格式才被解析翻译
（`:782-796`、`parseCompatiblePermissionUpdates` 里逐类校验 `setMode` `:844-848`、
`addDirectories|removeDirectories` `:851`），且 `:799` 对 `updatedInput`/`updatedPermissions` 走另一条分支。
⇒ **同一个危险输出，按"来源格式"选择性禁用，而不是全局允许或全局禁止**；
这与我们 AGENTS §1.2 那条禁形（**由面板侧来源的 L2「允许」**）用的是同一种思路：**来源决定它能不能产生授权**。
我们缺的是他们那种**机读的来源标签**（`sourceFormat` 三枚枚举＋`toolProvenance{kind:'plugin_mcp'|'plugin_host', pluginName}`，
`contracts.ts:18`、`:95-103`，注释逐字 `plugin MCP names are scoped only when the host proves which plugin owns the server`）。

## 6.4 还有一层：`PluginHookCoordinator`（1 289 行）与"提交事务"

`index.ts` 导出 `PluginHookCoordinator`、`type PluginHookAdmissionTransaction`、`type PluginHookSessionEndFence`
三个名字（`plugin-hooks/src/index.ts` 第 12-16 行区段）——**"钩子产出的准入是一枚事务"＋"会话收尾有一道 fence"**。
另有 `cleanupPluginHookSessionArtifacts`（`output-artifacts.ts` 359 行）与
`composePluginHookToolResultContent` / `mergePluginHookContext` / `mergePluginHookDecisions` /
`renderPluginHookRejectionReminder` 四枚合成函数（`index.ts` 第 3-11 行）。
⇒ **"多个 handler 的裁决要怎么合并"与"被拒时要给模型回一句什么"在他们那里是独立导出的函数**，
不是散在各 hook 里的 if。**coordinator.ts 1 289 行与 runner.ts 2 348 行我只读了抽检段落**（登记 §①）。

---

# 7. `goal`（19 枚）／`mcp`（13 枚）／`skills`（4 枚）—— 派单未列必查，本轮＝**抽检登记**，不是逐行读完

> 口径先说清：这三枚**不在派单的六枚必查里**，我只做了"文件清单＋导出名＋抽检若干段"。
> 下面每条都带行，但**不要当成完整覆盖**（覆盖度自述见 §①）。之所以留几条，
> 是因为这三枚里各有一条我认为比某些必查项更该进我们视野的东西。

## 7.1 `skills`（4 枚 `.ts`：`index.ts` 24／`directory-watcher.ts` 60／`types.ts` 124／`registry.ts` 1099）

- **六个来源层**：`SkillSourceKind = 'project' | 'workspace' | 'agent' | 'global' | 'user' | 'builtin'`
  （`skills/src/types.ts:1`），每层带 `priority?: number`、`external?: boolean`，
  以及一枚我们一定会踩的坑：`allowDirectorySymlinksOutsideRoot?: boolean`（`:3-13`，
  注释逐字 `Follow directory links outside the root for configured compatibility sources`）
  ⇒ 他们的答案是"默认不跟根外符号链接，只有登记为兼容来源的才允许"——**可配置但默认关**。
  我们 `internal/risk/pathresolver.go:104` 的规矩是"非豁免的 reparse 遍历即错误"，**方向一致、我们更硬**。
- **同名遮蔽要留败者名单**：`SkillLoser {entry, reason, winnerLocationUri}` ＋ `SkillViewEntry.losers[]` ＋
  `SkillSnapshot {version, generatedAt, roots, entries, winners, losers, diagnostics, metrics}`（`:49-79`）。
  ⇒ 这是 skills 里我最看重的一点：**"谁的 skill 被谁的盖掉了"是可查询、可版本号化的一等数据**。
  我们侧〔现读：`grep -rniE "loser|shadow|遮蔽" --include=*.go internal/tools internal/config | grep -v _test`＝**0 命中**
  （**口径只覆盖 `internal/tools` 与 `internal/config` 两目录**，全仓未搜，故这条是"限定范围内为 0"而不是"全仓为 0"）〕。
- 变更监听是**去抖**的：`WATCH_MAX_WAIT_MS = 1000`（`registry.ts:84`），watcher 有显式 `close()`（`:49`、`:210-213`），
  按 `(目标路径 \0 rootId)` 组合键复用/替换（`:216-217`、`:257-263`）。
- 诊断分两级 `'warning' | 'error'`（`types.ts:14`）；链接型 skill 的 `entryDir` 与 `skillDir` 会不同（`:33-34` 注释）；
  数据目录来自环境变量 `MINIMAX_DATA_DIR`（`registry.ts:65`）。

## 7.2 `goal`（19 枚 `.ts`，含测试共 4 259 行；`types.ts` 335／`store-port.ts` 341／`tool-impls.ts` 361／`verification/subagent.ts` 394／`verification/evaluator-adapter.ts` 423）

- 目标语句被 **sha256 摘要**钉住：`digestThreadGoalObjective(objective) = sha256(objective)`（`objective-digest.ts:4-5`），
  `objectiveDigest` 散在 `types.ts:70/78/131`、`store-port.ts:102/133`、`verification/evidence-brief.ts:24/38/57/89`、
  `verification/verifier-port.ts:53`。
- **"账过期了"的原因码是三枚枚举**：`ThreadGoalBoundUsageStaleReason = 'missing_goal' | 'goal_epoch' | 'objective_digest'`
  （`store-port.ts:89`）⇒ 目标没了／纪元换了／目标文本改了，三种陈旧各有一个码，不共用一个"stale"。
- **完成/受阻是写给模型的审计作业**，不是给人看的报表（§4.3 结论的出处）：`continuation.ts:50` `Completion audit:`、
  `:59` 逐字 `The audit must prove completion, not merely fail to find obvious remaining work.`、`:63` `Blocked audit:`、
  `:67` 同一阻塞条件要**连续三枚恢复后的 goal 轮**才允许报 blocked、`:73` `Do not mark a goal complete merely because
  you are stopping work.`；模板常量 `:90`/`:97`，渲染器 `renderTerminalAuditPrompt()` `:163-170`。
- 校验被拆成 6 枚文件：`verification/{subagent,evaluator-adapter,verification-policy,transcript-window,verifier-port,evidence-brief}.ts`
  （口径：`ls packages/agent-modules/goal/src/verification | wc -l`＝**6**）。
  §1.4 里 `runaway-guard` 的 `shouldRemind: turnIntent?.kind !== 'goal-verifier'`
  （`local-runtime-v2/src/service/turn-system/runaway-guard/extension.ts:21`）正是**给这套校验子轮开的豁免**。
- 另有 `reply-fingerprint.ts:12`（对归一化后的回复做指纹），与"是否重复回答"这一维相关。

## 7.3 `mcp`（13 枚 `.ts`，含测试共 2 947 行；`runtime/connection-pool.ts` 911／`runtime/name-registry.ts` 258）

- **工具名是一等持久资源**：两级分配（server 段＋tool 段），落盘＋加锁＋原子替换
  （`name-registry.ts:1-14` 用 `node:fs` 的 `mkdtempSync/renameSync/fsyncSync` ＋ `proper-lockfile` 的 `lockSync`），
  状态形状 `{version:1, servers:[{key, raw, segment, tools:[{raw, segment}]}]}`（`:29-34`）。
- 三条最要紧的语义（注释逐字，`:50-54`）：
  `Host-scoped, two-level allocation. **Never derives identity from a public name.**
  Removed identities retain their assignments so a new source cannot inherit a name.
  Disk-backed instances refresh inside an exclusive transaction before allocating.`
  ⇒ **"名字永不回收"**：删掉的 server/tool 仍占着它那段，新来源不能继承旧名。
  对照我们：票 175 正在把 `task.output` 送进 C25 名册，`internal/agent/spill_name_injectivity_test.go`（**路径在此，不在 `internal/tools/`**）
  钉的是**同一轮内的单射性**〔现读：`ls` 该文件存在，62 行级；`internal/tools/` 下无同名文件〕；
  **跨时间（删除后重建）的名字不继承这一条我们无对应物**〔现读：
  `grep -rniE "nameRegistry|NameRegistry|tool_name.*unique" --include=*.go internal/`＝**0 命中**；
  `grep -rniE "名册|roster" --include=*.go internal/memory | grep -v _test`＝**4 命中**，逐名全是**协程名册**
  （`internal/memory/retention.go:19` 的 `D38b roster`、`internal/memory/writer.go:22/24/25` 的 resident 6 枚），
  **不是工具名分配器**；我们的"D34 名册"活在文档侧——`docs/PLAN.md:1531` 逐字
  `名册在册但生产注册表零实现`〕。
- 段长上限具名：`MCP_RUNTIME_TOOL_NAME_MAX_LENGTH = 80`、`MCP_RUNTIME_SERVER_SEGMENT_MAX_LENGTH = 48`
  （`tool-name.ts:2`、`:4`），配 `normalizeMcpNameSegment(value, kind)`（`:7`）、
  `buildMcpToolRuntimeName`（`:19`）、`buildMcpServerRuntimeName`（`:29`）、`shortenMcpNameSegment`（`:39`）。
- **连接池自带"跑飞"防线**：`connection-pool.ts:15` 注释逐字 `runaway child cannot grow pool memory unboundedly`、
  `:656` `safety valve flushes a runaway non-newline buffer verbatim`（超长无换行缓冲有泄压阀）。
- transport 三种由工厂选：`createTransport`／`createHttpTransport`／`createStdioTransport`（`mcp/src/index.ts:25-27`）
  ——**层＝子进程(stdio) 与 网络(http)**。
  ⚠ 与 **D13「MCP 取舍」**直接相关：`docs/PLAN.md:1533` 把 `MCP client` 挂在 **RESERVED**，逐字
  `D13 与 host bridge 收口架构冲突（第四轮再次驳回审查的引入建议，理由见 §16.9 第 7 条）`。
  ⇒ 他们这 13 枚文件是一份**可用的 MCP 运行时参考实现**（含名字治理与池内存防线），
  但**引入它＝推翻 D13**，属人工批准级，不是票级。

---

# 8. 最值得先开的 3 枚票（只写"是什么／后果／在哪一层"，拆票交给编排者）

| # | 是什么 | 不开的后果（大白话） | 跑在哪一层 |
|---|---|---|---|
| ① | **把 C22 的重复检测从"整轮签名"降到"单条调用"，并补 `same_error_family` ＋ `abab_action_cycle` 两枚信号**（参照 `runaway-guard/src/step-view.ts:153-171` 算键、`signals.ts:174` 观测第 2 次、`:199-217` ABAB、`step-view.ts:388-409` 错误家族 8 类；我们侧改的对象是 `internal/agent/guard.go:165-198` 的 `ObserveTurn` ＋ `turnSignature` `:234-243`） | 现在的 `[3,5,8]` 只要模型每轮多带一枚无关小调用就**永远不会触发**；同一类错误换个参数就**重置计数**。也就是说"它自己在转圈"这件事**今天我们只在参数一字不变时才看得见**，最常见的转圈形态（换着写法重试、两个文件来回改）全在盲区。落到 Win32＋Go 只是 `map[string]int` ＋一枚 `crypto/hmac`，**不需要常驻进程、不需要网络** | 进程内（我们 `internal/agent` 的同一层） |
| ② | **token 估算别再 `len(s)/4` 一路用到底**：改成"以 provider 上报的 usage 为前缀和、只估尾部"，并让估算偏保守（宁早触发），同时给每次计数标来源 | Wisp 是中文优先：中文 UTF-8 3 字节／字 ⇒ `bytes/4≈0.75`，实际 BPE 约 1–2 token／字，**低估 1.3–2.7 倍**（推导是我算的，方向与 MiniMax 自己在 `context-manager/src/token-estimator.ts:6-11` 写下的注释同向）。低估的直接后果是**压缩触发线右移、真到线上被 provider 硬 4xx 拒掉**（他们逐字点名 `rejects when input > contextWindow - max_tokens`），受害的正是长中文会话。这条影响 `internal/agent/budgets.go:125-129`（被每条阈值消费）、`compress.go:110-121`（触发线）、`guard.go:148-152`（token 预算） | 进程内（估算）＋**服务端**（真实 usage 来自 LLM 响应，我们已有 C23 Usage 事件） |
| ③ | **后台任务与定时唤起的"投递账本"与"重启对账"补上**：任务对象要有 `deliveredAt`、查询要能只取"跑完了但还没告诉会话"的；重启时按宿主存活证据把孤儿判成 `lost`，而不是让它停在 `running`（参照 `background-task/src/types.ts:82`、`:108-110`、`:168-170`，`local-runtime/src/background-task/schema.ts:38-45` 的 `delivered_at_ms` ＋复合索引，`startup-recovery.ts:88-106`、`:147-149`；我们侧的落点是 `internal/tools/task.go:165-189` 的进程内 `TaskRoster`） | 现在崩溃/重启后，后台任务与子代理在账上**要么消失要么留在"还在跑"的假状态**（`recycle_windows.go:32` 的注释里已经在说"后台工具调用不许挂住"，但那张表本身不落盘）。同时"完成通知到底送没送出去"无人记账——这正是票 147／票 175 那一串 offset／续读问题的同一个根：**没有一条持久的事实源，回执就永远只能靠猜**。而且 D12 的"定时交给 OS"要真能落地，`restoredTaskExecution:'quarantined'`（`cron/src/registry.ts:42-48`）这条"恢复来的档期先只登记不执行"是**必答题不是可选项** | 进程内决策 ＋ **SQLite 持久层**（我们已有 `internal/memory`）＋ **OS API**（进程存活判定：Windows 侧要的是"宿主 PID 还在不在"这类查询，对应他们 `isOwnerAlive(ownerId)`） |

**次一等但我建议排队**（各一句，不拆票）：
- 授权宽度候选枚举（`permission/src/types.ts:151-195`，index 0 必为窄档的兼容律）——直接对 D45；
- "无审批通道 ⇒ 拒绝"那一支（`permission-core.ts:422-430`）——对 D2 深度睡眠拓扑与票 197 子代理都有效；
- 拒绝文案的**双语表＋类型穷尽**（`reason-format.ts:7-9`）与**给模型的固定英文来源标签**（`:306`、`:324-338`）；
- 名字永不回收（`mcp/src/runtime/name-registry.ts:50-54`）——对 C25 名册；
- `SessionLlmCallEnvelope` 这枚"每轮运行时配置证据"（`session-report/src/contracts.ts:52-67`）——对票 174／票 136。

---

# ① 我这腿**没读到**什么（具名到模块／文件，越界读物也一起列清）

**A. 12 枚里有 6 枚是"逐行读完"，6 枚只到"抽检"**（诚实分级，别混着引用）：

| 分级 | 模块 | 说明 |
|---|---|---|
| **逐行读完**（`src/` 全部 `.ts`） | `runaway-guard`（10 枚，含 `test/guard.test.ts` 未读） | 只读了 `src/` 9 枚，**测试 666 行未读** |
| | `session-report`（4 枚全部） | 全读 |
| | `conversation-contract`（2 枚） | `process-local.ts` 全读；`index.ts` 739 行读了 `:1-150`、`:315-434`，**中段 `:150-315`（消息/工具调用/file-change 观察形状）与 `:434-739` 未读** |
| **抽检**（读了关键类型与判据函数，未逐行） | `permission` | 读了 `types/context/permission-core/engine/ask-policy/reason-format(片段)/classifier/cloud-classify-client/http-cloud-gateway-client/cloud-gateway/conversation-renderer/windows-trash-execution/dangerous-patterns(片段)`。**未读的整枚文件（18 枚，按行数从大到小具名）**：`tools/path-capability.ts`（1 535，只读了 `decideUnknownToolPathCapabilities` 段与导出名）、`tools/bash-permission.ts`（1 520，**整枚未读**）、`tools/bash-checker.ts`（1 216，**未读**）、`tools/fs-permission.ts`（1 146，**未读**）、`tools/bash-split.ts`（704，**未读**）、`tools/bash-wrapper-unwrap.ts`（841，**未读**）、`tools/bash-write-target.ts`（600，**未读**）、`tools/slow-command-scan.ts`（453，**未读**）、`tools/bash-ast.ts`（436，**未读**）、`tools/safe-command.ts`（393，**未读**）、`host-utils.ts`（340，**未读**）、`execution-plan.ts`（490，**未读**）、`tools/bash-safety-scan.ts`（566，**未读**）、`tools/bash-rule-match.ts`／`bash-fast-allow.ts`／`bash-safe-first-words.ts`／`shell-tokenize.ts`／`fs-checker.ts`／`windows-native-delete.ts`／`written-files-registry.ts`／`path-resolver.ts`／`effective-input.ts`／`legacy-permission-adapter.ts`／`locale-detect.ts`／`permission-request-store.ts`／`mcp-runtime-name.ts`／`classifier/index.ts`／`in-memory-cloud-gateway-client.ts`／`index.ts`（195）——**其中 `windows-native-delete.ts`（162）与 `written-files-registry.ts`（105）我认为最该补读，因为直接对 Windows** |
| | `cron` | 读了 `types.ts` 全、`registry.ts:14-258`、`executor.ts:60-209`、`busy-queue.ts`/`active-hours.ts` 的导出名与常量。**未读**：`executor.ts:210-777`（含会话历史与保留策略 `:393/:427/:519` 的具体实现）、`registry.ts:259-698`（update/completeOnceAfterExecution 等）、`host-ports.ts` 374 行、`host-utils.ts` 292 行、`report-delivery.ts` 211 行、`retired-cleanup.ts` 70 行（只读了导出名） |
| | `background-task` | 读了 `types.ts`/`manager.ts` 全。**未读**：`in-memory-store.ts` 172、`in-memory-output-store.ts` 129、`errors.ts` 76、`status.ts` 41、`index.ts` 65、`id.ts` 11 |
| | `system-reminder` | 读了 `registry.ts` 全、`types.ts` 的 `:1-120` 与 `:360-444`、`evolution.ts:1-70`、`blocks.ts` 与 `providers.ts` 的**导出名清单**。**未读**：`blocks.ts` 907 行正文、`providers.ts` 854 行正文、`service.ts` 260、`dependencies.ts` 169、`plugin-reference.ts` 221、`todo-state.ts` 95、`collaborators.ts` 50、`mcode-tools-master-reminder.ts` 32、`index.ts` 162 |
| | `plugin-hooks` | 读了 `contracts.ts` 全 245 行＋`runner.ts` 抽检（`:250-300`、`:415-430`、`:495-529`、`:750-800`、`:840-860`、常量区 `:30-38`、`:1702`、`:1832`）＋`index.ts` 全。**未读**：`runner.ts` 其余约 2 000 行、`coordinator.ts` 1 289、`parser.ts` 417、`wire-tool-adapter.ts` 359、`output-artifacts.ts` 359、`command-invocation.ts` 128、`effort.ts` 10 |
| | `context-manager` | 读了 `types.ts`/`settings.ts`/`provider-budget.ts` 全、`token-estimator.ts` 的头注释与导出名。**未读**：`manager.ts` 528 行正文（**压缩主逻辑，最该补读的一枚**）、`count-tokens-body.ts` 239、`context-usage-estimator.ts` 115、`index.ts` 40 |
| | `goal`／`mcp`／`skills` | **只登记＋抽检**（见 §7），`goal` 19 枚中读了 `continuation.ts`/`objective-digest.ts`/`reply-fingerprint.ts` 的抽检段与 `store-port.ts`/`types.ts` 的个别行；`mcp` 读了 `name-registry.ts:1-80`、`tool-name.ts` 导出名、`index.ts` 导出名、`connection-pool.ts` 两枚注释行；`skills` 读了 `types.ts:1-79` 与 `registry.ts` 若干行 |

**B. 越界读物（在 `packages/agent-modules/` 之外，为了回答派单点名的问题才读的；每条都在正文里标了"越界"）**：
`packages/agent-extension/src/runaway-guard.ts`（262 行，全读）·
`packages/agent-extension/src/index.ts`（导出名）·
`packages/config/src/runaway-guard-config.ts`（22 行，全读）·
`packages/local-runtime-v2/src/service/turn-system/runaway-guard/{extension.ts,tool-policy.ts,observation.ts}`（全读）·
`packages/local-runtime-v2/src/service/turn-system/production-composition.ts`（**只 grep 到文件名，未读**）·
`packages/config/src/config.ts`（只读 `:820-860` 的 DailyDigest/Memory）·
`packages/agent-tools/src/shared/cron-schedule-input.ts`（读 `:1-60`）·
`packages/agent-tools/src/desktop/builtin-defs.ts`（**只 grep 了 cron 相关行**，`:714-996` 抽检）·
`packages/local-runtime/src/background-task/{startup-recovery.ts(读 182 行中 1-182 的关键段), schema.ts(全读 46 行)}`·
`packages/agent-core/src/pi-turn-runner/{events.ts,types.ts}`（**只看 grep 命中行**）。
未读但存在、且可能改变本轮结论的：`packages/agent-tools/src/desktop/local-mavis-cron-adapter.ts`（cron 工具适配层，
`deliveredAt` 出现在 `:70`、`:551`，**没读**）、`packages/local-runtime/src/background-task/{service.ts 478, delivery.ts 364,
terminal.ts 186, checkpoint-snapshot.ts 191, bash-runner.ts 446, runner.ts 313, conversation-delivery.ts 111}`（**全未读**）。

**C. 我按纪律没碰的**：`frontend/**`、`design/**`（派单禁读）；参照仓内任何写操作；`docs/reports/pending-and-issues.md`；
别人的证据件。**未做任何 commit**。

**D. 本轮未能定量的项（正文里已标的汇总，以及已回收的）**：
1. ~~同名遮蔽/败者名单的对应物~~ ⇒ **已量**，但口径受限：`grep -rniE "loser|shadow|遮蔽" --include=*.go internal/tools internal/config`＝**0 命中**，
   **全仓未搜**，所以这条只能写成"限定范围内为 0"（见 §7.1）；
2. ~~"删除后名字不回收"的名册持久层~~ ⇒ **已量**：`nameRegistry|NameRegistry|tool_name.*unique` 全仓 **0 命中**，
   `internal/memory` 里的 4 处 `roster` 全是协程名册（见 §7.3）；
3. "我们的提醒文案会不会被模型写进记忆"这一维**未量**（我只量了没有频控表，没量防污染条款；
   他们那句 anti-contamination 后缀在 `runaway-guard/src/reminder.ts:59-60`，我们有没有等价物本轮**没查**）；
4. 我们 D42 可靠性登记表与 §3.2 `lost`／`deliveredAt` 的逐条对应关系（**未做**，只写了"对应缺口未逐条比"）；
5. §4.3 那句"每 Turn 结构化摘要我们没有"⇒ 我量的是 `policy_version`/`strategy_version` 两个具体词为 0，
   **没有做"全仓搜同类摘要形状"的反向取证**，所以严格讲这条应是"两枚具体版本号缺失"，不是"完全没有摘要类读数"。
   **如果他要把这条写进差集，请按第 5 点的弱版本写。**

---

# ② 我认为**编排者会误读**的地方（他被推翻／将被推翻的结论逐条列）

1. **会被推翻（派单内部前提，不是他的账）**：派单写"我们有审批门但**没有"它自己转圈"这一维**"。
   ⇒ **错**。我们有 C22 LoopGuard，`internal/agent/guard.go:14-19` 注释逐字写着"梯度刹车／三级 `[3,5,8]`／
   顶级进 `Stuck` 并显式告知用户／禁止静默停"。
   **真实的差不是"有/无"，是"粒度"**：我们比的是**整轮签名**（`guard.go:234-243`），他们比的是**单条调用**＋错误家族＋
   ABAB＋进度证据（§1.1、§1.4）。如果按"我们没有"去开票，票会写成"新增防跑飞"，而正确票是"把已有那枚降粒度＋补信号"，
   **两者判据完全不同**（后者要动 `ObserveTurn` 的签名语义，会牵 D43 行 #19）。
2. **不会被推翻（他点名要我检的那条）**：`session-report` **不是**"今天它替我做了哪些决定"的日报，
   它是**诊断包／反馈工件收集器**（`session-report/package.json:5` 逐字 `Session diagnostics and feedback artifact
   collection.`；消费者在 `packages/tui/src/runtime/feedback/`）。
   ⇒ **`A409` 第九节与 `docs/reports/missing-features-2026-09-28-v2.md` 第九节那一行维持原判，不必改口。**
   ⚠ 但**极易被误读成推翻**，因为仓里真有一枚 `DailyDigestConfig`（`packages/config/src/config.ts:824-837`）。
   它**不是日报**：它写的是 `<dataDir>/agents/<name>/daily/` 下的**记忆文件**，**`enabled` 默认 false**（`:836`），
   且注释明确 `This only gates digest generation and daily TTL archival. It does NOT disable the scheduler itself`
   （`:829-832`）；更关键的是**被引用的实现文件 `daily-digest.ts` 在这份公开克隆里根本不存在**
   〔现读：`find . -name "*daily*" -type f` 只有 `packages/shared/src/daily-signin.ts`（签到活动）与两枚 tui 测试〕。
   ⇒ 如果他看到别人转述"MiniMax 有 daily digest"就以为被推翻，那是**把记忆整理读成了决策汇报**。
3. **会被部分推翻（他给 owner 的差集清单里"MCP"那一行）**：他在 `AGENTS.md`/`PLAN.md:1533` 把 MCP 记为
   **RESERVED、无实现计划**，理由是与 host bridge 架构冲突。
   ⇒ 事实层面要补一句：**MiniMax 家有一整份可读的 MCP 运行时（13 枚 `.ts`，含连接池 911 行与持久名字治理 258 行）**，
   所以"没有参考实现"这个隐含理由**不再成立**；但他们**并没有推翻 D13 的架构冲突论**——
   那是我们自己的取舍，属人工批准级。⇒ 差集清单里这一行应从"别人有我们没有"改写成
   **"别人有一份现成实现；我们不接是主动取舍，理由仍成立，但今后不能再写'无处可抄'"**。
4. **会被推翻（他可能据此判"三档制不够用"）**：他们的授权**不是三档**，是
   **6 枚 `PermissionMode` × 3 枚 behavior × 3 层规则来源 × 5 档授权宽度 × 4 轴 profile**（§2.1）。
   ⇒ 很容易读成"我们三档太糙"。**不能这么读**：他们 `mode ?? 'default'` 之后仍有
   `no-checker ⇒ allow`（`permission-core.ts:279-287`）这条**默认放行**兜底，
   而我们的 R9 是 `无判定器 ⇒ fail-closed 升 L2`（`SPEC-06:42`）。**在"最坏情况谁更保守"这一问上我们更严**。
   真正的差距只在**表达力**（能否描述"允许这个命令的前两个 argv"），不在**安全底线**。
5. **容易被他当成"他们已经解决了"的一条**：`cron` 的**错过档期不补跑**（§3.1 第 3 行，`lastRun` 重启即 `null`、
   无 overdue/misfire 逻辑）。如果他在差集里写"MiniMax 有完善的定时续跑"，**那是过度概括**——
   他们做的是**登记恢复**与**历史保留清扫**，不是**补偿执行**。
   我们 D12 若要落地"OS 定时"，缺的那一块**两边都没有**，只是我们连登记恢复也还没有。
6. **容易误读为"他们能取消后台任务所以我们应该照抄"的一条**：他们的 `stop()` **在 runner 停止失败时仍写 `canceled`**
   状态（只是事件里带 `stopFailed:true`，`manager.ts:139-147`）。
   ⇒ 若我们照抄状态机就会造出一枚**"账上已取消、实际进程还在"**的状态。他们靠 payload 里的 `stopFailed` 说实话；
   **抄状态不抄那枚旗标＝制造假终态**。这一点必须写进票的判据。
7. **我这条腿自己的结论也可能被误读**：我把 `session-report` 判成"非日报"，但**这不等于"五家里一定没有人做日报"**——
   我只复查了 MiniMax 这一家；腿 1（openchamber）与腿 2（DSH 的 52 枚 `ui-*`）如果翻出日报，**以他们的现量为准，
   不要用我这句"维持原判"去否掉他们的读数**。
8. **一处措辞陷阱**：本文件里所有"[runaway guard] ..."、"Auto classifier"之类引文都是**给模型看的话术**，
   不是给用户看的 UI 文案。派单问"给用户留什么话"，我在 §1.2 的答是**"给用户：没有直接文案"**——
   如果他据此写"我们要有 UI 文案"，那是**把"给模型的 steer"读成"给用户的提示"**（我们的 C22 反倒两样都有：
   `guard.go:192-194` 那句提醒是注入对话的，`loop.go:470-483` 的 `EvStuck` 才是给人看的）。





