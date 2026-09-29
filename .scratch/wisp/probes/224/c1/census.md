# 224-c1 普查：「会话」这一档在本仓的全部载体、判定链、三件套落点、三枚钉的射程

- 锚点 `git rev-parse --short HEAD` = **`6233dedc`**（分支 dev，现跑；本件的读数是**对着这颗 sha 的工作树**取的）
- ⚠ **锚点与"交回时 HEAD"双记**：`11:50` 现跑 `git rev-parse --short HEAD` 已前移为 **`e2594664`**（220-c1 证据件入库，采集期别人的写腿把链往前推了——与本仓 220-c1 那轮同形）。
  ⇒ 我随后**把本件引用的全部 `cmd/wisp/**` 与 `internal/agent/approval/**` 行号在 `e2594664` 上重跑一遍核对**：
  `run.go` 的 14 处（`:154/:212/:220/:229/:321/:381/:549/:676/:721/:751/:756/:758/:919`）、`main.go` 的 5 处（`:83/:85/:109/:145/:153`）、
  `approval_reply.go` 的 5 处（`:208/:235/:276/:389/:433`，全文仍 539 行）、`gate.go` 的 4 处（`:387/:596/:626/:638`）、`queue.go:343` —— **全部一致，无一漂移**。
  ⚠ 一处自查纠正：`nativeAPI.Allow` 的真位置是 **`gate.go:596`**（我一度误写成 `approval.go:596`），`Queue.allow` 是 **`:343`**（一度误写 `:342`），已在 C1 就地改正。
- 取证时间（现跑 `date "+%Y-%m-%d %H:%M %z"`）= **`2026-09-29 11:33 +0800`**
- 落盘完成（同一枚尺再跑一次，非推算）= **`2026-09-29 11:49 +0800`**
- ⚠ 工作树非干净：`git status --short` 现跑显示 `cmd/wisp/approval_reply.go` 与 `cmd/wisp/run.go` 处于 **M（未提交）**状态
  ⇒ 凡引 `cmd/wisp/**` 与 `internal/agent/approval/**` 的行号，**交回即可能已漂**（各节已就地标注）。
- 本轮**只读**：未跑 `go test` / `go build` / `go vet`，未改任何源码／工单／台账。
- `frontend/**` 与 `design/**` **零读零转述**（本文件不含其中任何内容）。
- 所有行号均为本锚点上现跑 `grep -n` / `sed -n` 所得；**未从 `.scratch/wisp/probes/**` 任何归档快照取行号**。

## 0. 本轮跑过的尺（可原样重跑，均在仓库根执行）

```bash
git rev-parse --short HEAD ; date "+%Y-%m-%d %H:%M %z"
grep -rn 'SessionID\|session_id' --include='*.go' internal cmd tools | grep -v '_test\.go'
grep -rni 'session' --include='*.go' internal cmd tools | grep -v '_test\.go'
grep -rn 'approval_grant' --include='*.go' .
grep -rn 'InsertGrant\|ListGrantsBySession\|RevokeGrant\|ListAllGrants\|DeleteGrant\|DecisionAllowGrant' --include='*.go' internal cmd tools
sed -n '85,105p' internal/memory/schema.go ; sed -n '85,115p' internal/memory/models.go
grep -n 'func \|resolveMode' internal/risk/mode.go ; sed -n '108,220p' internal/risk/mode.go
grep -rn '\.Screen(\|Silenced\|LoosestOf' --include='*.go' internal cmd tools | grep -v '_test\.go'
sed -n '30,200p' internal/tools/bridge.go ; sed -n '250,300p' internal/tools/bridge.go ; sed -n '355,420p' internal/tools/bridge.go ; sed -n '690,790p' internal/tools/bridge.go
grep -n 'func (b \*Bridge) Execute\|sil := mode.Screen\|switch {\|case sil.Silenced\|b\.route(' internal/tools/bridge.go
sed -n '1,40p' internal/perm/ticket90_persist_test.go ; sed -n '138,292p' internal/perm/ticket90_persist_test.go
sed -n '1,60p' internal/perm/store.go ; sed -n '400,470p' cmd/wisp/run_mode101_test.go
grep -rn 'func Test' --include='*_test.go' internal cmd tools | grep -iE 'session|grant'
sed -n '185,242p;670,715p;740,800p;1229,1275p;1330,1460p;1715,1760p;1870,1930p;1959,2000p' internal/panel/l2_grant_boundary_test.go
grep -n 'func knownComposerMethod' internal/panel/bridge.go ; sed -n '95,135p' internal/panel/bridge.go
sed -n '180,200p' internal/config/schema.go ; grep -rn '\.Session\b\|SessionSection' --include='*.go' internal cmd tools | grep -v '_test\.go'
ls internal/session/ ; cat internal/session/doc.go
sed -n '1536,1545p;1636,1648p;2140,2160p;2180,2195p;2700,2710p' docs/PLAN.md
sed -n '156,200p' docs/specs/SPEC-08-ui-ball-panel.md ; sed -n '75,105p;130,160p' docs/specs/SPEC-02-data-storage.md
grep -rn 'C31\|SessionScope' --include='*.go' internal cmd tools | grep -v '_test\.go'
grep -n 'defer rt.close()\|return rt.execute(task)\|func cmdRun\|runTextTask(' cmd/wisp/*.go ; sed -n '80,90p;145,160p' cmd/wisp/main.go
grep -n 'func \|Ticket' cmd/wisp/approval_reply.go
grep -rn 'GrantID\|grant_id' --include='*.go' internal cmd | grep -v '_test\.go'
sed -n '25,40p' internal/agent/journal.go ; sed -n '215,235p' internal/memory/retention.go
```

## A. 「会话」在本仓已有的全部候选载体

> 三态标记：**〔已证（现读）〕** = 这行代码在这颗 sha 上确实存在且我逐字看过；**〔建了但没接〕** = 结构／表／包边界在，生产没人写没人读；**〔仅文档写了、代码没有〕**。
> 「有没有明确的结束点」是本票的命门：**没有结束点＝永久通行证**；有结束点但没人执行＝同样没有结束点。

### A1 `wisp run` 的进程生命周期 —— 有明确结束点，但**它不铸造任何身份**

| 事实 | file:line | 三态 |
|---|---|---|
| 入口派发：`case "run":` → `os.Exit(cmdRun(args[1:]))` | `cmd/wisp/main.go:83`、`:85` | 〔已证（现读）〕 |
| `func cmdRun(args []string) int` → `return runTextTask(runSpec{…})` | `cmd/wisp/main.go:145`、`:153` | 〔已证（现读）〕 |
| 一次进程 = 一枚任务文本：`task := strings.TrimSpace(firstArg(s.argv))`，空则 `return 2` | `cmd/wisp/run.go:167`、`:171` | 〔已证（现读）〕 |
| 进程开始（装配）：`func runTextTask` 内 `defer rt.close()`（第 212 行注册），装配后 `return rt.execute(task)` | `cmd/wisp/run.go:154`、`:212`、`:220` | 〔已证（现读）〕 |
| **让它退出的东西**：`execute()` 返回一个 int → `runTextTask` return → `main` 的 `os.Exit` 结束进程 | `cmd/wisp/run.go:751`、`:220`、`cmd/wisp/main.go:85` | 〔已证（现读）〕 |
| 关停顺序：`func (rt *agentRuntime) close()` 只做两件事——`rt.replyRoot.Cancel()` 与 `rt.store.Close()` | `cmd/wisp/run.go:676-683` | 〔已证（现读）〕 |
| ⚠ 关停**不含任何「撤销本会话授权」的一步**：`close()` 里没有 grant 相关调用（全函数只有上面两行） | `cmd/wisp/run.go:676-683` | 〔已证（现读）〕 |

**结论 A1**：进程生命周期是**今天唯一真存在、且有明确结束点**的「会话」候选（进程退出即结束，OS 保证）。
但它**从头到尾没有铸造任何 id**：全仓非测试代码里 `SessionID` 的命中只有 `internal/memory` 自己（见 A3），
`run.go`／`main.go` 里没有任何一行为这次进程生成 session 身份。⇒ 拿它当会话＝有结束点、**没有身份**，
授权行写进去之后没有任何东西能把它和下一次进程区分开。

### A2 `internal/agent`：没有 session/conversation/thread 结构，只有 **task id**

| 事实 | file:line | 三态 |
|---|---|---|
| `type Loop struct` 的字段：`opt / b / guard / asm / sp / comp / reg / provider / mu / history / steer / current` —— **没有 session 字段，没有会话身份字段** | `internal/agent/loop.go`（`type Loop struct` 起，字段逐枚现读） | 〔已证（现读）〕 |
| Loop 唯一的「会话语义」是注释级的：`// Reset drops the conversation history (session end).` → 只清 `history`/`steer` | `internal/agent/loop.go:264-270` | 〔已证（现读）〕，⚠ **不产生身份、不通知任何授权面** |
| 身份铸造只有 task id：`// newTaskID returns a UUIDv4-shaped task id` → `func newTaskID() string`，调用者只有两枚：`:322`（`RunAsync`）、`:333`（`Run`） | `internal/agent/loop.go:1121-1123`、`:322`、`:333` | 〔已证（现读）〕 |
| task 作用域已有完整生命周期（**本仓今天最接近「有界会话」的机器**）：`scopes map[string]*risk.Scope`、`func (b *Bridge) OpenTask(taskID string)`、`func (b *Bridge) CloseTask(taskID string)` | `internal/tools/bridge.go:133`（字段）、`:699`、`:749` | 〔已证（现读）〕 |
| 它的关闭**确实挂在生产边界上**：`AdmitTask: rt.admitTask`（文本循环 owns 它）；`func (rt *agentRuntime) admitTask(taskID string) func()` | `cmd/wisp/run.go:733`、`:721` | 〔已证（现读）〕 |
| ⚠ `CloseTask` 的注释自陈**仍有没 owner 的腿**：「Today no production code dispatches on the bridge except loop.go」、「a host inventing its own id … nothing in this tree closes it」、该问题＝**`Q-56`** | `internal/tools/bridge.go:712-748` | 〔已证（现读）〕，`Q-56` 属待人定案 |
| 决策词表里「会话授权」那一档**只有常量、零写者**：`DecisionAllowGrant = "allow_session_grant"` | `internal/agent/journal.go:32` | 〔建了但没接〕（非测试命中＝**这一行自己**） |

**结论 A2**：`internal/agent` 里**没有**任何「一次会话」的结构体。可复用的、有界的身份只有 **task id**。
task id 有结束点（`CloseTask`），但它是 `tool_call.task_id`／`task_log.id` 的口径，
而 SPEC-02／`PLAN.md` 把授权绑的是 `session_id` 不是 `task_id` ⇒ **拿 task id 冒充 session id 是一枚契约面变更（要批准）**。

### A3 `internal/memory`：`SessionID` 的引用者**只有 DAO 自己＋schema＋retention**

`grep -rn 'SessionID|session_id' --include='*.go' internal cmd tools | grep -v _test` 的**全部**命中（现跑，逐枚）：

| file:line | 是什么 |
|---|---|
| `internal/memory/dao_misc.go:21`、`:22` | 插入前必填校验 `g.Tool == "" \|\| g.Pattern == "" \|\| g.SessionID == ""` → 错误文本 `memory: grant tool/pattern/session_id are required` |
| `internal/memory/dao_misc.go:33`、`:35` | `INSERT INTO approval_grant(scope, tool, pattern, session_id, created_at, expires_at, revoked_at)`，绑 `g.SessionID` |
| `internal/memory/dao_misc.go:80`、`:81` | `SELECT … FROM approval_grant WHERE session_id=? ORDER BY created_at DESC, id DESC`（`ListGrantsBySession`） |
| `internal/memory/dao_misc.go:92` | `ListGrants` 的 SELECT（**不带 session 过滤**，全表读，审计用） |
| `internal/memory/dao_misc.go:120` | `rows.Scan(&g.ID, &g.Scope, &g.Tool, &g.Pattern, &g.SessionID, …)` |
| `internal/memory/models.go:98` | `SessionID string`（`ApprovalGrant` 字段） |
| `internal/memory/schema.go:94` | DDL 列 `session_id TEXT NOT NULL` |
| `internal/memory/schema.go:99` | `CREATE INDEX idx_grant_session ON approval_grant(session_id)` |

**除 `dao_misc.go` 之外引用 `SessionID` 的只有 `models.go` 与 `schema.go`（声明与 DDL），没有任何第三个模块读它。**
另两枚顺带事实：
- `internal/memory/retention.go:219 func deleteExpiredGrants` → `retention.go:225 DELETE FROM approval_grant WHERE MAX(expires_at, COALESCE(revoked_at, 0)) < ?`：
  这是**唯一一段会在生产碰 `approval_grant` 的代码**，但它**只删不读不写**，判据是**墙上时间**，不是「会话是否还活着」。〔已证（现读）〕
- `tool_call.grant_id` **有管道、无来源**：`schema.go:83`（列）、`dao_toolcall.go:25`、`:31`、`:52`、`:160`、`:176`、`models.go:88`（`GrantID *int64`）。
  全仓非测试**没有任何一处把 `GrantID` 赋成非 nil** ⇒「哪一行 grant 放掉了这一次调用」今天写不出来。〔建了但没接〕

### A4 `approval_grant` 建表 DDL 与逐枚列

位置：**`internal/memory/schema.go:89-99`**（`schema.go:88` 逐字 `-- D45 作用域会话授权`）。列逐枚（现读 DDL）：

| 列 | DDL 原文（`schema.go:90-97`） | 生产写手 | 生产读者 |
|---|---|---|---|
| `id` | `INTEGER PRIMARY KEY AUTOINCREMENT` | 由 `InsertGrant` 返回 | — |
| `scope` | `TEXT NOT NULL, -- 'session'` | **0** | **0** |
| `tool` | `TEXT NOT NULL` | **0** | **0** |
| `pattern` | `TEXT NOT NULL, -- 路径模式 / 目标进程名（input.type）` | **0** | **0** |
| `session_id` | `TEXT NOT NULL` | **0** | **0** |
| `created_at` | `INTEGER NOT NULL` | **0** | **0** |
| `expires_at` | `INTEGER NOT NULL, -- 会话结束时间（失效后仅作审计）` | **0** | **0**（只有 retention 删） |
| `revoked_at` | `INTEGER`（可空） | **0**（`RevokeGrant` 存在，`dao_misc.go:49`） | **0** |
| 索引 | `schema.go:99 CREATE INDEX idx_grant_session ON approval_grant(session_id)` | — | — |

⚠ **DDL 的注释自己就把 `expires_at` 定义成「会话结束时间」**（`schema.go:96`）——这是仓里唯一一处把「结束」落到可执行字段的**文字**；代码侧没人执行它。〔仅文档写了（在 DDL 注释里）、代码没有〕

- `internal/memory/models.go:104` 逐字：`// GrantScopeSession is the only scope value (SPEC-02 §3 approval_grant.scope).`；`:105 const GrantScopeSession = "session"`。
- `internal/memory/models.go:91-92`：`// ApprovalGrant is one D45 scoped session grant. Expired/revoked rows stay for audit for GrantAuditTTL (30 days)`。
- `dao_misc.go:18-19` 拒绝非 `session` 的 scope：`memory: invalid grant scope %q (want session)` ⇒ **今天连「另一档」都写不进去**。
- DAO 方法清单（`internal/memory/dao_misc.go`）：`InsertGrant:17`、`RevokeGrant:49`、`ListGrantsBySession:78`、`ListGrants:91`、`DeleteGrant:102`。
  非测试调用者（`grep -rn` 现跑）：**全部为零**。测试调用者三处：`internal/memory/dao_test.go:383`、`internal/perm/ticket90_persist_test.go:230`、`cmd/wisp/run_mode101_test.go:411`。〔建了但没接〕

### A5 `internal/panel` / `internal/ball`：有没有「一次交互会话」的边界

| 侧 | 事实 | file:line | 三态 |
|---|---|---|---|
| panel | **没有会话身份**。`internal/panel` 非测试里所有 `session` 命中都属 C25 污染：`SessionOverrideBlocked bool \`json:"sessionOverrideBlocked"\`` | `internal/panel/approval.go:54`、`:85`；`internal/panel/pump.go:80`、`:100` | 〔已证（现读）〕⚠ **别把这枚当会话位读** |
| panel | 唯一的「会话」字样是给用户看的文案：`…请重启会话` | `internal/panel/workspace.go:121` | 〔已证（现读）〕，非结构 |
| panel | Go 今天真答复的入向方法只有四枚（`MethodModeRequest / MethodWorkspaceRequest / MethodAttachmentAdd / MethodMessageSend`），`knownComposerMethod` 是唯一的门 | `internal/panel/bridge.go:104-110`；常量 `:43` 一带 | 〔已证（现读）〕，**没有 grants 相关方法** |
| ball | **有「一次交互会话」的起止**，但它是**液态视觉**的起止、不携身份：`func (m *liquidMotion) beginSession()` / `endSession()`；边沿由 `func (b *Ball) motionStateChanged(from, to statemachine.State)` 按 `liquidDriven` 判定 | `internal/ball/liquid.go:88`、`:99`；`internal/ball/liquid_windows.go:153-167`（`:157`、`:163`、`:165`） | 〔已证（现读）〕 |
| ball | 会话族还有一枚热键归还：`func (b *Ball) ReleaseEscAfterSession()`；`// MUST be handed back at session end` | `internal/ball/ball_windows.go:866-868`；`internal/ball/hotkey_windows.go:7`、`:448`、`:465` | 〔已证（现读）〕 |
| ball | ⚠ 会话边界的**真判据在状态机**，而状态机是**冻结的 D43 转移表**：起＝`Sleeping --EvSummon--> Listening`、`Warm --EvSummon--> Listening`；止＝`Warm --EvWarmIdle--> Settling`、`Settling --EvSettleExpired--> Sleeping` | `internal/statemachine/table.go:51`、`:215`、`:227`、`:231`、`:236` | 〔已证（现读）〕 |

**结论 A5**：**球那一侧真有边界，但没有铸造点**。`PLAN.md:1279`／`:1381` 指定的铸造者 **C31 `SessionScope`** 在这颗 sha 上是**一枚空包**：

- `ls internal/session/` → **只有 `doc.go`**（没有其它 `.go`）。
- `internal/session/doc.go:1-3` 逐字：`// Package session owns SessionScope (C31) (SPEC-01 §3): the wake -> explicit` / `// end / 90s-idle lifecycle, Warm/Conversation timing, model reference` / `// counting, and D45 authorization registration.`（`:6-8` 另逐字写着 `session end MUST fully dispose, including debug.FreeOSMemory (C11 step 4)`）
- `internal/session/doc.go:10` 逐字：`//   - per-session grants ledger (D45-2)`
- `internal/session/doc.go:12-14` 逐字（非职责）：`// Non-responsibilities:` / `//   - no state transitions (statemachine), no model loading (speech/models),` / `//   - no audio device ownership (audio)`
- `internal/session/doc.go:16-17` 逐字：`// DEFERRED(SessionScope): implemented by ticket 28. This ticket only freezes` / `// the package boundary.`，末行 `package session`
- `grep -rn 'SessionScope|C31' --include='*.go' internal cmd tools | grep -v _test` 只有四命中：上面 `doc.go:1`/`:16`，加 `internal/plugin/disposal.go:33`（注释 `session scope (C31) > task scope > tool scope > plugin scope`）、`internal/statemachine/table.go:16`（注释：转移钩子 `executing them (SessionScope creation, model loading, panel pushes)` 不归状态机）。
- ⚠ **没有任何代码执行「建 SessionScope」／「Dispose SessionScope」**：`grep -rn 'Settling|Dispose|FreeOSMemory'` 非测试命中里 `internal/ball/*` 全是视觉（fade/breath 计时），`internal/observe/sampler.go:408` 是量测闸门。〔建了但没接（包边界）／仅文档写了（生命周期与授权登记）〕

### A6 `internal/config`：有没有 session 级生效的键

| 事实 | file:line | 三态 |
|---|---|---|
| `[session]` 这一节**存在**：`Session SessionSection \`toml:"session"\`` | `internal/config/schema.go:113` | 〔已证（现读）〕 |
| 节内三枚键，注释逐字 `// SessionSection is [session]; hot-tier.`：`warm_timeout_sec`(90)／`settling_sec`(3)／`conversation_idle_sec`(30) | `internal/config/schema.go:186-191` | 〔已证（现读）〕 |
| 热加载差异表认识这一节：`{"session", &cur.Session, &fresh.Session, func() { cur.Session = fresh.Session }}` | `internal/config/manager.go:205` | 〔已证（现读）〕 |
| ⛔ **三枚键在生产路径上零读者**：`grep -rn '\.Session\b\|ConversationIdleSec\|WarmTimeoutSec\|SettlingSec' --include='*.go' internal cmd tools \| grep -v _test` 的命中**只有上面这三行**（声明＋热加载表），没有第四处 | 同上 | 〔建了但没接〕——**键能写、没人读** |
| ⚠ 它们**不在**票 83 那套「说谎的键必须响亮失败」守卫表里（`internal/config/unwired.go:60 var unwiredKeys`、`:99-101 validateUnwired`；`unwired_test.go:172` 的 `[session]` 只是给守卫测射程的样例），所以今天**静默无效不会报错** | `internal/config/unwired.go:43-60` | 〔已证（现读）〕 |
| 真正管事的「会话」约束写在注释里而不是键里：`NOT persisted anywhere: a D45 session grant. GrantScopeSession stays` / `session-scoped and dies with the process (PLAN.md:1640, ticket 49)` | `internal/config/permmode.go:19-20`；同一句另见 `internal/config/schema.go:465-467`（`nothing else may borrow that channel`） | 〔已证（现读）〕⚠ **只有注释在管，没有代码在管** |

**结论 A6**：`[session]` 是**音频/交互计时**的档位（`warm_timeout_sec` 正是 `PLAN.md:2187` 那句会话定义点名的可配项），
**不是**授权 scope 的档位；而且今天**一枚读者都没有**。⇒ 会话级生效的键：**不存在**。

### A7 冻结文字里凡是写了 `session` 这一档的地方（**逐字抄，只引不改**）

**`docs/PLAN.md`**（现跑 `sed -n`）：

- `docs/PLAN.md:1279`：`| **\`session\`** | **第四轮新增**：**C31 \`SessionScope\`** —— 会话生命周期、\`Warm\`/\`Conversation\` 计时、模型引用计数、面板保活、D45 会话授权登记 | 具体模型推理（归 \`speech\`） |` —— **「会话」作为一枚职责层的出处**，它把「D45 会话授权登记」分给 C31。
- `docs/PLAN.md:1381`：`| **C31** | **\`SessionScope\`** | 会话定义（唤起 → 显式结束 / 90s 空闲）；持有 Warm 模型引用计数、面板保活、D45 会话授权、\`Settling\`/\`Warm\` 计时；**会话结束必须完整 Dispose（C11）**。BLOCKER B4 的实现契约 | D32/D43 |`
- `docs/PLAN.md:1535`：`**\`session\` 档已由 D45 实现并从本条移出**；只剩跨会话持久档仍无计划。理由：自用期会话级已够…`
  ⚠ **这句「已由 D45 实现」就是本票顶部指的那枚落差**：冻结文字说「已实现」，A1–A6 证的是〔建了但没接〕。
- `docs/PLAN.md:1536`：`| RESERVED | Codex 的 \`available_decisions\`（服务端告诉客户端渲染哪些按钮） | D31：当前确认动作固定（允许一次/拒绝/**本会话内允许**，第三项由 D45 加入），无需动态 | — **无实现计划** | C18 | 确认卡按钮是硬编码的三种 |`
- `docs/PLAN.md:1642`：`会话授权**不得**覆盖 C25 污染升级；会话结束后授权**必须**失效；`
- `docs/PLAN.md:2147-2149`（D45 第 2 条本体）：`「**本会话内允许 \`<工具>\` 于 \`<路径模式>\`**」。授权绑定 **(工具, 路径模式, 会话 ID)**，` / `**会话结束即失效**；面板「安全」页实时列出所有生效中的授权并可一键撤销；` / `每次使用被授权通道都写 \`tool_call\` 日志（可取证）。`
- `docs/PLAN.md:2152-2155`（三条不可逾越）：`- **D30 能力组合闸门（C25 污染命中）触发的升级，不受会话授权覆盖** ——` ／ `- **安全相关配置（\`[risk]\`/\`[fs]\`/\`[net]\`/\`[plugins]\`）的放宽不适用会话授权**，` + `必须走 D36 的重新确认`
- `docs/PLAN.md:2157-2159`：`⚠ **必须改写原条目而不是新增一条** —— 否则后续 agent 读到旧的 RESERVED 会认为「会话授权无实现计划」，把 D45 当遗漏补一遍或直接不做（这正是 M0 要防的事）。`
- `docs/PLAN.md:2187`：**会话的文档定义（仓里唯一一处成文定义）**：`- **会话**定义：从唤起到用户显式结束，或**空闲 90s**（可配 \`[session] warm_timeout_sec\`）。`
- `docs/PLAN.md:2193`：`- 90s 无交互 → \`Settling\` → Dispose SessionScope（走 C11，含 \`debug.FreeOSMemory()\`）→ \`Sleeping\`。`
- `docs/PLAN.md:2704`：`| \`approval_grant\` | \`id\` PK, \`scope\`, \`tool\`, \`pattern\`, \`session_id\`, \`created_at\`, \`expires_at\`, \`revoked_at\` | \`session_id\` | 会话结束即失效（行保留 30 天供审计） | **D45 作用域会话授权** |`
- `docs/PLAN.md:2703`：`… \`grant_id\` 关联 D45 的会话授权`（`tool_call` 那一行；对应 A3 末「有列无来源」）
- `docs/PLAN.md:2983`：`| R4 | **C25 污染命中**（敏感读 × 外泄通道） | → **L2**，且**不受 D45 会话授权覆盖** |`
- `docs/PLAN.md:1592`：`**D45 双管**：批量聚合（500 个 L1 → 一次确认）+ 作用域会话授权（(工具,路径模式,会话) 三元组）。**且 L2 永不聚合、永不进持久授权**`
- `docs/PLAN.md:2569`（`shell.session` 那枚工具，⚠ **同名不同物**，别混）：`| **\`shell.session\`** | **在一条活着的会话里跑下一条命令（工作目录与环境延续）** | **L2** | \`shell\` | S3 | 票 163（2026-09-27 owner 批准新增）。**一次性外观、常驻内核**；**会话寿命绑任务作用域**；超时＝真 deadline（禁墙钟差）…`
  ⇒ 这是仓里**唯一一处把「会话寿命」明确绑到「任务作用域」的成文文字**，但它绑的是 shell 子进程会话，不是授权 scope。
- `docs/PLAN.md:3115`（S7 切片）：`增：**D45-2 作用域会话授权** …` ⇒ **D45-2 的落点是 S7**，今天还没到那一格。
- `docs/PLAN.md:3311`：`├── security.html        ← 生效中的会话授权（D45）+ 黑名单 + 污染告警历史`

**`docs/specs/SPEC-02-data-storage.md`**：`:92-102` 的 `CREATE TABLE approval_grant` 与 `internal/memory/schema.go:89-99` **逐字一致**（表在代码里已建成）。
`:79` `decision TEXT, -- 'allow'|'allow_session_grant'|'reject'|'timeout'|'batch_aggregated'`（`journal.go:32` 的常量就是抄这一行）；
`:86` `grant_id INTEGER -- 关联 approval_grant（D45）`；`:158` `| \`approval_grant\` 失效行 | 保留 30 天供审计 | 同上 |`；
`:214` `不做跨进程共享 DB（单进程单实例语义由 per-session 互斥保证，SPEC-09 §5）`（此处 per-session 指**登录会话**，不是授权 scope）。

**`docs/specs/SPEC-06-security-gatekeeping.md`**：`:117` `2. **作用域会话授权（S7）**：确认卡第三选项「本会话内允许 \`<工具>\` 于 \`<路径模式>\`」；`
`:121-122` `3. **三条不可逾越**：L2 永不进持久授权（含会话级）；C25/R4 升级不受会话授权覆盖；` + `安全配置（\`[risk]\`/\`[fs]\`/\`[net]\`/\`[plugins]\`）放宽不适用会话授权（走 D36 重新确认）。`
`:37`（R4 行不受会话授权覆盖）、`:62`（A 档「任何授权/会话授权/配置放宽无效」）、`:152` `批量聚合：L1 聚合、L2 不聚合、R7（≥50 文件）升 L2、会话授权过期即失效。`

**`docs/specs/SPEC-05-agent-core.md:135-136`**：`- 会话结束 → \`Settling\` → Dispose SessionScope（C11 全程，含 \`FreeOSMemory\`）→ 卸模型/` + `隐藏面板/失效会话授权（D45-2 的授权随 session_id 过期）。`
⇒ **这一行给了「失效」的机制说法：随 \`session_id\` 过期**（不是删行、不是 revoke），与 `expires_at` 注释（A4）同构。〔仅文档写了、代码没有〕

**`docs/specs/SPEC-00-product-overview.md`**：`:71` 场景 30 `作为用户，我想在确认卡上选择「**本会话内允许此工具于此路径模式**」，以便同类重复操作只确认一次且会话结束自动失效（D45-2）`；`:75` 场景 34 `…面板「安全」页看到**所有生效中的会话授权并可一键撤销**，以便授权永远在掌握中（D45-2）`。

**C17 契约（`docs/specs/SPEC-08-ui-ball-panel.md` §5.2，`:156-176`）** —— ⚠ 题头逐字**【SPEC 提案，S5 定稿走契约批准】**，即**尚未定稿**：
- `:171` `| \`grants.list\` / \`grants.revoke\` | invoke | — |` ⇒ **C17 白名单里已经写了两枚 grants 方法，Go 侧零实现**（A5：`knownComposerMethod` 只认四枚）。〔仅文档写了、代码没有〕
  ⚠⚠ 这两枚名字**与 E 节那枚冻结钉直接对撞**（`carriesGrantWord` 归一化后含 `grant`）——见 E3。
- `:158` `每方法标注 capability 与是否需原生侧授权；未列出方法名 → 拒绝并记日志。`
- `:168` `| \`approval.decide\` | invoke | **「allow」拒绝一切面板来源（F2）；仅 \`reject\` 可面板发起** |`
- `:145` `单例 \`PanelManager\`：一会话至多一个 WebView2 窗口`（panel 侧的「会话」，无身份、无常量、无字段）
- `:92` `| 4 | \`Sleeping\` | 快捷键/单击球 | \`Listening\` | **建 SessionScope(C31)**；加载 VAD+ASR（冷会话 1–3s） |`
- `:121` `| 31 | \`Warm\` | 90s 无交互 | \`Settling\` | **Dispose SessionScope**；卸 ASR/TTS；FreeOSMemory；销毁或隐藏面板 |`
- `:100` `| 12 | \`Listening\` | 超时（首轮 15s/会话内 90s/Conversation 30s） | \`Sleeping\`(首轮)/\`Warm\`(会话内) | — |`
- `:226` `Esc 在 Confirming 期间临时接管为取消键，**会话结束必须归还**`（A5 的 `ReleaseEscAfterSession` 就是这条的实现）

**`docs/specs/SPEC-01-architecture.md:64`**：`│   ├── session/                       # SessionScope(C31)、Warm/Conversation 计时、模型引用计数、D45 授权登记`
**`docs/specs/SPEC-01:59`**：`│   ├── panel/                         # WebView 宿主(C27)、PanelBridge(C17)、webassets embed`

**`docs/specs/SPEC-12-roadmap-governance.md:83`**：`| RESERVED | 跨会话持久授权档（workspace/user） | D45 后仅剩跨会话档；持久授权×提示注入风险窗口大 | — | C18 决策类型可扩展 | 每会话重新授权 |`
⇒ 登记表这句「会话档已实现、只剩跨会话档」**与 A1–A6 的现读不一致**（现读：会话档本身零执行者）。只如实记，不在本件裁定。

### A8 汇总：候选载体 × 有没有明确结束点 × 有没有身份

| 候选 | 有明确结束点？ | 有可携带的身份？ | 今天被授权链读到过？ |
|---|---|---|---|
| `wisp run` 进程 | **有**（`cmd/wisp/main.go:85` 的 `os.Exit`） | **没有** | 否 |
| `agent.Loop` 的一次任务（`task_id`） | **有**（`internal/tools/bridge.go:749` `CloseTask`，挂在 `cmd/wisp/run.go:733` 的 `AdmitTask`） | **有**（`internal/agent/loop.go:1123 newTaskID()`） | 否（且 DDL 绑的是 `session_id`） |
| `Loop.Reset()`（清对话历史） | 形式上有（`loop.go:265`） | **没有** | 否 |
| 球/面板的一次交互会话（状态族） | **有**（`statemachine/table.go:215/:227/:231`；视觉侧 `ball/liquid_windows.go:153-167`） | **没有**（`beginSession()` 不产 id） | 否 |
| `C31 SessionScope`（`internal/session`） | 文档有（`PLAN.md:2187`） | **代码零实现**（包里只有 `doc.go`） | 否 |
| `approval_grant.expires_at` | 注释有（`schema.go:96`「会话结束时间」） | — | 只有 retention 按墙上时间删 |
| `[session]` 配置节 | 不是身份，且**三枚键零读者** | 没有 | 否 |

⇒ **A 的总结论**：全仓今天**没有任何一枚「会话身份」**。有结束点的（进程／task／状态族）都不铸身份；
铸了身份的（task id）在文档上不属会话这一档。**「会话」这一档的执行者＝0，这一点在票 224 的现量之外，本普查独立复现。**

## B. 判定链今天怎么走到「要不要问」

### B0 ⚠ 先更正尺名：**`internal/risk/mode.go` 里没有 `resolveMode` 这枚函数**

`grep -n 'func \|resolveMode' internal/risk/mode.go` 现跑得到全部函数：`DefaultMode:70`、`ModeNames:74`、`ParseMode:82`、`(*ModeError).Error:99`、`quote:105`、`Mode.Valid:108`、`Mode.String:112`、`LoosestOf:129`、**`Mode.Screen:157`**、`redLine:197`。
`grep -rn 'resolveMode\|ResolveMode' --include='*.go' internal cmd tools` **零命中**。
⇒ 票 224 表里那句「`resolveMode` 一带」指的是 **`func (m Mode) Screen(d Decision) Silenced`（`internal/risk/mode.go:157`）**。本件按现读函数名与行号引。

### B1 全链（现读，逐枚行号）

| # | 发生了什么 | file:line |
|---|---|---|
| 1 | 唯一的生产调用方：`func (b *Bridge) Execute(ctx, req agent.ToolRequest)` —— 它是 `agent.ToolProvider` 的实现，被文本循环调 | `internal/tools/bridge.go:245` |
| 2 | C3 能力检查（在任何风险工作之前）：`dec, why := b.checkCaps(req, entry)` | `:258`、`checkCaps:327` |
| 3 | 参数解码 `decodeArgs(req.Args)`；`pathArgs(params, entry.Decl.PathParams)` 取路径参数 | `:265`、`:275`、`pathArgs:1052` |
| 4 | **C19 唯一判定**：`verdict := b.assessorFor(req.TaskID).Assess(req.Name, params, b.factsFor(ctx, entry, params, rawPaths))` | `:276-277`、`assessorFor:776` |
| 5 | 判定结果搬进 `dec`：`dec.Level / dec.RulesHit / dec.Reason / dec.SessionOverrideBlocked` | `:278-280` |
| 6 | **档（mode）的读取**：`mode := b.permissionMode()` —— 每调用读一次 | `:287`、`permissionMode:35`（`internal/tools/mode.go`） |
| 7 | **屏（screen）这一步**：`sil := mode.Screen(verdict)` | `:288` ← **这一行就是「要不要问」的全部决策** |
| 8 | 记账用的回声：`dec.Mode / dec.ModeSilenced / dec.ModeKept` | `:289-291` |
| 9 | 审计行：`MODE-SILENCE` / `MODE-REDLINE`；⚠ 文案逐字「a silenced question, not an allow: no user click happened」 | `:295-303` |
| 10 | `b.emit(dec)`（把结构化判定交给宿主） | `:304` |
| 11 | **落到问不问**：`ok2, why := b.route(ctx, &dec, sil)` | `:307`、`func (b *Bridge) route:371` |
| 12 | `route` 按 `sil.Level` 分派：`switch sil.Level`（`:372`） | |
| 12a | `case risk.L0:`（`:373`）→ `dec.DecisionColumn = agent.DecisionAllow`、`return true, ""` ⇒ **不问、直接执行** | `:373-375` |
| 12b | `case risk.Deny:`（`:377`）→ `DecisionReject`，**根本不进闸门**（R3 A 档） | `:377-380` |
| 12c | `case risk.L1:`（`:381`）→ `a, why := b.gate.PendingWindow(ctx, *dec)`；`AnswerAllow`/`AnswerTimeout` 都放行（L1 到点没被否＝执行），`AnswerVeto` 拒 | `:381-395`（`:382` 真调 `PendingWindow`） |
| 12d | `case risk.L2:`（`:396`）→ `a, why := b.gate.PendingApproval(ctx, *dec)`；只有 `AnswerAllow` 放行，`AnswerTimeout` **判拒**（与 L1 反极性） | `:396-411`（`:397` 真调 `PendingApproval`） |
| 12e | `default:` 不可能等级 → fail-closed 拒 | `:413-416` |

**入参形状**：`Screen` 收 `risk.Decision`（判定器 C19 的产物，含 `Level`/`RulesHit`/`Reason`/`SessionOverrideBlocked`）。
**返回形状**：`risk.Silenced{Level, Silenced bool, Kept string, Mode Mode}`（`internal/risk/mode.go:138-152`）。

### B2 `Screen` 内部的三态（这一层就是今天全部的「上界」逻辑）

现读 `internal/risk/mode.go:157-189`，按分支顺序：

| 分支 | 条件 | 结果 | 行 |
|---|---|---|---|
| Deny 优先 | `d.Level == Deny` | `Silenced{Level: Deny, Kept: "拒绝级 verdict（Deny 级不受任何模式影响）"}` | `:160-161` |
| 未知级 fail-closed | `d.Level != L1 && d.Level != L2` | 保留询问 | `:163-164` |
| 最严档短路 | `m == ModeAskEveryStep \|\| !m.Valid()` | `Silenced{Level: d.Level}`（照问） | `:167-169` |
| L1 静默 | `d.Level == L1` | `Silenced{Level: L0, Silenced: true}` | `:174-175` |
| ask_high_risk 遇 L2 | `m == ModeAskHighRisk` | `Silenced{Level: L2, Kept: "ask_high_risk 只静默 L1，本判定是 L2"}` | `:178-182` |
| auto_approve 遇 L2 的红线 | `redLine(d)` 命中 | `Silenced{Level: L2, Kept: why}` | `:185-187`、`redLine:197` |
| auto_approve 遇 L2 非红线 | 否则 | `Silenced{Level: L0, Silenced: true}` | `:188` |

`redLine`（`:197-220`）是**排除表**，命中的规则名逐枚：`R4 污染升级（C25）`、`R1`、`R2`、`R3`、`R5`、`R8`、`R9`；
文件头 `:36-40` 逐字：`Silence-immune rule set: R1 (declared L2), R2 (outside the authorized allowlist = "writes outside the workspace"), R3 (sensitive path A/B), R4 (taint), R5 (network target = "外网"), R8 (irreversibility), R9 (fail-closed).`
⇒ **今天判定链里没有任何一处会去问「这一条 grant 命中了吗」。** 全链唯一的外部读入是 `permissionMode()`（一枚 `risk.Mode` 整数）与 `b.confirmations()`（B 档单文件集合）。〔已证（现读）〕

### B3 若加「grant 命中则不问」这一步：**最小落点**与它能拿到什么

| 候选落点 | 精确行 | 该位置手头的材料 | 评估 |
|---|---|---|---|
| **甲（推荐的最小品）**：在 `Screen` 之后、审计 switch 之前插一步「命中则把 `sil.Level` 降为 L0 并记 Kept」 | `internal/tools/bridge.go:288` 之后（即第 289 行前） | `req.Name`（工具）、`params`、`rawPaths`／`dec.Paths`、`req.TaskID`、`ctx`、`verdict`、`mode`、`sil` | **这是唯一一处同时拿着「档」与「判定」的地方**；`route()` 只拿 `sil`（`:371` 的签名把 `sil` 作参数是刻意的，注释 `:366-370` 逐字「sil is a parameter, not a lookup」），所以放在 `route` 里会破坏那枚 AC#1 论证 |
| 乙：在 `route()` 的 `case risk.L1:`／`case risk.L2:` 之前各加一次命中查询 | `internal/tools/bridge.go:381`、`:396` | `ctx`、`*dec`（含 `Tool`/`Paths`/`RulesHit`） | 能拿到工具名与路径，但**要改两处**且把「授权」塞进"分派器"，与 `:362-370` 的注释（route 只做 level→branch 映射）冲突 |
| 丙：在 `risk` 包里加一步 | `internal/risk/mode.go:157` 之内 | 只有 `Decision` 与 `Mode` | ⛔ **拿不到任何 store**，且 `mode.go:3-9` 逐字声明这一层「adds no rule, renumbers nothing and changes no verdict … it may only ever make the answer STRICTER than the mode would otherwise allow」——**授权是放松方向，写在这一层违反它自己的文件头** |

⇒ **最小落点结论（现读支撑）：`internal/tools/bridge.go:288` 与 `:289` 之间**。

### B4 那个位置**能不能拿到 memory 的 store**（今天谁持有 store 引用）

| 侧 | 具名 file:line | 事实 |
|---|---|---|
| `internal/tools`（判定链所在） | `internal/tools/bridge.go:32-97`（`Options` 全部字段）、`:135`（`modes ModeSource`）、`:90`（`Modes`）、`:96`（`Confirmations`） | **`Options` 里没有任何 store／grant 源字段**。今天只有两枚只读注入：`ModeSource`（`internal/tools/mode.go:27-30`，接口只有一个方法 `PermissionMode() risk.Mode`）与 `Confirmations func() map[string]bool`（`:96`）。⇒ **桥拿不到 `memory.Store`。**〔已证（现读）〕 |
| `internal/perm` | `internal/perm/store.go:47-54`：`type ConfigManager interface { Config() *config.Config; SetPermissionMode(mode risk.Mode) error }` —— **`Options.Manager` 只有配置这一面** | `grep -rn 'memory\.' internal/perm/*.go \| grep -v _test` **零命中**（只有 `store.go:45` 注释里那句英文单词）。⇒ **`internal/perm` 不持有 store 引用。**〔已证（现读）〕 |
| `internal/agent/approval` | `grep -rn 'memory\.' internal/agent/approval/*.go \| grep -v _test` **零命中**（唯一命中是 `pending_read.go:70` 注释里的英文 "fresh memory"） | ⇒ **闸门／队列也不持有 store 引用。**〔已证（现读）〕 |
| **生产装配根（唯一持有者）** | `cmd/wisp/run.go:229`：`store *memory.Store`（`agentRuntime` 的字段）；`:381 rt.store = mem`（开库）；`:680-682` 关停里 `rt.store.Close()`；`:919 type storeHealthSink struct{ store *memory.Store }` | ⇒ **今天 `memory.Store` 的引用只活在 `cmd/wisp` 的 `agentRuntime` 上。**〔已证（现读）〕 |
| 装配点（桥是在这里组的） | `cmd/wisp/run.go` 里 `func assembleRuntime(s runSpec) (*agentRuntime, int)`（`:321`），`opts := agent.Options{ Tools: rt.bridge … }`（`run.go:756-757`） | 组桥那一段就是注入 grant 源的位置——**它同一时刻手里既有 `rt.store`（`:381` 已经开好）又有 `rt.bridge`**。〔已证（现读）〕 |

**B 的总结论**：
1. 判定链今天**只读两样外部状态**：一枚 `risk.Mode` 整数、一枚 `map[string]bool`（B 档已确认文件集）。**没有第三处给授权留的口子。**
2. 加「grant 命中则不问」的最小落点是 `internal/tools/bridge.go:288` 之后一行，材料齐全（工具名＋参数＋路径＋task id）。
3. **但那一行今天拿不到 store**：`memory.Store` 只在 `cmd/wisp/run.go:229`/`:381` 被 `agentRuntime` 持有，`internal/tools`/`internal/perm`/`internal/agent/approval` 三个包都零引用 ⇒ 要接就得**新增一枚注入 seam（照 `ModeSource` 那形的只读接口）**，这本身就是一枚新的契约面（B 节末的 `Confirmations nil` 现状由 `cmd/wisp/run_mode101_test.go:429-430` 的注释逐字记录着：「the assembly hands the bridge a mode, never a grant source (Options.Confirmations nil)」）。〔已证（现读）〕

## C. 三件套各自的落点（写／读／失效）

### C0 前置事实：三样东西今天都不存在（现跑）

| 缺口 | 尺与读数 |
|---|---|
| **答复里没有「档」这一维** | `grep` 现读 `internal/tools/gate.go:77-88` 的 `Answer` 词表只有四枚：`AnswerAllow "allow"`（`:78`）、`AnswerReject "reject"`（`:80`）、`AnswerVeto "veto"`（`:83`）、`AnswerTimeout "timeout"`（`:88`）。**没有 "allow_session_grant" 这一支**——它在 `internal/agent/journal.go:32` 只是**记账词**，不是**答复词**，且零写者（见 A2） |
| **请求身上没有会话位** | `internal/agent/tools.go:54 type ToolRequest struct` 的字段（逐枚现读）：`TaskID / CorrelationID / CallID / Name / Args / Timeout`。**没有 `SessionID`、没有 `Session`** |
| **没有任何路径模式匹配器** | `grep -rn 'MatchString\|filepath.Glob\|fnmatch\|pathMatch\|MatchPattern' internal/risk internal/tools internal/memory --include='*.go' \| grep -v _test` → **零命中**。`approval_grant.pattern` 那一列（`schema.go:93`）今天**没有任何代码会去比对它** |

### C1 写：答复「本会话内允许」→ 真落一行

| 项 | 具名 file:line | 状态 |
|---|---|---|
| **今天答复的入口（⚠ 正被一枚写腿改，本件只读、不保证内容稳定）** | `cmd/wisp/approval_reply.go`（`grep -n 'func \|Ticket'` 现读）：`newNativeCards:113`、`(*nativeCards).bind:119`、`record:133`、`look:142`、`forget:150`、`pending:158`、`waitingState:168`；`(*replySurface).allow:208`、`reject:235`、`refuse:243`、`panelReject:266`、`panelAllow:276`、`veto:304`、`head:321`、`view:330`、`record:361`；`(*agentRuntime).attachReplyListener:389`、`runReplyLoop:433` | 〔已证（现读）〕，⚠ **不稳定** |
| 答复真正落到的闸门 API | `internal/agent/approval/gate.go:626 (*Gate) DecideFromNative`、`:638 (*Gate) DecideFromPanel`、`:387 (*Gate) Veto`；`internal/agent/approval/queue.go:343 func (q *Queue) allow(corr, nonce string) error`（⚠ **只有两枚参数**，没有 scope／reason 位）；`internal/agent/approval/gate.go:596 func (n nativeAPI) Allow(_ context.Context, corr, grant string) error` | 〔已证（现读）〕⚠ 这三枚 `gate.go` 行号与票 219 顶部记录的（`:622`/`:610`/`:371`）**已经不同**——写腿在动这个文件，引它们之前必须重跑 |
| **最近的半成品（写侧）** | `internal/memory/dao_misc.go:17 func (s *Store) InsertGrant(ctx, g ApprovalGrant) (int64, error)` —— 一行 INSERT 全好，带 scope 校验（`:18-19` 拒绝非 `session` 的 scope）与三必填校验（`:21-22`）；模型 `internal/memory/models.go:93-101 ApprovalGrant` | 〔建了但没接〕：**零生产调用者** |
| **缺什么（写侧，逐条）** | ① 答复词表里要有一枚带 scope 的答复（`internal/tools/gate.go:77-88` 现只有四枚）——**新增答复枚举值＝契约面**；② `Queue.allow` 的签名（`queue.go:342`）要能携带「这一答是会话档」，或答复侧另起一条落库腿；③ **一枚会话身份**（A 节证：全仓零铸造点）；④ `expires_at` 得有人填（`schema.go:96` 注释说它是「会话结束时间」，今天只有测试写 `now+3600`：`internal/perm/ticket90_persist_test.go:235`、`cmd/wisp/run_mode101_test.go:417`）；⑤ 写手要拿到 store——`memory.Store` 引用今天只在 `cmd/wisp/run.go:229`/`:381`，**闸门／队列／桥三个包零引用**（见 B4）；⑥ `tool_call.grant_id` 的联动今天被硬写死成 nil：`internal/agent/journal.go:93 func (t *taskJournal) decide` → `:101 _ = j.DecideToolCall(ctx, rowID, decision, nil)`——**接口有 grantID 形参（`journal.go:22`）而唯一的调用方递 nil**，且全仓非测试只有这一处调用 |
| ⚠ 一条已存在的**反向**证据 | `cmd/wisp/run_mode101_test.go:413`（`SessionID: session`，`const session = "session-before-restart"` 在 `:400`）＋ `:427-428` 注释逐字：「Refused with the grant live IN ITS OWN SESSION, too: the assembly hands the bridge a mode, never a grant source (Options.Confirmations nil)。」⇒ **这枚钉现在就在把「零读者」这件事钉成期望行为**（详见 E4） | 〔已证（现读）〕 |

### C2 读：下一次同类请求命中它

| 项 | 具名 file:line | 状态 |
|---|---|---|
| **最近的半成品（读侧，DAO 层）** | `internal/memory/dao_misc.go:78 func (s *Store) ListGrantsBySession(ctx, sessionID string) ([]ApprovalGrant, error)`，SQL 在 `:80-81`（`WHERE session_id=? ORDER BY created_at DESC, id DESC`）；另 `:91 ListGrants`（全表，审计） | 〔建了但没接〕：**零生产调用者**（现跑：非测试命中＝定义行自己） |
| **最近的半成品（读侧，链上路点）** | `internal/tools/bridge.go:287-288`：`mode := b.permissionMode()` → `sil := mode.Screen(verdict)`。这是**唯一一处同时拿着档、判定与工具身份**的位置，也是 B3 判定的最小落点 | 〔已证（现读）〕，**这里今天不查任何授权** |
| **注入 seam 的现成形状可抄** | `internal/tools/mode.go:27-30 type ModeSource interface { PermissionMode() risk.Mode }`，被 `internal/perm/store.go:149`（`PermissionMode implements tools.ModeSource`）满足；桥侧只读、**没有 setter**（`internal/tools/bridge.go:135` 字段、`:90` Option） | 〔已证（现读）〕——**照这个形加一枚只读 GrantSource 是改动面最小的路径，但它是一枚新的契约面** |
| **缺什么（读侧，逐条）** | ① 一枚带 session 身份的查询入参（`ToolRequest` 没有该字段：`internal/agent/tools.go:54`）；② **`(tool, pattern)` 的匹配器根本不存在**（C0 第三行：零 glob/前缀匹配代码）——`PLAN.md:2147` 要求的是「工具＋路径模式」双匹配，`PLAN.md:2569` 那一支还要按目标进程名分档；③ 红线不可覆盖的守卫：命中之后仍必须让 `R4/SessionOverrideBlocked`（`internal/risk/mode.go:198-199`）、A 档（`route` 的 `case risk.Deny:` `internal/tools/bridge.go:377`）、`安全配置放宽`（`PLAN.md:2154`）这三类**拒绝被授权放松**——`PLAN.md:1642`／`SPEC-06:121-122` 逐字；④ 「命中」要能被审计：`Silenced` 结构体（`internal/risk/mode.go:138-152`）今天只有 `Kept` 一个字符串槽位表达「为什么仍然问」，**没有「为什么没问（除 mode 之外）」的槽位**——`dec.ModeSilenced` 只说档，不说授权 |
| ⚠ 一条**方向性**约束（现读） | `internal/tools/bridge.go:362-370` 的注释逐字：「sil is a parameter, not a lookup: this is the one place a level becomes "ask the user" or "do not ask", so the mode that decided it has to be visible in the signature (AC#1). An implementation that reached for a global here would put the permission switch inside the enforcement layer.」⇒ **任何「命中即不问」的实现都必须把命中结果作参数传进 `route`，不能在 route 里自己去查** | 〔已证（现读）〕 |

### C3 失效：会话结束 → 查不到

| 项 | 具名 file:line | 状态 |
|---|---|---|
| **最近的半成品（事件源）** | 状态机那两条冻结边就是「会话结束」的成文事件：`internal/statemachine/table.go:227`（`Warm --EvWarmIdle--> Settling`）、`:231`/`:236`（`Settling --EvSettleExpired--> Sleeping`）；视觉侧已经在监听边沿：`internal/ball/liquid_windows.go:148-167 func (b *Ball) motionStateChanged`（`:157`/`:165` 调 `endSession()`）、`internal/ball/liquid.go:99 func (m *liquidMotion) endSession()` | **边沿有人听，但听的是动画**；「per-session grants ledger (D45-2)」这句话只写在 `internal/session/doc.go:10`，`internal/ball/*` 里零出现（现读 `grep -rn 'grants ledger' internal` 唯一命中就是那一行） |
| **最近的半成品（撤销 API）** | `internal/memory/dao_misc.go:49 func (s *Store) RevokeGrant(ctx, id int64) error`（`:53 UPDATE approval_grant SET revoked_at=? WHERE id=? AND revoked_at IS NULL`，幂等，`:61-65` 不存在则 `ErrNotFound`）；`:102 DeleteGrant` | 〔建了但没接〕：**零生产调用者** |
| **唯一的自动清理** | `internal/memory/retention.go:219 func deleteExpiredGrants`，SQL 在 `:225-227`：`DELETE FROM approval_grant WHERE MAX(expires_at, COALESCE(revoked_at, 0)) < ?`；注释 `:216-218` 逐字「an invalid grant stays 30 days for audit, then goes (SPEC-02 §4)」 | 〔已证（现读）〕⚠ **它按墙上时间删，不按「会话是否还活着」删**——`expires_at` 若没人按会话真填，这一步就永远轮不到 |
| **缺什么（失效侧，逐条）** | ① 有人把「会话结束」翻成一次**写**：`rt.close()`（`cmd/wisp/run.go:676-683`）今天只做 `replyRoot.Cancel()` 与 `store.Close()`，**不 revoke 任何 grant**；② 或者按 `SPEC-05:136` 说的「随 `session_id` 过期」——那要求**铸造点保证身份不重复**，而这正是 D 节要独立复核的那枚坑；③ C17 的「一键撤销」（`docs/specs/SPEC-08:171 grants.revoke`）今天零实现，而且**这枚方法名会被那枚冻结钉判死**（见 E3）；④ 撤销之后「同类请求重新弹卡」需要读侧先存在（C2），否则「失效」无从验证 | 〔仅文档写了、代码没有〕（①②③④ 全部） |

## D. 那三枚钉的射程（独立复核）

> ⚠ 本节**没有采用编排者给的结论**：我自己 `sed -n '1,40p'` 与 `sed -n '138,292p'` 逐枚读了三枚函数的每一条断言，
> 再回答最后一问。**复核结果：编排者那句判断成立**，但**它低估了一格**——守卫这个形状的钉**不是三枚而是四枚，四枚全在同一方向上瞎**（见 D5(b)）。

### D0 文件头的逐字声明（现跑 `sed -n '21,31p' internal/perm/ticket90_persist_test.go`）

```
:21 // Ticket 90, part 2: the persistence boundary R20/M3 drew, and the two AC
:22 // cases that must NEVER be merged into one.
:23 //
:24 //	AC#3(a)  a manually chosen mode survives a restart           (perm: YES)
:25 //	AC#3(b)  a never-touched config reads the default after start (perm: YES)
:26 //	AC#3b    a D45 session grant does NOT survive a restart      (grant: NO)
:27 //
:28 // They are three test functions on purpose. The orchestrator's ruling says a
:29 // single "persistence" case is how a permanent免审通行证 gets in: one green
:30 // checkbox that silently covers two opposite requirements can only ever be
:31 // satisfied by the looser one.
```
另有 `internal/perm/store.go:22-28` 把同一条边界写在生产侧（现读逐字）：
```
:22 //   - It does not touch authorizations. R20/M3 persists the MODE, one global
:23 //     preference; it does not make a D45 session grant persistent. A grant stays
:24 //     session-scoped and dies with the process (PLAN.md:1640), and ticket 49's
:25 //     GrantScopeSession keeps that meaning. The two are pinned by two separate
:26 //     tests (AC#3a and AC#3b) so nobody can merge them into one "persistence"
:27 //     case later - a permanent免审通行证 is exactly what this boundary exists to
:28 //     keep out.
```
⚠ 注意 `store.go:25` 说的是「**two** separate tests」，文件头 `:28` 说的是「**three** test functions」——
`grep -n 'func Test'` 现读这个文件共 **5 枚**函数：`:138`、`:181`、`:220`、`:291`、`:336`。三枚是 AC#3a/AC#3b/AC#3(b)，另两枚（`:291 TestTicket90ConfigKeyChangesWhatTheChainAsks`、`:336 TestTicket90HandEditLooseningGoesThroughD36`）管的是键不许说谎。这不是矛盾（"两枚独立的持久化用例" vs "三枚函数"），但**引用时别再写成两枚**。

### D1 钉 #1 —— `internal/perm/ticket90_persist_test.go:138 TestTicket90ManualSwitchSurvivesRestart`

| 它**看得见**什么（现读逐枚断言） | 它**看不见**什么 |
|---|---|
| `:150-152` 起始档必须是 `ModeAskEveryStep`；`:153-163` 循环切到 `AskHighRisk`、`AutoApprove`，每次都用 `openManager(t, path)` 重开一个 Manager（`:159`）再读 `PermissionMode()`（`:160`）；`:164-166` `confirmCalls != 1`（只有切进 auto_approve 要花一枚确认）；`:169-176` 直接读 `config.toml` 原文，断 `permission_mode` 与 `auto_approve` 两串字面都在文件里 | ⛔ **整个函数体一次都没碰 `memory.Store`、没碰 `approval_grant`、没碰桥／闸门／工具调用**。它导入 `internal/memory`（`:16`）是**同文件其它函数**用的。⇒ 会话授权怎么写、什么时候失效、身份从哪来，**这枚钉与它零关系** |

### D2 钉 #2 —— `internal/perm/ticket90_persist_test.go:181 TestTicket90UntouchedConfigStartsAtTheDefault`

| 看得见 | 看不见 |
|---|---|
| `:183-187` `config.NewDefaults()` 的 `Risk.PermissionMode` 必须是 `ask_every_step`；`:188-191` 存盘；`:192-205` 同一文件**冷启三次**，每次 `New(Options{Manager: mgr})`（`:194`）后断 `PermissionMode() == ModeAskEveryStep`（`:198-200`），并断 `askOnce(t, path)` 返回 0（`:202-204`，即从没够到过 L2 卡） | ⛔ 同样**零 grant、零 store、零判定链**。这枚钉只管"没人写过配置 ⇒ 读到默认"。**会话身份的铸造方式对它完全不可见** |

### D3 钉 #3 —— `internal/perm/ticket90_persist_test.go:220 TestTicket90SessionGrantDoesNotSurviveRestart`

**逐枚断言（现读，含它自己写死的字面量）**：

| 行 | 做的事 |
|---|---|
| `:225` | `store, err := memory.Open(dir, …)`（真库，真 temp dir） |
| `:230-237` | `store.InsertGrant(ctx, memory.ApprovalGrant{Scope: GrantScopeSession（`:231`）, Tool: "fs.write"（`:232`）, Pattern: filepath.Join(dir, "*")（`:233`）, **SessionID: "session-before-restart"（`:234`）**, CreatedAt: now（`:235`）, ExpiresAt: now+3600（`:236`，注释逐字 `still live by its own clock: only the SESSION dies`}）` |
| `:241-243` | `id == 0` → Fatal |
| `:244-247` | `store.ListGrantsBySession(ctx, "session-before-restart")` → 断 `len(live) == 1`（**同一枚字面量查询**） |
| `:248-250` | `store.Close()` |
| `:252` | 注释逐字：`// "Restart": a new process, therefore a new session id, over the same DB.` |
| `:253-257` | `reopened, err := memory.Open(dir, …)` 同一个目录 |
| `:258-264` | `reopened.ListGrantsBySession(ctx, **"session-after-restart"**)` → `len(next) != 0` 才 Errorf（`:263` 逐字 `a session grant survived the restart (%d rows for the new session)`） |
| `:265-273` | `reopened.ListGrants(ctx)`（**不带 session 过滤**）→ 断还剩 **1 行**、且 `all[0].SessionID == "session-before-restart"`（`:271`）——即**它明确要求那一行活着留在库里** |
| `:275-284` | 注释逐字：`// Deliberately NO mode assertion here. … this test fails ONLY if a session grant becomes durable … Wiring the contrast into this test would make AC#3b depend on AC#3a's mechanism, which is exactly the merge the ruling forbids` |

**看得见**：DAO 的 `WHERE session_id=?` 真的按字面量过滤；行留在库里供审计（这正是 `PLAN.md:2704`／`SPEC-02:158` 要的形状）。
**看不见**（逐条，全都能同时成立而此钉全绿）：
1. ⛔ **看不见生产怎么铸 id**——全树没有任何铸造点（A 节），而这枚测试从不向生产要 id，它自己打了两枚字符串（`:234`、`:244`、`:258`）。
2. ⛔ **看不见判定链**——函数体里没有 `tools.New`、没有 `Bridge.Execute`、没有 `Gate`、没有任何一次工具调用。它断的是"查得到/查不到行"，**不是"问不问"**。
3. ⛔ **看不见匹配器**——`Pattern` 写的是 `dir + "/*"`（`:233`），而**全仓没有任何代码比对过这一列**（C0 第三行：零 glob/前缀匹配）。⇒ 未来匹配器把模式放得再宽，这枚钉也不会红。
4. ⛔ **它的"失效"判据是"换一枚字面量"**，不是"会话真的结束了"。⇒ 生产只要**在两次启动之间造不出这个字面量**（这是必然，因为生产从不产这串字），这一支**永远绿**。

### D4 第四枚同类钉 —— `cmd/wisp/run_mode101_test.go:389 TestTicket101SessionGrantDoesNotCrossRestart`

这枚**在本票列的三枚之外**，而且**它才是唯一真的走了一遍判定链的会话授权测试**：

| 行 | 做的事 |
|---|---|
| `:400` | `const session = "session-before-restart"`（**和钉 #3 同一枚字面量**） |
| `:411-418` | `rt.store.InsertGrant(...)` 写一行活授权（`:415 SessionID: session`、`:417 ExpiresAt: now + 3600, // live by its OWN clock: only the session may die`） |
| `:422-426` | `rt.store.ListGrantsBySession(ctx, session)` → `sameSession = len(rows)` |
| `:427-428` | 注释逐字：`// Refused with the grant live IN ITS OWN SESSION, too: the assembly hands` / `// the bridge a mode, never a grant source (Options.Confirmations nil).` |
| `:429` | `_, firstWhy, firstErr = t101call(rt, 1, "fs.write", env)` —— **一次真工具调用** |
| `:434-437` | 断 `sameSession != 1` 则 Fatal（逐字 `every assertion below would be vacuous`，防"断言空转"） |
| `:438-441` | 断 `!firstErr` 则 Fatal（逐字 `boot 1 already let the B-tier .env write through … the grant boundary cannot be measured against a chain that never asked`）——即**"授权活着也照样被拒"是它要的正确答案** |
| `:453-456` | 重启后 `rt.store.ListGrantsBySession(ctx, "session-after-restart")` → `after = len(rows)`（`:457`） |
| `:459` | 第二次真调用：`_, why, refused = t101call(rt, 1, "fs.write", env)` |
| `:469-471` | `if after != 0 { t.Errorf("%d grant rows are visible to the new session, want 0 (PLAN.md:1640)") }` |

⚠ **它能看见判定链，但它的期望值恰好把"零执行者"钉成了正确答案**（`:438-441` 与 `:427-428` 那两行注释互相印证）。
由此得到两条**方向不同**的射程判断，别混着说：
- **对"匹配器丢掉 session 条件"这一形，它不瞎**（这是全树唯一有此射程的一件）：一旦 grant source 真被接上，
  若实现只比 `tool`+`pattern` 而不比 `session_id`，`:429` 那次调用就不再被拒 ⇒ `:438-441` Fatal。**能红。**
- **对"派生式 id 让两次启动铸出同一枚 id"这一形，它瞎**：`InsertGrant` 写进去的是字面量 `session`（`:415`），
  生产铸的 id 永远不等于那串字，重启后查 `"session-after-restart"`（`:453`）永远得 0 行，`:469` 永远绿。
  ⇒ **同一枚坑（D5 要答的那一问）在这枚钉上也是绿的。**

### D5 回答那一问：**「生产若用派生式 session id，会让哪一枚钉看不见永久通行证？」**

**直接答案：三枚里没有一枚能看见它——只有钉 #3（`:220`）名义上关这一块，而它正是看不见的那一枚。钉 #1（`:138`）与钉 #2（`:181`）连授权都不碰，谈不上"看见"。**

推理链（每步都有 D0–D4 的现读支撑，不接受"应该差不多"）：
1. 永久通行证 = 「进程重启后，同类请求仍然不问」。要证明它，必须**同时**跑到 (i) 生产的铸造点、(ii) 生产的匹配步、(iii) 一次真的"问不问"。
2. 钉 #1／#2：三样一样都没跑（D1、D2 的"看不见"列）。派生式 id 对它们零影响 → **必绿**。
3. 钉 #3：跑了 (i) 吗？**没有**——它用的是测试自己打的 `:234`／`:244`／`:258` 两枚字面量，全仓生产的铸造点是**零**（A8）。跑了 (ii) 吗？**没有**——匹配器根本不存在（C0）。跑了 (iii) 吗？**没有**——函数体里没有桥、没有 gate、没有工具调用。⇒ **派生式 id 让它必绿，而通行证可以完全真实生效。**
4. 更要紧的是它的 `:265-273` 一支**主动要求那一行留在库里**（`want the one old-session row kept for audit`）——
   也就是说：**"库里还有这行"不是这枚钉的失败信号**。任何"审计留痕 vs 通行证"的读法冲突，都被它写死成"留痕是对的"。

**与编排者那句话的差异（明写，不附和）**：
- ✅ **成立的部分**：「钉用的是测试自己写死的两枚字符串 `session-before-restart`／`session-after-restart`，看不见生产怎么铸 id」——
  现读 `internal/perm/ticket90_persist_test.go:234`、`:244`、`:258` 逐字核对，两枚字面量、无生产铸造参与。**这一句我独立复现，不推翻。**
- ⚠ **差异一（数目）**：把守卫数成"三枚"**低估了**。真实是**四枚同方向**（外加 `cmd/wisp/run_mode101_test.go:389`），
  而这第四枚**是三枚里唯一真跑过判定链的**，它对"派生式 id"这一形**同样瞎**（D4 末段论证）。
  ⇒ 汇报时应写「**四枚钉全绿而通行证实际生效是可能的**」，不是「三枚」。
- ⚠ **差异二（射程不是零）**：说"钉看不见生产怎么铸 id"**容易被读成"这些钉什么授权缺陷都抓不到"**，这不准确。
  `run_mode101_test.go:438-441`（断"授权活着仍被拒"）**是能抓"匹配器只比 tool+pattern、不比 session"这一形的**——
  前提是实现真把 grant source 接上了（`:427-428` 注释说明今天没接）。
  ⇒ 准确的写法是：**四枚钉对"身份派生"这条轴全瞎，其中一枚对"匹配器丢 session 条件"这条轴有射程。要防的是前者，现有的射程恰好不覆盖它。**
- ⚠ **差异三（可执行判据）**：AC#4（票 224 的反控）如果照"三枚钉"来设计，会在"派生式 id"这一支**假绿**。
  本件给的**唯一足以判红的形状**是三件事缺一必空转：
  ① 测试必须**调用生产的铸造函数本身**（不是打字符串），跨两次启动取两枚 id 并**断它们不等**；
  ② 必须用那枚**生产铸出的 id** 去 `InsertGrant`，再用第二次的**生产 id** 去 `ListGrantsBySession`，断 0 行；
  ③ 必须走一次**真工具调用**（`Bridge.Execute` 那条腿），断"授权活着的同一会话里不问、新会话里重新问"。
  今天这三件事**一件都没有**（A8：无铸造点；C0：无匹配器；D3「看不见」第 2 条：钉 #3 不跑链）。

## E. 会不会撞别的钉

### E1 与 session／grant／approval 相关的测试名（**只枚举名字，一枚都没跑**）

`grep -rn 'func Test' --include='*_test.go' internal cmd tools | grep -iE 'session|grant'` 现跑全量（18 命中，去掉 6 枚 `internal/winsec/*` 的 NT ACL「grant」同名异物）：

| file:line | 函数名 | 与"铸造会话身份"这步的关系 |
|---|---|---|
| `internal/perm/ticket90_persist_test.go:220` | `TestTicket90SessionGrantDoesNotSurviveRestart` | 见 D3（本票的三枚之一） |
| `cmd/wisp/run_mode101_test.go:389` | `TestTicket101SessionGrantDoesNotCrossRestart` | 见 D4（第四枚，同方向） |
| `internal/memory/dao_test.go:378` | `TestGrantCostPluginStateDAO` | 只测 DAO 自身（`:383` 写、`:391` 拒非 session scope、`:394` 查、`:398`/`:401` revoke 幂等、`:408` ErrNotFound）——**不碰链** |
| `internal/agent/approval/queue_test.go:150` | `TestGrantIsSingleUseAndBoundToItsItem` | ⚠ **这里的 "grant" 不是会话授权**，是原生一次性 nonce（见 E4） |
| `internal/agent/approval/queue_test.go:269` | `TestReplayRedisplaysUnderAFreshGrant` | 同上（nonce） |
| `internal/agent/approval/ticket87_veto_l2_test.go:203` | `TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant` | 同上（nonce） |
| `cmd/wisp/approval_reply_201_test.go:340` | `TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject` | 同上（nonce），且它管的是**面板侧不许 allow** |
| `internal/risk/assessor_test.go:78` | `TestFusionR4BlocksSessionOverride` | **直接管这一票**：`:89-91` 断 R4 命中时 `d.SessionOverrideBlocked` 必须为真，注释逐字 `R4 hit must block D45 session override` |
| `internal/risk/provenance_test.go:987` | `TestScopeIDGrantsNoPower` | ⚠ 名字最容易被误当"这枚管住了 id"：现读函数体（`:988-997`）它只断 `OpenScope("task-1").ID()=="task-1"` 与 close 后 `p.scopes` 里不再注册。**它管的是 C25 污染作用域的句柄，不是授权身份；一枚新铸的 session id 落在它射程之外** |
| `internal/ball/liquid_test.go:152` | `TestBorderGlidesInAfterTheSessionEnds` | 视觉侧（A5），与授权无关 |
| `internal/panel/l2_grant_boundary_test.go:1229／1530／1547／1595／1807／1959／2021` | 七枚（见 E2） | **冻结件，一字不许动** |
| `internal/winsec/*`（6 枚，同名异物） | `…InheritedGrant…`／`…OwnGrants…`／`…KindFieldSays…AGrant…`／`…LeavesARealGrantToAnotherAccount` | NT ACL 的 grant，**不属本票射程**（`internal/winsec/private_set_sid_windows_test.go:372` 等） |

另附 `internal/agent/approval` 包内与"谁能答"直接相关的两枚（`grep -n 'func Test' internal/agent/approval/*_test.go` 现读）：
`queue_test.go:89 TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis`、`ticket97_alias_direction_test.go:75 TestAnAliasCanNeverBuyAnAllow`。
⇒ **"会话档答复"若被写成一种 allow，会撞上这两枚**：它们钉的是"来源不可伪造的 allow"，不是"档"。

### E2 那枚冻结钉的禁名词根表（**现读**，`internal/panel/l2_grant_boundary_test.go`，2380 行）

**它扫哪一枚包**：`dir := filepath.Join(root, "internal", "panel")` —— `:1231`（facet 1）与 `:1848`（三面一致那枚）。⛔ **`internal/tools`／`internal/risk`／`internal/memory`／`internal/agent/**`／`cmd/wisp` 不在它的 AST 射程内**（现读这两处目录参数）。

**表一：入向 JSON 键黑名单**（facet 3，`:1882-1886` 逐字 `candidates`，**十二个词根 × 大小写两拼 = 24 串**）：
```
:1883  "outcome", "Outcome", "allow", "Allow", "allowOnce", "AllowOnce",
:1884  "approved", "Approved", "grant", "Grant", "verdict", "Verdict",
:1885  "decision", "Decision", "decide", "Decide", "bypass", "Bypass",
       "override", "Override", "permit", "Permit", "authorize", "Authorize",
```
判红的那一句在 `:1907`：`%s binds the wire key %q, so a panel can address an approval decision through it (D33/F2, AGENTS.md ban #6). If %s is inbound, this is the finding; if it is not inbound any more, drop it from inboundTypeRegistry instead of dropping the check`。

**豁免规则（现读，这就是 `reason` 安全的完整理由）**：
1. **只作用在"入向解码目标"上**：`inboundTypeRegistry()`（`:1720-1725`）逐枚只有四型——`ComposerRequest`、`ModeRequest`、`AttachmentPayload`、`AttachmentRef`；循环体在 `:1887` `for _, name := range sortedRegistryNames(reg)`。⇒ **Go→前端出向结构体不受这张表管**（`ApprovalCardView` 那族是出向，见 `internal/panel/approval.go:48`）。
2. **只比"这一层自己的 JSON 键"**：`astTopKeys`（`:1733-1747`）跳过 `json:"-"`（`:1736`）、匿名且无 tag 的提升字段（`:1739`）、非导出字段（`:1742`）；`reflectionTopKeys`（`:1752+`）同口径（tag 名，缺 tag 用 Go 字段名——`:1878` 注释逐字「encoding/json falls back to the Go field name, so a field called Outcome is addressable as "outcome" and as "Outcome"」）。
3. ⇒ **`reason` 为什么安全**：它**不在这 12 枚词根里**（表一逐字可验），而且它**语义是"用户在陈述什么"，不是"系统在决定什么"**；票 219 顶部第 3 格已现读它早已存在且是**出向**：`internal/agent/approval/ui.go:35 Prompt.Reason`／`:54 PanelItem.Reason`、`internal/panel/approval.go:48 ApprovalCardView.Reason json:"reason"`＋`:52 ReasonKnown`。⇒ 答复入向再加一枚 `reason` **既不落词表、又已有同名字段先例**，是这张表下唯一安全的答复字段名。

**表二：路由名黑名单**（facet 1，`:192-195` 逐字 `grantRouteWords`，**11 个词根**）：
```
"approve", "approval", "grant", "allow", "permit",
"ratify", "authorize", "authorised", "decide", "decision", "verdict",
```
归一化判据在 `:242-250 func carriesGrantWord`：`strings.ToLower(strings.NewReplacer("_", "", "-", "", " ", "", ".", "").Replace(name))` 然后 **`strings.Contains`**（**下划线、连字符、空格、点全被抹掉 ⇒ `"grants.list"` 归一化成 `grantslist`，含 `grant` ⇒ 命中**）。
判红句在 `:1244-1246`：`Go answers the inbound route %q, whose own name is an approval decision (D33/F2, R20, AGENTS.md §1.2 ban #6)`。
另有显式候选清单 `:1251-1262`（`panel.approval.grant`、`panel.l2.grant`、`panel.grant`、`approval.grant`… 十枚）与一张乘积网格：`wantGrantRoutePrefixes = 8`／`wantGrantRouteWords = 11`／`wantGrantRouteSuffixes = 3`（`:1334-1336`），`sweepAssembledNames`（`:1433`）把 8×11×3×2（词根与词根+s）**全部问一遍真闸门** `knownComposerMethod`（`internal/panel/bridge.go:104`），**任何一个被答就判红**。

**还有两枚会误伤的机制（现读）**：
- `:1855-1860`：**注册表腐蚀守卫**——`for seed := range pkg.inboundSeeds { if _, ok := reg[seed]; !ok { t.Fatalf("decode destination %s has no reflection twin in inboundTypeRegistry…") } }`；`inboundSeeds` 是从包里**自动派生**的（`classifyDecodes` `:675-708`，凡是 JSON 解码目标落在同包结构体上就入册 `:700`）。⇒ **在 `internal/panel` 新增任何一枚会被 JSON 解码的结构体，都必须同时把这枚名字加进 `inboundTypeRegistry()`——而那枚文件一字不许动 ⇒ 这条路今天走不通。**
- `:1828` 的 `f5Probe`、`:2021 TestPlantedGrantWiringGoesRedInASnapshot` 的植入样本（`:2097`、`:2112`、`:2130`、`:2180`、`:2220`、`:2266`、`:2294` 各建一份合成目录）——**它们只在自己的临时目录里跑，不会碰生产文件**，但说明这枚钉的正控是"造出违形状必须红"。
- `:1982` 另有一枚六词 `smuggled` 子集：`outcome`／`allowOnce`／`decision`／`verdict`／`approved`／`grant`，逐条问四枚真路由（`:1983`），**塞进去必须不生效**。

### E3 结论：**哪些名字今天会被这枚冻结钉直接判死**

| 想新增的东西 | 会不会判死 | 现读依据 |
|---|---|---|
| **面板入向方法名** `grants.list` / `grants.revoke`（`docs/specs/SPEC-08:171` 的 C17 白名单原文） | ⛔ **判死**（如果实现成"Go 应答的入向方法"） | `carriesGrantWord("grants.list")` → `grantslist` 含 `grant`；`:1244-1246` 直接 Errorf。**这是文档与冻结钉的一处正面冲突**，且 SPEC-08 §5.2 题头自认【SPEC 提案，S5 定稿走契约批准】（`:156`）⇒ 属待人定案，不在本件裁定 |
| 面板入向方法名 `panel.grant.*`／`session.grant`／`panel.decide`／`grantScope` 之类 | ⛔ 判死 | 同上，且 `:1251-1262` 已把这些拼法写成显式候选 |
| **面板入向结构体的字段**：`Grant`、`Allow`、`Scope`… | `Grant`／`Allow` ⛔ 判死（`:1883-1885`）；**`Scope` 不在词表里 ⇒ 安全**（`sessionScope` 也不在） | `:1882-1886` 表一逐字 + `:1907` 判据 |
| 面板入向**新类型**（例如为会话答复建一枚 `SessionGrantRequest`） | ⛔ **不是判死，是判"注册表腐蚀"**：`inboundSeeds` 自动入册 → `:1855-1860 t.Fatalf`，而修 registry 要动冻结件 | `classifyDecodes:675-708`（`:700 pkg.inboundSeeds[d.TypeName] = true`） |
| **`internal/tools` / `internal/risk` / `internal/memory` / `internal/agent/**` / `cmd/wisp` 里的新导出名**（如 `GrantSource`、`SessionIDSource`、`MatchGrant`、`mintSessionID`、`SessionGrant`） | ✅ **不在射程内**（AST 只看 `internal/panel`：`:1231`、`:1848`） | 现读两处 `filepath.Join(root, "internal", "panel")` |
| **Go→前端出向字段**（卡片上要显示"本会话内已授权"） | ✅ 不判死（表一只跑 `inboundTypeRegistry` 的四型） | `:1720-1725` |
| **答复字段名** | ✅ **只有 `reason` 安全**（票 219 顶部第 8 格同结论，本件独立复现：12 词根里无 `reason`） | `:1883-1885` |

⚠ 还有一枚**非冻结但相邻**的形状守卫要避：`tools/d22scan` 的 ban #6 扫 `frontend/` 里的 `approval.decide`（票 219 顶部表格第 2 行记录的现读位置 `tools/d22scan/main.go:23-24`、`:825-831`）——**本件未现读 `tools/d22scan`，此条只作指路，不作依据**（见 Z 节）。

### E4 **"grant"这个词在本仓已经是二手的**（命名前的既成事实，现读）

`internal/agent/approval` 里 `grant` **早已被占用为"一次性原生 nonce"**，与 D45 的"会话授权"**同名不同物**：
- `internal/agent/approval/approval.go:218` 逐字标题 `// Grant: the unforgeable native-source proof (F2 layer 3)`；`:221-223` `A grant is a single-use nonce minted by the gate AT THE MOMENT a native…`
- `:234 const grantBytes = 32`、`:237 const grantPrefix = "grant_"`、`:239 func mintGrant()`、`:273 type grantStore struct`、`:281 newGrantStore`、`:284 issue`、`:292 spend`、`:309 revoke`、`:316 live`
- `internal/agent/approval/queue.go:288`、`:336`、`:358`、`:384` 全在用这枚 nonce；`gate.go:478 nonce, err := g.q.grantNonce(it)`；`gate.go:596 func (n nativeAPI) Allow(_ context.Context, corr, grant string)`
⇒ **在同一棵树里再造一枚也叫 grant 的"会话授权"，会让 E1 那四枚测试名（`…UnderAFreshGrant`／`…ItsOwnGrant`／`…BurnsTheGrant…`／`TestGrantIsSingleUseAndBoundToItsItem`）读起来全部歧义。**
歧义本身不判红（它们是测试名，不是键名也不是路由名），但它是**下一位读者误判的入口**，属实的风险，不是风格问题。〔已证（现读）〕

## F. 「什么算同一次会话」的候选定义（只摊事实与后果，不替 owner 选）

> ⚠ **本节刻意不排序、不推荐、不替 owner 回答。**票 224 自己写着这一问"属待人定案（未定义即停）"，
> 且 `AGENTS.md` §2 的「未定义即停」清单里"会话＝什么"这一条**尚未定案**。下面每一支的后果全部挂在 A–E 的现读行号上；
> **owner 拍哪一支，是 owner 的事。**
> 三条共用的现读前提（先摆明，免得被当成某一支的私有缺点）：
> ① 全仓今天**没有任何铸造点**（A8）；② **没有任何匹配器**（C0 第三行）；③ **`internal/perm/ticket90_persist_test.go` 三枚钉之外还有第四枚 `cmd/wisp/run_mode101_test.go:389`，四枚对"身份派生"这条轴全瞎**（D4、D5）。

### F-甲 「会话 ＝ 一次 `wisp run` 进程」

| 维度 | 现读事实与后果 |
|---|---|
| 铸造点会长在哪 | `cmd/wisp/run.go:321 func assembleRuntime`（它同时是 `rt.store = mem` 的落点 `:381` 与 `rt.bridge = tools.New(...)` 的落点 `:549`）——**这是全树唯一同时拿着 store 与桥的位置** |
| 结束点在哪 | 有，且由 OS 保证：`cmd/wisp/main.go:85 os.Exit(cmdRun(...))`；库侧 `cmd/wisp/run.go:212 defer rt.close()` → `:676-683`（只做 `replyRoot.Cancel()`＋`store.Close()`，**不 revoke**） |
| **稀释成永久通行证的风险** | **中等偏高，且风险全在"id 怎么来"这一枚旋钮上**：若铸的是每次启动新的随机值（如 `crypto/rand`），进程一死身份就不可复用——形状安全；**若从 pid / 启动时刻 / 数据目录路径派生**，Windows 的 pid 会复用、路径跨启动恒定、时刻可预测 ⇒ **同一条 `INSERT` 在下一进程照样被查回 ⇒ 正是 D5 那一形，四枚钉全绿**。`approval_grant` 的行**本来就跨重启留着**（`internal/perm/ticket90_persist_test.go:265-273` 明确要求留一行供审计）⇒ **库不干净，只有身份能挡住复用** |
| ⚠ 另一支要一并摊的事实 | 常驻腿**也存在**：`cmd/wisp/main.go:55-61`（无参数 → `runResident()`）＋ `cmd/wisp/resident_windows.go:24 func runResident()`。现读 `cmd/wisp/approval_reply.go:48` 逐字：「`cmd/wisp/resident_windows.go` runs an event loop with **no gate, no bridge and no** …」⇒ **今天常驻腿上一枚工具调用都不会发生**，所以"进程＝会话"在常驻腿上是**空值**；**一旦常驻腿接上桥（票 201/222/223 那族的方向），常驻进程的"会话"＝整个应用寿命 ⇒ 那一支上"本会话内允许"逐字变成"本次开机内永久免审"** |
| 要新增的契约面 | ① `agent.ToolRequest`（`internal/agent/tools.go:54`）要带会话身份，或桥侧新增一枚只读 GrantSource seam（照 `internal/tools/mode.go:27-30 ModeSource` 那形）——**两者都是新契约面**；② `internal/tools/bridge.go:288` 之后那一步的"命中"语义要写进 SPEC-06／`PLAN.md`（现文只写"绑(工具,路径模式,会话 ID)"`PLAN.md:2147`，**没写身份从哪来**）；③ `Silenced` 结构体（`internal/risk/mode.go:138-152`）缺"为什么没问（非 mode）"的槽位 |
| 可不可撤销 | **可**，两条路：进程死（自动，但依赖 id 不复用）；显式撤销要 `RevokeGrant`（`internal/memory/dao_misc.go:49`，今天零生产调用者）＋ C17 的 `grants.revoke`（`docs/specs/SPEC-08:171`）——⚠ **后者作为面板入向方法名会被 E3 那枚冻结钉判死**，作为原生侧 API 名则不在射程内 |

### F-乙 「会话 ＝ 宿主显式铸造、绑状态机边沿的一次交互（球/面板从召起到结束／90s 空闲）」

| 维度 | 现读事实与后果 |
|---|---|
| 文档已经站在这一支上 | `PLAN.md:2187` 逐字给了定义（`从唤起到用户显式结束，或空闲 90s（可配 [session] warm_timeout_sec）`）；`PLAN.md:1279`／`:1381` 把"D45 会话授权登记"分给 **C31 `SessionScope`**；`SPEC-08:92`／`:121` 把建/销钉在两条 D43 转移上；`SPEC-05:135-136` 说失效是"随 `session_id` 过期" |
| 但它的 owner 是空包 | `internal/session/` 里**只有 `doc.go`**（`ls` 现跑），`:16-17` 逐字 `DEFERRED(SessionScope): implemented by ticket 28. This ticket only freezes the package boundary.`；`doc.go:12-14` 又逐字声明它**不管状态转移**（`no state transitions (statemachine)`）而 `internal/statemachine/table.go:16` 反向声明转移钩子**不归它执行** ⇒ **"谁在边沿上建/销 SessionScope"今天两不管**〔仅文档写了、代码没有〕 |
| 边沿本身是真的 | 起：`internal/statemachine/table.go:51`（`Sleeping --EvSummon--> Listening`）、`:215`（`Warm --EvSummon--> Listening`）；止：`:227`（`Warm --EvWarmIdle--> Settling`）、`:231`/`:236`（`Settling --EvSettleExpired--> Sleeping`）。视觉侧已有同一批边沿的监听者可以照抄：`internal/ball/liquid_windows.go:148-167`（`:157`/`:165` 调 `endSession()`） |
| **稀释风险** | **最低**（若每次 Summon 都新铸、每次 Settling 都销毁）。⚠ 但有三条**独有的**漏法必须写进判据：① 从状态计数/计时器**派生** id ⇒ 可预测，等价于 F-甲 的派生坑；② `Warm` 跨空闲**不换 id** ⇒ 会话实际寿命由 `warm_timeout_sec` 决定，而那三枚配置键**今天零读者**（A6）⇒ 计时器不存在，结束点也就不存在；③ 超时若用墙钟时间差实现 ⇒ 直接撞 `AGENTS.md` §1.2 的禁止形状（"用墙钟时间差实现超时"），这一条本件**未现读仪器源码**，只作指路 |
| 要新增的契约面 | **最多**。① 一枚"会话身份从哪来"的契约面（票 224 已点名：C17／SPEC-02 §3 只写了 `session` 这个 scope 值，**没写身份来源**）；② `internal/session` 的真实现与 owner（C31，B4 的 BLOCKER）；③ 授权登记要从这里进桥（B3/B4）；④ 面板「安全」页的实时列表＋一键撤销（`PLAN.md:2149`、`SPEC-00:75`）⇒ **撞 E3 的方法名冲突** |
| 可不可撤销 | **最强**：显式结束（用户操作）＋ 90s 空闲（自动）＋ `RevokeGrant` 三路；且"撤销"有用户可见的语义（球回 `Sleeping`）。前提是有人真执行 Dispose——今天**零执行**（A5） |

### F-丙 「会话 ＝ 一段有界时长（`expires_at` 就是 TTL，不依赖身份）」

| 维度 | 现读事实与后果 |
|---|---|
| 它在文档里有一半支撑 | `internal/memory/schema.go:96` 注释把 `expires_at` 写成"**会话结束时间**（失效后仅作审计）"；`internal/memory/retention.go:219-227 func deleteExpiredGrants` 已经在按 `MAX(expires_at, COALESCE(revoked_at,0)) < cutoff` 删行（30 天审计窗，`SPEC-02:158`）⇒ **TTL 的机器是现成的**，只有 retention 那一条腿且不查活不查死只看墙上时间 |
| **稀释风险** | **三支里最高，且是"字面意义上的永久通行证"**：`PLAN.md:1642` 与 `SPEC-05:136` 要求的失效条件是"**会话结束后必须失效**／随 `session_id` 过期"，**不是**"到点失效"。⇒ 用户显式结束会话之后，那一行**仍然活着**直到 TTL 到；如果 id 又是派生/恒定的，则 TTL 只是"多久以内必问、多久以后又变免审"的**再循环**，而不是撤销。⚠ 另外 `session_id` 列是 `TEXT NOT NULL`（`schema.go:94`）＋ `InsertGrant` 硬性要求非空（`dao_misc.go:21-22`）⇒ **这一支仍必须铸一枚身份才能落库**，它并没有省掉铸造点 |
| 要新增的契约面 | **与冻结文字正面冲突**：把 `expires_at` 从"会话结束时间"改读成"TTL"要动 `internal/memory/schema.go:96` 与 `docs/specs/SPEC-02:99` 的文字——**`PLAN.md`／`docs/specs/**` 属禁区，须人工批准**（票 224 禁区第一条）。代码面反而最省（只补一枚 TTL 检查＋匹配器） |
| 可不可撤销 | **可撤销且撤销得最不干净**：`RevokeGrant`（`dao_misc.go:49`）能立刻打标，但"到点自动失效"与"会话结束即失效"是两件事，前者不提供"用户做了 X ⇒ 授权没了"的可验证因果；AC#3"重启必失效"在这一支**根本不成立**（重启不改变时钟）⇒ **AC#3 会直接判红，除非实现成"身份＋TTL 双条件"** |

### F-丁（附一支，因为它几乎是免费的）「会话 ＝ 一次任务（复用 `task_id`）」

| 维度 | 现读事实与后果 |
|---|---|
| 唯一现成的身份 | `internal/agent/loop.go:1123 newTaskID()`（UUIDv4 形），铸造点 `:322`（`RunAsync`）／`:333`（`Run`）；有真结束点：`internal/tools/bridge.go:749 CloseTask`，挂在 `cmd/wisp/run.go:733 AdmitTask: rt.admitTask` 上 |
| **稀释风险** | **最低**（每轮随机、每任务一个）。⚠ 但它**语义上不是"会话"**：`Run`/`RunAsync` 每次调用都铸新 id ⇒ **一枚"本会话内允许"会在一次用户回合结束时失效**，`PLAN.md:2147` 的"同类重复操作只确认一次"在多轮对话里几乎不成立（`SPEC-00:71` 场景 30 要的是整个会话内不再问） |
| 要新增的契约面 | 把授权从 `session_id` 改绑 `task_id` ⇒ **直接改 `PLAN.md:2147`／`:2704` 与 `SPEC-02:92-102` 的三元组形状**＝契约变更，人工批准；且 `CloseTask` 注释自陈仍有没 owner 的腿（`internal/tools/bridge.go:712-748`，`Q-56`）⇒ **有结束点的保证今天只对 loop.go 自己铸的 id 成立** |
| 可不可撤销 | 任务结束自动撤销（若 `CloseTask` 真被走到），但**用户不可见、不可主动撤销**（没有"关掉一个任务"的产品动作） |

⇒ **F 的诚实收法**：四支各自的"结束点"只有 F-甲（进程）与 F-乙（状态边沿）是**产品可感知**的；只有 F-丁 的身份**今天已经铸好**；
**没有任何一支能在不动契约面的前提下落地**——差别只在动哪一枚。本件到此为止，**不替 owner 选，也不替他答**。

## Z. 我没查清／查不动的

1. **`cmd/wisp/approval_reply.go` 的内容我不作依据**：它是编排者点名的"写腿正在改"的文件（现读 539 行、函数清单见 C1）。
   ⚠ 本件在 C1 给的该文件行号是**这颗 sha 的工作树快照**，交回时极可能已漂移；`git status` 现跑显示 `cmd/wisp/approval_reply.go` 与 `cmd/wisp/run.go` 均为 **M（未提交修改）**，`gate.go`／`queue.go` 的行号也与票 219 顶部记录不同（`DecideFromNative` 现为 `:626`、`DecideFromPanel` 现为 `:638`、`Veto` 现为 `:387`）。**引用前必须重跑。**
2. **没现读 `tools/d22scan`**：E3 末行关于 ban #6 扫 `frontend/` 的那句，出处我只从**票 219 顶部的记录**得知（它引 `tools/d22scan/main.go:23-24`、`:825-831`），我**没有亲自 `sed` 那两枚位置** ⇒ 那一格在本件里**只作指路、不作证据**；要当证据请另跑。
3. **`internal/panel` 的 4 枚 FAIL 我没核实**：票 219 顶部第 8 格写着「`internal/panel 95 PASS／4 FAIL`，含 `tokens_fourway_test.go`」——**硬约束禁止我跑测试**，所以这 4 枚红在本件里**是未核实的转述**，我只把文件名出现与否当作事实（`tokens_fourway_test.go` 我只在禁区清单里读到名字，未打开）。
4. **`docs/specs/SPEC-09` 里 per-session 互斥（单实例）没查**：`SPEC-02:214` 那句"单进程单实例语义由 per-session 互斥保证，SPEC-09 §5"里的 **per-session 指 Windows 登录会话**，与本票的授权会话是两回事；我没有现读 SPEC-09 §5 去确认互斥量的命名是 `Local\` 还是 `Global\`（`PLAN.md:3045` 逐字要求"单实例互斥量必须 per-session（`Local\`）而非 global"）。⇒ **这一条会影响 F-甲 的一个真实漏洞**：同一登录会话里两个 `wisp` 进程会不会共享同一枚派生 id，我**没查**。
5. **`wisp panel-inbound` / `panel_pump.go` 这两条腿我没读**：票 219 顶部记录说 `panel-inbound` 走的是"入向受理"而不是答复。本件只核了 `cmd/wisp/main.go:109` 的派发存在，**没读 `panel_inbound.go`／`panel_pump.go` 的内容** ⇒ "面板那侧今天有没有任何一条腿会把会话身份带进原生"这一问**未答**。
6. **音频/语音侧的"会话"我按 A5 只扫了命名，没扫语义**：`internal/audio/wasapimic_windows.go:117`（"last capture session"）、`internal/config/schema.go:195`/`:276`（recording session、`MicMutedDefault starts every session muted`）、`internal/llm/systemproxy_windows.go:16`/`:50`（Windows user session）都在 `session` 这个词下说的是**另一件事**。我只确认它们与授权无关，**没有逐枚读**；`internal/config/schema.go:536 keep_alive_in_session` 我读了字段声明，**没查它的读者**。
7. **`internal/panel/l2_grant_boundary_test.go` 我只读了五段**（`:1-40`、`:185-242`、`:670-715`、`:1229-1275`、`:1330-1460`、`:1715-1760`、`:1870-2000`），全文 2380 行**未通读**。⇒ E2/E3 的结论对**我读到的 facet 1/2/3 成立**；`:1530 TestRealGuardRefusesEveryAssemblableApprovalRouteName`、`:1595 TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` 的**完整断言**我只看了名字与相邻注释，没逐行。若写腿要照 E3 避名，**建议它自己再核这两枚**。
8. **`tool_call.grant_id` 的"谁该写它"我没找到成文出处**：`PLAN.md:2703` 只说"`grant_id` 关联 D45 的会话授权"，`internal/agent/journal.go:22` 有形参、`:101` 硬写 nil——**没有任何 spec 写明是哪一步负责填**。这一格属"未定义"，按 `AGENTS.md` §2 应停手问，不在本件裁定。
9. **没跑任何编译／测试／vet（硬约束 1）**；**没读、没引用 `frontend/**` 与 `design/**` 的任何内容（硬约束 2）**——包括 `internal/panel/l2_grant_boundary_test.go` 文件头里那些指向别家树的路径文字，我一律未转述。
10. **票 224 顶部那张表里"`internal/agent/journal.go` 的 `DecisionAllowGrant` 决策类型零写者"我复现了**（A2），但**同表那句"判定链只认最严"我用的是 `Mode.Screen`（`internal/risk/mode.go:157`）**，仓里**不存在 `resolveMode` 这枚符号**（B0）⇒ 那一格属票面的旧命名，**不是我的读数分歧**，写腿照 `Screen` 引即可。
