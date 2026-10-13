# 票 145 · B 节 —— 十四态"缺哪个字段"现读表刷新（只读腿 `145-a2`）

- 锚 = dev `8e4f6f83`。表体尺：`sed -n '3479,3494p' docs/PLAN.md` ⇒ 表头 `:3479` + 分隔 `:3480` + **14 行数据 `:3481-3494`**。
  ⚠ 票面/普查引的 `PLAN.md:3473-3488` 已漂（现表体 = `:3479-3494`），原句不抹、就地更正。
- 14 行（按 `:3481-3494` 顺序）：思考中 / SSE 流式 / 推理过程 / 工具调用 / 工具调用(展开) / 审批等待(L2) / L2 确认卡 / L1 阻止窗口 / L2 原生降级卡 / 错误 / Stuck / 成本 / 已取消 / 注入检出。
- 「旧判」列 = `145-c2`（`probes/145/c2/census.md`，10-03）与本票"AC#1 交件后六处更正"①（缺 8 行 = 1/3/4/5/10/11/12/13）的口径。本腿**复量后不改旧判**，只刷新"今天快照上有没有那一维"。

---

## B1 逐态现状表（每条 = 态 / 需要的字段 / 今天快照有无该载体 / 现读判）

| # | 行 | 态 | 需要的字段（`PLAN.md` 视觉列点名） | 今天载体 | 现读判 |
|---|---|---|---|---|---|
| 1 | `:3481` | 思考中（等 LLM 首 token） | 已耗时计时（`1.2s` 实时走）+ phase 标记 | **无**（无 `phase`/`elapsed` 字段） | 无输入 |
| 2 | `:3482` | SSE 流式输出 | 逐段文字 + 实时 token 计数 | `results[].text`（活）有；**token 计数无** | 部分（半 fed） |
| 3 | `:3483` | 推理过程（reasoning delta） | 推理文本 + 推理耗时 | **无** | 无输入 |
| 4 | `:3484` | 工具调用（chip） | 工具名 + 参数摘要 + 16px 类别图标 + 4 值状态(运行/成功/失败/被拒) + 风险级徽标 | **无**（无 `tools[]`） | 无输入 |
| 5 | `:3485` | 工具调用（展开） | 完整参数 JSON + 结果摘要 + 耗时 + `correlationId` | **无** | 无输入 |
| 6 | `:3486` | 审批等待（L2） | chip 变 warn 描边 + **队列深度角标** | `pending[]`（卡列表）在；**per-call 队列深度/等待标记无** | 矛盾（旧判） |
| 7 | `:3487` | L2 确认卡（面板内） | 工具名 + `L2` 徽标 + 完整参数 + C19 命中规则 | `pending[] ApprovalCardView`（10 键，`approval.go:39-59`）**齐** | **已有**（旧判） |
| 8 | `:3488` | L1 阻止窗口 | 工具名 + 参数摘要 + 环形倒计时(`remainingMs`) | `pending[]` 无 L1；`PumpSources.L1Windows` **全仓零装配点** => 窗口维度黑 | 矛盾/半输入（旧判） |
| 9 | `:3489` | L2 原生降级卡（面板不可用） | 原生侧小卡 | **原生侧，面板契约够不着**（不在快照射程） | 射程外 |
| 10 | `:3490` | 错误 | 人话文案(`humanText`) + `error_class` + 技术细节 + 动作 | **无** | 无输入 |
| 11 | `:3491` | Stuck（梯度刹车） | 重复了什么 + 三按钮 | **无** | 无输入 |
| 12 | `:3492` | 成本（C23） | `IN/OUT/CACHED` token + 金额 | **无**（无 `cost`/`usage`） | 无输入 |
| 13 | `:3493` | 已取消 / 被打断 | 取消分隔线（取消原因） | **无** | 无输入 |
| 14 | `:3494` | 注入检出（C25/R4） | 危险条 + 来源 URL + 命中片段 | 来源链**今天通**：`risk/rules_gateway.go:112` -> `panel/approval.go:83` -> 卡 `rulesHit`/`reason`；**命中片段 `fragment` 无** | 部分（半 fed） |

**计数刷新（复量 `145-c2` 口径）**：完全无输入 = **8 行**（1/3/4/5/10/11/12/13）；加行 8 的半枚 = 9；再加行 9（射程外）= 10。**数字不变**（= `A244①`/票面① 的 8）。

## B2 为什么 299 面板宿主落地**没有**改变这张表

299 落地新增的载体是 `Tasks`（票 197）与 `Instructions`（票 200）两段，以及更早的 `Git`/`CurrentModel`/`Credential` 三枚 composer 子键 —— 本腿逐枚复认它们**全喂"设置 / 名册 / 工作区"侧，无一枚喂上面那 14 行工作态**：

- `Tasks`（`TaskRosterSection`，`subagent_roster_197.go:150`）：喂"任务监控"屏（票 182/197），与 14 行表**不同物**。
- `Instructions`（`InstructionsSection`，`instructions_200.go:71`）：喂上下文/指令加载态，与 14 行表**不同物**。
- `Git`/`CurrentModel`/`Credential`：工作区/模型/凭据侧，与 14 行表**不同物**。

=> 本票标题级前提"缺 8~10"的**方向仍成立**（`145-c2` 的再框定本腿复量一致）。299 改变的是"页面能不能画得出"（宿主在），**没有**改变"快照里有没有那 14 行要的数据"。

## B3 今天真卡点两类（复量 `145-c2` 口径，逐枚现读）

- **跨-seam（有活源、泵读不到）**：`cost`/`usage`/`reasoning`/`stuck` —— 活在 `internal/agent/loop.go` 的 run 局部（`145-c2` 报 `:367`；本腿未复跑该行，按〔腿报〕标），`Snapshot` 在 execute() 外也驱动 => run 局部不可达。
- **纯无源**（今天算不出来）：计时 elapsed / 人话文案 `humanText` / 命中片段 `fragment` / 实时剩余 `remainingMs` / `phase` 映射表。
- **恒空两枚**（载体在、泵不填）：`ComposerState.Attachments`（`pump.go:269` 传 nil 归一空，`composer.go:274-276`）、`AttachmentReason`（`pump.go` 从不赋值）。

## B4 现读"宁缺毋造"侧（新落地维逐枚带未知守卫）

- `currentModel` + `modelKnown`（`pump.go:280-286`）：无读口 -> `modelKnown=false`。
- `credentialState` + `credentialKnown`（`pump.go:287-297`）：无读口 -> `CredentialKnown=false`。
- `git`（`pump.go:275-279`）：无读口 -> `GitViewNotProbed()`。
- `tasks` / `instructions`：**指针 + omitempty**，无读口则整键缺席（`pump.go:304-328`）。
- 残余两枚无守卫（`145-c2` 报，本腿静态复认）：`AttachmentReason` 无 known 位；`StreamLog.TruncationFor` 未知键塌零值（`pump.go:692-712`）。
