# 267-a3 普查：`confirm_timeout_sec` 生效面 vs「300 秒」说谎面

> 腿型：只读普查（coordinator 派单代号 267-a3）。产物＝名册与射程判断，**不裁 `Q-77`**、不选甲乙。

## §0 起手锚

- 起手时刻：`09:53:08+0800`（`date "+%H:%M:%S%z"` 自取）。
- 起手锚点：`git rev-parse --short HEAD` = `c9600334`（自取，非抄派单）。
- 分支：`dev`，共享工作树。并发可见：同票另有 `267-a1`／`267-a2`／`267-r1`／`267-gate` 证据件已落；`git status --porcelain` 起手即见大量他腿改动（含 `design/**`、`.gitignore`、`.scratch/wisp/probes/151/**` 等），本腿**不 add 任何他枚文件**。
- 本腿禁跑清单（起手自约）：零 `go build`／`go test`／`go vet`／`gofumpt`／`wisp slo`／`sh scripts/*`；不读 `frontend/**`、`design/**`；不碰 `.scratch/wisp/issues/**` 的 AC 框；不碰 `docs/reports/pending-and-issues.md`；不改任何已跟踪文件（含冻结件 `docs/PLAN.md`、`docs/specs/**`）。
- 派单给的两枚消费点（`cmd/wisp/run.go:616`、`cmd/wisp/resident_approval_windows.go:460`）**本腿一律现量**，不抄行号。

## §1 用户会看见的字：写死「300 秒／5 分钟」的文案名册

> 尺：`"[^"]*秒[^"]*"`、`"[^"]*300[^"]*"`、`"[^"]*(分钟|minute)[^"]*"`（根＝`cmd`／`internal`／`tools`，正则要求秒数落在**双引号字符串内**＝产码可见文案，注释不计）；逐枚再 `Read` 现量喂值表达式。取数窗口＝09:53:45+0800 至本节 commit 时刻（逐尺行见 §6）。
> 判语：**Go 侧生产文案里零枚字面「300 秒」/「5 分钟」。** 两枚带秒数的句子都是 `fmt.Sprintf` 现算，数字来源就是队列自己的 `timeout`。

| # | file:line（现量） | 逐字文案 | 数字从哪来 | 用户配 `confirm_timeout_sec = 90` ⇒ 这枚会不会说谎 |
|---|---|---|---|---|
| 1 | `internal/agent/approval/queue.go:479` | `审批超时（%d 秒未确认），C18 一律判拒绝，已自动拒绝` | `int(q.timeout.Seconds())`（`queue.go:64` 字段，`NewQueue` 由调用方喂） | **不说谎**：配 90 后显示「90 秒未确认」。⚠ 前提是这一枚确实是队列实例的 timeout；`wisp run` 路上 `run.go:616` 喂的是 config 值 ⇒ 一致 |
| 2 | `internal/agent/approval/gate.go:563` | `审批将在 %d 秒后自动拒绝，请尽快确认` | `int(g.q.WarningLead().Seconds())`＝**提示提前量**（默认 `DefaultApprovalWarning = 30s`，`queue.go:109`），与 timeout 不是同一枚数 | **不说谎**，但射程要说清：这句永远讲「30」，因为它量的不是卡片开多久。若机主将来把 warn 也接进 config，这枚才会跟着动 |
| 3 | `cmd/wisp/run.go:1364` | `\n[确认 %s %s] %s`（Level／Tool／Reason） | 无秒数 | **不说谎＝但它是沉默面**：CLI 从不开口说这张卡会开多久 |
| 4 | `cmd/wisp/run.go:1394` | `[%s] %s\n`（`e.Kind` ＋ `e.Text`） | 转抄 #1／#2 的现算文本 | 不说谎（跟着 #1/#2 真值走） |
| 5 | `cmd/wisp/resident_approval_windows.go:829` | `wisp: 卡片挂起：%s %s（编号 %s）` | 无秒数 | 同 #3：常驻腿的 stdout 也是沉默面 |
| 6 | `cmd/wisp/resident_approval_windows.go:746` | `审批门已装配进本进程（取消通道：%s 已加载；等待中的确认项：%d）` | 无秒数（`%d`＝队列深度） | 不说谎 |
| 7 | `internal/tools/bridge.go:467` | `审批超时未确认，已自动拒绝`（`orDefault(why, …)`＝`why` 为空的兜底句） | 无秒数 | 不说谎；形状是「不含数字的兜底」，配任何值都不会与它矛盾 |
| 8 | `internal/agent/prompt.go:63` | `不可逆或越界操作必须先经用户确认，确认超时视为拒绝。`（进模型提示词，`prompt_test.go:268` 逐字钉 needle） | 无秒数 | 不说谎（这把尺不量秒数） |
| 9 | `internal/config/unwired.go:121` | `consumed: cmd/wisp/run.go and cmd/wisp/resident_approval_windows.go both build the approval timeout from it (the resident leg since ticket 256); ticket 267 bands it [31, 3600] at load` | 无秒数，但**断言「consumed」** | ★**判甲方向（"配了但不真生效"）唯一一枚会当场变成假话的现成文案**：它今天声称两枚消费点都从它建超时。反方向（判乙＝真生效）会说谎的不是文案，而是 §3 L3（日志字面 `DefaultApprovalTimeout=300s`）、§3 L5／§2.3（事件名 `approval.timeout-300s` 与状态机表里那枚 300s）。`Q-77` 若判甲，说谎的就是这一行；判乙则它为真。名册归本腿，**裁定归机主** |
| 10 | `cmd/wisp/config_reload.go:224`／`:229` | 放宽未确认／已确认两条款回执 | 无秒数 | 不说谎 |
| 11 | `internal/panel/pump.go:50-53`（注释自陈） | 逐字：`the no-source list is seven: approval.remainingMs joins it, because the L1 countdown has no producer either - approval.EventTick is declared and emitted nowhere in the tree, the Remaining that is emitted carries the static window length` | 面板快照**没有**审批剩余秒数这一栏 | **不说谎＝也是沉默面**：快照里根本没有 `remainingMs`，所以配 90 也不会有一栏写 300。本腿用第二把尺独立复认：`json:"(timeout\|window\|deadline\|remaining\|expire\|warn)[^"]*"` 在 `{cmd,internal}` **零命中** ⇒ 两把尺同向。（gate.go:274/:387/:510-513 那三枚行号是**该注释引用他票的读数**，本腿未复测，见 §7 G-4） |

**用户可见文案的说谎数＝0 枚（Go 侧），另 1 枚**条件性**说谎（#9，方向＝判甲）。** 但「不说谎」不等于「说全」：#3／#5 两片把「卡片还会等多久」这件事**完全不打印**，机主把 300 改成 90 之后，屏幕上既不会有假的 300，也不会有真的 90——这一格是**沉默**，不是谎。

⚠ **审批卡／托盘／悬浮球／面板 HTML 上的字（「5 分钟」这类人话最可能落点）属页面侧：`frontend/**` 与 `design/**` 两棵树本编队不许读，故"页面上有没有写死的 300 秒／5 分钟"这一格本腿量不到，需机主带去他用的那枚 agent 另查。** 派单点名的四类可见面里，球与托盘的**文字渲染**同样落在那两棵树（Go 侧只 `SetState` 一个状态名，`ball_windows.go:309-311` 注释逐字「the consumer drives this after its statemachine dispatch」＝球不自带计时器）。

## §2 常量与默认值的落点

### 2.1 `DefaultApprovalTimeout` 定义与引用者（现量尺＝`DefaultApprovalTimeout`，根＝`{cmd,internal,tools,scripts}`）

- 定义：`internal/agent/approval/queue.go:107` 逐字 `DefaultApprovalTimeout = 300 * time.Second`（同段 `:106` 注释逐字 `// DefaultApprovalTimeout: 「超时 300s 一律判拒绝」. Never an infinite wait.`，注释豁免 d22scan 的 emoji 尺，但它是**人读的 300**）。
- **生产引用者＝1 枚**：`internal/agent/approval/queue.go:86`（`NewQueue` 的 `timeout <= 0` 兜底；兜底**不看 config**，只看调用方给没给值）。
- 生产侧"非引用者"的 300 落点（**字面串**，不是常量）：`cmd/wisp/resident_approval_windows.go:456`（日志，见 §3）、`:445`／`:296`／`:308`／`:309`、`cmd/wisp/resident_windows.go:129`、`cmd/wisp/config_reload.go:50`、`internal/config/manager.go:26`、`internal/tools/gate.go:86`、`internal/perm/store.go:86`、`internal/risk/mode.go:32`、`internal/agent/approval/{doc.go:6,21; gate.go:30,499; approval.go:21; queue.go:59,82}`、`internal/panel/composer_handlers.go:23`＝**注释与包文档**，不进屏幕。
- 测试引用者（**与生产分开**）：`cmd/wisp/resident_approval_risk_256_windows_test.go:141,144,214,229,237`；`internal/agent/approval/ticket84_no_owner_test.go:105,108,110,113,114,232,233,253,254,256,257`；`internal/agent/approval/queue_test.go:48`。

### 2.2 schema 默认值

- `internal/config/schema.go:460` 逐字 `ConfirmTimeoutSec int \`toml:"confirm_timeout_sec" default:"300"\``（注释 `:449` 「ConfirmTimeoutSec is how long an L2 confirmation card stays open.」）。
- 落盘路径：`wisp` 首启 `NewDefaults()` 反射该 tag ⇒ 今天盘上的 `config.toml` 写 300（267-a2 现量 `firstrun.go:82` 那条路，本腿不重跑，引它并注明是**他腿读数**）。
- 值域门：`internal/config/validate.go:128-129` 逐字 `confirmTimeoutSecMin = 31`／`confirmTimeoutSecMax = 3600`；报错文案 `:147` 逐字 `"config.toml: risk.confirm_timeout_sec %d out of range [%d, %d]"`——**动态喂值，不说谎**。

### 2.3 ★第二处写死 300 的**行为**面（不是文案，是本腿新钉的一格）

- `internal/statemachine/timeouts.go:51` 逐字 `{state: StateAwaitingApproval}: {300 * time.Second, EvApprovalTimeout},`——状态机自带的**每状态超时表**里第二枚硬编码 300s，`DefaultTimeouts()`（`:43`）在 `Options.Timeouts == nil` 时被 `machine.go:70` 采用。
- `internal/statemachine/events.go:42` 逐字 `EvApprovalTimeout Event = "approval.timeout-300s"   // #24`——**事件 ID 的字面里嵌着 300s**。这条字符串会出现在 `machine.go:139-140` 的 `slog.Error("state machine rejected transition", … "event", string(ev) …)` 里，即：**任何一处非法转移的日志会把「300s」打进屏幕/台账**。
- 射程（现量尺＝`EvApprovalTimeout`，根＝`{cmd,internal,tools}`）：命中只有 `timeouts.go:51`、`table.go:172`（D43 第 24 行的转移定义）、`table_test.go:181`、`machine_test.go:24`——**零枚生产者手工 Dispatch 它**，只有那张表的定时器会发它。
- 本腿能判的：**球侧不自带机器定时器**（`internal/ball/ball_windows.go:309-311` 的 `SetState` 只换渲染状态；`internal/ball/doc.go:19` 逐字「no business logic: state transitions live in statemachine」）。本腿**量不到**的是「常驻/run 两条路里有没有哪条真的把 Machine 推进 `AwaitingApproval` 并 arm 这枚 300s」——见 §7 G-2。
- 为什么这一格重要：若机主判「配置真生效」，队列在 90s 关卡，而状态机表仍写着 300s ⇒ **同一件事有两枚数**，`approval.timeout-300s` 这枚事件名从此是假话（哪怕它没人发）。

### 2.4 C18 原文〔引自冻结件原文，本腿未改一字〕

- `docs/PLAN.md:1368`（契约表 C18 行，节选逐字）：`**超时 = 300s，一律判拒绝**（第四轮补数值，原文无数值 → agent 会自己拍一个或写成无限等待，见 §16.9 第 4 条）；**超时前 30s 醒目提示**；拒绝后**任务 root ctx 不取消，可一键重放**`。
- `docs/PLAN.md:3210`（§16 差异表 C18 行，节选逐字）：`| **C18** | 「超时一律判拒绝」**无数值** | **300s**；超时前 30s 醒目提示；…`。
- ⇒ **C18 把 300s 写成契约数值本身**，这正是 `Q-77` 不能被任何腿自行裁的原因：改「配了才生效」＝改 C18 的射程，属人工批准（`SPEC-12 §4.1`）。

### 2.5 消费点现量（不抄派单行号）

- `cmd/wisp/run.go:616` 逐字 `ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second,`（派单写的 616＝**仍对**，本腿 `Read` 复认；`run.go` 的 `grep -n ApprovalTimeout` 只这一枚）。
- 常驻腿**两枚落点**（派单只点了 460）：`cmd/wisp/resident_approval_windows.go:459-460` 逐字 `return time.Duration(c.Risk.L1WindowSec) * time.Second,` / `time.Duration(c.Risk.ConfirmTimeoutSec) * time.Second,`（在 `residentRiskGateValues` `:447` 内）⇒ 真正喂进接缝的是 `:372` 逐字 `ApprovalTimeout: timeout,`。
- 带外两条出口（`:448-449` dataDir 空、`:452-457` `config.LoadFile` 失败）都返回 `0, 0` ⇒ 落 `NewQueue` 的 `timeout <= 0` 兜底＝编译 300s。**267 的值域门把「带外配置」变成了这里的 `err` 分支**（`:453` `if err != nil || c == nil`），于是那条路既回落 300s 又打 §3 那行日志。
- 零枚消费点的反面也现量了：`grep -n "Risk\." cmd/wisp/*_test.go` 之外，全仓非测试对 `ConfirmTimeoutSec` 的读取**只有上述两枚**（尺＝`[Cc]onfirm[_]?[Tt]imeout`，根＝`{cmd,internal,tools,docs,scripts}`，命中行里生产者／消费者／断言者分列见 §5）。

## §3 日志与审计里带秒数的形状

> 尺：`logf\(|slog\.|auditf`（根＝`internal/agent/approval`，26 枚命中全读）＋`"[^"]*300[^"]*"`（§2 那把尺的日志子集）。逐枚再 `Read` 现量喂值。取数窗口＝09:53:45–10:01:38+0800（本节末行有 commit 时刻）。

| # | file:line | 逐字形状 | 秒数来源 | 配 90 ⇒ 会不会说谎 |
|---|---|---|---|---|
| L1 | `internal/agent/approval/queue.go:489-490` | `approval: ANSWER-EXPIRED corr=%s tool=%s decision=timeout->reject after=%ds` | 尾参 `int(q.timeout.Seconds())`（`:490`） | **真话**：配 90 就打 `after=90s`。审计行与屏幕文案（§1 #1）同源同值 ⇒ 两面对得上 |
| L2 | `internal/agent/approval/gate.go:338` | `approval: ANSWER-EXPIRED corr=%s tool=%s decision=timeout->execute after=%s` | L1 窗口 `Duration`（`%s`，非秒数整数） | 真话；且这枚是 **L1 窗口的 timeout→execute**，与 L2 的 reject 共用 `ANSWER-EXPIRED` 词头——审计里区分两者靠 `decision=` 那一对值，不靠数字（这是射程事实，不是缺陷判语） |
| L3 | `cmd/wisp/resident_approval_windows.go:454-456` | `slog.Warn("resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants", "path", cfgPath, "err", err, "fallback", "DefaultApprovalTimeout=300s / DefaultL1Window=3s")` | ★**日志里唯一一枚字面「300s」**，硬编码在 attr 值里 | 触发条件（`:448` dataDir 空 / `:453` `LoadFile` 返错）下门**确实**回落编译 300s ⇒ 今天为真。**两枚说谎风险**：① 兜底值一旦改动（或改接 config），这行字面立刻成假话，而它不在任何常量的引用链上＝**没有尺会因为它变红**（`queue.go:107` 改了它不知道）；② 267 的值域门让「用户配带外值」也走 `:453` 的 err 分支 ⇒ 日志只说"回落到编译常量"，**不说"你写的那个数被拒了"**；对结果是实话，对原因是残缺的实话。要补哪一句＝归机主，不归本腿 |
| L4 | `cmd/wisp/resident_approval_windows.go:382-383` | `"window_sec_read", window.Seconds(), "confirm_timeout_sec_read", timeout.Seconds()` ＋ `"gate_window", ra.gate.Window().String(), "gate_queue_timeout", ra.gate.Queue().Timeout().String()` | 现读的 config 值／现读的队列值 | 真话，而且是**唯一一枚能把"配的值"与"门实际拿到的值"并排打出来的日志**——`Q-77` 无论判甲判乙，机主要核验都得看这两行 |
| L5 | `internal/statemachine/machine.go:139-140` | `slog.Error("state machine rejected transition", "from", string(from), "event", string(ev), "err", err.Error())` | `string(ev)`；当事件是 `EvApprovalTimeout` 时字面打出 `approval.timeout-300s`（`events.go:42`） | ★日志面的**第二枚字面 300**：它把"300s"写进了**事件的名字**里。若机主判"配置生效"，那枚定时器/事件名的 300s 就成了假话；而它同时是 D43 冻结转移表第 24 行的事件身份（`table.go:172`）＝**改名＝改契约**，本腿只登记，不动 |
| L6 | `cmd/wisp/resident_approval_windows.go:824-828` | `slog.Info("approval: 常驻进程显示一张确认卡片", "corr", …, "level", …, "tool", …, "orb_state", …, "esc_borrowed", …, "channels", …)` | 零秒数 | 不谎；同样是沉默面（举卡时不宣告截止） |
| L7 | `internal/agent/approval/gate.go:567` | `approval: warning delivery failed corr=%s: %v` | 零秒数 | 不谎 |

审计面小结：**没有一枚审计行会把「300」写死**（L1 现算、L2 现算），字面 300 只活在 L3（回落日志的 attr 串）与 L5（事件名）两处。

## §4 文档（含冻结件标注）

> 〔冻结〕＝`docs/PLAN.md` 与 `docs/specs/**`：**只读判射程，一字未改**。

| # | file:line | 逐字（节选） | 标注 |
|---|---|---|---|
| D1 | `docs/PLAN.md:1368` | `**超时 = 300s，一律判拒绝**…；**超时前 30s 醒目提示**`（C18 契约行） | 〔冻结〕C18 本体 |
| D2 | `docs/PLAN.md:3210` | `\| **C18** \| 「超时一律判拒绝」**无数值** \| **300s**；超时前 30s 醒目提示；…`（§16 差异表） | 〔冻结〕C18 修订记录 |
| D3 | `docs/PLAN.md:2738` | `\| **\`[risk]\`** \| \`confirm_timeout_sec\`(300) \`l1_window_sec\`(2) …` \| 🔒 **变更即需 L2 级重新确认** \|`（D36 配置树） | 〔冻结〕——**注意它把 300 写成"默认值"而不是"固定值"**，这一枚是机主可引的原文弹药（判乙的文档依据） |
| D4 | `docs/specs/SPEC-06-security-gatekeeping.md:104` | `- **超时 300s 一律判拒绝**；超时前 30s 醒目提示；拒绝后任务 root ctx 不取消、可一键重放。` | 〔冻结〕门控规格 |
| D5 | `docs/specs/SPEC-06-security-gatekeeping.md:151` | `…队头单显、超时判拒绝、重放可续` | 〔冻结〕**不含秒数**⇒ 这一行永远不会因配置生效而说谎（名册里唯一一枚"无数字"的门控验收句） |
| D6 | `docs/specs/SPEC-03-config-secrets-envs.md:34` | `\| 🔒 \`[risk]\` \| \`confirm_timeout_sec(int)=300\` \`l1_window_sec(int)=2\` …` | 〔冻结〕配置规格——写的是 `=300`（默认值形状），与 D3 同族 |
| D7 | `docs/specs/SPEC-00-product-overview.md:74` | `33. 作为用户，我想让**审批超时一律判拒绝**（300s）且任务可一键重放…（C18）` | 〔冻结〕场景矩阵行——**这一枚最像"对用户承诺的字"**：场景文案把 300s 写成用户想要的东西 |
| D8 | `docs/specs/SPEC-08-ui-ball-panel.md:112` | `\| 24 \| \`AwaitingApproval\` \| 队头被拒绝/**超时 300s** \| \`Acting\`（取消该调用） \| 超时前 30s 醒目提示；可一键重放 \|` | 〔冻结〕UI 状态表（D43 第 24 行的 UI 视图）＝§2.3 那枚 `timeouts.go:51` 的**文档来源**：机器表与这行同源，改配置不改这张表的话两边都会讲 300 |
| D9 | `docs/specs/SPEC-02-data-storage.md:162` | `**进程启动后延迟 5 分钟跑一次 + 每 24h 一次**` | 〔冻结〕★**假阳性登记**：这是 `RetentionJob` 的 5 分钟，与审批超时无关。本腿把它列出来是因为同名尺（`5 ?分钟`）会命中它，别让它被当成第二枚说谎面 |

叙述件（非冻结，机主会读，本腿只列不改）：`docs/reports/missing-features-2026-09-29-v4.md:194`（`confirm_timeout_sec` 300 一行）；`docs/evidence/s1/246-resident-task-source-v2.md:112/:212`、`docs/evidence/s1/248-settings-write-path-r1.md:301-305`、`docs/evidence/s1/90-adversarial-acceptance.md:272`、`docs/evidence/s1/128-ac4-r1-acceptance.md:425`（引了一句红测试的原文，内含「审批超时（300 秒未确认）」）、`docs/evidence/s1/84-ac1-bounded-wait.md:25`。台账 `docs/reports/pending-and-issues.md:11761`（`Q-77` 本体）／`:11829`／`:10107`——**⛔ 本腿不碰台账**（`:11932` 那种派单记述同理只读）。
文档面判语：**把 300 写成数字的冻结件共 7 枚（D1／D2／D3／D4／D6／D7／D8）——其中 2 枚是"默认值"形状（D3 `\`confirm_timeout_sec\`(300)`、D6 `\`confirm_timeout_sec(int)=300\``），5 枚是"契约值"形状（D1／D2／D4／D7／D8）**。`Q-77` 判乙（配置真生效）之后，"契约值"那 5 枚并不会立刻假——300 仍是**默认**；真正会被变成假话的是把它们读成"固定值"的那类表述，其中**唯一一枚把 300s 写成用户愿望的是 D7（SPEC-00 场景 33）**。这一格判语归机主。

## §5 会因配置生效而红的测试钉

> 尺＝`DefaultApprovalTimeout`／`"[^"]*300[^"]*"`／`秒未确认|审批超时|确认超时|自动拒绝|已超时`／`confirm_timeout_sec`（含 `-C` 上下文读断言体），根＝`{cmd,internal,tools}`；生产 vs `_test.go` 分列。取数窗口＝09:53:45–10:03:50+0800。
> 派单担心的那一类（"断言提示文案里含 300"或"超时提示逐字等于某句"）——**本腿判它不存在**：全仓测试面里 `秒未确认` **零命中**，`"[^"]*300[^"]*"` 的测试行**全部是 `t.Errorf`／`t.Fatalf` 的解释语，不是被断言的 needle**。所以这一节按三种"红条件"分开登记，别把它们混成一锅。

### 5.1 A 类＝字面钉住「编译 300／schema 默认 300」的钉子（★改 `queue.go:107` 或 `schema.go:460` 才红；**用户配 90 不红**）

| file:line | 断言逐字（现量） | 红条件 | 用户配 90 ⇒ 红不红 |
|---|---|---|---|
| `internal/agent/approval/queue_test.go:48-49` | `if p.Deadline != approval.DefaultApprovalTimeout { t.Errorf("Deadline=%v，期望 300s", p.Deadline) }` | 常量动 | **不红**（该测试不读 config，量的是 `NewQueue` 兜底） |
| `internal/agent/approval/batch_test.go:106` | `t.Errorf("L2 卡片要挂 C18 的 300s 截止，得 %v", p.Deadline)` | 同上（断言体本腿未读到，见 §7 G-5） | 不红（同上形状） |
| `internal/agent/approval/ticket84_no_owner_test.go:105-114` | `{"explicit default", approval.DefaultApprovalTimeout}` ＋ `if approval.DefaultApprovalTimeout <= 0 { t.Fatal("DefaultApprovalTimeout must be a finite positive bound") }` | 常量被改成 `<=0`／无限 | 不红 |
| `internal/agent/approval/ticket84_no_owner_test.go:232-233,253-257` | `if got := noUI.Queue().Timeout(); got != approval.DefaultApprovalTimeout {` ＋ `if bel < approval.DefaultApprovalTimeout`／`if bel > 2*approval.DefaultApprovalTimeout` | ★**墙钟**钉：默认上界动它就红（跑时还要 `WISP_84_MEASURE=1`，`:224` 逐字 Skip 句「有意慢：300s 墙钟计量…」） | 不红（同样不经 config）；**但它是"300s 要改"这件事最贵的一枚代价钉** |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:141-144` | `if got := ra.gate.Queue().Timeout(); got != approval.DefaultApprovalTimeout { … "(300s is the contract default; a fourth number here means the fallback grew a value of its own)" }` | ★**这一枚就是"§3 L3 那句字面 300"的守门钉的反面**：兜底若长出第四个数就红；但 `resident_approval_windows.go:456` 那句**字面串不在它的射程里**（它测 `Queue().Timeout()`，不测日志文案）⇒ 改了 `:456` 的字面，**没有尺会红** | 不红 |
| `cmd/wisp/resident_approval_risk_256_windows_test.go:214,229,237` | `wantTime: approval.DefaultApprovalTimeout`（种子分别是 `l1_window_sec = 2`／`= 99`／`= 1`，都没写 `confirm_timeout_sec`） | ⇒ 它们断言的是等式「**schema 默认 300 ＝＝ 编译常量 300**」；`default:"300"` 或常量任一改动都红 | 不红（各自的 `t.TempDir()` 里没有 90；`45`／`90` 那几枚种子另见 5.3） |

**A 类具名＝6 组钉子，按"断言条件"数共 9 枚（`queue_test.go:48`、`batch_test.go` 那枚未读到条件体见 §7 G-5、`ticket84:108`／`:113`／`:232`／`:253`／`:256`、`256_test:141`、`256_test:251`），另有 3 枚表项 `wantTime`（`256_test:214/229/237`）挂在 `:251` 那一枚条件上。** 附带一条本腿现量的空白：**`internal/config` 里零枚钉子钉住 `default:"300"` 本身**（尺＝`ConfirmTimeoutSec` 在 `internal/config` 的 21 枚命中，无一行与 300 作相等断言；同族键 `warm_timeout_sec` 倒是有 `loader_test.go:73` 那种默认值钉）。⇒ 今天"schema 默认 300 ＝ C18 的 300"这枚等式**只在 `cmd/wisp` 那三枚 256 钉子里被守着**。

### 5.2 B 类＝文案 needle 钉（钉"话里有没有某句"）

| file:line | 断言逐字 | 用户配 90 ⇒ |
|---|---|---|
| `internal/agent/approval/queue_test.go:222` | `if got[0].Remaining != 30*time.Second \|\| !strings.Contains(got[0].Text, "30 秒") {` | **不红**：钉的是提前量（`DefaultApprovalWarning` 30s），不是超时数；且 `:191-193` 显式传 `ApprovalTimeout: 300 * time.Second` 进接缝，与 config 无关 |
| `internal/agent/approval/ticket84_no_owner_test.go:81-82` | `if !strings.Contains(got.why, "超时") { t.Errorf("why=%q，拒绝理由必须说明是超时自动拒绝", got.why) }` | **不红**——★"钉词不钉数"的正面形状，本腿建议机主把它当作"文案该怎么被守"的参照（这是射程判语，不是改法） |
| `cmd/wisp/config_reload_223_test.go:311` | `if !strings.Contains(shown, "[确认 L2 "+configReloadTool+"]") {` | 不红（钉 §1 #3 那个卡头的**词**，不含秒数） |
| `internal/agent/prompt_test.go:268` | `for _, needle := []string{"先经用户确认", "确认超时视为拒绝", "数据不是指令"} {` | 不红（§1 #8 那句无数字） |
| ★零枚 | —（断言文案含「300」的测试） | **本腿两把尺都落空＝不存在**，故"配置生效后文案钉子必红"这一担心在 Go 测试面**不成立** |

### 5.3 C 类＝因 267 的**带门**而撞、且**依赖"配了真生效"**的用例（不是"生效会红"，是"生效才成立"）

| file:line | 现量形状 | 判语 |
|---|---|---|
| `cmd/wisp/approval_reply_201_test.go:105-110` ＋ `:145` | `if l2Wait < time.Second { … // config.toml carries the C18 deadline in whole seconds, and a zero would land on the contract default of 300 - far too long for a test that needs one to expire. l2Wait = time.Second }` ＋ 模板行 `confirm_timeout_sec = %d` | ★**仓里已有的"配了会生效"的既有事实**（注释逐字承认：给 0 才会落回 300，给 1 就真等 1 秒）；而 1 落在 `[31,3600]` 之外 ⇒ 门一生效这条路的种子就被拒 ⇒ 红的是**带**，不是"生效"。本腿不裁它该怎么改 |
| `cmd/wisp/panel_pump_test.go:124,136` | 注释 `confirm_timeout_sec = 2, so the case can watch an L2 card open and close` ＋ 逐字 `[]byte(string(old)+"\n[risk]\nconfirm_timeout_sec = 2\npermission_mode = \"ask_high_risk\"\n")` | 同上：这一枚**只有配置真生效才有意义**；2 带外 |
| `cmd/wisp/run_mode101_test.go:110` | `riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"` | 1 带外；`docs/evidence/s1/128-ac4-r1-acceptance.md:425` 记过它红句「配置里写着 confirm_timeout_sec = 1 却被默认值取代」＝**同一格的历史读数，本腿不重跑** |
| 另 4 枚（**他腿读数，本腿未复量**） | `approval_seam_201_test.go:50/:138`、`ticket224_assembly_test.go:81/:234` 各喂 20 | 出处＝同票 `267-a2` 的 §0／§1／§2 普查（本腿只引不抄，其 commit 见 `git log`：`89976ed1`／`c9600334`） |
| 带内、门不撞（具名以免被误登记） | `resident_approval_risk_256_windows_test.go:159/202/221/286`（45）、`:293/:325`（90）；`unwired_test.go:66`（300）、`:188`（60）；`boundary_test.go:142`（60，断言体见 §7 G-5）；`validate_test.go:236`（31）／`:240`（30，**这条就是门的钉子本身**） | 全在 `[31,3600]` 内或本就是拒带外的正向钉 ⇒ ★**「配 90 真生效」这一格今天已被 `256_test:284-327` 钉住，但钉的是"重建后才生效"那一半**：`:286` 种 45 建门 → `:293` 把文件改成 90 → `:310` 逐字要求**已建好的门仍是 45s**（红句 `"the live gate moved from its construction value (%v) to %v after config.toml changed."`＝**构造期取值、不热生效**）→ `:324-325` 再要求重建的新门读到 90s。同段 `:311-314` 逐字承认：若将来给 `[risk]` 加了 re-apply 路径，这一枚就红，「ticket 256's stated limitation is obsolete and has to be re-adjudicated, and the 255/265 wording for the ⓑ sentence with it」。⇒ **生效时机（重启／重建 vs 热生效）是 `Q-77` 的第三维**，本腿在 §9 单独立一条给机主；本腿不裁 |

**§5 判语：因"用户配 90 且真生效"而必红的测试钉＝0 枚（Go 面）。** 会红的是另外两类：改那个**字面 300**（A 类 9 枚落点）与 267 的**带门**把带外种子挡住（C 类 4 枚具名＋他腿 4 枚）。⇒ **"配了要不要真生效"这件事不会撞上任何现有测试**；它撞上的只有 §4 那 7 枚冻结文档文字与 §2.3 那枚状态机表/事件名。

## §6 尺读数与未跑清单（门禁读数）

> 每把尺都带取数窗口；计数尺一律用检索工具的 count/content 模式（⛔ 未用 shell grep，避免 `grep -c` 零命中断 `&&` 链）。

| # | 尺（正则） | 根 | 读数 | 取数窗口 |
|---|---|---|---|---|
| R1 | `DefaultApprovalTimeout` | `{cmd,internal,tools,scripts}/**` | 22 行＝定义 1（`queue.go:107`）＋注释 2（`queue.go:106`、`resident_approval_windows.go:445`）＋生产引用 1（`queue.go:86`）＋日志字面 1（`:456`）＋测试 17（256_test 5／ticket84 11／queue_test 1） | 09:54+08 |
| R2 | `300 ?秒\|300秒\|5 ?分钟\|五分钟\|5min\|5 min\|300s` | `cmd` | 9 枚文件／15 行；**产码可见文案 0 枚**（唯一非注释落点＝`resident_approval_windows.go:456` 日志 attr） | 09:55+08 |
| R3 | 同 R2 | `internal` | 32 行；产码可见文案 0 枚（其余为注释／包文档／测试解释语） | 09:55+08 |
| R4 | `秒未确认\|审批超时\|确认超时\|自动拒绝\|已超时` | `{cmd,internal,tools}` | 10 行；产码 4 枚（`bridge.go:467`／`prompt.go:63`／`gate.go:563`／`queue.go:479`）＋测试 6 枚 | 09:56+08 |
| R5 | `"[^"]*秒[^"]*"`（字符串内的「秒」） | `{cmd,internal,tools}` | **3 行**＝`queue_test.go:222`、`queue.go:479`、`gate.go:563` ⇒ 产码秒数文案**只两枚**，且都 `Sprintf` 现算 | 09:56+08 |
| R6 | `"[^"]*300[^"]*"`（字符串内的「300」） | `{cmd,internal,tools}` | 24 行；生产可见串**零枚字面 300**；非注释落点＝`events.go:42`（事件名）／`resident_approval_windows.go:456`（日志 attr）／`schema.go:460`（toml tag） | 09:57+08 |
| R7 | `"[^"]*(分钟\|minute)[^"]*"` | `{cmd,internal,tools}` | 2 行（`cmd/llmrecord/main.go:46` 的 flag 默认 2 分钟、`loop_golden_test.go:22` 的语料「提醒我五分钟后关火」）⇒ **零枚"5 分钟"审批文案** | 09:55+08 |
| R8 | `[Cc]onfirm[_]?[Tt]imeout` | `{cmd,internal,tools,docs,scripts}` | **26 枚文件／98 行**＝`cmd` 21／`internal` 53／文档（PLAN＋SPEC＋reports）15／叙述件（`docs/evidence/**`）9；count 尺现量 | 10:05+08 |
| R9 | `json:"(timeout\|window\|deadline\|remaining\|expire\|warn)[^"]*"` | `{cmd,internal}` | **0 命中**（⇒ 面板快照无超时秒数栏，与 `pump.go:50-53` 注释同向） | 09:58+08 |
| R10 | `logf\(\|slog\.\|auditf` | `internal/agent/approval` | 26 行（含 3 行 `t.Logf`）⇒ 审计行带秒数的只 `queue.go:489` 一枚（现算） | 10:00+08 |
| R11 | `超时\|30 ?秒\|秒` | `docs/specs/SPEC-06-security-gatekeeping.md` | 2 行（`:104` 带 300s／`:151` 不带数） | 10:02+08 |
| R12 | `^\|?\s*\*?\*?C18` | `docs/PLAN.md` | 3 行（`:1368` 契约表／`:1444` 切片映射／`:3210` 差异表） | 09:57+08 |
| R13 | `statemachine\.` 与 `SetState\(\|Dispatch\(` | `cmd/wisp` | 生产落点：`models.go:303`（FirstRun Machine）、`resident_approval_windows.go:823`／`:900`（`b.SetState(…)`＝球的渲染态）、`approval_always.go:190`／`approval_reply.go:168`（只回状态名）⇒ 未见 Machine 被推进 `AwaitingApproval`（§7 G-2） | 10:00+08 |
| R14 | `Queue\(\)\.Timeout\(\)\|ApprovalTimeout\|approval\.NewQueue` | `cmd/wisp` | 23 行＝产码 3（`run.go:616`、`resident_approval_windows.go:372`、`:383` 日志）＋注释/日志 3（`config_reload.go:42`、`:445`、`:456`）＋测试 17 | 10:01+08 |

**未跑清单（本腿一票未跑，全数具名）**：`go build ./...`／`go test ./...`／`go vet`／`gofumpt`／`wisp slo`／`sh scripts/*`／任何 `go run`。原因＝派单硬闸①（编排者正在跑 `cmd/wisp` 整包取红名册，并发跑会洗掉它的读数）＋硬规矩①。⇒ **本文件全部判语＝静态读码**，没有任何一枚"今天红/绿"的断言出自本腿；哪枚测试今天真的红，只有编排者那一发的名册能答。
**门禁读数（自陈）**：写本节时本腿工具调用约 56 次、**终值约 74 次**（§8 的自我对抗与就地改读数占了尾程；`go` 命令**始终 0 次**）；`git add` 只带过本文件一枚 pathspec（骨架 `4e877958` 09:53:45＋0800／§1＋§2 `f9bd0eb7` 10:01:38／§3＋§4 `e1491088` ≈10:04／§5＋§6＋§7 `f3a8fab6` 10:08:16／本次 §8＋§9 与就地更正另落一笔，终值 `wc` 与时刻见 §9 末行与 `git log`）；未改任何他人已跟踪文件（本文件是本腿自建自填的交付件）；未碰 `.scratch/wisp/issues/**` 的 AC 框、未碰 `docs/reports/pending-and-issues.md`。

## §7 判不动／量不到

| 号 | 格子 | 本腿能判到哪 | 为什么判不动／量不到 |
|---|---|---|---|
| G-1 | **审批卡／托盘／悬浮球／面板页面上的字**（"5 分钟"这类人话最可能的落点） | Go 侧 R5／R6／R7 三把尺都是零 ⇒ 屏幕上的秒数若存在，其来源要么现算、要么在那两棵树里 | 硬规矩②：`frontend/**` 与 `design/**` 不许读、结论里不许转述 ⇒ **这一面属页面侧、本编队量不到，需机主带去他用的那枚 agent 另查** |
| G-2 | 生产里 `AwaitingApproval` 那枚 **300s 机器定时器到底有没有被 arm** | 球不自带计时器（`ball_windows.go:309-311`＋`ball/doc.go:19`）；`EvApprovalTimeout` 零枚手工 Dispatch（R1 尺与 `EvApprovalTimeout` 尺都只命中 `timeouts.go:51`／`table.go:172`／测试） | 本腿未读到"常驻/run 装配根把 Machine 推进 `AwaitingApproval` 并 rearm"的那一行；`resident_ball_windows.go:275` 那个 `Initial: statemachine.StateSleeping` 属于谁（ball.Options 还是 statemachine.Options）未复现 |
| G-3 | 面板快照是否用**别种形状**（手写 map／字符串拼接）把 timeout 换算成毫秒送进 JS | R9 零命中＋`panel_pump.go:408` 只打 `len(snap.Pending)`＋`pump.go:50-53` 注释自陈 `approval.remainingMs` 在 no-source 名单 | 未逐行读 `internal/panel` 的快照构造体；"没有 json 标签"不等于"没有手写键名"，本腿只到"两把尺同向"，不到"证否" |
| G-4 | `pump.go:53` 引的 `gate.go:274`／`:387`／`:510-513`（"Remaining 装的是静态窗口长度"） | 只复核到 `gate.go:562` `Remaining: g.q.WarningLead()` 与 `:308`／`:437` `Remaining: g.window` 三枚同族落点确实存在 | 那三枚行号是他票注释的读数，本腿未逐行 `Read`；行号会漂（派单硬规矩⑤），所以引而不认 |
| G-5 | `batch_test.go:106`、`boundary_test.go:142` 的**断言条件** | 只现量到落点与解释语／赋值行 | 未读到紧邻上文的 `if` ⇒ 5.1 A 类这两枚的"红条件"是形状推断；机主若拿它当凭据，请让实现腿补读 |
| G-6 | 这些钉子**今天**红不红 | — | 本腿禁跑 go（硬规矩①），全节不提红绿；只有编排者那一发 `cmd/wisp` 整包名册可答 |
| G-7 | **`Q-77` 本体**：C18 写死的 300s 与可配的 `confirm_timeout_sec` 谁优先 | 本腿只交名册：条件性说谎文案 1 枚（§1 #9）＋行为面写死 300 共 2 枚（§2.3）＋日志字面 2 枚（§3 L3／L5）＋冻结文字 7 枚（§4）＋**生效红 0 枚**（§5 判语） | 硬契约面（`SPEC-12 §4.1`＋票 267 AC#3 禁区）⇒ ⛔ 不裁、不选甲乙；**"配了要不要真生效"这一格归机主** |
| G-8 | 托盘 tooltip 的字 | Go 侧未找到承载位（R5 只 3 行） | 与 G-1 同格：若存在则在页面侧／原生文案资源里，本编队不读 |
| G-9 | 首启落盘那一刻屏上写的值（盘上现在真是 300 吗） | `schema.go:460` 的 `default:"300"` 是反射源 | 「盘上＝300」是同票 `267-a2` 的现量（其 §0 判 `firstrun.go:82` 走 `NewDefaults()`），本腿未自跑那把尺，只引注 |

## §8 我写错的读数（自我对抗，逐条真改）

| # | 我原本写成 | 错在哪 | 已就地改成（现量凭据） |
|---|---|---|---|
| E1 | §4 判语「**冻结件里 5 枚把 300s 写成数字（D1／D2／D3／D4／D6／D7／D8 共 7 枚…**」 | 同一个句子里既写 5 又列 7，把"带数字的枚数"和"契约形状的枚数"混成一锅——这是**会被机主当成数错**的那种句子 | 改成「共 **7** 枚带数字，其中 **2** 枚默认值形状（D3／D6）、**5** 枚契约值形状（D1／D2／D4／D7／D8）」，单位与集合都列名 |
| E2 | §1 #9「本腿判出的**唯一**一枚"配而不生效就会说谎"的现成文案」 | "唯一"没带方向，读者会以为是全案唯一说谎面——而 §3 L3／L5 与 §2.3 明明是**反方向**（判乙）会说谎的落点 | 改成「**判甲方向**唯一一枚……；反方向（判乙）会说谎的是 §3 L3、§3 L5／§2.3」，并同步把 §1 表后那句小结改成「另 1 枚**条件性**说谎（#9，方向＝判甲）」 |
| E3 | §2.1 那条注释名册里我写的 `internal/agent/approval/{doc.go:6,21,**82,59**; …; queue.go:59,82; **panel/composer_handlers.go:23**}` | 两处错置：`82`／`59` 属 `queue.go` 不属 `doc.go`；`composer_handlers.go:23` 在 `internal/panel` 不在 `internal/agent/approval`——**我把两把尺的输出串行错了包**，这种错会让人点进去找不到行 | 改成 `{doc.go:6,21; gate.go:30,499; approval.go:21; queue.go:59,82}` ＋ 具名 `internal/panel/composer_handlers.go:23`（两处落点都已在 09:55+08 那次 `internal` 尺输出里逐行复认） |
| E4 | §5.3 我先前写「★90 这一枚已经有测试在守（`256_test:293/:325` 种 90 并要求 `Queue().Timeout() == 90s`）」 | **漏了中间那一枚，方向就反了**：`:293` 只是把文件改成 90，真正守着 90 的是 `:324-325` 那枚"**重建**后的新门"；而 `:310` 守的是"**已建好的门仍是 45s**"。我那一句会让机主读成"改了文件当场生效"，恰好把 `Q-77` 最贵的那一维（生效时机）抹平 | 已重写：四步链（`:286` 建 45 → `:293` 改 90 → `:310` 旧门仍 45＝构造期取值／不热生效 → `:324-325` 新门读 90）＋把 `:311-314` 那句逐字承认（加了 re-apply 就红、256 的自陈限制须重裁、连带 255/265 的 ⓑ 文案）抄回，并在 §9 单列第三维 |
| E5 | §6 尺表 R14 我第一版草稿记「26 行」 | 手快把上一次 grep 的 head_limit 当行数；真数是**逐行点数 23**＝产码 3＋注释/日志 3＋测试 17 | 表中已写 23 并给出三段拆法（`cmd/wisp` 那次 `Queue\(\)\.Timeout\(\)\|ApprovalTimeout\|approval\.NewQueue` 尺输出可逐行核） |
| E6 | §3 L3 我最初的判语是「回落日志说的是编译常量＝**实话，无风险**」 | 读完 `resident_approval_windows.go:447-462` 才发现 267 的值域门会让"用户写了带外值"也走 `:453` 的 err 分支——那句话**对结果为真、对原因为残缺**（它不告诉用户"你写的那个数被拒了"） | 已改成两枚风险（① 该字面串不在任何常量引用链上，改 `queue.go:107` 它不会红；② 带外值被洗成"编译常量"），并顺带把 L4 那枚"能把配的值与门拿到的值并排打出来"的日志单独具名给机主当核验抓手 |
| E7 | §1 #2 我最初只写"与 timeout 无关，永远说 30" | 说得不够：`queue.go:88-90` 的兜底是 `warnBefore <= 0 \|\| warnBefore >= timeout → DefaultApprovalWarning(30s)`，所以当 timeout 落在带内最小 31s 时，提示**只剩 1s 提前量**——这正是 `validate.go:128` 定 31、`validate_267_test.go:63-80` 钉 `lead+1` 的理由 | §1 #2 那句"永远讲 30"保留（它对屏幕文案为真），但把 31／`lead+1` 这层因果补在 §5.3 的带内行与 §9 的第三条里，别让读者以为"提示提前量"也是可配数 |

## §9 交件判语

**这份名册摆给机主的三堆数（都不是裁定）：**

1. **会被这把配置变成假话的用户可见文案＝0 枚（Go 侧）。** 两枚会打出秒数的句子（`queue.go:479` 拒绝句、`gate.go:563` 提示句）都是 `fmt.Sprintf` 现算，数字来源就是队列自己的 `timeout`；配 90 之后它们说 90，说"5 分钟"的那种字面在 Go 生产面**根本不存在**（R2／R3／R5／R6／R7 五把尺同向）。**真正的谎面藏在非文案处**：判甲方向 1 枚（`unwired.go:121` 的 `consumed:` 断言），判乙方向 3 枚（`resident_approval_windows.go:456` 的日志字面 `DefaultApprovalTimeout=300s`、`events.go:42` 的事件名 `approval.timeout-300s`、`timeouts.go:51` 状态机表里那枚 300s 行为值）。另有**沉默面 3 枚**（`run.go:1364`／`resident_approval_windows.go:829`／面板快照零栏）——改了配置之后那里既不会写假的 300，也不会写真的 90。
2. **会因"配置真生效"而红的测试钉＝0 枚。** 派单担心的那类（文案含 300 的 needle）经两把尺确认为**不存在**。会红的另外两族与"生效"无关：改那个**字面 300**（A 类 6 组，其中 `ticket84_no_owner_test.go:232-257` 是 300s 墙钟、最贵）；以及 267 的**带门**挡住带外种子（具名 `run_mode101_test.go:110`＝1、`panel_pump_test.go:136`＝2、`approval_reply_201_test.go:105-110`＝1 下限；另 4 枚 20 是同票 267-a2 的读数，本腿只引注不重跑）。
3. **`Q-77` 还有第三维，本腿现量到并且它已经有钉子**：`resident_approval_risk_256_windows_test.go:310` 逐字钉住「**已建好的门保持构造期那个值，改了 config.toml 不会热生效**」。所以"配了要不要真生效"至少分成 **(甲) 要不要吃这个数**与 **(乙) 什么时候吃（重启／重建 vs 热生效）**——今天代码是「吃、但只在构造时吃」。机主只答前者会漏掉后者；而 `:311-314` 已经预告：一旦加了 re-apply 路径，这枚钉子会红，256 的自陈限制与 255／265 那句 ⓑ 文案都得重裁。

**本腿量不到的那一格（具名，供机主转给他用的那枚 agent）**：**审批卡／托盘／悬浮球／面板页面上的字**——`frontend/**` 与 `design/**` 两棵树不属本编队、本腿一律未读（硬规矩②）。Go 侧 R5 只 3 行、R6 生产可见串零枚字面 300 ⇒ 页面若真写着"5 分钟"，它要么现算、要么就在那里写死；**"屏幕上会不会出现假的 300"这一问，只有页面侧能答**。其余量不到／判不动见 §7 G-2（生产里那枚机器定时器有没有被 arm）、G-3（快照有无手写键）、G-4／G-5（他票行号与两枚断言体未补读）、G-6（今天红绿——本腿零枚 go，全部判语静态读码）。

**合规自陈**：⛔ 零枚 `go`／`wisp slo`／`sh scripts/*`；搜索根只 `cmd internal tools docs .scratch scripts`；冻结件（`docs/PLAN.md`、`docs/specs/**`）只读、一字未改，引用处均标〔冻结〕；未碰 issues AC 框与台账；`git add` 全程只带本文件一枚 pathspec；未用 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`restore`；`Q-77` 不裁、不选甲乙。

**体量与时刻（现量尺，`wc -l -c` ＋ `date`，`10:10:53+0800`）**：本文件写到这一行为止＝**207 行／43,018 字节**，最后一笔已落 commit＝`f3a8fab6`（10:08:16+0800，§5＋§6＋§7）。下面这一行本身会让终值再长 2 行左右，故交件终值以 §8／§9 那笔 commit 之后编排者复量为准（本腿在同一批把它报回去）。起手锚 `c9600334`、骨架 `4e877958`、§1＋§2 `f9bd0eb7`、§3＋§4 `e1491088`、§5＋§6＋§7 `f3a8fab6`——**五笔，全部只带 `.scratch/wisp/probes/267/a3/census.md` 这一枚 pathspec**。
