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
| 9 | `internal/config/unwired.go:121` | `consumed: cmd/wisp/run.go and cmd/wisp/resident_approval_windows.go both build the approval timeout from it (the resident leg since ticket 256); ticket 267 bands it [31, 3600] at load` | 无秒数，但**断言「consumed」** | ★**本腿判出的唯一一枚"配而不生效就会说谎"的现成文案**：它今天声称两枚消费点都从它建超时。`Q-77` 若判甲（配了不真生效），说谎的就是这一行；判乙则它为真。名册归本腿，**裁定归机主** |
| 10 | `cmd/wisp/config_reload.go:224`／`:229` | 放宽未确认／已确认两条款回执 | 无秒数 | 不说谎 |
| 11 | `internal/panel/pump.go:50-53`（注释自陈） | 逐字：`the no-source list is seven: approval.remainingMs joins it, because the L1 countdown has no producer either - approval.EventTick is declared and emitted nowhere in the tree, the Remaining that is emitted carries the static window length` | 面板快照**没有**审批剩余秒数这一栏 | **不说谎＝也是沉默面**：快照里根本没有 `remainingMs`，所以配 90 也不会有一栏写 300。本腿用第二把尺独立复认：`json:"(timeout\|window\|deadline\|remaining\|expire\|warn)[^"]*"` 在 `{cmd,internal}` **零命中** ⇒ 两把尺同向。（gate.go:274/:387/:510-513 那三枚行号是**该注释引用他票的读数**，本腿未复测，见 §7 G-4） |

**用户可见文案的说谎数＝0 枚（Go 侧），另 1 枚条件性说谎（#9）。** 但「不说谎」不等于「说全」：#3／#5 两片把「卡片还会等多久」这件事**完全不打印**，机主把 300 改成 90 之后，屏幕上既不会有假的 300，也不会有真的 90——这一格是**沉默**，不是谎。

⚠ **审批卡／托盘／悬浮球／面板 HTML 上的字（「5 分钟」这类人话最可能落点）属页面侧：`frontend/**` 与 `design/**` 两棵树本编队不许读，故"页面上有没有写死的 300 秒／5 分钟"这一格本腿量不到，需机主带去他用的那枚 agent 另查。** 派单点名的四类可见面里，球与托盘的**文字渲染**同样落在那两棵树（Go 侧只 `SetState` 一个状态名，`ball_windows.go:309-311` 注释逐字「the consumer drives this after its statemachine dispatch」＝球不自带计时器）。

## §2 常量与默认值的落点

### 2.1 `DefaultApprovalTimeout` 定义与引用者（现量尺＝`DefaultApprovalTimeout`，根＝`{cmd,internal,tools,scripts}`）

- 定义：`internal/agent/approval/queue.go:107` 逐字 `DefaultApprovalTimeout = 300 * time.Second`（同段 `:106` 注释逐字 `// DefaultApprovalTimeout: 「超时 300s 一律判拒绝」. Never an infinite wait.`，注释豁免 d22scan 的 emoji 尺，但它是**人读的 300**）。
- **生产引用者＝1 枚**：`internal/agent/approval/queue.go:86`（`NewQueue` 的 `timeout <= 0` 兜底；兜底**不看 config**，只看调用方给没给值）。
- 生产侧"非引用者"的 300 落点（**字面串**，不是常量）：`cmd/wisp/resident_approval_windows.go:456`（日志，见 §3）、`:445`／`:296`／`:308`／`:309`、`cmd/wisp/resident_windows.go:129`、`cmd/wisp/config_reload.go:50`、`internal/config/manager.go:26`、`internal/tools/gate.go:86`、`internal/perm/store.go:86`、`internal/risk/mode.go:32`、`internal/agent/approval/{doc.go:6,21,82,59; gate.go:30,499; approval.go:21; queue.go:59,82; panel/composer_handlers.go:23}`＝**注释与包文档**，不进屏幕。
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

## §4 文档（含冻结件标注）

## §5 会因配置生效而红的测试钉

## §6 尺读数与未跑清单（门禁读数）

## §7 判不动／量不到

## §8 我写错的读数（自我对抗）

## §9 交件判语
